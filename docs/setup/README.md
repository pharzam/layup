# Setup procedure — Armature for a new project

This is the procedure that set up this repository by hand, written as ordered
steps so that it can be repeated, measured, and later automated. The measured
run is [`record-T-n1hp.md`](record-T-n1hp.md). The same steps are in
[`steps.tsv`](steps.tsv) in a form a program can read; check `procedure` in
[`setup-check.sh`](setup-check.sh) keeps the two in step.

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
- A `no` row is done when its evidence holds; for most steps that evidence is one
  check of `setup-check.sh` (`--only <check> <root>` runs one check).
- The manual run gives the baseline to compare against: the elapsed time of each
  step, and each finding K-01 to K-08 where the kit left a decision open, in
  [`record-T-n1hp.md`](record-T-n1hp.md).

## Steps

### S01 — Operator decisions before the copy

- **Input:** PSB of the target project; Armature repository URL
- **Action:** Ask the Operator the project decisions that no file answers: stack, repository name and visibility, gate mode
- **Output:** Operator decisions in Git (setup record, Operator decisions table)
- **Evidence:** The Operator's answers, with the date
- **Human decision:** yes
- **Done in this run by:** session, before `T-n1hp`

### S02 — Copy the kit at one commit

- **Input:** Armature repository URL
- **Action:** Read the current main commit with gh api repos/<owner>/armature/commits/main; copy that exact commit with npx degit <owner>/armature#<sha> <dir>
- **Output:** A project directory with the kit files and no Git history
- **Evidence:** The gh api output (SHA); the degit command line
- **Human decision:** no
- **Done in this run by:** `T-n1hp` (timeline rows 2–3)

### S03 — Root commit and remote

- **Input:** The copied directory
- **Action:** git init; commit the unmodified copy as the root commit; create the remote; push main BEFORE the hooks are installed (the pre-push hook refuses a push to main)
- **Output:** Root commit = the kit tree; remote main
- **Evidence:** git rev-parse <root>^{tree} equals the Armature commit tree (GitHub API)
- **Human decision:** yes
- **Done in this run by:** `T-n1hp` (rows 3–4; finding K-05)

### S04 — Pin the version

- **Input:** Root commit, Armature SHA
- **Action:** sh .githooks/install.sh; write docs/setup/armature.pin (source, commit, tree, method, date); ADR for the pin
- **Output:** Pin file; ADR; check pin passes
- **Evidence:** setup-check pin OK
- **Human decision:** no
- **Done in this run by:** `T-r7zg`

### S05 — Remove the kit's own history

- **Input:** The kit's own history
- **Action:** Delete docs/decisions/, docs/audit/, kit docs/tasks/T-*.md, kit completed-log entries, kit backlog lines and notes; fix the links
- **Output:** No kit history; check kit-history passes
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
- **Action:** Rewrite docs/onboarding-for-engineers.md from the PSB, each claim cited as F-0001#n; keep the kit's process sections
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

- **Input:** All remaining adopter markers
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

- **Input:** The kit's active workflows
- **Action:** Keep the kit's workflows (they restore the checks from the default branch); replace their headers; add the setup-check job with fetch-depth 0 and a Restore step
- **Output:** CI runs every check from the default branch; check ci passes
- **Evidence:** CI run logs; setup-check ci OK
- **Human decision:** no
- **Done in this run by:** `T-q344`, `T-fvng`

### S13 — Require the checks

- **Input:** Job names of every workflow
- **Action:** Write docs/setup/branch-protection.json (one required check per job, pinned to the GitHub Actions app); review it; PUT it; read it back
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
- **Action:** Write this procedure and the measured record; replace the kit section How to adapt this kit with a link to it
- **Output:** docs/setup/README.md and steps.tsv; check procedure passes
- **Evidence:** setup-check procedure OK
- **Human decision:** no
- **Done in this run by:** `T-9mmm`
