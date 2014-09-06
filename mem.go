package cl11

import (
	"errors"

	clw "github.com/rdwilliamson/clw11"
)

// ErrInvalidRect is returned when a rectangle has a negative value or isn't
// completely set to 2D or 3D. Such as a layout having only one of Origin[2] or
// SlicePitch non-zero; or the rectangle layout is 2D and a Rect.Region[2] isn't
// 1.
var ErrInvalidRect = errors.New("cl: invalid rectangle")

// RectLayout defines a rectangle layout in bytes. For a 2D image SlicePitch and
// Origin[2] must be zero. The offset of pixel X,Y,Z in bytes is
// (Z-Origin[2])*SlicePitch + (Y-Origin[1])*RowPitch + (X-Origin[0]).
type RectLayout struct {
	Origin     [3]int64
	RowPitch   int64
	SlicePitch int64
}

// valid validates the a rectangle layout.
func (rl *RectLayout) valid() bool {
	// No negative numbers.
	if rl.Origin[0] < 0 || rl.Origin[1] < 0 || rl.Origin[2] < 0 || rl.RowPitch < 0 || rl.SlicePitch < 0 {
		return false
	}

	// Ambigious wether it's a 2D or 3D image.
	if (rl.Origin[2] != 0 && rl.SlicePitch == 0) || (rl.Origin[2] == 0 && rl.SlicePitch != 0) {
		return false
	}

	return true
}

// dimensions returns 2 for a 2D image and 3 otherwise. It assumes the
// RectLayout is valid.
func (rl *RectLayout) dimensions() int {
	if rl.Origin[2] == 0 && rl.SlicePitch == 0 {
		return 2
	}
	return 3
}

// A source and destination rectangular region.
type Rect struct {

	// The source offset in bytes. The offset in bytes is Src.Origin[2] *
	// Src.SlicePitch + Src.Origin[1] * Src.RowPitch + Src.Origin[0].
	Src RectLayout

	// The destination offset in bytes. The offset in bytes is Dst.Origin[2] *
	// Dst.SlicePitch + Dst.Origin[1] * Dst.RowPitch + Dst.Origin[0].
	Dst RectLayout

	// The (width, height, depth) in pixels of the region being copied. For a 2D
	// image the depth value given by Region[2] must be 1.
	Region [3]int64
}

// valid validates the a rectangle.
func (r *Rect) valid() bool {
	return r.Src.valid() && r.Dst.valid() && r.Region[0] >= 0 && r.Region[1] >= 0 && r.Region[2] >= 1
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

func (r *Rect) width() clw.Size {
	return clw.Size(r.Region[0])
}

func (r *Rect) height() clw.Size {
	return clw.Size(r.Region[1])
}

func (r *Rect) depth() clw.Size {
	return clw.Size(r.Region[2])
}

func (r *Rect) srcBytes() int64 {
	result := r.Src.RowPitch * r.Region[0]
	if r.Src.SlicePitch > 0 {
		result *= r.Src.SlicePitch
	}
	return result
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
