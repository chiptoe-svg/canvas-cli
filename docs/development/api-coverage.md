# API Coverage & Spec Compliance

Canvas CLI is validated against Canvas's **official API specification** so its
endpoints stay correct as both the CLI and the Canvas API evolve.

## How it works

Canvas publishes a machine-readable API spec (Swagger 1.2) at
`https://<canvas-host>/doc/api/api-docs.json` plus per-resource files. The CLI
commits a slimmed inventory of that spec under `testdata/spec/`:

- `canvas_endpoints.json` — every documented endpoint (`method` + `path`),
  1086 in total.
- `canvas_models.json` — documented response models (field names and types).

A network-free contract test (`internal/api/spec_contract_test.go`, part of the
normal `go test` suite and CI) harvests every `/api/v1/...` path the service
layer calls and asserts each one matches a documented Canvas endpoint. **If a
command is wired to a path Canvas doesn't document, the build fails.** This has
already caught several real path bugs.

## Current coverage

The CLI implements **738 of 1086** documented endpoint patterns (**67%**),
across 98 command groups. The remaining gap is mostly niche admin/integration
surface (LTI internals, some SIS edge cases).

## Working with the spec

```bash
# Refresh the committed manifest from a live Canvas host.
# Any Canvas instance that serves /doc/api works; canvas.instructure.com
# IP-blocks datacenter requests, so the default is learn.canvas.net.
make spec-sync
CANVAS_SPEC_HOST=https://myschool.instructure.com make spec-sync

# Show which documented endpoints aren't implemented yet, grouped by resource.
make spec-coverage
```

When you add a new endpoint, take its exact path and verb from the committed
manifest (and field names from `canvas_models.json`); the contract test
enforces correctness, and the project's ≥80% coverage gate means each new
command should ship with tests.
