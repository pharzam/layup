# T-4cep: the mutation audit

Each variant changes one source file by one replacement, runs `go build ./...` and `go vet -tags=integration ./internal/...`, then the unit level of
`internal/setup`, `internal/run` and `internal/git`, and, if the unit level passes, the focused integration tests of the landing, the tree row and the git calls. A variant is
valid only if the build and the vet pass on it; its kill is a `--- FAIL` line of its own output. Every exit status and output is kept by the helper
([`mutate.py`](mutate.py), run from the root of the worktree), the hash of each file before and after proves the restore, and a control run on the restored tree
is green. First run 2026-10-16T14:34:31Z to 14:43:20Z on the code of the commit `3069526` plus documents (no Go file differs); the control after it: build 0, unit 0,
integration 0, 0, 0. Tools: go1.26.9, git 2.56.0, gh 2.102.0.

**One variant was not valid and is not credited:** `m07-config-paths-ignored` in its first form passed `nil` for `configs` and left the variable unused, so `go build` refused it
(`rc={'build': 1}`). Its second form, `CodeOwners(login, configs[:0])`, builds and vets, passes the unit level, and is killed at the integration level (a second run of the helper with that tag alone, control green).

| Variant | The file | Killed by (its own `--- FAIL` lines) |
| ------- | -------- | ------------------------------------ |
| `m01-parent-not-compared` | `internal/setup/codeowner.go` | `TestCheckCodeOwnerLandingReadsOnlyTheObjectsOfTheConditions`, `TestClassifyLanding` |
| `m02-more-than-one-parent` | `internal/setup/codeowner.go` | `TestClassifyLanding` |
| `m03-tree-not-compared` | `internal/setup/codeowner.go` | `TestClassifyLanding` |
| `m04-mode-ignored` | `internal/setup/codeowner.go` | `TestClassifyLanding` |
| `m05-bytes-not-compared` | `internal/setup/codeowner.go` | `TestClassifyLanding` |
| `m06-last-may-have-the-file` | `internal/setup/codeowner.go` | `TestClassifyLanding` |
| `m08-login-is-the-first-row` | `internal/run/codeowner.go` | `TestLandingInputs` |
| `m09-record-gate-ignored` | `internal/run/codeowner.go` | `TestLandingInputs` |
| `m10-check-asked-of-the-last-commit` | `internal/run/forgecheck.go` | `TestADetailIsTheFirstLineOfItsReasonInValidUTF8WithDIRForTheHost`, `TestEachCheckThatFailsIsItsRowAndTheTableStillHasEightRows`, `TestTheStepForgeCheckCommitsEightRowsAndIsDoneWhenEachIsAPass`, `TestTheTreeRowAdmitsTheCodeOwnerCommitAndNothingElse` |
| `m11-check-error-is-a-pass` | `internal/run/forgecheck.go` | `TestEachCheckThatFailsIsItsRowAndTheTableStillHasEightRows`, `TestTheTreeRowAdmitsTheCodeOwnerCommitAndNothingElse` |
| `m12-check-asked-of-no-head` | `internal/run/forgecheck.go` | `TestTheTreeRowAdmitsTheCodeOwnerCommitAndNothingElse` |
| `m13-probe-on-the-last-commit` | `internal/run/forgecheck.go` | `TestEachCheckThatFailsIsItsRowAndTheTableStillHasEightRows`, `TestTheRulesOfLayupProbeAreReadOnceWhateverTheWorkAreaAndTheRulesetFile`, `TestTheTreeRowAdmitsTheCodeOwnerCommitAndNothingElse` |
| `m14-no-note-on-admit` | `internal/run/forgecheck.go` | `TestTheTreeRowAdmitsTheCodeOwnerCommitAndNothingElse` |
| `m15-detail-on-admit` | `internal/run/forgecheck.go` | `TestTheTreeRowAdmitsTheCodeOwnerCommitAndNothingElse` |
| `m16-script-refuses-again` | `internal/setup/codeowner.go` | `TestTheLandingScriptAsProducedIsNoRefusal` |
| `m17-script-adds-another-path` | `internal/setup/codeowner.go` | `TestTheLandingCommitIsAdmitted`, `TestTheLandingScript` |
| `m18-script-commit-on-the-wrong-parent` | `internal/setup/codeowner.go` | `TestTheLandingCommitIsAdmitted`, `TestTheLandingScript` |
| `m19-trap-missing` | `internal/setup/codeowner.go` | `TestTheLandingScript` |
| `m20-restore-the-before-object` | `internal/setup/codeowner.go` | `TestTheLandingCommitIsAdmitted`, `TestTheLandingScript` |
| `m21-no-disable-before-the-push` | `internal/setup/codeowner.go` | `TestTheLandingScript` |
| `m22-git-parents-include-the-commit` | `internal/git/git.go` | `TestTheOutputIsRead` |
| `m23-blob-header-without-the-length` | `internal/git/git.go` | `TestBlobID`, `TestCheckCodeOwnerLandingReadsOnlyTheObjectsOfTheConditions`, `TestClassifyLanding` |
| `m07-config-paths-ignored` (second form) | `internal/setup/codeowner.go` | `TestAHeadAboveTheCanonicalCommitIsRefused`, `TestCheckCodeOwnerLanding` |

23 variants, 23 valid, 23 killed; hashes after equal hashes before for each file, both runs. The first form of m07 is the only variant counted out.

Differences from the plan's list, said plainly: the configs read at the head and not at the last commit is an equivalent mutant (a head whose manifest differs is refused by the
tree comparison before the bytes are judged), so it has no row; the script that commits an extra file is not a source variant but a case of `TestTheLandingCommitIsAdmitted`, whose
mutated scripts (another path, a second file) the check refuses, and `m17` is the source variant of the wrong path; the three variants of T-jew3 that the plan names, the trap
missing, the BEFORE object restored and the disable not before the push, are `m19`, `m20` and `m21`, on the production path.
