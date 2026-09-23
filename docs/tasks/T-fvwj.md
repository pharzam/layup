# T-fvwj — store the PSB and the vision brief as raw facts

Issue: [#4](https://github.com/pharzam/layup/issues/4) (parent [#1](https://github.com/pharzam/layup/issues/1)). PR: [#16](https://github.com/pharzam/layup/pull/16).

## Verdict

Delivered: `F-0001` (PSB, 39 numbered verbatim facts) and `F-0002` (vision brief,
a solution document), both byte-identical, hashed in `docs/setup/facts.sha256`;
check `facts`. Round 1: material (a last hash line without a newline was skipped),
fixed in cycle 1. Round 2: nothing material. Lesson for `guardrails.md` §2 (with
`T-7ndb`): `while read` drops a last line that has no final newline.

## Resource record

Recorded, not budgeted (ADR-0007). Author tokens: `not reported`.

| Part | Expected tier | Model | Effort | Tokens | Elapsed |
| ---- | ------------- | ----- | ------ | ------ | ------- |
| The plan and its review | reasoning | Claude Opus 5.5; Claude Fable 5.1 (shared) | not reported | shared with #2 | shared with #2 |
| The decay review rounds | reasoning | Claude Fable 5.1 | not reported | 242,615 (103,666 + 138,949) | 15 min (7 + 8) |
| Writing the tests and the code | execution | Claude Opus 5.5 (one tier: limit recorded) | not reported | not reported | 25 min |
| Isolate, guardrails, docs, close-out | `—` | Claude Opus 5.5 | not reported | not reported | 10 min |
| **Total** | | | | 242,615 reported + not reported | about 50 min |
