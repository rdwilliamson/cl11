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
		_, err = platforms[i].GetDevices()
		if err != nil {
			t.Fatal(err)
		}
	}
}

func ExampleDevices() {
	platforms, err := GetPlatforms()
	if err != nil {
		fmt.Println(err)
		return
	}
	for i := range platforms {
		fmt.Printf("%+v\n", platforms[i])
		devices, err := platforms[i].GetDevices()
		if err != nil {
			fmt.Println(err)
			return
		}
		for j := range devices {
			fmt.Printf("%+v\n", devices[j])
		}
	}
}
