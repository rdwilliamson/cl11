package cl11

import (
	clw "github.com/rdwilliamson/clw11"
)

type Image struct {
	id clw.Mem

	Format ImageFormat

	Width int

	Height int

	Depth int

	RowPitch int

	SlicePitch int

	Flags MemFlags
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
	MemObjectImage2D = MemObjectType(clw.MemObjectImage2D)
	MemObjectImage3D = MemObjectType(clw.MemObjectImage3D)
)

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
	Float          = ChannelType(clw.Float)
)

func (c *Context) GetSupportedImageFormats(mf MemFlags, t MemObjectType) ([]ImageFormat, error) {

	var count clw.Uint
	err := clw.GetSupportedImageFormats(c.id, clw.MemFlags(mf), clw.MemObjectType(t), 0, nil, &count)
	if err != nil {
		return nil, err
	}

	formats := make([]clw.ImageFormat, int(count))

	err = clw.GetSupportedImageFormats(c.id, clw.MemFlags(mf), clw.MemObjectType(t), count, &formats[0], nil)
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

func (c *Context) CreateImage2D(flags MemFlags, fmt ImageFormat, width, height, pitch int,
	src interface{}) (*Image, error) {

	fmt2 := clw.CreateImageFormat(clw.ChannelOrder(fmt.ChannelOrder), clw.ChannelType(fmt.ChannelType))

	mem, err := clw.CreateImage2D(c.id, clw.MemFlags(flags), fmt2, clw.Size(width), clw.Size(height), clw.Size(pitch),
		nil)
	if err != nil {
		return nil, err
	}

	return &Image{id: mem, Format: fmt}, nil
}
