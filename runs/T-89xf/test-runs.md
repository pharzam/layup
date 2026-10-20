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

The stubs: `git.Merge`/`git.MergeBase` return zero values, `FetchBranch`, `BringUpToDate`, `UpToDate`, `TipOf`, `mergeEdge`, `redoRegister`, `redoFiles`, `intentBytes`, `parseIntent`, `closeOlderPulls` inert, `firstBase` returns its default. Each block gives the command, `date -u` start and end (taken after the step), the direct exit status and the first failing assertions. The first integration red run saw the git calls and the records rule still stubbed (they were written during it); the tests added or changed after it (the status repost, the issue composer, the driver base, the failed pull request read, the merge time, the intent resumed with a pull request closed meanwhile) were run red again on a copy of the tree with the run-package stubs (block "late tests"). Green negatives at the stubs, named: the refusals of `TestAMergeThatTouchesATaskPathOrLeavesTheDefaultIsNotCarried` (its control, a real merge, is red), the refusal count of `TestMergeAndMergeBaseRefuseANonObjectIDBeforeGitStarts` (`calls = 0`), the case "a conflict" of `TestMergeTellsAConflictFromEveryOtherFailure`, the no-request assertions of the failed-read, closed-pull-request and rendered-stop tests, and the second assertion of `TestTheFirstAttemptOfARedoneTaskStartsAtTheRebasedBase`. The unit tests of the replacement were first in `internal/run/redo_test.go` (the red lines below name it) and moved to `internal/run/repair_test.go` to keep the file count in the budget. The wrong-reason first run of the carry test (a harness failure: the base commit not fetched) was fixed in the test before the record below.

### records-U1

- Command: `go test -count=1 -run '^TestTheDetailOfClosedAndRebasedIsATaskID$' ./internal/records/`
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

- Command: `go test -count=1 -run '^(TestMergeAndMergeBaseRefuseANonObjectIDBeforeGitStarts|TestMergeTellsAConflictFromEveryOtherFailure)$' ./internal/git/`
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

- Command: `go test -count=1 -tags=integration -run '^TestMergeMakesOneDeterministicMergeCommitOrReportsAConflict$' ./internal/git/`
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

- Command: `go test -count=1 -run '^(TestTheReplacementTakesTheOldRowAndClearsItsPullRequest|TestEachTaskThatNamesTheOldTaskNamesTheNewOne|TestTheRedoIsOneCommitOfTheRegisterTheTwoEventsAndTheBindings|TestTheIntentOfARedoIsReadBackOrRefused|TestTheFirstAttemptOfARedoneTaskStartsAtTheRebasedBase)$' ./internal/run/`
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

- Command: `go test -count=1 -tags=integration -timeout 20m -run '^(TestAnUpToDateVerifiedTaskIsLeftAsItIs|TestADisjointMovedBaseIsMergedAndCarriesTheVerdict|TestAMergeThatTouchesATaskPathOrLeavesTheDefaultIsNotCarried|TestAFailedGitReadIsAnErrorNeverARedo|TestAnOverlappingMovedBaseRedoesTheTask|TestAConflictingMergeOfDisjointPathsRedoesTheTask|TestARestartAfterEachActOfTheRedoCompletesItOnce|TestAPullRequestMergedAfterItsRedoStops|TestAClosedUnmergedPullRequestWithNoRedoIsAWait|TestALostLeaseOpensNoIssueAndWritesNothingAfterIt|TestARedoThatWouldDropTheRenderedCommitIsAStop|TestTheOlderAttemptsPullRequestIsClosedOnce|TestAPushedMergeCommitIsAdoptedOnARestart|TestAMergeAnnouncedAndNotPushedIsMadeAgainWhenTheBaseMovesAgain|TestTheCloseOutHeadIsMergedWhenTheBaseLeavesItsFilesAndRedoneWhenItDoesNot)$' ./internal/run/`
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

- Command: `go test -count=1 -tags=integration -run '^TestAMergeThatTouchesATaskPathOrLeavesTheDefaultIsNotCarried$' ./internal/run/`
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

- Command: `go test -count=1 -tags=integration -timeout 20m -run '^(TestACarryWhoseStatusFailedIsPostedAtTheRestart|TestTheNewIssueIsComposedFromTheInputsTheRecordsKeep|TestTheDriverStartsTheFirstAttemptOfARedoneTaskAtTheRebasedBase|TestAFailedGitReadIsAnErrorNeverARedo|TestAMergeThatTouchesATaskPathOrLeavesTheDefaultIsNotCarried|TestARestartAfterEachActOfTheRedoCompletesItOnce|TestADisjointMovedBaseIsMergedAndCarriesTheVerdict)$' ./internal/run/`
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

## Gates on the head b07115bb3ad18fda9cde5e4f96527692aca9dd54

- `go build ./...`, `go vet ./...`, `go vet -tags=integration ./...`, `go vet -tags=e2e ./...`, `go test -count=1 ./...`: exit 0 each.
- `go test -count=1 -tags=e2e ./...`: exit 0 (run on 99b6059).
- The full integration level, `TMPDIR=/home/layup/.cache/layup-session/tmp-green go test -count=1 -tags=integration -timeout 30m ./...`, on 99b6059 (the later commits move the unit tests of the replacement to `internal/run/repair_test.go` and add records): exit 0, 21 packages ok, `internal/run` 374 s (output in `/home/layup/.cache/layup-session/T-89xf-full-integration.txt`). main had not moved (`df9bfac`), so no merge head.
- adr-lint, prd-lint, link-lint, run-discipline-tests, nested-checkout-check, setup-check: exit 0 each. `git diff --check origin/main HEAD`: exit 0.
- `sh docs/tests/freeze-check.sh --task T-89xf --base origin/main --budget 2800/30` (output at `/home/layup/.cache/layup-session/T-89xf-freeze.txt`), exit 0:

```text
ok K1: white space
ok K2: registration of new tests
ok K3: test selectors
ok K4: restore evidence of mutants
ok K5: projection: figure 2474/30, maximum 2800/30, margin 326/0 (the figure at landing is still measured at landing)
ok K6: conditions table
```

## The up-to-date helper, for the fold with K4 (T-ehb4, pushed at 08c5ec6, `upToDate` in `internal/run/finish.go`)

K3's helper is `(*Sessions).UpToDate(ctx, head) (current bool, def string, err error)` in `internal/run/baseup.go`. It fetches the default branch with `SessionStore.FetchBranch` (so it writes the ref `refs/heads/<default>` of the run's clone) and gives whether that head is an ancestor of `head`, with the head it fetched. `head` is the tip of the open attempt, `TipOf(events, attempt)`, the last event `bound` of any session: the verified head, the rendered commit, a merge of this task, or a close-out head. It does not look for the head that carries the verdict; `BringUpToDate` and `carries` decide that. The helper it replaces or is replaced by must therefore be given the carried head, as K4's D1 does. Whichever of K3 and K4 lands second merges main and folds the two into one home (condition 4 of the plan review of #425).

## Round 1 of #424 (the review a7c9d9f): red transcripts, before the code of the fixes

### round 1, findings 1 to 3: red against the code of a7c9d9f (merged with main d27b331)

- Command: `go test -count=1 -tags=integration -run '^(TestAnInterruptedCarryIsAdoptedWhenTheDefaultMovesAgain|TestOpenPullRequestAdoptsTheSamePullRequestAfterACarry|TestASecondRedoOfAShapeTaskKeepsTheOriginalPin)$' ./internal/run/`
- Start: 2026-10-20T19:52:28Z, end: 2026-10-20T19:52:29Z (`date -u`, after the step)
- Exit status (direct, not piped): 1
- First failing assertions (from the saved output):

```text
--- FAIL: TestASecondRedoOfAShapeTaskKeepsTheOriginalPin (0.82s)
    baseup_integration_test.go:1137: the second replacement issue = "LAYUP task T-zzzz: the target's architecture.\n\nRequirement IDs: REQ-005\nPinned draft: spec/draft.tsv, SHA-256 3e54755b5e4a5c19315a0a7213ac667aff194600ffaf84bcef7a5dd62d53c10a, read at the records commit 69fbce353d92933515bad196652b040517aa24e1.\n\nDefinition of Done: docs/spec/ and docs/adr/ hold each decision that has one survivor as decided, and a section \"Decision pending the bet\" for each decision that has more than one; docs/tasks/T-zzzz.md is the task's file. The head passes the four layup/ statuses, and a verifier on a harness outside the authors finds nothing material in scope.\n", <nil>; want the original pin "SHA-256 534bfca4b311333b887811de908c81e41fd1be78692bc32d3244812b748e2239, read at the records commit c0413aa959eea78ce766903e84b4d5d4e9dad6a8"
--- FAIL: TestOpenPullRequestAdoptsTheSamePullRequestAfterACarry (0.89s)
    baseup_integration_test.go:1075: OpenPullRequest 1 = {Number:0 NodeID: State: Draft:false Merged:false HeadSHA: MergeSHA: Body:}, the head of the pull request 1 of task/T-ab12/1 moved: f533a00226d812e46ec6e44dbe2f3dc03286777e, not the bound 485b63a9e7f3ecbb3a1081a8259ade516307c0ed; want pull request 1 at f533a00226d812e46ec6e44dbe2f3dc03286777e
--- FAIL: TestAnInterruptedCarryIsAdoptedWhenTheDefaultMovesAgain (0.91s)
    baseup_integration_test.go:1036: restart 1 = {Task: Verified:false Head: Wait: Stop: Replaced:}, the push of the session's head was refused: task/T-ab12/1 holds "c17609d55c19ee0b48a7fefaf86b9b4f3f7dba4d", not the head 397b28d55af35fa50129cb40fcd4f9e2ac02b64b or its merge 72d30dda0bdf803b76e411b9b5eba1db41a13abc; want a merge of D2 over the adopted c17609d55c19ee0b48a7fefaf86b9b4f3f7dba4d
--
FAIL	github.com/pharzam/layup/internal/run	0.917s
FAIL
```

### round 1, finding 4: the tests that characterize the code of a7c9d9f (green there; their mutants are in mutants.md)

- Command: `go test -count=1 -tags=integration -run '^(TestTheCarryGoesOverTheRenderedFirstParent|TestARestartAfterEachActOfTheRedoCompletesItOnce)$' ./internal/run/`
- Start: 2026-10-20T19:52:29Z, end: 2026-10-20T19:52:33Z (`date -u`, after the step)
- Exit status (direct, not piped): 0
- First failing assertions (from the saved output):

```text
```

## Gates on the code head of round 1 (late evidence)

- Code head: the commit `fix: T-89xf round 1`; main `d27b331` merged before it. `go build ./...`, `go vet ./...` and with `-tags=integration` and `-tags=e2e`, `gofmt -l .`, `go test -count=1 ./...`: exit 0 each.
- Full integration level, `TMPDIR=/home/layup/.cache/layup-session/tmp-green go test -count=1 -tags=integration -timeout 30m ./...`, on the code head: exit 0, 21 packages ok, `internal/run` 308.6 s (`/home/layup/.cache/layup-session/T-89xf-full-integration.txt`). `go test -count=1 -tags=e2e ./...`: exit 0.
- The six lints and `git diff --check origin/main HEAD`: exit 0.
- Size: excluding `runs/` (ADR-0028) the diff against `origin/main` is 2527 lines over 26 files; with the close-out projection of 70/2 the figure is inside 2800/30. The freeze check of this tree still counts `runs/` in K5 (its figure 2872/31, margin -72/-1), so K5 fails on the tool's old count; K1 to K4 and K6 pass.
- `sh docs/tests/freeze-check.sh --task T-89xf --base origin/main --budget 2800/30` (output at `/home/layup/.cache/layup-session/T-89xf-freeze.txt`), exit 1 for K5 alone:

```text
ok K1: white space
ok K2: registration of new tests
ok K3: test selectors
ok K4: restore evidence of mutants
FAIL K5: -:0: projection: figure 2872/31, maximum 2800/30, margin -72/-1 (the figure at landing is still measured at landing)
ok K6: conditions table
```

## Round 2 of #424 (the review 3f8ef98): red transcripts, before the code of the fixes

The findings F5 to F8 of the review; F8 (the three requirement rows of the PRD) is a record and has no test. The eight tests below are run against the code of 3f8ef98 with inert stubs for the pin (`Pin` field of the plan, `pinPath`, `pinBytes`, `parsePin`); all eight fail. The full output is `/home/layup/.cache/layup-session/k3/r4.txt`.

### round 2, findings F5 to F7: red against the code of 3f8ef98 (inert stubs: the plan field Pin, pinPath, pinBytes, parsePin)

- Command: `go test -count=1 -tags=integration -run '^(TestAnUnexplainedHeadOfTheBranchRefusesTheRecoveryBeforeAnyWrite|TestOpenPullRequestAdoptsTheSamePullRequestAfterARenderedCarry|TestAShapeTaskIsOpenedWithItsPinInTheRecords|TestTheReplacementOfAShapeTaskInheritsTheOpeningPinExactly|TestAShapeTaskWithNoPinIsNotRedone|TestThePinIsInheritedThroughAnyChainOfRedos|TestThePinIsReadBackOrRefused|TestASecondRedoOfAShapeTaskKeepsTheOriginalPin)$' ./internal/run/`
- Start: 2026-10-20T20:57:38Z, end: 2026-10-20T20:57:41Z (`date -u`, after the step)
- Exit status (direct, not piped): 1
- First failing assertions (from the saved output):

```text
--- FAIL: TestThePinIsInheritedThroughAnyChainOfRedos (0.00s)
    repair_test.go:461: redo 1: tasks/T-1000/pin.txt = ""; want the pin as it was
--- FAIL: TestThePinIsReadBackOrRefused (0.00s)
    repair_test.go:476: parsePin(pinBytes) = "", "", <nil>
--
--- FAIL: TestAShapeTaskIsOpenedWithItsPinInTheRecords (0.35s)
    handoffpost_integration_test.go:124: git show layup-records:tasks/T-uyst/pin.txt: exit status 128
--
--- FAIL: TestTheReplacementOfAShapeTaskInheritsTheOpeningPinExactly (0.72s)
    handoffpost_integration_test.go:124: git show layup-records:tasks/T-krgg/pin.txt: exit status 128
--
--- FAIL: TestASecondRedoOfAShapeTaskKeepsTheOriginalPin (0.73s)
    baseup_integration_test.go:1126: the first replacement issue lacks the original pin "SHA-256 534bfca4b311333b887811de908c81e41fd1be78692bc32d3244812b748e2239, read at the records commit 3b93b4f3d25676d72fef33a4b9bf8bc30c81f405":
--
```

## Gates on the code head of round 2 (late evidence)

- Code head `2a8a433` (the fix commit; the later commits are the mutants record). Main had not moved (`d27b331`). `go build ./...`, `go vet ./...` and with `-tags=integration` and `-tags=e2e`, `gofmt -l .`, `go test -count=1 ./...`: exit 0 each.
- Full integration level, `TMPDIR=/home/layup/.cache/layup-session/tmp-green go test -count=1 -tags=integration -timeout 30m ./...`: exit 0, 21 packages ok, `internal/run` 188.5 s (`/home/layup/.cache/layup-session/T-89xf-full-integration.txt`). `go test -count=1 -tags=e2e ./...`: exit 0.
- The six lints and `git diff --check origin/main HEAD`: exit 0.
- Size excluding `runs/T-89xf` (ADR-0028): 2715 lines over 27 files against `origin/main`; with the close-out projection of 70/2 it is 2785/29. K5 counts `runs/` (figure 3111/32) and fails as expected.
- `sh docs/tests/freeze-check.sh --task T-89xf --base origin/main --budget 2800/30` (`/home/layup/.cache/layup-session/T-89xf-freeze.txt`), exit 1 for K5 alone:

```text
ok K1: white space
ok K2: registration of new tests
ok K3: test selectors
ok K4: restore evidence of mutants
FAIL K5: -:0: projection: figure 3111/32, maximum 2800/30, margin -311/-2 (the figure at landing is still measured at landing)
ok K6: conditions table
```

## Round 3 of #424 (the review b9ed9b0): red transcripts, before the code of the fixes

F9 and F10, run against the code of b9ed9b0 (the test of F9 holds the heartbeat between the capture of the pin and the open of the issue by a wrapper of the token of the records store). The full output is `/home/layup/.cache/layup-session/k3/r5.txt`; the F10 cases of `TestThePinIsReadBackOrRefused` and `TestAMalformedOpeningPinRefusesTheRedoBeforeAnyWrite` fail for each malformed form. Evidence added after the code is labelled late.

### round 3, F9 and F10: red against the code of b9ed9b0

- Command: `go test -count=1 -tags=integration -run '^(TestTheOpeningPinIsCapturedOnceWhateverTheRecordsDoBeforeTheIssue|TestAMalformedOpeningPinRefusesTheRedoBeforeAnyWrite|TestThePinIsReadBackOrRefused)$' ./internal/run/`
- Start: 2026-10-20T21:32:50Z, end: 2026-10-20T21:32:59Z (`date -u`, after the step)
- Exit status (direct, not piped): 1
- First failing assertions (from the saved output):

```text
--- FAIL: TestThePinIsReadBackOrRefused (0.00s)
    repair_test.go:482: parsePin("not-a-digest not-a-commit\n") = "not-a-digest", "not-a-commit", nil; want an error
--
--- FAIL: TestTheOpeningPinIsCapturedOnceWhateverTheRecordsDoBeforeTheIssue (0.34s)
    baseup_integration_test.go:1389: tasks/T-33ej/pin.txt = "b95428303e68cc271142ed228363ab48fe4ee76bd9fa37abf64ce5d6f52a5e33 8a1bd83a29448c644802492cc52f78958dc23158\n"; want the pin of the issue, "b95428303e68cc271142ed228363ab48fe4ee76bd9fa37abf64ce5d6f52a5e33 2cc589335de624994cc4b741cf4cd1c5fa060e3b\n"
--- FAIL: TestAMalformedOpeningPinRefusesTheRedoBeforeAnyWrite (7.25s)
    baseup_integration_test.go:1406: the pin "not-a-digest not-a-commit\n": {Task:T-ab12 Verified:false Head: Wait: Stop: Replaced:T-gpkd}, <nil>; writes [PATCH /repos/acme/target/issues/3 POST /repos/acme/target/issues]; want an error and no write
--
FAIL	github.com/pharzam/layup/internal/run	7.256s
FAIL
```
