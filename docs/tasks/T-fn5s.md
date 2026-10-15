# T-fn5s — row 60 of the plan, the kinds of the task loop's records

Issue: [#347](https://github.com/pharzam/layup/issues/347), row 60 of [the tasks of
M2e](../plan/README.md#the-tasks-of-m2e), opened by `T-57vx` (#346). Serves
`F-0003#45` (Role Handoffs; `F-0003#59` is the related fact). Base `48f07ca` (the
merge of #371). Author: Claude Sonnet 5.5 on Claude Code. Evidence:
[`runs/T-fn5s/`](../../runs/T-fn5s/).

## Plan and plan review

The plan (R12), its review and the author's answer are comments on #347. The plan
review (GPT-6.1 Sol on Codex CLI, effort `xhigh`, a fresh sandboxed session in a
clone at `48f07ca`, 4 min 18 s; the harness is outside the routing rule of ADR-0026
and was allowed by the Operator on 2026-10-13) gave `approve-with-conditions`:
Budget maximum 400 lines added plus removed over 10 files, close-out inside; Cycle
cap 1; no panel. Its four conditions are applied: the encoding of a question's label,
`needs_human` word and records, and of a lesson's records, decided in `records.md`;
the forms of a handoff named in the tests; the nine commands of criterion 3; the
In-Scope fact named. The goal count is 1 (O-214).

## What was done

1. **T:** the cases in `internal/records/session_test.go`, seen red for the right
   reason before the code ([`test-runs.md`](../../runs/T-fn5s/test-runs.md)).
2. **D1:** `internal/records/session.go`: the enums of the two Go schemas, `CheckEvent`
   (`handoff`, `merged`), `CheckResultRow` (the six kinds) and `CheckResult` (the numbers
   of each kind).
3. **D2:** `records.md`: the two blocks and the rules the specification did not give
   (decided here, R10); one sentence of `loop.md` (`invalid pair`); the trace entry and
   the Test cells of `PRD-0001` §12.

**Revealed, off the path, as an issue:** #377 (the kind of `result` that holds an owner
role's answer, which rows 61 and 66d meet).

## Review rounds

The records are comments on #347. Both rounds were Claude Opus 5.5, effort `xhigh`, a
fresh read-only `claude -p` session in a clone, a model that wrote none of the change.

| Round | Head | Cycle | Cap | Material | Then |
| ----- | ---- | ----- | --- | -------- | ---- |
| 1 | `0c2a8ff` | 0 | 1 | 3 | fixes in `1438ced`; round 2 (cycle 1) is the last the cap allows |
| 2 | `1438ced` | 1 | 1 | 0 | `nothing material in scope` |

The verdict of round 1 was edited to `material`, as RR8 asks. Round 2 left six notes
that change nothing at this head: no case accepts `invalid pair` or an `invalid` detail
with a comment suffix (the code is right); the names of the items will have two homes
after row 61; the words about `pair` are loose; `merged` has a second writer in row
66k; the record's line count is 236, not 232; and the finding 1 above, as #377.

## Verdict

Delivered: the Go values of `result` and `events` hold the kinds of the task loop, with
their rules, and each reader refuses a row that breaks them. The review ended by decay
at cycle 1 of cap 1. The diff against `48f07ca` is inside 400 lines over 10 files. Next:
rows 61, 62a and the others that follow row 60.

## Resource record

Recorded, not budgeted (ADR-0007). UTC, 2026-10-15. Reviewers' tokens are as the
harness reports them; the author's are not reported.

| Part | Expected tier | Model | Effort | Tokens | Elapsed |
| ---- | ------------- | ----- | ------ | ------ | ------- |
| The plan | reasoning | Claude Sonnet 5.5 | not reported | not reported | 17:20 (the claim) to 17:25 (the plan's comment) |
| The plan review | reasoning | GPT-6.1 Sol | `xhigh` | 142,337 | 4 min 18 s |
| T and D1, D2 | execution | Claude Sonnet 5.5 | not reported | not reported | 17:26 to 17:50 (the full integration run, 19 min 15 s, is in it) |
| Round 1 | reasoning | Claude Opus 5.5 | `xhigh` | 4,736,643 (USD 3.10) | 12 min 22 s |
| The fix of round 1 | execution | Claude Sonnet 5.5 | not reported | not reported | 18:06 to 18:12 |
| Round 2 | reasoning | Claude Opus 5.5 | `xhigh` | 3,796,573 (USD 2.85) | 10 min 33 s |
| The close-out | execution | Claude Sonnet 5.5 | not reported | not reported | 18:25 to 18:30 |
