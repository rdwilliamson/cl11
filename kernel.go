package cl11

import (
	"fmt"
	"reflect"
	"unsafe"

	clw "github.com/rdwilliamson/clw11"
)

// A kernel is a function declared in a program.
type Kernel struct {
	id clw.Kernel

	// The number of arguments to kernel.
	Arguments int

	// The kernel function name.
	FunctionName string

	// The context associated with the kernel.
	Context *Context

	// The program associated with the kernel.
	Program *Program

	// Information about the kernel object that may be specific to a device.
	WorkGroupInfo []KernelWorkGroupInfo

	// Scratch space to store arugment.
	argScratch [][scratchSize]byte
}

// Information about the kernel object specific to a device.
type KernelWorkGroupInfo struct {

	// The device associated with the kernel.
	Device *Device

	// The maximum workgroup size that can be used by the kernel to execute on
	// the device.
	WorkGroupSize int

	// The work group size specified in the function qualifiers or all zeros.
	CompileWorkGroupSize [3]int

	// The amount of local memory used by the kernel.
	LocalMemSize int

	// The preferred multiple of workgroup size for launch. This is a
	// performance hint.
	PreferredWorkGroupSizeMultiple int

	// The minimum amount of private memory, in bytes, used by each workitem in
	// the kernel.
	PrivateMemSize int
}

type LocalSpaceArg int

// Creates a kernal object.
//
// A kernel is a function declared in a program. A kernel is identified by the
// __kernel qualifier applied to any function in a program. A kernel object
// encapsulates the specific __kernel function declared in a program and the
// argument values to be used when executing this __kernel function.
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

	k.Arguments = k.getUint(clw.KernelNumArgs)
	k.argScratch = make([][8]byte, k.Arguments)

	for i := range k.WorkGroupInfo {
		wgi := &k.WorkGroupInfo[i]

		wgi.WorkGroupSize = k.getWorkGroupSize(wgi.Device, clw.KernelWorkGroupSize)
		wgi.LocalMemSize = k.getWorkGroupUlong(wgi.Device, clw.KernelLocalMemSize)
		wgi.PreferredWorkGroupSizeMultiple =
			k.getWorkGroupSize(wgi.Device, clw.KernelPreferredWorkGroupSizeMultiple)
		wgi.PrivateMemSize = k.getWorkGroupUlong(wgi.Device, clw.KernelPrivateMemSize)
		wgi.CompileWorkGroupSize = k.getWorkGroupSize3(wgi.Device, clw.KernelCompileWorkGroupSize)
	}

	return
}

func (k *Kernel) getUint(paramName clw.KernelInfo) int {
	var param clw.Uint
	err := clw.GetKernelInfo(k.id, clw.KernelNumArgs, clw.Size(unsafe.Sizeof(param)), unsafe.Pointer(&param), nil)
	if err != nil {
		panic(err)
	}
	return int(param)
}

func (k *Kernel) getWorkGroupSize(d *Device, paramName clw.KernelWorkGroupInfo) int {
	var param clw.Size
	err := clw.GetKernelWorkGroupInfo(k.id, d.id, paramName, clw.Size(unsafe.Sizeof(param)), unsafe.Pointer(&param),
		nil)
	if err != nil {
		panic(err)
	}
	return int(param)
}

func (k *Kernel) getWorkGroupSize3(d *Device, paramName clw.KernelWorkGroupInfo) [3]int {
	var param [3]clw.Size
	err := clw.GetKernelWorkGroupInfo(k.id, d.id, paramName, clw.Size(unsafe.Sizeof(param)), unsafe.Pointer(&param),
		nil)
	if err != nil {
		panic(err)
	}
	return [3]int{int(param[0]), int(param[1]), int(param[2])}
}

func (k *Kernel) getWorkGroupUlong(d *Device, paramName clw.KernelWorkGroupInfo) int {
	var param clw.Ulong
	err := clw.GetKernelWorkGroupInfo(k.id, d.id, paramName, clw.Size(unsafe.Sizeof(param)), unsafe.Pointer(&param),
		nil)
	if err != nil {
		panic(err)
	}
	return int(param)
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
				localCopy.SetInt(int64(clw.True))
			} else {
				localCopy.SetInt(int64(clw.False))
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

func (k *Kernel) SetArguments(args ...interface{}) error {

	if len(args) != k.Arguments {
		return wrapError(fmt.Errorf("invalid number of arguments: expecting %d, got %d", k.Arguments, len(args)))
	}

	for i := range args {
		err := k.SetArgument(i, args[i])
		if err != nil {
			return wrapError(fmt.Errorf("setting argument %d: %s", i, err.Error()))
		}
	}

	return nil
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
