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
)

var (
	int32Type   = reflect.TypeOf(int32(0))
	int32Size   = int32Type.Size()
	uint32Type  = reflect.TypeOf(uint32(0))
	uint32Size  = uint32Type.Size()
	float32Type = reflect.TypeOf(float32(0))
	float32Size = float32Type.Size()
)

var (
	errNotAddressable    = errors.New("value not addressable")
	errValueSizeTooSmall = errors.New("value's size is too small")
)

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
		// Maybe a regexp to try and find a "\d.\d"? It works on nVidia, AMD,
		// and Intel atm.
		panic("could not parse version string \"" + version + "\"")
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
