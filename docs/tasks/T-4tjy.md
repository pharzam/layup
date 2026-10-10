# T-4tjy — row 39a of the plan, the review of the release of M2b

Issue: [#169](https://github.com/pharzam/layup/issues/169), row 39a of [the
tasks of M2b](../plan/README.md#the-tasks-of-m2b), opened by `T-fdaq` (#152).
Serves `F-0003#52` through `REQ-015` and `REQ-017`; takes
[#148](https://github.com/pharzam/layup/issues/148) (O-187 of #152). Base
`9490de9` (the merge of #197, row 38); the budget is read against it. Author:
Claude Opus 5.5 on Claude Code. Evidence: [`runs/T-4tjy/`](../../runs/T-4tjy/).

## Plan and plan review

The plan (R12), its review and the author's answer are comments on #169. The
plan review (Claude Fable 5.1, effort `xhigh`, a fresh read-only session in a
clone at `9490de9`, 8 min 30 s) gave `approve-with-conditions`: Budget maximum
750 lines added plus removed over 13 files, close-out inside; Cycle cap 1; no
panel. Its five conditions are applied: the fourth allowance of check (4),
`remote` in `CloneLocal`, is written in `session.md`; the release of `M2b` is
defined in the paragraph of row 39a, and the phase-1 sentence is unchanged; the
cites and the rows of rule 5 are corrected; a finding of the release review is
applied on the task's path, or opens an issue; each red run is its own commit
of a scratch clone. The goal count is 2, final by O-187.

## What was done

1. **The release of `M2b`** (D1): the code of `main` that the demo runs, at
   `9490de9`, the merge of the last build row; it holds phase 1 and `M2a`, so
   the review settles #148.
2. **The check, adapted in place** (D3):
   [`release-check.sh`](../../runs/T-efmy/release-check.sh). Check (3) fails on
   a call that starts a program out of a fixed list of a file and its program;
   check (4) allows a verb that reaches a remote only in its own function
   (`fetch` in `Fetch` and `FetchSession`, `push` in `Push`, `remote` in
   `CloneLocal`) and fails on one in no function; (2) names the packages that
   import a package of rule 5; (6) lists the calls of the forge adapter. Its
   output at `9490de9`, exit 0:
   [`release-check.txt`](../../runs/T-4tjy/release-check.txt). The output of
   phase 1 stays `runs/T-efmy/release-check.txt`.
3. **The release review** (D4, the uat):
   [`release-review.md`](../../runs/T-4tjy/release-review.md), posted on #169:
   a fresh session of Claude Fable 5.1, a model whose sessions wrote no commit
   of the release, records `holds` for `REQ-015` and for `REQ-017`, with
   seventeen findings; none is a defect. Its notes 9 to 16 are off the path of
   this task and need no change: each records a path that a stricter reading
   would catch, or a known limit (L-A1, the program of a register row, the
   raw system call). Note 17 is on the path: its first half is condition 1,
   applied; its second half (check (6) prints `path` for the token call, whose
   endpoint is the line above) needs no change, as the reviewer says. The
   reading of acceptance criterion 2: a finding that its reviewer closes with
   "no change" names no defect, so it is neither fixed nor an issue.
4. **The documents** (D5): `session.md` (the allowance of check (4) and the
   failing list of (3)); the plan (the release of `M2b`); the §12 cells of
   `REQ-015` and `REQ-017` and a §13 line of `PRD-0001`; a uat row of
   `traceability.md`.

**Tests:** [`test-runs.md`](../../runs/T-4tjy/test-runs.md): the check of phase
1 red at the base; a first draft that read no file, caught by its run; five
mutations, each on its own commit, each caught by its line of (3) or (4); green
at `9490de9`.

**The rejected alternatives:** a new script under `runs/T-4tjy/` (the
specification adapts the one script); a review of `M2a` on its own (O-187 gives
#148 to this review).
