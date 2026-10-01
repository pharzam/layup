# T-hbw8 — second search, searcher B, word for word

Searcher: Grok 4.7 (`grok-4-7-high`) on Devin, a fresh session, 2026-09-28 10:37:39Z to 10:50:20Z (12 min 41 s). Brief: [`../search-v2/sources.md`](sources.md) (the six lists in full and [`gh-search.tsv`](gh-search.tsv)), the problem statement, the idea owner's answers and the vision brief — nothing of the LAYUP architecture. The session ended with Devin's error "Your daily usage quota has been exhausted" (exit code 1) after it had written this file (415 lines, every section of the brief present); its last message said it was checking names in the reject lists. The output follows word for word inside a fenced block (its links are relative to the searcher's directory); each U+2039 and U+203A character is written by name.

~~~~text
# Public-solution search (searcher B)

Date of the activity data: 2026-09-28. GitHub API `pushed_at` and `stargazers_count` where shown. A feature is `not stated` when the source that this file quotes does not say it. This search does not design LAYUP.

The twelve In-Scope items of PSB section 6, short names:

1. Problem Statement Quality
2. Reproducible Discipline Setup (Armature)
3. Rule Protection
4. Stack-Dependent Gates
5. Role Handoffs
6. Autonomous Clarification
7. Verification on Every Change
8. Human-on-the-Loop
9. Stall Resolution
10. Cost Visibility
11. Specification Synthesis
12. Harness-Agent Neutrality

## 1. Queries

Each query uses words from the PSB or the vision brief. The candidates are the public items that the result pages named.

| # | Query | Source text | Candidates |
| --- | --- | --- | --- |
| 1 | human as message bus multi-agent software development clarifying questions | PSB Problem 1, "Humans are forced to serve as manual, synchronous message buses" | No product. Papers and posts only: BusMA (arxiv 2609.15054); "Beyond Messag... |
| 2 | multi-role agent team software delivery Product Owner architect QA Engineer autonomous | PSB section 2 role list | MetaGPT; ChatDev; DevTeam AI (blog, no kept repo); Deliver-Team-Agent (GitH... |
| 3 | agent orchestration budget approval gates token cost governance coding agents | PSB Problem 5 and Human Decision Point 4 | Paperclip; agentgateway cost-control post (Solo.io); GitHub "What is AI Age... |
| 4 | stalled agents deadlock resolution independent diagnosis fresh context | PSB Problem 4 and vision section 3.2 | No product that gives a fresh-context diagnosis. Closest names: Gas Town st... |
| 5 | derive requirements and acceptance criteria from problem statement software agents traceability | PSB Problem 6 | No agent product in the top results. Classical traceability pages only. BMA... |
| 6 | harness-neutral agent rules AGENTS.md Claude Code Codex same repository governance | PSB Problem 7 and Invariant 9 | GraylumAI PR 444; sisa-messaging-rs issue 41; agent-harness (furedea); harn... |
| 7 | detect gaps in problem statement before development starts agents batch questions | PSB Problem 3 and Decision Point 2 | LAYUP issue 37 (`layup psb check`) is the named implementation of this item... |
| 8 | git as system of record multi-agent software delivery discipline gates | Invariant 1 | GNAP; Gas Town beads; agent-delivery-gates; a NXTG.AI note on git as the on... |
| 9 | BMAD method multi-agent software development specification gates | Vision section 2.1, specification from a brief | BMAD-METHOD (bmad-code-org) |
| 10 | spec-kitty spec-driven development Claude Codex requirements trace | Problem 6, harness neutrality | Spec Kitty |
| 11 | gastown beads git-backed agent orchestration mayor polecat | Vision section 3.4, issues as the substrate | Gas Town (gastownhall/gastown) |
| 12 | gnap git-native agent protocol no server coordination | Invariant 1 and Invariant 2 | GNAP (farol-team/gnap) |
## 2. The pass over the lists and the search file

Every entry of the six lists and every row of gh-search.tsv was read. A keep is one row. A reject that touches the problem is named in the class line for its source. The count is the entries that do not touch the problem.

Reject class codes (one line each; the name list is every touching entry that is not a keep):

- UI: a session, board, or dashboard. No discipline record.
- prompts: role text or a skill pack. No machine gate.
- one-harness: one product only. Fails Invariant 9.
- library: a framework or SDK. No project system of record.
- state-outside-git: state in a server, SQLite, or a ledger outside the project repository. Fails Invariant 1.
- cost-screen: token or cost on a screen or in a proxy. Not a per-task record in the project repository.
- thin-description: the search row or the list line does not state a gate, a schema, or a Git record.
- stale: last push is months before 2026-09, or the list marks it resting.
- infra: a sandbox, a router, or a transport. Not a delivery discipline system.
- partial-no-discipline: it touches the problem (gates, approvals, handoffs, or telemetry) but does not state a schema, a PSB path, Armature, or a Git record.

### andyrewlee awesome-agent-orchestrators

List size: 248 product entries (307 lines). Not in contact with the problem: 188.

| Name | Link | Keep | Reason |
| --- | --- | --- | --- |
| AgentBridge | https://github.com/raysonmeng/agent-bridge | keep | Claude Code and Codex as live peers for mid-turn review; a component, not a... |
| crewAI | https://github.com/crewAIInc/crewAI | keep | Role framework. See the shortlist |
| gastown | https://github.com/gastownhall/gastown | keep | See the shortlist (Gas Town) |
| kodo | https://github.com/ikamensh/kodo | keep | A separate agent verifies each result; a pattern for independent verification |
| loki-mode | https://github.com/asklokesh/loki-mode | keep | See the shortlist |
| paperclip | https://github.com/paperclipai/paperclip | keep | See the shortlist |
| no_human | https://github.com/no-human-ai/no_human | keep | A second model in a fresh session reviews the diff and must refute "done";... |
| Concord MCP | https://github.com/Get-Concord-AI/concord-mcp | keep | Claim, collision detect, durable messages, review handoff across harnesses;... |
| foremerge | https://github.com/naw103/foremerge | keep | Intent before write, deterministic collision rules, acceptance on a verific... |
| Open Multi-Agent | https://github.com/open-multi-agent/open-multi-agent | keep | DAG, approvals, traces, evaluation, checkpoints; a pattern. State is not st... |
| gnap | https://github.com/farol-team/gnap | keep | See the shortlist. The list marks it resting (last commit 2026-03) |

Rejects that touch the problem:

- UI (13): agent-session-manager, agent-orchestrator, AI4Kanban, automaker, Garcon, octomux, Paseo, agent-kanban, openswarm, ralphex, sortie, OpenWorker, vibe-kanban
- cost-screen (5): Cyclops, termany, Fusion, fractal, aGiTrack
- partial-no-discipline (41): YYLO, AGX, Fletch, intentic, ivy-tendril, kandev, 5dive, Agon, ClawTeam, CompanyHelm, NXTG-Forge Orchestrator, OpenRig, orc, ORCH, Orkas, ruflo, shire, tutti, bernstein, Dex, LongHorizon-Harness, LoopTroop, MartinLoop, ordewell, ralph-orchestrator, aeon, Contrabass, gh-aw, NEEDLE, OpenHands, agent-runbook, Agentlas OS, AIWG, Archon, Claudexor, Crewplane, LionClaw, skillfold, Hivekeep, OpenMausBot, swarm-protocol
- one-harness (3): Aperant, corellis, squad
- state-outside-git (7): buzz, OtoDock, scion, future-os, symphony, guild, omnigent
- infra (3): centaur, Agent Messaging Protocol, AX

### kyrolabs awesome-agents

List size: 196 entries (224 lines). Not in contact with the problem: 155.

| Name | Link | Keep | Reason |
| --- | --- | --- | --- |
| MetaGPT | https://github.com/FoundationAgents/MetaGPT | keep | See the shortlist |
| CrewAI | https://github.com/crewAIInc/crewAI | keep | Same product as the crewAI row above |
| Cordum | https://github.com/cordum-io/cordum | keep | Out-of-process policy, approval gates, signed audit; a governance component... |
| zeroshot | https://github.com/the-open-engine/zeroshot | keep | Planner, implementer, independent validators in isolation; a verification p... |
| h5i | https://github.com/h5i-dev/h5i | keep | Peer review, then a neutral verifier replays and tests; a verification pattern |
| Paperclip | https://github.com/paperclipai/paperclip | keep | Same as paperclip above |
| Ouroboros (Q00) | https://github.com/Q00/ouroboros | keep | A Socratic interview gates the spec on an ambiguity score, then a budgeted... |
| Orbi | https://github.com/orbi-build/orbi | keep | Issue to implementation to an independent review of the frozen head; a veri... |

Rejects that touch the problem:

- partial-no-discipline (25): OpenClaw, AG2, Swarm, Swarms Framework, AgentField, hcom, Hive, Aeon, Agent Swarm, Tenuo, GPT Pilot, Aider, Devika, TypedAI, OpenCode, ReviewCerberus, Bernstein, Orca (echoVic), OpenCodeReview, CompozyOS, Ordewell, 5dive, SAGE, Caura, Agno
- infra (1): e2b
- library (6): Autogen, AgentVerse, AgentScope, agency-swarm, PraisonAI, llama-agents
- one-harness (3): Maestro, Plandex, SWE Agent
- state-outside-git (4): Astron, Lobu, Latitude, ctop
- UI (5): SwarmClaw, SandBase Harness, Manifest, ClawMetry, Busabase
- cost-screen (1): Maestro Orchestrate

### vijaythecoder awesome-claude-agents

List size: 24 agent files plus the repo itself. Not in contact with the problem: 21. The repo is one harness (Claude Code) by its own install text: `claude /agents`.

Rejects that touch the problem:

- prompts (2): awesome-claude-agents, Team Configurator
- UI (1): Tech Lead Orchestrator

### caramaschiHG awesome-ai-agents-2026

List size: 340+ named resources (672 lines). Categories with no contact (creative, voice, CRM, sales, local model runners, vector databases, healthcare, papers, newsletters, market stats): 248 entries, count only. Coding, frameworks, multi-agent, protocols, observability, and governance: the rows below. Entries already in an earlier table are not repeated unless this list adds a fact.
Not in contact inside the remaining sections (browser, voice, creative, CRM, research tools, local runners, models, papers, books, newsletters): 248.

| Name | Link | Keep | Reason |
| --- | --- | --- | --- |
| GNAP | https://github.com/farol-team/gnap | keep | "coordinate AI agent teams with 4 JSON files in a git repo. No server, no d... |

Rejects that touch the problem:

- state-outside-git (5): Kiro, Devin, AXME, Temporal, Langfuse
- partial-no-discipline (10): Caliber, Qodo, CodeRabbit, PR-Agent, Miyabi, DeerFlow, Portia AI, Guardrails AI, NeMo Guardrails, Credo AI, Microsoft Agent Governance Toolkit
- one-harness (2): Copilot Workspace, Cursor AI Automated Team
- library (2): LangGraph, Semantic Kernel
- UI (5): MagiC, Mission Control, Dify, Langflow, TeamHero
- infra (2): MCP, A2A
- cost-screen (1): Helicone

### bradAGI awesome-cli-coding-agents

List size: 312 named entries (884 lines). Not in contact with the problem: 246. A terminal agent, a sandbox, a memory store, a voice frontend, or a token compressor is not a row.

| Name | Link | Keep | Reason |
| --- | --- | --- | --- |
| gastown | https://github.com/steveyegge/gastown | keep | Same product family as gastownhall/gastown |
| zeroshot | https://github.com/the-open-engine/zeroshot | keep | Same as zeroshot above |
| Loki Mode | https://github.com/asklokesh/loki-mode | keep | Same as loki-mode above |
| h5i | https://github.com/h5i-dev/h5i | keep | Same as h5i above |
| Orbi | https://github.com/orbi-build/orbi | keep | Same as Orbi above |
| kodo | https://github.com/ikamensh/kodo | keep | Same as kodo above |
| Vibestrate | https://github.com/guyshonshon/vibestrate | keep | Fresh reviewer process, approval gates, local token and cost ledger; a patt... |
| Beads | https://github.com/gastownhall/beads | keep | Git-backed issue graph for agents; a component of Gas Town, not a disciplin... |
| Concord MCP | https://github.com/Get-Concord-AI/concord-mcp | keep | Same as Concord above |
| AgentPlane | https://github.com/basilisk-labs/agentplane | keep | Task, plan, approve, implement, verify, finish. "All state stays in `.agent... |
| Mneme | https://github.com/MnemeHQ/mneme | keep | ADRs become guardrails; enforces at the earliest boundary. A pattern for ru... |
| Hivelore | https://github.com/Doucs91/hivelore | keep | A lesson becomes a regex or AST guard that Git hooks refuse; a pattern for... |
| isitdone | https://github.com/raimondasl/isitdone | keep | Blocks "done" until test, typecheck, and lint pass; scans for weakened test... |
| Squelette | https://github.com/JyMinet/squelette | keep | Human decisions and scope recorded before work; pre-commit refuses out-of-s... |

Rejects that touch the problem:

- partial-no-discipline (32): Oh My OpenAgent, CLI Agent Orchestrator, hcom, ADHDev, Agent AFK, Podiom, TaskHandoff, postmortemthis, Symphony, Omnigent, Bernstein, MartinLoop, ORCH, OMK, Forge, fab-kit, Galley, Dahrk, Relay, the-perfect-orchestrator, gh-aw, HOL Guard, Wit, AgentDiff, Data Olympus, grite, TraceFold, SpecWave, PatchWarden, Project Tiny Context Harness, BlaBla, llm-panel
- UI (6): Multica, Agent Teams AI, Better Agent, CLITrigger, Even, context-bridge
- one-harness (3): defract, Kiro Crew, Team1-Factory
- prompts (5): oh-my-agent, great_cto, OpenCastle, TeDDy, Agentic Engineering Framework
- cost-screen (2): deepsec, tu
- state-outside-git (3): Okto Nexus, mycroft-class myc, Weaver

### e2b-dev awesome-ai-agents

List size: 215 open-source projects plus 80 closed products (5591 lines). Not in contact with the problem: 269. The list is a 2023-2024 catalog. The software entries that touch the problem:

| Name | Link | Keep | Reason |
| --- | --- | --- | --- |
| ChatDev | https://github.com/OpenBMB/ChatDev | keep | "CEO, CPO, CTO, programmer, reviewer, tester, and art designer" in seminars... |
| CrewAI | https://github.com/joaomdmoura/crewai | keep | Same product. "Each agent is akin to a team member, possessing specific rol... |
| MetaGPT | https://github.com/geekan/MetaGPT | keep | Same product. "given one line requirement, returns PRD, Design, Tasks, or R... |

Rejects that touch the problem:

- partial-no-discipline (8): AutoGen, AutoGPT, Devon, GPT Pilot, L2MAC, OpenDevin, Sweep, Devin
- library (1): CAMEL
- state-outside-git (1): DevGPT
- one-harness (3): Devika, SWE Agent, GitHub Copilot X
- UI (2): DevOpsGPT, Factory
- prompts (1): GPT Engineer

### gh-search.tsv

233 data rows. Not in contact with the problem: 168. A row is in contact when the description names roles, governance, budgets, handoffs, verification, stalls, or specification. Stars and last push are the file values.
Rows in the TSV that are awesome-lists, token-compressors, single coding agents, or tutorials, and that the table does not name: 168.

| Name | Link | Keep | Reason |
| --- | --- | --- | --- |
| open-multi-agent |  | keep | Same product. "Self-hosted TypeScript agent runtime with durable ap... |
| council-of-high-intelligence |  | keep | Councils, triads, debates across products; a pattern for a panel, n... |
| spec-kitty |  | keep | See the shortlist |
| MetaGPT |  | keep | Same product. Last push on the API is 2026-01-21 |
| eshwarvijay/secondmate |  | keep | "Your agent writes the code, a different model tries to break it."... |

Rejects that touch the problem:

- infra (4): hoangsonww/AI-Agents-Orchestrator, othmaneachirorionech/token-governance-protocol, sprix-sage-router, Trendifai/agentstate
- stale (2): Tarekivida/ai_development_team, CodeMachine-CLI
- partial-no-discipline (43): iFeel-is-a-mouse/openclaw-teamdev-flow, NeoDesign/NEOs-Paperclip-Agent-Templates, z4gunn/springer, Edcwsyh/microco, rupeshpoojary9/multi-agent-orchestrator, sshnaidm/spenda, synaptixs/spine, debasisdwivedy/Coderama, jonathanpopham/estate-agent, feedback-loop-ai/brokkr, marius-patrik/DarkFactory, the-vibey-project/vibey, ZeeshanAmjad0495/Forge-Loop, oh-my-claudecode, paseo, DeepCode, agent-orchestrator, omnigent, nexent, golutra, Bindu, munder-difflin, mission-control (builderz-labs), AgentsMesh, agency-orchestrator, ruflo, oh-my-claudecode (topic row), adk-python, camel, EvoAgentX, oinsio/gnomish-factory, Bhavya-Dhoot/Cohort, moai-adk, open-design and archify and awesome lists, deer-flow, CowAgent, wshobson/agents, agentscope, PraisonAI, astron-agent, ChaoYue0307/awesome-graph-engineering, luoyuejun9/dsh-workstate, caya8205-2/crewctl
- one-harness (5): Tekyz-Inc/get-stuff-done-teams, cmAIdx/headless-claude-automation-template, armelhbobdad/bmad-module-ultracode-goal, Mysti, yukinoshi/lazyskill
- prompts (6): jiyoon99/ai-development-team-playbook, pengpengliu1212-art/lobsterai-team, q3ok/coordinated-agent-team, ulpi-io/autonomous-engineering, rafmacalaba/armada, nbzhaosq/sdlc-team
- UI (12): brian-mwirigi/CostHQ, seanyao/roll, AntonioR199/Agentic-Scrum, Free-devloper/autonomous-software-delivery-organization, yao, kungfu, planning-with-files, eigent, harness-sdk, claude_codex_bridge, edict, AlphaMine-Tech/agent-os
- cost-screen (3): LeelaKrishna-R/agent-budget, Rumblingb/agent-cost-tracker-mcp, AlvaroGB/tokenwatch
- library (10): kosminus/reasonflow, mitchellsjohnson/agent0-pdlc, adk-go, swarms, solace-agent-mesh, koog, langroid, nanobot, microsoft/agent-framework, DaosPath/handoffkit
- state-outside-git (1): filipemmacedo/OctoGent
- thin-description (6): eduardocerqueira/ai-alpha-squad, launchapp-dev/ao, certified84/autonomous-sdd, Jiffy-Agnet/gateway, Nasiko-Labs/nasiko, cengebretson/orc

## 3. The shortlist

Eight candidates. Marks: `covers`, `partly`, `no`. Evidence is a quoted source. Activity is the GitHub API on 2026-09-28. Item numbers are the twelve In-Scope items in section 1 of this file. A mark of `no` with no extra note means the quoted source does not state the item.

### 3.1 Paperclip

- Link: https://github.com/paperclipai/paperclip
- License: MIT (API `license.spdx_id`).
- Language and stack: TypeScript. "Paperclip is a Node.js server and React UI."
- Infrastructure: a Node.js server and a dashboard. The default database is not stated in the README sentences that this search fetched. The GNAP comparison table says Paperclip uses PostgreSQL; that claim is GNAP's table, not Paperclip's README.
- Activity: 91437 stars. Last push 2026-09-28.

| Item | Mark | Evidence |
| --- | --- | --- |
| 1, 2, 4, 6, 9 | no | Not stated. No Armature. No fresh-context diagnosis. |
| 3 | partly | "Governance: who can do what." Project rule files are not stated. |
| 5 | partly | "Hierarchies, roles, reporting lines." Armature artifacts are not stated. |
| 7 | partly | "verify from diffs, screenshots & tests." A human verifies. Deterministic gates are not stated. |
| 8 | partly | "Approve hires, override strategy, pause or terminate any agent." Approvals are broad, not only the five Decision Points. |
| 10 Cost Visibility | partly | "Monthly budgets per agent. When they hit the limit, they stop." "Cost tracking surfaces token budgets." Latency and wall-clock duration in the project repository are not stated. |
| 11 Specification Synthesis | no | "Every task traces back to the organization mission." A trace to problem-statement text with an acceptance criterion is not stated. |
| 12 Harness-Agent Neutrality | partly | "Bring your own agents" and the logo row names OpenClaw, Claude Code, Codex, Cursor, Bash, HTTP. Rules of the target project are not stated as one neutral form. |

Conflicts with the invariants: work state lives in the Paperclip server ("sessions persist across reboots"), not only in the target project repository (Invariant 1). A human cannot continue from the target repository alone if the server is gone (Invariant 2). Armature pin is not stated (Invariant 8). The control plane is a product the project depends on (Invariant 9, for LAYUP's own rules).

Use: a pattern for budgets, org charts, and approvals. Not a base. The server is the system of record.

### 3.2 BMAD-METHOD

- Link: https://github.com/bmad-code-org/BMAD-METHOD
- License: README says "MIT License — see LICENSE." The API `license.spdx_id` is NOASSERTION. Treat the README as the author's claim and the API as not confirmed.
- Language and stack: Python (API `language`). Install path uses Node.js, npm, and uv. Skills install into the coding tool.
- Infrastructure: no server. Skills in the project and in the coding tool.
- Activity: 53576 stars. Last push 2026-09-28.

| Item | Mark | Evidence |
| --- | --- | --- |
| 1 | partly | "agents and workflows make the important decisions explicit." A batch gap check of a prose PSB is not stated. |
| 2, 3, 4, 6, 9, 10 | no | Not stated. No Armature. No rule-path ban. The BMad Loop module is a different repo. |
| 5 | partly | "Bring in product, architecture, UX, development, and testing expertise." A machine schema is not stated. |
| 7 | partly | "reviewed implementation, correction, and learning." Layout and contract gates are not stated. |
| 8 | partly | "without handing over judgment." The five Decision Points are not stated. |
| 11 | partly | "carry its briefs, specifications, and architecture." A trace to the source sentence is not stated. |
| 12 | partly | Install routes for a skills CLI, a Claude Code plugin, and a Codex plugin. One neutral rule form is not stated. |

Conflicts: durable context is "Carry product and technical decisions forward instead of re-explaining them in every chat." The store of that context is not stated as Git only (Invariant 1, not stated). Agents can install and update skills (`npx skills update`), so Invariant 3 is not shown. No pinned Armature (Invariant 8).

Use: a pattern for specification and role skills. Not a base. It does not set the discipline baseline.

### 3.3 MetaGPT

- Link: https://github.com/FoundationAgents/MetaGPT
- License: MIT (API and README badge).
- Language and stack: Python. Needs Node and pnpm. Config in `~/.metagpt/config2.yaml`.
- Infrastructure: local CLI. Output is a repo in `./workspace`. No server is required for the CLI path. A hosted product MGX is a different thing.
- Activity: 70662 stars. Last push 2026-01-21. No push in eight months.

| Item | Mark | Evidence |
| --- | --- | --- |
| 1 | no | Input is "a one line requirement." A gap batch is not stated. |
| 2, 3, 4, 6, 7, 8, 9, 10, 12 | no | Not stated. One model config. No Decision Points. |
| 5 | partly | "product managers / architects / project managers / engineers" and "carefully orchestrated SOPs." |
| 11 | partly | "outputs user stories / competitive analysis / requirements / data structures / APIs / documents." A trace to the source sentence is not stated. |

Conflicts: config and secrets live in `~/.metagpt/`, outside the project repository (Invariant 1). The generated repo can stand alone after the run (Invariant 2 is closer). No rule protection, no Armature pin.

Use: a pattern for SOP roles and document output. Not a base. Stale push.

### 3.4 Spec Kitty

- Link: https://github.com/spec-kitty/spec-kitty
- License: MIT (API and README).
- Language and stack: Python 3.11+. CLI `spec-kitty`.
- Infrastructure: "Spec Kitty is local-first and stores its core artifacts in your repo." Optional hosted tracker. Optional local dashboard.
- Activity: 1647 stars. Last push 2026-09-28. The README says `main` carries a 4.x release candidate and "stable launch acceptance remains pending."

| Item | Mark | Evidence |
| --- | --- | --- |
| 1, 2, 3, 4, 6, 10 | no | Not stated. Start is `specify`, not a PSB gap check. No Armature. No rule-path ban. |
| 5 | partly | "Work packages with lifecycle lanes." A machine schema is not stated. |
| 7 | partly | "Review, accept, merge, and retrospective gates." Layout and contract gates are not stated. |
| 8 | partly | "governed and human-in-loop." "reviewers accept, reject, or merge with an audit trail." The human still accepts each merge. The five Decision Points are not named. Not `covers`. |
| 9 | partly | spec-driven.md says the lane workflow "prevents work from stalling." A fresh-context diagnosis is not stated. |
| 11 | partly | "acceptance criteria" live with the spec. A trace to a problem-statement sentence is not stated. |
| 12 | partly | Skills for Claude Code, Codex, Cursor, Gemini, Copilot, Windsurf, OpenCode. One gate run under two harnesses is not stated. |

Invariant conflicts: core artifacts are in the repo (Invariant 1, close). "the repository remains the source of truth" (Invariant 2, close). Hosted sync is opt-in. Agents are not shown as unable to edit the workflow files (Invariant 3, not stated). No Armature (Invariant 8).

Use: the closest public workflow for spec, plan, tasks, review, and Git worktrees. A component pattern. Not a base: no Armature, no cost record, no stall diagnosis, no rule protection.

### 3.5 Gas Town

- Link: https://github.com/gastownhall/gastown
- License: MIT (API).
- Language and stack: Go. Needs Git, Beads (`bd`), Dolt, tmux, and a coding CLI.
- Infrastructure: a town directory `~/gt`, git worktrees, a Dolt-backed Beads ledger, tmux sessions. Not a single project repository.
- Activity: 18206 stars. Last push 2026-09-18.

| Item | Mark | Evidence |
| --- | --- | --- |
| 1, 2, 3, 4, 6, 10, 11 | no | Not stated. A scheduler limits concurrency. A per-task token record in the project repo is not stated. |
| 5 | partly | "Built-in mailboxes, identities, and handoffs." A machine schema is not stated. |
| 7 | partly | "runs verification gates, and merges to main using a Bors-style bisecting queue." Layout and contract gates are not stated. |
| 8 | partly | Escalation routes "through the Deacon, Mayor, and (if needed) Overseer." A business-forking rule is not stated. |
| 9 | partly | "autonomous stall detection and smart skip logic." "detects stuck agents, triggers recovery." A recorded fresh-context diagnosis is not stated. |
| 12 | partly | "Claude Code, GitHub Copilot, Codex, Gemini, and others." The town and Dolt ledger are extra systems. |

Conflicts: work state is in the Beads ledger and the town, not only in the target project repository (Invariant 1). Dolt is a database (Invariant 2: a human cannot continue from the project repo alone). No Armature pin.

Use: a pattern for git-backed work tracking, stuck-agent recovery, and a merge queue. Not a base.

### 3.6 GNAP

- Link: https://github.com/farol-team/gnap
- License: MIT (API and badge).
- Language and stack: a protocol. "Four JSON files. That's the entire protocol." No runtime language.
- Infrastructure: Git only. "No servers. No databases."
- Activity: 86 stars. Last push 2026-03-17. The orchestrator list marks it resting. The protocol is a draft.

| Item | Mark | Evidence |
| --- | --- | --- |
| 1, 2, 4, 6, 7, 11 | no | Not stated. |
| 3 | no | The README does not give a rule-path ban. |
| 5 | partly | Tasks, runs, and messages are JSON with required fields. Armature conventions are not stated. |
| 8 | partly | "humans and AI agents are both first-class participants." Decision Points are not stated. |
| 9 | no | `blocked` and `blocked_reason` exist. A fresh-context diagnosis is not stated. |
| 10 | partly | A run has `tokens` and `cost_usd`. "cost tracking (budget = sum of runs)." Latency and duration are not stated. The record is in Git. |
| 12 | covers | "Any agent, any runtime — if it can git push, it can participate." OpenClaw, Codex, Claude Code, or custom. The cover is the task protocol, not the gates. |

Conflicts: Git is the store (Invariant 1 and Invariant 2, close). Budgets and governance are "application layer — not part of the protocol," so the protocol itself does not enforce them (Invariant 3, not met). No Armature (Invariant 8). Stale.

Use: a pattern for harness-neutral task and cost records in Git. A component idea. Not a base. Inactive since March 2026.

### 3.7 Loki Mode

- Link: https://github.com/asklokesh/loki-mode
- License: BUSL-1.1 (README badge). API `spdx_id` is NOASSERTION. Source-available, not OSI-open. Commercial editions exist under Autonomi.
- Language and stack: Shell (API). CLI also ships as npm and Docker. Providers include Claude, aider, cline.
- Infrastructure: local CLI and optional Docker. Writes a Git repo with "source, tests, configs, and audit logs."
- Activity: 1072 stars. Last push 2026-09-28.

| Item | Mark | Evidence |
| --- | --- | --- |
| 1 | no | Input is a PRD, a GitHub issue, OpenAPI, or a one-line brief. A gap batch is not stated. |
| 2, 3, 4, 6 | no | Not stated. Gates are not stated as stack-selected. |
| 5 | partly | "41 role definitions." Machine-validated handoff artifacts are not stated. |
| 7 | partly | "8 quality gates" and "Code is not done until it passes automated verification." |
| 8 | partly | "shows the real cost and time estimate before spending anything." Ongoing Decision Points are not stated. |
| 9 | partly | "parallel review (blind council)." A recorded stall diagnosis is not stated. |
| 10 | partly | The sample receipt shows a cost line. Latency and duration per task are not stated. |
| 11 | partly | It accepts a spec and does not mark done on an empty diff. A trace to the problem statement is not stated. |
| 12 | partly | Several providers. One neutral rule form is not stated. |

Conflicts: the license is not MIT-class (a project constraint, not one of the nine invariants). Audit logs are in the output repo (Invariant 1, partly). Agents are not shown as barred from the gate files (Invariant 3, not stated). No Armature.

Use: a pattern for a completion gate and a blind review council. Not a base. License is BUSL-1.1.

### 3.8 AgentPlane (plus the verification pair)

- Link: https://github.com/basilisk-labs/agentplane
- License: MIT (list text).
- Language and stack: not stated in the list line.
- Infrastructure: "All state stays in `.agentplane/` inside the repo; no hosted runtime."
- Activity: 78 stars in the list text. API not fetched. Last push: not stated.

| Item | Mark | Evidence |
| --- | --- | --- |
| 1 | no | Not stated. |
| 2 | no | Not stated. |
| 3 | partly | The workflow is "task, plan, approve, implement, verify, finish." A ban on writes to rule paths is not stated. |
| 4 | no | Not stated. |
| 5 | partly | State is in the repo. A machine schema for role artifacts is not stated. |
| 6 | no | Not stated. |
| 7 | partly | A verify step is in the workflow. Deterministic layout gates are not stated. |
| 8 | partly | An approve step is in the workflow. The five Decision Points are not stated. |
| 9 | no | Not stated. |
| 10 | no | Not stated. |
| 11 | no | Not stated. |
| 12 | partly | "wraps Claude Code, Codex, Cursor, and Aider." One neutral rule file is not stated. |

The verification pair that this item does not replace:

- no_human (https://github.com/no-human-ai/no_human): "a second model in a fresh session reviews each diff and is told to refute done ... merging stays human."
- isitdone (https://github.com/raimondasl/isitdone): "blocks a coding agent's done until the repo's real test, typecheck and lint commands pass" and "scans the diff for weakened tests." "Zero LLM calls."

Conflicts: state in `.agentplane/` is inside the repo (Invariant 1, close). No Armature. Small and not fetched past the list line, so several marks stay `no` because the source does not say.

Use: a small Git-native workflow pattern, plus two verification components. Not a base.

## 4. What no candidate covers

No public solution in this pass covers all twelve In-Scope items. No candidate covers these items at all, on the sources this search quoted:

- Reproducible Discipline Setup of the Armature baseline, at a pinned version, with evidence for each configuration value (Invariants 4 and 8). No README names Armature.
- Stack-Dependent Gates that add layout, interface, contract, and test rules and that cannot weaken a baseline rule (Invariant 7).
- Rule Protection as a machine ban: the agents that do the work cannot write the rule paths, and a known-bad commit is detected (Invariant 3). Mneme, Hivelore, Squelette, and isitdone are patterns. None states this invariant.
- Stall Resolution as specified: a retry limit, an independent examination with a fresh context, and a diagnosis plus an outcome stored in the project repository. Gas Town detects stuck agents. Loki Mode has a blind council. Neither states the record.
- Cost Visibility as three fields per task in the project repository: token count, latency, and wall-clock duration. Paperclip and GNAP record cost or tokens. Latency and duration together are not stated.
- Specification Synthesis from an approved problem statement, with an identifier, a trace to the source text, and an acceptance criterion. BMAD, MetaGPT, and Spec Kitty produce specs. The trace to the problem-statement sentence is not stated.
- Problem Statement Quality as a check that emits one batch of gap questions before delivery. The only named implementation in the search results is LAYUP's own issue 37. Ouroboros (Q00) gates a spec on an ambiguity score. That is a pattern, not this check.
- Autonomous Clarification: a question that does not need a human gets an accepted answer from the responsible role, with no human. No README states this.
- Harness-Agent Neutrality as specified: the same rules and gates run under two harnesses, and the verifier is not the harness that made the change. GNAP is harness-neutral as a task protocol. It has no gates. Spec Kitty and Paperclip run many harnesses. They do not state one gate run under two harnesses.

Invariant 6 (a deterministic check is preferred to an LLM judgement) is a pattern in isitdone, bernstein, foremerge, and Hivelore. None of them is a full discipline system.

No candidate can be LAYUP's base. The usable public parts are patterns and components: Git as the task bus (GNAP, Beads), budgets and approvals (Paperclip), spec workflows (Spec Kitty, BMAD), a second-model review (no_human, kodo, h5i, secondmate), and a mechanical done-check (isitdone).
~~~~
