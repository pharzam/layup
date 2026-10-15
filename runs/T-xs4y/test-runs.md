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

## After round 1 (fixes of the review), 2026-10-15T12:44Z

Nine material findings and six notes of round 1 are fixed in documents only (`loop.md`, `records.md`, `forge.md`, `packages.md`, `scaffold.md`, `architecture.md`, the plan, `sections.tsv` with one new row). Re-run: `sh runs/T-xs4y/sections.sh` exit 0 (33 of 33); `adr-lint`, `prd-lint`, `link-lint`, `run-discipline-tests` and `setup-check` exit 0; `go test -tags=integration ./internal/tsv/ ./cmd/layup/` ok (the two packages that read the changed documents); `git diff --check origin/main` exit 0. No Go file changed since the full integration run above.

## After round 2 (the fixes that wait for the Operator's answer), 2026-10-15T13:01Z

Round 2 ended `not mergeable, findings recorded` at the cap 1 (six material findings, seven notes). The fixes are a local commit, not pushed until the Operator chooses. Red then green for the section added in round 1's fix: with its heading renamed in a copy, `sh runs/T-xs4y/sections.sh` prints `FAIL docs/spec/loop.md: ## REQ-005 — The steps of the loop after the Scaffold`; restored, 33 of 33 `ok`. `adr-lint`, `prd-lint`, `link-lint`, `run-discipline-tests` and `setup-check` exit 0; `go test -tags=integration ./internal/tsv/ ./cmd/layup/` ok; `git diff --check origin/main` exit 0.

## After round 3 (the fixes that wait for the Operator's answer), 2026-10-15T13:29Z

Round 3 ended `not mergeable, findings recorded` at the cap 2 (six material findings, six notes). The fixes are a local commit, not pushed until the Operator chooses. `sections.sh` 33 of 33; `adr-lint`, `prd-lint`, `link-lint`, `run-discipline-tests` and `setup-check` exit 0; `go test -tags=integration ./internal/tsv/ ./cmd/layup/` ok; `git diff --check origin/main` exit 0.

## After round 4 (the fixes that wait for the Operator's answer), 2026-10-15T14:01Z

Round 4 ended `not mergeable, findings recorded` at the cap 3 (four material findings, nine notes). The fixes are a local commit, not pushed until the Operator chooses. `sections.sh` 33 of 33; `adr-lint`, `prd-lint`, `link-lint`, `run-discipline-tests` and `setup-check` exit 0; `go test -tags=integration ./internal/tsv/ ./cmd/layup/` ok; `git diff --check origin/main` exit 0.
