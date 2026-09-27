# Git hooks

Local enforcement of the [quality gate](../docs/engineering-discipline.md). These
hooks catch a violation **before** it is committed, so it never reaches CI or a
reviewer. They are the fast-feedback twin of [`docs/ci/`](../docs/ci/), which is
the authority — both run the same rules.

Hooks live here (version-controlled), not in `.git/hooks` (local, untracked), so
the whole team shares one set.

## Install (one command)

```bash
sh .githooks/install.sh
```

This pins `core.hooksPath` to the relative `.githooks` for you — the script runs
`git config core.hooksPath .githooks` behind a guard and prints a confirmation.
Until you run it, the hooks are inert. Run it once per clone. To stop using them,
`git config --unset core.hooksPath`.

**Keep the path relative.** `.git/config` is **shared by every worktree**, so an
absolute value binds them all to one checkout's hooks — and a check added on a
branch then does not run on that branch's own commits. A relative value is
resolved per working tree, which is what you want. The `pre-commit` hook refuses
to run when the resolved path lies outside the tree being committed to; see
[`tests/provenance-check.sh`](tests/provenance-check.sh) for the proof, and
[`guardrails.md`](../docs/guardrails.md) for the trap it closes.

## What each hook does

| Hook | Runs | Set for LAYUP |
|------|------|---------------|
| [`pre-commit`](pre-commit) | A provenance check on where the hooks came from, then the ADR, PRD and link linters, the discipline self-tests, and the setup checks `adapted` and `markers` (`sh docs/setup/setup-check.sh --only adapted,markers .`), then, when `go.mod` exists, `gofmt -l` over the tracked Go files, `go vet ./...`, and the unit level `go test ./...`. Integration, end-to-end and the security scans run in CI only. | Set for Go (`T-t8qp`). |
| [`commit-msg`](commit-msg) | Conventional-Commits check on the subject line. | As shipped. |
| [`pre-push`](pre-push) | Refuses a direct push to `main` — use a branch and a PR instead. | LAYUP's default branch is `main`. |

## Protecting `main`

[`pre-push`](pre-push) blocks a direct push to `main` so changes go through a
branch and a pull request. It is a **local, fast-feedback guardrail only** —
advisory, bypassable with `git push --no-verify`, and absent on a fresh clone
until `core.hooksPath` is set. The lock that cannot be bypassed is the
**branch protection** of `main` on GitHub (a pull request before a merge, and the
required checks), enforced on the server; see
[Make the checks required](../docs/ci/README.md#make-the-checks-required). This
hook is its local twin, not a substitute.

## Changing a hook

1. [`pre-commit`](pre-commit) is set for Go (`T-t8qp`): lint and the unit level.
   A new step is cheap and fast; the full test suite runs in CI.
2. [`commit-msg`](commit-msg) takes the types that
   [§"Commit messages"](../docs/engineering-discipline.md#commit-messages) lists;
   a new type is agreed there first.

LAYUP uses these plain shell hooks, which need no dependency, and not the
[`pre-commit`](https://pre-commit.com) framework.
