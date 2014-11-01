package cl11

import (
	"fmt"
	"strings"
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
		if strings.Contains(platform.Name, "AMD") {
			fmt.Println("skipping amd")
			continue
		}
		if strings.Contains(platform.Name, "Intel") {
			fmt.Println("skipping intel")
			continue
		}
		// if strings.Contains(platform.Name, "NVIDIA") {
		// 	fmt.Println("skipping nvidia")
		// 	continue
		// }
		for _, device := range platform.Devices {
			results = append(results, device)
		}
	}
	return results
}
