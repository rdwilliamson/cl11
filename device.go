package cl11

import (
	"fmt"
	clw "github.com/rdwilliamson/clw11"
	"strings"
	"unsafe"
)

type Device struct {
	ID                       clw.DeviceID
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

	MaxWorkItemSizes []uint
}

// Bitfield.
const (
	DeviceTypeDefault     = clw.DeviceTypeDefault
	DeviceTypeCpu         = clw.DeviceTypeCpu
	DeviceTypeGpu         = clw.DeviceTypeGpu
	DeviceTypeAccelerator = clw.DeviceTypeAccelerator
	DeviceTypeAll         = clw.DeviceTypeAll
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

type MemCache uint8

type LocalMem uint8

type ExecCapabilities uint8

func (p *Platform) GetDevices() ([]Device, error) {

	if p.Devices != nil {
		return p.Devices, nil
	}

	var numEntries clw.Uint
	err := clw.GetDeviceIDs(p.ID, clw.DeviceTypeAll, 0, nil, &numEntries)
	if err != nil {
		return nil, err
	}

	deviceIDs := make([]clw.DeviceID, numEntries)
	err = clw.GetDeviceIDs(p.ID, clw.DeviceTypeAll, numEntries, &deviceIDs[0], nil)
	if err != nil {
		return nil, err
	}

	p.Devices = make([]Device, len(deviceIDs))
	for i := range p.Devices {

		p.Devices[i].ID = deviceIDs[i]

		err = p.Devices[i].getAllInfo()
		if err != nil {
			return nil, err
		}
	}

	return p.Devices, nil
}

func (d *Device) getAllInfo() (err error) {
	// defer func() {
	// 	if r := recover(); r != nil {
	// 		err = r.(error)
	// 	}
	// }()

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

	return
}

func (d *Device) getInfo(paramName clw.DeviceInfo) (interface{}, error) {

	switch paramName {

	// fp_config
	case clw.DeviceSingleFpConfig:

	// exec_capabilities
	case clw.DeviceExecutionCapabilities:

	// mem_cache_type
	case clw.DeviceGlobalMemCacheType:

	// size_t
	case clw.DeviceImage2dMaxHeight,
		clw.DeviceImage2dMaxWidth,
		clw.DeviceImage3dMaxDepth,
		clw.DeviceImage3dMaxHeight,
		clw.DeviceImage3dMaxWidth,
		clw.DeviceMaxParameterSize,
		clw.DeviceMaxWorkGroupSize,
		clw.DeviceMaxWorkItemSizes,
		clw.DeviceProfilingTimerResolution:

	// device_type
	case clw.DeviceTypeInfo:

	// command_queue_properties
	case clw.DeviceQueueProperties:

	// local_mem_type
	case clw.DeviceLocalMemTypeInfo:

	// platform_id
	case clw.DevicePlatform:
	}

	return nil, nil
}

func (d *Device) getBool(paramName clw.DeviceInfo) bool {
	var paramValue clw.Bool
	err := clw.GetDeviceInfo(d.ID, paramName, clw.Size(unsafe.Sizeof(paramValue)), unsafe.Pointer(&paramValue), nil)
	if err != nil {
		panic(err)
	}
	return clw.ToGoBool(paramValue)
}

func (d *Device) getUint(paramName clw.DeviceInfo) uint32 {
	var paramValue clw.Uint
	err := clw.GetDeviceInfo(d.ID, paramName, clw.Size(unsafe.Sizeof(paramValue)), unsafe.Pointer(&paramValue), nil)
	if err != nil {
		panic(err)
	}
	return uint32(paramValue)
}

func (d *Device) getString(paramName clw.DeviceInfo) string {
	var paramValueSize clw.Size
	err := clw.GetDeviceInfo(d.ID, paramName, 0, nil, &paramValueSize)
	if err != nil {
		panic(err)
	}

	buffer := make([]byte, paramValueSize)
	err = clw.GetDeviceInfo(d.ID, paramName, paramValueSize, unsafe.Pointer(&buffer[0]), nil)
	if err != nil {
		panic(err)
	}

	// Trim space and trailing \0.
	return strings.TrimSpace(string(buffer[:len(buffer)-1]))
}

func (d *Device) getUlong(paramName clw.DeviceInfo) uint64 {
	var paramValue clw.Ulong
	err := clw.GetDeviceInfo(d.ID, paramName, clw.Size(unsafe.Sizeof(paramValue)), unsafe.Pointer(&paramValue), nil)
	if err != nil {
		panic(err)
	}
	return uint64(paramValue)
}

func (d *Device) getSize(paramName clw.DeviceInfo) uint {
	var paramValue clw.Size
	err := clw.GetDeviceInfo(d.ID, paramName, clw.Size(unsafe.Sizeof(paramValue)), unsafe.Pointer(&paramValue), nil)
	if err != nil {
		panic(err)
	}
	return uint(paramValue)
}

func (d *Device) getSizeArray(paramName clw.DeviceInfo) []uint {
	var paramValueSize clw.Size
	err := clw.GetDeviceInfo(d.ID, paramName, 0, nil, &paramValueSize)
	if err != nil {
		panic(err)
	}

	var a clw.Size
	buffer := make([]clw.Size, paramValueSize/clw.Size(unsafe.Sizeof(a)))
	fmt.Println(paramValueSize, len(buffer))
	err = clw.GetDeviceInfo(d.ID, paramName, paramValueSize, unsafe.Pointer(&buffer[0]), nil)
	if err != nil {
		panic(err)
	}

	results := make([]uint, len(buffer))
	for i := range results {
		results[i] = uint(buffer[i])
	}

	return results
}
