# PRD-0001. LAYUP, the orchestrator product: full scope with phases

Date: 2026-09-25

## 1. Header & "In plain terms"

| Field        | Value                                                             |
| ------------ | ----------------------------------------------------------------- |
| PRD ID       | `PRD-0001`                                                        |
| Title        | LAYUP: an orchestrator that sets up a target's discipline and drives its delivery |
| Status       | `Accepted`                                                        |
| Date         | 2026-09-25                                                        |
| Author       | Claude Fable 5.1 (agent), for the Operator, task `T-wjq4`         |
| Derives from | `F-0001` and `F-0003` (the PSB, Revision 6), `F-0004` (the idea owner's answers to the PSB's gap questions). `F-0002` (the vision brief) is input for ADRs, not a source of requirements. |

### In plain terms

> LAYUP takes a problem statement from the person who owns an idea, sets up a
> new project repository with a ready-made engineering discipline (Armature),
> and then drives autonomous agents to deliver the product under that
> discipline. People watch, and give input only at a few planned points. Every
> decision, every cost and every stall is written into the project repository,
> and a second agent product can always check the work of the first.

## 2. Problem & evidence

The PSB states seven problems that stop multi-agent software delivery from
reaching its expected speed (`F-0003#1`–`#7`):

1. Humans serve as the message bus between agents (`F-0003#1`), because an agent cannot tell which kind of ambiguity it has or who owns the answer (`F-0003#9`), so every question goes to the Operator and throughput collapses to business days (`F-0003#10`).
2. Discipline decays across sessions (`F-0003#2`): free-form prompts do not bind an agent and no machine stops a change that breaks a convention (`F-0003#12`), so senior engineers act as human linters (`F-0003#13`).
3. Gaps in a problem statement are found one at a time during delivery, and the discipline of a new project is set up by hand (`F-0003#3`, `#15`), so questions arrive one at a time and placeholder values are guesses (`F-0003#16`).
4. A deadlock has no state and no procedure (`F-0003#4`, `#18`), so tokens and time go into repeated attempts and the evidence is scattered when a human finally sees the stall (`F-0003#19`).
5. The cost of a task is invisible (`F-0003#5`, `#21`), so the budget, a business-forking decision, has no number (`F-0003#22`).
6. Requirements are written by hand from prose (`F-0003#6`, `#24`), so intent degrades and a delivered result cannot be traced to a need (`F-0003#25`).
7. Work crosses several harness agents whose rule copies drift (`F-0003#7`, `#27`), so no independent product can check the work (`F-0003#28`).

The cost of inaction is compounding latency, architectural entropy, unviable
agent scaling, reliance on hero reviewers, non-reproducible project starts,
silent stalls, unbounded cost and lock-in to one product (`F-0003#33`–`#40`).

## 3. Goals & non-goals

**Goals.** The twelve In-Scope items of the PSB (`F-0003#41`–`#52`) and the
Generality check (`F-0003#67`): each is a functional requirement in §6.

**Non-goals.** The four Out-of-Scope items of the PSB (`F-0003#53`–`#56`): each
is a `Won't` row in §6 (REQ-015 to REQ-018), so the boundary is
machine-checkable.

## 4. Personas & jobs-to-be-done

- **Idea owner** (`F-0001#23`): owns the intent — the problem, the success criteria, the funding. Hires LAYUP to turn an approved problem statement into a delivered product, with questions in one batch and acceptance requirement by requirement (`F-0001#11`, `#12`).
- **Operator** (`F-0001#22`): starts, monitors and approves. Hires LAYUP to run delivery with input only at the planned points (`F-0001#24`) and to receive a stall with its evidence and diagnosis (`F-0001#14`).
- **Role agent** (`F-0001#21`): an autonomous agent with one function. Hires LAYUP for validated handoffs (`F-0003#45`), accepted answers to its questions (`F-0003#46`) and deterministic gates on its changes (`F-0003#47`).
- **Harness agent** (`F-0001#18`): the product that runs a role agent. LAYUP keeps the rules, the context and the state in a form that any harness agent can use and check (`F-0003#52`).

## 5. User stories

- As the idea owner, I get every gap question about my problem statement in one batch before delivery starts, so that delivery is not blocked one question at a time (`F-0003#41`, `F-0001#11`).
- As the Operator, I get a new project repository whose discipline is set up with evidence for every value, so that the setup is reproducible and no value is a guess (`F-0003#42`, `F-0001#4`).
- As the Operator, I know that the agents cannot change the rules or the gates that check their work (`F-0003#43`, `F-0001#3`).
- As a role agent, my change is checked by the gates of the target's stack before a human sees it (`F-0003#44`, `#47`).
- As a role agent, what I hand to the next role is a record a machine validates (`F-0003#45`).
- As a role agent, a question that needs no human decision gets an accepted answer without a human (`F-0003#46`).
- As the Operator, only business-forking decisions and stalls reach me, and every other input is counted as unplanned (`F-0003#48`, `F-0001#13`, `#28`).
- As the Operator, a stall stops at a limit, gets a fresh-context diagnosis, and reaches me with the evidence (`F-0003#49`, `F-0001#14`).
- As the idea owner, I see the token count, the latency and the duration of every task in the repository, so that the budget has a number (`F-0003#50`, `F-0001#38`).
- As the idea owner, I accept the product requirement by requirement, each traced to my problem statement (`F-0003#51`, `F-0001#31`).
- As the Operator, I can give a task to a different harness agent, and it can check the work of the first (`F-0003#52`, `F-0001#9`).

## 6. Functional requirements

| REQ       | Statement                              | MoSCoW | Phase | Facts             |
| --------- | -------------------------------------- | ------ | ----- | ----------------- |
| REQ-001 | Before delivery starts, `layup psb check` finds the rule gaps of a problem statement, a review session finds its gaps of meaning, and both go to the idea owner as one batch of questions, whose answers are stored as a raw fact. | Must | 1 | F-0003#41, F-0003#14, F-0003#15, F-0001#11 |
| REQ-002 | `layup setup` creates a target repository from the pinned Armature for the target's domain and stack, stops at each human decision, writes the answers to Git, and `layup setup verify` proves the setup with evidence for every value. | Must | 1 | F-0003#42, F-0003#15, F-0003#37, F-0001#4, F-0001#8 |
| REQ-003 | The rules and the gates of a target are protected from the agents that they govern: an agent cannot change a rule path without a control that the agents cannot pass by themselves (ADR-0017). | Must | 2 | F-0003#43, F-0001#3, F-0003#64 |
| REQ-004 | `layup gate` runs the stack-dependent gates of a target (repository layout, interface boundaries, contract checks, test quality) from outside the target and reports `pass`, `fail` or `not-active` per gate; the gates add rules and weaken none. | Must | 1 | F-0003#44, F-0003#11, F-0003#12, F-0003#13, F-0001#7 |
| REQ-005 | Information that passes between role agents is a record in the target that a machine validates against a schema based on Armature conventions. | Must | 2 | F-0003#45, F-0003#59 |
| REQ-006 | A question that does not need a human decision gets an accepted answer from the responsible role agent, without a human, and the answer is recorded in the target. | Must | 3 | F-0003#46, F-0003#8, F-0003#9, F-0003#10 |
| REQ-007 | Each change passes the gates for repository layout, interface boundaries and testing pyramids before it reaches human review or merges; a gate is deterministic where a rule can be checked mechanically. | Must | 1 | F-0003#47, F-0003#58, F-0001#6 |
| REQ-008 | An escalation rule that a machine applies selects the business-forking decisions; the agent stops on one and the idea owner decides; human input outside the Human Decision Points, and an answer to an escalation that the idea owner does not confirm as business-forking, are counted as unplanned. | Must | 3 | F-0003#48, F-0003#57, F-0001#12, F-0001#13, F-0001#24, F-0001#28 |
| REQ-009 | Every stall has a record in the target with its diagnosis and its outcome. | Must | 1 | F-0003#49, F-0003#61, F-0001#37 |
| REQ-010 | A stall (a disagreement between role agents, or a step that repeats without progress) stops at a stated limit, gets a diagnosis from an independent examination with a fresh context, and reaches the Operator with the evidence. | Must | 3 | F-0003#49, F-0003#17, F-0003#18, F-0003#19, F-0001#14 |
| REQ-011 | Every task has a record of its token count, its latency and its wall-clock duration in the target. | Must | 1 | F-0003#50, F-0003#20, F-0003#21, F-0003#22, F-0003#60, F-0001#38 |
| REQ-012 | LAYUP derives identified requirements and technical specifications from an approved problem statement; each requirement has a trace to the text that states it and an acceptance criterion. | Must | 2 | F-0003#51, F-0003#23, F-0003#24, F-0003#25, F-0003#62, F-0001#39 |
| REQ-013 | The rules, the context and the task state of a target stay in a form that belongs to no harness agent, so that a second harness agent can do the work and can check the work of the first. | Must | 4 | F-0003#52, F-0003#26, F-0003#27, F-0003#28, F-0003#66, F-0001#9 |
| REQ-014 | LAYUP shows its result on at least two problem statements with different technology stacks: from each brief it sets up a target and delivers the product. | Must | 4 | F-0003#67, F-0003#42 |
| REQ-015 | LAYUP modifies base LLM weights or trains a custom foundation model. | Won't | — | F-0003#53 |
| REQ-016 | LAYUP makes a decision that sets the intent of a project: which problem to solve, what counts as success, or the funding. | Won't | — | F-0003#54, F-0001#10 |
| REQ-017 | LAYUP modifies a downstream cloud infrastructure provider or hosting platform. | Won't | — | F-0003#55 |
| REQ-018 | LAYUP changes a rule of the Armature baseline for one project domain or technology stack. | Won't | — | F-0003#56, F-0001#7 |

## 7. Non-functional requirements

| REQ       | Statement                              | MoSCoW | Phase | Facts        |
| --------- | -------------------------------------- | ------ | ----- | ------------ |
| NFR-001 | Git is the system of record: no project state and no decision is kept only outside the project repository. | Must | 1 | F-0001#1 |
| NFR-002 | A target repository is independent: it passes its gates without LAYUP, and a human or a different agent continues the work without the automation that created it. | Must | 1 | F-0001#2, F-0003#65 |
| NFR-003 | No configuration value is set without evidence; a value that comes from a guess is a defect. | Must | 1 | F-0001#4, F-0003#63 |
| NFR-004 | A check that is not active does not count as passed. | Must | 1 | F-0001#5 |
| NFR-005 | A deterministic check is preferred to an LLM judgement wherever a rule can be checked mechanically; the engine checks make no model call, and the `layup` process calls a model only through the smart-if provider (the design decision of ADR-0015, which replaces ADR-0011 decision 8; not a clause of the fact). | Must | 1 | F-0001#6 |
| NFR-006 | LAYUP uses Armature at a pinned, recorded version. | Must | 1 | F-0001#8 |
| NFR-007 | LAYUP is written in Go with the standard library only, and calls Git as the `git` program. | Must | 1 | F-0004#1 |

### 7.1 Acceptance criteria

One criterion per requirement: the test a machine or a person runs to accept it.
The IDs are quoted so that this table adds no requirement rows.

| Requirement | Acceptance criterion |
| ----------- | -------------------- |
| `REQ-001` | `layup psb check docs/facts/problem-statement-brief.md` writes the batch `internal/psb/testdata/psb.tsv` byte for byte (`TestGoldenRealPSB`), and the idea owner's answers to it are the raw fact `F-0004`, one fact per question (check `facts`). On a pilot problem statement, the one Intake batch holds both the rule gaps and the gaps of meaning that the review session found. In the pilot, the Early Question Share (`F-0003#75`) is measured against its start value (`F-0004#19`). |
| `REQ-002` | `layup setup` on a pilot problem statement creates a target that passes Armature's discipline tests and `layup setup verify` with zero values without a citation, and an audit of every setup value finds that its cited source supports it (`F-0003#63`, `F-0001#4`); the steps follow the step table of [`docs/spec/setup.md`](../spec/setup.md#the-steps), which derives from `docs/setup/steps.tsv`, and each `human_decision` row's answers are in the target's Git before the next step. |
| `REQ-003` | On a pilot target, an attempted rule change by an agent without the control is refused; an audit of the history shows zero agent writes to a rule path; and a set of known-bad commits against the gates is detected in full (`F-0003#64`). |
| `REQ-004` | `layup gate` on a target of each pilot stack reports one verdict per gate; on each pilot stack, a seeded violation of each of the four gate kinds (layout, interface boundary, contract, test quality) reports `fail`; a gate that did not run reports `not-active`, which never counts as `pass` (reading, ADR-0016: a `pending` kind's `clear` is the pass of the rule "no product path may change while this kind is pending", not of its gate); the target's own baseline gates are unchanged (`F-0001#7`). |
| `REQ-005` | In a pilot task, every role transition carries a record that validates against its schema: validated transitions divided by all transitions equals 1 (`F-0003#59`); and a review finds that the schema encodes the Armature conventions it is based on (the record kinds and fields of Armature), so that a schema that accepts everything does not pass. |
| `REQ-006` | In the pilot, a question that needs no human decision gets an accepted answer without a human, given by the role agent responsible for the question's kind and recorded in the target; the Clarification Turnaround (`F-0003#71`) is measured against its start value (`F-0004#15`); the Reversal Rate (`F-0003#72`) is measured against its start value (`F-0004#16`). |
| `REQ-007` | In the pilot, the count of PRs with a failed gate that reach human review or merge is zero (`F-0003#58`) for each gate kind the target's stack selects; no PR reaches human review or merges with a selected gate that did not run (the reading of `clear` of `REQ-004` applies) (`not-active` counts as not passed, `F-0001#5`); and a repeated run of a gate on the same input gives the same verdict. |
| `REQ-008` | A seeded business-forking decision is selected by the escalation rule, the agent stops, the decision reaches the idea owner, and the idea owner makes it before the task continues; in the audit sample of the pilot, zero business-forking decisions were made by an agent (`F-0003#57`); each human input of the sample is classified as planned or unplanned and the classification is audited against `F-0001#28`; the Task Intervention Rate (`F-0003#70`) is measured against its start value (`F-0004#14`). |
| `REQ-009` | In the pilot, the count of stalls with no diagnosis record is zero (`F-0003#61`); each record holds the diagnosis and the outcome. |
| `REQ-010` | A seeded stall of each kind (a disagreement between role agents; a step that repeats without progress) stops at the stated limit, a fresh-context session of a party that did not do the stalled work writes the diagnosis, and the Operator receives the evidence; the Stall Rate and Resolution (`F-0003#73`) are measured against their start values (`F-0004#17`). |
| `REQ-011` | In the pilot, tasks with a complete telemetry record divided by all tasks equals 1 (`F-0003#60`); a token count a harness does not give is recorded as `not reported` (the ledger's status `unavailable`) and counts as incomplete. |
| `REQ-012` | Every delivered requirement of the pilot has an identifier, a trace to the text of the problem statement and an acceptance criterion: requirements with a complete trace divided by all delivered requirements equals 1 (`F-0003#62`), and each delivered requirement has a technical specification in the target that names the requirement it implements and that a review finds derived from it — the specification's content agrees with the requirement's text (`F-0001#39`). For LAYUP itself, this document is the first instance, and the technical specification in [`docs/spec/`](../spec/README.md) is the second, written one milestone at a time (O-114). |
| `REQ-013` | The same project rules and gates run under at least two harness agents, and each change gets at least one verification from a harness agent that did not make the change (`F-0003#66`). |
| `REQ-014` | For each of at least two problem statements with different technology stacks, `layup setup` created the target from the brief, the role agents that LAYUP orchestrated delivered the product, and the idea owner accepts it (`F-0003#67`). |
| `REQ-015` | No code path of LAYUP modifies base LLM weights or trains a custom foundation model; a code review of each release records it. |
| `REQ-016` | Every intent decision of a pilot is recorded as the idea owner's, at Decision Point 1 (`F-0001#10`); no agent made one. |
| `REQ-017` | No code path of LAYUP calls a cloud provider or hosting platform to modify it; a code review of each release records it. |
| `REQ-018` | A target's baseline rules are byte-identical to the pinned Armature's after setup, except the adapted values the setup records with evidence (`F-0001#7`). |
| `NFR-001` | An audit of a pilot task finds no project state and no decision kept only outside the target's Git: not on the forge, not in a chat, not in a harness's memory, not in a tool's database; every decision an issue cites is an ADR, a task record or a decision note in the tree, or a record on the target's records branch (ADR-0014). |
| `NFR-002` | A gate run on a pilot target with LAYUP removed passes (`F-0003#65`), and a fresh session of a different harness agent, with no LAYUP running, continues one open task of the target from the target's records alone and lands it under the target's gate (`F-0001#2`). |
| `NFR-003` | The count of configuration values with no citation in a target is zero, and an audit of every value finds that its cited source supports it; a value whose citation does not support it is a defect (`F-0003#63`, `F-0001#4`). |
| `NFR-004` | `layup gate` and `layup setup verify` report an inactive check as not passed, and a fixture proves that a check that does not run cannot produce `pass`. Reading (ADR-0016): a `pending` gate kind's job runs the rule "no product path may change while this kind is pending", and `clear` is that rule's pass, not a pass of the kind's gate. |
| `NFR-005` | Each gate verdict is reproducible: two runs on the same input give the same output. The engine checks start no model process and open no connection to a model service; the `layup` process opens a connection to a model service only from the smart-if client, and each such call writes one row to `decisions.tsv`; a call with no row, or at a point that is not named, fails the criterion (ADR-0015). Role sessions are harness processes and are not counted. |
| `NFR-006` | `docs/setup/armature.pin` names the commit and the tree, and check `pin` passes. |
| `NFR-007` | `go build ./...` succeeds with no module outside the standard library (`go list -deps ./...` names none), and Git is called only as the `git` program. |

## 8. Success metrics

Two layers, as the PSB sets them (`F-0001 §7`). Layer 1 holds the checks that
must pass; a pilot does not calibrate them. Layer 2 holds the numbers a pilot
calibrates: the printed value is the working hypothesis, the first pilot
measures the baseline with the current process, and then the idea owner sets
each start value in one batch (`F-0004#11`).

**Layer 1 — pass or fail** (`F-0003#57`–`#67`): missed escalations 0; no PR with
a failed gate reaches human review or merges; 100 % schema-validated role
transitions; every task has a telemetry record; every stall has a diagnosis
record; every delivered requirement has a complete trace; 100 % of new targets
pass the discipline tests with 0 values without evidence; 100 % detection of
known-bad commits and 0 agent writes to rule paths; the target passes its gates
without the automation; the gates run under at least two harness agents with at
least one independent verification per change; at least two pilot problem
statements with different stacks.

**Layer 2 — measured trend**, start values as the idea owner set them:

| Dimension | Start value | Set by | Fact |
| --------- | ----------- | ------ | ---- |
| Delivery Lead Time | median ≤ 50 % of the pilot baseline | the idea owner, after the baseline | `F-0003#68`, `F-0004#12` |
| First-Review Acceptance | ≥ 90 % of delivered requirements accepted at the first review | the same | `F-0003#69`, `F-0004#13` |
| Task Intervention Rate | ≤ 10 % of tasks with unplanned human input | the same | `F-0003#70`, `F-0004#14` |
| Clarification Turnaround | ≤ 120 seconds at the 95th percentile | the same | `F-0003#71`, `F-0004#15` |
| Reversal Rate | < 5 % of audited agent answers overturned within 30 days | the same | `F-0003#72`, `F-0004#16` |
| Stall Rate and Resolution | ≤ 5 % of tasks stall; ≥ 90 % of stalls close without human input | the same | `F-0003#73`, `F-0004#17` |
| Cost per Requirement | median cost per requirement ≤ the pilot baseline; the baseline's own median is the start value | the same | `F-0003#74`, `F-0004#18` |
| Early Question Share | ≥ 80 % of all human questions asked before delivery starts | the same | `F-0003#75`, `F-0004#19` |

## 9. Rollout & phases

The phase tags of §6 and §7. Each phase ships as tasks under the gate of this
repository; phase 1 is the core engine that ADR-0011 structures, which
[`architecture.md`](../architecture.md) extends into the orchestrator `layup run`.

| Phase | Name | Ships |
| ----- | ---- | ----- |
| 1 | The core engine | `layup psb check` (REQ-001: its rule gaps delivered, the review of meaning not yet), `layup setup` as a step runner with no forge call and no role session, and `layup setup verify` (REQ-002; the boundary is in [`docs/spec/setup.md`](../spec/setup.md#the-boundary-of-phase-1)), `layup gate` (REQ-004, REQ-007), the schemas of the telemetry record (REQ-011) and the stall record (REQ-009), whose writer `layup run` comes in a later phase; NFR-001 to NFR-007 hold from the first release. |
| 2 | Specification and handoffs | Specification synthesis (REQ-012), role handoff records (REQ-005), rule protection (REQ-003). |
| 3 | Autonomy | Autonomous clarification (REQ-006), the escalation rule and human-on-the-loop accounting (REQ-008), the stall procedure with fresh-context diagnosis (REQ-010). |
| 4 | Neutrality and the pilot | A second harness agent doing and checking work (REQ-013), the pilot on two problem statements with different stacks (REQ-014), the baseline measurement and the start values of §8. |

The four `Won't` rows (REQ-015 to REQ-018) hold in every phase.

## 10. Risks

- **A value from a guess.** Every setup value needs evidence (NFR-003); the pre-registered rule is [`guardrails.md` §1.1, Inv-4](../guardrails.md#11-layup-system-invariants-psb-6).
- **A check that is not active counted as passed.** NFR-004; [`guardrails.md` §1.1, Inv-5](../guardrails.md#11-layup-system-invariants-psb-6), and the pitfall of a check restored from `main` that runs a new rule only after the merge.
- **The agents change the rules they are checked by.** REQ-003's control is designed ([ADR-0017](../adr/0017-prevent-rule-changes-by-agents.md): no credential in a session, a check before the push, rulesets, rule batches) and not built; [`guardrails.md` §1.1, Inv-3](../guardrails.md#11-layup-system-invariants-psb-6) holds the invariant.
- **Process about process.** A rule for the builders can crowd out the product; [`guardrails.md` §2](../guardrails.md#2-known-pitfalls--the-traps-specific-to-this-domain) records the lesson, and ADR-0012 bounds the work until the first pilot.
- **Two checks that drift.** The setup check is `sh` and the engine is Go (ADR-0011); the engine's setup verification runs the same fixtures as `docs/setup/tests/run.sh`, the defence that [ADR-0011](../adr/0011-structure-the-core-engine-as-a-go-cli-over-repository-files.md)'s Consequences name; [`guardrails.md` §2](../guardrails.md#2-known-pitfalls--the-traps-specific-to-this-domain) holds the related reference-sweep pitfalls.

## 11. Open questions & assumptions

The architecture task (`T-hbw8`, #72) answered questions 1 to 5 in ADR-0013 to ADR-0025; question 6 stays open:

1. **The runner for stack gates.** Answered: the target's own native gates in its own CI, and `layup gate` runs the same commands from outside ([ADR-0016](../adr/0016-put-the-native-stack-gates-in-the-target.md)).
2. **The rule-protection control** (REQ-003). Answered: [ADR-0017](../adr/0017-prevent-rule-changes-by-agents.md).
3. **The list of roles.** Answered: the seven PSB §2 roles by default, a parameter per target (O-81, [ADR-0020](../adr/0020-route-role-sessions-over-registered-harnesses.md)).
4. **The escalation rule's form** (REQ-008). Answered: [ADR-0022](../adr/0022-screen-for-business-forking-decisions-before-the-work.md).
5. **The stall procedure's limits** (REQ-010). Answered: `stall.T`, `stall.N` (O-82) and the other limits of [ADR-0023](../adr/0023-stop-a-stall-at-a-limit-and-diagnose-it-with-a-fresh-context.md).
6. **The pilot's problem statements** (REQ-014): which two, with which stacks; the idea owner selects them (Decision Point 1).

Assumptions: the PSB, Revision 6, is approved and immutable (`F-0001`); the
stack of LAYUP is Go (`F-0004#1`, ADR-0010); a target's stack selects its gates
and is not LAYUP's to choose (`F-0004#1`).

## 12. Requirements traceability matrix

One row per requirement in §6 / §7. `—` marks a cell no record fills today; the
implementation plan (`T-55n2`, [`docs/plan/README.md`](../plan/README.md)) fills the
Task column: in each row, the tasks and the milestones that its table "What
phase 1 proves" names for that requirement, or its milestones (`M2a` to `M4c`)
for a requirement of a later phase. Each delivering task fills the Test column.

| REQ     | Facts                          | Guardrail   | ADR      | Task     | Test |
| ------- | ------------------------------ | ----------- | -------- | -------- | ---- |
| REQ-001 | F-0003#41, F-0003#14, F-0003#15, F-0001#11 | — | ADR-0011 | T-dq05, T-zmj6, T-5zmw, T-evad; `M2c`, `M4b` | TestGoldenRealPSB; check facts (F-0004); TestUsageErrors (cmd/layup, e2e: the usage errors of `layup psb check`); TestPSBCheckOnTheRealProblemStatement (cmd/layup, e2e: the binary gives `psb.tsv`); TestEveryGoldenIsARecordOfTheBlock (internal/psb, integration: each golden is a `psb-gaps` record); TestG1ReadsTheValueOfAStack, TestEdgeCases (internal/psb, unit: O-131 and the edge cases) |
| REQ-002 | F-0003#42, F-0003#15, F-0003#37, F-0001#4, F-0001#8 | §1.1 Inv-4 | ADR-0011, ADR-0016 | T-b97r (its child tasks), T-evad; `M2d`, `M4c` | TestTheEmbeddedTestEntry (internal/catalog, integration: the form of a catalog entry, its files with the module path) |
| REQ-003 | F-0003#43, F-0001#3, F-0003#64 | §1.1 Inv-3 | ADR-0017 | `M2b`, `M2d`, `M2e`, `M2f` | — |
| REQ-004 | F-0003#44, F-0003#11, F-0003#12, F-0003#13, F-0001#7 | §1.1 Inv-7 | ADR-0011, ADR-0016 | T-vk3k (its child tasks), T-d6q5, T-evad; `M2f`, `M4c` | TestTheEmbeddedTestEntry, TestTheSchemaBlocks (internal/catalog, integration: the form of a catalog entry, its kinds and its manifest); TestGateOnAGoRepository (cmd/layup, e2e: one verdict per kind of a Go repository, exit 0 only when each kind passes or is clear) |
| REQ-005 | F-0003#45, F-0003#59            | —           | ADR-0019 | `M2b`, `M2e`, `M2g`, `M4b` | — |
| REQ-006 | F-0003#46, F-0003#8, F-0003#9, F-0003#10 | — | ADR-0019, ADR-0020, ADR-0021 | `M3a`, `M3c`, `M4b` | — |
| REQ-007 | F-0003#47, F-0003#58, F-0001#6  | §1.1 Inv-6  | ADR-0016, ADR-0019 | T-vk3k, T-b97r, T-evad; `M2e`, `M2g` | TestGateOnAGoRepository (cmd/layup, e2e: the repeat rule); TestRunOnARealRepository (internal/gate, integration: the base's gate files judge the head) |
| REQ-008 | F-0003#48, F-0003#57, F-0001#12, F-0001#13, F-0001#24, F-0001#28 | — | ADR-0021, ADR-0022 | `M3a`, `M3b`, `M4b` | — |
| REQ-009 | F-0003#49, F-0003#61, F-0001#37 | —           | ADR-0023 | T-dgy7; `M3d` | — |
| REQ-010 | F-0003#49, F-0003#17, F-0003#18, F-0003#19, F-0001#14 | — | ADR-0023 | `M3d`, `M4b` | — |
| REQ-011 | F-0003#50, F-0003#20, F-0003#21, F-0003#22, F-0003#60, F-0001#38 | — | ADR-0007, ADR-0024 | T-tmhw; `M2b`, `M4b` | — |
| REQ-012 | F-0003#51, F-0003#23, F-0003#24, F-0003#25, F-0003#62, F-0001#39 | — | ADR-0002, ADR-0018 | T-wjq4, T-0drh; `M2c`, `M2e`, `M2g` | prd-lint (this document); the review of `docs/spec/` (#74) |
| REQ-013 | F-0003#52, F-0003#26, F-0003#27, F-0003#28, F-0003#66, F-0001#9 | §1.1 Inv-9 | ADR-0005, ADR-0012, ADR-0015, ADR-0020 | `M2b`, `M4a` | — |
| REQ-014 | F-0003#67, F-0003#42            | —           | —        | `M4c` | — |
| REQ-015 | F-0003#53                       | —           | —        | T-efmy | — |
| REQ-016 | F-0003#54, F-0001#10            | —           | —        | T-evad; `M4c` | — |
| REQ-017 | F-0003#55                       | —           | —        | T-efmy | — |
| REQ-018 | F-0003#56, F-0001#7             | §1.1 Inv-7  | ADR-0011 | T-evad; `M4c` | — |
| NFR-001 | F-0001#1                        | §1.1 Inv-1  | ADR-0011, ADR-0014 | T-b97r, T-evad; `M2a`, `M4c` | — |
| NFR-002 | F-0001#2, F-0003#65             | §1.1 Inv-2  | ADR-0013, ADR-0016 | T-b97r, T-evad; `M4a` | — |
| NFR-003 | F-0001#4, F-0003#63             | §1.1 Inv-4  | —        | T-nfh8, T-b97r, T-evad | check markers; TestHasNamesTheFilesOfTheEntryByTheirPathInIt (internal/catalog, unit: a `catalog` ref resolves to a file of the entry) |
| NFR-004 | F-0001#5                        | §1.1 Inv-5  | ADR-0011 | T-vk3k, T-b97r | TestExitCode (internal/cli, unit: a check that did not run never gives 0); TestGateNeverPassesACheckThatDidNotRun (cmd/layup, e2e: the fixture of item 5 of its section); TestACheckThatDidNotRunNeverPasses (internal/gate, unit); TestPSBCheckGivesTwoWhenItCannotWriteTheTable (internal/cli, unit) and TestPSBCheckWithAReadOnlyStandardOutput (cmd/layup, e2e): a gap table that is not written never gives 0 |
| NFR-005 | F-0001#6                        | §1.1 Inv-6  | ADR-0015 | T-2tc2, T-5sgt, T-b97r; `M3a` | TestPackageRules (cmd/layup, integration: no network package); TestInputRule (cmd/layup, integration: no input from the environment); TestVersion (cmd/layup, e2e: two runs give the same bytes); TestGateOnAGoRepository (cmd/layup, e2e: two runs of layup gate give the same bytes); TestPSBCheckOnTheRealProblemStatement (cmd/layup, e2e: two runs of layup psb check give the same bytes, and a run with no environment variable gives them too) |
| NFR-006 | F-0001#8                        | §1.1 Inv-8  | ADR-0009 | T-r7zg, T-b97r | check pin |
| NFR-007 | F-0004#1                        | —           | ADR-0010 | T-mtb9, T-2tc2 | TestPackageRules (cmd/layup, integration; it holds `internal/psb` to its one import, `internal/tsv`, task `T-5zmw`) |

## 13. Change log

| Date       | Change                     | Requirement(s) affected |
| ---------- | -------------------------- | ----------------------- |
| 2026-09-25 | First draft: the full scope of the PSB with four phases (task `T-wjq4`, Operator decision O-15) | REQ-001 to REQ-018, NFR-001 to NFR-007 |
| 2026-09-29 | The architecture (task `T-hbw8`, #72): REQ-001 names the review of meaning; REQ-003 names ADR-0017; NFR-005 and its criterion follow ADR-0015; §9, §10, §11 and §12 name ADR-0013 to ADR-0025 | REQ-001, REQ-003, NFR-005; the ADR column |
| 2026-10-01 | The technical specification of phase 1 (task `T-0drh`, #74, O-114): `REQ-012`'s criterion names `docs/spec/` as LAYUP's own specification; `REQ-002`'s criterion names the step table of `docs/spec/setup.md`; §9 phase 1 states the boundary of `layup setup` and that phase 1 gives the schemas of the telemetry and stall records | REQ-002, REQ-009, REQ-011, REQ-012 |
| 2026-10-01 | The implementation plan (task `T-55n2`, #76, O-121 to O-124): the §12 Task column names the plan's tasks for phase 1 and its milestones for phases 2 to 4 | REQ-001 to REQ-018, NFR-001 to NFR-007 (the Task column), and the Test cell of NFR-007 |
| 2026-10-01 | The test of the package rules (task `T-2tc2`, #79): the §12 Test cells of NFR-005 and NFR-007 name `TestPackageRules` | NFR-005, NFR-007 (the Test cells) |
| 2026-10-02 | The PDR (task `T-4wrw`, #99, decision O-129): the Operator approved this PRD at the commit `b48764f` ([`PDR-0001`](../pdr/PDR-0001.md)) | none (Status `Draft` to `Accepted`, PDR-0001) |
| 2026-10-02 | The command frame (task `T-2yw7`, #80): the §12 Test cells of REQ-001, NFR-004 and NFR-005 name the tests of the frame | REQ-001, NFR-004, NFR-005 (the Test cells) |
| 2026-10-02 | The stack catalog package (task `T-3jpx`, #81): the §12 Test cells of REQ-002, REQ-004 and NFR-003 name the tests of the catalog | REQ-002, REQ-004, NFR-003 (the Test cells) |
| 2026-10-02 | `layup gate` (task `T-5sgt`, #82): the §12 Test cells of REQ-004, REQ-007, NFR-004 and NFR-005 name its tests | REQ-004, REQ-007, NFR-004, NFR-005 (the Test cells) |
| 2026-10-02 | `layup psb check` to its specification (task `T-5zmw`, #83, O-131): the §12 Test cells of REQ-001, NFR-004, NFR-005 and NFR-007 name its tests | REQ-001, NFR-004, NFR-005, NFR-007 (the Test cells) |
