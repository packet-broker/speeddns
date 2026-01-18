# Makefile for speeddns

VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
DATE ?= $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')

LDFLAGS := -ldflags "-X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE) -s -w"

BINARY := speeddns
BUILD_DIR := build

.PHONY: all build clean test lint install run deps tidy

all: build

# Development build for current platform
build:
	go build $(LDFLAGS) -o $(BINARY) ./cmd/speeddns

# Cross-compilation for Debian targets
build-linux-amd64:
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY)-linux-amd64 ./cmd/speeddns

build-linux-arm64:
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY)-linux-arm64 ./cmd/speeddns

build-all: build-linux-amd64 build-linux-arm64

# Development
run:
	go run ./cmd/speeddns

test:
	go test -v -race ./...

test-cover:
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

lint:
	golangci-lint run ./...

fmt:
	go fmt ./...

vet:
	go vet ./...

# Installation
install:
	go install $(LDFLAGS) ./cmd/speeddns

# Cleanup
clean:
	rm -f $(BINARY)
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html

# Dependencies
deps:
	go mod download
	go mod verify

tidy:
	go mod tidy

# Release builds (static, stripped)
release-linux-amd64:
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -a -o $(BUILD_DIR)/$(BINARY)-linux-amd64 ./cmd/speeddns

release-linux-arm64:
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -a -o $(BUILD_DIR)/$(BINARY)-linux-arm64 ./cmd/speeddns

release: release-linux-amd64 release-linux-arm64
	@echo "Release binaries built in $(BUILD_DIR)/"
	@ls -lh $(BUILD_DIR)/

# Create distribution archive
dist: release
	cd $(BUILD_DIR) && tar -czvf $(BINARY)-$(VERSION)-linux-amd64.tar.gz $(BINARY)-linux-amd64
	cd $(BUILD_DIR) && tar -czvf $(BINARY)-$(VERSION)-linux-arm64.tar.gz $(BINARY)-linux-arm64
	@echo "Distribution archives created in $(BUILD_DIR)/"

# Quick test run
quick-test: build
	./$(BINARY) -p -n 2 -c 5

# Help
help:
	@echo "Available targets:"
	@echo "  build             - Build for current platform"
	@echo "  build-linux-amd64 - Build for Linux x86_64"
	@echo "  build-linux-arm64 - Build for Linux ARM64"
	@echo "  build-all         - Build for all Linux platforms"
	@echo "  run               - Run without building"
	@echo "  test              - Run tests"
	@echo "  lint              - Run linter"
	@echo "  install           - Install to GOPATH/bin"
	@echo "  clean             - Remove build artifacts"
	@echo "  deps              - Download dependencies"
	@echo "  tidy              - Tidy go.mod"
	@echo "  release           - Build static release binaries"
	@echo "  dist              - Create distribution archives"
	@echo "  quick-test        - Build and run quick test"
