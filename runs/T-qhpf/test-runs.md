# T-qhpf — the test runs

All runs on the branch `T-qhpf`, base `91e41b5`, on 2026-10-09, with `mawk 1.3.4`.

## Red (note 1)

At 19:54 UTC, the case `good-earlier-table` against the reader of `91e41b5`:
`T-aaa2` (row 2, in the first table, `After` 1, with `T-aaa1` not done) prints
`ready`, not `blocked`. `run.sh` exits 2, the right reason:

```
FAIL  task-state/good-earlier-table (wanted exit 0, got exit 2)
run-discipline-tests: 96 passed, 1 failed
```

The mutation that must fail the case is the reader as it was: the `After` of a
`## Now` task with no row in the last table taken as empty. That is the red above.

## Green

At 19:54 UTC, after the lookup of the row's `After` in any task table:
`run-discipline-tests: 97 passed, 0 failed`.

## Note 3: the files of main, not of the checkout

At 19:55 UTC, with `main` at `c15e2c8` on the forge (`gh api repos/pharzam/layup/commits/main`, the same
before and after the runs). A clone at `b1ff841` (the merge of `T-eep8`, seven merges
behind: #181, #183, #184, #187, #188, #190, #191; round 1, note 3) ran both scripts:

- the new script there printed the same bytes as the new script in this worktree
  (`cmp`: equal), the state of `main`;
- the old script there (the one of `b1ff841`, which copies the checkout's files)
  differed in eight lines: it read the old completed log, so it printed seven
  merged tasks as `running` (their branches are still on the forge) and
  `T-z027` blocked on `T-6sbe`, which is done.

```
4c4
< T-cht1	#156	running	—
---
> T-cht1	#156	done	—
7,12c7,12
< T-vxdg	#159	running	—
< T-6sbe	#160	running	after T-vxdg
< T-5pxd	#161	running	after T-vxdg
< T-bpxg	#162	running	after T-5pxd
< T-4c3q	#163	running	—
< T-d8t9	#164	running	after T-cht1,T-6sbe,T-5pxd
---
> T-vxdg	#159	done	—
> T-6sbe	#160	done	—
> T-5pxd	#161	done	—
> T-bpxg	#162	done	—
> T-4c3q	#163	done	—
> T-d8t9	#164	done	—
14c14
< T-z027	#166	blocked	T-6sbe,T-fsjp
---
> T-z027	#166	blocked	T-fsjp
```
