# Evaluation: AgentPlane

| Field | Value |
| ----- | ----- |
| Repository | https://github.com/basilisk-labs/agentplane |
| Pinned commit | 81fc89167d9d7655bd29664849af6fa2eb4bb054 (2026-09-27T19:09:19Z; 13,851 commits; package version 0.7.12-beta.1) |
| License | MIT, from the LICENSE file |
| Language and needs | TypeScript monorepo (Bun build, Node.js 24 or later at run time), npm package `agentplane`. No server and no database by default (a local `node:sqlite` cache). Optional `cloud` backend to `https://sync.agentplane.cloud` and a `redmine` backend. The managed runner drives Codex, a custom command or Hermes. Installs Git hooks and files into the target. 82 stars, last push 2026-09-28T10:09:37Z (API GET). |
| Evaluated by | Claude Opus 5.5 on Claude Code, 2026-09-28 |
| Elapsed | About 35 minutes of wall clock for all three candidates in one session (12:31 to 13:06 UTC); I did not time each candidate apart |
| Agent runs and cost | 6 Claude Code runs (`claude -p`, claude-haiku-4-5-20251001, `--max-budget-usd 0.60` each), 5 billed: 0.1678, 0.0972, 0.1873, 0.1620, 0.1470 USD; one failed at the API with no cost. Total 0.761 USD. |

## Verdict

Borrow the pattern. AgentPlane is the most complete Git-native control plane of the three, and several of its mechanisms answer deep-check findings directly: a single supervisor that alone makes formal transitions, typed and bound agent results, a plan dataflow check, telemetry with an explicit `unavailable` state, and a policy digest that stops a task when the rules change. It cannot be LAYUP's base or its records format: part of its state lives in `.git/agentplane/` and is never pushed, it installs hooks and files that the target then needs (O-76), it cannot run `go` or `make` as a check, and it reads tokens from Codex only.

## What it is (from the code)

Paths are at the pinned SHA; `A` = `packages/agentplane/src`, `C` = `packages/core/src`.

- **One writer of formal state.** The CLI issues a "work order" to an agent and admits a typed result. "A semantic episode cannot perform or claim formal lifecycle transitions" (README). Results are bound to a work order ID, a state fingerprint and a claim ID; a result for an older state is rejected as stale. Task README writes take a lock file and check an expected `revision` ([C/tasks/task-readme-io.ts#L65](https://github.com/basilisk-labs/agentplane/blob/81fc89167d9d7655bd29664849af6fa2eb4bb054/packages/core/src/tasks/task-readme-io.ts#L65), [A/commands/shared/task-store/store.ts#L104-L108](https://github.com/basilisk-labs/agentplane/blob/81fc89167d9d7655bd29664849af6fa2eb4bb054/packages/agentplane/src/commands/shared/task-store/store.ts#L104-L108)). An execution lease goes stale after 11 minutes ([A/commands/shared/supervisor-execution-lease.ts#L8](https://github.com/basilisk-labs/agentplane/blob/81fc89167d9d7655bd29664849af6fa2eb4bb054/packages/agentplane/src/commands/shared/supervisor-execution-lease.ts#L8)). All locks are local to one checkout.
- **State location.** Committed: `.agentplane/tasks/<id>/README.md` (YAML front matter plus Markdown, with the kernel record and an event list), `quality/objects/sha256/*.json`, `supervision/`, `verification/`, `pr/`. Not committed: `<git-common-dir>/agentplane/kernel/exchanges/<task>/<order>/` (work orders, results, validations), supervisor journals and runner traces ([A/commands/task/kernel-exchange.ts#L34-L47](https://github.com/basilisk-labs/agentplane/blob/81fc89167d9d7655bd29664849af6fa2eb4bb054/packages/agentplane/src/commands/task/kernel-exchange.ts#L34-L47)). Gitignored: `.agentplane/tasks/*/handoff`, `tasks.json`, `cache.sqlite` (target `.gitignore`, line 36 in this repository).
- **Task kernel.** States and transitions ([C/tasks/task-kernel/kernel.ts#L38-L56](https://github.com/basilisk-labs/agentplane/blob/81fc89167d9d7655bd29664849af6fa2eb4bb054/packages/core/src/tasks/task-kernel/kernel.ts#L38-L56)); `reject_plan` and `amend_plan` events ([#L146-L160](https://github.com/basilisk-labs/agentplane/blob/81fc89167d9d7655bd29664849af6fa2eb4bb054/packages/core/src/tasks/task-kernel/kernel.ts#L146-L160)); COMPLETED has no outgoing transition.
- **Plan dataflow check.** Each `required_input` of a work item must be the `expected_output` of exactly one other item; cycles and missing dependencies fail ([C/tasks/task-kernel/invariants.ts#L223-L240](https://github.com/basilisk-labs/agentplane/blob/81fc89167d9d7655bd29664849af6fa2eb4bb054/packages/core/src/tasks/task-kernel/invariants.ts#L223-L240)).
- **Scope and rule detection.** A result that changed paths outside the work item scope is refused ([A/commands/task/kernel-semantic-result.ts#L58](https://github.com/basilisk-labs/agentplane/blob/81fc89167d9d7655bd29664849af6fa2eb4bb054/packages/agentplane/src/commands/task/kernel-semantic-result.ts#L58)). The authority carries `policy_digests`; a change gives `native_policy_changed` ([A/runner/usecases/kernel-authority.ts#L310-L334](https://github.com/basilisk-labs/agentplane/blob/81fc89167d9d7655bd29664849af6fa2eb4bb054/packages/agentplane/src/runner/usecases/kernel-authority.ts#L310-L334)). The pre-commit hook protects `AGENTS.md`, `CLAUDE.md`, `.agentplane/policy`, `.agentplane/agents` and others, not `.agentplane/WORKFLOW.md`; the override is the environment variable `AGENTPLANE_ALLOW_POLICY` ([A/shared/protected-paths.ts#L51-L100](https://github.com/basilisk-labs/agentplane/blob/81fc89167d9d7655bd29664849af6fa2eb4bb054/packages/agentplane/src/shared/protected-paths.ts#L51-L100)).
- **Human approval.** `--by USER` passes when the string matches `^USER(:...)?$` and the call comes from a "manual operator" route ([A/runner/usecases/kernel-authority.ts#L175-L180](https://github.com/basilisk-labs/agentplane/blob/81fc89167d9d7655bd29664849af6fa2eb4bb054/packages/agentplane/src/runner/usecases/kernel-authority.ts#L175-L180)). A second route checks an Ed25519 signed receipt from a configured trusted issuer ([A/adapters/authority/user-approval-receipt.ts#L181-L223](https://github.com/basilisk-labs/agentplane/blob/81fc89167d9d7655bd29664849af6fa2eb4bb054/packages/agentplane/src/adapters/authority/user-approval-receipt.ts#L181-L223)); the default list of issuers is empty.
- **Telemetry.** `token_usage` per task with `state` `observed`, `partial` or `unavailable` and an `unavailable_reason` ([A/commands/task/task-token-usage.ts#L36-L140](https://github.com/basilisk-labs/agentplane/blob/81fc89167d9d7655bd29664849af6fa2eb4bb054/packages/agentplane/src/commands/task/task-token-usage.ts#L36-L140)). Tokens come only from Codex `turn.completed` events ([A/runner/adapters/codex-result-transport.ts#L310](https://github.com/basilisk-labs/agentplane/blob/81fc89167d9d7655bd29664849af6fa2eb4bb054/packages/agentplane/src/runner/adapters/codex-result-transport.ts#L310)); an external agent (for example Claude Code) is `external_host_turn_not_task_attributable` ([A/commands/task/external-agent-supervisor.ts#L247](https://github.com/basilisk-labs/agentplane/blob/81fc89167d9d7655bd29664849af6fa2eb4bb054/packages/agentplane/src/commands/task/external-agent-supervisor.ts#L247)). No money cost field. Durations are in execution receipts and journals under the Git directory.
- **Checks.** Commands run as argv with no shell, through an executable allowlist: `bash, bun, cat, chmod, gh, git, node, npm, ps, sh, tar, unzip, zip` ([C/process/run-process.ts#L175-L224](https://github.com/basilisk-labs/agentplane/blob/81fc89167d9d7655bd29664849af6fa2eb4bb054/packages/core/src/process/run-process.ts#L175-L224); no change since tag v0.7.11). Kernel final validation uses `allow_empty: true` ([A/commands/task/kernel-final-validation.ts#L208](https://github.com/basilisk-labs/agentplane/blob/81fc89167d9d7655bd29664849af6fa2eb4bb054/packages/agentplane/src/commands/task/kernel-final-validation.ts#L208)).
- **Limits.** `max_rework_attempts` default 3 ([C/config/schema.impl.ts#L544](https://github.com/basilisk-labs/agentplane/blob/81fc89167d9d7655bd29664849af6fa2eb4bb054/packages/core/src/config/schema.impl.ts#L544)); supervisor episode budgets exist but default to off, `max_episodes` = `Number.MAX_SAFE_INTEGER` ([A/commands/shared/supervisor-execution-episode.ts#L40-L50](https://github.com/basilisk-labs/agentplane/blob/81fc89167d9d7655bd29664849af6fa2eb4bb054/packages/agentplane/src/commands/shared/supervisor-execution-episode.ts#L40-L50)).
- **Harnesses.** Managed runner adapters: `codex`, `custom`, `hermes` ([A/runner/adapters/index.ts#L8-L20](https://github.com/basilisk-labs/agentplane/blob/81fc89167d9d7655bd29664849af6fa2eb4bb054/packages/agentplane/src/runner/adapters/index.ts#L8-L20)). Codex starts with `--ignore-user-config --strict-config --disable hooks --ephemeral` ([A/runner/adapters/codex-preparation.ts#L180-L190](https://github.com/basilisk-labs/agentplane/blob/81fc89167d9d7655bd29664849af6fa2eb4bb054/packages/agentplane/src/runner/adapters/codex-preparation.ts#L180-L190)). Any other agent uses the `task advance` exchange loop by hand or by script.

## What we ran

npm has no 0.7.12-beta.1, so I installed `agentplane@0.7.11` (latest, tag `v0.7.11` = `65b1c24`, 128 commits before the pin) with `npm install --prefix eval-work/agentplane/npm`, `HOME=eval-work/agentplane/home`, `GH_TOKEN` and `GITHUB_TOKEN` unset. Toy target: the same Go module as for GNAP, with a local bare remote.

1. `agentplane init --tool claude --workflow direct --backend local --hooks true --require-plan-approval true --yes`. It made commit `chore: install agentplane 0.7.11`: 37 files, 1,804 lines (`CLAUDE.md`, `.agentplane/WORKFLOW.md`, 15 agent JSON files, 13 policy files, a shim `.agentplane/bin/agentplane`, `.gitignore`), and it installed `pre-commit`, `commit-msg`, `pre-push` and `post-merge` hooks in `.git/hooks`.
2. `task create "Add calc.Sub ..." --verify "make check" --route direct`. Task `202609281241-756JVV`. `task advance --agent-json` gave a PLANNER packet; the work order and the result path are under `.git/agentplane/kernel/exchanges/...`.
3. PLANNER with Claude Code, three tries: run 1 could not write into `.git` (Claude Code protects it), 0.168 USD, no result; run 2 wrote a result of the wrong shape, and AgentPlane refused it (`Canonical result requires its issued work_order_id`), 0.097 USD; run 3 with `--json-schema` (the result schema minus its top-level `allOf`, which the API refuses) gave a valid shape, 0.187 USD. AgentPlane then refused it twice: `scope_roots ... "calc/"` (trailing slash), then `invalid_input:implement-sub-function:Go 1.23+ environment ...` (the dataflow check). I fixed both by hand. Then:
   ```
   {"action":{"kind":"approval_required","reason":"kernel_plan_approval_required"},
    "operator_action":{"argv":["agentplane","task","plan","approve","...","--by","USER"]}}
   ```
4. Approval: `--by CODER` gave `Canonical authority rejected: explicit_manual_operator_required`. `--by USER` from my non-interactive shell (stdin `/dev/null`) passed and wrote `updated_by: "USER"`. Any process with a shell can do the same.
5. EXECUTOR with Claude Code (0.162 USD) added `Sub`. AgentPlane committed `(emoji) 756JVV task: apply canonical agent result`, then native validation ran the plan's checks and failed:
   ```
   "command":"go build ./calc","exit_code":null,
   "stderr_tail":"process command is not in the allowed executable set: go"
   "status":"FAILED", "reason":"Declared check could not run: go build ./calc"
   ```
   A check that could not run counted as a failure, and the work item went to rework (attempt 2). `task create --verify 'sh -c "go vet ./... && go test ./..."'` is refused ("inline shell evaluation is not allowed"). The way that worked: a repository script `scripts/check.sh` and `--verify "sh scripts/check.sh"`.
6. Hooks and the human: `git commit` of `scripts/check.sh` by me was refused ("Mutating staged paths require an active AgentPlane task"), and so was `git push` ("pre-push blocked"). `--no-verify` passed both. A fresh clone has no hooks.
7. Rule change by a working agent (task `202609281251-XPT2XM`, plan and result written by hand as the agent): I edited `calc/` and also `.agentplane/WORKFLOW.md` (`require_plan: false`, `require_verify: false`) and `.agentplane/policy/dod.code.md`. Results:
   - `task advance --result` refused: `Canonical implementation changed paths outside its WorkItem scope`.
   - `git commit` of `WORKFLOW.md` with `AGENTPLANE_TASK_ID` set: accepted by the hook.
   - `git commit` of `policy/dod.code.md`: `is protected by agentplane hooks (set AGENTPLANE_ALLOW_POLICY=1 to override)`; with `AGENTPLANE_ALLOW_POLICY=1`: accepted.
   - The task then stopped for good: `Canonical authority rejected: native_policy_changed`.
   - A new task created after the change still stopped at `kernel_plan_approval_required`; `.agentplane/workflows/last-known-good.md` kept `require_plan: true`, so the edit did not take effect there.
8. Two writers (clones `c1` and `c2` of the bare remote, task `202609281253-B3QN3B`): both ran `task plan approve --by USER`, both got `revision: 4` locally, both committed. `c1` pushed; `c2` was rejected, and `git pull --rebase` gave `CONFLICT (content): Merge conflict in .agentplane/tasks/202609281253-B3QN3B/README.md` on `updated_at` and `note`. No merge driver exists.
9. Task 3 to review: EXECUTOR by hand in scope; `sh scripts/check.sh` passed; the task moved to an EVALUATOR episode. Claude Code as EVALUATOR (0.147 USD) returned `status: blocked` and found a real defect in my hand-made plan: "Task objective and acceptance criteria are inconsistent (Neg vs Mul)". AgentPlane refused that result as `Canonical implementation result is stale` and reissued the same packet. I did not find why; I stopped there.
10. Records written: the toy task README grew to 767 lines after two episodes. `token_usage` was not written, because no task reached finish. In AgentPlane's own repository, 3,436 task records exist; `token_usage` has `state: "unavailable"` in 163 of them.

Processes: none stay after a command. `ps` at 12:59:49Z showed no `agentplane`, `node` or `claude` process from the work directory.

## In-Scope items S1-S12

| Item | Mark | Evidence |
| ---- | ---- | -------- |
| S1 Problem Statement Quality | no | No gap check of a problem statement. |
| S2 Reproducible Discipline Setup | partly | `init` sets up its own policy pack the same way each time (step 1). Not Armature, and no evidence per value: WORKFLOW.md has about 100 defaults with no source. |
| S3 Rule Protection | partly | Hooks protect policy paths, scope reconciliation refuses out-of-scope results, and `native_policy_changed` stops the running task (step 7). But WORKFLOW.md is not protected, the override is an environment variable that the agent can set, hooks can be skipped, and `--by USER` is a string. |
| S4 Stack-Dependent Gates | no | No layout, boundary or contract gates for a target stack. `go` and `make` are not in the executable allowlist (step 5). |
| S5 Role Handoffs | partly | Typed work orders and results, validated by zod, bound to a fingerprint, with a dataflow check (steps 3, 9). Not based on Armature; the evidence of each exchange stays in `.git/agentplane/`; `tasks/*/handoff` is gitignored. |
| S6 Autonomous Clarification | partly | A `needs_context` result gets a bounded local knowledge search by the machine, not an answer from the responsible role agent. |
| S7 Verification on Every Change | partly | The supervisor runs declared checks and counts a check that cannot run as a failure (step 5). Checks are per task, not stack gates, and an empty check list passes kernel final validation. |
| S8 Human-on-the-Loop | partly | Typed stops for plan approval, verification approval and human review (step 3). No escalation rule for business-forking decisions, and the human identity is self-declared (step 4). |
| S9 Stall Resolution | partly | `max_rework_attempts` 3 then `blocked_external`; runner idle and wall-clock timeouts. No independent diagnosis with a fresh context; episode budgets are off by default. |
| S10 Cost Visibility | partly | Per-task `token_usage` in the task README with an honest `unavailable` state. Tokens only from Codex; no money; durations are in files under `.git`, not in the committed record. |
| S11 Specification Synthesis | no | The PLANNER writes work items, not requirements traced to a problem statement. |
| S12 Harness-Agent Neutrality | partly | One policy generates `AGENTS.md` or `CLAUDE.md`; the exchange loop works with any agent (I used Claude Code). The managed runner and the telemetry support Codex only. A separate EVALUATOR role exists, but nothing requires a different harness. |

## Invariants I1-I9

| Invariant | Effect | Reason |
| --------- | ------ | ------ |
| I1 Git is the system of record | conflicts | Work orders, results, validations, journals and runner traces are in `.git/agentplane/`, which is never pushed. Committed READMEs point to them with `../../../.git/agentplane/...` paths that break in a clone. |
| I2 Independent repository | conflicts | In the working clone, the installed hooks refused a human commit and push until `--no-verify` (step 6). A fresh clone has no hooks, and the committed files are readable, so a human can continue. The committed `CLAUDE.md` and shim tell agents to use AgentPlane, which O-76 does not allow in a target. |
| I3 Agents cannot change rules or gates | conflicts | An agent can commit WORKFLOW.md freely and a policy file with one environment variable (step 7). Detection is strong: the running task stops on `native_policy_changed`. |
| I4 No value without evidence | conflicts | `init` writes defaults (timeouts, 3 rework attempts, retry counts) with no evidence. It does support I4 for telemetry: missing tokens are `unavailable`, never zero. |
| I5 Inactive check is not a pass | supports | A check that could not run was `FAILED` (step 5); a required `not_run` check makes a run "unverified". Exception: `allow_empty: true` in kernel final validation passes a plan with no checks. |
| I6 Deterministic over LLM | supports | Transitions, scope, digests and the plan dataflow are code. The LLM only plans, implements and reviews. |
| I7 Domain changes content, not rules | neutral | No notion of domain rules. |
| I8 Pinned baseline | neutral | WORKFLOW.md records `framework.cli.expected_version: "0.7.11"`, a pin of its own version. No Armature. |
| I9 Replaceable harness | supports | The exchange protocol is harness-neutral (step 3 used Claude Code). Gap: telemetry and the managed runner are Codex only. |

## Deep-check findings it answers

- **Fable-M12 / Sol-3 (one writer has six writers).** Answered as a pattern: the supervisor CLI is the only writer of formal state; agents return typed results into an exchange directory, and the CLI admits or refuses them and commits. This is the finding's own fix ("each other producer gives its output to `layup run`, which commits it"). Limit: its locks are local; two clones still conflict in Git (step 8).
- **Author-8 / Fable-M14 / Sol-21 (handoff schema).** Answered as a pattern: a result schema per role and phase, a `status` of `completed`, `blocked`, `needs_context` or `failed`, a binding to work order and state fingerprint, stale rejection, the dataflow check, and "claimed" checks kept apart from "observed" checks. Validity is computed from the artifact and state, not from the agent's claim, as Sol-21 asks.
- **Sol-11 (unknown spend).** Partly: missing tokens are `unavailable` with a reason, never a guessed zero. It has no money bound.
- **Fable-M4 / Author-1 / Sol-22 (disagreement or repeat is never a stall).** Partly: evaluator rework rounds are counted per work item (`attempt`), and after `max_rework_attempts` the task is blocked. A new commit does not reset the counter.
- **Fable-M1 / Fable-M2 / Sol-1 (rule protection is detection only).** Partly: policy digests bound to the approved authority make the running task fail closed on a rule change. It is still not prevention (step 7).
- **Sol-5 (commit author is not proof).** Partly: the signed approval receipt route is a real authentication of a human decision, independent of Git metadata. It is off by default.
- **Fable-M23 (host user-level setup is a second rule source).** Answered for Codex: every run starts with `--ignore-user-config --strict-config --disable hooks --ephemeral`.
- **Fable-M15 / Sol-14 (measures with no record).** Partly: each transition is an event with a mutation ID, a revision and a time in the task README. No intervention, reversal or audit-sample record.
- **Author-9 / Sol-13 (reversal record), Sol-4 (dead-man job), Sol-12 (incomplete telemetry cannot pass), Author-6 (order of Intake and records), Fable-M6 (stall with no diagnosis).** None.

## Parts and their verdicts

| Part | Verdict | How LAYUP would take it | Reason |
| ---- | ------- | ----------------------- | ------ |
| Supervisor as the only writer; agents return typed results to an exchange | borrow the pattern | `layup run` issues a work file per session and admits a typed result file; only `layup run` commits records | Answers Fable-M12 and Sol-3. |
| Result binding: work order ID, state fingerprint, stale rejection | borrow the pattern | Hash of the task record in each work file; refuse a result for an older hash | Deterministic (I6); stops lost updates in one checkout. |
| Result statuses `completed / blocked / needs_context / failed` with required `blocker` or `knowledge_request` | borrow the pattern | LAYUP handoff and question records | A typed question from a session is the start of S6. |
| Plan dataflow check (inputs = outputs of exactly one dependency) | borrow the pattern | A check in LAYUP's plan validation, in Go | Cheap and deterministic. |
| `token_usage` with `observed / partial / unavailable` and `unavailable_reason` | borrow the pattern | LAYUP telemetry row | Honest missing data (I4, I5). LAYUP must add latency, wall clock, money and harnesses other than Codex. |
| Policy digests bound to the approved authority (`native_policy_changed`) | borrow the pattern | Record the hash of the rule files at plan approval; stop a task when it changes | Complements prevention (Fable-M1, M2). |
| Signed approval receipts from trusted issuers | borrow the pattern | Compare with the GitHub App marker of O-77 | Proves a human decision without trusting commit metadata. |
| Codex argv with no user config | borrow the pattern | LAYUP harness adapters start sessions with no user-level setup | Answers Fable-M23. |
| AgentPlane as a program next to LAYUP | reject | Not used | Installs hooks, a shim and `CLAUDE.md` into the target (O-76), keeps state in `.git/agentplane` (I1), cannot run `go` or `make` checks, reads tokens from Codex only, Node 24, pre-1.0 with fast churn. |
| Task README (YAML front matter plus event list) as LAYUP's records format | reject | Not used | 767 lines for one toy task after two episodes; one mutable file per task gives text conflicts with two writers (step 8); evidence links point into `.git`. |

How it deals with many writers to one Git branch: inside one checkout, by lock files, revision checks, leases and a Git mutation mutex. Across clones, only by Git itself and `--force-with-lease` when it publishes a branch; in `branch_pr` mode each task has its own branch and worktree. Two clones that change the same task README conflict (step 8).

Could its file format be LAYUP's records format? Not as it is (above). The typed JSON result objects and the `token_usage` block are good models for LAYUP's handoff and telemetry records.

## Where the searchers were wrong or incomplete

- Searcher B: "All state stays in `.agentplane/` inside the repo; no hosted runtime." Incomplete. The exchanges, journals and traces are in `.git/agentplane/`, which Git never pushes; a `cloud` backend exists.
- Searcher B: language and last push "not stated", S10 `no`. The code has per-task `token_usage` (partly), and the project is active (last push 2026-09-28).
- Searcher B: S3 `partly` because of the workflow words. The real mechanisms are the hooks, scope reconciliation and policy digests; step 7 shows their gaps.
- Searcher B: S5 "a machine schema for role artifacts is not stated". The code has zod schemas for every result, published as JSON Schemas in `schemas/`.
- Searcher B: "wraps Claude Code, Codex, Cursor, and Aider". The managed runner supports Codex, a custom command and Hermes only; Claude Code works only through the exchange loop.
- Searcher A lists it as "keep: Git-native plan/approve/verify" with no rating. The approval is a self-declared `--by USER` unless signed receipts are set up.

## Limits of this evaluation

- I ran 0.7.11 from npm; the code references are at the pinned SHA (0.7.12-beta.1). The executable allowlist is the same in both.
- I did not run `branch_pr` mode, the managed runner with Codex, hosted PR checks, or the `cloud` backend.
- No task reached `DONE`, so I did not see `token_usage`, the ACR or close-out records on the toy. The EVALUATOR result was refused as stale for a reason I did not find.
- Several plan and result files were written or fixed by hand (steps 3, 7, 9); haiku did not produce the strict shapes without help.
- A code-mapping subagent read parts of the source; I checked the cited lines above myself.
