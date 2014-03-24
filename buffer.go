package cl11

import (
	"fmt"
	"reflect"
	"unsafe"

	clw "github.com/rdwilliamson/clw11"
)

// A buffer object stores a one-dimensional collection of elements. Elements of
// a buffer object can be a scalar data type (such as an int, float), vector
// data type, or a user-defined structure.
type Buffer struct {
	id clw.Mem

	// The context the buffer was created on.
	Context *Context

	// The size of the buffer in bytes.
	Size int64

	// Usage information for the buffer from the device's point of view.
	Flags MemFlags

	// The host backed memory for the buffer, if applicable.
	Host interface{}

	// The parent buffer, if applicable.
	Buffer *Buffer

	// The offset in bytes within the parent buffer.
	Origin int64
}

// A source and destination rectangular region.
type Rect struct {

	// The source offset. For a 2D rectangle SrcOrigin[2] should be 0. The
	// offset in bytes is SrcOrigin[2] * SrcSlicePitch + SrcOrigin[1] *
	// SrcRowPitch + SrcOrigin[0].
	SrcOrigin     [3]int64
	SrcRowPitch   int64
	SrcSlicePitch int64

	// The destination offset. For a 2D rectangle DstOrigin[2] should be 0. The
	// offset in bytes is DstOrigin[2] * DstSlicePitch + DstOrigin[1] *
	// DstRowPitch + DstOrigin[0].
	DstOrigin     [3]int64
	DstRowPitch   int64
	DstSlicePitch int64

	// The (width, height, depth) in bytes of the 2D or 3D rectangle being
	// copied. For a 2D rectangle the depth value given by Region[2] should be
	// 1.
	Region [3]int64
}

type MappedBuffer struct {
	pointer unsafe.Pointer
	size    int64
	buffer  *Buffer
}

type MemFlags int

// Bitfield.
const (
	MemReadWrite = MemFlags(clw.MemReadWrite)
	MemWriteOnly = MemFlags(clw.MemWriteOnly)
	MemReadOnly  = MemFlags(clw.MemReadOnly)
)

type MapFlags int

// Bitfield.
const (
	MapRead  = MapFlags(clw.MapRead)
	MapWrite = MapFlags(clw.MapWrite)
)

type BlockingCall clw.Bool

const (
	Blocking    = BlockingCall(clw.True)
	NonBlocking = BlockingCall(clw.False)
)

func (c *Context) CreateDeviceBuffer(size int64, mf MemFlags) (*Buffer, error) {

	memory, err := clw.CreateBuffer(c.id, clw.MemFlags(mf), clw.Size(size), nil)
	if err != nil {
		return nil, err
	}

	return &Buffer{id: memory, Context: c, Size: size, Flags: mf}, nil
}

func (c *Context) CreateDeviceBufferInitializedBy(mf MemFlags, value interface{}) (*Buffer, error) {

	flags := clw.MemFlags(mf) | clw.MemCopyHostPointer

	var scratch [scratchSize]byte
	pointer, size := getPointerAndSize(value, unsafe.Pointer(&scratch[0]))

	memory, err := clw.CreateBuffer(c.id, flags, clw.Size(size), pointer)
	if err != nil {
		return nil, err
	}

	return &Buffer{id: memory, Context: c, Size: int64(size), Flags: mf}, nil
}

func (c *Context) CreateDeviceBufferFromHostMem(mf MemFlags, host interface{}) (*Buffer, error) {

	flags := clw.MemFlags(mf) | clw.MemUseHostPointer

	value := reflect.ValueOf(host)
	if kind := value.Kind(); kind != reflect.Ptr && kind != reflect.Slice {
		return nil, wrapError(fmt.Errorf("host value not addressable"))
	} else if kind == reflect.Ptr {
		for {
			value = value.Elem()
			if value.Kind() != reflect.Ptr {
				break
			}
		}
	}
	pointer, size := addressablePointerAndSize(value)

	memory, err := clw.CreateBuffer(c.id, flags, clw.Size(size), pointer)
	if err != nil {
		return nil, err
	}

	return &Buffer{id: memory, Context: c, Size: int64(size), Host: host, Flags: mf}, nil
}

func (c *Context) CreateHostBuffer(size int64, mf MemFlags) (*Buffer, error) {

	flags := clw.MemFlags(mf) | clw.MemAllocHostPointer

	memory, err := clw.CreateBuffer(c.id, flags, clw.Size(size), nil)
	if err != nil {
		return nil, err
	}

	return &Buffer{id: memory, Context: c, Size: size, Flags: mf}, nil
}

func (c *Context) CreateHostBufferInitializedBy(mf MemFlags, value interface{}) (*Buffer, error) {

	flags := clw.MemFlags(mf) | clw.MemAllocHostPointer | clw.MemCopyHostPointer

	var scratch [scratchSize]byte
	pointer, size := getPointerAndSize(value, unsafe.Pointer(&scratch[0]))

	memory, err := clw.CreateBuffer(c.id, flags, clw.Size(size), pointer)
	if err != nil {
		return nil, err
	}

	return &Buffer{id: memory, Context: c, Size: int64(size), Flags: mf}, nil
}

func (b *Buffer) CreateSubBuffer(mf MemFlags, origin, size int64) (*Buffer, error) {

	region := clw.BufferRegion{clw.Size(origin), clw.Size(size)}

	memory, err := clw.CreateSubBuffer(b.id, clw.MemFlags(mf), region)
	if err != nil {
		return nil, err
	}

	return &Buffer{id: memory, Context: b.Context, Size: size, Flags: mf, Buffer: b, Origin: origin}, nil
}

func (b *Buffer) Retain() error {
	return clw.RetainMemObject(b.id)
}

func (b *Buffer) Release() error {
	return clw.ReleaseMemObject(b.id)
}

func (cq *CommandQueue) CopyBuffer(src, dst *Buffer, srcOffset, dstOffset, size int64, waitList []*Event,
	e *Event) error {

	var event *clw.Event
	if e != nil {
		event = &e.id
		e.Context = cq.Context
		e.CommandType = CommandCopyBuffer
		e.CommandQueue = cq
	}

	return clw.EnqueueCopyBuffer(cq.id, src.id, dst.id, clw.Size(srcOffset), clw.Size(dstOffset), clw.Size(size),
		cq.toEvents(waitList), event)
}

func (cq *CommandQueue) CopyBufferRect(src, dst *Buffer, r *Rect, waitList []*Event, e *Event) error {

	var event *clw.Event
	if e != nil {
		event = &e.id
		e.Context = cq.Context
		e.CommandType = CommandCopyBuffer
		e.CommandQueue = cq
	}

	return clw.EnqueueCopyBufferRect(cq.id, src.id, dst.id, r.srcOrigin(), r.dstOrigin(), r.region(), r.srcRowPitch(),
		r.srcSlicePitch(), r.dstRowPitch(), r.dstSlicePitch(), cq.toEvents(waitList), event)
}

func (cq *CommandQueue) MapBuffer(b *Buffer, bc BlockingCall, flags MapFlags, offset, size int64, waitList []*Event,
	e *Event) (*MappedBuffer, error) {

	var event *clw.Event
	if e != nil {
		event = &e.id
		e.Context = cq.Context
		e.CommandType = CommandMapBuffer
		e.CommandQueue = cq
	}

	pointer, err := clw.EnqueueMapBuffer(cq.id, b.id, clw.Bool(bc), clw.MapFlags(flags), clw.Size(offset),
		clw.Size(size), cq.toEvents(waitList), event)
	if err != nil {
		return nil, err
	}

	return &MappedBuffer{pointer, size, b}, nil
}

func (cq *CommandQueue) UnmapBuffer(mb *MappedBuffer, waitList []*Event, e *Event) error {

	var event *clw.Event
	if e != nil {
		event = &e.id
		e.Context = cq.Context
		e.CommandType = CommandUnmapMemoryObject
		e.CommandQueue = cq
	}

	return clw.EnqueueUnmapMemObject(cq.id, mb.buffer.id, mb.pointer, cq.toEvents(waitList), event)
}

func (r *Rect) srcOrigin() [3]clw.Size {
	return [3]clw.Size{clw.Size(r.SrcOrigin[0]), clw.Size(r.SrcOrigin[1]), clw.Size(r.SrcOrigin[2])}
}

func (r *Rect) srcRowPitch() clw.Size {
	return clw.Size(r.SrcRowPitch)
}

func (r *Rect) srcSlicePitch() clw.Size {
	return clw.Size(r.SrcSlicePitch)
}

func (r *Rect) dstOrigin() [3]clw.Size {
	return [3]clw.Size{clw.Size(r.DstOrigin[0]), clw.Size(r.DstOrigin[1]), clw.Size(r.DstOrigin[2])}
}

func (r *Rect) dstRowPitch() clw.Size {
	return clw.Size(r.DstRowPitch)
}

func (r *Rect) dstSlicePitch() clw.Size {
	return clw.Size(r.DstSlicePitch)
}

func (r *Rect) region() [3]clw.Size {
	return [3]clw.Size{clw.Size(r.Region[0]), clw.Size(r.Region[1]), clw.Size(r.Region[2])}
}

func (bm *MappedBuffer) Float32Slice() []float32 {

	var header reflect.SliceHeader
	header.Data = uintptr(bm.pointer)
	size := int(bm.size / int64(float32Size))
	header.Len = size
	header.Cap = size

	return *(*[]float32)(unsafe.Pointer(&header))
}
