package cl11

import (
	"errors"
	"unsafe"

	clw "github.com/rdwilliamson/clw11"
)

// A 2D or 3D image (Depth and SlicePitch are 0 for 2D images).
type Image struct {
	id clw.Mem

	// The context the buffer was created on.
	Context *Context

	// The image format, the channel order and data type.
	Format ImageFormat

	// The width in pixels.
	Width int

	// The height in pixels.
	Height int

	// The depth in pixels. Zero for a 2D image.
	Depth int

	// Scan line width in bytes. Only valid if Host is not nil. If Host is not
	// nil then valid values are 0 for Width * size of element in bytes or
	// greater than or equal to Width * size of element in bytes.
	RowPitch int

	// The size in bytes of each 2D image in bytes. Only valid if Host is not
	// nil; valid values are 0 for RowPitch * Height or greater than or equal to
	// RowPitch * Height.
	SlicePitch int

	// Usage information for the buffer from the device's point of view.
	Flags MemFlags

	// The host backed memory for the buffer, if applicable.
	Host interface{}
}

type (
	MemObjectType clw.MemObjectType
	ChannelOrder  clw.ChannelOrder
	ChannelType   clw.ChannelType
)

type ImageFormat struct {
	ChannelOrder ChannelOrder
	ChannelType  ChannelType
}

const (
	R         = ChannelOrder(clw.R)
	A         = ChannelOrder(clw.A)
	RG        = ChannelOrder(clw.RG)
	RA        = ChannelOrder(clw.RA)
	RGB       = ChannelOrder(clw.RGB)
	RGBA      = ChannelOrder(clw.RGBA)
	BGRA      = ChannelOrder(clw.BGRA)
	ARGB      = ChannelOrder(clw.ARGB)
	Intensity = ChannelOrder(clw.Intensity)
	Luminance = ChannelOrder(clw.Luminance)
	Rx        = ChannelOrder(clw.Rx)
	RGx       = ChannelOrder(clw.RGx)
	RGBx      = ChannelOrder(clw.RGBx)
)

const (
	SnormInt8      = ChannelType(clw.SnormInt8)
	SnormInt16     = ChannelType(clw.SnormInt16)
	UnormInt8      = ChannelType(clw.UnormInt8)
	UnormInt16     = ChannelType(clw.UnormInt16)
	UnormShort565  = ChannelType(clw.UnormShort565)
	UnormShort555  = ChannelType(clw.UnormShort555)
	UnormInt101010 = ChannelType(clw.UnormInt101010)
	SignedInt8     = ChannelType(clw.SignedInt8)
	SignedInt16    = ChannelType(clw.SignedInt16)
	SignedInt32    = ChannelType(clw.SignedInt32)
	UnsignedInt8   = ChannelType(clw.UnsignedInt8)
	UnsignedInt16  = ChannelType(clw.UnsignedInt16)
	UnsignedInt32  = ChannelType(clw.UnsignedInt32)
	HalfFloat      = ChannelType(clw.HalfFloat)
	Float32        = ChannelType(clw.Float32)
)

var (
	UnsupportedImageFormatErr = errors.New("unsupported image format")
)

// Get the list of image formats supported by an OpenCL implementation.
func (c *Context) GetSupportedImage2DFormats(mf MemFlags) ([]ImageFormat, error) {

	var count clw.Uint
	err := clw.GetSupportedImageFormats(c.id, clw.MemFlags(mf), clw.MemObjectImage2D, 0, nil, &count)
	if err != nil {
		return nil, err
	}

	formats := make([]clw.ImageFormat, int(count))

	err = clw.GetSupportedImageFormats(c.id, clw.MemFlags(mf), clw.MemObjectImage2D, count, &formats[0], nil)
	if err != nil {
		return nil, err
	}

	results := make([]ImageFormat, int(count))

	for i := range formats {
		results[i].ChannelOrder = ChannelOrder(formats[i].ChannelOrder())
		results[i].ChannelType = ChannelType(formats[i].ChannelType())
	}

	return results, nil
}

// Creates a 2D image object.
//
// Creates an uninitialized buffer on the device.
func (c *Context) CreateDeviceImage2D(mf MemFlags, format ImageFormat, width, height int) (*Image, error) {

	cFormat := clw.CreateImageFormat(clw.ChannelOrder(format.ChannelOrder), clw.ChannelType(format.ChannelType))

	mem, err := clw.CreateImage2D(c.id, clw.MemFlags(mf), cFormat, clw.Size(width), clw.Size(height), 0, nil)
	if err != nil {
		return nil, err
	}

	return &Image{id: mem, Context: c, Format: format, Width: width, Height: height, Flags: mf}, nil
}

// func (c *Context) CreateDeviceImage2DInitializedBy(mf MemFlags, i image.Image) (*Image, error) {
// 	flags := clw.MemFlags(mf) | clw.MemCopyHostPointer
// 	return nil, nil
// }

// func (c *Context) CreateDeviceImage2DFromHostMem(mf MemFlags, i image.Image) (*Image, error) {
// 	flags := clw.MemFlags(mf) | clw.MemUseHostPointer
// 	return nil, nil
// }

// Creates a 2D image object.
//
// Creates an uninitialized buffer on the host.
func (c *Context) CreateHostImage2D(mf MemFlags, format ImageFormat, width, height int) (*Image, error) {

	flags := clw.MemFlags(mf) | clw.MemAllocHostPointer

	cFormat := clw.CreateImageFormat(clw.ChannelOrder(format.ChannelOrder), clw.ChannelType(format.ChannelType))

	mem, err := clw.CreateImage2D(c.id, flags, cFormat, clw.Size(width), clw.Size(height), 0, nil)
	if err != nil {
		return nil, err
	}

	return &Image{id: mem, Context: c, Format: format, Width: width, Height: height, Flags: mf}, nil
}

// func (c *Context) CreateHostImage2DInitializedBy(mf MemFlags, i image.Image) (*Image, error) {
// 	flags := clw.MemFlags(mf) | clw.MemAllocHostPointer | clw.MemCopyHostPointer
// 	return nil, nil
// }

// Increments the image object reference count.
//
// The OpenCL commands that return a buffer perform an implicit retain.
func (b *Image) Retain() error {
	return clw.RetainMemObject(b.id)
}

// Decrements the image object reference count.
//
// After the buffers reference count becomes zero and commands queued for
// execution that use the buffer have finished the buffer is deleted.
func (b *Image) Release() error {
	return clw.ReleaseMemObject(b.id)
}

// Return the image's reference count.
//
// The reference count returned should be considered immediately stale. It is
// unsuitable for general use in applications. This feature is provided for
// identifying memory leaks.
func (b *Image) ReferenceCount() (int, error) {

	var count clw.Uint
	err := clw.GetMemObjectInfo(b.id, clw.MemReferenceCount, clw.Size(unsafe.Sizeof(count)), unsafe.Pointer(&count),
		nil)
	if err != nil {
		return 0, err
	}

	return int(count), nil
}
