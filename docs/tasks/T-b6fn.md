# T-b6fn — a pre-freeze check of the defects a round keeps finding

Issue [#440](https://github.com/pharzam/layup/issues/440), child B2 of the M2e throughput investigation [#438](https://github.com/pharzam/layup/issues/438). Started by the Operator's instruction quoted on #438 (ADR-0026 rule 1: a check script). Serves F-0003#47 and R5.
Plan revision 2 (6292340511) and amendment 1 (6292629041), by Opus 5.5. The first plan was rejected by the plan review 6291820237; revision 2's review 6292502682 approved it with five conditions, all met. Maximum 1,300 lines over 170 files; cap 2; no panel. The Operator's rulings: goal count 6 in one task (6292311744), criterion 3 corrected (6292613396), end 3 at the cap (6295306362).
Evidence: [test runs](../../runs/T-b6fn/test-runs.md), [real heads](../../runs/T-b6fn/real-heads.md), [conditions](../../runs/T-b6fn/conditions.md).

## Delivered and boundaries

- [x] [`docs/tests/freeze-check.sh`](../tests/freeze-check.sh), run on a clean, committed task head (exit 0, 1, or 2 on a usage or environment error), with six checks:
  - **K1:** white space;
  - **K2:** each new test in `traceability.md` and in the PRD §12 Test cell of each requirement it covers (by name or by its file), with the PRD and §12 required whenever new tests exist;
  - **K3:** `-run` selectors quoted and resolving in the named package and build tag, read in prose, spans and fences;
  - **K4:** restore evidence of each mutant, exempted only by an explicit disposition;
  - **K5:** the budget projected at landing;
  - **K6:** each evidence span of `runs/<task>/conditions.md` resolving at `HEAD`.
- [x] Its suite, `docs/tests/freeze-check/`: 62 cases, inputs stored as `.in`, run by `run-all.sh` as a step of CI's `nested-checkout-check` job, restored from the default branch.
- [x] The gate sentence (`engineering-discipline.md`, at the freeze), `AGENTS.md`, the two READMEs, the enforcement table, `.gitattributes`, and one pitfall in `guardrails.md` §2.
- Not delivered:
  - a CI or hook check of real task heads (a written rule, run by the author);
  - K4's form of three hashes in one cell (T-h65e);
  - the reconciliation of older registrations with PRD §12, a follow-up the rounds named.

## Verification and verdict

**Red then green, per stage and per finding:**
- the frame, then each check (Devin);
- the K2 false positive found on `main`;
- rounds 1 to 3: 44/5 → 51/0; 51/5 → 56/0; 58/4 → 62/0.

**Real heads:** K3 finds #412's unquoted selector; K4 finds #415's rows without hashes at `fa7a476` and passes once fixed; K2 and K1 find #426's F7 and F8 and nothing else.

**Rounds** (GPT-6.1 Sol, xhigh, Codex):
- Round 1 (6294640592): **material**, F1 to F5.
- Round 2 (6295020505): **material**, F6 to F9.
- Round 3 (6295234110): **not mergeable, findings recorded**, F10 and F11.

The Operator chose end 3: F10 and F11 were fixed at `65960dc`, with cases red first and a mutant per rule. **That fix is one that no round read.**

**Verdict: not mergeable, findings recorded, with the merge permitted by the Operator's end 3 (6295306362), and the landing figure 2,363/238 approved by the Operator (6295483284).** Goal count 6. The author's coding run on Devin ended on its daily quota (14:46:46Z); the plan author finished the task. Harness exception: ADR-0026's routing rule names `claude`, `devin` and `opencode`; Codex is used on the Operator's message of 2026-10-13.

## Resource record

ADR-0007, recorded, not budgeted. Tokens as the harness reported them.

| Part | Tier | Model | Effort | Tokens | Elapsed |
| --- | --- | --- | --- | --- | --- |
| Plan, revision, amendment, the K2 fix, the fixes of rounds 1 to 3, evidence, close-out | reasoning | Claude Opus 5.5 | not recorded | not kept | 2026-10-20T12:55Z to 2026-10-20T16:08:45Z |
| Plan review 1 | reasoning | GPT-6.1 Sol | xhigh | 186,895 | 13:09:15Z to 13:17:02Z |
| Plan review 2 | reasoning | GPT-6.1 Sol | xhigh | 147,099 | 13:46:05Z to 13:53:17Z |
| Script, fixtures, CI step, docs, replays | execution | Gemini 4 Argon High (`gemini-4-argon-high`) on Devin CLI | high | not reported | 14:00:15Z to 14:46:46Z (ended on the daily quota) |
| Review round 1 | reasoning | GPT-6.1 Sol | xhigh | 183,984 | 15:12:55Z to 15:21:44Z |
| Review round 2 | reasoning | GPT-6.1 Sol | xhigh | 172,032 | 15:33:59Z to 15:43:05Z |
| Review round 3 | reasoning | GPT-6.1 Sol | xhigh | 204,213 | 15:51:35Z to 16:00:17Z |
| **Total** | | | | 894,223 reported | |
