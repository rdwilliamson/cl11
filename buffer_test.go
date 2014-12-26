package cl11

import (
	"math/rand"
	"reflect"
	"testing"
	"time"
)

func TestBuffer(t *testing.T) {
	rand.Seed(time.Now().Unix())
	allDevices := getDevices(t)
	for _, device := range allDevices {
		t.Log(device.Name, "on", device.Platform.Name)

		var toRelease []Object
		size := int64(1024 * 1024)

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

		host0, err := ctx.CreateHostBuffer(size, MemReadWrite)
		if err != nil {
			t.Error(err)
			releaseAll(toRelease, t)
			continue
		}
		toRelease = append(toRelease, host0)

		host1, err := ctx.CreateHostBuffer(size, MemReadWrite)
		if err != nil {
			t.Error(err)
			releaseAll(toRelease, t)
			continue
		}
		toRelease = append(toRelease, host1)

		device0, err := ctx.CreateDeviceBuffer(size, MemReadWrite)
		if err != nil {
			t.Error(err)
			releaseAll(toRelease, t)
			continue
		}
		toRelease = append(toRelease, device0)

		map0, err := cq.EnqueueMapBuffer(host0, Blocking, MapWrite, 0, size, nil, nil)
		if err != nil {
			t.Error(err)
			releaseAll(toRelease, t)
			continue
		}

		values := make([]float32, size/4)
		for i := range values {
			values[i] = rand.Float32()
		}
		copy(map0.Float32s(), values)

		err = cq.EnqueueUnmapBuffer(map0, nil, nil)
		if err != nil {
			t.Error(err)
			releaseAll(toRelease, t)
			continue
		}

		err = cq.EnqueueCopyBuffer(host0, device0, 0, 0, size, nil, nil)
		if err != nil {
			t.Error(err)
			releaseAll(toRelease, t)
			continue
		}

		err = cq.EnqueueCopyBuffer(device0, host1, 0, 0, size, nil, nil)
		if err != nil {
			t.Error(err)
			releaseAll(toRelease, t)
			continue
		}

		map0, err = cq.EnqueueMapBuffer(host0, Blocking, MapRead, 0, size, nil, nil)
		if err != nil {
			t.Error(err)
			releaseAll(toRelease, t)
			continue
		}

		map1, err := cq.EnqueueMapBuffer(host1, Blocking, MapRead, 0, size, nil, nil)
		if err != nil {
			t.Error(err)
			releaseAll(toRelease, t)
			continue
		}

		if !reflect.DeepEqual(map1.Float32s(), values) {
			t.Error("values mismatch")
			releaseAll(toRelease, t)
			continue
		}

		err = cq.EnqueueUnmapBuffer(map0, nil, nil)
		if err != nil {
			t.Error(err)
			releaseAll(toRelease, t)
			continue
		}

		err = cq.EnqueueUnmapBuffer(map1, nil, nil)
		if err != nil {
			t.Error(err)
			releaseAll(toRelease, t)
			continue
		}

		err = cq.Finish()
		if err != nil {
			t.Error(err)
			releaseAll(toRelease, t)
			continue
		}

		releaseAll(toRelease, t)
	}
}
