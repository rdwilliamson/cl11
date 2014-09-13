// Converts an image to gray scale.
package main

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"runtime"
	"text/tabwriter"
	"time"

	cl "github.com/rdwilliamson/cl11"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

var kernel = `
__constant sampler_t sampler = CLK_NORMALIZED_COORDS_FALSE | CLK_FILTER_NEAREST;

__kernel void toGray(__read_only image2d_t input, __write_only image2d_t output)
{
	int width  = get_image_width(input);
	int height = get_image_height(input);

	for (int y = get_global_id(1); y < height; y += get_global_size(1)) {
		for (int x = get_global_id(0); x < width; x += get_global_size(0)) {

			uint4 pixel = read_imageui(input, sampler, (int2)(x, y));
#ifdef INTEGER
			uint v = (19595*pixel.x + 38469*pixel.y + 7472*pixel.z) >> 16;
#else
			uint v = 0.298912*pixel.x + 0.586611*pixel.y + 0.114478*pixel.z;
#endif
			write_imageui(output, (int2)(x, y), (uint4)(v, v, v, 255));
		}
	}
}
`

func openImage(file string) (image.Image, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func main() {

	if len(os.Args) < 2 {
		fmt.Println("no input file")
		os.Exit(1)
	}
	input, err := openImage(os.Args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	width, height := input.Bounds().Dx(), input.Bounds().Dy()
	output := image.NewRGBA(image.Rect(0, 0, width, height))

	var count int
	base := filepath.Base(os.Args[1])
	ext := filepath.Ext(os.Args[1])
	base = base[:len(base)-len(ext)]

	outText := tabwriter.NewWriter(os.Stdout, 0, 1, 1, ' ', 0)

	platforms, err := cl.GetPlatforms()
	check(err)
	for _, p := range platforms {
		for _, d := range p.Devices {
			count++

			contextCallback := func(err string, data []byte, userData interface{}) {
				fmt.Println(userData.(*cl.Device).Name, err)
			}
			c, err := cl.CreateContext([]*cl.Device{d}, nil, contextCallback, d)
			check(err)

			progam, err := c.CreateProgramWithSource([]byte(kernel))
			check(err)

			var options string
			if d.Type == cl.DeviceTypeCpu {
				options = "-D INTEGER"
			}
			err = progam.Build([]*cl.Device{d}, options, nil, nil)
			if err != nil {
				buildStatus, err := progam.BuildStatus(d)
				check(err)
				if buildStatus == cl.BuildError {
					fmt.Println(progam.BuildLog(d))
					os.Exit(1)
				}
			}
			check(err)

			kernel, err := progam.CreateKernel("toGray")
			check(err)

			cq, err := c.CreateCommandQueue(d, cl.QueueProfilingEnable)
			check(err)

			outData, err := c.CreateDeviceImage(cl.MemWriteOnly, cl.ImageFormat{cl.RGBA, cl.UnsignedInt8}, width,
				height, 1)
			check(err)

			var readEvent, kernelEvent, writeEvent cl.Event

			rgba := input.(*image.RGBA)
			format := cl.ImageFormat{cl.RGBA, cl.UnsignedInt8}
			inData, err := c.CreateDeviceImage(cl.MemReadOnly, format, rgba.Rect.Dx(), rgba.Rect.Dy(), 1)
			check(err)

			err = kernel.SetArguments(inData, outData)
			check(err)

			localWidth := kernel.WorkGroupInfo[0].PreferredWorkGroupSizeMultiple
			if runtime.GOOS == "darwin" && d.Type == cl.DeviceTypeCpu {
				localWidth = 128
			}
			localHeight := d.MaxWorkGroupSize / localWidth
			if localHeight > d.MaxWorkItemSizes[1] {
				localHeight = d.MaxWorkItemSizes[1]
			}

			globalWidth, remainder := width/localWidth, width%localWidth
			if remainder > 0 {
				globalWidth++
			}
			globalWidth *= localWidth
			globalHeight, remainder := height/localHeight, height%localHeight
			if remainder > 0 {
				globalHeight++
			}
			globalHeight *= localHeight

			for i := 0; i < 2; i++ {

				start := time.Now()

				err = cq.EnqueueWriteImageFromImage(inData, cl.NonBlocking, input, nil, &writeEvent)
				check(err)

				err = cq.EnqueueNDRangeKernel(kernel, nil, []int{globalWidth, globalHeight}, []int{localWidth, localHeight},
					[]*cl.Event{&writeEvent}, &kernelEvent)
				check(err)

				err = cq.EnqueueReadImageToImage(outData, cl.NonBlocking, output, []*cl.Event{&kernelEvent}, &readEvent)
				check(err)

				err = readEvent.Wait()
				check(err)

				if i == 1 {

					wall := time.Since(start)

					err = writeEvent.GetProfilingInfo()
					check(err)
					err = kernelEvent.GetProfilingInfo()
					check(err)
					err = readEvent.GetProfilingInfo()
					check(err)

					size := float64(rgba.Rect.Dx()*rgba.Rect.Dy()*4) / 1024 / 1024 / 1024

					fmt.Fprintf(outText,
						"%s\ton\t%s\tWrite (%v, %.2fGB/s)\tIdle (%v)\tKernel (%v)\tIdle (%v)\tRead (%v, %.2fGB/s)\t"+
							"CPU Wall (%v)\n",
						d.Name, p.Name,
						time.Duration(writeEvent.End-writeEvent.Start),
						size/time.Duration(writeEvent.End-writeEvent.Start).Seconds(),
						time.Duration(kernelEvent.Start-writeEvent.End),
						time.Duration(kernelEvent.End-kernelEvent.Start),
						time.Duration(readEvent.Start-kernelEvent.End),
						time.Duration(readEvent.End-readEvent.Start),
						size/time.Duration(readEvent.End-readEvent.Start).Seconds(),
						wall)
				}
			}

			outFile, err := os.Create(fmt.Sprintf("%s%d%s", base, count, ext))
			check(err)
			defer outFile.Close()

			err = png.Encode(outFile, output)
			check(err)
		}
	}

	err = outText.Flush()
	check(err)
}
