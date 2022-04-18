SRCFILES := *.go go.sum go.mod

all: worker cni

worker: $(SRCFILES)
	go build

cni: 
	export GOWORK=off;\
	cd ../../plugins && ./build_linux.sh
	export GOWORK=off;\
	cd ../../tc-redirect-tap && $(MAKE) && mv tc-redirect-tap ../plugins/bin

clean:
	go clean
	cd ../../plugins && rm -rf bin

format:
	go fmt

.PHONY: all worker cni clean format
