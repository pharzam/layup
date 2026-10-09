# The test runs of T-3py1

## Red 1: no Go value (2026-10-09T11:46Z)

`internal/records/session_test.go` and the five names in `TestTheSchemasEqualTheirBlocks`, before `session.go`: `go test ./internal/records/` does not compile.

```
# github.com/pharzam/layup/internal/records [github.com/pharzam/layup/internal/records.test]
internal/records/session_test.go:48:61: undefined: ReadSessions
internal/records/session_test.go:53:15: undefined: ReadSessions
internal/records/session_test.go:59:15: undefined: ReadSessions
internal/records/session_test.go:84:62: undefined: ReadHarnesses
```

## Red 2: stubs (2026-10-09T11:46Z)

`session.go` as a stub: the five schemas with one column each, and readers that check no rule. Each valid file reads (the stubs accept everything); 76 refusal cases fail (`grep -c "want an error"`), in these tests:

```
--- FAIL: TestSessionsRefusesEachBrokenRule (0.00s)
--- FAIL: TestHarnessesRefusesEachBrokenRule (0.00s)
--- FAIL: TestRoutingRefusesEachBrokenRule (0.00s)
--- FAIL: TestEventsRefusesEachBrokenRule (0.00s)
--- FAIL: TestEventsAreOnlyAppended (0.00s)
--- FAIL: TestResultRefusesEachBrokenRule (0.00s)
```

`go test -tags=integration -run TestTheSchemasEqualTheirBlocks ./internal/records/` fails for each of the five blocks (9 lines name them).

## Green (2026-10-09T11:49Z)

Each with exit 0 on the tree of the commit `feat: T-3py1 …`: `go build ./...`, `go vet ./...`, `gofmt -l internal/` (empty), `go test ./...`, `go test -tags=integration ./...`; `adr-lint`, `prd-lint`, `link-lint`, `setup-check`, `run-discipline-tests`, `git diff --check`.

## The fix of round 1 (2026-10-09T12:01Z)

Finding 1: `CheckResult` checks the numbers of each kind as a set, 1 to k, in any order of the file; `TestAValidResultIsRead` reads a result whose rows are in another order (it failed at line 2 before the fix: "the next status is n 1"). Note 2, the two session IDs now pass the type and fail the rule: with the check of `sessionForm` taken out of `CheckHarness` and `CheckEvent` (a copy of the file, put back after), `go test ./internal/records/` fails both cases ("want an error in column session"), and passes with it. Each check of the green run above passes again on the fix tree.
