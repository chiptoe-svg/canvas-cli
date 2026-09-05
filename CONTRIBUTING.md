# Contributing to Canvas CLI (faculty edition)

The engineering guidance — layout, patterns, the tests that guard the command
surface and the API paths, the coverage gate, and the release procedure — lives
in **[AGENTS.md](AGENTS.md)** (also reachable as `CLAUDE.md`). Read that first;
this file covers only the process around it.

## Setup

Go 1.25 (matching `go.mod`), git, and make.

```bash
git clone https://github.com/chiptoe-svg/canvas-cli.git
cd canvas-cli
make deps
make build
make setup-hooks    # gofmt, golangci-lint, go vet, go test -short -race on commit
```

## Scope

This edition is for instructors managing and grading their own courses. It has
no account administration, no user provisioning, and no cross-instance tooling,
and that is the product decision — not a gap to fill. See
[`docs/superpowers/specs/2026-09-05-faculty-edition-design.md`](docs/superpowers/specs/2026-09-05-faculty-edition-design.md)
for what was removed and why, and `DECISIONS.md` for why it was deleted rather
than hidden behind a build tag.

A change that adds a command must justify it against that scope, and must add
the command to `facultySurface` in `commands/surface_test.go` in the same
commit. A change that writes to Canvas must read the object back and print
evidence (`commands/submissions_readback.go` is the pattern).

## Branches

- `main` — where work lands. Branch from it: `feature/*`, `fix/*`, `docs/*`,
  `refactor/*`, `test/*`.
- `release/audited` — what faculty install: `main` plus the commit pinning
  `install.sh`. Do not develop here.

Open pull requests against `main`.

## Before you push

```bash
make check    # go vet, golangci-lint, gofmt, gosec if installed,
              # go test -race, and the integration tests
```

CI runs those plus the OS matrix, govulncheck, and the ≥80% coverage gate. A
red `make check` is a red PR, but a green one is not the whole of CI.

New code needs tests: a new command needs cmdtest coverage of its run function
*and* its option struct's `Validate()`, or the ≥80% total coverage gate drops.

## Commits

Conventional commits — `<type>(<scope>): <subject>`, with types `feat`, `fix`,
`docs`, `style`, `refactor`, `test`, `chore`, `ci`. The body should say why, not
restate the diff.

```
fix(submissions): verify rubric rows on read-back

Canvas accepts a rubric assessment and can silently drop criteria whose
ids do not match the association. The read-back now diffs each requested
criterion, so a dropped row reports verified: no instead of success.
```

## Pull requests

A PR should have a clear title and description, reference related issues, carry
tests for new behaviour, pass CI, and update `README.md` or the bundled skill
under `skills/canvas-cli/` when user-facing behaviour changes. Anything telling
a user how to install or update names the audited installer only — see the last
section of [AGENTS.md](AGENTS.md).

## Reporting

- Bugs and feature requests: GitHub Issues.
- Security vulnerabilities: **not** a public issue — see [SECURITY.md](SECURITY.md).

## Code of conduct

Be respectful and inclusive, welcome newcomers, give constructive feedback, and
assume good faith. Maintainers clarify the standard, act on unacceptable
behaviour, and moderate contributions.

## License

Contributions are licensed under the [MIT License](LICENSE). This project is a
fork; the upstream attribution is in the [README](README.md).
