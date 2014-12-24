package cl11

import (
	"errors"
	"io"
	"reflect"
	"unsafe"
)

// MappedBuffer is a mapping of an OpenCL buffer into the host address space. It
// implements a number of IO interfaces and can be converted to slices of the
// numeric types.
type MappedBuffer struct {
	Buffer  *Buffer
	pointer unsafe.Pointer
	index   int64
	size    int64
}

// Read reads the next len(b) bytes from the buffer or until the buffer is
// drained. The return value is the number of bytes read. If the buffer has no
// data to return the error is io.EOF (unless len(b) is zero); otherwise it is
// nil.
func (mb *MappedBuffer) Read(b []byte) (int, error) {
	if len(b) == 0 {
		return 0, nil
	}
	if mb.index >= mb.size {
		return 0, io.EOF
	}
	n := copy(b, mb.Bytes()[mb.index:])
	mb.index += int64(n)
	return n, nil
}

// ReadAt reads len(b) bytes into b starting at offset off in the underlying
// input source. It returns the number of bytes read (0 <= n <= len(b)) and any
// error encountered. When ReadAt returns n < len(b), it returns io.EOF.
func (mb *MappedBuffer) ReadAt(b []byte, offset int64) (int, error) {
	if len(b) == 0 {
		return 0, nil
	}
	if offset < 0 {
		return 0, errors.New("cl: MappedBuffer.ReadAt: negative offset")
	}
	if offset >= mb.size {
		return 0, io.EOF
	}
	n := copy(b, mb.Bytes()[offset:])
	if n < len(b) {
		return n, io.EOF
	}
	return n, nil
}

// Seek sets the offset for the next Read or Write to offset, interpreted
// according to whence: 0 means relative to the origin, 1 means relative to the
// current offset, and 2 means relative to the end. Seek returns the new offset
// and an error, if any.
//
// Seeking to a negative offset is an error. Seeking to any positive offset is
// legal, but an offset larger than the buffer's size will result in io.EOF for
// reads and ErrBufferFull for writes.
func (mb *MappedBuffer) Seek(offset int64, whence int) (int64, error) {
	var newIndex int64
	switch whence {
	case 0:
		newIndex = offset
	case 1:
		newIndex = mb.index + offset
	case 2:
		newIndex = mb.size + offset
	default:
		return 0, errors.New("cl: MappedBuffer.Seek: invalid whence")
	}
	if newIndex < 0 {
		return 0, errors.New("cl: MappedBuffer.Seek: negative position")
	}
	mb.index = newIndex
	return 0, nil
}

// Write writes the next len(b) bytes from the buffer or as much as possible
// until the buffer is full. The return value is the number of bytes written. If
// the buffer has no room to write all the data, err is ErrBufferFull (unless
// len(b) is zero); otherwise it is nil.
func (mb *MappedBuffer) Write(b []byte) (int, error) {
	if len(b) == 0 {
		return 0, nil
	}
	if mb.index >= mb.size {
		return 0, ErrBufferFull
	}
	n := copy(mb.Bytes()[mb.index:], b)
	mb.index += int64(n)
	if n < len(b) {
		return n, ErrBufferFull
	}
	return n, nil
}

// WriteAt writes len(b) bytes from b to the buffer at offset. It returns the
// number of bytes written from b (0 <= n <= len(b)). WriteAt returns
// ErrBufferFull if it wrote less than len(b).
func (mb *MappedBuffer) WriteAt(b []byte, offset int64) (int, error) {
	if len(b) == 0 {
		return 0, nil
	}
	if offset < 0 {
		return 0, errors.New("cl: MappedBuffer.WriteAt: negative offset")
	}
	if offset >= mb.size {
		return 0, ErrBufferFull
	}
	n := copy(mb.Bytes()[offset:], b)
	if n < len(b) {
		return n, ErrBufferFull
	}
	return n, nil
}

// WriteTo writes data to w until the buffer is drained or an error occurs. The
// return value is the number of bytes written; it always fits into an int, but
// it is int64 to match the io.WriterTo interface. Any error encountered during
// the write is also returned.
func (mb *MappedBuffer) WriteTo(w io.Writer) (int64, error) {
	if mb.index >= mb.size {
		return 0, nil
	}
	b := mb.Bytes()[mb.index:]
	m, err := w.Write(b)
	if m > len(b) {
		panic("cl: MappedBuffer.WriteTo: invalid Write count")
	}
	n := int64(m)
	mb.index += n
	if err != nil {
		return n, err
	}
	if m != len(b) {
		return n, io.ErrShortWrite
	}
	return n, nil
}

// Returns a slice of bytes backed by the mapped buffer.
func (mb *MappedBuffer) Bytes() []byte {
	var header reflect.SliceHeader
	header.Data = uintptr(mb.pointer)
	size := int(mb.size)
	header.Len = size
	header.Cap = size
	return *(*[]byte)(unsafe.Pointer(&header))
}

// Returns a slice of float32s backed by the mapped buffer.
func (mb *MappedBuffer) Float32s() []float32 {
	var header reflect.SliceHeader
	header.Data = uintptr(mb.pointer)
	size := int(mb.size / int64(4))
	header.Len = size
	header.Cap = size
	return *(*[]float32)(unsafe.Pointer(&header))
}

var _ io.Reader = (*MappedBuffer)(nil)
var _ io.ReaderAt = (*MappedBuffer)(nil)
var _ io.Seeker = (*MappedBuffer)(nil)
var _ io.Writer = (*MappedBuffer)(nil)
var _ io.WriterAt = (*MappedBuffer)(nil)
var _ io.WriterTo = (*MappedBuffer)(nil)
