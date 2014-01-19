package cl11

import (
	"fmt"

	clw "github.com/rdwilliamson/clw11"
)

type Context struct {
	ID         clw.Context
	Devices    []*Device
	Properties ContextProperties
}

func (c Context) String() string {
	return fmt.Sprintf("%x", c.ID)
}

type ContextProperties struct {
	// TODO add options.
}

func CreateContext(d []*Device, cp ContextProperties, callback func(err string, data []byte)) (*Context, error) {

	devices := make([]clw.DeviceID, len(d))
	for i := range d {
		devices[i] = clw.DeviceID(d[i].ID)
	}
	var properties []clw.ContextProperties
	// TODO convert struct into C array of properties.

	context, err := clw.CreateContext(properties, devices, callback)
	if err != nil {
		return nil, err
	}
	return &Context{context, d, cp}, nil
}
