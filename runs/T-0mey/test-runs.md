# T-0mey — test runs

The completeness review (#221, row 45b).

## Red, the unit level (09:24Z, 2026-10-11)

`go test -run TestKeepOfCompleteness ./internal/run/`, against a stub `Keep`
that gave an error. The `pre-commit` hook runs the unit level, so the tests
landed with the code.

```
--- FAIL: TestKeepOfCompleteness (0.00s)
    completeness_test.go:28: Keep: not built
```

## Red, the integration level (09:25Z)

`go test -tags=integration -run
'TestTheStepCompleteness|TestTheCompletenessReviewWithNoOtherHarness'
./internal/run/`, with the world's second harness and its routes, before the
step: both harnesses are probed, and the restart ends at `draft` (the rows of
the table and of `sessions.tsv` cut).

```
--- FAIL: TestTheStepCompleteness (3.24s)
    run_integration_test.go:871: the restart: the steps [forge clone version lease probe rule-gaps meaning numbering spans draft], want [forge clone version lease probe rule-gaps meaning numbering spa
    run_integration_test.go:873: the detail of completeness ""
    run_integration_test.go:878: the review is not on devin, the harness that did not number or draft:
    run_integration_test.go:880: git show refs/heads/layup-records:intake/coverage.tsv: exit status 128
--- FAIL: TestTheCompletenessReviewWithNoOtherHarness (3.36s)
    run_integration_test.go:899: the last row {Name:draft Result:done Detail:2 rows, in spec/draft.tsv}; want completeness fail, not run
```

## Mutations (2026-10-11, 09:27Z to 09:30Z)

Each rule mutated alone, the unit test and, when it passes, the integration
tests of the step run. The first run (09:27Z to 09:28Z) caught 7 of 9:
`no-except` did not build (`except` unused), was not counted, and was rewritten
to pass over a harness of no register (`except+"-none"`); `session-error-ignored`
was not caught, as no test refused the review's file at the step;
`TestACoverageFileThatNamesNoFactIsRefused` was added, after the code, and its
red is that mutant's catch. The second run (09:29Z to 09:30Z): 9 of 9 caught,
each mutant built.

```
== except-ignored: caught run_integration_test.go:583: the last row of the restart with a harness lost {Name:completeness Result:done Detail:0 items, in intake/coverage.tsv}; want completeness fai
== fact-check-cut: caught completeness_test.go:47: a fact that is no fact: map["intake/coverage.tsv":"id\tkind\tfact\trequirement\tquestion\tsession\nC-001\tuncovered\t9\t—\tQ?\tS-0a0b0c0d\n"], <n
== req-check-cut: caught completeness_test.go:47: a requirement of no draft row: map["intake/coverage.tsv":"id\tkind\tfact\trequirement\tquestion\tsession\nC-001\tunsupported\t—\tR-007\tQ?\tS-0a0
== block-check-cut: caught completeness_test.go:47: an unsupported item with a fact: map["intake/coverage.tsv":"id\tkind\tfact\trequirement\tquestion\tsession\nC-001\tunsupported\t1\tR-001\tQ?\tS-0
== id-from-zero: caught completeness_test.go:28: Keep: line 2, column "id": the next ID is C-001
== no-except: caught run_integration_test.go:583: the last row of the restart with a harness lost {Name:completeness Result:done Detail:0 items, in intake/coverage.tsv}; want completeness fai
== not-run-cut: caught run_integration_test.go:583: the last row of the restart with a harness lost {Name:completeness Result:fail Detail:no admitted pair in the routing register for the role a
== done-before-cut: caught run_integration_test.go:885: the second restart's row {Name:completeness Result:done Detail:1 items, in intake/coverage.tsv}; want completeness done before
== session-error-ignored: caught run_integration_test.go:919: the last row {Name:completeness Result:fail Detail:line 1: no header row}; want completeness fail, naming the session and the row
```

## Green

At `283789e` with the documents (09:32Z to 09:35Z): `go build ./...`, `go vet
./...`, `go test ./...`, `go test -tags=integration ./...` and `go test
-count=1 -tags=e2e ./...` pass, every package `ok`.
