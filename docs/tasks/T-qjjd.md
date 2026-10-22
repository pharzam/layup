# T-qjjd — a crashed task session gets one new attempt

Issue [#475](https://github.com/pharzam/layup/issues/475), a child of row 68b ([#370](https://github.com/pharzam/layup/issues/370)), plan row 68b1. It serves `F-0003#45`, through REQ-005 (the task loop).
Budget maximum 900 lines over 24 files, close-out inside, under ADR-0028; cycle cap 1. Goal count 1.
Evidence: [test runs](../../runs/T-qjjd/test-runs.md), [conditions](../../runs/T-qjjd/conditions.md).

## Plans and reviews

- First plan [6334945091](https://github.com/pharzam/layup/issues/475#issuecomment-6334945091), rejected by the first plan review [6336038815](https://github.com/pharzam/layup/issues/475#issuecomment-6336038815).
- The Operator's decisions [6334968394](https://github.com/pharzam/layup/issues/475#issuecomment-6334968394), and the split of the stderr tail into #476 [6336044551](https://github.com/pharzam/layup/issues/475#issuecomment-6336044551).
- Revised plan [6336109085](https://github.com/pharzam/layup/issues/475#issuecomment-6336109085); second plan review [6336309085](https://github.com/pharzam/layup/issues/475#issuecomment-6336309085) (approve-with-conditions, four conditions, all met: `runs/T-qjjd/conditions.md`).
- Round 1 [6337025749](https://github.com/pharzam/layup/issues/475#issuecomment-6337025749).

## What is delivered

- A task session of the plan, the plan review, the developer or the verifier whose end is `crash` or `no-result`, and not refused, gets one new attempt n+1 at the base of the attempt it replaces; the bound is once per stage, counted by the eligible end that caused each replacement. `wall`, `output`, `start`, the close-out and the answer session stay stops.
- One rule (`retries`, `replacedOf`, `priorAttempt` in `internal/run/shapeloop.go`), read by `findingsOf`, `PromptFile` and `AnswerAccepted` (through `Position` and `OpenQuestions`).
- The rule in `docs/spec/loop.md` (The driver of a Shape task), the "Again?" table of `docs/spec/session.md`, plan row 68b1, the REQ-005 cells of the PRD, four traceability rows and one lesson in `guardrails.md` §2.

## Not delivered

The stderr tail of a crashed session is [#476](https://github.com/pharzam/layup/issues/476). The retry of the close-out and of the answer session, of `wall` and `output`, and every further retry are `M3d`'s.

## Findings of round 1, and their closing

Round 1 found nothing material in scope. Its notes:

- Note 1 (the three guard tests the plan listed): declined after the last round. The excluded stages are pinned by the single path of `outcome` and the existing `wall` and `failed` pins, and the plan retry is pinned at unit level. They can be added in a later task.
- Note 2 (the loop.md wording "the four retries add at most five task sessions"): applied as a document change.
- Note 3 (a crashed developer's push, if any, is on its own attempt's branch and in no later diff, so it is no author; the place of the `retry` row after the question and before the stages): applied as a document change.
- Note 4 (a shared constructor of the rule's input): closed with no change; condition 1 is met in substance, the three readers take the one `replacedOf(loopView)`.
- Note 5 (the red evidence): closed with no change; the red evidence is the record in `runs/T-qjjd/test-runs.md`, as the plan review allowed.

## Verification and verdict

Red first, then green, with `date -u` times in the test runs. Verified by running: build, vet (three tags), gofmt, unit, the full integration level (exit 0), the e2e level (exit 0), the six lints, `git diff --check`, and the freeze check (K1 to K6).

Round 1, comment 6337025749 (Claude Fable 5.1 on OpenCode): **nothing material in scope.**

**Verdict: nothing material in scope.** Goal count 1. Harness exception: ADR-0026's routing rule names `claude`, `devin` and `opencode`; Codex is used for the two plan reviews on the Operator's message of 2026-10-13.

## Resource record

ADR-0007, recorded, not budgeted. Tokens as the harness reported them; build and close-out times are commit times.

| Part | Tier | Model | Effort | Tokens | Elapsed |
| --- | --- | --- | --- | --- | --- |
| Plans (first, revised) | reasoning | Claude Opus 5.5 (Plan agent) | not recorded | not kept | not kept |
| Plan review 1 | reasoning | GPT-6.1 Sol (Codex) | xhigh | 224,924 | 2026-10-22T07:09:43Z to 07:19:03Z |
| Plan review 2 | reasoning | GPT-6.1 Sol (Codex) | xhigh | 212,888 | 2026-10-22T08:17:15Z to 08:27:45Z |
| Build | execution | Claude Sonnet 5.5 (subagent) | not recorded | not kept | 2026-10-22T08:38:58Z to 09:07:44Z (first to last commit) |
| Review round 1 | reasoning | Claude Fable 5.1 (OpenCode, max) | max | 17,636,050 (USD 11.11) | 2026-10-22T08:46:04Z to 09:06:47Z |
| Close-out | execution | Claude Sonnet 5.5 (subagent) | not recorded | not kept | not kept; its commit follows 4eb79c2 at 2026-10-22T09:07:44Z |
| **Total** | | | | 18,073,862 reported (the three numeric entries); the complete total is not reported | not reported |
