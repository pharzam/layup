# T-tetz — test runs

The records of the Scaffold (#270, row 50): the Go schemas of `markers`,
`setup-tree` and `forge-check` against their blocks, and the rules of their
readers. All times are UTC, 2026-10-13, read from the clock after each step.

## Red (04:17Z)

The first version of the tests before `scaffold.go` of either package, `go vet`
of the packages, each (the second pair with `-tags=integration`):

```
vet: internal/records/scaffold_test.go:40:16: undefined: ReadMarkers
vet: internal/run/scaffold_test.go:66:16: undefined: ReadSetupTree
vet: internal/records/records_integration_test.go:27:103: undefined: MarkersSchema
vet: internal/run/scaffold_integration_test.go:20:59: undefined: SetupTreeSchema
```

## Mutations

Each rule of `scaffold.go` of both packages mutated alone, on a copy of the tree
and never on the worktree: its condition joined with `&& false`, or a column
taken out of the list of those that never hold the empty value, so each mutant
builds; then the tests of its package (`go test ./internal/records/`, and
`go test ./internal/run/ -run 'SetupTree|ForgeCheck'`). A mutant that does not
build would be counted as not run, never as caught.

- **The first run** (04:19:29Z to 04:20:30Z): 37 mutants, 32 caught, **5 survived**:
  the empty value in `line`, `marker` and `session` of `markers`, in `step` of
  `setup-tree` and in `run` of `forge-check`. The form of each of these columns
  refuses the empty value too, at the same column, so a case that named only the
  column passed with the rule of the empty value taken out. The tests were
  written again so that each refusal names its line, its column and its reason
  (`refusedBy`), and a case was added for the valid row of a marker that held a
  tab.
- **The second run** (04:22:03Z to 04:23:20Z): 46 mutants, the 37 and 9 more of the
  line that an error names (`i + 2`, `Line: 1`): 46 of 46 caught.
- **The final run** on the finished tree (04:26:32Z to 04:28:45Z): 46 of 46 caught,
  each with its first failing line (each line cut at 170 characters). The guard `i >= len(forgeChecks)` of
  `CheckForgeChecks` is also the bound of an index: without it the ninth row
  panics, which the test reports as a failure.

```
== markers-required-all: caught scaffold_test.go:85: the empty value in file: read, want an error of line 2, column "file"
== markers-required-file: caught scaffold_test.go:85: the empty value in file: read, want an error of line 2, column "file"
== markers-required-line: caught scaffold_test.go:85: the empty value in line: line 2, column "line": a line is counted from 1, want line 2, column "line" and a reason with "never holds the empty value"
== markers-required-marker: caught scaffold_test.go:85: the empty value in marker: line 2, column "marker": a marker has its two angle quotes, want line 2, column "marker" and a reason with "never holds th
== markers-required-session: caught scaffold_test.go:85: the empty value in session: line 2, column "session": a session ID is S- and 8 lowercase hexadecimal characters, want line 2, column "session" and a
== markers-id-form: caught scaffold_test.go:87: an id with a letter that is not hexadecimal: read, want an error of line 2, column "id"
== markers-line-1: caught scaffold_test.go:88: a line 0: read, want an error of line 2, column "line"
== markers-open-quote: caught scaffold_test.go:89: a marker with no open quote: read, want an error of line 2, column "marker"
== markers-close-quote: caught scaffold_test.go:90: a marker with no close quote: read, want an error of line 2, column "marker"
== markers-value-start: caught scaffold_test.go:91: a value with no start and no end: read, want an error of line 2, column "start"
== markers-start-end: caught scaffold_test.go:93: an end with no start: read, want an error of line 2, column "end"
== markers-start-below-end: caught scaffold_test.go:94: start equal to end: read, want an error of line 2, column "end"
== markers-session-form: caught scaffold_test.go:96: a session that is not S- and 8 hexadecimal characters: read, want an error of line 2, column "session"
== markers-order-file: caught scaffold_test.go:97: a file that goes down: read, want an error of line 3, column "file"
== markers-order-line: caught scaffold_test.go:98: a line that goes down in a file: read, want an error of line 3, column "line"
== runid-form: caught start_test.go:172: a run ID of 15 characters: read, want an error
== tree-required-all: caught scaffold_test.go:81: the empty value in commit: read, want an error of line 2, column "commit"
== tree-required-commit: caught scaffold_test.go:81: the empty value in commit: read, want an error of line 2, column "commit"
== tree-required-tree: caught scaffold_test.go:81: the empty value in tree: read, want an error of line 2, column "tree"
== tree-required-step: caught scaffold_test.go:81: the empty value in step: line 2, column "step": the step is root exactly when the order is 0, want line 2, column "step" and a reason with "never hol
== tree-root-iff-order-0: caught scaffold_test.go:83: the root with the step of a commit: read, want an error of line 2, column "step"
== tree-step-form: caught scaffold_test.go:85: a commit of the step S03: read, want an error of line 3, column "step"
== tree-no-row: caught scaffold_test.go:89: a table with no row: read, want an error of line 1, column "order"
== tree-order-sequence: caught scaffold_test.go:87: a first row that is not the root: read, want an error of line 2, column "order"
== forge-required-all: caught scaffold_test.go:138: the empty value in result: read, want an error of line 2, column "result"
== forge-required-result: caught scaffold_test.go:138: the empty value in result: read, want an error of line 2, column "result"
== forge-required-run: caught scaffold_test.go:138: the empty value in run: line 2, column "run": a run ID is 16 lowercase hexadecimal characters, want line 2, column "run" and a reason with "never ho
== forge-required-at: caught scaffold_test.go:138: the empty value in at: read, want an error of line 2, column "at"
== forge-fail-detail: caught scaffold_test.go:140: a fail with no reason: read, want an error of line 5, column "detail"
== forge-pass-no-detail: caught scaffold_test.go:141: a pass of tree with a detail: read, want an error of line 2, column "detail"
== forge-pass-rules-exception: caught scaffold_test.go:128: a pass of rules with its extra types: line 3, column "detail": a pass has no detail, but a pass of a rules check, which may name its extra rule type
== forge-run-form: caught scaffold_test.go:145: a run 0123456789abcde: read, want an error of line 4, column "run"
== forge-too-many-rows: caught a panic in CheckForgeChecks: index out of range [8] with length 8
== forge-second-prefix: caught scaffold_test.go:153: the second check that is not a rules check: read, want an error of line 3, column "check"
== forge-second-branch: caught scaffold_test.go:152: the second check with no branch: read, want an error of line 3, column "check"
== forge-name-by-position: caught scaffold_test.go:150: a check that is not the first: read, want an error of line 2, column "check"
== forge-too-few-rows: caught scaffold_test.go:147: a table of seven checks: read, want an error of line 9, column "check"
== line-readtable: caught scaffold_test.go:81: the empty value in commit: line 1, column "commit": this column never holds the empty value, want line 2, column "commit" and a reason with "never ho
== line-markers-file: caught scaffold_test.go:97: a file that goes down: line 2, column "file": the rows are in the order of the stop table of S10: the file is not below the file before it, want line
== line-markers-line: caught scaffold_test.go:98: a line that goes down in a file: line 2, column "line": the rows are in the order of the stop table of S10: in a file the line is not below the line
== line-tree-no-row: caught scaffold_test.go:89: a table with no row: line 2, column "order": no row: the first row is the root commit, of order 0, want line 1, column "order" and a reason with "no
== line-tree-order: caught scaffold_test.go:87: a first row that is not the root: line 1, column "order": the order of this row is 0: 0 for the root commit, then 1, 2, ... in the order of the rows,
== line-forge-too-many: caught scaffold_test.go:148: a table of nine checks: line 9, column "check": the table has eight rows, one for each check, want line 10, column "check" and a reason with "eight
== line-forge-second: caught scaffold_test.go:152: the second check with no branch: line 2, column "check": the next check is rules:<default branch>, want line 3, column "check" and a reason with "th
== line-forge-position: caught scaffold_test.go:150: a check that is not the first: line 1, column "check": the next check is tree, want line 2, column "check" and a reason with "the next check is tree
== line-forge-few: caught scaffold_test.go:147: a table of seven checks: line 8, column "check": the next check is activity, want line 9, column "check" and a reason with "the next check is activi
46 caught, 0 survived, 0 not run, of 46
```

## Green (04:26:19Z to 04:33:46Z)

Each with exit 0, in this order; the time of a line is when its command ended. The
three runs of `go test` are with `-count=1`, then with `-tags=integration` and with
`-tags=e2e`. They read the working tree of 04:26Z, whose files are those of `a89030b`:
`1473df7` (the code and the tests) was committed at 04:29Z and `a89030b` (the
documents) at 04:30Z, with no file changed between; this record is the one file of
`89917d2` (04:34Z):

```
go build ./...: exit 0 (04:26:21Z)
go vet ./...: exit 0 (04:26:23Z)
gofmt -l internal cmd (empty): exit 0 (04:26:23Z)
go test ./...: exit 0 (04:26:36Z)
go test -tags=integration ./...: exit 0 (04:31:33Z)
go test -tags=e2e ./...: exit 0 (04:32:59Z)
adr-lint: exit 0 (04:33:00Z)
prd-lint: exit 0 (04:33:00Z)
link-lint: exit 0 (04:33:13Z)
run-discipline-tests: exit 0 (04:33:19Z)
nested-checkout-check: exit 0 (04:33:20Z)
setup-check: exit 0 (04:33:46Z)
git diff --check origin/main: exit 0 (04:33:46Z)
```
