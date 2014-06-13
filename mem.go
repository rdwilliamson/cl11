package cl11

import (
	clw "github.com/rdwilliamson/clw11"
)

// Rectangle layout in bytes. For a 2D image SlicePitch and Origin[2] are zero.
// The offset of pixel X,Y,Z in bytes is (Z-Origin[2])*SlicePitch +
// (Y-Origin[1])*RowPitch + (X-Origin[3]).
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

// Only the source and region are used from the rectangle.
func (cq *CommandQueue) EnqueueCopyImageToBuffer(src *Image, dst *Buffer, r *Rect, offset int, waitList []*Event,
	e *Event) error {

	var event *clw.Event
	if e != nil {
		event = &e.id
		e.Context = cq.Context
		e.CommandType = CommandCopyImageToBuffer
		e.CommandQueue = cq
	}

	return clw.EnqueueCopyImageToBuffer(cq.id, src.id, dst.id, r.srcOrigin(), r.region(), clw.Size(offset),
		cq.toEvents(waitList), event)
}

// Only the destination and region are used from the rectangle.
func (cq *CommandQueue) EnqueueCopyBufferToImage(src *Buffer, dst *Image, offset int, r *Rect, waitList []*Event,
	e *Event) error {

	var event *clw.Event
	if e != nil {
		event = &e.id
		e.Context = cq.Context
		e.CommandType = CommandCopyBufferToImage
		e.CommandQueue = cq
	}

	return clw.EnqueueCopyBufferToImage(cq.id, src.id, dst.id, clw.Size(offset), r.dstOrigin(), r.region(),
		cq.toEvents(waitList), event)
}
