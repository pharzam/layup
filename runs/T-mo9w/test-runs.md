# T-mo9w — test runs

The e2e of the steps of Intake (#226, row 48).

## Red (13:20Z, 2026-10-11)

`go test -count=1 -tags=e2e -run TestTheStepsOfIntake ./cmd/layup/`, the test
written against `main` with the fake forge serving a person's comment on the
Intake issue: the second restart fails at `answers-review`, as the world's
harness script answers the review prompt with a span table (the rows of the
table elided, the lines cut at 200 characters).

```
--- FAIL: TestTheStepsOfIntake (15.65s)
    intake_e2e_test.go:79: the second restart has no "\nanswers-review\tdone\t3 answers reviewed, 0 open, in intake/answers-review.tsv\n"
    intake_e2e_test.go:79: the second restart has no "\nfollow-up\tdone\t1 questions asked again, in questions.tsv and the comment of the follow-up\n"
    intake_e2e_test.go:83: the second restart: exit 1
        answers-review	fail	the session S-e2af1dfb: refused file: line 1, column "question": the header row has "start" in its place
```

Then the script answers the review prompt with `closed` for each answered
question of its input; green (13:21Z), with `TestRunNewThenRestart`.

## Mutations (2026-10-11, 13:21Z to 13:22Z)

Two mutations of the product (note 5 of the plan review), each with the binary
rebuilt by the test's own build: 2 of 2 caught, each mutant built.

```
== follow-up-asks-all: caught intake_e2e_test.go:79: the second restart has no "\nfollow-up\tdone\t1 questions asked again, in questions.tsv and the comment of the follow-up\n"
== follow-up-answers-records-nothing: caught intake_e2e_test.go:104: intake/answers.tsv of the records branch: <nil>; want "\nQ-005\tGo\t902\t303\t"
```

## Green

At `eb29d89`, with the documents in the tree and not yet committed (13:23Z to
13:28Z): `go build ./...`, `go vet ./...`, `go test ./...`, `go test
-tags=integration ./...` and `go test -count=1 -tags=e2e ./...` pass, every
package `ok`; `adr-lint`, `prd-lint`, `link-lint`, `setup-check` and
`run-discipline-tests` exit 0.
