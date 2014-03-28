package cl11

import (
	"testing"
	"unsafe"
)

var scratchSpace = unsafe.Pointer(new([scratchSize]byte))

func sizeCheck(want int, got uintptr, t *testing.T) {
	if want != int(got) {
		t.Error("size mismatch, wanted", want, "got", got)
	}
}

func pointerCheck(want, got unsafe.Pointer, t *testing.T) {
	if want != got {
		t.Error("pointer mismatch, wanted", want, "got", got)
	}
}

// Bool

func TestIntPointerAndSize(t *testing.T) {

	var scratchSpace [8]byte
	scratch := unsafe.Pointer(&scratchSpace[0])

	var anInt int = 1

	pointer, size := getPointerAndSize(anInt, scratch)

	sizeCheck(4, size, t)
	pointerCheck(scratch, pointer, t)
	if anInt != *(*int)(pointer) {
		t.Error("value mismatch, wanted", anInt, "got", *(*int)(scratch))
	}

	pointer, size = getPointerAndSize(&anInt, scratch)

	sizeCheck(4, size, t)
	// Expect scratch since int needs to be converted to int32.
	pointerCheck(scratch, pointer, t)
	if anInt != *(*int)(pointer) {
		t.Error("value mismatch, wanted", anInt, "got", *(*int)(scratch))
	}
}

// Int8
// Int16

func TestInt32PointerAndSize(t *testing.T) {

	var scratchSpace [8]byte
	scratch := unsafe.Pointer(&scratchSpace[0])

	var anInt int32 = 1

	pointer, size := getPointerAndSize(anInt, scratch)

	sizeCheck(4, size, t)
	pointerCheck(scratch, pointer, t)
	if anInt != *(*int32)(pointer) {
		t.Error("value mismatch, wanted", anInt, "got", *(*int)(scratch))
	}

	pointer, size = getPointerAndSize(&anInt, scratch)

	sizeCheck(4, size, t)
	pointerCheck(unsafe.Pointer(&anInt), pointer, t)
	if anInt != *(*int32)(pointer) {
		t.Error("value mismatch, wanted", anInt, "got", *(*int)(scratch))
	}
}

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
// Ptr
// Slice
// String
// Struct
