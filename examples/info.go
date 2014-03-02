package main

import (
	"fmt"

	cl "github.com/rdwilliamson/cl11"

	"github.com/davecgh/go-spew/spew"
)

func main() {
	platforms, err := cl.GetPlatforms()
	if err != nil {
		fmt.Println(err)
		return
	}
	spew.Dump(platforms)
}
