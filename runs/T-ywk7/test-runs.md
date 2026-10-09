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

## S4: the map of the parts to their headings

Red, before `docs/spec/session.md` existed (2026-10-09, after the S3 comment
6075865294): `sh runs/T-ywk7/sections.sh` printed `FAIL` for each of the 28
headings of [`sections.tsv`](sections.tsv), and exited 1.

## S6: the schema blocks of the records of a session

Red (06:58Z), with the eight new blocks in `docs/spec/records.md` and no name in
either list: `go test -count=1 -tags=integration -run TestEverySchemaBlockIsBuiltOrNotYetBuilt ./internal/tsv/`

```
--- FAIL: TestEverySchemaBlockIsBuiltOrNotYetBuilt (0.00s)
    blocks_integration_test.go:41: the block events is listed 0 times in built and notYetBuilt; want 1
    blocks_integration_test.go:41: the block harnesses is listed 0 times in built and notYetBuilt; want 1
    blocks_integration_test.go:41: the block models is listed 0 times in built and notYetBuilt; want 1
    blocks_integration_test.go:41: the block probe-result is listed 0 times in built and notYetBuilt; want 1
    blocks_integration_test.go:41: the block result is listed 0 times in built and notYetBuilt; want 1
    blocks_integration_test.go:41: the block routing is listed 0 times in built and notYetBuilt; want 1
    blocks_integration_test.go:41: the block routing-register is listed 0 times in built and notYetBuilt; want 1
    blocks_integration_test.go:41: the block sessions is listed 0 times in built and notYetBuilt; want 1
FAIL
```

Each block parses, so the parser of `internal/tsv` takes its form. Green, with
the eight names in `notYetBuilt` of `internal/tsv/blocks_integration_test.go`:
`ok  github.com/pharzam/layup/internal/tsv`. `go test -count=1 -tags=integration ./cmd/layup/ ./internal/route/ ./internal/records/`
passes with the rows of `internal/route`, `internal/run` and `internal/records`
changed in place.

## Green, on the head before the freeze (2026-10-09)

- `sh runs/T-ywk7/sections.sh`: 28 `ok`, exit 0.
- `sh runs/T-ywk7/survey.sh`: 11 `ok` (7 `take` with their heading, 4 `reject`), exit 0.
- `go test -count=1 -tags=integration ./...`: each package `ok`; `go build ./...` and `go vet ./...` clean.
- `adr-lint`, `prd-lint`, `link-lint`, `run-discipline-tests` (81 passed),
  `nested-checkout-check`, `setup-check`: exit 0; `git diff --check`: clean.
- The order of S1: `git log --reverse --format='%h %s' origin/main..HEAD -- runs/T-ywk7/search docs/spec`

```
0a01743 docs: T-ywk7 the public-solution search of M2b: two blind searchers, the sources and the summary
8fd76b0 docs: T-ywk7 the specification of M2b: role sessions, the records of a session, the table of M2b
```
