# T-nfh8 — fill the adopter markers with evidence or list them as open gaps

Issue: [#8](https://github.com/pharzam/layup/issues/8) (parent [#1](https://github.com/pharzam/layup/issues/1)). PR: [#20](https://github.com/pharzam/layup/pull/20).

## Verdict

Delivered: check `markers`; ADR-0010 (Go); record rows V-08 to V-20 and the
Operator decisions O-1 to O-8 in Git; 19 open-gap markers with 3 questions in
`setup/open-gaps.tsv`. Round 1: material (a mention hid a marker; fills with no
record row). Round 2: material (two stale test-layout claims; non-ASCII paths
skipped). Round 3: nothing material. Revealed follow-ups: #21.

## Resource record

Recorded, not budgeted (ADR-0007). Author tokens: `not reported`.

| Part | Expected tier | Model | Effort | Tokens | Elapsed |
| ---- | ------------- | ----- | ------ | ------ | ------- |
| The plan and its review | reasoning | Claude Opus 5.5; Claude Fable 5.1 (shared) | not reported | shared with #2 | shared with #2 |
| The decay review rounds | reasoning | Claude Fable 5.1 | not reported | 566,001 (177,489 + 228,558 + 159,954) | 34 min (14 + 12 + 8) |
| Writing the tests and the code | execution | Claude Opus 5.5 (one tier: limit recorded) | not reported | not reported | 40 min |
| Isolate, guardrails, docs, close-out | `—` | Claude Opus 5.5 | not reported | not reported | 15 min |
| **Total** | | | | 566,001 reported + not reported | about 89 min |
