# T-bng8 — test runs

The build tasks of `M2d` (#267), step T of the plan, with condition 2 of its
review: the coverage check `parts.sh`, red before the table of the plan, green
after.

## Red (2026-10-12, 06:14Z)

`sh runs/T-bng8/parts.sh`, before "The tasks of M2d" existed: exit 1. The other
checks pass: each heading, each acceptance test, each input state and each part
outside `scaffold.md` already has its row. Its 16 lines:

```
FAIL docs/plan/README.md has no table The tasks of M2d
FAIL row 50 is not a row of The tasks of M2d
FAIL row 51a is not a row of The tasks of M2d
FAIL row 51b is not a row of The tasks of M2d
FAIL row 52a is not a row of The tasks of M2d
FAIL row 52b is not a row of The tasks of M2d
FAIL row 52c is not a row of The tasks of M2d
FAIL row 53 is not a row of The tasks of M2d
FAIL row 54 is not a row of The tasks of M2d
FAIL row 55 is not a row of The tasks of M2d
FAIL row 56 is not a row of The tasks of M2d
FAIL row 57a is not a row of The tasks of M2d
FAIL row 57b is not a row of The tasks of M2d
FAIL row 58 is not a row of The tasks of M2d
FAIL row 59a is not a row of The tasks of M2d
FAIL row 59b is not a row of The tasks of M2d
```

## Green (2026-10-12, 06:18Z)

After D1 (the table "The tasks of M2d", 15 rows from `gen-issues.py --table`)
and D3 (the issues #270 to #284, opened by the App with the duplicate check of
each in its body):

- `sh runs/T-bng8/parts.sh`: exit 0, "ok   each part is in one row of The tasks
  of M2d, or named with its reason".
- `sh runs/T-bng8/issues.sh`: exit 1 before the issues ("FAIL
  runs/T-bng8/issues.tsv is missing: no issue is open"), and exit 0 after, with 15
  lines `ok   row <n>: #<issue> (<task>) by layup-agent[bot]`.
- `link-lint` (2,318 links), `prd-lint`, `setup-check` and `run-discipline-tests`
  (102 passed): exit 0.

## The fix of round 1 (2026-10-12, 06:26Z)

Finding 1 of round 1: `parts.sh` compared `parts.tsv` with the plan by row
number only. Check 6 adds that each listed part is named in the Task cell of its
own row, in the form that `gen-issues.py` writes it. The plan's path can come from
`PLAN`, so a mutant of the plan is checked:

- The first form of check 6 read the wrong field and failed on each part: a
  defect of the check, fixed before any record.
- With `PLAN` a copy of the plan whose row 57a drops "`scaffold.md`: The
  read-back": exit 1, `FAIL row 57a of The tasks of M2d does not name its part:
  \`scaffold.md\`: The read-back`.
- On the real plan: exit 0, "ok   each part is in one row of The tasks of M2d,
  or named with its reason".

## The fix of round 2 (2026-10-12, 12:40Z)

Findings 1 and 2 of round 2, and the note:

- Finding 1: the table of `M2d` in `packages.md` sat at `none`, reasoning that
  it adds no package while its changed cells had no owning Task cell. The part
  is now 13 clause parts, each in the row whose code changes the cell (the
  steps to 53 to 57a, the setup cells to 52a and 52b, the forge calls to 51a,
  the blocks to 50, `CommitTree` to 51b, the imports and the starts-no-program
  cell to 54). The Task cells and the bodies of issues #270 to #280 come from
  the same `parts.tsv` through `gen-issues.py`.
- Finding 2: check 5 asked for any `run-steps:` entry, so a clause dropped from
  the map and its row still passed. Check 5 now names each clause of each built
  block changed in words and each cell of the table of `M2d`, one `want` per
  clause.
- The note: `gen-issues.py --edit` re-ran the duplicate search with the issue
  already open, so the search found the issue itself and the words "before this
  issue opened" misdated the check. `dup()` now takes the issue number, drops
  it from the found set, and says "at this edit, this issue itself excluded".

Red before the fix — the mutant of finding 2, `run-steps: verify` and
`internal/run: the step verify` dropped from `parts.tsv` and row 55's Task
cell, against the old check: exit 0. Against the fixed check:

```
FAIL no part: docs/spec/records.md	## NFR-001 — The records of the Scaffold	run-steps: verify
FAIL no part: docs/spec/packages.md	### The table of M2d	internal/run: the step verify
```

Green after:

- `sh runs/T-bng8/parts.sh`: exit 0, "ok   each part is in one row of The tasks
  of M2d, or named with its reason".
- `go build ./...`, `go vet ./...`, `go test ./...`, `go test
  -tags=integration ./...`, `gofmt -l .` (empty), `prd-lint`, `link-lint`
  (2,318 links), `setup-check`, `run-discipline-tests` (102 passed),
  `adr-lint`, `git diff --check`: exit 0.
- `gen-issues.py --edit 50 51a 51b 52a 52b 53 54 55 56 57a` wrote the bodies of
  issues #270 to #280 with their new parts and the corrected duplicate check.

## The fix of round 3 (2026-10-12, 14:20Z)

Findings 1 and 2 of round 3:

- Finding 1: `setup/rulesets.txt` and `setup/inputs/files` were one `none`
  entry reasoning "a rule, no block; written by rows 56 and 54" — the reason
  named the rows but no Task cell named either part. They are now two parts,
  `setup/inputs/files` in row 54 and `setup/rulesets.txt` in row 56, each in
  the row that writes the file; the Task cells and the bodies of #277 and #279
  come from the same `parts.tsv`.
- Finding 2: "exactly one" was checked only against the mapped row: a copy of
  a part in a second row's Task cell passed. Check 6 now also counts the Task
  cells that name each part and fails unless the count is one.

Red before the fix — the mutant of finding 2, "`scaffold.md`: The step
markers" copied into row 54's Task cell, against the old check: exit 0.
Against the fixed check:

```
FAIL the part is named in 2 Task cells of The tasks of M2d, not one: `scaffold.md`: The step markers
```

Green after:

- `sh runs/T-bng8/parts.sh`: exit 0, "ok   each part is in one row of The tasks
  of M2d, or named with its reason".
- `go build ./...`, `go vet ./...`, `go test ./...`, `go test
  -tags=integration ./...`, `gofmt -l .` (empty), `prd-lint`, `link-lint`
  (2,318 links), `setup-check`, `run-discipline-tests` (102 passed),
  `adr-lint`, `git diff --check`: exit 0.
- `gen-issues.py --edit 54 56` wrote the bodies of issues #277 and #279 with
  their new parts.
