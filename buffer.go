package cl11

import (
	"fmt"
	"reflect"
	"unsafe"

	clw "github.com/rdwilliamson/clw11"
)

// A buffer object stores a one-dimensional collection of elements (though there
// are "rectangle" operations). Elements of a buffer object can be a scalar data
// type (such as an int, float), vector data type, or a user-defined structure.
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

// Rectangle layout in bytes. For a 2D image SlicePitch and Origin[2] are zero.
type RectLayout struct {
	Origin     [3]int64
	RowPitch   int64
	SlicePitch int64
}

// A source and destination rectangular region.
type Rect struct {

	// The source offset in bytes. For a 2D rectangle Src.Origin[2] should be 0.
	// The offset in bytes is Src.Origin[2] * Src.SlicePitch + Src.Origin[1] *
	// Src.RowPitch + Src.Origin[0].
	Src RectLayout

	// The destination offset in bytes. For a 2D rectangle Dst.Origin[2] should
	// be 0. The offset in bytes is Dst.Origin[2] * Dst.SlicePitch +
	// Dst.Origin[1] * Dst.RowPitch + Dst.Origin[0].
	Dst RectLayout

	// The (width, height, depth) in bytes of the 2D or 3D rectangle being
	// copied. For a 2D rectangle the depth value given by Region[2] should be
	// 1.
	Region [3]int64
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

// Creates a buffer object on the device.
//
// Creates an uninitialized buffer on the device. The size is in bytes.
func (c *Context) CreateDeviceBuffer(size int64, mf MemFlags) (*Buffer, error) {

	memory, err := clw.CreateBuffer(c.id, clw.MemFlags(mf), clw.Size(size), nil)
	if err != nil {
		return nil, err
	}

	return &Buffer{id: memory, Context: c, Size: size, Flags: mf}, nil
}

// Create a buffer object on the device.
//
// Create a buffer object initialized with the passed value. The size and
// contents are determined by the passed value.
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

// Create a buffer object for the device backed by host memory.
//
// Create a device accessible buffer object on the host. The host value must be
// addressable. The OpenCL implementation is allowed to cache the data on the
// device.
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

// Creates a buffer object that is host accessible.
//
// Creates an uninitialized buffer in pinned memory. The size is in bytes. This
// memory is not pageable and allows for DMA copies (which are faster).
func (c *Context) CreateHostBuffer(size int64, mf MemFlags) (*Buffer, error) {

	flags := clw.MemFlags(mf) | clw.MemAllocHostPointer

	memory, err := clw.CreateBuffer(c.id, flags, clw.Size(size), nil)
	if err != nil {
		return nil, err
	}

	return &Buffer{id: memory, Context: c, Size: size, Flags: mf}, nil
}

// Creates a buffer object that is host accessible.
//
// Creates an initialized buffer in pinned memory. The size and contents are
// determined by value. This memory is not pageable and allows for DMA copies
// (which are faster).
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

// Creates a buffer object from an existing object with the passed offset and
// size in bytes.
//
// If the size and offset are the same as another sub buffer the implementation
// may return the same sub buffer and increment the reference count.
func (b *Buffer) CreateSubBuffer(mf MemFlags, origin, size int64) (*Buffer, error) {

	region := clw.BufferRegion{clw.Size(origin), clw.Size(size)}

	memory, err := clw.CreateSubBuffer(b.id, clw.MemFlags(mf), region)
	if err != nil {
		return nil, err
	}

	return &Buffer{id: memory, Context: b.Context, Size: size, Flags: mf, Buffer: b, Origin: origin}, nil
}

// Increments the memory object reference count.
//
// The OpenCL commands that return a buffer perform an implicit retain.
func (b *Buffer) Retain() error {
	return clw.RetainMemObject(b.id)
}

// Decrements the memory object reference count.
//
// After the buffers reference count becomes zero and commands queued for
// execution that use the buffer have finished the buffer is deleted.
func (b *Buffer) Release() error {
	return clw.ReleaseMemObject(b.id)
}

// Return the buffer's reference count.
//
// The reference count returned should be considered immediately stale. It is
// unsuitable for general use in applications. This feature is provided for
// identifying memory leaks.
func (b *Buffer) ReferenceCount() (int, error) {

	var count clw.Uint
	err := clw.GetMemObjectInfo(b.id, clw.MemReferenceCount, clw.Size(unsafe.Sizeof(count)), unsafe.Pointer(&count),
		nil)
	if err != nil {
		return 0, err
	}

	return int(count), nil
}

// Enqueues a command to copy from one buffer object to another.
//
// Source offset, destination offset, and size are in bytes.
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

// Enqueues a command to copy a rectangular region from the buffer object to
// another buffer object.
//
// See Rect definition for how source and destination are defined.
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

// Enqueues a command to map a region of the buffer object given by buffer into
// the host address space.
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

	return &MappedBuffer{pointer, size, b.id}, nil
}

// Enqueue commands to read from a buffer object to host memory.
//
// Offset is in bytes. The destination must be addressable. If the buffer object
// is backed by host memory then all commands that use it and sub buffers must
// have finished execution and it must not be mapped otherwise the results are
// undefined.
func (cq *CommandQueue) ReadBuffer(b *Buffer, bc BlockingCall, offset int64, dst interface{}, waitList []*Event,
	e *Event) error {

	var event *clw.Event
	if e != nil {
		event = &e.id
		e.Context = cq.Context
		e.CommandType = CommandReadBuffer
		e.CommandQueue = cq
	}

	pointer, size, err := tryPointerAndSize(dst)
	if err != nil {
		return wrapError(err)
	}

	err = clw.EnqueueReadBuffer(cq.id, b.id, clw.Bool(bc), clw.Size(offset), clw.Size(size), pointer,
		cq.toEvents(waitList), event)
	if err != nil {
		return err
	}

	return nil
}

// Enqueue commands to write to a buffer object from host memory.
//
// Offset is in bytes. If the buffer object is backed by host memory then all
// commands that use it and sub buffers must have finished execution and it must
// not be mapped otherwise the results are undefined. If the source is
// addressable a reference will be held to prevent it being garbage collected,
// it is the user's responsibility to ensure that the source data is valid at
// the time the write is actually performed. If the source isn't addressable a
// copy will be created.
func (cq *CommandQueue) WriteBuffer(b *Buffer, bc BlockingCall, offset int64, src interface{}, waitList []*Event,
	e *Event) error {

	// Ensure we always have an event if not blocking. The event will be used to
	// register a callback. Thus the source data is guaranteed to be referenced
	// somewhere to preventing it from being garbage collected. Once the event
	// has completed and the callback is triggered (doing nothing) the reference
	// to the source data will be removed allowing it to be garbage collected.
	if e == nil && bc == NonBlocking {
		e = &Event{}
	}

	var event *clw.Event
	if e != nil {
		event = &e.id
		e.Context = cq.Context
		e.CommandType = CommandWriteBuffer
		e.CommandQueue = cq
	}

	pointer, size, err := tryPointerAndSize(src)

	if err != nil {
		// The source value is not addressable so create a local copy of it.
		var scratch [scratchSize]byte
		pointer, size = getPointerAndSize(src, unsafe.Pointer(&scratch[0]))
		src = &scratch
	}

	err = clw.EnqueueWriteBuffer(cq.id, b.id, clw.Bool(bc), clw.Size(offset), clw.Size(size), pointer,
		cq.toEvents(waitList), event)
	if err != nil {
		return err
	}

	// Set a no-op callback, just need hold a reference to the source data.
	if bc == NonBlocking {
		err = e.SetCallback(noOpEventCallback, src)
		if err != nil {
			return err
		}
	}

	return nil
}

// Enqueue commands to read from a rectangular region from a buffer object to
// host memory.
//
// See Rect definition for offset is defined. The destination must be
// addressable. If the buffer object is backed by host memory then all commands
// that use it and sub buffers must have finished execution and it must not be
// mapped otherwise the results are undefined.
func (cq *CommandQueue) ReadBufferRect(b *Buffer, bc BlockingCall, offset *Rect, dst interface{}, waitList []*Event,
	e *Event) error {

	var event *clw.Event
	if e != nil {
		event = &e.id
		e.Context = cq.Context
		e.CommandType = CommandReadBuffer
		e.CommandQueue = cq
	}

	pointer, _, err := tryPointerAndSize(dst)
	if err != nil {
		return wrapError(err)
	}

	err = clw.EnqueueReadBufferRect(cq.id, b.id, clw.Bool(bc), offset.dstOrigin(), offset.srcOrigin(), offset.region(),
		offset.dstRowPitch(), offset.dstSlicePitch(), offset.srcRowPitch(), offset.dstRowPitch(), pointer,
		cq.toEvents(waitList), event)
	if err != nil {
		return err
	}

	return nil
}

// Enqueue commands to write a rectangular region to a buffer object from host
// memory.
//
// See Rect definition for offset is defined. If the buffer object is backed by
// host memory then all commands that use it and sub buffers must have finished
// execution and it must not be mapped otherwise the results are undefined. If
// the source is addressable a reference will be held to prevent it being
// garbage collected, it is the user's responsibility to ensure that the source
// data is valid at the time the write is actually performed. If the source
// isn't addressable a copy will be created.
func (cq *CommandQueue) WriteBufferRect(b *Buffer, bc BlockingCall, offset *Rect, src interface{}, waitList []*Event,
	e *Event) error {

	// Ensure we always have an event if not blocking. The event will be used to
	// register a callback. Thus the source data is guaranteed to be referenced
	// somewhere to preventing it from being garbage collected. Once the event
	// has completed and the callback is triggered (doing nothing) the reference
	// to the source data will be removed allowing it to be garbage collected.
	if e == nil && bc == NonBlocking {
		e = &Event{}
	}

	var event *clw.Event
	if e != nil {
		event = &e.id
		e.Context = cq.Context
		e.CommandType = CommandWriteBuffer
		e.CommandQueue = cq
	}

	pointer, _, err := tryPointerAndSize(src)

	if err != nil {
		// The source value is not addressable so create a local copy of it.
		var scratch [scratchSize]byte
		pointer, _ = getPointerAndSize(src, unsafe.Pointer(&scratch[0]))
		src = &scratch
	}

	err = clw.EnqueueWriteBufferRect(cq.id, b.id, clw.Bool(bc), offset.dstOrigin(), offset.srcOrigin(), offset.region(),
		offset.dstRowPitch(), offset.dstSlicePitch(), offset.srcRowPitch(), offset.srcSlicePitch(), pointer,
		cq.toEvents(waitList), event)
	if err != nil {
		return err
	}

	// Set a no-op callback, just need hold a reference to the source data.
	if bc == NonBlocking {
		err = e.SetCallback(noOpEventCallback, src)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Rect) srcOrigin() [3]clw.Size {
	return [3]clw.Size{clw.Size(r.Src.Origin[0]), clw.Size(r.Src.Origin[1]), clw.Size(r.Src.Origin[2])}
}

func (r *Rect) srcRowPitch() clw.Size {
	return clw.Size(r.Src.RowPitch)
}

func (r *Rect) srcSlicePitch() clw.Size {
	return clw.Size(r.Src.SlicePitch)
}

func (r *Rect) dstOrigin() [3]clw.Size {
	return [3]clw.Size{clw.Size(r.Dst.Origin[0]), clw.Size(r.Dst.Origin[1]), clw.Size(r.Dst.Origin[2])}
}

func (r *Rect) dstRowPitch() clw.Size {
	return clw.Size(r.Dst.RowPitch)
}

func (r *Rect) dstSlicePitch() clw.Size {
	return clw.Size(r.Dst.SlicePitch)
}

func (r *Rect) region() [3]clw.Size {
	return [3]clw.Size{clw.Size(r.Region[0]), clw.Size(r.Region[1]), clw.Size(r.Region[2])}
}
