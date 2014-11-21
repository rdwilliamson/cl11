package cl11

import (
	"errors"
)

// NotAddressable is returned when attempting to read from or write to a value
// that is not addressable (see reflect.Value.CanAddr).
var NotAddressable = errors.New("cl: not addressable")

// UnsupportedImageFormat is returned when trying to use an unsupported Go image
// format in one of the convenience image methods (ends something along the
// lines on By/From Image).
var UnsupportedImageFormat = errors.New("cl: unsupported image format")
