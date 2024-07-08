NAME := mongo-streamer
BUILD_DIR := ./build
GO ?= go
BIN := $(abspath ./bin)

$(BIN)/wire:
	GOBIN=$(BIN) go install github.com/google/wire/cmd/wire@latest
$(BIN)/mockgen:
	GOBIN=$(BIN) go install go.uber.org/mock/mockgen@latest

.PHONY: build
build: VERSION := $(shell git describe --tags --always --dirty)
build: REVISION := $(shell git rev-parse --short HEAD)
build: TIMESTAMP := $(shell date -u '+%Y-%m-%d_%I:%M:%S%p')
build: STAMP_PKG := github.com/ucpr/mongo-streamer/pkg/stamp
build:
	$(GO) build -ldflags "-X $(STAMP_PKG).BuildVersion=$(VERSION) -X $(STAMP_PKG).BuildRevision=$(REVISION) -X $(STAMP_PKG).BuildTimestamp=$(TIMESTAMP)" -o $(BUILD_DIR)/$(NAME) ./cmd

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
generate: $(BIN)/wire $(BIN)/mockgen
generate: PKG ?= ./...
generate:
	GOBIN=$(BIN) $(GO) generate $(PKG)
