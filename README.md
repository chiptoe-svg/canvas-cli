# Canvas CLI — faculty edition

A command-line tool for instructors who manage and grade their own
[Canvas LMS](https://www.instructure.com/canvas) courses. It reaches only the
courses you are enrolled to teach or TA and needs no admin rights: every
account-level, provisioning and cross-instance surface has been removed rather
than hidden. Every release is signed with [Sigstore](https://www.sigstore.dev/)
and reproducible — you can rebuild the published binary byte-for-byte from the
tag. Based on
[jjuanrivvera/canvas-cli](https://github.com/jjuanrivvera/canvas-cli) (MIT),
trimmed and re-audited for faculty use.

## Install or update

macOS and Linux, amd64 and arm64:

```bash
curl -fsSL https://raw.githubusercontent.com/chiptoe-svg/canvas-cli/release/audited/install.sh | sh
```

The script downloads the release it is pinned to (a `v1.13.0+audited.N` tag),
checks the archive's SHA-256 against the release's `checksums.txt`, and installs
the binary to `INSTALL_DIR` if you set it, `/usr/local/bin` if that is writable,
otherwise `~/.local/bin`. Running it again is how you update. Confirm what you
got:

```bash
canvas version    # canvas-cli 1.13.0+audited.N
canvas doctor     # install, config, auth and connectivity in one pass
```

### Verify the signature

The release workflow signs `checksums.txt` keylessly, so the signer identity is
this repository's `release.yml` at the tagged ref. Download `checksums.txt`,
`checksums.txt.sig` and `checksums.txt.pem` from the release, then with
[cosign](https://github.com/sigstore/cosign) v2:

```bash
cosign verify-blob --certificate checksums.txt.pem --signature checksums.txt.sig \
  --certificate-oidc-issuer https://token.actions.githubusercontent.com \
  --certificate-identity-regexp '^https://github.com/chiptoe-svg/canvas-cli/\.github/workflows/release\.yml@refs/tags/v' \
  checksums.txt
```

That regexp is the one in [`.goreleaser.yaml`](.goreleaser.yaml). A signature
from any other workflow, repository or ref fails it.

### Reproduce the build

Builds use `-trimpath` and take the build date from the commit, so the tag plus
the Go version in [`go.mod`](go.mod) reproduce the released bytes exactly:

```bash
TAG=v1.13.0+audited.1     # the tag you are checking
git clone --depth 1 --branch "$TAG" https://github.com/chiptoe-svg/canvas-cli.git
cd canvas-cli
CGO_ENABLED=0 GOTOOLCHAIN=go1.25.13 go build -trimpath \
  -ldflags "-s -w -X main.Version=${TAG#v} \
            -X main.Commit=$(git rev-parse HEAD) \
            -X main.BuildDate=$(TZ=UTC0 git log -1 --format=%cd --date=format-local:%Y-%m-%dT%H:%M:%SZ)" \
  ./cmd/canvas
shasum -a 256 canvas
```

Compare it with `canvas` from the release archive (`GOOS`/`GOARCH` for others).

## First five commands

```bash
canvas auth login https://your-school.instructure.com   # OAuth 2.0 + PKCE
canvas courses list                                     # your courses
canvas assignments list --course-id 123                 # what is in one course
canvas submissions missing --course-id 123              # who is missing work
canvas users todo                                       # your grading to-do
```

`canvas <group> <sub> --help` is the authoritative flag reference.

## What it does

| Group | What it covers |
|---|---|
| `courses` | list, get, create, update, delete, and `courses settings` |
| `assignments`, `assignment-groups`, `overrides` | assignments, their groups, per-student and per-section overrides, `upcoming` |
| `submissions` | list, get, download, grade, bulk-grade, comments, `missing`, `excuse` |
| `grades` | grade-change history, the gradebook feed, custom columns |
| `grading-periods`, `grading-standards` | course-level reads; `grading-periods` updates and deletes, `grading-standards` creates and deletes |
| `rubrics`, `rubric-associations` | build rubrics and attach them to assignments |
| `quizzes` | quizzes, questions, submissions, `regrade`, `extensions`, `statistics`, `reports` |
| `course-extensions` | quiz and assignment accommodations for a named student |
| `modules`, `pages`, `files`, `folders` | course content and its publish state |
| `discussions`, `announcements`, `conversations` | course discussion, announcements, messaging students |
| `enrollments`, `sections`, `groups` | roster, sections, student groups; how TAs and co-instructors are added |
| `users` | `me`, `list`, `search`, `get`, `profile`, `todo`, `missing-submissions`, `upcoming-events`, `activity-stream` |
| `calendar`, `appointment-groups` | course calendar and Scheduler slots (office hours, exit interviews) |
| `analytics`, `outcomes`, `peer-reviews` | course analytics, outcomes, peer-review assignment |
| `content-shares`, `collaborations`, `course-features` | Direct Share with co-instructors, collaborations, course feature flags |
| `content-exports`, `content-migrations` | course copy and import; their `get` reports job progress |
| `schedule` | set available, due and close times in local time — one by `--id`, or in bulk by `--match` |
| `activity` | the local activity log (below) |
| `api` | `api get` only: a read-only escape hatch for reads no command covers |
| `agent`, `skills` | the agent guard hook installer and the bundled agent skill |
| `auth`, `config`, `context`, `alias`, `cache`, `completion`, `doctor`, `update`, `version` | tool-level plumbing |

The exact command surface is pinned by a test, `commands/surface_test.go`: it
cannot drift without a reviewed commit.

## Safety

- **`--dry-run` on every command.** It prints the exact HTTP request as a curl
  line with the token redacted, and sends nothing. For a `--match` or a CSV
  batch the dry run *is* the plan.
- **Writes read back and print evidence.** `submissions grade`, `add-comment`
  and `excuse` re-read the object after writing and print what changed —
  `grade: 88 → 95`, the new comment's id and author, `excused: not excused →
  excused` — plus a `verified:` line. `verified: no` exits non-zero. The write's
  own echo is not evidence; the read-back is.
- **`api get` is read-only.** There is no `api post`, `put`, `patch` or
  `delete`. The faculty build cannot write through the raw API at all.
- **Optional local activity log.** Off by default, and it never contains
  tokens. When enabled it writes one JSON line per invocation: the command, its
  arguments with secrets and free-text values redacted, every HTTP request with
  status and outcome, the objects touched, the exit code and the duration.

  ```bash
  canvas activity configure --enable     # persist in ~/.canvas-cli/config.yaml
  canvas activity list --since 7d        # what ran in the last week
  canvas activity list --writes -o json  # every invocation that changed something
  ```

  By default only writes are recorded, and the values of `--comment`,
  `--rubric-comment`, `--text`, `--message`, `--body` and `--student` are
  replaced by `[REDACTED]` — the log says who was graded and what score was
  posted, not what the feedback said.
  `configure --capture-bodies` keeps the payloads too; that file then holds
  student-directed text, so keep it where only you can read it. A write whose
  response was lost is always logged with `"verification_required": true`,
  because Canvas may have applied it. The directory is created 0700 and the file
  0600, and `configure --required` refuses to write to Canvas at all when the
  log cannot be opened. `canvas activity path` prints the effective settings.

## Using it with an AI agent

The binary ships with an agent skill — what the tool is, the disciplines that
apply to grading somebody's education record, and the workflow references.
Installing it from the CLI guarantees it matches the version you are running:

```bash
canvas skills install --global           # write the bundled skill
canvas skills install --agent cursor     # target a specific agent
```

An agent driving `canvas` can still issue destructive commands. The guard
generates permission rules and a PreToolUse hook that hard-block irreversible
operations and require approval for writes:

```bash
canvas agent guard --host claude-code            # review what would be written
canvas agent guard --host claude-code --write    # install into the project
canvas agent guard --host claude-code --all-writes --write   # also gate create/update/grade
```

Hosts: `claude-code`, `codex`, `opencode`. Regenerate after upgrading.

## Trust review

What this fork removed and why is written down, not summarized:

- [`docs/superpowers/specs/2026-09-05-faculty-edition-design.md`](docs/superpowers/specs/2026-09-05-faculty-edition-design.md)
  — the design: every command group kept, every one removed, and the reasoning
  for each judgment call.
- The **Faculty edition** section at the top of [`CHANGELOG.md`](CHANGELOG.md)
  — the same list in release-note form.

The claims above are checkable rather than asserted: the command surface is a
test, the release recipe is a SHA-pinned in-repo workflow
([`.github/workflows/release.yml`](.github/workflows/release.yml)) so the
signature identifies that workflow at the tag, and the build reproduces —
verify all of it with the commands under [Install or update](#install-or-update).

## Development

```bash
make build    # bin/canvas
make test     # the full suite
make check    # go vet, golangci-lint, gofmt, gosec if installed,
              # go test -race, and the integration tests
```

`main` is the development branch. `release/audited` is what faculty install:
`main` plus the commit that pins the installer's version. Releases are tags
`v1.13.0+audited.N` cut from `release/audited`. Contributor guidance is in
[AGENTS.md](AGENTS.md); see also [CONTRIBUTING.md](CONTRIBUTING.md).

## License

[MIT](LICENSE).
