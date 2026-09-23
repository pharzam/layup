# F-0003. LAYUP problem statement brief, Revision 6 — numbered facts, part 2

| Field | Value |
| ------------ | ----- |
| Fact ID | `F-0003` |
| Source | LAYUP Problem Statement Brief, Revision 6 — the same document as [`F-0001`](F-0001-layup-problem-statement-brief.md). |
| Collected by | Claude Opus 5.5 (agent), for the Operator, task `T-ertw` |
| Date collected | 2026-09-23 |
| Origin | File [`problem-statement-brief.md`](problem-statement-brief.md), the same file as `F-0001`, byte-identical; SHA-256 `3e96862b2578e74f42bfd24929926d0c64496f1a9a51f78dc8c6ca9ecf394ab8` ([`../setup/facts.sha256`](../setup/facts.sha256)). |
| Status | `Raw` |

## Facts as collected

> This record adds fact IDs to the clauses of the brief that `F-0001` does not
> number, so each requirement can cite `F-0003#n`. `F-0001` is immutable, so the
> numbering is a new record. Fact N is the Nth list line, bullet, or table row of
> its section, in source order: a numbered line without its `N. ` prefix, a bullet
> without its `* ` prefix, or a whole table row. Each fact is a byte-exact
> substring of the source file (check `facts` in `docs/setup/setup-check.sh`).

### §1 The seven problems (facts 1–7)

1. **The Human Context-Router Bottleneck:** Clarifying questions generated during execution stall autonomous agents. Humans are forced to serve as manual, synchronous message buses—triaging ambiguities, routing inquiries to specialist sub-agents, evaluating multi-turn responses, and feeding context back into prompts.
2. **Discipline Decay & Non-Deterministic Execution:** Human software engineering teams exhibit natural variance in discipline due to cognitive fatigue, seniority, and personal habits. While artificial agents can theoretically maintain absolute discipline, unconstrained agents default to drift-prone code, weak architectural boundaries, and superficial testing. Autonomous delivery needs discipline that does not depend on the agent or the session. Architectural boundaries, interface contracts, and validation pyramids must stay deterministic across every commit. The discipline baseline is **Armature**.
3. **Late Ambiguity Discovery & Manual Bootstrap:** Gaps in a problem statement are found one at a time during delivery, when each one already blocks an agent. The discipline system of a new project is set up by hand, so placeholders get guesses and gates stay partly unwired.
4. **Unresolved Deadlock:** Role agents disagree, or one step repeats without progress. Nothing stops the loop and nothing diagnoses it, so the stall reaches a human late and without the evidence.
5. **Invisible Cost:** Nobody knows the token count, the latency, and the duration of a task. The budget is a human decision, but no number supports it.
6. **Manual Specification Synthesis:** Humans convert the approved problem statement into requirements and technical specifications by hand. Delivery waits for these documents, and a delivered result cannot be traced to a need.
7. **Fragmented Harness Agents:** Work crosses more than one agent product. The rules are copied into the format of each product, and each copy drifts.

### §3 Symptom, Root Cause, Downstream Friction (facts 8–28)

8. **Symptom:** An agent encounters an architectural tradeoff or requirement ambiguity, halts execution, and generates a blocking ticket or prompt for the Operator.
9. **Root Cause:** An agent cannot tell which kind of ambiguity it has (domain/business, architectural boundary, interface contract, or environment/infra). It cannot tell which peer agent owns the answer, or if the question needs a human decision. So the agent sends every question to the Operator.
10. **Downstream Friction:** The Operator manually investigates, queries a specialist agent, synthesizes conflicting perspectives, and writes context back into the prompt stream. This collapses throughput from seconds to business days.
11. **Symptom:** Code delivered across sprints or by different agent sessions exhibits divergent package layouts, leaky domain abstractions, missing unit-test boundaries, and inconsistent API contracts.
12. **Root Cause:** Free-form text prompts do not bind an agent, and no machine stops a change that breaks a convention. Human teams already struggle to hold uniform discipline across projects; unconstrained agents naturally mirror the lowest common denominator found in training sets.
13. **Downstream Friction:** When no machine checks the conventions, senior engineers and architects act as human linters and manual code reviewers, negating productivity gains. Armature alone does not close this gap, because its domain-free gates do not cover package layout, interface boundaries, or test quality.
14. **Symptom:** A project starts from a problem statement that has gaps: metrics with no measurement method, terms with two readings, no named technology stack. Each gap becomes a blocking question during delivery. In parallel, the discipline system of the new project is copied and configured by hand.
15. **Root Cause:** Nothing checks a problem statement for gaps before delivery starts. People adapt the discipline template by hand to the domain and technology stack of each project. No record shows the evidence for each value.
16. **Downstream Friction:** Human questions arrive one at a time, each with a full wait cycle. Placeholder values are guesses. Gates that are not wired pass silently, so conformance values are false.
17. **Symptom:** Two role agents give different answers, and neither one gives way. Or one agent repeats a step that fails. The task does not reach its goal, and it does not fail cleanly. A human finds the stall by chance, hours or days later.
18. **Root Cause:** The system has no state for "no agreement" and no procedure for it. The number of retries has no limit. No independent party with a fresh context examines the evidence and gives a diagnosis.
19. **Downstream Friction:** Tokens and time go into repeated attempts. When the stall reaches the Operator, the evidence is spread across sessions, so the Operator collects it by hand.
20. **Symptom:** A task consumes an unknown number of tokens, an unknown amount of money, and an unknown quantity of time. Nobody can say which role, task, or retry loop consumes the budget.
21. **Root Cause:** No record holds the token count, the latency, and the wall-clock duration of each task. Cost is not part of the acceptance of a result.
22. **Downstream Friction:** The budget is a business-forking decision (§6), but the idea owner has no number for it. A cheap path and an expensive path cannot be compared. A retry loop can consume a budget before a human sees it.
23. **Symptom:** After the approval of the problem statement, humans write the requirement documents and the technical specifications by hand. Delivery waits for these documents.
24. **Root Cause:** A problem statement is prose. Nothing converts it into identified requirements that a machine can validate and trace back to the text that states them.
25. **Downstream Friction:** Each step states a requirement again in different words, so the intent degrades. A delivered result cannot be traced to a need, so the idea owner cannot accept the work requirement by requirement.
26. **Symptom:** Work crosses more than one harness agent. Each one has its own rule file, memory, context format, and limits. The project rules are copied into each format, and each copy drifts.
27. **Root Cause:** Rules, context, and task state are kept in the form of one harness agent, and not in one neutral form in the project repository. No procedure gives the same rules to a different harness agent.
28. **Downstream Friction:** An independent harness agent cannot check the result of the first one, so a verification stays inside the product that made the error. The project also loses its independence (Invariant 2).

### §4 Stakeholder rows (facts 29–32)

29. | **Idea Owners & Product (PO/PM)** | Requirements stall in discovery loops; intent degrades across manual clarification layers.                | Time-to-market extends; delivered features diverge from original business value.                                  |
30. | **Architects & Tech Leads**       | Context fatigue from continuous re-review; forced to manually enforce foundational conventions.           | Strategic technical direction is neglected in favor of policing package layout, interfaces, and testing pyramids. |
31. | **Domain Experts & Analysts**     | Repetitive synthesis of prerequisite domain knowledge across isolated agent tasks.                        | Context remains siloed in ephemeral chat logs rather than codified into structured project state.                 |
32. | **Engineers (Human & Agent)**     | Idling on blocked threads, reconciling conflicting architectural decisions, and refactoring brittle code. | Velocity degrades into technical debt remediation; confidence in autonomous pipelines plummets.                   |

### §5 Cost of inaction (facts 33–40)

33. **Compounding Delivery Latency:** Agent speed advantages are negated by asynchronous human wait-times for routine clarifications.
34. **Architectural Entropy:** Localized agent optimizations violate global system invariants, creating brittle, unmaintainable systems.
35. **Unviable Agent Scaling:** When role handoffs are unstructured and no machine checks the rules, more agents add organizational drag, not delivery capacity.
36. **Reliance on "Hero" Reviewers:** Quality remains dependent on who reviews the pull request rather than on a verifiable, automated standard.
37. **Non-Reproducible Project Starts:** Each new project has a different, hand-made discipline setup, so results from one project do not transfer to the next.
38. **Silent Stalls:** A task that does not converge consumes time and tokens until a human finds it by chance.
39. **Unbounded Cost:** Without a record of tokens, latency, and duration, an autonomous loop can consume a budget with no warning, and the next project gets no lesson from it.
40. **Lock-In to One Product:** When the rules exist only in the format of one harness agent, the project cannot move to another product, and no independent product can check the work.

### §6 In Scope (facts 41–52)

41. **Problem Statement Quality:** Detection of the gaps in a problem statement before delivery starts, so that the questions for the idea owner come in one batch and not one at a time.
42. **Reproducible Discipline Setup:** A repeatable, evidence-based setup of the Armature baseline for a new project, for the domain and technology stack of that project.
43. **Rule Protection:** Protection of the rules from the agents that the rules govern.
44. **Stack-Dependent Gates:** Gates for repository layout, interface boundaries, contract checks, and test quality, which the domain-free baseline does not give. The technology stack of a project selects these gates. They add rules; they do not change the baseline rules (Invariant 7).
45. **Role Handoffs:** Information that passes between role agents has a form that a machine can validate, based on Armature conventions.
46. **Autonomous Clarification:** A question that does not need a human decision gets an accepted answer from the responsible role agent, without a human.
47. **Verification on Every Change:** Each change passes the gates for repository layout, interface boundaries, and testing pyramids before it reaches human review or merges. The gates are deterministic where a rule can be checked mechanically (Invariant 6).
48. **Human-on-the-Loop:** Humans monitor delivery and give input only at the Human Decision Points below. A rule that a machine can apply (the escalation rule) selects the business-forking decisions, for example budget, legal compliance, and strategic trade-offs.
49. **Stall Resolution:** A disagreement between role agents, or a step that repeats without progress, gets a procedure with a limit. An independent examination with a fresh context gives a diagnosis. The project repository keeps the diagnosis and the outcome.
50. **Cost Visibility:** Each task gets a record of its token count, its latency, and its wall-clock duration in the project repository.
51. **Specification Synthesis:** Derivation of identified requirements and technical specifications from the approved problem statement. Each requirement has a trace to the text that states it, and an acceptance criterion.
52. **Harness-Agent Neutrality:** The rules, the context, and the task state stay in a form that does not belong to one harness agent. A second harness agent can do the work, and can check the work of the first one.

### §6 Out of Scope (facts 53–56)

53. Modifying base LLM weights or training custom foundation models.
54. Decisions that set the intent of a project: which problem to solve, what counts as success, and the funding. The idea owner makes these decisions.
55. Modifying downstream cloud infrastructure providers or hosting platforms.
56. Changing the rules of the Armature baseline for one project domain or technology stack.

### §7.1 Invariant checks (facts 57–67)

57. | **Missed Escalations** (§6 Decision Point 4) | Not measured. | 0 business-forking decisions made by an agent in the audit sample. | Business-forking decisions that an agent made without an escalation. The audit uses a random sample of agent decisions and the window of the Reversal Rate. |
58. | **Structural Conformance** (Inv. 5) | Variable; manually checked during PR review. | No PR with a failed gate reaches human review or merges. | Count of PRs with a failed gate that reach human review or merge. The first-attempt pass rate is monitored only, so that no one weakens a rule to get a better value. |
59. | **Inter-Role Communication Format** | Unstructured chat prompts and markdown notes. | $100\%$ schema-validated state artifacts across all role transitions. | Validated role transitions, divided by all role transitions. |
60. | **Telemetry Completeness** (§6 Cost Visibility) | Absent. | Every task has a token count, a latency, and a wall-clock duration in the project repository. | Tasks with a complete telemetry record, divided by all tasks. |
61. | **Stall Diagnosis** (§6 Stall Resolution) | Not recorded. | Every stall has a recorded diagnosis and a recorded outcome. | Count of stalls with no diagnosis record. |
62. | **Specification Traceability** | Manual, not measured. | Every delivered requirement has an identifier, a trace to the text of the problem statement, and an acceptance criterion. | Requirements with a complete trace, divided by all delivered requirements. |
63. | **Setup Correctness** (Inv. 4) | Manual setup, not checked. | $100\%$ of new project repositories pass the Armature discipline tests. 0 configuration values without evidence. | Run of the discipline tests of the baseline. Count of values with no citation. |
64. | **Gate Integrity** (Inv. 3) | Not measured. | $100\%$ detection of known-bad commits. 0 agent writes to rule paths. | A set of known-bad commits against the gates. Audit of the writes to rule paths. |
65. | **Independence** (Inv. 2) | Not applicable. | The project repository passes its gates without the automation. | Gate run with the automation removed. |
66. | **Harness-Agent Neutrality** (Inv. 9) | Not applicable. | The same project rules and gates run under $\ge 2$ harness agents. Each change gets $\ge 1$ verification from a harness agent that did not make the change. | Count of harness agents that pass the gate run. Changes with an independent verification, divided by all changes. |
67. | **Generality** | Not applicable. | The result is shown on $\ge 2$ problem statements with different technology stacks. | Count of pilot problem statements and stacks. |

### §7.2 Calibrated targets (facts 68–75)

68. | **Delivery Lead Time** | Not measured. | Median $\le 50\%$ *(start value)* of the pilot baseline. | Time from the start of the first task for a requirement to the acceptance of that requirement by the idea owner. The pilot measures the baseline with the current process on the pilot problem statements. |
69. | **First-Review Acceptance** | Not measured. | The idea owner accepts $\ge 90\%$ *(start value)* of delivered requirements at the first review. | Requirements that the idea owner accepts at the first review, divided by all delivered requirements. |
70. | **Task Intervention Rate** | Not measured per task. Estimate per question: humans answer $>80\%$ of clarifying questions. | $\le 10\%$ *(start value)* of tasks have unplanned human input. | Tasks with $\ge 1$ unplanned human input, divided by all tasks. An answer to an escalation that the idea owner confirms as business-forking does not count. |
71. | **Clarification Turnaround (autonomous resolution)** | 4 to 48 hours (human asynchronous cycle). | $\le 120$ seconds *(start value)* at the 95th percentile. | Time from the question to the accepted answer, for questions that are resolved without a human. |
72. | **Reversal Rate** | Not measured. | $< 5\%$ *(start value)* of agent answers are overturned later. | Overturned agent answers, divided by audited agent answers. The audit uses a random sample and a window of 30 days after the merge. |
73. | **Stall Rate and Resolution** | Not measured. | $\le 5\%$ *(start value)* of tasks stall. $\ge 90\%$ *(start value)* of stalls close without human input. | Tasks with $\ge 1$ stall, divided by all tasks. Stalls closed without human input, divided by all stalls. A stall is planned human input, so this metric is read together with the Task Intervention Rate. |
74. | **Cost per Requirement** | Not measured. | Median cost per requirement $\le$ the pilot baseline *(start value)*. | Sum of the token cost of all tasks of a requirement. The pilot measures the baseline cost with the current process. |
75. | **Early Question Share** | Gaps are found during delivery. | $\ge 80\%$ *(start value)* of all human questions for a project are asked before delivery starts. | Questions before delivery, divided by all human questions for the project. |

## Notes on capture (optional)

- Two records now number the same file: `F-0001` (the invariants, the Human
  Decision Points, the terms) and `F-0003` (§1, §3, §4, §5, the scope lists of
  §6, and §7). Neither supersedes the other. The two bullets of §2 are numbered
  by no record.
