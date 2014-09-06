package cl11

import (
	"errors"
	"image"
	"unsafe"

	clw "github.com/rdwilliamson/clw11"
)

// A 2D or 3D image (Depth is 1 and SlicePitch is 0 for 2D images).
//
// TODO If Depth is 1 and SlicePitch is 0 it is ambiguous whether it is a 2D or
// 3D image.
type Image struct {
	id clw.Mem

	// The context the buffer was created on.
	Context *Context

	// The image format, the channel order and data type.
	Format ImageFormat

	// Size of each element of the image memory object
	ElementSize int

	// The width in pixels.
	Width int

	// The height in pixels.
	Height int

	// The depth in pixels. One for a 2D image.
	Depth int

	// Scan line width in bytes. Only valid if Host is not nil. If Host is not
	// nil then valid values are 0 (which is the same as Width * ElementSize) or
	// a value greater than or equal to Width * ElementSize.
	RowPitch int

	// The size in bytes of each 2D image in bytes. Only valid if Host is not
	// nil; valid values are 0 (which is the same as RowPitch * Height) or a
	// value greater than or equal to RowPitch * Height.
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

const (
	MemObjectBuffer  = MemObjectType(clw.MemObjectBuffer)
	MemObjectImage2D = MemObjectType(clw.MemObjectImage2D)
	MemObjectImage3D = MemObjectType(clw.MemObjectImage3D)
)

var (
	ErrUnsupportedImageFormat        = errors.New("cl: unsupported image format")
	ErrImageRectToBufferSizeMismatch = errors.New("cl: image rectangle to buffer size mismatch")
	ErrInvalidImageFormat            = errors.New("cl: invalid image format")
)

// Invalid image formats will return 0.
func (i *ImageFormat) elementSize() int {

	var channels int
	switch i.ChannelOrder {
	case R, A, Intensity, Luminance:
		channels = 1
	case RG, RA:
		channels = 2
	case RGB:
		channels = 3
	case RGBA, BGRA, ARGB:
		channels = 4
	case Rx, RGx, RGBx:
		// TODO how many channels do each of these have?
		panic("unknown number of channels in format")
	}

	var channelBytes int
	switch i.ChannelType {
	case SnormInt8, UnsignedInt8, UnormInt8, UnormInt16, SignedInt8:
		channelBytes = 1
	case SnormInt16, UnormShort565, UnormShort555, SignedInt16, UnsignedInt16, HalfFloat:
		channelBytes = 2
	case UnormInt101010, SignedInt32, UnsignedInt32, Float32:
		channelBytes = 4
	}

	return channels * channelBytes
}

// Get the list of image formats supported by an OpenCL implementation.
func (c *Context) GetSupportedImageFormats(mf MemFlags, mot MemObjectType) ([]ImageFormat, error) {

	var count clw.Uint
	err := clw.GetSupportedImageFormats(c.id, clw.MemFlags(mf), clw.MemObjectType(mot), 0, nil, &count)
	if err != nil {
		return nil, err
	}

	formats := make([]clw.ImageFormat, int(count))

	err = clw.GetSupportedImageFormats(c.id, clw.MemFlags(mf), clw.MemObjectType(mot), count, &formats[0], nil)
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

// Creates an image object.
//
// Creates an uninitialized buffer on the device.
func (c *Context) CreateDeviceImage(mf MemFlags, format ImageFormat, width, height, depth int) (*Image, error) {

	cFormat := clw.CreateImageFormat(clw.ChannelOrder(format.ChannelOrder), clw.ChannelType(format.ChannelType))

	var mem clw.Mem
	var err error
	if depth == 1 {
		mem, err = clw.CreateImage2D(c.id, clw.MemFlags(mf), cFormat, clw.Size(width), clw.Size(height), 0, nil)
	} else {
		mem, err = clw.CreateImage3D(c.id, clw.MemFlags(mf), cFormat, clw.Size(width), clw.Size(height),
			clw.Size(depth), 0, 0, nil)
	}
	if err != nil {
		return nil, err
	}

	return &Image{
			id:          mem,
			Context:     c,
			Format:      format,
			ElementSize: format.elementSize(),
			Width:       width,
			Height:      height,
			Depth:       depth,
			Flags:       mf,
		},
		nil
}

// Only source and region are used from the rectangle (though the destination is
// still validated).
func (c *Context) CreateDeviceImageInitializedBy(mf MemFlags, format ImageFormat, r *Rect,
	value interface{}) (*Image, error) {

	// Validate input.
	if !r.valid() {
		return nil, ErrInvalidRect
	}
	dim := r.Src.dimensions()
	if dim == 2 && r.Region[2] != 1 {
		return nil, ErrInvalidRect
	}

	var scratch [scratchSize]byte
	pointer, size := getPointerAndSize(value, unsafe.Pointer(&scratch[0]))
	if uintptr(r.srcBytes()) > size {
		return nil, ErrImageRectToBufferSizeMismatch
	}

	if format.elementSize() == 0 {
		return nil, ErrInvalidImageFormat
	}

	// Create the image.
	cFormat := clw.CreateImageFormat(clw.ChannelOrder(format.ChannelOrder), clw.ChannelType(format.ChannelType))
	flags := clw.MemFlags(mf) | clw.MemCopyHostPointer

	var mem clw.Mem
	var err error
	if dim == 2 {
		mem, err = clw.CreateImage2D(c.id, flags, cFormat, r.width(), r.height(), r.Src.rowPitch(), pointer)
	} else {
		mem, err = clw.CreateImage3D(c.id, flags, cFormat, r.width(), r.height(), r.depth(),
			r.Src.rowPitch(), r.Src.slicePitch(), pointer)
	}
	if err != nil {
		return nil, err
	}

	return &Image{id: mem,
			Context:     c,
			Format:      format,
			ElementSize: format.elementSize(),
			Width:       int(r.width()),
			Height:      int(r.height()),
			Depth:       int(r.depth()),
			RowPitch:    int(r.Src.rowPitch()),
			SlicePitch:  int(r.Src.slicePitch()),
			Flags:       mf,
		},
		nil
}

// Creates a 2D image object.
//
// Creates an initialized buffer on the host. Currently only *image.RGBA format
// is supported.
func (c *Context) CreateDeviceImageInitializedByImage(mf MemFlags, i image.Image) (*Image, error) {

	var pointer unsafe.Pointer
	var width, height, imageRowPitch clw.Size
	var format ImageFormat

	switch v := i.(type) {

	case *image.RGBA:
		pointer = unsafe.Pointer(&v.Pix[v.Rect.Min.Y*v.Stride+v.Rect.Min.X*4])
		width = clw.Size(v.Rect.Dx())
		height = clw.Size(v.Rect.Dy())
		imageRowPitch = clw.Size(v.Stride)
		format.ChannelOrder = RGBA
		format.ChannelType = UnsignedInt8

	default:
		return nil, ErrUnsupportedImageFormat
	}

	cFormat := clw.CreateImageFormat(clw.ChannelOrder(format.ChannelOrder), clw.ChannelType(format.ChannelType))
	flags := clw.MemFlags(mf) | clw.MemCopyHostPointer

	mem, err := clw.CreateImage2D(c.id, flags, cFormat, width, height, imageRowPitch, pointer)
	if err != nil {
		return nil, err
	}

	return &Image{id: mem,
			Context:     c,
			Format:      format,
			ElementSize: format.elementSize(),
			Width:       int(width),
			Height:      int(height),
			Depth:       1,
			Flags:       mf,
		},
		nil
}

func (c *Context) CreateDeviceImageFromHostMem(mf MemFlags, format ImageFormat, r *Rect,
	value interface{}) (*Image, error) {
	// flags := clw.MemFlags(mf) | clw.MemUseHostPointer
	return nil, nil
}

func (c *Context) CreateDeviceImageFromHostImage(mf MemFlags, i image.Image) (*Image, error) {
	// flags := clw.MemFlags(mf) | clw.MemUseHostPointer
	return nil, nil
}

// Creates a image object.
//
// Creates an uninitialized buffer on the host.
func (c *Context) CreateHostImage(mf MemFlags, format ImageFormat, width, height, depth int) (*Image, error) {

	cFormat := clw.CreateImageFormat(clw.ChannelOrder(format.ChannelOrder), clw.ChannelType(format.ChannelType))
	flags := clw.MemFlags(mf) | clw.MemAllocHostPointer

	var mem clw.Mem
	var err error
	if depth == 1 {
		mem, err = clw.CreateImage2D(c.id, flags, cFormat, clw.Size(width), clw.Size(height), 0, nil)
	} else {
		mem, err = clw.CreateImage3D(c.id, flags, cFormat, clw.Size(width), clw.Size(height), clw.Size(depth), 0, 0,
			nil)
	}
	if err != nil {
		return nil, err
	}

	return &Image{id: mem,
			Context:     c,
			Format:      format,
			ElementSize: format.elementSize(),
			Width:       width,
			Height:      height,
			Depth:       depth,
			Flags:       mf,
		},
		nil
}

func (c *Context) CreateHostImageInitializedBy(mf MemFlags, format ImageFormat, r *Rect,
	value interface{}) (*Image, error) {
	// flags := clw.MemFlags(mf) | clw.MemAllocHostPointer | clw.MemCopyHostPointer
	return nil, nil
}

func (c *Context) CreateHostImageInitializedByImage(mf MemFlags, i image.Image) (*Image, error) {
	// flags := clw.MemFlags(mf) | clw.MemAllocHostPointer | clw.MemCopyHostPointer
	return nil, nil
}

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

func (cq *CommandQueue) EnqueueReadImage(dst *Image, bc BlockingCall, r *Rect, src interface{}, waitList []*Event,
	e *Event) error {

	var event *clw.Event
	if e != nil {
		event = &e.id
		e.Context = cq.Context
		e.CommandType = CommandReadImage
		e.CommandQueue = cq
	}

	pointer, _, err := tryPointerAndSize(src)
	if err != nil {
		return wrapError(err)
	}

	return clw.EnqueueWriteImage(cq.id, dst.id, clw.Bool(bc), r.Src.origin(), r.region(), r.Src.rowPitch(),
		r.Src.slicePitch(), pointer, cq.toEvents(waitList), event)
}

func (cq *CommandQueue) EnqueueReadImageFromImage(dst *Image, bc BlockingCall, src image.Image, waitList []*Event,
	e *Event) error {

	var rect Rect
	var actualSrc interface{}

	switch v := src.(type) {

	case *image.RGBA:
		rect.Region[0] = int64(v.Rect.Dx())
		rect.Region[1] = int64(v.Rect.Dy())
		rect.Region[2] = 1
		rect.Src.RowPitch = int64(v.Stride)
		actualSrc = v.Pix[v.Rect.Min.Y*v.Stride+v.Rect.Min.X*4 : (v.Rect.Max.Y-1)*v.Stride+v.Rect.Max.X-1]

	default:
		return ErrUnsupportedImageFormat
	}

	return cq.EnqueueReadImage(dst, bc, &rect, actualSrc, waitList, e)
}

// rect only uses the dst and region, src is ignored
func (cq *CommandQueue) EnqueueWriteImage(src *Image, bc BlockingCall, r *Rect, dst interface{}, waitList []*Event,
	e *Event) error {

	var event *clw.Event
	if e != nil {
		event = &e.id
		e.Context = cq.Context
		e.CommandType = CommandWriteImage
		e.CommandQueue = cq
	}

	pointer, _, err := tryPointerAndSize(dst)
	if err != nil {
		return wrapError(err)
	}

	return clw.EnqueueReadImage(cq.id, src.id, clw.Bool(bc), r.Dst.origin(), r.region(), r.Dst.rowPitch(),
		r.Dst.slicePitch(), pointer, cq.toEvents(waitList), event)
}

func (cq *CommandQueue) EnqueueWriteImageToImage(src *Image, bc BlockingCall, dst image.Image, waitList []*Event,
	e *Event) error {

	var rect Rect
	var actualDst interface{}

	switch v := dst.(type) {

	case *image.RGBA:
		rect.Region[0] = int64(v.Rect.Dx())
		rect.Region[1] = int64(v.Rect.Dy())
		rect.Region[2] = 1
		rect.Src.RowPitch = int64(v.Stride)
		actualDst = v.Pix[v.Rect.Min.Y*v.Stride+v.Rect.Min.X*4 : (v.Rect.Max.Y-1)*v.Stride+v.Rect.Max.X-1]

	default:
		return ErrUnsupportedImageFormat
	}

	return cq.EnqueueWriteImage(src, bc, &rect, actualDst, waitList, e)
}

func (cq *CommandQueue) EnqueueCopyImage(src, dst *Image, r *Rect, waitList []*Event, e *Event) error {

	var event *clw.Event
	if e != nil {
		event = &e.id
		e.Context = cq.Context
		e.CommandType = CommandCopyImage
		e.CommandQueue = cq
	}

	return clw.EnqueueCopyImage(cq.id, src.id, dst.id, r.Src.origin(), r.Dst.origin(), r.region(),
		cq.toEvents(waitList), event)
}
