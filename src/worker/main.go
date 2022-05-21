package main

import (
	"time"
	"fmt"
	//log "github.com/sirupsen/logrus"
)

//Basic list node to store channels
type Node struct{
	next *Node
	channel chan string
}

func main() {
	opts := newOptions()

	// These files must exist
	opts.FcKernelImage = "ext/alpine.bin"
	opts.FcRootDrivePath = "ext/rootfs.ext4"

	head := &Node{
		next: nil,
		channel: nil,
	}
	head.channel = createNewVM(opts)
	if(head.channel == nil){
		fmt.Errorf("Error creating machine")
	}
	time.Sleep(5 * time.Second)
	head.channel <- "shutdown"
}
