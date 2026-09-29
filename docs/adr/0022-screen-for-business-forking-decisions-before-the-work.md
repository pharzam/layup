# 0022. Screen for business-forking decisions before the work

Date: 2026-09-29

## Status

Proposed

## Context

The PSB says: "When the escalation rule selects a decision as business-forking,
the agent stops and the idea owner decides. An agent never makes a
business-forking decision" (`F-0001#13`); a rule that a machine can apply selects
them (`F-0003#48`); Missed Escalations must be 0 in the audit sample (`#57`); an
answer to an escalation that the idea owner does not confirm as business-forking
is unplanned input (`F-0001#28`). The review of #69 found that an escalation can
come after the work (B8).

The deep check found that the first draft screened after the work, on the
agent's own list, with a floor that claimed to read prose (Sol-8, Fable-M7,
Author-4). The search found no public rule that selects business-forking
*decisions*, and the public pattern of a deterministic floor that a model may
raise and never lower (Governed APA;
[`runs/T-hbw8/selection-v2.md`](../../runs/T-hbw8/selection-v2.md) §1.1, row
0022).

## Decision

We will select candidate business-forking decisions from three sources and
screen them before the work goes on:

1. **The session declares** a candidate: its prompt makes it stop with
   `decision_needed` before a choice on the budget, the legal or compliance
   position, the approved intent, or a trade between approved goals.
2. **The floor** (code) reads only fixed fields and proposed diffs: the plan's
   `new dependencies` field (by manifest identifier, with its alternatives); every
   new dependency in the stack's manifest, unless the idea owner put its
   identifier on the allowed-dependency list; a licence file; the PRD's requirement rows, priorities or
   criteria; a path that the Intake named as intent. It exempts a rendering of an
   approved bet or decision and the specification a bet approves. It runs on
   every plan and every handoff diff.
3. **The four questions** of smart-if point P1, one per axis, on prose: each plan,
   handoff and answer, as often as the frequency parameter says (O-79). Code
   combines them with OR; no answer removes a candidate. A question that needs a
   human is a candidate too.
4. **A selected candidate** stops the task in a wait state and goes to the idea
   owner as one brief with numbered options (from the declaration or the plan's
   alternatives, or written by a Product Owner session). The idea owner answers
   `option <N>; business-forking: yes` or `no`; a "no" is counted as unplanned
   input. The chosen and rejected options are recorded, a dependency by its
   manifest identifier; a later candidate with the identifier of a rejected
   option is a finding for the task, not a new brief. An option that changes a
   requirement, a priority or the band carries bet-form lines that code checks,
   and is rendered at once. Each screen writes one row, selected or not.

We reject: the agent's own list as the only source; a floor that claims to read
prose; a screen that runs only after the pull request exists.

## Consequences

- A decision that a session makes without declaring it, in prose that the four
  questions miss and in no file the floor reads, is found only by the Missed
  Escalations audit; under `off` or `shadow` at P1 and P2, that is every undeclared
  choice in prose (known limit L-E1 of [`architecture.md`](../architecture.md)).
- Each false candidate costs the idea owner one answer, which is counted as
  unplanned input; the calibration of P1 is how that cost goes down.
