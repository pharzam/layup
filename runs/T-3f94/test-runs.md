# T-3f94 - the test runs of row 66h

## S0 - the reconciliation of the plan (6249349020), the review (6249724875) and amendment 1 (6267136184) with the landed code

Read at `origin/main` = `652db13` (the branch had no work; the merge of `origin/main` was a fast-forward). Each name was checked in the tree.

| Name of the plan, the review or the amendment | Landed form | Note |
| --------------------------------------------- | ----------- | ---- |
| 66r: `AdvanceShapeTask`, `ShapeTaskInput`, `Advance` | `shapeloop.go:459`, `:22`, `:36` | `Advance` had no `Head`: added, set at `posVerified` from the bound head (amendment 1) |
| 66r: `openShapeTask(ctx, kind)`, `finalID`, `finalIDs` | `shapetask.go:77`, `:41`, `:54` | called, not rebuilt; the constraints use `finalID` (condition 3) |
| 66r: `drawTaskID`, `Sessions.taskID` | `probe.go:163`, `:182` | not touched (condition 1) |
| 66r: `loopPosition`, `readOnlyStep`, `startSession` | `shapeloop.go:254`, `:761`, `:517` | `runStage` passed `scratch` `false` for the plan, the plan review and the developer: now `true` (amendment 1) |
| 66r: scratch base through `TaskSpec.Branch` and `ProbeSpec.Branch` | `session.go` (`Branch`), `startSession` | the probe already takes the branch, so the optional change of `probe.go` of the amendment is not needed |
| 66f: `PostStatuses`, `HeadInput`, `holdsLease` | `statuses.go:202`, `handoffpost.go:418` | used through 66r's driver; `holdsLease` by the posts of this row |
| 66p: `filesAt` | `treeitems.go:30` | the decision records of `docs/adr/` |
| 62a: `ReadOptions`, `Standing`, `Filter`, `IDs`, `FactsOfSpans`, `RequirementsOfDraft` | `internal/spec/shape.go:47`, `:119`, `:187`, `:147`, `:75`, `:93` | as planned |
| plan: `PhaseSession` with the step `shape-architecture` | `phase.go:97` | the step of a phase is a label: `options` and `options-rewrite` (the name `shape-architecture` is the architecture task's step) |
| plan: the output of the options session through `Files`/`KeepFiles` | `session.ReadOutputFiles` needs the names before the session | **diverges** (condition 5): one file `options.tsv`, read by `ReadOutput`, split by code |
| plan: `recordFollowUpAnswers`, an extraction of `followUpAnswersStep` | not made | Shape has its own recorder (`shapeQuestions`), smaller; the two steps of the follow-up of Intake ignore the rows of kind `shape` (`intakeFollowUps`), which the plan did not see: a Shape row is a row of source `follow-up`, and without the filter `follow-up-answers` waited for it |
| plan: `scaffoldedWorld(t)` extracted from the forge-check restart test | not made | the real `Restart` is covered by the forge-check test (23 rows) and by the binary's Scaffold test (e2e); the step is run on 66r's `drvWorld` |
| `verify.InvariantIDs` | new, `target.go` | the matcher `guardrailsEntry`; each ID once |
| the kind `shape` of a follow-up row | `records/intake.go` `CheckQuestion` | one clause |

## S3 - the red runs, against the inert declarations (`91f5a2d`, `ba65fef`) with the final tests

Each command was run directly: the exit status is the one `go` gave. The long lines are cut at 210 characters and the repeated lines of one test are shown once.

### Unit level: `go test -count=1 ./internal/run/ ./internal/records/ ./internal/verify/`
```
start 2026-10-19T16:54:02Z
--- FAIL: TestStartAndRestartCallStepAtTheStartOfEachStep (0.00s)
    command_test.go:52: the steps of a restart with the Scaffold done: forge-check between lease and probe, and shape last: ["forge" "clone" "version" "lease" "forge-check" "probe" "rule-gaps" "meaning" "number
--- FAIL: TestTheConstraintsThatApplyAreListedByCodeFromTheRecordsAndTheTarget (0.00s)
    options_test.go:57: constraintsOf = []; want [{NFR-002 p99 under 200 ms} {F-2 } {F-3 } {Inv-1 } {Inv-2 } {ADR-0001 } {ADR-0003 } {ADR-0007 }] (the final ID of the non-functional row, the facts of class cons
    options_test.go:62: a record with no status line: err = <nil>; want "docs/adr/0004-x.md has no status line"
    options_test.go:62: a status heading with no line after it: err = <nil>; want "docs/adr/0004-x.md has no status line"
--- FAIL: TestTheOutputOfTheOptionsSessionIsCheckedBeforeItIsKept (0.00s)
    options_test.go:70: splitOptions(two decisions) = map[], <nil>; want the two tables, byte for byte
--- FAIL: TestTheRoundOfADecisionIsReadFromTheRecords (0.00s)
    options_test.go:90: optionTables = map[], map[]; want db and cache as current and db as round 1
--- FAIL: TestEveryOptionDroppedAsksOneRewriteAndASecondDropStops (0.00s)
    options_test.go:101: filterAll = map[]; want pg the survivor of db and none of cache
--- FAIL: TestTheCommentWithNoSurvivorNamesEachDroppedOptionWithItsRows (0.00s)
    options_test.go:136: noSurvivorComment =
--- FAIL: TestTheRewriteMustHoldExactlyTheDecisionsItReplaces (0.00s)
    options_test.go:144: rewriteFiles = map[], <nil>; want the new table and the table of round 1
    options_test.go:148: one option: err = nil; want a refusal
    options_test.go:148: a decision added: err = nil; want a refusal
    options_test.go:148: a decision omitted: err = nil; want a refusal
--- FAIL: TestThePromptsListTheConstraintsAndTheFiles (0.00s)
    options_test.go:159: the prompt of the options session lacks "The output file: /s/result/options.tsv\n":
    options_test.go:159: the prompt of the options session lacks "The draft: /s/input/draft.tsv\n":
    options_test.go:159: the prompt of the options session lacks "- NFR-002: p99 under 200 ms\n":
    options_test.go:159: the prompt of the options session lacks "- Inv-1\n":
    options_test.go:159: the prompt of the options session lacks "decision<TAB>option<TAB>constraint<TAB>result<TAB>evidence":
    options_test.go:159: the prompt of the options session lacks "at least two options":
    options_test.go:164: the prompt of the rewrite session lacks "The output file: /s/result/options.tsv\n":
    options_test.go:164: the prompt of the rewrite session lacks "The reasons: /s/input/reasons.txt\n":
    options_test.go:164: the prompt of the rewrite session lacks "- NFR-002: p99 under 200 ms\n":
    options_test.go:164: the prompt of the rewrite session lacks "cache: /s/input/cache.tsv\n":
    options_test.go:164: the prompt of the rewrite session lacks "decision<TAB>option":
--- FAIL: TestTheArchitectureTaskStartsAtTheVerifiedHeadOfTheSpecification (0.00s)
    shape_test.go:14: archStart(verified, no events) = "", "", <nil>; want 1111111111111111111111111111111111111111
    shape_test.go:17: archStart(two attempts at H) = "", "", <nil>; want 1111111111111111111111111111111111111111, no stop
    shape_test.go:21: archStart(an attempt at X) stop = "", <nil>; want "T-ar01: its base 2222222222222222222222222222222222222222 is not the verified head 1111111111111111111111111111111111111111 of T-sp01"
    shape_test.go:24: archStart(verified, no head) = nil; want an error
--- FAIL: TestTheQuestionsLeftOpenAreTheOpenAndTheUnansweredAsked (0.00s)
    shapequestions_test.go:37: leftOpen = []; want [Q-001 Q-002 Q-006] (the open and the unanswered asked; not the closed, the answered, the follow-ups, or one that Shape asks)
--- FAIL: TestAShapeRowAsksAgainInTheBlockOfItsQuestion (0.00s)
    shapequestions_test.go:50: shapeRows = []; want [[Q-003 follow-up Q-002 operator before-delivery shape   text of Q-002 2026-10-12T09:00:00Z layup run  asked] [Q-004 follow-up Q-001 idea-owner before-deliver
--- FAIL: TestTheCommentOfTheRoundIsTheFormOfTheBatch (0.00s)
    shapequestions_test.go:68: shapeComment = ""; want the comment of the follow-up under the header of Shape: "LAYUP Shape: the questions left open\n\n## For the idea owner\n\nNo question for the idea owner.\n
--- FAIL: TestOnlyTheRowsOfIntakeAreItsFollowUp (0.00s)
    shapequestions_test.go:75: intakeFollowUps = []; want the row Q-002 alone
--- FAIL: TestAnIntentAnswerThatParsesWritesBudgetWithShapesQuestionAsSource (0.00s)
    shapequestions_test.go:105: budget missing: []; want [[band.B 5.0 Q-005] [band.U 10.0 Q-006] [appetite.hours 40.0 Q-007]]
    shapequestions_test.go:109: budget missing: the state of Q-006 = asked; want answered
    shapequestions_test.go:109: budget missing: the state of Q-007 = asked; want answered
    shapequestions_test.go:109: budget missing: the state of Q-001 = open; want answered
    shapequestions_test.go:109: budget missing: the state of Q-002 = open; want answered
    shapequestions_test.go:109: budget missing: the state of Q-003 = open; want answered
    shapequestions_test.go:109: budget missing: the state of Q-005 = asked; want answered
    shapequestions_test.go:113: the times of Q-001 and Q-005 = "", ""; want the time of the answer
    shapequestions_test.go:127: a valid pair on a replace: [[band.B 5.0 Q-005] [band.U 10.0 Q-006] [appetite.hours 40.0 Q-007]]; want [[band.B 5.0 Q-005] [band.U 12.0 Q-008] [appetite.hours 40.0 Q-007]]
    shapequestions_test.go:130: a valid pair on a replace: the states map[Q-001:answered Q-002:answered Q-003:answered Q-004:answered Q-005:answered Q-006:answered Q-007:answered Q-008:asked]; want Q-002 and Q-
    shapequestions_test.go:134: a B not below U on a replace: [[band.B 5.0 Q-005] [band.U 10.0 Q-006] [appetite.hours 40.0 Q-007]], map[Q-001:answered Q-002:answered Q-003:answered Q-004:answered Q-005:answered
--- FAIL: TestAnUnparsedIntentAnswerLeavesTheQuestionOpenAndALaterCorrectionIsTheAnswerToANewRow (0.00s)
    shapequestions_test.go:143: an unparsed B: [], map[Q-001:open Q-002:open Q-003:open Q-004:answered Q-005:asked Q-006:asked Q-007:asked]; want no budget, Q-001 open, Q-005 answered, Q-002 answered
--- FAIL: TestAPartialAnswerLeavesTheOtherQuestionsAsked (0.00s)
    shapequestions_test.go:170: the state of Q-005 = asked; want answered
    shapequestions_test.go:170: the state of Q-001 = open; want answered
    shapequestions_test.go:177: the wait = []; want [Q-002 Q-003]
--- FAIL: TestOnlyTheBudgetFieldsBlockAndAShapeRowStillAskedBlocks (0.00s)
    shapequestions_test.go:185: intentWait = []; want [Q-001 Q-002]: an open field, and one that a Shape row asks; not success-criteria, not a question of another source
FAIL	github.com/pharzam/layup/internal/run	0.382s
--- FAIL: TestValidIntakeRecordsAreRead (0.00s)
    intake_test.go:80: questions: line 10, column "kind": kind is a kind of the question's source
FAIL	github.com/pharzam/layup/internal/records	0.061s
--- FAIL: TestTheInvariantsAreTheEntriesOfGuardrails (0.00s)
    target_test.go:253: the entries:
FAIL	github.com/pharzam/layup/internal/verify	0.038s
exit=1
end 2026-10-19T16:54:15Z
```

Green negatives at the stub, each beside its red control: in `internal/records`, the refusal of a follow-up of the kind `other` and of an intent question of the kind `shape` were refused before this task (the control is `TestValidIntakeRecordsAreRead`, red on the kind `shape`).

### Integration level: `TMPDIR=/home/layup/.cache/layup-session/tmp-green go test -count=1 -tags=integration -timeout 30m -run 'TestTheQuestionsLeftOpenAreAsked|TestTheRewriteAndTheRound|TestAStopDuringTheRewrite|TestTheStepShape|TestTheStepWaitsOnTheIntentForm|TestTheArchitectureHeadOfATarget|TestTheTextBlocksOfRunMd|TestTheStepForgeCheckChecks|TestValidIntakeRecordsAreRead' ./internal/run/ ./internal/records/`
```
start 2026-10-19T16:56:45Z
--- FAIL: TestTheStepForgeCheckChecksTheForgeAndCommitsItsTableAtEachRestart (77.06s)
    forgecheck_integration_test.go:266: the rows ["forge" "clone" "version" "lease" "forge-check" "probe" "rule-gaps" "meaning" "numbering" "spans" "draft" "completeness" "markers" "batch" "answers" "answers-re
    forgecheck_integration_test.go:281: the rows ["forge" "clone" "version" "lease" "forge-check"] and the notes of progress ["1/6 forge" "2/6 clone" "3/22 version" "4/22 lease" "5/22 forge-check"]; want the ro
    forgecheck_integration_test.go:289: the rows ["forge" "clone" "version" "lease" "forge-check"] and the notes of progress ["1/6 forge" "2/6 clone" "3/22 version" "4/22 lease" "5/22 forge-check"]; want the ro
    forgecheck_integration_test.go:298: the rows ["forge" "clone" "version" "lease" "forge-check"] and the notes of progress ["1/6 forge" "2/6 clone" "3/22 version" "4/22 lease" "5/22 forge-check"]; want the ro
--- FAIL: TestTheTextBlocksOfRunMd (0.01s)
    run_integration_test.go:668: schema run-steps: column 1 (step): the type: the block has "enum(forge|plan|baseline|root-push|read-back|records|issues|watch|clone|version|lease|probe|phase|rule-gaps|meaning|n
--- FAIL: TestTheQuestionsLeftOpenAreAskedAndTheIntentAnswerWritesBudget (1.81s)
    shape_integration_test.go:129: a post that fails: <nil>, 5 rows, 0 comments; want the error, the five Shape rows and no comment
--- FAIL: TestTheRewriteAndTheRoundWithNoProgress (1.82s)
    shape_integration_test.go:220: the step = {Name:shape Result: Detail:}; want {Name:shape Result:done Detail:waiting for the Operator: shape: no survivor}
--- FAIL: TestAStopDuringTheRewriteRunsItAgainOnce (1.86s)
    shape_integration_test.go:252: the rewrite that stops: wait "", <nil>, round-1 copy: false; want the error and no copy
--- FAIL: TestTheStepShape (1.94s)
    shape_integration_test.go:279: noPRD true: the step = {Name:shape Result: Detail:}; want done, with both tasks verified
--- FAIL: TestTheStepWaitsOnTheIntentFormAndFailsClosedWhenNothingIsOpen (1.95s)
    shape_integration_test.go:340: the step = {Name:shape Result: Detail:}; want {Name:shape Result:done Detail:waiting: the intent form: Q-001 Q-002 Q-003}
    shape_integration_test.go:347: the step with no budget.tsv and nothing open = {Name:shape Result: Detail:}; want the fail that names it
--- FAIL: TestTheArchitectureHeadOfATargetWithNoPRDPassesLayupSpec (24.44s)
    shapeloop_integration_test.go:602: noPRD true: the specification: {Task:T-lkrn Verified:true Head: Wait: Stop:}, <nil>; want a verified head that is its bound head
FAIL	github.com/pharzam/layup/internal/run	110.960s
--- FAIL: TestValidIntakeRecordsAreRead (0.00s)
    intake_test.go:80: questions: line 10, column "kind": kind is a kind of the question's source
FAIL	github.com/pharzam/layup/internal/records	0.018s
exit=1
end 2026-10-19T16:58:53Z
```

`TestTheArchitectureHeadOfATargetWithNoPRDPassesLayupSpec` runs its control (a target with a PRD, the architecture task at the default branch) first; it is green at the stub and the case with no PRD is the red one (its first assertion is the `Head` of the driver). The first form of `TestTheStepShape` stopped at a helper that reads `tasks.tsv` (not an assertion); it now asserts the row first.

### The end-to-end level: `go test -count=1 -tags=e2e -timeout 30m -run TestTheStepsOfTheScaffold ./cmd/layup/`
```
start 2026-10-19T16:58:53Z
--- FAIL: TestTheStepsOfTheScaffold (30.39s)
FAIL	github.com/pharzam/layup/cmd/layup	31.712s
exit=1
end 2026-10-19T16:59:28Z
```

## S9 - the compiled mutants of row 66h

Run 2026-10-19 from 17:41Z to 17:58Z by `mutants.py` (in the session's scratchpad) in a copy of the worktree at the committed green head `b653076`
(the copy also holds the documents and the one comment of `steps.go` that the next commit adds; no code differs). Each mutant is one edit of one
file, taken from a fresh backup (a copy made in a new directory at that mutant's start), built (`go build ./...`) and vetted
(`go vet -tags=integration` of `internal/run`, `internal/records` and `internal/verify`), then run with the test command of its row; a mutant that
does not build or vet counts as not run. The file is restored from the backup, its SHA-256 is checked to equal the hash before, and the same test
command is run on the restored file as the control (the column after the first failing assertion; 0 is green). The statuses are the exit
statuses of `go` read directly. `unit` is `go test -count=1 ./internal/run/ ./internal/records/ ./internal/verify/`; the others are
`TMPDIR=/home/layup/.cache/layup-session/tmp-green go test -count=1 -tags=integration -timeout 20m -run '<tests>' ./internal/run/`. The hashes are the
first 12 hexadecimal characters of the SHA-256 of the file before and during the mutation; the assertion is cut at 150 characters.

A first pass counted every mutant as `FAILED` at the build: the copy has no `.git`, and `go build` wanted the VCS status. No mutant had run; the pass was
repeated with `GOFLAGS=-buildvcs=false`. The first form of M32 left a variable unused and did not build (guardrails, "a mutant that does not build"); its row
is the second form, run alone at the end. Numbering: M1 to M7, M9 to M16 and M22 to M26 are the plan's; M17 to M21 went to row 66r with the driver; M8 is
"the facts of the spans are not read", since the inventory's `out` facts have no code in M2e (the rule is decided in `loop.md`, the reader is row 66j's);
M27 to M32 were added for the head given by the driver, the copy of the comment, the filter of Intake's steps, the kind `shape`, the field `success-criteria`
and the final ID.

| # | mutant | build | vet | test | first failing assertion | control | tests | sha256 before | sha256 during | restore |
| - | ------ | ----- | --- | ---- | ----------------------- | ------- | ----- | ------------- | ------------- | ------- |
| M1 | leftOpen takes closed too | ok | ok | 1 | TestTheQuestionsLeftOpenAreTheOpenAndTheUnansweredAsked: shapequestions_test.go:37: leftOpen = [Q-001 Q-002 Q-004 Q-006]; want [Q-001 Q-002 Q-006] (th | 0 | unit | 0275114802c2 | bc195c60744b | same |
| M2 | an asked question whose follow-up is answered is asked again | ok | ok | 1 | TestTheQuestionsLeftOpenAreTheOpenAndTheUnansweredAsked: shapequestions_test.go:37: leftOpen = [Q-001 Q-002 Q-003 Q-006]; want [Q-001 Q-002 Q-006] (th | 0 | unit | 0275114802c2 | 1abc9b711cea | same |
| M3 | a question that Shape already asks is asked again | ok | ok | 1 | TestTheQuestionsLeftOpenAreTheOpenAndTheUnansweredAsked: shapequestions_test.go:37: leftOpen = [Q-001 Q-002 Q-005 Q-006]; want [Q-001 Q-002 Q-006] (th | 0 | unit | 0275114802c2 | 5b7e4e84b85c | same |
| M4 | the source of budget.tsv is the batch's question (the Shape rows are not read) | ok | ok | 1 | TestAnIntentAnswerThatParsesWritesBudgetWithShapesQuestionAsSource: shapequestions_test.go:105: budget missing: []; want [[band.B 5.0 Q-005] [band.U 1 | 0 | unit | 0275114802c2 | d8caed567b39 | same |
| M5 | B not below U accepted on a replace | ok | ok | 1 | TestBudgetAfterFollowUp: followanswers_test.go:88: written, a budget of one row: "band.B  ; band.U  ; appetite.hours 50.0 Q-005", open []; want "appet | 0 | unit | 3ba07dee95c4 | fef5e6a4f08d | same |
| M6 | the wait is skipped for an open intent field | ok | ok | 1 | TestAnUnparsedIntentAnswerLeavesTheQuestionOpenAndALaterCorrectionIsTheAnswerToANewRow: shapequestions_test.go:146: the wait = []; want [Q-001] | 0 | unit | 0275114802c2 | 4948daaf4e9e | same |
| M7 | the question that a Shape row asks again keeps its state | ok | ok | 1 | TestAnIntentAnswerThatParsesWritesBudgetWithShapesQuestionAsSource: shapequestions_test.go:109: budget missing: the state of Q-001 = open; want answer | 0 | unit | 0275114802c2 | a4b6ed376010 | same |
| M8 | the facts of the spans are not read | ok | ok | 1 | TestTheConstraintsThatApplyAreListedByCodeFromTheRecordsAndTheTarget: options_test.go:57: constraintsOf = [{NFR-002 p99 under 200 ms} {Inv-1 } {Inv-2 | 0 | unit | 9a79d576c223 | 425fbd8f2222 | same |
| M9 | a record with no status is skipped, not an error | ok | ok | 1 | TestTheConstraintsThatApplyAreListedByCodeFromTheRecordsAndTheTarget: options_test.go:62: a record with no status line: err = <nil>; want "docs/adr/00 | 0 | unit | 9a79d576c223 | 58eb1c2548af | same |
| M10 | a superseded record applies | ok | ok | 1 | TestTheConstraintsThatApplyAreListedByCodeFromTheRecordsAndTheTarget: options_test.go:57: constraintsOf = [{NFR-002 p99 under 200 ms} {F-2 } {F-3 } {I | 0 | unit | 9a79d576c223 | e1d63dbf53eb | same |
| M11 | a table with one option is kept | ok | ok | 1 | TestTheOutputOfTheOptionsSessionIsCheckedBeforeItIsKept: options_test.go:82: one option: err = <nil>; want "decision db has fewer than two options" | 0 | unit | 9a79d576c223 | b06c1081afd0 | same |
| M12 | the table of round 1 is not kept | ok | ok | 1 | TestTheRewriteMustHoldExactlyTheDecisionsItReplaces: options_test.go:144: rewriteFiles = map[spec/options/cache.tsv:[111 112 116 105 111 110 9 99 111 | 0 | unit | 9a79d576c223 | 50a37dc3db95 | same |
| M13 | the rewrite runs a second time | ok | ok | 1 | TestEveryOptionDroppedAsksOneRewriteAndASecondDropStops: options_test.go:109: the rewrite dropped too: rewrite [cache], stop []; want no second rewrit | 0 | unit | 9a79d576c223 | 1f737a4753f5 | same |
| M14 | the rewrite is skipped (the stop at round 1) | ok | ok | 1 | TestEveryOptionDroppedAsksOneRewriteAndASecondDropStops: options_test.go:105: round 1, cache dropped: rewrite [], stop [cache]; want the rewrite of ca | 0 | unit | 9a79d576c223 | ba8cfef3501c | same |
| M15 | the comment with no survivor is posted without the read | ok | ok | 1 | TestTheRewriteAndTheRoundWithNoProgress: shape_integration_test.go:240: a restart: {Name:shape Result:done Detail:waiting for the Operator: shape: no | 0 | TestTheRewriteAndTheRound | 9a79d576c223 | 4103ff090174 | same |
| M16 | shape is done before with a decision that has no survivor | ok | ok | 1 | TestTheRewriteAndTheRoundWithNoProgress: shape_integration_test.go:240: a restart: {Name:shape Result:done Detail:done before}, head moved: false, 1 c | 0 | TestTheRewriteAndTheRound | 4bc3bddba925 | 1b07818825ca | same |
| M22 | shape is a row before the Scaffold is done | ok | ok | 1 | TestStartAndRestartCallStepAtTheStartOfEachStep: command_test.go:52: the steps of a restart with Start done and the Scaffold not done: forge-check las | 0 | unit | d444364e8531 | a90acf185efb | same |
| M23 | the architecture task starts at the run's base | ok | ok | 1 | TestTheStepShape: shape_integration_test.go:280: noPRD true: the step = {Name:shape Result:fail Detail:T-3cbu: layup/spec failed: covers: did not run: | 0 | TestTheStepShape | 4bc3bddba925 | ed4ba88189d7 | same |
| M24 | the base guard of the architecture task is skipped | ok | ok | 1 | TestTheArchitectureTaskStartsAtTheVerifiedHeadOfTheSpecification: shape_test.go:21: archStart(an attempt at X) stop = "", <nil>; want "T-ar01: its bas | 0 | unit | 4bc3bddba925 | 50f2f75cebb6 | same |
| M25 | runStage passes scratch false | ok | ok | 1 | TestTheArchitectureHeadOfATargetWithNoPRDPassesLayupSpec: shapeloop_integration_test.go:610: noPRD true: the architecture task: {Task: Verified:false | 0 | TestTheArchitectureHeadOfATarget | af88db2f9a9c | 48bc8abc5baa | same |
| M26 | the architecture task starts before the specification is verified | ok | ok | 1 | TestTheArchitectureTaskStartsAtTheVerifiedHeadOfTheSpecification: shape_test.go:28: archStart(the specification in a wait) = "", "", T-sp01: the task | 0 | unit | 4bc3bddba925 | 3eee15257d8b | same |
| M27 | the driver gives no head | ok | ok | 1 | TestTheArchitectureHeadOfATargetWithNoPRDPassesLayupSpec: shapeloop_integration_test.go:602: noPRD true: the specification: {Task:T-jcg6 Verified:true | 0 | TestTheArchitectureHeadOfATarget | af88db2f9a9c | d1fe88fd29cd | same |
| M28 | the comment of the questions is not copied | ok | ok | 1 | TestTheQuestionsLeftOpenAreAskedAndTheIntentAnswerWritesBudget: shape_integration_test.go:146: the restart: wait [Q-002 Q-003 Q-004], <nil>, 1 comment | 0 | TestTheQuestionsLeftOpenAreAsked | 0275114802c2 | afbc6045a13f | same |
| M29 | Intake's follow-up reads the Shape rows | ok | ok | 1 | TestOnlyTheRowsOfIntakeAreItsFollowUp: shapequestions_test.go:75: intakeFollowUps = [[Q-002 follow-up Q-001 idea-owner before-delivery follow-up   tex | 0 | unit | d268f8e6a98d | 08c3e127f8e6 | same |
| M30 | the kind shape is refused for a follow-up | ok | ok | 1 | TestAShapeRowAsksAgainInTheBlockOfItsQuestion: shapequestions_test.go:54: the row Q-003: column kind: kind is a kind of the question's source | 0 | unit | e00a8eee9ed0 | d58e16897862 | same |
| M31 | success-criteria blocks the step | ok | ok | 1 | TestOnlyTheBudgetFieldsBlockAndAShapeRowStillAskedBlocks: shapequestions_test.go:185: intentWait = [Q-001 Q-002 Q-004]; want [Q-001 Q-002]: an open fi | 0 | unit | 0275114802c2 | f0f698c1b26a | same |
| M32 | a requirement is named by its draft ID | ok | ok | 1 | TestTheConstraintsThatApplyAreListedByCodeFromTheRecordsAndTheTarget: options_test.go:57: constraintsOf = [{R-002 p99 under 200 ms} {F-2 } {F-3 } {Inv | 0 | unit | 9a79d576c223 | 10c797315341 | same |
