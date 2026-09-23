# T-fvng — restore the pr-link and review-record lint scripts from main

Issue: [#23](https://github.com/pharzam/layup/issues/23) (blocker, revealed by #9). PR: [#25](https://github.com/pharzam/layup/pull/25).

## Verdict

Delivered: `pr-link.yml` and `review-record.yml` restore their lint script from
`main` (CI logs show it); job `setup-check` restores the two `docs/ci` scripts;
check `ci` has cause `restore`. Rounds 1 and 2: material (the check could be
passed by a name line, other run forms, other indents, `.yaml`, CRLF); fixed by a
rule change in cycle 2. Round 3: nothing material. Follow-ups on #21 and #24.

## Resource record

Recorded, not budgeted (ADR-0007). Author tokens: `not reported`.

| Part | Expected tier | Model | Effort | Tokens | Elapsed |
| ---- | ------------- | ----- | ------ | ------ | ------- |
| The plan and its review | reasoning | Claude Opus 5.5; Claude Fable 5.1 | not reported | 108,245 (review) | 6 min (review) |
| The decay review rounds | reasoning | Claude Fable 5.1 | not reported | 449,478 (114,132 + 152,557 + 182,789) | 33 min (7 + 12 + 14) |
| Writing the tests and the code | execution | Claude Opus 5.5 (one tier: limit recorded) | not reported | not reported | 35 min |
| Isolate, guardrails, docs, close-out | `—` | Claude Opus 5.5 | not reported | not reported | 10 min |
| **Total** | | | | 557,723 reported + not reported | about 84 min |
