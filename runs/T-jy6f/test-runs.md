# T-jy6f - the test runs of row 66r

## S0 - the reconciliation of the plan (6250184301) with the landed code

Read at `origin/main` = `5a290c9` (the branch had no work, so the merge of `origin/main` was a fast-forward). Each name of the plan
was checked in the tree, not in a plan or a review (note 4 of 6250262119).

| Plan name | Landed form | Note |
| --------- | ----------- | ---- |
| 66f: `PostStatuses`, `HeadInput`, `PostStatus`, `Status`, `VerifyAt`, `VerifyInput`, `GatesTable`, `Sessions.Bot` | `internal/run/statuses.go:202`, `:176`, `:63`, `:33`; `verifystatus.go:33`, `:21`; `gates.go:167`; `session.go:65` | `PostStatuses` posts the four contexts, `layup/verify` included (a `not-active` failure too, `verifystatus.go:44-46`) |
| 66f: `holdsLease` | `handoffpost.go:418` | a method of `Sessions`; a `state` reaches it through `state.sessions()` |
| 66g: `OpenPullRequest`, `PullInput`, `Sessions.Web` | `pullrequest.go:104`, `:24`; `session.go:66` | the title is `Type` and `Description` (`checkTitle`), no task ID |
| 66q: the body edit | `(*Sessions).LinkReviewRecord(ctx, task, verifier) (bool, error)`, `reviewlink.go:30` | the plan had no name; it reads the latest `valid comment` handoff of the verifier |
| 66s: `PlanFile`, `PlanReviewFile`, `ReviewRecordFile`, `PromptIn.Options` | `instructions.go:22-24`; `promptfile.go:41-49` | landed |
| 66s: `PromptContext.base` (and `task`, `step`) | `promptfile.go:51-64`, set at `:176` | landed, so the guard needs no new field; `readOnlyStep` is not in `internal/session` (the plan cited `context.go:59-63`, which is the step table): it is to be decided from the private `PromptContext.step` |
| 66s: a reusable per-kind Definition of Done | none | 66s landed no exported text of it; the `shape-issue` text carries its own (two blocks, one for each kind) |
| 66v: `RepairTask` | `repair.go:30` | the driver calls it first |
| the base guard | none | `loop.md:469` leaves it to this row; `TaskSession` (`session.go:208`) has none |
| `probeTaskID` | `probe.go:162`, reads `sessions.tsv` only, `crypto/rand` direct | the seam `Sessions.random` and `drawTaskID` are to be made |
| `registerTask`, `registerFiles`, `parameterFiles`, `loopPermissionsStep` | `register.go:267`, `:213`, `:254`; `parameters.go` | as planned; `state.store` is the concrete `*recordsStore`, so the opening tests run in the real-git world (condition 5) |
| the pinned linter call | `spec.WaitingSpecAccepted(script, prdDir)`, `internal/spec/prdlint.go:61` | for the scope addition; the script reads the facts from `<dirname prdDir>/facts` |

### A landed form that the plan does not settle (STOP)

The verifier session of the plan starts at the bound head (`TaskSession{Base: bound}`, plan §2.5). `session.Make` clones **the
default branch only** (`git.CloneLocal`, `--single-branch`, `internal/git/git.go:507`; `session.md:52-58`) and then
`SwitchCreate(task/<task>/<attempt>, base)`. A developer's head is in the run's clone only as `refs/layup/sessions/<id>` (and, after
a restart, as `origin/task/<task>/<attempt>`), so it is not in the history of the default branch, and the start of the verifier
fails. The existing test world moves the default branch onto the head by hand (`reviewlink_integration_test.go:61`, "as the loop would
before a verifier"); no production code does. The close-out of row 66k's children needs the same.

## S3 - the red runs, S4-S6 - the green runs

Every command was run directly in a worktree, never piped for its status: the exit status is the one `go` gave (`exit=`), and the
start and end are `date -u` read before and after the run. The long lines of a failure are cut at 330 characters; nothing else is
changed.

**The stubs.** The red runs are against a worktree of `3f06e20` (the inert declarations of commit `b538f03`: `finalIDs` returning
nothing, `openShapeTask` returning `Task{}, nil`, `loopPosition` returning the zero step, `readOnlyStep` false, `UpdateRef` returning nil,
no guard) plus the two inert declarations that the final tests need to compile: `devStep` and `shapeIssueBlocks` (empty map). The final
test files were copied over it. An earlier red run, on the first form of the tests (before four corrections of the test world: the lease
that the world holds, a harness function that returned the status of its last test, a model context of 1000 tokens, and the layout of
the architecture head), failed for the same reasons; this one is the record.

Green negatives at the stub, each beside its red control: `TestTheShapeIssueBlocksOfLoopMd` was a vacuous green at the first stub
(the map of blocks was empty, so its loop ran nothing); it now asserts that the code holds both kinds, and is red. The
`TestADrawnIDIsDrawnAgainWhileASessionHoldsIt` of commit `3f06e20` is the green characterization U0 and stays green.

### Red, the unit level

Command: `go test -count=1 ./internal/run/ ./internal/git/`

```
start 2026-10-19T10:52:05Z
--- FAIL: TestADrawnIDIsDrawnAgainWhileTheRegisterHoldsIt (0.00s)
    probe_test.go:118: taskID = "T-0000", <nil>; want T-1111, the second draw
--- FAIL: TestAReadOnlyStepHasNoPush (0.00s)
    session_test.go:396: readOnlyStep(plan-of-a-task) = false; want true
    session_test.go:396: readOnlyStep(plan-review) = false; want true
    session_test.go:396: readOnlyStep(verification) = false; want true
--- FAIL: TestATaskSessionRefusesAPromptFileOfAnotherTaskOrBase (0.00s)
    session_test.go:424: another task: mkdir : no such file or directory; want "the prompt file is of the task T-cd34, the session's task is T-ab12"
    session_test.go:427: another task: the calls ["sweep" "make" "remove /h/sessions/S-1a2b3c4d"] and 0 commits; want nothing swept, made or written
    session_test.go:424: another base: mkdir : no such file or directory; want "the prompt file of T-ab12 was built at base cccccccccccccccccccccccccccccccccccccccc, the session's base is 0123456789abcdef0123456789abcdef01234567"
    session_test.go:427: another base: the calls ["sweep" "make" "remove /h/sessions/S-1a2b3c4d"] and 0 commits; want nothing swept, made or written
--- FAIL: TestThePositionIsTheFirstStepTheRecordsShowNotDone (0.00s)
    shapeloop_test.go:94: specification, 0 facts: position 0 (""); want 1
    shapeloop_test.go:94: specification, 1 facts: position 0 (""); want 2
    shapeloop_test.go:94: specification, 2 facts: position 0 (""); want 2
    shapeloop_test.go:94: specification, 3 facts: position 0 (""); want 3
    shapeloop_test.go:94: specification, 4 facts: position 0 (""); want 4
    shapeloop_test.go:94: specification, 5 facts: position 0 (""); want 4
    shapeloop_test.go:94: specification, 6 facts: position 0 (""); want 5
    shapeloop_test.go:94: specification, 7 facts: position 0 (""); want 6
    shapeloop_test.go:94: specification, 8 facts: position 0 (""); want 6
    shapeloop_test.go:94: specification, 9 facts: position 0 (""); want 13
    shapeloop_test.go:94: specification, 10 facts: position 0 (""); want 13
    shapeloop_test.go:94: specification, 11 facts: position 0 (""); want 7
    shapeloop_test.go:94: specification, 12 facts: position 0 (""); want 8
    shapeloop_test.go:94: specification, 13 facts: position 0 (""); want 8
    shapeloop_test.go:94: specification, 14 facts: position 0 (""); want 9
    shapeloop_test.go:94: specification, 15 facts: position 0 (""); want 10
    shapeloop_test.go:94: architecture, 0 facts: position 0 (""); want 1
    shapeloop_test.go:94: architecture, 1 facts: position 0 (""); want 2
    shapeloop_test.go:94: architecture, 2 facts: position 0 (""); want 2
    shapeloop_test.go:94: architecture, 3 facts: position 0 (""); want 3
    shapeloop_test.go:94: architecture, 4 facts: position 0 (""); want 4
    shapeloop_test.go:94: architecture, 5 facts: position 0 (""); want 4
    shapeloop_test.go:94: architecture, 6 facts: position 0 (""); want 5
    shapeloop_test.go:94: architecture, 7 facts: position 0 (""); want 6
    shapeloop_test.go:94: architecture, 8 facts: position 0 (""); want 6
    shapeloop_test.go:94: architecture, 9 facts: position 0 (""); want 13
    shapeloop_test.go:94: architecture, 10 facts: position 0 (""); want 13
    shapeloop_test.go:94: architecture, 11 facts: position 0 (""); want 7
    shapeloop_test.go:94: architecture, 12 facts: position 0 (""); want 8
    shapeloop_test.go:94: architecture, 13 facts: position 0 (""); want 8
    shapeloop_test.go:94: architecture, 14 facts: position 0 (""); want 9
    shapeloop_test.go:94: architecture, 15 facts: position 0 (""); want 10
    shapeloop_test.go:100: the position to the verifier = {Pos:0 Session: Msg: Attempt:0 Base: Head:}; want attempt 1, base, head and the developer's session
    shapeloop_test.go:103: the position to finish the plan = {Pos:0 Session: Msg: Attempt:0 Base: Head:}; want the plan session
--- FAIL: TestAStageThatEndedIsFinishedAndAStageThatStoppedIsRunAgain (0.00s)
    shapeloop_test.go:133: the plan session with no result: 0 ""; want 2 holding ""
    shapeloop_test.go:133: the verifier ended and posted nothing yet: 0 ""; want 9 holding ""
    shapeloop_test.go:133: a plan session that ended wall: 0 ""; want 13 holding "ended wall"
    shapeloop_test.go:133: a plan whose status is needs_context: 0 ""; want 13 holding "ended needs_context"
    shapeloop_test.go:133: a refused session: 0 ""; want 13 holding "was refused (artifact)"
    shapeloop_test.go:133: an invalid plan handoff: 0 ""; want 13 holding "task plan → plan review is invalid plan-fields"
    shapeloop_test.go:133: a plan review that rejects: 0 ""; want 13 holding "plan review → developer is invalid verdict-field"
    shapeloop_test.go:133: an invalid developer handoff: 0 ""; want 13 holding "is invalid spec-sections"
    shapeloop_test.go:133: a developer that ended with no commit pushed: 0 ""; want 13 holding "ended with no push"
    shapeloop_test.go:133: a push with no result: 0 ""; want 13 holding "a push with no result"
    shapeloop_test.go:133: a refused push: 0 ""; want 13 holding "was refused (push-refused)"
    shapeloop_test.go:133: a start that was refused: 0 ""; want 13 holding "the start of the session S-x1 was refused (context 1175 1000)"
    shapeloop_test.go:133: a start refused in an attempt before: 0 ""; want 2 holding ""
    shapeloop_test.go:133: an invalid verifier handoff: 0 ""; want 13 holding "the handoff verifier is invalid nothing-material"
--- FAIL: TestThePlanReviewIsBeforeTheFirstPushAndTheVerifierAfterItsBound (0.00s)
    shapeloop_test.go:143: a qa session before the first push: {Pos:0 Session: Msg: Attempt:0 Base: Head:}; want the plan review to finish
    shapeloop_test.go:148: an unclassified qa session: {Pos:0 Session: Msg: Attempt:0 Base: Head:}; want the stop
    shapeloop_test.go:153: a qa session before the plan: {Pos:0 Session: Msg: Attempt:0 Base: Head:}; want the stop
--- FAIL: TestOneMaterialRoundStartsAnAttemptAndTheSecondWaitsAtTheCap (0.00s)
    shapeloop_test.go:162: the first material round: {Pos:0 Session: Msg: Attempt:0 Base: Head:}; want the next attempt at the base, after the verdict of S-v1
    shapeloop_test.go:167: after the new attempt: {Pos:0 Session: Msg: Attempt:0 Base: Head:}; want the developer of attempt 2
--- FAIL: TestNotMergeableFindingsRecordedIsAWaitNeverAnAttempt (0.00s)
    shapeloop_test.go:187: not mergeable: {Pos:0 Session: Msg: Attempt:0 Base: Head:}; want the wait
--- FAIL: TestAVerifiedHeadNeedsNothingMaterialAtTheLastBoundHead (0.00s)
    shapeloop_test.go:194: the control: {Pos:0 Session: Msg: Attempt:0 Base: Head:}; want the verified head
--- FAIL: TestTheFindingsOfTheNextAttemptAreTheLastRoundsOrTheGatesTable (0.00s)
    shapeloop_test.go:206: findingsOf after a material round = [], ""; want its finding and no gate head
    shapeloop_test.go:212: findingsOf after a gate failure = [], ""; want no record and the bound head bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb
    shapeloop_test.go:219: findingsOf for attempt 3 = [], ""; want no finding of the older record and the gate head cccccccccccccccccccccccccccccccccccccccc
--- FAIL: TestTheFinalIDsAreTheKindsOfThePinnedDraft (0.00s)
    shapetask_test.go:16: finalIDs = []; want [REQ-001 NFR-002 REQ-004]
    shapetask_test.go:23: finalID(R-007, requirement) = "", false; want "REQ-007", true
    shapetask_test.go:23: finalID(R-012, non-functional) = "", false; want "NFR-012", true
--- FAIL: TestTheIssueOfAShapeTaskHasItsTitleAndTheFilledBlockOfItsKind (0.00s)
    shapetask_test.go:30: the title of a specification = ""
    shapetask_test.go:33: the title of an architecture task = ""
    shapetask_test.go:40: specification: the text lacks "LAYUP task T-ab12: ":
    shapetask_test.go:40: specification: the text lacks "Requirement IDs: REQ-001 NFR-002\n":
    shapetask_test.go:40: specification: the text lacks "SHA-256 eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee, read at the records commit f00d.\n":
    shapetask_test.go:40: specification: the text lacks "docs/tasks/T-ab12.md":
    shapetask_test.go:40: architecture: the text lacks "LAYUP task T-ab12: ":
    shapetask_test.go:40: architecture: the text lacks "Requirement IDs: REQ-001 NFR-002\n":
    shapetask_test.go:40: architecture: the text lacks "SHA-256 eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee, read at the records commit f00d.\n":
    shapetask_test.go:40: architecture: the text lacks "docs/tasks/T-ab12.md":
FAIL
FAIL	github.com/pharzam/layup/internal/run	0.398s
--- FAIL: TestUpdateRefRefusesItsInputBeforeGitStarts (0.00s)
    git_test.go:597: a ref with a colon: <nil>, 0 starts of git; want a *FailedError of code -1 and none
    git_test.go:597: a ref that ends in .lock: <nil>, 0 starts of git; want a *FailedError of code -1 and none
    git_test.go:597: an empty SHA: <nil>, 0 starts of git; want a *FailedError of code -1 and none
    git_test.go:597: a tag: <nil>, 0 starts of git; want a *FailedError of code -1 and none
    git_test.go:597: a ref with a space: <nil>, 0 starts of git; want a *FailedError of code -1 and none
    git_test.go:597: a ref that ends in a slash: <nil>, 0 starts of git; want a *FailedError of code -1 and none
    git_test.go:597: an empty name below the prefix: <nil>, 0 starts of git; want a *FailedError of code -1 and none
    git_test.go:597: a SHA that is a branch name: <nil>, 0 starts of git; want a *FailedError of code -1 and none
    git_test.go:597: a ref of the default branch: <nil>, 0 starts of git; want a *FailedError of code -1 and none
    git_test.go:597: a ref outside heads: <nil>, 0 starts of git; want a *FailedError of code -1 and none
    git_test.go:597: a ref that leaves the prefix: <nil>, 0 starts of git; want a *FailedError of code -1 and none
    git_test.go:604: 0 starts of git, want 1
FAIL
FAIL	github.com/pharzam/layup/internal/git	0.025s
FAIL
exit=1
end 2026-10-19T10:52:17Z
```

### Red, the integration level

Command: `TMPDIR=/home/layup/.cache/layup-session/tmp-green go test -count=1 -tags=integration -timeout 30m -run '^(TestAShapeTaskIsOpenedOnlyWithTheLoopPermissions|TestARefusedOpeningWritesNoIssue|TestAnOpenedShapeTaskIsReusedAndOpensNoSecondIssue|TestAStopBetweenTheIssueAndTheRowOpensASecondIssue|TestTheShapeIssueBlocksOfLoopMd|TestTheShapeTasksGoThroughTheLoopToTheirVerifiedHeads|TestAGateFailureOfAShapeTaskStartsANewAttempt|TestAMaterialRoundStartsOneAttemptAndTheSecondWaitsAtTheCap|TestNotMergeableStartsNoAttempt|TestAVerifierWithNoIndependentHarnessPostsNotActiveBeforeTheStop|TestAnInterruptedStageIsFinishedFromItsKeptFilesAndRunsNoSecondSession|TestAMissingSectionMakesTheHandoffInvalidAndStops|TestASpecificationThatPrdLintDoesNotAcceptStopsBeforeTheVerifier|TestADeveloperThatCommitsNothingMakesAnInvalidHandoffAndStops|TestARestartAfterTheRowStartsTheFirstAttemptOnce|TestTheCapOfMilestoneZeroIsAskedBeforeEachSession|TestUpdateRefSetsABranchOfLayupsNamespaceAndNoOther)$' ./internal/run/ ./internal/git/`

```
start 2026-10-19T10:52:17Z
--- FAIL: TestTheShapeTasksGoThroughTheLoopToTheirVerifiedHeads (2.11s)
    shapeloop_integration_test.go:85: AdvanceShapeTask() = {Task: Verified:false Wait: Stop:}, <nil>; want the verified head
--- FAIL: TestAGateFailureOfAShapeTaskStartsANewAttempt (1.87s)
    shapeloop_integration_test.go:137: AdvanceShapeTask() = {Task: Verified:false Wait: Stop:}, <nil>; want the verified head
--- FAIL: TestAMaterialRoundStartsOneAttemptAndTheSecondWaitsAtTheCap (1.97s)
    shapeloop_integration_test.go:168: AdvanceShapeTask = {Task: Verified:false Wait: Stop:}, <nil>; want the wait "waiting for the Operator: : the review-cycle cap 1 is reached with findings open"
--- FAIL: TestNotMergeableStartsNoAttempt (2.29s)
    shapeloop_integration_test.go:195: AdvanceShapeTask = {Task: Verified:false Wait: Stop:}, <nil>; want the wait for the Operator
--- FAIL: TestAVerifierWithNoIndependentHarnessPostsNotActiveBeforeTheStop (3.47s)
    shapeloop_integration_test.go:214: AdvanceShapeTask = {Task: Verified:false Wait: Stop:}, <nil>; want the stop for not-active
--- FAIL: TestAnInterruptedStageIsFinishedFromItsKeptFilesAndRunsNoSecondSession (15.29s)
    shapeloop_integration_test.go:270: the post of the plan: the first call = nil; want the interruption
    shapeloop_integration_test.go:270: the post of the review record: the first call = nil; want the interruption
    shapeloop_integration_test.go:270: the handoff of the verifier: the first call = nil; want the interruption
    shapeloop_integration_test.go:270: the edit of the body: the first call = nil; want the interruption
--- FAIL: TestAMissingSectionMakesTheHandoffInvalidAndStops (4.08s)
    shapeloop_integration_test.go:338: AdvanceShapeTask = {Task: Verified:false Wait: Stop:}, <nil>; want the stop for the invalid spec-sections
--- FAIL: TestASpecificationThatPrdLintDoesNotAcceptStopsBeforeTheVerifier (4.51s)
    shapeloop_integration_test.go:353: AdvanceShapeTask = {Task: Verified:false Wait: Stop:}, <nil>; want the stop that names prd-lint
--- FAIL: TestADeveloperThatCommitsNothingMakesAnInvalidHandoffAndStops (2.62s)
    shapeloop_integration_test.go:366: AdvanceShapeTask = {Task: Verified:false Wait: Stop:}, <nil>; want the stop for the invalid handoff
--- FAIL: TestARestartAfterTheRowStartsTheFirstAttemptOnce (1.59s)
    shapeloop_integration_test.go:379: AdvanceShapeTask with no harness = {Task: Verified:false Wait: Stop:}, <nil>; want the stop for the pair
--- FAIL: TestTheCapOfMilestoneZeroIsAskedBeforeEachSession (1.81s)
    shapeloop_integration_test.go:396: AdvanceShapeTask with milestone 0 over its time = <nil>; want ledger.ErrMilestoneCap
--- FAIL: TestAShapeTaskIsOpenedOnlyWithTheLoopPermissions (1.69s)
    shapetask_integration_test.go:229: openShapeTask without the permissions = <nil>; want the permissions named
--- FAIL: TestARefusedOpeningWritesNoIssue (8.32s)
    shapetask_integration_test.go:283: a draft with no requirement: <nil>; want an error holding "spec/draft.tsv has no row of kind requirement or non-functional"
    shapetask_integration_test.go:283: a table of transitions the dry run refuses: <nil>; want an error holding "transitions.tsv"
    shapetask_integration_test.go:283: a file of parameters the dry run refuses: <nil>; want an error holding "parameters.tsv"
    shapetask_integration_test.go:283: a lost lease: <nil>; want an error holding "lease"
--- FAIL: TestAnOpenedShapeTaskIsReusedAndOpensNoSecondIssue (1.82s)
    shapetask_integration_test.go:302: the specification  then , 2 issues; want the first reused and 4 issues in all, the 2 of Start and one for each task
    handoffpost_integration_test.go:124: git show layup-records:tasks.tsv: exit status 128
--- FAIL: TestAStopBetweenTheIssueAndTheRowOpensASecondIssue (1.61s)
    shapetask_integration_test.go:318: openShapeTask with the lease lost after the issue = nil; want the error of the lost lease
--- FAIL: TestTheShapeIssueBlocksOfLoopMd (0.00s)
    shapetask_integration_test.go:357: 0 texts in the code; want the two kinds of Shape task
FAIL
FAIL	github.com/pharzam/layup/internal/run	55.079s
--- FAIL: TestUpdateRefSetsABranchOfLayupsNamespaceAndNoOther (0.11s)
    git_integration_test.go:990: git rev-parse refs/heads/layup/base/T-ab12/1: exit status 128
FAIL
FAIL	github.com/pharzam/layup/internal/git	0.126s
FAIL
exit=1
end 2026-10-19T10:53:34Z
```

### Green, the same commands on the code

```
start 2026-10-19T10:53:59Z
ok  	github.com/pharzam/layup/internal/run	0.446s
ok  	github.com/pharzam/layup/internal/git	0.044s
exit=0
end 2026-10-19T10:54:01Z
start 2026-10-19T10:54:01Z
ok  	github.com/pharzam/layup/internal/run	142.506s
ok  	github.com/pharzam/layup/internal/git	0.218s
exit=0
end 2026-10-19T10:56:26Z
```

The record of the first red run, of the intermediate failures that the first green runs showed and of what each was, is in the commits:
the lease of the world (the opening failed with `the records push was refused; the lease names the run`, a `LostError` of `holdsLease` on a
world whose Start released the lease); a refused start with no start row that the first form of the position never saw (a hundred acts,
now a stop and a unit case); the architecture head of a target with no PRD (below).

### A finding of the work (R10), not fixed here

The full check of `layup/spec` judges every head that touches `docs/prd/` or `docs/spec/` before the first bet, and the rule `covers`
**did not run: no PRD** when the head has none. The architecture task's head touches `docs/spec/` and, on a target whose default branch
has no PRD, holds none: its `layup/spec` is a failure and the driver stops (`<task>: layup/spec failed: covers: did not run: no PRD`).
The two Shape tasks start from the same base, so the architecture task cannot see the specification task's PRD. The tests give the
target a PRD and a section in its first tree. This is issue #432, owned by row 66h; it is off the path of this
row's Definition of Done.

## S9 - the full integration level

Command (both runs): `TMPDIR=/home/layup/.cache/layup-session/tmp-green go test -count=1 -tags=integration -timeout 90m ./...`

**Run 1**, on the head `2b19ae7` plus the record of the mutants. One failure, in an existing test whose premise this task changes:

```
start 2026-10-19T11:17:46Z
ok  	github.com/pharzam/layup/internal/git	6.431s
--- FAIL: TestEachStepsPromptFileFromARealRecordsTree (3.03s)
    promptfile_integration_test.go:303: git show :seen-prompt.txt: exit status 128
FAIL
FAIL	github.com/pharzam/layup/internal/run	1559.080s
FAIL
exit=1
end 2026-10-19T11:44:29Z
```

`TestEachStepsPromptFileFromARealRecordsTree` read the file that a session committed from its **bound** head. A step that only reads
(`plan-of-a-task`, `plan-review`, `verification`) now has no push and so no `bound`, which is the rule of this task (loop.md, The driver
of a Shape task). The test now reads the head from the ref that the end of the session fetched into the run's clone
(`refs/layup/sessions/<id>`), for all seven steps; it asserts what it asserted. This is the second of the three honest moves of
`docs/engineering-discipline.md` ("update the requirement the test encodes, with a written reason"): the requirement is the push rule
of loop.md and session.md, written in the first commit of this branch.

**Run 2**, on the head `9223f1c` (the fix of that test), nothing else changed:

```
start 2026-10-19T11:46:14Z
ok  	github.com/pharzam/layup/internal/git	5.912s
ok  	github.com/pharzam/layup/internal/run	1578.775s
exit=0
end 2026-10-19T12:13:10Z
```

## S7 - the compiled mutants of row 66r (the file `mutants.md` of the first head, folded in here by the fix cycle)

Run 2026-10-19 from 10:57:57Z to 11:13:42Z by `mutants.py` (in the session's scratchpad), on the committed green head `2b19ae7`
(commit `feat: T-jy6f the Shape tasks through the task loop to their verified heads`). Each mutant is one edit of one file, taken from a
fresh backup of the file (a copy made in a new directory at that mutant's start, never a fixed path), built (`go build ./...`) and
vetted (`go vet -tags=integration ./internal/run/ ./internal/git/`), then run with the test command of its row; a mutant that does not
build or vet counts as not run. The file is restored from the backup, its SHA-256 is checked to equal the hash before, and **the same test
command is run on the restored file as the control** (the column after the test status; 0 is green). The statuses are the exit statuses of
`go` read directly. `unit` is `go test -count=1 ./internal/run/ ./internal/git/`; the others are
`TMPDIR=/home/layup/.cache/layup-session/tmp-green go test -count=1 -tags=integration -timeout 20m -run '<tests>' ./internal/run/`. The
hashes are the first 12 hexadecimal characters of the SHA-256 of the file before and during the mutation.

| # | mutant | build | vet | test | first failing assertion | control | tests | sha256 before | sha256 during | restore |
| - | ------ | ----- | --- | ---- | ----------------------- | ------- | ----- | ------------- | ------------- | ------- |
| SM1 | permissions not checked before the issue | 0 | 0 | 1 | `shapetask_integration_test.go:232: 1 issues opened and the records head moved: false; want none and no commit` | 0 | `'^(TestAShapeTaskIsOpenedOnlyWithTheLoopPermissions)$'` | 3d7607898ae6 | 78f5d23ba790 | restored | 2026-10-19T10:57:57Z
| SM2 | a success criterion kept as a requirement | 0 | 0 | 1 | `shapetask_test.go:16: finalIDs = [REQ-001 NFR-002 REQ-003 REQ-004]; want [REQ-001 NFR-002 REQ-004]` | 0 | `unit` | 3d7607898ae6 | 8831a310eee0 | restored | 2026-10-19T10:58:26Z
| SM3 | a non-functional row mapped to REQ | 0 | 0 | 1 | `shapetask_test.go:16: finalIDs = [REQ-001 REQ-002 REQ-004]; want [REQ-001 NFR-002 REQ-004]` | 0 | `unit` | 3d7607898ae6 | eb1d18e8f342 | restored | 2026-10-19T10:58:48Z
| SM4 | the draw ignores tasks.tsv | 0 | 0 | 1 | `probe_test.go:118: taskID = "T-0000", <nil>; want T-1111, the second draw` | 0 | `unit` | c1f2906627ac | 276029ee5005 | restored | 2026-10-19T10:59:10Z
| SM5 | an opened task is not reused | 0 | 0 | 1 | `shapetask_integration_test.go:302: the specification T-9ut9 then T-7516, 5 issues; want the first reused and 4 issues in all, the 2 of Start and one for each task` | 0 | `'^(TestAnOpenedShapeTaskIsReusedAndOpensNoSecondIssue)$'` | 3d7607898ae6 | 16e13bf8df78 | restored | 2026-10-19T10:59:31Z
| SM6 | the dry run of registerFiles skipped | 0 | 0 | 1 | `shapetask_integration_test.go:286: a table of transitions the dry run refuses: 1 issues opened; want none` | 0 | `'^(TestARefusedOpeningWritesNoIssue)$'` | 3d7607898ae6 | 28e4f15eb67a | restored | 2026-10-19T11:00:01Z
| SM7 | the dry run of parameterFiles skipped | 0 | 0 | 1 | `shapetask_integration_test.go:286: a file of parameters the dry run refuses: 1 issues opened; want none` | 0 | `'^(TestARefusedOpeningWritesNoIssue)$'` | 3d7607898ae6 | d3d511d9e664 | restored | 2026-10-19T11:00:44Z
| SM8 | the lease not read before the issue | 0 | 0 | 1 | `shapetask_integration_test.go:283: a lost lease: <nil>; want an error holding "lease"` | 0 | `'^(TestARefusedOpeningWritesNoIssue)$'` | 3d7607898ae6 | e5f773349cc0 | restored | 2026-10-19T11:01:28Z
| SM9 | readOnlyStep without verification | 0 | 0 | 1 | `session_test.go:396: readOnlyStep(verification) = false; want true` | 0 | `unit` | b0c27a7e09eb | 7605fecc42f0 | restored | 2026-10-19T11:02:11Z
| SM10 | the push rule removed in TaskSession | 0 | 0 | 1 | `shapeloop_integration_test.go:85: AdvanceShapeTask(T-ljze) = {Task:T-ljze Verified:false Wait: Stop:T-ljze: the session S-7a86d837 is neither the plan review nor a verifier}, <nil>; want the verified head` | 0 | `'^(TestTheShapeTasksGoThroughTheLoopToTheirVerifiedHeads)$'` | 492cf0ededaa | 4c129d4f5031 | restored | 2026-10-19T11:02:32Z
| SM11 | the guard compares the task only | 0 | 0 | 1 | `session_test.go:424: another base: mkdir : no such file or directory; want "the prompt file of T-ab12 was built at base cccccccccccccccccccccccccccccccccccccccc, the session's base is 0123456789abcdef0123456789abcdef01234567"` | 0 | `unit` | 492cf0ededaa | 690eaf4acd13 | restored | 2026-10-19T11:03:23Z
| SM12 | the guard after the sweep | 0 | 0 | 1 | `session_test.go:427: another task: the calls ["sweep"] and 0 commits; want nothing swept, made or written` | 0 | `unit` | 492cf0ededaa | 01e8e09671e2 | restored | 2026-10-19T11:03:43Z
| SM13 | the verifier taken as the plan review | 0 | 0 | 1 | `shapeloop_test.go:94: specification, 13 facts: position 4 (""); want 8` | 0 | `unit` | b0c27a7e09eb | 9d465f2f83c0 | restored | 2026-10-19T11:04:03Z
| SM14 | the cycle cap read as 2 | 0 | 0 | 1 | `shapeloop_test.go:173: the second material round: {Pos:11 Session:S-v2 Msg: Attempt:2 Base:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa Head:cccccccccccccccccccccccccccccccccccccccc}; want the wait at the cap` | 0 | `unit` | b0c27a7e09eb | 9f37532ab3ab | restored | 2026-10-19T11:04:24Z
| SM15 | not mergeable starts an attempt | 0 | 0 | 1 | `shapeloop_test.go:187: not mergeable: {Pos:11 Session:S-v1 Msg: Attempt:1 Base:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa Head:bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb}; want the wait` | 0 | `unit` | b0c27a7e09eb | de234b4d2b2f | restored | 2026-10-19T11:04:43Z
| SM16 | a push with no result treated as no push | 0 | 0 | 1 | `shapeloop_test.go:133: a push with no result: 13 "T-sp01: the session S-d1 ended with no push (its resume is the step loop's)"; want 13 holding "a push with no result"` | 0 | `unit` | b0c27a7e09eb | 878e20bcf8c2 | restored | 2026-10-19T11:05:03Z
| SM17 | a second developer session in the attempt | 0 | 0 | 1 | `shapeloop_test.go:94: specification, 9 facts: position 6 (""); want 13` | 0 | `unit` | b0c27a7e09eb | d2886f5328f8 | restored | 2026-10-19T11:05:26Z
| SM18 | not-active does not stop | 0 | 0 | 0 | — | 0 | `'^(TestAVerifierWithNoIndependentHarnessPostsNotActiveBeforeTheStop)$'` | b0c27a7e09eb | 327ba5ed8576 | restored | 2026-10-19T11:05:46Z
| SM19a | (the first form did not build: an unused variable; not run) | 1 | 1 | — | — | — | unit | — | — | — |
| SM19b | the verifier taken for any base of the attempt | 0 | 0 | 0 | — | 0 | `unit` | b0c27a7e09eb | ba397356e357 | restored | 2026-10-19T11:13:21Z
| SM20 | the cap of milestone 0 not asked | 0 | 0 | 1 | `shapeloop_integration_test.go:396: AdvanceShapeTask with milestone 0 over its time = <nil>; want ledger.ErrMilestoneCap` | 0 | `'^(TestTheCapOfMilestoneZeroIsAskedBeforeEachSession)$'` | b0c27a7e09eb | 02036b6e8a62 | restored | 2026-10-19T11:06:31Z
| SM21 | the findings of the last round dropped | 0 | 0 | 1 | `shapeloop_test.go:206: findingsOf after a material round = [], ""; want its finding and no gate head` | 0 | `unit` | b0c27a7e09eb | 090af48bdcc8 | restored | 2026-10-19T11:07:09Z
| SM22 | UpdateRef accepts any ref | 0 | 0 | 1 | `git_test.go:597: a ref outside heads: <nil>, 1 starts of git; want a *FailedError of code -1 and none` | 0 | `unit` | d35dd1ea277b | 8d26cbf1fec8 | restored | 2026-10-19T11:07:29Z
| SM23 | the verifier clones the default branch | 0 | 0 | 1 | `shapeloop_integration_test.go:85: AdvanceShapeTask(T-llwk) = {Task: Verified:false Wait: Stop:}, git switch -c task/T-llwk/1 75c21b9c31f63c8ad2dca94e9e49ecd0efbf3831: exit status 128: fatal: unable to read tree (75c21b9c31f63c8ad2dca94e9e49ecd0efbf3831); want the verified head` | 0 | `'^(TestTheShapeTasksGoThroughTheLoopToTheirVerifiedHeads)$'` | b0c27a7e09eb | dcd4ce62b4bd | restored | 2026-10-19T11:07:53Z
| SM24 | the acceptance of a waiting specification skipped | 0 | 0 | 1 | `shapeloop_integration_test.go:353: AdvanceShapeTask = {Task:T-rjit Verified:true Wait: Stop:}, <nil>; want the stop that names prd-lint` | 0 | `'^(TestASpecificationThatPrdLintDoesNotAcceptStopsBeforeTheVerifier)$'` | b0c27a7e09eb | 6156edd99b73 | restored | 2026-10-19T11:08:46Z
| SM25 | an ended stage is run again | 0 | 0 | 1 | `shapeloop_test.go:94: specification, 3 facts: position 2 (""); want 3` | 0 | `unit` | b0c27a7e09eb | 9d6a277e7cf8 | restored | 2026-10-19T11:09:29Z
| SM26 | the verdict not posted before the wait | 0 | 0 | 1 | `shapeloop_integration_test.go:176: the newest status of 0c73c30 = [map[context:layup/verify creator:map[id:9001 login:layup-agent[bot]] description:waits for its turn state:pending] map[context:layup/rules creator:map[id:9001 login:layup-agent[bot]] description:no rule path state:success] map[contex` | 0 | `'^(TestAMaterialRoundStartsOneAttemptAndTheSecondWaitsAtTheCap)$'` | b0c27a7e09eb | 4cc95b7c7893 | restored | 2026-10-19T11:09:49Z
| SM27 | a refused start not seen | 0 | 0 | 1 | `shapeloop_test.go:133: a start that was refused: 2 ""; want 13 holding "the start of the session S-x1 was refused (context 1175 1000)"` | 0 | `unit` | b0c27a7e09eb | a19c9c658a26 | restored | 2026-10-19T11:10:44Z
| SM28 | an interrupted post of the review record reruns the verifier | 0 | 0 | 1 | `shapeloop_integration_test.go:273: AdvanceShapeTask(T-qv8n) = {Task: Verified:false Wait: Stop:}, the cap of milestone 0: the time since Start passes 8.0 hours; want the verified head` | 0 | `'^(TestAnInterruptedStageIsFinishedFromItsKeptFilesAndRunsNoSecondSession)$'` | b0c27a7e09eb | 8c5951e6fbd4 | restored | 2026-10-19T11:11:06Z

**Result of this first sweep.** 26 of the 28 runs are killed (SM1 to SM17 and SM20 to SM28); SM19a did not build and is not a run. Two mutants were **not killed**:

- **SM18** (`not-active` does not stop in the list of statuses): green at first. The path that follows reaches `ReviewerPair`, which gives the same stop
  text, so the output was the same for a head with no gate failure. Round 1 asked for the combined case; it is added and SM18 is killed (round 1, below).
- **SM19b** (the verifiers of the attempt are not narrowed to the bound head): green, and the reviewer accepts it as equivalent on the paths this driver
  produces. The classification of a `qa-engineer` session as a verifier already requires its base to be a bound SHA of a developer's session of its own
  attempt, and an attempt binds one head in every path the driver takes. It stays as a guard.

SM1 is killed by the check of the permissions made twice: the opening fails the second time, in `registerTask`, after the issue is
opened. SM8 is killed by the case of the lost lease. SM11 is killed by the first message of the guard that it keeps, and its run proves
that the missing check of the base reaches `make` (the first failing line is the call of the stand-in, which the guard stops before).

## Round 1 (comment 6267186962): the fix cycle

The record of the round is `not mergeable, findings recorded` at cycle 0 of cap 1. Findings 1, 3 and 4 of the record are fixed below; its
finding 2 (the budget) is NOT met by this cycle: folding `mutants.md` into this file did not lower the count, which was 26 files before and after (see round 2, below); its notes are applied (the single home of the
requirement mapping and the task names in `loop.md`, issue #432 above).

### Finding 1 - the admission probe loses the scratch branch

Red (a version that no probe passed, 1.0.0 recorded and 2.0.0 reported, at a base only the scratch branch holds), against `ProbeSpec.Branch`
declared and used by nothing. The control at the default base passes (a green negative):
`go test -count=1 -tags=integration -run '^TestTheAdmissionProbeOfASessionClonesTheBranchOfItsBase$' ./internal/run/`

```
start 2026-10-19T12:28:05Z
--- FAIL: TestTheAdmissionProbeOfASessionClonesTheBranchOfItsBase (0.34s)
    probe_integration_test.go:196: a base on the scratch branch: git switch -c task/T-bvcs/1 8e694d23e8c96197206840dca70a32c1ececca47: exit status 128: fatal: unable to read tree (8e694d23e8c96197206840dca70a32c1ececca47); want the admission probe to clone its base
FAIL
FAIL	github.com/pharzam/layup/internal/run	0.357s
FAIL
exit=1
end 2026-10-19T12:28:24Z
```

Green after the fix (`Probe` sets `Branch: p.Branch` on its `TaskSpec`; `startSession` gives `Admit` the branch):

```
start 2026-10-19T12:28:49Z
ok  	github.com/pharzam/layup/internal/run	0.374s
exit=0
end 2026-10-19T12:29:08Z
```

### Findings 3 and 4 - the recovery cases and the combined refusal

The new cases are green negatives at the code, which already did what they require; their **semantic failing controls** are the mutants below.
`TestAnInterruptedStageIsFinishedFromItsKeptFilesAndRunsNoSecondSession` gained the plan-review case and counts the sessions by role and the
comments of the plan, the plan review and the review record; `TestAFailedPostOfTheVerdictStatusIsPostedBeforeTheAttemptOrTheWait` injects a failed
post of the failure of `layup/verify` at the end of the first attempt and at the terminal wait; `TestAGateFailureWithNoIndependentVerifierStopsAndStartsNoAttempt`
has a gate failure and no independent verifier. Command and result on the code of `45005ed`:

```
start 2026-10-19T12:29:47Z
--- PASS: TestAnInterruptedStageIsFinishedFromItsKeptFilesAndRunsNoSecondSession (47.89s)
--- PASS: TestAFailedPostOfTheVerdictStatusIsPostedBeforeTheAttemptOrTheWait (29.06s)
--- PASS: TestAGateFailureWithNoIndependentVerifierStopsAndStartsNoAttempt (7.41s)
ok  	github.com/pharzam/layup/internal/run	84.367s
exit=0
end 2026-10-19T12:31:29Z
```

| # | mutant | build | vet | test | first failing assertion | control | tests |
| - | ------ | ----- | --- | ---- | ----------------------- | ------- | ----- |
| SM29 | the admission probe drops the branch | 0 | 0 | 1 | `probe_integration_test.go:196: a base on the scratch branch: git switch -c task/T-lxll/1 40054b992fe7df7d547afd54e294e370a60b0270: exit status 128: fatal: unable to read tree (40054b992fe7df7d547afd54e294e370a60b0270); want the admission probe to clone its base` | 0 | `'^(TestTheAdmissionProbeOfASessionClonesTheBranchOfItsBase)$'` |
| SM30 | the verdict not settled before the next attempt | 0 | 0 | 1 | `shapeloop_integration_test.go:443: the end of the first attempt: 2 attempts after the failed post; want 1 (no attempt starts before the failure is posted)` | 0 | `'^(TestAFailedPostOfTheVerdictStatusIsPostedBeforeTheAttemptOrTheWait)$'` |
| SM31 | a plan or plan review that ended is run again | 0 | 0 | 1 | `shapeloop_integration_test.go:277: AdvanceShapeTask(T-e764) = {Task: Verified:false Wait: Stop:}, the cap of milestone 0: the time since Start passes 8.0 hours; want the verified head` | 0 | `'^(TestAnInterruptedStageIsFinishedFromItsKeptFilesAndRunsNoSecondSession)$'` |
| SM18 | not-active does not stop | 0 | 0 | 1 | `shapeloop_integration_test.go:482: 2 attempts; want one: the gate failure starts none when the verifier cannot run` | 0 | `'^(TestAVerifierWithNoIndependentHarnessPostsNotActiveBeforeTheStop|TestAGateFailureWithNoIndependentVerifierStopsAndStartsNoAttempt)$'` |

Run 2026-10-19 from 12:33:54Z to 12:39:53Z on the head `9b0ebee`, by the same runner and with the rules of the first table (the control is the same test
command on the restored file, exit 0; the files were restored by hash). The first run of SM18 with the old test command (12:33:55Z) still passed, since
that test has no gate failure; the row above is SM18 with the combined case, and it is **killed**. SM29 drops the branch on the probe, SM30 settles nothing
before the next attempt and SM31 runs a plan or plan review again after it ended. The sweep above has no measured hashes; round 2 repeated it with them (below), and its count replaces the one that stood here.

### Round 1 - the full integration level on the final head

Command: `TMPDIR=/home/layup/.cache/layup-session/tmp-green go test -count=1 -tags=integration -timeout 90m ./...`, on the head `2d2fd3b` (the code and
tests of the fix cycle; only this record follows). The unit level, the e2e level, `go vet` (both forms), `gofmt -l` and the lints were run on it too, all exit 0.

```
start 2026-10-19T12:43:57Z
ok  	github.com/pharzam/layup/internal/git	6.258s
ok  	github.com/pharzam/layup/internal/run	1600.540s
exit=0
end 2026-10-19T13:11:18Z
```

## Round 2 (comment from the Operator, cap 2): the fix cycle

Round 2 reads `fa7a476` as `not mergeable, findings recorded` at cycle 1 of cap 1; the Operator chose one more cycle (cap 2). Its findings:

- **Finding 1, the budget.** The statement above, that finding 2 of round 1 was met, was false: the diff was 26 files (over 25) before the fold and 26 after it.
  The file count is 26 before the close-out, over the reviewed maximum of 25, and the landing figure (files and lines, the close-out included) goes to
  the Operator as one exact measure. This cycle adds no file.
- **Finding 2, a bound of a session outside the developer role.** Red, on the unit test `TestABoundOfASessionOutsideTheDeveloperRoleStopsTheStep`
  (the ordinary completed-verifier ledger is the control and passes; the same ledger with the verifier's own push and bound, and the same for the plan and
  the plan review, returned the verified head):

```
start 2026-10-19T13:55:18Z
--- FAIL: TestABoundOfASessionOutsideTheDeveloperRoleStopsTheStep (0.00s)
    shapeloop_test.go:232: a bound of the verifier: {Pos:10 Session:S-v1 Msg: Attempt:1 Base:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa Head:bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb}; want the stop that names S-v1
    shapeloop_test.go:232: a bound of the plan: {Pos:10 Session:S-v1 Msg: Attempt:1 Base:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa Head:bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb}; want the stop that names S-p1
    shapeloop_test.go:232: a bound of the plan review: {Pos:10 Session:S-v1 Msg: Attempt:1 Base:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa Head:bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb}; want the stop that names S-r1
FAIL
FAIL	github.com/pharzam/layup/internal/run	0.009s
FAIL
exit=1
end 2026-10-19T13:55:31Z
```

  Green after the check in `loopPosition` (a `bound` whose session is not of the developer role stops the step, before any stage):

```
start 2026-10-19T13:55:36Z
ok  	github.com/pharzam/layup/internal/run	0.389s
ok  	github.com/pharzam/layup/internal/git	0.033s
exit=0
end 2026-10-19T13:55:47Z
```

- **Finding 3, the hashes.** SM18 (with the combined case), SM29, SM30 and SM31 are repeated with the new mutant SM33 (the check removed), by the runner of the
  first sweep, on the head `e3b11b2`, 2026-10-19T13:56:39Z to 2026-10-19T14:02:18Z. For each, the build and vet statuses, the test status
  and the first failing assertion are below, and then the SHA-256 of the whole source file before the mutation, with the mutation, and after the restore, with
  the equality of the first and the last; the control is the same test command on the restored file (exit 0 for each).

| # | mutant | build/vet/test | first failing assertion | control |
| - | ------ | -------------- | ----------------------- | ------- |
| SM18 | not-active does not stop | 0/0/1 | shapeloop_integration_test.go:482: 2 attempts; want one: the gate failure starts none when the verifier cannot run | 0 |
| SM29 | the admission probe drops the branch | 0/0/1 | probe_integration_test.go:196: a base on the scratch branch: git switch -c task/T-hcie/1 8b58ba6ab488290838d26a0e031be217ca1c1ec5: exit status 128: fatal: unable to read tree (8b58ba6ab488290838d26a0e031be217ca1c1ec5); want the admission pr | 0 |
| SM30 | the verdict not settled before the next attempt | 0/0/1 | shapeloop_integration_test.go:443: the end of the first attempt: 2 attempts after the failed post; want 1 (no attempt starts before the failure is posted) | 0 |
| SM31 | a plan or plan review that ended is run again | 0/0/1 | shapeloop_integration_test.go:277: AdvanceShapeTask(T-6j6a) = {Task: Verified:false Wait: Stop:}, the cap of milestone 0: the time since Start passes 8.0 hours; want the verified head | 0 |
| SM33 | a bound of a session outside the developer role not seen | 0/0/1 | shapeloop_test.go:232: a bound of the verifier: {Pos:10 Session:S-v1 Msg: Attempt:1 Base:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa Head:bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb}; want the stop that names S-v1 | 0 |

| # | sha256 before | sha256 during (mutated) | sha256 after the restore | before = after |
| - | ------------- | ----------------------- | ------------------------ | -------------- |
| SM18 | `a1fc3c879c557d54f80c80572f59a791036a3086cb62c4757f4cdad64367712c` | `7e03a254cdec9754afeb2e6522a35c52e86c4fe014db06656d78e2dda492756e` | `a1fc3c879c557d54f80c80572f59a791036a3086cb62c4757f4cdad64367712c` | equal |
| SM29 | `2cbf939eb1ad6e41bb4dc2571c2c76515ba2985db0a0a10c556d2115e3f9c207` | `d93f290491be01ca79c2817f0eaa87a7221026ae41f54751cc350e01e4ec3056` | `2cbf939eb1ad6e41bb4dc2571c2c76515ba2985db0a0a10c556d2115e3f9c207` | equal |
| SM30 | `a1fc3c879c557d54f80c80572f59a791036a3086cb62c4757f4cdad64367712c` | `3944cf84c397c1326874de59f05c82f2c5347e6f1f0405ba27fd6d344acfcebf` | `a1fc3c879c557d54f80c80572f59a791036a3086cb62c4757f4cdad64367712c` | equal |
| SM31 | `a1fc3c879c557d54f80c80572f59a791036a3086cb62c4757f4cdad64367712c` | `37999572542a0c03dd96d875f7a0de9cbd2e31e5bdba395980a8aaad46e46e61` | `a1fc3c879c557d54f80c80572f59a791036a3086cb62c4757f4cdad64367712c` | equal |
| SM33 | `a1fc3c879c557d54f80c80572f59a791036a3086cb62c4757f4cdad64367712c` | `c2db02be353c55567a9126d9eedc03842cd39141cb0584e60fd2b545f98aece2` | `a1fc3c879c557d54f80c80572f59a791036a3086cb62c4757f4cdad64367712c` | equal |

**The count, recounted.** The mutants of S7 are SM1 to SM28 and SM19 (a first form that did not build, SM19a, and a second, SM19b), SM29 to SM31 and SM33.
Killed: SM1 to SM17 and SM20 to SM28 (26) in the first sweep, SM18 (with the combined case), SM29, SM30, SM31 and SM33 (5): **31 killed**. Not killed:
**SM19b** alone, equivalent on the paths this driver produces (S7). SM19a did not build and is not a run. (SM32 was a duplicate of SM26 and was never run.)
