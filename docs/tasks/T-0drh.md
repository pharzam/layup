# T-0drh — the technical specification of LAYUP, phase 1

Issue: [#74](https://github.com/pharzam/layup/issues/74). Parent
[#42](https://github.com/pharzam/layup/issues/42), between the architecture
(`T-hbw8`, #72) and the implementation plan (`T-55n2`). Serves `F-0003#51` and
`REQ-012`. Base `648b37f`. Author: Claude Opus 5.5 on Claude Code. Evidence:
[`runs/T-0drh/`](../../runs/T-0drh/).

## The Operator's decisions

**O-114** (2026-10-01, [comment 5927429422 on #42](https://github.com/pharzam/layup/issues/42#issuecomment-5927429422),
"I agree.", to the proposal of comment 5927351177): a new child of #42, this
task, the technical specification of LAYUP for phase 1, in `docs/spec/`, one
section per requirement ID, before `T-55n2`; the same pull request corrects
`REQ-012`. The children table of #42 names it as row 4a.

**O-115** (2026-10-01, [comment 5928163366 on #74](https://github.com/pharzam/layup/issues/74#issuecomment-5928163366),
to the question of plan-review condition 2): option (a). In phase 1,
`layup setup` writes the setup record and the rule-path register as files and
prints one command; the Operator runs it to commit them to the target's
`layup-records` branch, as for the root commit. ADR-0014 gets a Consequences
line: `layup run` is the one writer from Start on. The §7.1 criteria of
`PRD-0001` stay as they are.

**O-116** (2026-10-01, [comment 5928828145 on #74](https://github.com/pharzam/layup/issues/74#issuecomment-5928828145),
to the decision request after review round 2): option (B). The cycle cap of
this task is 2. The author applies the fixes of round 2, and a fresh session
reviews only those commits.

## Plan and plan review

The plan (R12) and its review are comments on #74. A second reading before the
plan (Claude Fable 5.1, comment 5927721875) found three points, all accepted:
the count of the named records, the phase-1 boundary of `layup setup`, and the
records with no phase-1 writer. The plan review (Claude Fable 5.1 on Claude
Code, a fresh session, 12 min 16 s) gave `approve-with-conditions`: Budget
maximum 1,900 lines added plus removed over 20 files on the whole diff against
`648b37f`, close-out inside; Cycle cap 1. The author applied the six conditions
and the fourteen notes (comment "Plan: the author's answer to the plan review"
on #74).

## Decisions of the author

- **D1.** In phase 1, `layup setup` is a deterministic step runner: no forge
  call, no session. Start, the prose rows, the probes and the records-branch
  writes of `layup run` come later, each named in `docs/spec/setup.md`.
- **D2.** `docs/spec/` has one file per command or concern, one section per
  requirement ID.
- **D3.** A record's schema is a fenced block that a Go test can read; the test
  comes with the code of #29.
- **D4.** One exit-code rule for every engine check: 0 pass, 1 a finding or a
  fail, 2 a usage or input error; a check that did not run is never 0.
- **D5.** The stack catalog has a path in LAYUP's repository and is embedded in
  the binary.
- **`docs/setup/steps.tsv` does not change** (plan-review condition 6): it is
  the record of LAYUP's own setup.
