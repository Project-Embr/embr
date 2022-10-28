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
	mkdir -p $(PLUGIN_PATH)/bin
	export GOWORK=off; cd $(PLUGIN_PATH)/plugins/main/ptp && go build -o ${PLUGIN_PATH}/bin/ptp
	export GOWORK=off; cd $(PLUGIN_PATH)/plugins/ipam/host-local && go build -o ${PLUGIN_PATH}/bin/host-local
	export GOWORK=off; cd $(PLUGIN_PATH)/plugins/meta/firewall && go build -o ${PLUGIN_PATH}/bin/firewall
	export GOWORK=off; cd submodules/tc-redirect-tap \
		&& $(MAKE) \
		&& mv tc-redirect-tap ../plugins/bin

format:
	go fmt

dockerbuild:
	docker build -t project-embr/embr:worker-latest .

dockerrun:
	docker run \
	--device /dev/kvm:/dev/kvm:rw \
	-p 127.0.0.1:2379:2379 \
	-v $$(pwd):/embr \
	project-embr/embr:worker-latest \
	# --cap-add=net_admin \
	# 	--net host \

dockerinteractive:
	docker run \
	--device /dev/kvm:/dev/kvm:rw \
	-p 2379:2379 \
	--user root \
	--cap-add=NET_ADMIN \
	--cap-add=net_admin \
	-v $$(pwd):/embr \
	-it \
	--entrypoint "" \
	project-embr/embr:worker-latest \
	bash
	# --cap-add=net_admin \
	# 	--net host \

.PHONY: all worker cni clean format
