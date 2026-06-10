<div align="center" markdown>

![Canvas CLI](assets/images/logo.svg){ width="280" }

**A powerful command-line interface for Canvas LMS**

[![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![Release](https://img.shields.io/github/v/release/jjuanrivvera/canvas-cli)](https://github.com/jjuanrivvera/canvas-cli/releases)

</div>

## Features

- **Comprehensive API Coverage** - Access courses, assignments, submissions, users, modules, pages, and more
- **Multiple Output Formats** - Table, JSON, YAML, and CSV output
- **Bulk Operations** - Grade submissions in bulk from CSV files
- **Course Synchronization** - Sync content between Canvas instances
- **Intelligent Caching** - Fast responses with automatic cache invalidation
- **Secure Authentication** - OAuth 2.0 with PKCE flow, tokens stored in your system keyring
- **MCP Integration** - Use Canvas CLI as an MCP server for AI coding assistants
- **AI Agent Skill** - Bundled skill that teaches Claude Code, Cursor, and other agents how to drive the CLI
- **Docker Image** - Distroless image on GHCR for containerized and CI usage
- **Signed Releases** - cosign-signed checksums and SBOMs on every release

## Quick Start

=== "macOS (Homebrew)"

    ```bash
    brew tap jjuanrivvera/canvas-cli
    brew install canvas-cli
    ```

=== "Go Install"

    ```bash
    go install github.com/jjuanrivvera/canvas-cli/cmd/canvas@latest
    ```

=== "Docker"

    ```bash
    docker run --rm ghcr.io/jjuanrivvera/canvas-cli:latest version
    ```

=== "Binary Download"

    Download the latest release from [GitHub Releases](https://github.com/jjuanrivvera/canvas-cli/releases).
    Checksums are signed with cosign — see [Installation](getting-started/installation.md) for verification steps.

Then authenticate with Canvas:

```bash
canvas auth login
```

## Example Usage

```bash
# List your courses
canvas courses list

# List assignments for a course
canvas assignments list --course-id 123

# Grade a submission
canvas submissions grade 456 --course-id 123 --score 95

# Bulk grade from CSV
canvas submissions bulk-grade --course-id 123 --csv-file grades.csv

# Export data as JSON
canvas users list --course-id 123 --output json
```

## Documentation

<div class="grid cards" markdown>

-   :material-rocket-launch:{ .lg .middle } **Getting Started**

    ---

    Install Canvas CLI and authenticate with your Canvas instance

    [:octicons-arrow-right-24: Installation](getting-started/installation.md)

-   :material-book-open-variant:{ .lg .middle } **User Guide**

    ---

    Learn how to configure and use Canvas CLI effectively

    [:octicons-arrow-right-24: User Guide](user-guide/index.md)

-   :material-console:{ .lg .middle } **Command Reference**

    ---

    Complete reference for all available commands

    [:octicons-arrow-right-24: Commands](commands/index.md)

-   :material-school:{ .lg .middle } **Tutorials**

    ---

    Step-by-step guides for common workflows

    [:octicons-arrow-right-24: Tutorials](tutorials/index.md)

-   :material-robot:{ .lg .middle } **AI Agent Skill**

    ---

    Teach Claude Code, Cursor, and other AI agents to drive Canvas CLI

    [:octicons-arrow-right-24: Agent Skill](user-guide/agent-skill.md)

</div>

## Support

- **Issues**: [GitHub Issues](https://github.com/jjuanrivvera/canvas-cli/issues)
- **Discussions**: [GitHub Discussions](https://github.com/jjuanrivvera/canvas-cli/discussions)

## License

Canvas CLI is released under the [MIT License](https://opensource.org/licenses/MIT).
