# Public-solution search: M2e / T-xs4y

Search date: 2026-10-15. "Keep" means useful prior art, **not** approved as a Go dependency or proof that the full workflow is implemented. Repository-search dates and language/license labels below are the supplied snapshot, not independently audited. GitHub documentation is a platform contract, not open-source code. No candidate was run.

## Queries

| Query | Words from the problem | What it found |
|---|---|---|
| `GitHub App commit statuses required checks pending skipped merge protection documentation` | F-0001#5, F-0003#47/#58; four controller-posted checks | [Protected branches](https://docs.github.com/en/repositories/configuring-branches-and-merges-in-your-repository/managing-protected-branches/about-protected-branches), [statuses](https://docs.github.com/en/pull-requests/reference/status-checks), [troubleshooting](https://docs.github.com/en/pull-requests/how-tos/merge-and-close-pull-requests/troubleshooting-required-status-checks): required checks allow success, neutral or skipped; controller must independently reject missing/inactive checks. |
| `git repository state machine typed artifact role handoff schema agent workflow open source` | F-0003#45/#59, F-0001#1 | [agent-context-store](https://github.com/max9159/agent-context-store), [agentctl structured handoffs](https://github.com/opensourceops/agentctl/blob/main/docs/guides/STRUCTURED_HANDOFFS.md), TAKT; different persistence and assurance models. |
| `requirements traceability markdown specification validation acceptance criteria open source tool` | F-0003#23/#24/#25/#51/#62 | [reqmd](https://reqmd.dev/), [supersigil](https://github.com/jonisavo/supersigil), spec-lint; distinguish warnings from failing coverage. |
| `architecture decision matrix evidence criteria machine validated alternatives open source` | Architecture row/evidence elimination and one survivor | [Decision Studio](https://github.com/FieldstateNZ/workspec-decision-studio), [adrkit](https://github.com/mbeacom/adrkit); neither documents precisely this filter. |
| `GitHub pull request merge exact head SHA expectedHeadOid REST API` | Merge exact checked head after close-out | [REST merge `sha`](https://docs.github.com/en/rest/pulls/pulls), [GraphQL `expectedHeadOid`](https://docs.github.com/en/graphql/reference/pulls); both support head compare-and-swap. |
| `GitHub app branch protection restrict pushes app bypass rulesets status checks source app` | Controller-only state branch; rule protection | [Ruleset restrictions](https://docs.github.com/en/repositories/configuring-branches-and-merges-in-your-repository/managing-rulesets/available-rules-for-rulesets); app identity may be selected as check source; bypass is not the same as exclusive writer. |
| `site:github.com/opensourceops/agentctl "License" "Go" "state"` | Typed graph, durable transitions | [Architecture](https://github.com/opensourceops/agentctl/blob/main/docs/ARCHITECTURE.md): Rust, SQLite and CAS, not Git. |
| `site:github.com/max9159/agent-context-store README license handoff git schema` | Schema-validated Git handoffs | [README](https://github.com/max9159/agent-context-store/blob/main/README.md), [intro](https://github.com/max9159/agent-context-store/blob/main/docs/ACS_CONTEXT_STORE_INTRO_en.md): handoff schema documented as permissive. |
| `site:github.com/FieldstateNZ/workspec-decision-studio LICENSE validation evidence criteria` | Architecture choice as validated evidence | [README](https://github.com/FieldstateNZ/workspec-decision-studio/blob/main/README.md), [license](https://github.com/FieldstateNZ/workspec-decision-studio/blob/main/LICENSE): Git-native YAML and deterministic ADR rendering. |
| `site:github.com/reqmd/reqmd license requirements trace checks` | Requirement IDs and coverage | Ambiguous name: [reqmd.dev](https://reqmd.dev/about/) is GPL-3.0; [voedger/reqmd](https://github.com/voedger/reqmd) is a *different* Apache-2.0 Go project. |
| `site:docs.github.com/en/rest/commits/statuses create commit status permissions context state combined` | Controller statuses, four contexts | [REST commit statuses](https://docs.github.com/en/rest/commits/statuses): POST to SHA with context and state; combined state alone cannot establish that each named check ran. |
| `site:github.com/speckit/spec-kit license requirements checklist traceability` | Prose to specification and plan | [Spec Kit](https://github.com/github/spec-kit) templates and optional agent-driven consistency analysis, not a mechanical full-coverage proof. |
| `site:github.com/nrslib/takt license workflow transitions output contracts git` | Different harnesses, roles and transitions | [TAKT workflow docs](https://github.com/nrslib/takt/blob/f70a91e1/docs/workflows.md): explicit YAML rules, report contracts and command gates. |
| `site:github.com/mbeacom/adrkit license deterministic validation git decisions` | Rule protection and evidence-bearing decisions | [adrkit README](https://github.com/mbeacom/adrkit/blob/main/README.md): Git corpus, deterministic model-free evaluation pass; no automatic prescribed survivor. |
| `site:github.com/GaosCode/PlanWeave README file-backed git agents review` | Different agents, plans and review | [PlanWeave README](https://github.com/GaosCode/PlanWeave/blob/main/README.md): review blocks, local file-backed run history; not protected Git state. |
| `site:github.com/cloudposse/github-status-updater README github status` | Four named statuses | [Status-updater README](https://github.com/cloudposse/github-status-updater/blob/main/README.md): status state/context posting and separate branch-protection action. |
| `site:github.com/palantir/bulldozer README merge status checks` | Safe merge | [Bulldozer README](https://github.com/palantir/bulldozer/blob/develop/README.md): respects branch protection, relies on correct configuration. |
| `site:github.com/LuD1161/agentjail README policy guardrails` | Agents may not change governing rules | [Agentjail README](https://github.com/LuD1161/agentjail/blob/main/README.md): hook checks are cooperative; optional OS sandbox; fallback not invariably fail-closed. |
| `site:github.com/Fission-AI/OpenSpec README changes specs artifacts` | Specification synthesis | [Editing changes](https://github.com/Fission-AI/OpenSpec/blob/main/docs/editing-changes.md): editable Markdown proposals/specs, explicitly no locked approval gate. |

## Candidates

Dates marked "pushed" are from `sources/gh-search.tsv`; publication date is used only where cited. "Not stated" means the cited source does not establish the fact. "Git: yes" means artifacts *can* reside in the target repo, not that a protected controller-only branch exists. "Det." refers only to the named portion, not LLM-generated text.

| Candidate and own documentation | License | Language | Last activity | Part solved | State in project's Git? | Deterministic? | Keep as / why |
|---|---|---|---|---|---|---|---|
| [GitHub commit-status REST API](https://docs.github.com/en/rest/commits/statuses) | not stated (service API) | HTTPS/JSON | not stated | POST four named `context` statuses to a SHA, read status | statuses on forge, not Git | status aggregation defined | API usage; stdlib HTTP possible, but require explicit per-context inspection. |
| [GitHub protected branches](https://docs.github.com/en/repositories/configuring-branches-and-merges-in-your-repository/managing-protected-branches/about-protected-branches) / [troubleshooting](https://docs.github.com/en/pull-requests/how-tos/merge-and-close-pull-requests/troubleshooting-required-status-checks) | not stated (service docs) | configuration/API | not stated | required status source and head/merge-commit semantics | rules on forge, not Git | rule enforcement yes; no check execution proof | known limit; `skipped`/`neutral` acceptable to GitHub, unlike F-0001#5. |
| [GitHub rulesets](https://docs.github.com/en/repositories/configuring-branches-and-merges-in-your-repository/managing-rulesets/available-rules-for-rulesets) | not stated (service docs) | configuration/API | not stated | restrict state branch updates and expected app status source | rules on forge, not Git | enforcement yes | API usage; admin must provision rules; App lacks administration permission. |
| [GitHub PR merge REST](https://docs.github.com/en/rest/pulls/pulls) / [GraphQL](https://docs.github.com/en/graphql/reference/pulls) | not stated (service docs) | HTTPS/JSON | not stated | `sha` / `expectedHeadOid` reject head change during merge | PR on forge, not Git | compare-and-swap yes | API usage; read checked SHA again after close-out commit. |
| [agent-context-store](https://github.com/max9159/agent-context-store/blob/main/README.md) / [schema notes](https://github.com/max9159/agent-context-store/blob/main/docs/ACS_CONTEXT_STORE_INTRO_en.md) | MIT | TypeScript / Node | pushed 2026-07-06 | role artifacts, handoff records and validation | yes in-repo `.acs/`; dedicated repo also supported | validation yes; content no | pattern/format; permissive handoff schema is insufficient as sole gate. |
| [agentctl](https://github.com/opensourceops/agentctl/blob/main/docs/guides/STRUCTURED_HANDOFFS.md) / [architecture](https://github.com/opensourceops/agentctl/blob/main/docs/ARCHITECTURE.md) | [other](https://github.com/opensourceops/agentctl) (specific terms not verified) | Rust | [v0.3.0 2026-07-28](https://github.com/opensourceops/agentctl/releases/tag/v0.3.0) | compiled workflow graph, typed output, deterministic handoff task, audit | no: SQLite + CAS outside Git | compiler/handoff yes; agent content no | method; explicit typed routing, not a Go dependency. |
| [TAKT](https://github.com/nrslib/takt/blob/f70a91e1/docs/workflows.md) / [repo](https://github.com/nrslib/takt) | MIT | TypeScript | not stated | roles, provider routing, step rules, report contracts and machine command gates | report files/worktrees possible; full state not stated | command gates yes; agent decisions no | pattern; caution: text-described quality gates are not machine checks. |
| [reqmd.dev](https://reqmd.dev/) / [about](https://reqmd.dev/about/) | GPL-3.0 | Go | not stated | Markdown requirement IDs, trace graph and inbound coverage | yes: Markdown single source | schema/link checks yes; defaults can warn | method; make missing links/coverage fatal here; unrelated to voedger/reqmd. |
| [supersigil](https://github.com/jonisavo/supersigil) | not stated | not stated | not stated | typed criterion/evidence/test graph in repo Markdown | yes | structural verify yes; semantics not stated | format; promising trace graph; license/activity unchecked. |
| [GitHub Spec Kit](https://github.com/github/spec-kit) / [workflow docs](https://github.com/github/spec-kit/blob/main/docs/reference/agentic-sdd.md) | MIT | Python CLI / templates | not stated | specify-plan-tasks; consistency and checklist guidance | yes, feature Markdown artifacts | no: agent-based synthesis/analysis | method; do not confuse generated checklists with code-enforced trace. |
| [adrkit](https://github.com/mbeacom/adrkit/blob/main/README.md) | Apache-2.0 (schema [CC0](https://github.com/mbeacom/adrkit/blob/main/schema/LICENSE)) | TypeScript | pushed 2026-10-11 | typed ADRs in Git; model-free proposal evaluation | yes | `adr evaluate` / queue documented deterministic | format/method; does not eliminate every option lacking an evidence row. |
| [WorkSpec Decision Studio](https://github.com/FieldstateNZ/workspec-decision-studio/blob/main/README.md) | [Apache-2.0](https://github.com/FieldstateNZ/workspec-decision-studio/blob/main/LICENSE) | TypeScript | not stated | validated YAML options, costs and weighted decision; ADR | yes; no database | validation and ADR render yes | format; its scoring is not the required fail/missing/no-evidence filter. |
| [OpenSpec](https://github.com/Fission-AI/OpenSpec) | MIT (search snapshot) | TypeScript | pushed 2026-10-15 | spec-driven change artifacts | repo files (per documentation); complete controller state not stated | specification synthesis no | method; not a traceability checker. |
| [PlanWeave README](https://github.com/GaosCode/PlanWeave/blob/main/README.md) | MIT (search snapshot) | TypeScript | pushed 2026-09-30 | file-backed plans, claimable tasks, multi-provider review routing | workspace files; committed into target Git not stated | local scripts yes; agent review no | pattern; durability does not imply protected Git state. |
| [cloudposse/github-status-updater README](https://github.com/cloudposse/github-status-updater/blob/main/README.md) | Apache-2.0 (search snapshot) | Go | pushed 2026-09-04 | sets `pending`/`success`/`error`/`failure` with named context | GitHub statuses, not Git | request yes; gate verdict not stated | API usage; its branch-protection update needs privileges the App lacks. |
| [palantir/bulldozer README](https://github.com/palantir/bulldozer/blob/develop/README.md) | Apache-2.0 (search snapshot) | Go | pushed 2026-10-14 | waits for configured checks/reviews and auto-merges PR | PR on GitHub; internal state not stated | configuration-dependent | tool/pattern; no evidence of checked-SHA role workflow; must configure branch protection first. |
| [agentjail README](https://github.com/LuD1161/agentjail/blob/main/README.md) | Apache-2.0 (search snapshot) | Go, OPA/Rego | pushed 2026-09-24 | local per-tool agent policy checks | policy local; project Git not stated | policies yes if reachable; default ordinary-hook fallback may permit | pattern/known limit; not stdlib-only and local hooks do not protect remote Git refs. |

## The lists

List identity and pinned commits: `sources/lists.tsv`. Entries below address overlapping functionality only; broad chatbot and unrelated CI/infrastructure entries are omitted.

### bradAGI/awesome-cli-coding-agents

- [Roo Code CLI](https://github.com/RooCodeInc/Roo-Code): reject: architect/code modes, not independently checked handoffs.
- [gptme](https://github.com/gptme/gptme): reject: Git-backed memory is not controller-owned transitions.
- [Devon](https://github.com/entropy-research/Devon): reject: Git workflow agent, not gates.
- [Codebuff](https://github.com/CodebuffAI/codebuff): reject: multi-agent coding, not enforceable roles.
- [DvalinCode](https://github.com/arthurpanhku/dvalincode): reject: local policy, not forge enforcement.
- [Codex Infinity](https://github.com/lee101/codex-infinity): reject: CI repair loop, not checked merge SHA.
- [WorkGround2](https://github.com/KiddPhenix/WorkGround2): reject: checkpoints, not role artifacts.
- [molt](https://github.com/solvyxtech/molt): reject: done receipts but no stated inter-role schema.
- [GitClaw](https://github.com/open-gitagent/gitclaw): reject: versioned agent rules, not rules protected from agents.
- [Agentlas OS](https://github.com/agentlas-ai/Agentlas-OS): reject: multi-harness security gates, not protected project-Git transitions.

### e2b-dev/awesome-ai-agents

- [Aider](https://github.com/paul-gauthier/aider): reject: coding agent, not role transition controller.
- [AutoGen](https://github.com/microsoft/autogen): reject: agent conversation framework; no project-Git authority.
- [AutoPR](https://github.com/irgolic/AutoPR): reject: PR automation without specified checked-SHA gate.
- [CAMEL](https://github.com/camel-ai/camel): reject: multi-agent dialogue, not typed Git handoff.
- [ChatDev](https://github.com/OpenBMB/ChatDev): reject: software role simulation, not independently enforced gate.
- [CodeFuse-ChatBot](https://github.com/codefuse-ai/codefuse-chatbot): reject: repository assistant, not task controller.
- [CrewAI](https://github.com/joaomdmoura/crewai): reject: agent delegation, not controller-protected records.
- [DevOpsGPT](https://github.com/kuafuai/DevOpsGPT): reject: delivery agents, no demonstrated exact-SHA merge.
- [GPT Pilot](https://github.com/Pythagora-io/gpt-pilot): reject: builds apps, lacks stated cross-harness checks.
- [MetaGPT](https://github.com/geekan/MetaGPT): reject: role-based generation without required Git-gated handoff proof.
- [OpenDevin](https://github.com/OpenDevin/OpenDevin): reject: coding harness, not four controller statuses.
- [SWE Agent](https://github.com/princeton-nlp/SWE-agent): reject: issue solving, not specification trace checker.
- [Sweep](https://sweep.dev/): reject: issue-to-PR automation, no stated protected state branch.

### hesreallyhim/awesome-claude-code

- [Agent Collab Skills](https://github.com/WenyuChiou/agent-collab-skills): reject: collaboration prompts, not mandatory typed envelopes.
- [AgentSys](https://github.com/avifenesh/agentsys): reject: PR workflows, not code-owned validity.
- [Claude Code Harness](https://github.com/Chachamaru127/claude-code-harness): reject: plan-work-review loop, one harness.
- [Claude Code Hook Comms](https://github.com/aannoo/claude-hook-comms): reject: hook messages, not project-Git records.
- [Claude Code Workflows](https://github.com/OneRedOak/claude-code-workflows): reject: review workflow, not four statuses.
- [Consort](https://github.com/databricks-solutions/consort): reject: deterministic orchestration but Lakebase state, not Git.
- [Plannotator](https://github.com/backnotprop/plannotator): reject: plan review UI, not separate-model enforced review.
- [RIPER Workflow](https://github.com/tony/claude-code-riper-5): reject: phase guidance, not forge-enforced gates.
- [napkin-to-blueprint](https://github.com/tsmztech/napkin-to-blueprint): reject: idea-to-blueprint prompts, not validated trace.
- [Ralph for Claude Code](https://github.com/frankbria/ralph-claude-code): reject: iteration loop, not independent verifier.
- [localpr](https://github.com/gilles-g/localpr): reject: local review, not merge protocol.

### meirwah/awesome-workflow-engines

- [Cordum](https://github.com/cordum-io/cordum): reject: policy/approvals, external control-plane state.
- [Dagu](https://github.com/yohamta/dagu): reject: workflow executor, not SDLC handoff validity.
- [Django River](https://github.com/javrasya/django-river): reject: explicit transitions but Django dependency/state.
- [Restate](https://github.com/restatedev/restate): reject: durable execution, not Git as system of record.
- [Temporal](https://github.com/temporalio/temporal): reject: durable service state outside project Git.
- [Tork](https://github.com/runabol/tork): reject: Go workflow engine, not standard-library-only Git controller.
- [uTask](https://github.com/ovh/utask): reject: YAML business process, not Git status protection.
- [Symfony Workflow](https://github.com/symfony/workflow): reject: state-machine pattern in PHP, not Go dependency.

### sdras/awesome-actions

- [Branch Protection Bot](https://github.com/benjefferies/branch-protection-bot): reject: temporarily disables protection; violates rule protection.
- [Auto-Approve PRs](https://github.com/hmarr/auto-approve-action): reject: approval action is not independent verification.
- [Automatically merge PRs That Are Ready](https://github.com/pascalgn/automerge-action): reject: generic readiness is insufficient.
- [Automatically Update PRs with Outdated Checks and Squash and Merge](https://github.com/tibdex/autosquash): reject: head-changing operation requires fresh check.
- [Merge Pal](https://github.com/maxkomarychev/merge-pal-action): reject: generic PR merge without specified exact head.
- [Block PR merges when Checks for target branches are failing](https://github.com/cirrus-actions/branch-guard): reject: target checks alone not required four head contexts.
- [Pull Request Description Enforcer](https://github.com/derkinderfietsen/pr-description-enforcer): reject: PR text only, not handoff schema.
- [Label PR once it has a Specified Number of Approvals](https://github.com/pullreminders/label-when-approved-action): reject: label not commit-status gate.

## The search file

Source: `sources/gh-search.tsv`; one line per entry below, including a number of loosely related hits to make exclusions visible. Repository names link to the searched result. A keep here refers to the candidate table or a limited pattern, never approval to install it.
- [agentscope-ai/AgentTeams](https://github.com/agentscope-ai/AgentTeams): reject: task/agent loop lacks documented end-to-end role gates.
- [saltbo/agent-kanban](https://github.com/saltbo/agent-kanban): reject: task/agent loop lacks documented end-to-end role gates.
- [PageAI-Pro/ralph-loop](https://github.com/PageAI-Pro/ralph-loop): reject: task/agent loop lacks documented end-to-end role gates.
- [ray-r-ren/agent-apprenticeship](https://github.com/ray-r-ren/agent-apprenticeship): reject: task/agent loop lacks documented end-to-end role gates.
- [agentrq/agentrq](https://github.com/agentrq/agentrq): reject: task/agent loop lacks documented end-to-end role gates.
- [Harvil1/codeAgent](https://github.com/Harvil1/codeAgent): reject: task/agent loop lacks documented end-to-end role gates.
- [GaosCode/PlanWeave](https://github.com/GaosCode/PlanWeave): keep: file-backed plan/review loop pattern; Git commitment unproven.
- [AMAP-ML/LongHorizon-Harness](https://github.com/AMAP-ML/LongHorizon-Harness): reject: task/agent loop lacks documented end-to-end role gates.
- [OpenMOSS/OpenETA](https://github.com/OpenMOSS/OpenETA): reject: task/agent loop lacks documented end-to-end role gates.
- [baidu-baige/LoongFlow](https://github.com/baidu-baige/LoongFlow): reject: task/agent loop lacks documented end-to-end role gates.
- [dabit3/agent-handoff](https://github.com/dabit3/agent-handoff): reject: schema transfer but project-Git persistence not shown.
- [openenvelope/schema](https://github.com/openenvelope/schema): reject: schema format only; Git and transition enforcement not stated.
- [ZibbyDev/agent-workflow](https://github.com/ZibbyDev/agent-workflow): reject: schema-enforced cross-harness graph; Git authority not shown.
- [max9159/agent-context-store](https://github.com/max9159/agent-context-store): keep: Git-backed artifact/handoff format; schema permissiveness matters.
- [agentsentinel/agentsentinel-handoff](https://github.com/agentsentinel/agentsentinel-handoff): reject: handoff tooling lacks proven controller-owned Git transitions.
- [faresrafat3/agent-handoff](https://github.com/faresrafat3/agent-handoff): reject: handoff tooling lacks proven controller-owned Git transitions.
- [talocode/handofflane](https://github.com/talocode/handofflane): reject: semantic handoff validation but Git controller not shown.
- [moona3k/handoff](https://github.com/moona3k/handoff): reject: handoff tooling lacks proven controller-owned Git transitions.
- [roymcfarland/agent_handoff_guide](https://github.com/roymcfarland/agent_handoff_guide): reject: handoff tooling lacks proven controller-owned Git transitions.
- [shin4141/decision-os-v12-completion-integrity](https://github.com/shin4141/decision-os-v12-completion-integrity): reject: restartable handoffs, not every allowed transition.
- [statelyai/agent](https://github.com/statelyai/agent): reject: agent state/memory is not an authorized SDLC transition table.
- [MemMachine/MemMachine](https://github.com/MemMachine/MemMachine): reject: agent state/memory is not an authorized SDLC transition table.
- [statewright/statewright](https://github.com/statewright/statewright): reject: agent FSM guardrails, no project-Git authority stated.
- [adamterlson/AgenticStateMachines](https://github.com/adamterlson/AgenticStateMachines): reject: agent state/memory is not an authorized SDLC transition table.
- [jsz-05/LLM-State-Machine](https://github.com/jsz-05/LLM-State-Machine): reject: conversation FSM, not project delivery states.
- [Dicklesworthstone/frankenterm](https://github.com/Dicklesworthstone/frankenterm): reject: agent state/memory is not an authorized SDLC transition table.
- [mrgoonie/zuey-pi-setup](https://github.com/mrgoonie/zuey-pi-setup): reject: agent state/memory is not an authorized SDLC transition table.
- [gaorch85/Graph-of-States](https://github.com/gaorch85/Graph-of-States): reject: agent state/memory is not an authorized SDLC transition table.
- [robocorp/llmstatemachine](https://github.com/robocorp/llmstatemachine): reject: agent state/memory is not an authorized SDLC transition table.
- [mjunior/whatsapp-ai-pix-agent](https://github.com/mjunior/whatsapp-ai-pix-agent): reject: agent state/memory is not an authorized SDLC transition table.
- [jenseickmeyer/github-commit-status-bot](https://github.com/jenseickmeyer/github-commit-status-bot): reject: status notification alone does not prove four active checks.
- [chef-boneyard/dcob](https://github.com/chef-boneyard/dcob): reject: status notification alone does not prove four active checks.
- [Cadair/giles](https://github.com/Cadair/giles): reject: status notification alone does not prove four active checks.
- [funfair-tech/BuildBot](https://github.com/funfair-tech/BuildBot): reject: status notification alone does not prove four active checks.
- [arashstar1/bot-lua](https://github.com/arashstar1/bot-lua): reject: status notification alone does not prove four active checks.
- [rabbah/travis-to-slack](https://github.com/rabbah/travis-to-slack): reject: status notification alone does not prove four active checks.
- [eyobGere/git-bot](https://github.com/eyobGere/git-bot): reject: status notification alone does not prove four active checks.
- [tink-ab/four-eyes](https://github.com/tink-ab/four-eyes): reject: bors review count status; no specification/verification checks.
- [wetneb/hook-proxy](https://github.com/wetneb/hook-proxy): reject: status notification alone does not prove four active checks.
- [lydacious/GitLab-Jira-Webhook-Bot](https://github.com/lydacious/GitLab-Jira-Webhook-Bot): reject: status notification alone does not prove four active checks.
- [cloudposse/github-status-updater](https://github.com/cloudposse/github-status-updater): keep: Go status-posting API usage; not stdlib dependency.
- [BrianCArnold/RequireUserApproval](https://github.com/BrianCArnold/RequireUserApproval): reject: approval status only, not independent model check.
- [pkgdeps/automerge-gate](https://github.com/pkgdeps/automerge-gate): reject: merge enable gate, not full four-check protocol.
- [Nmishin/GitHub-PR-Status-Validator](https://github.com/Nmishin/GitHub-PR-Status-Validator): reject: post-merge detection cannot stop a bad merge.
- [Mustry-Academy/cicd-lab-03-github-actions](https://github.com/Mustry-Academy/cicd-lab-03-github-actions): reject: partial required-check facility; not the four controller gates.
- [etsy/combined-status-check-action](https://github.com/etsy/combined-status-check-action): reject: one aggregate check is not four separate statuses.
- [clusy-io/status-page](https://github.com/clusy-io/status-page): reject: partial required-check facility; not the four controller gates.
- [zammitjohn/required-status-checks](https://github.com/zammitjohn/required-status-checks): reject: partial required-check facility; not the four controller gates.
- [wechorg/required_status](https://github.com/wechorg/required_status): reject: partial required-check facility; not the four controller gates.
- [looshlee/merge_required_checks](https://github.com/looshlee/merge_required_checks): reject: partial required-check facility; not the four controller gates.
- [Spark-To-Paper-Skills/paperjury](https://github.com/Spark-To-Paper-Skills/paperjury): reject: review is not the independently enforced full task loop.
- [skainguyen1412/antigravity-superpowers](https://github.com/skainguyen1412/antigravity-superpowers): reject: review is not the independently enforced full task loop.
- [adshao/flounder](https://github.com/adshao/flounder): reject: review is not the independently enforced full task loop.
- [nurettincoban/ai-prd-workflow](https://github.com/nurettincoban/ai-prd-workflow): reject: review is not the independently enforced full task loop.
- [KimYx0207/Meta_Kim](https://github.com/KimYx0207/Meta_Kim): reject: review is not the independently enforced full task loop.
- [alibaba/aacr-bench](https://github.com/alibaba/aacr-bench): reject: review is not the independently enforced full task loop.
- [appautomaton/latex-arxiv-SKILL](https://github.com/appautomaton/latex-arxiv-SKILL): reject: review is not the independently enforced full task loop.
- [nexpeakcore/deepseek-harness-pr-review](https://github.com/nexpeakcore/deepseek-harness-pr-review): reject: review is not the independently enforced full task loop.
- [zhiyuwang720-dev/CodeSage](https://github.com/zhiyuwang720-dev/CodeSage): reject: review is not the independently enforced full task loop.
- [pjwang2022/Monstrare](https://github.com/pjwang2022/Monstrare): reject: review is not the independently enforced full task loop.
- [bors-ng/bors-ng](https://github.com/bors-ng/bors-ng): reject: auto-merge alone does not bind verification to exact SHA.
- [palantir/bulldozer](https://github.com/palantir/bulldozer): keep: auto-merge bot pattern; checked SHA still needs verification.
- [squalrus/merge-bot](https://github.com/squalrus/merge-bot): reject: auto-merge alone does not bind verification to exact SHA.
- [python/miss-islington](https://github.com/python/miss-islington): reject: auto-merge alone does not bind verification to exact SHA.
- [genuinetools/ghb0t](https://github.com/genuinetools/ghb0t): reject: auto-merge alone does not bind verification to exact SHA.
- [firstcontributions/merge-bot](https://github.com/firstcontributions/merge-bot): reject: auto-merge alone does not bind verification to exact SHA.
- [SnapInteractive/mergeatron](https://github.com/SnapInteractive/mergeatron): reject: auto-merge alone does not bind verification to exact SHA.
- [traefik/lobicornis](https://github.com/traefik/lobicornis): reject: auto-merge alone does not bind verification to exact SHA.
- [Schachte/Git-Enforcer](https://github.com/Schachte/Git-Enforcer): reject: auto-merge alone does not bind verification to exact SHA.
- [voyagegroup/popuko](https://github.com/voyagegroup/popuko): reject: auto-merge alone does not bind verification to exact SHA.
- [doorstop-dev/doorstop](https://github.com/doorstop-dev/doorstop): reject: trace-related but required source/criterion coverage not proven.
- [awslabs/duvet](https://github.com/awslabs/duvet): reject: trace-related but required source/criterion coverage not proven.
- [bmw-software-engineering/lobster](https://github.com/bmw-software-engineering/lobster): reject: trace-related but required source/criterion coverage not proven.
- [goeb/reqflow](https://github.com/goeb/reqflow): reject: trace-related but required source/criterion coverage not proven.
- [tobhey/finegrained-traceability](https://github.com/tobhey/finegrained-traceability): reject: trace-related but required source/criterion coverage not proven.
- [ociu/sphinx-traceability-extension](https://github.com/ociu/sphinx-traceability-extension): reject: trace-related but required source/criterion coverage not proven.
- [SAP-archive/quality-continuous-traceability-monitor](https://github.com/SAP-archive/quality-continuous-traceability-monitor): reject: trace-related but required source/criterion coverage not proven.
- [rtmx-ai/rtmx](https://github.com/rtmx-ai/rtmx): reject: trace-related but required source/criterion coverage not proven.
- [adrianpuiu/specification-document-generator](https://github.com/adrianpuiu/specification-document-generator): reject: trace-related but required source/criterion coverage not proven.
- [Ex-Origin/trace-syscall](https://github.com/Ex-Origin/trace-syscall): reject: trace-related but required source/criterion coverage not proven.
- [awslabs/duvet](https://github.com/awslabs/duvet): reject: trace-related but source/criterion coverage not proven.
- [goeb/reqflow](https://github.com/goeb/reqflow): reject: trace-related but source/criterion coverage not proven.
- [rtmx-ai/rtmx](https://github.com/rtmx-ai/rtmx): reject: trace-related but source/criterion coverage not proven.
- [konstantin-hatvan/traceability-tool](https://github.com/konstantin-hatvan/traceability-tool): reject: trace-related but source/criterion coverage not proven.
- [marquesarthur/traceability-tool](https://github.com/marquesarthur/traceability-tool): reject: trace-related but source/criterion coverage not proven.
- [lcdunstan/reqtrace](https://github.com/lcdunstan/reqtrace): reject: trace-related but source/criterion coverage not proven.
- [germanattanasio/traceability-assistant](https://github.com/germanattanasio/traceability-assistant): reject: trace-related but source/criterion coverage not proven.
- [lab9rinth/Traceability-Tool-Python](https://github.com/lab9rinth/Traceability-Tool-Python): reject: trace-related but source/criterion coverage not proven.
- [yesilzeytin/JanusTrace](https://github.com/yesilzeytin/JanusTrace): reject: trace-related but source/criterion coverage not proven.
- [simon3z/tmatrix](https://github.com/simon3z/tmatrix): reject: trace-related but source/criterion coverage not proven.
- [suzuking1192/samurai-agent](https://github.com/suzuking1192/samurai-agent): reject: prose synthesis is not deterministic full-trace validation.
- [wirelessr/kiro-workflow-prompts](https://github.com/wirelessr/kiro-workflow-prompts): reject: prose synthesis is not deterministic full-trace validation.
- [presidio-oss/specifai-mcp-server](https://github.com/presidio-oss/specifai-mcp-server): reject: prose synthesis is not deterministic full-trace validation.
- [hieubui2409/product-spec](https://github.com/hieubui2409/product-spec): reject: prose synthesis is not deterministic full-trace validation.
- [itential/builder-skills](https://github.com/itential/builder-skills): reject: prose synthesis is not deterministic full-trace validation.
- [ZTE-AICloud/Co-OmniSpec](https://github.com/ZTE-AICloud/Co-OmniSpec): reject: prose synthesis is not deterministic full-trace validation.
- [andrewstellman/quality-playbook](https://github.com/andrewstellman/quality-playbook): reject: prose synthesis is not deterministic full-trace validation.
- [mori-shin-x/artgraph](https://github.com/mori-shin-x/artgraph): reject: prose synthesis is not deterministic full-trace validation.
- [156181679-dev/prd-pilot](https://github.com/156181679-dev/prd-pilot): reject: prose synthesis is not deterministic full-trace validation.
- [LucivHuang/CodeSpec](https://github.com/LucivHuang/CodeSpec): reject: prose synthesis is not deterministic full-trace validation.
- [zmap/zlint](https://github.com/zmap/zlint): reject: linter does not prove every source need covered.
- [jonverrier/RequirementLinter](https://github.com/jonverrier/RequirementLinter): reject: AI critique, not provable requirement coverage.
- [python-coincidence/dep_checker](https://github.com/python-coincidence/dep_checker): reject: linter does not prove every source need covered.
- [brunocampos01/understanding-the-python-ecosystem](https://github.com/brunocampos01/understanding-the-python-ecosystem): reject: linter does not prove every source need covered.
- [numerodix/deplint](https://github.com/numerodix/deplint): reject: linter does not prove every source need covered.
- [almunnings/eslint-config](https://github.com/almunnings/eslint-config): reject: linter does not prove every source need covered.
- [atosatto/ansible-requirements-lint](https://github.com/atosatto/ansible-requirements-lint): reject: linter does not prove every source need covered.
- [kenderb/RCP_basic_SCSS_linter](https://github.com/kenderb/RCP_basic_SCSS_linter): reject: linter does not prove every source need covered.
- [salab/reqchecker](https://github.com/salab/reqchecker): reject: requirements lint; full text-to-ID coverage not shown.
- [walkerbrown/goodwrite](https://github.com/walkerbrown/goodwrite): reject: linter does not prove every source need covered.
- [debasisdwivedy/Coderama](https://github.com/debasisdwivedy/Coderama): reject: delivery workflow lacks demonstrated Git-only controller state.
- [synaptixs/spine](https://github.com/synaptixs/spine): reject: delivery workflow lacks demonstrated Git-only controller state.
- [digipulse-engineering/GAAI-framework](https://github.com/digipulse-engineering/GAAI-framework): reject: delivery workflow lacks demonstrated Git-only controller state.
- [q3ok/coordinated-agent-team](https://github.com/q3ok/coordinated-agent-team): reject: delivery workflow lacks demonstrated Git-only controller state.
- [cmAIdx/headless-claude-automation-template](https://github.com/cmAIdx/headless-claude-automation-template): reject: delivery workflow lacks demonstrated Git-only controller state.
- [Free-devloper/autonomous-software-delivery-organization](https://github.com/Free-devloper/autonomous-software-delivery-organization): reject: delivery workflow lacks demonstrated Git-only controller state.
- [eduardocerqueira/ai-alpha-squad](https://github.com/eduardocerqueira/ai-alpha-squad): reject: delivery workflow lacks demonstrated Git-only controller state.
- [launchapp-dev/ao](https://github.com/launchapp-dev/ao): reject: delivery workflow lacks demonstrated Git-only controller state.
- [seanyao/roll](https://github.com/seanyao/roll): reject: delivery workflow lacks demonstrated Git-only controller state.
- [jonathanpopham/estate-agent](https://github.com/jonathanpopham/estate-agent): reject: guarded issue delivery; Git-only state not shown.
- [initialcommit-com/git-sim](https://github.com/initialcommit-com/git-sim): reject: Git records do not establish protected controller-owned transitions.
- [entireio/cli](https://github.com/entireio/cli): reject: Git records do not establish protected controller-owned transitions.
- [Ian-stetsenko/lore-protocol](https://github.com/Ian-stetsenko/lore-protocol): reject: Git records do not establish protected controller-owned transitions.
- [codeprakhar25/agentdiff](https://github.com/codeprakhar25/agentdiff): reject: Git records do not establish protected controller-owned transitions.
- [mbeacom/adrkit](https://github.com/mbeacom/adrkit): keep: typed Git ADR and deterministic evaluation pattern.
- [pinsky-three/vibe-graph](https://github.com/pinsky-three/vibe-graph): reject: Git records do not establish protected controller-owned transitions.
- [kerbelp/metatron](https://github.com/kerbelp/metatron): reject: Git records do not establish protected controller-owned transitions.
- [crabbuild/trail](https://github.com/crabbuild/trail): reject: Git records do not establish protected controller-owned transitions.
- [termaxa/termaxa](https://github.com/termaxa/termaxa): reject: Git records do not establish protected controller-owned transitions.
- [yunaremaia/agent-undo](https://github.com/yunaremaia/agent-undo): reject: Git records do not establish protected controller-owned transitions.
- [kaplanelad/shellfirm](https://github.com/kaplanelad/shellfirm): reject: local guardrail cannot enforce remote Git protection alone.
- [clay-good/OpenLore](https://github.com/clay-good/OpenLore): reject: local guardrail cannot enforce remote Git protection alone.
- [SumonMSelim/agentguard](https://github.com/SumonMSelim/agentguard): reject: local guardrail cannot enforce remote Git protection alone.
- [KbWen/agentic-os](https://github.com/KbWen/agentic-os): reject: local guardrail cannot enforce remote Git protection alone.
- [Dicklesworthstone/ultimate_bug_scanner](https://github.com/Dicklesworthstone/ultimate_bug_scanner): reject: local guardrail cannot enforce remote Git protection alone.
- [LuD1161/agentjail](https://github.com/LuD1161/agentjail): keep: tool-call policy pattern, not branch protection.
- [melihmucuk/leash](https://github.com/melihmucuk/leash): reject: local guardrail cannot enforce remote Git protection alone.
- [aws-samples/sample-amazon-bedrock-for-beginners](https://github.com/aws-samples/sample-amazon-bedrock-for-beginners): reject: local guardrail cannot enforce remote Git protection alone.
- [chainloop-dev/chainloop](https://github.com/chainloop-dev/chainloop): reject: signed signals stored outside project Git.
- [stickerdaniel/saas-starter](https://github.com/stickerdaniel/saas-starter): reject: local guardrail cannot enforce remote Git protection alone.
- [awslabs/aws-config-engine-for-compliance-as-code](https://github.com/awslabs/aws-config-engine-for-compliance-as-code): reject: settings automation requires admin privilege absent from App.
- [Vivswan/github-settings-as-code](https://github.com/Vivswan/github-settings-as-code): reject: declarative rules useful but App cannot administer them.
- [theztefan/rulesets-as-code-action](https://github.com/theztefan/rulesets-as-code-action): reject: ruleset action needs privileges controller lacks.
- [nattcharoen-beep/agentic-ceo-protocol](https://github.com/nattcharoen-beep/agentic-ceo-protocol): reject: settings automation requires admin privilege absent from App.
- [CheckPointSW-Community/cloud-contracts-as-code](https://github.com/CheckPointSW-Community/cloud-contracts-as-code): reject: settings automation requires admin privilege absent from App.
- [noshiro-pf/github-settings-as-code](https://github.com/noshiro-pf/github-settings-as-code): reject: settings backup/apply needs admin permissions.
- [AlexAsAService/grac](https://github.com/AlexAsAService/grac): reject: GitHub ruleset automation needs admin permissions.
- [ucsb-cs156/rulesets](https://github.com/ucsb-cs156/rulesets): reject: settings automation requires admin privilege absent from App.
- [alvie-devsecops-lab/github-governance](https://github.com/alvie-devsecops-lab/github-governance): reject: settings automation requires admin privilege absent from App.
- [carpet-stain/infra](https://github.com/carpet-stain/infra): reject: settings automation requires admin privilege absent from App.
- [asyncapi-archived-repos/shape-up-process](https://github.com/asyncapi-archived-repos/shape-up-process): reject: bet/shaping content is not validated answer copied into Git.
- [rhyslawsn/shape-up-skills](https://github.com/rhyslawsn/shape-up-skills): reject: bet/shaping content is not validated answer copied into Git.
- [ramboz/shaper](https://github.com/ramboz/shaper): reject: bet/shaping content is not validated answer copied into Git.
- [liz02may-commits/How-Professional-Land-Development-Services-Shape-Better-Projects-From-the-Ground-Up](https://github.com/liz02may-commits/How-Professional-Land-Development-Services-Shape-Better-Projects-From-the-Ground-Up): reject: bet/shaping content is not validated answer copied into Git.
- [albertocruzdg/beto-learns-modern-js](https://github.com/albertocruzdg/beto-learns-modern-js): reject: bet/shaping content is not validated answer copied into Git.
- [florianilch/mermaid-hillchart](https://github.com/florianilch/mermaid-hillchart): reject: bet/shaping content is not validated answer copied into Git.
- [rafae2k/rafoflow](https://github.com/rafae2k/rafoflow): reject: bet/shaping content is not validated answer copied into Git.
- [medvelaszloka/ai-native-pm-os](https://github.com/medvelaszloka/ai-native-pm-os): reject: bet/shaping content is not validated answer copied into Git.
- [danderfer/Comp_Sci_Sem_2](https://github.com/danderfer/Comp_Sci_Sem_2): reject: bet/shaping content is not validated answer copied into Git.
- [rezwanuldp/Weight-Loss](https://github.com/rezwanuldp/Weight-Loss): reject: bet/shaping content is not validated answer copied into Git.
- [dohooo/helmor](https://github.com/dohooo/helmor): reject: multi-agent coding lacks required checked role handoffs.
- [kuafuai/DevOpsGPT](https://github.com/kuafuai/DevOpsGPT): reject: multi-agent coding lacks required checked role handoffs.
- [codefuse-ai/codefuse-chatbot](https://github.com/codefuse-ai/codefuse-chatbot): reject: multi-agent coding lacks required checked role handoffs.
- [untra/operator](https://github.com/untra/operator): reject: multi-agent coding lacks required checked role handoffs.
- [hoangsonww/AI-Agents-Orchestrator](https://github.com/hoangsonww/AI-Agents-Orchestrator): reject: multi-agent coding lacks required checked role handoffs.
- [jonnyzzz/run-agent](https://github.com/jonnyzzz/run-agent): reject: multi-agent coding lacks required checked role handoffs.
- [path-to-ai-com/just-ship](https://github.com/path-to-ai-com/just-ship): reject: multi-agent coding lacks required checked role handoffs.
- [LeviXDD/AgentDeskSkill](https://github.com/LeviXDD/AgentDeskSkill): reject: multi-agent coding lacks required checked role handoffs.
- [crewswarm/crewswarm](https://github.com/crewswarm/crewswarm): reject: multi-agent coding lacks required checked role handoffs.
- [pridiuksson/cursor-agents](https://github.com/pridiuksson/cursor-agents): reject: multi-agent coding lacks required checked role handoffs.
- [20AhmedRamadan04/Instagram_Scraper](https://github.com/20AhmedRamadan04/Instagram_Scraper): reject: comment handling lacks hashed typed handoff into Git.
- [dimitristheod01/Instagram-comment-data-cleaning](https://github.com/dimitristheod01/Instagram-comment-data-cleaning): reject: comment handling lacks hashed typed handoff into Git.
- [Kavyashree1456/ai-github-issue-assistant](https://github.com/Kavyashree1456/ai-github-issue-assistant): reject: comment handling lacks hashed typed handoff into Git.
- [iamAntimPal/Python-Programs](https://github.com/iamAntimPal/Python-Programs): reject: comment handling lacks hashed typed handoff into Git.
- [ankitsablok89/Stanford_Online_Design_and_Analysis_of_Algorithms_Part_2](https://github.com/ankitsablok89/Stanford_Online_Design_and_Analysis_of_Algorithms_Part_2): reject: comment handling lacks hashed typed handoff into Git.
- [a1yak/Comments-to-Tickets](https://github.com/a1yak/Comments-to-Tickets): reject: comment handling lacks hashed typed handoff into Git.
- [SelinaY123/OpinionRadar](https://github.com/SelinaY123/OpinionRadar): reject: comment handling lacks hashed typed handoff into Git.
- [DaiaPrado/llm-powered_market_trend_discovery](https://github.com/DaiaPrado/llm-powered_market_trend_discovery): reject: comment handling lacks hashed typed handoff into Git.
- [alexmaunu/collaborative-debug-loop](https://github.com/alexmaunu/collaborative-debug-loop): reject: comment handling lacks hashed typed handoff into Git.
- [AngelineVenicia/Automated-Code-Review-System-for-Github](https://github.com/AngelineVenicia/Automated-Code-Review-System-for-Github): reject: comment handling lacks hashed typed handoff into Git.
- [In-sp3ctr3/memPR](https://github.com/In-sp3ctr3/memPR): reject: PR provenance alone does not gate exact checked head.
- [twilio-labs/gordon](https://github.com/twilio-labs/gordon): reject: PR provenance alone does not gate exact checked head.
- [shapor/gitiq](https://github.com/shapor/gitiq): reject: PR provenance alone does not gate exact checked head.
- [semanticash/cli](https://github.com/semanticash/cli): reject: PR provenance alone does not gate exact checked head.
- [Dev-31/provenance-layer2](https://github.com/Dev-31/provenance-layer2): reject: PR provenance alone does not gate exact checked head.
- [m1insights/provenance](https://github.com/m1insights/provenance): reject: PR provenance alone does not gate exact checked head.
- [Kartik24Hulmukh/jittest](https://github.com/Kartik24Hulmukh/jittest): reject: PR provenance alone does not gate exact checked head.
- [jhinpan/ROCmKernelWiki](https://github.com/jhinpan/ROCmKernelWiki): reject: PR provenance alone does not gate exact checked head.
- [ryanjmichie-git/forgeproof-verify](https://github.com/ryanjmichie-git/forgeproof-verify): reject: PR provenance alone does not gate exact checked head.
- [P-r-e-m-i-u-m/ai-change-passport](https://github.com/P-r-e-m-i-u-m/ai-change-passport): reject: PR provenance alone does not gate exact checked head.
- [PerryLink/dsh-github](https://github.com/PerryLink/dsh-github): reject: generic approvals do not validate/copy one owner answer.
- [urstrulyemkay/ig-comment-dm-bot](https://github.com/urstrulyemkay/ig-comment-dm-bot): reject: generic approvals do not validate/copy one owner answer.
- [arturo-c/github_bot_approvals](https://github.com/arturo-c/github_bot_approvals): reject: generic approvals do not validate/copy one owner answer.
- [aayushcyy/slackBot](https://github.com/aayushcyy/slackBot): reject: generic approvals do not validate/copy one owner answer.
- [theimperious1/Reddit_Approval_Bot](https://github.com/theimperious1/Reddit_Approval_Bot): reject: generic approvals do not validate/copy one owner answer.
- [samuelemb/telegram-confession-bot](https://github.com/samuelemb/telegram-confession-bot): reject: generic approvals do not validate/copy one owner answer.
- [ArulAgnes/POST-COMMENTS-Social-Media-Automation-Platform-Bot](https://github.com/ArulAgnes/POST-COMMENTS-Social-Media-Automation-Platform-Bot): reject: generic approvals do not validate/copy one owner answer.
- [nootkan/ossn_account_approval](https://github.com/nootkan/ossn_account_approval): reject: generic approvals do not validate/copy one owner answer.
- [Zyrenkam/Anonymous_tg_bot](https://github.com/Zyrenkam/Anonymous_tg_bot): reject: generic approvals do not validate/copy one owner answer.
- [ranamahmod/rana-engagement-reply-bot](https://github.com/ranamahmod/rana-engagement-reply-bot): reject: generic approvals do not validate/copy one owner answer.
- [krishagarwal314/CodeJury](https://github.com/krishagarwal314/CodeJury): reject: multi-agent delivery but exact four-status enforcement not shown.
- [tiammomo/RepoSteward](https://github.com/tiammomo/RepoSteward): reject: policy gated issue-to-PR; protected Git records not shown.
- [ktwu01/agent-pr-gate](https://github.com/ktwu01/agent-pr-gate): reject: partial PR gate, not four mandatory active statuses.
- [caique-neves-ferreira/pr-review-agent](https://github.com/caique-neves-ferreira/pr-review-agent): reject: partial PR gate, not four mandatory active statuses.
- [maida-ai/maida-assert](https://github.com/maida-ai/maida-assert): reject: partial PR gate, not four mandatory active statuses.
- [xdlc-labs/airlock](https://github.com/xdlc-labs/airlock): reject: checks agent configuration drift, not complete task gates.
- [hamc/blastproof](https://github.com/hamc/blastproof): reject: partial PR gate, not four mandatory active statuses.
- [brett-buskirk/agent-gate](https://github.com/brett-buskirk/agent-gate): reject: partial PR gate, not four mandatory active statuses.
- [mojunrubest/agent-delivery-gate](https://github.com/mojunrubest/agent-delivery-gate): reject: host gate, not stated four contexts or checked-SHA merge.
- [FreePeak/devagent](https://github.com/FreePeak/devagent): reject: partial PR gate, not four mandatory active statuses.
- [jangtrinh/design-os-code2flow](https://github.com/jangtrinh/design-os-code2flow): reject: document lint is not total requirement/source coverage.
- [sohaibt/prdlinter](https://github.com/sohaibt/prdlinter): reject: document lint is not total requirement/source coverage.
- [rock3r/publish-spec](https://github.com/rock3r/publish-spec): reject: document lint is not total requirement/source coverage.
- [Abhishma/roadmap-linter](https://github.com/Abhishma/roadmap-linter): reject: document lint is not total requirement/source coverage.
- [DreamStars1/WenLint](https://github.com/DreamStars1/WenLint): reject: document lint is not total requirement/source coverage.
- [murtuzabuilds/prd-to-jira](https://github.com/murtuzabuilds/prd-to-jira): reject: document lint is not total requirement/source coverage.
- [Jaimin4727/prd-writer](https://github.com/Jaimin4727/prd-writer): reject: document lint is not total requirement/source coverage.
- [j0yen/atlas](https://github.com/j0yen/atlas): reject: document lint is not total requirement/source coverage.
- [nilukush/prd-as-a-code](https://github.com/nilukush/prd-as-a-code): reject: Git-native schema format; cannot prove uncovered prose needs.
- [shaikbasha1437/speclint](https://github.com/shaikbasha1437/speclint): reject: document lint is not total requirement/source coverage.
- [Fission-AI/OpenSpec](https://github.com/Fission-AI/OpenSpec): keep: specification workflow method, not mechanical trace proof.
- [https-deeplearning-ai/sc-spec-driven-development-files](https://github.com/https-deeplearning-ai/sc-spec-driven-development-files): reject: spec workflow not a fail-closed traceability checker.
- [gsd-build/get-shit-done](https://github.com/gsd-build/get-shit-done): reject: spec workflow not a fail-closed traceability checker.
- [Engineering4AI/awesome-spec-driven-development](https://github.com/Engineering4AI/awesome-spec-driven-development): reject: spec workflow not a fail-closed traceability checker.
- [buildermethods/agent-os](https://github.com/buildermethods/agent-os): reject: spec workflow not a fail-closed traceability checker.
- [formulahendry/mcp-server-spec-driven-development](https://github.com/formulahendry/mcp-server-spec-driven-development): reject: spec workflow not a fail-closed traceability checker.
- [zhu1090093659/spec_driven_develop](https://github.com/zhu1090093659/spec_driven_develop): reject: spec workflow not a fail-closed traceability checker.
- [mouredev/hello-sdd](https://github.com/mouredev/hello-sdd): reject: spec workflow not a fail-closed traceability checker.
- [marcelsud/spec-driven-agentic-development](https://github.com/marcelsud/spec-driven-agentic-development): reject: spec workflow not a fail-closed traceability checker.
- [Pimzino/spec-workflow-mcp](https://github.com/Pimzino/spec-workflow-mcp): reject: spec workflow not a fail-closed traceability checker.
- [oceanbase/powercontext](https://github.com/oceanbase/powercontext): reject: task memory/state lacks protected project-Git authority.
- [MaxHu-xuan/task-state-guard](https://github.com/MaxHu-xuan/task-state-guard): reject: task memory/state lacks protected project-Git authority.
- [AMAP-ML/LongHorizon-Harness](https://github.com/AMAP-ML/LongHorizon-Harness): reject: task memory/state lacks protected project-Git authority.
- [henryalps/OpenManus](https://github.com/henryalps/OpenManus): reject: task memory/state lacks protected project-Git authority.
- [thu-nmrc/OpenHarness](https://github.com/thu-nmrc/OpenHarness): reject: task memory/state lacks protected project-Git authority.
- [agentrebench/AgentRE-Bench](https://github.com/agentrebench/AgentRE-Bench): reject: task memory/state lacks protected project-Git authority.
- [manish-raana/openclaw-mission-control](https://github.com/manish-raana/openclaw-mission-control): reject: task memory/state lacks protected project-Git authority.
- [pengyue-polaron/deepseek-harness-genui](https://github.com/pengyue-polaron/deepseek-harness-genui): reject: task memory/state lacks protected project-Git authority.
- [svcvit/dify-plugin-tod_agent](https://github.com/svcvit/dify-plugin-tod_agent): reject: task memory/state lacks protected project-Git authority.
- [AndreaGriffiths11/agent-context-system](https://github.com/AndreaGriffiths11/agent-context-system): reject: task memory/state lacks protected project-Git authority.

## What no candidate covers

- No source establishes the *whole* sequence: issue -> plan -> review by a different model -> development -> draft PR -> automated gates -> different-harness verifier -> close-out commit -> check **that** final head -> conditional merge of precisely that SHA.
- No candidate supplies this project's closed transition table plus code-owned validity for *every* role handoff, including a forge-comment-only artifact identified by immutable location and content hash. Git-backed handoff files alone do not prove that the controller, rather than the agent, authorized a transition.
- No reviewed tool guarantees all four separately named `layup/` contexts (gates, specification, rules, verification) are **active**, fresh and successful on the exact required commit before marking ready. GitHub accepts `neutral`/`skipped` for required checks; a green branch rule alone cannot enforce F-0001#5. A source app restriction still needs administrator configuration.
- No candidate enforces the stated architecture-choice rule: discard a missing or failed row and an unevidenced pass, then require exactly one survivor. ADR scoring and ADR lint are useful but different contracts.
- No candidate implements the owner-facing brief's exact mandatory parts and length, binds exactly one owner's answering comment, then copies and hashes that answer into the project's Git.
- No candidate simultaneously fails on an unresolved requirement ID, a source need without any requirement, a requirement with no acceptance criterion, and a delivered requirement without a specification section. `reqmd` documents missing traces and coverage as *warnings* in some cases, so its default green exit is not this contract.
- No reviewed tool demonstrates a protected controller-only Git branch containing *all* controller records (task register, transition ledger, events, gate outcomes and parameters). Files in Git and protected Git refs are different guarantees; an App with no administration permission cannot establish its own rules.
- No source was found for automatic assurance that the demo's rendered specification retains both MoSCoW and Phase columns through the first bet; rendered-format and gate requirements remain project-specific.

## Limits

- The five curated lists are discovery indexes at the commits in `sources/lists.tsv`; index descriptions and the supplied repository-search snapshot are not proof of behavior. The entries are screened by what their descriptions say, not by execution or installation.
- The long indexes and some web-result displays were truncated by the terminal/search renderer; consequently the screening of entries outside the named lines may be incomplete. The per-repository documentation checks are limited to linked pages; no source code, tests, workflows, release binaries, or live branch settings were executed or audited.
- `sources/gh-search.tsv` reports pushed dates, not releases. For candidates absent from that file, activity is recorded as "not stated" unless the linked release supplied a date. Licenses marked "search snapshot" were not independently verified against a license file. `reqmd.dev` and `voedger/reqmd` are distinct projects despite their shared name.
- The cited GitHub docs describe platform capabilities, not a test of this particular App's permissions, installed status-source selection, repository ruleset, merge method, or actual status checks. Their permission/configuration constraints require administrator setup outside this search.
- No test repository, GitHub API mutations, candidate code, installations or external service state were created or run. The search does not establish a complete implementation or imply that an uninspected alternative is unsuitable.
