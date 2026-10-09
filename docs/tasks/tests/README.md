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
inputs have no extension (`backlog`, `completed`, `plan`, `forge/branches`,
`forge/pr/<n>/head` and `body`), so the Markdown checks do not read them as
documents. The base gives a done task, a ready one, two blocked ones and a
backlog task with no row in the plan, over two task tables of 14 and 13 columns.

| Case | Expected | Exercises |
| ---- | -------- | --------- |
| `good-done-ready-blocked` | exit 0 | the base: the rows of the last table only, then the backlog task with no row; `done` by the ID field of the completed log, not by a mention (`T-bbb2` is named in passing); an `After` of `—`, of one row and of two |
| `good-running` | exit 0 | a branch named the task ID; a `running` task with an unfinished predecessor shows it as `after T-bbb2` |
| `good-branch-prefix` | exit 0 | `T-bbb2-backlog` and `x-T-ccc1` are not claims |
| `good-review-head` | exit 0 | a pull request whose head is the task ID: `in review` before `running` and `blocked`; the smaller of two numbers |
| `good-review-closes` | exit 0 | `Closes #5` in the body of a pull request from another branch |
| `good-review-forms` | exit 0 | `FIXED #4` and `resolves: #9`: the forge's other forms, in any case, with a colon |
| `good-review-not-closing` | exit 0 | `Refs #4`, `Part of #5`, `Closes #44`, `prefixes #6` and `Closes #9x` close nothing |
| `good-done-first` | exit 0 | a done task with a branch and a pull request stays `done` |
| `bad-missing-completed` | exit 1 | no completed log |
| `bad-backlog-no-issue` | exit 1 | a backlog line under Now with no issue |
| `bad-after-unknown-row` | exit 1 | an `After` row that no task table holds |
| `bad-no-task-table` | exit 1 | a plan with no table whose header holds `#` and `Task ID` |
| `bad-row-cells` | exit 1 | a row with one cell more than its header |

Each `bad-*` case is otherwise valid, so it fails for its own single reason.
