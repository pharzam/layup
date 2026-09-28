# 0017. Put native stack gates in the target

Date: 2026-09-28

## Status

Proposed

## Context

A project repository passes its gates without the automation (Invariant 2,
`F-0001#2`; `F-0003#65`). The previous architecture ran the stack gates only
through LAYUP, from outside the target, so a target without LAYUP had no stack
gates (finding A2), and it did not say who writes a gate for a new stack
(finding B9; Generality, `F-0003#67`). It also counted a rule-path approval as a
planned approval point (finding A1), although a planned approval point "is not a
step inside each task" (`F-0001#25`). The Operator decided that LAYUP is only an
external binary and the target's stack and tools never depend on it (O-76,
which supersedes O-10 and O-11 where they put the stack gates outside the
target), and that rule-path changes happen only at the retrospective after each
milestone (O-69). The decisions are in
[`runs/T-hbw8/operator-decisions.md`](../../runs/T-hbw8/operator-decisions.md)
and [`runs/T-hbw8/inputs-from-pr-69.md`](../../runs/T-hbw8/inputs-from-pr-69.md);
the options and the selection in
[`runs/T-hbw8/selection.md`](../../runs/T-hbw8/selection.md), section 7.

## Decision

1. **The gates are the target's own.** The stack gates of a target (repository
   layout, interface boundaries, contract checks, test quality; `F-0003#44`)
   are the target's own native tools and their configuration, run by the
   target's own CI job on every pull request, and required by the target's
   branch protection. Public stack tools come first; a small check script in the
   target's own stack is written only where no public tool covers a gate kind;
   a gate kind with neither is `not-active`, which is not a pass (`F-0001#5`).
   Nothing in the target calls LAYUP, and the target passes its gates with
   LAYUP absent.
2. **Who writes them.** For each stack, LAYUP holds a gate recipe: the gate kind,
   the public tool and its pinned version, the configuration, a known-bad
   fixture per gate kind, and the evidence. At Scaffold, `layup setup` writes the
   recipe's configuration and CI job into the target as part of the pinned
   baseline copy ([`setup/armature.pin`](../setup/armature.pin)). A recipe for
   a new stack is a change to LAYUP, reviewed under LAYUP's gate, with evidence
   that each gate catches its fixture; Generality needs recipes for two stacks
   (`F-0003#67`). A governed coding agent never writes a recipe for its own
   target.
3. **Run from outside too.** `layup gate TARGET` runs the same native commands
   against a checkout of the target and writes the result to the records branch
   ([ADR-0016](0016-keep-the-run-records-in-git-and-tell-agent-from-human.md)).
   It adds no second gate implementation and is not a required check of the
   target. This keeps true [ADR-0012](0012-build-layup-in-bootstrap-mode.md)
   part 6: the pilot "has run that target's gate from outside".
4. **Rule paths change only at the retrospective.** The rule paths of a target
   are its `AGENTS.md`, its `docs/` rules, its gate configuration, its check
   scripts and its CI workflows. During a milestone, a proposed change to a rule
   path is written to the retrospective's batch and not applied. At the
   retrospective, one change holds the batch, and the approver of
   `approvers.tsv` approves it at that planned point, which Intake lists before
   delivery (`F-0001#25`; O-69). `layup run` does not merge a change that
   touches a rule path outside that batch, and the audit reports any such merge
   as unplanned input (`F-0001#28`); it is never counted as a planned approval
   (finding A1).

This decision amends ADR-0011 decisions 3, 4 and 7. Decision 3: the engine also
writes into a target the stack-gate configuration and the CI job of its recipe,
which are the target's own tools; it still writes no LAYUP code, script or
binary. Decision 4: the stack gates are the target's own and run inside it;
`layup gate` runs them from outside as a second observer. Decision 7: in a
target, CI runs the baseline's jobs and the stack-gate job.

We reject: LAYUP holding the gates with its runner as the only enforcer (this is
finding A2); and gates kept outside and exported into the target at each
retrospective (the target has no current gates between two exports).

## Consequences

- A target merges only through its own gates, with or without LAYUP.
- A gate recipe per stack is new LAYUP content with fixtures; the second stack
  is work for the pilot plan.
- With one GitHub account, the target's own rule protection cannot need an
  approval from a second person; LAYUP's merge rule and the audit carry
  Invariant 3 (`F-0001#3`), with the limit that
  [ADR-0016](0016-keep-the-run-records-in-git-and-tell-agent-from-human.md)
  states.
- A wrong gate that blocks all work in a milestone waits for the retrospective
  through a stall; that is a recorded stall, not a hidden rule change.
