package cl11

import "testing"

func TestGetDevices(t *testing.T) {
	platforms, err := GetPlatforms()
	if err != nil {
		t.Fatal(err)
	}
	for i := range platforms {
		_, err = platforms[i].GetDevices()
		if err != nil {
			t.Fatal(err)
		}
	}
}
