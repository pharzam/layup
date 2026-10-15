# T-92lm — the prompt of the prose session gives the rules of the checks

Issue: [#326](https://github.com/pharzam/layup/issues/326), a defect that the demo of `M2d` (row 59b, `T-sdtb`, #284) revealed, off its path, so a task of its own (gate step 5), as `T-gkvi` was. Serves `F-0003#42` through REQ-002. Base `b7f3f7d`. Author: Claude Sonnet 5.5 on Claude Code. Evidence: [`runs/T-92lm/`](../../runs/T-92lm/).

## What was done

The plan (R12), its review (GPT-6.1 Sol on Codex CLI, `high`, 4 min 36 s: `approve-with-conditions`, 800 lines over 18 files, cap 1, six conditions, all applied) and the author's answer are comments on #326; the plan is reasoning-tier work
that an execution-tier model wrote, a routing mismatch by choice. The demo of `M2d` failed twice at the step `setup` (runs 3 and 4, S08): the prompt of the prose session gave none of the rules of the five checks that accept its files, and
S14 would have refused 154 to 222 hits of the baseline's words. The prompt (`internal/run/prose.go`) now gives each rule with the files it governs, and a search command, called an aid; `internal/verify` gives it its words
(`AdaptedPhrases`, `AdaptedAllowed`, `AdaptedExempt`, `KitPhrases`), each probed against the check by a test. Documents: `scaffold.md`, the pitfall of `guardrails.md` §2, `traceability.md`, REQ-002. Tests:
[`test-runs.md`](../../runs/T-92lm/test-runs.md) (red runs, 52 mutants each alone on a copy, all caught, the gate). Evidence: [`evidence.md`](../../runs/T-92lm/evidence.md) (the protocol, the control that fails S08 as runs 3 and 4
did, attempts 1 and 2). **Known limits:** a model reads the text, so a session can still leave a hit, and two attempts are two samples; a failed check still ends `setup` as `fail`, with a lease held for 30 minutes.

## Review rounds and verdict

GPT-6.1 Sol on Codex CLI, `xhigh` (ADR-0026 does not name it; the Operator allowed it). Round 1 (`d7eff34`, comment 6176783005) `material`; round 2 (`4fbfc06`, comment 6177919651) `not mergeable, findings recorded` at the cap of 1; the Operator
chose one more cycle (Telegram message 255, cap 2, comment 6185512430); round 3 (`ec62a6d`, comment 6185638824) `nothing material in scope`. The pull request was then behind `main` (`T-dnge` landed first) and the repository requires an up-to-date branch, so the branch took `origin/main` by a merge (`f23e79e`), which changed two files of the branch: round 4 (comment 6186502305), scoped to the merge, cap 3 chosen by the Operator (comment 6186394994), gave `nothing material in scope` with one note, declined (the line of `T-92lm` in the completed log stands where it avoids a conflict, as commit `3f223ca` says; the log is out of date order by two entries). While the answer was awaited the author had Sol read the branch four times, privately (no round,
no record; five, one, one and one findings, all fixed, listed in `test-runs.md`). Attempts 1 (`4fbfc06`) and 2 (`ec62a6d`) each passed the step `verify` first time. Budget at landing: 799 lines over 18 files against `f1ce7a4`. Delivered; next: the demo (`T-sdtb`, #284).

## Resource record

- Recorded, not budgeted (ADR-0007); UTC, 2026-10-14 and 15; Sol's tokens are the "tokens used" of Codex CLI, the author's are not reported. The plan (reasoning, a routing mismatch): Claude Sonnet 5.5 `max`, 18:41 to 18:44. The plan review: Sol `high`, 158,441 tokens, 4 min 36 s. Rounds 1, 2, 3: Sol `xhigh`, 224,137, 209,886 and 216,495 tokens, 9 min 46 s, 11 min 47 s and 8 min 16 s; round 4 (scoped to the merge): Sol `high`, 76,245 tokens, 5 min 54 s.
- The tests, code, mutants, documents, gates and records (execution): Claude Sonnet 5.5 `max`, 18:54 on 2026-10-14 to 07:05 on 2026-10-15, with the waits for the Operator's answers for the Operator's answer at the cap.
- The four private reads: Sol `xhigh`, 242,830, 238,187, 250,509 and 223,826 tokens, 9 min 40 s, 8 min 12 s, 8 min 25 s and 7 min 53 s. The real sessions (execution, Claude Opus 5.5): the control 15 min 46 s, USD 3.76; attempt 1 9 min 31 s, USD 3.49; attempt 2 5 min 19 s, USD 3.43 (list basis; subscription).
