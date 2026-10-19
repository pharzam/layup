# Test runs of T-syxm (#422)

Every run is a direct command with its own exit status, never piped. Times are read from `date -u` around
the command. The host was loaded by another agent's integration run (times are not a measure).

## S2, the inert stubs and the world (commit bd1bcfa), all green

- `go build ./...` exit 0; `go vet ./...` exit 0; `go vet -tags=integration ./...` exit 0; `go vet -tags=uat ./internal/run/` exit 0.
- `go test ./...` exit 0.
- `TMPDIR=/home/layup/.cache/layup-session/tmp-green go test -count=1 -tags=integration -timeout 45m ./internal/git/ ./internal/run/ -run 'Session|Push|End|Comment|Probe|Phase|QA|Store|Remote|Ls'`
  start 2026-10-18T18:10:44Z, end 2026-10-18T18:22:01Z, exit 0 (`internal/run` 673 s). A first run of a wider pattern
  without `-timeout` was killed by the default 10-minute limit of `go test` (a panic of the timeout, not a failure of a test); it is not evidence.

## S3, the red run against the stubs (commit bd1bcfa plus the new test files, uncommitted)

Unit, start 2026-10-18T18:25:05Z, end 2026-10-18T18:25:06Z, exit 1:

`go test -count=1 ./internal/git/ ./internal/run/ -run 'TestRemoteHead|TestAPush|TestABranchHead|TestAFailedBranchRead|TestAnOpenAttempt|TestSeveralLost|TestAClosedOrRebased|TestTheNewAttempt'`

```
--- FAIL: TestRemoteHeadRunsItsVerb (0.00s)
    --- FAIL: TestRemoteHeadRunsItsVerb/ls-remote_of_a_branch (0.00s)
        remote_test.go:33: 0 starts of git, want 1
    --- FAIL: TestRemoteHeadRunsItsVerb/ls-remote_of_a_branch_with_a_token (0.00s)
        remote_test.go:33: 0 starts of git, want 1
--- FAIL: TestRemoteHeadRefusesABranchBeforeGitStarts (0.00s)
    remote_test.go:54: branch "": found false, err <nil>, 0 starts; want a *FailedError of code -1 and no start
    remote_test.go:54: branch "a:b": found false, err <nil>, 0 starts; want a *FailedError of code -1 and no start
    remote_test.go:54: branch "-x": found false, err <nil>, 0 starts; want a *FailedError of code -1 and no start
--- FAIL: TestRemoteHeadReadsOnlyTheExactRefLine (0.01s)
    --- FAIL: TestRemoteHeadReadsOnlyTheExactRefLine/the_exact_line (0.00s)
        remote_test.go:90: RemoteHead = "", false; want "0123456789abcdef0123456789abcdef01234567", true
    --- FAIL: TestRemoteHeadReadsOnlyTheExactRefLine/the_exact_line_after_a_longer_ref (0.00s)
        remote_test.go:90: RemoteHead = "", false; want "0123456789abcdef0123456789abcdef01234567", true
    --- FAIL: TestRemoteHeadReadsOnlyTheExactRefLine/exit_128 (0.00s)
        remote_test.go:87: err <nil>; want code 128
    --- FAIL: TestRemoteHeadReadsOnlyTheExactRefLine/a_short_ID_on_the_exact_line (0.00s)
        remote_test.go:87: err <nil>; want code 0
    --- FAIL: TestRemoteHeadReadsOnlyTheExactRefLine/an_upper-case_ID_on_the_exact_line (0.00s)
        remote_test.go:87: err <nil>; want code 0
--- FAIL: TestAPushWhoseBranchHoldsItsSHAIsBound (0.00s)
    repair_test.go:96: events = [attempt session result push], <nil>; want [attempt session result push bound]
--- FAIL: TestAPushWhoseBranchIsAbsentIsPushLostAndStartsOneAttempt (0.00s)
    repair_test.go:109: events = [attempt session result push], <nil>; want [attempt session result push refused attempt]
--- FAIL: TestABranchHeadThatDiffersFromThePushIsPushLost (0.00s)
    repair_test.go:127: events [["1" "attempt" "1" "" "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" "" "" "2026-10-09T12:00:00Z"] ["2" "session" "1" "S-1a2b3c4d" "" "" "" "2026-10-09T12:00:00Z"] ["3" "re [cut: the event rows follow]
--- FAIL: TestAFailedBranchReadIsAnErrorAndWritesNothing (0.00s)
    repair_test.go:135: err = <nil>; want the read's error
--- FAIL: TestAPushLostInAnOlderAttemptStartsNoAttempt (0.00s)
    repair_test.go:152: events [["1" "attempt" "1" "" "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" "" "" "2026-10-09T12:00:00Z"] ["2" "session" "1" "S-1a2b3c4d" "" "" "" "2026-10-09T12:00:00Z"] ["3" "re [cut: the event rows follow]
--- FAIL: TestAnOpenAttemptHoldingPushLostGetsExactlyOneAttempt (0.00s)
    repair_test.go:167: after call 1: 1 attempts; want 2
    repair_test.go:167: after call 2: 1 attempts; want 2
    repair_test.go:171: 0 commits and 0 reads; want one commit (the attempt) and no read
--- FAIL: TestSeveralLostPushesOfTheOpenAttemptGiveOneAttempt (0.00s)
    repair_test.go:182: events [attempt session result push session result push], Repair {Bound:[] Lost:[] Posted:[] Attempt:0}, <nil>; want two refused and attempt 2 once
--- FAIL: TestAClosedOrRebasedLatestAttemptGetsNoReplacement (0.00s)
    --- FAIL: TestAClosedOrRebasedLatestAttemptGetsNoReplacement/closed (0.00s)
        repair_test.go:196: events [attempt session result push closed], Repair {Bound:[] Lost:[] Posted:[] Attempt:0}, <nil>; want the refusal and no new attempt
    --- FAIL: TestAClosedOrRebasedLatestAttemptGetsNoReplacement/rebased (0.00s)
        repair_test.go:196: events [attempt session result push rebased], Repair {Bound:[] Lost:[] Posted:[] Attempt:0}, <nil>; want the refusal and no new attempt
--- FAIL: TestTheNewAttemptStartsAtTheLostAttemptsBase (0.00s)
    repair_test.go:225: events [["1" "attempt" "1" "" "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" "" "" "2026-10-09T12:00:00Z"] ["2" "session" "1" "S-1a2b3c4d" "" "" "" "2026-10-09T12:00:00Z"] ["3" "re [cut: the event rows follow]
```

Integration, start 2026-10-18T18:25:06Z, end 2026-10-18T18:27:51Z, exit 1:

`TMPDIR=/home/layup/.cache/layup-session/tmp-green go test -count=1 -tags=integration -timeout 20m ./internal/git/ ./internal/run/ -run 'TestRemoteHeadOf|TestRemoteHeadSends|TestARestart|TestAMoved|TestATaskSessionsComment|TestASessionComment|TestAMissingSession|TestNoRegister|TestTheGuards|TestTheSameText'`

```
--- FAIL: TestRemoteHeadOfABareRepository (0.29s)
    remote_integration_test.go:29: a branch that is there: "", false, <nil>; want 089ff9e9ab70fceaa7cc3a91110b0e7f895370ba, true, nil
    remote_integration_test.go:41: a remote that does not exist: <nil>, found false; want a *FailedError of code 128
--- FAIL: TestRemoteHeadSendsTheTokenAsAHeaderOnly (0.00s)
    remote_integration_test.go:60: a server with no repository: no error
    remote_integration_test.go:63: the server got Authorization "", want Basic of x-access-token and the token
--- FAIL: TestARestartBindsAPushThatTheForgeAccepted (4.73s)
    repair_integration_test.go:143: events [attempt session result push], Repair {Bound:[] Lost:[] Posted:[] Attempt:0}, <nil>; want [attempt session result push bound] bound
--- FAIL: TestARestartRefusesALostPushAndStartsOneAttempt (4.85s)
    repair_integration_test.go:155: events [attempt session result push], Repair {Bound:[] Lost:[] Posted:[] Attempt:0}, <nil>; want [attempt session result push refused attempt]
--- FAIL: TestARestartRepairsAShapePushBeforeAnyBet (5.76s)
    repair_integration_test.go:185: events [attempt session result push]; want [attempt session result push refused attempt]
--- FAIL: TestARestartBetweenPushLostAndTheAttemptStartsExactlyOne (5.39s)
    repair_integration_test.go:206: the stop: <nil>; want its error
--- FAIL: TestAMovedBranchHeadIsPushLost (5.84s)
    repair_integration_test.go:234: events [attempt session result push]; want [attempt session result push refused attempt]
--- FAIL: TestATaskSessionsCommentIsOnItsIssue (5.77s)
    repair_integration_test.go:242: the comments of the task's issue []; want ["S-84d5a440: developer session of T-ab12, attempt 1: done\n"]
    repair_integration_test.go:245: the control issue holds ["S-84d5a440: developer session of T-ab12, attempt 1: done\n"]; want no comment of a task session
--- FAIL: TestASessionCommentThatTheAppPostedIsNotPostedAgain (5.91s)
    repair_integration_test.go:254: Repair {Bound:[] Lost:[] Posted:[] Attempt:0}, <nil>, 0 comments; want none posted
--- FAIL: TestAMissingSessionCommentIsPostedAgainFromTheRecords (5.24s)
    repair_integration_test.go:271: call 1: Repair {Bound:[] Lost:[] Posted:[] Attempt:0}, <nil>, comments []; want ["S-eab708ca: developer session of T-ab12, attempt 1: done\n"] posted by the first call only
    repair_integration_test.go:271: call 2: Repair {Bound:[] Lost:[] Posted:[] Attempt:0}, <nil>, comments []; want ["S-eab708ca: developer session of T-ab12, attempt 1: done\n"] posted by the first call only
--- FAIL: TestNoRegisterRowOrABrokenRegisterPostsNowhere (16.57s)
    --- FAIL: TestNoRegisterRowOrABrokenRegisterPostsNowhere/no_row (5.35s)
        repair_integration_test.go:302: TaskSession: "S-bc751050", <nil>; want an error with "has no row of T-ab12"
    --- FAIL: TestNoRegisterRowOrABrokenRegisterPostsNowhere/refused (5.52s)
        repair_integration_test.go:302: TaskSession: "S-c53ed9c6", <nil>; want an error with "tasks.tsv:"
    --- FAIL: TestNoRegisterRowOrABrokenRegisterPostsNowhere/absent (5.69s)
        repair_integration_test.go:302: TaskSession: "S-5e4e9f9a", <nil>; want an error with "no tasks.tsv"
--- FAIL: TestTheGuardsOfThePostHoldOnBothPaths (33.40s)
    --- FAIL: TestTheGuardsOfThePostHoldOnBothPaths/a_failed_read_of_the_comments,_the_session's_path (5.24s)
        repair_integration_test.go:344: TaskSession "S-5ac598ad": <nil>, 0 comments; want the guard's error and no post
    --- FAIL: TestTheGuardsOfThePostHoldOnBothPaths/a_failed_read_of_the_comments,_the_repost (5.18s)
        repair_integration_test.go:358: RepairTask: <nil>, 0 comments; want the guard's error and no post
    --- FAIL: TestTheGuardsOfThePostHoldOnBothPaths/a_failed_read_of_the_register,_the_session's_path (5.18s)
        repair_integration_test.go:344: TaskSession "S-e8418edf": <nil>, 0 comments; want the guard's error and no post
    --- FAIL: TestTheGuardsOfThePostHoldOnBothPaths/a_failed_read_of_the_register,_the_repost (4.98s)
        repair_integration_test.go:358: RepairTask: <nil>, 0 comments; want the guard's error and no post
    --- FAIL: TestTheGuardsOfThePostHoldOnBothPaths/a_lease_lost_before_the_post,_the_session's_path (6.54s)
        repair_integration_test.go:344: TaskSession "S-0abcc179": <nil>, 0 comments; want the guard's error and no post
    --- FAIL: TestTheGuardsOfThePostHoldOnBothPaths/a_lease_lost_before_the_post,_the_repost (6.28s)
        repair_integration_test.go:358: RepairTask: <nil>, 0 comments; want the guard's error and no post
--- FAIL: TestTheSameTextFromAPersonOrAnotherAppDoesNotSuppressTheSummary (5.16s)
    repair_integration_test.go:386: after call 1: 2 comments of this line; want the two foreign ones and ours, once (session S-5f9677c9)
    repair_integration_test.go:386: after call 2: 2 comments of this line; want the two foreign ones and ours, once (session S-5f9677c9)
```

Observed outcomes against the plan's predictions:

- Green at the stub, as the green negatives and controls: `TestAPushThatIsAlreadyBoundOrRefusedIsNotReadAgain` (both fixtures, bound and refused: no read, no commit), and the subtests `only a longer ref` and `exit 2` of `TestRemoteHeadReadsOnlyTheExactRefLine` (absence is the stub's answer).
- `TestAFailedBranchReadIsAnErrorAndWritesNothing` is red on `err = <nil>`; its `commits = 0` half is green at the stub, as predicted.
- `TestASessionCommentThatTheAppPostedIsNotPostedAgain` (plan I7) is **red** at the stub, not green as the plan predicted: it asserts the one summary on the task's issue, which the control-issue route does not post there (`0 comments; want none posted`).
- `TestATaskSessionsCommentIsOnItsIssue` (I6) is red on both assertions: the task's issue holds no summary, and the control issue holds it (plan review condition 5: the control-issue summary is red, not a green negative).
- In `TestTheGuardsOfThePostHoldOnBothPaths` the "repost" subtests' first `TaskSession` still fails at the stub, because the forge's one-time `failNote` fires on the control-issue post; the red assertion is the guard's error that never comes.
- The assertion of `TestAPushLostInAnOlderAttemptStartsNoAttempt` named the 5th event; the refusal of attempt 1 is appended after attempt 2's event, so it is the 6th. The index was corrected before the code, and the assertion fails at the stub the same way (5 events, no refusal).
- The pre-commit hook ran the unit level on the working tree when the first commit of S2 was tried with the new test files present, and refused it with the same failures (`go test ./...` exit 1); the new files were set aside, S2 committed green, and the files put back.

## S4, green

- The new tests, run alone at fb771b6 (2026-10-18T18:30:47Z to 18:34:23Z): `TMPDIR=/home/layup/.cache/layup-session/tmp-green go test -count=1 -tags=integration -timeout 20m ./internal/git/ ./internal/run/ -run 'TestRemoteHeadOf|TestRemoteHeadSends|TestARestart|TestAMoved|TestATaskSessionsComment|TestASessionComment|TestAMissingSession|TestNoRegister|TestTheGuards|TestTheSameText'`, exit 0 (`internal/run` 185.7 s). The unit tests of the same commit: `go test -count=1 ./internal/git/ ./internal/run/` exit 0.
- **The first full run of the integration level failed, and is the cost of a plan that missed the lease.** At fb771b6, in a copy made with `git archive`:
  `TMPDIR=/home/layup/.cache/layup-session/tmp-green go test -count=1 -tags=integration -timeout 45m ./...`, 2026-10-18T18:38:27Z to 19:21:45Z, exit 1. 26 tests of `internal/run` failed with
  `TaskSession: "S-…", the records push was refused; the lease names the run 0123456789abcdef`: the world of the session tests wrote no lease, and the comment reads it before the post (the plan review's Q8, and the plan's S2, moved the register but not the lease). `TestTheScriptsOfLAYUP` (`internal/verify`) failed only because the `git archive` copy has no `.git` to clone; it is not a failure of the change. Fixed in ac05f6a: `sessionWorldWith` writes the lease (moved from `commentWorld`).
- The full integration level at ac05f6a, in a `git clone` of the worktree: `TMPDIR=/home/layup/.cache/layup-session/tmp-green go test -count=1 -tags=integration -timeout 45m ./...`, 2026-10-18T19:23:40Z to 19:57:06Z, **exit 0** (`internal/run` 1923.7 s, `internal/verify` 92.2 s, `internal/git` 8.4 s). `TestProcessTimesTheFirstOutput` did not fail.
- At ac05f6a in the worktree: `go build ./...`, `go vet ./...`, `go vet -tags=integration ./...`, `go vet -tags=e2e ./cmd/layup/`, `go vet -tags=uat ./internal/run/`, `gofmt -l .` (no output), `go test -count=1 ./...` (19:26:01Z to 19:26:22Z) and `go test -count=1 -tags=e2e ./cmd/layup/` (193.3 s) all exit 0; `adr-lint`, `prd-lint`, `link-lint`, `run-discipline-tests`, `nested-checkout-check`, `setup-check` and `git diff --check` exit 0.

## Round 1 fix cycle (#422, head 9260f17), test first

Finding 1, `Repair.Attempt` after a failed `StartAttempt`.
- Red, at 9260f17 plus the new test `TestAReplacementThatWasNotCommittedIsNotReportedAsWritten` (2026-10-18T20:43:24Z, the time of the output file):
  `go test -count=1 ./internal/run/ -run TestAReplacementThatWasNot`, exit 1; the first failing assertion, verbatim:
  `repair_test.go:241: Repair.Attempt = 2 with no replacement committed; want 0`. The earlier assertion (the error, no commit, one attempt event) held.
- Fix (`internal/run/repair.go`): `out.Attempt` is set only when `StartAttempt` returned no error.
- Green: the same command, exit 0; the retry commits attempt 2 once and a third call commits nothing.

Finding 3, a refused `sessions.tsv`. `TestARefusedSessionsTableStopsTheRepostAndTheRepair` (RepostComments and RepairTask) is **a green negative**: green on the
code as it stands (`go test -count=1 ./internal/run/ -run 'TestARefusedSessions|TestAReplacementThatWasNot'`, exit 0, written before any code change), with its red mutant M23
(`runs/T-syxm/mutants.md`): build, vet and vet-int exit 0, the unit test exit 1 at `repair_test.go:274: RepostComments: posted [], <nil>, 0 commits; want an error that names sessions.tsv, ...`.

Finding 2, the artifact selection. `commentsOf(w, body)` counts by the exact body, `commentByID` and `editComment(w, id, body)` select by the returned `p.ID`; an unrelated comment (ID 7001) is
put before the artifact in `TestAnEditedCommentMakesTheSecondCopyAndTheNewHandoffNamesIt`. This is **a green negative** on the new helpers; its red mutant M25 puts the old selection
back (first comment that is no summary): the test exits 1 at `handoffpost_integration_test.go:206: after the edit the unrelated comment is map[body:edited by a person ...`. The posting, copy, hash and adoption assertions are kept.
- Green: `TMPDIR=/home/layup/.cache/layup-session/tmp-green go test -count=1 -tags=integration -timeout 90m -run '^(<the ten tests of handoffpost_integration_test.go>)$' ./internal/run/`,
  2026-10-18T20:46:53Z to 20:48:29Z, exit 0 (94.3 s).
- Gates at 7640a2c (to 20:51:56Z), all exit 0: `go build ./...`, `go vet ./...`, `go vet -tags=integration ./...`, `go vet -tags=e2e ./cmd/layup/`, `go vet -tags=uat ./internal/run/`, `gofmt -l .` (no output),
  `go test -count=1 ./...`, the five text lints, `nested-checkout-check`, `setup-check`, `git diff --check origin/main HEAD`.

## End 3 fix (#422, comment 6260569452), test first

Round 2, finding 1: the refused-events path had no protecting repair test. `TestARefusedEventsTableStopsTheRepairAndTheRepost` (`RepairPushes` and `RepostComments`, an `events.tsv` of one unparseable line)
is **a green negative**: written before any code change, green on the code as it stands, red on its mutants M26 and M27 (`runs/T-syxm/mutants.md`), which are the red evidence.
- Green on the unchanged code, `go test -count=1 ./internal/run/ -run TestARefusedEventsTable`, 2026-10-19T05:46:40Z (start; the run took 0.016 s), exit 0, `ok  github.com/pharzam/layup/internal/run`.
- Red, M26 (both events reads of `RepairPushes` swallowed): build, vet and vet-int exit 0; the same command, exit 1; `--- FAIL: TestARefusedEventsTableStopsTheRepairAndTheRepost/RepairPushes (0.00s)`,
  `repair_test.go:299: RepairPushes: <nil>, 0 commits, 0 branch reads; want an error that names events.tsv, no commit, no read`.
- Red, M27 (the events read of `RepostComments` swallowed): the same, `.../RepostComments`, the same assertion with `RepostComments`.
- Not a kill: swallowing only the first events read of `RepairPushes` leaves the test green, since the second read of the function returns the same error; not counted.
- Green after the restore (SHA 1e3303f4... checked), the same command, exit 0, ended 2026-10-19T05:47:40Z. No production code changed.

## The full integration level at the merge head 9260f17, a failure

`TMPDIR=/home/layup/.cache/layup-session/tmp-green go test -tags integration -timeout 90m ./...` on the merge head 9260f17, 2026-10-18T20:28:40Z to 21:02:29Z, **exit 1**. Every package passed except
`internal/run` (`FAIL ... 1961.629s`), with two failing tests, both load flakes of this host and not of the change:
- `TestAProseSessionWhoseCommentFailsLeavesNoFileInTheRecords` (55.80 s): the `fastClock` flake logged on #291 (comment 6255116905): `internal/run/run_integration_test.go:173` multiplies the real elapsed time by 600,
  so about 48 s under host load crosses the 8-hour cap.
- `TestTheFullInputOfAHeadAfterIntakeAndTheScaffold` (52.71 s): `specinput_integration_test.go:142: git push -q .../acme/target.git layup-setup:main: exit status 128`, a push to a directory that was missing under load.

Both tests alone at 7f1ba6d (the head of the fix cycle, a merge-free tree of the same code), 2026-10-18T21:03:14Z to 21:04:16Z, **exit 0**: `--- PASS: TestAProseSessionWhoseCommentFailsLeavesNoFileInTheRecords (10.14s)`,
`--- PASS: TestTheFullInputOfAHeadAfterIntakeAndTheScaffold (50.00s)`, `ok  github.com/pharzam/layup/internal/run  60.163s`. The full level has therefore **not** been green in one run at the merge head;
its only green full run is at ac05f6a (above), and the merge head's failures are the two flakes with the green alone-runs.

## The full integration level at bf6d7ac (the end 3 fix), green

In a `git clone` of the worktree at bf6d7ac (the code of the final head; bbe2689 adds documents only): `TMPDIR=/home/layup/.cache/layup-session/tmp-green go test -count=1 -tags integration -timeout 90m ./...`,
2026-10-19T05:49:22Z to 06:21:42Z, **exit 0**, every package `ok` (`internal/run` 1832.049 s, `internal/setup` 68.557 s, `internal/verify` 71.237 s). Neither load flake of the merge head failed in this run.
Also exit 0 at bbe2689: `go build ./...`, `go vet ./...`, `go vet -tags=integration ./...`, `gofmt -l .` (no output), `go test ./...`, the five text lints, `nested-checkout-check`, `setup-check`, `git diff --check origin/main HEAD`.
