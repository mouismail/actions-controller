# Makefile for actions-control

# Variables
VERSION = v1.0.0
BUILD_TIME = $(shell date -u +%Y%m%d.%H%M%S)
COMMIT_ID = $(shell git rev-parse HEAD)
HTTP_PORT = 3000

# Targets
.PHONY: all test dev build start

.PHONY: all
all::
	go mod tidy

release:: all;

start: all bin/sap-actions-control --bind-addr 0.0.0.0 --log-level debug

test:
	go test -v ./...

dev:
	go run main.go

build:
	go test -v ./...
	if [ $$? -ne 0 ]; then exit 1; fi
	docker build -t actions-control:latest . \
	--build-arg VERSION=$(VERSION) \
	--build-arg BUILD_TIME=$(BUILD_TIME) \
	--build-arg COMMIT_ID=$(COMMIT_ID) \
	--build-arg HTTP_PORT=$(HTTP_PORT)

start:
	docker run -d -p 3000:3000 actions-control:latest
