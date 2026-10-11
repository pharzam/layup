# T-1u2v — test runs

The review of meaning (#218, row 43b).

## Red, the unit level (07:05Z, 2026-10-11)

`go test -run 'TestQuoteAt|TestKeepOfMeaning|TestSessionSchema' ./internal/run/`,
against stubs (`quoteAt` gave 0, 0, nil; `Keep` an error; `sessionSchema` the
block). Before the code, two expectations of the test itself were corrected,
the offsets of `safe.` (51 to 56, not 50 to 55) and the count of `s` (6, not
5), each checked with a byte search of the file. The `pre-commit` hook runs the
unit level, so the tests landed with the code.

```
--- FAIL: TestQuoteAt (0.00s)
    meaning_test.go:35: "# The problem": 0, 0, <nil>; want 0, 13
    meaning_test.go:35: "service is": 0, 0, <nil>; want 19, 29
    meaning_test.go:35: "safe.": 0, 0, <nil>; want 50, 55
    meaning_test.go:35: "The service is fast.": 0, 0, <nil>; want 15, 35
    meaning_test.go:37: "slow": <nil>; want it refused, "occurs nowhere"
    meaning_test.go:37: "fast": <nil>; want it refused, "occurs 2 times"
    meaning_test.go:37: "s": <nil>; want it refused, "occurs 5 times"
    meaning_test.go:37: "is\tfast": <nil>; want it refused, "holds a tab"
    meaning_test.go:37: "fast.\nIt": <nil>; want it refused, "holds a line break"
    meaning_test.go:37: "problem\r": <nil>; want it refused, "holds a line break"
    meaning_test.go:37: "": <nil>; want it refused, "is empty"
--- FAIL: TestKeepOfMeaning (0.00s)
    meaning_test.go:51: Keep: not built
--- FAIL: TestSessionSchema (0.00s)
    meaning_test.go:93: the column id of the session's file is a key
    meaning_test.go:97: the columns ["id" "kind" "start" "end" "quote" "question" "session"]; want kind quote question, and the block unchanged
```

## Red, the integration level (07:06Z to 07:08Z)

`go test -tags=integration -run 'TestTheStepMeaning|TestTheStepRuleGaps'
./internal/run/`, with the world of `run_integration_test.go` given its
harness, models and route (condition 1 of the plan review), before the step.
Two gaps of the test world came first, each fixed in the world, not in the
code:

1. The step `probe` failed at the clone of a session's `repo/`: "Remote branch
   main not found in upstream origin". The world's bare target was made with
   `git init --bare` and no `-b main`, so its `HEAD` named a branch that did
   not exist, and the restart's clone had no branch `main`; the forge's `HEAD`
   names its default branch, so the world now makes it with `-b main`.
2. The step `probe` failed at its comment: "Comment: status 404: Not Found".
   The world's fake GitHub took no posted comment; it now records them.

Then the red for the right reason, the restart ending at `rule-gaps` with no
step `meaning` (07:08Z):

```
--- FAIL: TestTheStepRuleGaps (1.89s)
    run_integration_test.go:568: the restart: the steps [forge clone version lease probe rule-gaps], want [forge clone version lease probe rule-gaps meaning]
    run_integration_test.go:570: the progress ["1/6 forge" "2/6 clone" "3/7 version" "4/7 lease" "5/7 probe" "6/7 rule-gaps"], want ["1/7 forge" "2/7 clone" "3/7 version" "4/7 lease" "5/7 probe" "6/7
    run_integration_test.go:586: the second restart's row {Name:probe Result:done Detail:probed and passed 0, probed and failed 0, skipped 1}, want rule-gaps done before
--- FAIL: TestTheStepMeaning (1.71s)
    run_integration_test.go:615: the restart: the steps [forge clone version lease probe rule-gaps], want [forge clone version lease probe rule-gaps meaning]
    run_integration_test.go:617: the detail of meaning "2 rule gaps, in intake/psb.tsv"
    run_integration_test.go:619: git show refs/heads/layup-records:intake/meaning.tsv: exit status 128
--- FAIL: TestTheStepMeaningFails (3.31s)
        run_integration_test.go:663: the last row {Name:rule-gaps Result:done Detail:3 rule gaps, in intake/psb.tsv}; want meaning fail with ": refused file: the quote \"fast\" occurs 2 times"
        run_integration_test.go:663: the last row {Name:rule-gaps Result:done Detail:3 rule gaps, in intake/psb.tsv}; want meaning fail with "pair: no admitted pair"
```

The block `run-steps`: the e2e run of the binary, before the enum gained
`meaning` (07:12Z): `layup: line 8, column "step": "meaning" is not a value of
enum(...|phase|rule-gaps)`, exit 2; the Go schema and the block of `run.md`
then changed in one edit, and `TestTheTextBlocksOfRunMd` compares them.

## Mutations (2026-10-11, 07:13Z to 07:15Z)

Each rule mutated alone, the unit tests and, when they pass, the two
integration tests run: 12 of 12 caught, each mutant built. A first run caught
11: the mutant that cleared the key of a kept column was equivalent for
`meaning`, whose one key `id` is dropped; the rule is now that a kept key stays
a key, with a case of the block `spans`, and that mutant is replaced by
`schema-drops-nothing`, whose first form did not build and was rewritten.

```
== empty-cut: caught meaning_test.go:38: "": the quote "" occurs 58 times; want it refused, "is empty"
== tab-cut: caught meaning_test.go:38: "is\tfast": the quote "is\tfast" occurs nowhere; want it refused, "holds a tab"
== break-cut: caught meaning_test.go:38: "fast.\nIt": <nil>; want it refused, "holds a line break"
== twice-kept: caught meaning_test.go:38: "fast": <nil>; want it refused, "occurs 2 times"
== offset-off-by-one: caught meaning_test.go:36: "# The problem": 0, 12, <nil>; want 0, 13
== quote-error-ignored: caught meaning_test.go:73: a quote that occurs twice: map[], row 1: column end: end is above start; want it refused, "the quote \"fast\" occurs 2 times"
== id-from-zero: caught meaning_test.go:52: Keep: line 2, column "id": the next ID is M-001
== session-dropped: caught meaning_test.go:58: the records ["intake/meaning.tsv"]:
== done-before-cut: caught run_integration_test.go:632: the second restart's last row {Name:meaning Result:done Detail:gaps of meaning: 1, in intake/meaning.tsv}, want meaning done before
== pair-fail-cut: caught run_integration_test.go:663: the last row {Name:meaning Result:fail Detail:no admitted pair in the routing register for the role and the tier}; want meaning fail with "pa
== session-error-ignored: caught run_integration_test.go:663: the last row {Name:meaning Result:fail Detail:line 1: no header row}; want meaning fail with ": refused file: row 1: the quote \"fast\" occur
== schema-drops-nothing: caught meaning_test.go:52: Keep: line 1, column "id": the header row has "kind" in its place
```

## Green

At `2bc7861` with the documents (07:16Z to 07:19Z): `go build ./...`, `go vet
./...`, `go test ./...`, `go test -tags=integration ./...` and `go test
-count=1 -tags=e2e ./...` pass, every package `ok`.
