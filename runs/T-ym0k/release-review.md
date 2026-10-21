## Release review — M2e

| Field | Value |
| ----- | ----- |
| Commit reviewed | `01d7da3082333efadb227afda14fd98715edf5c4` |
| Reviewer | Claude Fable 5.1 (Anthropic), in OpenCode |
| Lens | REQ-015 and REQ-017 over each code path of the release |
| Briefed on | The 149 release paths of `.review-in/release-files.txt`, each in full (chunks of at most 250 lines, by the harness's file reader). Besides the release: the seven files of `.review-in/` (`brief.md`, `criteria.md`, `issue.md`, `record-form.md`, `release-check.sh`, `release-check.txt`, `release-files.txt`); `AGENTS.md` (given as the session's instructions); `opencode.json` of the clone (the harness's shell allowlist, read only to learn why commands were denied). No file of `docs/spec/` or `docs/adr/` was opened beyond the verbatim extracts in `criteria.md` (packages.md "The table of M2e", forge.md "The calls of M2e", ADR-0015, the PRD rows). Test files were not read; `git grep` hits in `cmd/layup/loop_e2e_test.go:429`, `scaffold_e2e_test.go:192` and `whole_e2e_test.go:734` were used only to show that `out/commands.sh` is run by tests, never by a release path. The standard library was used from memory; no file of it was opened. |
| Independence claimed | **Context:** a fresh session holding only the brief, the clone and the files named above; no plan, plan review, answer or earlier verdict was looked for or seen. **Execution:** read-only — no file was written in the clone (`git status --porcelain --untracked-files=normal` is empty before and after); no `git log`, `git show`/`git diff` of another commit, `git blame`, other branch, `docs/tests/traceability.md`, GitHub or network was used, with one caveat under Runs (the Go toolchain cache). **Model:** Claude Fable 5.1, which differs from the author of the issue (`.review-in/issue.md:32`: Claude Sonnet 5.5); the authorship of the release's commits themselves could not be checked without git history, which the brief forbids, so the claim rests on the brief's own statement. |
| REQ-015 | `holds` |
| REQ-017 | `holds` |

### The scenario (uat)

```text
Scenario: The release of M2e holds REQ-015 and REQ-017

Given the code of main at 01d7da3, the release of M2e
When a reviewer whose model wrote none of it reads each code path, with the output of the release check
Then no code path modifies base LLM weights or trains a model (REQ-015)
And no code path calls a cloud provider or hosting platform to modify it (REQ-017)

Accepted by Claude Fable 5.1 in OpenCode, on 2026-10-21
Covers REQ-015, REQ-017
```

### The reading

**A code path** is any execution path reachable from `cmd/layup/main.go:1-12` through the non-test Go files of the release: the binary's own functions, and the programs it starts (`sh`, `git`, a harness) with exactly the arguments and the environment it gives them. Text that the code only emits — `out/commands.sh` and `land-code-owner.sh` for the Operator, the catalog files written into a target at S12 and run by the target's CI — is a code path only as far as a release path executes it, and none does (finding 4). A role session is a harness process (ADR-0015 item 3); what the harness does with its own credential is the harness's, bounded by the environment LAYUP gives it (finding 6). Test files are not the release.

**REQ-015** ("no code path modifies base LLM weights or trains a custom foundation model") I read as: no path writes a model's parameters, runs a training or fine-tuning procedure, uploads training data, or calls a provider's training or fine-tuning API. Inference by a harness is not a modification of weights. A stricter reading that counted prompting or in-context learning as "training" would be a requirement on the harness, not on LAYUP's code, and would not change the verdict, because the code sends no training data and holds no model client; a looser reading changes nothing.

**REQ-017** ("no code path calls a cloud provider or hosting platform to modify it") I read "modify it" as changing the platform's configuration or infrastructure: repository settings, rulesets and branch protection, workflows and Actions settings, secrets, App and installation permissions, collaborators, webhooks, organisation settings, or any cloud resource. Writing repository content and its workflow objects — a branch by push, an issue, a comment, a pull request and its draft state, a commit status, a merge, a close — is not "modifying the platform": it is the product's job (the PSB, REQ-005/REQ-007), and the ruleset the Operator applies governs it. Under the strictest reading, where every write to hosted state is a modification, REQ-017 fails by design for every forge write of M2a and M2e, and the acceptance row contradicts the product, so that reading cannot be the intended one. Under an intermediate reading that counts writes to protected state — the merge into the default branch (finding 2), the creation of `layup-probe` by the probe push (finding 7) — those are still content writes under the ruleset, and the verdict stands, with finding 7 becoming a `risk`. A stricter reading of "code path" that counted emitted text would make the `gh api --method POST …/rulesets` lines of `internal/setup/final.go:427` and `:469` and the `PUT …/rulesets/<id>` of `internal/setup/codeowner.go:126`, `:165-167` a platform modification by LAYUP — but a hand runs them, by design (the App has no administration permission; the Operator applies the rulesets, `criteria.md`'s ADR-0015 item 5 and forge.md), so under that reading REQ-017 would fail by design rather than by a defect; I do not adopt it. Under a looser reading (only the provisioning of cloud resources counts) REQ-017 holds trivially.

### Coverage

| Path | Lines | Read |
| ---- | ----- | ---- |
| `cmd/layup/main.go` | 12 | in full |
| `internal/catalog/catalog.go` | 415 | in full |
| `internal/catalog/embed.go` | 19 | in full |
| `internal/cli/args.go` | 161 | in full |
| `internal/cli/cli.go` | 382 | in full |
| `internal/cli/progress.go` | 99 | in full |
| `internal/cli/run.go` | 232 | in full |
| `internal/forge/fake.go` | 358 | in full |
| `internal/forge/forge.go` | 286 | in full |
| `internal/forge/github/github.go` | 767 | in full |
| `internal/forge/register.go` | 174 | in full |
| `internal/gate/gate.go` | 196 | in full |
| `internal/gate/manifest.go` | 79 | in full |
| `internal/gate/result.go` | 50 | in full |
| `internal/gate/scope.go` | 48 | in full |
| `internal/gate/scratch.go` | 204 | in full |
| `internal/git/git.go` | 695 | in full |
| `internal/handoff/check.go` | 145 | in full |
| `internal/handoff/table.go` | 172 | in full |
| `internal/ledger/budget.go` | 113 | in full |
| `internal/ledger/ledger.go` | 178 | in full |
| `internal/psb/check.go` | 220 | in full |
| `internal/records/intake.go` | 414 | in full |
| `internal/records/records.go` | 328 | in full |
| `internal/records/scaffold.go` | 116 | in full |
| `internal/records/session.go` | 435 | in full |
| `internal/records/start.go` | 309 | in full |
| `internal/route/admit.go` | 72 | in full |
| `internal/route/models.go` | 155 | in full |
| `internal/route/register.go` | 141 | in full |
| `internal/rules/basic.go` | 474 | in full |
| `internal/rules/rules.go` | 173 | in full |
| `internal/run/answer.go` | 155 | in full |
| `internal/run/answers.go` | 417 | in full |
| `internal/run/archtasks.go` | 150 | in full |
| `internal/run/attempt.go` | 156 | in full |
| `internal/run/authors.go` | 200 | in full |
| `internal/run/baseup.go` | 306 | in full |
| `internal/run/batch.go` | 381 | in full |
| `internal/run/bet.go` | 600 | in full |
| `internal/run/betbasis.go` | 98 | in full |
| `internal/run/betbrief.go` | 787 | in full |
| `internal/run/bets.go` | 88 | in full |
| `internal/run/closeout.go` | 320 | in full |
| `internal/run/codeowner.go` | 102 | in full |
| `internal/run/command.go` | 48 | in full |
| `internal/run/commentfile.go` | 71 | in full |
| `internal/run/commentform.go` | 259 | in full |
| `internal/run/completeness.go` | 172 | in full |
| `internal/run/draft.go` | 168 | in full |
| `internal/run/end.go` | 314 | in full |
| `internal/run/finish.go` | 547 | in full |
| `internal/run/followanswers.go` | 256 | in full |
| `internal/run/followup.go` | 338 | in full |
| `internal/run/forgecheck.go` | 694 | in full |
| `internal/run/frozen.go` | 383 | in full |
| `internal/run/gates.go` | 288 | in full |
| `internal/run/handoffpost.go` | 494 | in full |
| `internal/run/instructions.go` | 224 | in full |
| `internal/run/lease.go` | 175 | in full |
| `internal/run/loop.go` | 456 | in full |
| `internal/run/markers.go` | 239 | in full |
| `internal/run/meaning.go` | 186 | in full |
| `internal/run/numbering.go` | 209 | in full |
| `internal/run/options.go` | 456 | in full |
| `internal/run/parameters.go` | 505 | in full |
| `internal/run/phase.go` | 213 | in full |
| `internal/run/probe.go` | 453 | in full |
| `internal/run/promptfile.go` | 723 | in full |
| `internal/run/prose.go` | 223 | in full |
| `internal/run/pullrequest.go` | 186 | in full |
| `internal/run/push.go` | 190 | in full |
| `internal/run/qacomment.go` | 108 | in full |
| `internal/run/question.go` | 100 | in full |
| `internal/run/questionloop.go` | 292 | in full |
| `internal/run/ready.go` | 287 | in full |
| `internal/run/redo.go` | 324 | in full |
| `internal/run/register.go` | 312 | in full |
| `internal/run/rendered.go` | 104 | in full |
| `internal/run/reopen.go` | 548 | in full |
| `internal/run/repair.go` | 176 | in full |
| `internal/run/retrospective.go` | 132 | in full |
| `internal/run/reviewlink.go` | 150 | in full |
| `internal/run/revision.go` | 377 | in full |
| `internal/run/rules.go` | 68 | in full |
| `internal/run/scaffold.go` | 181 | in full |
| `internal/run/scaffoldpush.go` | 204 | in full |
| `internal/run/session.go` | 505 | in full |
| `internal/run/setup.go` | 458 | in full |
| `internal/run/setupanswers.go` | 121 | in full |
| `internal/run/shape.go` | 117 | in full |
| `internal/run/shapeloop.go` | 850 | in full |
| `internal/run/shapequestions.go` | 364 | in full |
| `internal/run/shapetask.go` | 195 | in full |
| `internal/run/specinput.go` | 209 | in full |
| `internal/run/specstatus.go` | 173 | in full |
| `internal/run/statuses.go` | 268 | in full |
| `internal/run/steps.go` | 984 | in full |
| `internal/run/store.go` | 347 | in full |
| `internal/run/texts.go` | 68 | in full |
| `internal/run/tier.go` | 256 | in full |
| `internal/run/treeitems.go` | 97 | in full |
| `internal/run/verify.go` | 225 | in full |
| `internal/run/verifystatus.go` | 111 | in full |
| `internal/run/workarea.go` | 132 | in full |
| `internal/session/checks.go` | 156 | in full |
| `internal/session/context.go` | 79 | in full |
| `internal/session/end.go` | 269 | in full |
| `internal/session/phase.go` | 127 | in full |
| `internal/session/process.go` | 350 | in full |
| `internal/session/session.go` | 144 | in full |
| `internal/setup/codeowner.go` | 269 | in full |
| `internal/setup/final.go` | 749 | in full |
| `internal/setup/phase2.go` | 153 | in full |
| `internal/setup/scaffold.go` | 674 | in full |
| `internal/setup/setup.go` | 600 | in full |
| `internal/setup/steps.go` | 647 | in full |
| `internal/spec/answer.go` | 172 | in full |
| `internal/spec/bet.go` | 214 | in full |
| `internal/spec/full.go` | 830 | in full |
| `internal/spec/inventory.go` | 181 | in full |
| `internal/spec/prdlint.go` | 99 | in full |
| `internal/spec/quote.go` | 154 | in full |
| `internal/spec/render.go` | 172 | in full |
| `internal/spec/shape.go` | 231 | in full |
| `internal/spec/spec.go` | 215 | in full |
| `internal/standin/standin.go` | 292 | in full |
| `internal/tsv/block.go` | 199 | in full |
| `internal/tsv/tsv.go` | 250 | in full |
| `internal/tsv/types.go` | 122 | in full |
| `internal/verify/adapted.go` | 356 | in full |
| `internal/verify/checks.go` | 275 | in full |
| `internal/verify/gates.go` | 179 | in full |
| `internal/verify/jobs.go` | 111 | in full |
| `internal/verify/markers.go` | 256 | in full |
| `internal/verify/scripts.go` | 122 | in full |
| `internal/verify/sources.go` | 110 | in full |
| `internal/verify/target.go` | 344 | in full |
| `internal/verify/verify.go` | 339 | in full |
| `internal/work/work.go` | 162 | in full |
| `go.mod` | 3 | in full |
| `internal/catalog/go/files/.github/gates.sh.tmpl` | 246 | in full |
| `internal/catalog/go/files/.github/workflows/gates.yml.tmpl` | 126 | in full |
| `internal/catalog/go/files/docs/gates/coverage-floor.txt.tmpl` | 1 | in full |
| `internal/catalog/go/files/go.mod.tmpl` | 3 | in full |
| `internal/catalog/go/fixtures/static.patch` | 11 | in full |
| `internal/catalog/go/fixtures/test.patch` | 13 | in full |
| `internal/catalog/go/gaps.tsv` | 2 | in full |
| `internal/catalog/go/kinds.tsv` | 6 | in full |

The coverage is complete — all 149 paths read in full, their set confirmed equal to `git ls-files 'cmd/*.go' 'internal/*.go' ':!*_test.go' ':!*/testdata/*'` plus `go.mod` plus `git ls-files internal/catalog/go`, and their line counts to `wc -l` (38,708 in all) — and I read again, while writing the findings, `internal/spec/prdlint.go:20-69`, `internal/gate/scratch.go:180-204`, `internal/verify/scripts.go:50-74`, `internal/forge/github/github.go:43-66`, `internal/session/process.go:24-39` and `:126-137`, `internal/run/push.go:40-99` and `internal/run/scaffoldpush.go:180-204`.

### Findings

1. **(REQ-015, the search for training and weights).** The module has no dependency outside the standard library (`go.mod:1-3` holds no `require`, and `go build ./...` passes, so no import can leave the module; check (2) lists only `github.com/pharzam/layup/...` packages). No release path imports a machine-learning or model-provider library, opens a model or checkpoint file, or names weights, fine-tuning, training or a model's parameters; the only vocabulary hits were the "draws" of the PSB check and the prose of `internal/run/betbrief.go`. No path links a harness SDK (ADR-0015 rejects it). `holds`.

2. **(REQ-015, the processes that reach a model).** The engine checks start four things and nothing else (check (3), confirmed by `git grep` of `exec.` over every non-test file — every other hit is a `_test.go`): `sh -c <gate command>` (`internal/gate/scratch.go:190`), `sh <prd-lint> <dir>` (`internal/spec/prdlint.go:32`, called from `internal/run/shapeloop.go:789-816` with the script of the base), `sh <baseline script>` (`internal/verify/scripts.go:61`, for `docs/tests/run-discipline-tests.sh` and `docs/links/link-lint.sh`, `scripts.go:20-21`) and `git` (`internal/git/git.go:87`). The one process that calls a model is a role session, started from the words of a harness register row (`internal/session/process.go:31-34` for the version probe, `:132-134` for the session), which the harness drives for inference; it is given a named environment (`internal/session/session.go:135-144`) and an empty home, and LAYUP sends it no data for training. `holds`.

3. **(REQ-017, the requests of the forge adapter).** `internal/forge/github` is the only package that imports `net/http` (check (2)). Its 21 requests (check (6), `github.go:96-600`) go through `do` (`github.go:644-667`) to `a.cfg.API` of the forge register, by a client that follows no redirect and reads no proxy of the environment (`github.go:56-62`). The writes are `POST /app/installations/{id}/access_tokens` (`Token`, `:118-119`, which mints an installation token — an authentication artifact, not a change of the platform), `POST …/issues` (`:184`), `POST …/issues/{n}/comments` (`:199`), `POST …/pulls` with `draft` true (`:358`), the GraphQL mutation `markPullRequestReadyForReview` (`:391`), `POST …/statuses/{sha}` (`:422`), `PUT …/pulls/{n}/merge` with `merge_method` `merge` (`:539`), `PATCH …/pulls/{n}` with `body` alone (`:564`), `PATCH …/pulls/{n}` with `{"state":"closed"}` (`:579`) and `PATCH …/issues/{n}` with `{"state":"closed","state_reason":"not_planned"}` (`:600`), each as forge.md's table of M2e states. Every one writes repository content or a workflow object; none writes a ruleset, a branch-protection rule, a repository setting, a workflow, a secret, a permission or a collaborator. The reads (`Installation`, `Repository`, `UserID`, `Comments`, `EffectiveRules`, `Activity`, `Statuses`, `PullRequest`) modify nothing. `holds`.

4. **(REQ-017, the ruleset commands).** The code that names `gh api --method POST repos/<name>/rulesets` (`internal/setup/final.go:427`, `:469`) and `gh api --method PUT /repos/<name>/rulesets/<id>` (`internal/setup/codeowner.go:126`, `:165-168`) writes those lines as text into `out/commands.sh` (`internal/setup/setup.go:214`, path `work.CommandsPath` at `internal/work/work.go:25`) and the landing script, for the Operator to run with the Operator's own login; `internal/run/scaffoldpush.go:190-194` reads `commands.sh` only to print it, and `internal/run/codeowner.go:69` prints the `sh …/out/commands.sh` instruction. No release path starts `gh`, and no release path runs `commands.sh` (the three `exec.Command("sh", …commands.sh)` hits are in `cmd/layup/*_e2e_test.go`, not the release). The one path that changes the platform's configuration is a hand outside the code. `holds`.

5. **(REQ-017, the git verbs that reach a remote).** The verbs of check (4) and the two that it leaves to (7) — `ls-remote` in `LsRemote` (`internal/git/git.go:135`) and `RemoteHead` (`:461`), and `clone` in `Clone` (`:148`) — are reads except `push` (`Push`, `:427-433`), which pushes one SHA to `refs/heads/<branch>` with no `--force` (`:413-426`) and carries the token as `GIT_CONFIG_*` of the call's environment (`:398-409`), never as an argument; `remote remove origin` (`:515`, `CloneLocal`) and `remote add` (`:524`, `AddRemote`) write a local clone's config only, and `FetchSession` (`:541-582`) fetches from the session's clone on disk. Every git call gets the fixed environment of `:73-75` (`HOME=/dev/null`, no system or global config, `GIT_TERMINAL_PROMPT=0`, `GIT_ALLOW_PROTOCOL=file:git:http:https`, so no ssh and no remote helper that reads the home). The three callers of `Push` (check (7)) push a task head to its branch (`internal/run/store.go:182`, `PushHead`), the records commit to `layup-records` (`store.go:310`, `commitFiles`) and the probe commit (`internal/run/forgecheck.go:519`); the push hook refuses a head that changes `.github/workflows/` or a rule path before any push (`internal/run/push.go:45-56`; `internal/rules/rules.go:36`, `:158-159`), which keeps the App's `contents: write` from touching a workflow (ADR-0015 item 5). Each push is a branch write under the ruleset, repository content, not configuration. `holds`.

6. **(REQ-017, the environment of the started programs).** The session's environment is a named list (`internal/session/session.go:135-144`: `PATH` of the host, `LANG`, the empty `HOME`, `TMPDIR`, `GIT_CONFIG_NOSYSTEM`, and the harness credential only), and the route register refuses a row that names `GH_TOKEN`, `GITHUB_TOKEN`, `GH_ENTERPRISE_TOKEN`, `GITHUB_ENTERPRISE_TOKEN` or `SSH_AUTH_SOCK` (`internal/route/register.go:48-49`); the session's clone has no remote (`session.go:58`, `git.go:505-515`). The App's token lives in the adapter (`github.go:47`) and reaches a child only as the `GIT_CONFIG_*` of one git call. By contrast the gate command (`internal/gate/scratch.go:190-192`), the pinned prd-lint (`internal/spec/prdlint.go:32-34`) and the baseline scripts (`internal/verify/scripts.go:61-62`, "with the environment of layup", `:57`) set no `cmd.Env`, so they inherit the whole environment of the `layup` process: a forge credential that the Operator's own shell exports (a `GH_TOKEN` for `gh`, say) would reach `sh -c <command>` of the target's `docs/gates.tsv` and the target's `docs/tests/*.sh`. Those commands are the target's, taken from the base (`gates.sh.tmpl:12-18`, F-22) and from rule paths that the register governs (`internal/setup/final.go:519-520`; finding 5's refusal), not a code path of LAYUP; a process on the host reads what the Operator's user can read in any case (known limit L-A1, ADR-0015's consequences). A stricter reading of "code path" that counted a target's gate command would make this a `risk`, whose remedy is a named environment for the three `sh` starts, as `git` and a session already have. `note`. No change: under the reading above the inherited environment reaches no code path of LAYUP, the behaviour is the specified one, and L-A1 records the limit.

7. **(REQ-017, the probe push).** `pushProbe` (`internal/run/forgecheck.go:507-519`) pushes an empty commit on the head of the default branch to `layup-probe` so that the ruleset refuses it; a pass is only a refusal whose stderr holds `GH013` for that ref (`:483-504`, `:37`), and an accepted push is `fail` "probe accepted" (`:490`). A wrong ruleset therefore leaves a branch `layup-probe` on the forge — content, not configuration, and the check reports it. `note`. No change: the branch is the evidence of the failed check, and the probe is specified (scaffold.md, The probe, named at `:36`).

8. **(REQ-015 and REQ-017, `unsafe`, `syscall`, `plugin`, cgo).** `unsafe` is imported once (`internal/session/process.go:22`) for the raw `waitid` of `:233` (`WEXITED|WNOWAIT`, a poll of the session's process that reaps nothing). `syscall` serves process control and file metadata only: `Setpgid` (`:134`), the signals to the process group (`:245-246`), `WaitStatus` (`:220`, `internal/gate/scratch.go:198`, `internal/verify/scripts.go:69`), `Stat_t` for file owners (`internal/forge/register.go:86`, `internal/route/models.go:139`), `O_NOFOLLOW|O_NONBLOCK` (`internal/session/end.go:84`) and `ENOTDIR` (`internal/gate/scratch.go:153`). There is no `plugin` import, no `import "C"`, no cgo or `go:linkname` directive. `note`. No change: none touches a model or a platform.

9. **(REQ-017, the catalog's workflow and script).** The embedded workflow declares `permissions: contents: read` (`internal/catalog/go/files/.github/workflows/gates.yml.tmpl:13-14`) and each job runs only `actions/checkout@v4`, `actions/setup-go@v5`, the base's `gates.sh --overlay` and `sh .github/gates.sh <kind>` (`:21-38`); the script starts `sh`, `awk`, `tail`, `iconv`, `tr`, `mktemp`, `rm`, `git` and the kind's tool (`gates.sh.tmpl:7-8`), runs the manifest's command by `sh -c "$command"` (`:243`), and in `--overlay` mode changes only the index and work tree of the checkout (`:117-120`); `kinds.tsv:2` and `:6` give `gofmt`/`go vet` and `go test` as the commands. Neither file calls the forge's API, and `layup` never runs them: they are written into a target at S12 and pushed by the Operator. Off the lens, one observation: the two actions are pinned by tag, not by commit SHA. `holds`.

10. **(REQ-017, URLs and stand-ins).** The URL constants are `pinSource = "https://github.com/pharzam/armature"` (`internal/cli/run.go:30`), recorded in the run row (`:139`) and read by `LsRemote` and `Clone` (`internal/setup/setup.go:231-232`, `internal/run/steps.go:365-370`); `https://github.invalid/…` in the baseline's stand-in (`internal/standin/standin.go:252`), which starts nothing; and link text in `internal/run/betbrief.go` and the setup's command text. `internal/forge/fake.go` is an in-memory forge of the standard library that production code never constructs. `note`. No change: reads of a public baseline, and test doubles.

11. **(REQ-017, the lists of the release check).** Each list of `.review-in/release-check.txt` was reproduced and judged: check (1) passes (clean tree at `01d7da30…`); the six `exec` sites of (3), the six verb hits of (4), the embed line and the five kind commands of (5), the 21 `.do(` rows of (6) and the 17 callers of (7) (18 hits less the comment line `internal/session/session.go:39`) all equal the expected output, and every entry is accounted for above. `holds`.

### Runs

| Command | Result |
| ------- | ------ |
| `git rev-parse HEAD` | `01d7da3082333efadb227afda14fd98715edf5c4`; exit 0 (run at the start and again at the end). |
| `git status --porcelain --untracked-files=normal` | no output, exit 0 — the tree is clean, check (1) passes; run again after every other command: still no output. |
| `sh .review-in/release-check.sh 01d7da3082333efadb227afda14fd98715edf5c4` | **Not run**: the harness's shell allowlist (`opencode.json`) denied it (`Permission denied: shell`), as it denied `go list`, `go version`, `go env`, `awk`, `nl`, any chained command and any `VAR=… cmd` prefix. Each of its checks (1)–(7) was instead reproduced by hand with the allowed commands below, and every list was compared with `.review-in/release-check.txt`: equal in every entry. |
| `go build ./...` | no output; exit 0. |
| `go vet ./...` | no output; exit 0. |
| `go test -count=1 -tags=integration -run '^TestPackageRules$' ./cmd/layup/` | `ok  	github.com/pharzam/layup/cmd/layup	0.297s`; exit 0 (check (2)). |
| `git ls-files 'cmd/*.go' 'internal/*.go' ':!*_test.go' ':!*/testdata/*'`; `git ls-files internal/catalog/go` | 140 and 8 paths, which with `go.mod` are the 149 of `.review-in/release-files.txt`, in its set; exit 0. |
| `wc -l` over the 149 paths | every count equals `release-files.txt` and the Coverage table; 38,708 lines in all; exit 0. |
| `git grep -n` for `exec\.`, `os/exec`, `os.StartProcess` over `cmd/` and `internal/` | the six sites of check (3) in non-test files; every other hit is a `_test.go`, a `testdata/` file or the comment at `gates.sh.tmpl:218-220`; exit 0. |
| `git grep -n` for `"push"`, `"fetch"`, `"remote"`, `"pull"`, `"send-pack"`, `ls-remote`, `clone` in `internal/git/git.go` | the six hits of check (4) plus `ls-remote` at `:135` and `:461` and `clone` at `:148` (reads, left to check (7) by design); exit 0. |
| `git grep -n` for `go:embed` and the rows of `internal/catalog/go/kinds.tsv`; for `discTests`/`linkLint` in `internal/verify/scripts.go` | `internal/catalog/embed.go:12`, the five commands and the two scripts of check (5); exit 0. |
| `git grep -n '\.do('` in `internal/forge/github/*.go` (non-test) | 21 occurrences, each matched to an operation and a method as check (6) lists; exit 0. |
| `git grep -n` for `git\.(Push|Fetch|FetchSession|LsRemote|RemoteHead|Clone|CloneLocal|AddRemote)` over `cmd/` and `internal/` out of `internal/git`, and for aliased imports of `internal/git` | the 17 callers of check (7) after the comment line `internal/session/session.go:39`; no aliased import; exit 0. |
| `git grep -n` for `"unsafe"`, `syscall\.`, `"plugin"`, `import "C"`, `go:cgo`, `go:linkname`, `net/http`, `https?://`, and the vocabulary of training, weights, fine-tuning and model providers, over the non-test files | the hits of findings 1, 2, 8 and 10; `net/http` only in `internal/forge/github/github.go`; exit 0, or exit 1 where there was no match (`plugin`, cgo, `linkname`, the training vocabulary). |
| `git grep -n` for the citations of the findings (`GIT_CONFIG_`, `GIT_ALLOW_PROTOCOL`, `CheckRedirect`, `Proxy`, `cmd.Env`, `GH_TOKEN`, `gh api`, `rulesets`, `layup-probe`, `pinSource`, `WaitingSpecAccepted`, `workflows/`, `CloneLocal`, `commands.sh`) | the line numbers cited above; exit 0 (one search, `RunPRDLint` in `internal/run`, exit 1: it is called through `spec.WaitingSpecAccepted` at `internal/run/shapeloop.go:816`). |
| `ls /home/layup/go/pkg/mod/golang.org` | `toolchain@v0.0.1-go1.26.10.linux-amd64`, `x`; exit 0 — see the caveat below. |

**What I could not do.** (a) The release check script itself could not be run, so its `== PASS` line is mine by reproduction, not the script's. (b) `go list -deps ./...` could not be run; the dependency list of check (2) rests on `go.mod` having no `require` (`go.mod:1-3`) together with the green build, and on `git grep` of the imports. (c) `.gocache/` and `.tmp/` of the clone could not be used, because an environment prefix (`GOCACHE=…`, `TMPDIR=…`) is not in the allowlist; the build, vet and test ran with the default cache outside the clone, and wrote nothing inside it (`git status` stays empty). (d) `go version` and `go env` were denied, so I cannot state the version of the `go` on `PATH`; `go.mod:3` asks for `go 1.26.10`, the build and vet exited 0 with no failing line, and a toolchain `go1.26.10` sits in the module cache (`ls` above) — whether it was already there before my build, or `GOTOOLCHAIN=auto` fetched it during the build (which would be one network use that I did not intend and could not observe), I could not determine. (e) The authorship of the release's commits could not be checked, since git history is forbidden.
