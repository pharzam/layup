# T-2tc2 — `internal/git` and the test of the package rules

Issue: [#79](https://github.com/pharzam/layup/issues/79), row 2 of the
[implementation plan](../plan/README.md); child of `layup gate` (#33). Started
before the PDR by the Operator's decision O-126. Serves `F-0003#44`. Base
`b48764f` (the branch was rebased onto the merge of the plan, #98). Author:
Claude Opus 5.5 (the author's session and the workflow agents). Evidence:
[`runs/T-2tc2/`](../../runs/T-2tc2/).

## Plan and plan review

The plan (R12), its review and the author's answer are comments on #79. The plan
review (Claude Fable 5.1, a fresh session) gave `approve-with-conditions`: Budget
maximum 1,600 lines added plus removed over 20 files against the base, close-out
inside; Cycle cap 1. The author applied the six conditions and the twelve notes.

## What was done

1. **Test first.** The unit tests of each call with a stub runner, the
   integration tests with real `git` under a hostile configuration, and the
   checker of the package rules, each run red before its code
   ([`red.md`](../../runs/T-2tc2/red.md)).
2. **`internal/git`:** the one caller of `git`, with a closed list of calls, a
   fixed environment, and two error kinds; **the test of the package rules** in
   `cmd/layup`, which reads the table of `docs/spec/packages.md`.
3. **`docs/spec/packages.md`:** the cell form and D1 to D9 as values decided
   here; the known limit L-A7 (a private baseline) in `docs/architecture.md`
   §15 and at S02; a lesson in `docs/guardrails.md` §2.
4. **A fresh verification** (Claude Opus 5.5) found 12 items. One was material:
   the host's credentials could still reach the clone of S02. The fixer fixed
   it and 9 of the notes.

**Two departures from the answer to the plan review**, posted on #79: D3 sets
`HOME=/dev/null` and `GIT_ALLOW_PROTOCOL=file:git:http:https` (with
`http.emptyAuth=false` and `core.precomposeUnicode=false`) in place of the host's
home and an ssh option (the status comment on #79 said "only `file` and `https`";
the value of the code and the specification is the one here, review round 1,
note 1);
D1 drops the call `branch`, which no step names.

**The requirements (condition 5):** `TestPackageRules` is the test of `NFR-005`
(no network package) and `NFR-007` (the standard library only; `git` only as a
program) in `PRD-0001` §12. This task adds no test of the criterion of
`NFR-001`; its tests come with rows 9, 13, 15, 16 and 20 of the plan.

**The rejected alternatives (note 7):** a Go copy of the table of the package
rules (it can drift from its one home); a new package for the test (it needs a
row of its own); a default identity for commits (the identity of the setup's
commits is the Operator's decision, for row 8); direct imports only for rule 5
(a dependency through another package would pass).

**The verification** is recorded in [`verification.md`](../../runs/T-2tc2/verification.md):
twelve findings; the fixer did not apply finding 8 (the close-out edits, which
the author made before the freeze) and finding 9 (the budget, the Operator's).

## The budget

The diff against `b48764f` is over the maximum of the plan review before the
close-out (reported on #79). The Operator's answer goes here.
