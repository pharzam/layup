# T-y10b — row 29 of the plan, the package rules of M2b

Issue: [#154](https://github.com/pharzam/layup/issues/154), row 29 of [the tasks
of M2b](../plan/README.md#the-tasks-of-m2b), opened by `T-fdaq` (#152). Serves
`F-0003#52` through `NFR-005` and `NFR-007`. Base `8d2fa2e` (the merge of #171).
Author: Claude Opus 5.5 on Claude Code. Evidence:
[`runs/T-y10b/`](../../runs/T-y10b/).

## Plan and plan review

The plan (R12, comment 6080306205), its review (6080486220) and the author's
answer (6080486845) are comments on #154. The plan review (Claude Fable 5.1,
effort `xhigh`, on the Claude Code CLI with stream output, a fresh read-only
session in a clone at `8d2fa2e`, 7 min 31 s) gave `approve-with-conditions`:
Budget maximum 700 lines added plus removed over 17 files against `8d2fa2e`,
close-out inside; Cycle cap 1 (`TestPackageRules` is not a gate); no panel. Its
three conditions are applied: `netimport` is also a fixture of the engine rule;
the Task and Test cells of `NFR-005` and `NFR-007`; one reading of the register
form. The goal count is 2, final by O-187 of #152.

## What was done

1. **The checker** (`cmd/layup/rules_checker_test.go`): `readTable` reads the
   table of `M2b`, whose rows are packages of their own; `readEngine` reads the
   line of the engine checks; `checkRules` gives "engine checks: P depends on D"
   for a package of that line that depends on `internal/session` or a package of
   rule 5, reading no cell; a register row (a cell that is exactly one code span
   and "(the command of a harness register row)") may import `os/exec` and start
   a program by a call whose program is not a string literal, and nothing else;
   a cell with those words and more text cannot be read. The scan marks a use of
   `exec.Command` that is not a call, so a register row does not allow it.
2. **The fixtures:** `enginedep` (`internal/gate` imports `internal/session`;
   `internal/session` starts its register's command by a variable), and
   `netimport`, which now also breaks the engine rule in `internal/psb`.
3. **The specification:** `packages.md` holds the line of the engine checks, its
   rule, the form of a register row with its known limit (rule 3 rests on
   review in a register row), and "The test of the package rules" says what the
   checker reads and what each fixture gives; `session.md` and `gate.md` state the
   check that exists.
4. **The documents:** `traceability.md`; the Task and Test cells of `NFR-005`
   and `NFR-007` in `PRD-0001` §12, with a §13 line.

**Tests:** [`test-runs.md`](../../runs/T-y10b/test-runs.md): red before the
checker (it did not compile), red with the two rules turned off on a copy, red
for `TestPackageRules` before `packages.md`, green after.

**The rejected alternatives:** a special code span for a register program; the
engine rule inside "May import" (it is on each dependency); the engine breach in
a package of `netimport` off the engine line (the plan review's first choice;
the chosen form also shows that the rule reads a file behind a build
constraint).
