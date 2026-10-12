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

## Green (2026-10-12, 12:18Z)

After the line in `runs/T-bng8/parts.tsv` (the test, as a part of row 53) and the
table of the plan from `python3 -I runs/T-bng8/gen-issues.py --table`:

- `--table` prints 17 lines, and each is in `docs/plan/README.md` verbatim but row
  53. So D1 changes the Task cell of row 53 alone, which gains "`scaffold.md`: The
  acceptance tests of M2d (The setup questions in the batch and the comment)".
- `sh runs/T-bng8/parts.sh`: exit 0, "ok   each part is in one row of The tasks of
  M2d, or named with its reason".

## The issue of row 53 (2026-10-12, 12:18Z)

`python3 -I runs/T-bng8/gen-issues.py --edit 53` wrote the body of #276 again, as
the App. Its new body was printed first and compared with the live one: one line
differs, "The parts of the specification", which gains the new test; the duplicate
check is the same. After it, the live body equals the printed one and names the test.

## The mutants (2026-10-12, 12:19Z)

On copies of the tree, never on the worktree:

- The new line dropped from the map: exit 1, the line of step T again, since check
  6 reads only the map.

  ```
  FAIL a row of ## The acceptance tests of M2d has no part: The setup questions in the batch and the comment
  ```

- The part dropped from the Task cell of row 53 in the plan: exit 1, two lines of
  check 6.

  ```
  FAIL row 53 of The tasks of M2d does not name its part: `scaffold.md`: The acceptance tests of M2d (The setup questions in the batch and the comment)
  FAIL the part is named in 0 Task cells of The tasks of M2d, not one: `scaffold.md`: The acceptance tests of M2d (The setup questions in the batch and the comment)
  ```

## The checks (2026-10-12, 12:20Z)

`adr-lint`, `prd-lint`, `link-lint` (2,329 links), `setup-check`,
`run-discipline-tests` (102 passed) and `sh runs/T-bng8/parts.sh`: exit 0.
`git diff --cached --check`, after `git add` of the new files: exit 0.
