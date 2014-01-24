package cl11

import clw "github.com/rdwilliamson/clw11"

type Buffer struct {
	ID      clw.Memory
	Context *Context
	Size    int
	Host    []byte
	Flags   MemoryFlags
}

type MappedBuffer struct {
	b     *Buffer // Buffer from which data was mapped.
	data  []byte  // Entire range of mapped data.
	start int     // Where data will be read from, -1 if not mapped for reading.
	end   int     // Where data will be written to, -1 of not mapped for writing.
}

type MemoryFlags struct {
	Read  bool // Can a kernel read from this buffer.
	Write bool // Can a kernel write to this buffer.
}

type MapFlags struct {
	Read  bool
	Write bool
}

const (
	Blocking    = true
	NonBlocking = false
)

func (mf MemoryFlags) toBits() clw.MemoryFlags {
	var flags clw.MemoryFlags
	if mf.Read {
		if mf.Write {
			flags = clw.MemoryReadWrite
		} else {
			flags = clw.MemoryReadOnly
		}
	} else if mf.Write {
		flags = clw.MemoryWriteOnly
	}
	return flags
}

func (c *Context) CreateDeviceBuffer(size int, mf MemoryFlags) (*Buffer, error) {
	memory, err := clw.CreateBuffer(c.ID, mf.toBits(), clw.Size(size), nil)
	if err != nil {
		return nil, err
	}
	return &Buffer{ID: memory, Context: c, Size: size, Flags: mf}, nil
}

func (c *Context) CreateDeviceBufferFromHost(mf MemoryFlags, host []byte) (*Buffer, error) {
	flags := mf.toBits() | clw.MemoryCopyHostPointer
	memory, err := clw.CreateBuffer(c.ID, flags, clw.Size(len(host)), host)
	if err != nil {
		return nil, err
	}
	return &Buffer{ID: memory, Context: c, Size: len(host), Flags: mf}, nil
}

func (c *Context) CreateDeviceBufferOnHost(mf MemoryFlags, host []byte) (*Buffer, error) {
	flags := mf.toBits() | clw.MemoryUseHostPointer
	memory, err := clw.CreateBuffer(c.ID, flags, clw.Size(len(host)), host)
	if err != nil {
		return nil, err
	}
	return &Buffer{ID: memory, Context: c, Size: len(host), Host: host, Flags: mf}, nil
}

func (c *Context) CreateHostBuffer(size int, mf MemoryFlags) (*Buffer, error) {
	flags := mf.toBits() | clw.MemoryAllocHostPointer
	memory, err := clw.CreateBuffer(c.ID, flags, clw.Size(size), nil)
	if err != nil {
		return nil, err
	}
	return &Buffer{ID: memory, Context: c, Size: size, Flags: mf}, nil
}

func (c *Context) CreateHostBufferFromHost(mf MemoryFlags, host []byte) (*Buffer, error) {
	flags := mf.toBits() | clw.MemoryAllocHostPointer | clw.MemoryCopyHostPointer
	memory, err := clw.CreateBuffer(c.ID, flags, clw.Size(len(host)), host)
	if err != nil {
		return nil, err
	}
	return &Buffer{ID: memory, Context: c, Size: len(host), Flags: mf}, nil
}

func (b *Buffer) Release() error {
	return clw.ReleaseMemObject(b.ID)
}

func (cq *CommandQueue) CopyBuffer(src, dst *Buffer, srcOffset, dstOffset, size int, waitList []Event, e *Event) error {
	return clw.EnqueueCopyBuffer(cq.ID, src.ID, dst.ID, clw.Size(srcOffset), clw.Size(dstOffset), clw.Size(size),
		toEvents(waitList), (*clw.Event)(e))
}

func (cq *CommandQueue) MapBuffer(b *Buffer, blocking bool, flags MapFlags, offset, size int, waitList []Event,
	e *Event) ([]byte, error) {

	var mapFlags clw.MapFlags
	if flags.Read {
		mapFlags |= clw.MapRead
	}
	if flags.Write {
		mapFlags |= clw.MapWrite
	}

	mapped, err := clw.EnqueueMapBuffer(cq.ID, b.ID, clw.ToBool(blocking), mapFlags, clw.Size(offset), clw.Size(size),
		toEvents(waitList), (*clw.Event)(e))
	if err != nil {
		return nil, err
	}
	return mapped, nil
}

func (cq *CommandQueue) UnmapBuffer(b *Buffer, mapped []byte, waitList []Event, e *Event) error {
	return clw.EnqueueUnmapMemObject(cq.ID, b.ID, mapped, toEvents(waitList), (*clw.Event)(e))
}
