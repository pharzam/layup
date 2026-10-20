# Test runs of T-b6fn (#440)

Every command below was run on this host (`sh` and `bash --posix`), with start and end times taken directly from `date -u +%FT%TZ` and exit statuses read directly from `$?`.

## Red-then-green progression by stage (`sh docs/tests/freeze-check/run-all.sh`)

| Stage | Commit | Start (UTC) | End (UTC) | Exit | Result and per-case diagnostics |
| --- | --- | --- | --- | --- | --- |
| Stage 0: exit-0 stub + runner + 6 frame cases | `a8f57c8` | `2026-10-20T14:17:06Z` | `2026-10-20T14:17:06Z` | `1` | `0 passed, 13 failed`: all 6 frame cases (`bad-base-not-a-commit`, `bad-not-a-worktree`, `bad-usage-bad-budget`, `bad-usage-missing-budget`, `bad-usage-no-args`, `bad-usage-unknown-flag`) RED (`exit 0 (wanted 2)`), plus 6 `K1`–`K6` floor failures and 1 `crlf` floor failure. |
| Stage 1: argument parsing and environment checks | `6e6e2f0` | `2026-10-20T14:18:04Z` | `2026-10-20T14:18:04Z` | `1` | `6 passed, 7 failed`: all 6 frame cases GREEN (`ok`); only the 6 `K1`–`K6` floor checks and 1 `crlf` floor check remain red until Stage 2 adds those fixtures. |
| Stage 2: six-check `ok` stub + `K1`–`K6` and `crlf` fixtures | `8db1ad7` | `2026-10-20T14:24:54Z` | `2026-10-20T14:24:56Z` | `1` | `14 passed, 28 failed`: the 6 frame cases and all 8 `good-*` cases (`good-K1-clean`, `good-K2-registered`, `good-K3-alternatives`, `good-K4-forms`, `good-K5-at-max`, `good-K6-no-conditions`, `good-K6-valid`, `good-crlf-all`) GREEN; all 28 `bad-K*` and `bad-crlf-conditions` cases RED (`exit 0 (wanted 1)`, getting `ok K<n>` in place of their expected `FAIL K<n>` line). |
| Stage 3, K1: white space | `a7265db` | `2026-10-20T14:25:50Z` | `2026-10-20T14:25:52Z` | `1` | `16 passed, 26 failed`: `bad-K1-blank-eof` and `bad-K1-trailing-space` turn GREEN; the remaining 26 `bad-*` cases stay RED. |
| Stage 3, K2: registration of new tests | `eec90f1` | `2026-10-20T14:27:07Z` | `2026-10-20T14:27:10Z` | `1` | `20 passed, 22 failed`: `bad-K2-missing-file`, `bad-K2-missing-prd`, `bad-K2-missing-traceability` and `bad-K2-two-missing` turn GREEN; the remaining 22 `bad-*` cases stay RED. |
| Stage 3, K3: test selectors | `d9a4ba3` | `2026-10-20T14:29:17Z` | `2026-10-20T14:29:20Z` | `1` | `27 passed, 15 failed`: `bad-K3-double-quoted`, `bad-K3-no-match`, `bad-K3-other-package`, `bad-K3-other-tag`, `bad-K3-subtest`, `bad-K3-unquoted` and `bad-K3-unsupported-form` turn GREEN; the remaining 15 `bad-*` cases stay RED. |
| Stage 3, K4: restore evidence of mutants | `0664f2d` | `2026-10-20T14:30:38Z` | `2026-10-20T14:30:41Z` | `1` | `30 passed, 12 failed`: `bad-K4-form-a-different`, `bad-K4-result-only` (with `skip K4` lines for `SM19a`, `SM29`, `SM30` and `FAIL K4` for `SM31`) and `bad-K4-two-hashes-no-restore` turn GREEN; the remaining 12 `bad-*` cases stay RED. |
| Stage 3, K5: budget at landing | `c07b7f8` | `2026-10-20T14:31:37Z` | `2026-10-20T14:31:40Z` | `1` | `33 passed, 9 failed`: `bad-K5-default-closeout`, `bad-K5-over-files` and `bad-K5-over-lines` turn GREEN; the remaining 9 `bad-*` cases stay RED. |
| Stage 3, K6: the conditions table | `1afd9f2` | `2026-10-20T14:32:58Z` | `2026-10-20T14:33:02Z` | `0` | `42 passed, 0 failed`: `bad-K6-empty-table`, `bad-K6-line-past-end`, `bad-K6-missing-commit`, `bad-K6-missing-file`, `bad-K6-missing-path`, `bad-K6-missing-test`, `bad-K6-mixed-valid-invalid`, `bad-K6-no-span` and `bad-crlf-conditions` turn GREEN. `bash --posix docs/tests/freeze-check/run-all.sh` also exited `0` (`42 passed, 0 failed`). |

## S3 real-head replays

Recorded with full SHAs, bases, commands, timestamps and outputs in [`runs/T-b6fn/real-heads.md`](real-heads.md).

## Gates and pre-freeze check on `T-b6fn`

| Command | Start (UTC) | End (UTC) | Exit | Output summary |
| --- | --- | --- | --- | --- |
| `sh docs/adr/adr-lint.sh` | `2026-10-20T14:45:48Z` | `2026-10-20T14:45:49Z` | `0` | `adr-lint: OK` |
| `sh docs/prd/prd-lint.sh` | `2026-10-20T14:45:49Z` | `2026-10-20T14:45:49Z` | `0` | `prd-lint: OK` |
| `sh docs/tests/run-discipline-tests.sh` | `2026-10-20T14:45:49Z` | `2026-10-20T14:45:52Z` | `0` | `run-discipline-tests: 138 passed, 0 failed` |
| `sh docs/tests/nested-checkout-check.sh` | `2026-10-20T14:45:52Z` | `2026-10-20T14:45:52Z` | `0` | `nested-checkout-check: OK  10 cases behaved` |
| `sh docs/tests/freeze-check/run-all.sh` | `2026-10-20T14:45:52Z` | `2026-10-20T14:45:55Z` | `0` | `freeze-check: 42 passed, 0 failed` |
| `go build ./...` | `2026-10-20T14:45:55Z` | `2026-10-20T14:45:56Z` | `0` | clean |
| `go vet ./...` | `2026-10-20T14:45:56Z` | `2026-10-20T14:45:57Z` | `0` | clean |
| `go test ./...` | `2026-10-20T14:45:57Z` | `2026-10-20T14:45:59Z` | `0` | all 21 packages `ok` |
| `sh docs/links/link-lint.sh` | `2026-10-20T14:47:10Z` | `2026-10-20T14:47:24Z` | `0` | `link-lint: OK  2904 links resolved` |
| `sh docs/setup/setup-check.sh` | `2026-10-20T14:47:24Z` | `2026-10-20T14:47:42Z` | `0` | all 13 checks (`pin` through `kit-linters`) `OK` |
| `git diff --check origin/main HEAD` | `2026-10-20T14:47:42Z` | `2026-10-20T14:47:42Z` | `0` | clean |
| `sh docs/tests/freeze-check.sh --task T-b6fn --base origin/main --budget 1300/170` | `2026-10-20T14:47:42Z` | `2026-10-20T14:47:43Z` | `0` | `ok K1`–`ok K6` (projection: figure `1178/154`, maximum `1300/170`, margin `122/16`) |

## The fix of K2 found on the merged `main` (by Claude Opus 5.5, after the Devin run ended on its daily quota)

The Devin run (Gemini 4 Argon High) ended at 2026-10-20T14:46:46Z with "Your daily usage quota has been exhausted", after commit `3e4ea34`. Its uncommitted edits, a compaction of the scripts, were set aside and not kept. The branch then merged `origin/main` (`349aef3`, with #443, #437 and #444; no conflict).

On the merged `main`, K2 failed on T-r7b9 as merged (`0246d9c..2866a3f`), because a §12 Test cell may name a test's file ("… and its 13 siblings of questionloop_test.go"). That was a false positive, so K2 now accepts the test's name or the basename of its file.

| Stage | Commit | Run | Time (UTC) | Exit | Result |
| --- | --- | --- | --- | --- | --- |
| Red: the cases `good-K2-prd-by-file` and `bad-K2-prd-other-file` on the old K2 | `ffb0a24` | `sh docs/tests/freeze-check/run-all.sh` | `2026-10-20T15:08:28Z` | `1` | `43 passed, 1 failed`: `good-K2-prd-by-file` red with `FAIL K2: internal/x/a_test.go:11: TestNew is not in … REQ-001 Test cell` |
| Green: K2 accepts the file's basename | `5cf545b` | the same, and `bash --posix` | `2026-10-20T15:08:52Z` | `0` | `44 passed, 0 failed` |

## The gates on the head before the freeze

| Command | Start (UTC) | Exit | Last line |
| --- | --- | --- | --- |
| `sh docs/adr/adr-lint.sh` | `2026-10-20T15:10:32Z` | `0` | adr-lint: OK |
| `sh docs/prd/prd-lint.sh` | `2026-10-20T15:10:33Z` | `0` | prd-lint: OK |
| `sh docs/links/link-lint.sh` | `2026-10-20T15:10:33Z` | `0` | link-lint: OK  2948 links resolved |
| `sh docs/tests/run-discipline-tests.sh` | `2026-10-20T15:10:46Z` | `0` | run-discipline-tests: 138 passed, 0 failed |
| `sh docs/tests/nested-checkout-check.sh` | `2026-10-20T15:10:49Z` | `0` | nested-checkout-check: OK  10 cases behaved |
| `sh docs/tests/freeze-check/run-all.sh` | `2026-10-20T15:10:49Z` | `0` | freeze-check: 44 passed, 0 failed |
| `bash --posix docs/tests/freeze-check/run-all.sh` | `2026-10-20T15:10:53Z` | `0` | freeze-check: 44 passed, 0 failed |
| `sh docs/setup/setup-check.sh` | `2026-10-20T15:10:57Z` | `0` | setup-check: kit-linters OK |
| `go build ./...` | `2026-10-20T15:11:19Z` | `0` |  |
| `go vet ./...` | `2026-10-20T15:11:19Z` | `0` |  |
| `go test ./...` | `2026-10-20T15:11:19Z` | `0` | ok  	github.com/pharzam/layup/internal/work	(cached) |

## The check on this branch itself

`sh docs/tests/freeze-check.sh --task T-b6fn --base origin/main --budget 1300/170`, at 2026-10-20T15:11:20Z, with `runs/T-b6fn/conditions.md` staged, exited `1`:

```
ok K1: white space
ok K2: registration of new tests
ok K3: test selectors
ok K4: restore evidence of mutants
FAIL K5: -:0: projection: figure 1684/162, maximum 1300/170, margin -384/8 (the figure at landing is still measured at landing)
ok K6: conditions table
```

K5 is a true finding. The branch is over its maximum (the script is 575 lines against the 320 of the plan, the fixtures 840 against 520), and the overrun is reported on #440 for the Operator's figure at landing.

## The fixes of round 1 (cycle 1 of cap 2; by Claude Opus 5.5)

Round 1 (GPT-6.1 Sol) found five material findings. Each of F1 to F4 has a fixture written red before its code:

| Finding | Fixture | Fix |
| --- | --- | --- |
| F1: an unreadable traceability or PRD file passed K2 | `bad-env-K2-unreadable`, `bad-env-K2-unreadable-prd` (`run.sh` makes the listed paths mode 000 after the commit, and refuses the case as root) | a read check before each input, and the status of the K2 parser propagated: exit 2, "cannot read …" |
| F2: a word of a mutant's description ("rule dropped") exempted a killed mutant | `bad-K4-dropped-in-description` (a hash-less killed row fails; "dropped as a duplicate" and "not run" still skip) | an exemption only from "not run" in a cell, or a cell that begins with "dropped" or "superseded"; a cell `KILLED` overrides both |
| F3: a backtick elsewhere on a line hid a `go test` command | `bad-K3-backtick-outside` | the text outside the spans is read, and fenced lines are read as they stand. The case's line ends in a comment, so it fails "unsupported selector form", one of the two failures the finding's falsifier names |
| F3, found while fixing it: an inline span that crosses a line break was paired wrongly (T-r7b9's `mutants.md:5-6` as merged) | `bad-K3-span-across-lines` (a bad selector in such a span was missed), `good-K3-span-across-lines` | spans are paired across the lines of a paragraph, and the joined span is checked at its first line |
| F4: an untracked file supplied K6's evidence | `bad-K6-untracked-evidence` (`run.sh` lays `work/` over the tree after the commit, untracked) | K6 reads its table and every cited path at `HEAD` (`git cat-file`, `git show`) |
| F5: the lesson not written back | — | `docs/guardrails.md` §2, "A check of the records that passes its own fixtures and fails the records of `main`" |

| Run | Commit or tree | Time (UTC) | Exit | Result |
| --- | --- | --- | --- | --- |
| Red of F1–F4 | `6b0ed68` | `2026-10-20T15:24:43Z` | `1` | `44 passed, 5 failed`: the five new cases |
| Red of the span cases | `0bb67fc` | `2026-10-20T15:30:26Z` | `1` | `45 passed, 6 failed`: the five above and `bad-K3-span-across-lines` |
| Green | `0d3fdd1` | `2026-10-20T15:27:59Z` and after | `0` | `51 passed, 0 failed`, under `sh` and `bash --posix` |

`bad-K3-backtick-outside` was red with an `EXPECT` that wanted "matches no test". The fix gives "unsupported selector form" for that line, because of its trailing `# see` comment, and the `EXPECT` was set to that in `0d3fdd1`.

**Merged tasks of `main`**, with `--no-conditions --closeout 0/0`:
- T-05f4 at `99ff862` (base `699d868`): exit 0.
- T-r7b9 at `2866a3f` (base `0246d9c`): one line, `FAIL K3: runs/T-r7b9/mutants.md:5: quote the selector with single quotes`. The prose template there, a command whose selector is a placeholder, spans two lines, and the K3 grammar refuses a placeholder selector, as the round noted for T-h65e. The old pairing missed that line. The new pitfall says how to write such a template.

## The gates after the fixes

| Command | Start (UTC) | Exit | Last line |
| --- | --- | --- | --- |
| `sh docs/adr/adr-lint.sh` | `2026-10-20T15:31:06Z` | `0` | adr-lint: OK |
| `sh docs/prd/prd-lint.sh` | `2026-10-20T15:31:07Z` | `0` | prd-lint: OK |
| `sh docs/links/link-lint.sh` | `2026-10-20T15:31:07Z` | `0` | link-lint: OK  2948 links resolved |
| `sh docs/tests/run-discipline-tests.sh` | `2026-10-20T15:31:19Z` | `0` | run-discipline-tests: 138 passed, 0 failed |
| `sh docs/tests/nested-checkout-check.sh` | `2026-10-20T15:31:22Z` | `0` | nested-checkout-check: OK  10 cases behaved |
| `sh docs/tests/freeze-check/run-all.sh` | `2026-10-20T15:31:23Z` | `0` | freeze-check: 51 passed, 0 failed |
| `bash --posix docs/tests/freeze-check/run-all.sh` | `2026-10-20T15:31:28Z` | `0` | freeze-check: 51 passed, 0 failed |
| `sh docs/setup/setup-check.sh` | `2026-10-20T15:31:32Z` | `0` | setup-check: kit-linters OK |
| `go build ./...` | `2026-10-20T15:31:54Z` | `0` |  |
| `go vet ./...` | `2026-10-20T15:31:54Z` | `0` |  |
| `go test ./...` | `2026-10-20T15:31:54Z` | `0` | ok  	github.com/pharzam/layup/internal/work	(cached) |
