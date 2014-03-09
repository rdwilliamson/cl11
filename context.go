package cl11

import clw "github.com/rdwilliamson/clw11"

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

type ContextCallback clw.ContextCallbackFunc

func CreateContext(d []*Device, cp []ContextProperties, cc ContextCallback, userData interface{}) (*Context, error) {

	devices := make([]clw.DeviceID, len(d))
	for i := range d {
		devices[i] = clw.DeviceID(d[i].id)
	}

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
