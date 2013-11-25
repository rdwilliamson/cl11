package cl11

import (
	"errors"
	"fmt"
	clw "github.com/rdwilliamson/clw11"
	"io"
	"strings"
)

type Platform struct {
	ID         clw.PlatformID
	Devices    []Device
	Profile    PlatformProfile
	Version    PlatformVersion
	Name       string
	Vendor     string
	Extensions []string
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
		platforms[i].ID = platformIDs[i]
		err = platforms[i].getAllInfo()
		if err != nil {
			return nil, err
		}
	}

	return platforms, nil
}

// Prints platform info in the form of "vendor, name, version, profile,
// extensions.""
func (p Platform) String() string {
	return fmt.Sprintf("%s, %s, %s, %s, %v", p.Vendor, p.Name, p.Version.String(), p.Profile.String(), p.Extensions)
}

func (p *Platform) getAllInfo() error {
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

func (p *Platform) getInfo(paramName clw.PlatformInfo) (string, error) {

	var paramValueSize clw.Size
	err := clw.GetPlatformInfo(p.ID, paramName, 0, nil, &paramValueSize)
	if err != nil {
		return "", err
	}

	buffer := make([]byte, paramValueSize)
	err = clw.GetPlatformInfo(p.ID, paramName, paramValueSize, clw.Pointer(buffer), nil)
	if err != nil {
		return "", err
	}

	// Trim space and trailing \0.
	return strings.TrimSpace(string(buffer[:len(buffer)-1])), nil
}

func (p *Platform) getProfile() error {
	profile, err := p.getInfo(clw.PlatformProfile)
	if err != nil {
		return err
	}
	switch profile {
	case "FULL_PROFILE":
		p.Profile = FullProfile
	case "EMBEDDED_PROFILE":
		p.Profile = EmbeddedProfile
	default:
		p.Profile = UnknownProfile
		return errors.New("unknown platform profile")
	}
	return nil
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
	version, err := p.getInfo(clw.PlatformVersion)
	if err != nil {
		return err
	}
	n, err := fmt.Sscanf(version, "OpenCL %d.%d %s", &p.Version.Major, &p.Version.Minor, &p.Version.Info)

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

func (pv PlatformVersion) String() string {
	if pv.Info != "" {
		return fmt.Sprint(pv.Major, ".", pv.Minor, " ", pv.Info)
	}
	return fmt.Sprint(pv.Major, ".", pv.Minor)
}

func (p *Platform) getName() error {
	var err error
	p.Name, err = p.getInfo(clw.PlatformName)
	if err != nil {
		return err
	}
	return nil
}

func (p *Platform) getVendor() error {
	var err error
	p.Vendor, err = p.getInfo(clw.PlatformVendor)
	if err != nil {
		return err
	}
	return nil
}

func (p *Platform) getExtensions() error {
	extensions, err := p.getInfo(clw.PlatformExtensions)
	if err != nil {
		return err
	}
	p.Extensions = strings.Split(extensions, " ")
	return nil
}

func (p *Platform) HasExtension(extension string) bool {
	for i := range p.Extensions {
		if p.Extensions[i] == extension {
			return true
		}
	}
	return false
}
