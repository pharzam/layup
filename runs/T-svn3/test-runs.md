# T-svn3 — the runs of the tests

Times from the clock of the host, UTC, 2026-10-15.

## Step T, red (17:30:26Z)

`go test ./internal/session -run "Context|StepTable|NoTwo|ACallers"` with `context_test.go` and no `context.go`: build failure, `undefined: Kind`, `undefined: ContextKinds`, `undefined: StepTable`. A compilation failure is a valid red start; it proves no single rule, so the mutants below do.

## Step C, green (17:30:49Z)

`go test ./internal/session`: ok. (A first build at 17:30:40Z failed on the name `steps`, which `process_test.go` already holds; the variable is `defaultSteps`.)

## Mutant runs, each rule broken alone (the code restored after each)

The mutant runs are of 17:3xZ (between the green at 17:30:49Z and the full runs, not timed one by one). The order of the rows and the copy of the rows are covered by the table test (row i against row i of the literal) and by the copy test, and have no mutant of their own. Each run is `go test ./internal/session -run "Context|StepTable|NoTwo|ACallers"`, exit 1 with the mutant, exit 0 restored:

| Mutant | Exit | The test that fails (by the test file's reading) |
| ------ | ---- | ---- |
| the order of the first two kinds swapped | 1 | `TestContextKindsAreInTheOrderOfTheSpecification` |
| a phase step (`gap-review`) with the kinds | 1 | `TestTheStepTableHasEachStepWithItsRoleTierAndKinds` |
| a task step (`verification`) with no kinds | 1 | `TestTheStepTableHasEachStepWithItsRoleTierAndKinds` |
| the row `close-out` removed | 1 | `TestTheStepTableHasEachStepWithItsRoleTierAndKinds` |
| the role of `plan-review` changed | 1 | `TestTheStepTableHasEachStepWithItsRoleTierAndKinds` |
| the tier of `implement` changed | 1 | `TestTheStepTableHasEachStepWithItsRoleTierAndKinds` |
| `ContextKinds` returns its own slice, not a copy | 1 | `TestACallersChangeDoesNotReachTheTableNorTheKinds` |
| a step name twice (`close-out` renamed `implement`) | 1 | `TestNoTwoStepsShareAName` and the table test |

## Verification before the freeze

- `go build ./...` and `go vet ./...`: exit 0. `gofmt -l .`: no file.
- `go test -count=1 ./...`: exit 0 (17:33:36Z).
- `go test -count=1 -tags=integration ./...` (TMPDIR on the disk), first run 17:33:36Z to 17:52:31Z: 19 packages ok and one failure at about 17:35Z, `TestProcessTimesTheFirstOutput` in `internal/session` (a timing test of 0.3 s and 0.2 s; the known flake of issue #291, with other tasks' integration runs on the host at the same time). `internal/session` alone, 17:52:56Z: ok. **The whole level again after the round, 18:06:42Z to 18:22:50Z: exit 0, 20 packages ok.**
- `adr-lint`, `prd-lint`, `link-lint`, `run-discipline-tests` and `setup-check`: exit 0. `git diff --check`: exit 0.

## After round 1 (notes only), 2026-10-15

Round 1 ended `nothing material in scope` with seven notes. The fixes are documents and one Go comment (`loop.md`'s pointer and "a step of a task", `session.md`'s reading of "for each step" and the words `by-class` and `owner-of-its-kind`, the comment of `StepRow`), and this record; no code changed. Re-run: `go vet ./internal/session` and `go test ./internal/session` ok; the whole integration level exit 0 (above); the five checks and `git diff --check` exit 0.
