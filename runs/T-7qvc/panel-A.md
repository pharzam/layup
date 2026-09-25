## Panel member A — CLI and Go architecture

These are options, not a vote or a selection. The proposed limits and role names below are **not** values set by the PSB. ADR-0011 fixes one Go CLI over repository files: its stack gates run outside the target, and no LAYUP code goes into the target. Today, the CLI has `version` and `psb check`; the five mechanisms below are not yet implemented.

### 1. Runner for stack gates — REQ-004, REQ-007

Both options need the same verdict rule. For every gate selected by the target’s recorded stack, publish `pass` as success and both `fail` and `not-active` as failure. Make the check required for the exact PR head, pin its reporting app, and require an up-to-date branch. A missing, cancelled, or skipped check must block merge. Keep gate-result events in the target’s Git; a forge check alone is not the system of record (`F-0001#1`, `#5`). If adding that evidence changes the PR head, run the gate again on the new head.

- **Option G1 — Event-driven external runner.** A target PR event starts a trusted job outside the target. The job checks out the exact proposed commit, runs the external `layup gate` binary without write credentials in that checkout, and a separate publisher posts each verdict to the target through a GitHub App.
  - **For, falsifiable:** It can report before review assignment. Open a PR with one seeded failure in each gate kind; if any reaches the managed review queue or can merge, this claim fails.
  - **Against, falsifiable:** The event relay and App add failure points. Drop and replay events, then check whether every current head gets one current result without a stale green result.
  - **Serves:** `F-0001#1`, `#2`, `#5`–`#7`, `#9`. **Strains:** `#4` if the App, gate set, or target registration has no recorded source. **Main risk:** A lost event leaves work blocked; a stale result tied to the wrong SHA can give a false green.

- **Option G2 — Polling external runner.** A scheduled job in LAYUP’s runner lists open PRs of registered targets and runs the same CLI against each new head. The App posts the results; the required check stays pending until the runner reports.
  - **For, falsifiable:** No target workflow or event relay is needed. Remove the relay and open a PR; if the poller does not find and check its head, the claim fails.
  - **Against, falsifiable:** Polling delays review and can miss a short-lived head. Measure time from PR update to check result, and test an update during a poll.
  - **Serves:** `F-0001#1`, `#2`, `#5`–`#7`, `#9`. **Strains:** `#4` if the target list has no evidence; `#5` if a missing result is shown as success. **Main risk:** A slow or failed poll becomes a delivery bottleneck.

| Comparison | G1: events | G2: polling |
| --- | --- | --- |
| Trigger | PR event and relay | Scheduled scan |
| New moving part | Event relay | Target list and poll schedule |
| Expected delay | Event delivery | Up to the next successful poll |
| Failure mode to test | Lost or duplicate event | Missed head or excessive delay |
| ADR-0011 boundary | CLI and gate code stay outside | CLI and gate code stay outside |

Neither a required check nor GitHub branch protection can stop a person from *looking at* a failing PR or posting an early review. Review assignment can wait for green checks, and the pilot can audit early reviews, but a literal prohibition on all human review needs an additional access or workflow control. Both runner options also depend on the rule-protection control in question 2: an agent that can change or bypass required checks defeats them.

### 2. Rule-protection control — REQ-003

The current protection record has `require_code_owner_reviews: false` and zero required approvals. The CI restore step protects some scripts, but a PR can change its workflow file; it does not close `F-0001#3`. With one GitHub account shared by agents and Operator, **neither option is a control the agents cannot pass** until authority or credentials are separated.

- **Option P1 — Separate approver and protected paths.** Give agents an identity without rule-administration rights; give an independent human identity control of branch rules and approval. Require its CODEOWNERS review for the complete rule-path set, including `.github/**`, kit rules and checks, target gate configuration, and LAYUP’s own gate source; an independent required check rejects protected-path commits authored by an agent rather than treating approval as permission for agent writes.
  - **For, falsifiable:** The forge can refuse an unapproved rule change. With only agent credentials, try a protected-path PR, a direct push, an altered workflow, and a check-name change; any merge or direct write disproves the claim.
  - **Against, falsifiable:** Path lists drift and human approval adds a wait. Add a new rule path in a test target and measure whether it is protected and how long approval takes.
  - **Serves:** `F-0001#3`, `#5`–`#7`. **Strains:** `#9` if approval depends on one agent harness; `#4` if the path list lacks evidence. **Main risk:** An omitted path or an agent-held administrator credential bypasses the control.

- **Option P2 — Independent policy App with operator authorization.** Keep the gate and path policy outside agent control. A separately administered App refuses every agent-authored protected-path change; a rule update needs an Operator-signed, Git-recorded authorization that names the paths and proposed commit, and the App checks that authorization before it reports success.
  - **For, falsifiable:** Approval is bound to an exact change, not to a general PR approval. Reuse an authorization for another SHA or path; success disproves the claim.
  - **Against, falsifiable:** Key custody, signature checking, and App operation are extra work. Disable the App or its credential in a test target; if a rule change can still merge, the design fails closed only in name.
  - **Serves:** `F-0001#1`, `#3`, `#5`, `#6`, `#9`. **Strains:** `#2` if ordinary target gates need the App; `#4` without recorded key and policy provenance. **Main risk:** The shared GitHub account can still change branch settings unless administration is moved outside agent reach.

| Comparison | P1: separate approver | P2: policy App |
| --- | --- | --- |
| Human act per rule change | Forge approval | Sign exact authorization |
| Main proof | Ownership, identity, protected-path test | Signature, identity, required App check |
| Path-list risk | CODEOWNERS and check coverage | App policy coverage |
| Required prerequisite | Separate approver and administrator credentials | Separate signing and administrator credentials |
| Shared-account situation today | Does not satisfy REQ-003 | Does not satisfy REQ-003 |

An audit must state whether “zero agent writes” means **no landed agent-authored rule change** or **no agent push to any remote ref containing one** (`F-0003#64`). GitHub merge controls address the first. They do not stop an agent from writing a rule file on its own branch; the second reading needs a push-time control outside ordinary CODEOWNERS review.

### 3. Roles and handoff schema — REQ-005

The PSB gives examples of functions but sets no role list (`F-0001#21`). In either option, a target keeps readable records under `runs/<task-id>/`; a Go validator runs outside it. A fixed-header, append-only `handoffs.tsv` indexes each transition, and its referenced Markdown record carries text too long for a TSV cell. The validator must check IDs, required fields, cited files at the named Git commit, allowed transitions, and acceptance evidence—not only field presence.

- **Option H1 — Four delivery roles.** Use **specification agent** (facts to requirements), **architect** (requirements to design and contracts), **implementer** (tests and code), and **verifier** (independent test and review); a specialist acts through the owning role. Each handoff row holds `id`, `task`, `from_role`, `to_role`, `artifact_kind`, `artifact_path`, `commit`, `requirement_ids`, `evidence_path`, and `state`; the Markdown record names the question, decision references, and receiver’s acceptance or refusal.
  - **For, falsifiable:** A small transition graph is easy to validate. Run a pilot task and count transitions with valid records; any missing or invalid transition disproves the claim of full coverage.
  - **Against, falsifiable:** Specialist questions may take an indirect route through the architect. Seed domain, contract, and environment questions; if the correct specialist cannot answer without a false role label, the model does not fit.
  - **Serves:** `F-0001#1`, `#2`, `#6`, `#9`. **Strains:** `#7` if one target domain changes role rules; `#5` if an unvalidated transition is treated as complete. **Main risk:** One broad architect role becomes the human message bus in agent form.

- **Option H2 — Roles by artifact and question kind.** Use **requirements owner**, **design owner**, **contract owner**, **delivery owner**, and **verification owner**; a target maps each question kind to one owner before work starts. The same TSV keys each handoff, but adds `question_kind` and `acceptance_criterion`; the validator checks an allowed artifact-kind-to-owner table, the requirement trace, and an independent verifier for a completed change.
  - **For, falsifiable:** Contract and requirement questions have direct recipients. Seed each kind and check that it reaches one named owner with no Operator routing.
  - **Against, falsifiable:** More roles create more transitions. Count transitions and invalid or bounced handoffs in a pilot; a higher bounce rate than H1 challenges the benefit.
  - **Serves:** `F-0001#1`, `#6`, `#9`. **Strains:** `#4` if owner mappings are guessed; `#7` if the mapping changes baseline rules rather than target content. **Main risk:** Over-specific routing makes a new question kind unowned.

| Comparison | H1: four roles | H2: artifact owners |
| --- | --- | --- |
| Routing key | Delivery stage | Artifact and question kind |
| Schema check that matters | Stage and evidence match | Owner, kind, trace, and evidence match |
| Likely failure | Architect bottleneck | Unowned kind or excess handoffs |
| Target record | TSV event and readable Markdown | TSV event and readable Markdown |

The schema should derive its record kinds from Armature’s task, requirement, decision, and review conventions. Neither a forge comment alone nor a harness-specific memory file is a valid handoff (`F-0001#1`, `#9`).

### 4. Escalation rule — REQ-008

Both options read the approved problem statement, PRD goals and acceptance criteria, recorded funding limit, and prior idea-owner decisions. They write a proposed decision and its classification in the target’s task record. A selected business fork stops the agent until the idea owner decides; that owner also confirms whether the escalation was business-forking. A rejected classification counts as unplanned human input (`F-0001#26`–`#28`). Neither option can infer the meaning of arbitrary prose with certainty; the seeded tests and decision audit remain necessary (`F-0003#57`).

- **Option E1 — Four-axis decision declaration.** Before an agent acts on a choice, it writes a typed proposal with evidence and `budget_change`, `legal_or_compliance_change`, `approved_intent_change`, and `strategic_goal_tradeoff` as `yes`, `no`, or `unknown`. A Go check selects any `yes`, rejects unsupported `no` fields, and routes an `unknown` to a different responsible role for classification; if it stays unknown, work stops without making the decision.
  - **For, falsifiable:** The four triggers directly match `F-0001#26`. Seed one decision of each kind and one internal architecture trade-off; a missed fork or escalated internal trade-off disproves the classification claim.
  - **Against, falsifiable:** An agent can misstate a semantic change as `no`. Audit a random sample against the approved sources; any business fork made without escalation exposes that limit.
  - **Serves:** `F-0001#1`, `#4`–`#6`, `#9`. **Strains:** `#6` if a free-text judgement is presented as machine proof. **Main risk:** A false declaration hides the fork.

- **Option E2 — Approved authority envelope.** The idea owner records permitted budget, intent, goals, and legal or compliance bounds before delivery. A proposed action must cite an applicable bound and show it stays inside it; a missing bound, an exceeded bound, or a trade-off between approved goals selects escalation, while a cited internal technical choice remains with agents.
  - **For, falsifiable:** Decisions outside explicit authority stop without a model call in the engine. Test a proposal just inside and just outside each recorded bound; a wrong result disproves the rule.
  - **Against, falsifiable:** An incomplete envelope can send routine questions to the idea owner. In a pilot, count escalations the owner says are not business-forking; a nonzero count measures that cost as unplanned input.
  - **Serves:** `F-0001#1`, `#4`–`#6`, `#9`. **Strains:** `#4` if bounds have no approved evidence; `#7` if project-specific bounds are mistaken for baseline rules. **Main risk:** Excess false escalations recreate the human routing bottleneck.

| Comparison | E1: declared axes | E2: authority envelope |
| --- | --- | --- |
| Machine reads | Typed proposal and cited evidence | Proposal and approved bounds |
| Uncertainty | Second role, then stop | No matching bound, then stop |
| Main audit target | False `no` declaration | False or missing bound |
| Owner’s final act | Confirm class and decide fork | Confirm class and decide fork |

### 5. Stall procedure — REQ-009, REQ-010

Both options make `no-agreement` and `repeating-without-progress` explicit states, not another retry. A progress event must name a changed acceptance result or a new evidence artifact; a repeated prompt, empty commit, or another run with the same result is not progress. Keep the trigger, attempts, elapsed time, opposing positions, input SHAs, logs, examiner identity and brief, diagnosis, outcome, and Operator package in `runs/<task-id>/` as an append-only event index plus readable record. An open stall with no diagnosis is **incomplete**, not a pass (`F-0003#61`). The limits below are candidates, not PSB values; review-cycle caps in bootstrap mode do not set product retry limits.

- **Option S1 — Fixed circuit breaker.** After two unsuccessful attempts at the same step **or** 30 minutes without a recorded progress event, stop; for a role disagreement, stop after each side has stated its position once and one evidence-based reconciliation attempt fails. Give the records, but not the performers’ reasoning threads, to a fresh session that did none of the stalled work; it writes a diagnosis and either a tested resolution or an Operator package with both positions and next choices.
  - **For, falsifiable:** Fixed limits are easy to test with an injected clock and seeded disagreement. If either case gets a third attempt or passes 30 minutes without entering a stall state, the claim fails.
  - **Against, falsifiable:** A long valid build can trip the timer. Run representative target-stack steps; if a step with genuine new evidence regularly stops, the limit is too small.
  - **Serves:** `F-0001#1`, `#5`, `#6`, `#9`. **Strains:** `#4` if the 30-minute value is adopted without pilot evidence. **Main risk:** A fixed time limit confuses slow progress with repetition.

- **Option S2 — Pre-registered per-step limit.** Before a task starts, its plan records an evidence-based retry count and no-progress duration for each long step, with a system ceiling of two identical-input retries; missing limits block the step from starting. On a limit or `no-agreement`, a fresh session from another harness, where available, examines only the task brief and evidence, writes the diagnosis, and either records a verified recovery or sends the full package to the Operator.
  - **For, falsifiable:** Step-specific durations can allow known slow work without an open-ended loop. Test a fast test step and a slow build step against their pre-registered limits; a limit changed after a failure disproves the claim.
  - **Against, falsifiable:** Setup needs evidence for each duration, and more values can be set badly. Audit the sources and count stops that a fresh examiner calls premature.
  - **Serves:** `F-0001#1`, `#4`–`#6`, `#9`. **Strains:** `#4` if durations are guesses; `#2` if recovery depends on one harness. **Main risk:** The plan itself becomes a way to grant an excessive retry budget.

| Comparison | S1: fixed | S2: per step |
| --- | --- | --- |
| Repeat limit | Two attempts or 30 minutes | Pre-registered count and duration; ceiling on identical-input retries |
| Disagreement | One reconciliation attempt | Stop at the recorded limit or unresolved position |
| Examiner | Fresh session, not a performer | Fresh session, different harness where available |
| Main failure | Premature stop | Bad or costly limit setting |

A fresh examiner is a **session with a role**, not a vote by the original agents. Its record must say what it was given and what it could not examine, as required by the independence standard in *Who may review*. If no independent examination is available, mark diagnosis pending and send that fact with the evidence to the Operator; do not report REQ-010 as met.

### Questions for the Operator

1. For REQ-007, does “reaches human review” mean entry into a **managed review queue**, or must the forge prevent any human from reviewing early? Which actor may install an App and administer required checks?
2. Will agents stop using the Operator’s GitHub credentials, and who will hold separate rule-administration, approval, or signing authority? Does “zero agent writes” include agent feature branches, or only landed history?
3. Which role split should be tested first, and who approves a target’s owner mapping when a question kind has no owner? These are proposed roles, not roles fixed by the PSB.
4. Should the idea owner approve a four-axis declaration rule or an authority envelope? Where is the approved funding limit recorded before either rule runs?
5. Which retry and no-progress limits may the pilot use, and may an independent examiner use the same harness in a fresh session when a second harness is unavailable?
