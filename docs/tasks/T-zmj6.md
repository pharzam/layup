# T-zmj6 — the gap batch of LAYUP's own PSB and the idea owner's answers (F-0004)

Issue: [#45](https://github.com/pharzam/layup/issues/45) (parent
[#42](https://github.com/pharzam/layup/issues/42), the PDR). The first task of
the product path of [ADR-0012](../adr/0012-build-layup-in-bootstrap-mode.md),
part 5; it serves `F-0003#41`. Evidence: [`runs/T-zmj6/`](../../runs/T-zmj6/drafts.md).

## Test runs

Run in the worktree of this task on 2026-09-25, base `aae7065`. Rows 1 and 2
are summaries: their raw output was not retained; from row 3 on the outputs are
in [`runs/T-zmj6/test-runs.txt`](../../runs/T-zmj6/test-runs.txt).

| # | State of the tree | Command | Exit | Output |
| - | ----------------- | ------- | ---- | ------ |
| 1 | Four fixture cases and the `good-passthrough` stand-ins added; check `facts` unchanged | `sh docs/setup/tests/run.sh` | 1 | the four new cases fail: `no output line: setup-check: facts FAIL numbering: F-0004 holds 18 distinct fact numbers in 1..19, expected 19` (and the three others) — the check does not know `F-0004` yet |
| 2 | Check `facts` extended (the answers block; `F-0004` in the index loop); the stand-in records first generated wrong (one line, a `zsh` word-splitting slip), then regenerated | `sh docs/setup/tests/run.sh` | 1, then 0 | `setup-check tests: 32 passed, 2 failed` (`bad-answers-missing`, `bad-answers-repeat`: `F-0004 holds 0 distinct fact numbers`), then `setup-check tests: 34 passed, 0 failed` |
| 3 | `F-0004` stored (19 facts) and its index row present; the batch `internal/psb/testdata/psb.tsv` unchanged (19 rows) | `sh docs/setup/setup-check.sh --only facts .` | 0 | `setup-check: facts OK` |
| 4 | A copy of the tree with fact 7 removed from `F-0004` | the same, in the copy | 1 | `setup-check: facts FAIL numbering: F-0004 holds 18 distinct fact numbers in 1..19, expected 19` |
| 5 | Frozen head `d93330d`, round 1 (cycle 0), first named harness | `devin -p --prompt-file .review-brief/PROMPT.md --model gpt-6-sol-xhigh --permission-mode dangerous --respect-workspace-trust false`, in a disposable clone | 0 | the review record, 493 s: `material`, 4 findings, 0 notes (posted on #45). Finding 1: 399 lines over 22 files against 320 / 18 — the Operator approved 540 / 27 once (O-45). Finding 2: a tab-only fact passed — the test now removes tabs, and the fixture `facts/bad-answers-blank` pins it. Findings 3 and 4: two sentences of the record. Revealed: the older loop of check `facts` has the same tab gap for `F-0001` and `F-0003` — #61, normal priority |
| 6 | The fixes present, before the fix commit: the tab-case fixture added | `sh docs/setup/tests/run.sh` | see `runs/T-zmj6/test-runs.txt` | expected `35 passed, 0 failed` (recorded there at the fix commit) |
| 7 | Fix head `457fe71`, round 2 (cycle 1, the cap), first named harness | the same `devin -p …` command, in a new disposable clone | 0 | the review record, 538 s: `not mergeable, findings recorded`, 2 findings and 1 note (posted on #45); the Operator ruled both findings notes (O-46) |
| 8 | The notes of O-46 applied (`.gitattributes`, the `F-0004` Notes); the verdict, the resource record and the completed-log line present (the close-out) | the five checks, `sh docs/setup/setup-check.sh`, `git diff --check aae7065 HEAD` | see `runs/T-zmj6/test-runs.txt` | recorded there at the close-out commit; expected: every check `OK`, the diff check exit 0 |

## Verdict

Delivered: [`F-0004`](../facts/F-0004-psb-gap-answers.md), the idea owner's
nineteen answers to the gap batch of the PSB, accepted in one batch (Decision
Point 2, `F-0001#11`; ADR-0012 part 5), and check `facts` extended to the
answers record with five fixture cases (`sh docs/setup/tests/run.sh`: 35
passed). Round 1 (GPT-6 Sol on Devin): `material`, 4 findings — the budget
(O-45), a tab-only fact, the fact-to-question binding, the Notes — fixed in
`457fe71`. Round 2 (GPT-6 Sol on Devin, the cap): 2 findings and 1 note; the
Operator settled both findings as notes (O-46, rule 3), applied in the close-out
commit. Revealed: [#61](https://github.com/pharzam/layup/issues/61), the same
tab gap in the older loop of check `facts`, normal priority. Next on the product
path: #42, `PRD-0001`.

## Resource record

Recorded, not budgeted (ADR-0007). Token counts: `not reported` (Devin prints
none; the author's session count was not read).

| Part | Expected tier | Model | Effort | Tokens | Elapsed |
| ---- | ------------- | ----- | ------ | ------ | ------- |
| The drafts of the 19 answers, the decision note, the batch acceptance | reasoning | Claude Fable 5.1 on Claude Code (author); the idea owner (acceptance) | not reported | not reported | about 15 min |
| The test slice, the record, the documents | execution | Claude Fable 5.1 — a reasoning-tier model on an execution part: the author's session did the edits | not reported | not reported | about 25 min |
| Round 1 | reasoning | GPT-6 Sol on Devin | xhigh | not reported | 493 s |
| The fixes of round 1, the fixture, #61 | execution | Claude Fable 5.1 | not reported | not reported | about 10 min |
| Round 2 | reasoning | GPT-6 Sol on Devin | xhigh | not reported | 538 s |
| The close-out (O-46, the notes, this record, the PR) | `—` | Claude Fable 5.1 | not reported | not reported | about 10 min |
| **Total** | | | | not reported | about 1 h 20 min |
