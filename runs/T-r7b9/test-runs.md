# Test runs of T-r7b9 (#426)

Every run is a direct command with its own exit status, never piped. Times are read from `date -u` after the step.

## S0, the reconciliation with what landed (main at 7919c71, merged into this branch by fast-forward)

Plan 6262327790 and review 6262463499 cite 66r's names as planned contracts. Each is checked against main (`internal/run`):

| Planned (plan 6262327790, 66r's plan 6250184301) | Landed form | Where | Effect on this child |
| --- | --- | --- | --- |
| `(*Sessions).AdvanceShapeTask`, `ShapeTaskInput{Task, Base; Records; Harnesses, Models; Options}` | as planned (`Records func() string`, `Options map[string]spec.Result`) | `shapeloop.go:27`, `:460` | none |
| `Advance{Task; Verified; Wait, Stop}` | `Advance{Task, Verified, Head, Wait, Stop}`; **no progress field** | `shapeloop.go:36` | the hook is not a function that returns an `Advance` (see D1) |
| `loopPosition` (pure) | `loopPosition(v loopView) loopStep`, positions `posAttempt` .. `posStop`; `loopView` holds the events, the rows of `sessions.tsv`, `Status` and the parsed review `Records` | `shapeloop.go:255`, `:72` | the view is widened for the question stages (D2) |
| `readOnlyStep` | `readOnlyStep(step)`, three steps; used by `TaskSession` (`session.go:343`) | `shapeloop.go:762` | gains `answer-to-a-question` (U6) |
| the stop `<task>: the session <id> of <step> ended <status>` | as planned, from `loopView.outcome`: any `done` result whose status is not `completed`, `needs_context` included | `shapeloop.go:166` | the stop for `needs_context` is where the hook goes (D1) |
| completed-stage recovery (`finishStage`) | as planned: post and handoff both read before they write | `shapeloop.go:591` | reused as the model of the question stages |
| `RepairTask` first | `AdvanceShapeTask` calls it first | `shapeloop.go:462`, `repair.go:30` | preserved: the question position is after it |
| the test seam | `drvWorld` and the scripted harness (`shapetask_integration_test.go`), the driver tests call `state.sessions()` | `shapetask_integration_test.go:44` | reused for I1 to I4 |
| `findingsOf` / an existing review-record selector (S0 case in the plan) | `findingsOf` exists but is attempt-bound and reads `loopView.Records` of the Shape driver; `Position` has no register row | `shapeloop.go:238` | **not** reused for `openMaterial` (D4) |
| the tier selection of the driver (review condition 5: "extend I5 through the driver's actual tier selection") | **none**: `runStage` and `ship` take the tier from the step table (`ResolveStep(step, "", "")`), and every Shape step is fixed `reasoning`; `Position` and `TierOf` have no caller in non-test code | `shapeloop.go:555`, `tier.go:50`, `:105` | I5 goes through the production `Position` and `TierOf` (D5) |

The other predecessors are as the plan read them (66m `answer.go:68`, `:98`; 66o `qacomment.go:25`, `:37`; 66l `commentform.go:198`, `:257`;
66d `question.go`, `handoffpost.go:143`, `:257`, `:418`; 66s `instructions.go`, `promptfile.go:301`, `:429`; 66v `repair.go:30`; 66e `session.go:187`, `:219`, `authors.go:130`, `tier.go:50`, `:163`).
The rows 66r, 66s and 66v are on main; row 66z is not (the R10 commit `006738d` of `T-9nbf` is not merged), so this branch adds row 66z.

**Divergences from the plan, each inside the plan's scope (a different form of the same behaviour), none a new goal:**

- **D1. The hook is a position and an act of the driver's loop, not `AdvanceQuestion(ctx, in) (Advance, error)`.** 66r's driver is a loop of one act and a fresh read
  (`maxActs` 100); `Advance` has no progress field. So the hook is the position `posQuestion` with `askStage(ctx, in, v, st) (stop, wait string, err error)`:
  an act is progress (the loop reads again), no applicable work is not a position, a wait is `Advance{Wait}`, a stop is `Advance{Stop}`. Review condition 2 (report an act, no work,
  a wait and a stop) is met by that shape; `RepairTask` stays the first call (`shapeloop.go:462`). The statuses `blocked`, `failed` and `decision_needed` keep the stop.
- **D2. `loopView` gains `Results`, `Copies` and `Issue`** (the typed rows of each `done` result, the copies of the task's issue), which `shapeView` already reads or can read;
  `questionPosition` is a pure function of the view, so U1 to U4 and U7 are unit tests.
- **D3. The plan stage and the plan review have no attempt filter in 66r's `loopPosition`** (`plans`, `reviews` are the sessions of the role over all attempts), so a plan that asked would stay
  the last session of its stage and stop the loop after the answer. The sessions that ended `needs_context` in an earlier attempt are therefore skipped from those lists, and the plan review is told apart
  from a verifier by a **comment** handoff only (a question or answer handoff of a plan-role session, whose detail has no `comment`, is no plan handoff). Review condition 3.
- **D4. `openMaterial` is its own selector over events, `sessions.tsv`, the typed results and the kept copy** (`Position` has no register row to build a `loopView`), with the qualification of condition 5.
- **D5. No tier selection exists in the Shape driver** (last row of the table). I5 is the production `Sessions.Position` (from `state.sessions()`) followed by `TierOf`; there is no caller to wire, and the spec text says so as a known limit.

## S2, the inert stubs and the extractions (commit 103a318), all green

The stubs: `questionPosition` returns no stage, `askStage` does nothing, `openMaterial` returns false, the two instruction constants are empty and in no map, `readOnlyStep` is
unchanged, nothing calls them; `postBody` is the real extraction of `PostComment` (the issue lookup to the read-back, `handoffpost.go`), `startAs` is `startSession` with the role given.

- `go build ./...` exit 0; `go vet ./...` exit 0; `go vet -tags=integration ./...` exit 0; `go vet -tags=e2e ./...` exit 0.
- `go test -count=1 ./...`: start 2026-10-20T05:04:30Z, end 2026-10-20T05:04:48Z, exit 0.
- `TMPDIR=/home/layup/.cache/layup-session/tmp-green go test -count=1 -tags=integration -timeout 30m ./internal/run/ -run 'Comment|Handoff|Position|Question'`:
  start 2026-10-20T05:04:48Z, end 2026-10-20T05:13:48Z, exit 0 (`internal/run` 498.891 s; the host was loaded by two other agents' runs). These are the characterization tests of the extraction
  (`handoffpost_test.go`, `handoffpost_integration_test.go`, `qacomment_*`, the position tests): green before and after it.

## S3, the semantic red run against the stubs (commit 103a318 plus the new and changed tests, uncommitted then)

The tests: `internal/run/questionloop_test.go` (new), `internal/run/questionloop_integration_test.go` (new), `internal/run/testdata/prompts/answer-to-a-question.txt` (new golden file, built
from the block of `loop.md`), and the changed `instructions_test.go`, `shapeloop_test.go` (the one case of `needs_context`) and `promptfile_integration_test.go` (the refusal list loses
`answer-to-a-question`; `TestTheInstructionBlocksOfLoopMd` also reads every `instruction-*` block of `loop.md`, review condition 6).

### Unit, run 1: start 2026-10-20T05:11:18Z, end 2026-10-20T05:11:39Z, exit 1

`go test -count=1 ./internal/run/ -run 'TestTheQuestionPositionIs|TestAnInvalidRecordedHandoff|TestAHumanQuestion|TestQuestionsOfTwoOwnerRoles|TestTheNthAnswer|TestAnOwnerSessionThatDoesNotAnswer|TestAnAskingSessionIs|TestTheStageThatAsked|TestTheAnswerStepIsReadOnly|TestTheLastQualifiedReview|TestTheAnswersToCite|TestThePromptFileOfEachStepIsItsGoldenFile|TestNoPlaceholder|TestAStageThatEndedIsFinished'`

```
--- FAIL: TestThePromptFileOfEachStepIsItsGoldenFile (0.01s)
    --- FAIL: TestThePromptFileOfEachStepIsItsGoldenFile/answer-to-a-question (0.00s)
        instructions_test.go:86: buildPrompt(answer-to-a-question): the step answer-to-a-question has no instruction in M2e
--- FAIL: TestNoPlaceholderIsLeftInAnyPrompt (0.00s)
    instructions_test.go:103: buildPrompt(answer-to-a-question): the step answer-to-a-question has no instruction in M2e
--- FAIL: TestTheAnswersToCiteAreWrittenOnlyWhenTheLastAskingAttemptHasThem (0.00s)
    instructions_test.go:590: the prompt with answers lacks "Their IDs are S-dddddddd.1 S-dddddddd.2."
    instructions_test.go:590: the prompt with answers lacks "question row"
--- FAIL: TestTheQuestionPositionIsTheFirstStageTheRecordsShowNotDone (0.00s)
    questionloop_test.go:158: after 0 facts: stage 0, session "", owner "", base "", question ""; want stage 1, "S-0a1b2c3d", "", "", ""
    questionloop_test.go:158: after 1 facts: stage 0, session "", owner "", base "", question ""; want stage 2, "S-0a1b2c3d", "", "", ""
    questionloop_test.go:158: after 2 facts: stage 0, session "", owner "", base "", question ""; want stage 4, "", "software-architect", "aaaaaaa", ""
    questionloop_test.go:158: after 3 facts: stage 0, session "", owner "", base "", question ""; want stage 4, "", "software-architect", "aaaaaaa", ""
    questionloop_test.go:158: after 4 facts: stage 0, session "", owner "", base "", question ""; want stage 5, "S-2c3d4e5f", "", "", "S-0a1b2c3d.1"
    questionloop_test.go:158: after 5 facts: stage 0, session "", owner "", base "", question ""; want stage 6, "S-2c3d4e5f", "", "", ""
    questionloop_test.go:158: after 6 facts: stage 0, session "", owner "", base "", question ""; want stage 7, "", "", "", ""
    questionloop_test.go:166: the row of the stage P2 = []; want the question row
--- FAIL: TestAnInvalidRecordedHandoffStaysAStop (0.00s)
    questionloop_test.go:183: the question's handoff, read 1: stage 0, ""; want the stop "T-sp01: the handoff asking role → owner role is invalid question-row"
    questionloop_test.go:183: the question's handoff, read 2: stage 0, ""; want the stop "T-sp01: the handoff asking role → owner role is invalid question-row"
    questionloop_test.go:183: the answer's handoff, read 1: stage 0, ""; want the stop "T-sp01: the handoff owner role → asking role is invalid answer-row"
    questionloop_test.go:183: the answer's handoff, read 2: stage 0, ""; want the stop "T-sp01: the handoff owner role → asking role is invalid answer-row"
--- FAIL: TestAHumanQuestionIsAWaitWithNoAnswerSession (0.00s)
    questionloop_test.go:194: a human question before its post: stage 0; want the post (P2)
    questionloop_test.go:198: a human question: stage 0, "", owner ""; want the wait "waiting for the Operator: T-sp01: question S-0a1b2c3d.1"
    questionloop_test.go:201: the control, a question of an owner role: stage 0; want the answer session
--- FAIL: TestQuestionsOfTwoOwnerRolesAreAWait (0.00s)
    questionloop_test.go:212: two owner roles: stage 0, ""; want the wait "waiting for the Operator: T-sp01: question S-0a1b2c3d.1"
    questionloop_test.go:216: the control, two questions of one owner role: stage 0, owner ""; want the answer session
--- FAIL: TestTheNthAnswerNamesTheNthQuestionAndAMissingOrExtraAnswerStops (0.00s)
    questionloop_test.go:235: two answers, none posted: stage 0, question "", row []; want the post of the answer 1 for S-0a1b2c3d.1
    questionloop_test.go:239: the first answer posted: stage 0, question ""; want the post of the answer 2 for S-1b2c3d4e.1 (two questions of n 1)
    questionloop_test.go:247: one answer for two questions: stage 0, ""; want the stop "T-sp01: the question S-1b2c3d4e.1 has no answer"
    questionloop_test.go:247: no answer row: stage 0, ""; want the stop "T-sp01: the question S-0a1b2c3d.1 has no answer"
    questionloop_test.go:247: three answers for two questions: stage 0, ""; want the stop "T-sp01: the answer S-2c3d4e5f.3 answers no question"
--- FAIL: TestAnOwnerSessionThatDoesNotAnswerStops (0.01s)
    questionloop_test.go:271: failed: stage 0, ""; want the stop "T-sp01: the session S-2c3d4e5f of answer-to-a-question ended wall"
    questionloop_test.go:271: status failed: stage 0, ""; want the stop "T-sp01: the session S-2c3d4e5f of answer-to-a-question ended failed"
    questionloop_test.go:271: asks again: stage 0, ""; want the stop "T-sp01: the session S-2c3d4e5f of answer-to-a-question ended needs_context"
    questionloop_test.go:271: a wrong role: stage 0, ""; want the stop "T-sp01: the session S-2c3d4e5f started after the question is of the role software-developer, not of the owner role software-architect"
--- FAIL: TestAnAskingSessionIsThePositionQuestionAndOtherStatusesStayStops (0.00s)
    questionloop_test.go:280: a plan that asked: position 13 ("T-sp01: the session S-0a1b2c3d of the plan ended needs_context"), stage 0; want the question at the handoff
--- FAIL: TestTheStageThatAskedRunsAgainInTheNextAttempt (0.00s)
    questionloop_test.go:295: a plan that asked, attempt 2: position 13 ("T-sp01: the session S-p1 of the plan ended needs_context"); want the plan again
    questionloop_test.go:298: a plan review that asked, attempt 2: position 13 ("T-sp01: the session S-r1 of the plan review ended needs_context"); want the plan review again
    questionloop_test.go:307: an owner handoff of the plan role: position 13 ("T-sp01: the session S-r1 is neither the plan review nor a verifier"); want the developer, the plan review unchanged
--- FAIL: TestTheAnswerStepIsReadOnly (0.00s)
    questionloop_test.go:314: readOnlyStep(answer-to-a-question) = false; want true
--- FAIL: TestTheLastQualifiedReviewRecordIsTheSourceOfAnOpenMaterialFinding (0.00s)
    questionloop_test.go:428: line 1: no header row
--- FAIL: TestAStageThatEndedIsFinishedAndAStageThatStoppedIsRunAgain (0.01s)
    shapeloop_test.go:136: a plan whose status is needs_context: 13 "T-sp01: the session S-p1 of the plan ended needs_context"; want 14 holding ""
FAIL
FAIL	github.com/pharzam/layup/internal/run	0.197s
FAIL
```

The first failing assertion of the one run that the first run could not reach (`TestTheLastQualifiedReviewRecord...` failed on a fixture defect of mine, `line 1: no header row` for a world with no
session; the fixture was corrected and the test run again):

### Unit, run 2 (the material source): start 2026-10-20T05:11:55Z, end 2026-10-20T05:12:17Z, exit 1

`go test -count=1 ./internal/run/ -run 'TestTheLastQualifiedReview'`

```
--- FAIL: TestTheLastQualifiedReviewRecordIsTheSourceOfAnOpenMaterialFinding (0.02s)
    questionloop_test.go:432: a material record: open = false, <nil>; want true
    questionloop_test.go:432: a later material record opens an earlier clean one (the last, not the first): open = false, <nil>; want true
    questionloop_test.go:432: a completed review awaiting its handoff: open = false, <nil>; want true
    questionloop_test.go:432: an invalid handoff: open = false, <nil>; want true
    questionloop_test.go:432: a missing kept file: open = false, <nil>; want true
    questionloop_test.go:432: a hash that is not the result's: open = false, <nil>; want true
    questionloop_test.go:432: another head: open = false, <nil>; want true
    questionloop_test.go:432: a form that does not parse: open = false, <nil>; want true
    questionloop_test.go:432: an unqualified later record never lets the older clean one stand: open = false, <nil>; want true
    questionloop_test.go:439: a failed read of the kept file: open = false, <nil>; want an error
FAIL
FAIL	github.com/pharzam/layup/internal/run	0.055s
FAIL
```

### Integration, run 3: start 2026-10-20T05:12:26Z, end 2026-10-20T05:14:28Z, exit 1

`TMPDIR=/home/layup/.cache/layup-session/tmp-green go test -count=1 -tags=integration -timeout 30m ./internal/run/ -run 'TestANeedsContextEnd|TestTheAnswerPromptHolds|TestAHarnessThatFollowsTheAnswer|TestAHumanRoutedQuestionAnd|TestARestartAtEachStage|TestAnOpenMaterialFinding|TestTheInstructionBlocksOfLoopMd'`

```
--- FAIL: TestTheInstructionBlocksOfLoopMd (0.01s)
    promptfile_integration_test.go:487: 9 texts in the code; want the eight steps, the shared result text and the conditional answers text
    promptfile_integration_test.go:493: loop.md has the block instruction-answer-to-a-question and the code has no text for it
    promptfile_integration_test.go:501: loop.md has the block instruction-answers: true; it is equal to the code: false
--- FAIL: TestANeedsContextEndIsAnsweredByTheOwnerRoleAndTheNextAttemptStarts (7.10s)
    questionloop_integration_test.go:102: AdvanceShapeTask = {Task:T-a19d Verified:false Head: Wait: Stop:T-a19d: the session S-5b9d24f6 of the plan ended needs_context}, <nil>; want the verified head after the answer
--- FAIL: TestTheAnswerPromptHoldsTheQuestionsAndTheNextPromptsTheAnswerIDs (6.59s)
    questionloop_integration_test.go:102: AdvanceShapeTask = {Task:T-gx6l Verified:false Head: Wait: Stop:T-gx6l: the session S-494868dc of the plan ended needs_context}, <nil>; want the verified head after the answer
--- FAIL: TestAHarnessThatFollowsTheAnswerInstructionGivesAValidAnswerHandoff (6.80s)
    questionloop_integration_test.go:102: AdvanceShapeTask = {Task:T-6l2r Verified:false Head: Wait: Stop:T-6l2r: the session S-dfb13eb1 of the plan ended needs_context}, <nil>; want the verified head after the answer
--- FAIL: TestAHumanRoutedQuestionAndAMixedOwnerQuestionWaitForTheOperator (6.54s)
    questionloop_integration_test.go:233: a human-routed question: AdvanceShapeTask = {Task:T-wap3 Verified:false Head: Wait: Stop:T-wap3: the session S-892ff335 of the plan ended needs_context}, <nil>; want the wait "waiting for the Operator: T-wap3: question S-892ff335.1"
--- FAIL: TestARestartAtEachStageOfTheQuestionWritesEachThingOnce (51.91s)
    questionloop_integration_test.go:274: the post of the question: the first call = nil; want the interruption
    questionloop_integration_test.go:274: the copy of the question (adoption): the first call = nil; want the interruption
    questionloop_integration_test.go:274: the handoff to the owner role: the first call = nil; want the interruption
    questionloop_integration_test.go:274: the end of the answer session: the first call = nil; want the interruption
    questionloop_integration_test.go:274: the post of the answer: the first call = nil; want the interruption
    questionloop_integration_test.go:274: the copy of the answer (adoption): the first call = nil; want the interruption
    questionloop_integration_test.go:274: the handoff back to the asking role: the first call = nil; want the interruption
    questionloop_integration_test.go:274: the start of the next attempt: the first call = nil; want the interruption
--- FAIL: TestAnOpenMaterialFindingKeepsTheTaskUphill (9.98s)
    questionloop_integration_test.go:335: an open material finding: downhill, execution; want uphill, reasoning
FAIL
FAIL	github.com/pharzam/layup/internal/run	89.058s
FAIL
```

### What the reds are, and the green negatives

- Every test above fails on the behaviour it names, not on a build or a fixture error: the driver returns the 66r stop `the session <id> of the plan ended needs_context` where the demo wants the
  verified head (I1, I6, I7) or the wait (I2, I3); no interruption occurs because nothing ran past the stop (I4: its first assertion is the interruption that the stage after the stop would take, so its
  eight cases are each red for the missing act, and the green negative of the same restart is the duplicate counts, which hold only after the acts exist); the position is `downhill` with a material
  finding (I5); `loop.md` has two blocks that the code does not know (the block test); the prompt of the answer step has no instruction (`the step answer-to-a-question has no instruction in M2e`, U8).
- **Green negatives**, passing against the stubs and kept as controls beside their red tests (they pin what must not change): U1 `no asking session: no stage`; U5 `no review`, `a clean record, valid handoff`,
  `a later clean record clears an earlier material one` and `a plan review is no review record`; U6 `readOnlyStep(implement) = false`; the loop test `a developer that asked, attempt 2 = the developer`
  (the attempt filter of 66r already does it) and `blocked`, `failed`, `decision_needed` stay stops; the changed case `a plan whose status is blocked`; `TestTheAnswersToCite...` first half (the prompt with no answers is
  byte for byte the golden file of `implement`); I5 `a frozen list and no review: downhill` (the control of the position); the seven existing golden files.
- **Red controls**: U1 each of the seven stages; U2 and U3 the waits (their control, a question of one owner role, is red too: `stage 0; want the answer session`); I1 is the red control of `owner sessions = 0` of I2 and I3.

## S4, green (the code, in the order of the plan: `openMaterial` and `Position`; the prompt and the blocks; `questionPosition` and `askStage`; the hook in the driver)

- `go test -count=1 ./internal/run/` (unit), after the code: exit 0 (0.6 s). Every test of S3 passes, the golden file of `answer-to-a-question` byte for byte.
- Integration, the same command as run 3: start 2026-10-20T05:17:37Z, end 2026-10-20T05:22:59Z, exit 0 (`internal/run` 298.556 s).
- `TestTheAnswerSessionStartsAtTheBaseOfTheLastAskingSession` was written **after** the code (review condition 3 names the base; the first plan had no test for it). It is green at once; its evidence is mutant M28, which it kills.
- One defect of the first draft of the code, found by the unit test before the integration run: none (the unit tests were green at the first run); the integration tests were green at the first run.

## S7 so far, and the resume point (host move, 2026-10-20)

Recorded at the head before this section: `go build`, `go vet` (plain, integration, e2e tags), `gofmt -l`, `go test -count=1 ./...`, the six lint scripts, `git diff --check origin/main` all exit 0;
`go test -count=1 -tags=e2e ./...` start 2026-10-20T06:20:01Z, end 2026-10-20T06:23:55Z, exit 0. The first full integration run (06:16:14Z to 07:05:56Z) wrote its output to a
scratch name that another agent's run shared; its output was truncated (no `internal/run` line), so it is **not evidence**. A rerun under the unique name `scratchpad/T-r7b9/T-r7b9-full-integration.txt` was stopped by the host move: not recorded.

Resume point:
1. Next slice: S7. Run `TMPDIR=/home/layup/.cache/layup-session/tmp-green go test -count=1 -tags=integration -timeout 90m ./...` under a T-r7b9-prefixed scratch name (about 50 min), and record start, end and exit.
2. Recorded: S0, S1 (docs), S2 (stubs; green), S3 (red runs 1 to 3), S4 (green), S5 (30 mutants, 30 killed at 8c2a05b). Not recorded: the full integration level; the close-out (task file, completed.md, backlog.md, resource record) is not written, as ordered.
3. Row 66z is in docs/plan/README.md, and 66z is in row 66k's After cell (`task-state.sh` exits 0, T-r7b9 `in progress`); expect a one-line conflict with T-tnxn (66aa) there, keep both labels.
4. Then `~/.local/bin/gitleaks git --no-banner` alone, on exit 0 `sh /home/layup/.cache/layup-session/harness/2026-10-16-m2e/push.sh T-r7b9`; no pull request.
5. Budget at this head: see `git diff --shortstat origin/main...HEAD` (about 1,750 lines over 22 files, limit 2,500 over 27, close-out inside).
