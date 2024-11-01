# Makefile for the glue project

# Variables
APP_NAME := glue
SRC := $(wildcard *.go)
PKG := ./...
GO_FILES := $(shell find . -name '*.go' -type f)
VERSION := $(shell git describe --tags --always --dirty)
LDFLAGS := -s -w -X main.version=$(VERSION)
BUILD_DIR := build
BIN_DIR := $(BUILD_DIR)/bin
BINARY := $(BIN_DIR)/$(APP_NAME)

# Default target
.PHONY: all
all: build

# Build the application
.PHONY: build
build: deps
	@echo "Building $(APP_NAME)..."
	@mkdir -p $(BIN_DIR)
	GO111MODULE=on go build -ldflags "$(LDFLAGS)" -o $(BINARY) $(SRC)
	@echo "Build complete."

# Install the application
.PHONY: install
install: build
	@echo "Installing $(APP_NAME)..."
	GO111MODULE=on go install -ldflags "$(LDFLAGS)" ./...
	@echo "Installation complete."

# Run the application
.PHONY: run
run: build
	@echo "Running $(APP_NAME)..."
	@$(BINARY)

# Clean build files
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@echo "Clean complete."

# Fetch dependencies
.PHONY: deps
deps:
	@echo "Fetching dependencies..."
	GO111MODULE=on go mod tidy
	@echo "Dependencies fetched."

# Format code
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	go fmt $(PKG)
	@echo "Code formatted."

# Lint code
.PHONY: lint
lint:
	@echo "Linting code..."
	@golangci-lint run
	@echo "Linting complete."

# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	@go test -v $(PKG)
	@echo "Tests complete."

# Build a release (compressed binary)
.PHONY: release
release: build
	@echo "Creating release archive..."
	@tar -czvf $(BUILD_DIR)/$(APP_NAME)-$(VERSION).tar.gz -C $(BIN_DIR) $(APP_NAME)
	@echo "Release archive created at $(BUILD_DIR)/$(APP_NAME)-$(VERSION).tar.gz"

# Help message
.PHONY: help
help:
	@echo "Makefile for the $(APP_NAME) project"
	@echo ""
	@echo "Usage:"
	@echo "  make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  all       - Default target (build)"
	@echo "  build     - Build the application"
	@echo "  install   - Install the application"
	@echo "  run       - Run the application"
	@echo "  clean     - Clean build artifacts"
	@echo "  deps      - Fetch dependencies"
	@echo "  fmt       - Format the code"
	@echo "  lint      - Lint the code (requires golangci-lint)"
	@echo "  test      - Run tests"
	@echo "  release   - Build a release archive"
	@echo "  help      - Display this help message"
