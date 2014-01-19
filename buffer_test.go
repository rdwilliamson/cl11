package cl11

import "testing"

func TestCreateBuffers(t *testing.T) {
	var host []byte

	contexts := createContexts(t)
	for _, context := range contexts {

		if cap(host) < int(context.Devices[0].MaxMemAllocSize) {
			host = make([]byte, int(context.Devices[0].MaxMemAllocSize))
		}
		host = host[:int(context.Devices[0].MaxMemAllocSize)]

		buffer, err := context.CreateDeviceBuffer(len(host), MemoryFlags{})
		if err != nil {
			t.Error(err)
			continue
		}
		err = buffer.Release()
		if err != nil {
			t.Error(err)
		}

		buffer, err = context.CreateDeviceBufferFromHost(MemoryFlags{}, host)
		if err != nil {
			t.Error(err)
			continue
		}
		err = buffer.Release()
		if err != nil {
			t.Error(err)
		}

		buffer, err = context.CreateDeviceBufferOnHost(MemoryFlags{}, host)
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

		buffer, err = context.CreateHostBuffer(len(host), MemoryFlags{})
		if err != nil {
			t.Error(err)
			continue
		}
		err = buffer.Release()
		if err != nil {
			t.Error(err)
		}

		buffer, err = context.CreateHostBufferFromHost(MemoryFlags{}, host)
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
