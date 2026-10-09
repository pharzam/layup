# The test runs of T-y10b

The order in which the steps ran: the tests and the fixture `enginedep` first; then the checker (`readTable` for the table of `M2b`, `readEngine`, the register form, the rule of the engine checks); then `packages.md`, `session.md` and `gate.md`, as `TestPackageRules` needs the line of the engine checks (D1 and D2 land together, as in `T-esfe`).

## Red 1: the tests before the checker (2026-10-09T12:04Z)

`go vet ./cmd/layup/` does not compile: the tests use the field `register`, `readEngine` and the new `checkRules`.

```
# github.com/pharzam/layup/cmd/layup
# [github.com/pharzam/layup/cmd/layup]
vet: cmd/layup/rules_test.go:98:64: unknown field register in struct literal of type row
```

## Red 2: the checker without its two rules (2026-10-09T12:07Z)

On a copy of `rules_checker_test.go` (put back after), with the rule of the engine checks turned off and the register form read as no register: the four engine cases and the three register cases of `TestEachRuleAndColumnFindsItsBreach` fail, and so do `TestReadTable`, `TestAGoodModuleKeepsTheRules` and each other case, as the good module's `internal/session` starts the command of its register by a variable; `TestPackageRules` gives 5 lines of failure.

## Red 3: `TestPackageRules` before `packages.md`

With the checker built and `packages.md` not yet changed: `rules_integration_test.go:38: 0 lines start with "The packages of the engine checks:", want 1`.

## Green (2026-10-09T12:07Z)

Each with exit 0 on the tree of the commit `feat: T-y10b …`: `go build ./...`, `go vet ./...`, `gofmt -l cmd internal` (empty), `go test ./...`, `go test -tags=integration ./...` (with `TestPackageRules` on the real module, `netimport` and `enginedep`); `adr-lint`, `prd-lint`, `link-lint`, `setup-check`, `run-discipline-tests`, `git diff --check`.
