# T-ye17 — the test runs

All runs on the branch `T-ye17`, base `95dccce`, on 2026-10-10, with `mawk 1.3.4`.
CI on the pull request runs `main`'s copy of the check and of `docs/setup/tests/`
(`.github/workflows/ci.yml`, job `setup-check`), so these local runs are the
evidence of the red and the green before the merge (plan review, condition 3).

## Red

At 16:16 UTC, `sh docs/setup/tests/run.sh` with the new cases against the check
of `95dccce`: the good case fails for the right reason, each job of a workflow
that no `pull_request` event runs read as a job that is not a required context.
`bad-pull-request-forms` passes already: the check of `95dccce` lists every job.

```
FAIL  protection/good-no-pull-request: exit 1, want 0
FAIL  protection/good-no-pull-request: no output line: setup-check: protection OK
      | setup-check: protection FAIL contexts: flow is a job but not a required context
      | setup-check: protection FAIL contexts: labels (not a gate) is a job but not a required context
      | setup-check: protection FAIL contexts: scalar is a job but not a required context
      | setup-check: protection FAIL contexts: seq is a job but not a required context
setup-check tests: 47 passed, 1 failed
```

## Green

At 16:17 UTC, after the reading of `on:`: `setup-check tests: 48 passed, 0 failed`;
`sh docs/setup/setup-check.sh` on this tree: `protection OK` (the four workflows
of the tree all run on `pull_request`, so no context changes).

## Each rule fails a case on its own

Each line is one change of the check, applied to a copy of `docs/` outside the
tree, and the cases that then fail. The copy is not a Git repository, so
`facts/good-autocrlf`, which clones the repository that holds the runner, fails
in every run of the copy, the unchanged one included; it is left out below.

| Mutation | Cases that fail |
| -------- | --------------- |
| `pull_request` read as a prefix (`pull_request_target` counts) | `protection/good-no-pull-request` |
| every workflow skipped | `frame/good-passthrough`, `protection/bad-contexts`, `protection/bad-pull-request-forms`, `protection/good-no-pull-request` |
| a workflow with no `on:` read as not run | `protection/bad-pull-request-forms` |
| a flow map or an alias read as not run | `protection/bad-pull-request-forms` |
| the block sequence not read | `protection/bad-pull-request-forms` |
| quotes kept on a flow-list item | `protection/bad-pull-request-forms` |
| a quoted `on` not read | `protection/good-no-pull-request` |

## Round 1 fixes (findings 1 and 2, note 4)

Four workflows added to `bad-pull-request-forms`: `seq-zero` (a block sequence
at the indent of `on:`, with `- pull_request`), `split-flow` (a flow list split
over two lines, with `pull_request_target`: read as run by `pull_request`, the
safe side), `anchor` (`on: &triggers […]`) and `merge-key` (`<<: *triggers`).

Red at 16:42 UTC, against the check of `a02031c`: the two findings, and only
they, are missing.

```
FAIL  protection/bad-pull-request-forms: no output line: setup-check: protection FAIL contexts: seq-zero is a job but not a required context
FAIL  protection/bad-pull-request-forms: no output line: setup-check: protection FAIL contexts: split-flow is a job but not a required context
setup-check tests: 47 passed, 1 failed
```

Green at 16:42 UTC: `setup-check tests: 48 passed, 0 failed`; `protection OK` on
this tree.

| Mutation (on a copy, as above) | Cases that fail |
| ------------------------------ | --------------- |
| a sequence at the indent of `on:` ends the block | `protection/bad-pull-request-forms` |
| a split flow list read as a one-line list | `protection/bad-pull-request-forms` |
| an anchor read | `protection/bad-pull-request-forms` |
| a merge key read | `protection/bad-pull-request-forms` |
| an alias read | `protection/bad-pull-request-forms` |
