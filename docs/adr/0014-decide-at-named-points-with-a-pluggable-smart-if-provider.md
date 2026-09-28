# 0014. Decide at named points with a pluggable smart-if provider

Date: 2026-09-28

## Status

Proposed

## Context

The previous architecture sent every question to the Operator (finding A3) and
selected an escalation only after the work was done (finding B8). The Operator
decided that high-level decisions are made by a "smart-if" model (O-67), that
the fallback when it cannot decide is a human (O-71), that the engine checks
make no model call and only the decision component calls one (O-72), that the
provider and its authority are parameters per target (O-78, O-79), and that
"if" means a conditional branch (O-84). The decisions are in
[`runs/T-hbw8/operator-decisions.md`](../../runs/T-hbw8/operator-decisions.md)
and [`runs/T-hbw8/inputs-from-pr-69.md`](../../runs/T-hbw8/inputs-from-pr-69.md).

The first provider is Jev, TypeSafe's System One model. Its documentation
([`runs/T-hbw8/jev-sources.md`](../../runs/T-hbw8/jev-sources.md)) says that it
returns typed answers and probabilities (Choice, Noul, Score), that it does not
write text or explanations, that it is weak at numbers, counting and date
comparison, and that thresholds depend on the consequences and must be tested
on the user's own data. The panel's options and the selection are in
[`runs/T-hbw8/selection.md`](../../runs/T-hbw8/selection.md), sections 2 and 4.

## Decision

We will add one decision component to `layup run`
([ADR-0013](0013-orchestrate-the-lifecycle-with-an-external-layup-run.md)),
with these rules:

1. **A conditional branch, nothing more.** At a named decision point, the flow
   of `layup run` takes one branch or another by the answer, as an `if`
   statement does (O-84). The component never writes text, plans work or edits
   a file.
2. **The order of deciding.** First a deterministic check (Invariant 6,
   `F-0001#6`); the smart-if provider is asked only when that check returns
   "undecided"; a human is asked only at a Human Decision Point or on the
   fallback. Code does each sum, count and date comparison; the provider gets
   only the semantic question.
3. **The named points.** Escalation selection (four literal yes/no questions,
   one per business-forking axis of `F-0001#26`, combined in code by OR,
   [ADR-0018](0018-escalate-before-the-work-and-stop-cost-at-the-band.md));
   question routing by ambiguity kind (domain, architecture boundary, interface
   contract, environment, or none; a table maps the kind to its owner role,
   [ADR-0015](0015-route-role-sessions-over-registered-harnesses.md)); the stall
   action among the allowed actions; the panel disposition (a reasoning session
   writes a panel's synthesis; the provider only scores or selects among the
   written options); the fit of a model and harness among the admitted pairs;
   and over budget, inside the band only. A call at any other point is a defect.
4. **The provider interface.** One Go interface with the three answer shapes
   (a choice from a closed set, the probability of yes, a position on ordered
   levels), each with its probabilities and the model version that answered.
   The provider, its model version (pinned, never a moving alias) and its
   endpoint are parameters of the target (O-78). Jev (`jev-1.13.0`) is the
   first provider.
5. **The authority level.** Per target and per point, a parameter selects
   `shadow` (the provider answers and is recorded, and the deterministic
   default or a human decides), `cautious` (the provider alone may only take
   the safer branch: escalate, stop, route to a human) or `delegate` (the
   provider decides above the threshold). Each threshold is a parameter with
   its source (O-79). The default of every point is `shadow` until a
   retrospective promotes it with the recorded agreement as its evidence
   (`F-0001#4`).
6. **The fallback.** A provider error, a rate-limit error after its retries, a
   state too large for its context, or an answer below the threshold takes the
   human branch, never "pass" (O-71). A target never needs the provider to pass
   its gates (`F-0001#2`).
7. **The record.** Each call adds one row to `runs/<task>/decisions.tsv` on the
   records branch: time, point, the deterministic result, the provider and the
   model version that answered, the question-set version, the answer and its
   probabilities, the threshold and its source, the authority level, the branch
   taken, the fallback reason, and the input tokens.

This decision replaces one sentence of ADR-0011 decision 8. "The engine makes no
model call; judgement stays with the harness agents that call it" becomes: the
engine **checks** (`layup psb check`, `layup setup verify`, `layup gate`,
`layup spec check`) make no model call; judgement stays with the harness agents,
and at the named points of `layup run`, with this decision component. The rest
of decision 8 stays: the interface is plain text and files, so any harness agent
can run the engine (Invariant 9).

We reject: a generative model that proposes a class for the provider to check
(a model call outside the named points); no smart-if on the live path (every
ambiguous question goes to a human again, finding A3); one choice over the whole
harness register (Invariant 9 would depend on a judgement); and the provider
writing a panel synthesis (the documentation says it does not generate text).

## Consequences

- A question that no human needs to decide can reach its owner role without the
  Operator (`F-0003#46`), and an escalation is selected before the work.
- In the first milestone of a target, every point is `shadow`, so the human
  load is higher; the records of that milestone are the evidence that promotes
  a point.
- The PRD's non-functional requirement NFR-005 ("the engine itself makes no
  model call") changes with this record to the engine checks.
- The provider is an external service: its price, its rate limits and its
  availability can change without notice (`jev-sources.md`); the fallback
  keeps the run safe, not fast.
