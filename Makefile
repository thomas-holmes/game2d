GO_BIN?=go
CC=$(GO_BIN)
GOBUILD=$(CC) build
GOPKG='github.com/thomas-holmes/game2d/cmd/game2d'
MKFILE_PATH=$(abspath $(lastword $(MAKEFILE_LIST)))
ROOT=$(shell git rev-parse --show-toplevel)
SHA=$(shell git rev-parse HEAD)
VERSION=$(shell git describe --tags --always)
BINARY=graphics

all: test build

build:
	$(GOBUILD) $(GOPKG)

binclean:
	rm graphics | true

run: binclean build
	./game2d

test:
	$(GO_BIN) test ./...
