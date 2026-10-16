# T-9et9 — the basic layup/rules judge

Issue [#352](https://github.com/pharzam/layup/issues/352), PR [#394](https://github.com/pharzam/layup/pull/394).
Serves F-0003#43 through REQ-003/NFR-004. Authors Sonnet5.5 and Codex GPT-6 family.
Plan confirmation6216617960: maximum1,600 added+removed lines/16files, cap1, goals1.
Evidence: [red/green and full integration](../../runs/T-9et9/test-runs.md),
[compiled mutation controls](../../runs/T-9et9/mutate.py).

## Delivered and verified

Basic reads the named base/head/register commits; Judge refuses changed rule
paths unless their regular head blob exactly equals code's rendering of the
named record. Render supports facts and PRD. Missing/malformed/empty registers,
non-ancestor bases and invalid pre-images refuse. Result reasons are closed.
Existing Check/ReadRegister/Match/GuardrailsAdditions and their frozen old tests
are unchanged. No copy hash is checked here; that remains the caller's duty.
These helpers are exported and unwired; later loop callers own production use.
The tests include byte differences, modes, deleted paths, register forgery,
pre-image ancestry and unusual filenames. Of38 compiled mutants,35 were killed
and3 were documented equivalent; invalid compilation is not a kill.
Full integration on f5693500d5a9d76b202a24fa54f2d7b6d89f18d7 passed21packages,
14:12:19.645234–14:26:39.749756UTC, direct0,860104.522ms. Later main changes
are covered by required final combined-SHA CI, not claimed in that earlier run.

## Review and verdict

The first SWE attempt is uncredited: its tools read an earlier unrelated review.
Qualifying round1 was fresh OpenCode DeepSeek V4.1 Flash free on40bfa5e,
record6218385008, cycle0; initial timeout, bounded completion and factual
prepublication correction were one review context, with cache counted once.
Its wording note was fixed in e098491; its defensive-precondition note is
informational, declined with the written valid-setup boundary intact.
Round2 SWE-2 Max reviewed e0984913a401deeb8abae646410c4bd9fe282c5e,
record6219271400, cycle1/cap1: not mergeable, findings recorded, budget only.
No material code defect remains. All main-merge resolutions and disjoint
incoming changes were judged, including the non-exempt PRD changelog.
Both qualifying rounds differ from every author. Round2 listed external audit
filenames but read no prior review content; that context limit is disclosed.
Its focused vet/unit/integration checks have captured direct exits of0.
Round1's piped checks report observed output, without invented direct exits.
Named reasoning-tier models were unavailable on the requested independent routes.
Round2 notes are declined: a read failure on an existing blob fails closed;
path-first ordering matches the spec; doubled table padding is cosmetic.
#390/#391 are pre-existing off-path findings; #392 is newly introduced cleanup
duplication, independently classified non-blocking. They remain separate issues.
Exact complete landing: 1701 lines/15 files; original maximum1,600/16 unchanged.
R12 requires the exact one-time approval on #352; final full-SHA CI is required.

## Resource record

ADR-0007: Sonnet I/O/Cw/Cr are uncached input/output/cache creation/read;
thinking is within output. SWE P/O/C are prompt/output/cached, cache inside P.
OpenCode reports its raw fields; inclusion/normalization is not inferred.

| Part | Expected tier | Model | Effort | Reported tokens | Elapsed |
| --- | --- | --- | --- | --- | --- |
| Plan | reasoning | Sonnet5.5 | not reported | I68/O69116/Cw201709/Cr4590433 | 585946ms |
| Plan review | reasoning | SWE-2 Max | not reported | P2043525/O35830/C1801766 | not reported |
| Build/tests/docs | execution | Sonnet5.5 | not reported | increment I196/O101557/Cw172629/Cr31070610 | 909884ms +2131ms final notification |
| Full integration | — | deterministic | — | not applicable | 860104.522ms |
| Uncredited review attempt | reasoning | SWE-2 Max | not reported | P3394332/O53344/C3283840 | not reported |
| Qualifying round1, all continuations | reasoning | DeepSeek V4.1 Flash free | not reported | raw I347317/O14615/R30029/Cr4289152/Cw0 | not reported |
| Round2 | reasoning | SWE-2 Max | not reported | P1371902/O30494/C1280365 | not reported |
| Merge/audit/close-out | — | Codex GPT-6 family | not reported | not reported | not reported |
| **Total reported subset** | | | | Sonnet 36206318; SWE 6929427; OpenCode raw vector above | 2358065.522ms |

Sonnet cumulative I264/O170673/Cw374338/Cr35661043 includes plan/build once;
the plan-review initial export is a subset of its final32 unique metric rows.
OpenCode's62 unique step-finish IDs cover timeout/completion/correction once.
Reported totals count cache once where defined; the unnormalized OpenCode
vector is not added to another model's scalar. Unreported parts remain missing.
