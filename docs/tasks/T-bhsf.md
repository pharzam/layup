# T-bhsf — require the three Go CI jobs on main

Issue: [#35](https://github.com/pharzam/layup/issues/35) (parent [#29](https://github.com/pharzam/layup/issues/29)). PR: [#41](https://github.com/pharzam/layup/pull/41).

## Verdict

Delivered: the second branch-protection `PUT` (19:49:55 UTC) makes `main` require
all 12 CI jobs, the three Go jobs included; the read-back equals the file and is in
the setup record (V-21). Round 1: material (a wrong phrase about the first `PUT`'s
nine contexts), fixed. Round 2: nothing material.

## Resource record

Recorded, not budgeted (ADR-0007). Author tokens: `not reported`.

| Part | Expected tier | Model | Effort | Tokens | Elapsed |
| ---- | ------------- | ----- | ------ | ------ | ------- |
| The plan and its review | reasoning | Claude Opus 5.5; Claude Fable 5.1 | not reported | 154,939 (review) | 8 min (review) |
| The decay review rounds | reasoning | Claude Fable 5.1 | not reported | 214,464 (104,187 + 110,277) | 12 min (6 + 6) |
| Applying and recording the setting | execution | Claude Opus 5.5 (one tier: limit recorded) | not reported | not reported | 10 min |
| Isolate, guardrails, docs, close-out | `—` | Claude Opus 5.5 | not reported | not reported | 5 min |
| **Total** | | | | 369,403 reported + not reported | about 35 min |
