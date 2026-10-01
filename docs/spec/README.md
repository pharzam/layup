# The technical specification of LAYUP

This directory is LAYUP's own technical specification (`REQ-012`): for each
requirement of [`PRD-0001`](../prd/PRD-0001-layup.md), a section whose heading
holds the requirement ID and that gives the exact contract that the code
implements. It has the form that [`architecture.md`](../architecture.md) §7 asks
of a target, so `layup spec check` can later read it.

It is written one milestone at a time (task `T-0drh`, #74; decision O-114).
This version covers **phase 1** of `PRD-0001` §9 only. A later milestone adds
its own sections; it does not rewrite the package table or the layout of the
records, which already name the later phases.

## The files

| File | Requirements | What it gives |
| ---- | ------------ | ------------- |
| [`packages.md`](packages.md) | `NFR-007` | the Go packages, their jobs, and the import rules that the boundary gate reads |
| [`records.md`](records.md) | `NFR-001`, `NFR-002`, `REQ-009`, `REQ-011` | where each record lives, who writes it, and the schemas of the phase-1 records |
| [`psb-check.md`](psb-check.md) | `REQ-001` | `layup psb check` |
| [`setup.md`](setup.md) | `REQ-002`, `NFR-003`, `NFR-006` | `layup setup`, `layup setup verify`, the stack catalog |
| [`gate.md`](gate.md) | `REQ-004`, `REQ-007`, `NFR-004`, `NFR-005` | `layup gate`, the gate manifest |

## How a section is written

- The heading holds the requirement ID: `## REQ-004 — …`.
- Its first lines name the requirement, and the architecture section and the
  ADR that it derives from, where one exists. A quote of a requirement writes
  the name of the baseline as `[the baseline]`, because outside the files that
  check `adapted` allows, the baseline is not named.
- A value that the architecture does not give, and that this specification
  sets, is marked **decided here**, with its reason. A value that needs a
  decision of the Operator or the idea owner is never set here: it stays a
  marker with a row in [`open-gaps.tsv`](../setup/open-gaps.tsv) (Invariant 4).
- A part of a requirement that a later phase delivers is named, with that
  phase, under **Not in phase 1**.

## Commands

Every command of LAYUP follows these rules. A section gives only what differs.

- **Arguments.** A command takes its inputs as arguments and flags. It reads no
  environment variable for an input, and asks no question on the terminal.
- **Output.** A command prints its result as one table on standard output, in
  the form of a record (below). It prints a diagnostic only on standard error.
- **Exit codes** (decided here; they extend the codes of the present
  `layup psb check`):

  | Code | Meaning |
  | ---- | ------- |
  | 0 | Every row of the table passed: `pass`, or `clear` where the section allows it. For `layup psb check`: no gap. |
  | 1 | At least one row is `fail` or `not-active`, or (for `layup psb check`) at least one gap. |
  | 2 | A usage or input error: an unknown flag, a missing or unreadable file, an input that does not match its schema. The table can be empty or incomplete. |
  | 3 | Only `layup setup`: the run stopped at a step that needs a human input; the table lists every missing input of that step. |

  A check that did not run is never 0 (`NFR-004`). Code 3 is not a failure and
  not a usage error: a later `layup run` reads it as "wait for the answers", so
  it gets its own code.
- **Determinism.** Two runs on the same input print the same bytes
  (`NFR-005`). So a result table holds no time, no duration and no path of a
  scratch directory, and its rows have a fixed order.

## Records

A record is a table of tab-separated values in UTF-8, with a header row and one
row per item, with line-feed line endings and a line feed after the last row
(ADR-0011 decision 2, ADR-0014). A field holds no tab, no line feed and no
carriage return: a writer replaces each with one space. An empty field is
written as `—` (U+2014), never as an empty string, so that a human sees it.

### The schema block

Each record has one schema block. A Go test of the code reads these blocks and
compares them with what the code writes and reads (the test comes with the
code, Bootstrap mode rule 1). The form:

````text
```tsv-schema <name> <location>
<column> <type> <key> <rule>
…
```
````

- `<name>` is the record's name, unique in `docs/spec/`.
- `<location>` is where the record lives: `stdout` (a command's table),
  `records:<path>` (the target's records branch `layup-records`),
  `target:<path>` (the target's default branch), `layup:<path>` (LAYUP's own
  repository) or `host:<path>` (the work area of the LAYUP host). A path may
  hold a pattern part in angle brackets, for example
  `records:tasks/<task>/events.tsv`.
- One line per column, in the column order. The fields of a line are separated
  by one or more spaces: the column name, its type, `key` or `-`, and the rest
  of the line is the rule in words.
- The block holds no tab. The header row of the record is the column names, in
  the order of the lines, joined by one tab.
- A key of more than one column marks each of its columns `key`.

### The types

The list is closed; a new type is added here first.

| Type | Values |
| ---- | ------ |
| `text` | any text of one line |
| `int` | a decimal integer, 0 or more, with no sign and no leading zero |
| `decimal` | a decimal number, 0 or more, with a point and no exponent, for example `0.42` |
| `bool` | `yes` or `no` |
| `time` | a time in UTC, RFC 3339, to the second: `2026-10-01T08:09:21Z` |
| `sha1` | a Git object name: 40 lowercase hexadecimal characters |
| `sha256` | 64 lowercase hexadecimal characters |
| `path` | a path relative to a repository's root, with `/` and no `..` |
| `enum(a\|b\|c)` | one of the listed words |
| `id(<pattern>)` | an identifier in the named form: `N` is a digit, `x` is a lowercase letter or digit as the column's rule says, and `<word>` is lowercase letters and `-`; for example `id(Q-NNN)` |
| `list(<type>)` | values of that type, separated by one space, in a fixed order |

A column that may be unknown says so in its rule, and has a status column
beside it; an unknown value is never written as 0 (FT2).

## Phase 1 and later phases

Phase 1 of `PRD-0001` §9 ships the engine checks and `layup setup`; it starts
no role session, runs no phase loop, and makes no forge call. The orchestrator
`layup run`, which is the one writer of a target's records branch (ADR-0014),
comes in a later phase. So a section of phase 1 that needs `layup run` gives the
schema and names `layup run` as the writer, and each part that waits for it is
under **Not in phase 1**.
