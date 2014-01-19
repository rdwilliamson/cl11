package cl11

import (
	"fmt"
	"log"

	"testing"
)

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
			_, err = CreateContext([]*Device{devices[j]}, ContextProperties{}, func(err string, data []byte) {})
			if err != nil {
				log.Fatalln(err)
			}
			log.Println("created context on", devices[j])
		}
	}
}

// Creates all possible contexts.
func createContexts(t *testing.T) []*Context {
	var contexts []*Context
	platforms, err := GetPlatforms()
	if err != nil {
		t.Fatal(err)
	}
	for i := range platforms {
		devices, err := platforms[i].GetDevices()
		if err != nil {
			t.Error(err)
			continue
		}
		for j := range devices {
			c, err := CreateContext([]*Device{devices[j]}, ContextProperties{}, func(err string, data []byte) {
				fmt.Println(err)
				fmt.Println(data)
			})
			if err != nil {
				t.Error(err)
				continue
			}
			contexts = append(contexts, c)
		}
	}
	return contexts
}
