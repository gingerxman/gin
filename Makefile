VERSION = $(shell grep 'const version' cmd/version.go | sed -E 's/.*"(.+)"$$/v\1/')

.PHONY: all test clean build install

GOFLAGS ?= $(GOFLAGS:)

all: build

build:
	go build $(GOFLAGS) gin.go

clean:
	go clean $(GOFLAGS) -i gin.go

install:
	rm `which gin`
	go install ${GOFLAGS} gin.go
