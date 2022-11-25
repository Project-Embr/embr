package main

import (
	log "github.com/sirupsen/logrus"
)

func main() {
	var runningVM []chan string

	etcdServer := startEtcd()
	defer etcdServer.Close()

	etcdClient := getClient()
	defer etcdClient.Close()

	SignalHandlers(etcdClient, etcdServer, &runningVM)

	startWatchers(etcdClient, &runningVM)

	log.Fatal(<-etcdServer.Err()) //Blocking statement
}
