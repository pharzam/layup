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

## The fixes of round 3

Review round 3 (`e4aabfe`, GPT-6 Sol on Devin) ended `not mergeable, findings
recorded` at the raised cap of 3, with four findings. The two code findings
are fixed in `ea43afe`; the review bookkeeping is done on the issue; the
budget stays the Operator's.

### Red

`bad-labels-workflow-tags-ignore` (new): a `push` bound `tags-ignore: ['*']`
beside `branches: [main]` — the key pattern `[a-z_]+` of round 2's fix could
not see a hyphenated bound, so `--check-workflow` printed `the workflow OK`;
`run.sh` exit 2, the fail for the finding's reason.

### Green

The workflow keys are recorded with `[a-z_-]+`, so `tags-ignore` and
`branches-ignore` are seen and rejected like `tags` and `branches`:
`bad-labels-workflow-tags-ignore`, `bad-labels-workflow-tags`,
`bad-labels-workflow-push` and `bad-labels-workflow-branches` each fail on
their own bound; `good-labels-workflow` and `good-labels-right` pass. The
workflow header, `glossary.md` and the `task-labels.sh` header now say an
open indexed issue absent from `task-state.sh` carries no `state:` label and
a closed issue's label comes from its close reason — the same claim
`AGENTS.md` was updated with in round 1. `run-discipline-tests`: 126 passed,
0 failed; `adr-lint`, `prd-lint`, `link-lint`, `setup-check`, `go build`,
`go vet`, `go test` (unit and `-tags=integration`), `gofmt -l` (empty) and
`git diff --check`: exit 0.

### The budget, measured

`git diff --shortstat a4d41cf` of this head: **1,609 lines, 95 files** —
recorded after the commit below.

## The fixes of round 4

Review round 4 (`8e7c48a`, GPT-6 Sol on Devin) ended `not mergeable, findings
recorded` at the raised cap of 4, with two material findings and a note. The
note is fixed here; the issue's own wording was fixed on the forge; the budget
stays the Operator's.

- Finding 1 (material, open): the reviewer measured `git diff --shortstat
  a4d41cf 8e7c48a` at **1,643 lines, 95 files** — the 1,609 recorded above
  missed the 34 lines of that evidence commit itself.
- Finding 2 (material, fixed on the forge): the issue's body still said "The
  six labels are created once" while the re-plan's table names seven (`stalled`
  from #202). The body now says seven.
- The note (fixed): `docs/ci/README.md`'s `discipline-tests` row now names the
  `task-labels.sh` fixtures beside `task-state.sh`'s.

### The budget, measured

`git diff --shortstat a4d41cf` of this head: **1,666 lines, 95 files** —
measured on the tree of this commit, the reviewer's 1,643/95 plus the 23 lines
of the note's fix and this section.

## The fixes of round 5

Review round 5 (`ea9aa4e`, GPT-6 Sol on Devin) ended `not mergeable, findings
recorded` at the raised cap of 5, with two material findings. One is fixed on
the forge; the budget stays the Operator's.

- Finding 1 (material, open): `git diff --shortstat a4d41cf ea9aa4e` reports
  **1,667 lines, 95 files** — the figure of the last section was one line
  short of the commit's own measure; it is left as recorded there.
- Finding 2 (material, fixed on the forge): the issue's Goal still said "Each
  task issue carries one state label" and the labels section "Exactly one
  `state:` label per task issue" — both miss the re-plan's exception that an
  open issue the command does not name carries none. The Goal now reads "Each
  task issue the command names carries one `state:` label … (an open issue it
  does not name carries none)", and the labels section "At most one `state:`
  label per task issue — an open issue the command does not name carries
  none".

### The budget, measured

`git diff --shortstat a4d41cf` of this head: **1,690 lines, 95 files** —
measured on the tree of this commit.

## The fixes of round 6

Review round 6 (`5acd43c`, GPT-6 Sol on Devin) ended `not mergeable, findings
recorded` at the raised cap of 6, with three material findings. Two are fixed
here; the budget stays the Operator's.

- Finding 1 (material, open): `git diff --shortstat a4d41cf 5acd43c` reports
  **1,690 lines, 95 files** — the overrun stays unapproved; it is the
  Operator's decision, not an implementation defect.
- Finding 2 (material, fixed): the `docs/ci/README.md` account of
  `task-state.yml` named neither qualification the plan review's note 2 asks
  for. It now says the task issues are the ones the command lists from the
  task tables, the completed log and `## Now` of `main`; that an open issue
  it does not list carries none; and that a closed issue takes
  `state: merged` when `state_reason` is `completed` — which stands in for
  "merged by its pull request" — and `state: rejected` otherwise.
- Finding 3 (material, fixed): a quoted push bound `'tags': ['*']` is valid
  YAML beside `branches: [main]`, but the checker's key pattern wanted an
  unquoted name and skipped the line, so the run printed `the workflow OK`.
  The bound keys of `on:` are now matched quoted or unquoted (single or
  double quote, and a space before the colon), and the name is recorded
  de-quoted, so a quoted `branches` still passes and every other bound fails.
  The mutant `bad-labels-workflow-tags-quoted` carries `'tags': ['*']`
  and fails with `push: tags: a bound of push is branches only`.

### The budget, measured

`git diff --shortstat a4d41cf` of this head: **1,803 lines, 98 files** —
measured on the tree of this commit.

## The fixes of round 7

Review round 7 (`8d1b77b`, Claude Fable 5.1 on Claude Code, effort `max`)
ended `not mergeable, findings recorded` at the raised cap of 7, with two
material findings and four notes. The material doc finding and three notes
are fixed here; the budget stays the Operator's, and one note is recorded as
a limit.

- Finding 1 (material, open): the landing figure, **1,803 lines / 98 files**
  at `8d1b77b`, is past every approved number; the Operator's approval is the
  only settling, and the head below measures **1,858 lines / 98 files**.
- Finding 2 (material, fixed): the new `docs/ci/README.md` sentence named
  `task-state.sh` as listing the three files; that command prints the last
  task table's rows and `## Now` only, and `task-labels.sh` is the reader of
  all three files. The sentence now says the task issues are the ones the
  task tables, the completed log and `## Now` of `main` name, which
  `task-labels.sh` reads.
- Note 3 (fixed): `DEFAULT_BRANCH` was `${{ github.event.repository.default_branch }}`,
  which a `schedule` run does not carry, so `git clone --branch ""` would
  fail under `set -eu` on every scheduled run. It is now
  `${{ github.event.repository.default_branch || 'main' }}`; the checkout's
  `ref:` is unchanged (an empty `ref:` falls back to the run's ref, the
  default branch on a schedule).
- Note 4 (fixed): the stated reason for dropping `delete` was wrong —
  `delete` has no ref left, so its run uses the default branch's file; it is
  dropped because a run on a deletion adds nothing the schedule gives. Said
  so in the workflow header, `docs/ci/README.md` and `docs/guardrails.md`.
- Note 5 (fixed): `--check-workflow` now names its limits in a `LIMIT`
  comment, as the checks of `setup-check.sh` do: a second checkout of
  another `ref:`, an event field outside a `run:` step, and YAML's
  explicit-key bound form pass it; each needs an edit of `main`'s file,
  which review of `.github/**` covers.
- Note 6 (fixed): the `good-labels-right` row of the tests README described
  a `duplicate`/`in progress` cross the case does not hold; it now names the
  `not_planned`/`in progress` cross it does hold, and the `duplicate` close
  as `rejected` besides.

### The budget, measured

`git diff --shortstat a4d41cf` of this head: **1,858 lines, 98 files** —
measured on the tree of this commit.
