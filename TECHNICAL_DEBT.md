# Technical Debt Tracking

This document tracks known technical debt in the Canvas CLI project.

**Last Updated:** 2026-06-10
**Status:** Updated after the June 2026 deep review. Earlier claims in this file
("0 global flag variables", "all 34 commands migrated") were inaccurate and
have been removed.

---

## Active Technical Debt

### Important

1. **Root-Level Global Flag Variables**
   - **Problem:** `commands/root.go` declares 13 package-level persistent-flag
     globals (`verbose`, `quiet`, `dryRun`, `noCache`, `outputFormat`,
     `filterText`, etc.) read roughly 380 times across the `commands/` package.
   - **Impact:** Commands share mutable global state; hard to test in
     isolation, unsafe for concurrent execution.
   - **Status:** Planned — needs a design decision first (see Next Steps);
     the per-command options migration (former item 2) is complete, this is
     the remaining root-level layer.
   - **Files/Areas:** `commands/root.go`, all command files reading the globals
   - **Next Steps:** (1) Decide the carrier: a `GlobalOptions` struct resolved
     once in `PersistentPreRun` and passed via `cmd.Context()`, vs. reading
     cobra flag values per call. (2) Migrate incrementally, one command file
     per PR, starting with files that have the fewest reads. Do not attempt a
     big-bang migration.
   - **Priority:** Important — **Target: incremental, start v1.11**

2. **Gosec Findings Backlog (CI is report-only)**
   - **Problem:** gosec reports ~300 findings (mostly G104 unchecked errors,
     e.g. `cmd.MarkFlagRequired(...)` return values). The CI gosec step runs
     with `-no-fail` so it cannot gate merges.
   - **Impact:** New security findings land unnoticed; the scanner provides no
     enforcement.
   - **Status:** Planned
   - **Files/Areas:** `.github/workflows/ci.yml` (security job), `commands/`
   - **Next Steps:** (1) Fix the G104 bulk with a small `mustMarkRequired`
     helper (panic on programmer error) — that alone removes most findings.
     (2) Triage the remainder: fix or `#nosec` with justification. (3) Flip
     CI to enforcing by removing `-no-fail`, optionally excluding rules that
     are consciously accepted.
   - **Priority:** Important — **Target: v1.10**

3. **Cosign Pinned to v2 Line**
   - **Problem:** `release.yml` pins `cosign-release: v2.6.3` because cosign
     v3 changed the sign-blob/verify-blob bundle format, and both
     `.goreleaser.yaml` (`signs:`) and the documented verification
     instructions use the v2 `.sig`/`.pem` flags.
   - **Impact:** Signing tooling ages; v2 will eventually stop receiving
     fixes.
   - **Status:** Identified (introduced deliberately, June 2026)
   - **Files/Areas:** `.github/workflows/release.yml`, `.goreleaser.yaml`,
     `docs/getting-started/installation.md`, `docs/security.md`
   - **Next Steps:** Migrate signs config and docs to the v3 bundle format
     together, then unpin. Verify with a snapshot release first.
   - **Priority:** Low — **Target: when cosign v3 bundles are the ecosystem
     default**

### Nice to Have

4. **Services Accept `*Client` Instead of `HTTPClient` Interface**
   - **Problem:** Every service in `internal/api` takes the concrete `*Client`
     (`func NewXxxService(client *Client)`) rather than the `HTTPClient`
     interface.
   - **Impact:** Services cannot be unit-tested with a lightweight fake; tests
     spin up `httptest` servers instead.
   - **Status:** Accepted trade-off — documented here. The `httptest`-based
     test pattern works and exercises real HTTP behavior; switching to the
     interface is a large mechanical change with modest payoff.
   - **Files/Areas:** `internal/api/*.go`
   - **Priority:** Low

5. **Duplicated Test Client Construction**
   - **Problem:** ~476 duplicated `NewClient(ClientConfig{...})` blocks across
     `internal/api` tests.
   - **Impact:** Boilerplate; config changes require mass edits.
   - **Status:** Planned
   - **Next Steps:** Add a shared `newTestClient(t, server)` helper; the
     call-site migration is regex-scriptable in one pass (verify with the full
     suite + `-race`).
   - **Files/Areas:** `internal/api/*_test.go`
   - **Priority:** Low — **Target: v1.10 test-debt batch**

6. **No Binary-Level Integration Tests**
   - **Problem:** There are no end-to-end tests that exercise the compiled
     `canvas` binary (no `test/integration` suite exists).
   - **Impact:** Flag parsing, alias expansion, exit codes, and output routing
     are only covered indirectly through package tests.
   - **Status:** Planned
   - **Next Steps:** Small `test/integration` suite (build tag `integration`):
     compile the binary once per run, exercise it with `os/exec` against an
     `httptest` mock Canvas — cover env-var auth, `-o json|csv`, exit codes
     0/1, alias expansion, `context set`, and `--dry-run` redaction. Wire into
     CI as a separate job.
   - **Priority:** Low — **Target: v1.10**

7. **No Golden-File Formatter Tests**
   - **Problem:** `internal/output` formatters (table/JSON/YAML/CSV) are tested
     with inline assertions, not golden files.
   - **Impact:** Output regressions (column order, truncation, spacing) are
     easy to miss and tedious to assert.
   - **Status:** Planned
   - **Next Steps:** `testdata/` golden files with an `-update` flag; cover one
     representative struct across all four formats plus nil/empty edge cases.
   - **Files/Areas:** `internal/output/`
   - **Priority:** Low — **Target: v1.10 test-debt batch**

8. **Two Command Test Frameworks Coexist**
   - **Problem:** The older `commands/testing` framework coexists with the
     newer `commands/internal/testing` package.
   - **Impact:** Confusing for contributors; duplicate maintenance. The old
     framework also doesn't route `getAPIClient()` to its mock server, so its
     tests don't exercise real HTTP dispatch.
   - **Status:** Planned
   - **Next Steps:** Migrate its few remaining usages to
     `commands/internal/testing`, then delete `commands/testing`.
   - **Priority:** Low — **Target: v1.10 test-debt batch**

9. **Benchmark Test Suite**
   - **Problem:** No automated performance regression detection.
   - **Impact:** Performance changes not caught until production.
   - **Status:** Identified
   - **Next Steps:** Add benchmarks for `GetAllPagesGeneric`, the adaptive
     rate limiter, and cache read/write. Run in CI as informational only
     (benchmarks on shared runners are too noisy to gate on).
   - **Priority:** Low — **Target: opportunistic**

10. **Additional Platform Coverage in Auth Tests**
    - **Problem:** Platform-specific auth code (macOS ioreg, Windows
      PowerShell) has limited coverage on Linux CI.
    - **Impact:** Some platform-specific code paths only exercised by the
      macOS/Windows legs of the CI matrix.
    - **Status:** Identified
    - **Priority:** Low

---

## Plan Summary

| When | What | Items |
|------|------|-------|
| v1.10 — security gate | G104 helper + triage, flip gosec to enforcing | 2 |
| v1.10 — test-debt batch (one PR, test-only, zero user risk) | shared test client helper, golden files, delete old framework | 5, 7, 8 |
| v1.10 | binary-level integration suite | 6 |
| v1.11+ — incremental | root-global flag migration (design first, file-by-file) | 1 |
| Opportunistic | benchmarks | 9 |
| When ecosystem ready | cosign v3 bundle migration | 3 |
| Accepted, no action | services take `*Client` (4), platform auth coverage (10) | 4, 10 |

## Resolved Debt

- **Package-level flag variables in remaining commands** (June 2026, #37) —
  `api`, `cache`, `sync`, `telemetry`, `repl`, `completion` migrated to the
  options-struct pattern; `shell` became an alias of `repl`.
- **Dependabot configured** (June 2026) — `github-actions` + `gomod` weekly;
  SECURITY.md had advertised this before any config existed.
- **GitHub Actions on Node 24** (June 2026) — all node-based actions bumped
  ahead of the June 16 forced migration.
- **Structured logging and options pattern introduced** (January 2026) —
  `commands/internal/options` and `commands/internal/logging` packages exist
  and most resource commands use them. Note: the original claim that the
  migration was *complete* was wrong; see Active items 1-2.
- **Auth module test coverage** raised from 48.9% to 71.7% (January 2026).
- **Configuration validation** added in `internal/config/validation.go`
  (January 2026).
- **GetAllPages generics optimization** replacing reflection (January 2026).
- **Silent error handling audit** completed; the standalone
  ERROR_HANDLING_AUDIT.md was removed in June 2026 as stale.
- **Overall coverage** raised to ~82% with a CI coverage gate (June 2026,
  #31).

---

## Debt Tracking Guidelines

### Adding New Debt Items

When adding technical debt, use this format:

```markdown
### [Priority Level]

N. **Short Title**
   - **Problem:** What is wrong?
   - **Impact:** Why does this matter?
   - **Status:** Current state (Planned/In Progress/Blocked)
   - **Files/Areas:** Where is this located?
   - **Next Steps:** What needs to happen next?
   - **Priority:** Critical/Important/Low
```

### Priorities

- **Critical:** Blocks new features, security issues, or severely impacts maintainability
- **Important:** Impacts code quality, performance, or developer experience
- **Nice to Have:** Improvements that would be beneficial but not urgent

### Status Values

- **Identified:** Problem recognized but not yet planned
- **Planned:** Accepted for future work
- **In Progress:** Actively being worked on
- **Blocked:** Cannot proceed due to dependencies
- **Resolved:** Completed and moved to "Resolved Debt"

---

**Next Review:** September 2026
**Maintained By:** Canvas CLI Development Team
