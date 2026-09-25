# T-7qvc — the software architecture of LAYUP and ADR-0013 to ADR-0017

Issue: [#66](https://github.com/pharzam/layup/issues/66) (parent
[#42](https://github.com/pharzam/layup/issues/42), the PDR, child 4). Serves
`F-0003#51`. Evidence: [`runs/T-7qvc/`](../../runs/T-7qvc/test-runs.txt) (the
panel's outputs, the selection table, the test runs).

## Test runs

Run in the worktree of this task on 2026-09-25, base `9beff4c`.

| # | State of the tree | Command | Exit | Output |
| - | ----------------- | ------- | ---- | ------ |
| 1 | The panel: three members in parallel, five minutes each (O-18, O-51) | `devin -p … gpt-6-sol-xhigh` (A); `devin -p … claude-opus-5-5-high` (B); `opencode run --standalone … grok-4.7#high` (C) | 0; 0; killed | A and B answered within 300 s (18,409 and 16,729 bytes); C failed at once (`Transport: Unable to connect`) and was dropped; its one rerun on the fallback binding, `claude -p --model claude-opus-5-5`, answered within 300 s (16,163 bytes). The three outputs are on #66 and under `runs/T-7qvc/` |
| 2 | Five ADRs present, no index rows | `sh docs/adr/adr-lint.sh` | 1 | five `FAIL … no row for it in README.md's index table` lines (in `runs/T-7qvc/test-runs.txt`) |
| 3 | The five index rows, `0018` next, ADR-0011 `Accepted; amended` | the same | 0 | `adr-lint: OK` |
| 4 | The architecture document, the PRD, the pointers, the glossary rows, the task record and the evidence present, before the commits | `sh docs/adr/adr-lint.sh`; `sh docs/prd/prd-lint.sh`; `sh docs/links/link-lint.sh`; `sh docs/tests/run-discipline-tests.sh`; `sh docs/setup/setup-check.sh`; `git diff --check` | 0 each | `adr-lint: OK` with no `WARN`; `prd-lint: OK`; links resolved; `81 passed, 0 failed`; every check `OK`; nothing (in `runs/T-7qvc/test-runs.txt`) |
| 5 | Frozen head `dc19860`, round 1 (cycle 0, cap 1), row 1 | `devin -p … gpt-6-sol-xhigh …`, in a disposable clone | 1 | no record after 186 s: `Your weekly usage quota has been exhausted … resource_exhausted`; skipped and recorded (rule 4); Devin exhausted for its weekly window |
| 6 | The same head, row 2 | `claude -p --model claude-opus-5-5 --dangerously-skip-permissions`, the same clone, a polling watchdog of 900 s | 0 | the review record, 496 s: `material`, 8 findings and 13 notes (posted on #66) |
| 7 | The fixes of round 1 present: the checks pinned to the App's identity (ADR-0013 D2), the row commit and the checked head (D3), the agents never the LAYUP App (ADR-0014 D1), the transition per kind (ADR-0015 D2–D3), the panel attributions (ADR-0013, 0014, 0015 Context; `selection.md`), the architecture citations (§4–§7), `operator-decisions.md` under `runs/T-7qvc/`, the glossary row; the thirteen notes | the six checks | see `runs/T-7qvc/test-runs.txt` | recorded there at the fix commit; expected: every check `OK`, `adr-lint` with no `WARN` |
