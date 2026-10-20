# T-3ff9 command evidence

Base and tested HEAD: `f45fffc78070d96c9ab3aab15bbbf1e4e31272bc`.
Times below were read from `date -u +%FT%TZ` around each run; statuses came from the process.

| Command | Start (UTC) | End (UTC) | Exit | Output |
| --- | --- | --- | --- | --- |
| `git add runs/T-3ff9/checks.sh && git commit -m "test: T-3ff9 check delegated landing text before implementation" -m "Co-Authored-By: GPT-6.1 Sol <noreply@openai.com>"` | 2026-10-20T18:39:36Z | 2026-10-20T18:39:36Z | 128 | `fatal: Unable to create '/home/layup/projects/layup/.git/worktrees/T-3ff9/index.lock': Read-only file system` |
| `sh -n runs/T-3ff9/checks.sh` | 2026-10-20T18:40:21Z | 2026-10-20T18:40:21Z | 0 | No output. |

The first run stopped before hooks ran and before the red run. Stage A delegates the commit to the coordinating session and authorizes the S1 red run here. Initial read-only repository inspection was not timed; no test result is inferred from it.

## Stage A: S1 only

The reviewed plan, amendment 1 and the Operator's option-1 decision govern this check. K5 implementation and fixtures belong to #449. No operative document was edited.

| Command | Start (UTC) | End (UTC) | Exit | Output |
| --- | --- | --- | --- | --- |
| `sh -n runs/T-3ff9/checks.sh` | 2026-10-20T18:43:41Z | 2026-10-20T18:43:41Z | 0 | No output. |
| `sh runs/T-3ff9/checks.sh` | 2026-10-20T18:43:41Z | 2026-10-20T18:43:41Z | 1 | Full red output below. |

```text
ok: unit
FAIL: evidence
FAIL: executable-inputs
FAIL: tests-count
FAIL: input-definition
FAIL: inspection-and-dual-use
FAIL: input-defined-once
FAIL: bookkeeping
FAIL: mixed-files
FAIL: delegation
FAIL: tolerance
FAIL: budget-post-and-overrule
FAIL: fixed-bound
FAIL: exact-approval
FAIL: budget-block
FAIL: successor
FAIL: growth
FAIL: verdict-budget
FAIL: finding-budget
FAIL: cap-intro
FAIL: unread-fix-and-verdict
FAIL: default-end-3
FAIL: records-definition
FAIL: records-examples
FAIL: records-exclusions
FAIL: operator-end
FAIL: end-post-and-overrule
FAIL: end-exception
FAIL: retired old-verdict
FAIL: retired old-finding
FAIL: retired old-budget
FAIL: retired old-cap
FAIL: retired old-approval
FAIL: retired old-end-3
ok: end-1-preserved
ok: end-2-preserved
ok: end-4-preserved
FAIL: entry-point
FAIL: retired old-entry-point
FAIL: dod-pointer
ok: dod-gap
FAIL: retired old-dod-budget
FAIL: plan-pointer
FAIL: retired old-plan-budget
FAIL: cap-pointer
FAIL: ceiling-pointer
FAIL: revealed-pointer
FAIL: retired old-ceiling
FAIL: retired old-cap-row
FAIL: retired old-revealed-budget
FAIL: 0007-record-task-resource-use-status
ok: 0007-record-task-resource-use-only-status-changed
FAIL: 0007-record-task-resource-use index
FAIL: 0026-keep-the-bootstrap-review-rules-as-the-standing-gate-status
ok: 0026-keep-the-bootstrap-review-rules-as-the-standing-gate-only-status-changed
FAIL: 0026-keep-the-bootstrap-review-rules-as-the-standing-gate index
ok: index-only-status-cells-changed
ok: parser and K5 scope unchanged
```

Red for the intended reasons: the new count, execution-data definition, tolerance, default end 3 and pointers are absent; the old budget and cap assertions and unamended Status/index cells remain. The unit, ends 1/2/4, missing-budget caveat, immutable ADR content, unrelated index content, parser and K5 are green regression controls, not red evidence. Semantic agreement remains for the independent clause-by-clause review after S2/S3. No green run is claimed in S1.

Stage A scope verification:

| Command | Start (UTC) | End (UTC) | Exit | Output |
| --- | --- | --- | --- | --- |
| `git diff --check` | 2026-10-20T18:44:36Z | 2026-10-20T18:44:36Z | 0 | No output. |
| `git diff --exit-code` | 2026-10-20T18:44:36Z | 2026-10-20T18:44:36Z | 0 | No output. |
| `git status --short --untracked-files=all` | 2026-10-20T18:44:36Z | 2026-10-20T18:44:36Z | 0 | `?? runs/T-3ff9/checks.sh`; `?? runs/T-3ff9/test-runs.md` |

Tracked files are unchanged. The only pending files are the two S1 artifacts named in the commit plan. Commit and hook execution are left to the coordinating session under the Stage A instruction.

## Stage B: S2 and S3

Tested the uncommitted Stage B documents on `9c19effc8ed2d457bb255964b354b56c57e6c41e`; the task check is unchanged from S1. The binding inputs are the supplied reviewed plan, amendment 1 and Operator option-1 decision on #447. [semantic.md](semantic.md) is the author's clause comparison; independent semantic review remains pending. [conditions.md](conditions.md) maps the three review conditions; [mentions.tsv](mentions.tsv) records the reference sweep. Read-only inspection and editing commands were not timed; no test result is inferred from them.

Output lines in the table are folded with semicolons; all statuses and UTC endpoints were read directly around the command.

| Command | Start (UTC) | End (UTC) | Exit | Output |
| --- | --- | --- | --- | --- |
| `sh runs/T-3ff9/checks.sh` | 2026-10-20T18:48:04Z | 2026-10-20T18:48:04Z | 0 | `ok: unit; ok: evidence; ok: executable-inputs; ok: tests-count; ok: input-definition; ok: inspection-and-dual-use; ok: input-defined-once; ok: bookkeeping; ok: mixed-files; ok: delegation; ok: tolerance; ok: budget-post-and-overrule; ok: fixed-bound; ok: exact-approval; ok: budget-block; ok: successor; ok: growth; ok: verdict-budget; ok: finding-budget; ok: cap-intro; ok: unread-fix-and-verdict; ok: default-end-3; ok: records-definition; ok: records-examples; ok: records-exclusions; ok: operator-end; ok: end-post-and-overrule; ok: end-exception; ok: retired old-verdict; ok: retired old-finding; ok: retired old-budget; ok: retired old-cap; ok: retired old-approval; ok: retired old-end-3; ok: end-1-preserved; ok: end-2-preserved; ok: end-4-preserved; ok: entry-point; ok: retired old-entry-point; ok: dod-pointer; ok: dod-gap; ok: retired old-dod-budget; ok: plan-pointer; ok: retired old-plan-budget; ok: cap-pointer; ok: ceiling-pointer; ok: revealed-pointer; ok: retired old-ceiling; ok: retired old-cap-row; ok: retired old-revealed-budget; ok: 0007-record-task-resource-use-status; ok: 0007-record-task-resource-use-only-status-changed; ok: 0007-record-task-resource-use index; ok: 0026-keep-the-bootstrap-review-rules-as-the-standing-gate-status; ok: 0026-keep-the-bootstrap-review-rules-as-the-standing-gate-only-status-changed; ok: 0026-keep-the-bootstrap-review-rules-as-the-standing-gate index; ok: index-only-status-cells-changed; ok: parser and K5 scope unchanged` |
| `sh docs/adr/adr-lint.sh` | 2026-10-20T18:50:11Z | 2026-10-20T18:50:11Z | 0 | `adr-lint: OK` |
| `sh docs/prd/prd-lint.sh` | 2026-10-20T18:50:11Z | 2026-10-20T18:50:11Z | 0 | `prd-lint: OK` |
| `sh docs/links/link-lint.sh` | 2026-10-20T18:50:11Z | 2026-10-20T18:50:25Z | 0 | `link-lint: OK  2971 links resolved` |
| `git diff --check` | 2026-10-20T18:50:25Z | 2026-10-20T18:50:25Z | 0 | `No output.` |

All 58 task checks are green after the S1 red run (exit 1 at the unchanged baseline). The green includes the input definition and its end-3 reference, both condition-3 sentences, pointers, Status/index edits and preservation controls. No check was weakened or changed. ADR bodies, K5, its fixtures and the review-record parser remain unchanged.

Inspection found `origin/main` at `df9bfacf33054741564dbf0254a741f0d4be1e0f` (PR #448), ahead of the reviewed baseline `f45fffc78070d96c9ab3aab15bbbf1e4e31272bc`. This stage does not integrate branches or claim a final landing figure. The coordinating stages must use the current-main base at landing. T-3ff9 counts all evidence and close-out under its baseline 650/20 maximum and cap 1; its own proposed tolerance and default end 3 do not authorize its landing. K5 retains the whole-diff plus 70/2 projection until #449; any failure must remain recorded as a failure.

Stage B leaves only the files in its commit plan. The staging instruction delegates commits to the coordinating session; no commit, hook run, freeze-check, task close-out or independent round is claimed here. The remaining clean-head validation and freeze-check belong to Stage C.

Final Stage B whitespace check after writing the records: `git diff --check`, 2026-10-20T18:51:39Z to 2026-10-20T18:51:40Z, exit 0, no output. Scope inspection found exactly the 13 files in the commit plan; each has one final newline and no trailing whitespace. The conditions table has exactly three rows; its path/line spans and S1 commit resolve in the prospective committed tree.

## Stage C: committed-head gates and pre-freeze check

Tested HEAD `43c1bd4ceb04a46f790b9d845128c3b65979964d` against `origin/main` at `df9bfacf33054741564dbf0254a741f0d4be1e0f` after the clean merge of #448. The worktree stayed clean through all nine commands, including freeze-check; output and timing metadata were written outside it. No document or check was changed. Read-only instruction and repository inspection was not timed; no test result is inferred from it.

UTC endpoints below came from `date -u +%FT%TZ` around each command; exit statuses came directly from its process. Long runs printed progress every ten seconds. Outputs are folded or summarized, except the complete freeze-check output below.

| Command | Start (UTC) | End (UTC) | Exit | Output |
| --- | --- | --- | --- | --- |
| `sh runs/T-3ff9/checks.sh` | 2026-10-20T18:54:11Z | 2026-10-20T18:54:11Z | 0 | 58 task checks passed; no failures. |
| `sh docs/adr/adr-lint.sh` | 2026-10-20T18:54:11Z | 2026-10-20T18:54:12Z | 0 | adr-lint: OK |
| `sh docs/prd/prd-lint.sh` | 2026-10-20T18:54:12Z | 2026-10-20T18:54:12Z | 0 | prd-lint: OK |
| `sh docs/links/link-lint.sh` | 2026-10-20T18:54:12Z | 2026-10-20T18:54:25Z | 0 | link-lint: OK  2983 links resolved |
| `sh docs/tests/run-discipline-tests.sh` | 2026-10-20T18:54:25Z | 2026-10-20T18:54:28Z | 0 | run-discipline-tests: 138 passed, 0 failed |
| `sh docs/tests/freeze-check/run-all.sh` | 2026-10-20T18:54:28Z | 2026-10-20T18:54:34Z | 0 | freeze-check: 62 passed, 0 failed |
| `sh docs/setup/setup-check.sh` | 2026-10-20T18:54:34Z | 2026-10-20T18:54:54Z | 0 | setup-check: pin OK; setup-check: kit-history OK; setup-check: facts OK; setup-check: onboarding OK; setup-check: glossary OK; setup-check: guardrails OK; setup-check: markers OK; setup-check: adapted OK; setup-check: ci OK; setup-check: protection OK; setup-check: identity OK; setup-check: procedure OK; run-discipline-tests: 138 passed, 0 failed; adr-lint: OK; prd-lint: OK; link-lint: OK  2983 links resolved; setup-check: kit-linters OK |
| `git diff --check origin/main HEAD` | 2026-10-20T18:54:54Z | 2026-10-20T18:54:54Z | 0 | No output. |
| `sh docs/tests/freeze-check.sh --task T-3ff9 --base origin/main --budget 650/20 > /home/layup/.cache/layup-session/tmp-green/t3ff9-fc.txt 2>&1` | 2026-10-20T18:54:54Z | 2026-10-20T18:54:55Z | 0 | Complete output below; saved first to the external file named in the command. |

```text
ok K1: white space
ok K2: registration of new tests
ok K3: test selectors
ok K4: restore evidence of mutants
ok K5: projection: figure 516/16, maximum 650/20, margin 134/4 (the figure at landing is still measured at landing)
ok K6: conditions table
```

Committed whole-diff numstat: 400 added + 46 removed = 446 lines over 14 files, including all evidence. Unchanged K5 adds its default 70/2, giving 516/16 against 650/20 and margin 134/4. This is a pre-freeze projection, not the final landing figure; subsequent records and actual close-out still count for this task. No proposed tolerance or default end 3 is used.

The Stage A task check was red (exit 1); Stage B and this committed-head run are green (exit 0). Every Stage C gate passed; no fix is required. This stage appends only this record and delegates its commit to the coordinating session. No independent review, close-out or post-record committed-head check is claimed.

Final record whitespace check: `git diff --check`, 2026-10-20T18:55:29Z to 2026-10-20T18:55:29Z, exit 0, no output.
