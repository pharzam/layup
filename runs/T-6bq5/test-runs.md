# T-6bq5 — the test runs

Evidence for row 24 (#129): each test red first, for the right reason, then green.

## `internal/forge`: red 1 (2026-10-08T08:40Z)

`forge_test.go` (a stand-in adapter that declares five capabilities, six, none; `CheckPermissions` on the full set, `write` for `read`, more than `M2a` uses, no `issues`, `contents: read`, none, an unknown level) before `forge.go`: `go test ./internal/forge/` did not compile (`undefined: Capability`, `undefined: Installation`).

## `internal/forge`: red 2, then green (2026-10-08T08:42Z)

With `Missing` and `CheckPermissions` as stubs that give nothing: six errors, each for its case (`Missing = [], want [commit statuses]`; `no issues: CheckPermissions = [], want [issues: write]`; …); the cases whose answer is "none" passed, as a stub must. Then the real functions: the package passes; `go list -deps ./internal/forge` names no `net`, `net/http`, `crypto/tls` or `crypto/x509`.

## `internal/forge/github`: red (2026-10-08T08:55Z)

`github_test.go` (unit: the JWT, the first line of a message) and `github_integration_test.go` (an `httptest` server on loopback that plays the calls of `M2a` and checks the JWT of each JWT call) against a stub adapter whose methods give an error: each test failed for its case (`the adapter does not declare [...]`, `Installation = {ID:0 ...}, stub`, `the JWT has not three parts`, `firstLine(...) = "", want "Not Found"`, each forge-error case `stub; want a forge.Error`, each rate-limit case `stub; want the call after the wait`).

## `internal/forge/github`: green, and three mutations (2026-10-08T09:05Z)

The real adapter: `go test -tags=integration ./internal/forge/...` passes. Mutations on a backup copy of `github.go`, each put back (`cmp` equal):
- a `403` with `x-ratelimit-reset` is a rate limit whatever `x-ratelimit-remaining` says → `a 403 with a reset time but requests left: <nil>; want a forge.Error`;
- the token kept while 3 minutes remain → `Token with four minutes left: <nil>; want a new one`;
- a `401` retried with no limit → first not caught (the fake gave two `401`s); the test then gives three and counts the requests → `4 requests after two 401, want 2 (one retry)`.

## The documents (2026-10-08T09:15Z)

Red before the edits: `sh runs/T-6bq5/docs.sh` gave sixteen `FAIL` lines, exit 1. After the edits of `forge.md`, `run.md`, `guardrails.md`, `packages.md`, the traceability, the PRD and `endpoints.md`: sixteen `ok`, exit 0.
