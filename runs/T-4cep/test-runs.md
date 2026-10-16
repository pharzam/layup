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

The counts are those of the red runs; the commit `a6ff498` later folded one case into the variant table (26 variants) and removed four duplicate rows of the unit table (12 refusals).

One defect of a test, found at the first green run of `TestParentsAndTheWholeTree` (14:23Z): the fixture added the submodule entry before `Add`, which drops it; fixed by adding it after (guardrails.md §2).

## Green, step by step (commits `e3db27a`, `3064bb8`, `e49e564`, `3069526`, `311ea87`; the hook ran gofmt, vet and the unit level on each, exit 0)

git calls 14:23:09Z; predicate: unit 14:26:04Z, the 25 variants 14:26:13Z–14:26:21Z; tree row: unit and integration (with the existing restart test, 113.9 s) 14:28:49Z–14:30:46Z; the flip: unit and
29.7 s of integration 14:32:13Z–14:32:48Z; every exit 0.

## Mutants: 23 valid, 23 killed

`runs/T-4cep/mutate.py`, one replacement per variant (`m01` to `m23`, named in the helper), `go build` and `go vet -tags=integration ./internal/...` first (a variant that fails either is not
credited), then the unit level and, if it survives, the focused integration tests; the kill is a `--- FAIL` line of its own output; the hash of the file before and after proves the restore.
Final run 15:03:28Z–15:07:45Z on the final code, exit 0, hashes equal, control (build, unit, three integrations) 0. An earlier run, on code and a helper that later changed, had one variant that
did not build (`m07` passed `nil` and left `configs` unused) and counted it out; its repaired form is in this run. The equivalent mutant (the configs read at the head, not the last commit: a head
whose manifest differs is refused by the tree comparison first) has no row.

| Variants | Killed by |
| -------- | --------- |
| `m01` | `TestCheckCodeOwnerLandingReadsOnlyTheObjectsOfTheConditions`, `TestClassifyLanding` |
| `m02`, `m03`, `m04`, `m05`, `m06` | `TestClassifyLanding` |
| `m07` | `TestCheckCodeOwnerLanding` |
| `m08`, `m09` | `TestLandingInputs` |
| `m10` | `TestADetailIsTheFirstLineOfItsReasonInValidUTF8WithDIRForTheHost`, `TestEachCheckThatFailsIsItsRowAndTheTableStillHasEightRows` (+2) |
| `m11` | `TestEachCheckThatFailsIsItsRowAndTheTableStillHasEightRows`, `TestTheTreeRowAdmitsTheCodeOwnerCommitAndNothingElse` |
| `m12`, `m14`, `m15` | `TestTheTreeRowAdmitsTheCodeOwnerCommitAndNothingElse` |
| `m13` | `TestEachCheckThatFailsIsItsRowAndTheTableStillHasEightRows`, `TestTheRulesOfLayupProbeAreReadOnceWhateverTheWorkAreaAndTheRulesetFile` (+1) |
| `m16` | `TestTheLandingScriptAsProducedIsNoRefusal` |
| `m17`, `m18`, `m20` | `TestTheLandingCommitIsAdmitted`, `TestTheLandingScript` |
| `m19`, `m21` | `TestTheLandingScript` |
| `m22` | `TestTheOutputIsRead` |
| `m23` | `TestBlobID`, `TestCheckCodeOwnerLandingReadsOnlyTheObjectsOfTheConditions` (+1) |

## The final pass

| Check | Start–End | Exit | Result |
| ----- | --------- | ---- | ------ |
| `go build ./...`, `go vet ./...`, `gofmt -l internal cmd` | 15:08:43Z–15:08:46Z | 0, 0, 0 | clean |
| `go test -count=1 ./...` | 15:08:46Z–15:08:55Z | 0 | 21 packages ok |
| `-tags=integration ./cmd/layup -run TestPackageRules` | 15:08:55Z–15:08:57Z | 0 | ok |
| focused integration (git calls, the variants, lockstep, landing script, update, rebuilt area, tree row, and the existing restart test of the forge check) | 15:08:57Z–15:10:01Z | 0 | ok (setup 35.7 s, run 60.8 s) |
| `adr-lint`, `prd-lint`, `link-lint`, `run-discipline-tests`, `setup-check`, `nested-checkout-check`, `git diff --check` | 15:10Z–15:15:42Z in two calls (the first was cut off at setup-check, exit 137, and that check and the next two were rerun, 15:15:11Z–15:15:42Z) | 0 each | OK; 2,771 links; 138 passed, 0 failed; 10 cases behaved |
| standalone `gitleaks git --redact` | after the last commit, see the report | 0 | no leaks |

**Not run, by the instruction:** the full module integration level, the e2e level, the close-out files, the push and the review. The token and time counters belong to the coordinator.

**Consolidation and refactors, said plainly.** Commit `53ceb15` held the evidence as `README.md`, `mutants.md` and this file; they were merged into this file to save files and lines. Everything they held survives here
or at that commit: the conditions of the plan review (above), the per-variant kills (grouped above; the ungrouped table of the first run is in `53ceb15:runs/T-4cep/mutants.md`). `579c035` removed the nil guards of the job's
two seams (`landing`, `note`) and their two tests: `forgeCheckStep` always gives both and no other job is built, so the guards were unreachable; the guarded refusal of a head asked of no check is still
shown by `TestTheTreeRowAdmitsTheCodeOwnerCommitAndNothingElse` and `m12`. `a6ff498` folded one case into the variant table and dropped four unit rows that the real-object table covers. No acceptance coverage was removed.


Coordinator before full integration: merge current maina3a4cdd (#360). Two textual
conflicts are resolved: PRD changelog keeps #360 then #384; packages.md retains
CheckCodeOwnerLanding ownership and the answer result kind settling #377.
The first independent review must judge both, plus the earlier f86 guardrail merge.
No source identity field is read as authorization; no real target was accessed.

Coordinator full integration: `go test -count=1 -tags=integration -timeout 25m ./...`
on `be0a38f96e979db32ddd27697bceb8e700d0af90`, 2026-10-16
15:22:49.182859–15:34:49.256158 UTC, direct exit 0, all 21 packages;
elapsed 720,073.299 ms. This includes current main's code.
Resource metadata: Sonnet plan I90/O69078/Cw224606/Cr5860322, 654549 ms;
build and bounded finish increment I356/O151617/Cw262937/Cr71545742.
The timed-out build reports no actual elapsed; bounded finish 127825 ms.
SWE-2 Max plan review cumulative P1522611/O40627/C1338011, 32 steps;
elapsed unreported. Its earlier 28-step export is a subset, never added.
Counters use each CLI's own fields; cached tokens are not added to prompt.
