package cl11

import (
	clw "github.com/rdwilliamson/clw11"
	"strings"
	"unsafe"
)

type Device struct {
	ID                     clw.DeviceID
	Available              bool
	CompilerAvailable      bool
	LittleEndian           bool
	ErrorCorrectionSupport bool
	ImageSupport           bool
	UnifiedHostMemory      bool
}

type VectorSizes struct {
	Char   uint32
	Short  uint32
	Int    uint32
	Long   uint32
	Float  uint32
	Double uint32
	Half   uint32
}

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

func (d Device) String() string {
	return ""
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

	return
}

func (d *Device) getInfo(paramName clw.DeviceInfo) (interface{}, error) {

	switch paramName {

	// bool
	case clw.DeviceAvailable,
		clw.DeviceCompilerAvailable,
		clw.DeviceEndianLittle,
		clw.DeviceErrorCorrectionSupport,
		clw.DeviceImageSupport,
		clw.DeviceHostUnifiedMemory:

	// uint
	case clw.DeviceAddressBits,
		clw.DeviceGlobalMemCachelineSize,
		clw.DeviceMaxClockFrequency,
		clw.DeviceMaxComputeUnits,
		clw.DeviceMaxConstantArgs,
		clw.DeviceMaxReadImageArgs,
		clw.DeviceMaxSamplers,
		clw.DeviceMaxWorkItemDimensions,
		clw.DeviceMaxWriteImageArgs,
		clw.DeviceMemBaseAddrAlign,
		clw.DeviceMinDataTypeAlignSize,
		clw.DevicePreferredVectorWidthChar,
		clw.DevicePreferredVectorWidthShort,
		clw.DevicePreferredVectorWidthInt,
		clw.DevicePreferredVectorWidthLong,
		clw.DevicePreferredVectorWidthFloat,
		clw.DevicePreferredVectorWidthDouble,
		clw.DevicePreferredVectorWidthHalf,
		clw.DeviceNativeVectorWidthChar,
		clw.DeviceNativeVectorWidthShort,
		clw.DeviceNativeVectorWidthInt,
		clw.DeviceNativeVectorWidthLong,
		clw.DeviceNativeVectorWidthFloat,
		clw.DeviceNativeVectorWidthDouble,
		clw.DeviceNativeVectorWidthHalf,
		clw.DeviceVendorId:

	// fp_config
	case clw.DeviceSingleFpConfig:

	// exec_capabilities
	case clw.DeviceExecutionCapabilities:

	// char[]
	case clw.DeviceExtensions,
		clw.DeviceName,
		clw.DeviceProfile,
		clw.DeviceVendor,
		clw.DeviceVersion,
		clw.DriverVersion,
		clw.DeviceOpenclCVersion:

	// ulong
	case clw.DeviceGlobalMemCacheSize,
		clw.DeviceGlobalMemSize,
		clw.DeviceLocalMemSize,
		clw.DeviceMaxConstantBufferSize,
		clw.DeviceMaxMemAllocSize:

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
		panic("not implemented yet")

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
	err = clw.GetDeviceInfo(d.ID, paramName, paramValueSize, clw.Pointer(buffer), nil)
	if err != nil {
		panic(err)
	}

	// Trim space and trailing \0.
	return strings.TrimSpace(string(buffer[:len(buffer)-1]))
}
