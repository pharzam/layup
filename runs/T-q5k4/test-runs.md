# T-q5k4 test runs (row 66d of M2e, #360)

Base `e16568c` (`main`). Author Claude Sonnet 5.5 on Claude Code, no other model and no subagent. Every time is
UTC from `date -u`, taken by the runner `x.sh` after the command; every exit code is the command's own (the
command writes its log to a file, so no pipe hides it). The runner `~/.cache/layup-session/q5k4/x.sh` and the mutation
runner `~/.cache/layup-session/q5k4/mut.py` are host scripts, not part of the change.

## Red and green, by slice

### Slice A: the kind `answer` (internal/records, records.md)
- RED 2026-10-16T12:15:47Z..12:15:48Z exit=1 `go test -count=1 ./internal/records/` at code reverted to base: TestAValidResultIsRead and TestResultRefusesEachBrokenRule FAIL; 38 assertion lines, all `"answer" is not a value of enum(...)` / want an error in column X (log /var/tmp/q-A-red.log)
- GREEN 2026-10-16T12:14:56Z exit=0 go test ./internal/records/ (code applied); 12:15:41Z exit=0 go test -count=1 -tags=integration ./internal/tsv/ (block equality with records.md)
### Slice B: session.CommentFile
- RED (stub = os.ReadFile, "has the step, not the behaviour") 2026-10-16T12:17:00Z..12:17:01Z exit=1: TestACommentFileThatIsNotARegularFileOfOneNameIsRefused fails on link.md (read 19 bytes), big.md (read 1048577 bytes), plan.md/ (read 5 bytes); TestACommentFileIsReadByteForByte green on arrival with the reading stub (declared; its mutants are M-B*). First stub (always error) gave 7 failures incl. the byte-for-byte cases at 12:16:54Z exit=1.
- GREEN 12:17:06Z exit=0 go test -count=1 ./internal/session/
### Slice C: the end check and the kept copy (internal/run commentfile.go, end.go)
- RED unit (stubs: KeptPath "", CheckCommentBody nil, keepComments nil/false) 12:18:21Z..12:18:25Z exit=1: TestTheKeptPathIsTheTaskTheSessionAndTheFile, TestACommentBodyIsRefusedBeforeThePost (4 refusals missing), TestTheKeptBytesAreTheBytesOfTheFile (7), TestACommentRowIsCheckedByItsHashAgainstTheResultFile (8) fail by assertion text
- RED integration 12:18:54Z..12:19:15Z exit=1: TestTheCommentRowOfASessionIsKeptWithTheResultInOneCommit (kept file does not exist), TestACommentFileThatFailsItsRowIsRefusedArtifactAndNotKept x3 (events lack `refused artifact`); TestAResultThatHoldsACommentIDIsNoResult x2 GREEN on arrival (row 60's rule, declared)
- GREEN unit 12:19:42Z..12:19:47Z exit=0; integration 12:19:47Z exit=1 (my expectation listed the commit's paths in the wrong git order; test fixed, not code); 12:20:23Z..12:21:13Z exit=0 (33.6 s)
### Slice D: the route of a question and the items question-row / answer-row (internal/run question.go)
- RED 2026-10-16T12:22:59Z exit=1 (stub DefaultOwners nil: 22 assertion lines then a nil-map panic in TestAReplacedOwnerMapRoutesByItsRoles, which hid two tests); stub made `map{}`; RED 12:23:15Z..12:23:27Z exit=1: seven tests fail by assertion (22 lines). Cases that expect false/Operator pass with a stub that returns false only for the cases they do not distinguish; their mutants M-D* carry them.
- GREEN 12:23:38Z..12:23:43Z exit=0
### Slice E1: PostComment and RecordHandoff, unit (internal/run handoffpost.go, handoffpost_test.go)
- RED (stubs return errors "not built") 2026-10-16T12:26:39Z..12:26:45Z exit=1: 17 of 18 tests fail (30 assertion lines); TestAHandoffOfASessionWithoutAStartRowIsAnError green on arrival with an erroring stub (declared; its mutant M-E-start carries it).
- GREEN first run 12:28:07Z..12:28:21Z exit=1: two assertions of mine expected a comment suffix on `invalid schema` where the result refused by the block names no comment row (expectation wrong, test fixed to `invalid schema`, not the code). An earlier expectation about the edit-back copy was also set by decision before the run: the handoff names the last copy of its bytes (seen 3). GREEN 12:28:27Z..12:28:39Z exit=0 go test -count=1 ./internal/run/ (the whole unit level of the package)
### Slice F: AnswerAccepted (internal/run answer.go, answer_test.go)
- RED (stub returns "", false) 2026-10-16T12:30:51Z..12:31:00Z exit=1: the accepted case and two accepting cases (a question citing another record; a second session of the attempt) fail by assertion; the 17 not-accepted cases pass with the stub (declared; their mutants M-F* carry them). An identifier clash (`asked`, setup_test.go) stopped the first compile at vet (exit 1) and was renamed.
- GREEN 12:31:15Z..12:31:25Z exit=0
### Slice G: integration tests of the post, the copy, the handoff, the fences and the route/acceptance from the records
- HONEST ORDER: these tests were written AFTER the code of slice E1/F (the unit slices were red first). To show them red for their own reason, handoffpost.go and answer.go were replaced by stubs (not committed) and the seven tests run: 2026-10-16T12:33:30Z..12:34:07Z exit=1, all seven (9 subtests incl. 3 kinds) FAIL; one assertion of mine (`a copy exists for a post that failed`) was a false positive of the test (Start already copied the dead-man notice, `copies/501-1.md`): fixed to count copies of the task's issue (3). Real code restored byte for byte (cp of the saved files).
- GREEN 12:34:17Z..12:35:03Z exit=0 (36.2 s): the seven integration tests with the real code


## The live probe of "to confirm at the build" (`loop.md`)

2026-10-16T12:36:09Z to 12:37:14Z, as `layup-agent[bot]` on issue #360 of `pharzam/layup`. Seven bodies were posted by
`POST /issues/{n}/comments`, read back by `GET /issues/comments/{id}`, compared with the bytes sent, and deleted by
`DELETE` (each delete exit 0). All seven came back byte for byte; the table gives the length in bytes, in characters,
the verdict and the SHA-256 of what was sent.

```
p1 111 bytes 104 chars equal f9932f213bc8763162d05b8ecb4b0c66d6acea934caba5bf0c22e30f6eba5a0d
p2 196608 bytes 65536 chars equal 957009e2bb57293a0b5805dc9602bb83ec943078e006d11cc644849751e1e542
p3 196611 bytes 65537 chars equal 3dfa7e064bf0dec8ed448e8acacd13ed54dd5d53163e3c8b1a6327c5e9d42beb
a1 65536 bytes 65536 chars equal bf718b6f653bebc184e1479f1935b8da974d701b893afcf49e701f3e2f9f9c5a
a2 65537 bytes 65537 chars equal 008ffc88d3c96a9f307524eb361e47c5222a887fc45fa0c1fb8d429c5c23b430
a3 131072 bytes 131072 chars equal b44ffb72fcc259676bd80495fef1b44b808ca8f1ffe1b1706a4d7911b0e31f11
e3 262143 bytes 87381 chars equal 973b0fd630c3f5d89138b92828894ee571bb0c601aaaccc5e27bf6e89e5660a3
```

Two bodies were refused or tested at the edge: 131,072 three-byte letters (393,216 bytes) were refused with HTTP 422 and
the message `Body is too long (maximum is 65536 characters)` (no comment was made), while 65,537 characters,
131,072 bytes and 262,143 bytes were accepted. The message is therefore not the enforced rule; code keeps the documented
65,536 characters, the safe side (`loop.md`, `guardrails.md` §2).

## Mutants

Method (`mut.py`): the target file's SHA-256 is compared with `git show HEAD:<file>` before the run (a differing file is
`SKIPPED`), a backup is copied bound to that hash, exactly one occurrence of the old text is replaced (else
`not-applied`), then `gofmt`, `go build ./...` and `go vet <package>` (with the tag for the integration cases) must pass
(else `build-broken`, which is not a kill), then the tests of the package run (`-count=1`); `killed` means a named test
of the mutant's list failed by its own assertion, `failed-other` a failure of another test, `survived` none failed. The
file is restored from the backup and the hash checked equal to the hash before and to the backup after each mutant
(`RESTORE-HASH-MISMATCH` would stop the run; none occurred), and `git status` was empty after each pass.

Result: 61 mutants. Pass 1: 50 killed, 8 build-broken (a variable left unused by the edit), 3 survived. The 8 were
rerun as variants that build (pass 2/3): all killed. Of the 3 survivors, `B2` and `R1` are **equivalent** (the
removed check is repeated by the next one: `readResult` refuses a directory, and `oneSpaced` refuses an empty reason),
kept in the code as defence in depth; `H2` was a **real test gap** (the lookup of `copies.tsv` could be skipped because the
forge read found the same comment), closed by two assertions (no forge call on the second call; no second post after an
edit) and killed on rerun. The variant `C7b` first **survived** (a missing file or a link with the hash of no bytes
passed the check: a real gap in the test), a case was added, and it was killed. No mutant is left surviving but the two
equivalent ones; none was `not-applied` or `failed-other`.

### Pass 1: 61 mutants, the source bound to its hash at HEAD

| id | outcome | detail |
| -- | ------- | ------ |
| B1 | killed | TestACommentFileThatIsNotARegularFileOfOneNameIsRefused |
| B2 | survived | the name checks of "", . and .. dropped (the directory is refused by readResult anyway: expected equivalent) |
| B3 | killed | TestACommentFileThatIsNotARegularFileOfOneNameIsRefused |
| C1 | killed | TestTheKeptBytesAreTheBytesOfTheFile,TestTheKeptPathIsTheTaskTheSessionAndTheFile |
| C2 | killed | TestACommentBodyIsRefusedBeforeThePost |
| C3 | killed | TestACommentBodyIsRefusedBeforeThePost |
| C4 | killed | TestACommentBodyIsRefusedBeforeThePost |
| C5 | killed | TestACommentBodyIsRefusedBeforeThePost |
| C6 | build-broken | # github.com/pharzam/layup/internal/run |
| C7 | build-broken | # github.com/pharzam/layup/internal/run |
| C8 | killed | TestTheKeptBytesAreTheBytesOfTheFile |
| C9 | killed | TestACommentRowIsCheckedByItsHashAgainstTheResultFile |
| C10 | killed | TestACommentRowIsCheckedByItsHashAgainstTheResultFile,TestTheKeptBytesAreTheBytesOfTheFile |
| E1 | killed | TestACommentFileThatFailsItsRowIsRefusedArtifactAndNotKept,TestTheCommentRowOfASessionIsKeptWithTheResultInOneCommit |
| E2 | build-broken | # github.com/pharzam/layup/internal/run |
| E3 | build-broken | # github.com/pharzam/layup/internal/run |
| Q1 | killed | TestAQuestionThatNeedsAHumanGoesToTheOperator |
| Q2 | killed | TestALabelOutsideTheFourGoesToTheOperator,TestAReplacedOwnerMapRoutesByItsRoles |
| Q3 | killed | TestALabelOutsideTheFourGoesToTheOperator,TestAQuestionThatNeedsAHumanGoesToTheOperator |
| Q4 | killed | TestTheRouteOfEachKindWithNoHumanNeeded |
| Q5 | build-broken | # github.com/pharzam/layup/internal/run |
| Q6 | killed | TestQuestionRowNeedsARecordAndAQuestion |
| Q7 | killed | TestAnAnswerRowNeedsAnAnswer |
| Q8 | killed | TestAnAnswerRowNeedsAnAnswer |
| Q9 | killed | TestAQuestionRowIsRead,TestQuestionRowNeedsARecordAndAQuestion |
| Q10 | killed | TestAQuestionRowIsRead |
| A1 | killed | TestAnAnswerIsNotAcceptedUnlessEachConditionHolds |
| A2 | killed | TestAnAnswerIsAcceptedByTheNextAttemptThatCitesIt,TestAnAnswerIsNotAcceptedUnlessEachConditionHolds |
| A3 | killed | TestAnAnswerIsNotAcceptedUnlessEachConditionHolds |
| A4 | build-broken | # github.com/pharzam/layup/internal/run |
| A5 | killed | TestAnAnswerIsNotAcceptedUnlessEachConditionHolds |
| A6 | killed | TestAnAnswerIsNotAcceptedUnlessEachConditionHolds |
| A7 | killed | TestAnAnswerIsNotAcceptedUnlessEachConditionHolds |
| A8 | killed | TestAnAnswerIsAcceptedByTheNextAttemptThatCitesIt |
| A9 | killed | TestAnAnswerIsNotAcceptedUnlessEachConditionHolds |
| A10 | killed | TestAnAnswerIsNotAcceptedUnlessEachConditionHolds |
| A11 | killed | TestAnAnswerIsNotAcceptedUnlessEachConditionHolds |
| H1 | killed | TestACommentThatHasACopyIsNotPostedAgain |
| H2 | survived | the lookup of copies.tsv skipped |
| H3 | killed | TestACommentThatTheForgeHoldsIsCopiedAndNotPostedAgain,TestAPostWhoseAnswerWasLostIsFoundAndCopiedNotPostedTwice |
| H4 | killed | TestACommentThatTheForgeHoldsIsCopiedAndNotPostedAgain |
| H5 | killed | TestALostLeaseStopsTheComment |
| H6 | killed | TestALostLeaseStopsTheComment |
| H7 | killed | TestTheCommentIsPostedFromTheKeptBytes |
| H8 | killed | TestThePostIsFencedByTheLease |
| H9 | killed | TestAnEditedCommentIsANewCopyAndTheFirstHandoffIsInvalid |
| H10 | killed | TestAnEditedCommentIsANewCopyAndTheFirstHandoffIsInvalid |
| H11 | build-broken | # github.com/pharzam/layup/internal/run |
| H12 | killed | TestAClaimOfTheCallerForWhatCodeComputesIsNotRead,TestTheQuestionAndAnswerHandoffsAreComputedFromTheRecords |
| H13 | killed | TestAClaimOfTheCallerForWhatCodeComputesIsNotRead |
| H14 | build-broken | # github.com/pharzam/layup/internal/run |
| H15 | killed | TestAFileThatCannotBeACommentIsRefusedBeforeThePost |
| H16 | killed | TestAHandoffOfASessionWithoutAStartRowIsAnError |
| H17 | killed | TestAnEditedCommentIsANewCopyAndTheFirstHandoffIsInvalid |
| H18 | killed | TestAFileThatCannotBeACommentIsRefusedBeforeThePost |
| H19 | killed | TestAnItemThatTheCallerDoesNotPassIsFalse,TestTheProgramItemsOfTheCallerAreHandedToTheCheck |
| H20 | killed | TestACommentThatWasNeverPostedIsInvalidArtifactsWithNoCopy |
| R1 | survived | the empty reason of an answer tested twice (expected equivalent) |
| R2 | killed | TestResultRefusesEachBrokenRule |
| R3 | killed | TestResultRefusesEachBrokenRule |
| R4 | killed | TestAValidResultIsRead,TestResultRefusesEachBrokenRule |

### Pass 2: the build-broken cases as variants that build, and H2 after its test was added

| id | outcome | detail |
| -- | ------- | ------ |
| H2 | killed | TestACommentThatHasACopyIsNotPostedAgain |
| C6b | killed | TestACommentRowIsCheckedByItsHashAgainstTheResultFile |
| C7b | survived | the read error not tested (variant that builds) |
| E2b | build-broken | # github.com/pharzam/layup/internal/run |
| E3b | killed | TestACommentFileThatFailsItsRowIsRefusedArtifactAndNotKept |
| Q5b | killed | TestQuestionRowNeedsARecordAndAQuestion |
| A4b | killed | TestAnAnswerIsNotAcceptedUnlessEachConditionHolds |
| H11b | killed | TestAHandoffIsNotRecordedTwice |
| H14b | killed | TestAFileThatCannotBeACommentIsRefusedBeforeThePost |

### Pass 3: C7b after its test was added, and a second variant of E2

| id | outcome | detail |
| -- | ------- | ------ |
| C7b | killed | TestACommentRowIsCheckedByItsHashAgainstTheResultFile |
| E2c | killed | TestTheCommentRowOfASessionIsKeptWithTheResultInOneCommit |

## Controls and checks on the clean tree

| Check | Start | End | Exit |
| ----- | ----- | --- | ---- |
| `go build ./...` | 12:55:07Z | 12:55:09Z | 0 |
| `go vet ./...` | 12:55:09Z | 12:55:09Z | 0 |
| `go test -count=1 ./...` (the unit level, 21 packages) | 12:55:09Z | 12:55:18Z | 0 |
| `go test -count=1 -tags=integration` of the ten new integration tests of `internal/run`, by exact name | 12:55:23Z | 12:56:22Z | 0 |

The full module integration level and the unfiltered `internal/run` integration level were **not** run here (the root runs
them once on the stable head).

| `sh docs/adr/adr-lint.sh` | 12:56:44Z | 12:56:45Z | 0 |
| `sh docs/prd/prd-lint.sh` | 12:56:45Z | 12:56:45Z | 0 |
| `sh docs/links/link-lint.sh` (2,737 links) | 12:56:45Z | 12:57:02Z | 0 |
| `sh docs/tests/run-discipline-tests.sh` (138 passed) | 12:57:02Z | 12:57:09Z | 0 |
| `sh docs/setup/setup-check.sh` | 12:57:09Z | 12:57:40Z | 0 |
| `git diff --cached --check` | | | 0 |
| `go test -count=1 -tags=integration -run TestPackageRules ./cmd/layup/` (after `packages.md`) | 12:40:41Z | 12:40:43Z | 0 |
| `go test -count=1 -tags=integration ./internal/tsv/ ./internal/records/ ./internal/session/ ./internal/handoff/` (the block equality) | 12:40:43Z | 12:40:54Z | 0 |

## Coordinator full integration, main merge and actual resource reports

Full module integration on abd7eadee4ae47a954da8023a7f00da7fced6e18:
`go test -count=1 -tags=integration -timeout 25m ./...` direct exit0,
all21packagesok. Wrapper UTC start2026-10-16T13:02:01.654318Z,
end2026-10-16T13:19:21.769262Z (1040114.944ms).
Current main cfcfc1c was then merged at4f05c2b. No task-authored Go changes
follow the full run; main adds the separately verified #356/#359 implementations.
Combined local checks and full-SHA CI are required on the merged landing tree.
The first independent code review must explicitly assess four documentation
conflicts: guardrails retains both sets of lessons; plan retains #359's real-git
Test cell and #360's real-git Test cell; PRD changelog retains #356 then #360;
packages retains #356 setup/migration ownership and #360's answer record kind.
No shared-index review exception is claimed for these resolutions.

Actual Sonnet5.5 CLI reports, effort high, no other model reported: plan
I120/O64449/cache creation224153/read7495702, duration594914ms; build increment
I244/O149625/cache creation231115/read45235778, duration2815648ms.
Thinking39652/46664 respectively is already included in output, not added again.
The end-of-build cumulative counters are not added to the plan a second time.
The author's remaining-context counter is not token usage and is not credited.
Independent SWE-2 Max plan review export: prompt3701351/completion47906/
cached3543055 over51steps; cached is inside prompt, not an additional input.
All exported step models are swe-2-max; elapsed is not reported.
Private exports stay outside git. Metadata audits are issue360-resource-metadata.json
and issue360-plan-review-resource.json under /tmp/layup-m2e-resume/.

## Round 1 of #360: finding 1 (the constructor of the run did not set `Sessions.App`)

All times UTC from `date -u` taken by the runner after the command; the exit code is the command's own.

| Step | Start | End | Exit | Result |
| ---- | ----- | --- | ---- | ------ |
| RED: `go test -count=1 -tags=integration -run '^TestTheSessionsOfTheRunAdoptTheCommentOfALostCopyAndPostOnce$' ./internal/run/` before the fix | 13:53:21Z | 13:53:31Z | 1 | fails by its assertion: `the restart: {ID:603 Seen:1 ...}, <nil>, 2 comments on the task's issue; want the comment adopted and copied, none posted again` |
| The one line `App: r.cfg.Register.AppSlug` in `state.sessions()` (`internal/run/meaning.go`) | | | | |
| GREEN: the same test and `TestAStopBetweenThePostAndTheCopyDoesNotPostTwice` | 13:53:37Z | 13:53:50Z | 0 | pass |
| Mutant `F1m` (`mut.py`: the assignment removed from `meaning.go`; `gofmt`, `go build ./...` and `go vet -tags=integration ./internal/run/` pass before the run) | 13:54:45Z | 13:54:50Z | 0 (runner) | **killed** by `TestTheSessionsOfTheRunAdoptTheCommentOfALostCopyAndPostOnce`; the file restored, `git hash-object` equals the blob of `HEAD`, `git status` empty |

The test builds a `state` from the world's real `Config` (`Register.AppSlug` `layup-agent`, the adapter, the `httptest`
forge) and the real records store, takes its `Sessions` from the production `state.sessions()`, runs the stop between the
post and the copy (the first `PostComment`'s copy commit fails) and a restart (a second `state.sessions()`), and asserts one
comment on the task's issue and one copy. The older tests that set `App` by hand stay (they test `PostComment`); this one
tests the constructor. The six notes of the round are not fixed here.

Checks after the fix (13:55:10Z to 13:56:15Z, all exit 0): `adr-lint`, `prd-lint`, `link-lint`, `run-discipline-tests`, `setup-check`,
`git diff --cached --check`, `go build ./...`, `go vet ./...`, `go test -count=1 ./...` (the unit level, 21 packages). No
integration command beyond the two named tests above was run.

Root full-module integration after the constructor fix: `go test -count=1
-tags=integration -timeout 25m ./...` at 7e10a292c9f20cfecff6183b4a0209ff02cdd1ab,
2026-10-16T14:00:52.928001Z–14:12:13.012882Z, direct exit0, all21 packages pass.
Elapsed680084.881ms; the wrapper captured the process return code directly.
Sonnet5.5 fix invocation reported an initial177ms acknowledgment with unchanged
counters, then278052ms for the actual fix turn. Fix increment I24/O8194/cache
creation16821/read5730253; thinking1440 is included in output. End cumulative
I388/O222268/cache creation472089/read58461733 counts plan, build and fix once.
No missing elapsed or unreported usage is inferred from context counters.

Before round2, merge current main f86f71c (the reviewed #358 landing). Five
conflicting files are resolved explicitly: guardrails keeps both sets of lessons;
PRD REQ005/NFR001 cells retain both sets of task/test citations; loop.md keeps
#360 artifact/answer acceptance tests and #358 schema refusal/replacement tests;
packages.md keeps both handoff and parameter ownership; traceability keeps both
sets of rows. Round2 must judge these resolutions and their code neighbors.
Only the PRD/traceability rows claim the named shared-index exception; other
conflicts receive scoped review. New main Go code follows the earlier full run,
so combined final full-SHA CI remains required.
SWE-2 Max round1 retained metrics: prompt1426076/completion18044/cached1301548
(cached inside prompt),17 unique metric rows. The completion export is a subset
of the corrected export, counted once. Earlier lost history and elapsed are not
reported; the timeout bound is not elapsed. Metadata: issue360-r1-resource-union-audit.json.

The merged unit test package initially failed vet: sum was declared by both the
parameter byte-slice helper and the handoff string helper. Rename only the
handoff helper and its callers to commentSum; no assertion or product code is
changed. The failed hook is not credited as a pass; the repeated hook must pass.
