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

	return &Program{id: program, Context: c}, nil
}

// Temp, figure out how to populate devices from program.
func (p *Program) GetDevices() ([]*Device, error) {

	var numDevices clw.Uint
	err := clw.GetProgramInfo(p.id, clw.ProgramNumDevices, clw.Size(unsafe.Sizeof(numDevices)),
		unsafe.Pointer(&numDevices), nil)
	if err != nil {
		return nil, err
	}

	deviceIDs := make([]clw.DeviceID, numDevices)
	err = clw.GetProgramInfo(p.id, clw.ProgramDevices, clw.Size(unsafe.Sizeof(deviceIDs[0])*uintptr(numDevices)),
		unsafe.Pointer(&deviceIDs[0]), nil)
	if err != nil {
		return nil, err
	}

	devices := make([]*Device, len(deviceIDs))
	for i := range devices {

		device := &Device{id: deviceIDs[i]}

		err = device.getAllInfo()
		if err != nil {
			return nil, err
		}

		devices[i] = device
	}

	return devices, nil
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
