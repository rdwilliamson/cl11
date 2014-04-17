package cl11

import (
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

func (p *Program) Build(d []*Device, options string) error {
	// TODO callback

	devices := make([]clw.DeviceID, len(d))
	for i := range d {
		devices[i] = d[i].id
	}

	err := clw.BuildProgram(p.id, devices, options, nil, nil)

	if err == nil {
		p.Devices = d
	}

	return err
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
