package cl11

import (
	"errors"

	clw "github.com/rdwilliamson/clw11"
)

// ErrNotAddressable is returned when attempting to read from or write to a
// value that is not addressable (see reflect.Value.CanAddr).
var ErrNotAddressable = errors.New("cl: not addressable")

// ErrUnsupportedImageFormat is returned when trying to use one of the
// convenience image methods and there is a mismatch or incompatibility between
// the Go image and the OpenCL image.
var ErrUnsupportedImageFormat = errors.New("cl: unsupported image format")

var (
	DeviceNotFound                     = clw.DeviceNotFound
	DeviceNotAvailable                 = clw.DeviceNotAvailable
	CompilerNotAvailable               = clw.CompilerNotAvailable
	MemObjectAllocationFailure         = clw.MemObjectAllocationFailure
	OutOfResources                     = clw.OutOfResources
	OutOfHostMemory                    = clw.OutOfHostMemory
	ProfilingInfoNotAvailable          = clw.ProfilingInfoNotAvailable
	MemCopyOverlap                     = clw.MemCopyOverlap
	ImageFormatMismatch                = clw.ImageFormatMismatch
	ImageFormatNotSupported            = clw.ImageFormatNotSupported
	BuildProgramFailure                = clw.BuildProgramFailure
	MapFailure                         = clw.MapFailure
	MisalignedSubBufferOffset          = clw.MisalignedSubBufferOffset
	ExecStatusErrorForEventsInWaitList = clw.ExecStatusErrorForEventsInWaitList
)

var (
	InvalidValue                 = clw.InvalidValue
	InvalidDeviceType            = clw.InvalidDeviceType
	InvalidPlatform              = clw.InvalidPlatform
	InvalidDevice                = clw.InvalidDevice
	InvalidContext               = clw.InvalidContext
	InvalidQueueProperties       = clw.InvalidQueueProperties
	InvalidCommandQueue          = clw.InvalidCommandQueue
	InvalidHostPtr               = clw.InvalidHostPtr
	InvalidMemObject             = clw.InvalidMemObject
	InvalidImageFormatDescriptor = clw.InvalidImageFormatDescriptor
	InvalidImageSize             = clw.InvalidImageSize
	InvalidSampler               = clw.InvalidSampler
	InvalidBinary                = clw.InvalidBinary
	InvalidBuildOptions          = clw.InvalidBuildOptions
	InvalidProgram               = clw.InvalidProgram
	InvalidProgramExecutable     = clw.InvalidProgramExecutable
	InvalidKernelName            = clw.InvalidKernelName
	InvalidKernelDefinition      = clw.InvalidKernelDefinition
	InvalidKernel                = clw.InvalidKernel
	InvalidArgIndex              = clw.InvalidArgIndex
	InvalidArgValue              = clw.InvalidArgValue
	InvalidArgSize               = clw.InvalidArgSize
	InvalidKernelArgs            = clw.InvalidKernelArgs
	InvalidWorkDimension         = clw.InvalidWorkDimension
	InvalidWorkGroupSize         = clw.InvalidWorkGroupSize
	InvalidWorkItemSize          = clw.InvalidWorkItemSize
	InvalidGlobalOffset          = clw.InvalidGlobalOffset
	InvalidEventWaitList         = clw.InvalidEventWaitList
	InvalidEvent                 = clw.InvalidEvent
	InvalidOperation             = clw.InvalidOperation
	InvalidGlObject              = clw.InvalidGlObject
	InvalidBufferSize            = clw.InvalidBufferSize
	InvalidMipLevel              = clw.InvalidMipLevel
	InvalidGlobalWorkSize        = clw.InvalidGlobalWorkSize
	InvalidProperty              = clw.InvalidProperty
)
