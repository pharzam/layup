# T-0xr7 — the survivors are the bet's in every judging step, and the Operator may grant one more review cycle at the cap

Issue [#493](https://github.com/pharzam/layup/issues/493), a child of row 68b ([#370](https://github.com/pharzam/layup/issues/370)), plan row 68b5. It serves `F-0003#45` and `F-0003#51`, through REQ-005 (the task loop) and REQ-012 (the specification of a target).
Budget maximum 900 lines over 21 files, close-out inside, under ADR-0028; cycle cap 1. Goal count 2: the contract of the survivors and the cycle command.
Evidence: [test runs](../../runs/T-0xr7/test-runs.md), [conditions](../../runs/T-0xr7/conditions.md).

## Plans and reviews

- Plan [6348094347](https://github.com/pharzam/layup/issues/493#issuecomment-6348094347).
- The Operator's decisions [6348122670](https://github.com/pharzam/layup/issues/493#issuecomment-6348122670): the form `cycle <task> because <reason>`, one grant per task, the copied comment is the record, the wait names the command.
- Plan review [6348315001](https://github.com/pharzam/layup/issues/493#issuecomment-6348315001) (approve-with-conditions, six conditions, all met: `runs/T-0xr7/conditions.md`).
- Round 1 [6349262719](https://github.com/pharzam/layup/issues/493#issuecomment-6349262719): material (F1, F2). Plan amendment [6349268480](https://github.com/pharzam/layup/issues/493#issuecomment-6349268480) for F2.
- Round 2 [6349703235](https://github.com/pharzam/layup/issues/493#issuecomment-6349703235): not mergeable, findings recorded (F3).
- The Operator's end 3 [6349703634](https://github.com/pharzam/layup/issues/493#issuecomment-6349703634): fix F3 test-first, no further round.

## What is delivered

- The plan, the plan review and the verification of a task of kind architecture carry one summary sentence of the contract of the survivors: one survivor is written as decided, several are a section `Decision pending the bet`; a plan that chooses is corrected, the deferral itself is never a defect, and a wrong survivor, a hidden choice and a contradiction are still checked. It is in `docs/spec/loop.md` (one home), `internal/run/instructions.go` and three golden prompt files.
- The command `cycle <task> because <reason>` of the Operator on the control issue: judged form, author, issue, task, cap, bound; eligible only after the cap round's record and never on a `not mergeable, findings recorded` round; one grant per task, so the cap is 2; the copy in `copies.tsv` is the record and the wait names the command. The code is `internal/run/cycle.go`.
- The demo: the records of `T-am8b` (`internal/run/testdata/t-am8b.txtar`, layup-records at `e65db3e30027cee5e74f466c9fc88897fd60a8b3`) with one grant give attempt 4 at `d03c865` with the open findings of attempt 3.
- `records.md`, the PRD cells of REQ-005 and REQ-012, traceability, plan row 68b5 and a lesson in `guardrails.md` section 2.

## Findings, and their closing

- F1 (round 1): a reason padded with a second space passed the form and granted a cycle. Fixed; the parser rejects an empty word before the reason.
- F2 (round 1): the plan lacked the execution order and the replacement budget, the plan row and the corrected figure were missing. Fixed: plan amendment 6349268480 (the order is reconstructed after the build), plan row 68b5, conditions rows 1 and 6.
- F3 (round 2): a reason that starts with a tab still granted a cycle. Fixed under end 3, test-first; that fix has had no round.

## Verification and verdict

Red first, then green, with `date -u` times in the test runs (the hook blocks a red commit, so a test and its code are one commit). Verified by running: build, vet (three tags), gofmt, unit, the full integration level and e2e at the build, the integration package touched at each fix, the six lints, `git diff --check`, and the freeze check (K1 to K6).

Round 1, comment 6349262719 (GPT-6.1 Sol on Codex): **material.** Round 2, comment 6349703235 (GPT-6.1 Sol on Codex): **not mergeable, findings recorded.**

**Verdict: not mergeable, findings recorded; F3 was fixed under the Operator's end 3, and that fix has had no round.** Goal count 2. Harness exception: ADR-0026's routing rule names `claude`, `devin` and `opencode`; Codex is used for the plan review and for GPT-6.1 Sol on the Operator's message of 2026-10-13.

## Resource record

ADR-0007, recorded, not budgeted. Tokens as the harness reported them; build and fix times are commit times.

| Part | Tier | Model | Effort | Tokens | Elapsed |
| --- | --- | --- | --- | --- | --- |
| Plan | reasoning | Claude Opus 5.5 (a Plan agent) | not recorded | not kept | not kept |
| Plan review | reasoning | GPT-6.1 Sol (Codex) | xhigh | 213,082 | 2026-10-22T18:24:47Z to 18:33:33Z |
| Build | execution | Claude Sonnet 5.5 (subagent) | not recorded | not kept | 2026-10-22T18:42:28Z to 19:08:43Z (first to last commit) |
| Review round 1 | reasoning | GPT-6.1 Sol (Codex) | xhigh | 166,580 | 2026-10-22T19:10:49Z to 19:19:33Z |
| Fix of round 1 | execution | Claude Sonnet 5.5 (subagent) | not recorded | not kept | 2026-10-22T19:27:07Z to 19:28:24Z (commits) |
| Review round 2 | reasoning | GPT-6.1 Sol (Codex) | xhigh | 177,788 | 2026-10-22T19:30:03Z to 19:41:31Z |
| Fix under end 3 | execution | Claude Sonnet 5.5 (subagent) | not recorded | not kept | 2026-10-22T19:43:23Z (commit) |
| Close-out | execution | Claude Sonnet 5.5 (subagent) | not recorded | not kept | not kept; its commit follows 4ec88074 (2026-10-22T19:43:23Z) |
| **Total** | | | | 557,450 reported (the three numeric entries); the complete total is not reported | not reported |
