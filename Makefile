GO ?= go
ROOT := $(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))
DIST = $(ROOT)/dist
TAG = 1.0
DOCKER_REPO = amikai

.PHONY: build/hello
build/hello:
	cd $(ROOT)/src && $(GO) build -o $(DIST)/hello/server hello/server/main.go

.PHONY: build/docker/hello
build/docker/hello:
	docker build -f $(ROOT)/docker/hello.dockerfile -t $(DOCKER_REPO)/hello:$(TAG) $(ROOT)/src

.PHONY: clean
clean:
	rm -rf $(DIST)
