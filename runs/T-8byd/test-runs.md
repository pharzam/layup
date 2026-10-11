# T-8byd — test runs

The review of the release of `M2c` (#227, row 49a). The green run of the check
at the release commit came first (13:46Z, before the plan, to learn whether the
script needed an adaptation), then the red runs (13:54Z), then the review; the
script is unchanged, so the order changes no result.

## Red (13:54Z, 2026-10-11)

`runs/T-efmy/release-check.sh`, unchanged, on two commits of a scratch clone of
`0b881fa`, each with one mutant: the check still bites on the tree of `M2c`.

1. A program out of the list of check (3): `exec.Command("curl", …)` in
   `internal/run/batch.go`. Exit 1:

```
FAIL
FAIL	github.com/pharzam/layup/cmd/layup	1.158s
FAIL
FAIL: TestPackageRules
FAIL: a call that starts a program, out of the list: internal/run/batch.go:293:func startsCurl() error { return exec.Command("curl", "https://example.invalid").Run() }
```

2. A verb that reaches a remote in a function of no allowance, check (4): a
   function `Leak` of `internal/git/git.go` that runs `git push`. Exit 1:

```
FAIL: a git verb that reaches a remote, out of its function: internal/git/git.go:495: push in Leak
```

## Green

`sh runs/T-efmy/release-check.sh 0b881fa…` (13:46Z): exit 0, `== PASS`, with no
adaptation of the script; its output is
[`release-check.txt`](release-check.txt). Against the output at the release of
`M2b` (`runs/T-4tjy/release-check.txt`), the new lines are `internal/spec` among
the packages of the module, and `CommitTime` among the functions of
`internal/git`; the sites of check (3), the hits of check (4) and the calls of
check (6) are the same.

## The review (the uat)

A fresh read-only session of Claude Fable 5.1 (13:54Z, 8 min 39 s;
`modelUsage`: `claude-fable-5-1` only, 2,185,023 tokens, USD 10.31; three
commands denied by its permission mode, `go list` among them, as its record
says) read the release at `0b881fa` and recorded `holds` for REQ-015 and for
REQ-017, eighteen findings, none a defect:
[`release-review.md`](release-review.md), posted on #227.

## Criterion 3 (14:04Z to 14:08Z)

With this task's files in the tree: `go build ./...`, `go vet ./...`, `go test
./...` and `go test -tags=integration ./...` pass, every package `ok`;
`adr-lint`, `prd-lint`, `link-lint`, `setup-check` and `run-discipline-tests`
exit 0.
