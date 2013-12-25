package cl11

import (
	"fmt"
	clw "github.com/rdwilliamson/clw11"
	"strings"
	"unsafe"
)

type Device struct {
	ID                       DeviceID
	Available                bool
	CompilerAvailable        bool
	LittleEndian             bool
	ErrorCorrectionSupport   bool
	ImageSupport             bool
	UnifiedHostMemory        bool
	AddressBits              uint32
	GlobalMemCachelineSize   uint32
	MaxClockFrequency        uint32
	MaxComputeUnits          uint32
	MaxConstantArgs          uint32
	MaxReadImageArgs         uint32
	MaxSamplers              uint32
	MaxWorkItemDimensions    uint32
	MaxWriteImageArgs        uint32
	MemBaseAddrAlign         uint32
	MinDataTypeAlignSize     uint32
	VendorID                 uint32
	PreferredVectorWidths    VectorWidths
	NativeVectorWidths       VectorWidths
	Extensions               string
	Name                     string
	Profile                  string
	Vendor                   string
	Version                  string
	OpenclCVersion           string
	DriverVersion            string
	GlobalMemCacheSize       uint64
	GlobalMemSize            uint64
	LocalMemSize             uint64
	MaxConstantBufferSize    uint64
	MaxMemAllocSize          uint64
	Image2dMaxHeight         uint
	Image2dMaxWidth          uint
	Image3dMaxDepth          uint
	Image3dMaxHeight         uint
	Image3dMaxWidth          uint
	MaxParameterSize         uint
	MaxWorkGroupSize         uint
	ProfilingTimerResolution uint
	MaxWorkItemSizes         []uint
	SingleFpConfig           FPConfig
	DoubleFpConfig           FPConfig
	ExecCapabilities         ExecCapabilities
}

type DeviceID clw.DeviceID

// Bitfield.
const (
	FPDenorm = iota
	FPFma
	FPInfNan
	FPRoundToInf
	FPRoundToNearest
	FPRoundToZero
	FPBits
)

type VectorWidths struct {
	Char   uint8
	Short  uint8
	Int    uint8
	Long   uint8
	Float  uint8
	Double uint8
	Half   uint8
}

type FPConfig uint8

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

type MemCache uint8

type LocalMem uint8

type ExecCapabilities uint8

// Bitfield.
const (
	ExecKernel = iota
	ExecNativeKernel
)

func (exec ExecCapabilities) String() string {
	var execStrings []string
	if exec&ExecKernel != 0 {
		execStrings = append(execStrings, "CL_EXEC_KERNEL")
	}
	if exec&ExecNativeKernel != 0 {
		execStrings = append(execStrings, "CL_EXEC_NATIVE_KERNEL")
	}
	return "(" + strings.Join(execStrings, "|") + ")"
}

func (p *Platform) GetDevices() ([]Device, error) {

	if p.Devices != nil {
		return p.Devices, nil
	}

	var numEntries clw.Uint
	err := clw.GetDeviceIDs(clw.PlatformID(p.ID), clw.DeviceTypeAll, 0, nil, &numEntries)
	if err != nil {
		return nil, err
	}

	deviceIDs := make([]clw.DeviceID, numEntries)
	err = clw.GetDeviceIDs(clw.PlatformID(p.ID), clw.DeviceTypeAll, numEntries, &deviceIDs[0], nil)
	if err != nil {
		return nil, err
	}

	p.Devices = make([]Device, len(deviceIDs))
	for i := range p.Devices {

		p.Devices[i].ID = DeviceID(deviceIDs[i])

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

	d.PreferredVectorWidths.Char = uint8(d.getUint(clw.DevicePreferredVectorWidthChar))
	d.PreferredVectorWidths.Short = uint8(d.getUint(clw.DevicePreferredVectorWidthShort))
	d.PreferredVectorWidths.Int = uint8(d.getUint(clw.DevicePreferredVectorWidthInt))
	d.PreferredVectorWidths.Long = uint8(d.getUint(clw.DevicePreferredVectorWidthLong))
	d.PreferredVectorWidths.Float = uint8(d.getUint(clw.DevicePreferredVectorWidthFloat))
	d.PreferredVectorWidths.Double = uint8(d.getUint(clw.DevicePreferredVectorWidthDouble))
	d.PreferredVectorWidths.Half = uint8(d.getUint(clw.DevicePreferredVectorWidthHalf))
	d.NativeVectorWidths.Char = uint8(d.getUint(clw.DeviceNativeVectorWidthChar))
	d.NativeVectorWidths.Short = uint8(d.getUint(clw.DeviceNativeVectorWidthShort))
	d.NativeVectorWidths.Int = uint8(d.getUint(clw.DeviceNativeVectorWidthInt))
	d.NativeVectorWidths.Long = uint8(d.getUint(clw.DeviceNativeVectorWidthLong))
	d.NativeVectorWidths.Float = uint8(d.getUint(clw.DeviceNativeVectorWidthFloat))
	d.NativeVectorWidths.Double = uint8(d.getUint(clw.DeviceNativeVectorWidthDouble))
	d.NativeVectorWidths.Half = uint8(d.getUint(clw.DeviceNativeVectorWidthHalf))

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

	d.ExecCapabilities = d.getExecCapabilities(clw.DeviceExecutionCapabilities)

	return
}

func (d *Device) getInfo(paramName clw.DeviceInfo) (interface{}, error) {

	switch paramName {

	// mem_cache_type
	case clw.DeviceGlobalMemCacheType:

	// device_type
	case clw.DeviceTypeInfo:

	// command_queue_properties
	case clw.DeviceQueueProperties:

	// local_mem_type
	case clw.DeviceLocalMemTypeInfo:
	}

	return nil, nil
}

func (d *Device) getBool(paramName clw.DeviceInfo) bool {
	var paramValue clw.Bool
	err := clw.GetDeviceInfo(clw.DeviceID(d.ID), paramName, clw.Size(unsafe.Sizeof(paramValue)),
		unsafe.Pointer(&paramValue), nil)
	if err != nil {
		panic(err)
	}
	return clw.ToGoBool(paramValue)
}

func (d *Device) getUint(paramName clw.DeviceInfo) uint32 {
	var paramValue clw.Uint
	err := clw.GetDeviceInfo(clw.DeviceID(d.ID), paramName, clw.Size(unsafe.Sizeof(paramValue)),
		unsafe.Pointer(&paramValue), nil)
	if err != nil {
		panic(err)
	}
	return uint32(paramValue)
}

func (d *Device) getString(paramName clw.DeviceInfo) string {
	var paramValueSize clw.Size
	err := clw.GetDeviceInfo(clw.DeviceID(d.ID), paramName, 0, nil, &paramValueSize)
	if err != nil {
		panic(err)
	}

	buffer := make([]byte, paramValueSize)
	err = clw.GetDeviceInfo(clw.DeviceID(d.ID), paramName, paramValueSize, unsafe.Pointer(&buffer[0]), nil)
	if err != nil {
		panic(err)
	}

	// Trim space and trailing \0.
	return strings.TrimSpace(string(buffer[:len(buffer)-1]))
}

func (d *Device) getUlong(paramName clw.DeviceInfo) uint64 {
	var paramValue clw.Ulong
	err := clw.GetDeviceInfo(clw.DeviceID(d.ID), paramName, clw.Size(unsafe.Sizeof(paramValue)),
		unsafe.Pointer(&paramValue), nil)
	if err != nil {
		panic(err)
	}
	return uint64(paramValue)
}

func (d *Device) getSize(paramName clw.DeviceInfo) uint {
	var paramValue clw.Size
	err := clw.GetDeviceInfo(clw.DeviceID(d.ID), paramName, clw.Size(unsafe.Sizeof(paramValue)),
		unsafe.Pointer(&paramValue), nil)
	if err != nil {
		panic(err)
	}
	return uint(paramValue)
}

func (d *Device) getSizeArray(paramName clw.DeviceInfo) []uint {
	var paramValueSize clw.Size
	err := clw.GetDeviceInfo(clw.DeviceID(d.ID), paramName, 0, nil, &paramValueSize)
	if err != nil {
		panic(err)
	}

	var a clw.Size
	buffer := make([]clw.Size, paramValueSize/clw.Size(unsafe.Sizeof(a)))
	err = clw.GetDeviceInfo(clw.DeviceID(d.ID), paramName, paramValueSize, unsafe.Pointer(&buffer[0]), nil)
	if err != nil {
		panic(err)
	}

	results := make([]uint, len(buffer))
	for i := range results {
		results[i] = uint(buffer[i])
	}

	return results
}

func (d *Device) getFpConfig(paramName clw.DeviceInfo) FPConfig {
	var paramValue clw.DeviceFPConfig
	err := clw.GetDeviceInfo(clw.DeviceID(d.ID), paramName, clw.Size(unsafe.Sizeof(paramValue)),
		unsafe.Pointer(&paramValue), nil)
	if err != nil {
		panic(err)
	}

	var result FPConfig
	if paramValue&clw.FPDenorm != 0 {
		result |= FPDenorm
	}
	if paramValue&clw.FPFma != 0 {
		result |= FPFma
	}
	if paramValue&clw.FPInfNan != 0 {
		result |= FPInfNan
	}
	if paramValue&clw.FPRoundToInf != 0 {
		result |= FPRoundToInf
	}
	if paramValue&clw.FPRoundToNearest != 0 {
		result |= FPRoundToNearest
	}
	if paramValue&clw.FPRoundToZero != 0 {
		result |= FPRoundToZero
	}
	return result
}

func (d *Device) getExecCapabilities(paramName clw.DeviceInfo) ExecCapabilities {
	var paramValue clw.DeviceExecCapabilities
	err := clw.GetDeviceInfo(clw.DeviceID(d.ID), paramName, clw.Size(unsafe.Sizeof(paramValue)),
		unsafe.Pointer(&paramValue), nil)
	if err != nil {
		panic(err)
	}

	var result ExecCapabilities
	if paramValue&clw.ExecKernel != 0 {
		result |= ExecKernel
	}
	if paramValue&clw.ExecNativeKernel != 0 {
		result |= ExecNativeKernel
	}
	return result
}
