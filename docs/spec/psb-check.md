# `layup psb check`

Conventions: [`README.md`](README.md). The code exists (`internal/psb`, task
`T-dq05`); this section states its contract, and the golden test
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
  [`README.md`](README.md#records); the present code does not yet replace a
  lone carriage return in `excerpt` with a space (#29 does).
- No flag.
- Exit codes: 0 no gap; 1 at least one gap; 2 a usage error (not exactly the
  two words `check FILE`), a file that cannot be read, or a `FILE` that is not
  valid UTF-8 ([`README.md`](README.md#commands)). The present code does not
  yet give 2 for a `FILE` that is not valid UTF-8: it reads the bytes as they
  are, and can copy them into `excerpt` (task `T-5zmw` adds the check).

### The rules

Each rule reads text only; none calls a model (`NFR-005`). Lines are counted
from 1.

| Rule | A gap when | Line | Question |
| ---- | ---------- | ---- | -------- |
| G1 | no line holds `technology stack` (any case), then only spaces or `*`, then `:`, then only spaces or `*`, then a character that is not white space | 0 | `Which technology stack does the product use (languages, frameworks, tools)?` |
| G2 | a row of a Markdown table whose header row has a column `Measurement` or `Verification` (any case) has, in that column, an empty cell or one of `—`, `-`, `tbd`, `not measured` (any case, one final `.` ignored) | the row | `How is this metric measured? The row gives no measurement method.` |
| G3 | the text has a table under a heading that contains `Terms`, and a line holds an abbreviation (2 to 6 characters: a capital, then capitals or digits, with ASCII word boundaries) that no bold span of that table's first column defines; each abbreviation once, at its first line | the line | `What does "<ABBR>" mean? The terms table does not define it.` |
| G4 | a line holds one of the words `fast`, `robust`, `soon`, `clean`, `better`, `handle` (whole word, any case) and no digit; the first such word of the line | the line | `Which number or threshold does "<word>" stand for here?`, the word in lowercase |
| G5 | a line holds `(start value)` (any case) | the line | `Which start value does the pilot use for this target, and who sets it?` |

With no terms table, G3 gives no gap. The present code is the reference for an
edge case that this table does not settle; such a case found later is added
here and to `internal/psb/testdata/edge.md`.

### The table

```tsv-schema psb-gaps stdout
id        id(Q-NNN)          key  the row number in the table, from Q-001, in the order of the rows
rule      enum(G1|G2|G3|G4|G5)  -  the rule that found the gap
line      int                -    the line of the gap; 0 when the gap is an absence (G1)
excerpt   text               -    the line, trimmed, a tab as a space, at most 80 characters; `—` for line 0
question  text               -    the question for the idea owner, from the table of rules
```

The rows are in the order of the line, then the rule, then the order on the
line. The header row is printed when there is no gap, too.

### Not in phase 1

- **The review of meaning** (§5 gap check step 2): a session; it needs `layup
  run` (phase 2).
- **One batch on the Intake issue**, merged with the review's rows and the setup
  questions, with a question ID per row (§5 step 5): `layup run` (phase 2). In
  phase 1 the `id` above is the question ID of the gap: `layup setup` asks it at
  step S01 ([`setup.md`](setup.md#the-steps)); the problem statement does not
  change during a setup, so the IDs stay.
- **The answers as a raw fact.** For LAYUP itself, the answers are `F-0004`.
  For a target, phase 1 writes the answers as a raw fact at setup step S06
  ([`setup.md`](setup.md#the-steps)); the copy of the answer comments is
  `layup run`'s (phase 2).
