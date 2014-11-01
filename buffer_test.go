package cl11

import (
	"fmt"
	"math/rand"
	"reflect"
	"testing"
)

func TestBuffer(t *testing.T) {
	allDevices := getDevices(t)
	for _, device := range allDevices {

		var toRelease []Object
		size := int64(1024 * 1024)

		ctx, err := CreateContext([]*Device{device}, []ContextProperties{}, nil, nil)
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

		values := map0.Float32Slice()
		for i := range values {
			values[i] = rand.Float32()
		}

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

		want := map0.Float32Slice()
		got := map1.Float32Slice()
		if !reflect.DeepEqual(want, got) {
			t.Error("values mismatch")
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

		want = make([]float32, int(size)/int(float32Size))
		got = make([]float32, int(size)/int(float32Size))
		for i := range want {
			want[i] = rand.Float32()
		}

		err = cq.EnqueueWriteBuffer(host0, NonBlocking, 0, want, nil, nil)
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

		err = cq.EnqueueReadBuffer(host1, Blocking, 0, got, nil, nil)
		if err != nil {
			t.Error(err)
			releaseAll(toRelease, t)
			continue
		}

		if !reflect.DeepEqual(want, got) {
			t.Error("values mismatch")
		}

		releaseAll(toRelease, t)
	}
}

func TestBufferRect(t *testing.T) {
	allDevices := getDevices(t)
	for _, device := range allDevices {

		var toRelease []Object
		elements := 25
		size := int64(elements * int(float32Size))

		ctx, err := CreateContext([]*Device{device}, []ContextProperties{}, nil, nil)
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

		// host1, err := ctx.CreateHostBuffer(size, MemReadWrite)
		// if err != nil {
		// 	t.Error(err)
		// 	releaseAll(toRelease, t)
		// 	continue
		// }
		// toRelease = append(toRelease, host1)

		// device0, err := ctx.CreateDeviceBuffer(size, MemReadWrite)
		// if err != nil {
		// 	t.Error(err)
		// 	releaseAll(toRelease, t)
		// 	continue
		// }
		// toRelease = append(toRelease, device0)

		want := make([]float32, elements)
		// got := make([]float32, size)
		for i := range want {
			want[i] = float32(i)
		}

		var rect Rect
		rect.Src.Origin[0] = 0
		rect.Src.Origin[1] = 0
		rect.Src.Origin[2] = 0
		rect.Src.RowPitch = 0 //5 * int64(float32Size)
		rect.Src.SlicePitch = 0
		rect.Dst.Origin[0] = 0
		rect.Dst.Origin[1] = 0
		rect.Dst.Origin[2] = 0
		rect.Dst.RowPitch = 0 //5 * int64(float32Size)
		rect.Dst.SlicePitch = 0
		rect.Region[0] = 5 * int64(float32Size)
		rect.Region[1] = 5 * int64(float32Size)
		rect.Region[2] = 1

		err = cq.EnqueueWriteBufferRect(host0, Blocking, &rect, want, nil, nil)
		if err != nil {
			t.Error(err)
			releaseAll(toRelease, t)
			continue
		}

		map0, err := cq.EnqueueMapBuffer(host0, Blocking, MapRead, 0, size, nil, nil)
		if err != nil {
			t.Error(err)
			releaseAll(toRelease, t)
			continue
		}

		fmt.Println("want:", want, "got:", map0.Float32Slice())

		err = cq.EnqueueUnmapBuffer(map0, nil, nil)
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

		// err = cq.EnqueueCopyBufferRect(host0, device0, &rect, nil, nil)
		// if err != nil {
		// 	t.Error(err)
		// 	releaseAll(toRelease, t)
		// 	continue
		// }

		// err = cq.EnqueueCopyBufferRect(device0, host1, &rect, nil, nil)
		// if err != nil {
		// 	t.Error(err)
		// 	releaseAll(toRelease, t)
		// 	continue
		// }

		// err = cq.EnqueueReadBufferRect(host1, Blocking, &rect, got, nil, nil)
		// if err != nil {
		// 	t.Error(err)
		// 	releaseAll(toRelease, t)
		// 	continue
		// }

		// for i := 1; i < 1023; i++ {
		// 	start := i*1024 + 1
		// 	end := start + 1022
		// 	if !reflect.DeepEqual(want[start:end], got[start:end]) {
		// 		t.Error("values mismatch")
		// 		break
		// 	}
		// }

		releaseAll(toRelease, t)
	}
}
