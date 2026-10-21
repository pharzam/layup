# T-89xf — a task whose base moved: the carry and the redo

Issue [#424](https://github.com/pharzam/layup/issues/424), row 66x of M2e (K3, a size child of row 66k, [#367](https://github.com/pharzam/layup/issues/367)). Serves REQ-005, REQ-007 and NFR-001 (F-0003#45, F-0003#47).
Plan 6262395410 (Opus 5.5); plan review 6262560771 (GPT-6.1 Sol, eight conditions, Q1 to Q14); the Operator's decisions 6262587821 (goal count 3; the rendered commit stays a stop in M2e, follow-on row 66t); the registration of the children 6252887262.
Evidence: [test runs](../../runs/T-89xf/test-runs.md), [mutants](../../runs/T-89xf/mutants.md) and [conditions](../../runs/T-89xf/conditions.md).

## Delivered and boundaries

- [x] **The carry.** `BringUpToDate` merges the default branch into a verified head by code (`git.Merge`, `git.MergeBase`), pushes the merge as code (events `push`, `bound`, session `—`) and carries `layup/verify` (`mergeEdge`, `CarryInput.Merge`); an outstanding push is resolved first (`internal/run/baseup.go`, `closeout.go`).
- [x] **The redo.** Overlap or conflict: an intent file, the closes, the new issue, then one records commit (register row replaced in place, `closed`, `rebased`, bindings, pin); the first attempt starts at the rebased base (`internal/run/redo.go`).
- [x] **The older attempt's pull request** is closed by `RepairTask` once the register names the newer one; `OpenPullRequest` adopts a pull request at a merge of code.
- [x] **Records and docs.** The detail of `closed` and `rebased` is a task ID; `FetchBranch`; `loop.md` (A task whose base moved), `records.md`, `packages.md`, row 66x, PRD-0001 §12 cells, `traceability.md`, one entry of `guardrails.md` §2.
  Not delivered: the redo of a task that holds the rendered commit (a stop in M2e, row 66t), the loop and restart scheduling (row 66k), the close-out of the merged task (row 66y, K4: not landed, so its `upToDate` and `UpToDate` are not folded; K4 folds them, condition 4 of #425's review).

## Verification and verdict

Tests and code were committed together (the hook runs `go test`); the red transcripts of the first version are in the commits after, and those of rounds 1 to 3 before their fix commits, in `runs/T-89xf/test-runs.md`. The 50 compiled mutants are late evidence (hashes, controls, in `runs/T-89xf/mutants.md`); four of the first set survived or did not build and were repaired. The full integration level was green on the final code head (exit 0, 21 packages), e2e, the six lints and gitleaks were green. The freeze check passes K1 to K4 and K6 and fails K5 only because it still counts `runs/` (ADR-0028 takes it out of the figure: 27 files and under 2,800 lines before the close-out).

Round 1 (GPT-6.1 Sol, xhigh, Codex), comment 6300001114: **material** (F1 to F4: a pushed merge with a lost `bound`, `OpenPullRequest` after a carry, the Shape pin, two integration paths), fixed. Round 2, comment 6300657349: **not mergeable, findings recorded** (F5 to F8: an unexplained branch head, the rendered edge, the opening pin and its link limit, three PRD rows) at cap 1; the Operator chose end 1, cap 2 (6301275363); fixed. Round 3, comment 6301826099: **not mergeable, findings recorded** (F9 and F10: the opening inputs captured once, a pin of full hashes) at cap 2; the Operator chose end 3 (6301934353): fixed with no further round, test-first, with mutants.
**Verdict: not mergeable, findings recorded; the findings were fixed under the Operator's end 3.** The Operator's goal count is 3. Harness exception: ADR-0026's routing rule names `claude`, `devin` and `opencode`; Codex is used on the Operator's message of 2026-10-13.

## Resource record

ADR-0007, recorded, not budgeted. Tokens as the harness reported them.

| Part | Tier | Model | Effort | Tokens | Elapsed |
| --- | --- | --- | --- | --- | --- |
| Plan | reasoning | Claude Opus 5.5 | not recorded | not kept | not kept (6262395410) |
| Plan review | reasoning | GPT-6.1 Sol | xhigh | not kept | not kept (6262560771) |
| Implementation and the fixes of rounds 1 and 2 | execution | Claude Sonnet 5.5 (subagent) | not recorded | 583,787 (cumulative, reported at the round-2 fix) | 2026-10-20T18:45Z to 21:10Z |
| Review round 1 | reasoning | GPT-6.1 Sol | xhigh | 254,929 | 2026-10-20T19:39:04Z to 19:49:20Z |
| Review round 2 | reasoning | GPT-6.1 Sol | xhigh | 451,762 | 2026-10-20T20:11:22Z to 20:22:21Z |
| Review round 3 | reasoning | GPT-6.1 Sol | xhigh | 257,341 | 2026-10-20T21:14:27Z to 21:25:18Z |
| Fix of round 3 and close-out | execution | Claude Sonnet 5.5 (subagent) | not recorded | not kept at writing | 2026-10-20T21:31:55Z to 2026-10-20T21:54:27Z |
