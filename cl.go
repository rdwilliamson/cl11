// A wrapper package that attempts to map the OpenCL 1.1 C API to idiomatic Go.
package cl11

import (
	"reflect"
	"unsafe"
)

// TODO move?
var (
	int32Type  = reflect.TypeOf(int32(0))
	int32Size  = int32Type.Size()
	uint32Type = reflect.TypeOf(uint32(0))
	uint32Size = uint32Type.Size()
)

func toByteSlice(p unsafe.Pointer, size int) []byte {

	var header reflect.SliceHeader
	header.Data = uintptr(p)
	header.Len = size
	header.Cap = size

	return *(*[]byte)(unsafe.Pointer(&header))
}

/*
// Scratch must be (at least) 8 bytes of scratch space.
func toBytes(x interface{}, scratch unsafe.Pointer) []byte {
	pointer, size := getPointerAndSize(x, scratch)
	return toByteSlice(pointer, size)
}

func getPointerAndSize(x interface{}, scratch unsafe.Pointer) (unsafe.Pointer, int) {

	value := reflect.ValueOf(x)
	switch value.Kind() {

	case reflect.Int, reflect.Int32:
		newValue := reflect.NewAt(int32Type, scratch).Elem()
		newValue.SetInt(value.Int())
		return addressablePointerAndSize(newValue)

	case reflect.Ptr:

		for {
			value = value.Elem()
			if value.Kind() != reflect.Ptr {
				break
			}
		}
		return addressablePointerAndSize(value)

	case reflect.Slice:
		return addressablePointerAndSize(value)

	}
	panic("unsupported kind")
}

func addressablePointerAndSize(v reflect.Value) (unsafe.Pointer, int) {

	switch v.Kind() {

	case reflect.Int:
		return unsafe.Pointer(v.UnsafeAddr()), int32Size

	case reflect.Int32:
		return unsafe.Pointer(v.UnsafeAddr()), int(v.Type().Size())

	case reflect.Slice:
		pointer := unsafe.Pointer(&(*reflect.SliceHeader)(unsafe.Pointer(v.Pointer())).Data)
		size := int(v.Type().Elem().Size()) * v.Len()
		return pointer, size

	}
	panic("unsupported kind")
}
*/
