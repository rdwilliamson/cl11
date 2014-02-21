package cl11

import (
	"fmt"
	"reflect"
	"unsafe"

	clw "github.com/rdwilliamson/clw11"
)

type Kernel struct {
	id            clw.Kernel
	Arguments     int
	FunctionName  string
	Context       *Context
	Program       *Program
	WorkGroupInfo []KernelWorkGroupInfo
	argScratch    [][scratchSize]byte
}

type KernelWorkGroupInfo struct {
	Device                         *Device
	WorkGroupSize                  int
	CompileWorkGroupSize           [3]int
	LocalMemSize                   int
	PreferredWorkGroupSizeMultiple int
	PrivateMemSize                 int
}

type LocalSpaceArg int

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

	k.Arguments = int(k.getUint(clw.KernelNumArgs))
	k.argScratch = make([][8]byte, k.Arguments)

	for i := range k.WorkGroupInfo {
		wgi := &k.WorkGroupInfo[i]

		wgi.WorkGroupSize = int(k.getWorkGroupSize(wgi.Device, clw.KernelWorkGroupSize))
		wgi.LocalMemSize = int(k.getWorkGroupUlong(wgi.Device, clw.KernelLocalMemSize))
		wgi.PreferredWorkGroupSizeMultiple =
			int(k.getWorkGroupSize(wgi.Device, clw.KernelPreferredWorkGroupSizeMultiple))
		wgi.PrivateMemSize = int(k.getWorkGroupUlong(wgi.Device, clw.KernelPrivateMemSize))

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

func (k *Kernel) SetArgument(index int, arg interface{}) error {

	var size uintptr
	var pointer unsafe.Pointer

	switch v := arg.(type) {

	case *Buffer:
		pointer = unsafe.Pointer(&v.id)
		size = unsafe.Sizeof(v.id)

	case LocalSpaceArg:
		pointer = nil
		size = uintptr(v)

	default:
		value := reflect.ValueOf(arg)
		kind := value.Kind()
		for kind == reflect.Ptr || kind == reflect.Interface {
			value = value.Elem()
			kind = value.Kind()
		}

		pointer = unsafe.Pointer(&k.argScratch[index][0])

		switch kind {

		case reflect.Bool:

			localCopy := reflect.NewAt(int32Type, pointer).Elem()
			if value.Bool() {
				localCopy.SetInt(1)
			} else {
				localCopy.SetInt(0)
			}

			size = int32Size

		case reflect.Int:

			localCopy := reflect.NewAt(int32Type, pointer).Elem()
			localCopy.SetInt(value.Int())

			size = int32Size

		case reflect.Uint:

			localCopy := reflect.NewAt(uint32Type, pointer).Elem()
			localCopy.SetUint(value.Uint())

			size = uint32Size

		case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint8, reflect.Uint16, reflect.Uint32,
			reflect.Uint64, reflect.Float32, reflect.Float64:

			localType := value.Type()
			localCopy := reflect.NewAt(localType, pointer).Elem()
			localCopy.Set(value)

			size = localType.Size()

		default:
			return wrapError(fmt.Errorf("invaild argument kind: %s", kind.String()))
		}
	}

	return clw.SetKernelArg(k.id, clw.Uint(index), clw.Size(size), pointer)
}

func (cq *CommandQueue) EnqueueNDRangeKernel(k *Kernel, globalOffset, globalSize, localSize []int,
	waitList []*Event, e *Event) error {

	var event *clw.Event
	if e != nil {
		event = &e.id
		e.Context = cq.Context
		e.CommandType = CommandNDRangeKernel
		e.CommandQueue = cq
	}

	dims := len(globalSize)
	sizes := make([]clw.Size, dims*3)
	if globalOffset != nil {
		for i := 0; i < dims; i++ {
			sizes[i] = clw.Size(globalOffset[i])
		}
	}
	for i := 0; i < dims; i++ {
		sizes[dims+i] = clw.Size(globalSize[i])
		sizes[2*dims+i] = clw.Size(localSize[i])
	}

	return clw.EnqueueNDRangeKernel(cq.id, k.id, sizes[:dims], sizes[dims:2*dims], sizes[2*dims:],
		cq.toEvents(waitList), event)
}
