package cl11

import (
	"bytes"
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

func TestChar2(t *testing.T) {

	got := NewChar2(Char(1), Char(2))

	for i := 0; i < 2; i++ {
		gotValue := got.Get(i)
		wantValue := Char(i + 1)
		if gotValue != wantValue {
			t.Fatalf("get/new failure: index %d, want %d, got %d", i, gotValue, wantValue)
		}
	}

	var want []int8
	for i := 0; i < 2; i++ {
		v := int8(i) + 3
		got.Set(i, Char(v))
		want = append(want, v)
	}

	gotBytes := toByteSlice(unsafe.Pointer(&got), unsafe.Sizeof(got))
	wantBytes := toByteSlice(unsafe.Pointer(&want[0]), uintptr(len(want))*unsafe.Sizeof(want[0]))
	if !bytes.Equal(gotBytes, wantBytes) {
		t.Fatalf("set failure:\nwant [% x]\ngot  [% x]", wantBytes, gotBytes)
	}
}

func TestDouble16(t *testing.T) {
	var got Double16
	var want []float64
	for i := 0; i < 16; i++ {
		v := float64(i) + 1
		got.Set(i, Double(v))
		want = append(want, v)
	}
	for i := 0; i < len(want); i++ {
		gotBytes := toByteSlice(unsafe.Pointer(&got), unsafe.Sizeof(got))
		wantBytes := toByteSlice(unsafe.Pointer(&want[0]), uintptr(len(want))*unsafe.Sizeof(want[0]))
		if !bytes.Equal(gotBytes, wantBytes) {
			t.Fatalf("set failure:\nwant [% x]\ngot  [% x]", wantBytes, gotBytes)
		}
	}
}
