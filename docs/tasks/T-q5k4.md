# T-q5k4 — the handoff in the loop

Issue [#360](https://github.com/pharzam/layup/issues/360), PR [#389](https://github.com/pharzam/layup/pull/389).
Serves F-0003#45/#59 through REQ-005 and NFR-001; resolves the in-path #377.
Authors: Claude Sonnet5.5 and Codex GPT-6 family; reviewers: SWE-2 Max.
Plan confirmation6215216194: maximum2,800 changed lines /27 files, cap1, goals2.
Evidence: [red/green, compiled mutants and integration](../../runs/T-q5k4/test-runs.md).

## Delivered and verified

Session comment files are checked against result rows, kept byte for byte and
posted from those bytes. Comment copying and handoff events bind the recorded
hashes, copy and live body. Invalid artifacts, edited comments and lost leases
refuse. Restarts adopt the App's own already-posted comment and copy it once,
including when the post succeeds and its copy commit is lost.
Question labels route to the owner role or the Operator; the new answer result
kind and the next attempt's recorded citation determine answer acceptance.
Caller-supplied program items remain caller-owned; #387 tracks the separate gap.
These helpers are unwired until the later loop driver. No end-to-end loop claim.

The production constructor regression first reproduced the duplicate post,
then passed with Sessions.App set from Register.AppSlug. A valid compiled mutant
removing that assignment fails the regression; restored controls pass.
Original full integration passed onabd7eadee4ae47a954da8023a7f00da7fced6e18,
13:02:01.654318–13:19:21.769262UTC, direct0,21packages,1040114.944ms.
Repaired full integration passed on7e10a292c9f20cfecff6183b4a0209ff02cdd1ab,
14:00:52.928001–14:12:13.012882UTC, direct0,21packages,680084.881ms.
Main f86f71c followed that run; final combined full-SHA CI remains required.
The test hash helper collision was observed by vet and fixed only by renaming
the handoff helper and its callers; the assertions and product behavior stay.

## Review and verdict

Round1 reviewed6dd598ae708d3b6b7b4489bb1a7af2805945491f atcycle0/cap1,
record6216956461, intermediate material. The sole material code finding was
missing App identity in the production Sessions constructor, fixed test-first.
Round2 reviewede4abfcf6c36e1b3b5293658fa6cc520541784563 atcycle1/cap1,
record6217822863: not mergeable, findings recorded, solely exact budget and CI
conditions. No material code finding remains; all resolutions are verified.
The five main-merge resolutions and the test-helper rename belong to round2:
guardrails both lesson sets; PRD both REQ005/NFR001 citations; loop both sets of
artifact/answer and transition replacement tests; packages both ownership sets;
traceability both test sets. No broad shared-index exception is claimed.
The six original notes are declined without source fixes: the documented dedup
key is session/detail even for a different pair; refusal detail stays in the
returned error while records require only the reason; the internal API relies
on upstream validated file names; code-written session identity is caller-owned;
two event reads are consistent under one held lease; the measured comment cap
is a conservative policy. These limits remain explicit, not extra passes.
Round2 notes are declined: one result end per session makes a later end
unreachable; answer association follows the specified positional rule; older
struct-literal tests complement the production-constructor regression; the
one-call lease window is documented. #387 was opened before merge; the PR body
closes #377. Context, method, execution and model independence held in both rounds.
Requested harness models limit access to the named reasoning tier; this is recorded.
No target was changed. Exact landing:2,707 lines/29 files; R12 thread6217886549.
The original maximum remains2,800/27; exact Operator approval and final CI are required.

## Resource record

ADR-0007: reported, not budgeted. Sonnet I/O/Cw/Cr = uncached input/output/
cache creation/cache read; thinking is included in output. SWE P/O/C = prompt/
output/cached, with cached inside prompt. Missing parts remain unreported.

| Part | Expected tier | Model | Effort | Reported tokens | Elapsed |
| ---- | ------------- | ----- | ------ | --------------- | ------- |
| Plan | reasoning | Sonnet5.5 | high | I120/O64449/Cw224153/Cr7495702 | 594914ms |
| Independent plan review | reasoning | SWE-2 Max | not reported | P3701351/O47906/C3543055 | not reported |
| Build/tests/docs | execution | Sonnet5.5 | high | increment I244/O149625/Cw231115/Cr45235778 | 2815648ms |
| Original full integration | — | not applicable | not applicable | not applicable | 1040114.944ms |
| Code round1 retained metrics | reasoning | SWE-2 Max | not reported | P1426076/O18044/C1301548 | not reported |
| Constructor fix | execution | Sonnet5.5 | high | increment I24/O8194/Cw16821/Cr5730253 | 278052ms +177ms acknowledgment |
| Repaired full integration | — | not applicable | not applicable | not applicable | 680084.881ms |
| Code round2 | reasoning | SWE-2 Max | not reported | P5294026/O50282/C5163648 | not reported |
| Merge/audit/close-out | — | Codex GPT-6 family | not reported | not reported | not reported |

Sonnet cumulative I388/O222268/Cw472089/Cr58461733 counts plan/build/fix once.
The fix acknowledgment had unchanged counters. Round1's completion export is a
subset of the corrected export,17 unique metric rows; earlier lost review history
is unavailable. No context-remaining counter or timeout bound is credited.
Private exports and metadata audits remain outside git; source names are in
run evidence. Reported retained subtotal69,694,163 tokens, including cache once,
and reported elapsed subtotal5,408,990.825ms; other parts stay unreported.
