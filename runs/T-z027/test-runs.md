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

## The fix of round 1 (2026-10-09T20:46Z to 20:50Z)

Finding 1: a head that deletes `docs/guardrails.md` stopped the hook with the
run's own error. `TestAHeadThatDeletesTheGuardrailsIsRefused`, on the code of
`9d53a8b`:

```
push_integration_test.go:128: TaskSession: "S-6ba9470d", git show --end-of-options d2a0de92975ed7dcf9f3f27c4329b5ecabafde79:docs/guardrails.md --: exit status 128: fatal: bad revision 'd2a0de92975ed7dcf9f3f27c4329b5ecabafde79:docs/guardrails.md'
```

The head's file is now looked for with `LsTree`, and a head that lacks it
keeps an empty `Head`, which no added line passes. Note 2: a records commit
with no register refuses, and a commit or a `git` that cannot be read is the
run's own error. Note 6: `TestPassed`, the unit test of `passed`. Note 4: the
demo's own diff runs through `gitOut`, with no configuration of the host. Note
5: the assertion "nothing is pushed" cannot fail until row 37b pushes. The
mutations of the fix, with `TestPassed` among the tests, each in full (cut at
400 characters):

```
== passed-refused: passed: a session refused at its end is checked (exit 1)
end_integration_test.go:150: TaskSession: "S-04c09b13", git rev-parse --verify --end-of-options refs/layup/sessions/S-04c09b13: exit status 128: fatal: Needed a single revision
push_integration_test.go:165: the refusals ["artifact" "rule-path 72f90d80b29c172552a5af5602d90fed8af2e131ec5e5969d5d679758b883788"]; want artifact only
session_test.go:323: done, then refused: true, want false
== passed-done: passed: a result of any class is checked (exit 1)
end_integration_test.go:130: TaskSession: "S-b2cbf609", git rev-parse --verify --end-of-options refs/layup/sessions/S-b2cbf609: exit status 128: fatal: Needed a single revision
session_test.go:323: another class: true, want false
== register-missing: pushTask: a records commit with no register is the run's own error (exit 1)
push_integration_test.go:147: TaskSession: "S-6bf2535a", no register
== guardrails-deleted: pushTask: a head that lacks docs/guardrails.md is read with Show (exit 1)
push_integration_test.go:128: TaskSession: "S-62021a55", git show --end-of-options 107db12a99a2ad32f1939dce169698a6b1febe78:docs/guardrails.md --: exit status 128: fatal: bad revision '107db12a99a2ad32f1939dce169698a6b1febe78:docs/guardrails.md'
```

## The fix of round 2, after O-190 (2026-10-09T21:00Z to 21:02Z; committed 2026-10-10T03:29Z)

Finding 1 of round 2: a head whose `docs/guardrails.md` is a gitlink to a
commit that the clone lacks stopped the hook with the run's own error.
`TestAHeadWhoseGuardrailsIsAGitlinkIsRefused`, on the code of `cd9cdaf`:

```
push_integration_test.go:149: TaskSession: "S-13cfeb0b", git show --end-of-options ab3e3028007f30ba596490041aa1631368c3403e:docs/guardrails.md --: exit status 128: fatal: bad object ab3e3028007f30ba596490041aa1631368c3403e:docs/guardrails.md
```

The head's entry must now be a file (`Type` `blob`); another entry keeps an
empty `Head`, which no added line passes. The fix was written while the
Operator's answer was open, and committed after O-190 (a). Its mutation, and
that of the deletion again, each in full (cut at 400 characters):

```
== guardrails-deleted: pushTask: a head that lacks docs/guardrails.md is read with Show (exit 1)
push_integration_test.go:128: TaskSession: "S-babdb5b3", git show --end-of-options 8ee4ddb83e41f0755ee32c0ef8947410b624005a:docs/guardrails.md --: exit status 128: fatal: bad revision '8ee4ddb83e41f0755ee32c0ef8947410b624005a:docs/guardrails.md'
push_integration_test.go:149: TaskSession: "S-97f1d6f0", git show --end-of-options b2c627a294f371f3b247a4634f2019851c2f9aa6:docs/guardrails.md --: exit status 128: fatal: bad object b2c627a294f371f3b247a4634f2019851c2f9aa6:docs/guardrails.md
== guardrails-gitlink: pushTask: a gitlink at docs/guardrails.md is read with Show (exit 1)
push_integration_test.go:149: TaskSession: "S-475aeec3", git show --end-of-options 0aae19473b640f79c607895e53bd191573d89510:docs/guardrails.md --: exit status 128: fatal: bad object 0aae19473b640f79c607895e53bd191573d89510:docs/guardrails.md
```

## Green (2026-10-09T20:34Z, again at 20:50Z after the fix of round 1, and on 2026-10-10 at 03:29Z after the fix of round 2)

Each with exit 0: `go build ./...`, `go vet ./...`, `gofmt -l internal cmd`
(empty), `go test ./...`, `go test -tags=integration ./...`; `adr-lint`,
`prd-lint`, `link-lint`, `setup-check`, `run-discipline-tests`,
`git diff --check`, `sh runs/T-fdaq/parts.sh`.
