# 0020. Derive the specification from numbered facts

Date: 2026-09-28

## Status

Proposed

## Context

Specification Synthesis is an In-Scope item: identified requirements and
technical specifications are derived from the approved problem statement, and
each requirement has a trace to the text that states it (`F-0003#51`); the
Specification Traceability measure needs that trace to be checkable
(`F-0003#62`). The previous architecture said only that role agents write the
PRD (finding B4); the vision brief asks for a preliminary design review, a PRD
and a phased plan (`F-0002` §2.1). LAYUP already numbers the facts of its own
problem statement and checks each numbered fact as a byte-exact substring of
the stored file (`docs/setup/setup-check.sh`, check `facts`), and `layup psb
check` finds the gaps of a problem statement before delivery (`F-0003#41`). The
smart-if provider does not generate text and is weak at counting
([`runs/T-hbw8/jev-sources.md`](../../runs/T-hbw8/jev-sources.md)). The options
and the selection are in
[`runs/T-hbw8/selection.md`](../../runs/T-hbw8/selection.md), section 10.

## Decision

1. **Intake numbers the facts.** The gap batch of Intake is the output of
   `layup psb check` (the rule-based gaps G1–G5, with no model call) plus the
   questions of a review session on a harness that reads the problem statement
   for gaps of meaning, checked by a session on a counterpart harness; each
   question names its source. After the batch and the idea owner's answers
   (Decision Points 1 and 2), `layup` stores the problem
   statement byte-identical in the target's `docs/facts/`, with a record that
   numbers each clause as a fact (`F-NNNN#n`), as LAYUP's own facts are stored.
2. **A generated skeleton.** In Design, `layup spec draft` writes one requirement
   row per In-Scope fact, with the fact's words as its text, the fact ID, and an
   empty acceptance criterion. A role session then adds the acceptance criteria,
   splits a row where one fact holds several needs (each child keeps the fact
   ID and a byte-exact quote of its part), and writes the preliminary design
   review, the specification (Architect) and the phased plan with its tasks
   (Plan). The smart-if provider writes nothing.
3. **The trace check, with no model.** `layup spec check TARGET` fails when a
   fact ID does not resolve, a quote is not a byte-exact substring of the stored
   problem statement, an In-Scope fact has no `Must` row, a requirement has no
   acceptance criterion, a task names no requirement, or a specification
   section names no requirement. A phase does not move on a failed check
   ([ADR-0013](0013-orchestrate-the-lifecycle-with-an-external-layup-run.md)).
4. **Agreement of meaning.** A verifier session on a counterpart harness
   ([ADR-0015](0015-route-role-sessions-over-registered-harnesses.md)) judges
   whether each requirement and each specification section agrees with the fact
   it cites. The smart-if provider may flag a row for the verifier; a flag
   never passes a row.
5. **The approvals.** The Architect phase ends at a planned approval point that
   Intake lists; the idea owner accepts each delivered requirement at Decision
   Point 3 ([ADR-0016](0016-keep-the-run-records-in-git-and-tell-agent-from-human.md)).
   An unapproved new intent goes to the idea owner, never to the smart-if.

We reject: the smart-if extracting the requirements (it does not generate, and
a model label is not a trace that a machine can check); and a draft checked only
at the end (a missing trace is found late).

## Consequences

- Every requirement and every task traces to a byte-exact span of the approved
  problem statement, and the check that proves it needs no model.
- A byte-exact quote does not prove that a paraphrase keeps the meaning; the
  counterpart verifier and the idea owner's acceptance carry that residual.
- A coarse In-Scope fact gives a large requirement; the split rule keeps the
  trace when a role session divides it.
- `layup spec draft` and `layup spec check` are new engine commands, planned in
  the implementation plan (#42).
