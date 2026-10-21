# T-4naq: the test runs

Task `T-4naq` (#428), row 66ab of M2e (working name K7). Times are `date -u` readings taken after each step; an exit status is the direct one of the command, never of a pipe.

## S0. The reconciliation

Main at the start: `ebbb581` (the merge of #455, K4), which holds K3 (`4259e05`) and K4. The branch `T-4naq` stood at the claim `d9852d2`; `git merge origin/main` was a fast-forward, so no merge commit was made.
Two things in the brief differ from the issue: it names the comments of #425 (K4's issue), whose IDs (6262367609, 6262491026, 6262589337) are the plan, the plan review and the Operator's decisions of #428 (this task), which were read.

| Cited by the plan (6262367609) | At `ebbb581` | Divergence |
| --- | --- | --- |
| `ReadBets` / `BetsSchema`, verdict `bet`, kind `first` | `internal/run/bets.go:14`, `:45` | none; `bet.go` reads the bets with it |
| the Intake class of a fact: a span fact from `spec/spans.tsv`, a `Q-NNN` by `factOfQuote`, else `context` | `answerFactClass` (`betbrief.go:196`) with `factRows` (`:219`), `factOfQuote` (`revision.go:34`); `(*state).factsOf` (`bet.go`) gives `[]spec.InventoryFact` with the Intake class | already extracted by 66u/66j: no extraction, K7 calls `factsOf` |
| `NextInventory` (`internal/spec/inventory.go:103`); version n > 1 refused | `inventory.go:103`; the bet wraps it as `nextVersion` (`rendered.go:46`) and refuses a later bet (`bet.go`, "a later bet's rendered task is row 66t's") | none |
| the precedent of a binding in the register commit | `archStage` (`archtasks.go:121`) with `openSpec.extra` | none |
| 66r's opening core `openTask(ctx, openSpec)` with the reuse-by-kind in `openShapeTask` | `shapetask.go:77`, `:87`, `:100`; `reuse` is a field of `openSpec`, false for the architecture tasks of the bet | none for the reuse. **One gap:** `openTask` draws its own ID (`s.taskID`), and K7's ID is reserved in `spec/reopenings.tsv` before the issue, so `openSpec` gains `id` |
| `taskID` / `drawTaskID` | `probe.go:182`, `:163`; draws avoid `sessions.tsv` and `tasks.tsv` | extended: draws also avoid `spec/reopenings.tsv`, and a batch of draws avoids itself (review condition 4) |
| `finalID` / `finalIDs` | `shapetask.go:40`, `:54` | none |
| `AdvanceShapeTask` | `shapeloop.go`; K7 calls none of it | none |
| `stepInstructions`, `inputsWanted`, `buildPrompt`, `PromptFile`, `taskInputs`, `factsRecords`, `nextFactsID` | `instructions.go:29`, `:158`; `promptfile.go:172`, `:437`, `:600`, `:312`, `:669` | lines moved only |
| `FinishTask`, `MergeTask`, `RecordMerged` | `finish.go:45`, `ready.go:195`, `:257` | **FinishTask treats every task of kind `specification` as the original one**: `finish.go:87` asks for the rendered commit of the first bet and `shapeloop.go:743` for the pinned `prd-lint` of a specification that waits for its bet. A reopened task is of kind `specification` and has no rendered commit; 66k tells the two apart with `reopenedTasks` (K7's) when it drives a reopened task. Not changed here (reported) |
| `registerTask` `register.go:267`, `registerFiles` `:213`, `AddTask` `:194`, `loopPermissionsStep` `:254`, `parameterFiles` `parameters.go:374`, `Fenced` `lease.go:118`, `specInput` `specinput.go:60`, `FullInput` `full.go:63`, `needs` `:494`, `measures` `:520`, the waits heading `full.go:352` | all at the lines in the column | none |
| K3 (`UpToDate` `baseup.go:43`, `TipOf` `:20`, the opening pin `pinPath` `redo.go:312`) | landed | K7 uses none. `redoIssue` (`redo.go:152`) composes a specification task's issue only from an opening pin; a reopened task has none (it is not `reuse`), so a redo of a reopened task is an error before any write: a known limit of M2e, stated in `loop.md` and tested |
| the `merged` event (`records.ReadEvents`) | `ready.go:257` writes it; column 5 is the merge commit, column 6 the pull request | none |

No step of the plan is impossible, and the plan, the review and the decisions settle every decision this task makes, with one exception that I record and settle by the plan review's own rule (condition 4, "an existing task row must agree with the reservation"): a restart that finds a register row of a reserved task ID compares its kind, requirements and predecessors with the reservation and stops on a difference.

## S3. The red run (against the stubs of `d491aa3`, before the code)

Every command was run alone, with its own status, output to a file and not through a pipe. The tests are in the working tree and uncommitted: the hook refuses a red commit (`guardrails.md` §2), so they are committed with the code. The first failing assertions are quoted from the outputs.

### unitrun

- Command: `TMPDIR=/home/layup/.cache/layup-session/tmp-green go test -count=1 ./internal/run/ -run '^(TestTheReclassedFactsAreTheVersionRowsWhoseClassDiffersFromIntake|TestTheNextRequirementIDIsOneMoreThanEveryKnownID|TestAReopeningsTableRefusesEachBrokenRule|TestTheRowsOfTheRecordAreBoundToTheSelection|TestALaterVersionIsRefused|TestTheExemptFactsFollowTheStatesOfTheTwoTasks|TestADrawAvoidsAnUnopenedReservationAndItsOwnBatch|TestTheReopenedSpecificationPromptIsItsGoldenFile|TestEveryStepOfAReopenedTaskSeesThePRDAndItsRow|TestTheReopenedPromptPinsTheNumberedRecord|TestTheReopenedFormPassesTheFullCheck|TestTheReopeningSetOfTheFullInputIsReadFromTheRecords)$'`
- Start to end: 23:34:37 to 23:34:37Z (`date -u`), exit 1.

First failing assertions, one test after the other:

```text
--- FAIL: TestTheReclassedFactsAreTheVersionRowsWhoseClassDiffersFromIntake (0.00s)
    reopen_test.go:41: facts = [], <nil>; want [3 7 12 Q-002] (a need from context, a constraint from need, a measure from context, a fact of the answers; not an equal class, an out-of-scope fact or a context)
    reopen_test.go:44: a fact of the version that Intake does not know: <nil>; want an error that names it
--- FAIL: TestTheNextRequirementIDIsOneMoreThanEveryKnownID (0.00s)
    reopen_test.go:62: all three sources: nextNumber = 0; want 6
    reopen_test.go:62: the draft alone: nextNumber = 0; want 9
    reopen_test.go:62: the register alone: nextNumber = 0; want 10
    reopen_test.go:62: the reopenings alone: nextNumber = 0; want 11
    reopen_test.go:62: none: nextNumber = 0; want 1
    reopen_test.go:70: ids = []; want [[1 3 need F-0002#3 REQ-006 T-aaaa 9] [1 7 constraint F-0002#7 NFR-007 T-bbbb 9] [1 12 measure F-0002#12  T-cccc 9] [1 Q-002 need Q-002 REQ-008 T-dddd 9]] (REQ for a need, NFR for a constraint
    reopen_test.go:73: the 1000th requirement = []; want REQ-1000
--- FAIL: TestAReopeningsTableRefusesEachBrokenRule (0.00s)
    reopen_test.go:83: the control: 0 rows, <nil>; want 2, a measure with the empty requirement
--- FAIL: TestTheRowsOfTheRecordAreBoundToTheSelection (0.00s)
    reopen_test.go:115: allocate gave 0 rows; want 2
--- FAIL: TestALaterVersionIsRefused (0.00s)
    reopen_test.go:164: the control: [], <nil>; want the row of version 1
--- FAIL: TestTheExemptFactsFollowTheStatesOfTheTwoTasks (0.00s)
    reopen_test.go:208: before the allocation, any head: [], <nil>; want [7 12]
    reopen_test.go:208: the head of A: [], <nil>; want [12]
    reopen_test.go:208: the head of B: [], <nil>; want [7]
    reopen_test.go:208: the head of another task: [], <nil>; want [7 12]
    reopen_test.go:208: after A has merged, another head: [], <nil>; want [12]
    reopen_test.go:208: a merged task of no fact: [], <nil>; want [7 12]
--- FAIL: TestADrawAvoidsAnUnopenedReservationAndItsOwnBatch (0.00s)
    reopen_test.go:223: taskID = "T-0000", <nil>; want T-1111, the draw after the one that the reservation holds
    reopen_test.go:227: taskIDs = [], <nil>; want T-1111, T-2222 and T-3333: a draw again for the one of the batch and for the reservation
--- FAIL: TestTheReopenedSpecificationPromptIsItsGoldenFile (0.00s)
    reopen_test.go:269: need: Text differs from its golden file: got 6005 bytes, want 4852; first difference at byte 578: got "/s/S-aaaaaaaa/input/draft.tsv 7eb2ca55b8", want "docs/prd/PRD-0001.md 937a536b787a2e0a725"
    reopen_test.go:269: measure: Text differs from its golden file: got 6449 bytes, want 5206; first difference at byte 2037: got "specification of the task T-ab12 (kind s", want "requirement that waits for a bet, for th"
--- FAIL: TestEveryStepOfAReopenedTaskSeesThePRDAndItsRow (0.00s)
    reopen_test.go:288: plan-of-a-task: the requirements = [spec/draft.tsv spec/spans.tsv start/problem-statement.md intake/answers.tsv spec/reopenings.tsv]; want the PRD with no row of REQ-006, the record of the reopenings and no
    reopen_test.go:288: plan-review: the requirements = [spec/draft.tsv spec/spans.tsv start/problem-statement.md intake/answers.tsv spec/reopenings.tsv]; want the PRD with no row of REQ-006, the record of the reopenings and no se
    reopen_test.go:288: specification: the requirements = [spec/draft.tsv spec/spans.tsv start/problem-statement.md intake/answers.tsv spec/reopenings.tsv]; want the PRD with no row of REQ-006, the record of the reopenings and no
    reopen_test.go:291: specification: holds the reopened instruction: false
    reopen_test.go:288: verification: the requirements = [spec/draft.tsv spec/spans.tsv start/problem-statement.md intake/answers.tsv spec/reopenings.tsv]; want the PRD with no row of REQ-006, the record of the reopenings and no s
--- FAIL: TestTheReopenedPromptPinsTheNumberedRecord (0.00s)
    reopen_test.go:321: another record in the token: <nil>; want an error holding "F-0005"
    reopen_test.go:321: no numbered record: <nil>; want an error holding "numbered facts record"
    reopen_test.go:321: two numbered records: <nil>; want an error holding "2 numbered facts records"
    reopen_test.go:325: the control: <nil>; want the line with the token and the numbered record
--- FAIL: TestTheReopenedFormPassesTheFullCheck (0.00s)
    reopen_test.go:367: the instruction shows no line for a need
--- FAIL: TestTheReopeningSetOfTheFullInputIsReadFromTheRecords (0.00s)
    reopen_test.go:437: no record yet: Reopening = [], <nil>; want [1]
    reopen_test.go:437: another task, the task of the fact not merged: Reopening = [], <nil>; want [1]
    reopen_test.go:443: a record the reader refuses: <nil>; want an error that names spec/reopenings.tsv
```

### unitspec

- Command: `TMPDIR=/home/layup/.cache/layup-session/tmp-green go test -count=1 ./internal/spec/ -run '^TestAReclassedFactIsExemptFromNeedsAndMeasuresWhileItIsListed$'`
- Start to end: 23:34:37 to 23:34:37Z (`date -u`), exit 1.

First failing assertions, one test after the other:

```text
--- FAIL: TestAReclassedFactIsExemptFromNeedsAndMeasuresWhileItIsListed (0.00s)
    full_test.go:772: the rule needs: result "fail", want "pass" (detail "the need fact 8 has no requirement")
    full_test.go:773: the rule measures: result "fail", want "pass" (detail "the measure fact 9 has no success-criterion row")
    full_test.go:776: the rule needs: result "fail", want "pass" (detail "the need fact 8 has no requirement")
    full_test.go:780: the rule measures: result "fail", want "pass" (detail "the measure fact 9 has no success-criterion row")
```

### integ

The selector of this run also named a test of `internal/tsv`, which the package argument did not reach, so it ran nothing and is left out of the command below.

- Command: `TMPDIR=/home/layup/.cache/layup-session/tmp-green go test -count=1 -tags=integration -timeout 30m ./internal/run/ -run '^(TestAReclassedFactOpensOneSpecificationTaskAfterTheSpecificationMerges|TestNothingOpensBeforeTheMergedEvent|TestARestartAfterTheAllocationOpensOneTaskWithTheSameIDs|TestAStopBetweenTheIssueAndTheRowOpensASecondIssueAndOneRow|TestTheReopeningIsNotTheSpecificationReusedByKind|TestTheRenderedHeadOfAReclassingBetPassesLayupSpec|TestALostLeaseOpensNoIssue|TestTheReopenedTasksPassTheFullCheckOneAfterTheOther|TestTheReopeningsBlockAndTheNewTextBlocksEqualTheCode|TestAHarnessThatFollowsTheReopenedInstructionGivesAValidHandoffAndAGreenSpecStatus|TestTheInstructionBlocksOfLoopMd)$'`
- Start to end: 23:34:37 to 23:35:01Z (`date -u`), exit 1.

First failing assertions, one test after the other:

```text
--- FAIL: TestARestartAfterTheAllocationOpensOneTaskWithTheSameIDs (4.66s)
    reopen_integration_test.go:147: the stop: <nil>, 0 rows, 0 issues, 3 tasks; want the stop with the allocation and nothing else
--- FAIL: TestAStopBetweenTheIssueAndTheRowOpensASecondIssueAndOneRow (4.80s)
    reopen_integration_test.go:171: the stop: <nil>, 22 issues, 3 tasks; want the stop after one issue, no row and the allocation
--- FAIL: TestALostLeaseOpensNoIssue (4.95s)
    reopen_integration_test.go:240: openReopenings = <nil>; 0 issues, 3 tasks; want a LostError, no issue and no task
--- FAIL: TestTheInstructionBlocksOfLoopMd (0.00s)
    promptfile_integration_test.go:509: loop.md has the block instruction-specification-reopened: true; it is equal to the code: false
--- FAIL: TestTheReopeningsBlockAndTheNewTextBlocksEqualTheCode (0.00s)
    reopen_integration_test.go:328: the block reopenings and ReopeningsSchema: schema reopenings: the number of columns: the block has "7"; the Go schema has "0"
    reopen_integration_test.go:338: loop.md has the block reopen-issue: true; it is equal to the code: false
    reopen_integration_test.go:338: loop.md has the block instruction-specification-reopened: true; it is equal to the code: false
--- FAIL: TestTheReopenedTasksPassTheFullCheckOneAfterTheOther (7.03s)
    reopen_integration_test.go:252: openReopenings = [], <nil>; want four tasks
--- FAIL: TestNothingOpensBeforeTheMergedEvent (7.06s)
    reopen_integration_test.go:134: the control, after merged: [], <nil>; want four tasks
--- FAIL: TestAReclassedFactOpensOneSpecificationTaskAfterTheSpecificationMerges (7.09s)
    reopen_integration_test.go:86: openReopenings = [], <nil>; want four tasks
--- FAIL: TestTheReopeningIsNotTheSpecificationReusedByKind (7.13s)
    reopen_integration_test.go:191: openReopenings = [], <nil>; want four tasks
--- FAIL: TestAHarnessThatFollowsTheReopenedInstructionGivesAValidHandoffAndAGreenSpecStatus (7.15s)
    reopen_integration_test.go:351: openReopenings = [], <nil>; want four tasks
--- FAIL: TestTheRenderedHeadOfAReclassingBetPassesLayupSpec (7.84s)
    reopen_integration_test.go:219: FinishTask of the specification = {Task:T-raoj Verified:false Merged:false Head: Attempt:0 Wait:waiting: T-raoj: not ready at bfa7eb39b13c41866e27417a86ed343d6fb9c4d4: layup/spec: failure Stop:
```

### Reading of the red

- Every unit test and every integration test above fails at the stub on its own assertion: an empty selection, the number 0, no rows, no exemption, the standard prompt in place of the golden file, no record, a stop that does not stop. Two are not red at the stub and are labelled so: `TestNothingOpensBeforeTheMergedEvent` is a **green negative** beside its red control (the same test after the event, and `TestAReclassedFactOpensOneSpecificationTaskAfterTheSpecificationMerges`, which fail on `openReopenings = [], <nil>; want four tasks`); the sub-check `a task with no row of reopenings` of `TestEveryStepOfAReopenedTaskSeesThePRDAndItsRow` (the standard prompt) is a characterization that passes at the stub and fails when the variant is selected for every task (mutant M11 of `mutants.md`).
- The red of `TestTheRenderedHeadOfAReclassingBetPassesLayupSpec` is the exemption itself: `FinishTask` waits with `layup/spec: failure` at the rendered head, the deadlock of Q1.
- `TestAHarnessThatFollowsTheReopenedInstructionGivesAValidHandoffAndAGreenSpecStatus` and `TestTheReopenedTasksPassTheFullCheckOneAfterTheOther` stop at their first step at the stub (`openReopenings = [], <nil>; want four tasks`), because both need the tasks first; their own assertions are reached by the green run, and the integration level alone kills the mutants M25, M42 and M12b-style states of the exemption at those heads (`mutants.md`).
- `TestTheInstructionBlocksOfLoopMd` and `TestTheReopeningsBlockAndTheNewTextBlocksEqualTheCode` fail on the empty constants and the schema with no column (`the block has "7"; the Go schema has "0"`).
- Two corrections of the tests were made after the first red run and before this one, none of the assertions: a guard against a panic on an empty result (a red must not abort the run), and the ADR files of the tree of the prompt test.

### Tests written after the code

`TestTheNumberedFactsRecordOfTheMergeCommitIsTheOneFile`, `TestARegisteredTaskThatDisagreesWithItsReservationIsRefused`, the register in `TestTheReopeningIsNotTheSpecificationReusedByKind` (a task of the register with the higher REQ-009) and the requirements assertion of the harness test (the PRD with no row and the record) were added after the first green run, to kill the mutants M37, M39, M5b and M23b, which no test killed before; they have no red run against the stubs and their red is the mutant of `mutants.md`. The Operator accepted this R8 exception (finding F5 of round 1) in the comment 6307600855 of #428.

## S4. Green, and the gates

The code is in `ba39f3f` and `c2f03bd`; the tests below ran on the head `7406949` (the later commits change records only). Each command alone, with its direct exit status.

### The tests of this task

- unitrun: 00:02:56 to 00:02:56Z, exit 0: `go test -count=1 ./internal/run/ -run '^(TestTheReclassedFactsAreTheVersionRowsWhoseClassDiffersFromIntake|TestTheNextRequirementIDIsOneMoreThanEveryKnownID|TestAReopeningsTableRefusesEachBrokenRule|TestTheRowsOfTheRecordAreBoundToTheSelection|TestALaterVersionIsRefused|TestTheExemptFactsFollowTheStatesOfTheTwoTasks|TestADrawAvoidsAnUnopenedReservationAndItsOwnBatch|TestTheReopenedSpecificationPromptIsItsGoldenFile|TestEveryStepOfAReopenedTaskSeesThePRDAndItsRow|TestTheReopenedPromptPinsTheNumberedRecord|TestTheReopenedFormPassesTheFullCheck|TestTheReopeningSetOfTheFullInputIsReadFromTheRecords|TestTheNumberedFactsRecordOfTheMergeCommitIsTheOneFile)$'`
- unitspec: 00:02:56 to 00:02:56Z, exit 0: `go test -count=1 ./internal/spec/ -run '^TestAReclassedFactIsExemptFromNeedsAndMeasuresWhileItIsListed$'`
- integ: 00:02:56 to 00:03:19Z, exit 0: `go test -count=1 -tags=integration -timeout 30m ./internal/run/ ./internal/tsv/ -run '^(TestAReclassedFactOpensOneSpecificationTaskAfterTheSpecificationMerges|TestNothingOpensBeforeTheMergedEvent|TestARestartAfterTheAllocationOpensOneTaskWithTheSameIDs|TestAStopBetweenTheIssueAndTheRowOpensASecondIssueAndOneRow|TestTheReopeningIsNotTheSpecificationReusedByKind|TestTheRenderedHeadOfAReclassingBetPassesLayupSpec|TestALostLeaseOpensNoIssue|TestTheReopenedTasksPassTheFullCheckOneAfterTheOther|TestTheReopeningsBlockAndTheNewTextBlocksEqualTheCode|TestAHarnessThatFollowsTheReopenedInstructionGivesAValidHandoffAndAGreenSpecStatus|TestTheInstructionBlocksOfLoopMd|TestEverySchemaBlockIsBuiltOrNotYetBuilt|TestARegisteredTaskThatDisagreesWithItsReservationIsRefused)$'`

### The gates

| Gate | Start to end | Exit | Command |
| ---- | ------------ | ---- | ------- |
| build | 23:57:45 to 23:57:46Z | 0 | `go build ./...` |
| vet | 23:57:46 to 23:57:46Z | 0 | `go vet ./...` |
| vetint | 23:57:46 to 23:57:46Z | 0 | `go vet -tags=integration ./...` |
| vete2e | 23:57:46 to 23:57:47Z | 0 | `go vet -tags=e2e ./...` |
| gofmt | 23:57:47 to 23:57:47Z | 0 | `sh -c out=$(gofmt -l .) && test -z "$out"` |
| unit | 23:57:47 to 23:57:48Z | 0 | `go test -count=1 ./...` |
| adr | 23:57:48 to 23:57:49Z | 0 | `sh docs/adr/adr-lint.sh` |
| prd | 23:57:49 to 23:57:49Z | 0 | `sh docs/prd/prd-lint.sh` |
| link | 23:57:49 to 23:58:02Z | 0 | `sh docs/links/link-lint.sh` |
| disc | 23:58:02 to 23:58:04Z | 0 | `sh docs/tests/run-discipline-tests.sh` |
| nested | 23:58:04 to 23:58:04Z | 0 | `sh docs/tests/nested-checkout-check.sh` |
| setup | 23:58:04 to 23:58:25Z | 0 | `sh docs/setup/setup-check.sh` |
| diffcheck | 23:58:25 to 23:58:25Z | 2 | `git diff --check origin/main HEAD` |
| taskstate | 23:58:25 to 23:58:33Z | 0 | `sh docs/tasks/task-state.sh` |
| e2e | 00:01:38 to 00:02:43Z | 0 | `go test -count=1 -tags=e2e ./...` |

- `diffcheck` ran first with exit 2 (`runs/T-4naq/mutants.md:66: new blank line at EOF`); the file was fixed and committed, and `git diff --check origin/main HEAD` then exited 0 (re-run in the freeze check below).
- `taskstate` printed `T-4naq #428 stalled` only because the branch has had no push yet; the command lists every task and exits 0.

### The full integration level

`go test -count=1 -tags=integration -timeout 30m ./...`, in the background, output only to `/home/layup/.cache/layup-session/T-4naq-full-integration.txt`: 2026-10-20T23:57:42Z to 2026-10-21T00:01:23Z, exit 0, 21 packages `ok`. An earlier full run on the first merge head `6c63a4f` (23:40:24Z to 23:44:11Z) also exited 0.

## The freeze check

On the clean committed head `33b9a14`, which has merged `origin/main` at `f9d46b0`: `sh docs/tests/freeze-check.sh --task T-4naq --base origin/main --budget 2600/30 --evidence runs/T-4naq/test-runs.md --evidence runs/T-4naq/mutants.md --evidence runs/T-4naq/conditions.md`, output kept outside the worktree in `/home/layup/.cache/layup-session/T-4naq-freeze.txt`, exit 0:

```text
ok K1: white space
ok K2: registration of new tests
ok K3: test selectors
ok K4: restore evidence of mutants
ok K5: projection: figure 1810/25, maximum 2600/30, tolerance 500/4, limit 3100/34, margin 1290/9 (the figure at landing is still measured at landing)
ok K6: conditions table
```

The figure is measured at landing; the figure above is the projection of K5 at the head named, with the evidence records of `runs/T-4naq/` out of the count (ADR-0028) and before the close-out.

## Round 1 of #428: the red of F1, F2 and F3, before the fixes

Round 1 (the comment on #428, reviewer GPT-6.1 Sol, head `977d2fb`) gave `material` with F1 to F5. The stubs of the fixes are the commit `d80a9bc` (the type `spec.PendingReq`, the field `FullInput.Pending`, `pendingOf` that gives nothing, the version in the signature of `bindReopenings` that nothing reads). The tests below were written and run against it, in the working tree, **before** any line of the fixes; the hook refuses a red commit, so they are committed with the fixes. Each command alone, `TMPDIR=/home/layup/.cache/layup-session/tmp-green` set, with its direct exit status.

### f1spec (F1)

- Command: `go test -count=1 ./internal/spec/ -run '^TestAReopenedFactNeedsItsOwnPendingLineEvenWhenARequirementCoversIt$'`
- Start to end: 00:17:00 to 00:17:00Z (`date -u`), exit 1.

First failing assertions:

```text
--- FAIL: TestAReopenedFactNeedsItsOwnPendingLineEvenWhenARequirementCoversIt (0.00s)
    full_test.go:799: the rule needs: ("pass", ""), want fail with "NFR-002" in the detail
    full_test.go:801: the rule needs: ("pass", ""), want fail with "NFR-002" in the detail
    full_test.go:802: the rule needs: ("pass", ""), want fail with "NFR-002" in the detail
    full_test.go:803: the rule needs: ("pass", ""), want fail with "NFR-002" in the detail
    full_test.go:807: the rule measures: ("pass", ""), want fail with "3" in the detail
```

### f2 (F2 (a panic, not an error))

- Command: `go test -count=1 ./internal/run/ -run '^TestAReopeningsTableRefusesEachBrokenRule$'`
- Start to end: 00:17:00 to 00:17:02Z (`date -u`), exit 1.

First failing assertions:

```text
--- FAIL: TestAReopeningsTableRefusesEachBrokenRule (0.00s)
panic: runtime error: invalid memory address or nil pointer dereference [recovered, repanicked]
```

### f3 (F3)

- Command: `go test -count=1 ./internal/run/ -run '^TestTheRowsOfTheRecordAreBoundToTheSelection$'`
- Start to end: 00:17:02 to 00:17:02Z (`date -u`), exit 1.

First failing assertions:

```text
--- FAIL: TestTheRowsOfTheRecordAreBoundToTheSelection (0.00s)
    reopen_test.go:136: another version: <nil>; want an error holding "version"
```

### f1run (F1 (the identity that specInput carries))

- Command: `go test -count=1 ./internal/run/ -run '^(TestTheReopeningSetOfTheFullInputIsReadFromTheRecords|TestThePendingIdentityIsTheStoredRowOfTheTask)$'`
- Start to end: 00:17:02 to 00:17:02Z (`date -u`), exit 1.

First failing assertions:

```text
--- FAIL: TestTheReopeningSetOfTheFullInputIsReadFromTheRecords (0.00s)
    reopen_test.go:441: its own task, not merged: Reopening = [], Pending = [], <nil>; want [] and [{1 NFR-002}]
--- FAIL: TestThePendingIdentityIsTheStoredRowOfTheTask (0.00s)
    reopen_test.go:473: pendingOf(T-aaaa) = []; want [{3 REQ-006}]
    reopen_test.go:473: pendingOf(T-bbbb) = []; want [{12 }]
```

### integ (F1 and F3 on the production path)

- Command: `go test -count=1 -tags=integration -timeout 30m ./internal/run/ -run '^(TestAReopenedConstraintThatARequirementCoversNeedsItsOwnPendingLine|TestARegisteredTaskThatDisagreesWithItsReservationIsRefused|TestTheReopenedTasksPassTheFullCheckOneAfterTheOther)$'`
- Start to end: 00:17:02 to 00:17:10Z (`date -u`), exit 1.

First failing assertions:

```text
--- FAIL: TestAReopenedConstraintThatARequirementCoversNeedsItsOwnPendingLine (5.15s)
    reopen_integration_test.go:361: no line, though REQ-001 covers the fact: needs = "pass "; want "fail the reopened constraint fact 1 has no pending requirement NFR-002 under waits for a bet"
    reopen_integration_test.go:361: the token and the ID of the old requirement: needs = "pass "; want "fail the reopened constraint fact 1 has no pending requirement NFR-002 under waits for a bet"
    reopen_integration_test.go:361: its ID and the token of another fact: needs = "pass "; want "fail the reopened constraint fact 1 has no pending requirement NFR-002 under waits for a bet"
--- FAIL: TestARegisteredTaskThatDisagreesWithItsReservationIsRefused (5.46s)
    reopen_integration_test.go:482: openReopenings with the version 2 = [], the task T-rsv1 of spec/reopenings.tsv is registered as specification for [REQ-009] after [T-erd2]; the reservation is a specification task for [NFR-002] after [T-e
```

F2 is red as a panic of the reader (`nil pointer dereference` at the expression of the requirement), which ends that run of the test; the case is `no class`. The integration test of F3 is red because the reservation of the version 2 is not refused for its version: the error it gives is the disagreement of the registered `T-rsv1`, the phase before it, which comes after the binding in the correct code.

## Round 1 of #428: green and the gates

The fixes are in `4aa800c`; the tests ran on the head `4aa800c` plus the record commit that follows it (records only). Each command alone, with its direct exit status.

- The tests of the fixes, green (unit level of the two packages, then the integration selection of the 14 integration tests of this task): 00:17:45Z to 00:17:46Z exit 0, and 00:18:18Z to 00:18:41Z exit 0, on the working tree that `4aa800c` committed (one expectation changed with the fix: the detail of the first failing head of `TestTheReopenedTasksPassTheFullCheckOneAfterTheOther` is now `the reopened need fact 2 has no pending requirement REQ-003 under waits for a bet`, because the stronger check of F1 comes first).

### The gates

| Gate | Start to end | Exit | Command |
| ---- | ------------ | ---- | ------- |
| build | 00:31:56 to 00:31:57Z | 0 | `go build ./...` |
| vet | 00:31:57 to 00:31:57Z | 0 | `go vet ./...` |
| vetint | 00:31:57 to 00:31:57Z | 0 | `go vet -tags=integration ./...` |
| vete2e | 00:31:57 to 00:31:58Z | 0 | `go vet -tags=e2e ./...` |
| gofmt | 00:31:58 to 00:31:58Z | 0 | `sh -c out=$(gofmt -l .) && test -z "$out"` |
| unit | 00:31:58 to 00:31:59Z | 0 | `go test -count=1 ./...` |
| adr | 00:31:59 to 00:32:00Z | 0 | `sh docs/adr/adr-lint.sh` |
| prd | 00:32:00 to 00:32:00Z | 0 | `sh docs/prd/prd-lint.sh` |
| link | 00:32:00 to 00:32:13Z | 0 | `sh docs/links/link-lint.sh` |
| disc | 00:32:13 to 00:32:15Z | 0 | `sh docs/tests/run-discipline-tests.sh` |
| nested | 00:32:15 to 00:32:15Z | 0 | `sh docs/tests/nested-checkout-check.sh` |
| setup | 00:32:15 to 00:32:36Z | 0 | `sh docs/setup/setup-check.sh` |
| diffcheck | 00:32:36 to 00:32:36Z | 0 | `git diff --check origin/main HEAD` |
| taskstate | 00:32:36 to 00:32:45Z | 0 | `sh docs/tasks/task-state.sh` |
| e2e | 00:36:03 to 00:37:10Z | 0 | `go test -count=1 -tags=e2e ./...` |

### The full integration level

`TMPDIR=/home/layup/.cache/layup-session/tmp-green go test -count=1 -tags=integration -timeout 30m ./...`, in the background, output only to `/home/layup/.cache/layup-session/T-4naq-full-integration.txt`: 2026-10-21T00:31:50Z to 2026-10-21T00:35:36Z, exit 0, 21 packages `ok`, on the code head `4aa800c` with `origin/main` at `f9d46b0` merged.

### The freeze check after the fixes of round 1

On the clean committed head, which has merged `origin/main` at `f9d46b0` (still its head at 00:40Z): `sh docs/tests/freeze-check.sh --task T-4naq --base origin/main --budget 2600/30 --evidence runs/T-4naq/test-runs.md --evidence runs/T-4naq/mutants.md --evidence runs/T-4naq/conditions.md`, output kept outside the worktree in `/home/layup/.cache/layup-session/T-4naq-freeze.txt`, exit 0 (the first run of it failed K1 on a blank line at the end of this record, which was fixed and committed):

```text
ok K1: white space
ok K2: registration of new tests
ok K3: test selectors
ok K4: restore evidence of mutants
ok K5: projection: figure 1963/25, maximum 2600/30, tolerance 500/4, limit 3100/34, margin 1137/9 (the figure at landing is still measured at landing)
ok K6: conditions table
```

## Round 2 of #428: the red of F6, F7 and F8, before the fixes

Round 2 (the comment 6305218412 on #428, GPT-6.1 Sol, head `9101153`) gave `not mergeable` with F6 to F8; the Operator raised the cap to 2 and accepted the R8 exception of F5 (comment 6307600855), so the late labels of F5 above stay as they are. The stub of round 2 is `ca5fc37` (`pendingOf` takes the merged tasks and reads none). The tests were run in the working tree against it **before** any line of the fixes, each command alone with its direct exit status, `TMPDIR=/home/layup/.cache/layup-session/tmp-green` set.

### f6unit (F6, the identity that specInput carries and pendingOf gives)

- Command: `go test -count=1 ./internal/run/ -run '^(TestThePendingIdentityIsTheStoredRowOfTheTask|TestTheReopeningSetOfTheFullInputIsReadFromTheRecords)$'`
- Start to end: 03:54:02 to 03:54:04Z (`date -u`), exit 1.

First failing assertions:

```text
--- FAIL: TestTheReopeningSetOfTheFullInputIsReadFromTheRecords (0.00s)
    reopen_test.go:444: another task, the task of the fact merged: Reopening = [], Pending = [], <nil>; want [] and [{1 NFR-002}]
--- FAIL: TestThePendingIdentityIsTheStoredRowOfTheTask (0.00s)
    reopen_test.go:489: pendingOf(T-cccc, merged [T-aaaa]) = []; want [{3 REQ-006}]
    reopen_test.go:489: pendingOf(T-bbbb, merged [T-aaaa]) = [{12 }]; want [{3 REQ-006} {12 }]
    reopen_test.go:489: pendingOf(T-aaaa, merged [T-aaaa T-bbbb]) = [{3 REQ-006}]; want [{3 REQ-006} {12 }]
```

### f6int (F6 on the production path)

- Command: `go test -count=1 -tags=integration -timeout 30m ./internal/run/ -run '^TestAReopenedConstraintThatARequirementCoversNeedsItsOwnPendingLine$'`
- Start to end: 03:54:04 to 03:54:12Z (`date -u`), exit 1.

First failing assertions:

```text
--- FAIL: TestAReopenedConstraintThatARequirementCoversNeedsItsOwnPendingLine (5.63s)
    reopen_integration_test.go:383: another task's head that removes only the merged line: needs = "pass "; want "fail the reopened constraint fact 1 has no pending requirement NFR-002 under waits for a bet"
    reopen_integration_test.go:383: another task's head that keeps the old requirement's citation only: needs = "pass "; want "fail the reopened constraint fact 1 has no pending requirement NFR-002 under waits for a bet"
```

### f7 (F7)

- Command: `go test -count=1 -tags=integration -timeout 30m ./internal/run/ -run '^TestAnEmptySelectionStillChecksAStoredReservation$'`
- Start to end: 03:54:12 to 03:54:17Z (`date -u`), exit 1.

First failing assertions:

```text
--- FAIL: TestAnEmptySelectionStillChecksAStoredReservation (5.05s)
    reopen_integration_test.go:405: an empty selection with a stored reservation = [], <nil>; 0 issues; want the refusal that names spec/reopenings.tsv and no issue
```

### f8 (F8, green on the code: no code change is needed)

- Command: `go test -count=1 -tags=integration -timeout 30m ./internal/run/ -run '^(TestTheReopenedSpecificationPromptIsItsGoldenFile|TestTheReopenedConstraintPromptCarriesItsStoredNFR|TestAHarnessThatFollowsTheReopenedInstructionGivesAValidHandoffAndAGreenSpecStatus)$'`
- Start to end: 03:54:17 to 03:54:26Z (`date -u`), exit 0.

First failing assertions:

```text
ok  	github.com/pharzam/layup/internal/run	8.468s
```

F8 asks for tests, not for code: the constraint prompt, its golden file and the output scenario for a constraint pass on the code, so their red is the mutant of the reviewer (a constraint prompt that substitutes `REQ-999` for the allocated NFR), run in `mutants.md` (M61 and M62). The output scenario is now a loop over a constraint (the fact 1, NFR-002, which REQ-001 already covers) and a need.

## Round 2 of #428: green and the gates

The fixes are in `3806384`; the tests ran on that head (the later commit changes records only). Each command alone, with its direct exit status.

- The tests of the fixes, green: unit level of the two packages 03:54:49Z to 03:54:49Z exit 0, and the integration selection (the 15 integration tests of this task) 03:55:20Z to 03:55:46Z exit 0. Two expectations changed with the fix of F6: the first failing head of `TestTheReopenedTasksPassTheFullCheckOneAfterTheOther` after the merge of the first task is now `the reopened need fact 2 has no pending requirement REQ-003 under waits for a bet` (the merged identity is read first), and the output scenario of the harness test now runs a constraint and a need.

### The gates

| Gate | Start to end | Exit | Command |
| ---- | ------------ | ---- | ------- |
| build | 04:10:02 to 04:10:02Z | 0 | `go build ./...` |
| vet | 04:10:02 to 04:10:02Z | 0 | `go vet ./...` |
| vetint | 04:10:02 to 04:10:03Z | 0 | `go vet -tags=integration ./...` |
| vete2e | 04:10:03 to 04:10:03Z | 0 | `go vet -tags=e2e ./...` |
| gofmt | 04:10:03 to 04:10:03Z | 0 | `sh -c out=$(gofmt -l .) && test -z "$out"` |
| unit | 04:10:03 to 04:10:04Z | 0 | `go test -count=1 ./...` |
| adr | 04:10:04 to 04:10:06Z | 0 | `sh docs/adr/adr-lint.sh` |
| prd | 04:10:06 to 04:10:06Z | 0 | `sh docs/prd/prd-lint.sh` |
| link | 04:10:06 to 04:10:18Z | 0 | `sh docs/links/link-lint.sh` |
| disc | 04:10:18 to 04:10:21Z | 0 | `sh docs/tests/run-discipline-tests.sh` |
| nested | 04:10:21 to 04:10:21Z | 0 | `sh docs/tests/nested-checkout-check.sh` |
| setup | 04:10:21 to 04:10:42Z | 0 | `sh docs/setup/setup-check.sh` |
| diffcheck | 04:10:42 to 04:10:42Z | 0 | `git diff --check origin/main HEAD` |
| taskstate | 04:10:42 to 04:10:50Z | 0 | `sh docs/tasks/task-state.sh` |
| e2e | 04:13:59 to 04:14:56Z | 0 | `go test -count=1 -tags=e2e ./...` |

### The full integration level

`TMPDIR=/home/layup/.cache/layup-session/tmp-green go test -count=1 -tags=integration -timeout 30m ./...`, in the background, output only to `/home/layup/.cache/layup-session/T-4naq-full-integration.txt`: 2026-10-21T04:09:58Z to 2026-10-21T04:13:47Z, exit 0, 21 packages `ok`, on the code head `3806384` with `origin/main` at `f9d46b0` merged (`go.mod` names go 1.26.9 and the toolchain is go1.26.9).
