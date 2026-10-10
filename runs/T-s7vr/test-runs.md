# T-s7vr — test runs

The fix of #199: a restart reads and writes `start.tsv` by the harnesses that
its rows name (plan of #199, with the conditions of its plan review).

## Red (2026-10-10, 16:17Z to 16:18Z)

The two tests before the code. D1, `go vet ./internal/records/`:

```
vet: internal/records/start_test.go:72:20: undefined: ReadStartAsWritten
```

D2, `go test -tags=integration -run TestARestartWithAHarnessRegisterThatChanged ./internal/run/`,
the error of #199 in each case:

```
run_integration_test.go:507: the restart with a harness gained: the step clone is fail: start/start.tsv: line 13, column "name": the next name is harness.claude.cap (the names of the block, in its order, with the harnesses of the register)
run_integration_test.go:517: the restart with a harness lost: the step clone is fail: start/start.tsv: line 15, column "name": the next name is pin.source (the names of the block, in its order, with the harnesses of the register)
```

## Mutations (2026-10-10, 16:19Z to 16:21Z)

Each rule mutated alone, each caught:

```
== harnessIDs-register: harnessIDs gives the register's IDs in a restart
run_integration_test.go:507: the restart with a harness gained: the step phase is fail: line 13, column "name": the next name is harness.claude.cap (...)
== restarted-false: cloneStep does not mark the run as a restart
run_integration_test.go:507: the restart with a harness gained: the step phase is fail: line 13, column "name": the next name is harness.claude.cap (...)
== ids-sorted: the harnesses sorted, not in the order of the rows
start_test.go:85: devin before claude: the harnesses [], line 13, column "name": the next name is harness.claude.cap (...); want devin claude
== no-word-check: no check of the form <word> of an ID
start_test.go:129: a harness ID that is not of the form <word>: read, want an error
start_test.go:129: an empty harness ID: read, want an error
```

`no-word-check` was not caught at first: the cases renamed the `.cap` row
alone, so the pair of rows was broken too and `CheckStart` refused it for
that. The cases now rename both rows of a harness, so only the form of the ID
is wrong.

## Green (2026-10-10, 16:23Z)

Each with exit 0: `go build ./...`, `go vet ./...`, `gofmt -l internal cmd`
(empty), `go test ./...`, `go test -tags=integration ./...`, `go test
-tags=e2e ./...`; `adr-lint`, `prd-lint`, `link-lint`, `setup-check`,
`run-discipline-tests`, `git diff --check`.
