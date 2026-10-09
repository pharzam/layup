# T-bpxg — row 34b of the plan, the end of a session

Issue: [#162](https://github.com/pharzam/layup/issues/162), row 34b of [the tasks
of M2b](../plan/README.md#the-tasks-of-m2b), opened by `T-fdaq` (#152). Serves
`F-0003#52` and `F-0003#50` through `REQ-013`, `REQ-005` and `REQ-011`. Base
`53f1bef` (the merge of #188, row 34a); the budget is read against it. Author:
Claude Opus 5.5 on Claude Code. Evidence: [`runs/T-bpxg/`](../../runs/T-bpxg/).

## Plan and plan review

The plan (R12), its review and the author's answer are comments on #162. The
plan review (Claude Fable 5.1, effort `xhigh`, a fresh read-only session at
row 34a's head `179646f`, 8 min 57 s) gave `approve-with-conditions`: Budget
maximum 1,000 lines added plus removed over 18 files, close-out inside; Cycle
cap 1; no panel. Its four conditions are applied: the token statuses are
total, with a reading of `partial` and of a report with no model; `token` at
most once, in the rules of `records.md`, with cases of zero rows and a wrong
header; the order of the classes, the refusal of a result file that is not a
regular file, and the sorted models, each **decided here** in `session.md`;
the recorded run names its models and version, and its stream was read before
anything of it was kept. The goal count is 2, final by O-187 of #152.

## What was done

1. **D1:** `Class`: the class of the end, in the order `start`; `wall` or
   `output`; `crash`; `no-result`; `done`. Another error of `Process` is no
   class.
2. **D2:** `ResultOf` (no link followed, a regular file of at most 1 MiB, by its
   block's reader), `ProbeResultSchema` and `ReadProbeResult`, the block
   `probe-result` moved to `built`.
3. **D3:** `Usage` (the fields of `ledger.Usage`, which a test converts),
   `UsageOf`: `claude-result` and `none`, on two recorded `result` events of
   Claude Code 2.1.295.
4. **D4:** `session.md`, `records.md`, `packages.md`, `traceability.md`, the
   Test cells of `REQ-013`, `REQ-005` and `REQ-011` and a §13 line.

**Tests:** [`test-runs.md`](../../runs/T-bpxg/test-runs.md): the recorded
events and their sums by hand; red before the functions; a mutation of each of
the twenty-four rules, each caught; green after.

**The rejected alternatives:** a class from the exit alone; a result file read
through a link; a float for the cost.
