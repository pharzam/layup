# T-msrd — the search of M2d, searcher A, word for word

Searcher: GPT-6 Sol (`gpt-6-sol`, a model of the reasoning tier) on Devin 3000.11.3 (`devin -p --permission-mode dangerous --respect-workspace-trust false --prompt-file brief.md`), a fresh session with a `HOME` of its own that held only copies of Devin's credential and configuration files (removed after the run), 2026-10-11T14:40:37Z to 2026-10-11T14:43:57Z (3 min 20 s), exit 0, on the first try. Devin reports no usage. Devin's paid models answered again on 2026-10-11 at 14:40Z (a one-word probe of `gpt-6-astra-low` and `gpt-6-sol`), so the free model of O-179 was not needed. Brief: [`brief.md`](brief.md), with the five lists of [`sources.md`](sources.md) and [`gh-search.tsv`](gh-search.tsv); nothing of the LAYUP architecture or the specification. The output follows word for word inside a fenced block; its paths are relative to the searcher's directory.

~~~~text
# M2d public-solution search — 2026-10-11

Scope: existing public artifacts, not an implementation proposal. “Keep” means useful as a precedent or documented interface, **not** approval to add a controller dependency. The controller is Go standard library plus the `git` executable; its App lacks Administration and Workflows permissions. Dates below are repository `pushed_at` dates from the linked GitHub repository API unless stated otherwise. No candidate was executed.

## Queries

| Query (run before reading supplied lists) | Words it comes from | What it found |
|---|---|---|
| `GitHub rulesets required status checks integration_id bypass actors evaluate ruleset API` | “required checks pinned to the app”; “no bypass” | [Repository rules REST API](https://docs.github.com/en/rest/repos/rules): `integration_id`, bypass actors, active versus evaluate rules. |
| `GitHub ruleset push probe rejected rule violation audit repository activity API` | “probe push … refuse for a rule violation”; “activity log” | [Rule suites](https://docs.github.com/en/rest/repos/rule-suites) list failures by actor/ref, with per-rule evaluations; not a substitute for the actual rejected push. |
| `reproducible repository template initialization questionnaire provenance evidence git human decisions resumable` | “repeatable, evidence-based setup”; “human decision” | [AgenticProjectInitializer](https://github.com/pbeens/AgenticProjectInitializer) prompts, [reproducible-project-template](https://github.com/alan-turing-institute/reproducible-project-template); neither search excerpt shows per-value cited evidence. |
| `Go standard library git executable project template setup resumable human approval provenance` | “Go, with the standard library only, and Git called as the git program” | [git-scaffold](https://github.com/stephenc/git-scaffold), an unusually close pinned-scaffold/update precedent, but not an importable dependency. |
| `GitHub App activity log repository ruleset branch push audit log API organization` | “read … activity log”; “controller’s own identity” | [Organization audit log docs](https://docs.github.com/en/organizations/keeping-your-organization-secure/managing-security-settings-for-your-organization/audit-log-events-for-your-organization): `git.push` has special access/retention, unlike repository activity. |
| `GitHub REST repository activity endpoint branch push activity actor activity type fine grained permissions` | “audit who updated each branch” | [Repository activity REST API](https://docs.github.com/en/rest/repos/repos#list-repository-activities): ref, actor, push/force-push/creation/deletion/merge. |
| `GitHub rulesets restrict updates actor integration bypass no bypass app repository branch` | “only the controller’s identity may update”; “no bypass for the controller” | [GitHub rules](https://docs.github.com/en/rest/repos/rules) distinguish restriction-rule bypass from bypass of other rules in a ruleset; a tension requiring distinct scopes/rulesets. |
| `Copier documentation answers file _commit update version pin skip questions` | “pinned, recorded version”; “second run” | [Copier answers](https://copier.readthedocs.io/en/stable/creating/) record `_commit` and `_src_path`; [updates](https://copier.readthedocs.io/en/stable/updating/) reuse answers, with conflicts possible. |
| `git scaffold stephenc git-scaffold README license lock --existing` | “pinned baseline commit”; “earlier version set up” | [git-scaffold specification](https://github.com/stephenc/git-scaffold/blob/main/DESIGN.md): lock SHA, explicit patches, `init --existing`, deterministic managed files. |
| `GitHub rulesets required pull request branch creation restrict updates probe branch ruleset creation direct push` | “PR required”; “push to layup-probe is refused” | [Available rules](https://docs.github.com/en/repositories/configuring-branches-and-merges-in-your-repository/managing-rulesets/available-rules-for-rulesets): creation and update are distinct; a PR rule does not by itself restrict creation of a probe branch. |

The supplied `sources/gh-search.tsv` contains **additional pre-run queries**, not my own. Its query-to-fact mappings are preserved in that file. The judgments below cover touching entries from those results and the five complete pinned lists in `sources/lists.tsv`.

## Candidates

| Kept candidate and own documentation | License | Language | Last activity | Part covered; state in project's Git; deterministic? | Keep as / why |
|---|---|---|---|---|---|
| [git-scaffold](https://github.com/stephenc/git-scaffold) ([spec](https://github.com/stephenc/git-scaffold/blob/main/DESIGN.md), [metadata](https://api.github.com/repos/stephenc/git-scaffold)) | Apache-2.0 | Go | 2026-08-28 | Managed baseline files, second run and adoption of an existing repo. `.git-scaffold/config.toml`, `lock`, patches **should** be committed in project Git. Spec explicitly says managed output deterministic given lock, arguments and patches. | **Pattern/format**, not dependency: literal substitution, full commit lock, drift check, explicit divergence. Does not attach citations to every argument or gate GitHub. |
| [Copier](https://github.com/copier-org/copier) ([answers](https://copier.readthedocs.io/en/stable/creating/), [update](https://copier.readthedocs.io/en/stable/updating/), [metadata](https://api.github.com/repos/copier-org/copier)) | MIT | Python | 2026-10-08 | Questionnaire plus template update. `.copier-answers.yml` in generated project, intended to be versioned in Git, stores source and `_commit`. **Not guaranteed deterministic**: Jinja/tasks, moving default tag and manual merge conflicts; can pin a Git ref. | **Pattern/format** for answers and versioned reruns, not dependency; no per-answer evidence or human stop/resume authority. |
| [GitHub repository rules REST API](https://docs.github.com/en/rest/repos/rules) ([available rules](https://docs.github.com/en/repositories/configuring-branches-and-merges-in-your-repository/managing-rulesets/available-rules-for-rulesets)) | not stated (hosted API documentation) | HTTP/JSON | not stated | Effective rules on a branch, repository/organization rulesets, `integration_id`, `pull_request`, `non_fast_forward`, `creation`, `update`, bypass. **Forge state**, not automatically project Git. API readback deterministic for a fixed current forge state, not immutable over time. | **API usage/known limit**: `GET /repos/{owner}/{repo}/rules/branches/{branch}` returns active applicable rules, even for a branch that does not yet exist; inspect full rulesets and bypass lists too. Writes require authority the App lacks. |
| [GitHub rule suites REST API](https://docs.github.com/en/rest/repos/rule-suites) | not stated (hosted API documentation) | HTTP/JSON | not stated | Rule evaluations by ref/actor and `pass`, `fail`, `bypass`, with individual rule details. Forge-held activity, not project Git unless saved separately. Deterministic report of recorded evaluations; **not stated** that every failed transport/auth attempt generates a suite. | **API usage** corroborating a rejection: filter `ref`, `actor_name`, `rule_suite_result=fail`, `evaluate_status=active`; do not equate a failed Git command with rule enforcement. |
| [GitHub repository activity REST API](https://docs.github.com/en/rest/repos/repos#list-repository-activities) ([audit-log limitation](https://docs.github.com/en/organizations/keeping-your-organization-secure/managing-security-settings-for-your-organization/audit-log-events-for-your-organization)) | not stated (hosted API documentation) | HTTP/JSON | not stated | `/repos/{owner}/{repo}/activity`: actor, ref, SHA pair, timestamp, type including push and force push. Forge state, not project Git until persisted. Deterministic read of a given returned activity page; retention/completeness across time **not stated** there. | **API usage/known limit**: audit accepted branch changes and paginate; not a log of rejected pushes. Organization `git.push` audit events require distinct access and have special retention. |
| [gh-migrate-rulesets](https://github.com/katiem0/gh-migrate-rulesets) ([README](https://github.com/katiem0/gh-migrate-rulesets/blob/main/README.md), [metadata](https://api.github.com/repos/katiem0/gh-migrate-rulesets)) | MIT | Go | 2026-08-04 | CSV export/import of rulesets; fields include bypass actors, `required_status_checks` context/integration and non-fast-forward. Report is a local CSV, **not necessarily** project Git. Determinism of full migration **not stated**. | **Format/tool** for a human-admin ruleset file and cross-check, not a dependency: requires org-owner/PAT admin scope for fullest use, which controller cannot have. README labels for update/deletion fields are inconsistent; API is authority. |
| [rulesets-as-code-action](https://github.com/theztefan/rulesets-as-code-action) ([README](https://github.com/theztefan/rulesets-as-code-action/blob/main/README.md), [metadata](https://api.github.com/repos/theztefan/rulesets-as-code-action)) | MIT | TypeScript | 2025-03-24 | JSON ruleset stored in governing Git repo, CODEOWNERS-approved then applied to **organization** via REST. State in a separate control repo, not necessarily project Git; deterministic application **not stated**. | **Pattern** for operator-controlled JSON + approval, not runnable controller dependency; privileged workflow/App deployment conflicts with controller permission boundary. |
| [repo-ruleset-bot](https://github.com/kuhlman-labs/repo-ruleset-bot) ([README](https://github.com/kuhlman-labs/repo-ruleset-bot/blob/main/README.md), [metadata](https://api.github.com/repos/kuhlman-labs/repo-ruleset-bot)) | MIT | Go | 2025-04-21 | Org ruleset JSON in bot repo/release; webhook detects and reverts unauthorized changes. **Separate** Git repo, not project state; deterministic reconciliation **not stated**. | **Known limit/pattern**: requires org Administration read/write, releases and running service; cannot be this App. Reversion after a change is not proof an unauthorized push was prevented. |
| [Allstar](https://github.com/ossf/allstar) ([README](https://github.com/ossf/allstar/blob/main/README.md), [metadata](https://api.github.com/repos/ossf/allstar)) | Apache-2.0 | Go | 2026-10-11 | Continuously checks repository policies; can issue alerts or repair some settings. `.allstar` control repository and optionally per-project `.allstar/`; not wholly in project's Git. Deterministic enforcement **not stated**. | **Tool/known limit**: hosted OpenSSF Allstar App **retired**, must self-host; issue/repair is a monitoring precedent, not proof of hard, pre-push gate integrity. |
| [Agentplane](https://github.com/basilisk-labs/agentplane) ([README](https://github.com/basilisk-labs/agentplane/blob/main/README.md), [metadata](https://api.github.com/repos/basilisk-labs/agentplane)) | MIT | TypeScript | 2026-10-11 | Human approvals, bounded agent episodes, supervisor-observed checks, typed records in `.agentplane/tasks/` beside project code; Git is durable review surface. CLI transitions advertised deterministic; model reasoning explicitly not. | **Pattern** for stop/resume and checked agent output; not Go stdlib dependency, and not a pinned discipline-baseline or independent forge guard. |
| [Squelette](https://github.com/JyMinet/squelette) ([README](https://github.com/JyMinet/squelette/blob/main/README.md), [metadata](https://api.github.com/repos/JyMinet/squelette)) | MIT | Python standard library | 2026-09-29 | Records human decisions, authorized paths, preflight and checked evidence in project Git; `status` reconstructs state from records. Deterministic controller refusal claimed; agent content **not stated** deterministic. | **Method** for explicit decisions and proof before DONE; local commit hook can be modified or bypassed by a writer, unlike GitHub rulesets. No documented GitHub rejection probe. |
| [Grite](https://github.com/neul-labs/grite) ([README](https://github.com/neul-labs/grite/blob/main/README.md), [metadata](https://api.github.com/repos/neul-labs/grite)) | MIT | Rust | 2026-07-02 | Agent tasks/decisions as append-only event log in project `refs/grite/wal`; optional local sled materialized cache. Docs claim deterministic CRDT resolution; not deterministic agent answers. | **Format/pattern** for Git-native resumption, not a Go dependency or protected setup evidence branch; its ordinary agent-writable refs are not authorization. |
| [agentdiff](https://github.com/codeprakhar25/agentdiff) ([README](https://github.com/codeprakhar25/agentdiff/blob/main/README.md), [metadata](https://api.github.com/repos/codeprakhar25/agentdiff)) | MIT OR Apache-2.0 (README; API labels Apache-2.0) | Rust | 2026-06-08 | Signed line-level *claimed agent authorship* via Git refs `refs/agentdiff/*`, optional CI; global hooks and local `.git/agentdiff/` capture also used. Attribution deterministic **not stated**. | **Pattern/known limit**: Git-native provenance is not forge-authenticated identity for each accepted branch update; hooks cannot prove no untracked writes. |
| [Conform](https://github.com/siderolabs/conform) ([README](https://github.com/siderolabs/conform/blob/main/README.md), [metadata](https://api.github.com/repos/siderolabs/conform)) | MPL-2.0 | Go | 2026-10-09 | Declarative `.conform.yaml` tests commit metadata and license headers; policy can live in project Git. Fixed checks deterministic in ordinary use, formal guarantee **not stated**. | **Tool/pattern** for independent CI gate execution; does not check all Armature discipline values, provenance, or forge rule configuration. List's old `talos-systems` link redirects to `siderolabs`. |
| [ChimeraTK/project-template](https://github.com/ChimeraTK/project-template) ([HowTo](https://github.com/ChimeraTK/project-template/blob/master/HowTo_project-template.md), [metadata](https://api.github.com/repos/ChimeraTK/project-template)) | LGPL-3.0 | CMake | 2026-09-29 | CMake baseline can be merged into an existing project; files in target Git. Deterministic second merge **not stated**; conflicts are possible. | **Method** illustrating second-run template merge; no README by design to avoid conflicts; no cited configuration evidence. |

### Judgments on the interface boundaries

- The rules API documents `required_status_checks.parameters.required_status_checks[].integration_id`: a check *name alone* does not pin its poster. A written JSON/CSV recipe is only a request to an admin, not proof the rule was applied.
- The branch-rules endpoint includes active rules inherited from higher levels and works for a nonexistent branch; the ruleset endpoint exposes `enforcement`, conditions and bypass actors. Both views matter. A disabled/evaluate ruleset is not an active protection.
- The same GitHub rules documentation says `update` allows **only bypass actors** to update the matching ref. Therefore “only controller updates the setup evidence branch” needs an explicit controller exception to *that restriction*, whereas “no controller bypass” of the **gate** rules is a different claim. A single ruleset whose bypass exempts the controller from its checks cannot satisfy both; none of the candidates supplies a tested configuration resolving the scopes.
- Repository rule suites include a `fail` result and per-rule details; they are stronger corroboration than an exit code alone. They do not replace matching the actual `git push` output and identity/ref/SHA to a **ruleset violation** rather than an authentication, network, missing-ref or permissions failure.
- A GitHub repository activity entry reports successful changes, not a full continuous audit archive and not the identities behind rejected attempts. The organization's audit-log `git.push` event has its own access and retention rules, documented separately. Git commit author is not reliable push actor identity.
- Requiring a PR governs updates to a targeted branch; a probe of a new `layup-probe` branch may be a branch *creation*, subject to different rules. No searched artifact proves the intended first-demo rejection without a live probe against the specific active ruleset.
- GitHub rulesets live on the forge. A file in project Git and the effective forge rule are two distinct states; readback and behavior must be checked independently. A separate policy check in the project cannot stop agents from changing its own source unless the forge controls that path/ref.

### Source-specific coverage notes

The following are narrower claims actually supported by the named candidate's own source; “not stated” is deliberate rather than an assumption of failure.

**Pinned setup and replay**

- [git-scaffold's specification](https://github.com/stephenc/git-scaffold/blob/main/DESIGN.md) describes its lock as a full Git commit SHA with a trailing newline.
- Its target configuration separately names the source Git URL/ref and literal argument values.
- Its source descriptor enumerates managed paths and literal substitution tokens; it does not render arbitrary template expressions.
- Its specification rejects missing required arguments rather than manufacturing values for them.
- It also treats an absent target value with a source default as acceptable; that default is not independently cited to a project fact or catalog entry.
- Its claim of determinism is about **managed file contents** for fixed inputs, not about an LLM deciding which inputs are right.
- Its documentation says `init --existing` records local differences as patches and makes the initial scaffold check clean.
- A clean scaffold check means files match their declared lock/patches; it does **not** mean Armature discipline tests or forge rules pass.
- Its explicit update action is relevant to the second run, but it does not say it understands an earlier unrelated setup version's record format.
- [Copier's updating guide](https://copier.readthedocs.io/en/stable/updating/) requires a valid prior answers file and versioned source and destination Git for the normal smart-update route.
- Its guide says an update can produce conflict markers or `.rej` files and asks for human review; “updatable” is not “automatic second run succeeds.”
- The guide explicitly says not to change the answers file by hand, since the update diff may then be computed against false prior state.
- [Copier's template docs](https://copier.readthedocs.io/en/stable/creating/) identify `_src_path` and `_commit` in `_copier_answers`, but contain no field meaning “evidence supporting this answer.”

**Human authority and agent-written output**

- [Agentplane's README](https://github.com/basilisk-labs/agentplane/blob/main/README.md) gives the human material-risk approval, the agent semantic judgment, and the CLI formal transitions.
- It says a bounded episode returns an objective, writable scope, result schema and exact resume command.
- It records missing provider telemetry as unavailable or partial, not a fabricated zero: a useful explicit-unknown precedent.
- Its trust boundary says the supervisor observes repository changes and executed checks but does not judge semantic quality.
- Its records concern coding tasks, not provenance citations for every new repository configuration value.
- [Squelette's README](https://github.com/JyMinet/squelette/blob/main/README.md) says `status` reconstructs state from project records after a new session.
- It describes explicit human decisions, scoped work items, checks before writes and evidence before DONE.
- It describes a commit hook installed outside the worktree, but does not document a forge-side rule denying an agent that can run its own Git client.
- [Grite's README](https://github.com/neul-labs/grite/blob/main/README.md) stores task events in `refs/grite/wal` and can rebuild a local materialized view.
- It claims deterministic conflict resolution when actors concurrently edit issues, not one-author-only Git ref updates.
- Its optional event signatures attest to event origin within its scheme; they do not prove the forge accepted a push from only the controller App.
- [agentdiff's README](https://github.com/codeprakhar25/agentdiff/blob/main/README.md) explicitly requires both global capture-hook configuration and per-repository initialization.
- Its Git refs preserve attribution records; GitHub's repository activity, by contrast, attributes the network operation updating a branch.
- The two kinds of provenance answer different questions, so neither can be substituted for the other.

**Forge settings and checks**

- [gh-migrate-rulesets' README](https://github.com/katiem0/gh-migrate-rulesets/blob/main/README.md) says its CSV includes enforcement status, bypass actors, required-check context/integration and non-fast-forward.
- Its README asks for organization-owner authorization and an admin-scoped token for full usage: evidence against treating it as controller-integrated code.
- Its CSV is a snapshot/report or import input; it does not claim to prove GitHub actually rejected a specific push.
- [rulesets-as-code-action's README](https://github.com/theztefan/rulesets-as-code-action/blob/main/README.md) requires a GitHub token to update org rulesets from a JSON file after CODEOWNERS review.
- Its example runs from a central rules-governing repo; it is not a project-local evidence branch created with limited App privileges.
- [repo-ruleset-bot's README](https://github.com/kuhlman-labs/repo-ruleset-bot/blob/main/README.md) explicitly requires Organization Administration read/write.
- It distributes JSON from release assets, then repairs modifications or deletions; repair is necessarily distinct from preventive server-side denial.
- [Allstar's README](https://github.com/ossf/allstar/blob/main/README.md) says its formerly hosted GitHub App is retired and that users must now run it as an Action or daemon.
- Allstar's configuration can come from a separate organization `.allstar` control repository or project-local `.allstar/`, depending on configured override policy.
- A project-local override allowed to maintainers is not, by itself, protection of policy from governed agents.
- [Conform's README](https://github.com/siderolabs/conform/blob/main/README.md) shows local `.conform.yaml`, an `enforce` invocation and optional commit-message hook.
- Its examples check commits and license headers, not the provenance of discipline parameters or who could alter its own YAML.
- GitHub's [rule-suite response](https://docs.github.com/en/rest/repos/rule-suites) includes `before_sha`, `after_sha`, `actor_name`, `ref` and per-rule evaluations.
- GitHub's [repository activity response](https://docs.github.com/en/rest/repos/repos#list-repository-activities) associates accepted operations with an actor and ref; this is a different source of evidence from an agent's signed local attribution.
- [GitHub's ruleset description](https://docs.github.com/en/rest/repos/rules) gives bypass modes `always`, `pull_request` and `exempt`; the last mode explicitly does not create a bypass audit entry.
- The same description says a ruleset in `evaluate` mode tests rules rather than enforcing them; a candidate that only reports evaluation cannot establish the required rejected push.

## The lists

Entries below are the entries touching setup, Git-state continuity, agent restrictions, policy checks or forge governance; unrelated package-vulnerability, Kubernetes-only, build-artifact signing and generic coding-agent entries do not need lines. List commit IDs: policy `af4e6354067b`; supply-chain `31252e6810c4`; CLI agents `1eb67af8afc7`; GitOps `82af5b15f445`; Actions `0edccf33480b` (see `sources/lists.tsv`). Duplicated items are marked duplicate rather than silently omitted.

### hysnsec/awesome-policy-as-code

- OPA: reject: general policy engine, no scoped GitHub ruleset/application proof; new runtime.
- Styra DAS: reject: commercial hosted policy state, not project Git.
- OPAL: reject: policy distribution rather than setup evidence or forge rules.
- OPCR: reject: policy-container supply chain, not repository gates.
- Topaz: reject: authorization data plane, not forge-enforced Git restrictions.
- HashiCorp Sentinel: reject: embedded policy language, extra runtime and wrong enforcement point.
- Regula: reject: cloud IaC validation, not setup provenance.
- Intercept: reject: static policy auditing, no GitHub push identity controls documented by list.
- Checkov: reject: IaC scanning, not repository setup proof.
- Terrascan: reject: cloud IaC compliance, not forge rules.
- kics: reject: IaC misconfiguration checker, not controller-only Git branch.
- Gatekeeper: reject: Kubernetes admission, not GitHub branch governance.
- Gatekeeper Policy Manager: reject: Kubernetes policy UI, not GitHub rulesets.
- Konstraint: reject: Kubernetes Gatekeeper policy management, wrong target.
- kube-mgmt: reject: Kubernetes OPA sidecar, wrong target.
- Kyverno: reject: Kubernetes admission, not GitHub branch governance.
- MagTape: reject: Kubernetes policy engine, wrong target.
- Cloud Custodian: reject: cloud resource remediation, not required-check App pinning.

### bureado/awesome-software-supply-chain-security

- GitBOM: reject: provenance of build inputs, not every setup configuration value.
- bomsh: reject: tooling for build-input provenance, not setup citations.
- GitXray: reject: GitHub repository security triage, no stated rule-violation push probe.
- Witness: reject: build-artifact attestations, not answer-by-answer setup citations.
- in-toto: reject: supply-chain step attestations, not human decisions during repo setup.
- SLSA: reject: artifact provenance requirements, not branch protection readback.
- Open Source Project Security Baseline: reject: useful project checklist, not a reproducible Armature setup runner.
- Conftest: reject: policy assertions over structured files, not trusted forge enforcement.
- Open Policy Agent: reject: generic policy evaluation, unnecessary controller dependency.
- Conforma: reject: verification of built artifacts, not branch-actor audit.
- kusari-oss/darnit: reject: engineering compliance attestations, not cited baseline inputs.
- scribe-public/gitgat: reject: source-control security posture, not documented refusal probe.
- Legit-Labs/legitify: reject: misconfiguration scan/remediation, not per-value setup evidence.
- crashappsec/github-analyzer: reject: organization security settings scan only.
- reposaur/reposaur: reject: development-platform compliance, not stated Git-native setup state.
- ossf/allstar: keep: GitHub policy monitoring/repair precedent, with hosted-App retirement caveat.
- kusaridev/skootrs: reject: secure repo creation, but no documented source-cited values or operator pause.
- chainloop-dev/chainloop: reject: signed build/agent attestations in external control plane, not project Git setup record.
- aquasecurity/chain-bench: reject: supply-chain compliance audit, not Armature setup tests.
- sethvargo/ratchet: reject: pins workflow dependencies, not discipline baseline/forge rules.
- apiiro/PRevent: reject: malicious-code PR scanner, not protected rule paths.
- SpecterOps/GitHound: reject: permission attack-path mapping, not setup proof.

### bradAGI/awesome-cli-coding-agents

- GitClaw: reject: Git-versioned agent rules and memory but agents can edit them.
- Agent Orchestrator (AO): reject: task worktrees and PR routing, not pinned project setup.
- h5i: reject: run metadata in Git refs and verifier, not per-value setup provenance.
- MartinLoop: reject: signed run receipts, not ruleset readback or evidence-only updater.
- Orbi: reject: reviewed-commit-bound merge gate, no setup baseline or forge rules.
- NEEDLE: reject: deterministic task state machine, but not protected project Git evidence.
- Crewplane: reject: restartable workflow, but local disk logs and no Git branch authority.
- Ralph Workflow: reject: restartable spec loop, not setup evidence/forge rule test.
- Beads: reject: task graph via Dolt sync, not specified project Git-only branch.
- GitHub Agentic Workflows: reject: safe outputs and lock files, but requires workflow authority and is not this Go controller.
- AgentPlane (listed as AgentPlane): keep: supervisor-owned human gates and Git task/evidence record pattern.
- AgentDiff: keep: Git-ref provenance format, with hook/identity limits.
- Grite (appears twice): keep: deterministic event-log Git refs for persistent tasks; one candidate.
- Mneme: reject: governance decisions, but repository Git storage and forge restriction not established by list.
- Tree Ring Memory: reject: evidence-backed SQLite memory, not project Git system of record.
- Okto Nexus: reject: SQLite-backed coordination rather than project Git.
- Project Tiny Context Harness: reject: repo-local context checks, not source citation for all values.
- Hivelore: reject: turns lessons into checks, but agent-editable hooks/CI are insufficient alone.
- Squelette: keep: recorded human decisions and local proof/check gates, limited by local enforcement.
- SpecWave: reject: approval gate and checks, but multi-host Node solution rather than pinned discipline setup.
- isitdone: reject: checks exact working tree, not GitHub rules or baseline initialization.
- Mantyl: reject: handover verification, not initial cited setup.
- Data Olympus (appears twice): reject: context governance, not guaranteed project Git record.
- AgentLint: reject: AI-repo hygiene score, not GitHub-ref immutability.
- Agentic Engineering Framework: reject: agent workflow guidance not external proof.
- schliff: reject: deterministic lint for instruction files, not all setup inputs.
- agent-runbook: reject: checkpointed YAML skill compiler, not forge identity controls.
- BlaBla: reject: executable Python contracts, not Go setup and rule readback.
- grite (duplicate): keep: same candidate, not a second solution.

### weaveworks/awesome-gitops

- GitOps working-group definition: reject: Git desired-state method, not a concrete GitHub ruleset check.
- Atlantis: reject: Terraform PR automation with additional service, not Go controller authority.
- Autoapply: reject: applies Git to infrastructure; GitHub branch rules not addressed.
- Flux: reject: Git-to-Kubernetes reconciliation, wrong target.
- ArgoCD: reject: Kubernetes deployment reconciliation, wrong target.
- Weave GitOps OSS: reject: Kubernetes operations, not repo setup evidence.
- git-secret: reject: encrypted secret storage, not needed for citation records; tokens must not enter Git.
- SOPS: reject: secret encryption, not ruleset verification.

### sdras/awesome-actions

- actions/starter-workflows: reject: general workflow examples, not independent Armature test suite.
- actions/checkout: reject: checkout helper, no baseline/evidence logic.
- actions/github-script: reject: JS API wrapper, not Go stdlib dependency.
- talos-systems/conform: keep: commit/license policy gate; current repository redirects to siderolabs.
- jpmcb/prow-github-actions: reject: PR/chatops enforcement, not evidence-backed setup.
- mszostok/codeowners-validator: reject: validates CODEOWNERS syntax, not forge effective policy.
- benjefferies/branch-protection-bot: reject: temporarily disables protection, directly opposed to no bypass.
- WyriHaximus/github-action-wait-for-status: reject: waits for checks; cannot pin check poster.
- deepakputhraya/action-branch-name: reject: branch-name lint, not server-side identity rule.
- boasiHQ/interactive-inputs: reject: runtime workflow inputs, not Git-persisted human decisions.
- dorny/paths-filter: reject: chooses CI steps based on diff; not a write restriction.
- peter-evans/create-pull-request: reject: opens PR, does not enforce PR-only protected branches.
- cirrus-actions/branch-guard: reject: blocks PR merge when branch checks fail, not server-side ref-actor rule.
- asgharlabs/github-terraform: reject: Terraform GitHub admin application, violates controller dependency and privilege bounds.

## The search file

Relevant hits below retain the repository names from `sources/gh-search.tsv`; a repeated hit is still judged. A repository search description alone does not establish the feature worked. Unrelated false hits (mirrored GitHub docs dumps, politics, IDE settings sync, Django CRUD, profile ranking, Kubernetes/HR apps) need no line.

- Vivswan/github-settings-as-code (rows 2, 13): reject: stateless Action applies privileged settings, not human-admin file handoff/readback.
- noshiro-pf/github-settings-as-code: reject: settings backup/apply, not externally proven push denial.
- theztefan/rulesets-as-code-action: keep: reviewed JSON ruleset in Git, but applies with privileged App.
- AlexAsAService/grac: reject: rulesets as code, no described operator handoff or probe.
- nega-xavier/github-configuration-as-code: reject: Terraform org administration, wrong runtime/authority.
- alvie-devsecops-lab/github-governance: reject: org governance lab, not cited per-value setup.
- pelotech/terraform-github-orgkit: reject: Terraform module, wrong controller dependency/authority.
- BuzzL/workforce-github: reject: fleet Terraform rulesets, not per-project Git evidence.
- carpet-stain/infra: reject: OpenTofu account governance, not setup proof.
- dryvist/tofu-github: reject: OpenTofu organization rules, not app-identity probe.
- m19c/protector: reject: branch protection as code, not operator-applied ruleset readback.
- tomdewildt/github-settings-as-code: reject: Terraform account settings, not permitted runtime.
- rdone44/gate: reject: CI branch protection action, can edit own gate.
- VeVarunSharma/github-agentic-sdlc-starter-pack (rows 16,129): reject: broad starter template, not source-cited configuration or verified rule probe.
- patel5d2/github-terraform-infra: reject: Terraform administration, not Go stdlib.
- cmichels-engineering/repo-governance (rows 18,51): reject: Terraform governance tiers, no per-project evidence/resume.
- thaw-app/platform: reject: Pulumi/Bun org settings, privileged external state.
- kbrashears5/github-action-repo-settings-sync: reject: automated settings sync, not operator handoff.
- joshjohanning/bulk-github-repo-settings-sync-action: reject: bulk setting changes, not effective gate proof.
- EndBug/export-label-config: reject: labels are not branch rules.
- EclipseFdn/eclipsefdn-github-sync: reject: org-wide settings synchronization, not project setup record.
- joshjohanning/sync-github-repo-settings: reject: central config sync, not evidence-cited setup.
- Wuodan/sync-github-repo-settings: reject: config sync, not rule-enforcement proof.
- ministryofjustice/repo-audit: reject: compliance reporting only, no negative push test.
- beast99-p3/Git-Branch-Protection-Auditor: reject: scans configuration, not effective rejection.
- turbot/tailpipe-mod-github-audit-log-detections: reject: org audit analytics, not controller-accessible repository activity proof.
- Uday1004/Activity-Audit-by-Github-Actions: reject: action history, not who pushed each ref.
- deepakvenkatram/GitHub-Organization-Audit-Script: reject: org setting/collaborator audit, not verified repo probe.
- lucashgrifoni/OSS-Security-Policy-as-Code-Starter-Kit: reject: general OSS hygiene, not pinned setup values.
- CloudRader/github-governance: reject: Terraform org IaC, not operator-applied file.
- creativ3395/OSS-Security-Policy-as-Code-Starter-Kit: reject: local governance audit only.
- BrianCArnold/RequireUserApproval: reject: App-posted status shows a gate producer, not setup/probe.
- OP5dev/Require-Team-Approval: reject: PR review status, not check-poster pin/config proof.
- ynsa/pr-review-gate-scratch: reject: scratch pin test description only; no documented reusable result.
- builtbysharath/gh-pr-check-explainer: reject: diagnostics for missing checks, not enforcement itself.
- LucianBuzzo/automerj: reject: automerge consumer of checks, not gate author or ruleset proof.
- sreedhargs89/pr-validator: reject: external Jira/CircleCI PR checks, not Armature setup.
- ChimeraTK/project-template: keep: template merging into an existing CMake project.
- yunier-rojas/sombra-cli: reject: Go scaffold-from-code, no recorded pinned input/checked evidence in search description.
- GoldstonianConcordanceBible/402: reject: provenance/human authorization described, but no concrete project setup or Git branch control established.
- kriasoft/cloudflare-starter-kit: reject: Cloudflare technology-specific boilerplate.
- stefanbuck/cookiecutter-template: reject: GitHub Action and template, not per-value evidence.
- nrjdalal/gitpick: reject: selective clone, no retained baseline lock or second-run behavior.
- terra-money/template-scaffolding: reject: clone/replace library, wrong language and no cited inputs.
- aspect-build/aspect-workflows-template: reject: Bazel-specific starter.
- microsoft/cookie-doh: reject: Python-specific project skeleton.
- rdmurphy/create-clone: reject: repository copying, no verified setup state.
- OpenSourceEconomics/econ-project-templates: reject: research-domain template, not evidence-cited Armature setup.
- TheLartians/ModernCppStarter: reject: C++ reproducibility, not repeatable domain/stack discipline selection.
- the-turing-way/reproducible-project-template: reject: reproducible research template, no operator approval stop.
- alan-turing-institute/reproducible-project-template: reject: predecessor research template, no value citations.
- AaronGullickson/research-template: reject: Quarto research template, wrong setup coverage.
- timtroendle/cookiecutter-reproducible-research: reject: research generator, not protected rule files.
- Robinlovelace/reproducible-project-template: reject: Quarto slide template, wrong domain.
- lachlandeer/snakemake-econ-r: reject: economics workflow, not controller setup state.
- katiem0/gh-migrate-rulesets: keep: operator-readable ruleset CSV and read/export pattern.
- kuhlman-labs/repo-ruleset-bot: keep: explicit privileged reconciliation boundary and JSON export.
- int128/wait-for-workflows-action: reject: wait for checks, no App pin or enforcement proof.
- davenicoll/github-goodies: reject: unspecific examples, not verified rules.
- peterdresslar/rulesets: reject: sample rulesets only, no readback/probe.
- ghsioux/cool-rulesets: reject: sample rulesets only, no proof.
- chainloop-dev/chainloop: reject: signed agent/build records in separate control plane, not project Git setup.
- sgroy10/speclock: reject: agent hook/CI guardrails, not immutable forge controls.
- TheArchitectit/AIGGP-Agentic-Framework: reject: generic quality gates, not controller identity branch restriction.
- blackwell-systems/agentic-project-scaffolding: reject: agent skill scaffolds/sets protection, but agent authority is not separated from admin.
- mak3r/claude-scaffolding: reject: multi-agent template, no independently verified protection.
- donggun-jung/multiagent-protocol: reject: local/self-built gate, not proven GitHub hard protection against same agent.
- aapp-protocol/aapp-kit: reject: worktree/git guardrails, no server-side rules readback.
- LevyBytes/AI-SKILL-github: reject: reference skill, not a setup solution.
- nilorg/ntpl: reject: Git-template sync, no documented pinned evidence/transaction proof; git-scaffold is stronger.
- actions-ext/copier-update: reject: wraps Copier update in Action, not independent controller.
- actions-ext/copier: reject: Action wrapper, not native Go/Git pattern beyond Copier.
- liblaf/copier-shared: reject: particular template, Copier is the general pattern.
- thesmart/repo-copier: reject: skill wrapper, not cited setup record.
- piotrgredowski/yet-another-python-project-template: reject: Python project template, wrong scope.
- EmanuelVogt/platform-template: reject: product platform template, not discipline evidence.
- RaynOlivia/copier-template-mcp: reject: MCP wrapper, wrong dependency and authority.
- pyfar/template-pyfar-package: reject: Python-package template, wrong scope.
- vladmesh/service-template: reject: microservice template, not source-cited discipline.
- kasi-x/python-copier-template-example: reject: Copier usage example, not separate solution.

## What no candidate covers

No checked candidate combines **all** of the following in one public solution:

1. A programmatic, pinned Armature baseline setup adapted to domain and stack, with **every** configuration value attributed to an answer, catalog item or fact, and every unknown left explicitly open.
2. A restartable, controller-owned sequence stopping at **each** material human decision and validating each agent-written free-text section in code; neither template answers nor generic task receipts prove the value-to-source invariant.
3. The setup's evidence and decisions committed in the **project's own Git** on a branch agents cannot update, while the controller can update it without bypassing the **separate gate rules** that apply to the controller's probe.
4. A human-admin-applied, program-produced GitHub ruleset with required checks pinned by `integration_id`, PR-only updates, non-fast-forward restriction, no controller bypass of those gate checks, then a readback of *effective* rules (including inherited rules/bypass actors).
5. A reproducible negative `git push` as the controller installation identity to `layup-probe`, definitively refused by a **rule violation**; matched failed rule-suite evidence; and paginated activity audit of who actually changed protected branches.
6. A project-repository gate run with the controller absent, plus a verified second setup run on a repository created by an earlier setup version. Template update tools cover parts of this but cannot establish those Armature-specific tests.

The most consequential uncovered distinction is enforcement location: a local hook or a CI definition in agent-writable Git is **not** an external authorization barrier. GitHub rulesets supply the barrier, but the public examples do not document the complete two-identity/two-scope setup and its first-demo proof.

## Limits

- This was documentation and source-metadata review, not a live GitHub test. No test repository, App installation, admin-issued token, real ruleset, rejected `layup-probe` push, or `setup/forge-check.tsv` existed here to inspect; therefore no gate is marked “passed.”
- Repository `pushed_at` is last pushed activity, not necessarily latest authored commit or release; GitHub documentation has no stated release date on the pages reviewed. A repository description is not evidence of working functionality.
- Public ruleset documentation describes interfaces, not whether this particular repository plan, visibility, organization policy and App installation are eligible. Read endpoints and organization audit-log endpoints have different authorization requirements; verify with the real installation and operator, not by extrapolating an App's contents permission to admin APIs.
- Effective-rules readback lists active rules; evaluate/disabled rules and unexpected broader or inherited bypasses require separately inspecting rulesets. API result completeness, pagination and transient failures need a live verification record.
- The complete curated lists were read as supplied; judgments are selective only for entries touching the brief. Search output is noisy and time-sensitive. No candidate code was installed, run, imported, or tested.
- For papers and hosted APIs with no license/language/date stated in their own docs, those fields are marked “not stated” instead of inferred from neighboring repositories.
~~~~
