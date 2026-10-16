# T-4cep: run evidence

Issue [#384](https://github.com/pharzam/layup/issues/384), parent [#356](https://github.com/pharzam/layup/issues/356). Serves F-0003#43 and F-0003#47
through REQ-003 and REQ-007. Author: Claude Sonnet 5.5 on Claude Code (no other model, tool model or agent); coordinator: Codex.
Plan 6216495237, independent plan review 6217115601, Operator ruling 6217133467 and approval 6217324567, audit note 6217334547 on #336,
author's acceptance of every condition 6217359299 (posted before any test was written). Plan base `82bb082`; `origin/main` `cfcfc1c` merged
as a fast-forward before any work. Maximum 1,300 added plus removed lines / 26 files, cap 2, goal classes 2 (admission, enablement).

| Item | Where |
| ---- | ----- |
| the red runs, the green runs, the clock times and exits | [`test-runs.md`](test-runs.md) |
| the mutation audit | [`mutants.md`](mutants.md), helper [`mutate.py`](mutate.py) |

## The conditions of the plan review, and where each is met

| # | Condition | Met by |
| - | --------- | ------ |
| 1 | no change before the review and the answers | the first commit is after 6217359299 (see `git log`) |
| 2 | one operative home | the conditions are written once, in `docs/spec/scaffold.md` "The tree check"; `loop.md` and the M2d-target subsection link to it |
| 3 | no identity reads | `TestCheckCodeOwnerLandingReadsOnlyTheObjectsOfTheConditions` (the list of reads) and the hostile-metadata control of `TestCheckCodeOwnerLanding` |
| 4 | inputs fixed | the Operator row (`TestLandingInputs`), the manifest at the last commit (the read list), `setup.CodeOwners` |
| 5 | seam discipline | `TestTheTreeRowAdmitsTheCodeOwnerCommitAndNothingElse`; `ForgeCheckSchema` is not in the diff |
| 6 | red against two stubs; valid mutants only | `test-runs.md`, `mutants.md` |
| 7 | the script diff is bounded | `git diff cfcfc1c HEAD -- internal/setup/codeowner.go` in the part of `landScript`: the signature, the doc comment, and the removed refusal block; nothing else |
| 8 | the recorded-task check by text only | `scaffold.md` and `loop.md` say row 66g's check is specified and not built |
| 9 | no live-target access | no forge call to `pharzam/layup-uat-m2d` was made, no generated command was run against a forge but the stub `gh` of the tests |
| 10 | goal classes stay 2 | the Operator's approval 6217324567 |
| 11 | budget | measured in the report of the build |

The reading closed by the author (R13) and posted in 6217359299: the ruling limits the act to "a target that M2d set up (its setup record has
no row `S13 codeowners.sha256`)", so the check of the run reads the record first and refuses a record that has the row; the tree predicate
independently refuses a parent that already holds the file. `TestLandingInputs` and `TestATargetSetUpWithTheRuleTakesNoCodeOwnerCommit` show each alone.

## What is not here

The full module integration level and the e2e level (the coordinator serializes them); the index and close-out files (`docs/tasks/T-4cep.md`, the
backlog and the completed log); the push, the pull request and the freeze; the token and time counters of the model (the coordinator reads
the CLI metadata); any read or write of the demo target, and the Operator-facing diff of the real target (presented at the demo).
