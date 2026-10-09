# T-bpxg — row 34b of the plan, the end of a session

Issue: [#162](https://github.com/pharzam/layup/issues/162), row 34b of [the tasks
of M2b](../plan/README.md#the-tasks-of-m2b), opened by `T-fdaq` (#152). Serves
`F-0003#52` and `F-0003#50` through `REQ-013`, `REQ-005` and `REQ-011`. Base
`53f1bef` (the merge of #188, row 34a); the budget is read against it. Author:
Claude Opus 5.5 on Claude Code. Evidence: [`runs/T-bpxg/`](../../runs/T-bpxg/).

## Plan and plan review

The plan (R12), its review and the author's answer are comments on #162. The
plan review (Claude Fable 5.1, effort `xhigh`, a fresh read-only session at
row 34a's head `179646f`, 8 min 57 s) gave `approve-with-conditions`: Budget
maximum 1,000 lines added plus removed over 18 files, close-out inside; Cycle
cap 1; no panel. Its four conditions are applied: the token statuses are
total, with a reading of `partial` and of a report with no model; `token` at
most once, in the rules of `records.md`, with cases of zero rows and a wrong
header; the order of the classes, the refusal of a result file that is not a
regular file, and the sorted models, each **decided here** in `session.md`;
the recorded run names its models and version, and its stream was read before
anything of it was kept. The goal count is 2, final by O-187 of #152.

## What was done

1. **D1:** `Class`: the class of the end, in the order `start`; `wall` or
   `output`; `crash`; `no-result`; `done`. Another error of `Process` is no
   class.
2. **D2:** `ResultOf` (no link followed, a regular file of at most 1 MiB, by its
   block's reader), `ProbeResultSchema` and `ReadProbeResult`, the block
   `probe-result` moved to `built`.
3. **D3:** `Usage` (the fields of `ledger.Usage`, which a test converts),
   `UsageOf`: `claude-result` and `none`, on two recorded `result` events of
   Claude Code 2.1.295.
4. **D4:** `session.md`, `records.md`, `packages.md`, `traceability.md`, the
   Test cells of `REQ-013`, `REQ-005` and `REQ-011` and a §13 line.

**Tests:** [`test-runs.md`](../../runs/T-bpxg/test-runs.md): the recorded
events and their sums by hand; red before the functions; a mutation of each of
the twenty-six rules, each caught, one added by the fix of round 1 and one
by the close-out; green after.

**The rejected alternatives:** a class from the exit alone; a result file read
through a link; a float for the cost.

**For row 36b** (note 3 of round 2): `ResultOf` refuses a link as the last
part of the path only (`O_NOFOLLOW`). A session that replaces `result/` by
a link to a directory of the host reaches a regular file there; the file still
has to pass its block, so it is one the session could have written. Row 36b,
which commits the file byte for byte, reads this limit.

## Review rounds

The records are comments on #162. Round 1 (Claude Fable 5.1, effort `xhigh`,
a fresh read-only session in a clone at `be96597`, cycle 0): `material`.
Finding 1: a token field of `null` was summed as 0. Finding 2: the sentence on
an error of `layup run`'s own was narrower than `Class`. Both fixed in
`0b706d5`, with notes 3 to 5 (row 10 and row 12 of Input states, the record
of the integration red). Round 2 (the same model and effort, a fresh session
at `0b706d5`, cycle 1, the cap): `nothing material in scope`, three notes.
Applied in the close-out: note 1 (the cell of `claude-result` points to the
decided paragraph) and note 2 (`none` reads nothing, with a case). Note 3 is
the paragraph above, for row 36b.

## Verdict

Delivered: the end of a session in `internal/session`: `Class`,
`ResultOf`, `ReadProbeResult` with the block `probe-result` built, `Usage`
and `UsageOf`. The review ended by decay at cycle 1 of cap 1. The diff against
`53f1bef`, the branch's base, is inside 1,000 lines over 18 files. Next: row
36b, after row 36a.

## Resource record

Recorded, not budgeted (ADR-0007). UTC, 2026-10-09; reviewers' tokens are
`modelUsage`; the author's are not reported.

| Part | Expected tier | Model | Effort | Tokens | Elapsed |
| ---- | ------------- | ----- | ------ | ------ | ------- |
| The plan | reasoning | Claude Opus 5.5 | not reported | not reported | to 17:43 |
| The plan review | reasoning | Claude Fable 5.1 | `xhigh` | 1,286,513 (USD 5.86) | 8 min 57 s |
| The recorded usage report (test data) | — | Claude Fable 5.1 and Claude Opus 5.5 | default | 52,666 (USD 0.27) | 10 s |
| The code, test first | execution | Claude Opus 5.5 | not reported | not reported | 18:13 to 18:23 |
| Round 1 | reasoning | Claude Fable 5.1 | `xhigh` | 1,086,767 (USD 4.60) | 9 min 13 s |
| The fix of round 1 | execution | Claude Opus 5.5 | not reported | not reported | 18:33 to 18:36 |
| Round 2 | reasoning | Claude Fable 5.1 | `xhigh` | 804,924 (USD 4.66) | 8 min 47 s |
| The close-out, with notes 1 and 2 | execution | Claude Opus 5.5 | not reported | not reported | 18:45 to 18:48 |
