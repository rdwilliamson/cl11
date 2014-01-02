package cl11

import (
	clw "github.com/rdwilliamson/clw11"
)

type Context clw.Context

func CreateContext(d []Device) (Context, error) {
	devices := make([]clw.DeviceID, len(d))
	for i := range d {
		devices[i] = clw.DeviceID(d[i].ID)
	}
	result, err := clw.CreateContext(nil, devices)
	return Context(result), err
}
