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
