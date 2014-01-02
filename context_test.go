package cl11

import "testing"

func TestCreateContext(t *testing.T) {
	platforms, err := GetPlatforms()
	if err != nil {
		t.Fatal(err)
	}
	for i := range platforms {
		devices, err := platforms[i].GetDevices()
		if err != nil {
			t.Fatal(err)
		}
		for j := range devices {
			_, err = CreateContext([]Device{devices[j]})
			if err != nil {
				t.Error(err)
			}
		}
	}
}
