package cl11

import (
	"reflect"
	"strings"
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
}

// Information about the kernel object specific to a device.
type KernelWorkGroupInfo struct {

	// The device associated with the kernel.
	Device *Device

	// The maximum work group size that can be used by the kernel to execute on
	// the device.
	WorkGroupSize int

	// The work group size specified in the function qualifiers or all zeros.
	CompileWorkGroupSize [3]int

	// The amount of local memory used by the kernel.
	LocalMemSize int

	// The preferred multiple of work group size for launch. This is a
	// performance hint.
	PreferredWorkGroupSizeMultiple int

	// The minimum amount of private memory, in bytes, used by each work item in
	// the kernel.
	PrivateMemSize int
}

// Type passed to SetArguments to allocate the set amount of local memory.
type LocalSpaceArg int

// Creates a kernel object.
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

	k := &Kernel{
		id:            kernel,
		Context:       p.Context,
		Program:       p,
		WorkGroupInfo: make([]KernelWorkGroupInfo, len(p.Devices)),
	}
	for i := range p.Devices {
		k.WorkGroupInfo[i].Device = p.Devices[i]
	}

	err = k.getAllInfo()
	if err != nil {
		return nil, err
	}

	return k, nil
}

// Creates kernel objects for all kernel functions in a program object.
//
// Creates kernel objects for all kernel functions in program. Kernel objects
// are not created for any __kernel functions in program that do not have the
// same function definition across all devices for which a program executable
// has been successfully built.
func (p *Program) CreateKernelsInProgram() ([]*Kernel, error) {

	var numKernels clw.Uint
	err := clw.CreateKernelsInProgram(p.id, nil, &numKernels)
	if err != nil {
		return nil, err
	}

	kernelIDs := make([]clw.Kernel, int(numKernels))
	err = clw.CreateKernelsInProgram(p.id, kernelIDs, nil)
	if err != nil {
		return nil, err
	}

	kernels := make([]*Kernel, int(numKernels))
	for i := range kernels {

		kernels[i] = &Kernel{
			id:            kernelIDs[i],
			Context:       p.Context,
			Program:       p,
			WorkGroupInfo: make([]KernelWorkGroupInfo, len(p.Devices)),
		}
		for j := range p.Devices {
			kernels[i].WorkGroupInfo[j].Device = p.Devices[j]
		}

		err = kernels[i].getAllInfo()
		if err != nil {
			return nil, err
		}
	}

	return kernels, nil
}

func (k *Kernel) getAllInfo() (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()

	k.FunctionName = k.getString(clw.KernelFunctionName)
	k.Arguments = k.getUint(clw.KernelNumArgs)

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

func (k *Kernel) getString(paramName clw.KernelInfo) string {

	var paramValueSize clw.Size
	err := clw.GetKernelInfo(k.id, paramName, 0, nil, &paramValueSize)
	if err != nil {
		panic(err)
	}

	buffer := make([]byte, paramValueSize)
	err = clw.GetKernelInfo(k.id, paramName, paramValueSize, unsafe.Pointer(&buffer[0]), nil)
	if err != nil {
		panic(err)
	}

	// Trim space and trailing \0.
	return strings.TrimSpace(string(buffer[:len(buffer)-1]))
}

// Increments the kernel object reference count.
//
// CreateKernel or CreateKernelsInProgram do an implicit retain.
func (k *Kernel) Retain() error {
	return clw.RetainKernel(k.id)
}

// Decrements the kernel reference count.
//
// The kernel object is deleted once the number of instances that are retained
// to kernel become zero and the kernel object is no longer needed by any
// enqueued commands that use kernel.
func (k *Kernel) Release() error {
	return clw.ReleaseKernel(k.id)
}

// Return the kernel reference count.
//
// The reference count returned should be considered immediately stale. It is
// unsuitable for general use in applications. This feature is provided for
// identifying memory leaks.
func (k *Kernel) ReferenceCount() (int, error) {
	var param clw.Uint
	err := clw.GetKernelInfo(k.id, clw.KernelReferenceCount, clw.Size(unsafe.Sizeof(param)), unsafe.Pointer(&param),
		nil)
	return int(param), err
}

// Set the argument value for a specific argument of a kernel.
//
// All OpenCL API calls are thread-safe except SetArg (and SetArguments), which
// is safe to call from any host thread, and is re-entrant so long as concurrent
// calls operate on different cl_kernel objects.
//
// A kernel object does not update the reference count for objects such as
// memory, sampler objects specified as argument values by clSetKernelArg. Users
// may not rely on a kernel object to retain objects specified as argument
// values to the kernel.
func (k *Kernel) SetArg(index int, arg interface{}) error {

	var size uintptr
	var pointer unsafe.Pointer

	switch v := arg.(type) {

	case *Buffer:
		pointer = unsafe.Pointer(&v.id)
		size = unsafe.Sizeof(v.id)

	case *Image:
		pointer = unsafe.Pointer(&v.id)
		size = unsafe.Sizeof(v.id)

	case LocalSpaceArg:
		pointer = nil
		size = uintptr(v)

	default:

		// Find the underlying type.
		value := reflect.ValueOf(arg)
		kind := value.Kind()
		for kind == reflect.Ptr || kind == reflect.Interface {
			value = value.Elem()
			kind = value.Kind()
		}

		// Create an addressable copy if required.
		if !value.CanAddr() {
			newvalue := reflect.New(value.Type()).Elem()
			newvalue.Set(value)
			value = newvalue
		}

		pointer = unsafe.Pointer(value.UnsafeAddr())
		size = value.Type().Size()
	}

	return clw.SetKernelArg(k.id, clw.Uint(index), clw.Size(size), pointer)
}

// Set all argument values of a kernel.
//
// This is a convenience wrapper around SetArg, consult it for more info.
func (k *Kernel) SetArguments(args ...interface{}) error {
	for i := range args {
		err := k.SetArg(i, args[i])
		if err != nil {
			return err
		}
	}
	return nil
}

// Enqueues a command to execute a kernel on a device.
//
// GlobalOffset is optional, if it omitted it is assumed to be all zeros. The
// dimensions of globalOffset, globalSize, and localSize must match and be less
// than or equal to the max work item dimensions.
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

	events := cq.createEvents(waitList)
	err := clw.EnqueueNDRangeKernel(cq.id, k.id, sizes[:dims], sizes[dims:2*dims], sizes[2*dims:],
		events, event)
	cq.releaseEvents(events)
	return err
}

// Enqueues a command to execute a kernel on a device.
//
// The kernel is executed using a single work-item.
func (cq *CommandQueue) EnqueueTask(k *Kernel, waitList []*Event, e *Event) error {

	var event *clw.Event
	if e != nil {
		event = &e.id
		e.Context = cq.Context
		e.CommandType = CommandTask
		e.CommandQueue = cq
	}

	events := cq.createEvents(waitList)
	err := clw.EnqueueTask(cq.id, k.id, events, event)
	cq.releaseEvents(events)
	return err
}
