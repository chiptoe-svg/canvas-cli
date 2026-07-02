# Development

Resources for contributing to Canvas CLI development.

## Overview

Canvas CLI is built with:

- **Go 1.25+** - Primary language
- **Cobra** - CLI framework
- **Viper** - Configuration management

## Topics

<div class="grid cards" markdown>

-   :material-folder-multiple:{ .lg .middle } **Architecture**

    ---

    Learn about Canvas CLI's internal structure and design decisions

    [:octicons-arrow-right-24: Architecture](architecture.md)

-   :material-source-pull:{ .lg .middle } **Contributing**

    ---

    Guidelines for contributing code, documentation, and bug reports

    [:octicons-arrow-right-24: Contributing](contributing.md)

-   :material-chart-donut:{ .lg .middle } **API Coverage**

    ---

    How coverage is measured against the official Canvas API spec, and what's left

    [:octicons-arrow-right-24: API Coverage](api-coverage.md)

</div>

## Quick Start for Contributors

### Prerequisites

- Go 1.25 or later
- Git
- Make (optional but recommended)

### Setup

```bash
# Clone the repository
git clone https://github.com/jjuanrivvera/canvas-cli.git
cd canvas-cli

# Install dependencies
go mod download

# Build
make build

# Run tests
make test
```

### Development Workflow

```bash
# Create a feature branch
git checkout -b feature/my-feature

# Make changes and run tests
make dev
make test

# Commit and push
git add .
git commit -m "feat: add my feature"
git push origin feature/my-feature
```

### Available Make Targets

| Target | Description |
|--------|-------------|
| `make build` | Build binary to `bin/canvas` |
| `make dev` | Build with fmt and vet |
| `make test` | Run all tests |
| `make test-integration` | Run binary-level integration tests |
| `make test-coverage` | Run tests with coverage |
| `make check` | Run everything CI runs (lint, security, tests, spec checks) |
| `make lint` | Run golangci-lint |
| `make fmt` | Format code |
| `make setup-hooks` | Install git pre-commit hooks |
| `make spec-sync` | Refresh the committed Canvas API spec manifest |
| `make spec-coverage` | Report documented-but-unimplemented endpoints |
| `make docs-gen` | Regenerate CLI reference docs and sync the changelog |
| `make install` | Install to `/usr/local/bin` |

## Project Structure

```
canvas-cli/
├── cmd/canvas/       # Entry point
├── commands/         # Cobra commands
├── internal/
│   ├── api/          # Canvas API client
│   ├── auth/         # Authentication
│   ├── config/       # Configuration
│   ├── cache/        # Response caching
│   ├── batch/        # Batch operations
│   └── output/       # Output formatters
├── tools/            # Development tools
└── docs/             # Documentation
```

## Integration Tests

A binary-level integration suite lives in `test/integration/` and is guarded by the `integration` build tag, so it is invisible to the normal `go test ./...` run. The suite compiles the binary once in `TestMain`, then exercises it as a black box against a lightweight `httptest` mock Canvas server. It covers version output, help, unknown-command errors, JSON/CSV/table output, auth failure, alias expansion, context propagation, `--dry-run` token redaction, and the REPL help surface.

Run it with:

```bash
go test -tags integration -v -timeout 5m ./test/integration/
# or
make test-integration
```

## Code Style

Canvas CLI follows standard Go conventions:

- Run `gofmt` before committing
- Use `golangci-lint` for static analysis
- Write tests for new features
- Document exported functions

## Getting Help

- [GitHub Issues](https://github.com/jjuanrivvera/canvas-cli/issues) - Bug reports and feature requests
- [GitHub Discussions](https://github.com/jjuanrivvera/canvas-cli/discussions) - Questions and ideas
