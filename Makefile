GOCMD=go
GOBUILD=$(GOCMD) build

BINARY_NAME=go-httpd

GIT_COMMIT := $(shell git rev-parse HEAD)

build:
	$(GOBUILD) -o $(BINARY_NAME) -ldflags="-X main.gitCommit=$(GIT_COMMIT)"
