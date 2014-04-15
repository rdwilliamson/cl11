package cl11

import (
	"unsafe"

	clw "github.com/rdwilliamson/clw11"
)

type Program struct {
	id      clw.Program
	Context *Context
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

func (c *Context) CreateProgramWithSource(sources ...[]byte) (*Program, error) {

	program, err := clw.CreateProgramWithSource(c.id, sources)
	if err != nil {
		return nil, err
	}

	return &Program{id: program, Context: c}, nil
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
