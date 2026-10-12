# T-msrd — the search of M2d, searcher B, word for word

Searcher: Claude Sonnet 5.5 (`claude-sonnet-5-5`) on Claude Code 2.1.296 (`claude -p --permission-mode dontAsk`, read-only tools, web search and fetch, writes only under `out/` by its absolute path; `--setting-sources project,local`, `--no-session-persistence`, `ANTHROPIC_DEFAULT_HAIKU_MODEL=claude-sonnet-5-5`), a fresh session, 2026-10-11T14:40:41Z to 2026-10-11T14:43:32Z (2 min 49 s by its result event), exit 0, on its first try on this harness; USD 1.61; `modelUsage`: `claude-sonnet-5-5` only, 1,733,482 tokens; no permission denied. A model other than the author's, of neither tier list (`docs/engineering-discipline.md`, Model tiers), as for `T-g49f`: OpenCode, tried first, failed with "Invalid credential" (2026-10-11 at 14:40Z). Brief: [`brief.md`](brief.md), with the five lists of [`sources.md`](sources.md) and [`gh-search.tsv`](gh-search.tsv). The output follows word for word inside a fenced block; its paths are relative to the searcher's directory.

~~~~text
# Public-solution search for M2d (task T-msrd), searcher B

Date: 2026-10-11. Scope: judge existing public solutions; no design.
Rule used for every value below: "not stated" means the source I read did not say it;
"not checked" means I did not open that source. Repository metadata (license, language,
last push) was read from `api.github.com/repos/...` on 2026-10-11. Metadata for entries
I did not fetch is taken from `sources/gh-search.tsv` and marked "(tsv)".

## Queries

My own queries first, then the fetches that checked candidates. "Words" = the fact words
the query comes from.

| # | Query / fetch | Words it comes from | What it found |
|---|---|---|---|
| Q1 | web: GitHub repository activity REST API list repository activities branch push force_push actor | F-0003#64 audit of writes; the activity-log read of the brief | `GET /repos/{owner}/{repo}/activity`; filters `ref`, `actor`, `activity_type`; third-party mirrors only, so I fetched the official page (Q9) |
| Q2 | web: GitHub rulesets REST API get rules for a branch rules/branches required_status_checks integration_id | F-0003#43; "required checks pinned to the app that posts them"; "effective rules read back" | official REST "rules" page; `GET /repos/{o}/{r}/rules/branches/{branch}`; `integration_id` seen in Q8 |
| Q3 | web: GH013 Repository rule violations found push declined | "probe push refused for a rule violation (any other failure is not a pass)" | only low-quality pages; no official page; refusal text not confirmed from an official source |
| Q4 | web: GitHub App permission required to create or update repository rulesets Administration | F-0004#1 token has no administration; "ruleset set by a human with admin rights" | summary of GHES docs: create ruleset needs Administration (write); read needs Metadata (read). Official github.com page did not state it (Q8) |
| Q5 | web: copier answers file .copier-answers.yml _commit pinned template version copier update | F-0001#8 pinned recorded version; "second run on a repo set up by an earlier version"; F-0003#42 | `_commit` and `_src_path` in answers file; `copier update` does a 3-way merge; `--vcs-ref` pins |
| Q6 | web: ossf allstar GitHub App branch protection policy enforcement permissions | F-0003#43; F-0001#3 | Allstar: Go GitHub App, policy YAML, actions log / issue / fix |
| Q7 | fetch api.github.com for 13 repositories (see Candidates) | step 3 of the brief | license, language, pushed_at, archived flag |
| Q8 | fetch docs.github.com/en/rest/repos/rules | F-0003#43 "no bypass for the controller", "pinned to the app" | `bypass_actors[].actor_type` includes `Integration`; `bypass_mode` always/pull_request/exempt; `required_status_checks[].integration_id` optional; rules/branches takes no wildcard and omits evaluate/disabled rulesets; permission for create not stated on page |
| Q9 | fetch docs.github.com/en/rest/repos/repos#list-repository-activities | F-0003#64 | six `activity_type` values, `actor` object, cursor paging, per_page max 100; permissions not stated |

Queries not run: I did not run web searches on the Go-only, git-notes, or evidence-per-value
themes (F-0001#4, F-0001#1). The tsv already has queries for those (all off-target, see
"The search file"). This is a gap, listed under Limits.

## Candidates

Columns: name and URL; license; language; last activity; part of the problem; state in Git;
deterministic; keep as; why.

| Name / URL | License | Language | Last activity | Part of the problem | State in Git | Deterministic | Keep as | Why |
|---|---|---|---|---|---|---|---|---|
| GitHub REST "rules" API, https://docs.github.com/en/rest/repos/rules | docs, not stated | not applicable (HTTP) | page current on 2026-10-11 | ruleset file the program writes; read-back of effective rules (F-0003#43) | not applicable | read calls are; result depends on forge state | **API usage** | `GET /repos/{o}/{r}/rules/branches/{branch}` returns the active rules for a branch, even if the branch does not exist. Create/update body has `bypass_actors` (`actor_type` Integration, OrganizationAdmin, RepositoryRole, Team, DeployKey, User) and `required_status_checks[]` with optional `integration_id`. The program can write that JSON; a human applies it. |
| Same page, known limits | | | | | | | **known limit** | Rules endpoint: no wildcard in the branch name (GraphQL for wildcards); rulesets in `evaluate` or `disabled` are not returned, so a read-back that finds the rule proves "active" only. `integration_id` is optional, so a ruleset without it does not pin the app: the read-back must compare it. `bypass_mode: exempt` leaves no bypass audit entry. |
| Same endpoint, permission | | | | | | | **known limit** | Create ruleset needs Administration (write) per a summary of GHES docs (Q4); the official page I fetched did not state it. This fits F-0004#1: the controller token cannot apply rules, so the human step is forced, not optional. Reading needs Metadata (read). |
| GitHub REST "List repository activities", https://docs.github.com/en/rest/repos/repos#list-repository-activities | docs, not stated | not applicable (HTTP) | page current on 2026-10-11 | audit of who updated each branch (F-0003#64) | not applicable | read is; log content is forge state | **API usage** | `GET /repos/{o}/{r}/activity` with `ref`, `actor`, `activity_type` (push, force_push, branch_creation, branch_deletion, pr_merge, merge_queue_merge), each item has `before`, `after`, `ref`, `timestamp`, `actor`. Fits "read the activity log to audit who updated each branch". Permissions and retention not stated on the page. |
| copier, https://github.com/copier-org/copier | MIT | Python | pushed 2026-10-08, not archived | setup from pinned baseline commit; second run on an older setup (F-0001#8, F-0003#42, F-0003#37) | answers file `.copier-answers.yml` in the project (`_src_path`, `_commit`, answers) | rendering is; update uses 3-way merge, conflicts need a person | **pattern** (cannot be a dependency: Python) | Records the template source and version next to the answers, and updates by 3-way merge between old template, new template, and the project. Copy the idea of a recorded-version file and a diff-based second run. Note `_tasks` need `--trust`: a template can run commands, which the brief's controller must not allow. |
| actions-ext/copier-update, https://github.com/actions-ext/copier-update | Apache-2.0 (tsv) | TypeScript (tsv) | pushed 2026-10-11 (tsv) | second run | not checked | not checked | **pattern** (low) | Shows an update run as a branch plus pull request, not as a direct push. Not read beyond tsv. |
| Vivswan/github-settings-as-code, https://github.com/Vivswan/github-settings-as-code | NOASSERTION (key "other") | TypeScript | pushed 2026-10-11, not archived | ruleset file applied from the repository (F-0003#43) | desired settings in the repo file; no state file stated | "stateless" by its description | **pattern** | Declarative file of rulesets, labels, branch protection, applied by an Action. Cannot be a dependency (TypeScript, unclear license). Useful for the file's field list and for "stateless apply". It applies with a token that has admin rights, which the controller does not have. |
| kuhlman-labs/repo-ruleset-bot, https://github.com/kuhlman-labs/repo-ruleset-bot | MIT | Go | pushed 2025-04-21, not archived | syncs org rulesets from a source; Go-language example of calling the rulesets API | not checked | not checked | **pattern** | Same language as the controller: a reference for ruleset API calls in Go. It is a Go module or App, so not a dependency. I did not read its code. It is an org-level tool; the brief's rules are per repository. |
| katiem0/gh-migrate-rulesets, https://github.com/katiem0/gh-migrate-rulesets | MIT | Go | pushed 2026-08-04, not archived | export/create of rulesets as files (F-0003#43) | not checked | not checked | **format** (low) | A `gh` extension that creates and reports rulesets from a file; its file layout is a worked example of a ruleset on disk. Not read beyond the description. |
| integrations/terraform-provider-github, https://github.com/integrations/terraform-provider-github | MIT | Go | pushed 2026-10-11, not archived | `github_repository_ruleset` resource field names (F-0003#43) | not applicable | plan/apply is, by design | **format** (reference only) | Not a dependency and not run. Its resource schema is a second public list of ruleset fields to check the program's file against. I did not open the schema page. |
| ossf/allstar, https://github.com/ossf/allstar | Apache-2.0 | Go | pushed 2026-10-11, not archived | continuous check and fix of repo settings (F-0001#3, F-0003#43) | policy YAML in a repo; no run state stated | not stated | **pattern** | Go GitHub App that compares settings with a policy and logs, opens an issue, or reverts. Its actions "log / issue / fix" map to read-back and report. It asks for write on issues and checks only, per the OpenSSF post: shows that an App can audit without administration, but then cannot fix. It is a hosted continuous service; the brief wants a one-time proof. |
| neul-labs/grite, https://github.com/neul-labs/grite | MIT | Rust | pushed 2026-07-02, not archived | event log kept in git refs, no server (F-0001#1) | yes, `refs/grite/wal`, per the tsv-style list text | claimed ("converge deterministically") | **pattern** (low) | Shows an append-only log in a git ref with `git push`. The brief wants the evidence on a branch that only the controller identity can update; a custom ref is not a branch, so branch rulesets would not cover it (my reading; the ruleset docs target branches and tags). Rust, not usable. |
| kusaridev/skootrs, https://github.com/kusaridev/skootrs | Apache-2.0 | Rust | pushed 2024-07-29, **archived** | creating secure repos with branch protection built in | not stated | not stated | **pattern** (low) | The closest "create a new repo from best practice" tool in the lists. Archived since 2024, so no maintenance. Use only to confirm that the gap is old. |
| basilisk-labs/agentplane, https://github.com/basilisk-labs/agentplane | MIT | TypeScript | pushed 2026-10-11, not archived | agent workflow with state in the repo (F-0001#1) | yes, `.agentplane/` inside the repo (list text) | not stated | **pattern** (low) | Task, plan, approve, verify, finish with all state in the repo. Matches "Git is the record" but has no ruleset or proof from outside. Not read beyond list and description. |
| github/gh-aw, https://github.com/github/gh-aw | MIT (list) | not checked | not checked (list stars 5.3k) | agent writes only through validated outputs (free-text parts checked by code) | not checked | not checked | **pattern** (low) | The list says agent jobs are read-only and writes pass "validated safe-outputs": the same shape as "free text written by an agent, each part checked by code". Not opened. |

Not kept although they were close, so the other searcher can compare: chainloop (state on a
server, not in project Git), h5i (the repository now describes a different product, see list
section), Terraform org-governance repos (HCL, need admin token and a state file outside Git).

## The lists

Format: `name: keep|reject: reason`. Only entries that touch the problem. Five lists, read in full.

### hysnsec/awesome-policy-as-code (af4e6354067b)

- OPA: reject: general policy engine; Go module cannot be a dependency; no GitHub ruleset support stated
- Sentinel: reject: HashiCorp-embedded, closed policy language
- Regula: reject: Terraform/CloudFormation checks only
- Intercept: reject: static analysis auditing, no forge rules
- Checkov: reject: IaC scanner, Python
- Terrascan: reject: IaC scanner
- kics: reject: IaC scanner
- Gatekeeper / Kyverno / Konstraint / kube-mgmt / MagTape: reject: Kubernetes admission only
- Cloud Custodian: reject: cloud resources, not git repos
- "Managing Azure Policy as Code with GitHub": reject: Azure, though it shows policy files kept in Git
- Blogs and videos sections: reject: no source text read that covers GitHub rulesets

### weaveworks/awesome-gitops (82af5b15f445)

- Atlantis: reject: Terraform PR automation; holds admin credentials; state outside project Git
- Flux, ArgoCD, Autoapply, PipeCD, Werf: reject: Kubernetes reconcile; same idea "Git is the record" but nothing for rulesets
- GitOps.tech / "Operations by pull request": reject: background reading only
- git-secret, SOPS: reject: secrets, not in the problem

### sdras/awesome-actions (0edccf33480b)

- Branch Protection Bot (benjefferies): reject: it temporarily turns off "Include administrators" in branch protection, the opposite of "no bypass"; an example of the bypass pattern the rules must refuse
- talos-systems/conform: reject: policy on commit messages and similar; Go but commit-style only; not checked further
- mszostok/codeowners-validator: reject: validates CODEOWNERS file; useful only if the program writes CODEOWNERS (not in the brief)
- Terraform GitHub administration example (asgharlabs/github-terraform): reject: Terraform, admin token, state outside Git
- cirrus-actions/branch-guard: reject: blocks merges on target-branch checks; not about app-pinned checks
- tibdex/autosquash, pascalgn/automerge-action: reject: merge automation; do not set rules
- WyriHaximus wait-for-status: reject: wait for checks; not a rule proof
- jpmcb/prow-github-actions: reject: ChatOps policy; no rule setting
- Copybara Action: reject: code move between repos, not a pinned-baseline setup
- Jinja2 action, Edit JSON File, YAML/JSON converter: reject: template or edit steps without source records
- Actions "Official Resources" / "Create your actions": reject: not the problem (the controller is not an Action)
- Everything else in this list (deploy, notifications, docker, k8s, aws, linting, testing, pages): reject: off problem

### bureado/awesome-software-supply-chain-security (31252e6810c4)

- ossf/allstar: keep: see Candidates
- kusaridev/skootrs: keep: see Candidates (archived)
- ossf/scorecard and scorecard-action: reject: scores repo settings read-only for supply-chain risk; its branch-protection check needs a token with admin read for some settings (not checked); no rule setting; useful only as a comparison of checks
- reposaur/reposaur: reject: "compliance tool for development platforms"; language and license not checked; policy engine style
- Legit-Labs/legitify: reject: detects GitHub misconfigurations; read-only report; not a proof by probe push; license/language not checked
- crashappsec/github-analyzer: reject: checks org security settings; org scope
- scribe-public/gitgat: reject: posture evaluation, report only
- aquasecurity/chain-bench: reject: CIS supply-chain checks, includes branch settings; report only; not checked further
- kusari-oss/darnit: reject: compliance audit framework with attestations and remediation; not read; unclear license; state not in project Git (not stated)
- sethvargo/ratchet, Pinning GitHub Actions guide: reject: pin Actions by SHA; relevant to F-0001#8 only by analogy
- in-toto, SLSA, sigstore (Cosign/Fulcio/Rekor), Witness, slsa-github-generator: reject: artifact provenance; no repo-settings proof; heavy dependencies
- chainloop-dev/chainloop: reject: attestation control plane with its own server (state not in project Git); Go, Apache-2.0, pushed 2026-10-10
- Conftest: reject: Rego tests on config files; an external tool could check the setup file but the brief has no Rego need; not checked
- "On Omitting Commits and Committing Omissions" (USENIX Security 2016): reject: about git metadata tampering by a server; background for the idea that the forge is trusted for the log
- safedep/gryph, audit layer for agent actions: reject: local SQLite log outside project Git
- Reproducible Builds material: reject: byte-for-byte builds; setup determinism is a different goal
- Dependency, SBOM, vulnerability sections (about 500 entries): reject: off problem

### bradAGI/awesome-cli-coding-agents (1eb67af8afc7)

- neul-labs/grite: keep (low): see Candidates (duplicated in the list)
- basilisk-labs/agentplane: keep (low): see Candidates
- github/gh-aw: keep (low): see Candidates
- h5i-dev/h5i: reject: list says "run metadata versioned under refs/h5i/*", but the repository's own description on 2026-10-11 is "agent-native web security workspace"; the list text no longer matches the repository
- Mneme (MnemeHQ): reject: ADR-to-guardrail; Python; enforces for agents "at the earliest boundary"; no forge proof; not opened
- Squelette (JyMinet): reject: records human decisions, pre-commit gate; local hook the agent can bypass; no forge rule
- Hivelore: reject: local hooks and CI guard; agents can edit them
- Data Olympus: reject: MCP knowledge server; state outside Git
- GitClaw: reject: agent framework, rules in files, no protection from the agent
- AgentDiff: reject: provenance in git refs of who wrote a line; not rule protection
- Bernstein, NEEDLE, Dahrk: reject: deterministic agent orchestration; not repository setup
- Orbi: reject: merge gate bound to reviewed commit is a good idea but fair-code license and Docker service
- Crewplane, Ralph Workflow, agx: reject: resume from files on disk, not from project Git
- Mantyl: reject: handover passport; not a setup
- molt, isitdone, CodeVetter: reject: completion checks local to the agent
- Beads: reject: Dolt store outside project Git
- kasetto: reject: agent environment manager
- numbat, HOL Guard, TraceFold: reject: endpoint hooks for agents; not forge rules
- Harness Starter Kit, AgentLint, schliff: reject: agent instruction files and lint; no setup from pinned baseline
- Remaining ~120 agents, orchestrators, session managers: reject: off problem

## The search file

`sources/gh-search.tsv`, 150 result rows, read in full. Metadata from the tsv unless noted.

Query "github ruleset as code":
- Vivswan/github-settings-as-code: keep: pattern; fetched, see Candidates
- noshiro-pf/github-settings-as-code: keep (low): backup/apply of settings and rulesets, Apache-2.0, TypeScript, pushed 2026-08-10 (tsv); pattern for dump-then-apply; not opened
- theztefan/rulesets-as-code-action: keep (low): MIT, TypeScript, pushed 2025-03-24 (tsv); org rulesets as code; pattern only
- AlexAsAService/grac: reject: no license, 0 stars, TypeScript, no stated detail
- nega-xavier/github-configuration-as-code, alvie-devsecops-lab/github-governance, pelotech/terraform-github-orgkit, BuzzL/workforce-github, carpet-stain/infra, dryvist/tofu-github: reject: HCL/Terraform, admin token and state file outside Git
- Branch protection queries: m19c/protector: reject: TypeScript, last push 2023-01, classic branch protection
- tomdewildt/github-settings-as-code, patel5d2/github-terraform-infra, cmichels-engineering/repo-governance: reject: Terraform
- rdone44/gate: reject: "quality gate", no detail; JavaScript
- thaw-app/platform: reject: Pulumi, org scope
- Snaatak-Saarthi/Authz-Setup, Abhinav-1901/Authz-Setup: reject: no content
- VeVarunSharma/github-agentic-sdlc-starter-pack: reject: Azure/Terraform template; description says branch protection as code; not opened further

Query "github repository settings sync":
- kbrashears5/github-action-repo-settings-sync: reject: Shell Action, last push 2023-09, classic settings only
- joshjohanning/bulk-github-repo-settings-sync-action, joshjohanning/sync-github-repo-settings, Wuodan/sync-github-repo-settings: reject: bulk repo settings (merge options etc.), not rulesets
- EclipseFdn/eclipsefdn-github-sync: reject: org sync, last push 2021
- EndBug/export-label-config, sadewole/antigravity-sync, YoraiLevi/sync-settings-with-github, cran/settingsSync, nomankhalil/...: reject: off problem

Query "repository compliance audit github rules" / "github audit log branch pushes" / "github repository activity audit":
- beast99-p3/Git-Branch-Protection-Auditor: reject: Python read-only scan of classic branch protection; pattern of reading settings only
- ministryofjustice/repo-audit: reject: Ruby, last push 2019
- turbot/tailpipe-mod-github-audit-log-detections: reject: uses the org audit log (not the repository activity endpoint); HCL; needs Powerpipe/Tailpipe
- VAIXLNS/VAIXLNS, Sfedfcv/redesigned-pancake, nMaas8388/github-ranking-audit, Sushma-1706/RepoRefine, exc-analyzer/exc, Uday1004/..., ism4il-04/..., deepakvenkatram/GitHub-Organization-Audit-Script, samkiva/github-health-cli, prem-the-dev/gh-repo-health, Maheshranjanpatra/Github-Roast-Rescue: reject: off problem (profile score, ML, dump of docs)

Query "policy as code repository governance":
- lucashgrifoni/OSS-Security-Policy-as-Code-Starter-Kit, creativ3395/(fork): reject: Python scan of OSS signals; not setup
- goutham04acharya/policy-as-code, squidstack/platform-policies, Boston-Code/identity-governance-policies, balaramr/spacelift-iac-orchestration, highlands-innovation/cic-public: reject: other domains
- CloudRader/github-governance: reject: Terraform org
- Acme-finance-demo/cursor-governance-demo: reject: vulnerability fleet demo

Query "required status checks github app":
- ynsa/pr-review-gate-scratch: keep (low): "GitHub App check run + required status check pin test", a scratch repo with no license; shows the same app-pinning question; not opened, do not rely
- builtbysharath/gh-pr-check-explainer: keep (low): MIT, JavaScript, pushed 2026-10-11; read-only evidence on why a required check is missing, including App and workflow source; pattern for the "check posted by the wrong app" read-back; not opened
- BrianCArnold/RequireUserApproval: reject: GitHub App posting a status; JavaScript; last push 2024-05
- OP5dev/Require-Team-Approval, LucianBuzzo/automerj, sreedhargs89/pr-validator: reject: approval or merge automation
- mercerheather476/turbo-garbanzo, klonnet23/helloy-word, Sfedfcv/redesigned-pancake, Nate0634034090/...: reject: noise (Android, Metasploit, docs dump)

Query "git as system of record configuration" / "git notes metadata branch":
- all rows (FTC SDK, pentestlab, OpenProject, build logs, bash completion): reject: noise

Query "pinned template version project generator":
- all rows (China propaganda repos, gege-circle, misc): reject: noise

Query "template update existing project" / "template sync existing repository":
- ChimeraTK/project-template: keep (low): LGPL-3.0, CMake, pushed 2026-09-29 (tsv); a template repo whose scripts are updated in existing projects by git merge; pattern for second run by merge; LGPL
- nilorg/ntpl: keep (low): MIT, Go, pushed 2026-04-16 (tsv); CLI syncing Git template repos into existing projects; not a dependency (module); pattern; not opened
- yunier-rojas/sombra-cli: reject: scaffold from existing code, 0 stars
- Django, todo, CRUD repos, belaja-akacija prototype: reject: noise

Query "configuration values without evidence":
- GoldstonianConcordanceBible/402: reject: governance prose repo, no tool, no content for the problem

Query "repository scaffolding template" / "project template reproducible":
- rdmurphy/create-clone: reject: JavaScript, copies a repo; no recorded version
- nrjdalal/gitpick: reject: clone part of repo; no record
- stefanbuck/cookiecutter-template, terra-money/template-scaffolding: reject: template fill without source records
- microsoft/cookie-doh: reject: Python package template, Jinja
- kriasoft/cloudflare-starter-kit, atapas/nextshad, airscripts/base, nwops/retrospec-templates, aspect-build/aspect-workflows-template: reject: stack-specific
- TheLartians/ModernCppStarter, econ-project-templates, Turing Way templates, cookiecutter-reproducible-research, snakemake-econ-r, NBISweden/project_template, others: reject: research or C++ templates; "reproducible" means results, not setup

Query "github rulesets":
- katiem0/gh-migrate-rulesets: keep: see Candidates
- kuhlman-labs/repo-ruleset-bot: keep: see Candidates
- int128/wait-for-workflows-action: reject: waits for workflows for ruleset checks (Action, TypeScript)
- 0xnyn/rulesets, davenicoll/github-goodies, peterdresslar/rulesets, ghsioux/cool-rulesets: reject (low): shared ruleset JSON samples; no code; license mixed; could serve as sample JSON, not opened
- UtrechtUniversity/yoda-ruleset, ensingm2/AI-threat-modeling-rulesets, Robotti-io/copilot-security-instructions: reject: other meaning of "ruleset"

Query "ai agent guardrails ci":
- chainloop-dev/chainloop: reject: see lists; state on its control plane
- sgroy10/speclock, JeongJaeSoon/agent-guard, TheArchitectit/AIGGP-Agentic-Framework, ctxwitch/ctxwitch, agent-axiom/agent-anvil: reject: local guardrails the agent can edit, or evals
- Azure-Samples/foundry-citadel-platform, Tirthanand17/*, vishal2505/bedrock-hr-agent: reject: unrelated

Query "coding agent branch protection":
- blackwell-systems/agentic-project-scaffolding: reject: Claude skill that creates repo and branch protection; the agent itself sets the rules, the opposite of F-0001#3
- mak3r/claude-scaffolding: reject: personas template, "branch protection" not shown as proof
- donggun-jung/multiagent-protocol: reject: self-built branch protection for solo developers, Python; no forge ruleset proof
- aapp-protocol/aapp-kit: reject: worktrees on orphan branches; no forge rule
- LevyBytes/AI-SKILL-github: reject: AGPL reference skill

Query "copier template update":
- actions-ext/copier-update, actions-ext/copier: keep (low): see Candidates
- liblaf/copier-shared, piotrgredowski/..., pyfar/template-pyfar-package, kasi-x/python-copier-template-example, EmanuelVogt/platform-template: reject: Copier templates for other stacks; EmanuelVogt notes "one copier copy generates the product; a semver tag updates" (a pinned-version pattern, same as copier itself)
- vladmesh/service-template: reject: Copier framework for AI-ready services; stack-specific
- thesmart/repo-copier, RaynOlivia/copier-template-mcp: reject: wrap copier for agents; the agent runs the update, not code with a record

## What no candidate covers

- A program that writes configuration values each with a source (answer, catalog entry, fact) and leaves a value without a source as an open gap. Copier records answers but not the source of each value and has no open-gap state. No list entry or search row does this.
- A setup that stops at each human decision and resumes from its own record in Git. Resume exists in agent runners, but their records are on disk or in a server, not on a branch only one identity can update.
- Free text written by an agent and each part checked by code before it is accepted. gh-aw "safe-outputs" is the nearest idea; I did not open it.
- An evidence branch that only the controller identity may update, and its ruleset. No public tool pairs these. Custom git refs (grite) are outside branch rulesets, by my reading.
- The probe push by the controller's own identity that the forge must refuse for a rule violation, with "any other failure is not a pass". No source found describes this. I also did not find an official page for the refusal text (GH013 query returned only low-quality pages). The pass test must therefore be built on the push result and the read-back together, not on text matching a page I could not confirm.
- A read-back that compares the effective rules, including `integration_id` pinning, with the file the program wrote. The API gives the data (rules/branches); no tool in the sources compares it.
- "Repository passes its own gates when the controller is removed" as a check. Nothing found.
- A second run on a repository set up by an earlier version. Copier and ChimeraTK show merge-based updates; neither knows about evidence, rulesets, or gaps.
- Tools usable as a dependency: none. The controller is Go with the standard library only; every working tool found is a Go module, TypeScript, Python, Rust or Terraform. Only API usage, patterns and formats transfer.

## Limits

- I did not open READMEs or code of any kept candidate except through `api.github.com` metadata and the list or tsv text. Statements about what a tool does come from its description or from the list text, and are marked "not checked" where they go further.
- The official page for "rules" did not state the permission for create ruleset. The Administration (write) claim comes from a search summary of GHES docs. It should be confirmed against the github.com page before use.
- The activity endpoint page did not state permissions, retention, or how pushes by a GitHub App appear in `actor`. I did not check whether an App push shows the App or a bot user.
- I did not confirm the exact text or exit form of a push refused by a ruleset (the "GH013" query found no official page).
- I did not check whether a ruleset covers refs that are not branches or tags (custom refs); the grite remark is my reading.
- I did not run a search for the F-0001#4 or F-0001#1 themes myself, nor for Go-only tools beyond the lists. Search engine coverage was one web search per theme, standard mode.
- License fields: NOASSERTION for Vivswan means the license is not recognised, not that none exists. Entries marked (tsv) have unverified metadata as of 2026-10-11.
- I ran no code and installed nothing. I read no file outside this directory.
- The h5i entry shows list text can drift from its repository; other list descriptions I did not re-check against their repositories may drift too.
~~~~
