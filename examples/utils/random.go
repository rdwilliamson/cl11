package utils

import (
	"math/rand"
)

func RandomFloat32(count int) []float32 {
	r := make([]float32, count)
	for i := range r {
		r[i] = rand.Float32()
	}
	return r
}
