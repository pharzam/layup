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
  github.com/zricethezav/gitleaks/v8@v8.30.1 git --redact --log-opts="origin/main..HEAD"`,
  and write a credential into an environment as `env["NAME"] = value`. Learned in
  `T-zwke`.
- ❌ **A one-line form that loses a failure.** The first lint command of the first
  pilot put the gate script of the brief into one pipe. When `git ls-files` failed
  (an invalid index), the pipe took the status of its last command, and the check
  passed with no file read. It is silent because the form gives the script's
  output in each case that works. **The check:** before a short form replaces a
  script, run both on a failure of each step, and keep the cases (13 cases in
  `T-evad`: F-9 of [`runs/T-evad/findings.md`](../runs/T-evad/findings.md)). Learned
  in `T-evad`.
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
  assert the line as well as the column. When a mutation record is cut to a few
  lines, say so, or it reads as fewer failing cases than ran.
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
