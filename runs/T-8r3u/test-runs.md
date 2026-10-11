# T-8r3u — test runs

The batch (#222, row 46).

## Red, the unit level (10:11Z, 2026-10-11)

`go test -run 'TestMergeQuestions|TestTheCommentOfTheBatch|TestBackquote'
./internal/run/`, against stubs (the merge gave no row, the comment the empty
text, `backquote` its input). The `pre-commit` hook runs the unit level, so the
tests landed with the code (the lines cut at 160 characters).

```
--- FAIL: TestMergeQuestions (0.00s)
    batch_test.go:80: the merge:
--- FAIL: TestTheCommentOfTheBatch (0.00s)
    batch_test.go:113: the comment has no "LAYUP Intake: the batch of questions\n":
    batch_test.go:113: the comment has no "\n## For the idea owner\n":
    batch_test.go:113: the comment has no "\nQ-001: Which stack?\n\n":
    batch_test.go:113: the comment has no "\nQ-002: What is TLA?\n`We use Go and TLA and XYZ.`\n":
    batch_test.go:113: the comment has no "\nQ-003: How fast?\n`It is fast.`\n":
    batch_test.go:113: the comment has no "\nQ-006: Which fact supports R-001?\n\n":
    batch_test.go:113: the comment has no "\nQ-007: \n":
    batch_test.go:113: the comment has no "\nQ-009: \n":
    batch_test.go:113: the comment has no "\n## For the Operator\n\nNo setup question in M2c.\n":
    batch_test.go:113: the comment has no "`Q-NNN: <answer>`":
    batch_test.go:117: the comment does not start with the header, or its blocks are out of order:
--- FAIL: TestBackquote (0.00s)
    batch_test.go:132: backquote("plain") = "plain", want "`plain`"
    batch_test.go:132: backquote("a`b") = "a`b", want "``a`b``"
    batch_test.go:132: backquote("`a") = "`a", want "`` `a ``"
    batch_test.go:132: backquote("a`") = "a`", want "`` a` ``"
    batch_test.go:132: backquote("a``b") = "a``b", want "```a``b```"
```

## Red, the integration level (10:12Z)

`go test -tags=integration -run
'TestTheStepBatch|TestABatchWithNoCommentIsPostedFromTheRecord'
./internal/run/`, with the world's fake keeping a posted comment, before the
step: each restart ends at `completeness` (the rows of the table elided).

```
--- FAIL: TestTheStepBatch (3.74s)
    run_integration_test.go:965: the restart: the steps [forge clone version lease probe rule-gaps meaning numbering spans draft completeness], want [forge clone version lease probe rule-gaps meaning
    run_integration_test.go:966: git show refs/heads/layup-records:questions.tsv: exit status 128
--- FAIL: TestABatchWithNoCommentIsPostedFromTheRecord (3.64s)
    run_integration_test.go:997: the restart: the steps [forge clone version lease probe rule-gaps meaning numbering spans draft completeness], want [forge clone version lease probe rule-gaps meaning
    run_integration_test.go:1002: git rev-parse refs/heads/layup-records:questions.tsv: exit status 128
```

## Mutations (2026-10-11, 10:14Z to 10:17Z)

Each rule mutated alone, the unit tests and, when they pass, the two
integration tests run. The first run (10:14Z to 10:15Z) caught 12 of 15:
`line-keeps-cr` was not caught (no line of a test ended in CR), so
`TestLineSpan` was added; `repost-cut` did not build (`rows` unused), was not
counted, and was rewritten (`repost-elsewhere`, the comment on another issue);
`app-not-checked` was not caught (no other comment on the Intake issue), so the
case of a restart with no batch now holds a person's comment that starts with
the header and another comment of the App. These two tests were written after
the code; their red is the catch of their mutants. The second run (10:16Z to
10:17Z): 15 of 15 caught, each mutant built.

```
== dedupe-cut: caught batch_test.go:80: the merge:
== dedupe-no-kind: caught batch_test.go:80: the merge:
== no-quote-deduped: caught batch_test.go:80: the merge:
== line-keeps-cr: caught batch_test.go:148: line 1: 0, 4, true; want 0, 3, true
== line-off-by-one: caught batch_test.go:80: the merge:
== coverage-no-span: caught batch_test.go:80: the merge:
== intent-cut: caught batch_test.go:80: the merge:
== order-coverage-first: caught batch_test.go:80: the merge:
== quote-line-cut: caught batch_test.go:113: the comment has no "\nQ-002: What is TLA?\n`We use Go and TLA and XYZ.`\n":
== operator-block-cut: caught batch_test.go:113: the comment has no "\n## For the Operator\n\nNo setup question in M2c.\n":
== backquote-space-cut: caught batch_test.go:132: backquote("`a") = "```a``", want "`` `a ``"
== backquote-short: caught batch_test.go:132: backquote("a`b") = "`a`b`", want "``a`b``"
== posted-ignored: caught run_integration_test.go:984: the second restart's row {Name:batch Result:done Detail:the comment of the batch posted from questions.tsv}, 2 comments; want batch done befo
== repost-elsewhere: caught run_integration_test.go:1016: the comment posted again ["another comment of the App\n"]; want the first's bytes ["LAYUP Intake: the batch of questions\n\n## For the idea
== app-not-checked: caught run_integration_test.go:1013: the row {Name:batch Result:done Detail:done before}; want batch done, posted from the record
```

## Green

At `4360955` with the documents (10:18Z to 10:22Z): `go build ./...`, `go vet
./...`, `go test ./...`, `go test -tags=integration ./...` and `go test
-count=1 -tags=e2e ./...` pass, every package `ok`.
