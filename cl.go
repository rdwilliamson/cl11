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
	float32Type = reflect.TypeOf(float32(0))
	float32Size = float32Type.Size()
)

// A reference counted OpenCL object.
type Object interface {
	Retain() error
	Release() error
	ReferenceCount() (int, error)
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

func pointerSize(value interface{}) (unsafe.Pointer, uintptr, error) {
	v := reflect.ValueOf(value)
	switch v.Kind() {
	// case reflect.Ptr:
	case reflect.Slice:
		pointer := unsafe.Pointer(&(*reflect.SliceHeader)(unsafe.Pointer(v.Pointer())).Data)
		size := v.Type().Elem().Size() * uintptr(v.Len())
		return pointer, size, nil
	}
	return unsafe.Pointer(uintptr(0)), 0, NotAddressable
}
