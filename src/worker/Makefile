SRCFILES := *.go go.sum go.mod
PLUGIN_PATH := $(shell pwd)/submodules/plugins
PLUGIN_FILES := $(PLUGIN_PATH)/plugins/main/ptp/*.go $(PLUGIN_PATH)/plugins/ipam/host-local/*.go $(PLUGIN_PATH)/plugins/meta/firewall/*.go

all: cni worker

worker: $(SRCFILES)
	go build

clean:
	go clean
	export GOWORK=off;\
	rm -rf $(PLUGIN_PATH)/bin

cni: $(PLUGIN_FILES)
	git submodule update
	mkdir -p $(PLUGIN_PATH)/bin
	export GOWORK=off; cd $(PLUGIN_PATH)/plugins/main/ptp && go build -o ${PLUGIN_PATH}/bin/ptp
	export GOWORK=off; cd $(PLUGIN_PATH)/plugins/ipam/host-local && go build -o ${PLUGIN_PATH}/bin/host-local
	export GOWORK=off; cd $(PLUGIN_PATH)/plugins/meta/firewall && go build -o ${PLUGIN_PATH}/bin/firewall
	export GOWORK=off; cd submodules/tc-redirect-tap \
		&& $(MAKE) \
		&& mv tc-redirect-tap ../plugins/bin

format:
	go fmt

etcdstart:
	etcdctl put /embr/config '{"VCpuCount":1,"MemSizeMib":128}'

.PHONY: all worker clean cni format etcdstart
