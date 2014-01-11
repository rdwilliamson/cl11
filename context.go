package cl11

import (
	"fmt"

	clw "github.com/rdwilliamson/clw11"
)

type Context struct {
	ID         ContextID
	Devices    []*Device
	Properties []ContextProperties
}

func (c Context) String() string {
	return fmt.Sprintf("%x", c.ID)
}

type (
	ContextID         clw.Context
	ContextProperties clw.ContextProperties
)

func CreateContext(p []ContextProperties, d []*Device, callback func(err string, data []byte)) (*Context, error) {

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
	if err != nil {
		return nil, err
	}

	return &Context{ContextID(result), d, p}, nil
}

// TODO CreateContextFromType RetainContext ReleaseContext
// GetContextInfo only needs get reference count
