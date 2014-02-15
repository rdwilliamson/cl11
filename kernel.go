package cl11

import (
	"unsafe"

	clw "github.com/rdwilliamson/clw11"
)

type Kernel struct {
	id            clw.Kernel
	Arguemnts     int
	FunctionName  string
	Context       *Context
	Program       *Program
	WorkGroupInfo []KernelWorkGroupInfo
}

type KernelWorkGroupInfo struct {
	Device                         *Device
	WorkGroupSize                  int
	CompileWorkGroupSize           [3]int
	LocalMemorySize                int
	PreferredWorkGroupSizeMultiple int
	PrivateMemorySize              int
}

func (p *Program) CreateKernel(name string) (*Kernel, error) {

	kernel, err := clw.CreateKernel(p.id, name)
	if err != nil {
		return nil, err
	}

	k := &Kernel{id: kernel, FunctionName: name, Context: p.Context, Program: p,
		WorkGroupInfo: make([]KernelWorkGroupInfo, len(p.Devices))}
	for i := range p.Devices {
		k.WorkGroupInfo[i].Device = p.Devices[i]
	}

	err = k.getAllInfo()
	if err != nil {
		return nil, err
	}

	return k, nil
}

func (k *Kernel) getAllInfo() (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()

	k.Arguemnts = int(k.getUint(clw.KernelNumArgs))

	for i := range k.WorkGroupInfo {
		wgi := &k.WorkGroupInfo[i]

		wgi.WorkGroupSize = int(k.getWorkGroupSize(wgi.Device, clw.KernelWorkGroupSize))
		wgi.LocalMemorySize = int(k.getWorkGroupUlong(wgi.Device, clw.KernelLocalMemSize))
		wgi.PreferredWorkGroupSizeMultiple =
			int(k.getWorkGroupSize(wgi.Device, clw.KernelPreferredWorkGroupSizeMultiple))
		wgi.PrivateMemorySize = int(k.getWorkGroupUlong(wgi.Device, clw.KernelPrivateMemSize))

		cwgs := k.getWorkGroupSize3(wgi.Device, clw.KernelCompileWorkGroupSize)
		wgi.CompileWorkGroupSize[0] = int(cwgs[0])
		wgi.CompileWorkGroupSize[1] = int(cwgs[1])
		wgi.CompileWorkGroupSize[2] = int(cwgs[2])
	}

	return
}

func (k *Kernel) getUint(paramName clw.KernelInfo) clw.Uint {
	var param clw.Uint
	err := clw.GetKernelInfo(k.id, clw.KernelNumArgs, clw.Size(unsafe.Sizeof(param)), unsafe.Pointer(&param), nil)
	if err != nil {
		panic(err)
	}
	return param
}

func (k *Kernel) getWorkGroupSize(d *Device, paramName clw.KernelWorkGroupInfo) clw.Size {
	var param clw.Size
	err := clw.GetKernelWorkGroupInfo(k.id, d.id, paramName, clw.Size(unsafe.Sizeof(param)), unsafe.Pointer(&param),
		nil)
	if err != nil {
		panic(err)
	}
	return param
}

func (k *Kernel) getWorkGroupSize3(d *Device, paramName clw.KernelWorkGroupInfo) [3]clw.Size {
	var param [3]clw.Size
	err := clw.GetKernelWorkGroupInfo(k.id, d.id, paramName, clw.Size(unsafe.Sizeof(param)), unsafe.Pointer(&param),
		nil)
	if err != nil {
		panic(err)
	}
	return param
}

func (k *Kernel) getWorkGroupUlong(d *Device, paramName clw.KernelWorkGroupInfo) clw.Ulong {
	var param clw.Ulong
	err := clw.GetKernelWorkGroupInfo(k.id, d.id, paramName, clw.Size(unsafe.Sizeof(param)), unsafe.Pointer(&param),
		nil)
	if err != nil {
		panic(err)
	}
	return param
}
