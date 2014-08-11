// Converts an image to gray scale.
package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"time"

	cl "github.com/rdwilliamson/cl11"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

var kernel = `
#define int32_t int

const sampler_t sampler = CLK_NORMALIZED_COORDS_FALSE | CLK_FILTER_NEAREST;

// __kernel void to_gray(__read_only image2d_t input, __write_only image2d_t output)
__kernel void to_gray(__read_only image2d_t input, __global float* output)
{
	int width = get_image_width(input);
	int height = get_image_height(input);

	for (int32_t y = get_global_id(1); y < height; y += get_global_size(1)) {
		for (int32_t x = get_global_id(0); x < width; x += get_global_size(0)) {

			uint4 pixel = read_imageui(input, sampler, (int2)(x, y));
			// write_imageui(output, (int2)(x, y), 0.298912*pixel.x + 0.586611*pixel.y + 0.114478*pixel.z);
			// write_imageui(output, (int2)(x, y), (uint4)(pixel.x, pixel.y, pixel.z, 255));
			output[y*width+x] = 0.298912*pixel.x + 0.586611*pixel.y + 0.114478*pixel.z;
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
	// output := image.NewGray(image.Rect(0, 0, width, height))
	output := image.NewRGBA(image.Rect(0, 0, width, height))

	var count int
	base := filepath.Base(os.Args[1])
	ext := filepath.Ext(os.Args[1])
	base = base[:len(base)-len(ext)]

	platforms, err := cl.GetPlatforms()
	check(err)
	for _, p := range platforms {
		for _, d := range p.Devices {
			count++

			c, err := cl.CreateContext([]*cl.Device{d}, nil, nil, nil)
			check(err)

			progam, err := c.CreateProgramWithSource([]byte(kernel))
			check(err)

			err = progam.Build([]*cl.Device{d}, "", nil, nil)
			if err != nil {
				buildStatus, err2 := progam.BuildStatus(d)
				check(err2)
				if buildStatus == cl.BuildError {
					fmt.Println(progam.BuildLog(d))
					os.Exit(1)
				}
			}
			check(err)

			// TODO modify kernel to copy from image to image
			kernel, err := progam.CreateKernel("to_gray")
			check(err)

			inData, err := c.CreateDeviceImage2DInitializedByImage(cl.MemReadOnly, input)
			check(err)

			// outData, err := c.CreateDeviceImage2D(cl.MemWriteOnly, cl.ImageFormat{cl.RGBA, cl.UnormInt8}, width, height)
			// check(err)

			outData, err := c.CreateDeviceBuffer(int64(width*height*4), cl.MemWriteOnly)
			check(err)

			err = kernel.SetArguments(inData, outData)
			check(err)

			cq, err := c.CreateCommandQueue(d, cl.QueueProfilingEnable)
			check(err)

			localWidth := kernel.WorkGroupInfo[0].PreferredWorkGroupSizeMultiple
			localHeight := d.MaxWorkGroupSize / localWidth
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

			var kernelEvent cl.Event
			err = cq.EnqueueNDRangeKernel(kernel, nil, []int{globalWidth, globalHeight}, []int{localWidth, localHeight},
				nil, &kernelEvent)
			check(err)

			// err = cq.EnqueueWriteImageToImage(outData, cl.Blocking, output, []*cl.Event{&kernelEvent}, nil)
			// check(err)

			mb, err := cq.MapBuffer(outData, cl.Blocking, cl.MapRead, 0, int64(width*height*4), nil, nil)
			check(err)

			values := mb.Float32Slice()
			for y := 0; y < height; y++ {
				for x := 0; x < width; x++ {
					v := uint8(values[y*width+x])
					output.Set(x, y, color.RGBA{R: v, G: v, B: v, A: 255})
				}
			}

			var event cl.Event
			check(cq.UnmapBuffer(mb, nil, &event))
			check(event.Wait())

			check(cq.Finish())

			check(kernelEvent.GetProfilingInfo())
			fmt.Println(d.Name, time.Duration(kernelEvent.End-kernelEvent.Start))

			outFile, err := os.Create(fmt.Sprintf("%s%d%s", base, count, ext))
			check(err)
			defer outFile.Close()

			err = png.Encode(outFile, output)
			check(err)
		}
	}
}
