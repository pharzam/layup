# T-0sj8: test runs

Task T-0sj8, issue #478: the shared result instruction says a comment file is never an artifact, and a session refused for `artifact` alone gets its own one new attempt per stage. Times are from `date -u` after each run. Go 1.26.10 from mise.

## 1. The scaffolding

The tests came first, with only what they need to compile:
- `eligibleEnd` (`internal/run/shapeloop.go`) returns a third value, the class, which is `uncertain` for a `crash` or `no-result` end and is inert otherwise; `retries` ignores it and keeps `spent` by stage;
- the harness control `artifact-<step>` (`internal/run/shapetask_integration_test.go`): the call numbers of that step whose comment file is also named by an `artifact` row.

The package compiled, and `go vet` (no tag and `integration`) was clean, before any run below, so no red below is a compile failure.

## 2. Red (assertion failures on the current behaviour)

Unit: `go test -count=1 ./internal/run`, exit 1, finished before 10:13:07Z.

| Test | First failing assertion |
| --- | --- |
| `TestARefusedArtifactIsAnEligibleEndOfEachStage` | all 4 stages: position 13 (`posStop`, "the session S-p1 was refused (artifact)"); want the retry at attempt 1, base `aaaaaaa`, and the stage's run position in attempt 2 |
| `TestTheDemoRecordsOfTUt9gGiveAttempt3` | `retries = map[1:S-38157561], ""; want attempt 1 replaced and S-141cc16c pending` |
| `TestAnArtifactRefusalHasItsOwnRetryPerStageInEitherOrder` | `crash, then artifact: {Pos:13 ... was refused (artifact)}; want the retry of S-r2 in attempt 2` |
| `TestAnArtifactRefusalOfTheVerifierCountsNoReviewRound` | `the verifier refused for an artifact: {Pos:13 ...}; want the retry at the task base aaaaaaa` |
| `TestTheResultInstructionSaysACommentFileIsNoArtifact` | the constant and the prompt of each of the 8 steps do not hold the sentence |

Guards (green before and after, they pin what stays a stop or a question): `TestARefusalOfAnotherKindOrStatusIsNoEligibleEnd` (closed-attempt, branch, push-lost, two refusal events, a crash/no-result/wall with `artifact`, the statuses `blocked` and none, the close-out, the answer sessions), `TestARefusedResultThatAsksAQuestionKeepsTheQuestionPath` (status `needs_context` keeps the question), and the replaced assertion of `TestEndsThatAreNotEligibleStayStops` (a `done` with `refused artifact` no longer stops; `branch` and `closed-attempt` still do).

Integration: the three tests below by name with the tag `integration`, exit 1, finished 2026-10-22T10:13:07Z.

| Test | First failing assertion |
| --- | --- |
| `TestAPlanReviewRefusedForAnArtifactRowGetsOneNewAttempt` | `AdvanceShapeTask = {... Stop:T-q9si: the session S-0a93a04d was refused (artifact)}, <nil>; want the verified head` |
| `TestAVerifierRefusedForAnArtifactIsRunAgainInANewAttemptFromTheTaskBase` | the same stop, `was refused (artifact)`; want the verified head |
| `TestARestartBeforeOrAfterTheArtifactRetryAttemptStartsExactlyOne` | `before the attempt: AdvanceShapeTask did not fail` (no attempt 2 is ever committed) |

## 3. Green

The change (`eligibleEnd` with the class `artifact`, `retries` keyed by stage and class, the sentence in the constant, `loop.md` and the 11 golden files) was committed with the tests, because the hook refuses a red commit (`f53ef58b`). Times are from `date -u` after each run.

- Unit, `go test -count=1 ./...`: every package ok, 2026-10-22T10:14:04Z (the hook run at the commit).
- The three new integration tests and the neighbours of the retry and the prompt file: ok, 10:14:19Z.
- The full integration level, `go test -count=1 -tags=integration -timeout 30m ./...` (output in `/home/layup/.cache/layup-session/T-0sj8-full-integration.txt`): one failure, finished 10:18:53Z, exit 1: `TestTheRulesResultOfAPushedHeadStaysThatOfItsBoundSessionWhenTheRecordsMoveOn` (`statuses_integration_test.go:42`), the step `records` of Start: `git push` to a `file://` remote, `fatal: Could not read from remote repository`. The test touches none of the changed code; run alone three times (`-count=3`) it passed. The `internal/run` package was then run again in full (`.../T-0sj8-full-integration-rerun.txt`): ok in 256.9 s, exit 0. All other packages of the first run were ok.
- E2E, `go test -count=1 -tags=e2e -timeout 10m ./...`: every package ok, `cmd/layup` 85.6 s, 10:24:57Z.
- `go build ./...`, `go vet ./...` (no tag, `integration`, `e2e`) and `gofmt -l .`: clean.
