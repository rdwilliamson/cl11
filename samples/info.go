package main

import (
	"fmt"
	"os"

	cl "github.com/rdwilliamson/cl11"
)

func check(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	platforms, err := cl.GetPlatforms()
	check(err)
	for _, v := range platforms {
		fmt.Printf("%+v\n", v)
		devices, err := v.GetDevices()
		check(err)
		for _, v := range devices {
			fmt.Printf("%+v\n", v)
		}
	}
}
