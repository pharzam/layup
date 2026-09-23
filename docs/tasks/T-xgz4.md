# T-xgz4 — merge the PSB terms into the glossary

Issue: [#6](https://github.com/pharzam/layup/issues/6) (parent [#1](https://github.com/pharzam/layup/issues/1)). PR: [#18](https://github.com/pharzam/layup/pull/18).

## Verdict

Delivered: `docs/glossary.md` §1 holds the 25 PSB §8 rows word for word, each
citing `F-0001#15`–`#39` once; check `glossary`. Round 1: material (a false
"same meaning" for Operator; an example that named an unwired check). Round 2:
material (the onboarding status went stale; now it points to #1). Round 3:
nothing material.

## Resource record

Recorded, not budgeted (ADR-0007). Author tokens: `not reported`.

| Part | Expected tier | Model | Effort | Tokens | Elapsed |
| ---- | ------------- | ----- | ------ | ------ | ------- |
| The plan and its review | reasoning | Claude Opus 5.5; Claude Fable 5.1 (shared) | not reported | shared with #2 | shared with #2 |
| The decay review rounds | reasoning | Claude Fable 5.1 | not reported | 410,660 (127,705 + 135,367 + 147,588) | 28 min (7 + 12 + 9) |
| Writing the tests and the code | execution | Claude Opus 5.5 (one tier: limit recorded) | not reported | not reported | 25 min |
| Isolate, guardrails, docs, close-out | `—` | Claude Opus 5.5 | not reported | not reported | 10 min |
| **Total** | | | | 410,660 reported + not reported | about 63 min |
