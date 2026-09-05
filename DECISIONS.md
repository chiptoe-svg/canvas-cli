# Decisions

Numbered record of ambiguous design/API assumptions and non-obvious engineering
choices for canvas-cli. Each entry captures the decision, the alternatives, and
the rationale so the reasoning survives beyond the diff.

## 1. Diagnostics runner injection seam for deterministic `canvas doctor` tests

**Context.** `canvas doctor` runs environment, configuration, connectivity,
authentication, API-access, disk-space, and permission checks. Several of these
touch live network, the host filesystem, and real credentials, so the command
tests (`TestDoctorCmd`) were flaky in local runs and were skipped entirely under
`CI=1` / `GITHUB_ACTIONS` — a crutch that traded determinism for coverage.

**Decision.** Introduce a `diagnostics.Runner` interface
(`Run(ctx) (*Report, error)`, already satisfied by `*diagnostics.Doctor`) and
route the command through an injectable package var
`newDoctorRunner(cfg, client) diagnostics.Runner` in `commands/doctor.go`
instead of calling `diagnostics.New` directly. The production default is
unchanged. Command tests override the seam with a `fakeRunner` that returns a
deterministic `*diagnostics.Report`, exercising every output format (human,
human-verbose, `--json`, and global `-o json`) plus failure and runner-error
handling — with no network, host-permission, or real-credential dependency.

**Alternatives considered.**
- *Inject each check's dependency (HTTP client, filesystem) into `Doctor`.*
  Rejected as heavier than needed: the diagnostics package tests
  (`internal/diagnostics`, ~97% covered) already drive individual checks against
  `httptest` and temp dirs deterministically. The flakiness lived at the
  *command* layer, so the seam belongs there.
- *Keep the `CI=1` skip.* Rejected — it leaves the command untested in CI and
  still flaky locally, which is exactly the bug (#28).

**Consequences.**
- `go test ./commands -run TestDoctorCmd` is deterministic and host-free by
  default (≈0.02s vs ≈6s), and no longer skips under CI.
- The real end-to-end path stays covered by a single opt-in smoke test,
  `TestDoctorCmd_Live`, gated behind `CANVAS_DOCTOR_LIVE=1` and skipped by
  default.
- Runtime behavior of `canvas doctor` for real users is unchanged; this is a
  testability refactor only.

Ref: [#28](https://github.com/chiptoe-svg/canvas-cli/issues/28)

## 2. Faculty edition: delete, not build-tag

**Context.** This fork narrows Canvas CLI to what an instructor managing and
grading their own courses can do with a standard teacher token. Roughly half of
the upstream command groups — account administration, provisioning, SIS import,
developer keys, cross-instance sync, and the personal/consumer surfaces — have
no place in that tool. The obvious cheap option was to keep the code and gate it
behind a build tag or a hidden-command flag, preserving the ability to rebuild
the full CLI from one tree.

**Decision.** Delete it. The command files, their option structs, their tests,
and every `internal/api` service left with no command caller are removed from
the tree, and `commands/surface_test.go` asserts the exact remaining list of
top-level commands as a test. Design record:
[`docs/superpowers/specs/2026-09-05-faculty-edition-design.md`](docs/superpowers/specs/2026-09-05-faculty-edition-design.md).

**Alternatives considered.**
- *Build tags (`//go:build !faculty`).* Rejected. The claim this edition makes
  to the people installing it is "the binary cannot do admin things" — and a
  build tag makes that claim contingent on a build flag nobody re-checks. It
  also doubles the test matrix, leaves dead code that lint and coverage still
  have to carry, and keeps the reproducible-build story ambiguous: two binaries
  from one tag is one binary too many.
- *Hidden commands (`cmd.Hidden = true`).* Rejected outright — hidden is not
  removed. The endpoints stay reachable by anyone who types the name, so it
  changes the documentation and nothing else.
- *A separate wrapper repository.* Rejected: it would have to track upstream
  forever to get security fixes, and the audited-release story needs one tree
  whose contents are what was reviewed.

**Consequences.**
- The scope is reviewable: `facultySurface` in `commands/surface_test.go` is the
  whole list, and adding a command means editing it in the same commit.
- Restoring a removed group means reverting a specific commit on this branch;
  the history has it, the binary does not.
- Merging upstream changes is manual rather than automatic. Accepted — the
  audited-release model requires reading the diff anyway.
- The spec contract test still covers every path the surviving services call,
  so the trim did not weaken the API-path guard.
