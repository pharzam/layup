# T-r7zg — record the Armature pin and start setup-check

Issue: [#2](https://github.com/pharzam/layup/issues/2) (parent [#1](https://github.com/pharzam/layup/issues/1)). PR: [#13](https://github.com/pharzam/layup/pull/13).

## Verdict

Delivered. `docs/setup/armature.pin` records Armature `a959655` and its tree
`8ffb250a`; check `pin` in `docs/setup/setup-check.sh` fails on an absent file, a
key not present exactly once, a commit that is not 40 hex characters, a shallow
clone, and a tree that differs from the one root commit. With no argument the
script also runs the four kit linters. Self-test: 7 cases pass under `sh` and
`dash`; round reviewers ran 17 mutations and each turned the suite red.
Round 1 found two material defects (a duplicate key passed; `F-0001` cited before
it exists), fixed in cycle 1; round 2 found nothing material. Non-material
findings of round 2 stay on the issue. The lesson "POSIX `sh` has no local
variables, so a check function can overwrite the caller's loop variable" goes to
`guardrails.md` §2 with `T-7ndb`.

## Resource record

Recorded, not budgeted (ADR-0007). Tokens of the author session are `not
reported` by the harness; subagent tokens are from the harness task report. The
plan review ran once for all child plans, so its figure is shared.

| Part | Expected tier | Model | Effort | Tokens | Elapsed |
| ---- | ------------- | ----- | ------ | ------ | ------- |
| The plan and its review | reasoning | Claude Opus 5.5 (plan); Claude Fable 5.1 (review, shared by #1–#12) | not reported | not reported (author); 195,984 (shared review run) | 12 min (review run, shared) |
| The decay review rounds | reasoning | Claude Fable 5.1 | not reported | 288,187 (140,377 + 147,810) | 18 min (10 + 8) |
| Writing the tests and the code | execution | Claude Opus 5.5 (no execution-tier model was used: one tier, recorded as a limit) | not reported | not reported | 35 min |
| Isolate, guardrails, docs, close-out | `—` | Claude Opus 5.5 | not reported | not reported | 15 min |
| **Total** | | | | 484,171 reported + not reported | about 80 min |
