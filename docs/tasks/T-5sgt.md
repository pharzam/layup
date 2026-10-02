# T-5sgt — `layup gate REPO --base REV --head REV`

Issue: [#82](https://github.com/pharzam/layup/issues/82), row 5 of the
[implementation plan](../plan/README.md); child of `layup gate` (#33). Serves
`F-0003#44` and `F-0003#47`. Base `f394433` (the merge of row 4, #104). Author:
Claude Opus 5.5 on Claude Code. Evidence: [`runs/T-5sgt/`](../../runs/T-5sgt/).

## Plan and plan review

The plan (R12), its review and the author's answer are comments on #82. The plan
review (Claude Fable 5.1 on Claude Code, a fresh session) gave
`approve-with-conditions`: Budget maximum 3,000 lines added plus removed over 32
files against the base, close-out inside; Cycle cap 1. The author applied the
condition (a reason that a failure gives is a fixed text, never a part of an
error, which can name a scratch path) and the ten notes, with one change of D4
(note 4: a `pending` row keeps its result from the diff when the scratch tree
fails), and added one decision: the overlay writes through an `os.Root` of the
scratch tree.

## What was done

1. **Test first** ([`test-runs.md`](../../runs/T-5sgt/test-runs.md)): one
   measurement of `git ls-tree`; the red run of each test before its code (the
   tests of `LsTree`, the unit and integration tests of `internal/gate`, the
   command's tests in `internal/cli`, the end-to-end scenarios on the binary of
   the base). Two tests that passed on the skeleton were made strict, one with a
   control run of plain `git`.
2. **`internal/git`:** the call `LsTree`, added first to the call table of
   `docs/spec/packages.md`.
3. **`internal/gate`:** the scope patterns (D1), the manifest reader (D2), the
   result rules (D5, D6), the scratch tree with the overlay of the base's gate
   files through an `os.Root` (D3), and the run, which gives input errors,
   fixed reasons for a failure (D4) and the held block of each kind (D7);
   `ManifestSchema` and `ResultSchema`, each compared with its block;
   `gate-result` moves to `built`.
4. **`internal/cli`:** the row `gate REPO --base REV --head REV`, with its
   progress lines and its exit codes (D8).
5. **`cmd/layup`:** the end-to-end scenarios on a small Go repository made in the
   test (the demo, the exit codes 0, 1 and 2, the repeat rule, the fixture of
   `NFR-004`), with the host's `GOCACHE` (note 7).
6. **The specification:** D1 to D8, the fixed reasons, the known limits and K19
   in `docs/spec/gate.md`; the Progress bullet, the exit map and the row of code 2
   in `docs/spec/README.md`; the row `LsTree` in `packages.md`; the rule for the
   authors of a fixture in `setup.md`.
7. **The other documents:** the commands in `README.md` and in the onboarding
   document; a lesson in `docs/guardrails.md` §2 (a reason that copies an error
   into a table); the traceability rows (the planned row of this task becomes
   eight rows of real tests); the `PRD-0001` §12 Test cells of `REQ-004`,
   `REQ-007` and `NFR-004`, with a §13 row.

**The rejected alternatives:** the rows of a failed scratch tree as `fail`
(a check that did not run is `not-active`, `NFR-004`); the text of an error in a
reason (it can name a scratch path, and the repeat rule breaks); a `git worktree
prune` at the start (it could remove another stale record of `REPO`); a clean
environment for a gate command (a target's tool needs `HOME` and its caches);
`git show` alone for the overlay (it gives no mode, and a directory gives a
listing, not files).

## Review round 1 and its fixes (cycle 1)

Review round 1 (Claude Fable 5.1, `a0435a6`, cycle 0) gave `material`, with
three material findings, all fixed: the overlay took the `config` paths in the
order of a Go map, so two overlapping paths could give two verdicts on one
input (now the sorted order, and a path under a file of the tree counts as
absent); a `config` path `.git` removed the `.git` file of the scratch tree and
left a work-tree record in `REPO` (now an input error); the `PRD-0001` §12 Test
cell of `NFR-005` was planned and not filled (now filled, with the §13 row).
Notes 4, 5, 6, 7 and 8 are applied with the fixes, so round 2 reads them: the
tests of the run with a fake of `internal/git` are at the integration level
(they write temporary files); `gate.md` names a symbolic link or a file of the
head at or above a `config` path, `exit -1` when `sh` does not start, and the
form `--base=-x`; a kind with no scope pattern is an input error, because it
could never run and would always be `clear`. Note 9 is a measurement and needs no
change.

## Review round 2, O-130 and the fix (cycle 2)

Review round 2 (Claude Fable 5.1, `bab6759`, cycle 1) found findings 1 to 3 of
round 1 fixed, and one new material finding: on a file system that folds case,
the `config` path `.GIT` is the file `.git` of the scratch tree, so the overlay
removed it and `REPO` kept a work-tree record. It gave `not mergeable, findings
recorded` under the cap of 1.

**O-130.** The Operator answered "#82 a" in the author's Claude Code session
(2026-10-02), copied to #82: the Cycle cap of #82 is 2. A `## Plan review`
comment on #82 records the new cap, and the verdict of round 2 is edited to
`material`, with the reason, because a fix followed (the lesson of `T-0drh`,
O-120).

The fix: a `config` path in `.git` in any case is an input error, and the
overlay never removes a path that is the same file as `.git`. Notes 2 and 3 of
round 2 are applied with it: `gate.md` says that a removal follows a symbolic
link that stays inside the tree, and that a kind with no scope pattern is a rule
for the catalog entries too.
