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
