package cl11

import (
	"io"
	"testing"
)

func TestMappedBufferRead(t *testing.T) {
	allDevices := getDevices(t)
	for _, device := range allDevices {
		t.Log(device.Name, "on", device.Platform.Name)

		var toRelease []Object
		size := 1024 * 1024

		ctx, err := CreateContext([]*Device{device}, nil, nil, nil)
		if err != nil {
			t.Error(err)
			continue
		}
		toRelease = append(toRelease, ctx)

		cq, err := ctx.CreateCommandQueue(device, 0)
		if err != nil {
			t.Error(err)
			releaseAll(toRelease, t)
			continue
		}
		toRelease = append(toRelease, cq)

		device0, err := ctx.CreateDeviceBuffer(int64(size), MemReadWrite)
		if err != nil {
			t.Error(err)
			releaseAll(toRelease, t)
			continue
		}
		toRelease = append(toRelease, device0)

		map0, err := cq.EnqueueMapBuffer(device0, Blocking, MapRead, 0, int64(size), nil, nil)
		if err != nil {
			t.Error(err)
			releaseAll(toRelease, t)
			continue
		}

		// Read nothing.
		var nothing []byte
		n, err := map0.Read(nothing)
		if n != 0 || err != nil {
			t.Error("reading nothing: want 0 <nil>, got", n, err)
			releaseAll(toRelease, t)
			continue
		}

		// Read everything.
		everything := make([]byte, size)
		n, err = map0.Read(everything)
		if n != size || err != nil {
			t.Error("reading everything: want", size, "<nil>, got", n, err)
			releaseAll(toRelease, t)
			continue
		}

		// Read past end.
		n, err = map0.Read(everything)
		if n != 0 || err != io.EOF {
			t.Error("reading past end of buffer: want 0 io.EOF, got", n, err)
			releaseAll(toRelease, t)
			continue
		}

		err = cq.EnqueueUnmapBuffer(map0, nil, nil)
		if err != nil {
			t.Error("failed to unmap buffer")
			releaseAll(toRelease, t)
			continue
		}

		releaseAll(toRelease, t)
	}
}
