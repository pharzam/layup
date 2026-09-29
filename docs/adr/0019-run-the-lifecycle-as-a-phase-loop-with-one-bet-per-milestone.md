# 0019. Run the lifecycle as a phase loop with one bet per milestone

Date: 2026-09-29

## Status

Proposed

## Context

The Operator decided that LAYUP advances a project "through distinct phases:
solution design, architecture, planning, implementation, and retrospectives",
and asks for approval "only when necessary (e.g., confirming architecture or core
solutions)" (O-67). The PSB puts human input only at the Human Decision Points,
with the acceptance of each delivered requirement always a planned point
(`F-0001#12`), and asks that each change pass the gates before it "reaches human
review or merges" (`F-0003#47`), and that information between roles have a form
a machine validates (`F-0003#45`).

The deep check found no issue and pull-request lifecycle, so the target's own gate
jobs would reject every agent pull request (Sol-27, Fable-M9); a review request
before the gates (Sol-2); handoff validity as a column that an agent writes
(Sol-21, Fable-M14, Author-8); no path from a question back to the session that
asked it (Sol-9, Author-2); and an architecture approval in each milestone
(Author-19). The evaluation gave the patterns of a transition table and typed
results (AgentPlane, ACS), a draft pull request until verified and a merge at the
verified head (AI-SDLC, Agent Orchestrator), the question answered back into the
asking role's work (Agent Orchestrator), and the appetite, the pitch and the bet
(Shape Up). The options compared are in
[`runs/T-hbw8/selection-v2.md`](../../runs/T-hbw8/selection-v2.md), row 0019.

## Decision

We will run each target's lifecycle as a phase loop:

1. **Phases.** Intake, Scaffold, Shape, Bet, Build, Accept, Retrospective, then
   the next Bet, until every `Must` requirement is accepted; back to Shape only
   when a bet proposes a change of the architecture.
2. **One bet per milestone.** A one-screen brief (problem, appetite, solution,
   rabbit holes, no-gos) that a Product Owner session writes and code checks, and
   one comment by the approver of that point. The bet approves the architecture,
   the priorities and milestones of its requirements, and its verified batches
   by their head SHA; no other approval runs inside the milestone except
   escalations and stalls. What the bet changes (priorities, classes) code writes
   as its own tasks.
3. **Every pull request is a task**, and the task loop gives each step of the
   target's own gate an actor: an issue by `layup run`; a plan and a plan review
   on another harness, posted on the issue in the forms the target's checks
   parse; the developer's test first; a draft pull request with no review
   request; the gates; a verifier on a harness that wrote no commit of the
   change, whose record goes on the issue and sets `layup/verify` at the head; a
   close-out commit that code limits to the task file and the completed log, and
   that carries `layup/verify`; the merge by `layup run` at that head SHA, one
   task at a time; a clean merge of the base carries the verification over with no
   new round.
4. **Handoffs** are typed results checked by code against a transition table:
   schema, artifacts with their hashes, and state that code computes.
5. **Questions** end the asking session's attempt; the owner role answers; the
   next attempt gets the answer; the answer is accepted when that attempt ends
      `completed`, cites the answer ID, and asks no question that cites it; an
   attempt that a question ends counts toward the attempt limit. A new attempt
   starts from the base commit, with the frozen tests applied.
6. **Accept.** The approver of the point (the idea owner by default) accepts or
   rejects each delivered requirement by one comment; a rejection is a need for the next bet.

We reject: an architecture approval as its own step in each milestone (human
load); a review request before the gates and the verification pass; handoff
validity from an agent's field; resuming a harness session across attempts
(state that a clone does not carry).

## Consequences

- The idea owner's planned input is: the Intake answers, one bet per milestone,
  one acceptance per delivered requirement, and one approval per retrospective.
- A question costs a new attempt, so the Clarification Turnaround start value
  (120 seconds) is out of reach while a session start takes minutes (known limit
  L-D1 of [`architecture.md`](../architecture.md)).
- The "one bet per milestone" adds planned approval points, which Intake lists
  before delivery starts (`F-0001#25`); the approval brief asks the Operator to
  accept this.
