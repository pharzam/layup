# The test runs of T-fdaq

## D1 red, at base `faeb1dd` (2026-10-09T10:49Z)

`sh runs/T-fdaq/parts.sh`, before the table: only check 3 fails, one line per row of `parts.tsv`; the headings of `session.md`, the parts of `runs/T-ywk7/sections.tsv` and the rows of "The acceptance tests of M2b" each have their part.

```
FAIL row 28 is not a row of The tasks of M2b
FAIL row 29 is not a row of The tasks of M2b
FAIL row 30a is not a row of The tasks of M2b
FAIL row 30b is not a row of The tasks of M2b
FAIL row 31 is not a row of The tasks of M2b
FAIL row 32 is not a row of The tasks of M2b
FAIL row 33a is not a row of The tasks of M2b
FAIL row 33b is not a row of The tasks of M2b
FAIL row 34a is not a row of The tasks of M2b
FAIL row 34b is not a row of The tasks of M2b
FAIL row 35 is not a row of The tasks of M2b
FAIL row 36a is not a row of The tasks of M2b
FAIL row 36b is not a row of The tasks of M2b
FAIL row 37a is not a row of The tasks of M2b
FAIL row 37b is not a row of The tasks of M2b
FAIL row 38 is not a row of The tasks of M2b
FAIL row 39a is not a row of The tasks of M2b
FAIL row 39b is not a row of The tasks of M2b
exit 1
```

## D2 red, before the issues (2026-10-09T10:50Z)

```
FAIL runs/T-fdaq/issues.tsv is missing: no issue is open
exit 1
```

## Green, on `ec6e3cd` (2026-10-09T10:56Z)

```
ok   each part is in one row of The tasks of M2b, or named with its reason
ok   row 28: #153 (T-3py1) by layup-agent[bot]
ok   row 29: #154 (T-y10b) by layup-agent[bot]
ok   row 30a: #155 (T-ysph) by layup-agent[bot]
ok   row 30b: #156 (T-cht1) by layup-agent[bot]
ok   row 31: #157 (T-z5dj) by layup-agent[bot]
ok   row 32: #158 (T-m1dx) by layup-agent[bot]
ok   row 33a: #159 (T-vxdg) by layup-agent[bot]
ok   row 33b: #160 (T-6sbe) by layup-agent[bot]
ok   row 34a: #161 (T-5pxd) by layup-agent[bot]
ok   row 34b: #162 (T-bpxg) by layup-agent[bot]
ok   row 35: #163 (T-4c3q) by layup-agent[bot]
ok   row 36a: #164 (T-d8t9) by layup-agent[bot]
ok   row 36b: #165 (T-fsjp) by layup-agent[bot]
ok   row 37a: #166 (T-z027) by layup-agent[bot]
ok   row 37b: #167 (T-e3sy) by layup-agent[bot]
ok   row 38: #168 (T-nxe4) by layup-agent[bot]
ok   row 39a: #169 (T-4tjy) by layup-agent[bot]
ok   row 39b: #170 (T-x7cs) by layup-agent[bot]
```

The checks of the third acceptance criterion, each with exit 0: `adr-lint`, `prd-lint`, `link-lint` (1842 links resolved), `run-discipline-tests` (81 passed, 0 failed), `nested-checkout-check` (10 cases), `setup-check`, and `git diff --check`.
