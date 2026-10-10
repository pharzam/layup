# T-mnmy — the test runs

All runs on the branch `T-mnmy`, base `40aeae9`, on 2026-10-10, with `mawk 1.3.4`.

## Red

At 17:42 UTC, the new and changed cases against the reader of `40aeae9`. Each
fails for the right reason: the reader prints `running` where `in progress` or
`stalled` is expected (exit 2, the output differs from `EXPECT`), and
`bad-missing-now` exits 2 because the old reader needs no `now`.

```
FAIL  task-state/bad-missing-now (wanted exit 1, got exit 2)
FAIL  task-state/good-comment-activity (wanted exit 0, got exit 2)
FAIL  task-state/good-crlf (wanted exit 0, got exit 2)
FAIL  task-state/good-in-progress (wanted exit 0, got exit 2)
FAIL  task-state/good-other-activity (wanted exit 0, got exit 2)
FAIL  task-state/good-window-boundary (wanted exit 0, got exit 2)
run-discipline-tests: 95 passed, 6 failed
```

## Green

At 17:43 UTC: `run-discipline-tests: 101 passed, 0 failed`.

## Each rule fails a case on its own

Each line is one change of `task-state.sh`, applied to a copy of `docs/` outside
the tree, and the cases that then fail.

| Mutation | Cases that fail |
| -------- | --------------- |
| the window strict (`>`, not `>=`) | `good-window-boundary` |
| comments ignored | `good-comment-activity` |
| the push wins when there is one (not the latest of the two) | `good-comment-activity` |
| a push of a branch that starts with the task ID counted | `good-other-activity` |
| any claimed task `in progress` | `good-in-progress`, `good-other-activity`, `good-window-boundary` |
| `now` not required | `bad-missing-now` |

## The fetcher on this repository

At 17:43 UTC, `sh docs/tasks/task-state.sh` in this worktree: five progress
lines (the pushes of the last `day`; the comments since `now` less 86400
seconds, `2026-10-09T17:43:08Z`), 4.6 s, and every row of `M2b` `done`, so no
row of the plan is claimed.

To see the two calls' inputs read, a second run, 15 seconds later, kept its
fetched directory (`now` 1791654203; the first run's `now` was 1791654188, whose
`since` is quoted above), and three claimed tasks were added under `## Now` of
its backlog copy (round 2, note 2):

| Task | Latest push | Latest comment on its issue | State |
| ---- | ----------- | --------------------------- | ----- |
| `T-mnmy` (#202) | 1791653521 (the claim, 17:32) | 1791654112 (#202) | `in progress` |
| `T-dm8m` (#185) | 1791617272 (the claim, 07:27) | 1791648432 (#185) | `in progress` |
| `T-eep8` (#175) | — | — | `done` (the completed log wins) |

No task of the forge was `stalled` at the time; the cases hold that state.

## Round 1 notes

At 17:55 UTC. Note 1: `good-window-boundary` gives `T-bbb2` an older comment
beside its push at the boundary; the mutation "the comment wins when there is
one" now fails `good-window-boundary`. Note 2: case `bad-now-not-a-number` was
red (exit 2: the reader read `yesterday` as 0 and printed with exit 0), and is
green after the check of `now`. Note 5: `good-review-head` holds a push of
`T-bbb4` in the window, so its README row is true. After the fixes:
`run-discipline-tests: 102 passed, 0 failed`.
