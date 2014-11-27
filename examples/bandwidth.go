package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	cl "github.com/rdwilliamson/cl11"
)

func check(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	platforms, err := cl.GetPlatforms()
	check(err)
	for _, p := range platforms {
		for _, d := range p.Devices {

			c, err := cl.CreateContext([]*cl.Device{d}, nil, nil, nil)
			check(err)

			cq, err := c.CreateCommandQueue(d, cl.QueueProfilingEnable)
			check(err)

			size := d.MaxMemAllocSize
			if size > 1024*1024*1024 {
				size = 1024 * 1024 * 1024
			}

			host, err := c.CreateHostBuffer(size, 0)
			check(err)

			random := make([]float32, size/4)
			for i := range random {
				random[i] = rand.Float32()
			}
			err = cq.EnqueueWriteBuffer(host, cl.NonBlocking, 0, random, nil, nil)
			check(err)

			device, err := c.CreateDeviceBuffer(size, 0)
			check(err)

			callbackChan := make(chan interface{})
			sendFunc := func(e *cl.Event, err error, userData interface{}) {

				if err != nil {
					callbackChan <- fmt.Sprintf("Error: %s", err.Error())
					return
				}

				check(e.GetProfilingInfo())

				name := userData.(string)
				transfered := float64(size) / 1024 / 1024
				duration := time.Duration(e.End - e.Start)
				transferSpeed := transfered / duration.Seconds() / 1024

				callbackChan <- fmt.Sprintf("Callback with Event Profiling: %s: %.2f MiB in %v (%.2f GiB/s)", name,
					transfered, duration, transferSpeed)
			}

			start := time.Now()

			var e cl.Event
			check(cq.EnqueueCopyBuffer(host, device, 0, 0, size, nil, &e))
			e.SetCallback(sendFunc, d.Name)

			check(cq.Finish())

			duration := time.Since(start)

			transfered := float64(size) / 1024 / 1024
			transferSpeed := transfered / duration.Seconds() / 1024

			fmt.Printf("Go Timer: %s: %.2f MiB in %v (%.2f GiB/s)\n", d.Name, transfered, duration, transferSpeed)
			fmt.Println(<-callbackChan)

			check(c.Release())
			check(cq.Release())
			check(host.Release())
			check(device.Release())
			check(e.Release())
		}
	}
}
