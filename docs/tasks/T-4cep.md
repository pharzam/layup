# T-4cep — admit the canonical code-owner setup completion

Issue [#384](https://github.com/pharzam/layup/issues/384), PR [#395](https://github.com/pharzam/layup/pull/395).
Serves F-0003#43/#47 through REQ-003/REQ-007; necessary child of #356.
Authors: Sonnet5.5 and Codex GPT-6 family; independent reviewer: SWE-2 Max.
Plan confirmation6217115601 sets maximum1,300 added+removed lines/26files,
cap2, goals2. Operator ruling6217324567 authorizes the narrow setup-completion
reading and coupled scope; it authorizes no target writes.
Evidence: [red/green, variants, mutants and integration](../../runs/T-4cep/test-runs.md).

## Delivered and verified

The tree check admits a canonical CODEOWNERS-only child of the recorded setup
head, including legacy setup records without an S13 hash. It verifies exactly
one parent, the whole-tree difference, regular mode100644 and exact bytes from
CodeOwners using the recorded Operator login and setup manifest at that head.
It reads no commit identity as authorization. The prepared Operator-run landing
script is enabled only after this predicate exists; no target write ran.
The first-merge follow-up remains #363, explicitly unbuilt here.
The original script's disable/trap/restore/read-back behavior is preserved.
Tests exercised wrong parent/history/path/mode/content, legacy records, real
Git objects, the script/check junction and the production forge-check row.
All23 final valid compiled mutants were killed with restored green controls.
The original invalid unused-variable mutant is not credited as a kill.
Full integration on be0a38f96e979db32ddd27697bceb8e700d0af90 passed all21packages,
15:22:49.182859–15:34:49.256158UTC, direct0,720073.299ms, including main a3a4cdd.

## Review and verdict

Round1 reviewed2563c04b5713206aabf3bff882d413ff2016ddb8, cycle0/cap2,
record6219171547; intermediate material when the next main merge needed review.
Round2 reviewedab9ebb851d37a83b9694e05423d7bb7ec1eb566b, cycle1/cap2,
record6220325988: nothing material in scope; no product fix followed.
The scoped review verified the #352 merge clause by clause, retaining both
implementations, all shared prose and index entries. The predicate is unchanged.
Round1 discloses filename exposure; round2 read only its own overflow outputs.
Piped-check limitations are explicit; root direct checks and hooks passed.
Exact complete landing: 1293 lines/25 files; maximum1,300/26 unchanged.
Final combined CI remains required: later main Go was not in the full run above.
An earlier security scan's database HTTP403 was a failure, not a credited pass.
The named reasoning tier was unavailable under the requested harness models;
model independence holds, and this tier limitation remains explicit.

## Resource record

ADR-0007 reported values: Sonnet I/O/Cw/Cr = uncached input/output/cache creation/
read; thinking is included in output. SWE P/O/C = prompt/output/cached (inside P).
Missing parts are unreported; no timeout ceiling or overlapping export is added.

| Part | Expected tier | Model | Effort | Reported tokens | Elapsed |
| --- | --- | --- | --- | --- | --- |
| Plan | reasoning | Sonnet5.5 | not reported | I90/O69078/Cw224606/Cr5860322 | 654549ms |
| Independent plan review | reasoning | SWE-2 Max | not reported | P1522611/O40627/C1338011 | not reported |
| Build/tests/docs and bounded finish | execution | Sonnet5.5 | not reported | increment I356/O151617/Cw262937/Cr71545742 | initial unreported; finish127825ms |
| Full integration | — | deterministic | — | not applicable | 720073.299ms |
| Code review including metadata correction | reasoning | SWE-2 Max | not reported | P3757434/O40629/C3650176 | not reported |
| Scoped merge review | reasoning | SWE-2 Max | not reported | P1352503/O28222/C1273344 | not reported |
| Merge/audit/close-out | — | Codex GPT-6 family | not reported | not reported | not reported |
| **Total reported subset** | | | | 84,856,774 tokens; cache counted once | 1,502,447.299ms |

Sonnet cumulative I446/O220695/Cw487543/Cr77406064 counts plan/build/finish once.
SWE plan review's earlier28step export is a subset of the final32steps, not added.
Initial author timeout ended without a result; bounded completion reports only
its own127825ms. No inferred elapsed or unreported context counter is credited.
No backlog entry existed; the completed log receives its direct entry.
The reported totals exclude unreported parts and do not claim full-session totals.
