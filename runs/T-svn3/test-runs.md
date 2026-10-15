# T-svn3 — the runs of the tests

Times from the clock of the host, UTC, 2026-10-15.

## Step T, red (17:30:26Z)

`go test ./internal/session -run "Context|StepTable|NoTwo|ACallers"` with `context_test.go` and no `context.go`: build failure, `undefined: Kind`, `undefined: ContextKinds`, `undefined: StepTable`. A compilation failure is a valid red start; it proves no single rule, so the mutants below do.

## Step C, green (17:30:49Z)

`go test ./internal/session`: ok. (A first build at 17:30:40Z failed on the name `steps`, which `process_test.go` already holds; the variable is `defaultSteps`.)

## Mutant runs, each rule broken alone (the code restored after each)

Each run is `go test ./internal/session -run "Context|StepTable|NoTwo|ACallers"`, exit 1 with the mutant, exit 0 restored:

| Mutant | Exit |
| ------ | ---- |
| the order of the first two kinds swapped | 1 |
| a phase step (`gap-review`) with the kinds | 1 |
| a task step (`verification`) with no kinds | 1 |
| the row `close-out` removed | 1 |
| the role of `plan-review` changed | 1 |
| the tier of `implement` changed | 1 |
| `ContextKinds` returns its own slice, not a copy | 1 |
| a step name twice (`close-out` renamed `implement`) | 1 |

## Verification before the freeze

- `go build ./...` and `go vet ./...`: exit 0. `gofmt -l .`: no file.
- `go test -count=1 ./...`: exit 0 (17:33:36Z).
- `go test -count=1 -tags=integration ./...` (TMPDIR on the disk): 19 packages ok and one failure at 17:35Z: `TestProcessTimesTheFirstOutput` in `internal/session` (a timing test of 0.3 s and 0.2 s, issue #291's flake: the host ran other tasks' integration runs at the same time). The package alone, run again at 17:52:56Z: ok.
- `adr-lint`, `prd-lint`, `link-lint`, `run-discipline-tests` and `setup-check`: exit 0. `git diff --check`: exit 0.
