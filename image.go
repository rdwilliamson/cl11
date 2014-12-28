package cl11

import (
	"fmt"
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
	// nil. valid values are 0 (which is the same as RowPitch * Height) or a
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

func (co ChannelOrder) String() string {
	switch co {
	case R:
		return "R"
	case A:
		return "A"
	case RG:
		return "RG"
	case RA:
		return "RA"
	case RGB:
		return "RGB"
	case RGBA:
		return "RGBA"
	case BGRA:
		return "BGRA"
	case ARGB:
		return "ARGB"
	case Intensity:
		return "Intensity"
	case Luminance:
		return "Luminance"
	case Rx:
		return "Rx"
	case RGx:
		return "RGx"
	case RGBx:
		return "RGBx"
	default:
		return fmt.Sprintf("unknown (%x)", int(co))
	}
}

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

func (ct ChannelType) String() string {
	switch ct {
	case SnormInt8:
		return "SnormInt8"
	case SnormInt16:
		return "SnormInt16"
	case UnormInt8:
		return "UnormInt8"
	case UnormInt16:
		return "UnormInt16"
	case UnormShort565:
		return "UnormShort565"
	case UnormShort555:
		return "UnormShort555"
	case UnormInt101010:
		return "UnormInt101010"
	case SignedInt8:
		return "SignedInt8"
	case SignedInt16:
		return "SignedInt16"
	case SignedInt32:
		return "SignedInt32"
	case UnsignedInt8:
		return "UnsignedInt8"
	case UnsignedInt16:
		return "UnsignedInt16"
	case UnsignedInt32:
		return "UnsignedInt32"
	case HalfFloat:
		return "HalfFloat"
	case Float32:
		return "Float32"
	default:
		return fmt.Sprintf("unknown (%x)", int(ct))
	}
}

const (
	MemObjectBuffer  = MemObjectType(clw.MemObjectBuffer)
	MemObjectImage2D = MemObjectType(clw.MemObjectImage2D)
	MemObjectImage3D = MemObjectType(clw.MemObjectImage3D)
)

// Invalid image formats will return 0.
func (i *ImageFormat) elementSize() int {

	var channels int
	switch i.ChannelOrder {
	case R, Rx, A, Intensity, Luminance:
		channels = 1
	case RG, RGx, RA:
		channels = 2
	case RGB, RGBx:
		channels = 3
	case RGBA, BGRA, ARGB:
		channels = 4
	}

	var channelBytes int
	switch i.ChannelType {
	case SnormInt8, UnormInt8, SignedInt8, UnsignedInt8:
		channelBytes = 1
	case SnormInt16, UnormInt16, UnormShort565, UnormShort555, SignedInt16, UnsignedInt16, HalfFloat:
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

// Creates an image object on the device.
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
	}, nil
}

// Creates a image object that is host accessible.
//
// Creates an uninitialized buffer on the host. This memory is not pageable and
// allows for DMA copies (which are faster).
func (c *Context) CreateHostImage(mf MemFlags, format ImageFormat, width, height, depth int) (*Image, error) {

	cFormat := clw.CreateImageFormat(clw.ChannelOrder(format.ChannelOrder), clw.ChannelType(format.ChannelType))

	var mem clw.Mem
	var err error
	if depth == 1 {
		mem, err = clw.CreateImage2D(c.id, clw.MemFlags(mf)|clw.MemAllocHostPointer, cFormat, clw.Size(width),
			clw.Size(height), 0, nil)
	} else {
		mem, err = clw.CreateImage3D(c.id, clw.MemFlags(mf)|clw.MemAllocHostPointer, cFormat, clw.Size(width),
			clw.Size(height), clw.Size(depth), 0, 0, nil)
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
	}, nil
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
	return int(count), err
}

// Enqueues a command to copy image objects.
func (cq *CommandQueue) EnqueueCopyImage(src, dst *Image, r *Rect, waitList []*Event, e *Event) error {

	var event *clw.Event
	if e != nil {
		event = &e.id
		e.Context = cq.Context
		e.CommandType = CommandCopyImage
		e.CommandQueue = cq
	}

	events := cq.createEvents(waitList)
	err := clw.EnqueueCopyImage(cq.id, src.id, dst.id, r.Src.origin(), r.Dst.origin(), r.region(), events, event)
	cq.releaseEvents(events)
	return err
}

// Enqueues a command to map a region of an image object into the host address
// space. Uses the Src and Region component of the rectangle.
func (cq *CommandQueue) EnqueueMapImage(i *Image, bc BlockingCall, flags MapFlags, r *Rect, waitList []*Event,
	e *Event) (*MappedImage, error) {

	var event *clw.Event
	if e != nil {
		event = &e.id
		e.Context = cq.Context
		e.CommandType = CommandMapImage
		e.CommandQueue = cq
	}

	var rowPitch, slicePitch clw.Size
	events := cq.createEvents(waitList)
	pointer, err := clw.EnqueueMapImage(cq.id, i.id, clw.Bool(bc), clw.MapFlags(flags), r.Src.origin(), r.region(),
		&rowPitch, &slicePitch, events, event)
	cq.releaseEvents(events)
	if err != nil {
		return nil, err
	}

	return &MappedImage{i, pointer, int64(rowPitch), int64(slicePitch)}, nil
}

// Enqueues a command to unmap a previously mapped image object.
func (cq *CommandQueue) EnqueueUnmapImage(mi *MappedImage, waitList []*Event, e *Event) error {

	var event *clw.Event
	if e != nil {
		event = &e.id
		e.Context = cq.Context
		e.CommandType = CommandUnmapMemoryObject
		e.CommandQueue = cq
	}

	events := cq.createEvents(waitList)
	err := clw.EnqueueUnmapMemObject(cq.id, mi.Image.id, mi.pointer, events, event)
	cq.releaseEvents(events)
	return err
}
