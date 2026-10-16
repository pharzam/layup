# T-cr8j — row 64 of the plan of M2e, the forge calls

Issue: [#354](https://github.com/pharzam/layup/issues/354), row 64 of [the tasks of
M2e](../plan/README.md#the-tasks-of-m2e), opened by `T-57vx` (#346). Serves `F-0003#58`
(`REQ-007`, `NFR-001`, `NFR-004`). Base `9e6393a`; `main` was merged twice, at `2aafa64`
and at `d24fae3`. Author: Claude Sonnet 5.5 on Claude Code. Evidence:
[`runs/T-cr8j/`](../../runs/T-cr8j/).

## Plan and plan review

The plan, its review and the author's answer are comments on #354. The plan review
(GPT-6.1 Sol on Codex CLI, `xhigh`; allowed by the Operator on 2026-10-13) gave
`approve-with-conditions`: Budget maximum 1,500 lines over 18 files, close-out inside;
Cycle cap 1; no panel. Its four conditions are applied.

## What was done

`OpenDraft`, `MarkReady`, `SetStatus`, `Statuses`, `PullRequest` and `Merge` in
`internal/forge` and its GitHub adapter, `CheckLoopPermissions`, and the stand-in
`forge.Fake`, tests first (26 red for the right reason, 29 mutants caught); `forge.md`,
`loop.md`, `packages.md` and `api-sources.md` follow. **Off the path:** #373, #380, #336.

## Review rounds

Both rounds: Claude Opus 5.5, `xhigh`, a fresh `claude -p` session; comments on #354.

| Round | Head | Cycle | Cap | Material | Then |
| ----- | ---- | ----- | --- | -------- | ---- |
| 1 | `8c2bd66` | 0 | 1 | 3 | fixes in `9ae6562` |
| 2 | `9ae6562` | 1 | 1 | 1 | `not mergeable, findings recorded`; O-216: fixes settled, no round |

The fixes of round 2 (`9916120`) and the two merges of `main` (O-219) are in the comments
"Fixes of round 2" and "The second merge of main". The budget is read against current
`main` (comments 6209446809 on #354 and 6209446202 on #355).

## Verdict

Delivered, as code and test result: the six calls, the permission check and the stand-in
forge, tested against the httptest forge; the full integration run at `7896b20` exited 0
([`test-runs.md`](../../runs/T-cr8j/test-runs.md)). Round 2's material finding is
settled as fixed by O-216. **Landing:** the diff against `main` (`d24fae3`), close-out
included, is 1,557 lines over 18 files, over the plan review's 1,500 over 18; the Operator
approved exactly that figure, once (comment 6209661198 on #354); it is no ceiling.

## Resource record

Recorded, not budgeted (ADR-0007). UTC; "not reported" is a value the harness did not give.

| Part | Expected tier | Model | Effort | Tokens | Elapsed |
| ---- | ------------- | ----- | ------ | ------ | ------- |
| The plan | reasoning | Claude Sonnet 5.5 | not reported | not reported | start not reported, to 18:33 on 2026-10-15 |
| The plan review | reasoning | GPT-6.1 Sol | `xhigh` | 177,104 | 5 min 10 s |
| The implementation (T and D) | execution | Claude Sonnet 5.5 | not reported | not reported | 18:36 to 18:58 on 2026-10-15 |
| Round 1 | reasoning | Claude Opus 5.5 | `xhigh` | not reported | not reported, ends 19:21 |
| The fixes of round 1 | execution | Claude Sonnet 5.5 | not reported | not reported | 19:25 to 19:38 |
| Round 2 | reasoning | Claude Opus 5.5 | `xhigh` | not reported | not reported, ends 19:53 |
| The fixes of round 2 | execution | Claude Sonnet 5.5 | not reported | not reported | 19:55 |
| The recovery: merges, checks, full integration, pushes | execution | Claude Sonnet 5.5 | `high` | input 122, output 32,926 (thinking 10,423 of it), cache creation 100,959, cache read 5,022,223 | 2,086,060 ms (34 min 46 s), about 05:53 to 06:28 on 2026-10-16 |
| This close-out | execution | Claude Sonnet 5.5 | `high` | not reported | 06:29 to 06:42 on 2026-10-16, the wait for the approval inside |
