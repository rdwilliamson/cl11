package cl11

import (
	"strings"
	"unsafe"

	clw "github.com/rdwilliamson/clw11"
)

// Platforms implement specific features and allow applications to query
// devices.
type Platform struct {
	id clw.PlatformID

	// The profile name supported by the implementation, either the full profile
	// or a subset of each OpenCL version (embedded profile).
	Profile Profile

	// A Platform's version.
	Version Version

	// Platform's name.
	Name string

	// Vendor's name.
	Vendor string

	// Extensions supported by all devices associated with the platform.
	Extensions []string

	// Devices available on the platform.
	Devices []*Device
}

const (
	zeroProfile     Profile = iota
	FullProfile     Profile = iota
	EmbeddedProfile Profile = iota
)

// Obtain the list of platforms and their information.
func GetPlatforms() ([]*Platform, error) {

	var numPlatforms clw.Uint
	err := clw.GetPlatformIDs(0, nil, &numPlatforms)
	if err != nil {
		return nil, err
	}

	platformIDs := make([]clw.PlatformID, numPlatforms)
	err = clw.GetPlatformIDs(numPlatforms, &platformIDs[0], nil)
	if err != nil {
		return nil, err
	}

	platforms := make([]*Platform, len(platformIDs))
	for i := range platforms {

		platform := &Platform{id: platformIDs[i]}

		err = platform.getAllInfo()
		if err != nil {
			return nil, err
		}

		err = platform.getDevices()
		if err != nil {
			return nil, err
		}

		platforms[i] = platform
	}

	return platforms, nil
}

func (p *Platform) getAllInfo() (err error) {

	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()

	p.Profile = toProfile(p.getString(clw.PlatformProfile))
	p.Version = toVersion(p.getString(clw.PlatformVersion))
	p.Name = p.getString(clw.PlatformName)
	p.Vendor = p.getString(clw.PlatformVendor)
	p.Extensions = strings.Fields(p.getString(clw.PlatformExtensions))

	return
}

func (p *Platform) getString(paramName clw.PlatformInfo) string {

	var paramValueSize clw.Size
	err := clw.GetPlatformInfo(p.id, paramName, 0, nil, &paramValueSize)
	if err != nil {
		panic(err)
	}

	buffer := make([]byte, paramValueSize)
	err = clw.GetPlatformInfo(p.id, paramName, paramValueSize, unsafe.Pointer(&buffer[0]), nil)
	if err != nil {
		panic(err)
	}

	// Trim space and trailing \0.
	return strings.TrimSpace(string(buffer[:len(buffer)-1]))
}

// Check if the platform supports the extension.
func (p *Platform) HasExtension(extension string) bool {

	for _, v := range p.Extensions {
		if v == extension {
			return true
		}
	}

	return false
}

// Returns the platform as a ContextProperties suitable for adding to a property
// list during context creation.
func (p *Platform) ToContextProperty() ContextProperties {
	return *(*ContextProperties)(unsafe.Pointer(&p.id))
}
