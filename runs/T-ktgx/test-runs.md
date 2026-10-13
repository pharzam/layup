# T-ktgx — test runs

The setup in phase 2 (#273, row 52a): `Phase2` and `StepsPhase2` of `internal/setup`,
the five steps that differ from phase 1. All times are UTC, 2026-10-13, read from the
clock after each step.

## Red (06:14Z)

The tests before `phase2.go`, `go vet` of `internal/setup` with and without the tag
`integration`, each:

```
vet: internal/setup/phase2_test.go:23:15: undefined: Phase2
vet: internal/setup/phase2_test.go:23:15: undefined: Phase2
```

The one case that the task adds to `internal/verify` (`TestTheSourcesOfARecord`, a
`computed` row whose ref is `records:start/start.tsv`) passes at the base, as the plan
says: it guards a claim of `records.md` that was true before the task, and has no red.

## Mutations

Each rule of `phase2.go` mutated alone, and the two edges of the shared body of
`runS04` in `steps.go`, on a copy of the tree and never on the worktree: its condition
joined with `&& false`, or a value, a ref, an actor or an argument changed, so each
mutant builds; then `go test ./internal/setup/` (the unit level). A mutant that does not
build would be counted as not run, never as caught.

- **The first run** (06:16:32Z to 06:17:46Z): 36 mutants, 35 caught, **1 survived**,
  `phase1-gap-table-error`. The refactor of `runS04` moved the refusal of a gap table of
  another form, which phase 1 gives as an input error, into a function that both phases
  share, and no test of phase 1 pinned its kind. `TestS04OfPhase1RefusesABadGapTable`
  was added.
- **The final run** on the finished tree (06:19:19Z to 06:20:34Z): 36 of 36 caught, each
  with its first failing line (each line cut at 170 characters). Two are caught by a
  panic, which the test reports as a failure: without the stop of S01 an absent answer
  is read as `stack[1]`, and with `Unchanged` removed the test calls a nil function.

```
== s01-no-stack-stop: caught panic: runtime error: index out of range [1] with length 0
== s01-empty-answer: caught phase2_test.go:55: an answer of the empty value: fail []; want a stop with [{Step:S01 Question:S01-stack Ask:Which technology stack does the target use? Give a stack of t
== s01-stack-entry: caught phase2_test.go:102: a stack of no entry: done "S01-stack answered; the stack has a catalog entry"; want fail with "the stack rust has no entry in the catalog"
== s01-name-form: caught phase2_test.go:102: a name that is not OWNER/NAME: done "S01-stack answered; the stack has a catalog entry"; want fail with "name: widget is not OWNER/NAME"
== s01-visibility: caught phase2_test.go:102: a visibility of another word: done "S01-stack answered; the stack has a catalog entry"; want fail with "visibility: internal is not public or private"
== s01-baseline-url: caught phase2_test.go:102: a baseline that is not a URL: done "S01-stack answered; the stack has a catalog entry"; want fail with "baseline: git@github.com:pharzam/armature.git
== s01-stack-source: caught phase2_test.go:81: an answer that cites a fact: done [["stack" "go" "answer" "S01-stack"] ["name" "acme/widget" "computed" "records:start/start.tsv"] ["visibility" "publi
== s01-name-row: caught phase2_test.go:74: the stack and two Q- rows: done "S01-stack answered; the stack has a catalog entry" [["stack" "go" "answer" "S01-stack"] ["name" "acme/widget" "answer"
== s01-visibility-row: caught phase2_test.go:74: the stack and two Q- rows: done "S01-stack answered; the stack has a catalog entry" [["stack" "go" "answer" "S01-stack"] ["name" "acme/widget" "compute
== s01-baseline-row: caught phase2_test.go:74: a row of S01-name that names another one: done "S01-stack answered; the stack has a catalog entry" [["stack" "go" "answer" "S01-stack"] ["name" "acme/w
== s01-brief-row: caught phase2_test.go:74: the stack and two Q- rows: done "S01-stack answered; the stack has a catalog entry" [["stack" "go" "answer" "S01-stack"] ["name" "acme/widget" "compute
== s01-evidence: caught phase2_test.go:74: the stack and two Q- rows: done "every answer present; the stack has a catalog entry" [["stack" "go" "answer" "S01-stack"] ["name" "acme/widget" "compu
== s01-unchanged: caught panic: runtime error: invalid memory address or nil pointer dereference
== s01-run: caught phase2_test.go:55: no answer: input []; want a stop with [{Step:S01 Question:S01-stack Ask:Which technology stack does the target use? Give a stack of the catalog of LAYU
== s02-tree-row: caught phase2_test.go:133: S02: done "the commit and the tree of start.tsv" [["pin.source" "https://github.com/pharzam/armature" "computed" "records:start/start.tsv"] ["pin.comm
== s02-time-row: caught phase2_test.go:133: S02: done "the commit and the tree of start.tsv" [["pin.source" "https://github.com/pharzam/armature" "computed" "records:start/start.tsv"] ["pin.comm
== s02-source-ref: caught phase2_test.go:133: S02: done "the commit and the tree of start.tsv" [["pin.source" "https://github.com/pharzam/armature" "computed" "x"] ["pin.commit" "a95965534b14b0bf1
== s02-commit-row: caught phase2_test.go:133: S02: done "the commit and the tree of start.tsv" [["pin.source" "https://github.com/pharzam/armature" "computed" "records:start/start.tsv"] ["pin.comm
== s02-evidence: caught phase2_test.go:133: S02: done "the commit and the tree" [["pin.source" "https://github.com/pharzam/armature" "computed" "records:start/start.tsv"] ["pin.commit" "a9596553
== s02-actor: caught phase2_test.go:136: the actor of S02 is "layup-setup"; want layup-run
== s03-no-row: caught phase2_test.go:162: no row pin.tree: fail "the root tree 8ffb250afd584da8b418bc220fb6d72e802924ce is not pin.tree "; want fail with "the record has no row S02 pin.tree"
== s03-git-error: caught phase2_test.go:162: an error of git: fail "the root tree  is not pin.tree 8ffb250afd584da8b418bc220fb6d72e802924ce"; want fail with "the tree of main: rev-parse refs/head
== s03-tree: caught phase2_test.go:162: another tree: done "root tree = pin.tree"; want fail with "the root tree 2222222222222222222222222222222222222222 is not pin.tree 8ffb250afd584da8b418
== s03-ref: caught phase2_test.go:147: S03: fail "the tree of main: unknown revision HEAD^{tree}" false, calls ["rev-parse HEAD^{tree}"]; want done, no commit and the one read of the tree
== s03-commands: caught phase2_test.go:166: S03 gives no command in phase 2, and its two commands stay in phase 1
== s03-actor: caught phase2_test.go:252: S03: same code false, actor "layup-setup", hand-off false, reads []; want the code of phase 1 exactly when it is not one of ["S01" "S02" "S03" "S04"],
== s04-run: caught phase2_test.go:196: S04: input "the gap table of inputs/briefs/problem-statement.md: line 1: no header row" false [], 0 files; want done, the values [["pin.adr" "docs/adr
== s04-q-prefix: caught phase2_test.go:211: a Q- row of no question of the batch: done "checks pin and facts", 6 files; want fail and no file
== s04-q-batch: caught phase2_test.go:211: a Q- row of no question of the batch: done "checks pin and facts", 6 files; want fail and no file
== s04-stack-text: caught phase2_test.go:196: S04: done "checks pin and facts" true [["pin.adr" "docs/adr/0009-pin-the-baseline.md" "computed" "the next free number of docs/adr/"] ["answers.record
== s04-refusal-reason: caught phase2_test.go:211: a Q- row of no question of the batch: fail "the question Q-009 is unknown", 0 files; want fail and no file
== s15-actor: caught phase2_test.go:252: S15: same code true, actor "layup-setup", hand-off false, reads []; want the code of phase 1 exactly when it is not one of ["S01" "S02" "S03" "S04"],
== s15-commands: caught phase2_test.go:220: S15 in phase 2 has Records false and Commands true; in phase 1 true and true; want nil in phase 2 and both in phase 1
== s15-records: caught phase2_test.go:220: S15 in phase 2 has Records true and Commands false; in phase 1 true and true; want nil in phase 2 and both in phase 1
== shared-s04-refusal: caught phase2_test.go:211: a Q- row of no question of the batch: done "checks pin and facts", 6 files; want fail and no file
== phase1-gap-table-error: caught phase2_test.go:281: a gap table of another form: fail "the gap table of inputs/briefs/problem-statement.md: line 1, column \"line\": the header row has no field for this
36 caught, 0 survived, 0 not run, of 36
```

## Green (06:21:44Z to 06:28:16Z)

Each with exit 0, in this order; the time of a line is when its command ended. The three
runs of `go test` are with `-count=1`, then with `-tags=integration` and with
`-tags=e2e`:

```
go build ./...: exit 0 (06:21:47Z)
go vet ./...: exit 0 (06:21:48Z)
gofmt -l internal cmd (empty): exit 0 (06:21:48Z)
go test ./...: exit 0 (06:21:54Z)
go test -tags=integration ./...: exit 0 (06:26:12Z)
go test -tags=e2e ./...: exit 0 (06:27:29Z)
adr-lint: exit 0 (06:27:30Z)
prd-lint: exit 0 (06:27:31Z)
link-lint: exit 0 (06:27:45Z)
run-discipline-tests: exit 0 (06:27:51Z)
nested-checkout-check: exit 0 (06:27:52Z)
setup-check: exit 0 (06:28:16Z)
git diff --check origin/main: exit 0 (06:28:16Z)
```
