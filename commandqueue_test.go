package cl11

import (
	"testing"
)

func TestCreateCommandQueue(t *testing.T) {
	contexts := createContexts(t)
	for _, context := range contexts {
		_, err := CreateCommandQueue(context, context.Devices[0], CommandQueueProperties{})
		if err != nil {
			t.Error(err)
			continue
		}
	}
}
