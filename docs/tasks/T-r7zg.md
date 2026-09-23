# T-r7zg — record the Armature pin and start setup-check

Issue: [#2](https://github.com/pharzam/layup/issues/2) (parent [#1](https://github.com/pharzam/layup/issues/1)). PR: [#13](https://github.com/pharzam/layup/pull/13).

## Verdict

Delivered: pin `a959655` / tree `8ffb250a` and check `pin` (7 fixtures; 17 reviewer
mutations each went red). Round 1: 2 material findings, fixed in cycle 1. Round 2:
nothing material. Lesson for `guardrails.md` §2 (lands with `T-7ndb`): POSIX `sh`
has no local variables, so a check function can overwrite its caller's loop variable.

## Resource record

Recorded, not budgeted (ADR-0007). Author tokens: `not reported`. Subagent tokens:
harness task report; the plan review ran once for #1–#12, so its figure is shared.

| Part | Expected tier | Model | Effort | Tokens | Elapsed |
| ---- | ------------- | ----- | ------ | ------ | ------- |
| The plan and its review | reasoning | Claude Opus 5.5 (plan); Claude Fable 5.1 (review, shared by #1–#12) | not reported | not reported (author); 195,984 (shared review run) | 12 min (review run, shared) |
| The decay review rounds | reasoning | Claude Fable 5.1 | not reported | 288,187 (140,377 + 147,810) | 18 min (10 + 8) |
| Writing the tests and the code | execution | Claude Opus 5.5 (no execution-tier model was used: one tier, recorded as a limit) | not reported | not reported | 35 min |
| Isolate, guardrails, docs, close-out | `—` | Claude Opus 5.5 | not reported | not reported | 15 min |
| **Total** | | | | 484,171 reported + not reported | about 80 min |
