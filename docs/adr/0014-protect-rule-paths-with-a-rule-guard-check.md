# 0014. Protect rule paths with a rule-guard check that requires the Operator's approval

Date: 2026-09-25

## Status

Accepted

## Context

Invariant 3 says the agents that do the work cannot change the rules or the gates that check the work (`F-0001#3`); REQ-003 of `PRD-0001` requires a control the agents cannot pass by themselves, and its criterion requires that an attempted rule change by an agent is refused. [ADR-0011](0011-structure-the-core-engine-as-a-go-cli-over-repository-files.md) named the missing control (a code-owners rule with an approval from an account the agents do not use) and deferred it (O-9: the agents and the Operator pushed with one GitHub account). Question 2 of `PRD-0001` §11.

The panel (`runs/T-7qvc/`) compared: `CODEOWNERS` with a required code-owner review and a second identity (A P1, B 2A, C 2-A); a rule-guard check posted by the App, satisfied by the Operator's approval or signed commit (B 2B, A P2); a push-time forge ruleset (C 2-B); detection only, a post-merge audit (B 2C, C 2-C; A names the audit as the measure of `F-0003#64`, not as a control). Every option needs two identities: with one account the forge cannot tell an agent write from a human write. The Operator decided (2026-09-25, #66): **O-53**, the agents get their own identity (an App or a machine account) and the Operator keeps the human account; **O-60**, no `CODEOWNERS` file: only the kit copy and the facts go into a target, so the control is the App's check; **O-61**, an approval of a rule-path change is a planned approval point (Decision Point 3), not unplanned input.

## Decision

We will protect a target's rule paths with **a required check, `layup/rule-guard`, posted by the LAYUP App** ([ADR-0013](0013-run-stack-gates-through-a-layup-github-app.md)):

1. **Three identities.** Agents push and open pull requests as **a separate agents' App or a machine account, never as the LAYUP App** that posts the verdicts ([ADR-0013](0013-run-stack-gates-through-a-layup-github-app.md)); the Operator's human account never runs an agent. The target's setup record names all three (Invariant 4).
2. **The rule-path list is pinned in LAYUP**, not in the target: `.github/**`, the kit's rule documents (`docs/engineering-discipline.md`, `docs/issue-workflow.md`, `docs/guardrails.md`, `docs/adr/**`), the kit's check scripts and their fixtures, `docs/setup/branch-protection.json`, `docs/setup/steps.tsv`, `docs/setup/stack.tsv` (the target's recorded stack and gate set, ADR-0013 decision 4), and the handoff transition table and kind list of [ADR-0015](0015-roles-are-the-kit-gate-roles-and-a-handoff-is-a-tsv-row.md). A pull request cannot shrink the list.
3. **The guard.** On each `pull_request` event and each `pull_request_review` event (a review submitted, changed or dismissed, by any identity — ADR-0013 decision 1) the job lists the changed paths. If none is a rule path, the check posts `success`. If one is, the check posts `failure` until **the Operator identity's latest review for the checked head `X` is an approval** (ADR-0013 decision 3: the App's own row commit `Y` neither starts a run nor dismisses the approval, and an approval given on `Y` counts for `X`), or a commit signed by the Operator's key (`git verify-commit`); any other new commit dismisses the approval and re-runs the guard; a review event — an approval given, changed to "request changes", or dismissed — starts a run that re-evaluates the guard, so a withdrawn approval turns the check red again. The Operator's identity is the approver that the setup record names for Decision Point 3, "a change to a rule path".
4. **The complement, an audit.** `layup audit rules <checkout>` lists every commit on the default branch that touches a rule path, with its author identity, signer and pull request, into `runs/audit/rule-changes.tsv`; an agent-authored row without a recorded approval is a defect. This is the measure of `F-0003#64`, "0 agent writes to rule paths".

We rejected: **`CODEOWNERS`** (a file LAYUP would write into the target, O-60; and with one account the author is its own owner); **a push-time forge ruleset** (a forge-plan feature that is not general across forges, `F-0003#67`, and it stops an agent from proposing a rule change at all); **detection only** (it does not refuse, so REQ-003's criterion fails).

## Consequences

- REQ-003 is met before the merge; its criterion is testable: an agent pull request that touches a rule path stays red until the Operator approves, and the audit counts zero unapproved agent writes.
- **A setup step is added:** the identities of O-53 (the App installation or the machine account, and the Operator's identity as the approver) are a `human_decision = yes` row of `docs/setup/steps.tsv` for a target; the task that builds `layup setup` adds it (a follow-up of this record).
- The human load is one approval per rule change, counted as a planned approval point (O-61), not as unplanned input; the pilot measures how many.
- **Two limits, stated.** The guard reads the forge's review state, which is forge state; the mirror is the audit row in the target's Git. And the guard depends on the App: while the App is removed, the required check keeps the pull request at "expected" (it blocks), which is the safe direction.
- ADR-0011's Consequences named this control as missing; this record supplies it, and O-9 is closed.
