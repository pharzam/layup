# T-dm8m — the test runs

All runs on the branch `T-dm8m`, base `a4d41cf`, on 2026-10-12, with `mawk 1.3.4`.

## Red

At 06:04 UTC, the sixteen `*-labels-*` cases against a stub
`task-labels.sh` that prints nothing: each fails for the right reason (`run.sh`
exits 2, the output or the stderr differs from `EXPECT`), and the runner's
case-name check asks the pin of `good-labels-crlf`.

```
FAIL  task-state/bad-labels-applier-issues-fail (wanted exit 1, got exit 2)
FAIL  task-state/bad-labels-applier-state-fails (wanted exit 1, got exit 2)
FAIL  task-state/bad-labels-disjoint (wanted exit 1, got exit 2)
FAIL  task-state/bad-labels-empty-set (wanted exit 1, got exit 2)
FAIL  task-state/bad-labels-missing-issues (wanted exit 1, got exit 2)
FAIL  task-state/bad-labels-unknown-state (wanted exit 1, got exit 2)
FAIL  task-state/bad-labels-workflow-dispatch (wanted exit 1, got exit 2)
FAIL  task-state/bad-labels-workflow-event (wanted exit 1, got exit 2)
FAIL  task-state/bad-labels-workflow-ref (wanted exit 1, got exit 2)
FAIL  task-state/bad-labels-workflow-restore (wanted exit 1, got exit 2)
FAIL  task-state/bad-labels-workflow-write (wanted exit 1, got exit 2)
FAIL  task-state/good-labels-applier-dry-run (wanted exit 0, got exit 2)
FAIL  task-state/good-labels-crlf (wanted exit 0, got exit 2)
FAIL  task-state/good-labels-each (wanted exit 0, got exit 2)
FAIL  task-state/good-labels-right (wanted exit 0, got exit 2)
FAIL  task-state/good-labels-workflow (wanted exit 0, got exit 2)
FAIL  crlf: docs/tasks/tests/good-labels-crlf/ needs Windows line endings and no .gitattributes rule pins them — copy the kit .gitattributes into your repository root
run-discipline-tests: 102 passed, 17 failed
```

## Green

At 06:06 UTC: `run-discipline-tests: 118 passed, 0 failed`; at 06:10, with the two
cases the mutations asked for (below), `119 passed, 0 failed`.
`sh docs/tasks/task-labels.sh --check-workflow .github/workflows/task-state.yml`:
`task-labels: the workflow OK`. `sh docs/setup/setup-check.sh`: every check OK
(check `protection` leaves the new workflow out by #201; check `ci` reads its
`Restore` step).

## Each rule fails a case on its own

Each line is one change of `task-labels.sh`, applied to a copy of the tree
outside it, and the cases that then fail. Two rules first had no case
(`the completed log not read`, `expressions in run: allowed`); `good-labels-each`
gained an issue only the completed log names, and `bad-labels-workflow-event-block`
was added; each is red under its mutation.

| Mutation | Cases that fail |
| -------- | --------------- |
| not_planned read as merged | good-labels-each good-labels-right |
| closed issue takes the command's state | good-labels-applier-dry-run good-labels-crlf good-labels-each good-labels-right |
| an unlisted open issue keeps its label | good-labels-right |
| labels outside state: removed | good-labels-each good-labels-right |
| issues outside the set labelled | bad-labels-disjoint good-labels-right |
| the completed log not read | good-labels-each |
| no empty-set guard | bad-labels-empty-set |
| no disjoint guard | bad-labels-disjoint |
| done not mapped to merged | good-labels-each |
| dry run writes the labels | good-labels-applier-dry-run |
| labels created before the issue read | bad-labels-applier-issues-fail good-labels-applier-dry-run |
| workflow_dispatch allowed | bad-labels-workflow-dispatch bad-labels-workflow-event-block bad-labels-workflow-event bad-labels-workflow-ref bad-labels-workflow-restore bad-labels-workflow-write good-labels-workflow |
| expressions in run: allowed | bad-labels-workflow-event-block |
| no ref check | bad-labels-workflow-ref |

## The fixes of round 1

On 2026-10-12, against the head `aafee79` that round 1 reviewed (reviewer GPT-6
Sol on Devin; verdict `material`, seven findings). The author of the fixes is
SWE-2 Max on Devin, continuing the task after the first author's usage limit.

### Red first

The new cases against `aafee79`: `good-labels-applier` exits 2 — the stand-in
records `issue edit 5 --add-label state: ready` where the removal belongs (the
empty second field of the edits file collapsed in `read`, and a removal was
applied as an addition); `bad-labels-workflow-push`, `-branches`, `-types` and
`-issues-types` exit 2 — the check answered `the workflow OK` on each bound it
did not read. `bad-labels-empty-set` and `bad-labels-disjoint` read now through
the applier and pass, their `--- calls` recording the reads and no write.

### Each new rule fails a case on its own

Each line is one change of `task-labels.sh`, applied to a copy of the tree
outside it, and the cases that then fail (exit 2 from `run.sh`).

| Mutation | Cases that fail |
| -------- | --------------- |
| labels created before the issue read | bad-labels-empty-set bad-labels-disjoint good-labels-applier good-labels-applier-dry-run |
| a removal applied as an addition (no `-` field) | good-labels-applier |
| the push branches bound not checked | bad-labels-workflow-push bad-labels-workflow-branches |
| the event types bounds not checked | bad-labels-workflow-types bad-labels-workflow-issues-types |

### Green

`run-discipline-tests`: 124 passed, 0 failed, the 119 of `aafee79` plus the five
new cases. `git diff --check` on the branch diff: clean (the ten fixture lines
that ended in an empty field now drop it).

## The applier on this repository, a dry run

At 06:06 UTC, `sh docs/tasks/task-labels.sh --dry-run` in this worktree: the
progress lines of `task-state.sh` and of the reads, 9.6 s, exit 0, and no write
(the forge held no `state:` label before and after). The changes it would make:
111, of which 110 `add … state: merged` (the closed and done task issues of the
plan, the completed log and the backlog) and one `add #228 state: in progress`
(`T-ekh5`, the demo of `M2c`, claimed by another session with recent activity).
No indexed task issue was closed as not planned, so no `state: rejected`.

## The fixes of round 2

Review round 2 (`33156a2`, GPT-6 Sol on Devin) ended `not mergeable, findings
recorded` at the cap of 2, with four material findings. Three are fixed here;
the fourth, the budget, is the Operator's.

### Red

`bad-labels-workflow-tags` (new): a `push` bound `tags: ['*']` beside
`branches: [main]` — `--check-workflow` printed `the workflow OK`; `run.sh`
exit 2, the fail for the finding's reason. `good-labels-right` gains a `plan`
naming #7 and the issue closed with an empty `state_reason`: the reader
already sent it to `state: rejected`, so the case adds the promised coverage
(D1) without a code change.

### Green

The check now rejects a `push` bound other than `branches`, and the mutant
fails on it alone. `glossary.md` and the `task-state.yml` comment now name the
indexed set and the close-reason precedence. `run-discipline-tests`: 125
passed, 0 failed.

### The budget, measured

`git diff --shortstat a4d41cf` of this head: **1,534 lines, 92 files** —
recorded after the commit below.
