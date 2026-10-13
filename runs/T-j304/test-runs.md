# T-j304 — test runs

Every open issue carries a state label (#303): the changes of `docs/tasks/task-state.sh`
and `docs/tasks/task-labels.sh`, and the cases of `docs/tasks/tests/`. All times are UTC,
2026-10-13, read from the clock after each step. A case gives its contract when a `good`
case exits 0 and a `bad` case exits 1 under `docs/tasks/tests/run.sh`, which gives 2 for
any other difference.

## Red (06:28Z)

The new and the changed cases with the two scripts and `run.sh` as at the base (57 cases):
12 do not give their contract, each for its own difference. `good-title-in-plan` passes at
the base, as the plan says: it guards that a title makes no second task, and has no red.
The base has no `forge/issues` and no `prs`, so a case of the base reads each as empty.

```
FAIL bad-labels-missing-prs (exit 2, want 1)
FAIL bad-missing-issues (exit 2, want 1)
FAIL good-labels-applier-prs (exit 2, want 0)
FAIL good-labels-applier (exit 2, want 0)
FAIL good-labels-closing (exit 2, want 0)
FAIL good-labels-hold (exit 2, want 0)
FAIL good-labels-listed-outside-set (exit 2, want 0)
FAIL good-labels-right (exit 2, want 0)
FAIL good-prs-table (exit 2, want 0)
FAIL good-title-review (exit 2, want 0)
FAIL good-title-tasks (exit 2, want 0)
FAIL good-title-two-issues (exit 2, want 0)
cases: 57, failing: 12
```

After the change of `task-state.sh`, the 7 cases of `task-labels.sh` still fail; after the change of
`task-labels.sh`, none (57 of 57).

## Mutations

Each rule of the two scripts mutated alone, on a copy of the tree and never on the
worktree, then every case of `docs/tasks/tests/` on the copy.

- **The first run** (06:34:43Z to 06:36:56Z): 27 mutants, 23 caught, **4 survived**:
  the number of a row of the plan and the number of a line of the backlog in the rule that a
  title makes no second task (`good-title-in-plan` named the issue of a plan row and of
  a backlog task by lines that another file also named, so each path hid the other), and two
  lines that no case runs (below). The case now has an issue of a row of the plan, and
  one of a backlog task, each under another ID.
- **The second run** (06:37:19Z to 06:39:21Z): 27 mutants, 25 caught, **2 survived**,
  and they are code that no fixture reaches: the `gsub` of the fetcher of `task-state.sh`
  that cleans a title, and the line of the applier that runs the real `task-state.sh` (the
  cases give the stand-in). The live read below runs both.
- **The final runs**, on the tree after a change that `git diff --check` asked for at the
  first commit (`task-state.sh --prs` writes no third field for a pull request that closes
  none, so no row of a fixture ends in a tab): 27 mutants again (06:59:39Z to 07:01:25Z), 25
  caught, 2 survived, the same two; then with `ts-prs-no-third-field` added (07:01:40Z to
  07:03:29Z), the run below: 28 mutants, 26 caught, **2 survived**, the same two.

```
== ts-title-colon: caught by 1 case(s), the first good-title-tasks
== ts-title-id-named: caught by 1 case(s), the first good-title-in-plan
== ts-title-issue-named: caught by 1 case(s), the first good-title-in-plan
== ts-title-lowest: caught by 1 case(s), the first good-title-two-issues
== ts-rows-issue-named: caught by 1 case(s), the first good-title-in-plan
== ts-rows-id-named: caught by 1 case(s), the first good-title-in-plan
== ts-now-id-named: caught by 1 case(s), the first good-title-in-plan
== ts-now-issue-named: caught by 1 case(s), the first good-title-in-plan
== ts-order-numeric: caught by 1 case(s), the first good-title-tasks
== ts-required-issues: caught by 1 case(s), the first bad-missing-issues
== ts-prs-written: caught by 1 case(s), the first good-prs-table
== ts-prs-no-third-field: caught by 1 case(s), the first good-prs-table
== ts-issues-fetch: SURVIVED
== tl-required-prs: caught by 13 case(s), the first bad-labels-disjoint
== tl-prs-numbers: caught by 2 case(s), the first good-labels-applier-prs
== tl-closed-stale: caught by 1 case(s), the first good-labels-hold
== tl-closed-untouched: caught by 1 case(s), the first good-labels-hold
== tl-stale-regex: caught by 1 case(s), the first good-labels-hold
== tl-closed-reason: caught by 3 case(s), the first good-labels-each
== tl-listed-any: caught by 1 case(s), the first good-labels-listed-outside-set
== tl-in-review: caught by 2 case(s), the first good-labels-applier-prs
== tl-on-hold: caught by 6 case(s), the first good-labels-applier-prs
== tl-seen: caught by 1 case(s), the first bad-labels-disjoint
== tl-applier-prs-arg: caught by 1 case(s), the first good-labels-applier-prs
== tl-applier-prs-real: SURVIVED
== tl-applier-prs-empty: caught by 4 case(s), the first bad-labels-disjoint
== tl-eighth-label: caught by 2 case(s), the first good-labels-applier-prs
== tl-label-text: caught by 2 case(s), the first good-labels-applier-prs
26 caught, 2 survived, of 28
```

## The live read (07:02:20Z to 07:02:33Z)

`sh docs/tasks/task-labels.sh --dry-run` on the forge, with the scripts of this branch,
exit 0, writes nothing. These are the changes that the first run on `main` will make: the
30 open issues with no label take `state: on hold`, and #303, which only its title names
as a task, takes the state of its branch. (14 open issues have a label already, each
right.)

```
add #68 state: on hold
add #178 state: on hold
add #186 state: on hold
add #189 state: on hold
add #195 state: on hold
add #205 state: on hold
add #207 state: on hold
add #230 state: on hold
add #232 state: on hold
add #235 state: on hold
add #236 state: on hold
add #238 state: on hold
add #240 state: on hold
add #242 state: on hold
add #244 state: on hold
add #245 state: on hold
add #247 state: on hold
add #249 state: on hold
add #251 state: on hold
add #252 state: on hold
add #254 state: on hold
add #256 state: on hold
add #258 state: on hold
add #260 state: on hold
add #262 state: on hold
add #264 state: on hold
add #291 state: on hold
add #296 state: on hold
add #298 state: on hold
add #300 state: on hold
add #303 state: in progress
```

## Green

**A first run** (06:39:30Z to 06:46:19Z) had every line at exit 0 but the level `integration`,
which exited 1: `internal/run` and `internal/verify` stopped at a write with "No space left
on device", because `/tmp` of the host, a tmpfs of 1.9 GB, was full of the clones of earlier
reviews. No test failed for a rule. The clones of merged tasks were removed (the free space
went from 46 MB to 975 MB) and the whole run was made again, with `TMPDIR` on the disk
(06:49:19Z to 06:56:20Z): every line at exit 0.

**A change followed that run.** `git diff --check` read the tracked files only, so it had
not seen the 11 new case directories, and at the first `git add` six fixture files showed
a trailing tab (a row whose last field is empty). `task-state.sh --prs` now writes no third
field for a pull request that closes none, the six files lost their tabs, one mutant was
added for the new branch, and the mutations and the live read above were made again.

**The final run on `012bba3`** (07:06:44Z to 07:14:27Z), each with exit 0, in this order; the
time of a line is when its command ended. The three runs of `go test` are with `-count=1`,
then with `-tags=integration` and with `-tags=e2e`, and `git diff --check` reads the
commits against `origin/main`:

```
go build ./...: exit 0 (07:06:46Z)
go vet ./...: exit 0 (07:06:46Z)
gofmt -l internal cmd (empty): exit 0 (07:06:46Z)
go test ./...: exit 0 (07:06:54Z)
go test -tags=integration ./...: exit 0 (07:11:34Z)
go test -tags=e2e ./...: exit 0 (07:13:37Z)
adr-lint: exit 0 (07:13:38Z)
prd-lint: exit 0 (07:13:38Z)
link-lint: exit 0 (07:13:52Z)
run-discipline-tests: exit 0 (07:14:00Z)
nested-checkout-check: exit 0 (07:14:00Z)
setup-check: exit 0 (07:14:27Z)
git diff --check origin/main HEAD: exit 0 (07:14:27Z)
```

**After the merge of `origin/main`** (`8bd6943`, the merge of #304, taken into the branch as
`129733a`) the checks that read the tree were made again, from 07:14:59Z to 07:16:05Z, each
with exit 0; the levels `integration` and `e2e` were not (the change has no Go code, and
CI runs them on the head of the pull request):

```
go build ./...: exit 0 (07:15:01Z)
go vet ./...: exit 0 (07:15:06Z)
gofmt -l internal cmd (empty): exit 0 (07:15:06Z)
go test ./...: exit 0 (07:15:15Z)
adr-lint: exit 0 (07:15:17Z)
prd-lint: exit 0 (07:15:17Z)
link-lint: exit 0 (07:15:30Z)
run-discipline-tests: exit 0 (07:15:37Z)
nested-checkout-check: exit 0 (07:15:38Z)
setup-check: exit 0 (07:16:05Z)
git diff --check origin/main HEAD: exit 0 (07:16:05Z)
```

## After round 1 (07:33:13Z to 07:34:11Z)

Round 1 gave `nothing material in scope` and eight notes; the notes commit `b87fd68` changes
documents, the usage text of a script and one fixture, and no code.

**One mutant for the new case** (note 7), on a copy of `docs/tasks` and never on the worktree
(07:32:32Z): in the applier of `task-labels.sh`, `[ -n "$_rem" ] && set -- ...` joined with
`&& [ -n "$_add" ]`, so that an edit with a removal and no addition drops the removal. It is
caught by `good-labels-applier` alone (exit 2, the call `issue edit 10 --remove-label
state: ready` missing); `good-labels-applier-prs` and `good-labels-applier-dry-run` stay at 0.
Before the case had #10 the three passed with the mutant.

**The checks that read the tree**, on `b87fd68`, each with exit 0 (the levels `integration`
and `e2e` are not repeated: the notes touch no Go code, and CI runs them on the head of the
pull request):

```
go build ./...: exit 0 (07:33:14Z)
go vet ./...: exit 0 (07:33:15Z)
gofmt -l internal cmd (empty): exit 0 (07:33:15Z)
go test ./...: exit 0 (07:33:21Z)
adr-lint: exit 0 (07:33:23Z)
prd-lint: exit 0 (07:33:23Z)
link-lint: exit 0 (07:33:36Z)
run-discipline-tests: exit 0 (07:33:43Z)
nested-checkout-check: exit 0 (07:33:44Z)
setup-check: exit 0 (07:34:10Z)
git diff --check origin/main HEAD: exit 0 (07:34:11Z)
```
