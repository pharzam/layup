# Guardrails — the known pitfalls, and the numbers you must not move

This document is the target of gate step 2, "Honor the guardrails", in
[`engineering-discipline.md`](engineering-discipline.md). It holds the pitfalls a
task author must take into account **before** writing code, so the team does not
re-derive a known trap every time.

It holds two kinds of guardrail in one file — **decision gates** (pre-registered
pass/fail rules) and **validation** (how you check you are not fooling yourself).
Both are read before work starts.

## In plain terms

The worst mistake this project can make is to fill a setup value with a guess
and then report the setup as done; the defense is a recorded source for each
value, an open question where no source exists, and a script that fails on each
gap with no open question (`F-0001#4`, `F-0001#5`).

## 1. Pre-registered decisions — or the goalposts move

A decision rule chosen **after** seeing the result is a fitted parameter, not a
rule. Write the pass/fail numbers first, somewhere they cannot be quietly edited.

- **What must be pre-registered:** the Layer 1 checks of PSB §7.1 (pass or fail, from the invariants), each
  Layer 2 target of PSB §7.2 before a pilot measures it, and each task's `Budget
  maximum` and `Cycle cap`.
- **Where the numbers freeze:** the PSB itself (fact `F-0001`, hashed in
  `docs/setup/facts.sha256`), and the plan-review comment on each issue.
- **The bands, not a single line:** prefer **Pass / Investigate / Fail** to a
  single pass line on a noisy measure. Define "Investigate" with a rule written
  before you look — for example, one re-examination whose scope is fixed in
  advance; landing there twice counts as Fail.

### 1.1 LAYUP System Invariants (PSB §6)

These nine rules bind every solution in this repository. They come from the
approved PSB and are frozen with it: a change needs a new PSB revision, not an
edit here. Each entry gives the rule with its fact ID, the trap that breaks it
in silence, and the check that catches a violation. A `Check:` value is a path
plus the gate that runs it (`hook` or `ci:<job>`), or the words `no check yet`.
A script that exists but that no gate runs is `no check yet` (Invariant 5).
Whether a named gate really runs the path is a review judgement.

- **Inv-1** — Git is the system of record (`F-0001#1`). Trap: a decision made in
  a chat, a dashboard, or a tool database, and never written to the repository.
  Check: no check yet
- **Inv-2** — The project repository is independent (`F-0001#2`). Trap: a gate
  that passes only while the automation that made the repository is present.
  Check: no check yet
- **Inv-3** — The agents that do the work cannot change the rules or the gates
  that check the work (`F-0001#3`). Trap: a change that edits the check that
  judges it. The CI restore step (`guardrails.md` §2, "A check the change
  supplies is not a control") covers each check script a workflow runs, except
  `docs/tests/nested-checkout-check.sh`. Check `ci` (cause `restore`) fails when a
  job names a `docs/...sh` script outside its Restore step, in any form, and the
  Restore step does not name it (a script run through a relative
  `working-directory:` is not seen). Neither covers the workflow file itself, which a branch can
  edit, nor the Go jobs `lint`, `tests` and `security`, which run the pull
  request's own Go code and tests; so this invariant has no check yet (O-9,
  ADR-0011; for a target, ADR-0017 designs a control, which is not built). Check: no check yet
- **Inv-4** — No configuration value without evidence (`F-0001#4`). Trap: an
  Armature example value accepted as a project value. Check: no check yet
- **Inv-5** — A check that is not active does not count as passed (`F-0001#5`).
  Trap: a script in the tree that no hook or CI job runs. Check: no check yet
- **Inv-6** — A deterministic check is preferred to an LLM judgement where a rule
  can be checked mechanically (`F-0001#6`). Trap: a review round asked to settle
  a claim that a script could settle. Check: no check yet
- **Inv-7** — The project domain changes content, never rules; the technology
  stack can add stack-dependent gates but cannot remove or weaken a baseline
  rule (`F-0001#7`).
  Trap: a baseline rule edited during setup. Check: no check yet
- **Inv-8** — Armature is used at a pinned, recorded version (`F-0001#8`). Trap:
  a copy with no record of the commit it came from. Check: docs/setup/setup-check.sh (ci:setup-check)
- **Inv-9** — A harness agent is replaceable (`F-0001#9`). Trap: rules kept only
  in one agent product's own file format. Check: no check yet

## 2. Known pitfalls — the traps specific to this domain

The traps below hurt this project. Each has the trap, why it is silent, and the
check that catches it.

- ❌ **A shell function overwrites its caller's variable.** POSIX `sh` has no
  local variables, so a loop variable in a check function can replace the loop
  variable of the main loop in `setup-check.sh`. The first try of a fix in
  `T-r7zg` used `c` and `k`; the fixture run went red before the commit. It is
  silent when the names happen to agree. **The check:**
  prefix each variable of a check function with the check name (`pin_`, `fa_`),
  and keep a fixture that runs two or more checks in one call
  (`frame/good-passthrough`). Learned in `T-r7zg`.
- ❌ **`while read` drops a last line that has no final newline.** A hash list
  whose last entry had no newline skipped that entry, and a changed file passed.
  It is silent because the loop ends normally. **The check:** write
  `while read -r a b || [ -n "$a" ]`, and keep a fixture whose last line has no
  newline (`facts/bad-hash`). Learned in `T-fvwj`.
- ❌ **A merge while checks are pending.** A pull request was merged while
  `gh pr checks` showed 3 of 8 CI jobs pending; all 8 passed, but a red would have
  merged the same way. It is silent because the merge command does not wait. **The check:**
  read `gh pr checks` until no job is pending before the merge, and make the jobs
  required on `main` (`T-afa5`, [#12](https://github.com/pharzam/layup/issues/12)).
  Learned in `T-xgz4`.
- ❌ **The hook refuses a red commit.** Test first (R8) needs a red run before
  the code, and a task often commits its tests as their own step; but the
  `pre-commit` hook runs `go vet` and the unit level, so a test-only commit does
  not compile or does not pass and is refused. It is silent
  in the history, because the red step then leaves no commit. **The check:** run
  the red step and paste its failing output into the task's decision note (R7),
  then commit the tests with the code; never bypass the hook (`--no-verify`).
  Learned in `T-dq05`.

- ❌ **A rule for the builders looks like a product task.** A routing policy for
  the agents that build LAYUP (which model reviews which comment) took a day, 9
  plan reviews, 11 review rounds and 26 of the project's 40 Operator decision
  numbers (19 decided, 7 open), and changed no product line; no PSB fact asked
  for it. It is silent because every step passed the gate: each issue had a plan
  and a plan review, and most rounds a record, so the process looked healthy
  while the product stood still. **The check:** a plan names the PSB
  In-Scope fact (`F-0003#41`–`#52`) it serves, or the task does not start
  ([the scope of a task](engineering-discipline.md#working-a-task-under-the-quality-gate),
  first [ADR-0012](adr/0012-build-layup-in-bootstrap-mode.md), part 1); and the
  product-to-process ratio of [`tasks/completed.md`](tasks/completed.md) is read
  at each pilot. Learned in `T-8ywj`.
- ❌ **A copied template keeps its voice.** After the setup, the documents still
  spoke as the Armature template to the person who copies it, left CI as a
  choice that was already made, and `AGENTS.md` said that the project has no
  product toolchain while the Go code built. It is
  silent because every setup check read values, markers and links, and none read
  the voice; a reader or an agent then acts on a choice that was already made.
  **The check:** check `adapted` of [`setup/setup-check.sh`](setup/setup-check.sh)
  in the `pre-commit` hook and in CI (see
  [Check `adapted`](setup/README.md#check-adapted)); a false positive is fixed by
  a better sentence, not by a list entry. Learned in `T-745n`
  ([#70](https://github.com/pharzam/layup/issues/70)).
- ❌ **A public-solution search that confirms the design.** The first search of
  `T-hbw8` ran after the solution shape was fixed, used the words of that
  solution, and passed the one list that named the missed competitor (Paperclip,
  line 163 of about 247) through a summarizer with a filter and a cap of 15. It is
  silent because the selection record shows a search section with sources, so
  the rule "search first" looks met. **The check:** the search runs before the
  decision that fixes the solution shape; queries come from the problem
  statement's own terms; each curated list is read in full, every relevant entry
  with a keep or reject reason; two searchers on different models, blind to each
  other ([`runs/T-hbw8/root-cause-missed-solution.md`](../runs/T-hbw8/root-cause-missed-solution.md)).
  Learned in `T-hbw8` ([#72](https://github.com/pharzam/layup/issues/72)).
- ❌ **A citation check that proves the path, not the claim.** In `T-w73g`,
  `cites.sh` was green on each of 66 source paths of the survey, and the review
  then found seven rows whose cited file said something else: a count, an outcome
  left out, a script said to "only warn" that exits 1, a second file that held
  the claim. The rows came from read-only agents' reports, which round a detail
  or name the nearest file. It is silent because a green path check reads as "the
  survey is true to its sources". **The check:** a path check proves only that a
  file exists; before an agent's report of a file becomes a record, a reviewer
  reads a sample of its rows against the files (a semantic pass of at least 15
  rows, as round 1 of `T-w73g` did), and the record says which claims were
  checked by hand. Learned in `T-w73g` ([#133](https://github.com/pharzam/layup/issues/133)).
- ❌ **Coverage by name.** A coverage table that maps each requirement to a
  component name looks complete while the component cannot do the work: the
  `T-hbw8` architecture passed its own tables and then got about 50 material
  findings when two reviewers walked concrete cases. **The check:** each coverage
  row points to a walkthrough of one concrete case, each step tagged `code`,
  `model` or `human`, and an independent review of the walkthroughs runs before a
  human is asked to approve. Learned in `T-hbw8`.
- ❌ **A tool under test writes into the host's home directory.** In the `T-hbw8`
  evaluation, three of eighteen agent tools wrote outside their work directory
  on their first calls: one installed a Claude Code plugin and turned it on in
  `~/.claude/settings.json`, one wrote 628 command and skill files for 13
  harnesses into `$HOME`, one created `~/.gt` and `~/.dolt`; a fourth tried to set
  branch protection on GitHub. A brief that said "do not change the global
  configuration" did not stop them. It is silent because the tool's install step
  runs by default and reports success; every later harness session on the host
  then loads the added hooks and skills. **The check:** start each tool and each
  harness it drives with `HOME` set to a directory under the work directory from
  its first call, and with no `gh` login; after the runs, list the files in the
  real home directory that changed in the run window
  ([`runs/T-hbw8/evaluation/summary.md`](../runs/T-hbw8/evaluation/summary.md),
  "Incidents during the evaluation"). Learned in `T-hbw8`.
- ❌ **A tool under test reads the host's home configuration.** The mirror of
  the trap above. The plan of `T-2tc2` (#79) isolated `git` with
  `GIT_CONFIG_NOSYSTEM=1` and an empty `GIT_CONFIG_GLOBAL`; its plan review
  measured five inputs that still reached `git` (the per-user attributes and
  ignore files, `GIT_CONFIG_COUNT`, `GIT_DIR`, `GIT_AUTHOR_NAME`), and the task
  found a sixth (the system attributes file). Its verification found three
  paths to a credential of the host, each through a program under `git`, that
  a fixed list with the host's `HOME` left open: libcurl read `$HOME/.netrc`;
  `ssh` reads the user's keys in the home of the password database, not of
  `HOME`; and a URL can start a remote helper of the `PATH`. It is silent
  because a test host has none of them, so each test passes while another host
  gets another tree or sends a credential. **The check:** give the tool an
  environment from a fixed list, with a home of its own that holds no file:
  never the host's `HOME`, and never no `HOME`, as a library then reads the
  home of the password database. Keep each program that finds the user's home
  without `HOME` from starting. Close each file with a documented setting, or
  with a measured one that the record names. Keep a test that seeds each
  input, shows that it changes a plain run, and proves that it changes nothing
  through the code (`TestAHostileHostChangesNothing` and
  `TestNoCallUsesACredentialOfTheHost` in
  [`internal/git`](../internal/git/git_integration_test.go)). Learned in `T-2tc2`.
- ❌ **A model on the "not used" list, used inside a test run.** The `T-hbw8`
  evaluation brief named Claude Haiku 4.5 as the cheap model inside the
  candidates' test runs, and eleven candidates ran with it. ADR-0012 part 3 lists
  Haiku as not used and has no exception for a test fixture; the plan review
  found it, and it became a reported deviation. It is silent because nothing
  checks a model name in a brief, and a short alias (`sonnet`, `opus`) hides which
  version ran. **The check:** before a brief names a model, compare the full model
  ID with the models not to use in
  [Model tiers](engineering-discipline.md#model-tiers) (first ADR-0012 part 3), and write full model IDs
  in every evidence file. Learned in `T-hbw8`.
- ❌ **A word-for-word copy with trailing spaces.** A comment or a tool output
  copied into a record keeps its trailing spaces. No hook runs `git diff --check`,
  so 17 such lines passed every commit of `T-hbw8` and failed only at close-out.
  **The check:** after you copy text into a record, run `git diff --check` on it;
  remove trailing spaces, which carry no meaning there. Learned in `T-hbw8`.
- ❌ **A mutation undone with `git checkout` undoes the work too.** To show that a
  check can fail, the author added a line to a file and then ran
  `git checkout <file>` to remove it. The file also held the task's own edits, not
  yet staged, and they were lost with the test line. It is silent because the
  checks pass on the old text. **The check:** commit or stage the work before a
  mutation, or run the mutation on a copy of the tree. Learned in `T-gd8q`.
- ❌ **A headless reviewer whose commands are denied.** Two review rounds of
  Claude Fable 5.1 (`claude -p --permission-mode dontAsk`) gave no record in
  fifteen minutes. The cause was the run, not the model: the RTK hook rewrote
  `git diff`, `ls`, `wc` and `diff` to `rtk …`, which the allow rules did not
  match, so they were denied; and text output is written only at the end, so the
  timeout lost the work. **The check:** allow `Bash(rtk git:*)`, `Bash(rtk ls:*)`,
  `Bash(rtk wc:*)` and `Bash(rtk diff:*)` (never `Bash(rtk:*)`: `rtk proxy` runs
  any command), use `--output-format stream-json --verbose`, and read
  `permission_denials` of the result. Learned in `T-gd8q`.
- ❌ **A variable name that reads as a secret to `gitleaks`.** A script of
  `runs/T-zwke/` held `env = dict(os.environ, GH_TOKEN=token, …)`; no secret was in
  it, but the rule `generic-api-key` matched it, and the CI job `security` failed.
  `gitleaks git` scans each commit, so a later edit does not clear a finding, and
  a published branch is never rewritten. The way out is one line of
  `.gitleaksignore` with the finding's fingerprint, a change to the input of a
  gate (O-166). **The check:** before the first push, run `go run
  github.com/zricethezav/gitleaks/v8@v8.30.1 git --redact --log-opts="origin/main..HEAD"`
  as a command of its own and push only on its exit 0, never through a pipe (a
  pipe into `tail` let a finding be pushed, O-191); write a credential into an
  environment as `env["NAME"] = value`, and name a test value `want`, not
  `token`. Learned in `T-zwke` and `T-nxe4`.
- ❌ **A one-line form that loses a failure.** The first lint command of the first
  pilot put the gate script of the brief into one pipe. When `git ls-files` failed
  (an invalid index), the pipe took the status of its last command, and the check
  passed with no file read. It is silent because the form gives the script's
  output in each case that works. **The check:** before a short form replaces a
  script, run both on a failure of each step, and keep the cases (13 cases in
  `T-evad`: F-9 of [`runs/T-evad/findings.md`](../runs/T-evad/findings.md)). Learned
  in `T-evad`.
- ❌ **A record read by the shape of a host file of now.** `start.tsv` holds a
  `harness.<id>.cap` and `.wall` row for each harness of the register at Start,
  and the restart read it against the register now, so a target started with no
  harness was refused at `clone` once the register had one (#199). It is silent
  until the restart, as Start and the register's readers pass. **The check:** a
  restart reads its records by what they hold, never by a host file that can
  change after they were written; a test restarts with that file changed.
  Learned in `T-x7cs`, fixed in `T-s7vr`.
- ❌ **A mutant that does not build, counted as caught.** A mutation run counts a
  mutant as caught when `go test` fails; a mutant that replaces a condition with
  `false` can leave a variable unused, so the package does not build and the
  test "fails" with no rule read. It is silent because "FAIL" reads the same.
  **The check:** a mutant keeps the code building (join the condition with
  `&& false`), and the run counts a "build failed" as not run, never as caught.
  Learned in `T-529i`.
- ❌ **A file list that is empty, read as a pass.** The first draft of check (4) of
  `release-check.sh` listed its files with `git ls-files 'internal/git/*.go'
  ':!*_test.go'`; with one positive pathspec and that exclusion, git 2.47 lists
  no file, so the check read nothing and printed `none`. It is silent because
  "no hit" and "no file read" give the same output. **The check:** a check that
  reads a list fails when the list is empty, and its first run is read for the
  files it names. Learned in `T-4tjy`.
- ❌ **A claim of a check that nothing runs.** In the first pilot, a setup answer
  filled a command into a comment of the hook, and the prose called the gate jobs
  required while the ruleset was only prepared. A reader takes both for checks
  that run. It is silent because each value is correct and in its place. **The
  check:** for each "X runs" or "X blocks a merge", name the line that runs X or
  the read-back that shows it, and test the claims (F-5, F-10 and F-12 of
  [`runs/T-evad/findings.md`](../runs/T-evad/findings.md)). Learned in `T-evad`.
- ❌ **A command file that goes on after a failure.** The `commands.sh` of S13 has
  no stop on an error: after a refused push it pushed `layup-records`, applied the
  ruleset and exited 0. It is silent because the last command passes. **The
  check:** run the commands one at a time, check each before the next, and read
  back the pushes before the ruleset
  ([`runs/T-evad/t3.sh`](../runs/T-evad/t3.sh), finding F-13). Learned in
  `T-evad`.
- ❌ **A value audit that reads only the values it is given.** In the first pilot,
  two passes of a model audit and the author's own check found each value at its
  place. They missed a marker that S11 filled in part, and the markers that the prose
  step removed (F-28). In the second run, an inventory that matched the places of
  two runs by their values put 15 new places on old lines, and it left out the values
  that no record row names (F-29). It is silent because each listed value is right.
  **The check:** make the inventory from the diff of the two trees, so that each added
  line has a row or only moved; check that each marker of the baseline is kept or has
  a record row; and let the auditor check the completeness against the diff
  ([`runs/T-evad/findings.md`](../runs/T-evad/findings.md)). Learned in `T-evad`.
- ❌ **A claim of a later state in the present tense.** The task file of the first
  pilot's target task said that a commit "copies the tree into `main`", and that the
  runs after the merge "are in" a record, while the branch was not merged yet. A
  reader takes such a claim for done work, and can skip the step. It is silent
  because the claim becomes true after the merge. **The check:** before the freeze,
  read each claim about the merge or a later run, and give it the future tense or its
  condition ("the merge will take it into `main`"); the line of the completed log
  says what the task delivers (round 4 of `T-vu2j`,
  [`pharzam/chat-orchestrator#6`](https://github.com/pharzam/chat-orchestrator/issues/6)).
  Learned in `T-evad`.
- ❌ **Every GitHub response carries `x-ratelimit-reset`.** A `403` for a missing
  permission carries a reset time too, so a reader that takes "a `403` with a
  reset time" for a rate limit waits up to an hour, then fails the same way. It
  is silent because the wait prints progress lines like a real limit. **The
  check:** a rate limit is `retry-after`, or `x-ratelimit-reset` when
  `x-ratelimit-remaining` is `0` ([`forge.md`](spec/forge.md#forge-errors)); keep
  the case of a `403` with requests left, which must fail with no wait
  (`internal/forge/github`). Learned in `T-6bq5`.
- ❌ **`git fetch` into the branch of a new repository.** `git init -b main`, then
  `git fetch -- URL refs/heads/main:refs/heads/main`, exits 128 with
  `fatal: refusing to fetch into branch 'refs/heads/main' checked out`, because
  `main` is the branch of `HEAD` even with no commit. It is silent in a test that
  fetches any other branch. **The check:** `Fetch` passes `--update-head-ok`
  ([`packages.md`](spec/packages.md#the-calls-of-internalgit)), and its
  integration test fetches `main` into a repository that `Init` made. Learned in
  `T-xhgz`.
- ❌ **A plan row sliced by package, with no goal count.** The slicing of `M2a`
  (`T-zwke`, #124) cut one row per package and did not count the goal classes of
  each row, so the count of R11 came only at the build plan of rows 22 and 25,
  and each needed a split by the Operator (O-170 of #127, O-173 of #130). It is
  silent because a row per package reads as one task, and the build plan comes
  after the issue is open. **The check:** before the issue of a row opens, apply
  the test of R11 ("either can fail while the other passes") to the goal classes
  of each row, and bring a row with more than one to the Operator then; a split
  that a slicing review offers is weighed by the goal count, not by a shared
  test. Learned in `T-trej`.
- ❌ **A question whose summary its own table contradicts.** The plan of
  `T-fdaq` (#152) counted the goal classes of each row in a table, then
  proposed a split "which gives 16 rows of one or two classes each"; by that
  table, three rows kept three classes. The Operator answered on the summary
  (O-186), and the plan review found it, so the question was asked again
  (O-187). It is silent because the table is right and only the sentence over
  it is wrong. **The check:** before a question goes to the Operator, derive
  each number of its summary from the table it summarises, row by row, and put
  the per-row result in the question. Learned in `T-fdaq`.
- ❌ **A hostile configuration whose keys hide each other.** The test of
  `FetchSession` (`T-z5dj`) writes each program-starting key of `git-config(1)`
  into a clone, and a control shows which fire on the host. Its first runs
  showed four ways a key is silent though the list holds it: `include.path`
  set a key that the file also set, so one of the two never ran; a process
  filter that fails aborts `git diff` and `git add` before the other drivers
  run; `git diff` stops at the first external diff that fails, and `textconv`
  is read only by the built-in diff; and a marker program that reads its
  standard input deadlocks a process filter and `upload-pack` (a five-minute
  hang). And a list written from memory left out about twenty keys that
  `git help --config` lists. **The check:** build the list from `git help
  --config` of the host's `git`; give each driver key a driver and an attribute
  of its own; let no include set a key the file sets; read one path per `git
  diff`, with `--no-ext-diff` for `textconv`; keep each marker program from
  reading its input; and make the control name each key that did not fire.
  Learned in `T-z5dj`.
- ❌ **A harness's own tool runs a model not to use.** In `T-ywk7`, searcher B
  ran on Claude Opus 5.5 (`claude -p`), and the `modelUsage` of its result event
  also named `claude-haiku-5-5`, with 4,759,374 input tokens (USD 1.25): Claude
  Code's `WebFetch` reads each page with that model, which is on the list of
  models not to use ([Model tiers](engineering-discipline.md#model-tiers)). No
  message of the session names it, and the top-level `usage` counts the main
  loop only (53,309 output tokens against 911,458 in `modelUsage`). It is
  silent because the session ends with success. **The check:** give each
  `claude -p` run `ANTHROPIC_DEFAULT_HAIKU_MODEL` set to an allowed model (a test
  run with it used no Haiku), read `modelUsage` after each run, name every model
  of it in the record, and take the tokens of a resource record from it. Learned
  in `T-ywk7` ([#147](https://github.com/pharzam/layup/issues/147)).
  A resumed Claude session also accumulates `modelUsage` across invocations:
  take per-model deltas between result events, preserving every model, rather
  than assigning the cumulative total to the latest part. Top-level `usage`
  is sufficient only after verifying that no child model ran. The plan-plus-build
  counters were mistaken for build-only counts in `T-cak3` (#357).
  Devin's console can omit usage while --export records ATIF final_metrics.
  Compare completed exports of a resumed conversation to attribute each part.
  ATIF total_cached_tokens is a subset of total_prompt_tokens, not an extra
  input count ([format](https://github.com/harbor-framework/harbor/blob/main/rfcs/0001-trajectory-format.md#metricsschema)).
  Keep private session exports outside version control and commit only reported
  metrics, preserving unavailable parts as unavailable. Learned in #359.
  A resumed export can lose earlier agent metrics after compaction while keeping
  the original system timestamp. Compare per-step timestamps and metric rows,
  not just session IDs or the final totals; union matching rows once, and label
  missing earlier usage unavailable. In #351 the interrupted export was a subset
  of the final export, while the plan metrics had been discarded.
- ❌ **A Go version that CI resolves at its own time.** `go.mod` said `go 1.26`,
  and `actions/setup-go` with `go-version-file: go.mod` installs the newest
  release of that line it knows. The job `security` passed on `main` at
  `ad3c400` and failed on #149, which changed no Go code: ten advisories of the
  standard library, found in `go1.26.8` and fixed in `go1.26.9`, were published
  between the two runs. It is silent until a release day, as each run of the
  same tree can resolve another toolchain. **The check:** the `go` directive
  names the full patch version (`go 1.26.9`), so CI builds what the commit
  names, and govulncheck's "Fixed in" names the version to move to; the red returns at each security release of
  the line, and the directive moves with it. A host
  below it then downloads that toolchain under `GOTOOLCHAIN=auto`, or refuses
  to build under `GOTOOLCHAIN=local`. Learned in `T-w89c`
  ([#150](https://github.com/pharzam/layup/issues/150)).
- ❌ **A review record without its heading.** `review-record-lint` takes a comment
  for a record only by a line that starts `## Review record` and names the round
  after the word `round`. A final message with no such line is no record to it:
  round 2 of `T-bng8` had none, so the linter read the rounds of the issue as 1, 3,
  4, 5 and 6, and the required check `review-record` was red on each push of the
  pull request until the heading was added. It is silent on the issue, where the rounds run and
  nothing shows that check: the record reads well, and no step that posts it looks
  at the heading. The check of the pull request runs on a pull-request event, so it
  shows the gap at the next push, not when the record is posted. **The check:** the
  brief of a round says that the first line of the final message is the heading
  `## Review record — round N`; and after each post, build the stream as the
  workflow does and run the linter on it:
  `gh api --paginate "repos/OWNER/NAME/issues/N/comments" --jq '.[] | "=== comment created=\(.created_at) updated=\(.updated_at) ===\n\(.body)"' > comments.txt`,
  then `sh docs/ci/review-record-lint.sh comments.txt`. Learned in `T-bng8`
  ([#267](https://github.com/pharzam/layup/issues/267)).
- ❌ **An allow rule whose glob reaches a script with a write mode.** A reviewer
  session runs read-only by its allow rules (`claude -p --permission-mode dontAsk`).
  The allow rules of a harness script, outside the repository, as they stood after
  a round with many denied commands, included `Bash(python3 runs/*)`, `Bash(PLAN=*)`
  and `Bash(sh runs/*)`, which let a reviewer run the check of a task with an option.
  The same globs also allow `gen-issues.py --post` and `--edit`, which open and edit
  issues through the App: writes to the forge that the brief barred only in words.
  It is silent because a reviewer that follows its brief never runs them, and the
  rules read as read-only. **The check:** allow the exact read-only forms (`--table`,
  `--body ROW`, `sh runs/T-xxxx/parts.sh`), never the glob of a directory, and before
  the round run a probe against the same rules, with a negative control: a command
  that must be denied (`gen-issues.py --help`, which prints only its usage and is
  not in the allow list). The rules of "A headless reviewer whose commands are
  denied" have the opposite fault, too narrow to run the commands of the brief.
  Learned in `T-bng8` ([#267](https://github.com/pharzam/layup/issues/267)).
- ❌ **A Start that names the wrong idea owner.** `layup run --new … --idea-owner
  ACCOUNT` writes `approvers.tsv` once, and an answer counts only from an approver in
  the role of its block: the block of the idea owner takes the answers of that account
  and of no other. No step changes the file. In the demo of `M2c` the person who owned
  the idea answered the follow-up from another account than the Start named: the restart
  copied those comments as input and read no answer in them, and alone they would have
  left the step at `waiting: no answer yet to the follow-up in the block idea-owner`,
  which reads as "nobody has answered". A quote of an answer with a `>` before each line
  (the quote-reply of the forge) holds none either, as a line counts only when it starts
  `Q-NNN: `. It is silent because the wait is the normal state of Intake and the comment
  reads well on the issue. **The check:** before Start, ask who will answer the batch
  and name that account in `--idea-owner`; if another account answers, the account that
  the Start named posts the answers in a comment of its own, one `Q-NNN: answer` line
  per question (a comment that held only a quote is edited to those lines before the
  next restart copies it: an edit before the first copy counts). Learned in `T-ekh5`
  ([#228](https://github.com/pharzam/layup/issues/228)).
- ❌ **A key recomputed from a field that the field rule changed.** The id of a marker in
  `intake/markers.tsv` is the SHA-256 of its file, a tab and its text, over the raw
  text, and the record holds the text after the field rule, which makes a tab or a
  carriage return one space. A reader that recomputes the hash from the row refuses the
  row of a marker that holds one; S11 already refuses such a marker for a `gap`, for the
  same reason. It is silent because each marker of a fixture is a word on one line.
  **The check:** the reader checks the form of a key that hashes raw text, and the
  writer computes it from the source; a test reads the row of a marker that held a tab,
  under the id of its raw text. Learned in `T-tetz`
  ([#270](https://github.com/pharzam/layup/issues/270)).

- ❌ **A path rule that passes a bullet it should refuse.** The register's `path` rule
  (`rules.ValidPath`) refuses `..`, a leading `/` and `.git`, but not a carriage return,
  and the free text of a finding is any text. A `- ` line of a plan or a review record
  that ends in `\r` (a CRLF checkout) therefore passed the path rule and the finding
  grammar and read as a valid bullet whose path or text held a hidden `\r`. It is silent
  because a fixture is written with `\n`. **The check:** the comment grammar refuses a
  bullet that holds a `\r` itself, a marker line with a `\r` is no marker, and a test pairs
  each such line with the same line without the `\r`. Learned in `T-4pje`
  ([#396](https://github.com/pharzam/layup/issues/396)).

### Writing a lesson back

A trap caught once should not be re-derived by the next task, so a lesson does not
stay on the issue that learned it. When a task ends, its author asks whether the task
taught a trap the next reader could hit; if it did, the lesson is written **here** as a
new `❌` pitfall — the trap, why it is silent, and the check that catches it — in the
same pull request. Gate step 7 asks the question, so the rule is applied rather than
merely written (see [Keeping documentation current](engineering-discipline.md#keeping-documentation-current)).

This is the one **cross-task** reach the gate adds on purpose.
[R6](issue-workflow.md#r6--agent-to-agent-communication-through-the-issue) and
[R7](issue-workflow.md#r7--decision-transparency-on-every-action) already keep the
coordination and the reasoning on the issue, and
[Honesty and evidence](engineering-discipline.md#honesty-and-evidence) already reports a
failure as a failure — but each is scoped to *one* issue thread. A lesson on issue #N is
discoverable only by someone who reads #N; §2 is where it reaches issue #N+1.

**The filter — or §2 grows until nobody reads it.** Write back only a trap that would
**catch the next reader**: a silent failure mode, a check that looked green for the
wrong reason, a footgun in the discipline system or the domain. Do **not** write back a one-off with
no general lesson, a restatement of a rule that already lives elsewhere, or the
blow-by-blow of the task — those belong to the issue thread and the commit history.
Volume is the failure mode here, not absence: a pitfall list nobody finishes reading
guards nothing.

### Gate pitfalls

The gate is only as real as the thing that runs it. These traps let it report
success without having done its job.

- ❌ **An absolute `core.hooksPath`.** Worktrees **share** `.git/config`, so an
  absolute path binds every worktree to one checkout's hooks. A check added on a
  branch then does not run on that branch's own commits: the hook reports success
  having run something other than what the branch says it runs. It is silent
  because the hook still runs, still passes, and still uses the *right* files —
  only the *set of checks* comes from elsewhere. **The check:** install with a
  **relative** path, `git config core.hooksPath .githooks`, which git resolves per
  working tree. The `pre-commit` hook's own block 0 then refuses to run when the
  resolved hooks directory lies outside the tree being committed to, and
  [`.githooks/tests/provenance-check.sh`](../.githooks/tests/provenance-check.sh)
  proves it against real worktrees, reporting how many cases it ran rather than a
  count written down here to go stale.
  **A relative path escapes just as surely:** `../elsewhere/.githooks` is as
  foreign as any absolute one, so the check resolves the value instead of trusting
  that relative means local.
  **No path at all is the same trap:** with `core.hooksPath` unset git falls back
  to `.git/hooks`, which a linked worktree reaches through the shared *common* git
  directory. So block 0 judges the resolved directory whatever set it, and names
  the source it actually found — telling an operator to fix a setting they never
  set is its own dishonest report.
  **Bound on the damage:** CI invokes each check script directly and never through
  `core.hooksPath`, so this costs a local round trip, not a landed bug — a
  developer-experience gap, not an open gate.
- ❌ **A check that cannot fail.** A grep whose pattern also matches its own error
  message, a fixture harness that compares only exit codes, a coverage floor that
  counts zero as success. The check: for every assertion, make it fail on purpose
  once and read the reason — a green nobody attacked is not evidence.
- ❌ **A check that runs but does not block.** CI is green, the pull request
  merges, and nothing connects the two: no check is required on the default
  branch, so the green was a run result, not a merge control, and a red would
  have merged the same way. It is silent because the run result looks identical
  either way. The check:
  [make the checks required](ci/README.md#make-the-checks-required) on the
  default branch, and take the branch API read as the evidence — not the green run.

- ❌ **A check the change supplies is not a control.** CI checks out the pull
  request's own head and then runs the check from that checkout, so **the script
  that judges the change comes from the change**. Measured: a branch that replaces
  `docs/links/link-lint.sh` with `exit 0` passes that job — and replacing
  `docs/tests/run-discipline-tests.sh` as well turns **every** required job green
  over a dead link in the tree. Gutting a linter alone does not, because the
  fixture harness asserts exit codes and 19 `bad-*` cases stop failing; the runner
  is the single point.
  **The trap inside the remedy:** each script roots its scan at its own directory
  (`dirname $0`), so running the default branch's copy *where it sits* lints the
  wrong tree and reports OK. That was measured too, while building the fix.
  **The check:** every job in [`ci.yml`](../.github/workflows/ci.yml) restores the
  check scripts from the default branch **in place** before running them, so the
  branch's copy is never the judge.
  **Bound on the damage, and it is not nil:** a `pull_request` event runs the
  workflow as the branch has it, so a branch that edits `ci.yml` removes the
  restore step — closed only by review of `.github/**`, which wants a `CODEOWNERS`
  entry and a second human the forge knows about. Restoring also stops a bypass,
  not a merge: once a weakened check lands it *is* the default branch's copy. And
  a change that *improves* a check is judged by the older copy, so it lands in two
  steps. These checks are a control against forgetting, not against an operator
  who edits the check.
- ❌ **A reason that copies an error into a result table.** A row's reason that
  holds the text of an error of `git` (for example `fatal: '<path>' already
  exists`) puts the path of a temporary directory into the table, so two runs
  on the same input print other bytes and the repeat rule of `REQ-007` breaks;
  it is silent because each single run looks right. Task `T-5sgt` (#82) planned
  such reasons, and its plan review found it before the code. **The check:** a
  reason in a table is a fixed text that names the failed part (for example
  `scratch tree: add failed`); the error goes to standard error; and a test of
  the failure asserts that the reason holds no path.
- ❌ **A file list without `-z`.** `git ls-files` quotes a name with `"`, `\` or
  a control character in its plain list, even with `core.quotePath=false`, so a
  script that reads the list as paths skips that file in silence: check
  `markers` did so until #21. **The check:** list the files with
  `git ls-files -z` (in sh, `| tr '\0' '\n'`, which still splits a name with a
  line feed), and give the fixture a name that git quotes (`markers/bad-quoted-name`
  has a DEL character, which a Windows checkout accepts).
- ❌ **A path in `awk -v`.** `awk -v f=<value>` reads escape sequences in the
  value, so a name with `\` reaches the program changed: `check_markers`
  printed `docs/bx.md` for `docs/b\x.md` until round 1 of #87, and failed a
  correct setup. **The check:** hand a value that holds a path to `awk` through
  the environment (`ENVIRON`), and test a name with `\` on a scratch repository
  at test time (no fixture can hold it: a Windows checkout refuses the name).
- ❌ **A status lost in a command substitution.** `test -z "$(cmd)"` tests only
  the output of `cmd`: when `cmd` is not found, or fails, it writes nothing to
  standard output, and the test passes. The form that §6 of the architecture
  names for the static kind of Go passed so, with no `gofmt` on `PATH` and a
  badly formatted file (measured, task `T-c06a`, #91). **The check:** keep the
  status, `out=$(cmd) && test -z "$out"`, and test the case where the program
  is not found. **Known limit:** LAYUP's own job `lint` keeps the form
  `test -z "$(gofmt -l .)"`; on its runner `setup-go` installs `gofmt` with
  `go`, and a change of `ci.yml` is K28 (note 7 of the plan review of #91).
- ❌ **A failed check that reads as a pass.** A script that maps the failure of
  a check to its pass, as `scan || result clear`, passes when the check does
  not run to its end. The job script of the Go entry gave `clear`, with its
  command not run, in about 5 runs of 100 with `bash` 5.3 on macOS, where a
  subshell of the script crashed (a segmentation fault; its cause is not
  known), and in none with `dash` (measured, task `T-c06a`, #91). **The
  check:** give a check three answers (a match, no match, a failure) and make
  a failure `not-active`; check the status of each read; run a script's tests
  many times with each `sh` of the host.
- ❌ **A hashed file that git may convert.** A list of hashes, such as
  `docs/setup/facts.sha256`, holds the bytes of each file it names; on a
  checkout with `core.autocrlf=true` (git's default on Windows) git gives a
  text file a carriage return per line, so each hash fails and a script does
  not run (#48). It is silent on the host that wrote the files. **The check:**
  each path of a hash list, the list itself and each script of a check have a
  rule in `.gitattributes` (`-text` for a file whose bytes count, `text
  eol=lf` for a script or a list that a script reads); the case
  `facts/good-autocrlf` of `docs/setup/tests/run.sh` clones `HEAD` with
  `core.autocrlf=true` and runs check `facts`.
- ❌ **A trigger that runs a branch's copy of a write-token workflow.**
  `workflow_dispatch` and `create` each can name a branch, and the run
  then uses **that branch's own file** — so a task branch's edited copy of the
  workflow executes with `permissions: issues: write`. (`delete` names a branch
  too, but the ref is gone: the run uses the default branch's file — it is left
  out of `task-state.yml` only because a run on a deletion adds nothing the
  schedule does not give.) Checking out the default
  branch inside the file does not constrain what the file itself does. It is
  silent because the trigger reads as a convenience and the write scope looks
  guarded by the checkout. **The check:** a workflow holding a write token takes
  no trigger a branch's copy can fire — the allowed set of `task-state.yml` is
  `push` to `main`, `pull_request_target`, `issues`, `issue_comment` and
  `schedule`, and `docs/tasks/task-labels.sh --check-workflow` fails any other,
  any `push` bound that is not `branches` of `main` alone, and any event
  `types:` outside the named ones
  (the `bad-labels-workflow-*` fixtures). Learned in `T-dm8m` (#185, the review
  of the re-plan, condition 1).
- ❌ **An `awk` check that reads the locale.** `awk` reads characters in a UTF-8
  locale and bytes in the C locale, so one sh check can give other lines on
  another host. macOS `awk` 20200816 in a UTF-8 locale stops ("towc: multibyte
  conversion failure") when `check_adapted` meets a word after a character of
  more than one byte, and the check then passes that file (task `T-8ya0`, #88);
  it is silent because the run prints `adapted OK`. **The check:** run an sh
  check with `LC_ALL=C` and in the host's locale on a case that is not ASCII;
  the Go port states the byte rule of its check. The same `awk` in a UTF-8
  locale missed a marker at the end of a line in `check_markers`, which runs its
  `awk` with `LC_ALL=C` since round 1 of #87.

- ❌ **Tests written while the plan review runs.** R12 puts the reviewer's
  confirmation before building begins, and the first thing built is the test slice, so
  a red run dated before the confirmation breaks the order, even when the review asks
  for nothing that touches a test. In task `T-9v92` (#274) the author used the nine
  minutes of the plan review to write the tests and ran them red at 08:07Z, seven
  minutes before the confirmation at 08:14Z; the reviewer of round 1 read the two times
  in the record and gave `material`, which cost a second round. **The check:** spend the
  wait of the plan review on reading, and on nothing that the plan names as a step; start
  the red slice after the answer is posted, and compare the time of the first red run
  with the time of the confirmation before the freeze.
- ❌ **A `state:` label — or its color — set by hand is overwritten at the next
  run.** `task-labels.sh`, run by `task-state.yml` on a 30-minute schedule and on
  issue and pull-request events, recreates the eight `state:` labels with
  `gh label create --force`: a description or a color edited on the forge holds
  until the next run, then reverts to the script's. It is silent because the edit
  applies — the label shows the new color for up to half an hour — and the revert
  lands on no issue and in no log. Task `T-p481` (#309) began from eight hand-set
  colors found back at `ededed`. **The check:** change the color in the one
  `LABEL=COLOR` list of `apply()` in [`task-labels.sh`](tasks/task-labels.sh), and
  let the workflow write it; never edit a `state:` label on the forge.

- ❌ **A merge of `main` moves the base of the budget.** The plan review names a base
  SHA, and the budget counts the lines added plus removed on the whole branch diff
  against it. A branch on the forge takes `main` by a merge before its freeze, so the
  diff against the old base also holds the tasks that landed meanwhile: in task
  `T-jpkc` (#271) it was 1,508 lines over 23 files at the head that round 1 read,
  against a maximum of 850 over 14, while the task's own diff against the `main` that
  the branch had taken was 614 over 10, and round 1 gave `material` for the first
  figure; at the head of round 2, after four merges, it was 4,823 over 61 against 700
  over 11. **The check:** when you merge `main`, say it on the issue in one comment,
  before the freeze: the new base SHA, the figure against it and the figure against
  the old base, and the tasks that landed between (`git log --merges OLD..NEW`); and
  ask the Operator to rule that the bound is read against the new base, as the plan
  review names the old one and an overrun against the named base is a finding
  ([the budget](engineering-discipline.md#reviewing-until-findings-decay)); the task file names the base that the
  ruling gives.

- ❌ **The pinned `prd-lint` fails on two streams.** Its awk `print` goes to
  stdout and its `err()` — the filename refusal — goes to stderr, so a reader of
  stdout alone sees an accepted run while a refused name stands unread. It is
  silent because every line starts `FAIL  ` and exit 1 is the expected verdict
  of a waiting specification. **The check:** `RunPRDLint` keeps the streams
  apart and the acceptance judges every non-empty line of both — only
  `MoSCoW "" not in Must|Should|Could|Won't` and `has no phase value` are an
  empty cell, a `—` MoSCoW is not (internal/spec, task `T-4wyg`).
- ❌ **A `clear` spec-check table that drops its rows.** Before the first bet, a
  head that touches neither the PRD nor `docs/spec/` is `clear` — but the table
  still holds the ten rule rows, `ran` included; a writer that emits nothing, or
  a partial table, leaves the head with a spec-check record `ran` never judged.
  It is silent because "clear" reads as "nothing to write". **The check:**
  `CheckFull` returns the ten rows in every mode and the writer round-trip reads
  them back (internal/spec, task `T-4wyg`).
- ❌ **A short table row panics the check.** A parser that indexes `cells[1]`
  crashes on a truncated row — and the pinned producer never reads that table,
  so nothing upstream catches it first. It is silent because every fixture row
  is well-formed: the suite stays green while one malformed line takes the
  whole check down with it. **The check:** each cell index carries a bounds
  guard, a truncated row records nothing (so its requirement fails the rule,
  never the process), and a test feeds the shortest forms — one cell and a
  bare `|` (internal/spec, task `T-4wyg`, round-1 finding F1).
- ❌ **Equal bytes are not an equal file.** A check that compares what `git show
  REV:PATH` gives with the bytes it expects (the basic `layup/rules` compares a
  rendered record) passes a symbolic link whose target text is the expected text,
  and an executable copy of the file: the blob of a link is its target, and the
  mode is not in the bytes. It is silent because every fixture is a regular file.
  **The check:** ask the tree for the entry first (`git ls-tree`, as `LsTree`
  does) and read only mode `100644` of type blob; give the test the same bytes as
  a link, as an executable and as a submodule, made with `update-index
  --cacheinfo`, which keeps them from being rewritten by `git add`. Learned in
  `T-9et9` (#352).
- ❌ **A fixture that holds one register cannot show which register a check read.**
  A test of "the register is read at the base" over a head and a base that hold
  the same `rule-paths.tsv` passes whichever it reads. **The check:** let the head
  carry a forged register that lists nothing, the base the true one, and keep a
  control that judges the same change by the forged register and passes, so the
  case can fail; a mutant that reads the head must fail it. Learned in `T-9et9`
  (#352).

- ❌ **Reading an earlier review to copy its record format.** An otherwise fresh
  reviewer read #359's verdict during #352, despite its blind brief; that invocation
  was not a qualifying review (issue comment 6217959614). **The check:** supply the
  required field names from the discipline text, and audit actual tool inputs against
  the brief before recording context independence.

### Testing pitfalls

These traps are not domain-specific: they hurt every project's test suite.

- ❌ **Testing after the code.** A test written to fit code that already "works"
  tends to encode the code's bugs as expected behaviour. The check: write the test
  first and watch it fail for the right reason
  ([strict TDD](engineering-discipline.md#requirements-traceability)).
- ❌ **Tests that depend on external state.** A test that reads a shared database, a
  live network, the wall clock, or another test's leftovers passes or fails for
  reasons unrelated to the code. The check: isolate and control every dependency,
  with a fresh fixture per run — see
  [`tests/scaling-checklist.md`](tests/scaling-checklist.md).
- ❌ **A fake clock made anew at each call of a world's helper.** The test world of
  `internal/run` made its fast clock inside `config`, so each restart of a run got a
  clock that began at 12:00:00 again, and the run kept its ID. A restart takes the
  lease over after 15 minutes of its own clock without a beat, and writes the row
  with its time cut to the second, so two restarts could write the same row in the
  same second, and `git commit` ended with exit 1, as it does when nothing changed.
  The test failed in 5 of 70 runs of CI (#291) and in 1 of 200 runs on a host of two
  CPUs, and its message names no clock. A real clock does not go back. Task `T-0p5p`
  (#293) met it. **The check:** make the fake clock once for each world, not for each
  call of the helper that reads the world's config, and assert that the clock of a
  second config is not before what the first read.
- ❌ **Tests that pass for the wrong reason.** A test that asserts nothing, asserts
  the wrong thing, or never actually exercises the path reports a safety that is not
  there — worse than no test. The check: confirm the test fails when the behaviour
  is broken; the red step is the proof. A red run on a skeleton proves only that a
  test needs some code, not each rule of it: in task `T-8ya0` (#88) two cases
  passed with a wrong lower-case rule and a wrong word position, because the test
  text was shorter in lower case and the extra finding had the text of a real one.
  For each rule, break that rule alone (a mutation) and watch its case fail. And
  give each case a fixture that breaks no second rule at the same column: in task
  `T-ysph` (#155) an empty `command` on a row with a cap was refused by the rule
  of `{cap}` too, so dropping the rule of the empty value left the case green,
  and a lone position `0` was refused by the rule of a gap as well as by the rule
  it named; move the case to a fixture where only its rule can refuse it, or
  assert the line as well as the column. A fixture can also be too small to tell
  two rules apart: in task `T-z027` (#166) the payload of a refused head was
  the diff of the head, and a demo that changed one file could not tell it from
  the diff of that file alone; give the fixture a second change that the wrong
  rule would leave out. The case of the empty value is one of these: the form of
  a column (a line from 1, a session, a run ID) refuses `—` too, at the same column,
  so a case that names only the column passed with its rule removed in 5 of 37
  mutants in task `T-tetz` (#270); assert the reason as well. A refusal that a
  refactor moves into a body that two callers share is another: in task `T-ktgx`
  (#273) the refusal of a gap table of another form moved into the shared body
  of S04, no test of phase 1 pinned its kind, and the mutant that dropped it
  survived; give the old caller a case that names the refusal before the refactor.
  When a mutation record is cut to a few lines, say so, or it reads as fewer
  failing cases than ran.
- ❌ **A step added to a list of steps, and the tests of the list at a tag that the
  author did not run.** The restart's steps are listed by name in tests of three
  levels: the unit test of `restartNames`, the integration world's lists and progress
  counts, and the end-to-end tests of `cmd/layup` that read the table that the binary
  prints. The last two need a build tag, so `go test ./...` and the hook leave them
  unread. Row 53 (`T-4c5c`, #276) added `markers` between `completeness` and `batch`,
  and the end-to-end lists were red only at `go test -tags=e2e ./...`; the same step
  made the Operator's block ask a question, so each test that answers the batch had to
  answer as the Operator too, which the plan review found for the integration tests and
  not for the end-to-end ones. **The check:** when a step is added to the restart, or a
  question to the batch, `git grep -n` the name of its neighbour step and the wait text
  ("waiting: no answer yet") in every `_test.go`, whatever its tag, and run the three
  levels before the freeze, not at the run of CI.
- ❌ **A package's integration level that needs minutes, run on a loaded host with the
  default limit.** `go test` stops a package after 10 minutes when no `-timeout` is given and
  prints the goroutines that still ran, which reads like a crash. On the host of the sessions
  (2 CPUs, a load average of about 6 when several sessions run their suites) the integration
  level of `internal/run` took 486 s alone and was stopped at 600 s inside a run of every
  package, while CI ran all three levels in 3 min 22 s (`T-4c5c`, #276). **The check:** run
  the level with `-timeout 30m`; when a package stops at the limit, read `uptime`, run that
  package alone with `-v`, record both runs, and edit no test for the stop. A package that
  needs more than 10 minutes alone on an idle host breaks
  [the rule of the scaling checklist](tests/scaling-checklist.md#fast-enough-for-hook-and-ci),
  and is a task.
- ❌ **A `git` verb that reads its standard input when its argument is empty.**
  `git commit-tree TREE -p PARENT -m ''` finds an empty message and reads the message
  from its standard input, as it does with no `-m`; with a standard input that stays
  open, the call hangs (measured with `git` 2.56.0: a pipe that stays open and a 5 s
  limit gave exit 124). `internal/git` leaves `cmd.Stdin` nil, so `git` reads the null
  device and the commit has an empty message; round 1 of `T-j4sd` (#272) found that no
  test passed an empty message. **The check:** keep `cmd.Stdin` nil in `call`, and give
  each call that passes a free text to `git` a test with the empty text, at the unit
  level (the arguments) and at the integration level (the result, and a return).
- ❌ **Stale tests after a requirement changes.** When a requirement changes but its
  test does not, the suite now guards the old behaviour and blocks the new. The
  check: the [old-tests conflict rule](engineering-discipline.md#testing) — fix the
  code, update the requirement with a written reason, or retire the test; never
  weaken a passing old test.
- ❌ **Tests that slow down as the project grows.** A suite that creeps past the
  hook's patience gets skipped, and a skipped gate is no gate. The check: keep the
  cheap levels fast and cheap-first, push slow ones to CI, and bound each with
  `-timeout 10m` — see [`tests/scaling-checklist.md`](tests/scaling-checklist.md).
- ❌ **A specification that fixes one input state at a time.** A step that reads
  an input file a human writes (an answers file, a list of rows) meets rows in
  states that the happy path never makes: rows given early, copied from an
  earlier run, or for a question that no step asked. A review finds one such
  state per round, and a fix per state gives one more round per state: task
  `T-0drh` (#74) needed rounds 2 to 4 for the one answers file of `layup setup`.
  The check: before the review, write one rule for every row of an input that
  no step of the run asked for (refuse it, with its exit code, or ignore it), and
  say which step writes each accepted row.
- ❌ **A plan check that reads the form of a plan, not the edges of its sources.**
  A check of a task plan can prove that each requirement has a task, that each
  predecessor is a task, and that the predecessors make no cycle, and still pass
  over a wrong order: an item whose source names a predecessor in another task, or
  a defect that one task settles and an earlier task reads. Task `T-55n2` (#76) had
  five such defects under a green check; an adversarial self-check found them.
  **The check:** compare each edge of the source inventory, and each "settled by" and
  "read by" pair, with the transitive predecessors of the tasks, and list each edge
  that the plan drops on purpose with its reason.
- ❌ **A glob that drops a read error.** Go's `fs.Glob` ignores an I/O error, so
  a directory that is missing or cannot be read gives no match and no error, and
  a test over "every file of the directory" passes with nothing read. Task
  `T-18v6` (#78) found it in its block reader. **The check:** read the directory
  with a call that returns its error (`fs.ReadDir`), and fail on zero matches
  where the test needs at least one.
- ❌ **A cycle cap raised after a last-round verdict.** `review-record-lint` reads
  the cap from a `## Plan review` comment, and a round that ended as the last
  round under the cap of its time carries a last-round verdict. When the Operator
  raises the cap and a fix follows, that round becomes an intermediate round and
  the check fails (task `T-0drh`, #74, O-120). The check: record a raised cap at
  once, as a `## Plan review` comment with the new `Cycle cap` row, and give the
  round that the raise reopened the verdict `material`, with the reason.
- ❌ **A run of the binary that does not test what it names.** A Go program that
  writes to a broken pipe on its standard output is stopped by the signal
  `SIGPIPE`: its write gives no error, so a scenario "exit 2 when the table
  cannot be written" through a closed pipe sees a signal, not the code. And an
  `exec.Cmd` with a nil `Env` gives the program the whole environment of the
  test, so a run "with no environment variable" that sets `Env = nil` reads
  them all. Task `T-5zmw` (#83) measured the first before it wrote its
  scenario. **The check:** give a write test an output that refuses each write
  (a file open for reading only, or a `Writer` that returns an error), and give
  an empty environment as an empty slice.
- ❌ **"White space" in the words, `\s` in the expression.** Go's `\s` is only a
  tab, a line feed, a form feed, a carriage return and a space: no vertical tab,
  no U+00A0, no U+2003. A rule that says "white space" and an expression with
  `\s` disagree on these characters, and tests with ASCII input do not show it.
  Review round 1 of task `T-5zmw` (#83) found it in G1, where a value of only
  U+00A0 named a stack. **The check:** write the class that the words name (the
  Unicode `White_Space` is `[\s\v\x{85}\pZ]`), or define the words by the
  class, and give the test a character outside ASCII.
- ❌ **A rule stated as a general claim.** A sentence such as "a link to
  another repository does not count" is a claim about every input, and a
  reviewer meets it with one counter-example per round: task `T-6x75` (#84)
  spent rounds 2 to 4, and three raises of the cycle cap, on one link rule.
  **The check:** state the exact rule (the classes of characters before and
  after the text), and let the examples follow it; test a neighbour of each
  side, and a character of another script.
- ❌ **A check that runs only when its row is there.** A resume check such as
  "compare the hash row of a done step with its inputs now" that runs *if* the
  row exists passes a record with no such row: a record made by hand, or by
  an older version, skips the check, and each test that seeds the row stays
  green. Review round 1 of task `T-79y7` (#85) found it in `answers.sha256`.
  **The check:** for each row that a later run compares, test a record that
  does not have it, and make its absence an error where the writer always
  writes it.
- ❌ **A commit identity with no time.** `git.Commit` gives the time of its
  `Identity` as both dates, `@<seconds> +0000`; the zero `time.Time` is
  `@-62135596800`, which `git` refuses (`fatal: invalid date format`). A test
  helper that commits with an `Identity{Name, Email}` literal fails only when
  it runs. Task `T-7s0y` (#86) met it in an integration test. **The check:**
  give each `Identity` a time (`time.Unix(0, 0)`, or `pin.time` for a setup
  commit).
- ❌ **A mutation run that restores a file with `git checkout`.** A script that
  changes a file, runs the tests, and restores the file with
  `git checkout -- <file>` puts back the committed version, so a change that is
  not committed yet is lost. The mutation script of task `T-c06a` (#91)
  restores so; task `T-7s0y` (#86) found it in its copy before the first run,
  on files that were not committed. **The check:** keep the text of the file
  before the change and write that text back, or commit first.
- ❌ **An escape that arrives as the character.** The input of an agent's tool
  call that holds the escape `\u2039` (in a heredoc of a shell command, in a
  script, or in a file that the tool writes) can reach the file as the angle
  quote itself, so a Go file gets a real marker, which LAYUP's own check
  `markers` refuses at the commit; tasks `T-8ya0`, `T-8vpw` and `T-b3r1` met
  it, and a `grep` with the escape in the shell of the host found nothing.
  **The check:** build the escape from its parts (`chr(92) + "u2039"` in
  Python), and search the changed files for the two characters by their bytes
  before the commit.
- ❌ **A known-bad fixture that fails for another reason.** A fixture run that
  gives `fail` proves the gate only when the kind's command fails on the
  fixture itself. The stand-in manifest of `internal/standin` ran `go vet ./...`
  on a tree with no `go.mod`, so the static fixture of the Go entry (a file that
  `gofmt` changes) would fail with "go.mod file not found", and the row
  `gate:static` would pass for that reason. Task `T-d6q5` (#92) found it before
  the first run and gave the stand-in a `go.mod` and the command of the entry.
  **The check:** give a test target each file that the command of the kind
  needs, and read the output of the fixture run once.
- ❌ **A background maintenance of git.** A commit runs `git maintenance run
  --auto` (since `git` 2.29), whose tasks go to the background (since 2.47), so
  a process of `git` can write into `.git/objects` after the call ends. In CI of
  the pull request #115 (`git` 2.55.0) such a task repacked into a work area
  while a test removed it (`unlinkat …/.git/objects: directory not empty`); the
  same test passed on the LAYUP host (`git` 2.54.0, where a hand run of review
  round 2 of #92 saw no repack after a commit). Task `T-d6q5` (#92) met it.
  **The check:** start each `git` of the product and of a test helper that
  commits with `-c maintenance.auto=false`, and test it with a repository whose
  own configuration asks for the maintenance at once, in the foreground.
- ❌ **A bare repository that receives a push starts its own maintenance.** The `-c
  maintenance.auto=false` of the pushing git (`T-d6q5`) does not reach the git that
  receives: a local transport clears the values of `-c` before it starts `git
  receive-pack`, which then runs `git maintenance run --auto --quiet --detach` in the
  target. From git 2.54 (measured on 2.55.0, CI's; 2.54 read from its source) that
  repacks the target in the background once `objects/17` holds two loose objects, and
  a test that removes its directory a moment after its last push fails in `TempDir
  RemoveAll cleanup` ("directory not empty"): 3 of 70 runs of CI (#291). In the made
  setup git 2.47.3 and 2.56.0 failed none of 100 removals. The receiving repository
  of a local push is not one of those that `docs/spec/packages.md` says no call
  leaves a process in. Task `T-2hm3` (#297) met it. **The check:** set
  `receive.autogc` and `maintenance.auto` to `false` in each bare repository that a
  test pushes into, and test it by pushing with `GIT_TRACE2_EVENT` set and reading
  that no `git maintenance run` starts.
- ❌ **A bare repository with no default branch.** `git init --bare` with no
  `-b` names `master` as its `HEAD` when no setting of the host gives another
  name, so after a push of `main` a clone of it has no checkout ("remote HEAD
  refers to nonexistent ref"), and a test that reads its files fails for a
  reason that is not its own. Task `T-dep6` (#93) met it. **The check:** make a
  bare repository that stands in for GitHub with `-b main`, the default branch
  of the target.
- ❌ **A shared fixture made in one test.** A run that several scenarios share
  (`sync.Once`) and that the first scenario makes in its own `t.TempDir()` is
  removed when that scenario ends; a `t.Fatal` inside `Once.Do` marks the
  once as done with a fixture made in part. Task `T-dep6` (#93) kept its whole
  setup in the directory of `TestMain`, and its maker gives an error, which each
  scenario then reports. **The check:** put a shared fixture where no test
  removes it, and give no `*testing.T` to the code that makes it.
- ❌ **A column that may hold the empty value.** `internal/tsv` takes `—` in
  each column that is not a key, whatever its type, and gives the check of a
  column whose rule forbids it to the package that owns the record
  (`docs/spec/README.md`, the schema block). A validator that checks only the
  rules in words passes a row with no status or no time. Round 1 of #94 found
  it in the telemetry record. **The check:** for each new schema, list the
  columns whose block rule has no clause for `—`, refuse `—` in each, and test
  one row per column.
- ❌ **A git date format that changed.** `git log --format=%aI` gives a UTC
  date as `Z` in git 2.54 and as `+00:00` in older versions, so a test that
  compares the text passes on one host only. Task `T-d6q5` (#92) met it.
  **The check:** compare `%at` and `%ct`, the seconds since the epoch.
- ❌ **A byte-for-byte fixture suite on a checkout with carriage returns.** A
  new suite whose runner compares the output with an `EXPECT` file, and whose
  reader matches a line exactly (`$0 == "## Now"`), passes on every host that
  keeps line feeds. With `core.autocrlf=true`, git's default on Windows, every
  input and `EXPECT` arrive with carriage returns: each case fails, and the
  command in real use drops the tasks under `## Now` with exit 0. Only the
  linters' own fixtures had a `*crlf*` case, so nothing caught it until round 1
  of `T-eep8` (#175). **The check:** a new reader strips carriage returns from
  its inputs, its runner strips them from `EXPECT`, and the suite has a case
  named `*crlf*`, CRLF throughout and pinned `eol=crlf` in `.gitattributes`,
  so `run-discipline-tests.sh` checks that it keeps them.
- ❌ **A fake whose child cannot see `SIGINT`.** A shell script that starts a
  child with `&` gives it `SIGINT` and `SIGQUIT` ignored (POSIX, with job
  control off), and a shell cannot trap a signal ignored at its entry; `trap ''`
  also passes the ignore to each program it starts. So a test of a stop that
  signals a whole process group sees no `SIGINT` in the child, and reads it as a
  defect of the code. Task `T-5pxd` (#161) met it. **The check:** run the child
  in the foreground, under a parent whose trap is `:`, not `''`.
- ❌ **A short time limit that a slow runner reaches.** A table of cases gave each
  case the same wall limit, 300 ms, which only the case `wall` tests. On a slow
  runner a shell that started late was stopped by the limit before it signalled
  itself, and the case "crash by a signal of its own" was classed `wall`:
  `TestEachEndGetsItsClass` failed twice in three hours of CI (#291), and never in
  300 runs of the old test on a host of two CPUs, idle or with a busy loop for each
  CPU, so a quiet host does not show it. Task `T-6rgp` (#292) met it. **The
  check:** give the short limit to the case that tests it and to each other case a
  limit that a slow start cannot reach (5 s), and add a case that runs longer than
  the short limit and expects its own class: it is red while the limit is shared.

- ❌ **A baseline of more than one commit in a world of `layup run`.** With a stand-in
  baseline of two commits the step `read-back` failed in one of the first four runs ("the
  default branch is not one root commit with the tree pin.tree") and passed in every run
  after the baseline was made one commit: the world's Operator pushes the `HEAD` of the root
  directory as soon as it has one, so a world of several commits is green by luck. Task
  `T-eahr` (#277). **The check:** make the baseline one commit (remove its `.git` and
  initialise it again after its files are written), and run a new world several times.
- ❌ **A path that is clean by its spelling still leaves through a link.** A name with no
  `..` and no empty part read or removed a file outside its directory when a parent was a
  link (`result/`, `DIR/work/OWNER`, a link of the records base). Task `T-eahr` (#277),
  round 1. **The check:** `Lstat` each directory above the file, or ask the tree for the
  path (a path below a link has no entry), before you read or remove it. The same trap
  met a **write** in task `T-6zj3` (#361): `applyFrozen` wrote the frozen test files with
  `MkdirAll` and `WriteFile` after `fs.ValidPath`, and a checkout whose parent directory or
  file was a symbolic link let the write overwrite a file outside it (found by the
  producing-code audit, compiled assertion: the outside bytes went from `original` to
  `frozen`). For a write, `Lstat` every existing part of the path, refuse a link (the file
  itself too, and one that points inside the checkout), and make this check for the whole
  list before the first file is written.
- ❌ **A record that no test read through the real check.** S12 in phase 2 wrote its row
  `module` as the answer `S01-name`, which phase 2 does not have; every step passed, and
  the first run of the real `layup setup verify` refused it (check `sources`). Task
  `T-79zk` (#278). **The check:** run the real check on the record of a world that a new
  step writes or that a step in a new phase writes, before the freeze.
- ❌ **A record that keeps a payload "byte for byte" and adds a byte to it.** The plan of
  `setup/rulesets.txt` added a line feed to an output that had none, so a reader could not
  tell the output's last byte from the separator. Task `T-hfn1` (#279), plan review. **The
  check:** write the marker and the separator as bytes of the record that are not the
  payload's, always the same, and test the whole record for a payload that ends with a line
  feed and for one that does not.
- ❌ **A pass test on a text that nobody has read from the real forge.** The probe of the forge check
  passed only for a refusal with `GH013` and the default branch's ruleset name in its standard error; the
  tests' forge refused so (a line `- Rulesets: layup: the default branch` that no source had), and no official page
  gave the text (task `T-msrd`'s search found none). GitHub's refusal names the ref and the violated rule and no
  ruleset: ten public reports show it, and the demo of `M2d` would have failed on a correct target. Task `T-wphp`
  (#280), found before the demo by `T-gkvi` (#323), which attributes the refusal by the ruleset ID of the effective
  rules. **The check:** a test double of an external system's message copies a text from a source that the test
  names (a report, a log of the real run), and the demo keeps the real text (`runs/T-sdtb/`); a rule never
  accepts a refusal because of its text alone.
- ❌ **A fake forge that echoes the file that it was given.** The read-back of the forge check compared each parameter
  of a ruleset file with the effective rules, and the fakes of the job, the world and the e2e gave back each rule as
  the file held it. The real forge does not: the rule `update` of the records ruleset, written with
  `update_allows_fetch_and_merge` `false`, came back with no `parameters`, and the row `rules:layup-records` was
  `fail` on a correct target in run 6 of the demo of `M2d`, with every test green. Task `T-j7c5` (#330). **The
  check:** a fake of the forge's answer takes its form from a first-hand answer, cites it beside the fixture and
  has an assertion on the fake itself; a comparison that accepts an omitted parameter lists the observed cases by
  rule type and parameter name, never "any missing `false`" (GitHub echoes the other `false` parameters).
- ❌ **A claim that a pushed commit is "on no ref" of the clone.** A `git push` to the URL of the clone's own
  `origin` makes `origin/<branch>` when the forge accepts it in the `git` of one host (2.56) and not in that of CI,
  so a test that wants the ref, or its absence, passed on the host and failed in CI. Task `T-wphp` (#280).
  **The check:** test `git for-each-ref --contains` after a refused push and an accepted one, and accept both forms.
- ❌ **A test that a plan review's condition names, written after the code.** The world scenarios of conditions 2
  and 3 were slow to write and went last, so they were never red first; round 1 of `T-wphp` (#280) took it as a
  failed acceptance criterion ("written red first"). **The check:** write each test that a condition names in its
  red slice; one that is late is run against the mutant of its rule, which shows the right reason, not the order.
- ❌ **A "done" test over the rows of a table that its reader did not accept.** "Each row is `pass`" is true of a
  table with no row, so a header alone would have called the Scaffold done and moved `forge-check` before
  `probe`. `scaffoldIsDone` reads the table by its schema first. Task `T-2j13` (#281), plan review. **The check:**
  a header and no row is a case of every test that decides from a table.
- ❌ **A test that compares the times of two rows of two runs as strictly later.** The `at` of a table has the
  resolution of a second, and the forge check of a restart that ends a run and the one of the restart that
  follows it can fall in the same second; the e2e of the Scaffold failed one run in two on it. Task `T-eevt`
  (#282). **The check:** let a second pass between the two runs, or compare the run IDs and the commits.
- ❌ **A mutant run against a script that already fails the clean tree.** The plan said that the script as it
  stands would pass the two push mutants and so show that it "does not bite"; it fails the clean release too
  (`remote in AddRemote`), so each exit 1 was the clean tree's. Task `T-5xjc` (#283), plan review. **The check:**
  run each mutant against a control, the script with only the fix of the clean tree, and read the `FAIL` line of
  each run, not its exit code.
- ❌ **A review brief that carries the plan.** The briefs of the release reviews of `T-4tjy` and `T-8byd` gave the
  reviewer the plan, its review and the author's answer, which the plan review of `T-5xjc` (#283) read as a breach
  of context independence. **The check:** give a release review the issue, the criteria verbatim, the list of the
  release with its line counts, the check with its output and a blank record form, in a clone with `runs/`,
  `docs/tasks/`, `docs/prd/` and `docs/plan/` out of its working tree (a sparse checkout; the history and
  `docs/tests/traceability.md` still hold earlier verdicts, so the brief says what it may read), and compare the
  Coverage table of its record with the commit.
- ❌ **A refusal case whose bad input an older check of the same program refuses too.** The first cases of
  `TestTheJobsReadTheBase` that a `config` path in `.git` or one that is not a path must fail the step made the
  head carry the bad manifest too, so the script as it stood (no mode) refused them with the same words, and
  the case was green before the code. Task `T-hz4m` (#355). **The check:** give the head a clean file and put the
  bad one only where the new code must read it (the base); run the case red against a stub that has the step and
  not the behaviour, and read each failing line, not the exit code.
- ❌ **A work area built again that cannot reproduce the recorded commits.** S13 gained a file in the tree of the target, so the
  area of a target recorded before it would have been rebuilt with another head, and `scaffold-push` and `forge-check` fail on the
  head of `setup/tree.tsv`. It is silent because a new target passes. **The check:** a change to what a setup step puts in the tree is
  a change to every recorded target; key the step to the record (`codeowners.sha256`), run the real step over a real work area in
  both modes and compare what its commit changed, and make the state with no way back a `fail` that names the recovery, never a guess
  (`TestARebuiltAreaReproducesTheRecordedCommits`, `TestSetupAreaChoosesTheMode`). Learned in `T-jew3` (#356).
- ❌ **A body regenerated from the file, sent over a live ruleset.** The file that S13 wrote is not the ruleset that the forge holds:
  the live object of the demo of `M2d` has fields the file never had (`required_reviewers`, `do_not_enforce_on_create`, `id`,
  `source`), and GitHub's page on the update says nothing of an omitted field or of `rules`. A `PUT` of the file would drop whatever
  the Operator added. **The check:** write the live object back whole with the one parameter changed, refuse unless it is the ruleset
  the records name, read it back and compare the reduced object byte for byte, and serve the live object of a first-hand record in the
  test's forge (`runs/T-sdtb/rulesets.txt`; `TestTheUpdateKeepsWhatIsUnknown`). A landing that disables a ruleset puts back the object
  that the update left, never the one before it, and arms its trap only after the disable succeeded. Learned in `T-jew3` (#356).
- ❌ **A mutant that does not build counts as caught.** The first mutation helper of `T-jew3` read only the lines `--- FAIL` and a
  word of `build failed`, and ignored the exit status of `go`: three variants that removed a `%s` and left its argument were refused by
  `go vet` and were counted as kills, so the count of kills was three too high. **The check:** a variant is valid only if `go build` and `go vet` pass on
  it, its kill is the named assertion in its own output, each run's exit status and output are kept, and the file's hash before and
  after proves the restore (`runs/T-jew3/mutants.md`). Learned in `T-jew3` (#356).
- ❌ **A check that admits a head by who made it.** The author, committer, message, date and signature of a commit are text the pusher writes,
  so an admission keyed to them proves nothing. **The check:** judge the parents, the whole tree and the bytes of the one file, and test both
  ways: the Operator's name with wrong content is refused, hostile metadata with the canonical content is admitted
  (`TestCheckCodeOwnerLanding`). Learned in `T-4cep` (#384).
- ❌ **A script and the check that judges its result, tested apart.** Each had a green test of its own and nothing showed that they agreed, so
  the script could land a commit that the next run refused. **The check:** one test runs the producer as production writes it and has the
  checker judge what it made, and a producer that adds another path or a second file is refused by the same call
  (`TestTheLandingCommitIsAdmitted`). Learned in `T-4cep` (#384).
- ❌ **A fixture that adds a submodule entry and then `git add --all`.** `Add` stages the tree as the work tree is, so an entry with no file is
  dropped and the case goes green without it. **The check:** add the entry after `Add` and read the listing (`TestParentsAndTheWholeTree`).
  Learned in `T-4cep` (#384).
- ❌ **A set that is "not yet applied" because its comment is in no row.** The plan of the parameters (`T-dx1h`,
  #358) judged each comment of the control issue again at every boundary and skipped a set only when its ID was the
  `comment` of some row. After a later set of the same name, the earlier one's ID is in no row, so it is applied again:
  a commit at every boundary, and the value goes back to the older one for a pass. It is silent until a second set of one
  name exists. Found by the plan review (C2), before any code. **The check:** decide "applied" by order and not by
  membership (apply a set only when its comment ID is greater than the row's recorded `comment`), and give the test two
  sets of one name, a second boundary that writes no commit, and the earlier set listed after the later one.
- ❌ **An owner check written for the rows of a table, not for the rows that the first writer writes.** The plan of
  the same task said "`evidence` and `comment` are `—` together" and also gave the three rows that take their value
  from `start.tsv` the evidence `command` and no comment, so the file the writer wrote would have been refused by its
  own reader (plan review, C1). **The check:** write the reader's rules as the kinds of row (here: a row with no
  comment is its own origin, a row with a comment is a set) and give the test the writer's output before any other
  case, so the writer and the reader cannot disagree.
- ❌ **A negative assertion on a records file that the world has already written to.** The test of a post that
  fails asserted that `copies.tsv` held no `copies/` text, and it was false on a correct tree: Start had copied the
  dead-man notice of the control issue (`copies/501-1.md`), so the file always held one. Task `T-q5k4` (#360),
  found when its integration tests were run against a stub. **The check:** count the rows that belong to the thing
  under test (here the rows whose `issue` is the task's), never the text of a file that other steps write, and
  run a "nothing was written" case against the real code once before the stub.
- ❌ **A limit quoted from the forge's error message and not measured.** The comment of a task is posted as one
  body, and `loop.md` took the forge's own words for its size ("maximum is 65536 characters"). Probes on
  2026-10-16 (`runs/T-q5k4/`) found 65,537 characters, 131,072 bytes and 262,143 bytes accepted, and 131,072
  three-byte letters (393,216 bytes) refused with that same message, so the message is not the rule that is
  enforced. Task `T-q5k4` (#360). **The check:** when code relies on a forge limit, read the answer of a probe at
  both sides of the number, say which unit it is counted in, and keep the documented number as the side that
  never posts a body the forge could refuse; do not write the limit as a fact on the strength of the message.
- ❌ **A fixture that sets a field on the struct and so bypasses the production constructor.** `Sessions.App` was
  added for the adoption of the App's own comment after a lost copy, and its test passed because the test world set
  `App` by hand; the only constructor that a run uses, `state.sessions()`, never set it, so the recovery was dead code
  and a restart posted a second comment. Task `T-q5k4` (#360), round 1, finding 1. **The check:** when a change adds a
  field that a run must carry, build the object under test in at least one behavioural test through the constructor
  that the run uses (here `state.sessions()` on a `Config` with the App's slug), and run that test against the mutant
  that removes the assignment; a test that only reads the copied field, or sets it itself, proves nothing about the run.

- ❌ **A display name used as a routing key.** `DefaultOwners()` of `internal/run` gave the owner of a
  question as `Systems Architect`, while `routing.tsv` and the step table key a role by its ID
  (`systems-architect`); the route of a question would have matched no routing row. It is silent because
  the route was tested alone, no caller used it with the register yet, and each value is right in its own
  file. Task `T-6zj3` (#361), found when its plan read the two sources side by side. **The check:** turn
  the name into the ID by one rule (`RoutingRole`: lower case, spaces as hyphens), and test that each
  default owner that the step table has as a role comes out as exactly that ID, with the roles read
  from `session.StepTable()`, never typed in the test.

- ❌ **A required intake field changes only the integration fixtures.** Adding
  success criteria and the revision step in #365 left the CI-only end-to-end
  worlds on the old question counts, setup-question index and scripted outputs.
  **The check:** inspect fixtures under every build tag when a required field or
  step changes; preserve their behavior assertions, supply the new answer and
  scripted output, and require the full CI ladder before landing. A local
  integration pass does not stand for the CI-only end-to-end level.

- ❌ **A doc comment that names a role the code never reads.** `AnswerAccepted` said an answer is "of a session of
  the owner role", took only `events`, `results` and the answer ID, and so accepted the answer of a developer session
  with every record valid (`CheckSession` takes `role` as free text, and `RecordHandoff` reads no session's role).
  Task `T-b3s3` (#401), revealed by the first review of #361. **The check:** when a contract names a role, the
  function takes the role from the records (`sessions.tsv`, written by `layup run`, never the result of the session
  itself) and a test runs two ledgers that differ only in that role; give the two maps of such a function different
  named types, so that a swap does not compile.

- ❌ **A tree reader that lists recursively takes links, gitlinks and nested files for files.** `git ls-tree -r`
  lists a symlink as a `blob` of mode `120000` whose bytes are its target text, a file below the directory by its
  full path, and a gitlink as type `commit`; a reader that keeps every blob under `docs/spec` would read all three
  as spec files. It is silent because the control, a regular file at depth one, passes. **The check:** filter on
  depth, type and mode, and give each shape a case where the section sits only in that shape, so a wrongly included
  file flips the verdict. Learned in `T-yeav` (#404).
- ❌ **An error test whose fixture fails at the first call.** Five failing calls of git were injected into a call that
  also rendered a PRD, so the pre-image read failed first and the mutant that swallowed the later read of the head's
  tree survived. It is silent because every row of the table was an error. **The check:** give each call its own row
  where the earlier calls succeed (here, a facts-only call), and keep the mutant that swallows that call.
  Learned in `T-yeav` (#404).

### Reference-sweep pitfalls

A change that edits references or a rule's wording across the tree has three silent
failure modes worth keeping.

- ❌ **A blanket find-and-replace over a renamed record's citations.** When a record
  moves or a directory is renumbered, the same bare token can name *different*
  records in two places — a bare `ADR-0005` is the living `docs/adr/` record to one
  reader and the archived `docs/decisions/` one to another, because the two sequences
  once shared numbers (setup step S05 deleted that archive from this repository;
  the citations of ADR-0001 to ADR-0008 to it stay). A global replace of the token silently rewrites the citations
  you must **not** touch alongside the ones you must; and the reverse — a citation the
  sweep's pattern never matched (a compound like `ADR-0003/0005`, a token in a code
  span or a `.sh`/`.yml` comment, one split across a line break) — is silently *left*
  pointing at the wrong record. It is silent because **no linter catches it**:
  `link-lint` checks only that a *link* resolves, and a bare textual mention resolves
  to nothing, so a citation that now sends a reader to the wrong record still passes
  every check. **The check:** classify each occurrence by its **link target**, not its
  token — a link into `../adr/` is the living record and stays, a link into
  `../decisions/` is the archive and is rewritten — and read every *bare* mention by
  hand, in every token shape, since it carries no path to classify it. A
  pre-registered grep that must finish returning only the intended survivors (the
  mapping table and deliberate historical prose) is the closest thing to a gate; run
  it against the whole tree, not only the files you expected to touch.
- ❌ **A repoint that orphans a bare back-reference.** Repointing a citation can
  strand a *different* reference that named the target only through it. A comment
  reading `section 6 says …` leaned on a nearby `D-0003 section 6` for its antecedent;
  repoint every `D-0003 section 6` and the bare `section 6` is left pointing at a
  structure only the deleted record holds — wrong on this tree, and sharing
  **no token** with the thing you renamed. It is silent because a grep keyed on the
  obvious token (`D-000N`) cannot match a bare `section 6`, so the pre-registered
  check goes green over the survivor. **The check:** grep for the *shapes* a reference
  takes, not only the token — a bare `section N`, a `§`, a pronoun (`that section`,
  `the record`) whose antecedent you removed — and read the neighbourhood of every
  citation you changed, not the citation alone.
- ❌ **Editing a rule whose decision record is archived.** A rule lives in two places
  — its operative statement in a living doc, and the immutable decision record that
  first set it (in the baseline, under `docs/decisions/`, which setup step S05
  deleted from this repository). Change the living one and the archived one
  still asserts the old, and **no check compares them** (`adr-lint` never reads
  `docs/decisions/`; `link-lint` checks resolution, not agreement). You cannot rewrite
  the immutable body to match; discharge the divergence with a `Status`-line
  **amendment pointer** on the archived record. **The check:** grep the whole tree —
  archive and forge templates included — for the old wording, and reconcile each living
  mirror or point each immutable one; a dated log entry recording history stays.
- ❌ **A hand-mirrored count or check-set that no linter guards.** The set of discipline
  linters — and how many there are — is spelled out by hand across many living
  docs, among them [`engineering-discipline.md`](engineering-discipline.md), this file,
  [`tests/test-levels.md`](tests/test-levels.md), the two `tests/README.md` files,
  [`ci/README.md`](ci/README.md) and
  [`.githooks/README.md`](../.githooks/README.md). Add or remove a check and every one
  can go stale, and a **removed** check leaves its name behind as a linter that no longer
  exists — a `link-lint` run stays green, because it resolves a *link*, not a claim. It is
  silent because the sentence still reads well and the count still looks deliberate: the
  documents once said `three`, `four` and `five` at once, and named an `agent-entry` linter that
  had been cut. **The check:** when you add or remove a discipline check, grep the whole
  tree for the check-set enumeration — the old name and each spelled count — and reconcile
  every living mirror in the same change; the immutable ADR and archived decision copies
  stay as history.
- ❌ **A rule on the clock, in a test world on another clock.** A rule that compares
  `Now` with a recorded time (the hours since Start of milestone 0) breaks every test
  world whose records were written on a clock of their own: the world of
  `internal/run` starts on 2026-10-08 by its fast clock, so with the real `Now` each
  session of a phase was refused "past 8.0 hours", and six green tests of another row
  went red at once. It is silent in a test that asserts only the refusal, which
  passes for the wrong reason. **The check:** in a test world, set `Now` from the
  world's own records (its time of Start), never from the wall clock, and keep one
  case inside the cap so a refusal for the wrong reason shows. Learned in `T-ybd5`.
- ❌ **`git diff --check` before `git add` reads no new file.** The check reads the
  diff of the tracked files; a new file, untracked until `git add`, is not in it, so
  a record pasted into a new file (a test run cut at a width, with its trailing
  spaces) passes the check and fails it after the commit. It is silent because the
  command exits 0. **The check:** run `git diff --check <base>` after `git add`, or
  on the commit, never on a tree with new untracked files. Learned in `T-hpid`, and
  met again in `T-j304` (#303): the green script ran it with 11 new case directories
  untracked, and six fixture files showed trailing tabs only at the first `git add`.
  A row of a TSV fixture whose last field is empty ends in a tab: write the row
  without it where the reader takes a missing field as empty (`awk -F'\t'` does),
  or let the writer of the table leave the tab out.
- ❌ **The writer of a command's table in `internal/cli`.** `internal/cli` may not
  import `internal/tsv` (`packages.md`, its row), so a command that prints a table
  by its schema needs the writer in the package that owns the table, as
  `psb.WriteTSV` and `spec.WriteTSV` are; `TestPackageRules` refuses the other
  form at the integration level only, after the unit level is green. **The
  check:** give the engine package its `WriteTSV` with the code that builds the
  rows. Learned in `T-hpid`.

- ❌ **A prompt that lists the columns of a row in one sentence.** "A span that is
  not a fact has "—" as its fact and its class, and its reason" reads as "`—` for
  each of the three" as easily as "`—` for two, and a reason". Three sessions of
  the demo of `M2c` took the first reading on every heading. The block refused
  each file, so `numbering` failed on every try, while every unit test passed:
  they write the file by hand. It is silent until a real session runs. **The
  check:** in a prompt, give each column of each kind of row its own clause, as
  "X in the column C"; say what a column is never; and test the prompt's text. Find
  a refused file by copying the session's `result/` before its directory is
  removed. Learned in `T-pzgw`.
- ❌ **A prompt that gives the producer the task and not the rules of the check that judges it.** The prompt of
  the prose session of the step `setup` gave none of the rules of the five checks that accept its files. On the
  first real target (the demo of `M2d`) S08 failed on a citation of an example fact of the baseline, and S14 then
  refused 222 hits of the baseline's words in 20 files and a `README.md` with no name: two runs, each leaving a
  lease held for 30 minutes, with every unit test green. Task `T-92lm` (#326). **The check:** such a prompt names
  each rule with the files it governs, from the check's own table where it has one, so that a test probes the table
  against the check (a first text that said "a letter" for "an ASCII letter" was looser in five places); a test ties
  each rule to the prompt; and one real session on a copy of the real input is the evidence before a real run.
- ❌ **A specification that changes a cell or a block and leaves the test that reads it unchecked.** The task that
  specified `M2e` changed the cell "Starts a program" of the row `internal/spec` with a second code span in it, and
  `TestPackageRules` refused it ("cannot read Starts a program"); a new schema block fails
  `TestEverySchemaBlockIsBuiltOrNotYetBuilt` until its name is listed; and `setup-check` `markers` reads only the
  tracked files, so a new file that holds a marker reports its row of `open-gaps.tsv` as stale until it is added.
  Task `T-xs4y` (#333). **The check:** after a change of `packages.md` run `go test -tags=integration ./cmd/layup/`,
  after a new block `go test -tags=integration ./internal/tsv/`, and `git add` a new file before `setup-check`.
- ❌ **A piped exit status that is the last command's, not the test's.** A run of
  `go test ./... | tail; echo exit=$?` reports the status of `tail`, normally 0
  when `tail` succeeds, so a red suite can read green. Earlier evidence also used
  unmeasured future timestamps from a draft instead of measured run times.
  The pipeline can succeed while the tests fail and its output looks like a result. **The check:** capture the
  real status — run the command directly, or `set -o pipefail`, or read
  `${PIPESTATUS[0]}` — and record only times `date -u` was read for; a file that
  cites a run names its start, end and exit, or says the run was interrupted
  and claims nothing. Task `T-3q0b` (#350); found by the root audit.
- ❌ **A mutation sweep whose restore brings back stale copies.** A harness
  that keeps backups at a fixed path (`/tmp/*.go.bak`) can restore copies
  from an older run and silently revert the fix under test; every mutant
  then runs against the pre-fix source and its results prove nothing —
  they are discarded, not counted. Task `T-3q0b` (#350). **The check:**
  take the backup fresh from the frozen head, verify what each restore
  wrote (content or hash against the head's file), and re-run a green
  control on the restored tree before crediting a single mutant.

## 3. Validation — how you check you are not fooling yourself

A result is **untrusted** until it passes the checks below, and the pass is a
recorded event, not a memory. Order the checks cheap-first, so a failure stops the
expensive ones.

| # | Check | Pass condition | Cost |
|---|-------|----------------|------|
| 1 | sh docs/setup/setup-check.sh | every check prints OK; exit 0 | minutes |
| 2 | sh docs/setup/tests/run.sh | every fixture case matches its EXPECT | seconds |
| 3 | the review rounds on a frozen head | the last round says `nothing material in scope` | one reviewer session per round |

Notes on how to read a failure: `setup-check.sh` catches a marker that is neither filled nor listed as an open
gap, and a missing pin, fact, or required section; the fixture
run catches a check that cannot fail; the review rounds catch a claim that no
script can settle. The first two are cheap enough for the hook and CI; the
rounds run once per change.

**The automated gate is this validation layer, mechanized.** The cheap, always-on
checks — the [discipline linters](engineering-discipline.md#testing) (ADR, PRD
and link) and their
[fixture self-tests](engineering-discipline.md#testing), the
[test levels](engineering-discipline.md#testing), lint, a security
scan, and the [commit-format](engineering-discipline.md#commit-messages)
check — run in the [`pre-commit` hook](engineering-discipline.md#git-hooks) for
fast local feedback and in [CI](engineering-discipline.md#continuous-integration)
as the authority. Treat those checks as pre-registered pass/fail rules under
section 1: they predate any single result and are not edited to make a change
pass. Wire the "cheap enough to wire into CI" checks from the table above into
both layers.

## 4. Mechanics

- **Frozen rules do not get edited.** Changing a guardrail after it is set means a
  new version with a written reason, the old one preserved. Legitimate reasons
  exist (a bug in the measure); silent edits do not.
- **This document holds the structure; Git (`F-0001#1`) holds the frozen
  values.** When real values exist, mirror them here as history, after the fact,
  never as the primary copy.

## Sources

- [`F-0001`](facts/F-0001-layup-problem-statement-brief.md) — the PSB, §6 and §7.
- [`setup/record-T-n1hp.md`](setup/record-T-n1hp.md) — the evidence for each setup value.
