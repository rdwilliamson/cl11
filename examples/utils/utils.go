package utils

import (
	"reflect"
	"unsafe"
)

func ToBytes(v interface{}) []byte {

	if v == nil {
		return nil
	}

	value := reflect.ValueOf(v)

	var result []byte
	header := (*reflect.SliceHeader)((unsafe.Pointer(&result)))

	size := value.Type().Size()
	header.Cap = int(size)
	header.Len = int(size)

	header.Data = value.UnsafeAddr()

	return result
}
