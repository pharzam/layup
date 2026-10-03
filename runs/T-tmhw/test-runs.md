# T-tmhw: the test runs

The runs of `T-tmhw` (#94). Host: macOS, `go1.27.1`, 2026-10-03, UTC. "…"
marks a cut part of an output.

## The red runs

Two skeletons, as condition 1 of the plan review asks: the block test is red
first, then each rule test.

1. **Run A**, the two schemas with their names and no column, and rules that
   check nothing:

   ```text
   $ go test -count=1 -tags=integration -run TestTheSchemasEqualTheirBlocks ./internal/records/
   --- FAIL: TestTheSchemasEqualTheirBlocks
       records_integration_test.go:27: schema telemetry: the number of columns: the block has "21"; the Go schema has "0"
       records_integration_test.go:27: schema prices: the number of columns: the block has "8"; the Go schema has "0"
   ```

2. **Run B**, the columns of the blocks, and rules that check nothing: the
   block test passes, the rows that must pass pass, and each of the 24 rule
   cases fails on its own rule:

   ```text
   $ go test -count=1 ./internal/records/
   --- FAIL: TestTelemetryRowsThatBreakARule
       records_test.go:106: a session ID that is not hexadecimal: <nil>; want an error of line 3, column "session"
       records_test.go:106: a latency with no first output: <nil>; want an error of line 3, column "latency_s"
       … (each of the 20 cases of the telemetry rows)
       records_test.go:125: CheckTelemetry of a computed row with no price: <nil>; want a *RowError of the column price
   --- FAIL: TestPriceRows
       records_test.go:163: a currency in lower case: <nil>; want an error of line 3, column "currency"
       records_test.go:163: a source of another scheme: <nil>; want an error of line 3, column "source"
       records_test.go:163: a source that is not a URL: <nil>; want an error of line 3, column "source"
   ```

   The first run B failed for another reason: the good row and the broken row
   of each case had one session ID, so the key rule of `internal/tsv` failed
   first ("the key … repeats the key of line 2"). The broken row now has its
   own ID; the test was fixed, not the code.

   The three cases of the types (a bad ID `P-1`, a class `total`, a price
   `1e3`) pass in run B, as the types of the columns refuse them.

## The green runs

`go test -count=1 ./internal/records/` and `go test -count=1 -tags=integration
./internal/records/ ./internal/tsv/`: `ok`. `TestPackageRules`
(`cmd/layup`, integration) passes with the new row of `internal/records`;
`net/url`, which the check of a source uses, depends on none of `net`,
`net/http` and `crypto/tls` (`go list -deps net/url`).

## The fix of review round 1 (cycle 1)

Round 1 (Claude Fable 5.1, on `9beb171`; the record is on #94) gave
`material`: (1) `seconds` gave whether both values were times, and each caller
dropped it, so a row with no `start` and no `end` and `duration_s` 0 passed;
(2) a column whose block rule has no clause for `—` could hold `—` (a row with
no `tokens_status`, no `money_status`, no `billing`, or a price row with no
price), which `docs/spec/README.md` gives the owner of a record to check. The
fix: a column with no such clause never holds `—` (checked first), `seconds`
reads the form of the type `time` and its result is used, and the clause of a
subscription in `records.md` is written as an explanation (note 4).

The red runs, with the new cases, on the code of `9beb171`:

```text
$ go test -count=1 ./internal/records/
    records_test.go:119: no start and no end: <nil>; want an error of line 3, column "start"
    records_test.go:119: the empty value of task: <nil>; want an error of line 3, column "task"
    … (role, harness, model, billing, tokens_status, money_status)
    records_test.go:142: CheckTelemetry of a start with an offset: <nil>; want a *RowError of the column start
    records_test.go:191: the empty value of harness: <nil>; want an error of line 3, column "harness"
    … (model, class, price, date)
```

14 cases fail on the missing rule; four more (an empty `end`, `duration_s`,
`currency` or `source`) gave their column already, through another rule, and
now through the new one. The green runs of the fix: `go test -count=1
./internal/records/`, and with `-tags=integration` `./internal/records/
./internal/tsv/`: `ok`. On the fix head `ed8ac6b`, 03:06:38Z to 03:08:56Z,
the 18 steps of the ladder, each exit 0; CI of PR #117 passed its job `tests`
on it before review round 2, which reproduced the 14 red cases on the code of
`9beb171`.
