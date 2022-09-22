GO ?= go
ROOT := $(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))
DIST = $(ROOT)/dist
TAG = 1.0
DOCKER_REPO = amikai

all: build lint

.PHONY: build
build: \
	build/docker/hello build/hello \
	build/buf

.PHONY: lint
lint: \
	lint/buf

.PHONY: build/hello
build/hello:
	cd $(ROOT)/src && $(GO) build -o $(DIST)/hello/server hello/server/main.go

.PHONY: build/docker/hello
build/docker/hello:
	docker build -f $(ROOT)/docker/hello.dockerfile -t $(DOCKER_REPO)/hello:$(TAG) $(ROOT)/src

.PHONY: build/buf
build/buf:
	docker run --mount type=bind,source=$(ROOT)/src,target=/build -w /build bufbuild/buf:1.8.0 build

.PHONY: lint/buf
lint/buf:
	docker run --mount type=bind,source=$(ROOT)/src,target=/build -w /build bufbuild/buf:1.8.0 lint

.PHONY: clean
clean:
	rm -rf $(DIST)
