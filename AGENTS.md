# AGENTS.md

Guidance for anyone — human or AI agent (Claude Code, Cursor, Copilot, …) —
working on this repository. `CLAUDE.md` is a symlink to this file.

This is the **faculty edition** of Canvas CLI: a tool for instructors managing
and grading their own courses. It has no admin capability, and that is a
property to preserve, not an accident to fix. Before adding anything, read
[`docs/superpowers/specs/2026-09-05-faculty-edition-design.md`](docs/superpowers/specs/2026-09-05-faculty-edition-design.md)
— it records every group that was kept, every one that was removed, and why.

## Build & test

```bash
make build              # bin/canvas
make dev                # fmt + vet + build

make test               # full suite
make test-coverage      # suite with a coverage report
make test-integration   # binary-level tests (-tags integration)
go test ./internal/api/...          # one package
go test -run TestName ./...         # one test

make fmt                # gofmt
make lint               # golangci-lint
make vet                # go vet
make check              # go vet, golangci-lint, gofmt, gosec if installed,
                        # go test -race, and the integration tests.
                        # CI adds govulncheck, gosec, the ≥80% coverage
                        # gate, and goreleaser check. Run before you push.

make setup-hooks        # pre-commit hook: gofmt, golangci-lint, go vet, go test -short -race
```

## Layout

```
cmd/canvas/     → entry point (main.go); alias expansion via internal/shellparse
commands/       → Cobra command definitions, one file per resource
  internal/
    options/    → option structs + Validate()
    logging/    → structured command logging
    testing/    → cmdtest helpers (mock server wired to getAPIClient)
internal/
  activity/     → the local activity log (redaction, permissions, audited mode)
  api/          → Canvas API client + service layer (Client, *Service)
  auth/         → OAuth 2.0 + PKCE, token storage (keyring / encrypted file)
  batch/        → concurrent batch operations (worker pool)
  cache/        → response caching with TTL
  config/       → Viper-based configuration
  diagnostics/  → canvas doctor checks
  dryrun/       → --dry-run curl rendering
  localtime/    → local-time parsing for --due / --available / --until
  output/       → table, JSON, YAML, CSV formatters
  progress/     → progress indicators
  resolve/      → name → id resolution (students, assignments)
  shellparse/   → shell-style argument parsing
  terminal/     → terminal capabilities
  update/       → self-update checks
skills/canvas-cli/  → the agent skill bundled into the binary
testdata/spec/  → committed Canvas API spec manifest
tools/speccheck → spec sync / coverage tool
```

## The command surface is a test

`commands/surface_test.go` asserts the exact sorted list of top-level commands,
the `users` subcommands, and that `api` has only `get`. **Adding or removing a
command means editing `facultySurface` in the same commit.** That is deliberate:
the scope of this tool is reviewable in one file and cannot drift in silently.

If a change makes that test fail, the question to answer is not "how do I fix
the test" but "is this command something an instructor can run with their own
token, and does the spec say it belongs here?" If the answer is yes, edit the
list and say so in the commit message.

## Write commands read back and print evidence

Every command that changes a student's record re-reads the object after writing
and prints what actually changed, plus a `verified:` line; `verified: no` exits
non-zero. The write's own response echo is **not** evidence — Canvas can accept
a request and store something different (a rubric row silently dropped, a grade
clamped by a grading standard).

`commands/submissions_readback.go` is the pattern: capture the object before,
write, re-read, diff the requested fields against the read-back, collect
mismatches, and render `before → after` per field. New write commands follow it.
A write command with no read-back is incomplete, not merely untested.

## Spec contract and coverage

`internal/api/spec_contract_test.go` is a network-free test that harvests every
`/api/v1/...` path the service layer calls and asserts each matches a documented
Canvas endpoint in `testdata/spec/canvas_endpoints.json`. **A wrong path fails
the build.** Take exact paths and verbs from the committed manifest, and field
names from `canvas_models.json`; the committed manifest is authoritative.

- `make spec-sync` refreshes the manifest from a live Canvas host
  (`-host`/`CANVAS_SPEC_HOST`, default `learn.canvas.net`).
- `make spec-coverage` lists documented-but-unimplemented endpoints. Treat it as
  information, not a target: this edition implements what faculty need.

CI enforces a **total coverage gate of ≥80%** (`go tool cover -func` over
`./...`) on the ubuntu job. Every new command needs cmdtest coverage — the run
function *and* the option struct's `Validate()` — or the gate drops. `commands`
and `commands/internal/options` are where coverage erodes fastest.

## Patterns

**Options struct, not package globals.** Define the struct in
`commands/internal/options/`, give it `Validate()`, bind flags to its fields in
the command constructor, and pass it to a `run…(ctx, client, opts)` function.
The persistent-flag globals in `commands/root.go` are a documented exception
(see [TECHNICAL_DEBT.md](TECHNICAL_DEBT.md)); do not add more.

**Structured logging.** `logging.NewCommandLogger(globalDebugFlag)`, then
`LogCommandStart(ctx, "resource.list", fields)` and
`LogCommandComplete(ctx, "resource.list", n)`.

**Service layer.** One service per Canvas resource in `internal/api/`:
`type ModulesService struct { client *Client }` with
`func NewModulesService(client *Client) *ModulesService`. The client handles
pagination (`GetAllPages`), adaptive rate limiting from Canvas quota headers,
and exponential-backoff retry.

**Tests** use `httptest.NewServer` mock servers against a real `*Client`. Use
`t.Fatal` (not `t.Error`) where a nil check would let a later line panic.

## Branches and releases

- **`main`** — the development branch. Work happens here (or on branches merged
  into it).
- **`release/audited`** — what faculty install: `main` plus the one commit that
  pins `install.sh` to the release being cut. Nothing else diverges.

Releases are tags `v1.13.0+audited.N` on `release/audited`. To cut one:

```bash
# on release/audited, merged up to main
vi install.sh                                  # VERSION="${CANVAS_VERSION:-v1.13.0+audited.N}"
git commit -am "release: pin installer to v1.13.0+audited.N"
git tag -a "v1.13.0+audited.N" -m "v1.13.0+audited.N"
git push origin release/audited "v1.13.0+audited.N"
```

`.github/workflows/release.yml` builds, signs and publishes. It is deliberately
self-contained and every action is pinned by commit SHA: that is what makes the
Sigstore signature meaningful, since the signer identity on `checksums.txt` is
this file in this repository at the tagged ref. Changing the recipe requires a
reviewed commit. After the release lands, **verify it as a user would** — the
cosign command and the byte-identical rebuild in the [README](README.md).

Before tagging: add the version's section to `CHANGELOG.md`, and update the
supported-versions table in `SECURITY.md`.

## Install and update instructions point at one place

Anything that tells a user how to get or update the binary — the README, the
bundled skill under `skills/canvas-cli/`, release notes, error messages — names
exactly one install path:

```bash
curl -fsSL https://raw.githubusercontent.com/chiptoe-svg/canvas-cli/release/audited/install.sh | sh
```

Never a Homebrew tap, never a package manager, never a build from source with
the Go toolchain, never a direct binary download without checksum verification. Those fetch builds that
are not the audited one, and the point of this edition is that the binary a
faculty member runs is the binary that was reviewed. The same rule applies to
docs an agent might read: if a file tells someone how to install `canvas`, it
tells them this.
