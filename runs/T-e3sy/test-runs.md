# The test runs of T-e3sy

The build ran on 2026-10-10 from 03:45Z (the claim of the branch) to 03:48Z, in
the order below; the steps inside it were not stamped one by one.

## Red 1: the order of the push and the bind, before the code

`TestThePushAndTheBindInOrder`, with the stand-in store logging each commit
and each push of a head in order:

```
# [github.com/pharzam/layup/internal/run]
vet: internal/run/session_test.go:352:15: undefined: errPushRefused
```

## The integration tests

`TestACleanHeadIsPushedThenBound` (the demo) and `TestARefusedPushBindsNothing`
were written after `pushAndBind`, so their red is the mutations below. Row
36b's demo now finds the commit of the result by its file, as the push and the
bind come after it.

## The mutations

A mutation of each rule, one at a time, run by a script that restores each
file from a copy after each; the unit order test for the rules of the order,
the integration tests for the push itself. Each test line in full, cut at 400
characters:

```
== push-first: pushAndBind: the push before its event (exit 1)
session_test.go:361: accepted: the calls ["push fedcba9876543210fedcba9876543210fedcba98 task/T-ab12/1" "commit push" "commit bound"], want ["commit push" "push fedcba9876543210fedcba9876543210fedcba98 task/T-ab12/1" "commit bound"]
session_test.go:361: refused: the calls ["push fedcba9876543210fedcba9876543210fedcba98 task/T-ab12/1" "commit push" "commit refused"], want ["commit push" "push fedcba9876543210fedcba9876543210fedcba98 task/T-ab12/1" "commit refused"]
== bound-first: pushAndBind: bound before the push (exit 1)
session_test.go:361: accepted: the calls ["commit push" "commit bound" "push fedcba9876543210fedcba9876543210fedcba98 task/T-ab12/1" "commit bound"], want ["commit push" "push fedcba9876543210fedcba9876543210fedcba98 task/T-ab12/1" "commit bound"]
session_test.go:368: accepted: the last event ["3" "bound" "1" "S-1a2b3c4d" "" "fedcba9876543210fedcba9876543210fedcba98" "task/T-ab12/1" "2026-10-09T12:00:00Z"], want ["2" "bound" "1" "S-1a2b3c4d" "" "fedcba9876543210fedcba9876543210fedcba98" "task/T-ab12/1" "2026-10-09T12:00:00Z"]
session_test.go:361: refused: the calls ["commit push" "commit bound" "push fedcba9876543210fedcba9876543210fedcba98 task/T-ab12/1" "commit refused"], want ["commit push" "push fedcba9876543210fedcba9876543210fedcba98 task/T-ab12/1" "commit refused"]
session_test.go:368: refused: the last event ["3" "refused" "1" "S-1a2b3c4d" "" "" "push-refused" "2026-10-09T12:00:00Z"], want ["2" "refused" "1" "S-1a2b3c4d" "" "" "push-refused" "2026-10-09T12:00:00Z"]
session_test.go:375: an error of the push: the remote cannot be reached, the events [["1" "push" "1" "S-1a2b3c4d" "" "fedcba9876543210fedcba9876543210fedcba98" "task/T-ab12/1" "2026-10-09T12:00:00Z"] ["2" "bound" "1" "S-1a2b3c4d" "" "fedcba9876543210fedcba9876543210fedcba98" "task/T-ab12/1" "2026-10-09T12:00:00Z"]]; want the error and the event push alone
== refused-event: pushAndBind: a refused push is the run's own error (exit 1)
session_test.go:358: refused: the push of the session's head was refused
== other-binds: pushAndBind: another error of the push binds (exit 1)
session_test.go:375: an error of the push: <nil>, the events [["1" "push" "1" "S-1a2b3c4d" "" "fedcba9876543210fedcba9876543210fedcba98" "task/T-ab12/1" "2026-10-09T12:00:00Z"] ["2" "bound" "1" "S-1a2b3c4d" "" "fedcba9876543210fedcba9876543210fedcba98" "task/T-ab12/1" "2026-10-09T12:00:00Z"]]; want the error and the event push alone
== no-push: pushAndBind: no push (exit 1)
push_integration_test.go:209: git rev-parse refs/heads/task/T-ab12/1: exit status 128
push_integration_test.go:251: the refusals [], want push-refused
push_integration_test.go:254: the events:
push_integration_test.go:257: the comment ["S-2325c474: developer session of T-ab12, attempt 1: done\n"]
== code1: PushHead: git's refusal not errPushRefused (exit 1)
push_integration_test.go:245: TaskSession: "S-35f775d8", git push --porcelain -- file:///tmp/TestARefusedPushBindsNothing4140080105/003/acme/target.git 35a660ebcfb5c53b06b5b0e9ba73407bb1dc1b65:refs/heads/task/T-ab12/1: exit status 1: error: failed to push some refs to 'file:///tmp/TestARefusedPushBindsNothing4140080105/003/acme/target.git'
== sha: appendEvents: the event push with no SHA (exit 1)
session_test.go:358: accepted: tasks/T-ab12/events.tsv: line 2, column "sha": sha is given exactly for push and bound
```

## Green (2026-10-10T03:48Z)

Each with exit 0: `go build ./...`, `go vet ./...`, `gofmt -l internal cmd`
(empty), `go test ./...`, `go test -tags=integration ./...`; `adr-lint`,
`prd-lint`, `link-lint`, `setup-check`, `run-discipline-tests`,
`git diff --check`, `sh runs/T-fdaq/parts.sh`.
