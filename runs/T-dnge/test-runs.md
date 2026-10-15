# T-dnge — the test runs

Task `T-dnge` (#322), base `b7f3f7d`. A docs-only change, so its tests are two commands that are red on the base and green after, with the check of the coverage map (condition 1 of the
plan review, comment 6181732821). UTC, 2026-10-15, from the clock after each step.

## Step 3: red first

The plan review (2 min 23 s) and the author's answer (6181735981, 00:44:48Z) came first.

| Command | At | Result |
| ------- | -- | ------ |
| `test "$(awk '/^## The acceptance tests of M2d/,0' docs/spec/scaffold.md \| grep -c '^\| The review of the release \|')" = 1` | 00:44:59Z | count 0, **exit 1**: red, for the right reason (the table has no such row). |
| `sh runs/T-bng8/parts.sh` on the base | 00:45:00Z | exit 1, two lines: a row of "Input states" has no part, for the two rows that `T-gkvi` reworded. A defect revealed off this task's path: #328 (the correction to the answer, comment 6181750245). |

## Step 3: green, and the second red

The row of the table and the cell of row 59a (`scaffold.md`: The acceptance tests of M2d (The review of the release)) are added in one step, then the line of the map after the second red.

| Command | At | Result |
| ------- | -- | ------ |
| the count command above | 00:45:58Z | count 1, **exit 0**: green. |
| `sh runs/T-bng8/parts.sh`, row and cell added, no line of the map yet | 00:46:00Z | exit 1, three lines: the two of the base and **a row of "The acceptance tests of M2d" has no part: The review of the release**: the second red. |
| the same, with the line `docs/spec/scaffold.md`, `## The acceptance tests of M2d`, `The review of the release`, `59a` in `runs/T-bng8/parts.tsv` | 00:46:06Z | exit 1, the two lines of the base and no other: the new row is a part of row 59a, and its name is in the Task cell of row 59a (check 6). |

## The gate

On the head of the change, 00:46:13Z to 00:47:03Z, each exit 0: `go build ./...`, `go vet ./...`, `adr-lint`, `prd-lint`, `link-lint` (2,474 links resolved: the anchor of the plan), `run-discipline-tests` (138 passed),
`nested-checkout-check`, `setup-check` (check `adapted` reads the specification) and `git diff --check`. The Go levels are the pre-commit hook's unit run and CI's, as no Go file changed.
