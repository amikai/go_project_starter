GO ?= go
ROOT := $(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))
DIST = $(ROOT)/dist

.PHONY: build/hello
build/hello:
	cd $(ROOT)/src && $(GO) build -o $(DIST)/hello/server hello/server/main.go

.PHONY: clean
clean:
	rm -rf $(DIST)
