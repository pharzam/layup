# T-89xf (#424, row 66x): the runs

## S0, the reconciliation with main (df9bfac)

`git fetch https://github.com/pharzam/layup.git +refs/heads/main:refs/remotes/origin/main` then `git merge origin/main`: a fast-forward from d9852d2 to df9bfac (no work was on the branch), so no merge commit.

| Name in the plan, review or decisions | What main holds | Consequence |
| --- | --- | --- |
| `lastBound` (pullrequest.go:73) | skips an event `bound` with no session (66j) | new `TipOf` reads the last `bound` of any session |
| `RepairPushes` and sessionless `push` | skips a `push` with no session (66j: its step resolves it) | the push of the merge follows the rendered commit's pattern (`push`, branch check, `PushHead`, `bound`); no clause added to `CheckEvent` for it |
| `CheckEvent` clause for session-less push/bound | present (push and bound exempt from the session rule) | none |
| `Advance` | `Task, Verified, Head, Wait, Stop` | adds `Replaced` |
| 66r `taskID`, `openTask`, `shapeIssueText` | `taskID` reads sessions.tsv and tasks.tsv; composers present | reused; the arch issue block is extracted to `archIssueTitle`/`archIssueBody` |
| `RepairTask` calls | `AdvanceShapeTask` calls it first | `closeOlderPulls` is added to `RepairTask`; `finishRedo` is not needed (all forge acts precede the one records commit, recorded by an intent file) |
| `SessionStore` | no `FetchBranch`, no `ReadFiles` | both added (recordsStore, standinStore) |
| binding reader (66j) | none; `tasks/<t>/options/<d>.tsv` read with `ReadFiles` and `spec.ReadOptions` | the redo copies the bytes and reads them back |
| K4 (T-ehb4) | not on main | `UpToDate` and `TipOf` are the shared helpers, named in loop.md |
| forge reads of an issue | none | the issue text is composed from what the records keep |

Divergences from the plan, all decided in loop.md (A task whose base moved): the intent file `tasks/<old>/redo.txt` (plan review condition 3); a closed/merged pull request found at the comparison is a Wait/Stop; the close-out head under a merge is an eligible first parent (condition 7).

## Semantic red runs, against the stubs (nothing committed)

The stubs: `git.Merge`/`git.MergeBase` return zero values, `FetchBranch`, `BringUpToDate`, `UpToDate`, `TipOf`, `mergeEdge`, `redoRegister`, `redoFiles`, `intentBytes`, `parseIntent`, `closeOlderPulls` inert, `firstBase` returns its default. Each block gives the command, `date -u` start and end (taken after the step), the direct exit status and the first failing assertions. The first integration red run saw the git calls and the records rule still stubbed (they were written during it); the tests added or changed after it (the status repost, the issue composer, the driver base, the failed pull request read, the merge time, the intent resumed with a pull request closed meanwhile) were run red again on a copy of the tree with the run-package stubs (block "late tests"). Green negatives at the stubs, named: the refusals of `TestAMergeThatTouchesATaskPathOrLeavesTheDefaultIsNotCarried` (its control, a real merge, is red), the refusal count of `TestMergeAndMergeBaseRefuseANonObjectIDBeforeGitStarts` (`calls = 0`), the case "a conflict" of `TestMergeTellsAConflictFromEveryOtherFailure`, the no-request assertions of the failed-read, closed-pull-request and rendered-stop tests, and the second assertion of `TestTheFirstAttemptOfARedoneTaskStartsAtTheRebasedBase`. The wrong-reason first run of the carry test (a harness failure: the base commit not fetched) was fixed in the test before the record below.

### records-U1

- Command: `go test -count=1 -run ^TestTheDetailOfClosedAndRebasedIsATaskID$ ./internal/records/`
- Start: 2026-10-20T18:59:20Z, end: 2026-10-20T18:59:20Z (`date -u`, after the step)
- Exit status (direct, not piped): 1
- First failing assertions (from the saved output):

```text
--- FAIL: TestTheDetailOfClosedAndRebasedIsATaskID (0.00s)
    session_test.go:254: a closed with the detail —: read, want an error in column detail
--
FAIL	github.com/pharzam/layup/internal/records	0.003s
FAIL
```

### git-unit

- Command: `go test -count=1 -run ^(TestMergeAndMergeBaseRefuseANonObjectIDBeforeGitStarts|TestMergeTellsAConflictFromEveryOtherFailure)$ ./internal/git/`
- Start: 2026-10-20T18:59:20Z, end: 2026-10-20T18:59:20Z (`date -u`, after the step)
- Exit status (direct, not piped): 1
- First failing assertions (from the saved output):

```text
--- FAIL: TestMergeAndMergeBaseRefuseANonObjectIDBeforeGitStarts (0.00s)
    git_test.go:621: an empty message: <nil>; want a *FailedError of Code -1
--
--- FAIL: TestMergeTellsAConflictFromEveryOtherFailure (0.00s)
    git_test.go:666: exit 1 with no unmerged entry: sha "", clean false, err <nil>; want clean false, an error true
--
FAIL	github.com/pharzam/layup/internal/git	0.006s
FAIL
```

### git-integration

- Command: `go test -count=1 -tags=integration -run ^TestMergeMakesOneDeterministicMergeCommitOrReportsAConflict$ ./internal/git/`
- Start: 2026-10-20T18:59:20Z, end: 2026-10-20T18:59:20Z (`date -u`, after the step)
- Exit status (direct, not piped): 1
- First failing assertions (from the saved output):

```text
--- FAIL: TestMergeMakesOneDeterministicMergeCommitOrReportsAConflict (0.02s)
    git_integration_test.go:1020: MergeBase = "", <nil>; want 089ff9e9ab70fceaa7cc3a91110b0e7f895370ba
--
FAIL	github.com/pharzam/layup/internal/git	0.017s
FAIL
```

### run-unit

- Command: `go test -count=1 -run ^(TestTheReplacementTakesTheOldRowAndClearsItsPullRequest|TestEachTaskThatNamesTheOldTaskNamesTheNewOne|TestTheRedoIsOneCommitOfTheRegisterTheTwoEventsAndTheBindings|TestTheIntentOfARedoIsReadBackOrRefused|TestTheFirstAttemptOfARedoneTaskStartsAtTheRebasedBase)$ ./internal/run/`
- Start: 2026-10-20T18:59:26Z, end: 2026-10-20T18:59:28Z (`date -u`, after the step)
- Exit status (direct, not piped): 1
- First failing assertions (from the saved output):

```text
--- FAIL: TestTheReplacementTakesTheOldRowAndClearsItsPullRequest (0.00s)
    redo_test.go:30: redoRegister = [], <nil>; want four rows
--- FAIL: TestEachTaskThatNamesTheOldTaskNamesTheNewOne (0.00s)
    redo_test.go:55: redoRegister = [], <nil>; want four rows
--- FAIL: TestTheRedoIsOneCommitOfTheRegisterTheTwoEventsAndTheBindings (0.00s)
    redo_test.go:93: the files = []; want [tasks.tsv tasks/T-new1/events.tsv tasks/T-new1/options/db.tsv tasks/T-old1/events.tsv]
--- FAIL: TestTheIntentOfARedoIsReadBackOrRefused (0.00s)
    redo_test.go:130: parseIntent(intentBytes) = "", "", <nil>; want T-new1 and 2222222222222222222222222222222222222222
--
--- FAIL: TestTheFirstAttemptOfARedoneTaskStartsAtTheRebasedBase (0.00s)
    redo_test.go:145: firstBase(rebased) = 1111111111111111111111111111111111111111; want the base of the event, 2222222222222222222222222222222222222222
--
FAIL	github.com/pharzam/layup/internal/run	0.003s
FAIL
```

### run-integration

- Command: `go test -count=1 -tags=integration -timeout 20m -run ^(TestAnUpToDateVerifiedTaskIsLeftAsItIs|TestADisjointMovedBaseIsMergedAndCarriesTheVerdict|TestAMergeThatTouchesATaskPathOrLeavesTheDefaultIsNotCarried|TestAFailedGitReadIsAnErrorNeverARedo|TestAnOverlappingMovedBaseRedoesTheTask|TestAConflictingMergeOfDisjointPathsRedoesTheTask|TestARestartAfterEachActOfTheRedoCompletesItOnce|TestAPullRequestMergedAfterItsRedoStops|TestAClosedUnmergedPullRequestWithNoRedoIsAWait|TestALostLeaseOpensNoIssueAndWritesNothingAfterIt|TestARedoThatWouldDropTheRenderedCommitIsAStop|TestTheOlderAttemptsPullRequestIsClosedOnce|TestAPushedMergeCommitIsAdoptedOnARestart|TestAMergeAnnouncedAndNotPushedIsMadeAgainWhenTheBaseMovesAgain|TestTheCloseOutHeadIsMergedWhenTheBaseLeavesItsFilesAndRedoneWhenItDoesNot)$ ./internal/run/`
- Start: 2026-10-20T18:59:36Z, end: 2026-10-20T18:59:42Z (`date -u`, after the step)
- Exit status (direct, not piped): 1
- First failing assertions (from the saved output):

```text
--- FAIL: TestAMergeThatTouchesATaskPathOrLeavesTheDefaultIsNotCarried (0.91s)
    baseup_integration_test.go:312: git rev-parse f7eafedee9776c7ae17f240ec954cc955b2fe155^{tree}: exit status 128
--
--- FAIL: TestAnUpToDateVerifiedTaskIsLeftAsItIs (0.93s)
    baseup_integration_test.go:233: BringUpToDate = {Task: Verified:false Head: Wait: Stop: Replaced:}, <nil>; want {T-ab12 verified at 73979491463610f379ba1f6ed8f1cd1790595d1c}
--- FAIL: TestADisjointMovedBaseIsMergedAndCarriesTheVerdict (0.93s)
    baseup_integration_test.go:253: BringUpToDate = {Task: Verified:false Head: Wait: Stop: Replaced:}, <nil>; want a verified merge commit
--- FAIL: TestAPushedMergeCommitIsAdoptedOnARestart (0.94s)
    baseup_integration_test.go:711: BringUpToDate with the bound lost: <nil>; want the stop
--- FAIL: TestAMergeAnnouncedAndNotPushedIsMadeAgainWhenTheBaseMovesAgain (0.99s)
    baseup_integration_test.go:740: BringUpToDate with the push lost: <nil>, the branch at 33f639f69afd42c7db90b2b7fabe27f8aea586e3; want the stop and the branch at 33f639f69afd42c7db90b2b7fabe27f8aea586e3
--- FAIL: TestAConflictingMergeOfDisjointPathsRedoesTheTask (1.00s)
    baseup_integration_test.go:489: BringUpToDate = {Task: Verified:false Head: Wait: Stop: Replaced:}, <nil>; closed events 0; want a redo
```

### run-integration-I4-rerun

- Command: `go test -count=1 -tags=integration -run ^TestAMergeThatTouchesATaskPathOrLeavesTheDefaultIsNotCarried$ ./internal/run/`
- Start: 2026-10-20T19:00:15Z, end: 2026-10-20T19:00:19Z (`date -u`, after the step)
- Exit status (direct, not piped): 1
- First failing assertions (from the saved output):

```text
--- FAIL: TestAMergeThatTouchesATaskPathOrLeavesTheDefaultIsNotCarried (2.00s)
    baseup_integration_test.go:350: the control merge: CarryVerify = "", no session and no rendered commit are named; want it carried
--
FAIL	github.com/pharzam/layup/internal/run	2.008s
FAIL
```

### run-integration-late-tests (stub tree: the stubs of the first red run, the git calls and the records rule real)

- Command: `go test -count=1 -tags=integration -timeout 20m -run ^(TestACarryWhoseStatusFailedIsPostedAtTheRestart|TestTheNewIssueIsComposedFromTheInputsTheRecordsKeep|TestTheDriverStartsTheFirstAttemptOfARedoneTaskAtTheRebasedBase|TestAFailedGitReadIsAnErrorNeverARedo|TestAMergeThatTouchesATaskPathOrLeavesTheDefaultIsNotCarried|TestARestartAfterEachActOfTheRedoCompletesItOnce|TestADisjointMovedBaseIsMergedAndCarriesTheVerdict)$ ./internal/run/`
- Start: 2026-10-20T19:05:38Z, end: 2026-10-20T19:05:42Z (`date -u`, after the step)
- Exit status (direct, not piped): 1
- First failing assertions (from the saved output):

```text
--- FAIL: TestTheDriverStartsTheFirstAttemptOfARedoneTaskAtTheRebasedBase (0.44s)
    baseup_integration_test.go:942: the base of the first attempt is "fe396aa7f0e9a97e3088cb553642ad89a6274ee7"; want the base of the event rebased, cccccccccccccccccccccccccccccccccccccccc
--- FAIL: TestADisjointMovedBaseIsMergedAndCarriesTheVerdict (0.55s)
    baseup_integration_test.go:299: BringUpToDate = {Task: Verified:false Head: Wait: Stop: Replaced:}, <nil>; want a verified merge commit
--- FAIL: TestACarryWhoseStatusFailedIsPostedAtTheRestart (0.58s)
    baseup_integration_test.go:862: BringUpToDate with the status refused = {Task: Verified:false Head: Wait: Stop: Replaced:}, nil; want an error
--- FAIL: TestTheNewIssueIsComposedFromTheInputsTheRecordsKeep (0.62s)
    baseup_integration_test.go:896: an architecture task of a decision: {Task: Verified:false Head: Wait: Stop: Replaced:}, <nil>; want a redo
--- FAIL: TestARestartAfterEachActOfTheRedoCompletesItOnce (0.68s)
    baseup_integration_test.go:569: stop intent: <nil>; want the stop
--- FAIL: TestAMergeThatTouchesATaskPathOrLeavesTheDefaultIsNotCarried (1.54s)
    baseup_integration_test.go:404: the control merge: CarryVerify = "", no session and no rendered commit are named; want it carried
```

