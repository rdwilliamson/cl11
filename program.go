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

	err := clw.BuildProgram(p.id, devices, options, nil)

	if err == nil {
		p.Devices = d
	}

	return err
}

func (p *Program) GetProgramBinaries() error {

	devices := make([]clw.DeviceID, len(p.Devices))
	err := clw.GetProgramInfo(p.id, clw.ProgramDevices, clw.Size(unsafe.Sizeof(devices[0])*uintptr(len(devices))),
		unsafe.Pointer(&devices[0]), nil)
	if err != nil {
		return err
	}

	// Reorder p.Devices to match the order returned by GetProgramInfo.
	for i := 0; i < len(devices)-1; i++ {
		j := 0
		for ; j < len(p.Devices); j++ {
			if p.Devices[j].id == devices[i] {
				break
			}
		}
		if i != j {
			p.Devices[i], p.Devices[j] = p.Devices[j], p.Devices[i]
		}
	}

	sizes := make([]clw.Size, len(p.Devices))
	err = clw.GetProgramInfo(p.id, clw.ProgramBinarySizes, clw.Size(unsafe.Sizeof(sizes[0])*uintptr(len(sizes))),
		unsafe.Pointer(&sizes[0]), nil)
	if err != nil {
		return err
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
		return err
	}

	// TODO what to actually do with binaries.

	return nil
}
