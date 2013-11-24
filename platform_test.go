package cl11

import (
	"fmt"
	"testing"
)

func TestGetPlatforms(t *testing.T) {
	platforms, err := GetPlatforms()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(len(platforms), "platforms:")
	for i := range platforms {
		err = platforms[i].GetInfo()
		if err != nil {
			t.Error(err)
			continue
		}
		fmt.Printf("\t%s\n", platforms[i].Name())
	}
}
