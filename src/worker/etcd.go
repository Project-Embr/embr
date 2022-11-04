package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	client "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/server/v3/embed"
)

func startEtcd() *embed.Etcd {
	cfg := embed.NewConfig()
	cfg.Dir = "/tmp/etcd"
	cfg.LogLevel = "debug"

	etcdServer, err := embed.StartEtcd(cfg)
	if err != nil {
		log.Fatal(err)
	}

	select {
	case <-etcdServer.Server.ReadyNotify():
		log.Info("etcd server is ready!")
	case <-time.After(60 * time.Second):
		log.Fatal("server took too long to start, stopping!")
	}

	return etcdServer
}

func getClient() *client.Client {
	etcdClient, err := client.New(client.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}

	return etcdClient
}

func setupNode() {
	log.Info(workerID)
}

func watchEmbrs(etcdClient *client.Client, runningVM *[]chan string) {
	go func() {
		watcher := etcdClient.Watch(context.Background(), etcdEmbrPrefix, client.WithPrefix())
		for resp := range watcher {
			for _, ev := range resp.Events {
				log.Info(fmt.Sprintf("%s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value))
				(*runningVM) = append((*runningVM), startVM(etcdClient, ev.Kv.Value, runningVM))
			}
		}
	}()
}

func startWatchers(etcdClient *client.Client, runningVM *[]chan string) {
	watchEmbrs(etcdClient, runningVM)
}
func startVM(etcdClient *client.Client, inputOps []byte, runningVM *[]chan string) chan string {
	opts := newOptions()

	// These files must exist
	opts.FcKernelImage = "/tmp/images/alpine.bin"
	opts.FcRootDrivePath = "/tmp/images/rootfs.ext4"
	opts.CNIConfigPath = "cni/conf.d/"
	opts.CNIPluginsPath = []string{"submodules/plugins/bin/"}
	opts.CNINetnsPath = "/tmp/netns"

	err := json.Unmarshal(inputOps, &opts)
	command := make(chan string, 1)
	errChan := make(chan error, 1)
	if err != nil {
		fmt.Println("Unable to convert the JSON string to a struct")
		return nil
	} else {
		go runVM(context.Background(), opts, errChan, command)
	}

	if <-errChan == nil {
		log.Info("Machine Started Sucessfully")
	} else {
		log.Warn("Failed To Create Machine")
	}

	return command
}
