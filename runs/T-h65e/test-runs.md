# T-h65e - the test runs of row 66u

## S0 - the reconciliation of the plan (6250830299) and its review (6250925715) with the landed code

Read at `origin/main` = `7919c71` (the branch was at `68a8a74` with no work; the merge was a fast-forward). Each name was checked in the tree.

| Name of the plan or the review | Landed form | Note |
| ------------------------------ | ----------- | ---- |
| 66h: `constraintsAt`, `filterAll` | `options.go:73`, `:190` | called as planned, with `optionTables` (`:171`) over `store.ReadFiles("spec/options")`; no session is run |
| 66h: the asking predicate `leftOpen` | `shapequestions.go:25`, one function with the clause "no Shape row asks it" | **`stillOpen` extracted** (`shapequestions.go`), `leftOpen` is `stillOpen` without the questions that a Shape row asks; the existing test of `leftOpen` is the characterization net |
| 66h: the fixture (I7, branch A) | `TestTheStepShape` (`shape_integration_test.go`) on `drvWorld` | the integration world is `newDrvWorldOpts` + `settled()` + `shapeStep`; no copy of the world was made |
| 66h: `options.tsv` with the `decision` column | `spec/options/<decision>.tsv` (current) and `spec/options/round-1/` | the brief reads the current tables only |
| 66r: `finalID`, `finalIDs` | `shapetask.go:40`, `:54` | called; the register row's `requirements` is `finalIDs(draft)` |
| 66r: the verified-head predicate | **none stand-alone**: `Advance.Head` comes out of `AdvanceShapeTask`, which acts | **diverges, adapted**: `Sessions.verifiedHead` (`betbrief.go`) reads `shapeView` and `loopPosition` and gives `st.Head` and `st.Base` at `posVerified`; it drives nothing |
| 66r: the register rows of the two Shape tasks | `Task.PR`, first row of each kind (as `openShapeTask` returns) | as planned |
| 66r: the pinned draft mapping | `shapeIssueBlocks`: "Pinned draft: spec/draft.tsv, SHA-256 {sha}, read at the records commit {records}" | **diverges**: the SHA is in the issue body, which the `forge.Forge` interface cannot read; the pin is the register row's `requirements` equal to `finalIDs(draft)` and the PRD at the verified head (loop.md, The brief of the first bet) |
| 62a: `spec.Result`, `Chosen`, `Option.Reasons` | `internal/spec/shape.go:158-179` | as planned |
| plan: the PRD content (condition 1) | `spec` parses the PRD privately (`parsePRDs`) | **added** `spec.ContentOfPRDs` (a small export of the same parse) for the statement, the criterion and the success metrics |
| plan: `ReadTelemetry` / `ReadBudget` | `records.go:323` / `intake.go:414` | as the plan review's note 4 corrects |
| plan: `forge.PullRequestQuery{Number}` | `forge.go:90` | as the plan review's note 4 corrects |
| plan: I14 / I12 | the comparison is `I12` | as the plan review's note 4 corrects |
| plan: success criteria `- <R-NNN>: <statement>` | the statement of a success-criterion row is constant ("The delivery's success criteria."), its measure is the `criterion` column | **diverges**: the part shows the PRD's metrics, which is the verified content (condition 1) |

No landed form makes a step impossible or changes the scope: stop condition not met.

## S2 - the inert declarations and the extractions (`4dd75e4`)

Direct statuses; the commit hook ran `go test ./...` of the unit level (exit 0) before the commit. The integration subset of the plan was run twice at once by a
mistake of the launch (the first was backgrounded inside the shell and survived); both runs are the same command and both are `ok`:

```
$ TMPDIR=/home/layup/.cache/layup-session/tmp-green go test -tags=integration -timeout 30m -count=1 ./internal/run/ -run 'Phase|Revision|Shape|Questions'
start 2026-10-20T05:02:03Z (first) and 2026-10-20T05:02:55Z (second)
ok  	github.com/pharzam/layup/internal/run	383.799s
ok  	github.com/pharzam/layup/internal/run	393.581s
exit=0 (both)
end 2026-10-20T05:09:31Z
```
`go test -count=1 ./internal/ledger/` (the net of `MilestoneCap`, which `KnownMoney` was extracted from): `ok`, exit 0, in the hook run.

## S3 - the red runs, against the inert declarations, before the code

Every command was run directly; the exit status is the one `go` gave, written to a file with `echo $?` and never read through a pipe. The times are
`date -u` taken after the step. The stub of the working tree differs from the commit `4dd75e4` in one line: the commit has `checkMilestone(milestone)`, the
tests name `checkBriefRequest(milestone, kind)` (the plan review's condition 5 asked for the kind case of U11), so the stub of the red runs has that name
and nothing else changed. A first unit run (2026-10-20T05:09:43Z to 05:10:04Z, exit 1) stopped U7 and U12 at their controls, which were `t.Fatalf`; the
controls became `t.Errorf` so that every case runs red, and the run below is the second. The list is cut at 205 characters.

### Unit level, against the stub (`betbrief.go` and `betbasis.go` of `4dd75e4`, with the one stub `checkBriefRequest` that the tests name; `spec.ContentOfPRDs` of the same commit)
```
$ go test -count=1 ./internal/run/ ./internal/spec/ ./internal/ledger/
start 2026-10-20T05:10:27Z
exit=1
end 2026-10-20T05:10:44Z
```
- `TestTheSolutionPartIsCodesBytes`: `betbrief_test.go:89: part = ""; want "### The verified heads\n- T-sp01 (specification): [pull request #7](https://github.com/acme/target/pull/7), head 4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f\n- T-ar01 (ar`
- `TestTheFactsPartHoldsEveryFactWithItsClassAndText`: `betbrief_test.go:97: part = ""; want "- 1 (need): The tool runs offline.\n- 2 (constraint): The data stays in Postgres.\n- Q-001 (context): 5.0\n- Q-004 (constraint): Postgres only\n"`
- `TestTheClassOfAnAnswersFactIsThatOfItsQuote`: `betbrief_test.go:115: answerFactClass(Q-002) = ""; want constraint`
- `TestTheOpenPartNamesEachQuestionStillOpen`: `betbrief_test.go:134: part = ""; want "### The questions still open\n- Q-002: The question Q-002\n- Q-003: The question Q-003\n- Q-007: The question Q-007\n"`
- `TestTheRequirementsPartUsesTheRegisteredFinalIDs`: `betbrief_test.go:149: part = "", <nil>; want "### The requirements\n- REQ-001: Alpha Criterion: a check\n- NFR-002: p99 under 200 ms Criterion: measured\n### The success criteria\n- faster\n- cheaper\n"`
- `TestTheAssembledBriefIsTheGoldenAndPassesCheckBrief`: `betbrief_test.go:188: assemble differs from testdata/briefs/first.md: got 0 bytes, want 1084`
- `TestTheBriefIsRefusedForEachFinding`: `betbrief_test.go:201: control: kept 0 files, <nil>; want [bets/1/attempts/S-1a2b3c4d.md] and no error`
- `TestTheAppetiteCountsTheSessionsOwnSpend`: `betbrief_test.go:279: with the session's row: [], <nil>; want the finding "the money caps 40 plus the milestone-0 spend 11 are over band.U 50"`
- `TestThePartsOfCodeOverTheLimitWaitBeforeAnySession`: `betbrief_test.go:308: brief.lines 28: err = <nil>; want "the parts of code are 29 lines, over brief.lines 28" wrapping ErrBriefOverLimit`
- `TestAMilestoneThatIsNoPathPartIsRefused`: `betbrief_test.go:317: checkBriefRequest("") = <nil>; want the milestone refused`
- `TestTheBasisTableRefusesEachBrokenRule`: `betbrief_test.go:353: control: rows = 0, schema bet-basis: no column; want 12`
- `TestTheRefusedPromptHoldsTheFinding`: `betbrief_test.go:394: the prompt lacks "the bet of milestone 1, for LAYUP.":`
- `TestTheLinkSetIsThePullRequestsAndTheChangedFilesOfTheHeads`: `betbrief_test.go:411: links = map[]; want map[https://github.com/acme/target/blob/4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f/docs/prd/PRD-0001-x.md:true https://github.com/acme/target/blob/a1a1a1a1a1a1a1a1a1`
- `TestContentOfPRDsGivesTheWordsOfTheRows`: `prdcontent_test.go:16: Statements = map[]; want REQ-001 Alpha and NFR-002 with its escaped bar`

### Integration level, against the same stub
```
$ TMPDIR=/home/layup/.cache/layup-session/tmp-green go test -count=1 -tags=integration -timeout 30m ./internal/run/ ./internal/tsv/ -run 'TestTheBriefIsKeptFromOneSessionOfTheRun|TestARefusedBriefIsTriedAgainWithItsFinding|TestAMissingOutputFileEndsNoResult|TestTheAppetiteIsJudgedAfterTheBriefSessionsTelemetry|TestAStopBeforeThePromotionPromotesWithNoSession|TestAKeptBriefIsDoneBeforeWithNoWrite|TestBriefLinesUnsetStartsNoSession|TestAFreshBuildRefusesWhatItCannotBuildBefore|TestTheCapOfMilestoneZeroEndsTheRetry|TestTheLinkSetIsTheChangedFilesAndPullRequestsOfTheVerifiedHeads|TestTwoSurvivorsGiveABriefWithNoRecommendation|TestAQuestionStillOpenIsShownInTheRabbitHoles|TestTheBasisRoundTripsThroughBriefBasisInAFreshState|TestTheBetBasisBlockEqualsItsGoSchema|TestTheBetBriefBlocksOfLoopMd|TestEverySchemaBlockIsBuiltOrNotYetBuilt'
start 2026-10-20T05:10:44Z
exit=1
end 2026-10-20T05:17:52Z
```
- `TestTheBriefIsKeptFromOneSessionOfTheRun`: `betbrief_integration_test.go:120: betBrief = 0 bytes, done false, <nil>; want the bytes of bets/1/brief.md`
- `TestARefusedBriefIsTriedAgainWithItsFinding`: `betbrief_integration_test.go:170: betBrief = 0 bytes, <nil> after 0 sessions; want a brief after two`
- `TestAMissingOutputFileEndsNoResultAndTheNextRunStartsAgain`: `betbrief_integration_test.go:191: betBrief = <nil> after 0 sessions; want a *PhaseError no-result for the missing file after one`
- `TestTheAppetiteIsJudgedAfterTheBriefSessionsTelemetry`: `betbrief_integration_test.go:214: betBrief = "", <nil> after 0 sessions; want the brief of attempt 2 after two`
- `TestAStopBeforeThePromotionPromotesWithNoSession`: `betbrief_integration_test.go:241: a stop at the comment: 0 bytes, <nil>, brief.md false, candidates []; want the error, no brief and the candidate`
- `TestAFreshBuildRefusesWhatItCannotBuildBefore`: `betbrief_integration_test.go:297: brief.lines 10: err = <nil> after 0 sessions; want ErrBriefOverLimit and none`
- `TestTheCapOfMilestoneZeroEndsTheRetry`: `betbrief_integration_test.go:317: err = <nil> after 0 sessions; want ErrMilestoneCap after one, with its one telemetry row`
- `TestTheLinkSetIsTheChangedFilesAndPullRequestsOfTheVerifiedHeads`: `betbrief_integration_test.go:332: betBrief = <nil> after 0 sessions; want attempt 2 kept, attempt 1 refused for the link to a file the head did not change`
- `TestTwoSurvivorsGiveABriefWithNoRecommendation`: `betbrief_integration_test.go:354: betBrief = <nil> after 0 sessions; want the second attempt kept, with both survivors unmarked`
- `TestAQuestionStillOpenIsShownInTheRabbitHoles`: `betbrief_integration_test.go:373: betBrief = <nil>:`
- `TestTheBasisRoundTripsThroughBriefBasisInAFreshState`: `betbrief_integration_test.go:400: briefBasis = {Requirements:[] Survivors:map[] Dropped:map[] Facts:[] Open:[] Heads:map[] Appetite:{USD: Hours:}}, <nil>`
- `TestTheBetBasisBlockEqualsItsGoSchema`: `betbrief_integration_test.go:418: schema bet-basis: the number of columns: the block has "3"; the Go schema has "0"`
- `TestTheBetBriefBlocksOfLoopMd`: `betbrief_integration_test.go:432: the block bet-brief of loop.md:`

Two of the tests failed on a defect of their helper (`setParameter` set `—` through a set row, which `CheckParameters` refuses), not on the brief; the helper was fixed and they ran again:

### The two again
```
$ TMPDIR=... go test -count=1 -tags=integration -timeout 30m ./internal/run/ -run 'TestAKeptBriefIsDoneBeforeWithNoWrite|TestBriefLinesUnsetStartsNoSession'
start 2026-10-20T05:18:12Z
exit=1
end 2026-10-20T05:20:06Z
```
- `TestAKeptBriefIsDoneBeforeWithNoWrite`: `betbrief_integration_test.go:270: a kept brief: 0 bytes, done false, <nil>; head moved true, POST 31 (was 31), PATCH 2 (was 2), 0 sessions; want done, the same bytes and no write`
- `TestBriefLinesUnsetStartsNoSession`: `betbrief_integration_test.go:288: err = <nil>; want ErrBriefLinesUnset`

`TestAKeptBriefIsDoneBeforeWithNoWrite` took its `head` before the two pushes of its own arrangement, so it would have failed at green too; the order was fixed and it ran again:

### `TestAKeptBriefIsDoneBeforeWithNoWrite` again
```
$ TMPDIR=... go test -count=1 -tags=integration -timeout 30m ./internal/run/ -run 'TestAKeptBriefIsDoneBeforeWithNoWrite'
start 2026-10-20T05:20:19Z
exit=1
end 2026-10-20T05:21:29Z
```
- `TestAKeptBriefIsDoneBeforeWithNoWrite`: `betbrief_integration_test.go:270: a kept brief: 0 bytes, done false, <nil>; head moved false, POST 31 (was 31), PATCH 2 (was 2), 0 sessions; want done, the same bytes and no write`

`TestAShapeTaskWithNoVerifiedHeadIsAnErrorBeforeAnySession` was added after the code, from the review of the coverage of `verifiedHead` (mutant M30). It was run red afterwards against the stub (the working tree's `betbrief.go`, `betbasis.go` and `full.go` were replaced by the stub and restored; SHA-256 checked, `sha256sum -c`: OK for the three files):

### The added test
```
$ TMPDIR=... go test -count=1 -tags=integration -timeout 30m ./internal/run/ -run 'TestAShapeTaskWithNoVerifiedHeadIsAnErrorBeforeAnySession'
start 2026-10-20T05:33:02Z
exit=1
end 2026-10-20T05:33:48Z
```
- `TestAShapeTaskWithNoVerifiedHeadIsAnErrorBeforeAnySession`: `betbrief_integration_test.go:357: betBrief on an unverified task = <nil> after 0 sessions; want an error that names T-5sdt, no session and no write`

### Green at the stub, each beside its red control (labelled, none counted as red)

- `TestKnownMoneyIsTheSumOfTheKnownRows` (U9) is a **characterization** of the extraction from `MilestoneCap`: green at `4dd75e4` and after; `TestMilestoneCap` stays green.
- In `TestTheAppetiteCountsTheSessionsOwnSpend`: "without the session's row" (49 of 50), "USD 30" and "a cap below band.B" pass at the stub; their red controls are "with the session's row", "hours over" and "a cap above band.U".
- In `TestThePartsOfCodeOverTheLimitWaitBeforeAnySession`: the line limit 29 (the boundary) passes at the stub; the red controls are 28 and 20.
- In `TestAMilestoneThatIsNoPathPartIsRefused`: the valid milestones `1`, `M2e`, `a.b-c_d` and the kind `first` pass at the stub; the red controls are the six names that are no path part and the kinds `later` and `other`.
- In `TestTheBriefIsRefusedForEachFinding`: a one-survivor option and a word that holds an option pass at the stub; their red controls are the option-word cases. The case "an indirect recommendation" is the **known limit** of the option-word rule, characterised and not a proof: it passes the rule at the stub and at green (condition 2 of the plan review).
- `TestEverySchemaBlockIsBuiltOrNotYetBuilt` is green at the stub (`bet-basis` is in `notYetBuilt` since `11aa82b`) and at green; its red control is `TestTheBetBasisBlockEqualsItsGoSchema`.
- `TestBriefLinesUnsetStartsNoSession`: "0 sessions and no write" is green at the stub (a green negative); its red control is the missing `ErrBriefLinesUnset`.
- `TestTheAppetiteIsJudgedAfterTheBriefSessionsTelemetry`: "the known money before the call is 0.0" is an assertion that holds at the stub and at green (without the session's own spend the 8 fits U 10).

## S4 to S6 - green, in the order of the plan (the parts and the link set, the assembly and the keep and the basis, the appetite, the promotion, the retry and the recovery)

One `feat:` commit holds the tests, the code and this log. The code was built until the unit level passed, and then until the integration list passed. The first integration
run of the list failed on two defects of the test arrangement (the scripted harness exited 1 when no cost file was left, which the class `crash` reports; a question row of the
kind `meaning` that `CheckQuestion` refuses) and passed after both were fixed.

```
$ go test -count=1 ./internal/run/ ./internal/spec/ ./internal/ledger/
start 2026-10-20T05:46:58Z
ok  	github.com/pharzam/layup/internal/run
ok  	github.com/pharzam/layup/internal/spec
ok  	github.com/pharzam/layup/internal/ledger
exit=0
end 2026-10-20T05:47:01Z

$ TMPDIR=/home/layup/.cache/layup-session/tmp-green go test -count=1 -tags=integration -timeout 30m ./internal/run/ ./internal/tsv/ -run '<the list above and TestAShapeTaskWithNoVerifiedHeadIsAnErrorBeforeAnySession>'
start 2026-10-20T05:33:56Z
ok  	github.com/pharzam/layup/internal/run	667.583s
ok  	github.com/pharzam/layup/internal/tsv	0.088s
exit=0
end 2026-10-20T05:45:46Z
```
