// A wrapper package that attempts to map the OpenCL 1.1 C API to idiomatic Go.
//
// TODO describe how empty interfaces are interpreted.
package cl11

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"unsafe"

	clw "github.com/rdwilliamson/clw11"
)

const (
	scratchSize = unsafe.Sizeof(Double16{})

	charSize   = unsafe.Sizeof(Char(0))
	shortSize  = unsafe.Sizeof(Short(0))
	intSize    = unsafe.Sizeof(Int(0))
	longSize   = unsafe.Sizeof(Long(0))
	floatSize  = unsafe.Sizeof(Float(0))
	doubleSize = unsafe.Sizeof(Double(0))
)

var (
	int32Type   = reflect.TypeOf(int32(0))
	int32Size   = int32Type.Size()
	uint32Type  = reflect.TypeOf(uint32(0))
	uint32Size  = uint32Type.Size()
	float32Type = reflect.TypeOf(float32(0))
	float32Size = float32Type.Size()
)

var errNotAddressable = errors.New("value not addressable")

type (
	Char   clw.Char
	Uchar  clw.Uchar
	Short  clw.Short
	Ushort clw.Ushort
	Int    clw.Int
	Uint   clw.Uint
	Long   clw.Long
	Ulong  clw.Ulong

	Half   clw.Half
	Float  clw.Float
	Double clw.Double

	Char2  clw.Char2
	Char4  clw.Char4
	Char8  clw.Char8
	Char16 clw.Char16

	Uchar2  clw.Uchar2
	Uchar4  clw.Uchar4
	Uchar8  clw.Uchar8
	Uchar16 clw.Uchar16

	Short2  clw.Short2
	Short4  clw.Short4
	Short8  clw.Short8
	Short16 clw.Short16

	Ushort2  clw.Ushort2
	Ushort4  clw.Ushort4
	Ushort8  clw.Ushort8
	Ushort16 clw.Ushort16

	Int2  clw.Int2
	Int4  clw.Int4
	Int8  clw.Int8
	Int16 clw.Int16

	Uint2  clw.Uint2
	Uint4  clw.Uint4
	Uint8  clw.Uint8
	Uint16 clw.Uint16

	Long2  clw.Long2
	Long4  clw.Long4
	Long8  clw.Long8
	Long16 clw.Long16

	Ulong2  clw.Ulong2
	Ulong4  clw.Ulong4
	Ulong8  clw.Ulong8
	Ulong16 clw.Ulong16

	Float2  clw.Float2
	Float4  clw.Float4
	Float8  clw.Float8
	Float16 clw.Float16

	Double2  clw.Double2
	Double4  clw.Double4
	Double8  clw.Double8
	Double16 clw.Double16
)

func checkIndex(i, max int) {
	if i < 0 || i >= max {
		panic("index out of range")
	}
}

// TODO use go generate

func NewChar2(v1, v2 Char) Char2 {
	return *(*Char2)(unsafe.Pointer(&[2]Char{v1, v2}))
}

func (v *Char2) Get(i int) Char {
	// checkIndex(i, 2)
	// return *(*Char)(unsafe.Pointer(uintptr(unsafe.Pointer(v)) + uintptr(i)*charSize))

	return Char((*[2]clw.Char)(unsafe.Pointer(v))[i])
}

func (v2 *Char2) Set(i int, v Char) {
	checkIndex(i, 2)
	*(*Char)(unsafe.Pointer(uintptr(unsafe.Pointer(v2)) + uintptr(i)*charSize)) = v
}

func (v4 *Char4) Get(i int) Char {
	checkIndex(i, 4)
	return *(*Char)(unsafe.Pointer(uintptr(unsafe.Pointer(v4)) + uintptr(i)*charSize))
}

func (v4 *Char4) Set(i int, v Char) {
	checkIndex(i, 4)
	*(*Char)(unsafe.Pointer(uintptr(unsafe.Pointer(v4)) + uintptr(i)*charSize)) = v
}

func (v8 *Char8) Get(i int) Char {
	checkIndex(i, 8)
	return *(*Char)(unsafe.Pointer(uintptr(unsafe.Pointer(v8)) + uintptr(i)*charSize))
}

func (v8 *Char8) Set(i int, v Char) {
	checkIndex(i, 8)
	*(*Char)(unsafe.Pointer(uintptr(unsafe.Pointer(v8)) + uintptr(i)*charSize)) = v
}

func (v16 *Char16) Get(i int) Char {
	checkIndex(i, 16)
	return *(*Char)(unsafe.Pointer(uintptr(unsafe.Pointer(v16)) + uintptr(i)*charSize))
}

func (v16 *Char16) Set(i int, v Char) {
	checkIndex(i, 16)
	*(*Char)(unsafe.Pointer(uintptr(unsafe.Pointer(v16)) + uintptr(i)*charSize)) = v
}

func (v2 *Uchar2) Get(i int) Uchar {
	checkIndex(i, 2)
	return *(*Uchar)(unsafe.Pointer(uintptr(unsafe.Pointer(v2)) + uintptr(i)*charSize))
}

func (v2 *Uchar2) Set(i int, v Uchar) {
	checkIndex(i, 2)
	*(*Uchar)(unsafe.Pointer(uintptr(unsafe.Pointer(v2)) + uintptr(i)*charSize)) = v
}

func (v4 *Uchar4) Get(i int) Uchar {
	checkIndex(i, 4)
	return *(*Uchar)(unsafe.Pointer(uintptr(unsafe.Pointer(v4)) + uintptr(i)*charSize))
}

func (v4 *Uchar4) Set(i int, v Uchar) {
	checkIndex(i, 4)
	*(*Uchar)(unsafe.Pointer(uintptr(unsafe.Pointer(v4)) + uintptr(i)*charSize)) = v
}

func (v8 *Uchar8) Get(i int) Uchar {
	checkIndex(i, 8)
	return *(*Uchar)(unsafe.Pointer(uintptr(unsafe.Pointer(v8)) + uintptr(i)*charSize))
}

func (v8 *Uchar8) Set(i int, v Uchar) {
	checkIndex(i, 8)
	*(*Uchar)(unsafe.Pointer(uintptr(unsafe.Pointer(v8)) + uintptr(i)*charSize)) = v
}

func (v16 *Uchar16) Get(i int) Uchar {
	checkIndex(i, 16)
	return *(*Uchar)(unsafe.Pointer(uintptr(unsafe.Pointer(v16)) + uintptr(i)*charSize))
}

func (v16 *Uchar16) Set(i int, v Uchar) {
	checkIndex(i, 16)
	*(*Uchar)(unsafe.Pointer(uintptr(unsafe.Pointer(v16)) + uintptr(i)*charSize)) = v
}

func (v2 *Short2) Get(i int) Short {
	checkIndex(i, 2)
	return *(*Short)(unsafe.Pointer(uintptr(unsafe.Pointer(v2)) + uintptr(i)*shortSize))
}

func (v2 *Short2) Set(i int, v Short) {
	checkIndex(i, 2)
	*(*Short)(unsafe.Pointer(uintptr(unsafe.Pointer(v2)) + uintptr(i)*shortSize)) = v
}

func (v4 *Short4) Get(i int) Short {
	checkIndex(i, 4)
	return *(*Short)(unsafe.Pointer(uintptr(unsafe.Pointer(v4)) + uintptr(i)*shortSize))
}

func (v4 *Short4) Set(i int, v Short) {
	checkIndex(i, 4)
	*(*Short)(unsafe.Pointer(uintptr(unsafe.Pointer(v4)) + uintptr(i)*shortSize)) = v
}

func (v8 *Short8) Get(i int) Short {
	checkIndex(i, 8)
	return *(*Short)(unsafe.Pointer(uintptr(unsafe.Pointer(v8)) + uintptr(i)*shortSize))
}

func (v8 *Short8) Set(i int, v Short) {
	checkIndex(i, 8)
	*(*Short)(unsafe.Pointer(uintptr(unsafe.Pointer(v8)) + uintptr(i)*shortSize)) = v
}

func (v16 *Short16) Get(i int) Short {
	checkIndex(i, 16)
	return *(*Short)(unsafe.Pointer(uintptr(unsafe.Pointer(v16)) + uintptr(i)*shortSize))
}

func (v16 *Short16) Set(i int, v Short) {
	checkIndex(i, 16)
	*(*Short)(unsafe.Pointer(uintptr(unsafe.Pointer(v16)) + uintptr(i)*shortSize)) = v
}

func (v2 *Ushort2) Get(i int) Ushort {
	checkIndex(i, 2)
	return *(*Ushort)(unsafe.Pointer(uintptr(unsafe.Pointer(v2)) + uintptr(i)*shortSize))
}

func (v2 *Ushort2) Set(i int, v Ushort) {
	checkIndex(i, 2)
	*(*Ushort)(unsafe.Pointer(uintptr(unsafe.Pointer(v2)) + uintptr(i)*shortSize)) = v
}

func (v4 *Ushort4) Get(i int) Ushort {
	checkIndex(i, 4)
	return *(*Ushort)(unsafe.Pointer(uintptr(unsafe.Pointer(v4)) + uintptr(i)*shortSize))
}

func (v4 *Ushort4) Set(i int, v Ushort) {
	checkIndex(i, 4)
	*(*Ushort)(unsafe.Pointer(uintptr(unsafe.Pointer(v4)) + uintptr(i)*shortSize)) = v
}

func (v8 *Ushort8) Get(i int) Ushort {
	checkIndex(i, 8)
	return *(*Ushort)(unsafe.Pointer(uintptr(unsafe.Pointer(v8)) + uintptr(i)*shortSize))
}

func (v8 *Ushort8) Set(i int, v Ushort) {
	checkIndex(i, 8)
	*(*Ushort)(unsafe.Pointer(uintptr(unsafe.Pointer(v8)) + uintptr(i)*shortSize)) = v
}

func (v16 *Ushort16) Get(i int) Ushort {
	checkIndex(i, 16)
	return *(*Ushort)(unsafe.Pointer(uintptr(unsafe.Pointer(v16)) + uintptr(i)*shortSize))
}

func (v16 *Ushort16) Set(i int, v Ushort) {
	checkIndex(i, 16)
	*(*Ushort)(unsafe.Pointer(uintptr(unsafe.Pointer(v16)) + uintptr(i)*shortSize)) = v
}

func (v2 *Int2) Get(i int) Int {
	checkIndex(i, 2)
	return *(*Int)(unsafe.Pointer(uintptr(unsafe.Pointer(v2)) + uintptr(i)*intSize))
}

func (v2 *Int2) Set(i int, v Int) {
	checkIndex(i, 2)
	*(*Int)(unsafe.Pointer(uintptr(unsafe.Pointer(v2)) + uintptr(i)*intSize)) = v
}

func (v4 *Int4) Get(i int) Int {
	checkIndex(i, 4)
	return *(*Int)(unsafe.Pointer(uintptr(unsafe.Pointer(v4)) + uintptr(i)*intSize))
}

func (v4 *Int4) Set(i int, v Int) {
	checkIndex(i, 4)
	*(*Int)(unsafe.Pointer(uintptr(unsafe.Pointer(v4)) + uintptr(i)*intSize)) = v
}

func (v8 *Int8) Get(i int) Int {
	checkIndex(i, 8)
	return *(*Int)(unsafe.Pointer(uintptr(unsafe.Pointer(v8)) + uintptr(i)*intSize))
}

func (v8 *Int8) Set(i int, v Int) {
	checkIndex(i, 8)
	*(*Int)(unsafe.Pointer(uintptr(unsafe.Pointer(v8)) + uintptr(i)*intSize)) = v
}

func (v16 *Int16) Get(i int) Int {
	checkIndex(i, 16)
	return *(*Int)(unsafe.Pointer(uintptr(unsafe.Pointer(v16)) + uintptr(i)*intSize))
}

func (v16 *Int16) Set(i int, v Int) {
	checkIndex(i, 16)
	*(*Int)(unsafe.Pointer(uintptr(unsafe.Pointer(v16)) + uintptr(i)*intSize)) = v
}

func (v2 *Uint2) Get(i int) Uint {
	checkIndex(i, 2)
	return *(*Uint)(unsafe.Pointer(uintptr(unsafe.Pointer(v2)) + uintptr(i)*intSize))
}

func (v2 *Uint2) Set(i int, v Uint) {
	checkIndex(i, 2)
	*(*Uint)(unsafe.Pointer(uintptr(unsafe.Pointer(v2)) + uintptr(i)*intSize)) = v
}

func (v4 *Uint4) Get(i int) Uint {
	checkIndex(i, 4)
	return *(*Uint)(unsafe.Pointer(uintptr(unsafe.Pointer(v4)) + uintptr(i)*intSize))
}

func (v4 *Uint4) Set(i int, v Uint) {
	checkIndex(i, 4)
	*(*Uint)(unsafe.Pointer(uintptr(unsafe.Pointer(v4)) + uintptr(i)*intSize)) = v
}

func (v8 *Uint8) Get(i int) Uint {
	checkIndex(i, 8)
	return *(*Uint)(unsafe.Pointer(uintptr(unsafe.Pointer(v8)) + uintptr(i)*intSize))
}

func (v8 *Uint8) Set(i int, v Uint) {
	checkIndex(i, 8)
	*(*Uint)(unsafe.Pointer(uintptr(unsafe.Pointer(v8)) + uintptr(i)*intSize)) = v
}

func (v16 *Uint16) Get(i int) Uint {
	checkIndex(i, 16)
	return *(*Uint)(unsafe.Pointer(uintptr(unsafe.Pointer(v16)) + uintptr(i)*intSize))
}

func (v16 *Uint16) Set(i int, v Uint) {
	checkIndex(i, 16)
	*(*Uint)(unsafe.Pointer(uintptr(unsafe.Pointer(v16)) + uintptr(i)*intSize)) = v
}

func (v2 *Long2) Get(i int) Long {
	checkIndex(i, 2)
	return *(*Long)(unsafe.Pointer(uintptr(unsafe.Pointer(v2)) + uintptr(i)*longSize))
}

func (v2 *Long2) Set(i int, v Long) {
	checkIndex(i, 2)
	*(*Long)(unsafe.Pointer(uintptr(unsafe.Pointer(v2)) + uintptr(i)*longSize)) = v
}

func (v4 *Long4) Get(i int) Long {
	checkIndex(i, 4)
	return *(*Long)(unsafe.Pointer(uintptr(unsafe.Pointer(v4)) + uintptr(i)*longSize))
}

func (v4 *Long4) Set(i int, v Long) {
	checkIndex(i, 4)
	*(*Long)(unsafe.Pointer(uintptr(unsafe.Pointer(v4)) + uintptr(i)*longSize)) = v
}

func (v8 *Long8) Get(i int) Long {
	checkIndex(i, 8)
	return *(*Long)(unsafe.Pointer(uintptr(unsafe.Pointer(v8)) + uintptr(i)*longSize))
}

func (v8 *Long8) Set(i int, v Long) {
	checkIndex(i, 8)
	*(*Long)(unsafe.Pointer(uintptr(unsafe.Pointer(v8)) + uintptr(i)*longSize)) = v
}

func (v16 *Long16) Get(i int) Long {
	checkIndex(i, 16)
	return *(*Long)(unsafe.Pointer(uintptr(unsafe.Pointer(v16)) + uintptr(i)*longSize))
}

func (v16 *Long16) Set(i int, v Long) {
	checkIndex(i, 16)
	*(*Long)(unsafe.Pointer(uintptr(unsafe.Pointer(v16)) + uintptr(i)*longSize)) = v
}

func (v2 *Ulong2) Get(i int) Ulong {
	checkIndex(i, 2)
	return *(*Ulong)(unsafe.Pointer(uintptr(unsafe.Pointer(v2)) + uintptr(i)*longSize))
}

func (v2 *Ulong2) Set(i int, v Ulong) {
	checkIndex(i, 2)
	*(*Ulong)(unsafe.Pointer(uintptr(unsafe.Pointer(v2)) + uintptr(i)*longSize)) = v
}

func (v4 *Ulong4) Get(i int) Ulong {
	checkIndex(i, 4)
	return *(*Ulong)(unsafe.Pointer(uintptr(unsafe.Pointer(v4)) + uintptr(i)*longSize))
}

func (v4 *Ulong4) Set(i int, v Ulong) {
	checkIndex(i, 4)
	*(*Ulong)(unsafe.Pointer(uintptr(unsafe.Pointer(v4)) + uintptr(i)*longSize)) = v
}

func (v8 *Ulong8) Get(i int) Ulong {
	checkIndex(i, 8)
	return *(*Ulong)(unsafe.Pointer(uintptr(unsafe.Pointer(v8)) + uintptr(i)*longSize))
}

func (v8 *Ulong8) Set(i int, v Ulong) {
	checkIndex(i, 8)
	*(*Ulong)(unsafe.Pointer(uintptr(unsafe.Pointer(v8)) + uintptr(i)*longSize)) = v
}

func (v16 *Ulong16) Get(i int) Ulong {
	checkIndex(i, 16)
	return *(*Ulong)(unsafe.Pointer(uintptr(unsafe.Pointer(v16)) + uintptr(i)*longSize))
}

func (v16 *Ulong16) Set(i int, v Ulong) {
	checkIndex(i, 16)
	*(*Ulong)(unsafe.Pointer(uintptr(unsafe.Pointer(v16)) + uintptr(i)*longSize)) = v
}

func (v2 *Float2) Get(i int) Float {
	checkIndex(i, 2)
	return *(*Float)(unsafe.Pointer(uintptr(unsafe.Pointer(v2)) + uintptr(i)*floatSize))
}

func (v2 *Float2) Set(i int, v Float) {
	checkIndex(i, 2)
	*(*Float)(unsafe.Pointer(uintptr(unsafe.Pointer(v2)) + uintptr(i)*floatSize)) = v
}

func (v4 *Float4) Get(i int) Float {
	checkIndex(i, 4)
	return *(*Float)(unsafe.Pointer(uintptr(unsafe.Pointer(v4)) + uintptr(i)*floatSize))
}

func (v4 *Float4) Set(i int, v Float) {
	checkIndex(i, 4)
	*(*Float)(unsafe.Pointer(uintptr(unsafe.Pointer(v4)) + uintptr(i)*floatSize)) = v
}

func (v8 *Float8) Get(i int) Float {
	checkIndex(i, 8)
	return *(*Float)(unsafe.Pointer(uintptr(unsafe.Pointer(v8)) + uintptr(i)*floatSize))
}

func (v8 *Float8) Set(i int, v Float) {
	checkIndex(i, 8)
	*(*Float)(unsafe.Pointer(uintptr(unsafe.Pointer(v8)) + uintptr(i)*floatSize)) = v
}

func (v16 *Float16) Get(i int) Float {
	checkIndex(i, 16)
	return *(*Float)(unsafe.Pointer(uintptr(unsafe.Pointer(v16)) + uintptr(i)*floatSize))
}

func (v16 *Float16) Set(i int, v Float) {
	checkIndex(i, 16)
	*(*Float)(unsafe.Pointer(uintptr(unsafe.Pointer(v16)) + uintptr(i)*floatSize)) = v
}

func (v2 *Double2) Get(i int) Double {
	checkIndex(i, 2)
	return *(*Double)(unsafe.Pointer(uintptr(unsafe.Pointer(v2)) + uintptr(i)*doubleSize))
}

func (v2 *Double2) Set(i int, v Double) {
	checkIndex(i, 2)
	*(*Double)(unsafe.Pointer(uintptr(unsafe.Pointer(v2)) + uintptr(i)*doubleSize)) = v
}

func (v4 *Double4) Get(i int) Double {
	checkIndex(i, 4)
	return *(*Double)(unsafe.Pointer(uintptr(unsafe.Pointer(v4)) + uintptr(i)*doubleSize))
}

func (v4 *Double4) Set(i int, v Double) {
	checkIndex(i, 4)
	*(*Double)(unsafe.Pointer(uintptr(unsafe.Pointer(v4)) + uintptr(i)*doubleSize)) = v
}

func (v8 *Double8) Get(i int) Double {
	checkIndex(i, 8)
	return *(*Double)(unsafe.Pointer(uintptr(unsafe.Pointer(v8)) + uintptr(i)*doubleSize))
}

func (v8 *Double8) Set(i int, v Double) {
	checkIndex(i, 8)
	*(*Double)(unsafe.Pointer(uintptr(unsafe.Pointer(v8)) + uintptr(i)*doubleSize)) = v
}

func (v16 *Double16) Get(i int) Double {
	checkIndex(i, 16)
	return *(*Double)(unsafe.Pointer(uintptr(unsafe.Pointer(v16)) + uintptr(i)*doubleSize))
}

func (v16 *Double16) Set(i int, v Double) {
	checkIndex(i, 16)
	*(*Double)(unsafe.Pointer(uintptr(unsafe.Pointer(v16)) + uintptr(i)*doubleSize)) = v
}

type Profile int

func toProfile(profile string) Profile {

	switch profile {
	case "FULL_PROFILE":
		return FullProfile
	case "EMBEDDED_PROFILE":
		return EmbeddedProfile
	}

	panic(errors.New("unknown profile"))
}

func (pp Profile) String() string {
	switch pp {
	case zeroProfile:
		return ""
	case FullProfile:
		return "full profile"
	case EmbeddedProfile:
		return "embedded profile"
	}
	panic("unreachable")
}

type Version struct {
	Major int
	Minor int
	Info  string
}

func toVersion(version string) Version {

	var result Version
	var err error

	if strings.HasPrefix(version, "OpenCL C") {
		_, err = fmt.Sscanf(version, "OpenCL C %d.%d", &result.Major, &result.Minor)
		result.Info = strings.TrimSpace(version[len(fmt.Sprintf("OpenCL C %d.%d", result.Major, result.Minor)):])

	} else if strings.HasPrefix(version, "OpenCL") {
		_, err = fmt.Sscanf(version, "OpenCL %d.%d", &result.Major, &result.Minor)
		result.Info = strings.TrimSpace(version[len(fmt.Sprintf("OpenCL %d.%d", result.Major, result.Minor)):])

	} else {
		_, err = fmt.Sscanf(version, "%d.%d", &result.Major, &result.Minor)
	}

	if err != nil {
		panic(err)
	}

	return result
}

func (v Version) String() string {
	if v.Info != "" {
		return fmt.Sprintf("%d.%d %s", v.Major, v.Minor, v.Info)
	}
	return fmt.Sprintf("%d.%d", v.Major, v.Minor)
}

func toByteSlice(p unsafe.Pointer, size uintptr) []byte {

	var header reflect.SliceHeader
	header.Data = uintptr(p)
	header.Len = int(size)
	header.Cap = int(size)

	return *(*[]byte)(unsafe.Pointer(&header))
}

func toBytes(x interface{}, scratch unsafe.Pointer) []byte {
	pointer, size := getPointerAndSize(x, scratch)
	return toByteSlice(pointer, size)
}

func tryPointerAndSize(x interface{}) (pointer unsafe.Pointer, size uintptr, err error) {

	value := reflect.ValueOf(x)
	switch value.Kind() {

	case reflect.Ptr:
		for {
			value = value.Elem()

			if value.Kind() != reflect.Ptr {
				break
			}
		}
		pointer, size = addressablePointerAndSize(value)
		return

	case reflect.Slice:
		pointer, size = addressablePointerAndSize(value)
		return

	default:
		err = errNotAddressable
		return
	}

	panic("unreachable")
}

func getPointerAndSize(x interface{}, scratch unsafe.Pointer) (unsafe.Pointer, uintptr) {

	value := reflect.ValueOf(x)
	switch value.Kind() {

	case reflect.Int, reflect.Int32:
		newValue := reflect.NewAt(int32Type, scratch).Elem()
		newValue.SetInt(value.Int())
		return addressablePointerAndSize(newValue)

	case reflect.Ptr:
		for {
			value = value.Elem()

			if kind := value.Kind(); kind == reflect.Int {
				newValue := reflect.NewAt(int32Type, scratch).Elem()
				newValue.SetInt(value.Int())
				value = newValue
				break

			} else if kind != reflect.Ptr {
				break
			}
		}
		return addressablePointerAndSize(value)

	case reflect.Slice:
		return addressablePointerAndSize(value)

	}
	panic("unsupported kind")
}

func addressablePointerAndSize(v reflect.Value) (unsafe.Pointer, uintptr) {

	switch v.Kind() {

	case reflect.Int32:
		return unsafe.Pointer(v.UnsafeAddr()), v.Type().Size()

	case reflect.Slice:
		pointer := unsafe.Pointer(&(*reflect.SliceHeader)(unsafe.Pointer(v.Pointer())).Data)
		size := v.Type().Elem().Size() * uintptr(v.Len())
		return pointer, size

	}
	panic("unsupported kind")
}
