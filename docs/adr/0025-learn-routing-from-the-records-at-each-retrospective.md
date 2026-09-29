# 0025. Learn routing from the records at each retrospective

Date: 2026-09-29

## Status

Proposed

## Context

The vision brief asks for automated retrospectives that commit lessons at the
close of every milestone, and for reward signals that tune routing from objective
completion, verification, efficiency, and human decisions (`F-0002` §2.3, §3.3);
O-67 keeps these in scope. The Operator decided that rule changes happen only at
the retrospective, approved there, and feed learning (O-69), and that the weight
adjustments, their bounds and the trigger are parameters (O-83). The PSB forbids
changing model weights (`F-0003#53`) and notes that "the next project gets no
lesson" (§5).

The deep check found that a learned weight acted only on a tie once a point was
`delegate` (Fable-M17), that the reward had no cost or verification term and
named "test pass rates" against `F-0003#58` (Author-12), and that no lesson
reached the next project (Sol-25, Fable-N13). The search found bounded weights
that a human adopts, clamped offsets from approvals, and reviewed lessons synced
to projects in Git, and no public reward from delivery records
([`runs/T-hbw8/selection-v2.md`](../../runs/T-hbw8/selection-v2.md) §1.4, row
0025).

## Decision

We will learn routing from the records at each retrospective:

1. **The reward** is computed in code by `layup learn`, per route, from terms whose
   weights are parameters: first-review acceptance and first-round verification
   add; material findings, stalls, unplanned human input, confirmed reversals and
   cost subtract, each against the milestone's median. A missing input leaves its
   term out; the gates' first-attempt pass rate is not a term.
2. **The update** moves each route's weight by a bounded step, only with enough
   tasks behind it; by default it is a proposal that the retrospective's approver
   adopts, and a parameter can let it apply within its bounds.
3. **The weight ranks** the admitted pairs in routing; the smart-if's fit point is
   asked only on a tie or with no weight.
4. **Lessons** come from a retrospective session with their evidence, scoped to
   the target or to LAYUP; the approver keeps or drops each in the one
   retrospective comment, together with the rule batch (O-69).
5. **A LAYUP lesson** becomes an issue on LAYUP's repository; a LAYUP task decides
   it under LAYUP's gate; a new target starts from the defaults of the LAYUP
   version that its Intake records.

We reject: a change of model weights; exploration trials before a baseline
exists (they spend budget on weaker routes); the smart-if proposing weights; a
lesson that changes another project without a reviewed LAYUP release.

## Consequences

- Routing changes only at a retrospective, and only within the bounds.
- A route with little evidence keeps its weight; a new target learns from LAYUP's
  priors, not from another target's records.
- The next project gets a lesson only as fast as LAYUP's own gate lands it.
