package cl11

import (
	clw "github.com/rdwilliamson/clw11"
)

type (
	Context           clw.Context
	ContextProperties clw.ContextProperties
)

func CreateContext(p []ContextProperties, d []Device, callback func(err string, data []byte)) (Context, error) {

	devices := make([]clw.DeviceID, len(d))
	for i := range d {
		devices[i] = clw.DeviceID(d[i].ID)
	}
	var properties []clw.ContextProperties
	if p != nil {
		properties = make([]clw.ContextProperties, len(p))
		for i := range p {
			properties[i] = clw.ContextProperties(p[i])
		}
	}

	result, err := clw.CreateContext(properties, devices, callback)

	return Context(result), err
}

// TODO CreateContextFromType RetainContext ReleaseContext
// GetContextInfo only needs get reference count
