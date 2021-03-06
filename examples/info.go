package main

import (
	"fmt"
	"os"

	cl "github.com/rdwilliamson/cl11"

	"github.com/davecgh/go-spew/spew"
)

func main() {
	platforms, err := cl.GetPlatforms()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	spew.Dump(platforms)
}
