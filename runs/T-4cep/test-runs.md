# T-4cep: test runs

Every time is the UTC clock of the host (`date -u`) read right after the command; every exit is the direct status of the command. Tools: go1.26.9, git 2.56.0,
gh 2.102.0 (the integration tests need `gh` on PATH and fail, never skip, without it). The raw output of each run is kept by the session; the lines below are copied from it.
The pre-commit hook runs the unit level of the whole module, so a red slice cannot be committed alone: each red was observed on the working tree, against the stubs below, and the tests
were committed with the green code of their step (a hook is never skipped; guardrails.md §2, "A red slice that the pre-commit hook will not take").

## The stubs of the red slice

Signatures that compile and no behaviour: `git.Parents` returned an error "not implemented" and `git.BlobID` returned "", `setup.CheckCodeOwnerLanding` returned "not implemented"
(stub A: refuses everything, as the tree check did before this task), `setup.classifyLanding` returned "not implemented", `landScript` was a shim over the old flag (the script still
refused at its first command), the job had the two new fields and consulted neither, and the two run helpers returned "not implemented".

## Red A: the stubs refuse everything (the right reason: the behaviour is missing)

| Run | Start | End | Exit | Red |
| --- | ----- | --- | ---- | --- |
| `go test -count=1 ./internal/git ./internal/setup ./internal/run` (unit) | 14:20:13Z | 14:20:35Z | 1 | git: `TestEachCallRunsItsVerb/rev-list_--parents` ("0 starts of git, want 1"), `TestTheOutputIsRead` (the three outputs of Parents), `TestBlobID`, `TestParentsRefusesAnOutputOfAnotherForm`; setup: `TestClassifyLanding` (all 18 cases "got not implemented"), `TestCheckCodeOwnerLandingReadsOnlyTheObjectsOfTheConditions`, `TestTheLandingScriptAsProducedIsNoRefusal` (the script holds "not applicable" and "Q1 of #356"); run: `TestTheMigrationNamesItsFiles`, `TestLandingInputs`, `TestTheTreeRowAdmitsTheCodeOwnerCommitAndNothingElse` |
| `-tags=integration ./internal/git -run TestParentsAndTheWholeTree`; `./internal/run -run TestTheTreeRowOnRealObjects` | 14:20:56Z | 14:21:26Z | 1, 1 | git: Parents "not implemented"; run: the canonical commit "fail the head of main is …, and setup/tree.tsv ends at …" with no reason (the job did not ask the check), each refusal without its reason |
| `-tags=integration ./internal/setup -run 'TestCheckCodeOwnerLanding\|TestAHeadAbove\|TestATargetSetUp\|TestTheLandingCommit\|TestTheLandingScript$'` | 14:21:35Z | 14:21:52Z | 1 | `TestCheckCodeOwnerLanding` (25 variants and both controls "got not implemented"), `TestAHeadAboveTheCanonicalCommitIsRefused`, `TestATargetSetUpWithTheRuleTakesNoCodeOwnerCommit`, `TestTheLandingCommitIsAdmitted` and `TestTheLandingScript` ("the landing exited 1: not applicable: …(Q1 of #356, open)") |

## Red B: the stub admits everything (each refusal must be red, each control green)

A first attempt at 14:22:21Z did not build (an unused import in the stub) and is not counted. Run at 14:22:32Z to 14:22:48Z, exit 1 for each of the two commands:
`TestClassifyLanding` failed on its 16 refusal cases (the two admit cases passed), `TestCheckCodeOwnerLandingReadsOnlyTheObjectsOfTheConditions` failed on its wrong-parent and login cases,
and `TestCheckCodeOwnerLanding` failed on its 25 refusal variants (both controls passed), `TestAHeadAboveTheCanonicalCommitIsRefused` and `TestATargetSetUpWithTheRuleTakesNoCodeOwnerCommit` failed.

## The steps, green

| Step | Command | Start–End | Exit | Result |
| ---- | ------- | --------- | ---- | ------ |
| git calls | `go test -count=1 ./internal/git`; focused integration | 14:23:09Z | 0 | ok. **A defect of the test, not of the code:** the first green run of `TestParentsAndTheWholeTree` failed because the fixture added the submodule entry before `Add`, which drops a path with no file; the entry is now added after `Add` (guardrails.md §2) and the test passes |
| the predicate | unit `./internal/setup`; `-tags=integration -run 'TestCheckCodeOwnerLanding$\|TestAHeadAbove\|TestATargetSetUp'` | 14:26:04Z; 14:26:13Z–14:26:21Z | 0, 0 | ok; the 25 variants and the controls |
| the tree row | unit `./internal/run`; `-tags=integration ./internal/run -run 'TestTheTreeRowOnRealObjects\|TestTheStepForgeCheckChecksTheForgeAndCommitsItsTableAtEachRestart'` | 14:28:49Z–14:30:46Z | 0, 0 | ok (113.9 s, the second test is the existing restart test: the "one commit ahead" case still fails on the row `tree`) |
| the flip | unit of the three packages; `-tags=integration ./internal/setup -run 'TestTheLanding\|TestTheUpdate\|TestARebuilt\|TestCheckCodeOwnerLanding$\|TestAHeadAbove\|TestATargetSetUp\|TestTheStubForm'` | 14:32:13Z–14:32:48Z | 0, 0 | ok (29.7 s) |

Commits (the hook ran `gofmt`, `go vet ./...` and `go test ./...` on each): `e3db27a` (14:24Z), `3064bb8` (14:27Z), `e49e564` (14:31Z), `3069526` (14:32Z), `311ea87` (14:45Z), exit 0.

## Mutants

[`mutants.md`](mutants.md): 23 valid variants, all killed; the first form of one variant did not build and is not credited. 14:34:31Z to 14:43:20Z, control green.

## The final pass, on the committed code (head `311ea87` plus this evidence; no Go file changes after `3069526`)

| Check | Start–End | Exit | Result |
| ----- | --------- | ---- | ------ |
| `go build ./...`, `go vet ./...`, `gofmt -l internal cmd` | 14:45:53Z–14:45:57Z | 0, 0, 0 | clean |
| `go test -count=1 ./...` (the unit level) | 14:45:57Z–14:46:06Z | 0 | 21 packages ok |
| `go test -count=1 -tags=integration ./cmd/layup -run TestPackageRules` | 14:46:14Z–14:46:17Z | 0 | ok |
| the focused integration set (git, setup, run: the calls, the 25 variants, the lockstep, the landing script, the update, the rebuilt area, the tree row) | 14:46:17Z–14:46:47Z | 0 | ok (`internal/setup` 28.9 s) |
| `sh docs/adr/adr-lint.sh`, `prd-lint.sh`, `links/link-lint.sh` | 14:36:53Z–14:36:58Z, on the documents of `311ea87` | 0, 0, 0 | OK; 2,766 links resolved |
| `sh docs/tests/run-discipline-tests.sh`, `sh docs/setup/setup-check.sh`, `git diff --check` | 14:37:42Z–14:38:59Z | 0, 0, 0 | 138 passed, 0 failed; setup-check OK |
| `sh docs/tests/nested-checkout-check.sh` | 14:39:03Z | 0 | 10 cases behaved |
| standalone `gitleaks git --redact` | 14:47:06Z–14:47:17Z | 0 | 1,023 commits scanned, no leaks |

**Not run, by the instruction:** the full module integration level (`go test -tags=integration ./...`), the e2e level, and the index and close-out files; the coordinator serializes the first and the
last round. The integration tests that ran are the focused ones named above.
