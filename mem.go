package cl11

import (
	clw "github.com/rdwilliamson/clw11"
)

// Rectangle layout in bytes. For a 2D image SlicePitch and Origin[2] are zero.
type RectLayout struct {
	Origin     [3]int64
	RowPitch   int64
	SlicePitch int64
}

// A source and destination rectangular region.
type Rect struct {

	// The source offset in bytes. For a 2D rectangle Src.Origin[2] should be 0.
	// The offset in bytes is Src.Origin[2] * Src.SlicePitch + Src.Origin[1] *
	// Src.RowPitch + Src.Origin[0].
	Src RectLayout

	// The destination offset in bytes. For a 2D rectangle Dst.Origin[2] should
	// be 0. The offset in bytes is Dst.Origin[2] * Dst.SlicePitch +
	// Dst.Origin[1] * Dst.RowPitch + Dst.Origin[0].
	Dst RectLayout

	// The (width, height, depth) in bytes of the 2D or 3D rectangle being
	// copied. For a 2D rectangle the depth value given by Region[2] should be
	// 1.
	Region [3]int64
}

type MemFlags int

// Bitfield.
const (
	MemReadWrite = MemFlags(clw.MemReadWrite)
	MemWriteOnly = MemFlags(clw.MemWriteOnly)
	MemReadOnly  = MemFlags(clw.MemReadOnly)
)

type MapFlags int

// Bitfield.
const (
	MapRead  = MapFlags(clw.MapRead)
	MapWrite = MapFlags(clw.MapWrite)
)

type BlockingCall clw.Bool

const (
	Blocking    = BlockingCall(clw.True)
	NonBlocking = BlockingCall(clw.False)
)

func (r *Rect) srcOrigin() [3]clw.Size {
	return [3]clw.Size{clw.Size(r.Src.Origin[0]), clw.Size(r.Src.Origin[1]), clw.Size(r.Src.Origin[2])}
}

func (r *Rect) srcRowPitch() clw.Size {
	return clw.Size(r.Src.RowPitch)
}

func (r *Rect) srcSlicePitch() clw.Size {
	return clw.Size(r.Src.SlicePitch)
}

func (r *Rect) dstOrigin() [3]clw.Size {
	return [3]clw.Size{clw.Size(r.Dst.Origin[0]), clw.Size(r.Dst.Origin[1]), clw.Size(r.Dst.Origin[2])}
}

func (r *Rect) dstRowPitch() clw.Size {
	return clw.Size(r.Dst.RowPitch)
}

func (r *Rect) dstSlicePitch() clw.Size {
	return clw.Size(r.Dst.SlicePitch)
}

func (r *Rect) region() [3]clw.Size {
	return [3]clw.Size{clw.Size(r.Region[0]), clw.Size(r.Region[1]), clw.Size(r.Region[2])}
}
