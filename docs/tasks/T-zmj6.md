# T-zmj6 — the gap batch of LAYUP's own PSB and the idea owner's answers (F-0004)

Issue: [#45](https://github.com/pharzam/layup/issues/45) (parent
[#42](https://github.com/pharzam/layup/issues/42), the PDR). The first task of
the product path of [ADR-0012](../adr/0012-build-layup-in-bootstrap-mode.md),
part 5; it serves `F-0003#41`. Evidence: [`runs/T-zmj6/`](../../runs/T-zmj6/drafts.md).

## Test runs

Run in the worktree of this task on 2026-09-25, base `aae7065`.

| # | State of the tree | Command | Exit | Output |
| - | ----------------- | ------- | ---- | ------ |
| 1 | Four fixture cases and the `good-passthrough` stand-ins added; check `facts` unchanged | `sh docs/setup/tests/run.sh` | 1 | the four new cases fail: `no output line: setup-check: facts FAIL numbering: F-0004 holds 18 distinct fact numbers in 1..19, expected 19` (and the three others) — the check does not know `F-0004` yet |
| 2 | Check `facts` extended (the answers block; `F-0004` in the index loop); the stand-in records first generated wrong (one line, a `zsh` word-splitting slip), then regenerated | `sh docs/setup/tests/run.sh` | 1, then 0 | `setup-check tests: 32 passed, 2 failed` (`bad-answers-missing`, `bad-answers-repeat`: `F-0004 holds 0 distinct fact numbers`), then `setup-check tests: 34 passed, 0 failed` |
| 3 | `F-0004` stored (19 facts) and its index row present; the batch `internal/psb/testdata/psb.tsv` unchanged (19 rows) | `sh docs/setup/setup-check.sh --only facts .` | 0 | `setup-check: facts OK` |
| 4 | A copy of the tree with fact 7 removed from `F-0004` | the same, in the copy | 1 | `setup-check: facts FAIL numbering: F-0004 holds 18 distinct fact numbers in 1..19, expected 19` |
| 5 | Frozen head `d93330d`, round 1 (cycle 0), first named harness | `devin -p --prompt-file .review-brief/PROMPT.md --model gpt-6-sol-xhigh --permission-mode dangerous --respect-workspace-trust false`, in a disposable clone | 0 | the review record, 493 s: `material`, 4 findings, 0 notes (posted on #45). Finding 1: 399 lines over 22 files against 320 / 18 — the Operator approved 540 / 27 once (O-45). Finding 2: a tab-only fact passed — the test now removes tabs, and the fixture `facts/bad-answers-blank` pins it. Findings 3 and 4: two sentences of the record. Revealed: the older loop of check `facts` has the same tab gap for `F-0001` and `F-0003` — #61, normal priority |
| 6 | The fixes present, before the fix commit: the tab-case fixture added | `sh docs/setup/tests/run.sh` | see `runs/T-zmj6/test-runs.txt` | expected `35 passed, 0 failed` (recorded there at the fix commit) |
