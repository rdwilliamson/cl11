package cl11

import (
	"image"
	"reflect"
	"unsafe"
)

// A mapped image. Has convenience functions for common data types.
type MappedImage struct {
	Image      *Image
	pointer    unsafe.Pointer
	rowPitch   int64 // Scan line width in bytes.
	slicePitch int64 // The size in bytes of each 2D image (0 for 2D image).
}

func (mi *MappedImage) GoImage() (image.Image, error) {
	if mi.Image.Format.ChannelOrder == RGBA && mi.Image.Format.ChannelType == UnsignedInt8 {
		return mi.rgbaImage(), nil
	}
	return nil, UnsupportedImageFormat
}

func (mi *MappedImage) rgbaImage() *image.RGBA {

	var header reflect.SliceHeader
	header.Data = uintptr(mi.pointer)
	size := int(mi.rowPitch) * mi.Image.Height
	header.Len = size
	header.Cap = size

	var result image.RGBA
	result.Pix = *(*[]uint8)(unsafe.Pointer(&header))
	result.Stride = int(mi.rowPitch)
	result.Rect = image.Rectangle{image.Point{}, image.Point{mi.Image.Width, mi.Image.Height}}
	return &result
}
