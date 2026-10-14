## Release review — M2d

| Field | Value |
| ----- | ----- |
| Commit reviewed | `0a104b8fc710f45cfe3e73778518616374bae1db` |
| Reviewer | GPT-6.1 Sol on Codex CLI |
| Lens | REQ-015 and REQ-017 over each code path of the release |
| Briefed on | The Operator’s brief; `.review-in/issue.md`, `criteria.md`, `release-files.txt`, `release-check.sh` and `release-check.txt`. Read all 90 release paths: 81 Go files, 19,854 lines; `go.mod`, 3 lines; eight embedded files, 297 lines. Total: 20,154 release lines. Also read repository instructions, engineering discipline, issue workflow, ADR-0015, the package specification, supporting checker code and the two local baseline-script copies. |
| Independence claimed | Context: fresh session; no earlier release review, release verdict, task plan or author decision record opened. Execution: separate review run, with this record; no release code written or changed. Model: GPT-6.1 Sol differs from both stated authors, Claude Opus 5.5 and Claude Sonnet 5.5. |
| REQ-015 | `holds` |
| REQ-017 | `holds` |

### The scenario (uat)

```text
Scenario: The release of M2d holds REQ-015 and REQ-017

Given the code of main at 0a104b8, the release of M2d
When a reviewer whose model wrote none of it reads each code path, with the output of the release check
Then no code path modifies base LLM weights or trains a model (REQ-015)
And no code path calls a cloud provider or hosting platform to modify it (REQ-017)

Accepted by GPT-6.1 Sol on Codex CLI, on 2026-10-14
Covers REQ-015, REQ-017
```

This fresh session stands for the person of the uat. Its model wrote none of the release. I read the whole specified release, reran the supplied check, searched independently, and built and vetted the code. I made no live forge or harness call. I opened no earlier release review, task plan or author reasoning record. The supplied script’s header references plans and older records (`.review-in/release-check.sh:2`, `:10`); I did not follow those references or use them as a verdict.

### The reading

For REQ-015, “modifies” means changing a base model’s learned parameters, including fine-tuning them, or requesting such a change from a service. Training a custom foundation model is also excluded. Selecting an existing model, constructing prompts, running inference and recording usage do not change its weights. This follows the distinction in the requirement’s words (`.review-in/criteria.md:8`, `:10`).

For REQ-017, “to modify it” means a call whose requested operation changes the provider or hosting platform’s configuration or infrastructure. This includes administrative settings, protection rules, provisioning and deployment. Ordinary use of hosted repository data—reading it, posting issues and comments, issuing an access token, and pushing target-repository commits or refs—is outside that meaning. Such operations do change hosted state. They do not configure the platform. This reading gives “it” its stated object, the downstream provider or platform (`.review-in/criteria.md:9`, `:11`). The supplied specification explicitly anticipates remote pushes and the M2d probe (`.review-in/criteria.md:23`, `:42`).

I include operations selected by LAYUP’s code, its embedded executable content and instructions it deliberately gives a harness. I do not treat every possible action of an Operator-supplied executable, target test or arbitrary harness response as authored release behavior. That boundary is material: the supplied criteria retain the session’s reach as the Operator’s user (`.review-in/criteria.md:48`), and ADR-0015 records the same limit.

A stricter reading that counts every write to hosted state would make REQ-017 fail on token creation, issues, comments and accepted Git pushes, including the probe. Counting generated commands executed by the Operator would also catch the ruleset applies. A stricter guarantee over every possible child-process behavior could not be established here: harnesses and target commands are not confined to these two rows. These paths are notes below. A looser reading limited to direct Go operations would discard those boundaries; it would change neither verdict and would understate the release’s reach.

Since M2c, 29 non-test Go files changed: 3,163 lines added and 115 removed. M2d adds marker proposals, setup and prose sessions, work-area reconstruction, setup verification and records, Operator push/apply commands, forge rule and activity reads, and the probe push. It adds no package (`docs/spec/packages.md:208`). This review covers the older release code too. **Zero findings name a defect under the reading above.**

### Coverage

| Path | Lines | Read |
| ---- | ----- | ---- |
| `cmd/layup/main.go` | 12 | `in full` |
| `internal/catalog/catalog.go` | 415 | `in full` |
| `internal/catalog/embed.go` | 19 | `in full` |
| `internal/cli/args.go` | 161 | `in full` |
| `internal/cli/cli.go` | 382 | `in full` |
| `internal/cli/progress.go` | 99 | `in full` |
| `internal/cli/run.go` | 232 | `in full` |
| `internal/forge/forge.go` | 158 | `in full` |
| `internal/forge/github/github.go` | 460 | `in full` |
| `internal/forge/register.go` | 174 | `in full` |
| `internal/gate/gate.go` | 196 | `in full` |
| `internal/gate/manifest.go` | 79 | `in full` |
| `internal/gate/result.go` | 50 | `in full` |
| `internal/gate/scope.go` | 48 | `in full` |
| `internal/gate/scratch.go` | 204 | `in full` |
| `internal/git/git.go` | 537 | `in full` |
| `internal/ledger/budget.go` | 93 | `in full` |
| `internal/ledger/ledger.go` | 178 | `in full` |
| `internal/psb/check.go` | 220 | `in full` |
| `internal/records/intake.go` | 411 | `in full` |
| `internal/records/records.go` | 328 | `in full` |
| `internal/records/scaffold.go` | 116 | `in full` |
| `internal/records/session.go` | 362 | `in full` |
| `internal/records/start.go` | 309 | `in full` |
| `internal/route/admit.go` | 72 | `in full` |
| `internal/route/models.go` | 155 | `in full` |
| `internal/route/register.go` | 141 | `in full` |
| `internal/rules/rules.go` | 172 | `in full` |
| `internal/run/answers.go` | 221 | `in full` |
| `internal/run/batch.go` | 351 | `in full` |
| `internal/run/command.go` | 48 | `in full` |
| `internal/run/completeness.go` | 172 | `in full` |
| `internal/run/draft.go` | 168 | `in full` |
| `internal/run/end.go` | 256 | `in full` |
| `internal/run/followanswers.go` | 239 | `in full` |
| `internal/run/followup.go` | 325 | `in full` |
| `internal/run/forgecheck.go` | 437 | `in full` |
| `internal/run/lease.go` | 175 | `in full` |
| `internal/run/markers.go` | 239 | `in full` |
| `internal/run/meaning.go` | 186 | `in full` |
| `internal/run/numbering.go` | 209 | `in full` |
| `internal/run/phase.go` | 242 | `in full` |
| `internal/run/probe.go` | 446 | `in full` |
| `internal/run/prose.go` | 108 | `in full` |
| `internal/run/push.go` | 176 | `in full` |
| `internal/run/rules.go` | 68 | `in full` |
| `internal/run/scaffold.go` | 181 | `in full` |
| `internal/run/scaffoldpush.go` | 204 | `in full` |
| `internal/run/session.go` | 394 | `in full` |
| `internal/run/setup.go` | 428 | `in full` |
| `internal/run/setupanswers.go` | 121 | `in full` |
| `internal/run/steps.go` | 871 | `in full` |
| `internal/run/store.go` | 319 | `in full` |
| `internal/run/texts.go` | 68 | `in full` |
| `internal/run/verify.go` | 225 | `in full` |
| `internal/run/workarea.go` | 132 | `in full` |
| `internal/session/checks.go` | 156 | `in full` |
| `internal/session/end.go` | 257 | `in full` |
| `internal/session/phase.go` | 127 | `in full` |
| `internal/session/process.go` | 350 | `in full` |
| `internal/session/session.go` | 144 | `in full` |
| `internal/setup/final.go` | 723 | `in full` |
| `internal/setup/phase2.go` | 151 | `in full` |
| `internal/setup/scaffold.go` | 674 | `in full` |
| `internal/setup/setup.go` | 598 | `in full` |
| `internal/setup/steps.go` | 647 | `in full` |
| `internal/spec/spec.go` | 215 | `in full` |
| `internal/standin/standin.go` | 292 | `in full` |
| `internal/tsv/block.go` | 199 | `in full` |
| `internal/tsv/tsv.go` | 250 | `in full` |
| `internal/tsv/types.go` | 122 | `in full` |
| `internal/verify/adapted.go` | 276 | `in full` |
| `internal/verify/checks.go` | 270 | `in full` |
| `internal/verify/gates.go` | 179 | `in full` |
| `internal/verify/jobs.go` | 111 | `in full` |
| `internal/verify/markers.go` | 256 | `in full` |
| `internal/verify/scripts.go` | 122 | `in full` |
| `internal/verify/sources.go` | 110 | `in full` |
| `internal/verify/target.go` | 332 | `in full` |
| `internal/verify/verify.go` | 339 | `in full` |
| `internal/work/work.go` | 162 | `in full` |
| `go.mod` | 3 | `in full` |
| `internal/catalog/go/files/.github/gates.sh.tmpl` | 173 | `in full` |
| `internal/catalog/go/files/.github/workflows/gates.yml.tmpl` | 88 | `in full` |
| `internal/catalog/go/files/docs/gates/coverage-floor.txt.tmpl` | 1 | `in full` |
| `internal/catalog/go/files/go.mod.tmpl` | 3 | `in full` |
| `internal/catalog/go/fixtures/static.patch` | 11 | `in full` |
| `internal/catalog/go/fixtures/test.patch` | 13 | `in full` |
| `internal/catalog/go/gaps.tsv` | 2 | `in full` |
| `internal/catalog/go/kinds.tsv` | 6 | `in full` |

Coverage is complete. Every release chunk was read; truncated output was reread. The tracked embedded-file list matches these eight entries.

### Findings

1. **REQ-015 and REQ-017, dependencies.** `go.mod:1` names this module and `go.mod:3` names Go 1.26.9. It has no dependency requirement. The rerun’s dependency list contains only this module and standard-library packages. No linked model-training or provider-management library appears. `internal/standin/standin.go:103` builds local test repositories; no production file imports it. `holds`.

2. **REQ-015, the word model.** Models are names, context limits, admission choices and usage labels (`internal/route/models.go:17`, `internal/route/admit.go:35`). The name substitutes into one command argument (`internal/session/process.go:47`, `internal/run/session.go:288`). Usage is parsed and priced (`internal/session/end.go:169`, `internal/ledger/ledger.go:120`). None of these paths loads or updates learned parameters. Whole-word searches found no `train` or `weights`; substring hits for `train` were `constraint`. The full read found no training operation. `holds`.

3. **REQ-015 and REQ-017, programs started.** There are five start sites: Git (`internal/git/git.go:85`); a gate shell (`internal/gate/scratch.go:190`); a baseline-script shell (`internal/verify/scripts.go:61`); the register’s version command and harness command (`internal/session/process.go:33`, `:132`). My search matches check 3’s list. Their release-selected work is repository management, checks and role sessions. The external-command boundary is finding 21. `holds`.

4. **REQ-015 and REQ-017, what sessions are sent.** The harness probe requests a token and loaded-file list (`internal/run/probe.go:29`). Intake requests meaning gaps, fact spans, requirements, coverage review and answer review (`internal/run/meaning.go:113`, `internal/run/numbering.go:64`, `internal/run/draft.go:94`, `internal/run/completeness.go:71`, `internal/run/followup.go:93`). M2d requests quoted marker proposals and specified prose files (`internal/run/markers.go:104`, `internal/run/prose.go:78`). Inputs are copied into the session (`internal/session/phase.go:28`). These prompts request neither parameter changes nor provider changes. `holds`.

5. **REQ-015 and REQ-017, session output.** Results are bounded regular files and typed records (`internal/session/end.go:83`, `internal/records/session.go:353`). Phase output passes its step’s reader and check before being committed (`internal/run/phase.go:179`). Prose output must match the requested file set and be valid UTF-8 (`internal/session/phase.go:63`); its kept paths are setup inputs (`internal/run/prose.go:48`). Output is not evaluated as a command by these paths. Task artifacts are checked by hash before the push path (`internal/run/end.go:198`). `holds`.

6. **REQ-015 and REQ-017, files written.** Writes create session directories, prompts, inputs, logs, copied harness credentials, local Git objects, work areas, target source/documents and records. Representative sinks are `internal/session/session.go:55`, `:70`, `:82`, `:86`; `internal/session/process.go:124`; `internal/git/git.go:493`; `internal/run/workarea.go:72`; `internal/run/store.go:252`; and `internal/setup/setup.go:254`. Setup writes catalog files and local protection descriptions (`internal/setup/final.go:112`, `:373`). No release-selected destination is a model-weight file or provider configuration endpoint. `holds`.

7. **REQ-015 and REQ-017, environment and credentials.** Git retains only host `PATH` and `TMPDIR`, then adds explicit configuration and any App token (`internal/git/git.go:71`, `:376`). A session gets host `PATH`, fixed values, its own home/tmp, the named harness credential and register variables (`internal/session/session.go:134`). The register refuses known forge-token and SSH-agent variable names (`internal/route/register.go:48`, `:95`, `:117`). This gives the supplied operations credentials for their stated roles; it adds no weight-training or infrastructure-management operation. `holds`.

8. **REQ-015 and REQ-017, unsafe and system calls.** The release’s `unsafe` use supplies a buffer to Linux `waitid` (`internal/session/process.go:229`). Other system calls inspect ownership, restrict file opens, read exit status or signal a process group (`internal/forge/register.go:86`, `internal/route/models.go:139`, `internal/session/end.go:84`, `internal/session/process.go:246`). I found no cgo import, plugin load or hidden process/network operation through these sites. `holds`.

9. **REQ-015 and REQ-017, connections and URLs.** The sole direct networking import is `net/http` in `internal/forge/github/github.go:19`. Requests leave through `:364`. `net/url` elsewhere parses source citations (`internal/records/records.go:297`, `internal/route/models.go:56`). The production client disables environment proxies and redirects (`internal/forge/github/github.go:56`). Its base address comes from the forge register (`internal/cli/run.go:202`); the release supplies the operations below, not a model-service or deployment request. Citation URLs and workflow comments are not fetched by those parsers. `holds`.

10. **REQ-017, Installation and Token calls.** Installation reads the existing App installation and permissions (`internal/forge/github/github.go:96`). Token posts to its access-token endpoint and caches the returned credential (`:119`, `:125`). The run calls these for identity checks and Git authentication (`internal/run/steps.go:257`, `:279`, `:285`); the adapter also calls them internally (`internal/forge/github/github.go:110`, `:348`). Neither changes installation permissions or platform configuration. Token issuance is a hosted-state write under the stricter reading. `holds`.

11. **REQ-017, Repository and UserID calls.** Repository reads metadata and every branch page, including heads (`internal/forge/github/github.go:137`, `:149`). Its callers check initial state, wait for the root push, read back that push, wait for the scaffold and check the forge (`internal/run/steps.go:298`, `:402`, `:438`; `internal/run/scaffoldpush.go:184`; `internal/run/forgecheck.go:421`). UserID reads account IDs for the Operator, idea owner and App bot (`internal/forge/github/github.go:169`; `internal/run/steps.go:317`, `:322`). All are reads. `holds`.

12. **REQ-017, OpenIssue and Comment calls.** OpenIssue posts the Intake and control issues (`internal/forge/github/github.go:184`; `internal/run/steps.go:631`, `internal/run/texts.go:24`). Comment posts question batches and follow-ups (`internal/run/batch.go:312`, `:347`; `internal/run/followup.go:298`, `:321`) and task, phase and harness-probe outcomes (`internal/run/end.go:142`, `internal/run/phase.go:230`, `internal/run/probe.go:299`). The adapter sends each comment body at `internal/forge/github/github.go:199`. These change repository discussion data, not platform settings. `holds`.

13. **REQ-017, Comments calls.** The adapter reads every comment page (`internal/forge/github/github.go:224`). Callers collect answers, detect an already-posted batch or follow-up, locate the setup batch and watch for a notice (`internal/run/answers.go:154`, `internal/run/batch.go:261`, `internal/run/followanswers.go:152`, `internal/run/followup.go:289`, `internal/run/setup.go:229`, `internal/run/steps.go:655`). Human input is copied into records before its decision use (`internal/run/steps.go:687`). No provider setting is written. `holds`.

14. **REQ-017, EffectiveRules and Activity calls.** EffectiveRules reads applicable branch rules; Activity reads one activity page (`internal/forge/github/github.go:266`, `:301`). Their callers are `internal/run/forgecheck.go:422` and `:430`. The rule reads cover the default branch, `layup-probe` and `layup-records` (`:127`). Activity is checked for a successful read (`:430`). There is no adapter operation that applies a ruleset. My search matches all ten request sites in check 6. `holds`.

15. **REQ-017, Git reads and local remote setup.** `LsRemote` and `Clone` read the baseline in `internal/run/steps.go:355`, `:360`. Setup reaches the same calls through `internal/setup/setup.go:230`, `:231`, used at `internal/setup/steps.go:269`, `:274`. The restart clones the target (`internal/run/steps.go:731`); read-back and lease reads fetch target refs (`:435`, `internal/run/store.go:50`). `FetchSession` reads local session objects (`internal/run/end.go:203`, `internal/git/git.go:469`). Local clones are made at `internal/session/session.go:58` and `internal/run/workarea.go:81`; `AddRemote` at `:102` only adds local configuration. These are all ten non-push references in check 7. None writes a provider. `holds`.

16. **REQ-017, records pushes.** `commitFiles` pushes only to `layup-records`, with the installation token and no force (`internal/run/store.go:278`, `:282`; `internal/git/git.go:403`). It is reached through commit and removal wrappers (`internal/run/store.go:208`, `:215`), including lease writes (`:83`), Start (`internal/run/steps.go:581`), fenced step writes (`:544`), session records (`internal/run/session.go:184`), prose removal (`internal/run/setup.go:364`), verification, applies and forge-check records (`internal/run/verify.go:222`, `internal/run/scaffoldpush.go:201`, `internal/run/forgecheck.go:433`). These are target-repository record writes. `holds`.

17. **REQ-017, task-head pushes.** The only caller of `PushHead` is `pushAndBind` (`internal/run/push.go:148`). It selects `task/<task>/<attempt>` and records the push before binding an accepted head (`:137`, `:145`, `:154`). Its upstream checks require an accepted result, ancestry and permitted paths (`:61`, `:69`, `:111`). Workflow and rule-path changes are refused (`internal/rules/rules.go:155`). `PushHead` uses the App token for the target URL (`internal/run/store.go:178`, `:182`). This writes source-control data, not platform configuration. `holds`.

18. **REQ-017, the forge probe push.** `pushProbe` makes a new commit with the default head’s exact tree and that head as parent, then pushes its SHA to `refs/heads/layup-probe` (`internal/run/forgecheck.go:28`, `:388`, `:392`, `:396`; `internal/git/git.go:404`). The caller supplies the target’s own repository URL, the App bot identity and an installation token (`internal/run/forgecheck.go:423`, `:428`; `internal/run/steps.go:251`, `:256`). A successful push produces `fail`, with the accepted SHA (`internal/run/forgecheck.go:369`). The table is committed and the step fails (`:76`, `:80`). No deletion follows. A refused push passes only when its nonzero error contains `GH013` and the expected ruleset name (`:374`). An accepted ref changes the target repository under my reading, not the platform. `holds`.

19. **REQ-017, hosted writes under the stricter reading.** Token issuance, issues, comments and the three push paths really write remote state (`internal/forge/github/github.go:119`, `:184`, `:199`; `internal/run/store.go:182`, `:282`; `internal/run/forgecheck.go:396`). The probe’s expected refusal does not make its accepted path disappear. Counting every hosted-state mutation would turn these into defects. **No change** under the stated reading. `note`.

20. **REQ-017, commands written for the Operator.** Setup constructs default-branch and records-branch rulesets (`internal/setup/final.go:215`, `:278`). Generated commands push setup commits and execute `gh api --method POST .../rulesets`, using the Operator’s login (`:403`, `:404`, `:428`, `:443`). Other hand-offs push root and records commits (`internal/setup/steps.go:352`, `:353`; `internal/setup/final.go:718`; `internal/run/steps.go:399`). LAYUP writes the command file (`internal/setup/setup.go:460`) and prints it while waiting (`internal/run/scaffoldpush.go:70`); it does not execute it. The ruleset applies would count if Operator execution of generated instructions were attributed to LAYUP. **No change** under the stated reading. `note`.

21. **REQ-015 and REQ-017, external-command reach.** Harness commands and fixed variables are Operator-supplied text (`internal/route/register.go:25`, `:34`). The process runs as the Operator’s user; its separate home and process group are not an access-control sandbox (`internal/session/session.go:43`, `internal/session/process.go:132`, `:134`; `docs/adr/0015-keep-model-calls-out-of-the-engine-checks.md:71`). Gate and baseline shells inherit LAYUP’s environment (`internal/gate/scratch.go:190`, `internal/verify/scripts.go:61`). A harness, target test or supplied shell command could train weights or call a provider using that user’s access. The release does not prevent every such action. **No change** for the retained reach limit; the verdict is not a confinement guarantee. `note`.

22. **REQ-015 and REQ-017, gate and catalog commands.** The gate reads its command from the base manifest and runs it in a scratch tree with base gate files (`internal/gate/gate.go:105`, `:128`, `:151`; `internal/gate/result.go:33`). The embedded Go catalog supplies formatting/vetting and ordinary Go tests (`internal/catalog/go/kinds.tsv:2`, `:6`). Its fixtures add a formatting fault and a failing test (`internal/catalog/go/fixtures/static.patch:7`, `internal/catalog/go/fixtures/test.patch:7`). These supplied commands perform no model training or infrastructure modification. Arbitrary target commands remain finding 21. `holds`.

23. **REQ-015 and REQ-017, baseline scripts.** Verify invokes `docs/tests/run-discipline-tests.sh` and `docs/links/link-lint.sh` from the target tree (`internal/verify/scripts.go:20`, `:21`, `:61`; `internal/verify/verify.go:262`). The local copies dispatch fixture checks and resolve local links (`docs/tests/run-discipline-tests.sh:288`, `:305`; `docs/links/link-lint.sh:137`, `:266`). They are not embedded release bytes. The fetched baseline and their subprocesses remain an external-code boundary; reading these local copies does not establish every future baseline’s behavior. **No change**. `note`.

24. **REQ-015 and REQ-017, the workflow setup writes.** S12 writes the embedded catalog workflow (`internal/setup/final.go:56`, `:112`; `internal/catalog/embed.go:12`). It declares `pull_request`, `contents: read`, checkout, Go setup and five gate jobs (`internal/catalog/go/files/.github/workflows/gates.yml.tmpl:7`, `:10`, `:18`, `:22`, `:25`, `:85`). Its shell runs the manifest’s command (`internal/catalog/go/files/.github/gates.sh.tmpl:170`). It contains no deployment, provider configuration or training command. Hosted check execution and target-test behavior have the boundary in finding 21. `holds`.

25. **REQ-015 and REQ-017, release-check rules 1 and 2.** Checking the reviewed HEAD and dirty tree is right (`.review-in/release-check.sh:75`). This run used the full supplied SHA and a clean tree. `TestPackageRules` is useful and passed (`:83`). It checks module dependencies, permitted imports, program starts and network reach (`cmd/layup/rules_checker_test.go:322`). It also reads ignored Go files (`cmd/layup/rules_integration_test.go:110`). Allowing network dependencies in the adapter, CLI and entry point is right: only the adapter imports networking directly; the latter two depend on it (`docs/spec/packages.md:127`, `:129`, `:130`). These checks do not decide endpoint semantics. **No change**. `note`.

26. **REQ-015 and REQ-017, release-check rules 3 and 4.** Rule 3’s five file/program allowances match my search and the actual starts (`.review-in/release-check.sh:97`, `:101`). Rule 4 correctly confines `push` to `Push`, `fetch` to `Fetch` and `FetchSession`, and `remote` to `CloneLocal` and `AddRemote` (`:125`, `:135`). The latter two perform local `remote remove` and `remote add`, without fetching (`internal/git/git.go:443`, `:452`). The allowance is right. Rule 4 is a literal-verb check; it does not enumerate every Git capability. **No change**. `note`.

27. **REQ-015 and REQ-017, release-check lists 5 and 6.** Listing embedded entries, catalog commands and baseline-script names is right as evidence, not a pass condition (`.review-in/release-check.sh:147`). The lists match my search and full read. List 6 matches the adapter’s ten request sites (`:157`); the listed operations are assessed in findings 10–14. Neither list reads arbitrary target scripts, harness executables or Operator execution of generated commands. **No change**. `note`.

28. **REQ-017, release-check rule 7.** The three allowed `git.Push` references are exactly the two store sinks and the probe sink (`.review-in/release-check.sh:180`; `internal/run/store.go:182`, `:282`; `internal/run/forgecheck.go:396`). Refusing aliased `internal/git` imports is right for this selector-based check (`.review-in/release-check.sh:189`). The read-reference list also matches my search. Its two `no function` entries are the setup system initializer, whose callers I traced in finding 15. Local clone/config calls in that list do not actually contact the forge. **No change**. `note`.

29. **REQ-015 and REQ-017, what the check cannot see.** The shell checks match lines and literals, not full behavior (`.review-in/release-check.sh:97`, `:125`, `:170`). A changed operation inside an allowed function, a constructed Git verb, a differently written selector, or another HTTP endpoint could escape their intended meaning. Program and network allowances do not inspect child behavior. The package checker improves start detection with an import-aware syntax scan (`cmd/layup/rules_checker_test.go:432`), but excludes cgo and raw system calls from that analysis (`docs/spec/packages.md:255`). Its documentation records the dynamic harness limit (`docs/spec/packages.md:175`). Full reading closes these gaps for this release’s code; it does not turn the check into proof over external programs. **No change**. `note`.

### Runs

Searches below used the non-test release Go pathspecs, excluded `testdata`, and included `internal/catalog/go/*` where relevant. Search patterns are shown without repeating those pathspecs.

| Command | Result |
| ------- | ------ |
| `git rev-parse HEAD` | Exit 0; exact reviewed SHA. |
| `git status --porcelain --untracked-files=normal` | Exit 0; empty before and after review. |
| `nl -ba PATH` with `sed -n` chunks of at most 250 lines | All 90 release paths read in full. Truncated output reread. |
| `GOCACHE="$PWD/.gocache" GOTMPDIR="$PWD/.tmp" sh .review-in/release-check.sh 0a104b8fc710f45cfe3e73778518616374bae1db` | Exit 0; `== PASS`; `TestPackageRules` passed in 2.419 s. Lists match the supplied output and my searches. |
| `go build ./...` | Exit 0. Used local cache/tmp, cleared `GOFLAGS`, disabled `GOENV`, `GOWORK`, toolchain switching and the module proxy. |
| `go vet ./...` | Exit 0, with the same environment. |
| `git diff --stat 0b881fa 0a104b8fc710f45cfe3e73778518616374bae1db -- cmd internal go.mod` | Exit 0; includes tests: 75 files, 9,696 additions, 228 removals. |
| Same diff with non-test Go and `testdata` exclusions | Exit 0; 29 files, 3,163 additions, 115 removals. |
| `git grep -n -E 'exec\.|os\.StartProcess|syscall\.'` and the focused starter search | Five start sites; other system-call uses assessed in finding 8. |
| `git grep -n -E '"(net(/[^"]*)?\|crypto/tls)"\|\b(net\|http\|tls)\.'` | Networking confined to the adapter; URL parsers also found. |
| `git grep` for environment readers, `cmd.Env`, forge tokens and `SSH_AUTH_SOCK` | Two explicit host-environment readers; shell inheritance and register restrictions assessed. |
| `git grep` for file creation, writes, removal, rename and mode changes | Write sinks traced; findings 5–7. |
| `git grep` for URLs, `unsafe`, cgo, `go:linkname`, plugins and reflection | URL uses and `waitid` assessed; no cgo import or dynamic code-loading path found. |
| `git grep -n -i -E 'model\|train\|weights'` | Model labels, routing and usage; no training operation found. |
| `git grep -n -i -w -E 'train\|weights'` | Exit 1: no whole-word matches. |
| `git grep` for training, fine-tuning, foundation models, provider tools and deployment | No relevant operation; incidental substring matches only. |
| `git grep` for every forge method, `a.do`, request creation and `Client.Do` | Ten request sites; all callers traced in findings 10–14. |
| `git grep` for remote Git calls, Git imports, internal command wrappers and their callers | Thirteen references outside `internal/git`, including three pushes; findings 15–18. |
| `git ls-files 'internal/catalog/go/*'` | Exactly the eight embedded paths in Coverage. |

No command was denied by the sandbox. No tracked file was changed. Go used `.gocache` and `.tmp`. The full test suite, discipline suite and live release demo were not run by this reviewer.

*Written by GPT-6.1 Sol on Codex CLI (codex exec), the reviewer of the release of M2d (#283).*

*Posted for a reviewer session (GPT-6.1 Sol on Codex CLI, codex exec, effort xhigh, a fresh session in a clone at 0a104b8 with `runs/`, `docs/tasks/`, `docs/prd/` and `docs/plan/` out of its working tree, from 09:53:22 UTC, 16 min 35 s; tokens used 443,462) through the layup-agent App, by a Claude Sonnet 5.5 session on Claude Code, the author of this task. The harness codex is not one of the three that the routing rule of ADR-0026 names (claude, devin, opencode); the Operator added Codex CLI to the host on 2026-10-13 and said it may be used. The text above is the session's final message, unchanged.*
