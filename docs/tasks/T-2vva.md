# T-2vva — the step bet

Issue [#366](https://github.com/pharzam/layup/issues/366), row 66j of M2e. Serves REQ-012 and REQ-003 (F-0003#51, F-0003#54), and REQ-005.
Plan 6250593948 (Opus 5.5); plan review 6250697855 (GPT-6.1 Sol, eight conditions, Q1 to Q18); the split note 6250701540 (66u, #418; 66t, #417, moved to M2g by the Operator, 6251940370); Q15 allocated to row 66k (6250705338); the note of #416 (6251478039); plan amendment 1 (6293953072); the author's note 6295169965; the goal count 2 (O-214).
Evidence: [test runs](../../runs/T-2vva/test-runs.md), [mutants](../../runs/T-2vva/mutants.md) and [conditions](../../runs/T-2vva/conditions.md).

## Delivered and boundaries

- [x] **The step.** `bet` follows `shape` (`internal/run/bet.go`): the milestone from the first tag of the PRD's `## 9.` table, `betBrief` of row 66u, the post by `postOnIssue`, the answer of the approver (the earliest valid first copy by comment ID, after the brief), the settled decisions, the version of the inventory, the rendered commit, the architecture tasks and `bets.tsv` last (`bets.go`, `rendered.go`, `archtasks.go`).
- [x] **The rendered commit and its push by code.** One deterministic commit on the verified head (the PRD and `docs/facts/<F-NNNN>-numbered-facts.md`), the events `push` and `bound` with no session, the branch tip checked (`errPushRefused`), `layup/verify` carried; the readers of a `bound` with no session are guarded.
- [x] **The architecture tasks.** `openTask` parameterized from `openShapeTask`; one task per settled decision with its binding `tasks/<task>/options/<decision>.tsv` in the register commit.
- [x] **Recovery.** A stop after any act resumes inside `bet`; the answer the inventory names is kept; a pending push reaches `bet` before `shape` settles the head again.
- [x] **Plan amendment 1.** The specification instruction asks for `## 6.`, `## 7.`, `## 8. Success metrics`, `## 9. Rollout & phases` (validated by `phaseTagsOf`) and a `**Non-goals.**` paragraph; the last of these is decided here.
- [x] **Docs.** `loop.md` (The step bet), `records.md`, `run.md`, `packages.md`, row 66j in `docs/plan/README.md`, PRD-0001 §12 cells, `traceability.md`, one entry of `guardrails.md` §2.
  Not delivered: the rendered task of a later bet (row 66t, #417, in M2g; a later bet fails closed), the new specification task a changed class reopens and the "waits for a bet" representation (row 66k), the brief itself (row 66u).

## Verification and verdict

Tests and code were committed together (the hook runs `go test`); the red transcripts of the first version were recorded in the commits that followed, and those of the fixes of rounds 2 and 3 after their commits, as `runs/T-2vva/test-runs.md` says. The mutants of the first version, of the fixes and of the cases of round 3 are late evidence (hashes before, during and after, with controls, in `runs/T-2vva/mutants.md`); G6 and G7 survived the first run and were killed after round 1. The full integration level was green on the final head (2026-10-20T18:15:14Z to 18:18:18Z, exit 0); e2e, the six lints and gitleaks were green. The freeze check passes K1 to K4 and K6 and fails K5 on the files (42 against 35).

Round 1 (GPT-6.1 Sol, xhigh, Codex), comment 6295520491: **material** (F1 to F4: the recovery after the rendered push, the coverage, the phase table, the branch tip), fixed. Round 2, comment 6296445019: **not mergeable, findings recorded** (F5 and F6: the separator row and the answer already applied) at cap 1; the Operator chose end 1, cap 2 (6296467310); fixed. Round 3, comment 6296840335: **not mergeable, findings recorded** (one coverage finding of three refusals of `bets`) at cap 2; the Operator chose end 3 (6297864154): fixed with no further round. Production code did not change in that fix.
**Verdict: not mergeable, findings recorded; the finding was fixed under the Operator's end 3.** The Operator's goal count is 2. Harness exception: ADR-0026's routing rule names `claude`, `devin` and `opencode`; Codex is used on the Operator's message of 2026-10-13.

## Resource record

ADR-0007, recorded, not budgeted. Tokens as the harness reported them.

| Part | Tier | Model | Effort | Tokens | Elapsed |
| --- | --- | --- | --- | --- | --- |
| Plan and amendment 1 | reasoning | Claude Opus 5.5 | not recorded | not kept | not kept (posted with 6250593948 and 6293953072) |
| Plan review | reasoning | GPT-6.1 Sol | xhigh | not kept | not kept (posted with 6250697855) |
| Build attempt | execution | Gemini 4 Argon High (`gemini-4-argon-high`) on Devin CLI | high | not kept | 2026-10-20T14:37:57Z to 14:46:07Z, stopped on the daily quota, no work |
| Implementation and the fixes of rounds 1 and 2 | execution | Claude Sonnet 5.5 (subagent) | not recorded | 625,346 (cumulative, reported at the round-2 fix) | 2026-10-20T14:47Z to 17:11Z |
| Review round 1 | reasoning | GPT-6.1 Sol | xhigh | 280,638 | 2026-10-20T16:00:36Z to 16:14:06Z |
| Review round 2 | reasoning | GPT-6.1 Sol | xhigh | 284,061 | 2026-10-20T16:46:32Z to 17:00:07Z (end from review.md's time; the launcher died after the review) |
| Review round 3 | reasoning | GPT-6.1 Sol | xhigh | 238,737 | 2026-10-20T17:13:00Z to 17:20:35Z |
| Fix of round 3 and close-out | execution | Claude Sonnet 5.5 (subagent) | not recorded | not kept at writing | 2026-10-20T18:13:54Z to 2026-10-20T18:21:04Z |
