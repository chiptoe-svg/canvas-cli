<p align="center">
  <img src="docs/assets/images/logo.svg" alt="Canvas CLI" width="280">
</p>

<p align="center">
  <a href="https://github.com/chiptoe-svg/canvas-cli/actions/workflows/ci.yml"><img src="https://github.com/chiptoe-svg/canvas-cli/actions/workflows/ci.yml/badge.svg" alt="CI"></a>
  <a href="https://github.com/chiptoe-svg/canvas-cli/actions/workflows/ci.yml"><img src="https://img.shields.io/badge/coverage-%E2%89%A580%25-brightgreen" alt="Coverage"></a>
  <a href="https://go.dev/"><img src="https://img.shields.io/badge/Go-1.25+-00ADD8?style=flat&logo=go" alt="Go Version"></a>
  <a href="LICENSE"><img src="https://img.shields.io/badge/License-MIT-blue.svg" alt="License"></a>
  <a href="https://github.com/chiptoe-svg/canvas-cli/releases"><img src="https://img.shields.io/github/v/release/jjuanrivvera/canvas-cli" alt="Release"></a>
  <a href="https://pkg.go.dev/github.com/chiptoe-svg/canvas-cli"><img src="https://pkg.go.dev/badge/github.com/chiptoe-svg/canvas-cli.svg" alt="Go Reference"></a>
  <a href="https://deepwiki.com/jjuanrivvera/canvas-cli"><img src="https://deepwiki.com/badge.svg" alt="Ask DeepWiki"></a>
  <a href="https://cliwright.jjuanrivvera.com"><img src="https://img.shields.io/badge/built_with-cliwright-1f6feb" alt="Built with cliwright"></a>
</p>

<p align="center">
  A powerful command-line interface for <a href="https://www.instructure.com/canvas">Canvas LMS</a>, built with Go.
</p>

<p align="center">
  <a href="https://jjuanrivvera.github.io/canvas-cli/"><strong>Documentation</strong></a> ·
  <a href="https://jjuanrivvera.github.io/canvas-cli/getting-started/installation/"><strong>Installation</strong></a> ·
  <a href="https://jjuanrivvera.github.io/canvas-cli/commands/"><strong>Commands</strong></a>
</p>

---

## Features

- **Secure Authentication** - OAuth 2.0 with PKCE (confidential or secret-less public clients), system keyring integration
- **Multi-Instance** - Manage multiple Canvas instances from one CLI
- **Smart Rate Limiting** - Adaptive throttling based on API quotas
- **Multiple Outputs** - Table, JSON, YAML, and CSV formats
- **Interactive Mode** - REPL shell with command history and completion
- **80% API Coverage** - 876 of Canvas's 1086 documented endpoints implemented in the service layer, exposed through 93 command groups
- **Spec-Verified** - Every endpoint path validated against Canvas's official API spec in CI ([details](#api-spec-compliance))
- **MCP Server** - Use as an AI agent tool via Model Context Protocol
- **AI Agent Skill** - Bundled skill for Claude Code, Cursor, and other agents (`canvas skills install`)
- **Signed Releases** - cosign-signed checksums, SBOMs, and a distroless Docker image

## Installation

### Quick Install (macOS/Linux)

Auto-detects your platform and verifies the checksum:

```bash
curl -fsSL https://raw.githubusercontent.com/jjuanrivvera/canvas-cli/main/install.sh | sh
```

### Homebrew (macOS/Linux)

```bash
brew tap jjuanrivvera/canvas-cli
brew install canvas-cli
```

### Scoop (Windows)

```powershell
scoop install https://raw.githubusercontent.com/jjuanrivvera/scoop-canvas-cli/main/canvas-cli.json
```

### Go Install

Requires Go 1.25+ (or Go 1.24+ with automatic toolchain download).

```bash
go install github.com/chiptoe-svg/canvas-cli/cmd/canvas@latest
```

### Docker

```bash
docker run --rm ghcr.io/jjuanrivvera/canvas-cli:latest version
```

Images are published to GHCR for every release (`:latest`, `:v1`, and exact
versions). Pass credentials via environment variables:

```bash
docker run --rm -e CANVAS_URL -e CANVAS_TOKEN ghcr.io/jjuanrivvera/canvas-cli:latest courses list
```

### Binary Download

Download from [GitHub Releases](https://github.com/chiptoe-svg/canvas-cli/releases).
Release checksums are signed with [cosign](https://github.com/sigstore/cosign)
(keyless) and archives ship with SBOMs — verification instructions are in the
[`.goreleaser.yaml`](.goreleaser.yaml) `signs` section.

## Quick Start

```bash
# Authenticate with your Canvas instance
canvas auth login https://your-school.instructure.com

# List your courses
canvas courses list

# Get assignments for a course
canvas assignments list --course-id <course-id>

# Start interactive mode (alias: canvas shell)
canvas repl
```

## Command Overview

| Category | Commands |
|----------|----------|
| **Auth** | `login`, `logout`, `status` |
| **Courses** | `list`, `get`, `create`, `update`, `delete` |
| **Assignments** | `list`, `get`, `create`, `update`, `delete`, `upcoming` (what is due within a window or by a date, per course, local time) |
| **Submissions** | `list`, `get`, `grade`, `bulk-grade`, `comments`, `missing` (who is missing work, read-only), `excuse` (by student/assignment name, with read-back) |
| **Users** | `me`, `list`, `get`, `create`, `update` |
| **Enrollments** | `list`, `get`, `create`, `conclude`, `reactivate`, `accept`, `reject` |
| **Modules** | `list`, `get`, `create`, `update`, `delete`, `publish`, `items` |
| **Pages** | `list`, `get`, `create`, `update`, `delete`, `front`, `revisions` |
| **Discussions** | `list`, `get`, `create`, `entries`, `post`, `reply`, `subscribe` |
| **Announcements** | `list`, `get`, `create`, `update`, `delete` |
| **Quizzes** | `list`, `get`, `create`, `update`, `delete`, `questions`, `submissions` |
| **Schedule** | `schedule` — set available/due/closed times on quizzes and assignments in local time, one by `--id` or in bulk by `--match` |
| **Grades** | `history`, `feed`, `columns` |
| **Groups** | `list`, `get`, `create`, `update`, `delete`, `users`, `categories` |
| **Outcomes** | `list`, `get`, `create`, `update`, `link`, `unlink`, `groups`, `results` |
| **Rubrics** | `list`, `get`, `create`, `update`, `delete`, `associate` |
| **Conversations** | `list`, `get`, `create`, `reply`, `archive`, `star`, `mark-read` |
| **Calendar** | `list`, `get`, `create`, `update`, `delete`, `reserve` |
| **Files** | `list`, `get`, `upload`, `download`, `delete` |
| **Sections** | `list`, `get`, `create`, `update`, `delete`, `crosslist` |
| **Polls** | `polls`, `choices`, `sessions`, `submissions` |
| **Folders / Files** | `folders`, `files`, multi-context (course / group / user) |
| **Favorites / Bookmarks** | `favorites`, `bookmarks`, `course-nicknames` |
| **Account Admin** | `auth-providers`, `csp-settings`, `account-reports`, `enrollment-terms`, `developer-keys`, `grading-period-sets` |
| **Grading** | `grading-periods`, `grading-standards`, `rubric-associations`, `live-assessments` |
| **Content** | `content-exports`, `content-migrations`, `blackout-dates`, `course-pacing`, `blueprint` |
| **Comms / Misc** | `conversations`, `comm-channels`, `content-shares`, `media`, `conferences`, `eportfolios`, `audit` |
| **Admin** | `admins`, `roles`, `analytics`, `sis-imports`, `appointment-groups` |
| **Utilities** | `repl`, `doctor`, `webhook`, `api`, `version` |

This is a sample — the CLI has **93 command groups**, and its service layer
implements **80%** of Canvas's official API. See the
[full command reference](https://jjuanrivvera.github.io/canvas-cli/commands/)
for all 540+ commands and their flags.

## API Spec Compliance

Canvas CLI's API paths are validated in CI against Canvas's **official API
spec** (Swagger 1.2, committed under `testdata/spec/`). A network-free contract
test harvests every endpoint the CLI calls and asserts it matches a documented
Canvas endpoint — so a wrong path fails the build. The committed manifest is
refreshed from a live Canvas host with `make spec-sync`, and `make spec-coverage`
reports which documented endpoints aren't implemented yet.

## Configuration

```yaml
# ~/.canvas-cli/config.yaml
default_instance: myschool
instances:
  myschool:
    url: https://myschool.instructure.com
    client_id: your-client-id
settings:
  default_output_format: table
  cache_enabled: true
  timezone: America/New_York   # zone for local-time flags such as --due "4:50pm" (default: $TZ, then system)
```

See [Authentication Guide](https://jjuanrivvera.github.io/canvas-cli/getting-started/authentication/) for detailed setup.

## Activity Log

An optional local log of what the CLI did: one JSON line per invocation with
the command, its arguments (secrets redacted), every HTTP request with its
status and outcome, the Canvas objects touched, the exit code and the
duration. It is **off by default** and never contains tokens.

```bash
canvas activity configure --enable       # persistently, in ~/.canvas-cli/config.yaml
canvas activity path                     # resolved path and the effective settings
canvas activity list --since 7d          # what ran in the last week
canvas activity list --writes -o json    # every invocation that changed something, in full
canvas activity archive                  # rename to activity-<UTC timestamp>.jsonl
canvas activity clear --force            # truncate (refuses without --force)
```

`configure` writes the `activity_log` block of the config file; the same keys
can be edited by hand:

```yaml
activity_log:
  enabled: true
  path: ~/.canvas-cli/activity.jsonl   # optional (default)
  writes_only: true                    # optional, default true
  capture_bodies: false                # optional
  required: false                      # optional
  max_size_mb: 10                      # optional: archive before exceeding
```

Environment variables take precedence over the file: `CANVAS_ACTIVITY_LOG=<path>`
enables the log and sets its path; `CANVAS_ACTIVITY_WRITES_ONLY`,
`CANVAS_ACTIVITY_CAPTURE_BODIES` and `CANVAS_ACTIVITY_REQUIRED` (`1`/`true` or
`0`/`false`) override the matching keys.

- **The log records writes; reads and dry-runs leave no entry unless
  `writes_only` is turned off** (`configure --writes-only=false` or
  `CANVAS_ACTIVITY_WRITES_ONLY=0`). Each entry carries the command line with
  secrets and free-text values redacted — the value of `--comment`,
  `--rubric-comment`, `--text`, `--message`, `--body` and `--student` is
  replaced by `[REDACTED]` — plus the method, path and outcome of every
  request. Ids, scores and switches are kept, so the log says who was
  graded and what score was posted, not what the feedback said. A write whose response was lost (timeout,
  connection dropped) is always logged, with `"outcome": "unknown"` on the
  request and `"verification_required": true` on the entry, because Canvas
  may have applied it; other writes carry `"outcome": "ok"` or `"rejected"`.
- **`capture_bodies`** adds to every write the exact payload sent (`body`)
  and the parsed response (`response`) — for example the full text of a
  comment, a message or an announcement the agent wrote. Reads are never
  captured. Secret-looking keys (token, secret, password, authorization,
  access_code, api_key — at any depth) and Canvas-shaped tokens or Bearer
  credentials inside strings are redacted. File- and stdin-driven commands
  (`submissions bulk-grade`, `assignments|users create|update --json-file`)
  also record their parsed input under `details.input`, and `quizzes
  regrade` its per-student verification table under `details`; with
  `capture_bodies` off, `details` is not written at all. This mode stores
  student-directed text and names; it is meant for operator audit logs, so
  keep the file where only the operator can read it.
- **`required`** is audited mode: before the first write of an invocation the
  log is prepared (directory created 0700, file created 0600, owned by you,
  writable), and if that fails the write is refused with
  `audit log unavailable: <reason>; refusing to write to Canvas` and a
  non-zero exit. Reads are never blocked. `required` implies `enabled`.
  Without it, a log that cannot be written only produces a warning.
- The log directory is created owner-only (0700) and the file 0600. On each
  run that logs, looser permissions on either — when you own them — are
  tightened to that and reported once on stderr; they are never loosened.
  Put the log in a directory of its own: the tightening applies to the
  parent directory of the file.
- When the file grows past `max_size_mb` it is archived automatically before
  the next write. A failure to write the log never fails the command (unless
  `required` is on).

## MCP Server Mode

Canvas CLI can also run as an [MCP](https://modelcontextprotocol.io/) server, exposing nearly all of its commands as tools for AI coding agents (Claude Code, Cursor, VS Code Copilot). Only the `canvas mcp` management commands themselves are excluded.

```bash
# Start as STDIO MCP server
canvas mcp start

# Start as HTTP MCP server
canvas mcp stream --port 8080

# Export all tool schemas to JSON
canvas mcp tools

# Auto-configure in your editor
canvas mcp claude enable
canvas mcp vscode enable
canvas mcp cursor enable
```

The same binary, two interfaces. When used as an MCP server, each CLI command becomes an MCP tool with typed parameters derived from the command's flags. Required flags become required schema properties. All output goes through structured JSON.

Sensitive flags (`--show-token`, `--config`) are automatically excluded from MCP exposure.

> **Note for `go install` users**: MCP support requires Go 1.25+ due to the [MCP Go SDK](https://github.com/modelcontextprotocol/go-sdk) dependency. Homebrew and binary downloads are not affected by this requirement.

For full setup (Claude Desktop, Claude Code CLI, Cursor, VS Code, OpenCode, Codex), auth precedence, and troubleshooting, see:

- [MCP Integration Guide](https://jjuanrivvera.github.io/canvas-cli/user-guide/mcp/)

## AI Agent Skill

Canvas CLI ships an **agent skill** that teaches AI coding agents (Claude Code,
Cursor, Codex, Gemini CLI, Windsurf, Copilot, …) how to drive it — commands,
flags, safety rules (`--dry-run` previews), and common grading/content
workflows. The skill is bundled in the binary, so installing it from the CLI
guarantees it matches the version you are running:

```bash
canvas skills install --global              # write the bundled skill
canvas skills install --agent cursor        # target a specific agent
```

The skill wraps this binary, so install the CLI (above) and authenticate first.
For structured tool access (Claude Desktop, etc.) use the MCP server mode
described above — the two can coexist.

See the [Agent Skill guide](https://jjuanrivvera.github.io/canvas-cli/user-guide/agent-skill/) for details.

## Agent Safety

When an AI agent drives `canvas`, it can issue destructive commands (`courses delete`, `sections crosslist`, `enrollments conclude`). Use the built-in guard to generate permission rules and a PreToolUse hook that hard-block irreversible operations and require approval for writes:

```bash
# Review what will be generated
canvas agent guard --host claude-code

# Install config and hook script into the current project
canvas agent guard --host claude-code --write

# Strictest: also block create/update/grade
canvas agent guard --host claude-code --all-writes --write
```

Supported hosts: `claude-code`, `codex`, `opencode`. Regenerate after upgrading `canvas` so new commands are covered.

See the [Agent Safety guide](https://jjuanrivvera.github.io/canvas-cli/user-guide/agent-safety/) for the hard-block vs ask model, obfuscation caveats, and MCP-only recommendations.

## Contributing

We welcome contributions! See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

[MIT License](LICENSE)
