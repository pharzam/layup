# The technical specification of LAYUP

This directory is LAYUP's own technical specification (`REQ-012`): for each
requirement of [`PRD-0001`](../prd/PRD-0001-layup.md), a section whose heading
holds the requirement ID and that gives the exact contract that the code
implements. It has the form that [`architecture.md`](../architecture.md) §7 asks
of a target, so `layup spec check` can later read it.

It is written one milestone at a time (task `T-0drh`, #74; decision O-114).
This version covers **phase 1** of `PRD-0001` §9 and milestone **`M2a`** of
phase 2 (task `T-zck8`, #123). A later milestone adds its own sections; it does
not rewrite the package table or the layout of the records, which already name
the later phases.

## The files

| File | Requirements | What it gives |
| ---- | ------------ | ------------- |
| [`packages.md`](packages.md) | `NFR-007` | the Go packages, their jobs, and the import rules that the boundary gate reads |
| [`records.md`](records.md) | `NFR-001`, `NFR-002`, `REQ-009`, `REQ-011` | where each record lives, who writes it, and the schemas of the phase-1 records |
| [`psb-check.md`](psb-check.md) | `REQ-001` | `layup psb check` |
| [`setup.md`](setup.md) | `REQ-002`, `NFR-003`, `NFR-006` | `layup setup`, `layup setup verify`, the stack catalog |
| [`gate.md`](gate.md) | `REQ-004`, `REQ-007`, `NFR-004`, `NFR-005` | `layup gate`, the gate manifest |
| [`run.md`](run.md) | `NFR-001`, `NFR-002`, `NFR-006`, `REQ-002` | `layup run`: the Start of a target and the restart (`M2a`) |
| [`forge.md`](forge.md) | `NFR-001`, `NFR-007` | the forge interface and the GitHub adapter (`M2a`) |

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

- **`layup version`** is exempt from the table: it prints one line,
  `layup <version>`, and exits 0. **Decided here** (K34 of the
  [defect register](../plan/README.md#the-defect-register)): it takes no
  argument, so `layup version extra` is a usage error (exit 2). Reason: one
  argument rule for every command; an extra word that a command ignores hides a
  typing error.
- **Arguments.** A command takes its inputs as arguments and flags. It reads no
  environment variable for an input, does not read the standard input, and asks
  no question on the terminal. **Decided here** (task `T-2yw7`, #80):
  - `internal/cli` holds one command table: the words of each command, its
    positional arguments, its flags and its line of the usage. The dispatch, the
    check of the arguments and the usage come from it. The longest match of the
    words decides the command: `layup setup verify WORK` is `setup verify`, and
    `layup setup verify` with no `WORK` is a usage error, not a setup of a
    directory named `verify` (such a directory is given as `./verify`).
  - A flag is `--name VALUE` or `--name=VALUE`, before or after the positional
    arguments. Each flag of phase 1 is required and is given once. There is no
    short flag and no `--`; a positional argument that starts with `-` is given
    as `./-name`.
  - A usage error is: no command, an unknown command, an unknown flag, a flag
    with no value (also `--name=` and `--name` before a word that starts with
    `-`), a missing or repeated flag, and a missing or extra positional
    argument.
  - The test `TestInputRule` reads the non-test Go files of the module and finds
    no read of an environment variable or of the standard input, except in the
    function `environ` of `internal/git`, which gives `PATH` and `TMPDIR` of the
    host to `git` ([`packages.md`](packages.md#the-calls-of-internalgit), D3 of
    #79). A program that `layup` starts gets its environment from `layup`; that
    is not an input of `layup`, so the test does not cover it.

  Reason: a value that a command reads from its environment has no record, and
  two runs on the same arguments and files could differ (Determinism, below).
- **Output.** A command prints its result as one table on standard output, in
  the form of a record (below). It prints a diagnostic only on standard error.
  **Decided here** (task `T-2yw7`): a usage error prints `layup: <reason>`, an
  empty line and the usage on standard error, and nothing on standard output;
  an input error prints only `layup: <reason>` on standard error. The usage
  lists each command with its arguments and its line, then the exit codes.
- **Progress** (decided here, K29): a command whose run can take more than ten
  seconds (`layup gate`, `layup setup`, `layup setup verify`) prints its
  progress on standard error, never on standard output: before each step, one
  line `layup <command>: [<i>/<n>] <step>`, and while that step runs, one more
  line every ten seconds, `layup <command>: [<i>/<n>] <step>: <s> s`. A step is
  a row of the command's own table: a kind, a step, a check. `layup version` and
  `layup psb check` print none. `internal/cli` prints the lines; the package
  that runs the steps gets a function that it calls at the start of each step,
  and the command stops the lines of its last step before it prints its table,
  so no progress line comes after the table. The output of a command that a
  step runs (a gate command, a baseline script) is held and printed after its
  step ends, so no progress line comes inside it (decided here, task `T-5sgt`,
  #82). Reason: the
  [progress rule](../engineering-discipline.md#progress-indicators-for-long-running-operations)
  asks which step runs, how much remains, and that the work is alive; a line
  per step and a line every ten seconds answer the three in a terminal and in a
  CI log, where a status line that redraws itself would not.
- **Exit codes** (decided here; they extend the codes of the present
  `layup psb check`):

  | Code | Meaning |
  | ---- | ------- |
  | 0 | Every row of the table passed: `pass`, or `clear` where the section allows it. For `layup psb check`: no gap. |
  | 1 | At least one row is `fail` or `not-active`, or (for `layup psb check`) at least one gap. |
  | 2 | A usage or input error: an unknown flag, a missing or unreadable file, an input that does not match its schema, an input that the command parses as text and that is not valid UTF-8 (below). For `layup gate` also `git` or `sh` not found, `git` older than 2.32, and a scratch work tree that it could not remove ([`gate.md`](gate.md#the-command)). For `layup psb check` also a table that it cannot write to standard output ([`psb-check.md`](psb-check.md#the-command)). For `layup setup verify` also `git` older than 2.32, a `docs/gates.tsv` at the head of `layup-setup` that is missing or malformed, a temporary directory in `WORK`, and a scratch work tree that it could not remove ([`setup.md`](setup.md#the-checks-of-layup-setup-verify)). The table can be empty or incomplete. |
  | 3 | Only `layup setup`: the run stopped at a step that needs a human input; the table lists every missing input of that step. |

  A check that did not run is never 0 (`NFR-004`). Code 3 is not a failure and
  not a usage error: a later `layup run` reads it as "wait for the answers", so
  it gets its own code.

  **Decided here** (task `T-2yw7`): one map in `internal/cli` gives the code of
  a table from its result column: 0 only when the table has at least one row
  and each row is `pass`, `clear`, `done` or `operator`; each other word, and a
  table with no row, give 1. It takes the result words of the three tables
  (`gate-result`, `setup-verify`, `setup-steps`); a word that a table's schema
  refuses (for example `done` in a `gate-result` row) cannot reach it, because
  the writer of the record refuses it first. Codes 2 and 3 do not come from a
  table. The rule for a table with no row is the default of the frame; for
  `layup gate`, a manifest with no row is an input error (exit 2,
  [`gate.md`](gate.md#the-command)).

  **Decided here** (K32 of the [defect register](../plan/README.md#the-defect-register)):
  a command that parses an input as text (a record, a problem statement, an
  answers file, a Markdown file that a step checks) and finds bytes that are not
  valid UTF-8 gives exit 2, with a diagnostic that names the file. A file that a
  command only copies or hashes is not under this rule. Reason: a record is
  UTF-8, and a repair would put a value into a record that no input holds.
  `layup psb check` applies it (task `T-5zmw`, row 6 of the plan), with the
  line of the first byte that is not valid
  ([`psb-check.md`](psb-check.md#the-command)).
- **Determinism.** Two runs on the same input print the same bytes
  (`NFR-005`). So a result table holds no time, no duration and no path of a
  scratch directory, and its rows have a fixed order. **Decided here** (task
  `T-2yw7`): the same bytes means the same exit code and the same standard
  output. Standard error is not compared: its diagnostics and progress lines can
  hold a time. The end-to-end harness of `cmd/layup` compares the two in its
  repeat helper, which each command's scenarios use.

## Records

A record is a table of tab-separated values in UTF-8, with a header row and one
row per item, with line-feed line endings and a line feed after the last row
(ADR-0011 decision 2, ADR-0014). A field holds no tab, no line feed and no
carriage return: a writer replaces each with one space. An empty field is
written as `—` (U+2014), never as an empty string, so that a human sees it.

**Decided here** (task `T-18v6`, #78), for the reader and the writer of
`internal/tsv` ([`packages.md`](packages.md)):

- **The empty field.** `—` is valid in each column that is not a key column,
  whatever its type, and the reader gives it as the empty value. The package
  that owns a record checks a column whose rule forbids `—`. A key column never
  holds `—`. The writer writes `—` for an empty value and for the text `—`, so
  the text `—` reads back as the empty value. Apart from the field rule, this is
  the one case where a read does not give back the value that the writer took.
  Reason: the types have no mark for "may be empty", and the rule of a column
  belongs to its owner.
- **The reader** refuses a byte-order mark at its start, a carriage return anywhere (a
  carriage return before a line feed included), bytes that are not valid UTF-8,
  an empty line, no line feed after the last line, a header row that is not the
  column names, the column names as the first line of a record with no header
  row, a row with too few or too many fields, a field that is an empty string, a
  field that is not a value of its column's type, `—` in a key column, and a key
  that an earlier row has. Each error names the line. The error for one field
  also names its column: a carriage return or bytes that are not valid UTF-8 in
  the field, an empty string, a value that is not of the column's type, `—` in
  a key column, a wrong name in the header row, or a missing field. The error
  for a carriage return also names the fix: line-feed endings. A command gives
  exit 2 for such an input. Reason: code writes each record with line feeds
  only, and a person who writes a record by hand (`answers.tsv`) gets an error
  that says what to fix, not a silent repair.
- **The writer** applies the field rule and writes `—` for an empty value. Then
  it makes each check of the reader, so it never writes a record that the reader
  refuses. A value that is not valid UTF-8 is an error, not a repair. The writer
  cannot tell one list value that holds a space from two values, so code makes
  the field of a `list(<type>)` column with `tsv.JoinList`, which refuses such a
  value.

### The schema block

Each record has one schema block. A Go test of the code reads these blocks and
compares them with what the code writes and reads (the test comes with the
code, [the scope of a task](../engineering-discipline.md#working-a-task-under-the-quality-gate)). The form:

````text
```tsv-schema <name> <location>
<column> <type> <key> <rule>
…
```
````

- `<name>` is the record's name, unique in `docs/spec/`. A third word
  `no-header` marks a record whose form another check fixes with no header row
  (for example `open-gaps.tsv`); the column lines still give the order.
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
- A key of more than one column marks each of its columns `key`. **Decided
  here:** the key of a record is the tuple of the values of all its `key`
  columns, and no two rows have the same key; the reader and the writer check
  it. Reason: a key that repeats is a defect that no other check finds.
- **Decided here:** the parser of `internal/tsv` refuses a block with a tab, an
  empty line, a column line with no rule, a type that is not in the list below,
  a word other than `key` or `-` in the place of the key, a column name that
  repeats, a location with another prefix or with a path that is not in the form
  of the type `path`, a name that another block has, no column, or no closing
  fence. A fence inside a longer fence is text, not a block, as in the example
  above. A block starts a line, after at most three spaces. The parser refuses
  a `tsv-schema` fence with other text before it on its line (a blockquote mark,
  a list marker, a tab, a date, or four spaces or more), and a directory that it cannot read or that holds no block, so that no
  block escapes the test. Reason: each block has one reading. Known limit: the
  parser reads a block in an HTML comment as a block.

**Decided here:** how "the test comes with the code" is read. Reason:
`internal/tsv` imports no package of this module, so its test cannot see the Go
schema of a record's owner.

- The package that owns a record holds the record's schema as a Go value. Its
  code writes and reads the record with that value, and its own test compares
  the value with the block (`tsv.Compare`).
- One test of `internal/tsv`, with the tag `integration`, parses each block of
  the Markdown files at the root of `docs/spec/` by the form above. It holds two lists of block names, `built` and
  `notYetBuilt`. It checks that each block is in exactly one list, and that each
  listed name is a block. An owner moves the name of its block from
  `notYetBuilt` to `built` in the same change as its own test.
- A later section adds the names of its new blocks to `notYetBuilt`;
  otherwise `notYetBuilt` only becomes shorter. This is a rule for the reviewer
  of each owner task, not a check: a test cannot read the list of its base.
- A block that a later section adds makes that test fail until its name is in
  one of the two lists.

### The types

The list is closed; a new type is added here first.

| Type | Values |
| ---- | ------ |
| `text` | any text of one line |
| `int` | a decimal integer, 0 or more, with no sign and no leading zero |
| `decimal` | a decimal number, 0 or more, with a point and no exponent: an `int`, a point, and one or more digits, for example `0.42` |
| `bool` | `yes` or `no` |
| `time` | a time in UTC, RFC 3339, to the second: `2026-10-01T08:09:21Z` |
| `sha1` | a Git object name: 40 lowercase hexadecimal characters |
| `sha256` | 64 lowercase hexadecimal characters |
| `path` | a path relative to a repository's root, with `/`; no empty, `.` or `..` part, and no `/` at its start or end |
| `enum(a\|b\|c)` | one of the listed words |
| `id(<pattern>)` | an identifier in the named form. In the pattern, a run of `N` is that many digits, or more digits with no leading zero, so `id(Q-NNN)` takes `Q-001` and `Q-1000`, and not `Q-0001`; each `x` is one lowercase letter or digit; `<word>` is one or more lowercase letters or `-`; each other character stands for itself. A column's rule can name a smaller set for `x` (for example hexadecimal), which the owner of the record checks |
| `list(<type>)` | values of that type, separated by one space, in a fixed order; a value holds no space, and an empty list is `—` |

**Decided here:** the readings of `decimal`, `path`, `id` and `list` in the
table. Reasons: `layup psb check` writes `Q-%03d`, which has no maximum and one
form for each number, so a run of `N` is a minimum, and a longer number has no
leading zero; an ID and a path have one form each, so two rows that name one
thing have the same text, and the key check finds them; a decimal follows the
rule of `int` before its point, so `.5`, `5.` and `05.5` are not values; one
space separates the values of a list, so a value cannot hold one.

A column that may be unknown says so in its rule, and has a status column
beside it; an unknown value is never written as 0 (FT2).

## Phase 1 and later phases

Phase 1 of `PRD-0001` §9 ships the engine checks and `layup setup`; it starts
no role session, runs no phase loop, and makes no forge call. The orchestrator
`layup run`, which is the one writer of a target's records branch (ADR-0014),
comes in a later phase. So a section of phase 1 that needs `layup run` gives the
schema and names `layup run` as the writer, and each part that waits for it is
under **Not in phase 1**. Milestone `M2a` specifies the Start of `layup run`
([`run.md`](run.md), [`forge.md`](forge.md)); a part that a later milestone
delivers is under **Not in M2a**, with that milestone.
