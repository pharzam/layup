# T-8ywj — ADR-0012: build LAYUP in bootstrap mode; revert F-0005; start product work

Issue: [#56](https://github.com/pharzam/layup/issues/56). Decision record:
[ADR-0012](../adr/0012-build-layup-in-bootstrap-mode.md). Evidence:
[`runs/T-8ywj/`](../../runs/T-8ywj/diagnosis.md).

## Test runs

Run in the worktree of this task on 2026-09-25, base `43a05c8`.

| # | State of the tree | Command | Exit | Output |
| - | ----------------- | ------- | ---- | ------ |
| 1 | The base | `git grep -l 'F-0005' -- . ':!runs/' ':!docs/adr/'` | 0 | 8 files: `README.md`, `docs/facts/F-0005-operator-routing-policy.md`, `docs/facts/README.md`, `docs/onboarding-for-engineers.md`, `docs/setup/setup-check.sh`, the `good-passthrough` fixture index, `docs/tasks/T-q1x6.md`, `docs/tasks/completed.md` |
| 2 | `git revert -m 1 --no-commit 43a05c8` applied (12 files, +6 −176) | `sh docs/setup/setup-check.sh --only facts .` | 0 | `setup-check: facts OK` |
| 3 | The same | the `git grep` of row 1 | 1 | no file |
| 4 | Commits `44f6552` (the revert) and `4c1bd86` (the evidence); the ADR file present, no index row | `sh docs/adr/adr-lint.sh` | 1 | `FAIL 0012-build-layup-in-bootstrap-mode.md: no row for it in README.md's index table`; `WARN … nothing outside adr/ LINKS it` |
| 5 | The index row, the ADR-0005 and ADR-0006 statuses, and every document edit present, before the three commits that hold them | `sh docs/adr/adr-lint.sh`; `sh docs/prd/prd-lint.sh`; `sh docs/links/link-lint.sh`; `sh docs/tests/run-discipline-tests.sh`; `git diff --check` | 0 each | `adr-lint: OK`; `prd-lint: OK`; `link-lint: OK  827 links resolved`; `run-discipline-tests: 81 passed, 0 failed`; nothing |
| 6 | The same | `sh docs/setup/setup-check.sh` | 1 | 14 checks `OK`; the one line `setup-check: kit-history FAIL orphan: docs/tasks/T-8ywj.md has no line with T-8ywj in backlog.md or completed.md`, which the close-out commit ends (the completed-log line cannot exist before the round finishes) |
