package main

import (
	"fmt"

	cl "github.com/rdwilliamson/cl11"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

var kernel = `
__kernel void copy(__global float* in, __global float* out, long size)
{
	for (long id = get_global_id(0); id < size; id += get_global_size(0)) {
		out[id] = in[id];
	}
}
`

func main() {
	platforms, err := cl.GetPlatforms()
	check(err)
	for _, platform := range platforms {
		for _, device := range platform.Devices {

			context, err := cl.CreateContext([]*cl.Device{device}, nil, nil, nil)
			check(err)

			program, err := context.CreateProgramWithSource([]byte(kernel))
			check(err)

			programDevices, err := program.GetDevices()
			check(err)
			for _, v := range programDevices {
				fmt.Println(v.Name)
			}
		}

		context, err := cl.CreateContext(platform.Devices, nil, nil, nil)
		check(err)

		program, err := context.CreateProgramWithSource([]byte(kernel))
		check(err)

		programDevices, err := program.GetDevices()
		check(err)
		for _, v := range programDevices {
			fmt.Println(v.Name)
		}
	}
}
