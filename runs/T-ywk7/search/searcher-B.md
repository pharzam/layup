# T-ywk7 — the search of M2b, searcher B, word for word

Searcher: Claude Opus 5.5 (`claude-opus-5-5`, effort `xhigh`) on Claude Code 2.1.295 (`claude -p`, `--permission-mode dontAsk`, read-only tools, web search and fetch, writes only under `out/`; `--setting-sources project,local`, `--strict-mcp-config`, `--no-session-persistence`, no `gh` login), a fresh session, by the Operator's O-178 (b): the author's model, blind to the author's context. 2026-10-09 06:15:37Z to 06:37:27Z (7 min 28 s of the session's own time); USD 19.61 by its result event. Brief: [`brief.md`](brief.md), with the five lists of [`sources.md`](sources.md) and [`gh-search.tsv`](gh-search.tsv); nothing of the LAYUP architecture, the specification or the survey of `T-w73g`. It read the lists through five subagents of the same model. **A model of the "not used" list ran inside it:** its result event's `modelUsage` names `claude-haiku-5-5` (4,759,374 input and 478,125 output tokens, USD 1.25) besides `claude-opus-5-5`; no message of the session or of its subagents names Haiku, and a test run showed that Claude Code's `WebFetch` reads a page with that model unless `ANTHROPIC_DEFAULT_HAIKU_MODEL` names another. One shell command was denied (an append to this file; it used the Edit tool instead). The output follows word for word inside a fenced block; its paths are relative to the searcher's directory.

~~~~text
# Public-solution search of M2b (task T-ywk7), searcher B

Date: 2026-10-09. Sources: my own web queries, the five lists in `sources/lists/` (each read in
full: 324, 884, 770, 494 and 801 lines), and all 200 rows of `sources/gh-search.tsv`. I ran no code
of a candidate and installed nothing. Each kept candidate was checked on its own repository or
documentation page. Letters (a) to (g) name the parts of the problem:

- (a) a session that the controller starts, isolates, limits in time and money, stops, and reads the result of;
- (b) the first check of a registered harness agent (its version and what it can do);
- (c) the choice of a harness agent and a model for each piece of work;
- (d) the record of tokens, latency, wall-clock duration and money (F-0001#38, F-0003#50, #60);
- (e) the rules in a form that no one harness agent owns (F-0003#26, #27, #52, #66);
- (f) no change by an agent to the paths of the rules and the gates (F-0001#3, F-0003#43, #64);
- (g) a machine-validatable form for role handoffs (F-0003#45, #59).

## Queries

| # | Query (web) | Words it comes from | What it found |
|---|---|---|---|
| 1 | orchestrate multiple coding agent CLIs headless Claude Code Codex Gemini CLI controller worktree session timeout budget | "a session ... that the controller starts, isolates, limits in time and money"; F-0003#52 | quintet-headless-worktrees skill (one-shot `codex exec`/`gemini -p` per worktree); Bernstein; a Codex issue on headless use; nothing on time/money caps |
| 2 | "deterministic orchestrator" CLI coding agents "no model in the coordination loop" git worktrees | F-0001#6 "deterministic check is preferred" | Bernstein (deterministic, 40+ CLI agents, HMAC audit chain); ORCH; OpenTree |
| 3 | kill coding agent CLI session process group timeout wall-clock limit headless orchestrator "max-turns" codex claude gemini | "limits in time", "stops" | process-group kill pattern (Setsid, SIGTERM then SIGKILL); reports that `--max-turns` is not a time limit; no session turn cap stated for `codex exec` |
| 4 | claude -p --output-format json total_cost_usd duration_ms usage --max-budget-usd | F-0003#50; "limits in ... money" | result fields `total_cost_usd`, `duration_ms`, `usage`; subtype `error_max_budget_usd`; the cap uses an estimate, also on a subscription |
| 5 | codex exec --json events turn.completed usage input_tokens output schema --output-schema | F-0003#50; "reads the result of" | `turn.completed` carries `usage.input_tokens`, `cached_input_tokens`, `output_tokens`; `turn.failed` |
| 6 | gemini cli headless --output-format json stats tokens | F-0003#50 | one JSON object `response`, `stats` (per-model tokens, API latency), `error` |
| 7 | opencode run --format json non-interactive tokens cost output events | F-0003#50 (a fourth harness) | JSONL `step_finish` with `tokens` and `cost` per step |
| 8 | Agent Client Protocol ACP coding agent editor JSON-RPC specification | "first check ... its version, and what it can do"; "starts ... stops" | ACP: JSON-RPC 2.0 over stdio, `initialize` negotiates `protocolVersion`; schema per version |
| 9 | ACP registry agents list claude-code-acp codex-acp gemini --experimental-acp capabilities initialize | M2b "a probe session on each registered harness" | ACP adapters: `claude-agent-acp`, `codex-acp`, `gemini --acp`, `opencode acp`, `copilot --acp`; flags differ by source |
| 10 | probe coding agent CLI capabilities version detection before dispatch orchestrator "doctor" harness health check multi-agent | (b); M2b probe session | agentabi (capabilities per agent CLI), CodingAgentRunner (.NET), version checkers; no tool that runs a live probe task |
| 11 | Harbor framework terminal-bench installed agents claude-code codex gemini-cli opencode adapters token cost trajectory | M2b probe; F-0003#66 "≥ 2 harness agents" | Harbor installed-agent adapters; ATIF trajectories with token and cost |
| 12 | Agent Trajectory Interchange Format ATIF schema | F-0001#38 Telemetry | ATIF v1.7 (NVIDIA NeMo docs use it); fields session_id, agent, steps |
| 13 | OpenTelemetry GenAI semantic conventions gen_ai.usage.input_tokens agent span | F-0001#38 | `gen_ai.usage.input_tokens`/`output_tokens`, `invoke_agent` span; status disputed (Development vs stable) |
| 14 | Claude Code OpenTelemetry monitoring claude_code.token.usage claude_code.cost.usage OTEL metrics events | F-0003#60 Telemetry Completeness | `claude_code.token.usage`, `claude_code.cost.usage`, `api_request` event with `cost_usd`, `duration_ms` |
| 15 | LLM model price list JSON tokens cost per model open source data | F-0003#22 "no number for" the budget; "money" | OpenLIT pricing JSON, llm-price, LLM Abacus prices.json (attribution licence) |
| 16 | LiteLLM model_prices_and_context_window.json license | F-0003#22 | community-edited price file; licence not shown in results (checked on the repo) |
| 17 | models.dev open-source database AI model pricing context limits api.json repository | (c) "the choice of ... a model" | models.dev: TOML per provider/model, generated `api.json`, schema-checked PRs |
| 18 | git notes record AI agent session tokens cost per commit tool | F-0003#50 "in the project repository"; F-0001#1 | commit trailers (copilot-session-usage), pi-per-commit-spend, Ingot (joins logs to commits) |
| 19 | Agent Trace specification AI code attribution git open spec | F-0001#1 "Git is the system of record" | Cursor's Agent Trace RFC: JSON trace records of AI contributions per file/line |
| 20 | AGENTS.md open format agent instructions harness neutral rules multiple coding agents | F-0003#26, #27, #52 | AGENTS.md (Agentic AI Foundation, Linux Foundation); Claude Code needs an import or a symlink |
| 21 | sync rules across AI coding agents one source CLAUDE.md GEMINI.md .cursor/rules generator tool ruler | F-0003#27 "No procedure gives the same rules to a different harness agent" | Ruler, rulesync, agentsync, crossby, block/ai-rules |
| 22 | vendor-neutral spec agent rules skills hooks across harnesses "skills" SKILL.md open standard agentskills | F-0003#52 | Agent Skills (SKILL.md, agentskills.io), `.agents/skills/`; no neutral standard for hooks |
| 23 | prevent coding agent from modifying protected files CI gate paths sandbox read-only claude code codex | F-0001#3; F-0003#43, #64 | Codex bubblewrap re-binds `.git`, `.codex` read-only; Claude Code sandbox covers Bash only, file tools need deny rules; ralon `agent.lock` |
| 24 | pre-receive hook protected paths agent cannot modify tests CI config tamper detection LLM agent reward hacking tests | F-0003#64 "known-bad commits", "0 agent writes to rule paths" | practitioner guides: tests outside the writable set, CI diff check, server-side hooks; local hooks can be deleted by the agent |
| 25 | GitHub push rulesets restrict file paths block push modifying .github/workflows | F-0003#64 | push rulesets with file-path restriction (private/internal repos, Team/Enterprise plans) |
| 26 | anthropic sandbox-runtime srt filesystem network restrictions bubblewrap seatbelt any process | "isolates"; F-0001#3 | sandbox-runtime (`srt`): OS-level write/network limits on any process, research preview |
| 27 | container-use dagger environments coding agents each agent own container git branch | "isolates"; M2b "commit lands on its task branch" | container-use: one container and one `container-use/<env>` branch per agent |
| 28 | JSON Schema validated handoff artifacts between AI agents roles pipeline structured state file | F-0003#45, #59 | practitioner patterns (versioned JSON Schema per handoff, validate on write and read); no standard |
| 29 | A2A Agent2Agent protocol task artifact JSON schema specification Linux Foundation | F-0003#45 | A2A Task/Artifact model, JSON Schema 2020-12 generated from `a2a.proto`; an HTTP protocol |
| 30 | route coding tasks to different coding agents and models by cost difficulty open source router claude code codex | (c) | claude-code-router, OpenCodex (provider proxies), KlaatCode (own agent); none chooses per task outside a harness |
| 31 | golang orchestrator spawns claude code codex CLI subprocess stream-json parse go library | F-0004#1 Go, standard library only | Go wrappers claude-go, agent-sdk-go, claudecli-go; cloadex runs `claude -p` and `codex exec --json` |
| 32 | Armature conventions AI agents software project rules gates | F-0001#6 "Armature R5"; F-0003#45 "Armature conventions" | nothing named Armature |

## Candidates

Licence, language and last activity are from each candidate's own page on 2026-10-09 (see Limits
for how). "State in Git" says where the candidate keeps its own state. No candidate can be a Go
dependency (F-0004#1); "keep as" says how it can still serve.

| Name and URL | Licence | Language | Last activity | Part | State in Git | Deterministic | Keep as | Why |
|---|---|---|---|---|---|---|---|---|
| Claude Code headless mode, https://code.claude.com/docs/en/headless | proprietary (© Anthropic) | not stated | release v2.1.295, 2026-10-08 | a b c d f g | settings in `.claude/` (yes); transcripts in `~/.claude` (no) | no | protocol, format | `-p --output-format json/stream-json`; `system/init` gives `claude_code_version`, `capabilities[]`; result has `total_cost_usd` (client estimate), `duration_ms`, `duration_api_ms`, `usage`, `num_turns`, `subtype` (`error_max_budget_usd`...); `--max-turns`, `--max-budget-usd`; deny rules `Edit(path)`; `--json-schema`; no wall-clock flag |
| Codex CLI `codex exec`, https://github.com/openai/codex | Apache-2.0 | not stated (`codex-rs` dir) | release 0.162.0, 2026-10-08 | a c d f g | rollouts in `~/.codex` (no) | no | protocol, format | `--json` JSONL; `turn.completed.usage` (`input_tokens`, `cached_input_tokens`, `output_tokens`, `reasoning_output_tokens`), no cost, no duration; `--sandbox workspace-write` keeps `.git`, `.codex`, `.agents` read-only; `--output-schema`; hooks (beta); no time, turn or budget flag |
| Gemini CLI headless mode, https://geminicli.com/docs/cli/headless | Apache-2.0 | not stated | release v0.63.0, 2026-10-06 | a b c d | `.gemini/settings.json` (yes) | no | protocol, format | `-o json` gives `response`, `stats` ("token usage and API latency"), `error`; inner `stats` fields not in the official schema; exit 53 on turn limit (`model.maxSessionTurns`); `--version`; no cost; unpaid tier moved to Antigravity CLI |
| Agent Client Protocol (ACP), https://agentclientprotocol.com | Apache-2.0 | Rust | release v1.10.2, 2026-10-01 | a b d g | n/a | n/a | protocol | stdio JSON-RPC; `initialize` returns `protocolVersion`, `agentCapabilities`, `agentInfo{name,version}`; `session/cancel`; stop reasons; `usage_update` (`used`, `size`, `cost`); JSON Schema published; adapters for Claude Code, Codex, Gemini CLI, OpenCode |
| NEEDLE, https://github.com/jedarden/NEEDLE | MIT | Rust | pushed 2026-10-08 | a b c d | SQLite in workspace; work committed (partly) | yes ("a state machine") | pattern | `needle test-agent` checks executable, version probe and prompt transport without a model call; exit-code table (0 close, 124 timeout defers, >128 alert); headless flags for `claude -p`, `codex exec`, `opencode run` |
| mco, https://github.com/mco-org/mco | MIT | Python | pushed 2026-10-07 | a b c g | `run.json` in an artifact dir (no) | ordering only | pattern | adapter contract "detect, run, poll, cancel"; statuses success/failed/timeout/cancelled; capability-probe spec C0–C6 with "Capture version" and `min_version`; `mco doctor --json` |
| gate4agent, https://github.com/ZENG3LD/gate4agent | MIT | Rust | pushed 2026-10-04 | a b c d | `~/.gate4agent` (no) | not stated | pattern | `probe_all()` caches `CliCapabilities` (models, permission modes, features) per CLI; NDJSON parsers for Claude Code, Codex, Kimi, Grok; 9 stars |
| agents-cli (now agi-cli), https://github.com/phnx-labs/agents-cli | FSL-1.1-Apache-2.0 | not stated | pushed 2026-10-09 | a b e | `~/.agents` (no) | partly | pattern | version pins per harness and a per-version feature table "enforced at sync time"; `--timeout`; spend caps; one AGENTS.md to each harness; FSL bars code reuse |
| scion, https://github.com/GoogleCloudPlatform/scion | Apache-2.0 | Go | pushed 2026-10-09 | a b c d | `.scion/` (gitignored), Hub DB (no) | not stated | pattern | container per agent; `max_turns`, `max_duration`; static capability matrix of 9 harnesses; tier aliases small..extra-large; normalized `scion.usage.tokens`, `gen_ai.api.duration`; pre-1.0 |
| sortie, https://github.com/sortie-ai/sortie | Apache-2.0 | Go | release 1.26.0, 2026-10-01 | a c d | SQLite next to the workflow file (no) | not stated | pattern | `max_turns`, `turn_timeout_ms`, `max_tokens`, `claude-code.max_budget_usd`; run history: tokens, "usage measured" flag, start/end, configured vs reported model, exit type |
| Galley, https://github.com/shinpr/galley | MIT | Go | pushed 2026-09-15 | a c f g | `~/.galley` (no); result as a PR | not stated | pattern | isolated checkout per task; executor/model/effort per task; `scope.forbidden_paths`; verdicts accepted/needs revision/needs review/hard stop; run-evidence folder |
| Coder Eval, https://github.com/UiPath/coder_eval | Apache-2.0 | Python | release v0.12.12, 2026-10-06 | a c d g | run directory (no) | no | pattern, format | harness is one field (`agent.type`); YAML tasks; report schema (run.json, task.json); "Run-Limit Parity" across harnesses; tokens and cost per run |
| Harbor, https://github.com/harbor-framework/harbor | Apache-2.0 | Python | PyPI 0.24.0, 2026-10-05 | a b c d g | local `jobs/` (no) | not stated | pattern | installed-agent adapters for claude-code, codex, gemini-cli, opencode (`install`, `run`, `version()`); container per trial; timeouts; `n_input_tokens`, `n_output_tokens`, `cost_usd` "as reported by the agent"; `started_at`/`finished_at` |
| Bernstein, https://github.com/sipyourdrink-ltd/bernstein | Apache-2.0 | Python | PyPI 3.21.0, 2026-10-05 | a b c d f | `.sdd/` in project (partly committed) | yes ("No LLM in the coordination loop") | pattern | adapter per CLI; worktree per task; `max_agent_runtime_s` with SIGTERM then SIGKILL; `--max-cost-usd`; `role_model_policy` per role; gate fails a diff that touches run-config paths; beta, solo-maintained |
| MartinLoop, https://github.com/Keesan12/martin-loop | Apache-2.0 | TypeScript | pushed 2026-10-08 | a b c d f | `~/.martin` (no) | no | pattern | `--budget-usd`, `--max-tokens`, `--max-iterations`; cost provenance actual/calculated/estimated/unavailable; `--allow-path`/`--deny-path`; signed run receipts; no wall-clock cap |
| headless-coder-sdk, https://github.com/OhadAssulin/headless-coder-sdk | MIT | TypeScript | pushed 2026-06-23 | a d g | not stated | not stated | pattern | one event list over Codex, Claude, Gemini (init, message, tool_use, usage, error, cancelled, done); `maxBudgetUsd`; `outputSchema`; 39 stars |
| ordewell, https://github.com/ordewell/ordewell | Apache-2.0 | TypeScript | pushed 2026-10-08 | a b c g | `.ordewell/sessions` in project | yes (completion gate) | pattern | runner, model, effort per task (Claude Code, Codex, OpenCode); done only on `task_complete` evidence, "never decided by a model's opinion" |
| sandbox-runtime (`srt`), https://github.com/anthropics/sandbox-runtime | Apache-2.0 | TypeScript | release v0.0.79, 2026-10-07 | a f | settings file, can live in the repo | yes (OS-enforced) | tool | wraps any process; writes denied unless `allowWrite`, `denyWrite` takes rule and gate paths back; domain allowlist; bubblewrap on Linux; "research preview" |
| container-use, https://github.com/dagger/container-use | Apache-2.0 | Go | release v0.4.2, 2025-08-19; pushed 2026-09-21 | a | branch per agent in the project repo | not stated | pattern | "a fresh container in its own git branch"; MCP server, so the agent must call it; experimental; needs Dagger |
| Brood Box, https://github.com/stacklok/brood-box | Apache-2.0 | Go | pushed 2026-09-16 | a f | `~/.config/broodbox`; ephemeral (no) | not stated | tool | microVM per session (libkrun/KVM); writes reach the tree only through a diff engine with hash re-check; egress profiles; "EXPERIMENTAL" |
| ccusage, https://github.com/ccusage/ccusage | MIT | TypeScript | release v20.0.26, 2026-09-27 | d | reads home-dir logs (no) | not stated | tool, data source | `session --json` with tokens and `totalCost` for Claude Code, Codex, OpenCode, Gemini CLI and more; prices from a pinned LiteLLM file; no duration |
| CodeBurn, https://github.com/getagentseal/codeburn | MIT | TypeScript | pushed 2026-10-08 | d | `~/.cache/codeburn` (no) | partly | data source | one docs page per harness with log path, fields and dedup rules; prices from LiteLLM; some tokens estimated |
| agenttrace, https://github.com/luoyuctl/agenttrace | MIT | Rust | pushed 2026-10-06 | d | none stated | not stated | data source | parser guide for ~14 harness log formats; per session wall-clock, latency, tokens; evidence level Detailed/Aggregate/Limited; costs are estimates |
| agentacct, https://github.com/mikehasa/agentacct | MIT | Python | pushed 2026-10-03 | d | `~/.local/state` (no) | not stated | pattern | marks each number `client_reported` or estimate, "Missing beats wrong"; per-task steps with tokens; early alpha |
| LiteLLM price file `model_prices_and_context_window.json`, https://github.com/BerriAI/litellm | MIT (root; `enterprise/` excluded) | JSON | pushed 2026-10-09 | c d | a file to vendor at a pinned commit | yes when pinned | data source | `input_cost_per_token`, `output_cost_per_token`, `cache_read_input_token_cost`, `max_input_tokens`, `supports_*`; USD per token; used by ccusage, CodeBurn, tokscale |
| models.dev, https://github.com/anomalyco/models.dev | MIT | TOML, TypeScript | pushed 2026-10-09 | c d | TOML files to vendor | yes when pinned | data source | `cost.input`, `cost.output`, `cost.cache_read` in USD per million tokens; `limit.context`, `tool_call`, `reasoning`; submissions checked against a schema |
| OpenLIT `assets/pricing.json`, https://github.com/openlit/openlit | Apache-2.0 | not stated | pushed 2026-10-08 | d | a file to vendor | n/a | data source | `chat.<model>.promptPrice`, `completionPrice`, cache prices; units not stated (values match USD per 1K tokens) |
| OpenTelemetry GenAI conventions, https://github.com/open-telemetry/semantic-conventions-genai | Apache-2.0 | YAML/Markdown | commit 2026-10-07 | d | n/a | n/a | format | field names `gen_ai.usage.input_tokens`, `gen_ai.request.model`, `gen_ai.provider.name`, `gen_ai.invoke_agent.duration`; all "Status: Development"; no cost attribute |
| ATIF (Agent Trajectory Interchange Format), https://github.com/harbor-framework/harbor/blob/main/rfcs/0001-trajectory-format.md | not stated (host repo Apache-2.0) | JSON | RFC "Active", v1.7 in use | d g | a file, place chosen by the user | n/a | format | `final_metrics.total_prompt_tokens`, `total_completion_tokens`, `total_cached_tokens`, `total_cost_usd`; `agent{name, version, model_name}`; step timestamps; a validator exists (Python) |
| aGiTrack, https://github.com/core-aix/agitrack | Apache-2.0 | Python | pushed 2026-10-08 | a d | yes: commit body, `refs/agitrack/*`, git notes | partly | pattern, format | per-turn commit with "# aGiTrack Metadata": backend, backend version, model, tokens, started/ended times; no cost; worktree per session |
| shift-log, https://github.com/re-cinq/shift-log | AINAL (not SPDX) | Go | pushed 2026-09-12 | d | yes: `refs/notes/shiftlog` | not stated | pattern | JSON note per commit: `version`, `session_id`, `agent`, `model`, `effort{turns, input_tokens, output_tokens, cache_*}`; Claude Code, Codex, Gemini CLI, OpenCode; no cost or duration |
| copilot-session-usage, https://github.com/gsemet/copilot-session-usage | MIT | Python | PyPI 0.9.2, 2026-10-08 | d | yes: commit trailers | not stated | pattern | writes usage into the commit message as trailers (`Copilot-Session-Usage-Acc: <model>,in:..,out:..,cache:..`, `...-AIC: 232`); Copilot only; no duration |
| agentabi, https://github.com/Oaklight/agentabi | MIT | Python | PyPI 0.4.0, 2026-09-08 | a b d | none (library) | not stated | pattern | `detect_agents()`; `get_agent_capabilities()` with `supports_streaming`, `supports_mcp`, `supports_session_resume`...; `SessionResult{usage, cost_usd, duration_ms}`; Claude Code, Codex, OpenCode; alpha |
| AGENTS.md, https://agents.md | MIT (site repo) | Markdown | pushed 2026-09-10 | e | yes | n/a | format | plain Markdown, nearest file wins; read by Codex, OpenCode, Gemini CLI (`context.fileName`), Aider (`read:`), Claude Code natively from v2.1.277 when no CLAUDE.md exists |
| Rulesync, https://github.com/dyoshikawa/rulesync | MIT | not stated | commit 2026-10-08 | e f | yes (`.rulesync/`) | not stated | tool, format | one source to ~65 targets: rules, ignore, MCP, hooks, permissions; hooks translated per harness |
| AIWG, https://github.com/jmagly/aiwg | MIT | TypeScript | commit 2026-10-05 | b e | `.aiwg/` (yes) | no | data source | table of native rule paths for 18 harnesses; `aiwg doctor` checks version and deployed providers |
| agnix, https://github.com/agent-sh/agnix | MIT OR Apache-2.0 | Rust | release v0.57.0, 2026-10-08 | e | none | not stated (static linter) | tool | lints CLAUDE.md, AGENTS.md, SKILL.md, GEMINI.md, hooks, MCP and Codex/OpenCode config; GitHub Action |
| agents-md-cookbook, https://github.com/Taiizor/agents-md-cookbook | MIT | Markdown, TypeScript | commit 2026-08-08 | e | yes | n/a | data source | dated matrix of how each harness reads AGENTS.md (native, adapter, config); 21 stars |
| chock, https://github.com/open-coder-ai/chock | Apache-2.0 | Python | commit 2026-10-09 | e f | yes (`.agents/policies/`, `chock.lock`) | yes ("no model") | tool, pattern | one policy compiled to an AGENTS.md block, native deny rules per client, a git hook and a CI gate; `protect-agent-config`; SARIF; 9 stars |
| turbo-flow, https://github.com/marcuspat/turbo-flow | MIT | Bash | commit 2026-10-09 | e g | yes (`rig-lite/`) | first stage only | pattern | constitution file; fail-closed gate; reviewer of another model family must return parseable APPROVED/REVISE; refuses same-family review (F-0003#66) |
| gh-aw, https://github.com/github/gh-aw | MIT | Go | release v0.91.6, 2026-10-08 | a c d f g | workflows yes; run data in Actions (no) | not stated | pattern | protected by default: AGENTS.md, CLAUDE.md, GEMINI.md, `.github/`, `.agents/`, CODEOWNERS; safe-outputs typed JSON applied by a separate job; audit: wall time, turns, tokens, cost |
| CC Safety Net, https://github.com/kenryu42/claude-code-safety-net | MIT | TypeScript | release v2.6.1, 2026-10-08 | f | `.cc-safety-net/policy.json` (yes) | not stated | tool | `deny_paths` hard stop across 16 harnesses; the policy file protects itself; fails open on a broken config |
| loopgate_harness, https://github.com/rxdt/loopgate_harness | MIT | Python | commit 2026-09-30 | a f | yes | gates yes | pattern | forbidden paths are unstaged at commit; `max_iterations`, `max_minutes`; "Only humans can bypass triggered gates"; 23 stars |
| no_human, https://github.com/no-human-ai/no_human | MIT | Python | commit 2026-10-05 | f g | `~/.no_human` (no) | guard yes | pattern | tamper guard fails on fewer tests or assertions, "No model judgement"; reviewer of another model returns a boolean verdict, fails closed |
| DOS Kernel, https://github.com/anthony-chaudhary/dos-kernel | MIT | Python | commit 2026-09-27 | f g | Git history, `dos.toml` | from artifacts only | pattern | `dos verify PLAN PHASE` from commit-subject stamps (exit 0/1); `scope-gate` refuses writes outside the declared tree |
| Git server-side hooks `pre-receive`/`update`, https://git-scm.com/docs/githooks | Git's licence (not checked) | part of Git | docs of Git 2.54.0, 2026-04-20 | f | hook scripts beside the receiving repo, outside the tree | yes (exit code) | pattern | a non-zero exit rejects a push or one ref; the hook sees `<old> <new> <ref>` and can diff the pushed objects for rule and gate paths; guards pushes only, not local writes |
| GitHub push rulesets and CODEOWNERS, https://docs.github.com/en/repositories/configuring-branches-and-merges-in-your-repository/managing-rulesets/about-rulesets | proprietary service | n/a | changelog 2026-08-25 | f | rulesets on GitHub's server (no); CODEOWNERS in the repo | yes (server) | tool | "Restrict file paths" blocks pushes that change listed `fnmatch` paths, private/internal repos, Team/Enterprise; "Require review from Code Owners" |
| JSON Schema 2020-12, https://json-schema.org/specification | IETF Trust; code under Revised BSD | n/a | draft of 2022-06-16 (current) | g | schema files in the repo | yes (validation) | format | the common form for validatable handoffs (used by agent-context-store, ACP, Codex `--output-schema`, Claude `--json-schema`); Go `encoding/json` has no schema validation, only `DisallowUnknownFields` |
| agent-context-store, https://github.com/max9159/agent-context-store | MIT | TypeScript | release v0.3.1, 2026-07-06 | g | yes (`.acs/`) | not stated | format | YAML handoff record per role (ba, sa, dev, qa) validated against a JSON Schema; `acs validate`; 6 stars |
| Crewplane, https://github.com/crewplaneai/crewplane | Apache-2.0 | Python | commit 2026-10-09 | a g | `.crewplane/` in project | `mock` provider only | pattern | Markdown workflows with `schema_version`; `events.ndjson`, `manifests/run.json`; eight harnesses; does not sandbox |
| AgentPlane, https://github.com/basilisk-labs/agentplane | MIT | not stated (Node) | commit 2026-10-08 | d e f g | yes (`.agentplane/`) | yes (control plane) | pattern, format | `acr.json` Agent Change Record; WorkOrder with result schema, writable roots; missing telemetry marked "partial or unavailable" |

## The lists

An entry touches the problem when its own text shows that it does one of (a) to (g) for coding-agent
harnesses. Three kinds of entry get no line of their own, but are named once per list, and all are
rejects: a human-facing UI (desktop app, TUI, tmux grid, dashboard, widget), which runs sessions for a
person and not for a controller; a plain harness agent whose limits or telemetry work only inside
itself; and an LLM API gateway or hosted observability service, which measures API calls and not
what the harness reports. "keep" means the entry is in the Candidates table.

### andyrewlee/awesome-agent-orchestrators (324 lines)

Human UI, no line: agent-console, agent-deck, agent-manager, agent-session-manager, agent-of-empires, agterm, amux, Calyx, claude-squad, cmux, dmux, openkanban, pappardelle, Prowl, repomon, Rove, termany, thurbox, Vigil, agent-orchestrator, agent-session-manager-desktop, agent-squid, AGX, ai-maestro, AI4Kanban, aizen, alas, Alethe, Aperant, ateam, automaker, bb, Berd, Claude Command Center (CCC), clave, clideck, CodeNomad, constellagent, diri, dorothy, Dray, Emdash, EvoFlux, Fletch, Garcon, Ghostex, GraphCode, Helicon, IM.codes, Intent, intentic, ivy-tendril, jean, kandev, MonoCode, MulmoTerminal, MultiCC, muxel, nimbalyst, Nuphos, octomux, omg.dev, OpenChamber, Open Session, Orca, Paseo, Ouijit, parallel-code, Podium ADE, Pragma, Proliferate, qm, Reemoat, Runner, supacode, superset, synara, t3code, takopi, Tempest, tlbx, Tortie, Traycer, vibe-tree, vibecraft, Vicoa, Waku, wmux, xum, Zaivern Code, Zeron, zuse, Agent Teams, claude_codex_bridge, openswarm, Rovai AI, ralph-tui, Taskuary, codecast, PersonalJarvis, OpenMausBot, vibe-kanban.

- agentbox: reject: VM/cloud sandbox, heavier than Brood Box
- ai-devkit: reject: tmux control plane; AIWG covers rule sharing
- Cyclops: reject: tmux mailbox, no validatable schema
- herdr: reject: terminal runtime; no limits, telemetry or result reading
- YYLO: reject: Pi and Codex only; same as Crewplane
- humanlayer: reject: deprecated per the list
- AgentBridge: reject: live peer session, not a controller
- agentsmesh: reject: remote-workstation platform
- ClawTeam: reject: agent-spawned teams; inbox is not a schema
- CompanyHelm: reject: distributed platform; chat between agents
- Fusion: reject: multi-node platform; gates not specified
- gastown: reject: Git-backed tracking, but no schema or telemetry stated
- hcom: reject: agent messaging, no schema
- kodo: reject: verify loop; no format stated
- loki-mode: reject: BUSL-1.1 licence
- multi-agent-shogun: reject: tmux hierarchy of agents
- NXTG-Forge Orchestrator: reject: formats not described
- OpenRig: reject: YAML teams for persistent tmux, not headless
- ORCH: reject: same as Crewplane
- Orkas: reject: an LLM commander dispatches; not deterministic
- OtoDock: reject: FSL-1.1 licence
- paperclip: reject: platform; budget caps same as MartinLoop
- Raven: reject: an LLM host plans; not deterministic
- ruflo: reject: swarm meta-harness, vague
- scion: keep: Go; `max_turns`/`max_duration`, capability matrix, token metrics
- shire: reject: mailboxes, not a schema
- squad: reject: Copilot only
- tutti: reject: its "typed artifacts" have no schema
- bernstein: keep: deterministic coordination over many CLI agents
- Dex: reject: human-gated, LLM reviewers
- fractal: reject: same as MartinLoop
- LoopTroop: reject: LLM council; OpenCode only
- MartinLoop: keep: USD/token/iteration caps, cost provenance, receipts
- ordewell: keep: runner, model, effort per task; evidence gate
- ralph-claude-code: reject: Claude Code only
- ralph-orchestrator: reject: same limits as scion and sortie; cost-cap docs inconsistent
- ralphex: reject: same as ordewell
- aeon: reject: tied to GitHub Actions
- background-agents: reject: cloud sandboxes, hosted triggers
- centaur: reject: Kubernetes/Slack platform
- claude-code-action: reject: Claude only GitHub Action
- codex-action: reject: Codex only GitHub Action
- Contrabass: reject: issue runner; no formats described
- cyrus: reject: issue-watcher service
- gh-aw: keep: protected-files list, safe-outputs, audit fields
- lalph: reject: vague issue runner
- Machinist: reject: Codex only
- multica: reject: managed platform
- NEEDLE: keep: deterministic dispatch, `test-agent` probe, exit-code table
- no_human: keep: deterministic test-tamper guard, fail-closed reviewer
- open-swe: reject: cloud sandbox service
- OpenHands: reject: large platform (it shows ACP driving)
- Orbi: reject: Pi only
- remote-swe-agents: reject: AWS serverless
- run-gemini-cli: reject: Gemini only GitHub Action
- sortie: keep: Go; per-agent limits, run-history fields
- SuperPlane: reject: platform; no formats stated
- symphony: reject: Codex App Server only
- Tale: reject: sandbox platform
- Agent Messaging Protocol: reject: signed messaging, not a handoff schema
- agent-runbook: reject: same as AIWG
- Agentlas OS: reject: an LLM orchestrates each task
- agenttier: reject: Kubernetes runtime
- aGiTrack: keep: per-turn tokens in commit message, refs, notes
- AIWG: keep: one rule source to 18 harness paths
- Archon: reject: same as Crewplane
- AX: reject: Kubernetes runtime
- Claudexor: reject: subscription rotation, not model choice
- Concord MCP: reject: MCP coordination server, not a file schema
- Cotal: reject: NATS pub/sub
- Crewplane: keep: on-disk handoffs and run manifest, eight harnesses
- foremerge: reject: SQLite, not Git
- handoff: reject: in-session delegation
- Heym: reject: OpenCode only
- oh-my-codex: reject: Codex only
- omnigent: reject: state in `~/.omnigent`, no time limit; MartinLoop and Brood Box cover it
- Open Multi-Agent: reject: TypeScript runtime library
- OrcaReplay: reject: socket recorder, not harness-reported
- sandbox-agent: reject: cloud-sandbox daemon
- skillfold: reject: skills only; same as AIWG
- sub-agents-skills: reject: thin; no limits or telemetry
- 1code: reject: archived
- CodexMonitor: reject: inactive, Codex only
- gnap: reject: inactive since 2026-03
- opengoat: reject: inactive
- ralphy: reject: inactive bash loop
- subtask: reject: inactive Claude skill
- swarm-protocol: reject: inactive MCP coordination
- wreckit: reject: inactive loop

### bradAGI/awesome-cli-coding-agents (884 lines)

Human UI, no line: Yaw, Orca (Stably), Multica, herdr, AionUi, vibe-kanban, cmux, Paseo, Superset, Agent Orchestrator (AO), Claude Squad, Emdash, CodexMonitor, Toad, agent-of-empires, Crystal, qwen-audio-agent, supacode, Agent Teams AI, Cate, Orkas, mux, Nimbalyst, jean, Parallel Code, agent-deck, Agent Sessions, agent-manager, Proliferate, amux (mixpeek), Catnip, Vicoa, ntm, Calyx, vibe-tree, MulmoTerminal, Ivy Tendril, Tempest, amux (andyrewlee), CliDeck, GridBash, tlbx, ADHDev, Garcon, intentic, run-kit, Better Agent, Clave, showagent, vibepanel, repomon, Claudescope, construct, Podiom, CLITrigger, cliclaw, Claudette, tring, Agent Workbench, mix2, PATAPIM, Bwee, Even, Unpeel, defract, CodeAgentSwarm, TermLink, VibeFuse, Traycer, 5dive, LoopTroop, Untether.
Harness agents, no line: Oh My OpenAgent, Codewhale, Kimi Code, Maki, zot, Kolega Code, Kolkrabbi, molt, GitClaw, pool, DeerFlow, SandBase Harness.

- VibePod: reject: licence not stated; same as Brood Box
- Martty: reject: ACP client TUI
- CLI Agent Orchestrator (CAO): reject: tmux plus SQLite; no limits or telemetry stated
- hcom: reject: unstructured message bus
- AgentBox: reject: VM/cloud sandbox, heavier than needed
- Polter: reject: a lead agent types into terminals
- tmuxlet: reject: drives interactive TUIs through tmux
- TaskHandoff: reject: Codex-only Docker; format not specified
- postmortemthis: reject: LLM verdict
- AgentX: reject: 0 stars; same as agents-cli
- clisweave: reject: same as agents-cli
- claude-flow: reject: no format or limits stated
- Symphony: reject: Codex App Server only
- Omnigent: reject: Python server, state outside Git
- Kiro Crew: reject: kiro-cli only
- ralph-orchestrator: reject: same as scion/sortie limits
- AgentsMesh: reject: BSL; remote platform
- zeroshot: reject: same as Crewplane
- Bernstein: keep: deterministic orchestrator, worktrees, budgets
- oh-my-agent: reject: LLM judges
- Loki Mode: reject: BUSL
- fractal: reject: same as MartinLoop
- Aeon: reject: GitHub Actions; LLM scoring
- h5i: reject: README now describes a web-security workspace
- MartinLoop: keep: budgets, caps, stop conditions, receipts
- Orbi: reject: fair-code licence
- Ordewell: keep: runner/model per task (licence is Apache-2.0)
- ORCH: reject: licence not stated
- outsourcerer: reject: non-commercial licence
- Crewplane: keep: on-disk handoffs, run manifest
- fab-kit: reject: per-stage model tiers not in its README
- NEEDLE: keep: deterministic dispatch, adapter probe
- ralph-harness: reject: Claude only; same as Galley
- Galley: keep: Go; forbidden paths, verdicts, run evidence
- Dahrk: reject: needs a hosted hub
- Ralph Workflow: reject: same as MartinLoop
- Relay: reject: runs inside the agent via MCP
- sage: reject: bash; same as NEEDLE
- Vibestrate: reject: 5 stars; ledger fields not stated
- baya-cli: reject: LLM-planned DAG
- DevPilot: reject: Claude-only runner
- Team1-Factory: reject: Claude only
- the-perfect-orchestrator: reject: Claude-only tmux
- OpenViking: reject: context database outside Git; AGPL
- claude-code-router: reject: provider proxy, Claude Code only
- Beads: reject: Dolt state, not Git
- NemoClaw: reject: OpenClaw only
- OpenWiki: reject: LLM-written AGENTS.md
- OpenCodex: reject: provider proxy, not choice per task
- Agent Executor (AX): reject: Antigravity only
- GitHub Agentic Workflows: keep: same as gh-aw
- numbat: reject: observes only; no tokens or cost; blocking off by default
- HOL Guard: reject: licence not stated
- AgentSight: reject: syscall traces, not tokens or cost
- Coasts: reject: service isolation beyond the need
- Concord MCP: reject: MCP server state
- AgentBridge: reject: live peer chat
- OrcaReplay: reject: TLS-interception replay
- cc-router: reject: provider proxy
- codex-profiles: reject: Codex only
- agenttrace: keep: wall-clock and latency per session from harness logs
- ax: reject: same as agenttrace
- ActPlane: reject: eBPF policy DSL heavier than path checks
- RoleCraft: reject: same as agents-cli
- AgentPlane: keep: state in `.agentplane/`, Agent Change Record
- brood-box: keep: microVM per session
- Okto Nexus: reject: SQLite hub
- AgentLint: reject: Claude plugin; same as agnix
- claudebox: reject: Claude-only sandbox
- AgentTier: reject: Kubernetes sandbox
- pi-reflect: reject: LLM-evolved rules
- AgentDiff: reject: line provenance; same as aGiTrack
- Loadout: reject: same as agents-cli
- Data Olympus: reject: MCP rule store, not repo files
- Mneme: reject: same as DOS Kernel
- grite: reject: Git-refs tracker; same as aGiTrack
- TraceFold: reject: Cedar/DSSE stack heavier than needed
- agents-cli: keep: version pins, per-version feature table
- schliff: reject: 16 stars; same as agnix
- machine: reject: Lima VM sandbox
- distro-rig-vps: reject: VM sandbox; GPL
- cowork-to-code-bridge: reject: Claude only
- myc: reject: SQLite
- agent-runbook: reject: compiles to SKILL.md
- agent-terminal: reject: PTY library
- agy-auto: reject: Antigravity only
- gate4agent: keep: `probe_all()`, `CliCapabilities`
- clu: reject: SQLite
- shim-cli: reject: secret masking
- SpecWave: reject: 6 stars; AIWG and Rulesync cover it
- PatchWarden: reject: same as gh-aw
- OSOP: reject: inactive since 2026-04; records written by the agent
- context-bridge: reject: unstructured handoff
- tu: reject: same as agenttrace
- Project Tiny Context Harness: reject: 4 stars
- agent-trace: reject: same as agenttrace
- Squelette: reject: 2 stars
- llm-panel: reject: LLM review panel
- CodeVetter: reject: same as no_human
- isitdone: reject: 1 star; the agent can remove its hook
- Agent Fleet: reject: closed

### hesreallyhim/awesome-claude-code (770 lines)

Human UI, no line: Claude Threads, Cate, CloudCLI (Claude Code UI), Happy Coder, Nimbalyst, Sidekick for Max, Vibeyard, Claude Squad, CC Harness, seedeep, c9watch, cctop, Claude Code Agent Monitor, claude-control, so-agentbar, AgentWatch, CCDash, ClaudeBar, Claumon, Pacer, claude-esp, Multi-Agent Observability, agents-observe.

- Claude Code Hooks: Complete Guide: reject: hook docs; official docs suffice
- Claude Code GitHub Action: reject: Claude only, GitHub CI
- Steering Claude Code: Skills, Hooks, Rules, Subagents and More: reject: article
- AI Research Skills: reject: research-domain skills; agent-context-store closer
- Codex Skill: reject: Claude starts Codex
- Fusion Harness: reject: Pi extension
- llm-router: reject: per-prompt proxy routing
- Agent Guard: reject: secret leaks, not rule paths
- aicontainer: reject: same as Brood Box
- Airut: reject: Claude only
- Brood Box: keep: microVM per session, diff-gated writes
- Claude Code Safety Net: keep: `deny_paths` across 16 harnesses
- Claude Code Safety Guard: reject: Claude Code only
- Cleat: reject: same as Brood Box
- Code on Incus: reject: same as Brood Box
- compass: reject: personal config bundle
- Container Use: keep: container and Git branch per agent
- Dippy: reject: bash auto-approval
- GouvernAI: reject: Claude only; same as Safety Net
- machine: reject: same as Brood Box
- Node9: reject: vague proxy governance
- SkilLock: reject: pins skills, not rule paths
- Agent Collab Skills: reject: no schema
- AgentSys: reject: Claude workflow bundle
- Claude Code Harness: reject: Claude Code only
- Claude Code Hook Comms (HCOM): reject: chat, no schema
- claude-intercom: reject: Claude only relay
- OpenRig: reject: interactive tmux
- Ralph for Claude Code: reject: Claude Code only
- Ralph Wiggum Plugin: reject: loop inside the session
- ralph-harness: reject: Claude only; Galley shows the pattern
- ralph-orchestrator: reject: same as scion/sortie limits
- The Ralph Playbook: reject: guide
- Callimachus: reject: history search
- claude-context-optimizer: reject: context plugin
- faf-cli: reject: context, not rules
- presence: reject: memory tool
- better-ccflare: reject: proxy dashboard
- ccusage: keep: per-session tokens and cost as JSON
- ccxray: reject: API proxy
- goccc: reject: same as ccusage
- toktrack: reject: same as ccusage
- Claude Code Observability Stack: reject: Docker Grafana stack
- OrcaReplay: reject: proxy record/replay
- clausona: reject: Claude only
- Rulesync: keep: one source to ~65 harness targets
- tweakcc: reject: patches the Claude Code install
- TDD Guard: reject: one harness; same as Safety Net
- agents-md-cookbook: keep: AGENTS.md compatibility matrix
- agnix: keep: linter for rule, hook, MCP files
- BlockWatch: reject: generic sync linter
- Ctxlint: reject: same as agnix
- Schliff: reject: same as agnix
- Upkeep: reject: LLM audit
- Claude CodePro: reject: Claude Code only
- claude-code-tools: reject: handoff without schema

### RoggeOhta/awesome-codex-cli (494 lines)

Human UI, no line: farion1231/cc-switch, The-Vibe-Company/companion, cnlimiter/codex-manager, dou-jiang/codex-console, milisp/codexia, xintaofei/codeg, zhu1090093659/CodeConductor, Dimillian/CodexMonitor, Cocoanetics/CodexMonitor, waskosky/agent-cli-farm, mixpeek/amux, kbwo/ccmanager, lsm1103/session-dashboard, michaelversus/BuildrAIApp, onewesong/codex-viz, SeemSeam/claude_codex_bridge, AgentsMesh/AgentsMesh, lorytek/PulseMeter, graykode/abtop, NihilDigit/waybar-ai-usage, steipete/CodexBar, GreenSheep01201/claw-empire, GreenSheep01201/Claw-Kanban, PleasePrompto/ductor, standardagents/dmux.

- Config Reference: reject: Codex-only config
- Subagents Docs: reject: fan-out inside Codex
- MCP Docs: reject: heavier than `codex exec`
- Hooks Docs: keep: part of the Codex CLI row (beta)
- Security & Sandboxing: keep: part of the Codex CLI row
- Non-interactive Mode: keep: part of the Codex CLI row
- Multi-Agent Guide: reject: inside the harness
- Exec Policy: reject: the sandbox protects paths better
- Speed / Fast Mode: reject: cost multiplier only
- AGENTS.md Guide: keep: part of the AGENTS.md row
- Codex Action: reject: GitHub Actions wrapper
- Codex Plugin for Claude Code: reject: delegation by an agent
- agents.md (Open Standard): keep: same as AGENTS.md
- codex-cli-best-practice: reject: Codex-only templates
- claude-codex-settings: reject: two hand-kept rule files
- caliber-ai-org/ai-setup: reject: same as Rulesync
- agentsmd/agents.md: keep: the neutral rule format
- Ischca/awesome-agents-md: reject: list of examples
- danielrosehill/Agents.md-Templates: reject: templates only
- CoderMageFox/claudecode-codex-subagents: reject: role prompts
- waltstephen/ArgusBot: reject: LLM supervisor
- atticus98/codex-turbo: reject: inside the harness
- sehoon787/my-codex: reject: Codex-only bundle
- shinpr/sub-agents-skills: reject: thin; run by the agent
- RBraga01/a-team: reject: role prompts
- Untrivial-ai/agent-orchestrator: reject: LLM-planned
- mco-org/mco: keep: adapter contract, capability-probe spec
- aannoo/hcom: reject: messaging, no schema
- obra/external-subagents: reject: launched by the agent
- affaan-m/ECC: reject: config bundle
- runkids/skillshare: reject: skills, not rules
- xmm/codex-bmad-skills: reject: no schema
- regenrek/codex-1up: reject: installer
- EveryInc/compound-engineering-plugin: reject: workflow prompts
- Yeachan-Heo/oh-my-codex: reject: Codex only
- shanraisshan/codex-cli-hooks: reject: starter examples
- vanthienha199/agent-cost-mcp: reject: cost MCP for the agent
- xiaolai/codex-octopus: reject: Codex-only MCP
- salparadi-labs/brain: reject: SQLite
- Mr-Tomahawk/codex-cli-mcp-tool: reject: MCP wrapper
- agency-ai-solutions/openai-codex-mcp: reject: same as codex-cli-mcp-tool
- Leonard013/BigBrain: reject: agent-to-agent bridge
- ching-kuo/claude-codex: reject: routing by Claude
- Dekelelz/let-them-talk: reject: chat broker
- just-every/code: reject: Codex fork
- rxdt/loopgate_harness: keep: forbidden paths unstaged at commit
- digipulse-engineering/GAAI-framework: reject: Elastic-2.0; "no programmatic enforcement"
- yazcaleb/rses: reject: transcripts, no validation
- gotalab/cc-sdd: reject: spec prompts
- re-cinq/shift-log: keep: session record in Git notes
- 2ue/ccman: reject: provider switcher
- Dicklesworthstone/coding_agent_session_search: reject: history search
- feiskyer/codex-settings: reject: provider configs
- ben-vargas/ai-sdk-provider-codex-cli: reject: JS library
- router-for-me/CLIProxyAPI: reject: model API proxy
- lich0821/ccNexus: reject: gateway proxy
- maksimzayats/acodex: reject: Python SDK of `codex exec`
- Adashuai5/quota-autopilot: reject: subscription quotas
- onurkanbakirci/awesome-codex-automations: reject: recipe list
- DeepMyst/Mysti: reject: same as mco
- athola/skrills: reject: skill format
- jcputney/agent-peer-review: reject: review run by the agent
- abhishekgahlot2/codex-claude-bridge: reject: same as agent-peer-review
- lbb00/ai-rules-sync: reject: same as Rulesync
- dsifry/metaswarm: reject: inside the harness
- pchalasani/claude-code-tools: reject: human tools
- Leoyang183/sync-agents-settings: reject: MCP configs, not rules
- oil-oil/codex: reject: same as Codex Plugin for Claude Code
- josstei/maestro-orchestrate: reject: inside the harness
- catlog22/maestro-flow: reject: same as maestro-orchestrate
- fahd09/watchtower: reject: API proxy
- yoavf/ai-sessions-mcp: reject: MCP for the agent
- ccusage/ccusage: keep: same as ccusage
- junhoyeo/tokscale: reject: same as CodeBurn; optional public leaderboard upload
- Codex Universal: reject: Codex-only image
- DeepBlueDynamics/gnosis-container: reject: Codex-only bundle
- libops/cli-sandbox: reject: no licence; 5 stars
- itscooleric/clide: reject: same as cli-sandbox
- Ducksss/codex-profiles: reject: account separation
- Codex Windows Permissions Guide: reject: Windows only
- Codex CLI Quick Start (JP Caparas): reject: blog
- Codex CLI Automation: 3 Workflow Patterns (SmartScope): reject: blog
- How to Run Codex CLI Safely in GitHub Actions (SmartScope): reject: blog
- Codex CLI vs Claude Code vs Gemini CLI (comparison table): reject: unsourced table
- duanyytop/agents-radar: reject: benchmark, no routing format
- Testing AI Coding Agents Benchmark (Render): reject: one-off article

### tensorchord/Awesome-LLMOps (801 lines)

Most entries (serving, training, vector stores, RAG and agent frameworks) do not touch the problem.
LLM API gateways and hosted or server observability, no line: Bifrost, SentryNode Gateway, Cordum, AgentCost, budget-guard, ClevAgent, Helicone (twice), Traceloop OpenLLMetry, Langfuse (twice), Lookspan, Maxim AI, Noveum, onWatch, traceAI, Future AGI (twice), Token Police, Weco Observe, witness, AgentField, Agnos, AI studio, AISIX, Arize-Phoenix, BitRouter, boldrouter, CrewDock, Doubleword Control Layer, FerryAPI, Glide, GoModel, Kunavo, MLflow, Laminar, LangKit, Lunary, Opik, OrcaRouter Lite, Quotaflow, TeamoRouter, TokenMix, Portkey, SAVI SDK, Spendline, TensorZero, ThinkWatch Lite, UnoRouter, XiuRouter, Relay, Nora.

- brood-box: keep: microVM per session
- DOS Kernel: keep: done claims checked against Git; scope gate
- CodeBurn: keep: per-harness log paths and prices
- Inferrail: reject: 1 star; API-layer gateway
- AgentMark: reject: prompt-file format, not handoffs
- AIWG: keep: same as AIWG
- claude-router: reject: embedding routing, Claude only
- depdesk: reject: code scanner
- LiteLLM: reject as a library; its price file is a candidate
- lintlang: reject: lints prompts, not rules
- magentic: reject: Python structured-output library
- OpenLIT: keep: its pricing JSON and GenAI field names
- Roundtable: reject: MCP hub; limits vague
- Trinity: reject: server platform
- Modelglass: reject: proprietary terms bar reuse of the data
- AgentsMesh: reject: heavy platform
- Bernstein: keep: same as bernstein
- Coder Eval: keep: harness as one field, report schema
- Cotal: reject: needs a NATS server
- fractal: reject: SQLite; same as MartinLoop
- promptext: reject: context packer
- AI Agent Token Cost Calculator: reject: browser estimate
- Awesome AI Coding Sandboxes: reject: a list, not a solution (see Limits)

## The search file

Each row is judged by its description in `sources/gh-search.tsv`, and for keeps also by its repository.
Off-topic rows get no line: pbi-cli, rikkahub-agent-pure, pxpipe, liteflow, ai-agent-unity-rules,
windmill-cli-docs, zalo-agent-cli, mirofish-cli, baguette, awesome-agentic-ai-zh, chatgpt-cli,
browser-act/skills, aionrs, analog-agents, open-reviewer, agentless-DaC, agent-journal, sftp,
KeyNanny, decrypt-tsd-files, portal-mcp-server, feedback-agent, local-protected-file-reader,
better-github-skill, agent-crm, headless-gtm, crm.cli, hope-agent, obsidian-llm-wiki,
subagent-cortex-code, emperor-agent, both genpark edge-inference rows, Light-Agentic-Template,
mercury-agent, AGI-Super-Team.
Human UI, no line: AgentLimits, TokenStep, Tokdash, tokentop, token-ledger (bisheshabramhacharya),
agent-token-ledger, sportnak5/token-ledger, Claude-Code-Usage-Monitor (both), CodexBar, claude-usage,
claude-hud, Clawdmeter, ClaudeBar, codenotch, codexU, vibe-usage-app, ClaudeCodeUsage,
codex-usage-badge, token-monitor, tokentelemetry, Tokens-4-Breakfast, llm-quota,
OrchestratorInc/agent-orchestrator, paseo, superset, claude_codex_bridge, claude-command-center,
cc-switch-web, VibeChard, DrSeedon/orchestra, agentic-tasktrees, 2code, herdr-worktree-agents.

- junhoyeo/tokscale: reject: same as CodeBurn; leaderboard upload option
- getagentseal/codeburn: keep: per-harness log paths, fields, prices
- DeepAgentLabs/agenticlens: reject: profiles LLM workflows, not harness sessions
- agentic-os-org/ANOLISA: reject: agentic OS platform
- splunk/token-meter: reject: dashboard; same as CodeBurn
- xiufengsun/TokenTracker: reject: native apps; same as CodeBurn
- stormzhang/token-tracker: reject: statusline, two harnesses
- ramtinJ95/opencode-tokenscope: reject: OpenCode only
- mag123c/toktrack: reject: same as ccusage
- Han-1413141/dsh-cost-meter: reject: DeepSeek Harness only
- mazzzystar/api-usage: reject: OpenAI API only; inactive since 2024
- juliantanx/aiusage: reject: same as CodeBurn
- oluwajubelo1/otellix: reject: Go SDK for LLM backends, not harness sessions
- inferock/inferock-bench: reject: API proxy, not harness-reported
- fjgbue/claude-delegator-deepseek-mcp: reject: delegation by the agent
- voly-codes/voly: reject: no Codex or Gemini executor; cost source not stated
- harpd-dev/llm-cost-benchmark: reject: benchmark tables, not a per-session record
- deeflect/smart-spawn: reject: 3 stars; inactive since 2026-03
- harpd-dev/cost-per-successful-task: reject: a metric, not a record format
- catalystneuro/llm-frontier: reject: price-capability tracker site
- adididitagain/tokentab: reject: in-app tracking
- bobbyhalljr/llm-cost-per-task: reject: calculator
- guevae2/llm-cost-per-outcome: reject: one-off benchmark
- OptimNow/cost-per-task: reject: API proxy
- Rylaa/fable5-opus5.5-orchestrator: reject: Claude Code plugin
- BrianWong05/TokenLedger: reject: 3 stars; same as CodeBurn
- camposvinicius/llm-gateway: reject: gateway server
- ethanplusai/agent-ledger: reject: SQLite, not Git
- JamesAnderson6217/nonprofit-agent-usage-ledger: reject: one Java loop, 0 stars
- hollis-labs/go-usage-ledger: reject: 1 star; one organisation's Go module
- realiti4/claude-swap: reject: account rotation
- vibe-cafe/vibe-usage: reject: same as CodeBurn
- douglasmonsky/codex-usage-tracker: reject: Codex only
- toroleapinc/claude-brain: reject: Claude only, across machines
- PanisHandsome/ai-rules-sync: reject: same as Rulesync
- sampleXbro/agentsmesh: reject: same as Rulesync, 25 stars
- Chemaclass/agnostic-ai: reject: same as Rulesync, 23 stars
- dhruv-anand-aintech/agent-rules-sync: reject: real-time daemon, 7 stars
- airulefy/Airulefy: reject: inactive since 2025-05
- martinmose/agentlink: reject: inactive since 2025-08
- saqibameen/agent-dotfiles: reject: same as Rulesync
- BayramAnnakov/claude-reflect: reject: LLM-learned rules
- iannuttall/source-agents: reject: inactive since 2025-10
- Signet-AI/signetai: reject: memory and secrets store
- intellectronica/claude-agentsmd: reject: Claude plugin; native AGENTS.md reading now exists
- Daliagents-com/agentsge-1: reject: memory/MCP sync, vague
- westonplatter/aps: reject: personal prompt collection
- shriramkv/agentsmd-sync: reject: same as Rulesync
- bensyverson/agents: reject: syncs across repos, not harnesses
- HKUDS/DeepCode: reject: its own harness
- Yeachan-Heo/oh-my-claudecode: reject: Claude Code only
- 21st-dev/1code: reject: inactive since 2026-03
- mattpocock/sandcastle: reject: TypeScript library; Brood Box covers isolation
- vijaythecoder/awesome-claude-agents: reject: Claude subagent prompts
- bradAGI/awesome-cli-coding-agents: reject: a list (read as a source)
- Danau5tin/multi-agent-coding-system: reject: its own harness
- awslabs/cli-agent-orchestrator: reject: tmux plus SQLite; no limits or telemetry
- wanshuiyin/Auto-claude-code-research-in-sleep: reject: research skills
- sno-ai/sno-station: reject: shared memory and messaging, no schema
- marcuspat/turbo-flow: keep: constitution, fail-closed cross-family review
- majiayu000/harness: reject: Rust server plus Postgres; same scope as Bernstein
- WenyuChiou/ai-research-skills: reject: research-domain skills
- gmickel/flow-next: reject: LLM cross-model review plugin
- EvilFreelancer/crossreview: reject: 6 stars; an LLM verifies
- Design-Arena/agent-runner: reject: its own harness
- darkzOGx/darkzloop: reject: its own runner, inactive
- TheArchitectit/AIGGP-Agentic-Framework: reject: 5 stars; generic CI gates
- getrunkite/runkite: reject: Postgres/Redis control plane
- theycallme-eric/agent-runner: reject: 0 stars
- phyothihakyaw/skillbench: reject: 0 stars
- panchew/local-agent-runner: reject: Ollama loop
- schuligan/agentic-test-runner: reject: 0 stars
- suchithnarayan/agent-security-hooks: reject: same as CC Safety Net, 8 stars
- Serhioromano/pi-defender: reject: Pi only
- ckckck/agent-safe-delete: reject: delete-to-archive only
- obstalabs/bulwark: reject: read gate, not write gate
- jmstar85/oh-my-githubcopilot: reject: Copilot only
- kenryu42/cc-safety-net: keep: same as CC Safety Net
- JeongJaeSoon/agent-guard: reject: secret leaks
- paudley/coding-ethos: reject: 6 stars; MCP plus CEL stack
- moerasermax/AgentCharter: reject: 2 stars
- Offsend/Offsend: reject: secret sealing
- open-coder-ai/chock: keep: one policy to deny rules, git hook, CI gate
- bram-wq/agent-guardrails: reject: 0 stars; same as CC Safety Net
- Quarkgluonmixture/coding-agent-guardrails: reject: 1 star bundle
- dabit3/agent-handoff: reject: in-process library; inactive since 2026-02
- openenvelope/schema: reject: tied to one runtime; fields not documented
- ZibbyDev/agent-workflow: reject: Zod, not JSON Schema; cloud product
- max9159/agent-context-store: keep: JSON-Schema-validated handoffs in Git
- agentsentinel/agentsentinel-handoff: reject: 0 stars
- faresrafat3/agent-handoff: reject: 0 stars
- moona3k/handoff: reject: handoff by URL through a hosted worker
- talocode/handofflane: reject: npm library, 0 stars
- roymcfarland/agent_handoff_guide: reject: playbook, no schema
- shin4141/decision-os-v12-completion-integrity: reject: 2 stars
- GDWN-BLDR/stateweave: reject: framework state migration
- yunaremaia/agent-workspace: reject: 1 star
- Lians-ai/Lians: reject: same as DOS Kernel
- hoeggsoftware/multi-agent-planning: reject: 2 stars; a model only
- llblab/pi-state-flow: reject: Pi only
- ychamel/RepoResident: reject: Claude-centred add-on
- saint0x/tl: reject: checkpoint streams, not records
- vekexasia/pi-extensible-workflows: reject: Pi only
- a5c-ai/babysitter: reject: description vague; NEEDLE states the pattern
- luckeyfaraday/athena-loops: reject: same as NEEDLE
- Extra-Chill/homeboy: reject: 16 stars, vague
- ray-amjad/claude-code-workflow-creator: reject: Claude Code only
- RchGrav/astraeus: reject: a prompt; inactive
- omar-os/omar: reject: same as NEEDLE
- NTCoding/autonomous-claude-agent-team: reject: Claude only
- mjasnikovs/pi-task: reject: Pi only
- Ali-hey-0/ai-runtime-lab: reject: learning lab
- OhadAssulin/headless-coder-sdk: keep: one event list over three harnesses
- nkhdiscovery/headless-agentic-codebase: reject: 10 stars; self-merging loop
- ClipboardHealth/groundcrew: reject: interactive agents
- DnzzL/herdr-automations: reject: herdr plugin
- errhythm/copse: reject: 2 stars
- vraj00222/agent-farm: reject: Claude only
- mutheejj/worktree-orchestrator: reject: 0 stars
- FabianoArthur/claude-code-kitchen: reject: Claude only
- ex3del/mast: reject: Claude plugin
- mikehasa/agentacct: keep: reported vs estimated vs missing labels
- genai-telemetry/genai-telemetry: reject: SDK exporting to servers
- PxA-Labs/AgentsScope: reject: in-process Python
- SahilSelokar/Tokensense: reject: in-process SDK

## What no candidate covers

- (b) A probe session. NEEDLE (`test-agent`, no model call), the mco probe spec, gate4agent and
  agents-cli check binaries, versions or static feature tables. None runs one small real task on each
  harness and commits a dated probe record (version, flags, the fields the harness reports) before
  work is given. Codex documents no version field over exec or app-server; Gemini CLI has only
  `--version`.
- A wall-clock limit inside the harness. `codex exec` and Gemini CLI have no time flag; Claude Code
  has `--max-turns` only. Every orchestrator enforces time from outside (scion `max_duration`,
  sortie `turn_timeout_ms`, NEEDLE exit 124), so the controller must own the timer and kill the
  process group.
- A money limit and a money record that every harness reports. Only Claude Code reports
  `total_cost_usd`, and it is a client estimate. Codex and Gemini report tokens and no cost; ACP
  `usage_update.cost` is optional. Cost for those must be computed from a price list, and no source
  defines the cost of a session on a subscription plan.
- One complete telemetry row per task in the repository. aGiTrack and shift-log write tokens into
  Git, but no cost and no latency. No candidate writes one row with tokens, latency, wall-clock
  duration and money, with a "reported / computed / missing" mark per field, for every task.
- A deterministic choice of harness and model. ordewell, Galley, sortie and Coder Eval take a manual
  field per task or role; scion has tier aliases. None chooses from recorded cost and outcome.
- F-0003#66 counted. turbo-flow refuses a reviewer of the same model family and no_human uses
  another model, but no candidate guarantees and counts a verification by a different harness per
  change.
- F-0003#64 measured. chock, gh-aw, CC Safety Net and loopgate block or unstage writes to protected
  paths, but none ships a set of known-bad commits or an audit count of agent writes to rule paths.
- F-0003#59 at 100 %. agent-context-store and AgentPlane validate some handoffs; none covers every
  role transition. Go's standard library cannot validate a JSON Schema (`encoding/json` offers only
  `DisallowUnknownFields`), and no candidate gives a stdlib way to do it.
- Anything named "Armature": no public source.
- A Go-standard-library implementation. Every candidate is a whole program or a pattern; none can be
  vendored into a stdlib-only controller.

## Limits

- The GitHub API answered 403 (rate limit) for most checks. Last activity then comes from the
  repository's commit or release Atom feed or the ungh.cc mirror, not from `pushed_at`; languages
  come from READMEs when the language bar did not load ("not stated" in the table).
- Repository checks were made through web fetches by helper agents that read only the candidates'
  own pages; a few facts come from mirrors or release-note pages, as noted in the table.
- I ran nothing: no harness, no probe, no candidate. Flags and JSON fields are as the docs state on
  2026-10-09, and they change often (Claude Code releases daily; Gemini CLI headless fields are
  only shown by examples).
- OpenCode's `run --format json` (`step_finish` tokens and cost) is known only from third-party pages
  (query 7); I did not reach its own documentation, so it has no row.
- Rejected entries were judged by the list or search-file text alone; their repositories were not
  opened. Human UIs, plain harnesses and API gateways were excluded by the rule above.
- The "Awesome AI Coding Sandboxes" list and other lists named inside the five lists were not read.
- Price data was not compared with provider invoices.
- The file is about 700 lines, over the 300–450 asked: the rule "one line per entry that touches
  the problem" gives about 500 entry lines, even after the grouped "no line" sentences.
- Web finds checked and rejected, so not in the table: Agent Trace (RFC v0.1.0, CC BY 4.0; its repo
  answered 404; no token or cost fields), pi-per-commit-spend (state in `~/.pi`, 1 star, no
  trailers), Ingot (SQLite, joins by week, 0 stars), Ruler, crossby, block/ai-rules (same as
  Rulesync; not checked further), A2A (an HTTP protocol, not a file handoff), claude-code-router and
  OpenCodex (provider proxies), the Go wrappers claude-go, agent-sdk-go, claudecli-go (Go modules,
  so not dependencies), Modelglass (proprietary data terms).

~~~~
