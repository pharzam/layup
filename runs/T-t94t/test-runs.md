# T-t94t: test runs

Task T-t94t, issue #483: a plan review that rejects the plan sends the task back to its planner, and the plan-review instruction states the target's plan form. Times are from `date -u` after each run. Go 1.26.10 from mise. Base `ae2997a6` (main moved from `5d13785e` by #485 before the first test; merged by fast-forward).

## 1. The scaffolding

Committed first (`ed16bfe9`, green): inert stubs only. `handoff.InForce` returns its rows unchanged; `PlanReviewReasons` returns nil and `RejectReasonsOK` false; `planRejects` returns 0 and `rejectionOf` nil; the position `posReplan` is declared and never chosen; the view field `Rejects` exists. The signature of `amendmentPending` gained its third argument, which was not read yet; that inert change is in the working tree before the red run and landed in the commit `588c145b` (it is not in `ed16bfe9`, which does not touch `frozen.go`). The item `reject-reasons`, the 14th default row, the reopened lists, the stable `qaKind`, the route of `finishStage`, the instruction block, its constant and `plan-review.txt` are not in the scaffolding: they are the operative edits of the slices below. The package compiled and `go vet` (no tag, `integration`) was clean before the red runs, so no red below is a compile failure.

## 2. Red (assertion failures on the current behaviour)

Unit: `go test -count=1 ./internal/handoff ./internal/run`, exit 1, finished 2026-10-22T12:19:55Z.

| Test | First failing assertion |
| --- | --- |
| `TestTheDefaultRowsAreTheFourteenOfTheSpecification`, `TestEachDefaultPairNeedsItsOwnItemsAndKinds`, `TestWriteThenReadGivesTheDefaultTable` | `13 default rows, want 14`; the read-back table has 13 rows |
| `TestTheListsAreTheWordsOfTheSpecification` | the items end at `answer-row`; want `reject-reasons` after it |
| `TestATableWrittenBeforeALaterPairTakesItsDefaultRow` | a reject on the table of 13 rows in force: `invalid pair comment 12-1`; want valid; `InForce` added 0 rows |
| `TestASetThatLacksALaterPairIsRefused` | `CheckTable of 13 rows: <nil>`; want the missing pair named |
| `TestTheReasonsOfAPlanReview` | `one reason: []; want ["no ordered steps"]` (3 of 5 cases) |
| `TestRejectReasonsIsTrueOnlyForARejectWithAReason` | `a reject with a reason: false`; want true |
| `TestARejectingReviewIsNoReviewOfAnAmendment` | `a review that rejected is the review of an amendment` |
| `TestTheCommentItemsAreComputedFromTheKeptCopy`, `TestAKeptCopyTheStoreLacksOrCannotReadMakesTheItemsFalse`, `TestTheKeptFileThatIsNotTheNamedHashIsNoComment` (the case `reject-reasons`, absent, empty, not UTF-8 and mismatched kept bytes) | `{Valid:false Detail:invalid pair comment 701-1}`; want `valid comment 701-1` or `invalid reject-reasons comment 701-1` |
| `TestACodeWrittenHandoffHasNoCommentItem` | `transitions.tsv: line 9: column items: "reject-reasons" is not an item` |
| `TestAReviewThatRejectsIsAValidHandoffToThePlan` | all six cases of the two tables (14 and 13 rows): `invalid pair comment 701-1`; the case to the developer is green (see the controls) |
| `TestASetOfTransitionsIsJudgedByTheCopyItNames` (the case of 13 rows) | `Change:0x...`; want the item `bound` and the pair named |
| `TestThePlanReviewIsToldThePlanFormOfTheTarget` | the constant and the prompt of each of 3 kinds lack each of the four sentences (16 failures) |
| `TestARejectedPlanStartsANewAttemptWithThePlanReopened` | `a valid reject: 6 ""; want the replan of attempt 1` (the developer) |
| `TestTheRejectsAreBoundedByAnOwnCountOfTheCycleCap` | `planRejects = 0; want 2`; the stop is `the session S-r1 is neither the plan review nor a verifier` |
| `TestARejectAfterARetryReopensThePlan` | `{Pos:6 ...}; want the replan` |
| `TestTheDemoRecordsOfTUt9gGiveAttempt4` | `{Pos:13 Session:S-0ef59a25 Msg:...is invalid verdict-field comment 6339776504-1}`; want to finish the handoff |
| `TestAPlanReviewKeepsItsStageWhenALaterPlanIsHandedOff` | `stageOf(S-38157561) = ""` (and the two others); `replacedOf` is `map[]` after the new plan, `map[1:S-38157561 2:S-141cc16c]` before; the stop `S-38157561 is neither the plan review nor a verifier` |
| `TestThePlanOfTheAttemptAfterARejectIsGivenTheReasons` | `rejectionOf(2) = []; want the two reasons`; the position after the crashed plan is the review S-r1 |

Controls that are green before and after (they pin what the change must not move): `TestARefusedOrReplacedReviewWhoseCopyRejectsReopensNothingAndSuppliesNothing` (the refused `S-141cc16c` with other reasons gives the retry, then the plan review, no replan, no reasons); `TestABareRejectIsAStopAndStartsNoAttempt` (an `invalid reject-reasons` handoff is a stop); `TestInForceAddsNoOriginalPairAndKeepsACustomizedRow` (a missing original pair still fails at each of the 13 positions; a customized row 14 is kept); `TestARegistrationKeepsATableWrittenBeforeALaterPair` (the table of 13 rows keeps its bytes, a table lacking an original pair is refused, a customized row 14 is kept); the case "a reject with a reason to the developer: invalid verdict-field" of `TestAReviewThatRejectsIsAValidHandoffToThePlan`.

Integration: the new tests by name with the tag `integration`, exit 1, finished 2026-10-22T12:20:07Z. In each, the scripted `reject` gives, on the current code, the stop `the handoff plan review → developer is invalid verdict-field comment 604-1`:

| Test | First failing assertion |
| --- | --- |
| `TestARejectedPlanGoesBackToThePlannerWithTheReasons` | `AdvanceShapeTask = ... Stop:...invalid verdict-field comment 604-1; want the verified head` |
| `TestASecondRejectStopsAtTheCapOfItsOwn` | the same stop; want `the plan review rejected the plan 2 times` |
| `TestABareRejectIsAnInvalidHandoffAndStartsNoAttempt` | the stop is `invalid verdict-field`, not `invalid reject-reasons` (a bare reject goes to the developer today) |
| `TestARejectAfterACrashedPlanReviewReopensThePlan` | the same stop at attempt 2; want the verified head |
| `TestARestartBeforeOrAfterTheReplanAttemptStartsExactlyOne` | `before the attempt: AdvanceShapeTask did not fail` (no attempt 2 is ever committed) |
| `TestTheDemoRecordsOfTUt9gGiveAttempt4WithNoRepair` | on the earlier binary's records (13-row table, `invalid verdict-field comment 606-1`): the stop; want the verified head |

`TestTheDefaultRowsAreTheTableOfTheSpecification` and `TestTheListsAreThoseOfTheSpecification` (integration, `internal/handoff`) and the golden-file test pass on the current code: they compare the defaults and the constants with `loop.md`, so they fail as soon as one side is edited alone, and the operative edit and the `loop.md` edit are made in one slice.

## 3. Green

The operative edits were made in their slices before this run (`internal/handoff`: the item, the 14th row, `InForce`; `internal/records`: the item word; `internal/run`: the readers, `commentItems`, `InForce` at the two reads, the driver, the amendment reader; `docs/spec/loop.md`, the constant and `plan-review.txt`: the table, the items, the instruction block), because the integration tests compare the defaults and the instruction constants with `loop.md`. They were committed with the tests (`588c145b`), because the hook refuses a red commit. The `loop.md` edits include the old acceptance sentence of thirteen default rows (now fourteen).

- Unit, `go test -count=1 ./...`: every package ok (the hook run at the commit, 2026-10-22T12:26Z); the 11 failing groups of section 2 are green, and so are the controls.
- The new integration tests and the comparisons with `loop.md`, by name with the tag `integration`: ok, 2026-10-22T12:23:59Z.
- The full integration level, `go test -count=1 -tags=integration -timeout 30m ./...` (output in `/home/layup/.cache/layup-session/T-t94t-full-integration.txt`): every package ok, exit 0, finished 2026-10-22T12:30:31Z (`internal/run` 242.4 s, `cmd/layup` 88.4 s).
- E2E, `go test -count=1 -tags=e2e -timeout 10m ./...` (`/home/layup/.cache/layup-session/T-t94t-e2e.txt`): every package ok, exit 0, finished 2026-10-22T12:32:12Z. The scenario `cmd/layup/loop_e2e_test.go` is unchanged; it asserts each handoff `valid` on the table of 14 rows, and the reject route is proved by the integration tests above, not by the e2e contract.
- `go build ./...`, `go vet ./...` (no tag, `integration`, `e2e`), `gofmt -l .`, `git diff --check origin/main HEAD`, the six lints (`adr-lint`, `prd-lint`, `link-lint`, `run-discipline-tests` 138 passed, `nested-checkout-check`, `setup-check`) and `freeze-check/run-all.sh` (85 passed): clean.

## 4. Round 1 (GPT-6 Astra, material; findings F1 to F3)

Red, written before the code (`go test -count=1 ./internal/run`, 2026-10-22T13:04:50Z, and the integration tests by name):

| Test | First failing assertion |
| --- | --- |
| `TestReviewerUnreadReviewIsNoApprovalInThePosition` | 4 cases (missing and mismatched, before and after the new plan): position 6, the developer; want the stop `the plan review S-r1 cannot be read from its kept copy` |
| `TestReviewerValidRejectLostCopyCannotAuthorizeDeveloper` (integration) | 4 cases: `Verified:true`, one developer session; want the stop and no developer. The controls (copy intact, before and after the new plan) are green. |
| `TestReviewerLostRejectCannotApproveAmendment` | `a missing copy: map[]; want named as not approving: true`, and the mismatched copy; the intact reject and the intact approve are green |
| `TestReviewerBareRejectNamesActualPair`, `TestABareRejectIsAnInvalidHandoffAndStartsNoAttempt` | the stop reads `plan review → developer is invalid reject-reasons`; want `plan review → task plan is invalid reject-reasons` |

Green: `readPlanReview` tells a checked verdict from an unread review (the view fields `Unread` and `Verdicts`), the driver stops on a valid handoff over an unread review, the label of the stage follows the verdict, and the amendment reader names an unreadable review with the rejecting ones. Unit, the new integration tests and the neighbours: ok, 2026-10-22T13:05:23Z. The full levels are recorded below.

Round 1 gates: the full integration level (`/home/layup/.cache/layup-session/T-t94t-full-integration-2.txt`) every package ok, exit 0, 2026-10-22T13:10:14Z; e2e every package ok, 13:11:55Z; build, vet (no tag, `integration`, `e2e`), gofmt, the six lints and `freeze-check/run-all.sh` (85 passed) clean.

## 5. End 3 (round 2: F4 and F5; the Operator chose end 3, no further round)

Red, written before the code (`go test -count=1 ./internal/run` and the integration tests by name, 2026-10-22T13:42:55Z):

| Test | First failing assertion |
| --- | --- |
| `TestReviewerCustomizedRejectRowCannotApproveBareReject` | a bare reject whose handoff is valid: position 6, the developer; want the stop. The three controls (an approve, a reject with a reason, an unparseable copy) are green. |
| `TestReviewerCustomizedRejectRowCannotApproveBareRejectInTheLoop` (integration, a customized row 14 with `schema` and `artifacts`) | `Verified:true`; want the stop for a reject with no reason |
| `TestReviewerLegacyBareRejectKeepsRecordedPair`, `TestReviewerLegacyBareRejectKeepsRecordedPairInTheLoop` | the stop reads `plan review → task plan is invalid verdict-field`; want `plan review → developer` |

Green (2026-10-22T13:43:22Z): the position stops a readable reject with no reason (`Verdicts` is `reject`, no reasons) whatever the table, keeping the table as the Operator wrote it; the label is the planner pair only when the handoff is not an earlier binary's `invalid verdict-field`. Recovery is not broadened. The stale count in `commentform.go:4` is removed. The full gates follow in the commit of the close-out.

End-3 gates: the integration packages touched (`go test -count=1 -tags=integration -timeout 30m ./internal/run ./internal/handoff`, `/home/layup/.cache/layup-session/T-t94t-integration-3.txt`) ok, exit 0, 2026-10-22T13:48:28Z (`internal/run` 242.1 s); unit ok at the commit hook; build, vet, gofmt, the six lints, `freeze-check/run-all.sh` and the freeze check clean in the close-out commit.
