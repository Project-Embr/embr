package main

import (
	log "github.com/sirupsen/logrus"
	"time"
)

func main() {
	var runningVM []chan string

	etcdServer := startEtcd()
	defer etcdServer.Close()

	etcdClient := getClient()
	defer etcdClient.Close()

	SignalHandlers(etcdClient, etcdServer, &runningVM)
	startWatchers(etcdClient, &runningVM)

	time.Sleep(2 * time.Second)
	log.Fatal(<-etcdServer.Err()) //Blocking statement
}
