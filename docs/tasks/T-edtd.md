# T-edtd — ADR-0011: architecture of the LAYUP core engine

Issue: [#30](https://github.com/pharzam/layup/issues/30) (parent [#29](https://github.com/pharzam/layup/issues/29)). PR: [#31](https://github.com/pharzam/layup/pull/31).

## Verdict

Delivered: [ADR-0011](../adr/0011-structure-the-core-engine-as-a-go-cli-over-repository-files.md)
after a three-member panel and the Operator decisions O-9 to O-13 (LAYUP is an
orchestrator product; LAYUP's machinery never goes into a target). Rounds 1 and 2:
material (a misquoted Operator decision; an unsettled reading, settled by O-13;
prose steps and step S12 with no owner), fixed. Round 3: nothing material.
Follow-ups: #33 (stack-gate runner), #34 (wording).

## Resource record

Recorded, not budgeted (ADR-0007). Author tokens: `not reported`.

| Part | Expected tier | Model | Effort | Tokens | Elapsed |
| ---- | ------------- | ----- | ------ | ------ | ------- |
| The plan and its review | reasoning | Claude Opus 5.5; Claude Fable 5.1 | not reported | 111,446 (review) | 5 min (review) |
| The panel (options, no verdict) | reasoning | Claude Opus 5.5 (A, C); Claude Fable 5.1 (B) | not reported | 284,035 (73,776 + 144,433 + 65,826) | 6 min (runs in parallel) |
| The decay review rounds | reasoning | Claude Fable 5.1 | not reported | 486,028 (141,651 + 186,931 + 157,446) | 29 min (9 + 11 + 8) |
| Writing the decision record | reasoning | Claude Opus 5.5 | not reported | not reported | 30 min |
| Isolate, guardrails, docs, close-out | `—` | Claude Opus 5.5 | not reported | not reported | 10 min |
| **Total** | | | | 881,509 reported + not reported | about 80 min |
