# T-fnsr — the test runs

Each run is on this host (macOS, `sh` of the system), in the worktree of the
task. The tables are the output of the scripts, shortened where marked.

## 1. `check.sh` on the empty target (red)

`sh runs/T-fnsr/check.sh https://github.com https://api.github.com pharzam/layup-uat`,
2026-10-08, before the run: exit 1.

```text
FAIL: no branch main
FAIL: no branch layup-records
2 failed
```

## 2. `check-test.sh`, the offline test of `check.sh` (condition 2)

First run: exit 1. Three faults of the test (the commit helper lost its `-m`;
`good` reused the loop variable `n`) and one of `check.sh`: `[bot]` in the
`grep` pattern of the authors is a bracket expression, so the check of the
records commits could not pass on a correct target
(`FAIL: 7 commits of layup-records are not by layup-agent[bot]`). Fixed with a
fixed-string comparison (`grep -cvxF`).

Run after the fixes: exit 0, 36 cases. A correct target and one with a
heartbeat commit pass; each other case gives its `FAIL` line once:

```text
ok: a correct target: exit 0
ok: a heartbeat commit: exit 0
ok: no repository: FAIL: the clone of
ok: no records branch: FAIL: no branch layup-records
ok: no bot response: FAIL: the ID of layup-agent[bot]: no answer
ok: two commits on main: FAIL: main has 2 commits
ok: a root commit by the Operator: FAIL: the root commit is by pharzam
ok: a root subject: FAIL: the root commit subject
ok: records by the Operator: FAIL: 7 commits of layup-records are not by
ok: a missing commit: FAIL: the commits of Start
ok: a vision brief: FAIL: the files:
ok: another README: FAIL: README.md is not the text
ok: another brief: FAIL: start/problem-statement.md is not the brief
ok: the header of start.tsv: FAIL: the header of start.tsv
ok: start.tsv <name>: FAIL: start.tsv <name> is 'x'   (14 cases, one per value)
ok: start.tsv app.permissions: FAIL: start.tsv app.permissions is empty
ok: start.tsv pin.time: FAIL: start.tsv pin.time is empty
ok: a harness row: FAIL: start.tsv has a harness row
ok: approvers.tsv: FAIL: approvers.tsv:
ok: the header of lease.tsv: FAIL: the header of lease.tsv
ok: a held lease: FAIL: lease.tsv: 6|0.1.0-dev|held
ok: the intake issue: FAIL: the intake issue #1 is by 'pharzam'
ok: the control issue: FAIL: the control issue #2 is by 'pharzam'
all cases pass
```

The response of `GET /users/layup-agent[bot]` is saved in
[`fixtures/users-bot.json`](fixtures/users-bot.json) (2026-10-08); the two
issue responses are written by the test.
