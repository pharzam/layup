# T-8ywj — ADR-0012: build LAYUP in bootstrap mode; revert F-0005; start product work

Issue: [#56](https://github.com/pharzam/layup/issues/56). Decision record:
[ADR-0012](../adr/0012-build-layup-in-bootstrap-mode.md). Evidence:
[`runs/T-8ywj/`](../../runs/T-8ywj/diagnosis.md).

## Test runs

Run in the worktree of this task on 2026-09-25, base `43a05c8`.

| # | State of the tree | Command | Exit | Output |
| - | ----------------- | ------- | ---- | ------ |
| 1 | The base | `git grep -l 'F-0005' -- . ':!runs/' ':!docs/adr/'` | 0 | 8 files: `README.md`, the `F-0005` record, `docs/facts/README.md`, `docs/onboarding-for-engineers.md`, `docs/setup/setup-check.sh`, the `good-passthrough` fixture index, `docs/tasks/T-q1x6.md`, `docs/tasks/completed.md`. The sweep of rows 1 and 3 does not exclude `docs/tasks/`; this file, which is dated history, names the fact it removed, so the later sweep in the ADR and the fixes reply excludes `docs/tasks/` as well |
| 2 | `git revert -m 1 --no-commit 43a05c8` applied (12 files, +6 −176) | `sh docs/setup/setup-check.sh --only facts .` | 0 | `setup-check: facts OK` |
| 3 | The same | the `git grep` of row 1 | 1 | no file |
| 4 | Commits `44f6552` (the revert) and `4c1bd86` (the evidence); the ADR file present, no index row | `sh docs/adr/adr-lint.sh` | 1 | `FAIL 0012-build-layup-in-bootstrap-mode.md: no row for it in README.md's index table`; `WARN … nothing outside adr/ LINKS it` |
| 5 | The index row, the ADR-0005 and ADR-0006 statuses, and every document edit present, before the three commits that hold them | `sh docs/adr/adr-lint.sh`; `sh docs/prd/prd-lint.sh`; `sh docs/links/link-lint.sh`; `sh docs/tests/run-discipline-tests.sh`; `git diff --check` | 0 each | `adr-lint: OK`; `prd-lint: OK`; `link-lint: OK  827 links resolved`; `run-discipline-tests: 81 passed, 0 failed`; nothing |
| 6 | The same | `sh docs/setup/setup-check.sh` | 1 | 14 checks `OK`; the one line `setup-check: kit-history FAIL orphan: docs/tasks/T-8ywj.md has no line with T-8ywj in backlog.md or completed.md`, which the close-out commit ends (the completed-log line cannot exist before the round finishes) |
| 7 | Frozen head `07b0467`, round 1, first harness | `claude -p --model claude-opus-5-5 --dangerously-skip-permissions` on the brief, in a disposable clone, 900 s watchdog | 142 | no output in 900 s (stdout and stderr empty); a 5 s probe of the same command answered `OK`, so the harness was up. Skipped and recorded (Bootstrap mode rule 4) |
| 8 | The same head, round 1, second harness | `devin -p --prompt-file .review-brief/PROMPT.md --model gpt-6-sol-xhigh --permission-mode dangerous --respect-workspace-trust false`, in the same clone | 0 | the review record, 406 s: `material`, 13 findings (posted on #56) |
| 9 | Round 1, finding 1, reproduced | `sh /dev/fd/3 --only facts . 3< <(git show 43a05c8:docs/setup/setup-check.sh)` on the branch tree (the check as CI restores it from `main`) | 1 | `listed: docs/facts/operator-routing-policy.md is not in docs/setup/facts.sha256`; `index: … no row for F-0005`. Fix: branch `T-8ywj-pre` (`35819c6`) drops both names from `main`'s check first; `sh docs/setup/tests/run.sh` there: `30 passed, 0 failed` |
| 10 | Fix head `1771c38`, round 2 (cycle 1, the cap), first named harness | `devin -p … --model gpt-6-sol-xhigh …`, in a new disposable clone | 0 | the review record, 730 s: `not mergeable, findings recorded`, 5 findings (posted on #56) |

## Verdict

Ended at the cycle cap (1) with `not mergeable, findings recorded`; split to the
successor [`T-7sbn`](T-7sbn.md) ([#59](https://github.com/pharzam/layup/issues/59)),
which carries the branch at `1771c38` and takes all five open findings. Delivered
on the branch: [ADR-0012](../adr/0012-build-layup-in-bootstrap-mode.md) (Accepted),
the revert of PR #50 (`44f6552`), the evidence under `runs/T-8ywj/`, the
`## Bootstrap mode` section and its pointers, the backlog lines of the product
path, and the prerequisite PR #57. Round 1 (GPT-6 Sol on Devin): `material`, 13
findings — 12 fixed in `1771c38`, one (the CI restore of `main`'s check) answered
by PR #57. Round 2 (GPT-6 Sol on Devin): 5 findings, all text or a recorded limit,
fixed in the successor.

## Resource record

Recorded, not budgeted (ADR-0007). Token counts: `not reported` (Devin prints none;
the author's session count was not read).

| Part | Expected tier | Model | Effort | Tokens | Elapsed |
| ---- | ------------- | ----- | ------ | ------ | ------- |
| The audit (three read-only subagents in a disposable clone, and the author's reads) | reasoning | Claude Fable 5.1 on Claude Code (author and subagents) | not reported | not reported | about 40 min |
| The plan and its review | reasoning | Claude Fable 5.1 (plan); the Operator (review) | not reported | — | about 10 min |
| The build: the revert, the evidence, the ADR, the documents | execution | Claude Fable 5.1 — a reasoning-tier model on an execution part: the author's session did the edits | not reported | not reported | about 35 min |
| Round 1 | reasoning | Claude Opus 5.5 on Claude Code (no output in 900 s, skipped); GPT-6 Sol on Devin (the record) | — ; xhigh | not reported | 15 min + 406 s |
| The fixes after round 1, PR #57 | execution | Claude Fable 5.1 | not reported | not reported | about 20 min |
| Round 2 | reasoning | GPT-6 Sol on Devin | xhigh | not reported | 730 s |
| **Total** | | | | not reported | about 2 h 20 min |
