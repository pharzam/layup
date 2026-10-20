# T-tnxn — the findings of `layup/spec` and `layup/rules` start a new attempt

Issue [#427](https://github.com/pharzam/layup/issues/427), row 66aa, a size child of row 66k ([#367](https://github.com/pharzam/layup/issues/367)). Serves REQ-007, REQ-005 (F-0003#47, F-0003#45).
Plan 6262334701 (Opus 5.5); plan review 6262459454 (GPT-6.1 Sol, seven conditions, all met); the Operator's decisions 6262588979 (goal count 3, kept at the maximum; the acceptance of a waiting specification
before the verifier is 66r's); maximum 1,350 lines over 20 files, cap 1.
Evidence: [test runs](../../runs/T-tnxn/test-runs.md) and [mutants](../../runs/T-tnxn/mutants.md).

## Delivered and boundaries

- [x] **The spec findings.** The `fail` rows of a head's table, the row `ran` included, as `spec:<rule>` (`SpecFindings`, `internal/run/attempt.go`).
- [x] **The rules finding.** The failure of the basic rule as `rules:<reason>` with its path as the file and the text of the status (`RulesFindings`; `RulesResult` extracted from `RulesAt`).
- [x] **The expected failures.** For the register kind `batch`, in M2e, the rules failure and the pending gate rows are no finding, an active row stays one (`expectedFailure`, `RetryGateFindings`);
  the waiting specification starts no attempt (I4).
- [x] **The retry.** A failure of `layup/spec` or `layup/rules` starts one new attempt at the same base with the findings of the gates, the spec and the rules; an error, `not-active` or a Go error
  comes first (`judgeStatuses`, `CheckFindings`, `NextAttempt` with the given findings, `headFindings`; the act `ship`); a restart starts exactly one.
- [x] **Docs.** `loop.md`, `packages.md`, row 66aa (and 66aa in the After cell of 66k), PRD §12 cells, traceability, one entry of `guardrails.md` §2.
  Not delivered: the `Restart`-level composition (66h, 66k), the native job's readiness (66y).

## Verification and verdict

Test first (inert stubs at `c4153a5`), 18 compiled mutants killed with hashes and controls (M16 by the unit test alone). Full integration level green on the code head
(2026-10-20T05:59:11Z to 07:05:56Z, exit 0); its start file was overwritten by a parallel session's run in the shared scratchpad, so the start time is the one read before that, as the reviewers noted;
e2e green; gitleaks exit 0 before each push.

Round 1 (GPT-6.1 Sol, xhigh, Codex), comment 6285589802: **material** (F1: 66aa missing from the After cell of 66k), fixed. Round 2, comment 6285869039: **nothing material in scope**.
**Verdict: nothing material in scope.** The Operator's goal count is 3. Harness exception: ADR-0026's routing rule names `claude`, `devin` and `opencode`; Codex is used on the Operator's message of 2026-10-13.

## Resource record

ADR-0007, recorded, not budgeted. Tokens as the harness reported them.

| Part | Tier | Model | Effort | Tokens | Elapsed |
| --- | --- | --- | --- | --- | --- |
| Plan | reasoning | Claude Opus 5.5 | not recorded | not kept | not kept (posted with comment 6262334701) |
| Plan review | reasoning | GPT-6.1 Sol | xhigh | not kept | not kept (posted with comment 6262459454) |
| Implementation | execution | Claude Sonnet 5.5 | not recorded | 296,291 | 140 min |
| Fix of round 1 | execution | Claude Sonnet 5.5 | not recorded | 299,975 | 6 min |
| Close-out | execution | Claude Sonnet 5.5 | not recorded | not kept at writing | 2026-10-20T07:36:08Z to 2026-10-20T07:36:29Z |
| Review round 1 | reasoning | GPT-6.1 Sol | xhigh | 177,684 | 2026-10-20T07:10:11Z to 07:18:31Z |
| Review round 2 | reasoning | GPT-6.1 Sol | xhigh | 339,037 | 07:25:29Z to 07:35:33Z |
