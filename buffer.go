package cl11

import (
	"fmt"
	"reflect"
	"unsafe"

	clw "github.com/rdwilliamson/clw11"
)

type Buffer struct {
	id      clw.Memory
	Context *Context
	Size    int64
	Flags   MemoryFlags
	Host    interface{} // The host backed memory for the buffer.
}

type MappedBuffer struct {
	pointer unsafe.Pointer
	size    int64
	buffer  *Buffer
}

type MemoryFlags int

// Bitfield.
const (
	MemoryReadWrite = MemoryFlags(clw.MemoryReadWrite)
	MemoryWriteOnly = MemoryFlags(clw.MemoryWriteOnly)
	MemoryReadOnly  = MemoryFlags(clw.MemoryReadOnly)
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

func (c *Context) CreateDeviceBuffer(size int64, mf MemoryFlags) (*Buffer, error) {

	memory, err := clw.CreateBuffer(c.id, clw.MemoryFlags(mf), clw.Size(size), nil)
	if err != nil {
		return nil, err
	}

	return &Buffer{id: memory, Context: c, Size: size, Flags: mf}, nil
}

func (c *Context) CreateDeviceBufferInitializedBy(mf MemoryFlags, value interface{}) (*Buffer, error) {

	flags := clw.MemoryFlags(mf) | clw.MemoryCopyHostPointer

	var scratch [scratchSize]byte
	pointer, size := getPointerAndSize(value, unsafe.Pointer(&scratch[0]))

	memory, err := clw.CreateBuffer(c.id, flags, clw.Size(size), pointer)
	if err != nil {
		return nil, err
	}

	return &Buffer{id: memory, Context: c, Size: int64(size), Flags: mf}, nil
}

func (c *Context) CreateDeviceBufferFromHostMemory(mf MemoryFlags, host interface{}) (*Buffer, error) {

	flags := clw.MemoryFlags(mf) | clw.MemoryUseHostPointer

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

func (c *Context) CreateHostBuffer(size int64, mf MemoryFlags) (*Buffer, error) {

	flags := clw.MemoryFlags(mf) | clw.MemoryAllocHostPointer

	memory, err := clw.CreateBuffer(c.id, flags, clw.Size(size), nil)
	if err != nil {
		return nil, err
	}

	return &Buffer{id: memory, Context: c, Size: size, Flags: mf}, nil
}

func (c *Context) CreateHostBufferInitializedBy(mf MemoryFlags, value interface{}) (*Buffer, error) {

	flags := clw.MemoryFlags(mf) | clw.MemoryAllocHostPointer | clw.MemoryCopyHostPointer

	var scratch [scratchSize]byte
	pointer, size := getPointerAndSize(value, unsafe.Pointer(&scratch[0]))

	memory, err := clw.CreateBuffer(c.id, flags, clw.Size(size), pointer)
	if err != nil {
		return nil, err
	}

	return &Buffer{id: memory, Context: c, Size: int64(size), Flags: mf}, nil
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

func (bm *MappedBuffer) Float32Slice() []float32 {

	var header reflect.SliceHeader
	header.Data = uintptr(bm.pointer)
	size := int(bm.size / int64(float32Size))
	header.Len = size
	header.Cap = size

	return *(*[]float32)(unsafe.Pointer(&header))
}
