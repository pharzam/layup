# T-57vx — the build tasks of M2e

Issue: [#346](https://github.com/pharzam/layup/issues/346), the slicing task that
`T-xs4y` (#333) named, under the Operator's directive ("go ahead with M2E until
finish it"). Serves `F-0003#51` and the requirements of `M2e`. Base `8f034a7`
(the merge of #334). Author: Claude Sonnet 5.5 on Claude Code. Evidence:
[`runs/T-57vx/`](../../runs/T-57vx/).

## Plan and plan review

The plan (R12), its review and the author's answer are comments on #346. The plan
review (GPT-6.1 Sol on Codex CLI, effort `xhigh`, a fresh sandboxed session in a
clone at `8f034a7`, 7 min 11 s; the harness is outside the routing rule of ADR-0026
and was allowed by the Operator on 2026-10-13) gave `approve-with-conditions`:
Budget maximum 1,400 lines added plus removed over 18 files, close-out inside;
Cycle cap 1; no panel. Its five conditions are applied: the package owners
separated (the check of a handoff in `internal/handoff`; the posting, the events
and the routing in `internal/run`; the tier and the prompt file in `internal/run`);
each cut clause has its own entry; the order repaired (the task-loop machinery
before Shape); each row bounded (demo, count, tests, size, lines, cap); the
validation named. **The Operator's decision:** O-214 (comment 6197053800): the 24
rows, 60 to 68b, with the count of each final.

## What was done

1. **T:** [`parts.tsv`](../../runs/T-57vx/parts.tsv) puts each part of the
   specification of `M2e`, and each clause that the slicing cuts, in one row, or
   names it with a reason; [`parts.sh`](../../runs/T-57vx/parts.sh) fails before the
   table and passes after, and six mutant runs fail as they should;
   [`issues.sh`](../../runs/T-57vx/issues.sh) reads the forge.
2. **D1:** [The tasks of M2e](../plan/README.md#the-tasks-of-m2e), rows 60 to 68b,
   with the release of `M2e` defined; the first paragraph, the milestone row, the
   hosts of three items and the row of #68 of the plan.
3. **D3:** issues #347 to #370, opened by the App through
   [`gen-issues.py`](../../runs/T-57vx/gen-issues.py), copied from `T-bng8`'s;
   row 68b's issue lists the Operator's inputs (R6), row 65b's the question about
   `CODEOWNERS`. **The backlog** has the 24 lines.

**Revealed, off the path, as issues:** #372 (the record of the Retrospective's event)
and #373 (the forge calls to edit and close a pull request, and the home of the
stand-in).

## Review rounds

The records are comments on #346. Both rounds were Claude Opus 5.5, effort `xhigh`,
a fresh read-only `claude -p` session in a clone, a model that wrote none of the
change.

| Round | Head | Cycle | Cap | Material | Then |
| ----- | ---- | ----- | --- | -------- | ---- |
| 1 | `f17f880` | 0 | 1 | 4 | fixes in `9553dd8`; round 2 (cycle 1) is the last the cap allows |
| 2 | `9553dd8` | 1 | 1 | 2 | **O-215**: settled as fixed by `d586682` |

Round 2 keeps its verdict, `not mergeable, findings recorded`. By **O-215**, the
Operator settled its findings as fixed by `d586682`, as O-207 (b) and O-213 did, and
no round reads that commit. The verdict of round 1 was edited to `material`, as RR8
asks.

## Verdict

Delivered: the 24 build tasks of `M2e`, with their table, their issues and their
checks.

- **How the review ended:** at cycle 1 of the cap 1, by O-215 and not by decay: the
  material findings were 4 and 2.
- **The size:** the diff against `8f034a7` is inside 1,400 lines over 18 files.
- **Next:** rows 60, 62a, 62b, 63b, 64 and 65a, which can start at once.

## Resource record

Recorded, not budgeted (ADR-0007). UTC, 2026-10-15. Reviewers' tokens are as the
harness reports them; the author's are not reported.

| Part | Expected tier | Model | Effort | Tokens | Elapsed |
| ---- | ------------- | ----- | ------ | ------ | ------- |
| The plan | reasoning | Claude Sonnet 5.5 | not reported | not reported | 14:40 (the issue) to 14:41 (the plan's comment) |
| The plan review | reasoning | GPT-6.1 Sol | `xhigh` | 207,843 | 7 min 11 s |
| T, the parts map and the generator | execution | Claude Sonnet 5.5 | not reported | not reported | 14:49 to 15:20 (the wait for O-214 is in it) |
| D1 and D3 | execution | Claude Sonnet 5.5 | not reported | not reported | 16:15 to 16:35 |
| Round 1 | reasoning | Claude Opus 5.5 | `xhigh` | 7,580,396 (USD 4.82) | 12 min 27 s |
| The fix of round 1 | execution | Claude Sonnet 5.5 | not reported | not reported | 16:48 to 16:53 |
| Round 2 | reasoning | Claude Opus 5.5 | `xhigh` | 7,019,253 (USD 5.32) | 14 min 37 s |
| The fix of round 2 | execution | Claude Sonnet 5.5 | not reported | not reported | 17:09 to 17:16 |
| The close-out | execution | Claude Sonnet 5.5 | not reported | not reported | 17:25 to 17:30 |
