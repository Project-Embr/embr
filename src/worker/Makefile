SRCFILES := *.go go.sum go.mod

all: worker cni

worker: $(SRCFILES)
	go build

clean:
	go clean
	cd ../../plugins && rm -rf bin

format:
	go fmt

.PHONY: all worker clean format
