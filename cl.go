// A wrapper package that attempts to map the OpenCL 1.1 C API to idiomatic Go.
package cl11

import (
	"reflect"
	"unsafe"
)

var (
	int32Type = reflect.TypeOf(int32(0))
	int32Size = int32Type.Size()
)

// Scratch must be (at least) 8 bytes of scratch space.
func toBytes(x interface{}, scratch unsafe.Pointer) []byte {

	pointer, size := toPointerAndSize(x, scratch)

	var header reflect.SliceHeader
	header.Data = uintptr(pointer)
	header.Len = int(size)
	header.Cap = int(size)

	return *(*[]byte)(unsafe.Pointer(&header))
}

func toPointerAndSize(x interface{}, scratch unsafe.Pointer) (unsafe.Pointer, uintptr) {

	value := reflect.ValueOf(x)
	switch value.Kind() {

	case reflect.Int, reflect.Int32:
		newValue := reflect.NewAt(int32Type, scratch).Elem()
		newValue.SetInt(value.Int())
		return safePointerAndSize(newValue)

	case reflect.Ptr:

		for {
			value = value.Elem()
			if value.Kind() != reflect.Ptr {
				break
			}
		}
		return safePointerAndSize(value)

	case reflect.Slice:
		return safePointerAndSize(value)

	}
	panic("unsupported kind")
}

func safePointerAndSize(v reflect.Value) (unsafe.Pointer, uintptr) {

	switch v.Kind() {

	case reflect.Int:
		return unsafe.Pointer(v.UnsafeAddr()), int32Size

	case reflect.Int32:
		return unsafe.Pointer(v.UnsafeAddr()), v.Type().Size()

	case reflect.Slice:
		pointer := unsafe.Pointer(&(*reflect.SliceHeader)(unsafe.Pointer(v.Pointer())).Data)
		size := v.Type().Elem().Size() * uintptr(v.Len())
		return pointer, size

	}
	panic("unsupported kind")
}
