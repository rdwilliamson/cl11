package cl11

import (
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

	// The offset in bytes within the parent buffer, if applicable.
	Origin int64
}

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

// Creates a buffer object from an existing object with the passed offset and
// size in bytes.
//
// If the size and offset are the same as another sub buffer the implementation
// may return the same sub buffer and increment the reference count.
func (b *Buffer) CreateSubBuffer(mf MemFlags, origin, size int64) (*Buffer, error) {

	region := clw.BufferRegion{Origin: clw.Size(origin), Size: clw.Size(size)}

	memory, err := clw.CreateSubBuffer(b.id, clw.MemFlags(mf), region)
	if err != nil {
		return nil, err
	}

	return &Buffer{id: memory, Context: b.Context, Size: size, Flags: mf, Buffer: b, Origin: origin}, nil
}

// Increments the buffer object reference count.
//
// The OpenCL commands that return a buffer perform an implicit retain.
func (b *Buffer) Retain() error {
	return clw.RetainMemObject(b.id)
}

// Decrements the buffer object reference count.
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
func (cq *CommandQueue) EnqueueCopyBuffer(src, dst *Buffer, srcOffset, dstOffset, size int64, waitList []*Event,
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
func (cq *CommandQueue) EnqueueCopyBufferRect(src, dst *Buffer, r *Rect, waitList []*Event, e *Event) error {

	var event *clw.Event
	if e != nil {
		event = &e.id
		e.Context = cq.Context
		e.CommandType = CommandCopyBuffer
		e.CommandQueue = cq
	}

	return clw.EnqueueCopyBufferRect(cq.id, src.id, dst.id, r.Src.origin(), r.Dst.origin(), r.region(),
		r.Src.rowPitch(), r.Src.slicePitch(), r.Dst.rowPitch(), r.Dst.slicePitch(), cq.toEvents(waitList), event)
}

// Enqueues a command to map a region of the buffer object given by buffer into
// the host address space.
func (cq *CommandQueue) EnqueueMapBuffer(b *Buffer, bc BlockingCall, flags MapFlags, offset, size int64,
	waitList []*Event, e *Event) (*MappedBuffer, error) {

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

// Enqueues a command to unmap a previously mapped buffer object.
func (cq *CommandQueue) EnqueueUnmapBuffer(mb *MappedBuffer, waitList []*Event, e *Event) error {

	var event *clw.Event
	if e != nil {
		event = &e.id
		e.Context = cq.Context
		e.CommandType = CommandUnmapMemoryObject
		e.CommandQueue = cq
	}

	return clw.EnqueueUnmapMemObject(cq.id, mb.memID, mb.pointer, cq.toEvents(waitList), event)
}

// Enqueue commands to read from a buffer object to host memory.
//
// Offset is in bytes. The destination must be addressable. If the buffer object
// is backed by host memory then all commands that use it and sub buffers must
// have finished execution and it must not be mapped otherwise the results are
// undefined.
func (cq *CommandQueue) EnqueueReadBuffer(b *Buffer, bc BlockingCall, offset int64, dst interface{}, waitList []*Event,
	e *Event) error {

	var event *clw.Event
	if e != nil {
		event = &e.id
		e.Context = cq.Context
		e.CommandType = CommandReadBuffer
		e.CommandQueue = cq
	}

	pointer, size, err := pointerSize(dst)
	if err != nil {
		return err
	}

	err = clw.EnqueueReadBuffer(cq.id, b.id, clw.Bool(bc), clw.Size(offset), clw.Size(size), pointer,
		cq.toEvents(waitList), event)
	if err != nil {
		return err
	}

	// Set a no-op callback, just need hold a reference to the source data.
	if bc == NonBlocking {
		err = e.SetCallback(noOpEventCallback, dst)
		if err != nil {
			return err
		}
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
func (cq *CommandQueue) EnqueueWriteBuffer(b *Buffer, bc BlockingCall, offset int64, src interface{}, waitList []*Event,
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

	pointer, size, err := pointerSize(src)
	if err != nil {
		return err
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
// See Rect definition for how r is defined. The destination must be
// addressable.
func (cq *CommandQueue) EnqueueReadBufferRect(b *Buffer, bc BlockingCall, r *Rect, dst interface{},
	waitList []*Event, e *Event) error {

	var event *clw.Event
	if e != nil {
		event = &e.id
		e.Context = cq.Context
		e.CommandType = CommandReadBuffer
		e.CommandQueue = cq
	}

	pointer, _, err := pointerSize(dst)
	if err != nil {
		return err
	}

	err = clw.EnqueueReadBufferRect(cq.id, b.id, clw.Bool(bc), r.Dst.origin(), r.Src.origin(), r.region(),
		r.Dst.rowPitch(), r.Dst.slicePitch(), r.Src.rowPitch(), r.Dst.rowPitch(), pointer, cq.toEvents(waitList), event)
	if err != nil {
		return err
	}

	// Set a no-op callback, just need hold a reference to the source data.
	if bc == NonBlocking {
		err = e.SetCallback(noOpEventCallback, dst)
		if err != nil {
			return err
		}
	}

	return nil
}

// Enqueue commands to write a rectangular region to a buffer object from host
// memory.
//
// See Rect definition for how r is defined. If the buffer object is backed by
// host memory then all commands that use it and sub buffers must have finished
// execution and it must not be mapped otherwise the results are undefined. If
// the source is addressable a reference will be held to prevent it being
// garbage collected, it is the user's responsibility to ensure that the source
// data is valid at the time the write is actually performed. If the source
// isn't addressable a copy will be created.
func (cq *CommandQueue) EnqueueWriteBufferRect(b *Buffer, bc BlockingCall, r *Rect, src interface{},
	waitList []*Event, e *Event) error {

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

	pointer, _, err := pointerSize(src)
	if err != nil {
		return err
	}

	err = clw.EnqueueWriteBufferRect(cq.id, b.id, clw.Bool(bc), r.Dst.origin(), r.Src.origin(), r.region(),
		r.Dst.rowPitch(), r.Dst.slicePitch(), r.Src.rowPitch(), r.Src.slicePitch(), pointer, cq.toEvents(waitList),
		event)
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
