package cl11

import (
	"testing"
	"unsafe"
)

func sizeCheck(want, got int, t *testing.T) {
	if want != got {
		t.Error("size mismatch, wanted", want, "got", got)
	}
}

func pointerCheck(want, got unsafe.Pointer, t *testing.T) {
	if want != got {
		t.Error("pointer mismatch, wanted", want, "got", got)
	}
}

// Bool

func TestIntToBytes(t *testing.T) {

	var scratchSpace [8]byte
	scratch := unsafe.Pointer(&scratchSpace[0])

	var anInt int = 1
	bytes := ToBytes(anInt, scratch)

	sizeCheck(4, len(bytes), t)
	pointerCheck(scratch, unsafe.Pointer(&bytes[0]), t)
	if anInt != *(*int)(scratch) {
		t.Error("value mismatch, wanted", anInt, "got", *(*int)(scratch))
	}
}

// Int8
// Int16
// Int32
// Int64
// Uint
// Uint8
// Uint16
// Uint32
// Uint64
// Uintptr
// Float32
// Float64
// Complex64
// Complex128
// Array
// Interface

func TestPtrToBytes(t *testing.T) {

	var scratchSpace [8]byte
	scratch := unsafe.Pointer(&scratchSpace[0])

	var anInt int
	bytes := ToBytes(&anInt, scratch)

	sizeCheck(4, len(bytes), t)
	pointerCheck(unsafe.Pointer(&anInt), unsafe.Pointer(&bytes[0]), t)
}

func TestSliceToBytes(t *testing.T) {

	var scratchSpace [8]byte
	scratch := unsafe.Pointer(&scratchSpace[0])

	slice := []float32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	bytes := ToBytes(slice, scratch)

	sizeCheck(len(slice)*4, len(bytes), t)
	pointerCheck(unsafe.Pointer(&slice[0]), unsafe.Pointer(&bytes[0]), t)
}

// String
// Struct
