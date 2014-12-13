package cl11

import (
	"encoding/binary"
	"math/rand"
	"reflect"
	"testing"
)

var kernel = `
#define int64_t long
__kernel void copy(__global float* in, __global float* out, int64_t size)
{
	for (int64_t id = get_global_id(0); id < size; id += get_global_size(0)) {
		out[id] = in[id];
	}
}
`

func TestKernel(t *testing.T) {
	allDevices := getDevices(t)
	for _, device := range allDevices {
		t.Log(device.Name, "on", device.Platform.Name)

		var toRelease []Object
		elements := int64(1024)
		size := elements * 4

		ctx, err := CreateContext([]*Device{device}, nil, contextCallback, t)
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

		host0, err := ctx.CreateHostBuffer(size*4, MemReadWrite)
		if err != nil {
			t.Error(err)
			releaseAll(toRelease, t)
			continue
		}
		toRelease = append(toRelease, host0)

		device0, err := ctx.CreateDeviceBuffer(size*4, MemReadWrite)
		if err != nil {
			t.Error(err)
			releaseAll(toRelease, t)
			continue
		}
		toRelease = append(toRelease, device0)

		map0, err := cq.EnqueueMapBuffer(host0, Blocking, MapWrite, 0, size*4, nil, nil)
		if err != nil {
			t.Error(err)
			releaseAll(toRelease, t)
			continue
		}

		values := make([]float32, size/4)
		for i := range values {
			values[i] = rand.Float32()
		}
		err = binary.Write(map0, device.ByteOrder, values)
		if err != nil {
			t.Error(err)
			releaseAll(toRelease, t)
			continue
		}

		err = cq.EnqueueUnmapBuffer(map0, nil, nil)
		if err != nil {
			t.Error(err)
			releaseAll(toRelease, t)
			continue
		}

		program, err := ctx.CreateProgramWithSource([]byte(kernel))
		if err != nil {
			t.Error(err)
			releaseAll(toRelease, t)
			continue
		}
		toRelease = append(toRelease, program)

		err = program.Build([]*Device{device}, "", nil, nil)
		if err != nil {
			t.Error(err)
			releaseAll(toRelease, t)
			continue
		}

		kernel, err := program.CreateKernel("copy")
		if err != nil {
			t.Error(err)
			releaseAll(toRelease, t)
			continue
		}
		toRelease = append(toRelease, kernel)

		err = kernel.SetArguments(host0, device0, elements)
		if err != nil {
			t.Error(err)
			releaseAll(toRelease, t)
			continue
		}

		err = cq.EnqueueNDRangeKernel(kernel, []int{0}, []int{int(elements)}, []int{1}, nil, nil)
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

		map1, err := cq.EnqueueMapBuffer(device0, Blocking, MapRead, 0, size, nil, nil)
		if err != nil {
			t.Error(err)
			releaseAll(toRelease, t)
			continue
		}

		want := values
		got := make([]float32, len(values))
		err = binary.Read(map1, device.ByteOrder, got)
		if err != nil {
			t.Error(err)
			releaseAll(toRelease, t)
			continue
		}
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

		releaseAll(toRelease, t)
	}
}
