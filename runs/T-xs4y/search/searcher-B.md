# Public-solution search for M2e (task T-xs4y)

Date: 2026-10-15. Searcher: Sonnet 5.5 (one of two). Judged only by what a source says. No candidate code was run or installed. "Not stated" means the source page I fetched did not say. For many repositories the fetched page showed a commit count but no commit date and no release list; those cells say "not stated (N commits)".

Method note: I ran the web searches and fetches myself. The five curated lists were each read in full by a sub-agent that I briefed (the three big lists were read in chunks; no filter or cap). I did not re-read the lists line by line myself, so the list verdicts are the sub-agents' judgments, spot-checked only through the candidates I fetched. The search file I read in full myself.

## Queries

| Query | Words it comes from | What it found |
|---|---|---|
| GitHub commit status API required status checks branch protection rulesets GitHub App statuses permission | F-0003#47, F-0003#58, F-0001#5 | Old GHES docs, ruleset API pages; `checks[].app_id` / `integration_id` binds a required check to an App. Then fetched GitHub docs (see next rows). |
| fetch docs.github.com/en/rest/commits/statuses | F-0003#58 | Create status needs push access; states error/failure/pending/success; context is case-insensitive; limit of 1000 statuses per SHA and context. |
| fetch docs.github.com ... available-rules-for-rulesets | F-0003#58, F-0001#3, F-0003#43 | A required-status app source needs `statuses:write`; "any source" lets anyone with write set it; bypass actors may be Apps; "Restrict updates" limits push to bypass actors. Page does not say how a never-reported check is treated (blocks, by implication). |
| multi-agent coding orchestrator typed handoff schema validated transitions Git-backed state controller Go | F-0003#45, F-0003#59 | norma, Sortie, hawk/mission, Orca, Gas Town. |
| LLM agent handoff transition table validated by code JSON schema Go standard library only | F-0003#59, F-0004#1 | Only hawk/mission (Go, unclear); OpenAI Agents SDK handoffs (JSON-schema tool args; Python/TS). Nothing Go-stdlib. Note: Go stdlib has no JSON Schema validator, so schema checks would be hand-written. |
| requirements traceability as code markdown checker unresolved ID uncovered requirement acceptance criterion CI | F-0003#24, F-0003#51, F-0003#62 | spec-trace, cargo-ovft, jamb, TRLC, SpecTrace; Tessl doc-traceability skill. |
| GitHub ruleset restrict updates branch only GitHub App bypass actor store state in orphan branch | F-0001#1, F-0001#3, F-0003#43 | App bypass is whole-ruleset, not path-scoped; PR-based alternative; no source on orphan-branch state. |
| architecture decision options evaluation matrix automated filter evidence required CI ADR tool | the options-filter need | MADR/Nygard ADR templates, decision matrices; nothing that filters options by code on evidence. |
| fetches of repositories (rtmx, agent-context-store, norma, sortie, gastown, github-status-updater, RepoSteward, agent-delivery-gate, doorstop, OpenFastTrace, duvet, entireio/cli, ZibbyDev/agent-workflow, statewright, lobster, faresrafat3/agent-handoff, Orbi, h5i, AgentPlane, Galley, grite, gh-aw, Symphony, isitdone, Vibestrate, Mneme) | the kept candidates | Details in Candidates. |

Queries of the search file (`sources/gh-search.tsv`) are listed under "The search file".

## Candidates

Columns: name and URL; license; language; last activity; part of the problem; state in Git; deterministic; keep as and why.

| Name / URL | License | Language | Last activity | Part of problem | State in Git | Deterministic | Keep as, and why |
|---|---|---|---|---|---|---|---|
| GitHub REST commit statuses + rulesets (docs.github.com/en/rest/commits/statuses; .../managing-rulesets/available-rules-for-rulesets) | n/a (forge docs) | n/a | live docs | controller-posted statuses, required checks, who can update a branch | n/a | n/a | **API usage and known limits.** POST /repos/{o}/{r}/statuses/{sha}; states error/failure/pending/success; context per check class (`layup/gates`, etc.); max 1000 statuses per SHA+context; required-check source can be pinned to the App (needs `statuses:write`); with "any source" anyone with write can set the status, so pin the App. "Restrict updates" lets only bypass actors push, and an App bypass is for the whole ruleset, not path-scoped (community threads; not in GitHub's page). The ruleset itself needs admin to set, which the App lacks: a human sets it once. The docs page does not say how a never-reported check is treated; a missing required check blocks merge by implication, so test it. |
| norma https://github.com/normahq/norma | MIT | Go (1.26+) | not stated (732 commits, no releases) | task loop (Plan-Do-Check-Act), worktree per task | Task branch in Git; run state in SQLite, tasks in Beads | not stated | **Pattern.** One-task loop with a worktree and branch per task. Reject as a base: SQLite plus Beads hold state, and the check step is an agent verdict. |
| Orbi https://github.com/orbi-build/orbi | AGPL-3.0 or Sustainable Use License | Python 3.14+ (inferred) | v0.5.47 linked; no date | task loop: issue label, frozen base SHA, worktree, PR, independent review, "Only the reviewed head can merge" | Issues are the only state store ("no database, queue, or daemon") | not stated | **Pattern.** Closest published match to "merge the exact checked head": review bound to a commit SHA, state in GitHub issues. Not a dependency (Python, AGPL). Its state is in the forge, not Git. |
| Galley https://github.com/shinpr/galley | MIT | Go (inferred from `cmd/`, `go.mod`) | not stated (450 commits) | executor plus separate supervisor against acceptance criteria; draft PR; bounded retries | No: run evidence under a local `runs/<run-id>/`, queue config in `~/.galley` | not stated | **Pattern.** Different-model supervisor and bounded retry in Go. External `go.mod` deps not shown; local state, so not a base. |
| AgentPlane https://github.com/basilisk-labs/agentplane | MIT | TypeScript (Node 24+) | not stated (14,268 commits, pre-1.0) | task state, plan, verify, finish as in-repo records | Yes: `.agentplane/tasks/<id>/README.md`, `acr.json`, `WORKFLOW.md` | control plane yes ("does not make an LLM deterministic") | **Format.** Per-task machine record (`acr.json`) in the repo, lifecycle transitions that "semantic episodes cannot perform or claim". Page does not say whether agents can edit the workflow file. |
| grite https://github.com/neul-labs/grite | MIT | Rust | not stated (126 commits) | controller records in Git | Yes: append-only event log in `refs/grite/wal`, BLAKE2b-chained events, optional Ed25519 signatures | claims same result in any merge order; LWW tie-break not specified | **Pattern.** Hash-chained event records in a dedicated ref, synced by fetch/push. For the "controller-only branch" the project would use a branch with a ruleset, not CRDT merge. |
| entireio/cli https://github.com/entireio/cli | MIT | Go (inferred) | not stated (10,312 commits) | agent sessions recorded in Git | Yes: `refs/entire/checkpoints/<shard>/<id>`, trailer `Entire-Checkpoint:` on commits; never writes to the active branch | not stated | **Pattern.** Records kept in a private ref, linked from commits by trailer. Its secret redaction is "best-effort". |
| GitHub Agentic Workflows (gh-aw) https://github.com/github/gh-aw | MIT | Go/Node (inferred) | not stated | agent writes buffered and validated in separate scoped jobs | n/a | not stated | **Pattern.** Closest published statement of "what an agent writes never makes a handoff valid by itself": the agent job is read-only; validated "safe outputs" are applied by a separate job. |
| OpenAI Symphony https://github.com/openai/symphony | Apache-2.0 | spec; Elixir reference | not stated (47 commits) | tracker issue to isolated run to proof of work | not stated | not stated | **Pattern / spec to read.** I could not see SPEC.md's contents, so what it specifies is unchecked. |
| max9159/agent-context-store https://github.com/max9159/agent-context-store | MIT | TypeScript/Node (inferred; page empty) | not stated (107 commits) | typed handoffs by role (BA, SA, Dev, QA) | Yes, git-tracked store | not stated | **Format and pattern.** YAML handoff records validated against JSON schema; `validate` rejects handoffs that break role-to-role rules in workflow files. Closest to the handoff need; Node, so only a design to copy. Matches "a transition not in the table is invalid". |
| faresrafat3/agent-handoff https://github.com/faresrafat3/agent-handoff | MIT | Python 3, no deps | not stated (13 commits) | typed handoff | `.agent-workspace/` in project | `doctor --strict` exits non-zero on problems | **Pattern.** Handoff has 11 sections; the next action must name a file, command, ID or check; claims such as `git_head` and `tree_state` are checked against Git. Good idea for "an artifact has a place and a hash". Tiny project. |
| ZibbyDev/agent-workflow https://github.com/ZibbyDev/agent-workflow | MIT | JS/TS | not stated (202 commits) | schema-enforced handoffs in a graph | No: local session folders | not stated | **Pattern.** Each node declares a Zod `outputSchema`; edges are explicit. Not Go. |
| statewright https://github.com/statewright/statewright | Apache-2.0 plus FSL-1.1-ALv2 (gateway) | Rust | not stated (736 commits) | state machine guardrails per phase | session cache / PocketBase / cloud | not stated | **Known limit, pattern.** Tools exposed to the agent include deactivate and create-workflow, and the page does not say they are blocked, which breaks F-0001#3. Cursor mode is advisory. |
| isitdone https://github.com/raimondasl/isitdone | MIT | TypeScript (inferred) | not stated (61 commits) | verification gate | receipts, not stated where | "zero-LLM", exit codes | **Pattern.** PASS/FAIL/STALE receipt bound to a tree hash; an edit after PASS makes it STALE; HMAC-signed. Matches "merge the exact head that was checked". |
| agent-delivery-gate https://github.com/mojunrubest/agent-delivery-gate | not shown | TypeScript/Node (inferred) | not stated (31 commits, experimental) | host-owned verifier, agent cannot declare delivery | not stated | verifier runs a canonical command without a shell | **Pattern and threat model.** Has `THREAT_MODEL.md` and JSON schemas. Says it does not stop a malicious project attacking a same-host runner. Unclear whether it posts statuses itself. |
| Mneme https://github.com/MnemeHQ/mneme | MIT | Python 3.11+ | not stated (761 commits) | ADRs as deterministic guardrails | `.mneme/project_memory.json` | yes, "core allow/warn/fail decisions don't depend on an LLM judge" | **Pattern.** Rules become active only after explicit validate and activate. Nearest to an option filter, but it checks proposals against ADRs rather than filtering option rows. |
| cloudposse/github-status-updater https://github.com/cloudposse/github-status-updater | Apache-2.0 | Go | not stated (20 commits, no releases) | posts commit statuses and sets required checks | n/a | n/a | **API usage (reference).** Uses a user token with `repo:status`, so it is not App-shaped; likely `go-github` (not confirmed). Read for the request shape only. |
| rtmx https://github.com/rtmx-ai/rtmx | Apache-2.0 | Go (2 deps: Cobra, YAML) | not stated (373 commits); Python version end-of-life 2026-09-25 | requirements traceability matrix | Yes: `.rtmx/database.csv` | status derived from test results; `health` lints orphans, cycles, stale refs | **Tool or format.** Go, git-stored. Test-centred rather than problem-statement-text-centred; its two deps rule it out as a dependency, but the CSV matrix and the health checks are a model. |
| Doorstop https://github.com/doorstop-dev/doorstop | LGPLv3 | Python 3.10+ | not stated | requirements as YAML items in version control, linked | Yes: one YAML file per item | not stated | **Format.** One file per requirement with links; tree integrity check. No trace to prose text or acceptance criterion by itself. |
| OpenFastTrace https://github.com/itsallcode/openfasttrace | GPL-3.0 | Java 17+ | v4.10.0 in example; no date | requirement tracing, markdown items | in-repo files | not stated | **Format (read only).** The page I fetched did not say how it detects uncovered or broken links; I did not confirm. A search result says cargo-ovft reads this format; I did not open it. |
| Duvet https://github.com/awslabs/duvet | Apache-2.0 | Rust | not stated (126 commits) | traces code to quoted specification text | annotations in source | `report` for CI | **Method.** `//=` spec URL/section and `//#` quoted text: a trace to the stating text. Closest published method for F-0003#51 "trace to the text that states it". Source-comment based, not requirement-ID based. |
| LOBSTER https://github.com/bmw-software-engineering/lobster | AGPL-3.0 | Python | not stated (716 commits) | tracing policy and coverage report | files | not stated | **Format idea.** Interchange schema plus a policy file (`lobster.conf`) saying which levels must trace. AGPL, Python: not a dependency. |

I did not open spec-trace, cargo-ovft, jamb, TRLC or SpecTrace beyond one search result each, so they are not kept as candidates; a later pass might.

## The lists

Verdicts below come from sub-agents who read each list in full; I fetched only the ones marked (fetched). None of the lists shows license, language or date per entry (only stars or badges), so those were checked only for fetched entries.

### hesreallyhim/awesome-claude-code (848 lines)
- Orbi-like merge gate: not in this list; see bradAGI.
- Consort: keep (borderline): spec-first TDD loop under a deterministic orchestrator with human gates; unfetched.
- Agents Can Communicate (ACC): keep (borderline): saved handoffs and verification results across providers; unfetched.
- AI Research Skills: keep (borderline): stage deliverables are schemas, enforced by skills, not code.
- Fusion Harness: keep (borderline): architect plus builders on different models; unfetched.
- ralph-harness (vkuprin): keep (pattern): CLI owns the loop, fresh session per iteration.
- Gryph: keep (pattern): YAML policy blocks agent actions, with an audit log; hook-level.
- Node9: keep (borderline): deterministic governance proxy with audit log.
- SkilLock: keep (pattern): blocks unapproved drift of skills in CI; fits "agents cannot change rules".
- TDD Guard: keep (pattern): hook blocks rule-violating edits.
- Claude Code GitHub Action / Security Review (Anthropic): keep (reference): CI agent on issues and PRs; possible verifier input.
- claude-code-tools: keep (borderline): Claude Code to Codex handoff, but context only, not typed.
- Container Use, Airut: keep (borderline): isolated environments for dev sessions.
- Claude Code PM, AB Method, RIPER, Charrette, OctoShell, Agent 007, napkin-to-blueprint, Agent Collab Skills, Claude Code Harness, AgentSys, OpenRig, Plannotator, Claude CodePro, Superpowers, Compound Engineering, gstack, Harness (revfactory), HCOM, Repowire, Claude Squad, Multi-Agent Observability, ralph-orchestrator, Ralph for Claude Code, Awesome Ralph: keep (low value) / prompt- or skill-enforced, agent-decided gates; none gives code-validated handoffs.
- Anthropic "Building Effective Agents", multi-agent research system, Steering Claude Code, Encyclopedia of Agentic Coding Patterns, Agentic Workflow Patterns, Dive into Claude Code: keep (reading only): pattern references.
- reject: undeadlist agents, infrastructure showcase, Ralph Wiggum plugin (self-judged completion), SuperClaude, ECC, safety hooks (Safety Net, Dippy, Agent Guard), presence (advisory), Selvedge (SQLite), sandboxes (Brood Box, aicontainer, Cleat, Incus, machine), claude-control, Vibeyard, Happy Coder, Claude Threads, agnix, Schliff/Ctxlint, localpr, MDXG Redline, Story Skills, Archify, AgentWatch/OrcaReplay, goccc, TermaGITchi, claude-esp, Claumon, Netresearch skills, OSS Autopilot, pedrohcgs workflow: not a controller, or prompt-only, or UI.

### bradAGI/awesome-cli-coding-agents (884 lines)
- Orbi: keep (fetched; see Candidates).
- h5i: keep (borderline). The sub-agent's description said run metadata lives in `refs/h5i/*`; the fetched README did not mention that and described formal verification of Rust code instead. Unconfirmed; treat as unproven.
- AgentPlane: keep (fetched).
- Galley: keep (fetched).
- grite: keep (fetched).
- AgentDiff: keep (borderline): provenance in Git refs, signed attributions; Rust (per list).
- NEEDLE: keep (borderline): explicit deterministic state machine, no model picks the next step; Rust; unfetched.
- Vibestrate: keep (borderline; fetched): role seats and reviewer in a fresh process; README shows no per-seat output validation or ledger location; beta.
- zeroshot, kodo, ORCH, Team1-Factory, MartinLoop, TaskHandoff, Omnigent, Multica, Agent Orchestrator (Untrivial-ai), Dahrk, fab-kit, Forge (LucasDuys), Crewplane, Relay, oh-my-agent, Bernstein, OMK, Ivy Tendril, ADHDev: keep (borderline): multi-agent delivery loops with independent verifier or merge gates; none fetched; patterns only.
- GitHub Agentic Workflows: keep (fetched).
- OpenAI Symphony: keep (fetched).
- isitdone: keep (fetched). CodeVetter, molt, BlaBla, Squelette, SpecWave, PatchWarden, Hivelore, Grinta, Project Tiny Context Harness, Agentic Engineering Framework, TraceFold (Cedar policy, signed receipts), Loki Mode, great_cto: keep (borderline) / unfetched gate ideas.
- Mneme: keep (fetched).
- postmortemthis, CLITrigger: keep (borderline): cross-CLI review.
- reject: agy-auto (permission layer only), OpenCastle (prompt-driven), defract (closed), Beads (Dolt), clu (SQLite), Okto Nexus, Concord MCP, myc, Agent Executor (AX), Aeon, OSOP, mix2, llm-panel, LoopTroop, Codewhale, ralph-tui, py_ralph_frame, Ralph Workflow, wreckit, Ordewell, baya-cli, sage, agx, fractal, the-perfect-orchestrator, DevPilot, TeDDy, Kiro Crew, Docker Agent, claude-flow, gastown, Traycer, AgentsMesh, Agent Teams AI, Agentlas OS, loom, The Factory, Agon, 5dive, OpenHands, SWE-agent, Roo Code, Plandex, RA.Aid, Reasonix, San, OpenCodeReview, Better Agent, AgentBridge: single agent, closed, SQLite/DB state, swarm without gates, or not a workflow.
- Note: the list calls gastown "multi-agent work tracking"; I fetched gastownhall/gastown (see below) and it keeps work state in git-backed Beads.

### e2b-dev/awesome-ai-agents (5591 lines; entries date from 2023-24)
All reject. Closest: MetaGPT (role SOPs, LLM-written handoffs; prior art only), ChatDev, CrewAI, AutoGen, AutoPR (issue label to PR, no gates), Ellipsis (closed), Lemon Agent (LLM validation), Mutahunter (test generation), data-to-paper (rule-checked steps), English Compiler (spec to code, no IDs), Langroid, Agents (aiwaves; SOP layer, LLM controller). Others: single coding agents (Devika, OpenDevin, SWE-agent, GPT Pilot, Aider, Mentat), closed SaaS (Devin, Factory, Sweep, Codegen), no-code builders. Reason for all: LLM-driven flow, closed source, or no commit-status or validated-handoff mechanism.

### meirwah/awesome-workflow-engines (134 lines)
- Temporal, Cadence, Restate, Dapr Workflows, Inngest, Kitaru: keep (reference only): durable state machines; heavy servers or databases; no stdlib-only fit.
- Cordum: keep (borderline): "policy before dispatch, approval gates, audit trail" for agents; unfetched.
- Tork, Dagu, Goflow, SciPipe, go-taskflow, uTask: keep (borderline): Go workflow engines; they would be Go modules or services, so not dependencies.
- Symfony Workflow, Django River, Petri Flow, easy-rules: keep (design reference): transition tables, authorization on transitions, rule engines.
- reject: Flyte, River Pro, Argo, Brigade, CDS, Conductor, Zeebe, Camunda, Activiti, Flowable, jBPM, Bonita, Kogito, nFlow, Imixs, Flowret, Workflow Core, Elsa, Viewflow, SpiffWorkflow, StackStorm, Huginn, Windmill, n8n, Activepieces, Kestra, Airflow, Prefect, Dagster, DBOS, Durable Workflow, Durable Task Framework, Oban, Taskade MCP, Codehooks, Step Functions, Google Workflows, Logic Apps: heavy, JVM/.NET/Python/SaaS, or data pipelines.

### sdras/awesome-actions (572 lines)
- WyriHaximus/github-action-wait-for-status: keep (reference): reads commit statuses and checks.
- talos-systems/conform: keep (reference): policy checks on commits; Go (likely).
- cirrus-actions/branch-guard: keep (reference): blocks merges when target-branch checks fail.
- tibdex/autosquash, maxkomarychev/merge-pal-action, pascalgn/automerge-action: keep (reference): merge bots that respect branch protection.
- jpmcb/prow-github-actions: keep (reference): chat-ops and merge policy.
- benjefferies/branch-protection-bot: keep (borderline): changes branch protection from automation (needs admin; the App lacks it).
- actions/github-script: keep (borderline): can post statuses inside a workflow, but the controller must post them itself, and the App has no workflows permission.
- peter-evans/create-pull-request: keep (borderline): PR creation from automation.
- reject: the rest (checkout, labeler, stale, toolkit, approve/label actions, locks, commitlint, ticket checks, CODEOWNERS validator, act, and similar): not controller-posted gates or not related.

## The search file

`sources/gh-search.tsv`: 240 result rows from 24 queries, read in full by me. Query to fact words, then verdicts. Rows I judge only from the description are marked "(desc)"; unmarked entries I fetched.

**agent task loop (F-0003#45)**
- agentscope-ai/AgentTeams: reject: Matrix-room collaboration OS (desc).
- saltbo/agent-kanban: reject: task board (desc).
- PageAI-Pro/ralph-loop: reject: shell loop until done (desc).
- GaosCode/PlanWeave: keep (borderline): file-backed loop with claimable tasks and review agents, "records every run" (desc); unfetched, TypeScript, MIT.
- AMAP-ML/LongHorizon-Harness: reject: computer-use harness, not the loop here (desc).
- agentrq, agent-apprenticeship, codeAgent, OpenETA, LoongFlow: reject: not the problem (desc).

**agent handoff schema (F-0003#59)**
- max9159/agent-context-store: keep (see Candidates).
- faresrafat3/agent-handoff: keep (see Candidates).
- ZibbyDev/agent-workflow: keep (see Candidates).
- dabit3/agent-handoff: reject: TypeScript, schema validation and chain tracking, 14 stars (desc); same idea as above.
- openenvelope/schema: keep (borderline): versioned document defining agents, roles, handoffs, access policies (desc); no code, unfetched, Apache-2.0.
- talocode/handofflane: reject: npm, "semantic" validation (desc).
- shin4141/decision-os-v12-completion-integrity: keep (borderline): schema, CLI and GitHub Action against "incomplete work called done" (desc); unfetched, Python, MIT.
- agentsentinel-handoff, moona3k/handoff, roymcfarland/agent_handoff_guide: reject: serialization schema only, URL sharing, or a playbook (desc).

**agent state machine (F-0003#59)**
- statewright/statewright: keep (known limit; see Candidates).
- statelyai/agent: reject: XState LLM agents (desc).
- AgenticStateMachines, LLM-State-Machine, llmstatemachine, Graph-of-States, frankterm, MemMachine, zuey-pi-setup, whatsapp-ai-pix-agent: reject: agent patterns, memory or unrelated (desc).

**commit status bot (F-0003#47)**
- tink-ab/four-eyes: reject: Go bot that sets a bors status on review count; 2018, archived-looking (desc).
- jenseickmeyer/github-commit-status-bot, dcob, giles, BuildBot, travis-to-slack, git-bot, hook-proxy, GitLab-Jira bot, bot-lua: reject: unrelated or old (desc; bot-lua row is a pasted file).

**required status checks (F-0003#58)**
- cloudposse/github-status-updater: keep (see Candidates).
- pkgdeps/automerge-gate: keep (borderline): one required check that gates Enable Auto Merge (desc); TypeScript, MIT, 2026-10-06; unfetched.
- Nmishin/GitHub-PR-Status-Validator: keep (borderline): detects force-merged PRs by reading statuses and check runs (desc); audit idea for "a check did not run".
- BrianCArnold/RequireUserApproval: reject: App emitting an approval status (desc); 2024; shows the pattern only.
- etsy/combined-status-check-action: reject: Action that combines checks into one (desc).
- others (status-page, lab, test repos): reject.

**ai code review verifier (F-0003#45)**
- nexpeakcore/deepseek-harness-pr-review: reject: LLM PR review (desc).
- zhiyuwang720-dev/CodeSage, adshao/flounder, aacr-bench: reject: review products or benchmark (desc).
- pjwang2022/Monstrare: keep (borderline): gate-based workflow where "nothing is called done without evidence" (desc); HTML/skills; unfetched.
- nurettincoban/ai-prd-workflow, KimYx0207/Meta_Kim, skainguyen1412/antigravity-superpowers, paperjury, latex-arxiv-SKILL: reject: skills and prompts.

**merge bot pull request (F-0003#58)**
- palantir/bulldozer: keep (reference): Go GitHub App merging PRs by policy; Apache-2.0, pushed 2026-10-14 (desc). Read its config for how it decides to merge on statuses.
- bors-ng/bors-ng: keep (reference): merge queue that tests the merged state; Elixir, Apache-2.0, last push 2024-04-04 (desc); relevant to "merge the exact head that was checked".
- traefik/lobicornis: keep (borderline): Go bot to update and merge PRs, Apache-2.0, 2026-09-07 (desc).
- voyagegroup/popuko: reject: older Go merge bot (desc).
- others (miss-islington, squalrus/merge-bot, ghb0t, mergeatron, Git-Enforcer): reject: different aims or stale.

**requirements traceability / trace tool (F-0003#51, #62)**
- rtmx-ai/rtmx: keep. awslabs/duvet: keep. doorstop: keep. bmw lobster: keep (see Candidates).
- SAP-archive/quality-continuous-traceability-monitor: reject: Go, Apache-2.0, links requirements to tests and results; archived (desc).
- goeb/reqflow: keep (borderline): requirements tracing in C++ with "other" license (desc); unfetched.
- konstantin-hatvan/traceability-tool: reject: markdown requirements tool, 2023 (desc).
- sphinx-traceability-extension, tobhey/finegrained-traceability, Traceability-Tool-Python, JanusTrace, simon3z/tmatrix, others: reject: Sphinx-bound, research, or unrelated.
- adrianpuiu/specification-document-generator: reject: prompt framework claiming "100% coverage" (desc).

**spec to requirements ai (F-0003#23)**
- mori-shin-x/artgraph: keep (borderline): "deterministic spec-to-code context", graph of requirements, docs, code, tests (desc); TypeScript, MIT, 2026-10-05; unfetched.
- andrewstellman/quality-playbook: reject: skill that derives requirements with an LLM (desc).
- hieubui2409/product-spec: reject: AGPL critique harness (desc).
- others (samurai-agent, kiro prompts, specifai-mcp, builder-skills, Co-OmniSpec, prd-pilot, CodeSpec): reject: prompt or tool workflows without a checker.

**requirements linter / prd lint (F-0003#24, #51)**
- salab/reqchecker, jonverrier/RequirementLinter, goodwrite, speclint, prdlinter, Jaimin4727/prd-writer: reject: prose-quality linters (desc); none checks IDs, trace or coverage.
- nilukush/prd-as-a-code: keep (borderline): PRD as schema-validated text compiled to artifacts (desc); JavaScript, MIT, 2026-09-22; unfetched.
- zlint, dep linters, SCSS linter: reject: not requirements.

**autonomous software delivery (F-0001#1)**
- synaptixs/spine: keep (borderline): "governed, provenance-grounded autonomous software delivery" (desc); Python, MIT; unfetched.
- q3ok/coordinated-agent-team: reject: prompt-driven "deterministic workflow", strict I/O contracts; AGPL; no code gate (desc).
- digipulse-engineering/GAAI-framework: reject: Markdown + YAML + bash prompts (desc).
- cmAIdx/headless-claude-automation-template, Coderama, Free-devloper/autonomous-software-delivery-organization, ai-alpha-squad, launchapp-dev/ao, seanyao/roll: reject: templates or platform claims (desc).
- jonathanpopham/estate-agent: reject: Go, "guarded agentic delivery" for issues (desc), 5 stars; unfetched.

**agents git records (F-0001#1)**
- entireio/cli: keep (see Candidates).
- codeprakhar25/agentdiff: keep (borderline): Rust, Apache-2.0 here; signed attributions in Git history (desc).
- Ian-stetsenko/lore-protocol: keep (borderline): decision records in commit messages (desc); TypeScript.
- mbeacom/adrkit: keep (borderline): machine-readable ADRs with lint, CI and MCP (desc); TypeScript, Apache-2.0; relevant to the option-filter record format.
- kerbelp/metatron: reject: markdown knowledge files for agents.
- git-sim, vibe-graph, trail, termaxa, agent-undo: reject: different aim.

**coding agent guardrails (F-0001#3)**
- chainloop-dev/chainloop: keep (borderline): Go, Apache-2.0; signed attestations of agent sessions and CI checked against policies (desc); heavy platform, unfetched.
- clay-good/OpenLore: reject: memory and guardrails with no LLM in hot path (desc); local rules.
- LuD1161/agentjail: reject: Go, per-tool-call policy before it runs (desc); permission layer, not workflow gates.
- KbWen/agentic-os, SumonMSelim/agentguard, shellfirm, leash, ultimate_bug_scanner, saas-starter, aws bedrock samples: reject.

**rulesets as code (F-0003#43)**
- Vivswan/github-settings-as-code, theztefan/rulesets-as-code-action, noshiro-pf/github-settings-as-code, AlexAsAService/grac: keep (known limit): declarative rulesets apply through the API, which needs repository administration; the controller's App has none. A human or separate admin identity must apply the ruleset; the agents then cannot weaken it.
- aws config engine, Dome9 contracts, cs156 rulesets, terraform repos: reject: unrelated.

**shape up bet (F-0003 shape and bet)**
- ramboz/shaper, rhyslawsn/shape-up-skills, rafoflow, florianilch/mermaid-hillchart, ai-native-pm-os: reject: Claude skills or plugins for Shape Up; no code checks (desc).
- remaining rows (land development, Bee Movie script, weight loss): reject: noise.

**multi agent software development (F-0003#45)**
- untra/operator: reject: Rust kanban-shaped orchestrator (desc).
- jonnyzzz/run-agent, hoangsonww/AI-Agents-Orchestrator, helmor, DevOpsGPT, codefuse-chatbot, just-ship, AgentDeskSkill, crewswarm, cursor-agents: reject: multi-agent frameworks without code-validated handoffs or statuses (desc).

**issue comment structured data (#68)**
All ten rows are Instagram scrapers, labelers or learning repos: reject. Nothing found for "an artifact that exists only as an issue comment has a place and a hash". Orbi (issues as the only state store) is the nearest.

**pull request provenance**
- Kartik24Hulmukh/jittest: keep (borderline): differential test gate for agent PRs, Ed25519 signed receipts re-verifiable (desc); Python, Apache-2.0; unfetched.
- semanticash/cli: reject: Go, provenance and attribution (desc).
- Dev-31/provenance-layer2: reject: Go GitHub App verifying attestations on PRs (desc); 1 star, 2026-04.
- twilio-labs/gordon: reject: status-check app for `about.yaml` (desc); 2025.
- P-r-e-m-i-u-m/ai-change-passport, forgeproof-verify, gitiq, memPR, provenance, ROCmKernelWiki: reject.

**approval comment bot (bet approval)**
All reject (social-media bots, chat approvals, a 2014 Ruby PR-approvals app). PerryLink/dsh-github: reject: DeepSeek review bot with status gate (desc). Nothing found for "one comment that the controller copies into Git".

**agent pull request gate (F-0003#47)**
- tiammomo/RepoSteward: keep (borderline; fetched): policy-gated issue-to-PR control plane; publish and merge are separate human decisions; state kept outside the repo; Python, MIT, 0.1.0.
- mojunrubest/agent-delivery-gate: keep (see Candidates).
- krishagarwal314/CodeJury: keep (borderline): Python, MIT; "gate pull requests with deterministic QA and ensemble code review" (desc); unfetched.
- FreePeak/devagent: keep (borderline): Go, MIT; ticket to tested PR, "evidence-gated orchestration" (desc); unfetched.
- xdlc-labs/airlock: reject: Go CI release gate for prompt/permission changes (desc); different problem.
- ktwu01/agent-pr-gate, maida-assert, blastproof, agent-gate, pr-review-agent: reject: contribution policy or review products.

**spec driven development (F-0003#51)**
- Fission-AI/OpenSpec (72k stars), gsd-build/get-shit-done, buildermethods/agent-os, Pimzino/spec-workflow-mcp, spec_driven_develop, formulahendry/mcp-server-spec-driven-development, awesome-spec-driven-development, hello-sdd, spec-driven-agentic-development: reject (desc): prompt-and-markdown workflows run by the agent; none has IDs with trace to problem-statement text and a code checker. OpenSpec is the best-known format reference for spec deltas; I did not open it, so I cannot say whether it checks anything.

**agent task state (task register)**
- MaxHu-xuan/task-state-guard: reject: SQLite reconciliation of stuck tasks (desc).
- oceanbase/powercontext: reject: memory/context layer (desc).
- LongHorizon-Harness, thu-nmrc/OpenHarness (persistent state, strict task contracts): reject: Python harnesses; state not in Git (desc).
- others: reject.

## What no candidate covers

- **Typed handoffs in Go stdlib.** The only handoff-schema projects (agent-context-store, agent-handoff, agent-workflow) are TypeScript or Python, and none ties validity to a controller-owned transition table plus a hash of an artifact held as an issue comment. gh-aw's safe-outputs idea is the nearest statement of "agent text does not make a handoff valid". Go has no JSON Schema validator in the stdlib, so validation would be hand-written.
- **A place and hash for an artifact that exists only as an issue comment.** I found nothing. Orbi uses issues as state but I saw no hash check.
- **Per-class commit statuses posted by the controller and required by rules,** with "not active is not passed". Only the GitHub docs (and a Go status-updater using a user token) cover it. No agent-delivery project in the sources shows it done with an App token for four named classes.
- **Code filter of architecture options** (failed or missing row dropped, pass with no evidence dropped, one survivor). I found ADR templates with decision matrices, Mneme (checks proposals against ADRs), and nothing that filters option rows by evidence.
- **A brief to the idea owner checked for parts and length, answered by one comment copied into Git.** Nothing found.
- **Requirements derived from a problem statement,** with ID, trace to the stating text, acceptance criterion, and a checker failing on unresolved ID, uncovered need, criterion-less requirement, or delivered requirement with no specification section. Duvet (quote-based trace), Doorstop, rtmx and LOBSTER cover parts; none starts from a prose problem statement with the four failure rules together.
- **Controller state as records on a branch only the controller may update,** in Git. grite and entireio/cli keep records in dedicated refs; AgentPlane keeps them in the working tree. None targets "only the controller's identity may write". That depends on a ruleset; the App's lack of administration means a human must create it.
- **A full loop with a verifier on a different harness from the author and a plan review by a different model** with code-checked transitions. Orbi, Galley, zeroshot and Vibestrate come near with agent-judged steps.

## Limits

- I did not run, install or read the source code of any candidate; all facts are from README pages and GitHub docs as fetched on 2026-10-15 by a summarizing fetch tool. That tool sometimes returned no commit date, no release list, or an empty language list; those cells are "not stated" or "inferred" and could be wrong.
- No licence, language or activity could be taken from the five lists, which show none; I checked those only for fetched entries. Many list-kept entries (NEEDLE, zeroshot, kodo, ORCH, MartinLoop, Cordum, Consort, TraceFold and others) were not fetched.
- The five lists were read by sub-agents, not line by line by me. Their summaries may contain errors; one is known: the sub-agent's note on h5i (refs under `refs/h5i/*`) was not confirmed by the h5i README I fetched.
- SPEC.md of Symphony, the docs of Orbi, AgentPlane `acr.json` schema, grite's data model, and OpenSpec's check behaviour were not read.
- Search-file rows marked "(desc)" were judged from a one-line description only.
- GitHub behaviours I could not confirm from a source: how a required check that never reports is treated (the docs page is silent); whether a ruleset can exempt one branch by ref condition for an App; whether an App can create or update a ref on a restricted branch without bypass. One community thread says App bypass is whole-ruleset.
- I did not open spec-trace, cargo-ovft, jamb, TRLC, SpecTrace, bulldozer, bors-ng or chainloop beyond search results or the tsv description.
- The web search tool is US-only and returns few results per query; I ran 6 searches, so coverage is not exhaustive.
