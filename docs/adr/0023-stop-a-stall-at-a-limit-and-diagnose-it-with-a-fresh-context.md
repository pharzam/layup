# 0023. Stop a stall at a limit and diagnose it with a fresh context

Date: 2026-09-29

## Status

Proposed

## Context

The PSB asks that "a disagreement between role agents, or a step that repeats
without progress, gets a procedure with a limit", that "an independent
examination with a fresh context gives a diagnosis", and that the repository
keeps the diagnosis and the outcome (`F-0003#49`, `#61`); at the limit, the agent
packages the evidence and the diagnosis for the Operator, who can answer or give
the task to another harness (`F-0001#14`). The Operator decided that `T` and `N`
are parameters, with defaults `T` = 10 minutes (a maximum) and `N` = 1 (O-82). The
vision brief asks for blind panels with a hypothesis posture, and for a circuit
breaker with a diagnostic package and external answers (`F-0002` §3.2, §3.5).

The deep check found that a new record reset the stall clock and that a
disagreement never stalled (Sol-22, Fable-M4, Author-1); that the clock fired on
every wait for a human (Fable-M5); that the Operator could get a stall with no
diagnosis (Fable-M6); that "a panel" had no mechanism (Sol-24, Fable-M18); that
the dead-man job had no run list, credential or package (Sol-4, Author-10); and
that a wrong gate blocked its own fix (Fable-M19). The evaluation gave round and
attempt limits (Paperclip, Gas Town), wait states apart from "no signal" (Agent
Orchestrator), a fresh read-only reviewer (no_human), a deterministic panel
aggregator with a quorum (AgentJury), and Shape Up's circuit breaker and hill. The
options compared are in
[`runs/T-hbw8/selection-v2.md`](../../runs/T-hbw8/selection-v2.md), row 0023.

## Decision

We will stop a stall at a limit and diagnose it with a fresh context:

1. **Progress is computed** from the records, by sets: a round makes progress
   when it closes an unknown that the last round left open, or a test of the
   frozen list (by ID and source hash, frozen at the first valid handoff to the
   verifier) newly passes. The first review sets the baseline. A task is uphill
   until its list is frozen and while it has an open unknown.
2. **Five triggers:** `stall.N` rounds without progress (a repeated finding
   between two roles included); rounds past the target's cycle cap or attempts
   past `stall.attempts`, questions included; no output for `stall.T` (the first
   session gets no more); a required check with no result for `ci.T`; a takeover
   of the lease. Waits for a human do not count; a clean failure and the circuit
   breaker have their own rows and are not stalls.
3. **A fresh diagnosis first:** code builds the evidence package without the
   sessions' reasoning, and an examiner session writes a diagnosis in a fixed
      form before any other action; an examiner that fails writes "diagnosis
   failed", and the package goes to the Operator at once.
4. **The ladder:** a retry with the diagnosis, then a blind panel (members on two
   or more harnesses with a sealed input and a hypothesis posture, a synthesis
      session that merges text options, which AgentJury's code vote cannot, and a
   quorum in code), then the Operator, by smart-if point P3 or its deterministic
   ladder. The panel rung is skipped when fewer than two admitted harnesses are
   free of the diagnosed failure.
5. **The Operator's package** has a fixed answer form (answer, reroute, stop, or
   an external answer; a parameter changes on the control issue); a reroute starts the next attempt on
   the named harness.
6. **One outcome row** per stall.
7. **The circuit breaker** stops a milestone at its cap; one extension only when
   every open task is downhill, as code computes it, and it fits the band; no
   requirement is dropped.
8. **A wrong gate** opens an early retrospective. **The dead-man job** lives in a
   control repository of the Operator, as a separate App with only read and issue
   permissions, reads the forge's time of the last records update, and only adds
   a notice to the target's control issue.

We reject: a clock reset by any new record; a stall that reaches the Operator
without a diagnosis; a panel of one harness; a job inside a target (FT5).

## Consequences

- The Stall Rate and Stall Diagnosis come from the stall and outcome rows.
- A question loop and a review loop stop after `stall.N` rounds without progress.
- The panel needs two admitted harnesses free of the diagnosed failure; with
  fewer, its rung is skipped and the Operator gets the package.
- A scheduled run can be late or dropped, so a dead host's notice can be late or
  missing (known limit L-F1 of [`architecture.md`](../architecture.md)).
- The Operator makes a second, narrow App for the dead-man job.
