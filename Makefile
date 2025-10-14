.PHONY: all build test clean install uninstall fmt vet lint release-build

# Variables
BINARY_NAME=stamp
ALIAS_NAME=nid
GO=go
GOFLAGS=
INSTALL_PATH=/usr/local/bin

# Get version info
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
DATE := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)

# Build flags
LDFLAGS := -ldflags "-X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)"

# Default target
all: build

# Build the binary
build:
	$(GO) build $(GOFLAGS) $(LDFLAGS) -o $(BINARY_NAME) cmd/stamp/main.go

# Build for multiple platforms
release-build:
	@echo "Building for multiple platforms..."
	@mkdir -p dist

	@echo "Building for macOS (amd64)..."
	GOOS=darwin GOARCH=amd64 $(GO) build $(LDFLAGS) -o dist/$(BINARY_NAME)-darwin-amd64 cmd/stamp/main.go

	@echo "Building for macOS (arm64)..."
	GOOS=darwin GOARCH=arm64 $(GO) build $(LDFLAGS) -o dist/$(BINARY_NAME)-darwin-arm64 cmd/stamp/main.go

	@echo "Building for Linux (amd64)..."
	GOOS=linux GOARCH=amd64 $(GO) build $(LDFLAGS) -o dist/$(BINARY_NAME)-linux-amd64 cmd/stamp/main.go

	@echo "Building for Linux (arm64)..."
	GOOS=linux GOARCH=arm64 $(GO) build $(LDFLAGS) -o dist/$(BINARY_NAME)-linux-arm64 cmd/stamp/main.go

	@echo "Building for Windows (amd64)..."
	GOOS=windows GOARCH=amd64 $(GO) build $(LDFLAGS) -o dist/$(BINARY_NAME)-windows-amd64.exe cmd/stamp/main.go

	@echo "Build complete! Binaries are in ./dist/"

# Run tests
test:
	$(GO) test -v -race -coverprofile=coverage.out ./...

# Run tests with coverage report
test-coverage: test
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Clean build artifacts
clean:
	@rm -f $(BINARY_NAME)
	@rm -rf dist/
	@rm -f coverage.out coverage.html
	@echo "Clean complete"

# Format code
fmt:
	$(GO) fmt ./...

# Run go vet
vet:
	$(GO) vet ./...

# Run golangci-lint (if installed)
lint:
	@which golangci-lint > /dev/null 2>&1 || (echo "golangci-lint not installed. Install from https://golangci-lint.run/usage/install/" && exit 1)
	golangci-lint run

# Install binary to system
install: build
	@echo "Installing $(BINARY_NAME) to $(INSTALL_PATH)..."
	@sudo cp $(BINARY_NAME) $(INSTALL_PATH)/$(BINARY_NAME)
	@sudo chmod 755 $(INSTALL_PATH)/$(BINARY_NAME)

	@echo "Creating symlink $(ALIAS_NAME) -> $(BINARY_NAME)..."
	@sudo ln -sf $(INSTALL_PATH)/$(BINARY_NAME) $(INSTALL_PATH)/$(ALIAS_NAME)

	@echo "Installation complete!"
	@echo "You can now use '$(BINARY_NAME)' or '$(ALIAS_NAME)' from anywhere"

# Uninstall binary from system
uninstall:
	@echo "Removing $(BINARY_NAME) and $(ALIAS_NAME) from $(INSTALL_PATH)..."
	@sudo rm -f $(INSTALL_PATH)/$(BINARY_NAME)
	@sudo rm -f $(INSTALL_PATH)/$(ALIAS_NAME)
	@echo "Uninstall complete"

# Development build (quick rebuild)
dev:
	$(GO) build -o $(BINARY_NAME) cmd/stamp/main.go

# Run the binary
run: build
	./$(BINARY_NAME)

# Check dependencies
deps:
	$(GO) mod download
	$(GO) mod tidy

# Update dependencies
update-deps:
	$(GO) get -u ./...
	$(GO) mod tidy

# Show help
help:
	@echo "Available targets:"
	@echo "  make build         - Build the binary"
	@echo "  make test          - Run tests"
	@echo "  make test-coverage - Run tests with coverage report"
	@echo "  make clean         - Remove build artifacts"
	@echo "  make install       - Install binary to $(INSTALL_PATH)"
	@echo "  make uninstall     - Remove binary from $(INSTALL_PATH)"
	@echo "  make fmt           - Format code"
	@echo "  make vet           - Run go vet"
	@echo "  make lint          - Run golangci-lint"
	@echo "  make release-build - Build for multiple platforms"
	@echo "  make dev           - Quick development build"
	@echo "  make run           - Build and run"
	@echo "  make deps          - Download dependencies"
	@echo "  make update-deps   - Update dependencies"
	@echo "  make help          - Show this help message"