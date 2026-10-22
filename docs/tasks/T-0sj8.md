# T-0sj8 — a comment file is no artifact, and a refused artifact gets one new attempt

Issue [#478](https://github.com/pharzam/layup/issues/478), a child of row 68b ([#370](https://github.com/pharzam/layup/issues/370)), plan row 68b2. It serves `F-0003#45`, through REQ-005 (the task loop).
Budget maximum 400 lines over 24 files, close-out inside, under ADR-0028; cycle cap 1. Goal count 1.
Evidence: [test runs](../../runs/T-0sj8/test-runs.md), [conditions](../../runs/T-0sj8/conditions.md).

## Plans and reviews

- Plan [6338082980](https://github.com/pharzam/layup/issues/478#issuecomment-6338082980).
- The Operator's decision on the separate bound [6338065517](https://github.com/pharzam/layup/issues/478#issuecomment-6338065517).
- Plan review [6338214816](https://github.com/pharzam/layup/issues/478#issuecomment-6338214816) (approve-with-conditions, four conditions, all met: `runs/T-0sj8/conditions.md`).
- Round 1 [6338807230](https://github.com/pharzam/layup/issues/478#issuecomment-6338807230).

## What is delivered

- The shared block `instruction-result` says that a comment file of `result/` is never an artifact and that a step whose work makes no commit writes no artifact row: in `docs/spec/loop.md`, `internal/run/instructions.go` and the 11 golden prompt files.
- A `done` session with the status `completed` and exactly one refusal, `artifact`, is an eligible end of its stage with its own retry: the bound is once per stage and class (`crash` and `no-result` share one, an artifact refusal has one), so the eight retries add at most ten task sessions. A refused result with the status `needs_context` keeps the question path.
- The rule in `docs/spec/loop.md` (The driver of a Shape task), the "Again?" table of `docs/spec/session.md`, plan row 68b2, the REQ-005 cells of the PRD, traceability and one lesson in `guardrails.md` §2.

## Findings of round 1, and their closing

Round 1 found nothing material in scope. Its notes:

- Finding 1 (the test-inventory row of `loop.md` reads two ways about "classes" and "a refused session"): applied as a document change.
- Finding 2 (plan row 68b2 lacks the `instruction-result` block): applied as a document change.
- Finding 3 (the "Again?" pointer sits in the `crash` row): applied as a document change; it is now in the `done` row.
- Finding 4 (the comment at `internal/run/shapeloop.go:74` on `posRetry` still says crash or no-result): declined, because it is a code-file edit after the last round.
- Finding 5 and the notes (the red case of the question precedence is a guard; the first assertion of the restart test): no fix needed.
- The flake of the integration level, `TestTheRulesResultOfAPushedHeadStaysThatOfItsBoundSessionWhenTheRecordsMoveOn`, is [#480](https://github.com/pharzam/layup/issues/480).

## Verification and verdict

Red first, then green, with `date -u` times in the test runs. Verified by running: build, vet (three tags), gofmt, unit, the full integration level (one flake, #480, then the package in full exit 0), the e2e level (exit 0), the six lints, `git diff --check`, and the freeze check (K1 to K6).

Round 1, comment 6338807230 (Claude Fable 5.1 on OpenCode): **nothing material in scope.**

**Verdict: nothing material in scope.** Goal count 1. Harness exception: ADR-0026's routing rule names `claude`, `devin` and `opencode`; Codex is used for the plan review on the Operator's message of 2026-10-13.

## Resource record

ADR-0007, recorded, not budgeted. Tokens as the harness reported them; build and close-out times are commit times.

| Part | Tier | Model | Effort | Tokens | Elapsed |
| --- | --- | --- | --- | --- | --- |
| Plan | reasoning | Claude Opus 5.5 (the coordinator) | not recorded | not kept | not kept |
| Plan review | reasoning | GPT-6.1 Sol (Codex) | xhigh | 190,312 | 2026-10-22T10:03:28Z to 10:09:54Z |
| Build | execution | Claude Sonnet 5.5 (subagent) | not recorded | not kept | 2026-10-22T10:13:40Z to 10:26:26Z (first to last commit) |
| Review round 1 | reasoning | Claude Fable 5.1 (OpenCode, max) | max | 12,077,857 (USD 8.38) | 2026-10-22T10:27:45Z to 10:42:39Z |
| Close-out | execution | Claude Sonnet 5.5 (subagent) | not recorded | not kept | not kept; its commit f86eca2 at 2026-10-22T10:43:52Z |
| **Total** | | | | 12,268,169 reported (the two numeric entries); the complete total is not reported | not reported |
