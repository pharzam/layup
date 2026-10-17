# T-4wyg — the full specification check

Issue [#351](https://github.com/pharzam/layup/issues/351), PR [#393](https://github.com/pharzam/layup/pull/393).
Serves F-0003#51/#58 through REQ-012/NFR-004. Authors: SWE-2 Max and
Codex GPT-6 family; independent plan/code reviewer: fresh Sonnet5.5 sessions.
Plan confirmation6215968646 sets maximum2,100 changed lines/20files and goals2.
Operator ruling6220109258 raises only the cycle cap to2 for a tests-only repair;
current-main measurement follows #3556209446202. No standing maximum changed.
Evidence: [all runs, repairs and compiled mutants](../../runs/T-4wyg/test-runs.md).

## Delivered and verified

`RunPRDLint` reads both streams, accepts exit0 and the source-bound empty
MoSCoW/Phase refusal, and refuses other failures. `CheckFull` always returns
all ten rule rows; it binds the numbered facts, answer quotes, PRD, inventory,
rendered copy and task register, with clear/no-run behavior as the frozen spec
states. The functions remain caller APIs; no new full-check CLI is claimed.
The `task-spec-check` block and second `tasks` value are owned by internal/spec.
Traceability, PRD cells, package ownership and known pitfalls accompany the code.
Short malformed criteria rows refuse without panic. Missing version bytes make
`ran` fail. Measures keep their out-of-scope/waits exemptions. The final repair
adds exact/case/level/boundary/token heading tests, latest-version selection in
both index orders, the empty-MoSCoW/set-Phase case and the PRD-path mode case.
All ten targeted compiled mutants fail assertions, with original source restored;
this cap2 repair changed tests/evidence, no #351 production code.
Full integration before the proposed test-only fix on3f406f80e20ac3f4ad0bef6f35b2f87f0c5d4094 passed21packages,
17:11:54.482887–17:25:34.613722UTC, direct0,820130.835ms, including main d09b1c0.
All main merges preserve both sides; the final changelog keeps #384 then #351.

## Review and verdict

Round1 record6217333827 found four material defects; the author repaired them.
Round2 record6218710801 found missing falsifying contract tests and a budget
condition. Its table is intermediate material; raw findings remain unchanged.
The Operator authorized the bounded additional cycle; the author answered6220274712
before changes. Round3 reviewed0e823863a2c469a68df280db83c4602040cf8280, cycle2/cap2,
record6221306780: `not mergeable, findings recorded`. Its material findings were
F1 (leading-text heading near miss and Markdown levels4–6 falsifiers) and
F2 (bare answers-record acceptance falsifier). The approved bounded fix adds
those tests, updates their traceability/PRD cells, and corrects one comment.
All three corresponding compiled mutants fail assertions; restored source is green.
Production logic is unchanged. The SWE proposal used32 reported metric rows.
This is a fix that no further round read. The last review verdict remains
`not mergeable, findings recorded`; Operator ruling6232930664 permits this
exact ending and approves the complete2,370-line/18-file landing. This record
identifies the final fix as unreviewed; final exact-head CI is required.
Complete current-main landing: 2370 changed lines/18 files.
This exceeds the unchanged2,100-line maximum; one-time R12 approval is recorded.
Final combined-head CI is required before merge, with evidence on the issue.
Notes N1/N2 from round3 are included in the ordinary fix, before close-out.
Earlier N3/N4 document expansion was declined within the tests-only scope;
N5's stale recap is replaced here; N6 test names are now in traceability;
N7 exit0 with stderr follows the existing contract.
The requested models could not reach the named reasoning tier; independence and
that limitation are explicit in every review. Tier availability is not overstated.

## Resource record

ADR-0007 reported subset only. SWE P/O/C = prompt/output/cached (inside P).
Sonnet I/O/Cw/Cr = uncached input/output/cache creation/read; thinking inside O.

| Part | Expected tier | Model | Effort | Reported tokens | Elapsed |
| --- | --- | --- | --- | --- | --- |
| Plan | reasoning | SWE-2 Max | not reported | P3186764/O50035/C3015704 | not reported |
| Independent plan review | reasoning | Sonnet5.5 | not reported | I66/O37166/Cw143918/Cr3137456 | 365791ms |
| Initial build/test/docs | execution | SWE-2 Max | not reported | P6847466/O38553/C6549504 | not reported |
| First material-fix cycle | execution | SWE-2 Max | not reported | P2318995/O10270/C2266176 | not reported |
| Cap2 tests-only repair, unique increment | execution | SWE-2 Max | not reported | P3812309/O38163/C3695655 | not reported |
| Code review1, including metadata correction | reasoning | Sonnet5.5 | not reported | I50/O30880/Cw108210/Cr1908949 | 296349ms |
| Code review2 | reasoning | Sonnet5.5 | not reported | I44/O22750/Cw104040/Cr1620001 | 301764ms |
| Code review3 | reasoning | Sonnet5.5 | not reported | I48/O19064/Cw80511/Cr1549852 | 294850ms |
| Approved bounded final fix | execution | SWE-2 Max | not reported | P1651898/O21888/C1600000 | not reported |
| Full integration, three recorded runs | — | deterministic | — | not applicable | 2480324.539ms |
| Producing audits/merge/close-out | reasoning/execution | Codex GPT-6 family | not reported | not reported | not reported |
| **Total reported subset** | | | | 26739346 tokens; cached counted once | 3739078.539ms |

SWE union202 unique metric rows: P17817432/O158909/C17127039. Cap2 export includes
30 earlier-fix rows; its full counter is not added again. Earlier pre-compaction
build history is unrecovered. No missing time, timeout ceiling or cache is inferred.
Totals exclude unreported parts; they are not claimed as whole-session totals.
