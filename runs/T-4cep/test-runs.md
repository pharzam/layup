# T-4cep: test runs and mutants

Issue [#384](https://github.com/pharzam/layup/issues/384), child of [#356](https://github.com/pharzam/layup/issues/356); serves F-0003#43 and F-0003#47 (REQ-003, REQ-007).
Author: Claude Sonnet 5.5 (no other model, tool model or agent). Plan 6216495237, plan review 6217115601, ruling 6217133467, approval 6217324567, my acceptance of every condition
6217359299 (posted before any test). Maximum 1,300 lines / 26 files, cap 2, goal classes 2 (admission, enablement). All times are the host clock (`date -u`), UTC, 2026-10-16;
every exit is direct. Tools: go1.26.9, git 2.56.0, gh 2.102.0 (the integration tests need `gh` and fail, never skip, without it).

**The conditions of the plan review, where met.** 2: the three conditions are written once, in `scaffold.md` "The tree check". 3: `TestCheckCodeOwnerLandingReadsOnlyTheObjectsOfTheConditions` (the list
of reads) and the hostile-metadata control of `TestCheckCodeOwnerLanding`. 4: `TestLandingInputs`, the manifest read at the last commit, `setup.CodeOwners`. 5: `TestTheTreeRowAdmitsTheCodeOwnerCommitAndNothingElse`;
`ForgeCheckSchema` is not in the diff. 6: this file. 7: in `landScript` the diff against `cfcfc1c` is the signature, its comment and the removed refusal block. 8: `scaffold.md` and `loop.md` say
row 66g's check is specified and not built. 9: no call to `pharzam/layup-uat-m2d`; the only `gh` that ran is the stub of the tests. The reading closed in 6217359299: the ruling limits the act to a target whose
setup record lacks `S13 codeowners.sha256`, so the run refuses a record that has it (`TestLandingInputs`) and the tree predicate independently refuses a parent that has the file
(`TestATargetSetUpWithTheRuleTakesNoCodeOwnerCommit`).

## Red

The pre-commit hook runs the unit level of the whole module, so a red slice cannot be committed alone: each red was observed on the working tree against stubs, with the output kept, and the tests
were committed with the green code of their step; a hook is never skipped. The stubs compile and do nothing: `git.Parents` and the two functions of `internal/setup` answered "not implemented",
`BlobID` answered "", `landScript` still refused at its first command, the job did not ask the check.

| Run | Start–End | Exit | Red (the right reason) |
| --- | --------- | ---- | ---------------------- |
| A, stubs refuse all: unit of git, setup, run | 14:20:13Z–14:20:35Z | 1 | the Parents rows, `TestBlobID`; all 18 cases of `TestClassifyLanding` ("got not implemented"); the reads test; the script "holds not applicable"; `TestLandingInputs`; the tree-row test |
| A: integration git, run | 14:20:56Z–14:21:26Z | 1, 1 | Parents not implemented; the canonical head "fail … ends at …" with no reason |
| A: integration setup (variants, lockstep, landing) | 14:21:35Z–14:21:52Z | 1 | 25 variants and 2 controls "got not implemented"; `TestTheLandingScript`, `TestTheLandingCommitIsAdmitted` "exited 1: not applicable … Q1 of #356" |
| B, stub admits all: classifier, variants | 14:22:32Z–14:22:48Z | 1, 1 | the 16 refusals of the table and the 25 variants fail, both controls pass (a first try at 14:22:21Z did not build and is not counted) |

One defect of a test, found at the first green run of `TestParentsAndTheWholeTree` (14:23Z): the fixture added the submodule entry before `Add`, which drops it; fixed by adding it after (guardrails.md §2).

## Green, step by step (commits `e3db27a`, `3064bb8`, `e49e564`, `3069526`, `311ea87`; the hook ran gofmt, vet and the unit level on each, exit 0)

git calls 14:23:09Z; predicate: unit 14:26:04Z, the 25 variants 14:26:13Z–14:26:21Z; tree row: unit and integration (with the existing restart test, 113.9 s) 14:28:49Z–14:30:46Z; the flip: unit and
29.7 s of integration 14:32:13Z–14:32:48Z; every exit 0.

## Mutants: 23 valid, 23 killed

`runs/T-4cep/mutate.py`, one replacement per variant, `go build` and `go vet -tags=integration ./internal/...` first (a variant that fails either is not credited), then the unit level and, if it
survives, the focused integration tests; the kill is a `--- FAIL` line of its own output; the hash of the file before and after proves the restore. Run 14:51:09Z–14:57:07Z on the final code, exit 0,
hashes equal, control (build, unit, three integrations) 0. An earlier run with an earlier helper had one variant that did not build (`m07` passed `nil` and left `configs` unused) and counted it out;
its repaired form is in this run. The equivalent mutant (the configs read at the head, not the last commit: a head whose manifest differs is refused by the tree comparison first) has no row.

| Variant | Killed by |
| ------- | --------- |
| `m01-parent-not-compared` | `TestCheckCodeOwnerLandingReadsOnlyTheObjectsOfTheConditions`, `TestClassifyLanding` |
| `m02-more-than-one-parent` | `TestClassifyLanding` |
| `m03-tree-not-compared` | `TestClassifyLanding` |
| `m04-mode-ignored` | `TestClassifyLanding` |
| `m05-bytes-not-compared` | `TestClassifyLanding` |
| `m06-last-may-have-the-file` | `TestClassifyLanding` |
| `m07-config-paths-ignored` | `TestAHeadAboveTheCanonicalCommitIsRefused`, `TestCheckCodeOwnerLanding` |
| `m08-login-is-the-first-row` | `TestLandingInputs` |
| `m09-record-gate-ignored` | `TestLandingInputs` |
| `m10-check-asked-of-the-last-commit` | `TestADetailIsTheFirstLineOfItsReasonInValidUTF8WithDIRForTheHost`, `TestEachCheckThatFailsIsItsRowAndTheTableStillHasEightRows` (+2) |
| `m11-check-error-is-a-pass` | `TestEachCheckThatFailsIsItsRowAndTheTableStillHasEightRows`, `TestTheTreeRowAdmitsTheCodeOwnerCommitAndNothingElse` |
| `m12-check-asked-of-no-head` | `TestTheTreeRowAdmitsTheCodeOwnerCommitAndNothingElse` |
| `m13-probe-on-the-last-commit` | `TestEachCheckThatFailsIsItsRowAndTheTableStillHasEightRows`, `TestTheRulesOfLayupProbeAreReadOnceWhateverTheWorkAreaAndTheRulesetFile` (+1) |
| `m14-no-note-on-admit` | `TestTheTreeRowAdmitsTheCodeOwnerCommitAndNothingElse` |
| `m15-detail-on-admit` | `TestTheTreeRowAdmitsTheCodeOwnerCommitAndNothingElse` |
| `m16-script-refuses-again` | `TestTheLandingScriptAsProducedIsNoRefusal` |
| `m17-script-adds-another-path` | `TestTheLandingCommitIsAdmitted`, `TestTheLandingScript` |
| `m18-script-commit-on-the-wrong-parent` | `TestTheLandingCommitIsAdmitted`, `TestTheLandingScript` |
| `m19-trap-missing` | `TestTheLandingScript` |
| `m20-restore-the-before-object` | `TestTheLandingCommitIsAdmitted`, `TestTheLandingScript` |
| `m21-no-disable-before-the-push` | `TestTheLandingScript` |
| `m22-git-parents-include-the-commit` | `TestTheOutputIsRead` |
| `m23-blob-header-without-the-length` | `TestBlobID`, `TestCheckCodeOwnerLandingReadsOnlyTheObjectsOfTheConditions` (+1) |

## The final pass

| Check | Start–End | Exit | Result |
| ----- | --------- | ---- | ------ |
| `go build ./...`, `go vet ./...`, `gofmt -l internal cmd` | 14:57:48Z–14:57:51Z | 0, 0, 0 | clean |
| `go test -count=1 ./...` | 14:57:51Z–14:57:58Z | 0 | 21 packages ok |
| `-tags=integration ./cmd/layup -run TestPackageRules` | 14:57:58Z–14:58:01Z | 0 | ok |
| focused integration (git calls, 25 variants, lockstep, landing script, update, rebuilt area, tree row) | 14:58:01Z–14:58:30Z | 0 | ok |
| the existing restart test of the forge check | 14:58:33Z–14:59:27Z | 0 | ok (52.6 s) |
| `adr-lint`, `prd-lint`, `link-lint`, `run-discipline-tests`, `setup-check`, `nested-checkout-check`, `git diff --check` | 14:59:57Z–15:00:49Z | 0 each | OK; 2,771 links; 138 passed, 0 failed; 10 cases behaved |
| standalone `gitleaks git --redact` | see the report | 0 | run after the last commit |

**Not run, by the instruction:** the full module integration level, the e2e level, the close-out files, the push and the review. The token and time counters belong to the coordinator.
