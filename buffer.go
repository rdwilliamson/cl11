package cl11

import clw "github.com/rdwilliamson/clw11"

// 3 cases:
// 1. Device memory
// 2. Host memory alloced by the CL
// 3. Host memory alloced by go
type Buffer struct {
	ID clw.Memory
	// mapped  []byte
	// Read    bool
	// Write   bool
	// Host    bool // is it on the host or device
	// alloced bool // if the implementation alloced the memory
}

func CreateBuffer(c *Context, size int, read, write, useHost, alloc, copyHost bool, host []byte) (*Buffer, error) {
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

	memory, err := clw.CreateBuffer(clw.Context(c.ID), flags, clw.Size(size), host)
	if err != nil {
		return nil, err
	}

	return &Buffer{memory}, nil
}
