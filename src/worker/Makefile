SRCFILES := *.go go.sum go.mod

all: worker

worker: $(SRCFILES)
	go build

clean:
	go clean

format:
	go fmt

etcdstart:
	etcdctl put /embr/config '{"VCpuCount":1,"MemSizeMib":128}'

.PHONY: all worker clean format
