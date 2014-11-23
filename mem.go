package cl11

import (
	clw "github.com/rdwilliamson/clw11"
)

// RectLayout defines a rectangle layout in bytes. For a 2D rectangle SlicePitch
// and Origin[2] must be zero. The offset of pixel X,Y,Z in bytes is
// (Z-Origin[2])*SlicePitch + (Y-Origin[1])*RowPitch + (X-Origin[0]).
type RectLayout struct {
	Origin     [3]int64
	RowPitch   int64
	SlicePitch int64
}

func (rl *RectLayout) origin() [3]clw.Size {
	return [3]clw.Size{clw.Size(rl.Origin[0]), clw.Size(rl.Origin[1]), clw.Size(rl.Origin[2])}
}

func (rl *RectLayout) rowPitch() clw.Size {
	return clw.Size(rl.RowPitch)
}

func (rl *RectLayout) slicePitch() clw.Size {
	return clw.Size(rl.SlicePitch)
}

// A source and destination rectangular region.
type Rect struct {

	// The source offset in bytes. The offset in bytes is Src.Origin[2] *
	// Src.SlicePitch + Src.Origin[1] * Src.RowPitch + Src.Origin[0].
	Src RectLayout

	// The destination offset in bytes. The offset in bytes is Dst.Origin[2] *
	// Dst.SlicePitch + Dst.Origin[1] * Dst.RowPitch + Dst.Origin[0].
	Dst RectLayout

	// The (width, height, depth) in bytes of the region being copied. For a 2D
	// rectangle the depth value given by Region[2] must be 1.
	Region [3]int64
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

func (r *Rect) size() uintptr {
	return uintptr(r.Region[0] * r.Region[1] * r.Region[2])
}

type MemFlags uint

// Bit field.
const (
	MemReadWrite = MemFlags(clw.MemReadWrite)
	MemWriteOnly = MemFlags(clw.MemWriteOnly)
	MemReadOnly  = MemFlags(clw.MemReadOnly)
)

type MapFlags uint

// Bit field.
const (
	MapRead  = MapFlags(clw.MapRead)
	MapWrite = MapFlags(clw.MapWrite)
)

type BlockingCall clw.Bool

const (
	Blocking    = BlockingCall(clw.True)
	NonBlocking = BlockingCall(clw.False)
)

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

	events := cq.createEvents(waitList)
	err := clw.EnqueueCopyImageToBuffer(cq.id, src.id, dst.id, r.Src.origin(), r.region(), clw.Size(offset), events,
		event)
	cq.releaseEvents(events)
	return err
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

	events := cq.createEvents(waitList)
	err := clw.EnqueueCopyBufferToImage(cq.id, src.id, dst.id, clw.Size(offset), r.Dst.origin(), r.region(), events,
		event)
	cq.releaseEvents(events)
	return err
}
