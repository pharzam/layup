# T-dx1h — the parameters

Issue [#358](https://github.com/pharzam/layup/issues/358), PR [#388](https://github.com/pharzam/layup/pull/388).
Serves F-0003#51 through REQ-012 and NFR-001; the handoff check also traces REQ-005.
Author: Claude Sonnet 5.5; coordinator: Codex GPT-6 family.
Goal count 2 under O-214; maximum 1,600 added-plus-removed lines /21 files, cap 1.
Plan confirmation 6214506834 sets that maximum; all four conditions are satisfied.
Evidence: [red/green, mutants and integration](../../runs/T-dx1h/test-runs.md).

## Delivered and verified

The first task writes six source-complete parameter rows with the register and
transition table in one fenced commit, keeping existing valid files. Start rows
use start.tsv and command evidence; other defaults remain unset. No brief.lines
value is guessed and the existing brief-limit gap remains registered.
The boundary copies first human comments, judges ordered set decisions by role,
name, value and bound, and commits an accepted row once. A later comment supersedes
older decisions. Edits and App comments are excluded as decisions. A transitions
set stores the exact copied table and its hash with the parameter row atomically.
The rendered-task transition cannot require schema; other rows retain that item.
The boundary is an unwired helper; the later loop driver is not delivered here.
Specs, schema ownership, package table, plan Test cells, PRD and traceability agree.

Tests failed by intended assertions before implementation; the initial stub-setup
failure is discarded and not credited as the red proof. Sweep 1: 39 attempts,
30 kills, four real survivors, one equivalent and four invalid builds. Added
regressions and valid re-expressions give eight kills in sweep 2. I7 is killed
by the helper's missing-file fatal, not a text assertion; that limit is retained.
Restored controls, unit and real-git/stand-in-forge integration pass.
Full module integration on 0c5e9f626433054347ddb4bf2d8a98a1509d4bc9:
wrapper 2026-10-16T12:20:34.367100Z–12:34:54.471589Z, direct exit 0,21 packages ok.
Later task changes are run evidence and the shared implemented-block inventory;
main contributes separately tested code. Combined frozen-head CI at 319ce9c has
successful tests, lint, security and discipline jobs; the review-record job ran
before round 2 was published. All final full-SHA checks remain required at landing.

## Review and verdict

| Round | Head | Cycle/cap | Material | Verdict |
| ----- | ---- | --------- | -------- | ------- |
| 1 | 7095100 | 0/1 | budget | material (intermediate) |
| 2 | 319ce9c | 1/1 | budget; final-CI condition | not mergeable, findings recorded |

Round 1 (6216689483) used one originally fresh independent SWE-2 Max context,
interrupted and then completed in a bounded resumed turn; no outside verdict was
supplied. Round 2 (6217067473) is a separate fresh context on a new head. Both
held context, method, execution and model independence from both change authors.
No material code defect was found. Round 2 explicitly verified all five main-merge
resolutions: both guardrail sets, both plan Test cells, both PRD Task/Test additions,
all parameter and retrospective owner predicates, and both implemented blocks.
The non-index overlaps consumed the allowed cycle; no broader O-219 claim is made.

Round 1's six notes are declined without source fixes: a leading reason-space is
verbatim suffix data; role-invalid old comments still refuse before supersession;
the pure parameter reader has no Start input and the writer binds those values;
App comments are copied by watch and excluded as decisions; schema-checked IDs
make discarded parse errors unreachable; the review's focused-coverage limit is
recorded, with full integration/CI evidence above. Notes cause no extra round.
The separate AGY draft is rejected (6216608013): invented governance quotations
and false per-file figures; its quota failure and resources remain recorded.

Code and tests are delivered; exact size is approved and final CI gates landing.
The complete close-out is 2,038 changed lines /20 files against main cfcfc1c,
including all three bookkeeping files and this resource record. Approval6217323927
is the one-time R12 exception; the original maximum remains 1,600/21.
No target was changed. The review's CI condition resolves only with final checks.

## Resource record

ADR-0007: reported, not budgeted. Thinking is included in output. Claude I/O/Cw/Cr
means uncached input/output/cache creation/cache read. SWE P/O/C means prompt/
output/cached, with cached inside prompt. AGY cache read is separate from its I/O.
Missing parts and elapsed figures stay unreported; a timeout limit is not elapsed.

| Part | Expected tier | Model | Effort | Reported tokens | Elapsed |
| ---- | ------------- | ----- | ------ | --------------- | ------- |
| Plan | reasoning | Sonnet 5.5 | high | I100/O61047/Cw202435/Cr6486617 | 553343ms |
| Independent plan review | reasoning | SWE-2 Max | not reported | P4885023/O65701/C4410496 | not reported |
| Build, tests, mutations and docs | execution | Sonnet 5.5 | high | increment I256/O158002/Cw237934/Cr46287147 | 2333621ms |
| Full module integration | — | not applicable | not applicable | not applicable | 860104.489ms |
| Code round 1, retained same context | reasoning | SWE-2 Max | not reported | P6332685/O44936/C6012860 | not reported |
| Rejected AGY attempt | reasoning | Gemini 4 Argon High | high | I376576/O31350/Cr2763622 | 792851.551ms |
| Code round 2 | reasoning | SWE-2 Max | not reported | P1695292/O15268/C1548288 | not reported |
| Main merge, audit and close-out | — | Codex GPT-6 family | not reported | not reported | not reported |

Sonnet plan-plus-build cumulative I356/O219049/Cw440369/Cr52773764 is counted once;
the build row is its increment, not another cumulative total. Small plan-review
metadata correction and any missing interrupted portions are not separately reported.
SWE round 1 retains metrics from its first invocation and bounded completion;
those invocations are not added again. Failed AGY exit 3 is RESOURCE_EXHAUSTED.
Private exports remain outside git; actual report sources are in the run record
and coordinator metadata audits. Exact retained-field subtotal is recorded below.
Reported token subtotal 69,643,991, including cache once, plus unreported parts.
Reported elapsed subtotal 4,539,920.040ms, plus unreported parts.
