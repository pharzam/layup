# The test runs of T-z027

The build ran from 2026-10-09T20:26Z (the claim of the branch) to 20:34Z, in the
order below; the steps inside it were not stamped one by one.

## The world of rows 36a and 36b

`sessionWorld` now gives a records commit that holds a rule-path register
(`AGENTS.md`, `docs/guardrails.md` with its exception, `docs/rules/`), as
the push hook reads the register of the session's records commit; the cases
of row 36b passed with it before the hook was written.

## Red 1: the checks before the hook

The tests of `push_integration_test.go`, with no push hook (the cases of
nothing refused, the guardrails exception and a session refused at its end,
passed; their mutations below show they can fail):

```
push_integration_test.go:44: the refusals []; want one rule-path with its payload
    push_integration_test.go:72: the refusals []; want one workflow with its payload
    push_integration_test.go:85: the refusals []; want base
    push_integration_test.go:116: the refusals []; want rule-paths
FAIL
FAIL
```

## Two cases added before the mutations

A register that its reader refuses (a second form of the case with no
register), and the payload and its event in one records commit (an assertion
of the demo).

## The mutations

A mutation of each rule of `push.go`, one at a time, run by a script that
restores the file from a copy after each, with the integration tests. Each
test line in full, cut at 400 characters. The first run of `payload-diff` was
not caught: the demo changed one file, whose own diff equals the head's. The
demo now changes two files, and `payload-diff` above ran with it; the others
ran with the demo of one file.

```
== passed-refused: passed: a session refused at its end is checked (exit 1)
push_integration_test.go:143: the refusals ["artifact" "rule-path 72f90d80b29c172552a5af5602d90fed8af2e131ec5e5969d5d679758b883788"]; want artifact only
== passed-done: passed: a result of any class is checked (exit 1)
end_integration_test.go:130: TaskSession: "S-302b7720", git rev-parse --verify --end-of-options refs/layup/sessions/S-302b7720: exit status 128: fatal: Needed a single revision
== base: pushTask: no check of the base (exit 1)
push_integration_test.go:90: the refusals []; want base
== register-missing: pushTask: a records commit with no register reads as an empty one (exit 1)
push_integration_test.go:128: a register that is broken false: the refusals []; want rule-paths
== register-broken: pushTask: a register that its reader refuses reads as an empty one (exit 1)
push_integration_test.go:128: a register that is broken true: the refusals []; want rule-paths
== check: pushTask: the result of rules.Check ignored (exit 1)
push_integration_test.go:44: the refusals []; want one rule-path with its payload
push_integration_test.go:77: the refusals []; want one workflow with its payload
== guardrails: pushTask: the change of docs/guardrails.md not given (exit 1)
push_integration_test.go:109: the refusals ["rule-path 878f184a659aebf2a8c154b4a0b7304688eb6c631e61ee5479f14c78d4b79d70"]; want none
== payload: refuse: the payload not written (exit 1)
push_integration_test.go:47: git show layup-records:payloads/5a404c9f0ae0e24afc3be7ac0caca46c226a0668380e8b79ca666871a3a94458: exit status 128
== payload-name: refuse: the event names no payload (exit 1)
push_integration_test.go:44: the refusals ["rule-path"]; want one rule-path with its payload
push_integration_test.go:77: the refusals ["workflow"]; want one workflow with its payload
== payload-diff: pushTask: the payload is the diff of names only (exit 1)
push_integration_test.go:53: the payload:
== one-commit: refuse: the payload in a commit of its own (exit 1)
push_integration_test.go:59: the commit of the payload changed ["payloads/5a404c9f0ae0e24afc3be7ac0caca46c226a0668380e8b79ca666871a3a94458"]; want the payload and the events
```

## Green (2026-10-09T20:34Z)

Each with exit 0: `go build ./...`, `go vet ./...`, `gofmt -l internal cmd`
(empty), `go test ./...`, `go test -tags=integration ./...`; `adr-lint`,
`prd-lint`, `link-lint`, `setup-check`, `run-discipline-tests`,
`git diff --check`, `sh runs/T-fdaq/parts.sh`.
