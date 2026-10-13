# T-j4sd — row 51b of the plan, the call CommitTree

Issue: [#272](https://github.com/pharzam/layup/issues/272), row 51b of [the tasks
of M2d](../plan/README.md#the-tasks-of-m2d), opened by `T-bng8` (#267). Serves
`F-0003#43` through REQ-003. The plan was made against `9144e4e` (the merge of #305);
the branch took `38ee78e` (the merge of #306, `T-9v92`) by a fast-forward before it
had a commit of its own, and `main` three times more by a merge: `6734a23` (#307, `T-4c5c`)
before the freeze, `4687b39` (#312, `T-p481`) after round 1 and `b28171b` (#311, `T-56dj`) after
round 2. The budget is read against `b28171b`, as the change of this task alone: three comments
"The base of the budget" on #272 (6143368806, 6143768991 and 6144155883) give the figures.
Author: Claude Sonnet 5.5 on Claude Code. Evidence: [`runs/T-j4sd/`](../../runs/T-j4sd/).

## Plan and plan review

The plan (R12), its review and the author's answer are comments on #272. The plan
review (GPT-6.1 Sol on Codex CLI, effort `high`, a fresh session in a clone at
`9144e4e`, 4 min 19 s, 217,752 tokens) gave `approve-with-conditions`: Budget maximum
550 lines added plus removed over 11 files, close-out inside; Cycle cap 1; no panel.
Its three conditions are applied: the green run of every check is a step of the plan
with its record; the exception to the rule of the delimiter and the result of the call
are written into `packages.md`, not left in the code; and the places of the documents
(the Job cell of the table of phase 1, the prose of the table of `M2d`) and the
tally of the files are corrected. The harness `codex` is not one of the three that
ADR-0026 names (`claude`, `devin`, `opencode`): the Operator added Codex CLI to the host
on 2026-10-13 and said it may be used. No test of this task was written or run before the
posted review and the answer (09:25:12Z and 09:25:50Z): the red run is of 09:32Z. The goal
count is 1, final by O-209 of #267. The answer named `git` 2.47.3 as the version that gave
the code 128 for a tree that no object has; no run of 2.47.3 was made, and `git` 2.56.0, the
version of this host, gives 128 (the integration test asserts a code above 0, and the version is
not in the code or the documents); the answer was not edited.

## What was done

1. **The call:** `CommitTree(dir, tree, parent, message, who)` in `internal/git/git.go`
   runs `git commit-tree TREE -p PARENT -m MESSAGE` and gives the ID that `git` prints. It
   refuses a `tree` or a `parent` that is not a full object ID before `git` starts (a
   `*FailedError` of `Code` -1), reads the output as an object ID (else a `*FailedError` of
   `Code` 0), and makes the object only: no ref moves and `HEAD` stays. An empty commit is
   the call with the tree of its parent, which the caller gives.
2. **The identity:** the six variables of the author and the committer are built by one
   function, `commitEnv`, which `Commit` and `CommitTree` share, so "as `Commit` takes
   them" has one home.
3. **The documents:** `packages.md` (the row of `CommitTree`, its refusal and its output,
   the exception to the rule of the delimiter, the empty message, the Job cell of `internal/git`
   and the sentence of the table of `M2d`), `traceability.md`, the Task and Test cells of
   REQ-003, a line of section 13 of `PRD-0001`, and a pitfall in `guardrails.md` §2 (a `git`
   verb that reads its standard input when its argument is empty).

**Tests:** [`test-runs.md`](../../runs/T-j4sd/test-runs.md): red before the code
(`undefined: CommitTree`); each rule of the call and of `commitEnv` mutated alone on a copy
of the tree, 15 of 15 caught, at the unit level and again by the integration level, and the
two mutants of the empty message (note 2 of round 1) caught, the refusal by both tests and the
argument left out by the unit test only (`git` makes the same commit without `-m`); every check of the gate
exits 0 on the merged head before the round, and the checks that the note can change exit 0
after it and after the first merge; the whole script exits 0 again on the head of round 3 (the
integration level of the other packages and the end-to-end level were not run on the two heads in
between: they changed tests of `internal/git`, documents and a script; CI runs them).

**Known limits:** the call does not read the tree of the parent, so a caller that wants an
empty commit gives the parent's tree itself; the call does not check that `tree` names a
tree or that `parent` names a commit, which `git` refuses with its own code (the integration
test asserts a code above 0; `git` 2.56.0 gave 128). An empty message reaches `git` as `-m ''`
and `git` reads the message from its standard input, which `call` leaves at the null device: a
`call` that gave `git` an open standard input would make that call wait. The probe itself, its
push and its checks are row 57a's.

**The rejected alternatives:** reading the parent's tree inside the call (`RevParse` is its
own call, and the caller has the tree); the message on standard input (a call has none); a
ref made by the call (the probe pushes the ID, and the run's clone has no new ref after
the push); a second copy of the identity variables in `CommitTree`; a refusal of an empty
message (the specification names none, and `git` makes the commit).

## Review rounds

Round 1 (GPT-6.1 Sol on Codex CLI, effort `xhigh`, a fresh session in a clone at `4bddf73`,
cycle 0 of the cap of 1; its record is comment 6143600554 on #272): `nothing material in
scope`, two notes. The record names no finding on the budget: it read 285 lines over 7 files
at the head against `6734a23`, 2,928 over 41 against `9144e4e`, and the expected close-out
inside the maximum.

- **Applied** in `3563b56`, a commit that no round read: note 2 (no test passed an empty
  message). The unit test and the integration test now pass `""`; doing it showed that `git
  commit-tree -m ''` reads the message from its standard input (measured: an open pipe made
  it wait), which `call` leaves at the null device, so the commit has an empty message. That is
  a sentence in `packages.md` and a pitfall in `guardrails.md` §2. Two mutants of the rule
  (the empty message refused, and its argument left out), each alone on a copy of the tree,
  are caught: the first by the unit and the integration test, the second by the unit test, which
  reads the arguments (the integration test cannot tell, as `git` makes the same commit).
- **No change** for note 1: the reviewer's sandbox opens no socket and finds the parent
  repository from its temporary directories, so packages of the integration level and
  `TestFlagged` failed there, none by a change of this branch (as the record says, the package
  `internal/git` panicked before `TestCommitTree` ran); the level exited 0 on this host (10:05Z to
  10:16Z, the record of the runs) and CI is the authority.

Round 2 (GPT-6.1 Sol on Codex CLI, effort `high`, a fresh session in a clone at `af5992d`;
comment 6143852255): `origin/main` at `4687b39`, the merge of #312 (`T-p481`, the colors of the
state labels), was merged into the branch after round 1, and changed two files that this branch
had touched, `docs/guardrails.md` and `docs/tasks/completed.md` (the one conflict was the completed
log, where both branches appended a line, kept in the order of landing). So a round scoped to those
files followed, as [the discipline](../engineering-discipline.md#reviewing-until-findings-decay)
says; it consumes no cycle, and its record has the cycle `1 — a round scoped to the merge of
origin/main`. The verdict row of round 1 was edited to `material`, with that reason, so that the
lint reads the two records in their positions (precedent: `T-4c5c`). Round 2 found nothing:
`nothing material in scope`, no finding and no note.

Round 3 (GPT-6.1 Sol on Codex CLI, effort `high`, a fresh session in a clone at `0daa935`;
comment 6144146264): `main` moved again, to `b28171b`, the merge of #311 (`T-56dj`), and the
branch took it by a merge; it changed five files that this branch had touched (`PRD-0001`,
`packages.md`, `traceability.md`, `backlog.md`, `completed.md`), so a round scoped to them
followed, with no cycle. It found nothing: `nothing material in scope`, no finding. A third
numbered record would fail `review-record-lint.sh` (RR7: the k-th record carries cycle k-1, and
none is past the cap of 1), so the reviewer's message is posted under a heading that the check
does not count, with a note that says so, and the check reads rounds 1 and 2. The disagreement
between the rule and the check is issue #313.

## Verdict

Delivered: the call `CommitTree` of `internal/git`, which makes a commit of a tree on a
parent with the identity of `Commit`, on no ref, and gives its ID, and refuses a tree or a
parent that is no full object ID before `git` starts; `commitEnv` is the one home of that
identity. The tests were red first, each rule was mutated alone and caught (15 of 15 at two
levels, and two more after the note), and every check of the gate exits 0. The review ended by
decay at cycle 0 of cap 1, with no material finding, and the round scoped to the merge of `main`
found nothing. The diff against `b28171b` is inside 550
lines over 11 files (536 lines added plus removed over 11 files, measured with the close-out inside; the comment on the base
estimated about 410 lines over 10 files, before the note added `guardrails.md`, its tests and a section
of the record); against `9144e4e`, the base of the plan review, it holds the work of `T-9v92`, `T-4c5c`,
`T-p481` and `T-56dj` too, as the comments "The base of the budget" say. Next: row 57a, `T-wphp`, the step `forge-check`, makes the
probe's commit with this call.

## Resource record

Recorded, not budgeted (ADR-0007). UTC, 2026-10-13; reviewers' tokens are as the harness
reports them; the author's are not reported. The harness exception is the one above.

| Part | Expected tier | Model | Effort | Tokens | Elapsed |
| ---- | ------------- | ----- | ------ | ------ | ------- |
| The plan | reasoning | Claude Sonnet 5.5 | `max` | not reported | 09:13 to 09:15 (the claim to the post) |
| The plan review | reasoning | GPT-6.1 Sol (Codex CLI, `codex exec`) | `high` | 217,752 (as Codex reports; no USD) | 4 min 19 s, 09:16 to 09:20 |
| The tests and the code, test first | execution | Claude Sonnet 5.5 | `max` | not reported | 09:26 to 10:18 (the answer to the pull request) |
| Round 1 | reasoning | GPT-6.1 Sol (Codex CLI, `codex exec`) | `xhigh` | 184,544 (as Codex reports; no USD) | 11 min 16 s, 10:18 to 10:30 |
| The note | execution | Claude Sonnet 5.5 | `max` | not reported | 10:32 to 10:36 |
| The close-out before the merge | execution | Claude Sonnet 5.5 | `max` | not reported | 10:36 to 10:38 |
| The merge of `main` and its checks | execution | Claude Sonnet 5.5 | `max` | not reported | 10:38 to 10:41 |
| Round 2, scoped to the merge | reasoning | GPT-6.1 Sol (Codex CLI, `codex exec`) | `high` | 70,434 (as Codex reports; no USD) | 4 min 58 s, 10:41 to 10:46 |
| The close-out of round 2 | execution | Claude Sonnet 5.5 | `max` | not reported | 10:47 to 10:49 |
| The second merge of `main`, its brief and the full run | execution | Claude Sonnet 5.5 | `max` | not reported | 10:53 to 11:08 |
| Round 3, scoped to the second merge | reasoning | GPT-6.1 Sol (Codex CLI, `codex exec`) | `high` | 64,962 (as Codex reports; no USD) | 5 min 50 s, 10:56 to 11:02 |
| The close-out of round 3 | execution | Claude Sonnet 5.5 | `max` | not reported | 11:02 to 11:09 (the commit of this file) |
