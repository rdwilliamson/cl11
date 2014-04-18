package cl11

import (
	"strings"
	"unsafe"

	clw "github.com/rdwilliamson/clw11"
)

// A program object created from either source or a binary.
type Program struct {
	id clw.Program

	// The context the program is associated with.
	Context *Context

	// The devices associated with the program.
	Devices []*Device

	// The build options specified during Build.
	Options string
}

// A program binary for a device.
type ProgramBinary struct {

	// The program from which the binary was compiled.
	Program *Program

	// The device for which the binary was compiled.
	Device *Device

	// The binary. Can be an implementation specific intermediate
	// representation, device specific bits, or both.
	Binary []byte
}

type BuildStatus clw.BuildStatus

const (
	BuildSuccess    = BuildStatus(clw.BuildSuccess)
	BuildNone       = BuildStatus(clw.BuildNone)
	BuildError      = BuildStatus(clw.BuildError)
	BuildInProgress = BuildStatus(clw.BuildInProgress)
)

type ProgramCallback func(p *Program, userData interface{})

func (bs BuildStatus) String() string {
	switch bs {
	case BuildSuccess:
		return "success"
	case BuildNone:
		return "none"
	case BuildError:
		return "error"
	case BuildInProgress:
		return "in progress"
	}
	return ""
}

// Create a program object for a context with the specified source.
//
// The devices associated with the program are initially all the context
// devices.
func (c *Context) CreateProgramWithSource(sources ...[]byte) (*Program, error) {

	program, err := clw.CreateProgramWithSource(c.id, sources)
	if err != nil {
		return nil, err
	}

	return &Program{id: program, Context: c, Devices: c.Devices}, nil
}

// Creates a program object for a context for the specified devices with the
// passed binaries.
//
// Length of devices, binaries, and status must match (though status can be
// nil). The status will contain an error, or not, for each device.
func (c *Context) CreateProgramWithBinary(d []*Device, binaries [][]byte, status []error) (*Program, error) {

	deviceIDs := make([]clw.DeviceID, len(d))
	for i := range deviceIDs {
		deviceIDs[i] = d[i].id
	}

	program, err := clw.CreateProgramWithBinary(c.id, deviceIDs, binaries, status)
	if err != nil {
		return nil, err
	}

	return &Program{id: program, Context: c, Devices: d}, nil
}

// Builds (compiles and links) a program executable from the program source or
// binary.
//
// The build options are categorized as pre-processor options, options for math
// intrinsics, options that control optimization and miscellaneous options. This
// specification defines a standard set of options that must be supported by an
// OpenCL compiler when building program executables online or offline. These
// may be extended by a set of vendor- or platform-specific options.
//
// The callback is optional. If it is supplied Build will return immediately, if
// it isn't then Build will block until the program has been build (successfully
// or unsuccessfully).
func (p *Program) Build(d []*Device, options string, pc ProgramCallback, userData interface{}) error {

	devices := make([]clw.DeviceID, len(d))
	for i := range d {
		devices[i] = d[i].id
	}

	var callback clw.ProgramCallbackFunc
	if pc != nil {
		callback = func(programID clw.Program, userData interface{}) {
			pc(p, userData)
		}
	}

	p.Options = options

	return clw.BuildProgram(p.id, devices, options, callback, userData)
}

// Returns the binaries for each device associated with the program.
//
// The binaries can be the source of CreateProgramWithBinary or the result of
// Build from either source or binaries. The returned bits can be an
// implementation specific intermediate representation, device specific bits, or
// both.
func (p *Program) GetProgramBinaries() ([]ProgramBinary, error) {

	devices := make([]clw.DeviceID, len(p.Devices))
	err := clw.GetProgramInfo(p.id, clw.ProgramDevices, clw.Size(unsafe.Sizeof(devices[0])*uintptr(len(devices))),
		unsafe.Pointer(&devices[0]), nil)
	if err != nil {
		return nil, err
	}

	sizes := make([]clw.Size, len(p.Devices))
	err = clw.GetProgramInfo(p.id, clw.ProgramBinarySizes, clw.Size(unsafe.Sizeof(sizes[0])*uintptr(len(sizes))),
		unsafe.Pointer(&sizes[0]), nil)
	if err != nil {
		return nil, err
	}

	binaries := make([][]byte, len(p.Devices))
	binaryPointers := make([]unsafe.Pointer, len(p.Devices))
	for i := range binaries {
		binaries[i] = make([]byte, int(sizes[i]))
		binaryPointers[i] = unsafe.Pointer(&binaries[i][0])
	}

	err = clw.GetProgramInfo(p.id, clw.ProgramBinaries,
		clw.Size(unsafe.Sizeof(binaryPointers[0])*uintptr(len(binaryPointers))), unsafe.Pointer(&binaryPointers[0]),
		nil)
	if err != nil {
		return nil, err
	}

	programBinaries := make([]ProgramBinary, len(p.Devices))
	for i := range programBinaries {
		programBinary := &programBinaries[i]

		programBinary.Program = p

		for j := range p.Devices {
			device := p.Devices[j]

			if device.id == devices[j] {
				programBinary.Device = device
				programBinary.Binary = binaries[j]
				break
			}
		}
	}

	return programBinaries, nil
}

// Returns build status for the device in the program object.
func (p *Program) BuildStatus(d *Device) (BuildStatus, error) {
	var buildStatus clw.BuildStatus
	err := clw.GetProgramBuildInfo(p.id, d.id, clw.ProgramBuildStatusInfo, clw.Size(unsafe.Sizeof(buildStatus)),
		unsafe.Pointer(&buildStatus), nil)
	return BuildStatus(buildStatus), err
}

// Returns build log for the device in the program object.
func (p *Program) BuildLog(d *Device) (string, error) {

	var paramValueSize clw.Size
	err := clw.GetProgramBuildInfo(p.id, d.id, clw.ProgramBuildLog, 0, nil, &paramValueSize)
	if err != nil {
		return "", err
	}

	buffer := make([]byte, paramValueSize)
	err = clw.GetProgramBuildInfo(p.id, d.id, clw.ProgramBuildLog, paramValueSize, unsafe.Pointer(&buffer[0]), nil)
	if err != nil {
		return "", err
	}

	// Trim space and trailing \0.
	return strings.TrimSpace(string(buffer[:len(buffer)-1])), nil
}

// Increments the program reference count.
//
// The OpenCL commands that return a program perform an implicit retain.
func (p *Program) Retain() error {
	return clw.RetainProgram(p.id)
}

// Decrements the program reference count.
//
// The program object is deleted after all kernel objects associated with
// program have been deleted and the program reference count becomes zero.
func (p *Program) Release() error {
	return clw.ReleaseProgram(p.id)
}

// The program reference count.
//
// The reference count returned should be considered immediately stale. It is
// unsuitable for general use in applications. This feature is provided for
// identifying memory leaks.
func (p *Program) ReferenceCount() (int, error) {
	var param clw.Uint
	err := clw.GetProgramInfo(p.id, clw.ProgramReferenceCount, clw.Size(unsafe.Sizeof(param)), unsafe.Pointer(&param),
		nil)
	return int(param), err
}

// Allows the implementation to release the resources allocated by the OpenCL
// compiler.
//
// This is a hint from the application and does not guarantee that the compiler
// will not be used in the future or that the compiler will actually be unloaded
// by the implementation. Calls to clBuildProgram after clUnloadCompiler will
// reload the compiler, if necessary, to build the appropriate program
// executable.
func UnloadCompiler() error {
	return clw.UnloadCompiler()
}
