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
