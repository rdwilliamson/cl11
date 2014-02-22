package main

import (
	"fmt"
	"reflect"
	"time"

	cl "github.com/rdwilliamson/cl11"
	"github.com/rdwilliamson/cl11/examples/utils"
	"github.com/rdwilliamson/snippets"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

var kernel = `
__kernel void copy(__global float* in, __global float* out, int size)
{
	for (int id = get_global_id(0); id < size; id += get_global_size(0)) {
		out[id] = in[id];
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

			progam, err := c.CreateProgramWithSource([]byte(kernel))
			check(err)

			err = progam.Build([]*cl.Device{d}, "")
			check(err)

			kernel, err := progam.CreateKernel("copy")
			check(err)

			size := 256 * 1024 * 1024 / 4
			if size*4 > int(d.MaxMemAllocSize)/2 {
				size = int(d.MaxMemAllocSize) / 2 / 4
			}
			values := utils.RandomFloat32(size)

			inData, err := c.CreateDeviceBufferInitializedBy(cl.MemoryReadOnly, values)
			check(err)
			outData, err := c.CreateDeviceBuffer(size*4, cl.MemoryWriteOnly)
			check(err)

			err = kernel.SetArguments(inData, outData, size)
			check(err)

			cq, err := c.CreateCommandQueue(d, cl.QueueProfilingEnable)
			check(err)

			localSize := kernel.WorkGroupInfo[0].PreferredWorkGroupSizeMultiple
			globalSize := size
			if size%localSize > 0 {
				globalSize = (globalSize/localSize + 1) * localSize
			}

			var kernelEvent cl.Event
			err = cq.EnqueueNDRangeKernel(kernel, nil, []int{globalSize}, []int{localSize}, nil, &kernelEvent)
			check(err)

			check(kernelEvent.Wait())
			check(kernelEvent.GetProfilingInfo())

			mb, err := cq.MapBuffer(outData, cl.Blocking, cl.MapRead, 0, size*4, nil, nil)
			check(err)

			equal := reflect.DeepEqual(values, mb.Float32Slice())

			var event cl.Event
			check(cq.UnmapBuffer(mb, nil, &event))
			check(event.Wait())

			if equal {
				fmt.Println(d.Name, time.Duration(kernelEvent.End-kernelEvent.Start), snippets.PrintBytes(int64(size*4)))
			} else {
				fmt.Println(d.Name, "values do not match")
			}
		}
	}
}
