# The test runs of T-ywk7

## S2: the survey table (2026-10-09)

Red, with no table (06:33Z): `sh runs/T-ywk7/survey.sh`

```
FAIL the table runs/T-ywk7/survey-m2b.tsv does not exist
exit 1
```

With the table, before the specification (06:58Z): each of the 11 rows has one
decision, matched by its row number and its first words; each `take` fails, as
its heading is not in `docs/spec/` yet.

```
FAIL row 1 (A process runner: one): take, but docs/spec/session.md has no line "### The limit of a session"
FAIL row 2 (Failure classes: "fails the): take, but docs/spec/session.md has no line "### The end of a session"
FAIL row 3 (The environment cleaned: the): take, but docs/spec/session.md has no line "### The environment and the harness credential"
ok   row 4 (The tool policy checked): reject
FAIL row 5 (The session ID chosen): take, but docs/spec/session.md has no line "### The start of a session"
FAIL row 6 (The launch of `claude): take, but docs/spec/session.md has no line "### The start of a session"
ok   row 7 (The argument lists of): reject
ok   row 8 (The headless flags of): reject
ok   row 9 (A token ledger from): reject
FAIL row 10 (A price table of): take, but docs/spec/session.md has no line "## REQ-011 — The writer of the telemetry record"
FAIL row 11 (Receipts of one line): take, but docs/spec/session.md has no line "## REQ-011 — The writer of the telemetry record"
exit 1
```
