# T-7ndb — encode the 9 System Invariants into guardrails

Issue: [#7](https://github.com/pharzam/layup/issues/7) (parent [#1](https://github.com/pharzam/layup/issues/1)). PR: [#19](https://github.com/pharzam/layup/pull/19).

## Verdict

Delivered: `docs/guardrails.md` §1.1 holds Inv-1 to Inv-9, each with its
`F-0001#N`, the trap, and `Check: no check yet` (honest: no gate runs
`setup-check.sh` yet); §2 holds three lessons from earlier tasks; check
`guardrails`. Round 1: material (an untested cause; a directory passed as a check
path), fixed in cycle 1. Round 2: nothing material.

## Resource record

Recorded, not budgeted (ADR-0007). Author tokens: `not reported`.

| Part | Expected tier | Model | Effort | Tokens | Elapsed |
| ---- | ------------- | ----- | ------ | ------ | ------- |
| The plan and its review | reasoning | Claude Opus 5.5; Claude Fable 5.1 (shared) | not reported | shared with #2 | shared with #2 |
| The decay review rounds | reasoning | Claude Fable 5.1 | not reported | 275,330 (155,548 + 119,782) | 21 min (10 + 11) |
| Writing the tests and the code | execution | Claude Opus 5.5 (one tier: limit recorded) | not reported | not reported | 25 min |
| Isolate, guardrails, docs, close-out | `—` | Claude Opus 5.5 | not reported | not reported | 10 min |
| **Total** | | | | 275,330 reported + not reported | about 56 min |
