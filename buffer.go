package cl11

import clw "github.com/rdwilliamson/clw11"

type Buffer struct {
	id      clw.Memory
	Context *Context
	Size    int
	Host    []byte // The host backed memory for the buffer.
	Flags   MemoryFlags
}

type MemoryFlags uint8

const (
	MemoryReadWrite MemoryFlags = MemoryFlags(clw.MemoryReadOnly)
	MemoryWriteOnly MemoryFlags = MemoryFlags(clw.MemoryReadOnly)
	MemoryReadOnly  MemoryFlags = MemoryFlags(clw.MemoryReadOnly)
)

type MapFlags uint8

const (
	MapRead  MapFlags = MapFlags(clw.MapRead)
	MapWrite MapFlags = MapFlags(clw.MapWrite)
)

const (
	Blocking    = true
	NonBlocking = false
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

	memory, err := clw.CreateBuffer(c.id, flags, clw.Size(len(host)), host)
	if err != nil {
		return nil, err
	}

	return &Buffer{id: memory, Context: c, Size: len(host), Flags: mf}, nil
}

func (c *Context) CreateDeviceBufferOnHost(mf MemoryFlags, host []byte) (*Buffer, error) {

	flags := clw.MemoryFlags(mf) | clw.MemoryUseHostPointer

	memory, err := clw.CreateBuffer(c.id, flags, clw.Size(len(host)), host)
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

	memory, err := clw.CreateBuffer(c.id, flags, clw.Size(len(host)), host)
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
		e.CommandType = CommandCopyBuffer
	}

	return clw.EnqueueCopyBuffer(cq.id, src.id, dst.id, clw.Size(srcOffset), clw.Size(dstOffset), clw.Size(size),
		cq.toEvents(waitList), event)
}

func (cq *CommandQueue) MapBuffer(b *Buffer, blocking bool, flags MapFlags, offset, size int, waitList []*Event,
	e *Event) ([]byte, error) {

	var event *clw.Event
	if e != nil {
		event = &e.id
		e.CommandType = CommandMapBuffer
	}

	mapped, err := clw.EnqueueMapBuffer(cq.id, b.id, clw.ToBool(blocking), clw.MapFlags(flags), clw.Size(offset),
		clw.Size(size), cq.toEvents(waitList), event)
	if err != nil {
		return nil, err
	}
	return mapped, nil
}

func (cq *CommandQueue) UnmapBuffer(b *Buffer, mapped []byte, waitList []*Event, e *Event) error {

	var event *clw.Event
	if e != nil {
		event = &e.id
		e.CommandType = CommandUnmapMemoryObject
	}

	return clw.EnqueueUnmapMemObject(cq.id, b.id, mapped, cq.toEvents(waitList), event)
}
