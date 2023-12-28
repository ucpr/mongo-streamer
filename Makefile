NAME := mongo-streamer
BUILD_DIR := ./build
GO ?= go

.PHONY: build
build: VERSION := $(shell git describe --tags --always --dirty)
build: REVISION := $(shell git rev-parse --short HEAD)
build: TIMESTAMP := $(shell date -u '+%Y-%m-%d_%I:%M:%S%p')
build:
	$(GO) build -ldflags "-X main.BuildVersion=$(VERSION) -X main.BuildRevision=$(REVISION) -X main.BuildTimestamp=$(TIMESTAMP)" -o $(BUILD_DIR)/$(NAME) ./cmd

.PHONY: test
test: PKG ?= ./...
test: FLAGS ?=
test:
	$(GO) test -race $(PKG) $(FLAGS)

.PHONY: integration-test
integration-test: PKG ?= ./...
integration-test:
	$(GO) test -race $(PKG) -tags=integration

.PHONY: generate
generate: PKG ?= ./...
generate:
	$(GO) generate $(PKG)
