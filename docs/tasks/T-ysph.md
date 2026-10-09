# T-ysph — row 30a of the plan, the host registers of M2b

Issue: [#155](https://github.com/pharzam/layup/issues/155), row 30a of [the tasks
of M2b](../plan/README.md#the-tasks-of-m2b), opened by `T-fdaq` (#152). Serves
`F-0003#52` through `REQ-013` and `NFR-001`. Base `ad8774d` (the merge of #173;
the plan was reviewed at `8d2fa2e`). Author: Claude Opus 5.5 on Claude Code.
Evidence: [`runs/T-ysph/`](../../runs/T-ysph/).

## Plan and plan review

The plan (R12, comment 6080340002), its review (6080580136) and the author's
answer (6080580748) are comments on #155. The plan review (Claude Fable 5.1,
effort `xhigh`, on the Claude Code CLI with stream output, a fresh read-only
session in a clone at `8d2fa2e`, 7 min 45 s) gave `approve-with-conditions`:
Budget maximum 1,100 lines added plus removed over 20 files against `8d2fa2e`,
close-out inside; Cycle cap 1; no panel. Its four conditions are applied: the
columns of `harness-register` that never hold `—`, one test per column; one
reading of a `PATH` with `..` and the form of `NAME`; a home in `records.md` for
each rule beyond the types; `packages.md` and the evidence file. The goal count
is 2, final by O-187 of #152.

## What was done

1. **D1:** the ten columns of `M2b` in the block `harness-register` and in
   `HarnessRegisterSchema`, after `wall`; `ReadHarnesses` checks their rules
   (the columns that never hold `—`; `{cap}` exactly when `cap`; one space
   between words; `credential`, `credential_to`, `policy`, `vars`); the three
   fixtures that write the register.
2. **D2:** `ModelsSchema` and `ReadModels`, `RoutingRegisterSchema` and
   `ReadRoutingRegister` (the positions of each role and tier 1 to k, in any
   order, as `routing`), `CheckRegisters` across the three files; `models` and
   `routing-register` moved to `built`.
3. **D3:** `CheckCredential`: an absolute path, a regular file of mode 0600,
   owned by the user of the run.
4. **D4:** `layup run` reads `models.tsv` and `routing.tsv`, checks the three
   registers across their files and each credential, and exits 2 naming the
   file. `Config` keeps no field for them yet: row 30b or 38 adds it.
5. **D5:** `records.md` (the block, and the rules of the three host registers
   of `M2b`), `session.md` (the form of `NAME`, the prohibitions of
   `var:NAME`, one reading of a `PATH` with `..`), `run.md` (the four
   registers), `packages.md` (the Job cells of `internal/cli` and
   `internal/route`), `traceability.md`, the Test cell of `REQ-013` and a §13
   line.

**Tests:** [`test-runs.md`](../../runs/T-ysph/test-runs.md): red before the
readers (no compile); red for `layup run` before it read the registers of
`M2b`; green after.

**Decided here** (note 6 of the plan review): the `NAME` of `var:NAME` takes
the prohibitions of a `vars` name, as a `var:PATH` would take the place of the
host's `PATH`.

**The rejected alternatives:** one check of a file's mode and owner shared by
`internal/forge` and `internal/route` (a new package or an import that neither
row allows); the credential check at the start of a session (`session.md` gives
it to the input states of `layup run`); a closed list of roles.

## Review rounds

The records and the Fixes reply are comments on #155. Round 1 (Claude Fable
5.1; `089a357`, cycle 0): `material`, two findings (a position `0` passed the
routing reader when one place of 1 to k was missing; the red record had no run
for each rule), fixed in `298e947` with notes 3 to 7; the same check of
`CheckRouting` in `internal/records` is revealed and not reachable yet, issue
#178. Round 2 (Fable; `298e947`, cycle 1): `nothing material in scope`, four
notes. Applied in the close-out: note 1 (the two mutation groups that the
record cut to four lines, in full; the lone position `0` asserts its line);
note 2 (`records.md` says a `PATH` is not empty); note 3 (the lesson of a
fixture that breaks a second rule, in `guardrails.md` §2). Declined: note 4
(the input-state rows of `models.tsv` and `routing.tsv` are in `session.md`,
to which `run.md` now points, so `run.md` keeps its one row).

## Verdict

Delivered: the ten columns of `M2b` in the harness register, the readers of
`models.tsv` and `routing.tsv`, the checks across the three registers and of a
credential file, and `layup run` reading the four registers. The review ended
by decay at cycle 1, the cap. The diff against `ad8774d`, the branch's base, is
inside 1,100 lines over 20 files. Next: row 30b (admission and the routing
order) can start; the Operator's host needs the ten columns, `models.tsv` and
`routing.tsv` before `layup run` runs on it again.

## Resource record

Recorded, not budgeted (ADR-0007). Times are UTC, 2026-10-09; the reviewers'
tokens are `modelUsage` of the CLI's `result` event; the author's are not
reported.

| Part | Expected tier | Model | Effort | Tokens | Elapsed |
| ---- | ------------- | ----- | ------ | ------ | ------- |
| The plan | reasoning | Claude Opus 5.5 | not reported | not reported | 11:52 to 11:54 |
| The plan review | reasoning | Claude Fable 5.1, Claude Code CLI | `xhigh` | 876,120 (USD 6.54) | 7 min 45 s, from 11:54 |
| The answer; the tests, red; the code; the documents | execution | Claude Opus 5.5, a reasoning-tier model on an execution part | not reported | not reported | 12:24 to 12:34 |
| Round 1 | reasoning | Claude Fable 5.1, Claude Code CLI | `xhigh` | 1,000,652 (USD 5.51) | 11 min 43 s, from 12:35 |
| The fix of round 1 | execution | Claude Opus 5.5 | not reported | not reported | 12:50 to 12:55 |
| Round 2 | reasoning | Claude Fable 5.1, Claude Code CLI | `xhigh` | 1,104,484 (USD 4.71) | 10 min 28 s, from 12:55 |
| The close-out, with notes 1 to 3 | execution | Claude Opus 5.5 | not reported | not reported | 13:08 to 13:12 |
