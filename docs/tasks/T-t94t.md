# T-t94t — a rejected plan goes back to the planner, and the plan review is told the target's plan form

Issue [#483](https://github.com/pharzam/layup/issues/483), a child of row 68b ([#370](https://github.com/pharzam/layup/issues/370)), plan row 68b3. It serves `F-0003#45`, through REQ-005 (the task loop).
Budget maximum 1100 lines over 28 files, close-out inside, under ADR-0028; cycle cap 1. Goal count 2, by the Operator's ruling: the replan route and the review's form.
Evidence: [test runs](../../runs/T-t94t/test-runs.md), [conditions](../../runs/T-t94t/conditions.md).

## Plans and reviews

- Plan [6340212486](https://github.com/pharzam/layup/issues/483#issuecomment-6340212486).
- The Operator's decisions [6340236265](https://github.com/pharzam/layup/issues/483#issuecomment-6340236265): a count of its own (`shapeCycleCap`), the plan-review instruction aligned in this task, a stop at the cap.
- Plan review [6340412848](https://github.com/pharzam/layup/issues/483#issuecomment-6340412848) (approve-with-conditions, six conditions, all met: `runs/T-t94t/conditions.md`).
- A first attempt at round 1 by Claude Fable failed on API credit and wrote no record.
- Round 1 [6341423977](https://github.com/pharzam/layup/issues/483#issuecomment-6341423977) (material: F1 to F3). Round 2 [6342237563](https://github.com/pharzam/layup/issues/483#issuecomment-6342237563) (not mergeable, findings recorded: F4, F5).
- The Operator's end 3 [6342237997](https://github.com/pharzam/layup/issues/483#issuecomment-6342237997): fix F4 and F5 test-first, no further round.

## What is delivered

- A plan review whose kept copy rejects with at least one reason is a valid handoff `plan review, task plan` (the item `reject-reasons`, the 14th default row put in force for older tables by `handoff.InForce`). The driver starts a new attempt at the same base with the plan stage reopened, gives the reasons to the planner as findings, counts rejections on its own, and stops at the second: `<task>: the plan review rejected the plan <n> times`. One rule of a rejection serves the count, the reopening and the reasons. The records of the earlier binary are recorded again with no repair.
- `qaKind` classifies by the first plan handoff, so a later plan moves no earlier review.
- The block `plan-review` (and its constant and golden file) states the target's plan form, that an empty `tests:` list is never a reason to reject, that LAYUP's own items are no part of it, and that a reject gives its reasons.
- The rule in `docs/spec/loop.md` (one home, with the table in force), `records.md`, the PRD and traceability cells, plan row 68b3 and a lesson in `guardrails.md` §2.

## Findings, and their closing

- F1 (round 1): a valid handoff over a lost kept copy authorized the developer. Fixed: an unread review is a stop that names it, in the driver and in the amendment reader.
- F2, F3 (round 1): the bare-reject stop named the wrong pair; `records.md` read as if every stored table needs fourteen pairs. Fixed.
- F4 (round 2): a customized row for the planner pair let a bare reject through as valid. Fixed under end 3: a reject with no reason stops whatever the table.
- F5 (round 2): the label fix misnamed an earlier binary's bare-reject handoff. Fixed under end 3: that failure keeps `plan review → developer`; recovery is not broadened.
- The stale count in `internal/run/commentform.go` was removed.

## Verification and verdict

Red first, then green, with `date -u` times in the test runs. Verified by running: build, vet (three tags), gofmt, unit, the full integration level and e2e at the round 1 fix, the integration packages touched at the end-3 fix, the six lints, `git diff --check`, and the freeze check (K1 to K6).

Round 1, comment 6341423977 (GPT-6 Astra on Codex): **material.** Round 2, comment 6342237563 (GPT-6 Astra on Codex): **not mergeable, findings recorded.**

**Verdict: not mergeable, findings recorded; F4 and F5 were fixed under the Operator's end 3, and that fix has had no round.** Goal count 2. Harness exception: ADR-0026's routing rule names `claude`, `devin` and `opencode`; Codex is used for the plan review on the Operator's message of 2026-10-13, and GPT-6 Astra on the Operator's message of 2026-10-22.

## Resource record

ADR-0007, recorded, not budgeted. Tokens as the harness reported them; build and fix times are commit times.

| Part | Tier | Model | Effort | Tokens | Elapsed |
| --- | --- | --- | --- | --- | --- |
| Plan | reasoning | Claude Opus 5.5 (a Plan agent) | not recorded | not kept | not kept |
| Plan review | reasoning | GPT-6.1 Sol (Codex) | xhigh | 245,229 | 2026-10-22T12:02:03Z to 12:10:45Z |
| Build | execution | Claude Sonnet 5.5 (subagent) | not recorded | not kept | 2026-10-22T12:15:53Z to 12:34:15Z (first to last commit) |
| Review round 1, failed attempt | reasoning | Claude Fable (API credit ran out) | not recorded | 7,337,830 (USD 5.68) | 2026-10-22T12:35:43Z to 12:44:06Z (no record) |
| Review round 1 | reasoning | GPT-6 Astra (Codex) | xhigh | 183,161 | 2026-10-22T12:48:01Z to 13:02:38Z |
| Fix of round 1 | execution | Claude Sonnet 5.5 (subagent) | not recorded | not kept | 2026-10-22T13:05:48Z to 13:12:11Z (commits) |
| Review round 2 | reasoning | GPT-6 Astra (Codex) | xhigh | 234,361 | 2026-10-22T13:14:37Z to 13:34:31Z |
| Fix under end 3 | execution | Claude Sonnet 5.5 (subagent) | not recorded | not kept | 2026-10-22T13:43:33Z (commit) |
| Close-out | execution | Claude Sonnet 5.5 (subagent) | not recorded | not kept | not kept; its commit follows 6dda3c1d (2026-10-22T13:43:33Z) |
| **Total** | | | | 8,000,581 reported (the four numeric entries); the complete total is not reported | not reported |
