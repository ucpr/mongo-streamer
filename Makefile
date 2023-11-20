NAME := mongo-streamer
BUILD_DIR := ./build
GO ?= go

.PHONY: build
build:
	$(GO) build -o $(BUILD_DIR)/$(NAME) ./cmd/main.go

.PHONY: test
test:
	go test -race ./...
