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
	opts.CNIConfigPath = "../cni/conf.d/fcnet.conflist"
	opts.CNIPluginsPath = []string{"../cni/plugins"}
	opts.CNINetnsPath = "ext/netns"

	if err := runVM(context.Background(), opts); err != nil {
		log.Fatalf(err.Error())
	}
}
