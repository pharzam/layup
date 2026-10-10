# T-ye17 — no required context for a workflow that no pull_request event runs

Issue: [#201](https://github.com/pharzam/layup/issues/201), a child of #185 (the
state labels), split out by the Operator on the plan review of #185. Serves
`F-0003#52` (ADR-0026 rule 1). Base `95dccce`. Author: Claude Opus 5.5 on Claude
Code. Evidence: [`runs/T-ye17/`](../../runs/T-ye17/). The branch was pushed at
the base before the plan.

## Plan and plan review

The plan, its review and the author's answer are comments on #201. The plan
review (Claude Fable 5.1, `xhigh`, 7 min 28 s) gave `approve-with-conditions`:
320 lines over 20 files against `95dccce`, close-out inside; cap 2 (a check
script); no panel. Its four conditions are applied: the overlay of
`frame/good-passthrough` gets its `on:` block; every form not read is named,
each with a case; the local runs are the evidence, as CI runs `main`'s copy of
the check; a job of a skipped workflow must not be listed in the body.

## What was done

Check `protection` of `docs/setup/setup-check.sh` skips a workflow that no
`pull_request` event runs. It reads `on:` as a scalar, a one-line flow list, a
map or a block sequence; a key or an item counts only as a plain event name, and
anything else sends the workflow to the safe side (read as run by
`pull_request`, its jobs required). Cases `good-no-pull-request` and
`bad-pull-request-forms`; `docs/ci/README.md` states the rule. The four
workflows of the tree all run on `pull_request`, so no context of `main` moves.

## Review rounds

The records and the Fixes replies are comments on #201. Round 1 (`a02031c`,
cycle 0): `material`, a block sequence at the indent of `on:` and a split flow
list taken to the unsafe side, and the file count over budget (the Operator
approved 320 over 30). Round 2 (`3ee51b4`, cycle 1): `material`, an anchor in a
flow list and a flow collection on the next line; the reading turned into an
allow-list. Round 3 (`6843d2b`, cycle 2, the cap): `not mergeable, findings
recorded`, the budget only (325 lines at the head); its notes 2 to 5 are #207.

## Verdict

The Operator approved 400 lines over 30 files at landing (comment 6100100915),
so the task merges with that figure. Next: #202, then #185.

## Resource record

Recorded, not budgeted (ADR-0007). Times are UTC, 2026-10-10, from the clock and
the tools' own timestamps; the reviewers' tokens are `modelUsage`; the author's
are not reported.

| Part | Expected tier | Model | Effort | Tokens | Elapsed |
| ---- | ------------- | ----- | ------ | ------ | ------- |
| The claim and the plan | reasoning | Claude Opus 5.5 | not reported | not reported | 16:07 to 16:08 |
| The plan review | reasoning | Claude Fable 5.1, Claude Code CLI | `xhigh` | 815,350 (USD 4.05) | 7 min 28 s, from 16:08 |
| The tests, red; the check; the documents | execution | Claude Opus 5.5 | not reported | not reported | 16:16 to 16:21 |
| Round 1 | reasoning | Claude Fable 5.1, Claude Code CLI | `xhigh` | 486,283 (USD 3.92) | 8 min 33 s, from 16:21 |
| The fix of round 1 | execution | Claude Opus 5.5 | not reported | not reported | 16:41 to 16:46 |
| Round 2 | reasoning | Claude Fable 5.1, Claude Code CLI | `xhigh` | 728,712 (USD 4.35) | 9 min 10 s, from 16:46 |
| The fix of round 2 | execution | Claude Opus 5.5 | not reported | not reported | 16:55 to 17:01 |
| Round 3 | reasoning | Claude Fable 5.1, Claude Code CLI | `xhigh` | 652,783 (USD 4.18) | 7 min 59 s, from 17:01 |
| The close-out | execution | Claude Opus 5.5 | not reported | not reported | 17:12 to 17:14 |
