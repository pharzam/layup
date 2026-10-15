# T-57vx — the runs of the tests

Times from the clock of the host, UTC, 2026-10-15; the runs of step T are from before 16:35, the mutant runs and the re-runs from after round 1.

## Step T: the coverage check and the issue check, red

- `sh runs/T-57vx/parts.sh` before the table of the plan existed (the parts table was written first): exit 1, 259 lines of FAIL: the plan had no table "The tasks of M2e", each of the 24 rows of `parts.tsv` was not a row of it, and each part was named in no Task cell. Rule 1 (each heading of `loop.md` is in the list) held; rules 4 and 5 (each acceptance row, and each clause that the slicing cuts, has its entry) held; rules 3 and 6 failed.
- `sh runs/T-57vx/issues.sh` with no `issues.tsv` (re-run after round 1, as the file moved aside): exit 1, `FAIL runs/T-57vx/issues.tsv is missing: no issue is open`. (The first version of this record quoted a run that named #347; that run came after the issues opened and before the table linked them, and it is not a red run of the check.)
- The branch has one commit before round 1 (`f17f880`): the order of step T is the order of the files that this record names, not a history of commits.

## After D1 and D3

- `sh runs/T-57vx/parts.sh`: exit 0, "each part is in one row of The tasks of M2e, or named with its reason".
- `sh runs/T-57vx/issues.sh`: exit 0, 24 lines `ok`: each row links its issue (#347 to #370) and the App's bot opened each.

## The checks of acceptance criterion 3, before the freeze

`sh docs/prd/prd-lint.sh`, `sh docs/links/link-lint.sh`, `sh docs/setup/setup-check.sh`, `sh docs/tests/run-discipline-tests.sh` and `sh docs/adr/adr-lint.sh`: exit 0 each. `git diff --check`: exit 0.

## Mutant runs of `parts.sh`, after round 1

Each on a copy of the table or of the plan, `parts.sh` exit 1 with the line shown: the forge part dropped (`FAIL no part: docs/spec/forge.md	## The calls of M2e	—`); a clause dropped (`FAIL no part: ... ### The revision of the draft	the check of the quote's fact`); a heading dropped (`FAIL a heading of docs/spec/loop.md has no part: ### The task register`); a part listed twice (`FAIL a part is listed twice`); a part copied into a second Task cell (`FAIL the part is named in 2 Task cells of The tasks of M2e, not one`); a row renamed (`FAIL row 66a is not a row of The tasks of M2e`). After the fixes of round 1 the unchanged table: exit 0; `issues.sh`: exit 0, 24 lines `ok`.

## After round 2 (the fixes that wait for the Operator's answer), 2026-10-15T17:09Z

Round 2 ended `not mergeable, findings recorded` at the cap 1 (two material findings, five notes). The fixes are a local commit, not pushed until the Operator chooses: the After cell of 66d gains 62b; the end check of a `comment` row is row 66d's; five clauses are cut or moved (the Retrospective step of the `internal/run` bullet; the lost merge event split between 66g and 66k; the open question shown in the brief, 66j; the two-survivors clause, 66j); the Tests cell of 65b; the labels of the run-steps and of the layout path; `issues.sh` says D3. `parts.sh` exit 0; `issues.sh` exit 0 (24 `ok`); the five checks exit 0; `git diff --check` exit 0.
