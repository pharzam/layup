# T-7sbn — ADR-0012 (bootstrap mode): successor of T-8ywj after its cycle cap

Issue: [#59](https://github.com/pharzam/layup/issues/59); predecessor
[`T-8ywj`](T-8ywj.md) ([#56](https://github.com/pharzam/layup/issues/56), split).
Decision record: [ADR-0012](../adr/0012-build-layup-in-bootstrap-mode.md).
Evidence: [`runs/T-8ywj/`](../../runs/T-8ywj/diagnosis.md) (the audit and the
Operator's routing text) and [`runs/T-7sbn/test-runs.txt`](../../runs/T-7sbn/test-runs.txt)
(this task's check outputs and round runs).

## Test runs

Run in the worktree of this task on 2026-09-25, base `43a05c8`, branch carried
from `1771c38`.

| # | State of the tree | Command | Exit | Output |
| - | ----------------- | ------- | ---- | ------ |
| 1 | The fixes of #56 round 2 and O-43 applied; `T-8ywj.md` with its verdict; this file absent | `sh docs/adr/adr-lint.sh`; `sh docs/prd/prd-lint.sh`; `sh docs/tests/run-discipline-tests.sh`; `git diff --check` | 0 each | `adr-lint: OK`; `prd-lint: OK`; `run-discipline-tests: 81 passed, 0 failed`; nothing |
| 2 | The same | `sh docs/links/link-lint.sh` | 1 | `FAIL  L1: docs/tasks/T-8ywj.md:27 links T-7sbn.md, but that path does not exist` — the right reason: this file did not exist yet |
| 3 | The same | `sh docs/setup/setup-check.sh` | 1 | 11 checks `OK`; the one line `setup-check: kit-history FAIL orphan: docs/tasks/T-8ywj.md has no line with T-8ywj in backlog.md or completed.md` |
| 4 | This file present; before the two commits `650904a` and `7e4c6e3` (the frozen head) | the five checks of row 1 and 2, and `sh docs/setup/setup-check.sh` | 0; 1 | `adr-lint: OK`; `prd-lint: OK`; `link-lint: OK  836 links resolved`; `run-discipline-tests: 81 passed, 0 failed`; `git diff --check` nothing; setup: 11 checks `OK` and the two orphan lines for `T-7sbn.md` and `T-8ywj.md` (outputs in `runs/T-7sbn/test-runs.txt`) |
| 5 | Frozen head `7e4c6e3`, round 1 (cycle 0), first named harness | `devin -p --prompt-file .review-brief/PROMPT.md --model gpt-6-sol-xhigh --permission-mode dangerous --respect-workspace-trust false`, in a disposable clone | 0 | the review record: `material`, 3 material findings and 4 notes (posted on #59) |
| 6 | The fixes of round 1 present, before the fix commit `0393289` | the five checks, `sh docs/setup/setup-check.sh`, `sh docs/tests/nested-checkout-check.sh` | 0; 1; 0 | all `OK` (837 links; 81 passed); setup: 11 checks `OK` and the two orphan lines; `nested-checkout-check: OK  10 cases behaved` (outputs in `runs/T-7sbn/test-runs.txt`) |
| 7 | Frozen head `0393289`, round 2 (cycle 1, the cap), first named harness | the same `devin -p …` command, in a new disposable clone | 0 | the review record, 385 s: `not mergeable, findings recorded`, 2 findings and 2 notes (posted on #59); the Operator ruled both findings notes (O-44) |
| 8 | The two notes of O-44 and the reviewer's glossary note applied; the verdict, the resource record and the completed-log lines present (the close-out) | the five checks, `sh docs/setup/setup-check.sh`, `git diff --check` | see `runs/T-7sbn/test-runs.txt` | recorded there at the close-out commit; expected: every check `OK`, no orphan line |

## Verdict

Delivered: [ADR-0012](../adr/0012-build-layup-in-bootstrap-mode.md) (Accepted;
O-41 to O-44), the `## Bootstrap mode` section of `docs/engineering-discipline.md`
with the Operator's materiality test and one home per rule (O-43), the revert of
PR #50, the evidence under `runs/T-8ywj/`, the marked summaries, and the product
path in the backlog. Round 1 (GPT-6 Sol on Devin): `material`, 3 findings and 4
notes, fixed in `0393289`. Round 2 (GPT-6 Sol on Devin, the cap): `not mergeable,
findings recorded`, 2 findings and 2 notes; the Operator settled both findings as
notes (O-44, rule 3), applied in the close-out commit with the glossary note. The
reviewer's note 2 stands as a recorded limit: the brief's task record named round
1's verdict, so round 2 did not hold full Context independence. Predecessor
[`T-8ywj`](T-8ywj.md) (#56) closed as split to this task. The prerequisite PR #57
merges first; this branch then merges `origin/main` in.

## Resource record

Recorded, not budgeted (ADR-0007). Token counts: `not reported` (Devin prints
none; the author's session count was not read). The predecessor's record is in
[`T-8ywj.md`](T-8ywj.md).

| Part | Expected tier | Model | Effort | Tokens | Elapsed |
| ---- | ------------- | ----- | ------ | ------ | ------- |
| The issue, the plan and its review | reasoning | Claude Fable 5.1 on Claude Code (author); the Operator (review) | not reported | — | about 10 min |
| The fixes of #56's findings, O-43 applied, the #56 record | execution | Claude Fable 5.1 — a reasoning-tier model on an execution part: the author's session did the edits | not reported | not reported | about 25 min |
| Round 1 | reasoning | GPT-6 Sol on Devin | xhigh | not reported | 350 s |
| The fixes of round 1, the run evidence | execution | Claude Fable 5.1 | not reported | not reported | about 15 min |
| Round 2 | reasoning | GPT-6 Sol on Devin | xhigh | not reported | 385 s |
| The close-out (O-44, the notes, this record, the PR) | `—` | Claude Fable 5.1 | not reported | not reported | about 15 min |
| **Total** | | | | not reported | about 1 h 20 min |
