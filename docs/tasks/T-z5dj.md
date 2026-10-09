# T-z5dj — row 31 of the plan, the calls of M2b of internal/git

Issue: [#157](https://github.com/pharzam/layup/issues/157), row 31 of [the tasks
of M2b](../plan/README.md#the-tasks-of-m2b), opened by `T-fdaq` (#152). Serves
`F-0003#43` through `REQ-003` and `NFR-001`. Base `8d2fa2e` (the merge of #171).
Author: Claude Opus 5.5 on Claude Code. Evidence:
[`runs/T-z5dj/`](../../runs/T-z5dj/).

## Plan and plan review

The plan (R12, comment 6080323845), its review (6080575151) and the author's
answer (6080575807) are comments on #157. The plan review (Claude Fable 5.1,
effort `xhigh`, on the Claude Code CLI with stream output, a fresh read-only
session in a clone at `8d2fa2e`, 12 min 27 s) gave `approve-with-conditions`:
Budget maximum 750 lines added plus removed over 12 files against `8d2fa2e`,
close-out inside; Cycle cap 1; no panel. Its four conditions are applied: one
error shape (`ErrNotACommit` in a `*FailedError`); `core.alternateRefsCommand`
in the list and a control that shows each live key; both Test cells; the
scratch directories removed on each return. The goal count is 3, final by O-187
of #152.

## What was done

1. **D1:** `CloneLocal`, `IsAncestor`, `DiffFile` and `DiffBinary` in
   `internal/git/git.go`, each with its input checked before `git` starts.
2. **D2:** `FetchSession` and `ErrNotACommit`: a scratch bare repository whose
   alternate is the session's object directory, `cat-file -t`, `update-ref`, and
   the fetch with an empty hooks directory; no `git` runs in the session's clone.
   `TestAHostileSessionRunsNothing` writes the 47 program-starting keys that
   `git help --config` of git 2.47.3 lists into the session's configuration,
   each with a marker named for it; a control shows, by plain reads in a copy,
   that 15 of them fire on this host, and `FetchSession` fires none.
3. **D3:** `packages.md` ("The calls of `M2b`": `ErrNotACommit`, the scratch
   directories, the directories of the commands, the refusal of a branch, the
   known limit of `--update-head-ok`); `traceability.md`; the Test cells of
   `REQ-003` and `NFR-001`, with a §13 line.

**Tests:** [`test-runs.md`](../../runs/T-z5dj/test-runs.md): red before the
calls (no compile); the control red three times until each key had its own
driver; red with a `FetchSession` that ran `git` in the session (the hostile
test found `core.fsmonitor` and `filter.y.clean`); green after.

**Known limits.** The list of keys is that of `git help --config` of 2.47.3; a
key that a later `git` adds is not in it. Thirty-two keys are asserted, not shown
live (the test's `notLive`), and `remote.<name>.vcs` names a helper of the
`PATH`, so it has no marker. A
`DST` that is the branch of `HEAD` of `dir` is refused by `git`.

**The rejected alternatives:** a fetch from the session's clone (its
`upload-pack` reads the session's configuration); `git bundle` (it runs `git` in
the session's clone); copying the objects by hand.

## Review rounds and verdict

The records, the Fixes replies and O-188 are comments on #157. Round 1
(`6958da4`, cycle 0): `material`, the list of keys was not each key of git's
documentation; fixed in `9e8e201` (45 keys, a mutation red for each call, the
lesson). Round 2 (`9e8e201`, cycle 1): two keys of `git help --config` still
missing; **O-188** (a) raised the cap to 2; fixed in `812aeb1` (47 keys).
Round 3 (`812aeb1`, cycle 2): `nothing material in scope`; notes 1 to 4 applied
in the close-out (the source of the list in `session.md`, the head checked
after the armed fetch, one helper for a full object ID, the Operator's host in
`packages.md`); note 5, revealed, went to #159. All rounds by Claude Fable 5.1.

Delivered: the five calls of `M2b` of `internal/git`. The review ended by decay
at cycle 2 of the cap that O-188 raised. Next: rows 33a and 37a use them.

## Resource record

Recorded, not budgeted (ADR-0007). UTC, 2026-10-09; reviewers' tokens are
`modelUsage`; the author's are not reported.

| Part | Expected tier | Model | Effort | Tokens | Elapsed |
| ---- | ------------- | ----- | ------ | ------ | ------- |
| The plan | reasoning | Claude Opus 5.5 | not reported | not reported | 11:51 to 11:53 |
| The plan review | reasoning | Claude Fable 5.1 | `xhigh` | 960,079 (USD 8.18) | 12 min 27 s |
| The code, test first | execution | Claude Opus 5.5 | not reported | not reported | 12:10 to 12:21 |
| Round 1 | reasoning | Claude Fable 5.1 | `xhigh` | 1,078,662 (USD 5.73) | 12 min 8 s |
| The fix of round 1 | execution | Claude Opus 5.5 | not reported | not reported | 12:34 to 12:39 |
| Round 2 | reasoning | Claude Fable 5.1 | `xhigh` | 941,598 (USD 6.13) | 13 min 8 s |
| O-188; the fix of round 2 | execution | Claude Opus 5.5 | not reported | not reported | 12:53 to 13:02 |
| Round 3 | reasoning | Claude Fable 5.1 | `xhigh` | 746,298 (USD 5.93) | 15 min 1 s |
| The close-out | execution | Claude Opus 5.5 | not reported | not reported | 13:25 to 13:30 |
