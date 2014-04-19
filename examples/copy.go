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
	for _, p := range platforms {
		for _, d := range p.Devices {

			c, err := cl.CreateContext([]*cl.Device{d}, nil, nil, nil)
			check(err)

			progam, err := c.CreateProgramWithSource([]byte(kernel))
			check(err)

			err = progam.Build([]*cl.Device{d}, "", nil, nil)
			check(err)

			kernel, err := progam.CreateKernel("copy")
			check(err)

			size := int64(256 * 1024 * 1024 / 4)
			if size*4 > d.MaxMemAllocSize/2 {
				size = d.MaxMemAllocSize / 2 / 4
			}
			values := utils.RandomFloat32(int(size))

			inData, err := c.CreateDeviceBufferInitializedBy(cl.MemReadOnly, values)
			check(err)
			outData, err := c.CreateDeviceBuffer(size*4, cl.MemWriteOnly)
			check(err)

			err = kernel.SetArguments(inData, outData, size)
			check(err)

			cq, err := c.CreateCommandQueue(d, cl.QueueProfilingEnable)
			check(err)

			localSize := kernel.WorkGroupInfo[0].PreferredWorkGroupSizeMultiple
			globalSize := int(size)
			if globalSize%localSize > 0 {
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
				duration := time.Duration(kernelEvent.End - kernelEvent.Start)
				fmt.Printf("%s copied %s in %v (%.2f GiB/s)\n", d.Name, snippets.PrintBytes(int64(size*4)),
					duration, float64(size*4)/1024/1024/1024/duration.Seconds())
			} else {
				fmt.Println(d.Name, "values do not match")
			}
		}
	}
}
