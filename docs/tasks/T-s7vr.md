# T-s7vr — the restart of a target whose harness register changed after Start

Issue: [#199](https://github.com/pharzam/layup/issues/199), a defect that the
demo of `M2b` revealed (`T-x7cs`, #170); taken first by the Operator's
directive after `M2b` (comment 6099589342 of #199). Serves `NFR-001` and
`F-0003#52`. Base `95dccce` (the merge of #200); the budget is read against
it. Author: Claude Opus 5.5 on Claude Code. Evidence:
[`runs/T-s7vr/`](../../runs/T-s7vr/).

## Plan and plan review

The plan (R12), its review and the author's answer are comments on #199. The
plan review (Claude Fable 5.1, effort `xhigh`, a fresh read-only session in a
clone at `95dccce`, 9 min 1 s) gave `approve-with-conditions`: Budget maximum
400 lines added plus removed over 12 files, close-out inside; Cycle cap 1; no
panel. Its three conditions are applied: the row of `run.md` decides only what
is new and points to `session.md` for the values a session runs under; the
block `start` of `records.md` says "at Start" and holds the reader's rule; the
rejected rewrite cites the record of the Start, not the events table. The goal
count is 1.

## What was done

1. **The reader:** `records.ReadStartAsWritten` reads `start.tsv` by the
   harnesses that its rows `harness.<id>.cap` name, in their order, each of the
   form `<word>`, and checks the rows with `CheckStart`; `ReadStart`, whose one
   caller was the restart, goes, and its cases move to the new reader's test.
2. **The restart:** `cloneStep` reads `start.tsv` so and keeps the Start's
   harnesses; `harnessIDs` gives them in a restart, so `writeStart` keeps them
   when the step `phase` writes `start.tsv` again.
3. **The documents:** `run.md` (Input states, a register that gained or lost a
   row after Start); `records.md` (the block `start`, "at Start", and the
   reader's rule); the §2 pitfall of #199 rewritten as the general trap; the
   open-issues row of #199 in the plan; `traceability.md`.

**Tests:** [`test-runs.md`](../../runs/T-s7vr/test-runs.md): red before the
code, a mutation of each of the four rules, each caught, green after.

**The rejected alternatives:** rewriting `start.tsv` with the register of now
(the record of the Start would say what the register says now); refusing a
changed register with a clearer message (the demo of `M2b` needed the change).

## Review rounds

The record is a comment on #199. Round 1 (Claude Fable 5.1, effort `xhigh`, a
fresh read-only session in a clone at `7c34fcb`, cycle 0, the cap): `nothing
material in scope`, four notes. Applied in the close-out: note 1 (the row of
`run.md` links the section of the records of Start); note 2 (the row points to
`session.md` for the values a session runs under, with no rule of its own).
Declined: note 3 (the lost case reads `start.tsv` back): a change of a test
after the last round is read by no round; in that case the Start was complete,
so `phase` writes nothing. Note 4, revealed and off the path, is #205.

## Verdict

Delivered: a restart reads and writes `start.tsv` by the harnesses that its
rows name, so a target whose harness register gained or lost a row after Start
restarts (#199). The review ended by decay at cycle 0 of cap 1. The diff against
`95dccce`, the branch's base, is inside 400 lines over 12 files. Next: the
specification task of `M2c`.

## Resource record

Recorded, not budgeted (ADR-0007). UTC; reviewers' tokens are `modelUsage`;
the author's are not reported.

| Part | Expected tier | Model | Effort | Tokens | Elapsed |
| ---- | ------------- | ----- | ------ | ------ | ------- |
| The plan | reasoning | Claude Opus 5.5 | not reported | not reported | 2026-10-10 16:00 to 16:07 |
| The plan review | reasoning | Claude Fable 5.1 | `xhigh` | 1,104,554 (USD 5.40) | 9 min 1 s |
| The code, test first | execution | Claude Opus 5.5 | not reported | not reported | 16:17 to 16:26 |
| Round 1 | reasoning | Claude Fable 5.1 | `xhigh` | 537,192 (USD 3.63) | 7 min 50 s |
| The close-out, with notes 1 and 2 | execution | Claude Opus 5.5 | not reported | not reported | 16:35 to 16:35 |
