package main

import (
	"context"
	"time"
	"fmt"
	log "github.com/sirupsen/logrus"
	client "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/server/v3/embed"
)

func startEtcd() *embed.Etcd {
	cfg := embed.NewConfig()
	cfg.Dir = "default.etcd"

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

func watchEmbrs(etcdClient *client.Client, runningVM *[]chan string) {
	go func() {
		watcher := etcdClient.Watch(context.Background(), etcdEmbrPrefix, client.WithPrefix())
		for resp := range watcher {
			log.Info("cakled")
			for _, ev := range resp.Events {
				log.Info(fmt.Sprintf("%s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value))
				(*runningVM) = append((*runningVM), createNewVM(etcdClient, ev.Kv.Value))
			}
		}
	}()
}

func startWatchers(etcdClient *client.Client, runningVM *[]chan string) {
	watchEmbrs(etcdClient, runningVM)
}
