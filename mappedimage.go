package cl11

import (
	"unsafe"

	clw "github.com/rdwilliamson/clw11"
)

// A mapped image. Has convenience functions for common data types.
type MappedImage struct {
	pointer    unsafe.Pointer
	RowPitch   int64 // Scan line width in bytes.
	SlicePitch int64 // The size in bytes of each 2D image (0 for 2D image).
	memID      clw.Mem
}
