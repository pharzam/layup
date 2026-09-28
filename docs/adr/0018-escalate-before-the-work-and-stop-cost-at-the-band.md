# 0018. Escalate before the work and stop cost at the band

Date: 2026-09-28

## Status

Proposed

## Context

A business-forking decision stops the agent, and the idea owner decides
(Decision Point 4, `F-0001#13`); business-forking means a change to the
budget, to the legal or compliance position, to the approved intent, or a
strategic trade-off between approved goals (`F-0001#26`). The previous
architecture found such a decision only in a pull request, after the work
(finding B8), and wrote the cost after the work, so a retry loop could spend a
budget before a human saw it (finding B7; PSB Problem 5). The Operator decided
that the idea owner writes a tolerance band before delivery, the smart-if
decides inside it, and outside it the decision escalates (O-68); that the
thresholds, the escalation rules and the screen frequency are parameters (O-79);
that these parameters cannot switch off the PSB invariants (O-84); and that a
harness that does not report tokens may work with a wall-clock limit as the
proxy by default (O-80). The options and the selection are in
[`runs/T-hbw8/selection.md`](../../runs/T-hbw8/selection.md), section 8; the
decisions in
[`runs/T-hbw8/operator-decisions.md`](../../runs/T-hbw8/operator-decisions.md).

## Decision

1. **Screen before the work.** Before a role session starts, and, by default,
   at each handoff on the handoff's list of the choices the role made (the
   parameter `escalation.screen` is `task` or `task-and-handoff`, O-79),
   `layup run` screens each planned decision:
   - **The floor, in code.** A decision that changes a number of the budget
     file, a line of the approved problem statement, or a `Must` requirement is
     business-forking with no model call.
   - **The four questions.** Otherwise, the smart-if provider answers four
     literal yes/no questions, one per axis of `F-0001#26`
     ([ADR-0014](0014-decide-at-named-points-with-a-pluggable-smart-if-provider.md)).
     Code combines them by OR. An answer above the threshold, or in the band
     between the two thresholds of the point, escalates.
   - **The stop.** An escalated decision stops the work that needs it; the idea
     owner's answer on the issue is copied into Git before the work continues
     ([ADR-0016](0016-keep-the-run-records-in-git-and-tell-agent-from-human.md)).
2. **The budget file.** At Intake the idea owner writes `budget.tsv` on the
   records branch: per milestone, the budget `B` and the upper edge of the band
   `U`, and per task an estimate. It is a planned approval point (Decision Point
   1 or 2).
3. **The cost stop, in code.** Before each new action, code sums the prices of
   `telemetry.tsv` and adds a projection of the remaining work. Below `B`, the
   work continues. Between `B` and `U`, the smart-if provider answers one
   semantic question (does the remaining work still serve the approved goal of
   the task); the arithmetic stays in code. At `U`, or with a projection above
   `U`, the work stops and escalates with no model call; the smart-if never
   approves money outside the band (`F-0001#26`). Per task, a session that
   reaches the parameter `cost.task_factor` times its estimate is stopped.
4. **Harnesses with no token report.** By default such a harness works with a
   wall-clock limit as the proxy for its cost stop, and its telemetry is
   incomplete (O-80); the parameter `harness.<id>.paid` can make it ineligible.
   Where a harness or a provider can enforce a spend cap per action, the action
   reserves that cap before it starts, and an action with no enforceable cap is
   recorded as such.

We reject: the harness estimating a price for the smart-if to judge (money and
numbers are the provider's known weak points, and it would judge money); and a
screen only at the pull request (finding B8).

## Consequences

- A business-forking decision is found before the work that depends on it; one
  that appears inside a role session is found at the next handoff, not before
  (a stated limit of the default).
- A sensitive screen lowers missed escalations (`F-0003#57`) and raises the Task
  Intervention Rate (`F-0003#70`): each escalation that the idea owner does not
  confirm as business-forking is unplanned input (`F-0001#28`); the pilot reads
  the two measures together.
- A session can still spend inside one action before the next check; only an
  enforceable per-action cap closes that, and not every harness has one.
- `B`, `U`, `cost.task_factor` and the thresholds are values the idea owner and
  the Operator set, with their sources, in the budget file and the parameter
  register.
