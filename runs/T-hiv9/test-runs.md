# T-hiv9 — test runs

The steps numbering and spans (#219, row 44).

## Red, the unit level (08:02Z, 2026-10-11)

`go test -run TestKeepOfNumbering ./internal/run/`, against a stub `Keep` that
gave an error. The `pre-commit` hook runs the unit level, so the tests landed
with the code.

```
--- FAIL: TestKeepOfNumbering (0.00s)
    numbering_test.go:19: Keep: not built
```

## Red, the integration level (08:03Z)

`go test -tags=integration -run
'TestTheStepsNumberingAndSpans|TestAFailedSpanCheck|TestTheStepMeaning$'
./internal/run/`, with the worlds' scripts and routes, before the two steps:
each restart ends at `meaning` (the rows of the table elided).

```
--- FAIL: TestTheStepMeaning (2.35s)
    run_integration_test.go:625: the restart: the steps [forge clone version lease probe rule-gaps meaning], want [forge clone version lease probe rule-gaps meaning numbering spans]
--- FAIL: TestTheStepsNumberingAndSpans (2.22s)
    run_integration_test.go:710: the restart: the steps [forge clone version lease probe rule-gaps meaning], want [forge clone version lease probe rule-gaps meaning numbering spans]
    run_integration_test.go:712: the details "" and ""
    run_integration_test.go:715: git show refs/heads/layup-records:spec/spans.tsv: exit status 128
--- FAIL: TestAFailedSpanCheckRunsNumberingAgain (2.16s)
    run_integration_test.go:751: the last row {Name:meaning Result:done Detail:0 gaps of meaning, in intake/meaning.tsv}; want spans fail, the first fail row of the check
```

## Mutations (2026-10-11, 08:05Z to 08:10Z)

Each rule mutated alone, the unit test and, when it passes, the integration
tests of the steps run. The first run (08:05Z to 08:07Z) caught 8 of 12:

- `text-not-field-rule` was not caught, as the writer of the table applies the
  field rule of the column `text` itself; the step's own replace was taken out,
  so the rule has one home, and the mutant with it.
- `rerun-cut-numbering` and `rerun-cut-spans` did not build (`failed` unused),
  were not counted, and were rewritten (`&& (!failed || true)`).
- `pair-fail-cut` was not caught: no test ran `numbering` with no route for
  `product-owner`; `TestNumberingWithNoPair` was added.

The second run (08:07Z to 08:09Z) caught 11 of 11. Its mutant `done-before-cut-spans`
showed a defect: a rerun whose table equals the last `spec/check.tsv` commits
nothing new, and git refused the empty commit, so `spans` failed with the error
of git. A case of `TestAFailedSpanCheckRunsNumberingAgain` (a second rerun with
the same bad table) was red first (08:09Z):

```
--- FAIL: TestAFailedSpanCheckRunsNumberingAgain (4.74s)
    run_integration_test.go:763: a rerun with the same table: the last row {Name:spans Result:fail Detail:git commit -m layup run: the span check of the facts: exit status 1}; want spans fail as befor
```

The step now keeps an equal table as it is (green at 08:10Z), and its mutant
`same-table-committed` is caught: 12 of 12.

```
== text-clamp-cut: caught numbering_test.go:32: a span partly out of the file: "", line 3, column "text": this column never holds the empty value; want its text in the file
== block-check-cut: caught numbering_test.go:41: a span with no byte in the file: map["spec/spans.tsv":"start\tend\tfact\tclass\treason\ttext\n20\t25\t1\tneed\t—\t—\n"], <nil>; want it refused, "te
== rerun-cut-numbering: caught run_integration_test.go:761: the restart after the fail: the step spans is fail: git commit -m layup run: the span check of the facts: exit status 1
== done-before-cut-numbering: caught run_integration_test.go:734: the second restart's rows {Name:numbering Result:done Detail:2 spans, 2 facts, in spec/spans.tsv} and {Name:spans Result:done Detail:done bef
== rerun-cut-spans: caught run_integration_test.go:766: spec/check.tsv after the rerun:
== done-before-cut-spans: caught run_integration_test.go:734: the second restart's rows {Name:numbering Result:done Detail:done before} and {Name:spans Result:fail Detail:git commit -m layup run: the spa
== check-not-committed: caught run_integration_test.go:725: git show refs/heads/layup-records:spec/check.tsv: exit status 128
== fail-cut: caught run_integration_test.go:751: the last row {Name:spans Result:done Detail:the span check passed, in spec/check.tsv}; want spans fail, the first fail row of the check
== checkFailed-false: caught run_integration_test.go:763: numbering after a failed check: "done before"; want a new session's spans
== facts-count-cut: caught run_integration_test.go:712: the details "2 spans, 0 facts, in spec/spans.tsv" and "the span check passed, in spec/check.tsv"
== pair-fail-cut: caught run_integration_test.go:782: the last row {Name:numbering Result:fail Detail:no admitted pair in the routing register for the role and the tier}; want numbering fail (pai
== same-table-committed: caught run_integration_test.go:763: a rerun with the same table: the last row {Name:spans Result:fail Detail:git commit -m layup run: the span check of the facts: exit status 1}
```

## Green

At `f75feaa` (08:10Z to 08:13Z, before the gofmt of the test file, which
changed no code): `go build ./...`, `go vet ./...`, `go test ./...`, `go test
-tags=integration ./...` and `go test -count=1 -tags=e2e ./...` pass, every
package `ok`.
