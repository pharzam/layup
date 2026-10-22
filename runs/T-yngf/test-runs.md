# T-yngf: test runs

Task T-yngf, issue #479: a restart that fails while it holds the lease releases it. Times are from `date -u` after each run. Go 1.26.10 from mise.

## 1. The scaffolding

Inert only: the fields `lostLease`, `viaLease` and `endTimeout` of `state`, `leaseStore()` (returns `r.store`, used by `beat` and `leaseStep`), `leaseEndTimeout`, and `endLease`, which is the base behaviour (release after a last row `done` when this run holds the lease; nothing after a `fail`). `go vet` (no tag and `integration`) was clean before any run: **no red below is a compile failure**, each is an assertion failure. Also in the first slice: the fake forge answers a `defBranch` (so `probe` can fail), and the five old held-after-failure assertions (condition 4) now expect `released`.

## 2. Red (assertion failures on the base)

Unit, the new tests by name, 2026-10-22T10:32:14Z, exit 1:

| Test | First failing assertion |
| --- | --- |
| `TestAFailedRowKeepsItsDetailAndTheLeaseIsReleased` | the lease `held`, the commits `[]`; want the lease released by one commit |
| `TestAFailedReleaseIsNamedAfterTheFailure/a_failed_row` | row detail `own`; want `own; the lease: disk full` |
| `TestARunThatLostTheLeaseNamesTheLossAndPublishesNoRelease` | detail `own`; want `own; the lease: the records push was refused; the lease names the run cccccccccccccccc`; the loss `<nil>`, not remembered |
| `TestTwoRefusalsAreTerminalAndNoReleaseIsTriedAfterThem` | detail `own`; want the same suffix naming the run itself; 0 writes tried, want 2 |
| `TestTheEndTriesAgainAfterAReleaseOfThePhaseThatFailed` | the lease `held`, 1 write; want released, 2 writes |
| `TestTheHeartbeatEndsBeforeTheRelease` | `the end returned while a beat was in flight; it must join the heartbeat first` |
| `TestAHeartbeatErrorThatIsNotALossStillReleases` | both rows: the lease `held`; want released |
| `TestTheReleaseOutlivesTheContextOfTheStepsAndIsBounded` | `a cancelled context of the steps stopped the release`; and no deadline named |

Integration, the new tests by name with the integration tag, 2026-10-22T10:32:04Z to 10:32:14Z, exit 1:

| Test | First failing assertion |
| --- | --- |
| `TestAFailedRestartReleasesTheLeaseAndTheNextTakesItAtOnce` (3 cases: a step of Intake, before phase, inside phase) | `the lease after the fail is not released by the last commit` (lease `held`); and the next restart `after 68 reads of the lease; want done with none` |
| `TestARestartThatLostTheLeaseToAnotherRunPublishesNoRelease` | last row `own detail`; want the loss named (`the lease names the run cccccccccccccccc`) |
| `TestAPartialRecordsStepReleasesAndKeepsItsPrefix` | the log does not start with `lease: released`; the next restart `lease is held by the run fedcba9876543210` after 67 reads |
| the five old assertions (condition 4): `TestTheStepMeaningFails` (2 cases), `TestAFailedSpanCheckRunsNumberingAgain`, `TestTheCompletenessReviewWithNoOtherHarness` (2 cases), `TestTheStepMarkersFails` (6 cases) | `the lease is not released after a fail`; the e2e `scaffold_e2e_test.go:290` expects `released` (run at the e2e level in section 4) |

## 3. Controls (green on the base, by design)

- Unit: `TestAFailedReleaseIsNamedAfterTheFailure/a_done_row` (the rule of a row `done` is as it was); `TestARunThatHoldsNothingReleasesNothing` (no lease, another run's lease, a released lease: no write tried).
- Integration: `TestARestartThatRunsPhaseReleasesOnce` (one `lease: released` commit after a restart whose `phase` releases); `TestTwoRefusedReleasesMakeNoThirdAttempt` (a pre-receive hook refuses the release twice and would pass the third push: phase fails naming the loss, 2 pushes only, the lease stays `held`). On the base `phaseStep` makes no second attempt, so it is green; it binds the new end against a retry after a known loss.

## 4. Green

The code: `endLease` (join the heartbeat, read the lease, release on a context that the steps do not cancel, one-minute bound, error named after the failure), `noteLoss` (a typed `LostError` kept in `lostLease`, from a beat, a write and a release), in `internal/run/steps.go`. Commit 3d42b91e holds the tests and the code together because the hook refuses a red commit (guardrails §2).

- Unit, `go test -race -count=3` on the nine new tests: exit 0. `go test -count=1 ./...`: exit 0.
- Full integration, `TMPDIR=... go test -count=1 -tags=integration -timeout 30m ./...`, 2026-10-22T10:36:49Z to 10:40:50Z, exit 0 (`internal/run` 239.9 s); output in `/home/layup/.cache/layup-session/T-yngf-full-integration.txt`.
- e2e, `go test -count=1 -tags=e2e -timeout 10m ./...`, 2026-10-22T10:36:59Z to 10:40:39Z: `cmd/layup` ok (121.7 s), the rest ok.
- build, vet (no tag, integration, e2e) and gofmt clean; adr, prd, link, discipline, setup and nested-checkout lints OK; `freeze-check/run-all.sh` 85 passed.
- The old-tests rule (condition 4): the requirement changed from "a failed restart leaves the lease held until the takeover" to "a failed restart releases it", because the held lease made the next restart wait `3 × lease.H`; the five assertions changed only in that expectation and its comment.

## 5. Round 1, F1: every fenced records path notes the loss

Red (an assertion failure on 7c0d788, before the fix; the fix's scaffolding is only the test), the integration run of `TestALossFoundByAProbeWriteIsKnownAtTheEnd` and `TestTwoRefusedReleasesMakeNoThirdAttempt`, 2026-10-22T11:10:39Z, exit 1: the probe row failed naming the loss and the hook counted 2 refused pushes, but `lease.tsv` ended `released` and the log began with `lease: released`, the third push after the loss. `TestTwoRefusedReleasesMakeNoThirdAttempt` stayed green (a control, now on the shared hook helper).

Green: the fix (`Sessions.Lost`, `(*state).fenced`, `noteLoss` in `Sessions.commit` and `holdsLease`): unit with `-race -count=3` exit 0; the integration tests of the restart world exit 0 (times in the output of the full run below). Note 4 of the review: the heartbeat test now lands the in-flight beat (`land`), and asserts `lease: heartbeat 1` before `lease: released` as the last write.

The first form of the oracle looked for `lease: released` in the whole log, which holds Start's own release; the test was corrected to count exactly one (Start's) and the lease `held`. The red above stands (lease `released`, two releases in the log). Green after the correction: the restart-world tests (`TestALossFound...`, `TestTwoRefused...`, `TestAFailedRestart...`, `TestAPartial...`, `TestARestartThat...`), 2026-10-22T11:12:11Z to 11:12:19Z, exit 0.

Round 1 gates on the merged head: full integration 2026-10-22T11:12:52Z to 11:16:58Z, exit 0 (`internal/run` 245.5 s), output in `/home/layup/.cache/layup-session/T-yngf-full-integration-2.txt`; e2e exit 0 (`cmd/layup` 101.5 s); unit, vet (no tag, integration, e2e), gofmt and the lints clean.
