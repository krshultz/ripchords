.PHONY: build test vet

build:
	go build -ldflags "-X main.version=$(shell git describe --tags --always)" -o ripchords .

test:
	go test ./...

vet:
	go vet ./...
