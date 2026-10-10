# T-nxe4 — row 38 of the plan, the step probe

Issue: [#168](https://github.com/pharzam/layup/issues/168), row 38 of [the tasks
of M2b](../plan/README.md#the-tasks-of-m2b), opened by `T-fdaq` (#152). Serves
`F-0003#52` through `REQ-013` and `NFR-001`. Base `6af022c` (the merge of #196,
row 37b); the budget is read against it. Author: Claude Opus 5.5 on Claude
Code. Evidence: [`runs/T-nxe4/`](../../runs/T-nxe4/).

## Plan and plan review

The plan (R12), its review and the author's answer are comments on #168. The
plan review (Claude Fable 5.1, effort `xhigh`, a fresh read-only session at row
37b's head `5586d66`, 9 min 31 s) gave `approve-with-conditions`: Budget
maximum 1,600 lines added plus removed over 24 files, close-out inside; Cycle
cap 1; no panel. Its three conditions are applied: the probe before a task
session (`Admit`, with no sweep inside a task session); a refused probe start
counts as probed and failed and posts no comment, written in `session.md`; the
e2e worlds register only scripted harnesses. The goal count is 2, final by
O-187 of #152. **O-191** (a): the pushed commit `e4cf2ce` held a test
constant that the rule `generic-api-key` of `gitleaks` matched, a false
positive; `.gitleaksignore` gains its fingerprint, a change to a gate's input,
so the cycle cap is 2.

## What was done

1. **The probe:** `Sessions.Probe`, a session of the role `probe` on the frame of
   `TaskSession`: its own task ID, the fixed prompt with the result path and a
   token, no attempt check, a start row with no event, the end `endProbe` (the
   pass rules of `probeReason`, one records commit of its row of
   `harnesses.tsv` and its telemetry row, then its comment), and a refused start
   written as its row. A version whose last probe passed ends the session with
   no record.
2. **The step `probe`** of the restart, after `lease`: the copy of the routing
   register when the tables differ, the skip, the probes, the three counts;
   `probe` in the enum of `run-steps`; `Config` takes the models, routing and
   price rows, which `internal/cli` hands over (`run.ReadPrices`).
3. **The probe before a task session:** `Sessions.Admit`.
4. **The tests and the documents:** the unit, integration and e2e tests; `session.md`, `run.md`, `packages.md`,
   `traceability.md`, the Test cells of `REQ-013` and `NFR-001` and a §13 line.

**Tests:** [`test-runs.md`](../../runs/T-nxe4/test-runs.md): red before the
code (the unit tests); a mutation of each of the eighteen rules, each caught;
green after.

**A deviation:** the integration and e2e tests were written after the code;
their red is the mutations.

**The rejected alternatives:** the version command run outside a session to
decide the skip; one records commit for the whole step.

## Review rounds

The records are comments on #168. Round 1 (Claude Fable 5.1, effort `xhigh`, a
fresh read-only session in a clone at `918a2f4`, cycle 0): `material`, one
finding and six notes. Finding 1: the skip and a probe's refused start were
unit acceptance clauses that this row tests at integration; they moved to a
row of their own at integration. The fix also reads the `records` column per
probe (note 3), qualifies the sentence on the counts (note 2), names the cost
of the early end (note 5) and asserts the two sentences with no test (note 7);
note 4 is declined with its reason (a routing table is read by its positions,
in any order), and note 6 is the deviation above. Round 2 (the same reviewer
type, a fresh session at `b61f585`, cycle 1, the cap): `nothing material in
scope`, four notes. Applied in the close-out: note 1 (row 38 of the plan names
the moved clauses); note 2 (the unit row reads "the result of a fake
harness"), with the new row renamed "The skip, a refused probe and the early
end" apart from the e2e row "The step `probe`"; note 3 (a row of Input states
for a `host:prices.tsv` that its reader refuses). Note 4 (a `TaskSpec` with no
`Admit` panics) is declined: the field was called the same way at the base,
each caller sets it, and a code change after the last round is read by no
round.

## Verdict

Delivered: the step `probe` in `internal/run`: the probe of each harness and
its pass rules, the skip, the three counts, the copy of the routing register,
and the probe before a task session (`Admit`). The review ended by decay at
cycle 1 of a cap that O-191 raised to 2. The diff against `6af022c`, the
branch's base, is inside 1,600 lines over 24 files. Next: row 39a.

## Resource record

Recorded, not budgeted (ADR-0007). UTC; reviewers' tokens are `modelUsage`;
the author's are not reported.

| Part | Expected tier | Model | Effort | Tokens | Elapsed |
| ---- | ------------- | ----- | ------ | ------ | ------- |
| The plan | reasoning | Claude Opus 5.5 | not reported | not reported | to 2026-10-10 03:51 |
| The plan review | reasoning | Claude Fable 5.1 | `xhigh` | 1,423,679 (USD 6.24) | 9 min 31 s |
| The code, test first | execution | Claude Opus 5.5 | not reported | not reported | 04:05 to 04:24 |
| O-191 asked and applied | execution | Claude Opus 5.5 | not reported | not reported | 04:25; 06:39 to 06:40 |
| Round 1 | reasoning | Claude Fable 5.1 | `xhigh` | 1,919,954 (USD 6.66) | 12 min 24 s |
| The fix of round 1 | execution | Claude Opus 5.5 | not reported | not reported | 06:53 to 06:59 |
| Round 2 | reasoning | Claude Fable 5.1 | `xhigh` | 1,230,093 (USD 6.61) | 12 min 21 s |
| The close-out, with notes 1 to 3 | execution | Claude Opus 5.5 | not reported | not reported | 07:11 to 07:13 |
