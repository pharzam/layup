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

## S1b - plan amendment 1 (comment 6293953072) and the instruction

The Operator's coordinator selected (A): the `instruction-specification` block (`loop.md`, `instructions.go`, `testdata/prompts/specification.txt`) now shows `## 6.`, `## 7.`,
`## 8. Success metrics`, the matrix, `## 9. Rollout & phases`, and asks for a `**Non-goals.**` paragraph. The last two additions beyond the amendment (the headings 6 and 7 and the
paragraph) are **decided here**: `RenderPRD` finds the rows of the requirements in `## 6.`/`## 7.` and the paragraph it appends to, and the specimen of the instruction had none of them.
`TestTheNewInputsAreInTheirKindsAfterTheExistingItems` and the specimen test (now five forms) follow; `internal/spec/prdcontent_test.go` cites `## 8.`.
The merge of `origin/main` at `99ff862` (#444, T-05f4) was made as a plain merge (`eb506cc`; its message is git's default, the hook accepted it).

## S2 - the inert stubs and the extractions (`423bf25`)

The commit hook ran `go test ./...` (exit 0). The characterization net of the two extractions (`postOnIssue` from `postBody`; `verifiedStep` out of `verifiedHead`), run on `423bf25`:
`TMPDIR=/home/layup/.cache/layup-session/tmp-green go test -tags=integration -timeout 30m -count=1 ./internal/run/ -run 'Comment|Post|Register|Revision|Brief|Shape'`, start
2026-10-20T14:55:38Z, `ok ... 114.286s`, exit=0, end 2026-10-20T14:57:34Z (`characterization.txt`).

## S3 - the red runs, before the code (not committed alone: the hook runs `go test`, guardrails.md:101)

Direct statuses (`echo "exit=$?"` after the command, no pipe); the times are `date -u` taken after each step. The files hold the command, the times and the output.

- Unit, first set (`red-unit-1.txt`, 2026-10-20T14:57:05Z to 14:57:07Z, exit=1): U1 `TestABetsTableRefusesEachBrokenRule` (`err = <nil>; want line 3, column "kind"` for each of the five cases;
  the control `0 rows; want 3`), U2 `betRow = []; want [...]`, U5 `optionsColumn = []`, U6 `tags = []; want [1 alpha 3]`, U3 `betAnswer = {...}, [], false; want the copy 99`, U13 `settled = []; want [db]`,
  U14 `ids = []; want [NFR-002 REQ-004]` and the three refusals, U11 `paths = []`, the version test `0 files, 0 rows`, U15 `the last rows = [scaffold-push shape]; want shape, bet`, and U12 in `internal/records`
  (`line 5, column "session": session is the empty value exactly for ...; want read`). The negatives (`phaseTagsOf` of no section, `settledDecisions` of no option line, `betAnswer` of copies that answer nothing,
  the `bet` row absent where the Scaffold is not done, the refusals of `session`/`refused` in U12) are green at the stub by design, beside their red twins.
- Unit, the readers of an event bound of code (`red-unit-2.txt`, 2026-10-20T15:08:26Z to 15:08:27Z, exit=1): `lastBound = bbbb..., ""; want the last bound of a session`; `registerRevOf ... the event bound of  has no
  start row in sessions.tsv`; `authorsFrom ... the same`; `loopPosition ... the session  has a bound and is not of the developer role`; `renderedBound = false, false; want true ... and false without`. The guards of
  `RepairPushes`, `LinkReviewRecord` and the driver are red in the integration tests below (the stop at the push event, and `TestShapeStaysDoneAfterTheRenderedPush`).
- Integration (`red-integration-1.txt`: `TMPDIR=/home/layup/.cache/layup-session/tmp-green go test -count=1 -tags=integration -timeout 60m ./internal/run/ -run "TestTheBetsBlock|TestTheStepBet$|TestBriefLinesUnset|TestOnlyTheApprover|TestADecisionWhose|TestANoBet|TestALaterBet|TestTheBetStepIsAtomic|TestShapeStaysDone"`,
  2026-10-20T15:08:34Z to 15:09:06Z, exit=1): every test failed at its first assertion against `betStep` that answered `done not wired` (`the step = {Name:bet Result:done Detail:not wired}; want the wait for the answer`),
  and `TestTheBetsBlockEqualsItsGoSchema` with `the number of columns: the block has "11"; the Go schema has "0"`. The red is the stub's, not a typo of the test: the same worlds are green below.

## S4 - green

- Unit: `go test -count=1 ./...` green after the code (the hook of the commit). The first integration run of the new tests (2026-10-20T15:12:03Z to 15:12:43Z, exit=1) had three failures, **all errors of the tests, not of the code**: `TestALaterBetFailsClosed`
  wrote an empty field for `—` in its fixture; `TestADecisionWhose...` and `TestANoBetWritesItsRowAndWaits` compared the branch of the specification task with "none" while it holds the verified head, and counted the session's `push`
  as the push of code. They were corrected in the tests (second run of the three: ok) with no change of code.
- The full integration level, started 2026-10-20T15:13:19Z: `internal/run` failed one test, `TestTheStepForgeCheckChecksTheForgeAndCommitsItsTableAtEachRestart` (the notes of progress count 24 rows, not 23, once the Scaffold is done: the expected change of the plan,
  `forgecheck_integration_test.go`), which was corrected; the final run is recorded at the end of this file.
