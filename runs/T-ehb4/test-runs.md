# T-ehb4 - the test runs of row 66y

## S0 - the reconciliation of the plan (6262348152), its review (6262459073) and the decisions (6262588229, 6264869454) with the landed code

Read at `origin/main` = `df9bfac` (the merge of #448, `T-2vva`), merged into the branch (the branch had no work). Each name was checked in the tree, not in a plan or a review (note 4 of 6262459073).
All After edges that hold rows are on main: 66f, 66g, 66j, 66q, 66r, 66s. Row 66x (K3, `T-89xf`, #424) is not on main: the branch `T-89xf` is at `d9852d2`, so the up-to-date helper is named here for it (D9).

| Plan name | Landed form | Note |
| --------- | ----------- | ---- |
| 66f `PostStatus`, `HeadInput`, `PostStatuses`, `StatusPost`, `registerRevOf`, `VerifyAt` | `statuses.go:63`, `:185`, `:250`, `:123`; `verifystatus.go:33` | as planned; `PostChecks` (`statuses.go:209`) is the planned extraction, with `postChecks` keeping the delayed cleanup error apart |
| 66g `PushCloseOut`, `CarryInput`, `CheckRendered`, `CarryVerify`, `MergeTask`, `MarkReady`, `RecordMerged`, `lastBound` | `closeout.go:169`, `:120`, `:201`; `ready.go:195`, `:161`, `:255`; `pullrequest.go:87` | as planned; `MergeTask` marks the draft ready itself, and a closed pull request is `Readiness{Reason: "the pull request is not open"}` |
| 66s `PromptFile`, `PromptIn` | `promptfile.go:437` | the plan is the copy of the last valid comment handoff of the plan role, inlined in the loop of events: extracted as `planHandoff` and `copyPathOf` (a `refactor:` commit, guarded by the prompt-file and driver tests) |
| 66e `TierOf`, `Position`, `ResolveStep`, `ParsePlanComment` | `tier.go:50`, `:106`, `:170`; `commentform.go:128` | as planned; `Position` of a Shape task is always uphill (no frozen list), so the close-out of a Shape task is on the reasoning route; `TierOf(plan, position)` is exercised by the unit table |
| 66v `RepairTask` | `repair.go:30` | as planned; it writes before it reads a pull request, hence the read-only checks first (D7, condition 2) |
| 66r `Advance`, the verified-head predicate, the rerun, the base guard, `readOnlyStep` | `Advance` in `shapeloop.go:35`; the predicate is **`loopPosition`** (`shapeloop.go:266`) giving `posVerified` with `Head` and `Session`; `startAs` (`shapeloop.go:552`) with the scratch branch; `readOnlyStep` (`shapeloop.go`) does not cover `close-out` | there is no exported predicate and no `Verified` argument; `loopPosition` stops on a `bound` of a session that is not the developer's and on a `refused` with no start row, which the close-out's own events are, so `FinishTask` reads it over the events up to the last handoff of the verification role (`finishCut`, D1 position 5). `Advance` gains `Merged` and `Attempt` |
| 66r plan amendment 1 (the scratch base) | `startAs` sets `layup/base/<task>/<attempt>` by `git.UpdateRef` (`git.go:616`) | reused through `startAs`, which takes the push hook as one more parameter (the three callers pass `nil`); no second mechanism (6264869454) |
| 66j `BetInputs(ctx) (copy, records, preRev, registerRev, err)` | **absent**: 66j landed the step bet (`bet.go`) with `renderedStage` pushing `R` as code's `push`/`bound` (session `—`) and `CarryVerify` at `R`, and no reader of the typed inputs | the plan allows an extraction in place (S0, third bullet): `(*state).renderedBet` rebuilds the inputs from `bets.tsv`, the latest inventory version and the verified head, from `answerCopies`, `factsOf`, `inventoryOf` and `renderedAt`, extracted from `bet()` and `renderedStage` in a `refactor:` commit (guarded by `TestTheStepBet` and the atomicity tests). It is passed as `FinishInput.Bet` (D10), as 66k will pass it |
| 66j the consumers of a sessionless `bound` | `loopPosition` (`shapeloop.go:272`), `RepairPushes` (`repair.go:55`), `authorsFrom` (`authors.go:72`), `registerRevOf` (`statuses.go:123`), `lastBound` (`pullrequest.go:87`) skip `e[3] == ""` | as planned; `renderedPushed` (`bet.go:129`) is the driver's own test of `R` |
| 66q `LinkReviewRecord` | `reviewlink.go:31` | not called by 66y (Q12): 66r calls it before `Verified` |
| `forge.Fake` and `ClosePullRequest` | `forge/fake.go:304` | landed with K2 (66w); used by the unit tests |
| K3 (66x) up-to-date predicate | not landed | `(*Sessions).upToDate(ctx, base, head)` is the home (D9), named for K3 to consume with the head that carries the verdict |

### The divergences from the plan, and what they change

- **No `BetInputs`**: `renderedBet` is written here (D10) as an extraction of what the bet step computes. The scope grows by about 40 lines in `bet.go`; the budget of 2,300 lines holds it.
- **The predicate is `loopPosition` on a cut view**, not an exported function: the close-out's session, its bound and its refusals come after the verifier's last handoff and are cut off before the position is told. The plan's U5 and I3 positions are `closeOutPosition` over the events after the cut.
- **A close-out that ended with no push after a restart**: the run's clone is made anew at each restart (`cloneStep`), so the kept head `refs/layup/sessions/<id>` is gone. `FinishTask` resumes the push when the ref reads, and starts a new session when it does not (the old row stays, the last session decides). Decided here (D5), reported on #425.
- **The tier**: `Position` is uphill for every Shape task, so the integration test checks the reasoning route and the unit table the four combinations (the downhill route needs a frozen list, M2g).
- **Fixtures**: `ghPull` gains `closed`, `pullRoutes` gains `follow` (a pull request's head is the tip of its branch, as the forge's is), the check runs of the key `*` answer for every SHA, and the scripted close-out takes `closeout-extra` to touch a file outside its scope.

## S3 - the red runs

Every command was run directly; the status is the one `go` gave (`exit=`), and the start and the end are `date -u` read before and after the run. The runs are of the final test files against the head `6bdf16b` (the inert stubs of
`FinishTask` and its helpers, with `PostChecks` and the extractions real), in a copy of that tree made with `git archive 6bdf16b` outside the worktree, to which the five test files of this row (`finish_test.go`, `finish_integration_test.go`
and the three fixture files) were copied. An earlier run of nine unit and nine integration tests in the worktree at the same head failed the same way before the tests of the stop at the merge, the kept head, the missing verified head and the rendered inputs were added; those are in this run.
The long lines are cut at 260 characters; nothing else is changed. Each integration test reached its call: the setup (the task verified by the driver, the route of the close-out) succeeded.

### Red, the unit level

Command: `go test -count=1 ./internal/run/ -run '^(TestAClosedUnmergedPullRequestOfTheOpenAttemptIsAWaitWithNoWrite|TestAMergedPullRequestWithNoEventGetsItsMergedEvent|TestAMergedPullRequestWithItsEventWritesNothing|TestThePullRequestOfTheOpenAttemptMustBeTheRegisters|TestAnAttemptWithNoVerifierHandoffHasNoVerifiedHead|TestTheCloseOutStepIsResolvedAtTierOfThePlanAndThePosition|TestTheCutIsTheLastHandoffOfASessionOfTheVerificationRole|TestTheCloseOutSessionOfThePositionIsTheLastOneAfterTheCutAtTheCarriedBase|TestNotReadyIsAWaitAndAMovedHeadIsAStop|TestTheDefaultHeadMustBeTheBaseOrAnAncestorOfTheCarriedHead)$'` (start 2026-10-20T19:05:43Z, end 2026-10-20T19:05:45Z, exit=1)

| Test | First failing assertion |
| ---- | ----------------------- |
| `TestAClosedUnmergedPullRequestOfTheOpenAttemptIsAWaitWithNoWrite` | `finish_test.go:48: FinishTask = {Task: Verified:false Merged:false Head: Attempt:0 Wait: Stop:}, <nil>; want the wait "waiting for the Operator: T-ab12: the pull request #1 of task/T-ab12/2 is closed and not merged"` |
| `TestAMergedPullRequestWithNoEventGetsItsMergedEvent` | `finish_test.go:73: FinishTask = {Task: Verified:false Merged:false Head: Attempt:0 Wait: Stop:}, <nil>; want Merged at 2222222222222222222222222222222222222222` |
| `TestAMergedPullRequestWithItsEventWritesNothing` | `finish_test.go:93: FinishTask = {Task: Verified:false Merged:false Head: Attempt:0 Wait: Stop:}, <nil>; want Merged at 2222222222222222222222222222222222222222` |
| `TestThePullRequestOfTheOpenAttemptMustBeTheRegisters` | `finish_test.go:107: FinishTask = {Task: Verified:false Merged:false Head: Attempt:0 Wait: Stop:}, <nil>; want the error "T-ab12: the pull request 1 of task/T-ab12/2 is not the register's 5"` |
| `TestAnAttemptWithNoVerifierHandoffHasNoVerifiedHead` | `finish_test.go:124: FinishTask = <nil> with 0 commits more; want the error that attempt 2 has no verified head and no write` |
| `TestTheCloseOutStepIsResolvedAtTierOfThePlanAndThePosition` | `finish_test.go:145: large and downhill: closeOutStep = {Step: Role: Tier:}, <nil>; want {Step:close-out Role:software-developer Tier:execution}` |
| `TestTheCutIsTheLastHandoffOfASessionOfTheVerificationRole` | `finish_test.go:172: finishCut = 0; want 2, the handoff of the verifier` |
| `TestTheCloseOutSessionOfThePositionIsTheLastOneAfterTheCutAtTheCarriedBase` | `finish_test.go:216: ended with no push: closeOutPosition = {Pos:0 Session: Head: Msg:}; want {Pos:1 Session:S-c0000001 Head: Msg:}` |
| `TestNotReadyIsAWaitAndAMovedHeadIsAStop` | `finish_test.go:243: notReady({Ready:false Reason:layup/gates: pending}) = {Task: Verified:false Merged:false Head: Attempt:0 Wait: Stop:}; want {Task:T-ab12 Verified:false Merged:false Head: Attempt:0 Wait:waiting: T-ab12: not ready at 222222222222222222222222 ...` |
| `TestTheDefaultHeadMustBeTheBaseOrAnAncestorOfTheCarriedHead` | `finish_test.go:258: upToDate with the tip at the base = false, <nil>; want true` |

Green at the stub, each beside its red control: the cases `none`, `running`, `another base`, `started before the cut` and `the last decides` of `TestTheCloseOutSessionOfThePositionIsTheLastOneAfterTheCutAtTheCarriedBase` (the zero `closeStep` is a start) are green negatives,
with the cases `ended with no push`, `pushed and bound`, `refused`, `pushed with no bound`, `ended failed` and `a start refused before its row` as their red controls; in `TestAClosedUnmergedPullRequestOfTheOpenAttemptIsAWaitWithNoWrite`
the counts of commits and operations and the loop that allows reads only are green negatives beside the red wait.

### Red, the integration level

Command: `TMPDIR=/home/layup/.cache/layup-session/tmp-green go test -count=1 -tags=integration -timeout 30m ./internal/run/ -run '^(TestAVerifiedTaskIsClosedOutCarriedAndMergedAtItsHead|TestARestartAfterEachWriteAdoptsIt|TestAClosedUnmergedPullRequestIsAWaitWithNoWriteAndARepairThatStartsAnAttemptIsAHandback|TestACloseOutOutsideItsScopeIsRefusedAndNotRetried|TestNotReadyWaitsAndARestartMerges|TestNoVerifyIsComputedAtTheCloseOutHead|TestADefaultBranchThatMovedStopsTheCloseOutBeforeAnyWrite|TestACloseOutThatEndedWithNoPushResumesItsCheckedPush|TestTheRenderedHeadGetsItsStatusesAndTheCloseOutGoesOnIt|TestACloseOutWhoseKeptHeadIsGoneIsReplacedByANewSession|TestAPullRequestClosedDuringTheCallIsTheWaitOfTheOperator)$'` (start 2026-10-20T19:05:45Z, end 2026-10-20T19:05:52Z, exit=1)

| Test | First failing assertion |
| ---- | ----------------------- |
| `TestAVerifiedTaskIsClosedOutCarriedAndMergedAtItsHead` | `finish_integration_test.go:125: FinishTask = {Task: Verified:false Merged:false Head: Attempt:0 Wait: Stop:}, <nil>; want Merged at the bound head of the close-out ""` |
| `TestARestartAfterEachWriteAdoptsIt` | `finish_integration_test.go:217: call 1 = {Task: Verified:false Merged:false Head: Attempt:0 Wait: Stop:}; want its stop` |
| `TestAClosedUnmergedPullRequestIsAWaitWithNoWriteAndARepairThatStartsAnAttemptIsAHandback` | `finish_integration_test.go:261: the first call gave no error; want the push that went down` |
| `TestACloseOutOutsideItsScopeIsRefusedAndNotRetried` | `finish_integration_test.go:295: FinishTask = {Task: Verified:false Merged:false Head: Attempt:0 Wait: Stop:}, <nil>; want the stop that names the refused session (scope)` |
| `TestNotReadyWaitsAndARestartMerges` | `finish_integration_test.go:314: FinishTask = {Task: Verified:false Merged:false Head: Attempt:0 Wait: Stop:}, <nil>; want the wait "waiting: T-u8ry: not ready at : job build: in_progress"` |
| `TestNoVerifyIsComputedAtTheCloseOutHead` | `finish_integration_test.go:331: FinishTask = {Task: Verified:false Merged:false Head: Attempt:0 Wait: Stop:}, <nil>; want Merged at` |
| `TestADefaultBranchThatMovedStopsTheCloseOutBeforeAnyWrite` | `finish_integration_test.go:370: FinishTask = {Task: Verified:false Merged:false Head: Attempt:0 Wait: Stop:}, <nil>; want the error "T-f42o: the base moved (row 66x's)"` |
| `TestACloseOutThatEndedWithNoPushResumesItsCheckedPush` | `finish_integration_test.go:391: the first call gave no error; want the stop at the push` |
| `TestTheRenderedHeadGetsItsStatusesAndTheCloseOutGoesOnIt` | `finish_integration_test.go:413: renderedBet before the bet gave no error` |
| `TestACloseOutWhoseKeptHeadIsGoneIsReplacedByANewSession` | `finish_integration_test.go:482: the first call gave no error; want the stop at the push` |
| `TestAPullRequestClosedDuringTheCallIsTheWaitOfTheOperator` | `finish_integration_test.go:514: FinishTask = {Task: Verified:false Merged:false Head: Attempt:0 Wait: Stop:}, <nil>; want the wait "waiting for the Operator: T-fj7k: the pull request #1 of task/T-fj7k/1 is closed and not merged" and no merge` |

Four tests fail on their first stop (`want its stop`, `the first call gave no error`): the stub returns an empty `Advance` and nothing runs, so the later asserts of those tests (the order of the writes, the single session, the adoption after each write) have no red of their own; they are shown by the mutants.

## S4 - green

The head is the one of the commit that holds this record: the code of `finish.go` for the stubs, with the same test files. The same commands as above, run directly in the worktree.

- Unit: `go test -count=1 ./internal/run/ -run '^(TestAClosedUnmergedPullRequestOfTheOpenAttemptIsAWaitWithNoWrite|TestAMergedPullRequestWithNoEventGetsItsMergedEvent|TestAMergedPullRequestWithItsEventWritesNothing|TestThePullRequestOfTheOpenAttemptMustBeTheRegisters|TestAnAttemptWithNoVerifierHandoffHasNoVerifiedHead|TestTheCloseOutStepIsResolvedAtTierOfThePlanAndThePosition|TestTheCutIsTheLastHandoffOfASessionOfTheVerificationRole|TestTheCloseOutSessionOfThePositionIsTheLastOneAfterTheCutAtTheCarriedBase|TestNotReadyIsAWaitAndAMovedHeadIsAStop|TestTheDefaultHeadMustBeTheBaseOrAnAncestorOfTheCarriedHead)$'` (start 2026-10-20T19:06:10Z, end 2026-10-20T19:06:11Z, exit=0): ok.
- Integration: `TMPDIR=/home/layup/.cache/layup-session/tmp-green go test -count=1 -tags=integration -timeout 30m ./internal/run/ -run '^(TestAVerifiedTaskIsClosedOutCarriedAndMergedAtItsHead|TestARestartAfterEachWriteAdoptsIt|TestAClosedUnmergedPullRequestIsAWaitWithNoWriteAndARepairThatStartsAnAttemptIsAHandback|TestACloseOutOutsideItsScopeIsRefusedAndNotRetried|TestNotReadyWaitsAndARestartMerges|TestNoVerifyIsComputedAtTheCloseOutHead|TestADefaultBranchThatMovedStopsTheCloseOutBeforeAnyWrite|TestACloseOutThatEndedWithNoPushResumesItsCheckedPush|TestTheRenderedHeadGetsItsStatusesAndTheCloseOutGoesOnIt|TestACloseOutWhoseKeptHeadIsGoneIsReplacedByANewSession|TestAPullRequestClosedDuringTheCallIsTheWaitOfTheOperator)$'` (start 2026-10-20T19:06:11Z, end 2026-10-20T19:06:20Z, exit=0): ok.
- Characterization of the extractions (`PostChecks`, `planHandoff`, `startAs`, the bet step's reads): the tests of the step bet, the statuses, the prompt file, the pull request, the question loop, the Shape driver and the repairs (those with Bet, Rendered, Status, Prompt, Ship, Pull, Question, Shape or Repair in their names) were run at the commit `cae4328` of the extractions (start 2026-10-20T18:54:11Z, end 2026-10-20T18:55:24Z, exit=0, 71 s) and again in the full integration level below. The same tests were green on `origin/main` (CI of #448).

## S7 - the gates

The code is frozen at `02a5f0f`; the commits after it hold only records and tables. Each command was run directly, with its own exit status; `origin/main` was fetched again before the gates and had not moved (`df9bfac`), so there is no merge to rerun on.

- `go build ./...`, `go vet ./...`, `go vet -tags=integration ./...`, `go vet -tags=e2e ./...`, `gofmt -l .` (empty), `go test -count=1 ./...` (21 packages ok), `git diff --check origin/main HEAD`: exit=0 each.
- The full integration level: `TMPDIR=/home/layup/.cache/layup-session/tmp-green go test -count=1 -tags=integration -timeout 30m ./...` (head `02a5f0f`, start 2026-10-20T19:20:25Z, end 2026-10-20T19:23:41Z, exit=0), output in `/home/layup/.cache/layup-session/T-ehb4-full-integration.txt`: 21 packages ok, `internal/run` and `cmd/layup` among them.
- The end-to-end level: `go test -count=1 -tags=e2e ./...` (start 2026-10-20T19:34:36Z, end 2026-10-20T19:35:45Z, exit=0): ok.
- `sh docs/adr/adr-lint.sh`, `sh docs/prd/prd-lint.sh`, `sh docs/links/link-lint.sh`, `sh docs/tests/run-discipline-tests.sh`, `sh docs/setup/setup-check.sh`, `sh docs/tests/nested-checkout-check.sh`: exit=0 each. `bash docs/tasks/task-state.sh`: exit=0, `T-ehb4 #425 stalled` before the push (the branch has no pull request and no push since the claim).

## The freeze check

Command: `sh docs/tests/freeze-check.sh --task T-ehb4 --base origin/main --budget 2300/22`, run on the clean committed head `3437207` (the output is kept outside the worktree, in `/home/layup/.cache/layup-session/T-ehb4-freeze.txt`), exit=0:

```text
ok K1: white space
ok K2: registration of new tests
ok K3: test selectors
ok K4: restore evidence of mutants
ok K5: projection: figure 1801/22, maximum 2300/22, margin 499/0 (the figure at landing is still measured at landing)
ok K6: conditions table
```

The figure is the added plus removed lines of `git diff --numstat origin/main HEAD` with the 70 lines and 2 files that the check keeps for the close-out; the close-out of this row is not written here.
