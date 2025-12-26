.PHONY: build test lint lint-fix fmt clean help

# Build variables
BINARY_NAME=bury-it
BUILD_DIR=bin
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS=-ldflags "-X github.com/deanhigh/bury-it/cmd.Version=$(VERSION)"

## build: Build the binary to bin/
build:
	@mkdir -p $(BUILD_DIR)
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) .

## test: Run all tests
test:
	go test -race -v ./...

## lint: Run golangci-lint
lint:
	golangci-lint run ./...

## lint-fix: Run golangci-lint with auto-fix
lint-fix:
	golangci-lint run --fix ./...

## fmt: Format code
fmt:
	go fmt ./...

## clean: Remove build artifacts
clean:
	rm -rf $(BUILD_DIR)
	go clean

## help: Show this help message
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@sed -n 's/^## //p' $(MAKEFILE_LIST) | column -t -s ':'
