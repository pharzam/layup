# T-c826 S1 test runs

Stage A of amendment 1. No checker, runner, workflow or canonical-text implementation
is included. The binding review sets 1,300 lines / 170 files, close-out inside,
and cap 2. The supplied Operator decisions authorize two goal classes in one task
and the checker-only, numbered-PR exception; removal issue #452 (T-szvp).

## Tested source and recording method

- Committed baseline and tested HEAD: `d27b33180a1aa82c23f3ab66e757e06de356eb35`.
- Tests are the uncommitted S1 inputs and EXPECT files in this worktree.
- Every timestamp below was read directly with `date -u +%FT%TZ` immediately
  before and after the named command. Status is the direct subprocess return code,
  with stdout and stderr captured together, never the status of a pipe.
- All disposable repositories and scratch output used
  `TMPDIR=/home/layup/.cache/layup-session/tmp-green/T-c826-stage-a`.
- A fixture harness status of 2 means an EXPECT mismatch. Its checker status is
  read from the runner diagnostic; a matching control includes the checker status
  in the runner's full-output/status comparison. Suite status 1 means failed cases.
- The actual migration PR number is not supplied for Stage A. The task check runs
  its general assertions; exact-number assertions remain for the coordinating
  session before S2. No PR number is invented.

| Unchanged producing file | SHA-256 before and after S1 |
| --- | --- |
| `docs/tests/freeze-check.sh` | `a847cc17bffd2ea5957f65159068d595e43b34ac5a43adc5113c06355a522b63` |
| `docs/tests/freeze-check/run.sh` | `9443c510830f705fdc3f81732270db9f406cb0302cb2f27e18b414392899a4b0` |
| `docs/tests/freeze-check/run-all.sh` | `b892114bdc1719b724c46ba2242941e1f8bcf00d3945ce95877dd265e7a26bb8` |
| `.github/workflows/ci.yml` | `a94883b91635695b31812843f22c524cb3094705febc498177042c174e37a8d9` |

## Requirement-driven replacements (R8, R9)

ADR-0028 and the supplied Operator margin decision replace the prior allowance
and direct-maximum comparison. This is a requirement change, recorded before
implementation, rather than a weakened legacy assertion.

- Rename `bad-K5-default-closeout` to `good-K5-no-closeout`: retain its six
  committed lines and assert 6/1, removing the synthetic 70/2.
- `good-K5-at-max` retains its five lines, uses maximum 5/1 and drops --closeout.
  The committed figure equals the maximum; the margin now includes tolerance.
- `bad-K5-over-lines` uses 161 committed lines against maximum 100 (limit 160).
  `good-K5-line-limit` supplies exactly 160, so changing > to >= is observable.
- `bad-K5-over-files` grows from two to six files against maximum 1 (limit 5).
  `good-K5-file-limit` supplies exactly five files with lines within their limit.
- `bad-usage-closeout` explicitly rejects the formerly accepted option with exit 2.
- All 55 original output contracts affected by K5/usage are updated: 54 retain
  their names and one is renamed. This includes both CRLF expectations, whose
  carriage returns are retained. The seven unchanged frame cases remain controls.
- Nineteen added cases plus the renamed case yield 81 cases. The fractional
  percentage case has one line, maximum 309, tolerance 61 and limit 370;
  no 370- or 371-line fixture body is added.

## Coverage reserved for independent mutations in S5

| Rule | Named case(s) |
| --- | --- |
| Each of .md, .tsv and .txt; repeated declarations; exclusion from lines and files | `good-K5-evidence-kinds`: gross 6/5 exceeds the file limit 4; net 2/2 passes. |
| Audited record excluded; script, execution data and dual-use record retained | `good-K5-evidence-inputs`: net 6/3, including all three unlisted consumers/inputs. |
| Added, modified and deleted committed endpoints | `good-K5-evidence-added`, `good-K5-evidence-modified`, `good-K5-evidence-deleted`; each net 1/1. Modified record replaces all three old lines with one, requiring removal and addition subtraction. |
| Exact task directory including prefix neighbour | `bad-usage-evidence-other-task`, `bad-usage-evidence-task-prefix`. |
| Traversal, suffix, committed presence, missing option argument | `bad-usage-evidence-traversal`, `bad-usage-evidence-suffix`, `bad-usage-evidence-missing`, `bad-usage-evidence-no-value`. |
| Literal selection | `good-K5-evidence-literal`: audit.md.txt and auditXmd remain counted beside declared audit.md. Prefix and unescaped-regex selection change net 3/2. |
| Duplicate declaration subtracts once | `good-K5-evidence-duplicate`: repeat report.md; net 1/1. |
| Floor, percentage, ceiling, integer truncation | `good-K5-line-limit`, `good-K5-percentage`, `good-K5-ceiling`, `good-K5-truncation`. |
| Inclusive line boundary and beyond; positive, zero and negative line margin | `good-K5-at-max`, `good-K5-line-limit`, `bad-K5-over-lines`. |
| File +4, inclusive boundary and beyond; positive, zero and negative file margin | `good-K5-at-max`, `good-K5-file-limit`, `bad-K5-over-files`. |
| Synthetic line and file allowances independently absent | `good-K5-no-closeout` complete output. |
| --closeout acceptance removed | `bad-usage-closeout`. |
| Canonical sentence, runner preservation, restore list, guarded one-path/one-number exception, removal issue | `runs/T-c826/checks.sh`; exact identity requires the actual PR number before S2. |

No validator, exclusion or mutant is credited by a baseline unsupported-flag
diagnostic. Each requires the later green and independent mutation evidence.
No mode, symlink, runner-metadata, SHA-pin or restore-replay case is included.

## Baseline control before test edits

| Command | Start UTC | End UTC | Exit | Observed result |
| --- | --- | --- | --- | --- |
| `sh docs/tests/freeze-check/run-all.sh` | 2026-10-20T20:51:59Z | 2026-10-20T20:52:05Z | 0 | freeze-check: 62 passed, 0 failed |

## Every individual case after writing S1

All 74 changed/new expectations mismatch for the intended new contract. All seven
unchanged controls match. No case fails to construct its disposable repository.

| Command | Start UTC | End UTC | Checker exit | Harness exit | Result and reason |
| --- | --- | --- | --- | --- | --- |
| `sh docs/tests/freeze-check/run.sh docs/tests/freeze-check/bad-K1-blank-eof` | 2026-10-20T20:55:17Z | 2026-10-20T20:55:17Z | 1 | 2 | K5 complete-output mismatch: legacy 70/2 and maximum-minus-figure margin; other check diagnostics preserved. |
| `sh docs/tests/freeze-check/run.sh docs/tests/freeze-check/bad-K1-trailing-space` | 2026-10-20T20:55:17Z | 2026-10-20T20:55:17Z | 1 | 2 | K5 complete-output mismatch: legacy 70/2 and maximum-minus-figure margin; other check diagnostics preserved. |
| `sh docs/tests/freeze-check/run.sh docs/tests/freeze-check/bad-K2-missing-file` | 2026-10-20T20:55:17Z | 2026-10-20T20:55:17Z | 1 | 2 | K5 complete-output mismatch: legacy 70/2 and maximum-minus-figure margin; other check diagnostics preserved. |
| `sh docs/tests/freeze-check/run.sh docs/tests/freeze-check/bad-K2-missing-prd` | 2026-10-20T20:55:18Z | 2026-10-20T20:55:18Z | 1 | 2 | K5 complete-output mismatch: legacy 70/2 and maximum-minus-figure margin; other check diagnostics preserved. |
| `sh docs/tests/freeze-check/run.sh docs/tests/freeze-check/bad-K2-missing-traceability` | 2026-10-20T20:55:18Z | 2026-10-20T20:55:18Z | 1 | 2 | K5 complete-output mismatch: legacy 70/2 and maximum-minus-figure margin; other check diagnostics preserved. |
| `sh docs/tests/freeze-check/run.sh docs/tests/freeze-check/bad-K2-no-prd-frame-test` | 2026-10-20T20:55:18Z | 2026-10-20T20:55:18Z | 1 | 2 | K5 complete-output mismatch: legacy 70/2 and maximum-minus-figure margin; other check diagnostics preserved. |
| `sh docs/tests/freeze-check/run.sh docs/tests/freeze-check/bad-K2-no-prd-mapped-test` | 2026-10-20T20:55:18Z | 2026-10-20T20:55:18Z | 1 | 2 | K5 complete-output mismatch: legacy 70/2 and maximum-minus-figure margin; other check diagnostics preserved. |
| `sh docs/tests/freeze-check/run.sh docs/tests/freeze-check/bad-K2-no-s12-frame-test` | 2026-10-20T20:55:18Z | 2026-10-20T20:55:18Z | 1 | 2 | K5 complete-output mismatch: legacy 70/2 and maximum-minus-figure margin; other check diagnostics preserved. |
| `sh docs/tests/freeze-check/run.sh docs/tests/freeze-check/bad-K2-no-s12-mapped-test` | 2026-10-20T20:55:18Z | 2026-10-20T20:55:18Z | 1 | 2 | K5 complete-output mismatch: legacy 70/2 and maximum-minus-figure margin; other check diagnostics preserved. |
| `sh docs/tests/freeze-check/run.sh docs/tests/freeze-check/bad-K2-prd-other-file` | 2026-10-20T20:55:18Z | 2026-10-20T20:55:18Z | 1 | 2 | K5 complete-output mismatch: legacy 70/2 and maximum-minus-figure margin; other check diagnostics preserved. |
| `sh docs/tests/freeze-check/run.sh docs/tests/freeze-check/bad-K2-two-missing` | 2026-10-20T20:55:18Z | 2026-10-20T20:55:18Z | 1 | 2 | K5 complete-output mismatch: legacy 70/2 and maximum-minus-figure margin; other check diagnostics preserved. |
| `sh docs/tests/freeze-check/run.sh docs/tests/freeze-check/bad-K3-backtick-outside` | 2026-10-20T20:55:18Z | 2026-10-20T20:55:19Z | 1 | 2 | K5 complete-output mismatch: legacy 70/2 and maximum-minus-figure margin; other check diagnostics preserved. |
| `sh docs/tests/freeze-check/run.sh docs/tests/freeze-check/bad-K3-double-quoted` | 2026-10-20T20:55:19Z | 2026-10-20T20:55:19Z | 1 | 2 | K5 complete-output mismatch: legacy 70/2 and maximum-minus-figure margin; other check diagnostics preserved. |
| `sh docs/tests/freeze-check/run.sh docs/tests/freeze-check/bad-K3-no-match` | 2026-10-20T20:55:19Z | 2026-10-20T20:55:19Z | 1 | 2 | K5 complete-output mismatch: legacy 70/2 and maximum-minus-figure margin; other check diagnostics preserved. |
| `sh docs/tests/freeze-check/run.sh docs/tests/freeze-check/bad-K3-other-package` | 2026-10-20T20:55:19Z | 2026-10-20T20:55:19Z | 1 | 2 | K5 complete-output mismatch: legacy 70/2 and maximum-minus-figure margin; other check diagnostics preserved. |
| `sh docs/tests/freeze-check/run.sh docs/tests/freeze-check/bad-K3-other-tag` | 2026-10-20T20:55:19Z | 2026-10-20T20:55:19Z | 1 | 2 | K5 complete-output mismatch: legacy 70/2 and maximum-minus-figure margin; other check diagnostics preserved. |
| `sh docs/tests/freeze-check/run.sh docs/tests/freeze-check/bad-K3-span-across-lines` | 2026-10-20T20:55:19Z | 2026-10-20T20:55:19Z | 1 | 2 | K5 complete-output mismatch: legacy 70/2 and maximum-minus-figure margin; other check diagnostics preserved. |
| `sh docs/tests/freeze-check/run.sh docs/tests/freeze-check/bad-K3-subtest` | 2026-10-20T20:55:19Z | 2026-10-20T20:55:19Z | 1 | 2 | K5 complete-output mismatch: legacy 70/2 and maximum-minus-figure margin; other check diagnostics preserved. |
| `sh docs/tests/freeze-check/run.sh docs/tests/freeze-check/bad-K3-unclosed-span-eof` | 2026-10-20T20:55:19Z | 2026-10-20T20:55:19Z | 1 | 2 | K5 complete-output mismatch: legacy 70/2 and maximum-minus-figure margin; other check diagnostics preserved. |
| `sh docs/tests/freeze-check/run.sh docs/tests/freeze-check/bad-K3-unclosed-span-paragraph` | 2026-10-20T20:55:19Z | 2026-10-20T20:55:19Z | 1 | 2 | K5 complete-output mismatch: legacy 70/2 and maximum-minus-figure margin; other check diagnostics preserved. |
| `sh docs/tests/freeze-check/run.sh docs/tests/freeze-check/bad-K3-unquoted` | 2026-10-20T20:55:19Z | 2026-10-20T20:55:19Z | 1 | 2 | K5 complete-output mismatch: legacy 70/2 and maximum-minus-figure margin; other check diagnostics preserved. |
| `sh docs/tests/freeze-check/run.sh docs/tests/freeze-check/bad-K3-unsupported-form` | 2026-10-20T20:55:19Z | 2026-10-20T20:55:20Z | 1 | 2 | K5 complete-output mismatch: legacy 70/2 and maximum-minus-figure margin; other check diagnostics preserved. |
| `sh docs/tests/freeze-check/run.sh docs/tests/freeze-check/bad-K4-dropped-in-description` | 2026-10-20T20:55:20Z | 2026-10-20T20:55:20Z | 1 | 2 | K5 complete-output mismatch: legacy 70/2 and maximum-minus-figure margin; other check diagnostics preserved. |
| `sh docs/tests/freeze-check/run.sh docs/tests/freeze-check/bad-K4-form-a-different` | 2026-10-20T20:55:20Z | 2026-10-20T20:55:20Z | 1 | 2 | K5 complete-output mismatch: legacy 70/2 and maximum-minus-figure margin; other check diagnostics preserved. |
| `sh docs/tests/freeze-check/run.sh docs/tests/freeze-check/bad-K4-form-b-equal-hashes` | 2026-10-20T20:55:20Z | 2026-10-20T20:55:20Z | 1 | 2 | K5 complete-output mismatch: legacy 70/2 and maximum-minus-figure margin; other check diagnostics preserved. |
| `sh docs/tests/freeze-check/run.sh docs/tests/freeze-check/bad-K4-result-only` | 2026-10-20T20:55:20Z | 2026-10-20T20:55:20Z | 1 | 2 | K5 complete-output mismatch: legacy 70/2 and maximum-minus-figure margin; other check diagnostics preserved. |
| `sh docs/tests/freeze-check/run.sh docs/tests/freeze-check/bad-K4-two-hashes-no-restore` | 2026-10-20T20:55:20Z | 2026-10-20T20:55:20Z | 1 | 2 | K5 complete-output mismatch: legacy 70/2 and maximum-minus-figure margin; other check diagnostics preserved. |
| `sh docs/tests/freeze-check/run.sh docs/tests/freeze-check/bad-K5-over-files` | 2026-10-20T20:55:20Z | 2026-10-20T20:55:20Z | 1 | 2 | 6 files must fail against limit 5, margin -1; baseline reports synthetic 8 and omits tolerance/limit. |
| `sh docs/tests/freeze-check/run.sh docs/tests/freeze-check/bad-K5-over-lines` | 2026-10-20T20:55:20Z | 2026-10-20T20:55:20Z | 1 | 2 | 161 lines must fail against limit 160, margin -1; baseline reports synthetic 231 and omits tolerance/limit. |
| `sh docs/tests/freeze-check/run.sh docs/tests/freeze-check/bad-K6-empty-table` | 2026-10-20T20:55:20Z | 2026-10-20T20:55:20Z | 1 | 2 | K5 complete-output mismatch: legacy 70/2 and maximum-minus-figure margin; other check diagnostics preserved. |
| `sh docs/tests/freeze-check/run.sh docs/tests/freeze-check/bad-K6-line-past-end` | 2026-10-20T20:55:20Z | 2026-10-20T20:55:21Z | 1 | 2 | K5 complete-output mismatch: legacy 70/2 and maximum-minus-figure margin; other check diagnostics preserved. |
| `sh docs/tests/freeze-check/run.sh docs/tests/freeze-check/bad-K6-line-zero` | 2026-10-20T20:55:21Z | 2026-10-20T20:55:21Z | 1 | 2 | K5 complete-output mismatch: legacy 70/2 and maximum-minus-figure margin; other check diagnostics preserved. |
| `sh docs/tests/freeze-check/run.sh docs/tests/freeze-check/bad-K6-missing-commit` | 2026-10-20T20:55:21Z | 2026-10-20T20:55:21Z | 1 | 2 | K5 complete-output mismatch: legacy 70/2 and maximum-minus-figure margin; other check diagnostics preserved. |
| `sh docs/tests/freeze-check/run.sh docs/tests/freeze-check/bad-K6-missing-file` | 2026-10-20T20:55:21Z | 2026-10-20T20:55:21Z | 1 | 2 | K5 complete-output mismatch: legacy 70/2 and maximum-minus-figure margin; other check diagnostics preserved. |
| `sh docs/tests/freeze-check/run.sh docs/tests/freeze-check/bad-K6-missing-path` | 2026-10-20T20:55:21Z | 2026-10-20T20:55:21Z | 1 | 2 | K5 complete-output mismatch: legacy 70/2 and maximum-minus-figure margin; other check diagnostics preserved. |
| `sh docs/tests/freeze-check/run.sh docs/tests/freeze-check/bad-K6-missing-test` | 2026-10-20T20:55:21Z | 2026-10-20T20:55:21Z | 1 | 2 | K5 complete-output mismatch: legacy 70/2 and maximum-minus-figure margin; other check diagnostics preserved. |
| `sh docs/tests/freeze-check/run.sh docs/tests/freeze-check/bad-K6-mixed-valid-invalid` | 2026-10-20T20:55:21Z | 2026-10-20T20:55:21Z | 1 | 2 | K5 complete-output mismatch: legacy 70/2 and maximum-minus-figure margin; other check diagnostics preserved. |
| `sh docs/tests/freeze-check/run.sh docs/tests/freeze-check/bad-K6-no-span` | 2026-10-20T20:55:21Z | 2026-10-20T20:55:21Z | 1 | 2 | K5 complete-output mismatch: legacy 70/2 and maximum-minus-figure margin; other check diagnostics preserved. |
| `sh docs/tests/freeze-check/run.sh docs/tests/freeze-check/bad-K6-range-zero` | 2026-10-20T20:55:21Z | 2026-10-20T20:55:21Z | 1 | 2 | K5 complete-output mismatch: legacy 70/2 and maximum-minus-figure margin; other check diagnostics preserved. |
| `sh docs/tests/freeze-check/run.sh docs/tests/freeze-check/bad-base-not-a-commit` | 2026-10-20T20:55:21Z | 2026-10-20T20:55:21Z | 2 | 0 | Unchanged frame assertion; green regression control. |
| `sh docs/tests/freeze-check/run.sh docs/tests/freeze-check/bad-crlf-conditions` | 2026-10-20T20:55:21Z | 2026-10-20T20:55:22Z | 1 | 2 | K5 complete-output mismatch: legacy 70/2 and maximum-minus-figure margin; other check diagnostics preserved. |
| `sh docs/tests/freeze-check/run.sh docs/tests/freeze-check/bad-env-K2-unreadable` | 2026-10-20T20:55:22Z | 2026-10-20T20:55:22Z | 2 | 0 | Unchanged frame assertion; green regression control. |
| `sh docs/tests/freeze-check/run.sh docs/tests/freeze-check/bad-env-K2-unreadable-prd` | 2026-10-20T20:55:22Z | 2026-10-20T20:55:22Z | 2 | 0 | Unchanged frame assertion; green regression control. |
| `sh docs/tests/freeze-check/run.sh docs/tests/freeze-check/bad-env-K3-unreadable-test` | 2026-10-20T20:55:22Z | 2026-10-20T20:55:22Z | 2 | 0 | Unchanged frame assertion; green regression control. |
| `sh docs/tests/freeze-check/run.sh docs/tests/freeze-check/bad-env-dirty-worktree` | 2026-10-20T20:55:22Z | 2026-10-20T20:55:22Z | 2 | 0 | Unchanged frame assertion; green regression control. |
| `sh docs/tests/freeze-check/run.sh docs/tests/freeze-check/bad-env-untracked-evidence` | 2026-10-20T20:55:22Z | 2026-10-20T20:55:22Z | 2 | 0 | Unchanged frame assertion; green regression control. |
| `sh docs/tests/freeze-check/run.sh docs/tests/freeze-check/bad-not-a-worktree` | 2026-10-20T20:55:22Z | 2026-10-20T20:55:22Z | 2 | 0 | Unchanged frame assertion; green regression control. |
| `sh docs/tests/freeze-check/run.sh docs/tests/freeze-check/bad-usage-bad-budget` | 2026-10-20T20:55:22Z | 2026-10-20T20:55:22Z | 2 | 2 | Checker 2 is preserved; complete usage text still advertises --closeout and lacks --evidence. |
| `sh docs/tests/freeze-check/run.sh docs/tests/freeze-check/bad-usage-closeout` | 2026-10-20T20:55:22Z | 2026-10-20T20:55:22Z | 0 | 2 | Baseline accepts retired --closeout (checker 0); expected usage and checker 2. |
| `sh docs/tests/freeze-check/run.sh docs/tests/freeze-check/bad-usage-evidence-missing` | 2026-10-20T20:55:22Z | 2026-10-20T20:55:22Z | 2 | 2 | Baseline rejects unsupported --evidence with its old usage text; initial red only, not validation or exclusion proof. |
| `sh docs/tests/freeze-check/run.sh docs/tests/freeze-check/bad-usage-evidence-no-value` | 2026-10-20T20:55:22Z | 2026-10-20T20:55:22Z | 2 | 2 | Baseline rejects unsupported --evidence with its old usage text; initial red only, not validation or exclusion proof. |
| `sh docs/tests/freeze-check/run.sh docs/tests/freeze-check/bad-usage-evidence-other-task` | 2026-10-20T20:55:22Z | 2026-10-20T20:55:22Z | 2 | 2 | Baseline rejects unsupported --evidence with its old usage text; initial red only, not validation or exclusion proof. |
| `sh docs/tests/freeze-check/run.sh docs/tests/freeze-check/bad-usage-evidence-suffix` | 2026-10-20T20:55:22Z | 2026-10-20T20:55:22Z | 2 | 2 | Baseline rejects unsupported --evidence with its old usage text; initial red only, not validation or exclusion proof. |
| `sh docs/tests/freeze-check/run.sh docs/tests/freeze-check/bad-usage-evidence-task-prefix` | 2026-10-20T20:55:22Z | 2026-10-20T20:55:22Z | 2 | 2 | Baseline rejects unsupported --evidence with its old usage text; initial red only, not validation or exclusion proof. |
| `sh docs/tests/freeze-check/run.sh docs/tests/freeze-check/bad-usage-evidence-traversal` | 2026-10-20T20:55:22Z | 2026-10-20T20:55:23Z | 2 | 2 | Baseline rejects unsupported --evidence with its old usage text; initial red only, not validation or exclusion proof. |
| `sh docs/tests/freeze-check/run.sh docs/tests/freeze-check/bad-usage-missing-budget` | 2026-10-20T20:55:23Z | 2026-10-20T20:55:23Z | 2 | 2 | Checker 2 is preserved; complete usage text still advertises --closeout and lacks --evidence. |
| `sh docs/tests/freeze-check/run.sh docs/tests/freeze-check/bad-usage-no-args` | 2026-10-20T20:55:23Z | 2026-10-20T20:55:23Z | 2 | 2 | Checker 2 is preserved; complete usage text still advertises --closeout and lacks --evidence. |
| `sh docs/tests/freeze-check/run.sh docs/tests/freeze-check/bad-usage-unknown-flag` | 2026-10-20T20:55:23Z | 2026-10-20T20:55:23Z | 2 | 2 | Checker 2 is preserved; complete usage text still advertises --closeout and lacks --evidence. |
| `sh docs/tests/freeze-check/run.sh docs/tests/freeze-check/good-K1-clean` | 2026-10-20T20:55:23Z | 2026-10-20T20:55:23Z | 0 | 2 | K5 complete-output mismatch: legacy 70/2 and maximum-minus-figure margin; other check diagnostics preserved. |
| `sh docs/tests/freeze-check/run.sh docs/tests/freeze-check/good-K2-prd-by-file` | 2026-10-20T20:55:23Z | 2026-10-20T20:55:23Z | 0 | 2 | K5 complete-output mismatch: legacy 70/2 and maximum-minus-figure margin; other check diagnostics preserved. |
| `sh docs/tests/freeze-check/run.sh docs/tests/freeze-check/good-K2-registered` | 2026-10-20T20:55:23Z | 2026-10-20T20:55:23Z | 0 | 2 | K5 complete-output mismatch: legacy 70/2 and maximum-minus-figure margin; other check diagnostics preserved. |
| `sh docs/tests/freeze-check/run.sh docs/tests/freeze-check/good-K3-alternatives` | 2026-10-20T20:55:23Z | 2026-10-20T20:55:23Z | 0 | 2 | K5 complete-output mismatch: legacy 70/2 and maximum-minus-figure margin; other check diagnostics preserved. |
| `sh docs/tests/freeze-check/run.sh docs/tests/freeze-check/good-K3-span-across-lines` | 2026-10-20T20:55:23Z | 2026-10-20T20:55:23Z | 0 | 2 | K5 complete-output mismatch: legacy 70/2 and maximum-minus-figure margin; other check diagnostics preserved. |
| `sh docs/tests/freeze-check/run.sh docs/tests/freeze-check/good-K4-forms` | 2026-10-20T20:55:23Z | 2026-10-20T20:55:23Z | 0 | 2 | K5 complete-output mismatch: legacy 70/2 and maximum-minus-figure margin; other check diagnostics preserved. |
| `sh docs/tests/freeze-check/run.sh docs/tests/freeze-check/good-K5-at-max` | 2026-10-20T20:55:23Z | 2026-10-20T20:55:24Z | 1 | 2 | Actual diff is exactly maximum 5/1; baseline adds 70/2 and fails. Expected positive tolerance margin. |
| `sh docs/tests/freeze-check/run.sh docs/tests/freeze-check/good-K5-ceiling` | 2026-10-20T20:55:24Z | 2026-10-20T20:55:24Z | 0 | 2 | Small figure isolates the required tolerance, limit and margin; baseline omits those fields and adds 70/2. |
| `sh docs/tests/freeze-check/run.sh docs/tests/freeze-check/good-K5-evidence-added` | 2026-10-20T20:55:24Z | 2026-10-20T20:55:24Z | 2 | 2 | Baseline rejects unsupported --evidence with its old usage text; initial red only, not validation or exclusion proof. |
| `sh docs/tests/freeze-check/run.sh docs/tests/freeze-check/good-K5-evidence-deleted` | 2026-10-20T20:55:24Z | 2026-10-20T20:55:24Z | 2 | 2 | Baseline rejects unsupported --evidence with its old usage text; initial red only, not validation or exclusion proof. |
| `sh docs/tests/freeze-check/run.sh docs/tests/freeze-check/good-K5-evidence-duplicate` | 2026-10-20T20:55:24Z | 2026-10-20T20:55:24Z | 2 | 2 | Baseline rejects unsupported --evidence with its old usage text; initial red only, not validation or exclusion proof. |
| `sh docs/tests/freeze-check/run.sh docs/tests/freeze-check/good-K5-evidence-inputs` | 2026-10-20T20:55:24Z | 2026-10-20T20:55:24Z | 2 | 2 | Baseline rejects unsupported --evidence with its old usage text; initial red only, not validation or exclusion proof. |
| `sh docs/tests/freeze-check/run.sh docs/tests/freeze-check/good-K5-evidence-kinds` | 2026-10-20T20:55:24Z | 2026-10-20T20:55:24Z | 2 | 2 | Baseline rejects unsupported --evidence with its old usage text; initial red only, not validation or exclusion proof. |
| `sh docs/tests/freeze-check/run.sh docs/tests/freeze-check/good-K5-evidence-literal` | 2026-10-20T20:55:24Z | 2026-10-20T20:55:24Z | 2 | 2 | Baseline rejects unsupported --evidence with its old usage text; initial red only, not validation or exclusion proof. |
| `sh docs/tests/freeze-check/run.sh docs/tests/freeze-check/good-K5-evidence-modified` | 2026-10-20T20:55:24Z | 2026-10-20T20:55:24Z | 2 | 2 | Baseline rejects unsupported --evidence with its old usage text; initial red only, not validation or exclusion proof. |
| `sh docs/tests/freeze-check/run.sh docs/tests/freeze-check/good-K5-file-limit` | 2026-10-20T20:55:24Z | 2026-10-20T20:55:24Z | 1 | 2 | Expected inclusive 5-file limit passes with zero file margin; legacy projection fails. |
| `sh docs/tests/freeze-check/run.sh docs/tests/freeze-check/good-K5-line-limit` | 2026-10-20T20:55:24Z | 2026-10-20T20:55:24Z | 1 | 2 | Expected inclusive 160-line limit passes with zero line margin; legacy projection fails. |
| `sh docs/tests/freeze-check/run.sh docs/tests/freeze-check/good-K5-no-closeout` | 2026-10-20T20:55:24Z | 2026-10-20T20:55:25Z | 1 | 2 | Baseline adds 70/2 and fails at 76/3; expected committed 6/1, no allowance, with tolerance. |
| `sh docs/tests/freeze-check/run.sh docs/tests/freeze-check/good-K5-percentage` | 2026-10-20T20:55:25Z | 2026-10-20T20:55:25Z | 0 | 2 | Small figure isolates the required tolerance, limit and margin; baseline omits those fields and adds 70/2. |
| `sh docs/tests/freeze-check/run.sh docs/tests/freeze-check/good-K5-truncation` | 2026-10-20T20:55:25Z | 2026-10-20T20:55:25Z | 0 | 2 | Small figure isolates the required tolerance, limit and margin; baseline omits those fields and adds 70/2. |
| `sh docs/tests/freeze-check/run.sh docs/tests/freeze-check/good-K6-no-conditions` | 2026-10-20T20:55:25Z | 2026-10-20T20:55:25Z | 0 | 2 | K5 complete-output mismatch: legacy 70/2 and maximum-minus-figure margin; other check diagnostics preserved. |
| `sh docs/tests/freeze-check/run.sh docs/tests/freeze-check/good-K6-valid` | 2026-10-20T20:55:25Z | 2026-10-20T20:55:25Z | 0 | 2 | K5 complete-output mismatch: legacy 70/2 and maximum-minus-figure margin; other check diagnostics preserved. |
| `sh docs/tests/freeze-check/run.sh docs/tests/freeze-check/good-crlf-all` | 2026-10-20T20:55:25Z | 2026-10-20T20:55:25Z | 0 | 2 | K5 complete-output mismatch: legacy 70/2 and maximum-minus-figure margin; other check diagnostics preserved. |

## Suite and task-check commands

| Command | Start UTC | End UTC | Exit | Observed result |
| --- | --- | --- | --- | --- |
| `sh docs/tests/freeze-check/run-all.sh` | 2026-10-20T20:55:25Z | 2026-10-20T20:55:33Z | 1 | freeze-check: 7 passed, 74 failed |
| `sh runs/T-c826/checks.sh` | 2026-10-20T20:55:33Z | 2026-10-20T20:55:33Z | 1 | Seven intended missing-implementation assertions fail; preservation controls pass; actual PR number pending. |
| `sh -n runs/T-c826/checks.sh` | 2026-10-20T20:55:33Z | 2026-10-20T20:55:33Z | 0 | No diagnostic. |
| `git diff --check` | 2026-10-20T20:55:33Z | 2026-10-20T20:55:33Z | 0 | No diagnostic. |

The suite reports `freeze-check: 7 passed, 74 failed` with exit 1. Its 81 case
labels agree with the individual results above; coverage and CRLF floors do not
fail. Actual fixture mismatch diagnostics follow, once per individual command.
Random temporary paths in diff headers are retained as observed.

### Task-check actual output

```text
FAIL: canonical margin sentence in the tolerance section
ok: unchanged runner: docs/tests/freeze-check/run.sh
ok: unchanged runner: docs/tests/freeze-check/run-all.sh
ok: all nine baseline paths stay in the restore list
FAIL: kept checker requires the pull_request event
FAIL: one literal PR number and only the checker path can be kept
pending: supply the actual migration PR number before S2
FAIL: event guard receives the actual GitHub event name
FAIL: number guard receives the actual pull request number
FAIL: only the guarded checker skips restoration; both runners and other paths restore for every event
ok: restoration still reads the default branch of this repository
ok: default branch still comes from the repository event
FAIL: exception comment names the task, literal PR and removal issue
```

## Per-case actual mismatch diagnostics

### bad-K1-blank-eof

```text
run: bad-K1-blank-eof: exit 1 (wanted 1); diff (-EXPECT +got):
--- /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.0LwpO7Xhzg/expect	2026-10-20 23:55:17.780073451 +0300
+++ /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.0LwpO7Xhzg/got	2026-10-20 23:55:17.779966457 +0300
@@ -3,5 +3,5 @@
 ok K2: registration of new tests
 ok K3: test selectors
 ok K4: restore evidence of mutants
-ok K5: projection: figure 2/1, maximum 100/10, tolerance 60/4, limit 160/14, margin 158/13 (the figure at landing is still measured at landing)
+ok K5: projection: figure 72/3, maximum 100/10, margin 28/7 (the figure at landing is still measured at landing)
 ok K6: conditions table
```

### bad-K1-trailing-space

```text
run: bad-K1-trailing-space: exit 1 (wanted 1); diff (-EXPECT +got):
--- /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.2VzwrME7vG/expect	2026-10-20 23:55:17.888945906 +0300
+++ /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.2VzwrME7vG/got	2026-10-20 23:55:17.888142973 +0300
@@ -3,5 +3,5 @@
 ok K2: registration of new tests
 ok K3: test selectors
 ok K4: restore evidence of mutants
-ok K5: projection: figure 1/1, maximum 100/10, tolerance 60/4, limit 160/14, margin 159/13 (the figure at landing is still measured at landing)
+ok K5: projection: figure 71/3, maximum 100/10, margin 29/7 (the figure at landing is still measured at landing)
 ok K6: conditions table
```

### bad-K2-missing-file

```text
run: bad-K2-missing-file: exit 1 (wanted 1); diff (-EXPECT +got):
--- /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.nECQaW8BlZ/expect	2026-10-20 23:55:17.995222235 +0300
+++ /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.nECQaW8BlZ/got	2026-10-20 23:55:17.994074378 +0300
@@ -3,5 +3,5 @@
 FAIL K2: docs/tests/traceability.md:0: missing traceability file
 ok K3: test selectors
 ok K4: restore evidence of mutants
-ok K5: projection: figure 8/2, maximum 100/10, tolerance 60/4, limit 160/14, margin 152/12 (the figure at landing is still measured at landing)
+ok K5: projection: figure 78/4, maximum 100/10, margin 22/6 (the figure at landing is still measured at landing)
 ok K6: conditions table
```

### bad-K2-missing-prd

```text
run: bad-K2-missing-prd: exit 1 (wanted 1); diff (-EXPECT +got):
--- /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.oNyd2xTSMR/expect	2026-10-20 23:55:18.108074872 +0300
+++ /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.oNyd2xTSMR/got	2026-10-20 23:55:18.107683140 +0300
@@ -3,5 +3,5 @@
 FAIL K2: internal/x/a_test.go:11: TestNew is not in docs/prd/PRD-0001-layup.md §12 REQ-001 Test cell
 ok K3: test selectors
 ok K4: restore evidence of mutants
-ok K5: projection: figure 2/2, maximum 100/10, tolerance 60/4, limit 160/14, margin 158/12 (the figure at landing is still measured at landing)
+ok K5: projection: figure 72/4, maximum 100/10, margin 28/6 (the figure at landing is still measured at landing)
 ok K6: conditions table
```

### bad-K2-missing-traceability

```text
run: bad-K2-missing-traceability: exit 1 (wanted 1); diff (-EXPECT +got):
--- /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.LQL5kGWZdY/expect	2026-10-20 23:55:18.218045409 +0300
+++ /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.LQL5kGWZdY/got	2026-10-20 23:55:18.217251761 +0300
@@ -3,5 +3,5 @@
 FAIL K2: internal/x/a_test.go:11: TestNew is not in docs/tests/traceability.md
 ok K3: test selectors
 ok K4: restore evidence of mutants
-ok K5: projection: figure 1/1, maximum 100/10, tolerance 60/4, limit 160/14, margin 159/13 (the figure at landing is still measured at landing)
+ok K5: projection: figure 71/3, maximum 100/10, margin 29/7 (the figure at landing is still measured at landing)
 ok K6: conditions table
```

### bad-K2-no-prd-frame-test

```text
run: bad-K2-no-prd-frame-test: exit 1 (wanted 1); diff (-EXPECT +got):
--- /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.nDIq11ANef/expect	2026-10-20 23:55:18.332746971 +0300
+++ /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.nDIq11ANef/got	2026-10-20 23:55:18.332112683 +0300
@@ -3,5 +3,5 @@
 FAIL K2: docs/prd/PRD-0001-layup.md:0: missing PRD file
 ok K3: test selectors
 ok K4: restore evidence of mutants
-ok K5: projection: figure 12/3, maximum 100/10, tolerance 60/4, limit 160/14, margin 148/11 (the figure at landing is still measured at landing)
+ok K5: projection: figure 82/5, maximum 100/10, margin 18/5 (the figure at landing is still measured at landing)
 ok K6: conditions table
```

### bad-K2-no-prd-mapped-test

```text
run: bad-K2-no-prd-mapped-test: exit 1 (wanted 1); diff (-EXPECT +got):
--- /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.WQognjPdYV/expect	2026-10-20 23:55:18.445076331 +0300
+++ /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.WQognjPdYV/got	2026-10-20 23:55:18.444076327 +0300
@@ -3,5 +3,5 @@
 FAIL K2: docs/prd/PRD-0001-layup.md:0: missing PRD file
 ok K3: test selectors
 ok K4: restore evidence of mutants
-ok K5: projection: figure 12/3, maximum 100/10, tolerance 60/4, limit 160/14, margin 148/11 (the figure at landing is still measured at landing)
+ok K5: projection: figure 82/5, maximum 100/10, margin 18/5 (the figure at landing is still measured at landing)
 ok K6: conditions table
```

### bad-K2-no-s12-frame-test

```text
run: bad-K2-no-s12-frame-test: exit 1 (wanted 1); diff (-EXPECT +got):
--- /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.yhIIeoS2WZ/expect	2026-10-20 23:55:18.561085025 +0300
+++ /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.yhIIeoS2WZ/got	2026-10-20 23:55:18.560332260 +0300
@@ -3,5 +3,5 @@
 FAIL K2: docs/prd/PRD-0001-layup.md:0: missing §12 table
 ok K3: test selectors
 ok K4: restore evidence of mutants
-ok K5: projection: figure 9/3, maximum 100/10, tolerance 60/4, limit 160/14, margin 151/11 (the figure at landing is still measured at landing)
+ok K5: projection: figure 79/5, maximum 100/10, margin 21/5 (the figure at landing is still measured at landing)
 ok K6: conditions table
```

### bad-K2-no-s12-mapped-test

```text
run: bad-K2-no-s12-mapped-test: exit 1 (wanted 1); diff (-EXPECT +got):
--- /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.60o4USOwqK/expect	2026-10-20 23:55:18.678352283 +0300
+++ /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.60o4USOwqK/got	2026-10-20 23:55:18.677077336 +0300
@@ -3,5 +3,5 @@
 FAIL K2: docs/prd/PRD-0001-layup.md:0: missing §12 table
 ok K3: test selectors
 ok K4: restore evidence of mutants
-ok K5: projection: figure 9/3, maximum 100/10, tolerance 60/4, limit 160/14, margin 151/11 (the figure at landing is still measured at landing)
+ok K5: projection: figure 79/5, maximum 100/10, margin 21/5 (the figure at landing is still measured at landing)
 ok K6: conditions table
```

### bad-K2-prd-other-file

```text
run: bad-K2-prd-other-file: exit 1 (wanted 1); diff (-EXPECT +got):
--- /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.MxUTbiDvbB/expect	2026-10-20 23:55:18.796077852 +0300
+++ /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.MxUTbiDvbB/got	2026-10-20 23:55:18.795982109 +0300
@@ -3,5 +3,5 @@
 FAIL K2: internal/x/a_test.go:11: TestNew is not in docs/prd/PRD-0001-layup.md §12 REQ-001 Test cell
 ok K3: test selectors
 ok K4: restore evidence of mutants
-ok K5: projection: figure 4/3, maximum 100/10, tolerance 60/4, limit 160/14, margin 156/11 (the figure at landing is still measured at landing)
+ok K5: projection: figure 74/5, maximum 100/10, margin 26/5 (the figure at landing is still measured at landing)
 ok K6: conditions table
```

### bad-K2-two-missing

```text
run: bad-K2-two-missing: exit 1 (wanted 1); diff (-EXPECT +got):
--- /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.lCkLuV1WuR/expect	2026-10-20 23:55:18.906078328 +0300
+++ /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.lCkLuV1WuR/got	2026-10-20 23:55:18.906032139 +0300
@@ -4,5 +4,5 @@
 FAIL K2: internal/x/a_test.go:12: TestNewTwo is not in docs/tests/traceability.md
 ok K3: test selectors
 ok K4: restore evidence of mutants
-ok K5: projection: figure 2/1, maximum 100/10, tolerance 60/4, limit 160/14, margin 158/13 (the figure at landing is still measured at landing)
+ok K5: projection: figure 72/3, maximum 100/10, margin 28/7 (the figure at landing is still measured at landing)
 ok K6: conditions table
```

### bad-K3-backtick-outside

```text
run: bad-K3-backtick-outside: exit 1 (wanted 1); diff (-EXPECT +got):
--- /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.p32zh7gIC3/expect	2026-10-20 23:55:19.014336775 +0300
+++ /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.p32zh7gIC3/got	2026-10-20 23:55:19.013078791 +0300
@@ -3,5 +3,5 @@
 ok K2: registration of new tests
 FAIL K3: runs/T-0000/test-runs.md:4: unsupported selector form
 ok K4: restore evidence of mutants
-ok K5: projection: figure 4/1, maximum 100/10, tolerance 60/4, limit 160/14, margin 156/13 (the figure at landing is still measured at landing)
+ok K5: projection: figure 74/3, maximum 100/10, margin 26/7 (the figure at landing is still measured at landing)
 ok K6: conditions table
```

### bad-K3-double-quoted

```text
run: bad-K3-double-quoted: exit 1 (wanted 1); diff (-EXPECT +got):
--- /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.V2ligTQtDd/expect	2026-10-20 23:55:19.121330657 +0300
+++ /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.V2ligTQtDd/got	2026-10-20 23:55:19.120079255 +0300
@@ -3,5 +3,5 @@
 ok K2: registration of new tests
 FAIL K3: runs/T-0000/test-runs.md:3: quote the selector with single quotes
 ok K4: restore evidence of mutants
-ok K5: projection: figure 2/1, maximum 100/10, tolerance 60/4, limit 160/14, margin 158/13 (the figure at landing is still measured at landing)
+ok K5: projection: figure 72/3, maximum 100/10, margin 28/7 (the figure at landing is still measured at landing)
 ok K6: conditions table
```

### bad-K3-no-match

```text
run: bad-K3-no-match: exit 1 (wanted 1); diff (-EXPECT +got):
--- /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.orWxUhEEzZ/expect	2026-10-20 23:55:19.226079714 +0300
+++ /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.orWxUhEEzZ/got	2026-10-20 23:55:19.225079709 +0300
@@ -3,5 +3,5 @@
 ok K2: registration of new tests
 FAIL K3: runs/T-0000/test-runs.md:3: selector TestNoSuch matches no test in the named packages
 ok K4: restore evidence of mutants
-ok K5: projection: figure 2/1, maximum 100/10, tolerance 60/4, limit 160/14, margin 158/13 (the figure at landing is still measured at landing)
+ok K5: projection: figure 72/3, maximum 100/10, margin 28/7 (the figure at landing is still measured at landing)
 ok K6: conditions table
```

### bad-K3-other-package

```text
run: bad-K3-other-package: exit 1 (wanted 1); diff (-EXPECT +got):
--- /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.0A6CufCSIi/expect	2026-10-20 23:55:19.335209256 +0300
+++ /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.0A6CufCSIi/got	2026-10-20 23:55:19.334080181 +0300
@@ -3,5 +3,5 @@
 ok K2: registration of new tests
 FAIL K3: runs/T-0000/test-runs.md:3: selector TestOtherPkg matches no test in the named packages
 ok K4: restore evidence of mutants
-ok K5: projection: figure 2/1, maximum 100/10, tolerance 60/4, limit 160/14, margin 158/13 (the figure at landing is still measured at landing)
+ok K5: projection: figure 72/3, maximum 100/10, margin 28/7 (the figure at landing is still measured at landing)
 ok K6: conditions table
```

### bad-K3-other-tag

```text
run: bad-K3-other-tag: exit 1 (wanted 1); diff (-EXPECT +got):
--- /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.aQ88YL6Gyu/expect	2026-10-20 23:55:19.441080645 +0300
+++ /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.aQ88YL6Gyu/got	2026-10-20 23:55:19.440080640 +0300
@@ -3,5 +3,5 @@
 ok K2: registration of new tests
 FAIL K3: runs/T-0000/test-runs.md:3: selector TestIntegOnly matches no test in the named packages
 ok K4: restore evidence of mutants
-ok K5: projection: figure 2/1, maximum 100/10, tolerance 60/4, limit 160/14, margin 158/13 (the figure at landing is still measured at landing)
+ok K5: projection: figure 72/3, maximum 100/10, margin 28/7 (the figure at landing is still measured at landing)
 ok K6: conditions table
```

### bad-K3-span-across-lines

```text
run: bad-K3-span-across-lines: exit 1 (wanted 1); diff (-EXPECT +got):
--- /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.VnXipQ4LtY/expect	2026-10-20 23:55:19.547241562 +0300
+++ /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.VnXipQ4LtY/got	2026-10-20 23:55:19.546081100 +0300
@@ -3,5 +3,5 @@
 ok K2: registration of new tests
 FAIL K3: runs/T-0000/test-runs.md:3: selector TestNoSuch matches no test in the named packages
 ok K4: restore evidence of mutants
-ok K5: projection: figure 3/1, maximum 100/10, tolerance 60/4, limit 160/14, margin 157/13 (the figure at landing is still measured at landing)
+ok K5: projection: figure 73/3, maximum 100/10, margin 27/7 (the figure at landing is still measured at landing)
 ok K6: conditions table
```

### bad-K3-subtest

```text
run: bad-K3-subtest: exit 1 (wanted 1); diff (-EXPECT +got):
--- /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.octPLLZAWM/expect	2026-10-20 23:55:19.655060136 +0300
+++ /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.octPLLZAWM/got	2026-10-20 23:55:19.654286365 +0300
@@ -3,5 +3,5 @@
 ok K2: registration of new tests
 FAIL K3: runs/T-0000/test-runs.md:3: unsupported selector form: 'TestAlpha/sub'
 ok K4: restore evidence of mutants
-ok K5: projection: figure 2/1, maximum 100/10, tolerance 60/4, limit 160/14, margin 158/13 (the figure at landing is still measured at landing)
+ok K5: projection: figure 72/3, maximum 100/10, margin 28/7 (the figure at landing is still measured at landing)
 ok K6: conditions table
```

### bad-K3-unclosed-span-eof

```text
run: bad-K3-unclosed-span-eof: exit 1 (wanted 1); diff (-EXPECT +got):
--- /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.4hb66ySeb2/expect	2026-10-20 23:55:19.763837372 +0300
+++ /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.4hb66ySeb2/got	2026-10-20 23:55:19.763065247 +0300
@@ -3,5 +3,5 @@
 ok K2: registration of new tests
 FAIL K3: runs/T-0000/test-runs.md:3: selector TestMissing matches no test in the named packages
 ok K4: restore evidence of mutants
-ok K5: projection: figure 3/1, maximum 100/10, tolerance 60/4, limit 160/14, margin 157/13 (the figure at landing is still measured at landing)
+ok K5: projection: figure 73/3, maximum 100/10, margin 27/7 (the figure at landing is still measured at landing)
 ok K6: conditions table
```

### bad-K3-unclosed-span-paragraph

```text
run: bad-K3-unclosed-span-paragraph: exit 1 (wanted 1); diff (-EXPECT +got):
--- /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.XNpOLQNtce/expect	2026-10-20 23:55:19.871928005 +0300
+++ /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.XNpOLQNtce/got	2026-10-20 23:55:19.871152234 +0300
@@ -3,5 +3,5 @@
 ok K2: registration of new tests
 FAIL K3: runs/T-0000/test-runs.md:3: selector TestMissing matches no test in the named packages
 ok K4: restore evidence of mutants
-ok K5: projection: figure 5/1, maximum 100/10, tolerance 60/4, limit 160/14, margin 155/13 (the figure at landing is still measured at landing)
+ok K5: projection: figure 75/3, maximum 100/10, margin 25/7 (the figure at landing is still measured at landing)
 ok K6: conditions table
```

### bad-K3-unquoted

```text
run: bad-K3-unquoted: exit 1 (wanted 1); diff (-EXPECT +got):
--- /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.e6DlKjQ20S/expect	2026-10-20 23:55:19.979744455 +0300
+++ /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.e6DlKjQ20S/got	2026-10-20 23:55:19.979048297 +0300
@@ -3,5 +3,5 @@
 ok K2: registration of new tests
 FAIL K3: runs/T-0000/test-runs.md:3: quote the selector with single quotes
 ok K4: restore evidence of mutants
-ok K5: projection: figure 2/1, maximum 100/10, tolerance 60/4, limit 160/14, margin 158/13 (the figure at landing is still measured at landing)
+ok K5: projection: figure 72/3, maximum 100/10, margin 28/7 (the figure at landing is still measured at landing)
 ok K6: conditions table
```

### bad-K3-unsupported-form

```text
run: bad-K3-unsupported-form: exit 1 (wanted 1); diff (-EXPECT +got):
--- /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.q8hiuP8wK8/expect	2026-10-20 23:55:20.085827418 +0300
+++ /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.q8hiuP8wK8/got	2026-10-20 23:55:20.085049177 +0300
@@ -3,5 +3,5 @@
 ok K2: registration of new tests
 FAIL K3: runs/T-0000/test-runs.md:3: unsupported selector form: 'Test.*'
 ok K4: restore evidence of mutants
-ok K5: projection: figure 2/1, maximum 100/10, tolerance 60/4, limit 160/14, margin 158/13 (the figure at landing is still measured at landing)
+ok K5: projection: figure 72/3, maximum 100/10, margin 28/7 (the figure at landing is still measured at landing)
 ok K6: conditions table
```

### bad-K4-dropped-in-description

```text
run: bad-K4-dropped-in-description: exit 1 (wanted 1); diff (-EXPECT +got):
--- /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.frvFtox911/expect	2026-10-20 23:55:20.192083897 +0300
+++ /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.frvFtox911/got	2026-10-20 23:55:20.192070784 +0300
@@ -7,5 +7,5 @@
 FAIL K4: runs/T-0000/mutants.md:5: M2 has no restore evidence
 FAIL K4: runs/T-0000/mutants.md:8: M5 has no restore evidence
 FAIL K4: runs/T-0000/mutants.md:9: M6 has no restore evidence
-ok K5: projection: figure 10/1, maximum 100/10, tolerance 60/4, limit 160/14, margin 150/13 (the figure at landing is still measured at landing)
+ok K5: projection: figure 80/3, maximum 100/10, margin 20/7 (the figure at landing is still measured at landing)
 ok K6: conditions table
```

### bad-K4-form-a-different

```text
run: bad-K4-form-a-different: exit 1 (wanted 1); diff (-EXPECT +got):
--- /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.iIwsGZ6xAO/expect	2026-10-20 23:55:20.301880601 +0300
+++ /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.iIwsGZ6xAO/got	2026-10-20 23:55:20.301066477 +0300
@@ -3,5 +3,5 @@
 ok K2: registration of new tests
 ok K3: test selectors
 FAIL K4: runs/T-0000/mutants.md:5: M1 has no restore evidence
-ok K5: projection: figure 2/1, maximum 100/10, tolerance 60/4, limit 160/14, margin 158/13 (the figure at landing is still measured at landing)
+ok K5: projection: figure 72/3, maximum 100/10, margin 28/7 (the figure at landing is still measured at landing)
 ok K6: conditions table
```

### bad-K4-form-b-equal-hashes

```text
run: bad-K4-form-b-equal-hashes: exit 1 (wanted 1); diff (-EXPECT +got):
--- /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.BgeDCQeeBQ/expect	2026-10-20 23:55:20.408169680 +0300
+++ /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.BgeDCQeeBQ/got	2026-10-20 23:55:20.407084828 +0300
@@ -3,5 +3,5 @@
 ok K2: registration of new tests
 ok K3: test selectors
 FAIL K4: runs/T-0000/mutants.md:5: M1 has no restore evidence
-ok K5: projection: figure 6/1, maximum 100/10, tolerance 60/4, limit 160/14, margin 154/13 (the figure at landing is still measured at landing)
+ok K5: projection: figure 76/3, maximum 100/10, margin 24/7 (the figure at landing is still measured at landing)
 ok K6: conditions table
```

### bad-K4-result-only

```text
run: bad-K4-result-only: exit 1 (wanted 1); diff (-EXPECT +got):
--- /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.kcQ3T0WocL/expect	2026-10-20 23:55:20.516085300 +0300
+++ /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.kcQ3T0WocL/got	2026-10-20 23:55:20.515866354 +0300
@@ -6,5 +6,5 @@
 skip K4: runs/T-0000/test-runs.md:9: SM29 (dropped)
 skip K4: runs/T-0000/test-runs.md:10: SM30 (superseded)
 FAIL K4: runs/T-0000/test-runs.md:14: SM31 has no restore evidence
-ok K5: projection: figure 11/1, maximum 100/10, tolerance 60/4, limit 160/14, margin 149/13 (the figure at landing is still measured at landing)
+ok K5: projection: figure 81/3, maximum 100/10, margin 19/7 (the figure at landing is still measured at landing)
 ok K6: conditions table
```

### bad-K4-two-hashes-no-restore

```text
run: bad-K4-two-hashes-no-restore: exit 1 (wanted 1); diff (-EXPECT +got):
--- /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.RCZO8r63NA/expect	2026-10-20 23:55:20.624136043 +0300
+++ /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.RCZO8r63NA/got	2026-10-20 23:55:20.623085763 +0300
@@ -3,5 +3,5 @@
 ok K2: registration of new tests
 ok K3: test selectors
 FAIL K4: runs/T-0000/mutants.md:5: M1 has no restore evidence
-ok K5: projection: figure 6/1, maximum 100/10, tolerance 60/4, limit 160/14, margin 154/13 (the figure at landing is still measured at landing)
+ok K5: projection: figure 76/3, maximum 100/10, margin 24/7 (the figure at landing is still measured at landing)
 ok K6: conditions table
```

### bad-K5-over-files

```text
run: bad-K5-over-files: exit 1 (wanted 1); diff (-EXPECT +got):
--- /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.AxKwY6qw5p/expect	2026-10-20 23:55:20.745921113 +0300
+++ /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.AxKwY6qw5p/got	2026-10-20 23:55:20.745151707 +0300
@@ -3,5 +3,5 @@
 ok K2: registration of new tests
 ok K3: test selectors
 ok K4: restore evidence of mutants
-FAIL K5: -:0: projection: figure 6/6, maximum 100/1, tolerance 60/4, limit 160/5, margin 154/-1 (the figure at landing is still measured at landing)
+FAIL K5: -:0: projection: figure 76/8, maximum 100/1, margin 24/-7 (the figure at landing is still measured at landing)
 ok K6: conditions table
```

### bad-K5-over-lines

```text
run: bad-K5-over-lines: exit 1 (wanted 1); diff (-EXPECT +got):
--- /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.3xtYrtpd5Z/expect	2026-10-20 23:55:20.853171955 +0300
+++ /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.3xtYrtpd5Z/got	2026-10-20 23:55:20.852086755 +0300
@@ -3,5 +3,5 @@
 ok K2: registration of new tests
 ok K3: test selectors
 ok K4: restore evidence of mutants
-FAIL K5: -:0: projection: figure 161/1, maximum 100/10, tolerance 60/4, limit 160/14, margin -1/13 (the figure at landing is still measured at landing)
+FAIL K5: -:0: projection: figure 231/3, maximum 100/10, margin -131/7 (the figure at landing is still measured at landing)
 ok K6: conditions table
```

### bad-K6-empty-table

```text
run: bad-K6-empty-table: exit 1 (wanted 1); diff (-EXPECT +got):
--- /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.K73xuXeKSU/expect	2026-10-20 23:55:20.948187195 +0300
+++ /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.K73xuXeKSU/got	2026-10-20 23:55:20.947297458 +0300
@@ -3,5 +3,5 @@
 ok K2: registration of new tests
 ok K3: test selectors
 ok K4: restore evidence of mutants
-ok K5: projection: figure 1/1, maximum 100/10, tolerance 60/4, limit 160/14, margin 159/13 (the figure at landing is still measured at landing)
+ok K5: projection: figure 71/3, maximum 100/10, margin 29/7 (the figure at landing is still measured at landing)
 FAIL K6: runs/T-0000/conditions.md:1: expected one | Condition | Evidence | table with at least one row
```

### bad-K6-line-past-end

```text
run: bad-K6-line-past-end: exit 1 (wanted 1); diff (-EXPECT +got):
--- /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.tn7tdeWMWz/expect	2026-10-20 23:55:21.053991423 +0300
+++ /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.tn7tdeWMWz/got	2026-10-20 23:55:21.053216400 +0300
@@ -3,5 +3,5 @@
 ok K2: registration of new tests
 ok K3: test selectors
 ok K4: restore evidence of mutants
-ok K5: projection: figure 2/1, maximum 100/10, tolerance 60/4, limit 160/14, margin 158/13 (the figure at landing is still measured at landing)
+ok K5: projection: figure 72/3, maximum 100/10, margin 28/7 (the figure at landing is still measured at landing)
 FAIL K6: runs/T-0000/conditions.md:5: line range in internal/x/a_test.go:99 is past the end of internal/x/a_test.go
```

### bad-K6-line-zero

```text
run: bad-K6-line-zero: exit 1 (wanted 1); diff (-EXPECT +got):
--- /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.VS1LLNBqYp/expect	2026-10-20 23:55:21.155088067 +0300
+++ /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.VS1LLNBqYp/got	2026-10-20 23:55:21.154088062 +0300
@@ -3,5 +3,5 @@
 ok K2: registration of new tests
 ok K3: test selectors
 ok K4: restore evidence of mutants
-ok K5: projection: figure 2/1, maximum 100/10, tolerance 60/4, limit 160/14, margin 158/13 (the figure at landing is still measured at landing)
+ok K5: projection: figure 72/3, maximum 100/10, margin 28/7 (the figure at landing is still measured at landing)
 FAIL K6: runs/T-0000/conditions.md:5: line range in internal/x/a_test.go:0 starts before line 1
```

### bad-K6-missing-commit

```text
run: bad-K6-missing-commit: exit 1 (wanted 1); diff (-EXPECT +got):
--- /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.AWCRBgMKrA/expect	2026-10-20 23:55:21.257088508 +0300
+++ /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.AWCRBgMKrA/got	2026-10-20 23:55:21.256785230 +0300
@@ -3,5 +3,5 @@
 ok K2: registration of new tests
 ok K3: test selectors
 ok K4: restore evidence of mutants
-ok K5: projection: figure 2/1, maximum 100/10, tolerance 60/4, limit 160/14, margin 158/13 (the figure at landing is still measured at landing)
+ok K5: projection: figure 72/3, maximum 100/10, margin 28/7 (the figure at landing is still measured at landing)
 FAIL K6: runs/T-0000/conditions.md:5: commit deadbee does not exist
```

### bad-K6-missing-file

```text
run: bad-K6-missing-file: exit 1 (wanted 1); diff (-EXPECT +got):
--- /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.bPhT9e5vVu/expect	2026-10-20 23:55:21.344088885 +0300
+++ /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.bPhT9e5vVu/got	2026-10-20 23:55:21.343088881 +0300
@@ -3,5 +3,5 @@
 ok K2: registration of new tests
 ok K3: test selectors
 ok K4: restore evidence of mutants
-ok K5: projection: figure 5/1, maximum 100/10, tolerance 60/4, limit 160/14, margin 155/13 (the figure at landing is still measured at landing)
+ok K5: projection: figure 75/3, maximum 100/10, margin 25/7 (the figure at landing is still measured at landing)
 FAIL K6: runs/T-0000/conditions.md:0: missing file
```

### bad-K6-missing-path

```text
run: bad-K6-missing-path: exit 1 (wanted 1); diff (-EXPECT +got):
--- /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.9RSW4bbbgJ/expect	2026-10-20 23:55:21.458205023 +0300
+++ /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.9RSW4bbbgJ/got	2026-10-20 23:55:21.457294078 +0300
@@ -3,5 +3,5 @@
 ok K2: registration of new tests
 ok K3: test selectors
 ok K4: restore evidence of mutants
-ok K5: projection: figure 2/1, maximum 100/10, tolerance 60/4, limit 160/14, margin 158/13 (the figure at landing is still measured at landing)
+ok K5: projection: figure 72/3, maximum 100/10, margin 28/7 (the figure at landing is still measured at landing)
 FAIL K6: runs/T-0000/conditions.md:5: path internal/x/no_such.go does not exist
```

### bad-K6-missing-test

```text
run: bad-K6-missing-test: exit 1 (wanted 1); diff (-EXPECT +got):
--- /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.lm2RRjqPe4/expect	2026-10-20 23:55:21.564198968 +0300
+++ /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.lm2RRjqPe4/got	2026-10-20 23:55:21.563089833 +0300
@@ -3,5 +3,5 @@
 ok K2: registration of new tests
 ok K3: test selectors
 ok K4: restore evidence of mutants
-ok K5: projection: figure 2/1, maximum 100/10, tolerance 60/4, limit 160/14, margin 158/13 (the figure at landing is still measured at landing)
+ok K5: projection: figure 72/3, maximum 100/10, margin 28/7 (the figure at landing is still measured at landing)
 FAIL K6: runs/T-0000/conditions.md:5: test TestDoesNotExist does not exist
```

### bad-K6-mixed-valid-invalid

```text
run: bad-K6-mixed-valid-invalid: exit 1 (wanted 1); diff (-EXPECT +got):
--- /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.cAKFLTZb6v/expect	2026-10-20 23:55:21.680948612 +0300
+++ /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.cAKFLTZb6v/got	2026-10-20 23:55:21.680016385 +0300
@@ -3,5 +3,5 @@
 ok K2: registration of new tests
 ok K3: test selectors
 ok K4: restore evidence of mutants
-ok K5: projection: figure 2/1, maximum 100/10, tolerance 60/4, limit 160/14, margin 158/13 (the figure at landing is still measured at landing)
+ok K5: projection: figure 72/3, maximum 100/10, margin 28/7 (the figure at landing is still measured at landing)
 FAIL K6: runs/T-0000/conditions.md:5: test TestDoesNotExist does not exist
```

### bad-K6-no-span

```text
run: bad-K6-no-span: exit 1 (wanted 1); diff (-EXPECT +got):
--- /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.rcfC6cVDr0/expect	2026-10-20 23:55:21.786090799 +0300
+++ /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.rcfC6cVDr0/got	2026-10-20 23:55:21.785884494 +0300
@@ -3,5 +3,5 @@
 ok K2: registration of new tests
 ok K3: test selectors
 ok K4: restore evidence of mutants
-ok K5: projection: figure 2/1, maximum 100/10, tolerance 60/4, limit 160/14, margin 158/13 (the figure at landing is still measured at landing)
+ok K5: projection: figure 72/3, maximum 100/10, margin 28/7 (the figure at landing is still measured at landing)
 FAIL K6: runs/T-0000/conditions.md:5: Evidence cell has no code span
```

### bad-K6-range-zero

```text
run: bad-K6-range-zero: exit 1 (wanted 1); diff (-EXPECT +got):
--- /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.PSCzPifdNo/expect	2026-10-20 23:55:21.901211952 +0300
+++ /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.PSCzPifdNo/got	2026-10-20 23:55:21.900091292 +0300
@@ -3,5 +3,5 @@
 ok K2: registration of new tests
 ok K3: test selectors
 ok K4: restore evidence of mutants
-ok K5: projection: figure 2/1, maximum 100/10, tolerance 60/4, limit 160/14, margin 158/13 (the figure at landing is still measured at landing)
+ok K5: projection: figure 72/3, maximum 100/10, margin 28/7 (the figure at landing is still measured at landing)
 FAIL K6: runs/T-0000/conditions.md:5: line range in internal/x/a_test.go:0-3 starts before line 1
```

### bad-crlf-conditions

```text
run: bad-crlf-conditions: exit 1 (wanted 1); diff (-EXPECT +got):
--- /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.KzarIWeNG4/expect	2026-10-20 23:55:22.077092059 +0300
+++ /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.KzarIWeNG4/got	2026-10-20 23:55:22.076092054 +0300
@@ -3,5 +3,5 @@
 ok K2: registration of new tests
 ok K3: test selectors
 ok K4: restore evidence of mutants
-ok K5: projection: figure 10/1, maximum 100/10, tolerance 60/4, limit 160/14, margin 150/13 (the figure at landing is still measured at landing)
+ok K5: projection: figure 80/3, maximum 100/10, margin 20/7 (the figure at landing is still measured at landing)
 FAIL K6: runs/T-0000/conditions.md:5: path internal/x/no_such.go does not exist
```

### bad-usage-bad-budget

```text
run: bad-usage-bad-budget: exit 2 (wanted 2); diff (-EXPECT +got):
--- /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.FdKfz7st0u/expect	2026-10-20 23:55:22.568897472 +0300
+++ /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.FdKfz7st0u/got	2026-10-20 23:55:22.568108080 +0300
@@ -1,2 +1,2 @@
 2
-freeze-check: usage: freeze-check.sh --task T-xxxx --base REV --budget LINES/FILES [--evidence PATH] [--no-conditions]
+freeze-check: usage: freeze-check.sh --task T-xxxx --base REV --budget LINES/FILES [--closeout LINES/FILES] [--no-conditions]
```

### bad-usage-closeout

```text
run: bad-usage-closeout: exit 0 (wanted 2); diff (-EXPECT +got):
--- /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.1K8dLjvrMn/expect	2026-10-20 23:55:22.683083956 +0300
+++ /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.1K8dLjvrMn/got	2026-10-20 23:55:22.682135371 +0300
@@ -1,2 +1,7 @@
-2
-freeze-check: usage: freeze-check.sh --task T-xxxx --base REV --budget LINES/FILES [--evidence PATH] [--no-conditions]
+0
+ok K1: white space
+ok K2: registration of new tests
+ok K3: test selectors
+ok K4: restore evidence of mutants
+ok K5: projection: figure 70/2, maximum 100/10, margin 30/8 (the figure at landing is still measured at landing)
+ok K6: conditions table
```

### bad-usage-evidence-missing

```text
run: bad-usage-evidence-missing: exit 2 (wanted 2); diff (-EXPECT +got):
--- /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.ogzYOHdQUJ/expect	2026-10-20 23:55:22.741060499 +0300
+++ /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.ogzYOHdQUJ/got	2026-10-20 23:55:22.740219539 +0300
@@ -1,2 +1,2 @@
 2
-freeze-check: --evidence runs/T-0000/missing.md: is absent from both base and HEAD
+freeze-check: usage: freeze-check.sh --task T-xxxx --base REV --budget LINES/FILES [--closeout LINES/FILES] [--no-conditions]
```

### bad-usage-evidence-no-value

```text
run: bad-usage-evidence-no-value: exit 2 (wanted 2); diff (-EXPECT +got):
--- /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.09QEFePDp4/expect	2026-10-20 23:55:22.799095184 +0300
+++ /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.09QEFePDp4/got	2026-10-20 23:55:22.798095180 +0300
@@ -1,2 +1,2 @@
 2
-freeze-check: usage: freeze-check.sh --task T-xxxx --base REV --budget LINES/FILES [--evidence PATH] [--no-conditions]
+freeze-check: usage: freeze-check.sh --task T-xxxx --base REV --budget LINES/FILES [--closeout LINES/FILES] [--no-conditions]
```

### bad-usage-evidence-other-task

```text
run: bad-usage-evidence-other-task: exit 2 (wanted 2); diff (-EXPECT +got):
--- /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.wnboF67rfB/expect	2026-10-20 23:55:22.861095453 +0300
+++ /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.wnboF67rfB/got	2026-10-20 23:55:22.860095449 +0300
@@ -1,2 +1,2 @@
 2
-freeze-check: --evidence runs/T-0001/report.md: must be under runs/T-0000/
+freeze-check: usage: freeze-check.sh --task T-xxxx --base REV --budget LINES/FILES [--closeout LINES/FILES] [--no-conditions]
```

### bad-usage-evidence-suffix

```text
run: bad-usage-evidence-suffix: exit 2 (wanted 2); diff (-EXPECT +got):
--- /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.wo2ShiqUvC/expect	2026-10-20 23:55:22.925147261 +0300
+++ /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.wo2ShiqUvC/got	2026-10-20 23:55:22.924198782 +0300
@@ -1,2 +1,2 @@
 2
-freeze-check: --evidence runs/T-0000/checks.sh: must end in .md, .tsv or .txt
+freeze-check: usage: freeze-check.sh --task T-xxxx --base REV --budget LINES/FILES [--closeout LINES/FILES] [--no-conditions]
```

### bad-usage-evidence-task-prefix

```text
run: bad-usage-evidence-task-prefix: exit 2 (wanted 2); diff (-EXPECT +got):
--- /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.wF7t6M1U9s/expect	2026-10-20 23:55:22.987338412 +0300
+++ /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.wF7t6M1U9s/got	2026-10-20 23:55:22.986095994 +0300
@@ -1,2 +1,2 @@
 2
-freeze-check: --evidence runs/T-0000-extra/report.md: must be under runs/T-0000/
+freeze-check: usage: freeze-check.sh --task T-xxxx --base REV --budget LINES/FILES [--closeout LINES/FILES] [--no-conditions]
```

### bad-usage-evidence-traversal

```text
run: bad-usage-evidence-traversal: exit 2 (wanted 2); diff (-EXPECT +got):
--- /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.85BoyvT9Rs/expect	2026-10-20 23:55:23.044826304 +0300
+++ /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.85BoyvT9Rs/got	2026-10-20 23:55:23.044075099 +0300
@@ -1,2 +1,2 @@
 2
-freeze-check: --evidence runs/T-0000/../T-0000/conditions.md: path traversal is not allowed
+freeze-check: usage: freeze-check.sh --task T-xxxx --base REV --budget LINES/FILES [--closeout LINES/FILES] [--no-conditions]
```

### bad-usage-missing-budget

```text
run: bad-usage-missing-budget: exit 2 (wanted 2); diff (-EXPECT +got):
--- /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.ojhbHyRJwP/expect	2026-10-20 23:55:23.101923116 +0300
+++ /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.ojhbHyRJwP/got	2026-10-20 23:55:23.101058492 +0300
@@ -1,2 +1,2 @@
 2
-freeze-check: usage: freeze-check.sh --task T-xxxx --base REV --budget LINES/FILES [--evidence PATH] [--no-conditions]
+freeze-check: usage: freeze-check.sh --task T-xxxx --base REV --budget LINES/FILES [--closeout LINES/FILES] [--no-conditions]
```

### bad-usage-no-args

```text
run: bad-usage-no-args: exit 2 (wanted 2); diff (-EXPECT +got):
--- /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.FRfw2kJ4l1/expect	2026-10-20 23:55:23.159906886 +0300
+++ /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.FRfw2kJ4l1/got	2026-10-20 23:55:23.159022252 +0300
@@ -1,2 +1,2 @@
 2
-freeze-check: usage: freeze-check.sh --task T-xxxx --base REV --budget LINES/FILES [--evidence PATH] [--no-conditions]
+freeze-check: usage: freeze-check.sh --task T-xxxx --base REV --budget LINES/FILES [--closeout LINES/FILES] [--no-conditions]
```

### bad-usage-unknown-flag

```text
run: bad-usage-unknown-flag: exit 2 (wanted 2); diff (-EXPECT +got):
--- /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.43HPWIiIRY/expect	2026-10-20 23:55:23.217852445 +0300
+++ /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.43HPWIiIRY/got	2026-10-20 23:55:23.216096990 +0300
@@ -1,2 +1,2 @@
 2
-freeze-check: usage: freeze-check.sh --task T-xxxx --base REV --budget LINES/FILES [--evidence PATH] [--no-conditions]
+freeze-check: usage: freeze-check.sh --task T-xxxx --base REV --budget LINES/FILES [--closeout LINES/FILES] [--no-conditions]
```

### good-K1-clean

```text
run: good-K1-clean: exit 0 (wanted 0); diff (-EXPECT +got):
--- /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.6RKBZpNR8H/expect	2026-10-20 23:55:23.338097518 +0300
+++ /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.6RKBZpNR8H/got	2026-10-20 23:55:23.337925804 +0300
@@ -3,5 +3,5 @@
 ok K2: registration of new tests
 ok K3: test selectors
 ok K4: restore evidence of mutants
-ok K5: projection: figure 1/1, maximum 100/10, tolerance 60/4, limit 160/14, margin 159/13 (the figure at landing is still measured at landing)
+ok K5: projection: figure 71/3, maximum 100/10, margin 29/7 (the figure at landing is still measured at landing)
 ok K6: conditions table
```

### good-K2-prd-by-file

```text
run: good-K2-prd-by-file: exit 0 (wanted 0); diff (-EXPECT +got):
--- /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.8BdjpK2bq3/expect	2026-10-20 23:55:23.470255215 +0300
+++ /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.8BdjpK2bq3/got	2026-10-20 23:55:23.469307383 +0300
@@ -3,5 +3,5 @@
 ok K2: registration of new tests
 ok K3: test selectors
 ok K4: restore evidence of mutants
-ok K5: projection: figure 4/3, maximum 100/10, tolerance 60/4, limit 160/14, margin 156/11 (the figure at landing is still measured at landing)
+ok K5: projection: figure 74/5, maximum 100/10, margin 26/5 (the figure at landing is still measured at landing)
 ok K6: conditions table
```

### good-K2-registered

```text
run: good-K2-registered: exit 0 (wanted 0); diff (-EXPECT +got):
--- /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.fmqulN9cHC/expect	2026-10-20 23:55:23.601098656 +0300
+++ /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.fmqulN9cHC/got	2026-10-20 23:55:23.600098652 +0300
@@ -3,5 +3,5 @@
 ok K2: registration of new tests
 ok K3: test selectors
 ok K4: restore evidence of mutants
-ok K5: projection: figure 4/3, maximum 100/10, tolerance 60/4, limit 160/14, margin 156/11 (the figure at landing is still measured at landing)
+ok K5: projection: figure 74/5, maximum 100/10, margin 26/5 (the figure at landing is still measured at landing)
 ok K6: conditions table
```

### good-K3-alternatives

```text
run: good-K3-alternatives: exit 0 (wanted 0); diff (-EXPECT +got):
--- /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.mWgyeQV5eq/expect	2026-10-20 23:55:23.718099163 +0300
+++ /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.mWgyeQV5eq/got	2026-10-20 23:55:23.717982569 +0300
@@ -3,5 +3,5 @@
 ok K2: registration of new tests
 ok K3: test selectors
 ok K4: restore evidence of mutants
-ok K5: projection: figure 4/1, maximum 100/10, tolerance 60/4, limit 160/14, margin 156/13 (the figure at landing is still measured at landing)
+ok K5: projection: figure 74/3, maximum 100/10, margin 26/7 (the figure at landing is still measured at landing)
 ok K6: conditions table
```

### good-K3-span-across-lines

```text
run: good-K3-span-across-lines: exit 0 (wanted 0); diff (-EXPECT +got):
--- /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.rqwbKZvgAO/expect	2026-10-20 23:55:23.837250470 +0300
+++ /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.rqwbKZvgAO/got	2026-10-20 23:55:23.836217352 +0300
@@ -3,5 +3,5 @@
 ok K2: registration of new tests
 ok K3: test selectors
 ok K4: restore evidence of mutants
-ok K5: projection: figure 7/1, maximum 100/10, tolerance 60/4, limit 160/14, margin 153/13 (the figure at landing is still measured at landing)
+ok K5: projection: figure 77/3, maximum 100/10, margin 23/7 (the figure at landing is still measured at landing)
 ok K6: conditions table
```

### good-K4-forms

```text
run: good-K4-forms: exit 0 (wanted 0); diff (-EXPECT +got):
--- /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.eMplWnfR8F/expect	2026-10-20 23:55:23.961100215 +0300
+++ /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.eMplWnfR8F/got	2026-10-20 23:55:23.960100210 +0300
@@ -3,5 +3,5 @@
 ok K2: registration of new tests
 ok K3: test selectors
 ok K4: restore evidence of mutants
-ok K5: projection: figure 22/2, maximum 100/10, tolerance 60/4, limit 160/14, margin 138/12 (the figure at landing is still measured at landing)
+ok K5: projection: figure 92/4, maximum 100/10, margin 8/6 (the figure at landing is still measured at landing)
 ok K6: conditions table
```

### good-K5-at-max

```text
run: good-K5-at-max: exit 1 (wanted 0); diff (-EXPECT +got):
--- /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.8xMhMX55Up/expect	2026-10-20 23:55:24.084100747 +0300
+++ /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.8xMhMX55Up/got	2026-10-20 23:55:24.083100743 +0300
@@ -1,7 +1,7 @@
-0
+1
 ok K1: white space
 ok K2: registration of new tests
 ok K3: test selectors
 ok K4: restore evidence of mutants
-ok K5: projection: figure 5/1, maximum 5/1, tolerance 60/4, limit 65/5, margin 60/4 (the figure at landing is still measured at landing)
+FAIL K5: -:0: projection: figure 75/3, maximum 5/1, margin -70/-2 (the figure at landing is still measured at landing)
 ok K6: conditions table
```

### good-K5-ceiling

```text
run: good-K5-ceiling: exit 0 (wanted 0); diff (-EXPECT +got):
--- /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.x9J3fCI0c2/expect	2026-10-20 23:55:24.202271546 +0300
+++ /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.x9J3fCI0c2/got	2026-10-20 23:55:24.201403959 +0300
@@ -3,5 +3,5 @@
 ok K2: registration of new tests
 ok K3: test selectors
 ok K4: restore evidence of mutants
-ok K5: projection: figure 1/1, maximum 3000/10, tolerance 500/4, limit 3500/14, margin 3499/13 (the figure at landing is still measured at landing)
+ok K5: projection: figure 71/3, maximum 3000/10, margin 2929/7 (the figure at landing is still measured at landing)
 ok K6: conditions table
```

### good-K5-evidence-added

```text
run: good-K5-evidence-added: exit 2 (wanted 0); diff (-EXPECT +got):
--- /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.Ri8NhQJEGM/expect	2026-10-20 23:55:24.267101539 +0300
+++ /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.Ri8NhQJEGM/got	2026-10-20 23:55:24.266101535 +0300
@@ -1,7 +1,2 @@
-0
-ok K1: white space
-ok K2: registration of new tests
-ok K3: test selectors
-ok K4: restore evidence of mutants
-ok K5: projection: figure 1/1, maximum 100/10, tolerance 60/4, limit 160/14, margin 159/13 (the figure at landing is still measured at landing)
-ok K6: conditions table
+2
+freeze-check: usage: freeze-check.sh --task T-xxxx --base REV --budget LINES/FILES [--closeout LINES/FILES] [--no-conditions]
```

### good-K5-evidence-deleted

```text
run: good-K5-evidence-deleted: exit 2 (wanted 0); diff (-EXPECT +got):
--- /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.nFP3yusrhI/expect	2026-10-20 23:55:24.332103250 +0300
+++ /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.nFP3yusrhI/got	2026-10-20 23:55:24.331101882 +0300
@@ -1,7 +1,2 @@
-0
-ok K1: white space
-ok K2: registration of new tests
-ok K3: test selectors
-ok K4: restore evidence of mutants
-ok K5: projection: figure 1/1, maximum 100/10, tolerance 60/4, limit 160/14, margin 159/13 (the figure at landing is still measured at landing)
-ok K6: conditions table
+2
+freeze-check: usage: freeze-check.sh --task T-xxxx --base REV --budget LINES/FILES [--closeout LINES/FILES] [--no-conditions]
```

### good-K5-evidence-duplicate

```text
run: good-K5-evidence-duplicate: exit 2 (wanted 0); diff (-EXPECT +got):
--- /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.kDddAbX6hm/expect	2026-10-20 23:55:24.396184753 +0300
+++ /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.kDddAbX6hm/got	2026-10-20 23:55:24.395244580 +0300
@@ -1,7 +1,2 @@
-0
-ok K1: white space
-ok K2: registration of new tests
-ok K3: test selectors
-ok K4: restore evidence of mutants
-ok K5: projection: figure 1/1, maximum 100/10, tolerance 60/4, limit 160/14, margin 159/13 (the figure at landing is still measured at landing)
-ok K6: conditions table
+2
+freeze-check: usage: freeze-check.sh --task T-xxxx --base REV --budget LINES/FILES [--closeout LINES/FILES] [--no-conditions]
```

### good-K5-evidence-inputs

```text
run: good-K5-evidence-inputs: exit 2 (wanted 0); diff (-EXPECT +got):
--- /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.7ssoWW3cya/expect	2026-10-20 23:55:24.466202232 +0300
+++ /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.7ssoWW3cya/got	2026-10-20 23:55:24.465102396 +0300
@@ -1,7 +1,2 @@
-0
-ok K1: white space
-ok K2: registration of new tests
-ok K3: test selectors
-ok K4: restore evidence of mutants
-ok K5: projection: figure 6/3, maximum 0/0, tolerance 60/4, limit 60/4, margin 54/1 (the figure at landing is still measured at landing)
-ok K6: conditions table
+2
+freeze-check: usage: freeze-check.sh --task T-xxxx --base REV --budget LINES/FILES [--closeout LINES/FILES] [--no-conditions]
```

### good-K5-evidence-kinds

```text
run: good-K5-evidence-kinds: exit 2 (wanted 0); diff (-EXPECT +got):
--- /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.1urBStNvgl/expect	2026-10-20 23:55:24.540102721 +0300
+++ /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.1urBStNvgl/got	2026-10-20 23:55:24.539904098 +0300
@@ -1,7 +1,2 @@
-0
-ok K1: white space
-ok K2: registration of new tests
-ok K3: test selectors
-ok K4: restore evidence of mutants
-ok K5: projection: figure 2/2, maximum 0/0, tolerance 60/4, limit 60/4, margin 58/2 (the figure at landing is still measured at landing)
-ok K6: conditions table
+2
+freeze-check: usage: freeze-check.sh --task T-xxxx --base REV --budget LINES/FILES [--closeout LINES/FILES] [--no-conditions]
```

### good-K5-evidence-literal

```text
run: good-K5-evidence-literal: exit 2 (wanted 0); diff (-EXPECT +got):
--- /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.0Ai1p6EQmc/expect	2026-10-20 23:55:24.608129959 +0300
+++ /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.0Ai1p6EQmc/got	2026-10-20 23:55:24.607201269 +0300
@@ -1,7 +1,2 @@
-0
-ok K1: white space
-ok K2: registration of new tests
-ok K3: test selectors
-ok K4: restore evidence of mutants
-ok K5: projection: figure 3/2, maximum 100/10, tolerance 60/4, limit 160/14, margin 157/12 (the figure at landing is still measured at landing)
-ok K6: conditions table
+2
+freeze-check: usage: freeze-check.sh --task T-xxxx --base REV --budget LINES/FILES [--closeout LINES/FILES] [--no-conditions]
```

### good-K5-evidence-modified

```text
run: good-K5-evidence-modified: exit 2 (wanted 0); diff (-EXPECT +got):
--- /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.uHNCJzj6AT/expect	2026-10-20 23:55:24.671031101 +0300
+++ /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.uHNCJzj6AT/got	2026-10-20 23:55:24.670250421 +0300
@@ -1,7 +1,2 @@
-0
-ok K1: white space
-ok K2: registration of new tests
-ok K3: test selectors
-ok K4: restore evidence of mutants
-ok K5: projection: figure 1/1, maximum 100/10, tolerance 60/4, limit 160/14, margin 159/13 (the figure at landing is still measured at landing)
-ok K6: conditions table
+2
+freeze-check: usage: freeze-check.sh --task T-xxxx --base REV --budget LINES/FILES [--closeout LINES/FILES] [--no-conditions]
```

### good-K5-file-limit

```text
run: good-K5-file-limit: exit 1 (wanted 0); diff (-EXPECT +got):
--- /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.ZqhuCrQHvj/expect	2026-10-20 23:55:24.797032397 +0300
+++ /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.ZqhuCrQHvj/got	2026-10-20 23:55:24.796069971 +0300
@@ -1,7 +1,7 @@
-0
+1
 ok K1: white space
 ok K2: registration of new tests
 ok K3: test selectors
 ok K4: restore evidence of mutants
-ok K5: projection: figure 5/5, maximum 100/1, tolerance 60/4, limit 160/5, margin 155/0 (the figure at landing is still measured at landing)
+FAIL K5: -:0: projection: figure 75/7, maximum 100/1, margin 25/-6 (the figure at landing is still measured at landing)
 ok K6: conditions table
```

### good-K5-line-limit

```text
run: good-K5-line-limit: exit 1 (wanted 0); diff (-EXPECT +got):
--- /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.iwPa8aN1o9/expect	2026-10-20 23:55:24.917287314 +0300
+++ /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.iwPa8aN1o9/got	2026-10-20 23:55:24.916414117 +0300
@@ -1,7 +1,7 @@
-0
+1
 ok K1: white space
 ok K2: registration of new tests
 ok K3: test selectors
 ok K4: restore evidence of mutants
-ok K5: projection: figure 160/1, maximum 100/10, tolerance 60/4, limit 160/14, margin 0/13 (the figure at landing is still measured at landing)
+FAIL K5: -:0: projection: figure 230/3, maximum 100/10, margin -130/7 (the figure at landing is still measured at landing)
 ok K6: conditions table
```

### good-K5-no-closeout

```text
run: good-K5-no-closeout: exit 1 (wanted 0); diff (-EXPECT +got):
--- /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.lDaoFPLFAg/expect	2026-10-20 23:55:25.038085354 +0300
+++ /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.lDaoFPLFAg/got	2026-10-20 23:55:25.037149422 +0300
@@ -1,7 +1,7 @@
-0
+1
 ok K1: white space
 ok K2: registration of new tests
 ok K3: test selectors
 ok K4: restore evidence of mutants
-ok K5: projection: figure 6/1, maximum 75/10, tolerance 60/4, limit 135/14, margin 129/13 (the figure at landing is still measured at landing)
+FAIL K5: -:0: projection: figure 76/3, maximum 75/10, margin -1/7 (the figure at landing is still measured at landing)
 ok K6: conditions table
```

### good-K5-percentage

```text
run: good-K5-percentage: exit 0 (wanted 0); diff (-EXPECT +got):
--- /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.J3sLwfPAFa/expect	2026-10-20 23:55:25.155982518 +0300
+++ /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.J3sLwfPAFa/got	2026-10-20 23:55:25.155026456 +0300
@@ -3,5 +3,5 @@
 ok K2: registration of new tests
 ok K3: test selectors
 ok K4: restore evidence of mutants
-ok K5: projection: figure 1/1, maximum 1000/10, tolerance 200/4, limit 1200/14, margin 1199/13 (the figure at landing is still measured at landing)
+ok K5: projection: figure 71/3, maximum 1000/10, margin 929/7 (the figure at landing is still measured at landing)
 ok K6: conditions table
```

### good-K5-truncation

```text
run: good-K5-truncation: exit 0 (wanted 0); diff (-EXPECT +got):
--- /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.3Se73jIUuI/expect	2026-10-20 23:55:25.276079575 +0300
+++ /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.3Se73jIUuI/got	2026-10-20 23:55:25.275020067 +0300
@@ -3,5 +3,5 @@
 ok K2: registration of new tests
 ok K3: test selectors
 ok K4: restore evidence of mutants
-ok K5: projection: figure 1/1, maximum 309/10, tolerance 61/4, limit 370/14, margin 369/13 (the figure at landing is still measured at landing)
+ok K5: projection: figure 71/3, maximum 309/10, margin 238/7 (the figure at landing is still measured at landing)
 ok K6: conditions table
```

### good-K6-no-conditions

```text
run: good-K6-no-conditions: exit 0 (wanted 0); diff (-EXPECT +got):
--- /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.kgyFHWTcyO/expect	2026-10-20 23:55:25.370085106 +0300
+++ /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.kgyFHWTcyO/got	2026-10-20 23:55:25.369254330 +0300
@@ -3,5 +3,5 @@
 ok K2: registration of new tests
 ok K3: test selectors
 ok K4: restore evidence of mutants
-ok K5: projection: figure 5/1, maximum 100/10, tolerance 60/4, limit 160/14, margin 155/13 (the figure at landing is still measured at landing)
+ok K5: projection: figure 75/3, maximum 100/10, margin 25/7 (the figure at landing is still measured at landing)
 skip K6: --no-conditions
```

### good-K6-valid

```text
run: good-K6-valid: exit 0 (wanted 0); diff (-EXPECT +got):
--- /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.9XFva6KtUF/expect	2026-10-20 23:55:25.507894281 +0300
+++ /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.9XFva6KtUF/got	2026-10-20 23:55:25.507178751 +0300
@@ -3,5 +3,5 @@
 ok K2: registration of new tests
 ok K3: test selectors
 ok K4: restore evidence of mutants
-ok K5: projection: figure 2/1, maximum 100/10, tolerance 60/4, limit 160/14, margin 158/13 (the figure at landing is still measured at landing)
+ok K5: projection: figure 72/3, maximum 100/10, margin 28/7 (the figure at landing is still measured at landing)
 ok K6: conditions table
```

### good-crlf-all

```text
run: good-crlf-all: exit 0 (wanted 0); diff (-EXPECT +got):
--- /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.YgAEkZhQZc/expect	2026-10-20 23:55:25.644107500 +0300
+++ /home/layup/.cache/layup-session/tmp-green/T-c826-stage-a/tmp.YgAEkZhQZc/got	2026-10-20 23:55:25.643107495 +0300
@@ -3,5 +3,5 @@
 ok K2: registration of new tests
 ok K3: test selectors
 ok K4: restore evidence of mutants
-ok K5: projection: figure 82/6, maximum 200/10, tolerance 60/4, limit 260/14, margin 178/8 (the figure at landing is still measured at landing)
+ok K5: projection: figure 152/8, maximum 200/10, margin 48/2 (the figure at landing is still measured at landing)
 ok K6: conditions table
```

## Final S1 record checks

The task assertion was tightened to require a numeric boundary in its PR
comment reference; the missing exception remains the observed red result. No
fixture, checker, runner or workflow changed after the per-case/suite runs.

| Command | Start UTC | End UTC | Exit | Observed result |
| --- | --- | --- | --- | --- |
| `sh runs/T-c826/checks.sh` | 2026-10-20T20:58:26Z | 2026-10-20T20:58:26Z | 1 | Same seven intended failures; preservation/source assertions pass; actual PR number pending. |
| `sh -n runs/T-c826/checks.sh` | 2026-10-20T20:58:26Z | 2026-10-20T20:58:26Z | 0 | No diagnostic. |
| `sh docs/links/link-lint.sh` | 2026-10-20T20:58:26Z | 2026-10-20T20:58:39Z | 0 | link-lint: OK  2991 links resolved |
| `sh docs/tests/run-discipline-tests.sh` | 2026-10-20T20:58:39Z | 2026-10-20T20:58:42Z | 0 | run-discipline-tests: 138 passed, 0 failed |
| `sh docs/setup/setup-check.sh --only adapted,markers .` | 2026-10-20T20:58:42Z | 2026-10-20T20:58:46Z | 0 | setup-check: markers OK |
| `git diff --check` | 2026-10-20T20:58:46Z | 2026-10-20T20:58:46Z | 0 | No diagnostic. |

Final task-check SHA-256: `a42c1ca3c241e4ae5f9192058c0e46efb73d5e3bfefd3246518e1f6a5b2aedef`.

Final task-check output:

```text
FAIL: canonical margin sentence in the tolerance section
ok: unchanged runner: docs/tests/freeze-check/run.sh
ok: unchanged runner: docs/tests/freeze-check/run-all.sh
ok: all nine baseline paths stay in the restore list
FAIL: kept checker requires the pull_request event
FAIL: one literal PR number and only the checker path can be kept
pending: supply the actual migration PR number before S2
FAIL: event guard receives the actual GitHub event name
FAIL: number guard receives the actual pull request number
FAIL: only the guarded checker skips restoration; both runners and other paths restore for every event
ok: restoration still reads the default branch of this repository
ok: default branch still comes from the repository event
FAIL: exception comment names the task, literal PR and removal issue
```
