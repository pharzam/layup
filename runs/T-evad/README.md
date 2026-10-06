# T-evad: the first pilot

The first pilot of LAYUP (O-122): `layup` set up one target repository from a Go
problem statement, and ran the target's gate from outside. Task `T-evad`
([#97](https://github.com/pharzam/layup/issues/97)), row 20 of the
[implementation plan](../../docs/plan/README.md). The verdict and the resource
record are in [`docs/tasks/T-evad.md`](../../docs/tasks/T-evad.md).

**The verdict of the first pilot is a Fail** (O-145 = a): at T5 the idea owner
rejected `REQ-002` and `NFR-003`, as the setup lost markers. LAYUP was fixed in this
task, and a second run set the target up again. The result of the second run is its
own result, in [`acceptance.md`](acceptance.md).

## The target

| Item | Value |
| ---- | ----- |
| Repository | [`pharzam/chat-orchestrator`](https://github.com/pharzam/chat-orchestrator), public |
| Problem statement | `PSB-CHAT-001`, revision 3.0, SHA-256 `97c657b6…`: the idea owner's choice (Decision Point 1) |
| Baseline | `pharzam/armature` at `a959655`, its newest commit (O-140 6a), in both runs |
| First run | root `242a205` (the unchanged baseline); setup head `cec749a`, the Operator approved its prose; `layup-records` `2a339bb`; the binary built from `6206335` (the fixes of F-8) |
| Ruleset | 24509051, `layup: the default branch`, active |
| `main` after D7 | `7d7b387`, the merge of the target's pull request 5 |
| Second run | root `1ced83a` (the tree of the baseline, `8ffb250`); setup head `e2b402b`, on GitHub as `layup-setup-2`; its records `8bc105c`, as `layup-records-2`; the binary built from `006fbee` (fixes 1 to 3) |

## The first run, in the order of `setup.md`

| Step | Result | Who |
| ---- | ------ | --- |
| S01 to S06 | done; the five answers of S01 came from the target's issue `Setup answers` (T1) | the author ran it; the Operator answered |
| The prose step | a stop for 25 files (S07, S08, S09 and 22 of S14); the author wrote them | the author; the Operator reviewed (3 corrections) and approved them |
| S10 | a stop with 107 questions (T2) | the Operator approved the answers, after 6 corrections, and posted them on the target |
| S10 to S14, then S15 | S15 stopped for `verify`; `layup setup verify`: exit 0, 15 `pass` and 3 `clear`; then S15 done, with S13 handed to the Operator | the author |
| T3 | [`t3.sh`](t3.sh) ran the five commands of `commands.sh`, fail-fast, and read back `main`, `layup-records` and the ruleset | the Operator |
| D7 | the task `T-a0rt` of the target, by the target's own process: an issue, a plan and its review, two review rounds, the hooks, a pull request with 13 checks, a merge under the ruleset; `layup gate` from outside on its base and head, twice, the same bytes | the author; the Operator posted and merged |
| D8 (a) | 66 of 66 values supported, and both planted rows caught ([`value-audit.md`](value-audit.md)) | GPT-6 Sol |
| D8 (b) | [`rules-diff.sh`](rules-diff.sh) on the setup head: PASS | the author |
| T5 | the Operator rejected V-040 and V-061 of D8 (a), and `REQ-002` and `NFR-003`, and asked for a fix ([`acceptance.md`](acceptance.md)) | the Operator |

The gate run of D7 (`layup gate`, base `cec749a`, head `ea28ed4`, before the merge):

```text
kind      state    result  reason
static    active   clear   no product path
layout    pending  clear   pending: no product path
boundary  pending  clear   pending: no product path
contract  pending  clear   pending: no product path
test      active   clear   no product path
```

The change had no Go file, so the run does not show an active gate on product code.
`verify` ran the active kinds on a clean run and on a known-bad fixture, which failed
as it must (`gate:static`, `gate:test`: `pass`).

`rules-diff.sh` on the setup head, with the record of `layup-records`: of the 231
baseline paths, 171 are byte-identical, 29 changed (20 adapted, 5 written for the
target, 16 with marker rows of S11, 2 task indexes, 1 index with added rows), 31
deleted by S05; 13 paths are new; 10 setup commits.

## The fix and the second run

| Step | Result | Who |
| ---- | ------ | --- |
| The fix | three fixes of LAYUP, test first: one rule for a marker, also over more lines (`c68b390`); S14 refuses an input that loses a place of a marker (`c68b390`, `cc31402`, `006fbee`); check `markers` names a quote with no pair (`c68b390`) | the author; the design by the Operator (O-145, O-146) |
| Rehearsals | A: S14 refuses the first inputs (27 places lost); B: the corrected inputs pass, and the seven places of T5 are right ([`test-runs.md`](test-runs.md)) | the author |
| The inputs | the prose of 8 files keeps the 27 places; 14 new questions of S10 | the Operator approved the prose and posted the 14 answers on the target |
| S01 to S15 | a new work area, the same baseline commit; `verify`: 15 `pass` and 3 `clear` | the author |
| D8 (b) | `rules-diff.sh` on `e2b402b`: PASS | the author |
| D8 (a) | 67 of 67 changed values supported, after a second inventory (F-29) ([`value-audit.md`](value-audit.md)) | GPT-6 Sol |
| T6 | `t6.sh` pushed `layup-setup-2` and `layup-records-2`, fail-fast, and read them back | the Operator |

`rules-diff.sh` on the second setup head: of the 231 baseline paths, 171 are
byte-identical, 29 changed (20 adapted, 5 written, 17 with marker rows of S11, 2
task indexes, 1 index with added rows), 31 deleted, 13 new; 10 setup commits.

## The files here

| File | What |
| ---- | ---- |
| [`findings.md`](findings.md), [`numbers.md`](numbers.md) | The findings F-1 to F-30, and the numbers of D10 |
| [`test-runs.md`](test-runs.md) | The red and green runs, and the ladder |
| [`value-audit.md`](value-audit.md), [`acceptance.md`](acceptance.md) | D8 (a) of both runs, and T5 |
| [`rules-diff.sh`](rules-diff.sh), [`rules-diff-test.sh`](rules-diff-test.sh) | The check of D8 (b), and its 20 fixture cases |
| [`t3.sh`](t3.sh) | The fail-fast procedure of T3 (F-13) |
| [`tree-equal.sh`](tree-equal.sh), [`tree-equal-test.sh`](tree-equal-test.sh) | The check of condition 2 of O-146, and its 5 cases |

The evidence of the runs is outside this repository: the target's Git (the prose,
the diffs, the answer records `F-0001` and `F-0002`, `layup-records` and
`layup-records-2`), the target's issues and pull requests, the comments of #97, and
the logs, the tools and their tests on the host of the pilot (`evidence/`,
`tools/`). The tools name the paths of that host.
