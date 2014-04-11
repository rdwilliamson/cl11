package cl11

import (
	"reflect"
	"unsafe"

	clw "github.com/rdwilliamson/clw11"
)

// A mapped buffer or image. Has convenience functions for common data types.
type MappedBuffer struct {
	pointer unsafe.Pointer
	size    int64
	memID   clw.Mem
}

func (cq *CommandQueue) UnmapBuffer(mb *MappedBuffer, waitList []*Event, e *Event) error {

	var event *clw.Event
	if e != nil {
		event = &e.id
		e.Context = cq.Context
		e.CommandType = CommandUnmapMemoryObject
		e.CommandQueue = cq
	}

	return clw.EnqueueUnmapMemObject(cq.id, mb.memID, mb.pointer, cq.toEvents(waitList), event)
}

// Returns a slice of float32s backed by the mapped buffer.
func (bm *MappedBuffer) Float32Slice() []float32 {
	var header reflect.SliceHeader
	header.Data = uintptr(bm.pointer)
	size := int(bm.size / int64(float32Size))
	header.Len = size
	header.Cap = size
	return *(*[]float32)(unsafe.Pointer(&header))
}
