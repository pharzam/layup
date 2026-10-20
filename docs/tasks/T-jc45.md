# T-jc45 — the integration level's CI bound is 20 minutes, a stopgap until #439

Issue [#442](https://github.com/pharzam/layup/issues/442), a gate task outside the plan. The Operator chose option (c) on 2026-10-20: this stopgap now, and #439's parallel tests later. Serves F-0003#47.
Plan 6291749149, revision 1 6292722941 (Opus 5.5); plan review 6292012272 (GPT-6.1 Sol, five conditions, all met); maximum 200 lines over 11 files, cap 2, no panel.
Evidence: [test runs](../../runs/T-jc45/test-runs.md).

## Delivered and boundaries

- [x] **ADR-0027** amends ADR-0010: the integration invocation's per-test-binary limit is `-timeout 20m`; unit, e2e and the hook keep 10 minutes. Only ADR-0010's Status line and the index changed.
- [x] **`ci.yml:316`** carries `-timeout 20m`.
- [x] **The restatements** in `test-levels.md`, `scaling-checklist.md`, `engineering-discipline.md` and `guardrails.md` link to ADR-0027, and one pitfall entry in §2 says to read a package's integration time in CI's log.
  Not delivered: the parallel tests (#439).

## Verification and verdict

The red was recorded first: CI timed out at 600.012 s (`2866a3f`, run 39422675099) and at 600.014 s (`614bef4`, run 39425720339). The green: PR #443 at `79169f4` (run 39432743311) passed `internal/run` in 588.190 s under the 20-minute limit, with the 13 check runs named; the local checks exit 0; gitleaks exit 0 before each push.

Round 1 (GPT-6.1 Sol, xhigh, Codex), comment 6292708351: **material** (F1, the plan without the review's corrections; F2, the 13 checks not named), fixed by revision 1 and the named runs. Round 2, comment 6292830195: **nothing material in scope**.
**Verdict: nothing material in scope.** Goal count 1. Harness exception: ADR-0026's routing rule names `claude`, `devin` and `opencode`; Codex is used on the Operator's message of 2026-10-13.

## Resource record

ADR-0007, recorded, not budgeted. Tokens as the harness reported them.

| Part | Tier | Model | Effort | Tokens | Elapsed |
| --- | --- | --- | --- | --- | --- |
| Plan, revision 1, green record, close-out | reasoning | Claude Opus 5.5 | not recorded | not kept | 2026-10-20T13:15Z to 14:09:42Z |
| Plan review | reasoning | GPT-6.1 Sol | xhigh | 115,847 | 13:22:56Z to 13:29:16Z |
| Implementation | execution | Gemini 4 Argon High (`gemini-4-argon-high`) on Devin CLI | high | not kept | 13:30:29Z to 13:42:13Z |
| Review round 1 | reasoning | GPT-6.1 Sol | xhigh | 121,207 | 13:56:32Z to 14:03:20Z |
| Review round 2 | reasoning | GPT-6.1 Sol | xhigh | 136,410 | 14:05:13Z to 14:09:17Z |
