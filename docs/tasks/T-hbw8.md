# T-hbw8 — the soft reset of the architecture

Issue: [#72](https://github.com/pharzam/layup/issues/72). The Operator's
decision O-72 (option A, the soft reset) on PR #69, which closed without a
merge; its content is kept at the tag `archive/t-0kn4-architecture`. Supersedes
#67. Parent #42. Serves `F-0003#51` and gives each In-Scope item `F-0003#41`–`#52`
a component. Base `e08f584`. Author: Claude Opus 5.5 on Claude Code. Evidence:
[`runs/T-hbw8/`](../../runs/T-hbw8/).

## Plan and plan review

The plan (R12) and its review are comments on #72. The plan review (Claude
Fable 5.1 on Claude Code, a fresh session, 10 min 42 s) gave
`approve-with-conditions`: Budget maximum 3,200 lines added plus removed over 36
files on the whole diff against `e08f584`, close-out inside; Cycle cap 1. The
author applied the eight conditions and the eleven notes (comment "Plan: the
author's answer to the plan review" on #72).

## What was done

1. **Inputs in Git.** The two #69 comments word for word
   ([`inputs-from-pr-69.md`](../../runs/T-hbw8/inputs-from-pr-69.md)) and the
   Jev documentation the panel read
   ([`jev-sources.md`](../../runs/T-hbw8/jev-sources.md)).
2. **The panel** (ADR-0006, O-70): one pass per member, no member saw another's
   output, each at most 300 lines
   ([`panel-brief.md`](../../runs/T-hbw8/panel-brief.md)):
   A — orchestration, Git state, the forge (GPT-6 Sol on Devin, 6 min 54 s);
   B — decisions, routing, learning (Grok 4.7: on OpenCode no output in
   15 minutes, skipped; then on Devin, 4 min 28 s); C — governance (Claude Opus
   5.5 on `claude -p`, 4 min 40 s). Authors under Bootstrap mode rule 4: Claude
   Opus 5.5, GPT-6 Sol, Grok 4.7.
3. **The selection** ([`selection.md`](../../runs/T-hbw8/selection.md)): the
   search for public solutions first, then the considerations of Solution
   selection per question; eight choices went to the Operator in one batch.
4. **The Operator's decisions** O-76 to O-84
   ([`operator-decisions.md`](../../runs/T-hbw8/operator-decisions.md)): native
   stack gates in the target and LAYUP strictly external (O-76, which supersedes
   O-10 and O-11 where they put the gates outside); agents under the Operator's
   account with the App badge, no second account (O-77); the smart-if provider,
   its authority, the thresholds, the screen frequency, the harness eligibility,
   the role matrix, `T`, `N` and the learning policy are parameters with
   defaults (O-78 to O-83); the parameters cannot switch off a PSB invariant,
   and "smart-if" is a conditional branch (O-84).
5. **The first draft, withdrawn:** ADR-0013 to ADR-0020 and
   `docs/architecture.md`, with the lines they added to the glossary, `PRD-0001`,
   `docs/setup/README.md`, `README.md`, the onboarding and the ADR index. The
   deep check found about 50 material defects in it; plan v2 (O-92) removed it
   in step 0 and rewrites it. Its text stays in the history of this branch.
   This file is rewritten at the close-out of plan v2.

## Decisions of the author (the withdrawn draft)

- **ADR-0011 is amended, not superseded:** decision 2 (the records branch,
  ADR-0016), decisions 3, 4 and 7 (native stack gates, ADR-0017), and one
  sentence of decision 8 (ADR-0014). Its other decisions stay.
- **The records branch** answers ADR-0011's rejection of a separate branch
  ("outside the tree a reviewer reads") by fetching it by name and citing its
  commit in each result; the reason for it is that no agent and no product pull
  request can change a record.
- **The ambiguity-kind map** of ADR-0015 is a default with no fact behind it:
  the PSB names the roles and the kinds but maps none, and panel members A and C
  gave two different maps.
- **The panel synthesis** is written by a reasoning session; the smart-if only
  scores the written options (the provider does not generate text; panel member
  C, point P1).
- **Eight ADRs,** the limit that the plan review set (note 9); the parameter
  register is in ADR-0013, not a ninth record.

## Known limits (the withdrawn draft)

- GitHub is the only forge for the pilot.
- With one GitHub account, the agent–human separation holds only while agent
  sessions cannot reach the Operator's own credentials; a merge or a rule-path
  change by another actor is detected, not prevented (ADR-0016, ADR-0017).
- The cost stop acts between actions; spend inside one action stops only where a
  harness enforces a cap (ADR-0018).
- The start values of the thresholds, `H`, `cost.task_factor` and the reward
  terms are set with their evidence by the implementation plan and the pilot.

## Lessons

The lessons that the next reader could hit are in
[`docs/guardrails.md`](../guardrails.md) §2, each learned in this task:

- "A public-solution search that confirms the design" and "Coverage by name":
  from the withdrawn first draft and its deep check
  ([`runs/T-hbw8/root-cause-missed-solution.md`](../../runs/T-hbw8/root-cause-missed-solution.md)).
- "A tool under test writes into the host's home directory": from the hands-on
  evaluation, where three tools wrote into the Operator's home directory and one
  tried to write to GitHub
  ([`runs/T-hbw8/evaluation/summary.md`](../../runs/T-hbw8/evaluation/summary.md)).
- "A model on the 'not used' list, used inside a test run": from the same
  evaluation, found by the plan review of plan v2 (Q-10 on #72).

## Test runs

The output is in [`runs/T-hbw8/test-runs.txt`](../../runs/T-hbw8/test-runs.txt).

## Verdict

Written at the close-out, after the Operator's approval and the review round.

## Resource record

Written at the close-out, after the review round.
