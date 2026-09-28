# 0019. Learn routing from the records at each retrospective

Date: 2026-09-28

## Status

Proposed

## Context

PSB §5 says that the next project gets no lesson (`F-0003#39`). The vision brief
asks for milestones with retrospectives and lessons (`F-0002` §2.3) and for a
learning loop on routing (`F-0002` §3.3); the Operator keeps both in scope
(O-67), puts rule-path changes and the learning input at the retrospective
after each milestone (O-69), and makes the weight adjustments, their bounds and
their triggers parameters (O-83). No model weights change (`F-0003#53`). The
first-attempt gate pass rate is monitored only, so that no one weakens a rule
(`F-0003#58`). The search found learned routers (RouteLLM, bandit routers) that
choose a model per request from quality and cost on large volumes
([`runs/T-hbw8/selection.md`](../../runs/T-hbw8/selection.md), sections 1 and
9).

## Decision

1. **The retrospective phase.** After the Accept phase of each milestone,
   `layup run`
   ([ADR-0013](0013-orchestrate-the-lifecycle-with-an-external-layup-run.md))
   runs the Retrospective. It is a planned approval point that Intake lists
   before delivery (`F-0001#25`).
2. **The reward, in code.** `layup learn` reads only the records branch
   ([ADR-0016](0016-keep-the-run-records-in-git-and-tell-agent-from-human.md))
   and computes, per row of `routing.tsv`
   ([ADR-0015](0015-route-role-sessions-over-registered-harnesses.md)), a reward
   from the records of the milestone: a requirement accepted at its first review
   counts for the route, a reversal, a stall that needed the Operator and an
   unplanned input count against it, and a task with incomplete telemetry
   counts as missing evidence. The first-attempt gate pass rate is not in the
   reward (`F-0003#58`). The terms and their weights are parameters with a
   cited start value (O-83). The result is `runs/<milestone>/reward.tsv`.
3. **The update.** A routing weight changes only when the evidence count of its
   row reaches the parameter `learn.min_evidence`, and by at most the parameter
   `learn.max_step`. A route that breaks an invariant is not eligible whatever
   its reward. The proposed `routing.tsv` goes into the retrospective's batch
   with the batch of rule-path changes
   ([ADR-0017](0017-put-native-stack-gates-in-the-target.md)), and applies from
   the next milestone only after the approval at that planned point. Between two
   retrospectives the routing table does not change.
4. **Lessons as records.** The retrospective's role session writes
   `runs/<milestone>/lessons.tsv`: the lesson, the records it rests on, the
   proposed change and the file it changes. The smart-if provider may score how
   far a lesson is tied to a record; the reward, not the score, changes a
   routing weight, and a proposed rule change goes into the batch.

We reject: exploration trials on weaker routes in the pilot (they spend budget
before a baseline exists); the smart-if proposing the new weights (no auditable
formula); and frozen weights with no learning (O-67 keeps the loop, and O-83
makes the policy a parameter: an Operator can set `learn.max_step` to 0).

## Consequences

- Each routing change is traceable to reward rows, and each reward row to the
  records of one milestone; the same records give the same update.
- A pilot has few tasks per milestone, so a weight moves slowly; the minimum
  evidence count and the step bound limit the effect of one lucky task, they do
  not remove it.
- The lessons of one target reach the next project only through LAYUP's own
  default parameters and recipes, which change under LAYUP's gate; that path is
  work for the implementation plan.
