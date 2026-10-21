# T-36xw — the bet brief fits `brief.lines` on a real Intake

Issue [#471](https://github.com/pharzam/layup/issues/471). It serves `F-0003#51` through REQ-012.
Plan [6318912584](https://github.com/pharzam/layup/issues/471#issuecomment-6318912584); plan review [6319073363](https://github.com/pharzam/layup/issues/471#issuecomment-6319073363); the Operator's decisions [6319102451](https://github.com/pharzam/layup/issues/471#issuecomment-6319102451).
Budget maximum 600 lines over 17 files, close-out inside, under ADR-0028; cycle cap 1. Goal count 1.
Rounds: 1 [6320087805](https://github.com/pharzam/layup/issues/471#issuecomment-6320087805) (material); 2 [6320550732](https://github.com/pharzam/layup/issues/471#issuecomment-6320550732) (nothing material in scope).
Evidence: [test runs](../../runs/T-36xw/test-runs.md), [conditions](../../runs/T-36xw/conditions.md).

## What is delivered

- The brief selection by block: span facts and non-Operator answers inline, the Operator's answers and open questions as exact ID ranges with a link.
- The annex (`bets/<m>/facts.md`, `open.md`) written before the session and pinned at the last commit that changed it; `git.LastChange`.
- The spec, records, packages, architecture, PRD, traceability and plan changes, and one guardrails lesson.
- The tests: the M2d-sized unit fixture, `TestLastChange`, the integration tests and the e2e stage 5 link and rendered facts check.

## Not delivered

The stale `open` setup questions (intake.md:435) are a known limit by the Operator's decision 5: no issue, bounded by the open-part rule.

## Notes closed with no change

- Round 2 note 4 (no test that a file is not written when nothing is left out): a test change after the last round.
- Note 5 (`packages.md:259` does not name this task's piece): the row stays true and the mechanism is registered at `packages.md:385`.
- Note 6 (`<the files of the brief>` in the pin's command): the wording has one reading in context.

## Verification and verdict

Red first, then green, with `date -u` times in the test runs; the full integration and e2e levels, the six lints, `git diff --check` and the freeze check passed.

**Verdict: nothing material in scope.** Goal count 1. Harness exception for the plan review: ADR-0026's routing rule names `claude`, `devin` and `opencode`; Codex is used on the Operator's message of 2026-10-13.

## Resource record

ADR-0007, recorded, not budgeted. Tokens as the harness reported them.

| Part | Tier | Model | Effort | Tokens | Elapsed |
| --- | --- | --- | --- | --- | --- |
| Plan | reasoning | Claude Opus 5.5 (Plan agent) | not recorded | 205,792 | 2026-10-21T14:36Z to 14:46Z |
| Plan review | reasoning | GPT-6.1 Sol (Codex) | xhigh | 190,952 | 2026-10-21T14:47:39Z to 14:55:17Z |
| Build | execution | Claude Sonnet 5.5 (subagent) | not recorded | not kept | 2026-10-21T15:00:04Z to 15:45:25Z (first and last commit) |
| Review round 1 | reasoning | Claude Fable 5.1 | max | 15,513,217 | 2026-10-21T15:26:54Z to 15:44:12Z |
| Review round 2 | reasoning | Claude Fable 5.1 | max | 10,363,310 | 2026-10-21T15:47:23Z to 16:06:06Z |
| Close-out | execution | Claude Sonnet 5.5 (subagent) | not recorded | not kept | 2026-10-21T16:07Z to 16:12Z |
| **Total** | | | | 26,273,271 reported (the four numeric entries and both rounds) | not reported (the rows have no common start and end) |
