# T-jew3 — the code-owner rule

Issue [#356](https://github.com/pharzam/layup/issues/356), PR [#385](https://github.com/pharzam/layup/pull/385).
Serves F-0003#43 and F-0003#47 through REQ-003 and REQ-007.
Author: Claude Sonnet5.5 on Claude Code; coordinator: Codex GPT-6 family.
Plan base8e26eca; mainE16568c merged at1e9a5d5. Frozen head987c19d.
Evidence: [runs](../../runs/T-jew3/README.md), [tests](../../runs/T-jew3/test-runs.md), [mutants](../../runs/T-jew3/mutants.md).

## Plan and scope

The plan6210782039, independent Devin SWE-2 Max confirmation6214007155 and
author answer6214046545 are on #356. Maximum1000 added-plus-removed lines /24
files, close-out and fixes inside; cap2 (branch protection), goal2 under O-214.
The Operator selected .github/CODEOWNERS and prepared, reviewed Operator-run
application, updating an existing recorded ruleset ID while preserving protections
(6209820437). Current-main budgeting is the standing M2e ruling.

Both review conditions are satisfied: missing/unreadable/nonpositive old apply
records refuse before touching the area and name recovery, with tests; the
In-Scope facts are named. The six notes are addressed in source and tests:
reduced write-field readback, trap after successful disable, loud restore failure,
migration command paths, real-gh dependency, shared login grammar, no push guard.

## Delivered and verified

S13 writes the four gate-path classes owned by the Operator account and the
code-owner rule. Legacy detection rebuilds the exact old setup commits while
preparing update-by-ID commands and a reviewed landing script. Unknown rules,
parameters, conditions and bypass actors survive the update and readback.
A never-applied legacy setup refuses with recovery; no ID is guessed.
The schema/layout, setup/scaffold/loop/package specifications, traceability,
PRD cells and known pitfalls are updated in the same PR.

Tests were seen red against stubs before implementation. Both unit and real
sh/git/gh integration seams are covered; the fixtures use local TLS and bare git.
Mutation pass1 had19 named assertion kills and3 invalid build failures, excluded.
The three valid re-expressions and one distinct stronger BEFORE-file variant
also fail named assertions:23 valid syntactic variants overall, including the
explicitly weak original unbound-before variant. Restored hashes and controls pass.
The script helper's invalid-counting flaw and interrupted invocation are retained.

Full module integration on Go-identical5564b10 passed directly: wrapper
2026-10-16T11:41:10.871415Z–11:53:50.955797Z, exit0,21packagesok.
Later987c19d changes run evidence only. Final build/vet/unit, five discipline/setup
checks, diff checks and standalone gitleaks pass. Final full-SHA CI is required.
The main merge's sole PRD changelog conflict retains #350 then #356 exactly once;
round1 explicitly verified it, without claiming an O-219 exemption.

## Review and verdict

| Round | Head | Cycle/cap | Material | Verdict |
| ----- | ---- | --------- | -------- | ------- |
| 1 | 987c19d | 0/2 | 1, budget | not mergeable, findings recorded |

Fresh independent Devin SWE-2 Max round1 is comment6215471562, with context,
method, execution and model independence from both authors. It found no
material code defect. Notes F2–F5 are declined without a fix: phase1 nil coupling
is unreachable as frozen; temporary-index cleanup, quoting the diagnostic hint
and concurrent administrative changes are non-material limits of the prepared
Operator-run path. Only close-out bookkeeping follows the reviewed head.

The earlier AGY draft is uncredited (6215310494): wrong cycle and fabricated
source/rule details made it unreliable. Its resources are still recorded below.
This is one qualifying review round, not two.

Delivered as code and tests; the Operator approved the exact size; CI remains.
Whole staged landing against mainE16568c: 1617 changed lines /28 files,
including this complete close-out. The maximum remains1000/24; the Operator
approved exactly1617/28 once on #356 in response to request6215671968.
No target was changed. Production landing remains disabled pending #384;
actual GitHub binding at approval count0 is left to the #370 demo.

## Resource record

Recorded, not budgeted (ADR-0007). Thinking is included in output, not added again.
Claude I/O/Cw/Cr means uncached input/output/cache creation/cache read; AGY I/O/Cr
are its reported fields. SWE P/O/C means prompt/output/cached, with C inside P.
Missing portions remain not reported; a timeout limit is not measured elapsed.

| Part | Expected tier | Model | Effort | Reported tokens | Elapsed |
| ---- | ------------- | ----- | ------ | --------------- | ------- |
| Plan | reasoning | Sonnet5.5 | high | I92/O76025/Cw236223/Cr6202154 | 692240ms |
| Unauthorized planning tool incident | reasoning | Haiku5.5 | not reported | I18410/O1537 | not reported |
| Failed AGY plan/partial clarification, uncredited | reasoning | Gemini4ArgonHigh | high | I570486/O41460/Cr4664751 | 986227.546ms cumulative |
| Independent plan review | reasoning | SWE-2 Max | not reported | not reported | not reported |
| Build, fixes, tests and docs | execution | Sonnet5.5 | high | reported increment I188/O129700/Cw448480/Cr34455668 | finish357510ms; interrupted part not reported |
| Full integration | — | not applicable | not applicable | not applicable | 760084.382ms |
| Failed code-review attempt, uncredited | reasoning | Gemini4ArgonHigh | high | I539318/O33975/Cr5399770 | 653532.949ms |
| Qualifying code round1 | reasoning | SWE-2 Max | not reported | P2979705/O29386/C2864064 | not reported |
| Merge, evidence audit, close-out | — | Codex GPT-6 family | not reported | not reported | not reported |
| **Total** | | | | reported subtotal55827328 including cache, plus unreported portions | reported subtotal3449594.877ms, plus unreported portions |

The build row is the reported cumulative Sonnet increment from plan to finish,
not that increment plus the finish invocation again. The interrupted request may
have an unreported portion. Haiku is a policy failure, not an authorized model
exception; its counters never increase after the Sonnet tool override.
Raw private sessions stay outside git. Actual CLI report sources and incidents
are named in runs; SWE round1 metadata is retained in the coordinator audit.
