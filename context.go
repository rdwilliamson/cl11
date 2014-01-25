package cl11

import clw "github.com/rdwilliamson/clw11"

type Context struct {
	id         clw.Context
	Devices    []*Device
	Properties ContextProperties
}

type ContextProperties struct {
	// TODO add options.
}

func CreateContext(d []*Device, cp ContextProperties, callback func(err string, data []byte)) (*Context, error) {

	devices := make([]clw.DeviceID, len(d))
	for i := range d {
		devices[i] = clw.DeviceID(d[i].id)
	}
	var properties []clw.ContextProperties
	// TODO convert struct into C array of properties.

	context, err := clw.CreateContext(properties, devices, callback)
	if err != nil {
		return nil, err
	}
	return &Context{context, d, cp}, nil
}
