# T-bng8 — the build tasks of M2d

Issue: [#267](https://github.com/pharzam/layup/issues/267), the slicing task that
`T-msrd` (#265) named for `M2d`, as `T-2r8u` (#211) did for `M2c`. Serves `F-0003#51`
and the requirements of `M2d`. Base `347d845` (the merge of #266). Authors: Claude
Opus 5.5 on Claude Code (the plan, the first head, the fix of round 1) and SWE-2 Max
on Devin (the fixes of rounds 2 to 6); the close-out is by Claude Sonnet 5.5 on Claude
Code. Evidence: [`runs/T-bng8/`](../../runs/T-bng8/).

## Plan and plan review

The plan (R12), its review and the author's answer are comments on #267. The plan
review (GPT-6 Sol on Devin, a fresh session in a clone at `347d845`, 1 min 25 s;
Claude Fable 5.1, first in the order, was at its usage limit) gave
`approve-with-conditions`: Budget maximum 1,000 lines added plus removed over 17 files
against `347d845`, close-out inside; Cycle cap 1; no panel. Its four conditions are
applied: row 54 is after 52b, through 52c, and `setup-record` has its two clauses in
rows 52a and 52b; the check accounts for the records of the Scaffold, the calls of
`M2d`, the table of `M2d`, `CommitTree` and each built block changed in words; rows 53
and 56 carry a test of the merge and the comment and of the work area built again;
each issue records its duplicate check, and the plan's list of sliced milestones is
current. The answer split two rows (52b into 52b and 52c, 57 into 57a and 57b): 15 rows.

**The Operator's decisions:** O-209 (comment 6118820298): the 15 rows, 50 to 59b,
with the count of each final. At each of the six caps the Operator chose (a), one
more cycle; the cap is now 7 (the `## Plan review` comments of #267).

## What was done

1. **T:** [`parts.tsv`](../../runs/T-bng8/parts.tsv) puts each part of the
   specification of `M2d` in one row, or names it with a reason: the 30 headings of
   `scaffold.md`, its 21 input states and 14 acceptance tests (each clause that the
   slicing cuts), the 12 clauses of the records of the Scaffold, the calls of `M2d`,
   and the table of `M2d` with `CommitTree` (95 entries).
   [`parts.sh`](../../runs/T-bng8/parts.sh) has six checks; it failed before the table
   (`77fac1c`) and passes after.
2. **D1:** [The tasks of M2d](../plan/README.md#the-tasks-of-m2d), rows 50 to 59b, with
   the rule of the slicing and the release of `M2d`; the plan's list of milestones, the
   row of `M2d` and the hosts of `later-p2-scaffold` name the tasks.
3. **D2:** none owed; `scaffold.md` is unchanged.
4. **D3:** issues #270 to #284, opened by the App through
   [`gen-issues.py`](../../runs/T-bng8/gen-issues.py), copied from `T-2r8u`'s and
   changed for `M2d`, each body with its duplicate check (R2); row 59b's issue lists
   the Operator's inputs (R6). [`issues.sh`](../../runs/T-bng8/issues.sh) failed before
   the issues and passes after (it reads the forge, so it runs by hand).
5. **The backlog** has the 15 lines, each with its issue.

## Review rounds

The records are comments on #267. Rounds 1 to 5 were GPT-6 Sol on Devin, and rounds 6
and 7 Claude Fable 5.1 on `claude -p`, effort `max` (read-only by its allow rules);
each is a fresh session in a clone, on a model that differs from both authors. The verdict of each of rounds 1
to 6 was edited to `material` when the Operator raised the cap. The record of round 2
had no heading `## Review record — round 2`, so `review-record-lint` read the rounds as
1, 3, 4, 5 and 6 and its check was red on each push from `39d9478` to `53b4c78`; the
heading was added before round 7 (11:12Z), the text under it unchanged.

| Round | Head | Cycle | Cap | Material | The findings | Then |
| ----- | ---- | ----- | --- | -------- | ------------ | ---- |
| 1 | `ffae841` | 0 | 1 | 1 | the check compared the map with the plan by row number, so a part could leave its row's cell | (a): cap 2 |
| 2 | `b1b7066` | 1 | 2 | 3 | the changed cells of the table of `M2d` had no row; check 5 passed with a clause of a built block dropped; the thread had no cap of 2 and no `material` verdict of round 1 | (a): cap 3 |
| 3 | `39d9478` | 2 | 3 | 2 | `setup/rulesets.txt` and `setup/inputs/files` had no row; a part copied into a second Task cell passed | (a): cap 4 |
| 4 | `d3df594` | 3 | 4 | 1 | the `--body` example of `gen-issues.py` named row 40, which `M2d` does not have | (a): cap 5 |
| 5 | `7ae4a30` | 4 | 5 | 1 | a clause of a split acceptance test or input state could be dropped | (a): cap 6 |
| 6 | `cc2cac5` | 5 | 6 | 2 | row 53 and #276 named a clause that its acceptance test lacks; a clause of a split heading could be dropped | (a): cap 7 |
| 7 | `53b4c78` | 6 | 7 | 0 | `nothing material in scope` | — |

Each material finding was fixed test first: a red run of the mutant, then green
([`test-runs.md`](../../runs/T-bng8/test-runs.md)). Rounds 1 to 6 also had six notes:
five fixed with their round, and the budget note of round 6 recorded, no change.
The gap that round 6 revealed, no acceptance test of the merge into the batch and the
comment, is #287. Round 7 ran with the allow rules narrowed to the read-only forms (#288).

Round 7 ended `nothing material in scope`, with four notes. The close-out may change
only the task indexes and this file, and the fix of each note edits a file that a round
read (`parts.tsv`, `parts.sh`, the plan or `issues.sh`), which would re-freeze the head.
None changes an exit code or leaves a part unassigned, so each is declined, as `T-2r8u`
declined its notes 3 to 6:

1. **Three `none` entries with thin reasons** (`parts.tsv` lines 13, 25 and 29): the
   map is this task's evidence, not an operative text; the reviewer found each rule
   reachable from its row (52b; 52b; 56, 55 and 57b), and each build task reads its
   section whole.
2. **The clause form of "The answers of the setup: the rows of answers.tsv"**: the
   clause names the whole test, so nothing is unassigned; the fix edits the check, the
   map, the Task cell of row 54 and the body of #277.
3. **"D2" in the header of `issues.sh`**: a label in a comment, copied from
   `T-2r8u`'s `issues.sh`, whose note 5 was declined alike (the issues are D3).
4. **The summary of the specification in the plan** names four of its five parts: the
   Task cell of row 51b is the operative text and names `CommitTree`.

## The lesson of the task (gate step 7)

Two traps are owed to `docs/guardrails.md` §2, and are
[#288](https://github.com/pharzam/layup/issues/288): a review record without its
heading is no record to the linter (round 2: the required check was red on five
pushes), and a reviewer's allow rule that names a script with a write mode
(`gen-issues.py --post`). The close-out may change only the task indexes and this
file, so the entries land by #288, each read by its own round.

## Verdict

Delivered: the build tasks of `M2d`, rows 50 to 59b of the plan, each with its issue
(#270 to #284), by O-209. The review ended by decay at cycle 6 of cap 7, after ten
material findings in rounds 1 to 6, each fixed. The four criteria of #267 hold and are
ticked. The diff against `origin/main` (`5dbd7af`, after the merge of #269) is 979
lines over 11 files with the close-out, inside the Budget maximum of 1,000 lines
over 17 files. No round read the merge of `origin/main` (clean, with none of the 9
files of the branch) or this close-out. Next: the build of `M2d`, rows 50, 51a, 51b and
52a first.

## Resource record

Recorded, not budgeted (ADR-0007). UTC, 2026-10-12 unless named. Devin reports no
usage; the author's tokens are not reported; the reviewers' of rounds 6 and 7 are
`modelUsage`, with Claude Code's own estimate of the money. Claude Sonnet 5.5 is the
model the Operator chose for the last session (`/model`); Model tiers names `Claude
Sonnet 5` for the execution tier. A total would need the tokens that no harness reported.

| Part | Expected tier | Model | Effort | Tokens | Elapsed |
| ---- | ------------- | ----- | ------ | ------ | ------- |
| The plan | reasoning | Claude Opus 5.5 | not reported | not reported | 2026-10-11, to 21:08 |
| The plan review | reasoning | GPT-6 Sol | not reported | not reported by Devin | 1 min 25 s |
| The answer and the counts to the Operator | reasoning | Claude Opus 5.5 | not reported | not reported | 2026-10-11, 21:10 to 21:11 |
| The wait for O-209 | — | — | — | — | 9 h 3 min, 2026-10-11 21:10 to 06:13 |
| T, D1 and D3, after O-209 | execution | Claude Opus 5.5 | not reported | not reported | 06:13 to 06:19 |
| Round 1 | reasoning | GPT-6 Sol | not reported | not reported by Devin | 6 min 25 s |
| The fix of round 1 | execution | Claude Opus 5.5 | not reported | not reported | 06:26 to 06:28 |
| Rounds 2 to 5 | reasoning | GPT-6 Sol | not reported | not reported by Devin | 8 min 41 s, 7 min 55 s, 6 min 46 s, 8 min 39 s |
| The fixes of rounds 2 to 6 and their replies | execution | SWE-2 Max | not reported | not reported by Devin | commits at 08:41, 09:10, 09:47, 10:16 and 10:51 |
| Round 6 | reasoning | Claude Fable 5.1 | `max` | 1,361,577 (USD 7.21) | 15 min 55 s, 10:26 to 10:42 |
| Round 7 | reasoning | Claude Fable 5.1 | `max` | 1,846,300 (USD 8.12) | 17 min 27 s, 11:16 to 11:34 |
| The cap record, the repair of round 2's heading, the narrowed rules and the probe, the launch of round 7 | execution | Claude Sonnet 5.5 | `max` | not reported | 11:11 to 11:16 |
| The merge of `origin/main`, #288 and the close-out | — | Claude Sonnet 5.5 | `max` | not reported | 11:34 to 11:44 |
