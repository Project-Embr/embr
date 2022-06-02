SRCFILES := *.go go.sum go.mod

all: cni worker

worker: $(SRCFILES)
	go build

clean:
	go clean
	cd ../../plugins && rm -rf bin

cni:
	export GOWORK=off;\
	cd ../../submodules/plugins && ./build_linux.sh
	# go build -o bin/ptp 
	export GOWORK=off;\
	cd ../../submodules/tc-redirect-tap && $(MAKE) && mv tc-redirect-tap ../plugins/bin

format:
	go fmt

.PHONY: all worker cni clean format
