package cl11

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"unsafe"

	clw "github.com/rdwilliamson/clw11"
)

type Platform struct {
	id         clw.PlatformID
	Devices    []*Device
	Profile    PlatformProfile
	Version    PlatformVersion
	Name       string
	Vendor     string
	Extensions []string
}

type PlatformProfile int

const (
	FullProfile     PlatformProfile = iota
	EmbeddedProfile PlatformProfile = iota
)

type PlatformVersion struct {
	Major int
	Minor int
	Info  string
}

// Get all the platforms on the system.
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

		platforms[i] = &Platform{id: platformIDs[i]}

		err = platforms[i].getAllInfo()
		if err != nil {
			return nil, err
		}

		platforms[i].Devices, err = platforms[i].GetDevices()
		if err != nil {
			return nil, err
		}
	}

	return platforms, nil
}

func (p *Platform) getAllInfo() (err error) {

	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()

	p.getProfile()
	p.getVersion()
	p.Name = p.getString(clw.PlatformName)
	p.Vendor = p.getString(clw.PlatformVendor)
	p.Extensions = strings.Split(p.getString(clw.PlatformExtensions), " ")

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

func (p *Platform) getProfile() {
	profile := p.getString(clw.PlatformProfile)
	switch profile {
	case "FULL_PROFILE":
		p.Profile = FullProfile
	case "EMBEDDED_PROFILE":
		p.Profile = EmbeddedProfile
	default:
		panic(errors.New("unknown platform profile"))
	}
}

func (pp PlatformProfile) String() string {
	switch pp {
	case FullProfile:
		return "full profile"
	case EmbeddedProfile:
		return "embedded profile"
	}
	panic("unreachable")
}

func (p *Platform) getVersion() {

	version := p.getString(clw.PlatformVersion)
	n, err := fmt.Sscanf(version, "OpenCL %d.%d %s", &p.Version.Major, &p.Version.Minor, &p.Version.Info)

	// May encounter EOF and only scan 2 items if there is no "info".
	if err == io.EOF && n == 2 {
		return
	}

	if err != nil {
		panic(err)
	}
	if n != 3 {
		panic(errors.New("could not parse OpenCL platform version"))
	}
}

func (pv PlatformVersion) String() string {
	if pv.Info != "" {
		return fmt.Sprint(pv.Major, ".", pv.Minor, " ", pv.Info)
	}
	return fmt.Sprint(pv.Major, ".", pv.Minor)
}

func (p *Platform) HasExtension(extension string) bool {
	for i := range p.Extensions {
		if p.Extensions[i] == extension {
			return true
		}
	}
	return false
}
