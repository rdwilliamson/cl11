package cl11

import (
	"unsafe"

	clw "github.com/rdwilliamson/clw11"
)

type Buffer struct {
	id      clw.Memory
	Context *Context
	Size    int
	Host    []byte // The host backed memory for the buffer.
	Flags   MemoryFlags
}

type MemoryFlags int

// Bitfield.
const (
	MemoryReadWrite MemoryFlags = MemoryFlags(clw.MemoryReadWrite)
	MemoryWriteOnly MemoryFlags = MemoryFlags(clw.MemoryWriteOnly)
	MemoryReadOnly  MemoryFlags = MemoryFlags(clw.MemoryReadOnly)
)

type MapFlags int

// Bitfield.
const (
	MapRead  MapFlags = MapFlags(clw.MapRead)
	MapWrite MapFlags = MapFlags(clw.MapWrite)
)

type BlockingCall clw.Bool

const (
	Blocking    = BlockingCall(clw.True)
	NonBlocking = BlockingCall(clw.False)
)

func (c *Context) CreateDeviceBuffer(size int, mf MemoryFlags) (*Buffer, error) {

	memory, err := clw.CreateBuffer(c.id, clw.MemoryFlags(mf), clw.Size(size), nil)
	if err != nil {
		return nil, err
	}

	return &Buffer{id: memory, Context: c, Size: size, Flags: mf}, nil
}

func (c *Context) CreateDeviceBufferFromHost(mf MemoryFlags, host []byte) (*Buffer, error) {

	flags := clw.MemoryFlags(mf) | clw.MemoryCopyHostPointer

	memory, err := clw.CreateBuffer(c.id, flags, clw.Size(len(host)), unsafe.Pointer(&host[0]))
	if err != nil {
		return nil, err
	}

	return &Buffer{id: memory, Context: c, Size: len(host), Flags: mf}, nil
}

func (c *Context) CreateDeviceBufferOnHost(mf MemoryFlags, host []byte) (*Buffer, error) {

	flags := clw.MemoryFlags(mf) | clw.MemoryUseHostPointer

	memory, err := clw.CreateBuffer(c.id, flags, clw.Size(len(host)), unsafe.Pointer(&host[0]))
	if err != nil {
		return nil, err
	}

	return &Buffer{id: memory, Context: c, Size: len(host), Host: host, Flags: mf}, nil
}

func (c *Context) CreateHostBuffer(size int, mf MemoryFlags) (*Buffer, error) {

	flags := clw.MemoryFlags(mf) | clw.MemoryAllocHostPointer

	memory, err := clw.CreateBuffer(c.id, flags, clw.Size(size), nil)
	if err != nil {
		return nil, err
	}

	return &Buffer{id: memory, Context: c, Size: size, Flags: mf}, nil
}

func (c *Context) CreateHostBufferFromHost(mf MemoryFlags, host []byte) (*Buffer, error) {

	flags := clw.MemoryFlags(mf) | clw.MemoryAllocHostPointer | clw.MemoryCopyHostPointer

	memory, err := clw.CreateBuffer(c.id, flags, clw.Size(len(host)), unsafe.Pointer(&host[0]))
	if err != nil {
		return nil, err
	}

	return &Buffer{id: memory, Context: c, Size: len(host), Flags: mf}, nil
}

func (b *Buffer) Release() error {
	return clw.ReleaseMemObject(b.id)
}

func (cq *CommandQueue) CopyBuffer(src, dst *Buffer, srcOffset, dstOffset, size int, waitList []*Event,
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

func (cq *CommandQueue) MapBuffer(b *Buffer, bc BlockingCall, flags MapFlags, offset, size int, waitList []*Event,
	e *Event) ([]byte, error) {

	var event *clw.Event
	if e != nil {
		event = &e.id
		e.Context = cq.Context
		e.CommandType = CommandMapBuffer
		e.CommandQueue = cq
	}

	mapped, err := clw.EnqueueMapBuffer(cq.id, b.id, clw.Bool(bc), clw.MapFlags(flags), clw.Size(offset),
		clw.Size(size), cq.toEvents(waitList), event)
	if err != nil {
		return nil, err
	}

	return toByteSlice(mapped, size), nil
}

func (cq *CommandQueue) UnmapBuffer(b *Buffer, mapped []byte, waitList []*Event, e *Event) error {

	var event *clw.Event
	if e != nil {
		event = &e.id
		e.Context = cq.Context
		e.CommandType = CommandUnmapMemoryObject
		e.CommandQueue = cq
	}

	return clw.EnqueueUnmapMemObject(cq.id, b.id, unsafe.Pointer(&mapped[0]), cq.toEvents(waitList), event)
}
