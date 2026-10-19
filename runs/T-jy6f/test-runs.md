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
target a PRD and a section in its first tree. This is reported on #415 for the plan of rows 66h and 66k; it is off the path of this
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
