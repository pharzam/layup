# T-hbw8 — the Operator's decisions on #72

The Operator's answers to the batch "Questions for the Operator — one batch (plan D4)" on #72. The answer comment is copied word for word; the numbers O-76 to O-83 and the author's reading of each follow (R13). If a reading is wrong, the Operator says so on #72 and this file is corrected.

## The answer, word for word

Author: pharzam, 2026-09-27T19:41:27Z, comment 5859201917 on #72.

--- begin
Q-1Definitive decision: LAYUP exists for the target strictly as an external binary whose sole responsibility is orchestrating the end-to-end transformation of the PSB into working software. The target's technology stack and its internal tooling have zero dependency on LAYUP. Anything in the architecture or ADRs that contradicts this clean separation must be removed and aligned with this scenario.Q-2I do not want to create separate GitHub accounts.Is it possible for agents to post comments using distinct visual identifiers (e.g., specific agent badges, icons, signatures, or GitHub App bot identities) while operating under my access level / permissions, where I remain the responsible actor?What viable technical solutions exist for this within GitHub?If this model (agents operating under my account/access but with clearly tagged agent identities) is technically feasible, base all human vs. agent decision-point tracking on this approach.Q-3The smart-if engine must be configurable and parametric rather than hardcoded to a specific model or fixed pattern:During the setup and adaptation phase of a target repository, the operator should be able to select the underlying smart-if provider (which may or may not be JEV).How decisions are made by the smart-if engine—and the authority level granted to it (e.g., shadow mode vs. autonomous delegation)—must be fully parametric. An operator might choose to trust the smart-if engine heavily in one project and keep it tightly restricted in another.Q-4As noted in Q-3, all smart-if intervention points must be parametric:The confidence thresholds, escalation rules, and check frequencies (e.g., pre-task only vs. at every handoff) must be operator-configurable parameters set when running LAYUP against the target, rather than static architectural invariants.Q-5Choose option (b) as the baseline, but ensure this is not hardcoded. It must be exposed as a configurable parameter so the operator can adjust or override the behavior based on the specific harness and environment.Q-6Option (a) is the preferred direction, with option (a) acting as the default configuration.However, please ensure this is not hardcoded. Investigate and implement this so role mapping and ambiguity ownership are parametric/configurable. If a project requires a custom or simplified role matrix, the operator should be able to define it.Q-7These values must not be hardcoded.Both $T$ and $N$ must be fully parametric and editable, allowing the operator to calibrate and refine them across runs based on empirical observation rather than guesswork.Set the default values to $T = 10$ minutes (maximum) and $N = 1$.Q-8This must also not be hardcoded. The learning/routing weight adjustments, boundaries, and trigger policies must be parametric and adjustable by the operator.
--- end

## The decisions and their reading

| No. | Question | Reading |
| --- | -------- | ------- |
| O-76 | Q-1, stack gates | LAYUP is only an external binary that orchestrates the target from the PSB to working software; the target's stack and its own tools never depend on LAYUP. The stack gates of a target are the target's own native tools (for Go: `go vet`, a layout test, a boundary test) and its own CI job, written at setup as part of the Armature copy; they run with LAYUP absent (Invariant 2, finding A2). No LAYUP program, library, script or LAYUP-only file that the target needs to build, test or pass its gates goes into the target. `layup gate` runs the same native commands from outside. This supersedes O-10 and O-11 where they put the stack gates outside the target, and keeps their intent (no LAYUP code in a target). Records that LAYUP writes (the records branch) are data about the run, not a dependency of the target's tools. Anything in an ADR that contradicts this is removed. |
| O-77 | Q-2, human identity | No second GitHub account. Agents act under the Operator's account through a GitHub App's user access token: GitHub shows the Operator's avatar with the app's badge, and the API field `performed_via_github_app` names the app. A comment by the Operator's account with no app is a human decision; a comment through the app is an agent's. The Operator stays the responsible actor. Condition: agent sessions get only the app's token, never the Operator's own token or SSH key; where that is not so, the record says that the separation is by convention only. |
| O-78 | Q-3, the smart-if engine | The decision component has a provider interface; the Operator selects the provider per target at setup (Jev or another). The authority level (for example shadow or autonomous) is a parameter per target and per decision point. |
| O-79 | Q-4, the intervention points | The thresholds, the escalation rules and the check frequency (before each task only, or also at each handoff) are parameters that the Operator sets when running LAYUP against a target, not fixed in the architecture. |
| O-80 | Q-5, harness telemetry | Default (b): a harness that does not report tokens or cannot enforce a spend cap may do paid autonomous work with a wall-clock limit as the proxy, and its telemetry counts as incomplete. A parameter per harness can change this. |
| O-81 | Q-6, roles | Default (a): the seven PSB §2 functions, each ambiguity kind owned by one of them. The role matrix and the owner of each ambiguity kind are a parameter that the Operator can replace per project. |
| O-82 | Q-7, stall values | T and N are parameters, editable per run. Defaults: T = 10 minutes (a maximum), N = 1. |
| O-83 | Q-8, learning | The weight adjustments, their bounds and the trigger policy of the learning loop are parameters. |

**One boundary the author reads into O-79, confirmed by O-84.** The PSB makes some rules invariant: a business-forking decision goes to the idea owner (`F-0001#13`, `#26`); a check that did not run is not a pass (`F-0001#5`); O-71 says a failed or below-threshold smart-if answer goes to a human, never to "pass". The parameters of O-78 and O-79 tune how and when these rules apply (thresholds, the screen frequency, the authority level); they do not switch the rules off. The architecture states this bound.

## O-84: the readings approved, and "smart-if"

The author posted the readings above on #72 (comment 5859227091) and asked the Operator to confirm the boundary. The Operator's answer, word for word (pharzam, 2026-09-27T20:41:41Z, comment 5859648760):

--- begin
The text is accurate and approved. One clarification on smart-if: here, "if" simply means a conditional "if" (an if-condition).

--- end

Reading: the readings of O-76 to O-83 and the boundary are approved. The smart-if engine is a conditional branch point in the orchestrator's flow: at a named point, the flow takes one branch or another by the engine's answer, as an `if` statement does; it does not write text or plan work.

## O-85 and O-86: the budget overrun, and the full evaluation before the rewrite

The Operator's answer to "Second public-solution search — result, and the revised question" on #72, word for word (pharzam, 2026-09-28T12:23:57Z, comment 5869765615):

--- begin
The budget overrun is approved.

Regarding the final question, proceed with option (a). However, do not keep the hands-on evaluation time-bounded or limited—ensure it is a thorough, full evaluation of the reuse candidates before moving forward with the rewrite plan.
--- end

| No. | Decision | Reading |
| --- | -------- | ------- |
| O-85 | The budget overrun is approved | The branch past the Budget maximum of the plan review (3,200 lines over 36 files) is approved as it stands at `22e65b3` (37 files, +4,055 −32). The rewrite gets its own budget from the plan review of the new plan. |
| O-86 | Option (a), with a full evaluation | LAYUP keeps its own Go orchestrator. Before the rewrite plan, each reuse candidate of `search-v2/summary.md` gets a thorough, hands-on evaluation, with no time bound: Gas Town and Beads first, then GNAP, Spec Kitty, AI-SDLC, Paperclip's budget and approval model, and the review and done-check components; each ends with use, borrow the pattern, or reject, with the evidence. Then the new plan for the rewrite, by the fixed method. |

## O-87 to O-91: the answers to the batch of plan v2

The Operator's answer to "Questions for the Operator — one batch (plan v2)" on #72 (comment 5871619023). The answer was given in the author's Claude Code session on 2026-09-29, not on #72; the author copied it to #72 word for word. Word for word:

--- begin
q-9 a , q-10 a , q-11 a, q-12 , i ment basecamp shapup method may be help us, JUST check I , Q-13 I'll reinstall it by myself
--- end

| No. | Question | Reading |
| --- | -------- | ------- |
| O-87 | Q-9, a review record and authorship | Option a: a review record in Git (for example `deep-check-fable.md`) does not make its writer an author under Bootstrap mode rule 4. Claude Fable 5.1 on `claude` may do the gate round of this task. |
| O-88 | Q-10, Haiku in the evaluation | Option a: the evidence of the runs that used Claude Haiku 4.5 as the model inside the candidates stays valid. The use stays a reported deviation from ADR-0012 part 3; no later run of this task uses a model on its "not used" list. |
| O-89 | Q-11, the budget growth of the evaluation | Option a: the growth `22e65b3..946edfc` (22 files, +2,897) is approved as it stands. The Budget maximum of plan v2 (7,200 lines over 64 files) counts from `946edfc`. |
| O-90 | Q-12, "Shape Up" | "Shape Up" names Basecamp's Shape Up method (Ryan Singer, the book at `https://basecamp.com/shapeup`), not a tool. "JUST check it": evaluate the method, as one more candidate, for what it can give LAYUP; it is not a decision to adopt it. The result is `evaluation/19-shape-up.md`. |
| O-91 | Q-13, `~/.local/bin/agy` | The Operator reinstalls `agy` personally. The author runs no `agy` session in this task until the Operator says on #72 that it is reinstalled. |

## O-92: plan v2 confirmed; a fresh session; the agent badge from now

The Operator's answer to "Plan v2 as amended — for the Operator's confirmation" on #72 (comment 5872132790), given in the author's Claude Code session on 2026-09-29 and copied to #72 word for word. Word for word:

--- begin
confirmed, Just I want to start it in fresh session and apply the Hanness-agent badge from NOW
--- end

| No. | Decision | Reading |
| --- | -------- | ------- |
| O-92 | Plan v2 confirmed, with two conditions | (1) Plan v2 as amended is confirmed; the rewrite starts at step 0. (2) The rewrite starts in a fresh session, not in the session that wrote the plan. (3) "The harness-agent badge" is the marker of O-77: from now, every GitHub write of an agent session (issue comments, pull requests, pushes) goes through the LAYUP GitHub App's user access token, so that GitHub shows the App's badge and the API field `performed_via_github_app` names the App. The App does not exist yet; the Operator creates and installs it. Until the author has its token, the author writes nothing more to GitHub except the record of this decision and the setup request; the fresh session starts only with the token. The O-77 condition stays: while an agent session can also reach the Operator's own `gh` login, the separation is by convention only, and each record says so. |
