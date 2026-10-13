# T-p481 — test runs

Each `state:` label in the color of its state (#309): the change of `apply()` in
`docs/tasks/task-labels.sh`, and the `label create` lines the two applier cases
of `docs/tasks/tests/` expect. All times UTC, 2026-10-13.

## Red (10:25Z)

The two `EXPECT` files at the new colors, `task-labels.sh` and `run.sh` as at
the base `6734a23` — run on a copy of the tree, because the `pre-commit` hook
keeps a red commit off the branch (the check of guardrails.md §2, "The hook
refuses a red commit"). Both cases exit 2 — an output/EXPECT difference — each
on the `ededed` byte alone, the right reason:

```
6,13c6,13
< label create state: ready --force --color 1d76db --description Set by .github/workflows/task-state.yml; never by hand
< label create state: in progress --force --color fbca04 --description Set by .github/workflows/task-state.yml; never by hand
< label create state: stalled --force --color d93f0b --description Set by .github/workflows/task-state.yml; never by hand
< label create state: in review --force --color 5319e7 --description Set by .github/workflows/task-state.yml; never by hand
< label create state: blocked --force --color b60205 --description Set by .github/workflows/task-state.yml; never by hand
< label create state: merged --force --color 0e8a16 --description Set by .github/workflows/task-state.yml; never by hand
< label create state: rejected --force --color 57606a --description Set by .github/workflows/task-state.yml; never by hand
< label create state: on hold --force --color c5def5 --description Set by .github/workflows/task-state.yml; never by hand
---
> label create state: ready --force --color ededed --description Set by .github/workflows/task-state.yml; never by hand
> ... the same for the seven other labels
```

`good-labels-applier-prs` fails the same way, on the same eight lines of its own
`EXPECT`.

## Mutation (10:26Z)

Two single-value mutants, each on a copy of the tree, to prove a wrong color
cannot pass: `state: ready=ffffff` fails `good-labels-applier` on the `ready`
line, and `state: on hold=eeeeee` fails `good-labels-applier-prs` on the
`on hold` line. Both exit 2, each diff on the one mutated byte.

## Green (10:27Z)

`sh docs/tests/run-discipline-tests.sh` on the branch: **138 passed, 0 failed**.
The `pre-commit` hook of the feat commit ran the same battery plus `go vet ./...`
and the unit level `go test ./...`, all green.

## The forge read (post-merge)

The plan review's evidence boundary: the stub-recorded `label create` calls
above are the pre-merge evidence; `gh label list` after the first run of
`task-state.yml` on `main` is the confirmation recorded on the issue, not a
pre-merge gate.
