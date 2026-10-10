# Continuous integration

CI runs the [quality gate](../engineering-discipline.md) on every change, so the
forge enforces the gate and not memory. It is the **authority**: its checks are
the ones that the branch protection makes required before a merge. The
[`.githooks/`](../../.githooks/) run the same rules locally for fast feedback.

LAYUP's CI is on GitHub Actions, in [`.github/workflows/`](../../.github/workflows/):

| Workflow | Jobs |
|----------|------|
| [`ci.yml`](../../.github/workflows/ci.yml) | `adr-lint`, `prd-lint`, `discipline-tests`, `link-lint`, `nested-checkout-check`, `setup-check`, and the Go jobs `lint`, `tests`, `security` ([ADR-0010](../adr/0010-use-go-as-the-technology-stack.md), since `T-t8qp`) |
| [`pr-title.yml`](../../.github/workflows/pr-title.yml) | `conventional-title` |
| [`pr-link.yml`](../../.github/workflows/pr-link.yml) | `pr-link` |
| [`review-record.yml`](../../.github/workflows/review-record.yml) | `review-record` |

This directory holds the two scripts that read forge artifacts, and their
fixtures under [`tests/`](tests/).

## What the jobs run

| Job | What it checks |
|-----|----------------|
| `adr-lint` | `docs/adr/` discipline, via [`adr-lint.sh`](../adr/adr-lint.sh). |
| `prd-lint` | `docs/prd/` discipline, via [`prd-lint.sh`](../prd/prd-lint.sh). |
| `discipline-tests` | Each discipline linter against its good and bad fixtures, and the fixtures of the command [`task-state.sh`](../tasks/task-state.sh), which is not a linter, via [`run-discipline-tests.sh`](../tests/run-discipline-tests.sh). |
| `link-lint` | Every in-tree Markdown link and heading anchor, via [`link-lint.sh`](../links/link-lint.sh). |
| `nested-checkout-check` | The linters skip a nested checkout, via [`nested-checkout-check.sh`](../tests/nested-checkout-check.sh). |
| `setup-check` | The setup of this repository and its evidence, via [`setup-check.sh`](../setup/setup-check.sh) and its fixtures. |
| `lint` | `gofmt` and `go vet` over the Go code. |
| `tests` | The test ladder, cheap to expensive: unit, integration, end-to-end (see [`test-levels.md`](../tests/test-levels.md)). |
| `security` | `govulncheck v1.8.0`, `go vet` and `gitleaks` over the full history (see [`security-checklist.md`](../tests/security-checklist.md)). |
| `conventional-title` | Conventional Commits on the pull-request title. |
| `pr-link` | The pull-request body links an issue (R1), via [`pr-link-lint.sh`](pr-link-lint.sh). |
| `review-record` | The linked issue carries a plan, a plan review with its budget and cycle cap, and a parseable review record per round whose chronology holds (see [What a round records](../engineering-discipline.md#what-a-round-records)), via [`review-record-lint.sh`](review-record-lint.sh). |

**Each discipline job restores its check scripts from the default branch before
it runs them**, with one exception: the job `nested-checkout-check` runs the
branch's own `nested-checkout-check.sh`, for the reason that a comment in
[`ci.yml`](../../.github/workflows/ci.yml) gives. The Go jobs `lint`, `tests` and
`security` run the pull request's own Go code and tests. Without the restore, a pull request that replaces a linter with `exit 0` passes
that linter's own required check, and a pull request that also replaces
`run-discipline-tests.sh` turns the whole required set green over a real defect.
Both were measured. Thus a check that a pull request adds runs in CI only after
it merges. The limit that remains is that the workflow file itself comes from the
pull request, so a branch that edits `ci.yml` removes the step; only review of
`.github/**` closes that, and [`docs/guardrails.md`](../guardrails.md) states all
three residual limits. A `CODEOWNERS` entry over `.github/**`, `.githooks/**` and
`*-lint.sh` needs two human operators; LAYUP has one, so it has no such entry.

A discipline test that lints files in the repository runs both in a CI job and in
the [`pre-commit`](../../.githooks/pre-commit) hook, as `adr-lint` and `prd-lint`
do. A check whose input is a forge artifact, not a repository file, has no local
hook to run in and runs in CI only. Two do: [`pr-link-lint.sh`](pr-link-lint.sh)
reads the pull-request body, which the event gives it;
[`review-record-lint.sh`](review-record-lint.sh) reads the linked issue's
comments, which the workflow fetches. In both, the forge call is in the
**workflow** and the parsing is in the **script**, so
[`run-discipline-tests.sh`](../tests/run-discipline-tests.sh) self-tests both
offline, with no network and no token.

## Make the checks required

A check that runs but does not block is a run result, not a merge control. Until
a check is required on the default branch, a red run and a green one merge alike.
This is the step that the branch-protection column of the
[enforcement table](../issue-workflow.md#what-is-enforced-where) names. LAYUP
keeps the setting as a file,
[`docs/setup/branch-protection.json`](../setup/branch-protection.json) (`T-afa5`):
twelve checks, one per job of a workflow that a `pull_request` event runs, each
pinned to `"app_id": 15368`, GitHub Actions. A bare `contexts` list would let any
app or token satisfy a name by posting a status under it. A context is the job's
`name:`, or its id when it has none, which is why `conventional-title` has no
parenthesis. Check `protection` in [`setup-check.sh`](../setup/setup-check.sh)
keeps the contexts equal to the names of those jobs. A workflow that no
`pull_request` event runs (by the key `pull_request`, so `pull_request_target`
is not it) judges no pull request: its jobs are not required, and the body must
not list one (#201). A workflow that judges a pull request uses `pull_request`.
The forms of `on:` the check reads, and its limit, are in its comment.

Apply it with:

```bash
gh api -X PUT repos/pharzam/layup/branches/main/protection --input docs/setup/branch-protection.json
```

**The body is the whole setting, so a partial body is destructive.** The `PUT`
replaces the entire protection object, and `required_status_checks`,
`enforce_admins`, `required_pull_request_reviews` and `restrictions` are required
parameters: a body that omits `required_conversation_resolution` turns it off, and
one that sends `required_pull_request_reviews` as `null` removes the pull-request
requirement while it adds the check requirement. `"restrictions": null` is
required, and must be null on a user-owned repository, because push restrictions
exist for organizations only. `required_signatures` is a separate endpoint that
the `PUT` does not touch. Read the setting back with
`gh api repos/pharzam/layup/branches/main/protection`; that output, not the green
run, is the evidence a close-out records.

**Limits.** To write or read the setting needs an administration-scoped token,
which `secrets.GITHUB_TOKEN` does not carry, so no text-only check proves it: the
verification is the command above, run by an operator, and each close-out records
its output. A renamed or removed **job** blocks every merge until the setting
follows it (the check name is the job's `name:` or its id, never the workflow's
name), which is the right direction of failure, loud and where the gate lives.
`enforce_admins: true` binds the operator too. `strict: true` enforces the rule
to bring a branch up to the latest `origin/main` before a merge (by a merge of
`origin/main`, for a branch on the forge:
[Integrating branches](../engineering-discipline.md#integrating-branches)), at a
price: with a plain merge, every merge to the default branch makes the up-to-date
status of every other open pull request stale, so a queue of pull requests
becomes update, re-run, merge. And a required check is only as trustworthy as the script it runs: the
restore step above closes that for the check scripts, but the Go jobs run the pull
request's own Go code and tests (the Invariant 3 gap, ADR-0011, O-9).

**What one wrong context costs.** A required context that no workflow reports
never arrives, so the check stays pending and no pull request merges. The body
sets `"enforce_admins": true`, so no pull request merges its way out. An
administrator edits the protection in the repository's settings, under Branches,
at no token cost (`enforce_admins` binds merges, not the setting itself). The
scripted route is a second `PUT` of the whole corrected body, which needs the
administration-scoped token that **Limits** names.
