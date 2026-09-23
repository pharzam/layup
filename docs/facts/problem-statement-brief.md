# PROBLEM STATEMENT BRIEF

**Title:** LAYUP  
**Status:** Approved (Revision 6, 2026-09-23)  
**Target Discipline Baseline:** [Armature Engineering Discipline](https://github.com/pharzam/armature)  
**Stakeholders:** Engineering Leadership, Architecture Board, Product Operations  

> This document states the problem only. It does not describe a solution, and it does not depend on a solution document.

---

## 1. Executive Summary

Current multi-agent software development systems fail to realize expected velocity gains due to seven interdependent structural bottlenecks:

1. **The Human Context-Router Bottleneck:** Clarifying questions generated during execution stall autonomous agents. Humans are forced to serve as manual, synchronous message buses—triaging ambiguities, routing inquiries to specialist sub-agents, evaluating multi-turn responses, and feeding context back into prompts.
2. **Discipline Decay & Non-Deterministic Execution:** Human software engineering teams exhibit natural variance in discipline due to cognitive fatigue, seniority, and personal habits. While artificial agents can theoretically maintain absolute discipline, unconstrained agents default to drift-prone code, weak architectural boundaries, and superficial testing. Autonomous delivery needs discipline that does not depend on the agent or the session. Architectural boundaries, interface contracts, and validation pyramids must stay deterministic across every commit. The discipline baseline is **Armature**.
3. **Late Ambiguity Discovery & Manual Bootstrap:** Gaps in a problem statement are found one at a time during delivery, when each one already blocks an agent. The discipline system of a new project is set up by hand, so placeholders get guesses and gates stay partly unwired.
4. **Unresolved Deadlock:** Role agents disagree, or one step repeats without progress. Nothing stops the loop and nothing diagnoses it, so the stall reaches a human late and without the evidence.
5. **Invisible Cost:** Nobody knows the token count, the latency, and the duration of a task. The budget is a human decision, but no number supports it.
6. **Manual Specification Synthesis:** Humans convert the approved problem statement into requirements and technical specifications by hand. Delivery waits for these documents, and a delivered result cannot be traced to a need.
7. **Fragmented Harness Agents:** Work crosses more than one agent product. The rules are copied into the format of each product, and each copy drifts.

---

## 2. Context & Background

Modern engineering organizations are shifting from isolated developer copilots to multi-role autonomous agent teams with dedicated functions: Product Owner, Systems Architect, Domain Expert, Software Architect, Software Engineer, Software Developer, and QA Engineer. 

While foundation models generate code rapidly in isolation, enterprise software delivery is an exercise in **interface contracts, structural boundaries, and asynchronous state synchronization**. When agents cannot coordinate without a human, more agents multiply coordination overhead, not delivery throughput.

A project also does not stay inside one agent product. Different tasks use different products, and each product has its own rule file, memory, context format, and limits. This document calls such a product a **harness**, or a **harness agent**. Claude Code and Codex are examples. It calls the set of rules and gates in the project repository the **discipline system**, and not a harness.

Armature is a template, not a product. It gives a project a ready-made discipline system: quality gate, decision records, requirement documents, test levels, and an issue-first workflow. It has two limits that are important here:

* Its `‹…›` placeholders must be filled for each project, and its roles must be run. Today a human does the two functions.
* It is domain-free and contains no product code. It thus has no gates for package layout, module boundaries, or interface contracts. These gates depend on the technology stack.

---

## 3. The Core Problems & Root Causes

### Problem 1: The Human-as-a-Message-Bus Bottleneck
* **Symptom:** An agent encounters an architectural tradeoff or requirement ambiguity, halts execution, and generates a blocking ticket or prompt for the Operator.
* **Root Cause:** An agent cannot tell which kind of ambiguity it has (domain/business, architectural boundary, interface contract, or environment/infra). It cannot tell which peer agent owns the answer, or if the question needs a human decision. So the agent sends every question to the Operator.
* **Downstream Friction:** The Operator manually investigates, queries a specialist agent, synthesizes conflicting perspectives, and writes context back into the prompt stream. This collapses throughput from seconds to business days.

### Problem 2: Discipline Decay & Structural Drift
* **Symptom:** Code delivered across sprints or by different agent sessions exhibits divergent package layouts, leaky domain abstractions, missing unit-test boundaries, and inconsistent API contracts.
* **Root Cause:** Free-form text prompts do not bind an agent, and no machine stops a change that breaks a convention. Human teams already struggle to hold uniform discipline across projects; unconstrained agents naturally mirror the lowest common denominator found in training sets.
* **Downstream Friction:** When no machine checks the conventions, senior engineers and architects act as human linters and manual code reviewers, negating productivity gains. Armature alone does not close this gap, because its domain-free gates do not cover package layout, interface boundaries, or test quality.

### Problem 3: Late Ambiguity Discovery & Manual Bootstrap
* **Symptom:** A project starts from a problem statement that has gaps: metrics with no measurement method, terms with two readings, no named technology stack. Each gap becomes a blocking question during delivery. In parallel, the discipline system of the new project is copied and configured by hand.
* **Root Cause:** Nothing checks a problem statement for gaps before delivery starts. People adapt the discipline template by hand to the domain and technology stack of each project. No record shows the evidence for each value.
* **Downstream Friction:** Human questions arrive one at a time, each with a full wait cycle. Placeholder values are guesses. Gates that are not wired pass silently, so conformance values are false.

### Problem 4: Unresolved Deadlock
* **Symptom:** Two role agents give different answers, and neither one gives way. Or one agent repeats a step that fails. The task does not reach its goal, and it does not fail cleanly. A human finds the stall by chance, hours or days later.
* **Root Cause:** The system has no state for "no agreement" and no procedure for it. The number of retries has no limit. No independent party with a fresh context examines the evidence and gives a diagnosis.
* **Downstream Friction:** Tokens and time go into repeated attempts. When the stall reaches the Operator, the evidence is spread across sessions, so the Operator collects it by hand.

### Problem 5: Invisible Cost of Autonomous Work
* **Symptom:** A task consumes an unknown number of tokens, an unknown amount of money, and an unknown quantity of time. Nobody can say which role, task, or retry loop consumes the budget.
* **Root Cause:** No record holds the token count, the latency, and the wall-clock duration of each task. Cost is not part of the acceptance of a result.
* **Downstream Friction:** The budget is a business-forking decision (§6), but the idea owner has no number for it. A cheap path and an expensive path cannot be compared. A retry loop can consume a budget before a human sees it.

### Problem 6: Manual Specification Synthesis
* **Symptom:** After the approval of the problem statement, humans write the requirement documents and the technical specifications by hand. Delivery waits for these documents.
* **Root Cause:** A problem statement is prose. Nothing converts it into identified requirements that a machine can validate and trace back to the text that states them.
* **Downstream Friction:** Each step states a requirement again in different words, so the intent degrades. A delivered result cannot be traced to a need, so the idea owner cannot accept the work requirement by requirement.

### Problem 7: Fragmented Harness Agents
* **Symptom:** Work crosses more than one harness agent. Each one has its own rule file, memory, context format, and limits. The project rules are copied into each format, and each copy drifts.
* **Root Cause:** Rules, context, and task state are kept in the form of one harness agent, and not in one neutral form in the project repository. No procedure gives the same rules to a different harness agent.
* **Downstream Friction:** An independent harness agent cannot check the result of the first one, so a verification stays inside the product that made the error. The project also loses its independence (Invariant 2).

---

## 4. Stakeholder Impact Matrix

| Stakeholder Role                  | Day-to-Day Friction Point                                                                                 | Systemic Breakdown                                                                                                |
| :-------------------------------- | :-------------------------------------------------------------------------------------------------------- | :---------------------------------------------------------------------------------------------------------------- |
| **Idea Owners & Product (PO/PM)** | Requirements stall in discovery loops; intent degrades across manual clarification layers.                | Time-to-market extends; delivered features diverge from original business value.                                  |
| **Architects & Tech Leads**       | Context fatigue from continuous re-review; forced to manually enforce foundational conventions.           | Strategic technical direction is neglected in favor of policing package layout, interfaces, and testing pyramids. |
| **Domain Experts & Analysts**     | Repetitive synthesis of prerequisite domain knowledge across isolated agent tasks.                        | Context remains siloed in ephemeral chat logs rather than codified into structured project state.                 |
| **Engineers (Human & Agent)**     | Idling on blocked threads, reconciling conflicting architectural decisions, and refactoring brittle code. | Velocity degrades into technical debt remediation; confidence in autonomous pipelines plummets.                   |

---

## 5. Cost of Inaction (COI)

* **Compounding Delivery Latency:** Agent speed advantages are negated by asynchronous human wait-times for routine clarifications.
* **Architectural Entropy:** Localized agent optimizations violate global system invariants, creating brittle, unmaintainable systems.
* **Unviable Agent Scaling:** When role handoffs are unstructured and no machine checks the rules, more agents add organizational drag, not delivery capacity.
* **Reliance on "Hero" Reviewers:** Quality remains dependent on who reviews the pull request rather than on a verifiable, automated standard.
* **Non-Reproducible Project Starts:** Each new project has a different, hand-made discipline setup, so results from one project do not transfer to the next.
* **Silent Stalls:** A task that does not converge consumes time and tokens until a human finds it by chance.
* **Unbounded Cost:** Without a record of tokens, latency, and duration, an autonomous loop can consume a budget with no warning, and the next project gets no lesson from it.
* **Lock-In to One Product:** When the rules exist only in the format of one harness agent, the project cannot move to another product, and no independent product can check the work.

---

## 6. Scope & System Invariants

### In Scope
* **Problem Statement Quality:** Detection of the gaps in a problem statement before delivery starts, so that the questions for the idea owner come in one batch and not one at a time.
* **Reproducible Discipline Setup:** A repeatable, evidence-based setup of the Armature baseline for a new project, for the domain and technology stack of that project.
* **Rule Protection:** Protection of the rules from the agents that the rules govern.
* **Stack-Dependent Gates:** Gates for repository layout, interface boundaries, contract checks, and test quality, which the domain-free baseline does not give. The technology stack of a project selects these gates. They add rules; they do not change the baseline rules (Invariant 7).
* **Role Handoffs:** Information that passes between role agents has a form that a machine can validate, based on Armature conventions.
* **Autonomous Clarification:** A question that does not need a human decision gets an accepted answer from the responsible role agent, without a human.
* **Verification on Every Change:** Each change passes the gates for repository layout, interface boundaries, and testing pyramids before it reaches human review or merges. The gates are deterministic where a rule can be checked mechanically (Invariant 6).
* **Human-on-the-Loop:** Humans monitor delivery and give input only at the Human Decision Points below. A rule that a machine can apply (the escalation rule) selects the business-forking decisions, for example budget, legal compliance, and strategic trade-offs.
* **Stall Resolution:** A disagreement between role agents, or a step that repeats without progress, gets a procedure with a limit. An independent examination with a fresh context gives a diagnosis. The project repository keeps the diagnosis and the outcome.
* **Cost Visibility:** Each task gets a record of its token count, its latency, and its wall-clock duration in the project repository.
* **Specification Synthesis:** Derivation of identified requirements and technical specifications from the approved problem statement. Each requirement has a trace to the text that states it, and an acceptance criterion.
* **Harness-Agent Neutrality:** The rules, the context, and the task state stay in a form that does not belong to one harness agent. A second harness agent can do the work, and can check the work of the first one.

### Out of Scope
* Modifying base LLM weights or training custom foundation models.
* Decisions that set the intent of a project: which problem to solve, what counts as success, and the funding. The idea owner makes these decisions.
* Modifying downstream cloud infrastructure providers or hosting platforms.
* Changing the rules of the Armature baseline for one project domain or technology stack.

### Human Decision Points
Humans can monitor a project at any time. Monitoring is not input. A human who reads a PR monitors. A change request, a correction, or a required approval of a PR is unplanned human input. These are the only planned points for human input. Any other human input to a task is unplanned human input.

1. **Intent decisions.** Before delivery starts, the idea owner sets the problem, the success criteria, and the funding.
2. **Problem statement answers.** Before delivery starts, the idea owner answers the gap questions about the problem statement, in one batch.
3. **Planned approval points.** Each planned approval point names its approver: the Operator or the idea owner. A routine human review of each PR is not a planned approval point. The acceptance of each delivered requirement is always a planned approval point.
4. **Escalations.** When the escalation rule selects a decision as business-forking, the agent stops and the idea owner decides. An agent never makes a business-forking decision.
5. **Stalls.** When the stall procedure reaches its limit, the agent stops, packages the evidence and the diagnosis, and sends them to the Operator. The Operator can answer, or can give the task to a different harness agent. A stall is a planned point, so the Stall Rate (§7) is monitored together with the Task Intervention Rate. This prevents interventions that hide in stalls.

### System Invariants
Each solution must obey these constraints.

1. **Git is the system of record.** No project state and no decision is kept only outside the project repository.
2. **The project repository is independent.** A human or a different agent can continue the work without the automation that created the repository.
3. **The agents that do the work cannot change the rules or the gates that check the work.**
4. **No configuration value without evidence.** A value that comes from a guess is a defect.
5. **A check that is not active does not count as passed.**
6. **A deterministic check is preferred to an LLM judgement** where a rule can be checked mechanically (Armature R5).
7. **The project domain changes content, never rules.** The technology stack can add stack-dependent gates, but it cannot remove or weaken a baseline rule.
8. **Armature is used at a pinned, recorded version.**
9. **A harness agent is replaceable.** The rules, the context, and the state of a project do not depend on one harness agent.

---

## 7. Success Criteria & Evaluation

The criteria have two layers. Layer 1 holds the checks that must pass. They come from the System Invariants (§6), and a pilot does not calibrate them. Layer 2 holds the numbers that a pilot calibrates. The Layer 2 numbers are working hypotheses for the evaluation of a pilot. They are not a service-level agreement and not a contract. A value marked *(start value)* is set after the pilot measures the baseline.

### 7.1 Invariant checks — pass or fail

| Check | Today | Requirement | Verification |
| :--- | :--- | :--- | :--- |
| **Missed Escalations** (§6 Decision Point 4) | Not measured. | 0 business-forking decisions made by an agent in the audit sample. | Business-forking decisions that an agent made without an escalation. The audit uses a random sample of agent decisions and the window of the Reversal Rate. |
| **Structural Conformance** (Inv. 5) | Variable; manually checked during PR review. | No PR with a failed gate reaches human review or merges. | Count of PRs with a failed gate that reach human review or merge. The first-attempt pass rate is monitored only, so that no one weakens a rule to get a better value. |
| **Inter-Role Communication Format** | Unstructured chat prompts and markdown notes. | $100\%$ schema-validated state artifacts across all role transitions. | Validated role transitions, divided by all role transitions. |
| **Telemetry Completeness** (§6 Cost Visibility) | Absent. | Every task has a token count, a latency, and a wall-clock duration in the project repository. | Tasks with a complete telemetry record, divided by all tasks. |
| **Stall Diagnosis** (§6 Stall Resolution) | Not recorded. | Every stall has a recorded diagnosis and a recorded outcome. | Count of stalls with no diagnosis record. |
| **Specification Traceability** | Manual, not measured. | Every delivered requirement has an identifier, a trace to the text of the problem statement, and an acceptance criterion. | Requirements with a complete trace, divided by all delivered requirements. |
| **Setup Correctness** (Inv. 4) | Manual setup, not checked. | $100\%$ of new project repositories pass the Armature discipline tests. 0 configuration values without evidence. | Run of the discipline tests of the baseline. Count of values with no citation. |
| **Gate Integrity** (Inv. 3) | Not measured. | $100\%$ detection of known-bad commits. 0 agent writes to rule paths. | A set of known-bad commits against the gates. Audit of the writes to rule paths. |
| **Independence** (Inv. 2) | Not applicable. | The project repository passes its gates without the automation. | Gate run with the automation removed. |
| **Harness-Agent Neutrality** (Inv. 9) | Not applicable. | The same project rules and gates run under $\ge 2$ harness agents. Each change gets $\ge 1$ verification from a harness agent that did not make the change. | Count of harness agents that pass the gate run. Changes with an independent verification, divided by all changes. |
| **Generality** | Not applicable. | The result is shown on $\ge 2$ problem statements with different technology stacks. | Count of pilot problem statements and stacks. |

### 7.2 Calibrated targets — measured trend

A pilot measures each baseline with the current process before a target is judged.

| Dimension | Baseline (estimate) | Target *(calibrated in the pilot)* | Measurement |
| :--- | :--- | :--- | :--- |
| **Delivery Lead Time** | Not measured. | Median $\le 50\%$ *(start value)* of the pilot baseline. | Time from the start of the first task for a requirement to the acceptance of that requirement by the idea owner. The pilot measures the baseline with the current process on the pilot problem statements. |
| **First-Review Acceptance** | Not measured. | The idea owner accepts $\ge 90\%$ *(start value)* of delivered requirements at the first review. | Requirements that the idea owner accepts at the first review, divided by all delivered requirements. |
| **Task Intervention Rate** | Not measured per task. Estimate per question: humans answer $>80\%$ of clarifying questions. | $\le 10\%$ *(start value)* of tasks have unplanned human input. | Tasks with $\ge 1$ unplanned human input, divided by all tasks. An answer to an escalation that the idea owner confirms as business-forking does not count. |
| **Clarification Turnaround (autonomous resolution)** | 4 to 48 hours (human asynchronous cycle). | $\le 120$ seconds *(start value)* at the 95th percentile. | Time from the question to the accepted answer, for questions that are resolved without a human. |
| **Reversal Rate** | Not measured. | $< 5\%$ *(start value)* of agent answers are overturned later. | Overturned agent answers, divided by audited agent answers. The audit uses a random sample and a window of 30 days after the merge. |
| **Stall Rate and Resolution** | Not measured. | $\le 5\%$ *(start value)* of tasks stall. $\ge 90\%$ *(start value)* of stalls close without human input. | Tasks with $\ge 1$ stall, divided by all tasks. Stalls closed without human input, divided by all stalls. A stall is planned human input, so this metric is read together with the Task Intervention Rate. |
| **Cost per Requirement** | Not measured. | Median cost per requirement $\le$ the pilot baseline *(start value)*. | Sum of the token cost of all tasks of a requirement. The pilot measures the baseline cost with the current process. |
| **Early Question Share** | Gaps are found during delivery. | $\ge 80\%$ *(start value)* of all human questions for a project are asked before delivery starts. | Questions before delivery, divided by all human questions for the project. |

---

## 8. Terms

| Term | Meaning |
| :--- | :--- |
| **LAYUP** | The name of this initiative. The name does not select a solution. |
| **PSB** | Problem Statement Brief. A document of this kind. |
| **Armature** | The discipline template that is the baseline. It is used at a pinned version. |
| **Harness**, **harness agent** | An agent product that runs a role agent, for example Claude Code or Codex. It gives the model, the session, the tools, and the context window. In this document, the word "harness" always has this meaning, in the title and in each section. A harness is not the discipline system. |
| **Discipline system** | The set of rules and gates that a machine enforces in a project repository. Armature is its baseline. |
| **Project repository** | The Git repository of one project. It is the system of record. |
| **Role agent** | An autonomous agent with one function, for example one of the functions in §2. This document does not set the list of roles. |
| **Operator** | The human who starts, monitors, and approves. This is the Armature meaning of the word. |
| **Idea owner** | The human who owns the intent of a project: the problem, the success criteria, and the funding. The idea owner answers the gap questions about the problem statement, makes the business-forking decisions, and accepts the delivered requirements. The idea owner and the Operator can be the same person. |
| **Human-on-the-Loop** | Humans monitor delivery and give input only at the Human Decision Points (§6). No human approves each task. |
| **Planned approval point** | A project-level approval step that is recorded in the project repository before delivery starts, for example the approval to start delivery. It is not a step inside each task. |
| **Business-forking decision** | A decision that changes the budget, the legal or compliance position, or the approved intent of a project, or that makes a strategic trade-off between approved goals (for example, scope against date). A trade-off inside the approved intent, for example an architectural trade-off, is not business-forking. |
| **Escalation** | A decision that an agent stops on and sends to the idea owner, because the escalation rule selects it as business-forking. The idea owner confirms each escalation as business-forking or not. |
| **Unplanned human input** | Human input to a task that is not at a Human Decision Point (§6), for example an answer, a correction, a restart, or a change to a gate. An answer to an escalation that the idea owner does not confirm as business-forking is unplanned human input. |
| **Task** | One unit of delivery work with one goal. It starts when an agent takes it. It ends when its result merges or when the task is closed. |
| **Delivery** | The work from the start of the first task to the acceptance of the last requirement. |
| **Requirement** | One need that the problem statement states. The idea owner accepts or rejects each delivered requirement. |
| **Project domain** | The business area of a project, with its vocabulary and requirements. The technology stack is not part of the domain. |
| **Technology stack** | The languages, frameworks, and tools that a project uses. |
| **Rule** | A statement of what the work must satisfy, for example a layout, boundary, contract, or test requirement. |
| **Gate** | A check that applies one or more rules to a change and gives pass or fail. |
| **Content** | The project-specific text and values that the rules check or use: terms, requirements, decisions, and configuration values. |
| **Stall** | A task that does not reach its goal and does not fail cleanly, because role agents do not agree, or because a step repeats without progress. |
| **Telemetry** | The record of the token count, the latency, and the wall-clock duration of a task. |
| **Specification** | The requirement and design documents that a project derives from the approved problem statement. Each requirement in them has an identifier, a trace to the problem statement, and an acceptance criterion. |
