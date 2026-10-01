# T-0drh — the local checks before the freeze (plan step 10)

Head `6597c14`, 2026-10-01T09:03:16Z, on the LAYUP host (macOS, go1.27.1).

| Command | Exit | Last line |
| ------- | ---- | --------- |
| `sh docs/adr/adr-lint.sh` | 0 | adr-lint: OK (FAIL lines: 0) |
| `sh docs/prd/prd-lint.sh` | 0 | prd-lint: OK (FAIL lines: 0) |
| `sh docs/links/link-lint.sh` | 0 | link-lint: OK  1154 links resolved (FAIL lines: 0) |
| `sh docs/tests/run-discipline-tests.sh` | 0 | run-discipline-tests: 81 passed, 0 failed (FAIL lines: 0) |
| `sh docs/tests/nested-checkout-check.sh` | 0 | nested-checkout-check: OK  10 cases behaved (FAIL lines: 0) |
| `sh docs/setup/setup-check.sh` | 0 | setup-check: kit-linters OK (FAIL lines: 0) |
| `git diff --check 648b37f HEAD` | 0 | — (FAIL lines: 0) |
| `go build ./...` | 0 | — (FAIL lines: 0) |
| `go vet ./...` | 0 | — (FAIL lines: 0) |
| `go test -count=1 ./...` | 0 | ok  	github.com/pharzam/layup/internal/psb	0.227s (FAIL lines: 0) |

The run is on the head that the commit of this file has as its parent; this file changes no checked text.
Diff against the base `648b37f` at that head: 15 files changed, 1189 insertions(+), 7 deletions(-). The plan review's Budget maximum is 1,900 lines added plus removed over 20 files.
