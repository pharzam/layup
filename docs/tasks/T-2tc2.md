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

The diff against `b48764f` went over the maximum of the plan review before the
close-out: 2,043 lines over 22 files after the notes of review round 1. The
author asked the Operator on #79 (comment "The budget, after the notes of review
round 1"): "**The new maximum to approve:** 2,150 lines over 25 files, close-out
inside. The alternative is a child issue for the growth. The Operator decides."

**O-128.** The Operator answered "#78 a; #79 a" in the author's Claude Code
session (2026-10-01), copied to
[#79](https://github.com/pharzam/layup/issues/79#issuecomment-5938474670): option
(a), the new maximum. The Budget maximum of this task is 2,150 lines added plus
removed over 25 files against `b48764f`, close-out inside. The Cycle cap stays 1.

## Verdict

Delivered: `internal/git`, the one caller of `git` (17 calls in a closed list, a
fixed environment with no input from the host's configuration or credentials,
and two error kinds); `TestPackageRules` in `cmd/layup`, which reads the table
of `docs/spec/packages.md` and holds each package to its rules, with a fixture
module that imports `net/http`; the cell form and D1 to D9 in
`docs/spec/packages.md`, which settle K7, K8, K31 and K39; the known limit L-A7
in `docs/architecture.md` §15 and at S02; the `PRD-0001` §12 Test cells of
`NFR-005` and `NFR-007`; the traceability row `green`; three glossary rows; one
lesson in `docs/guardrails.md` §2.

The plan review (Claude Fable 5.1) gave `approve-with-conditions`; the author
applied its six conditions and twelve notes, with the two departures above. A
fresh verification (Claude Opus 5.5) found 12 items, one material (the host's
credentials could still reach the clone of S02); the fixer fixed it and 9 of the
notes. Review round 1 (Claude Fable 5.1, `daa1ded`, cycle 0) gave `nothing
material in scope`, with seven notes: notes 1 to 6 are applied (note 1 also by a
correction comment on #79; note 3 by rows for the abbreviations APFS, GSS, NFC
and NFD); note 7 is O-128. At the head, `go build`, `go vet` and `go test`
(untagged, `integration` and `e2e`) pass, and so do all local checks, also
after the merge of row 1 (`673e289`): the merge resolved two conflicts, in the
package table and in the traceability table, and a probe showed that
`TestPackageRules` now holds `internal/tsv` too. The diff of this task against
`origin/main` after that merge is 2,104 lines over 24 files, inside the Budget
maximum of O-128.

Next: the PDR (`T-4wrw`, #99) waits for the Operator's approval; then the tasks
of phase 1 continue in the order of the plan, from row 3 (`T-2yw7`, #80).

## Resource record

Recorded, not budgeted (ADR-0007). Times are 2026-10-01, UTC. Token counts are
`not reported` where neither the harness nor `claude -p` in text mode gives them.
The plan and the plan review of rows 1 and 2 ran side by side; each row records
its own.

| Part | Expected tier | Model | Effort | Tokens | Elapsed |
| ---- | ------------- | ----- | ------ | ------ | ------- |
| The plan (after O-126) | reasoning | Claude Opus 5.5 | max | not reported | within 13:43 to 14:02, shared with the plan of row 1 and the fixes of `T-55n2` |
| The plan review | reasoning | Claude Fable 5.1 | not reported | not reported | 13 min, 14:25 to 14:37; a first run at 14:02 stopped at the usage limit, with no record |
| The answer to the plan review | reasoning | Claude Opus 5.5 | max | not reported | 14:37 to 14:39 |
| The implementation (a workflow agent) | execution | Claude Opus 5.5, a reasoning-tier model on an execution part | max | 810,254 (the workflow's count) | 156 min, 14:41 to 17:17; the first run stopped answering at 15:00, two restarts gave no answer, the fourth run stopped answering at 15:57, the fifth run went 16:13 to 17:17 |
| The verification (a fresh workflow agent) | reasoning | Claude Opus 5.5 | max | 285,112 (the workflow's count) | 25 min, 17:17 to 17:43 |
| The fixes of the verification (a workflow agent) | execution | Claude Opus 5.5 | max | 314,776 (the workflow's count) | 32 min, 17:43 to 18:14 |
| The task record, the traceability row, the §12 Test cells and the freeze | execution | Claude Opus 5.5 | max | not reported | 18:14 to 18:18 |
| Review round 1 | reasoning | Claude Fable 5.1 | not reported | not reported | 15 min, 18:19 to 18:34 |
| The notes of round 1, the correction on #79 and the budget question | execution | Claude Opus 5.5 | max | not reported | 18:34 to 18:37 |
| Close-out, with the merge of `origin/main` after row 1 | — | Claude Opus 5.5 | max | not reported | 19:00 to 19:12 |
