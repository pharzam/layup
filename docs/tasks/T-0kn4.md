# T-0kn4 — the architecture and ADR-0013 to ADR-0017: successor of T-7qvc after its cycle cap

Issue: [#67](https://github.com/pharzam/layup/issues/67); predecessor
[`T-7qvc`](T-7qvc.md) ([#66](https://github.com/pharzam/layup/issues/66), split,
O-63); parent [#42](https://github.com/pharzam/layup/issues/42), child 4. Serves
`F-0003#51`. Evidence: [`runs/T-0kn4/test-runs.txt`](../../runs/T-0kn4/test-runs.txt)
(this task) and [`runs/T-7qvc/`](../../runs/T-7qvc/selection.md) (the panel, the
selection, the Operator's decisions).

## Test runs

Run in the worktree of this task on 2026-09-26 (UTC), base `9beff4c`, branch
carried from `2b7ed6a` with the two findings and nine notes of #66's round 2
applied.

| # | State of the tree | Command | Exit | Output |
| - | ----------------- | ------- | ---- | ------ |
| 1 | The fixes applied; `T-7qvc.md` with its verdict and resource record; this file, the backlog line and the evidence file present | `sh docs/adr/adr-lint.sh`; `sh docs/prd/prd-lint.sh`; `sh docs/links/link-lint.sh`; `sh docs/tests/run-discipline-tests.sh`; `sh docs/setup/setup-check.sh`; `git diff --check` | 0 each | `adr-lint: OK` with no `WARN`; `prd-lint: OK`; links resolved; `81 passed, 0 failed`; every check `OK`; nothing (in `runs/T-0kn4/test-runs.txt`) |
| 2 | Frozen head `eea9bf6`, round 1 (cycle 0, cap 1), row 1 | `claude -p --model claude-opus-5-5 --dangerously-skip-permissions`, in a disposable clone, a polling watchdog of 900 s | 0 | the review record, 471 s: `material`, 1 finding and 8 notes (posted on #67): the `question` kind cited a clarification ADR-0016 does not define and the reply had no pair |
| 3 | The fix of round 1 present: `question` goes to the operator (REQ-006's route is a later ADR's), the stall pairs complete (a verified recovery examiner → the continuing role), the pair check against the artifact's own conditions, one wording for the review events in ADR-0013, ADR-0014, the architecture and the glossary (a withdrawn approval starts a run), the progress event as a handoff row, `stall examine` twice, the T-7qvc total as wall-clock | the six checks | see `runs/T-0kn4/test-runs.txt` | recorded there at the fix commit; expected: every check `OK` |
