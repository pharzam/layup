# T-hpid — row 41b of the plan, the command layup spec check --facts

Issue: [#214](https://github.com/pharzam/layup/issues/214), row 41b of [the tasks
of M2c](../plan/README.md#the-tasks-of-m2c), opened by `T-2r8u` (#211). Serves
`F-0003#51` through REQ-012. Base `bec9093` (the merge of #239, row 42b); the
budget is read against it. Author: Claude Opus 5.5 on Claude Code. Evidence:
[`runs/T-hpid/`](../../runs/T-hpid/).

## Plan and plan review

The plan (R12), its review and the author's answer are comments on #214. The
plan review (Claude Fable 5.1, effort `xhigh`, a fresh read-only session in a
clone at `bec9093`, 5 min 38 s) gave `approve-with-conditions`: Budget maximum
600 lines added plus removed over 14 files, close-out inside; Cycle cap 1; no
panel. By O-199 (comment 6106919969) the files rose to 16 once and the cap to
2 (comment 6106920131). Its two conditions are applied: each exit-2 case asserts its reason, so
the red is `unknown command "spec"`; a table that cannot be written has a unit
test. The goal count is 1, final by O-196.

## What was done

1. **`internal/cli`:** the row `spec check --facts SPANS --psb FILE` of the
   command table, and `specCheck`: both files read, a file missing, unreadable
   or not valid UTF-8 is exit 2 naming it; a span table that its block refuses
   is exit 2, the reason prefixed with the file; else the table of
   `spec.Check`, and the exit code of its column `result` (0 or 1); a table
   that cannot be written is exit 2.
2. **`internal/spec`:** `WriteTSV`, the writer of the table by `CheckSchema`,
   as `internal/cli` may not import `internal/tsv` (`TestPackageRules` refused
   the first form).
3. **The documents:** `intake.md` (the decided-here on a table that cannot be
   written), `README.md` (Progress now names `layup spec check` among the
   commands that print none; the exit code 2 row names it), `traceability.md`,
   the Task and Test cells of REQ-012 and a §13 line.

**Tests:** [`test-runs.md`](../../runs/T-hpid/test-runs.md): red before the code;
7 of 7 mutations caught; green at the three levels.

## Review rounds

The records are comments on #214. Round 1 (Claude Fable 5.1, effort `xhigh`, a
fresh read-only session in a clone at `0aab2d0`, cycle 0): `material`, two
findings: trailing spaces in the pasted red of the test record (the check had
run before `git add`, on an untracked file), and the NFR-005 cell of PRD §12
not naming the repeat test. The fix `2de3ec0` (cycle 1) took them out, named
the command in `README.md` and the onboarding (note 5), and wrote two lessons of
§2. Round 2 (the same reviewer, at `2de3ec0`, cycle 1, the cap): `not
mergeable, findings recorded`, edited to `material` after O-199: the files at
landing passed 14, and one clause of the record, from a wrong note of round 1,
was false. The fix `cae30ac` (cycle 2) took the clause out and named the
write-failure test for NFR-004. Round 3 (the same reviewer, at `cae30ac`, cycle
2, the cap): `nothing material in scope`; its notes 2 to 4, with note 5 of
round 1, are #245.

## Verdict

Delivered: the command `layup spec check --facts SPANS --psb FILE`, the table of
the span check, exit 0, 1 or 2, and the writer of its table in
`internal/spec`. The review ended by decay at cycle 2 of cap 2 (O-199). The diff
against `bec9093`, the branch's base, is inside 600 lines over 16 files (O-199).
CI runs the three levels on the landing head. Next: row 44.

## Resource record

Recorded, not budgeted (ADR-0007). UTC; reviewers' tokens are `modelUsage`;
the author's are not reported.

| Part | Expected tier | Model | Effort | Tokens | Elapsed |
| ---- | ------------- | ----- | ------ | ------ | ------- |
| The plan | reasoning | Claude Opus 5.5 | not reported | not reported | 2026-10-11 06:37 to 06:37 (the claim to the post) |
| The plan review | reasoning | Claude Fable 5.1 | `xhigh` | 1,116,563 (USD 5.62) | 5 min 38 s |
| The code, test first | execution | Claude Opus 5.5 | not reported | not reported | 06:43 to 06:52 |
| Round 1 | reasoning | Claude Fable 5.1 | `xhigh` | 924,765 (USD 4.64) | 9 min 26 s |
| The fix of round 1 | execution | Claude Opus 5.5 | not reported | not reported | 07:02 to 07:04 |
| Round 2 | reasoning | Claude Fable 5.1 | `xhigh` | 1,487,356 (USD 5.40) | 11 min 0 s |
| The question to the Operator | — | the Operator | — | — | 07:21 to 08:00 (O-199) |
| The fix of round 2 | execution | Claude Opus 5.5 | not reported | not reported | 08:00 to 08:02 |
| Round 3 | reasoning | Claude Fable 5.1 | `xhigh` | 1,128,918 (USD 5.31) | 10 min 29 s |
| The close-out | execution | Claude Opus 5.5 | not reported | not reported | 08:14 to 08:15 |
