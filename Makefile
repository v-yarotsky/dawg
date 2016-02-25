.PHONY: all test build alfred clean

all: test build alfred

test:
	go test

build:
	mkdir -p build
	go build -o build/dawg ./cmd/dawg/dawg.go

alfred: build
	zip -j DAWG.alfredworkflow build/dawg alfred_workflow/*
clean:
	rm -f bin/*
	rm -f DAWG.alfredworkflow
