# T-vbwc — remove the kit's own history and backlog

Issue: [#3](https://github.com/pharzam/layup/issues/3) (parent [#1](https://github.com/pharzam/layup/issues/1)). PR: [#14](https://github.com/pharzam/layup/pull/14).

## Verdict

Delivered: kit step 4 done (31 files, the kit log entries, 6 backlog lines and the
pivot note); check `kit-history` with 5 fixtures; `--only` for fixture runs.
Rounds 1 and 2 found material defects (untested no-argument path, degenerate
orphan fixture, pointers to deleted records, a prose link that masked the ADR
index check); fixed in cycles 1 and 2. Round 3: nothing material. Revealed: #15.

## Resource record

Recorded, not budgeted (ADR-0007). Author tokens: `not reported`.

| Part | Expected tier | Model | Effort | Tokens | Elapsed |
| ---- | ------------- | ----- | ------ | ------ | ------- |
| The plan and its review | reasoning | Claude Opus 5.5; Claude Fable 5.1 (review shared by #1–#12) | not reported | not reported; shared | shared with #2 |
| The decay review rounds | reasoning | Claude Fable 5.1 | not reported | 424,770 (115,820 + 163,349 + 145,601) | 29 min (9 + 10 + 10) |
| Writing the tests and the code | execution | Claude Opus 5.5 (one tier: limit recorded) | not reported | not reported | 30 min |
| Isolate, guardrails, docs, close-out | `—` | Claude Opus 5.5 | not reported | not reported | 10 min |
| **Total** | | | | 424,770 reported + not reported | about 70 min |
