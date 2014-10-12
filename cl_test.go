package cl11

import (
	"testing"
)

func releaseAll(o []Object, t *testing.T) {
	for _, v := range o {
		err := v.Release()
		if err != nil {
			t.Error(err)
		}
	}
}

func getDevices(t *testing.T) []*Device {
	allPlatforms, err := GetPlatforms()
	if err != nil {
		t.Fatal(err)
	}
	var results []*Device
	for _, platform := range allPlatforms {
		for _, device := range platform.Devices {
			results = append(results, device)
		}
	}
	return results
}
