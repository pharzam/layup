# T-8rm6 — point approvers, criteria and draft revision

Issue [#365](https://github.com/pharzam/layup/issues/365), PR [#399](https://github.com/pharzam/layup/pull/399).
Also closes the one-answer marker-fixture repair [#397](https://github.com/pharzam/layup/issues/397).
Serves REQ-012 and F-0003. Author SWE-2 Max; Codex performed producing-code
and evidence audits, main integration and close-out. Plan6219011001;
independent plan confirmation6219108719; author answer6219173043 before code.
Initial maximum1500 changed lines/24files, cap1, goals2, no panel.
Current-main budget basis follows #3556209446202; the maximum is unchanged.
Evidence: [test runs](../../runs/T-8rm6/test-runs.md).

## Delivered and boundaries

Counted idea-owner answers bind approval points. Named GitHub accounts resolve
to numeric IDs before deciding whether an existing Operator/idea-owner row already
serves the point; only404 is an input fallback, other forge errors stop the step.
Success criteria come from the follow-up answer when present, otherwise the batch.
Revision keeps the checked session draft, adds measured criteria, checks answered
question and quote-fact coverage through follow-up references, and commits once.
Restart adopts the completed revision and reposts without another commit; refused
files retry within the milestone cap. Positional requirement IDs and historical
pre-revision coverage references are documented rather than treated as stable IDs.
The actual later runtime approval wiring belongs to #366, not this change.
#401's pre-existing answer-owner helper remains a separately tracked prerequisite.

The #397 fixture adds only its now-required success-criteria answer; all original
marker assertions remain. No unrelated production repair is hidden in that test.

## Verification and verdict

Own semantic red/green and compiled-mutant evidence is in the linked run record.
The producing audit found numeric-login identity handling and the marker fixture;
the bounded repairs were tested before independent final review. Initial full
integration failed on the marker fixture; it remains recorded as a failure.
Root repaired full runs passed all21packages:960092.005ms, then840078.196ms
after the independent findings were repaired. The original failed run took
1080138.239ms. The last full pass preceded the #361 main union; named intersection
checks passed afterward. Final exact-head CI remains required for landing.

Independent initial round6233209670 found four material defects; its raw record
is retained. The author answered before the cap1 repair: quote-fact inheritance,
follow-up criteria, positional-ID contract clarification and restart/retry prose.
Final blind round6233744064 read frozen ba3b70c against current main9d8bb9a:
**code, tests and documents: nothing material in scope**, cycle1/cap1.
The formal verdict is **not mergeable, findings recorded**, solely because the
complete landing size exceeds the unchanged initial maximum without approval.
The original final CI then exposed two stale end-to-end fixtures. Operator
6235055886 authorized one additional cycle, cap2; the prior final record is
intermediate. The bounded three-test-file repair6233887395 and guardrails lesson
receive a fresh independent round. No production-code fix enters this cycle.
The previous2,364/28 request is withdrawn; the new exact complete figure and
any one-time R12 ruling are recorded on the issue;
no estimate, reserve or implicit approval replaces that ruling.

Later rounds: round 3 (Sonnet 5.5) ended `not mergeable`, F1-F3; the Operator set
cap 3 (6246378213); the fix c7db4f7 repaired F2 (the restart counts 22 steps) and F3
(the empty draft). Round 4 (Claude Opus 5.5 `max` on OpenCode, 6247229055) found F1,
F4 and F5 (R1 tracked as #408) and ended `not mergeable, findings recorded`. The
Operator chose end 3 (6247381864), a fix with no further round: F5 fixed (the test
is in the traceability and PRD-0001 rows) and F4 a known limit of the task, its
rule #409 (see loop.md, Done before). Delivered; the merge waits for the Operator's
one approval of the landing figure.

Final notes disposition (no code correction or additional review round):
N1 retain the first-byte implementation and test; the imprecise spec citation is
non-material and no new operative sentence is introduced after the final freeze.
N2 retain the identical intake summary and its link to the canonical loop list;
no divergent rule exists. N3 the context-only empty-draft input was outside the
DoD and tracked as #402; it is now fixed by F3 of round 3 (an empty draft is not
revised until a session row exists; test TestTheRevisionOfAnEmptyDraft). N4 the planned named no-result
test was not delivered; shared PhaseError handling is tested elsewhere, and the
review found no failing acceptance behavior. This coverage limit is explicit.
N5 these helper callers remain #366's responsibility; no runtime claim is made.

## Resource record

ADR-0007/ADR-0024 actual reported subset only. SWE P/O/C means prompt/output/cache;
cache is inside prompt. Sonnet I/O/Cw/Cr means uncached input/output/cache creation/
read, with thinking inside output. Missing exports and elapsed values are not inferred.

| Part | Expected tier | Model | Effort | Reported tokens | Elapsed |
| --- | --- | --- | --- | --- | --- |
| Plan plus retained implementation/test/repair turns, unique union | reasoning/execution | SWE-2 Max | not reported | P28279539/O162206/C27237080;258 unique metric rows | not reported |
| Independent plan confirmation | reasoning | Sonnet5.5 | not reported | I44/O23895/Cw83262/Cr1326345 | 245419ms |
| Initial code round | reasoning | Sonnet5.5 | not reported | I80/O33180/Cw136317/Cr3476486 | 371326ms |
| Final code round | reasoning | Sonnet5.5 | not reported | I72/O26025/Cw123785/Cr2776856 | 497168ms |
| Root full integration, one failed and two passed | — | deterministic | — | not applicable | 2880308.440ms |
| Root audits, coordination, main merges, focused checks and close-out | reasoning/execution | Codex GPT-6 family | not reported | not reported | not reported |
| Cap-3 fix (F2, F3) | execution | Claude Sonnet 5.5 | not reported | 88,483 tokens | 59 min (2026-10-18 06:04Z-07:03Z) |
| Round 3 | reasoning | Claude Sonnet 5.5 | not reported | not recorded | not recorded |
| Round 4 | reasoning | Claude Opus 5.5, variant max, on OpenCode | max | 15,436,326 tokens with cache, USD 7.01 | 33 min (2026-10-18T07:05:06Z-07:38:10Z) |
| Orchestration | reasoning | Claude Opus 5.5 on Claude Code | not reported | not separately recorded | not recorded |
| This close-out | execution | Claude Sonnet 5.5 | not reported | not known to the author | 2026-10-18 07:58Z-08:01Z |
| **Total reported subset** | | | | 36448092 tokens; cache counted once | 3994221.440ms |

The SWE original build export is missing; unique retained increments exclude
repeated cumulative rows. Plan and execution cannot be split reliably from those
retained aggregates. Sonnet and SWE do not reach the repository's named reasoning
tier; user model restrictions and independent model separation remain explicit.
These figures are a known subset, not a whole-session or wall-clock total.
