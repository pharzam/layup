# The test runs of T-8kqn

## Red 1: no Go value (2026-10-07T14:28Z)

`go test -tags=integration ./internal/records/`:

```
# github.com/pharzam/layup/internal/records [github.com/pharzam/layup/internal/records.test]
internal/records/records_integration_test.go:22:12: undefined: StartSchema
internal/records/records_integration_test.go:22:38: undefined: ApproversSchema
internal/records/records_integration_test.go:22:64: undefined: LeaseSchema
internal/records/records_integration_test.go:22:87: undefined: CopiesSchema
internal/records/start_test.go:72:15: undefined: ReadStart
```

## Red 2: the stub, with a wrong column in `StartSchema` and readers that check no rule

```
    records_integration_test.go:29: schema start: column 2 (value): the type: the block has "text"; the Go schema has "int"
    start_test.go:73: line 2, column "value": "v0.1.0" is not a value of int
    start_test.go:118: an empty register and no harness row: line 2, column "value": "v0.1.0" is not a value of int
    start_test.go:136: an empty since: read, want an error
    start_test.go:136: a source that is not start or an ID: read, want an error
    start_test.go:136: an empty source: read, want an error
    start_test.go:136: an empty login: read, want an error
    start_test.go:160: no row: read, want an error
    start_test.go:160: two rows: read, want an error
    start_test.go:160: a run ID of 15 characters: read, want an error
    start_test.go:160: an empty started: read, want an error
    start_test.go:160: an empty heartbeat: read, want an error
    start_test.go:160: an empty state: read, want an error
    start_test.go:160: a run ID in capitals: read, want an error
    start_test.go:160: an empty host: read, want an error
    start_test.go:160: an empty version: read, want an error
    start_test.go:186: seen 0: read, want an error
    start_test.go:186: a body at another path: read, want an error
    start_test.go:186: an empty author login: read, want an error
    start_test.go:186: an empty sha256: read, want an error
    start_test.go:186: seen does not start at 1: read, want an error
    start_test.go:186: a gap in seen: read, want an error
FAIL
FAIL	github.com/pharzam/layup/internal/records	0.197s
FAIL
```

## Red 3: `StartSchema` right, the readers still check no rule

The comparison test passes, and the valid start record reads. Each of the 21 rule cases of start fails with `read, want an error` (`start_test.go:114`), as do the cases of approvers, lease and copies of red 2.

## D3 and D4 red: `sh runs/T-8kqn/docs.sh` (2026-10-07T14:30Z)

```
FAIL the Job cell of internal/records does not name the records of Start
FAIL traceability has no row of TestAValidStartIsRead with T-8kqn
FAIL traceability has no row of TestStartRefusesEachBrokenRule with T-8kqn
FAIL traceability has no row of TestApproversRefusesEachBrokenRule with T-8kqn
FAIL traceability has no row of TestLeaseRefusesEachBrokenRule with T-8kqn
FAIL traceability has no row of TestCopiesRefusesEachBrokenRule with T-8kqn
FAIL the row of TestTheSchemasEqualTheirBlocks of internal/records does not name NFR-001 and T-8kqn
FAIL the Test cell of NFR-001 does not name the tests of the records of Start (T-8kqn)
exit 1
```

## Green: `go test -tags=integration ./internal/records/ ./internal/tsv/`

With `start.go` and the four names in `built`: `ok` for both packages; `go build ./...`, `go vet ./...`, `go test ./...` and `go test -tags=integration ./...` pass.

## D3 and D4 green (2026-10-07T14:30Z)

`sh runs/T-8kqn/docs.sh`: the eight rules `ok`, exit 0.

## The notes of round 1, in the close-out (2026-10-07T14:40Z)

Two cases added (note 1: a harness of the register with no rows; note 2: a register ID that is not of the form `<word>`). With the `<word>` check of `ReadStart` turned off on a backup copy of `start.go`, the second fails with `a register ID that is not of the form <word>: read, want an error`; the file was put back, and `go test ./internal/records/` passes.
