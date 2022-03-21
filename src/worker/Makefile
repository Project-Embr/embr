SRCFILES := *.go go.sum go.mod

all: worker

worker: $(SRCFILES)
	go build

clean:
	go clean

format:
	go fmt

.PHONY: all worker clean format
