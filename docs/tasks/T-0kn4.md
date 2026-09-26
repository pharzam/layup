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
| 3 | The fix of round 1 present: `question` goes to the operator (REQ-006's route is a later ADR's), the stall pairs complete (a verified recovery examiner → the continuing role), the pair check against the artifact's own conditions, one wording for the review events in ADR-0013, ADR-0014, the architecture and the glossary (a withdrawn approval starts a run), the progress event as a handoff row, `stall examine` twice, the T-7qvc total as wall-clock | the six checks | 0 each | `adr-lint: OK` with no `WARN`; `prd-lint: OK`; `link-lint: OK  914 links resolved`; `81 passed, 0 failed`; 12 `setup-check: … OK` lines, no `FAIL`; nothing (in `runs/T-0kn4/test-runs.txt`) |
| 4 | Frozen head `99c5c1e`, round 2 (cycle 1, the cap), row 1 | `claude -p --model claude-opus-5-5 --dangerously-skip-permissions`, in a new disposable clone, a polling watchdog of 900 s | 0 | the review record, 486 s: `not mergeable, findings recorded`, 2 findings (1 in the change, 1 revealed) and 7 notes (posted on #67); the Operator settled finding 1 as a note (O-65); the revealed finding is [#68](https://github.com/pharzam/layup/issues/68) |
| 5 | The note of round 2 applied (`layup handoff check` reads the plan-review and review-record verdicts itself, ADR-0015 D3; the `decision` fix pair, D2) and the seven notes; this record's verdict and resource record; the two completed-log lines; the two backlog lines removed | the six checks | 0 each | `adr-lint: OK` with no `WARN`; `prd-lint: OK`; `link-lint: OK  920 links resolved`; `81 passed, 0 failed`; 12 `setup-check: … OK` lines, no `FAIL`; nothing (in `runs/T-0kn4/test-runs.txt`) |

## Verdict

Mergeable. Round 2 (the cap) ended `not mergeable, findings recorded` with two
findings and seven notes; the Operator settled finding 1 as a note (O-65), which
this close-out applies, and finding 2 is a revealed defect off this task's path,
opened as [#68](https://github.com/pharzam/layup/issues/68) (the Git path of the
records the kit keeps as issue comments). Delivered on the branch, from
[`T-7qvc`](T-7qvc.md): [ADR-0013](../adr/0013-run-stack-gates-through-a-layup-github-app.md)
to [ADR-0017](../adr/0017-stop-a-stall-at-a-counted-limit-and-examine-it-fresh.md),
[`docs/architecture.md`](../architecture.md), the PRD's §11 and §12, the
registration, and the panel's evidence under `runs/T-7qvc/`; in this task, the
two findings and nine notes of #66's round 2, the one finding and eight notes of
round 1 (Claude Opus 5.5, 471 s, `material`), and the note and seven notes of
round 2 (Claude Opus 5.5, 486 s). Devin was not tried: its weekly quota was
exhausted on 2026-09-25 (`T-7qvc` row 5), so the order named on #67 started at
`claude -p`.

## Resource record

Recorded, not budgeted (ADR-0007). Token counts: `not reported` (`claude -p`
printed none; the author's session count was not read).

| Part | Expected tier | Model | Effort | Tokens | Elapsed |
| ---- | ------------- | ----- | ------ | ------ | ------- |
| The fixes of #66's round 2, the issue, the plan and its review, the carry (`4f10c1e`, `eea9bf6`) | reasoning (plan), execution (fixes) | Claude Fable 5.1 on Claude Code (author); the Operator (plan review) | not reported | — | not measured: the session paused between #66's round 2 (2026-09-25, about 18:50 UTC) and the commits (2026-09-26 06:10 UTC) |
| Round 1 | reasoning | Claude Opus 5.5 on Claude Code | not reported | not reported | 471 s |
| The fix of round 1 (`99c5c1e`) | execution | Claude Fable 5.1 | not reported | not reported | about 3 min |
| Round 2 | reasoning | Claude Opus 5.5 on Claude Code | not reported | not reported | 486 s |
| O-65, #68, the note and seven notes, this record, the index lines, the pull request | execution | Claude Fable 5.1; the Operator (O-65) | not reported | not reported | about 20 min |
| **Total** | | | | not reported | about 45 min wall-clock from the issue (06:07 UTC) to the close-out commit; the first part is outside it |
