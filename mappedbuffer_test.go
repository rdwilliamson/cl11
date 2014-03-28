package cl11

import (
	"reflect"
	"testing"
	"unsafe"
)

func TestBufferMapping(t *testing.T) {
	values := []float32{1, 2, 3, 4}
	mapped := MappedBuffer{unsafe.Pointer(&values[0]), int64(len(values)) * int64(float32Size), nil}
	slice := mapped.Float32Slice()
	if !reflect.DeepEqual(values, slice) {
		t.Errorf("float buffer mismatch")
	}
}
