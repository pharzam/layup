# T-h17s — row 61 of the plan of M2e, the handoff

Issue: [#348](https://github.com/pharzam/layup/issues/348), row 61 of [the tasks of
M2e](../plan/README.md#the-tasks-of-m2e), opened by `T-57vx` (#346). Serves
`F-0003#45` and `F-0003#59` (`REQ-005`, `NFR-004`). Base `9e6393a` (the merge of #375,
row 60). Author: Claude Sonnet 5.5 on Claude Code. Evidence:
[`runs/T-h17s/`](../../runs/T-h17s/).

## Plan and plan review

The plan (R12), its review and the author's answer are comments on #348. The plan
review (GPT-6.1 Sol on Codex CLI, effort `xhigh`, a fresh sandboxed session in a clone
at `9e6393a`, 6 min 29 s; the harness is outside the routing rule of ADR-0026 and was
allowed by the Operator on 2026-10-13) gave `approve-with-conditions`: Budget maximum
1,100 lines over 16 files, close-out inside; Cycle cap 1; no panel. Its five conditions
are applied: the schema and every named artifact guard each handoff whatever the row
lists; the tests of the three comment kinds and the table refusals; `Items()` against the
event details (`pair` is no table item); the commands named in the record; `F-0003#45`.

## What was done

1. **`internal/handoff`** (new, importing `internal/tsv` only): the thirteen default rows,
   `ReadTable`, `CheckTable`, `WriteTable`, the block `transitions` (moved to `built`),
   and `Check`, whose verdict is the detail of the event `handoff` (`valid`,
   `invalid <item>`, `invalid pair`, with ` comment <ID>-<seen>`). No field that an agent
   writes decides: the hashes named are compared with the hashes found.
2. **Registration:** the heading of the table of `M2e` in `TestPackageRules` (`cmd/layup`),
   its fixtures, `packages.md`, `traceability.md` and the Test cells of `PRD-0001` §12.

**Tests:** [`test-runs.md`](../../runs/T-h17s/test-runs.md): red against a stub, 37
mutants caught, the full integration level green.

## Review rounds

Both rounds: Claude Opus 5.5, effort `xhigh`, a fresh read-only session in a clone.
Round 1 (`1a5aeff`, cycle 0): `material` (a code-written handoff skipped a `schema` item
that its row requires), six notes, all fixed in `7b554e2`. Round 2 (`7b554e2`, cycle 1,
the cap): `nothing material in scope`, seven notes; the notes that are documents and tests
are fixed in `8c460c0` with no round; a question for row 66b is posted on #358.

## Verdict

Delivered: `internal/handoff`, with its table, its block and its check. The review ended
by decay at cycle 1 of cap 1. The diff against `9e6393a` is inside 1,100 lines over 16
files. Next: rows 66a, 66d and 66e use it.

## Resource record

Recorded, not budgeted (ADR-0007). UTC, 2026-10-15; the author's tokens are not reported.

| Part | Expected tier | Model | Effort | Tokens | Elapsed |
| ---- | ------------- | ----- | ------ | ------ | ------- |
| The plan | reasoning | Claude Sonnet 5.5 | not reported | not reported | 18:26 to 18:33 |
| The plan review | reasoning | GPT-6.1 Sol | `xhigh` | 160,790 | 6 min 29 s |
| Tests, code, mutants, runs | execution | Claude Sonnet 5.5 | not reported | not reported | 18:34 to 19:17 |
| Round 1 | reasoning | Claude Opus 5.5 | `xhigh` | 3,010,907 (USD 2.38) | 8 min 47 s |
| The fix of round 1 | execution | Claude Sonnet 5.5 | not reported | not reported | 19:26 to 19:30 |
| Round 2 | reasoning | Claude Opus 5.5 | `xhigh` | 3,649,299 (USD 2.67) | 8 min 48 s |
| The notes of round 2, the close-out | execution | Claude Sonnet 5.5 | not reported | not reported | 19:39 to 19:50 |
