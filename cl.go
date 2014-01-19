// A wrapper package that attempts to map the OpenCL 1.1 C API to idiomatic Go.
package cl11

import (
	"unsafe"

	clw "github.com/rdwilliamson/clw11"
)

func Flush(cq CommandQueue) error {
	return clw.Flush(clw.CommandQueue(cq.ID))
}

func Finish(cq CommandQueue) error {
	return clw.Finish(clw.CommandQueue(cq.ID))
}

func toEvents(in []Event) []clw.Event {
	return *(*[]clw.Event)(unsafe.Pointer(&in))
}
