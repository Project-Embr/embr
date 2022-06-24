package main

import (
	"github.com/google/uuid"
)

var etcdPrefix = "/"
var etcdEmbrPrefix = etcdPrefix + "embr/"

var workerID = uuid.New().String()[:8]
