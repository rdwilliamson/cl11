package cl11

import "testing"

func TestCreateDeviceBuffer(t *testing.T) {
	var host []byte

	contexts := createContexts(t)
	for _, context := range contexts {

		if cap(host) < int(context.Devices[0].MaxMemAllocSize) {
			host = make([]byte, int(context.Devices[0].MaxMemAllocSize))
		}
		host = host[:int(context.Devices[0].MaxMemAllocSize)]

		buffer, err := CreateDeviceBuffer(context, len(host), true, true, nil, false)
		if err != nil {
			t.Error(err)
			continue
		}
		err = buffer.Release()
		if err != nil {
			t.Error(err)
		}

		if context.Devices[0].Type == DeviceTypeCpu {
			continue
		}

		buffer, err = CreateDeviceBuffer(context, len(host), true, true, host, false)
		if err != nil {
			t.Error(err)
			continue
		}
		err = buffer.Release()
		if err != nil {
			t.Error(err)
		}

		buffer, err = CreateDeviceBuffer(context, len(host), true, true, host, true)
		if err != nil {
			t.Error(err)
			continue
		}
		err = buffer.Release()
		if err != nil {
			t.Error(err)
		}
	}
}

func TestCreateHostBuffer(t *testing.T) {
	contexts := createContexts(t)
	for _, context := range contexts {
		if context.Devices[0].Type == DeviceTypeCpu {
			continue
		}

		buffer, err := CreateHostBuffer(context, int(context.Devices[0].MaxMemAllocSize), true, true)
		if err != nil {
			t.Error(err)
			continue
		}
		err = buffer.Release()
		if err != nil {
			t.Error(err)
		}
	}
}
