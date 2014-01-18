package cl11

import (
	"log"

	"testing"
)

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
			_, err = CreateContext(nil, []*Device{devices[j]}, func(err string, data []byte) {})
			if err != nil {
				t.Error(err)
			}
		}
	}
}

func ExampleContext() {
	platforms, err := GetPlatforms()
	if err != nil {
		log.Fatalln(err)
	}
	for i := range platforms {
		devices, err := platforms[i].GetDevices()
		if err != nil {
			log.Fatalln(err)
		}
		for j := range devices {
			_, err = CreateContext(nil, []*Device{devices[j]}, func(err string, data []byte) {})
			if err != nil {
				log.Fatalln(err)
			}
			log.Println("created context on", devices[j])
		}
	}
}
