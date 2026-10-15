# T-xs4y — the runs of the tests

## Step T: the section check, red (before D1), 2026-10-15T11:55Z

`sh runs/T-xs4y/sections.sh` at `c2fecde`: exit 1, 32 of 32 headings missing, none present (`docs/spec/loop.md` does not exist; `records.md` has no `NFR-001 — The records of the task loop`). First lines:

```
FAIL docs/spec/loop.md: ## REQ-005 — The task loop (the task loop, step by step (later-p2-phase-loop-handoffs; decision 2))
FAIL docs/spec/loop.md: ### The attempts (the attempts, the retry and the new attempt (decision 2; session.md))
FAIL docs/spec/loop.md: ### The independence of the reviewers (independence of the plan review and the verifier; the authors of a change (decision 4; session.md))
FAIL docs/spec/loop.md: ### The tier and the context of a task (the tier of a task from its plan and the record kinds of each step's context (condition 1; session.md))
```

## The block test, red then green, 2026-10-15T12:0xZ

`go test -tags=integration ./internal/tsv/` with the ten new blocks (`bets`, `inventory`, `inventory-index`, `options`, `parameters`, `task-gates`, `task-spec-check`, `task-tests`, `tasks`, `transitions`) in `docs/spec/records.md` and `notYetBuilt` empty: FAIL, `the block bets is listed 0 times in built and notYetBuilt; want 1` for each. With the ten names in `notYetBuilt`: ok.

## Green, at the head before the freeze

- `sh runs/T-xs4y/sections.sh`: exit 0, 32 of 32 headings present (the heading of the statuses section was written with its backticks in `sections.tsv`, a correction of the test file, not of the check).
- `adr-lint`, `prd-lint`, `link-lint`, `run-discipline-tests` and `setup-check`: exit 0 each. `git diff --check origin/main`: exit 0.
- `go build ./...` and `go vet ./...`: exit 0.
- `go test -count=1 -tags=integration ./...` (TMPDIR on the disk, as `/tmp` is small): exit 0, 20 packages ok, 11 min 42 s (12:07:33 to 12:19:15Z).
- Two defects of the author's own text were found by these checks and fixed before the freeze: a second code span in the cell "Starts a program" of `internal/spec` (`TestPackageRules`: "cannot read Starts a program"), and a marker whose file was not yet tracked (`setup-check` `markers`: "stale"), as the lesson in `docs/guardrails.md` §2 says.
