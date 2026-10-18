# T-8bwd — the handoff item contracts and implementation ownership

Issue [#387](https://github.com/pharzam/layup/issues/387), PR [#398](https://github.com/pharzam/layup/pull/398).
One specification artifact, serving REQ-005/F-0003#45. Authors: SWE-2 Max;
Codex coordinated evidence, main integration and close-out. Initial plan6220167894,
plan review6220326750 (approve-with-conditions), author answer6220408093 before edits.
Maximum1,000 changed lines/14files, cycle cap1, goals1, no panel remain unchanged.
The measurement basis is current main, confirmed in the plan review.
Evidence: [the one-off checker and controls](../../runs/T-8bwd/).

## Delivered

The single operative section in `loop.md` maps all fourteen handoff items to
inputs, rules, owning implementation and milestone. Field lines and the
`tests:`/`findings:` block markers have one pinned grammar, required cardinality,
block boundaries, zero-test spec-task form and normative malformed cases.
`internal/handoff` computes schema/artifacts; the existing result recorder owns
question-row/answer-row. Other results are computed by their declared callers;
a missing result is false. #396 owns the four comment items, #404 the two tree
items and #405 the question/answer body composition; #367 owns their posting
and step-loop wiring.
The plan and package table identify those owners and dependencies.

The typed per-test result remains #361's M2e contract; its runtime producer and
ID-to-test mapping are later work. No M2e build-task handoff is claimed valid.
Patch-per-kind, register-rows and heading-ids retain the specified later-milestone
owners. Existing Operator-approved goal counts are unchanged. #396 gets its own
initial plan and independently judged goal count; this task does not approve it.
The evidence checker is one-off, with no hook, CI job or production gate.
PRD Test cells/traceability name tests; this documentation task adds none there.

## Evidence and review

The missing mapping produced44 initial assertions; the completed artifact passes.
Independent controls expose duplicate mappings, missing Task/Demo ownership,
missing owner tokens, wrong milestones and changed verdict words, with original
bytes restored. The current-main union passed all10 local checks, five selected
Go document-reader integration packages, the item checker and standalone gitleaks.
The earlier unfiltered integration attempt was interrupted, never credited.
No whole-module integration/e2e run is claimed for this documentation-only change.
Final combined-head CI remains a required landing condition, recorded on the issue.

Round1, record6221128932, found two material grammar gaps; its table is intermediate
material and its raw findings remain. The author answered6221153143 before fixing
the block boundaries, marker grammar, empty-list form and acceptance examples.
Round2, record6233114681, onac84165dabe795ff8fa9a8e0055daad277ddfbdb,
cycle1/cap1: **nothing material in scope**. Main3c9c4be (#351) was merged before
that round. No source fix followed; this last commit is bookkeeping only.
Complete current-main size: 492 changed lines/10 files, within1,000/14.

Notes N1–N6: literal marker/CRLF behavior stays fail-closed; #396's own plan pins
normalization. The section finder is now known to be unexported, so #396 exports
and reuses the predicate per the approved ownership choice; the later heading-ID
contract stays later work. Composer/poster ownership remains explicit. Dependency
notifications went to #362/#367/#368. The typed-results ownership-cell exception
remains disclosed. Temporary checker files stay one-off and unwired. No new
operative rule or scope was introduced through these dispositions.

## Resource record

Reported subset only. SWE P/O/C = prompt/output/cached; C is inside P. Sonnet
I/O/Cw/Cr = uncached input/output/cache creation/read; thinking is inside O.

| Part | Expected tier | Model | Effort | Reported tokens | Elapsed |
| --- | --- | --- | --- | --- | --- |
| Plan/build/repairs, unique retained union | reasoning/execution | SWE-2 Max | not reported | P14564221/O84120/C14058651 | not reported |
| Independent plan review and same-context correction | reasoning | Sonnet5.5 | not reported | I72/O31331/Cw128723/Cr3209850 | 301511ms |
| Code round1 | reasoning | Sonnet5.5 | not reported | I36/O11939/Cw79441/Cr1104706 | 217607ms |
| Code round2 | reasoning | Sonnet5.5 | not reported | I80/O31899/Cw131333/Cr3151815 | 407567ms |
| Root claim/pre-review/final-freeze checks | — | deterministic | — | not applicable | 289087.370ms |
| Coordination/producing audit/integration/close-out | reasoning/execution | Codex GPT-6 family | not reported | not reported | not reported |
| **Total reported subset** | | | | 22529566 tokens; cached counted once | 1215772.370ms |

SWE union132 retained metric rows includes overlaps only once; the interrupted
build export was unavailable, with missing usage/time not inferred. The requested
models did not reach the named reasoning tier; reviewer independence and actual
model use remain explicit. The final reviewer disclosed its own harness overflow
files and the checker's transient files; no earlier review/private transcript or
sibling-task read was found. Direct producer exits were retained. Final CI and
unreported coordination are excluded from these subset totals.
