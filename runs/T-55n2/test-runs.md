# T-55n2 — the local checks before the freeze (plan step 8)

Head `3e30c37`, 2026-10-01T13:57:34Z, on the LAYUP host (macOS, go1.27.1).

| Command | Exit | Last line |
| ------- | ---- | --------- |
| `sh docs/adr/adr-lint.sh` | 0 | adr-lint: OK |
| `sh docs/prd/prd-lint.sh` | 0 | prd-lint: OK |
| `sh docs/links/link-lint.sh` | 0 | link-lint: OK  1231 links resolved |
| `sh docs/tests/run-discipline-tests.sh` | 0 | run-discipline-tests: 81 passed, 0 failed |
| `sh docs/tests/nested-checkout-check.sh` | 0 | nested-checkout-check: OK  10 cases behaved |
| `sh docs/setup/setup-check.sh` | 0 | setup-check: kit-linters OK |
| `git diff --check 7cdd346 HEAD` | 2 | +- Refuter: C(note): I could not refute the finding. It is real, and it is a note.  1. The cited cell is word for word at docs/plan/README.md:210: "\/ K25 \/ O-121 \/ row 20 \/ The pilot rule in row 20's plan review; the count in the ADR task. \/". Line 181 says that the "Settled by" column names "the task that settles it".  |
| `python3 runs/T-55n2/plan-check.py` | 0 | PASS  each defect is settled before it is read |
| `go build ./...` | 0 | — |
| `go vet ./...` | 0 | — |
| `go test -count=1 ./...` | 0 | ok  	github.com/pharzam/layup/internal/psb	0.313s |

This file is written after the run, in the commit that follows `3e30c37`. Diff against `7cdd346` at `3e30c37`: 20 files changed, 4071 insertions(+), 36 deletions(-); the Budget maximum is 4,800 lines added plus removed over 36 files.
