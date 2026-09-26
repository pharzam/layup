# 0016. Escalate by a deterministic floor over the approved intent and a declared class

Date: 2026-09-25

## Status

Accepted

## Context

The PSB reserves business-forking decisions to the idea owner: a rule that a machine can apply selects them, the agent stops, and the idea owner decides (`F-0003#48`, `F-0001#13`, REQ-008); a business-forking decision changes the budget, the legal or compliance position, or the approved intent, or makes a strategic trade-off between approved goals (`F-0001#26`); an answer to an escalation that the idea owner does not confirm as business-forking is unplanned input (`F-0001#28`). REQ-008's criterion requires a seeded fork to be selected, the agent to stop, and the idea owner to make the decision. Question 4 of `PRD-0001` §11.

The panel (`runs/T-7qvc/`) compared: a deterministic floor that reads the approved-intent files and recorded bounds (B 4A, A E2, C 4-B); a class declared by the agent on four axes and enforced by the engine (A E1, B 4B, C 4-A); both stacked, the stricter result winning (B, C 4-C). The floor cannot see a fork no file shows; the declaration can be under-declared. The Operator decided (#66): **O-56**, both stacked; **O-57**, the intent set is the PSB facts, the PRD, and a budget record `docs/setup/budget.md` that the idea owner writes before the pilot; **O-62**, the idea owner approves the trigger list as part of the PDR, and when the rule is unsure it escalates.

## Decision

We will select business-forking decisions by **two rules, the stricter result winning**:

1. **The floor, deterministic.** The **intent set** of a target is: the facts records of its problem statement and of the idea owner's answers to its gap questions (for LAYUP: `F-0001`, `F-0003`, `F-0004`; other facts, such as a vision brief, are not in the set), its PRD (scope, MoSCoW, phases, success metrics), and `docs/setup/budget.md` (the idea owner's limits: tokens and money per requirement and per task, and the dates it commits to). The engine, `layup escalation check <checkout> <change>`, selects a change as business-forking when it touches an intent path; when a decision record drops or defers a `Must` requirement, moves a phase or a committed date, or changes a dependency's licence class; or when the telemetry of a requirement or task passes a budget value ([REQ-011](../prd/PRD-0001-layup.md)). The trigger list is content of the target's PDR that the idea owner accepts (O-62, Decision Point 1).
2. **The declared class.** Every decision record — a plan, a solution note, an ADR — carries `class: business-forking | in-intent` and the four axes `budget`, `legal_or_compliance`, `approved_intent`, `goal_tradeoff`, each `yes`, `no` or `unknown`, with the evidence for each `no`. The engine stops on any `yes` and on any `unknown` (O-62: when unsure, escalate); a record without the fields is refused.
3. **The stop and the record.** A selection writes a row to `runs/<task>/escalations.tsv` (`id, trigger, decision_path, stopped_at_sha, answer_path, confirmed_business_forking, answered_at`) and the App's check `layup/escalation` ([ADR-0013](0013-run-stack-gates-through-a-layup-github-app.md)) posts `failure` on the pull request until the idea owner's answer is a record in the target's Git; then the row closes and the task continues. The idea owner confirms each escalation as business-forking or not; a not-confirmed one is counted as unplanned input (`F-0001#28`).
4. **The audit.** A random sample of agent decisions is read against the intent set in the window of the Reversal Rate (`F-0003#57`); a fork made without an escalation is the failure the criterion of REQ-008 names.

We rejected: **the floor alone** — a fork made in code and in no intent file is missed; **the declaration alone** — an agent can under-declare, and the judgement is not deterministic (Invariant 6); **the agent as owner of the trigger list** — the agent would define its own boundary.

## Consequences

- REQ-008 is testable: a seeded fork of each of the four kinds and one internal architectural trade-off; the trade-off must not escalate (`F-0001#26`).
- A cost, accepted: each `unknown` escalates, and each escalation the idea owner does not confirm raises the Task Intervention Rate; the pilot measures the confirmation ratio and the idea owner may narrow the trigger list in one batch.
- `docs/setup/budget.md` is a new record of a target and of this repository, written by the idea owner before the pilot (O-57); until it exists, a budget fork is caught by the declared class only, and the record says so.
- The decision records of this repository gain the class fields; a target's decision records carry them from the task that builds `layup escalation check` on. The kit's own templates are not changed (O-10).
