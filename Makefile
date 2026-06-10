.PHONY: help build test clean install uninstall fmt lint vet run deps setup-hooks docs-gen docs-serve docs-build docs-deploy

# Variables
BINARY_NAME=canvas
VERSION?=dev
COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
BUILD_DATE=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.Commit=$(COMMIT) -X main.BuildDate=$(BUILD_DATE)"

# Default target
help:
	@echo "Canvas CLI - Makefile targets:"
	@echo ""
	@echo "  make build        - Build the binary"
	@echo "  make install      - Install the binary to /usr/local/bin"
	@echo "  make uninstall    - Remove the binary from /usr/local/bin"
	@echo "  make test         - Run tests"
	@echo "  make test-coverage - Run tests with coverage"
	@echo "  make clean        - Remove build artifacts"
	@echo "  make fmt          - Format code"
	@echo "  make lint         - Run linter"
	@echo "  make vet          - Run go vet"
	@echo "  make run          - Build and run the CLI"
	@echo "  make deps         - Download dependencies"
	@echo "  make release      - Build snapshot release for all platforms (GoReleaser)"
	@echo "  make setup-hooks  - Install git pre-commit hooks"
	@echo ""
	@echo "Documentation:"
	@echo "  make docs-gen     - Generate CLI reference documentation"
	@echo "  make docs-serve   - Serve documentation locally"
	@echo "  make docs-build   - Build documentation site"
	@echo "  make docs-deploy  - Deploy documentation to GitHub Pages"
	@echo ""

# Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	@go build $(LDFLAGS) -o bin/$(BINARY_NAME) ./cmd/canvas
	@echo "✓ Build complete: bin/$(BINARY_NAME)"

# Install the binary
install: build
	@echo "Installing $(BINARY_NAME) to /usr/local/bin..."
	@sudo cp bin/$(BINARY_NAME) /usr/local/bin/
	@echo "✓ Installed: /usr/local/bin/$(BINARY_NAME)"

# Uninstall the binary
uninstall:
	@echo "Removing $(BINARY_NAME) from /usr/local/bin..."
	@sudo rm -f /usr/local/bin/$(BINARY_NAME)
	@echo "✓ Uninstalled"

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	@go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...
	@go tool cover -html=coverage.txt -o coverage.html
	@echo "✓ Coverage report: coverage.html"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin/
	@rm -rf dist/
	@rm -f coverage.txt coverage.html
	@echo "✓ Clean complete"

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...
	@echo "✓ Format complete"

# Run linter (requires golangci-lint)
lint:
	@echo "Running linter..."
	@which golangci-lint > /dev/null || (echo "golangci-lint not installed. Install: https://golangci-lint.run/usage/install/" && exit 1)
	@golangci-lint run
	@echo "✓ Lint complete"

# Run go vet
vet:
	@echo "Running go vet..."
	@go vet ./...
	@echo "✓ Vet complete"

# Build and run
run: build
	@./bin/$(BINARY_NAME)

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy
	@echo "✓ Dependencies downloaded"

# Snapshot release build via GoReleaser (same pipeline as tagged releases:
# ldflags, archives, checksums) without publishing anything.
release:
	@which goreleaser > /dev/null || (echo "goreleaser not installed. Install: https://goreleaser.com/install/" && exit 1)
	@echo "Building snapshot release with GoReleaser..."
	@goreleaser release --snapshot --clean
	@echo "✓ Snapshot release built in dist/"

# Development build with verbose output
dev: fmt vet
	@go build -v $(LDFLAGS) -o bin/$(BINARY_NAME) ./cmd/canvas

# Setup git hooks
setup-hooks:
	@echo "Setting up git hooks..."
	@git config core.hooksPath .githooks
	@echo "✓ Git hooks installed (.githooks/pre-commit)"

# Documentation targets
docs-gen:
	@echo "Generating CLI reference documentation..."
	@go run ./tools/gendocs/main.go
	@cp CHANGELOG.md docs/changelog.md
	@echo "✓ CLI docs generated in docs/commands/ (changelog synced)"

docs-serve: docs-gen
	@echo "Serving documentation at http://localhost:8000..."
	@mkdocs serve

docs-build: docs-gen
	@echo "Building documentation site..."
	@mkdocs build --strict
	@echo "✓ Documentation built in site/"

docs-deploy: docs-gen
	@echo "Deploying documentation to GitHub Pages..."
	@mkdocs gh-deploy --force
	@echo "✓ Documentation deployed"
