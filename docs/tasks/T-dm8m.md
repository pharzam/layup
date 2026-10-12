# T-dm8m — the state of each task as an issue label

Issue: [#185](https://github.com/pharzam/layup/issues/185), the child of #175 (`T-eep8`)
that it left out on purpose. Serves `F-0003#52` (ADR-0026 rule 1). Base `a4d41cf`; the first
plan was at `27efb4e`. Authors, as the record of round 7 names them: Claude Opus 5.5 on
Claude Code and SWE-2 Max on Devin. The merge of `origin/main` and this close-out are by
Claude Sonnet 5.5 on Claude Code. Evidence: [`runs/T-dm8m/`](../../runs/T-dm8m/).

## Plan and plan reviews

The plan (R12, 2026-10-10) and its review (Claude Fable 5.1, effort `xhigh`, 11 min 20 s)
gave `approve-with-conditions`: Budget maximum 900 lines added plus removed over 60 files,
close-out inside; Cycle cap 2, as the task changes CI; no panel. The Operator's rulings of
2026-10-10 split two tasks off, #201 (the exemption of check `protection`) and #202 (the
state `stalled`), both merged first, and the task was planned again on `a4d41cf`. The
review of the re-plan (GPT-6 Sol on Devin, 3 min 23 s) gave `approve-with-conditions` with
the same maximum and cap; its four conditions were applied. An approval of 2026-10-12,
relayed on the issue by the author, named 1,150 lines over 80 files. The rounds read it as a
ceiling, which the gate does not allow: an approval is the figure at landing.

## What was done

1. **T:** sixteen fixture cases under `docs/tasks/tests/`, red at 06:04Z against a stub that
   prints nothing, each for its own reason ([`test-runs.md`](../../runs/T-dm8m/test-runs.md)).
2. **D1, D2:** [`task-labels.sh`](task-labels.sh). Its reader (`--plan DIR`) gives `add` and
   `remove` lines from the state that `task-state.sh` derives and the issues of the forge; its
   applier writes them only after the command, the three files, the issues and the reader have
   all succeeded; `--dry-run` writes nothing; `--check-workflow` holds the rules of the workflow's
   header. [`task-state.yml`](../../.github/workflows/task-state.yml) is the only writer: `push`
   to `main`, `pull_request_target`, `issues`, `issue_comment` and a schedule of 30 minutes, with
   `issues: write` as its only write permission (and `contents: read`, `pull-requests: read`),
   and the scripts restored from the default branch.
3. **D4:** `docs/ci/README.md`, the enforcement table of `docs/issue-workflow.md`, Starting a
   task, `AGENTS.md`, the glossary and the tests README name the labels.
4. **The lesson (step 7):** `docs/guardrails.md`, "A trigger that runs a branch's copy of a
   write-token workflow".
5. **The seven labels** are created by the applier at each run (`gh label create --force`): the
   first run on `main` after the merge creates them, and no step is by hand. Criterion 1 of #185
   is shown on the forge by that run, as the plan records.

## Review rounds

The records are comments on #185. Rounds 1 to 6 were GPT-6 Sol on Devin and round 7 Claude
Fable 5.1 on `claude -p`, effort `max`. The Operator chose (a), one more cycle, five times.

| Round | Head | Cycle | Cap | Material | The findings | Then |
| ----- | ---- | ----- | --- | -------- | ------------ | ---- |
| 1 | `aafee79` | 0 | 2 | 7 | a removal applied as an addition; the workflow check read no event bound; two no-write cases never ran the writer; the documented diff check failed; the instructions contradicted the population; the budget; the lesson of the trigger missing | fixed, `33156a2` |
| 2 | `33156a2` | 1 | 2 | 4 | the budget; a push bound could be a tag; the empty close reason uncovered; summaries named the wrong population | (a): cap 3; `e4aabfe` |
| 3 | `e4aabfe` | 2 | 3 | 4 | the budget; the record would fail `review-record`; the push-bound check had a false green; a summary asserted a label where none is due | (a): cap 4; `ea43afe` |
| 4 | `8e7c48a` | 3 | 4 | 2 | the budget; the issue gave the old label count and creation step | (a): cap 5; `ea9aa4e` |
| 5 | `ea9aa4e` | 4 | 5 | 2 | the budget; the issue promised a label on every task issue | (a): cap 6; `5acd43c` |
| 6 | `5acd43c` | 5 | 6 | 3 | the budget; the CI account omitted the population and the precedence; a quoted push bound was missed | (a): cap 7; `8d1b77b` |
| 7 | `8d1b77b` | 6 | 7 | 2 | the budget, 1,803 lines over 98 files; the CI account named the wrong script | the Operator's approval of the figure; `997d523` |
| 8 | `cbc17be` | 7 | 7 | 0 | scoped to `docs/guardrails.md`, which the merge of `origin/main` changed and the branch also touched; it consumes no cycle | — |

The finding of the landing size stood in each of the seven rounds, 7 of the 24 material findings:
the approvals were figures set before the work ended, and the figure at landing passed each of
them. No author and no reviewer can settle it; only the Operator can, and it is settled once,
below, on the final figure with this close-out counted. The other 17 were fixed, each with its
evidence in `test-runs.md`, and round 7 added four notes, all fixed in `997d523`.

Round 8 (Claude Fable 5.1, effort `xhigh`, a fresh read-only session in a clone at `cbc17be`, its record a
comment on #185): `nothing material in scope`, two notes, both declined, since the close-out may change only
the task indexes and this file: the `Learned in` of the new entry names one of the two plan reviews that
taught the lesson, and the lead-in of "Reference-sweep pitfalls" is stale on `main` before this branch.

## The end

The Operator chose (b), settle, in the session of 2026-10-12 ("approve"), and approved **1,966 lines
added plus removed over 100 files**, the figure at landing with this close-out inside, once
(comment 6125160894 of #185). That settles the finding of the budget. Round 7 ended `not mergeable,
findings recorded` and no round read its fixes; the Operator's comment permits the merge.

The fixes of `997d523` (the sentence of the CI account, the `|| 'main'` fallback of the schedule's
Restore step, the `LIMIT` comment of the check, the reason for leaving out `delete`, one row of the
tests README) were read by no round. Their checks are green: 127 discipline cases, `--check-workflow`,
and the battery.

## Verdict

Delivered: the labels of #185, set from the derived state by `task-labels.sh` and the workflow
`task-state.yml`, under 127 discipline cases that were red first. The review ended at the cap of 7
by the Operator's approval of the landing figure. The diff against `origin/main` (`f8f125f`, after
the merge of #290) is 1,966 lines over 100 files with the close-out, the figure that the
Operator approved. Next: the first run of the workflow on `main`, which creates the labels and shows
criterion 1 on the forge.

## Resource record

Recorded, not budgeted (ADR-0007). UTC, 2026-10-12 unless named. Devin reports no usage; the
reviewers' tokens of Claude Code are `modelUsage`, with its estimate of the money; the author's
are not reported. The Operator chose Claude Sonnet 5.5 for the last session; Model tiers names
`Claude Sonnet 5` for the execution tier.

| Part | Expected tier | Model | Effort | Tokens | Elapsed |
| ---- | ------------- | ----- | ------ | ------ | ------- |
| The plan | reasoning | Claude Opus 5.5 | not reported | not reported | 2026-10-10, to 07:28 |
| The plan review | reasoning | Claude Fable 5.1 | `xhigh` | 1,209,273 (USD 6.85) | 11 min 20 s, 2026-10-10 07:28 to 07:40 |
| The re-plan | reasoning | Claude Opus 5.5 | not reported | not reported | 2026-10-11, to 20:52 |
| The review of the re-plan | reasoning | GPT-6 Sol | `xhigh` | not reported by Devin | 3 min 23 s, 2026-10-11 20:53 |
| T, D1, D2, D4 and the first head | execution | Claude Opus 5.5 | not reported | not reported | 06:04 to 06:11 |
| Rounds 1 to 6 | reasoning | GPT-6 Sol | not reported | not reported by Devin | 7 min 26 s, 7 min 56 s, 9 min 16 s, 9 min 18 s, 7 min 27 s, 8 min 6 s |
| The fixes of rounds 1 to 7 | execution | SWE-2 Max | not reported | not reported by Devin | commits at 07:01, 07:16, 08:30, 09:08, 09:57, 10:14 and 10:48 |
| Round 7 | reasoning | Claude Fable 5.1 | `max` | 1,652,684 (USD 8.60) | 19 min 35 s, 10:26 to 10:45 |
| Round 8, scoped | reasoning | Claude Fable 5.1 | `xhigh` | 555,448 (USD 3.66) | 4 min 55 s, 13:15 to 13:20 |
| The merge of `origin/main`, the scoped round and the close-out | — | Claude Sonnet 5.5 | `max` | not reported | 13:10 to 13:23 |
