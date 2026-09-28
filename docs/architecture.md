# LAYUP architecture

Status: **Proposed** — it is approved by the Operator as a comment on issue
[#72](https://github.com/pharzam/layup/issues/72) before it merges (O-73).
Task `T-hbw8`, the soft reset of the architecture (O-72). Written from the
problem statement ([`F-0001`](facts/F-0001-layup-problem-statement-brief.md),
[`F-0003`](facts/F-0003-layup-psb-numbered-facts-part-2.md), the idea owner's
answers [`F-0004`](facts/F-0004-psb-gap-answers.md)), the vision brief
([`F-0002`](facts/F-0002-architectural-vision-brief.md)), the requirements
([`PRD-0001`](prd/PRD-0001-layup.md)) and the Operator's decisions O-66 to O-84
([`runs/T-hbw8/inputs-from-pr-69.md`](../runs/T-hbw8/inputs-from-pr-69.md),
[`runs/T-hbw8/operator-decisions.md`](../runs/T-hbw8/operator-decisions.md)).
The options came from one panel
([`runs/T-hbw8/selection.md`](../runs/T-hbw8/selection.md)). The decisions this
document rests on are [ADR-0010](adr/0010-use-go-as-the-technology-stack.md),
[ADR-0011](adr/0011-structure-the-core-engine-as-a-go-cli-over-repository-files.md),
[ADR-0012](adr/0012-build-layup-in-bootstrap-mode.md) and ADR-0013 to ADR-0020
(section 2).

## In plain terms

> LAYUP is one program that a person runs on their own computer. It takes a
> problem statement, makes a new repository for the product, sets up the
> engineering rules and the checks in it, and then runs AI coding agents
> through fixed phases — design, architecture, plan, build, acceptance, lessons —
> until the product is accepted. It asks a person only at planned points. A
> small "smart-if" model helps it choose at a few named branch points; if that
> model is unsure or down, a person decides. Everything it decides is written
> into Git. The new repository never needs LAYUP to check its own work.

## 1. LAYUP and a target

- **LAYUP** is the Go program `layup` (ADR-0010, ADR-0011). It runs outside every
  target and orchestrates it from the problem statement to working software
  (O-67, O-76). This repository, `pharzam/layup`, is LAYUP's own repository.
- **A target** is the repository of one new product that LAYUP delivers. It
  holds its copy of the pinned baseline ([`setup/armature.pin`](setup/armature.pin)),
  its facts, its own stack gates, its product, and a records branch that LAYUP
  writes. The target's stack and tools never depend on LAYUP: it builds, tests
  and passes its gates with LAYUP absent (Invariant 2, O-76).
- **The Operator** runs LAYUP and approves at the planned points; **the idea
  owner** decides the intent, the budget and each escalation (`F-0001#10`–`#14`).

## 2. Components

| Component | What it is | Decision |
| --------- | ---------- | -------- |
| Engine checks | `layup psb check`, `layup setup verify`, `layup gate`, `layup spec check`: deterministic commands; no model call | [ADR-0011](adr/0011-structure-the-core-engine-as-a-go-cli-over-repository-files.md), [ADR-0014](adr/0014-decide-at-named-points-with-a-pluggable-smart-if-provider.md) |
| The orchestrator | `layup run TARGET`: the phase state machine in the foreground, the role sessions, the clock, the heartbeat, the parameter register | [ADR-0013](adr/0013-orchestrate-the-lifecycle-with-an-external-layup-run.md) |
| The smart-if component | A conditional branch at named points; a provider interface (Jev first); authority per point; the fallback to a human | [ADR-0014](adr/0014-decide-at-named-points-with-a-pluggable-smart-if-provider.md) |
| Squads and routing | The harness register, the role register, the routing table; admission in code; verification on a counterpart harness | [ADR-0015](adr/0015-route-role-sessions-over-registered-harnesses.md) |
| Records and identities | The records branch; one writer; the forge copy; agent and human told apart by the App badge | [ADR-0016](adr/0016-keep-the-run-records-in-git-and-tell-agent-from-human.md) |
| The target's gates | Native stack gates in the target's own CI; the same commands from outside with `layup gate`; rule paths changed only at the retrospective | [ADR-0017](adr/0017-put-native-stack-gates-in-the-target.md) |
| Escalation and budget | The screen before the work; the budget band; the cost stop in code | [ADR-0018](adr/0018-escalate-before-the-work-and-stop-cost-at-the-band.md) |
| The learning loop | The reward in code at each retrospective; a bounded routing update; lessons as records | [ADR-0019](adr/0019-learn-routing-from-the-records-at-each-retrospective.md) |
| Specification synthesis | Numbered facts; a generated requirement skeleton; the trace check; a counterpart verifier | [ADR-0020](adr/0020-derive-the-specification-from-numbered-facts.md) |
| The LAYUP GitHub App | Owned by the Operator's account: its user access token for the agents (the badge), its installation token for the orchestrator (the only writer of the records branch) | [ADR-0016](adr/0016-keep-the-run-records-in-git-and-tell-agent-from-human.md) |
| The dead-man job | A scheduled job in LAYUP's own repository that reads each active run's heartbeat | [ADR-0013](adr/0013-orchestrate-the-lifecycle-with-an-external-layup-run.md) |

## 3. The engine commands

| Command | Does | Model call |
| ------- | ---- | ---------- |
| `layup psb check FILE` | Finds the gaps of a problem statement as one batch of questions (delivered, REQ-001) | no |
| `layup setup TARGET` | Creates the target from the pinned baseline and the stack's gate recipe, per `setup/steps.tsv`; stops at each human decision | no |
| `layup setup verify TARGET` | Proves the setup with evidence for every value | no |
| `layup run TARGET` | The phase loop | only through the smart-if component |
| `layup gate TARGET` | Runs the target's own gate commands from outside; writes the result to the records branch | no |
| `layup spec draft TARGET` / `layup spec check TARGET` | Writes the requirement skeleton from the numbered facts; checks the trace | no |
| `layup learn TARGET` | Computes the reward and the proposed routing update at a retrospective | no |
| `layup audit TARGET` | Checks the records branch history, the merges and the rule-path changes | no |

## 4. The phase loop

| Phase | What happens | Exit check | Human decision |
| ----- | ------------ | ---------- | -------------- |
| Intake | The problem statement and the vision brief are stored; `layup psb check` gives the gap batch; the idea owner's answers are stored; the facts are numbered; `approvers.tsv`, `budget.tsv` and the list of planned approval points are written | the facts check; every gap answered or listed | Decision Points 1 and 2 (planned) |
| Scaffold | `layup setup`: the baseline copy, the stack's native gates and CI job, the records branch, the App installed | `layup setup verify` | the setup questions of `steps.tsv` (planned) |
| Design | `layup spec draft`; role sessions write the requirements, the preliminary design review; a panel where ADR-0006 requires one | `layup spec check` | — |
| Architect | The architecture and its decisions | the trace check; the review | the architecture approval (planned) |
| Plan | Milestones and tasks, each task naming its requirements | every `Must` requirement has a task | — |
| Implement | Per task: the screen, the role session, the target's gates, the verification on a counterpart harness, the merge by the orchestrator | the gates and the verification pass | an escalation (Decision Point 4), only when the screen selects one |
| Accept | The idea owner accepts or rejects each delivered requirement | every delivered requirement has a decision | Decision Point 3 (planned) |
| Retrospective | `layup learn`; the lessons; one batch of rule-path and routing changes | the batch approved | the retrospective approval (planned, listed at Intake) |

After a Retrospective, the next milestone starts at Design, until every `Must`
requirement is accepted. A stall in any phase follows ADR-0013 decision 5.

## 5. The order of deciding

At each branch point: **a deterministic check first** (Invariant 6), **then the
smart-if** where the check returns "undecided", **then a human** only at a Human
Decision Point or on the fallback (O-67, O-71, O-84).

| Point | Deterministic check first | Smart-if question | Default authority | Fallback |
| ----- | ------------------------- | ----------------- | ----------------- | -------- |
| Escalation selection | the floor: a changed budget number, a changed line of the problem statement or a changed `Must` row | four yes/no questions, one per axis of `F-0001#26`, OR in code | shadow | the idea owner |
| Question routing | a question that names its kind | one choice: domain, architecture boundary, interface contract, environment, none | shadow | the Operator |
| Stall action | the retry count against `N` | one choice among the allowed actions | shadow | the Operator |
| Panel disposition | every option cites a fact | one score per written option | shadow | the Operator |
| Model and harness fit | admission of the pair (ADR-0015) | one score per admitted pair | shadow | the highest weight, recorded as a hold |
| Over budget | the sum against `B` and `U` | inside the band only: does the remaining work serve the task's goal | shadow | the idea owner |

Every point starts in `shadow` in a new target; a retrospective promotes a point
to `cautious` or `delegate` with the recorded agreement as its evidence
(ADR-0014 decision 5).

## 6. Identities and the forge

| Actor | Acts as | Seen on GitHub as | Allowed to |
| ----- | ------- | ----------------- | ---------- |
| The Operator, the idea owner | their own account (`approvers.tsv`) | the account, no App | decide at a Human Decision Point |
| A role agent | the Operator's account through the App's user access token | the account's avatar with the App's badge; `performed_via_github_app` = the App | push a task branch, open a pull request, comment |
| `layup run` | the App's installation | the App's bot | write the records branch, merge a verified change |

A comment or a review is copied into the records branch before any action reads
it (ADR-0016 decision 4). The separation of the rows above holds only while each
harness session gets the App's user token and not the Operator's own
credentials (O-77).

## 7. The records branch

The branch `layup-records` of each target (ADR-0016). Every table is
tab-separated with a fixed header row; the headers below are the start set, and
the implementation plan fixes their column order.

| File | Kind | Columns |
| ---- | ---- | ------- |
| `parameters.tsv` | register | name, value, default, source, set_at |
| `approvers.tsv` | register | decision point, account, role, set_at, source |
| `harnesses.tsv` | register | harness, command, version, models, reports_tokens, enforces_cap, probe_result, probed_at |
| `roles.tsv` | register | role, phase, answers_kinds, source |
| `routing.tsv` | register | role, task_class, author_harness, author_model, verifier_harness, verifier_model, weight, evidence, set_at_retro |
| `prices.tsv` | register | harness_or_provider, model, unit, price, source, effective_date |
| `budget.tsv` | register | milestone, B, U, task, estimate, source |
| `stack.tsv` | register | gate_id, kind, command, tool_version, fixture, evidence |
| `phase.tsv` | events | seq, time, milestone, phase, event, attempt, actor, head |
| `heartbeat.tsv` | events | time, run, phase, attempt |
| `acceptance.tsv` | events | requirement, milestone, delivered_commit, decision, account, comment_id, review_count, time |
| `forge/<id>.md` | payload | a copied comment or review: body, author, id, `performed_via_github_app`, time, SHA-256 |
| `runs/<task>/handoffs.tsv` | events | seq, time, from_role, to_role, kind, artifact, artifact_hash, condition, valid |
| `runs/<task>/questions.tsv` | events | qid, asked_by, kind, kind_source, owner, asked, answered, accepted, human, planned |
| `runs/<task>/decisions.tsv` | events | time, point, deterministic_result, provider, model_version, question_set, answer, probabilities, threshold, threshold_source, authority, branch, fallback_reason, input_tokens |
| `runs/<task>/telemetry.tsv` | events | time, task, requirements, cost_share, role, harness, model, tokens_in, tokens_out, latency_ms, wall_s, price, price_source |
| `runs/<task>/stall.md` | payload | the stall record: limit reached, attempts, examiner, diagnosis, outcome |
| `runs/<milestone>/reward.tsv`, `lessons.tsv` | events | the reward per routing row; the lessons with their evidence |

`not reported` in a numeric column is incomplete, never zero (`F-0003#60`).

## 8. The parameter register

Each row of `parameters.tsv` has a default that LAYUP ships and its source. A
parameter tunes how and when a PSB rule applies; it cannot switch one off
(O-84).

| Parameter | Default | Source |
| --------- | ------- | ------ |
| `smartif.provider`, `smartif.model` | `jev`, `jev-1.13.0` (pinned) | O-78; `runs/T-hbw8/jev-sources.md` |
| `smartif.<point>.authority` | `shadow` | O-78, ADR-0014 |
| `smartif.<point>.threshold` | a conservative start value per point, recorded with its evidence before the first delegated use | O-79; `F-0001#4` |
| `escalation.screen` | `task-and-handoff` | O-79, ADR-0018 |
| `stall.T` | 10 minutes (a maximum) | O-82 |
| `stall.N` | 1 | O-82 |
| `heartbeat.H` | set at the first run, with its reason | ADR-0013 |
| `harness.<id>.paid` | allowed with a wall-clock proxy, telemetry incomplete | O-80 |
| `roles.matrix` | the seven PSB §2 functions and the owner map of ADR-0015 | O-81 |
| `learn.min_evidence`, `learn.max_step`, the reward terms | a cited start value at the first retrospective; `learn.max_step` = 0 freezes the weights | O-83, ADR-0019 |
| `cost.task_factor` | set with the budget file, with its source | O-68, ADR-0018 |

## 9. Coverage

### 9.1 The In-Scope items of the PSB

| In-Scope item | Component |
| ------------- | --------- |
| `F-0003#41` Problem Statement Quality | `layup psb check` in Intake (delivered) |
| `F-0003#42` Reproducible Discipline Setup | `layup setup`, `layup setup verify`, the stack's gate recipe (ADR-0017) |
| `F-0003#43` Rule Protection | rule paths changed only in the retrospective batch; the orchestrator's merge rule; `layup audit` (ADR-0017, ADR-0016) |
| `F-0003#44` Stack-Dependent Gates | the target's native stack gates and `layup gate` (ADR-0017) |
| `F-0003#45` Role Handoffs | `handoffs.tsv` and its check (ADR-0016, ADR-0013) |
| `F-0003#46` Autonomous Clarification | question routing by ambiguity kind to the owner role (ADR-0014, ADR-0015) |
| `F-0003#47` Verification on Every Change | the target's gates and the counterpart verification before the merge (ADR-0017, ADR-0015) |
| `F-0003#48` Human-on-the-Loop | the order of deciding; the planned approval points listed at Intake; the escalation screen (ADR-0014, ADR-0018) |
| `F-0003#49` Stall Resolution | the clock, the limits `T` and `N`, the examiner, the stall package, the dead-man job (ADR-0013) |
| `F-0003#50` Cost Visibility | `telemetry.tsv` per action with the requirement and the price (ADR-0016) |
| `F-0003#51` Specification Synthesis | numbered facts, `layup spec draft`, `layup spec check` (ADR-0020) |
| `F-0003#52` Harness-Agent Neutrality | the records branch readable from a clean clone; one neutral rule source; the counterpart harness (ADR-0016, ADR-0015) |

### 9.2 The requirements of `PRD-0001`

| Requirement | Phase | Component |
| ----------- | ----- | --------- |
| REQ-001 | 1 | `layup psb check` (delivered) |
| REQ-002 | 1 | `layup setup`, `layup setup verify` |
| REQ-003 | 2 | the retrospective batch, the merge rule, `layup audit`; with the limit of section 10 |
| REQ-004 | 1 | `layup gate` runs the target's native gates from outside (ADR-0017) |
| REQ-005 | 2 | `handoffs.tsv` and its check |
| REQ-006 | 3 | question routing and `questions.tsv` |
| REQ-007 | 1 | the target's CI gates required before the merge; the orchestrator merges only a verified change |
| REQ-008 | 3 | the escalation screen and the human-input accounting (`questions.tsv` human and planned columns) |
| REQ-009 | 1 | `runs/<task>/stall.md` |
| REQ-010 | 3 | the stall procedure of ADR-0013 |
| REQ-011 | 1 | `telemetry.tsv` |
| REQ-012 | 2 | ADR-0020 |
| REQ-013 | 4 | ADR-0015, ADR-0016 |
| REQ-014 | 4 | two gate recipes and two pilot targets |
| REQ-015 to REQ-018 | — | not built: no weight change (ADR-0019); the intent stays with the idea owner (ADR-0018); no cloud change; a gate adds rules and weakens none (ADR-0017) |
| NFR-001 | 1 | the records branch and the forge copy |
| NFR-002 | 1 | native gates in the target (ADR-0017) |
| NFR-003 | 1 | every parameter and price with its source; `layup setup verify` |
| NFR-004 | 1 | `not-active` is not a pass, in `layup gate` and in the counterpart verification |
| NFR-005 | 1 | the engine checks make no model call; only the smart-if component calls one (ADR-0014) |
| NFR-006 | 1 | the pin ([`setup/armature.pin`](setup/armature.pin), ADR-0009) |
| NFR-007 | 1 | Go, standard library: the provider and the forge through `net/http` |

### 9.3 The findings of the review of #69

| Finding | Answer |
| ------- | ------ |
| A1 a rule-path approval counted as planned | rule paths change only in the retrospective batch, a planned point listed at Intake; any other change is unplanned input (ADR-0017) |
| A2 the target cannot pass its gates without LAYUP | the stack gates are the target's own (ADR-0017, O-76) |
| A3 every question goes to the Operator | question routing to the owner role (ADR-0014, ADR-0015) |
| B1 nothing starts the role agents | `layup run` starts and stops the role sessions (ADR-0013) |
| B2 the timed stall is not detected | the clock in `layup run` and the dead-man job (ADR-0013) |
| B3 no record of requirement acceptance | `acceptance.tsv`; telemetry per requirement with a price (ADR-0016) |
| B4 no specification synthesis | ADR-0020 |
| B5 no counterpart harness | admission in code: the verifier harness differs from the author's (ADR-0015) |
| B6 decisions stay on the forge | the forge copy before any action; the records branch (ADR-0016) |
| B7 no cost stop before the budget is spent | the cost stop in code between actions; the per-action cap where a harness enforces one (ADR-0018); the stated limit of section 10 |
| B8 an escalation after the work | the screen before each task and at each handoff (ADR-0018) |
| B9 the source of the stack gates | a gate recipe per stack in LAYUP, with fixtures and evidence (ADR-0017) |
| The forge question | GitHub only for the pilot: a known limit (section 10) |

### 9.4 The vision brief (table C of the review)

| Vision item | Answer |
| ----------- | ------ |
| 2.1 PDR, PRD and a phased plan | the Design, Architect and Plan phases; ADR-0020 |
| 2.2 squads, model allocation by complexity | the routing table and the fit score among admitted pairs (ADR-0015) |
| 2.2 cross-verification by a counterpart harness | ADR-0015 decisions 4 and 5 |
| 2.3 milestones, retrospectives, lessons | the Retrospective phase; `lessons.tsv` (ADR-0019) |
| 3.1 model, context and solution routing | model and harness fit (ADR-0014, ADR-0015); context: the prompt file names the records a session reads (ADR-0013); solution: the panel and the selection |
| 3.2 a blind panel with a hypothesis posture | a panel as a stall action and in Design (ADR-0006); members see no other member's output; the examiner gets the evidence, not the performers' reasoning |
| 3.3 reinforcement learning on routing | the reward and the bounded update, with no weight change (ADR-0019) |
| 3.4 communication through issues | issues carry the human decisions; the forge copy puts them into Git (ADR-0016) |
| 3.4 telemetry per action | `telemetry.tsv` per action (ADR-0016) |
| 3.5 circuit breaker, diagnostic package, external answers | the stall limits, the stall package, the Operator's answer copied into Git (ADR-0013, ADR-0016) |

## 10. Known limits

- **GitHub only.** The forge copy, the App and the badge are GitHub mechanisms;
  another forge is out of the pilot.
- **One account.** Agent and human are told apart by the App badge, which holds
  only while the agent sessions cannot reach the Operator's own credentials
  (O-77). GitHub cannot stop an actor with the Operator's rights from merging or
  from changing a rule path on a repository owned by one user; the merge rule
  and `layup audit` detect it, they do not prevent it.
- **Spend inside one action.** The cost stop acts between actions; only a
  harness or provider that enforces a cap per action stops the spend inside one.
- **The host.** `layup run` needs its host on during delivery; the dead-man job
  detects its loss, it does not prevent it.
- **The smart-if provider** is an external service whose price and rate limits
  can change without notice; the fallback keeps a run safe, not fast.
- **Meaning.** A byte-exact trace does not prove that a paraphrase keeps the
  meaning; the counterpart verifier and the idea owner's acceptance carry it.

## 11. For the implementation plan

The implementation plan (#42) slices this architecture into tasks. It fixes the
column order of each table of section 7; the harness adapters for the pilot
harnesses; the first gate recipe (Go) and the second stack; the App's
permissions and the ruleset probe; the start values of the thresholds, `H`,
`cost.task_factor` and the reward terms, each with its evidence; and the order
in which the phases are built, `layup setup` first.
