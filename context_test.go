package cl11

import "testing"

func contextCallback(err string, data []byte, userData interface{}) {
	t := userData.(*testing.T)
	t.Log("Error:", err)
}
