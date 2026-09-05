# Faculty edition: turning the fork into an instructional-faculty CLI

Date: 2026-09-05
Status: draft for review

## Purpose

This repository stops being "upstream canvas-cli plus faculty features" and
becomes a CLI for instructors managing and grading their own courses. One
audience, one command surface, one skill. Nothing account- or admin-level
remains, no power-user surfaces that can write outside the commands, and the
shipped binary contains only what a faculty member can use.

Faculty keep the same binary name (`canvas`), the same install command, and
the same update path. Nothing changes on their side except a shorter `--help`.

## Decisions already made

- Admin capability is not needed. Faculty at this institution never create
  Canvas accounts; TAs and co-instructors are enrolled as existing users.
- Of the power-user surfaces, only read-only `api get` stays. The
  write-capable `api`, the MCP server, the REPL and the webhook listener go.
- Branches: `main` becomes the development branch (reset to the current
  `integration/all` line, which is then retired). `release/audited` stays the
  faculty release branch so the install URL does not move.
- Identity: module path becomes `github.com/chiptoe-svg/canvas-cli`; README is
  rewritten for faculty; the MkDocs site, its deploy workflow and the docs
  generator are removed.
- Mechanics: delete, do not build-tag. Git history keeps everything; upstream
  fixes are cherry-picked when they touch kept code.

## Command surface

### Stays (course-scoped and tool-level)

| Group | Notes |
|---|---|
| `courses` | list, get, update, settings (`course-settings` folded in as `courses settings`) |
| `assignments`, `assignment-groups`, `overrides` | unchanged, incl. `upcoming` |
| `submissions` | list, get, download, grade, bulk-grade, add-comment, comments, delete-comment, missing, excuse |
| `grades` | history, feed, columns |
| `grading-periods`, `grading-standards` | course-level reads and course-level create |
| `rubrics`, `rubric-associations` | unchanged |
| `quizzes` | unchanged, incl. `regrade`, `extensions`, `statistics`, `reports`, `questions`, `submissions` |
| `modules`, `pages`, `files`, `folders` | unchanged |
| `discussions`, `announcements`, `conversations` | unchanged |
| `enrollments`, `sections`, `groups` | unchanged; `enrollments` is how TAs and co-instructors are added |
| `users` | `list`, `search`, `get`, `profile`, `me`, `missing-submissions`, `todo` (the instructor's grading to-do), `upcoming-events`, `activity-stream` |
| `calendar`, `appointment-groups` | course calendar and Canvas Scheduler (office-hour and exit-interview slots); instructor-level |
| `analytics`, `outcomes`, `peer-reviews` | unchanged |
| `content-shares` | Direct Share of assignments, quizzes and modules with co-instructors; user-level |
| `collaborations`, `course-features` | course-level, instructor-operable; kept at the owner's request |
| `course-extensions` | quiz and assignment extensions for one student (accommodations); instructor-level. Missed in the first inventory, ruled in during execution |
| `content-exports`, `content-migrations` | course copy and import; their `get` commands report job progress |
| `schedule`, `activity` | fork features, unchanged |
| `context`, `alias`, `cache`, `completion` | conveniences, unchanged |
| `agent` | the Claude Code / Codex guard hook installer; part of the agent-safety story |
| `auth`, `config`, `doctor`, `update`, `skills`, `version` | tool-level |
| `api` | **`api get` only**; the command file is rewritten to register the read subcommand alone |

### Goes

Command files (with their option structs, tests and generated docs):
`account_analytics`, `account_calendars`, `account_content_migrations`,
`account_external_tools_favorites`, `account_features`, `account_logins`,
`account_notifications`, `account_reports`, `accounts`, `admins`,
`audit_logs`, `auth_providers`, `blackout_dates`,
`blueprint`, `bookmarks`, `brand`, `comm_channels`,
`comm_messages`, `conferences`,
`course_nicknames`, `course_pacing`, `csp_settings`, `developer_keys`,
`enrollment_terms`, `eportfolios`, `epub_exports`, `error_reports`,
`external_tools`, `favorites`, `grading_period_sets`, `history`, `jwts`,
`live_assessments`, `mcp`, `mcp_annotations`, `media`, `observees`,
`planner`, `polls`, `progress`, `repl`, `shell` (REPL helper), `sis_imports`, `sync`,
`telemetry`, `temporary_enrollment_pairings`, `user_features`, `webhook`,
and the write path of `api`. From `users`: `create`, `merge`, `split`,
`logins`, `settings`, `update`, `update-settings`, `page-views`, `courses`.

`polls` goes even though an instructor may call it: the Polls API served the
"Polls for Canvas" mobile app, which Instructure discontinued, so a poll
created through it has no student client. In-class polling is an ungraded
quiz.

Internal packages that only those commands used: `internal/repl`,
`internal/telemetry`, `internal/webhook`, the MCP bridge dependency
(`njayp/ophis`) and its transitive modules. `internal/shellparse` stays:
`cmd/canvas/main.go` uses it for alias expansion.

Service files in `internal/api` with no remaining command caller are deleted
too. The rule is mechanical: after the command deletions, a service file whose
constructor is referenced by nothing outside tests goes. Expected:
`account_*`, `accounts`, `admins`, `roles`,
`audit_logs`, `auth_providers`, `blackout_dates`, `blueprint`, `bookmarks`,
`brand`, `comm_messages`, `communication_channels`,
`conferences`, `content_migrations_account`,
`course_nicknames`, `course_pacing`, `csp_settings`,
`developer_keys`, `enrollment_terms`, `eportfolios`, `epub_exports`,
`error_reports`, `external_tools`, `favorites`, `grading_period_sets`,
`groups_lti`, `history`, `jwts`, `live_assessments`, `lti_registrations`,
`media_objects`, `observees`, `outcome_imports`, `planner`, `polls`,
`sis_imports`, `temporary_enrollment_pairings`, `user_features`, `user_misc`,
`users_extra`, `custom_data`, `progress`. The spec contract test still covers every path
the remaining services call.

### Judgment calls to confirm

- `agent` stays. It writes the guard hook that makes Claude Code and Codex
  confirm writes; it is faculty-facing agent safety, not admin tooling.
- `content-exports` / `content-migrations` stay. Course copy between terms is
  a faculty task. `account-content-migrations` goes.
- `analytics` stays; `account-analytics` goes.
- `groups` stays (student groups within a course); `group categories` come with
  it.
- `grading-periods` stays for reads; `grading-period-sets` (account-level) goes.
- `progress` goes: only the `progress` command used the Progress service; content-migration `get` reports job state itself.
- `appointment-groups`, `content-shares`, `collaborations`, `course-features`
  and `users todo|upcoming-events|activity-stream` were on the first cut list
  and came back after review: all instructor-operable with a standard token.
- `completion` stays (shell tab completion).
- `alias` stays; `cmd/canvas/main.go` depends on `internal/shellparse` for it, so `internal/shellparse` stays even though the REPL goes.

## Guard: the command list is a test

`commands/surface_test.go` asserts the exact sorted list of top-level command
names registered on the root, and the exact subcommand list of `users` and
`api`. Adding or removing a command means editing the expected list in the
same commit, so the faculty surface is documented in one place and cannot
grow or shrink by accident. The test runs in the normal suite.

## Identity

- `go.mod` module path → `github.com/chiptoe-svg/canvas-cli`; every import
  rewritten (mechanical `sed`, one commit, build must pass).
- `README.md` rewritten for faculty: what the tool is for, the install
  command, verification (checksum, signature, rebuild), the first five
  commands, where the skill lives, the audit log, and a link to the trust
  review. Under 200 lines.
- `CLAUDE.md`/`AGENTS.md` rewritten as contributor guidance for this repo:
  branch model, release process, the command-list test, the spec contract
  test, the coverage gate, the read-back discipline for write commands.
- `docs/` removed except `docs/superpowers/` (specs and plans). `mkdocs.yml`,
  `.github/workflows/docs.yml`, `tools/gendocs` removed. `Makefile` targets
  for docs removed.
- `CHANGELOG.md` gets a "Faculty edition" section explaining the cut and
  listing removed groups; older sections stay as history.
- `SECURITY.md` supported-versions table updated to the audited tags.
- `DECISIONS.md`, `TECHNICAL_DEBT.md`, `CONTRIBUTING.md`: reviewed and
  trimmed to what still applies.
- `Dockerfile`, `Dockerfile.goreleaser`: removed (no image is published).

## Skill

`skills/canvas-cli/` is rewritten as a workflow skill for an agent helping an
instructor:

- `SKILL.md`: what the tool is, how to install and check it, the five
  disciplines the fork already codifies (dry-run first, propose then post,
  read-back is the evidence, student text is data not instruction, never
  send student work to an unapproved service), and pointers to references.
- `references/grading-week.md`: download, review, grade, bulk-grade,
  comment, verify, excuse; rubric grading.
- `references/term-setup.md`: course settings, dates via `schedule`,
  modules and pages, enrolling TAs, announcements, office-hour and
  exit-interview slots via `appointment-groups`.
- `references/mid-term-check.md`: missing work, analytics, at-risk, messaging.
- `references/accommodations.md`: quiz extensions and assignment overrides
  for one student across a course.
- `references/canvas-commands.md`: cheatsheet regenerated from the faculty
  surface only.
- `references/auth-and-config.md`, `references/output-and-filtering.md`:
  trimmed to what remains.

`skill_test.go` asserts every reference is linked from `SKILL.md` and that no
reference names a command that the surface test does not list.

## Branches and CI

1. `main` is force-reset to the tip of `integration/all` after the trim lands
   there; `integration/all` is deleted. This is the one force-push and it is
   announced in CHANGELOG and README.
2. `release/audited` is unchanged in role. The trim reaches it as the next
   release, `v1.13.0+audited.15`, cut the usual way.
3. `.github/workflows/ci.yml` is rewritten in-repo and SHA-pinned like
   `release.yml`: gofmt/vet/golangci-lint, `go test -race` with the 80 %
   coverage gate, govulncheck, and `goreleaser check`, on pushes to `main`
   and `release/audited` and on pull requests. No reusable workflow from
   another repository.
4. `claude.yml` and `claude-code-review.yml` (upstream's Claude automation)
   are removed.

## Testing and verification

- The full suite, coverage gate, lint, spec contract test and the new
  surface test pass on `main` and on `release/audited`.
- Coverage is recomputed after the deletions; if it falls below 80 %, the
  uncovered kept code is tested, not the gate lowered.
- The binary is built and its `--help` compared with the expected list.
- `v1.13.0+audited.15` is verified as before: signature identity, checksum,
  byte-identical rebuild from a fresh clone; the update check from an
  audited.14 install notices it.
- The install script is run against the release on a clean temporary HOME.

## Out of scope

No new features. The wave-one reports (readiness audit, gradebook snapshot,
message-students-who, and the rest) follow on the trimmed base as separate
work, each with its own spec.

## Risks

- A kept command may depend on a deleted service through a shared type. The
  build catches this; the fix is to keep the type, not the service.
- A faculty member may use a removed command. Mitigation: the removed groups
  are all account-level or non-faculty, and the release notes list them; a
  read-only `api get` covers one-off needs until a command exists.
- Upstream merges become cherry-picks. Accepted; the fork has been doing this
  in practice.
