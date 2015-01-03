// Converts an image to gray scale.
package main

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"os"
	"path/filepath"

	cl "github.com/rdwilliamson/cl11"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

var kernelSrc = []byte(`
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
			uint v = 0.298912f*pixel.x + 0.586611f*pixel.y + 0.114478f*pixel.z;
#endif
			write_imageui(output, (int2)(x, y), (uint4)(v, v, v, 255));
		}
	}
}
`)

func readImage(file string) (image.Image, error) {
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

func writeImage(file string, img image.Image) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	err = png.Encode(f, img)
	if err != nil {
		f.Close()
		return err
	}
	return f.Close()
}

func main() {

	// Read the input file.
	if len(os.Args) < 2 {
		fmt.Println("no input file")
		os.Exit(1)
	}
	input, err := readImage(os.Args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	width, height := input.Bounds().Dx(), input.Bounds().Dy()
	base := filepath.Base(os.Args[1])
	ext := filepath.Ext(os.Args[1])
	base = base[:len(base)-len(ext)]

	// Run the kernel on every device of every platform.
	var count int
	allPlatforms, err := cl.GetPlatforms()
	check(err)
	for _, platform := range allPlatforms {
		for _, device := range platform.Devices {
			count++

			// Create the OpenCL context.
			c, err := cl.CreateContext([]*cl.Device{device}, nil, nil, nil)
			check(err)

			// Create the device buffers.
			format := cl.ImageFormat{cl.RGBA, cl.UnsignedInt8}
			outData, err := c.CreateDeviceImage(cl.MemWriteOnly, format, width, height, 1)
			check(err)
			inData, err := c.CreateDeviceImage(cl.MemReadOnly, format, width, height, 1)
			check(err)

			// Create the kernel and set its arguments.
			progam, err := c.CreateProgramWithSource(kernelSrc)
			check(err)
			var options string
			if device.Type == cl.DeviceTypeCpu {
				options = "-D INTEGER"
			}
			err = progam.Build([]*cl.Device{device}, options, nil, nil)
			check(err)
			kernel, err := progam.CreateKernel("toGray")
			check(err)
			err = kernel.SetArguments(inData, outData)
			check(err)

			// Create the command queue then use it to copy the image to the
			// device and run the kernel.
			cq, err := c.CreateCommandQueue(device, 0)
			check(err)
			inMapped, err := cq.EnqueueMapImage(inData, cl.Blocking, cl.MapWrite, cl.EntireImage, nil, nil)
			check(err)
			inRgba := inMapped.RGBA()
			draw.Draw(inRgba, inRgba.Bounds(), input, input.Bounds().Min, draw.Src)
			err = cq.EnqueueUnmapImage(inMapped, nil, nil)
			check(err)
			err = cq.EnqueueNDRangeKernel(kernel, nil, []int{width, height}, []int{128, 1}, nil, nil)
			check(err)

			// Copy the result to a file.
			outMapped, err := cq.EnqueueMapImage(outData, cl.Blocking, cl.MapRead, cl.EntireImage, nil, nil)
			check(err)
			err = writeImage(fmt.Sprintf("%s%d%s", base, count, ext), outMapped.RGBA())
			check(err)
			err = cq.EnqueueUnmapImage(outMapped, nil, nil)
			check(err)
		}
	}
}
