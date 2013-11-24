package cl11

import (
	"errors"
	"fmt"
	clw "github.com/rdwilliamson/clw11"
	"io"
	"strings"
)

type Platform struct {
	id         clw.PlatformID
	profile    PlatformProfile
	version    PlatformVersion
	name       string
	vendor     string
	extensions []string
}

type PlatformProfile int8

const (
	UnsetProfile PlatformProfile = iota

	FullProfile     PlatformProfile = iota
	EmbeddedProfile PlatformProfile = iota

	UnknownProfile PlatformProfile = iota
)

type PlatformVersion struct {
	Major int8
	Minor int8
	Info  string
}

// Get all the platforms on the system.
func GetPlatforms() ([]Platform, error) {

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

	platforms := make([]Platform, len(platformIDs))
	for i := range platforms {
		platforms[i].id = platformIDs[i]
	}

	return platforms, nil
}

// Retrieve all platform information. Required if your going to print platform
// out or avoid panics in the convenience calls for info.
func (p *Platform) GetInfo() error {
	err := p.getProfile()
	if err != nil {
		return err
	}
	err = p.getVersion()
	if err != nil {
		return err
	}
	err = p.getName()
	if err != nil {
		return err
	}
	err = p.getVendor()
	if err != nil {
		return err
	}
	err = p.getExtensions()
	if err != nil {
		return err
	}
	return nil
}

// Prints platform info in the form of "vendor, name, version, profile,
// extensions.""
func (p Platform) String() string {
	return fmt.Sprintf("%s, %s, %s, %s, %v", p.vendor, p.name, p.version.String(), p.profile.String(), p.extensions)
}

func (p *Platform) getPlatformInfo(paramName clw.PlatformInfo) (string, error) {

	var paramValueSize clw.Size
	err := clw.GetPlatformInfo(p.id, paramName, 0, nil, &paramValueSize)
	if err != nil {
		return "", err
	}

	buffer := make([]byte, paramValueSize)
	err = clw.GetPlatformInfo(p.id, paramName, paramValueSize, buffer, nil)
	if err != nil {
		return "", err
	}

	// Trim space and trailing \0.
	return strings.TrimSpace(string(buffer[:len(buffer)-1])), nil
}

func (p *Platform) getProfile() error {
	profile, err := p.getPlatformInfo(clw.PlatformProfile)
	if err != nil {
		return err
	}
	switch profile {
	case "FULL_PROFILE":
		p.profile = FullProfile
	case "EMBEDDED_PROFILE":
		p.profile = EmbeddedProfile
	default:
		p.profile = UnknownProfile
		return errors.New("unknown platform profile")
	}
	return nil
}

func (p *Platform) Profile() PlatformProfile {

	if p.profile == UnsetProfile {
		err := p.getProfile()
		if err != nil {
			panic(err)
		}
	}

	return p.profile
}

func (pp PlatformProfile) String() string {
	switch pp {
	case UnsetProfile:
		return "profile not yet quiered"
	case FullProfile:
		return "full profile"
	case EmbeddedProfile:
		return "embedded profile"
	case UnknownProfile:
		return "unknown profile"
	}
	panic("unreachable")
}

func (p *Platform) getVersion() error {
	version, err := p.getPlatformInfo(clw.PlatformVersion)
	if err != nil {
		return err
	}
	n, err := fmt.Sscanf(version, "OpenCL %d.%d %s", &p.version.Major, &p.version.Minor, &p.version.Info)

	// May encounter EOF and only scan 2 items if there is no "info".
	if err == io.EOF && n == 2 {
		return nil
	}

	if err != nil {
		return err
	}
	if n != 3 {
		return errors.New("could not parse OpenCL platform version")
	}

	return nil
}

func (p *Platform) Version() PlatformVersion {

	if p.version.Major == 0 {
		err := p.getVersion()
		if err != nil {
			panic(err)
		}
	}

	return p.version
}

func (pv PlatformVersion) String() string {
	if pv.Info != "" {
		return fmt.Sprint(pv.Major, ".", pv.Minor, " ", pv.Info)
	}
	return fmt.Sprint(pv.Major, ".", pv.Minor)
}

func (p *Platform) getName() error {
	var err error
	p.name, err = p.getPlatformInfo(clw.PlatformName)
	if err != nil {
		return err
	}
	return nil
}

func (p *Platform) Name() string {

	if p.name == "" {
		err := p.getName()
		if err != nil {
			panic(err)
		}
	}

	return p.name
}

func (p *Platform) getVendor() error {
	var err error
	p.vendor, err = p.getPlatformInfo(clw.PlatformVendor)
	if err != nil {
		return err
	}
	return nil
}

func (p *Platform) Vendor() string {

	if p.vendor == "" {
		err := p.getVendor()
		if err != nil {
			panic(err)
		}
	}

	return p.vendor
}

func (p *Platform) getExtensions() error {
	extensions, err := p.getPlatformInfo(clw.PlatformExtensions)
	if err != nil {
		return err
	}
	p.extensions = strings.Split(extensions, " ")
	return nil
}

func (p *Platform) Extensions() []string {

	if p.extensions == nil {
		err := p.getExtensions()
		if err != nil {
			panic(err)
		}
	}

	return p.extensions
}

func (p *Platform) HasExtension(extension string) bool {

	if p.extensions == nil {
		err := p.getExtensions()
		if err != nil {
			panic(err)
		}
	}

	for i := range p.extensions {
		if p.extensions[i] == extension {
			return true
		}
	}
	return false
}
