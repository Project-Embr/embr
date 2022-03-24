package main

import (
	"context"
	"os"

	log "github.com/sirupsen/logrus"
)

func main() {
	opts := newOptions()
	args := os.Args[1:]
	arguments, err := newArgs(args)
	if err != nil {
		log.Fatal(err.Error())
	}
	// These files must exist
	opts.FcKernelImage = "ext/alpine.bin"
	opts.FcRootDrivePath = "ext/rootfs.ext4"

	if err := runVM(context.Background(), opts, *arguments); err != nil {
		log.Fatalf(err.Error())
	}
}
