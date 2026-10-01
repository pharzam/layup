# T-2tc2 — the local checks before the freeze

Head `bb22c0f`, 2026-10-01T18:17:41Z, on the LAYUP host (macOS, go1.27.1, git version 2.54.0 (Apple Git-157)).

| Command | Exit | Last line |
| ------- | ---- | --------- |
| `sh docs/adr/adr-lint.sh` | 0 | adr-lint: OK |
| `sh docs/prd/prd-lint.sh` | 0 | prd-lint: OK |
| `sh docs/links/link-lint.sh` | 0 | link-lint: OK  1248 links resolved |
| `sh docs/tests/run-discipline-tests.sh` | 0 | run-discipline-tests: 81 passed, 0 failed |
| `sh docs/tests/nested-checkout-check.sh` | 0 | nested-checkout-check: OK  10 cases behaved |
| `sh docs/setup/setup-check.sh` | 0 | setup-check: kit-linters OK |
| `git diff --check b48764f HEAD` | 0 | — |
| `go build ./...` | 0 | — |
| `go vet ./...` | 0 | — |
| `go test -count=1 ./...` | 0 | ok  	github.com/pharzam/layup/internal/psb	0.344s |
| `go test -count=1 -tags=integration ./...` | 0 | ok  	github.com/pharzam/layup/internal/psb	0.226s |
| `go test -count=1 -tags=e2e -timeout 10m ./...` | 0 | ok  	github.com/pharzam/layup/internal/psb	0.210s |

This file is written after the run, in the commit that follows `bb22c0f`. Diff against `b48764f` at `bb22c0f`: 18 files changed, 1864 insertions(+), 12 deletions(-).
