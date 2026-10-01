# 0025. Learn routing from the records at each retrospective

Date: 2026-09-29

## Status

Accepted

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

1. **The order:** the audit and the idea owner's confirmations, then `layup
   learn`, then the lessons, then one brief to the Operator; `learn.trigger` (a
   parameter) says at which retrospectives `layup learn` runs, and each run reads,
   per route, the records since that route's last adopted update.
2. **The reward** is computed in code per implementing route: per task, first-
   review acceptance and first-round verification add; material findings, stalls,
   unplanned human input, confirmed reversals and cost subtract; money and
   wall-clock each against their own median; the mean is taken within one role and
   tier. The term weights are parameters with the Operator's comment as evidence.
   A task of unknown money is left out of the money term; a route with no known
   money gets no upward step;
   the gates' first-attempt pass rate is not a term.
3. **The update** moves each route's weight by a bounded step, only with enough
   tasks behind it; by default it is a proposal that the Operator adopts, and can
   later revert; a parameter can let it apply within its bounds. A route's first
   weight is LAYUP's prior, else the mean of its role and tier's weights, else the
   middle of the bounds. A share
   `learn.explore` of the implementing role's tasks goes to the other pairs in
   turn.
4. **The weight ranks** the admitted pairs in routing; the smart-if's fit point
   or the table's order decides only when no pair has a weight, or on a tie.
5. **Lessons** come from a retrospective session on a harness outside the
   milestone's authors, with their evidence, scoped to the target or to LAYUP. Each
   proposed rule change becomes a verified rule batch task before the brief; a batch that changes a gate kind must fail on every recorded known-bad patch of
   that kind and its own before it merges.
6. **A LAYUP lesson** becomes an issue on LAYUP's repository with no content of
   the target; a LAYUP task decides it under LAYUP's gate; a target starts from the
   defaults of the LAYUP version that its Intake records; `layup run` stops on a
   version that differs, and a bet adopts a new one.

We reject: a change of model weights; the smart-if proposing weights; a lesson
that changes another project without a reviewed LAYUP release; a reward that
rewards a route for an unknown cost.

## Consequences

- Routing changes only at a retrospective, and only within the bounds.
- With `learn.explore` at zero, the weights stop moving once one route leads
  (known limit L-H1); only implementing routes learn (L-H2).
- A route with little evidence keeps its weight; a new target learns from LAYUP's
  priors, not from another target's records.
- The next project gets a lesson only as fast as LAYUP's own gate lands it.
