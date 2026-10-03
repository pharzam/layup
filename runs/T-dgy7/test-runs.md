# T-dgy7: the test runs

The runs of `T-dgy7` (#95). Host: macOS, `go1.27.1`, 2026-10-03, UTC. "…"
marks a cut part of an output.

## The red runs

Two skeletons, as row 17's condition gives (D4 of the plan): the block test is
red first, then each rule test.

1. **Run A** (03:36:54Z), the schema `StallsSchema` with its name and no
   column, and rules that check nothing:

   ```text
   $ go test -count=1 -tags=integration ./internal/records/ ./internal/tsv/
   --- FAIL: TestTheSchemasEqualTheirBlocks (0.00s)
       records_integration_test.go:27: schema stalls: the number of columns: the block has "11"; the Go schema has "0"
   FAIL	github.com/pharzam/layup/internal/records
   ok  	github.com/pharzam/layup/internal/tsv
   ```

   The test of `internal/tsv` passes, as the name `stalls` moved from the list
   `notYetBuilt` to `built` in the same change as the block test of its owner.

2. **Run B** (03:37:58Z), the 11 columns of the block, and rules that check
   nothing: the block test passes, the three stalls of `stallRows` (closed,
   open, of the orchestrator with a `diagnosis-failed` row) and one open stall
   alone pass, and each of the 30 rule cases fails on its own rule:

   ```text
   $ go test -count=1 ./internal/records/
   --- FAIL: TestStallRowsThatBreakARule (0.00s)
       stalls_test.go:100: no time: <nil>; want an error of line 2, column "time"
       stalls_test.go:100: no task: <nil>; want an error of line 2, column "task"
       stalls_test.go:100: a stall row with no trigger: <nil>; want an error of line 2, column "trigger"
       … (each of the 21 cases of the columns of a row)
       stalls_test.go:100: a diagnosis before its stall row: <nil>; want an error of line 2, column "stall"
       stalls_test.go:100: both diagnosis rows of one stall: <nil>; want an error of line 5, column "kind"
       stalls_test.go:100: an outcome before the second row: <nil>; want an error of line 3, column "kind"
       stalls_test.go:100: ST-002 first: <nil>; want an error of line 2, column "stall"
       stalls_test.go:100: ST-001 then ST-003: <nil>; want an error of line 3, column "stall"
       stalls_test.go:100: a row before the time of the row before it: <nil>; want an error of line 4, column "time"
       stalls_test.go:100: a row of a stall with another task: <nil>; want an error of line 4, column "task"
       stalls_test.go:100: an orchestrator stall of a task: <nil>; want an error of line 3, column "task"
       stalls_test.go:100: a stall of project with another trigger: <nil>; want an error of line 2, column "task"
   FAIL
   ```

   Two cases of the order were first written so that the key `(stall, kind)`
   of `internal/tsv` would refuse them before the order check: "a diagnosis
   before its stall row" and "ST-002 before ST-001" (now "ST-002 first") each
   held a second row of the same key. A read of the table before run B found
   them; they now hold one row of each key, so each fails on its own rule (the
   lesson "A known-bad fixture that fails for another reason" of
   `guardrails.md` §2). The test was fixed, before the code was written.

## The green runs

`go test -count=1 ./internal/records/` and `go test -count=1 -tags=integration
./internal/records/ ./internal/tsv/`: `ok`. The per-row rules run on all the
rows before the order rules, so a row that breaks both gives the error of the
row ("an orchestrator stall of a task" gives line 3, column `task`, from
`CheckStall`).

## The freeze

On the frozen head `6e3fa70`, 03:45:06Z to 03:47:29Z, the 18 steps of the
ladder, each exit 0: the eight local checks of `AGENTS.md` (with `git diff
--check` against `origin/main`), `go build`, `go vet` with each tag, `gofmt
-l` (no file), the three test levels, the two runs with `-race` (with
`./internal/records/`), and the harness of the fixtures of `setup-check.sh`
(44 passed). CI of PR #118 passed its job `tests` on it before review round 1.
