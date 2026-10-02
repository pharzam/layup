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
value, a script that fails on each gap, and an open question where no source
exists (`F-0001#4`, `F-0001#5`).

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
  ([ADR-0012](adr/0012-build-layup-in-bootstrap-mode.md), part 1); and the
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
  ID with the "not used" list of
  [ADR-0012](adr/0012-build-layup-in-bootstrap-mode.md), and write full model IDs
  in every evidence file. Learned in `T-hbw8`.
- ❌ **A word-for-word copy with trailing spaces.** A comment or a tool output
  copied into a record keeps its trailing spaces. No hook runs `git diff --check`,
  so 17 such lines passed every commit of `T-hbw8` and failed only at close-out.
  **The check:** after you copy text into a record, run `git diff --check` on it;
  remove trailing spaces, which carry no meaning there. Learned in `T-hbw8`.

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
  is broken; the red step is the proof.
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
