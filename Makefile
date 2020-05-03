VERSION = $(shell grep 'const version' cmd/version.go | sed -E 's/.*"(.+)"$$/v\1/')

.PHONY: all test clean build install

GOFLAGS ?= $(GOFLAGS:)

all: build

build:
	go build $(GOFLAGS)

clean:
	go clean $(GOFLAGS) -i
