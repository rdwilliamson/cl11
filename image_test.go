package cl11

import (
	"image"
	"image/color"
	"math/rand"
	"reflect"

	"testing"
)

func TestImage(t *testing.T) {
	allDevices := getDevices(t)
	for _, device := range allDevices {

		var toRelease []Object
		img0 := image.NewRGBA(image.Rect(0, 0, 256, 256))
		for y, height := 0, img0.Bounds().Dy(); y < height; y++ {
			for x, width := 0, img0.Bounds().Dx(); x < width; x++ {
				r := uint8(rand.Intn(256))
				g := uint8(rand.Intn(256))
				b := uint8(rand.Intn(256))
				a := uint8(rand.Intn(256))
				img0.Set(x, y, color.NRGBA{r, g, b, a})
			}
		}
		img1 := image.NewRGBA(img0.Bounds())

		ctx, err := CreateContext([]*Device{device}, nil, nil, nil)
		if err != nil {
			t.Error(err)
			continue
		}
		toRelease = append(toRelease, ctx)

		cq, err := ctx.CreateCommandQueue(device, 0)
		if err != nil {
			t.Error(err)
			releaseAll(toRelease, t)
			continue
		}
		toRelease = append(toRelease, cq)

		allFormats, err := ctx.GetSupportedImageFormats(MemReadOnly, MemObjectImage2D)
		if err != nil {
			t.Error(err)
			releaseAll(toRelease, t)
			continue
		}

		var format ImageFormat
		for _, v := range allFormats {
			if v.ChannelOrder == RGBA && v.ChannelType == UnsignedInt8 {
				format = v
				break
			}
		}
		if format.ChannelOrder == 0 || format.ChannelType == 0 {
			t.Error("could not find desired image format")
			releaseAll(toRelease, t)
			continue
		}

		host0, err := ctx.CreateHostImage(MemReadOnly, format, img0.Bounds().Dx(), img0.Bounds().Dy(), 1)
		if err != nil {
			t.Error(err)
			releaseAll(toRelease, t)
			continue
		}
		toRelease = append(toRelease, host0)

		host1, err := ctx.CreateHostImage(MemWriteOnly, format, img0.Bounds().Dx(), img0.Bounds().Dy(), 1)
		if err != nil {
			t.Error(err)
			releaseAll(toRelease, t)
			continue
		}
		toRelease = append(toRelease, host1)

		device0, err := ctx.CreateDeviceImage(MemReadWrite, format, img0.Bounds().Dx(), img0.Bounds().Dy(), 1)
		if err != nil {
			t.Error(err)
			releaseAll(toRelease, t)
			continue
		}
		toRelease = append(toRelease, device0)

		err = cq.EnqueueWriteImageFromImage(host0, NonBlocking, img0, nil, nil)
		if err != nil {
			t.Error(err)
			releaseAll(toRelease, t)
			continue
		}

		var rect Rect
		rect.Src.RowPitch = int64(img0.Stride)
		rect.Dst.RowPitch = int64(img0.Stride)
		rect.Region[0] = int64(img0.Bounds().Dx())
		rect.Region[1] = int64(img0.Bounds().Dy())
		rect.Region[2] = 1
		err = cq.EnqueueCopyImage(host0, device0, &rect, nil, nil)
		if err != nil {
			t.Error(err)
			releaseAll(toRelease, t)
			continue
		}

		err = cq.EnqueueReadImageToImage(device0, NonBlocking, img1, nil, nil)
		if err != nil {
			t.Error(err)
			releaseAll(toRelease, t)
			continue
		}

		err = cq.Finish()
		if err != nil {
			t.Error(err)
			releaseAll(toRelease, t)
			continue
		}

		if !reflect.DeepEqual(img0, img1) {
			t.Error("images don't match")
			releaseAll(toRelease, t)
			continue
		}
	}
}
