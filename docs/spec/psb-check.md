# `layup psb check`

Conventions: [`README.md`](README.md). The code exists (`internal/psb`, tasks
`T-dq05` and `T-5zmw`); this section states its contract, and the golden test
`TestGoldenRealPSB` holds it to `internal/psb/testdata/psb.tsv`.

## REQ-001 — The rule gaps of a problem statement, in one batch

Requirement: "Before delivery starts, `layup psb check` finds the rule gaps of a
problem statement, a review session finds its gaps of meaning, and both go to
the idea owner as one batch of questions, whose answers are stored as a raw
fact." Derives from `architecture.md` §5 (the gap check, step 1), ADR-0011 and
ADR-0015. Phase 1 delivers the rule gaps (`PRD-0001` §9).

### The command

```text
layup psb check FILE
```

- `FILE`: a problem statement in Markdown, UTF-8. A carriage return before a
  line feed is ignored. A field of the table follows the field rule of
  [`README.md`](README.md#records), so a tab or a lone carriage return in
  `excerpt` becomes one space.
- No flag.
- Exit codes: 0 no gap; 1 at least one gap; 2 a usage error (not exactly the
  two words `check FILE`), a file that cannot be read, a `FILE` that is not
  valid UTF-8, or a table that the command cannot write
  ([`README.md`](README.md#commands)). After a usage error, a file that cannot
  be read, or a `FILE` that is not valid UTF-8, standard output holds nothing.

**Decided here** (task `T-5zmw`, #83):

- **A `FILE` that is not valid UTF-8** (K32): no rule runs, and standard error
  holds `layup: FILE: line <n> is not valid UTF-8`, where `<n>` is the line of
  the first byte that is not valid. Reason: a record is UTF-8, a repair would
  put into `excerpt` a value that the file does not hold, and the line lets the
  author find the byte.
- **A table that the command cannot write** (for example on a full disk, or to
  a standard output that is open for reading only): `layup: <the error>` on
  standard error, exit 2. Reason: the run cannot give its result, and a check
  that did not give its result never gives 0 (`NFR-004`). A closed pipe stops
  `layup`, as any Go program, with the signal `SIGPIPE`, so it gives no code.

### The rules

Each rule reads text only; none calls a model (`NFR-005`). Lines are counted
from 1.

| Rule | A gap when | Line | Question |
| ---- | ---------- | ---- | -------- |
| G1 | no line holds `technology stack` (any case), then only spaces or `*`, then `:`, then only spaces or `*`, then a character that is neither white space nor `*` (O-131); after the `:`, white space other than a tab, a form feed or a carriage return counts as a space | 0 | `Which technology stack does the product use (languages, frameworks, tools)?` |
| G2 | a row of a Markdown table whose header row has a column `Measurement` or `Verification` (any case) has, in that column, an empty cell or one of `—`, `-`, `tbd`, `not measured` (any case, one final `.` ignored) | the row | `How is this metric measured? The row gives no measurement method.` |
| G3 | the text has a table under a heading that contains `Terms`, and a line holds an abbreviation (2 to 6 characters: a capital, then capitals or digits, with ASCII word boundaries) that no bold span of that table's first column defines; each abbreviation once, at its first line | the line | `What does "<ABBR>" mean? The terms table does not define it.` |
| G4 | a line holds one of the words `fast`, `robust`, `soon`, `clean`, `better`, `handle` (whole word, any case) and no digit; the first such word of the line | the line | `Which number or threshold does "<word>" stand for here?`, the word in lowercase |
| G5 | a line holds `(start value)` (any case) | the line | `Which start value does the pilot use for this target, and who sets it?` |

With no terms table, G3 gives no gap. In the table, a space is U+0020, and
white space is a character of the Unicode property `White_Space` (for example a
tab, a vertical tab, U+00A0 or U+2003).

**Decided by the Operator** (O-131, #83): `**Technology stack:**` with no value
is a G1 gap, because the value of a named stack is a character that is neither
white space nor `*`; so an empty label of a template asks for the stack.
**Decided here** (review round 1 of #83), as the reading of O-131: a value that
is only white space (for example U+00A0) is no value, and after the `:` the
white space other than a tab, a form feed or a carriage return counts as a
space, so `Technology stack:` with U+00A0 and then `Go` still names a stack.
Reason: the words of O-131 say "white space", and Go's `\s` in its example
expression is only a tab, a line feed, a form feed, a carriage return and a
space; with the second part, a result of the code before O-131 changes only
where that code read `*` or a white space character as the value.

The present code is the reference for an edge case that this table does not
settle; such a case found later is added here and to
`internal/psb/testdata/edge.md`, or to a Go test when it needs a file of its
own: a case of G1 to `TestG1ReadsTheValueOfAStack`, and a case of G3 to
`TestEdgeCases`, because these two rules read the whole file.
**Decided here** (task `T-5zmw`, D3 of #83), as the code reads them:

- G1: a tab before the `:` is not a named stack (a gap).
- G2: a bold header (`**Measurement**`) is not a measurement column; when a
  header row has both `Measurement` and `Verification`, the first of them is the
  column; an escaped `\|` splits a cell all the same.
- G3: a heading `terms` in lower case is not a terms heading; the first heading
  that holds `Terms` (for example `## Payment Terms`) gives the terms table, and
  a later one is not read; a byte-order mark before a `# Terms` on the first
  line stops that heading; an abbreviation inside the terms table is read as on
  any other line.
- Each rule: a `#` line inside a code fence is a heading, and G4 applies inside
  a fence.

The golden `internal/psb/testdata/values.md` pins the values of this table that
the other goldens do not: the column name `Verification`; `—`, `-`, `tbd` and
`not measured` in mixed case, one with a final `.`; each word of G4 in mixed
case, two of them on one line, and one on a line with a digit; `(start value)`
in mixed case; and `technology stack` in capitals.

### The table

```tsv-schema psb-gaps stdout
id        id(Q-NNN)          key  the row number in the table, from Q-001, in the order of the rows
rule      enum(G1|G2|G3|G4|G5)  -  the rule that found the gap
line      int                -    the line of the gap; 0 when the gap is an absence (G1)
excerpt   text               -    the line, trimmed, at most 80 characters, by the field rule; `—` for line 0
question  text               -    the question for the idea owner, from the table of rules
```

The rows are in the order of the line, then the rule, then the order on the
line. The header row is printed when there is no gap, too. **Decided here**
(task `T-5zmw`, D1 of #83): `internal/psb` holds the Go schema of this block,
`GapsSchema`, and writes the table with `tsv.Write`, so the field rule has one
owner; its integration test compares the schema with the block and reads each
golden table by the block, `psb.tsv` included.

### Not in phase 1

- **The review of meaning** (§5 gap check step 2): a session; it needs `layup
  run` (phase 2).
- **One batch on the Intake issue**, merged with the review's rows and the setup
  questions, with a question ID per row (§5 step 5): `layup run` (phase 2). In
  phase 1 the `id` above is the question ID of the gap: `layup setup` asks it at
  step S01 ([`setup.md`](setup.md#the-steps)); the problem statement does not
  change during a setup, so the IDs stay.
- **The answers as a raw fact.** For LAYUP itself, the answers are `F-0004`.
  For a target, phase 1 writes the answers as a raw fact at setup step S04, in
  the first commit on the setup branch (O-124;
  [`setup.md`](setup.md#the-steps)); the copy of the answer comments is
  `layup run`'s (phase 2).
