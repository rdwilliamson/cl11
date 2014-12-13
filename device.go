package cl11

import (
	"encoding/binary"
	"strings"
	"unsafe"

	clw "github.com/rdwilliamson/clw11"
)

// Device configuration information, also used to create Contexts.
type Device struct {
	id clw.DeviceID

	// The platform from which the device was created.
	Platform *Platform

	// A device type can be a CPU, a host processor that runs the OpenCL
	// implementation. A GPU, a device can also be used to accelerate a 3D API
	// such as OpenGL or DirectX. A dedicated OpenCL accelerator (for example
	// the IBM CELL Blade), these devices communicate with the host processor
	// using a peripheral interconnect such as PCIe.
	Type DeviceType

	// True if the device is available.
	Available bool

	// False if the implementation does not have a compiler available to compile
	// the program source. Can only be false on an embedded platform.
	CompilerAvailable bool

	// Either binary.BigEndian or binary.LittleEndian.
	ByteOrder binary.ByteOrder

	// True if the device implements error correction for all accesses to
	// compute device memory (global and constant).
	ErrorCorrectionSupport bool

	// True if if images are supported.
	ImageSupport bool

	// True if the device and the host have a unified memory subsystem.
	UnifiedHostMemory bool

	// The default compute device address space size in bits. Currently
	// supported values are 32 or 64 bits.
	AddressBits int

	// Size of global memory cache line in bytes.
	GlobalMemCachelineSize int

	// Maximum configured clock frequency of the device in MHz.
	MaxClockFrequency int

	// The number of parallel compute cores on the OpenCL device. The minimum
	// value is 1.
	MaxComputeUnits int

	// Max number of arguments declared with the __constant qualifier in a
	// kernel. The minimum value is 8.
	MaxConstantArgs int

	// Max number of simultaneous image objects that can be read by a kernel.
	// The minimum value is 128 if ImageSupport is true.
	MaxReadImageArgs int

	// Maximum number of samplers that can be used in a kernel. The minimum
	// value is 16 if ImageSupport is true.
	MaxSamplers int

	// Maximum dimensions that specify the global and local work-item IDs used
	// by the data parallel execution model. (Refer to EnqueueNDRangeKernel).
	// The minimum value is 3.
	MaxWorkItemDimensions int

	// Max number of simultaneous image objects that can be written to by a
	// kernel. The minimum value is 8 if ImageSupport is true.
	MaxWriteImageArgs int

	// Describes the alignment in bits of the base address of any allocated
	// memory object.
	MemBaseAddrAlign int

	// The smallest alignment in bytes which can be used for any data type.
	MinDataTypeAlignSize int

	// Preferred native vector width size for built-in scalar types that can be
	// put into vectors. The vector width is defined as the number of scalar
	// elements that can be stored in the vector.
	PreferredVectorWidths VectorWidths

	// Returns the native ISA vector width. The vector width is defined as the
	// number of scalar elements that can be stored in the vector.
	NativeVectorWidths VectorWidths

	// Extensions supported by the device.
	Extensions []string

	// Device name string.
	Name string

	// Either FulleProfile, if the device supports the OpenCL specification
	// (functionality defined as part of the core specification and does not
	// require any extensions to be supported). Or EmbedddedProfile, if the
	// device supports the OpenCL embedded profile.
	Profile Profile

	// Vendor name string.
	Vendor string

	// A unique device vendor identifier. An example of a unique device
	// identifier could be the PCIe ID.
	VendorID int

	// The OpenCL version supported by the device.
	Version Version

	// The highest OpenCL C version supported by the compiler for this device.
	OpenCLCVersion Version

	// OpenCL software driver version.
	DriverVersion Version

	// Size of global memory cache in bytes.
	GlobalMemCacheSize int64

	// Size of global device memory in bytes.
	GlobalMemSize int64

	// Size of local memory arena in bytes. The minimum value is 32 KB.
	LocalMemSize int64

	// Max size in bytes of a constant buffer allocation. The minimum value is
	// 64 KB.
	MaxConstantBufferSize int64

	// Max size of memory object allocation in bytes. The minimum value is
	// max(1/4th of GlobalMemSize, 128*1024*1024).
	MaxMemAllocSize int64

	// Max height of 2D image in pixels. The minimum value is 8192 if
	// ImageSupport is true.
	Image2DMaxHeight int

	// Max width of 2D image in pixels. The minimum value is 8192 if
	// ImageSupport is true.
	Image2DMaxWidth int

	// Max depth of 3D image in pixels. The minimum value is 2048 if
	// ImageSupport is true.
	Image3DMaxDepth int

	// Max height of 3D image in pixels. The minimum value is 2048 if
	// ImageSupport is true.
	Image3DMaxHeight int

	// Max width of 3D image in pixels. The minimum value is 2048 if
	// ImageSupport is true.
	Image3DMaxWidth int

	// Max size in bytes of the arguments that can be passed to a kernel. The
	// minimum value is 1024. For this minimum value, only a maximum of 128
	// arguments can be passed to a kernel.
	MaxParameterSize int

	// Maximum number of work-items in a work-group executing a kernel using the
	// data parallel execution model. (Refer to EnqueueNDRangeKernel). The
	// minimum value is 1.
	MaxWorkGroupSize int

	// Describes the resolution of device timer. This is measured in nanoseconds.
	ProfilingTimerResolution int

	// Maximum number of work-items that can be specified in each dimension of
	// the work-group to EnqueueNDRangeKernel. Returns n entries, where n is the
	// value returned by the query for MaxWorkItemDimensions. The minimum value
	// is (1, 1, 1).
	MaxWorkItemSizes []int

	// Describes single precision floating-point capability of the device. The
	// mandated minimum capability is round to nearest and infinity and NaN
	// supported. Basic operations can be implemented in software (soft floats).
	SingleFpConfig FPConfig

	// Describes optional double precision floating-point capability of the
	// device. The mandated minimum capability is fused multiply-add, round the
	// nearest, round to zero, round to infinity, infinity and NaN support, and
	// denormalized numbers.
	DoubleFpConfig FPConfig

	// HalfFpConfig   FPConfig

	// Describes the execution capabilities of the device. This is a bit-field
	// that describes one or more of the following values: ExecKernel - The
	// OpenCL device can execute OpenCL kernels. ExecNativeKernel - The OpenCL
	// device can execute native kernels. The mandated minimum capability is
	// ExecKernel.
	ExecCapabilities ExecCapabilities

	// Describes the command-queue properties supported by the device. This is a
	// bit-field that describes one or more of the following values:
	// QueueOutOfOrderExecution. QueueProfilingEnable. These
	// properties are described in the table for CreateCommandQueue. The
	// mandated minimum capability is QueueProfilingEnable.
	CommandQueueProperties CommandQueueProperties

	// Type of global memory cache supported. Valid values are: None,
	// ReadOnlyCache, and ReadWriteCache.
	GlobalMemCacheType GlobalMemCacheType

	// Type of local memory supported. This can be set to Local implying
	// dedicated local memory storage such as SRAM, or Global.
	LocalMemTypeInfo LocalMemTypeInfo
}

type DeviceType uint

// Bit field.
const (
	DeviceTypeDefault     = DeviceType(clw.DeviceTypeDefault)
	DeviceTypeCpu         = DeviceType(clw.DeviceTypeCpu)
	DeviceTypeGpu         = DeviceType(clw.DeviceTypeGpu)
	DeviceTypeAccelerator = DeviceType(clw.DeviceTypeAccelerator)
	DeviceTypeAll         = DeviceType(clw.DeviceTypeAll)
)

func (dt DeviceType) String() string {
	switch dt {
	case DeviceTypeDefault:
		return "default"
	case DeviceTypeCpu:
		return "CPU"
	case DeviceTypeGpu:
		return "GPU"
	case DeviceTypeAccelerator:
		return "accelerator"
	case DeviceTypeAll:
		return "all"
	}
	return ""
}

type VectorWidths struct {
	Char   int
	Short  int
	Int    int
	Long   int
	Float  int
	Double int
	// Half   int
}

type FPConfig uint

// Bit field.
const (
	FPDenorm         = FPConfig(clw.FPDenorm)
	FPFma            = FPConfig(clw.FPFma)
	FPInfNan         = FPConfig(clw.FPInfNan)
	FPRoundToInf     = FPConfig(clw.FPRoundToInf)
	FPRoundToNearest = FPConfig(clw.FPRoundToNearest)
	FPRoundToZero    = FPConfig(clw.FPRoundToZero)
	FPSoftFloat      = FPConfig(clw.FPSoftFloat)
)

func (fpConfig FPConfig) String() string {
	var configStrings []string
	if fpConfig&FPDenorm != 0 {
		configStrings = append(configStrings, "denormalized numbers")
	}
	if fpConfig&FPFma != 0 {
		configStrings = append(configStrings, "fused multiply-add")
	}
	if fpConfig&FPInfNan != 0 {
		configStrings = append(configStrings, "inf and nan")
	}
	if fpConfig&FPRoundToInf != 0 {
		configStrings = append(configStrings, "round to inf")
	}
	if fpConfig&FPRoundToNearest != 0 {
		configStrings = append(configStrings, "round to nearest")
	}
	if fpConfig&FPRoundToZero != 0 {
		configStrings = append(configStrings, "round to zero")
	}
	if fpConfig&FPSoftFloat != 0 {
		configStrings = append(configStrings, "soft float")
	}
	return "{" + strings.Join(configStrings, ", ") + "}"
}

type GlobalMemCacheType int

const (
	None           = GlobalMemCacheType(clw.None)
	ReadOnlyCache  = GlobalMemCacheType(clw.ReadOnlyCache)
	ReadWriteCache = GlobalMemCacheType(clw.ReadWriteCache)
)

func (gmct GlobalMemCacheType) String() string {
	switch gmct {
	case None:
		return "none"
	case ReadOnlyCache:
		return "read only"
	case ReadWriteCache:
		return "read and write"
	}
	return ""
}

type LocalMemTypeInfo int

const (
	Global = LocalMemTypeInfo(clw.Global)
	Local  = LocalMemTypeInfo(clw.Local)
)

func (lmti LocalMemTypeInfo) String() string {
	switch lmti {
	case Global:
		return "global"
	case Local:
		return "local"
	}
	return ""
}

type ExecCapabilities uint

// Bit field.
const (
	ExecKernel       = ExecCapabilities(clw.ExecKernel)
	ExecNativeKernel = ExecCapabilities(clw.ExecNativeKernel)
)

func (ec ExecCapabilities) String() string {
	var execStrings []string
	if ec&ExecKernel != 0 {
		execStrings = append(execStrings, "kernel")
	}
	if ec&ExecNativeKernel != 0 {
		execStrings = append(execStrings, "native kernel")
	}
	return "{" + strings.Join(execStrings, ", ") + "}"
}

func (p *Platform) getDevices() error {

	var numEntries clw.Uint
	err := clw.GetDeviceIDs(p.id, clw.DeviceTypeAll, 0, nil, &numEntries)
	if err != nil {
		return err
	}

	deviceIDs := make([]clw.DeviceID, numEntries)
	err = clw.GetDeviceIDs(p.id, clw.DeviceTypeAll, numEntries, &deviceIDs[0], nil)
	if err != nil {
		return err
	}

	p.Devices = make([]*Device, len(deviceIDs))
	for i := range p.Devices {

		device := &Device{id: deviceIDs[i], Platform: p}

		err = device.getAllInfo()
		if err != nil {
			return err
		}

		p.Devices[i] = device
	}

	return nil
}

func (d *Device) getAllInfo() (err error) {

	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()

	d.Available = d.getBool(clw.DeviceAvailable)
	d.CompilerAvailable = d.getBool(clw.DeviceCompilerAvailable)
	d.ErrorCorrectionSupport = d.getBool(clw.DeviceErrorCorrectionSupport)
	d.ImageSupport = d.getBool(clw.DeviceImageSupport)
	d.UnifiedHostMemory = d.getBool(clw.DeviceHostUnifiedMemory)

	d.AddressBits = d.getUint(clw.DeviceAddressBits)
	d.GlobalMemCachelineSize = d.getUint(clw.DeviceGlobalMemCachelineSize)
	d.MaxClockFrequency = d.getUint(clw.DeviceMaxClockFrequency)
	d.MaxComputeUnits = d.getUint(clw.DeviceMaxComputeUnits)
	d.MaxConstantArgs = d.getUint(clw.DeviceMaxConstantArgs)
	d.MaxReadImageArgs = d.getUint(clw.DeviceMaxReadImageArgs)
	d.MaxSamplers = d.getUint(clw.DeviceMaxSamplers)
	d.MaxWorkItemDimensions = d.getUint(clw.DeviceMaxWorkItemDimensions)
	d.MaxWriteImageArgs = d.getUint(clw.DeviceMaxWriteImageArgs)
	d.MemBaseAddrAlign = d.getUint(clw.DeviceMemBaseAddrAlign)
	d.MinDataTypeAlignSize = d.getUint(clw.DeviceMinDataTypeAlignSize)
	d.VendorID = d.getUint(clw.DeviceVendorID)

	d.PreferredVectorWidths.Char = d.getUint(clw.DevicePreferredVectorWidthChar)
	d.PreferredVectorWidths.Short = d.getUint(clw.DevicePreferredVectorWidthShort)
	d.PreferredVectorWidths.Int = d.getUint(clw.DevicePreferredVectorWidthInt)
	d.PreferredVectorWidths.Long = d.getUint(clw.DevicePreferredVectorWidthLong)
	d.PreferredVectorWidths.Float = d.getUint(clw.DevicePreferredVectorWidthFloat)
	d.PreferredVectorWidths.Double = d.getUint(clw.DevicePreferredVectorWidthDouble)
	// d.PreferredVectorWidths.Half = d.getUint(clw.DevicePreferredVectorWidthHalf)
	d.NativeVectorWidths.Char = d.getUint(clw.DeviceNativeVectorWidthChar)
	d.NativeVectorWidths.Short = d.getUint(clw.DeviceNativeVectorWidthShort)
	d.NativeVectorWidths.Int = d.getUint(clw.DeviceNativeVectorWidthInt)
	d.NativeVectorWidths.Long = d.getUint(clw.DeviceNativeVectorWidthLong)
	d.NativeVectorWidths.Float = d.getUint(clw.DeviceNativeVectorWidthFloat)
	d.NativeVectorWidths.Double = d.getUint(clw.DeviceNativeVectorWidthDouble)
	// d.NativeVectorWidths.Half = d.getUint(clw.DeviceNativeVectorWidthHalf)

	d.Extensions = strings.Fields(d.getString(clw.DeviceExtensions))

	if d.getBool(clw.DeviceEndianLittle) {
		d.ByteOrder = binary.LittleEndian
	} else {
		d.ByteOrder = binary.BigEndian
	}

	d.Name = d.getString(clw.DeviceName)
	d.Vendor = d.getString(clw.DeviceVendor)

	d.Profile = toProfile(d.getString(clw.DeviceProfile))

	d.Version = toVersion(d.getString(clw.DeviceVersion))
	d.DriverVersion = toVersion(d.getString(clw.DriverVersion))

	d.OpenCLCVersion = toVersion(d.getString(clw.DeviceOpenCLCVersion))

	d.GlobalMemCacheSize = d.getUlong(clw.DeviceGlobalMemCacheSize)
	d.GlobalMemSize = d.getUlong(clw.DeviceGlobalMemSize)
	d.LocalMemSize = d.getUlong(clw.DeviceLocalMemSize)
	d.MaxConstantBufferSize = d.getUlong(clw.DeviceMaxConstantBufferSize)
	d.MaxMemAllocSize = d.getUlong(clw.DeviceMaxMemAllocSize)

	d.Image2DMaxHeight = d.getSize(clw.DeviceImage2dMaxHeight)
	d.Image2DMaxWidth = d.getSize(clw.DeviceImage2dMaxWidth)
	d.Image3DMaxDepth = d.getSize(clw.DeviceImage3dMaxDepth)
	d.Image3DMaxHeight = d.getSize(clw.DeviceImage3dMaxHeight)
	d.Image3DMaxWidth = d.getSize(clw.DeviceImage3dMaxWidth)
	d.MaxParameterSize = d.getSize(clw.DeviceMaxParameterSize)
	d.MaxWorkGroupSize = d.getSize(clw.DeviceMaxWorkGroupSize)
	d.ProfilingTimerResolution = d.getSize(clw.DeviceProfilingTimerResolution)

	d.MaxWorkItemSizes = d.getSizeArray(clw.DeviceMaxWorkItemSizes)

	d.SingleFpConfig = d.getFpConfig(clw.DeviceSingleFpConfig)
	d.DoubleFpConfig = d.getFpConfig(clw.DeviceDoubleFpConfig)
	// d.HalfFpConfig = d.getFpConfig(clw.DeviceHalfFpConfig)

	d.ExecCapabilities = d.getExecCapabilities(clw.DeviceExecutionCapabilities)

	d.Type = d.getType(clw.DeviceTypeInfo)

	d.CommandQueueProperties = d.getCommandQueueProperties(clw.DeviceQueueProperties)

	d.GlobalMemCacheType = d.getGlobalMemCacheType(clw.DeviceGlobalMemCacheType)

	d.LocalMemTypeInfo = d.getLocalMemTypeInfo(clw.DeviceLocalMemTypeInfo)

	return
}

func (d *Device) getBool(paramName clw.DeviceInfo) bool {
	var paramValue clw.Bool
	err := clw.GetDeviceInfo(d.id, paramName, clw.Size(unsafe.Sizeof(paramValue)),
		unsafe.Pointer(&paramValue), nil)
	if err != nil {
		panic(err)
	}
	return paramValue != clw.False
}

func (d *Device) getUint(paramName clw.DeviceInfo) int {
	var paramValue clw.Uint
	err := clw.GetDeviceInfo(d.id, paramName, clw.Size(unsafe.Sizeof(paramValue)),
		unsafe.Pointer(&paramValue), nil)
	if err != nil {
		panic(err)
	}
	return int(paramValue)
}

func (d *Device) getString(paramName clw.DeviceInfo) string {

	var paramValueSize clw.Size
	err := clw.GetDeviceInfo(d.id, paramName, 0, nil, &paramValueSize)
	if err != nil {
		panic(err)
	}

	buffer := make([]byte, paramValueSize)
	err = clw.GetDeviceInfo(d.id, paramName, paramValueSize, unsafe.Pointer(&buffer[0]), nil)
	if err != nil {
		panic(err)
	}

	// Trim space and trailing \0.
	return strings.TrimSpace(string(buffer[:len(buffer)-1]))
}

func (d *Device) getUlong(paramName clw.DeviceInfo) int64 {
	var paramValue clw.Ulong
	err := clw.GetDeviceInfo(d.id, paramName, clw.Size(unsafe.Sizeof(paramValue)),
		unsafe.Pointer(&paramValue), nil)
	if err != nil {
		panic(err)
	}
	return int64(paramValue)
}

func (d *Device) getSize(paramName clw.DeviceInfo) int {
	var paramValue clw.Size
	err := clw.GetDeviceInfo(d.id, paramName, clw.Size(unsafe.Sizeof(paramValue)),
		unsafe.Pointer(&paramValue), nil)
	if err != nil {
		panic(err)
	}
	return int(paramValue)
}

func (d *Device) getSizeArray(paramName clw.DeviceInfo) []int {

	var paramValueSize clw.Size
	err := clw.GetDeviceInfo(d.id, paramName, 0, nil, &paramValueSize)
	if err != nil {
		panic(err)
	}

	var a clw.Size
	buffer := make([]clw.Size, paramValueSize/clw.Size(unsafe.Sizeof(a)))
	err = clw.GetDeviceInfo(d.id, paramName, paramValueSize, unsafe.Pointer(&buffer[0]), nil)
	if err != nil {
		panic(err)
	}

	results := make([]int, len(buffer))
	for i := range results {
		results[i] = int(buffer[i])
	}

	return results
}

func (d *Device) getFpConfig(paramName clw.DeviceInfo) FPConfig {
	var paramValue clw.DeviceFPConfig
	err := clw.GetDeviceInfo(d.id, paramName, clw.Size(unsafe.Sizeof(paramValue)),
		unsafe.Pointer(&paramValue), nil)
	if err != nil {
		panic(err)
	}
	return FPConfig(paramValue)
}

func (d *Device) getExecCapabilities(paramName clw.DeviceInfo) ExecCapabilities {
	var paramValue clw.DeviceExecCapabilities
	err := clw.GetDeviceInfo(d.id, paramName, clw.Size(unsafe.Sizeof(paramValue)),
		unsafe.Pointer(&paramValue), nil)
	if err != nil {
		panic(err)
	}
	return ExecCapabilities(paramValue)
}

func (d *Device) getType(paramName clw.DeviceInfo) DeviceType {
	var paramValue clw.DeviceType
	err := clw.GetDeviceInfo(d.id, paramName, clw.Size(unsafe.Sizeof(paramValue)),
		unsafe.Pointer(&paramValue), nil)
	if err != nil {
		panic(err)
	}
	return DeviceType(paramValue)
}

func (d *Device) getCommandQueueProperties(paramName clw.DeviceInfo) CommandQueueProperties {
	var paramValue clw.CommandQueueProperties
	err := clw.GetDeviceInfo(d.id, paramName, clw.Size(unsafe.Sizeof(paramValue)),
		unsafe.Pointer(&paramValue), nil)
	if err != nil {
		panic(err)
	}
	return CommandQueueProperties(paramValue)
}

func (d *Device) getGlobalMemCacheType(paramName clw.DeviceInfo) GlobalMemCacheType {
	var paramValue clw.DeviceMemCacheType
	err := clw.GetDeviceInfo(d.id, paramName, clw.Size(unsafe.Sizeof(paramValue)),
		unsafe.Pointer(&paramValue), nil)
	if err != nil {
		panic(err)
	}
	return GlobalMemCacheType(paramValue)
}

func (d *Device) getLocalMemTypeInfo(paramName clw.DeviceInfo) LocalMemTypeInfo {
	var paramValue clw.DeviceLocalMemType
	err := clw.GetDeviceInfo(d.id, paramName, clw.Size(unsafe.Sizeof(paramValue)),
		unsafe.Pointer(&paramValue), nil)
	if err != nil {
		panic(err)
	}
	return LocalMemTypeInfo(paramValue)
}

// Check if the device supports the extension.
func (d *Device) HasExtension(extension string) bool {
	for _, v := range d.Extensions {
		if v == extension {
			return true
		}
	}
	return false
}
