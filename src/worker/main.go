package main

import (
	"context"
	"time"

	//log "github.com/sirupsen/logrus"
)

func main() {
	opts := newOptions()

	// These files must exist
	opts.FcKernelImage = "ext/alpine.bin"
	opts.FcRootDrivePath = "ext/rootfs.ext4"

	command := make(chan string, 1)
	err := make(chan error, 1)
	go runVM(context.Background(), opts, err, command)
	time.Sleep(5 * time.Second)
	command <- "shutdown"
	//create an error channel and halt the program until VM is started
	//possibly impliment an occaasional statuscheck routine aswell
}
