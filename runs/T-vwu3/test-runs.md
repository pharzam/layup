# T-vwu3 — test runs

The review of the answers and the follow-up (#224, row 47b). The order: the
red of D1 to D3 before the code, then the code, then the documents with it.
D4, the Go schema of `run-steps`, changed with the enum in one edit, and its
red is not recorded.

## Red, the unit level (11:30Z, 2026-10-11)

`go test -run 'TestKeepOfTheReview|TestFollowUpRows|TestTheCommentOfTheFollowUp'
./internal/run/`, against stubs (`Keep` an error, no follow-up row, an empty
comment). The first run held one fault of the fixture (a question `asked` with
the time of an answer, which `ReadQuestions` refuses); the fixture was fixed,
and the red below is its second run. The `pre-commit` hook runs the unit level,
so the tests landed with the code.

```
--- FAIL: TestKeepOfTheReview (0.00s)
    followup_test.go:46: Keep: not built
--- FAIL: TestFollowUpRows (0.00s)
    followup_test.go:103: missing, open, B not below U: ; want Q-007<Q-002 Q-008<Q-003 Q-009<Q-004 Q-010<Q-005
    followup_test.go:103: an unparsed hours field: ; want Q-007<Q-003 Q-008<Q-006
    followup_test.go:103: an open intent field with budget.tsv written: ; want Q-007<Q-003 Q-008<Q-006
    followup_test.go:103: an intent field with no answer, asked once: ; want Q-007<Q-003 Q-008<Q-006
--- FAIL: TestTheCommentOfTheFollowUp (0.00s)
    followup_test.go:131: the comment has no "LAYUP Intake: the follow-up\n":
    followup_test.go:131: the comment has no "\n## For the idea owner\n":
    followup_test.go:131: the comment has no "\nQ-007: What is TLA?\n`We use Go and TLA and XYZ.`\n":
    followup_test.go:131: the comment has no "\nQ-008: Which stack?\n\n":
    followup_test.go:131: the comment has no "\n## For the Operator\n":
    followup_test.go:131: the comment has no "`Q-NNN: <answer>`":
```

## Red, the integration level (11:31Z)

`go test -tags=integration -run TestTheReviewOfTheAnswersAndTheFollowUp
./internal/run/`, with the world's script answering the review prompt, before
the steps: the restart ends at `answers` (the rows of the table elided).

```
--- FAIL: TestTheReviewOfTheAnswersAndTheFollowUp (4.64s)
    run_integration_test.go:1150: the restart with the answers: the steps [forge clone version lease probe rule-gaps meaning numbering spans draft completeness batch answers], want [forge clone versio
    run_integration_test.go:1152: the detail of follow-up ""
    run_integration_test.go:1166: the follow-up asks again []; want the missing, the two unparsed and the open ["Q-001" "Q-002" "Q-003" "Q-004"]
    run_integration_test.go:1168: git show refs/heads/layup-records:intake/answers-review.tsv: exit status 128
```

After the code, `TestKeepOfTheReview` asked for "Q-001 twice" where the kept key
of the review's file refuses the repeat first ("repeats the key"); the case now
asks for that message. Rows 47a and 43a had tests that end at `answers`: the
restart now runs on into the review and the follow-up, so `TestTheStepAnswers`
asserts the review's closed states and "nothing to ask again", and the progress
of `TestTheStepRuleGaps` counts 15 rows.

## Mutations (2026-10-11, 11:36Z to 11:40Z)

Each rule mutated alone, the unit tests and, when they pass, the two
integration tests run. The first run (11:36Z to 11:37Z) caught 12 of 13: the
guard "when `budget.tsv` is not written" was an equivalent mutant, as
`budget.tsv` is written exactly when the three fields parse; the guard and its
parameter were taken out, so the code says the specification's words. The
second run (11:39Z to 11:40Z): 12 of 12 caught, each mutant built.

```
== review-unanswered: caught followup_test.go:65: a question with no answer: map[], the answered question Q-001 has no row; want it refused, "Q-003 is not an answered question"
== review-missing-row: caught followup_test.go:65: an answered question with no row: map["intake/answers-review.tsv":"question\tstate\treason\tsession\nQ-001\tclosed\tx\tS-0a0b0c0d\n" "questions.tsv":
== review-state-not-set: caught followup_test.go:56: the states ["answered" "answered" "asked"]; want closed, open, and asked for the question with no answer
== follow-asked-cut: caught followup_test.go:102: missing, open, B not below U: Q-007<Q-002 Q-008<Q-004 Q-009<Q-005; want Q-007<Q-002 Q-008<Q-003 Q-009<Q-004 Q-010<Q-005
== follow-open-cut: caught followup_test.go:102: missing, open, B not below U: Q-007<Q-003 Q-008<Q-004 Q-009<Q-005; want Q-007<Q-002 Q-008<Q-003 Q-009<Q-004 Q-010<Q-005
== follow-unparsed-cut: caught followup_test.go:102: missing, open, B not below U: Q-007<Q-002 Q-008<Q-003; want Q-007<Q-002 Q-008<Q-003 Q-009<Q-004 Q-010<Q-005
== b-u-pair-cut: caught followup_test.go:102: missing, open, B not below U: Q-007<Q-002 Q-008<Q-003; want Q-007<Q-002 Q-008<Q-003 Q-009<Q-004 Q-010<Q-005
== quote-of-row: caught followup_test.go:130: the comment has no "\nQ-007: What is TLA?\n`We use Go and TLA and XYZ.`\n":
== header-of-batch: caught followup_test.go:130: the comment has no "LAYUP Intake: the follow-up\n":
== posted-ignored: caught run_integration_test.go:1189: the third restart's row {Name:follow-up Result:done Detail:the comment of the follow-up posted from questions.tsv}; want follow-up done befo
== review-done-before-cut: caught run_integration_test.go:1189: the third restart's row {Name:answers-review Result:done Detail:no answered question; no session}; want answers-review done before
== empty-review-not-header: caught followup_test.go:70: the record of no review "question\tstate\treason\tsession\nQ-001\tclosed\tx\tS-00000000\n"
```

## Green

At `dceb86c`, with the documents in the tree and not yet committed (11:41Z to
11:45Z): `go build ./...`, `go vet
./...`, `go test ./...`, `go test -tags=integration ./...` and `go test
-count=1 -tags=e2e ./...` pass, every package `ok`.

## The fix of round 1 (2026-10-11)

Finding 1: the session of the review got no `problem-statement.md`. The
world's script now refuses a session of a phase with no problem statement in its
input (exit 3), and `TestTheReviewOfTheAnswersAndTheFollowUp` went red first
(12:00Z):

```
--- FAIL: TestTheReviewOfTheAnswersAndTheFollowUp (5.02s)
    run_integration_test.go:1164: the restart with the answers: the step answers-review is fail: the session S-3426d66f: crash
    run_integration_test.go:1164: the restart with the answers: the steps [forge clone version lease probe rule-gaps meaning numbering spans draft completeness batch answers answers-review], want [for
    run_integration_test.go:1166: the detail of follow-up ""
    run_integration_test.go:1180: the follow-up asks again []; want the missing, the two unparsed and the open ["Q-001" "Q-002" "Q-003" "Q-004"]
    run_integration_test.go:1182: git show refs/heads/layup-records:intake/answers-review.tsv: exit status 128
```

Then the review's input holds `problem-statement.md` and its prompt names it;
green (12:01Z). Note 4: `follow-up` fails on a missing problem statement, as
`batch` does.

With the fix (12:01Z to 12:05Z): `go build ./...`, `go vet ./...`, `go test ./...`,
`go test -tags=integration ./...` and `go test -count=1 -tags=e2e ./...` pass,
every package `ok`.
