# Setup record — T-n1hp, the manual Armature setup of LAYUP (baseline)

This is the record of one manual run of the Armature setup, on this repository,
on 2026-09-23. It is PSB Problem 3 measured once (the PSB is fact `F-0001`; see D-03):
the baseline for "Reproducible Discipline Setup" (PSB §6 In Scope) and its first
test case. Parent issue: [#1](https://github.com/pharzam/layup/issues/1). Each
child task adds its own rows in the same PR that lands its work.

## In plain terms

> One person and one agent set up the discipline template by hand. This file
> writes down, for each value that was filled, where the value came from. A value
> with no source is not filled; it is listed as a question for the idea owner.
> The file also writes down how long each step took, so a later automated setup
> can be compared with this manual one.

## How to read this record

- **Evidence** names the source of a value: a file and line, a command and its
  output, a fact ID (`F-0001#n`), or an Operator decision with its date. A value
  whose only source is a guess is a defect (PSB Invariant 4).
- **Status** is `active` when a gate reads the value today, `not active` when a
  gate must read it but none does yet (PSB Invariant 5), and `recorded` when the
  value is provenance that no gate reads by design.
- Times are UTC, from the Git commit, the GitHub API `created_at`, or the file
  system, as the Source column says.

## Timeline

| # | Event | Time (UTC) | Source |
|---|-------|-----------|--------|
| 1 | Operator decisions: stack Go; public repo `pharzam/layup`; full gate | before 10:36:52 | session answers, 2026-09-23 |
| 2 | Kit copied with `npx degit` into `~/projects/layup` | 10:21:59 | directory modification time |
| 3 | Root commit `d2516fd`, the unmodified kit | 10:36:52 | commit date |
| 4 | Repository `pharzam/layup` created (public), `main` pushed | 10:36:55 | GitHub API `created_at` |
| 5 | Issue #1 opened | 10:37:42 | GitHub API |
| 6 | Plan (R12) on #1 | 10:38:21 | GitHub API |
| 7 | Plan review on #1: `approve-with-conditions`, split into goal classes (R11) | 10:45:55 | GitHub API |
| 8 | Children #2–#11 opened, each with its plan | 10:49:20–10:49:50 | GitHub API |
| 9 | Plan reviews of the children | 11:01:54 (for #2) | GitHub API |
| 10 | Operator approval of the parent budget, 7,000 lines over 145 files | 11:05:53 | GitHub API |
| 11 | Child #12 opened (branch protection, split out of #9) | 11:05:51 | GitHub API |
| 12 | `T-r7zg` worktree created | 11:06:44 | shell `date` |
| 13 | `T-r7zg` merged (PR #13), after 2 review rounds | 11:34:13 | GitHub API `mergedAt` |
| 14 | `T-vbwc` worktree created | 11:34:46 | shell `date` |
| 15 | `T-vbwc` merged (PR #14), after 3 review rounds | 12:08:57 | GitHub API `mergedAt` |
| 16 | `T-fvwj` worktree created | 12:09:37 | shell `date` |
| 17 | `T-fvwj` merged (PR #16), after 2 review rounds | 12:29:05 | GitHub API `mergedAt` |
| 18 | `T-vpty` worktree created | 12:29:08 | shell `date` |
| 19 | `T-vpty` merged (PR #17), after 2 review rounds | 12:50:25 | GitHub API `mergedAt` |
| 20 | `T-xgz4` worktree created | 12:50:28 | shell `date` |
| 21 | `T-xgz4` merged (PR #18), after 3 review rounds | 13:22:02 | GitHub API `mergedAt` |
| 22 | `T-7ndb` worktree created | 13:22:05 | shell `date` |
| 23 | `T-7ndb` merged (PR #19), after 2 review rounds | 13:48:38 | GitHub API `mergedAt` |
| 24 | `T-nfh8` worktree created | 13:48:41 | shell `date` |
| 25 | Operator answered one batch of 8 setup questions (O-1 to O-8) | about 13:50 | comment on [#8](https://github.com/pharzam/layup/issues/8) |

## Values

| ID | Value | File | Evidence | Status | Task |
|----|-------|------|----------|--------|------|
| V-01 | `source=https://github.com/pharzam/armature` | [`armature.pin`](armature.pin) | PSB header "Target Discipline Baseline" (`F-0001`, see D-03) | recorded | `T-r7zg` |
| V-02 | `commit=a95965534b14b0bf14ad74da0c9a45b5f4aedf88` | [`armature.pin`](armature.pin) | `gh api repos/pharzam/armature/commits/main` printed this SHA just before the copy (timeline row 2); the degit command named it | active | `T-r7zg` |
| V-03 | `tree=8ffb250afd584da8b418bc220fb6d72e802924ce` | [`armature.pin`](armature.pin) | `git rev-parse d2516fd^{tree}`; equal to the Armature commit tree (plan review of #1, GitHub API) | active | `T-r7zg` |
| V-04 | `method=npx degit pharzam/armature#<commit>` | [`armature.pin`](armature.pin) | the command that was run (kit README, step 1, with the commit added) | recorded | `T-r7zg` |
| V-05 | `date=2026-09-23` | [`armature.pin`](armature.pin) | timeline row 2 | recorded | `T-r7zg` |
| V-06 | PSB SHA-256 `3e96862b2578e74f42bfd24929926d0c64496f1a9a51f78dc8c6ca9ecf394ab8` | [`facts.sha256`](facts.sha256) | `shasum -a 256` of the approved Revision 6 file in the local research repository (`4cdf9d0`, no remote) | active | `T-fvwj` |
| V-07 | Vision brief SHA-256 `b1354c71b3d686331a2cfd1f37a175f5a8ec44fb8dcd863e2cac0138d2d9f625` | [`facts.sha256`](facts.sha256) | same command, same repository | active | `T-fvwj` |
| V-08 | default branch `main`; owner `pharzam`; repo `layup` | `docs/ci/README.md` | `gh repo view pharzam/layup --json defaultBranchRef` | recorded | `T-nfh8` |
| V-09 | worktree directory `.worktree` (gitignored) | `AGENTS.md`, `engineering-discipline.md`, `.gitignore` | O-1, Operator decision on [#8](https://github.com/pharzam/layup/issues/8) | recorded | `T-nfh8` |
| V-10 | evidence store `runs/` | `engineering-discipline.md` | O-2, Operator decision on [#8](https://github.com/pharzam/layup/issues/8) | not active | `T-nfh8` |
| V-11 | reasoning tier: Claude Opus 5.5, Claude Fable 5.1; execution tier: Claude Sonnet 5, Claude Haiku 4.5 | `engineering-discipline.md` | O-3, Operator decision on [#8](https://github.com/pharzam/layup/issues/8) | recorded | `T-nfh8` |
| V-12 | panel-worthy: each new ADR | `engineering-discipline.md` | O-4, Operator decision on [#8](https://github.com/pharzam/layup/issues/8) | recorded | `T-nfh8` |
| V-13 | test runner `go test`; unit `go test ./...`; integration `go test -tags=integration ./...`; e2e `go test -tags=e2e ./...`; test directory `*_test.go` beside the code, root `tests/` for e2e fixtures | `docs/tests/*`, `.githooks/pre-commit`, `tests/README.md`, `docs/prd/README.md` | O-5, Operator decision on [#8](https://github.com/pharzam/layup/issues/8); `go help test`, `go help buildconstraint` (build tags) | not active | `T-nfh8` |
| V-14 | lint `test -z "$(gofmt -l .)" && go vet ./...` | `.githooks/*`, `engineering-discipline.md`, `docs/ci/README.md` | `gofmt -h` (`-l` lists files that differ), `go help vet`; stack from [ADR-0010](../adr/0010-use-go-as-the-technology-stack.md) | not active | `T-nfh8` |
| V-15 | test timeout `-timeout 10m` | `docs/tests/*`, `guardrails.md` | `go help testflag`: "The default is 10 minutes (10m)" | not active | `T-nfh8` |
| V-16 | security: govulncheck `v1.8.0`, `go vet`, gitleaks (`gitleaks git --redact`) | `docs/tests/*`, `.githooks/pre-commit`, `docs/ci/README.md` | O-6, Operator decision on [#8](https://github.com/pharzam/layup/issues/8); `go list -m golang.org/x/vuln@latest` gave `v1.8.0`, which needs Go 1.26; the gitleaks command line is not yet run | not active | `T-nfh8` |
| V-17 | task-ID scheme `T-` + 4 of `0-9 a-z` without `i l o u` | `backlog.md`, `issue-workflow.md`, `engineering-discipline.md` | O-8, Operator decision on [#8](https://github.com/pharzam/layup/issues/8) (confirms the kit example in use since `T-n1hp`) | recorded | `T-nfh8` |
| V-18 | record of record: Git; plain-terms units: tokens, seconds, "N of M"; failure modes; harness report source | `guardrails.md`, `engineering-discipline.md` | `F-0001#1`, `F-0001 §7`, `guardrails.md` §1.1–§2, this session's task reports | recorded | `T-nfh8` |

The decision behind V-01 to V-05 is
[ADR-0009](../adr/0009-pin-armature-at-a-recorded-commit.md).

## Findings about the kit's setup procedure

Each finding is a place where the manual setup needed a decision that the kit
does not state. They are inputs for the automated setup; none of them is a change
to an Armature rule.

| ID | Finding | Where in the kit | Effect on this run |
|----|---------|------------------|--------------------|
| K-01 | The Operator gave the count "17 domain terms"; PSB Revision 6 §8 holds 25 rows (one row holds two names). | not a kit gap: a gap between the instruction and the fact | All 25 rows are merged into [`glossary.md`](../glossary.md) §1 (`T-xgz4`). |
| K-02 | degit copies the kit's own active workflows under `.github/workflows/`; kit step 4 does not name them. | `engineering-discipline.md` How to adapt, step 4 | Handled by `T-q344`. |
| K-03 | `docs/tasks/backlog.md` holds the kit's own tasks and a pivot note that links the kit's issues. Kit step 1 says to fill `backlog.md` with your own tasks, but step 4, the removal list, names only `completed.md` entries and the `T-*.md` files, so the kit backlog lines are easy to miss. | How to adapt, step 4 | `T-vbwc` removed 31 kit files (3,133 lines: `docs/decisions/` 10, `docs/audit/` 2, `docs/tasks/T-*.md` 19), the kit entries of `completed.md`, and the 6 kit lines and the pivot note of `backlog.md`. |
| K-04 | The CI template `docs/ci/github-actions-ci.yml` has no "restore the checks from the default branch" step, which `guardrails.md` §2 requires and the kit's own `ci.yml` has. | `docs/ci/github-actions-ci.yml` | `T-q344` keeps the kit's own workflows. |
| K-05 | The first push to `main` must come before the `pre-push` hook is installed, because the hook refuses a direct push to `main`. | README step 1 and step 3 order | The root commit was pushed before `sh .githooks/install.sh`. |
| K-06 | The adoption steps do not say to rewrite the "what this is" text of `README.md` and `AGENTS.md`, which name the kit. | README "Using it as a template" | Handled by `T-6rg3`. |
| K-07 | The pin tree check needs the full history; `actions/checkout` is shallow by default. | not in the kit (LAYUP check) | `T-q344` sets `fetch-depth: 0`. |

## Deviations

| ID | Deviation | Reason | Ends |
|----|-----------|--------|------|
| D-01 | Worktrees were under `../layup-worktrees/`, outside the repository. | The worktree-directory value was not yet set. | Ended by `T-nfh8`: O-1 set `.worktree`, and the `T-nfh8` worktree moved there. |
| D-02 | `setup-check.sh` and its self-test `tests/run.sh` are not yet run by the `pre-commit` hook or by CI. | CI wiring is `T-q344`; no child names the hook. | `T-q344` |
| D-03 | The PSB was cited as `F-0001` before the facts document existed. | Facts were stored by a later child. | Ended by `T-fvwj` ([#4](https://github.com/pharzam/layup/issues/4)): [`F-0001`](../facts/F-0001-layup-problem-statement-brief.md). |
| D-04 | PR #18 was merged while `gh pr checks` showed 3 of its 8 CI jobs pending; all 8 passed. | Human error: the merge did not wait. Lesson in `guardrails.md` §2. | `T-afa5` ([#12](https://github.com/pharzam/layup/issues/12)) makes the jobs required. |

Open gaps: the values with no source are listed in [`open-gaps.tsv`](open-gaps.tsv) with their question (`T-nfh8`). Check `markers` fails when a gap is not listed there.
