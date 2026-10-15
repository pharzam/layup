# T-57vx — the runs of the tests

Times from the clock of the host, UTC, 2026-10-15.

## Step T: the coverage check and the issue check, red

- `sh runs/T-57vx/parts.sh` before the table of the plan existed (the parts table was written first): exit 1, 259 lines of FAIL: the plan had no table "The tasks of M2e", each of the 24 rows of `parts.tsv` was not a row of it, and each part was named in no Task cell. Rule 1 (each heading of `loop.md` is in the list) held; rules 4 and 5 (each acceptance row, and each clause that the slicing cuts, has its entry) held; rules 3 and 6 failed.
- `sh runs/T-57vx/issues.sh` before the issues were open (it reads the forge, so it runs by hand): exit 1, `FAIL FAIL row 60 of The tasks of M2e does not link #347 with T-fn5s`.

## After D1 and D3

- `sh runs/T-57vx/parts.sh`: exit 0, "each part is in one row of The tasks of M2e, or named with its reason".
- `sh runs/T-57vx/issues.sh`: exit 0, 24 lines `ok`: each row links its issue (#347 to #370) and the App's bot opened each.

## The checks of acceptance criterion 3, before the freeze

`sh docs/prd/prd-lint.sh`, `sh docs/links/link-lint.sh`, `sh docs/setup/setup-check.sh`, `sh docs/tests/run-discipline-tests.sh` and `sh docs/adr/adr-lint.sh`: exit 0 each. `git diff --check`: exit 0.
