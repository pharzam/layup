# Real-head replays of `freeze-check.sh` (S3, T-b6fn #440)

Run on 2026-10-20 in disposable clones under `/home/layup/.cache/layup-session/tmp-green/t-b6fn-replay/` (cloned from `/home/layup/projects/layup` and removed after the runs). Each head and base was read from the forge review-record comments (`gh issue view N -R pharzam/layup --comments`) and the task records (`docs/tasks/T-cav3.md`, `runs/T-jy6f/test-runs.md`, `docs/tasks/T-r7b9.md`). All four commits are present in the repository history.

##Summary table

| Issue / task | Head (full SHA) | Base (full SHA) | Defect named in the review | Caught by | Named FAIL / ok line |
| --- | --- | --- | --- | --- | --- |
| #412 (`T-cav3`), round 2 (comment 6249923251) | `97a993b32b2d3bbf39d7f225c03bccd0029f2d84` | `b5f83fc53e46a328fe204152a599a8d13bbcbddf` | F1: unquoted `-run` selector at `runs/T-cav3/test-runs.md:92` | **Yes (K3)** | `FAIL K3: runs/T-cav3/test-runs.md:92: quote the selector with single quotes` |
| #415 (`T-jy6f`), round 2 (comment 6268303209) | `fa7a476e07815ce43004170a84b53ab3c55ee100` | `5a290c9b6fd85d348f00dea233dea297e8a399ab` | F3: fix-cycle mutant rows without restore hashes (`runs/T-jy6f/test-runs.md:396-398`); F1: file budget overrun | **Yes (K4, K5)** | `FAIL K4: runs/T-jy6f/test-runs.md:396: SM29 has no restore evidence` (`:397: SM30`, `:398: SM31`); `FAIL K5: -:0: projection: figure 2977/28, maximum 3100/25, margin 123/-3 ...` |
| #415 (`T-jy6f`), round 3 after fix `57fec7f` (comment 6269963495) | `38b5176c633493aec76e11763193c622883c4520` | `5a290c9b6fd85d348f00dea233dea297e8a399ab` | F3 settled: separate hash table added at `runs/T-jy6f/test-runs.md:466-472` | **Yes (K4 clears)** | `skip K4: runs/T-jy6f/test-runs.md:320: SM19a (not run)` then `ok K4: restore evidence of mutants` |
| #426 (`T-r7b9`), round 2 (comment 6289957559) | `9bb845e7252186388ee63feb38547db86de2ba5c` | `0246d9c72ecae7c3509a7522990cbf39706f3fd4` | F7: two new tests missing from `docs/tests/traceability.md`; F8: `new blank line at EOF` at `runs/T-r7b9/test-runs.md:293` | **Yes (K2, K1)** | `FAIL K2: internal/run/questionloop_integration_test.go:163: TestAVerifierThatAsksIsAnsweredAtTheHeadItWasBoundTo is not in docs/tests/traceability.md`, `:203: TestTheQuestionOfAVerifierAtItsBoundHeadLeavesTheFindingOpenWhileItWaits ...`; `FAIL K1: runs/T-r7b9/test-runs.md:293: new blank line at EOF.` |

## 1. #412 (`T-cav3`) round-2 head

- **Head:** `97a993b32b2d3bbf39d7f225c03bccd0029f2d84` (from `Commit reviewed` in comment 6249923251).
- **Base:** `b5f83fc53e46a328fe204152a599a8d13bbcbddf` (comment 6249923251).
- **Command:** `sh /home/layup/projects/layup/.worktree/T-b6fn/docs/tests/freeze-check.sh --task T-cav3 --base b5f83fc53e46a328fe204152a599a8d13bbcbddf --budget 1100/17`
- **Start / end / exit:** `2026-10-20T14:36:49Z` to `2026-10-20T14:36:49Z`, `exit=1`.
- **Output:**

```
ok K1: white space
FAIL K2: internal/run/specinput_integration_test.go:40: TestTheTreeFieldsOfTheFullInputAreTheRegularFilesOfTheHead is not in docs/prd/PRD-0001-layup.md §12 REQ-012 Test cell
... (36 K2 lines where traceability.md lists REQ-012, REQ-007, NFR-004 on the shared row while PRD §12 distributes the tests across the three cells)
FAIL K3: runs/T-cav3/test-runs.md:92: quote the selector with single quotes
ok K4: restore evidence of mutants
ok K5: projection: figure 1100/15, maximum 1100/17, margin 0/2 (the figure at landing is still measured at landing)
FAIL K6: runs/T-cav3/conditions.md:0: missing file
```

- **Assessment:** K3 reports `FAIL K3: runs/T-cav3/test-runs.md:92: quote the selector with single quotes`, which is the exact defect named in round 2 finding 1 of #412 (`docs/tasks/T-cav3.md:47`).

## 2. #415 (`T-jy6f`) round-2 head `fa7a476` and round-3 head `38b5176`

- **Defective head:** `fa7a476e07815ce43004170a84b53ab3c55ee100` (`runs/T-jy6f/test-runs.md:421`).
- **Later head holding the restore evidence:** `38b5176c633493aec76e11763193c622883c4520` (round 3 `Commit reviewed`, comment 6269963495, after commit `57fec7fc804f1cf215abdf8cd737041c4a9c4842` added the hash table at `runs/T-jy6f/test-runs.md:466-472`).
- **Base:** `5a290c9b6fd85d348f00dea233dea297e8a399ab` (comments 6268303209, 6269963495).
- **Command:** `sh /home/layup/projects/layup/.worktree/T-b6fn/docs/tests/freeze-check.sh --task T-jy6f --base 5a290c9b6fd85d348f00dea233dea297e8a399ab --budget 3100/25`
- **On `fa7a476e07815ce43004170a84b53ab3c55ee100`** (`2026-10-20T14:36:49Z` to `2026-10-20T14:36:50Z`, `exit=1`):

```
ok K1: white space
... (K2 lines for tests omitted from traceability.md or PRD §12 cells, and K3 on line 297's prose placeholder '<tests>')
skip K4: runs/T-jy6f/test-runs.md:320: SM19a (not run)
FAIL K4: runs/T-jy6f/test-runs.md:396: SM29 has no restore evidence
FAIL K4: runs/T-jy6f/test-runs.md:397: SM30 has no restore evidence
FAIL K4: runs/T-jy6f/test-runs.md:398: SM31 has no restore evidence
FAIL K5: -:0: projection: figure 2977/28, maximum 3100/25, margin 123/-3 (the figure at landing is still measured at landing)
FAIL K6: runs/T-jy6f/conditions.md:0: missing file
```

- **On `38b5176c633493aec76e11763193c622883c4520`** (`2026-10-20T14:36:50Z` to `2026-10-20T14:36:51Z`, `exit=1`):

```
ok K1: white space
... (K2 and K3 lines as above)
skip K4: runs/T-jy6f/test-runs.md:320: SM19a (not run)
ok K4: restore evidence of mutants
FAIL K5: -:0: projection: figure 3068/28, maximum 3100/25, margin 32/-3 (the figure at landing is still measured at landing)
FAIL K6: runs/T-jy6f/conditions.md:0: missing file
```

- **Assessment:** K4 fails on `fa7a476e07815ce43004170a84b53ab3c55ee100` for the hash-less fix-cycle rows (`SM29`, `SM30`, `SM31` at `:396-398`, while `SM18` had restore evidence in the first sweep at `:319` and `SM19a` is skipped as `not run`), and passes (`ok K4: restore evidence of mutants`) on `38b5176c633493aec76e11763193c622883c4520` once the separate hash table at `:466-472` is present. K5 also catches finding 1 of rounds 2 and 3 (the file budget overrun).

## 3. #426 (`T-r7b9`) round-2 head

- **Head:** `9bb845e7252186388ee63feb38547db86de2ba5c` (from `Commit reviewed` in comment 6289957559).
- **Base:** `0246d9c72ecae7c3509a7522990cbf39706f3fd4` (comment 6289957559).
- **Command:** `sh /home/layup/projects/layup/.worktree/T-b6fn/docs/tests/freeze-check.sh --task T-r7b9 --base 0246d9c72ecae7c3509a7522990cbf39706f3fd4 --budget 2500/27`
- **Start / end / exit:** `2026-10-20T14:36:51Z` to `2026-10-20T14:36:51Z`, `exit=1`.
- **Output:**

```
FAIL K1: runs/T-r7b9/test-runs.md:293: new blank line at EOF.
FAIL K2: internal/run/questionloop_integration_test.go:163: TestAVerifierThatAsksIsAnsweredAtTheHeadItWasBoundTo is not in docs/tests/traceability.md
FAIL K2: internal/run/questionloop_integration_test.go:203: TestTheQuestionOfAVerifierAtItsBoundHeadLeavesTheFindingOpenWhileItWaits is not in docs/tests/traceability.md
FAIL K2: internal/run/questionloop_integration_test.go:258: TestTheAnswerPromptHoldsTheQuestionsAndTheNextPromptsTheAnswerIDs is not in docs/prd/PRD-0001-layup.md §12 REQ-005 Test cell
FAIL K2: internal/run/questionloop_integration_test.go:279: TestAHarnessThatFollowsTheAnswerInstructionGivesAValidAnswerHandoff is not in docs/prd/PRD-0001-layup.md §12 REQ-005 Test cell
FAIL K2: internal/run/questionloop_integration_test.go:293: TestAHumanRoutedQuestionAndAMixedOwnerQuestionWaitForTheOperator is not in docs/prd/PRD-0001-layup.md §12 REQ-005 Test cell
FAIL K2: internal/run/questionloop_integration_test.go:328: TestARestartAtEachStageOfTheQuestionWritesEachThingOnce is not in docs/prd/PRD-0001-layup.md §12 REQ-005 Test cell
FAIL K2: internal/run/questionloop_integration_test.go:389: TestAnOpenMaterialFindingKeepsTheTaskUphill is not in docs/prd/PRD-0001-layup.md §12 REQ-005 Test cell
FAIL K2: internal/run/questionloop_test.go:172: TestAnInvalidRecordedHandoffStaysAStop is not in docs/prd/PRD-0001-layup.md §12 REQ-005 Test cell
FAIL K2: internal/run/questionloop_test.go:191: TestAHumanQuestionIsAWaitWithNoAnswerSession is not in docs/prd/PRD-0001-layup.md §12 REQ-005 Test cell
FAIL K2: internal/run/questionloop_test.go:206: TestQuestionsOfTwoOwnerRolesAreAWait is not in docs/prd/PRD-0001-layup.md §12 REQ-005 Test cell
FAIL K2: internal/run/questionloop_test.go:232: TestTheNthAnswerNamesTheNthQuestionAndAMissingOrExtraAnswerStops is not in docs/prd/PRD-0001-layup.md §12 REQ-005 Test cell
FAIL K2: internal/run/questionloop_test.go:258: TestAnOwnerSessionThatDoesNotAnswerStops is not in docs/prd/PRD-0001-layup.md §12 REQ-005 Test cell
FAIL K2: internal/run/questionloop_test.go:282: TestAnAskingSessionIsThePositionQuestionAndOtherStatusesStayStops is not in docs/prd/PRD-0001-layup.md §12 REQ-005 Test cell
FAIL K2: internal/run/questionloop_test.go:296: TestTheStageThatAskedRunsAgainInTheNextAttempt is not in docs/prd/PRD-0001-layup.md §12 REQ-005 Test cell
FAIL K2: internal/run/questionloop_test.go:316: TestTheAnswerStepIsReadOnly is not in docs/prd/PRD-0001-layup.md §12 REQ-005 Test cell
FAIL K2: internal/run/questionloop_test.go:432: TestTheLastQualifiedReviewRecordIsTheSourceOfAnOpenMaterialFinding is not in docs/prd/PRD-0001-layup.md §12 REQ-005 Test cell
FAIL K2: internal/run/questionloop_test.go:538: TestTheAnswerSessionStartsAtTheBaseOfTheLastAskingSession is not in docs/prd/PRD-0001-layup.md §12 REQ-005 Test cell
ok K3: test selectors
ok K4: restore evidence of mutants
ok K5: projection: figure 2145/23, maximum 2500/27, margin 355/4 (the figure at landing is still measured at landing)
FAIL K6: runs/T-r7b9/conditions.md:0: missing file
```

- **Assessment:** Both named round-2 defects of #426 are caught:
  - F7: K2 reports `TestAVerifierThatAsksIsAnsweredAtTheHeadItWasBoundTo` (`:163`) and `TestTheQuestionOfAVerifierAtItsBoundHeadLeavesTheFindingOpenWhileItWaits` (`:203`) missing from `docs/tests/traceability.md` (and reports the sibling tests that `PRD-0001-layup.md` §12 `REQ-005` Test cell abbreviated as "and its ten siblings" / "and the five tests beside it").
  - F8: K1 reports `FAIL K1: runs/T-r7b9/test-runs.md:293: new blank line at EOF.`.
  - K3 and K4 pass on `#426`'s records (`ok K3: test selectors`, `ok K4: restore evidence of mutants`).

## 4. #426 again, after the fix of K2 (`5cf545b`)

The run of §3 also printed 16 `FAIL K2 … not in docs/prd/PRD-0001-layup.md §12 REQ-005 Test cell` lines. Those were false positives. The merged `main` names those tests by their file ("TestTheQuestionPositionIsTheFirstStageTheRecordsShowNotDone and its 13 siblings of questionloop_test.go"), and the same check on T-r7b9 as merged (`0246d9c..2866a3f`) failed K2 the same way. K2 now accepts a Test cell that names the test or the basename of its file (cases `good-K2-prd-by-file` and `bad-K2-prd-other-file`). On T-r7b9 as merged it prints `ok K2`. The rerun on the round-2 head:

- Command (in a detached worktree of `9bb845e7252186388ee63feb38547db86de2ba5c`): `sh docs/tests/freeze-check.sh --task T-r7b9 --base 0246d9c72ecae7c3509a7522990cbf39706f3fd4 --budget 2500/27`
- Start: `2026-10-20T15:09:26Z`, exit `1`

```
FAIL K1: runs/T-r7b9/test-runs.md:293: new blank line at EOF.
FAIL K2: internal/run/questionloop_integration_test.go:163: TestAVerifierThatAsksIsAnsweredAtTheHeadItWasBoundTo is not in docs/tests/traceability.md
FAIL K2: internal/run/questionloop_integration_test.go:203: TestTheQuestionOfAVerifierAtItsBoundHeadLeavesTheFindingOpenWhileItWaits is not in docs/tests/traceability.md
ok K3: test selectors
ok K4: restore evidence of mutants
ok K5: projection: figure 2145/23, maximum 2500/27, margin 355/4 (the figure at landing is still measured at landing)
FAIL K6: runs/T-r7b9/conditions.md:0: missing file
```

**Catch: yes, and only the named defects.** F7 (the two tests missing from `traceability.md`) is the pair of K2 lines, and F8 (the blank line at the end of `test-runs.md`) is the K1 line. K6 fails only because conditions tables did not exist then.
