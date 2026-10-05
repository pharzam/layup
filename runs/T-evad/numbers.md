# The numbers of the first pilot

Task `T-evad` ([#97](https://github.com/pharzam/layup/issues/97)). D10 of the plan:
the defects with their class, the stops, the human inputs and the time of each step,
for the ADR that ends bootstrap mode (ADR-0012 part 6). That ADR counts the stalls of
LAYUP's own phase-1 tasks as its first step (`gov-pilot-numbers`); this file gives
the numbers of the pilot only. The findings are in [`findings.md`](findings.md).

## The defects, by class

| Class | Count | Findings |
| ----- | ----- | -------- |
| A defect of LAYUP's phase-1 code that blocked the pilot; fixed here, test first (one Investigate each) | 2 | F-8: S11 (row 13), check `sources` (row 10) |
| A defect or limit of LAYUP's phase-1 code that did not block; not fixed | 6 | F-1 (row 14, a known limitation), F-3 (row 13), F-6 (row 15), F-13 (row 15), F-16 (row 13), F-19 (row 15) |
| A defect of this task's own check; fixed here, test first | 1 | F-18 |
| A defect of the baseline that LAYUP copies | 1 | F-15 (bootstrap mode suspends the rule in LAYUP) |
| An error in the author's inputs, found by the Operator | 9 | F-5 (4), F-9, F-10, F-12 (3) |
| An error in the author's inputs, found by a check | 1 | F-7 (check `glossary` of S08) |
| An error in the author's brief to a reviewer | 1 | F-14 |
| A material finding of a review round of the target's own task (D7) | 1 | Round 1 of `T-a0rt`: the assertion read words, not the cells; fixed in cycle 1 |

No check that did not run counted as a pass. Of the 22 defects and errors above,
the Operator found 10 (F-1 and the 9 input errors), a check or the step runner 3
(F-7 and the two of F-8), a reviewer 2 (F-16, in the plan review of `T-a0rt`, and
the finding of its round 1), and the author 7 (F-3, F-6, F-13, F-14, F-15, F-18
and F-19).

## The stops

| Stop | Count | Notes |
| ---- | ----- | ----- |
| Planned input stops of one work area | 4 | The prose step (25 files: S07, S08, S09 and 22 of S14), S10 (107 questions), S15 (`O-verify`), and S13 handed to the Operator (T3) |
| Work areas of the real run | 5 | Run 1, then runs 2 to 5 after the T2 corrections, T2 revision 2 and the prose review; LAYUP refuses a changed input of a done step (F-2) |
| Rehearsals of S10 to S15 | 3 | The first found F-8 |
| A step that failed on an input error | 1 | S08, F-7 |
| Investigate events (the pilot rule) | 2 | S11 and check `sources`, one each (F-8) |
| Decisions asked of the Operator | 6 | O-139 to O-144 |
| A reviewer harness that failed | 1 | Devin with a new home (F-17); the second try worked |

## The human inputs

| Input | Count | Where |
| ----- | ----- | ----- |
| The Operator's comments on #97 | 8 | 5979471031, 5993893238, 5995076794, 5995693012, 5996154495, 5996579611, 5996877039, 5998876808; T5 is the ninth |
| The Operator's posts on the target | 2 | The S01 answers (5989769582) and the S10 answers (5995217834) of `pharzam/chat-orchestrator#1` |
| Decisions given in the session | 2 | O-141 ("29 a") and T2 revision 2 |
| Corrections asked before an approval | 10 | T2: 4, then 2; the prose: 3, and 1 known limitation |
| Scripts that the Operator ran | 4 | `t3.sh`, `t4-1.sh`, `t4-2.sh`, `t4-3.sh`; each passed its checks and exited 0 |
| Other actions of the Operator | 2 | The choice of the brief (Decision Point 1), and the issue `Setup answers`; the author made the target repository with the Operator's login, at the Operator's request in the session |

## The review rounds and the audits

| Round | Reviewer | Result |
| ----- | -------- | ------ |
| The plan review of this task | GPT-6 Sol on the Devin CLI | `approve-with-conditions`: 8 conditions, applied |
| The plan review of `T-a0rt` (D7) | GPT-6 Sol on the Devin CLI | `approve-with-conditions`: 120 lines over 4 files, cap 2 |
| Round 1 of `T-a0rt` | GPT-6 Sol on the Devin CLI | `material` (one finding), by O-144 |
| Round 2 of `T-a0rt` | GPT-6 Sol on the Devin CLI | `nothing material in scope` |
| D8 (a), pass 1 | GPT-6 Sol on the Devin CLI | 56 supported, the 2 planted rows caught, 8 cannot check |
| D8 (a), pass 2 | GPT-6 Sol on the Devin CLI | the 10 rows supported |

## The time of each step

The times are UTC, from the times of the comments and of the runs. "Work" is the
author or a reviewer; "wait" is a turn of the Operator. The author also did some
work during a wait (for example, `t4-3.sh` while T4-2 waited).

| From | To | Kind | Elapsed | Step |
| ---- | -- | ---- | ------- | ---- |
| 10-03 16:15 | 10-03 16:17 | work | 2 min 18 s | the question of the inputs |
| 10-03 16:17 | 10-04 11:30 | wait | 19 h 12 min | O-139 and O-140 |
| 10-04 11:30 | 10-04 12:05 | work | 34 min 45 s | the plan, its review (5 min 29 s), the answer, the T1 request |
| 10-04 12:05 | 10-05 07:07 | wait | 19 h 02 min | T1: the target repository, O-141, the issue `Setup answers`, the S01 answers |
| 10-05 07:07 | 10-05 07:47 | work | 39 min 52 s | run 1: S01 to S06; the prose step (4 helper agents); S07 to S14; the T2 request |
| 10-05 07:47 | 10-05 08:13 | work | 25 min 24 s | `rules-diff.sh`, red then green; the rehearsal (F-8: two fixes) |
| 10-05 08:13 | 10-05 11:53 | wait | 3 h 40 min | T2: the review of the table of S10 |
| 10-05 11:53 | 10-05 12:29 | work | 36 min 09 s | the four corrections of T2; runs 2 and 3; a rehearsal |
| 10-05 12:29 | 10-05 12:47 | work | 17 min 36 s | T2 revision 2; run 4; a rehearsal |
| 10-05 12:47 | 10-05 13:15 | wait | 28 min 12 s | the verification of T2, and its post on the target |
| 10-05 13:15 | 10-05 13:17 | work | 1 min 39 s | S10 to S15, `verify`, `rules-diff.sh`, the claim checks |
| 10-05 13:17 | 10-05 13:42 | wait | 25 min 15 s | the prose review |
| 10-05 13:42 | 10-05 13:48 | work | 5 min 49 s | three prose corrections and one limit; run 5 |
| 10-05 13:48 | 10-05 14:08 | wait | 19 min 49 s | the prose approval |
| 10-05 14:08 | 10-05 14:13 | work | 5 min 12 s | `t3.sh` and its tests |
| 10-05 14:13 | 10-05 14:21 | wait | 7 min 50 s | T3 |
| 10-05 14:21 | 10-05 14:24 | work | 3 min 01 s | the read-back of T3; the question O-142 |
| 10-05 14:24 | 10-05 14:31 | wait | 6 min 51 s | O-142 |
| 10-05 14:31 | 10-05 14:45 | work | 13 min 58 s | D7: the plan and the plan review of `T-a0rt`; the question O-143 |
| 10-05 14:45 | 10-05 14:47 | wait | 2 min 23 s | O-143 |
| 10-05 14:47 | 10-05 14:50 | work | 3 min 11 s | `t4-1.sh` and its tests |
| 10-05 14:50 | 10-05 14:51 | wait | 1 min 09 s | T4-1 |
| 10-05 14:51 | 10-05 15:02 | work | 11 min 05 s | `T-a0rt`: cycle 0, round 1, the fix; the question O-144 |
| 10-05 15:02 | 10-05 16:00 | work | 57 min 39 s | `T-a0rt`: round 2, the close-out, `layup gate`, `t4-2.sh` and its tests |
| 10-05 16:00 | 10-05 17:19 | wait | 1 h 18 min | O-144, then T4-2 |
| 10-05 17:19 | 10-05 17:21 | work | 2 min 09 s | the checks of the pull request; `t4-3.sh` |
| 10-05 17:21 | 10-05 18:00 | wait | 39 min 29 s | T4-3 |
| 10-05 18:00 | 10-05 18:02 | work | 1 min 21 s | the read-back of the merge and of CI |
| 10-05 18:02 | 10-05 18:31 | work | 29 min 07 s | D8 (a): the inventory, two audit passes, the T5 request |

From the question to the T5 request: 50 h 15 min. Of it, about 4 h 50 min was work
and about 45 h 25 min was a wait for the Operator, with two nights.
