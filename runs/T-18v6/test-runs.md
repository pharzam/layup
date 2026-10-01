# T-18v6 — the local checks before the freeze

Head `b2557cc`, 2026-10-01T18:16:52Z, on the LAYUP host (macOS, go1.27.1, git version 2.54.0 (Apple Git-157)).

| Command | Exit | Last line |
| ------- | ---- | --------- |
| `sh docs/adr/adr-lint.sh` | 0 | adr-lint: OK |
| `sh docs/prd/prd-lint.sh` | 0 | prd-lint: OK |
| `sh docs/links/link-lint.sh` | 0 | link-lint: OK  1244 links resolved |
| `sh docs/tests/run-discipline-tests.sh` | 0 | run-discipline-tests: 81 passed, 0 failed |
| `sh docs/tests/nested-checkout-check.sh` | 0 | nested-checkout-check: OK  10 cases behaved |
| `sh docs/setup/setup-check.sh` | 0 | setup-check: kit-linters OK |
| `git diff --check b48764f HEAD` | 0 | — |
| `go build ./...` | 0 | — |
| `go vet ./...` | 0 | — |
| `go test -count=1 ./...` | 0 | ok  	github.com/pharzam/layup/internal/tsv	0.328s |
| `go test -count=1 -tags=integration ./...` | 0 | ok  	github.com/pharzam/layup/internal/tsv	0.120s |
| `go test -count=1 -tags=e2e -timeout 10m ./...` | 0 | ok  	github.com/pharzam/layup/internal/tsv	0.221s |

This file is written after the run, in the commit that follows `b2557cc`. Diff against `b48764f` at `b2557cc`: 13 files changed, 1421 insertions(+), 10 deletions(-).
