# T-yngf — a restart that fails releases its lease

Issue [#479](https://github.com/pharzam/layup/issues/479), a child of row 68b ([#370](https://github.com/pharzam/layup/issues/370)). It serves `F-0003#45`, through REQ-005 (the task loop). It has no plan row (see round 2, note 7).
Budget maximum 500 lines over 16 files, close-out inside, under ADR-0028; cycle cap 1. Goal count 1: after a failed restart that still owns its lease, the next fresh restart takes the released lease without a takeover wait.
Evidence: [test runs](../../runs/T-yngf/test-runs.md), [conditions](../../runs/T-yngf/conditions.md).

## Plans and reviews

- Plan [6338354590](https://github.com/pharzam/layup/issues/479#issuecomment-6338354590), and the author's correction [6338360777](https://github.com/pharzam/layup/issues/479#issuecomment-6338360777) of its last risk (the takeover rule times the heartbeat by the waiting run's own clock).
- Plan review [6338496385](https://github.com/pharzam/layup/issues/479#issuecomment-6338496385) (approve-with-conditions, six conditions, all met: `runs/T-yngf/conditions.md`).
- Round 1 [6339265866](https://github.com/pharzam/layup/issues/479#issuecomment-6339265866): material, finding F1.
- The Operator chose to fix F1 and run round 2, the last round the cap allows, in the coordinator's session; the Operator approves the exact landing figure at merge.
- Round 2 [6339953673](https://github.com/pharzam/layup/issues/479#issuecomment-6339953673): nothing material in scope.

## What is delivered

- `endLease` (`internal/run/steps.go`) ends every restart: it stops and joins the heartbeat, reads the lease, and releases it through Fencing when the row names this run, says `held` and no loss is known, whatever the last row is and whether or not Start was done. It uses a context that the steps' cancellation does not cancel, bounded by one minute. An error of the release is named once after the failure (`…; the lease: …`), a `LostError` too; a last row `done` becomes `fail` as before.
- `noteLoss` keeps the typed `LostError` (two refusals while the row names the run count). Every fenced records path of the run notes it: `(*state).fenced`, `Sessions.commit` and `holdsLease` through `Sessions.Lost`, the beat and the release.
- The rule has one home, `docs/spec/run.md` (The lease and fencing); the restart prose refers to it, `intake.md`, `scaffold.md` and the known limit (narrowed to Start) agree. The tests are in the acceptance table, the PRD (REQ-005), traceability, and one lesson in `guardrails.md` §2.
- The five assertions that expected a held lease after a failure now expect `released` (the old-tests rule: the requirement changed, because the held lease made the next restart wait `3 × lease.H`).

## Findings, and their closing

Round 1, F1 (material): a loss found by a fenced write outside `r.write` (the probe's, through `Sessions`) was not remembered, so the end tried a release after a terminal loss. Fixed in the same branch, test first (`TestALossFoundByAProbeWriteIsKnownAtTheEnd`, red then green in the test runs, section 5). Round 1's notes 1 to 4 and 7 are also round 2's notes 1 to 4 and 7 below.

Round 2 raised no finding. Its notes:

- Notes 1 to 4 (a direct `Fenced` call survives in the unreachable `retrospectiveStep`; only the probe has a falsifier among the loss-noting paths; a loss first found by the joined beat is not named on a `fail` row; a timing margin in two `midRestart` tests): recorded with no change, after the last round. The follow-up issue is [#484](https://github.com/pharzam/layup/issues/484).
- Note 5 (the acceptance row of `run.md` did not name a loss found by a session's write): applied as a document change.
- Notes 6, 8 and 9 (the in-flight beat cancelled by the stop is the old rule of `leaseStep`; the landing figure; history that still says "held"): closed with no change. The landing figure is 640 lines over 20 files against the limit 600/20, so the Operator approves it at merge.
- Note 7 (no plan row for `T-yngf`): no plan row.
- Round 1's note 1 (the world's operator resolves `HEAD` in a directory with no `.git`, and git can walk up into an enclosing repository when `TMPDIR` is inside the checkout) is a revealed, off-path defect, filed by the coordinator as its own issue and not fixed here.

## Verification and verdict

Red first, then green, with `date -u` times in the test runs. Verified by running: build, vet (three tags), gofmt, unit, the full integration level (twice: exit 0), the e2e level (exit 0), the six lints, `git diff --check`, and the freeze check (K1 to K4 and K6 ok; K5 reports 640/20, the figure above).

**Verdict: nothing material in scope** (round 2, comment 6339953673, Claude Fable 5.1 on OpenCode). Goal count 1. Harness exception: ADR-0026's routing rule names `claude`, `devin` and `opencode`; Codex is used for the plan review on the Operator's message of 2026-10-13.

## Resource record

ADR-0007, recorded, not budgeted. Tokens as the harness reported them; build, fix and close-out times are commit times.

| Part | Tier | Model | Effort | Tokens | Elapsed |
| --- | --- | --- | --- | --- | --- |
| Plan | reasoning | Claude Opus 5.5 (the coordinator) | not recorded | not kept | not kept |
| Plan review | reasoning | GPT-6.1 Sol (Codex) | xhigh | 156,069 | 2026-10-22T10:18:21Z to 10:25:21Z |
| Build | execution | Claude Sonnet 5.5 (subagent) | not recorded | not kept | 2026-10-22T10:33:48Z to 10:42:31Z (first to last commit) |
| Review round 1 | reasoning | Claude Fable 5.1 (OpenCode, max) | max | 16,247,091 (USD 10.53) | 2026-10-22T10:43:54Z to 11:07:50Z |
| Fix of F1 | execution | Claude Sonnet 5.5 (subagent) | not recorded | not kept | 2026-10-22T11:09:50Z to 11:17:46Z (merge of main to last commit) |
| Review round 2 | reasoning | Claude Fable 5.1 (OpenCode, max) | max | 26,304,334 (USD 14.69) | 2026-10-22T11:20:08Z to 11:45:40Z |
| Close-out | execution | Claude Sonnet 5.5 (subagent) | not recorded | not kept | not kept; its commit follows 76bdbc3 at 2026-10-22T11:49:30Z |
| **Total** | | | | 42,707,494 reported (the three numeric entries); the complete total is not reported | not reported |
