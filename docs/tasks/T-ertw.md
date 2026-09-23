# T-ertw — numbered facts for the rest of the PSB (F-0003)

Issue: [#43](https://github.com/pharzam/layup/issues/43) (parent [#42](https://github.com/pharzam/layup/issues/42), the PDR). PR: [#44](https://github.com/pharzam/layup/pull/44).

## Verdict

Delivered: [`F-0003`](../facts/F-0003-layup-psb-numbered-facts-part-2.md), 75
byte-exact facts over the PSB (§1, §3, §4, §5, the scope lists of §6, §7), so
PRD-0001 can cite them; check `facts` covers a list of records. Round 1: material
(the new index row hid a missing `F-0001` row from the check), fixed by anchoring
the ID. Round 2: nothing material. Both rounds re-extracted the 75 facts
independently and matched.

## Resource record

Recorded, not budgeted (ADR-0007). Author tokens: `not reported`.

| Part | Expected tier | Model | Effort | Tokens | Elapsed |
| ---- | ------------- | ----- | ------ | ------ | ------- |
| The plan and its review | reasoning | Claude Opus 5.5; Claude Fable 5.1 | not reported | 106,991 (review) | 7 min (review) |
| The decay review rounds | reasoning | Claude Fable 5.1 | not reported | 244,463 (111,682 + 132,781) | 18 min (9 + 9) |
| Writing the tests and the record | execution | Claude Opus 5.5 (one tier: limit recorded) | not reported | not reported | 25 min |
| Isolate, guardrails, docs, close-out | `—` | Claude Opus 5.5 | not reported | not reported | 5 min |
| **Total** | | | | 351,454 reported + not reported | about 55 min |
