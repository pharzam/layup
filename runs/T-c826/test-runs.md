# T-c826 test runs

Stage A below records amendment 1 before checker, runner, workflow or canonical-text
implementation. Stage B records the numbered workflow exception; the checker
remains unchanged. The binding review sets 1,300 lines / 170 files, close-out inside,
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

## Stage B — S2 numbered CI exception

Tested HEAD: `c456dc0df4d4ab9e32e6ad1c31a19deeb50aa700` (committed
Stage A), with only the workflow and evidence records changed during this stage.
The binding Stage B input identifies PR #453 as the PR that closes #449 and
supplies R4 approvals 6301028661 and 6301084598. Removal issue: #452 (T-szvp).
The read-only forge query below failed to connect; no fresh PR-body verification
or hosted CI result is claimed here. The coordinating session will push this
stage and record the PR's required red `nested-checkout-check` run before S3.

The same direct timestamp/status method as Stage A was used. Scratch logs and
disposable fixture repositories were confined to
`/home/layup/.cache/layup-session/tmp-green/T-c826-stage-b` via `TMPDIR`.
The setup command printed a running step and a heartbeat every five seconds.
No repository commit or push was attempted.

| Command | Start UTC | End UTC | Exit | Observed result |
| --- | --- | --- | --- | --- |
| `git rev-parse HEAD` | 2026-10-20T21:02:31Z | 2026-10-20T21:02:31Z | 0 | Confirmed the committed Stage A HEAD above. |
| `gh pr view 453 --repo pharzam/layup --json number,body,url` | 2026-10-20T21:02:31Z | 2026-10-20T21:02:31Z | 1 | Network failure: error connecting to api.github.com; supplied PR identity used, body verification deferred to coordinator. |
| `sh runs/T-c826/checks.sh 453` | 2026-10-20T21:02:31Z | 2026-10-20T21:02:31Z | 1 | Before workflow edits: all numbered-exception assertions red; baseline runner, path-list and default-source controls green. Canonical margin remains red for S4. |
| `sh runs/T-c826/checks.sh 453` | 2026-10-20T21:03:01Z | 2026-10-20T21:03:01Z | 1 | Intermediate edit: four failures, including margin, event/number environment bindings and the exception comment. The scoped assertions caught those fields in an earlier job; moved them into nested-checkout-check. |
| `sh docs/setup/setup-check.sh` | 2026-10-20T21:03:01Z | 2026-10-20T21:03:24Z | 0 | Intermediate workflow: all setup checks OK, including ci; 138 discipline tests passed and all kit linters OK. Does not establish that event bindings were in the right job. |
| `sh docs/tests/freeze-check/run-all.sh` | 2026-10-20T21:03:01Z | 2026-10-20T21:03:09Z | 1 | Unchanged checker and runners: freeze-check: 7 passed, 74 failed. EXPECT mismatches reproduce Stage A; no fixture-construction error. |
| `sh runs/T-c826/checks.sh 453` | 2026-10-20T21:03:51Z | 2026-10-20T21:03:51Z | 1 | Final workflow: all 12 CI/preservation assertions pass. Sole failure is the canonical margin sentence, intentionally reserved for S4. |
| `sh docs/setup/setup-check.sh` | 2026-10-20T21:03:51Z | 2026-10-20T21:04:15Z | 0 | Final workflow: setup-check: ci OK; all other setup checks OK, 138 discipline tests passed, all kit linters OK; exit 0. |

### PR-number assertions before implementation

```text
FAIL: canonical margin sentence in the tolerance section
ok: unchanged runner: docs/tests/freeze-check/run.sh
ok: unchanged runner: docs/tests/freeze-check/run-all.sh
ok: all nine baseline paths stay in the restore list
FAIL: kept checker requires the pull_request event
FAIL: one literal PR number and only the checker path can be kept
FAIL: exception matches the supplied actual migration PR
FAIL: event guard receives the actual GitHub event name
FAIL: number guard receives the actual pull request number
FAIL: only the guarded checker skips restoration; both runners and other paths restore for every event
ok: restoration still reads the default branch of this repository
ok: default branch still comes from the repository event
FAIL: exception comment names the task, literal PR and removal issue
```

### Assertions after correcting the workflow scope

```text
FAIL: canonical margin sentence in the tolerance section
ok: unchanged runner: docs/tests/freeze-check/run.sh
ok: unchanged runner: docs/tests/freeze-check/run-all.sh
ok: all nine baseline paths stay in the restore list
ok: kept checker requires the pull_request event
ok: one literal PR number and only the checker path can be kept
ok: exception matches the supplied actual migration PR
ok: event guard receives the actual GitHub event name
ok: number guard receives the actual pull request number
ok: only the guarded checker skips restoration; both runners and other paths restore for every event
ok: restoration still reads the default branch of this repository
ok: default branch still comes from the repository event
ok: exception comment names the task, literal PR and removal issue
```

The intermediate assertions were not credited as green; the final invocation
above reads the corrected `nested-checkout-check` restore step. The workflow
comment names both R4 approvals. Its only keep condition requires event
`pull_request`, literal number `453` and path `docs/tests/freeze-check.sh`.
All nine paths remain in the original restore list. Both runners and every
other listed script retain the default-branch restoration and log message.

### Final full setup result

```text
setup-check: pin OK
setup-check: kit-history OK
setup-check: facts OK
setup-check: onboarding OK
setup-check: glossary OK
setup-check: guardrails OK
setup-check: markers OK
setup-check: adapted OK
setup-check: ci OK
setup-check: protection OK
setup-check: identity OK
setup-check: procedure OK
run-discipline-tests: 138 passed, 0 failed
adr-lint: OK
prd-lint: OK
link-lint: OK  2991 links resolved
setup-check: kit-linters OK
```

### Local red suite and preservation

The suite still reports `freeze-check: 7 passed, 74 failed` (exit 1), using the
unchanged checker through the unchanged runner. For example, the actual K5 line
in `bad-K1-blank-eof` remains:

```text
ok K5: projection: figure 72/3, maximum 100/10, margin 28/7 (the figure at landing is still measured at landing)
```

The corresponding S1 EXPECT requires figure 2/1, tolerance 60/4, limit 160/14
and margin 158/13. Evidence cases still encounter the unsupported `--evidence`
option; no validator is credited by this red result. This local run is not the
required hosted red run. The next-stage evidence must include its URL, tested
SHA, checker identity, suite failure and these log lines:

```text
kept as written for PR 453: docs/tests/freeze-check.sh
restored from main: docs/tests/freeze-check/run.sh
restored from main: docs/tests/freeze-check/run-all.sh
```

SHA-256 comparisons before and after the workflow edit confirmed that the checker,
both runners, `runs/T-c826/checks.sh` and `docs/setup/setup-check.sh` were unchanged:

| Unchanged file | SHA-256 before and after S2 |
| --- | --- |
| `docs/tests/freeze-check.sh` | `a847cc17bffd2ea5957f65159068d595e43b34ac5a43adc5113c06355a522b63` |
| `docs/tests/freeze-check/run.sh` | `9443c510830f705fdc3f81732270db9f406cb0302cb2f27e18b414392899a4b0` |
| `docs/tests/freeze-check/run-all.sh` | `b892114bdc1719b724c46ba2242941e1f8bcf00d3945ce95877dd265e7a26bb8` |
| `runs/T-c826/checks.sh` | `a42c1ca3c241e4ae5f9192058c0e46efb73d5e3bfefd3246518e1f6a5b2aedef` |
| `docs/setup/setup-check.sh` | `caad7e8587a75164dc8966cb6afe791f896f9eeea88cebb9aa76d1dce34f75a6` |

### Final whitespace and scope check

| Command | Start UTC | End UTC | Exit | Observed result |
| --- | --- | --- | --- | --- |
| `git diff --check` | 2026-10-20T21:05:12Z | 2026-10-20T21:05:12Z | 0 | No diagnostics. |

Only `.github/workflows/ci.yml`, `runs/T-c826/conditions.md` and this record
are modified. Each has no trailing spaces and exactly one final newline. The
checker, both runners, task check and setup check retain their recorded hashes.
No task file, completed-log line or fixture was changed in Stage B.

## Stage C — S2 hosted red evidence, before checker implementation

The supplied CI evidence was corroborated by reading GitHub job 118887031319's
steps and decoded log through the connected GitHub tools before editing the
checker. [Run 39497659349, nested-checkout-check](https://github.com/pharzam/layup/actions/runs/39497659349/job/118887031319)
concluded **failure**. Checkout tested merge SHA
`1470b85f4f34cc30c40e804c7ab204e03d56c7ca`, merging Stage B head
`fc3a65a84e6f7b62353a8a1081769dda70266553` into
`d27b33180a1aa82c23f3ab66e757e06de356eb35`. The kept checker has SHA-256
`a847cc17bffd2ea5957f65159068d595e43b34ac5a43adc5113c06355a522b63`.

The log directly records these results (hosted times are log timestamps, not
local `date` measurements):

```text
2026-10-20T21:07:12.1826151Z restored from main: docs/adr/adr-lint.sh
2026-10-20T21:07:12.1845169Z restored from main: docs/prd/prd-lint.sh
2026-10-20T21:07:12.1855071Z restored from main: docs/links/link-lint.sh
2026-10-20T21:07:12.1870459Z restored from main: docs/tests/run-discipline-tests.sh
2026-10-20T21:07:12.1890997Z restored from main: docs/ci/pr-link-lint.sh
2026-10-20T21:07:12.1911758Z restored from main: docs/ci/review-record-lint.sh
2026-10-20T21:07:12.1912679Z kept as written for PR 453: docs/tests/freeze-check.sh
2026-10-20T21:07:12.1932099Z restored from main: docs/tests/freeze-check/run.sh
2026-10-20T21:07:12.1947524Z restored from main: docs/tests/freeze-check/run-all.sh
2026-10-20T21:07:12.5889937Z Run sh docs/tests/freeze-check/run-all.sh
2026-10-20T21:07:22.6598294Z freeze-check: 7 passed, 74 failed
2026-10-20T21:07:22.6612549Z Process completed with exit code 1.
```

The restore and nested-checkout fixture steps succeeded; the freeze suite step
failed. This is the intended complete-output mismatch against the unchanged
checker, rather than a checkout, YAML or environment failure. No hosted green is
claimed in this uncommitted build stage; the coordinating session must obtain it
after committing and pushing this implementation.

## Stage C local verification log

Source: Stage B parent `fc3a65a84e6f7b62353a8a1081769dda70266553` plus
this stage's uncommitted checker and documentation. The checker SHA-256 is
`04c5918711f4f07e5fdd423424cbf14b99a43aefd77e76603c183b98cc96e2cf`.
The canonical sentence assertion was observed red before the documentation edit;
all S1 expectations remained unchanged. Both green suites report **81 passed,
0 failed**, under `sh` and `bash --posix` respectively. No fixture, runner,
workflow, setup-check, review-record parser or task-check source changed in Stage C.
Existing attributes preserve the CRLF fixtures; no new attribute is needed.

All command rows below record actual direct subprocess exits, with UTC endpoints
read immediately before and after using `date -u +%FT%TZ`. Output is captured
without a pipeline. `$W` denotes this worktree and `$C` denotes
`/home/layup/.cache/layup-session/tmp-green/T-c826-stage-c`; paths in commands
are abbreviated only for readability. Every run inherits `TMPDIR=$C`.
Checker mutants run the full unchanged suite from a fresh copy of `docs/tests`
under `$C/Mxx`, because run.sh has no alternate-checker hook. Each restored
control runs its named case after a verified byte-for-byte restore. Task-assertion
copies use the original read-only Git directory only for `git show` of the two
baseline runners. Scratch drivers and disposable trees stay outside the worktree; this record holds observations only.

The mutation driver stopped once before constructing M37: a requested workflow
edit matched six jobs, so its uniqueness assertion refused the edit. No test ran
for that attempt and no kill is credited. The resumed driver scopes edits to
`nested-checkout-check`; M53–M66 record its actual assertion runs. M51 and M52
are the two independent limit-addition mutants. Labels retain their actual run IDs.

| Run | Command | Working directory | Start UTC | End UTC | Exit | Observed result |
| --- | --- | --- | --- | --- | --- | --- |
| before-docs | `sh runs/T-c826/checks.sh 453` | `$W` | 2026-10-20T21:17:20Z | 2026-10-20T21:17:21Z | 1 | FAIL: canonical margin sentence; every other task assertion passed. |
| s3-suite | `sh docs/tests/freeze-check/run-all.sh` | `$W` | 2026-10-20T21:17:47Z | 2026-10-20T21:17:55Z | 0 | freeze-check: 81 passed, 0 failed |
| s5-task-check | `sh runs/T-c826/checks.sh 453` | `$W` | 2026-10-20T21:18:31Z | 2026-10-20T21:18:31Z | 0 | All 13 task assertions passed with actual PR 453. |
| s5-sh-suite | `sh docs/tests/freeze-check/run-all.sh` | `$W` | 2026-10-20T21:18:31Z | 2026-10-20T21:18:39Z | 0 | freeze-check: 81 passed, 0 failed |
| s5-bash-suite | `bash --posix docs/tests/freeze-check/run-all.sh` | `$W` | 2026-10-20T21:19:01Z | 2026-10-20T21:19:09Z | 0 | freeze-check: 81 passed, 0 failed |
| M01-syntax | `sh -n $C/M01/docs/tests/freeze-check.sh` | `$C/M01` | 2026-10-20T21:20:04Z | 2026-10-20T21:20:04Z | 0 | Valid POSIX shell syntax; no diagnostics. |
| M01-mutant | `sh docs/tests/freeze-check/run-all.sh` | `$C/M01` | 2026-10-20T21:20:04Z | 2026-10-20T21:20:11Z | 1 | freeze-check: 73 passed, 8 failed; good-K5-evidence-kinds failed as recorded in mutants.md. |
| M01-restored | `sh docs/tests/freeze-check/run.sh docs/tests/freeze-check/good-K5-evidence-kinds` | `$C/M01` | 2026-10-20T21:20:11Z | 2026-10-20T21:20:11Z | 0 | Restored hash matched; named control green. |
| M02-syntax | `sh -n $C/M02/docs/tests/freeze-check.sh` | `$C/M02` | 2026-10-20T21:20:12Z | 2026-10-20T21:20:12Z | 0 | Valid POSIX shell syntax; no diagnostics. |
| M02-mutant | `sh docs/tests/freeze-check/run-all.sh` | `$C/M02` | 2026-10-20T21:20:12Z | 2026-10-20T21:20:20Z | 1 | freeze-check: 80 passed, 1 failed; good-K5-evidence-kinds failed as recorded in mutants.md. |
| M02-restored | `sh docs/tests/freeze-check/run.sh docs/tests/freeze-check/good-K5-evidence-kinds` | `$C/M02` | 2026-10-20T21:20:20Z | 2026-10-20T21:20:20Z | 0 | Restored hash matched; named control green. |
| M03-syntax | `sh -n $C/M03/docs/tests/freeze-check.sh` | `$C/M03` | 2026-10-20T21:20:20Z | 2026-10-20T21:20:20Z | 0 | Valid POSIX shell syntax; no diagnostics. |
| M03-mutant | `sh docs/tests/freeze-check/run-all.sh` | `$C/M03` | 2026-10-20T21:20:20Z | 2026-10-20T21:20:29Z | 1 | freeze-check: 80 passed, 1 failed; good-K5-evidence-kinds failed as recorded in mutants.md. |
| M03-restored | `sh docs/tests/freeze-check/run.sh docs/tests/freeze-check/good-K5-evidence-kinds` | `$C/M03` | 2026-10-20T21:20:29Z | 2026-10-20T21:20:29Z | 0 | Restored hash matched; named control green. |
| M04-syntax | `sh -n $C/M04/docs/tests/freeze-check.sh` | `$C/M04` | 2026-10-20T21:20:29Z | 2026-10-20T21:20:29Z | 0 | Valid POSIX shell syntax; no diagnostics. |
| M04-mutant | `sh docs/tests/freeze-check/run-all.sh` | `$C/M04` | 2026-10-20T21:20:29Z | 2026-10-20T21:20:37Z | 1 | freeze-check: 80 passed, 1 failed; good-K5-evidence-kinds failed as recorded in mutants.md. |
| M04-restored | `sh docs/tests/freeze-check/run.sh docs/tests/freeze-check/good-K5-evidence-kinds` | `$C/M04` | 2026-10-20T21:20:37Z | 2026-10-20T21:20:37Z | 0 | Restored hash matched; named control green. |
| M05-syntax | `sh -n $C/M05/docs/tests/freeze-check.sh` | `$C/M05` | 2026-10-20T21:20:37Z | 2026-10-20T21:20:37Z | 0 | Valid POSIX shell syntax; no diagnostics. |
| M05-mutant | `sh docs/tests/freeze-check/run-all.sh` | `$C/M05` | 2026-10-20T21:20:37Z | 2026-10-20T21:20:46Z | 1 | freeze-check: 74 passed, 7 failed; good-K5-evidence-kinds failed as recorded in mutants.md. |
| M05-restored | `sh docs/tests/freeze-check/run.sh docs/tests/freeze-check/good-K5-evidence-kinds` | `$C/M05` | 2026-10-20T21:20:46Z | 2026-10-20T21:20:46Z | 0 | Restored hash matched; named control green. |
| M06-syntax | `sh -n $C/M06/docs/tests/freeze-check.sh` | `$C/M06` | 2026-10-20T21:20:46Z | 2026-10-20T21:20:46Z | 0 | Valid POSIX shell syntax; no diagnostics. |
| adr | `sh docs/adr/adr-lint.sh` | `$W` | 2026-10-20T21:20:50Z | 2026-10-20T21:20:51Z | 0 | adr-lint: OK |
| prd | `sh docs/prd/prd-lint.sh` | `$W` | 2026-10-20T21:20:51Z | 2026-10-20T21:20:51Z | 0 | prd-lint: OK |
| M06-mutant | `sh docs/tests/freeze-check/run-all.sh` | `$C/M06` | 2026-10-20T21:20:46Z | 2026-10-20T21:20:54Z | 1 | freeze-check: 74 passed, 7 failed; good-K5-evidence-kinds failed as recorded in mutants.md. |
| M06-restored | `sh docs/tests/freeze-check/run.sh docs/tests/freeze-check/good-K5-evidence-kinds` | `$C/M06` | 2026-10-20T21:20:54Z | 2026-10-20T21:20:55Z | 0 | Restored hash matched; named control green. |
| M07-syntax | `sh -n $C/M07/docs/tests/freeze-check.sh` | `$C/M07` | 2026-10-20T21:20:55Z | 2026-10-20T21:20:55Z | 0 | Valid POSIX shell syntax; no diagnostics. |
| M07-mutant | `sh docs/tests/freeze-check/run-all.sh` | `$C/M07` | 2026-10-20T21:20:55Z | 2026-10-20T21:21:03Z | 1 | freeze-check: 78 passed, 3 failed; good-K5-evidence-inputs failed as recorded in mutants.md. |
| M07-restored | `sh docs/tests/freeze-check/run.sh docs/tests/freeze-check/good-K5-evidence-inputs` | `$C/M07` | 2026-10-20T21:21:03Z | 2026-10-20T21:21:03Z | 0 | Restored hash matched; named control green. |
| M08-syntax | `sh -n $C/M08/docs/tests/freeze-check.sh` | `$C/M08` | 2026-10-20T21:21:03Z | 2026-10-20T21:21:03Z | 0 | Valid POSIX shell syntax; no diagnostics. |
| links | `sh docs/links/link-lint.sh` | `$W` | 2026-10-20T21:20:51Z | 2026-10-20T21:21:03Z | 0 | link-lint: OK  2992 links resolved |
| discipline | `sh docs/tests/run-discipline-tests.sh` | `$W` | 2026-10-20T21:21:03Z | 2026-10-20T21:21:06Z | 0 | run-discipline-tests: 138 passed, 0 failed |
| M08-mutant | `sh docs/tests/freeze-check/run-all.sh` | `$C/M08` | 2026-10-20T21:21:03Z | 2026-10-20T21:21:11Z | 1 | freeze-check: 76 passed, 5 failed; good-K5-evidence-added failed as recorded in mutants.md. |
| M08-restored | `sh docs/tests/freeze-check/run.sh docs/tests/freeze-check/good-K5-evidence-added` | `$C/M08` | 2026-10-20T21:21:11Z | 2026-10-20T21:21:11Z | 0 | Restored hash matched; named control green. |
| M09-syntax | `sh -n $C/M09/docs/tests/freeze-check.sh` | `$C/M09` | 2026-10-20T21:21:11Z | 2026-10-20T21:21:11Z | 0 | Valid POSIX shell syntax; no diagnostics. |
| M09-mutant | `sh docs/tests/freeze-check/run-all.sh` | `$C/M09` | 2026-10-20T21:21:11Z | 2026-10-20T21:21:19Z | 1 | freeze-check: 80 passed, 1 failed; good-K5-evidence-deleted failed as recorded in mutants.md. |
| M09-restored | `sh docs/tests/freeze-check/run.sh docs/tests/freeze-check/good-K5-evidence-deleted` | `$C/M09` | 2026-10-20T21:21:19Z | 2026-10-20T21:21:19Z | 0 | Restored hash matched; named control green. |
| M10-syntax | `sh -n $C/M10/docs/tests/freeze-check.sh` | `$C/M10` | 2026-10-20T21:21:19Z | 2026-10-20T21:21:19Z | 0 | Valid POSIX shell syntax; no diagnostics. |
| M10-mutant | `sh docs/tests/freeze-check/run-all.sh` | `$C/M10` | 2026-10-20T21:21:19Z | 2026-10-20T21:21:27Z | 1 | freeze-check: 75 passed, 6 failed; good-K5-evidence-modified failed as recorded in mutants.md. |
| M10-restored | `sh docs/tests/freeze-check/run.sh docs/tests/freeze-check/good-K5-evidence-modified` | `$C/M10` | 2026-10-20T21:21:27Z | 2026-10-20T21:21:27Z | 0 | Restored hash matched; named control green. |
| M11-syntax | `sh -n $C/M11/docs/tests/freeze-check.sh` | `$C/M11` | 2026-10-20T21:21:27Z | 2026-10-20T21:21:27Z | 0 | Valid POSIX shell syntax; no diagnostics. |
| setup | `sh docs/setup/setup-check.sh` | `$W` | 2026-10-20T21:21:06Z | 2026-10-20T21:21:27Z | 0 | setup-check: kit-linters OK |
| whitespace | `git diff --check` | `$W` | 2026-10-20T21:21:27Z | 2026-10-20T21:21:27Z | 0 | No diagnostics. |
| M11-mutant | `sh docs/tests/freeze-check/run-all.sh` | `$C/M11` | 2026-10-20T21:21:27Z | 2026-10-20T21:21:36Z | 1 | freeze-check: 79 passed, 2 failed; good-K5-evidence-modified failed as recorded in mutants.md. |
| M11-restored | `sh docs/tests/freeze-check/run.sh docs/tests/freeze-check/good-K5-evidence-modified` | `$C/M11` | 2026-10-20T21:21:36Z | 2026-10-20T21:21:36Z | 0 | Restored hash matched; named control green. |
| M12-syntax | `sh -n $C/M12/docs/tests/freeze-check.sh` | `$C/M12` | 2026-10-20T21:21:36Z | 2026-10-20T21:21:36Z | 0 | Valid POSIX shell syntax; no diagnostics. |
| M12-mutant | `sh docs/tests/freeze-check/run-all.sh` | `$C/M12` | 2026-10-20T21:21:36Z | 2026-10-20T21:21:44Z | 1 | freeze-check: 79 passed, 2 failed; bad-usage-evidence-other-task failed as recorded in mutants.md. |
| M12-restored | `sh docs/tests/freeze-check/run.sh docs/tests/freeze-check/bad-usage-evidence-other-task` | `$C/M12` | 2026-10-20T21:21:44Z | 2026-10-20T21:21:44Z | 0 | Restored hash matched; named control green. |
| M13-syntax | `sh -n $C/M13/docs/tests/freeze-check.sh` | `$C/M13` | 2026-10-20T21:21:44Z | 2026-10-20T21:21:44Z | 0 | Valid POSIX shell syntax; no diagnostics. |
| T-3ff9-clone | `git clone --quiet --shared --no-checkout $W $C/replay-T-3ff9` | `$W` | 2026-10-20T21:21:49Z | 2026-10-20T21:21:49Z | 0 | Disposable clone operation succeeded. |
| T-3ff9-checkout | `git checkout --quiet --detach 0dfac539061cda59b233ab222dca76d111c72db8` | `$C/replay-T-3ff9` | 2026-10-20T21:21:49Z | 2026-10-20T21:21:49Z | 0 | Disposable clone operation succeeded. |
| T-3ff9-consumers | `git grep -n -F runs/T-3ff9 -- ':!runs/T-3ff9/*.md' ':!runs/T-3ff9/*.tsv'` | `$C/replay-T-3ff9` | 2026-10-20T21:21:49Z | 2026-10-20T21:21:49Z | 0 | Task-record citations found; no execution consumer at the named task paths. |
| T-3ff9-records | `git show 0dfac539061cda59b233ab222dca76d111c72db8:runs/T-3ff9/test-runs.md 0dfac539061cda59b233ab222dca76d111c72db8:runs/T-3ff9/conditions.md 0dfac539061cda59b233ab222dca76d111c72db8:runs/T-3ff9/semantic.md 0dfac539061cda59b233ab222dca76d111c72db8:runs/T-3ff9/mentions.tsv` | `$C/replay-T-3ff9` | 2026-10-20T21:21:49Z | 2026-10-20T21:21:49Z | 0 | Committed contents inspected for evidence role and consumers. |
| T-3ff9-check-source | `git show 0dfac539061cda59b233ab222dca76d111c72db8:runs/T-3ff9/checks.sh` | `$C/replay-T-3ff9` | 2026-10-20T21:21:49Z | 2026-10-20T21:21:49Z | 0 | Committed contents inspected for evidence role and consumers. |
| T-3ff9-gross | `git diff --no-renames --numstat df9bfacf33054741564dbf0254a741f0d4be1e0f HEAD` | `$C/replay-T-3ff9` | 2026-10-20T21:21:49Z | 2026-10-20T21:21:49Z | 0 | Committed numstat captured; arithmetic in real-heads.md. |
| T-3ff9-projection | `sh $W/docs/tests/freeze-check.sh --task T-3ff9 --base df9bfacf33054741564dbf0254a741f0d4be1e0f --budget 650/20 --evidence runs/T-3ff9/test-runs.md --evidence runs/T-3ff9/conditions.md --evidence runs/T-3ff9/semantic.md --evidence runs/T-3ff9/mentions.tsv` | `$C/replay-T-3ff9` | 2026-10-20T21:21:49Z | 2026-10-20T21:21:50Z | 0 | ok K5: projection: figure 318/12, maximum 650/20, tolerance 130/4, limit 780/24, margin 462/12 (the figure at landing is still measured at landing) |
| T-b6fn-clone | `git clone --quiet --shared --no-checkout $W $C/replay-T-b6fn` | `$W` | 2026-10-20T21:21:50Z | 2026-10-20T21:21:50Z | 0 | Disposable clone operation succeeded. |
| T-b6fn-checkout | `git checkout --quiet --detach 78386339e297ed1e1a4bd8ff35d2af5189bbfd6f` | `$C/replay-T-b6fn` | 2026-10-20T21:21:50Z | 2026-10-20T21:21:50Z | 0 | Disposable clone operation succeeded. |
| T-b6fn-consumers | `git grep -n -F runs/T-b6fn -- ':!runs/T-b6fn/*.md' ':!runs/T-b6fn/*.tsv'` | `$C/replay-T-b6fn` | 2026-10-20T21:21:50Z | 2026-10-20T21:21:50Z | 0 | Task-record citations found; no execution consumer at the named task paths. |
| T-b6fn-records | `git show 78386339e297ed1e1a4bd8ff35d2af5189bbfd6f:runs/T-b6fn/test-runs.md 78386339e297ed1e1a4bd8ff35d2af5189bbfd6f:runs/T-b6fn/conditions.md 78386339e297ed1e1a4bd8ff35d2af5189bbfd6f:runs/T-b6fn/real-heads.md` | `$C/replay-T-b6fn` | 2026-10-20T21:21:50Z | 2026-10-20T21:21:50Z | 0 | Committed contents inspected for evidence role and consumers. |
| T-b6fn-gross | `git diff --no-renames --numstat 99ff86235b726f8842120dafb2356b1aff035238 HEAD` | `$C/replay-T-b6fn` | 2026-10-20T21:21:50Z | 2026-10-20T21:21:50Z | 0 | Committed numstat captured; arithmetic in real-heads.md. |
| T-b6fn-projection | `sh $W/docs/tests/freeze-check.sh --task T-b6fn --base 99ff86235b726f8842120dafb2356b1aff035238 --budget 1300/170 --evidence runs/T-b6fn/test-runs.md --evidence runs/T-b6fn/conditions.md --evidence runs/T-b6fn/real-heads.md` | `$C/replay-T-b6fn` | 2026-10-20T21:21:50Z | 2026-10-20T21:21:51Z | 1 | FAIL K5: -:0: projection: figure 2023/235, maximum 1300/170, tolerance 260/4, limit 1560/174, margin -463/-61 (the figure at landing is still measured at landing) |
| M13-mutant | `sh docs/tests/freeze-check/run-all.sh` | `$C/M13` | 2026-10-20T21:21:44Z | 2026-10-20T21:21:53Z | 1 | freeze-check: 80 passed, 1 failed; bad-usage-evidence-task-prefix failed as recorded in mutants.md. |
| M13-restored | `sh docs/tests/freeze-check/run.sh docs/tests/freeze-check/bad-usage-evidence-task-prefix` | `$C/M13` | 2026-10-20T21:21:53Z | 2026-10-20T21:21:53Z | 0 | Restored hash matched; named control green. |
| M14-syntax | `sh -n $C/M14/docs/tests/freeze-check.sh` | `$C/M14` | 2026-10-20T21:21:53Z | 2026-10-20T21:21:53Z | 0 | Valid POSIX shell syntax; no diagnostics. |
| M14-mutant | `sh docs/tests/freeze-check/run-all.sh` | `$C/M14` | 2026-10-20T21:21:53Z | 2026-10-20T21:22:01Z | 1 | freeze-check: 80 passed, 1 failed; bad-usage-evidence-traversal failed as recorded in mutants.md. |
| M14-restored | `sh docs/tests/freeze-check/run.sh docs/tests/freeze-check/bad-usage-evidence-traversal` | `$C/M14` | 2026-10-20T21:22:01Z | 2026-10-20T21:22:01Z | 0 | Restored hash matched; named control green. |
| M15-syntax | `sh -n $C/M15/docs/tests/freeze-check.sh` | `$C/M15` | 2026-10-20T21:22:01Z | 2026-10-20T21:22:01Z | 0 | Valid POSIX shell syntax; no diagnostics. |
| M15-mutant | `sh docs/tests/freeze-check/run-all.sh` | `$C/M15` | 2026-10-20T21:22:01Z | 2026-10-20T21:22:10Z | 1 | freeze-check: 80 passed, 1 failed; bad-usage-evidence-suffix failed as recorded in mutants.md. |
| M15-restored | `sh docs/tests/freeze-check/run.sh docs/tests/freeze-check/bad-usage-evidence-suffix` | `$C/M15` | 2026-10-20T21:22:10Z | 2026-10-20T21:22:10Z | 0 | Restored hash matched; named control green. |
| M16-syntax | `sh -n $C/M16/docs/tests/freeze-check.sh` | `$C/M16` | 2026-10-20T21:22:10Z | 2026-10-20T21:22:10Z | 0 | Valid POSIX shell syntax; no diagnostics. |
| M16-mutant | `sh docs/tests/freeze-check/run-all.sh` | `$C/M16` | 2026-10-20T21:22:10Z | 2026-10-20T21:22:18Z | 1 | freeze-check: 80 passed, 1 failed; bad-usage-evidence-missing failed as recorded in mutants.md. |
| M16-restored | `sh docs/tests/freeze-check/run.sh docs/tests/freeze-check/bad-usage-evidence-missing` | `$C/M16` | 2026-10-20T21:22:18Z | 2026-10-20T21:22:18Z | 0 | Restored hash matched; named control green. |
| M17-syntax | `sh -n $C/M17/docs/tests/freeze-check.sh` | `$C/M17` | 2026-10-20T21:22:18Z | 2026-10-20T21:22:18Z | 0 | Valid POSIX shell syntax; no diagnostics. |
| M17-mutant | `sh docs/tests/freeze-check/run-all.sh` | `$C/M17` | 2026-10-20T21:22:18Z | 2026-10-20T21:22:27Z | 1 | freeze-check: 80 passed, 1 failed; good-K5-evidence-literal failed as recorded in mutants.md. |
| M17-restored | `sh docs/tests/freeze-check/run.sh docs/tests/freeze-check/good-K5-evidence-literal` | `$C/M17` | 2026-10-20T21:22:27Z | 2026-10-20T21:22:27Z | 0 | Restored hash matched; named control green. |
| M18-syntax | `sh -n $C/M18/docs/tests/freeze-check.sh` | `$C/M18` | 2026-10-20T21:22:27Z | 2026-10-20T21:22:27Z | 0 | Valid POSIX shell syntax; no diagnostics. |
| M18-mutant | `sh docs/tests/freeze-check/run-all.sh` | `$C/M18` | 2026-10-20T21:22:27Z | 2026-10-20T21:22:35Z | 1 | freeze-check: 80 passed, 1 failed; good-K5-evidence-duplicate failed as recorded in mutants.md. |
| M18-restored | `sh docs/tests/freeze-check/run.sh docs/tests/freeze-check/good-K5-evidence-duplicate` | `$C/M18` | 2026-10-20T21:22:35Z | 2026-10-20T21:22:36Z | 0 | Restored hash matched; named control green. |
| M19-syntax | `sh -n $C/M19/docs/tests/freeze-check.sh` | `$C/M19` | 2026-10-20T21:22:36Z | 2026-10-20T21:22:36Z | 0 | Valid POSIX shell syntax; no diagnostics. |
| M19-mutant | `sh docs/tests/freeze-check/run-all.sh` | `$C/M19` | 2026-10-20T21:22:36Z | 2026-10-20T21:22:44Z | 1 | freeze-check: 21 passed, 60 failed; good-K5-line-limit failed as recorded in mutants.md. |
| M19-restored | `sh docs/tests/freeze-check/run.sh docs/tests/freeze-check/good-K5-line-limit` | `$C/M19` | 2026-10-20T21:22:44Z | 2026-10-20T21:22:44Z | 0 | Restored hash matched; named control green. |
| M20-syntax | `sh -n $C/M20/docs/tests/freeze-check.sh` | `$C/M20` | 2026-10-20T21:22:44Z | 2026-10-20T21:22:44Z | 0 | Valid POSIX shell syntax; no diagnostics. |
| M20-mutant | `sh docs/tests/freeze-check/run-all.sh` | `$C/M20` | 2026-10-20T21:22:44Z | 2026-10-20T21:22:53Z | 1 | freeze-check: 78 passed, 3 failed; good-K5-percentage failed as recorded in mutants.md. |
| M20-restored | `sh docs/tests/freeze-check/run.sh docs/tests/freeze-check/good-K5-percentage` | `$C/M20` | 2026-10-20T21:22:53Z | 2026-10-20T21:22:53Z | 0 | Restored hash matched; named control green. |
| M21-syntax | `sh -n $C/M21/docs/tests/freeze-check.sh` | `$C/M21` | 2026-10-20T21:22:53Z | 2026-10-20T21:22:53Z | 0 | Valid POSIX shell syntax; no diagnostics. |
| M21-mutant | `sh docs/tests/freeze-check/run-all.sh` | `$C/M21` | 2026-10-20T21:22:53Z | 2026-10-20T21:23:01Z | 1 | freeze-check: 80 passed, 1 failed; good-K5-ceiling failed as recorded in mutants.md. |
| M21-restored | `sh docs/tests/freeze-check/run.sh docs/tests/freeze-check/good-K5-ceiling` | `$C/M21` | 2026-10-20T21:23:01Z | 2026-10-20T21:23:01Z | 0 | Restored hash matched; named control green. |
| M22-syntax | `sh -n $C/M22/docs/tests/freeze-check.sh` | `$C/M22` | 2026-10-20T21:23:01Z | 2026-10-20T21:23:01Z | 0 | Valid POSIX shell syntax; no diagnostics. |
| M22-mutant | `sh docs/tests/freeze-check/run-all.sh` | `$C/M22` | 2026-10-20T21:23:01Z | 2026-10-20T21:23:10Z | 1 | freeze-check: 80 passed, 1 failed; good-K5-truncation failed as recorded in mutants.md. |
| M22-restored | `sh docs/tests/freeze-check/run.sh docs/tests/freeze-check/good-K5-truncation` | `$C/M22` | 2026-10-20T21:23:10Z | 2026-10-20T21:23:10Z | 0 | Restored hash matched; named control green. |
| M23-syntax | `sh -n $C/M23/docs/tests/freeze-check.sh` | `$C/M23` | 2026-10-20T21:23:10Z | 2026-10-20T21:23:10Z | 0 | Valid POSIX shell syntax; no diagnostics. |
| M23-mutant | `sh docs/tests/freeze-check/run-all.sh` | `$C/M23` | 2026-10-20T21:23:10Z | 2026-10-20T21:23:19Z | 1 | freeze-check: 18 passed, 63 failed; good-K5-file-limit failed as recorded in mutants.md. |
| M23-restored | `sh docs/tests/freeze-check/run.sh docs/tests/freeze-check/good-K5-file-limit` | `$C/M23` | 2026-10-20T21:23:19Z | 2026-10-20T21:23:19Z | 0 | Restored hash matched; named control green. |
| M24-syntax | `sh -n $C/M24/docs/tests/freeze-check.sh` | `$C/M24` | 2026-10-20T21:23:19Z | 2026-10-20T21:23:19Z | 0 | Valid POSIX shell syntax; no diagnostics. |
| M24-mutant | `sh docs/tests/freeze-check/run-all.sh` | `$C/M24` | 2026-10-20T21:23:19Z | 2026-10-20T21:23:27Z | 1 | freeze-check: 80 passed, 1 failed; bad-K5-over-lines failed as recorded in mutants.md. |
| M24-restored | `sh docs/tests/freeze-check/run.sh docs/tests/freeze-check/bad-K5-over-lines` | `$C/M24` | 2026-10-20T21:23:27Z | 2026-10-20T21:23:27Z | 0 | Restored hash matched; named control green. |
| M25-syntax | `sh -n $C/M25/docs/tests/freeze-check.sh` | `$C/M25` | 2026-10-20T21:23:27Z | 2026-10-20T21:23:27Z | 0 | Valid POSIX shell syntax; no diagnostics. |
| M25-mutant | `sh docs/tests/freeze-check/run-all.sh` | `$C/M25` | 2026-10-20T21:23:27Z | 2026-10-20T21:23:36Z | 1 | freeze-check: 80 passed, 1 failed; good-K5-line-limit failed as recorded in mutants.md. |
| M25-restored | `sh docs/tests/freeze-check/run.sh docs/tests/freeze-check/good-K5-line-limit` | `$C/M25` | 2026-10-20T21:23:36Z | 2026-10-20T21:23:36Z | 0 | Restored hash matched; named control green. |
| M26-syntax | `sh -n $C/M26/docs/tests/freeze-check.sh` | `$C/M26` | 2026-10-20T21:23:36Z | 2026-10-20T21:23:36Z | 0 | Valid POSIX shell syntax; no diagnostics. |
| M26-mutant | `sh docs/tests/freeze-check/run-all.sh` | `$C/M26` | 2026-10-20T21:23:36Z | 2026-10-20T21:23:44Z | 1 | freeze-check: 80 passed, 1 failed; bad-K5-over-files failed as recorded in mutants.md. |
| M26-restored | `sh docs/tests/freeze-check/run.sh docs/tests/freeze-check/bad-K5-over-files` | `$C/M26` | 2026-10-20T21:23:44Z | 2026-10-20T21:23:45Z | 0 | Restored hash matched; named control green. |
| M27-syntax | `sh -n $C/M27/docs/tests/freeze-check.sh` | `$C/M27` | 2026-10-20T21:23:45Z | 2026-10-20T21:23:45Z | 0 | Valid POSIX shell syntax; no diagnostics. |
| M27-mutant | `sh docs/tests/freeze-check/run-all.sh` | `$C/M27` | 2026-10-20T21:23:45Z | 2026-10-20T21:23:53Z | 1 | freeze-check: 80 passed, 1 failed; good-K5-file-limit failed as recorded in mutants.md. |
| M27-restored | `sh docs/tests/freeze-check/run.sh docs/tests/freeze-check/good-K5-file-limit` | `$C/M27` | 2026-10-20T21:23:53Z | 2026-10-20T21:23:53Z | 0 | Restored hash matched; named control green. |
| M28-syntax | `sh -n $C/M28/docs/tests/freeze-check.sh` | `$C/M28` | 2026-10-20T21:23:53Z | 2026-10-20T21:23:53Z | 0 | Valid POSIX shell syntax; no diagnostics. |
| M28-mutant | `sh docs/tests/freeze-check/run-all.sh` | `$C/M28` | 2026-10-20T21:23:53Z | 2026-10-20T21:24:02Z | 1 | freeze-check: 18 passed, 63 failed; good-K5-line-limit failed as recorded in mutants.md. |
| M28-restored | `sh docs/tests/freeze-check/run.sh docs/tests/freeze-check/good-K5-line-limit` | `$C/M28` | 2026-10-20T21:24:02Z | 2026-10-20T21:24:02Z | 0 | Restored hash matched; named control green. |
| M29-syntax | `sh -n $C/M29/docs/tests/freeze-check.sh` | `$C/M29` | 2026-10-20T21:24:02Z | 2026-10-20T21:24:02Z | 0 | Valid POSIX shell syntax; no diagnostics. |
| M29-mutant | `sh docs/tests/freeze-check/run-all.sh` | `$C/M29` | 2026-10-20T21:24:02Z | 2026-10-20T21:24:11Z | 1 | freeze-check: 18 passed, 63 failed; good-K5-file-limit failed as recorded in mutants.md. |
| M29-restored | `sh docs/tests/freeze-check/run.sh docs/tests/freeze-check/good-K5-file-limit` | `$C/M29` | 2026-10-20T21:24:11Z | 2026-10-20T21:24:11Z | 0 | Restored hash matched; named control green. |
| M30-syntax | `sh -n $C/M30/docs/tests/freeze-check.sh` | `$C/M30` | 2026-10-20T21:24:11Z | 2026-10-20T21:24:11Z | 0 | Valid POSIX shell syntax; no diagnostics. |
| M30-mutant | `sh docs/tests/freeze-check/run-all.sh` | `$C/M30` | 2026-10-20T21:24:11Z | 2026-10-20T21:24:19Z | 1 | freeze-check: 18 passed, 63 failed; good-K5-no-closeout failed as recorded in mutants.md. |
| M30-restored | `sh docs/tests/freeze-check/run.sh docs/tests/freeze-check/good-K5-no-closeout` | `$C/M30` | 2026-10-20T21:24:19Z | 2026-10-20T21:24:20Z | 0 | Restored hash matched; named control green. |
| M31-syntax | `sh -n $C/M31/docs/tests/freeze-check.sh` | `$C/M31` | 2026-10-20T21:24:20Z | 2026-10-20T21:24:20Z | 0 | Valid POSIX shell syntax; no diagnostics. |
| M31-mutant | `sh docs/tests/freeze-check/run-all.sh` | `$C/M31` | 2026-10-20T21:24:20Z | 2026-10-20T21:24:28Z | 1 | freeze-check: 18 passed, 63 failed; good-K5-no-closeout failed as recorded in mutants.md. |
| M31-restored | `sh docs/tests/freeze-check/run.sh docs/tests/freeze-check/good-K5-no-closeout` | `$C/M31` | 2026-10-20T21:24:28Z | 2026-10-20T21:24:28Z | 0 | Restored hash matched; named control green. |
| M32-syntax | `sh -n $C/M32/docs/tests/freeze-check.sh` | `$C/M32` | 2026-10-20T21:24:28Z | 2026-10-20T21:24:28Z | 0 | Valid POSIX shell syntax; no diagnostics. |
| M32-mutant | `sh docs/tests/freeze-check/run-all.sh` | `$C/M32` | 2026-10-20T21:24:28Z | 2026-10-20T21:24:37Z | 1 | freeze-check: 80 passed, 1 failed; bad-usage-closeout failed as recorded in mutants.md. |
| M32-restored | `sh docs/tests/freeze-check/run.sh docs/tests/freeze-check/bad-usage-closeout` | `$C/M32` | 2026-10-20T21:24:37Z | 2026-10-20T21:24:37Z | 0 | Restored hash matched; named control green. |
| M33-syntax | `sh -n $C/M33/docs/tests/freeze-check.sh` | `$C/M33` | 2026-10-20T21:24:37Z | 2026-10-20T21:24:37Z | 0 | Valid POSIX shell syntax; no diagnostics. |
| M51-syntax | `sh -n docs/tests/freeze-check.sh` | `$C/M51` | 2026-10-20T21:24:41Z | 2026-10-20T21:24:41Z | 0 | Valid POSIX shell syntax; no diagnostics. |
| M33-mutant | `sh docs/tests/freeze-check/run-all.sh` | `$C/M33` | 2026-10-20T21:24:37Z | 2026-10-20T21:24:45Z | 1 | freeze-check: 80 passed, 1 failed; bad-usage-evidence-no-value failed as recorded in mutants.md. |
| M33-restored | `sh docs/tests/freeze-check/run.sh docs/tests/freeze-check/bad-usage-evidence-no-value` | `$C/M33` | 2026-10-20T21:24:45Z | 2026-10-20T21:24:45Z | 0 | Restored hash matched; named control green. |
| M34-mutant | `sh runs/T-c826/checks.sh 453` | `$C/M34` | 2026-10-20T21:24:45Z | 2026-10-20T21:24:45Z | 1 | Named task assertion failed; canonical margin sentence failed as recorded in mutants.md. |
| M34-restored | `sh runs/T-c826/checks.sh 453` | `$C/M34` | 2026-10-20T21:24:45Z | 2026-10-20T21:24:46Z | 0 | Restored hash matched; named control green. |
| M35-syntax | `sh -n $C/M35/docs/tests/freeze-check/run.sh` | `$C/M35` | 2026-10-20T21:24:46Z | 2026-10-20T21:24:46Z | 0 | Valid POSIX shell syntax; no diagnostics. |
| M35-mutant | `sh runs/T-c826/checks.sh 453` | `$C/M35` | 2026-10-20T21:24:46Z | 2026-10-20T21:24:46Z | 1 | Named task assertion failed; unchanged runner: docs/tests/freeze-check/run.sh failed as recorded in mutants.md. |
| M35-restored | `sh runs/T-c826/checks.sh 453` | `$C/M35` | 2026-10-20T21:24:46Z | 2026-10-20T21:24:46Z | 0 | Restored hash matched; named control green. |
| M36-syntax | `sh -n $C/M36/docs/tests/freeze-check/run-all.sh` | `$C/M36` | 2026-10-20T21:24:46Z | 2026-10-20T21:24:46Z | 0 | Valid POSIX shell syntax; no diagnostics. |
| M36-mutant | `sh runs/T-c826/checks.sh 453` | `$C/M36` | 2026-10-20T21:24:46Z | 2026-10-20T21:24:46Z | 1 | Named task assertion failed; unchanged runner: docs/tests/freeze-check/run-all.sh failed as recorded in mutants.md. |
| M36-restored | `sh runs/T-c826/checks.sh 453` | `$C/M36` | 2026-10-20T21:24:46Z | 2026-10-20T21:24:46Z | 0 | Restored hash matched; named control green. |
| M51-mutant | `sh docs/tests/freeze-check/run-all.sh` | `$C/M51` | 2026-10-20T21:24:41Z | 2026-10-20T21:24:49Z | 1 | freeze-check: 18 passed, 63 failed; good-K5-line-limit failed as recorded in mutants.md. |
| M51-restored | `sh docs/tests/freeze-check/run.sh docs/tests/freeze-check/good-K5-line-limit` | `$C/M51` | 2026-10-20T21:24:49Z | 2026-10-20T21:24:49Z | 0 | Restored hash matched; named control green. |
| M52-syntax | `sh -n docs/tests/freeze-check.sh` | `$C/M52` | 2026-10-20T21:24:49Z | 2026-10-20T21:24:49Z | 0 | Valid POSIX shell syntax; no diagnostics. |
| M52-mutant | `sh docs/tests/freeze-check/run-all.sh` | `$C/M52` | 2026-10-20T21:24:49Z | 2026-10-20T21:24:58Z | 1 | freeze-check: 18 passed, 63 failed; good-K5-file-limit failed as recorded in mutants.md. |
| M52-restored | `sh docs/tests/freeze-check/run.sh docs/tests/freeze-check/good-K5-file-limit` | `$C/M52` | 2026-10-20T21:24:58Z | 2026-10-20T21:24:58Z | 0 | Restored hash matched; named control green. |
| M53-mutant | `sh runs/T-c826/checks.sh 453` | `$C/M53` | 2026-10-20T21:25:15Z | 2026-10-20T21:25:15Z | 1 | Named task assertion failed; all nine baseline paths failed as recorded in mutants.md. |
| M53-restored | `sh runs/T-c826/checks.sh 453` | `$C/M53` | 2026-10-20T21:25:15Z | 2026-10-20T21:25:15Z | 0 | Restored hash matched; named control green. |
| M54-mutant | `sh runs/T-c826/checks.sh 453` | `$C/M54` | 2026-10-20T21:25:15Z | 2026-10-20T21:25:15Z | 1 | Named task assertion failed; all nine baseline paths failed as recorded in mutants.md. |
| M54-restored | `sh runs/T-c826/checks.sh 453` | `$C/M54` | 2026-10-20T21:25:15Z | 2026-10-20T21:25:15Z | 0 | Restored hash matched; named control green. |
| M55-mutant | `sh runs/T-c826/checks.sh 453` | `$C/M55` | 2026-10-20T21:25:15Z | 2026-10-20T21:25:15Z | 1 | Named task assertion failed; all nine baseline paths failed as recorded in mutants.md. |
| M55-restored | `sh runs/T-c826/checks.sh 453` | `$C/M55` | 2026-10-20T21:25:15Z | 2026-10-20T21:25:15Z | 0 | Restored hash matched; named control green. |
| M56-mutant | `sh runs/T-c826/checks.sh 453` | `$C/M56` | 2026-10-20T21:25:15Z | 2026-10-20T21:25:15Z | 1 | Named task assertion failed; kept checker requires the pull_request event failed as recorded in mutants.md. |
| M56-restored | `sh runs/T-c826/checks.sh 453` | `$C/M56` | 2026-10-20T21:25:15Z | 2026-10-20T21:25:15Z | 0 | Restored hash matched; named control green. |
| M57-mutant | `sh runs/T-c826/checks.sh 453` | `$C/M57` | 2026-10-20T21:25:15Z | 2026-10-20T21:25:15Z | 1 | Named task assertion failed; one literal PR number and only the checker path failed as recorded in mutants.md. |
| M57-restored | `sh runs/T-c826/checks.sh 453` | `$C/M57` | 2026-10-20T21:25:15Z | 2026-10-20T21:25:15Z | 0 | Restored hash matched; named control green. |
| M58-mutant | `sh runs/T-c826/checks.sh 453` | `$C/M58` | 2026-10-20T21:25:15Z | 2026-10-20T21:25:15Z | 1 | Named task assertion failed; one literal PR number and only the checker path failed as recorded in mutants.md. |
| M58-restored | `sh runs/T-c826/checks.sh 453` | `$C/M58` | 2026-10-20T21:25:15Z | 2026-10-20T21:25:16Z | 0 | Restored hash matched; named control green. |
| M59-mutant | `sh runs/T-c826/checks.sh 453` | `$C/M59` | 2026-10-20T21:25:16Z | 2026-10-20T21:25:16Z | 1 | Named task assertion failed; exception matches the supplied actual migration PR failed as recorded in mutants.md. |
| M59-restored | `sh runs/T-c826/checks.sh 453` | `$C/M59` | 2026-10-20T21:25:16Z | 2026-10-20T21:25:16Z | 0 | Restored hash matched; named control green. |
| M60-mutant | `sh runs/T-c826/checks.sh 453` | `$C/M60` | 2026-10-20T21:25:16Z | 2026-10-20T21:25:16Z | 1 | Named task assertion failed; event guard receives the actual GitHub event name failed as recorded in mutants.md. |
| M60-restored | `sh runs/T-c826/checks.sh 453` | `$C/M60` | 2026-10-20T21:25:16Z | 2026-10-20T21:25:16Z | 0 | Restored hash matched; named control green. |
| M61-mutant | `sh runs/T-c826/checks.sh 453` | `$C/M61` | 2026-10-20T21:25:16Z | 2026-10-20T21:25:16Z | 1 | Named task assertion failed; number guard receives the actual pull request number failed as recorded in mutants.md. |
| M61-restored | `sh runs/T-c826/checks.sh 453` | `$C/M61` | 2026-10-20T21:25:16Z | 2026-10-20T21:25:16Z | 0 | Restored hash matched; named control green. |
| M62-mutant | `sh runs/T-c826/checks.sh 453` | `$C/M62` | 2026-10-20T21:25:16Z | 2026-10-20T21:25:16Z | 1 | Named task assertion failed; only the guarded checker skips restoration failed as recorded in mutants.md. |
| M62-restored | `sh runs/T-c826/checks.sh 453` | `$C/M62` | 2026-10-20T21:25:16Z | 2026-10-20T21:25:16Z | 0 | Restored hash matched; named control green. |
| M63-mutant | `sh runs/T-c826/checks.sh 453` | `$C/M63` | 2026-10-20T21:25:16Z | 2026-10-20T21:25:16Z | 1 | Named task assertion failed; only the guarded checker skips restoration failed as recorded in mutants.md. |
| M63-restored | `sh runs/T-c826/checks.sh 453` | `$C/M63` | 2026-10-20T21:25:16Z | 2026-10-20T21:25:16Z | 0 | Restored hash matched; named control green. |
| M64-mutant | `sh runs/T-c826/checks.sh 453` | `$C/M64` | 2026-10-20T21:25:16Z | 2026-10-20T21:25:16Z | 1 | Named task assertion failed; restoration still reads the default branch of this repository failed as recorded in mutants.md. |
| M64-restored | `sh runs/T-c826/checks.sh 453` | `$C/M64` | 2026-10-20T21:25:16Z | 2026-10-20T21:25:16Z | 0 | Restored hash matched; named control green. |
| M65-mutant | `sh runs/T-c826/checks.sh 453` | `$C/M65` | 2026-10-20T21:25:16Z | 2026-10-20T21:25:16Z | 1 | Named task assertion failed; default branch still comes from the repository event failed as recorded in mutants.md. |
| M65-restored | `sh runs/T-c826/checks.sh 453` | `$C/M65` | 2026-10-20T21:25:16Z | 2026-10-20T21:25:16Z | 0 | Restored hash matched; named control green. |
| M66-mutant | `sh runs/T-c826/checks.sh 453` | `$C/M66` | 2026-10-20T21:25:16Z | 2026-10-20T21:25:16Z | 1 | Named task assertion failed; exception comment names the task, literal PR and removal issue failed as recorded in mutants.md. |
| M66-restored | `sh runs/T-c826/checks.sh 453` | `$C/M66` | 2026-10-20T21:25:16Z | 2026-10-20T21:25:16Z | 0 | Restored hash matched; named control green. |
| final-task-check | `sh runs/T-c826/checks.sh 453` | `$W` | 2026-10-20T21:27:05Z | 2026-10-20T21:27:05Z | 0 | All 13 task assertions passed with actual PR 453. |
| final-suite | `sh docs/tests/freeze-check/run-all.sh` | `$W` | 2026-10-20T21:27:05Z | 2026-10-20T21:27:14Z | 0 | freeze-check: 81 passed, 0 failed |
| final-whitespace | `git diff --check` | `$W` | 2026-10-20T21:27:36Z | 2026-10-20T21:27:36Z | 0 | No diagnostics. |
| final-links | `sh docs/links/link-lint.sh` | `$W` | 2026-10-20T21:27:35Z | 2026-10-20T21:27:48Z | 0 | link-lint: OK  2995 links resolved |
| final-audit | `python3 $C/final-audit.py` | `$W` | 2026-10-20T21:28:01Z | 2026-10-20T21:28:01Z | 0 | Working diff against d27b331, excluding only the five evidence records: 844 lines / 143 files. This is a scope check, not a committed-head freeze or landing audit. |

The final audit confirms that only the ten planned files changed and that the two
runners, task check, workflow, setup-check and review-record parser retain their
committed hashes. Against d27b331 the working diff, excluding only the five
evidence records, is **844 lines / 143 files**, below the reviewed 1300/170 bound.
No task file or completed-log line is written. This is a working-tree scope check;
clean committed-head freeze, hosted green and independent review remain later.

The complete baseline red mismatches are recorded in Stage A; hosted red is above.
The mutation ledger quotes each named failure, its checker/harness statuses and
restoration hashes; other failing cases are intentionally omitted there. Real-head
replays and their consumer inspections are in real-heads.md. These records are
observational evidence and are not fixture/configuration inputs.

## Stage D — gates on the committed head and pre-freeze projection

Tested clean committed HEAD: `ecf98718028c99e3997d111f682f1ed3b87d917e`
(Stage C commits `7cdf80f`, `34ff09b`, `ecf9871`). Named base `origin/main`
resolved to `d27b33180a1aa82c23f3ab66e757e06de356eb35`; it was not fetched or
changed during this stage. The worktree remained clean through both checker
invocations. Only this evidence record was appended afterward. No task detail,
completed-log line, checker, runner, workflow, task check or conditions record
was changed.

Each local command below has its UTC endpoints read directly with
`date -u +%FT%TZ` and its direct subprocess exit status. Output was captured
without a pipeline. All temporary repositories and scratch output used
`TMPDIR=/home/layup/.cache/layup-session/tmp-green/T-c826-stage-d`.
The two explicitly redirected checker outputs were captured at the named paths
outside the worktree. The gate driver printed each step and an eight-second
heartbeat while a command was still running.

### Hosted green after implementation

[PR 453, run 39500930110, job 118897922221](https://github.com/pharzam/layup/actions/runs/39500930110/job/118897922221)
is the hosted green paired with Stage B's red run 39497659349 above.
The connected GitHub tools independently returned all nine jobs of run
39500930110 as completed/success, including Go tests, lint, security and setup.
The nested-checkout job's restore, nested-checkout and freeze-suite steps all
concluded success. The supplied stage input separately reports that
`review-record` waits for the independent round; this is not a claim of an
all-checks green or a completed review.

The hosted checkout is synthetic merge
`a037930c65445302e1faf88ca7c0fe32619e7afd`, merging branch head
`ecf98718028c99e3997d111f682f1ed3b87d917e` into base
`d27b33180a1aa82c23f3ab66e757e06de356eb35`, as the decoded checkout log states.
The decoded job log starts at `2026-10-20T21:30:38.8247449Z`; its last line is
`2026-10-20T21:30:53.6102549Z`. The hosted nested-checkout command starts at
`2026-10-20T21:30:42.4853427Z` and reports 10 behaved cases at
`2026-10-20T21:30:42.8176936Z`; the hosted freeze suite starts at
`2026-10-20T21:30:42.8312654Z` and reports 81 passed, 0 failed at
`2026-10-20T21:30:53.3403200Z`. Both steps succeeded (exit 0 inferred from
GitHub's successful step conclusions; the log does not print numeric exits).

The read-only `github_fetch_workflow_run_jobs` and
`github_fetch_workflow_job_logs` calls were made between direct local UTC
reads `2026-10-20T21:42:35Z` and `2026-10-20T21:43:02Z`.
Both returned `isError: false`; these tool calls have no shell exit status.
Selected decoded log lines preserve the actual keep and every restoration:

```text
2026-10-20T21:30:41.1717356Z HEAD is now at a037930 Merge ecf98718028c99e3997d111f682f1ed3b87d917e into d27b33180a1aa82c23f3ab66e757e06de356eb35
2026-10-20T21:30:42.4054744Z restored from main: docs/adr/adr-lint.sh
2026-10-20T21:30:42.4072385Z restored from main: docs/prd/prd-lint.sh
2026-10-20T21:30:42.4088726Z restored from main: docs/links/link-lint.sh
2026-10-20T21:30:42.4105392Z restored from main: docs/tests/run-discipline-tests.sh
2026-10-20T21:30:42.4121226Z restored from main: docs/ci/pr-link-lint.sh
2026-10-20T21:30:42.4136765Z restored from main: docs/ci/review-record-lint.sh
2026-10-20T21:30:42.4137556Z kept as written for PR 453: docs/tests/freeze-check.sh
2026-10-20T21:30:42.4153306Z restored from main: docs/tests/freeze-check/run.sh
2026-10-20T21:30:42.4168924Z restored from main: docs/tests/freeze-check/run-all.sh
2026-10-20T21:30:42.4853427Z ##[group]Run sh docs/tests/nested-checkout-check.sh
2026-10-20T21:30:42.8176936Z nested-checkout-check: OK  10 cases behaved
2026-10-20T21:30:42.8312654Z ##[group]Run sh docs/tests/freeze-check/run-all.sh
2026-10-20T21:30:53.3403200Z freeze-check: 81 passed, 0 failed
```

### Local commands on ecf9871

| Command | UTC start | UTC end | Exit | Observation |
| --- | --- | --- | --- | --- |
| `sh docs/adr/adr-lint.sh` | 2026-10-20T21:42:31Z | 2026-10-20T21:42:32Z | 0 | adr-lint: OK |
| `sh docs/prd/prd-lint.sh` | 2026-10-20T21:42:32Z | 2026-10-20T21:42:32Z | 0 | prd-lint: OK |
| `sh docs/links/link-lint.sh` | 2026-10-20T21:42:32Z | 2026-10-20T21:42:44Z | 0 | link-lint: OK  2995 links resolved |
| `sh docs/tests/run-discipline-tests.sh` | 2026-10-20T21:42:44Z | 2026-10-20T21:42:47Z | 0 | run-discipline-tests: 138 passed, 0 failed |
| `sh docs/tests/nested-checkout-check.sh` | 2026-10-20T21:42:47Z | 2026-10-20T21:42:47Z | 0 | nested-checkout-check: OK  10 cases behaved |
| `sh docs/tests/freeze-check/run-all.sh` | 2026-10-20T21:42:47Z | 2026-10-20T21:42:55Z | 0 | freeze-check: 81 passed, 0 failed |
| `bash --posix docs/tests/freeze-check/run-all.sh` | 2026-10-20T21:42:55Z | 2026-10-20T21:43:03Z | 0 | freeze-check: 81 passed, 0 failed |
| `sh docs/setup/setup-check.sh` | 2026-10-20T21:43:03Z | 2026-10-20T21:43:25Z | 0 | All setup checks OK, including ci; 138 discipline cases pass and all kit linters pass. |
| `sh runs/T-c826/checks.sh` | 2026-10-20T21:43:25Z | 2026-10-20T21:43:25Z | 0 | All 12 non-number-specific assertions pass; prints pending for the optional actual PR-number argument. |
| `git diff --check origin/main HEAD` | 2026-10-20T21:43:25Z | 2026-10-20T21:43:25Z | 0 | No diagnostics. |
| `sh docs/tests/freeze-check.sh --task T-c826 --base origin/main --budget 1300/170 > /home/layup/.cache/layup-session/tmp-green/c826-fc.txt 2>&1` | 2026-10-20T21:43:25Z | 2026-10-20T21:43:25Z | 1 | K1–K4 pass; K5 fails on the gross figure; K6 fails at conditions.md:9. Full output below. |
| `sh docs/tests/freeze-check.sh --task T-c826 --base origin/main --budget 1300/170 --evidence runs/T-c826/test-runs.md --evidence runs/T-c826/mutants.md --evidence runs/T-c826/real-heads.md --evidence runs/T-c826/conditions.md --evidence runs/T-c826/ci-exception.md > /home/layup/.cache/layup-session/tmp-green/c826-fc-net.txt 2>&1` | 2026-10-20T21:43:25Z | 2026-10-20T21:43:26Z | 1 | K1–K5 pass; K6 fails at conditions.md:9. Full output below. |
| `sh runs/T-c826/checks.sh 453` | 2026-10-20T21:43:26Z | 2026-10-20T21:43:26Z | 0 | All 13 assertions pass, including actual migration PR 453. |

The no-argument task-check invocation is the exact Stage D request. Its pending
message concerns the optional number-specific assertion, so the additional
invocation with `453` verifies that assertion explicitly.

### Requested checker invocation, without declarations

Complete output from
`/home/layup/.cache/layup-session/tmp-green/c826-fc.txt`, exit **1**:

```text
ok K1: white space
ok K2: registration of new tests
ok K3: test selectors
ok K4: restore evidence of mutants
FAIL K5: -:0: projection: figure 4110/150, maximum 1300/170, tolerance 260/4, limit 1560/174, margin -2550/24 (the figure at landing is still measured at landing)
FAIL K6: runs/T-c826/conditions.md:9: Evidence cell has no code span
```

### ADR-0028 classification and declared-evidence invocation

Inspected the five records and their repository consumers before the explicit
`--evidence` invocation. `test-runs.md` records commands, times and diagnostics;
`mutants.md` records mutation observations and restoration hashes;
`real-heads.md` records historical replay results and consumer inspections;
`conditions.md` maps review conditions to evidence; and `ci-exception.md` records
authorization and CI observations. They are not execution data, configuration,
fixtures or expected outputs. The task check reads operative documents and
committed runner objects, not these records as execution inputs. The checker
reads task Markdown to audit selectors and mutation evidence (K3/K4), and the
conditions table to audit citations (K6). These evidence audits do not make the
records execution inputs under the canonical definition. The declaration is the
author's classification of inspected consumers, not proof about unknown consumers.

| Excluded evidence record | Added plus removed lines at ecf9871 |
| --- | --- |
| `runs/T-c826/test-runs.md` | 1739 |
| `runs/T-c826/mutants.md` | 1377 |
| `runs/T-c826/real-heads.md` | 88 |
| `runs/T-c826/conditions.md` | 11 |
| `runs/T-c826/ci-exception.md` | 37 |
| Total excluded | 3252 lines / 5 files |

The executable `runs/T-c826/checks.sh` remains counted: **84 lines / 1 file**.
No execution input or close-out bookkeeping was excluded. K5 deliberately uses
`git diff --no-renames --numstat origin/main HEAD`; its gross figure is
**4110 lines / 150 files**, and its net figure is **858 lines / 145 files**.
The earlier Stage C scope estimate of 844/143 used rename folding. The difference
is the unchanged renamed ARGS file (1 line at each endpoint) and fixture body
(6 lines at each endpoint), which K5 counts as delete/add: another 14 lines and
2 touched files. The figures here follow the branch's own checker consistently.

| Measure | Lines | Files |
| --- | --- | --- |
| Figure, evidence out and execution inputs in | 858 | 145 |
| Reviewed maximum | 1300 | 170 |
| Tolerance | 260 | 4 |
| Limit, maximum plus tolerance | 1560 | 174 |
| Margin, limit minus figure | 702 | 29 |

The figure is 442 lines and 25 files below the reviewed maximum. The append to
this evidence record also leaves the count under ADR-0028. This remains a
projection against the named base, not the final landing audit.

Complete declared-evidence output from
`/home/layup/.cache/layup-session/tmp-green/c826-fc-net.txt`, exit **1**:

```text
ok K1: white space
ok K2: registration of new tests
ok K3: test selectors
ok K4: restore evidence of mutants
ok K5: projection: figure 858/145, maximum 1300/170, tolerance 260/4, limit 1560/174, margin 702/29 (the figure at landing is still measured at landing)
FAIL K6: runs/T-c826/conditions.md:9: Evidence cell has no code span
```

### Stage outcome

All ten requested gates exit 0; both 81-case suites and the actual-number task
check are green. Both committed-head pre-freeze invocations exit 1. Declaring
eligible evidence resolves K5's gross overrun, but **K6 still fails** because the
Evidence cell on committed `runs/T-c826/conditions.md:9` has no code span.
The checker is therefore not green and this head has not passed pre-freeze.
The stage author leaves the defect unchanged because Stage D expressly permits
only this evidence-record append. A later authorized stage must correct the
condition citation, commit it and rerun pre-freeze on the clean committed head.
No independent review or close-out is claimed here.

Final evidence-append verification:

| Command | UTC start | UTC end | Exit | Observation |
| --- | --- | --- | --- | --- |
| `git diff --check` | 2026-10-20T21:44:57Z | 2026-10-20T21:44:57Z | 0 | No diagnostics after appending the Stage D evidence. |

Read-only scope checks also confirm that the committed record is an unchanged
prefix, this record is the only modified file, and no untracked file exists.

## Stage F — clean-head pre-freeze after Stage E

All four commands below ran on clean committed head
`b11d8d494b780b33233850ae3f091df37eff37fa`, containing Stage E fixes
`569229b` and `b11d8d4`, against `origin/main` at
`d27b33180a1aa82c23f3ab66e757e06de356eb35`.
The worktree remained clean through all four runs; this append followed them.
Times come directly from `date -u +%FT%TZ`; exits are the observed command
statuses. `TMPDIR` was set to
`/home/layup/.cache/layup-session/tmp-green`, which also holds captured output.

| Command | UTC start | UTC end | Exit | Observation |
| --- | --- | --- | --- | --- |
| `sh docs/tests/freeze-check.sh --task T-c826 --base origin/main --budget 1300/170 --evidence runs/T-c826/ci-exception.md --evidence runs/T-c826/conditions.md --evidence runs/T-c826/mutants.md --evidence runs/T-c826/real-heads.md --evidence runs/T-c826/test-runs.md > /home/layup/.cache/layup-session/tmp-green/c826-fc3.txt 2>&1` | 2026-10-20T21:50:16Z | 2026-10-20T21:50:17Z | 0 | K1–K6 pass; complete output below. |
| `sh docs/tests/freeze-check/run-all.sh` | 2026-10-20T21:50:17Z | 2026-10-20T21:50:25Z | 0 | freeze-check: 81 passed, 0 failed. |
| `sh docs/links/link-lint.sh` | 2026-10-20T21:50:25Z | 2026-10-20T21:50:37Z | 0 | link-lint: OK  2997 links resolved. |
| `sh docs/setup/setup-check.sh` | 2026-10-20T21:50:37Z | 2026-10-20T21:51:02Z | 0 | Every setup check passes, including ci; 138 discipline cases pass and all kit linters pass. |

Complete freeze-check output, captured before any worktree edit in
`/home/layup/.cache/layup-session/tmp-green/c826-fc3.txt`, exit **0**:

```text
ok K1: white space
ok K2: registration of new tests
ok K3: test selectors
ok K4: restore evidence of mutants
ok K5: projection: figure 871/146, maximum 1300/170, tolerance 260/4, limit 1560/174, margin 689/28 (the figure at landing is still measured at landing)
ok K6: conditions table
```

The declared-evidence projection is **871 lines / 146 files**, below the
reviewed maximum of **1300 / 170**. Tolerance is **260 / 4**, limit is
**1560 / 174**, and margin is **689 / 28**. K6 now accepts the condition
citations. This is the clean-head pre-freeze result; it does not claim an
independent review or close-out.

## Stage G — round 1 fixes (cycle 1 of cap 2)

Round 1 is the supplied Claude Opus 5.5 review, comment
[6303062417](https://github.com/pharzam/layup/issues/449#issuecomment-6303062417),
material on F1 and F2 at `fdb8d626008762b2e13cbab5fe53c4d5d3226700`.
This stage started on that head with a clean worktree. It writes the fixes and
records for the coordinator's ordered test/fix/docs commits. No commit,
pre-freeze invocation, task detail or completed-log entry is made in this stage.
Both runners and `runs/T-c826/checks.sh` remain unchanged.
All scratch files, captured outputs and temporary fixture repositories are under
`/home/layup/.cache/layup-session/tmp-green`; `TMPDIR` selects that directory.
Times below are direct reads of `date -u +%FT%TZ`; exits are the direct child
statuses, not the exit of a pipeline or of the recording helper.

### Cases before code

- **F1, HEAD tree:** `bad-usage-evidence-directory-head` declares a `.md`
  directory containing a 61-line unlisted script. The current checker accepts
  it and excludes the script, printing figure 0/0 and exiting 0. The fixed
  expectation is exit 2 before any K1–K6 output.
- **F1, base tree:** `bad-usage-evidence-directory-base` replaces a directory
  at base with a one-line record at HEAD, declaring that path. The base tree's
  unlisted script must not be subtracted. It currently also prints figure 0/0
  and exits 0, against the expected exit 2.
- **F2, rename into evidence:** `good-K5-evidence-renamed` deletes the baseline's
  three-line `test-runs.md` and adds the identical `moved.md`, declaring only
  the destination. The original checker already passes: figure 3/1, margin
  57/3 at maximum 0/0. This is a green regression control, not claimed as red.
  M69 removes only `--no-renames` on the gross diff. The case then fails on
  figure -3/0 and margin 63/4, with checker exit 0 and harness exit 2. Restoring
  the flag restores the named green control before any Stage G checker edit.
- **Note 3, diagnostic:** `bad-usage-evidence-trailing-slash` expects an empty
  path-component diagnostic instead of calling the trailing slash traversal.
  Its observed pre-fix mismatch has checker exit 2 and harness exit 2.

The added base-tree input uses the existing shared `base/` fixture, `.in`
overlays and `ABSENT.in`, without runner metadata. It changes the deterministic
base commit from `8f46ad9` to `1d1b4a0512eff698c4107a45939f9f8c38f060aa`.
The first suite therefore also failed `good-K6-valid`, which cited the old
base commit. That unrelated failure is not credited as red evidence for F1.
`G-base-sha` reconstructs the shared base using the runner's overlay, identity
and fixed dates. Updating only that case's commit citation to `1d1b4a0` restores
its original valid-commit assertion; its EXPECT is unchanged. The corrected
pre-code suite has exactly the two intended F1 failures: **82 passed, 2 failed**.

Observed commands through the first green suite and mutation controls:

| Run | Command | UTC start | UTC end | Exit |
| --- | --- | --- | --- | --- |
| G-red | `sh docs/tests/freeze-check/run-all.sh` | 2026-10-20T22:21:22Z | 2026-10-20T22:21:31Z | 1 |
| G-base-sha | `python3 /home/layup/.cache/layup-session/tmp-green/T-c826-stage-g/base-sha.py` | 2026-10-20T22:21:58Z | 2026-10-20T22:21:58Z | 0 |
| G-red-corrected | `sh docs/tests/freeze-check/run-all.sh` | 2026-10-20T22:22:03Z | 2026-10-20T22:22:11Z | 1 |
| M69-syntax | `sh -n /home/layup/.cache/layup-session/tmp-green/T-c826-stage-g/M69/docs/tests/freeze-check.sh` | 2026-10-20T22:22:22Z | 2026-10-20T22:22:22Z | 0 |
| M69-mutant | `sh /home/layup/.cache/layup-session/tmp-green/T-c826-stage-g/M69/docs/tests/freeze-check/run.sh /home/layup/.cache/layup-session/tmp-green/T-c826-stage-g/M69/docs/tests/freeze-check/good-K5-evidence-renamed` | 2026-10-20T22:22:22Z | 2026-10-20T22:22:23Z | 2 |
| M69-restored | `sh /home/layup/.cache/layup-session/tmp-green/T-c826-stage-g/M69/docs/tests/freeze-check/run.sh /home/layup/.cache/layup-session/tmp-green/T-c826-stage-g/M69/docs/tests/freeze-check/good-K5-evidence-renamed` | 2026-10-20T22:22:23Z | 2026-10-20T22:22:23Z | 0 |
| G-red-trailing-slash | `sh docs/tests/freeze-check/run.sh docs/tests/freeze-check/bad-usage-evidence-trailing-slash` | 2026-10-20T22:22:46Z | 2026-10-20T22:22:46Z | 2 |
| G-green | `sh docs/tests/freeze-check/run-all.sh` | 2026-10-20T22:22:55Z | 2026-10-20T22:23:04Z | 0 |
| M67-syntax | `sh -n /home/layup/.cache/layup-session/tmp-green/T-c826-stage-g/M67/docs/tests/freeze-check.sh` | 2026-10-20T22:23:19Z | 2026-10-20T22:23:19Z | 0 |
| M67-mutant | `sh /home/layup/.cache/layup-session/tmp-green/T-c826-stage-g/M67/docs/tests/freeze-check/run.sh /home/layup/.cache/layup-session/tmp-green/T-c826-stage-g/M67/docs/tests/freeze-check/bad-usage-evidence-directory-head` | 2026-10-20T22:23:19Z | 2026-10-20T22:23:20Z | 2 |
| M67-restored | `sh /home/layup/.cache/layup-session/tmp-green/T-c826-stage-g/M67/docs/tests/freeze-check/run.sh /home/layup/.cache/layup-session/tmp-green/T-c826-stage-g/M67/docs/tests/freeze-check/bad-usage-evidence-directory-head` | 2026-10-20T22:23:20Z | 2026-10-20T22:23:20Z | 0 |
| M68-syntax | `sh -n /home/layup/.cache/layup-session/tmp-green/T-c826-stage-g/M68/docs/tests/freeze-check.sh` | 2026-10-20T22:23:20Z | 2026-10-20T22:23:20Z | 0 |
| M68-mutant | `sh /home/layup/.cache/layup-session/tmp-green/T-c826-stage-g/M68/docs/tests/freeze-check/run.sh /home/layup/.cache/layup-session/tmp-green/T-c826-stage-g/M68/docs/tests/freeze-check/bad-usage-evidence-directory-base` | 2026-10-20T22:23:20Z | 2026-10-20T22:23:20Z | 2 |
| M68-restored | `sh /home/layup/.cache/layup-session/tmp-green/T-c826-stage-g/M68/docs/tests/freeze-check/run.sh /home/layup/.cache/layup-session/tmp-green/T-c826-stage-g/M68/docs/tests/freeze-check/bad-usage-evidence-directory-base` | 2026-10-20T22:23:20Z | 2026-10-20T22:23:20Z | 0 |
| M70-syntax | `sh -n /home/layup/.cache/layup-session/tmp-green/T-c826-stage-g/M70/docs/tests/freeze-check.sh` | 2026-10-20T22:23:20Z | 2026-10-20T22:23:20Z | 0 |
| M70-mutant | `sh /home/layup/.cache/layup-session/tmp-green/T-c826-stage-g/M70/docs/tests/freeze-check/run.sh /home/layup/.cache/layup-session/tmp-green/T-c826-stage-g/M70/docs/tests/freeze-check/bad-usage-evidence-trailing-slash` | 2026-10-20T22:23:20Z | 2026-10-20T22:23:20Z | 2 |
| M70-restored | `sh /home/layup/.cache/layup-session/tmp-green/T-c826-stage-g/M70/docs/tests/freeze-check/run.sh /home/layup/.cache/layup-session/tmp-green/T-c826-stage-g/M70/docs/tests/freeze-check/bad-usage-evidence-trailing-slash` | 2026-10-20T22:23:20Z | 2026-10-20T22:23:20Z | 0 |

Complete corrected pre-code suite output:

```text
ok    bad-K1-blank-eof
ok    bad-K1-trailing-space
ok    bad-K2-missing-file
ok    bad-K2-missing-prd
ok    bad-K2-missing-traceability
ok    bad-K2-no-prd-frame-test
ok    bad-K2-no-prd-mapped-test
ok    bad-K2-no-s12-frame-test
ok    bad-K2-no-s12-mapped-test
ok    bad-K2-prd-other-file
ok    bad-K2-two-missing
ok    bad-K3-backtick-outside
ok    bad-K3-double-quoted
ok    bad-K3-no-match
ok    bad-K3-other-package
ok    bad-K3-other-tag
ok    bad-K3-span-across-lines
ok    bad-K3-subtest
ok    bad-K3-unclosed-span-eof
ok    bad-K3-unclosed-span-paragraph
ok    bad-K3-unquoted
ok    bad-K3-unsupported-form
ok    bad-K4-dropped-in-description
ok    bad-K4-form-a-different
ok    bad-K4-form-b-equal-hashes
ok    bad-K4-result-only
ok    bad-K4-two-hashes-no-restore
ok    bad-K5-over-files
ok    bad-K5-over-lines
ok    bad-K6-empty-table
ok    bad-K6-line-past-end
ok    bad-K6-line-zero
ok    bad-K6-missing-commit
ok    bad-K6-missing-file
ok    bad-K6-missing-path
ok    bad-K6-missing-test
ok    bad-K6-mixed-valid-invalid
ok    bad-K6-no-span
ok    bad-K6-range-zero
ok    bad-base-not-a-commit
ok    bad-crlf-conditions
ok    bad-env-K2-unreadable-prd
ok    bad-env-K2-unreadable
ok    bad-env-K3-unreadable-test
ok    bad-env-dirty-worktree
ok    bad-env-untracked-evidence
ok    bad-not-a-worktree
ok    bad-usage-bad-budget
ok    bad-usage-closeout
FAIL  bad-usage-evidence-directory-base
      run: bad-usage-evidence-directory-base: exit 0 (wanted 2); diff (-EXPECT +got):
      --- /home/layup/.cache/layup-session/tmp-green/tmp.XxVftHtVeK/expect	2026-10-21 01:22:08.172232035 +0300
      +++ /home/layup/.cache/layup-session/tmp-green/tmp.XxVftHtVeK/got	2026-10-21 01:22:08.170864025 +0300
      @@ -1,2 +1,7 @@
      -2
      -freeze-check: --evidence runs/T-0000/base-tree.md: must name a file at base
      +0
      +ok K1: white space
      +ok K2: registration of new tests
      +ok K3: test selectors
      +ok K4: restore evidence of mutants
      +ok K5: projection: figure 0/0, maximum 0/0, tolerance 60/4, limit 60/4, margin 60/4 (the figure at landing is still measured at landing)
      +ok K6: conditions table
FAIL  bad-usage-evidence-directory-head
      run: bad-usage-evidence-directory-head: exit 0 (wanted 2); diff (-EXPECT +got):
      --- /home/layup/.cache/layup-session/tmp-green/tmp.8vberUnnMc/expect	2026-10-21 01:22:08.293864358 +0300
      +++ /home/layup/.cache/layup-session/tmp-green/tmp.8vberUnnMc/got	2026-10-21 01:22:08.293780883 +0300
      @@ -1,2 +1,7 @@
      -2
      -freeze-check: --evidence runs/T-0000/d.md: must name a file at HEAD
      +0
      +ok K1: white space
      +ok K2: registration of new tests
      +ok K3: test selectors
      +ok K4: restore evidence of mutants
      +ok K5: projection: figure 0/0, maximum 0/0, tolerance 60/4, limit 60/4, margin 60/4 (the figure at landing is still measured at landing)
      +ok K6: conditions table
ok    bad-usage-evidence-missing
ok    bad-usage-evidence-no-value
ok    bad-usage-evidence-other-task
ok    bad-usage-evidence-suffix
ok    bad-usage-evidence-task-prefix
ok    bad-usage-evidence-traversal
ok    bad-usage-missing-budget
ok    bad-usage-no-args
ok    bad-usage-unknown-flag
ok    good-K1-clean
ok    good-K2-prd-by-file
ok    good-K2-registered
ok    good-K3-alternatives
ok    good-K3-span-across-lines
ok    good-K4-forms
ok    good-K5-at-max
ok    good-K5-ceiling
ok    good-K5-evidence-added
ok    good-K5-evidence-deleted
ok    good-K5-evidence-duplicate
ok    good-K5-evidence-inputs
ok    good-K5-evidence-kinds
ok    good-K5-evidence-literal
ok    good-K5-evidence-modified
ok    good-K5-evidence-renamed
ok    good-K5-file-limit
ok    good-K5-line-limit
ok    good-K5-no-closeout
ok    good-K5-percentage
ok    good-K5-truncation
ok    good-K6-no-conditions
ok    good-K6-valid
ok    good-crlf-all
freeze-check: 82 passed, 2 failed
```

Complete pre-code diagnostic mismatch:

```text
run: bad-usage-evidence-trailing-slash: exit 2 (wanted 2); diff (-EXPECT +got):
--- /home/layup/.cache/layup-session/tmp-green/tmp.501FIUMVb6/expect	2026-10-21 01:22:46.279283441 +0300
+++ /home/layup/.cache/layup-session/tmp-green/tmp.501FIUMVb6/got	2026-10-21 01:22:46.278603956 +0300
@@ -1,2 +1,2 @@
 2
-freeze-check: --evidence runs/T-0000/: empty path components are not allowed
+freeze-check: --evidence runs/T-0000/: path traversal is not allowed
```

### Fix and mutation outcome

The validator now requires `git cat-file -t` to return `blob` at each committed
endpoint where the declaration exists; a tree is refused with exit 2 before
results print. Existing absent-at-both-endpoints and path/suffix checks remain.
This adds an object-kind check, with no executable-mode or symlink restriction.
Both F1 cases and every prior case pass: **85 passed, 0 failed**.
M67 removes HEAD validation alone; M68 removes base validation alone. Each
reproduces its named false acceptance and returns harness exit 2. M70 restores
only the incorrect diagnostic and fails its exact-output assertion. All syntax
checks and all four restored controls pass. Fresh source/mutant/restored
SHA-256 values and the exact mismatches are in [mutants.md](mutants.md#stage-g--round-1-fixes-m67m70).

### Round 1 notes

| Note | Disposition and reason |
| --- | --- |
| 1. One home per rule | Applied. The test README links the canonical input definition and budget rule instead of restating the input carve-out and tolerance constants/inclusive limit. It retains K5's whole-line truncation, literal selection, duplicate handling, committed endpoints and output behavior. The CI README links the guardrail's migration lesson. The test README also adds the F1 file-kind and F2 delete-plus-add rename clauses; the existing freeze-check pitfall records both lessons. |
| 2. Leading zeros in budget | Accepted as a revealed, pre-existing defect outside this task's DoD, as the reviewer agrees. The parser is unchanged here; normal follow-up must refuse leading zeros. The supplied review gives no follow-up issue number, so this record does not claim an issue was opened. |
| 3. Empty/multiline guard and trailing slash | The diagnostic fix is applied with a pre-code mismatch, green case and M70 kill. Decline runner changes to encode empty/multiline ARGS: amendment 1 and plan-review condition 2 keep both runners unchanged; the guard is retained as defense in depth. No figure or validator coverage is claimed for removing that guard. |
| 4. PR identity | Applied by read-only `github_get_pr_info` for `pharzam/layup`, PR 453. The response names head T-c826 at fdb8d626008762b2e13cbab5fe53c4d5d3226700, base main at d27b33180a1aa82c23f3ab66e757e06de356eb35, and the body contains its own line `Closes #449`. This resolves the earlier inability to verify the body. Final hosted pr-link evidence belongs to the coordinator's later CI/close-out stage; no final-CI pass is claimed here. |
| 5. T-szvp removal list | The checklist below is preserved for the coordinator. Decline posting an external issue comment in this repository-change stage: it returns a local commit plan, and no direct instruction authorizes that issue-thread message. Read-only `github_fetch_issue` confirms #452 is open and its body does not yet include the checklist. |

Removal checklist for [T-szvp (#452)](https://github.com/pharzam/layup/issues/452):

- `.github/workflows/ci.yml`: the R4/removal and approval comments, event/PR
  environment values, checker-only keep/continue branch, and suite comment.
- `docs/ci/README.md`: the job row and temporary exception explanation.
- `docs/guardrails.md`: the Inv-3 exception and restore-pitfall exception text.
- `docs/tests/README.md`: the T-c826 migration-check invocation and exception text.
- `runs/T-c826/checks.sh`: the exception-specific assertions.

The forge observations above are read-only connector responses in this stage,
not shell test commands or evidence of a later hosted run. No issue was changed.

### Stage G validation

| Run | Command | UTC start | UTC end | Exit | Observation |
| --- | --- | --- | --- | --- | --- |
| G-syntax | `sh -n docs/tests/freeze-check.sh` | 2026-10-20T22:25:25Z | 2026-10-20T22:25:25Z | 0 | POSIX shell syntax accepted. |
| G-task-check | `sh runs/T-c826/checks.sh 453` | 2026-10-20T22:25:25Z | 2026-10-20T22:25:25Z | 0 | All 13 assertions pass with actual PR number 453. |
| G-adr | `sh docs/adr/adr-lint.sh` | 2026-10-20T22:25:25Z | 2026-10-20T22:25:26Z | 0 | adr-lint: OK. |
| G-prd | `sh docs/prd/prd-lint.sh` | 2026-10-20T22:25:26Z | 2026-10-20T22:25:26Z | 0 | prd-lint: OK. |
| G-links | `sh docs/links/link-lint.sh` | 2026-10-20T22:25:26Z | 2026-10-20T22:25:40Z | 0 | 3001 links resolved. |
| G-discipline | `sh docs/tests/run-discipline-tests.sh` | 2026-10-20T22:25:40Z | 2026-10-20T22:25:43Z | 0 | 138 passed, 0 failed. |
| G-nested | `sh docs/tests/nested-checkout-check.sh` | 2026-10-20T22:25:43Z | 2026-10-20T22:25:44Z | 0 | All 10 cases behaved. |
| G-setup-fixtures | `sh docs/setup/tests/run.sh` | 2026-10-20T22:25:44Z | 2026-10-20T22:25:47Z | 0 | 48 passed, 0 failed. |
| G-setup | `sh docs/setup/setup-check.sh` | 2026-10-20T22:25:47Z | 2026-10-20T22:26:10Z | 0 | Every setup check passes, including ci, discipline cases and kit linters. |
| G-diff | `git diff --check` | 2026-10-20T22:26:10Z | 2026-10-20T22:26:10Z | 2 | One extra blank line at EOF in this record; removed before the final recheck. |

The initial whitespace failure and verification-driver assertion are failures, not green evidence. Only this record's final empty line required correction. The substantive checker, fixture and documentation checks above passed; the final whitespace and scope checks follow. No pre-freeze run is included.

| Run | Command | UTC start | UTC end | Exit | Observation |
| --- | --- | --- | --- | --- | --- |
| G-final-diff | `git diff --check` | 2026-10-20T22:26:22Z | 2026-10-20T22:26:22Z | 0 | No whitespace diagnostics. |
| G-scope | `python3 /home/layup/.cache/layup-session/tmp-green/T-c826-stage-g/final-scope.py` | 2026-10-20T22:26:53Z | 2026-10-20T22:26:53Z | 0 | All 21 changed paths match the three-commit plan; unchanged HEAD/index, runners, checks.sh, workflow and close-out files. All planned files have one final newline and no trailing whitespace. |

Stage G outcome: both material findings have their cases and fixes, all 85 fixtures pass, four named mutants fail for the intended reasons and their restored controls pass, and all relevant lints/setup/task checks pass. The whitespace failure is corrected. The coordinator can commit the test, checker and documentation/record slices in that order. Clean-head pre-freeze and the next independent round remain for the following stages.

## Stage H — clean-head runs after round 1 fixes

Tested clean committed head: `0fbdefb77055e1124fe23831c6977f9acbd35ff4` on `T-c826`.
Base: `origin/main`, resolved to `4259e05f113dec28b60ada61fa245b5f94f72980`.
All nine commands ran before this record was edited. UTC start/end values were
read directly with `date -u +%FT%TZ`; exit statuses were read from each command.
`TMPDIR=/home/layup/.cache/layup-session/tmp-green` confined temporary files to
scratch. Each long-running command printed its step and an eight-second heartbeat.

| Run | Command | UTC start | UTC end | Exit | Observation |
| --- | --- | --- | --- | --- | --- |
| H-freeze | `sh docs/tests/freeze-check.sh --task T-c826 --base origin/main --budget 1300/170 --evidence runs/T-c826/ci-exception.md --evidence runs/T-c826/conditions.md --evidence runs/T-c826/mutants.md --evidence runs/T-c826/real-heads.md --evidence runs/T-c826/test-runs.md > /home/layup/.cache/layup-session/tmp-green/c826-fc4.txt 2>&1` | 2026-10-20T22:29:49Z | 2026-10-20T22:29:50Z | 1 | FAIL K5; K1–K4 and K6 pass. Complete output below. |
| H-suite-sh | `sh docs/tests/freeze-check/run-all.sh` | 2026-10-20T22:29:50Z | 2026-10-20T22:29:59Z | 0 | 85 passed, 0 failed. |
| H-suite-bash-posix | `bash --posix docs/tests/freeze-check/run-all.sh` | 2026-10-20T22:29:59Z | 2026-10-20T22:30:07Z | 0 | 85 passed, 0 failed. |
| H-adr | `sh docs/adr/adr-lint.sh` | 2026-10-20T22:30:07Z | 2026-10-20T22:30:08Z | 0 | adr-lint: OK. |
| H-links | `sh docs/links/link-lint.sh` | 2026-10-20T22:30:08Z | 2026-10-20T22:30:20Z | 0 | link-lint: OK; 3001 links resolved. |
| H-discipline | `sh docs/tests/run-discipline-tests.sh` | 2026-10-20T22:30:20Z | 2026-10-20T22:30:23Z | 0 | 138 passed, 0 failed. |
| H-setup | `sh docs/setup/setup-check.sh` | 2026-10-20T22:30:23Z | 2026-10-20T22:30:45Z | 0 | All 12 setup checks and kit-linters pass; 138 discipline cases pass, ADR/PRD lints pass, 3001 links resolve. |
| H-task-check | `sh runs/T-c826/checks.sh` | 2026-10-20T22:30:45Z | 2026-10-20T22:30:45Z | 0 | 12 assertions pass; no PR argument supplied, so exact migration-PR identity remains pending in this invocation. Complete output below. |
| H-diff | `git diff --check origin/main HEAD` | 2026-10-20T22:30:45Z | 2026-10-20T22:30:45Z | 0 | No whitespace diagnostics for the committed diff. |

Complete freeze-check output, captured in
`/home/layup/.cache/layup-session/tmp-green/c826-fc4.txt`:

```text
ok K1: white space
ok K2: registration of new tests
ok K3: test selectors
ok K4: restore evidence of mutants
FAIL K5: -:0: projection: figure 4172/192, maximum 1300/170, tolerance 260/4, limit 1560/174, margin -2612/-18 (the figure at landing is still measured at landing)
ok K6: conditions table
```

Complete task-check output for the requested invocation without a PR argument:

```text
ok: canonical margin sentence in the tolerance section
ok: unchanged runner: docs/tests/freeze-check/run.sh
ok: unchanged runner: docs/tests/freeze-check/run-all.sh
ok: all nine baseline paths stay in the restore list
ok: kept checker requires the pull_request event
ok: one literal PR number and only the checker path can be kept
pending: supply the actual migration PR number before S2
ok: event guard receives the actual GitHub event name
ok: number guard receives the actual pull request number
ok: only the guarded checker skips restoration; both runners and other paths restore for every event
ok: restoration still reads the default branch of this repository
ok: default branch still comes from the repository event
ok: exception comment names the task, literal PR and removal issue
```

Stage H outcome: **failed**. The required freeze-check exit 0 was not met:
K5 reports figure **4172/192**, maximum **1300/170**, tolerance **260/4**,
limit **1560/174** and margin **-2612/-18**. All eight remaining commands exit 0.
No fix, budget change, base change, task close-out or completed-log edit was made.
This append records the failure; it does not claim a successful freeze.

| Run | Command | UTC start | UTC end | Exit | Observation |
| --- | --- | --- | --- | --- | --- |
| H-record-diff | `git diff --check` | 2026-10-20T22:31:13Z | 2026-10-20T22:31:13Z | 0 | No whitespace diagnostics in the evidence append. |
