# T-1u2v — row 43b of the plan, the review of meaning

Issue: [#218](https://github.com/pharzam/layup/issues/218), row 43b of [the tasks
of M2c](../plan/README.md#the-tasks-of-m2c), opened by `T-2r8u` (#211). Serves
`F-0003#41` through REQ-001. Base `bec9093` (the merge of #239, row 42b); the
budget is read against it, the base of this branch's own work (as O-198 read
#213's). Author: Claude Opus 5.5 on Claude Code. Evidence:
[`runs/T-1u2v/`](../../runs/T-1u2v/).

## Plan and plan review

The plan (R12), its review and the author's answer are comments on #218. The
plan review (Claude Fable 5.1, effort `xhigh`, a fresh read-only session in a
clone at `bec9093`, 8 min 37 s) gave `approve-with-conditions`: Budget maximum
1,200 lines added plus removed over 22 files, close-out inside; Cycle cap 1; no
panel. Its two conditions are applied: the world of `run_integration_test.go`
gained a harness, its models and a route, so the restart tests keep their
assertions; the lease after a `fail` keeps its one home, the known limit of
`run.md`. Note 7 is #242, a question for the Operator. The goal count is 1,
final by O-196.

## What was done

1. **The quote check:** `quoteAt` finds a quote once in the problem statement
   and gives its byte offsets; an empty quote, a tab, a line break, none or
   more than one occurrence refuse it.
2. **The file of the session:** `sessionSchema` derives its form from the block
   `meaning` by dropping the columns that code adds (`kind`, `quote`,
   `question`); a kept key stays a key, for the later rows' blocks.
3. **`Keep` of the step:** the quotes found, `M-NNN` in order, the offsets and
   the session ID, as `intake/meaning.tsv`; a refused quote refuses the file.
4. **The step `meaning`** after `rule-gaps`: "done before" with the record; else
   the pair of `systems-architect` (none: `fail`, `pair`), one session of a
   phase with its fixed prompt, and `done` with the count of the gaps. The
   enum of `run-steps` and its Go schema gained `meaning`.
5. **The tests' worlds:** the world of `run_integration_test.go` makes its bare
   target with `HEAD` on `main`, takes a posted comment, and has a scripted
   harness with a route; the e2e world of M2a gained the route and a meaning
   branch in its script.
6. **The documents:** `intake.md` (the derived form of a session's file; the
   lease after a `fail`, by `run.md`), `run.md` (the block `run-steps`; the
   known limit extended to the rows after `lease`), the Job cell of
   `internal/run`, `traceability.md`, the Task and Test cells of REQ-001 and a
   §13 line.

**Tests:** [`test-runs.md`](../../runs/T-1u2v/test-runs.md): red before the code;
12 of 12 mutations caught; green at the three levels.

**A known limit met:** the quote of one em dash, which reads back as the empty
value (row 40, `T-529i`), is refused as an empty quote.
