# T-4c5c — test runs

The step markers (#276, row 53): the scan of the pinned tree, one session of `systems-architect`
that proposes values for the markers, `intake/markers.tsv`, and the setup questions of the batch
and its comment. All times are UTC, 2026-10-13, read from the clock after each step.

## Red (08:28:45Z to 08:28:48Z)

The tests before the code, each package alone. The questions of a setup row are refused for the
enum of `source`; `internal/setup` and `internal/run` do not build, for the functions that the
tests name:

```
== go test ./internal/records/ (the questions of M2d: the rows of a setup question)
    intake_test.go:77: questions: line 10, column "source": "setup" is not a value of enum(rule|meaning|coverage|intent|follow-up)
    intake_test.go:163: questions: a setup question whose ref is not S01-stack: line 2, column "source": "setup" is not a value of enum(rule|meaning|coverage|intent|follow-up), want an error in column ref
    intake_test.go:163: questions: a setup question whose ref is a gap of meaning's ID: line 2, column "source": "setup" is not a value of enum(rule|meaning|coverage|intent|follow-up), want an error in column r
    intake_test.go:163: questions: a setup question whose ref is a marker's key of 7 characters: line 2, column "source": "setup" is not a value of enum(rule|meaning|coverage|intent|follow-up), want an error in
    intake_test.go:163: questions: a setup question whose ref has an uppercase hexadecimal character: line 2, column "source": "setup" is not a value of enum(rule|meaning|coverage|intent|follow-up), want an err
    intake_test.go:163: questions: a setup question of a kind that is not setup: line 2, column "source": "setup" is not a value of enum(rule|meaning|coverage|intent|follow-up), want an error in column kind
    intake_test.go:163: questions: S01-stack not asked by layup run: line 2, column "source": "setup" is not a value of enum(rule|meaning|coverage|intent|follow-up), want an error in column asker
    intake_test.go:163: questions: a marker's question not asked by a session: line 2, column "source": "setup" is not a value of enum(rule|meaning|coverage|intent|follow-up), want an error in column asker
    intake_test.go:163: questions: S01-stack with a quote: line 2, column "source": "setup" is not a value of enum(rule|meaning|coverage|intent|follow-up), want an error in column start
    intake_test.go:163: questions: a setup question in the idea owner's block: line 2, column "source": "setup" is not a value of enum(rule|meaning|coverage|intent|follow-up), want an error in column block
== go vet ./internal/setup/
vet: internal/setup/scaffold_test.go:577:12: undefined: MarkerPlace
== go vet ./internal/run/ and go vet -tags=integration ./internal/run/
vet: internal/run/markers_test.go:31:29: undefined: setup.MarkerPlace
vet: internal/run/markers_test.go:31:29: undefined: setup.MarkerPlace
== go test ./internal/run/ (the package does not build; the list of the steps of a restart is asserted by TestStartAndRestartCallStepAtTheStartOfEachStep)
# github.com/pharzam/layup/internal/run [github.com/pharzam/layup/internal/run.test]
internal/run/markers_test.go:31:29: undefined: setup.MarkerPlace
internal/run/batch_test.go:42:3: unknown field Markers in struct literal of type batchRecords
```

## Mutations

Each rule of the new code mutated alone, on a copy of the tree and never on the worktree: its
condition joined with `&& false`, or a value, a ref, an actor or an argument changed, so that
each mutant builds; then the tests that the plan names for it (`go test` of the package, or
`-tags=integration` with `-run` for the scan and the step). A mutant that does not build is
counted as not run, never as caught.

- **The first run** (08:52:00Z to 08:58:35Z): 61 mutants, 54 caught, **0 survived**, and 7 not
  run: they did not build (an unused variable or import after the change).
- **The seven, repaired** (08:59:28Z to 09:00:27Z), each alone: each caught. The mutant that gives
  the session no list (`step-list-input`) is caught because the scripted harness of the world
  now reads the first line of the list that the session is given and exits with 4 when it is not
  the header of the block, as it did not before.
- **The final table** (each line cut at 150 characters), 61 of 61:

```
== rec-enum-setup: caught --- FAIL: TestValidIntakeRecordsAreRead (0.00s)
== rec-ref-stack: caught --- FAIL: TestValidIntakeRecordsAreRead (0.00s)
== rec-ref-marker-form: caught --- FAIL: TestValidIntakeRecordsAreRead (0.00s)
== rec-ref-any: caught --- FAIL: TestIntakeRecordsRefuseEachBrokenRule (0.00s)
== rec-kind: caught --- FAIL: TestIntakeRecordsRefuseEachBrokenRule (0.00s)
== rec-asker-stack: caught --- FAIL: TestIntakeRecordsRefuseEachBrokenRule (0.00s)
== rec-asker-marker: caught --- FAIL: TestIntakeRecordsRefuseEachBrokenRule (0.00s)
== rec-no-quote-stack: caught --- FAIL: TestIntakeRecordsRefuseEachBrokenRule (0.00s)
== rec-block: caught --- FAIL: TestIntakeRecordsRefuseEachBrokenRule (0.00s)
== set-id: caught --- FAIL: TestMarkerPlaces (0.00s)
== set-line: caught --- FAIL: TestMarkerPlaces (0.00s)
== set-marker: caught --- FAIL: TestMarkerPlaces (0.00s)
== list-line: caught --- FAIL: TestTheListOfMarkers (0.00s)
== list-marker: caught --- FAIL: TestTheListOfMarkers (0.00s)
== keep-key-lacks: caught --- FAIL: TestKeepOfMarkersRefusesEachBrokenRule (0.00s)
== keep-empty-value: caught --- FAIL: TestKeepOfMarkersRefusesEachBrokenRule (0.00s)
== keep-quote: caught --- FAIL: TestKeepOfMarkersRefusesEachBrokenRule (0.00s)
== keep-session: caught --- FAIL: TestKeepOfMarkers (0.00s)
== keep-offsets: caught --- FAIL: TestKeepOfMarkers (0.00s)
== keep-value: caught --- FAIL: TestKeepOfMarkers (0.00s)
== keep-final-read: caught --- FAIL: TestKeepOfMarkersReadsTheTableByItsBlock (0.00s)
== keep-scan-order: caught --- FAIL: TestKeepOfMarkers (0.00s)
== text-list: caught --- FAIL: TestTheTextOfTheMarkersPrompt (0.00s)
== text-header: caught --- FAIL: TestTheTextOfTheMarkersPrompt (0.00s)
== pinned-one-root: caught --- FAIL: TestThePinnedTreeOfTheMarkers (0.25s)
== pinned-pin-tree: caught --- FAIL: TestThePinnedTreeOfTheMarkers (0.28s)
== pinned-root-tree: caught --- FAIL: TestThePinnedTreeOfTheMarkers (0.16s)
== pinned-remove-worktree: caught --- FAIL: TestThePinnedTreeOfTheMarkers (0.25s)
== pinned-remove-dir: caught --- FAIL: TestThePinnedTreeOfTheMarkers (0.25s)
== step-done-before: caught --- FAIL: TestTheStepMarkers (8.06s)
== step-no-place: caught --- FAIL: TestTheStepMarkersWithNoMarker (6.18s)
== step-role: caught --- FAIL: TestTheStepMarkersFails (6.57s)
== step-pair-text: caught --- FAIL: TestTheStepMarkersFails (5.57s)
== step-proposals: caught --- FAIL: TestTheStepMarkers (7.76s)
== step-list-input: caught --- FAIL: TestTheStepMarkers (6.67s)
== merge-stack-text: caught --- FAIL: TestMergeQuestions (0.00s)
== merge-stack-asker: caught --- FAIL: TestMergeQuestions (0.00s)
== merge-stack-absent: caught --- FAIL: TestMergeQuestions (0.00s)
== merge-marker-quote: caught --- FAIL: TestMergeQuestions (0.00s)
== merge-marker-text: caught --- FAIL: TestMergeQuestions (0.00s)
== merge-marker-asker: caught --- FAIL: TestMergeQuestions (0.00s)
== merge-block: caught --- FAIL: TestMergeQuestions (0.00s)
== merge-kind: caught --- FAIL: TestMergeQuestions (0.00s)
== merge-source: caught --- FAIL: TestMergeQuestions (0.00s)
== merge-never-dropped: caught --- FAIL: TestMergeQuestions (0.00s)
== comment-block-rows: caught --- FAIL: TestTheCommentOfTheBatch (0.00s)
== comment-empty-text: caught --- FAIL: TestTheBlocksOfAComment (0.00s)
== comment-proposed: caught --- FAIL: TestTheCommentOfTheBatch (0.00s)
== comment-forms-flag: caught --- FAIL: TestTheCommentOfTheBatch (0.00s)
== comment-forms-always: caught --- FAIL: TestTheBlocksOfAComment (0.00s)
== about-stack: caught --- FAIL: TestTheBlocksOfAComment (0.00s)
== about-value: caught --- FAIL: TestTheCommentOfTheBatch (0.00s)
== batch-held-markers: caught --- FAIL: TestABatchWithNoCommentIsPostedFromTheRecord (7.43s)
== batch-read-markers: caught --- FAIL: TestTheStepBatchAsksTheSetupQuestions (31.07s)
== batch-comment-markers: caught --- FAIL: TestTheStepBatchAsksTheSetupQuestions (46.10s)
== followup-markers: caught --- FAIL: TestTheFollowUpAsksTheMarkersAgain (25.00s)
== followup-held-markers: caught --- FAIL: TestTheFollowUpAsksTheMarkersAgain (13.96s)
== followup-block: caught --- FAIL: TestTheFollowUpOfASetupQuestion (0.00s)
== names-markers: caught --- FAIL: TestStartAndRestartCallStepAtTheStartOfEachStep (0.00s)
== tail-markers: caught --- FAIL: TestTheStepMarkers (9.70s)
== enum-markers: caught --- FAIL: TestTheTextBlocksOfRunMd (0.02s)
61 caught, 0 survived, 0 not run, of 61
```


## Green (09:02:10Z to 09:10:43Z, head `cf2fc05`)

Each with exit 0, in this order; the time of a line is when its command ended. The three runs of
`go test` are with `-count=1`, then with `-tags=integration` and with `-tags=e2e`, and
`git diff --check` reads the commits against `origin/main`:

```
go build ./...: exit 0 (09:02:12Z)
go vet ./...: exit 0 (09:02:12Z)
gofmt -l internal cmd (empty): exit 0 (09:02:12Z)
go test ./...: exit 0 (09:02:18Z)
go test -tags=integration ./...: exit 0 (09:08:21Z)
go test -tags=e2e ./...: exit 0 (09:09:50Z)
adr-lint: exit 0 (09:09:52Z)
prd-lint: exit 0 (09:09:52Z)
link-lint: exit 0 (09:10:06Z)
run-discipline-tests: exit 0 (09:10:13Z)
nested-checkout-check: exit 0 (09:10:14Z)
setup-check: exit 0 (09:10:43Z)
git diff --check origin/main HEAD: exit 0 (09:10:43Z)
```
