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

**O-117** (2026-10-01, the Operator's answer "B" in the author's Claude Code session, copied to #74): option (B) of the request after review round 3. The cycle cap of this task is 3; a fresh session reviews only the fix of round 3.

**O-118** (2026-10-01, the Operator's answer "Keep" in the author's session, copied to [#74](https://github.com/pharzam/layup/issues/74#issuecomment-5929172285)): round-1 note 6. The `REQ-002` criterion of `PRD-0001` §7.1 keeps the words "the steps follow the step table of `docs/spec/setup.md`, which derives from `docs/setup/steps.tsv`"; O-115 does not bar it.

**O-119** (2026-10-01, the Operator's answer "go ahead with the close-out" in the author's session, copied to #74): option (A) of the request after review round 4. The author commits the fix of round 4 and closes out with no fifth round.

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

## What was done

1. **The checklist first** ([`checklist.md`](../../runs/T-0drh/checklist.md)):
   13 requirement sections, the package table, 56 record kinds of the whole
   architecture, the phase-1 boundaries, the conventions; every row open, then
   closed against `docs/spec/`.
2. **`docs/spec/`**: `README.md` (the conventions, the exit codes, the schema
   block and its types), `packages.md` (`NFR-007`), `records.md` (the layout of
   the records branch for phases 1 to 4; `NFR-001`, `NFR-002`; `REQ-009` and
   `REQ-011` as schemas only), `psb-check.md` (`REQ-001`), `setup.md`
   (`REQ-002`, `NFR-003`, `NFR-006`; the step table S01 to S15; the stack
   catalog), `gate.md` (`REQ-004`, `REQ-007`, `NFR-004`, `NFR-005`). Fourteen
   schema blocks that a Go test can read.
3. **Registered:** `PRD-0001` (`REQ-012` names `docs/spec/`; the `REQ-002`
   criterion names the step table; §9 phase 1 states the boundary; §12; §13),
   ADR-0014 (a Consequences line, O-115), the architecture (one link; the §3
   sentence on the first records commit), the glossary (Phase; `TSV`, `UTF-8`,
   `RFC 3339`, `ISO 4217`), README, onboarding, and a lesson in
   `guardrails.md` §2.
4. **Not changed:** `docs/setup/steps.tsv` (plan-review condition 6), any Go
   code, any check script.

## Verdict

Delivered: LAYUP's technical specification for phase 1 in `docs/spec/`, one
section per phase-1 requirement, and its registration. Four review rounds
(Claude Fable 5.1, fresh read-only sessions): round 1 on `d0634a4`, 5 material
findings and 21 notes; round 2 on `32304b9` (cycle 1), 1 material and 8 notes;
round 3 on `f83320e` (cycle 2, O-116), 1 material and 4 notes; round 4 on
`db960ef` (cycle 3, O-117), 1 material and 2 notes. Each record is in
`runs/T-0drh/`. The last round ended `not mergeable, findings recorded`; its
finding and both notes are fixed in `d3458f7`, which no round reviewed, by the
Operator's decision O-119. Every finding of every round is fixed; none became an
issue. **Budget:** about 1,600 lines added plus removed, inside the maximum of 1,900;
21 files, one over the maximum of 20. The extra files are the four review
records under `runs/T-0drh/`, which the two rounds beyond the plan (O-116,
O-117) added; they are kept as separate records, not joined to meet the count. All local checks pass on the landing head
([`test-runs.md`](../../runs/T-0drh/test-runs.md)). Next: `T-55n2`, the
implementation plan (#42), then #29.

## Resource record

Recorded, not budgeted (ADR-0007). Times are 2026-10-01, UTC. Token counts are
`not reported`: neither the author's harness nor `claude -p` in text mode
gives them.

| Part | Expected tier | Model | Effort | Tokens | Elapsed |
| ---- | ------------- | ----- | ------ | ------ | ------- |
| The plan, its answer, the questions | reasoning | Claude Opus 5.5 | auto | not reported | 08:15 to 08:35, and the answers |
| The plan review | reasoning | Claude Fable 5.1 | not reported | not reported | 12 min 16 s |
| Writing the specification | execution | Claude Opus 5.5, a reasoning-tier model on an execution part | auto | not reported | 08:50 to 09:03 |
| Review rounds 1 to 4 | reasoning | Claude Fable 5.1 | not reported | not reported | 13 min 38 s; about 11 min; about 11 min; about 6 min |
| Fixes of rounds 1 to 4 | execution | Claude Opus 5.5 | auto | not reported | about 25 min in all |
| Isolate, guardrails, docs, close-out | `—` | Claude Opus 5.5 | auto | not reported | 10:10 to 10:25 |
