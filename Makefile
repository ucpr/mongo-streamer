NAME := mongo-streamer
BUILD_DIR := ./build
GO ?= go

.PHONY: build
build:
	$(GO) build -o $(BUILD_DIR)/$(NAME) ./cmd

.PHONY: test
test: PKG ?= ./...
test:
	$(GO) test -race $(PKG)

.PHONY: generate
generate: PKG ?= ./...
generate:
	$(GO) generate $(PKG)
