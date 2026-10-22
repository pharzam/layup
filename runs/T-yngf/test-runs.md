# T-yngf: test runs

Task T-yngf, issue #479: a restart that fails while it holds the lease releases it. Times are from `date -u` after each run. Go 1.26.10 from mise.

## 1. The scaffolding

Inert only: the fields `lostLease`, `viaLease` and `endTimeout` of `state`, `leaseStore()` (returns `r.store`, used by `beat` and `leaseStep`), `leaseEndTimeout`, and `endLease`, which is the base behaviour (release after a last row `done` when this run holds the lease; nothing after a `fail`). `go vet` (no tag and `integration`) was clean before any run: **no red below is a compile failure**, each is an assertion failure. Also in the first slice: the fake forge answers a `defBranch` (so `probe` can fail), and the five old held-after-failure assertions (condition 4) now expect `released`.

## 2. Red (assertion failures on the base)

Unit, `go test -count=1 -run '<the tests>' ./internal/run`, 2026-10-22T10:32:14Z, exit 1:

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

Integration, `TMPDIR=... go test -count=1 -tags=integration -run '<the tests>' ./internal/run`, 2026-10-22T10:32:04Z to 10:32:14Z, exit 1:

| Test | First failing assertion |
| --- | --- |
| `TestAFailedRestartReleasesTheLeaseAndTheNextTakesItAtOnce` (3 cases: a step of Intake, before phase, inside phase) | `the lease after the fail is not released by the last commit` (lease `held`); and the next restart `after 68 reads of the lease; want done with none` |
| `TestARestartThatLostTheLeaseToAnotherRunPublishesNoRelease` | last row `own detail`; want the loss named (`the lease names the run cccccccccccccccc`) |
| `TestAPartialRecordsStepReleasesAndKeepsItsPrefix` | the log does not start with `lease: released`; the next restart `lease is held by the run fedcba9876543210` after 67 reads |
| the five old assertions (condition 4): `TestTheStepMeaningFails` (2 cases), `TestAFailedSpanCheckRunsNumberingAgain`, `TestTheCompletenessReviewWithNoOtherHarness` (2 cases), `TestTheStepMarkersFails` (6 cases) | `the lease is not released after a fail`; the e2e `scaffold_e2e_test.go:290` expects `released` (run at the e2e level in section 4) |

## 3. Controls (green on the base, by design)

- Unit: `TestAFailedReleaseIsNamedAfterTheFailure/a_done_row` (the rule of a row `done` is as it was); `TestARunThatHoldsNothingReleasesNothing` (no lease, another run's lease, a released lease: no write tried).
- Integration: `TestARestartThatRunsPhaseReleasesOnce` (one `lease: released` commit after a restart whose `phase` releases); `TestTwoRefusedReleasesMakeNoThirdAttempt` (a pre-receive hook refuses the release twice and would pass the third push: phase fails naming the loss, 2 pushes only, the lease stays `held`). On the base `phaseStep` makes no second attempt, so it is green; it binds the new end against a retry after a known loss.
