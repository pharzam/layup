# T-5csx — test runs

The acceptance test of the merge into the batch and the comment (#287). The test is
the check `runs/T-bng8/parts.sh`: red when the specification gains a test that no row
of the plan delivers, green after the map, the plan and the issue of row 53 name it.

## Baseline (2026-10-12, 12:16Z)

On `45f546c`, `sh runs/T-bng8/parts.sh`: exit 0, "ok   each part is in one row of The
tasks of M2d, or named with its reason". `python3 -I runs/T-bng8/gen-issues.py --table`
prints 17 lines, and each is in `docs/plan/README.md` verbatim.

## Red (2026-10-12, 12:16Z)

The row "The setup questions in the batch and the comment" is added to "The acceptance
tests of M2d" in `docs/spec/scaffold.md`, and the map is not changed. `sh
runs/T-bng8/parts.sh`: exit 1, one line, from check 4; the other five checks pass.

```
FAIL a row of ## The acceptance tests of M2d has no part: The setup questions in the batch and the comment
```
