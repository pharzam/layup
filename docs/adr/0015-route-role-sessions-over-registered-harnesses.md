# 0015. Route role sessions over registered harnesses

Date: 2026-09-28

## Status

Proposed

## Context

The PSB requires at least one verification from a harness agent that did not
make the change (Invariant 9, `F-0001#9`; `F-0003#66`); the previous
architecture required only a different model (finding B5). The rules were
copied into the rule file of each harness, and each copy drifted (PSB Problem 7,
`F-0003#26`, `#27`). The vision brief asks for squads of counterpart harnesses
and model allocation by complexity (`F-0002` §2.2); the Operator keeps that item
in scope (O-67). The Operator made the role matrix a parameter, with the seven
functions of PSB §2 as the default (O-81), and made harness eligibility a
parameter (O-80). The options and the selection are in
[`runs/T-hbw8/selection.md`](../../runs/T-hbw8/selection.md), section 5.

## Decision

We will route each role session of `layup run`
([ADR-0013](0013-orchestrate-the-lifecycle-with-an-external-layup-run.md))
through three registers on the target's records branch:

1. **The harness register**, `harnesses.tsv`: harness ID, the command line of
   its adapter, its version, the models it serves, whether it reports tokens,
   whether it can enforce a spend cap, its probe result and the date of the
   probe. A harness enters the register only after a recorded probe run.
2. **The role register**, `roles.tsv`: each role, its phase, and the ambiguity
   kinds it answers. The default is the seven functions of PSB §2 (Product
   Owner, Domain Expert, Systems Architect, Software Architect, Software
   Engineer, Software Developer, QA Engineer), with domain questions owned by
   the Domain Expert, architecture-boundary questions by the Systems Architect,
   interface-contract questions by the Software Architect, and environment
   questions by the Software Engineer. The PSB names the seven roles and the
   four ambiguity kinds but maps no kind to a role (PSB §2, Problem 1), so this
   map is a default with no fact behind it, and the panel members gave two
   different maps (`panel-A.md`, `panel-C.md`). The Operator can replace the
   matrix per target (O-81); a question kind with no owner goes to a human.
3. **The routing table**, `routing.tsv`: role, task class, author harness and
   model, verifier harness and model, weight, evidence, and the retrospective
   that set the row
   ([ADR-0019](0019-learn-routing-from-the-records-at-each-retrospective.md)).
4. **Admission in code.** A pair is admitted only when its harness is
   registered and probed, the model is served by it, the role matches, and, for
   a verification, the verifier harness differs from the author harness of the
   change (Invariant 9). The smart-if provider may then score the fit of each
   admitted pair to the task
   ([ADR-0014](0014-decide-at-named-points-with-a-pluggable-smart-if-provider.md));
   code takes the best pair above the threshold and uses the weight to break a
   tie. With fewer than two working harnesses, a verification is `not-active`,
   and the change does not merge (`F-0001#5`).
5. **Verification.** Before a change moves on, a fresh session on the verifier
   harness reads a clean clone, runs the target's gates and writes a
   verification record bound to the change's commit. A verification record
   whose harness equals the author's fails the check.
6. **One neutral rule source.** The target's rules live only in `AGENTS.md` and
   `docs/`. Each harness entry file (for example `CLAUDE.md`) holds only a
   pointer to `AGENTS.md`; a check fails any other content (`F-0003#27`).
7. **Harness eligibility for paid work.** By default a harness that does not
   report tokens or cannot enforce a spend cap may do paid work with a
   wall-clock limit as the proxy, and its telemetry counts as incomplete (O-80);
   a parameter per harness can make it ineligible.

We reject: a fixed pair per role with no smart-if (it drops the allocation by
complexity that O-67 keeps); one smart-if choice over the whole register (a
judgement would decide Invariant 9); and rules copied into each harness file
(the drift of Problem 7).

## Consequences

- Invariant 9 is checked by code on every change, not by a reviewer's memory.
- A target needs two working harnesses before delivery starts; the pilot must
  prove them (ADR-0012 records 18 failed external harness runs).
- Each harness adapter is new code in LAYUP; a harness whose command line
  changes needs a new probe before it is admitted again.
- The role matrix is content of the target, not a rule of LAYUP; a replaced
  matrix is recorded with its source.
