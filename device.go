package cl11

import (
	clw "github.com/rdwilliamson/clw11"
	"unsafe"
)

type Device struct {
	id clw.DeviceID
}

func (p *Platform) GetDevices() ([]Device, error) {

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

	p.Devices = make([]Device, len(deviceIDs))
	for i := range p.Devices {
		p.Devices[i].id = deviceIDs[i]
	}

	return p.Devices, nil
}

func (d *Device) GetInfo() error {
	return nil
}

func (d Device) String() string {
	return ""
}

func getDeviceUint(id clw.DeviceID, paramName clw.DeviceInfo) (clw.Uint, error) {
	var paramValue clw.Uint
	err := clw.GetDeviceInfo(id, paramName, clw.Size(unsafe.Sizeof(paramValue)), unsafe.Pointer(&paramValue), nil)
	if err != nil {
		return 0, err
	}
	return paramValue, nil
}

func (d *Device) getInfo(paramName clw.DeviceInfo) error {

	switch paramName {

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

	// bool
	case clw.DeviceAvailable,
		clw.DeviceCompilerAvailable,
		clw.DeviceEndianLittle,
		clw.DeviceErrorCorrectionSupport,
		clw.DeviceImageSupport,
		clw.DeviceHostUnifiedMemory:

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

	return nil
}
