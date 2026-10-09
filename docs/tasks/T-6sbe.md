# T-6sbe — row 33b of the plan, the checks before the start of a session

Issue: [#160](https://github.com/pharzam/layup/issues/160), row 33b of [the tasks
of M2b](../plan/README.md#the-tasks-of-m2b), opened by `T-fdaq` (#152). Serves
`F-0003#52` through `REQ-013` and `REQ-003`. Base `cee0e4d` (the merge of #184,
row 33a). Author: Claude Opus 5.5 on Claude Code. Evidence:
[`runs/T-6sbe/`](../../runs/T-6sbe/).

## Plan and plan review

The plan (R12, comment 6082571206), its review (6085494971) and the author's
answer (6085495672) are comments on #160. The plan review (Claude Fable 5.1,
effort `xhigh`, a fresh read-only session at row 33a's head `0e8f15b`, 8 min 2
s) gave `approve-with-conditions`: Budget maximum 650 lines added plus removed
over 13 files, close-out inside; Cycle cap 1; no panel. Its two conditions are
applied: the rule of a `packed-refs` line (its ref field, its first field of 40
characters), written in `session.md`; the rules of the head at unit through a
stand-in reader, the two forms of the real `git` at integration. The goal count
is 3, final by O-187 of #152.

## What was done

1. **D1:** `Refusal` (the reason and its value, the form of the event
   `refused`), `NewID`, `CheckContext` (the estimate then the size, decided
   here), `CheckPromptSize`.
2. **D2:** `CheckRuleFiles`: from the session directory up to `/`; the policy
   paths that exist.
3. **D3:** `HeadOf` and `headOf`: the loose ref, else the line of
   `packed-refs`, each other form `Refusal{branch}`.
4. **D4:** `session.md`, `packages.md`, `traceability.md`, the Test cells of
   `REQ-013` and `REQ-003`, a §13 line.

**Tests:** [`test-runs.md`](../../runs/T-6sbe/test-runs.md): red before the
functions; a mutation of each of the fifteen rules, each caught (two cases now
assert the value that names their rule); green after.

**The rejected alternatives:** `git rev-parse` in `repo/`; one type per reason
of a refusal.
