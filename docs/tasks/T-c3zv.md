# T-c3zv — the Shape options prompts name the empty evidence of a fail

Issue [#473](https://github.com/pharzam/layup/issues/473), a child of row 68b ([#370](https://github.com/pharzam/layup/issues/370)). It serves `F-0003#51`, through Shape of REQ-012.
Plan [6333505461](https://github.com/pharzam/layup/issues/473#issuecomment-6333505461); plan review [6333583535](https://github.com/pharzam/layup/issues/473#issuecomment-6333583535) (approve-with-conditions).
Budget maximum 120 lines over 7 files, close-out inside, under ADR-0028; cycle cap 1.
Goal count 1. Round 1: comment 6333827667.
Evidence: [test runs](../../runs/T-c3zv/test-runs.md), [conditions](../../runs/T-c3zv/conditions.md).

## What is delivered

- The two prompts of the options session (`optionsPromptText`, `rewritePromptText` in `internal/run/options.go`) say `the single character "—" (U+2014), the empty value, and nothing else; a hyphen - or any other text is refused`. The pass clause no longer says "never a dash".
- `TestEachOptionsPromptNamesTheEmptyEvidenceOfAFail` renders both prompts, takes the bytes between the prompt's quotes, reads a fail row with them through `spec.ReadOptions`, and keeps a `-` control the reader refuses.
- The complete rule in the options session of `docs/spec/loop.md`, the test in the traceability row and the REQ-005 and REQ-012 cells of PRD §12, and one lesson in `guardrails.md` §2.

The reader is unchanged: `-` stays refused, as `—` is the empty value of the TSV rules for every block.

## Findings of round 1, and their closing

- F2 (a records note): applied; the REQ-005 Test cell names the test in its own group, `(task T-c3zv, #473)`.
- F1 (the test does not bind the quoted character to the fail clause): closed with no change, because a test change after the last round would need another round.
- Note 1 (the apposition could be read as a list): closed with no change, because a prompt wording change after the last round is not made, and the plan review approved the wording as written.
- Note 2 ("both prompts of the session"): closed with no change; the parenthesis keeps the reading clear.
- Note 3 (the lesson's stricter check, the code point): closed with no change; the stricter check applies to new prompts.
- Notes 4 to 6 (the figure, the red in the green commit, no plan row): informational. The red test could not land alone because the pre-commit hook runs the unit level; the red run is recorded first in the test runs.

## Verification and verdict

Red first, then green, with `date -u` times in the test runs. Verified by running: build, vet (three tags), gofmt, unit, the integration subset `Options|Shape`, the six lints, `git diff --check`, and the freeze check (K1 to K6).

Round 1, comment 6333827667 (Claude Fable 5.1 on OpenCode): **nothing material in scope.**

**Verdict: nothing material in scope.** Goal count 1. Harness exception: ADR-0026's routing rule names `claude`, `devin` and `opencode`; Codex is used for the plan review on the Operator's message of 2026-10-13.

## Resource record

ADR-0007, recorded, not budgeted. Tokens as the harness reported them; build and close-out times are commit times.

| Part | Tier | Model | Effort | Tokens | Elapsed |
| --- | --- | --- | --- | --- | --- |
| Plan | reasoning | Claude Opus 5.5 (coordinator) | not recorded | not kept | not kept |
| Plan review | reasoning | GPT-6.1 Sol (Codex) | xhigh | 163,039 | 2026-10-22T05:29:58Z to 05:35:25Z |
| Build | execution | Claude Sonnet 5.5 (subagent) | not recorded | not kept | 2026-10-22T05:38:29Z to 05:40:32Z (first to last commit) |
| Review round 1 | reasoning | Claude Fable 5.1 (OpenCode, max) | max | 7,901,971 (USD 6.30) | 2026-10-22T05:41:48Z to 05:53:43Z |
| Close-out | execution | Claude Sonnet 5.5 (subagent) | not recorded | not kept | not kept; its commit d8fffc7 at 2026-10-22T05:55:24Z |
| **Total** | | | | 8,065,010 reported (the two numeric entries); the complete total is not reported | not reported |
