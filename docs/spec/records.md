# The records

Where each record of LAYUP lives, who writes it, and the schemas of the
records that phase 1 needs. Conventions: [`README.md`](README.md). The records
of a target are on its branch `layup-records`, an orphan branch that only
`layup run` writes ([`architecture.md`](../architecture.md) §3, ADR-0014).

## The layout of the records branch

One row per record kind of the whole architecture (§1 to §15), so that a later
milestone adds schemas and does not move a file. The inventory that this table
was checked against is [`runs/T-0drh/checklist.md`](../../runs/T-0drh/checklist.md),
part C. "Phase" is the phase of the milestone whose code first writes the
record, or first reads it for a record that a human writes (K38, task
`T-zck8`); "1, schema only" marks a record whose schema phase 1 gives and whose
writer or first reader comes later; a schema is in this directory for a record of
phase 1 and for a record of each milestone whose specification task has run
(`M2a`: [the records of Start](#nfr-001--the-records-of-start); `M2b`: [the
records of a session](#nfr-001--the-records-of-a-session)).

**Decided here:** every path that the architecture does not name (it names 18
files, most without a directory), and the phase of each record. The paths
follow two rules: a table that the whole target shares is at the root; a table
or payload of one task is under `tasks/<task>/`.

### On the records branch (`records:`)

| Path | Record | § | Writer | Phase | Schema |
| ---- | ------ | - | ------ | ----- | ------ |
| `README.md` | what the branch is, and how to read it with no tool | 3 | `layup setup` (O-115); `layup run` for a target that Start makes (`M2a`), replaced at the Scaffold (`M2d`) | 1 | a fixed text per writer: [`setup.md`](setup.md#the-readme-of-the-records-branch); [`run.md`](run.md#the-readme-of-a-target-that-start-makes); [`scaffold.md`](scaffold.md#the-readme-of-the-scaffold) |
| `start/problem-statement.md`, `start/vision.md` | the two briefs, byte for byte | 5 | `layup run` | 2 | — |
| `start/start.tsv` | each brief's SHA-256, the LAYUP version, the forge plan, the values of the Start command | 5, 13 | `layup run` | 2 | [below](#nfr-001--the-records-of-start) |
| `setup/record.tsv` | the setup record: each value with its source, the pin rows included | 5 | `layup setup`; committed by `layup run` in phase 2 (`M2d`) | 1 | [`setup.md`](setup.md#the-setup-record) |
| `setup/verify.tsv` | the table of `layup setup verify` | 5 | `layup setup verify`, committed by `layup setup` (O-115); by `layup run` in phase 2 (`M2d`) | 1 | [`setup.md`](setup.md#the-checks-of-layup-setup-verify) |
| `setup/rulesets.txt` | the output of the Operator's ruleset commands (the bypass lists) | 3 | `layup run` | 2 (`M2d`) | a rule, no block: [below](#nfr-001--the-records-of-the-scaffold) |
| `setup/tree.tsv` | the root commit and each setup commit that `layup setup verify` read | 5 | `layup run` | 2 (`M2d`) | [below](#nfr-001--the-records-of-the-scaffold) |
| `setup/inputs/files/<path>` | the text of a prose row of the setup, written by a role session | 5 | `layup run` | 2 (`M2d`) | a rule, no block: [below](#nfr-001--the-records-of-the-scaffold) |
| `setup/forge-check.tsv` | the effective rules read back, and the probe results | 3, 5 | `layup run` | 2 (`M2d`) | [below](#nfr-001--the-records-of-the-scaffold) |
| `rule-paths.tsv` | the rule-path register | 6 | `layup setup`; committed by `layup run` in phase 2 (`M2d`) | 1 | [`setup.md`](setup.md#the-rule-path-register) |
| `approvers.tsv` | each human account by its numeric ID, with its role | 3 | `layup run` | 2 | [below](#nfr-001--the-records-of-start) |
| `lease.tsv` | the lease row: run ID, host, start, heartbeat counter | 2 | `layup run` | 2 | [below](#nfr-001--the-records-of-start) |
| `copies.tsv`, `copies/<comment-id>-<seen>.md` | each copied comment: the row (author ID and login, comment ID, App field, time, SHA-256) and the body | 3 | `layup run` | 2 | [below](#nfr-001--the-records-of-start) |
| `intake/psb.tsv` | the table of `layup psb check` on the problem statement | 5 | `layup run` | 2 (`M2c`) | the block `psb-gaps` ([`psb-check.md`](psb-check.md)) |
| `intake/meaning.tsv` | the rows of the review of meaning | 5 | `layup run` | 2 (`M2c`) | [below](#nfr-001--the-records-of-intake) |
| `intake/coverage.tsv` | the items of the completeness review | 7 | `layup run` | 2 (`M2c`) | [below](#nfr-001--the-records-of-intake) |
| `intake/markers.tsv` | the markers of the pinned tree and the proposed values | 5 | `layup run` | 2 (`M2d`) | [below](#nfr-001--the-records-of-the-scaffold) |
| `intake/answers.tsv`, `intake/answers-review.tsv` | the answers of the batch and of its follow-up, one per question (task `T-tmm2`), and the review of the batch's | 5, 7 | `layup run` | 2 (`M2c`) | [below](#nfr-001--the-records-of-intake) |
| `questions.tsv` | every question and its answer, with the phase, the asker and the times | 5, 8, 12 | `layup run` | 2 (`M2c`, K38) | [below](#nfr-001--the-records-of-intake); `M3c` and `M4b` add columns |
| `spec/spans.tsv` | the fact spans of the problem statement | 7 | `layup run` | 2 (`M2c`) | [below](#nfr-001--the-records-of-intake) |
| `spec/check.tsv` | the table of the span check | 7 | `layup run` | 2 (`M2c`) | the block `spec-check`, [below](#nfr-001--the-records-of-intake) |
| `spec/draft.tsv` | the requirement draft of Intake | 7 | `layup run` | 2 (`M2c`) | [below](#nfr-001--the-records-of-intake) |
| `spec/inventory-<n>.tsv` | a version of the confirmed inventory, with its SHA-256 in `spec/inventory.tsv` | 7 | `layup run` | 2 | later |
| `batches.tsv` | each rule batch: head SHA, rule-file hash, approval | 6 | `layup run` | 2 | later |
| `payloads/<sha256>` | a payload by its SHA-256: a known-bad patch, a refused diff, a stall package | 4, 6, 11 | `layup run` | 2 | a rule, no block: [below](#nfr-001--the-records-of-a-session) |
| `tasks.tsv` | the task register | 8 | `layup run` | 2 | later |
| `transitions.tsv` | the transition table | 8 | `layup run` | 2 | later |
| `sessions.tsv` | the session start rows | 4, 9, 12 | `layup run` | 2 (`M2b`) | [below](#nfr-001--the-records-of-a-session) |
| `tasks/<task>/events.tsv` | the events of one task | 3, 8 | `layup run` | 2 (`M2b`, K38) | [below](#nfr-001--the-records-of-a-session) |
| `tasks/<task>/results/<session>.tsv` | the typed result (handoff) of one session, after its check | 4, 8 | `layup run` | 2 (`M2b`) | [below](#nfr-001--the-records-of-a-session) |
| `tasks/<task>/tests.tsv` | the frozen test list | 11 | `layup run` | 2 (`M2e`, K38) | later |
| `tasks/<task>/gates.tsv` | the table of `layup gate` for each head of the task | 6, 8 | `layup run` | 2 | [`gate.md`](gate.md#the-table), with the head as part of the key |
| `acceptance.tsv` | one row per review of a requirement | 8, 12 | `layup run` | 2 | later |
| `parameters.tsv` | every parameter | 10 | `layup run` | 2 | later |
| `harnesses.tsv` | the admitted harnesses, their versions and probe results | 9 | `layup run` | 2 (`M2b`) | [below](#nfr-001--the-records-of-a-session) |
| `owners.tsv` | the confirmed owner map | 9 | `layup run` | 3 | later |
| `routing.tsv` | the routing register and its weights | 9, 13 | `layup run` | 2 (`M2b`, K38) | [below](#nfr-001--the-records-of-a-session), the order only |
| `overrides.tsv` | a harness override of a task (a reroute) | 11 | `layup run` | 3 | later |
| `providers.tsv` | the admitted smart-if providers and their probe results | 10 | `layup run` | 3 | later |
| `decisions.tsv` | one row per smart-if call | 4, 10 | `layup run` | 3 | later |
| `screens.tsv` | one row per run of the escalation screen | 10 | `layup run` | 3 | later |
| `candidates.tsv` | the selected candidates | 10 | `layup run` | 3 | later |
| `escalations.tsv` | the idea owner's escalation decisions | 10 | `layup run` | 3 | later |
| `budget.tsv` | the band `B`, `U`, the appetite, accepted bounds | 10, 12 | `layup run` | 2 (`M2c`, K38) | [below](#nfr-001--the-records-of-intake); `M3b` adds the accepted bounds |
| `milestones.tsv` | each milestone: its bet, cap, clock, and a "milestone stopped" row | 11, 12 | `layup run` | 3 | later |
| `telemetry.tsv` | one row per session | 12 | `layup run` | 1, schema only | [below](#req-011--the-telemetry-record) |
| `stalls.tsv` | the stall, diagnosis and outcome rows | 11 | `layup run` | 1, schema only | [below](#req-009--the-stall-record) |
| `human-inputs.tsv` | every human action | 12 | `layup run` | 3 | later |
| `audit.tsv` | each audited item, with the seed and the population hash | 12 | `layup run` | 3 | later |
| `baseline.tsv`, `start-values.tsv` | the pilot's baseline and start values | 5, 12 | `layup run` | 4 | later |
| `learn/<retrospective>.tsv` | the reward table and the proposed weights | 13 | `layup learn`, then `layup run` | 4 | later |
| `lessons.tsv` | each lesson, with its records and scope | 13 | `layup run` | 4 | later |

A "failed" step (§11, "Not stalls") is a row of `tasks/<task>/events.tsv`, not
of `stalls.tsv`, so the Stall Rate does not count it.

### In other places

| Location | Record | § | Writer | Phase | Schema |
| -------- | ------ | - | ------ | ----- | ------ |
| `target:docs/gates.tsv` | the gate manifest | 6 | `layup setup` | 1 | [`gate.md`](gate.md#the-gate-manifest) |
| `target:docs/setup/armature.pin` | the target's pin file | 5 | `layup setup` | 1 | [`setup.md`](setup.md#nfr-006--the-baseline-at-a-pinned-recorded-version) |
| `target:docs/facts/` | the problem statement and the answers as raw facts | 5, 7 | `layup setup` | 1 | the baseline's facts convention |
| `target:docs/setup/facts.sha256` | the hash of each raw facts file | 5 | `layup setup` | 1 | [`setup.md`](setup.md#the-files-of-docssetup-in-a-target) |
| `target:docs/setup/open-gaps.tsv` | each kept marker with its question | 5 | `layup setup` | 1 | [`setup.md`](setup.md#the-files-of-docssetup-in-a-target) |
| `layup:docs/setup/steps.tsv` | LAYUP's own setup steps; the engine does not read it, and the step table of `setup.md` derives from it | 5 | a LAYUP task | — | its header |
| `layup:internal/catalog/<stack>/` | the stack catalog | 6 | a LAYUP task | 1 | [`setup.md`](setup.md#the-stack-catalog) |
| `host:<work>/` | the work area of `layup setup` | 5 | `layup setup` | 1 | [`setup.md`](setup.md#the-command-layup-setup) |
| `host:registers/harnesses.tsv` | the harness register: the columns that Start reads, and those that `M2b` adds (K40) | 9 | the Operator | 2 | [below](#nfr-001--the-records-of-start); the columns of `M2b` [below](#nfr-001--the-records-of-a-session) |
| `host:registers/models.tsv` | each model of a harness, with its context size and its source, and the "not used" list | 9 | the Operator | 2 (`M2b`) | [below](#nfr-001--the-records-of-a-session) |
| `host:registers/routing.tsv` | the Operator's routing register, which `layup run` copies into `records:routing.tsv` | 9 | the Operator | 2 (`M2b`) | [below](#nfr-001--the-records-of-a-session) |
| `host:registers/forge.tsv` | the forge register: the App, its key file, the API (K40) | 1, 3 | the Operator | 2 | [below](#nfr-001--the-records-of-start) |
| `host:registers/providers.tsv` | the provider register | 10 | the Operator | 3 | later |
| `host:prices.tsv` | the price list | 12 | the Operator | 1, schema only | [below](#req-011--the-telemetry-record) |
| the control repository | the dead-man job's last run time | 11 | the dead-man job | 3 | later |

## NFR-001 — Git is the system of record

Requirement: "Git is the system of record: no project state and no decision is
kept only outside the project repository." Derives from `architecture.md` §3,
ADR-0014 and ADR-0011 decision 2.

1. **From Start on**, `layup run` is the one writer of the records branch
   (ADR-0014), and each producer hands it a typed result (§3). That is phase 2
   and later; [the records of Start](#nfr-001--the-records-of-start) and
   [`run.md`](run.md) give the part of milestone `M2a`.
2. **In phase 1** (O-115, comment 5928163366 on #74): `layup setup` writes the
   setup record, the rule-path register and the records branch's `README.md`
   as files, and makes the first commit of `layup-records` in its work area;
   the Operator pushes it with one printed command, as the Operator pushes the
   root commit ([`setup.md`](setup.md#where-the-records-go-in-phase-1)). So the
   setup record and the evidence of every value are in the target's Git.
3. **The answers** of phase 1 are in the target's tree as two raw fact records (S04, S11; O-124), and
   each answer row names the comment that holds it (`source`). The copy of the
   comments themselves is `layup run`'s (phase 2).
4. **No other state.** The work area of `layup setup` on the host is rebuilt
   from the inputs and the records; nothing in it is the only copy of a value
   once the Operator ran `commands.sh`.
5. **The test of phase 1** (task `T-dep6`, #93, D5 of its plan): on a work
   area that `layup setup` made, `layup-records` has no common commit with
   `main`; its `README.md` is the fixed text of
   [`setup.md`](setup.md#the-readme-of-the-records-branch); its
   `setup/record.tsv`, `setup/verify.tsv` and `rule-paths.tsv` are valid by
   their schemas and equal the files of `WORK/out/`; each answer of
   `answers.tsv` is a fact of a raw answers record on `layup-setup`, with its
   `by` and its `source`; and `commands.sh` holds the four commands in the
   order of [Where the records go in phase 1](setup.md#where-the-records-go-in-phase-1),
   item 4. The audit of a pilot task (`PRD-0001` §7.1) is a measure of
   a later phase.

**Not in phase 1:** the one-writer rule as a ruleset (only the LAYUP App
updates `layup-records`), applied by the Operator at Scaffold 6 (`M2d`,
[`scaffold.md`](scaffold.md#the-ruleset-of-the-records-branch); ADR-0014
decision 6, ADR-0017 decision 3); and `layup audit`, which
then allows the Operator's phase-1 commit of the records branch by its SHA, as
it allows the setup commits (§3).

## NFR-002 — A target is independent of LAYUP

Requirement: "A target repository is independent: it passes its gates without
LAYUP, and a human or a different agent continues the work without the
automation that created it." Derives from `architecture.md` §1 and §6, ADR-0013
and ADR-0016.

1. Only three things go into a target (§1): the setup output, the work of the
   role sessions through pull requests, and the records on `layup-records`. In
   phase 1, only the first and the third.
2. The setup writes no file that the target needs LAYUP to build, test or pass
   its gates: no `layup` binary, no LAYUP script, no LAYUP CI job
   ([`setup.md`](setup.md#the-steps), S12). The gate jobs run the commands of
   the target's own `docs/gates.tsv`. **Decided here** (task `T-c06a`, #91,
   condition 1 of its plan review): the workflow of the gate jobs and their job
   script, which the setup writes from the stack catalog
   ([`setup.md`](setup.md#the-stack-catalog)), are the target's own files: they
   start no LAYUP program, fetch no LAYUP file, and run with LAYUP absent.
   "LAYUP script" and "LAYUP CI job" name LAYUP's own `setup-check.sh` and its
   job (ADR-0016 decision 5). **Decided here** (task `T-dep6`, #93, D6 of its
   plan with condition 1 of its plan review), the mechanical test of this item
   on the head of `layup-setup`: no tracked file is named `layup` or
   `setup-check.sh`, or starts with the magic bytes of an ELF, Mach-O or PE
   program; in each tracked file under `.github/`, on each line with the text
   from its first `#` removed, neither of the words `layup` and `setup-check`
   appears, where a word is delimited by the start or the end of the line or a
   character that is not a letter, a digit, `-` or `_`, so `./layup gate` and a
   URL of LAYUP's repository are refused, and `layup-records` is not; and the
   two JSON files of S13 name no context that starts with `layup/`. In a plain
   clone, with no `layup` on `PATH`, the job script of each kind (the run line
   of its job) gives `pass` or `clear` on the setup head; on a commit of a Go
   file, each active kind gives `pass`, and each pending kind `fail`, `pending:
   product path changed` (§6: a pending kind waits for the first bet). The other
   half of the requirement, another harness that continues a task
   (`PRD-0001` §7.1), is a measure of a later phase.
3. In phase 1, the default branch's ruleset requires only the target's own
   jobs, not the `layup/` checks, which no phase-1 command posts
   ([`setup.md`](setup.md#the-steps), S13). So a target set up in phase 1 merges
   with LAYUP absent from the start.
   In phase 2 the ruleset also requires the four `layup/` checks, and the
   setup record tells the Operator how to remove them when LAYUP is absent
   ([`scaffold.md`](scaffold.md#nfr-002--a-target-without-layup)).
4. Every record is a table or a Markdown file that a human reads with no tool;
   a plain `git clone` carries `origin/layup-records`.

## NFR-001 — The records of Start

Milestone `M2a` (task `T-zck8`, #123): the records that `layup run` writes at
Start, and the two host registers that it reads. The steps that write them are
in [`run.md`](run.md). Each block has one owner, which holds its Go schema and reads or writes the
record with it ([`packages.md`](packages.md#the-table-of-m2a)):
`internal/records` owns `start`, `approvers`, `lease` and `copies`;
`internal/route` owns `harness-register`; `internal/forge` owns
`forge-register`; `internal/run` owns `run-steps` ([`run.md`](run.md#the-command)).

**`start/start.tsv`** holds one row per value of the Start (**decided here**: a
table of names, not one wide row, so that a value of a later milestone is a new
name, not a new column). The names are `layup.version`, `psb.sha256`,
`vision.sha256`, `forge.plan`, `forge.visibility`, `app.permissions`,
`operator.id`, `idea-owner.id`, `intake.cap`, `lease.H`, `watch.T`,
`harness.<id>.cap` and `harness.<id>.wall` for each row of the harness
register at Start, `pin.source`, `pin.commit`, `pin.tree`, `pin.time`, `issue.intake`,
`issue.control` and `watch`. The rows are in this order, and the harness rows
in the order of the register at Start. **Decided here** (task `T-s7vr`, #199):
the harnesses of `start.tsv` are those that its rows `harness.<id>.cap` name,
in their order, each `<id>` of the form `<word>`, and each followed by its row
`harness.<id>.wall`; a restart reads them so, whatever the register is now.

```tsv-schema start records:start/start.tsv
name text key one of the names above, and no other
value text - the value as text: a flag's value as given (`forge.plan`, `intake.cap` as `MONEY,HOURS`); a time in the form of the type `time`; a SHA in the form of `sha1` or `sha256`; a number in the form of `int` or `decimal`; `app.permissions` as `name:level` pairs in the form of `list(text)`; `issue.*` an `int` or `opening`; `watch` `confirmed` or `not-confirmed`; `—` only for `vision.sha256` with no vision brief, `harness.<id>.cap` of a harness with no spend cap, and `issue.*` and `watch` before their steps
source enum(command|register|forge|run) - `command` for a flag's value and for `psb.sha256` and `vision.sha256` (the hash of a flag's file); `register` for `harness.<id>.*`; `forge` for `forge.visibility`, `app.permissions`, `operator.id`, `idea-owner.id` and `issue.*`; `run` for `layup.version`, `pin.*` and `watch`
```

```tsv-schema approvers records:approvers.tsv
id int key the forge's numeric user ID; a login is never a key, as a login can change
role enum(operator|idea-owner|approver) key the role; the Operator and the idea owner can be one account, with one row per role
login text - the login when the row was written
since time - when the row was written
source text - `start` for the rows of the Start command; else the ID of the comment that named the approver (`M2e`, intake.md, Not in M2c)
```

```tsv-schema lease records:lease.tsv
run text key the run ID: 16 lowercase hexadecimal characters, random, made when the run starts
host text - the host name of the run
version text - the LAYUP version of the run
started time - when the run took the lease
heartbeat int - the counter; the run adds 1 every `lease.H`
state enum(held|released) - `released` when the run ended and gave the lease back
```

The lease table has exactly one row; Git history is its log of takeovers.

```tsv-schema copies records:copies.tsv
comment int key the forge's comment ID
seen int key 1 for the first copy; each later read that finds an edit adds a row with the next number, and the first copy stays
issue int - the issue number of the comment
author_id int - the author's numeric ID
author_login text - the author's login at the copy
app text - the App field: the slug of the App that made the comment, or `—` when no App made it
created time - the time of the comment, or of its last edit, as the forge gives it
copied time - when `layup run` copied it
sha256 sha256 - the SHA-256 of the body
body path - the file of the body on the records branch: `copies/<comment>-<seen>.md`
```

```tsv-schema harness-register host:registers/harnesses.tsv
harness id(<word>) key the harness ID, for example `claude`
cap decimal - the spend cap of one session in US dollars; `—` for a harness with no token count or no spend cap
wall int - the wall-clock limit of one session in minutes, 1 or more; never `—`
command text - the command line, words separated by one space; `{model}`, `{cap}` and `{prompt}` stand for the model, the cap and the prompt; `{cap}` exactly when `cap` is not `—`
prompt enum(file|arg|stdin) - how the prompt reaches the harness (`session.md`, The start of a session)
version text - the command that prints the harness's version on the first line of its standard output, words separated by one space
credential text - the absolute path of the credential file, mode 0600; `—` for none
credential_to text - `var:NAME` or `file:PATH`; `—` exactly when `credential` is `—`
rules list(text) - each rule-file name that the harness loads, for example `CLAUDE.md`
policy list(text) - each absolute path of a system-wide policy file of the harness; `—` for none
usage enum(claude-result|none) - the format of its usage report (`session.md`, The usage report of a harness)
billing enum(api|subscription) - the billing type of its account
vars list(text) - each fixed variable of its sessions, `NAME=VALUE`; `—` for none
```

```tsv-schema forge-register host:registers/forge.tsv
forge enum(github) key the forge; GitHub is the only adapter (L-A2)
app_id int - the numeric ID of the LAYUP App
app_slug text - the App's slug, for example `layup-agent`; its bot is `<slug>[bot]`
key_file text - the absolute path of the App's private key (PEM) on the host; mode 0600, owned by the user that runs `layup run`
watch_slug text - the slug of the dead-man job's App, for example `layup-watch`; `—` when there is none
api text - the base URL of the API: `https://api.github.com`; a test names its loopback server
web text - the base URL of the Git remotes: `https://github.com`; a test names a local directory URL
```

**Decided here** (K40, task `T-zck8`): the App ID, the key file and the URLs are
a register of the host, which `--host DIR` names, not flags and not an
environment variable. Reason: `README.md` reads no environment variable for an
input, and a restart reads the same row, so two runs on one target cannot take
different values. The API and Git URLs are columns so that the end-to-end test
runs against a local fake forge with no secret. The harness credential is not
in `M2a`: the architecture gives a session "a variable, or a file that the
register row names" (§4), and `README.md` forbids a variable as an input of
`layup`; `M2b` settles it, with the columns it adds to the harness register
(task `T-ywk7`: [the credential](session.md#the-environment-and-the-harness-credential)).

## NFR-001 — The records of a session

Milestone `M2b` (task `T-ywk7`, #147): the records that `layup run` writes for a
session and a probe, and the host registers that it reads. The steps that write
them are in [`session.md`](session.md). Each block has one owner, which holds its
Go schema ([`packages.md`](packages.md#the-table-of-m2b)): `internal/records`
owns `sessions`, `harnesses`, `routing`, `events` and `result`;
`internal/route` owns `models` and `routing-register`; `internal/session` owns
`probe-result`. One writer (§3): a session hands over its result file and its
commits, and only `layup run` commits a record of it, from its own clone.

```tsv-schema sessions records:sessions.tsv
session id(S-xxxxxxxx) key the session ID, as in `telemetry.tsv`
task id(T-xxxx) - the task; for a probe, the probe's own task ID
attempt int - the attempt, 1 or more; 1 for a probe
role text - the role of the step table (§9), or `probe`
harness id(<word>) - the harness ID of the register
version text - the version that the harness's version command printed before the start
model text - the model ID as the harness names it
base sha1 - the commit that the session's clone starts at
records sha1 - the records commit that the prompt was built from
prompt_bytes int - the size of `prompt.md` in bytes
prompt_tokens int - the estimate: `prompt_bytes` over four, rounded up
context int - the model's context size from `models.tsv`; never below `prompt_tokens`
cap decimal - the spend cap of the session in US dollars; `—` for a harness with none
wall int - the wall-clock limit of the session in minutes, 1 or more
vars list(text) - each fixed variable given to the session, `NAME=VALUE`, `NAME` a letter or `_` and then letters, digits and `_`; `—` for none
policy list(text) - each policy path of the register row that exists, absolute; `—` for none
```

```tsv-schema harnesses records:harnesses.tsv
session id(S-xxxxxxxx) key the probe's session ID
harness id(<word>) - the harness ID of the register
version text - the version that the probe ran; `—` when the version check refused its start
model text - the model of the probe
result enum(passed|failed) - the probe's result
reason text - why it failed: the reason of a refused start (`version`, `rules`, `context`, `prompt`), and after one space its value where it has one (the path of a rule file, the two numbers of the context, the size of the prompt); the class of its end other than `done`; `outside`, `not-used`, `token` or `files`; `—` when `passed`
files list(text) - each instruction file that the harness reported, as it wrote it; `—` for none
models list(text) - each model of its usage report; `—` when the report names none
end time - when the probe ended, or when its start was refused
```

```tsv-schema routing-register host:registers/routing.tsv
role text key the role of the step table (§9)
tier enum(reasoning|execution) key the tier of a task (§9)
position int key the place of the pair in the role's list, from 1
harness id(<word>) - the harness ID of the register
model text - a model of that harness in `models.tsv`
```

```tsv-schema routing records:routing.tsv
role text key the role of the step table (§9)
tier enum(reasoning|execution) key the tier of a task (§9)
position int key the place of the pair in the role's list, from 1
harness id(<word>) - the harness ID of the register
model text - a model of that harness in `models.tsv`
```

```tsv-schema models host:registers/models.tsv
harness id(<word>) key the harness ID of the register
model text key the model ID as the harness names it
context int - the model's context size in tokens
source text - the `http` or `https` URL that states the size
date time - when the source was read
use enum(yes|no) - `no` puts the model on the host's "not used" list
reason text - why `no`, with its source; `—` when `yes`
```

```tsv-schema events records:tasks/<task>/events.tsv
n int key the event's number, from 1, in the order of the events
kind enum(attempt|session|refused|result|push|bound|closed|rebased) - the event
attempt int - the attempt that the event concerns, 1 or more
session id(S-xxxxxxxx) - the session; `—` for `attempt`, `closed` and `rebased`, and for a `refused` whose reason is `pair`, which comes before a session
base sha1 - the base commit, for `attempt` and `rebased`; `—` otherwise
sha sha1 - the commit, for `push` and `bound`; `—` otherwise
detail text - one line: for `result`, the class of the end; for `refused`, the reason, and after one space its value where it has one: the path of a rule file, the two numbers of the context, the size of the prompt, or the SHA-256 of the payload of a refused diff (`payloads/<sha256>`, task `T-z027`); for `push` and `bound`, the branch `task/<task>/<attempt>`, with the row's attempt; `—` for `attempt` and `session`; for `closed` and `rebased`, as `M2e` sets
time time - when `layup run` recorded the event
```

```tsv-schema result records:tasks/<task>/results/<session>.tsv
kind enum(status|artifact) key the kind of the row; `M2e` adds the kinds of a handoff
n int key the row's number within its kind, from 1
value text - for `status`: `completed`, `blocked`, `needs_context`, `decision_needed` or `failed`; for `artifact`: its path in the target
sha256 sha256 - for `artifact`: the SHA-256 of the file at the session's head; `—` for `status`
reason text - for `status`: one line; `—` for `artifact`
```

```tsv-schema probe-result host:sessions/<session>/result/probe.tsv
kind enum(token|file) key `token` at most once; `file` for each instruction file
value text key the token, or the path of the file as the harness wrote it
```

**Decided here**, the row rules that the blocks give in words, which each owner
checks after the types, an error naming the line and the column (task `T-tmhw`'s
form): a column whose rule has no clause for `—` never holds `—`; the kinds
of `events`, each with its writer: `attempt`, the call that starts an attempt
(the uat test of `M2b`, the task loop of `M2e`); `session`, the start of a
session; `result`, its end; `refused`, a start, a result or a push of a task
session that `layup run` refuses (a probe's refused start is its row of
`harnesses.tsv`); `push` and `bound`, the push of a session's SHA; `closed` and
`rebased`, the task loop of `M2e`; each row is one event, and no row is changed
(§3); a result holds one `status` row, with `n` 1, and an
`artifact` value of the type `path`; a `probe-result` holds at most one
`token` row (**decided here**, task `T-bpxg`, condition 2 of the plan review of
#162: a `probe.tsv` with none passes the block, and the probe fails with
`token`, [`session.md`](session.md#the-probe)); `harnesses.tsv` holds one row per probe,
and a harness is admitted by its last row at the version that the version check
reads ([`session.md`](session.md#admission)).

**The columns that `M2b` adds to `host:registers/harnesses.tsv`** (K40) follow
`wall` in the block `harness-register`, with their Go schema in
`internal/route` (task `T-ysph`, row 30a of [the
plan](../plan/README.md#the-tasks-of-m2b); condition 1 of the plan review of
#147). **Decided here** (task `T-ysph`), the rules of the host registers of
`M2b` that their blocks give in words, which `internal/route` checks after the
types, an error naming the line and the column: in `harness-register`,
`command`, `prompt`, `version`, `rules`, `usage` and `billing` never hold `—`;
the rules of `credential_to` and `vars` are those of
[`session.md`](session.md#the-environment-and-the-harness-credential) (a
`NAME` of the form `[A-Za-z_][A-Za-z0-9_]*`, not of the named list, not of a
forge credential or an agent socket; a `PATH` not empty, relative, with no part `..`); in
`models`, no column but `reason` holds `—`, `reason` is `—` exactly when `use`
is `yes`, and `source` is an `http` or `https` URL; in `routing-register`, no
column holds `—`, and the positions of each role and tier are 1 to k, each
once, in any order of the file, as in `routing`. Across the files, which
`layup run` checks when it reads them: each row of `models.tsv` and
`routing.tsv` names a harness of the register, and each model of `routing.tsv`
has a row of its harness in `models.tsv`.

**`payloads/<sha256>`** has a rule, not a block, as a payload is not a table:
its name is the SHA-256 of its bytes, in lowercase hexadecimal; `M2b` writes one
kind, the refused diff of a session (`git diff --binary` from the base to the
head), which an event `refused` names in its `detail`.

**The session directory** (`host:sessions/<session>/`) is not a record: it is
thrown away after the session's records are pushed
([`session.md`](session.md#the-session-directory)).

## NFR-001 — The records of Intake

The records that Intake writes (`M2c`, [`intake.md`](intake.md)); task `T-g49f`
(#206) specified them. `internal/records` owns eight of their blocks, which row
40 (task `T-529i`) compares with their Go schemas; `internal/spec` owns
`spec-check` and holds a second Go value of `spans`, which row 41a (task
`T-4ywt`) compares. A session of Intake writes the file of
its block in `result/` without the columns that code adds (`id`, `start`, `end`,
`text`, `session`), and `layup run` commits the rows with them.
`records:intake/psb.tsv` is the table of `layup psb check`, byte for byte, in
the form of the block `psb-gaps` ([`psb-check.md`](psb-check.md)); it has no
block of its own.

```tsv-schema meaning records:intake/meaning.tsv
id        id(M-NNN)  key  the gap's ID, from M-001, in the order of the rows of the session's file
kind      enum(two-readings|unmeasurable|missing-intent|new-work|fit|no-design|decision)  -  the kind of the gap of meaning (intake.md, The review of meaning)
start     int        -    the byte offset of the quote in start/problem-statement.md, from 0, found by code
end       int        -    the byte offset after the quote; the quote is the file's bytes from start to end
quote     text       -    the quote, byte for byte, one line of the file or part of one, with no tab; it occurs once in the file
question  text       -    the question for the idea owner
session   id(S-xxxxxxxx)  -  the session that wrote the row
```

```tsv-schema spans records:spec/spans.tsv
start   int   key  the byte offset of the span, from 0
end     int   -    the byte offset after the span; start below end
fact    int   -    the fact's number, 1 to n in the order of the spans; `—` for a span that is not a fact
class   enum(need|constraint|measure|context)  -  the fact's class; `—` exactly when fact is `—`
reason  text  -    why the span is not a fact (a heading, a rule line); `—` exactly when fact is not `—`
text    text  -    the span's text as code reads it from the file, with each tab, CR and LF written as one space; written by code, never by the session
```

```tsv-schema spec-check stdout
rule    enum(cover|bounds|text|ids|classes)  key  the rule of the span check (intake.md, The span check), in this order
result  enum(pass|fail)  -  pass, or fail
detail  text  -  for fail, the first broken place: a byte offset or a span's start, and what breaks; `—` for pass
```

`layup run` commits the table of the span check byte for byte as
`records:spec/check.tsv`.

```tsv-schema draft records:spec/draft.tsv
id         id(R-NNN)  key  the requirement's or success criterion's ID in the draft, from R-001, in the order of the rows
kind       enum(requirement|non-functional|success-criterion)  -  a requirement (a need), a non-functional requirement (a constraint), or a success criterion (a measure)
statement  text       -    the statement
covers     list(int)  -    the numbers of the facts that it covers, each a fact of spec/spans.tsv
criterion  text       -    the acceptance criterion of a requirement, or the measure of a success criterion
session    id(S-xxxxxxxx)  -  the session that wrote the row
```

```tsv-schema coverage records:intake/coverage.tsv
id           id(C-NNN)  key  the item's ID, from C-001, in the order of the rows
kind         enum(uncovered|unsupported|class)  -  a fact that nothing covers, a requirement that no fact supports, or a fact whose class the review doubts
fact         int        -    the fact's number; `—` for `unsupported`
requirement  id(R-NNN)  -    the requirement of the draft; `—` for `uncovered` and `class`
question     text       -    the question for the idea owner
session      id(S-xxxxxxxx)  -  the session that wrote the row
```

```tsv-schema questions records:questions.tsv
id        id(Q-NNN)  key  the question's ID, from Q-001, in the order of the merge; never renumbered
source    enum(rule|meaning|coverage|intent|follow-up|setup)  -  where the question comes from; setup is the Operator's question of M2d (scaffold.md, The answers of the setup; task `T-4c5c`)
ref       text       -    its source row: the ID of the rule gap in intake/psb.tsv (a Q-NNN of that table, which the batch's ID replaces), M-NNN, C-NNN, the field of the intent form (band.B, band.U, appetite.hours), the Q-NNN of the batch that a follow-up asks again, or for a setup question (M2d) its setup ID, S01-stack or a marker's M-<x8> (M- and 8 lowercase hexadecimal characters, the id of its row of intake/markers.tsv)
block     enum(idea-owner|operator)  -  the block of the batch that asks it; a setup question is asked in the operator's block (task `T-4c5c`)
phase     enum(before-delivery)  -  the phase of the question (Early Question Share, architecture.md §12); a later milestone adds its phases
kind      text       -    the kind of the gap (a rule G1 to G5, a kind of meaning, a kind of coverage, intent, follow-up, or setup)
start     int        -    the byte offset of its quote; `—` for a question with no quote (an absence, the intent form, an `unsupported` requirement of the draft, S01-stack, a marker with no proposal, a follow-up, whose quote is that of the question it asks again, task `T-vwu3`); the quote of a marker with a proposal is the proposal's (task `T-4c5c`)
end       int        -    the byte offset after its quote; `—` exactly when start is `—`
question  text       -    the question as the batch asks it
asked     time       -    the time that the step read from the run's clock before it wrote the rows, which comes before the records commit and the post of the comment that asks it (task `T-8r3u`)
asker     text       -    `layup psb check`, the session that wrote its source row, or `layup run` for the intent form and a follow-up; for a setup question, `layup run` for S01-stack, which has no source row, and the `session` of its row of intake/markers.tsv for a marker's (task `T-4c5c`)
answered  time       -    the time of the copy that holds its answer; `—` until it has one
state     enum(asked|answered|closed|open)  -  asked until an answer; answered when it has one that no review read (each answer to a follow-up); closed or open as the review of the answers says, or open by code for an intent field still unparsed after the follow-up (intake.md, The review of the answers and the follow-up)
```

A later milestone adds the columns of the measures here first: accepted, the
accepting actor (`M3c`, Clarification Turnaround), and a confirmed reversal's
time and auditor (`M4b`, Reversal Rate) (`architecture.md` §12).

```tsv-schema answers records:intake/answers.tsv
question  id(Q-NNN)  key  the question
answer    text       -    the answer line after `Q-NNN: `, byte for byte
comment   int        -    the numeric ID of the comment that holds it, as copies.tsv records it
author    int        -    the numeric ID of its author, an approver in the role of the question's block
time      time       -    the time of the comment, as copies.tsv records it
```

```tsv-schema answers-review records:intake/answers-review.tsv
question  id(Q-NNN)  key  the question
state     enum(closed|open)  -  closed when the answer closes the gap; open otherwise
reason    text       -    why, one line
session   id(S-xxxxxxxx)  -  the session that wrote the row
```

```tsv-schema budget records:budget.tsv
name    text  key  band.B, band.U (the money appetite of the delivery, architecture.md §12) or appetite.hours; a later milestone adds its names here first (the accepted bounds of an unknown part, M3b)
value   decimal  -  the value: USD for the band, hours for the appetite
source  id(Q-NNN)  -  the question whose answer gives it: the question of the intent form, or the follow-up that asks it again
```

**Decided here:** `budget.tsv` is a table of names, as `start.tsv` is, so that a
value of a later milestone is a new name, not a new column (note 6 of the plan
review of #206).

## NFR-001 — The records of the Scaffold

Milestone `M2d` (task `T-msrd`, #265): the records that `layup run` writes for
the setup questions, the setup and the forge check
([`scaffold.md`](scaffold.md)). `internal/records` owns `markers`;
`internal/run` owns `setup-tree` and `forge-check`, which it alone writes and
reads; row 50 (task `T-tetz`) compares the three with their Go schemas.
`setup/record.tsv`, `setup/verify.tsv` and `rule-paths.tsv` keep their
blocks and owners of phase 1; `layup run` commits them
([`scaffold.md`](scaffold.md#the-step-verify-and-the-records-of-the-setup)).

```tsv-schema markers records:intake/markers.tsv
id      id(M-xxxxxxxx)  key  the marker's key of S10: M- and the first 8 hexadecimal characters of the SHA-256 of its file, a tab and its text; in questions.tsv the source setup tells it from a gap of meaning (M-NNN), whose source is meaning; the reader checks its form, `M-` and 8 lowercase hexadecimal characters, and does not compute the hash again, as the row holds the marker after the field rule, which makes a tab or a carriage return one space
file    path   -    the file that holds the marker, in the pinned tree
line    int    -    the line of its open quote, from 1; for a text that the file holds on more lines (one key), the first
marker  text   -    the marker, with its two angle quotes; a marker over more lines has its line breaks written as spaces; never `—`
value   text   -    the proposed value; `—` for no proposal
start   int    -    the byte offset of the proposal's quote in start/problem-statement.md, from 0, found by code; `—` exactly when value is `—`
end     int    -    the byte offset after the quote, above start; `—` exactly when start is `—`
session text   -    the session that proposed the values (S-xxxxxxxx); `—` only for a table with no row, which has none
```

The rows are in the order of S10's stop table, by file and then by line: the
file of a row is not below the file before it, in the bytes of the path, and in
one file the line is not below the line before it, so two markers on one line
are two rows with one line. A table with no row is a scan that found no marker,
and no session started.

```tsv-schema setup-tree records:setup/tree.tsv
order   int   key  0 for the root commit, then 1, 2, … for each setup commit of layup-setup, in order; the rows are in the order of this column, so the first row is the root commit, and a table has at least that row
commit  sha1  -    the commit
tree    sha1  -    its tree
step    text  -    `root` for order 0; else the step of the commit, S04 to S14, from its message `chore: setup <step>`; never `—`
```

The last row is the head of `layup-setup` that `layup setup verify` read and
that the Operator pushes; no column is ever `—`.

```tsv-schema forge-check records:setup/forge-check.tsv
check   text   key  `tree`, `rules:<branch>` for the default branch, `layup-probe` and `layup-records`, `bypass:default`, `bypass:records`, `probe`, `activity`, in that order, one row for each whatever its result, so the table has eight rows; the branch of the default branch is not empty, and is neither of the next two names, which would repeat a key
result  enum(pass|fail)  -  the result by the rule of the check (scaffold.md, REQ-003, Scaffold 7)
detail  text   -    what the check found, one line: for a `fail` the first missing or different rule or the reason; for a `pass` of `rules:<branch>` the extra rule types or `—`; `—` for another `pass`
run     text   -    the run ID that made the check, in the form of lease.tsv (16 lowercase hexadecimal characters); never `—`
at      time   -    when the check ran, from the run's clock; never `—`
```

The table is written again at each start after the Scaffold
([`scaffold.md`](scaffold.md#at-each-later-start)); Git history keeps each.

**The built blocks that `M2d` changes in words** (condition 4 of the plan
review of #265): each change is written beside its rule, and the build task
that writes it changes the block and its Go schema in one change.

- `questions`: the enum of `source` gains `setup`, and `ref` the setup IDs
  ([`scaffold.md`](scaffold.md#the-answers-of-the-setup));
- `run-steps`: the enum of `step` gains `markers`, `setup`, `verify`,
  `scaffold-push` and `forge-check`
  ([`scaffold.md`](scaffold.md#the-steps-of-the-scaffold));
- `setup-record`: the row name `without-layup`
  ([`scaffold.md`](scaffold.md#nfr-002--a-target-without-layup)), and two forms
  of the `ref` of a `computed` row, `records:start/start.tsv` (S01, S02 in
  phase 2) and `out/ruleset-default.json` (`without-layup`), which check
  `sources` accepts as it is (a ref that is not empty);
- `setup-steps`: none; the actor `layup-run` exists.

**Not a block:** `setup/rulesets.txt`, the Operator's two apply outputs, each
after a line `== <file>`, byte for byte (JSON, not a table); and
`setup/inputs/files/<path>`, a prose text at its path in the target, which the
check of its setup row reads; it is committed with the end of the prose session
that wrote it, and removed again, in one commit, when the step `setup` is `fail`
after that session ([`scaffold.md`](scaffold.md#the-stops-of-the-runner)).

## REQ-011 — The telemetry record

Requirement: "Every task has a record of its token count, its latency and its
wall-clock duration in the target." Derives from `architecture.md` §12 (The
ledger), ADR-0024 and ADR-0007.

**Schema only.** One row per session (§12: the unit of an action); a session
is started only by `layup run`, so no phase-1 command writes a row. The writer
is `layup run` (phase 2). A task's record is the set of its rows (the `task`
column); a task's record is complete when each of its rows has `observed`
tokens, a latency and a duration (`F-0003#60`).

```tsv-schema telemetry records:telemetry.tsv
session        id(S-xxxxxxxx)             key  the session ID: `S-` and 8 lowercase hexadecimal characters, random (decided here)
task           id(T-xxxx)                 -    the task ID of the target's scheme
requirements   list(text)                 -    the requirement IDs of the task, from the task register
role           text                       -    the role of the step table (§9)
harness        text                       -    the harness ID of the register
model          text                       -    the model ID as the harness names it
billing        enum(api|subscription)     -    the billing type of the harness's account
start          time                       -    when `layup run` started the harness process
first_output   time                       -    the first output of the process; `—` when it gave none
end            time                       -    when the process exited or was killed; when a child holds an output after that, when the outputs closed (task `T-5pxd`)
latency_s      int                        -    first_output − start, in seconds; `—` when first_output is `—`, which is its status column
duration_s     int                        -    end − start, in seconds
tokens_in      int                        -    input tokens; `—` when not reported
tokens_out     int                        -    output tokens; `—` when not reported
tokens_cache   int                        -    cache read and write tokens; `—` when not reported
tokens_status  enum(observed|partial|unavailable)  -  `unavailable` is the criterion's `not reported`; a row that is not `observed` makes its task incomplete
tokens_reason  text                       -    why the tokens are `partial` or `unavailable`; `—` when `observed`
money          decimal                    -    the cost; `—` when `unknown`, never 0 for an unknown value (FT2)
currency       text                       -    ISO 4217 code, for example `USD`; `—` when `unknown`
money_status   enum(reported|computed|unknown)  -  `reported` by the harness; `computed` from the tokens and a price row; `unknown`
price          text                       -    the `id` of the `prices.tsv` row for `computed`, else `—`; a subscription session is `computed` from the list price
```

```tsv-schema prices host:prices.tsv
id        id(P-NNN)  key  the row ID
harness   text       -    the harness or provider ID
model     text       -    the model ID
class     enum(in|out|cache)  -  the token class
price     decimal    -    the price per million tokens
currency  text       -    ISO 4217 code
source    text       -    the URL of the price list that states it
date      time       -    when the source was read
```

**Decided here:** the session ID's form, the column names, and three token
classes (the harness reports differ; a fourth class is a new column, added
here first). The architecture gives the content (§12) and not the names.

**Decided here** (task `T-tmhw`, #94, D2 and D3 of its plan), the row rules
that the two blocks give in words, which `internal/records` checks after the
types, on each row (`CheckTelemetry`, `CheckPrice`, so a writer checks a row
before it writes it), an error naming the line and the column:

- A column whose block rule or type has no clause for `—` never holds `—`, as
  [`README.md`](README.md#the-schema-block) gives that check to the owner of a
  record (round 1 of #94): in `telemetry.tsv`, `task`, `role`, `harness`,
  `model`, `billing`, `start`, `end`, `duration_s`, `tokens_status` and
  `money_status`; in `prices.tsv`, each column that is not the key (a key is
  never `—` by `internal/tsv`). `requirements` may be `—`, as its type `list`
  writes an empty list so.
- `session` is `S-` and 8 lowercase hexadecimal characters (the type gives
  letters and digits).
- `latency_s` is `—` exactly when `first_output` is `—`; `start ≤ first_output
  ≤ end`; `latency_s` is `first_output − start` and `duration_s` is `end −
  start`, in whole seconds, as the times are to the second. Reason: a reader
  of a later phase sums these columns, and a sum that the times contradict is
  a false record.
- `tokens_reason` is `—` exactly when `tokens_status` is `observed`;
  `observed` has the three token columns, `unavailable` none of them, and
  `partial` at least one and not all three (a report with a part missing).
- `money` and `currency` are `—` exactly when `money_status` is `unknown`, so
  an unknown cost is never 0 (FT2, ADR-0024); a `currency` is three capital
  letters, the form of ISO 4217; `price` is a row ID `P-NNN` exactly when
  `money_status` is `computed`; a `subscription` session is never `reported`:
  it is `computed` from the list price, or `unknown`, for example when it has
  no tokens.
- In `prices.tsv`, `currency` has the same form, and `source` is an `http` or
  `https` URL ("the URL of the price list").

**Known limits:** the list of the codes of ISO 4217 is not checked, only
their form; that a `price` ID is a row of the host's `prices.tsv` is a check
of two files, which the writer makes (phase 2).

**Not in phase 1:** the writer (`M2b`:
[`session.md`](session.md#req-011--the-writer-of-the-telemetry-record)), the
budget check before each session start (`M2c` and `M3d`), and the Telemetry
Completeness report (`layup report`, phase 4).

## REQ-009 — The stall record

Requirement: "Every stall has a record in the target with its diagnosis and its
outcome." Derives from `architecture.md` §11, ADR-0023.

**Schema only.** A stall is found and handled by `layup run` (the triggers and
the procedure of §11, phase 3), so no phase-1 command writes a row. A closed
stall has three rows, in this order: a `stall` row; a `diagnosis` or a
`diagnosis-failed` row; an `outcome` row (§11 procedure steps 2 and 6). An
open stall has the rows of the steps that ended. A new trigger in the same
task is a new stall, with its own ID and rows.

```tsv-schema stalls records:stalls.tsv
stall     id(ST-NNN)    key  the stall ID, from ST-001, in order
kind      enum(stall|diagnosis|diagnosis-failed|outcome)  key  the row kind; a stall has at most one row of each of its three steps
time      time          -    when the row was written
task      text          -    the task ID; `project` for a stall of the orchestrator (trigger 5)
trigger   enum(no-progress|too-many-rounds|hang|no-report|orchestrator)  -  on a `stall` row; `—` on the others
evidence  sha256        -    on a `stall` row: the payload of the package; on a `diagnosis` row: the payload of the diagnosis in its fixed form, with each open unknown and its evidence (§11); `—` on the others (`payloads/<sha256>`)
cause     enum(disagreement|missing-information|wrong-gate|harness-failure|task-too-large|other)  -  on a `diagnosis` row; `—` on the others
rung      enum(retry|panel|operator)  -  on a `diagnosis` row: the rung it recommends; `—` on the others
examiner  text          -    on a `diagnosis` row: the examiner's session ID; on a `diagnosis-failed` row: the examiner's session ID, or `—` when there was none; `—` on the others
outcome   enum(closed-without-human|closed-by-operator|task-stopped)  -  on an `outcome` row; `—` on the others
note      text          -    one line; on a `diagnosis-failed` row: its reason; on an `outcome` row: the reason of the outcome; `—` on the others
```

Stall Diagnosis (`F-0003#61`) counts the stalls whose second row is
`diagnosis-failed` or missing; it must be zero (the criterion of `REQ-009`).

**Decided here:** one file for the three row kinds, so that a stall's rows are read together and Stall Diagnosis is one count over one file; and the column names, which §11 does not give.

**Decided here** (task `T-dgy7`, #95, D2 and D3 of its plan, with the
conditions and the notes of its plan review), the rules that the block gives
in words, which `internal/records` checks after the types: `CheckStall` on
each row, so that a writer checks a row before it writes it, and `CheckStalls`
on the order of the rows; `ReadStalls` reads a file with both. An error names
the line and the column.

- `time` and `task` never hold `—`, as their block rules have no clause for
  it (the lesson of row 17, `guardrails.md` §2).
- Each kind has the columns that the block names for it, and `—` in the
  others: a `stall` row has `trigger` and `evidence`; a `diagnosis` row has
  `evidence`, `cause`, `rung` and `examiner`; a `diagnosis-failed` row has
  `note`, and `examiner` or `—`; an `outcome` row has `outcome` and `note`.
  Reason for `examiner`: an examiner session writes the diagnosis (§11
  procedure step 2), so "`—` when there was none" is a `diagnosis-failed` row,
  for example with no admitted harness.
- `examiner`, when it is not `—`, is a session ID: `S-` and 8 lowercase
  hexadecimal characters, the form of `session` in `telemetry.tsv`.
- On a `stall` row, `task` is `project` exactly when `trigger` is
  `orchestrator` (trigger 5); each other row of a stall names the `task` of its
  `stall` row, so all the rows of an orchestrator stall name `project`.
- The `stall` rows have the IDs `ST-001`, `ST-002`, … in the order of the
  file, with no gap ("from ST-001, in order"). Each other row comes after the
  `stall` row of its ID; a stall has one `diagnosis` or `diagnosis-failed` row,
  not both; its `outcome` row comes after that row (§11: step 6 follows step
  2); and each row is at or after the time of the row before it of the same
  stall. The time order is per stall, as the block gives no rule across
  stalls. The key `(stall, kind)` already refuses two rows of one kind for one
  stall (`internal/tsv`), so the order check adds only what the key does not
  hold: both second kinds, the order, the IDs with no gap, the task and the
  time.
- A stall with only its `stall` row, or with no `outcome` row, is valid in the
  file: it is open, as the writer appends a row when its step ends (§11). That
  a second row is missing is what Stall Diagnosis counts, a measure of the
  report and not a rule of the file.

**Known limits:** the writer (phase 3) holds the rules of §11 that no block
sentence gives: after a `diagnosis-failed` row the package goes to the
Operator at once, so that outcome is never `closed-without-human`; and a stall
of the orchestrator has only the diagnosis and the Operator, so its outcome is
never `closed-without-human` and its `rung` is never `retry` or `panel`.
`task` is `project` or the target's own task ID, whose scheme the validator
leaves open. That an `evidence` hash names a file of `payloads/` is a check of
two places, which the writer makes.

**Not in phase 1:** the writer, the triggers, the examiner, the panel and the
Operator's answer form (§11; `REQ-010`, phase 3).
