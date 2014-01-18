package cl11

import (
	"testing"
)

func TestCreateUserEvent(t *testing.T) {
	contexts := createContexts(t)
	for _, context := range contexts {
		_, err := CreateUserEvent(context)
		if err != nil {
			t.Error(err)
		}
	}
}
