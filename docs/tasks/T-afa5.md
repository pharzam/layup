# T-afa5 — require the LAYUP CI jobs on main and record the setting

Issue: [#12](https://github.com/pharzam/layup/issues/12) (parent [#1](https://github.com/pharzam/layup/issues/1)). PR: [#26](https://github.com/pharzam/layup/pull/26).

## Verdict

Delivered: `main` requires all nine CI jobs (app 15368, `strict`,
`enforce_admins`, a PR with 0 approvals); the body is
`docs/setup/branch-protection.json`, applied at 16:35:40 UTC after round 1 checked
the names; the read-back equals the file. Check `protection` keeps the contexts
equal to the job names. Rounds 1 and 2: material (the CI README still offered the
kit's six-check body for this repository), fixed. Round 3: nothing material.
Revealed: added to #24.

## Resource record

Recorded, not budgeted (ADR-0007). Author tokens: `not reported`.

| Part | Expected tier | Model | Effort | Tokens | Elapsed |
| ---- | ------------- | ----- | ------ | ------ | ------- |
| The plan and its review | reasoning | Claude Opus 5.5; Claude Fable 5.1 (shared) | not reported | shared with #2 | shared with #2 |
| The decay review rounds | reasoning | Claude Fable 5.1 | not reported | 376,908 (96,017 + 121,460 + 159,431) | 25 min (7 + 9 + 9) |
| Writing the tests and the code | execution | Claude Opus 5.5 (one tier: limit recorded) | not reported | not reported | 30 min |
| Isolate, guardrails, docs, close-out | `—` | Claude Opus 5.5 | not reported | not reported | 10 min |
| **Total** | | | | 376,908 reported + not reported | about 65 min |
