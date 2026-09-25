# T-7sbn — ADR-0012 (bootstrap mode): successor of T-8ywj after its cycle cap

Issue: [#59](https://github.com/pharzam/layup/issues/59); predecessor
[`T-8ywj`](T-8ywj.md) ([#56](https://github.com/pharzam/layup/issues/56), split).
Decision record: [ADR-0012](../adr/0012-build-layup-in-bootstrap-mode.md).
Evidence: [`runs/T-8ywj/`](../../runs/T-8ywj/diagnosis.md) (the audit and the
Operator's routing text; this task adds none of its own).

## Test runs

Run in the worktree of this task on 2026-09-25, base `43a05c8`, branch carried
from `1771c38`.

| # | State of the tree | Command | Exit | Output |
| - | ----------------- | ------- | ---- | ------ |
| 1 | The fixes of #56 round 2 and O-43 applied; `T-8ywj.md` with its verdict; this file absent | `sh docs/adr/adr-lint.sh`; `sh docs/prd/prd-lint.sh`; `sh docs/tests/run-discipline-tests.sh`; `git diff --check` | 0 each | `adr-lint: OK`; `prd-lint: OK`; `run-discipline-tests: 81 passed, 0 failed`; nothing |
| 2 | The same | `sh docs/links/link-lint.sh` | 1 | `FAIL  L1: docs/tasks/T-8ywj.md:27 links T-7sbn.md, but that path does not exist` — the right reason: this file did not exist yet |
| 3 | The same | `sh docs/setup/setup-check.sh` | 1 | 12 checks `OK`; the one line `setup-check: kit-history FAIL orphan: docs/tasks/T-8ywj.md has no line with T-8ywj in backlog.md or completed.md` |
