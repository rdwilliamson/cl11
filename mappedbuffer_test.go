package cl11

import (
	"bytes"
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

		// Reset relative to start.
		nn, err := map0.Seek(0, 0)
		if nn != 0 || err != nil {
			t.Error("seeking: want 0 <nil>, got", nn, err)
			releaseAll(toRelease, t)
			continue
		}

		// Write to, the full buffer.
		var buf bytes.Buffer
		nn, err = map0.WriteTo(&buf)
		if nn != int64(size) || err != nil {
			t.Error("writing to: want", size, "<nil>, got", nn, err)
			releaseAll(toRelease, t)
			continue
		}

		// Write to, but buffers already consumed.
		nn, err = map0.WriteTo(&buf)
		if nn != 0 || err != nil {
			t.Error("writing to: want 0 <nil>, got", nn, err)
			releaseAll(toRelease, t)
			continue
		}

		// Read at everything.
		n, err = map0.ReadAt(everything, 0)
		if n != size || err != nil {
			t.Error("read at 0 everything: want", size, "<nil>, got", n, err)
			releaseAll(toRelease, t)
			continue
		}

		// Read at running past end of buffer.
		n, err = map0.ReadAt(everything, 1)
		if n != size-1 || err != io.EOF {
			t.Error("read at 0 everything: want", size-1, "<nil>, got", n, err)
			releaseAll(toRelease, t)
			continue
		}

		// Read at starting past the end of the buffer.
		n, err = map0.ReadAt(everything, int64(size))
		if n != 0 || err != io.EOF {
			t.Error("read at 0 everything: want 0 io.EOF, got", n, err)
			releaseAll(toRelease, t)
			continue
		}

		// Read at starting at a negative value.
		n, err = map0.ReadAt(everything, -1)
		if n != 0 || err == nil {
			t.Error("read at -1 everything: want 0 cl: MappedBuffer.ReadAt: negative offset, got", n, err)
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

func TestMappedBufferWrite(t *testing.T) {
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

		map0, err := cq.EnqueueMapBuffer(device0, Blocking, MapWrite, 0, int64(size), nil, nil)
		if err != nil {
			t.Error(err)
			releaseAll(toRelease, t)
			continue
		}

		// Write nothing.
		var nothing []byte
		n, err := map0.Write(nothing)
		if n != 0 || err != nil {
			t.Error("writing nothing: want 0 <nil>, got", n, err)
			releaseAll(toRelease, t)
			continue
		}

		// Write everything.
		everything := make([]byte, size)
		n, err = map0.Write(everything)
		if n != size || err != nil {
			t.Error("writting everything: want", size, "<nil>, got", n, err)
			releaseAll(toRelease, t)
			continue
		}

		// Write past end of buffer.
		n, err = map0.Write(everything)
		if n != 0 || err != ErrBufferFull {
			t.Error("writting past end of buffer: want 0", ErrBufferFull, "got", n, err)
			releaseAll(toRelease, t)
			continue
		}

		// Reset relative to start.
		nn, err := map0.Seek(0, 0)
		if nn != 0 || err != nil {
			t.Error("seeking: want 0 <nil>, got", nn, err)
			releaseAll(toRelease, t)
			continue
		}

		// Write at everything.
		n, err = map0.WriteAt(everything, 0)
		if n != size || err != nil {
			t.Error("write at 0 everything: want", size, "<nil>, got", n, err)
			releaseAll(toRelease, t)
			continue
		}

		// Write at running past the end of the buffer.
		n, err = map0.WriteAt(everything, 1)
		if n != size-1 || err != ErrBufferFull {
			t.Error("write at 1 everything: want", size-1, "<nil>, got", n, err)
			releaseAll(toRelease, t)
			continue
		}

		// Write at starting past the end of the buffer.
		n, err = map0.WriteAt(everything, int64(size))
		if n != 0 || err != ErrBufferFull {
			t.Error("write at end of buffer: want 0 <nil>, got", n, err)
			releaseAll(toRelease, t)
			continue
		}

		// Write at starting at a negative value.
		n, err = map0.WriteAt(everything, -1)
		if n != 0 || err == nil {
			t.Error("write at -1 everything: want 0 cl: MappedBuffer.WriteAt: negative offset, got", n, err)
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
