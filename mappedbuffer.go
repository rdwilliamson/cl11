package cl11

import (
	"reflect"
	"unsafe"
)

// A mapped buffer. Has convenience functions for common data types.
type MappedBuffer struct {
	Buffer  *Buffer
	pointer unsafe.Pointer
	size    int64
}

// Returns a slice of float32s backed by the mapped buffer.
func (mb *MappedBuffer) Float32Slice() []float32 {
	var header reflect.SliceHeader
	header.Data = uintptr(mb.pointer)
	size := int(mb.size / int64(float32Size))
	header.Len = size
	header.Cap = size
	return *(*[]float32)(unsafe.Pointer(&header))
}
