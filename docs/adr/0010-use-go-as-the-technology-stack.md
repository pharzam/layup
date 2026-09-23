# 0010. Use Go as the technology stack

Date: 2026-09-23

## Status

Accepted

## Context

Armature is domain-free and names no stack; its test, lint and security commands
are `‹…›` markers that an adopter fills (PSB `F-0001` §2). PSB Invariant 4 says a
value that comes from a guess is a defect, so the stack must come from a
decision with a source.

The Operator selected Go on 2026-09-23, in the session, from three options:
Go, TypeScript on Node, and Python with uv. The reasons given with the options:
Go builds one static binary with no runtime, so it installs into a target
repository with nothing else; a terminal UI library exists for the later cockpit;
Go was installed on the machine. The Operator then set the Go values in one batch
(decisions O-5 and O-6 on [#8](https://github.com/pharzam/layup/issues/8)).

## Decision

We will build LAYUP in Go. Tests use the Go toolchain's `go test` with no external
framework. Unit tests are `*_test.go` files beside the code
(`go test ./...`); integration and end-to-end tests carry the build tags
`integration` and `e2e` (`go test -tags=integration ./...`,
`go test -tags=e2e ./...`). Lint is `gofmt -l` (fails on any listed file) and
`go vet ./...`. The security track is govulncheck `v1.8.0`, `go vet`, and gitleaks.
The per-test-binary timeout is the Go default, `-timeout 10m`
(`go help testflag`).

We reject TypeScript on Node and Python with uv: each needs a runtime in every
target repository.

No panel was convened (ADR-0006). The Operator selected Go before the rule "each
new ADR gets a panel" was set (decision O-4, the same day); this record documents
a selection already made and does not claim a panel's comparison.

## Consequences

- The Go commands are recorded now but are **not active**: no Go code exists, so
  the hook lines stay commented and CI has no Go job. The task that lands the
  first Go code turns them on, in the same change.
- govulncheck `v1.8.0` needs Go 1.26 or later. The machine measured `go1.25.3`
  and later `go1.27.1` on 2026-09-23; `go.mod` must state a version of 1.26 or later.
- The coverage threshold stays open (decision O-7) until a baseline exists.
