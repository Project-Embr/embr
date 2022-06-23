SRCFILES := *.go go.sum go.mod
PLUGIN_PATH := $(shell pwd)/../../submodules/plugins
PLUGIN_FILES := $(PLUGIN_PATH)/plugins/main/ptp/*.go $(PLUGIN_PATH)/plugins/ipam/host-local/*.go $(PLUGIN_PATH)/plugins/meta/firewall/*.go

all: cni worker

worker: $(SRCFILES)
	go build

clean:
	go clean
	export GOWORK=off;\
	cd ../../submodules/plugins && go clean

cni: $(PLUGIN_FILES)
	mkdir -p ../../submodules/plugins/bin
	export GOWORK=off;\
	cd ../../submodules/plugins/plugins/main/ptp && go build -o ${PLUGIN_PATH}/bin/ptp && \
	cd ../../ipam/host-local && go build -o ${PLUGIN_PATH}/bin/host-local &&\
	cd ../../meta/firewall && go build -o ${PLUGIN_PATH}/bin/firewall
	export GOWORK=off;\
	cd ../../submodules/tc-redirect-tap && $(MAKE) && mv tc-redirect-tap ../plugins/bin

format:
	go fmt

.PHONY: all worker cni clean format
