# T-2r8u — test runs

The slicing of `M2c` (#211): `parts.sh` (each part in one row) and `issues.sh`
(each row links its issue, opened by the App), red before the table and the
issues, green after.

## Red (2026-10-10, 20:31Z and 23:19Z)

`sh runs/T-2r8u/parts.sh` with an empty `parts.tsv` and no table, at 20:31Z:
exit 1, 55 `FAIL` lines, the first ones:

```
FAIL a heading of docs/spec/intake.md has no part: ## REQ-001 — The gap check of Intake, in one batch
FAIL a heading of docs/spec/intake.md has no part: ### The sessions of Intake
FAIL a heading of docs/spec/intake.md has no part: ### The budget of Intake
FAIL a heading of docs/spec/intake.md has no part: ### The rule gaps
```

`sh runs/T-2r8u/issues.sh` before the issues, at 23:19Z: exit 1,

```
FAIL runs/T-2r8u/issues.tsv is missing: no issue is open
```

## Green (2026-10-10, 23:21Z)

`sh runs/T-2r8u/parts.sh`: exit 0, "ok   each part is in one row of The tasks
of M2c, or named with its reason". `python3 -I runs/T-2r8u/gen-issues.py --post`
opened #212 to #228 as the App; `sh runs/T-2r8u/issues.sh` (it reads the forge):
exit 0, 17 rows `ok`, each "by layup-agent[bot]". Each with exit 0:
`git diff --check b54748b`, `adr-lint`, `prd-lint`, `link-lint`, `setup-check`,
`run-discipline-tests`.
