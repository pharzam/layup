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

Command: `go test -count=1 ./internal/run/ -run '^(TestAClosedUnmergedPullRequestOfTheOpenAttemptIsAWaitWithNoWrite|TestAMergedPullRequestWithNoEventGetsItsMergedEvent|TestAMergedPullRequestWithItsEventWritesNothing|TestThePullRequestOfTheOpenAttemptMustBeTheRegisters|TestAnAttemptWithNoVerifierHandoffHasNoVerifiedHead|TestTheCloseOutStepIsResolvedAtTierOfThePlanAndThePosition|TestTheCutIsTheLastHandoffOfASessionOfTheVerificationRole|TestTheCloseOutSessionOfThePositionIsTheLastOneAfterTheCutAtTheCarriedBase|TestNotReadyIsAWaitAndAMovedHeadIsAStop)$'` (start 2026-10-20T19:05:43Z, end 2026-10-20T19:05:45Z, exit=1)

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

Command: `TMPDIR=/home/layup/.cache/layup-session/tmp-green go test -count=1 -tags=integration -timeout 30m ./internal/run/ -run '^(TestAVerifiedTaskIsClosedOutCarriedAndMergedAtItsHead|TestARestartAfterEachWriteAdoptsIt|TestAClosedUnmergedPullRequestIsAWaitWithNoWriteAndARepairThatStartsAnAttemptIsAHandback|TestACloseOutOutsideItsScopeIsRefusedAndNotRetried|TestNotReadyWaitsAndARestartMerges|TestNoVerifyIsComputedAtTheCloseOutHead|TestADefaultBranchThatMovedStopsTheCloseOutBeforeAnyWrite|TestACloseOutThatEndedWithNoPushResumesItsCheckedPush|TestTheRenderedHeadGetsItsStatusesAndTheCloseOutGoesOnIt|TestAPullRequestClosedDuringTheCallIsTheWaitOfTheOperator)$'` (start 2026-10-20T19:05:45Z, end 2026-10-20T19:05:52Z, exit=1)

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

- Unit: `go test -count=1 ./internal/run/ -run '^(TestAClosedUnmergedPullRequestOfTheOpenAttemptIsAWaitWithNoWrite|TestAMergedPullRequestWithNoEventGetsItsMergedEvent|TestAMergedPullRequestWithItsEventWritesNothing|TestThePullRequestOfTheOpenAttemptMustBeTheRegisters|TestAnAttemptWithNoVerifierHandoffHasNoVerifiedHead|TestTheCloseOutStepIsResolvedAtTierOfThePlanAndThePosition|TestTheCutIsTheLastHandoffOfASessionOfTheVerificationRole|TestTheCloseOutSessionOfThePositionIsTheLastOneAfterTheCutAtTheCarriedBase|TestNotReadyIsAWaitAndAMovedHeadIsAStop)$'` (start 2026-10-20T19:06:10Z, end 2026-10-20T19:06:11Z, exit=0): ok.
- Integration: `TMPDIR=/home/layup/.cache/layup-session/tmp-green go test -count=1 -tags=integration -timeout 30m ./internal/run/ -run '^(TestAVerifiedTaskIsClosedOutCarriedAndMergedAtItsHead|TestARestartAfterEachWriteAdoptsIt|TestAClosedUnmergedPullRequestIsAWaitWithNoWriteAndARepairThatStartsAnAttemptIsAHandback|TestACloseOutOutsideItsScopeIsRefusedAndNotRetried|TestNotReadyWaitsAndARestartMerges|TestNoVerifyIsComputedAtTheCloseOutHead|TestADefaultBranchThatMovedStopsTheCloseOutBeforeAnyWrite|TestACloseOutThatEndedWithNoPushResumesItsCheckedPush|TestTheRenderedHeadGetsItsStatusesAndTheCloseOutGoesOnIt|TestAPullRequestClosedDuringTheCallIsTheWaitOfTheOperator)$'` (start 2026-10-20T19:06:11Z, end 2026-10-20T19:06:20Z, exit=0): ok.
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

## Round 2 - the four findings of round 1 (review.md of issue425-r1), red first

The head is `646691b` (`origin/main` d27b331 merged), the tests below uncommitted. Each command was run directly; the statuses and the times are those of the files named.

- Finding 3, unit: `go test -count=1 ./internal/run/ -run '^(TestAMergedPullRequestIsRecordedAndStartsNoReplacementAttempt)$'` (start 2026-10-20T19:53:30Z, end 2026-10-20T19:53:32Z, exit=1). The first form of the test failed on the App's slug that a repair needs; the slug was set, and the failure is the one of the review: a replacement attempt.
- Findings 1 and 4, integration, on the code of `02a5f0f` (with the inert field `FinishInput.Copy` and `betCopy`, so that the file compiles): `TMPDIR=/home/layup/.cache/layup-session/tmp-green go test -count=1 -tags=integration -timeout 30m ./internal/run/ -run '^(TestACloseOutThatEndedWithNoPushResumesAfterTheCloneIsMadeAnew|TestACloseOutWhoseKeptHeadIsGoneIsAnErrorAndNoSecondSession|TestAnArchitectureCloseOutAfterTheBetCarriesTheAnswerCopy|TestAPullRequestClosedAtTheRefusedMergeIsTheWaitOfTheOperator)$'` (start 2026-10-20T19:53:18Z, end 2026-10-20T19:53:25Z, exit=1).
- Finding 2 and the new git call, in a copy of that head made with `git archive` outside the worktree, with the stubs `Refs`, `betCopy` and the field `ShapeTaskInput.Copy` (inert): `go test -count=1 ./internal/git/ -run '^TestRefsListsTheRefsBelowAPrefixAndRefusesAnythingElse$'` and `TMPDIR=/home/layup/.cache/layup-session/tmp-green go test -count=1 -tags=integration -timeout 30m ./internal/run/ -run '^TestAnArchitectureCloseOutAfterTheBetCarriesTheAnswerCopy$'` (start 2026-10-20T19:56:36Z, end 2026-10-20T19:56:44Z, exit=1 each). An earlier form of the test of finding 2, on the architecture task of Shape, failed with `layup/spec: failure` and showed that such a head cannot pass the full check: the rendered commit is not in its tree, so the test uses the architecture task that the bet opens, based on the rendered commit, and the driver takes the copy too.

| Test | First failing assertion |
| ---- | ----------------------- |
| `TestAMergedPullRequestIsRecordedAndStartsNoReplacementAttempt` | `finish_test.go:302: FinishTask = {Task:T-ab12 Verified:false Merged:false Head: Attempt:3 Wait: Stop:}, <nil>; want Merged at 2222222222222222222222222222222222222222 with no new attempt` |
| `TestAPullRequestClosedAtTheRefusedMergeIsTheWaitOfTheOperator` | `finish_integration_test.go:564: FinishTask = {Task:T-po3z Verified:false Merged:false Head: Attempt:0 Wait:waiting: T-po3z: not ready at 27fde9ba09b904ce63c7195e3eb14057cdb018d2:  Stop:}, the merge of the pull request 1 at 27fde9ba09b904ce63c7195e3eb14057cdb018d2 was refused and its head did not mov ...` |
| `TestACloseOutWhoseKeptHeadIsGoneIsAnErrorAndNoSecondSession` | `finish_integration_test.go:525: FinishTask = <nil> with 2 sessions; want the error that names S-c1ce7a15 and no second session` |
| `TestACloseOutThatEndedWithNoPushResumesAfterTheCloneIsMadeAnew` | `finish_integration_test.go:503: the second call = {Task:T-1jgl Verified:false Merged:true Head:50ac34ba03dbc003535bbc541c736979dd2a84c3 Attempt:0 Wait: Stop:}, <nil> with the close-out sessions [[S-685e8f86 T-1jgl 1 software-developer claude 1.0.0 m1 904f248c2b37d7db931483e0beeb60883c8a6265 e56f3967 ...` |
| `TestAnArchitectureCloseOutAfterTheBetCarriesTheAnswerCopy` | `finish_integration_test.go:537: FinishTask = {Task:T-8h6p Verified:false Merged:false Head: Attempt:0 Wait:waiting: T-8h6p: not ready at 8a2df8ec762c24bcdcc51f6fdeae59b8c77661d0: layup/spec: failure Stop:}, <nil>; want Merged at the head of the close-out ([layup/rules success layup/spec failure layu ...` |
| `TestRefsListsTheRefsBelowAPrefixAndRefusesAnythingElse` | `git_test.go:614: Refs = [], <nil>; want the two refs` |
| `TestAnArchitectureCloseOutAfterTheBetCarriesTheAnswerCopy` | `finish_integration_test.go:541: betCopy = "", <nil>; want the answer's copy` |

`TestACloseOutWhoseKeptHeadIsGoneIsReplacedByANewSession` of round 1 is replaced by the first two tests of finding 1: it required two producing sessions, which condition 3 of the plan review forbids.

### Round 2 - green and the gates

- Green: the same commands as above on the fixes, exit=0: the unit set (12 tests) and the integration set (15 tests) of `/home/layup/.cache/layup-session/T-ehb4-un.txt` and `-in.txt` ran in the mutant controls (exit 0 in each of 40 controls) and in the full integration level. `TestAMergedPullRequestWhosePushWasNeverBoundGetsItsBoundAndItsMergedEvent` and `TestAFailedReadOfTheKeptHeadsIsTheErrorOfTheReadAndNotAnAbsence` have no red transcript: they were written after the code (late evidence), and their mutants M36 and M34 show what they hold.
- The code is frozen at `dcea749`; the head `7cc2f8a` and the later ones add tests, tables and records only. `origin/main` (d27b331) was merged first (`646691b`) and had not moved at the gates.
- Full integration level: `TMPDIR=/home/layup/.cache/layup-session/tmp-green go test -count=1 -tags=integration -timeout 30m ./...` (start 2026-10-20T20:19:14Z, end 2026-10-20T20:22:24Z, exit=0), 21 packages ok, on the head `f4c7da1`; an earlier run on `dcea749` (19:58:49Z to 20:02:15Z) was exit=0 too.
- `go build ./...`, `go vet ./...` with no tag, `integration` and `e2e`, `gofmt -l .` (empty), `go test -count=1 ./...` (21 ok), `go test -count=1 -tags=e2e ./...` (start 2026-10-20T20:22:31Z, end 2026-10-20T20:23:36Z, exit=0), the six lints, `bash docs/tasks/task-state.sh` and `git diff --check origin/main HEAD`: exit=0 each.

### Round 2 - the freeze check

Command: `sh docs/tests/freeze-check.sh --task T-ehb4 --base origin/main --budget 2300/22`, on the clean committed head `a38488a` (output kept in `/home/layup/.cache/layup-session/T-ehb4-freeze2.txt`), exit=1 on K5 alone:

```text
ok K1: white space
ok K2: registration of new tests
ok K3: test selectors
ok K4: restore evidence of mutants
FAIL K5: -:0: projection: figure 2209/27, maximum 2300/22, margin 91/-5 (the figure at landing is still measured at landing)
ok K6: conditions table
```

K5 of this script still counts the files of `runs/` (3 files here: test-runs.md, mutants.md and conditions.md) and the 70 lines and 2 files of the close-out allowance. Measured as ADR-0028 reads it (`git diff --numstat origin/main HEAD` without `runs/`): 1,920 lines (1,838 added, 82 removed) over 22 files, with the close-out of 2 more files (the task file and the completed log) 24 files against the maximum of 22: the lines fit, the files exceed the maximum by 2. The cut is the Operator's. Run with `--budget 2300/32`, K1 to K6 are ok.

## Round 3 - the two findings of round 2 (review.md of issue425-r2)

### Red

The head is `49855aa`, the tests uncommitted. Command: `TMPDIR=/home/layup/.cache/layup-session/tmp-green go test -count=1 -tags=integration -timeout 30m ./internal/run/ -run '^(TestARestartOfATargetTouchesNoOtherTargetsCloneAndKeepsItsOwnCandidate|TestARestartThatStopsAfterTheSetAsideStillFindsTheKeptHeads)$'` (start 2026-10-20T20:56:55Z, end 2026-10-20T20:56:59Z, exit=1).

| Test | First failing assertion |
| ---- | ----------------------- |
| `TestARestartOfATargetTouchesNoOtherTargetsCloneAndKeepsItsOwnCandidate` | `finish_integration_test.go:538: the clone of acme/target.prev after the restart of acme/target: "", open .../targets/acme/target.prev/mine.txt: no such file or directory; want its file and its HEAD unchanged` |

`TestARestartThatStopsAfterTheSetAsideStillFindsTheKeptHeads` is green at `49855aa` (a green negative: the first form of the fix already finds a clone set aside by an earlier restart); it guards the interrupted migration of the new form, and its mutant shows it. A first form of it used a cancelled context, which the clone step does not honour; the forge's token call fails instead.
Finding 6 (the specification) has no test: it is a clause-by-clause reading of `docs/spec/run.md` step 2 and D5.

### Green

The same command on the fix (`cloneStep` sets the old clone aside in `DIR/recovery/OWNER/NAME/`, which is no target's clone; `docs/spec/run.md` step 2 and D5 of `loop.md` state one contract), start 2026-10-20T20:57:58Z, end 2026-10-20T20:58:03Z, exit=0: ok.

### Round 3 - the gates

The code is frozen at `b74fb19`; later commits add tables and records only. `origin/main` had not moved. Full integration level: `TMPDIR=/home/layup/.cache/layup-session/tmp-green go test -count=1 -tags=integration -timeout 30m ./...` (head `5473bd4`, start 2026-10-20T21:02:02Z, end 2026-10-20T21:05:44Z, exit=0), 21 packages ok. `go build ./...`, `go vet ./...` with no tag, `integration` and `e2e`, `gofmt -l .` (empty), `go test -count=1 ./...` (21 ok), `go test -count=1 -tags=e2e ./...` (start 2026-10-20T21:05:47Z, end 2026-10-20T21:06:54Z, exit=0), the six lints and `git diff --check origin/main HEAD`: exit=0 each.

### Round 3 - the freeze check

Command: `sh docs/tests/freeze-check.sh --task T-ehb4 --base origin/main --budget 2300/22` on the clean committed head `99f85f4` (output in `/home/layup/.cache/layup-session/T-ehb4-freeze3.txt`), exit=1 on K5 alone, which counts the files of `runs/` and the old close-out allowance:

```text
ok K1: white space
ok K2: registration of new tests
ok K3: test selectors
ok K4: restore evidence of mutants
FAIL K5: -:0: projection: figure 2349/28, maximum 2300/22, margin -49/-6 (the figure at landing is still measured at landing)
ok K6: conditions table
```

Without `runs/` (ADR-0028): 1926 added, 84 removed, 2010 lines over 23 files, with the close-out of 2 more files still to come.

## Round 4 - finding 7 of round 3 (review.md of issue425-r3)

### Red

The code was written before the tests were run (the order of this finding broke the rule), so the red run is in a copy of the head `42153e8` made with `git archive` outside the worktree, to which the new test file was copied. Command: `TMPDIR=/home/layup/.cache/layup-session/tmp-green go test -count=1 -tags=integration -timeout 30m ./internal/run/ -run '^(TestACleanupThatFailsAfterTheImportLosesNoHead|TestACleanupInterruptedAfterAPartialDeleteLosesNoHead)$'` (start 2026-10-20T21:33:23Z, end 2026-10-20T21:33:27Z, exit=1).

| Test | First failing assertion |
| ---- | ----------------------- |
| `TestACleanupThatFailsAfterTheImportLosesNoHead` | `finish_integration_test.go:621: the clone step of the restart: {Name:clone Result:fail Detail:git for-each-ref --format=%(refname) refs/layup/sessions/: exit status 128: fatal: not a git repository (or any parent up to mount point /)}` |
| `TestACleanupInterruptedAfterAPartialDeleteLosesNoHead` | `finish_integration_test.go:639: the clone step of the restart: {Name:clone Result:fail Detail:git for-each-ref ... fatal: not a git repository ...}` |

The restart that follows the failed cleanup discarded the current clone, which held the imported head, and then read the damaged place. Case (a) removes write permission on the objects directories of the old clone (as the review did), so that the removal deletes the rest and stops; case (b) deletes `HEAD`, `refs` and `config` of the old clone as well, as an interrupted removal would.

### Green and the gates

Green: the same command as the red one, with the tests of finding 5, on `eb0e718` (start 2026-10-20T21:33:41Z, end 2026-10-20T21:33:46Z, exit=0): ok. The code is frozen at `eb0e718`; later commits add tables and records only; `origin/main` had not moved. Full integration level: `TMPDIR=/home/layup/.cache/layup-session/tmp-green go test -count=1 -tags=integration -timeout 30m ./...` (head `bc96d13`, start 2026-10-20T21:37:17Z, end 2026-10-20T21:40:57Z, exit=0), 21 packages ok. build, vet with no tag, `integration` and `e2e`, gofmt, `go test -count=1 ./...` (21 ok), `go test -count=1 -tags=e2e ./...` (start 2026-10-20T21:40:59Z, end 2026-10-20T21:42:06Z, exit=0), the six lints and `git diff --check origin/main HEAD`: exit=0 each.

### Round 4 - the freeze check

Command: `sh docs/tests/freeze-check.sh --task T-ehb4 --base origin/main --budget 2300/22` on the clean committed head `ea6b895` (output in `/home/layup/.cache/layup-session/T-ehb4-freeze4.txt`), exit=1 on K5 alone, which counts the files of `runs/`:

```text
ok K1: white space
ok K2: registration of new tests
ok K3: test selectors
ok K4: restore evidence of mutants
FAIL K5: -:0: projection: figure 2517/28, maximum 2300/22, margin -217/-6 (the figure at landing is still measured at landing)
ok K6: conditions table
```

Without `runs/T-ehb4`: 2050 added, 84 removed, 2134 lines over 23 files, the close-out of 2 more files still to come.

## Round 5 - finding 8 of round 4 (end 3 of the Operator, 6302892665)

### Red

Strict red on the unchanged code `c68e1a3`, the new test file alone, in the worktree. Command: `TMPDIR=/home/layup/.cache/layup-session/tmp-green go test -count=1 -tags=integration -timeout 30m ./internal/run/ -run '^(TestAPullRequestMergedOutsideBeforeTheClosingPushIsMergedAtItsOwnHead)$'` (start 2026-10-20T22:10:33Z, end 2026-10-20T22:10:37Z, exit=1).

| Test | First failing assertion |
| ---- | ----------------------- |
| `TestAPullRequestMergedOutsideBeforeTheClosingPushIsMergedAtItsOwnHead` | `finish_integration_test.go:689: the call that finds the pull request merged = {Task:T-holl Verified:false Merged:true Head:b367bf6c... Attempt:0}, <nil>; want Merged at the head 04e995b8... that was merged` (Head is the close-out head that was never merged) |

### Green

The fix: `mergeAt` gives the head of the pull request that the forge holds after `MergeTask` (`Head: pr.HeadSHA`), not the close-out head. The same command, the new test and the integration set of the round 1 to 4 tests, start 2026-10-20T22:10:44Z, end 2026-10-20T22:10:56Z, exit=0: ok.

### Round 5 - the merge of main, the fold and the gates

`origin/main` (4259e05, row 66x landed) was merged (`abd8166`). Conflicts: `internal/git/git.go`, `internal/git/git_test.go` (both sides appended: kept both), `internal/run/pullroutes_integration_test.go` (the field `closed` and the fields of the route: both kept), `docs/spec/loop.md` (D-text of this row and the section of row 66x: both kept), `docs/spec/packages.md` (the cell of `internal/run`: both additions kept), `docs/tests/traceability.md` and `docs/prd/PRD-0001-layup.md` (rows of both sides, each row of section 12 merged cell by cell); no `+` line of either side was dropped. The fold of the two predicates was made in the same commit as the merge (a deviation from the order asked): `FinishTask` calls `UpToDate(ctx, carried)` of `baseup.go` and `upToDate` of `finish.go` and its unit test `TestTheDefaultHeadMustBeTheBaseOrAnAncestorOfTheCarriedHead` were deleted; D9 of `loop.md` names the one predicate.

Full integration level: `TMPDIR=/home/layup/.cache/layup-session/tmp-green go test -count=1 -tags=integration -timeout 30m ./...` (head `bfba9ec`, start 2026-10-20T22:17:47Z, end 2026-10-20T22:21:09Z, exit=0), 21 packages ok. build, vet with no tag, `integration` and `e2e`, gofmt, `go test -count=1 ./...` (21 ok), `go test -count=1 -tags=e2e ./...` (start 2026-10-20T22:21:16Z, end 2026-10-20T22:22:21Z, exit=0), the six lints and `git diff --check origin/main HEAD`: exit=0 each. The code is frozen at `1fc7253` plus the merge and the fold (`abd8166`).

## Round 6 - finding 10 of round 5 (cap 5, the scoped round 6 follows)

### Red

Strict red on the unchanged code `c75505f`, the test file alone changed. Command: `TMPDIR=/home/layup/.cache/layup-session/tmp-green go test -count=1 -tags=integration -timeout 30m ./internal/run/ -run '^(TestAVerifiedTaskBroughtCurrentByAMergeIsClosedOutFromTheMerge|TestAVerifiedSpecificationBroughtCurrentAboveItsRenderedCommitIsClosedOutFromTheMerge)$'` (start 2026-10-20T22:38:00Z, end 2026-10-20T22:38:07Z, exit=1).

| Test | First failing assertion |
| ---- | ----------------------- |
| `TestAVerifiedTaskBroughtCurrentByAMergeIsClosedOutFromTheMerge` | `finish_integration_test.go:390: FinishTask = {...}, T-zjl4: the base moved (row 66x's); want Merged at the head of the close-out` |
| `TestAVerifiedSpecificationBroughtCurrentAboveItsRenderedCommitIsClosedOutFromTheMerge` | `finish_integration_test.go:416: FinishTask = {...}, 63333940... is not a permitted rendered commit above the verified head b3dde440...; want Merged at the head of the close-out` (the merge was taken for the rendered commit) |

### Green

The fix: `FinishTask` takes as the head that carries the verdict the last event `bound` of code of the attempt (the rendered commit, or K3's merge above V or above R), `V` with none; checks it with `carries`; asks `UpToDate` of it; and `renderedHead` finds R by its own record (the bound of code with the one parent V). The same command with the two new tests and the integration set of the earlier rounds, start 2026-10-20T22:38:43Z, end 2026-10-20T22:38:57Z, exit=0: ok. `TestTheRenderedHeadGetsItsStatusesAndTheCloseOutGoesOnIt` now meets the records of another rendering at the check `carries`, which names the verified head too, and its assertion was widened to that text.

### Round 6 - the gates

The code is frozen at \`b0bfde5\`; later commits add a test, tables and records only (\`origin/main\` had not moved). Full integration level: \`TMPDIR=/home/layup/.cache/layup-session/tmp-green go test -count=1 -tags=integration -timeout 30m ./...\` (head \`53e8588\`, start 2026-10-20T22:44:08Z, end 2026-10-20T22:47:44Z, exit=0), 21 packages ok. build, vet with no tag, \`integration\` and \`e2e\`, gofmt, \`go test -count=1 ./...\` (21 ok), \`go test -count=1 -tags=e2e ./...\` (start 2026-10-20T22:47:46Z, end 2026-10-20T22:48:51Z, exit=0), the six lints and \`git diff --check origin/main HEAD\`: exit=0 each.
