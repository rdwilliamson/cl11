package cl11

import (
	"reflect"
	"testing"

	"math/rand"
)

type bufferTestData struct {
	cq         *CommandQueue
	in         *Buffer
	out        *Buffer
	k          *Kernel
	localSize  []int
	globalSize []int
}

func setupBuffers(d *Device, t *testing.T) bufferTestData {
	c, err := CreateContext([]*Device{d}, nil, nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	progam, err := c.CreateProgramWithSource([]byte(`
#define int64_t long
__kernel void copy(__global float* in, __global float* out, int64_t size)
{
	for (int64_t id = get_global_id(0); id < size; id += get_global_size(0)) {
		out[id] = in[id];
	}
}`))
	if err != nil {
		t.Fatal(err)
	}

	err = progam.Build([]*Device{d}, "", nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	k, err := progam.CreateKernel("copy")
	if err != nil {
		t.Fatal(err)
	}

	size := int64(1024 * 1024 / 4)
	if size*4 > d.MaxMemAllocSize {
		size = d.MaxMemAllocSize / 4
	}

	in, err := c.CreateDeviceBuffer(size*4, MemReadOnly)
	if err != nil {
		t.Fatal(err)
	}
	out, err := c.CreateDeviceBuffer(size*4, MemWriteOnly)
	if err != nil {
		t.Fatal(err)
	}

	err = k.SetArguments(in, out, size)
	if err != nil {
		t.Fatal(err)
	}

	cq, err := c.CreateCommandQueue(d, 0)
	if err != nil {
		t.Fatal(err)
	}

	localSize := k.WorkGroupInfo[0].PreferredWorkGroupSizeMultiple
	globalSize := int(size)
	if globalSize%localSize > 0 {
		globalSize = (globalSize/localSize + 1) * localSize
	}

	return bufferTestData{cq, in, out, k, []int{localSize}, []int{globalSize}}
}

func TestEnqueueReadWriteBuffer(t *testing.T) {
	platforms, err := GetPlatforms()
	if err != nil {
		t.Fatal(err)
	}
	for _, p := range platforms {
		for _, d := range p.Devices {

			vars := setupBuffers(d, t)

			want := make([]float32, vars.in.Size/4)
			got := make([]float32, vars.in.Size/4)
			for i := range want {
				want[i] = rand.Float32()
			}

			err = vars.cq.EnqueueWriteBuffer(vars.in, NonBlocking, 0, want, nil, nil)
			if err != nil {
				t.Fatal(err)
			}
			err = vars.cq.EnqueueNDRangeKernel(vars.k, nil, vars.globalSize, vars.localSize, nil, nil)
			if err != nil {
				t.Fatal(err)
			}
			err = vars.cq.EnqueueReadBuffer(vars.out, Blocking, 0, got, nil, nil)
			if err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(want, got) {
				t.Error("copy does not match")
			}
		}
	}
}
