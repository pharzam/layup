# T-v2sv — the Retrospective

Issue [#359](https://github.com/pharzam/layup/issues/359), PR [#386](https://github.com/pharzam/layup/pull/386); settles #372 under R10.
Serves F-0003#51 through REQ-012; goal count1, cycle cap1.
Author: Devin SWE-2 Max; coordinator: Codex GPT-6 family.
Plan6214476246, independent confirmation6214608102, answer6214713954.
Maximum700 added-plus-removed lines /15 files, including close-out and fixes.
Base main e16568c; frozen head9d77920ba1ae25df160be5b3e91e3400a49dfdd1.
Evidence: [test and mutation runs](../../runs/T-v2sv/test-runs.md).

## Delivered and verified

The helper appends milestone/time to retrospectives.tsv and commits the event
through Fenced once. Existing M2d/M2e rows survive. A repeated milestone returns
done-before without a write; malformed, empty or dashed owner data refuses.
It starts no session, spends no budget and returns the bet step.
This is an unwired helper: no actual loop driver or caller is delivered here.
The records schema and tasks Go schema equal the implemented shared block;
#372's missing block, source owner and inventory agreement are resolved.
Loop, records, packages, traceability and PRD test cells are synchronized.

Nine unit tests were observed failing against stubs; two refusal guards were
already green on arrival and are not claimed as red demonstrations. Eleven
compiled, assertion-killed mutants cover writing, preservation, idempotence,
next step, clock, lease fencing, no session, schema, order and owner validation.
The Truncate-only mutation is observationally equivalent under second-precision
formatting and is not credited. Both real-git/schema integration tests pass.
Full module integration at Go-identical2566c75: wrapper
2026-10-16T11:59:07.405138Z–12:12:47.517735Z, direct exit0,21packagesok.
Later evidence/guardrail changes do not modify Go. Final local checks and exact
landing-head CI remain required before plain merge.

## Review and verdict

| Round | Head | Cycle/cap | Material | Verdict |
| ----- | ---- | --------- | -------- | ------- |
| 1 | 9d77920 | 0/1 | 0 | nothing material in scope |

Fresh independent Sonnet5.5 review6215666115 held context, method, execution
and model independence from both authors. Its six notes are declined without
source fixes: N1 clock wiring is directly helper-tested, driver remains unwired;
N2 a historical timestamp typo remains, not credited as a precise time;
N3 resource prose spacing is readable and private audit paths are not public
artifacts; N4 complete resources follow below; N5 control-character normalization
requires a noncanonical milestone and no driver supplies one; N6 helper-only
scope is explicit above. Notes do not cause a new round. Only this close-out
and the completed/backlog indexes follow the reviewed head.
Delivered and reviewed with no material findings. No target changes performed.
Main82bb082 merged without conflicts; changes are disjoint or shared indexes.
The completion entry uses the parser-required bullet; final size exactly668/15.

## Resource record

ADR-0007: reported usage, not budgets. SWE P/O/C means prompt/output/cached;
cached is included in prompt. Sonnet I/O/Cw/Cr means uncached input/output/cache
creation/cache read. Thinking is already included in output. Missing figures
remain not reported; interrupted invocation limits are not measured duration.

| Part | Expected tier | Model | Effort | Reported tokens | Elapsed |
| ---- | ------------- | ----- | ------ | --------------- | ------- |
| Plan, including resumed timeout | reasoning | SWE-2 Max | not reported | P2093066/O35128/C1901893 | not reported |
| Independent plan review | reasoning | Sonnet5.5 | high | I54/O33428/Cw123215/Cr2211811 | 341527ms |
| Build, tests, mutations and docs | execution | SWE-2 Max | not reported | increment P8841569/O60825/C8669336 | not reported |
| Full integration | — | not applicable | not applicable | not applicable | 820112.597ms |
| Independent code round1 | reasoning | Sonnet5.5 | high | I38/O17292/Cw72479/Cr1087529 | 363010ms |
| Evidence, guardrails and close-out | — | Codex GPT-6 family | not reported | not reported | not reported |

Combined SWE conversation reports P10934635/O95953/C10571229 over86steps;
plan26steps plus build increment60. Do not add cached tokens again. Sonnet totals
I92/O50720/Cw195694/Cr3299340. Reported token subtotal14576434, plus unknown
portions; reported elapsed subtotal1524649.597ms, plus unknown portions.
The first planning invocation exited124; same-conversation resumes are retained,
not represented as independent sessions. Private full exports stay outside git;
actual metadata is documented in runs, and the export lesson is in guardrails.
The model-tier limitation for Sonnet reviews is recorded; reviewer independence
from the models authoring this change holds.
