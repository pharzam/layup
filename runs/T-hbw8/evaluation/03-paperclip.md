# Evaluation: Paperclip

| Field | Value |
| ----- | ----- |
| Repository | https://github.com/paperclipai/paperclip |
| Pinned commit | 0f14d261233c545aa6a8a38ec253c498a5130fff (2026-09-27T06:42:19-05:00, "fix(runtime): allow bounded large untracked workspace snapshots (#14194)") |
| License | MIT, "Copyright (c) 2025 Paperclip AI" (LICENSE file) |
| Language and needs | TypeScript. Node.js server (Express) and a React UI, pnpm workspace (about 1,300 packages). Database: PostgreSQL. By default the server starts an embedded PostgreSQL 18 (npm `embedded-postgres`) on 127.0.0.1:54329; an external PostgreSQL is optional. The dev start also builds a Rust binary (`paperclip-runnerd`, Rust 1.97.1), because the experimental "native runner" is on by default. The Claude adapter needs the `claude` CLI and, on its default path, the `claude-agent-acp` bridge (Node >= 24.11). No cloud account is necessary. Telemetry to a Paperclip endpoint is on by default (opt-out). |
| Evaluated by | Claude Opus 5.5 (`claude-opus-5-5`) on Claude Code, 2026-09-28 |
| Elapsed | about 35 minutes wall clock |
| Agent runs and cost | 3 Paperclip heartbeat runs of Claude Code, model `claude-haiku-4-5-20251001`. Claude reported USD 0.0715, 0.0644 and 0.0366 (the third run was cancelled by the budget stop), total USD 0.1725. Paperclip recorded all three as 0 cents (subscription login). One direct harness probe outside Paperclip: USD 0.016. |

## Verdict

**Borrow the pattern; reject as a component to run.** The budget model (scope, window, warn percent, hard stop, incident, override approval, check before each run, cancel of running work), the typed approval record, the atomic issue checkout and the review-round limit are good, working designs. But all state is in a PostgreSQL database outside the target's Git (I1, I2). In its default local mode any local process without a key acts as the human board, and our agent did so (I3). A subscription login records every run as 0 cents, so the budget stop never trips and reads "ok" (I5).

## What it is (from the code)

Paperclip is a "control plane for AI-agent companies" (`AGENTS.md`). A company has agents in an org chart, projects, issues, goals, budgets and approvals. A heartbeat scheduler wakes agents; an adapter starts a harness (Claude Code, Codex, Cursor, Gemini, OpenCode and others) for one run; the agent works on its checked-out issue and calls the Paperclip REST API back.

Budget enforcement, `server/src/services/budgets.ts`:

- Policy row per scope (`company`, `agent`, `project`), metric and window (`calendar_month_utc` or `lifetime`), with `warnPercent` (default 80), `hardStopEnabled` (default true): https://github.com/paperclipai/paperclip/blob/0f14d261233c545aa6a8a38ec253c498a5130fff/packages/db/src/schema/budget_policies.ts#L4
- The observed amount is the sum of `cost_events.cost_cents` in the window. A metric other than `billed_cents` always gives 0: https://github.com/paperclipai/paperclip/blob/0f14d261233c545aa6a8a38ec253c498a5130fff/server/src/services/budgets.ts#L147
- After each cost event, `evaluateCostEvent` makes a soft incident at the warn percent, and at 100 % a hard incident, pauses the scope and cancels its work: https://github.com/paperclipai/paperclip/blob/0f14d261233c545aa6a8a38ec253c498a5130fff/server/src/services/budgets.ts#L693
- A hard incident creates an approval of type `budget_override_required` with the numbers and the guidance "Raise the budget and resume the scope, or keep the scope paused.": https://github.com/paperclipai/paperclip/blob/0f14d261233c545aa6a8a38ec253c498a5130fff/server/src/services/budgets.ts#L385
- Before each run, `getInvocationBlock` recomputes company, agent and project spend and refuses the start: https://github.com/paperclipai/paperclip/blob/0f14d261233c545aa6a8a38ec253c498a5130fff/server/src/services/budgets.ts#L718
- `resolveIncident` raises the amount (it must exceed the observed spend) and resumes, or dismisses; it marks the linked approval approved or rejected: https://github.com/paperclipai/paperclip/blob/0f14d261233c545aa6a8a38ec253c498a5130fff/server/src/services/budgets.ts#L866
- Cancel of running work for the scope: https://github.com/paperclipai/paperclip/blob/0f14d261233c545aa6a8a38ec253c498a5130fff/server/src/services/heartbeat.ts#L29221
- A second, per-agent daily cap in the heartbeat policy (`maxDailyRuns`, `maxDailyCostCents`): https://github.com/paperclipai/paperclip/blob/0f14d261233c545aa6a8a38ec253c498a5130fff/server/src/services/heartbeat.ts#L16578

Cost ledger:

- Table `cost_events`: company, agent, issue, project, goal, heartbeat run, billing code, provider, biller, `billingType`, `costStatus` (`reported` or `unpriced`), model, input, cached input and output tokens, `cost_cents`, time. No latency and no duration column: https://github.com/paperclipai/paperclip/blob/0f14d261233c545aa6a8a38ec253c498a5130fff/packages/db/src/schema/cost_events.ts#L9
- The heartbeat writes one cost event per run, at the end of the run, from the adapter result: https://github.com/paperclipai/paperclip/blob/0f14d261233c545aa6a8a38ec253c498a5130fff/server/src/services/heartbeat.ts#L19725
- `subscription_included` billing gives 0 cents, whatever the reported USD: https://github.com/paperclipai/paperclip/blob/0f14d261233c545aa6a8a38ec253c498a5130fff/server/src/services/heartbeat.ts#L5213 . The Claude adapter sets "subscription" when `ANTHROPIC_API_KEY` is absent: https://github.com/paperclipai/paperclip/blob/0f14d261233c545aa6a8a38ec253c498a5130fff/packages/adapters/claude-local/src/server/execute.ts#L165
- Duration: `heartbeat_runs.started_at` and `finished_at`; the issue cost summary adds `runtimeMs`: https://github.com/paperclipai/paperclip/blob/0f14d261233c545aa6a8a38ec253c498a5130fff/packages/db/src/schema/heartbeat_runs.ts#L32
- An agent can report its own cost events (only for itself): https://github.com/paperclipai/paperclip/blob/0f14d261233c545aa6a8a38ec253c498a5130fff/server/src/routes/costs.ts#L114

Approval workflow:

- Table `approvals` (type, requester agent or user, status, JSON payload, decision note, `decided_by_user_id`, time), `approval_comments`, `issue_approvals` (link to issues): https://github.com/paperclipai/paperclip/blob/0f14d261233c545aa6a8a38ec253c498a5130fff/packages/db/src/schema/approvals.ts#L5
- Four types only: `hire_agent`, `approve_ceo_strategy`, `budget_override_required`, `request_board_approval`: https://github.com/paperclipai/paperclip/blob/0f14d261233c545aa6a8a38ec253c498a5130fff/packages/shared/src/constants.ts#L689
- Statuses pending, revision_requested, approved, rejected, cancelled; the resolve is a conditional update, so a double decision cannot apply twice: https://github.com/paperclipai/paperclip/blob/0f14d261233c545aa6a8a38ec253c498a5130fff/server/src/services/approvals.ts#L45
- Approve and reject need a "board" actor: https://github.com/paperclipai/paperclip/blob/0f14d261233c545aa6a8a38ec253c498a5130fff/server/src/routes/approvals.ts#L286 . After approval the server wakes the requesting agent.
- A hire approval creates the agent and its monthly budget policy: https://github.com/paperclipai/paperclip/blob/0f14d261233c545aa6a8a38ec253c498a5130fff/server/src/services/approvals.ts#L154
- Issue-level execution policy with `review` and `approval` stages and `maxReviewRounds`: https://github.com/paperclipai/paperclip/blob/0f14d261233c545aa6a8a38ec253c498a5130fff/packages/shared/src/validators/issue.ts#L434

Identity: in the default `local_trusted` mode each request without a bearer credential is the board user `local-board`: https://github.com/paperclipai/paperclip/blob/0f14d261233c545aa6a8a38ec253c498a5130fff/server/src/middleware/auth.ts#L227 . The CSRF guard lets `local_implicit` mutations through: https://github.com/paperclipai/paperclip/blob/0f14d261233c545aa6a8a38ec253c498a5130fff/server/src/middleware/board-mutation-guard.ts#L92 . The mode default is `local_trusted`: https://github.com/paperclipai/paperclip/blob/0f14d261233c545aa6a8a38ec253c498a5130fff/server/src/config.ts#L175 . In `authenticated` mode a board actor needs a login session or a board key; an agent uses a hashed agent key (`agent_api_keys`) and an optional run header `x-paperclip-run-id`.

Atomic checkout: one conditional `UPDATE issues ... WHERE status IN (expected) AND (assignee IS NULL OR same assignee without a run lock) AND execution lock free`: https://github.com/paperclipai/paperclip/blob/0f14d261233c545aa6a8a38ec253c498a5130fff/server/src/services/issues.ts#L11412 . A stale lock of a finished run is adopted inside a transaction with `SELECT ... FOR UPDATE`. Route: https://github.com/paperclipai/paperclip/blob/0f14d261233c545aa6a8a38ec253c498a5130fff/server/src/routes/issues.ts#L15076

Heartbeat scheduler: an in-process timer, every `HEARTBEAT_SCHEDULER_INTERVAL_MS` (default 30 s, minimum 10 s): https://github.com/paperclipai/paperclip/blob/0f14d261233c545aa6a8a38ec253c498a5130fff/server/src/config.ts#L366 . `tickTimers` wakes each agent whose policy has `enabled` and `intervalSec`: https://github.com/paperclipai/paperclip/blob/0f14d261233c545aa6a8a38ec253c498a5130fff/server/src/services/heartbeat.ts#L29620 . Wakes also come from assignment, comments, approvals and "automation" (continuations).

Stalls and recovery: each run gets a liveness state (`completed`, `advanced`, `plan_only`, `empty_response`, `blocked`, `failed`, `needs_followup`). Bounded automatic continuations are fixed constants: `DISPOSITION_REPAIR_MAX_ATTEMPTS = 5` (https://github.com/paperclipai/paperclip/blob/0f14d261233c545aa6a8a38ec253c498a5130fff/server/src/services/recovery/disposition-repair.ts#L18), `DEFAULT_MAX_LIVENESS_CONTINUATION_ATTEMPTS = 2` (https://github.com/paperclipai/paperclip/blob/0f14d261233c545aa6a8a38ec253c498a5130fff/server/src/services/recovery/run-liveness-continuations.ts#L9). A review stage counts consecutive agent "changes requested" rounds; at `maxReviewRounds` (default 3) it escalates to the responsible human, and a human decision resets the count: https://github.com/paperclipai/paperclip/blob/0f14d261233c545aa6a8a38ec253c498a5130fff/server/src/services/issue-execution-policy.ts#L69 . A "stranded assigned issue" becomes `blocked` with a board escalation. An optional "task watchdog" agent reads a stopped issue tree and decides if the stop is legitimate (`doc/TASK-WATCHDOG.md`).

Claude Code adapter, `packages/adapters/claude-local`: the default engine is ACP (`claude-agent-acp`), not `claude --print`: https://github.com/paperclipai/paperclip/blob/0f14d261233c545aa6a8a38ec253c498a5130fff/packages/adapters/claude-local/src/server/acp.ts#L84 . `dangerouslySkipPermissions` defaults to true: https://github.com/paperclipai/paperclip/blob/0f14d261233c545aa6a8a38ec253c498a5130fff/packages/adapters/claude-local/src/server/execute.ts#L434 . The CLI engine builds `--print --output-format stream-json --verbose`, `--model`, `--max-turns`, and `extraArgs`: https://github.com/paperclipai/paperclip/blob/0f14d261233c545aa6a8a38ec253c498a5130fff/packages/adapters/claude-local/src/server/execute.ts#L880 . The adapter has no money cap per run; `extraArgs` is documented as "additional CLI args", so on the ACP default our `--max-budget-usd 0.30` had no effect that we could see. Usage and `total_cost_usd` are parsed from the final result: https://github.com/paperclipai/paperclip/blob/0f14d261233c545aa6a8a38ec253c498a5130fff/packages/adapters/claude-local/src/server/parse.ts#L116 . The adapter writes skills under `$HOME/.claude/skills`.

Where state lives (all under `PAPERCLIP_HOME`, default `~/.paperclip/instances/default/`): `db/` (the PostgreSQL data directory: companies, agents, issues, comments, approvals, budgets, cost events, heartbeat runs, activity log, config revisions, more than 100 tables in `packages/db/src/schema/`), `data/run-logs/`, `data/workspace-operation-logs/`, `data/backups/`, `secrets/` (a signing key), `workspaces/<agent-id>/`, `companies/<id>/` (ACP state). The target repository gets only the agent's commits, and Claude's own `.claude/settings.local.json`.

Export: `POST /api/companies/:id/export` gives a file package: `COMPANY.md`, `agents/<slug>/AGENTS.md`, `projects/<slug>/PROJECT.md`, `tasks/<slug>/TASK.md`, `.paperclip.yaml`. Include flags are only company, agents, projects, issues, skills: https://github.com/paperclipai/paperclip/blob/0f14d261233c545aa6a8a38ec253c498a5130fff/packages/shared/src/types/company-portability.ts#L7

## What we ran

Setup (all under `EVAL_WORK/paperclip`, host `$HOME` untouched):

```
git clone --depth 1 https://github.com/paperclipai/paperclip.git   # a full clone failed twice: "curl 56 Recv failure"
HOME=$W/home pnpm install                                          # Done in 1m 6.4s
curl https://sh.rustup.rs | RUSTUP_HOME=$W/rust/rustup CARGO_HOME=$W/rust/cargo sh -s -- --default-toolchain 1.97.1
HOME=$W/home PAPERCLIP_HOME=$W/pchome PAPERCLIP_TELEMETRY_DISABLED=1 DO_NOT_TRACK=1 HOST=127.0.0.1 PORT=3100 pnpm dev:once
```

- First start without Rust failed: `sh: cargo: command not found ... paperclip runner native binary build failed`. The native runner is on by default (`enableNativeRunner ?? true` in `server/src/services/instance-settings.ts:227`).
- Claude Code under a sandbox `HOME` said "Not logged in". A symbolic link from the sandbox `Library/Keychains` to the user's keychain fixed the login read; no token was copied or printed.
- The log said: `Migrating database via embedded-postgres@54329`, `Applying 283 pending migration(s)`, `Server listener bound on 127.0.0.1:3100`. Health: `"deploymentMode":"local_trusted"`.
- Toy target: `toy/` (module `example.com/toy`, package `calc`, one test, `check.sh` = `go vet ./... && go test ./...`), a local bare remote `toy-remote.git`.

Org chart and budgets (curl, no credential):

```
POST /api/companies {"name":"ToyCo","budgetMonthlyCents":500}      -> company policy 500 cents created
POST /api/companies/$C/agents Ada (ceo, process), Linus (cto, reportsTo Ada), Grace (engineer, claude_local,
     model claude-haiku-4-5-20251001, cwd toy, maxTurnsPerRun 20, timeoutSec 600, budgetMonthlyCents 100)
GET  /api/companies/$C/org -> Ada > Linus > Grace
```

Atomic checkout: two parallel `POST /api/issues/$I/checkout` for two agents on one `todo` issue.

```
a86239d1-... 409  {"error":"Issue checkout conflict"}
03ba6f2c-... 200  {"status":"in_progress","assigneeAgentId":"03ba6f2c-..."}
```

Approval gate, with an agent key for Linus:

```
agent POST /companies/$C/approvals request_board_approval  -> status "pending"
agent POST /approvals/$AP/approve                           -> 403 {"error":"Board access required"}
agent POST /approvals/$AP/comments "approved by me"         -> stored with authorAgentId; status unchanged
no credential, same host: POST /approvals/$AP/approve      -> {"status":"approved","decidedByUserId":"local-board"}
```

One ticket through a heartbeat: TOY-2 "Add Sub to package calc", assigned to Grace.

- Run 1 (`invocationSource: assignment`), 12:54:30 to 12:55:09: the agent checked out the issue, added `Sub` and a table test, ran `./check.sh`, and committed `d7832e4 calc: add Sub` in `toy/`. It did not set the issue status. Liveness: "advanced", "1 workspace operation(s), 11 tool/action event(s)".
- Run 2 (`automation`, wake reason `issue_disposition_repair`), 12:55:09 to 12:55:44: the platform woke the agent again to set a disposition. The agent ran `curl -X PATCH http://127.0.0.1:3100/api/issues/<id> -d '{"status":"done"}'` with no `Authorization` header. The activity log recorded this as `actorType: user, actorId: local-board`, with no agent and no run id.
- Ledger: `GET /issues/$I2/cost-summary` -> `{"costCents":0,"inputTokens":33341,"cachedInputTokens":428951,"outputTokens":4379,"runCount":2,"runtimeMs":74806}`. `GET /costs/summary` -> `"spendCents":0,"utilizationPercent":0`. Claude reported USD 0.136 for the two runs.

Hard stop:

```
POST cost-events (manual, metered_api, 14 cents)          -> agent policy 14/100 "ok"
PATCH /agents/$ENG/budgets {"budgetMonthlyCents":10}       -> policy "hard_stop", agent "paused" pauseReason "budget",
                                                              incident open, approval budget_override_required "pending"
POST /agents/$ENG/heartbeat/invoke                        -> 409 "Agent is paused because its budget hard-stop was reached."
POST /agents/$ENG/resume  (board)                         -> agent "idle"
POST /agents/$ENG/heartbeat/invoke                        -> 409 "Agent cannot start because its budget hard-stop is still exceeded."
POST budget-incidents/$ID/resolve raise_budget_and_resume 20 -> incident "resolved", approval "approved" by local-board
```

Stop of a running run: a new ticket TOY-4 started run 3 at 12:56:46. At 12:56:55 we posted a 10-cent metered cost event (24 >= 20). Run 3: `"status":"cancelled","error":"Cancelled due to budget pause"`, finished 12:56:56.8; agent `paused`. No `claude` child process stayed. The cancelled run still reported USD 0.0366 of usage.

Export: `POST /api/companies/$C/export` gave 14 files. `TASK.md` has the name, assignee, project and description; it has no status, comments, runs, costs or approvals. `.paperclip.yaml` has adapter config and `budgetMonthlyCents: 20`.

Stopped: the dev runner, the server (tsx/node), esbuild, and the embedded PostgreSQL. After stop, no process with `EVAL_WORK/paperclip` in its command line, and no listener on 3100, 13100 or 54329. No Docker container was used.

## In-Scope items S1–S12

| Item | Mark | Evidence |
| ---- | ---- | -------- |
| S1 Problem Statement Quality | no | No component reads a problem statement for gaps. Goals are free text. |
| S2 Reproducible Discipline Setup | no | A company package (`COMPANY.md`, `AGENTS.md`, `.paperclip.yaml`) can be exported and imported, so an org can be repeated. No Armature baseline, no evidence per value. |
| S3 Rule Protection | partly | Approve, budget and policy routes need a board actor (`assertBoard`), and the agent key got 403. But the default `local_trusted` mode makes any keyless local call the board, and our agent used that path. Target-repository rules are not protected; agents run with `--dangerously-skip-permissions` by default. |
| S4 Stack-Dependent Gates | no | No gate runs on a change. The agent ran `./check.sh` only because the ticket said so. |
| S5 Role Handoffs | partly | Typed issues, one assignee, sub-issues, execution-policy stages (`review`, `approval`) with typed decisions. No schema based on Armature conventions. |
| S6 Autonomous Clarification | no | Agents can comment or assign sub-issues to other agents, but no mechanism routes a question by kind to its owner role and records an accepted answer. |
| S7 Verification on Every Change | no | Review stages are agent or human judgements. No deterministic check before a merge. |
| S8 Human-on-the-Loop | partly | Four approval types; budget overrides go to the board; review loops escalate to a human after `maxReviewRounds`. No escalation rule selects business-forking decisions; an agent decides itself to ask (`request_board_approval`). |
| S9 Stall Resolution | partly | Liveness state per run, bounded continuations (fixed constants), review-round limit, stranded issue to `blocked` with board escalation, optional watchdog agent. No recorded diagnosis in Git; disagreement outside a review stage is not detected. |
| S10 Cost Visibility | partly | Per run: tokens by class, model, provider, cents, cost status; per issue: `runtimeMs`. No latency. Not in the project repository. Subscription runs record 0 cents. |
| S11 Specification Synthesis | no | Issues trace to goals ("every task traces back to the mission"), not to problem-statement text with an acceptance criterion. |
| S12 Harness-Agent Neutrality | partly | Twelve adapter packages (Claude, Codex, Cursor local and cloud, Gemini, Grok, Hermes and Hermes gateway, Kimi, OpenClaw gateway, OpenCode, Pi), plus built-in `http` and `process` adapters. One instructions bundle per agent. No rule that a second harness checks the first. |

## Invariants I1–I9

| Invariant | Effect | Reason |
| --------- | ------ | ------ |
| I1 Git is the system of record | conflicts | Issues, comments, approvals, budgets, cost events, runs and the activity log live only in PostgreSQL under `PAPERCLIP_HOME`. The export carries definitions only (company, agents, projects, task text), not decisions or costs. |
| I2 Independent repository | conflicts | The code in the target stays plain Git (a normal commit), so a human can continue the code. The task state, decisions and cost history are lost without the server and its database. |
| I3 Agents cannot change rules | conflicts | Shown: the agent's keyless `curl` was logged as the human board; a keyless call approved an approval. In `authenticated` mode the board routes are closed to agent keys, but the agent still has full permissions on the target. |
| I4 No value without evidence | conflicts where its defaults are used | Built-in values with no stated evidence: `warnPercent` 80, `DEFAULT_MAX_REVIEW_ROUNDS` 3, disposition repair 5 attempts, liveness continuations 2, scheduler 30 s. |
| I5 Inactive check is not a pass | conflicts | A budget on a subscription login shows "ok", 0 % utilization, after USD 0.17 of reported spend. A policy with a metric other than `billed_cents` always observes 0. Unpriced tokens count as 0 cents. |
| I6 Deterministic over LLM | supports | Budget stop, checkout lock and round limit are plain code and SQL. |
| I7 Domain changes content only | neutral | No rule layer. |
| I8 Pinned Armature | neutral | No Armature. |
| I9 Harness replaceable | supports | The adapter is a per-agent setting; state does not live in a harness. (Paperclip itself becomes the dependency.) |

## Deep-check findings it answers

These are the findings where a Paperclip mechanism gives LAYUP a usable answer or a tested pattern.

- **Author-7, Fable-M11 (budget before tasks exist):** budgets attach to scopes that exist (company at start, then project, agent) with a window, not to tasks. Pattern: take one total at Intake, add scope budgets when Plan makes them. Raising a budget is a typed approval with the old and observed amounts.
- **Sol-11 (unknown token data in the money stop):** answered as a warning, not as a fix. Paperclip has the right fields (`costStatus: unpriced`, `billingType`), but it adds 0 cents for unknown or subscription cost, and the stop never trips. LAYUP must treat unknown spend as unknown and stop or escalate (O-80).
- **Sol-12 (incomplete telemetry):** partly. Tokens are always recorded, even when cents are 0; `runtimeMs` per issue. No latency field. The same completeness gap remains.
- **Fable-N3 (the unit of an action):** the unit is one heartbeat run (one adapter invocation), with `usageSource: per_run`. Pattern for LAYUP's "action = one session".
- **Fable-M1 (agent token has orchestrator rights):** pattern only. Two principal types (board, agent key) and board-only routes for approve, budget and policy. Paperclip's own default mode breaks it, as the run showed. LAYUP must never have a keyless trusted path.
- **Fable-M3, Sol-5 (who did it):** pattern: the actor comes from the credential that the server verified (hashed agent key, run id header), never from text the client writes. Our run showed the failure when a keyless path exists.
- **Sol-30 (a review is not an issue comment):** pattern: a decision is one typed record (`approvals` row with `decided_by_user_id`, status, time); a comment on it never changes its status (our agent's "approved by me" did nothing). No path accepts a GitHub comment as a decision.
- **Author-1, Fable-M4, Sol-22 (disagreement never stalls):** answered for review stages. Consecutive agent "changes requested" rounds on one stage are counted; at `maxReviewRounds` the stage goes to a human; a human decision resets the count. LAYUP needs a fresh diagnosis before the human (Fable-M6).
- **Fable-M6 (stall without diagnosis):** partly. The optional task watchdog is a separate agent that reads a stopped tree and decides. It is not a required step before the human.
- **Author-10, Sol-4, Fable-N5 (dead-man recovery):** pattern: recovery at the next start and on each scheduler tick marks stranded issues `blocked` and escalates to the board with a recovery record. It runs only while the server runs.
- **Fable-M22 (parameters the Operator can set):** answered for the interface. Budget amount, warn percent, hard stop, heartbeat interval, concurrency, daily run and cost caps, and `maxReviewRounds` are set through the REST API at any time, and each change is an activity-log row with its actor. The recovery limits are fixed constants, so not all of them are parameters.
- **Sol-3, Fable-M12 (one writer):** pattern: every write goes through one server with an activity-log row. Not in Git.

## Parts and their verdicts

| Part | Verdict | How LAYUP would take it | Reason |
| ---- | ------- | ----------------------- | ------ |
| Whole product next to LAYUP | reject | none | Second system of record in PostgreSQL (I1, I2); keyless board in the default mode (I3); heavy stack (Node 24, PostgreSQL, Rust build, about 1,300 npm packages, telemetry on by default). |
| Budget policy, incident and override approval | borrow the pattern | Git records: a budget file per scope and window, an incident row, an override decision as an issue comment copied into Git; check before each session and after each cost row | Worked in the run: pause, 409 on start, 409 after a manual resume, resolve only above observed spend. Fix its defect: unknown or subscription spend is not 0. |
| Cancel of running work at the hard stop | borrow the pattern | `layup run` kills the session process group when a cost row crosses the stop | Cancelled in about 2 s in the run. Only useful when cost arrives during a run; Claude reports cost at the end, so one run can overshoot. Add a per-session cap before start. |
| Cost ledger schema | borrow the pattern | TSV columns: task, role, run, harness, provider, biller, billing type, cost status, model, token classes, cents, plus latency and duration | Clear columns; `costStatus` and `billingType` make unknown cost visible. Missing latency. |
| Approval record and statuses | borrow the pattern | One typed decision record per planned point (pending, revision requested, approved, rejected, cancelled), decided by a verified human | Conditional update prevents a double decision. Comments are separate from decisions. |
| Actor model (board or agent key, run header) | borrow the pattern | Map to O-77: human = Operator account without the App; agent = App token; never a keyless path | Shown to fail in `local_trusted`; the pattern holds only with no unauthenticated path. |
| Atomic issue checkout | borrow the pattern | A claim as a compare-and-set in Git (for example a claim ref pushed with `--force-with-lease`, or an issue assignment check) | One 200 and one 409 in the race test. |
| Review-round limit and escalation | borrow the pattern | Parameter K for repeated rejections between the same two roles; then a fresh examiner, then the Operator | Directly answers the disagreement stall class. |
| Liveness states and bounded continuations | borrow the pattern | Classify each session result; limit automatic continuations by a parameter, not a constant | Run 2 was an automatic "disposition repair" wake; useful, but its limit is hard-coded. |
| Heartbeat scheduler | reject | none | An in-process 30 s timer in a long-lived server; LAYUP runs as a binary per phase. |
| Claude Code adapter | reject | none (LAYUP starts `claude -p` itself) | Default engine ACP; no money cap per run; `--dangerously-skip-permissions` by default; writes into `$HOME/.claude`. The stream-json usage parse is a small, known pattern. |
| Org chart (`reportsTo`, role enum) | borrow the pattern | A role matrix file with a "reports to" column (O-81) | Simple and enough for escalation paths. |
| Company export package | reject as a format dependency | none | Leaves out decisions, costs, runs and status; not a record of the work. |

## Where the searchers were wrong or incomplete

- **Database.** Both said the database is not stated or is PostgreSQL by GNAP's table. The code uses PostgreSQL: an embedded PostgreSQL 18 by default (`server/src/config.ts:124`, `embedded-postgres` in `server/package.json:85`), or an external one. The repository's own `AGENTS.md` says "embedded PGlite"; that text is stale.
- **"Monthly budgets per agent. When they hit the limit, they stop" (B), "enforces per-agent monthly spend" (A).** True only for metered billing. With a Claude subscription login each run costs 0 cents in the ledger, so the stop never trips. The stop acts between runs; a run's cost arrives at its end, so one run can overshoot. Budgets exist for company, agent and project, with a monthly or lifetime window.
- **"Governance: who can do what" (B, S3 partly).** In the default mode a keyless local call is the board. Our agent used it, and the audit log named the human.
- **Approvals (B, S8).** Exactly four approval types. The approval is a database record; it cannot be an issue comment or a Git record.
- **Cost visibility (B, S10).** Tokens and duration are recorded; latency is not. Nothing goes to the project repository.
- **Harness neutrality (B, S12).** Confirmed: twelve adapter packages. Not stated by either searcher: the Claude adapter defaults to ACP, where `extraArgs` do not apply.
- **Not found by either searcher:** the atomic checkout, the review-round limit with human escalation, the budget incident with an override approval, the export package and its gaps, the Rust build needed at start, and the default telemetry to a Paperclip endpoint.

## Limits of this evaluation

- We ran the dev runner from source at the pinned commit, not the npm release. The clone was shallow (depth 1), so we did not read the history.
- We did not run `authenticated` mode. The claim that agent keys cannot approve there rests on the code and on the 403 that the agent key got in local mode.
- The host has a Claude subscription login and no API key, so no real run was metered. We triggered the hard stop with manual metered cost events; the mid-run cancel is real, but its trigger was manual.
- We did not test the Codex or other adapters, the CLI engine of the Claude adapter, the task watchdog, the GitHub channel (`chat-github-*.ts`), or company import.
- We did not measure how often the automatic continuations fire over many tickets; we saw one extra run per ticket.
- `heartbeat.ts` has about 29,800 lines and `issues.ts` more than 11,000; we read the paths named above, not all of the code.
