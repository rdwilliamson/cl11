package cl11

import (
	"unsafe"

	clw "github.com/rdwilliamson/clw11"
)

// An OpenCL context is created with one or more devices. Contexts are used by
// the OpenCL runtime for managing objects such as command-queues, memory,
// program and kernel objects and for executing kernels on one or more devices
// specified in the context.
type Context struct {
	id clw.Context

	// The devices in the context.
	Devices []*Device

	// The properties used to create the context.
	Properties []ContextProperties
}

type ContextProperties clw.ContextProperties

const (
	ContextPlatform = ContextProperties(clw.ContextPlatform)
)

type ContextCallback func(err string, data []byte, userData interface{})

// Ensure the definitions of context callbacks match. ContextCallback isn't an
// alias for clw.ContextCallbackFunc just so the user doesn't have to look into
// clw.
var _ clw.ContextCallbackFunc = clw.ContextCallbackFunc(ContextCallback(func(err string, data []byte,
	userData interface{}) {
}))

// Creates an OpenCL context.
//
// An OpenCL context is created with one or more devices. Contexts are used by
// the OpenCL runtime for managing objects such as command-queues, memory,
// program and kernel objects and for executing kernels on one or more devices
// specified in the context.
func CreateContext(d []*Device, cp []ContextProperties, cc ContextCallback, userData interface{}) (*Context, error) {

	devices := make([]clw.DeviceID, len(d))
	for i := range d {
		devices[i] = clw.DeviceID(d[i].id)
	}

	// "Convert" type and appending terminating zero.
	var properties []clw.ContextProperties
	if cp != nil {
		properties = make([]clw.ContextProperties, len(cp)+1)
		for i := range cp {
			properties[i] = clw.ContextProperties(cp[i])
		}
		properties[len(cp)] = 0
	}

	context, err := clw.CreateContext(properties, devices, clw.ContextCallbackFunc(cc), userData)
	if err != nil {
		return nil, err
	}

	return &Context{id: context, Devices: d, Properties: cp}, nil
}

// Create an OpenCL context from a device type that identifies the specific
// device(s) to use.
func CreateContextFromType(cp []ContextProperties, dt DeviceType, cc ContextCallback,
	userData interface{}) (*Context, error) {

	// "Convert" type and appending terminating zero.
	var properties []clw.ContextProperties
	if cp != nil {
		properties = make([]clw.ContextProperties, len(cp)+1)
		for i := range cp {
			properties[i] = clw.ContextProperties(cp[i])
		}
		properties[len(cp)] = 0
	}

	context, err := clw.CreateContextFromType(properties, clw.DeviceType(dt), clw.ContextCallbackFunc(cc), userData)
	if err != nil {
		return nil, err
	}

	// Get devices.
	var numDevices clw.Uint
	err = clw.GetContextInfo(context, clw.ContextNumDevices, clw.Size(unsafe.Sizeof(numDevices)),
		unsafe.Pointer(&numDevices), nil)
	if err != nil {
		return nil, err
	}

	var devicePtrs []*Device
	if numDevices > 0 {
		devices := make([]clw.DeviceID, int(numDevices))
		devicePtrs = make([]*Device, int(numDevices))

		err = clw.GetContextInfo(context, clw.ContextDevices, clw.Size(uintptr(numDevices)*unsafe.Sizeof(devices[0])),
			unsafe.Pointer(&devices[0]), nil)
		if err != nil {
			return nil, err
		}

		for i := range devices {
			d := &Device{id: devices[i]}
			err = d.getAllInfo()
			if err != nil {
				return nil, err
			}
			devicePtrs[i] = d
		}
	}

	return &Context{context, devicePtrs, cp}, nil
}

// Increment the context reference count.
func (c *Context) Retain() error {
	return clw.RetainContext(c.id)
}

// Decrement the context reference count.
func (c *Context) Release() error {
	return clw.ReleaseContext(c.id)
}

// The context reference count.
//
// The reference count returned should be considered immediately stale. It is
// unsuitable for general use in applications. This feature is provided for
// identifying memory leaks.
func (c *Context) ReferenceCount() (int, error) {

	var count clw.Uint
	err := clw.GetContextInfo(c.id, clw.ContextReferenceCount, clw.Size(unsafe.Sizeof(count)), unsafe.Pointer(&count),
		nil)
	if err != nil {
		return 0, err
	}

	return int(count), nil
}
