# T-36xw: test runs

Task T-36xw, issue #471: the bet brief fits `brief.lines` on a real Intake. Times are from `date -u` after each run. Go 1.26.10 from mise (`GOROOT=$HOME/.local/share/mise/installs/go/1.26.10`).

## 1. Red (R8): the tests of the plan with the conditions 3, 4 and 5 of the plan review, against stubs

The stubs have the new signatures and the old behaviour (`operatorBlock` is false, `idRanges` joins with commas, `briefFactsPart` and `briefOpenPart` are `factsPart` and `openPart`, `LastChange` gives the empty string, no annex is written), so each red below is an assertion failure and not a compile error. The pre-commit hook runs the unit level, so the four red unit tests carry `heldRed(t)` in the commit of the red, and the runs below were made **without** it (the green commit removes it).

### 1.1 Unit, `internal/run`

Command: `go test -count=1 -run 'TestTheM2dFixture|TestAnIntakeOfTheM2dSize|TestTheOperatorsAnswers|TestTheOverLimitWait|TestTheLeftOutIDs' ./internal/run/`, 15:05:06Z to 15:05:07Z, exit 1, 4 of 5 fail (the fixture-count test passes: 38 spans, 25 facts, 319 answers, 305 of the Operator's, 150 open).

- `TestAnIntakeOfTheM2dSizeFitsTheLimitWithRoomForTheText`: first assertion `overLimit = the parts of code are 531 lines, over brief.lines 100; want none` (the plan's 531); then `room -431; want 78 and at least 20`; then `keep of 20 authored lines = 0 files, limit: the brief is 551 lines, over the limit 100`.
- `TestTheOperatorsAnswersAreRangesWithALinkAndEveryFactIsReachable`: the facts part holds one line `- Q-010 (need): answer Q-010` for each Operator answer, and not the five range-and-link lines (every class, a gap, a kept ID between left-out IDs, an unknown question, an idea owner's `context` answer, a follow-up of an Operator question); `the ranges name map[]; want exactly the left-out set {14 IDs}`; the open part lists `Q-051` to `Q-055` one by one and has no link.
- `TestTheOverLimitWaitFiresWhenTheInlinedPartAloneIsOver`: `err = the parts of code are 566 lines, over brief.lines 100; want "the parts of code are 113 lines ..."`. Its control (no Operator row, 112 lines) is the second assertion and is not reached by the first failing path.
- `TestTheLeftOutIDsAreNamedInExactRanges`: six of eight cases fail, for example `idRanges([Q-001 Q-002 Q-004]) = "Q-001, Q-002, Q-004"; want "Q-001 to Q-002, Q-004"`.

### 1.2 Integration, `internal/git`

Command: `go test -count=1 -tags=integration -run 'TestLastChange$' ./internal/git/`, 15:05:07Z, exit 1. First failing assertion: `a path not touched since the first: LastChange(HEAD, [a.txt]) = "", <nil>; want "0d0225d3ceee..."`.

### 1.3 Integration, `internal/run`

Command: `TMPDIR=/home/layup/.cache/layup-session/tmp-green go test -count=1 -tags=integration -timeout 30m -run 'TestTheLinkedFactsAreKept|TestAChangedTextLeftOut|TestBriefLinesUnsetStartsNoSession|TestAFreshBuildRefuses' ./internal/run/`, 15:05:07Z to 15:05:14Z, exit 1.

- `TestTheLinkedFactsAreKeptAtAPinnedRecordsCommitBeforeTheSession`: `the annex commit "", 0 annex commits, 1 sessions; want one commit and one session`.
- `TestAChangedTextLeftOutOfTheBriefRefusesTheOldCandidate` (both cases, an omitted answer and an omitted open question): `the annex SHA  (was ) after 0 annex commits; want a new commit, linked by the new brief and not the old SHA`.
- `TestBriefLinesUnsetStartsNoSession` and `TestAFreshBuildRefusesWhatItCannotBuildBefore`, now with the Operator's rows seeded: **pass** on the stubs, as they must: a wait that writes nothing is the current behaviour, and the seeded rows make the green code's wait the test of it (the annex would be needed).
- Stated limit: the checks of the same test after its first `Fatalf` (the `ParseAnswer` line `fact Q-005: need`, the basis row of the omitted fact and the rendered facts record) are guards that hold today, since the basis and `factRows` already hold every fact; they are reached only once the annex exists, and are shown green in section 2.

### 1.4 End to end, `cmd/layup`

Command: `go test -count=1 -tags=e2e -timeout 10m -run 'TestTheSpecificationTaskEndToEnd$' ./cmd/layup/`, 15:05:19Z to 15:05:38Z, exit 1. First failing assertion: `stage 5: the posted brief holds no link to blob/<sha>/bets/1/facts.md`.

## 2. Green: the code of the plan's section 4

The code: `git.LastChange`, `(*recordsStore).LastChange`, and in `internal/run/betbrief.go` `operatorIDs`, `idRanges`, `leftOutFacts`, `leftOutOpen`, `briefFactsPart`, `briefOpenPart`, `annexFiles`, `(*briefInputs).compose` and `(*state).pinAnnex` (called by `betBrief` after the two waits and the pair). `heldRed` and its four calls are removed. One correction of a test between the red and the green: the character class of the range-line pattern of the selection test lacked `Q`, and its "none in the brief" check matched the prefix of the summary line of a single ID; both made the test blind to the ranges, not to the code, and the test still fails on the stubs (the ranges name `map[]`).

| Level | Command | Times (UTC) | Exit | Result |
| ----- | ------- | ----------- | ---- | ------ |
| Unit | `go test -count=1 -v -run 'TestTheM2dFixture\|TestAnIntakeOfTheM2dSize\|TestTheOperatorsAnswers\|TestTheOverLimitWait\|TestTheLeftOutIDs' ./internal/run/` | 15:09:01Z | 0 | the five pass |
| Integration, git | `go test -count=1 -tags=integration -run 'TestLastChange$' ./internal/git/` | 15:08:21Z | 0 | `ok` |
| Integration, run | `TMPDIR=/home/layup/.cache/layup-session/tmp-green go test -count=1 -tags=integration -timeout 30m -run 'TestTheLinkedFactsAreKept\|TestAChangedTextLeftOut\|TestBriefLinesUnsetStartsNoSession\|TestAFreshBuildRefuses' ./internal/run/` | 15:08:21Z to 15:08:28Z | 0 | `ok` (4.3 s) |
| End to end | `go test -count=1 -tags=e2e -timeout 10m -run 'TestTheSpecificationTaskEndToEnd$' ./cmd/layup/` | 15:08:31Z to 15:08:53Z | 0 | `ok` (21.6 s) |

## 3. The change trimmed to the budget, and the gates on the final tree

After the green of section 2 the change was over the budget (969 lines against 600 and a tolerance of 120), so the tests, the spec text and the comments were tightened (commit `chore: T-36xw trim the change to the budget`): the fixture classes are one list, the range check expands its ranges in place, the rendered-facts check of the integration test is left to the e2e (which compares the file line by line with the records), and the guardrails lesson is shorter. No assertion of the conditions was removed. The figure after the trim is 718 lines over 15 files, `runs/T-36xw` excluded. The results below are of the trimmed tree.

| Level | Command | Times (UTC) | Exit | Result |
| ----- | ------- | ----------- | ---- | ------ |
| Unit | `go test -count=1 ./...` (the pre-commit hook of the trim commit) | 15:18Z | 0 | all packages `ok` |
| Integration, full | `TMPDIR=/home/layup/.cache/layup-session/tmp-green go test -count=1 -tags=integration -timeout 30m ./...` (output: `/home/layup/.cache/layup-session/T-36xw-full-integration.txt`) | 15:18:12Z to 15:22:30Z | 0 | 21 packages `ok`, no FAIL; `internal/run` 254.5 s, `internal/git` 1.1 s |
| End to end, full | `go test -count=1 -tags=e2e -timeout 10m ./...` | 15:22:44Z to 15:24:11Z | 0 | all `ok` (`cmd/layup` with the stage 5 link and the rendered facts check) |
