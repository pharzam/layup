# T-jc45 — test runs

Task `T-jc45` (#442, serves `F-0003#47`): bound CI's integration level at `-timeout 20m` per test binary (`docs/adr/0027-bound-the-integration-level-at-20-minutes-in-ci.md`).

## The red (CI integration timeouts at 10m)

In `.github/workflows/ci.yml:316`, the `integration` step of job `tests (unit → integration → e2e)` ran `go test -tags=integration -timeout 10m ./...`. Two CI runs timed out in `internal/run` with no assertion failure:

1. **`main` at `2866a3febf99130a074a5608df98b32509a70536`** (merge of PR #436, row 66z), CI run `39422675099`, command `go test -tags=integration -timeout 10m ./...`:
   ```
   tests (unit → integration → e2e)	integration	2026-10-20T12:40:37.8087237Z panic: test timed out after 10m0s
   tests (unit → integration → e2e)	integration	2026-10-20T12:40:37.8187973Z FAIL	github.com/pharzam/layup/internal/run	600.012s
   tests (unit → integration → e2e)	integration	2026-10-20T12:40:37.8191349Z FAIL
   tests (unit → integration → e2e)	integration	2026-10-20T12:40:37.8400605Z ##[error]Process completed with exit code 1.
   ```
2. **PR #437 head `614bef4adc4b1114a34eb77848c3d0b7aea9e948`** (row 66u), CI run `39425720339`, command `go test -tags=integration -timeout 10m ./...`:
   ```
   tests (unit → integration → e2e)	integration	2026-10-20T13:03:33.4050603Z panic: test timed out after 10m0s
   tests (unit → integration → e2e)	integration	2026-10-20T13:03:33.4317944Z FAIL	github.com/pharzam/layup/internal/run	600.014s
   tests (unit → integration → e2e)	integration	2026-10-20T13:03:33.4320805Z FAIL
   tests (unit → integration → e2e)	integration	2026-10-20T13:03:33.4324642Z ##[error]Process completed with exit code 1.
   ```

For contrast, the two preceding CI passes under the same command measured:
- `0246d9c72ecae7c3509a7522990cbf39706f3fd4`, run `39389141724`: `ok  	github.com/pharzam/layup/internal/run	519.198s`
- `51d949b3fd2e7dd1a1d275585988fab0cc651055` (PR #436 head), run `39421233368`: `ok  	github.com/pharzam/layup/internal/run	493.632s`

## Local checks (head `8646392`)

Each command was run directly (never piped through `tail` or `grep`), with `date -u` read before and after each step:

| Command | Start | End | Exit |
| --- | --- | --- | --- |
| `sh docs/adr/adr-lint.sh` | 2026-10-20T13:39:18Z | 2026-10-20T13:39:18Z | exit 0 |
| `sh docs/prd/prd-lint.sh` | 2026-10-20T13:39:18Z | 2026-10-20T13:39:18Z | exit 0 |
| `sh docs/links/link-lint.sh` | 2026-10-20T13:39:18Z | 2026-10-20T13:39:31Z | exit 0 |
| `sh docs/tests/run-discipline-tests.sh` | 2026-10-20T13:39:31Z | 2026-10-20T13:39:34Z | exit 0 |
| `sh docs/tests/nested-checkout-check.sh` | 2026-10-20T13:39:34Z | 2026-10-20T13:39:34Z | exit 0 |
| `sh docs/setup/setup-check.sh` | 2026-10-20T13:39:34Z | 2026-10-20T13:39:55Z | exit 0 |
| `git diff --check origin/main HEAD` | 2026-10-20T13:39:55Z | 2026-10-20T13:39:55Z | exit 0 |

## The green (CI on the PR head)

PR #443 at head `79169f4428582db17600716b10ea2e6f49271583`, CI run `39432743311`, job `tests (unit → integration → e2e)`. The log lines are quoted from that run:

```
unit         2026-10-20T13:42:39Z  Run go test -timeout 10m ./...
integration  2026-10-20T13:43:10Z  Run go test -tags=integration -timeout 20m ./...
integration  2026-10-20T13:53:57Z  ok  	github.com/pharzam/layup/internal/run	588.190s
end-to-end   2026-10-20T13:53:57Z  Run go test -tags=e2e -timeout 10m ./...
```

- The integration step permitted each test binary up to 20 minutes. `internal/run` passed uncached in **588.190 s**, below the 1,200 s limit. That is the same product code that timed out at 600 s on `2866a3f`: this head changes no Go file.
- Unit and end-to-end kept `-timeout 10m` and passed.
- 13 check runs: 12 `success` and `review-record` `failure`. The failure was expected: the first review record of #442 comes after this head, and the check is read again on the frozen head.
- Recorded by the coordinator (Claude Opus 5.5), from `gh run view 39432743311 --log` and `gh api repos/pharzam/layup/commits/79169f4…/check-runs`.
