package main

import (
	"fmt"
	"strconv"
)

//Parses program arguments, Takes program args as a param and returns an object containing args
func newArgs(args []string) (*arguments, error) {
	a := arguments{VcpuCount: 1, MemSizeMib: 512}
	for x := 0; x < len(args); x++ {
		switch args[x] {
		case "-v":
			VcpuCount, err := strconv.ParseInt(args[x+1], 10, 64)
			if err != nil {
				fmt.Println("Invalid argument: usage -v integer")
				return &arguments{}, err
			}
			a.VcpuCount = VcpuCount
		case "-m":
			MemSizeMib, err := strconv.ParseInt(args[x+1], 10, 64)
			if err != nil {
				fmt.Println("Invalid argument: usage -m integer")
				return &arguments{}, err
			}
			a.MemSizeMib = MemSizeMib
		}
	}
	return &a, nil
}

type arguments struct {
	VcpuCount  int64
	MemSizeMib int64
}
