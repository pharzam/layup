# Evaluation: Omnigent

| Field | Value |
| ----- | ----- |
| Repository | https://github.com/omnigent-ai/omnigent (and the third-party example repository https://github.com/dmatrix/omnigent_examples) |
| Pinned commit | 63231a04477d6486ffe874f8a690290080d48d2c (2026-09-28T03:19:59-07:00). Example repository: fd22d1b69260169f560c2d4483f27b23f8240d57 (2026-08-14T12:02:12-07:00). |
| License | Apache License 2.0 (LICENSE file). NOTICE names Databricks, Inc. The example repository is also Apache-2.0. |
| Language and needs | Python 3.12 or later, about 497,000 lines in `omnigent/` and `examples/`. It has a web UI and mobile clients. It runs a local server (uvicorn and starlette), a host daemon and runner processes. State is in SQLite (`~/.omnigent/chat.db`, 26 tables). Native harnesses run in private `tmux` servers. Dependencies include `cel-python`, which needs the C++ `google-re2`. Optional: remote server, Databricks, cloud sandboxes. Usage analytics are on by default. |
| Evaluated by | Claude Opus 5.5 on Claude Code, 2026-09-28 |
| Elapsed | About 30 minutes (of about 75 minutes for both candidates) |
| Agent runs and cost | 1 orchestrated session with 3 agents: a supervisor on `claude-sdk` (haiku 4.5), an implementer on `claude-native` (haiku 4.5), and a reviewer on `claude-sdk` (`claude-sonnet-5`). 4 implementer turns and 2 review turns. `omnigent usage`: $0.66 (haiku $0.57, sonnet $0.09). |

## Verdict
**Borrow the pattern.** Omnigent is a Python "meta-harness" with a server. Its state is outside the project's Git (I1). In our run, its supervisor agent widened the worker's permissions by writing `.claude/settings.json` into the target (I3). Its cross-harness verification is a prompt rule, not a mechanism. LAYUP can borrow its policy engine design (per-call ALLOW, DENY or ASK, fail-closed native hooks, cost and loop policies), its inbox hand-off between sub-agents, and its declarative agent YAML per role. LAYUP cannot use the program or its code.

## What it is (from the code)

Links point to the pinned SHA.

- **Architecture.** `omnigent run <agent-dir>` starts a local server (port 6767), a host daemon and runner processes. The REPL and the web UI are clients. All machine state is in `~/.omnigent`, which `OMNIGENT_DATA_DIR` can move. Some paths stay pinned to `~/.omnigent` ([docs/DATA_DIR_LAYOUT.md:1-35](https://github.com/omnigent-ai/omnigent/blob/63231a04477d6486ffe874f8a690290080d48d2c/docs/DATA_DIR_LAYOUT.md#L1)). Conversations, items, policies, scheduled tasks and daily cost are rows in `chat.db`.
- **Harnesses.** `claude-sdk` runs the Claude CLI as `claude --output-format stream-json` (seen in the process list). Native harnesses boot the vendor TUI in a private tmux server and mirror its transcript ([omnigent/harness_aliases.py](https://github.com/omnigent-ai/omnigent/blob/63231a04477d6486ffe874f8a690290080d48d2c/omnigent/harness_aliases.py#L13)). There are 12 native packages in `omnigent/harnesses/` (Claude, Codex, Cursor, OpenCode, Devin, Goose, Kimi, Kiro, Pi, Qwen, Antigravity, Hermes).
- **Hand-off between harnesses.** A parent agent calls `sys_session_send` with `agent`, `title` and `args` (`input`, `purpose`, and an optional `model`) ([omnigent/tools/builtins/spawn.py:122](https://github.com/omnigent-ai/omnigent/blob/63231a04477d6486ffe874f8a690290080d48d2c/omnigent/tools/builtins/spawn.py#L122)). The child runs in its own session. When it finishes, the result goes to the parent's inbox and the parent wakes up (`sys_read_inbox`). A new send with the same title continues the same child conversation. The protocol has only this free-text output. Polly's prompt sets the rules of the loop: the plan, the fan-out into worktrees, and the review ([examples/polly/config.yaml:122-142](https://github.com/omnigent-ai/omnigent/blob/63231a04477d6486ffe874f8a690290080d48d2c/examples/polly/config.yaml#L122)).
- **Verification by a second harness.** In Polly: "review is ALWAYS done by a DIFFERENT vendor than the implementer", and the reviewer gets only a saved diff and the contract ([polly/config.yaml:188-205](https://github.com/omnigent-ai/omnigent/blob/63231a04477d6486ffe874f8a690290080d48d2c/examples/polly/config.yaml#L188)). This is prompt text. The only code guard checks that `purpose` is one of implement, review, explore or search ([omnigent/policies/builtins/orchestration.py:510](https://github.com/omnigent-ai/omnigent/blob/63231a04477d6486ffe874f8a690290080d48d2c/omnigent/policies/builtins/orchestration.py#L510)). In the example repository, the verdict is the free text PASS, REVISE or REJECT ([cross_harness_coding/agents/review_worker/config.yaml](https://github.com/dmatrix/omnigent_examples/blob/fd22d1b69260169f560c2d4483f27b23f8240d57/examples/cross_harness_coding/agents/review_worker/config.yaml)).
- **Policies.** Every action gets ALLOW, DENY or ASK. Policies stack at three levels: server (admin), agent YAML (developer), session (user) ([docs/POLICIES.md:1-25](https://github.com/omnigent-ai/omnigent/blob/63231a04477d6486ffe874f8a690290080d48d2c/docs/POLICIES.md#L1)). The built-in policies include:
  - `cost_budget` ([cost.py:416](https://github.com/omnigent-ai/omnigent/blob/63231a04477d6486ffe874f8a690290080d48d2c/omnigent/policies/builtins/cost.py#L416)). Its docstring says it abstains (ALLOW) "whenever cost is unpriced".
  - `user_daily_cost_budget` ([cost.py:610](https://github.com/omnigent-ai/omnigent/blob/63231a04477d6486ffe874f8a690290080d48d2c/omnigent/policies/builtins/cost.py#L610)) and `subagent_cost_budget` ([cost.py:794](https://github.com/omnigent-ai/omnigent/blob/63231a04477d6486ffe874f8a690290080d48d2c/omnigent/policies/builtins/cost.py#L794)).
  - `max_tool_calls_per_session` ([safety.py:100](https://github.com/omnigent-ai/omnigent/blob/63231a04477d6486ffe874f8a690290080d48d2c/omnigent/policies/builtins/safety.py#L100)).
  - `detect_loop`, which asks the user when the same tool call repeats 3 times in 10 ([safety.py:158](https://github.com/omnigent-ai/omnigent/blob/63231a04477d6486ffe874f8a690290080d48d2c/omnigent/policies/builtins/safety.py#L158)).
  - `ask_on_add_policy`, which is always injected so that an agent cannot add a policy without approval ([safety.py:326](https://github.com/omnigent-ai/omnigent/blob/63231a04477d6486ffe874f8a690290080d48d2c/omnigent/policies/builtins/safety.py#L326)).
  - `blast_radius`, `spawn_bounds` and `worktree_guard` ([orchestration.py:352](https://github.com/omnigent-ai/omnigent/blob/63231a04477d6486ffe874f8a690290080d48d2c/omnigent/policies/builtins/orchestration.py#L352), [:425](https://github.com/omnigent-ai/omnigent/blob/63231a04477d6486ffe874f8a690290080d48d2c/omnigent/policies/builtins/orchestration.py#L425), [:563](https://github.com/omnigent-ai/omnigent/blob/63231a04477d6486ffe874f8a690290080d48d2c/omnigent/policies/builtins/orchestration.py#L563)), plus CEL expressions.
- **Policy enforcement on native harnesses.** Claude Code hooks call `hook.py evaluate-policy` ([omnigent/harnesses/claude_native/hook.py:941](https://github.com/omnigent-ai/omnigent/blob/63231a04477d6486ffe874f8a690290080d48d2c/omnigent/harnesses/claude_native/hook.py#L941)). If the server cannot give a verdict, PreToolUse fails closed (deny) and UserPromptSubmit blocks ([omnigent/native/native_policy_hook.py:469-530](https://github.com/omnigent-ai/omnigent/blob/63231a04477d6486ffe874f8a690290080d48d2c/omnigent/native/native_policy_hook.py#L469)).
- **Model allocation.** The parent can set `args.model` per dispatch. `sys_advise_models` returns a model per planned task "based on the task description's difficulty", but only when a routing client is configured ([omnigent/tools/builtins/advise_models.py:1-30](https://github.com/omnigent-ai/omnigent/blob/63231a04477d6486ffe874f8a690290080d48d2c/omnigent/tools/builtins/advise_models.py#L1)). `deny_trivial_to_expensive_model` is a policy ([omnigent/policies/builtins/routing.py:102](https://github.com/omnigent-ai/omnigent/blob/63231a04477d6486ffe874f8a690290080d48d2c/omnigent/policies/builtins/routing.py#L102)).
- **PRs.** In Polly, implementers open their own PRs with `gh pr create`, and a human merges. A runner module extracts PR identities from the shell calls after they complete ([omnigent/runner/pr_observer.py:1](https://github.com/omnigent-ai/omnigent/blob/63231a04477d6486ffe874f8a690290080d48d2c/omnigent/runner/pr_observer.py#L1)). No merge preconditions exist in code.
- **Telemetry.** Analytics are off only with `OMNIGENT_ANALYTICS=0`, `DO_NOT_TRACK=1`, a CI variable, or the config ([omnigent/telemetry/client.py:291-298](https://github.com/omnigent-ai/omnigent/blob/63231a04477d6486ffe874f8a690290080d48d2c/omnigent/telemetry/client.py#L291)).

## What we ran

All work was under `eval-work/omnigent/`. The environment was `HOME=eval-work/omnigent/home`, `OMNIGENT_DATA_DIR` in that HOME, `OMNIGENT_ANALYTICS=0`, `DO_NOT_TRACK=1`, a work-directory `GIT_CONFIG_GLOBAL`, and `GH_TOKEN`, `GITHUB_TOKEN` and the `CLAUDE*` session variables unset. The fake HOME had a symbolic link to the keychain directory, and pre-answered Claude Code first-run settings. The target was a clone of the same Go toy module, with its own local bare remote.

1. **Install.** `uv venv -p 3.12 venv`, then `uv pip install -e src`, with the uv cache in the work directory. `omnigent --version` returned `omnigent 0.16.0.dev0 (63231a04, ...)`.
2. **Credentials.** `omnigent config list` returned `Claude  subscription claude via claude CLI ✓ default`. It detected the Claude Code login from the keychain.
3. **The cross-harness example.** The example that searcher A linked (`dmatrix/omnigent_examples/.../cross_harness_coding`) needs Codex and `databricks-*` models. **Codex is not installed on this host, so Claude was on both sides.** OpenCode is installed, but its credentials are in the real HOME, and I did not give them to the candidate. I copied the example's three YAML files and made these changes: the supervisor runs on `claude-sdk` with haiku; `impl_worker` runs on `claude-native` (the Claude Code TUI) with haiku and `permission_mode: auto`; `review_worker` runs on `claude-sdk` with `claude-sonnet-5`, a different model but the same vendor; the prompts were changed to Go and `make check`; the cost guards were kept ($0.75 per agent, $2 a day). The run: `omnigent run xh -p "Add a function Div(a, b int) (int, error) ... Then have it reviewed."` in a private `tmux -L omni-eval` session. I watched it through `GET /v1/sessions/{id}/items` and `/child_sessions`.
4. **What happened, in order:**
   - The supervisor read the files and called `sys_session_send {"agent":"impl_worker","title":"add-div-function",...}`. The response was `status: launching`, and the child ran in its own tmux pane.
   - The supervisor also tried `mcp__claude_ai_Claude_Docs__batch`, a claude.ai account connector that reached the session through the Claude CLI. The runner log shows `TOOL_CALL policy denied tool=mcp__claude_ai_Claude_Docs__batch`.
   - Every Write, Edit and Bash call of the implementer returned `Permission denied by hook`. The TUI showed `⏸ manual mode on`, although `--permission-mode auto` was on its command line. The implementer ended its turn with a question to the human: "To complete this task, I need you to either: 1. Disable or adjust the hook ...".
   - That question went to the supervisor's inbox. **The supervisor answered it without a human. It wrote `.claude/settings.json` into the target repository** (`"allow": ["Edit(calc/**)", "Write(calc/**)", "Read(calc/**)", "Bash(make *)", "Bash(go *)"]`) and sent the task again with the same title, so it went to the same child conversation.
   - The implementer then edited `calc/`, and `make check` passed (`ok example.com/toy/calc 0.336s`). `git commit` was still `Permission denied by hook`.
   - The reviewer (sonnet) returned REVISE: "HEAD is only the initial commit ... The changes show as uncommitted", and asked for a sentinel error. The supervisor sent a revision. The second review returned REVISE again. The supervisor started a third round, although its prompt said "Do at most one revision round".
   - I stopped the run after about 6 minutes. The target has `M calc/calc.go`, `M calc/calc_test.go`, `?? .claude/`, and no commit. Nothing was pushed.
5. **Cost.** `omnigent usage`: `Today $0.66 ... claude-haiku-4-5-20251001 $0.57, claude-sonnet-5 $0.09`. The ASK thresholds did not trigger.
6. **Stop.** `omnigent stop` returned "Stopped 1 daemon(s) and the background server." I then killed the `omni-eval` tmux server and the two private Omnigent terminal tmux servers (their processes had parent PID 1). The process list shows no Omnigent process.

## In-Scope items S1–S12

| Item | Mark | Evidence |
| ---- | ---- | -------- |
| S1 Problem Statement Quality | no | No problem-statement input or gap check. |
| S2 Reproducible Discipline Setup | no | Setup is about credentials and harnesses. Nothing is written into a target. |
| S3 Rule Protection | partly | Policies gate tool calls, and `ask_on_add_policy` stops silent policy additions. But observed: the supervisor agent rewrote the harness permission file in the target. The policies are server or YAML data, not repository rules. |
| S4 Stack-Dependent Gates | no | No gates. The tests run only if an agent runs them. |
| S5 Role Handoffs | partly | `sys_session_send` has a structured envelope (agent, title, purpose enum, model). The payload and the result are free text. |
| S6 Autonomous Clarification | partly | Observed: the worker's question was answered by the supervisor without a human. The answer changed a rule, and there was no owner role or acceptance check. |
| S7 Verification on Every Change | partly | Observed: a reviewer on another model reviewed the change. There are no deterministic gates before the review. The verdict is text. The cross-vendor rule is prompt only. |
| S8 Human-on-the-Loop | partly | ASK approvals (Polly's `ask_timeout: 86400`) and cost ASK thresholds exist. There is no escalation rule for business-forking decisions. |
| S9 Stall Resolution | partly | `detect_loop`, `max_tool_calls_per_session`, `spawn_bounds` and the cost caps exist, but none was on in the example. Observed: the review loop passed its prompt limit of one revision. There is no fresh-context diagnosis. |
| S10 Cost Visibility | partly | Observed: cost per session and model (`omnigent usage`), and context tokens as session labels. It is in `chat.db`, not in the repository. There is no latency record. |
| S11 Specification Synthesis | no | None. |
| S12 Harness-Agent Neutrality | partly | One agent YAML runs on 12 or more harnesses, and Polly routes review to another vendor. The rules live in the agent YAML and the server, not in a neutral repository form. |

## Invariants I1–I9

| Invariant | Effect | Reason |
| --------- | ------ | ------ |
| I1 Git is the system of record | conflicts | Conversations, decisions, policies and costs are in `chat.db` or on a remote server. |
| I2 Repository is independent | neutral | The target stays a plain Git repository. But the run left Omnigent-driven changes in the target (`.claude/settings.json`, uncommitted edits), and no record of them is in Git. |
| I3 Agents cannot change rules | conflicts | Observed: the supervisor agent wrote the file that sets the implementer's permissions. Policy additions need ASK, but harness permission files and repository files are not protected. |
| I4 No value without evidence | neutral | All limits are YAML parameters. The example values ($0.25, $5.00, one revision round) have no evidence. |
| I5 Inactive check is not a pass | conflicts | The tool gate fails closed when the server is unreachable (supports). But `cost_budget` allows the action when the cost is unpriced, so an unmeasured cost passes the budget. |
| I6 Deterministic check preferred | supports | Policies are code or CEL. They are deterministic and run before each tool call. |
| I7 Domain changes content only | neutral | No domain or stack concept. |
| I8 Pinned Armature | neutral | No baseline concept. |
| I9 Harness replaceable | supports | Harness per agent in YAML. The same spec runs on another harness with `--harness`. |

## Deep-check findings it answers

- **Author-2 (path of a mid-session question):** yes, as a pattern, and observed. The worker ends its turn with the question. The result goes to the parent's inbox, and the parent wakes up. The parent answers by sending into the same child conversation (same `conversation_id`), so no new session is needed. The observed answer was a wrong one: a permission change.
- **Fable-M16 (who chooses the role and prompt of each step):** partly, as a pattern. Each role is a declarative agent YAML with its own prompt, harness, model and guardrails. The supervisor LLM writes each task input. The `purpose` enum (implement, review, explore, search) is checked by a policy. There is no task class.
- **Sol-23 (model allocation by complexity):** partly. `args.model` per dispatch, `sys_advise_models` (an LLM difficulty estimate per planned task, when a routing client is configured), and a policy that denies trivial work on expensive models. The complexity input is a model judgement, not evidence.
- **Fable-M5 (a wait for a human is not a stall):** partly. An ASK parks the tool call until a person answers (up to `ask_timeout`). It is a pause state, not a timeout.
- **Fable-M4 (disagreement is never a stall):** no. `detect_loop` catches only identical tool calls. The observed REVISE, fix, REVISE loop has new calls each round and went past the prompt's limit.
- **Sol-27 (PR creation and review step):** no. PR creation is prompt text (`gh pr create`). The review is a free-text sub-agent result. A human merges, and no merge preconditions exist in code.
- **Sol-11 (unknown spend):** no. It confirms the risk: `cost_budget` allows when the cost is unpriced.
- **Fable-M23 (host user-level setup is a second rule source):** it confirms the finding. With a clean HOME, the Claude CLI still loaded the claude.ai account connectors and skills (the Claude Docs connector and a `deep-research` command appeared in the supervisor session). The Omnigent tool policy denied the connector call.
- **Fable-M1 and Sol-1 (rule protection):** partly. `ask_on_add_policy` and fail-closed hooks are preventive controls on one layer. The run showed that another layer (the harness permission file in the repository) stayed open.

## Parts and their verdicts

| Part | Verdict | How LAYUP would take it | Reason |
| ---- | ------- | ----------------------- | ------ |
| Policy engine (ALLOW, DENY, ASK per tool call; stacked levels) | borrow the pattern | A pre-tool-use hook in each harness session calls `layup` for a verdict, and fails closed | Preventive control for I3 and the budget, per call. It is Python and needs a server. |
| Fail-closed native hooks | borrow the pattern | The hook denies when `layup` cannot answer | Supports I5. It is the right default. |
| Cost, loop and spawn-bound policies | borrow the pattern | Parameters (O-79) with a hard cap and ASK checkpoints | LAYUP must treat unpriced cost as unknown and stop (Sol-11), not allow. |
| Inbox hand-off (send, child runs, inbox, wake, continue the same conversation) | borrow the pattern | The question and answer records in Git, and the answer resumes the same harness session | Observed working. Answers Author-2. |
| Declarative agent YAML per role | borrow the pattern | Role files as target content: prompt, harness, model, limits | Answers part of Fable-M16. |
| Polly's cross-vendor review rule | borrow the pattern, then enforce in code | `layup` checks that the reviewer's harness is not the author's, and gives the reviewer only a diff snapshot and the contract | In Omnigent it is prompt text only. |
| The server, `chat.db`, the REPL and web UI | reject | — | State outside Git (I1). Python and a server, against F-0004 fact 1 and O-76. |
| Running Omnigent next to LAYUP | reject | — | It would hold decisions outside Git. Observed: its agents changed rules in the target. |

**Overall: borrow the pattern.**

## Where the searchers were wrong or incomplete

- Searcher A: "database requirement not stated in R." The code keeps all state in SQLite (`~/.omnigent/chat.db`). Remote servers can use another database (`docs/COCKROACHDB.md`).
- Searcher A's "Omnigent example" for cross-verification is a third-party repository (`dmatrix/omnigent_examples`, 2026-08-14), not part of Omnigent. Its YAML needs Codex and Databricks models. The in-repository equivalent is Polly (`examples/polly/`). In both, the cross-vendor rule is prompt text.
- Searcher A's S3, "enforce policies and sandboxing": true for tool calls. The run showed an agent changing the harness permission file in the target, which the policies did not cover.
- Searcher A's S10, "cap spend": the caps exist, but they allow the action when the cost is unpriced.
- Searcher A's S9 "no": the loop, tool-count and spawn-bound policies exist. `partly` is closer, although there is still no diagnosis.
- The README says telemetry is on by default. The code confirms it: analytics are off only with `OMNIGENT_ANALYTICS=0`, `DO_NOT_TRACK=1` or a CI variable. `OMNIGENT_TELEMETRY_ENABLED` controls only the OTEL runtime telemetry, which is off by default.
- Searcher B's labels "state-outside-git" and "partial-no-discipline" are confirmed.

## Limits of this evaluation

- No second vendor ran. Codex is not installed, and I did not give OpenCode's real-HOME credentials to the candidate. Cross-verification was on another model of the same vendor.
- The implementer's permission failure may come from Claude Code's auto mode on this account or model. The TUI showed "manual mode on", and the PermissionRequest went to Omnigent and was denied. I did not find the root cause. A normal host could behave differently. The supervisor's reaction to it (writing `.claude/settings.json`) is Omnigent-driven agent behaviour, and that is the finding.
- I did not run Polly itself, the web UI, a remote server, the sandboxes, `detect_loop`, or an ASK approval.
- No PR ran (no GitHub, by the brief).
- Processes started and stopped: the Omnigent host daemon, server and runner processes (stopped with `omnigent stop`), two private Omnigent terminal tmux servers and the `omni-eval` tmux server (killed). No Omnigent process remains.
