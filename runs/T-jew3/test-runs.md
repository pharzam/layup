# T-jew3: the test runs

Row 65b of `M2e` ([#356](https://github.com/pharzam/layup/issues/356)); facts `F-0003#43` and `F-0003#47`. Times are UTC, read from the clock
(`date -u`) after each run; an exit status is the one of `go test`, taken from the shell. Base: `8e26eca7d01d26e47ed508e26f673a0d34bf7c9c`.

## The spike (plan review, note N4), before any broad code

A throwaway Go program outside the tree (`/tmp/t-jew3-spike`, not committed) served the object `default.json` of `runs/T-sdtb/rulesets.txt` from
`httptest.NewTLSServer`, trusted by `SSL_CERT_FILE`, with `GH_HOST`, `GH_ENTERPRISE_TOKEN` and `GH_CONFIG_DIR` set, and ran the real `gh` 2.102.0.
At 2026-10-16T11:02:00Z, exit 0: `gh api … --jq '<filter>'` gave the live object reduced to the fields of a write with `require_code_owner_review`
`true`, keys sorted, compact, one line feed at the end; `gh api --method PUT --input FILE` sent the file's bytes; and `--jq 'if … then error("…") else … end'`
exited 1 with the message on standard error. The dependency is available (host: `gh` 2.102.0, `jq` 1.7 present and not used); the tests that run `gh` fail,
and never skip, when it is absent.

## The red slice

The stubs of `internal/setup/codeowner.go` and `internal/run/codeowner.go` have the names and none of the behaviour, so each test fails on its own assertion.

- Unit, `go test -count=1 ./internal/setup ./internal/run`, 2026-10-16T11:09:33Z to 11:09:34Z, exit 1: `TestCodeOwnersText`, `TestCodeOwnersRefusals`,
  `TestS13WritesCodeOwnersInTheTree`, `TestS13InMigrateModeWritesTheFilesOfM2d`, `TestTheRulesetRequiresTheCodeOwner`, `TestSetupAreaChoosesTheMode`,
  `TestOperatorLoginIsTheRowOfTheRole`, `TestTheMigrationNamesItsFiles`, and the three tests of the earlier rows that the rule changes
  (`TestTheRulesetOfTheDefaultBranchInPhase2`, `TestTheRulesAreThoseOfTheTables`, `TestPhase2S13WritesBothRulesetsAndTheRow`). `TestTheStubFormIsTheLiveForm`
  is green already: it tests the fixture (the first-hand object), not the code.
- Integration, `go test -count=1 -tags=integration ./internal/setup -run 'TestTheUpdate|TestTheLandingScript|TestARebuiltArea|TestTheStubForm'`,
  11:09:34Z to 11:09:43Z, exit 1: `TestTheUpdateKeepsWhatIsUnknown`, `TestTheUpdateRefusals` (the stub's commands are those of M2d, which read
  `out/ruleset-default.json`), `TestTheLandingScript` (the script as produced exits 0), `TestARebuiltAreaReproducesTheRecordedCommits` (S13 adds no CODEOWNERS).

## The tests, by the dependencies they use

Unit (no program, file beyond a temporary directory, or network): `TestCodeOwnersText`, `TestCodeOwnersRefusals`, `TestTheRulesetRequiresTheCodeOwner`,
`TestS13WritesCodeOwnersInTheTree`, `TestS13InMigrateModeWritesTheFilesOfM2d`, `TestSetupAreaChoosesTheMode`, `TestTheStubFormIsTheLiveForm` (the seven of the plan);
the plan's mode table also holds two cases that I wrote as their own unit tests, `TestOperatorLoginIsTheRowOfTheRole` and `TestTheMigrationNamesItsFiles`.
Integration (a real program across a real seam): `TestARebuiltAreaReproducesTheRecordedCommits` (real `git`, the step runner), `TestTheUpdateKeepsWhatIsUnknown`
and `TestTheUpdateRefusals` (real `sh` and `gh`, an `httptest` TLS server, `git` for the guard), `TestTheLandingScript` (real `sh`, `git` with a bare repository, `gh`).

## The green runs (focused; the full integration level is the coordinator's)

All times UTC from the clock; statuses are those of `go test`.

- Unit, `go test -count=1 ./internal/setup ./internal/run`: exit 0 after the code (2026-10-16T11:10Z to 11:12Z, repeated at the controls below).
- Integration of the new tests, `go test -count=1 -tags=integration ./internal/setup -run 'TestTheUpdate|TestTheLandingScript|TestARebuiltArea|TestTheStubForm'`,
  11:12:55Z to 11:13:16Z, exit 0 (the first green attempt, 11:12:36Z to 11:12:47Z, exit 1: the fixture's `source` was `pharzam/layup-uat-m2d`, not the test
  repository's `acme/widget`, and the update's refusal worked; the tests got `liveFor`).
- Integration of the part of `internal/run` that the change touches, `go test -count=1 -tags=integration -timeout 20m ./internal/run -run 'Scaffold|ForgeCheck|TheStepSetup|TheStepVerify|Setup|Probe|Markers|Prose|Phase|Lease|Session|Code|Operator|Migration'`,
  11:14:24Z to 11:19:49Z, exit 0, `ok internal/run 319.432s` (a first attempt at 11:14:03Z did not build: the unit file imported `work` under the name of a function of an
  integration test; fixed by an alias). It is not repeated. The full module integration level was not run by the author.

## The controls and the checks after the documents

See [`mutants.md`](mutants.md) for the controls of the mutation audit (exit 0 for build, vet, unit and the three integration tests at 11:33Z).
