# T-qhpf — the three notes of round 2 of T-eep8

Issue: [#182](https://github.com/pharzam/layup/issues/182), opened by `T-eep8`
(#175) for the notes of its round 2 that change the command or the rule. Serves
`F-0003#52`, as a child of #175 (ADR-0026 rule 1). Base `91e41b5` (the merge of
#190). Author: Claude Opus 5.5 on Claude Code. Evidence:
[`runs/T-qhpf/`](../../runs/T-qhpf/). The branch `T-qhpf` was pushed at the base
before the plan, the reading that note 2 writes down.

## Plan and plan review

The plan (R12, comment 6087977850), its review (6088094461) and the author's
answer (6088094760) are comments on #182. The plan review (Claude Fable 5.1,
effort `xhigh`, on the Claude Code CLI with stream output, a fresh read-only
session in a clone at `91e41b5`, 6 min 54 s) gave `approve-with-conditions`:
Budget maximum 400 lines added plus removed over 16 files against `91e41b5`,
close-out inside; Cycle cap 2 (note 2 changes the order of the gate's preamble);
no panel. Its conditions: the goal count of 3 goes to the Operator (R11), the
case drops `T-aaa2` from its completed log too, so it has a red; and plan step
2 is stated in the future tense. The Operator's ruling, **one task of three
goals**, was given in the session `issue-status` and relayed to #182 by the
author (comment 6088204010); the Operator has not confirmed it on the issue.

## What was done

1. **Note 1:** a `## Now` task with no row in the last task table keeps the
   `After` of its row in any earlier table, so `task-state.sh` can print it
   `blocked` (case `good-earlier-table`).
2. **Note 2:** the claim (the worktree and its pushed branch) is the first act
   of taking a task, before its plan comment: [Starting a task](../engineering-discipline.md#starting-a-task),
   the gate's preamble and step 1, and `AGENTS.md` (the preamble, step 1, the
   branches).
3. **Note 3:** the fetcher reads `docs/tasks/backlog.md`, `docs/tasks/completed.md`
   and `docs/plan/README.md` from the forge's `main` with `gh api …/contents`,
   not from the checkout; `AGENTS.md` and the onboarding say so.

**Tests:** [`test-runs.md`](../../runs/T-qhpf/test-runs.md): `good-earlier-table`
red against the reader of `91e41b5` (`T-aaa2` printed `ready`), then 97 passed;
for note 3, a run from a clone seven merges behind printed the same bytes as the
run from this branch, where the old script differed in eight lines.

**Known limits.** The fetcher still needs a checkout, to name the repository for
`gh`. The claim stays a written rule: nothing refuses a plan by a session that
has not pushed its branch.

## Review rounds

The record is a comment on #182. Round 1 (Claude Fable 5.1; `b4cb3fb`, cycle 0):
`nothing material in scope`, four notes, applied in the close-out: the comment
of `fetch` (note 1), the onboarding line the plan named (note 2), the count of
the stale clone, seven merges and not eight (note 3), and this file's sentence
that the Operator's ruling was relayed (note 4).

## Verdict

Delivered: the three notes of #182. The review ended by decay at cycle 0, under
the cap of 2. The diff against `91e41b5` is inside 400 lines over 16 files.
Next: #185, the state labels that a CI job sets from this command.

## Resource record

Recorded, not budgeted (ADR-0007). Times are UTC, 2026-10-09, from the clock
after each step; the reviewers' tokens are `modelUsage` of the CLI's `result`
event; the author's are not reported.

| Part | Expected tier | Model | Effort | Tokens | Elapsed |
| ---- | ------------- | ----- | ------ | ------ | ------- |
| The claim and the plan | reasoning | Claude Opus 5.5 | not reported | not reported | 19:39 |
| The plan review | reasoning | Claude Fable 5.1, Claude Code CLI | `xhigh` | 696,222 (USD 4.61) | 6 min 54 s, from 19:40 |
| The answer; the Operator's ruling; the test, red; the code; the documents | execution | Claude Opus 5.5, a reasoning-tier model on an execution part | not reported | not reported | 19:47 to 19:57 |
| Round 1 | reasoning | Claude Fable 5.1, Claude Code CLI | `xhigh` | 718,686 (USD 4.23) | 8 min 51 s, from 19:57 |
| The close-out | execution | Claude Opus 5.5 | not reported | not reported | 20:06 to 20:08 |

No harness was skipped: the first in the reviewer order (`claude`) returned a
record each time.
