# Forge templates

Issue and pull-request templates that follow the [issue-first
workflow](../issue-workflow.md): single-goal scope (R11), the duplicate check
(R2), solution selection (R3), and the action/why/tradeoff comment (R7).

LAYUP's forge is GitHub. The two templates are here, in
[`github/`](github/), and are **not active**: `.github/ISSUE_TEMPLATE/` and
`.github/PULL_REQUEST_TEMPLATE.md` do not exist, because an issue or pull-request
template changes the live forge interface when it lands. A LAYUP task issue and
its pull request follow the shape of these files by hand:

- [`github/ISSUE_TEMPLATE/task.md`](github/ISSUE_TEMPLATE/task.md), the task issue;
  its metadata header (`name`, `about`, `labels`) is the part the forge reads.
- [`github/PULL_REQUEST_TEMPLATE.md`](github/PULL_REQUEST_TEMPLATE.md), the pull
  request.

Each `‹…›` marker in them is a field of a new issue or pull request.

> These differ from the per-record templates [`adr/template.md`](../adr/template.md),
> [`facts/template.md`](../facts/template.md) and [`prd/template.md`](../prd/template.md),
> which are copied *inside* their own directories, not into a forge path.
