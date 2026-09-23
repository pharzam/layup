# T-9mmm — record the setup procedure and close the baseline

Issue: [#11](https://github.com/pharzam/layup/issues/11) (parent [#1](https://github.com/pharzam/layup/issues/1)). PR: [#28](https://github.com/pharzam/layup/pull/28).

## Verdict

Delivered: `docs/setup/README.md` and `steps.tsv` (15 steps, 4 with a human
decision) for a later CLI and TUI; the baseline summary in the setup record; the
kit section "How to adapt this kit" replaced; check `procedure`. Rounds 1 and 2:
material (a step that would re-create blocker #23; missing hooks evidence; a wrong
count), fixed. Round 3: nothing material.

## Resource record

Recorded, not budgeted (ADR-0007). Author tokens: `not reported`.

| Part | Expected tier | Model | Effort | Tokens | Elapsed |
| ---- | ------------- | ----- | ------ | ------ | ------- |
| The plan and its review | reasoning | Claude Opus 5.5; Claude Fable 5.1 (shared) | not reported | shared with #2 | shared with #2 |
| The decay review rounds | reasoning | Claude Fable 5.1 | not reported | 493,858 (179,245 + 142,713 + 171,900) | 33 min (12 + 12 + 9) |
| Writing the tests and the code | execution | Claude Opus 5.5 (one tier: limit recorded) | not reported | not reported | 35 min |
| Isolate, guardrails, docs, close-out | `—` | Claude Opus 5.5 | not reported | not reported | 10 min |
| **Total** | | | | 493,858 reported + not reported | about 78 min |
