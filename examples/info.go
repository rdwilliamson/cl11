package main

import (
	"fmt"
	"os"

	"github.com/davecgh/go-spew/spew"
	cl "github.com/rdwilliamson/cl11"
)

func main() {
	platforms, err := cl.GetPlatforms()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	spew.Dump(platforms)
}
