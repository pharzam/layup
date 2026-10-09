# T-ywk7 — the search of M2b, searcher A, word for word

Searcher: GPT-6 Sol (`gpt-6-sol-high`) on Devin 3000.11.3 (`devin -p --permission-mode dangerous`), a fresh session with a `HOME` of its own, 2026-10-09 06:08:16Z to 06:13:14Z (4 min 58 s). Brief: [`brief.md`](brief.md), with the five lists of [`sources.md`](sources.md) and [`gh-search.tsv`](gh-search.tsv); nothing of the LAYUP architecture, the specification or the survey of `T-w73g`. The session ended with Devin's error "Your weekly usage quota has been exhausted" (exit code 1) after it had written this file (370 lines, each section of the brief present). The output follows word for word inside a fenced block; its paths are relative to the searcher's directory.

~~~~text
# Public-solution search — M2b / T-ywk7 — searcher A

2026-10-09. Scope: a Go-standard-library controller, Git as project record; registered coding-agent harnesses, pre-work probes, branch-bound sessions, harness-reported telemetry, protected project rules/gates and validated role handoffs. This is a search, not an implementation proposal.

Evidence convention: **verified** means a linked candidate's own README, documentation, or GitHub metadata explicitly says it; **list/search lead** is only a discovery claim. “Not stated” does not mean impossible. License/language/last activity below use repository metadata and default-branch latest-commit API on 2026-10-09; latest commit date is not evidence of a release. No candidate was run.

## Queries

| Web query (verbatim) | Brief words motivating it | What it found |
| --- | --- | --- |
| `coding agent CLI headless JSONL session token usage latency cost` | F-0001#38; F-0003#50/#60 | Claude CLI headless docs, Pi RPC, Headless CLI normalized usage. |
| `multi coding agent harness controller capability probe version model selection git worktree` | F-0001#9; M2b “probe session” | Orchestrate live profiles, worktree runners; many merely recommend probing. |
| `Go agent orchestrator coding agents worktrees git session telemetry standard library` | F-0004#1; M2b task branch | Go daemons using SQLite; no evidenced standard-library-only, Git-only complete match. |
| `coding agent protected rules gate paths prevent writes git CI` | F-0001#3; F-0003#64 | Chock's CI/hooks, tool-specific path guards; distinction between pre-write and commit rejection. |
| `agent role handoff JSON schema git repository coding agents` | F-0003#45/#59 | Agent Context Store and schema-backed role records. |
| `harness neutral AGENTS.md CLAUDE.md rules sync coding agents` | F-0003#26–#28/#52 | AGENTS.md and agentsync canonical rules plus generated copies. |
| `coding agents independent cross model verification deterministic gates` | F-0001#6; F-0003#66 | Independent-review tools; deterministic evidence differs from model verdict. |
| `agent client protocol initialize capability negotiation protocol version` | “first check … version, and what it can do” | ACP v1 initialize specifies version/capability negotiation. |
| `coding agent per task token latency wall clock duration money retries telemetry repo` | F-0003#20–#22/#60 | Benchmark JSONL and cost lenses; task-level completeness mostly unproven. |
| `Claude Code Codex Gemini OpenCode headless usage JSON output cost CLI` | M2b probe and developer session | Official headless formats; Headless CLI adapter, with explicit billing caveats. |

Also read **all** five supplied list READMEs at their supplied snapshots and **all** 200 rows of `sources/gh-search.tsv`. The 20 repository-search queries are represented by the source/fact and query columns in that file; web queries above are this searcher's separate searches. Duplicated hits are identified below, not independent confirmations.

## Candidates

“Git?” means documented project-repository state, **not** merely that the candidate's source code is hosted on GitHub. “Deterministic?” refers to the relevant check/control, not whether the agent's answers are deterministic. All 2026 dates are default-branch latest commits except Claude's undated documentation. A candidate classified as a tool is a separate executable, **not** a permissible Go module dependency.

| Kept candidate and primary evidence | License; language; last commit | Relevant behavior; project Git state?; deterministic? | Keep as / boundary |
| --- | --- | --- | --- |
| [ACP](https://github.com/agentclientprotocol/agent-client-protocol); [v1 initialization](https://agentclientprotocol.com/protocol/v1/initialization) | Apache-2.0; Rust schema; 2026-10-09 | JSON-RPC `initialize` negotiates protocol version, agent info/version and optional capabilities before sessions; Git state **no** (wire protocol); negotiation **yes**. | **Format/pattern** for actual capability checks where ACP is supported. It is editor↔agent, not a universal CLI probe, task telemetry store, or Git gate. Stable protocol v1 despite separate schema-release versions. |
| [Codex CLI](https://github.com/openai/codex); [noninteractive docs](https://developers.openai.com/codex/noninteractive.md) | Apache-2.0; Rust; 2026-10-09 | `codex exec --json` emits JSONL; `turn.completed.usage` contains token counts; `--output-schema` requests structured final output; sandbox and resume documented. Git state **not stated for controller artifacts**; JSON event decoding **yes**, model response **no**. | **Tool/adapter pattern**; no documented universal latency/money field or pre-work capability registry. `--output-schema` is not a repository handoff validator. |
| [Gemini CLI](https://github.com/google-gemini/gemini-cli); [headless spec](https://raw.githubusercontent.com/google-gemini/gemini-cli/main/docs/cli/headless.md) | Apache-2.0; TypeScript; 2026-10-08 | `-p --output-format json` has `stats` token usage/API latency; stream JSONL has `init` (model/session), `result` (aggregate/per-model usage), errors/exit codes. Git project state **not stated**; event parse **yes**. | **Tool/adapter pattern**; API latency is not necessarily wall-clock task time; harness-reported money not evidenced here. |
| [OpenCode](https://github.com/anomalyco/opencode); [CLI docs](https://opencode.ai/docs/cli/) | MIT; TypeScript; 2026-10-08 | `opencode run` supports programmatic invocation; `--model`, resume/session, stats/export and ACP are documented. Git project state **not stated**; output parsing **yes**. | **Tool/adapter pattern**; verify exact version's JSON event/usage schema before claiming complete telemetry. |
| [Claude Code headless](https://code.claude.com/docs/en/headless); [usage semantics](https://code.claude.com/docs/en/agent-sdk/cost-tracking) | License **not stated in these docs**; CLI language **not stated**; docs last-update **not stated** | `claude -p --output-format json/stream-json` reports session/cost; SDK documents per-step tokens, whole-tree cost and resumed cumulative totals. Git artifact storage **not stated**; field aggregation **yes**. | **Tool/adapter pattern**. `total_cost_usd` is an **estimate**, not billing; resumed totals and repeated message IDs can double-count; crashed sessions may lack result. `--bare` skips rules/hooks unless explicitly supplied. |
| [Headless CLI](https://github.com/RobertTLange/headless-cli); [usage guide](https://github.com/RobertTLange/headless-cli/blob/main/docs/usage.md) | Apache-2.0; TypeScript; 2026-09-12 | Normalizes CLI/model/cwd/sessions/output for seven agents and ACP; `--usage` appends tokens/cost provenance and attempt records; `--check` validates local setup. Project Git state **not stated**; adapter/usage normalization **yes**. | **Pattern/separate tool**. Missing usage stays missing and cost null; estimated API price is not subscription spend. Automatic paid fallback is documented; no local dollar cap. Node.js 22+ prevents in-process Go dependency. |
| [Orchestrate](https://github.com/bestagentkits/orchestrate) | MIT; GitHub primary language HTML (Agent Skills/README); 2026-09-19 | Claims live runtime/model/flag reprobe, capability/risk routing, isolated worktrees, bounded output, deterministic pre-checks, independent arbiter. `.orchestrate/benchmarks.json` in repo; run `state.json` location relative to project **not stated**. Hard filters **yes**; arbiter **no**. | **Pattern**, not a verified Go library. README admits R0/R1 calibration exception without independent arbiter; no evidenced mandatory complete per-task Git telemetry row or protected-rule path isolation. |
| [Galley](https://github.com/shinpr/galley) | MIT; Go; 2026-09-15 | Local daemon: task YAML, executor CLI/model/effort, separate supervisor, retry budget, worktree/task branch, checks and per-attempt run evidence under local workflow root. Project Git state **partial** (task branches/PR; run-root commits not stated); checks **yes**, supervisor verdict **no**. | **Pattern/separate tool**, closest Go lifecycle analogue. README also describes provider keys and optional Python/gh tooling; standard-library-only dependency and complete harness-reported tokens/latency/duration/money not established. |
| [Superharness](https://github.com/artificemachine/superharness) | GitHub license **NOASSERTION** (README claims Apache-2.0); Python; 2026-10-01 | Adapter registry for Claude/Codex/Gemini/OpenCode/Pi; queue/SQLite, idle+absolute watchdog, worktree fan-out, benchmark JSONL cost/duration/outcome. Project Git state **no as sole record** (SQLite); watchdog/health **yes**, model review **no**. | **Pattern/separate tool**; license discrepancy requires checking LICENSE. README's large test-count claim was not independently reproduced; token+latency+duration+money in one committed task row not shown. |
| [Conduit](https://github.com/conduit-cli/conduit) | MIT; Rust; 2026-03-19 | Four-agent worktree/checkout UI, sessions, provider/model selection, raw events and token/cost status. Local session state; Git workspaces **yes**, committed ledger **not stated**; status aggregation **yes**. | **Pattern/separate tool**; README warns cost uses configurable Claude pricing defaults and is an approximation; no complete telemetry or probe/gates. |
| [aGiTrack](https://github.com/core-aix/agitrack) | Apache-2.0; Python; 2026-10-08 | Claude/Codex/OpenCode sessions, per-turn auto-commits recording model and token cost, own worktree. Git **yes**; commit creation **yes**, code judgment **no**. | **Pattern/separate tool** for trace-to-Git. README says Windows agent write sandbox absent and only warns if base repo edited; not evidence of 0 writes to rule paths or task latency/wall-clock coverage. |
| [Agent Context Store](https://github.com/max9159/agent-context-store) | MIT; TypeScript; 2026-07-06 | `acs validate`, role artifacts/confirmed handoffs, in-repo `.acs/` default or local/dedicated store option; Claude/Codex/Cursor skills. Project Git **yes only in in-repo mode**; validation **yes**. | **Format/pattern/separate tool** for handoff. README says schema validation; exact percent coverage of all transitions not demonstrated; no harness session supervisor or usage row. |
| [AI-SDLC framework](https://github.com/ai-sdlc-framework/ai-sdlc); [draft agent spec](https://raw.githubusercontent.com/ai-sdlc-framework/ai-sdlc/main/spec/agents.md) | Apache-2.0; TypeScript; 2026-10-07 | **Draft normative specification** says versioned JSON-Schema handoffs, audit, `blockedPaths`, deterministic constraints and optional A2A discovery. Project Git storage **not stated**; prescribed validation **yes if implemented**. | **Format/pattern** only: specification MUST ≠ verified runtime enforcement, not a standard-library Go JSON Schema validator. |
| [Chock](https://github.com/open-coder-ai/chock) | Apache-2.0; Python; 2026-10-09 | Canonical policy compiled to native hooks, git hooks and CI; baseline rejects policy weakening; deterministic scripts/evals. Project Git policy **yes**; checks **yes**. | **Pattern/separate tool**. README explicitly calls agent plugins best-effort/fail-open and notes hooks not cloned; commit/CI rejection ≠ preventing all filesystem writes by an agent. |
| [AGENTS.md](https://github.com/agentsmd/agents.md) | MIT; TypeScript website; 2026-09-10 | Open project-local instruction format for coding agents. Git **yes when committed**; instructions **no** (not enforced). | **Format** for neutral source of rules, not protected paths, tests, or a handoff schema. |
| [agentsync](https://github.com/PanisHandsome/ai-rules-sync) | MIT; JavaScript; 2026-06-03 | Canonical `AGENTS.md`→Claude/Gemini/Cursor/etc.; `sync --check` detects drift and reports lossy translations. Git **yes when files committed**; drift check **yes**. | **Pattern/separate tool**; `--auto` explicitly permits reverse edits to any rule copy, incompatible with a strict protected canonical source unless disabled. Sync is not rule integrity. |

No candidate above is approved as a Go module dependency: F-0004#1 forbids even a third-party Go module. “Keep” means investigate the documented pattern/format or use an external executable if permitted, **not** adopt it unchanged. Primary descriptions do not prove deployment performance, security, completeness or license compatibility of transitive components.

## The lists

Each line below is a distinct entry **touching session management, routing, rule portability/protection, verification, handoff, Git state, or telemetry**; generic coding agents are represented when they are actual proposed harnesses, not every model/UI/editor from the large directories. List descriptions alone are not implementation verification. Repeated list entries are noted as duplicates. Kept names are the primary-checked candidates in the table; other promising leads are rejected **for this report** pending primary verification, not asserted unsuitable forever.

### `bradAGI_awesome-cli-coding-agents.md`

- OpenCode: keep: documented CLI harness/model selection.
- Codex CLI: keep: documented headless JSONL and usage.
- Gemini CLI: keep: documented headless structured stats.
- Claude Code: keep: documented headless cost and session behavior.
- Pi: reject: RPC lead; current repository/licensing/activity not checked.
- OpenHands: reject: broader agent platform; no checked task-ledger evidence.
- Oh My OpenAgent: reject: source-available multi-harness claim unverified here.
- Every Code: reject: Codex fork; no checked controller telemetry.
- Letta Code: reject: memory-first agent, not controller record.
- Tau: reject: JSONL-agent lead; no checked Git task ledger.
- Maki: reject: headless agent claim; narrower than controller.
- hax: reject: one-shot harness; not cross-agent controller.
- Waveloom: reject: Go harness is not Go controller; external modules unverified.
- San (two identical entries): reject: Go harness, not verified as controller.
- openHarness: reject: alternate harness; no verified shared task record.
- Symphony: reject: isolated-run orchestration lead, primary not checked.
- Omnigent: reject: multiple agent orchestration claim, primary not checked.
- zeroshot: reject: verification-loop claim, primary not checked.
- h5i: reject: list describes orchestration but [current repository](https://api.github.com/repos/h5i-dev/h5i) describes web-security workspace; mismatched lead.
- oh-my-agent: reject: multi-agent workflow claim, primary not checked.
- Aeon: reject: cross-agent routing claim, primary not checked.
- Agent Executor (AX): reject: cluster-scale runtime, not verified small Go/Git solution.
- GitHub Agentic Workflows: reject: GitHub-hosted workflow, not verified local controller.
- AgentBox: reject: sandbox/worktree partial lead, primary not checked.
- Parallel Code: reject: worktree/session partial lead, primary not checked.
- agent-deck: reject: session UI partial lead, primary not checked.
- agent-manager: reject: session management lead, primary not checked.
- amux: reject: multiplexing lead, primary not checked.
- AgentsMesh: reject: multi-harness lead, primary not checked.
- Traycer: reject: task routing lead, primary not checked.
- Bernstein: reject: quality-gate lead, primary not checked.
- fractal: reject: orchestration lead, primary not checked.
- Kiro Crew: reject: specialized team workflow, primary not checked.
- aGiTrack: keep: verified Git commit/token linkage.
- AIWG: reject: rule-deployment lead, primary not checked.
- Archon: reject: deterministic workflow claim, primary not checked.
- Claudexor: reject: quota-aware routing claim, primary not checked.
- numbat: reject: orchestration lead, primary not checked.
- AgentSight: reject: observability lead, primary not checked.
- HOL Guard: reject: gate integrity lead, primary not checked.
- SandBase Harness: reject: sandbox lead, primary not checked.
- Data Olympus: reject: context governance lead, primary not checked.
- Mneme: reject: context persistence lead, primary not checked.
- AgentPack: reject: portable state claim, primary not checked.
- grite: reject: Git-native task state lead, primary not checked.

### `andyrewlee_awesome-agent-orchestrators.md`

- sortie: reject: Go/SQLite ticket runner; not verified Git-only or stdlib-only.
- SuperPlane: reject: CI/review orchestration claim, primary not checked.
- symphony: reject: isolated agent runs claim, primary not checked.
- Tale: reject: sandbox and review claim, primary not checked.
- Taskuary: reject: local multi-CLI queue claim, primary not checked.
- Agent Messaging Protocol: reject: signed transport, no verified schema-valued project handoffs.
- agenttier: reject: Kubernetes sandbox infra, not local Git controller.
- aGiTrack: keep: verified worktree/turn commits and token linkage.
- AIWG: reject: project-owned rules claim, primary not checked.
- Archon: reject: test/worktree lead, primary not checked.
- AX: reject: Kubernetes workload runtime, not demonstrated Go-stdlib fit.
- Claudexor: reject: cross-harness review/quota lead, primary not checked.

### `hesreallyhim_awesome-claude-code.md`

- Claude Code / official docs: keep: actual headless and cost semantics checked.
- Claude Code GitHub Action: reject: hosted CI slice; not full local lifecycle.
- Claude Code Security Review: reject: model review alone is not deterministic gate.
- Agent Skills: reject: skills format is not validated role handoff.
- ccusage: reject: historical usage reports, no checked complete Git task rows.
- ccstatusline: reject: UI token display, no persisted task ledger.
- cc-probeline: reject: status/probe display, not verified harness capability probe.
- CCDash: reject: usage dashboard, not Git system of record.
- Rulesync: reject: rule-copy portability lead, primary not checked.
- Container Use: reject: isolation slice, no checked full telemetry.
- Dippy: reject: permissions slice, no checked path-write guarantee.
- GouvernAI: reject: governance lead, primary not checked.
- TDD Guard: reject: test-discipline gate, not whole controller.
- agents-md-cookbook: reject: instruction examples, not enforced protected rules.

### `tensorchord_Awesome-LLMOps.md`

- Cordum: reject: general policy/orchestration platform, coding-harness integration not checked.
- brood-box: reject: microVM isolation slice, Git/telemetry not checked.
- DOS Kernel: reject: evidence-gate lead, primary not checked.
- Helicone (duplicated list entry): reject: API observability, not checked per coding-task Git row.
- Traceloop OpenLLMetry: reject: OTel tracing slice, not Git controller.
- Langfuse (duplicated list entry): reject: external observability state, not Git-only.
- Lookspan: reject: local trace UI, not verified Git task records.
- traceAI: reject: LLM spans alone not session lifecycle.
- Token Police: reject: SDK/API-call budget, no verified CLI/subscription coverage.
- witness: reject: proxy/journal pattern, primary not checked.
- AgentField: reject: general agent control plane, not small local Git solution.
- Agnos: reject: LLM gateway, not harness/session controller.
- AIWG: reject: portable rules lead, primary not checked.
- Arize-Phoenix: reject: broad ML observability, no verified Git task ledger.
- BitRouter: reject: model gateway ≠ harness selection/session control.
- LiteLLM: reject: model API adapter, not coding-harness CLI adapter.
- OpenLIT: reject: token/cost traces, not project Git record.
- Kitaru: reject: durable workflow lead, primary not checked.

### `RoggeOhta_awesome-codex-cli.md`

- Codex CLI / official docs: keep: headless JSONL and structured output verified.
- mco: reject: neutral orchestration lead, primary not checked.
- external-subagents: reject: external process/sandbox slice, not checked ledger.
- shinpr/sub-agents-skills: reject: cross-model delegation skill, not controller.
- Untrivial-ai/agent-orchestrator: reject: daemon with external state; not checked Git-only.
- codex-cli-hooks: reject: hook/cost claims unverified, hooks can be bypassed.
- agent-cost-mcp: reject: live cost UI, not checked complete Git task ledger.
- openai-codex-mcp: reject: MCP bridge, not end-to-end controller.
- SandBase CLI: reject: model/API discovery, not harness probe verification.
- CodexMonitor: reject: desktop sessions, no verified Git ledger.
- agent-cli-farm: reject: tmux/session manager, telemetry not checked.
- ccNexus: reject: provider gateway, not coding agent task ledger.
- Mysti: reject: cross-agent routing claim, primary not checked.
- agent-peer-review: reject: cross-agent review claim, primary not checked.
- codex-claude-bridge: reject: collaboration link, no verified independent gate.
- claude_codex_bridge: reject: shared session UI, not verified independent gate.
- metaswarm: reject: TDD/team orchestration lead, primary not checked.
- maestro-orchestrate: reject: specialist workflow lead, primary not checked.
- maestro-flow: reject: multi-CLI workflow lead, primary not checked.
- waybar-ai-usage: reject: status widget, not repository row.
- ai-sessions-mcp: reject: usage history service, not checked task coverage.
- ccusage: reject: post-hoc usage reports, not guaranteed task join.

## The search file

`source/query` groups are shown in their order in `sources/gh-search.tsv`; each line below is **one relevant repository entry**. Duplicate hits in multiple queries remain visible. Reject is a relevance/evidence decision for this search, not a technical impossibility. Search descriptions, stars and pushed timestamps are **not** primary evidence. Entries about unrelated dashboards, CRM, prediction/ML research, browser automation or provider serving that do not touch this controller are omitted per brief.

### Cost visibility / token usage / ledger (rows 2–41)

- junhoyeo/tokscale: reject: usage tracker; not checked task/latency/Git record.
- getagentseal/codeburn: reject: multi-CLI usage dashboard, not checked task ledger.
- DeepAgentLabs/agenticlens: reject: profiler lead, no checked Git completeness.
- splunk/token-meter: reject: local usage dashboard, not verified role/task association.
- Nihondo/AgentLimits: reject: quota widget, not task telemetry.
- Backtthefuture/TokenStep: reject: desktop token widget, not task record.
- JingbiaoMei/Tokdash: reject: dashboard, no checked Git task row.
- tokentopapp/tokentop: reject: usage display, not committed task record.
- xiufengsun/TokenTracker: reject: multi-tool usage UI; task latency not checked.
- stormzhang/token-tracker: reject: status-line token tracker, not complete ledger.
- ramtinJ95/opencode-tokenscope: reject: one-harness session analytics.
- mag123c/toktrack: reject: cost tracker, no checked task correlation.
- getagentseal/codeburn (second hit): reject: duplicate usage-dashboard lead.
- Han-1413141/dsh-cost-meter: reject: DeepSeek-specific spend/quota view.
- juliantanx/aiusage: reject: usage sessions, no checked full task row.
- oluwajubelo1/otellix: reject: Go LLM spans, not coding-agent Git sessions.
- inferock/inferock-bench: reject: provider proxy, not CLI task controller.
- fjgbue/claude-delegator-deepseek-mcp: reject: model delegation, not task ledger.
- voly-codes/voly: reject: per-task cost claim; primary not checked.
- deeflect/smart-spawn: reject: model router, not harness sessions.
- adididitagain/tokentab: reject: spend/quality dashboard, no checked Git row.
- OptimNow/cost-per-task: reject: proxy-cost estimator, not CLI native reports.
- BrianWong05/TokenLedger: reject: local-log counts, latency/Git not checked.
- ethanplusai/agent-ledger: reject: SQLite workspace, not Git-only.
- bisheshabramhacharya/token-ledger: reject: menu-bar usage, not task ledger.
- hollis-labs/go-usage-ledger: reject: ledger lead; coding CLI integration unverified.
- wakeup595626-cmyk/agent-token-ledger: reject: desktop scan, not per-task gate.

### Claude/Codex usage (rows 42–61)

- Maciek-roboblog/Claude-Code-Usage-Monitor: reject: quota monitoring, not repo task row.
- steipete/CodexBar: reject: usage menu bar, not repo state.
- phuryn/claude-usage: reject: one-harness session dashboard.
- jarrodwatts/claude-hud: reject: in-session UI, not Git row.
- CodeZeno/Claude-Code-Usage-Monitor: reject: quota widget.
- realiti4/claude-swap: reject: account switching rather than task accounting.
- tddworks/ClaudeBar: reject: quota UI.
- vinzdg/codenotch: reject: usage limits UI.
- shanggqm/codexU: reject: quota widget, not task telemetry.
- junhoyeo/tokscale (second hit): reject: duplicate tracker.
- vibe-cafe/vibe-usage: reject: token CLI, full record not checked.
- douglasmonsky/codex-usage-tracker: reject: one-harness history; Git row not checked.
- getagentseal/codeburn (third hit): reject: duplicate usage-dashboard lead.
- ClaudeCodeUsage/ClaudeCodeUsage: reject: editor usage estimates, not task repository.
- xiufengsun/TokenTracker (second hit): reject: duplicate tracker.
- jaykinhoo9/codex-usage-badge: reject: quota/indicator, not ledger.
- Javis603/token-monitor: reject: desktop tokens/cost, not verified task join.

### Neutral rules / synchronizers (rows 62–81)

- toroleapinc/claude-brain: reject: Claude memory sync, not neutral protected rules.
- PanisHandsome/ai-rules-sync: keep: verified canonical AGENTS.md converter/drift check.
- Common-ka/ai-agent-unity-rules: reject: Unity-specific rule copy.
- sampleXbro/agentsmesh: reject: multi-tool config lead, primary not checked.
- Chemaclass/agnostic-ai: reject: sync claim, primary not checked.
- dhruv-anand-aintech/agent-rules-sync: reject: sync claim, primary not checked.
- airulefy/Airulefy: reject: sync claim, primary not checked.
- martinmose/agentlink: reject: sync claim, primary not checked.
- saqibameen/agent-dotfiles: reject: dotfile sync claim, primary not checked.
- BayramAnnakov/claude-reflect: reject: memory learning may mutate rules.
- iannuttall/source-agents: reject: file-sync lead, primary not checked.
- Signet-AI/signetai: reject: shared identity/memory service, Git protection unverified.
- PanisHandsome/ai-rules-sync (second hit): keep: same verified rules converter.
- intellectronica/claude-agentsmd: reject: one-way Claude shim, scope narrower.
- Daliagents-com/agentsge-1: reject: sync/memory lead, primary not checked.
- westonplatter/aps: reject: instruction composition, no checked enforcement.
- shriramkv/agentsmd-sync: reject: AGENTS propagation lead, primary not checked.
- bensyverson/agents: reject: cross-repo sync lead, no checked gate.

### Orchestration and independent review (rows 82–111)

- OrchestratorInc/agent-orchestrator: reject: multi-harness platform; Git-only state not verified.
- getpaseo/paseo: reject: multi-agent desktop/mobile controller; ledger unverified.
- HKUDS/DeepCode: reject: standalone agent harness rather than neutral CLI supervisor.
- 21st-dev/1code: reject: orchestrator lead; primary not checked.
- mattpocock/sandcastle: reject: sandboxed TypeScript runner, Git row unverified.
- bradAGI/awesome-cli-coding-agents: reject: discovery list, not a solution.
- Danau5tin/multi-agent-coding-system: reject: model-team orchestration, task ledger unverified.
- superset-sh/superset: reject: IDE/worktree sessions; rule-gate evidence unverified.
- awslabs/cli-agent-orchestrator: reject: tmux sessions; Git-only task state unverified.
- SeemSeam/claude_codex_bridge: reject: collaboration UI, no checked independent gate.
- wanshuiyin/Auto-claude-code-research-in-sleep: reject: ML research, not code gate.
- sno-ai/sno-station: reject: shared agents/memory; gate proof unverified.
- aAAaqwq/AGI-Super-Team: reject: multi-agent org, no checked telemetry.
- marcuspat/turbo-flow: reject: cross-model review/gate claim, primary not checked.
- majiayu000/harness: reject: Rust governance controller claim, primary not checked.
- gmickel/flow-next: reject: cross-model workflow lead, primary not checked.
- EvilFreelancer/crossreview: reject: independent reviews only, gate/ledger unchecked.

### Agent runners, protected paths, gates (rows 112–141)

- Design-Arena/agent-runner: reject: model-agnostic loop, not CLI harness orchestration.
- darkzOGx/darkzloop: reject: model runner, not verified task record.
- TheArchitectit/AIGGP-Agentic-Framework: reject: quality-gate lead, primary not checked.
- getrunkite/runkite: reject: Postgres/Redis control plane, not Git-only.
- theycallme-eric/agent-runner: reject: verified coding controller lead, primary not checked.
- panchew/local-agent-runner: reject: Ollama loop, not cross-harness CLIs.
- schuligan/agentic-test-runner: reject: tracker/tests platform, no checked telemetry.
- suchithnarayan/agent-security-hooks: reject: harness hooks, not enforced cross-harness filesystem protection.
- Serhioromano/pi-defender: reject: Pi-specific deny hooks, not universal boundary.
- ckckck/agent-safe-delete: reject: reversible deletion only.
- obstalabs/bulwark: reject: read gate, not protected-rule write gate.
- jmstar85/oh-my-githubcopilot: reject: Copilot-only automation/gates.
- kenryu42/cc-safety-net: reject: pre-exec guard, not verified all write surfaces.
- JeongJaeSoon/agent-guard: reject: secret-leak guard, not rule-path guarantee.
- paudley/coding-ethos: reject: policy/CI lead, primary not checked.
- moerasermax/AgentCharter: reject: charter/hooks claim, primary not checked.
- Offsend/Offsend: reject: secrets/context protection, not rule-path proof.
- open-coder-ai/chock: keep: checked policy compiler, CI/hooks and fail-open caveat.
- bram-wq/agent-guardrails: reject: deterministic hooks lead, primary not checked.
- Quarkgluonmixture/coding-agent-guardrails: reject: rule/hook lead, primary not checked.

### Handoff and Git state (rows 142–161)

- dabit3/agent-handoff: reject: schema claim, primary not checked.
- openenvelope/schema: reject: versioned workflow schema claim, primary not checked.
- ZibbyDev/agent-workflow: reject: schema handoff claim, primary not checked.
- max9159/agent-context-store: keep: in-repo artifacts and `acs validate` documented.
- agentsentinel/agentsentinel-handoff: reject: handoff schema lead, primary not checked.
- faresrafat3/agent-handoff: reject: typed handoff lead, primary not checked.
- moona3k/handoff: reject: short-URL handoff is not Git-only.
- talocode/handofflane: reject: validator lead, primary not checked.
- roymcfarland/agent_handoff_guide: reject: playbook, not validated artifact.
- shin4141/decision-os-v12-completion-integrity: reject: schema/CI lead, primary not checked.
- GDWN-BLDR/stateweave: reject: agent memory state, not project task state.
- yunaremaia/agent-workspace: reject: worktree manager, telemetry unverified.
- Lians-ai/Lians: reject: evidence-based verification lead, primary not checked.
- hoeggsoftware/multi-agent-planning: reject: Git shared artifacts lead, primary not checked.
- minhtran3124/herdr-worktree-agents: reject: worktree sidebar, no checked ledger.
- ychamel/RepoResident: reject: Git project-memory harness, no checked cost row.
- saint0x/tl: reject: checkpoint primitive, not telemetry/runner.

### Deterministic workflow, headless, task worktrees and final row (rows 162–201)

- vekexasia/pi-extensible-workflows: reject: Pi-only deterministic flow.
- a5c-ai/babysitter: reject: gate orchestration lead, primary not checked.
- luckeyfaraday/athena-loops: reject: Python review loop, no checked Git ledger.
- Extra-Chill/homeboy: reject: deterministic orchestration claim, primary not checked.
- omar-os/omar: reject: orchestration specification claim, primary not checked.
- NTCoding/autonomous-claude-agent-team: reject: Claude-only hook workflow.
- mjasnikovs/pi-task: reject: Pi-only prompt workflow.
- Ali-hey-0/ai-runtime-lab: reject: generic FSM examples, not coding CLI task.
- amirfish1/claude-command-center: reject: multi-CLI sessions; primary not checked.
- OhadAssulin/headless-coder-sdk: reject: multi-CLI SDK lead; primary not checked.
- 2233admin/obsidian-llm-wiki: reject: wiki roles, not task controller.
- cp-yu/cc-switch-web: reject: provider/config UI, not branch/telemetry.
- Snowflake-Labs/subagent-cortex-code: reject: Snowflake delegation, not task sessions.
- nkhdiscovery/headless-agentic-codebase: reject: container workflow lead, primary not checked.
- ClipboardHealth/groundcrew: reject: per-task sandbox/worktree lead, primary not checked.
- DnzzL/herdr-automations: reject: scheduled fresh worktree, task ledger unchecked.
- errhythm/copse: reject: worktree CLI, no checked harness runner.
- Maples7/VibeChard: reject: Apple per-task worktree UI, no checked telemetry.
- vraj00222/agent-farm: reject: Claude-only worktrees, no cross-agent check.
- mutheejj/worktree-orchestrator: reject: worktree/model routing claim, primary not checked.
- FabianoArthur/claude-code-kitchen: reject: Claude-only gates/task tmux.
- ex3del/mast: reject: Claude plugin/worktrees, no verified ledger.
- DrSeedon/orchestra: reject: task branch/team lead, primary not checked.
- 2-ring/agentic-tasktrees: reject: VS Code worktree UI, no telemetry.
- VasiHemanth/tokentelemetry: reject: telemetry dashboard, Git task join unverified.
- mikehasa/agentacct: reject: per-step cost/time lead, primary not checked.
- genai-telemetry/genai-telemetry: reject: SDK traces to external backend.
- PxA-Labs/AgentsScope: reject: agent traces, no verified Git task record.
- SahilSelokar/Tokensense: reject: in-process model router, not harness controller.
- onekapisch/Tokens-4-Breakfast: reject: quota/spend monitor, not task ledger.

## What no candidate covers

- **Complete combined demonstration:** no verified candidate both probes **every registered** harness, then runs a developer on a task branch and commits one row containing **harness-reported** token count, latency, task wall-clock duration **and money** in the project repository.
- **Go boundary:** Galley is Go, but standard-library-only was not established; ACP, ACS, Chock and Headless CLI are protocols/other-language tools, not permitted in-process dependencies. Go's `encoding/json` does not itself validate arbitrary JSON Schema.
- **Probe semantics:** ACP has negotiated capabilities but only for ACP-speaking agents; CLI `--version` and a successful test prompt are not proof of standardized telemetry, model support, sandbox behavior or protected-write enforcement across non-ACP CLIs.
- **Latency vs duration:** Gemini's API latency and Codex/Claude result timing are not the controller's end-to-end wall-clock task duration. The latter must be measured separately; truncation/timeouts require explicit missing-data status, not fabricated numbers.
- **Money provenance:** Claude's `total_cost_usd` is estimated; Conduit's cost is approximate; Headless CLI distinguishes native and price-list estimates and leaves unavailable cost null. Subscription spend is not equivalent to list-price token arithmetic.
- **Retry attribution:** resumptions can include earlier cumulative totals; per-step duplicate IDs and subagent tokens can double-count or undercount; a timeout/crash can have no final result. None proves every retry maps once to a Git-committed role/task ledger.
- **Rule integrity:** AGENTS.md and sync utilities express/replicate instructions, not deny writes. A git hook catches commits, not the rule-path write itself, and an agent with shell access may alter its local hook/config. Chock's native hooks explicitly fail open. CI protected by permissions and an out-of-agent filesystem boundary are distinct from a prompt prohibition; no kept project demonstrates the 0-write invariant end to end.
- **Validated transitions:** ACS documents validation, while AI-SDLC specifies contracts; neither demonstrates 100% of all role transitions are schema-validated under this Go controller. Codex structured generation still needs trusted validation.
- **Independent review:** Galley has a separately selectable supervisor and Orchestrate an arbiter for higher risk, but README text does not prove reviewer harness differs on every change; model review alone does not replace deterministic gate checks.
- **State location:** worktree branches and Git history do not imply queue, telemetry and decisions are committed. Superharness uses SQLite; Conduit local persistence; Galley's workflow-root evidence is not documented as necessarily committed.

## Limits

- Only the named kept candidates received primary-source/metadata checks. All other list and search-file entries are discovery leads; “reject” frequently means evidence or scope insufficient **here**, not that their implementation has been disproven.
- Metadata's latest commit is from each repository's **default branch** as observed 2026-10-09, not the installed executable's version or a verified release. Claude headless docs did not state a code license, language or dated release. Superharness README badge calls itself Apache-2.0 while GitHub repository metadata reports NOASSERTION; the LICENSE text was not independently inspected.
- Public README/spec claims (including test counts, security guarantees and routing behavior) were not reproduced by source audit or execution; candidate code was not run or installed as the brief requires.
- Large curated lists include duplicated entries and broad LLMOps tools. Items with no material overlap with the specified controller (model training/serving, CRM, unrelated research, display-only novelty) are omitted. Some list descriptions can drift from current repositories (notably h5i).
- A web search is not an exhaustive registry; authentication, paid provider settings, platform differences and changing CLI JSON schemas were not exercised. Money billed by subscription cannot be established from a reported token-cost estimate.
- The source lists are snapshots, `sources/gh-search.tsv` is a ranked search sample rather than a census, and neither popularity nor a `pushed` timestamp confirms feature coverage. The verified primary URLs above supply the basis for kept rows; the remaining judgments are explicitly provisional.
~~~~
