# T-jy6f - the compiled mutants of row 66r

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

**Result.** 26 of the 28 runs are killed (SM1 to SM17 and SM20 to SM28); SM19a did not build and is not a run. Two mutants are **not killed**, and neither is
called equivalent by assertion:

- **SM18** (`not-active` does not stop in the list of statuses): the test is green. The path that follows reaches `ReviewerPair`, which
  gives the same error and the same stop text, so the output is the same for a head whose statuses are otherwise good. What the earlier
  case adds is the precedence of the stop over the retry of a gates failure and over the acceptance of the specification; no test
  makes a head that fails the gates and has no independent harness.
- **SM19b** (the verifiers of the attempt are not narrowed to the bound head): the test is green. The classification of a `qa-engineer` session
  as a verifier already requires its base to be a bound SHA of a developer's session of its own attempt, and an attempt binds one head in
  every path the driver takes, so the narrowing cannot change a result there. It stays as a guard.

SM1 is killed by the check of the permissions made twice: the opening fails the second time, in `registerTask`, after the issue is
opened. SM8 is killed by the case of the lost lease. SM11 is killed by the first message of the guard that it keeps, and its run proves
that the missing check of the base reaches `make` (the first failing line is the call of the stand-in, which the guard stops before).
