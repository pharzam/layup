# The Operator's research review by GPT-6 on Codex

The Operator gave this report to the author's Claude Code session on 2026-09-30, as a file (SHA-256 `865749237893281611aba81a64212acef06d42d5281770adffc1911674d26702`): a detailed state-of-the-art research and review by GPT-6 on the Codex desktop harness, head `0ee5c8b`. Its section 6 repeats the four findings of [the Operator's Codex review](operator-codex-review.md), which are answered there. It is copied here word for word, so a fresh session can read it from Git (O-73). Not a gate round.

~~~~text
# LAYUP architecture: detailed SOTA research and comparison

- Research date: 2026-09-30
- Reviewer: GPT-6, Codex desktop harness
- Reviewed branch: T-hbw8, issue #72
- Reviewed head: `0ee5c8b943a5e39cb82538af3773c47104a7e2f9`
- Comparison base: `e08f584538ddf5f2cc8df10efec2d62f48800669`
- Status: research and review record; recommendations are not adopted architecture decisions

## 1. Research conclusion

LAYUP's central direction is consistent with current practice for dependable agent systems: a deterministic coordinator, explicit workflow transitions, durable project records, bounded model decisions, replaceable execution harnesses, native repository gates, and human decisions at declared points.

The most consequential weaknesses are in the contracts around that core:

1. Recovery after an external action has an uncertain outcome.
2. Ownership transfer while an earlier controller can still perform an action.
3. Isolation between agent-controlled execution and privileged credentials.
4. Empirical qualification of semantic decisions before delegation.
5. Evidence that routing-weight updates improve outcomes rather than reward easier assignments.
6. Complete accounting of context, external operations, costs and recovery behavior.

The recommended direction is to keep the core and strengthen these contracts. Replacing it wholesale with an agent framework is not justified by the evidence reviewed. A durable workflow platform could reduce implementation burden, but adopting one would require an explicit decision about dependencies, infrastructure and Git's authority.

This assessment separates two classes of result:

- **Bootstrap material findings:** contradictions, omissions or acceptance failures under the user's specified review test.
- **Research recommendations:** changes that would improve operational dependability or provide evidence for optimization. A disclosed limitation does not automatically become a new bootstrap blocker because a more mature production design exists.

## 2. Scope, method and evidence limits

### 2.1 Repository evidence

The design comparison uses these documents at the reviewed SHA:

| Document | Role |
|---|---|
| [docs/architecture.md][A] | Architecture under review; all fifteen sections considered. |
| [docs/facts/problem-statement-brief.md][P] | Product scope, invariants, acceptance criteria and human authority. |
| [docs/facts/architectural-vision-brief.md][V] | Architectural intent, decision routing, learning, telemetry and circuit breaking. |
| [docs/prd/PRD-0001-layup.md][PRD] | Derived requirements, including the current Go/standard-library constraint. |
| [runs/T-hbw8/inputs-from-pr-69.md][INPUTS] | The explicitly permitted supplementary input record. |

Issue comments and author reasoning records were not used as evidence. Referenced walkthroughs and ADRs were not independently audited; a coverage-table pointer is evidence of a declared mapping, not proof that its implementation or walkthrough passes.

Section numbering in the initial review request differs from the frozen document. At this head, records are in §3, decisions/defaults in §10, coverage in §14, and known limits in §15.

### 2.2 Research method

The research compared the design with primary sources from workflow platforms, agent-runtime developers, protocol maintainers, model providers and research authors. It focused on mechanisms that address LAYUP's requirements rather than treating popularity or novelty as evidence of suitability.

For this report, SOTA means the strongest relevant approaches supported by the reviewed evidence. It does not mean that one framework is universally best.

Three evidence levels are kept separate:

- **Documented behavior:** a source states what its system or protocol does.
- **Reported experiment:** authors report measured results within a particular evaluation.
- **LAYUP inference:** this review derives a recommendation from those facts and the repository design.

Provider-reported results are not independent validation. Neither Jev nor Laya was benchmarked on LAYUP tasks. No runtime, security, recovery or performance tests were executed as part of this design research. Proposed tests below are acceptance criteria for future work, not reported passes.

### 2.3 Constraints on recommendations

The following requirements constrain any adoption or redesign:

- Git remains the project system of record.
- A project remains usable without LAYUP.
- Working agents cannot change the rules or gates that judge their work.
- Configuration values require evidence.
- Inactive checks cannot be represented as passing.
- Deterministic checks precede model judgment when the rule is mechanically checkable.
- Domain-specific content cannot weaken baseline rules.
- The Armature version is pinned and recorded.
- Harnesses remain replaceable.
- Human authority and escalation remain within the PSB's decision points.
- PRD NFR-007 currently specifies Go, the standard library and the Git program.

These constraints come from the [PSB invariants][P-138], its [pass/fail criteria][P-152], and [PRD NFR-007][PRD-105]. A recommendation that changes them requires an explicit project decision; a research report cannot authorize that change.

## 3. Comparison across all architecture sections

| Section | Existing design | Research assessment | Recommended disposition |
|---|---|---|---|
| §1 Boundaries and technology | LAYUP operates outside the target; Go core; forge capability interface; native target gates. | A small control layer with replaceable execution is a sound fit. Supporting only GitHub initially is a disclosed product limit. | Retain. Require capability probes to report unsupported behavior honestly. |
| §2 Runtime and coordination | One controller per target; Git lease/heartbeats; announce writes before forge effects. | Clear ownership intent, but uncertain external outcomes and stale-controller effects need fuller contracts. | Specify operation recovery, in-flight effects and fencing. |
| §3 Git records | Dedicated records branch, structured registers, copied human decisions and a single writer. | Compatible with Git as authority. Event storage alone does not define replay or recovery. | Add operation receipts, record versions and reconstruction rules. |
| §4 Harness execution | Separate clones, homes, prompts and result folders; environment filtering. | Good session hygiene, but the shared OS identity leaves privileged files accessible. | Enforce isolation between execution and control authority. |
| §5 Intake/setup | Structured gap questions, cited configuration, pinned baseline and setup verification. | Strong provenance. Compatibility failure is recognized but could be detected earlier. | Retain; move feasible compatibility checks before target mutations. |
| §6 Gates/rule protection | Trusted gate definitions, protected paths, known-bad fixtures, approved rule batches. | Strong design for rule integrity. Execution of untrusted product tests also needs containment. | Retain gate ownership; isolate execution of tested code. |
| §7 Specification synthesis | Byte-exact source spans, requirement identifiers and semantic counterpart review. | Traceability is strong; structural completeness cannot prove semantic correctness. | Retain both mechanical and semantic checks; evaluate omissions and misclassification. |
| §8 Lifecycle/handoffs | Explicit phases, schema-validated handoffs, bounded attempts, independent verification. | Good workflow structure. Acceptance authority and transition recovery need clarification. | Correct authority; define resumable states at each boundary. |
| §9 Roles/routing/context | Seven default roles, capability registry, tier routing and linked context bundles. | Sensible responsibilities. Optimal topology and context sizing remain empirical questions. | Keep roles configurable; evaluate routes and improve context accounting. |
| §10 Smart-if/escalation | Fixed decision points, deterministic floor, shadow/cautious/delegate modes and fallback. | Correct control order. Threshold admission does not yet establish decision quality. | Require evaluation evidence before delegated authority. |
| §11 Stalls/panels | Progress measures, retry ladder, blind panels, circuit breaker and dead-man notice. | Bounded escalation is appropriate. Detection/recovery guarantees depend on infrastructure. | Test failure states; state achievable detection and recovery bounds. |
| §12 Cost/telemetry | Session ledger, explicit unknowns, reservations, budget band and retrospective audit. | Honest accounting foundation. Some caps are estimates and per-action causality is incomplete. | Distinguish enforced and estimated limits; connect action records and receipts. |
| §13 Learning | Explicit reward formula, route-weight adjustment, bounded exploration and propose default. | Useful experiment, but raw route means are vulnerable to assignment bias. | Validate benefit before automatic adoption. |
| §14 Coverage | All twelve PSB In-Scope items have mappings or limits. | Broad coverage; the vision's solution-routing clause lacks an explicit answer. | Add the missing solution-routing contract or limitation. |
| §15 Known limits | Shared host, stale write, watchdog, telemetry and learning limitations are listed. | Disclosure is useful and materially better than an implicit guarantee. | Turn important limits into deployment restrictions and measurable exit criteria. |

## 4. Detailed research results

### 4.1 Deterministic orchestration remains the right foundation

**Existing design.** Architecture §10 fixes the order as deterministic checks, then semantic smart-if decisions, then humans at a decision point. §4 keeps ordinary engine checks free of model calls. §8 gives delivery an explicit state machine.

**External evidence.** Anthropic's workflow guidance distinguishes prescribed workflows from more autonomous agents and recommends starting with simple, composable patterns. That supports using models selectively where judgment is needed, rather than making every transition a model decision. [Building effective agents][S01].

**LAYUP assessment.** Preserve the deterministic coordinator. The PSB's invariants, approval authority, budgets and gate outcomes are a particularly poor fit for unconstrained conversational coordination.

Recommended refinements:

1. Give every transition explicit preconditions and a durable result.
2. Make invalid transitions fail mechanically.
3. Keep model outputs as typed inputs to the transition logic.
4. Record which rule, decision or human authority caused a transition.
5. Treat replacement of a harness as an adapter change, not a change to project policy.

The seven named roles are useful responsibilities. They should not imply that every task benefits from seven independent sessions or concurrent agents. Google's controlled agent-system research found that benefits depend on task structure and coordination costs. This is evidence for evaluating topology, not evidence for a particular LAYUP role count. [Google Research][S02].

**Proposed validation.** Run representative tasks through the required role structure and compare permitted routing/topology variants. Measure accepted outcomes, intervention, cost, latency and failure recovery. Never remove required independent verification solely to improve throughput.

### 4.2 Durable records need an external-operation recovery protocol

**Existing design.** Architecture §2 records the intention to write to the forge before performing the write. §3 makes Git authoritative and limits record mutation to the controller.

**Gap.** The document does not fully define the crash window after a remote effect succeeds but before its receipt is durably recorded. The same issue occurs when an HTTP timeout hides whether GitHub accepted the request.

Examples include:

- A comment is created, but the response is lost.
- A branch push succeeds, but the controller crashes before binding the SHA.
- A merge succeeds, but the local step still appears pending.
- A paid invocation begins, but its session-start acknowledgement is lost.

Retrying blindly can duplicate work. Assuming success can lose work or skip a check.

**External evidence.** Temporal explains why external activities need idempotency even when workflow execution is durable. A workflow engine's history does not make arbitrary third-party effects execute exactly once. [Temporal durability and idempotency][S03].

**Recommended contract.** For each externally visible operation, specify:

| Element | Purpose |
|---|---|
| Stable operation ID | Recognize the same logical operation across retries and restarts. |
| Target and expected version | Prevent an operation from silently applying to changed state. |
| Durable intent | Record authorization and requested effect before execution. |
| Completion receipt | Bind the result to the operation and relevant external identifier. |
| Unknown-outcome state | Preserve ambiguity instead of falsely reporting success or failure. |
| Reconciliation procedure | Inspect remote state before deciding whether retry is safe. |
| Operation-specific retry policy | Account for differences between comments, pushes, merges and paid calls. |

This is a proposed design contract, not an existing LAYUP schema. Exact field names and formats should be selected in the architecture change.

A fresh clone should be sufficient to reconstruct the next valid action. A runtime cache may accelerate that reconstruction, but cannot hold the only copy of project state or decisions.

**Proposed validation.** Inject a crash before the intent, after the intent, during the remote request, after remote success and before the receipt. Recover from Git and prove that the final result is correct and duplicate effects are controlled.

### 4.3 Lease ownership is not sufficient fencing

**Existing design.** The controller uses locally observed heartbeat changes to detect a stale owner. A records push based on an obsolete commit is refused. The architecture explicitly permits one previously announced forge write after takeover; see [§2][A-75] and known limit L-A3.

**External evidence.** Kubernetes' leader-election implementation explicitly says that leader election does not guarantee that only one client acts as leader: it does not itself provide fencing. [client-go leader-election documentation][S04].

**LAYUP inference.** Preventing a stale Git records push does not necessarily prevent a stale GitHub action. The two systems are separate effect destinations.

A satisfactory design must address both new and already in-flight actions:

- Establish a single trusted path for privileged forge effects.
- Associate operations with an ownership epoch or equivalent authority.
- Reject obsolete authority at that execution boundary.
- Reconcile or drain in-flight effects before allowing incompatible successor work.
- Use external conditional operations, such as expected-head checks, where available.

A read of the lease immediately before an API call is not a complete solution: authority can change between the check and the effect. Nor should the design claim that an already accepted remote operation can always be cancelled.

For a simpler initial deployment, automatic takeover can remain unavailable until the old executor is demonstrably fenced. That is a product behavior decision to document, not a silent workaround.

**Proposed validation.** Pause controller A immediately before an external effect, let controller B acquire ownership, then resume A. Check each operation type, including merges and status updates.

### 4.4 Isolation must cover agents and verification workloads

**Existing design.** Architecture §4 separates clones and session homes, filters environment variables and withholds known forge credentials. However, [the documented shared-user limit][A-270] allows a session to read the App private key or the Operator's own forge login. The latter can make an agent-generated comment appear to satisfy the human-decision check.

**External evidence.** Current production runtime designs separate durable orchestration, replaceable agent execution and sandbox resources, while keeping credentials outside the sandbox. Anthropic also explains why filesystem and network restrictions need to work together. [Managed Agents architecture][S05], [sandboxing design][S06].

**Recommended boundary.**

- Controller state, credential material and approval validation remain outside agent-accessible execution.
- Each workload receives only the files and capabilities it needs.
- A trusted mediator performs allowed privileged operations.
- Network access is restricted or mediated according to the task.
- Cleanup removes session capabilities as well as temporary directories.
- Isolation applies to product tests, build scripts and other code executed during verification.

The last point extends beyond the role harness: §6 runs product tests on a host scratch worktree. Trusted policy files still invoke untrusted product code. GitHub's privileged-workflow guidance illustrates the same distinction between trusted control configuration and untrusted code execution. [GitHub security guidance][S07].

A container label alone does not establish the boundary; mounts, host sockets, identity, capabilities, network access and secrets determine what it protects. gVisor is one implementation option for a Linux executor, not a mandatory choice or a portable answer for every host. [gVisor architecture][S08].

**Proposed validation.** Use harmless adversarial fixtures that try to read a canary secret, write controller state, invoke an unauthorized forge action and submit a forged human approval. All must be prevented or rejected with recorded evidence.

### 4.5 Smart-if authority requires empirical admission

**Existing design.** §10 limits provider calls to named points, checks response shape and model version, pins the selected version, offers shadow mode and resets delegated points after a model change. These are useful controls.

The calibration clause at [architecture line 961][A-961] allows either an Operator comment or agreement with observed ground truth to serve as threshold evidence. An Operator comment establishes a configuration decision; it does not measure predictive quality.

**External evidence.**

- TypeSafe recommends selecting thresholds according to risk and evaluating them against the application's own data. Its confidence statistic is derived from the output distribution, and it explicitly allows alternative statistics. Consequently, LAYUP's use of the highest option probability is not, by itself, a documented API misuse. [TypeSafe confidence][S09].
- Laya's own model card reports 0.362 accuracy for its base English checkpoint and 0.766 for a specialized checkpoint on its typed-decisions benchmark. It also warns about confident errors and deployment-specific limits. These are author-reported results on that benchmark, not LAYUP performance measurements. [Laya model card][S10].

**Recommended admission record.**

1. Decision point and question-pack version.
2. Provider, checkpoint and runtime configuration.
3. Version of the context-construction logic.
4. Held-out evaluation population and labeling method.
5. False-negative and false-positive rates relevant to that point.
6. Error rate among accepted decisions and abstention coverage.
7. Threshold choice and its evidence.
8. Conditions that return the point to shadow mode.

For P1 and P2, missed business-forking decisions deserve direct measurement; aggregate accuracy can conceal a rare but consequential miss. For P3 and P4, observed outcomes are affected by the action chosen, so the evaluation design must address incomplete or biased labels.

Do not infer correctness from a high returned probability. Test missing context, ambiguous choices, contradictory evidence, multilingual inputs and oversized inputs.

These recommendations do not require changing foundation-model weights. Provider/checkpoint selection, output calibration and evaluation can remain within the current scope.

**Proposed validation.** Run each point in shadow against held-out and pilot cases, report errors and abstentions separately, and require a documented admission decision before delegation. Preserve deterministic floors throughout.

### 4.6 Learning needs a defensible comparison between routes

**Existing design.** §13 scores implementing routes using acceptance, verification, findings, stalls, interventions, reversals and cost. It compares route means within role and tier, adjusts weights within bounds, and defaults to proposing changes for adoption.

**Gap.** Grouping by role and tier does not necessarily equalize task difficulty. A route assigned difficult tasks can look worse than a weaker route assigned easy tasks. Delayed acceptance and reversal signals further complicate the comparison.

**External evidence.** Contextual-bandit methods explicitly represent context, action and the reward observed for the selected action. Off-policy evaluation requires appropriate logged data and assumptions; selection probabilities matter when using propensity-based estimators. [Contextual bandits][S11], [off-policy evaluation][S12].

**Recommended approach.**

- Keep the initial update mechanism explicitly experimental.
- Record relevant task characteristics and the complete eligible route set.
- Record the selection policy and, for randomized selection, the selection probability.
- Compare policies on matched held-out tasks or a bounded randomized experiment.
- Keep delayed outcomes linked to the original assignment.
- Separate quality constraints from optimization of cost and latency.
- Specify rollout and rollback criteria before automatic adoption.

Deterministic round-robin exploration does not automatically supply the statistical support needed for every off-policy estimator. Logging a nominal probability cannot repair a missing comparison population.

The present `propose` default is appropriate. The evidence does not establish that the reward formula or update rule is wrong in every setting; it establishes that improved performance has not yet been demonstrated.

**Proposed validation.** Compare candidate routing with a fixed baseline on the same task distribution. Show that measured savings do not come from reduced quality, more missed escalations, incomplete telemetry or easier assignments.

### 4.7 Context routing needs complete accounting and reproducible provenance

**Existing design.** §9 constructs context from linked requirements, plans, decision records, handoffs, answers and rules. It estimates size as bytes divided by four and refuses oversized starts.

**External evidence.** Current context-engineering guidance emphasizes selective retrieval, structured notes and managing finite context rather than indiscriminately supplying more text. [Context engineering][S13].

**Recommended refinement.**

- Record a context manifest identifying source paths, revisions, selected spans and generated inputs.
- Account for system instructions, tool definitions, message framing, retrieved content and reserved output.
- Use a provider tokenizer or counting API when available; otherwise use a documented conservative estimator whose error is evaluated.
- Support bounded retrieval from approved sources when static bundles omit necessary information.
- Preserve verified artifacts across recovery while maintaining fresh, independent review sessions.
- Treat retrieved text as task data, not as authority to change orchestration rules.

The current refusal behavior is safer than silent truncation. However, `bytes / 4` is not a dependable universal token bound, especially across languages and structured source material.

**Proposed validation.** Exercise multilingual documents, large tool schemas, long source files, retrieval expansion and output reservations. Verify both the total budget and the ability to reconstruct the exact context supplied.

### 4.8 Harness neutrality benefits from protocols, but protocols do not enforce policy

**Existing design.** The harness register records capabilities, versions, instruction behavior, usage reporting and spending controls. Session results are imported through controlled handoffs.

**External evidence.** Agent Client Protocol defines initialization, capability negotiation, sessions, updates, cancellation and permission interactions. [ACP overview][S14]. OpenHands' SDK research presents another approach to model-independent execution, structured events and sandboxed agent lifecycles. [OpenHands SDK paper][S15].

**LAYUP inference.** ACP is a useful optional adapter where the chosen harness supports it. A CLI adapter remains appropriate for other harnesses.

Keep the LAYUP contract above any transport:

- Validate actual capabilities rather than accepting labels.
- Record harness and model versions separately.
- Normalize cancellation, usage and completion semantics.
- Preserve the same rules and evidence regardless of transport.
- Do not treat a permission-request protocol as proof of OS isolation.

Different harnesses satisfy the PSB's stated independence criterion when the required conditions hold. They do not necessarily eliminate correlated errors if they share the same underlying model or inputs. Record that distinction and evaluate observed verification effectiveness.

### 4.9 Specifications and native gates are strengths to preserve

**Existing design.** §§5-7 bind configuration to cited facts, use exact source spans, generate identifiable requirements and keep target gates executable without LAYUP.

**External evidence.** Spec-driven toolchains such as GitHub Spec Kit similarly organize work around specifications, plans and implementation, but their existence does not prove semantic completeness. [Spec Kit][S16]. Agent-evaluation guidance emphasizes checking outcomes rather than trusting an agent's account of what it did. [Agent evaluations][S17].

**Recommended refinement.**

- Preserve mechanical source coverage and semantic counterpart review.
- Test adversarial briefs: contradictions, implicit constraints, duplicated needs and negative requirements.
- Check baseline layout compatibility before setup starts mutating a target.
- Keep required gate definitions separate from changes under test.
- Preserve known-bad fixtures and verify that they fail for the intended reason.
- Continue validating target independence with LAYUP removed.

Do not infer that every source byte being assigned to a category means the resulting requirements express the correct intent.

**Proposed validation.** Use a versioned set of briefs with independently established expected requirements, ambiguities and escalation points. Measure omissions and incorrect transformations as well as trace coverage.

### 4.10 Stalls need reliable observation and tested recovery

**Existing design.** §11 distinguishes progress from activity, limits attempts, requires diagnosis and outcomes, escalates through a blind panel and provides a dead-man workflow.

**External evidence.** GitHub documents that scheduled workflows may be delayed or dropped under load. Therefore the architecture's disclosed watchdog limitation is real. [GitHub scheduling][S18].

**Recommended refinement.**

- Distinguish workload inactivity, controller failure, forge outage and planned human waiting.
- Record when detection became possible and when recovery actually began.
- Use an independent supervisor if a required detection bound cannot be supported by GitHub scheduling.
- Ensure the diagnostic package can be reconstructed after controller failure.
- Test that a blind panel cannot silently enlarge scope or reset budgets.

A foreground controller is not inherently wrong for a bootstrap product. Its availability and recovery expectations must match the deployment claim.

**Proposed validation.** Exercise a dead controller, delayed CI, dropped watcher execution, provider outage, frozen harness and planned human wait. Verify the correct classification, clock behavior and eventual outcome record.

### 4.11 Telemetry and budgets need causal completeness

**Existing design.** §12 records sessions, tokens, latency, duration and cost, distinguishes unknown values from zero, reserves costs before work, and audits outcomes.

**External evidence.** OpenTelemetry's GenAI conventions define structured agent and model-operation spans. The agent conventions remain marked Development, so their schema should not be treated as an immutable standard. [OpenTelemetry GenAI agent spans][S19].

**Recommended refinement.**

Use stable identifiers to connect:

- Project, run, requirement and task.
- Attempt, role session and harness/model version.
- Smart-if decision and context revision.
- Tool or command execution.
- Forge intent, external result and recovery attempt.
- Budget reservation and eventual usage reconciliation.

Git remains authoritative. An exporter can provide operational visibility without becoming the only store of project decisions or state. Pin an exporter mapping so external schema evolution does not rewrite historical meaning.

For spending, distinguish:

| Limit type | What the design may claim |
|---|---|
| Enforced provider/executor cap | Work is stopped according to the documented enforcement semantics. |
| Conservative upper bound | Cost stays within a justified bound under stated assumptions. |
| Estimate or time proxy | Expected exposure is estimated; a monetary ceiling is not guaranteed. |
| Unknown usage | Completeness is not established; unknown is not zero. |

A wall-clock limit cannot establish a monetary ceiling without a valid bound relating elapsed time to cost. Account for concurrent work, late usage reports and calls already accepted when cancellation occurs.

**Proposed validation.** Simulate missing usage, late receipts, concurrent reservations, provider timeout and cap overshoot. Verify honest accounting and refusal/escalation behavior.

## 5. Candidate systems and adoption trade-offs

This table is an architecture-fit comparison, not a completed procurement, license, vulnerability or operational-support review.

| Candidate | Relevant capability | Constraint or trade-off | Recommendation |
|---|---|---|---|
| Current Go core plus Git | Direct fit to repository independence, explicit policy and the current dependency constraint. | LAYUP must specify and implement recovery, versioning and effect coordination itself. | Retain as the baseline; make those contracts explicit. |
| Temporal | Durable workflow execution, retries, timers and human waits; Go support. | Adds a service, SDK and workflow history; adoption must reconcile Git authority and NFR-007. | Evaluate if recovery implementation becomes a dominant cost. [Temporal][S20] |
| Restate | Durable services, keyed exclusive access and workflows; Go SDK. | Introduces runtime infrastructure and another state/history system. External effects still need appropriate semantics. | Compare against Temporal for a narrowly scoped prototype if adoption is authorized. [Restate][S21] |
| LangGraph | Persistent checkpoints and resumable agent workflows. | Python/JavaScript ecosystem and external persistence are less direct fits for the current Go core. | Use as a reference for state semantics, not a default replacement. [LangGraph][S22] |
| OpenHands SDK | Structured agent execution, lifecycle events and sandbox integrations. | An execution runtime does not own LAYUP's governance, acceptance or repository invariants. | Consider as a replaceable execution adapter only. [OpenHands][S15] |
| Agent Client Protocol | Standardized lifecycle and capability negotiation between clients and agents. | Support varies; protocol conformance does not establish isolation or policy compliance. | Add where it reduces adapter complexity without weakening capability checks. [ACP][S14] |
| Jev / TypeSafe | Typed, option-driven decisions through a provider API. | Quality and appropriate thresholds are task-specific; a provider claim is not an admission test. | Keep behind the provider interface; qualify each decision point. [TypeSafe][S09] |
| Laya | Open-weight typed decision models and a compatible serving interface. | Checkpoint, language, context and deployment details matter; self-hosting still consumes resources. | Keep optional; pin and evaluate the exact deployment. [Laya][S10] |
| OpenTelemetry exporter | Correlated operational traces and integration with observability tools. | GenAI conventions are evolving; external backends must not become the sole project record. | Optional projection from stable internal records. [OpenTelemetry][S19] |

Before adopting a new platform, assess license terms, release health, security posture, dependency footprint, hosting requirements, failure behavior and migration/reconstruction cost. The present research establishes relevant capabilities and fit questions, not a final vendor selection.

## 6. Bootstrap review record and material findings

- Reviewer: GPT-6, Codex desktop harness
- Head: `0ee5c8b943a5e39cb82538af3773c47104a7e2f9`
- Lens: correctness and acceptance criteria
- Cycle: round 1 of 1 for a document change in bootstrap mode

### F1. Clarification Turnaround measures the wrong endpoint

**Observation:** Architecture §12 calculates the 95th percentile of `answered - asked`, restricted to answers later accepted.

**Contradicted clause:** The PSB's Clarification Turnaround measure is time from the question to the accepted answer.

**Effect:** Time spent reaching acceptance is omitted, so the reported metric can appear to satisfy its target even when the required metric does not.

**Required correction:** Calculate `accepted - asked` for the eligible autonomous questions.

**Basis:** [architecture line 1280][A-1280]; [PSB line 179][P-179].

### F2. Requirement acceptance authority is ambiguous

**Observation:** Intake permits naming an approver for each planned approval point. The Accept phase then says the approver accepts or rejects the requirement.

**Contradicted clause:** The PSB defines the idea owner as the person who accepts delivered requirements, and repeats that authority in the definition of Requirement.

**Effect:** The architecture does not clearly prevent another configured approver from exercising authority reserved by the PSB.

**Required correction:** Explicitly bind requirement acceptance to the idea owner. Any delegation beyond the PSB would need an authorized change to that requirement.

**Basis:** [architecture Intake][A-343], [architecture Accept][A-721], [PSB idea owner and Requirement definitions][P-199].

### F3. The vision's solution-routing clause has no explicit coverage answer

**Observation:** The §14 row for vision §3.1 covers model and context routing.

**Contradicted clause:** Vision §3.1 also requires systematic evaluation, scoring and ranking of competing technical pathways against architectural constraints.

**Effect:** Selecting a model or harness does not specify how a technical solution is selected.

**Required correction:** Identify the component or phase, candidate inputs, deterministic exclusions, comparison criteria, recorded result and authority boundary; alternatively name the limitation explicitly in the coverage response.

**Basis:** [architecture coverage row][A-1440]; [vision §3.1][V-39].

A suitable proposed flow is: enumerate feasible alternatives; eliminate violations deterministically; gather evidence for the remaining trade-offs; apply semantic judgment where required; record the choice and rejected alternatives; escalate only if it changes a PSB-defined business boundary. Numeric weights must have evidence rather than being invented for the report.

### F4. The reviewed diff exceeds the supplied budget

**Observation:** The comparison from the recorded base to the reviewed head contains 90 changed files, 12,036 added lines and 43 removed lines: 12,079 lines added plus removed.

**Contradicted constraint:** The review request approved at most 3,200 lines added plus removed over 36 files.

**Effect:** The required budget assertion cannot honestly be reported as passing.

**Required correction:** Reduce or split the change, or obtain an explicit revised budget. Excluding `runs/` is not part of the supplied rule and does not establish compliance; the previously measured remainder was still 3,249 lines across 38 files.

**Basis:** Git diff statistics at the stated base and head.

### Notes that do not independently cause another round

- All twelve PSB In-Scope items, F-0003#41 through #52, have coverage rows naming a component, walkthrough, check or known limit.
- The decision order is explicitly deterministic checks first, smart-if for semantic questions, and humans at defined points or fallback.
- The document says parameter settings cannot disable PSB rules, but a declared invariant is not proof that every execution path satisfies it.
- Known limits are explicitly collected in §15. Security, stale-write, token-visibility and watchdog limitations must not be represented as passed acceptance criteria.
- Git records can satisfy the system-of-record intent when host configuration actually used and human decisions are copied as specified. Recovery and reconstruction semantics still need refinement.
- Coverage mappings and disclosed limits were reviewed as design evidence; implementation acceptance tests were not run.
- The SOTA recommendations in §4 do not independently require another bootstrap review round.

**Verdict:** not mergeable, findings recorded.

Before approval, resolve F1-F3 and the budget failure F4. The research recommendations should be prioritized separately according to the intended deployment and their acceptance implications.

## 7. Proposed sequence for architecture improvements

The following is a proposed sequence for future work. Saving this report does not adopt it or change operative rules.

| Phase | Scope | Reviewable output | Exit evidence |
|---|---|---|---|
| A. Restore consistency | Metric, acceptance authority, solution-routing coverage and budget. | Corrected architecture clauses and an in-budget change or approved budget. | Clause-level review against the cited requirements. |
| B. Define recovery | Operation identities, intents, receipts, unknown outcomes, reconstruction and ownership transfer. | State/operation contracts and crash-point scenarios. | No lost or uncontrolled duplicate effects in fault-injection cases. |
| C. Enforce isolation | Credential boundary, agent execution and product-test execution. | Threat model, capability contract and chosen isolation mechanism. | Adversarial fixtures cannot access control authority. |
| D. Qualify decisions | Per-point evaluation sets, threshold evidence and promotion/reset rules. | Versioned evaluation and admission record. | Reported error, abstention and escalation behavior on held-out cases. |
| E. Establish operational truth | Context manifests, causal telemetry, budget classifications and watchdog expectations. | Observable action chain and explicit operational limits. | Recovery, accounting and context-limit scenarios pass. |
| F. Validate optimization | Routing experiments and learning adoption criteria. | Comparison with fixed routing and rollback plan. | Improvement under the same quality and authority constraints. |

Keep the implementation scope of each phase small enough to produce independent evidence. A learning or routing improvement should not obscure unresolved recovery or authorization failures.

## 8. Proposed validation matrix

These are proposed acceptance scenarios, not tests executed by this research.

| ID | Scenario | Expected evidence |
|---|---|---|
| REC-01 | Crash after durable intent, before external request. | Recovery executes or safely cancels the pending operation once its state is established. |
| REC-02 | External action succeeds; response or receipt is lost. | Reconciliation identifies the actual effect before retry. |
| REC-03 | Recover on a fresh host using project records and authorized credentials. | Next valid action is reconstructed without undocumented local state. |
| REC-04 | Restart after a record-schema or controller-version change. | Compatible replay, explicit migration or an honest refusal. |
| OWN-01 | Old owner resumes after successor takeover. | Obsolete authority cannot issue an incompatible new effect; in-flight effects are accounted for. |
| ISO-01 | Agent searches outside its task files. | Canary secrets and controller state remain inaccessible. |
| ISO-02 | Product test attempts the same access. | Verification execution has equivalent protection. |
| AUTH-01 | Agent tries to impersonate a human acceptance. | Approval validation rejects it; requirement authority remains with the idea owner. |
| GATE-01 | Known-bad patch modifies or evades gate behavior. | The intended gate fails and cannot be bypassed by the working agent. |
| SPEC-01 | Brief contains contradictions and implicit constraints. | Mechanical trace checks and semantic review surface the expected issues. |
| DEC-01 | Provider confidently misses a business-forking decision. | Evaluation records the miss; admission criteria and deterministic floors remain explicit. |
| DEC-02 | Provider fails, abstains or receives oversized input. | The documented fallback occurs and is recorded. |
| DEC-03 | Provider or relevant question/context logic changes. | Re-evaluation or a return to shadow follows the version policy. |
| CTX-01 | Multilingual and tool-heavy session approaches its context limit. | Full input/output budget is considered; no silent truncation or false fit. |
| COST-01 | Usage arrives late or is missing during concurrent work. | Reservations reconcile honestly; missing values are not treated as zero. |
| STALL-01 | Controller, harness or watcher stops independently. | Correct classification, diagnosis and eventual outcome record. |
| LEARN-01 | Routes receive tasks of different difficulty. | Evaluation does not mistake assignment differences for route quality. |
| IND-01 | LAYUP is removed and another harness continues the project. | Native gates and required project context remain usable. |

## 9. Source register

All sources below were consulted during the research associated with this report on 2026-09-30. Publication dates and live documentation state can differ; repository comparisons are pinned to the reviewed SHA.

| ID | Primary source | Evidence used | Important limit |
|---|---|---|---|
| S01 | [Anthropic: Building effective agents][S01] | Simple workflows, selective autonomy and composable patterns. | Practitioner guidance, not a benchmark of LAYUP. |
| S02 | [Google Research: science of scaling agent systems][S02] | Controlled evidence that topology benefits depend on task structure. | Results do not identify LAYUP's optimal topology. |
| S03 | [Temporal: idempotency and durable execution][S03] | External-effect retries and idempotency requirements. | Platform guidance does not prove LAYUP's recovery behavior. |
| S04 | [Kubernetes client-go leader election][S04] | Leader election does not itself provide fencing. | LAYUP uses a similar pattern, not this library. |
| S05 | [Anthropic: Managed Agents architecture][S05] | Separation of session, execution and sandbox/credential concerns. | Vendor production architecture, not a required implementation. |
| S06 | [Anthropic: Claude Code sandboxing][S06] | Filesystem and network containment. | Actual protection depends on configuration and host platform. |
| S07 | [GitHub: secure use of pull_request_target][S07] | Privileged workflow authority and untrusted code separation. | The source addresses Actions; application to LAYUP's host is an architectural inference. |
| S08 | [gVisor architecture guide][S08] | An example of an isolation boundary for untrusted workloads. | Not a portability or deployment recommendation by itself. |
| S09 | [TypeSafe confidence documentation][S09] | Distribution-derived confidence and risk-specific thresholds. | No demonstrated accuracy on LAYUP decisions. |
| S10 | [Laya model card][S10] | Checkpoint-specific results, context/language limits and caveats. | Author-reported measurements; cross-provider comparisons are not controlled LAYUP tests. |
| S11 | [Vowpal Wabbit contextual-bandit tutorial][S11] | Context, action, observed reward and action probability. | A learning method reference, not a required dependency. |
| S12 | [Vowpal Wabbit off-policy evaluation][S12] | Requirements for evaluating policies from logged behavior. | Estimator assumptions and support must actually hold. |
| S13 | [Anthropic: effective context engineering][S13] | Selective context, retrieval and structured persistent state. | LAYUP-specific token and quality measurements are still needed. |
| S14 | [Agent Client Protocol overview][S14] | Sessions, negotiation, cancellation and permission interactions. | Protocol support does not imply security or governance compliance. |
| S15 | [OpenHands Software Agent SDK paper][S15] | Modular execution, structured events and sandboxed lifecycles. | Framework results are not evidence for replacing LAYUP's coordinator. |
| S16 | [GitHub Spec Kit][S16] | Specification-to-plan-to-implementation workflow. | Traceability structure is not a semantic correctness proof. |
| S17 | [Anthropic: demystifying agent evals][S17] | Outcome-oriented evaluation and repeated trials. | LAYUP's task population and graders must be designed separately. |
| S18 | [GitHub workflow schedule documentation][S18] | Scheduled jobs may be delayed or dropped. | Does not supply a reliable maximum detection interval. |
| S19 | [OpenTelemetry GenAI agent conventions][S19] | Structured agent spans and explicit Development status. | Pin mappings; do not assume a stable schema. |
| S20 | [Temporal AI documentation][S20] | Durable agent workflows and operational capabilities. | Adoption adds infrastructure and dependencies. |
| S21 | [Restate service foundations][S21] | Durable services, virtual objects and workflows. | External-effect semantics and record authority still need design. |
| S22 | [LangGraph persistence documentation][S22] | Checkpoints and persistent workflow state. | Storage and language choices must fit the project constraints. |

## 10. Status of this artifact

This report preserves the research and review results for the frozen head. It does not alter the architecture, adopt an ADR, change the PSB, authorize a dependency or claim that proposed validation has passed.

The report was saved as a separate research artifact. The branch budget in §6 is the measured budget of the reviewed head, not a claim that this additional report belongs to or has been approved within that frozen diff.

[A]: https://github.com/pharzam/layup/blob/0ee5c8b943a5e39cb82538af3773c47104a7e2f9/docs/architecture.md
[P]: https://github.com/pharzam/layup/blob/0ee5c8b943a5e39cb82538af3773c47104a7e2f9/docs/facts/problem-statement-brief.md
[V]: https://github.com/pharzam/layup/blob/0ee5c8b943a5e39cb82538af3773c47104a7e2f9/docs/facts/architectural-vision-brief.md
[PRD]: https://github.com/pharzam/layup/blob/0ee5c8b943a5e39cb82538af3773c47104a7e2f9/docs/prd/PRD-0001-layup.md
[INPUTS]: https://github.com/pharzam/layup/blob/0ee5c8b943a5e39cb82538af3773c47104a7e2f9/runs/T-hbw8/inputs-from-pr-69.md
[A-75]: https://github.com/pharzam/layup/blob/0ee5c8b943a5e39cb82538af3773c47104a7e2f9/docs/architecture.md#L75
[A-270]: https://github.com/pharzam/layup/blob/0ee5c8b943a5e39cb82538af3773c47104a7e2f9/docs/architecture.md#L270
[A-343]: https://github.com/pharzam/layup/blob/0ee5c8b943a5e39cb82538af3773c47104a7e2f9/docs/architecture.md#L343
[A-721]: https://github.com/pharzam/layup/blob/0ee5c8b943a5e39cb82538af3773c47104a7e2f9/docs/architecture.md#L721
[A-961]: https://github.com/pharzam/layup/blob/0ee5c8b943a5e39cb82538af3773c47104a7e2f9/docs/architecture.md#L961
[A-1280]: https://github.com/pharzam/layup/blob/0ee5c8b943a5e39cb82538af3773c47104a7e2f9/docs/architecture.md#L1280
[A-1440]: https://github.com/pharzam/layup/blob/0ee5c8b943a5e39cb82538af3773c47104a7e2f9/docs/architecture.md#L1440
[P-138]: https://github.com/pharzam/layup/blob/0ee5c8b943a5e39cb82538af3773c47104a7e2f9/docs/facts/problem-statement-brief.md#L138
[P-152]: https://github.com/pharzam/layup/blob/0ee5c8b943a5e39cb82538af3773c47104a7e2f9/docs/facts/problem-statement-brief.md#L152
[P-179]: https://github.com/pharzam/layup/blob/0ee5c8b943a5e39cb82538af3773c47104a7e2f9/docs/facts/problem-statement-brief.md#L179
[P-199]: https://github.com/pharzam/layup/blob/0ee5c8b943a5e39cb82538af3773c47104a7e2f9/docs/facts/problem-statement-brief.md#L199
[V-39]: https://github.com/pharzam/layup/blob/0ee5c8b943a5e39cb82538af3773c47104a7e2f9/docs/facts/architectural-vision-brief.md#L39
[PRD-105]: https://github.com/pharzam/layup/blob/0ee5c8b943a5e39cb82538af3773c47104a7e2f9/docs/prd/PRD-0001-layup.md#L105
[S01]: https://www.anthropic.com/engineering/building-effective-agents
[S02]: https://research.google/blog/towards-a-science-of-scaling-agent-systems-when-and-why-agent-systems-work/
[S03]: https://temporal.io/blog/idempotency-and-durable-execution
[S04]: https://pkg.go.dev/k8s.io/client-go/tools/leaderelection
[S05]: https://www.anthropic.com/engineering/managed-agents
[S06]: https://www.anthropic.com/engineering/claude-code-sandboxing
[S07]: https://docs.github.com/en/actions/reference/security/securely-using-pull_request_target
[S08]: https://gvisor.dev/docs/architecture_guide/intro/
[S09]: https://docs.typesafe.ai/confidence
[S10]: https://huggingface.co/convaiinnovations/laya
[S11]: https://vowpalwabbit.org/docs/vowpal_wabbit/python/latest/tutorials/python_Contextual_bandits_and_Vowpal_Wabbit.html
[S12]: https://vowpalwabbit.org/docs/vowpal_wabbit/python/latest/tutorials/off_policy_evaluation.html
[S13]: https://www.anthropic.com/engineering/effective-context-engineering-for-ai-agents
[S14]: https://agentclientprotocol.com/protocol/v1/overview
[S15]: https://arxiv.org/abs/2511.03690
[S16]: https://github.github.com/spec-kit/
[S17]: https://www.anthropic.com/engineering/demystifying-evals-for-ai-agents
[S18]: https://docs.github.com/en/actions/reference/workflows-and-actions/events-that-trigger-workflows#schedule
[S19]: https://github.com/open-telemetry/semantic-conventions-genai/blob/main/docs/gen-ai/gen-ai-agent-spans.md
[S20]: https://docs.temporal.io/ai
[S21]: https://docs.restate.dev/foundations/services
[S22]: https://docs.langchain.com/oss/python/langgraph/persistence
~~~~

## The author's answer

Findings F1 to F4 (report section 6) are the findings of the Operator's Codex review; they are answered in [operator-codex-review.md](operator-codex-review.md) (F1 to F3 fixed in `f20d5a3`, F4 answered with O-85, O-89 and O-99). At `94748c5` the branch is 7,562 lines over 64 files against `946edfc`, inside the O-99 maximum of 8,500 over 72.

The research recommendations (report sections 4, 7 and 8) are not adopted as design; each is mapped to a known limit of `docs/architecture.md` §15 (O-110). Six were not stated there and are added as limits:

| Report | Limit | Status |
| ------ | ----- | ------ |
| 4.1 deterministic core, role topology | none | Already the design (§9, §10); no limit. |
| 4.2 recovery of external writes | L-A3 | Added: a forge write with an unknown result has no receipt and no check of the forge before a retry. REC-01 to REC-04. |
| 4.3 fencing | L-A3, §2 Fencing | Already stated. OWN-01. |
| 4.4 isolation of sessions | L-A1 | Already stated. ISO-01, AUTH-01. |
| 4.4 isolation of product tests | L-A1 | Added: `layup gate` runs the product's tests and build scripts on the host under the same user. ISO-02. |
| 4.5 admission of a decision point | L-E2, L-E1 | Added: the threshold can rest only on the Operator's comment at every point, and P1 and P2 have no held-out measure of missed business-forking decisions. DEC-01 to DEC-03. |
| 4.6 fair comparison of routes | L-H1, L-H2 | Added to L-H1: exploration goes in turn, not by chance, and the reward ignores task difficulty. LEARN-01. |
| 4.7 context accounting | L-D2 | Added: bytes divided by four is not measured, and it counts only the prompt file. CTX-01. |
| 4.8 harness protocols, correlated errors | new L-D4 | Added: independence by harness does not stop two harnesses on one model from making the same error. |
| 4.9 specifications and gates | L-C1, L-B3, L-A6 | Already stated. SPEC-01, GATE-01, IND-01. |
| 4.10 stalls | L-F1, L-A3 | Already stated. STALL-01. |
| 4.11 telemetry and budgets | L-G1, L-G2, L-G3 | Already stated. COST-01. |

These are limits, not new tasks (the Operator's rule: a finding that is not fixed becomes a known limit, not an issue). Adding them changes no behaviour, so it is applied as notes after round 6 without another round (Bootstrap mode rule 3), as the round-6 notes were.

Round 7 (O-111) moved the new limit from L-B4 to L-D4 (its note N3), because the B limits belong to §6; the comment on #72 of 2026-09-30 still says L-B4.
