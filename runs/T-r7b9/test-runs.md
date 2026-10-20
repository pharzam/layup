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
