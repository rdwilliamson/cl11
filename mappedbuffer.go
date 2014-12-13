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

// Returns a slice of bytes backed by the mapped buffer.
func (mb *MappedBuffer) Bytes() []byte {
	var header reflect.SliceHeader
	header.Data = uintptr(mb.pointer)
	size := int(mb.size)
	header.Len = size
	header.Cap = size
	return *(*[]byte)(unsafe.Pointer(&header))
}

// Returns a slice of float32s backed by the mapped buffer.
func (mb *MappedBuffer) Float32s() []float32 {
	var header reflect.SliceHeader
	header.Data = uintptr(mb.pointer)
	size := int(mb.size / int64(4))
	header.Len = size
	header.Cap = size
	return *(*[]float32)(unsafe.Pointer(&header))
}
