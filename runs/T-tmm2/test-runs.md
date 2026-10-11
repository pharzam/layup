# T-tmm2 — test runs

The answers to the follow-up (#225, row 47c). The order: the red of D1 and D2
before the code, then the code, then the documents with it. D3, the Go schema
of `run-steps`, changed with the enum in one edit, and its red is not recorded.

## Red, the unit level (12:30Z, 2026-10-11)

`go test -run TestBudgetAfterFollowUp ./internal/run/`, against a stub that
gave no row and no question open. The `pre-commit` hook runs the unit level, so
the tests landed with the code.

```
--- FAIL: TestBudgetAfterFollowUp (0.00s)
    followanswers_test.go:78: not yet written, B from the follow-up: "", open []; want "band.B 5.0 Q-005; band.U 10.0 Q-006; appetite.hours 40.0 Q-004", open ""
    followanswers_test.go:78: not yet written, the hours still unparsed: "", open []; want "", open "Q-004"
    followanswers_test.go:78: not yet written, B not below U after the follow-up: "", open []; want "", open "Q-002 Q-003"
    followanswers_test.go:78: not yet written, a field asked again with no answer: "", open []; want "", open "Q-004"
    followanswers_test.go:78: written, the hours asked again and replaced: "", open []; want "band.B 5.0 Q-002; band.U 10.0 Q-003; appetite.hours 50.0 Q-005", open ""
    followanswers_test.go:78: written, an answer that does not parse: "", open []; want "band.B 5.0 Q-002; band.U 10.0 Q-003; appetite.hours 40.0 Q-004", open "Q-004"
    followanswers_test.go:78: written, a replacement that leaves B not below U: "", open []; want "band.B 5.0 Q-002; band.U 10.0 Q-003; appetite.hours 40.0 Q-004", open "Q-003"
```

## Red, the integration level (12:30Z)

`go test -tags=integration -run TestTheAnswersToTheFollowUp ./internal/run/`,
before the step: each restart ends at `follow-up`.

```
--- FAIL: TestTheAnswersToTheFollowUpReplaceAValueOfTheBudget (5.09s)
    run_integration_test.go:1257: the detail of follow-up-answers ""; want waiting
--- FAIL: TestTheAnswersToTheFollowUpWriteTheBudget (5.43s)
    run_integration_test.go:1307: the detail of follow-up-answers ""
    run_integration_test.go:1309: git show refs/heads/layup-records:budget.tsv: exit status 128
```

## Mutations (2026-10-11, 12:32Z to 12:38Z)

Each rule mutated alone, the unit test and, when it passes, the integration
tests run. The first run (12:32Z to 12:35Z) caught 11 of 13: `open-not-set`
was not caught (no integration test left an answer to the follow-up unparsed),
so `TestAnUnparsedAnswerToTheFollowUpLeavesItsQuestionOpen` was added; and
`app-copied` was not caught, so the test of the replacement now asserts that
`copies.tsv` holds no comment of the App. Both were written after the code;
their red is the catch of their mutants. The second run (12:35Z to 12:38Z): 13
of 13 caught, each mutant built.

```
== followup-not-preferred: caught followanswers_test.go:78: not yet written, B from the follow-up: "", open ["Q-002" "Q-003"]; want "band.B 5.0 Q-005; band.U 10.0 Q-006; appetite.hours 40.0 Q-004", open "
== unwritten-open-cut: caught followanswers_test.go:78: not yet written, the hours still unparsed: "band.B 5.0 Q-002; band.U 10.0 Q-003; appetite.hours  ", open []; want "", open "Q-004"
== unwritten-pair-cut: caught followanswers_test.go:78: not yet written, B not below U after the follow-up: "band.B 12.0 Q-005; band.U 11.0 Q-006; appetite.hours 40.0 Q-004", open []; want "", open "Q
== replace-cut: caught followanswers_test.go:78: written, the hours asked again and replaced: "band.B 5.0 Q-002; band.U 10.0 Q-003; appetite.hours 40.0 Q-004", open ["Q-004"]; want "band.B 5.0
== replace-pair-cut: caught followanswers_test.go:78: written, a replacement that leaves B not below U: "band.B 5.0 Q-002; band.U 4.0 Q-005; appetite.hours 40.0 Q-004", open []; want "band.B 5.0 Q-0
== replace-unparsed-not-open: caught followanswers_test.go:78: written, an answer that does not parse: "band.B 5.0 Q-002; band.U 10.0 Q-003; appetite.hours 40.0 Q-004", open []; want "band.B 5.0 Q-002; band.
== no-follow-done-cut: caught run_integration_test.go:1086: the restart with the answers: the step follow-up-answers is fail: git commit -m layup run: the answers to the follow-up: exit status 1
== done-before-cut: caught run_integration_test.go:1289: the third restart's row {Name:follow-up-answers Result:fail Detail:line 7: the key (question "Q-005") repeats the key of line 6}; want follo
== waiting-cut: caught run_integration_test.go:1257: the detail of follow-up-answers "git commit -m layup run: the answers to the follow-up: exit status 1"; want waiting
== answered-not-set: caught run_integration_test.go:1279: the follow-up ["Q-005" "follow-up" "Q-004" "idea-owner" "before-delivery" "follow-up" "" "" "The wall-clock appetite of the delivery, in hou
== open-not-set: caught run_integration_test.go:1341: the question of B ["Q-002" "intent" "band.B" "idea-owner" "before-delivery" "intent" "" "" "B, in USD: below it the delivery goes on without
== budget-not-replaced: caught run_integration_test.go:1270: budget.tsv:
== app-copied: caught run_integration_test.go:1285: copies.tsv holds a comment of the App
```

## Green

At `0a45fde`, with the documents in the tree and not yet committed (12:40Z to
12:44Z): `go build ./...`, `go vet ./...`, `go test ./...`, `go test
-tags=integration ./...` and `go test -count=1 -tags=e2e ./...` pass, every
package `ok`.
