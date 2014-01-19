package cl11

import (
	"testing"
)

func TestCreateCommandQueue(t *testing.T) {
	contexts := createContexts(t)
	for _, context := range contexts {
		_, err := context.CreateCommandQueue(context.Devices[0], CommandQueueProperties{})
		if err != nil {
			t.Error(err)
			continue
		}
	}
}
