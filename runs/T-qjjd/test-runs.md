# T-qjjd: test runs

Task T-qjjd, issue #475: a task session that ends `crash` or `no-result` gets one new attempt. Times are from `date -u` after each run. Go 1.26.10 from mise.

## 1. The scaffolding, and the three kinds of result (condition 3 of the second plan review)

The first slice holds the tests and only what they need to compile:
- inert stubs in `internal/run/shapeloop.go`: the position `posRetry`, `type Replaced`, `replacedOf` (returns nil) and `priorAttempt` (returns b-1), wired nowhere;
- a trailing parameter `Replaced` on `AnswerAccepted` and `OpenQuestions` (`answer.go`, `tier.go`), ignored; the callers of the tests pass `nil`;
- the harness controls in `harnessSteps` (`shapetask_integration_test.go`): `crash-<step>` and `noresult-<step>` hold the call numbers of that step to end `exit 3` or `exit 0` with no result; the counter is `calls-<step>`.

The package compiled and `go vet` (no tag and `integration`) was clean before any run below, so **no red below is a compile failure.** Each is an assertion failure on current behaviour (the stub), listed by kind.

## 2. Red, behavioural (assertion failures)

Unit: `go test -count=1 ./internal/run`, 08:35:02Z to 08:35:03Z, exit 1.

| Test | First failing assertion |
| --- | --- |
| `TestAnEligibleEndOfAStageIsTheRetryPosition` | all 4 stages x 2 classes: position 13 (`posStop`, "the session S-p1 of the plan ended crash"); want the retry at attempt 1, base `aaaaaaa`; and the stage's run position in attempt 2 |
| `TestAnAttemptIsReplacedOnlyByTheRetryOfItsEligibleEnd` | `a crash and its retry: replaced map[]; want map[1:S-r1]`, and the same for `a no-result and its retry` |
| `TestTheBoundIsOncePerStageAcrossClasses` | `a developer crash after a retried plan review: ... Pos:13 ...; want the retry in attempt 2`; `a verifier crash, its retry, then the developer's first crash: ... Pos:13; want the retry in attempt 2` |
| `TestTheFindingsPassOverTwoReplacedAttempts` | `findingsOf for attempt 4 = [], "cccc..."; want the finding of S-v1 and no gate head`; `priorAttempt(4) = 3; want 1` |
| `TestAnAnswerIsAcceptedByTheAttemptThatReplacesTwoCrashedOnes` | `after two replaced attempts: "", false; want accepted by S-cccccccc of attempt 4`; `OpenQuestions after two replaced attempts = 1; want 0` |

Integration: `TMPDIR=... go test -count=1 -tags=integration ./internal/run -run '<the five names below>'`, 08:34:52Z to 08:34:57Z, exit 1.

| Test | First failing assertion |
| --- | --- |
| `TestACrashedPlanReviewGetsOneNewAttemptAndASecondEndStops` | `AdvanceShapeTask = {... Stop:T-5ms5: the session S-90f30344 of the plan review ended crash}, <nil>; want the verified head` |
| `TestACrashedDeveloperOrVerifierIsRunAgainInANewAttemptFromTheTaskBase` | `... Stop:T-6oer: the session S-6fa57c15 of specification ended crash ...; want the verified head` |
| `TestTheFindingsReachTheDeveloperAcrossARetry` | `... Stop:T-oqhb: the session S-b08beb97 of specification ended crash ...; want the verified head` |
| `TestARestartBeforeOrAfterTheRetryAttemptStartsExactlyOne` | `before the attempt: AdvanceShapeTask did not fail` (the driver stops at the crash and starts no attempt) |
| `TestTheAnswerReachesTheAttemptAfterTwoRetries` | `... Stop:T-4ewt: the session S-fb7d597c of specification ended crash ...; want the verified head` |

## 3. Guards (green before and after, by design)

- `TestEndsThatAreNotEligibleStayStops` (a crash with a refusal, a done with a refusal, `wall`, `output`, `start`): passes against the stub, as it must; it pins the boundary.
- `TestTheReviewCycleCapCountsNoRetry`: passes against the stub. Its records give each attempt explicitly, so the position of the latest attempt never looks at the replaced one; it guards that the retry adds no review round (the cap count) once the rule is wired.
- The negative rows of `TestAnAttemptIsReplacedOnlyByTheRetryOfItsEligibleEnd` (refused, closed, rebased, lost push, other base, ordering, not the last session, close-out, crashed answer session of the plan role and of the developer role, `wall`, `output`, `start`, gate failure, material round, answer handback) are `nil` against the stub by construction; they bind only once the rule is implemented, and each differs from the positive row in one fact.
- In `TestTheBoundIsOncePerStageAcrossClasses` the three mixed-class cases (crash/no-result, no-result/crash, crash/crash) pass against the stub (the position already stops "ended <class>"); they pin the message `ended no-result` after a crash.
- In `TestAnAnswerIsAcceptedByTheAttemptThatReplacesTwoCrashedOnes` the two guard asserts (one of the two replaced; none replaced) pass against the stub.

## 4. Green

The rule: `retries`, `replacedOf`, `priorAttempt`, `stageOf`, `eligibleEnd` and `posRetry` with its act in `internal/run/shapeloop.go`; `PromptFile` (`promptfile.go`), `AnswerAccepted` (`answer.go`) and `Position` -> `OpenQuestions` (`tier.go`) take the same `Replaced`, built from one `loopView` (events, session rows with roles and bases, the status of each kept result, the task kind).

- Unit level, `go test -count=1 ./...`: exit 0 (all five red tests of section 2 and both guards pass).
- Integration subset of section 2 plus `Position|Question`: `go test -count=1 -tags=integration ./internal/run -run '...'`, 08:36:12Z to 08:36:38Z, exit 0.
- The full gates are listed in section 5.
