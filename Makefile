.PHONY: build clean install help

# Variables
BINARY_NAME=wagmi
VERSION=1.0.0
BUILD_TIME=$(shell date +%FT%T%z)
LDFLAGS=-ldflags "-X github.com/hauchu1196/wagmi/cmd.version=$(VERSION) -X github.com/hauchu1196/wagmi/cmd.buildTime=$(BUILD_TIME)"

# Default target
all: build

# Build the project
build:
	@echo "Building $(BINARY_NAME)..."
	go build $(LDFLAGS) -o bin/$(BINARY_NAME) main.go

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -rf bin/

# Install dependencies
install:
	@echo "Installing dependencies..."
	go mod download

# Show help
help:
	@echo "Available commands:"
	@echo "  make build    - Build the project"
	@echo "  make clean    - Clean build artifacts"
	@echo "  make install  - Install dependencies"
	@echo "  make help     - Show this help message"

# Build for multiple platforms
build-all: build-linux build-darwin-intel build-darwin-arm build-windows

build-linux:
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o bin/$(BINARY_NAME)-linux-amd64 main.go

build-darwin-intel:
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o bin/$(BINARY_NAME)-darwin-intel main.go

build-darwin-arm:
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o bin/$(BINARY_NAME)-darwin-arm main.go

build-windows:
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o bin/$(BINARY_NAME)-windows-amd64.exe main.go

# Create release package
release: clean build-all
	@echo "Creating release package..."
	tar -czf bin/$(BINARY_NAME)-$(VERSION).tar.gz bin/$(BINARY_NAME)-*

# Default target
.DEFAULT_GOAL := help 