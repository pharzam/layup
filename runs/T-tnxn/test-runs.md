# T-tnxn - the test runs of row 66aa

## S0 - the reconciliation of the plan (6262334701) with the landed code

Read at `origin/main` = `7919c71` (the branch had no work, so the merge of `origin/main` was a fast-forward). Each name of the plan was
checked in the tree, not in a plan or a review (note 4 of 6262459454). Every After edge of the row is on main: 66f, 66r (#415, `652db13`).

| Plan name | Landed form | Note |
| --------- | ----------- | ---- |
| 66e: `Finding`, `GateFindings`, `NextAttemptInput`, `NextAttempt`, `ErrNoFailure` | `internal/run/attempt.go:19-71`, as planned | `NextAttemptInput{Task, Base; Table}` is unchanged on main |
| 66f: `SpecTable`, `SpecStatus`, `specComplete` | `specstatus.go:99`, `:51`, `:33` | `SpecTable` rows are `(rule, result, detail)` and always complete: the reader refuses any other table |
| 66f: `RulesAt`, `RulesStatus`, `RulesInput`, `registerRevOf` | `statuses.go:155`, `:96`, `:113`, `:123` | as planned |
| 66f: `HeadInput`, `PostStatuses`, `StatusPost`, the order `statusContexts` | `statuses.go:176`, `:202`, `:190`, `:29` | as planned |
| 66f: `GatesTable` | `gates.go:167` | as planned |
| 66v: `RepairTask`, "each driver calls first" | `repair.go:30`; `AdvanceShapeTask` calls it at entry (`shapeloop.go:462`) | Q10 holds: nothing to raise on #415 |
| 66a: `taskRow`, the kind enum | `handoffpost.go:72`; `register.go:37` | as planned; the kind `batch` is in the enum |
| `internal/rules`: `Result`, `Success`/`Failure`/`Error` | `basic.go:25-52` | as planned |
| 66r: `AdvanceShapeTask`, `ShapeTaskInput`, `Advance{Task, Verified, Head, Wait, Stop}` | `shapeloop.go:27-41`, `:460` | `Advance` has the five fields planned |
| 66r: `loopPosition` | `shapeloop.go:255` | landed, and **has no position for the statuses of a bound head**: it gives `posShip` while no verifier of the attempt has the bound head as its base |
| 66r: the decision after the statuses | **`(*Sessions).ship`**, `shapeloop.go:649` | it posts the four statuses (`s.head`) and judges them in a loop: an `error`, a failure of `layup/spec` or `layup/rules`, and `not-active` of `layup/verify` are stops; a gates failure calls `NextAttempt` with `GatesTable`; the waiting specification is accepted by `waitingSpecAccepted` before the verifier |
| 66r: `findingsOf` | `shapeloop.go:238`, returns `(found, gateHead)`; `runStage` rebuilds the gate findings from `GatesTable(gateHead)` | it never uses an older record (`s[sessionAttempt] == prev && s[sessionBase] == gateHead`): Q10 holds |
| 66r: the row "the four contexts not posted at the bound head" | none as a position: `ship` re-runs `s.head`, which reads the forge before it posts | a restart re-derives the same decision |
| 66r: `drawTaskID`, `openShapeTask`, the scratch base `layup/base/<task>/<attempt>` | `shapetask.go`; `startSession` (`shapeloop.go:518`) | the new attempt's sessions clone the scratch branch of the base |
| 66k: its act calls `AdvanceShapeTask` | not landed (66k is not on main) | the `Restart`-level composition is 66h's and 66k's (Q8) |

### The divergence from the plan, and what it changes

The plan (D6, S4, U8) put the new decision in a "position next attempt" that replaces 66r's stop, and the driver's "status row". The landed form has no
such position: the single function that decides what follows the statuses of a bound head is `ship`, which the plan allowed for ("K6's change goes
into whatever single function decides what follows the statuses of a bound head"). The consequences, none of which changes the scope:

- D6 is stated as "the act `ship` re-derives the decision on every call", as it already does for the gates (`loop.md`, The driver of a Shape task). The restart
  property (I3) holds by the same mechanism: a stop between the post and `StartAttempt` leaves the attempt the latest and no verifier of its head.
- U8 (`TestAFailedSpecOrRulesCheckIsTheNextAttemptPosition`) tests the judgement of the posted statuses, extracted from `ship` as `judgeStatuses`
  (a genuine extraction: the loop over the posts, behaviour kept), because `loopPosition` is not where the decision lies. M14 ("without A is the
  latest") has no code to mutate; the property is held by `loopPosition` (the attempt it reads is the latest) and is tested by I3's third run.
- U9 (`findingsOf`) cannot hold the checks: `findingsOf` is pure and gives the gate head. The checks of the previous head are joined in `runStage` by
  `headFindings`, extracted from it (the gate-table part, behaviour kept). It needs real git for the rules, so its test is I1's prompt file.
- The existing `TestAFailedStatusOfTheHeadStopsTheStepWithNoVerifierAndNoNewAttempt` asserted the stops of 66r for a failed spec or rules. Those two cases
  are the behaviour this row changes, and move to the new tests; the error cases and the control stay.
- 66k's After cell: the R10 commit `006738d` of `T-9nbf` already names 66aa; this change does not touch the cell (no edit of 66k's row).

## S3 - the red runs

Every command was run directly in the worktree, never piped for its status: the exit status is the one `go` gave (`exit=`), and the start and
end are `date -u` read before and after the run. The stubs are the head `c4153a5` (inert `SpecFindings`, `RulesFindings`, `CheckFindings`,
`checkResultOf`, `expectedFailure` and `RetryGateFindings`, `NextAttemptInput.Kind` and `.Findings` ignored; `judgeStatuses` and `headFindings`
extracted with the behaviour of #415 kept), plus the tests of this row, uncommitted. The long lines are cut at 330 characters; nothing else is changed.

### Red, the unit level

Command: `go test -count=1 ./internal/run/` (start 2026-10-20T05:16:24Z, end 2026-10-20T05:16:46Z, exit=1)

First failing assertion of each test:

| Test | First failing assertion |
| ---- | ----------------------- |
| U1 `TestTheSpecFindingsAreTheFailRowsInTheOrderOfTheRules` | `attempt_test.go:94: SpecFindings = []; want [{Rule:spec:needs ... Text:detail of needs} {Rule:spec:ran ... Text:detail of ran}] (the rule ran last)` |
| U2 `TestTheRulesFindingIsItsReasonWithItsPath` | `attempt_test.go:103: RulesFindings = []; want [{Rule:rules:rule path File:.github/workflows/ci.yml Line:0 Text:rule path: .github/workflows/ci.yml}]` |
| U4 `TestTheNextAttemptTakesTheGateFindingsThenTheGivenOnes` | `attempt_test.go:128: NextAttempt = {N:1 Findings:[{Rule:gate:exits ...}]}, <nil>; want attempt 1 with [gate:exits spec:criteria rules:rule path]` |
| U5 `TestAnExpectedFailureOfABatchTaskIsNoFindingAndAnActiveRowStaysOne` | `attempt_test.go:160: batch: a pending row that failed gives [{Rule:gate:docs ...}]; want 0 finding`; its control `architecture: the rules failure gives {Findings:[] Stop:}; want 1 finding` |
| U6 `TestACheckThatCouldNotRunStopsBeforeAnyFinding` | `attempt_test.go:198: both failed: {Findings:[] Stop:}; want the stop "" and 2 findings` (and `a rules error beside a spec failure: ...; want the stop "layup/rules is error: register" and 0 findings`) |
| U8 `TestAFailedSpecOrRulesStatusIsACheckToRetryAndAnErrorStopsFirst` | `shapeloop_test.go:264: layup/spec failed: judgeStatuses = "T-sp01: layup/spec failed: criteria: x", gates false, checks false; want "", false, true` |

Green at the stub, each beside its red control: `TestAPassClearSuccessOrErrorGivesNoFinding` (U3, a green negative; red controls U1 and U2); in U8 the
cases `all pass`, `layup/gates failed` and `layup/gates error beside a spec failure` are green negatives (the gates context is judged first), and the
cases `layup/spec failed`, `layup/rules failed`, `gates and spec failed`, `layup/rules error after a spec failure` and `not-active after a rules failure` are
their red controls; `TestTheRulesResultOfAPushedHeadStaysThatOfItsBoundSessionWhenTheRecordsMoveOn` (U7) is the green characterization of the commit
`5347d5a`, kept green over the extraction `494e591`.

### Red, the integration level

Command: `TMPDIR=/home/layup/.cache/layup-session/tmp-green go test -count=1 -tags=integration -timeout 30m ./internal/run/ -run 'TestASpecFailureDuringShape|TestARulesFailureStartsANewAttempt|TestARestartAfterTheFailedStatuses|TestTheSpecificationThatWaitsForItsBet|TestAStatusErrorBesideASpec|TestCheckFindingsReads|TestAnErrorStatusOfTheHead|TestTheRulesResultOfAPushedHead'`
(start 2026-10-20T05:18:17Z, end 2026-10-20T05:21:20Z, exit=1)

| Test | First failing assertion |
| ---- | ----------------------- |
| I1 `TestASpecFailureDuringShapeStartsANewAttemptWithItsFindings` | `checkretry_integration_test.go:56: AdvanceShapeTask(T-g7lj) = {... Stop:T-g7lj: layup/spec failed: criteria: the requirement REQ-009 has no criterion}, <nil>; want the verified head` |
| I2 `TestARulesFailureStartsANewAttemptWithItsPath` | `checkretry_integration_test.go:84: AdvanceShapeTask = {... Stop:T-wui2: layup/rules failed: rule path: docs/tasks/T-wui2.md}, <nil>; want the verified head after the second attempt` |
| I3 `TestARestartAfterTheFailedStatusesStartsExactlyOneAttempt` | `checkretry_integration_test.go:118: call 1 = <nil>; want the refusal of the attempt commit` (the stop of #415 starts no attempt, so nothing is refused) |
| `TestAStatusErrorBesideASpecFailureStopsAndStartsNoAttempt` (the rules case) | `checkretry_integration_test.go:191: layup/rules error: {... Stop:T-u7nb: layup/spec failed: criteria: the requirement REQ-009 has no criterion}, <nil>; want a stop holding ": layup/rules is error: "` |
| `TestCheckFindingsReadsTheSpecTableTheRulesResultAndTheKindOfTheTask` | `statuses_integration_test.go:251: CheckFindings of a head with no row = {Findings:[] Stop:}, <nil>; want the stop for the spec status error` (then the three later assertions fail the same way) |

Green negatives at the stub: the gates case of `TestAStatusErrorBesideASpecFailureStopsAndStartsNoAttempt` (the gates are judged first; its red control is the rules case),
`TestAnErrorStatusOfTheHeadStopsTheStepWithNoVerifierAndNoNewAttempt` (the control and the gates error of #415, kept) and U7.

I4 `TestTheSpecificationThatWaitsForItsBetStartsNoAttempt` was in that run with a wrong assertion of mine (`waitingSpecAccepted` gives no refused line for an accepted
head, and I had asked for lines); it is corrected and run alone at the same stubs, **green, a negative** (its red control is I1): `go test -count=1 -tags=integration -timeout 30m
./internal/run/ -run 'TestTheSpecificationThatWaitsForItsBet'`, start 2026-10-20T05:21:42Z, end 2026-10-20T05:22:34Z, exit=0.

## S4 - green

The functions of the plan (`SpecFindings`, `RulesFindings`, `expectedFailure`, `RetryGateFindings`, `checkResultOf`, `CheckFindings`, `NextAttempt` with the
given findings and the kind), then the driver: `judgeStatuses` gives a failure of `layup/spec` or `layup/rules` as a check to retry with and an error first,
`ship` calls `CheckFindings` and `NextAttempt{Table, Findings}`, and `headFindings` joins the gate table and the checks of the failed head for the next developer's prompt file.

- `go test -count=1 ./internal/run/` (unit): `ok`, exit=0 (the same command as the red run).
- The integration command of the red run, unchanged: start 2026-10-20T05:24:16Z, end 2026-10-20T05:27:53Z, exit=0, `ok github.com/pharzam/layup/internal/run 184.094s`.
