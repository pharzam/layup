# T-9nbf: the test runs

Task `T-9nbf` (#367), row 66k of M2e: the step loop and the restart. Times are `date -u` readings taken after each step; an exit status is the direct one of the command, never of a pipe. Go is 1.26.10 (`GOROOT` and `PATH` of mise, in every command).

## S0. The reconciliation

Main at the start: `be6a0c7` (the merge of #458, K7). The branch `T-9nbf` stood at the claim `68a8a74` with the one registration commit `006738d`, which is not an ancestor of main; `git merge origin/main` conflicted in `docs/plan/README.md` and `docs/spec/loop.md` only, and was resolved to main's rows of the children plus the narrowed Task cell of 66k (with the After cell `66c, 66f, 66g, 66h, 66j, 66l, 66m, 66n, 66o, 66q, 66r, 66s, 66u, 66v, 66w, 66x, 66y, 66z, 66aa, 66ab`, every label of which now has a row) and main's side of `loop.md` (the "Not in M2e" sentence on the execution of a frozen-test amendment was already there from the merge of `006738d`'s unconflicted lines). Commit `4e41374`.

| Cited by the plan (6304795183) | At `be6a0c7` | Divergence |
| --- | --- | --- |
| `RepairTask` / `RepairPushes` / `RepostComments` `repair.go:30,51,144`; called first by `AdvanceShapeTask` `shapeloop.go:484` | same lines (`AdvanceShapeTask` `:482`) | none |
| `ClosePullRequest`, `CloseIssue` (66w) | `forge.Forge`, `internal/forge/forge.go` | none; not rebuilt |
| `BringUpToDate` `baseup.go:65`, `UpToDate` `:43`, redo `redo.go:173`; a redo that would drop the rendered commit is a Stop `baseup.go:150-156` | `:65`, `:43`; the Stop at `:156`; a merged pull request is a Stop at `:128` | none (the line of the Stop moved by 4) |
| `FinishTask` `finish.go:45`; the terminal read, repair, verified head, carried head, `UpToDate`, close-out, merge | `finish.go:45`; the cut `:69-79` | the plan's `finish.go:70-79` is `:69-79`; the demand of the rendered commit is at `:83` (plan `:87`), its error at `:141` (the review says `:145`) |
| the only production `Position` -> `TierOf` consumer `finish.go:252,391` | `finish.go:248,252` (`startCloseOut`) | none |
| `openReopenings` and `reopenedTasks` of `origin/T-4naq` (`reopen.go:398,125`) | landed: `reopen.go:401`, `:125`; `reopeningOf` `:134`; the prompt reads the record at `promptfile.go:577` | the call is `openReopenings(ctx, spec)` of `*state`; it gives `[]Task` and is idempotent |
| `loopPosition` stops on a bound of a session that is not the developer's `shapeloop.go:282-286` | `:282` | none (F1) |
| `verifiedStep` `betbrief.go:420-429` called at `bet.go:178` | `betbrief.go:420`, called at `bet.go:172` | the call is at `:172` (as the review says) |
| `archStart` `shape.go:23-34`; `tasksStage` `shape.go:72-80` | `shape.go:22`, `:69` | lines moved only |
| `parametersBoundary` `parameters.go:404-412`, error for an absent table `:421-423` | `:412`, `:422` | none (F3) |
| `loopPermissionsStep` `register.go:254` | `:254` | none |
| `toVerifier` passes no `Freeze` (`shapeloop.go:695`; `frozen.go:203`); `Position` is uphill `tier.go:121` | `shapeloop.go:675` (`toVerifier`); `tier.go:106` (`Position`), the frozen-list decision `:132`, the uphill return `:141` | none: the review's citations (`tier.go:132`, `:141`) are the lines; the plan's `:121` was stale (F5) |
| `optionsKept` `instructions.go:178-201`; `Result.Options` must equal the whole table | `instructions.go:170` | the reader is unchanged and shared: `bindingOptions` gives it every option of the table and the bound one as the only survivor |
| close-out role `software-developer` `internal/session/context.go:61` | `context.go:61` | none |
| `waitingSpecAccepted` asked by `ship` `shapeloop.go:743` | `shapeloop.go:744` | none |

No step of the plan is impossible, and the plan, its review and the decisions settle every decision this task makes. Three small divergences, recorded as the plan review's own rules direct:

1. `loopStep` and the scheduler are written over a `loopEnv` of injected reads and acts (`runLoop`), so that U1 to U8 test the pure scheduler (the plan's "pure scheduler with injected reads and acts"); `actOn` and `dispatch` hold the act; `wrapBoundary` and `boundaryRule` hold the permission and the boundary (U10 is a unit test of them, not of a state with a store).
2. A test seam: `state.sessionsHook` (a function of the Sessions that the step makes) lets a test inject a fault into the Sessions of `loop` (a push that went down, a refused commit); it is nil in production. `afterScaffold()` is the tail of `Restart` after the Scaffold, so a test composes the three steps without the whole restart: a `Restart` of the world of these tests cannot run the steps of Intake (the world has no `meaning.tsv` session).
3. Condition 6: **the original conditions 1 and 5 of 6252684822 are reconciled here, not discharged by the plan.** Condition 1 (the cuts, the barrier, the After cell) is met by the registration `006738d`, its merge `4e41374` and the After cell above, which names every child. Condition 5 (K3 and K7 contracts) is met by the landed `BringUpToDate` (`baseup.go:65`) and `openReopenings` (`reopen.go:401`), used as they are. The citation `F-0003#45` is kept in the Task cell of 66k (`docs/plan/README.md`), and the "D13" of the plan is D11 (unsupported kinds) in `loop.md`, The step loop, item 11.

## S2 and the stubs

`ba18e81` (`verifiedView` extracted with its one caller `FinishTask`, `optionResults` and `defaultConstraints` extracted from `optionsStage`, the stubs of `loop.go`), `4515549` (`afterScaffold`, the same two steps), `e8aaf1d` (the stubs of `dispatch`, `boundaryRule` and `wrapBoundary`). Each was built, vetted (with and without the tag `integration`) and unit-tested green. Wiring `verifiedView` into `verifiedStep` is **not** in them: it changes that caller, so it comes after its red test (condition 5).

## S3. The red run (against the stubs, before the code)

Every command was run alone, with its own status, output to a file and not through a pipe. The tests are in the working tree and uncommitted: the hook refuses a red commit (`guardrails.md` section 2), so they are committed with the code. The first failing assertions are quoted from the outputs.

### unit

- Command: `TMPDIR=/home/layup/.cache/layup-session/tmp-green go test -count=1 ./internal/run/ -run '^(TestTheMergeOrderIsSpecificationReopenedShapeArchitectureThenBetArchitecture|TestAKindThatM2eDoesNotRunFailsClosed|TestATaskWithAnUnmergedPredecessorIsNotStarted|TestOnlyTheFirstUnmergedTaskIsActedOn|TestTheRegisterIsReadAgainAfterAReplacementOrAReopening|TestDoneBeforeOnlyWhenThisInvocationActedOnNoTask|TestAWaitIsTheLastRowAndAStopIsAFail|TestAnActThatMakesNoProgressFails|TestTheRestartNamesEndWithLoopOnlyAfterTheScaffold|TestBetFollowsShapeOnlyOnceTheScaffoldIsDone|TestTheBoundaryIsSkippedOnlyForAnAbsentTable|TestTheActDispatchesByTheRecords|TestAReopenedSpecificationIsNotTheOriginal)$'`
- Start to end: 08:09:28 to 08:09:29Z, exit 1; 13 tests fail, all 13.

First failing assertion of each:

- `TestBetFollowsShapeOnlyOnceTheScaffoldIsDone`: `bet_test.go:244: the last rows = [scaffold-push shape bet]; want shape, bet, loop`
- `TestTheMergeOrderIsSpecificationReopenedShapeArchitectureThenBetArchitecture`: `loop_test.go:39: the order = [], <nil>; want [S R1 R2 A B1 B2] (the order of the rows of the reopenings, then Shape's architecture, then the bet's in register order)`
- `TestAKindThatM2eDoesNotRunFailsClosed`: `loop_test.go:50: kind batch: <nil>; want the refusal, naming row 66t for rendered only`
- `TestATaskWithAnUnmergedPredecessorIsNotStarted`: `loop_test.go:59: pick = {ID: Issue:0 Kind: Requirements:[] Size: Predecessors:[] PR:0 Opened:0001-01-01 00:00:00 +0000 UTC}, false, <nil>; want the failure that names the predecessor`
- `TestOnlyTheFirstUnmergedTaskIsActedOn`: `loop_test.go:117: acts [] and the row {Name: Result: Detail:}; want [S S A], never A before S has merged, and the merged tasks named`
- `TestTheRegisterIsReadAgainAfterAReplacementOrAReopening`: `loop_test.go:142: acts [] after 0 openings, row {Name: Result: Detail:}; want [S R Aold Anew]: the reopened task is picked after the merge of S, the replaced ID never again, the new one`
- `TestDoneBeforeOnlyWhenThisInvocationActedOnNoTask`: `loop_test.go:157: every task merged: {Name: Result: Detail:}; want done before`
- `TestAWaitIsTheLastRowAndAStopIsAFail`: `loop_test.go:185: {Task: Verified:false Merged:false Head: Attempt:0 Wait:no word that says it Stop: Replaced:}, <nil>: row {Name: Result: Detail:} after 0 acts; want {Name:loop Result:done Detail:no word that says it} after 2`
- `TestAnActThatMakesNoProgressFails`: `loop_test.go:195: row {Name: Result: Detail:} after 0 acts; want the failure at the second act, and exactly 2 acts`
- `TestTheRestartNamesEndWithLoopOnlyAfterTheScaffold`: `loop_test.go:212: the last rows after the Scaffold = [scaffold-push shape bet]; want shape, bet, loop`
- `TestTheBoundaryIsSkippedOnlyForAnAbsentTable`: `loop_test.go:230: parameters "", tasks "": skip false, err <nil>; want skip true, failure false`
- `TestTheActDispatchesByTheRecords`: `loop_test.go:281: merged: dispatch = 0; want 1`
- `TestAReopenedSpecificationIsNotTheOriginal`: `loop_test.go:301: originalSpec: S true, R true, A false; want true, false, false`

The cases of `TestTheActDispatchesByTheRecords` that pass at the stub (`no attempt`, `an attempt that is not verified`, `a handback`: the driver) are the stub's own answer; they are characterization with their red composition control in the cases that do not pass (merged, closed, terminal, verified, close-out bound).

### integration

- Command: `TMPDIR=/home/layup/.cache/layup-session/tmp-green go test -count=1 -tags=integration ./internal/run/ -run '^(TestAfterTheBetLoopMergesTheShapeTasksAndTheBetsTaskInOrder|TestARestartAfterTheFirstCloseOutReachesLoop|TestARestartAfterTheArchitectureTaskIsReplacedReachesLoop|TestALostCloseOutPushAfterTheBetIsRefusedRetriedAndMerged|TestASessionCommentOfATaskAfterTheBetIsOnItsIssue|TestALostMergedEventIsRecordedByLoopWithNoSecondMerge|TestARestartAfterEachWriteOfTheLoopAdoptsIt|TestASetIsAppliedAtEachBoundaryOfTheTail|TestAnAbsentParametersTableIsSkippedOnlyBeforeTheFirstTask|TestAMissingLoopPermissionFailsTheRowOfEachStepBeforeAnyWrite|TestTheTaskLoopWithTheStandInForgeThroughLoop|TestACloseOutThatTouchesTheReadmeIsRefusedThroughLoop|TestAnOverlappingMovedBaseIsRedoneThroughLoop|TestAReopenedTaskIsMergedByLoopWithNoRenderedCommit|TestAnAttemptAfterTheBetWithNoRenderedCommitIsTheStopOfTheOperator|TestTheCloseOutOfALoopTaskIsRoutedFromItsPositionNotItsPlan|TestTheTextBlocksOfRunMd)$'`
- Start to end: 08:08:27 to 08:08:40Z, exit 1; 17 tests fail.

First failing assertion of each:

- `TestTheTextBlocksOfRunMd`: `run_integration_test.go:676: schema run-steps: column 1 (step): the type: the block has "enum(forge|plan|baseline|root-push|read-back|records|issues|watch|clone|version|lease|probe|phase|rule-gaps|meaning|numbering|spans|draft|completeness|markers|batch|answers|answers-review|follow-up|follow-up-ans`
- `TestAnAbsentParametersTableIsSkippedOnlyBeforeTheFirstTask`: `loop_integration_test.go:389: the boundary before the first task = parameters.tsv is not in the records: the first task of the loop writes it (records moved: false); want a skip`
- `TestALostCloseOutPushAfterTheBetIsRefusedRetriedAndMerged`: `loop_integration_test.go:219: the first loop = {Name: Result: Detail:}; want the fail of the push that went down`
- `TestACloseOutThatTouchesTheReadmeIsRefusedThroughLoop`: `loop_integration_test.go:467: loop = {Name: Result: Detail:}; want a fail that names the refused session (scope) and no merge`
- `TestASessionCommentOfATaskAfterTheBetIsOnItsIssue`: `loop_integration_test.go:237: the act on T-zn3t = {Task: Verified:false Merged:false Head: Attempt:0 Wait: Stop: Replaced:}, <nil>; want its driver to a verified head`
- `TestAnOverlappingMovedBaseIsRedoneThroughLoop`: `loop_integration_test.go:484: loop = {Name: Result: Detail:}; want the specification, the task that redoes T-6p1x, and T-ozz3 merged`
- `TestTheTaskLoopWithTheStandInForgeThroughLoop`: `loop_integration_test.go:445: loop = {Name: Result: Detail:}; want the tasks merged`
- `TestTheCloseOutOfALoopTaskIsRoutedFromItsPositionNotItsPlan`: `loop_integration_test.go:572: the act on T-1pi7 = {Task: Verified:false Merged:false Head: Attempt:0 Wait: Stop: Replaced:}, <nil>; want it merged`
- `TestALostMergedEventIsRecordedByLoopWithNoSecondMerge`: `loop_integration_test.go:252: the first loop = {Name: Result: Detail:} with the events []; want a fail after the merge and before its event`
- `TestARestartAfterEachWriteOfTheLoopAdoptsIt`: `loop_integration_test.go:282: stop 1: the loop = {Name: Result: Detail:}; want the fail of its fault`
- `TestASetIsAppliedAtEachBoundaryOfTheTail`: `loop_integration_test.go:364: the rows = [{Name:shape Result:done Detail:done before} {Name:bet Result:done Detail:done before}]; want the three steps done`
- `TestAfterTheBetLoopMergesTheShapeTasksAndTheBetsTaskInOrder`: `loop_integration_test.go:142: the rows = [{Name:shape Result:done Detail:done before} {Name:bet Result:done Detail:done before}]; want [{Name:shape Result:done Detail:done before} {Name:bet Result:done Detail:done before} {Name:loop Result:done Detail:the loop merged T-waus, T-44z0, T-aemy}]`
- `TestAMissingLoopPermissionFailsTheRowOfEachStepBeforeAnyWrite`: `loop_integration_test.go:429: shape = {Name:shape Result:done Detail:done before}; want its row, fail, forge: the installation lacks the permissions of the task loop: statuses: write`
- `TestAnAttemptAfterTheBetWithNoRenderedCommitIsTheStopOfTheOperator`: `loop_integration_test.go:557: loop 1 after the loss = {Name: Result: Detail:} with 0 sessions more; want {Name:loop Result:fail Detail:waiting for the Operator: T-qeb9: an attempt after the bet has no rendered commit} and no session of a driver`
- `TestARestartAfterTheFirstCloseOutReachesLoop`: `loop_integration_test.go:176: the rows = [{Name:shape Result:fail Detail:T-9cuu: the session S-edb1abe6 has a bound and is not of the developer role}]; want shape and bet done before, then the wait of loop`
- `TestAReopenedTaskIsMergedByLoopWithNoRenderedCommit`: `loop_integration_test.go:532: the act on the reopened task = {Task: Verified:false Merged:false Head: Attempt:0 Wait: Stop: Replaced:}, <nil>; want it merged with no rendered commit asked ([attempt session result handoff session result handoff session result push bound handoff session result handoff`
- `TestARestartAfterTheArchitectureTaskIsReplacedReachesLoop`: `loop_integration_test.go:207: shape after the redo = {Name:shape Result:fail Detail:T-mxoy: its base 46ab02b40c89badee255033a2ee53e2d57ad0499 is not the verified head 8da4c9a16d6029d67db471161758df3d85a5ae03 of T-tdqf}; want done before`

Two further tests (condition 2: a closed unmerged pull request with an unresolved push and a missing comment, and the opening of the reopenings of a merge recorded before the invocation):

- Command: `TMPDIR=/home/layup/.cache/layup-session/tmp-green go test -count=1 -tags=integration ./internal/run/ -run '^(TestAClosedUnmergedPullRequestIsTheWaitOfLoopWithNoWrite|TestLoopOpensTheReopeningsOfAMergeThatWasRecordedBefore)$'`
- Start to end: 08:09:00 to 08:09:08Z, exit 1.

- `TestLoopOpensTheReopeningsOfAMergeThatWasRecordedBefore`: `loop_integration_test.go:612: 0 tasks and 0 reservations opened by loop; want the 4 that the bet re-classed`
- `TestAClosedUnmergedPullRequestIsTheWaitOfLoopWithNoWrite`: `loop_integration_test.go:594: loop 1 = {Name: Result: Detail:}; want {Name:loop Result:done Detail:waiting for the Operator: T-rl0z: the pull request #1 of task/T-rl0z/1 is closed and not merged}`

An earlier run of the same set (08:04:16 to 08:04:27Z, exit 1) showed two defects of the tests themselves, fixed before this record: `TestARestartAfterTheArchitectureTaskIsReplacedReachesLoop` and `TestAnOverlappingMovedBaseIsRedoneThroughLoop` moved the default branch before the specification had merged, so the base of the architecture task of Shape (the specification's verified head) was not yet an ancestor of the default branch (`the base ... is not an ancestor of the default branch`); the move now follows the first merge (`afterMerge`), and the first test merges the specification outside, as a person does, so that no close-out masks F1 in it. The new task of that test also needs the copy of the answer (D10): it is given as `loop` gives it.

### F1 to F5

- **F1 reproduced.** `TestARestartAfterTheFirstCloseOutReachesLoop`: `shape` fails with `... the session S-... has a bound and is not of the developer role`, and `verifiedStep` of the task fails with `Shape has no verified head for it` (the output above).
- **F2 reproduced, isolated.** `TestARestartAfterTheArchitectureTaskIsReplacedReachesLoop`, the specification merged outside so that F1 does not mask it: `shape` fails with `T-...: its base ... is not the verified head ... of T-...`.
- **F3 reproduced.** `TestAnAbsentParametersTableIsSkippedOnlyBeforeTheFirstTask`: the boundary before the first task is `parameters.tsv is not in the records: the first task of the loop writes it`.
- **F4 reproduced (the premise).** `TestAnAttemptAfterTheBetWithNoRenderedCommitIsTheStopOfTheOperator` makes the lost close-out push of the specification, the repair (the handback of attempt 2) and then asserts that `FinishTask` of that attempt is an error while the step bet is done; it passes at the stub, and the test fails at its next assertion, the stop that `loop` is to give.
- **F5** is a reachability limit (no Shape task freezes a list), not a defect to observe: `TestTheCloseOutOfALoopTaskIsRoutedFromItsPositionNotItsPlan` (I12) is the consumer test, as the Operator decided (6307601453).

### Characterization and their controls

- `TestALostCloseOutPushAfterTheBetIsRefusedRetriedAndMerged` (I3) and `TestASessionCommentOfATaskAfterTheBetIsOnItsIssue` (I4) test what 66v and the drivers already do; each has its red composition control, the call of `loop` / `actOn` that does not exist at the stub.
- The premise lines of I13 (the repair and the handback) and the setup of I3 are 66v's and 66y's behaviour, green at the stub.

## S4 and S5. Green

The code and the tests were committed together (`3ae88e2`), as the hook runs `go test`. Against the stubs 13 of the 13 unit tests and 17 of the 17 integration tests that were red are green; the two tests of condition 2, the carry of the rendered commit, the latest-attempt check and the case before the bet were added in the same pass, each red at the stub by the same rule (they call a function or a step that the stub does not give).

Findings of the first full run of the integration level on the new code (08:14:51 to 08:19:15Z, exit 1, before the fixes below), each an existing test that the specification changed:

- `TestTheStepForgeCheckChecksTheForgeAndCommitsItsTableAtEachRestart`: the count of rows is 25 once the Scaffold is done (the row `loop`); the test said 24 (a stale comment and a number, `forgecheck_integration_test.go`).
- `TestALaterBetFailsClosed`: its `bets.tsv` held a row with the verdict `bet` for another milestone, which D1 now hands to `loop`; the later bet it tests is told by a first bet that was not made, so its row is a `no-bet` of the other milestone.
- (the e2e level, below) the fake installation of `cmd/layup/run_e2e_test.go` granted only the three permissions of M2a; `shape` now checks the permissions of the loop before it does anything (condition 4), so the fake grants them.

## S6. The mutants

`runs/T-9nbf/mutants.md`: 28 mutants (M1 to M18 of the plan and G1 to G10, one for each new guard) on a copy of the tree, each built and vetted, 26 killed, 2 not killed (G1 equivalent, G5 not reached by a fixture), every restore matched its hash and every control was green.

## S8. The checks (head `3ae88e2` plus the change of the e2e fake; each command alone, output to a file outside the worktree)

| Command | Start to end (Z) | Exit |
| --- | --- | --- |
| `go build ./...` | 08:26:22 to 08:26:23 | 0 |
| `go vet ./...` | 08:26:23 to 08:26:23 | 0 |
| `go vet -tags=integration ./...` | 08:26:23 to 08:26:23 | 0 |
| `go vet -tags=e2e ./...` | 08:26:23 to 08:26:24 | 0 |
| `gofmt -l .` | 08:26:24 | 0 files |
| `go test -count=1 ./...` | 08:26:24 to 08:26:25 | 0 |
| `TMPDIR=/home/layup/.cache/layup-session/tmp-green go test -count=1 -tags=integration -timeout 30m ./...` (the full integration level, in the background, output to `/home/layup/.cache/layup-session/T-9nbf-full-integration.txt`) | 08:26:30 to 08:30:37 | 0 |
| `go test -count=1 -tags=e2e ./...` | 08:30:48 to 08:31:54 | 1: `TestTheStepsOfTheScaffold`, `shape` fail `forge: the installation lacks the permissions of the task loop: pull_requests: write, statuses: write, checks: read` (the fake of the e2e level granted the three of M2a) |
| `go test -count=1 -tags=e2e ./...` after the fake grants the permissions of the loop | 08:33:20 to 08:34:26 | 0 |
| `sh docs/adr/adr-lint.sh`, `sh docs/prd/prd-lint.sh`, `sh docs/links/link-lint.sh`, `sh docs/tests/run-discipline-tests.sh`, `sh docs/tests/nested-checkout-check.sh`, `sh docs/setup/setup-check.sh`, `sh docs/tests/freeze-check/run-all.sh` | 08:31:54 to 08:32:45 | 0 each |
| `git diff --check origin/main HEAD` | 08:32:45 | 0 |
| `sh docs/tasks/task-state.sh` | 08:32:45 to 08:32:54 | 0 (T-9nbf `in progress`) |

The full integration level was run on the code head before the change of the e2e fake, which no integration test reads; the unit and integration levels of `cmd/layup` are in that run.

## S8, final. The head `0db3555c` (the merge of `origin/main` at `c0f691e` included), each command alone, output to a file outside the worktree

| Command | Start to end (Z) | Exit |
| --- | --- | --- |
| `go build ./...` | 08:36:25 to 08:36:25 | 0 |
| `go vet ./...` | 08:36:25 to 08:36:25 | 0 |
| `go vet -tags=integration ./...` | 08:36:25 to 08:36:26 | 0 |
| `go vet -tags=e2e ./...` | 08:36:26 to 08:36:26 | 0 |
| `go test -count=1 ./...` | 08:36:26 to 08:36:27 | 0 |
| full integration | 08:36:27 to 08:40:34 | 0 |
| `go test -count=1 -tags=e2e ./...` | 08:40:34 to 08:41:40 | 0 |
| `sh docs/adr/adr-lint.sh` | 08:41:40 to 08:41:42 | 0 |
| `sh docs/prd/prd-lint.sh` | 08:41:42 to 08:41:42 | 0 |
| `sh docs/links/link-lint.sh` | 08:41:42 to 08:41:56 | 0 |
| `sh docs/tests/run-discipline-tests.sh` | 08:41:56 to 08:41:59 | 0 |
| `sh docs/tests/nested-checkout-check.sh` | 08:41:59 to 08:41:59 | 0 |
| `sh docs/setup/setup-check.sh` | 08:41:59 to 08:42:22 | 0 |
| `sh docs/tests/freeze-check/run-all.sh` | 08:42:22 to 08:42:31 | 0 |
| `git diff --check origin/main HEAD` | 08:42:31 to 08:42:31 | 0 |
| `sh docs/tasks/task-state.sh` | 08:42:31 to 08:42:39 | 0 |

The full integration level is the command `TMPDIR=/home/layup/.cache/layup-session/tmp-green go test -count=1 -tags=integration -timeout 30m ./...`, output in `/home/layup/.cache/layup-session/T-9nbf-full-integration.txt`. `gofmt -l .` listed 0 files.

### The pre-freeze check

`sh docs/tests/freeze-check.sh --task T-9nbf --base origin/main --budget 2600/28 --evidence runs/T-9nbf/test-runs.md --evidence runs/T-9nbf/mutants.md --evidence runs/T-9nbf/conditions.md`, run on the clean committed head 0db3555c at 08:35Z (output written to `/home/layup/.cache/layup-session/T-9nbf-freeze-check.txt` first), exit 0:

```text
ok K1: white space
ok K2: registration of new tests
ok K3: test selectors
ok K4: restore evidence of mutants
ok K5: projection: figure 1672/27, maximum 2600/28, tolerance 500/4, limit 3100/32, margin 1428/5 (the figure at landing is still measured at landing)
ok K6: conditions table
```

## Round 1 of #367 (6312207004): the findings, test first (late evidence, after the freeze at `30ca2a8`)

Go 1.26.10 in every command. Each red ran before its code change.

### Finding 1: a successful unresolved close-out push when the base moved

- Test: `TestASuccessfulCloseOutPushWhoseAnswerWasLostIsBoundByTheRestart` (both cases, the base current and the base moved by a file the task does not touch, through `afterScaffold` and a second restart).
- Command: `TMPDIR=/home/layup/.cache/layup-session/tmp-green go test -count=1 -tags=integration ./internal/run/ -run '^(TestASuccessfulCloseOutPushWhoseAnswerWasLostIsBoundByTheRestart|TestAnAttemptAfterTheBetWithNoRenderedCommitIsTheStopOfTheOperator|TestAReopenedTaskReachesItsVerifierWithoutTheWaitingSpecCheckAndTheJobStillBlocksItsMerge)$'`.
- 09:05:52 to 09:06:09Z, exit 1. First failing assertion (the current-base case passes, as the reviewer said): `moved true: the restart = [shape done done before, bet done done before, loop fail the push of the session's head was refused: task/T-.../1 holds "...", not the head ... or its merge ...]`.
- Fix: in the act, a task whose view is verified has `resolvePushes` called and its view read again before `BringUpToDate`; `loop.md` item 3 states the order. Green after it, with one change of the scripted harness: the close-out step now writes a task file and a completed log that differ at each run (a second close-out at the merge that carries the head wrote the same bytes and was refused for its scope, an empty change; a real close-out at a new base has a change). The mutant M19 kills it.

### Finding 2: G1 (D1 in the step `bet`)

`TestAnAttemptAfterTheBetWithNoRenderedCommitIsTheStopOfTheOperator` now runs the recovered attempt through `afterScaffold` (two restarts), asserting the rows `shape` and `bet` `done before`, the exact stop and no new session. It passes at the code of `30ca2a8` (the reviewer's probe did too), so its redness is the mutant: G1 is killed (`runs/T-9nbf/mutants.md`; its earlier "equivalent" was wrong, since `bet` read `verifiedStep` before its row of `bets.tsv`).

### Finding 3: G5 (the exemption of a reopened task in the driver)

`TestAReopenedTaskReachesItsVerifierWithoutTheWaitingSpecCheckAndTheJobStillBlocksItsMerge` makes the pinned `prd-lint.sh` of the default branch fail; the reopened task reaches its verifier, and a failing required job still blocks its merge (the wait "not ready"). It passes at the code of `30ca2a8`; G5 is killed by it.

### Notes

- The comment of `loopPermissionsStep` (`register.go`) now says that the steps of the restart call the check (D6).
- F4's direct assertion is narrowed and says so: with no pull request for attempt 2, `FinishTask` fails at `attemptPull` (`is not on the forge`), before the rendered-head error can be reached; the test asserts that error. The stop of Q2 is the loop's and comes before any call of `FinishTask`.

### Mutants of the round (late evidence)

M19 (finding 1), G1 (finding 2) and G5 (finding 3) are killed; G11 (the reread of the view after the pushes are resolved) is not, and is recorded as such. Every restore matched its hash and every control was green (`runs/T-9nbf/mutants.md`, 30 mutants, 29 killed). The line "2 not killed (G1 equivalent, G5 not reached)" in S6 above is superseded by it.

### Gates of the round, on the head `9dfea6da` (late evidence; each command alone, output outside the worktree)

| Command | Start to end (Z) | Exit |
| --- | --- | --- |
| `go build ./...` | 09:10:57 to 09:10:57 | 0 |
| `go vet ./...` | 09:10:57 to 09:10:57 | 0 |
| `go vet -tags=integration ./...` | 09:10:57 to 09:10:57 | 0 |
| `go vet -tags=e2e ./...` | 09:10:57 to 09:10:58 | 0 |
| `go test -count=1 ./...` | 09:10:58 to 09:10:59 | 0 |
| full integration | 09:10:59 to 09:15:08 | 0 |
| `go test -count=1 -tags=e2e ./...` | 09:15:08 to 09:16:14 | 0 |
| `sh docs/adr/adr-lint.sh` | 09:16:14 to 09:16:16 | 0 |
| `sh docs/prd/prd-lint.sh` | 09:16:16 to 09:16:16 | 0 |
| `sh docs/links/link-lint.sh` | 09:16:16 to 09:16:30 | 0 |
| `sh docs/tests/run-discipline-tests.sh` | 09:16:30 to 09:16:33 | 0 |
| `sh docs/tests/nested-checkout-check.sh` | 09:16:33 to 09:16:33 | 0 |
| `sh docs/setup/setup-check.sh` | 09:16:33 to 09:16:57 | 0 |
| `sh docs/tests/freeze-check/run-all.sh` | 09:16:57 to 09:17:06 | 0 |
| `git diff --check origin/main HEAD` | 09:17:06 to 09:17:06 | 0 |

The full integration level: `TMPDIR=/home/layup/.cache/layup-session/tmp-green go test -count=1 -tags=integration -timeout 30m ./...`, on the code head (the commits after it change records and documents). The pre-freeze check follows.

Run on the clean committed head `2941df0c`, merged with `origin/main` at `c0f691e`, exit 0:

```text
ok K1: white space
ok K2: registration of new tests
ok K3: test selectors
ok K4: restore evidence of mutants
ok K5: projection: figure 1773/28, maximum 2600/28, tolerance 500/4, limit 3100/32, margin 1327/4 (the figure at landing is still measured at landing)
ok K6: conditions table
```
