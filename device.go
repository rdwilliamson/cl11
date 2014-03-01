package cl11

import (
	"strings"
	"unsafe"

	clw "github.com/rdwilliamson/clw11"
)

type Device struct {
	id                       clw.DeviceID
	Type                     DeviceType
	Available                bool
	CompilerAvailable        bool
	LittleEndian             bool
	ErrorCorrectionSupport   bool
	ImageSupport             bool
	UnifiedHostMemory        bool
	AddressBits              int
	GlobalMemCachelineSize   int
	MaxClockFrequency        int
	MaxComputeUnits          int
	MaxConstantArgs          int
	MaxReadImageArgs         int
	MaxSamplers              int
	MaxWorkItemDimensions    int
	MaxWriteImageArgs        int
	MemBaseAddrAlign         int
	MinDataTypeAlignSize     int
	VendorID                 int
	PreferredVectorWidths    VectorWidths
	NativeVectorWidths       VectorWidths
	Extensions               string
	Name                     string
	Profile                  string
	Vendor                   string
	Version                  string
	OpenclCVersion           string
	DriverVersion            string
	GlobalMemCacheSize       int64
	GlobalMemSize            int64
	LocalMemSize             int64
	MaxConstantBufferSize    int64
	MaxMemAllocSize          int64
	Image2dMaxHeight         int
	Image2dMaxWidth          int
	Image3dMaxDepth          int
	Image3dMaxHeight         int
	Image3dMaxWidth          int
	MaxParameterSize         int
	MaxWorkGroupSize         int
	ProfilingTimerResolution int
	MaxWorkItemSizes         []int
	SingleFpConfig           FPConfig
	DoubleFpConfig           FPConfig
	// HalfFpConfig             FPConfig
	ExecCapabilities       ExecCapabilities
	CommandQueueProperties CommandQueueProperties
	GlobalMemCacheType     GlobalMemCacheType
	LocalMemTypeInfo       LocalMemTypeInfo
}

type DeviceType uint32

// Bitfield.
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
	panic("unknown device type")
}

type VectorWidths struct {
	Char   int
	Short  int
	Int    int
	Long   int
	Float  int
	Double int
	Half   int
}

type FPConfig uint32

// Bitfield.
const (
	FPDenorm         = FPConfig(clw.FPDenorm)
	FPFma            = FPConfig(clw.FPFma)
	FPInfNan         = FPConfig(clw.FPInfNan)
	FPRoundToInf     = FPConfig(clw.FPRoundToInf)
	FPRoundToNearest = FPConfig(clw.FPRoundToNearest)
	FPRoundToZero    = FPConfig(clw.FPRoundToZero)
)

func (fpConfig FPConfig) String() string {
	var configStrings []string
	if fpConfig&FPDenorm != 0 {
		configStrings = append(configStrings, "CL_FP_DENORM")
	}
	if fpConfig&FPFma != 0 {
		configStrings = append(configStrings, "CL_FP_FMA")
	}
	if fpConfig&FPInfNan != 0 {
		configStrings = append(configStrings, "CL_FP_INF_NAN")
	}
	if fpConfig&FPRoundToInf != 0 {
		configStrings = append(configStrings, "CL_FP_ROUND_TO_INF")
	}
	if fpConfig&FPRoundToNearest != 0 {
		configStrings = append(configStrings, "CL_FP_ROUND_TO_NEAREST")
	}
	if fpConfig&FPRoundToZero != 0 {
		configStrings = append(configStrings, "CL_FP_ROUND_TO_ZERO")
	}
	return "(" + strings.Join(configStrings, "|") + ")"
}

type GlobalMemCacheType uint32

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
		return "read write"
	}
	panic("unknown global mem cache type")
}

type LocalMemTypeInfo uint32

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
	panic("unknown local mem type")
}

type ExecCapabilities uint32

// Bitfield.
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
	return "(" + strings.Join(execStrings, "|") + ")"
}

func (p *Platform) GetDevices() ([]*Device, error) {

	if p.Devices != nil {
		return p.Devices, nil
	}

	var numEntries clw.Uint
	err := clw.GetDeviceIDs(p.id, clw.DeviceTypeAll, 0, nil, &numEntries)
	if err != nil {
		return nil, err
	}

	deviceIDs := make([]clw.DeviceID, numEntries)
	err = clw.GetDeviceIDs(p.id, clw.DeviceTypeAll, numEntries, &deviceIDs[0], nil)
	if err != nil {
		return nil, err
	}

	p.Devices = make([]*Device, len(deviceIDs))
	for i := range p.Devices {

		p.Devices[i] = &Device{id: deviceIDs[i]}

		err = p.Devices[i].getAllInfo()
		if err != nil {
			return nil, err
		}
	}

	return p.Devices, nil
}

func (d *Device) getAllInfo() (err error) {

	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()

	d.Available = d.getBool(clw.DeviceAvailable)
	d.CompilerAvailable = d.getBool(clw.DeviceCompilerAvailable)
	d.LittleEndian = d.getBool(clw.DeviceEndianLittle)
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
	d.PreferredVectorWidths.Half = d.getUint(clw.DevicePreferredVectorWidthHalf)
	d.NativeVectorWidths.Char = d.getUint(clw.DeviceNativeVectorWidthChar)
	d.NativeVectorWidths.Short = d.getUint(clw.DeviceNativeVectorWidthShort)
	d.NativeVectorWidths.Int = d.getUint(clw.DeviceNativeVectorWidthInt)
	d.NativeVectorWidths.Long = d.getUint(clw.DeviceNativeVectorWidthLong)
	d.NativeVectorWidths.Float = d.getUint(clw.DeviceNativeVectorWidthFloat)
	d.NativeVectorWidths.Double = d.getUint(clw.DeviceNativeVectorWidthDouble)
	d.NativeVectorWidths.Half = d.getUint(clw.DeviceNativeVectorWidthHalf)

	d.Extensions = d.getString(clw.DeviceExtensions)
	d.Name = d.getString(clw.DeviceName)
	d.Profile = d.getString(clw.DeviceProfile)
	d.Vendor = d.getString(clw.DeviceVendor)
	d.Version = d.getString(clw.DeviceVersion)
	d.OpenclCVersion = d.getString(clw.DeviceOpenclCVersion)
	d.DriverVersion = d.getString(clw.DriverVersion)

	d.GlobalMemCacheSize = d.getUlong(clw.DeviceGlobalMemCacheSize)
	d.GlobalMemSize = d.getUlong(clw.DeviceGlobalMemSize)
	d.LocalMemSize = d.getUlong(clw.DeviceLocalMemSize)
	d.MaxConstantBufferSize = d.getUlong(clw.DeviceMaxConstantBufferSize)
	d.MaxMemAllocSize = d.getUlong(clw.DeviceMaxMemAllocSize)

	d.Image2dMaxHeight = d.getSize(clw.DeviceImage2dMaxHeight)
	d.Image2dMaxWidth = d.getSize(clw.DeviceImage2dMaxWidth)
	d.Image3dMaxDepth = d.getSize(clw.DeviceImage3dMaxDepth)
	d.Image3dMaxHeight = d.getSize(clw.DeviceImage3dMaxHeight)
	d.Image3dMaxWidth = d.getSize(clw.DeviceImage3dMaxWidth)
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
