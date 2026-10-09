# T-eep8 — the state of each task, derived from the plan and the forge

Issue: [#175](https://github.com/pharzam/layup/issues/175), opened by the
Operator's request in the session `issue-status`. Serves `F-0003#52` (the task
state in a form that does not belong to one harness agent), by the Operator's
placement on #175 (comment 6081029848, ADR-0026 rule 1). Base `ad8774d` (the
merge of #173). Author: Claude Opus 5.5 on Claude Code. Evidence:
[`runs/T-eep8/`](../../runs/T-eep8/). The branch `T-eep8` was pushed at the base
before the plan, as the rule of this task asks.

## Plan and plan review

The plan (R12, comment 6081102065), its review (6081280764) and the author's
answer (6081281115) are comments on #175. The plan review (Claude Fable 5.1,
effort `xhigh`, on the Claude Code CLI with stream output, a fresh read-only
session in a clone at `ad8774d`, 10 min 15 s) gave `approve-with-conditions`:
Budget maximum 900 lines added plus removed over 140 files against `ad8774d`,
close-out inside; Cycle cap 2 (a dispatch line of a check script and a rule of
gate step 1); no panel. Its four conditions are applied: `done` is printed, over
the rows of the last task table and the backlog tasks with no row there; the
cases are `good-*` and `bad-*` with `run.sh` as the runner's linter, exiting 0,
1 or 2; a branch on the forge takes `origin/main` by a merge; the mirrors of the
check set say the suite is a command's, not a sixth linter's. The goal count is
2 (the command; the claim rule).

## What was done

1. **D1:** `docs/tasks/tests/`: 15 cases over a shared `base/`, each holding
   only the files it changes, run by `run.sh` as the suite `task-state` of
   `run-discipline-tests.sh`; a suite README lists them.
2. **D2:** `docs/tasks/task-state.sh`. The reader takes `backlog`, `completed`,
   `plan`, `forge/branches` and `forge/pr/<n>/head` and `body` from a directory
   (`--read DIR`), with carriage returns removed; the fetcher copies the three
   files of the checkout and reads the forge with `gh` (one progress line per
   step). The states, the first that holds: `done` (the ID field of a line of
   the completed log), `in review` (a pull request whose head is the task ID, or
   whose body closes its issue with one of the forge's nine keywords, in any
   case), `running` (a branch named exactly the task ID), `blocked`, `ready`.
   `.gitattributes` pins it at line feeds, and the case `good-crlf` at carriage
   returns.
3. **D3:** the claim rule in [Starting a task](../engineering-discipline.md#starting-a-task);
   a branch on the forge merges `origin/main` in
   ([Integrating branches](../engineering-discipline.md#integrating-branches),
   and `docs/ci/README.md` on `strict: true`); the row of the enforcement table;
   `AGENTS.md` (step 1, the command, the branches); the onboarding;
   `docs/tests/README.md`, `test-levels.md`, `docs/ci/README.md`; `README.md`;
   the glossary row **Claim**.

**Tests:** [`test-runs.md`](../../runs/T-eep8/test-runs.md): red against a stub
(all 13 first cases exit 2, an output that differs from `EXPECT`); green; 15
mutations, each failing a case; after round 1, `good-crlf` red then green and
`good-columns-order` pinned by two mutations; a run of the fetcher on this
repository (18 lines, rows 28 to 39b, 1.3 s).

**Two decisions taken while building**, recorded on #175 (Fixes — round 1):
the branches are read with `gh api repos/{owner}/{repo}/branches`, not `git
ls-remote --heads origin` as the plan said, so one tool reads the whole forge
(on this host `origin` is an SSH URL with no key, and `git ls-remote` fails);
and the check `adapted` is kept off the fixtures by inputs with no extension,
not by the exemption lists of `setup-check.sh`, which is unchanged (the
answer's note 5 said the lists; the Fixes reply corrects it).

**Known limits**, each in [#182](https://github.com/pharzam/layup/issues/182):
a backlog task whose row is in an earlier task table is never `blocked`; the
rule does not say whether the plan phase is claimed; the fetcher reads the
checkout's copies of the three files, so a checkout behind `origin/main` gives a
stale answer. The command enforces nothing: a task whose branch is not pushed is
invisible to it, and a branch left after an abandoned task reads `running`.

**The rejected alternatives:** labels or assignees set by each session; a
`layup` subcommand in Go; `git worktree list` (one host only); a CI job that
writes labels (a second goal, left to a child issue if wanted).

## Review rounds

The records and the Fixes reply are comments on #175. Round 1 (Claude Fable
5.1; `5517013`, cycle 0): `material`, three findings (the branch not named the
task ID in the rule's home; two readings of "take only one that it prints as
`ready`"; the suite and the reader broken by carriage returns) and six notes,
fixed in `1185d2c`. Round 2 (Fable; `1185d2c`, cycle 1): `nothing material in
scope`, four notes: note 2 (the set of tasks the command does not print is the
backlog's `## Now`) is applied in the close-out; notes 1, 3 and 4 change the
command or the rule, so they are #182.

## Verdict

Delivered: the command that prints the state of each task, its fixture suite,
and the claim rule that makes a started task visible to it. The review ended by
decay at cycle 1, under the cap of 2. The diff against `ad8774d` is inside 900
lines over 140 files. Next: #182; PR #177 (the backlog line of this task) is
superseded by this pull request, which records the task in the completed log.

## Resource record

Recorded, not budgeted (ADR-0007). Times are UTC, 2026-10-09; the reviewers'
tokens are `modelUsage` of the CLI's `result` event; the author's are not
reported.

| Part | Expected tier | Model | Effort | Tokens | Elapsed |
| ---- | ------------- | ----- | ------ | ------ | ------- |
| The plan | reasoning | Claude Opus 5.5 | not reported | not reported | 12:38 to 12:42 |
| The plan review | reasoning | Claude Fable 5.1, Claude Code CLI | `xhigh` | 760,919 (USD 6.03) | 10 min 15 s, from 12:42 |
| The answer; the tests, red; the code; the documents | execution | Claude Opus 5.5, a reasoning-tier model on an execution part | not reported | not reported | 12:53 to 13:01 |
| Round 1 | reasoning | Claude Fable 5.1, Claude Code CLI | `xhigh` | 1,810,297 (USD 6.84) | 15 min 14 s, from 13:02 |
| The fix of round 1 | execution | Claude Opus 5.5 | not reported | not reported | 13:17 to 13:20 |
| Round 2 | reasoning | Claude Fable 5.1, Claude Code CLI | `xhigh` | 1,102,852 (USD 5.52) | 11 min 23 s, from 13:20 |
| The close-out | execution | Claude Opus 5.5 | not reported | not reported | 13:33 to 13:38 |

No harness was skipped: the first in the reviewer order (`claude`) returned a
record each time.
