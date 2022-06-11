package main

import (
	log "github.com/sirupsen/logrus"
	"context"
	client "go.etcd.io/etcd/client/v3"
	"encoding/json"
)

//Creates a new firecracker VM and returns the command channel
func createNewVM(etcdClient *client.Client, inputOps []byte) chan string{
	opts := newOptions()

	// These files must exist
	opts.FcKernelImage = "ext/alpine.bin"
	opts.FcRootDrivePath = "ext/rootfs.ext4"
	err := json.Unmarshal(inputOps, &opts)
	command := make(chan string, 1)
	errChan := make(chan error, 1)
	if(err == nil){
		go runVM(context.Background(), opts, errChan, command)
	}else{
		log.Info("Invalid opts file")
		return nil
	}
	if (<- errChan == nil){
		log.Info("machine started successfully")
	} else{
		log.Warn("Failed to create machine")
	}
	return command
}