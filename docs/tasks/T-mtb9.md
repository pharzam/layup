# T-mtb9 — Go skeleton and gates (successor of T-t8qp after the split)

Issue: [#38](https://github.com/pharzam/layup/issues/38) (parent [#29](https://github.com/pharzam/layup/issues/29); source [#32](https://github.com/pharzam/layup/issues/32), `T-t8qp`). PR: [#39](https://github.com/pharzam/layup/pull/39).

## Verdict

Delivered: the first Go code, `layup version` (unit and e2e tests); the Go gates
active — CI jobs `lint`, `tests` and `security`, and the hook's lint and unit
level; the three contexts in `docs/setup/branch-protection.json` (the live `PUT`
is #35). #32 reached its cycle cap with one material doc claim open and was
split; this successor carried its work from `9578c17` and fixed that claim.
Round 1: nothing material.

## Resource record

Recorded, not budgeted (ADR-0007). Author tokens: `not reported`. The #32 figures
are carried, because this task continues its work.

| Part | Expected tier | Model | Effort | Tokens | Elapsed |
| ---- | ------------- | ----- | ------ | ------ | ------- |
| The plan and its review | reasoning | Claude Opus 5.5; Claude Fable 5.1 | not reported | 222,544 (134,991 for #32 + 87,553 for #38) | 14 min (8 + 6) |
| The decay review rounds | reasoning | Claude Fable 5.1 | not reported | 543,932 (123,460 + 129,574 + 150,319 on #32; 140,579 on #38) | 43 min (9 + 10 + 12 on #32; 9 on #38; rounded) |
| Writing the tests and the code | execution | Claude Opus 5.5 (one tier: limit recorded) | not reported | not reported | 45 min |
| Isolate, guardrails, docs, close-out | `—` | Claude Opus 5.5 | not reported | not reported | 15 min |
| **Total** | | | | 766,476 reported + not reported | about 117 min |
