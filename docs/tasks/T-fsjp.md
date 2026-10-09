# T-fsjp — row 36b of the plan, the result and the end of a task session

Issue: [#165](https://github.com/pharzam/layup/issues/165), row 36b of [the tasks
of M2b](../plan/README.md#the-tasks-of-m2b), opened by `T-fdaq` (#152). Serves
`F-0003#45` and `F-0003#50` through `REQ-005`, `REQ-011` and `NFR-001`. Base
`c15e2c8` (the merge of #191, row 36a); the budget is read against it. Author:
Claude Opus 5.5 on Claude Code. Evidence: [`runs/T-fsjp/`](../../runs/T-fsjp/).

## Plan and plan review

The plan (R12), the Operator's decision O-189, the plan review and the author's
answer are comments on #165. The plan asked the Operator whether `M2b` or `M2c`
gives the forge call `Comment`; **O-189** (a): this row specifies and builds it.
The plan review (Claude Fable 5.1, effort `xhigh`, a fresh read-only session
at a local merge of row 36a into `main`, 8 min 7 s) gave
`approve-with-conditions`: Budget maximum 1,200 lines added plus removed over
21 files, close-out inside; Cycle cap 1; no panel. Its three conditions are
applied: no command of `layup` runs in `repo/` (the head read from its files,
fetched into the run's clone with `FetchSession`, the artifacts read there;
the clause moved from row 37a); each value this row sets is marked **decided
here** in `session.md`; `Pair` carries the row's `Billing` and `Usage`. The
goal count is 2, final by O-187 of #152.

## What was done

1. **The forge call** `Comment(issue, body)`: `forge.md`, "The calls of M2b";
   the interface, the GitHub adapter, and the two stand-ins of the tests.
2. **D1:** `openAttempt`, the open attempt of a session's result.
3. **D2:** `endTask`, the end of a task session: the result file from one read
   (`session.ResultFile`), the class, the head and the artifacts, the open
   attempt, and one records commit of the result, the events and the
   telemetry row.
4. **D3:** the comment on the control issue, and the removal of the directory
   after the call once the process started.
5. **D4:** `session.md`, `forge.md`, `packages.md`, the plan's rows 36b and
   37a and `runs/T-fdaq/parts.tsv` (the fetch moved), `traceability.md`, the
   Test cells of `REQ-005`, `REQ-011` and `NFR-001` and a §13 line.

**Tests:** [`test-runs.md`](../../runs/T-fsjp/test-runs.md): red before the
forge call and before the open attempt; a mutation of each of the twenty-one
rules of the end, each caught, three of them added by the fix of round 1 and
one by the close-out; green after.

**A deviation:** the integration tests of the end were written after
`end.go`, so they were not red before the code; their red is the mutations.

**A second deviation** (note 8 of round 1): the end's comment is tested
against the stand-in forge, and the adapter's `Comment` against `httptest`;
no one test runs the whole path from `Sessions` to HTTP.

**For row 37a:** the head is fetched into the run's clone at the end of a
`done` session of an open attempt, as `refs/layup/sessions/<session>`; the check of the base
(`IsAncestor`) and the rule-path check read it there.

**The rejected alternatives:** reading the result file twice; committing the
result of a closed attempt.

## Review rounds

The records are comments on #165. Round 1 (Claude Fable 5.1, effort `xhigh`,
a fresh read-only session in a clone at `a7d4ff4`, cycle 0): `material`.
Finding 1: the comment was posted inside the end, before the push hook could
record a refusal of the result. Fixed in `0738278`: `TaskSession` posts it
after the call, from the events, with notes 3 to 6 and two more cases; note 9
declined (row 33b decided that a `branch` refusal has no value of the event).
Round 2 (the same model and effort, a fresh session at `0738278`, cycle 1, the
cap): `nothing material in scope`, seven notes, all applied in the close-out:
no comment after a lost lease; the note of the second deviation above; the two
comment tests in the Test cell of `REQ-005`; the comment of `endTask`; "of an
open attempt" in the paragraph for row 37a; an error of `FetchSession` other
than `ErrNotACommit` is the run's own; a case for a `stdout` that cannot be
read.

## Verdict

Delivered: the forge call `Comment` (O-189), and the end of a task session in
`internal/run`: the result committed byte for byte with its events and its
telemetry row in one records commit, the head and the artifacts in the run's
clone, the open attempt, the comment after every record, and the removal of
the directory. The review ended by decay at cycle 1 of cap 1. The diff against
`c15e2c8`, the branch's base, is inside 1,200 lines over 21 files. Next: row
37a.

## Resource record

Recorded, not budgeted (ADR-0007). UTC, 2026-10-09; reviewers' tokens are
`modelUsage`; the author's are not reported.

| Part | Expected tier | Model | Effort | Tokens | Elapsed |
| ---- | ------------- | ----- | ------ | ------ | ------- |
| The plan | reasoning | Claude Opus 5.5 | not reported | not reported | to 18:52 |
| The plan review | reasoning | Claude Fable 5.1 | `xhigh` | 1,021,510 (USD 6.66) | 8 min 7 s |
| The code, test first | execution | Claude Opus 5.5 | not reported | not reported | 19:41 to 19:54 |
| Round 1 | reasoning | Claude Fable 5.1 | `xhigh` | 1,170,768 (USD 5.57) | 9 min 47 s |
| The fix of round 1 | execution | Claude Opus 5.5 | not reported | not reported | 20:05 to 20:10 |
| Round 2 | reasoning | Claude Fable 5.1 | `xhigh` | 856,904 (USD 5.52) | 10 min 17 s |
| The close-out, with the notes of round 2 | execution | Claude Opus 5.5 | not reported | not reported | 20:20 to 20:23 |
