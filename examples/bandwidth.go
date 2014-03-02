package main

import (
	"fmt"
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
		devices, err := p.GetDevices()
		check(err)
		for _, d := range devices {

			c, err := cl.CreateContext([]*cl.Device{d}, cl.ContextProperties{}, nil)
			check(err)

			cq, err := c.CreateCommandQueue(d, cl.QueueProfilingEnable)
			check(err)

			size := d.MaxMemAllocSize

			host, err := c.CreateHostBuffer(size, 0)
			check(err)

			device, err := c.CreateDeviceBuffer(size, 0)
			check(err)

			callbackChan := make(chan interface{})
			sendFunc := func(e *cl.Event, userData interface{}) {

				check(e.GetProfilingInfo())

				name := userData.(string)
				transfered := float64(size) / 1024 / 1024
				duration := time.Duration(e.End - e.Start)
				transferSpeed := transfered / duration.Seconds() / 1024

				callbackChan <- fmt.Sprintf("Callback: %s: %.2f MiB in %v (%.2f GiB/s)", name, transfered, duration,
					transferSpeed)
			}
			recieveFunc := func() {
				fmt.Println(<-callbackChan)
				callbackChan <- ""
			}
			go recieveFunc()

			start := time.Now()

			var e cl.Event
			check(cq.CopyBuffer(host, device, 0, 0, size, nil, &e))
			e.SetCallback(sendFunc, d.Name)

			check(cq.Finish())

			duration := time.Since(start)

			transfered := float64(size) / 1024 / 1024
			transferSpeed := transfered / duration.Seconds() / 1024

			fmt.Printf("Go: %s: %.2f MiB in %v (%.2f GiB/s)\n", d.Name, transfered, duration, transferSpeed)

			check(host.Release())
			check(device.Release())

			<-callbackChan
		}
	}
}
