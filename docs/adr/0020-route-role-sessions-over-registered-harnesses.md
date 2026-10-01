# 0020. Route role sessions over registered harnesses

Date: 2026-09-29

## Status

Accepted

## Context

The PSB asks that a second harness can do the work and check the work of the
first (`F-0003#52`), and that each change gets at least one verification from a
harness that did not make it (`F-0003#66`, Invariant 9). The vision brief asks for
squads from registered harnesses, models allocated by task complexity, and
context routed per invocation (`F-0002` §2.2, §3.1); O-67 keeps these in scope.
The Operator decided that the roles are the seven PSB §2 functions by default,
each ambiguity kind owned by one of them, and that the matrix is a parameter
(O-81).

The deep check found no role per step, no task class, and no author of the prompt
(Fable-M16, Sol-23); context routing by name only (Sol-29); and an owner map with
no evidence (Sol-32). The search found no public pattern for an owner per
ambiguity kind ([`runs/T-hbw8/selection-v2.md`](../../runs/T-hbw8/selection-v2.md)
§1.2); the evaluation gave a step table (Spec Kitty, MetaGPT, Omnigent), a harness
table (Agent Orchestrator), a fresh read-only verifier (no_human), and Shape Up's
hill position as a class. The options compared are in `selection-v2.md`, row 0020.

## Decision

We will route role sessions over registered harnesses:

1. **Roles and steps.** The seven PSB §2 roles and a default step table (step,
   role, tier), which the Operator can replace per target.
2. **The owner map** of the four ambiguity kinds has a default that the Intake
   form shows and the Operator confirms or changes; the confirmed map, with its
   comment ID as evidence, is a register.
3. **The harness register** on the LAYUP host; a probe session per harness at
   Intake; the admitted harnesses in the records.
4. **Admission in code:** a probed harness, a model not on the "not used" list,
   and for a plan review or a verification a harness that is not an author: for a
   plan review, the harnesses of the task's plan sessions; for a verification,
   every harness bound, in the ledger, to a commit in the diff from the base, or whose diff went into an attempt's prompt, or that wrote
   the task's frozen tests. With no admitted harness outside the authors, a
   verification is `not-active`. A harness whose version changed is probed again.
5. **Choice among admitted pairs:** the learned weight first; with none or a tie,
   the smart-if's fit point or the table's order.
6. **The tier** of an implementing task: at its start from its plan (the
   reasoning tier for a large task, an open question or a new interface; the
   execution tier otherwise; `new dependencies` counts as a list too); later from its
   computed progress position.
7. **Context:** the records that each step's row names, selected by their links,
   and a start refused when the estimated size exceeds the model's context.

We reject: a fixed pair per role (drops vision 2.2); one provider judgement over
the whole register (a judgement would decide Invariant 9); an owner map with no
evidence.

## Consequences

- Invariant 9 is checked by code at each verification.
- Context is chosen by link, not by meaning; a record that no link reaches is not
  in the prompt (known limit L-D2 of [`architecture.md`](../architecture.md)).
- A single-harness host cannot merge a change.
