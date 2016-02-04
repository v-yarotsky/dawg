.PHONY: all test build

all: test build

test:
	go test

build:
	mkdir -p build
	go build -o build/dawg ./cmd/dawg/dawg.go
