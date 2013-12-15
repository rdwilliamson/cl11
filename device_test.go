package cl11

import (
	"fmt"
	"testing"
)

func TestGetDevices(t *testing.T) {
	platforms, err := GetPlatforms()
	if err != nil {
		t.Fatal(err)
	}
	for i := range platforms {
		fmt.Println(platforms[i].Name)
		devices, err := platforms[i].GetDevices()
		if err != nil {
			t.Fatal(err)
		}
		for j := range devices {
			fmt.Println(" ", devices[j].Name)
			fmt.Printf("%+v\n", devices[j])
		}
	}
}
