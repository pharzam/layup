# Setup procedure — Armature for a new project

This is the procedure that set up this repository by hand, written as ordered
steps so that it can be repeated, measured, and later automated. The measured
run is [`record-T-n1hp.md`](record-T-n1hp.md). In this run each step was a task
under the full quality gate (an issue, a reviewed plan, test first, independent
review rounds, a pull request); the automation keeps that gate. The same steps are in
[`steps.tsv`](steps.tsv) in a form a program can read; check `procedure` in
[`setup-check.sh`](setup-check.sh) keeps their step IDs in step.

## In plain terms

> Fifteen steps turn a copy of the Armature template into a project whose every
> setup value has a source. Four steps need a person to decide; the other eleven
> can be done by a program. One script, `sh docs/setup/setup-check.sh`, proves the
> result, and CI runs it on every change.

## How to read a step

Each step names its input, its action, its output, the evidence that proves it
was done, and whether a human must decide (`yes`) or a program can do it
(`no`). A step with `human_decision` `yes` is a stop point: the automation asks
and waits. A value with no source is never filled; it becomes an open gap
([`open-gaps.tsv`](open-gaps.tsv)) and a question in the next batch (step S10).

## For the automation and the interactive TUI

- `steps.tsv` is the contract: one row per step, the header
  `id, input, action, output, evidence, human_decision`, tab-separated.
- The `yes` rows (S01, S03, S10, S13) are the interactive screens; each one
  collects its answers in one batch and writes them to Git before the next step
  (PSB Invariant 1, Decision Point 2).
- A `no` row is done when its evidence holds. In this run the evidence of most
  steps is one check of LAYUP's own `setup-check.sh` (`--only <check> <root>` runs
  one check). Armature does not ship that script, and nothing from LAYUP goes into a
  target repository beyond the setup output (O-10 on
  [#30](https://github.com/pharzam/layup/issues/30), O-76): for a target, `layup setup verify` does these checks from outside, the
  Operator pushes the unmodified copy as the root commit (S03), S02 copies with
  `git clone`, S12 adds no `setup-check` job, and S13's checks go into rulesets
  ([architecture](../architecture.md) §5,
  [ADR-0016](../adr/0016-put-the-native-stack-gates-in-the-target.md)).
- The manual run gives the baseline to compare against: the elapsed time of each
  step, and each finding K-01 to K-08 where Armature left a decision open, in
  [`record-T-n1hp.md`](record-T-n1hp.md).
- **Known limit** (the added scope of [#21](https://github.com/pharzam/layup/issues/21),
  kept as a known limit by the Operator's decision O-138 there; the task that
  next changes `check_ci` fixes it): `check_ci` of `setup-check.sh` types the
  list of the pass-through linters and the two `docs/ci` suites of the runner by
  hand, and four of its rule branches have no fixture that can fail (the
  expansion of the four linters, the exemption of the nested checkout, the
  `ci_named` half of the runner rule, and a quoted Restore name).

## Steps

### S01 — Operator decisions before the copy

- **Input:** PSB of the target project; Armature repository URL
- **Action:** Ask the Operator the project decisions that no file answers: the technology stack; the repository name and visibility; the gate mode (full: each task gets a reviewed plan and review rounds; light: branch, tests and one PR per step)
- **Output:** Operator decisions in Git (setup record, Operator decisions table)
- **Evidence:** The Operator's answers, with the date
- **Human decision:** yes
- **Done in this run by:** the session, after S02 and before S03 (in this run the copy came first; the decisions do not depend on it)

### S02 — Copy Armature at one commit

- **Input:** Armature repository URL
- **Action:** Read the current main commit with gh api repos/<owner>/armature/commits/main; copy that exact commit with npx degit <owner>/armature#<sha> <dir>
- **Output:** A project directory with the Armature files and no Git history
- **Evidence:** The gh api output (SHA); the degit command line
- **Human decision:** no
- **Done in this run by:** `T-n1hp` (timeline rows 2–3)

### S03 — Root commit and remote

- **Input:** The copied directory
- **Action:** git init; commit the unmodified copy as the root commit; create the remote under the name and visibility from S01 (the Operator authorizes the create); push main BEFORE the hooks are installed (the pre-push hook refuses a push to main)
- **Output:** Root commit = the Armature tree; remote main
- **Evidence:** git rev-parse <root>^{tree} equals the Armature commit tree (GitHub API)
- **Human decision:** yes
- **Done in this run by:** `T-n1hp` (rows 3–4; finding K-05)

### S04 — Pin the version

- **Input:** Root commit, Armature SHA
- **Action:** sh .githooks/install.sh (it sets core.hooksPath to the relative .githooks); write docs/setup/armature.pin (source, commit, tree, method, date); ADR for the pin
- **Output:** Pin file; ADR; check pin passes
- **Evidence:** git config core.hooksPath prints .githooks; setup-check pin OK
- **Human decision:** no
- **Done in this run by:** `T-n1hp` (the hooks install, record V-22) and `T-r7zg` (the pin)

### S05 — Remove Armature's own history

- **Input:** Armature's own history
- **Action:** Delete docs/decisions/, docs/audit/, Armature's docs/tasks/T-*.md, completed-log entries, backlog lines and notes; fix the links
- **Output:** No Armature history; check kit-history passes
- **Evidence:** setup-check kit-history OK; link-lint OK
- **Human decision:** no
- **Done in this run by:** `T-vbwc`

### S06 — Store the facts

- **Input:** PSB file and other source briefs
- **Action:** Copy each file byte-identical into docs/facts/; write F-NNNN records with numbered verbatim facts; hash the files into docs/setup/facts.sha256; index rows
- **Output:** Raw facts with IDs; check facts passes
- **Evidence:** sha256 values; setup-check facts OK
- **Human decision:** no
- **Done in this run by:** `T-fvwj`

### S07 — Bind the onboarding

- **Input:** F-0001
- **Action:** Rewrite docs/onboarding-for-engineers.md from the PSB, each claim cited as F-0001#n; keep the process sections
- **Output:** Onboarding bound to the PSB; check onboarding passes
- **Evidence:** Semantic-agreement review round; setup-check onboarding OK
- **Human decision:** no
- **Done in this run by:** `T-vpty`

### S08 — Merge the domain terms

- **Input:** PSB terms section
- **Action:** Add each PSB term row to docs/glossary.md word for word, with collision notes and its F-0001#n
- **Output:** Glossary domain section; check glossary passes
- **Evidence:** Semantic-agreement review round; setup-check glossary OK
- **Human decision:** no
- **Done in this run by:** `T-xgz4`

### S09 — Encode the invariants

- **Input:** PSB System Invariants
- **Action:** Add one guardrails entry per invariant: statement with F-0001#n, trap, and Check (path plus gate, or no check yet)
- **Output:** Guardrails invariants; check guardrails passes
- **Evidence:** setup-check guardrails OK; review of each Check value
- **Human decision:** no
- **Done in this run by:** `T-7ndb`

### S10 — Ask the gap questions in one batch

- **Input:** All remaining markers
- **Action:** List every marker; separate those a file or command answers from those that need a decision; ask the Operator every open question in ONE batch
- **Output:** Operator decisions in Git; the questions that stay open
- **Evidence:** The Operator's batch answers, copied into the setup record
- **Human decision:** yes
- **Done in this run by:** `T-nfh8` (Operator decisions O-1 to O-8)

### S11 — Fill the markers or list the gaps

- **Input:** Markers, decisions, command references
- **Action:** Fill each marker that has a source and write its record row (value, file, evidence, status active / not active / recorded); list the rest in docs/setup/open-gaps.tsv with a question
- **Output:** No unlisted marker; check markers passes
- **Evidence:** Record rows; setup-check markers OK
- **Human decision:** no
- **Done in this run by:** `T-nfh8`

### S12 — Activate CI

- **Input:** Armature's active workflows
- **Action:** Keep Armature's workflows and replace their headers; add a Restore step to every workflow that runs a check script and lacks one (at the pin, pr-link.yml and review-record.yml lack it: finding K-08); add the setup-check job with fetch-depth 0 and a Restore step
- **Output:** CI runs every check from the default branch; check ci passes
- **Evidence:** CI run logs show each script restored from the default branch; setup-check ci OK (cause restore)
- **Human decision:** no
- **Done in this run by:** `T-q344` (headers, setup-check job) and `T-fvng` (the missing Restore steps, found by review as blocker #23)

### S13 — Require the checks

- **Input:** Job names of every workflow
- **Action:** Write docs/setup/branch-protection.json (one required check per job, pinned to the GitHub Actions app); wait until each job has reported once; review the body; PUT it; read it back
- **Output:** main requires every job; check protection passes
- **Evidence:** gh api read-back equals the file
- **Human decision:** yes
- **Done in this run by:** `T-afa5`

### S14 — Rewrite the identity

- **Input:** README.md, AGENTS.md
- **Action:** Rewrite the identity text to name the project and link the pin; keep every rule text
- **Output:** Entry files describe the project; check identity passes
- **Evidence:** Semantic-agreement review round; setup-check identity OK
- **Human decision:** no
- **Done in this run by:** `T-6rg3`

### S15 — Record the procedure

- **Input:** The steps above
- **Action:** Write this procedure and the measured record; replace Armature's adoption section of docs/engineering-discipline.md with a link to it
- **Output:** docs/setup/README.md and steps.tsv; check procedure passes
- **Evidence:** setup-check procedure OK
- **Human decision:** no
- **Done in this run by:** `T-9mmm`

## Check `adapted`

Check `adapted` of [`setup-check.sh`](setup-check.sh) reads the documents of this
repository for the rules 1 to 3 of [#70](https://github.com/pharzam/layup/issues/70):
a document speaks about LAYUP, names Armature only where it states a fact about
LAYUP, and states each choice as a result. Check `markers` enforces rule 4. Both
run in the `pre-commit` hook and in CI. They are not a step of this procedure;
the same checks become part of `layup setup verify` for a target later (#70
rule 5).

**Known limit of check `adapted`.** It reads the tracked Markdown files only. A
comment in a script, a row of `steps.tsv` or a line of a workflow can still speak
as the Armature template; on `9beff4c` the plan review of #70 counted 27 such
lines in `setup-check.sh`, 7 rows of `steps.tsv`, and lines in the three linters
and `ci.yml`. The records (facts, `runs/`, task records, the completed log, the
accepted ADRs 0001 to 0012, `record-T-n1hp.md`) are not read, because a record is
not edited.
