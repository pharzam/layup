# F-0001. LAYUP problem statement brief, Revision 6

| Field | Value |
| ------------ | ----- |
| Fact ID | `F-0001` |
| Source | LAYUP Problem Statement Brief, Revision 6, status "Approved (Revision 6, 2026-09-23)". Stakeholders named in the brief: Engineering Leadership, Architecture Board, Product Operations. |
| Collected by | Claude Opus 5.5 (agent), for the Operator, task `T-fvwj` |
| Date collected | 2026-09-23 |
| Origin | File [`problem-statement-brief.md`](problem-statement-brief.md), kept next to this record, byte-identical. SHA-256 `3e96862b2578e74f42bfd24929926d0c64496f1a9a51f78dc8c6ca9ecf394ab8` (also in [`../setup/facts.sha256`](../setup/facts.sha256)). Provenance note: copied from the local research repository `layup-research-governner` at commit `4cdf9d0`, which has no remote. |
| Status | `Raw` |

## Facts as collected

> Each fact below is a byte-exact copy of one part of the source file. For a
> list line, the fact is the line without its list number (`1. `). For a §8
> term, the fact is the whole table row. `docs/setup/setup-check.sh` (check
> `facts`) proves that each fact is a substring of the source file. Text of the
> brief that has no number here is cited by section, for example `F-0001 §7.1`.

### §6 System Invariants (facts 1–9; fact N is Invariant N)

1. **Git is the system of record.** No project state and no decision is kept only outside the project repository.
2. **The project repository is independent.** A human or a different agent can continue the work without the automation that created the repository.
3. **The agents that do the work cannot change the rules or the gates that check the work.**
4. **No configuration value without evidence.** A value that comes from a guess is a defect.
5. **A check that is not active does not count as passed.**
6. **A deterministic check is preferred to an LLM judgement** where a rule can be checked mechanically (Armature R5).
7. **The project domain changes content, never rules.** The technology stack can add stack-dependent gates, but it cannot remove or weaken a baseline rule.
8. **Armature is used at a pinned, recorded version.**
9. **A harness agent is replaceable.** The rules, the context, and the state of a project do not depend on one harness agent.

### §6 Human Decision Points (facts 10–14; fact 9+N is Decision Point N)

10. **Intent decisions.** Before delivery starts, the idea owner sets the problem, the success criteria, and the funding.
11. **Problem statement answers.** Before delivery starts, the idea owner answers the gap questions about the problem statement, in one batch.
12. **Planned approval points.** Each planned approval point names its approver: the Operator or the idea owner. A routine human review of each PR is not a planned approval point. The acceptance of each delivered requirement is always a planned approval point.
13. **Escalations.** When the escalation rule selects a decision as business-forking, the agent stops and the idea owner decides. An agent never makes a business-forking decision.
14. **Stalls.** When the stall procedure reaches its limit, the agent stops, packages the evidence and the diagnosis, and sends them to the Operator. The Operator can answer, or can give the task to a different harness agent. A stall is a planned point, so the Stall Rate (§7) is monitored together with the Task Intervention Rate. This prevents interventions that hide in stalls.

### §8 Terms (facts 15–39, one per table row, in source order)

15. | **LAYUP** | The name of this initiative. The name does not select a solution. |
16. | **PSB** | Problem Statement Brief. A document of this kind. |
17. | **Armature** | The discipline template that is the baseline. It is used at a pinned version. |
18. | **Harness**, **harness agent** | An agent product that runs a role agent, for example Claude Code or Codex. It gives the model, the session, the tools, and the context window. In this document, the word "harness" always has this meaning, in the title and in each section. A harness is not the discipline system. |
19. | **Discipline system** | The set of rules and gates that a machine enforces in a project repository. Armature is its baseline. |
20. | **Project repository** | The Git repository of one project. It is the system of record. |
21. | **Role agent** | An autonomous agent with one function, for example one of the functions in §2. This document does not set the list of roles. |
22. | **Operator** | The human who starts, monitors, and approves. This is the Armature meaning of the word. |
23. | **Idea owner** | The human who owns the intent of a project: the problem, the success criteria, and the funding. The idea owner answers the gap questions about the problem statement, makes the business-forking decisions, and accepts the delivered requirements. The idea owner and the Operator can be the same person. |
24. | **Human-on-the-Loop** | Humans monitor delivery and give input only at the Human Decision Points (§6). No human approves each task. |
25. | **Planned approval point** | A project-level approval step that is recorded in the project repository before delivery starts, for example the approval to start delivery. It is not a step inside each task. |
26. | **Business-forking decision** | A decision that changes the budget, the legal or compliance position, or the approved intent of a project, or that makes a strategic trade-off between approved goals (for example, scope against date). A trade-off inside the approved intent, for example an architectural trade-off, is not business-forking. |
27. | **Escalation** | A decision that an agent stops on and sends to the idea owner, because the escalation rule selects it as business-forking. The idea owner confirms each escalation as business-forking or not. |
28. | **Unplanned human input** | Human input to a task that is not at a Human Decision Point (§6), for example an answer, a correction, a restart, or a change to a gate. An answer to an escalation that the idea owner does not confirm as business-forking is unplanned human input. |
29. | **Task** | One unit of delivery work with one goal. It starts when an agent takes it. It ends when its result merges or when the task is closed. |
30. | **Delivery** | The work from the start of the first task to the acceptance of the last requirement. |
31. | **Requirement** | One need that the problem statement states. The idea owner accepts or rejects each delivered requirement. |
32. | **Project domain** | The business area of a project, with its vocabulary and requirements. The technology stack is not part of the domain. |
33. | **Technology stack** | The languages, frameworks, and tools that a project uses. |
34. | **Rule** | A statement of what the work must satisfy, for example a layout, boundary, contract, or test requirement. |
35. | **Gate** | A check that applies one or more rules to a change and gives pass or fail. |
36. | **Content** | The project-specific text and values that the rules check or use: terms, requirements, decisions, and configuration values. |
37. | **Stall** | A task that does not reach its goal and does not fail cleanly, because role agents do not agree, or because a step repeats without progress. |
38. | **Telemetry** | The record of the token count, the latency, and the wall-clock duration of a task. |
39. | **Specification** | The requirement and design documents that a project derives from the approved problem statement. Each requirement in them has an identifier, a trace to the problem statement, and an acceptance criterion. |

## Notes on capture (optional)

- The Operator's instruction named "17 domain terms". Revision 6 §8 holds 25
  rows; one row ("Harness, harness agent") holds two names. All 25 rows are
  facts 15–39. The difference is finding K-01 in
  [`../setup/record-T-n1hp.md`](../setup/record-T-n1hp.md).
- The brief says of itself: "This document states the problem only. It does not
  describe a solution, and it does not depend on a solution document."
