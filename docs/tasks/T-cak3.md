# T-cak3 — row 66a of the plan of M2e, the task register

Issue: [#357](https://github.com/pharzam/layup/issues/357), row 66a of [the tasks of
M2e](../plan/README.md#the-tasks-of-m2e), opened by `T-57vx` (#346). Serves `F-0003#45`
(`REQ-005`, `NFR-001`). Base `f12da6a`; `main` merged once, at `8e26eca` (`eada039`). Pull
request [#383](https://github.com/pharzam/layup/pull/383). Author: Claude Sonnet 5.5 on Claude
Code. Evidence: [`runs/T-cak3/`](../../runs/T-cak3/).

## Plan and plan review

The plan, its review and the author's answer are comments on #357 (6210010288, 6210141165,
6210153718). The plan review (Devin SWE-2 Max, a fresh session in a clone at `f12da6a`) gave
`approve`, no condition: Budget maximum 1,100 lines over 16 files, close-out and fixes inside;
Cycle cap 1; no panel; goal count 1 (O-214). The models follow the Operator's
mandate (Claude Sonnet 5.5, Devin SWE-2 Max; for the review, Antigravity Gemini was added), a tier
override. The budget is read against current `main` at landing (the Operator's
confirmation, 6209446809 on #354, recording the answer to 6209355485 on #355).

## What was done

`internal/run/register.go`: `Task`, `TasksSchema` (equal to the block `tasks`, moved to `built`),
`WriteRegister`, `ReadRegister`, `CheckRegister` (requirements never `—`, each predecessor a
task, no cycle, each refusal at its line and column) and `AddTask`; `registerTask` (the
permission guard, then one fenced records commit of `tasks.tsv` and, only when absent, the
default `transitions.tsv`; an existing table is checked and kept; a refusal writes nothing);
`loopPermissionsStep` and the field `granted` (Start and the restart still only record the
permissions). Registered in `packages.md`, `forge.md` (R10), `traceability.md` and `PRD-0001`
§12. **Not built, by the plan's readings:** the step `loop` and the restart (row 66k), the
check that a predecessor has merged (66g, 66k), a row's `pr` (66g), `parameters.tsv` (66b).

**Definition of Done:** 1 the demo, by 11 unit and 6 integration tests; 2 each part, by the
register, files, permission and schema tests and `TestPackageRules`; 3 the commands, below; 4 the
documents. **Red first:** both sets failed against a stub by assertion (unit 07:14:46 to
07:14:57Z, integration 07:15:00 to 07:15:50Z, exit 1). **Mutants:** 30, each alone with the code
building, all caught by assertion ([`test-runs.md`](../../runs/T-cak3/test-runs.md)).

## Review rounds

| Round | Head | Cycle | Cap | Material | Then |
| ----- | ---- | ----- | --- | -------- | ---- |
| 1 | `dae1023` | 0 | 1 | 0 | `nothing material in scope`, one note |

Round 1 (comment 6213845433): Google Gemini 4 Argon (High) on Antigravity CLI, a fresh session
briefed on the source, the issue body and the plan review only. It differs from the author (Claude
Sonnet 5.5) and from the GPT-6-family coordinator, who edited `test-runs.md` and `guardrails.md`
(`dae1023`). **The note** (`register.go:76-92, 126-169`): a `—` in `kind` is accepted, and one in
`issue` or `opened` is refused without its column. True, but `records.md:736` gives the owner of
`tasks` four checks and no rule for these columns, and `layup run` is the only writer: declined, no
fix, no issue; a change of the checks starts in `records.md`.

## Verdict

Delivered: the register, its checks, the first write and the permission step, with their tests. The
review ended by decay at cycle 0 of cap 1. The full integration level at `e16b7c9` (the
coordinator's `go test -count=1 -tags=integration -timeout 25m ./...`) exited 0 from
07:33:19.112569Z to 07:46:59.169359Z (13 min 40 s), 21 packages ok. No Go file of this branch has
changed since: the later commits are the merge of `main` (disjoint hunks in `PRD-0001` and
`traceability.md`, no conflict), the evidence and the counter correction. Build, vet, unit level, `git
diff --check` and the five checks exit 0 at the head. **Budget:** maximum 1,100 lines over 16 files;
the diff against `main` (`8e26eca`) is 1,099 over 15 with this close-out (1,018 over 12 before it).

## Resource record

Recorded, not budgeted (ADR-0007). UTC, 2026-10-16. Thinking is in output, not added again. The
Claude rows are per invocation, checked against the differences of the cumulative `modelUsage`
(one model, no child usage); cost is the list price from the difference of `costUSD`;
"unavailable" is a value no log gives.

| Part | Expected tier | Model | Effort | Input / output (thinking) / cache creation / cache read | Elapsed |
| ---- | ------------- | ----- | ------ | ------------------------------------------------------- | ------- |
| The plan | reasoning | Claude Sonnet 5.5 | high | 98 / 38,005 (22,980) / 176,936 / 5,591,107; USD 1.6471 | 359,780 ms, 06:56 to 07:02 |
| The plan review | reasoning | Devin SWE-2 Max | unavailable | unavailable | unavailable |
| Tests, code, mutants, documents, first push | execution | Claude Sonnet 5.5 | high | 82 / 52,988 (15,356) / 87,644 / 9,817,788; USD 1.8624 | 1,231,322 ms, 07:11 to 07:32 |
| The merge of `main`, the evidence | — | Claude Sonnet 5.5 | high | 22 / 7,563 (1,520) / 12,071 / 3,097,698; USD 0.4337 | 255,941 ms, 07:49 to 07:54 |
| The counter correction (`dae1023`) | — | GPT-6 family (Codex) | unavailable | unavailable | unavailable |
| Round 1 | reasoning | Gemini 4 Argon (High), `agy` | high (explicit CLI flag) | 405,037 / 46,727 (30,692) / not reported / 3,839,062; total 451,764 without the cache read | 779.9 s, 1 turn |
| This close-out | — | Claude Sonnet 5.5 | high | 38 / 21,397 (6,907) / 323,548 / 5,732,804 | 233,810 ms |
| Full integration | — | not applicable | not applicable | not applicable | 820,057 ms |
| **Total** | | | | reported subtotal 29,250,615 including cache; other portions unavailable | reported subtotal 3,680,855 ms; other portions unavailable |
