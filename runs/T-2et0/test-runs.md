# T-2et0 — test runs

The draft (#220, row 45a).

## Red, the unit level (08:49Z, 2026-10-11)

`go test -run TestKeepOfDraft ./internal/run/`, against a stub `Keep` that gave
an error. The `pre-commit` hook runs the unit level, so the tests landed with
the code.

```
--- FAIL: TestKeepOfDraft (0.00s)
    draft_test.go:34: Keep: not built
```

## Red, the integration level (08:49Z)

`go test -tags=integration -run 'TestTheStepDraft|TestADraftThatCoversNoFact'
./internal/run/`, with the world's script answering the draft prompt, before
the step: each restart ends at `spans` (the rows of the table elided).

```
--- FAIL: TestTheStepDraft (2.63s)
    run_integration_test.go:808: the restart: the steps [forge clone version lease probe rule-gaps meaning numbering spans], want [forge clone version lease probe rule-gaps meaning numbering spans dra
    run_integration_test.go:810: the detail of draft ""
    run_integration_test.go:812: git show refs/heads/layup-records:spec/draft.tsv: exit status 128
--- FAIL: TestADraftThatCoversNoFactIsRefused (2.55s)
    run_integration_test.go:836: the last row {Name:spans Result:done Detail:the span check passed, in spec/check.tsv}; want draft fail, naming the session and the row
```

## Mutations (2026-10-11, 08:51Z to 08:52Z)

Each rule mutated alone, the unit test and, when it passes, the two
integration tests run: 11 of 11 caught, each mutant built.

```
== empty-covers-cut: caught draft_test.go:55: a row that covers nothing: map["spec/draft.tsv":"id\tkind\tstatement\tcovers\tcriterion\tsession\nR-001\trequirement\tThe tool does 1.\t1\tIt does 1.\tS
== no-fact-cut: caught draft_test.go:55: a number above the facts: map["spec/draft.tsv":"id\tkind\tstatement\tcovers\tcriterion\tsession\nR-001\trequirement\tThe tool does 1.\t1\tIt does 1.\tS-
== need-uncovered-cut: caught draft_test.go:55: a need covered by no row: map["spec/draft.tsv":"id\tkind\tstatement\tcovers\tcriterion\tsession\nR-001\tnon-functional\tK.\t2\tK.\tS-1a2b3c4d\nR-002\tsu
== constraint-not-checked: caught draft_test.go:55: a constraint covered by no row: map["spec/draft.tsv":"id\tkind\tstatement\tcovers\tcriterion\tsession\nR-001\trequirement\tD.\t1\tD.\tS-1a2b3c4d\nR-002\
== measure-cut: caught draft_test.go:55: a measure with no success criterion: map["spec/draft.tsv":"id\tkind\tstatement\tcovers\tcriterion\tsession\nR-001\trequirement\tD.\t1 3\tD.\tS-1a2b3c4d\
== measure-any-kind: caught draft_test.go:55: a measure with no success criterion: map["spec/draft.tsv":"id\tkind\tstatement\tcovers\tcriterion\tsession\nR-001\trequirement\tD.\t1 3\tD.\tS-1a2b3c4d\
== check-ignored: caught draft_test.go:55: a row that covers nothing: map["spec/draft.tsv":"id\tkind\tstatement\tcovers\tcriterion\tsession\nR-001\trequirement\tThe tool does 1.\t1\tIt does 1.\tS
== id-from-zero: caught draft_test.go:34: Keep: line 2, column "id": the next ID is R-001
== done-before-cut: caught run_integration_test.go:823: the second restart's row {Name:draft Result:done Detail:2 rows, in spec/draft.tsv}; want draft done before
== session-error-ignored: caught run_integration_test.go:836: the last row {Name:draft Result:fail Detail:line 1: no header row}; want draft fail, naming the session and the row
== no-spans-input: caught run_integration_test.go:808: the restart: the step draft is fail: the session S-01f0f3d1: refused file: the need fact 1 is covered by no row
```

## Green

The first full run (08:53Z to 08:56Z) failed at the integration level:
`TestAFailedSpanCheckRunsNumberingAgain` (row 44) gave "4 sessions of
product-owner; want 3", as the step draft now runs one more session of that
role after the span check passes. The test now counts the three of the
numbering and the draft's.

Then, with the count fixed (08:57Z to 09:00Z): `go build ./...`, `go vet ./...`,
`go test ./...`, `go test -tags=integration ./...` and `go test -count=1
-tags=e2e ./...` pass, every package `ok`.
