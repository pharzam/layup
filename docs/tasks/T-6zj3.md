# T-6zj3 — attempts, reviewer independence and frozen tests

Issue [#361](https://github.com/pharzam/layup/issues/361), PR [#400](https://github.com/pharzam/layup/pull/400).
Serves REQ-005/NFR-004/F-0003#45/#58. Author: Sonnet5.5; Codex coordinated
producing audits, main integration and bookkeeping. Plan6218735978;
independent plan review6219157382; author answer6219170199 before edits.
Maximum4,300 changed lines/36files, cycle cap1, goals3 and no panel unchanged.
Current-main basis follows #3556209446202, not the older named issue base.
Evidence: [runs](../../runs/T-6zj3/test-runs.md) and [mutants](../../runs/T-6zj3/mutants.md).

## Delivered and boundaries

A failed typed gate table supplies findings to a new attempt. The author ledger
and commit ancestry determine independent admitted reviewer pairs; PhasePair
uses the same selector. Task-tests has a checked reader/writer and built block.
Typed red/green results freeze the first valid verifier handoff's source hashes
atomically with its event under the lease fence. Source drift yields findings;
an approved event-ordered amendment can replace the list. Restarts deduplicate.
The next attempt receives recovered frozen bytes after whole-list path/symlink
preflight, with one commit only when changed. Existing file modes are not retained
(the stated0644 limit). Tier selection and prompt-file construction use the
plan/records and include source hashes, accepted handoffs/answers and findings.

These are caller APIs, not a new standalone runtime producer. #396 parses the
fixed comment forms and material-finding count; #362/#367 wire the table/callers.
No M2e runtime produces the typed per-test results or a valid build-task handoff.
Position's question count is not claimed as the whole later material-finding
count: #367 must combine its #396 input. The pre-existing AnswerAccepted helper's
owner-role omission is tracked separately as #401 before runtime wiring; this
branch does not absorb that off-path repair or claim it fixed.

## Verification and verdict

Own test-first unit/real-git integration and55 valid compiled mutants cover the
three approved goal classes. Initial compile-invalid/surviving mutants are
reported honestly and replaced with valid falsifiers, never counted as kills.
The producing-code audit found a symlink escape; own regressions and five link
mutants proved its whole-list repair. Root independently reproduced red/green.
Root full integration passed all21packages/direct0 on dff729b after that repair,
920091.157ms; the later main unions are covered by focused current-head checks
and required final combined-head CI, not retroactively by that earlier full run.
All recorded local checks and named union document-reader tests pass; final
exact-head CI remains a landing condition, with its result recorded on the issue.

Round1 on a48ec3e, record6233168862: no material source defect. Its publication
footer corrects piped-exit and scratch-location claims; only directly captured
producer exits are credited. The code verdict is preserved as raw evidence.
Round2 on a753aabc9831f64b9221766ef9b6d1503c34adcd, record6233462565:
**nothing material in scope**, cycle1/cap1, scoped to main integration after #351/#387.
No production-code fix followed the first freeze. Both sides' clauses and shared
indexes are retained; the final round reads the ordinary prose resolutions.
Only task detail/backlog/completed bookkeeping follows the final review.
Complete current-main landing: 3983 changed lines/33 files, within4,300/36.

First-round notes retain the source-backed first-bound event, missing-result
tolerance and restart recomputation; no runtime owner-role guarantee is inferred
from AnswerAccepted. Its separate #401 prerequisite and material-count boundary
remain explicit. Final-round notes: the separately tracked #401 limit remains off-path; no source correction follows.

## Resource record

ADR-0007 reported subset only. Sonnet I/O/Cw/Cr are uncached input/output/cache
creation/read; thinking is inside output. SWE P/O/C are prompt/output/cache,
with C already inside P. No cumulative-plus-increment double counting.

| Part | Expected tier | Model | Effort | Reported tokens | Elapsed |
| --- | --- | --- | --- | --- | --- |
| Plan | reasoning | Sonnet5.5 | not reported | I114/O83725/Cw316502/Cr11162910 | 693657ms |
| Independent plan review | reasoning | SWE-2 Max | not reported | P997615/O14251/C946136 | not reported |
| Implementation/test/docs/repairs | execution | Sonnet5.5 | not reported | I334/O196514/Cw282433/Cr79638196 | 2362490ms |
| Code round1, retained subset | reasoning | SWE-2 Max | not reported | P4044887/O31452/C3925696 | not reported |
| Scoped main-integration round2 | reasoning | SWE-2 Max | not reported | P6871138/O50418/C6713631 | not reported |
| Root full integration | — | deterministic | — | not applicable | 920091.157ms |
| Root claim/pre-freeze/final-union checks | — | deterministic | — | not applicable | 202537.992ms |
| Producing audits/coordination/merge/close-out | reasoning/execution | Codex GPT-6 family | not reported | not reported | not reported |
| **Total reported subset** | | | | 103690489 tokens; cache counted once | 4178776.149ms |

Sonnet totals match both the reported turn increments and final cumulative usage;
interrupted/idle/unreported coordination time is not inferred. SWE round1's earlier
compacted history is unrecovered; own extracted input copies matched supplied
bodies, with missing earlier tool history not claimed audited. Actual requested
models do not reach the named reasoning tier; the routing limitation and reviewer
independence remain explicit. Subset totals are not whole-session totals.
