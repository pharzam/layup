# T-hxb9 — test runs

The answers (#223, row 47a).

## Red, the unit level (10:49Z, 2026-10-11)

`go test -run 'TestCollectAnswers|TestBudgetOf' ./internal/run/`, against
stubs (no answer and no block waiting; `budgetOf` none and true). The
`pre-commit` hook runs the unit level, so the tests landed with the code. Before
the code, the case of spaces around an intent value was added (decided: the
parse trims them).

```
--- FAIL: TestCollectAnswers (0.00s)
    answers_test.go:41: no answer line of an approver in the role: [], waiting []; want none, the idea owner's block waiting
--- FAIL: TestBudgetOf (0.00s)
    answers_test.go:75: 9.5, 10.0, 40.0: [], true; want the three rows (9.5 below 10.0 as numbers)
    answers_test.go:84: B not below U: [], true; want budget.tsv not written
    answers_test.go:84: B equal to U: [], true; want budget.tsv not written
    answers_test.go:84: a value that is no decimal: [], true; want budget.tsv not written
    answers_test.go:84: a value with a word: [], true; want budget.tsv not written
    answers_test.go:88: a field with no answer: [], true; want budget.tsv not written
```

## Red, the integration level (10:50Z)

`go test -tags=integration -run TestTheStepAnswers ./internal/run/`, with
scripted comments whose times come from the batch's `asked`, before the step:
each restart ends at `batch` (the rows of the table elided).

```
--- FAIL: TestTheStepAnswers (4.18s)
    run_integration_test.go:1043: the restart: the steps [forge clone version lease probe rule-gaps meaning numbering spans draft completeness batch], want [forge clone version lease probe rule-gaps m
    run_integration_test.go:1045: the detail of answers ""; want waiting
    run_integration_test.go:1077: the restart with the answers: the steps [forge clone version lease probe rule-gaps meaning numbering spans draft completeness batch], want [forge clone version lease
    run_integration_test.go:1079: the detail of answers ""
    run_integration_test.go:1081: git show refs/heads/layup-records:intake/answers.tsv: exit status 128
```

## Mutations (2026-10-11, 10:52Z to 10:55Z)

Each rule mutated alone, the unit tests and, when they pass, the integration
test of the step run. The first run (10:52Z to 10:53Z) caught 12 of 17:
`cr-kept`, `empty-answers` and `dash-answers` were not caught, as a later line
of the fixture hid each; the fixture now puts a CR on a line whose answer is
kept, and an empty answer and one of `—` after the last good line of Q-001.
`no-time-check` did not build (`asked` unused), was not counted, and was
rewritten. `waiting-not-stop` was not caught, as no step follows `answers` yet;
`TestAWaitingStepIsTheLastRow` (unit) was added, after the code, and its red is
that mutant's catch. The second run (10:54Z to 10:55Z): 17 of 17 caught, each
mutant built.

```
== cr-kept: caught answers_test.go:63: answer 1: ["Q-002" "9.5\r" "707" "303" "2026-10-11T10:05:00Z"], want ["Q-002" "9.5" "707" "303" "2026-10-11T10:05:00Z"]
== empty-answers: caught answers_test.go:63: answer 0: ["Q-001" "" "708" "303" "2026-10-11T10:06:00Z"], want ["Q-001" "second" "708" "303" "2026-10-11T10:06:00Z"]
== dash-answers: caught answers_test.go:63: answer 0: ["Q-001" "—" "708" "303" "2026-10-11T10:06:00Z"], want ["Q-001" "second" "708" "303" "2026-10-11T10:06:00Z"]
== not-strictly-after: caught answers_test.go:41: no answer line of an approver in the role: [["Q-001" "at the time of the batch" "702" "303" "2026-10-11T10:00:00Z"]], waiting []; want none, the idea
== no-time-check: caught answers_test.go:41: no answer line of an approver in the role: [["Q-001" "at the time of the batch" "702" "303" "2026-10-11T10:00:00Z"]], waiting []; want none, the idea
== any-role: caught answers_test.go:41: no answer line of an approver in the role: [["Q-001" "the Operator, in the other block's role" "703" "101" "2026-10-11T10:01:00Z"]], waiting []; want
== first-line-wins: caught answers_test.go:63: answer 0: ["Q-001" "first" "707" "303" "2026-10-11T10:05:00Z"], want ["Q-001" "second" "708" "303" "2026-10-11T10:06:00Z"]
== ack-counts: caught answers_test.go:41: no answer line of an approver in the role: [], waiting []; want none, the idea owner's block waiting
== trim-cut: caught answers_test.go:82: spaces around the values: [], false; want them parsed
== b-u-text: caught answers_test.go:77: 9.5, 10.0, 40.0: [], false; want the three rows (9.5 below 10.0 as numbers)
== b-equal-u: caught answers_test.go:91: B equal to U: [["band.B" "10.0" "Q-002"] ["band.U" "10.0" "Q-003"] ["appetite.hours" "40.0" "Q-004"]], true; want budget.tsv not written
== waiting-cut: caught run_integration_test.go:1045: the detail of answers "0 answers, in intake/answers.tsv; budget.tsv not written (a field of the intent form does not parse)"; want waiting
== app-copied: caught run_integration_test.go:1100: copies.tsv holds the App's own comment:
== state-not-updated: caught run_integration_test.go:1087: questions.tsv: each question answered
== budget-not-written: caught run_integration_test.go:1089: git show refs/heads/layup-records:budget.tsv: exit status 128
== done-before-cut: caught run_integration_test.go:1104: the third restart's row {Name:answers Result:fail Detail:git commit -m layup run: the answers of the batch: exit status 1}; want answers don
== waiting-not-stop: caught steps_test.go:189: the rows [{Name:a Result:done Detail:the first} {Name:b Result:done Detail:waiting: no answer yet} {Name:c Result:done Detail:never run}], 3 steps ran;
```

## Green

The first full run (10:57Z to 11:01Z) failed at the integration level:
`TestABatchWithNoCommentIsPostedFromTheRecord` (row 46) gave "the restart
changed questions.tsv", as its person's comment held an answer line and the
step `answers` now reads it. That comment needs only to start with the header,
so its answer line was taken out. Then (11:01Z to 11:05Z): `go build ./...`,
`go vet ./...`, `go test ./...`, `go test -tags=integration ./...` and `go test
-count=1 -tags=e2e ./...` pass, every package `ok`.
