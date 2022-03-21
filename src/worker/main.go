package main

import (
	"context"

	log "github.com/sirupsen/logrus"
)

func main() {
	opts := newOptions()

	// These files must exist
	opts.FcKernelImage = "ext/alpine.bin"
	opts.FcRootDrivePath = "ext/rootfs.ext4"

	if err := runVM(context.Background(), opts); err != nil {
		log.Fatalf(err.Error())
	}
}
