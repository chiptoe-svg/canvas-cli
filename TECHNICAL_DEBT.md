# Technical Debt Tracking

This document tracks known technical debt in the Canvas CLI project.

**Last Updated:** 2026-06-10
**Status:** Re-triaged after the v1.9.0 release. Every item here is either
actively planned with a target, deliberately dormant, or an accepted design
choice — items that would sit at "Planned" forever have been removed (see
"Removed Items" for why, so they don't get re-added by a future review).

---

## Active Items

1. **Gosec Findings Backlog (CI is report-only)**
   - **Problem:** gosec reports ~300 findings (mostly G104 unchecked errors,
     e.g. `cmd.MarkFlagRequired(...)` return values). The CI gosec step runs
     with `-no-fail` so it cannot gate merges — the scan is decorative.
   - **Impact:** New, genuine security findings would land unnoticed. This is
     the only item with permanent, ongoing payoff: every future PR gets a real
     gate.
   - **Next Steps:** (1) Add a `mustMarkRequired(cmd, name)` helper (panic on
     programmer error) and sweep the G104 bulk. (2) Triage the remainder: fix
     or `#nosec` with written justification. (3) Remove `-no-fail`.
   - **Effort:** ~half a day.
   - **Priority:** Important — **Target: v1.10**

2. **No Binary-Level Integration Tests**
   - **Problem:** Nothing tests the compiled `canvas` binary. Package tests
     cannot see binary-surface regressions — the June review's findings that
     weren't code bugs (documented exit codes that didn't exist, a quickstart
     command that didn't parse, duplicate REPL commands) were all of this kind.
   - **Impact:** For a CLI, the binary is the product; it currently ships
     untested as a whole.
   - **Next Steps:** A deliberately small suite in `test/integration` (build
     tag `integration`): compile once, drive with `os/exec` against an
     `httptest` mock Canvas. Cap at ~10–15 cases: env-var auth, `-o json`,
     exit codes 0/1, alias expansion, `context set`, `--dry-run` redaction.
     Resist letting it grow into a parallel test universe.
   - **Effort:** ~a day.
   - **Priority:** Important — **Target: v1.10**

3. **Old Command Test Framework Is a Trap**
   - **Problem:** `commands/testing` coexists with the newer
     `commands/internal/testing` — and the old one never routes
     `getAPIClient()` to its mock server, so tests written against it look
     like integration tests but exercise no HTTP dispatch. The next
     contributor (or AI agent) will copy it.
   - **Impact:** Actively misleading; the only item where doing nothing
     misinforms people.
   - **Next Steps:** Migrate its few remaining usages to
     `commands/internal/testing`, delete the package.
   - **Effort:** ~an hour.
   - **Priority:** Important — **Target: v1.10 (do first)**

4. **Duplicated Test Client Construction** *(timeboxed)*
   - **Problem:** ~476 duplicated `NewClient(ClientConfig{...})` blocks across
     `internal/api` tests. Only bites when `ClientConfig` changes shape, which
     is rare — this is tidiness, not compounding debt.
   - **Next Steps:** One timeboxed attempt: add `newTestClient(t, serverURL)`
     and migrate call sites with a scripted regex pass, validated by the full
     suite + `-race`. **If it isn't done in about an hour, abandon it** and
     move this to Removed Items.
   - **Priority:** Low — **Target: opportunistic, alongside the v1.10 batch**

## Dormant (deliberate — do not act early)

5. **Cosign Pinned to the v2 Line**
   - **Problem:** `release.yml` pins `cosign-release: v2.6.3` because cosign
     v3 changed the sign-blob/verify-blob bundle format, and both
     `.goreleaser.yaml` (`signs:`) and the documented verification
     instructions use the v2 `.sig`/`.pem` flags.
   - **Why dormant:** v2 verification remains fully supported (v3 clients
     verify v2 signatures), the pin is commented in the workflow, and
     dependabot keeps surfacing installer updates as a reminder. Migrating
     early buys nothing and risks breaking published verify instructions.
   - **Trigger to act:** cosign v2 EOL announcement, or v3 bundles becoming
     the GoReleaser-documented default. Then migrate `signs:` config and both
     docs pages together, validate on a snapshot release, unpin.
   - **Files/Areas:** `.github/workflows/release.yml`, `.goreleaser.yaml`,
     `docs/getting-started/installation.md`, `docs/security.md`

## Accepted Design Choices (not debt — do not re-add)

- **Root-level persistent-flag globals.** `commands/root.go` declares 13
  package-level globals read ~380 times. This is the standard Cobra pattern;
  a CLI process runs one command, so the shared state causes no production
  bugs — the only cost is that command tests can't run in parallel, and the
  suite takes ~27s anyway. A 380-call-site migration trades theoretical
  cleanliness for real regression risk. **Rule:** when a command file is
  substantially rewritten for other reasons, thread an options struct then;
  never as a standalone project. Reconsider only if a genuinely concurrent
  surface (TUI, server-mode dispatch) is built.
- **Services accept `*Client`, not the `HTTPClient` interface.** The
  `httptest`-based test pattern exercises real HTTP, real Link headers, and
  real error bodies — converting ~30 services to interface fakes would make
  tests faster and weaker. This is a design choice, not debt.

## Removed Items (assessed June 2026 — won't do, and why)

- **Golden-file formatter tests** — the formatter is at ~94% coverage with
  substantive assertions; JSON/CSV correctness (the scriptable contract) is
  already asserted, and table spacing is not a contract. Goldens would mostly
  generate `-update` churn. Revisit only if a real formatting regression ever
  ships.
- **Benchmark suite** — this CLI's performance is dominated by network
  latency and Canvas's rate-limit quota, which the client deliberately
  throttles to; local benchmarks measure nothing a user can feel, and CI
  runner noise makes them un-gateable. Sat at "Planned" since January with no
  constituency. Revisit only if a concrete performance issue is reported.
- **Additional platform coverage in auth tests** — the 3-OS CI matrix *is*
  the coverage mechanism for platform-specific code; there was nothing to do.

---

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
  migration was *complete* was wrong; corrected June 2026.
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
N. **Short Title**
   - **Problem:** What is wrong?
   - **Impact:** Why does this matter?
   - **Next Steps:** What needs to happen next?
   - **Priority:** Important/Low — **Target: <milestone or trigger>**
```

Every item must carry a target milestone or an explicit trigger condition.
An item that would sit at "Planned" indefinitely belongs in Removed Items
with its rationale — a plan that won't be executed is misinformation.

### Priorities

- **Important:** Impacts security posture, correctness guarantees, or
  actively misleads contributors
- **Low:** Beneficial but not urgent; must be timeboxed or opportunistic

---

**Next Review:** September 2026
**Maintained By:** Canvas CLI Development Team
