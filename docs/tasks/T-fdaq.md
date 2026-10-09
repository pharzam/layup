# T-fdaq — the build tasks of M2b

Issue: [#152](https://github.com/pharzam/layup/issues/152), the slicing task that
`T-ywk7` (#147) named, under the Operator's decision O-164. Serves `F-0003#51`
and the requirements of `M2b`. Base `faeb1dd` (the merge of #149). Author: Claude
Opus 5.5 on Claude Code. Evidence: [`runs/T-fdaq/`](../../runs/T-fdaq/).

## Plan and plan review

The plan (R12, comment 6079120704), its review and the author's answer
(6079317854) are comments on #152. The plan review (Claude Fable 5.1, effort
`xhigh`, on the Claude Code CLI with stream output, a fresh read-only session in
a clone at `faeb1dd`, 10 min 54 s; comment 6079310864) gave
`approve-with-conditions`: Budget maximum 750 lines added plus removed over 14
files against `faeb1dd`, close-out inside; Cycle cap 1; no panel. Its three
conditions are applied: the honest goal count of each row went to the Operator
again; `parts.tsv` holds each clause of an acceptance row that the slicing cuts,
and `parts.sh` checks the rows of "The acceptance tests of M2b"; the duty of the
build tasks for `release-check.sh` is dropped. Its notes are applied as the
answer says.

**The Operator's decisions:** O-186 (comment 6079171344) split rows 30, 33, 34
and 36; O-187 (6079376251), after condition 1, split rows 37 and 39 too, and made
the goal count of each row final.

## What was done

1. **D1:** [The tasks of M2b](../plan/README.md#the-tasks-of-m2b), rows 28 to
   39b, one per package boundary of the specification of `M2b`, in the order of
   its imports, with the goal count of each row in its issue; the stale
   sentences of the plan, the milestone row, the item hosts and the row of #148.
   [`parts.tsv`](../../runs/T-fdaq/parts.tsv) puts each part of the
   specification, and each clause of an acceptance row that the slicing cuts, in
   one row, or names it with a reason; [`parts.sh`](../../runs/T-fdaq/parts.sh)
   failed before the table and passes after.
2. **Two moves in the specification:** the extension of `TestPackageRules` to
   the table of `M2b` is row 29's, before any package of that table exists (as
   O-170 did for `M2a`), in `session.md`, `packages.md` and `gate.md`; the review
   of the release of `M2b` is row 39a's, before the demo (O-187), in `session.md`.
3. **D2:** issues #153 to #170, opened by the App through
   [`gen-issues.py`](../../runs/T-fdaq/gen-issues.py), copied from `T-zwke`'s and
   changed for `M2b` (it reads the parts of each row from `parts.tsv`, and
   renders the table from the same rows); row 39b's issue lists the Operator's
   inputs (R6), and row 39a's takes #148.
   [`issues.sh`](../../runs/T-fdaq/issues.sh) failed before the issues and passes
   after (it reads the forge, so it runs by hand).
4. **The backlog** has the 18 lines, each with its issue; **the lesson** "A
   question whose summary its own table contradicts" is in
   [`guardrails.md`](../guardrails.md) §2.

**The rejected alternatives:** one task per package without the splits (rows 33
to 38 as two tasks, of about 1,800 and 2,400 lines); the package rules inside
the first task of `internal/session`, as the specification said; the step
`probe` inside the records of a task session; each row cut to one class (about
34 rows, each with its own gate); a task of its own for #148 now, before the
release exists.
