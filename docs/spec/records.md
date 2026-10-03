# The records

Where each record of LAYUP lives, who writes it, and the schemas of the
records that phase 1 needs. Conventions: [`README.md`](README.md). The records
of a target are on its branch `layup-records`, an orphan branch that only
`layup run` writes ([`architecture.md`](../architecture.md) §3, ADR-0014).

## The layout of the records branch

One row per record kind of the whole architecture (§1 to §15), so that a later
milestone adds schemas and does not move a file. The inventory that this table
was checked against is [`runs/T-0drh/checklist.md`](../../runs/T-0drh/checklist.md),
part C. "Phase" is the earliest `PRD-0001` phase whose requirement needs the
record; a schema is in this directory only for a record of phase 1.

**Decided here:** every path that the architecture does not name (it names 18
files, most without a directory), and the phase of each record. The paths
follow two rules: a table that the whole target shares is at the root; a table
or payload of one task is under `tasks/<task>/`.

### On the records branch (`records:`)

| Path | Record | § | Writer | Phase | Schema |
| ---- | ------ | - | ------ | ----- | ------ |
| `README.md` | what the branch is, and how to read it with no tool | 3 | `layup setup` (O-115) | 1 | a fixed text: [`setup.md`](setup.md#the-readme-of-the-records-branch) |
| `start/problem-statement.md`, `start/vision.md` | the two briefs, byte for byte | 5 | `layup run` | 2 | — |
| `start/start.tsv` | each brief's SHA-256, the LAYUP version, the forge plan, the values of the Start command | 5, 13 | `layup run` | 2 | later |
| `setup/record.tsv` | the setup record: each value with its source, the pin rows included | 5 | `layup setup` | 1 | [`setup.md`](setup.md#the-setup-record) |
| `setup/verify.tsv` | the table of `layup setup verify` | 5 | `layup setup verify`, committed by `layup setup` (O-115) | 1 | [`setup.md`](setup.md#the-checks-of-layup-setup-verify) |
| `setup/rulesets.txt` | the output of the Operator's ruleset command (the bypass list) | 3 | `layup run` | 2 | — |
| `setup/forge-check.tsv` | the effective rules read back, and the probe results | 3, 5 | `layup run` | 2 | later |
| `rule-paths.tsv` | the rule-path register | 6 | `layup setup` | 1 | [`setup.md`](setup.md#the-rule-path-register) |
| `approvers.tsv` | each human account by its numeric ID, with its role | 3 | `layup run` | 2 | later |
| `lease.tsv` | the lease row: run ID, host, start, heartbeat counter | 2 | `layup run` | 2 | later |
| `copies.tsv`, `copies/<comment-id>.md` | each copied comment: the row (author ID and login, comment ID, App field, time, SHA-256) and the body | 3 | `layup run` | 2 | later |
| `intake/meaning.tsv` | the rows of the review of meaning | 5 | `layup run` | 2 | later |
| `questions.tsv` | every question and its answer, with the phase, the asker and the times | 5, 8, 12 | `layup run` | 3 | later |
| `spec/spans.tsv` | the fact spans of the problem statement | 7 | `layup run` | 2 | later |
| `spec/inventory-<n>.tsv` | a version of the confirmed inventory, with its SHA-256 in `spec/inventory.tsv` | 7 | `layup run` | 2 | later |
| `batches.tsv` | each rule batch: head SHA, rule-file hash, approval | 6 | `layup run` | 2 | later |
| `payloads/<sha256>` | a payload by its SHA-256: a known-bad patch, a refused diff, a stall package | 4, 6, 11 | `layup run` | 2 | — |
| `tasks.tsv` | the task register | 8 | `layup run` | 2 | later |
| `transitions.tsv` | the transition table | 8 | `layup run` | 2 | later |
| `sessions.tsv` | the session start rows | 4, 9, 12 | `layup run` | 2 | later |
| `tasks/<task>/events.tsv` | the events of one task | 3, 8 | `layup run` | 2 | later |
| `tasks/<task>/results/<session>.tsv` | the typed result (handoff) of one session, after its check | 4, 8 | `layup run` | 2 | later |
| `tasks/<task>/tests.tsv` | the frozen test list | 11 | `layup run` | 3 | later |
| `tasks/<task>/gates.tsv` | the table of `layup gate` for each head of the task | 6, 8 | `layup run` | 2 | [`gate.md`](gate.md#the-table), with the head as part of the key |
| `acceptance.tsv` | one row per review of a requirement | 8, 12 | `layup run` | 2 | later |
| `parameters.tsv` | every parameter | 10 | `layup run` | 2 | later |
| `harnesses.tsv` | the admitted harnesses, their versions and probe results | 9 | `layup run` | 2 | later |
| `owners.tsv` | the confirmed owner map | 9 | `layup run` | 3 | later |
| `routing.tsv` | the routing register and its weights | 9, 13 | `layup run` | 3 | later |
| `overrides.tsv` | a harness override of a task (a reroute) | 11 | `layup run` | 3 | later |
| `providers.tsv` | the admitted smart-if providers and their probe results | 10 | `layup run` | 3 | later |
| `decisions.tsv` | one row per smart-if call | 4, 10 | `layup run` | 3 | later |
| `screens.tsv` | one row per run of the escalation screen | 10 | `layup run` | 3 | later |
| `candidates.tsv` | the selected candidates | 10 | `layup run` | 3 | later |
| `escalations.tsv` | the idea owner's escalation decisions | 10 | `layup run` | 3 | later |
| `budget.tsv` | the band `B`, `U`, the appetite, accepted bounds | 10, 12 | `layup run` | 3 | later |
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
| `host:registers/harnesses.tsv` | the harness register | 9 | the Operator | 2 | later |
| `host:registers/providers.tsv` | the provider register | 10 | the Operator | 3 | later |
| `host:prices.tsv` | the price list | 12 | the Operator | 1, schema only | [below](#req-011--the-telemetry-record) |
| the control repository | the dead-man job's last run time | 11 | the dead-man job | 3 | later |

## NFR-001 — Git is the system of record

Requirement: "Git is the system of record: no project state and no decision is
kept only outside the project repository." Derives from `architecture.md` §3,
ADR-0014 and ADR-0011 decision 2.

1. **From Start on**, `layup run` is the one writer of the records branch
   (ADR-0014), and each producer hands it a typed result (§3). That is phase 2
   and later.
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
updates `layup-records`), applied at Start (phase 2); and `layup audit`, which
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
4. Every record is a table or a Markdown file that a human reads with no tool;
   a plain `git clone` carries `origin/layup-records`.

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
end            time                       -    when the process exited or was killed
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

**Not in phase 1:** the writer, the budget check before each session start, and
the Telemetry Completeness report (`layup report`, phase 4).

## REQ-009 — The stall record

Requirement: "Every stall has a record in the target with its diagnosis and its
outcome." Derives from `architecture.md` §11, ADR-0023.

**Schema only.** A stall is found and handled by `layup run` (the triggers and
the procedure of §11, phase 3), so no phase-1 command writes a row. Each stall
has three rows, in this order: a `stall` row; a `diagnosis` or a
`diagnosis-failed` row; an `outcome` row (§11 procedure steps 2 and 6). A new
trigger in the same task is a new stall, with its own ID and rows.

```tsv-schema stalls records:stalls.tsv
stall     id(ST-NNN)    key  the stall ID, from ST-001, in order
kind      enum(stall|diagnosis|diagnosis-failed|outcome)  key  the row kind; each stall has one row of each of its three steps
time      time          -    when the row was written
task      text          -    the task ID; `project` for a stall of the orchestrator (trigger 5)
trigger   enum(no-progress|too-many-rounds|hang|no-report|orchestrator)  -  on a `stall` row; `—` on the others
evidence  sha256        -    on a `stall` row: the payload of the package; on a `diagnosis` row: the payload of the diagnosis in its fixed form, with each open unknown and its evidence (§11); `—` on the others (`payloads/<sha256>`)
cause     enum(disagreement|missing-information|wrong-gate|harness-failure|task-too-large|other)  -  on a `diagnosis` row; `—` on the others
rung      enum(retry|panel|operator)  -  on a `diagnosis` row: the rung it recommends; `—` on the others
examiner  text          -    on a `diagnosis` or `diagnosis-failed` row: the examiner's session ID; `—` when there was none
outcome   enum(closed-without-human|closed-by-operator|task-stopped)  -  on an `outcome` row; `—` on the others
note      text          -    one line: the reason of a `diagnosis-failed` row, or of the outcome
```

Stall Diagnosis (`F-0003#61`) counts the stalls whose second row is
`diagnosis-failed` or missing; it must be zero (the criterion of `REQ-009`).

**Decided here:** one file for the three row kinds, so that a stall's rows are read together and Stall Diagnosis is one count over one file; and the column names, which §11 does not give.

**Not in phase 1:** the writer, the triggers, the examiner, the panel and the
Operator's answer form (§11; `REQ-010`, phase 3).
