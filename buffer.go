package cl11

import clw "github.com/rdwilliamson/clw11"

type Buffer struct {
	ID      clw.Memory
	Context *Context
	Size    int
	Host    []byte
	Flags   MemoryFlags
}

type MemoryFlags struct {
	Read  bool // Can a kernel read from this buffer.
	Write bool // Can a kernel write to this buffer.
}

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

func CreateDeviceBuffer(c *Context, size int, mf MemoryFlags) (*Buffer, error) {
	memory, err := clw.CreateBuffer(c.ID, mf.toBits(), clw.Size(size), nil)
	if err != nil {
		return nil, err
	}
	return &Buffer{ID: memory, Context: c, Size: size, Flags: mf}, nil
}

func CreateDeviceBufferFromHost(c *Context, mf MemoryFlags, host []byte) (*Buffer, error) {
	flags := mf.toBits() | clw.MemoryCopyHostPointer
	memory, err := clw.CreateBuffer(c.ID, flags, clw.Size(len(host)), host)
	if err != nil {
		return nil, err
	}
	return &Buffer{ID: memory, Context: c, Size: len(host), Flags: mf}, nil
}

func CreateDeviceBufferOnHost(c *Context, mf MemoryFlags, host []byte) (*Buffer, error) {
	flags := mf.toBits() | clw.MemoryUseHostPointer
	memory, err := clw.CreateBuffer(c.ID, flags, clw.Size(len(host)), host)
	if err != nil {
		return nil, err
	}
	return &Buffer{ID: memory, Context: c, Size: len(host), Host: host, Flags: mf}, nil
}

func CreateHostBuffer(c *Context, size int, mf MemoryFlags) (*Buffer, error) {
	flags := mf.toBits() | clw.MemoryAllocHostPointer
	memory, err := clw.CreateBuffer(c.ID, flags, clw.Size(size), nil)
	if err != nil {
		return nil, err
	}
	return &Buffer{ID: memory, Context: c, Size: size, Flags: mf}, nil
}

func CreateHostBufferFromHost(c *Context, mf MemoryFlags, host []byte) (*Buffer, error) {
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
