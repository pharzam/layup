# T-0drh — the checklist of the specification (plan step 1)

Written by the author (Claude Opus 5.5) on 2026-10-01, before any file of
`docs/spec/`. It is the red half of this document task: every row is open, and
step 10 closes each one with the place in `docs/spec/` that answers it. The
plan, the plan review and the author's answer are comments on
[#74](https://github.com/pharzam/layup/issues/74). "AC *n*" is the *n*-th
acceptance criterion of the issue body; "C*n*" and "N*n*" are the conditions and
notes of the plan review. Sections (§) are those of
[`docs/architecture.md`](../../docs/architecture.md) at `648b37f`.

Status values: `open`; `done: <file>#<section>`; `out: <reason>`.

## A. One section per phase-1 requirement (AC 1, AC 3, AC 4, AC 5)

| Row | ID | File | Derives from (§, ADR) | Contract it must give | Status |
| --- | -- | ---- | --------------------- | --------------------- | ------ |
| A1 | `REQ-001` | `psb-check.md` | §5 gap check step 1; ADR-0011 | `layup psb check`: argument, columns of the gap table, exit codes; rules G1 to G5; the review of meaning named as a later phase | done: psb-check.md#req-001 |
| A2 | `REQ-002` | `setup.md` | §5 Start and Scaffold; ADR-0011 decisions 6 and 7, ADR-0016 | `layup setup`: flags, inputs, the step table S01 to S15 (C1), the stop table (C3), exit codes; `layup setup verify`: flags, the check list in phase 1 (C4), result columns, exit codes | done: setup.md#req-002 |
| A3 | `REQ-004` | `gate.md` | §6 The stack gates; ADR-0016 | `layup gate`: flags, the manifest it reads, result columns per kind, verdict values `pass`, `fail`, `not-active`, `clear`, exit codes | done: gate.md#req-004 |
| A4 | `REQ-007` | `gate.md` | §6, §8 task loop step 5; ADR-0016, ADR-0019 | phase 1: the deterministic verdict and the repeat rule; out: the status `layup/gates` and the merge block, with the phase that delivers them (C4) | done: gate.md#req-007 |
| A5 | `REQ-009` | `records.md` | §11; ADR-0023 | schema only (AC 5): the stall row, the diagnosis row (with "diagnosis failed"), the outcome row (N10); key; writer `layup run` | done: records.md#req-009 |
| A6 | `REQ-011` | `records.md` | §12 The ledger; ADR-0024, ADR-0007 | schema only (AC 5): `telemetry.tsv`, key the session ID, a task's record is its rows (N9); the token status `unavailable` is the criterion's `not reported`; writer `layup run` | done: records.md#req-011 |
| A7 | `NFR-001` | `records.md` | §3; ADR-0014, ADR-0011 | where each phase-1 record lives in Git, and who writes it there (C2) | done: records.md#nfr-001 |
| A8 | `NFR-002` | `records.md` | §1, §6; ADR-0013, ADR-0016 | what goes into a target (three kinds only); the target's gates with LAYUP absent | done: records.md#nfr-002 |
| A9 | `NFR-003` | `setup.md` | §5 Scaffold step 1; no ADR, guardrails §1.1 Inv-4 (N7) | each value of the setup record has a source; the source kinds; a value with no source keeps its marker | done: setup.md#nfr-003 |
| A10 | `NFR-004` | `gate.md`, linked from `setup.md` (N6) | §6 table of `layup/gates`; ADR-0016, ADR-0011 | a check that did not run is never 0, in both commands; `clear` is 0 (N4) | done: gate.md#nfr-004 |
| A11 | `NFR-005` | `gate.md` | §4; ADR-0015 | the engine checks open no model connection; `decisions.tsv`: name, place, phase only (N6) | done: gate.md#nfr-005 |
| A12 | `NFR-006` | `setup.md` | §5 Start step 2; ADR-0009 | LAYUP's own pin and the target's pin: the pin file's fields | done: setup.md#nfr-006 |
| A13 | `NFR-007` | `packages.md` | §1; ADR-0010, ADR-0013 | the standard library only, `git` as a program; the package table | done: packages.md#nfr-007 |

## B. The package table (AC 2, AC 6)

| Row | Item | Status |
| --- | ---- | ------ |
| B1 | Each package of phase 1: path, job, may import (the input of the boundary rule of §6) | done: packages.md, the table of phase 1 |
| B2 | Each package of phases 2 to 4, by name, job and phase only (AC 6) | done: packages.md, the packages of later phases |
| B3 | A component row for `layup setup`, marked "decided here" (N2): what it may write | done: packages.md, the components of phase 1 |
| B4 | The stack catalog embedded with `embed` (N3), the package that holds it | done: packages.md (`internal/catalog`); setup.md, the stack catalog |

## C. The record inventory of the whole architecture (C5, AC 4, AC 6)

Drawn from §1 to §15, not only from the sections that the issue names. The 18
files of the `grep` are marked `(file)`; the other kinds have no file name in
the architecture. "Phase" is the `PRD-0001` phase whose requirement first needs
the record; `records.md` sets the final value. "Where": `records` is the target's
branch `layup-records`; `target` is the target's default branch; `host` is the
LAYUP host; `LAYUP` is LAYUP's own repository; `forge` is a forge object.

| Row | Record | § | Where | Writer | Phase | Status |
| --- | ------ | - | ----- | ------ | ----- | ------ |
| C1 | the target's pin record (source, commit, tree, time) | 5 | records | `layup setup` (C1); later `layup run` | 1 | done: records.md, the layout |
| C2 | the problem statement and the vision brief, byte for byte, each with its SHA-256 | 5 | records | `layup run` | 1 (the facts by S06: target) | done: records.md, the layout |
| C3 | `approvers.tsv` (file) | 3, 5 | records | `layup run` | 2 | done: records.md, the layout |
| C4 | the lease row | 2 | records | `layup run` | 2 | done: records.md, the layout |
| C5 | the copy of a comment (body, author ID and login, comment ID, App field, time, SHA-256) | 3 | records | `layup run` | 2 | done: records.md, the layout |
| C6 | the output of the Operator's ruleset command (the bypass list) | 3 | records | `layup run` | 2 | done: records.md, the layout |
| C7 | the rules read back and the probe results | 3, 5 | records | `layup run` | 2 | done: records.md, the layout |
| C8 | the LAYUP version and the plan of the forge (`--plan`) | 5, 13 | records | `layup run` | 2 | done: records.md, the layout |
| C9 | the setup record: each value with its source; the steps with their evidence | 5 | records (C2: the home in phase 1) | `layup setup` | 1 | done: records.md, the layout; schema in setup.md |
| C10 | the gap table of `layup psb check` | 5 | standard output; then a question row | `layup psb check` | 1 | done: psb-check.md, the table |
| C11 | the rows of the review of meaning (quote, kind, question) | 5 | records | `layup run` | 2 | done: records.md, the layout |
| C12 | `questions.tsv` (file): Intake questions, follow-ups, questions during the work, answers, accepted time | 5, 8, 12 | records | `layup run` | 3 | done: records.md, the layout |
| C13 | `docs/gates.tsv` (file): gate kind, command, scope, state | 6 | target | `layup setup` | 1 | done: records.md, the layout; schema in gate.md |
| C14 | the stack catalog entry: per kind, tool, version, command, scope, configuration, known-bad fixture, evidence | 6 | LAYUP | a LAYUP task | 1 | done: records.md, the layout; form in setup.md |
| C15 | the rule-path register | 6 | records (C2) | `layup setup` | 1 written, 2 read | done: records.md, the layout; schema in setup.md |
| C16 | a known-bad patch, as a payload | 6 | records | `layup run` | 2 | done: records.md, the layout |
| C17 | the approval of a rule batch: head SHA, rule-file hash | 6 | records | `layup run` | 2 | done: records.md, the layout |
| C18 | the result table of `layup setup verify` | 5 | standard output; then records | `layup setup verify` | 1 | done: records.md, the layout; schema in setup.md |
| C19 | the result table of `layup gate`, one row per kind | 6 | standard output; then records | `layup gate` | 1 | done: gate.md, the table; records.md `tasks/<task>/gates.tsv` |
| C20 | the fact spans (ID, class, byte offsets; "not a fact" with a reason) | 7 | records | `layup run` | 2 | done: records.md, the layout |
| C21 | a version of the confirmed inventory, with its SHA-256 | 7 | records | `layup run` | 2 | done: records.md, the layout |
| C22 | the copy of a bet; the copy of a decision | 7, 8 | records | `layup run` | 2 | done: records.md, the layout |
| C23 | the diff of a refused change, as a payload | 4 | records | `layup run` | 2 | done: records.md, the layout |
| C24 | the session start row (attempt, base commit, caps, context size estimate) | 4, 9, 12 | records | `layup run` | 2 | done: records.md, the layout |
| C25 | the typed result of a session (the handoff) | 4, 8 | session directory, then records | a role session; checked by `layup run` | 2 | done: records.md, the layout |
| C26 | `tasks/<task>/events.tsv` (file) | 3, 8 | records | `layup run` | 2 | done: records.md, the layout |
| C27 | the task register | 8 | records | `layup run` | 2 | done: records.md, the layout |
| C28 | the transition table | 8 | records | `layup run` | 2 | done: records.md, the layout |
| C29 | the frozen test list (test ID, SHA-256 of its source) | 11 | records | `layup run` | 3 | done: records.md, the layout |
| C30 | the harness register | 9 | host | the Operator | 2 | done: records.md, the layout |
| C31 | the harness probe results | 9 | records | `layup run` | 2 | done: records.md, the layout |
| C32 | the owner map | 9 | records | `layup run` | 3 | done: records.md, the layout |
| C33 | the routing register | 9, 13 | records | `layup run` | 3 | done: records.md, the layout |
| C34 | a harness override of a task (a reroute) | 11 | records | `layup run` | 3 | done: records.md, the layout |
| C35 | the provider register | 10 | host | the Operator | 3 | done: records.md, the layout |
| C36 | the provider probe results | 10 | records | `layup run` | 3 | done: records.md, the layout |
| C37 | `decisions.tsv` (file) | 4, 10 | records | `layup run` | 3 | done: records.md, the layout |
| C38 | `screens.tsv` (file) | 10 | records | `layup run` | 3 | done: records.md, the layout |
| C39 | `candidates.tsv` (file) | 10 | records | `layup run` | 3 | done: records.md, the layout |
| C40 | `escalations.tsv` (file) | 10 | records | `layup run` | 3 | done: records.md, the layout |
| C41 | `parameters.tsv` (file) | 10 | records | `layup run` | 2 | done: records.md, the layout |
| C42 | `budget.tsv` (file) | 10, 12 | records | `layup run` | 3 | done: records.md, the layout |
| C43 | `telemetry.tsv` (file) | 12 | records | `layup run` | 1 (schema only) | done: records.md, the layout; schema in records.md |
| C44 | `prices.tsv` (file) | 12 | host | the Operator | 1 (schema only) | done: records.md, the layout; schema in records.md |
| C45 | the stall row, the diagnosis row, the outcome row | 11 | records | `layup run` | 1 (schema only) | done: records.md, the layout; schema in records.md |
| C46 | the "failed" row and the "milestone stopped" row | 11 | records | `layup run` | 3 | done: records.md, the layout |
| C47 | the stall package, as a payload | 11 | records | `layup run` | 3 | done: records.md, the layout |
| C48 | `human-inputs.tsv` (file) | 12 | records | `layup run` | 3 | done: records.md, the layout |
| C49 | `audit.tsv` (file), with the seed and the population hash | 12 | records | `layup run` | 3 | done: records.md, the layout |
| C50 | `acceptance.tsv` (file) | 8, 12 | records | `layup run` | 2 | done: records.md, the layout |
| C51 | `baseline.tsv` (file) | 5, 12 | records | `layup run` | 4 | done: records.md, the layout |
| C52 | `start-values.tsv` (file) | 12 | records | `layup run` | 4 | done: records.md, the layout |
| C53 | the reward table and the proposed weights | 13 | records | `layup learn`, then `layup run` | 4 | done: records.md, the layout |
| C54 | a lesson, with its records and scope | 13 | records | `layup run` | 4 | done: records.md, the layout |
| C55 | `setup/steps.tsv` (file): LAYUP's own setup steps | 5 | LAYUP | a LAYUP task | 1 (read) | done: records.md, the layout; unchanged (C6) |
| C56 | the last run time of the dead-man job | 11 | the control repository | the dead-man job | 3 | done: records.md, the layout |

## D. The phase-1 boundaries (AC 7, C1, C4, C6)

| Row | Item | Status |
| --- | ---- | ------ |
| D1 | One row per step S01 to S15 with the phase-1 actor, inputs, outputs, evidence (C1) | done: setup.md, the steps |
| D2 | Each forge call of Start and the Scaffold (§3, §5) named in or out | done: setup.md, the boundary of phase 1 |
| D3 | Each session of Start and the Scaffold (the review of meaning, the specification sessions, the marker sources, the prose rows) named in or out | done: setup.md, the boundary of phase 1 |
| D4 | Each check of `layup setup verify` (§5 step 5) in or out; out is not in the list (C4) | done: setup.md, the checks of layup setup verify |
| D5 | Each job of `layup gate` (§6) in or out (C4) | done: gate.md, REQ-004 not in phase 1; REQ-007 |
| D6 | `docs/setup/steps.tsv` does not change; `setup.md` is the one home of the phase-1 rows of a target (C6) | done: setup.md, the steps (first paragraph) |
| D7 | The home of the setup record and the rule-path register in phase 1: the Operator's answer to the question of the author's answer (C2) | done: records.md NFR-001; setup.md, where the records go in phase 1 (O-115) |
| D8 | A stop prints every missing input of its row in one table (C3) | done: setup.md, the stop table |

## E. Conventions and the other criteria

| Row | Item | AC | Status |
| --- | ---- | -- | ------ |
| E1 | `docs/spec/README.md`: the schema block (one column per line, no tab; the header row is the names joined by one tab; the info string may hold a path pattern, N5), the type names, the exit-code rule, the index | 4 | done: README.md, records |
| E2 | Each value that the specification decides is marked "decided here" with its reason; a value that waits for the Operator is a marker with a row in `open-gaps.tsv` (N11) | 8 | done: README.md, how a section is written; each "decided here" |
| E3 | `PRD-0001`: `REQ-012` criterion names `docs/spec/`; §9 phase 1 per D1; the §7.1 words of `REQ-002` on `steps.tsv` (C6); §12; change log | 7, 9 | done: PRD-0001 §7.1, §9, §12, §13 |
| E4 | Stale documents: one link in `docs/architecture.md` or none (N8); glossary "Phase" (N13); README; onboarding | 11 | done: architecture.md intro; glossary Phase; README; onboarding |
| E5 | All local checks and `git diff --check` | 10 | done: `runs/T-0drh/test-runs.md` |
| E6 | The review brief names the two reads of AC 8 (N14) | 8 | open: step 11 |
