# T-v2sv — the Retrospective: test runs

Row 66c of `M2e`, issue #359; the record settles #372 under R10. Times are UTC
from `date -u`. The mutants ran against a backup of `internal/run/retrospective.go`
at sha256 `ea7e659ad2d6293b08f332e11f60ed258d261199bc20139f3693a3935f380386`,
restored and hash-verified after each run, each run taken from `go test`'s own
exit status and output, unfiltered except for the `-run` anchor of the catching
tests.

## Red first (each for its own reason)

- `2026-10-16T11:45:04Z` — `go test ./internal/run/ -run '<the eleven names>' -v`
  → FAIL (exit 1): nine tests fail on the refusing stub ("the Retrospective step
  is not built" / "names no column milestone|time" / "names no file").
  Green on arrival: `TestAReadErrorIsReturned` and
  `TestATableWithADuplicateMilestoneIsRefused` — refusal-path guards that a stub
  that refuses everything satisfies; their bite is shown by the mutants M10 and
  M11 (a swallowed error, a check that never fires).
- `2026-10-16T11:45:24Z` — `go test -tags=integration ./internal/run/ -run
  '^(TestTheSchemaOfRetrospectivesEqualsItsBlock|TestTheRetrospectiveIsOneCommitAndDoneBeforeOnARestart)$' -v`
  → FAIL (exit 1): "docs/spec/ has no block retrospectives" (the block lands in
  the same change as the `built` entry and the schema) and "the Retrospective
  step is not built".

## Green

- `2026-10-16T11:46:56Z` — the eleven unit tests: `ok github.com/pharzam/layup/internal/run`.
- `2026-10-16T11:47:08Z` — the two integration tests: `ok`, 2.3s (the real
  adapter on the `httptest` forge and real git: one commit of
  `retrospectives.tsv` alone, the row read back, the restarted state done
  before — the records head does not move).
- `2026-10-16T11:51:43Z` — control after the last mutant restore: both runs
  green; `go build ./...`, `go vet ./...`, `go vet -tags=integration
  ./internal/run/`: exit 0.
- `2026-10-16T11:52:21Z` — `go test ./...` (unit, whole module): no failures.
- `2026-10-16T11:52:5_Z` — `go test -tags=integration ./internal/tsv/ -run
  TestEverySchemaBlockIsBuiltOrNotYetBuilt`: PASS; `... ./cmd/layup/ -run
  TestPackageRules` (after the `packages.md` change, guardrails §2): PASS.

## Mutants (all caught; none build-broken or not applied)

| # | Mutation (kept building) | Caught by | Exit |
| - | ------------------------ | --------- | ---- |
| M1 | `CheckRetrospectives`: `r[1] == "" && false` (the owner check never fires) | `TestADashInTimeIsRefusedByTheOwner` | FAIL |
| M2 | `ReadRetrospectives`: the owner check's `err != nil && false` (the call is skipped) | `TestADashInTimeIsRefusedByTheOwner` | FAIL |
| M3 | `retrospectiveFiles`: `r[0] == milestone && false` (done-before never matches) | `TestARetrospectiveDoneBeforeWritesNothing` | FAIL |
| M4 | `nextBet = "loop"` | `TestTheStepWritesItsEventAndGivesTheNextBet` and the two row-table tests | FAIL |
| M5 | `stamp(time.Time{})` for the event's time | `TestTheEventTimeIsTheRunsClockInUTC` (and the write test) | FAIL |
| M6 | `Fenced` removed — a bare `Commit` | `TestTheWriteIsFencedByTheLease` (no retry, no `LostError`) | FAIL |
| M7 | the commit also writes `sessions.tsv` | `TestTheStepStartsNoSession` | FAIL |
| M8 | `RetrospectivesSchema` column `time` retyped `text` | `TestTheSchemaOfRetrospectivesEqualsItsBlock` | FAIL |
| M9 | the new row prepended instead of appended | `TestAnExistingTableKeepsItsRowsAndAppends` | FAIL |
| M10 | the error of `retrospectiveFiles` dropped (`_ =`) | `TestATableThatReadRefusesStopsTheStep`, `TestAnEmptyOrDashedMilestoneIsRefusedAndNothingCommits`, `TestADashInTimeIsRefusedByTheOwner` | FAIL |
| M11 | the owner check reads `r[0]` instead of `r[1]` | `TestADashInTimeIsRefusedByTheOwner` | FAIL |

`2026-10-16T11:48:52Z`–`11:51:4_Z`. A `Truncate`-only mutation is not listed:
`stamp` formats UTC second precision, so it cannot change the written bytes —
it is recorded as unobservable, not as caught.

## Checks

`2026-10-16T11:53:01Z` — `sh docs/adr/adr-lint.sh`: OK; `sh docs/prd/prd-lint.sh`:
OK; `sh docs/links/link-lint.sh`: OK (2732 links); `sh
docs/tests/run-discipline-tests.sh`: 138 passed, 0 failed; `sh
docs/setup/setup-check.sh`: OK (after `git add` of the new files); `git diff
--check` and `git diff --cached --check`: clean.

Not run here, by instruction: `go test -tags=integration ./...` — the full
integration level runs once before the freeze, serialized by the integration
coordinator after the Go head is stable.
