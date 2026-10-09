# T-eep8 — the test runs

All runs on the branch `T-eep8`, base `ad8774d`, on 2026-10-09, with `mawk 1.3.4`.

## Red (D1 against a stub that prints nothing)

`sh docs/tests/run-discipline-tests.sh`, with `docs/tasks/task-state.sh` a stub
that exits 0 and prints nothing. Every case of the suite fails for the right
reason: `run.sh` exits 2 because the output (or the stderr) differs from `EXPECT`.

```
FAIL  task-state/bad-after-unknown-row (wanted exit 1, got exit 2)
FAIL  task-state/bad-backlog-no-issue (wanted exit 1, got exit 2)
FAIL  task-state/bad-missing-completed (wanted exit 1, got exit 2)
FAIL  task-state/bad-no-task-table (wanted exit 1, got exit 2)
FAIL  task-state/bad-row-cells (wanted exit 1, got exit 2)
FAIL  task-state/good-branch-prefix (wanted exit 0, got exit 2)
FAIL  task-state/good-done-first (wanted exit 0, got exit 2)
FAIL  task-state/good-done-ready-blocked (wanted exit 0, got exit 2)
FAIL  task-state/good-review-closes (wanted exit 0, got exit 2)
FAIL  task-state/good-review-forms (wanted exit 0, got exit 2)
FAIL  task-state/good-review-head (wanted exit 0, got exit 2)
FAIL  task-state/good-review-not-closing (wanted exit 0, got exit 2)
FAIL  task-state/good-running (wanted exit 0, got exit 2)
run-discipline-tests: 81 passed, 13 failed
```

## Green (D2)

```
run-discipline-tests: 94 passed, 0 failed
```

## Each rule fails a case on its own

Each line is one mutation of `task-state.sh`, applied with `sed` to a copy of
`docs/` outside the tree, and the cases that then fail.

| Mutation | Cases that fail |
| -------- | --------------- |
| `done` read from any mention of an ID, not the ID field | 8 good cases (the base names `T-bbb2` in passing) |
| a keyword not at the start of a word (`prefixes #6`) | `good-review-not-closing` |
| a number followed by a letter (`Closes #9x`) | `good-review-not-closing` |
| R1's three forms only | `good-review-forms` |
| the case of the keyword kept | `good-review-closes`, `good-review-forms` |
| a branch that starts with the task ID counts | `good-branch-prefix` |
| `running` only with no predecessor | `good-running` |
| `done` not first | 8 good cases |
| the largest pull request, not the smallest | `good-review-head` |
| every backlog section, not only `## Now` | 8 good cases (`T-ddd1` of Next appears) |
| the first task table, not the last | 8 good cases and `bad-after-unknown-row` |
| a row of `After` not checked | `bad-after-unknown-row` |
| the count of cells not checked | `bad-row-cells` |
| the detail of a `running` task without its predecessors | `good-running` |
| the head branch of a pull request not matched | `good-review-head` |

## The fetcher on this repository

`sh docs/tasks/task-state.sh` in the worktree, at 12:58 UTC; 1.3 s. `T-eep8`
is not listed: its backlog line is in PR #177, not yet merged.

```
task-state: reading the branches of the forge
task-state: reading the open pull requests
T-3py1	#153	done	—
T-y10b	#154	done	—
T-ysph	#155	in review	#176
T-cht1	#156	blocked	T-ysph
T-z5dj	#157	in review	#174
T-m1dx	#158	ready	—
T-vxdg	#159	blocked	T-ysph,T-z5dj
T-6sbe	#160	blocked	T-vxdg
T-5pxd	#161	blocked	T-ysph,T-vxdg
T-bpxg	#162	blocked	T-5pxd
T-4c3q	#163	ready	—
T-d8t9	#164	blocked	T-cht1,T-6sbe,T-5pxd
T-fsjp	#165	blocked	T-bpxg,T-4c3q,T-d8t9
T-z027	#166	blocked	T-z5dj,T-m1dx,T-6sbe,T-fsjp
T-e3sy	#167	blocked	T-z027
T-nxe4	#168	blocked	T-fsjp
T-4tjy	#169	blocked	T-e3sy,T-nxe4
T-x7cs	#170	blocked	T-4tjy
```
