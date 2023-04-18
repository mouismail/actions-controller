# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=actions-controller
IMAGE_NAME = $(BINARY_NAME)
IMAGE_TAG =latest
GHES_APP_PRIVATE_KEY=keys/actions-control.2023-03-30.private-key.pem

all: build

build:
	$(GOBUILD) -o $(BINARY_NAME) -v

docker-build:
	docker build -t $(IMAGE_NAME):$(IMAGE_TAG) .

docker-run:
	docker run -it -p 3000:3000 -e GHES_APP_PRIVATE_KEY=$(GHES_APP_PRIVATE_KEY) -e GHES_APP_WEBHOOK_SECRET=$(GHES_APP_WEBHOOK_SECRET) $(IMAGE_NAME):$(IMAGE_TAG)

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

test:
	$(GOTEST) -v ./...

run:
	$(GOBUILD) -o $(BINARY_NAME) -v
	./$(BINARY_NAME)

deps:
	$(GOGET) -v ./...

.PHONY: all build clean test run deps
