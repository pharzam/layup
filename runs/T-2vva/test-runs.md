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

## Appendix - the transcripts of the red runs (the output cut at 330 characters a line; the gate noise of the scripted harness left out)

### red-unit-1
```
2026-10-20T14:57:05Z
--- FAIL: TestABetsTableRefusesEachBrokenRule (0.00s)
    bet_test.go:39: a second first bet: err = <nil>; want line 3, column "kind"
    bet_test.go:39: a later bet before the first: err = <nil>; want line 2, column "kind"
    bet_test.go:39: no brief hash: err = <nil>; want line 2, column "brief_sha256"
    bet_test.go:39: a bet with no appetite: err = <nil>; want line 2, column "appetite_usd"
    bet_test.go:39: a no-bet with an appetite: err = <nil>; want line 2, column "appetite_usd"
    bet_test.go:44: control: 0 rows, <nil>; want 3 and no error
--- FAIL: TestTheBetRowIsItsAnswerAndItsBrief (0.00s)
    bet_test.go:57: betRow = []; want ["1" "first" "901" "a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1" "40.0" "30.5" "905" "db=pg web=go" "bet" "77" "2026-10-08T14:00:00Z"]
    bet_test.go:62: betRow of a no-bet = []; want ["1" "first" "901" "a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1" "" "" "905" "" "no-bet" "77" "2026-10-08T14:00:00Z"]
--- FAIL: TestTheOptionLinesAreTheColumnOptions (0.00s)
    bet_test.go:69: optionsColumn = []; want [db=pg web=go]
--- FAIL: TestThePhaseTagsAreTheFirstColumnOfSection9 (0.00s)
    bet_test.go:80: tags = []; want [1 alpha 3]
--- FAIL: TestTheAnswerIsTheEarliestFirstCopyOfTheApproverAfterTheBrief (0.00s)
    bet_test.go:111: betAnswer = {Reqs:map[] Facts:map[] Options:map[] Verdict: Milestone:}, [], false; want the copy 99 with the option db=pg
--- FAIL: TestTheSettledDecisionsAreTheOptionLines (0.00s)
    bet_test.go:126: settled = []; want [db] (web had one survivor)
--- FAIL: TestTheRequirementsOfAnOptionAreItsMappedRows (0.00s)
    bet_test.go:139: ids = [], <nil>; want [NFR-002 REQ-004] once each (NFR-009 is no requirement of the draft)
    bet_test.go:147: rows ["F-3" "Inv-1" "ADR-0001" "Q-002"]: err = <nil>; want the decision db: its option pg names no requirement
    bet_test.go:147: rows ["R-005"]: err = <nil>; want the decision db: its option pg names no requirement
    bet_test.go:147: rows []: err = <nil>; want the decision db: its option pg names no requirement
--- FAIL: TestTheRenderedRecordsAreOnePRDAndTheFactsFile (0.00s)
    bet_test.go:169: paths = []; want the facts record F-0005 and the PRD
--- FAIL: TestTheNextVersionOfTheInventory (0.00s)
    bet_test.go:201: the first version: 0 files, 0 rows, <nil>
--- FAIL: TestBetFollowsShapeOnlyOnceTheScaffoldIsDone (0.00s)
    bet_test.go:220: the last rows = [scaffold-push shape]; want shape, bet
FAIL
FAIL	github.com/pharzam/layup/internal/run	0.184s
--- FAIL: TestAPushAndABoundOfCodeHaveNoSession (0.00s)
    session_test.go:374: a push and a bound with no session: line 5, column "session": session is the empty value exactly for attempt, closed, rebased, merged and a refused whose reason is pair (a handoff has one unless code wrote it for a rendered task); want read
FAIL
FAIL	github.com/pharzam/layup/internal/records	0.010s
FAIL
exit=1
2026-10-20T14:57:07Z
```

### red-unit-2
```
$ go test -count=1 ./internal/run/ ./internal/records/
start 2026-10-20T15:08:26Z
exit=1
end 2026-10-20T15:08:27Z
--- FAIL: TestABetsTableRefusesEachBrokenRule (0.00s)
    bet_test.go:39: a second first bet: err = <nil>; want line 3, column "kind"
    bet_test.go:39: a later bet before the first: err = <nil>; want line 2, column "kind"
    bet_test.go:39: no brief hash: err = <nil>; want line 2, column "brief_sha256"
    bet_test.go:39: a bet with no appetite: err = <nil>; want line 2, column "appetite_usd"
    bet_test.go:39: a no-bet with an appetite: err = <nil>; want line 2, column "appetite_usd"
    bet_test.go:44: control: 0 rows, <nil>; want 3 and no error
--- FAIL: TestTheBetRowIsItsAnswerAndItsBrief (0.00s)
    bet_test.go:57: betRow = []; want ["1" "first" "901" "a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1" "40.0" "30.5" "905" "db=pg web=go" "bet" "77" "2026-10-08T14:00:00Z"]
    bet_test.go:62: betRow of a no-bet = []; want ["1" "first" "901" "a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1" "" "" "905" "" "no-bet" "77" "2026-10-08T14:00:00Z"]
--- FAIL: TestTheOptionLinesAreTheColumnOptions (0.00s)
    bet_test.go:69: optionsColumn = []; want [db=pg web=go]
--- FAIL: TestThePhaseTagsAreTheFirstColumnOfSection9 (0.00s)
    bet_test.go:80: tags = []; want [1 alpha 3]
--- FAIL: TestTheAnswerIsTheEarliestFirstCopyOfTheApproverAfterTheBrief (0.00s)
    bet_test.go:111: betAnswer = {Reqs:map[] Facts:map[] Options:map[] Verdict: Milestone:}, [], false; want the copy 99 with the option db=pg
--- FAIL: TestTheSettledDecisionsAreTheOptionLines (0.00s)
    bet_test.go:126: settled = []; want [db] (web had one survivor)
--- FAIL: TestTheRequirementsOfAnOptionAreItsMappedRows (0.00s)
    bet_test.go:139: ids = [], <nil>; want [NFR-002 REQ-004] once each (NFR-009 is no requirement of the draft)
    bet_test.go:147: rows ["F-3" "Inv-1" "ADR-0001" "Q-002"]: err = <nil>; want the decision db: its option pg names no requirement
    bet_test.go:147: rows ["R-005"]: err = <nil>; want the decision db: its option pg names no requirement
    bet_test.go:147: rows []: err = <nil>; want the decision db: its option pg names no requirement
--- FAIL: TestTheRenderedRecordsAreOnePRDAndTheFactsFile (0.00s)
    bet_test.go:169: paths = []; want the facts record F-0005 and the PRD
--- FAIL: TestTheNextVersionOfTheInventory (0.00s)
    bet_test.go:201: the first version: 0 files, 0 rows, <nil>
--- FAIL: TestBetFollowsShapeOnlyOnceTheScaffoldIsDone (0.00s)
    bet_test.go:220: the last rows = [scaffold-push shape]; want shape, bet
--- FAIL: TestACodeBoundWithNoSessionIsSkippedByItsReaders (0.00s)
    bet_test.go:235: lastBound = bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb, ""; want the last bound of a session
    bet_test.go:240: registerRevOf of the rendered commit = "", the event bound of  has no start row in sessions.tsv; want the commit that the caller gives
    bet_test.go:243: authorsFrom with a bound of code: the event bound of  has no start row in sessions.tsv; want none
    bet_test.go:249: loopPosition with a bound of code = {Pos:13 Session: Msg:T-aaaa: the session  has a bound and is not of the developer role Attempt:1 Base:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa Head: Ask:{Stage:0 Session: Row:[] Question: Owner: Base: Msg:}}; want no stop
    bet_test.go:252: renderedBound = false, false; want true with the bound of code and false without
FAIL
FAIL	github.com/pharzam/layup/internal/run	0.176s
--- FAIL: TestAPushAndABoundOfCodeHaveNoSession (0.00s)
    session_test.go:374: a push and a bound with no session: line 5, column "session": session is the empty value exactly for attempt, closed, rebased, merged and a refused whose reason is pair (a handoff has one unless code wrote it for a rendered task); want read
FAIL
FAIL	github.com/pharzam/layup/internal/records	0.011s
FAIL
```

### red-integration-1
```
$ TMPDIR=/home/layup/.cache/layup-session/tmp-green go test -count=1 -tags=integration -timeout 60m ./internal/run/ -run "TestTheBetsBlock|TestTheStepBet$|TestBriefLinesUnset|TestOnlyTheApprover|TestADecisionWhose|TestANoBet|TestALaterBet|TestTheBetStepIsAtomic|TestShapeStaysDone"
start 2026-10-20T15:08:34Z
exit=1
end 2026-10-20T15:09:06Z
--- FAIL: TestTheBetStepIsAtomicToItsLastAct (26.66s)
    --- FAIL: TestTheBetStepIsAtomicToItsLastAct/version (3.75s)
        bet_integration_test.go:346: the stopped step = {Name:bet Result:done Detail:not wired}; want the stop
    --- FAIL: TestTheBetStepIsAtomicToItsLastAct/pushevent (3.80s)
        bet_integration_test.go:346: the stopped step = {Name:bet Result:done Detail:not wired}; want the stop
    --- FAIL: TestTheBetStepIsAtomicToItsLastAct/pushed (3.80s)
        bet_integration_test.go:346: the stopped step = {Name:bet Result:done Detail:not wired}; want the stop
    --- FAIL: TestTheBetStepIsAtomicToItsLastAct/bound (3.84s)
        bet_integration_test.go:346: the stopped step = {Name:bet Result:done Detail:not wired}; want the stop
    --- FAIL: TestTheBetStepIsAtomicToItsLastAct/carry (3.86s)
        bet_integration_test.go:346: the stopped step = {Name:bet Result:done Detail:not wired}; want the stop
    --- FAIL: TestTheBetStepIsAtomicToItsLastAct/issue (3.83s)
        bet_integration_test.go:346: the stopped step = {Name:bet Result:done Detail:not wired}; want the stop
    --- FAIL: TestTheBetStepIsAtomicToItsLastAct/task (3.77s)
        bet_integration_test.go:346: the stopped step = {Name:bet Result:done Detail:not wired}; want the stop
--- FAIL: TestTheBetsBlockEqualsItsGoSchema (0.00s)
    bet_integration_test.go:92: the block bets and BetsSchema: schema bets: the number of columns: the block has "11"; the Go schema has "0"
--- FAIL: TestTheStepBet (3.54s)
    bet_integration_test.go:106: the step = {Name:bet Result:done Detail:not wired}; want the wait for the answer
--- FAIL: TestANoBetWritesItsRowAndWaits (3.55s)
    bet_integration_test.go:300: the step = {Name:bet Result:done Detail:not wired}; want {Name:bet Result:done Detail:waiting for the Operator: no bet}
--- FAIL: TestADecisionWhoseOptionNamesNoRequirementIsRefusedBeforeAnyWrite (3.62s)
    bet_integration_test.go:280: run 0: {Name:bet Result:done Detail:not wired}; want the refusal
--- FAIL: TestShapeStaysDoneAfterTheRenderedPush (3.62s)
    bet_integration_test.go:402: the step = {Name:bet Result:done Detail:not wired}
--- FAIL: TestALaterBetFailsClosed (3.65s)
    bet_integration_test.go:323: the step = {Name:bet Result:done Detail:not wired}; want a fail naming row 66t and no write
--- FAIL: TestOnlyTheApproverOfTheBetAnswersAndALookupThatFailsChoosesNoOne (3.68s)
    bet_integration_test.go:247: with dave named: {Name:bet Result:done Detail:not wired}; want the idea owner's answer to answer nothing
--- FAIL: TestBriefLinesUnsetWaitsWithNoSession (3.69s)
    bet_integration_test.go:225: the step = {Name:bet Result:done Detail:not wired} after 0 sessions; want the wait for the Operator, no session and no post
FAIL
FAIL	github.com/pharzam/layup/internal/run	30.373s
FAIL
```

### characterization
```
ok  	github.com/pharzam/layup/internal/run	114.286s
exit=0
2026-10-20T14:57:34Z
```

## Final gates on the head `5735e4e` (origin/main `99ff862` merged)

- `go build ./...`, `go vet ./...`, `go vet -tags=integration ./...`, `go vet -tags=e2e ./...`: exit 0; `gofmt -l .`: empty; `go test -count=1 ./...`: ok; `go test -count=1 -tags=e2e ./...`: ok.
- The full integration level, `TMPDIR=/home/layup/.cache/layup-session/tmp-green go test -count=1 -tags=integration -timeout 90m ./...`: 2026-10-20T15:50:50Z to 15:54:11Z, exit=0 (`/home/layup/.cache/layup-session/T-2vva-full-integration.txt`). No load flake occurred.
- `adr-lint`, `prd-lint`, `link-lint`, `run-discipline-tests` (138 passed), `nested-checkout-check`, `setup-check`, `git diff --check origin/main HEAD`: all OK.
