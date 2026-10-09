# T-vxdg — row 33a of the plan, the directory and the environment of a session

Issue: [#159](https://github.com/pharzam/layup/issues/159), row 33a of [the tasks
of M2b](../plan/README.md#the-tasks-of-m2b), opened by `T-fdaq` (#152). Serves
`F-0003#52` through `REQ-013` and `REQ-003`. Base `d5d0f7f` (the merge of #181;
the plan was reviewed at `5fbe7e5`). Author: Claude Opus 5.5 on Claude Code.
Evidence: [`runs/T-vxdg/`](../../runs/T-vxdg/).

## Plan and plan review

The plan (R12, comment 6082093535), its review (6082403220) and the author's
answer (6082403842) are comments on #159. The plan review (Claude Fable 5.1,
effort `xhigh`, a fresh read-only session in a clone at `5fbe7e5`, 8 min 28 s)
gave `approve-with-conditions`: Budget maximum 800 lines added plus removed
over 16 files, close-out inside; Cycle cap 1; no panel. Its five conditions are
applied: `Sweep` tells a stopped run's directory from another target's by a
file `target`; `home/` holds `.gitconfig` and a `file:` credential and nothing
else; the work tree of `repo/` is the base; `Remove` is left to row 36b; the
allowed reader `environ` builds the whole list. The goal count is 2, final by
O-187 of #152.

## What was done

1. **D1:** `environ` (the named list, the host's `PATH` the one variable it
   reads; the credential's variable less one final line feed; the fixed
   variables) and `Environ`, which reads the credential file.
2. **D2:** `Make` (the file `target`, `repo/` by `CloneLocal` and
   `SwitchCreate`, `home/` with `.gitconfig` and a `file:` credential of mode
   0600, `tmp/`, `result/`, `prompt.md`; an existing ID refused) and `Sweep`
   (the directories of its own target, and those with no file `target`).
3. **D3:** `TestInputRule` allows `environ` of `internal/session`, with a case
   of the checker's own test; `session.md` (the file `target`, `Sweep`, what
   `home/` holds, the clone's remote removed by `CloneLocal`); `packages.md`;
   `traceability.md`; the Test cells of `REQ-013` and `REQ-003`, a §13 line.

**Tests:** [`test-runs.md`](../../runs/T-vxdg/test-runs.md): red before the
package; a mutation of each rule, each caught; green after.

**The rejected alternatives:** importing `internal/route`; the credential read
by `layup run`; keeping a stopped run's directory to resume it; one run per
host as the reading of `Sweep`.

## Review rounds

The record is a comment on #159. Round 1 (Claude Fable 5.1, effort `xhigh`, a
fresh read-only session in a clone at `0e8f15b`, cycle 0): `nothing material in
scope`, seven notes. Applied in the close-out: note 1 (the record says its
mutation lines are cut); note 2 (a commit after the base in the fixture, so "at
the base" can fail on its own rule); note 3 (the clause of `CloneLocal` beside
its noun); note 4 (the known limit of the window between the first two writes
of `Make`); note 5 (`Sweep` skips an entry that is not a directory, with a
case); note 6 (the fourth test in the Test cell). Declined: note 7, revealed
(the word "empty" for `home/` in `architecture.md`, ADR-0015 and the glossary
means "nothing of the host", which `architecture.md` itself shows, and those
documents are off this row's path).

## Verdict

Delivered: the session's files in `internal/session`: `Make`, `Sweep`,
`Environ`. The review ended by decay at cycle 0 of cap 1. The diff against
`d5d0f7f`, the branch's base, is inside 800 lines over 16 files. Next: rows 33b
and 34a.

## Resource record

Recorded, not budgeted (ADR-0007). UTC, 2026-10-09; reviewers' tokens are
`modelUsage`; the author's are not reported.

| Part | Expected tier | Model | Effort | Tokens | Elapsed |
| ---- | ------------- | ----- | ------ | ------ | ------- |
| The plan | reasoning | Claude Opus 5.5 | not reported | not reported | 13:36 to 13:40 |
| The plan review | reasoning | Claude Fable 5.1 | `xhigh` | 1,262,363 (USD 5.94) | 8 min 28 s |
| The code, test first | execution | Claude Opus 5.5 | not reported | not reported | 13:56 to 14:02 |
| Round 1 | reasoning | Claude Fable 5.1 | `xhigh` | 870,263 (USD 4.83) | 10 min 30 s |
| The close-out, with notes 1 to 6 | execution | Claude Opus 5.5 | not reported | not reported | 14:19 to 14:23 |
