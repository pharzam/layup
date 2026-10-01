# T-0drh — the local checks at the close-out

Head `e8456cd`, 2026-10-01T10:13:43Z, on the LAYUP host (macOS, go1.27.1).

| Command | Exit | Last line |
| ------- | ---- | --------- |
| `sh docs/adr/adr-lint.sh` | 0 | adr-lint: OK (FAIL lines: 0) |
| `sh docs/prd/prd-lint.sh` | 0 | prd-lint: OK (FAIL lines: 0) |
| `sh docs/links/link-lint.sh` | 0 | link-lint: OK  1165 links resolved (FAIL lines: 0) |
| `sh docs/tests/run-discipline-tests.sh` | 0 | run-discipline-tests: 81 passed, 0 failed (FAIL lines: 0) |
| `sh docs/tests/nested-checkout-check.sh` | 0 | nested-checkout-check: OK  10 cases behaved (FAIL lines: 0) |
| `sh docs/setup/setup-check.sh` | 0 | setup-check: kit-linters OK (FAIL lines: 0) |
| `git diff --check 648b37f HEAD` | 0 | — (FAIL lines: 0) |
| `go build ./...` | 0 | — (FAIL lines: 0) |
| `go vet ./...` | 0 | — (FAIL lines: 0) |
| `go test -count=1 ./...` | 0 | ok  	github.com/pharzam/layup/internal/psb	0.348s (FAIL lines: 0) |

This file is written after the run, in the commit that follows `e8456cd`; it changes no checked text. Diff against `648b37f` at `e8456cd`: 21 files changed, 1582 insertions(+), 8 deletions(-). The Budget maximum is 1,900 lines added plus removed over 20 files.
