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
