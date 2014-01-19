package cl11

import clw "github.com/rdwilliamson/clw11"

// 3 cases:
// 1. Device memory
// 2. Host memory alloced by the CL
// 3. Host memory alloced by go
type Buffer struct {
	ID     clw.Memory
	Device bool   // Is the memory on the device.
	Host   []byte // Host backed memory (if applicable).

	// mapped  []byte
	// Read    bool
	// Write   bool
	// Host    bool // is it on the host or device
	// alloced bool // if the implementation alloced the memory
}

func CreateDeviceBuffer(c *Context, size int, read, write bool, host []byte, copyHost bool) (*Buffer, error) {
	var flags clw.MemoryFlags
	if read {
		if write {
			flags = clw.MemoryReadWrite
		} else {
			flags = clw.MemoryReadOnly
		}
	} else if write {
		flags = clw.MemoryWriteOnly
	}

	if copyHost {
		flags |= clw.MemoryCopyHostPointer
	}

	memory, err := clw.CreateBuffer(c.ID, flags, clw.Size(size), host)
	if err != nil {
		return nil, err
	}
	return &Buffer{ID: memory, Device: true}, nil
}

func CreateHostBuffer(c *Context, size int, read, write bool) (*Buffer, error) {
	flags := clw.MemoryAllocHostPointer

	if read {
		if write {
			flags |= clw.MemoryReadWrite
		} else {
			flags |= clw.MemoryReadOnly
		}

	} else {
		flags |= clw.MemoryWriteOnly
	}

	memory, err := clw.CreateBuffer(c.ID, flags, clw.Size(size), nil)
	if err != nil {
		return nil, err
	}
	return &Buffer{ID: memory, Device: false}, nil
}

func (b *Buffer) Release() error {
	return clw.ReleaseMemObject(b.ID)
}
