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
