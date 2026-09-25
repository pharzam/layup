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
