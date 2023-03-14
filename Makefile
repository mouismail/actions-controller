# Makefile for actions-control

# Variables
VERSION = v1.0.0
BUILD_TIME = $(shell date -u +%Y%m%d.%H%M%S)
COMMIT_ID = $(shell git rev-parse HEAD)

# Targets
.PHONY: dev build start

dev:
	go run main.go

build:
	docker build -t  actions-control:latest . \
	--build-arg VERSION=$(VERSION) \
	--build-arg BUILD_TIME=$(BUILD_TIME) \
	--build-arg COMMIT_ID=$(COMMIT_ID)

start:
	docker run -d -p 3000:3000 actions-control:latest