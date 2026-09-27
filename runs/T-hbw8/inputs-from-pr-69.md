# T-hbw8 — the inputs from PR #69, word for word

Copied on 2026-09-27 from the two comments on PR #69 (closed without a merge by
O-72) that issue #72 names as its inputs, so that a fresh session or another
harness reads them from Git (O-73, Invariant 1). Each comment follows as the
forge holds it, between the lines `--- begin` and `--- end`.

**Reading of two counts (R13; plan review of #72, note 1).** The target design
says "B1 to B10" and "C(1) to C(4)". The review has the findings B1 to B9 and a
forge question, and its table C has ten rows. This task answers A1–A3, B1–B9,
the forge question and each of the ten rows of table C.

## 1. The review of #69 against the PSB and the vision brief

Author: pharzam, 2026-09-26T12:49:14Z.

--- begin
# Review of PR #69 against the PSB and the vision brief

**Scope of this review.** References used: the PSB (`docs/facts/problem-statement-brief.md`, F-0001 and F-0003) and the vision brief (`docs/facts/architectural-vision-brief.md`, F-0002). Not used as a reference: the Operator decisions O-52 to O-62, the panel records under `runs/T-7qvc/`, and the PRD decisions. Reviewed: `docs/architecture.md` and ADR-0013 to ADR-0017 at branch `T-0kn4`. Reviewer: Claude Opus 5.5, for the idea owner.

**Verdict.** The architecture is good on gates, stall records and telemetry records. But it contradicts the PSB in three places. It also has no component for five PSB In-Scope items. The largest problem: PSB Problem 1, the human message bus, is the first bottleneck in the PSB, and the architecture does not solve it.

## A. Direct contradictions with the PSB

**A1. A rule-path approval is counted as a "planned approval point".** (ADR-0014 decision 3 and its Consequences)
- The PSB §8 says a planned approval point is "recorded in the project repository before delivery starts … It is not a step inside each task."
- PSB §6 says "a required approval of a PR is unplanned human input".
- PSB §8 term 28 gives "a change to a gate" as an example of unplanned input.
- Result: the architecture moves these approvals out of the Task Intervention Rate. PSB Decision Point 5 names this risk: "interventions that hide".

**A2. The target cannot pass its stack gates without LAYUP.** (Invariant 2, §7.1 "Independence")
- The PSB says: "The project repository passes its gates without the automation."
- In the architecture, stack gates run only through the `layup` binary and the App, from outside the target. The rule-path list and the handoff transition table are also kept in LAYUP, not in the target (ADR-0014 decision 2, ADR-0015 decision 4).
- The ADR calls this a "stated limit": the Operator removes the check names. After that, the target has no stack gates at all.
- PSB Problem 7 says the rules must be "in one neutral form in the project repository".

**A3. Every question goes to the Operator.** (PSB Problem 1, In Scope "Autonomous Clarification")
- ADR-0015 decision 2 gives the `question` kind one route only: any role → operator. The route to a responsible role agent is left to "a later ADR".
- The PSB root cause of Problem 1 says: "So the agent sends every question to the Operator." The architecture builds this failure as the only path.
- The roles are process steps (planner, implementer, reviewer). No role owns an answer by ambiguity kind (domain, architecture boundary, interface contract, environment). The PSB names these four kinds.
- No record supports three §7.2 measures: Clarification Turnaround, Early Question Share and Reversal Rate.

## B. PSB items with no component

**B1. Nothing starts or moves the role agents.**
- The engine makes no model call and has no daemon. The App reacts only to pull-request events.
- So no component starts a role session, selects a harness, or sends a handoff to the next role. A human must do this, and that makes the human the message bus again (PSB §1 item 1). Each restart is unplanned input (PSB §8 term 28).

**B2. The time-limited stall cannot be detected.**
- ADR-0017 sets T = 15 minutes without a progress event. But no process watches the clock, so the stalled agent must open its own stall.
- A retry loop inside one session that never pushes is also not seen.
- This keeps the PSB Problem 4 symptom: "A human finds the stall by chance."

**B3. Requirement acceptance has no record.**
- PSB Decision Point 3 says: "The acceptance of each delivered requirement is always a planned approval point." The architecture has no table and no command for it.
- `telemetry.tsv` has no requirement ID and no price.
- So three §7.2 measures cannot be calculated: Delivery Lead Time, First-Review Acceptance and Cost per Requirement.

**B4. Specification synthesis has no component.** (PSB Problem 6, §7.1 Specification Traceability)
- The PRD is only "written by the role agents".
- No engine command checks the trace from a requirement to the PSB text.

**B5. Verification by a different harness is not enforced.** (Invariant 9, §7.1)
- The PSB requires "≥ 1 verification from a harness agent that did not make the change". The architecture requires a different *model* for the reviewer, not a different harness.
- Nothing solves the root cause of Problem 7: the rules are copied into the rule file of each harness, and each copy drifts.

**B6. Some decisions stay on the forge.** (Invariant 1)
- The plan, the plan review and the review record are issue comments. How they are copied into Git is not decided (ADR-0015 decision 2).
- The rule-guard approval is review state on the forge. The audit copies it into Git only later.

**B7. Cost has no stop before the budget is spent.** (PSB Problem 5)
- A harness or a person writes telemetry after the work. A value can be "not reported".
- The budget trigger of the escalation floor depends on this data.
- The PSB risk stays: "A retry loop can consume a budget before a human sees it."

**B8. An escalation can come after the work.** (Decision Point 4)
- The deterministic floor selects a *change*, which it finds in a pull request. By then, the agent has already done the work.
- The PSB says: "the agent stops and the idea owner decides."

**B9. The source of the stack gates is not stated.**
- `stack.tsv` has a `command` column, but no document says who writes the gate for a new stack.
- The Generality criterion requires two or more stacks.

**Open question.** The PSB names no forge. The runner is a GitHub App only. Is a forge lock-in acceptable, or is the forge part of the technology stack?

## C. Coverage of the vision brief (solution input, not a requirement, as F-0002 says)

| Vision item | In the architecture? |
|---|---|
| 2.1 Derivation of PDR and PRD, phased plan | No component |
| 2.2 Agent squads, model allocation by complexity | No |
| 2.2 Cross-verification by a counterpart harness | Partly: a different model only |
| 2.3 Milestones, retrospectives, lessons | No. The PSB §5 also says "the next project gets no lesson". |
| 3.1 Model, context and solution routing (Jev/Laya) | No |
| 3.2 Blind panel with a hypothesis posture | Partly: one fresh examiner that does not get the performers' reasoning |
| 3.3 Reinforcement learning on routing | No |
| 3.4 Communication through issues | Replaced by Git tables. This agrees with Invariant 1, but it is a deviation from the vision. |
| 3.4 Telemetry per action | Per gate part only |
| 3.5 Circuit breaker, diagnostic package, external answers | Yes |

## D. What agrees with the PSB

- A check that did not run blocks the merge (Invariant 5).
- The required checks are pinned to the App identity, and agents never act as the App (Invariant 3).
- An `unknown` class escalates.
- The handoff schema rejects a bad record: a bad transition, an unmet condition, a missing artifact or an empty record fails.
- A telemetry value that is "not reported" counts as incomplete.

--- end

## 2. The Operator's decisions O-66 to O-75 and the target design

Author: pharzam, 2026-09-27T14:59:48Z.

--- begin
## Operator decisions O-66 to O-75 (R7): the review of #69 against the PSB and the vision brief, and the soft reset

Given in a Claude Code session (Claude Opus 5.5) on 2026-09-26 and 2026-09-27, after the review comment above. The Operator's words are quoted exactly. Where the Operator chose an option that the agent offered, the option text is quoted too. These decisions are the input of the soft reset. They are posted here so that a fresh session or a different harness can read them from Git and the forge, not from a chat.

**O-66: fix findings in place, with no follow-up issues.** The Operator: "opening another issues and following the progress from there take us to the endless loop. And this is not my favor. I want to work on this PR". Review findings are fixed in the change under review. A finding that is not fixed is recorded as a known limit in the document, not as a new issue.

**O-67: LAYUP is a deterministic orchestrator of the whole lifecycle.** The Operator: "the workflow for any new project should strictly follow Armature principles, with LAYUP acting as a deterministic orchestrator driving the entire lifecycle. Rather than just running disjointed checks, it should advance the process through distinct phases: solution design, architecture, planning, implementation, and retrospectives—continuing this loop until working software that genuinely solves the problem is delivered. Along the way, it should solicit user approval only when necessary (e.g., confirming architecture or core solutions) and manage specialized agent squads (counterpart harnesses)." And: "high-level decisions should be driven by a "smart-if" model like JEV to steer the flow—specifically for handling escalations, convening adjudication panels, or deciding whether to proceed with a task that has exceeded its allocated budget." And, on the vision brief items (reinforcement learning, Jev/Laya routing, squads): "while ignoring it might seem acceptable right now, doing so risks compromising the core objective of the project in the long run." Result: the vision brief (`F-0002`) is an input of the architecture, and its items 3.1 (Jev routing), 3.2 (blind panels), 3.3 (the learning loop) and 2.2 (squads) are in scope. Its note "not as requirements" still holds for the PRD.

**O-68: over budget.** Question: "Over budget: the PSB makes a budget change business-forking. How should Jev handle a task past its budget?" The Operator chose **"Tolerance band (Recommended)"**: "Idea owner writes a band in budget.md before delivery; Jev decides inside it; outside it escalates."

**O-69: rule-path changes happen at the retrospective, and feed the learning loop. Supersedes O-61.** Question: "Rule-path changes (A1): when can they happen?" The Operator wrote: "retro, after every milestone donw should be done , and learned in a way help as RL to working progress". Read as: rule-path changes are made only in the retrospective after each milestone, approved at that planned point, and the lessons of the retrospective are an input of the learning loop that tunes the routing. A rule change at any other time counts as unplanned input (`F-0001#28`).

**O-70: one panel for the new product-architecture ADRs.** Question: "The three new ADRs (0018–0020) are product architecture, and bootstrap mode needs a panel. How do we proceed?" The Operator chose **"One panel, all three (Recommended)"**: "One panel for the three ADRs together, then one review round on the PR." The ADR numbers are provisional.

**O-71: the fallback when Jev cannot decide.** Question: "When Jev is unavailable or below threshold, what is the default?" The Operator chose **"Escalate to human (Recommended)"**: "Safe direction; the target continues without Jev (Invariant 2)."

**O-72: soft reset of the architecture (option A).** The agent offered A (soft reset), B (fix #69 in place) and C (delete the repository and start again). The Operator: "A, is ok". Result: #69 closes **without a merge** when the soft reset starts. The facts, the kit copy, the Go stack (ADR-0010), `layup psb check`, `PRD-0001` and ADR-0012 stay. A new `docs/architecture.md` and new ADRs from ADR-0013 are written from the PSB, the vision brief and O-66 to O-75. ADR-0011 decision 8 ("the engine makes no model call") gets a successor: the engine checks make no model call, and only the decision component calls Jev. The panel records under `runs/T-7qvc/` stay as reference material, not as decisions.

**O-73: approval of the architecture, and communication through Git.** The Operator: "JUST after Architecuter is ready should get my approved as comment on issue, Any commiunication should be through git not agent prompt, answer option". Result: the new architecture does not merge before the Operator's approval as a comment on its issue. All questions, decisions and approvals go through Git and issue comments. The agent copies each decision into a record in Git. This is also a design rule for LAYUP: at each human decision point, a person writes an issue comment, and the orchestrator copies it into Git (Invariant 1).

**O-74: remove the Armature template text before the architecture.** The Operator: "the adoptation of Layout from Armature is not what I expected, Agents.md readme.md, docs/engineering-discipline.md and … have a templete part of Armateure or even name Armature and befor Architecture may it should be fixed". Result: issue #70 (`T-745n`) runs before the soft reset.

**O-75: #70 runs in a fresh session, and that session also gets the plan review.** The Operator: "70 is ok, JUST I wanted to run it from another fresh session" and "just approval part is also assign to that session". Result: the fresh session writes the ordered plan (R12). It then gets the plan review from a second fresh session whose model differs from its own, because a session cannot review its own plan. The reviewer sets `Budget maximum` and `Cycle cap`.

### The target design for the soft reset (the agent's proposal of 2026-09-27; O-67 to O-71 accept it)

- **Components.** The engine checks (deterministic, no model call). The orchestrator `layup run TARGET`: a phase state machine that runs in the foreground, starts the squads and reads the handoff rows. The decision component: Jev (Choice, Noul, Score), called only at named decision points. The squad manager: `docs/harnesses.tsv` and `docs/routing.tsv`. The learning loop: a reward from the records, applied to `routing.tsv` at each retrospective, with no change to model weights. The runner (the LAYUP App). The target repository, with its stack gates written into it, so that it passes its gates without LAYUP (Invariant 2).
- **The phase loop.** Intake (PSB and vision brief, gap batch, answers; Decision Points 1 and 2) → Scaffold (`layup setup`, adapted as #70 defines it) → Design (options, blind panel) → Architect (planned approval) → Plan (milestones and tasks) → Implement (the Armature gate per task, with review by a counterpart harness) → Accept (the idea owner accepts each requirement; Decision Point 3) → Retrospective (lessons, reward, rule changes in one batch; planned approval) → the next milestone, until every `Must` requirement is accepted.
- **Order of decisions.** A deterministic check first (Invariant 6), then Jev where the judgement is semantic, then a human only at a decision point. The Jev decision points: escalation selection (the floor, then one Noul for each of the four PSB axes; this replaces the agent's own declaration), question routing by ambiguity kind, the action for a stall (retry, one examiner, a blind panel, or the Operator), the synthesis of a panel, model and harness routing, and over budget (inside the band only). Each Jev call writes a row to `runs/<task>/decisions.tsv`.
- **Limits from the PSB.** Jev never approves money outside the band (`F-0001#26`). Jev never replaces a check that a machine can apply without judgement. Each threshold is a start value with evidence, and the pilot calibrates it. If Jev is not available, the decision goes to a human, never to "pass".
- **The findings of the review above, all fixed in the new architecture:** A1 (by O-69), A2, A3, B1 to B10, C(1) to C(4). The only known limit: GitHub is the only forge for the pilot.

--- end
