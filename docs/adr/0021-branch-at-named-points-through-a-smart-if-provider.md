# 0021. Branch at named points through a smart-if provider

Date: 2026-09-29

## Status

Proposed

## Context

The Operator decided that "high-level decisions should be driven by a
'smart-if' model like JEV to steer the flow" (O-67); that the provider and the
authority level are parameters per target and per point (O-78), as are the
thresholds, the escalation rules and the check frequency (O-79); that a failed or
below-threshold answer goes to a human (O-71); that parameters tune the PSB rules
and never switch one off, and that "smart-if" means a conditional branch (O-84).
The PSB prefers a deterministic check where a rule can be checked mechanically
(`F-0001#6`).

The deep check found mandatory `shadow` against O-78 (Sol-6), routine decisions
put on a human by `shadow` and the fallback (Sol-7), calls at points that the
list did not name (Fable-M13), a `cautious` level with no safer branch (Fable-N7),
and Laya never looked at (Author-13). The search found public patterns for a
classifier with ordered tiers, fail-closed escalation and recorded thresholds
(jevals, JevLoop), a monitor or enforce switch per rule (numbat), and found that
Laya exists ([`runs/T-hbw8/selection-v2.md`](../../runs/T-hbw8/selection-v2.md)
§1.1, row 0021). ADR-0015 lets the `layup` process call a model only here.

## Decision

We will branch at named points through a provider interface:

1. **The provider** is registered on the LAYUP host in a provider register, in
   the harness register's form (kind, endpoint, credential route, model versions,
   size limit, price source), and probed before use; the admitted providers go
   to the records (O-106). The Operator chooses one per target at Intake: Jev,
   Laya, another registered provider, or none. A
   request holds the point, its fixed literal questions, the options and a state
   text built by code; code does all arithmetic, dates and counts.
2. **Five named points:** escalation (P1), question kind and need of a human
   (P2), stall action (P3), routing fit (P4), over budget inside the band (P5).
   A call anywhere else is a defect.
3. **Authority per point:** `off`, `shadow`, `cautious` (only toward the point's
   safer branch; not for P4) or `delegate`. A yes-or-no answer decides yes at
   p ≥ t and no at p ≤ 1 − t, and is undecided in between; a choice or a score
   decides when its top option reaches t. P1 never removes a candidate. Each
   point has a deterministic branch that decides under `off` and `shadow`. The
   Intake form offers `shadow` and says what it leaves open; the Operator
   chooses. A threshold's evidence is the Operator's comment, or a calibration
   record against the point's named ground truth, never against the default;
   for P3 and P4, whose branches are not run under `shadow`, only the Operator's
   comment.
   The model version is pinned; a new version resets `delegate` to `shadow`.
4. **Bounds:** under `off`, `shadow` and `cautious`, a failure takes the
   deterministic branch and is recorded; at a `delegate` point, a failure or an
   undecided answer goes to a human, never to a pass. No parameter turns off a
   PSB rule.
5. **One `decisions.tsv` row per call,** with the state hash, the model version,
   the answers, the tokens and price, the authority, the branch taken and who
   decided.
6. **Parameters** are a register; a value changes by a comment in a fixed form
   on the target's control issue, by the role that its row allows, applied at the
   next step boundary; a change that reaches an open task is input to that task.

We reject: a mandatory `shadow` period (O-78); a generative model on the live
path (it would be a model call that ADR-0015 does not name); a provider that
decides a count or a sum.

## Consequences

- Under `shadow`, the deterministic branches decide; a human is called only at a
  decision point or on a failure, and the load of P1 falls on the floor and on
  the sessions' own declarations.
- The Operator sees, per point, how often the provider agreed before choosing
  `delegate`.
- Laya is usable only after a calibration of its own; its thresholds are not
  Jev's.
