# T-8r3u — row 46 of the plan, the batch

Issue: [#222](https://github.com/pharzam/layup/issues/222), row 46 of [the tasks
of M2c](../plan/README.md#the-tasks-of-m2c), opened by `T-2r8u` (#211). Serves
`F-0003#41` through REQ-001. Base `06dff9c` (the merge of #250, row 45b); the
budget is read against it, the base of this branch's own work. Author: Claude
Opus 5.5 on Claude Code. Evidence: [`runs/T-8r3u/`](../../runs/T-8r3u/).

## Plan and plan review

The plan (R12), its review and the author's answer are comments on #222. The
plan review (Claude Fable 5.1, effort `xhigh`, a fresh read-only session in a
clone at `06dff9c`, 9 min 11 s) gave `approve-with-conditions`: Budget maximum
1,400 lines added plus removed over 24 files, close-out inside; Cycle cap 1; no
panel. Its four conditions are applied: the dropped case of the test is two rows
of one source; the rule of `asked` is in `records.md`; `run.md` says that the
App's own comment is read with no copy; the second rule gap of one line and one
rule is a known limit. The goal count is 2, final by O-196: the merge and the
IDs; the comment, its two blocks and the intent form.

## What was done

1. **The merge:** the rule gaps (a gap's quote is its line, trimmed of its line
   break and a CR), the gaps of meaning, the items of the completeness review
   (an `uncovered` or `class` item quoted by its fact's span) and the three
   fields of the intent form, in that order; a repeat of quote and kind
   dropped, a question with no quote never; `Q-NNN` in order.
2. **The comment:** the header line, the block for the idea owner (each
   question and its quote in backquotes), the block for the Operator (none in
   `M2c`), the rule of the answers.
3. **The step `batch`** after `completeness`: `questions.tsv` in one records
   commit, then the comment on the Intake issue; "done before" with both; the
   comment posted from the record when only the record is there. `run-steps`
   gains `batch`.
4. **The test forges** keep a posted comment on its issue (the world of
   `internal/run`) and take a comment on the Intake issue (the e2e world).
5. **The documents:** `intake.md`, `records.md` (`asked`), `run.md` (Copy
   before read; the block), the Job cell of `internal/run`, `traceability.md`,
   the Task and Test cells of REQ-001 and a §13 line.

**Tests:** [`test-runs.md`](../../runs/T-8r3u/test-runs.md): red before the code;
15 of 15 mutations caught; green at the three levels.

## Known limits

- **Two rule gaps of one line and one rule** (G3, two undefined abbreviations of
  one line): the second is dropped by the merge's rule, so its abbreviation is
  not asked (condition 4 of the plan review; `intake.md`).

## Review rounds

The record is a comment on #222. Round 1 (Claude Fable 5.1, effort `xhigh`, a
fresh read-only session in a clone at `08140f4`, cycle 0, the cap): `nothing
material in scope`, four notes (a quote with a blank line, one sentence of
Copy before read, the match of the header, a bounds check of the post from the
record), which are #254.

## Verdict

Delivered: the step `batch`, the merge of the questions into `questions.tsv`
with their IDs, and one comment on the Intake issue with its two blocks and the
intent form, posted after the records commit, or from the record after a stop.
The review ended by decay at cycle 0 of cap 1. The diff against `06dff9c`, the
branch's base, is inside 1,400 lines over 24 files. CI runs the three levels on
the landing head. Next: row 47a.

## Resource record

Recorded, not budgeted (ADR-0007). UTC; reviewers' tokens are `modelUsage`;
the author's are not reported.

| Part | Expected tier | Model | Effort | Tokens | Elapsed |
| ---- | ------------- | ----- | ------ | ------ | ------- |
| The plan | reasoning | Claude Opus 5.5 | not reported | not reported | 2026-10-11 10:00 to 10:01 (the claim to the post) |
| The plan review | reasoning | Claude Fable 5.1 | `xhigh` | 1,062,375 (USD 5.55) | 9 min 11 s |
| The code, test first | execution | Claude Opus 5.5 | not reported | not reported | 10:11 to 10:24 |
| Round 1 | reasoning | Claude Fable 5.1 | `xhigh` | 968,068 (USD 4.64) | 11 min 10 s |
| The close-out | execution | Claude Opus 5.5 | not reported | not reported | 10:36 to 10:36 |
