package main

import (
	"fmt"
	"time"

	"github.com/rdwilliamson/snippets"

	cl "github.com/rdwilliamson/cl11"
	"github.com/rdwilliamson/cl11/examples/utils"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

var kernel = `
__kernel void reduce(__global float* data, int size)
{
	if (get_global_id(0) < size) {
		data[get_global_id(0)] += 1;
	}
}
`

func main() {

	platforms, err := cl.GetPlatforms()
	check(err)
	for _, p := range platforms {
		devices, err := p.GetDevices()
		check(err)
		for _, d := range devices {

			c, err := cl.CreateContext([]*cl.Device{d}, cl.ContextProperties{}, nil)
			check(err)

			progam, err := c.CreateProgramWithSource([][]byte{[]byte(kernel)})
			check(err)

			err = progam.Build([]*cl.Device{d}, "")
			check(err)

			kernel, err := progam.CreateKernel("reduce")
			check(err)

			size := 256 * 1024 * 1024 / 4
			if size > int(d.MaxMemAllocSize) {
				size = int(d.MaxMemAllocSize)
			}
			values := utils.RandomFloat32(size)

			data, err := c.CreateDeviceBufferFromHost(cl.MemoryReadOnly, cl.ToBytes(values, nil))
			check(err)

			var resultMem [4]byte
			_, err = c.CreateDeviceBufferOnHost(cl.MemoryWriteOnly, resultMem[:])
			check(err)

			err = kernel.SetArgument(0, data)
			check(err)
			err = kernel.SetArgument(1, size)
			check(err)

			cq, err := c.CreateCommandQueue(d, cl.QueueProfilingEnable)
			check(err)

			var event cl.Event
			err = cq.EnqueueNDRangeKernel(kernel, nil, []int{size}, []int{int(d.MaxWorkGroupSize)}, nil, &event)
			check(err)

			check(event.Wait())

			check(event.GetProfilingInfo())

			fmt.Println(d.Name, time.Duration(event.End-event.Start), snippets.PrintBytes(int64(size*4)))
		}
	}
}
