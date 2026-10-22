# T-dtw6: test runs

Task T-dtw6, issue #491: a session refused for a rule path alone gets one new attempt, and every committing step is told the rule paths. Times are from `date -u` after each run. Go 1.26.10 from mise. Slices follow the plan review (condition 1): tests and red evidence; parser and findings; prompt rendering; retry classification; documentation; close-out.

## 1. Scaffolding and red

The tests were written first with stubs that compile and do nothing (`rules.PayloadPaths`, `rulePathFindings`, `rulePathGuidance` return nothing). `go vet` (no tag and `integration`) was clean before every run, so no red below is a compile failure.

Defect tests (fail at db71ef5), `go test -count=1 ./internal/rules ./internal/run`, exit 1, finished before 2026-10-22T17:12:18Z:

| Test | First failing assertion | Acceptance criterion (slice) |
| --- | --- | --- |
| `TestPayloadPathsReadsTheHeadersOfABinaryDiff` | all 8 cases: `PayloadPaths = [], <nil>` | the paths come from the payload (slice 2) |
| `TestPayloadPathsFailsClosedOnAnUnreadableHeader` | all 12 cases: `PayloadPaths = [], nil; want an error` | the parser fails closed (slice 2) |
| `TestRulePathFindingsKeepTheRulePathsOfThePayload` | `the findings name []; want [docs/engineering-discipline.md docs/guardrails.md scripts/run.sh docs/setup/new-file.md]` | the findings name the rule paths, with the register's exception (slice 2) |
| `TestEachCommittingStepNamesTheRulePathsAndTheirRefusal` | the 6 rendered prompts do not hold the sentence; the 4 constants hold no `{rule-paths}` | each committing step names the rule paths (slice 3) |
| `TestARefusedRulePathIsAnEligibleEndOfEachStage` | `position 13 (posStop), "the session S-p1 was refused (rule-path 0123…)"`; want the retry | the retry of the class (slice 4) |
| `TestARulePathRefusalHasItsOwnRetryPerStage` | `[0 1 2]: Pos 13 ... was refused (rule-path`; want the retry of the third end | the bound of the class (slice 4) |
| `TestTheDemoRecordsOfTAm8bGiveAttempt2AndTheRenderedPromptNamesThePaths` | `Pos:13 ... S-06aeffc9 was refused (rule-path a506…)`; want the retry at `d03c865` | the demo's recovery (slice 4) |
| `TestTheRulePathGuidanceSurvivesConsecutiveReplacementsOfOtherClasses` | `[], <nil>; want the two findings` | the guidance across replacements (slice 4) |
| `TestTheRefusedPayloadIsReadAtTheRecordsBoundaryOrTheRunErrs` | 8 cases (one is a read that fails): `[], <nil>; want an error` | the records boundary (slice 4) |

Integration (`-tags=integration`, the tests below by name), exit 1, finished 2026-10-22T17:12:20Z:

| Test | First failing assertion |
| --- | --- |
| `TestPayloadPathsReadsTheRealDiffBinary` | `PayloadPaths = [], <nil>; want [a b.md café.md gone.md img.bin x.md]` (the real `git.DiffBinary` producer, `--no-renames`) |
| `TestADeveloperRefusedForARulePathGetsOneNewAttemptThatNamesThePath` | `AdvanceShapeTask = {... Stop: the session S-428356cd was refused (rule-path d0c6…)}; want the verified head` |
| `TestARestartBeforeOrAfterTheRulePathRetryAttemptStartsExactlyOne` | `before the attempt: AdvanceShapeTask did not fail` (no attempt 2 is ever committed) |
| `TestARefusedPayloadThatCannotBeReadStopsTheRunBeforeAnotherDeveloper` | `AdvanceShapeTask = <nil>; want the run's error that names payloads/<sha>` |

Regression controls (green before and after, they pin what stays a stop or a question): `TestARefusalOfAnotherKindOrStatusIsNoRulePathRetry` (`rule-paths` with and without a payload, `rule-pathx`, `workflow`, `base`, `push-refused`, `rule-path` with `artifact` in both orders, two `rule-path` refusals, a crash or wall refused, the statuses `blocked` and none, the close-out, the answer sessions, and the `needs_context` question).

The goldens are not red tests of their own: `TestThePromptFileOfEachStepIsItsGoldenFile` and `TestReopenedSpecificationPromptIsItsGoldenFile` go red when the instruction changes (slice 3) and green when the six files are updated in the same commit.

## 2. Green

The change was committed in slices, each with its tests, because the hook refuses a red commit: the parser (`0f484e5f`), the findings (`3c97ec34`), the prompt rendering with the six goldens and the four blocks of `loop.md` (`a8a94b74`), the retry class, the guidance and the records boundary (`437b50e1`), then the documents. Times are from `date -u` after each run.

- Unit, `go test -count=1 ./...`: every package ok (the hook run at each commit).
- The new integration tests: ok, 2026-10-22T17:17:05Z.
- The first full integration run failed once, at `TestAPromptFileIsBuiltFromTheRecordsOfARealTaskAndTheSessionReadsTheFiles`: `refused: context 1091 1000`, because the sentence made the prompt longer than the 1000 tokens of the stand-in pair. The test now gives the pair `Context = 100000`, as its neighbours do (`fae66d7a`).
- The full integration level, `TMPDIR=/home/layup/.cache/layup-session/tmp-green go test -count=1 -tags=integration -timeout 30m ./...` (output in `/home/layup/.cache/layup-session/T-dtw6-full-integration.txt`), on the final code: every package ok, exit 0, finished 2026-10-22T17:30:07Z.
- E2E, `go test -count=1 -tags=e2e -timeout 10m ./...`: every package ok, exit 0, 2026-10-22T17:31:43Z.
- `go build ./...`, `go vet ./...` (no tag, `integration`, `e2e`), `gofmt -l .`, the six lints and `git diff --check origin/main HEAD`: clean. `setup-check.sh` reports only `kit-history FAIL orphan: docs/tasks/T-dtw6.md has no line in backlog.md or completed.md`, which the close-out's completed-log line ends.
