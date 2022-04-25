package main

import log "github.com/sirupsen/logrus"

func main() {
	etcdServer := startEtcd()
	defer etcdServer.Close()

	etcdClient := getClient()
	defer etcdClient.Close()

	watchKey(etcdClient, "config")

	log.Fatal(<-etcdServer.Err()) //Blocking statement
}
