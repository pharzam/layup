# T-zck8 — the technical specification of M2a, the Start of `layup run`

Issue: [#123](https://github.com/pharzam/layup/issues/123). Milestone `M2a` of the
[implementation plan](../plan/README.md#milestones), which starts with its
specification task (O-114). Serves `F-0003#51` (`REQ-012`) for `NFR-001`,
`NFR-002`, `NFR-006` and `REQ-002`. Base `9e980b9` (the merge of #122). Author:
Claude Opus 5.5 on Claude Code. Evidence: [`runs/T-zck8/`](../../runs/T-zck8/).

## Plan and plan review

The plan (R12, comment 6034135134), its review and the author's answer
(6034605324) are comments on #123. The plan review (Claude Fable 5.1, effort
`xhigh`, on the Claude Code CLI with stream output, a fresh read-only session in
a clone at `9e980b9`, 10 min 47 s; comment 6034314980) gave
`approve-with-conditions`: Budget maximum 1,300 lines added plus removed over 18
files against `9e980b9`, close-out inside; Cycle cap 1; no panel. Its eight
conditions and its notes are applied.

**O-163** (2026-10-07, the Operator's comment 6034592068 on #123): (a), `M2a`
specifies an empty repository only; the adoption of a phase-1 target goes to
`M2d` (K41). **O-164** (the same comment): (b), a separate task slices the build
tasks of `M2a` after this task merges; this task's close-out opens its issue.

## What was done

1. **D1, [`docs/spec/run.md`](../spec/run.md):** the command (`--new`, the
   restart, `--host DIR`, the optional `--vision`, the table `run-steps`, the
   exit codes), the nine steps of Start, the restart, the lease and fencing, a
   human decision, copy before read, the Intake and control issues, the input
   states, the sections of `NFR-006`, `NFR-002` and `REQ-002`, the acceptance
   test of each part, and the pilot findings that `M2a` does not change.
2. **D2, [`docs/spec/forge.md`](../spec/forge.md):** the six capabilities and the
   milestone that first uses each, the App identity (the key file, the JWT, the
   installation token, the token to `git` in the environment of one call), the
   GitHub calls of `M2a`, the forge errors (FT1) and the test of the adapter.
3. **D3, [`docs/spec/records.md`](../spec/records.md#nfr-001--the-records-of-start):**
   the schema blocks `start`, `approvers`, `lease`, `copies`,
   `harness-register` and `forge-register` (and `run-steps` in `run.md`), in
   `notYetBuilt`; the Phase rule of the layout and the four records that K38
   moves.
4. **D4, [`docs/spec/packages.md`](../spec/packages.md#the-table-of-m2a):** the
   table of `M2a` with the column "Connects" (rule 5: only
   `internal/forge/github`), the changed rows of `internal/cli` and
   `internal/records`, and the calls `Fetch` and `Push` of `internal/git`.
5. **D5:** `docs/spec/README.md` (the version line, the files, the rule of
   `notYetBuilt`, the later phases), `docs/spec/setup.md` (the boundary rows of
   Start 3), `docs/plan/README.md` (K38, K40, K41; the hosts of the nine
   findings and of O-112), `PRD-0001` §12 and §13, and the glossary (`JWT`,
   installation token).

**Tests:** [`sections.sh`](../../runs/T-zck8/sections.sh) failed for each of its 18
headings at the base and passes on the head; the schema-block test failed for
the six new blocks of D3 until they were listed
([`test-runs.md`](../../runs/T-zck8/test-runs.md)).

**The rejected alternatives:** K40 by two flags or an environment variable; K41
(b), `--adopt` in `M2a`; O-164 (a), the task rows of `M2a` in this task (a second
goal, condition 1 of the plan review); a new column in the table of phase 1,
which would change the checker in this task.

**Known limits:** the harness credential stays open for `M2b` (K40). A run that
stops between the open of an issue and the record of its number leaves a second
issue on a restart ([`run.md`](../spec/run.md#input-states)). The token
reaches `git` in the environment of one call, which the same user of the host can
read. The GitHub calls are read from the architecture; the build task reads the
documentation of each at its date.

## Review rounds

The records, the Fixes replies and the decision at the cap are comments on #123.
Round 1 (Claude Fable 5.1; `832a3f2`, cycle 0): `material`, five findings (the two
forms of `run`, the exit code of a forge state, the issue rows and fencing, the
plan check after the root push, the backlog line), fixed. Round 2 (Fable;
`2ebf22a`, cycle 1): three material findings at the cap of 1 (rule 5 through the
adapter, the readers of the registers, the README of a Start target); **O-165**
(a): one more cycle, cap 2. Round 3 (Fable; `5706440`, cycle 2): `nothing
material in scope`, nine notes; notes 1 to 6 and 9 are applied in the close-out,
note 7 is declined (`NFR-005` is a requirement of `M2b`, whose specification task
names the check that replaces the phase-1 import rule), note 8 needs no change.
`review-record-lint` on the comments of #123 gives `OK  11 comments; 3 round(s);
cap 2`.

## Verdict

Delivered: the technical specification of `M2a` in `docs/spec/` (`run.md`,
`forge.md`, the records of Start, the table of `M2a`), with K38, K40 and K41
settled and the nine pilot findings moved to their milestones. The review ended
by decay at cycle 2 of a cap that the Operator raised once (O-165). The diff
against `9e980b9` is inside 1,300 lines over 18 files.

Next: `T-zwke` ([#124](https://github.com/pharzam/layup/issues/124)), the build
tasks of `M2a` sliced in the plan (O-164 (b)). Before the demo of `M2a`, the
Operator puts the App's private key on the host and installs the App on a test
target (O-112).

## Resource record

Recorded, not budgeted (ADR-0007). Times are UTC, 2026-10-07; tokens are the
`result` event of the Claude Code CLI (stream runs); `not reported` otherwise.

| Part | Expected tier | Model | Effort | Tokens | Elapsed |
| ---- | ------------- | ----- | ------ | ------ | ------- |
| The issue, the sources, the plan | reasoning | Claude Opus 5.5 | max | not reported | 08:15 to 08:32 |
| The plan review | reasoning | Claude Fable 5.1, Claude Code CLI | `xhigh` | not reported | 10 min 47 s, from 08:32 |
| The decision brief (English and Persian, for O-163 and O-164) | — | Claude Opus 5.5 | max | not reported | 08:45 to 09:00 |
| The answer; the work, test first; the fixes of rounds 1 and 2 | execution | Claude Opus 5.5, a reasoning-tier model on an execution part | max | not reported | 09:01 to 09:11, 09:24 to 09:26, 09:53 to 09:56 |
| Round 1 | reasoning | Claude Fable 5.1, Claude Code CLI | `xhigh` | 1,533,646 (USD 7.01) | 11 min 51 s, from 09:11 |
| Round 2 | reasoning | the same | `xhigh` | 1,083,759 (USD 6.90) | 12 min 27 s, from 09:26 |
| Round 3 | reasoning | the same | `xhigh` | 1,474,687 (USD 7.01) | 10 min 31 s, from 09:56 |
| The close-out, with the notes of round 3; the issue of `T-zwke` | execution | Claude Opus 5.5 | max | not reported | 10:07 to 10:15 |
