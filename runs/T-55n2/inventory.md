# T-55n2 — the inventory of the phase-1 work (plan step 1)

Evidence for the plan of [#76](https://github.com/pharzam/layup/issues/76). The workflow `t55n2-understand` (run `wf_051087d7-99c`, 2026-10-01, 58 min) ran 22 agents of Claude Opus 5.5, read-only, in a detached worktree at `7cdd346`. Seven areas each ran a reader, an adversarial critic and a reviser; one cross-area critic then checked all of them. This file is the result, transcribed by code with no change of content: the angle quotes of a marker are written as HTML entities, so that check `markers` reads no marker here.

The column **Host** was empty at plan step 1 (the red half). Plan step 4 filled it by code from [`docs/plan/README.md`](../../docs/plan/README.md) with [`fill-hosts.py`](fill-hosts.py): a row of the task table, a milestone, a task, a decision of the Operator, or `out`. The one home of the hosts is the plan; [`plan-check.py`](plan-check.py) checks that this column agrees with it.

| Area | Items | Open questions | External inputs | Constraints | Risks |
| ---- | ----- | -------------- | --------------- | ----------- | ----- |
| `psb` | 10 | 14 | 5 | 14 | 6 |
| `gate` | 12 | 33 | 7 | 24 | 14 |
| `setup` | 18 | 60 | 15 | 19 | 13 |
| `verify` | 16 | 38 | 6 | 23 | 12 |
| `found` | 10 | 35 | 9 | 23 | 10 |
| `gov` | 22 | 24 | 23 | 38 | 13 |
| `later` | 15 | 23 | 36 | 20 | 10 |

## Area `psb`

`layup psb check` exists. go build, go vet and go test pass at 7cdd346. The code obeys every value of the rule table at psb-check.md:34-40, the sort order, the exit codes 0/1/2, the header row and the version exemption. But the table has readings that only the code settles. One of them is the G1 case: `**Technology stack:**` with no value gives no G1 gap (I ran it), because the regex at check.go:25 accepts a `*` as the value character. The phase-1 deltas are small. (1) A lone CR reaches `excerpt` (psb-field-rule-cr). This fix lands first in internal/psb, then psb-tsv moves the field rule to internal/tsv, so the rule has one owner. (2) The writer moves to internal/tsv, and a Go test reads the `psb-gaps` schema block (psb-tsv). (3) internal/cli hands the gap table to internal/setup for S01 and S06 through a seam that the plan must choose. internal/cli may not import internal/tsv (packages.md:37) (psb-batch-api). (4) An e2e test of the REQ-001 command. The tests that touch files move off the unit level. The doc comment of Run (cli.go:24-25) omits exit 1, but the usage text (cli.go:20-21) already names 0, 1 and 2 (psb-cli-contract). (5) The edge cases that psb-check.md:42-44 requires to be added to the table and to edge.md (psb-rule-edges). (6) Golden cases for each value of the rule table (psb-rule-values). (7) A rule for a FILE that is not UTF-8 (psb-utf8-input). The parts of REQ-001 for later phases are the review of meaning (phase 2), the Intake batch (phase 2), and the "before delivery" question rows with the Early Question Share (phases 3 and 4).

| Key | Title | Requirements | Size | After | Host |
| --- | ----- | ------------ | ---- | ----- | ---- |
| `psb-field-rule-cr` | Replace a lone carriage return in excerpt with a space | REQ-001 | small | — | row 6 |
| `psb-tsv` | Write the gap table through internal/tsv, checked against its schema block | REQ-001, NFR-007 | small | psb-field-rule-cr, found-tsv (foundation area, best key name: the internal/tsv writer, the schema-block reader and the type checker for id(Q-NNN), enum, int and text) | row 6 |
| `psb-batch-api` | Hand the gap table from internal/psb to internal/setup through internal/cli for S01 and S06 | REQ-001, REQ-002 | small | psb-tsv, setup-runner (setup area, best key name: the internal/setup step runner that receives the table at S01; needed only when the chosen design converts the rows into a type of internal/setup) | row 9 |
| `psb-cli-contract` | The command `layup psb check FILE` end to end: exit codes, streams, determinism | REQ-001, NFR-004, NFR-005 | small | — | row 6 |
| `psb-rule-edges` | Add the unsettled edge cases of G1 to G5 to the rule table and to edge.md | REQ-001, NFR-005 | small | psb-field-rule-cr | row 6 |
| `psb-rule-values` | Pin each value of the rule table with a golden case | REQ-001, NFR-005 | small | psb-field-rule-cr | row 6 |
| `psb-utf8-input` | A rule and a test for a FILE that is not valid UTF-8 | REQ-001 | small | psb-field-rule-cr | row 6 |
| `later-psb-meaning-review` | The review of meaning: a session writes the gaps of meaning as typed rows | REQ-001 | large | — | `M2c` |
| `later-psb-intake-batch` | One batch on the Intake issue, with the answers as a raw fact | REQ-001 | large | later-psb-meaning-review, psb-batch-api | `M2c` |
| `later-psb-early-question-share` | The rows 'before delivery' in questions.tsv and the measure of the Early Question Share | REQ-001 | large | later-psb-intake-batch | `M4b` |

### `psb-field-rule-cr` — Replace a lone carriage return in excerpt with a space

- **What:** A lone CR inside a line of the problem statement becomes one space in the `excerpt` field, as a tab does today. CRLF endings stay ignored. The output for a file with no CR byte does not change by one byte. The existing goldens become embedded test data (//go:embed), so the unit level opens no file at run time.
- **Specification:** docs/spec/README.md: Records (the field rule); docs/spec/psb-check.md: The command; docs/spec/psb-check.md: The table
- **Packages:** internal/psb
- **Present state:** internal/psb/check.go:41 changes only CRLF to LF. check.go:110-116 excerpt() replaces only a tab. For a file that holds 'The API\ris fast.', Q-001 is the G1 row with excerpt `—` (check.go:43-48), and Q-002 (G4) holds the raw CR. internal/psb/check_test.go:14,20 reads the goldens with os.ReadFile at run time. docs/spec/psb-check.md:23-24 records the defect ('the present code does not yet replace a lone carriage return in `excerpt` with a space (#29 does)'), and the same PR deletes that sentence.
- **Tests:**
  - unit: Check on the Go literal "Technology stack: Go\nThe API\ris fast.\n" gives a G4 row whose excerpt is "The API is fast." and holds no CR byte. It fails today. A Go literal keeps a raw CR byte out of a tracked file.
  - unit: triggers.md converted to CRLF endings in the test (strings.ReplaceAll) gives triggers.tsv byte for byte (psb-check.md:21-22). No test gives CRLF input today.
  - unit: The goldens triggers, clean and edge, read through //go:embed, stay byte for byte the same (test-levels.md:43: unit tests touch no file).
  - integration: TestGoldenRealPSB still matches internal/psb/testdata/psb.tsv byte for byte. The real PSB holds no CR, so F-0004's mapping of facts 1 to 19 to Q-001 to Q-019 stays true.
- **Open questions:**
  1. docs/spec/psb-check.md:52 says the excerpt is 'the line, trimmed, a tab as a space, at most 80 characters'. The code trims first and then cuts (check.go:111-113), so an excerpt can end with a space. psb.tsv pins this (Q-012 ends in '*(start value)* of '). psb-check.md:42-43 makes the present code the reference. Keep the code unless the Operator decides otherwise.
  2. docs/spec/README.md:108 says `text` is 'any text of one line', and the field rule (README.md:69-70) names only tab, LF and CR. A vertical tab, a form feed, U+0085 or U+2028 inside a line goes into excerpt unchanged (I ran a VT case). Is such a field 'one line'? No source says.

### `psb-tsv` — Write the gap table through internal/tsv, checked against its schema block

- **What:** internal/psb writes the gap table with the record writer of internal/tsv. That writer gives the header, the field rule (tab, LF and CR become a space; an empty field becomes —), LF endings and a final LF. The CR replacement of psb-field-rule-cr moves out of internal/psb into internal/tsv, so the field rule has one owner. A Go test reads the `tsv-schema psb-gaps stdout` block and checks every golden against it. The output bytes do not change.
- **Specification:** docs/spec/psb-check.md: The table; docs/spec/README.md: Records; docs/spec/README.md: The schema block; docs/spec/README.md: The types; docs/spec/packages.md: NFR-007 — The table of phase 1 (internal/psb, internal/tsv)
- **Packages:** internal/psb, internal/tsv
- **Present state:** internal/psb/check.go:101-108 WriteTSV writes the header and the rows with fmt.Fprintf and applies no field rule. internal/psb imports no package of this module. docs/spec/packages.md:46-47 says 'it moves to `internal/tsv` when that package exists (#29)', and the same PR deletes or rewrites that note. internal/tsv does not exist. No test reads a tsv-schema block. The import rule of packages.md:40 is left to the boundary rule of the future Go gate (packages.md:31-32), not to a test of this item.
- **Tests:**
  - unit: The embedded goldens triggers, clean and edge, and the lone-CR literal of psb-field-rule-cr, give the same bytes through the tsv writer as before.
  - integration: The schema test reads docs/spec/psb-check.md:48-54. The header that the code writes equals the column names joined by one tab (README.md:98-99). Every row of every golden .tsv, psb.tsv included, passes the types: id is Q-NNN in row order, rule is G1..G5, line is an int with no leading zero, and excerpt and question are text with no tab, LF or CR.
  - integration: The same test checks the column rules that code can check: excerpt has at most 80 characters; excerpt is `—` when line is 0; no field is empty (README.md:70-71); question matches one of the five texts of the rule table, with <ABBR> or <word> filled in.
- **Open questions:**
  1. docs/spec/README.md:70-71 says an empty field is written as `—`. Does the tsv writer put `—` in an empty field by itself, or does internal/psb keep writing `—` for line 0 (check.go:44)? No rule of G2 to G5 can give an empty excerpt today, so the bytes are the same either way. The plan must say which package owns that rule.

### `psb-batch-api` — Hand the gap table from internal/psb to internal/setup through internal/cli for S01 and S06

- **What:** internal/cli reads inputs/briefs/problem-statement.md one time. It gives the same bytes to psb.Check and to internal/setup for the row `brief.sha256`. Then it hands the gap table to internal/setup, with each gap's Q-NNN ID, rule, line, excerpt and question. internal/setup does not import internal/psb (packages.md:43), and internal/cli does not import internal/tsv (packages.md:37). So one of two designs crosses the seam: cli hands the psb-gaps table bytes and setup reads them with internal/tsv against the psb-gaps schema, or cli converts the psb rows into a type of internal/setup. The plan chooses one.
- **Specification:** docs/spec/setup.md: The steps (S01, S06); docs/spec/setup.md: The stop table; docs/spec/psb-check.md: Not in phase 1 (the `id` is the question ID of the gap); docs/spec/packages.md: NFR-007 — The table of phase 1 (internal/cli, internal/psb, internal/setup); docs/architecture.md: 5. Intake and setup (the gap check, step 1)
- **Packages:** internal/psb, internal/cli
- **Present state:** IDs exist only in the output of WriteTSV (internal/psb/check.go:105-106, format 'Q-%03d'). Gap (check.go:16-21) has no ID field. internal/cli/cli.go:48-54 reads the file and calls psb.Check and psb.WriteTSV, and passes nothing to another package. cli can already write the table into a buffer and hand the bytes, so a change to internal/psb is not necessary. internal/setup does not exist.
- **Tests:**
  - unit: For the Go literal of triggers.md, the handed batch holds the IDs Q-001 to Q-007 in row order, and two calls on the same text give equal batches (the IDs stay stable for S06).
  - integration: internal/cli with the real internal/psb on a real file: the IDs and the question texts that cli prints for `layup psb check` equal those that cli hands on, and the bytes that were hashed equal the bytes that the rules read (setup.md:86). The S01 stop-table test and the S06 fact test belong to the setup items, which depend on this item.
- **Open questions:**
  1. Which type crosses the seam: the table bytes read by internal/setup with internal/tsv, or a type of internal/setup that cli fills? setup.md:86 says only 'hands the gap table to `internal/setup`, as it does for the step checks'.
  2. docs/spec/setup.md:114 gives the stop table's `where` column as '`—` otherwise' for a Q-NNN gap. Then the idea owner sees no line and no excerpt. G4 gives the same question text for each line that has `fast` (psb-check.md:39), so two G4 rows in the stop table cannot be told apart. Should `where` hold `<line> <excerpt>` for a gap?
  3. S06 needs the question text of each Q- answer (setup.md:91). A resumed run skips S01 when S01 has a done row (setup.md:57-60). No line says whether cli runs the rules again on each run, or whether S01 persists the batch. Running them again is safe only while the brief hash (S06) and the rules both stay the same.
  4. The format 'Q-%03d' (check.go:106) gives Q-1000 for a 1000th gap, which breaks the type id(Q-NNN) (psb-check.md:49). F-0004:9 says 'the question IDs have three digits'. The spec gives no cap and no exit code for more than 999 gaps.

### `psb-cli-contract` — The command `layup psb check FILE` end to end: exit codes, streams, determinism

- **What:** The built binary runs `layup psb check docs/facts/problem-statement-brief.md`, prints internal/psb/testdata/psb.tsv byte for byte and exits 1. Two runs print the same bytes. The binary gives the same bytes with an empty environment and a closed stdin. A usage error or a file that cannot be read exits 2, with stdout empty and a diagnostic on stderr. The doc comment of Run (cli.go:24-25) names exit 1. The usage line can also say that exit 2 covers a file that cannot be read. The cli tests that create files move to the integration level.
- **Specification:** docs/spec/psb-check.md: The command; docs/spec/README.md: Commands; docs/spec/packages.md: NFR-007 — The table of phase 1 (cmd/layup, internal/cli)
- **Packages:** internal/cli, cmd/layup
- **Present state:** internal/cli/cli.go:43-59 implements the command. Not exactly `check FILE` gives exit 2 (cli.go:44-46). A read error gives exit 2 (cli.go:48-52). Gaps give exit 1. The usage text at cli.go:16-22 already says 'exit 0 with no gap, 1 with gaps, 2 on a usage error'. Only the doc comment of Run (cli.go:24-25) omits exit 1. internal/cli/cli_test.go:48-70 covers 0, 1 and 2 with temp files (helpers_test.go:10). cmd/layup/main_e2e_test.go:12-24 runs only `layup version`. cmd/layup/main.go:10-12 holds one line of wiring and no logic, so the e2e tests cover it. The e2e test reads internal/psb/testdata/psb.tsv where it is: PRD-0001-layup.md:114 names that path, and tests/README.md:5-7 puts the first e2e test beside its package.
- **Tests:**
  - e2e: The built binary on docs/facts/problem-statement-brief.md gives psb.tsv byte for byte and exit 1. This is the REQ-001 acceptance criterion as written (PRD-0001-layup.md:114). Today only TestGoldenRealPSB calls Check and WriteTSV, and it does not run the command.
  - e2e: Two runs of the binary on the same file give identical stdout (README.md:61-63).
  - e2e: A clean file gives the header row only and exit 0.
  - e2e: The binary run with an empty environment and a closed stdin gives the same bytes and exit code as a normal run (README.md:44-45).
  - unit: `psb check` with no file, `psb` alone, and `psb check a b` each give exit 2, an empty stdout and a message on stderr. These cases touch no file. Today the test does not check that stdout is empty.
  - integration: A missing file and a directory each give exit 2, an empty stdout and a message on stderr. A clean temp file gives exit 0 and an empty stderr. These tests move from the untagged level (cli_test.go:48-70 writes temp files) to -tags=integration.
- **Open questions:**
  1. docs/spec/README.md:42-43: '`layup version` is exempt: it prints one line, `layup <version>`, and exits 0'. Today `layup version extra` prints the line and exits 0 (cli.go:32-34). README.md:55 makes an unknown flag a usage error (exit 2). Does the exemption accept extra words?
  2. docs/tests/test-levels.md:48 says 'every component has at least one unit test'. cmd/layup has none (`cmd/layup [no test files]`), and main.go holds no logic of its own. Is the e2e test sufficient for cmd/layup?
  3. No spec section gives the usage text of the whole `layup` command. Other areas add `setup`, `setup verify` and `gate` lines to the same constant (cli.go:16-22). Who owns that text?

### `psb-rule-edges` — Add the unsettled edge cases of G1 to G5 to the rule table and to edge.md

- **What:** psb-check.md:42-44 requires that 'such a case found later is added here and to internal/psb/testdata/edge.md'. This item writes a row or a sentence for each edge case below into the rule table of docs/spec/psb-check.md, and adds a case to edge.md and edge.tsv that pins the present behaviour. If the Operator changes a behaviour (for example G1), the code changes in the same item. psb.tsv must not change.
- **Specification:** docs/spec/psb-check.md: The rules
- **Packages:** internal/psb
- **Present state:** triggers.md and triggers.tsv cover G1 to G5 (internal/psb/testdata/triggers.tsv). edge.md:3-12 covers a header-less data row and a short row. clean.md covers no gap. The rule table at psb-check.md:34-40 settles none of the cases above.
- **Tests:**
  - unit: edge.md and edge.tsv (embedded) pin these cases. G1: `**Technology stack:**` and `Technology stack: *` with no value give no G1 gap (the regex at check.go:25 backtracks, so a `*` is the value character); a tab between the words and `:` gives a G1 gap. G2: a `Measurement` header in bold is not found (check.go:147); with both Measurement and Verification, the first column wins (check.go:146-150); an escaped pipe `\|` splits a cell (check.go:124). G3: a heading `terms` in lower case is not a terms heading (check.go:179); a heading that only contains `Terms`, such as `## Payment Terms`, makes the terms table, and a later `## Terms` table is ignored, so its abbreviations become G3 gaps (I ran it: API gives Q-001 G3 on line 5); a UTF-8 byte order mark before a first-line `# Terms` stops the heading, so G3 gives no gap (I ran it: exit 0 with an undefined SLA); an abbreviation inside the terms table itself. All rules: a `#` line inside a code fence is read as a heading (check.go:175), and G4 applies inside a fence (I ran it: `the cache is fast` in a fence gives a G4 gap).
  - integration: TestGoldenRealPSB still gives psb.tsv byte for byte. The real PSB has no `technology stack:` line, so a G1 change does not move a row.
- **External inputs:**
  - A decision of the Operator for the G1 case `**Technology stack:**` with no value: keep it (no gap), or change the regex to `[ *]*[^\s*]` so it is a G1 gap. It is the usual form of an empty template label.
- **Open questions:**
  1. docs/spec/psb-check.md:36 (G1) says 'then only spaces or `*`, then `:`, then only spaces or `*`, then a character that is not white space'. A `*` is a character that is not white space, so the table gives two readings of `**Technology stack:**`. The code reads it as a named stack.
  2. docs/spec/psb-check.md:36 says 'only spaces'. The regex at check.go:25 uses a literal space, so a tab is not a space. Is that the intent?
  3. docs/spec/psb-check.md:38 (G3) says 'a heading that contains `Terms`'. Does a UTF-8 byte order mark before `#` stop a heading? Does a `#` line inside a code fence count as a heading?

### `psb-rule-values` — Pin each value of the rule table with a golden case

- **What:** A golden case for each value that the rule table names, so that a regression in one value fails a test. The output of the code does not change, and psb.tsv does not change.
- **Specification:** docs/spec/psb-check.md: The rules
- **Packages:** internal/psb
- **Present state:** triggers.md pins only an empty cell and `Not measured.` for G2 (triggers.md:9-10), only `fast` for G4 (triggers.md:3), and one G5 row (triggers.md:9). No golden has a `Verification` header. The critic ran these cases, and the code obeys the table.
- **Tests:**
  - unit: An embedded golden (for example testdata/values.md and its .tsv) pins: a `Verification` header (G2); `—`, `-`, `tbd`, `TBD.` and `NOT MEASURED` in a measure column (G2, any case, one final `.` ignored); each of robust, soon, clean, better and handle in mixed case (G4, the question word in lowercase); two G4 words on one line (only the first one, psb-check.md:39); a G4 word on a line with a digit (no gap); `(Start Value)` (G5, any case); `TECHNOLOGY STACK: Go` (G1, any case).

### `psb-utf8-input` — A rule and a test for a FILE that is not valid UTF-8

- **What:** `layup psb check` on a file that holds bytes that are not valid UTF-8 never writes invalid bytes into the psb-gaps table. The spec gets one rule: exit 2 with a diagnostic (an input that does not match), or each invalid byte written as U+FFFD. The code and a test follow that rule.
- **Specification:** docs/spec/psb-check.md: The command; docs/spec/README.md: Commands (exit code 2); docs/spec/README.md: Records (a record is in UTF-8)
- **Packages:** internal/psb, internal/cli
- **Present state:** internal/cli/cli.go:48-53 passes the file bytes to psb.Check with no check. I ran it: a line with the bytes 0xff 0xfe gives those bytes in `excerpt` and exit 1. psb-check.md:21 says '`FILE`: a problem statement in Markdown, UTF-8'. README.md:67 says 'A record is a table of tab-separated values in UTF-8'.
- **Tests:**
  - unit: The Go literal "Technology stack: Go\n\xff\xfe fast\n" gives the chosen result: exit 2 with an empty stdout, or an excerpt with U+FFFD. Today the raw bytes 0xff 0xfe reach stdout and the exit is 1.
  - integration: With the chosen rule, utf8.Valid is true for the stdout of every golden and of the invalid case.
- **External inputs:**
  - A decision on the rule: exit 2 or U+FFFD. README.md:31-34 says a value that needs a decision of the Operator is never set by the spec; if the plan sets it, it is marked 'decided here' with its reason.
- **Open questions:**
  1. docs/spec/README.md:55 gives exit 2 for 'an input that does not match its schema', but a problem statement has no schema block. psb-check.md:21 says the input is UTF-8 and gives no exit code for one that is not. Two readings: exit 2, or a repair to U+FFFD.

### `later-psb-meaning-review` — The review of meaning: a session writes the gaps of meaning as typed rows

- **What:** A review session reads the problem statement and writes typed rows to records:intake/meaning.tsv: a byte-exact quote, the kind and the question. Code checks each quote against the file and refuses a row whose quote is not in it. Phase 2, under `layup run`.
- **Specification:** docs/spec/psb-check.md: Not in phase 1; docs/spec/records.md: the row intake/meaning.tsv; docs/architecture.md: 5. Intake and setup (the gap check, step 2)
- **Packages:** internal/run, internal/session, internal/records
- **Present state:** none. docs/spec/records.md:36 lists intake/meaning.tsv as phase 2, 'later'.
- **Tests:**
  - unit: The quote check refuses a quote that is not byte-exact in the file and accepts one that is.
  - integration: The rows that a session writes validate against the schema of intake/meaning.tsv and go to the records branch through internal/records.
- **External inputs:**
  - The harness and the role routing of architecture.md §9 (phase 2).

### `later-psb-intake-batch` — One batch on the Intake issue, with the answers as a raw fact

- **What:** Code merges the rule gaps, the meaning gaps and the setup questions. It removes a question whose quote and kind repeat, gives each question an ID, posts one Intake comment, copies the answers, and asks again only for the missing or open answers. Phase 2, under `layup run`. In phase 1, setup S01 and S06 do this with answers.tsv.
- **Specification:** docs/spec/psb-check.md: Not in phase 1; docs/architecture.md: 5. Intake and setup (the gap check, steps 5 and 6)
- **Packages:** internal/run, internal/forge, internal/records
- **Present state:** none. For LAYUP itself, F-0004 holds the answers to Q-001 to Q-019 (docs/facts/F-0004-psb-gap-answers.md:9).
- **Tests:**
  - integration: A session's rows and the rule gaps merge into one table. The merge removes a question whose quote and kind repeat, and every ID gets exactly one answer line.
  - e2e: On a pilot target, one Intake comment holds both kinds of gap, and the answers become a raw fact.
  - uat: On a pilot problem statement, the idea owner sees the gaps of meaning beside the rule gaps in the one Intake batch (PRD-0001-layup.md:114). This rides on the e2e path above.
- **External inputs:**
  - A pilot repository and its problem statement (PRD-0001-layup.md:114).

### `later-psb-early-question-share` — The rows 'before delivery' in questions.tsv and the measure of the Early Question Share

- **What:** Each Intake question and each follow-up is a row of the project question table records:questions.tsv with the flag 'before delivery' (phase 3). `layup report` computes the Early Question Share (human questions before delivery over all human questions) and compares it with its start value F-0004#19 (phase 4). This is the last clause of the REQ-001 acceptance criterion.
- **Specification:** docs/spec/records.md: the row questions.tsv; docs/architecture.md: 5. Intake and setup (the gap check, step 6); docs/architecture.md: the measure Early Question Share (§12)
- **Packages:** internal/records, internal/report
- **Present state:** none. docs/spec/records.md:37 gives questions.tsv to phase 3. docs/spec/packages.md:90 gives internal/report to phase 4. docs/architecture.md:375 and :1322 define the flag and the measure.
- **Tests:**
  - unit: The measure gives the share of rows flagged 'before delivery' among all rows where a human was asked, with its population.
  - integration: The rows that layup run writes for the Intake questions carry the flag 'before delivery' until the first task starts, and layup report reads them.
  - e2e: In the pilot, layup report prints the Early Question Share beside its start value (PRD-0001-layup.md:114, 169).
- **External inputs:**
  - The pilot and its baseline measurement (PRD-0001-layup.md §9, phase 4).

### Constraints of area `psb`

1. **docs/prd/PRD-0001-layup.md:114** — `layup psb check docs/facts/problem-statement-brief.md` writes the batch `internal/psb/testdata/psb.tsv` byte for byte (`TestGoldenRealPSB`), and the idea owner's answers to it are the raw fact `F-0004`, one fact per question *Effect:* Every psb item must keep psb.tsv byte for byte the same. A change to a rule or to the writer that moves a row breaks the REQ-001 acceptance criterion and the mapping of F-0004 facts 1-19 to Q-001 to Q-019.
2. **docs/spec/README.md:67-71** — A field holds no tab, no line feed and no carriage return: a writer replaces each with one space. An empty field is written as `—` (U+2014) *Effect:* The lone-CR fix is phase-1 work. The field rule must have one owner: it lands in internal/psb first (psb-field-rule-cr), and psb-tsv moves it to internal/tsv.
3. **docs/spec/README.md:75-77** — A Go test of the code reads these blocks and compares them with what the code writes and reads (the test comes with the code, Bootstrap mode rule 1) *Effect:* The schema-block test belongs to psb-tsv. It cannot be a separate later issue.
4. **docs/spec/packages.md:37,40,43** — `internal/cli` may import `internal/psb`, `internal/setup`, `internal/verify`, `internal/gate`; `internal/psb` may import `internal/tsv`; `internal/setup` may import `internal/tsv`, `internal/git`, `internal/catalog` *Effect:* setup may not import psb, and cli may not import tsv. The hand-off of psb-batch-api crosses the seam as table bytes or as a type of internal/setup. psb-batch-api comes before the S01 item of the setup area.
5. **docs/spec/packages.md:31-32** — An import that the table does not allow is a defect, and the boundary rule of the future Go gate of LAYUP reads this table. *Effect:* No psb item adds its own import test. The gate area owns the import check.
6. **docs/spec/packages.md:46-47** — `internal/psb` today imports no package of this module and writes its table itself; it moves to `internal/tsv` when that package exists (#29). *Effect:* psb-tsv depends on the foundation item that creates internal/tsv. The same PR removes this note (AGENTS.md gate step 7).
7. **docs/spec/psb-check.md:23-24** — the present code does not yet replace a lone carriage return in `excerpt` with a space (#29 does) *Effect:* The PR of psb-field-rule-cr deletes this sentence.
8. **docs/spec/psb-check.md:31** — Each rule reads text only; none calls a model (`NFR-005`). *Effect:* internal/psb stays pure: text in, rows out, no file and no network. packages.md:22-23 forbids net imports in phase 1.
9. **docs/spec/psb-check.md:42-44** — The present code is the reference for an edge case that this table does not settle; such a case found later is added here and to `internal/psb/testdata/edge.md`. *Effect:* psb-rule-edges is phase-1 work and not optional. Each edge case found gets a row in the rule table and a golden case in edge.md. A change of behaviour is never silent.
10. **docs/spec/setup.md:86,91** — A record row `brief.sha256` (`computed`) holds the SHA-256 of the problem statement that the rules read. … Refuses a problem statement whose SHA-256 differs from the row `brief.sha256` of S01 (exit 2), so the `Q-NNN` IDs stay true *Effect:* The rules must be deterministic and the order of the batch stable across runs. The bytes that are hashed must be the bytes that internal/psb reads, so the file is read one time.
11. **docs/spec/README.md:44-45,61-63** — It reads no environment variable for an input, and asks no question on the terminal. … Two runs on the same input print the same bytes (`NFR-005`). *Effect:* The psb e2e tests check two runs for identical output, and a run with an empty environment and a closed stdin.
12. **docs/tests/test-levels.md:43** — Unit tests touch no file, network, or clock, so they are fast and deterministic and run first, on every commit. *Effect:* The psb goldens at the unit level use //go:embed or Go literals. The cli tests that create temp files move to -tags=integration. A unit test stays in each package.
13. **docs/tests/test-levels.md:48,60,74** — every component has at least one unit test … every interface or workflow has an integration test … every user-facing scenario has an E2E test *Effect:* psb-cli-contract needs e2e tests of `layup psb check`. Today the e2e level covers only `layup version`.
14. **docs/tests/template-e2e.md:22** — One scenario per test. *Effect:* The e2e tests of psb-cli-contract are separate tests: the golden, the determinism, the clean file, and the empty environment.

### Risks of area `psb`

- The question IDs depend on the rules as well as on the brief. S06 checks only the brief's SHA-256 (setup.md:91). If LAYUP is upgraded between S01 and S06 of a resumed setup and a rule changes, the Q-NNN IDs in answers.tsv can point at other questions, and nothing detects it. The setup record holds no engine version for S01.
- Moving the writer to internal/tsv can change bytes by accident, for example through the handling of `—`, the final LF or a trimmed field. That breaks TestGoldenRealPSB, the REQ-001 acceptance criterion and F-0004's mapping.
- A change of the G1 regex after an Operator decision changes the result on a target's problem statement that holds an empty template label. psb.tsv does not change, but an earlier answers.tsv of a target can lose its Q-NNN alignment.
- Several areas edit the same usage constant and the dispatch in Run (internal/cli/cli.go:16-40), so parallel issues can conflict.
- Moving the file-reading cli tests to -tags=integration changes what the pre-commit hook runs. The hook then runs fewer psb tests unless unit tests with embedded data replace them.
- More than 999 gaps give Q-1000, which breaks the id(Q-NNN) type. This is unlikely for a real brief, but the spec has no limit.

### Critic points that the reviser of `psb` rejected

- psb-cli-contract tests (c): move the e2e fixture psb.tsv into the root tests/ directory. *Basis:* PRD-0001-layup.md:114 names internal/psb/testdata/psb.tsv as the batch of the acceptance criterion, so the e2e test reads that path. tests/README.md:5-7 says 'the first end-to-end test sits beside its package (`cmd/layup`)'. The golden belongs to internal/psb and is not a cross-package fixture. The item now gives this reason.
- cmd/layup needs a unit test of its own (test-levels.md:48). *Basis:* cmd/layup/main.go:10-12 holds one line of wiring and no logic, and packages.md:36 gives it only the job 'passes the arguments to internal/cli and exits with its code'. The e2e tests drive it. I did not add a unit test item. psb-cli-contract keeps this as an open question for the plan.
- The field rule must give a rule for a vertical tab, a form feed, U+0085 or U+2028 (as work in psb-field-rule-cr). *Basis:* docs/spec/README.md:69-70 names only tab, LF and CR, and no source gives a value for the other characters. Phase-1 work cannot set one. It stays an open question in psb-field-rule-cr, with no code change.

## Area `gate`

The final gate inventory has eleven items. Nine are phase 1 and two are later-phase.  Changes from the first inventory: - I removed three items and moved their tests into the items that hold the code. gate-nfr004-fixture and gate-repeat-rule went into gate-results and gate-command. gate-catalog-fixture-ci went into gate-catalog-go-entry. Each was a test written after its code, and gate.md:144-146 and setup.md:301-302 say that the test "comes with the code" (AGENTS.md gate step 3). - I added gate-run. It is the exported entry point of internal/gate: it resolves the revisions, reads the base manifest, runs the kinds and writes the table gate-result. packages.md:37 lets internal/cli import only internal/psb, internal/setup, internal/verify and internal/gate, so internal/cli cannot call internal/git or internal/tsv. internal/verify calls the same entry point. - I added gate-verify-kind for the row gate:<kind> of `layup setup verify` (setup.md:256). If the verify area has the same item, merge the two. - I added two later items: later-gate-status (phase 2) and later-gate-rule-batch (phase 2). They come from the "Not in phase 1" lists (gate.md:94-102, 125-127; setup.md:312). - I applied these critic points: NFR-002 on the target workflow, a test for a pending kind with a missing tool, peeling a revision to a commit, no exit 0 when git or the scratch tree fails, no skip filter in a job, the Go case of a product path with no package, hooks off for `git worktree add`, and the catalog lookups that internal/setup and internal/verify need.  I measured these facts with go1.27.1 at 7cdd346: - `go list -deps -test ./...` names os/exec, so the NFR-005 import check cannot use -test for rule 4. - `go vet ./...` and `go test -count=1 ./...` exit 1 with "no packages to vet" and "no packages to test" on a tree whose only .go file is in testdata/. - `//go:embed` fails on a tree that holds a go.mod, and without the `all:` prefix it leaves out `.github/`.  No code of this area exists today. internal/cli has only `version` and `psb` (internal/cli/cli.go:31-40), and go.mod has no require line.  External inputs: - Each tool version with its documentation (Invariant 4): Go, golangci-lint/depguard, actions/checkout and actions/setup-go. - The coverage floor (L-B2). - The toolchains on the host (L-B1). - A pilot repository on GitHub for the UATs.

| Key | Title | Requirements | Size | After | Host |
| --- | ----- | ------------ | ---- | ----- | ---- |
| `gate-manifest` | Gate manifest reader and scope patterns | REQ-004, NFR-004 | small | found-tsv (other area: the internal/tsv reader with schema and type checks), found-schema-test (other area: the test that compares each tsv-schema block with the code) | row 5 |
| `gate-scratch-tree` | Scratch work tree of the head with the base's gate files | REQ-004, REQ-007, NFR-004 | small | found-git (other area: internal/git with worktree add/remove, show, rev-parse, diff --name-only), gate-manifest | row 5 |
| `gate-results` | The result rules, the command run and the stderr blocks | REQ-004, NFR-004 | small | gate-manifest, gate-scratch-tree | row 5 |
| `gate-run` | internal/gate entry point and the table gate-result | REQ-004, REQ-007, NFR-004, NFR-005 | small | gate-manifest, gate-scratch-tree, gate-results, found-tsv (other area: the writer that replaces tab and line feed and writes `—`), found-git (other area) | row 5 |
| `gate-command` | The command `layup gate REPO --base REV --head REV` | REQ-004, REQ-007, NFR-004, NFR-005 | small | gate-run, setup-steps (other area: a target set up by `layup setup` S01-S15, for the UAT only) | row 5 |
| `gate-nfr005-imports` | NFR-005 import rule: no network package in phase 1 | NFR-005, NFR-007 | small | found-import-rules (other area: the check of the package table of packages.md, rules 2 to 4; this item can merge into it) | row 2 |
| `gate-catalog-package` | internal/catalog: the embedded stack catalog | REQ-002, REQ-004, NFR-003 | small | found-tsv (other area), found-schema-test (other area) | row 4 |
| `gate-catalog-go-entry` | The Go entry: kinds.tsv, files, fixtures and the fixture test | REQ-004, REQ-007, REQ-002, NFR-003 | small | gate-catalog-package, gate-command, gate-manifest | row 14 |
| `gate-catalog-go-workflow` | The target's CI workflow: one job per gate kind | REQ-007, REQ-004, NFR-004, NFR-002 | large | gate-catalog-go-entry, gate-command, setup-s12 (other area: S12 writes the entry files and refuses to overwrite a baseline path) | row 14 |
| `gate-verify-kind` | The check gate:<kind> of `layup setup verify` | REQ-002, REQ-004, NFR-004 | small | gate-run, gate-catalog-package, gate-catalog-go-entry, verify-runner (other area: the check runner and the table setup-verify) | row 15 |
| `later-gate-status` | Phase 2: the status `layup/gates` and the ready/merge rule | REQ-004, REQ-007 | large | gate-run, later-run (other area: `layup run` and the forge adapter, phase 2) | `M2e` |
| `later-gate-rule-batch` | Phase 2: activation and the rule batch's own gate files | REQ-003, REQ-004 | large | gate-run, later-gate-status, later-rules (other area: rule batches of REQ-003, phase 2) | `M2f` |

### `gate-manifest` — Gate manifest reader and scope patterns

- **What:** internal/gate reads the bytes of docs/gates.tsv against the schema gate-manifest: the header, the field count, kind as id(<word>) and unique, state active|pending, scope list(text), config list(path) or `—`. It gives the scope matcher and the product paths of a kind. A missing or malformed manifest is a typed input error (exit 2 at the command).
- **Specification:** docs/spec/gate.md: The gate manifest; docs/spec/README.md: Records; docs/spec/README.md: The types; docs/spec/gate.md: NFR-004 — A check that is not active is not passed (item 4); docs/architecture.md: 6. Native gates and rule protection / The stack gates
- **Packages:** internal/gate, internal/tsv
- **Present state:** none. No internal/gate or internal/tsv exists (docs/spec/packages.md:38, 42). internal/psb writes its table itself (docs/spec/packages.md:46-47).
- **Tests:**
  - unit: A valid manifest parses into rows in file order. Each malformed form gives a typed input error: a wrong header, a wrong field count, a bad state, a bad kind, a duplicate kind, a `..` in a config path, an empty field in place of `—`.
  - unit: Table tests of the scope patterns: `./*.go` matches `a/b/c.go` and `x.go`; `internal/*.go` matches `internal/a/b.go` and not `cmd/x.go`; a pattern with no `*` matches by prefix; `go.mod` is not a product path of `./*.go`.
  - integration: The gate-manifest block in docs/spec/gate.md has the same columns and types as the reader (the schema-block test).
- **Open questions:**
  1. docs/spec/gate.md:38: "A pattern with no `*` matches each file whose path starts with it". Is this a string prefix or a path-segment prefix? For example, does `internal` match `internal2/x.go`?
  2. docs/spec/gate.md:36-38 defines only `P/*.E` and a pattern with no `*`. A pattern such as `*.go`, `P/*` or `P/*/x.E` has no rule. Is it a schema error (exit 2) or a literal?
  3. docs/spec/gate.md:22 and README.md:117: `<word>` is "lowercase letters and `-`", so a kind with a digit is not valid. Confirm.

### `gate-scratch-tree` — Scratch work tree of the head with the base's gate files

- **What:** internal/gate makes a detached scratch work tree of the head (`git worktree add --detach`) outside REPO's working files. It writes the base's docs/gates.tsv and each config path of each base row into the tree, from `git show <base>:<path>`, and it removes each config path that the base does not have. So the head's own gate files never judge the head (FT4). It removes the tree and its worktree metadata on every exit path. Proposal (open question): it runs these git commands with hooks off, as architecture §4 does for `layup run`, so that a post-checkout hook of REPO cannot run. A failure of git, the add, the overlay or the removal is a typed error that never maps to exit 0.
- **Specification:** docs/spec/gate.md: The command (The run, steps 2 and 4); docs/spec/README.md: Commands (exit codes); docs/architecture.md: 6. Native gates and rule protection / The stack gates (FT4); docs/adr/0016-put-the-native-stack-gates-in-the-target.md: Decision 3
- **Packages:** internal/gate, internal/git
- **Present state:** none. internal/git does not exist (docs/spec/packages.md:39 gives its job).
- **Tests:**
  - unit: With stand-ins for git and the file system, the overlay plan writes each base config path and the base manifest, and removes a config path that is absent at the base. A head edit of a config file or of docs/gates.tsv is replaced. The unit test touches no real file (test-levels.md:43).
  - integration: With a real repository in t.TempDir, the head changes a gate test and the manifest, and the scratch tree holds the base bytes. After the run, `git worktree list` shows only the main tree, the scratch directory is gone, and the refs, the index and `git status` of REPO are unchanged. Removal also occurs when a command fails or the overlay fails.
  - integration: REPO has a post-checkout hook (through a relative core.hooksPath) that writes a marker file. After the run, the marker file does not exist. This test applies only if the hooks-off proposal is accepted.
  - integration: When `git worktree add` fails (for example, a head that is not a commit) or the overlay write fails, the error is typed, the scratch tree is removed, and the mapped exit code is not 0.
- **Open questions:**
  1. docs/spec/gate.md:59 says only "Make a scratch work tree of the head (`git worktree add --detach`)". It does not say that hooks are off. architecture.md:263-264 turns hooks off only for `layup run`'s fetch ("`core.hooksPath` set to an empty directory"). Is hooks-off decided for `layup gate`?
  2. docs/spec/gate.md:52-54 gives exit 2 for usage and input errors only. README.md:55-58 says "A check that did not run is never 0" and does not give a code for `git` not found or a failed `worktree add`. Is it 2, 1 (not-active rows), or a new code?
  3. docs/spec/gate.md:47-48: "changes nothing in it except a scratch work tree, which it removes before it exits". The spec does not say what occurs on SIGINT or SIGTERM, or where the scratch directory lives. A killed run leaves a registered worktree. Is a `git worktree prune` at the start allowed?
  4. docs/spec/gate.md:27 types config as `list(path)`. Can an item be a directory, so that the overlay copies and removes a whole subtree?
  5. docs/spec/setup.md:256: if a fixture patch touches a config path, the base overlay reverts it, and the fixture does not test that change. Is this a rule for fixture authors?

### `gate-results` — The result rules, the command run and the stderr blocks

- **What:** For each base row, in manifest order and one at a time, internal/gate gives one result. The first matching line of the table decides: 1. pending, with no product path in `git diff --name-only <base> <head>`: clear. 2. pending, with a product path: fail, with the first path. 3. active, with the tool not found by exec.LookPath: not-active. 4. active, with no product path in the scratch tree: clear. 5. The command, run with `sh -c` in the scratch tree, exits 0: pass. 6. The command exits with another code: fail, `exit <code>`. It is killed by a signal: fail, `signal <name>`. A pending command never runs. The command's output goes to stderr in one block per kind, with the kind as its heading. This item holds the NFR-004 tests: no input gives pass for a check that did not run.
- **Specification:** docs/spec/gate.md: The command (The run, step 3, and the result table); docs/spec/gate.md: The table (the stderr blocks); docs/spec/gate.md: NFR-004 — A check that is not active is not passed (items 1 to 3, 5); docs/architecture.md: 6. Native gates and rule protection / The stack gates (the table of layup/gates); docs/architecture.md: 15. Known limits (L-B1)
- **Packages:** internal/gate, internal/git
- **Present state:** none.
- **Tests:**
  - unit: With stand-ins for the runner, LookPath and the diff, a table test covers each of the six lines with the exact reason strings. Order wins: an active kind with a missing tool and no product path gives not-active, not clear. A pending row never calls the runner (zero recorded calls).
  - unit: An NFR-004 property table over all rule inputs: no input with a missing tool or a pending state gives pass. A pending row with a missing tool gives the pending result (clear or fail), never pass. This follows the table order; see the open question on gate.md:140-141.
  - integration: With a real `sh -c`: `true` gives pass with `—`; `exit 3` gives fail with `exit 3`; `kill -TERM $$` gives fail with `signal <name>`. The working directory is the scratch tree. Stderr holds one block per kind, headed by the kind, in manifest order. A real `git diff` on a rename and on a delete of a .go file gives the pending fail.
- **Open questions:**
  1. docs/spec/gate.md:140-141 says "A missing toolchain is `not-active` (`tool not found`), never `pass` or `clear`". The table (gate.md:69-71) decides a pending row before it checks the tool, so a pending kind with a missing tool gives clear. My reading is that item 3 applies to active rows, because a pending command never runs (gate.md:76). Confirm this reading.
  2. docs/spec/gate.md:74 gives "`signal <name>`" with no form of the name. Go's Signal.String() gives `terminated`, not `SIGTERM`.
  3. docs/spec/gate.md:69-70 uses `git diff --name-only <base> <head>`. With rename detection, a rename out of the scope shows only the new name. diff.renames and core.quotePath of the host also change the output. Is `--no-renames -z` (or an equivalent) decided for internal/git?
  4. docs/spec/gate.md:70: "`<first path>`". Is this the first path in git's output order (byte order)?
  5. docs/spec/gate.md:71 covers only a missing `tool`. If a second program is missing (`gofmt` beside `go`), sh exits 127 and the result is fail, not not-active. Is fail the intended result?
  6. No timeout is set for a gate command. A `go test` that hangs blocks `layup gate` and `layup setup verify` with no limit. Is a timeout decided, and which result does it give?
  7. The spec does not say which environment the command gets. GOFLAGS, GOTOOLCHAIN or GOPATH of the host can change a verdict (README.md:44-45 forbids an environment variable as an input to layup). Does the command inherit the environment, or get a clean one?

### `gate-run` — internal/gate entry point and the table gate-result

- **What:** One exported function of internal/gate takes REPO, base and head. It does these steps: 1. It peels each revision to a commit (`<rev>^{commit}`) through internal/git. 2. It reads the base manifest. 3. It runs the scratch tree and the result rules. 4. It returns the rows, the stderr blocks and a typed error. It writes the table gate-result through internal/tsv with full 40-character commit SHAs. internal/cli and internal/verify both call it. A revision that does not resolve to a commit, and a missing or malformed base manifest, are typed input errors. The table holds no time and no scratch path, and its rows follow the manifest (the repeat rule).
- **Specification:** docs/spec/gate.md: The command (REPO, --base, --head; The run, steps 1 and 4); docs/spec/gate.md: The table; docs/spec/gate.md: REQ-007 — Each change passes the gates before review or merge (item 3); docs/spec/README.md: Commands (Determinism); docs/spec/packages.md: The table of phase 1 (internal/gate, internal/cli, internal/verify)
- **Packages:** internal/gate, internal/git, internal/tsv
- **Present state:** none.
- **Tests:**
  - unit: With stand-ins for git and the rules, the rows go to the writer in manifest order with base and head as given. Each typed error (unresolved revision, missing manifest, malformed manifest) comes back with no rows.
  - integration: On a real repository: an annotated tag as --head gives the commit SHA, not the tag object's SHA. A revision that names a tree or a blob gives the typed input error. A base with no docs/gates.tsv gives the typed input error. The gate-result schema block matches the writer.
- **Open questions:**
  1. docs/spec/gate.md:50-51 says "any revision that `git rev-parse` resolves", and gate.md:83-84 types base and head as "the base commit" and "the head commit". The spec does not say that a revision is peeled to a commit. I read it as `^{commit}`, and a tree or a blob as an input error. Confirm.
  2. docs/spec/gate.md:142-143 makes a missing or malformed manifest exit 2. It does not cover a manifest with a valid header and zero rows. Is that an input error, or an empty table with exit 0 (the FT1 trap of passing on nothing)?

### `gate-command` — The command `layup gate REPO --base REV --head REV`

- **What:** internal/cli parses `layup gate REPO --base REV --head REV` and calls the gate-run entry point. It prints the table on stdout and the stderr blocks on stderr. It maps the result to an exit code: 0 when each row is pass or clear, 1 when a row is fail or not-active, and 2 on a usage error or a typed input error. In phase 1 it posts no status, takes no rule-batch or approval-hash input, and runs no known-bad fixture (gate.md:94-102). Its e2e tests hold the NFR-004 fixture and the repeat rule.
- **Specification:** docs/spec/gate.md: The command; docs/spec/gate.md: Not in phase 1; docs/spec/gate.md: NFR-004 — A check that is not active is not passed (items 4 and 5); docs/spec/gate.md: REQ-007 — Each change passes the gates before review or merge (item 3); docs/spec/README.md: Commands; docs/spec/packages.md: The table of phase 1 (cmd/layup, internal/cli)
- **Packages:** internal/cli, cmd/layup
- **Present state:** internal/cli/cli.go:31-40 dispatches only `version` and `psb`. The usage text at internal/cli/cli.go:16-22 has no `gate`. cmd/layup/main.go:10-12 exits with the code of cli.Run. cmd/layup/main_e2e_test.go is the e2e pattern.
- **Tests:**
  - unit: internal/cli, with a stand-in for the gate entry point: no REPO, a missing --base or --head, an unknown flag, or an extra argument gives exit 2 and a usage line on stderr only. The result-to-exit-code map is 0, 1 or 2, and a typed error never gives 0.
  - e2e: The built binary runs on a fixture repository with rows that give pass, fail, pending clear, pending fail and active clear. Stdout equals a golden table and the exit code is 1. The fixture commits have fixed author and committer dates and a fixed identity, or the golden table is built from the resolved SHAs (guardrails.md:276-280). With all rows pass or clear, the exit code is 0.
  - e2e: The NFR-004 fixture (gate.md:144-146): a manifest whose tool is `layup-no-such-tool-xyz` (active, with a product path) gives the row `not-active` with `tool not found: layup-no-such-tool-xyz`, and exit 1. A base with no manifest, or a revision that does not resolve, gives exit 2, never an empty table with exit 0.
  - e2e: The repeat rule: two runs on one fixture repository give byte-equal stdout and the same exit code. Stdout holds no path of the temporary directory or of REPO.
  - uat: The Operator runs `layup gate` on a pilot target that `layup setup` set up. The Operator reads one verdict per kind and finds the reason of a failure in the stderr block.
- **External inputs:**
  - A pilot target for the UAT. It is set up by the setup area, and the Operator pushes it.
- **Open questions:**
  1. docs/spec/gate.md:44 puts REPO before the flags. Go's flag package stops at the first non-flag argument, so a plain flag.FlagSet cannot parse `layup gate REPO --base X`. Is the order fixed (flags first), or is a hand parser expected?
  2. docs/spec/gate.md:142-143 does not cover a manifest with a header and zero rows. Keep that e2e assertion pending until there is a decision (see gate-run).

### `gate-nfr005-imports` — NFR-005 import rule: no network package in phase 1

- **What:** A Go test fails when a phase-1 package of the module imports net, net/http or crypto/tls, so the engine checks open no connection. It is the mechanical check of NFR-005 item 1 (packages.md rule 5).
- **Specification:** docs/spec/gate.md: NFR-005 — No model call in the engine checks (item 1); docs/spec/packages.md: NFR-007 — Go, the standard library only, and Git as the `git` program (rule 5)
- **Packages:** none (a test only; its home is open)
- **Present state:** none. At 7cdd346, `go list -deps -test ./...` names os/exec in every test binary (through testing's internal packages) and names no net, net/http or crypto/tls (measured, go1.27.1). So -test works for rule 5 and cannot work for rule 4 (os/exec).
- **Tests:**
  - integration: `go list -deps ./...` (it starts the Go toolchain and reads the module files) names none of net, net/http, crypto/tls. A fixture module outside ./... (its own go.mod under testdata/) that imports net/http makes the check fail.
- **Open questions:**
  1. docs/spec/packages.md:22: "No package of phase 1 imports `net`, `net/http` or `crypto/tls`". It does not say if this means direct or transitive imports, or if _test.go files count.
  2. docs/spec/packages.md:34-44 has no package for a repository-wide import test. Which package holds it?

### `gate-catalog-package` — internal/catalog: the embedded stack catalog

- **What:** internal/catalog embeds internal/catalog/<stack>/ and lists the stacks. It reads each kinds.tsv against the schema catalog-kinds. For an entry, it gives these: - the files, with `{{module}}` replaced - the fixture patch of a kind - the derived manifest (the kinds table without version, fixture and evidence) - a lookup that says if `<stack>/<path>` is a file of the entry (for the check `sources`) - the config paths of each kind (for the S15 rule-path register) An active kind must name an existing fixtures/<kind>.patch. A pending kind has `—` as its fixture.
- **Specification:** docs/spec/setup.md: The stack catalog; docs/spec/setup.md: The steps (S12, S15); docs/spec/setup.md: The setup record (ref `catalog`: `<stack>/<path>`); docs/spec/setup.md: The checks of `layup setup verify` (sources); docs/spec/packages.md: The table of phase 1 (internal/catalog, internal/setup, internal/verify); docs/architecture.md: 6. Native gates and rule protection / The stack gates
- **Packages:** internal/catalog, internal/tsv
- **Present state:** none. No internal/catalog exists.
- **Tests:**
  - unit: With an fs.FS test entry: the rows load in order. A malformed kinds.tsv errors. An active kind with no fixture file errors. The derived manifest has exactly the six gate-manifest columns in order, byte for byte. `{{module}}` is replaced in every file and in no other token. The `<stack>/<path>` lookup is true only for files of the entry. The config paths of each kind come back in row order.
  - integration: The catalog-kinds schema block matches the code. The real embedded FS holds `.github/workflows/*` and the go.mod template. The derived manifest is parsed by the gate-manifest reader in gate-catalog-go-entry's test, because internal/catalog may not import internal/gate (packages.md:41-42).
- **Open questions:**
  1. docs/spec/setup.md:285 puts each file "at `internal/catalog/<stack>/files/<path>`", and S12 writes `go.mod` (setup.md:97). I measured that `//go:embed` fails on a tree that holds a go.mod (a separate module), and that without the `all:` prefix it leaves out `.github/`. The catalog needs a rename rule (for example `go.mod.tmpl` to `go.mod`) and an `all:` pattern. The spec does not give one.
  2. A `.go` file under internal/catalog/go/files/ (for example a layout-test template) becomes a package of LAYUP's module for go vet, go test and go list -deps. Is a template suffix required for Go sources in files/?

### `gate-catalog-go-entry` — The Go entry: kinds.tsv, files, fixtures and the fixture test

- **What:** internal/catalog/go/ holds kinds.tsv with five kinds. static (`test -z "$(gofmt -l .)"` and `go vet ./...`) and test (`go test -count=1 ./...`) are active. layout, boundary and contract are pending. files/ holds the go.mod template with `{{module}}` and the tools' configuration. fixtures/static.patch and fixtures/test.patch each make their kind fail. In the same change, a Go test renders a temporary target from each embedded entry. It runs `layup gate` on the clean commit (each active kind gives pass or clear) and on a commit with each active fixture applied (that kind gives fail). This is LAYUP's own CI run of each fixture, in the existing `tests` job, so ci.yml does not change. The test fails when no entry exists.
- **Specification:** docs/spec/setup.md: The stack catalog (LAYUP's own CI runs each fixture); docs/spec/setup.md: The steps (S12); docs/architecture.md: 6. Native gates and rule protection / The stack gates (The Go entry; LAYUP's own CI runs every fixture); docs/architecture.md: 15. Known limits (L-B2); docs/adr/0016-put-the-native-stack-gates-in-the-target.md: Decision 1 and Decision 4; docs/walkthroughs/W-04-stack-dependent-gates.md: row 1
- **Packages:** internal/catalog, cmd/layup (e2e test of the fixtures)
- **Present state:** none. LAYUP's CI runs the same static commands (.github/workflows/ci.yml:302-303). Its test step is `go test -timeout 10m ./...` (ci.yml:314) with tags at :316-318, and no step uses `-count=1`. The CI `tests` job installs Go from go.mod (ci.yml:310-312), and go.mod:3 says `go 1.26`.
- **Tests:**
  - unit: The Go entry loads with exactly the kinds static, test, layout, boundary and contract. static and test are active with fixtures. The other three are pending with `—`. Each evidence is a URL, and each version of an active kind is not `—`.
  - e2e: Through the built binary, on a temporary target rendered from the entry: the clean commit gives pass or clear for static and test, and pending kinds give `clear` with `pending: no product path`. static.patch makes static fail, and test.patch makes test fail, each with exit 1. A fixture that does not apply fails the test. The derived manifest parses with the gate-manifest reader. With zero embedded entries, the test fails and does not pass on nothing.
  - integration: A new embedded entry is picked up with no change to the test: the test iterates over all embedded stacks.
- **External inputs:**
  - The Go version for the target's go.mod `go` directive and for the CI setup-go step, with its release documentation (Invariant 4; docs/spec/README.md:31-34). It must not be newer than LAYUP's CI Go (go.mod:3, `go 1.26`). If it is newer, go downloads a toolchain (GOTOOLCHAIN auto), so the fixture test uses the network, or fails with GOTOOLCHAIN=local.
  - The coverage floor of the test kind. It stays an open gap until the idea owner or the Operator sets it with evidence (architecture.md:1548-1551, L-B2; ADR-0016:85-86).
- **Open questions:**
  1. docs/spec/gate.md:25 gives the static command as `test -z "$(gofmt -l .)"` only, and architecture.md:462-463 adds `go vet ./...`. One row has one command. Is the command `test -z "$(gofmt -l .)" && go vet ./...`?
  2. docs/spec/setup.md:296: the fixture is "`—` for a `pending` kind until its activation". No value is given for tool, version, command and config of layout, boundary and contract before activation. Are they `—` at setup, or the planned tool?
  3. docs/walkthroughs/W-04-stack-dependent-gates.md:17 names only static as active. architecture.md:451-452 makes only the architecture-dependent kinds pending, so I read test as active at setup. Confirm.
  4. architecture.md:1550-1551: "the setup record lists the floor as an open gap". setup.md:254 (sources) needs each `gap` row's marker in the tree and its open-gaps.tsv row. But S10/S11 handle markers before S12 writes the catalog files, so no step writes this gap. A `&lsaquo;…&rsaquo;` marker inside internal/catalog/ also fails LAYUP's own check `markers` (docs/setup/setup-check.sh:305, MK_EXEMPT does not exempt it). Which step writes the floor's gap, and in which form?
  5. architecture.md:456-458: an active kind with no path in its scope gives clear, "for Go, `./...` that matches no package". The scope `./*.go` is not `./...`. A tree whose only .go files are in testdata/, or start with `_` or `.`, or are excluded by build constraints, has a product path and no package. I measured `go vet ./...` and `go test -count=1 ./...` exit 1 ("no packages to vet", "no packages to test"), so the kind gives fail, not clear. Is this a known case, or does the scope or the command change?
  6. docs/spec/setup.md:286: a fixture is applied "to the setup head", and the baseline is the latest commit at setup (O-101, setup.md:340-341). Must each fixture add only new files, or touch only files that the entry writes, so that it applies to any baseline?

### `gate-catalog-go-workflow` — The target's CI workflow: one job per gate kind

- **What:** internal/catalog/go/files/.github/workflows/<file>.yml has one job per kind, with the kind as the job name. Each job reads docs/gates.tsv of the pull request head and decides in its own step, with no `paths:` filter and no job-level `if:`. A job that is skipped reports success. A missing manifest, or no row for the job's kind, fails the job. The jobs apply the same rules as gate-results, in shell: - A pending kind fails when the diff from the pull request base changes a product path. It passes with `pending: no product path` when it changes none. - An active kind with no product path passes with `no product path`. - Other active kinds run the manifest command. The workflow starts no `layup` binary and fetches no LAYUP file (NFR-002). Its file name differs from each baseline workflow, and the baseline's own workflows stay byte for byte. The jobs of all kinds exist from the setup.
- **Specification:** docs/spec/gate.md: The gate manifest (one job per row, the kind as the job's name); docs/spec/gate.md: REQ-004 (the gates add rules and weaken none); docs/spec/gate.md: REQ-007 — Each change passes the gates before review or merge (item 1); docs/spec/setup.md: The steps (S12); docs/spec/setup.md: The checks of `layup setup verify` (jobs); docs/spec/records.md: NFR-002 — A target is independent of LAYUP (item 2); docs/architecture.md: 6. Native gates and rule protection / The stack gates (In the target); docs/architecture.md: 15. Known limits (L-B3); docs/adr/0016-put-the-native-stack-gates-in-the-target.md: Decision 2 and Decision 5
- **Packages:** internal/catalog
- **Present state:** none. .github/workflows/ci.yml:295-318 holds LAYUP's own lint and tests jobs, not a target's.
- **Tests:**
  - unit: A Go test of internal/catalog checks these points: - The job names equal the kinds of kinds.tsv, one to one. - No job has a `paths:` filter or a job-level `if:`. - No job is a setup-check job, and no job starts `layup` or fetches a LAYUP file. - Each `uses:` names a version that has evidence (Invariant 4, README.md:31-34). - The form is plain enough for the stdlib line reader of the check `jobs`.
  - integration: The run script of each job (or the shared script it calls) runs with sh on fixture repositories for each line of the result table, for a missing manifest, and for a manifest with no row for the kind. It gives the same pass or fail as `layup gate` on the same base and head.
  - uat: On a pilot target on GitHub, a pull request that adds a .go file fails the layout, boundary and contract jobs with the reason, and a docs-only pull request passes them. The static and test jobs run their commands. There is no e2e test for this path, because phase 1 makes no forge call (docs/spec/README.md:125-126). The integration test is the nearest automated level.
- **External inputs:**
  - The versions of actions/checkout and actions/setup-go, with their documentation (Invariant 4). LAYUP uses @v4 and @v5 today (.github/workflows/ci.yml:298-301).
  - If the boundary job installs golangci-lint from the setup on (the activation changes no CI file, architecture.md:458-459): its pinned version and documentation.
  - A pilot repository on GitHub for the UAT.
- **Open questions:**
  1. The spec gives no workflow file name, no trigger (pull_request only, or also push to main), and no way to get the base SHA (checkout fetch-depth). The pending rule needs `git diff` against the pull request base.
  2. architecture.md:458-459: "The jobs of all kinds exist from the setup, so the activation changes no CI file." So the boundary job must install its tool at setup, or its command must fetch the tool (for example `go run …@version`). The sources do not choose.
  3. The scope rule and the result table exist twice: in Go (internal/gate) and in the job's shell. Does the job logic live in one script file in the target? That file is a target file and a rule path (architecture.md:514-518).

### `gate-verify-kind` — The check gate:<kind> of `layup setup verify`

- **What:** internal/verify runs the gate-run entry point twice for each active kind. The first run has --base and --head at the setup head and must give pass or clear. The second run has --head at a commit of the kind's fixture, applied on the setup head with `git apply`. That commit is an object on no ref, and the run must give fail. The first rule that matches decides: - The clean run gives fail: fail. - Either run gives not-active: not-active. - The fixture run does not give fail: fail, with the reason `fixture not detected`. - Otherwise: pass. A pending kind gives clear with `pending: fixture not run`. If the verify area owns the same item, merge them.
- **Specification:** docs/spec/setup.md: The checks of `layup setup verify` (gate:<kind>); docs/spec/packages.md: The table of phase 1 (internal/verify) and the split of internal/setup, internal/verify and internal/gate; docs/spec/gate.md: Not in phase 1 (the known-bad fixtures at setup); docs/adr/0016-put-the-native-stack-gates-in-the-target.md: Decision 4
- **Packages:** internal/verify, internal/git, internal/gate, internal/catalog
- **Present state:** none. No internal/verify exists.
- **Tests:**
  - unit: With a stand-in for the gate entry point, a table test covers each first-match line and its reason: clean fail gives fail; not-active in either run gives not-active; a fixture run that is not fail gives `fixture not detected`; a pending kind gives `pending: fixture not run` with zero gate calls.
  - integration: On a real temporary target from the Go entry, the output of `git for-each-ref` and the index are the same before and after the check. The fixture commit exists as an object and is on no ref. A fixture that does not apply gives fail.
- **Open questions:**
  1. docs/spec/setup.md:256 does not give the result when the fixture does not apply (`git apply` fails). Is it fail with its own reason, or `fixture not detected`?

### `later-gate-status` — Phase 2: the status `layup/gates` and the ready/merge rule

- **What:** `layup run` posts the commit status `layup/gates` from the gate-result table: a success only when each kind is pass or clear. In phase 2 the status becomes a required check, and no pull request is marked ready or merged while a selected kind did not run. These are not in phase 1.
- **Specification:** docs/spec/gate.md: Not in phase 1 (the commit status `layup/gates`); docs/spec/gate.md: REQ-007 (Not in phase 1); docs/spec/setup.md: Not in phase 1 (the `layup/` required checks); docs/architecture.md: 6. Native gates and rule protection / The stack gates (the table of layup/gates)
- **Packages:** internal/run, internal/forge, internal/forge/github
- **Present state:** none. Phase 1 makes no forge call (docs/spec/README.md:125-126).
- **Tests:**
  - integration: With a stand-in forge, each row combination maps to the status by the table of architecture §6. not-active is never a success.
  - e2e: A pull request whose head has a failed kind is not marked ready and does not merge.
- **External inputs:**
  - The App credentials and a pilot repository on the forge.

### `later-gate-rule-batch` — Phase 2: activation and the rule batch's own gate files

- **What:** For an approved rule batch, `layup gate` runs the batch's own gate files on its head, which must pass. It also runs the head with each recorded known-bad patch of each kind the batch touches, and each must fail. It binds this to the approval hash. The activation of the pending kinds at the first bet is such a batch. These are not in phase 1.
- **Specification:** docs/spec/gate.md: Not in phase 1 (a rule batch's own gate files); docs/spec/setup.md: Not in phase 1 (the activation of the pending kinds); docs/architecture.md: 6. Native gates and rule protection / The stack gates (Activation); docs/adr/0016-put-the-native-stack-gates-in-the-target.md: Decision 4; docs/walkthroughs/W-04-stack-dependent-gates.md: rows 2 to 5
- **Packages:** internal/gate, internal/rules, internal/records
- **Present state:** none.
- **Tests:**
  - integration: A batch head passes its own gate files. Each recorded patch fails its kind. A patch that passes, or that no longer applies and is not replaced, refuses the merge.

### Constraints of area `gate`

1. **docs/spec/packages.md:37** — internal/cli "May import" only `internal/psb`, `internal/setup`, `internal/verify`, `internal/gate`. *Effect:* Revision resolution, the manifest read and the gate-result writer live in internal/gate (gate-run). internal/cli only parses arguments and maps exit codes.
2. **docs/spec/packages.md:41-44** — internal/catalog may import only internal/tsv. internal/gate may import internal/tsv and internal/git. internal/setup and internal/verify import internal/catalog. *Effect:* internal/gate cannot read the catalog. Tests that join the gate and the catalog live in cmd/layup (e2e) or internal/verify. The catalog gives internal/setup and internal/verify the lookups they need.
3. **docs/spec/packages.md:16-21** — go.mod has no `require` line. Only internal/git starts `git`. Only the packages that the table marks "starts a program" import os/exec. *Effect:* There is no YAML library, so the target workflow keeps a simple line form. internal/gate starts only `sh -c`.
4. **docs/spec/packages.md:22-23** — "No package of phase 1 imports `net`, `net/http` or `crypto/tls`" *Effect:* gate-nfr005-imports is an integration test in the existing `tests` job. It cannot use `-test` dependencies for the os/exec rule (measured).
5. **docs/spec/gate.md:61-62** — "So the head's own gate files never judge the head (FT4)." *Effect:* gate-scratch-tree has an integration test where the head edits a gate file and the base version runs.
6. **docs/spec/gate.md:47-48** — "`layup gate` changes nothing in it except a scratch work tree, which it removes before it exits." *Effect:* Each gate test checks that REPO's refs, index and worktree list do not change, also on error paths. This also supports running git with hooks off.
7. **docs/spec/gate.md:52-54; docs/spec/README.md:55-58** — Exit 0 when each row is pass or clear, 1 for fail or not-active, and 2 for a usage or input error. "A check that did not run is never 0 (`NFR-004`)." *Effect:* gate-command has an e2e case for each code. Each failure path of git or of the scratch tree maps to a code other than 0.
8. **docs/spec/gate.md:63-64** — "For each row, in the order of the manifest ... the first line of the table that matches the row decides." *Effect:* gate-results is one ordered rule list that runs rows one at a time. Its unit tests cover the order cases.
9. **docs/spec/gate.md:76-78** — "A `pending` kind's command never runs" *Effect:* A unit test proves zero runner calls for pending rows.
10. **docs/spec/gate.md:94-102, 125-127, 164-165; docs/spec/setup.md:312** — The status `layup/gates`, a rule batch's own gate files, the required status, the ready/merge rule, the smart-if client and the activation of pending kinds are not in phase 1. *Effect:* The traceability matrix sends these parts of REQ-004, REQ-007 and NFR-005 to later-gate-status and later-gate-rule-batch (and to the smart-if item of phase 3).
11. **docs/spec/gate.md:119-123; docs/spec/README.md:59-61** — Two runs with the same REPO content, base and head print the same bytes. The table holds no time and no scratch path. *Effect:* gate-command's e2e test compares stdout and checks that it holds no temporary path. The fixture commits use fixed dates and a fixed identity.
12. **docs/spec/gate.md:144-146; docs/spec/setup.md:301-302** — The fixture test "comes with the code (#29)"; "LAYUP's own CI runs each fixture of each entry (§6); that test comes with the code (#29)." *Effect:* The NFR-004 fixture test is in gate-command. The catalog fixture test is in gate-catalog-go-entry, not in later test-only items.
13. **docs/spec/setup.md:276-280** — The catalog is `internal/catalog/<stack>/`, embedded in the binary with `embed`. *Effect:* The embed rules apply: no go.mod in the embedded tree, and an `all:` prefix for `.github/` (measured).
14. **docs/spec/setup.md:97** — S12: "The baseline's own workflows stay byte for byte (`REQ-018`). Adds no `setup-check` job." *Effect:* The catalog's workflow file name must not be a baseline file name, and S12 must refuse to overwrite a baseline path (setup area).
15. **docs/spec/records.md:127-130** — "The setup writes no file that the target needs LAYUP to build, test or pass its gates: no `layup` binary, no LAYUP script, no LAYUP CI job" *Effect:* gate-catalog-go-workflow serves NFR-002 and has a test that no job starts `layup` or fetches a LAYUP file.
16. **docs/architecture.md:453-459** — A pending kind's job fails a pull request that changes a product path and passes one that changes none, with the reason. An active kind never passes silently on nothing. "The jobs of all kinds exist from the setup, so the activation changes no CI file." *Effect:* Each job decides in its own step, with no skip filter. The pending-kind jobs are complete at setup.
17. **docs/architecture.md:461-463** — "`test -z "$(gofmt -l .)"` (bare `gofmt -l` exits 0 on a bad file)" *Effect:* The static command uses the test -z form. The static fixture is a badly formatted .go file.
18. **docs/adr/0016-put-the-native-stack-gates-in-the-target.md:41-42** — "A new entry is a LAYUP change under LAYUP's gate, never the work of a session on a target." *Effect:* The catalog items are LAYUP tasks with LAYUP's full gate.
19. **docs/spec/README.md:31-34** — A value that needs a decision of the Operator or the idea owner stays a marker with a row in open-gaps.tsv (Invariant 4). *Effect:* Each tool version needs its documentation as evidence before the Go entry merges. The coverage floor stays a gap.
20. **go.mod:3; .github/workflows/ci.yml:310-312** — `go 1.26`; setup-go with `go-version-file: go.mod`. *Effect:* The catalog's go.mod template must not need a Go newer than LAYUP's CI Go, or the fixture test downloads a toolchain.
21. **docs/architecture.md:267-268** — "the App has no workflows permission (O-92), so LAYUP never delivers a workflow change" *Effect:* Put the new Go tests in the existing `tests` job tags, so that ci.yml does not change.
22. **docs/tests/test-levels.md:43, 48, 60, 74, 77-81** — A unit test touches no file, network or clock. Every component has a unit test, every interface or workflow an integration test, and every user-facing scenario an e2e test. A UAT rides on the E2E path. *Effect:* Unit tests use stand-ins. Tests that start git, sh or go are integration or e2e tests. A UAT with no e2e test (the forge path) names its exception.
23. **docs/guardrails.md:272-280** — Testing after the code, and tests that depend on the wall clock, are pitfalls. *Effect:* No test-only item lands after its code. The golden tables do not depend on commit times.
24. **AGENTS.md (Bootstrap mode); docs/adr/0012-build-layup-in-bootstrap-mode.md** — A task is a PSB In-Scope item or what one needs, with one plan-review comment and one review round. *Effect:* Each gate slice is small enough for a one-round review.

### Risks of area `gate`

- The embed trap: the spec form `files/go.mod` and `files/.github/workflows/*` does not embed as written (measured). If the plan does not fix a rename rule, the first catalog task fails at build.
- The result rules exist twice, in Go (internal/gate) and in shell (the target's jobs). Without the parity test, a pull request can pass the native job and fail `layup gate`, or the reverse.
- Rename detection and host git configuration (diff.renames, core.quotePath) can hide a product path from the pending rule. The pending kind then gives a false clear (FT1).
- Determinism depends on the target's tools. Module fetches, GOTOOLCHAIN downloads and the inherited environment can change a verdict between runs.
- The Go scope `./*.go` is not `./...`. A tree with .go files and no package makes an active kind fail, not clear (measured). Setup verify or a pilot can then fail for a reason that the spec does not name.
- PRD-0001:117 asks that a seeded violation of each of the four kinds (layout, boundary, contract, test quality) reports fail. In phase 1, three kinds are pending with no fixture (setup.md:296), and their fail comes only from the pending rule. The plan must say if this meets the criterion in phase 1 or moves part of it to phase 2.
- If `git worktree add` runs REPO's hooks, a rule path (`.githooks/`) can run inside `layup gate` and change REPO (FT4). The spec does not turn hooks off for `layup gate`.
- A killed `layup gate` leaves a registered worktree in REPO. The spec has no recovery rule.
- No timeout exists for gate commands. A `go test` that hangs blocks `layup gate` and `layup setup verify` with no limit and no sign of life.
- The exit code of a failure of git or of the scratch tree is not given. A wrong map could give 0 for a check that did not run (NFR-004).
- Issue #33 (T-vk3k) still says that the gates run outside the target per ADR-0011 decision 4 (.ctx/issues/33.md). ADR-0016 superseded this. #33 must be re-scoped to the plan's issues or closed.
- A `&lsaquo;…&rsaquo;` coverage-floor marker or a Go template file inside internal/catalog/ can trip LAYUP's own check `markers` (setup-check.sh:305), gofmt, go vet or go list -deps.
- No setup step writes the gap of the coverage floor, so the check `sources` or the record can be wrong for the floor.
- External inputs (tool versions with documentation, the golangci-lint choice, the coverage floor, a pilot repository) can block the Go entry and the workflow items.

### Critic points that the reviser of `gate` rejected

- gate-catalog-package depends_on: add gate-manifest, because the integration test parses the derived manifest with the gate-manifest reader. *Basis:* docs/spec/packages.md:41-42: internal/catalog may import only internal/tsv. So I did not add the dependency. I moved that parse to the e2e test of gate-catalog-go-entry, which already depends on gate-manifest through gate-command. This is the critic's own alternative fix.
- gate-catalog-go-workflow tests: remove the assertion that each `uses:` is pinned to a version, as it has no stated basis. *Basis:* docs/spec/README.md:31-34 (Invariant 4) and docs/spec/setup.md NFR-003 item 2: each value that the setup writes needs a source. An action version is such a value. I kept the assertion and replaced its basis (it is no longer the ruleset pinning).
- gate-command size: 'small' is doubtful (more than 400 lines). *Basis:* docs/spec/packages.md:37 forced the split. Revision resolution, the manifest read and the table writer moved to gate-run. gate-command now holds only argument parsing, the exit-code map and the e2e tests, so small holds.
- gate-results size: 'small' is doubtful (about 450 lines). *Basis:* The table writer and the revision resolution moved to gate-run, and the scope matcher is in gate-manifest (docs/spec/gate.md:35-39, 80-89). The rule table, the sh runner and the stderr blocks remain, so I keep small.
- The gate:<kind> row has no owner: add it here with the full first-match order. *Basis:* Partly rejected. docs/spec/packages.md:44 and :51-53 put the check in internal/verify, which belongs to the verify area. I added gate-verify-kind with the exported interface and a merge note, so that it is not lost. It can move to the verify area.

## Area `setup`

Phase 1 makes `layup setup WORK` a step runner. It does S01 to S15 in one work area per target, with `git` as its only external program. It makes no forge call. It stops with exit 3 and one stop table when an input is missing. It resumes from WORK/out/record.tsv. It writes the setup commits on `layup-setup`, the orphan records commit on `layup-records`, and WORK/out/commands.sh for the Operator. No code for it exists today: internal/cli/cli.go:31-40 knows only `version` and `psb`, and internal/setup, internal/catalog, internal/tsv and internal/git do not exist. The final inventory has 18 items: the runner core and the answers rule; the step groups S01, S02+S03, commands.sh, S04, S05, S06, the prose copies S07/S08/S09/S14, S10, S11, the catalog package, S12, S13, the rule-path register, S15; one offline e2e test; one pilot UAT. Each step's `done` row needs a check of `layup setup verify`, which internal/cli runs after internal/setup does the step (packages.md:49-56), so most items depend on a verify item. The e2e test runs a whole setup with no network: it generates a fixture baseline repository in t.TempDir (git init, one commit), gives `S01-baseline` as `file://<dir>` so `git ls-remote` and `git clone` read it locally, points `origin` at a local bare repository, and runs only the `git push` lines of commands.sh. The ruleset apply line needs `gh` and the network, so the test checks only its text. The fixture is generated, not checked in, because LAYUP's checks `markers`, `adapted` and link-lint read checked-in testdata outside their exemptions. After the critic round, two measured blockers on the real baseline are on record. I measured a copy of LAYUP's root commit d2516fd (the baseline at LAYUP's pin) with S05's deletions and the five S07-S09/S14 files removed. Check `adapted` gives 235 failures in 21 files that no step writes. Check `kit-history` still fails on kit-link lines that S05's rule leaves: 7 in backlog.md and 56 in completed.md. On the real baseline, S05 and S15 cannot reach `done`, so the pilot cannot pass until the spec gives a rule. The sources leave about 40 points open, listed per item.

| Key | Title | Requirements | Size | After | Host |
| --- | ----- | ------------ | ---- | ----- | ---- |
| `setup-runner` | The `layup setup WORK` command and the resumable step runner | REQ-002, NFR-003, NFR-001 | large | found-tsv (foundation area: internal/tsv, the record reader and writer), found-git (foundation area: internal/git, with the verbs listed in setup-s15-records-commit and a fixed git configuration), found-cli-exit-codes (foundation area: command dispatch and the 0/1/2/3 mapping), found-schema-block-test (foundation area: the Go test that reads the tsv-schema blocks) | row 8 |
| `setup-answers` | answers.tsv: its schema and the rule for each row state | REQ-002, NFR-003 | small | setup-runner, found-tsv (foundation area) | row 8 |
| `setup-s01-questions` | S01: the Operator's answers and the gap questions in one stop table | REQ-002, REQ-001, NFR-003 | small | setup-runner, setup-answers, setup-catalog, psb-on-tsv (psb area: internal/psb moves to internal/tsv, packages.md:46-47) | row 9 |
| `setup-s02-s03-baseline` | S02 and S03: resolve the baseline once, copy it, make the root commit | REQ-002, NFR-006, NFR-003 | small | setup-runner, setup-s01-questions, setup-commands-file, found-git (foundation area) | row 9 |
| `setup-commands-file` | WORK/out/commands.sh in its fixed order | REQ-002, NFR-001 | small | setup-runner | row 8 |
| `setup-s04-pin` | S04: the target's pin file and the pin's decision record | NFR-006, REQ-002, NFR-003 | small | setup-s02-s03-baseline, verify-check-pin (verify area: check `pin` in the target's form) | row 9 |
| `setup-s05-history` | S05: delete the baseline's own history, with link fixes as input files | REQ-002 | large | setup-s04-pin, verify-check-kit-history (verify area), verify-check-link-lint (verify area: runs the baseline's link-lint.sh) | row 13 |
| `setup-s06-facts` | S06: the briefs and the first answers record as raw facts | REQ-002, NFR-001, NFR-003 | small | setup-s05-history, setup-s01-questions, verify-check-facts (verify area: check `facts` in a target's form) | row 13 |
| `setup-prose-inputs` | S07, S08, S09, S14: copy the prose input files and run their checks | REQ-002 | small | setup-s06-facts, verify-check-prose (verify area: onboarding, glossary, guardrails in a target's form, and identity) | row 13 |
| `setup-s10-markers` | S10: list every marker outside MK_EXEMPT, and stop for the missing answers | REQ-002, NFR-003 | small | setup-prose-inputs, setup-answers, found-git (foundation area: a verb that lists the tracked files) | row 13 |
| `setup-s11-fill` | S11: fill the markers, keep the gaps, write the second answers record | REQ-002, NFR-003, NFR-001 | small | setup-s10-markers, setup-s06-facts, verify-check-markers-sources (verify area: checks `markers` and `sources`) | row 13 |
| `setup-catalog` | internal/catalog: the embedded stack catalog and its reader | REQ-002, NFR-007 | small | found-tsv (foundation area) | row 4 |
| `setup-s12-stack-files` | S12: the stack's files, go.mod, the gate manifest and one CI job per kind | REQ-002, NFR-002, NFR-003 | small | setup-catalog, gate-catalog-go (gate area: the content of the Go entry, its kinds, tool versions, evidence URLs and known-bad fixtures; #29), setup-s11-fill, gate-manifest-reader (gate area: the gate-manifest schema), verify-check-jobs-gates (verify area: checks `jobs` and `gate:<kind>`) | row 15 |
| `setup-s13-rulesets` | S13: branch-protection.json, the default-branch ruleset and the push commands | REQ-002, NFR-002 | small | setup-s12-stack-files, setup-commands-file | row 15 |
| `setup-rule-paths` | The rule-path register from the catalog and the baseline | REQ-002, NFR-001 | small | setup-s12-stack-files | row 15 |
| `setup-s15-records-commit` | S15: the verify table, the last record rows and the orphan records commit | NFR-001, REQ-002, NFR-002 | small | setup-rule-paths, setup-s13-rulesets, setup-commands-file, verify-command (verify area: `layup setup verify WORK` and its table), found-git (foundation area: verbs for an orphan commit, add, branch and ls-files) | row 15 |
| `setup-e2e-offline` | End-to-end: a full setup with no network | REQ-002, NFR-001, NFR-003, NFR-006 | large | setup-s15-records-commit, setup-prose-inputs, setup-s05-history, gate-catalog-go (gate area: the only embedded entry is Go, so the e2e uses it), verify-command (verify area) | row 16 |
| `setup-pilot-acceptance` | REQ-002 acceptance on a pilot problem statement | REQ-002, NFR-003 | large | setup-e2e-offline | row 20 |

### `setup-runner` — The `layup setup WORK` command and the resumable step runner

- **What:** `layup setup WORK` checks the work area, reads WORK/out/record.tsv, skips each step that has a `done` row, and runs the next steps in order S01 to S15. With no stop it prints the step table, one row per step, in step order. At a stop it prints only the stop table, in a fixed row order. It exits 0, 1, 2 or 3. Each step that changes the tree is one commit `chore: setup <step>` on `layup-setup`, made with LAYUP's own git configuration and not the host's.
- **Specification:** docs/spec/setup.md: The command `layup setup`; docs/spec/setup.md: The steps (the commit-per-step rule, line 102); docs/spec/setup.md: The stop table; docs/spec/setup.md: The step table; docs/spec/setup.md: The setup record; docs/spec/setup.md: The boundary of phase 1; docs/spec/README.md: Commands (Arguments, Exit codes, Determinism); docs/spec/README.md: Records; docs/spec/packages.md: The components of phase 1; docs/architecture.md: 5. Intake and setup (Scaffold 1)
- **Packages:** internal/setup, internal/cli, internal/tsv, internal/git
- **Present state:** none. internal/cli/cli.go:31-40 dispatches only `version` and `psb`; internal/setup does not exist.
- **Tests:**
  - unit: With a stubbed git and a stubbed step checker, a record.tsv with `done` rows for S01..S04 makes the runner start at S05 and call no earlier step.
  - unit: Exit code mapping: a stop gives 3 and only the stop table on stdout; a failed check gives 1 with a `fail` row; a not-runnable check gives 1 with `not-active`; a schema error in record.tsv gives 2; all steps done or `operator` gives 0.
  - unit: Determinism: two runs of the runner on the same stubbed inputs give the same bytes for the step table and for the stop table, and the rows hold no time and no temp path.
  - integration: The setup-stop, setup-steps and setup-record outputs match their tsv-schema blocks in docs/spec/setup.md (the schema-block test reads a file, so it is integration, docs/tests/test-levels.md:43).
  - integration: With a real git in a temp work area, a step that changes the tree makes exactly one commit `chore: setup <step>` on `layup-setup`; a second run makes no new commit for a done step; a host with a global hook, commit.gpgsign=true and core.autocrlf=true does not change the commit or its tree.
  - e2e: The built binary run as `layup setup` with no WORK, with a flag, or with a WORK that does not exist exits 2 and writes a diagnostic only on stderr.
- **External inputs:**
  - An Operator decision on the git identity (author, committer, and whether dates are fixed) of the setup commits and the records commit.
- **Open questions:**
  1. The step table has no result for a step that did not run. docs/spec/setup.md:128 gives `enum(done|fail|not-active|operator)`, and :123 says "one row per step". When S05 fails, the sources give no value for the rows of S06 to S15.
  2. The actor of a step done by layup setup with a push by the Operator is not clear. docs/spec/setup.md:127 gives `enum(layup-setup|operator|layup-run)`, and :88 gives S03 as "`layup setup`; the push by the Operator". It is not clear whether S03, S13 and S15 are `layup-setup` with result `operator`, or actor `operator`.
  3. S13 cannot get a `done` row. Its evidence is "the Operator's run of `commands.sh`" (docs/spec/setup.md:98), which phase 1 cannot see. The sources do not say which row S13 writes to record.tsv, whether the runner goes on to S14 and S15 without a `done` row for S13, or whether each rerun does S13 again.
  4. How the runner finds that an input changed is not clear. docs/spec/setup.md:55 makes exit 2 "an input that changed after a step read it", but only S06 (setup.md:91, `brief.sha256`) has a mechanism. There is no hash row for answers.tsv rows already used or for inputs/files/<path>.
  5. The rule for absent inputs is not given. docs/spec/README.md:55 makes "a missing or unreadable file" exit 2, and docs/spec/setup.md:56 makes a missing input exit 3. The sources do not say whether an absent WORK/inputs/answers.tsv on the first run is exit 2 or a stop with all S01 questions, or what an absent inputs/briefs/problem-statement.md gives (docs/guardrails.md:294-303 asks for one rule per input state).
  6. Commit identity and dates are not given. No source names the author and committer of the root commit, the setup commits and the records commit, or whether their dates are fixed. A host with no git user.email cannot commit, and the commit SHAs (the evidence) change per run, against docs/spec/README.md:61-63.
  7. The host's git configuration is not ruled out. docs/setup/tests/run.sh:44 runs git with `-c core.hooksPath=/dev/null -c commit.gpgsign=false` and a fixed identity. No spec says that layup setup does the same, or sets core.autocrlf=false so the root tree equals pin.tree.
  8. The row order of the stop table is not given. docs/spec/README.md:63 asks for "a fixed order". The sources do not give the order of the mixed `S01-` and `Q-` rows of S01, the `F-` rows of S05, or the `M-` rows of S10.

### `setup-answers` — answers.tsv: its schema and the rule for each row state

- **What:** Reads WORK/inputs/answers.tsv against the setup-answers schema. A row that does not match the schema exits 2. A row for a question that no step of this run asked exits 2 and names the row: a stale `M-` row, an `F-` row, an `O-` row, an unknown `S01-` ID, a `Q-` ID that the gap table does not have. S01 checks the `S01-` and `Q-` rows; S10 checks the `M-` rows and the rest. The answer `gap` must carry `question_text`.
- **Specification:** docs/spec/setup.md: The answers; docs/spec/setup.md: The stop table (question ID forms; the M-<x8> rule); docs/guardrails.md: 2 (A specification that fixes one input state at a time)
- **Packages:** internal/setup, internal/tsv
- **Present state:** none
- **Tests:**
  - unit: Each row state gives its result, on in-memory rows: an asked row is accepted; an `F-`, `O-`, unknown `S01-`, unknown `Q-` or unlisted `M-` row is exit 2 with the row named; a `gap` row with `—` as question_text is exit 2; a bad `by` value is exit 2; two rows with one question ID give the result that the open question below decides.
  - unit: The marker ID is `M-` plus the first 8 hex characters of SHA-256 of `<file>\t<marker>`, and it does not change when the marker moves to another line.
  - unit: An `S01-` or `Q-` row added after S01 has its `done` row is refused at S10 ("S10 checks ... the rest"), so S06 cannot write an unchecked row as a fact.
  - integration: A real answers.tsv file is read through internal/tsv and the rule, and the reader matches the setup-answers tsv-schema block of docs/spec/setup.md.
- **Open questions:**
  1. Can a question ID appear on two rows? docs/spec/setup.md:143 makes `question` the key, but the rule text does not say if a second row for the same ID is exit 2 or a replacement.
  2. The answer `gap` to a non-marker question is not given. docs/spec/setup.md:144 says "`gap` keeps a marker as an open gap". The sources do not say whether `gap` for an `S01-` or `Q-` question is refused (exit 2) or accepted.
  3. A changed `S01-` or `Q-` row after S01 is done is not covered. S01 does not run again on resume (docs/spec/setup.md:57-58), and only brief.sha256 detects a change (:91). The sources do not say whether S10's "the rest" (:139-140) re-checks these rows or how a changed value is found.

### `setup-s01-questions` — S01: the Operator's answers and the gap questions in one stop table

- **What:** Reads the answers `S01-stack`, `S01-name`, `S01-visibility` and `S01-baseline`. internal/cli runs the `layup psb check` rules on inputs/briefs/problem-statement.md and gives the gap table (`Q-NNN` with its question text) to internal/setup. The run stops once and lists every missing answer of both kinds. When all are present and the stack has a catalog entry, it writes the record rows `stack`, `name`, `visibility`, `baseline` (source `answer`, or `fact` when the answer's source is an accepted `F-NNNN#n`) and `brief.sha256` (source `computed`).
- **Specification:** docs/spec/setup.md: The steps (S01); docs/spec/setup.md: The setup record (the source kinds); docs/spec/psb-check.md: The table; docs/architecture.md: 5. Intake and setup (The gap check, item 4)
- **Packages:** internal/setup, internal/cli, internal/psb, internal/catalog
- **Present state:** internal/psb/check.go exists and internal/cli/cli.go:43-59 runs it for `layup psb check`; nothing hands its table to a setup step.
- **Tests:**
  - unit: With no answers and two gaps, the stop table lists the four S01- questions and both Q- IDs in one table, in a fixed order, with exit 3.
  - unit: An `S01-stack` answer that has no catalog entry fails the step's evidence check.
  - unit: An answer whose source is `F-NNNN#n` gives a record row with source `fact` and ref `F-NNNN#n`; an answer whose source is a comment URL gives source `answer` and ref the question ID (the reading of the open question below).
  - integration: internal/cli hands the real psb gap table to internal/setup, and the `Q-NNN` IDs in the stop table equal the IDs of `layup psb check` on the same file.
- **Open questions:**
  1. The form of `S01-name` is not given. docs/spec/setup.md:97 derives the go.mod module path from `name`, and S13 needs OWNER/REPO for the ruleset command. The sources do not say whether `name` is `REPO`, `OWNER/REPO` or a full module path.
  2. The use of `visibility` in phase 1 is not clear. docs/spec/setup.md:86 records it as an answer, but no phase-1 step reads it.
  3. When a record row gets source `fact` is not given. docs/spec/setup.md:146 lets an answer's source be "a fact citation `F-NNNN#n` that the Operator accepted", and :159-160 has the record source `fact` with ref `F-NNNN#n`. The sources do not say whether such an answer gives a `fact` row or an `answer` row.

### `setup-s02-s03-baseline` — S02 and S03: resolve the baseline once, copy it, make the root commit

- **What:** S02 runs `git ls-remote <url> HEAD`, clones and checks out that commit, takes its tree with `git rev-parse <commit>^{tree}`, removes `.git`, and writes the record rows `pin.source`, `pin.commit`, `pin.tree` and `pin.time`. git never asks for a login on the terminal; a URL that needs one fails with a diagnostic. S03 runs `git init` on `main` and makes one commit of the unmodified copy (`chore: the unmodified baseline at <commit>`). It checks that the root tree equals `pin.tree` and adds the push command of the root commit to commands.sh. After the S02 done row exists, no run resolves the commit again.
- **Specification:** docs/spec/setup.md: The steps (S02, S03); docs/spec/setup.md: The setup record (source `computed`); docs/spec/setup.md: NFR-006 — The baseline at a pinned, recorded version; docs/spec/README.md: Commands (Arguments: "asks no question on the terminal"); docs/architecture.md: 5. Intake and setup (Start, item 2)
- **Packages:** internal/setup, internal/git
- **Present state:** none. LAYUP's own setup used `npx degit` (docs/setup/steps.tsv:3, docs/setup/armature.pin:8).
- **Tests:**
  - unit: With a stubbed git, S02 writes pin.commit and pin.tree with source `computed` and ref the command, and pin.source with source `answer` and ref `S01-baseline` (the reading of the open question below); S03 refuses a root tree that differs from pin.tree (exit 1).
  - integration: With a local fixture repository and a file:// URL, S02 and S03 give a WORK/target whose root commit tree equals the fixture commit's tree, and `.git` of the clone is gone.
  - integration: Resolve once: after S02 is done, a new commit in the fixture repository does not change pin.commit on the next run.
  - integration: A baseline URL that asks for credentials makes S02 fail fast with a diagnostic on stderr and no prompt (git run with no terminal input and GIT_TERMINAL_PROMPT=0 or an equal setting).
- **External inputs:**
  - The real baseline repository URL for a non-test run (LAYUP's own pin names https://github.com/pharzam/armature, docs/setup/armature.pin:5; a target's pin is the latest commit at setup time, O-101).
- **Open questions:**
  1. The source kind of pin.source is not clear. docs/spec/setup.md:87 takes it from the answer `S01-baseline`, but :159 says "the pin values come from `git`, not from a person" (`computed`). S01 already writes the row `baseline` from the same answer.
  2. The meaning of `pin.time` is not clear. docs/architecture.md:308 says "source, commit, tree and time go into the first records commit", and docs/spec/setup.md:348 has `date=<the date of pin.time, YYYY-MM-DD>`. The sources do not say whether it is the clock time of the resolve or the commit time of the baseline commit.
  3. A crash in the middle of S02 is not handled. It can leave a clone with `.git` removed but no `done` row. The next run would then resolve the commit again, against docs/spec/setup.md:57-60 ("the pin is resolved once"). The sources give no rule for a partial step.
  4. The printed push command and the remote are not clear. docs/architecture.md:313-315 asks the printed push command to show "the resolved commit and its difference from LAYUP's own pin", and docs/spec/setup.md:88 does not repeat this. No source says who adds the `origin` remote that `git push origin` needs.
  5. The credential prompt is not ruled out in words. docs/spec/README.md:44-45 says a command "asks no question on the terminal", but no source says how internal/git stops git from asking.

### `setup-commands-file` — WORK/out/commands.sh in its fixed order

- **What:** Writes commands.sh with a comment line before each command, in the order of the spec. The order is: push of the root commit (S03); push of `layup-setup` onto `main` (S13); push of `layup-records` (S15); apply of the default-branch ruleset (S13). A later step adds its command, and the file is written again in this order. It is not appended in step order.
- **Specification:** docs/spec/setup.md: The command `layup setup` (the work area); docs/spec/setup.md: Where the records go in phase 1
- **Packages:** internal/setup
- **Present state:** none
- **Tests:**
  - unit: When S03, S13 and S15 add their commands in step order, the text holds the S15 push before the ruleset apply, each with its comment line.
  - integration: `sh -n` accepts the written file.
- **Open questions:**
  1. The ruleset apply command is not given. docs/spec/setup.md:98 says only "the apply command of the ruleset". docs/architecture.md:185-186 asks that it print the bypass list. The sources give no command text, for example `gh api` with its method and path.

### `setup-s04-pin` — S04: the target's pin file and the pin's decision record

- **What:** Makes the branch `layup-setup` from the root commit. Writes docs/setup/armature.pin from the pin rows, with the keys and the order of LAYUP's own pin, each key once, and `method=git clone <source>, checkout <commit>, .git removed`. Writes the decision record of the pin in the baseline's ADR form from a fixed text with the pin values, at the next free ADR number, with its row in the index table of docs/adr/README.md. It installs no hooks.
- **Specification:** docs/spec/setup.md: The steps (S04); docs/spec/setup.md: NFR-006 — The baseline at a pinned, recorded version; docs/spec/setup.md: NFR-003 — No value without evidence (item 1)
- **Packages:** internal/setup, internal/git
- **Present state:** LAYUP's own pin docs/setup/armature.pin:5-9 (degit method); no writer.
- **Tests:**
  - unit: The pin text equals the NFR-006 form for given pin rows: five keys in order, each once, and `date` taken from pin.time.
  - unit: The ADR text has the next free number after the baseline's highest ADR (0009 at LAYUP's pin), a `Date:` line, an allowed `## Status`, the sections `## Context`, `## Decision`, `## Consequences`, and an index row.
  - integration: After S04, check `pin` of the verify area passes on WORK/target, and it fails when the pin's tree differs from the root tree.
  - integration: The baseline's own `sh docs/adr/adr-lint.sh`, run by the test on the generated fixture target after S04, exits 0 (no phase-1 verify check runs adr-lint, so this test is the only guard before the target's own CI).
- **External inputs:**
  - An Operator decision, or a spec fix, on the fixed text of the pin ADR.
- **Open questions:**
  1. The ADR text is not specified. docs/spec/setup.md:89 says "The decision record of the pin is the baseline's own ADR form, written from a fixed text with the pin values", but gives no file name or text. Issue #34 item 4 (.ctx/issues/34.md:12) says "S04 says 'ADR for the pin' as if each target needs one". The baseline's adr-lint needs contiguous numbering and an index row (baseline docs/adr/README.md:56-59 at d2516fd), so the number is the next free one (0009 at LAYUP's pin).
  2. The ADR and the pin file put values into the target (the ADR number, its `Date:`, its status). docs/spec/setup.md:320-322 asks a record row with a source for each value that layup setup writes. The sources do not say which rows S04 writes for these values.

### `setup-s05-history` — S05: delete the baseline's own history, with link fixes as input files

- **What:** Deletes docs/decisions/, docs/audit/, each docs/tasks/T-*.md, and their lines in backlog.md and completed.md. When the deletion breaks a link, the run stops with an `F-<path>` question for each such file. The fixed text then comes as WORK/inputs/files/<path>, and the step copies it in.
- **Specification:** docs/spec/setup.md: The steps (S05); docs/spec/setup.md: The stop table
- **Packages:** internal/setup, internal/git
- **Present state:** none. LAYUP did it by hand (docs/setup/steps.tsv:6).
- **Tests:**
  - unit: On an in-memory tree, the step removes the three history groups and only the backlog and completed lines that name a deleted T-ID.
  - unit: A link into a deleted file gives a stop row `F-<path>` with exit 3; with the input file present, the step uses it and goes on.
  - integration: After S05 on a fixture target in a temp directory, checks `kit-history` and `link-lint` pass, and fail when one T-*.md is left.
  - integration: A fixture backlog line that links github.com/pharzam/armature/ and names no T-*.md gives the result that the open question below decides (today: check `kit-history` fails, so S05 has no `done` row).
- **External inputs:**
  - An Operator decision, or a spec fix, on the lines and notes that S05 deletes beyond the lines of a deleted T-ID.
- **Open questions:**
  1. How S05 finds a broken link is not clear. docs/spec/packages.md:43 says internal/setup starts no program, so it cannot run the baseline's link-lint.sh. docs/spec/setup.md:90 says "A link that the deletion breaks is a missing input". The sources do not say whether the engine has its own link parser (the baseline's link-lint.sh is about 350 lines of sh) or whether internal/cli runs verify's link-lint before the commit. The size depends on this answer.
  2. The rule "their lines" fails at LAYUP's pin. check_kit_history fails on any line of backlog.md or completed.md that links `github.com/pharzam/armature/` (docs/setup/setup-check.sh:106-109). I measured a copy of root commit d2516fd: of the lines that link it, 7 in backlog.md and 56 in completed.md name none of the 19 deleted T-IDs. docs/spec/setup.md:90 names only "their lines"; steps.tsv S05 adds "notes" (docs/setup/steps.tsv:6). So on the real baseline, S05 cannot get its `done` row. The sources give no rule for these lines.

### `setup-s06-facts` — S06: the briefs and the first answers record as raw facts

- **What:** Refuses a problem statement whose SHA-256 differs from `brief.sha256` (exit 2). Copies each brief byte for byte into docs/facts/. Writes the `S01-` and `Q-` answers as one raw fact record, one fact per question ID with the ID, the question text and the answer. internal/cli gives the gap table to internal/setup again at S06 for the `Q-NNN` question texts. Writes docs/setup/facts.sha256 in sha256sum form and the index rows of docs/facts/README.md. S10 and S11 never edit docs/facts/.
- **Specification:** docs/spec/setup.md: The steps (S06); docs/spec/setup.md: The files of `docs/setup/` in a target; docs/spec/records.md: NFR-001 — Git is the system of record (item 3, the answers); docs/spec/records.md: The layout of the records branch (In other places: docs/facts/, facts.sha256); docs/spec/psb-check.md: The table (the `question` column); docs/architecture.md: 5. Intake and setup (Scaffold 3)
- **Packages:** internal/setup, internal/cli
- **Present state:** LAYUP's own raw facts exist (docs/facts/, docs/setup/facts.sha256); no writer.
- **Tests:**
  - unit: On in-memory inputs, a changed problem statement after S01 gives exit 2; an unchanged one gives byte-equal copies and a facts.sha256 line `<64 hex>  <path>` per raw file.
  - unit: The answers record holds exactly one fact per `S01-` and `Q-` question ID, with the question text of the gap table for each `Q-` ID, and no `M-` ID.
  - integration: internal/cli hands the gap table to S06, and after S06 check `facts` passes on the fixture target; it fails when one byte of a brief changes.
- **Open questions:**
  1. The file names in the target are not given. The inputs are inputs/briefs/problem-statement.md and vision.md (docs/spec/setup.md:67-68), but the facts.sha256 example uses `docs/facts/problem-statement-brief.md` (setup.md:198).
  2. The answers record is not fully given. The baseline template asks for "Collected by" and "Date collected" (baseline docs/facts/template.md), and the index has a "Collected" column. answers.tsv has no date column. NFR-003 forbids a value with no source. The F-NNNN number of the answers record is also not given.
  3. The handoff of the question texts is not given. Only the `question` column of the psb gap table (docs/spec/psb-check.md:53) holds the text of a `Q-NNN` question, and internal/setup cannot import internal/psb (docs/spec/packages.md:43). The sources do not say whether internal/cli hands the table again at S06 or S01 keeps the texts in a file of WORK.
  4. The answer's source and `by` are not in the record. docs/spec/records.md:105-106 says "each answer row names the comment that holds it (`source`)", and :108-110 says nothing in the work area is the only copy of a value. docs/spec/setup.md:91 gives each fact only "the question ID, the question text and the answer", and no step commits answers.tsv. The sources do not say whether the fact also carries `source` and `by`.
  5. A marker inside a copied brief is not handled. MK_EXEMPT does not exempt docs/facts/ (docs/setup/setup-check.sh:305), so a `&lsaquo;...&rsaquo;` in the problem statement becomes an S10 question, and a fill by S11 would change an immutable raw fact and break its facts.sha256 line.

### `setup-prose-inputs` — S07, S08, S09, S14: copy the prose input files and run their checks

- **What:** Copies inputs/files/docs/onboarding-for-engineers.md, docs/glossary.md, docs/guardrails.md, README.md and AGENTS.md into the tree, one step each (S14 copies two files). A missing file stops the run with an `F-<path>` question. Each step is done only when its check passes: `onboarding`, `glossary`, `guardrails`, `identity`.
- **Specification:** docs/spec/setup.md: The steps (S07, S08, S09, S14); docs/spec/setup.md: The boundary of phase 1 (Scaffold 2); docs/architecture.md: 5. Intake and setup (Scaffold 2); docs/architecture.md: 3. Records and identities (Where the records live: the README names the branch)
- **Packages:** internal/setup, internal/cli
- **Present state:** none
- **Tests:**
  - unit: On an in-memory tree, a missing input gives a stop row `F-docs/glossary.md` with the path in `where`, exit 3; a present input is placed at its path byte for byte.
  - integration: internal/cli runs the matching verify check after each copy; an onboarding file that holds a marker makes S07 `fail` with exit 1 and no done row.
- **External inputs:**
  - The prose texts themselves for a real target: written by a human or an agent in phase 1 (docs/spec/setup.md:69).
- **Open questions:**
  1. The order of S14 is not clear. It is after S13 in the table, but it only needs its input files. The stop table lists only "every missing input of that step" (docs/spec/setup.md:107-108), so S07, S08, S09 and S14 each stop in their own run when their file is missing.
  2. S14 runs after S10 and S11. The baseline's AGENTS.md:104 at d2516fd holds the marker `&lsaquo;worktree dir&rsaquo;`. S10 asks for it and S11 fills it or keeps it with an open-gaps.tsv row; then S14 overwrites the file. A `gap` answer leaves a stale open-gaps row, so check `markers` fails; a value answer leaves a record row for a value that is not in the tree. A marker inside an S14 input is never asked by S10, so check `markers` fails as unlisted. The sources do not say whether S10 skips the paths that S14 replaces, or whether S14 inputs must hold no marker.
  3. The README's content rule is not checked. docs/architecture.md:111 says "The target's README, written at setup, names the branch" (`layup-records`). S14 copies the input as it is and runs only check `identity`. The sources do not say whether S14 refuses a README that does not name the branch.

### `setup-s10-markers` — S10: list every marker outside MK_EXEMPT, and stop for the missing answers

- **What:** Lists the tracked files of the tree as `git ls-files` does, and scans them with the rule of check_markers: it skips a `&lsaquo;` in a code span and the `&lsaquo;…&rsaquo;` mention, and a marker that does not close on its line runs to the line end. It skips the paths that LAYUP's MK_EXEMPT matches, embedded at the engine's version. Each (file, marker) pair becomes one question `M-<x8>`, one stop-table row per ID. All markers with no answer row stop the run in one table. Then it checks the `M-` rows and refuses every other unasked row (exit 2).
- **Specification:** docs/spec/setup.md: The steps (S10); docs/spec/setup.md: The stop table (the M-<x8> rule); docs/spec/setup.md: The answers; docs/spec/setup.md: NFR-003 — No value without evidence; docs/architecture.md: 5. Intake and setup (The gap check, item 4)
- **Packages:** internal/setup, internal/git
- **Present state:** The rule exists only as shell: docs/setup/setup-check.sh:294-324 (MK_EXEMPT and the awk scan).
- **Tests:**
  - unit: On in-memory text, the scan skips the code-span mention and `&lsaquo;…&rsaquo;`, takes an unclosed marker to the line end, and gives one ID for a marker that occurs on four lines of one file.
  - unit: With three unanswered markers, one stop table lists all three in a fixed order with `where` = `<file>:<line> <marker>`; a stale `M-` row for a marker not in the tree gives exit 2.
  - integration: Parity: the scan gives the same (file, marker) set as check_markers on each case of docs/setup/tests/markers/, built as run.sh builds it (a temporary Git repository), including the quoted-name case of issue #21 (bad-unlisted/overlay/docs/é.md).
  - integration: The Go copy of MK_EXEMPT equals the value of MK_EXEMPT in docs/setup/setup-check.sh:305 (embed cannot read that file, docs/spec/setup.md:277).
  - integration: On a generated fixture baseline, S10 lists the expected marker IDs, and the IDs stay the same after a line is added above a marker.
- **Open questions:**
  1. A marker on many lines of one file is not fully given. Its ID is one (docs/spec/setup.md:117-119), and `question` is the key of the stop table (:112), so the table has one row per ID; but `where` names one line (:114). The baseline repeats markers: `&lsaquo;security test command&rsaquo;` 4 times in docs/ci/github-actions-ci.yml and `&lsaquo;worktree dir&rsaquo;` 4 times in docs/engineering-discipline.md (d2516fd). The sources do not say which line `where` names. The record has one row per line, `marker:<file>:<line>` (:96).
  2. The list verb is not in internal/git. check_markers lists files with `git ls-files` (docs/setup/setup-check.sh:307), but docs/spec/packages.md:39 has no `ls-files` verb.
  3. See setup-prose-inputs: S10 runs before S14 replaces README.md and AGENTS.md.

### `setup-s11-fill` — S11: fill the markers, keep the gaps, write the second answers record

- **What:** Replaces each marker that has a value with that value, and writes a record row `marker:<file>:<line>` with source `answer` (or `fact`). A marker whose answer is `gap` keeps its marker, gets a record row with source `gap` and ref `docs/setup/open-gaps.tsv`, and gets a row in docs/setup/open-gaps.tsv (no header) with the answer's question_text. When S10 listed at least one marker, it writes the `M-` answers as a second raw fact record with its own index row and facts.sha256 line. The S06 record does not change, and S11 edits nothing under docs/facts/ except the new record.
- **Specification:** docs/spec/setup.md: The steps (S11); docs/spec/setup.md: The setup record (sources `answer`, `fact`, `gap`); docs/spec/setup.md: The files of `docs/setup/` in a target; docs/spec/setup.md: NFR-003 — No value without evidence
- **Packages:** internal/setup
- **Present state:** LAYUP's own docs/setup/open-gaps.tsv has the no-header form; no writer.
- **Tests:**
  - unit: On in-memory text, a value answer replaces the marker text exactly; a `gap` answer leaves the marker, writes `<file>\t<marker>\t<question>` and a record row with source `gap` and ref `docs/setup/open-gaps.tsv`; with no marker, there is no second record and no new facts.sha256 line.
  - unit: A marker that does not close on its line is replaced by the rule that the open question below decides.
  - integration: After S11, checks `markers`, `sources` and `facts` pass on the fixture target; a removed open-gaps row makes `markers` fail; the second answers record holds no marker of its own.
- **Open questions:**
  1. The same provenance, source and date questions as S06 apply to the second answers record (its F-NNNN number, its "Date collected", the answer's `source` and `by`).
  2. The second record can fail check `markers`. MK_EXEMPT does not exempt docs/facts/ (docs/setup/setup-check.sh:305). The question of an `M-` ID is the marker, so the record will likely quote it with its angle quotes. The sources do not say how the record names a marker.
  3. The replacement of an unclosed marker is not given. check_markers says one that "does not close on its line runs to the line end" (docs/setup/setup-check.sh:295-296). The baseline's docs/tasks/backlog.md:20 at d2516fd holds `&lsaquo;State your exact scheme`, which closes on a later line. The sources do not say what S11 replaces.

### `setup-catalog` — internal/catalog: the embedded stack catalog and its reader

- **What:** Embeds internal/catalog/<stack>/ with `embed`. It reads kinds.tsv against the catalog-kinds schema. It lists files/<path> with `{{module}}` replaced, and gives each kind's fixture path. S01 asks it whether a stack has an entry. S12 and S15 read its files and config paths. Its form comes first; the Go entry's content (gate-catalog-go, #29) fills it.
- **Specification:** docs/spec/setup.md: The stack catalog; docs/spec/packages.md: The table of phase 1 (internal/catalog)
- **Packages:** internal/catalog, internal/tsv
- **Present state:** none. internal/catalog does not exist.
- **Tests:**
  - unit: An in-memory test entry's kinds.tsv is read against the catalog-kinds schema; a bad state value or a missing fixture for an `active` kind is an error.
  - unit: `{{module}}` in a catalog file is replaced with the given module path, and no other text changes.
  - integration: The embedded entry is read through internal/tsv, and the catalog-kinds tsv-schema block of docs/spec/setup.md equals the header that the reader expects (the schema-block test reads a file, so it is integration).

### `setup-s12-stack-files` — S12: the stack's files, go.mod, the gate manifest and one CI job per kind

- **What:** Writes the catalog entry's files into the tree. For Go, that is go.mod with the module path from `name` and the tools' configuration. It writes docs/gates.tsv (the catalog table without `version`, `fixture` and `evidence`) and one CI job per gate kind, named as the kind. It keeps the baseline's own workflows byte for byte and adds no setup-check job. Each value gets a record row with source `catalog` or `answer`.
- **Specification:** docs/spec/setup.md: The steps (S12); docs/spec/setup.md: The stack catalog; docs/spec/gate.md: the gate manifest (gate-manifest schema); docs/spec/records.md: NFR-002 — A target is independent of LAYUP; docs/architecture.md: 5. Intake and setup (Scaffold 4)
- **Packages:** internal/setup, internal/catalog
- **Present state:** none
- **Tests:**
  - unit: On an in-memory tree, the written docs/gates.tsv matches the gate-manifest schema and has one row per catalog kind; the files of .github/workflows/ that the baseline had are unchanged.
  - unit: go.mod holds the module path derived from `name` and no `require` line from LAYUP; no file names `layup` as a tool the target needs.
  - integration: After S12 on the fixture target, the baseline's workflows are byte-equal, check `jobs` finds a job per kind, and `gate:<kind>` for an active kind gives pass on the clean head and fail on its fixture.
- **External inputs:**
  - The Go entry's content from #29 / gate area (tool configuration files, the CI workflow of the gate jobs), with each tool's version and the URL of its documentation at that version (the `evidence` column, docs/spec/setup.md:297).
- **Open questions:**
  1. How the module path is derived is not given. docs/spec/setup.md:97 says "go.mod with the module path from `name`" (see setup-s01-questions). The rule from `name` to a module path (a forge host prefix or not) is not given.
  2. The CI file is not named. "The baseline's own workflows stay byte for byte" (setup.md:97) suggests a new file; the catalog row `files/<path>` (setup.md:285) says "the CI workflow of the gate jobs", but it gives no file name.
  3. The evidence check name differs. docs/spec/setup.md:97 names "checks `jobs` and `gates`", but the verify table (:256) has no check `gates`, only `gate:<kind>`, one per kind. This inventory reads `gates` as every `gate:<kind>` row.

### `setup-s13-rulesets` — S13: branch-protection.json, the default-branch ruleset and the push commands

- **What:** Writes docs/setup/branch-protection.json in the form of LAYUP's own file. It writes WORK/out/ruleset-default.json for the default branch and the ref `layup-probe`. The ruleset requires a pull request and each gate job as a required check pinned to GitHub Actions. It forbids force push and deletion, and its bypass list is empty. It requires no `layup/` check. It adds to commands.sh the push `git push origin layup-setup:main` and the ruleset apply. The step result is `operator`.
- **Specification:** docs/spec/setup.md: The steps (S13); docs/spec/setup.md: Not in phase 1; docs/spec/records.md: NFR-002 — A target is independent of LAYUP; docs/architecture.md: 5. Intake and setup (Scaffold 6); docs/architecture.md: 3. Records and identities (What the forge enforces)
- **Packages:** internal/setup
- **Present state:** LAYUP's own docs/setup/branch-protection.json:1-70 (classic protection body, app_id 15368); no writer; no ruleset JSON exists.
- **Tests:**
  - unit: The ruleset JSON includes the default branch and refs/heads/layup-probe, one required check per gate kind with integration_id 15368, non_fast_forward and deletion rules, an empty bypass_actors list, and no `layup/` context.
  - unit: branch-protection.json has the key set of LAYUP's docs/setup/branch-protection.json and one check per gate job.
  - integration: With a real git on the fixture target, S13 makes one commit `chore: setup S13` on `layup-setup` that holds branch-protection.json, writes ruleset-default.json outside the tree, rewrites commands.sh in its fixed order, and a rerun makes no second commit.
  - uat: The Operator applies ruleset-default.json to a test repository with the printed command, and GitHub accepts it with no edit (no e2e covers the apply, because it needs the network).
- **External inputs:**
  - The GitHub REST API documentation of repository rulesets (body form, rule types, integration_id) at a stated version.
  - The GitHub Actions app ID (LAYUP's file uses 15368, docs/setup/branch-protection.json:7) confirmed by its source.
  - A test repository on GitHub and the Operator's login for the UAT.
- **Open questions:**
  1. The set of required checks is not clear. docs/spec/setup.md:98 says "each gate job a required check". docs/setup/steps.tsv:14 says "one required check per job". docs/architecture.md:413 says the rulesets "hold the same required checks and LAYUP's". The sources do not say whether the baseline's own CI jobs (adr-lint, link-lint, ...) are required checks too.
  2. The ruleset name and the pull_request rule parameters (required approvals, dismiss stale reviews) are not given in any source.
  3. The purpose of branch-protection.json is not given. It is in the classic form, which the rulesets replace (setup.md:257), and no phase-1 check reads it.
  4. The values of branch-protection.json and the ruleset (each check context, app ID 15368, strict, the review settings) have no record rows. docs/spec/setup.md:320-322 asks a row with one of five source kinds for each value written into the target; a constant in LAYUP's code fits none unless the catalog holds it.
  5. See setup-runner: S13 cannot get a `done` row in phase 1.

### `setup-rule-paths` — The rule-path register from the catalog and the baseline

- **What:** Writes WORK/out/rule-paths.tsv with the phase-1 entries. These are .github/, .githooks/, .gitattributes, AGENTS.md, CLAUDE.md, docs/engineering-discipline.md, docs/issue-workflow.md, docs/ci/, docs/tests/, each `*.sh` file of the tree, docs/gates.tsv, each `config` path of the manifest, docs/setup/, docs/facts/, and docs/guardrails.md with its exception. Each entry has its source.
- **Specification:** docs/spec/setup.md: The rule-path register; docs/spec/records.md: The layout of the records branch (rule-paths.tsv); docs/architecture.md: 6. Native gates and rule protection (Rule protection)
- **Packages:** internal/setup, internal/catalog, internal/tsv
- **Present state:** none
- **Tests:**
  - unit: For an in-memory tree and manifest, the register holds each fixed entry once, each .sh file, and each config path with source `catalog`; docs/guardrails.md has exception `added lines in section 2`; the row order is fixed.
  - integration: On the fixture target after S12, the register reads the real tree, the written docs/gates.tsv and the catalog config paths, and the output matches the rule-paths tsv-schema block.
- **Open questions:**
  1. The source value of each entry is not given. docs/spec/setup.md:171 has `enum(baseline|catalog|architecture)`, but the sources do not say which value each fixed entry gets. For example, docs/gates.tsv and docs/setup/ come from neither the baseline nor the catalog.

### `setup-s15-records-commit` — S15: the verify table, the last record rows and the orphan records commit

- **What:** When WORK/out/verify.tsv is missing, the run stops with the question `O-verify`, which asks the Operator to run `layup setup verify WORK > WORK/out/verify.tsv`. A verify.tsv with a row that is not `pass` or `clear` is exit 1. Otherwise S15 writes the last record rows. It makes in WORK/target the first commit of the orphan branch `layup-records` with README.md (a fixed text: what the branch is, that only `layup run` writes it from Start on, and how to read it with no tool), setup/record.tsv, setup/verify.tsv and rule-paths.tsv. It adds the push of that branch to commands.sh. `steps.tsv` does not go into the target.
- **Specification:** docs/spec/setup.md: The steps (S15); docs/spec/setup.md: Where the records go in phase 1; docs/spec/records.md: NFR-001 — Git is the system of record; docs/spec/records.md: The layout of the records branch (README.md, setup/record.tsv, setup/verify.tsv, rule-paths.tsv); docs/architecture.md: 3. Records and identities (Where the records live; The actor on the forge)
- **Packages:** internal/setup, internal/git, internal/tsv
- **Present state:** none
- **Tests:**
  - unit: No verify.tsv gives a stop row `O-verify` (exit 3); a `fail` or `not-active` row gives exit 1; a verify.tsv that does not match the setup-verify schema gives exit 2.
  - unit: The README text names the branch, says that only `layup run` writes it from Start on, and says how to read it with no tool.
  - integration: The `layup-records` commit has no parent, its tree holds exactly the four paths, its record.tsv equals WORK/out/record.tsv with the S15 done row, and `main` and `layup-setup` do not move.
  - integration: A rerun after S15 makes no second records commit and exits 0 with the step table.
- **Open questions:**
  1. The real baseline cannot pass check `adapted`, so S15 cannot pass. docs/spec/setup.md:251 gives the target the rule `check_adapted`, and :212-213 refuses a verify.tsv with a row that is not `pass` or `clear`. I measured a copy of root commit d2516fd with S05's deletions and the five S07-S09/S14 files removed: `sh docs/setup/setup-check.sh --only adapted` gives 235 failures in 21 files that no phase-1 step writes (docs/engineering-discipline.md 69, docs/ci/README.md 37, docs/tests/README.md 21, docs/issue-workflow.md 11, docs/adr/README.md 10, and others), mostly the words kit, adopter, optional, adapt and your project. No step and no input file of the spec changes these files. A spec fix (input files like the S05 link fixes, or a target form of `adapted`) is needed.
  2. The fixed text of the records branch README.md is not given. docs/spec/setup.md:215-216 gives "what the branch is, that only `layup run` writes it from Start on", and docs/spec/records.md:25 adds "how to read it with no tool". No source gives the words.
  3. The git verbs needed are not all listed. docs/spec/packages.md:39 lists "clone, `ls-remote`, `init`, commit, `rev-parse`, `worktree`, `show`, `diff --name-only`, `apply`". The list has no verb for an orphan branch (`checkout --orphan`, or `commit-tree`/`update-ref`), `add`/`rm`, `checkout` (S02), branch creation (S04), `ls-files` (S10), or a remote for the `origin` that commands.sh pushes to. The plan must say if the list is open or must change.
  4. The commit order of the record is not clear. verify.tsv is made from record.tsv before S15 adds its rows, so the committed setup/verify.tsv judged an earlier record. The sources do not say whether this is accepted.
  5. "The record's last rows" are not given (docs/spec/setup.md:100). docs/architecture.md:189-199 says `layup audit` reads activity "from the setup commits that the records name onward" and allows "the setup commits ... by their SHA". No source says whether S15 writes the SHAs of the setup commits or the head of `layup-setup` into record.tsv.
  6. verify.tsv is not tied to a head. The setup-verify schema (docs/spec/setup.md:268-272) has no head SHA and no record hash, so S15 cannot refuse a verify.tsv made on an older head or in another work area, and it commits that file as evidence.

### `setup-e2e-offline` — End-to-end: a full setup with no network

- **What:** One e2e test builds the binary. It generates a fixture baseline repository in t.TempDir from Go code (git init and one commit). The fixture has history paths, markers (one repeated, one unclosed), a broken-link case, stub discipline and link-lint scripts, an ADR directory, and workflows. The test gives `S01-baseline` as a file:// URL. It runs `layup setup` again and again, answers each stop table, runs `layup setup verify`, then S15. Last, it points `origin` at a local bare repository and runs the git push lines of commands.sh.
- **Specification:** docs/spec/setup.md: Where the records go in phase 1; docs/spec/setup.md: The steps; docs/spec/setup.md: The command `layup setup`; docs/architecture.md: 5. Intake and setup (Scaffold)
- **Packages:** cmd/layup
- **Present state:** cmd/layup/main_e2e_test.go:11-24 builds the binary and tests only `layup version`.
- **Tests:**
  - e2e: With no input files given in advance, the stop sequence is S01 (S01- and Q- questions), S05 (F- link fixes), S07, S08, S09 (one F- file each), S10 (M- markers), S14 (F- files), S15 (O-verify). After the answers, the last run exits 0, and verify.tsv has only pass or clear rows.
  - e2e: After the git push lines run against a local bare repository, it has `main` = the setup head (a fast-forward from the root commit) and an orphan `layup-records`. The root tree equals pin.tree.
  - e2e: Resume and resolve once: a new commit added to the fixture repository between runs does not change pin.commit; two runs of a stop on the same inputs print the same bytes.
- **External inputs:**
  - The tools of the Go entry's active kinds on the CI host, at the versions of the catalog, so that `gate:<kind>` runs with no network.
- **Open questions:**
  1. The fixture must be generated in Go code. MK_EXEMPT (docs/setup/setup-check.sh:305) and AD_EXCLUDE (:359) do not cover internal/setup/testdata/. LAYUP's link-lint reads every tracked *.md except good/ and bad-*/ case directories (docs/links/link-lint.sh:141, :368). Checked-in fixture Markdown with markers or broken links would make LAYUP's own checks fail. The only other choice is to extend these lists, and MK_EXEMPT is also what S10 embeds for targets.
  2. The real baseline cannot be used in CI. The real baseline tree is LAYUP's root commit d2516fd (tree 8ffb250…, docs/setup/armature.pin:7), but the CI job `tests` checks out with the default depth (.github/workflows/ci.yml:305-312, no fetch-depth), so the root commit is absent there.
  3. This test does not cover the failures of `adapted` and `kit-history` on the real baseline (see setup-s05-history and setup-s15-records-commit). The generated fixture can only hold a case that the spec's steps can fix.
  4. Rule 3 of NFR-007 and test code are not clear. docs/spec/packages.md:18-21 says only internal/git starts `git`. The fixture builder starts `git` from test code, and the present e2e test already uses os/exec in cmd/layup (cmd/layup/main_e2e_test.go:6, :14). The sources do not say whether the rule and the future boundary gate apply to test files, or whether a shared fixture builder must go through internal/git.

### `setup-pilot-acceptance` — REQ-002 acceptance on a pilot problem statement

- **What:** The Operator runs `layup setup` and `layup setup verify` on a real pilot problem statement against the real baseline URL. The target passes the baseline's discipline tests and verify, with zero values without a citation. A human audit checks that the cited source of every setup value supports it.
- **Specification:** docs/spec/setup.md: REQ-002 — Set up a target, and prove the setup; docs/spec/setup.md: NFR-003 — No value without evidence (items 1 and 4)
- **Packages:** —
- **Present state:** none
- **Tests:**
  - uat: On the pilot, the count of record rows with no source is zero, and a human audit finds that every row's source supports its value (docs/prd/PRD-0001-layup.md:115 and :134; docs/spec/setup.md:329-330).
- **External inputs:**
  - A pilot problem statement and a pilot target repository (F-0004#11; PRD-0001 criterion of REQ-002, docs/prd/PRD-0001-layup.md:115).
  - The real baseline repository URL and network access on the LAYUP host.
  - The Operator's GitHub login and gh, to run commands.sh.
  - Answers with sources for every marker: on a copy of d2516fd after S05, with the S07-S09/S14 files replaced, check `markers` lists 113 (file, marker) pairs.
  - The five prose files (S07, S08, S09, S14) written by a human or an agent.
- **Open questions:**
  1. When a pilot problem statement exists is not clear. REQ-002 is phase 1 (docs/prd/PRD-0001-layup.md:179), and its criterion names "a pilot problem statement" (:115), but phase 4 ships "the pilot on two problem statements with different stacks" (:184). The plan must say whether this UAT is a phase-1 exit criterion and which problem statement it uses.
  2. The PRD criterion and the spec order disagree. docs/prd/PRD-0001-layup.md:115 asks that "each `human_decision` row's answers are in the target's Git before the next step"; the human_decision rows are S01, S03, S10 and S13 (docs/setup/steps.tsv, last column). The S01 answers go into Git only at S06 (docs/spec/setup.md:91), after S02 to S05, and the S10 answers only at S11 (:96). An Operator decision or a spec fix must settle this before the UAT can pass.
  3. The pilot is blocked on the real baseline by checks `adapted` and `kit-history` (see setup-s15-records-commit and setup-s05-history).

### Constraints of area `setup`

1. **docs/spec/setup.md:17-22** — `layup setup` is a step runner ... with `git` as its only external program. It makes no forge call and starts no role session. *Effect:* No setup issue may add a forge client or gh call. Every forge action is a printed line in commands.sh, and the tests run offline.
2. **docs/spec/packages.md:43 and :49-56** — internal/setup may import internal/tsv, internal/git, internal/catalog; starts no program; internal/cli runs the check of each step from internal/verify after internal/setup did the step. *Effect:* Setup steps take a checker that internal/cli injects, and internal/cli hands the psb gap table to S01 and S06. A setup step's `done` row depends on a verify check, so the verify-area checks must land before or with the steps that cite them.
3. **docs/spec/packages.md:18-22 and :39** — Only internal/git starts `git`; no phase-1 package imports net, net/http or crypto/tls; go.mod has no require line. internal/git verbs: clone, ls-remote, init, commit, rev-parse, worktree, show, diff --name-only, apply. *Effect:* Clone, ls-remote, ls-files, branch and commit go through internal/git only. The verb list must grow (orphan commit, add, checkout, branch, ls-files, remote) in the foundation item, or the spec must say the list is open.
4. **docs/spec/setup.md:51-52; docs/spec/README.md:44-45** — No flag. Every value comes from a file of WORK. A command reads no environment variable for an input, and asks no question on the terminal. *Effect:* The only argument of `layup setup` is WORK. The tests give every value through answers.tsv and inputs/. git must run with no terminal prompt.
5. **docs/spec/setup.md:57-60 and docs/architecture.md:311** — A run is resumable ... the pin is resolved once ("no run resolves the commit again"). *Effect:* The runner and S02 must be built and tested together for resume. Each step needs an idempotent `done` row write.
6. **docs/spec/setup.md:134-140; docs/guardrails.md:294-303** — Any other row ... is an input error: exit 2 ... write one rule for every row of an input that no step of the run asked for. *Effect:* The answers rule is its own item, with a test per row state, written before review. It must exist before S01 and S10.
7. **docs/spec/setup.md:218-223** — commands.sh order: push root (S03), push layup-setup onto main (S13), push layup-records (S15), apply the ruleset (S13). *Effect:* commands.sh is written in a fixed order, not in step order. One writer serves S03, S13 and S15.
8. **docs/spec/setup.md:95 and :277** — LAYUP's MK_EXEMPT ... which the engine embeds at its version; `embed` reads only files under the package's own directory. *Effect:* The Go copy of MK_EXEMPT is a second copy. An integration test must hold it equal to docs/setup/setup-check.sh:305.
9. **docs/spec/setup.md:235-239** — The engine's version of each such check passes and fails on the same fixtures as docs/setup/tests/run.sh. *Effect:* The S10 marker scan and the verify checks need parity tests on docs/setup/tests/<check>/ fixtures, built as temporary Git repositories (integration level).
10. **docs/spec/README.md:75-77** — A Go test of the code reads these blocks and compares them with what the code writes and reads (the test comes with the code, Bootstrap mode rule 1). *Effect:* Each setup item that reads or writes a record (setup-stop, setup-steps, setup-answers, setup-record, rule-paths, open-gaps, catalog-kinds) ships its schema-block test in the same issue, at the integration level.
11. **docs/tests/test-levels.md:43, :48, :53-56, :60, :74** — Unit tests touch no file, network, or clock. Every component has a unit test; every interface or workflow an integration test (`go test -tags=integration`, CI only); every user-facing scenario an E2E test. *Effect:* Unit tests use in-memory trees and a stubbed git. Each step group has an integration test with a real git in a temp directory. The full setup path needs one e2e test.
12. **docs/spec/README.md:61-63** — Two runs on the same input print the same bytes (NFR-005) ... no time, no duration, no scratch path, and its rows have a fixed order. *Effect:* The step and stop tables need a defined row order and a two-run byte test. Commit SHAs in the evidence need a fixed commit identity and date, or the plan must accept that they differ per run.
13. **docs/guardrails.md:276-280** — A test that reads ... a live network ... passes or fails for reasons unrelated to the code. *Effect:* The e2e uses a generated local fixture repository and a file:// URL for git ls-remote and git clone, never the real baseline URL.
14. **docs/setup/setup-check.sh:305, :359; docs/links/link-lint.sh:141, :368** — MK_EXEMPT and AD_EXCLUDE do not cover internal/setup/testdata/; link-lint reads every tracked *.md except good/ and bad-*/ case directories. *Effect:* Fixture baselines with markers, baseline-name mentions or broken links must be generated at test time, or the plan must extend those lists with a reason.
15. **docs/spec/setup.md:79-82 and docs/spec/records.md:82** — steps.tsv is the record of LAYUP's own setup, and it does not change for a target; the engine does not read it. *Effect:* The step table is code that derives from docs/spec/setup.md. No issue edits docs/setup/steps.tsv or reads it at run time.
16. **docs/spec/setup.md:251 and :212-213** — Check `adapted`: `check_adapted`. S15 refuses a verify.tsv with a row that is not `pass` or `clear` (exit 1). *Effect:* On the real baseline, S15 cannot pass until the spec gives a step or input that clears the 235 measured `adapted` failures. The plan must put this spec fix before the pilot.
17. **docs/spec/setup.md:301-302** — LAYUP's own CI runs each fixture of each entry (§6); that test comes with the code (#29). *Effect:* The catalog fixture test belongs to gate-catalog-go (#29), not to a setup item.
18. **docs/spec/setup.md:304-312** — Not in phase 1: the records-branch ruleset, the layup/ required checks, the activation of pending kinds. *Effect:* S13 writes only the default-branch ruleset. These parts go to later-area items.
19. **AGENTS.md (Bootstrap mode) / docs/adr/0012** — A task is a PSB In-Scope item or what one needs; one plan-review comment; one review round. *Effect:* Slice the setup work into issues each small enough for one review round (about 400 lines of diff).

### Risks of area `setup`

- The real baseline blocks the setup. Measured on a copy of LAYUP's root commit d2516fd: check `adapted` gives 235 failures in 21 files that no step writes, and check `kit-history` fails on kit-link lines that S05's rule leaves (7 in backlog.md, 56 in completed.md). S05 and S15 cannot reach `done`, so the pilot UAT cannot pass until the spec changes. The generated e2e fixture hides this.
- The PRD criterion of REQ-002 (answers in Git before the next step) conflicts with the spec order (S01 answers in Git at S06). The UAT cannot pass as written without an Operator decision or a spec fix.
- S14 runs after S10 and S11, so a marker in the baseline's AGENTS.md (line 104 at d2516fd) is asked, filled or kept, and then overwritten, which leaves a stale open-gaps row or an orphan record row and fails `markers` or `sources`.
- Hidden judgement in S05: which backlog and completed lines and notes to delete, and how to find broken links without starting a program in internal/setup. S05 may need its own link checker (the baseline's link-lint.sh is about 350 lines of sh) and may exceed one review-sized issue.
- About 40 open questions block exact code. The biggest are the form of `name` and the remote, `pin.time`, the pin ADR text, the raw-fact names, dates and sources, the required-check set, the git verb list, the S13 `done` row, and the step-table value for an unreached step. Each needs an Operator decision or a spec fix before its issue starts, or the review will find them one round at a time (docs/guardrails.md:294-303).
- A crash in the middle of a step can leave a partial state with no `done` row. A rerun of S02 then resolves the baseline again and breaks the resolve-once rule.
- The host's git configuration can change the result: no user.email fails S03, a global hook or commit.gpgsign blocks the commits, and core.autocrlf can make the root tree differ from pin.tree. Time-based commit SHAs make the evidence differ per run. The target's raw facts also have no line-ending pin (issue #48), so check `facts` can fail on a Windows clone.
- A baseline URL that needs a login can make git ask on the terminal, and the run hangs.
- Raw facts and markers interact: MK_EXEMPT does not exempt docs/facts/, so the second answers record that quotes a marker, or a brief that holds `&lsaquo;`, fails check `markers` or tempts a fill that breaks an immutable fact.
- The CI test job checks out with the default depth (.github/workflows/ci.yml:305-312), so a test cannot use LAYUP's root commit (the real baseline tree) as a fixture without a CI change. A CI change touches a rule path.
- Parity drift: the Go marker scan and the embedded MK_EXEMPT can drift from docs/setup/setup-check.sh. Each needs a test that compares them.
- The real baseline can change between plan and pilot (the latest commit is resolved at setup, O-101). A newer baseline can add rule documents that the fixed rule-path list does not name (known limit, docs/spec/setup.md:180-183), can change the paths of S05, or can add ADRs that change the pin ADR number.
- The setup items depend on many verify-area checks, because every `done` row needs its check. If the verify area slips, the setup steps cannot close except with a stub checker. The e2e also depends on the Go catalog entry (#29) and on its tools being on the CI host.

### Critic points that the reviser of `setup` rejected

- Uncovered (adapted): add an item with `F-<path>` input files for each file that fails `adapted`, or a target form of the check. *Basis:* I confirmed the measurement (235 failures, same per-file counts) on a copy of d2516fd. But no source gives such a mechanism: docs/spec/setup.md:90 gives `F-` inputs only for links that S05 breaks, and :92-99 only for the five prose files; :251 gives `adapted` the plain rule `check_adapted`. So I did not add a build item. I added the point as an open question to setup-s15-records-commit and setup-pilot-acceptance, as a constraint, and as the first risk.
- Uncovered (S13 done row): the same question applies to the record rows of S03 and S15, because their pushes also wait for the Operator. *Basis:* Partly rejected. The evidence of S03 is "root tree = `pin.tree`" (docs/spec/setup.md:88) and the evidence of S15 is "every row of `verify.tsv` is `pass` or `clear`" (:100). Phase 1 can check both, so S03 and S15 can get `done` rows. Only S13's evidence, "the Operator's run of `commands.sh`" (:98), cannot be seen. The S13 question is added to setup-runner and setup-s13-rulesets.
- Uncovered (catalog fixture test): if the gate area does not own the test that applies each fixtures/<kind>.patch, add it to setup-catalog. *Basis:* docs/spec/setup.md:301-302: "LAYUP's own CI runs each fixture of each entry (§6); that test comes with the code (#29)." #29 is the Go entry of the gate area (gate-catalog-go). The test belongs there; I recorded this as a constraint and did not add it to setup-catalog.
- Correction (setup-s05-history open question): completed.md still links github.com/pharzam/armature/ on 32 lines whose T-IDs have no T-*.md file. *Basis:* The point is accepted, but not the number. My count on d2516fd of completed.md lines that link github.com/pharzam/armature/ and name none of the 19 deleted T-IDs is 56 (backlog.md: 7, the same as the critic). The open question gives my numbers and my method.

## Area `verify`

This area has no Go code yet. At 7cdd346 there is only internal/cli (version and psb check) and internal/psb (internal/cli/cli.go:26-41). The reference rules are the sh functions in docs/setup/setup-check.sh:58-571. Their 40 EXPECT cases are under docs/setup/tests/: adapted 4, ci 3, facts 7, frame 3, glossary 2, guardrails 2, identity 1, kit-history 5, markers 4, onboarding 1, pin 4, procedure 2, protection 2. run.sh:44-84 builds each case as a temp repo. The kit/ tree is the root commit and overlay/ is a second commit. Each commit uses a fixed identity with hooks and signing off. A case runs in mode root, self or shallow. The spec (setup.md:225-272) defines 13 named checks, one gate:<kind> row per manifest kind, an output table {check,result,reason}, and exit codes 0/1/2. Nine checks reuse sh rules. Each of these is split in two parts. The first is a root-only core that runs on the sh fixtures, as `sh setup-check.sh ROOT` does. The second is the target addition that the frame runs (record pin rows, NFR-006 method, target name). pin, kit-history, markers, adapted and identity use their own fixture group. They also use the three frame overlays run in root mode, because those overlays hold the only `pin OK` and `identity OK` lines. facts, onboarding, glossary and guardrails have a target form, so 12 sh cases apply only in their generic lines. The harness and kit-history land as one pair, because a harness with no check is not a pass (run.sh:35-36). These are 16 items. One item is new: a stand-in baseline for tests, because CI job `tests` makes a depth-1 clone and the root commit d2516fd is not in it. The largest open issues are these. (1) The unchanged baseline holds rule-1/rule-3 words of check `adapted` in kept files that no phase-1 step rewrites. (2) S12's evidence names a check `gates` that the verify table does not define. (3) The form of the answers record is not fixed, and check `facts` and S06 both need it.

| Key | Title | Requirements | Size | After | Host |
| --- | ----- | ------------ | ---- | ----- | ---- |
| `verify-frame` | The command `layup setup verify WORK`: inputs, scratch tree, table, exit codes, one-check call | REQ-002, NFR-004, NFR-005 | large | verify-harness, found-tsv (foundation area), found-git (foundation area: worktree add --detach, worktree remove, rev-parse), found-cli (foundation area: dispatch and exit-code mapping), found-schema-test (foundation area: Go test that reads the tsv-schema blocks) | row 7 |
| `verify-harness` | The shared fixture harness: run docs/setup/tests through the Go check cores | REQ-002 | small | verify-kit-history (same pull request: the harness is not green with no check, run.sh:35-36, :96), found-git (foundation area: init, add, commit --allow-empty, clone --depth 1, with -c user.name/user.email/core.hooksPath/commit.gpgsign) | row 7 |
| `verify-pin` | Check `pin`: the check_pin core, and the target's record rows and method | REQ-002, NFR-006 | small | verify-frame, verify-harness, found-git (foundation area: rev-list --max-parents=0, rev-parse --is-shallow-repository, rev-parse <commit>^{tree}; packages.md:39 does not list rev-list) | row 7 |
| `verify-kit-history` | Check `kit-history`: the first check, landed with the harness | REQ-002 | small | found-git (foundation area) | row 7 |
| `verify-markers` | Check `markers`: the scanner with LAYUP's MK_EXEMPT, and open-gaps by its schema | REQ-002, NFR-003 | small | verify-frame, verify-harness, found-git (foundation area: ls-files -z or with core.quotePath=false; packages.md:39 does not list ls-files) | row 10 |
| `verify-adapted` | Check `adapted` (rules 1 to 3), with AD_EXCLUDE and ad_allowed | REQ-002 | large | verify-frame, verify-harness, found-git (foundation area: ls-files) | row 11 |
| `verify-identity` | Check `identity`, with the target's name | REQ-002 | small | verify-frame, verify-harness | row 7 |
| `verify-facts` | Check `facts` in a target's form, the answers-record reader, and the fact resolver | REQ-002, NFR-003 | large | verify-frame, verify-harness | row 12 |
| `verify-onboarding-glossary` | Checks `onboarding` and `glossary` in a target's form | REQ-002, NFR-003 | small | verify-facts | row 12 |
| `verify-guardrails` | Check `guardrails` in a target's form | REQ-002, NFR-003 | small | verify-facts | row 12 |
| `verify-baseline-scripts` | Checks `discipline-tests` and `link-lint`: the baseline's own scripts | REQ-002, NFR-004 | small | verify-frame, verify-test-baseline | row 10 |
| `verify-test-baseline` | A stand-in baseline for the integration and e2e tests, with no network | REQ-002 | small | found-git (foundation area) | row 7 |
| `verify-sources` | Check `sources`: every value of the setup record has a source that resolves | NFR-003, REQ-002 | small | verify-frame, verify-facts, verify-markers, setup-catalog (setup area: internal/catalog with the embedded Go entry) | row 10 |
| `verify-jobs` | Check `jobs`: one CI job per gate kind | REQ-002, REQ-007 | small | verify-frame, gate-manifest (gate area: the manifest schema reader), setup-catalog (setup area: the CI workflow file of the Go entry) | row 15 |
| `verify-gate-fixtures` | Rows `gate:<kind>`: the clean run and the known-bad fixture | REQ-002, REQ-004, REQ-007, NFR-004, NFR-005 | large | verify-frame, gate-run (gate area: an internal/gate call that runs a manifest on base and head and returns the rows), gate-manifest (gate area: the reader of docs/gates.tsv with kinds and states), found-git (foundation area: apply, commit on a detached scratch tree or commit-tree, rev-parse), setup-catalog (setup area: the Go entry's kinds.tsv and fixtures/<kind>.patch) | row 15 |
| `verify-acceptance` | End to end: `layup setup verify` on a WORK made by `layup setup` | REQ-002, NFR-003, NFR-004 | large | verify-frame, verify-pin, verify-kit-history, verify-facts, verify-onboarding-glossary, verify-guardrails, verify-markers, verify-adapted, verify-identity, verify-baseline-scripts, verify-sources, verify-jobs, verify-gate-fixtures, verify-test-baseline, setup-runner (setup area: `layup setup WORK` through S15) | row 16 |

### `verify-frame` — The command `layup setup verify WORK`: inputs, scratch tree, table, exit codes, one-check call

- **What:** `layup setup verify WORK` reads WORK/inputs/answers.tsv and WORK/out/record.tsv by their schemas. It opens a scratch work tree of the head of `layup-setup`, runs the 13 checks and the gate:<kind> rows in a fixed order, and prints the setup-verify table (no row for the rulesets read back, the probes or the layup/ checks, which are not in phase 1). The exit is 0 when each row is pass or clear, 1 on fail or not-active, and 2 on a usage or input error. It also gives internal/cli a call that runs one or more named checks, as the evidence of a step of `layup setup` (S04 to S14), including the name `gates` of S12.
- **Specification:** docs/spec/setup.md — The checks of `layup setup verify` (command, reads, no ref changed, schema setup-verify, row 'the rulesets read back; the probes of §3'); docs/spec/setup.md — Not in phase 1 (the layup/ required checks); docs/spec/setup.md — The answers (schema setup-answers); docs/spec/setup.md — The setup record (schema setup-record); docs/spec/setup.md — The steps (evidence column of S04 to S14); docs/spec/setup.md — Where the records go in phase 1 (item 2: O-verify); docs/spec/README.md — Commands (arguments, output, exit codes, determinism); docs/spec/packages.md — The table of phase 1 (internal/verify, internal/cli; the Decided-here paragraph); docs/spec/gate.md — NFR-004 (items 1, 2, 4); docs/spec/gate.md — NFR-005 (item 2); docs/architecture.md — 5. Intake and setup (Scaffold 5)
- **Packages:** internal/verify, internal/cli, internal/tsv, internal/git
- **Present state:** none. internal/cli/cli.go:31-40 knows only `version` and `psb`. internal/verify, internal/tsv and internal/git do not exist.
- **Tests:**
  - unit: The table writer prints the header check<TAB>result<TAB>reason, `—` for a pass reason, the rows in the fixed order, and the exit code from the rows (0 only when every row is pass or clear).
  - unit: The table has only the 13 named checks and the gate:<kind> rows. No row exists for ci, procedure, protection, the rulesets read back, the probes or layup/verify.
  - integration: The schema test: the setup-verify, setup-answers and setup-record blocks of docs/spec/setup.md equal what the code writes and reads (it reads a repository file, so it is integration, as internal/psb/check_integration_test.go:1).
  - integration: On a real temp WORK repo the command makes and removes a scratch work tree. After the run, `git for-each-ref` of WORK/target and the hashes of the files in WORK/out and WORK/inputs are unchanged.
  - integration: An input error gives exit 2: WORK/target missing, no `layup-setup` branch, or an answers.tsv or record.tsv that does not match its schema.
  - integration: The one-check call that internal/cli uses: one name, two names (S05: kit-history and link-lint; S11: markers, sources and facts; S12: jobs and gates), and an unknown name (an error, not pass).
  - e2e: The built binary run twice on one WORK prints the same bytes (NFR-005), and `layup setup verify` with an extra argument or a flag exits 2.
- **Open questions:**
  1. setup.md:97 gives S12 the evidence "checks `jobs` and `gates`", but the verify table (setup.md:243-258) has no check `gates`, only `gate:<kind>` rows. architecture.md:385-386 gives S12 "the native gate jobs exist, one per kind". Proposed reading: `gates` means every gate:<kind> row. No source states this.
  2. architecture.md:385 says the evidence of S04 to S14 is "`layup setup verify <check>` OK", but setup.md:231 says "No flag" and the command is `layup setup verify WORK`. Is one check a CLI argument, or only the Go call that internal/cli makes (packages.md:53-56)?
  3. The row order is not stated. setup.md:269 says only "the check name of the table above". Is it the order of the table at setup.md:243-256, with the gate:<kind> rows in manifest order?
  4. setup.md:231 says verify reads record.tsv. When internal/cli runs one check during `layup setup`, record.tsv is partial. What do checks `sources` and `facts` assert on a partial record (see verify-facts for the S06 case)?

### `verify-harness` — The shared fixture harness: run docs/setup/tests through the Go check cores

- **What:** A Go integration test reads each docs/setup/tests/<check>/<case>/EXPECT. It builds the same temp repo as run.sh: the kit/ tree is the root commit and overlay/ is a second commit, each with a fixed identity, hooks off and signing off, and a --depth 1 clone for mode=shallow. Then it calls the root-only core of the Go check of that name and asserts the EXPECT exit and each `line=`. It also runs the overlays of frame/bad-missing-linter, bad-passthrough and good-passthrough in root mode. For these it asserts only the lines of pin, kit-history, markers, adapted and identity, not the case exit. Cases that the Go engine does not run (ci, procedure, protection, the kit-linters lines) and the LAYUP-form lines are skipped from a fixed list that gives the reason. It lands in one pull request with verify-kit-history, its first check.
- **Specification:** docs/spec/setup.md — The checks of `layup setup verify` (paragraph: the engine's version of each check passes and fails on the same fixtures as docs/setup/tests/run.sh); docs/prd/PRD-0001-layup.md — 10. Risks (Two checks that drift); docs/adr/0011-structure-the-core-engine-as-a-go-cli-over-repository-files.md — Consequences
- **Packages:** internal/verify, internal/git
- **Present state:** docs/setup/tests/run.sh:1-98 (sh runner; fixed identity at :44-45). 40 EXPECT cases. The shared root tree is docs/setup/tests/kit/README.md. frame/good-passthrough/overlay holds a full LAYUP-form tree. No Go harness exists.
- **Tests:**
  - integration: Each applicable case gives the EXPECT exit, and each `line=` is among the Go findings when they are written in the sh form `setup-check: <check> FAIL <cause>: <detail>` (or `<check> OK`). The cases are pin 4, kit-history 5, markers 4, adapted 4 and identity 1, and the generic lines of facts, onboarding and guardrails.
  - integration: The three frame overlays in root mode give `pin OK`, and good-passthrough also gives `kit-history OK`, `markers OK`, `adapted OK` and `identity OK`. So a Go pin or identity check that always fails does not pass the harness.
  - integration: With no case directory the harness fails, and a skip-list entry that names a case that does not exist fails.
  - integration: The commits of the harness succeed on a host with no global git identity and with a global hook or signing set (HOME set to an empty temp dir).
- **Open questions:**
  1. How can the Go table match the fixtures? setup.md:271 gives `reason` as "the first failure" only, but EXPECT files assert every FAIL line (for example identity/bad-kit has 4 lines). Proposed reading: each Go check core returns all findings in the sh line form, the harness compares them, and the table shows the first. No source states this.
  2. setup.md:238 says "passes and fails on the same fixtures" with no exception, but setup.md:260-266 drops LAYUP's counts for facts, onboarding, glossary and guardrails. Are the 12 LAYUP-form cases of these checks (or their LAYUP-form lines) skipped, or run against a Go LAYUP-form variant?
  3. packages.md:18-19 says "Only the package `internal/git` starts the `git` program". Does that rule cover _test.go files that build fixture repos? If it does, internal/git needs `add`, `clone --depth 1` and `commit --allow-empty`, and packages.md:39 does not list `add`.
  4. frame/bad-missing-linter expects `kit-linters FAIL absent`, but setup.md:270 calls a missing script `not-active`. The harness skips the kit-linters lines. Is that the intended reading?
  5. CI job setup-check restores docs/setup/setup-check.sh and docs/setup/tests/ from the default branch (ci.yml:270-288), but job tests runs the pull request's own Go harness and fixtures. So a fixture change can pass the Go harness on the pull request and break the sh run on main only after the merge. Is this a recorded known limit?

### `verify-pin` — Check `pin`: the check_pin core, and the target's record rows and method

- **What:** The core is a port of check_pin on a repository root. Each of source, commit, tree, method and date appears once. The commit is 40 hex characters. The clone is not shallow. There is one root commit, and its tree equals the pinned tree. The frame adds the target part: the pin values equal the record rows pin.source, pin.commit and pin.tree and the date of pin.time, and the method equals the NFR-006 text with the source and the commit.
- **Specification:** docs/spec/setup.md — The checks of `layup setup verify` (row pin); docs/spec/setup.md — NFR-006 — The baseline at a pinned, recorded version; docs/spec/setup.md — The steps (S02, S04); docs/architecture.md — 5. Intake and setup (Scaffold 5: the root tree equals the pinned tree)
- **Packages:** internal/verify, internal/git
- **Present state:** docs/setup/setup-check.sh:58-91 (check_pin). Fixtures docs/setup/tests/pin/ (4 cases). LAYUP's own pin has method=npx degit (docs/setup/armature.pin:8), which is not the target form.
- **Tests:**
  - unit: Pin-file parsing (key counts, the commit format), and the target part (a comparison to the record rows and to the method text) with stubbed git answers.
  - integration: Through the harness, the core gives the lines of the four pin cases (bad-commit, bad-missing, bad-shallow, bad-tree) and `pin OK` on the three frame overlays, whose method is npx degit.
  - integration: New target-form fixtures through the frame: a pin whose commit differs from record row pin.commit fails, and a method in the npx-degit form fails.
- **Open questions:**
  1. setup.md:244 says "the pin file's values equal the record's pin rows". The pin file has `date=YYYY-MM-DD` and the record has `pin.time`. Is the date the UTC date of pin.time?

### `verify-kit-history` — Check `kit-history`: the first check, landed with the harness

- **What:** Port of check_kit_history on a repository root: docs/decisions/ and docs/audit/ are absent, each docs/tasks/T-*.md has its ID in backlog.md or completed.md, and no task index links the baseline's repository. It also defines the shape of a check core that the other checks use: a root in, findings in the sh line form out.
- **Specification:** docs/spec/setup.md — The checks of `layup setup verify` (row kit-history); docs/spec/setup.md — The steps (S05)
- **Packages:** internal/verify
- **Present state:** docs/setup/setup-check.sh:93-110. Fixtures docs/setup/tests/kit-history/ (5 cases).
- **Tests:**
  - unit: The orphan-task and kit-link rules on in-memory trees.
  - integration: The five kit-history cases give the same exit and lines through the harness, and good-passthrough gives `kit-history OK`.
- **Open questions:**
  1. setup-check.sh:107 matches the literal `github.com/pharzam/armature/`. setup.md:245 says only `check_kit_history`. For a target whose S01-baseline is another repository, does the kit-link rule use the record's pin.source or the literal?

### `verify-markers` — Check `markers`: the scanner with LAYUP's MK_EXEMPT, and open-gaps by its schema

- **What:** Port of check_markers. It lists each marker in git-tracked files outside MK_EXEMPT and skips the code-span `&lsaquo;` mention and `&lsaquo;…&rsaquo;`. Each found marker needs a row in docs/setup/open-gaps.tsv, and each row needs its marker in the tree. open-gaps.tsv is read by its schema (no header; key file and marker). A row whose question is empty or `—` fails. The same scanner gives the marker list of S10.
- **Specification:** docs/spec/setup.md — The checks of `layup setup verify` (row markers); docs/spec/setup.md — The files of `docs/setup/` in a target (schema open-gaps: question never `—`); docs/spec/README.md — the empty-field form `—`; docs/spec/setup.md — The steps (S10, S11); docs/spec/setup.md — NFR-003 — No value without evidence (item 3)
- **Packages:** internal/verify, internal/git, internal/tsv
- **Present state:** docs/setup/setup-check.sh:294-337 (MK_EXEMPT at :305; the question test at :328 fails only an empty third field). Fixtures docs/setup/tests/markers/ (4 cases).
- **Tests:**
  - unit: The marker scanner: a closed marker, a marker that runs to the line end, the code-span mention, `&lsaquo;…&rsaquo;`, two markers on one line, and UTF-8 paths.
  - unit: The open-gaps rules: a question that is empty, `—`, a space or `\r` fails; a duplicate key fails.
  - integration: Drift guard: the embedded MK_EXEMPT equals the MK_EXEMPT line of docs/setup/setup-check.sh:305, read at test time.
  - integration: The open-gaps tsv-schema block of setup.md is in the schema test and equals what the code reads.
  - integration: The four markers cases, including docs/é.md, give the same exit and lines through the harness, and good-passthrough gives `markers OK`.
- **Open questions:**
  1. S10 (setup.md:95) lists markers in internal/setup, which imports only tsv, git and catalog (packages.md:43). Where does the one scanner live so that S10 and check markers cannot drift? One reading: internal/cli runs verify's scanner and hands the list to setup, "as it does for the step checks" (setup.md:86).
  2. Open issue #21 (.ctx/issues/21.md items 1 and 5): the sh check skips a file whose name git quotes, and accepts a space or `\r` as a question. Should the Go port be correct (and differ from sh on these inputs), or should #21 fix sh first?

### `verify-adapted` — Check `adapted` (rules 1 to 3), with AD_EXCLUDE and ad_allowed

- **What:** Port of check_adapted over the git-tracked *.md files outside AD_EXCLUDE. Lines are joined per paragraph, and a marker counts as one unit. There are 16 hit patterns. Rule-1 has 3: kit, adopter, the template. Rule-2 has 1: Armature, except in a file of ad_allowed. Rule-3 has 12. Each hit is reported at the line where it starts.
- **Specification:** docs/spec/setup.md — The checks of `layup setup verify` (row adapted)
- **Packages:** internal/verify, internal/git
- **Present state:** docs/setup/setup-check.sh:339-435 (16 hit() calls at :398-413). Fixtures docs/setup/tests/adapted/ (4 cases).
- **Tests:**
  - unit: Each of the 16 hit patterns, with its word boundaries, the kit-history and kit-linters exemption, the `#`/`-` anchor exemption, and the line number of a phrase that crosses a line end.
  - integration: Drift guard: the embedded AD_EXCLUDE and ad_allowed equal setup-check.sh:359 and :361-373, read at test time.
  - integration: The four adapted cases (bad-rule-1, bad-rule-2, bad-rule-3, good-allowed) give the same exit and lines through the harness, and good-passthrough gives `adapted OK`.
- **External inputs:**
  - A decision of the Operator on the unchanged baseline (first open question).
- **Open questions:**
  1. `git grep` on LAYUP's root commit d2516fd (the unchanged baseline) finds rule-1/rule-3 words in kept files that no phase-1 step rewrites (docs/engineering-discipline.md, docs/issue-workflow.md, docs/tests/*.md, docs/ci/README.md, .githooks/README.md and others). Phase-1 setup writes only S07-S09 and S14 prose (setup.md:92-99). So setup.md:251 ("`adapted` | yes | `check_adapted`") may fail every correct phase-1 setup. Is a newer baseline already adapted, does setup need more prose inputs, or is adapted a later-phase check for a target?
  2. AD_EXCLUDE (setup-check.sh:359) and ad_allowed (:361-373) name LAYUP paths and LAYUP reasons. Does a target use these lists as they are? setup.md:251 says only `check_adapted`.

### `verify-identity` — Check `identity`, with the target's name

- **What:** The core is a port of check_identity: README.md and AGENTS.md hold none of the three phrases that say the repository is the baseline template, and README.md links docs/setup/armature.pin. The frame adds what the target's name requires (open question).
- **Specification:** docs/spec/setup.md — The checks of `layup setup verify` (row identity); docs/spec/setup.md — The steps (S14)
- **Packages:** internal/verify
- **Present state:** docs/setup/setup-check.sh:559-571. Fixture docs/setup/tests/identity/bad-kit (the only identity case; the only `identity OK` line is in frame/good-passthrough/EXPECT).
- **Tests:**
  - unit: Each phrase and the pin link, on in-memory files.
  - integration: Drift guard: the phrases equal setup-check.sh:565, read at test time.
  - integration: identity/bad-kit gives the same exit and four lines through the harness, and good-passthrough gives `identity OK`.
- **Open questions:**
  1. setup.md:252 says "`check_identity`, with the target's name", but check_identity uses no name (setup-check.sh:562-571). What does the record row `name` change: one more rule that README.md names the target, or nothing?

### `verify-facts` — Check `facts` in a target's form, the answers-record reader, and the fact resolver

- **What:** Each brief is in docs/facts/ and in docs/setup/facts.sha256 with a hash that matches. Each answers record also has its line, and a path with two lines fails. Each question ID of answers.tsv is a fact in exactly one answers record: S01- and Q- IDs in the S06 record, M- IDs in the S11 record. The facts index has a row per record. The item fixes the text form of an answers record, which S06 and S11 write. It also gives the F-NNNN#n resolver that onboarding, glossary, guardrails and sources use.
- **Specification:** docs/spec/setup.md — The checks of `layup setup verify` (row facts; Decided here: the target's form); docs/spec/setup.md — The files of `docs/setup/` in a target (facts.sha256: one line per raw facts file, the key is the path); docs/spec/setup.md — The steps (S06, S11); docs/facts/README.md — Fact IDs; docs/architecture.md — 5. Intake and setup (Scaffold 3)
- **Packages:** internal/verify, internal/tsv
- **Present state:** docs/setup/setup-check.sh:112-207. LAYUP-specific parts: the counts F-0001:39 and F-0003:75 (:150), F-0004 with 19 rows of internal/psb/testdata/psb.tsv (:177-199), the index IDs F-0001 to F-0004 (:202), and the required names problem-statement-brief.md and architectural-vision-brief.md (:142-145). The generic part is the hash loop (:128-141). The record form today is docs/facts/template.md, and F-0004 is LAYUP's own answers record. Fixtures docs/setup/tests/facts/ (7 cases, 6 of them LAYUP-form).
- **Tests:**
  - unit: The sha256sum-line parser (a last line with no newline, a line with no path, a path on two lines), the answers-record parser, the exactly-one-record rule per question ID, the placement rule of S01-, Q- and M- IDs, and the resolver, on in-memory records.
  - integration: The generic lines of facts/bad-hash (the hash and path lines) match through the harness.
  - integration: New target-form fixtures: a question ID in two records fails, an M- ID in the S06 record fails, a record with no index row fails, an answers record with no facts.sha256 line fails, and a correct S06-only target (no marker) passes.
- **Open questions:**
  1. The text form of an answers record is not fixed. setup.md:91 says only "one fact per question ID, each with the question ID, the question text and the answer". Where in a numbered fact does the question ID stand, so that the check can read it? The check and S06 (setup area) both need this form. If this item fixes the form, S06 depends on verify-facts and the cycle breaks.
  2. When internal/cli runs check `facts` as the evidence of S06 (setup.md:91), S10 and S11 have not run. An M- row given early in answers.tsv has no S11 record yet, so the full rule of setup.md:246 fails a correct S06. What does the check assert at S06?
  3. How does verify find "the record of S06" and "the record of S11"? setup.md:91 and :96 give no F-NNNN number or file name.
  4. setup.md:91 copies each brief "into docs/facts/" but gives no file names. check_facts requires problem-statement-brief.md and architectural-vision-brief.md (setup-check.sh:142-145), and a target's vision brief exists only "when there is one" (setup.md:68). Which names does the target check expect?
  5. Open issue #61: the sh empty-fact test misses a tab. Does an emptiness rule apply to an answer fact in the target form (setup.md:246 states none)?
  6. Open issue #48: on an autocrlf=true host, a checkout adds CR bytes and the hash fails. Should verify hash the blob (`git show`) and not the work-tree file?

### `verify-onboarding-glossary` — Checks `onboarding` and `glossary` in a target's form

- **What:** onboarding: docs/onboarding-for-engineers.md exists, holds no marker, links the problem statement, and each F-NNNN#n that it cites resolves. glossary: each F-NNNN#n that docs/glossary.md cites resolves. LAYUP's glossary heading and its 25 rows do not apply.
- **Specification:** docs/spec/setup.md — The checks of `layup setup verify` (rows onboarding, glossary; Decided here: the target's form); docs/spec/setup.md — The steps (S07, S08)
- **Packages:** internal/verify
- **Present state:** docs/setup/setup-check.sh:209-249. LAYUP-specific parts: the F-0001 record only (:219-224), the heading '## 1. LAYUP domain (PSB §8)' (:235), 25 rows (:242), and F-0001#15 to #39 (:243-248). Generic parts: the marker test (:216) and the link test (:217).
- **Tests:**
  - unit: Citation extraction for any F-NNNN#n, the marker test and the link test, on in-memory files.
  - integration: onboarding/bad-all gives its marker and link lines through the harness. New target fixtures: a citation to a fact that does not exist fails, and a glossary with any heading and any row count passes.
- **Open questions:**
  1. setup.md:247 says "links the problem statement", and the sh test looks for the exact text `](facts/problem-statement-brief.md)` (setup-check.sh:217). The target brief's file name is not given (see verify-facts). Is the link test on the path that S06 writes?
  2. glossary/bad-no-section and bad-rows are LAYUP-form only. Are they on the harness skip list?

### `verify-guardrails` — Check `guardrails` in a target's form

- **What:** Each entry's `Check:` value in docs/guardrails.md is `no check yet`, or `<existing path> (<gate>)` with gate `hook` or `ci:<job>`, as in check_guardrails. Each F-NNNN#n citation resolves. LAYUP's count of 9 invariants does not apply.
- **Specification:** docs/spec/setup.md — The checks of `layup setup verify` (row guardrails; Decided here: the target's form); docs/spec/setup.md — The steps (S09)
- **Packages:** internal/verify
- **Present state:** docs/setup/setup-check.sh:251-292. LAYUP-specific parts: the heading '### 1.1 LAYUP System Invariants (PSB §6)' (:259), Inv-1 to Inv-9 (:272-279), and entry N cites F-0001#N (:280). Generic part: the Check: rule (:281-290).
- **Tests:**
  - unit: Parsing a Check: value: `no check yet`, a path with hook, a path with ci:<job>, a path that does not exist, a bad gate, and no value.
  - integration: The Check: lines of guardrails/bad-entries match through the harness (cron, `maybe later`, docs/nothing.sh, docs). New target fixture: an entry whose citation does not resolve fails.
- **Open questions:**
  1. Without LAYUP's heading and its `- **Inv-N**` entry pattern (setup-check.sh:259, :269), what is an "entry" of a target's guardrails.md (setup.md:249)? Every line with `Check:`, every bullet of §1, or the baseline's own entry form?

### `verify-baseline-scripts` — Checks `discipline-tests` and `link-lint`: the baseline's own scripts

- **What:** In the scratch tree, verify runs `sh docs/tests/run-discipline-tests.sh` and `sh docs/links/link-lint.sh` of the target. Exit 0 gives pass, and another exit gives fail with reason `exit <code>` or `signal <name>`. A missing script, or no `sh`, gives not-active with a fixed reason, never pass. The script output goes to standard error and never into the table. A run of more than ten seconds shows on standard error which script runs.
- **Specification:** docs/spec/setup.md — The checks of `layup setup verify` (rows discipline-tests, link-lint; schema setup-verify: not-active for a missing script); docs/spec/packages.md — The table of phase 1 (internal/verify starts the baseline's own check scripts with sh); docs/spec/gate.md — NFR-004 (items 1, 3, 5); docs/spec/README.md — Commands (no scratch path in a table); docs/architecture.md — 5. Intake and setup (Scaffold 5)
- **Packages:** internal/verify
- **Present state:** The sh pass-through runs these scripts with no argument, and fails an absent one (setup-check.sh:614-624). Fixtures frame/bad-missing-linter, bad-passthrough and good-passthrough (mode=self).
- **Tests:**
  - unit: The result mapping: exit 0 is pass, exit 1 is fail `exit 1`, a signal is fail `signal <name>`, a missing file is not-active, and sh not found is not-active. No reason holds a path of the scratch tree.
  - integration: With stub scripts in a temp scratch tree: an OK stub gives pass, a FAIL stub gives fail, and an absent stub gives not-active with exit 1 (the NFR-004 fixture: a check that does not run cannot produce pass).
  - integration: On the stand-in baseline of verify-test-baseline, both scripts run from a git checkout that keeps .gitattributes, so the crlf fixtures of run-discipline-tests.sh pass.
- **External inputs:**
  - The baseline's run-discipline-tests.sh and link-lint.sh at a target's pin (the latest baseline commit, O-101). Their run time and needs are known only for LAYUP's pin.
- **Open questions:**
  1. setup.md:271 does not give the reason of a fail row of these two checks. Proposed: `exit <code>`, as gate.md:74 does for a gate kind. No source states this for verify.
  2. The scratch tree comes from `git worktree add`. On a host with core.autocrlf=true, the checkout changes line ends (#48). Should verify set core.autocrlf=false on its scratch checkout?

### `verify-test-baseline` — A stand-in baseline for the integration and e2e tests, with no network

- **What:** The tests get a local baseline repository without the network and without LAYUP's history. CI job `tests` makes a depth-1 clone, so LAYUP's root commit d2516fd is not there. The stand-in is either a test-data tree or a tree that the test builds. It is the input of `S01-baseline` as a file:// URL.
- **Specification:** docs/spec/setup.md — The steps (S01, S02: git ls-remote and git clone of S01-baseline); docs/tests/test-levels.md — 2. Integration, 3. End-to-end
- **Packages:** internal/verify
- **Present state:** none. .github/workflows/ci.yml:309 (job tests) uses actions/checkout@v4 with no fetch-depth; only jobs setup-check (:263) and security set fetch-depth: 0.
- **Tests:**
  - integration: The stand-in builds in a temp dir in job tests (depth-1 clone, no network), and `git ls-remote file://<path> HEAD` gives one commit.
  - discipline: LAYUP's own `sh docs/setup/setup-check.sh` (markers, adapted) still passes with the stand-in in the tree.
- **External inputs:**
  - A decision of the Operator: change ci.yml job tests to fetch-depth: 0 (a change to a gate, with its own cycle cap), or keep a stand-in tree.
- **Open questions:**
  1. A test-data tree that copies baseline files holds markers and kit words. MK_EXEMPT (setup-check.sh:305) and AD_EXCLUDE (:359) exempt only internal/psb/testdata/ and docs/*/tests/. So LAYUP's own `markers` and `adapted` fail unless setup-check.sh changes (a gate rule path), or the test writes the tree at run time. Which one?
  2. Is a file:// URL an allowed `S01-baseline` answer in tests (setup.md:87 names `git ls-remote <url>`)?

### `verify-sources` — Check `sources`: every value of the setup record has a source that resolves

- **What:** Every value row of record.tsv (not a done row) has a source. An answer ref is a row of answers.tsv. A catalog ref `<stack>/<path>` is a file of the embedded catalog entry. A fact ref resolves in docs/facts/. A gap row has its marker in the tree and its row in open-gaps.tsv. A computed ref has a command or `sha256 <path>`.
- **Specification:** docs/spec/setup.md — The checks of `layup setup verify` (row sources); docs/spec/setup.md — The setup record (schema setup-record: source and ref); docs/spec/setup.md — NFR-003 — No value without evidence (items 1 to 3); docs/architecture.md — 5. Intake and setup (Scaffold 1, Scaffold 5)
- **Packages:** internal/verify, internal/tsv, internal/catalog
- **Present state:** none (no sh counterpart).
- **Tests:**
  - unit: Each source kind resolves or fails: an answer ref with no answers row, a catalog ref to a missing file, a fact ref to a missing fact, and a gap row whose marker is gone or has no open-gaps row.
  - integration: A WORK whose record.tsv has one value row with an empty or unresolved source gives `sources fail` with that row in the reason, and exit 1.
- **Open questions:**
  1. setup.md:254 gives rules for answer, catalog, fact and gap, and none for `computed` (setup.md:159-160). Is a computed row checked (for example by re-hashing a `sha256 <path>`), or accepted as it is?

### `verify-jobs` — Check `jobs`: one CI job per gate kind

- **What:** For each kind of docs/gates.tsv, a workflow under .github/workflows/ has a job with the kind's name. A kind with no job gives fail and names the kind. A missing or malformed manifest gives an input error.
- **Specification:** docs/spec/setup.md — The checks of `layup setup verify` (row jobs); docs/spec/setup.md — The steps (S12); docs/spec/gate.md — The gate manifest; docs/spec/gate.md — REQ-007 (In phase 1, item 1); docs/architecture.md — 6. Native gates and rule protection (The stack gates)
- **Packages:** internal/verify, internal/tsv
- **Present state:** none. The nearest sh code is the job-name reader of check_protection (setup-check.sh:538-550), which takes the job's `name:`, else its id.
- **Tests:**
  - unit: The job-name reader on workflow text (`jobs:` keys, comments, CR line ends), with no YAML library (NFR-007).
  - integration: A target tree with one manifest kind and no job fails, and the Go catalog's workflow, as S12 writes it, passes.
- **Open questions:**
  1. "a CI job with the kind's name" (setup.md:255, gate.md:19): is that the job id under `jobs:`, or the job's `name:` (the check-run name that S13's ruleset requires)? check_protection uses `name:`, else the id.
  2. The check does not test NFR-002 item 2 (records.md:127-130: the gate jobs run the commands of the target's own docs/gates.tsv). Is that part of `jobs`, or only a review?

### `verify-gate-fixtures` — Rows `gate:<kind>`: the clean run and the known-bad fixture

- **What:** For each manifest kind there is one row. An active kind runs internal/gate with base and head both equal to the setup head, which must give pass or clear. Then the kind's catalog fixture is applied with `git apply` on the setup head, in a separate tree that the other checks do not read. It is committed with a fixed identity, hooks off and signing off, as an object on no ref. A run with that commit as head must give fail. The first rule that matches decides (setup.md:256). A pending kind gives clear, reason `pending: fixture not run`. Progress for each run goes to standard error.
- **Specification:** docs/spec/setup.md — The checks of `layup setup verify` (row gate:<kind>); docs/spec/setup.md — The stack catalog (fixtures/<kind>.patch, schema catalog-kinds); docs/spec/gate.md — The command (The run, step 2: config paths from the base; the result table; output to standard error); docs/spec/gate.md — Not in phase 1 (the known-bad fixtures at setup are run by layup setup verify); docs/spec/gate.md — NFR-004; docs/spec/packages.md — The table of phase 1 (internal/verify imports internal/gate); docs/architecture.md — 5. Intake and setup (Scaffold 5); docs/architecture.md — 6. Native gates and rule protection (Detection, as the complement)
- **Packages:** internal/verify, internal/gate, internal/git, internal/catalog
- **Present state:** none. internal/gate and internal/catalog do not exist (issue #33 is open, .ctx/issues/33.md).
- **Tests:**
  - unit: The decision table from a pair of gate results to a row: each line of setup.md:256 in order, including clean fail with fixture not-active (fail wins).
  - integration: On a temp target with a stub manifest (kind `t`, command `! test -f bad`) and a fixture patch that adds `bad`, the row is pass. After the run, `git for-each-ref` is unchanged, and the fixture commit is an object that no ref names.
  - integration: The NFR-004 fixture: a manifest whose tool does not exist gives `gate:<kind> not-active`, reason `tool not found: <tool>`, and exit 1. A fixture that does not make the kind fail gives fail, reason `fixture not detected`. A pending kind gives clear, reason `pending: fixture not run`.
  - integration: A fixture patch that changes only a `config` path of its kind gives `fixture not detected`, because gate.md:59-62 writes the config from the base (the setup head). Each catalog fixture changes a product path.
  - integration: A kind whose run passes ten seconds prints the kind and a live sign on standard error, and nothing of that output is in the table.
  - e2e: The binary on a WORK set up with the Go catalog entry prints one gate:<kind> row per manifest kind, in manifest order, and the same bytes on a second run.
- **External inputs:**
  - The Go toolchain and each gate tool at the version that the catalog's `evidence` column documents, on the CI host and the LAYUP host.
- **Open questions:**
  1. setup.md:256 does not say what happens when `git apply` of the fixture fails on the setup head. Is it fail, not-active, or an input error (exit 2)?
  2. catalog-kinds (setup.md:296) allows `fixture` `—` only for a pending kind. Is an active kind with `—` an input error of the catalog (exit 2) or a fail row?
  3. How is the commit on no ref made: `git commit` on a detached scratch tree, or `git commit-tree`? packages.md:39 lists commit and apply, not commit-tree or write-tree.
  4. The clean run has base equal to head, so the pending rule and the config of the scratch tree come from the setup head. Is that the intended base, given that gate.md:49-50 defines base as "the base of the change"?
  5. The e2e test needs a WORK through S12: does it use setup-runner, or a WORK made by hand?

### `verify-acceptance` — End to end: `layup setup verify` on a WORK made by `layup setup`

- **What:** A user-path test builds the binary and runs `layup setup WORK` on the stand-in baseline and a stand-in problem statement, with all answers and input files given. Then it runs `layup setup verify WORK` and checks that every row is pass or clear with exit 0. It breaks one value at a time and checks the failing row. It also checks that S15 refuses a verify.tsv with a fail row. A human then does the UAT and the NFR-003 audit.
- **Specification:** docs/spec/setup.md — Where the records go in phase 1 (items 2 and 3: O-verify; S15 refuses a non-pass row); docs/spec/setup.md — The checks of `layup setup verify`; docs/spec/setup.md — REQ-002 — Set up a target, and prove the setup; docs/spec/setup.md — NFR-003 — No value without evidence (item 4: the audit is a review); docs/prd/PRD-0001-layup.md — 7.1 Acceptance criteria (REQ-002, NFR-003, NFR-004)
- **Packages:** cmd/layup, internal/cli, internal/verify
- **Present state:** none. cmd/layup/main_e2e_test.go:11-24 tests only `layup version`.
- **Tests:**
  - e2e: A correct setup gives every row pass or clear and exit 0. The table holds only the 13 checks and the gate:<kind> rows. The table, redirected to WORK/out/verify.tsv, lets the next `layup setup WORK` do S15.
  - e2e: Each seeded defect gives fail on its own row and exit 1: a value row with no source, an unlisted marker, a pin tree that differs, and a removed gate job.
  - e2e: A target with a missing baseline script gives not-active and exit 1, and S15 then refuses that verify.tsv with exit 1.
  - uat: The Operator runs `layup setup` and `layup setup verify` on a pilot problem statement and reads the table. Each row is understandable and each fail names its file (the REQ-002 criterion, PRD-0001 §7.1).
  - uat: The NFR-003 audit (PRD-0001:134; setup.md:329-330): a person reads each value row of record.tsv against its source and records whether the source supports the value. It is a recorded review, not a check.
- **External inputs:**
  - A pilot problem statement and its answers for the UAT (PRD-0001 §7.1 REQ-002).
  - A reviewer for the NFR-003 audit.
- **Open questions:**
  1. Can the stand-in baseline pass `adapted` after a phase-1 setup? See verify-adapted. If it cannot, a correct e2e run cannot exit 0.

### Constraints of area `verify`

1. **docs/spec/setup.md:231-233** — No flag. It reads `WORK/target` at the head of `layup-setup`, `WORK/out/record.tsv` and `WORK/inputs/answers.tsv`, in a scratch work tree. It changes no ref of `WORK/target` and no file of `WORK/out` or `WORK/inputs`; the commit of a fixture run is an object on no ref. *Effect:* Every verify item works in a scratch work tree. Each integration test asserts that refs and input files are unchanged.
2. **docs/spec/setup.md:235-239** — The engine's version of each such check passes and fails on the same fixtures as `docs/setup/tests/run.sh` (the defence that ADR-0011's Consequences name). *Effect:* Each reused check has a root-only core that the harness calls on the sh fixtures. The target additions run only in the frame. verify-harness lands with verify-kit-history, before the other checks.
3. **docs/adr/0011-structure-the-core-engine-as-a-go-cli-over-repository-files.md:123-126** — A check that exists in both can drift apart ... the engine's setup verification must run the same fixtures as `setup/tests/run.sh`. *Effect:* Embedded MK_EXEMPT, AD_EXCLUDE, ad_allowed and identity phrases need drift guards that read setup-check.sh at test time. These are integration tests.
4. **docs/prd/PRD-0001-layup.md:192** — Two checks that drift. ... the engine's setup verification runs the same fixtures as `docs/setup/tests/run.sh` *Effect:* The traceability matrix maps this risk to verify-harness.
5. **docs/setup/tests/run.sh:35-36, :96** — A run that found no case fails: a harness that tested nothing is not a pass. *Effect:* The Go harness cannot land alone. It lands with its first check (verify-kit-history).
6. **docs/setup/tests/run.sh:44-45** — g() { git -c user.name=fixture -c user.email=fixture@invalid -c core.hooksPath=/dev/null -c commit.gpgsign=false "$@"; } *Effect:* Each commit of the harness and of a gate fixture run uses a fixed identity, hooks off and signing off.
7. **docs/spec/gate.md:135-146** — A result is `pass` only when the check ran and passed. A check that did not run is `not-active` ... A missing or malformed manifest is an input error (exit 2) *Effect:* verify-baseline-scripts and verify-gate-fixtures each need a fixture where the check cannot run and gives not-active with exit 1.
8. **docs/spec/setup.md:256** — The first rule that matches decides: the clean run `fail` gives `fail`; either run `not-active` gives `not-active`; a fixture run that is not `fail` gives `fail`, reason `fixture not detected`. A `pending` kind: `clear`, reason `pending: fixture not run` *Effect:* The gate row logic is a fixed decision table with a unit test for each line.
9. **docs/spec/gate.md:59-62** — Write into it, from the base, the manifest and each path of each row's `config`; a path that the base does not have is removed from the scratch tree. So the head's own gate files never judge the head (FT4). *Effect:* A catalog fixture that changes only a config path is never detected. Each fixture must change a product path, and a test proves the rule.
10. **docs/spec/setup.md:257-258** — `ci`, `procedure`, `protection` | not a check for a target ... the rulesets read back (S13); the probes of §3 | no | `layup run`, phase 2: not in the list, so a correct setup exits 0 *Effect:* There are no items or rows for these. The harness skips the 7 ci/procedure/protection cases from a written list, and an e2e test asserts the table has no such row.
11. **docs/spec/setup.md:260-266** — LAYUP's functions count LAYUP's own facts (39 facts of `F-0001`, 25 glossary rows, 9 invariants) ... these checks hold the rules that do not depend on LAYUP's numbers. *Effect:* facts, onboarding, glossary and guardrails are new target-form code. They need their own fixtures in addition to the generic sh lines.
12. **docs/spec/setup.md:190-194** — open-gaps target:docs/setup/open-gaps.tsv no-header ... question text - the question that its value needs; never `—` *Effect:* verify-markers reads open-gaps.tsv by its schema and fails a `—` question. The block is in the schema test.
13. **docs/spec/setup.md:97** — | the gate files | checks `jobs` and `gates` | *Effect:* The one-check call of verify-frame must accept the name `gates`, or the spec must change. It is an open question.
14. **docs/spec/packages.md:18-21, :39, :42, :44** — Only the package `internal/git` starts the `git` program ... Only the packages that the table below marks "starts a program" import `os/exec`. *Effect:* Only internal/git, internal/gate and internal/verify import os/exec, and only internal/git starts git. found-git must add ls-files and rev-list (packages.md:39 does not list them). Only verify-sources, verify-jobs and verify-gate-fixtures need setup-catalog, gate-manifest or gate-run.
15. **docs/spec/packages.md:53-56** — `internal/setup` writes a tree and must not depend on the checks that judge it, so `internal/cli` runs the check of each step from `internal/verify` after `internal/setup` did the step. *Effect:* The step evidence of `layup setup` (S04 to S14) depends on the verify checks. So the verify checks come before, or in the same item as, the setup steps that cite them. S06 depends on verify-facts, which fixes the answers-record form.
16. **docs/spec/packages.md:22-23** — No package of phase 1 imports `net`, `net/http` or `crypto/tls` *Effect:* The job reader and all other code use the standard library only. Tests use no network, so the stand-in baseline is local (verify-test-baseline).
17. **docs/spec/README.md:61-63** — Two runs on the same input print the same bytes (`NFR-005`). So a result table holds no time, no duration and no path of a scratch directory, and its rows have a fixed order. *Effect:* verify-frame fixes the row order. The reasons of the baseline-script rows are `exit <code>`, never script output. An e2e test runs twice and compares the bytes.
18. **docs/spec/setup.md:212-213** — it refuses a `verify.tsv` with a row that is not `pass` or `clear` (exit 1) *Effect:* The setup area's S15 reads the setup-verify schema. Both areas share one schema through internal/tsv.
19. **docs/tests/test-levels.md:43-48, :60, :74; internal/psb/check_integration_test.go:1** — Unit tests touch no file, network, or clock ... every interface or workflow has an integration test ... every user-facing scenario has an E2E test *Effect:* Schema tests and drift guards that read a repository file are integration tests. The one-check call is an interface and has an integration test.
20. **.github/workflows/ci.yml:309 (job tests); :263 (only setup-check sets fetch-depth: 0)** — - uses: actions/checkout@v4 (no fetch-depth in job tests) *Effect:* No Go test may need LAYUP's root commit d2516fd. The stand-in baseline is test data or built at run time, unless the Operator changes ci.yml (a gate change).
21. **.github/workflows/ci.yml:270-288** — rm -rf docs/setup/tests; cp -R .checks-from-default/docs/setup/tests docs/setup/tests *Effect:* The sh harness runs on the fixtures of the default branch, but the Go harness runs on those of the pull request. A fixture change can break the sh run only after the merge. Record this as a known limit.
22. **AGENTS.md (quality gate step 4); docs/spec/gate.md:91-92** — Anything that can run over ten seconds shows which step runs and that it lives. *Effect:* verify-baseline-scripts and verify-gate-fixtures print progress to standard error, never into the table.
23. **docs/setup/setup-check.sh:23-25** — A later setup task adds its check here, with fixtures under docs/setup/tests/<name>/, in the same change. *Effect:* After the Go harness exists, a change to a sh check must also keep the Go harness green. The plan says this in each later setup task.

### Risks of area `verify`

- Check `adapted` may fail every correct phase-1 setup. The unchanged baseline (LAYUP's root commit d2516fd) holds rule-1/rule-3 words in kept files that no phase-1 step rewrites. An Operator decision is needed before verify-acceptance can pass.
- S12's evidence names a check `gates` (setup.md:97) that the verify table does not define. If the plan does not fix the name, the setup area and the verify area will implement different evidence for S12.
- The answers-record form is not fixed (setup.md:91). Check `facts` and S06 both need it. If each area guesses its own form, S06 writes a record that check `facts` cannot read.
- The same fixtures cannot prove the target-form checks. 12 of the 40 cases are LAYUP-form. So the drift defence covers fewer checks than ADR-0011 suggests, and the target-form logic has only new fixtures to guard it.
- Output-form mismatch. The sh check prints every FAIL line, and the Go table prints only the first failure (setup.md:271). If the Go cores do not return all findings in the sh line form, the harness can compare only exit codes.
- The sh harness and the Go harness run on different fixture trees in CI (ci.yml:270-288). A fixture change can pass on the pull request and fail on main after the merge.
- Test data under internal/verify/testdata/ that holds markers or kit words fails LAYUP's own markers and adapted checks, because MK_EXEMPT and AD_EXCLUDE do not exempt it. Changing setup-check.sh is a gate change.
- Known sh defects (#21: quoted file names and a blank question; #61: a tab-only fact; #48: autocrlf) make a correct Go port disagree with the sh check on some inputs. Each needs a decision: fix both, or copy the defect.
- The baseline is not fixed. A target pins the latest baseline commit (O-101), so a newer baseline can break kit-history, adapted, discipline-tests or link-lint. Tests on a stand-in can hide such a break.
- Cross-area order: verify-gate-fixtures needs gate-run, gate-manifest and setup-catalog. The setup area's step evidence needs the verify checks. A wrong order in the plan blocks both areas.
- A catalog fixture that changes only a config path is never detected, because the gate writes the config from the base (gate.md:59-62). A wrong fixture gives `fixture not detected` on every correct setup.
- The baseline's run-discipline-tests.sh asserts line endings. A scratch checkout on an autocrlf host can fail it for a reason that is not a setup defect.

### Critic points that the reviser of `verify` rejected

- Correction to verify-facts size: "medium, not small". *Basis:* The output schema allows only `small` or `large`. The item is set to large. The rest of the correction (the item holds more than about 400 lines) is accepted.
- Unsupported: verify-frame lists NFR-007 with no import test; the critic suggests a `go list -deps` test in this item. *Basis:* NFR-007 is removed from verify-frame and no import test is added here. docs/spec/packages.md:14-25 states the import rules for every package, so the module-wide test belongs to the foundation area, not to one verify item. NFR-001 is also removed, as the critic says (setup.md:73; records.md:99-104).

## Area `found`

At 7cdd346 the code is small. cmd/layup/main.go hands the arguments to internal/cli. internal/cli/cli.go:26-59 dispatches `version` and `psb check` and returns 0, 1 or 2. internal/psb/check.go:103-108 writes its own table with fmt. internal/tsv and internal/git do not exist. The code has no schema-block test, no import-boundary test, no telemetry, stall or price type. The e2e level has one test, and that test builds the binary in the test itself (cmd/layup/main_e2e_test.go:12-24). `go list -f {{.Imports}} ./...` names only stdlib and the 3 module packages. CI already runs the integration and e2e tags (ci.yml:313-318). The foundations are 10 items: found-tsv (large), found-spec-schema-test, found-git (large), found-boundary-test, found-cli, found-telemetry-schema, found-stalls-schema, found-e2e-harness, found-nfr001-tests, found-nfr002-tests. Changes after the critic: found-spec-schema-test now depends only on found-tsv, and every schema owner registers its schema with it in its own change. found-cli now serves NFR-004 and not NFR-001, and gives only the frame (version, psb check, usage, exit mapping, argument rules); each command item adds its own dispatch. found-git now makes each call independent of the host's git config and terminal. found-boundary-test now reads the direct non-test imports and scans the source for program starts. The telemetry completeness test and the Stall Diagnosis helper are removed, because both are measures of `layup report` (phase 4). The NFR-001 and NFR-002 tests now depend on a full setup run and on verify. Their pilot UATs move to a later phase. The main open decisions are these. Which package holds the telemetry, price and stall types in phase 1 (packages.md has no row)? Do the import rules apply to _test.go files? What identity signs the setup commits? How does a typed column accept `—`? How does the isolation of git config agree with the credentials that a private baseline needs?

| Key | Title | Requirements | Size | After | Host |
| --- | ----- | ------------ | ---- | ----- | ---- |
| `found-tsv` | internal/tsv: read and write a record by its schema | NFR-005, NFR-007, NFR-001, NFR-002 | large | — | row 1 |
| `found-spec-schema-test` | The test that compares each tsv-schema block of docs/spec with its Go schema | REQ-001, REQ-002, REQ-004, REQ-007, REQ-009, REQ-011, NFR-003 | small | found-tsv | row 1 |
| `found-git` | internal/git: the one caller of the git program | NFR-007, NFR-005, NFR-001 | large | — | row 2 |
| `found-boundary-test` | The import-boundary and program-start test of the package table | NFR-007, NFR-005 | small | — | row 2 |
| `found-cli` | internal/cli: the frame of every command (usage, argument rules, exit-code mapping) | NFR-004, REQ-001, REQ-002, REQ-004 | small | found-e2e-harness | row 3 |
| `found-telemetry-schema` | The telemetry and price schemas in code (REQ-011, schema only) | REQ-011 | small | found-tsv, found-spec-schema-test | row 17 |
| `found-stalls-schema` | The stall schema in code (REQ-009, schema only) | REQ-009 | small | found-tsv, found-spec-schema-test | row 18 |
| `found-e2e-harness` | The e2e harness: build the binary once and run commands as a user | NFR-005, REQ-001 | small | — | row 3 |
| `found-nfr001-tests` | The tests of NFR-001 in phase 1: the records are in the target's Git | NFR-001 | small | found-e2e-harness, found-tsv, setup-s01-s14 (setup area: every step of a full run), setup-s13-ruleset (setup area: the order of commands.sh), setup-s15-records-branch (setup area), setup-baseline-fixture (setup area: a fixture baseline that passes its own discipline-tests and link-lint), verify-checks (verify area: every check, and the setup-verify schema) | row 16 |
| `found-nfr002-tests` | The tests of NFR-002 in phase 1: a target is independent of LAYUP | NFR-002 | small | found-e2e-harness, setup-s01-s14 (setup area: a full run), setup-s12-gates (setup area: catalog files, docs/gates.tsv, one CI job per kind), setup-s13-ruleset (setup area), setup-s15-records-branch (setup area), setup-catalog-go (setup area: the Go entry), setup-baseline-fixture (setup area), verify-checks (verify area) | row 16 |

### `found-tsv` — internal/tsv: read and write a record by its schema

- **What:** A new package, internal/tsv, reads and writes a record against a Schema value: the header row, the field count, and the type of each field from the closed list of 11 types. The writer applies the field rules: one space for each tab, LF and CR; `—` for an empty field; LF endings and a final LF. It also reads and writes the no-header form, and it parses one tsv-schema block into a Schema. psb, gate, setup, verify and catalog then write all their tables through it.
- **Specification:** docs/spec/README.md: Records; docs/spec/README.md: Records > The schema block; docs/spec/README.md: Records > The types; docs/spec/README.md: Commands (Output, Determinism); docs/spec/packages.md: The table of phase 1 (row internal/tsv); docs/architecture.md: 3. Records and identities
- **Packages:** internal/tsv
- **Present state:** None. internal/psb/check.go:103-108 writes its table with fmt.Fprintf. check.go:111 replaces tabs in the excerpt but not a lone CR (psb-check.md:22-24 says #29 fixes this). The move of psb onto tsv is psb-area work (packages.md:46-47).
- **Tests:**
  - unit: Each type of the closed list accepts its valid values and rejects invalid ones: int with a sign or a leading zero, decimal with an exponent, time with no Z or with fractions, uppercase sha1, a path with `..` or a leading `/`, an id that does not match its N/x/<word> pattern, an enum word not in the list, a list with two spaces.
  - unit: The writer puts one space for each tab, LF and CR in a field, writes `—` for an empty field, joins the header names with one tab, and ends each row with one LF, the last row included.
  - unit: The reader rejects a wrong header, a row with too few or too many fields, a field of the wrong type, a CR, invalid UTF-8, a byte-order mark, an empty field (an empty string, not `—`), a missing final LF and an empty line. Each error names the row and the column, so the caller can return exit 2.
  - unit: The no-header form reads rows by the column order of the schema and does not take a header row as data.
  - unit: Write then read gives the same rows, and two writes of the same rows give the same bytes (NFR-005).
  - unit: The block parser reads the README form from an in-memory string: name, location, an optional no-header word; one line per column split on one or more spaces; the rest of the line is the rule. It rejects a tab in a block, an unknown type, a key word that is not `key` or `-`, and a duplicate column name.
- **Open questions:**
  1. How does a typed column accept `—`? README.md:70-71: 'An empty field is written as `—`'. The types (README.md:106-118) have no 'may be empty' mark; the permission is only in rule words (records.md:159-171 '`—` when not reported'; setup.md:193 'never `—`'). Reading A: `—` passes in every column. Reading B: a machine-readable mark, which is a README change first ('a new type is added here first', README.md:104).
  2. Is N in id(<pattern>) a fixed width or a minimum? README.md:117 gives `id(Q-NNN)`, but psb writes `Q-%03d` (check.go:106), which gives Q-1000 for the 1000th gap. stalls 'from ST-001, in order' (records.md:204) has the same question.
  3. Is a decimal with no point (for example `3`) valid? README.md:110: 'a decimal number, 0 or more, with a point and no exponent'.
  4. Does the reader reject a CR before LF, or strip it? README.md:69-70 gives only the writer's rule. psb-check.md:21-22 ignores a CR before LF in a Markdown input. answers.tsv is written by the Operator (setup.md:66), possibly on Windows.
  5. What does the writer do with invalid UTF-8 in a field? README.md:67 says a record is UTF-8, and psb copies excerpts from any input file (psb-check.md:21). No source says whether the writer replaces the bytes or returns an error.
  6. Does tsv check that key columns are unique? packages.md:38 lists header, field count and types only. README.md:100 defines `key` but no check.
  7. list(text) separates values with one space (README.md:118), so a value with a space (for example a scope pattern, gate.md:26) cannot be written. Is that an error?

### `found-spec-schema-test` — The test that compares each tsv-schema block of docs/spec with its Go schema

- **What:** A Go test reads every tsv-schema block in docs/spec/*.md and compares it with the registered Go schema of the same name: location, no-header mark, column names and order, types and key marks. It fails on two blocks with the same name, on a location with no valid prefix, and on a block or a Go schema with no partner. A list of block names that have no code yet starts with the 14 blocks less those that land with this item; it may only shrink, and it is empty at the end of phase 1. Each schema owner registers, in its own change, the same Schema value that its writer and reader use.
- **Specification:** docs/spec/README.md: Records > The schema block; docs/spec/psb-check.md: The table; docs/spec/gate.md: The gate manifest; docs/spec/gate.md: The table; docs/spec/setup.md: The stop table; docs/spec/setup.md: The step table; docs/spec/setup.md: The answers; docs/spec/setup.md: The setup record; docs/spec/setup.md: The rule-path register; docs/spec/setup.md: The files of `docs/setup/` in a target; docs/spec/setup.md: The checks of `layup setup verify`; docs/spec/setup.md: The stack catalog; docs/spec/records.md: REQ-011 — The telemetry record; docs/spec/records.md: REQ-009 — The stall record
- **Packages:** internal/tsv (the comparer and the block scanner), the home of the registry is open (see open_questions)
- **Present state:** None. The 14 blocks: psb-check.md:48; gate.md:21, :82; setup.md:110, :125, :142, :155, :168, :190, :268, :288; records.md:150, :174, :203. The README form example at README.md:80 is inside a four-backtick fence. No block has Go code; psb-gaps is written by hand (check.go:104).
- **Tests:**
  - unit: The comparer, given an in-memory fixture block and a fixture schema, reports a difference for each changed attribute: name, location, no-header, column order, type, key.
  - unit: The block scanner follows fence nesting: a tsv-schema fence inside a four-backtick text fence (the form example at README.md:78-84) is not a block. It reports a duplicate name and a location whose prefix is not stdout, records:, target:, layup: or host: (README.md:86-93).
  - integration: Each block of docs/spec at the current commit either matches its registered Go schema or is on the not-yet-built list, and each registered schema has a block. Renaming a column on either side makes the test fail.
  - integration: Each owner's own test reads that owner's real output (for example the stdout of a command) with its registered schema, so a declared schema cannot pass while the command writes another form. This is a rule for each owner item, checked in this item for the first owner.
- **Open questions:**
  1. Where does the registry live? A central test must import every owner (psb, gate, setup, verify, catalog, the home of telemetry and stalls), but cli may import only psb, setup, verify and gate (packages.md:37). Option A: a test in each owner package plus one completeness test over the block names. Option B: one test file in internal/cli, if test files are exempt from the import table (see found-boundary-test).
  2. Is a shrinking list of not-yet-built blocks acceptable? README.md:76-77: 'the test comes with the code'. Without the list the test is red until every owner lands.

### `found-git` — internal/git: the one caller of the git program

- **What:** A new package, internal/git, is the only package that starts `git`, with os/exec. It gives typed calls for the phase-1 operations and returns a typed error that tells 'git not found' from 'git failed' (with stderr). Each call is independent of the host: no system or global config changes a branch name, a byte of a tree, a hook or a signing step, and no call can ask a question on the terminal. init makes the branch `main`.
- **Specification:** docs/spec/packages.md: NFR-007 — Go, the standard library only, and Git as the `git` program (rules 3 and 4); docs/spec/packages.md: The table of phase 1 (row internal/git); docs/spec/README.md: Commands (Arguments, Determinism); docs/spec/setup.md: The command `layup setup` (the work area: root commit on `main`); docs/spec/setup.md: The steps (S02, S03, S04, S15); docs/spec/setup.md: Where the records go in phase 1; docs/spec/setup.md: The checks of `layup setup verify` (gate:<kind>: git apply; 'an object on no ref'); docs/spec/gate.md: The command (The run, steps 1, 2, 4; the diff row); docs/architecture.md: 1. LAYUP and a target; docs/architecture.md: 3. Records and identities
- **Packages:** internal/git
- **Present state:** None. No package imports os/exec at 7cdd346 (`go list -f {{.Imports}} ./...`); only the test cmd/layup/main_e2e_test.go:6 does.
- **Tests:**
  - unit: Each call builds the right argument list (with the fixed -c values and environment) and parses the output (ls-remote line, rev-parse SHA, diff --name-only list), with a stub runner and no file or process.
  - integration: With real git in temp directories, HOME isolated: ls-remote and clone of a local file:// bare repository; init on `main`, add, commit; rev-parse <commit>^{tree}; worktree add --detach and remove leave no worktree entry; show <rev>:<path> of a present and a missing path; diff --name-only; apply of a patch; an orphan-branch commit with no shared history; a commit object on no ref.
  - integration: With a hostile global config (core.autocrlf=true, init.defaultBranch=master, core.hooksPath to a hook that fails, commit.gpgsign=true), the root commit is on main, its tree equals the tree of the source files, and no hook runs.
  - integration: No call can ask a question: GIT_TERMINAL_PROMPT=0, stdin not connected, ssh in batch mode; a clone of a URL that needs a password fails at once with the 'failed' kind and does not hang.
  - integration: The error kinds with real processes: a PATH with no git gives 'not found' (exec.ErrNotFound); a git that exits non-zero gives 'failed' with stderr kept (*exec.ExitError).
- **External inputs:**
  - The git version on the LAYUP host and on ubuntu-latest, with the git documentation of init --initial-branch, worktree add --detach, ls-remote, orphan branches and commit-tree at that version.
  - A decision of the Operator on the author and committer identity of the setup commits (S03, the 'chore: setup <step>' commits, the S15 records commit). No source gives it.
  - A decision of the Operator on how S02 reaches a private baseline repository: which credential helper or ssh key the clone may use, when the global config is not read.
- **Open questions:**
  1. Is the operation list closed? packages.md:39 lists 'clone, ls-remote, init, commit, rev-parse, worktree, show, diff --name-only, apply'. But setup.md:87 needs 'git checkout', S03 needs add, S04 needs a branch 'from the root commit' (setup.md:89), S15 needs an orphan branch (setup.md:213-217), and setup.md:232 needs a commit that is 'an object on no ref'. The table row changes in the same PR.
  2. Which identity signs the commits? setup.md:88, 102-103 and 213-217 give messages but no author. A CI runner has no user.name, so `git commit` fails unless the code sets one. The commit time also makes SHAs differ between runs.
  3. How does the config isolation (GIT_CONFIG_NOSYSTEM, an empty GIT_CONFIG_GLOBAL, or explicit -c values) agree with S02's network clone, which may need the Operator's credential helper from the global config? setup.md:87 gives the commands only.
  4. What does a command do when git is missing? NFR-004 (gate.md:140-141) makes a missing tool `not-active` for a gate row, but setup and gate need git themselves. Exit 2 (README.md:55) or a `not-active` row?
  5. Is there a minimum git version? No source gives one.

### `found-boundary-test` — The import-boundary and program-start test of the package table

- **What:** A Go test enforces the rules of packages.md. It reads `go list -json ./...` (the direct, non-test .Imports of each package) for rules 1, 2, 4 and the may-import column, `go list -deps ./...` for rule 2 and rule 5, and go.mod for 'no require line'. A go/ast scan of the non-test source enforces rule 3 and the program column: only internal/git starts `git`, internal/gate starts only `sh -c`, internal/verify starts only `sh`, and no other package calls exec.Command, os.StartProcess or a syscall exec function. Each package must have a table row. It is the mechanical check of NFR-007 and NFR-005, and it fills the Test cell of NFR-007 in PRD §12.
- **Specification:** docs/spec/packages.md: NFR-007 — Go, the standard library only, and Git as the `git` program (rules 1 to 5); docs/spec/packages.md: The table of phase 1; docs/spec/gate.md: NFR-005 — No model call in the engine checks (point 1); docs/architecture.md: 6. Native gates and rule protection > The stack gates (boundary from the package table)
- **Packages:** none new (a _test.go file; the home is an open question)
- **Present state:** None. go.mod:1-3 has no require line. Measured: cmd/layup -> internal/cli, os; internal/cli -> internal/psb, fmt, io, os; internal/psb -> stdlib only. Measured: `go list -deps -test ./internal/psb` names os/exec through testing/internal/testdeps -> internal/fuzz, so a test-including scan breaks rule 4 for every package. `go list -deps os/exec embed path/filepath` names no net package. PRD-0001:241 gives NFR-007 the Test 'CI job lint (gofmt, go vet)', which does not prove the criterion.
- **Tests:**
  - unit: The rule checker, given an in-memory import graph, reports each broken rule: a non-stdlib dependency, a require line, os/exec in a 'no' package, net, net/http or crypto/tls in any package, a module import the table does not allow, a package with no row.
  - unit: The program-start checker, given in-memory Go source, reports exec.Command("git") outside internal/git, a program other than sh in internal/gate or internal/verify, and os.StartProcess or syscall.Exec in any package.
  - integration: The real `go list` output and the real non-test source of this module pass every rule at the current commit. A package that the table names but that does not exist yet is not an error.
- **Open questions:**
  1. Do the rules bind _test.go files? cmd/layup is 'Starts a program: no' (packages.md:36), but main_e2e_test.go:6 imports os/exec, and this test needs os/exec for `go list`. The e2e harness needs git, and found-nfr001-tests needs internal/tsv from a cmd/layup test file (packages.md:36 allows only internal/cli). One decision answers all of them.
  2. Does the test parse the Markdown table (packages.md:34-44; 'the boundary rule of the future Go gate of LAYUP reads this table', packages.md:31-32), or hold a Go copy? A copy can drift; a parser breaks on a prose edit.
  3. Does rule 5 mean direct imports or -deps? packages.md:22 says 'No package of phase 1 imports net'. Today both readings pass.
  4. Where does the test live? A new package needs a table row (packages.md:31), and root tests/ holds fixtures only (engineering-discipline.md:733-734).

### `found-cli` — internal/cli: the frame of every command (usage, argument rules, exit-code mapping)

- **What:** internal/cli keeps `version` and `psb check` and gives the shared frame that each command item plugs into: one usage text, the argument rules (no environment variable and no terminal question for an input), and one mapping from result words to exit codes. 0: every row is pass, clear, done or operator, or no gap. 1: any row is fail or not-active, or a gap. 2: a usage or input error, with empty stdout. 3: only a setup stop. A check that did not run is never 0. Diagnostics and progress lines go only to stderr. The dispatch of `setup`, `setup verify` and `gate` lands with each command's item; the setup step loop in internal/cli (setup, then its verify check, then the done row, packages.md:54-56; the hand-over of the psb gap table, setup.md:86) is setup-area work (best key: setup-step-loop).
- **Specification:** docs/spec/README.md: Commands; docs/spec/packages.md: The table of phase 1 (rows cmd/layup and internal/cli); docs/spec/gate.md: The command (exit codes); docs/spec/gate.md: NFR-004 — A check that is not active is not passed (points 1 and 4); docs/spec/setup.md: The command `layup setup` (exit codes); docs/spec/setup.md: The step table (result enum); docs/spec/setup.md: The checks of `layup setup verify` (exit codes)
- **Packages:** internal/cli, cmd/layup
- **Present state:** internal/cli/cli.go:16-22 has a usage for version and psb check only; cli.go:26-41 dispatches the two; cli.go:43-59 maps psb to 0/1/2. cmd/layup/main.go:11 calls os.Exit(cli.Run(...)). Tests: internal/cli/cli_test.go:15-70, cmd/layup/main_e2e_test.go:12-24. No os.Getenv or os.Stdin in the code today. README.md:12-13 and docs/onboarding-for-engineers.md:84 name only `version` and `psb check`; each command item updates them in its PR.
- **Tests:**
  - unit: The mapping gives 0 for rows that are all pass or clear (gate, verify) and for rows that are all done or operator (setup steps, setup.md:53-54, :128); 1 for any fail or not-active, never 0 for a check that did not run; 2 for an input error; 3 only for a setup stop.
  - unit: An unknown command, an unknown flag, a missing argument or an extra argument gives exit 2, the usage on stderr and empty stdout. `version` prints 'layup 0.1.0-dev' and exits 0 (cli_test.go:15-26 stays).
  - unit: The argument parser accepts positional arguments before flags (`layup gate REPO --base REV --head REV`, gate.md:44); Go's flag package alone stops at the first non-flag.
  - integration: A go/parser scan of the non-test Go files of internal/ and cmd/ finds no os.Getenv, os.LookupEnv or read of os.Stdin (README.md:44-45). It reads files, so it is integration.
  - e2e: The built binary prints the usage on stderr, prints nothing on stdout and exits 2 with no arguments (through found-e2e-harness).
- **Open questions:**
  1. How is `layup setup verify` with no WORK parsed: as verify with a missing WORK (exit 2), or as setup of a directory named `verify`? setup.md:46 'layup setup WORK' and setup.md:228 'layup setup verify WORK' give no rule.
  2. internal/cli may import only psb, setup, verify and gate (packages.md:37), and setup may not import psb (packages.md:43), but setup.md:86 says cli 'hands the gap table to internal/setup'. Which type carries it: a tsv record, or a type that setup defines?
  3. May internal/git read the environment to pass it to the child git (os.Environ)? README.md:44-45 forbids only reading an environment variable 'for an input'.
  4. Does the progress rule (engineering-discipline.md:832-834, 'Any task, batch job, network operation, or loop that takes more than 10 seconds must show an active, informative progress indicator') bind the product's commands? If yes, setup (the S02 clone) and setup verify need stderr progress lines; gate.md:91-92 gives a stderr block per kind for gate only.

### `found-telemetry-schema` — The telemetry and price schemas in code (REQ-011, schema only)

- **What:** A Go row type and a registered schema for records:telemetry.tsv and host:prices.tsv, with a validator of the types and of the row rules that the blocks state. Examples: latency_s is `—` exactly when first_output is `—`; tokens_reason is `—` exactly when tokens_status is observed; money and currency are `—` when money_status is unknown, and money is never 0 then; price names a P-NNN row only for computed. No command writes or reads a row in phase 1, and the completeness measure is not built.
- **Specification:** docs/spec/records.md: REQ-011 — The telemetry record; docs/spec/records.md: The layout of the records branch (rows telemetry.tsv and host:prices.tsv); docs/architecture.md: 12. Cost, budget and the measures > The ledger
- **Packages:** the home is open: internal/ledger (phase 2, packages.md:85), internal/records (phase 2, packages.md:76), or a new phase-1 row
- **Present state:** None. The blocks are at docs/spec/records.md:150-172 and :174-183.
- **Tests:**
  - unit: A complete observed row and a row of each status (partial, unavailable; reported, computed, unknown) pass. Each broken row rule gives an error that names the column: a latency with no first_output, money 0 with unknown, a computed row with no price id, tokens_reason `—` with partial.
  - unit: A prices row with a bad id, a class outside in|out|cache, or a price with an exponent is rejected.
  - integration: The registered schemas equal the telemetry and prices blocks (records.md:150, :174), through found-spec-schema-test; both names leave the not-yet-built list.
- **External inputs:**
  - A decision of the Operator on the phase-1 home of the type, and the packages.md row with its import rules in the same change (packages.md:70-72).
- **Open questions:**
  1. Which package holds the type in phase 1? packages.md:70-72: a package 'gets its row ... in the same change that adds its requirement's section'. REQ-011's section is in phase 1 (records.md:138) but no row was added, and packages.md:85 puts internal/ledger in phase 2.
  2. Must the validator enforce the row rules that the block gives only in words (records.md:159-171), or only the types?
  3. Must the validator recompute latency_s and duration_s from the times (records.md:161-162 'first_output − start')?

### `found-stalls-schema` — The stall schema in code (REQ-009, schema only)

- **What:** A Go row type, a registered schema and a validator for records:stalls.tsv. The validator checks the types, the per-kind rules (trigger, evidence, cause, rung, examiner and outcome are `—` except on the rows the block names) and the order: per stall ID a `stall` row, then one `diagnosis` or `diagnosis-failed` row, then an `outcome` row; IDs from ST-001 in order. No command writes or reads a row in phase 1, and the Stall Diagnosis count is not built.
- **Specification:** docs/spec/records.md: REQ-009 — The stall record; docs/spec/records.md: The layout of the records branch (row stalls.tsv); docs/architecture.md: 11. Stalls
- **Packages:** the home is open: internal/stall (phase 3, packages.md:88) or internal/records (phase 2, packages.md:76)
- **Present state:** None. The block is at docs/spec/records.md:203-215.
- **Tests:**
  - unit: A file of complete stalls passes. Each of these gives an error: a diagnosis before its stall row; both diagnosis and diagnosis-failed for one stall; an outcome before a diagnosis; a cause on a non-diagnosis row; ST-002 before ST-001.
  - integration: The registered schema equals the stalls block (records.md:203), through found-spec-schema-test.
- **External inputs:**
  - A decision of the Operator on the phase-1 home of the type, and its packages.md row (packages.md:70-72).
- **Open questions:**
  1. Is a stall with only its `stall` row valid? records.md:217-218 counts stalls whose second row is 'missing', but records.md:199-200 says 'Each stall has three rows'.
  2. Key (stall, kind) at records.md:204-205 allows both diagnosis and diagnosis-failed for one stall. Must the validator forbid it ('one row of each of its three steps')?
  3. Which package holds the type in phase 1 (packages.md:88 puts internal/stall in phase 3)?

### `found-e2e-harness` — The e2e harness: build the binary once and run commands as a user

- **What:** A shared e2e helper in cmd/layup (tag e2e). TestMain builds the layup binary once. A run helper returns stdout, stderr and the exit code, with HOME and the git config isolated and PATH controlled. Helpers make a git repository from a fixture directory under root tests/, and a repeat helper runs a command twice and compares the bytes (NFR-005). Fixture Go code under tests/ has its own go.mod, or is a patch or a file with no .go suffix, so that it is not a package of LAYUP's module and gofmt does not read a bad-on-purpose file. tests/README.md is updated in the same PR.
- **Specification:** docs/spec/README.md: Commands (Exit codes, Determinism); docs/spec/gate.md: REQ-007 — Each change passes the gates before review or merge (the repeat rule); docs/spec/gate.md: NFR-005 — No model call in the engine checks (point 2); docs/spec/packages.md: NFR-007 (rule 1, one binary)
- **Packages:** cmd/layup (test files only)
- **Present state:** cmd/layup/main_e2e_test.go:12-24 builds the binary inside the test and checks one output. tests/README.md:6-7 says the directory 'is empty so far'. CI runs `gofmt -l .` over the whole tree (ci.yml:302); pre-commit checks tracked .go files (.githooks/pre-commit:126-130).
- **Tests:**
  - unit: The byte-compare part of the repeat helper fails on a difference of one byte and passes on equal bytes (pure code, no file).
  - e2e: The binary is built once, and the existing scenario (layup version prints 'layup 0.1.0-dev', exit 0) passes through the harness.
  - e2e: layup psb check on a fixture under tests/ gives the golden table, and two runs give the same bytes.
- **External inputs:**
  - A baseline fixture repository for the setup e2e paths: a local file:// repository built from a fixture under tests/. Its content is setup-area work; the Operator decides whether a real baseline clone may run in CI.
- **Open questions:**
  1. The harness makes fixture repositories with git through os/exec in a cmd/layup test file, but cmd/layup may import only internal/cli (packages.md:36), and only internal/git starts git (packages.md:18-19). Are test files exempt (found-boundary-test)?
  2. The PATH control is proven through a real path (a gate row `not-active`, gate area). Which gate item owns that e2e scenario?

### `found-nfr001-tests` — The tests of NFR-001 in phase 1: the records are in the target's Git

- **What:** These e2e tests show that, after a full setup run on the fixture (S01 to S14, `layup setup verify` into WORK/out/verify.tsv, then S15), WORK/target has the orphan branch layup-records. It holds README.md (present), setup/record.tsv, setup/verify.tsv and rule-paths.tsv, each equal to its copy in WORK/out and valid by its schema, and the answers are fact records in the target tree. So no value is held only in the work area. The PRD criterion (an audit of a pilot task, PRD-0001:132) is a pilot measure of a later phase and does not close this item.
- **Specification:** docs/spec/records.md: NFR-001 — Git is the system of record; docs/spec/records.md: The layout of the records branch; docs/spec/setup.md: Where the records go in phase 1; docs/spec/setup.md: The steps (S06, S11, S15); docs/architecture.md: 3. Records and identities
- **Packages:** cmd/layup (test files only)
- **Present state:** None.
- **Tests:**
  - e2e: layup-records exists in WORK/target, `git merge-base main layup-records` finds no common commit, README.md is present, and the three tables are valid by their schemas.
  - e2e: setup/record.tsv, setup/verify.tsv and rule-paths.tsv on layup-records equal WORK/out/record.tsv, verify.tsv and rule-paths.tsv byte for byte.
  - e2e: Each question ID of WORK/inputs/answers.tsv is a fact, with its answer, in a raw fact record of the layup-setup tree (S06, S11).
  - e2e: commands.sh lists, in order, the root push (S03), the layup-setup push (S13), the layup-records push (S15) and the ruleset apply (setup.md:218-223).
- **External inputs:**
  - A pilot target and a person for the later-phase audit (PRD-0001:132).
- **Open questions:**
  1. records.md:105-107 says 'each answer row names the comment that holds it (`source`)', but S06 writes 'the question ID, the question text and the answer' (setup.md:91), and the record's `source` holds a kind and `ref` a question ID (setup.md:159-160). The comment URL of answers.tsv (setup.md:145) seems to stay only on the host. Is that a spec defect?
  2. How does a cmd/layup test file check a file against its schema, when cmd/layup may import only internal/cli (packages.md:36)? Through a layup command, or by the test-file exemption?
  3. README.md of layup-records is 'a fixed text' (setup.md:215-216) with no given text. The test checks presence only.

### `found-nfr002-tests` — The tests of NFR-002 in phase 1: a target is independent of LAYUP

- **What:** These e2e tests show that the target that a full setup run makes (S01 to S15, with verify) holds no layup binary, no LAYUP script and no CI step that runs layup; the ruleset requires no `layup/` check; a plain clone carries origin/layup-records; and, in that clone at layup-setup with no layup on PATH, what each gate job runs passes. The second half of the PRD criterion (a different harness continues a task, PRD-0001:133) is a pilot measure of a later phase.
- **Specification:** docs/spec/records.md: NFR-002 — A target is independent of LAYUP; docs/spec/setup.md: The steps (S12, S13, S15); docs/spec/gate.md: The gate manifest; docs/spec/gate.md: The command (the table of the run: pending and active rows); docs/architecture.md: 1. LAYUP and a target
- **Packages:** cmd/layup (test files only)
- **Present state:** None.
- **Tests:**
  - e2e: No file of the layup-setup head is a layup binary, and no workflow run line names `layup`. WORK/out/ruleset-default.json and docs/setup/branch-protection.json name no `layup/` context (records.md:127-134).
  - e2e: A `git clone` of WORK/target has origin/layup-records (records.md:135-136).
  - e2e: In that clone, with layup-setup checked out explicitly and a PATH with no layup: each active kind's command of docs/gates.tsv exits 0 with sh -c, and each pending kind's job rule (no product path changed) passes.
- **External inputs:**
  - The tool of each Go catalog kind on the CI runner (only Go is installed, ci.yml:310-312).
  - A pilot target and a second harness for the later-phase half (PRD-0001:133).
- **Open questions:**
  1. Which strings count as 'a LAYUP script or CI job'? records.md:127-129 names them but gives no mechanical test. Is a search for `layup` in workflow run lines enough?
  2. How does a test run 'what each job runs' with no YAML library (packages.md:16-17)? From docs/gates.tsv, which the jobs run (records.md:129-130), or from the workflow file?

### Constraints of area `found`

1. **docs/spec/packages.md:16-17** — go.mod has no require line. `go list -deps ./...` names only packages of the standard library and of this module. *Effect:* No test library, no depguard, no YAML or TSV library. The boundary test uses `go list`, go/parser and go/ast only.
2. **docs/spec/packages.md:18-21** — Only internal/git starts git, with os/exec. Only packages marked 'starts a program' import os/exec. *Effect:* internal/git lands before setup and gate. The boundary test needs a source scan for rule 3, because gate and verify may import os/exec. The test-file question needs one decision before the e2e harness grows.
3. **docs/spec/packages.md:22-25** — No package of phase 1 imports net, net/http or crypto/tls. *Effect:* The boundary test is the mechanical check of NFR-005. No phase-1 item adds a network client.
4. **docs/spec/packages.md:34-44** — May import: cli -> psb, setup, verify, gate; setup -> tsv, git, catalog; verify -> tsv, git, catalog, gate; psb, catalog -> tsv; gate -> tsv, git. *Effect:* Order: tsv and git first; then catalog, gate, the psb move; then setup and verify; the cli dispatch of each command lands with that command. Data across packages uses types of the importee or of tsv.
5. **docs/spec/packages.md:46-47** — internal/psb ... moves to internal/tsv when that package exists (#29). *Effect:* The psb move (psb area) depends on found-tsv, and registers psb-gaps with found-spec-schema-test.
6. **docs/spec/packages.md:70-72** — Each later milestone gives a package its row in the table above, with its import rules, in the same change that adds its requirement's section. *Effect:* The home of the telemetry, price and stall types needs an Operator decision and a packages.md row before or with their code.
7. **docs/spec/README.md:75-77** — A Go test of the code reads these blocks and compares them with what the code writes and reads (the test comes with the code, Bootstrap mode rule 1). *Effect:* found-spec-schema-test lands with found-tsv as a comparer and a shrinking not-yet-built list; each owner registers its schema in its own PR.
8. **docs/spec/README.md:48-60** — Exit codes 0, 1, 2 and 3 (3 only for layup setup). A check that did not run is never 0 (NFR-004). *Effect:* found-cli owns one mapping over the result words (pass, clear, done, operator; fail, not-active) that every command uses.
9. **docs/spec/README.md:44-47** — A command reads no environment variable for an input, asks no question on the terminal, prints its result as one table on stdout and diagnostics only on stderr. *Effect:* internal/git runs git with no prompt and no host config; the cli tests assert empty stdout on error; a source scan finds no os.Getenv or os.Stdin.
10. **docs/spec/README.md:61-63** — Two runs on the same input print the same bytes; a result table holds no time, no duration and no scratch path. *Effect:* tsv writes one canonical form; the e2e harness gives a repeat helper; no table prints a setup commit SHA that depends on a commit time.
11. **docs/spec/README.md:104** — The list is closed; a new type is added here first. *Effect:* tsv implements exactly the 11 types. A nullable mark or a new type is a spec change before the code.
12. **docs/spec/records.md:144-146, :189-190, :197-198, :222-223** — No phase-1 command writes a telemetry or a stall row; the writer, the Telemetry Completeness report (phase 4), the triggers and the examiner are not in phase 1. *Effect:* The REQ-009 and REQ-011 items give types, validators and block tests only; no command, no e2e test, no measure.
13. **docs/prd/PRD-0001-layup.md:179** — Phase 1 ships ... the schemas of the telemetry record (REQ-011) and the stall record (REQ-009) ...; NFR-001 to NFR-007 hold from the first release. *Effect:* The schema items and the NFR tests are phase-1 work.
14. **docs/prd/PRD-0001-layup.md:211-213** — the implementation plan (`T-55n2`) fills the Task column, and each delivering task fills the Test column. *Effect:* found-boundary-test (NFR-007, NFR-005), found-nfr001-tests, found-nfr002-tests, found-telemetry-schema (REQ-011) and found-stalls-schema (REQ-009) each fill their Test cell in their own PR. The NFR-007 cell 'CI job lint (gofmt, go vet)' (PRD-0001:241) is replaced.
15. **docs/engineering-discipline.md:190-193** — A task is a PSB In-Scope item (F-0003#41–#52) or a child of one ... Its plan names the In-Scope fact it serves. No new process rule, routing rule, check script or policy capture starts on its own. *Effect:* Each found issue is a child of an In-Scope item (for example REQ-001/002/004 work under #42) and names that fact. NFR-only items (boundary test, NFR-001 and NFR-002 tests) cannot start on their own.
16. **docs/tests/test-levels.md:43** — Unit tests touch no file, network, or clock. *Effect:* Tests that start git, run `go list`, scan source files or read docs/spec are integration or e2e. The unit level runs in the pre-commit hook (.githooks/pre-commit:141).
17. **docs/engineering-discipline.md:733-734; tests/README.md:1-7** — Unit and integration tests are *_test.go beside the code; the repo-root tests/ holds cross-package end-to-end fixtures only. *Effect:* The harness is Go test code in cmd/layup; fixtures go under tests/; tests/README.md changes in the same PR.
18. **docs/engineering-discipline.md:738-745** — Every component has a unit test; every interface or workflow an integration test; every user-facing scenario an e2e test; every REQ/NFR at least one test. *Effect:* Each item lists tests at these levels, and the traceability matrix gets one test name per requirement.
19. **.github/workflows/ci.yml:310-318** — CI installs only Go and runs go test ./..., -tags=integration and -tags=e2e, each with -timeout 10m. *Effect:* New tests with these tags need no CI change. But if a Go catalog kind needs a tool that ubuntu-latest lacks, the full-setup e2e paths fail (not-active, then S15 refuses); a step that installs it is a gate change with cycle cap 2 (engineering-discipline.md:202-204).
20. **.github/workflows/ci.yml:302; .githooks/pre-commit:126-130** — CI runs `test -z "$(gofmt -l .)"` over the tree; pre-commit runs gofmt on tracked .go files. *Effect:* A bad-on-purpose Go fixture is a patch or a file with no .go suffix; a Go fixture target has its own go.mod.
21. **docs/guardrails.md:150-161** — Start each tool under test with HOME set to a directory under the work directory. *Effect:* The git integration tests and the e2e harness set HOME and the git config to temp values.
22. **docs/engineering-discipline.md:747** — Old tests do not get weakened. *Effect:* cli_test.go:23-25, main_e2e_test.go:21 and the psb golden files stay; the harness keeps the version scenario.
23. **AGENTS.md (Placeholders)** — Never invent a command, path or number; never replace a gap with a guess. *Effect:* The commit identity, the minimum git version, the credential path of S02 and the home of the telemetry and stall types stay open until the Operator decides.

### Risks of area `found`

- The schema-block test needs a not-yet-built list. If the list is never emptied, the test is a silent pass for those blocks.
- If the boundary test parses the Markdown table of packages.md, a prose edit breaks it; if it holds a Go copy, the copy drifts.
- git commit fails on a host or runner with no user.name, unless internal/git sets an identity; the commit time makes setup SHAs differ between runs.
- Isolating the host's git config can remove the credential helper that S02 needs for a private baseline; not isolating it lets autocrlf, init.defaultBranch, hooksPath or gpgsign change the target.
- The tagged integration and e2e files get only the vet subset that go test runs (measured: -atomic, -printf, -stdversion, -tests and others); the full `go vet ./...` (ci.yml:303) does not build them.
- The local toolchain is go1.27.1 and go.mod says go 1.26. stdversion in go vet catches a newer std symbol, but other toolchain differences can still pass locally and fail in CI.
- `go test -timeout 10m` limits the whole test binary of a package. All e2e tests in cmd/layup (full setup runs, verify with two gate runs per kind, the NFR-001 and NFR-002 runs) share one 10-minute limit.
- An autocrlf checkout (issue #48) adds CR to testdata .tsv goldens and to future catalog files and patches; .gitattributes pins none of them.
- A local baseline fixture must be complete enough to pass the baseline's own discipline-tests and link-lint, or no full-setup e2e path reaches S15.
- If `—` passes every type, a wrong `—` in a required column is not found; if no column accepts `—`, valid rows fail.

### Critic points that the reviser of `found` rejected

- Lower-confidence part of the unsupported findings on found-telemetry-schema and found-stalls-schema: no phase-1 source needs Go code for telemetry, prices or stalls at all. *Basis:* PRD-0001-layup.md:179 says phase 1 ships 'the schemas of the telemetry record (REQ-011) and the stall record (REQ-009)', and docs/spec/README.md:75-77 says 'Each record has one schema block. A Go test of the code reads these blocks'. A block with no Go schema leaves the block test incomplete. The items stay, reduced to type, schema and validator, with the home package as an Operator decision (packages.md:70-72).

## Area `gov`

The gov area has no Go code. It holds the rules that set the shape of the plan, and the non-code work that the plan must put in order. (1) T-55n2 fills the PRD §12 Task column for every row, the four Won't rows included, and adds a §13 row. (2) T-55n2 also makes a test-side traceability table with a row for every REQ and NFR. Its DoD part stays an open gap for the idea owner. (3) T-55n2 re-scopes #29, which must also map `layup setup verify` to an ADR-0012 part 5 deliverable and give the new order. It also re-scopes #33, which predates ADR-0016. (4) T-4wrw, the PDR, sets PRD-0001 to Accepted and adds a PDR glossary row. #42 must close before any build issue starts. (5) An Operator decision must say which pilot ends bootstrap mode. Today "pilot" has two meanings (ADR-0012 part 6, and PRD §8/§9 phase 4), and the glossary has no row for it. (6) The pilot inputs are a large item. They need a selected Go problem statement, a target that enforces rulesets, answers with comment-URL sources, and five prose files. The M- answers wait for the setup runner. (7) The pilot run follows the order of setup.md:207-223 and is a UAT, not an E2E test. (8) The pilot numbers come from LAYUP's own records: the defects that one review round missed, the stall events of LAYUP's tasks, and the product-to-process ratio. The rule that reads them is written down before the pilot. (9) The ADR that supersedes ADR-0012 must re-home the part 3 tier names and the "not used" list. It must also fix docs/adr/README.md:42 and every bootstrap summary. (10) A code review of the phase-1 release records REQ-015 and REQ-017. (11) Each milestone of phases 2 to 4 starts with its own specification task. (12) The Operator has setup items: the O-112 items for phase 2 or 3, and a branch-protection change if a new CI job is added in phase 1. (13) Each open issue gets a decision. None blocks a phase-1 item. #15, #21, #34, #48 and #61 are batched into the product items that touch their files, not opened as separate issues. #24 is check-and-close. #68 goes to the phase-2 handoff item. #49 waits for the end of bootstrap mode. LAYUP's own coverage gap stays open, as a constraint. It is not an item.

| Key | Title | Requirements | Size | After | Host |
| --- | ----- | ------------ | ---- | ----- | ---- |
| `gov-prd-task-column` | Fill the Task column of the PRD-0001 section-12 matrix | REQ-012, REQ-001, REQ-002, REQ-004, REQ-007, REQ-009, REQ-011, NFR-001, NFR-002, NFR-003, NFR-004, NFR-005, NFR-006, NFR-007, REQ-015, REQ-016, REQ-017, REQ-018 | small | gov-rescope-29 | `T-55n2` |
| `gov-test-traceability` | A test-side traceability table for LAYUP, and LAYUP's Definition of Done | REQ-012, REQ-001, REQ-002, REQ-004, REQ-007, REQ-009, REQ-011, NFR-001, NFR-002, NFR-003, NFR-004, NFR-005, NFR-006, NFR-007 | small | gov-prd-task-column, gov-rescope-29, found-stdlib-only-test (other area, packages: the test that runs go list -deps; best-guess key) | `T-55n2` |
| `gov-rescope-29` | Re-scope #29's children table to the plan's phase-1 issues | REQ-002, REQ-004, REQ-007, REQ-009, REQ-011 | small | — | `T-55n2` |
| `gov-rescope-33` | Close or rewrite #33 (T-vk3k), whose text predates ADR-0016 | REQ-004, REQ-007, NFR-004 | small | gov-rescope-29 | `T-55n2` |
| `gov-pdr-approval` | T-4wrw: the PDR record and the Operator's approval | REQ-012 | small | gov-prd-task-column, gov-test-traceability, gov-rescope-29, gov-rescope-33 | `T-4wrw` |
| `gov-pilot-definition` | Decide which pilot ends bootstrap mode, and what phase 1 is accepted against | REQ-001, REQ-002, REQ-004, REQ-007, REQ-009, REQ-011, NFR-002 | small | — | O-122 |
| `gov-pilot-inputs` | Prepare the inputs of the first pilot | REQ-002, NFR-003, REQ-016 | large | gov-pilot-definition, setup-runner (other area; best-guess key: the step runner, whose stop table gives the M- and S05 questions) | row 20 |
| `gov-first-pilot` | Run the first pilot: set up one target and run its gate from outside | REQ-002, REQ-004, REQ-007, REQ-016, REQ-018, NFR-001, NFR-002, NFR-003, NFR-004, NFR-006 | large | gov-pilot-inputs, gov-pdr-approval, setup-runner, verify-setup-verify, gate-command, setup-catalog-go (other-area keys are best guesses) | row 20 |
| `gov-pilot-numbers` | Pre-register and record the pilot's defect and stall numbers | — | small | gov-pilot-definition, gov-first-pilot | the ADR that supersedes ADR-0012 |
| `gov-supersede-adr-0012` | The ADR that supersedes ADR-0012 and ends bootstrap mode | — | large | gov-first-pilot, gov-pilot-numbers | the ADR that supersedes ADR-0012 |
| `gov-release-review` | The code review of the phase-1 release for the Won't rows | REQ-015, REQ-017 | small | gov-first-pilot | row 19 |
| `gov-milestone-spec` | A specification task at the start of each later milestone | REQ-012 | large | gov-first-pilot | `T-55n2` |
| `gov-operator-setup-o112` | The Operator's setup items from O-112 (T-hbw8 verdict) | REQ-010, REQ-007 | small | — | `M2a` |
| `gov-operator-ci-protection` | The Operator applies branch protection if phase 1 adds a CI job | REQ-004 | small | setup-catalog-go (other area; best-guess key) | out |
| `gov-issue-15` | #15: stale text about the deleted docs/decisions/ (batched; not a separate issue) | — | small | — | `T-55n2` |
| `gov-issue-21` | #21: the markers check skips quoted file names; stale sentences (batched with the Go port of markers) | NFR-003 | small | verify-setup-verify (other area; best-guess key: the item that ports markers) | row 10 |
| `gov-issue-24` | #24: docs/ci/README.md and AGENTS.md (check and close) | — | small | — | `T-55n2` |
| `gov-issue-34` | #34: ADR-0011 wording, a glossary row for 'rule path', the S04 wording (batched with the setup rule-path register) | REQ-002 | small | setup-rule-paths (other area; best-guess key: the item that writes rule-paths.tsv) | row 15 |
| `gov-issue-48` | #48: autocrlf breaks check facts and setup-check.sh (batched; not blocking on the macOS host) | — | small | — | row 12 |
| `gov-issue-61` | #61: check facts accepts a blank fact made of a tab (batched with the verify item that ports facts) | REQ-001 | small | verify-setup-verify (other area; best-guess key: the item that ports facts) | row 12 |
| `gov-issue-68` | #68: handoff artifacts that are issue comments (phase 2; not blocking phase 1) | REQ-005, NFR-001 | small | gov-milestone-spec | `M2e` |
| `gov-issue-49` | #49: how R11 counts the goal classes of a decision record (waits for the end of bootstrap mode) | — | small | gov-supersede-adr-0012 | out |

### `gov-prd-task-column` — Fill the Task column of the PRD-0001 section-12 matrix

- **What:** T-55n2 writes the task ID of each planned issue into the Task column of PRD-0001 §12, for each phase-1 row, for the later rows that the plan slices, and for the four Won't rows (REQ-015 to REQ-018, through gov-release-review and the setup items). It adds a §13 change-log row. The Test cells stay for each delivering task, but the NFR-007 cell is marked as not proving its criterion (see gov-test-traceability).
- **Specification:** docs/architecture.md §8 The phase loop (Build: each Must requirement of the milestone has a task, each task a requirement); docs/spec/README.md (Phase 1 and later phases)
- **Packages:** —
- **Present state:** PRD:211-213 says T-55n2 fills the Task column. Rows REQ-002 to REQ-011 (except REQ-001), REQ-013 to REQ-018 and NFR-001, 002, 004, 005 have '—' (PRD:218-239). REQ-001 has 'T-dq05, T-zmj6' (PRD:217). REQ-012 has 'T-wjq4, T-0drh'. NFR-003 has T-nfh8, NFR-006 has T-r7zg, NFR-007 has T-mtb9 (PRD:237-241). These are tasks of LAYUP's own repository, not of the phase-1 contracts of docs/spec/.
- **Tests:**
  - discipline: sh docs/prd/prd-lint.sh exits 0: each §12 row still has 6 cells, and the matrix ID set equals the requirement set (prd-lint.sh:109-142)
  - discipline: A one-off command over the §12 rows and the plan's task list shows that each phase-1 Must row has at least one task ID and each task ID is a task of the plan. Its output goes under runs/T-55n2/ as evidence. It is not a new check script (bootstrap rule 1, engineering-discipline.md:190-193; R5, dod-checklist.md:90-91)
  - discipline: sh docs/links/link-lint.sh and sh docs/setup/setup-check.sh exit 0 after the PRD edit
- **External inputs:**
  - New task IDs: 'T-' plus four random characters, with no tasks/<id>.md that exists already (docs/tasks/backlog.md:15-22)
- **Open questions:**
  1. Do the old IDs in the NFR-003, NFR-006 and NFR-007 cells stay, with the new phase-1 build tasks added? PRD:237 'T-nfh8 | check markers' is about LAYUP's own markers, but docs/spec/setup.md:314-330 makes NFR-003 a rule of a target.
  2. Does the plan fill the Task cells of the phase 2 to 4 rows now, or leave '—' until each milestone's own plan? PRD:212 says only 'the implementation plan (`T-55n2`) fills the Task column'.
  3. Does filling the Task column need a §13 row? PRD README:39-40 asks for a row only for 'a changed requirement'. T-hbw8 and T-0drh added a row for edits to §12 (PRD:248-249).

### `gov-test-traceability` — A test-side traceability table for LAYUP, and LAYUP's Definition of Done

- **What:** A copy of the table of docs/tests/traceability-template.md, with one row per test and a row for every REQ and NFR of PRD-0001 (phases 2 to 4 and the four Won't rows included), each with status 'planned' where no test exists. It states how a §12 Test cell that is a discipline check (check facts, check markers, check pin, prd-lint, the CI job lint) gets a row, because 'discipline' is not an allowed Level. It also records that the NFR-007 criterion (`go list -deps ./...` names no module) has no test yet, and names the phase-1 item that owns that test.
- **Specification:** docs/spec/README.md (the schema block: the Go test that reads the schema blocks); docs/spec/packages.md (NFR-007)
- **Packages:** —
- **Present state:** No table exists. docs/tests/traceability-template.md:34-36 is a template only, with Level 'unit | integration | e2e | uat' (traceability-template.md:40-41). dod-checklist.md:24-29 and 80-81 require a green or frozen row for every REQ/NFR, so today no task can tick that box. Tests that exist: TestGoldenRealPSB (internal/psb/check_integration_test.go:12), TestGolden (internal/psb/check_test.go:29), four cli tests (internal/cli/cli_test.go:15,28,38,48), TestBinaryPrintsItsVersion (cmd/layup/main_e2e_test.go:12). No Go file and no CI job runs go list (git grep; ci.yml:294-303 runs gofmt and go vet only).
- **Tests:**
  - discipline: setup-check markers passes: the copy holds no '&lsaquo;…&rsaquo;' cell and not the header 'Task (&lsaquo;task-ID&rsaquo;)', because MK_EXEMPT exempts only docs/tests/traceability-template.md (setup-check.sh:305); each '&lsaquo;…&rsaquo;' left in dod-checklist.md keeps its open-gaps.tsv row (open-gaps.tsv:3)
  - discipline: link-lint passes on the new table; a one-off command, with its output under runs/, shows that the table's REQ/NFR ID set equals PRD-0001's set (dod-checklist.md:80-81) and that each task ID matches a backlog line (traceability-template.md:64-66)
- **External inputs:**
  - The idea owner's answer to 'What is LAYUP's Definition of Done?' (open-gaps.tsv:3)
- **Open questions:**
  1. Where does the table live? traceability-template.md:15-16 says 'A new traceability table copies the table below' and names no path.
  2. Which Level does a discipline check get in the table? traceability-template.md:40-41 allows only unit, integration, e2e or uat, but the §12 Test cells name discipline checks (PRD:217, 237, 240, 241).

### `gov-rescope-29` — Re-scope #29's children table to the plan's phase-1 issues

- **What:** #29's table gets the plan's phase-1 issues in the new order. The order is setup and gate first, then verify, because verify needs layup gate (setup.md:256; internal/verify imports internal/gate, packages.md:44). It maps `layup setup verify` to one of the four ADR-0012 part 5 deliverables, because part 5 does not name verify, and 'the kit's checks run from outside' is now verify (ADR-0016). Rows 1-3 are done. The old rows 4-6 (T-b97r, T-vk3k, T-tmhw) are replaced, a stall-record row is added, and each new task gets a backlog line. Each new issue text follows the task-issue template, with a Solution note (R3). This closes #42's criterion 3.
- **Specification:** docs/spec/setup.md The boundary of phase 1; docs/spec/setup.md The checks of `layup setup verify`; docs/spec/records.md REQ-011 — The telemetry record; docs/spec/records.md REQ-009 — The stall record; docs/spec/packages.md; docs/architecture.md §8 The phase loop (Build)
- **Packages:** —
- **Present state:** #29 (.ctx/issues/29.md:15-22) lists six children. T-b97r is the 'program steps (human_decision = no) of docs/setup/steps.tsv', but docs/spec/setup.md:77-100 now owns the step table. T-tmhw has 'a completeness check', but PRD:179 and records.md:144-148 give schemas only. There is no stall-record row and no verify row. #29:7 gives the old order 'Armature setup → stack-dependent gates'. backlog.md:34 names the four deliverables. The issue template is at docs/templates/github/ISSUE_TEMPLATE/task.md:16-29, used by hand (docs/templates/README.md:10-13).
- **Tests:**
  - discipline: link-lint passes after the backlog lines are added; backlog.md has one line per open task (backlog.md:8-13)
- **External inputs:**
  - Forge edits of #29 with the layup-agent App token, done by the T-55n2 author (no gh run in this read)
  - The Operator's confirmation on #42 that each remaining #29 child is an issue of the plan (#42:38)
- **Open questions:**
  1. Is each new phase-1 issue a direct child of #29, or a child of one of the four part 5 deliverable tasks? ADR-0012:48 says 'Each deliverable named here … is one task under part 2; #29 stays the parent issue of the four engine deliverables'.
  2. Do the planned IDs T-b97r and T-tmhw (never opened) stay for the new issues? backlog.md:15 says an ID is 'assigned once and never reused'.
  3. Does T-55n2 open the issues, or only list them so that each is opened in order after T-4wrw? #29:13 and #42:22 say 'opened in order; each has its own plan and plan review'.

### `gov-rescope-33` — Close or rewrite #33 (T-vk3k), whose text predates ADR-0016

- **What:** #33 says that the Go gates run outside the target, in LAYUP's own runner (ADR-0011 decision 4, O-11). ADR-0016 replaced that: native gates go in the target, and layup gate runs them from outside. The item either closes #33 as superseded by the plan's gate and Go-catalog issues, or rewrites #33 as one of them.
- **Specification:** docs/spec/gate.md REQ-004 — The stack gates of a target, run from outside; docs/spec/gate.md NFR-004 — A check that is not active is not passed; docs/spec/setup.md The stack catalog
- **Packages:** —
- **Present state:** .ctx/issues/33.md:5 says 'they run outside the target, in LAYUP's own runner'. 33.md:9 carries the runner fork. 33.md:17 says 'To be set by this task's plan (R12) when it starts, after T-dq05 and T-b97r'. ADR-0016:14-18 supersedes O-10 and O-11 for the gates. ADR-0016:50-56 and docs/spec/gate.md answer how a missing runner shows as not-active.
- **Tests:**
  - discipline: link-lint passes on any backlog or task-file change that the close or the rewrite makes
- **External inputs:**
  - A forge edit or close of #33 (App token)
- **Open questions:**
  1. Close #33 and open new issues, or keep #33 and its ID T-vk3k for the layup gate issue? R2 (issue-workflow.md:38-40) prefers not to duplicate an existing thread.

### `gov-pdr-approval` — T-4wrw: the PDR record and the Operator's approval

- **What:** docs/pdr/PDR-0001.md records the Operator's approval (O-14) of PRD-0001, the architecture and the implementation plan, and that each requirement traces to a PSB fact. T-4wrw sets the Status of PRD-0001 to Accepted, changes the index row and the onboarding sentence, adds a §13 row, and adds a glossary row for PDR. When it merges, #42 can close and #29's build issues can start.
- **Specification:** docs/architecture.md §14 Coverage; docs/architecture.md §15 Known limits
- **Packages:** —
- **Present state:** docs/pdr/ does not exist. #42:32 lists T-4wrw as child 6. PRD-0001's Status is 'Draft' (PRD:11; prd/README.md:64). onboarding-for-engineers.md:135 says 'Draft until the PDR records the Operator's acceptance'. The glossary has no PDR row, but #42, backlog.md:33, architecture.md:1476 and onboarding-for-engineers.md:135 use the term (engineering-discipline.md:602-615 requires a row).
- **Tests:**
  - discipline: prd-lint, link-lint and setup-check (adapted, markers, glossary) exit 0 with the new docs/pdr/ file, the Status 'Accepted' (an allowed value, prd/template.md:11) and the PDR glossary row
- **External inputs:**
  - The Operator's approval (O-14: 'I approve it myself', #42:15), as an issue comment copied into Git (#42:37, Invariant 1)
- **Open questions:**
  1. Does the PDR also record the approval of docs/spec/ (child 4a, T-0drh)? #42:5 says 'the specification, the architecture, and the implementation plan, approved by the Operator', but the demo sentence of #42:5 names only 'PRD-0001, the architecture document, and the plan'.

### `gov-pilot-definition` — Decide which pilot ends bootstrap mode, and what phase 1 is accepted against

- **What:** An Operator decision on an issue, copied into Git, with a glossary row for 'pilot'. The decision says three things. First, which pilot ends bootstrap mode: the first pilot of ADR-0012 part 6 after phase 1, or the phase-4 pilot of PRD §9. Second, whether that pilot also measures the PRD §8 baseline. Third, against which §7.1 criteria the idea owner accepts phase 1. O-115 keeps the §7.1 criteria as they are, so the plan cannot rewrite them to fit phase 1 without a new decision.
- **Specification:** docs/spec/setup.md The boundary of phase 1 (row: the pilot baseline, phase 4); docs/architecture.md §14 Coverage (row: ADR-0012 part 6)
- **Packages:** —
- **Present state:** ADR-0012:49 says 'The first pilot is done when layup has set up one target repository from a problem statement and has run that target's gate from outside'. PRD:144-146 says 'the first pilot measures the baseline with the current process'. setup.md:41 gives 'The pilot baseline | no | layup run, phase 4'. PRD:182 puts 'the pilot on two problem statements' in phase 4. The phase-1 §7.1 criteria say 'In the pilot' and need later work: REQ-001 needs the review of meaning (PRD:114; psb-check.md:58-62, phase 2). REQ-004 needs all four kinds, but three stay pending in phase 1 (PRD:117; architecture.md:451-452; setup.md:312). REQ-007 needs PRs on the target (PRD:120). REQ-009 and REQ-011 need a writer from a later phase (records.md:144-146, 197-198). NFR-002 needs a fresh session to continue a task (PRD:133). T-0drh.md:22-23 records O-115: 'The §7.1 criteria of PRD-0001 stay as they are.' The glossary has no Pilot row.
- **Tests:**
  - discipline: setup-check glossary and link-lint pass with the new 'Pilot' row; the decision record links the issue comment
- **External inputs:**
  - The Operator's decision (and the idea owner's, for the acceptance criteria of phase 1)
- **Open questions:**
  1. Does the Operator want a new decision here, or does R13 require that the two readings be rewritten in the sources? R13 says that a text with two honest readings is rewritten, not settled by a guess.

### `gov-pilot-inputs` — Prepare the inputs of the first pilot

- **What:** The idea owner selects one pilot problem statement whose stack has a catalog entry, which in phase 1 means Go only. The Operator creates the empty target repository, which must enforce rulesets (a public repository, or a plan that has them). The Operator names the issue where each S01-, Q-NNN and M- answer is posted, so that each answers.tsv row has a comment URL as its source. The Operator also supplies the prose files for S07, S08, S09 and S14, which must pass the target-form checks. The M- answers and the S05 link files are known only from a layup setup stop table.
- **Specification:** docs/spec/setup.md The boundary of phase 1; docs/spec/setup.md The command `layup setup`; docs/spec/setup.md The steps; docs/spec/setup.md The answers; docs/spec/setup.md The stack catalog; docs/spec/setup.md The checks of `layup setup verify`; docs/architecture.md §5 Intake and setup; docs/architecture.md §15 Known limits (L-A1, L-A6, L-B1)
- **Packages:** —
- **Present state:** None. PRD:203 (§11 question 6) is open: 'which two, with which stacks; the idea owner selects them'. The host has go1.27.1, and go.mod says 'go 1.26'. LAYUP's own versions of the five prose files have 164, 162, 412, 73 and 185 lines. The manual setup of the same rows in LAYUP took four tasks (T-vpty 46 min, T-xgz4 63 min, T-7ndb 56 min, T-6rg3 30 min), and markers took 89 min (docs/setup/record-T-n1hp.md:37-50).
- **Tests:**
  - discipline: Each answers.tsv row has 'by' and a 'source' that is a comment URL or an accepted fact citation (setup-answers schema, setup.md:142-148), checked by the first run of layup setup S01 and S10 on the inputs
  - uat: Given the inputs, when the Operator runs layup setup, then no stop names an input that is missing; the idea owner's choice of the problem statement is the only intent decision, at Decision Point 1 (REQ-016, PRD:129)
- **External inputs:**
  - The idea owner's choice of a pilot problem statement (Decision Point 1)
  - A target repository on GitHub that the Operator owns (L-A1, architecture.md:1515), on a plan that enforces rulesets (architecture.md:324-328, K15)
  - The baseline's repository URL for S01-baseline
  - The answers to the S01-, Q-NNN and M- questions, posted as issue comments (Decision Point 2)
  - The prose files for S07, S08, S09 and S14 (docs/spec/setup.md:69, 92-94, 99)
  - The Go toolchain and git on the host (L-B1, architecture.md:1546-1547)
- **Open questions:**
  1. PRD §11 question 6 asks for the two problem statements of REQ-014 (phase 4), but ADR-0012 part 6 needs one. Must the first pilot's brief be one of the two REQ-014 briefs?
  2. Must the App be installed on the pilot target in phase 1? setup.md:26 gives 'Start 1: the empty repository on the forge, the App installed | no | the Operator, by hand, as today', but no phase-1 command uses the App.
  3. Which issue holds the pilot answers? Phase 1 has no Intake issue (setup.md:31, 34).

### `gov-first-pilot` — Run the first pilot: set up one target and run its gate from outside

- **What:** The Operator follows the order of setup.md:207-223. First, layup setup WORK runs S01 to S14 with its stops. It stops at S15 with O-verify. The Operator then runs layup setup verify WORK > WORK/out/verify.tsv. The next layup setup WORK does S15 and refuses any row that is not pass or clear. Then the Operator runs commands.sh: the root-commit push, the layup-setup push, the layup-records push, and the ruleset apply, in that order. After that, layup gate --base --head runs on the target from outside. This meets ADR-0012 part 6 through W-02 step 5 and one layup gate run. W-04 steps 5 and 9 are phase-2 work (gate.md:96-100). The evidence goes under runs/.
- **Specification:** docs/spec/setup.md Where the records go in phase 1; docs/spec/setup.md The checks of `layup setup verify`; docs/spec/gate.md The command; docs/spec/gate.md Not in phase 1; docs/spec/records.md NFR-002 — A target is independent of LAYUP; docs/architecture.md §14 Coverage (row: ADR-0012 part 6); docs/architecture.md §5 Intake and setup; docs/architecture.md §6 Native gates and rule protection; docs/architecture.md §15 Known limits (L-A6)
- **Packages:** —
- **Present state:** None. Only `layup version` and `layup psb check` exist (internal/cli/cli.go). ADR-0016:72-73 says 'ADR-0012 part 6 stays true: layup sets up the target and runs its gate from outside.'
- **Tests:**
  - uat: Given the pilot inputs, when the Operator runs the steps of setup.md:207-223 and layup gate --base --head on the target, then verify.tsv holds only pass or clear and layup gate prints one verdict per kind. The run is recorded under runs/ as evidence. It reads the live network, so it is not an automated E2E test (test-levels.md:63-81; guardrails.md:276-280)
  - uat: On the pilot target, a pull request whose own CI jobs run and pass merges with LAYUP absent (NFR-002, records.md:131-134)
  - uat: An audit of every setup value finds that its cited source supports it (REQ-002, NFR-003; PRD:115, 134). A byte comparison shows that the baseline's rules are unchanged except the recorded adapted values (REQ-018, PRD:131)
  - uat: The idea owner accepts or rejects each phase-1 requirement, against the criteria that gov-pilot-definition fixes, by an issue comment copied into Git
  - e2e: The automated twin of the pilot path (setup, verify, gate) runs with go test -tags=e2e against a local fixture baseline and a local target, with no network. It is owned by the setup, verify or gate area, and listed here as a precondition
- **External inputs:**
  - The pilot inputs (gov-pilot-inputs)
  - The Operator runs commands.sh with the Operator's own login, from a clone with no hooks (docs/spec/setup.md:218-223)
  - Network access to the baseline and to GitHub
- **Open questions:**
  1. L-A6 (architecture.md:1522-1525): S02 pins the baseline's latest commit, not LAYUP's pin, and S10 uses the MK_EXEMPT that the engine embeds at its own version. If the baseline changed, does the pilot pin the baseline to a known commit, or accept a stop?

### `gov-pilot-numbers` — Pre-register and record the pilot's defect and stall numbers

- **What:** Before the pilot runs, the plan of the pilot task writes down the rule by which the ADR that ends bootstrap mode reads the numbers: what counts as a defect or a stall, and which values restore or re-decide the full gate. At the pilot, the task counts three things: the defects that the pilot finds in phase-1 code that passed its one review round, the stall events of LAYUP's own phase-1 tasks (cycle-cap ends, 'not mergeable' verdicts, Operator decisions per task, from docs/tasks/T-*.md and runs/), and the product-to-process ratio of docs/tasks/completed.md. The counts go under runs/.
- **Specification:** —
- **Packages:** —
- **Present state:** None. ADR-0012:56 says 'Fewer review rounds find fewer defects. That risk is accepted for the bootstrap phase and measured by the pilot'. engineering-discipline.md:203-204 says 'the pilot measures what one pass misses', and :237-238 says 'The ADR that ends this mode reads the pilot's defect and stall numbers'. guardrails.md:116-120 says 'the product-to-process ratio of tasks/completed.md is read at each pilot'.
- **Tests:**
  - discipline: link-lint passes on the record under runs/; the pre-registered rule is on the pilot issue or in its plan-review comment before the pilot's first run (guardrails.md:21-26 'Write the pass/fail numbers first')
- **External inputs:**
  - The Operator's acceptance of the pre-registered rule
- **Open questions:**
  1. No source defines how the defects and stalls are counted. Is a 'stall' of LAYUP's own work a cycle-cap end, a 'not mergeable' verdict, or the target's stall record of REQ-009? engineering-discipline.md:237 says only 'the pilot's defect and stall numbers'.

### `gov-supersede-adr-0012` — The ADR that supersedes ADR-0012 and ends bootstrap mode

- **What:** The task that closes the first pilot opens a new ADR (the next free number is 0026). The ADR re-decides the full gate with the pilot's numbers and sets ADR-0012's Status to Superseded. It re-homes the tier names, the 'not used' model list and the reviewer rule of part 3, or states that they end, because two documents point to part 3 as their home. It updates every summary of bootstrap mode, the ADR-0005 and ADR-0006 Status lines, and the stale line 'the next constitutional ADR is 0013' (docs/adr/README.md:42).
- **Specification:** —
- **Packages:** —
- **Present state:** ADR-0012 is Accepted (docs/adr/0012-build-layup-in-bootstrap-mode.md:7). The operative text is engineering-discipline.md:177-239. Other places that state or cite the rules: AGENTS.md (Bootstrap mode paragraph); issue-workflow.md:143-148, 208-212; README.md:21-26; onboarding-for-engineers.md:156-164; glossary.md:87 (Cycle cap), :93 (Material), :144 (Bootstrap mode); guardrails.md:116-120, 166-173 (the 'not used' list of part 3); engineering-discipline.md:146-151 (the tier names come from part 3); docs/spec/README.md:77; docs/adr/README.md:42 and :151; the Status lines of ADR-0005 and ADR-0006.
- **Tests:**
  - discipline: adr-lint, link-lint and setup-check (glossary, adapted) exit 0; ADR-0012's Status names the new ADR; git grep -n -i bootstrap shows each remaining mention as history or as a link to the new rule
- **External inputs:**
  - The pilot's numbers (gov-pilot-numbers)
  - The Operator's decision on the full gate and the acceptance of the ADR
- **Open questions:**
  1. Does this ADR need a panel? Bootstrap rule 5 (engineering-discipline.md:225-227) allows a panel only for an ADR that changes the product architecture; this ADR changes the process.
  2. Is it in the phase-1 milestone, or the first item after it? ADR-0012:49 says only 'the task that closes the pilot opens that ADR'.
  3. Does the bootstrap substitution end for each issue already in flight when the ADR is accepted, or only for new plans?

### `gov-release-review` — The code review of the phase-1 release for the Won't rows

- **What:** A code review of the phase-1 release records that no code path modifies base LLM weights or trains a model (REQ-015), and that no code path calls a cloud provider or hosting platform to modify it (REQ-017). The review is recorded on an issue and copied into Git. Its task ID fills the §12 Task cells of REQ-015 and REQ-017.
- **Specification:** docs/spec/packages.md (the import rules)
- **Packages:** —
- **Present state:** None. PRD:128 and :130 say 'a code review of each release records it'. PRD:184 says 'The four Won't rows (REQ-015 to REQ-018) hold in every phase.' PRD:231-234 has '—' in their Task cells.
- **Tests:**
  - discipline: The review record passes review-record-lint in CI; a one-off grep over the release's imports and its exec calls (only git, the gate commands, sh) is kept under runs/ as the deterministic part of the claim (R5)
- **External inputs:**
  - A reviewer whose model differs from the authors' (bootstrap rule 3)
- **Open questions:**
  1. No source defines 'release' for LAYUP. Is it the end of phase 1, or each tagged build?

### `gov-milestone-spec` — A specification task at the start of each later milestone

- **What:** Each milestone of phases 2 to 4 starts with a specification task, as T-0drh did for phase 1. The task writes the docs/spec/ sections for the milestone's requirement IDs and the packages.md rows of its packages, before the milestone's build issues start.
- **Specification:** docs/spec/README.md (written one milestone at a time); docs/spec/packages.md
- **Packages:** —
- **Present state:** docs/spec/README.md:9-12 says 'It is written one milestone at a time (task T-0drh, #74; decision O-114). This version covers phase 1 of PRD-0001 §9 only. A later milestone adds its own sections'. PRD:125 (REQ-012) says the same. packages.md:70-72 names the later phases.
- **Tests:**
  - discipline: prd-lint and link-lint pass; the Go test that reads the schema blocks of docs/spec/ passes on the new sections

### `gov-operator-setup-o112` — The Operator's setup items from O-112 (T-hbw8 verdict)

- **What:** The Operator creates the GitHub App layup-watch for the dead-man job. The Operator also gives layup-agent a private key for the LAYUP host and the commit-statuses permission, and installs it on each target. No phase-1 command needs these items. The plan puts them before the first item of a later phase that posts a status or runs the dead-man job.
- **Specification:** docs/architecture.md §11 Stalls (the dead-man job); docs/architecture.md §3 Records and identities; docs/spec/gate.md Not in phase 1 (the commit status layup/gates); docs/spec/setup.md The boundary of phase 1
- **Packages:** —
- **Present state:** docs/tasks/T-hbw8.md:137-139: 'Open for the Operator: the setup work of O-112 (layup-watch, a private key and commit statuses for layup-agent) when LAYUP is built'. runs/T-hbw8/operator-decisions.md:317 has the exact text. Today layup-agent is installed on pharzam/layup only, and agent writes use its user access token from the device flow (operator-decisions.md:82). architecture.md:139-141 says that layup run makes an installation token from the App's private key. The status layup/gates is phase 2 (gate.md:96-97).
- **Tests:**
- **External inputs:**
  - The Operator's GitHub App settings, keys and installations, and the Operator's read-back of them before the first issue that needs them
- **Open questions:**
  1. Which phase needs layup-watch? setup.md:31 puts the dead-man job's first notice in Start 3 (phase 2), but records.md:88 puts the dead-man job's record in phase 3, and REQ-010 is phase 3 (PRD:85).

### `gov-operator-ci-protection` — The Operator applies branch protection if phase 1 adds a CI job

- **What:** If the catalog fixtures (or any other phase-1 test) run in a new CI job, the same pull request changes docs/setup/branch-protection.json, and the Operator applies it with an administration-scoped token and records the read-back. If the fixtures run inside the existing job 'tests', no Operator item is needed. The plan states which choice it makes.
- **Specification:** docs/spec/setup.md The stack catalog
- **Packages:** —
- **Present state:** setup.md:301-302 says 'LAYUP's own CI runs each fixture of each entry (§6); that test comes with the code (#29).' docs/ci/README.md:62-72 says that the protection file has twelve checks, one per job. docs/ci/README.md:88-95 says that writing the setting needs an administration-scoped token, run by an operator.
- **Tests:**
  - discipline: setup-check protection passes: the contexts equal the job names (setup-check.sh:553)
- **External inputs:**
  - The Operator's administration-scoped token and the read-back output

### `gov-issue-15` — #15: stale text about the deleted docs/decisions/ (batched; not a separate issue)

- **What:** Most of #15 is fixed in the tree. The guardrails §2 pitfalls (docs/guardrails.md:319, 347-349) still use the archive as their example. The decision is to keep them with a reason on #15, or to reword them. This goes into the first product task that writes a guardrails §2 lesson (ADR-0012:60), and that task then closes #15. It does not block a phase-1 item.
- **Specification:** —
- **Packages:** —
- **Present state:** engineering-discipline.md:533 and glossary.md:40 say that the directory was deleted. guardrails.md:319 and 347-349 remain. The issue is open.
- **Tests:**
  - discipline: link-lint passes after the edit; the grep for 'docs/decisions/' is kept as evidence under runs/, not as a test

### `gov-issue-21` — #21: the markers check skips quoted file names; stale sentences (batched with the Go port of markers)

- **What:** Fix `markers` in setup-check.sh for file names that git quotes, with a fixture, and fix the stale sentences and the question check. #21 does not block phase 1. But layup setup verify must pass and fail on the same fixtures as docs/setup/tests/run.sh, so the fix goes into, or before, the verify item that ports markers.
- **Specification:** docs/spec/setup.md The checks of `layup setup verify`
- **Packages:** —
- **Present state:** .ctx/issues/21.md:13-17 lists five residuals. The issue is open.
- **Tests:**
  - discipline: A fixture with a quoted file name and an unlisted marker fails markers (#21:21), through docs/setup/tests/run.sh
  - integration: The Go markers check fails on the same fixture repository; it lists the tracked files through internal/git, so it crosses a real seam (test-levels.md:51-61; packages.md:39, 44)

### `gov-issue-24` — #24: docs/ci/README.md and AGENTS.md (check and close)

- **What:** Both acceptance criteria of #24 are met in the tree. The item closes #24 with the evidence and one note: docs/ci/README.md:37 says 'Every job restores its check scripts', but nested-checkout-check.sh is not restored, on purpose (ci.yml:246-251). The fix of that sentence is batched into the first product task that edits docs/ci/README.md or ci.yml.
- **Specification:** —
- **Packages:** —
- **Present state:** AGENTS.md:114 lists sh docs/setup/setup-check.sh. docs/ci/README.md:37-38 says that the jobs restore their checks from the default branch. The README has no example array of contexts now: it names docs/setup/branch-protection.json (README:62-72), and check protection keeps the contexts equal to the job names (setup-check.sh:553). The issue is open.
- **Tests:**
  - discipline: link-lint passes after the one-sentence fix
- **External inputs:**
  - A forge close of #24 with the evidence (App token)

### `gov-issue-34` — #34: ADR-0011 wording, a glossary row for 'rule path', the S04 wording (batched with the setup rule-path register)

- **What:** Add a glossary row for 'rule path', because phase 1 writes rule-paths.tsv (setup.md:163, S15). Reword the glossary's 'runs role agents'. Answer S04 'ADR for the pin' with the spec (setup.md:89). Answer point 2 (the header row of a register) with docs/spec/README.md:86-88, which gives 'no-header' for open-gaps.tsv. A change to ADR-0011's body, including its unwrapped line, needs a new ADR or a Status pointer (#34:15). Not blocking.
- **Specification:** docs/spec/setup.md The rule-path register; docs/spec/setup.md The steps; docs/spec/README.md (the schema block, no-header)
- **Packages:** —
- **Present state:** glossary.md:140 still says 'runs role agents'. The glossary has no 'rule path' row. docs/setup/README.md:78 and steps.tsv:5 still say 'ADR for the pin'. The sentence 'for ADR-0011 to decide' is no longer in docs/setup/README.md (the only matches of 'decide' are lines 14 and 21, with other text). The issue is open.
- **Tests:**
  - discipline: setup-check glossary and adapted, adr-lint and link-lint exit 0 with the new row and wording
- **Open questions:**
  1. Step 2 is still defined in no file (#34:11), and backlog.md:34 and #42:20 use 'Step 2'. Does the plan rename it, or define it?

### `gov-issue-48` — #48: autocrlf breaks check facts and setup-check.sh (batched; not blocking on the macOS host)

- **What:** Pin LAYUP's raw facts files and setup-check.sh with a line-ending rule in .gitattributes, with a test that fails without the pin. The fix goes into the first product task that edits .gitattributes or setup-check.sh (ADR-0012:60). It changes nothing in a target.
- **Specification:** —
- **Packages:** —
- **Present state:** No facts path and not setup-check.sh has a line-ending pin: git check-attr -a prints nothing for docs/setup/setup-check.sh or the PSB file. .gitattributes:201 is 'docs/setup/tests/facts/bad-answers-blank/** -whitespace', which is a different attribute (.gitattributes:198-200). The issue is open.
- **Tests:**
  - discipline: On a clone with core.autocrlf=true, sh docs/setup/setup-check.sh --only facts . exits 0 (#48:5), through a fixture of docs/setup/tests/run.sh

### `gov-issue-61` — #61: check facts accepts a blank fact made of a tab (batched with the verify item that ports facts)

- **What:** The one-character fix in setup-check.sh (tr -d ' \t' for the F-0001 and F-0003 loop), with a fixture. It serves F-0003#41. The Go facts check is in the target's form and does not count LAYUP's facts (setup.md:246), so the fix goes with the verify item that ports facts, to keep the shared fixtures equal.
- **Specification:** docs/spec/setup.md The checks of `layup setup verify`
- **Packages:** —
- **Present state:** docs/setup/setup-check.sh:164 still uses tr -d ' '. The issue is open.
- **Tests:**
  - discipline: A fixture under docs/setup/tests/facts/ with a fact of one tab fails check facts
- **Open questions:**
  1. setup.md:238-239 says that the engine passes and fails 'on the same fixtures as docs/setup/tests/run.sh', but setup.md:246 says that LAYUP's counts do not apply to a target. Which facts fixtures must the Go check share?

### `gov-issue-68` — #68: handoff artifacts that are issue comments (phase 2; not blocking phase 1)

- **What:** The task that builds layup handoff check (REQ-005, phase 2) decides how the plan, the plan review and the review record get a Git path, or a short ADR decides it. The plan puts #68 into the phase-2 handoff item.
- **Specification:** docs/architecture.md §8 The phase loop (Handoffs)
- **Packages:** —
- **Present state:** .ctx/issues/68.md:3 is open, normal priority, 'no handoff row exists yet'.
- **Tests:**
  - integration: In phase 2, a handoff row for each of the kinds plan, plan-review and review-record validates against its artifact

### `gov-issue-49` — #49: how R11 counts the goal classes of a decision record (waits for the end of bootstrap mode)

- **What:** Make R11 state that a decision record is one goal (O-15 on #46). ADR-0012:47 says that #49 waits for the end of bootstrap mode, so the plan puts it after the ADR that supersedes ADR-0012.
- **Specification:** —
- **Packages:** —
- **Present state:** issue-workflow.md:123-134 does not have the rule. .ctx/issues/49.md is open.
- **Tests:**
  - discipline: link-lint passes, and AGENTS.md and the other R11 summaries agree with the new text (R10)
- **Open questions:**
  1. Does the ADR that supersedes ADR-0012 take over #49? Both change R11.

### Constraints of area `gov`

1. **docs/adr/0012-build-layup-in-bootstrap-mode.md:48** — Each deliverable named here — the #45 batch, PRD-0001, layup setup, layup gate, the telemetry record, the stall record — is one task under part 2; #29 stays the parent issue of the four engine deliverables (layup setup, layup gate, the telemetry record, the stall record), and #45 stays a child of #42. *Effect:* #29 stays the parent of all phase-1 engine issues. If the plan cuts one deliverable into many issues, this conflicts with 'is one task', unless each issue is a child of that task (bootstrap rule 1 allows 'a child of one') or the Operator decides. Part 5 does not name layup setup verify, so the plan must map it to a deliverable.
2. **docs/adr/0012-build-layup-in-bootstrap-mode.md:49; docs/adr/0016-put-the-native-stack-gates-in-the-target.md:72-73** — The first pilot is done when layup has set up one target repository from a problem statement and has run that target's gate from outside … the task that closes the pilot opens that ADR … / ADR-0012 part 6 stays true: layup sets up the target and runs its gate from outside. *Effect:* The plan needs a pilot milestone after layup setup, layup setup verify and layup gate, read as layup gate (ADR-0016). The pilot task opens the ADR that supersedes ADR-0012.
3. **docs/adr/0012-build-layup-in-bootstrap-mode.md:56; docs/engineering-discipline.md:203-204, 237-238** — Fewer review rounds find fewer defects. That risk is accepted for the bootstrap phase and measured by the pilot … / the pilot measures what one pass misses … The ADR that ends this mode reads the pilot's defect and stall numbers. *Effect:* The plan schedules the count of these numbers from LAYUP's own task records and issues (gov-pilot-numbers), not from a target's stalls.tsv.
4. **docs/guardrails.md:21-26** — A decision rule chosen after seeing the result is a fitted parameter, not a rule. Write the pass/fail numbers first, somewhere they cannot be quietly edited. *Effect:* The rule that reads the pilot numbers is written on the pilot issue or in its plan review before the pilot runs.
5. **docs/guardrails.md:116-120** — a plan names the PSB In-Scope fact (F-0003#41–#52) it serves, or the task does not start (ADR-0012, part 1); and the product-to-process ratio of tasks/completed.md is read at each pilot. *Effect:* Each planned issue names its In-Scope fact. The pilot task reads the product-to-process ratio.
6. **docs/engineering-discipline.md:190-193; docs/adr/0012-build-layup-in-bootstrap-mode.md:44** — A task is a PSB In-Scope item (F-0003#41–#52) or a child of one, a defect that blocks such a task, or a documentation fix that a task leaves stale. … No new process rule, routing rule, check script or policy capture starts on its own. *Effect:* An NFR-only item (for example internal/tsv or internal/git) must be a child of an In-Scope item. A deterministic claim of a gov item is settled by a one-off command with its output under runs/, not by a new check script.
7. **docs/setup/open-gaps.tsv:2; docs/issue-workflow.md:250** — &lsaquo;add a coverage gate&rsaquo; … Operator decision O-7 (#8): none yet; set it after the first Go code exists and a baseline is measured. *Effect:* LAYUP's own coverage gate serves no PRD requirement, and bootstrap rule 1 bars a new check unless a product task needs it. So the plan keeps the marker as an open gap. The target's coverage floor (L-B2, architecture.md:1548-1551) is a different gap, and it belongs to the catalog and setup items.
8. **docs/engineering-discipline.md:194-212** — The plan review is one comment … with Verdict, Budget maximum and Cycle cap … The review is one round … the cycle cap is 1, and 2 when the change touches a gate (a check script, a hook, CI or branch protection). *Effect:* Each issue gets one plan-review comment and one review round. An issue that changes CI, a hook or setup-check.sh (for example the batched #21, #48 or #61 fixes, or a new CI job for the catalog fixtures) has cycle cap 2. Sizes must fit one round.
9. **docs/engineering-discipline.md:197-199; docs/issue-workflow.md:106-134** — the goal-class count of R11 is not applied to a task whose deliverable is one artifact and its registration (a record, a decision, a document, a fact). / One issue is one actionable, demoable goal … The plan names the one demo … without an 'and' joining two outcomes. *Effect:* The R11 count applies to every build issue. layup setup verify has 14 checks marked yes in phase 1, and 18 rows with the five Go gate kinds (setup.md:241-258; gate.md:22), so one verify issue can be counted as several goals. Each issue has one demo sentence with no 'and'. The Operator decides a reject that is about the count only.
10. **docs/issue-workflow.md:152-206** — Before the first test, turn the issue into an ordered plan … Review the plan once … The plan-review comment carries a Budget maximum and a Cycle cap. *Effect:* The implementation plan does not replace each issue's own R12 plan and review. It gives each issue its goal, its DoD and its order.
11. **docs/templates/README.md:10-13; docs/templates/github/ISSUE_TEMPLATE/task.md:16-29; docs/issue-workflow.md:172-180** — A LAYUP task issue and its pull request follow the shape of these files by hand … ## Solution note (R3) … Chosen … Rejected … Important tradeoffs. *Effect:* Each issue text that the plan writes uses the task-issue template with a Solution note. T-55n2 records its selected slicing and the rejected one (one issue per part 5 deliverable, against finer R11 children), with the tradeoffs.
12. **docs/engineering-discipline.md:231-235; docs/tests/dod-checklist.md:82-85** — An operative rule has one canonical home. Another document may link to it and may give a clearly marked non-operative summary, but it must not restate the rule as an independent requirement. *Effect:* The plan document and its issue texts link to the bootstrap gate, R11, R12 and the phase-1 boundary, or give a marked non-operative summary that gets a semantic-agreement review.
13. **docs/issue-workflow.md:19-40** — Every change needs an open issue before any commit or pull request … Before opening an issue, search the open and closed issues. If the work is part of a larger one, open it as a child/sub-issue and link the parent. *Effect:* Each planned issue is a child of #29, with a duplicate check. #33 already exists for the gates. Task IDs go in commit subjects, and Closes/Refs #N goes in PR bodies.
14. **.ctx/issues/42.md:38** — #29's remaining children are re-scoped to the plan's issues. *Effect:* The #29 children table is rewritten before #42 closes, by T-55n2 or by T-4wrw.
15. **.ctx/issues/42.md:5; .ctx/issues/42.md:20** — approved by the Operator before the core engine is built further … #29 (Step 2, on hold until this closes) *Effect:* No phase-1 build issue starts before T-4wrw merges and #42 closes.
16. **docs/prd/PRD-0001-layup.md:211-213** — the implementation plan (T-55n2) fills the Task column, and each delivering task fills the Test column. *Effect:* T-55n2 writes task IDs in §12 and leaves the Test cells to the build tasks. So the plan gives each issue its task ID.
17. **docs/prd/prd-lint.sh:109-142** — A row whose first cell is REQ/NFR-NNN with 6 cells is a matrix row; the matrix ID set must equal the requirement set. The Task and Test cells are not checked. *Effect:* Every §12 row keeps exactly 6 cells, with no unescaped '|' in a Task or Test cell. prd-lint does not prove the Task column is right, so a one-off command does.
18. **docs/prd/README.md:39-40** — a changed requirement is a new ## 13. Change log entry, never a silent edit. *Effect:* An edit of PRD-0001 by the plan or by the PDR adds a §13 row, as T-hbw8 and T-0drh did (PRD:248-249).
19. **docs/onboarding-for-engineers.md:135; docs/prd/template.md:11** — PRD-0001 … (Draft until the PDR records the Operator's acceptance). / Status: Draft / Accepted / Superseded by PRD-NNNN *Effect:* T-4wrw sets PRD-0001's Status to Accepted and updates the index row (prd/README.md:64) and the onboarding sentence.
20. **docs/prd/PRD-0001-layup.md:179, 184** — Phase 1 ships layup psb check (rule gaps), layup setup as a step runner with no forge call and no role session, layup setup verify, layup gate (REQ-004, REQ-007), the schemas of the telemetry record (REQ-011) and the stall record (REQ-009) … NFR-001 to NFR-007 hold from the first release. / The four Won't rows (REQ-015 to REQ-018) hold in every phase. *Effect:* The phase-1 milestone has these contents and nothing more. REQ-009 and REQ-011 are schemas with a Go test only. The NFRs and the Won't rows bind every phase-1 issue, and a release review records REQ-015 and REQ-017 (PRD:128, 130).
21. **docs/tasks/T-0drh.md:22-23** — The §7.1 criteria of PRD-0001 stay as they are. (O-115) *Effect:* The plan cannot edit the §7.1 criteria to fit phase 1. What phase 1 is accepted against needs an Operator decision (gov-pilot-definition).
22. **docs/spec/README.md:9-12** — It is written one milestone at a time (task T-0drh, #74; decision O-114). This version covers phase 1 of PRD-0001 §9 only. A later milestone adds its own sections. *Effect:* Each later milestone starts with a specification task before its build issues.
23. **docs/adr/0012-build-layup-in-bootstrap-mode.md:60** — The documentation issues #15, #21, #24, #34 and #48 stay open as documentation fixes under part 1; they are batched into product tasks that touch the same files. *Effect:* The plan names, for each of these issues, the product task that touches its files, and adds no separate documentation tasks for them.
24. **docs/adr/0012-build-layup-in-bootstrap-mode.md:47** — #49 (the R11 count) waits for the end of bootstrap mode. *Effect:* #49 goes after the ADR that supersedes ADR-0012, not in phase 1.
25. **docs/spec/setup.md:235-239** — The engine's version of each such check passes and fails on the same fixtures as docs/setup/tests/run.sh (the defence that ADR-0011's Consequences name). *Effect:* A fix to setup-check.sh or its fixtures (#21, #61) changes what the Go verify checks must do, so it goes with or before the verify item that ports the same check.
26. **docs/spec/setup.md:276-280** — One directory per stack … internal/catalog/<stack>/ … The Go entry itself is the work of #29. *Effect:* Phase 1 has only the Go catalog entry, so the first pilot's problem statement must have a Go stack.
27. **docs/architecture.md:324-328; docs/spec/setup.md:29** — a private repository on GitHub Free has neither, and on GitHub Pro no drafts: the Operator makes it public or moves it to a plan that has both, K15 / Start 3: read back … visibility; the plan check … | no | layup run, phase 2 *Effect:* Phase 1 does not check the plan. The pilot target must enforce rulesets, or the S13 ruleset, on which REQ-007's merge control rests, does not hold.
28. **docs/architecture.md:1515-1516** — L-A1 Close when: before LAYUP runs on a target that is not the Operator's; each session and each gate run gets a container or another user. *Effect:* The first pilot's target must be the Operator's, or the plan schedules the isolation work before the pilot.
29. **docs/architecture.md:1546-1547** — L-B1. The toolchains on the host. layup gate needs each target stack's toolchain on the LAYUP host. *Effect:* The pilot inputs include the Go toolchain on the host (go1.27.1 today).
30. **docs/tests/test-levels.md:63-81; docs/guardrails.md:276-280** — Command: go test -tags=e2e ./… Timeout: -timeout 10m … Where: CI only … UAT is a human layer on top of E2E. / A test that reads … a live network … passes or fails for reasons unrelated to the code. *Effect:* The pilot run on GitHub is a UAT with evidence under runs/. Each user-facing phase-1 path also needs an automated E2E test against a local fixture baseline and a local target.
31. **docs/tests/dod-checklist.md:24-29, 80-81; docs/tests/traceability-template.md:40-41; docs/setup/setup-check.sh:305** — Every requirement (REQ/NFR) and every DoD item maps to at least one traceability row whose status is green or frozen. … The traceability table's ID set matches the requirement set. / Level — one of unit, integration, e2e, or uat. / MK_EXEMPT … ^docs/tests/traceability-template\.md$ *Effect:* Phase-1 build tasks cannot tick this box until a table exists with a row for every REQ and NFR. The copy carries no &lsaquo;…&rsaquo; cell, and the plan says how a discipline check gets a row.
32. **docs/prd/PRD-0001-layup.md:138, 241; .github/workflows/ci.yml:294-303** — go build ./... succeeds with no module outside the standard library (go list -deps ./... names none) / NFR-007 … T-mtb9 | CI job lint (gofmt, go vet) *Effect:* The NFR-007 Test cell names a test that does not check its criterion. A phase-1 item must own a test that runs go list -deps, and the §12 cells change with it.
33. **docs/engineering-discipline.md:602-615** — Every abbreviation that appears in any conversation, context, prompt, reply, or response must have an entry in glossary.md. *Effect:* T-4wrw adds a PDR row. T-55n2 adds a row for each term it gives to LAYUP's own plan that the glossary lacks (for example pilot, milestone, size class).
34. **docs/tasks/backlog.md:8-22** — Each task has a stable ID assigned once and never reused … confirm tasks/<id>.md does not already exist. One line per task. *Effect:* New task IDs for the plan's issues are random and checked. Each open task gets one backlog line.
35. **docs/engineering-discipline.md:49-53; docs/spec/README.md:46-47, 61-63** — If an operation or loop can run longer than 10 seconds, give an explicit, lightweight progress indicator. *Effect:* The issues for layup setup (the clone of the baseline) and layup gate (test runs) need a progress line on standard error, because standard output must stay deterministic.
36. **docs/architecture.md:739-743** — A plan task splits the milestone into build tasks, each with the requirement IDs it serves, the tests that will show it done, a size class, and its predecessors; … the predecessors form no cycle. *Effect:* This form binds a target's plan, not LAYUP's own. If the plan uses it, the plan can be checked the same way: every phase-1 Must ID has a task, and the predecessors form no cycle.
37. **docs/engineering-discipline.md:225-227; docs/setup/open-gaps.tsv:1** — A panel is convened only for an ADR that changes the product architecture. / &lsaquo;the domains of a panel&rsaquo; *Effect:* T-55n2, T-4wrw and the ADR that supersedes ADR-0012 need no panel. A phase-1 issue that needs a new product-architecture ADR needs a panel, whose domains are still an open gap.
38. **docs/ci/README.md:62-72, 88-95; docs/setup/setup-check.sh:553** — twelve checks, one per workflow job … To write or read the setting needs an administration-scoped token … run by an operator. *Effect:* A new CI job in phase 1 needs a change of branch-protection.json in the same PR and an Operator step (gov-operator-ci-protection).

### Risks of area `gov`

- The word 'pilot' has two meanings. ADR-0012:49 means one target set up and its gate run from outside, after phase 1. PRD:144-146 and PRD:182 mean the pilot that measures the baseline, and the phase-4 pilot on two problem statements. If no one decides (gov-pilot-definition), bootstrap mode ends either too early or not until phase 4.
- Several phase-1 §7.1 criteria cannot pass after phase 1 alone. REQ-001 needs the review of meaning, REQ-004 needs four active kinds, REQ-007 needs PRs on the target, REQ-009 and REQ-011 need a writer, and NFR-002 needs a session to continue a task. O-115 keeps those criteria as they are, so the idea owner's acceptance of phase 1 has no clear test without a new decision.
- No source defines how the pilot's defect and stall numbers are counted. If the rule is chosen after the pilot, it is a fitted parameter (guardrails.md:21-26).
- If the plan slices the work very finely, the gate cost grows: ADR-0012:27 measured at least four model sessions for each small task. If it slices too coarsely, the R11 count can reject plans many times, which is the deadlock that ADR-0012 records. A verify issue with 14 checks (18 rows) is the likely case.
- 'Is one task' (ADR-0012:48) and R11 slicing can conflict. A plan review can reject the plan on this alone, and the Operator must decide.
- The engine checks can drift from setup-check.sh. If #21 or #61 changes setup-check.sh after the Go port, the shared fixtures differ (setup.md:238-239).
- The pilot inputs are large: five prose files of a target and every answer with a comment-URL source. The manual setup of the same rows in LAYUP took four tasks and about 195 minutes (record-T-n1hp.md:37-50).
- L-A6: the pilot pins the baseline's latest commit. If the baseline changed after LAYUP's pin, a setup step can break, and the pilot stops until a LAYUP change follows.
- A private pilot target on GitHub Free does not enforce rulesets, and phase 1 does not check the plan. A failed gate job could then merge.
- Forge work (opening and re-scoping issues, approvals, closing #24 and #33) depends on the layup-agent App token helper. If it fails, an agent writes nothing to GitHub, and the plan waits for the Operator.
- The first pilot needs network access (git ls-remote and clone of the baseline, GitHub). External harness runs failed often (ADR-0012:28).
- L-A1: layup gate runs the target's own code on the shared host, under the Operator's user. A pilot target that is not the Operator's is not allowed until isolation exists.
- At session start, the Operator's main checkout had an uncommitted change to docs/spec/packages.md. This read is at 7cdd346 and does not see that change, so the package rows in the plan can be out of date.

### Critic points that the reviser of `gov` rejected

- Correction to the tests of gov-first-pilot (part of the low-confidence blanket correction on uat entries): an approval or a decision is evidence, not a test level, so remove the uat entries. *Basis:* I applied it to the approvals of gov-prd-task-column, gov-rescope-29, gov-rescope-33, gov-pdr-approval, gov-supersede-adr-0012 and gov-operator-setup-o112. I kept uat for the pilot run and the idea owner's acceptance on it, because a UAT is 'the acceptance step that rides on the E2E path', judged by a person (docs/tests/test-levels.md:77-81).

## Area `later`

Phases 2 to 4 have no code today. `go list ./...` at 7cdd346 names only cmd/layup, internal/cli and internal/psb, and internal/cli/cli.go:31-40 dispatches only `version` and `psb`. docs/spec/ covers phase 1 only (docs/spec/README.md:9-12). The architecture designs all later work. docs/spec/packages.md:74-91 and docs/spec/records.md:23-88 give the later packages and records with their phases.  I checked each critic point against the sources and accepted most of them. The final inventory has 14 large milestones: seven for phase 2, four for phase 3 and three for phase 4.  Changes from the first inventory: (1) New milestone later-p2-scaffold. `layup run` drives the phase-1 step runner and reads exit code 3 as a wait. It brings the phase-2 changes of S01 to S04 and S15, Scaffold 6 (both rulesets, with the four layup/ checks) and Scaffold 7 (read-back, the GH013 probe, the activity read). Scaffold 7 moves here from run-start, because a probe can be refused only after the ruleset is applied. (2) New milestone later-p2-build-accept, which takes Build and Accept out of the phase-loop milestone. This breaks the cycle with rule-protection: a pending kind fails every product-path change (docs/architecture.md:451-457), so no build task can merge before the activation batch. The order of phase 2 is now: run-start, sessions-ledger, intake-spec, scaffold, phase-loop-handoffs, rule-protection, build-accept. (3) The basic layup/rules status moves into phase-loop-handoffs. The Scaffold 6 ruleset requires all four layup/ checks (docs/architecture.md:413-418, :546-550), so all four must be posted before the first merge, which is the specification task at the first bet. (4) The pre-push rule-path refusal moves into sessions-ledger (ADR-0017 decision 2). (5) The full `layup spec check` moves into phase-loop-handoffs. intake-spec keeps `--facts` only. (6) later-p4-measures-learning is split into later-p4-measures and later-p4-learning. The learning loop has no PRD requirement and no In-Scope fact. (7) In phase 3, the dependency order is now smartif, then escalation, then clarification and stalls-budget. Escalation writes the questions.tsv rows of its escalations itself. stalls-budget applies changes of the band. (8) e2e tests run on a local fake forge (httptest) with scripted harness programs. CI job `tests` has contents: read and no secret, and a 10-minute timeout. Every run on a real forge, harness or provider, and every human step, is UAT.  The architecture does not give the build order inside a phase. It also leaves some phase assignments unsettled: questions.tsv, routing.tsv, the budget check, the frozen test list, the dead-man job, the Shape panel and the retrospective. Each is an open question with file:line. A depends_on key without the later- prefix is my best guess at a key of another area (psb-, gate-, setup-, verify-, found-).

| Key | Title | Requirements | Size | After | Host |
| --- | ----- | ------------ | ---- | ----- | ---- |
| `later-p2-run-start` | Phase 2: layup run Start and restart, the records writer, the lease and the forge adapter | NFR-001, NFR-002, NFR-006, REQ-002 | large | found-tsv, found-git, found-cli, found-record-schemas, setup-runner | `M2a` |
| `later-p2-sessions-ledger` | Phase 2: role sessions, the harness probe and admission, the pre-push checks, and the ledger writer | REQ-011, REQ-005, REQ-003, REQ-013, NFR-005, NFR-001 | large | later-p2-run-start, found-record-schemas, setup-rule-path-register | `M2b` |
| `later-p2-intake-spec` | Phase 2: the Intake gap check in one batch, and the fact spans | REQ-001, REQ-012, NFR-003 | large | later-p2-run-start, later-p2-sessions-ledger, psb-check | `M2c` |
| `later-p2-scaffold` | Phase 2: layup run drives the setup, the rulesets and the forge checks (Scaffold) | REQ-002, NFR-001, NFR-002, NFR-003, NFR-006, REQ-003 | large | later-p2-run-start, later-p2-sessions-ledger, later-p2-intake-spec, setup-runner, setup-records-branch, verify-checks | `M2d` |
| `later-p2-phase-loop-handoffs` | Phase 2: the task loop, typed handoffs, Shape and the first bet | REQ-005, REQ-007, REQ-012, REQ-003, NFR-004, NFR-001 | large | later-p2-scaffold, later-p2-sessions-ledger, later-p2-intake-spec, gate-command | `M2e` |
| `later-p2-rule-protection` | Phase 2: rule batches, gate activation, batch approval and layup audit | REQ-003, REQ-004, NFR-001, NFR-004 | large | later-p2-phase-loop-handoffs, later-p2-scaffold, gate-command, setup-catalog-go | `M2f` |
| `later-p2-build-accept` | Phase 2: the milestone plan, build tasks with merge order, and Accept | REQ-005, REQ-007, REQ-004, REQ-012, NFR-004 | large | later-p2-rule-protection, later-p2-phase-loop-handoffs, setup-catalog-go | `M2g` |
| `later-p3-smartif` | Phase 3: the smart-if provider client and decisions.tsv | NFR-005, REQ-006, REQ-008 | large | later-p2-run-start, later-p2-phase-loop-handoffs | `M3a` |
| `later-p3-escalation` | Phase 3: the escalation screen, human-input accounting and the audit | REQ-008, NFR-001 | large | later-p3-smartif, later-p2-build-accept, later-p2-intake-spec | `M3b` |
| `later-p3-clarification` | Phase 3: autonomous clarification by the owner role | REQ-006 | large | later-p3-escalation, later-p3-smartif, later-p2-build-accept, later-p2-sessions-ledger | `M3c` |
| `later-p3-stalls-budget` | Phase 3: stalls, the panel, the circuit breaker, the budget and the dead-man job | REQ-010, REQ-009, REQ-011 | large | later-p3-escalation, later-p3-smartif, later-p2-build-accept, later-p2-sessions-ledger, found-record-schemas | `M3d` |
| `later-p4-neutrality` | Phase 4: a second harness agent does and checks the work, and a target continues without LAYUP | REQ-013, NFR-002 | large | later-p2-build-accept, later-p2-sessions-ledger, later-p3-stalls-budget | `M4a` |
| `later-p4-measures` | Phase 4: layup report, the pilot baseline and the start values | REQ-011, REQ-005, REQ-008, REQ-006, REQ-009, REQ-010, REQ-001 | large | later-p3-escalation, later-p3-clarification, later-p3-stalls-budget, later-p2-build-accept | `M4b` |
| `later-p4-learning` | Phase 4: the retrospective, layup learn and the lessons | — | large | later-p4-measures, later-p2-rule-protection, later-p3-stalls-budget | a decision of the Operator |
| `later-p4-pilot` | Phase 4: the pilot on two problem statements with different stacks | REQ-014, REQ-001, REQ-002, REQ-003, REQ-004, REQ-005, REQ-006, REQ-007, REQ-008, REQ-009, REQ-011, REQ-012, NFR-001, NFR-002, NFR-003, REQ-015, REQ-016, REQ-017, REQ-018 | large | later-p4-neutrality, later-p4-measures, later-p2-rule-protection, later-p3-escalation, setup-catalog-go | `M4c` |

### `later-p2-run-start` — Phase 2: layup run Start and restart, the records writer, the lease and the forge adapter

- **What:** `layup run --new OWNER/NAME ...` does Start: it resolves and clones the baseline, prints the root push, reads back the default branch, the root tree, the visibility and the plan, and stops when the plan or one of the six forge capabilities is missing. It pushes the first records commit (the pin, the briefs with their SHA-256, approvers.tsv, the lease row, start.tsv with the LAYUP version and each harness's cap and wall-clock limit read from host:registers/harnesses.tsv), opens the Intake and control issues, holds the lease with heartbeat and fencing, and copies each comment before it acts on it. `layup run TARGET` restarts: it rebuilds its clones from the forge and the records, takes the lease over, and stops when its own version differs from the recorded one. First demo: on an empty GitHub repository, after the Operator's printed root push, the records branch shows the first commit by layup-agent[bot], and the Intake and control issues exist. Before the build, docs/spec must add: the NFR-001 one-writer section for layup run, a forge-interface section, the schemas of start.tsv, approvers.tsv, lease.tsv, copies.tsv and host:registers/harnesses.tsv, how layup run gets the App ID, the key file and harness credentials under the 'Arguments' rule, and the package rows of internal/run, internal/records, internal/forge and internal/forge/github.
- **Specification:** docs/spec/setup.md: The boundary of phase 1 (rows Start 2; Start 3: read-back, plan check, six capabilities; first records commit, approvers.tsv, lease; Intake and control issues); docs/spec/records.md: NFR-001 — Git is the system of record (item 1: from Start on, layup run is the one writer); docs/spec/records.md: The layout of the records branch (rows start/, start.tsv, approvers.tsv, lease.tsv, copies.tsv; In other places: host:registers/harnesses.tsv); docs/spec/packages.md: The packages of later phases (internal/run, internal/records, internal/forge, internal/forge/github); docs/spec/packages.md: NFR-007 (rule 5: net only in the forge adapter); docs/spec/README.md: Commands (Arguments; Exit codes) and Records (The schema block); docs/architecture.md: 1. LAYUP and a target (the LAYUP host; the forge and its six capabilities); docs/architecture.md: 2. The components (one run per target, the lease, fencing); docs/architecture.md: 3. Records and identities (one writer; identities; a human decision; copy before read); docs/architecture.md: 5. Intake and setup (Start 1 to 3)
- **Packages:** internal/run, internal/records, internal/forge, internal/forge/github, internal/route, internal/git, internal/cli, internal/tsv
- **Present state:** None in code. internal/cli/cli.go:31-40 dispatches only version and psb, and `go list ./...` names only cmd/layup, internal/cli and internal/psb. The design is at docs/architecture.md:295-337 (Start), :77-97 (lease and fencing) and :115-182 (one writer, identities, copy). The records rows are 'later' at docs/spec/records.md:26-35. The package names are at docs/spec/packages.md:76-79.
- **Tests:**
  - unit: With a stand-in clock and a stand-in git, the lease is taken over only after the heartbeat counter has not moved for 3 x lease.H by the run's own clock, and a refused records push stops the run until it reads the lease again.
  - unit: A human decision is a comment whose author ID is in approvers.tsv and whose App field is empty. A review, a review comment, a commit, a reaction or an edit is never a decision. Each rule takes only its role: Intake answers only from the Operator or the idea owner, an acceptance only from the idea owner.
  - unit: Copy-before-read writes the body, the author ID and login, the comment ID, the App field, the time and the SHA-256. An edit is a new input event, and the first copy stays.
  - unit: internal/records refuses a table whose header does not match its record kind. For each new record kind, a schema-block test compares the spec block with what the code writes (docs/spec/README.md:75-77).
  - integration: The GitHub adapter, against an httptest server on loopback, makes an App JWT and an installation token and maps each of the six capabilities. A server that lacks one makes Start stop and name the capability.
  - integration: Start, against the httptest forge and a real local bare repository with the real internal/git, writes the first records commit with the pin, the briefs byte for byte, approvers.tsv, start.tsv and the lease row. A restart with a different LAYUP version stops.
  - e2e: In CI with no secret: `layup run --new` and then `layup run TARGET` (restart), against a local fake forge and local bare repositories, from the command line to the records branch.
  - uat: On a real GitHub test target, the Operator runs Start and the printed root push, sees each LAYUP write as layup-agent[bot], and reads the records branch with a plain git clone and no tool.
- **External inputs:**
  - The GitHub App layup-agent with the commit-statuses permission, installed on a test target, and its private key on the host in mode 0600 (docs/architecture.md:139-144; the Operator).
  - A GitHub plan that enforces rulesets and gives draft pull requests (a public repository, or a paid plan; docs/architecture.md:322-328, K15).
  - The GitHub REST and GraphQL documentation for the App JWT, the installation token, issues and comments with performed_via_github_app, at the date of use.
  - The Operator's harness register rows with harness.<id>.cap and harness.<id>.wall, set before Start (docs/architecture.md:1108), and the Start values --plan, --intake-cap, --lease-h, --watch-t and the two logins (docs/architecture.md:297-304).
- **Open questions:**
  1. How does layup run get the App ID, the private-key file and the harness credentials? docs/spec/README.md:44-45 says 'It reads no environment variable for an input', but docs/architecture.md:241-243 gives a session 'the harness's own credential: a variable, or a file'. The Start command at :297-300 names no flag for the App.
  2. Which phase delivers the dead-man job's first notice? docs/spec/setup.md:31 puts 'the dead-man job's first notice' under '`layup run`, phase 2'. docs/spec/records.md:88 gives the dead-man job's record phase 3. This plan builds the job in later-p3-stalls-budget. Until then, Start records 'watch not confirmed' (docs/architecture.md:335-337).

### `later-p2-sessions-ledger` — Phase 2: role sessions, the harness probe and admission, the pre-push checks, and the ledger writer

- **What:** layup run starts a role session in a session directory: a --no-local clone at the base commit, an empty home, a prompt file built by code, and an environment from a named list with no forge credential. It refuses a rule file above the session directory, kills the session at harness.<id>.wall, and records the cap and the wall-clock limit in the start row. It reads the typed result and refuses it when its attempt is no longer open. It fetches the branch by SHA with hooks off, refuses a workflow change, and refuses a change to a rule path (from the rule-path register of the base) before any push. It pushes with the App token, binds the SHA to the session only after the forge accepts the push, and starts each comment for a session with the session ID. It probes each registered harness, admits pairs by code in the order of the routing register, and writes one sessions.tsv row and one telemetry.tsv row per session (the REQ-011 writer). First demo: a probe session on each registered harness, then one developer session whose commit lands on task/<task>/<attempt>, with one complete telemetry row. Before the build, docs/spec must add: a session-contract section (ADR-0015 decisions 4 and 5, with the REQ-013 parts 'rules only from the target' and the probe of the loaded instruction files), the schemas of sessions.tsv, harnesses.tsv, routing.tsv (order only), tasks/<task>/results/<session>.tsv and payloads/, and the package rows of internal/session, internal/route, internal/ledger and the pre-push part of internal/rules.
- **Specification:** docs/spec/records.md: REQ-011 — The telemetry record (Not in phase 1: the writer); docs/spec/records.md: The layout of the records branch (rows sessions.tsv, harnesses.tsv, routing.tsv, tasks/<task>/results/<session>.tsv, payloads/<sha256>; In other places: host:registers/harnesses.tsv, host:prices.tsv); docs/spec/setup.md: The rule-path register (read before a push); docs/spec/packages.md: The packages of later phases (internal/session, internal/route, internal/ledger, internal/rules); docs/architecture.md: 3. Records and identities (one writer: the open-attempt check; role sessions); docs/architecture.md: 4. Role sessions and model calls; docs/architecture.md: 6. Native gates and rule protection (Prevention layer 2); docs/architecture.md: 9. Squads and routing (the harness register, the probe, admission, the context of a session); docs/architecture.md: 12. Cost, budget and the measures > The ledger
- **Packages:** internal/session, internal/route, internal/ledger, internal/rules, internal/records, internal/run, internal/git
- **Present state:** None. The design is at docs/architecture.md:224-285 (sessions), :540-545 (pre-push), :899-935 (register, routing, context) and :1251-1271 (ledger). The schemas of telemetry.tsv and prices.tsv exist at docs/spec/records.md:150-183, with no writer. The rule-path register schema exists at docs/spec/setup.md:163-183, with no reader.
- **Tests:**
  - unit: The session environment holds only the named list (PATH, locale, HOME, GIT_CONFIG_NOSYSTEM, the harness credential), with no GH_TOKEN and no SSH agent socket.
  - unit: Telemetry rows use '—' and the status column for an unknown value, never 0. tokens_status is 'unavailable' when the harness gives no count, and money_status 'computed' cites a prices.tsv row.
  - unit: Admission: a model on the 'not used' list, or an unprobed harness, is not admitted. A new harness version is probed again.
  - unit: With a stand-in events table, a result whose attempt was closed, replaced or rebased by a later event is refused. The attempt and the base come from the start row, never from the result.
  - unit: The pre-push check refuses a change to a register path (a directory ending in '/', each .sh file), and allows added lines in section 2 of docs/guardrails.md.
  - integration: In a real temporary directory tree, the start is refused when a rule-file name of the harness is in any directory from the session directory up to the root.
  - integration: With a fake harness program, the real git and a local bare repository: the branch is fetched by SHA with hooks off and must descend from the base. A diff that touches .github/workflows/ or a rule path is refused before any push, and the diff goes to payloads/. The SHA is bound only after the push is accepted.
  - integration: A fake harness that runs past harness.<id>.wall is killed, and the start row records the cap and the wall-clock limit. A start whose estimated prompt size (bytes / 4) passes the model's context size is refused, with both numbers in the start row.
  - e2e: In CI with no secret, on a local fake forge with a scripted harness program: the probe, one session, the App push, a comment that starts with the session ID, and a complete ledger row (W-12 steps 2-3 and 10; W-10 step 4).
  - uat: With the Operator's credentials, one real session on each registered harness: the probe, the push, and a telemetry row whose tokens are 'observed'.
- **External inputs:**
  - The harness CLIs at a pinned version, with documentation for the permission flag, the model flag, the output format and the usage report (claude, devin, opencode).
  - The harness credentials on the host.
  - A prices.tsv with source URLs for each model used (docs/spec/records.md:174-183).
- **Open questions:**
  1. Where do routes live in phase 2? docs/spec/packages.md:89 gives internal/route ('the harness register, the probe, admission and routing') phase 2, and docs/architecture.md:908 says 'Routing is a register on the records branch'. But docs/spec/records.md:53 gives routing.tsv ('the routing register and its weights') phase 3. This plan writes routing.tsv with the order only in phase 2.

### `later-p2-intake-spec` — Phase 2: the Intake gap check in one batch, and the fact spans

- **What:** In Intake, sessions do the review of meaning, the numbering into byte-offset spans, the draft and a completeness review on another harness. `layup spec check --facts` checks the spans. Code merges these rows with the psb rule gaps and the setup questions into one Intake comment, copies the answers, writes them as a raw fact, and asks one follow-up for each missing or open answer. layup run commits the gap table and the span check table. First demo: on a pilot problem statement, one Intake comment holds the rule gaps and the gaps of meaning, and `layup spec check --facts` passes on the numbered spans. Before the build, docs/spec must add: the rest of REQ-001 (the review of meaning, the one batch), the Intake part of REQ-012, the schemas of intake/meaning.tsv and spec/spans.tsv, the `layup spec check --facts` table, the records paths of the committed gap and check tables, and the internal/spec package row.
- **Specification:** docs/spec/psb-check.md: REQ-001 > Not in phase 1 (the review of meaning; one batch on the Intake issue); docs/spec/setup.md: The boundary of phase 1 (rows Gap check 2 to 4; Gap check 5 and 6); docs/spec/records.md: The layout of the records branch (rows intake/meaning.tsv, spec/spans.tsv, questions.tsv); docs/spec/packages.md: The packages of later phases (internal/spec); docs/architecture.md: 5. Intake and setup (The gap check, in one batch); docs/architecture.md: 7. The specification (steps 1 to 5)
- **Packages:** internal/spec, internal/psb, internal/run, internal/session, internal/forge, internal/records
- **Present state:** Only the psb rule gaps exist: internal/psb (layup psb check, the REQ-001 rule part). There is no spec check and no session. The design is at docs/architecture.md:339-375 and :580-640.
- **Tests:**
  - unit: The span check fails when a non-whitespace byte is in zero spans or two, when a span's text differs from the file at its offsets, when a fact ID repeats, or when a class is not one of the four.
  - unit: The batch merge drops a question whose quote and kind repeat, and gives each question an ID. Each question carries its byte-exact quote.
  - integration: With scripted harness programs, the real git and the httptest forge: an Intake run that fails the span check posts no batch and stops. With one admitted harness, the completeness review does not run, and layup run records 'not run' and stops.
  - e2e: In CI on a local fake forge, with scripted sessions and scripted answer comments from stand-in approver IDs: one Intake comment, the answers copied and written as a raw fact, and a follow-up only for missing or open answers (W-01 steps 2-11; W-11 steps 1-5).
  - uat: On a pilot problem statement, the idea owner reads the one batch and finds no question asked twice and no gap of meaning missing.
- **External inputs:**
  - At least two admitted harnesses, for the completeness review (docs/architecture.md:606-608).
  - A pilot problem statement chosen by the idea owner (docs/prd/PRD-0001-layup.md:203).
- **Open questions:**
  1. REQ-013 (two harnesses) is phase 4 (docs/prd/PRD-0001-layup.md:182), but Intake already needs a second harness: docs/architecture.md:606-608 says that with only one admitted harness the completeness review 'does not run, and `layup run` records "not run" and stops'.
  2. Which phase holds questions.tsv? docs/architecture.md:373-375 and W-01 step 6 (docs/walkthroughs/W-01-problem-statement-quality.md:21) write each Intake question as a row of the project question table (phase 2). docs/spec/records.md:37 gives questions.tsv phase 3.

### `later-p2-scaffold` — Phase 2: layup run drives the setup, the rulesets and the forge checks (Scaffold)

- **What:** layup run writes WORK/inputs/answers.tsv and inputs/files from the copied Intake answers and the session results, runs the phase-1 step runner, reads exit code 3 as 'wait for the answers', and commits setup/record.tsv, setup/verify.tsv and rule-paths.tsv itself. Role sessions write the prose rows S07, S08, S09 and S14. It writes the rulesets of both branches: on the default branch the native jobs, the four layup/ checks pinned to the App, the 'up to date' rule and an empty bypass list; on the records branch, updates only by the App. It prints the Operator's push and apply commands and records the apply output. Then Scaffold 7: the pushed tree equals the verified one, the rules read back, the probe push refused with GH013, and the activity read, all in setup/forge-check.tsv; each later start reads the rules back and probes again. First demo: after the Operator's printed pushes and ruleset apply, the push to layup-probe is refused (GH013), and setup/forge-check.tsv shows each rule and probe as passed. Before the build, docs/spec must add: the phase-2 changes of the step table (S01 no longer asks the name, the visibility or the baseline; S02 and S03 move to Start 2; S04 reads the pin from the records; S15 no longer makes the orphan commit; a new writer of the records README), a record of the verified tree and the setup-commit SHAs, the schema of setup/forge-check.tsv, the setup-record line that tells the Operator to remove the layup/ checks when LAYUP is absent, and how layup run adopts a target that phase 1 set up.
- **Specification:** docs/spec/setup.md: The boundary of phase 1 (rows Scaffold 2, Scaffold 6, Scaffold 7); docs/spec/setup.md: The steps (S01 to S04, S13, S15); docs/spec/setup.md: Where the records go in phase 1; docs/spec/setup.md: The checks of layup setup verify (row 'the rulesets read back (S13); the probes of §3'); docs/spec/setup.md: Not in phase 1 (the ruleset of the records branch; the layup/ required checks); docs/spec/records.md: NFR-001 — Git is the system of record > Not in phase 1 (the one-writer ruleset); docs/spec/records.md: NFR-002 — A target is independent of LAYUP (item 3); docs/spec/records.md: The layout of the records branch (rows README.md, setup/rulesets.txt, setup/forge-check.tsv); docs/spec/README.md: Commands (exit code 3); docs/architecture.md: 3. Records and identities (What the forge enforces, and what LAYUP checks); docs/architecture.md: 5. Intake and setup (Scaffold 1 to 7); docs/architecture.md: 6. Native gates and rule protection (Prevention layer 3; When LAYUP is absent)
- **Packages:** internal/run, internal/setup, internal/verify, internal/session, internal/forge, internal/forge/github, internal/records, internal/git
- **Present state:** The phase-1 step runner and the records commit of S15 are phase-1 work of #29 (docs/spec/setup.md:77-104, :201-224). No code drives them. The design is at docs/architecture.md:166-199 and :377-423.
- **Tests:**
  - unit: The exit-code map of layup setup: 0 goes on, 3 waits for the listed answer IDs, and 1 or 2 stops the run.
  - unit: The default-branch ruleset file names the native jobs pinned to GitHub Actions, the four layup/ checks pinned to the App's integration_id, the 'up to date' rule and an empty bypass list. The records-branch ruleset has the App as its only bypass actor.
  - integration: With a real local bare repository whose pre-receive hook refuses with the GH013 text, and the real internal/git: a GH013 refusal counts as 'refused'. Any other failure is 'probe not run', never a pass.
  - integration: layup run with the real step runner on a fixture baseline: answers.tsv is written from copied answers, exit 3 waits, a later run finishes, and layup run commits setup/record.tsv, setup/verify.tsv and rule-paths.tsv.
  - e2e: In CI on a local fake forge and local bare repositories: Start, scripted Intake answers, the setup, the printed commands run by the test in place of the Operator, and the Scaffold 7 checks (W-02 steps 1-7).
  - uat: On a real GitHub test target, the Operator runs the printed pushes and the ruleset apply, and sees the probe refused and the activity read work.
- **External inputs:**
  - The App's integration_id, and the Operator's admin login for the ruleset apply (the App has no administration permission, docs/architecture.md:139-142).
  - The GitHub documentation for rulesets, GET rules/branches, the GH013 refusal and the repository activity API, at the date of use.
- **Open questions:**
  1. How does layup run adopt a target that phase 1 set up? docs/spec/setup.md:309-311 says the Operator applies the records ruleset and the layup/ checks 'with the commands that `layup run` prints when it first starts on the target (phase 2)'. But docs/architecture.md:297-300 starts only from 'an empty repository' with `layup run --new`, and docs/spec/records.md:99-103 says that the phase-1 records branch already has the Operator's first commit.
  2. When is the records-branch ruleset applied? docs/spec/records.md:112-113 says 'applied at Start (phase 2)'. docs/architecture.md:413-414 applies 'the two rulesets of §6' at Scaffold 6.

### `later-p2-phase-loop-handoffs` — Phase 2: the task loop, typed handoffs, Shape and the first bet

- **What:** Each pull request is a task: an issue, a plan and a plan review on another harness, a developer session, a draft pull request whose body links the records commit it started from, the statuses layup/gates, layup/spec and layup/rules (the basic rule: a rule-path change fails unless it is a rendered record equal to code's rendering), a verifier on a non-author harness that sets layup/verify, the close-out commit, and the merge at the head SHA. Each handoff is a typed result that code checks against the transition table. Shape filters the architecture options by code. The Bet brief is checked by code, and the bet comment is copied; at the first bet, code adds the rendered commit and writes the confirmed inventory version. The full `layup spec check`, the expected prd-lint failure of a specification that waits for its bet, and the parameter change `set <name> <value> because <reason>` on the control issue are also in this milestone. First demo: the specification task merges after the first bet: each handoff is valid, the rendered commit holds the MoSCoW and Phase columns, and the four layup/ checks are green. Before the build, docs/spec must add: REQ-005, the status part of REQ-007, the rest of REQ-012 (the full layup spec check, its table and records path), the schemas of tasks.tsv, transitions.tsv, tasks/<task>/events.tsv, tasks/<task>/gates.tsv with the head in the key, spec/inventory-<n>.tsv and parameters.tsv, and the package rows of internal/handoff, internal/spec and the layup/rules part of internal/rules, with the 'starts a program' mark for the package that runs the pinned prd-lint.
- **Specification:** docs/spec/gate.md: REQ-004 > Not in phase 1 (the commit status layup/gates); docs/spec/gate.md: REQ-007 — Each change passes the gates before review or merge (Not in phase 1); docs/spec/records.md: The layout of the records branch (rows tasks.tsv, transitions.tsv, tasks/<task>/events.tsv, tasks/<task>/gates.tsv, spec/inventory-<n>.tsv, parameters.tsv); docs/spec/packages.md: The packages of later phases (internal/handoff, internal/spec, internal/rules, internal/run); docs/spec/packages.md: NFR-007 (rule 4: starts a program); docs/architecture.md: 7. The specification (steps 6 to 8; layup spec check); docs/architecture.md: 8. The phase loop (Shape; Bet; The task loop; Handoffs); docs/architecture.md: 6. Native gates and rule protection (layup/rules; layup gate and its status table); docs/architecture.md: 10. Decisions > The parameters
- **Packages:** internal/run, internal/handoff, internal/spec, internal/rules, internal/gate, internal/forge, internal/forge/github, internal/records, internal/session
- **Present state:** None. layup gate is specified (docs/spec/gate.md) but not built: `go list ./...` at 7cdd346 names no internal/gate, and internal/cli/cli.go:31-40 has no gate command. The design is at docs/architecture.md:684-869, with the transition table at :826-839.
- **Tests:**
  - unit: A handoff is valid only when the schema is right, each named artifact is in the commits with its hash, and each required item is true. A pair with no table row is invalid, and a field that an agent writes never makes it valid.
  - unit: layup spec check fails on a covers ID that does not resolve, an uncovered need or constraint, a measure with no success-criterion row, a requirement with no criterion, a MoSCoW or Phase value that differs from the latest copy, a bad inventory hash, an answer's fact without its quote's fact, a delivered requirement with no docs/spec section, a task with no requirement, and 'did not run'. Before the first bet, a task that does not touch the PRD or docs/spec is 'clear'.
  - unit: layup/gates is a success only when every kind passed or is 'clear'; 'not-active', or a pending kind with a product path, gives failure.
  - unit: The Shape filter drops an option with a failed row, a missing row or a pass without evidence. One survivor is chosen. When all are dropped, one rewrite is asked, and a second full drop records a round with no progress.
  - unit: The rendered commit may change only the MoSCoW and Phase columns, the non-goals and the numbered facts record. layup/rules passes a rendered record only when it equals code's rendering of the named copy.
  - unit: A 'set <name> <value> because <reason>' comment is refused for a wrong role, an unknown name or a value out of bound, and is applied at the next step boundary. A prd-lint head is accepted only when each failure line names an empty MoSCoW or Phase cell.
  - integration: With the httptest forge, fake harness programs and the real git: a gate failure goes back to the developer as a new attempt; the close-out commit is refused when it changes a file other than the task file and the completed log; the pull request body links the records commit; the merge is at the head SHA.
  - e2e: In CI on a local fake forge, with scripted harnesses and scripted approver comments: Shape, the first bet, and the specification task merged with the four layup/ checks green (W-05 steps 1-3; W-11 steps 6-9; W-07 steps 1-3 and 6-10 for a task with no product path).
  - uat: On a real test target, the idea owner reads a bet brief of five parts within the line limit and answers it with one comment.
- **External inputs:**
  - At least two admitted harnesses, so that the plan review and the verifier run outside the authors (docs/architecture.md:908-918).
  - The brief's line limit, set at Intake with evidence (docs/architecture.md:1107).
  - The Operator's confirmation of the reading of F-0003#58 ('reaches human review' = review requested or marked ready, docs/architecture.md:813-816).
- **Open questions:**
  1. How does a handoff name an artifact that lives only as an issue comment? docs/architecture.md:775 says layup run 'posts both on the issue', but :841-842 says 'A handoff is valid only when code finds ... each named artifact in the commits with its hash'. Issue #68 (.ctx/issues/68.md) records this gap and says 'The task that builds `layup handoff check` decides it, or a short ADR does'.
  2. Does phase-2 Shape run with no panel? docs/architecture.md:714-717 allows 'the brief then shows the survivors with no recommendation', but the panel of §11 is phase 3 (REQ-010).
  3. Where do the bet's caps live in phase 2? The bet brief has an appetite with 'its cap of money and wall-clock' (docs/architecture.md:722-723), but docs/spec/records.md:60-61 gives budget.tsv and milestones.tsv phase 3.
  4. What does the Retrospective phase (docs/architecture.md:689-690) do in phases 2 and 3? Lessons and layup learn are phase 4 (docs/spec/records.md:67-68), but rule batches are proposed 'at each retrospective' (docs/architecture.md:522-523).

### `later-p2-rule-protection` — Phase 2: rule batches, gate activation, batch approval and layup audit

- **What:** A rule batch goes on batch/<point> as a task. Before its approval, its own gate files run only in a scratch work tree, as evidence. The bet approves a batch by its head SHA and rule-file hash; code refuses the approval when the head has moved. After the approval, layup gate runs the batch's gate files (they must pass) and every recorded known-bad patch of each kind it touches (each must fail). layup/rules then passes the approved batch. At the first bet, the activation batch makes the pending kinds active. A batch that changes .github/workflows/ is pushed by the Operator and merged by its approver. Refused rule-path changes are collected as proposals for the next batch. `layup audit` lists each rule-path change with its merge actor from the activity API, and layup run commits its table. First demo: the activation batch of the first bet merges only after each known-bad patch fails its kind, and layup audit shows zero agent writes to rule paths. Before the build, docs/spec must add: REQ-003, the batch part of REQ-004, the schema of batches.tsv, the layup audit table and its records path, and the package rows of the batch part of internal/rules and of internal/audit.
- **Specification:** docs/spec/gate.md: REQ-004 > Not in phase 1 (a rule batch's own gate files; the run with each recorded known-bad patch and the approval hash); docs/spec/setup.md: Not in phase 1 (the activation of the pending kinds at the first bet); docs/spec/records.md: NFR-001 > Not in phase 1 (layup audit allows the Operator's phase-1 records commit by its SHA); docs/spec/records.md: The layout of the records branch (rows batches.tsv, payloads/<sha256>, rule-paths.tsv); docs/spec/packages.md: The packages of later phases (internal/rules, internal/audit); docs/architecture.md: 6. Native gates and rule protection (layup gate; Activation; Rule protection; Detection); docs/architecture.md: 3. Records and identities (The actor on the forge); docs/architecture.md: 8. The phase loop (Bet: approval of the verified batches; merge order after a bet)
- **Packages:** internal/rules, internal/audit, internal/gate, internal/run, internal/forge, internal/forge/github, internal/records
- **Present state:** The rule-path register schema exists at docs/spec/setup.md:163-183 (written by layup setup in phase 1, with no reader). There is no code. The design is at docs/architecture.md:476-578 and ADR-0017 decisions 1-5.
- **Tests:**
  - unit: An approval is refused when the batch head has moved since the request. The basic layup/rules check fails when the rule-file hash at the head differs from the hash recorded with the approval.
  - unit: With a stand-in activity list, layup audit marks a push to the default branch after setup as a bypass and an update with no actor as a failure. It allows the setup commits and the first records commit by their SHA, and it never reads a commit author.
  - integration: layup gate in batch mode, with the real git on a fixture Go target: before the approval, the scratch run gives evidence and posts no status. After the approval, the gate files pass and each recorded patch fails. A patch that passes, or that no longer applies and was not replaced, refuses the merge.
  - e2e: In CI on a local fake forge with scripted sessions and a scripted bet comment: the activation batch from Shape through its approval to its merge, and the audit table committed (W-03 steps 2-5; W-04 steps 1-5).
  - uat: On a real test target, the Operator pushes a workflow batch and the approver merges it. The Operator confirms the reading of F-0003#64 (docs/architecture.md:560-564). A set of known-bad commits is detected in full (the REQ-003 criterion).
- **External inputs:**
  - The Operator's push and merge of a batch that changes .github/workflows/ (the App has no workflows permission, docs/architecture.md:526-533).
  - The versions and documentation of the Go activation tools (golangci-lint depguard, docs/architecture.md:464-469), unless the stack catalog already pins them.
  - A coverage floor with evidence; otherwise the gap stays open (L-B2, docs/architecture.md:466-469).

### `later-p2-build-accept` — Phase 2: the milestone plan, build tasks with merge order, and Accept

- **What:** A milestone-plan task writes the task register. Code checks that each Must requirement has a task, each task has a requirement, and the predecessors form no cycle. A build task starts when its predecessors have merged, with at most build.parallel at once. Its layup/verify stays pending ('waits for its turn') until it is next in the merge order and up to date. After each merge, layup run merges the base into the next task's branch, and a conflict starts a new attempt. A verdict carries over a clean base merge that touches no file of the task; otherwise the task is closed and redone. The developer-to-verifier handoff needs the plan's tests to fail at the base and pass at the head, which layup gate runs. At Accept, the idea owner's 'accept' or 'reject' is copied to acceptance.tsv. First demo: one build task that changes a product path merges after the activation with every gate and the four layup/ checks green, and the idea owner's 'accept' is a row of acceptance.tsv. Before the build, docs/spec must add: the build-task part of REQ-005 and REQ-007, the task-register rules, the schema of acceptance.tsv, a named-test mode of layup gate and a per-stack named-test command in the catalog form, and the parameter build.parallel.
- **Specification:** docs/spec/gate.md: REQ-007 — Each change passes the gates before review or merge (Not in phase 1: no ready or merge while a selected kind did not run); docs/spec/gate.md: The command (a new named-test mode); docs/spec/setup.md: The stack catalog (a named-test command); docs/spec/records.md: The layout of the records branch (rows tasks.tsv, acceptance.tsv, tasks/<task>/tests.tsv); docs/architecture.md: 8. The phase loop (Build; Accept; the task loop step 6; Handoffs: developer (build task) to verifier, developer (milestone plan) to verifier); docs/architecture.md: 9. Squads and routing (admission of a verifier outside the authors)
- **Packages:** internal/run, internal/handoff, internal/gate, internal/route, internal/catalog, internal/records, internal/forge
- **Present state:** None. The design is at docs/architecture.md:740-768 (Build, Accept) and :804-811 (task loop step 6).
- **Tests:**
  - unit: The milestone-plan check finds a Must requirement with no task, a task with no requirement, a predecessor that is not in the register, and a cycle.
  - unit: The merge order is first ready, with ties by register row order. A waiting task does not hold up a ready one. No more than build.parallel tasks start at once.
  - unit: A verdict carries over a clean base merge that changes no file of the task. Otherwise the task is redone as a new task, which takes the old task's row and predecessor links.
  - unit: The author set of a verification has each harness bound in the ledger to a commit in the diff, each harness whose diff went into a prompt, and the writer of the frozen tests.
  - integration: layup gate in named-test mode, with the real git and go on a fixture Go target: the tests fail at the base and pass at the head, and a test that passes at the base refuses the handoff.
  - integration: With the httptest forge, fake harness programs and the real git: a verifier from an author harness is not admitted, and with only author harnesses the verification is 'not-active' and nothing merges. A base-merge conflict starts a new attempt.
  - e2e: In CI on a local fake forge, with scripted harnesses and scripted comments: a milestone plan, two build tasks that merge in order, a gate failure that goes back to the developer, and an accept row (W-04 steps 6-10; W-07; W-11 steps 10-12).
  - uat: On a real test target, the idea owner accepts a delivered requirement with one comment. A review finds that the transition table encodes the baseline's record kinds and fields (the REQ-005 criterion, docs/prd/PRD-0001-layup.md:118).
- **External inputs:**
  - build.parallel, set at Intake with the Operator's evidence (docs/architecture.md:1106).
  - At least two admitted harnesses, for a verifier outside the authors (docs/architecture.md:908-918).
- **Open questions:**
  1. Which phase holds the frozen test list? docs/architecture.md:776-777 says 'the list is frozen at the task's first valid handoff to the verifier', and the phase-2 transition row at :830 checks that 'the plan's tests fail at the base and pass at the head'. But docs/spec/records.md:47 gives tasks/<task>/tests.tsv phase 3.

### `later-p3-smartif` — Phase 3: the smart-if provider client and decisions.tsv

- **What:** The Operator registers providers on the host, a paid probe request runs before use, and the Operator chooses one provider per target, or none. A client asks only at the five named points P1 to P5, with the authority off, shadow, cautious or delegate and a threshold. It writes one decisions.tsv row per call. First demo: at P2 under shadow, a question goes to the registered provider, the answer is recorded with its model version, tokens and price, and the deterministic branch decides. Before the build, docs/spec must add: the smart-if part of NFR-005, the schemas of decisions.tsv, providers.tsv and host:registers/providers.tsv, and the internal/smartif package row, which is the second net exception.
- **Specification:** docs/spec/gate.md: NFR-005 — No model call in the engine checks (Not in phase 1: the smart-if client and decisions.tsv); docs/spec/packages.md: NFR-007 (rule 5: internal/smartif, the later exception); docs/spec/records.md: The layout of the records branch (rows decisions.tsv, providers.tsv; In other places: host:registers/providers.tsv); docs/architecture.md: 10. Decisions > The smart-if; docs/architecture.md: 4. Role sessions and model calls
- **Packages:** internal/smartif, internal/records, internal/run
- **Present state:** None. The design is at docs/architecture.md:945-1018 and ADR-0021 decisions 1-6. decisions.tsv is a 'later' row at docs/spec/records.md:56.
- **Tests:**
  - unit: A yes-or-no answer decides yes at p >= t and no at p <= 1 - t, and is undecided between. A choice decides only when its top option reaches t.
  - unit: Under off, shadow and cautious, a failure (an error, a time-out, a 429 after backoff, a state too large) takes the deterministic branch and is recorded. At delegate, a failure or an undecided answer goes to a human. cautious never leaves the safer branch, and is refused for P4.
  - unit: A call at a point outside P1 to P5 is refused. Each call writes exactly one decisions.tsv row, and unknown tokens or price are 'unknown', never 0.
  - integration: Against an httptest provider: the probe records the model version and the response form, and a new model version resets every delegate point to shadow.
  - integration: An import check reads the direct imports of each package (the boundary rule that reads the package table, docs/spec/packages.md:29-32): only internal/smartif and internal/forge/github import net, net/http or crypto/tls.
  - e2e: In CI with no secret, a scripted httptest provider: layup run makes one shadow call at P2 on a fake-forge target and records it in decisions.tsv (W-06 step 3).
  - uat: With the real chosen provider and its credential: a paid probe and one shadow call at P2 on the test target, recorded in decisions.tsv (W-08 steps 3-4).
- **External inputs:**
  - A provider and its credential: Jev (TypeSafe System One API) or Laya at an endpoint, with its API documentation and price source (docs/architecture.md:950-962, O-106).
  - The Operator's authority and threshold for each point, set at Intake with evidence (docs/architecture.md:989-995, L-E2).
- **Open questions:**
  1. The text says 'one' exception but names two: docs/spec/packages.md:24-25 says 'The one later exception is `internal/smartif` (phase 3), the smart-if client, and the forge adapter (phase 2).' Which package rows may import net?
  2. Must the smart-if come before the escalation screen and clarification? Both have a deterministic branch under off (docs/architecture.md:977-978), so the architecture does not force this order. This plan puts it first.

### `later-p3-escalation` — Phase 3: the escalation screen, human-input accounting and the audit

- **What:** Each plan, handoff and answer passes the screen: the session's declaration, the code floor (the new-dependencies field, a new entry in the stack's manifest, a licence file, the PRD rows, the intent paths) and the four P1 questions, combined with OR. A selected candidate stops the task in a wait state and posts one escalation brief. The idea owner's 'option N; business-forking: yes|no' is copied to escalations.tsv and is also a row of questions.tsv. Every human action goes into human-inputs.tsv with its class, planned or unplanned. The audit samples screens and answers with a recorded seed at each end-of-milestone retrospective and once more 30 days after the last merge. First demo: a seeded plan that adds a dependency stops, one brief reaches the idea owner, and the answer is classed as planned input. Before the build, docs/spec must add: REQ-008, the schemas of screens.tsv, candidates.tsv, escalations.tsv, human-inputs.tsv, audit.tsv and questions.tsv (if phase 2 does not add it), the per-stack manifest rule and a licence-tool column in the catalog form, and the internal/escalate package row.
- **Specification:** docs/spec/records.md: The layout of the records branch (rows screens.tsv, candidates.tsv, escalations.tsv, human-inputs.tsv, audit.tsv, questions.tsv); docs/spec/setup.md: The stack catalog (a manifest rule and a licence tool per stack); docs/spec/packages.md: The packages of later phases (internal/escalate); docs/architecture.md: 10. Decisions > The escalation screen; docs/architecture.md: 12. Cost, budget and the measures > The records of the measures (Task Intervention Rate, Early Question Share, Missed Escalations; the audit)
- **Packages:** internal/escalate, internal/smartif, internal/catalog, internal/run, internal/forge, internal/forge/github, internal/records
- **Present state:** None. The design is at docs/architecture.md:1020-1079 and :1317-1342, and ADR-0022 decisions 1-4.
- **Tests:**
  - unit: The floor reads only fixed fields and diffs. Each new require in go.mod, indirect ones too, is a candidate unless it is on the allowed list. A rendering of an approved bet is exempt. No answer of P1 removes a floor candidate.
  - unit: A later candidate with the manifest identifier of a rejected option gives a finding and no new brief.
  - unit: The class: planned only for the answer at a planned point listed at Intake, a push that such a point asks for, a stall answer, or an escalation confirmed as business-forking. A reaction or a close is 'other' and unplanned.
  - unit: The audit sample, with a recorded seed and population hash, takes only items whose task merged at least 30 days before and that no earlier audit sampled. An empty sample is 'not measured'.
  - integration: With the httptest forge and the real records writer: each screen run writes one screens.tsv row, selected or not, and a selected task waits. A decision that changes a requirement is rendered as a task.
  - e2e: In CI on a local fake forge, with a scripted plan field and a scripted idea-owner comment: a seeded business-forking decision is selected, stops the agent and reaches the idea owner, and a later attempt that adds the rejected module is refused (W-08 steps 1-9).
  - uat: On a real test target, the idea owner decides a seeded escalation before the task goes on, and confirms the audit's positives. The sample has zero business-forking decisions made by an agent.
- **External inputs:**
  - The idea owner's allowed-dependency list and the intent paths, set at Intake (docs/architecture.md:1034-1037).
  - audit.n and the escalation frequency, set at Intake with evidence (docs/architecture.md:1099, :1108).
  - The GitHub documentation for the permission that lists workflow runs and their triggering_actor.
- **Open questions:**
  1. The REQ-008 criterion measures the Task Intervention Rate against its start value (docs/prd/PRD-0001-layup.md:121), but `layup report` and start-values.tsv are phase 4 (docs/spec/packages.md:90, docs/spec/records.md:66). Can REQ-008 be accepted at the end of phase 3?
  2. Can the App read workflow runs? human-inputs.tsv reads 'workflow runs whose `triggering_actor` is a human' (docs/architecture.md:1319), but the App's permissions are 'contents, issues and pull requests: write; metadata: read' plus commit statuses (docs/architecture.md:140-142), with no Actions permission.

### `later-p3-clarification` — Phase 3: autonomous clarification by the owner role

- **What:** A session that meets a question ends with needs_context and the question. layup run records it, posts it on the task issue, and gets its kind and needs_human (the asker's label, or P2). It sends the question to a session of that kind's owner role, from the confirmed owner map. The answer passes the escalation screen and goes into the asker's next attempt, which is recorded as an attempt. The answer is accepted when that attempt ends 'completed', cites the answer ID, and asks no question that cites it. First demo: a seeded interface-contract question gets an answer from a Software Architect session, accepted without a human, with its asked, answered and accepted times in questions.tsv. Before the build, docs/spec must add: REQ-006, the schema of owners.tsv, the question and answer rows of the transition table, and the clarification columns of questions.tsv.
- **Specification:** docs/spec/records.md: The layout of the records branch (rows questions.tsv, owners.tsv); docs/architecture.md: 8. The phase loop (A question during the work; Handoffs: asking role to owner role, owner role to asking role); docs/architecture.md: 9. Squads and routing (The owner map); docs/architecture.md: 10. Decisions (P2)
- **Packages:** internal/run, internal/route, internal/handoff, internal/session, internal/smartif, internal/escalate
- **Present state:** None. The design is at docs/architecture.md:847-859 and :892-897. W-06 covers it (docs/walkthroughs/W-06-autonomous-clarification.md).
- **Tests:**
  - unit: An answer is accepted only when the next attempt ends completed, cites the answer ID, and none of its questions cites it. The attempt that a question ends is an attempt row that counts toward stall.attempts.
  - unit: The owner map routes each of the four kinds to its owner role. A needs_human question becomes an escalation candidate.
  - integration: With fake harness programs, the real git and the real records writer: question, owner answer, screen, next attempt, accepted. The asking role waits, and its clock does not run during the wait.
  - e2e: In CI on a local fake forge, with a scripted harness that ends with needs_context: the question is answered without a human, and the Clarification Turnaround row exists (W-06 steps 1-8).
  - uat: The Operator confirms or changes the default owner map in the Intake form, and the confirmed map has the comment ID as its evidence. A real session's question gets an accepted answer on a test target.
- **External inputs:**
  - The Operator's confirmation of the owner map at Intake (no fact supports a default, docs/architecture.md:894-897).
- **Open questions:**
  1. The start value of 120 s at the 95th percentile is out of reach while a session start takes minutes (L-D1, docs/architecture.md:1562-1568). Is REQ-006 accepted when the measure is only recorded, not met?

### `later-p3-stalls-budget` — Phase 3: stalls, the panel, the circuit breaker, the budget and the dead-man job

- **What:** Code computes progress from the records (open unknowns and the frozen test list). Five triggers open a stall: no progress, too many rounds, a hang (the session is killed at stall.T), a check that does not report, and an orchestrator takeover. A stall gets a package, a diagnosis by a fresh-context examiner or 'diagnosis-failed', the ladder (retry, a panel with a quorum and a synthesis, the Operator at P3), and one outcome row in stalls.tsv. A diagnosis that names a gate opens an early retrospective. Before each session start, code checks the band and the milestone cap; the circuit breaker stops a milestone at its cap; the bet's caps are checked against U and the appetite; an escalation decision that changes the band writes budget.tsv. A dead-man job in the Operator's control repository notices a silent target. First demo: a seeded no-progress task opens a stall with a stall row, a diagnosis row from another harness, and an outcome row. Before the build, docs/spec must add: REQ-010, the writer and procedure part of REQ-009, the budget check of REQ-011, the schemas of tasks/<task>/tests.tsv, overrides.tsv, budget.tsv and milestones.tsv, the dead-man job (where its workflow, its App-token step, its target list and its last-run commit live, and how they are tested), and the internal/stall package row.
- **Specification:** docs/spec/records.md: REQ-009 — The stall record (Not in phase 1: the writer, the triggers, the examiner, the panel and the Operator's answer form); docs/spec/records.md: REQ-011 — The telemetry record (Not in phase 1: the budget check before each session start); docs/spec/records.md: The layout of the records branch (rows tasks/<task>/tests.tsv, overrides.tsv, budget.tsv, milestones.tsv, stalls.tsv; In other places: the control repository); docs/spec/setup.md: The boundary of phase 1 (row Start 3: the dead-man job's first notice); docs/spec/packages.md: The packages of later phases (internal/stall, internal/ledger); docs/architecture.md: 11. Stalls (the triggers, the ladder, the panel, the circuit breaker, a wrong gate, the dead-man job); docs/architecture.md: 12. Cost, budget and the measures > The budget
- **Packages:** internal/stall, internal/ledger, internal/run, internal/session, internal/smartif, internal/escalate, internal/records
- **Present state:** The stalls.tsv and telemetry.tsv schemas exist with no writer (docs/spec/records.md:144-223). There is no code. The design is at docs/architecture.md:1113-1241 and :1273-1313.
- **Tests:**
  - unit: A round makes progress only when it closes an unknown that was open at the end of the last round, or a test of the frozen list newly passes. A new commit or record alone is not progress, and a test whose source hash changed is no longer passed.
  - unit: Each stall has its three rows in order. A diagnosis-failed row sends the package to the Operator at once and counts as no diagnosis.
  - unit: The panel quorum: a majority of members with a valid output from at least two harnesses; otherwise 'insufficient panel', and the next rung is the Operator. Each member gets the same sealed input.
  - unit: The budget table: past the milestone cap the breaker fires; below B with no unknown part, the work goes on; between B and U, P5 decides or the run escalates; an unknown part never counts as below B; at U it is a hard stop. A bet whose caps pass U or the appetite is refused.
  - integration: With a fake harness that writes no output for stall.T, the session is killed and a hang stall opens. The panel rung is skipped when fewer than two admitted harnesses are free of the diagnosed failure. Waits for a human stop the clock and the round count.
  - e2e: In CI on a local fake forge, with scripted harnesses and a scripted Operator comment: a disagreement and a hang each stop at their limit, a fresh-context examiner writes the diagnosis, and a reroute closes the stall (W-09 steps 1-15; W-10 steps 3, 5, 6).
  - uat: The Operator receives one stall comment with the evidence and answers it in the fixed form (answer, reroute, stop, external). On a real control repository, the dead-man job adds 'no LAYUP run' to a silent target's control issue.
- **External inputs:**
  - A control repository owned by the Operator, and a second GitHub App, layup-watch, with only metadata, contents read and issues write (docs/architecture.md:1228-1241).
  - stall.attempts, ci.T and panel.K, set at Intake with evidence. The band B and U and the wall-clock appetite, from the idea owner (docs/architecture.md:1102-1103, :1275-1279).
- **Open questions:**
  1. Which phase holds the budget check? docs/spec/packages.md:85 gives internal/ledger 'the cost ledger (`telemetry.tsv`) and the budget' phase 2. docs/spec/records.md:60 gives budget.tsv phase 3. docs/spec/records.md:189 names 'the budget check before each session start' as not in phase 1, with no phase.
  2. What does an early retrospective do in phase 3? docs/architecture.md:1224-1226 says a stall that names a gate 'opens an **early retrospective** ... whose rule batch may fix the gate (§13)', but lessons and layup learn are phase 4 (docs/spec/records.md:67-68).

### `later-p4-neutrality` — Phase 4: a second harness agent does and checks the work, and a target continues without LAYUP

- **What:** The criterion run of REQ-013: the same target rules and gates run under at least two harness agents, and each change gets at least one verification from a harness that did not make it. The NFR-002 continuation: with LAYUP removed, the target's gate passes, and a fresh session of a different harness continues one open task from the records alone and lands it under the target's gate. The REQ-013 mechanisms (rules only from the target, the probe of the loaded instruction files, the author set) are specified and built in phase 2 (later-p2-sessions-ledger, later-p2-build-accept). First demo: W-12 on the test target, with a task built on one harness, verified on another, and continued by a third harness with LAYUP removed. Before the build, docs/spec must add: a REQ-013 section that maps the phase-2 parts to the criterion, and the NFR-002 part for a continued task.
- **Specification:** docs/spec/records.md: NFR-002 — A target is independent of LAYUP; docs/architecture.md: 4. Role sessions and model calls (Rules only from the target); docs/architecture.md: 9. Squads and routing (admission outside the authors); docs/architecture.md: 6. Native gates and rule protection (When LAYUP is absent); docs/architecture.md: 14. Coverage (S12 Harness-Agent Neutrality, W-12)
- **Packages:** internal/route, internal/session, internal/run
- **Present state:** None. W-12 (docs/walkthroughs/W-12-harness-agent-neutrality.md) is the design test. ADR-0012 part 3 lists the routed models and harnesses for LAYUP's own work, not for a target.
- **Tests:**
  - integration: With two fake harness programs registered, the real git and the httptest forge: both sessions read the same AGENTS.md from their clones, and the change made on one harness is verified only on the other.
  - e2e: In CI: on a fixture target with no layup/ checks and no LAYUP process, the target's own gate commands pass on a clean tree and fail on a seeded violation (the first part of the NFR-002 criterion).
  - uat: With LAYUP removed, a fresh session of a different harness continues an open task from the records alone and lands it under the target's gate (the NFR-002 criterion, docs/prd/PRD-0001-layup.md:133).
  - uat: The Operator removes the layup/ checks from the ruleset by the setup record's instruction, and checks that the harness entry files in the target are only pointers to AGENTS.md (W-12 steps 12-15).
- **External inputs:**
  - At least two harness agents with working credentials and their documentation (claude, devin, opencode).
  - The 'not used' model list, which also binds test runs (docs/adr/0012-build-layup-in-bootstrap-mode.md:46).
- **Open questions:**
  1. What is left for REQ-013 in phase 4? Phase 2 already needs verification outside the authors (docs/architecture.md:910-918) and a completeness review on another harness (:606-608). PRD §9 puts 'A second harness agent doing and checking work (REQ-013)' in phase 4 (docs/prd/PRD-0001-layup.md:182). This plan reads phase 4 as the criterion run only.

### `later-p4-measures` — Phase 4: layup report, the pilot baseline and the start values

- **What:** `layup report` computes each measure from the records with its population. These are the Layer-1 measures (Missed Escalations #57, Inter-Role Communication Format #59, Telemetry Completeness #60, Stall Diagnosis #61) and the PSB §7.2 measures. It prints 'not comparable' when a measure or its start value is missing, and 'partial' for an unknown input. layup run records the pilot baseline from the Intake answer and the idea owner's start values from the control issue, and renders them into the pilot's PRD as a rendered task. First demo: `layup report` on the test target's records prints Telemetry Completeness and the Stall Rate, and 'not comparable' where no start value exists. Before the build, docs/spec must add: a section for the measures, the schemas of baseline.tsv and start-values.tsv, the table of layup report, and the internal/report package row.
- **Specification:** docs/spec/records.md: REQ-011 — The telemetry record (Not in phase 1: the Telemetry Completeness report, layup report, phase 4); docs/spec/setup.md: The boundary of phase 1 (row The pilot baseline: layup run, phase 4); docs/spec/records.md: The layout of the records branch (rows baseline.tsv, start-values.tsv); docs/spec/packages.md: The packages of later phases (internal/report); docs/architecture.md: 12. Cost, budget and the measures > The records of the measures; docs/architecture.md: 5. Intake and setup (The pilot baseline)
- **Packages:** internal/report, internal/run, internal/records
- **Present state:** None. The design is at docs/architecture.md:1315-1357 and ADR-0024 decision 6.
- **Tests:**
  - unit: Each measure of docs/architecture.md:1317-1329 is computed from fixture records. A missing start value gives 'not comparable', never a pass. An empty audit sample is 'not measured'. Cost per Requirement is 'partial' when a part is unknown.
  - unit: layup report reads files only, and prints the same bytes on two runs (ADR-0013 decision 2).
  - integration: layup report reads a fixture records branch through the real git and prints its table. The rendered start-values task equals code's rendering of the copy.
  - e2e: In CI: layup report on a records branch that the fake-forge runs of earlier milestones made (W-05 step 4; W-10 step 7).
  - uat: The idea owner posts the baseline and the start values in one batch, and reads them back in the rendered PRD.
- **External inputs:**
  - The baseline numbers, measured with the current process, each with evidence, from the idea owner (F-0004#11, docs/architecture.md:425-429).
  - The idea owner's start values in one batch after the baseline (docs/prd/PRD-0001-layup.md:142-146).

### `later-p4-learning` — Phase 4: the retrospective, layup learn and the lessons

- **What:** At each retrospective: the audit, then `layup learn` (the reward per implementing route and a bounded weight proposal, with exploration by learn.explore and a first weight for a route with no weight), then the lessons of a retrospective session, then one brief to the Operator. Under 'propose', the routing register changes only after the Operator's adopt line. A LAYUP lesson becomes an issue on LAYUP's repository. First demo: at the end of a milestone on the test target, one retrospective brief with the reward table and a routing proposal, and the adopted weights in routing.tsv. Before the build, docs/spec must add: a learning section, the schemas of learn/<retrospective>.tsv, lessons.tsv and the weights column of routing.tsv, the table of layup learn, and the internal/learn package row.
- **Specification:** docs/spec/records.md: The layout of the records branch (rows learn/<retrospective>.tsv, lessons.tsv, routing.tsv: the weights); docs/spec/packages.md: The packages of later phases (internal/learn); docs/architecture.md: 13. Retrospective and learning
- **Packages:** internal/learn, internal/route, internal/run, internal/records
- **Present state:** None. The design is at docs/architecture.md:1359-1449, ADR-0025.
- **Tests:**
  - unit: A task of unknown money is left out of the money term, and a route with no known money gets no upward step. The update stays in [learn.min_weight, learn.max_weight], and a route with fewer than learn.min tasks gets no change.
  - unit: A share learn.explore of the implementing role's tasks goes to the other admitted pairs in turn. A route with no weight gets LAYUP's prior, else the mean of its role and tier, else the middle of the bounds.
  - unit: layup learn reads files only and prints the same bytes on two runs (ADR-0013 decision 2).
  - integration: A retrospective runs in its order (audit, learn, lessons, brief). Under 'propose', the routing register changes only after the Operator's adopt line.
  - e2e: In CI on a local fake forge, with scripted sessions and a scripted Operator comment: W-13 steps 1-8.
  - uat: The Operator adopts a routing proposal at a real retrospective, and a LAYUP lesson becomes an issue on LAYUP's repository (W-13 step 9).
- **External inputs:**
  - The reward term weights and the learn.* parameters, with the Operator's comment as evidence (docs/architecture.md:1104).
- **Open questions:**
  1. Which requirement does layup learn serve? PRD-0001 §6 has no requirement for the learning loop; its source is the vision brief, which 'is input for ADRs, not a source of requirements' (docs/prd/PRD-0001-layup.md:14). ADR-0012 part 1 (docs/adr/0012-build-layup-in-bootstrap-mode.md:44) starts a task only when 'its plan names the In-Scope fact it serves', and F-0003#41-#52 names no learning. docs/spec/packages.md:93-94 sets a package's phase by 'the earliest PRD-0001 phase whose requirement needs it', but gives internal/learn phase 4 with no requirement. Does the Operator want this milestone before the pilot?

### `later-p4-pilot` — Phase 4: the pilot on two problem statements with different stacks

- **What:** LAYUP sets up a target from each of two pilot problem statements with different stacks, and the role agents that it orchestrates deliver each product until the idea owner accepts it. The second stack gets its catalog entry (gate kinds, tools, versions, fixtures, evidence, the manifest rule and the licence tool of the escalation floor) under LAYUP's gate. The pilot measures each criterion of PRD-0001 §7 that says 'in the pilot' or 'on a pilot target'. First demo: on a target of the second stack, `layup setup verify` passes and each active kind's fixture fails; after the first bet's activation batch, a seeded violation of each of the four gate kinds gives 'fail'. Before the build, docs/spec must add: REQ-014, the second stack's catalog entry in the form of 'The stack catalog', and the pilot's audit for the setup values (REQ-002, NFR-003) and the four Won't rows (REQ-015 to REQ-018).
- **Specification:** docs/spec/setup.md: The stack catalog; docs/spec/setup.md: NFR-003 — No value without evidence (item 4: the audit is a review); docs/spec/gate.md: REQ-004 — The stack gates of a target, run from outside; docs/architecture.md: 6. Native gates and rule protection > The stack gates (a new stack gets an entry through a LAYUP change); docs/architecture.md: 10. Decisions > The escalation screen (the floor per stack); docs/architecture.md: 15. Known limits (L-A1, L-A3 close before a target that is not the Operator's)
- **Packages:** internal/catalog, internal/escalate
- **Present state:** None. The catalog form is specified at docs/spec/setup.md:274-302, and the Go entry is phase-1 work of #29. PRD §11 question 6 is open (docs/prd/PRD-0001-layup.md:203).
- **Tests:**
  - unit: The second stack's kinds.tsv validates against the catalog-kinds schema, and each active kind has a fixture.
  - integration: LAYUP's CI runs each fixture of the second entry, and each fixture makes its kind fail.
  - e2e: In CI: layup setup and layup setup verify on a fixture problem statement of the second stack, with check 'adapted' passing (the baseline rules are byte-identical apart from the recorded adapted values, REQ-018).
  - uat: For each pilot problem statement, layup setup creates the target, layup run delivers the product, and the idea owner accepts it requirement by requirement (the REQ-014 criterion).
  - uat: An audit of every setup value finds that its cited source supports it (REQ-002, NFR-003). Each intent decision is recorded as the idea owner's (REQ-016). A code review of the release records no path for REQ-015 or REQ-017.
- **External inputs:**
  - The two pilot problem statements and their stacks, chosen by the idea owner at Decision Point 1 (docs/prd/PRD-0001-layup.md:203).
  - The second stack's gate tools, their versions and their documentation at that version.
  - The second stack's toolchain on the LAYUP host (L-B1) and in LAYUP's CI, to run its fixtures.
  - The pilot targets on GitHub with the App installed, and a decision whether they are the Operator's own (L-A1, L-A3 'Close when').
- **Open questions:**
  1. 'Pilot' has two readings. docs/adr/0012-build-layup-in-bootstrap-mode.md:49 says 'The first pilot is done when `layup` has set up one target repository from a problem statement and has run that target's gate from outside', which is reachable at the end of phase 1. docs/prd/PRD-0001-layup.md:182 makes 'the pilot on two problem statements with different stacks' a phase-4 item. Which one ends bootstrap mode, and which one do the 'in the pilot' criteria of §7 mean?
  2. The REQ-004 criterion needs 'a target of each pilot stack' (docs/prd/PRD-0001-layup.md:117), but phase 1 has only the Go entry (docs/spec/setup.md:280). Is REQ-004 accepted in phase 1 on Go alone?

### Constraints of area `later`

1. **docs/spec/README.md:9-12** — It is written one milestone at a time (task T-0drh, #74; decision O-114). This version covers phase 1 of PRD-0001 §9 only. A later milestone adds its own sections; it does not rewrite the package table or the layout of the records. *Effect:* Each later milestone starts with a specification task that adds its docs/spec sections and record schemas. The plan cannot slice a later milestone into build issues before that specification exists.
2. **docs/spec/packages.md:70-72** — Each later milestone gives a package its row in the table above, with its import rules, in the same change that adds its requirement's section. *Effect:* The package row and import rules of each later package are deliverables of the milestone's specification task, not of a build task.
3. **docs/spec/records.md:10-14** — One row per record kind of the whole architecture ... so that a later milestone adds schemas and does not move a file. *Effect:* Later milestones keep the paths of the layout. A record that the walkthroughs commit but the layout does not list (the psb gap table, the spec check tables, the audit table) needs a row there first.
4. **docs/spec/README.md:75-77** — Each record has one schema block. A Go test of the code reads these blocks and compares them with what the code writes and reads (the test comes with the code, Bootstrap mode rule 1). *Effect:* Each milestone that writes a new record kind adds its schema block and this comparison test.
5. **docs/spec/packages.md:16-25** — go.mod has no require line ... Only the package internal/git starts the git program ... No package of phase 1 imports net, net/http or crypto/tls ... The one later exception is internal/smartif (phase 3), the smart-if client, and the forge adapter (phase 2). *Effect:* The GitHub adapter and the smart-if client use only the standard library (net/http, crypto for the App JWT). No SDK is allowed. The GH013 probe goes through internal/git, because only that package starts git.
6. **docs/adr/0013-orchestrate-a-target-from-outside-with-layup-run.md:41-45** — The engine checks are pure. layup psb check, layup setup verify, layup gate, layup spec check, layup report and layup learn read files, and layup audit reads files and the forge's read-only API; each prints a typed table. *Effect:* layup spec check, layup report and layup learn are separate pure commands with unit and integration tests on fixture records. Only layup run writes.
7. **docs/adr/0014-keep-the-records-in-the-target-with-one-writer.md:42-47** — Only layup run commits to the records branch, from its own clone, with pushes that are never forced. ... a result whose attempt is no longer the task's open attempt is refused *Effect:* Every later record writer goes through internal/records inside layup run, so every phase-2 milestone depends on later-p2-run-start.
8. **docs/architecture.md:413-418, 546-550** — the Operator applies the rulesets, which hold the same required checks and LAYUP's ... the required checks are the native gate jobs, pinned to GitHub Actions, and layup/gates, layup/spec, layup/verify and layup/rules, pinned to the LAYUP App *Effect:* After Scaffold 6, a merge needs all four layup/ statuses. So the basic layup/rules status is built with the task loop (later-p2-phase-loop-handoffs), before the first merge at the first bet.
9. **docs/architecture.md:451-457** — A gate kind whose rules depend on the architecture ... has the state pending until its activation ... Its job then fails a pull request that changes a path in the product's scope *Effect:* No build task that changes a product path can merge before the activation batch. So later-p2-build-accept comes after later-p2-rule-protection.
10. **docs/architecture.md:31-39** — Only three things go into a target ... No file goes in that the target needs LAYUP to build, test or pass its gates. *Effect:* No later milestone may add a LAYUP binary, script or CI job to a target (NFR-002).
11. **docs/architecture.md:60-62** — All six are required: an adapter that cannot give one stops the setup at Start or at the probes of §3 ... No other adapter is designed (known limit L-A2). *Effect:* Phase 2 builds only the GitHub adapter. The pilots must run on GitHub.
12. **.github/workflows/ci.yml:10-12, 305-318; docs/tests/test-levels.md:63-81** — permissions: contents: read ... name: end-to-end run: go test -tags=e2e -timeout 10m ./... ... It is judged by a person, not asserted by a command, so it is not a rung of the automated ladder *Effect:* The e2e tests of later milestones run in CI with no secret and within 10 minutes, so they use a local fake forge and scripted harness programs. A run on a real forge, harness or provider, and every human step, is UAT.
13. **docs/tests/template-unit.md:20-23** — A stand-in replaces a real collaborator (a database, a network call, the clock) with something predictable and local, so the test touches no real file, network, or clock. *Effect:* A test of the GitHub adapter against httptest, or of a real directory tree, is at the integration level, not the unit level.
14. **docs/prd/PRD-0001-layup.md:184** — The four Won't rows (REQ-015 to REQ-018) hold in every phase. *Effect:* Each milestone's release review records that no code path trains a model or changes a cloud provider, and that the baseline rules are unchanged apart from the adapted values.
15. **docs/adr/0012-build-layup-in-bootstrap-mode.md:44** — A task in this repository is one of: a PSB In-Scope item (F-0003#41–#52) or a child of one; ... No new process rule, routing rule, check script or policy capture starts unless a product task needs it and its plan names the In-Scope fact it serves. *Effect:* Each later milestone names the In-Scope fact that it serves. later-p4-learning has none, so it needs the Operator's decision before it starts.
16. **docs/adr/0012-build-layup-in-bootstrap-mode.md:49** — The mode ends when the ADR that supersedes this record is accepted ... the task that closes the pilot opens that ADR, which re-decides the full gate with the pilot's numbers. *Effect:* The review process of later phases can change after the first pilot. The plan must not tie the gate of phases 2 to 4 to the bootstrap rules.
17. **docs/adr/0012-build-layup-in-bootstrap-mode.md:46** — Not used: the models the Operator listed as not to use ... Haiku; deprecated models; -fast and -priority variants *Effect:* The harness register, the admission tests and the UAT runs of later phases must not use these models, also inside test runs.
18. **docs/spec/README.md:31-34** — A value that needs a decision of the Operator or the idea owner is never set here: it stays a marker with a row in open-gaps.tsv (Invariant 4). *Effect:* The parameters that later milestones need (stall.attempts, ci.T, panel.K, audit.n, build.parallel, thresholds, learn.*) stay open gaps until a human sets them with evidence. Tests use fixture values only.
19. **docs/architecture.md:1515-1516, 1536-1537** — Close when: before LAYUP runs on a target that is not the Operator's; each session and each gate run gets a container or another user ... each forge write gets an ID, and a retry reads the forge first. *Effect:* A pilot on a target that is not the Operator's first needs session isolation and idempotent forge writes. No PRD phase has that work, so the plan must add it or keep the pilot targets as the Operator's own.
20. **docs/prd/PRD-0001-layup.md:211-213** — the implementation plan (T-55n2) fills the Task column, and each delivering task fills the Test column. *Effect:* Each later milestone maps to the requirement rows of §12, so the Task column can name it. A criterion measured 'in the pilot' maps to later-p4-pilot as well as to its building milestone.

### Risks of area `later`

- Phase 2 is much larger than PRD §9 shows. PRD §9 names three requirements, but the specification puts the whole orchestrator into phase 2: layup run, the forge adapter, sessions, the ledger writer, routing, the Scaffold, rule batches and audit. That is seven large milestones. A plan that sizes phase 2 by its three requirements will underestimate it.
- Phase 2 needs two admitted harnesses: for the completeness review, the plan review and the verification outside the authors. REQ-013 is phase 4. Without a second working harness, phase 2 can merge no task.
- Phases 2 and 3 run sessions with no escalation screen until later-p3-escalation, so a business-forking choice has no guard. A hang has only the wall-clock limit of harness.<id>.wall until the stall triggers of phase 3.
- The real-world checks of phases 2 to 4 need a real GitHub target, the App key, paid harness sessions and a smart-if provider. They are UAT runs that cost money, are slow and are not deterministic. CI e2e uses fakes, so a mismatch between the fake forge and real GitHub can stay hidden until a UAT run.
- The shared host (L-A1): sessions and gate runs share the Operator's user, so a session can reach the App key and the Operator's gh login. A pilot on a target that is not the Operator's needs isolation that no phase delivers.
- A newer baseline (L-A6) can break the setup steps between phase 1 and the pilot, because each target pins the baseline's latest commit.
- The measures of phases 3 and 4 need a baseline measured 'with the current process' and start values from the idea owner. If these come late, layup report prints only 'not comparable', and the Layer-2 criteria cannot be accepted.
- The specification is written one milestone at a time, so each milestone adds a specification task and its review before the build. This adds approval load on the Operator and can delay the build.
- The open phase conflicts (questions.tsv, routing.tsv, the budget check, the frozen test list, the dead-man job, the Shape panel, the retrospective) can each move work between phases. They must be decided before the milestone specifications are written, or the specifications will contradict records.md and packages.md.
- A target that phase 1 set up has no layup/ checks in its ruleset and a records branch first pushed by the Operator. Until the adoption path is specified, phase 2 can drive only new targets.

### Critic points that the reviser of `later` rejected

- later-p3-clarification depends_on: add later-p3-stalls-budget, because its unit test uses the stall.attempts count. *Basis:* docs/architecture.md:851-853 says an attempt that a question ends 'counts toward `stall.attempts` (§11)'. The count is an attempt row of tasks/<task>/events.tsv, which is phase 2 (docs/spec/records.md:44). The limit is enforced by the stall trigger in later-p3-stalls-budget, which reads that row. So clarification writes the row and does not depend on stalls-budget. I changed the unit test to say this.
- later-p3-escalation depends_on: add later-p3-clarification, because each escalation to the idea owner is a row of questions.tsv. *Basis:* docs/architecture.md:1322 makes each escalation a row of questions.tsv, but no source says the clarification milestone must build that writer. The cycle breaks when escalation writes its own questions.tsv rows (or phase 2 does, if the open question puts questions.tsv there, docs/spec/records.md:37 vs docs/architecture.md:373-375). Clarification then depends on escalation, which it needs anyway for the screen (docs/architecture.md:855).
- later-p3-escalation depends_on: add later-p3-stalls-budget, because a decision that changes the band writes budget.tsv. *Basis:* budget.tsv belongs to the budget work (docs/spec/records.md:60), and the budget check needs escalation (docs/architecture.md:1310-1313, W-10 step 6). So stalls-budget depends on escalation and applies band changes itself. The reverse dependency would make a cycle and is not needed.

## The cross-area check

### Duplicates

- `setup-catalog`, `gate-catalog-package` → `gate-catalog-package` (area `gate`). Both items build internal/catalog: embed the entries, read kinds.tsv by the schema catalog-kinds, replace {{module}}, give the fixture paths. gate-catalog-package also gives the derived manifest, the <stack>/<path> lookup for check sources, the config paths for the S15 register, and the embed rules (no go.mod in the embedded tree, the all: prefix for .github/). The gate area names the stack catalog in its scope. Point setup-s01-questions, setup-s12-stack-files, verify-sources, verify-jobs and verify-gate-fixtures to gate-catalog-package.
- `gate-verify-kind`, `verify-gate-fixtures` → `verify-gate-fixtures` (area `verify`). Same check gate:<kind>, with the same first-match rules (setup.md:256). packages.md:44 puts it in internal/verify, and gate-verify-kind asks for the merge. Keep from gate-verify-kind the test for a fixture that does not apply and its open question.
- `gate-nfr005-imports`, `found-boundary-test` → `found-boundary-test` (area `found`). Both check rule 5 of packages.md (no net, net/http or crypto/tls). found-boundary-test also checks rules 1 to 4 and the table. Keep from gate-nfr005-imports the fixture module with its own go.mod that imports net/http, and the measured fact that go list -deps -test names os/exec, so rule 4 cannot use -test.
- `later-gate-status`, `later-p2-phase-loop-handoffs`, `later-p2-build-accept` → `later-p2-phase-loop-handoffs` (area `later`). later-p2-phase-loop-handoffs posts the status layup/gates. later-p2-build-accept holds the rule that no PR is ready or merged while a selected kind did not run. later-gate-status adds no other work.
- `later-gate-rule-batch`, `later-p2-rule-protection` → `later-p2-rule-protection` (area `later`). Same work: the batch's own gate files on its head, each recorded known-bad patch, the approval hash, and the activation at the first bet.
- `later-psb-meaning-review`, `later-psb-intake-batch`, `later-p2-intake-spec` → `later-p2-intake-spec` (area `later`). later-p2-intake-spec holds the review of meaning, the merge into one Intake batch, and the answers as a raw fact. later-psb-meaning-review lists no dependency, but it needs layup run and sessions (later-p2-run-start, later-p2-sessions-ledger).
- `later-psb-early-question-share`, `later-p4-measures` → `later-p4-measures` (area `later`). layup report computes the Early Question Share in later-p4-measures. The questions.tsv rows with the flag 'before delivery' come from later-p3-escalation, or from later-p2-intake-spec if questions.tsv moves to phase 2 (records.md:37 against architecture.md:373-375).
- `setup-pilot-acceptance`, `gov-first-pilot`, `gov-pilot-inputs` → `gov-first-pilot` (area `gov`). Same pilot run and same NFR-003 audit (PRD-0001:115, :134). The inputs that setup-pilot-acceptance lists are gov-pilot-inputs. Also move into gov-first-pilot the UATs that use the same pilot target: verify-acceptance (the Operator reads the table; the audit), gate-command (layup gate on the pilot target), gate-catalog-go-workflow (a PR on GitHub) and setup-s13-rulesets (the ruleset apply). One human session then gives all the evidence.
- `setup-e2e-offline`, `verify-acceptance`, `found-nfr001-tests`, `found-nfr002-tests` → `setup-e2e-offline` (area `setup`). All four make one full offline setup: stand-in baseline, all answers, layup setup to S15, layup setup verify. setup-e2e-offline owns the full run and the scenario 'a correct setup exits 0 with only pass or clear rows', which verify-acceptance repeats. verify-acceptance keeps its seeded defects; the NFR-001 and NFR-002 items keep their assertions. Each scenario stays its own test (template-e2e.md:22), but all use one fixture builder and share full runs where they can, because all e2e tests of cmd/layup share one 10-minute limit (ci.yml:318). found-nfr001-tests repeats the commands.sh order (setup-commands-file) and the orphan-branch test (setup-s15-records-commit); found-nfr002-tests repeats the 'no layup/ check' test of setup-s13-rulesets.
- `verify-test-baseline`, `setup-e2e-offline`, `found-e2e-harness` → `verify-test-baseline` (area `verify`). Three items build the stand-in baseline repository: verify-test-baseline, the generator inside setup-e2e-offline, and the fixture-repository helper of found-e2e-harness (the dangling key setup-baseline-fixture also means it). Build it once. internal/verify integration tests and cmd/layup e2e tests both need it, so the builder needs a package with a row in packages.md, or each test package makes it again (see conflicts).
- `verify-markers`, `setup-s10-markers` → `verify-markers` (area `verify`). Both port the scan of check_markers, embed MK_EXEMPT, test that copy against setup-check.sh:305, and run the parity cases of docs/setup/tests/markers/. Keep one scanner in verify-markers. setup-s10-markers keeps the stop rows and the M-<x8> IDs, and gets the marker list through internal/cli, as for the gap table, because internal/setup cannot import internal/verify (packages.md:43). Add the edge setup-s10-markers -> verify-markers.
- `psb-batch-api`, `setup-s01-questions`, `setup-s06-facts`, `found-cli` → `psb-batch-api` (area `psb`). Four items hold the same seam: internal/cli hands the gap table to internal/setup (setup.md:86), as bytes or as a setup type, at S01 and again at S06. The integration test of setup-s01-questions (the Q-NNN IDs equal those of layup psb check) repeats the test of psb-batch-api. Decide the seam once in psb-batch-api. Add the edges setup-s01-questions -> psb-batch-api and setup-s06-facts -> psb-batch-api.
- `psb-cli-contract`, `found-e2e-harness`, `found-cli` → `psb-cli-contract` (area `psb`). The psb scenarios (the golden on the real PSB, two equal runs, exit 2 with an empty stdout) are in psb-cli-contract. found-e2e-harness repeats the golden and repeat scenario, and found-cli repeats the psb exit-2 cases. found-e2e-harness keeps the helpers and the version scenario. found-cli owns the usage text, which answers an open question of psb-cli-contract. Add the edges psb-cli-contract -> found-e2e-harness and psb-cli-contract -> found-cli, so that two items do not edit cli.go in parallel.
- `psb-field-rule-cr`, `psb-tsv`, `found-tsv` → `found-tsv` (area `found`). The field rule (a tab, LF or CR becomes one space; an empty field is the dash U+2014) has one owner: the writer of found-tsv. psb-field-rule-cr puts the CR rule into internal/psb only for psb-tsv to remove it. found-tsv has no dependency and can land first; psb-field-rule-cr can then keep only its tests and the deletion of psb-check.md:23-24. Two more questions are asked twice: who writes the dash for line 0 (psb-tsv, found-tsv), and the width of N in id(Q-NNN) (psb-batch-api, found-tsv).
- `gov-issue-21`, `verify-markers` → `verify-markers` (area `verify`). ADR-0012:60 batches #21 into the product task that touches the same files. verify-markers changes the same check and has the open question on #21 (fix sh first, or make the Go port correct). Put the sh fix and its fixture into verify-markers.
- `gov-issue-61`, `gov-issue-48`, `verify-facts` → `verify-facts` (area `verify`). #61 (a fact of one tab) and #48 (autocrlf breaks check facts and setup-check.sh) touch check facts and setup-check.sh. verify-facts has open questions on both. Batch both into verify-facts (ADR-0012:60).
- `gov-issue-34`, `setup-rule-paths`, `setup-s04-pin` → `setup-rule-paths` (area `setup`). The glossary row 'rule path' of #34 goes with setup-rule-paths, which writes rule-paths.tsv. The 'ADR for the pin' wording of #34 goes with setup-s04-pin, which has the same open question (setup.md:89).
- `found-git`, `setup-runner`, `verify-harness` → `found-git` (area `found`). Three items isolate git from the host configuration and need a commit identity: found-git (hostile-config test; the identity is an Operator decision), setup-runner (the same test and the same decision as an external input), verify-harness (no global identity, empty HOME). Make one mechanism in internal/git and one Operator decision; the other items use it.
- `found-cli`, `gate-command`, `setup-runner`, `verify-frame` → `found-cli` (area `found`). found-cli owns the map from result words to exit codes and the parser that accepts REPO before the flags. gate-command, setup-runner and verify-frame each test the same map again. Each command item uses the map of found-cli and tests only its own result words. Add the edge gate-command -> found-cli.
- `gate-catalog-go-workflow`, `verify-jobs` → `verify-jobs` (area `verify`). Both need a reader of job names in workflow text: a unit test of gate-catalog-go-workflow (job names equal the kinds; a form that the reader of check jobs can read) and check jobs. internal/catalog cannot import internal/verify (packages.md:41). Keep one reader in internal/verify and move that unit test there.
- `gov-milestone-spec`, `later-p2-run-start`, `later-p2-sessions-ledger`, `later-p2-intake-spec`, `later-p2-scaffold`, `later-p2-phase-loop-handoffs`, `later-p2-rule-protection`, `later-p2-build-accept`, `later-p3-smartif`, `later-p3-escalation`, `later-p3-clarification`, `later-p3-stalls-budget`, `later-p4-neutrality`, `later-p4-measures`, `later-p4-learning`, `later-p4-pilot` → `gov-milestone-spec` (area `gov`). Each later item lists the docs/spec sections that it needs 'before the build'. That is the work of gov-milestone-spec, one specification task per milestone (docs/spec/README.md:9-12). Make each later milestone depend on its specification task, and use the lists in the later items as input to it.

### Dangling dependencies

- `gate-manifest, gate-catalog-package, verify-frame` names `found-schema-test`: found-spec-schema-test.
- `setup-runner` names `found-schema-block-test`: found-spec-schema-test.
- `setup-runner` names `found-cli-exit-codes`: found-cli.
- `gate-nfr005-imports` names `found-import-rules`: found-boundary-test. gate-nfr005-imports merges into it (see duplicates).
- `gov-test-traceability` names `found-stdlib-only-test`: found-boundary-test. Use the key only as a name in the traceability table. A build edge from a gov item that must merge before gov-pdr-approval makes a cycle (see cycles).
- `later-p2-run-start, later-p2-sessions-ledger, later-p3-stalls-budget` names `found-record-schemas`: No item has this key. For later-p2-sessions-ledger it is found-telemetry-schema. For later-p3-stalls-budget it is found-stalls-schema. For later-p2-run-start it is found-spec-schema-test; the schemas of start.tsv, approvers.tsv and lease.tsv have no item and come from the phase-2 specification task (gov-milestone-spec).
- `gate-command` names `setup-steps`: The setup chain to S15 (setup-s15-records-commit), for the UAT only. This edge closes a cycle through setup-s12-stack-files and gate-catalog-go-entry. Remove it and move the UAT to gov-first-pilot.
- `gate-catalog-go-workflow` names `setup-s12`: setup-s12-stack-files. S12 writes the workflow file from the catalog (setup.md:97, :285), so this edge makes a cycle. Remove it: the refusal to overwrite a baseline path is a test of setup-s12-stack-files.
- `found-nfr002-tests` names `setup-s12-gates`: setup-s12-stack-files.
- `found-nfr001-tests, found-nfr002-tests` names `setup-s13-ruleset`: setup-s13-rulesets.
- `found-nfr001-tests, found-nfr002-tests` names `setup-s15-records-branch`: setup-s15-records-commit.
- `later-p2-scaffold` names `setup-records-branch`: setup-s15-records-commit.
- `found-nfr001-tests, found-nfr002-tests` names `setup-s01-s14`: No single key. Use setup-e2e-offline: it holds the full run and depends on every step through setup-s15-records-commit.
- `found-nfr001-tests, found-nfr002-tests` names `setup-baseline-fixture`: verify-test-baseline (the one stand-in baseline; see duplicates).
- `found-nfr002-tests, gov-first-pilot, gov-operator-ci-protection, later-p2-rule-protection, later-p2-build-accept, later-p4-pilot` names `setup-catalog-go`: gate-catalog-go-entry (the gate area owns the Go entry). found-nfr002-tests and gov-first-pilot also need gate-catalog-go-workflow for the CI jobs.
- `setup-s12-stack-files, setup-e2e-offline` names `gate-catalog-go`: gate-catalog-go-entry. setup-s12-stack-files also needs gate-catalog-go-workflow, because S12 writes the CI workflow of the gate jobs from the catalog (setup.md:97, :285).
- `setup-s12-stack-files` names `gate-manifest-reader`: gate-manifest.
- `later-p2-sessions-ledger` names `setup-rule-path-register`: setup-rule-paths.
- `setup-s01-questions` names `psb-on-tsv`: psb-tsv. setup-s01-questions and setup-s06-facts also need psb-batch-api (the hand-off of the gap table), which no setup item lists.
- `later-p2-intake-spec` names `psb-check`: No item has this key. The phase-1 command is psb-cli-contract, and the seam for the batch is psb-batch-api. Use both.
- `gate-verify-kind` names `verify-runner`: verify-frame. gate-verify-kind merges into verify-gate-fixtures (see duplicates).
- `setup-s15-records-commit` names `verify-command`: verify-frame (the command and its table).
- `setup-e2e-offline` names `verify-command`: verify-frame and every check item: verify-pin, verify-kit-history, verify-markers, verify-adapted, verify-identity, verify-facts, verify-onboarding-glossary, verify-guardrails, verify-baseline-scripts, verify-sources, verify-jobs, verify-gate-fixtures. Only this edge brings in verify-adapted, because adapted is the evidence of no step (setup.md:86-100).
- `setup-s04-pin` names `verify-check-pin`: verify-pin.
- `setup-s05-history` names `verify-check-kit-history`: verify-kit-history.
- `setup-s05-history` names `verify-check-link-lint`: verify-baseline-scripts. It gives only exit <code>, not the list of files with broken links that S05 needs (see conflicts).
- `setup-s06-facts` names `verify-check-facts`: verify-facts.
- `setup-prose-inputs` names `verify-check-prose`: verify-onboarding-glossary, verify-guardrails and verify-identity.
- `setup-s11-fill` names `verify-check-markers-sources`: verify-markers and verify-sources.
- `setup-s12-stack-files` names `verify-check-jobs-gates`: verify-jobs and verify-gate-fixtures.
- `found-nfr001-tests, found-nfr002-tests, later-p2-scaffold` names `verify-checks`: No aggregate key. It means verify-frame and the 12 check items. For the two found items, setup-e2e-offline already brings them all.
- `gov-first-pilot` names `verify-setup-verify`: verify-acceptance, which depends on verify-frame and every check.
- `gov-issue-21` names `verify-setup-verify`: verify-markers (gov-issue-21 merges into it).
- `gov-issue-61` names `verify-setup-verify`: verify-facts (gov-issue-61 merges into it).
- `later-gate-status` names `later-run`: later-p2-run-start.
- `later-gate-rule-batch` names `later-rules`: later-p2-rule-protection.

### Cycles

- `gate-command` → `setup-s15-records-commit` → `setup-rule-paths` → `setup-s12-stack-files` → `gate-catalog-go-entry` → `gate-command`
- `gate-catalog-go-workflow` → `setup-s12-stack-files` → `gate-catalog-go-workflow`
- `gov-pdr-approval` → `gov-test-traceability` → `found-boundary-test` → `gov-pdr-approval`
- `verify-gate-fixtures` → `setup-s12-stack-files` → `verify-gate-fixtures`
- `verify-jobs` → `setup-s12-stack-files` → `verify-jobs`

### Layers (a topological order of the phase-1 items)

0. `gov-rescope-29`, `gov-pilot-definition`
1. `gov-prd-task-column`, `gov-rescope-33`
2. `gov-test-traceability`
3. `gov-pdr-approval`
4. `found-tsv`, `found-git`, `found-boundary-test`, `found-e2e-harness`, `psb-field-rule-cr`
5. `found-spec-schema-test`, `found-cli`, `psb-rule-edges`, `psb-rule-values`, `psb-utf8-input`, `verify-kit-history`, `verify-test-baseline`
6. `found-telemetry-schema`, `found-stalls-schema`, `psb-tsv`, `psb-cli-contract`, `gate-manifest`, `gate-catalog-package`, `setup-runner`, `verify-harness`
7. `gate-scratch-tree`, `psb-batch-api`, `setup-answers`, `setup-commands-file`, `verify-frame`
8. `gate-results`, `setup-s01-questions`, `verify-pin`, `verify-markers`, `verify-adapted`, `verify-identity`, `verify-facts`, `verify-baseline-scripts`
9. `gate-run`, `setup-s02-s03-baseline`, `verify-onboarding-glossary`, `verify-guardrails`, `verify-sources`
10. `gate-command`, `setup-s04-pin`, `verify-gate-fixtures`
11. `gate-catalog-go-entry`, `setup-s05-history`
12. `gate-catalog-go-workflow`, `setup-s06-facts`
13. `verify-jobs`, `setup-prose-inputs`
14. `setup-s10-markers`
15. `setup-s11-fill`, `gov-pilot-inputs`
16. `setup-s12-stack-files`
17. `setup-s13-rulesets`, `setup-rule-paths`
18. `setup-s15-records-commit`
19. `setup-e2e-offline`
20. `verify-acceptance`, `found-nfr001-tests`, `found-nfr002-tests`
21. `gov-first-pilot`
22. `gov-pilot-numbers`, `gov-release-review`
23. `gov-supersede-adr-0012`

### Coverage per requirement

- **REQ-001** (15 items): Each heading of psb-check.md has an item: the command (psb-cli-contract, psb-utf8-input), the rules (psb-rule-edges, psb-rule-values), the table (psb-tsv, psb-field-rule-cr), Not in phase 1 (later-p2-intake-spec; the answers as a raw fact in setup-s06-facts). Gaps: (1) The second and third sentences of the §7.1 criterion (one Intake batch with the gaps of meaning; the Early Question Share) need phases 2 to 4; no phase-1 item meets them. (2) psb-rule-edges waits for an Operator decision on G1 for '**Technology stack:**' with no value. (3) No item sets one rule for invalid UTF-8 for both psb-utf8-input and found-tsv. (4) No item limits the gap table to 999 rows, the width of id(Q-NNN).
- **REQ-002** (39 items): Each heading of setup.md under REQ-002 has an item. Gaps: (1) No item holds the spec fix that lets the real baseline pass check adapted and check kit-history; without it the pilot fails at S05 and S15 (see conflicts). (2) No source gives the fixed texts that steps write: the pin ADR (S04), the README of layup-records (S15), the form, number and file names of the answers records (S06, S11), the module path from 'name' (S12). Each item carries them as open questions, and no item collects the Operator decisions before the issues start. (3) The evidence of S12 names a check 'gates' that does not exist. (4) The evidence of S13 (the Operator's run of commands.sh) cannot be seen, so S13 gets no done row, and the step table has no result for a step that did not run. (5) The §7.1 clause 'answers in Git before the next step' disagrees with the step order. (6) Check identity 'with the target's name' has no rule (setup.md:252), and no check reads that the README names layup-records (architecture.md:111).
- **REQ-004** (15 items): Each heading of gate.md under REQ-004 has an item. Gaps: (1) The §7.1 criterion needs a seeded violation of each of four kinds on each pilot stack. Phase 1 has only Go, and layout, boundary and contract stay pending with no fixture (setup.md:296), so only the pending rule fails them. (2) The state of the test kind at setup is not settled. (3) No item settles the exit code when git is missing or the scratch tree fails (gate.md:52-54). (4) No gate item gives a progress line for a run over ten seconds.
- **REQ-007** (11 items): Items 1 to 3 of 'In phase 1' have items. Gaps: (1) Item 1 rests on the ruleset that the Operator applies. Phase 1 has no plan check (setup.md:29), and only UATs (collected in gov-first-pilot) show that a failed job blocks the merge. (2) Item 3 says that a kind with two verdicts on one input is a catalog defect; gate-catalog-go-entry has no repeat test for its kinds, only the repeat rule of gate-command on a fixture. (3) The §7.1 criterion counts PRs in the pilot, which needs layup run (phase 2).
- **REQ-009** (3 items): Schema only; found-stalls-schema and found-spec-schema-test cover the block. Gap: no package row of phase 1 holds the type (internal/stall is phase 3, packages.md:88). The §7.1 criterion needs the writer of phase 3.
- **REQ-011** (5 items): Schema only; found-telemetry-schema covers telemetry.tsv and host:prices.tsv. Gap: no package row of phase 1 holds the types (internal/ledger is phase 2, packages.md:85). The §7.1 criterion needs the writer (phase 2) and layup report (phase 4).
- **NFR-001** (14 items): Item 2 (the setup record, the register, the README, the first records commit): setup-s15-records-commit and found-nfr001-tests. Item 3: setup-s06-facts and setup-s11-fill write the answers, but no item puts the 'source' and 'by' of each answer into the target (records.md:105-106 against setup.md:91). Item 4: no test rebuilds a work area from the inputs and the target's Git. Item 1 and 'Not in phase 1': later-p2-run-start, later-p2-scaffold, later-p2-rule-protection. The §7.1 audit needs a pilot task.
- **NFR-002** (7 items): Items 1 to 4 have items. Gaps: the README text of item 4 has no source, so 'how to read it with no tool' (records.md:25) is checked only for presence. The second half of the §7.1 criterion (a fresh session of another harness continues a task) is phase 4 (later-p4-neutrality).
- **NFR-003** (16 items): Items 2 and 3 have items. Gaps: (1) Item 1: no item gives a record row with a source for some written values: the number, date and status of the pin ADR (S04), the values of branch-protection.json and the ruleset (S13), the date and number of the answers records (S06, S11), and the README text (S15). (2) Check sources has no rule for computed rows (setup.md:254). Item 4 (the audit) is a review in gov-first-pilot.
- **NFR-004** (12 items): Items 1, 2, 4 and 5 have items. Gaps: (1) Item 3 disagrees with the result table for a pending kind with a missing tool (see conflicts). (2) No item settles the exit code when git itself is missing or fails (README.md:55-58); found-git, gate-scratch-tree, setup-runner and verify-frame each leave it open. (3) A manifest with a header and no row has no rule (gate-run); it can give an empty table and exit 0 (FT1).
- **NFR-005** (11 items): Item 1 (the import rule): found-boundary-test, with gate-nfr005-imports merged into it. Item 2: the repeat tests of gate-command, verify-frame and psb-cli-contract. Gaps: (1) The identity and dates of the setup commits are not decided, so an output that holds a setup SHA differs per run. (2) No item says which environment a gate command gets (GOFLAGS, GOTOOLCHAIN), which can change a verdict. The smart-if part is phase 3 (later-p3-smartif).
- **NFR-006** (5 items): Item 1 (LAYUP's own pin) holds today; no item is needed. Item 2: setup-s02-s03-baseline, setup-s04-pin, verify-pin. Gaps: the meaning of pin.time, and so of the date key, is not given; no rule covers a crash inside S02, after which a rerun resolves the commit again (setup.md:57-60).
- **NFR-007** (5 items): Rules 1 to 5 and the table: found-boundary-test. Gaps: (1) The verb list of internal/git (packages.md:39) is too short; found-git must change it. (2) The rules do not say if test files are bound. (3) A shared test-fixture package has no row. (4) The §12 Test cell 'CI job lint (gofmt, go vet)' does not prove the criterion; found-boundary-test must fill it (gov-test-traceability).

### Sections of `docs/spec/` that no item names

- `docs/spec/README.md` ## The files: no item needed. An index of the spec files and their requirements. No build work. The specification task of each later milestone (gov-milestone-spec) adds its rows.
- `docs/spec/README.md` ## How a section is written: no item needed. Rules for each change to docs/spec: a 'decided here' value with its reason, or a marker with an open-gaps.tsv row. They bind the items that change the spec (psb-rule-edges, psb-utf8-input, the verb row of found-git, the S12 'gates' fix) and gov-milestone-spec. No build work.
- `docs/spec/gate.md` # `layup gate` and the gate manifest: no item needed. A title with two sentences that repeat REQ-004. The REQ-004 items cover them.
- `docs/spec/packages.md` # The packages: no item needed. A title. The table under NFR-007 is the work (found-boundary-test).
- `docs/spec/psb-check.md` # `layup psb check`: no item needed. A title. It names TestGoldenRealPSB, which psb-tsv and psb-cli-contract keep green.
- `docs/spec/records.md` # The records: no item needed. A title.
- `docs/spec/records.md` ### On the records branch (`records:`): no item needed. No spec_ref names this heading; items name its rows through 'The layout of the records branch'. Each phase-1 row has an item: README.md and the records commit (setup-s15-records-commit), setup/record.tsv (setup-runner), setup/verify.tsv (verify-frame), rule-paths.tsv (setup-rule-paths), telemetry.tsv (found-telemetry-schema), stalls.tsv (found-stalls-schema). The fixed text of README.md has no source (an open question of setup-s15-records-commit).
- `docs/spec/setup.md` # `layup setup`, `layup setup verify` and the stack catalog: no item needed. A title.

### Conflicts

Each conflict K*n* gets one host in step 3.

1. **K1.** Build order. #29 orders the work: gap check, setup, gates, telemetry. gov-rescope-29 proposes setup and gate first, then verify. The inventories need another order in two places. Each setup step needs its verify check first (packages.md:53-56; setup-s04-pin to setup-s12-stack-files). S12 needs gate-run, the Go entry and check gate:<kind> (setup.md:97, :256). So the order is: foundations; psb, the gate core and the verify frame; the setup steps with their checks; the Go entry; then S12 to S15. gov-rescope-29 must write this order. *Sources:* .ctx/issues/29.md:7; docs/adr/0012-build-layup-in-bootstrap-mode.md:48; docs/spec/packages.md:44; docs/spec/packages.md:53-56; docs/spec/setup.md:97; docs/spec/setup.md:256; gov-rescope-29; setup-s04-pin; setup-s12-stack-files
2. **K2.** Task count. ADR-0012:48 says that each engine deliverable (layup setup, layup gate, the telemetry record, the stall record) is one task under #29, and it does not name layup setup verify. The inventories cut phase 1 into about 60 items. The foundation items (found-tsv, found-git, found-boundary-test, found-cli, found-e2e-harness) serve no deliverable by name, but bootstrap rule 1 needs each task to name its In-Scope fact. Make each issue a child of one deliverable task, with its fact, or get an Operator decision. *Sources:* docs/adr/0012-build-layup-in-bootstrap-mode.md:48; docs/engineering-discipline.md:190-193; docs/issue-workflow.md:106-134; gov-rescope-29
3. **K3.** The PDR hold. #42 stops the build until the PDR closes, so every build item waits for gov-pdr-approval; the layers apply this edge. gov-test-traceability depends on found-boundary-test, a build item, and this makes a cycle. gov-test-traceability needs only the name of that test. *Sources:* .ctx/issues/42.md:5; .ctx/issues/42.md:20; gov-test-traceability; found-boundary-test; gov-pdr-approval
4. **K4.** Edges that items need but do not list; the layers add them: setup-s01-questions and setup-s06-facts -> psb-batch-api; setup-s10-markers -> verify-markers; setup-s12-stack-files -> gate-catalog-go-workflow; verify-jobs -> gate-catalog-go-workflow; gate-command -> found-cli; psb-cli-contract -> found-e2e-harness and found-cli; psb-tsv -> found-spec-schema-test; setup-e2e-offline -> verify-adapted and verify-test-baseline; verify-acceptance, found-nfr001-tests and found-nfr002-tests -> setup-e2e-offline. The layers remove two edges: gate-command -> setup-steps and gate-catalog-go-workflow -> setup-s12. *Sources:* setup-s01-questions; setup-s06-facts; setup-s10-markers; setup-s12-stack-files; verify-jobs; gate-command; psb-cli-contract; psb-tsv; setup-e2e-offline; verify-acceptance; found-nfr001-tests; found-nfr002-tests; gate-catalog-go-workflow
5. **K5.** Three edges name setup-runner but mean later steps: verify-acceptance ('layup setup WORK through S15'), gov-pilot-inputs (the stop tables of S05 and S10) and gov-first-pilot (a full setup). setup-runner alone gives none of these. The layers use setup-e2e-offline, setup-s10-markers and setup-s05-history. *Sources:* verify-acceptance; gov-pilot-inputs; gov-first-pilot; docs/spec/setup.md:207-208
6. **K6.** Test-only cycles. The e2e test of verify-gate-fixtures needs a WORK through S12, and the integration test of verify-jobs reads the workflow 'as S12 writes it'. But setup-s12-stack-files needs both checks first. Make those WORK trees by hand, read the catalog file directly, or move these scenarios to setup-e2e-offline. *Sources:* verify-gate-fixtures; verify-jobs; setup-s12-stack-files; docs/spec/setup.md:97
7. **K7.** The verbs of internal/git. packages.md:39 lists clone, ls-remote, init, commit, rev-parse, worktree, show, diff --name-only and apply. The items also need checkout (S02), add (S03 and the harness), a branch (S04), ls-files (S10, checks markers and adapted), rev-list (check pin), an orphan commit (S15), a commit on no ref (gate:<kind>), clone --depth 1 (the harness) and a remote (commands.sh). found-git must change the row of packages.md in the same PR. *Sources:* docs/spec/packages.md:39; docs/spec/setup.md:87-89; docs/spec/setup.md:213-217; docs/spec/setup.md:232; docs/setup/setup-check.sh:307; found-git; setup-s15-records-commit; verify-pin; verify-markers; verify-harness; verify-gate-fixtures
8. **K8.** Test files and the package rules. packages.md does not say if rules 3 and 4 and the may-import column bind _test.go files. cmd/layup/main_e2e_test.go already imports os/exec. The e2e harness, the fixture builders and the catalog fixture test start git from test code, and found-nfr001-tests needs internal/tsv from a cmd/layup test. found-boundary-test reads only non-test imports. Make one decision before found-boundary-test lands. *Sources:* docs/spec/packages.md:18-21; docs/spec/packages.md:36; cmd/layup/main_e2e_test.go:6; found-boundary-test; found-e2e-harness; found-nfr001-tests; verify-harness; setup-e2e-offline; gate-catalog-go-entry
9. **K9.** Schemas that two packages share. internal/setup and internal/verify both read or write the blocks setup-answers, setup-record, setup-verify and open-gaps. Neither package may import the other, and internal/tsv holds no named schema. The Go values of these schemas need one home, or two copies that the block test pins. found-spec-schema-test leaves the home of its registry open. *Sources:* docs/spec/packages.md:38; docs/spec/packages.md:43-44; docs/spec/setup.md:142; docs/spec/setup.md:155; docs/spec/setup.md:190; docs/spec/setup.md:212-213; docs/spec/setup.md:268; setup-answers; setup-s15-records-commit; verify-frame; verify-markers; found-spec-schema-test
10. **K10.** S05 and link-lint. S05 stops with one F-<path> row for each file whose links the deletion breaks. internal/setup starts no program. verify-baseline-scripts gives link-lint only as 'exit <code>', with the script output on stderr. So no item gives S05 the list of files. S05 needs its own link parser, or a verify call that returns the files. *Sources:* docs/spec/setup.md:90; docs/spec/packages.md:43; setup-s05-history; verify-baseline-scripts
11. **K11.** The stand-in baseline. found-e2e-harness makes fixture repositories from a directory under root tests/. setup-e2e-offline and verify-test-baseline make the tree at run time. tests/ is not in MK_EXEMPT or AD_EXCLUDE, and link-lint reads every tracked .md file except paths with good/ or bad-*/. So a checked-in baseline with markers, kit words or broken links fails LAYUP's own checks. Also, one builder for internal/verify tests and cmd/layup tests must be a package with a row in packages.md, which found-boundary-test requires for each package. *Sources:* docs/setup/setup-check.sh:305; docs/setup/setup-check.sh:359; docs/links/link-lint.sh:366-370; tests/README.md; found-e2e-harness; setup-e2e-offline; verify-test-baseline; found-boundary-test
12. **K12.** The real baseline cannot pass. Check adapted fails on kept files that no step writes (the setup and verify areas measured 235 failures in 21 files). adapted is the evidence of no step, so it fails first at S15. Check kit-history also fails after S05. Measured on d2516fd: completed.md has 56 lines that link github.com/pharzam/armature/, and 24 of them name one of the 19 deleted task IDs, so 32 lines stay; backlog.md keeps all 7. setup-s05-history says 56 for completed.md; the correct count is 32. gov-first-pilot plans a pilot that cannot pass. A spec fix needs an item before the pilot. *Sources:* docs/spec/setup.md:90; docs/spec/setup.md:212-213; docs/spec/setup.md:245; docs/spec/setup.md:251; docs/setup/setup-check.sh:106-109; git show d2516fd:docs/tasks/completed.md; setup-s05-history; setup-s15-records-commit; verify-adapted; gov-first-pilot
13. **K13.** The evidence of S12 is 'checks jobs and gates'. The verify table has no check gates, only one gate:<kind> row per kind. setup-s12-stack-files and verify-frame read 'gates' as every gate:<kind> row, but no source says so. *Sources:* docs/spec/setup.md:97; docs/spec/setup.md:255-256; setup-s12-stack-files; verify-frame
14. **K14.** PRD-0001:115 asks that the answers of each human_decision row are in the target's Git before the next step. The step order puts the S01 answers into Git at S06, after S02 to S05, and the S10 answers at S11. The pilot UAT cannot pass as written. *Sources:* docs/prd/PRD-0001-layup.md:115; docs/spec/setup.md:91; docs/spec/setup.md:96; gov-first-pilot; setup-pilot-acceptance
15. **K15.** records.md says that each answer in the target names the comment that holds it (source), and that no value stays only in the work area. setup.md:91 gives each answer fact only the question ID, the question text and the answer. So the source and by columns of answers.tsv stay only on the host. *Sources:* docs/spec/records.md:105-110; docs/spec/setup.md:91; docs/spec/setup.md:142-148; setup-s06-facts; setup-s11-fill; verify-facts; found-nfr001-tests
16. **K16.** S14 runs after S10 and S11. The baseline's AGENTS.md holds the marker '&lsaquo;worktree dir&rsaquo;' (line 104 at d2516fd). S10 asks it, S11 fills it or keeps it with an open-gaps row, and then S14 overwrites the file. This leaves a stale open-gaps row or a record row for a value that is not in the tree, so check markers or check sources fails. gov-pilot-inputs must also give S14 files with no marker. *Sources:* docs/spec/setup.md:95-99; git show d2516fd:AGENTS.md (line 104); setup-prose-inputs; setup-s10-markers; verify-markers; verify-sources; gov-pilot-inputs
17. **K17.** Markers inside docs/facts/. S11 writes the M- answers as a raw facts record, and its text can quote each marker with its angle quotes. A brief that S06 copies can also hold '&lsaquo;'. MK_EXEMPT does not exempt docs/facts/, so check markers, which is the evidence of S11, fails. *Sources:* docs/spec/setup.md:91; docs/spec/setup.md:96; docs/spec/setup.md:250; docs/setup/setup-check.sh:305; setup-s06-facts; setup-s11-fill; verify-markers
18. **K18.** Check facts as the evidence of S06. The rule of check facts puts each M- ID in the record of S11, which does not exist at S06. An M- row that the Operator gives early then fails a correct S06. verify-frame also reads a partial record.tsv when a check runs as step evidence. *Sources:* docs/spec/setup.md:91; docs/spec/setup.md:246; setup-s06-facts; verify-facts; verify-frame
19. **K19.** A pending kind with a missing tool. The result table decides 'pending' before it looks at the tool, so the result is clear. NFR-004 item 3 says that a missing toolchain is never clear. gate-results and verify-gate-fixtures both follow the table. *Sources:* docs/spec/gate.md:67-74; docs/spec/gate.md:140-141; gate-results; verify-gate-fixtures
20. **K20.** The state of the test kind at setup. gate-catalog-go-entry makes it active. W-04 step 1 names only the static kind as active. One reading is necessary before the Go entry and its fixtures. *Sources:* docs/architecture.md:451-466; docs/walkthroughs/W-04-stack-dependent-gates.md:17; gate-catalog-go-entry; verify-gate-fixtures
21. **K21.** Phase-1 criteria that need later phases. O-115 keeps the §7.1 criteria as they are. REQ-001 needs the Intake batch and the Early Question Share. REQ-004 needs a seeded violation of each of four kinds on each pilot stack, but three kinds stay pending with no fixture and only Go exists. REQ-007 needs PRs in the pilot. REQ-009 and REQ-011 need a writer. NFR-001 and NFR-002 need a pilot task. gov-pilot-definition must say what phase 1 is accepted against. *Sources:* docs/tasks/T-0drh.md:22-23; docs/prd/PRD-0001-layup.md:114-133; docs/spec/setup.md:296; docs/spec/setup.md:312; gov-pilot-definition; gov-first-pilot; gate-catalog-go-entry
22. **K22.** Two meanings of 'pilot'. ADR-0012:49 ends bootstrap mode when one target is set up and its gate runs from outside, which is possible after phase 1. PRD-0001 puts the pilot baseline and the pilot on two problem statements in phase 4. gov-first-pilot, later-p4-pilot and setup-pilot-acceptance use different meanings. *Sources:* docs/adr/0012-build-layup-in-bootstrap-mode.md:49; docs/prd/PRD-0001-layup.md:144-145; docs/prd/PRD-0001-layup.md:182; gov-pilot-definition; gov-first-pilot; later-p4-pilot
23. **K23.** The home of the telemetry, price and stall types. PRD-0001:179 and the schema-block rule need phase-1 Go code for these schemas. packages.md gives internal/ledger phase 2 and internal/stall phase 3, and has no phase-1 row for them, although packages.md:70-72 adds a row with the requirement's section, which exists. later-p2-sessions-ledger and later-p3-stalls-budget then need the same types. *Sources:* docs/prd/PRD-0001-layup.md:179; docs/spec/README.md:75-77; docs/spec/packages.md:70-72; docs/spec/packages.md:85; docs/spec/packages.md:88; found-telemetry-schema; found-stalls-schema; later-p2-sessions-ledger; later-p3-stalls-budget
24. **K24.** The coverage floor. architecture.md says that the setup record lists the floor as an open gap. No step writes it: S10 and S11 run before S12 writes the catalog files. A marker in internal/catalog/ fails LAYUP's own check markers, because MK_EXEMPT does not exempt that path. *Sources:* docs/architecture.md:1548-1551; docs/spec/setup.md:95-97; docs/setup/setup-check.sh:305; gate-catalog-go-entry; setup-s12-stack-files; verify-sources
25. **K25.** gov-pilot-numbers. Its rule must be written before the pilot runs, but the item depends on gov-first-pilot. Split it: the rule before gov-first-pilot, the count after it. *Sources:* docs/guardrails.md:21-26; gov-pilot-numbers; gov-first-pilot
26. **K26.** #33 says that the Go gates run outside the target, in LAYUP's own runner (ADR-0011 decision 4, O-11). ADR-0016 puts the native gates in the target and runs them from outside with layup gate. gov-rescope-33 must close or rewrite #33 before a gate issue starts. *Sources:* .ctx/issues/33.md:5; docs/adr/0016-put-the-native-stack-gates-in-the-target.md:16-18; gov-rescope-33; gate-run
27. **K27.** The e2e time limit. All e2e tests of cmd/layup run in one test binary with -timeout 10m. setup-e2e-offline, verify-acceptance, found-nfr001-tests, found-nfr002-tests, gate-command, gate-catalog-go-entry and psb-cli-contract all add to it. Each full setup runs layup gate twice per active kind, and the test kind runs go test on the target. Share one full run, or plan the time. *Sources:* .github/workflows/ci.yml:318; docs/tests/test-levels.md:71; setup-e2e-offline; verify-acceptance; found-nfr001-tests; found-nfr002-tests; gate-catalog-go-entry
28. **K28.** gov-operator-ci-protection applies only when phase 1 adds a CI job. gate-catalog-go-entry runs the fixture test in the existing job tests, so ci.yml does not change and the item is void; the layers leave it out. verify-test-baseline offers fetch-depth: 0 in job tests as an option. That is a CI change (cycle cap 2), against the gate inventory's choice to keep ci.yml unchanged. *Sources:* gov-operator-ci-protection; gate-catalog-go-entry; verify-test-baseline; .github/workflows/ci.yml:306-318; docs/engineering-discipline.md:202-203
29. **K29.** Progress lines. The discipline asks a progress indicator for any operation over ten seconds. verify-baseline-scripts and verify-gate-fixtures print progress on stderr. gate-results, gate-command and setup-s02-s03-baseline (the clone) do not. found-cli leaves open whether the rule binds the product commands. One rule for all commands is necessary. *Sources:* docs/engineering-discipline.md:832-834; docs/spec/gate.md:91-92; found-cli; gate-results; setup-s02-s03-baseline; verify-baseline-scripts; verify-gate-fixtures
30. **K30.** Job name or job id. gate.md and setup.md name each CI job 'as the kind'. check_protection reads the name: of a job, else its id. The ruleset of S13 requires check contexts, and GitHub takes a context from name: when it is set. gate-catalog-go-workflow, setup-s13-rulesets and verify-jobs must use one reading, or a required check never reports. *Sources:* docs/spec/gate.md:19; docs/spec/setup.md:97-98; docs/spec/setup.md:255; docs/setup/setup-check.sh:538-550; gate-catalog-go-workflow; setup-s13-rulesets; verify-jobs
31. **K31.** Git configuration and credentials. found-git and setup-runner isolate git from the host configuration. S02 clones the baseline, and a private baseline needs the Operator's credential helper, which the isolation removes. No source says how S02 gets a credential with no prompt. *Sources:* docs/spec/setup.md:87; docs/spec/README.md:44-45; found-git; setup-runner; setup-s02-s03-baseline
32. **K32.** A file that is not UTF-8. psb-utf8-input chooses exit 2 or U+FFFD for layup psb check. found-tsv leaves open what its writer does with invalid UTF-8, and its reader rejects it. Make one decision; psb-tsv then writes through found-tsv with that rule. *Sources:* docs/spec/psb-check.md:21; docs/spec/README.md:55; docs/spec/README.md:67; psb-utf8-input; found-tsv; psb-tsv
33. **K33.** The stop table hides the line of a gap. setup.md:114 gives 'where' as the dash for a Q-NNN question, and G4 gives the same question for each line with the same word. Two G4 rows then look the same to the idea owner, who answers by ID at S01. *Sources:* docs/spec/setup.md:114; docs/spec/psb-check.md:39; psb-batch-api; setup-s01-questions
34. **K34.** layup version with extra words. README.md:42-43 exempts layup version from the command rules, and today 'layup version extra' exits 0. found-cli plans 'an extra argument gives exit 2'. psb-cli-contract leaves this open. *Sources:* docs/spec/README.md:42-43; internal/cli/cli.go:32-34; found-cli; psb-cli-contract
35. **K35.** The embed rules and the catalog form. setup.md puts each target file at internal/catalog/<stack>/files/<path>, so files/go.mod and files/.github/. The gate area measured that //go:embed refuses a tree with a go.mod and leaves out .github/ without the all: prefix. The catalog form needs a rename rule (for example go.mod.tmpl) as a 'decided here' value. gate-catalog-package and setup-s12-stack-files must use the same rule. *Sources:* docs/spec/setup.md:276-286; gate-catalog-package; gate-catalog-go-entry; setup-s12-stack-files
36. **K36.** Batched issues with no host item. ADR-0012:60 puts #15, #21, #24, #34 and #48 on product tasks that touch the same files. #21, #34, #48 and #61 have hosts (see duplicates). gov-issue-15 (a guardrails §2 example) and gov-issue-24 (one sentence of docs/ci/README.md) have no host in phase 1, because no phase-1 item edits those files. The layers leave both out. #24 can close with its evidence in the forge edits of gov-rescope-29. *Sources:* docs/adr/0012-build-layup-in-bootstrap-mode.md:60; gov-issue-15; gov-issue-24; gov-rescope-29
37. **K37.** Items outside phase 1 that have no later- prefix: gov-milestone-spec, gov-issue-68, gov-issue-49 and gov-operator-setup-o112. A filter by prefix puts them in phase 1. The layers leave them out. *Sources:* gov-milestone-spec; gov-issue-68; gov-issue-49; gov-operator-setup-o112
38. **K38.** Later phases: the phase of some records. questions.tsv is phase 3, but Intake writes its rows in phase 2. routing.tsv is phase 3, but internal/route with routing is phase 2. budget.tsv is phase 3, but internal/ledger holds the budget in phase 2. The frozen test list is phase 3, but the phase-2 transition checks the plan's tests. The first notice of the dead-man job is phase 2, but its record is phase 3. Phase 2 needs a second admitted harness, but REQ-013 is phase 4. Settle these before the phase-2 specification. *Sources:* docs/spec/records.md:37; docs/architecture.md:373-375; docs/spec/records.md:53; docs/spec/packages.md:89; docs/spec/records.md:60; docs/spec/packages.md:85; docs/spec/records.md:47; docs/architecture.md:776-777; docs/architecture.md:830; docs/spec/setup.md:31; docs/spec/records.md:88; docs/architecture.md:606-608; docs/prd/PRD-0001-layup.md:182; later-p2-intake-spec; later-p2-sessions-ledger; later-p2-build-accept; later-p3-stalls-budget
39. **K39.** packages.md:24-25 says 'the one later exception' and names two packages: internal/smartif and the forge adapter. found-boundary-test and later-p3-smartif need the exact list of packages that may import net. *Sources:* docs/spec/packages.md:22-25; found-boundary-test; later-p3-smartif; later-p2-run-start
40. **K40.** Inputs of layup run. README.md:44-45 forbids an environment variable as an input. The architecture gives a session its harness credential as a variable or a file, and the Start command names no flag for the App key. later-p2-run-start needs a rule. *Sources:* docs/spec/README.md:44-45; docs/architecture.md:241-243; docs/architecture.md:297-300; later-p2-run-start
41. **K41.** A target that phase 1 set up. setup.md says that the Operator applies the later rulesets with commands that layup run prints when it first starts on the target. The architecture starts layup run only on an empty repository (--new), and the records branch of a phase-1 target already holds the Operator's commit. later-p2-scaffold has no adoption path. *Sources:* docs/spec/setup.md:309-311; docs/architecture.md:297-300; docs/spec/records.md:99-103; later-p2-scaffold; later-p2-run-start
