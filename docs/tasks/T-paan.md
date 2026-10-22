# T-paan — the plan, its review and the verification state that the cells and the facts record are the bet's

Issue [#488](https://github.com/pharzam/layup/issues/488), a child of row 68b ([#370](https://github.com/pharzam/layup/issues/370)). It serves `F-0003#51`, through REQ-012 (the specification of a target).
Budget maximum 160 lines over 12 files, close-out inside, under ADR-0028; cycle cap 1. Goal count 1.
Evidence: [test runs](../../runs/T-paan/test-runs.md), [conditions](../../runs/T-paan/conditions.md).

## Plans and reviews

- Plan [6343358889](https://github.com/pharzam/layup/issues/488#issuecomment-6343358889).
- Plan review [6343519315](https://github.com/pharzam/layup/issues/488#issuecomment-6343519315) (approve-with-conditions, four conditions, all met: `runs/T-paan/conditions.md`).
- Round 1 [6343833497](https://github.com/pharzam/layup/issues/488#issuecomment-6343833497): material (M1).
- Round 2 [6344192225](https://github.com/pharzam/layup/issues/488#issuecomment-6344192225): nothing material in scope.

## What is delivered

- The blocks `instruction-plan-of-a-task`, `instruction-plan-review` and `instruction-verification` carry one sentence, a summary of the contract: the initial specification task leaves the MoSCoW and Phase cells blank (not —) until the first bet and writes no numbered facts record, because code supplies both at the rendered commit; a reopened task keeps what exists. It is in `docs/spec/loop.md`, `internal/run/instructions.go` and the three golden prompt files.
- The complete contract and its review-scope limit have one home, the rendered-commit section of `docs/spec/loop.md`.
- The PRD REQ-012 cells, traceability and one lesson in `guardrails.md` §2.

## Findings, and their closing

- Round 1, M1 (the new unit test read `docs/spec/loop.md`, and a unit test touches no file): fixed. The four prompt assertions stay in the unit test; the document check moved to `TestTheInstructionContractOfTheCellsAndTheRecordHasItsHomeInLoopMd` at integration level.
- Round 2 notes. Note 1 (the comment at `internal/run/instructions_test.go:646` reverses the rule): declined, because it is a code-file edit after the last round. Note 2 (the unit-test file read that the change revealed): [#489](https://github.com/pharzam/layup/issues/489).

## Verification and verdict

Red first, with `date -u` times in the test runs (the hook blocks a red commit, so the test and the code are one commit). Verified by running: build, vet (three tags), gofmt, unit, the selected integration level, the six lints, `git diff --check`, and the freeze check (K1 to K6).

Round 2, comment 6344192225 (GPT-6 Astra on Codex): **nothing material in scope.**

**Verdict: nothing material in scope.** Goal count 1. Harness exception: ADR-0026's routing rule names `claude`, `devin` and `opencode`; Codex is used for the plan review and for GPT-6 Astra on the Operator's message of 2026-10-22.

## Resource record

ADR-0007, recorded, not budgeted. Tokens as the harness reported them; build, fix and close-out times are commit times.

| Part | Tier | Model | Effort | Tokens | Elapsed |
| --- | --- | --- | --- | --- | --- |
| Plan | reasoning | Claude Opus 5.5 (the coordinator) | not recorded | not kept | not kept |
| Plan review | reasoning | GPT-6.1 Sol (Codex) | xhigh | 126,390 | 2026-10-22T14:34:48Z to 14:41:30Z |
| Build | execution | Claude Sonnet 5.5 (subagent) | not recorded | not kept | 2026-10-22T14:45:36Z to 14:46:52Z (first to last commit) |
| Review round 1 | reasoning | GPT-6 Astra (Codex) | xhigh | 146,234 | 2026-10-22T14:48:43Z to 14:56:00Z |
| Fix of M1 | execution | Claude Sonnet 5.5 (subagent) | not recorded | not kept | 2026-10-22T14:57:12Z to 14:58:34Z (first to last commit) |
| Review round 2 | reasoning | GPT-6 Astra (Codex) | xhigh | 187,596 | 2026-10-22T15:00:08Z to 15:12:16Z |
| Close-out | execution | Claude Sonnet 5.5 (subagent) | not recorded | not kept | not kept; its commit follows f2d847f at about 2026-10-22T15:13:24Z |
| **Total** | | | | 460,220 reported (the three numeric entries); the complete total is not reported | not reported |
