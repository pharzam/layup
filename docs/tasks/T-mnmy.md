# T-mnmy — in progress with recent activity, stalled without

Issue: [#202](https://github.com/pharzam/layup/issues/202), a child of #185 (the
state labels), from the Operator's meaning of "in progress" (comment 6099491323
on #185). Serves `F-0003#52` (ADR-0026 rule 1). Base `40aeae9`. Author: Claude
Opus 5.5 on Claude Code. Evidence: [`runs/T-mnmy/`](../../runs/T-mnmy/). The
branch was pushed at the base before the plan.

## Plan and plan review

The plan, its review and the author's answer are comments on #202. The plan
review (Claude Fable 5.1, `xhigh`, 8 min 47 s) gave `approve-with-conditions`:
450 lines over 40 files against `40aeae9`, close-out inside; cap 2; no panel.
Its four conditions are applied: criterion 1 edited on the issue to the push
time (a branch claimed at the base holds the base's commit, so its committer
date says nothing of the work) and criterion 2 to "at most 86400 seconds"; both
fetch bounds follow from the one window value; the glossary separates the state
`stalled` from **Stall**; a case gives one task both a push and a comment.

## What was done

`task-state.sh` splits `running`: a claimed task is `in progress` while its
branch has a push (a push, a force push or the branch's creation) or its issue a
comment created within `window` (86400 seconds) before `now`, and `stalled`
after. The reader takes `now`, `forge/pushes` and `forge/comments`; the fetcher
writes `now` from the clock and reads the forge's activity of the smallest
period that covers the window and the issue comments since `now` less the
window. Seven cases red first; Starting a task, `AGENTS.md`, the glossary rows
Claim and Stall, the suite README.

**Known limits.** Any comment counts, so a question posted to a stopped session
keeps its task `in progress` for a day. The host's clock sets `now` and the
forge's clock stamps the activity, so a skew moves the boundary by the skew. The
forge's period and the reader's window differ by the call's latency. A
`stalled` task is shown, not freed: what a session may do with one is a new
process rule, not this task's.

## Review rounds

The records and the Fixes reply are comments on #202. Round 1 (`b56cfb5`, cycle
0): `nothing material in scope`, five notes, applied in `4400a44` (a case where
the later push decides, a `now` that is not a number exits 1, the period
`quarter`, two comments). Round 2 (`4400a44`, cycle 1): `nothing material in
scope`, two notes applied in the close-out (the indent of a comment line; the
two runs of the run record).

## Verdict

Delivered: `in progress` and `stalled` in `task-state.sh`, with their cases and
documents. The review ended by decay at cycle 1, under the cap of 2. The diff
against `40aeae9` is inside 450 lines over 40 files. Next: #185, re-planned with
`state: stalled`.

## Resource record

Recorded, not budgeted (ADR-0007). Times are UTC, 2026-10-10, from the clock and
the tools' own timestamps; the reviewers' tokens are `modelUsage`; the author's
are not reported.

| Part | Expected tier | Model | Effort | Tokens | Elapsed |
| ---- | ------------- | ----- | ------ | ------ | ------- |
| The claim and the plan | reasoning | Claude Opus 5.5 | not reported | not reported | 17:32 |
| The plan review | reasoning | Claude Fable 5.1, Claude Code CLI | `xhigh` | 690,976 (USD 4.50) | 8 min 47 s, from 17:32 |
| The answer; the tests, red; the code; the documents | execution | Claude Opus 5.5 | not reported | not reported | 17:41 to 17:46 |
| Round 1 | reasoning | Claude Fable 5.1, Claude Code CLI | `xhigh` | 727,858 (USD 4.32) | 8 min 21 s, from 17:46 |
| The notes of round 1 | execution | Claude Opus 5.5 | not reported | not reported | 17:55 to 17:56 |
| Round 2 | reasoning | Claude Fable 5.1, Claude Code CLI | `xhigh` | 817,911 (USD 4.30) | 7 min 17 s, from 17:56 |
| The close-out | execution | Claude Opus 5.5 | not reported | not reported | 18:04 to 18:06 |
