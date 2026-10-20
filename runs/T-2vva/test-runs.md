# T-2vva - the test runs of row 66j

## S0 - the reconciliation of the plan (6250593948), its review (6250697855) and the notes with the landed code

Read at `origin/main` = `699d868` (the branch was at `699d868` with no work; the merge was a no-op). T-05f4 (#444) had not merged (`gh pr view 444`: OPEN), so the integration tests follow the pattern of main today.

| Name of the plan, the review or a note | Landed form | Note |
| -------------------------------------- | ----------- | ---- |
| 66u: the brief session, inputs, prompt, checks, `ledger.KnownMoney` | `betBrief(ctx, name, kind)`, `briefBasis`, `ErrBriefLinesUnset`, `ErrBriefOverLimit` (`betbrief.go`), `BetBasisSchema` (`betbasis.go`) | **dropped from this branch**: steps S6 (brief), `briefInputs`, `briefKeep`, `KnownMoney`, U7, U8, I13 are 66u's; this row calls `betBrief` and `briefBasis` |
| 66t: the rendered task of a later bet | none; moved to M2g by the Operator (#417) | **dropped**: a `later` kind is a fail-closed refusal, and the plan row 66t is not registered here (no row on main, so 66j's After cell does not list it) |
| 66k: the new specification task a changed class reopens, the "waits for a bet" representation, the other three contexts on the rendered head, the close-out, `push-lost` of a session | not landed | **dropped**: this row records the changed class in the inventory and stops |
| 66r: `openShapeTask`, `finalID`, `finalIDs`, `drawTaskID` | `shapetask.go:40,54,87`; the ID is drawn by `Sessions.taskID` | `openShapeTask` is **parameterized** into `openTask` (reuse by kind stays the Shape wrapper's rule); `registerTask` takes the extra files of a commit (the binding) |
| 66r: the verified-head predicate | `loopPosition(v).Pos == posVerified` gives `Head`, `Base`, `Attempt` and (as `Session`) the verifier | adapted: `verifiedStep` reads `shapeView` + `loopPosition`, drives nothing; `betbrief.go`'s `verifiedHead` becomes a caller of it |
| 66h: the step `shape`, `Restart` wiring, `restartNames`, `filterAll`, `constraintsAt` | `steps.go:137-139,204-213`, `shape.go` | `bet` is appended after `shape` when the Scaffold is done; `shape` reads **no inventory**, so the recovery of condition 5 rests on the basis, not on a new filter |
| 66h/Q2: the milestone's issue | `r.value("issue.intake")` | as planned |
| 66s: `nextFactsID` | `promptfile.go:659` `nextFactsID(ids)` over `factRecordIDs(names)`; `factsNames(entries)` | reused as the one home (note of #416); no computation of its own |
| 66s: `PromptIn.Options` | landed | its consumer for the new tasks is 66k's (the binding is written here) |
| 66f/66g: `PostStatus`, `holdsLease`, `CheckRendered`, `CarryVerify(ctx, in, "")`, `Sessions.Bot/Web` | `closeout.go:120,237`, `handoffpost.go:427` | called as planned; `CarryInput.Verifier` is the verifier session of `loopPosition` at `posVerified` |
| 66i: `pointBindings`, `approverFor`, `approvalPoints`, `Decision` | `answers.go:110,180`, `batch.go:49`, `rules.go:25` | as planned |
| #412: `specInput`, `filesAt`, `isPRDName`, `isNumberedFacts` | `specinput.go`, `treeitems.go:30` | `filesAt` is the only tree selector used |
| 66z: `postBody` | `handoffpost.go:169` (the posting core, by task) | **extraction**: `postOnIssue(ctx, issue, data)` is its core; `postBody` calls it |
| plan: `pushAndBind`'s event pattern | `push.go:148` (needs a session start row) | not reused: `pushRendered` writes `push`/`bound` with session `—` itself |
| plan Q9 consumers of `bound` | `loopPosition` (stop on a `bound` of a non-developer), `RepairPushes` (resolves a `push` by session), `authors.go:73`, `statuses.go:126`, `pullrequest.go` `lastBound`, `LinkReviewRecord` (the PR head must be the reviewed head) | **audited**: each skips a `bound` with no session; the driver does not settle a verified head again once the rendered commit is bound (I12) |
| plan 2.4 "created strictly after the brief" | `copies.tsv` `created` (one-second resolution) | kept as the plan has it |
| tests: `drvWorld`, `shapedWorld`, `mustVerify`-style helpers, `ctlBrief` | `shapetask_integration_test.go`, `betbrief_integration_test.go` | reused; no second world |

No landed form makes a step impossible. The consumers of `bound` are the one unplanned cost; they are small guards, each with a test.
