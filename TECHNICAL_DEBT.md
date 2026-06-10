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
   - **Status:** In Progress — a concurrent work stream is migrating these to
     the options-struct pattern.
   - **Files/Areas:** `commands/root.go`, all command files reading the globals
   - **Next Steps:** Finish migration to `commands/internal/options`.
   - **Priority:** Important

2. **Package-Level Flag Variables in Remaining Commands**
   - **Problem:** `api`, `cache`, `sync`, `telemetry`, `repl`, `shell`, and
     `completion` commands still use package-level command/flag variables
     instead of the options-struct pattern documented in AGENTS.md.
   - **Impact:** Same testability/state issues as item 1; the documented
     pattern is not consistently applied.
   - **Status:** In Progress (same migration work stream)
   - **Files/Areas:** `commands/api.go`, `commands/cache.go`,
     `commands/sync.go`, `commands/telemetry.go`, `commands/repl.go`,
     `commands/shell.go`, `commands/completion.go`
   - **Priority:** Important

3. **Gosec Findings Backlog (CI is report-only)**
   - **Problem:** gosec reports ~300 findings (mostly G104 unchecked errors,
     e.g. `cmd.MarkFlagRequired(...)` return values). The CI gosec step runs
     with `-no-fail` so it cannot gate merges.
   - **Impact:** New security findings land unnoticed; the scanner provides no
     enforcement.
   - **Status:** Identified
   - **Files/Areas:** `.github/workflows/ci.yml` (security job), `commands/`
   - **Next Steps:** Burn down findings (handle or `#nosec`-annotate with
     justification), then remove `-no-fail`.
   - **Priority:** Important

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
   - **Status:** Identified
   - **Next Steps:** Add a shared `newTestClient(t, server)` helper and migrate
     call sites incrementally.
   - **Files/Areas:** `internal/api/*_test.go`
   - **Priority:** Low

6. **No Binary-Level Integration Tests**
   - **Problem:** There are no end-to-end tests that exercise the compiled
     `canvas` binary (no `test/integration` suite exists).
   - **Impact:** Flag parsing, alias expansion, exit codes, and output routing
     are only covered indirectly through package tests.
   - **Status:** Identified
   - **Next Steps:** Add a small suite that builds the binary and runs it
     against a mock Canvas server.
   - **Priority:** Low

7. **No Golden-File Formatter Tests**
   - **Problem:** `internal/output` formatters (table/JSON/YAML/CSV) are tested
     with inline assertions, not golden files.
   - **Impact:** Output regressions (column order, truncation, spacing) are
     easy to miss and tedious to assert.
   - **Status:** Identified
   - **Files/Areas:** `internal/output/`
   - **Priority:** Low

8. **Benchmark Test Suite**
   - **Problem:** No automated performance regression detection.
   - **Impact:** Performance changes not caught until production.
   - **Status:** Planned
   - **Next Steps:** Add benchmark tests for critical paths (GetAllPages, rate
     limiter, cache).
   - **Priority:** Low

9. **Additional Platform Coverage in Auth Tests**
    - **Problem:** Platform-specific auth code (macOS ioreg, Windows
      PowerShell) has limited coverage on Linux CI.
    - **Impact:** Some platform-specific code paths only exercised by the
      macOS/Windows legs of the CI matrix.
    - **Status:** Identified
    - **Priority:** Low

---

## Resolved Debt

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
- **Misleading `commands/testing` framework removed** (June 2026) — the package
  never wired `getAPIClient()` to its mock server, making it a trap for
  contributors. Deleted in favour of `commands/internal/testing`.

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
