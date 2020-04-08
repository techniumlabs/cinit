NAME    := cinit
PACKAGE := github.com/techniumlabs/$(NAME)
GIT     := $(shell git rev-parse --short HEAD)
SOURCE_DATE_EPOCH ?= $(shell date +%s)
DATE    := $(shell date -u +%FT%T%Z)
VERSION  ?= v0.0.1
IMG_NAME := techniumlabs/cinit
IMAGE    := ${IMG_NAME}:${VERSION}

default: help

test:      ## Run all tests
	@go clean --testcache && go test ./...

cover:     ## Run test coverage suite
	@go test -race -covermode atomic -coverprofile=profile.cov ./...

build:     ## Builds the CLI
	@go build \
	-ldflags "-w -s -X ${PACKAGE}/cmd.version=${VERSION} -X ${PACKAGE}/cmd.commit=${GIT} -X ${PACKAGE}/cmd.date=${DATE}" \
	-a -tags netgo -o execs/${NAME} main.go

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":[^:]*?## "}; {printf "\033[38;5;69m%-30s\033[38;5;38m %s\033[0m\n", $$1, $$2}'
