# task-state self-tests

These are fixtures for [`../task-state.sh`](../task-state.sh), the command that
prints the state of each task (#175). It is a command, not a discipline linter,
so [`run.sh`](run.sh) gives it the exit contract of the
[discipline-test runner](../../tests/run-discipline-tests.sh): a `good*` case
exits 0 when the command exits 0 and its output equals `EXPECT` byte for byte; a
`bad*` case exits 1 when the command exits 1 and its stderr equals `EXPECT`; any
other outcome exits 2, so a bad case never passes for the wrong reason.

A case holds only the files it changes. The input is [`base/`](base/) with the
case laid over it, and a file `ABSENT` names the inputs a case removes. The
inputs have no extension (`backlog`, `completed`, `plan`, `now`,
`forge/branches`, `forge/pushes`, `forge/comments`, `forge/pr/<n>/head` and
`body`), so the Markdown checks do not read them as documents. The base gives a
done task, a ready one, two blocked ones and a backlog task with no row in the
plan, over two task tables of 14 and 13 columns, at a fixed `now` with no push
and no comment.

| Case | Expected | Exercises |
| ---- | -------- | --------- |
| `good-done-ready-blocked` | exit 0 | the base: the rows of the last table only, then the backlog task with no row; `done` by the ID field of the completed log, not by a mention (`T-bbb2` is named in passing); an `After` of `—`, of one row and of two |
| `good-in-progress` | exit 0 | a branch named the task ID with a push in the window is `in progress`; a branch with no activity is `stalled`, and shows its unfinished predecessor as `after T-bbb2` |
| `good-comment-activity` | exit 0 | a comment on the task's issue in the window, beside an old push of its branch: the latest counts, `in progress` |
| `good-window-boundary` | exit 0 | activity at `now - 86400` is `in progress`, at `now - 86401` `stalled`; `T-bbb2` also has an older comment, so the later push decides (round 1, note 1) |
| `good-other-activity` | exit 0 | a push of `T-bbb2-notes` and a comment on another issue are not the task's activity: `stalled` |
| `good-branch-prefix` | exit 0 | `T-bbb2-backlog` and `x-T-ccc1` are not claims |
| `good-review-head` | exit 0 | a pull request whose head is the task ID, with a push in the window: `in review` before `in progress` and `blocked`; the smaller of two numbers |
| `good-review-closes` | exit 0 | `Closes #5` in the body of a pull request from another branch |
| `good-review-forms` | exit 0 | `FIXED #4` and `resolves: #9`: the forge's other forms, in any case, with a colon |
| `good-review-not-closing` | exit 0 | `Refs #4`, `Part of #5`, `Closes #44`, `prefixes #6` and `Closes #9x` close nothing |
| `good-done-first` | exit 0 | a done task with a branch and a pull request stays `done` |
| `good-crlf` | exit 0 | every input and `EXPECT` with carriage returns, as a checkout with `core.autocrlf=true` gives them: `## Now`, a branch, a push, `now` and a closing keyword still match. Endings pinned by [`.gitattributes`](../../../.gitattributes) |
| `good-columns-order` | exit 0 | the last task table of 14 columns, with `After` first and `Task ID` last: the columns are found by the header, not the position; an `After` of four rows |
| `good-earlier-table` | exit 0 | a backlog task whose row is in an earlier task table keeps that row's `After`: `T-aaa2` is `blocked` on `T-aaa1`, listed after the rows of the last table (#182) |
| `bad-missing-now` | exit 1 | no `now`, the time of the run |
| `bad-now-not-a-number` | exit 1 | a `now` that is not a number of seconds (round 1, note 2) |
| `bad-missing-completed` | exit 1 | no completed log |
| `bad-backlog-no-issue` | exit 1 | a backlog line under Now with no issue |
| `bad-after-unknown-row` | exit 1 | an `After` row that no task table holds |
| `bad-no-task-table` | exit 1 | a plan with no table whose header holds `#` and `Task ID` |
| `bad-row-cells` | exit 1 | a row with one cell more than its header |

Each `bad-*` case is otherwise valid, so it fails for its own single reason.

## The cases of task-labels.sh (#185)

`run.sh` runs [`../task-labels.sh`](../task-labels.sh) when a case holds `state`
(the reader, `--read`), `APPLIER` (the applier, with the arguments the file
holds, against the stand-in [`gh-stub.sh`](gh-stub.sh) and a stand-in
`task-state.stub`; the calls the stand-in records are appended to the output) or
`WORKFLOW` (the check of a workflow: the case's `workflow`, or the repository's
own `task-state.yml` when the file says `real`).

| Case | Expected | Exercises |
| ---- | -------- | --------- |
| `good-labels-each` | exit 0 | each label added: a closed issue `completed` (`merged`) and `not_planned` (`rejected`), an open issue the command says is `done` (`merged`), `in review`, `in progress`, `stalled`, `ready`; an issue only the completed log names (`#7`) |
| `good-labels-right` | exit 0 | nothing for an issue already right; a hand-set label replaced (`blocked`); a second `state:` label removed; an open task issue the command does not list loses its `state:` label; an issue closed as `duplicate` while the command says `in progress` is `rejected`; an issue closed with an empty reason is `rejected` (`#7`, named by a `plan` of the case); labels outside `state:` and an issue outside the set untouched |
| `good-labels-crlf` | exit 0 | `state`, `issues` and `EXPECT` with carriage returns. Endings pinned by [`.gitattributes`](../../../.gitattributes) |
| `good-labels-applier-dry-run` | exit 0 | `--dry-run`: the four reads, the changes printed, no write |
| `good-labels-applier` | exit 0 | the writes: the seven labels created, then one `issue edit` a task — an addition, a removal only, an addition with a removal; the calls the stand-in records prove each |
| `good-labels-workflow` | exit 0 | the repository's `task-state.yml` passes its check |
| `bad-labels-empty-set` | exit 1 | the three files name no task issue; the calls show the reads and no write |
| `bad-labels-disjoint` | exit 1 | an issue read that holds none of the task issues; the calls show the reads and no write |
| `bad-labels-missing-issues` | exit 1 | no `issues` |
| `bad-labels-unknown-state` | exit 1 | a state the labels do not know |
| `bad-labels-applier-state-fails` | exit 1 | `task-state.sh` fails: no call to the forge at all |
| `bad-labels-applier-issues-fail` | exit 1 | the issue read fails after the three files: no write |
| `bad-labels-workflow-dispatch` | exit 1 | a `workflow_dispatch` trigger |
| `bad-labels-workflow-push` | exit 1 | a `push` trigger with no `branches:` bound |
| `bad-labels-workflow-branches` | exit 1 | a `push` to a branch that is not `main` |
| `bad-labels-workflow-tags` | exit 1 | a `push` bound that is not `branches` (`tags`) |
| `bad-labels-workflow-tags-ignore` | exit 1 | a `push` bound that is not `branches` (`tags-ignore`, a hyphenated key) |
| `bad-labels-workflow-tags-quoted` | exit 1 | a `push` bound that is not `branches` (`'tags'`, a quoted key) |
| `bad-labels-workflow-types` | exit 1 | a named type missing from `pull_request_target` |
| `bad-labels-workflow-issues-types` | exit 1 | a `types:` of `issues` with a foreign and a missing value |
| `bad-labels-workflow-write` | exit 1 | `pull-requests: write` |
| `bad-labels-workflow-ref` | exit 1 | a checkout with no `ref:` of the default branch |
| `bad-labels-workflow-event` | exit 1 | an event field in a one-line `run:` step |
| `bad-labels-workflow-event-block` | exit 1 | an event field in a `run: |` block |
| `bad-labels-workflow-restore` | exit 1 | no step named `Restore` that names `docs/tasks` |
