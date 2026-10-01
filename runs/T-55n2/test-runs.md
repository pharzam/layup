# T-55n2 — the local checks before the freeze (plan step 8)

Head `ee69da6` (the fix of review round 1), 2026-10-01T14:44:14Z, on the LAYUP host (macOS, go1.27.1). The first freeze was `a88f99f`.

| Command | Exit | Last line |
| ------- | ---- | --------- |
| `sh docs/adr/adr-lint.sh` | 0 | adr-lint: OK |
| `sh docs/prd/prd-lint.sh` | 0 | prd-lint: OK |
| `sh docs/links/link-lint.sh` | 0 | link-lint: OK  1232 links resolved |
| `sh docs/tests/run-discipline-tests.sh` | 0 | run-discipline-tests: 81 passed, 0 failed |
| `sh docs/tests/nested-checkout-check.sh` | 0 | nested-checkout-check: OK  10 cases behaved |
| `sh docs/setup/setup-check.sh` | 0 | setup-check: kit-linters OK |
| `git diff --check 7cdd346 HEAD` | 0 | — |
| `python3 runs/T-55n2/plan-check.py` | 0 | PASS  each defect is settled before it is read |
| `go build ./...` | 0 | — |
| `go vet ./...` | 0 | — |
| `go test -count=1 ./...` | 0 | ok  	github.com/pharzam/layup/internal/psb	0.378s |

This file is written after the run, in the commit that follows `ee69da6`. Diff against `7cdd346` at `ee69da6`: 21 files changed, 4180 insertions(+), 36 deletions(-); the Budget maximum is 4,800 lines added plus removed over 36 files.
