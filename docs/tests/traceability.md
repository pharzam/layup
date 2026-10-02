# Test traceability

LAYUP's test-side view of the traceability line, in the form of
[`traceability-template.md`](traceability-template.md): one row per test, so a
reader goes from any requirement of [`PRD-0001`](../prd/PRD-0001-layup.md) to the
test that covers it. The [DoD checklist](dod-checklist.md) asks that the ID set of
this table equals the requirement set: each `REQ` and `NFR` has at least one row.

Task `T-55n2` ([#76](https://github.com/pharzam/layup/issues/76)) wrote this table
from the [implementation plan](../plan/README.md). Each row carries its status at
`7cdd346`: `green` for a test that exists and passes, `planned` for a test that a
task of the plan will write. A planned row names its test by a handle of the form
`<task>/<level>/<name>`; the delivering task puts the real test name in its place,
with its status, in its own pull request. A row of a later milestone names the
milestone, whose specification task names the test.

A discipline check (for example check `pin` of
[`setup-check.sh`](../setup/setup-check.sh) for LAYUP's own pin) is not a row here,
because a row's level is one of `unit`, `integration`, `e2e` and `uat`;
`PRD-0001` §12 names such a check in its Test column. A test of the command
frame of [`spec/README.md`](../spec/README.md), "Commands", that proves no
requirement of `PRD-0001` (the version, the parser, the usage, the progress
lines) has `—` in Covers and in Fact, and its row names the frame (task `T-2yw7`,
#80). A row that covers more than one requirement gives the fact and the ADR of
the first.

| Test ID | Level | Covers (REQ/NFR) | Fact (F-NNNN#n) | Guardrail | ADR | Task | Status |
|---------|-------|------------------|-----------------|-----------|-----|------|--------|
| `TestGolden` (`internal/psb/check_test.go`) | unit | REQ-001 | F-0003#41 | — | ADR-0011 | T-dq05, T-5zmw | green |
| `TestCRLFGivesTheSameTable`, `TestALoneCarriageReturnBecomesASpace`, `TestG1ReadsTheValueOfAStack`, `TestEdgeCases`, `TestWriteTSVGivesTheErrorOfItsOutput` (`internal/psb/check_test.go`) | unit | REQ-001 | F-0003#41 | — | ADR-0011 | T-5zmw | green |
| `TestPSBCheckRefusesAFileThatIsNotUTF8`, `TestPSBCheckGivesTwoWhenItCannotWriteTheTable` (`internal/cli/psb_test.go`) | unit | REQ-001, NFR-004 | F-0003#41 | — | ADR-0011 | T-5zmw | green |
| `TestPSBCheckExitCodes` (`internal/cli/psb_integration_test.go`) | integration | REQ-001 | F-0003#41 | — | ADR-0011 | T-dq05, T-5zmw | green |
| `TestGoldenRealPSB` (`internal/psb/check_integration_test.go`) | integration | REQ-001 | F-0003#41 | — | ADR-0011 | T-dq05 | green |
| `TestEveryGoldenIsARecordOfTheBlock` (`internal/psb/check_integration_test.go`) | integration | REQ-001 | F-0003#41 | — | ADR-0011 | T-5zmw | green |
| `TestPSBCheckOnTheRealProblemStatement`, `TestPSBCheckOnAStatementWithNoGap`, `TestPSBCheckInputErrors` (`cmd/layup/psb_e2e_test.go`) | e2e | REQ-001, NFR-005 | F-0003#41 | — | ADR-0011 | T-5zmw | green |
| `TestPSBCheckWithAReadOnlyStandardOutput` (`cmd/layup/psb_e2e_test.go`) | e2e | REQ-001, NFR-004 | F-0003#41 | — | ADR-0011 | T-5zmw | green |
| `TestEverySchemaBlockIsBuiltOrNotYetBuilt` (`internal/tsv/blocks_integration_test.go`) | integration | NFR-001 | F-0001#1 | guardrails.md §1.1 Inv-1 | ADR-0014 | T-18v6 | green |
| `TestPackageRules` (`cmd/layup/rules_integration_test.go`) | integration | NFR-005, NFR-007 | F-0001#6 | guardrails.md §1.1 Inv-6 | ADR-0015 | T-2tc2 | green |
| `TestExitCode` (`internal/cli/cli_test.go`) | unit | NFR-004 | F-0001#5 | guardrails.md §1.1 Inv-5 | ADR-0011 | T-2yw7 | green |
| `TestUsageErrors` (`cmd/layup/usage_e2e_test.go`) | e2e | REQ-001 | F-0003#41 | — | ADR-0011 | T-2yw7 | green |
| `TestVersion` (`cmd/layup/usage_e2e_test.go`) | e2e | NFR-005 | F-0001#6 | guardrails.md §1.1 Inv-6 | ADR-0015 | T-2yw7 | green |
| `TestSameBytes` (`cmd/layup/bytes_test.go`) | unit | NFR-005 | F-0001#6 | guardrails.md §1.1 Inv-6 | ADR-0015 | T-2yw7 | green |
| `TestInputRule` (`cmd/layup/rules_integration_test.go`) | integration | NFR-005 | F-0001#6 | guardrails.md §1.1 Inv-6 | ADR-0015 | T-2yw7 | green |
| `TestCheckInputsFindsEachReadOfTheEnvironmentOrTheStandardInput`, `TestCheckInputsPassesAModuleThatReadsNoInput` (`cmd/layup/rules_test.go`) | unit | NFR-005 | F-0001#6 | guardrails.md §1.1 Inv-6 | ADR-0015 | T-2yw7 | green |
| `TestTheEmbeddedTestEntry` (`internal/catalog/catalog_integration_test.go`) | integration | REQ-004, REQ-002 | F-0003#44 | guardrails.md §1.1 Inv-7 | ADR-0016 | T-3jpx | green |
| `TestTheSchemaBlocks` (`internal/catalog/catalog_integration_test.go`) | integration | REQ-004 | F-0003#44 | guardrails.md §1.1 Inv-7 | ADR-0016 | T-3jpx | green |
| `TestReadGivesTheKindsInTheOrderOfTheFile`, `TestReadRefusesAnEntryThatBreaksARule`, `TestReadRefusesANameThatIsNotAnEntry`, `TestManifestIsTheKindsWithoutVersionFixtureAndEvidence`, `TestFixtureAndConfigOfAKind`, `TestStacksListsTheDirectoriesWithAKindsFile` (`internal/catalog/catalog_test.go`) | unit | REQ-004, REQ-002 | F-0003#44 | guardrails.md §1.1 Inv-7 | ADR-0016 | T-3jpx | green |
| `TestFilesReplaceTheModuleAndNothingElse` (`internal/catalog/catalog_test.go`) | unit | REQ-002 | F-0003#42 | guardrails.md §1.1 Inv-4 | ADR-0016 | T-3jpx | green |
| `TestHasNamesTheFilesOfTheEntryByTheirPathInIt` (`internal/catalog/catalog_test.go`) | unit | NFR-003 | F-0001#4 | guardrails.md §1.1 Inv-4 | — | T-3jpx | green |
| The frame of `spec/README.md`, "Commands": the four tests of `internal/cli/args_test.go`, the three of `internal/cli/progress_test.go`, and each test of `internal/cli/cli_test.go` other than `TestExitCode` | unit | — | — | — | ADR-0011 | T-t8qp, T-2yw7 | green |
| `TestGateOnAGoRepository` (`cmd/layup/gate_e2e_test.go`) | e2e | REQ-004, REQ-007, NFR-005 | F-0003#44 | guardrails.md §1.1 Inv-7 | ADR-0016 | T-5sgt | green |
| `TestGateNeverPassesACheckThatDidNotRun` (`cmd/layup/gate_e2e_test.go`) | e2e | NFR-004 | F-0001#5 | guardrails.md §1.1 Inv-5 | ADR-0011 | T-5sgt | green |
| `TestRunOnARealRepository`, `TestAHookOfTheRepositoryDoesNotRun`, `TestTheRevisionsAndTheManifestOnARealRepository`, `TestARenameAndADeleteChangeAProductPath`, `TestAFailedScratchTreeGivesAReasonWithNoPath`, `TestTheSchemaBlocks` (`internal/gate/gate_integration_test.go`) | integration | REQ-004, REQ-007 | F-0003#44 | guardrails.md §1.1 Inv-7 | ADR-0016 | T-5sgt | green |
| `TestACheckThatDidNotRunNeverPasses` (`internal/gate/result_test.go`) | unit | NFR-004 | F-0001#5 | guardrails.md §1.1 Inv-5 | ADR-0011 | T-5sgt | green |
| The other unit tests of `internal/gate` (`scope_test.go`, `manifest_test.go`, `result_test.go`) | unit | REQ-004 | F-0003#44 | guardrails.md §1.1 Inv-7 | ADR-0016 | T-5sgt | green |
| The tests of the run with a fake of `internal/git` (`internal/gate/run_test.go`; it writes temporary files) | integration | REQ-004, NFR-004 | F-0003#44 | guardrails.md §1.1 Inv-7 | ADR-0016 | T-5sgt | green |
| `TestGateUsageErrors`, `TestGatePrintsTheTableAndGivesItsExitCode`, `TestGateGivesTwoOnAnInputErrorOrALeftoverScratchTree`, `TestTheUsageListsGate` (`internal/cli/gate_test.go`) | unit | REQ-004, NFR-004 | F-0003#44 | guardrails.md §1.1 Inv-7 | ADR-0016 | T-5sgt | green |
| `TestLsTreeReadsEachEntry` (`internal/git/git_test.go`) | unit | REQ-004 | F-0003#44 | guardrails.md §1.1 Inv-7 | ADR-0016 | T-5sgt | green |
| `TestLsTree` (`internal/git/git_integration_test.go`) | integration | REQ-004 | F-0003#44 | guardrails.md §1.1 Inv-7 | ADR-0016 | T-5sgt | green |
| `TestARecordValueByStepAndName`, `TestARecordOfAnotherFormIsAnError` (`internal/work/work_test.go`) | unit | REQ-002 | F-0003#42 | guardrails.md §1.1 Inv-4 | ADR-0016 | T-6x75 | green |
| `TestTheSchemasEqualTheirBlocks`, `TestReadTheTwoFilesOfAWorkArea` (`internal/work/work_integration_test.go`) | integration | REQ-002 | F-0003#42 | guardrails.md §1.1 Inv-4 | ADR-0016 | T-6x75 | green |
| `TestPinFindings`, `TestKitHistoryFindings`, `TestIdentityFindings`, `TestThePinOfATarget`, `TestTheKitLinkOfATarget`, `TestTheChecksOfATarget` (`internal/verify/checks_test.go`) | unit | REQ-002, NFR-006 | F-0003#42 | guardrails.md §1.1 Inv-4 | ADR-0016 | T-6x75 | green |
| `TestRunGivesEachRowInTheOrderOfTheTable`, `TestCheckRunsTheNamedChecks`, `TestTheInputErrors`, `TestTheScratchTree`, `TestTheProgressLinesCoverTheScratchTree`, `TestTheTable`, `TestWithin` (`internal/verify/verify_test.go`) | unit | REQ-002, NFR-004 | F-0003#42 | guardrails.md §1.1 Inv-4 | ADR-0016 | T-6x75 | green |
| `TestTheFixturesOfSetupCheck` (`internal/verify/harness_integration_test.go`) | integration | REQ-002 | F-0003#42 | guardrails.md §1.1 Inv-4 | ADR-0016 | T-6x75 | green |
| `TestTheHitsOfAdapted`, `TestTheTextOfAdapted`, `TestTheFilesOfAdapted` (`internal/verify/adapted_test.go`) | unit | REQ-002 | F-0003#42 | guardrails.md §1.1 Inv-4 | ADR-0016 | T-8ya0 | green |
| `TestTheHashesOfFacts`, `TestTheFactsOfATarget`, `TestTheFactResolver`, `TestTheOnboardingOfATarget`, `TestTheGlossaryOfATarget`, `TestTheGuardrailsOfATarget` (`internal/verify/target_test.go`) | unit | REQ-002, NFR-003 | F-0003#42 | guardrails.md §1.1 Inv-4 | ADR-0016 | T-9t1q | green |
| `TestTheMarkersOfALine`, `TestScanMarkers`, `TestTheOpenGaps` (`internal/verify/markers_test.go`); `TestTheSourcesOfARecord` (`internal/verify/sources_test.go`); `TestTheBaselineScripts`, `TestTheBrokenLinks` (`internal/verify/scripts_test.go`) | unit | REQ-002, NFR-003, NFR-004 | F-0003#42 | guardrails.md §1.1 Inv-4 | ADR-0016 | T-8vpw | green |
| `TestTheScriptsOnATree`, `TestTheScriptsOfLAYUP` (`internal/verify/scripts_integration_test.go`); `TestTheExemptionsOfMarkersEqualTheSh`, `TestTheShAndTheGoFormOfMarkersAgree` (`internal/verify/harness_integration_test.go`) | integration | REQ-002, NFR-004 | F-0003#42 | guardrails.md §1.1 Inv-4 | ADR-0016 | T-8vpw | green |
| `TestTheListsOfAdaptedEqualTheSh`, `TestFlagged` (`internal/verify/harness_integration_test.go`) | integration | REQ-002 | F-0003#42 | guardrails.md §1.1 Inv-4 | ADR-0016 | T-8ya0 | green |
| `TestRunOnAStandInWorkArea`, `TestEachFindingOfATarget`, `TestTheInputErrorsOfAWorkArea`, `TestTheTableSchemaEqualsItsBlock` (`internal/verify/verify_integration_test.go`) | integration | REQ-002, NFR-004 | F-0003#42 | guardrails.md §1.1 Inv-4 | ADR-0016 | T-6x75 | green |
| `TestThePinText`, `TestTheRecordRows` (`internal/standin/standin_test.go`) | unit | REQ-002 | F-0003#42 | guardrails.md §1.1 Inv-4 | ADR-0016 | T-6x75 | green |
| `TestAStandInBaseline`, `TestAStandInWorkArea` (`internal/standin/standin_integration_test.go`) | integration | REQ-002 | F-0003#42 | guardrails.md §1.1 Inv-4 | ADR-0016 | T-6x75 | green |
| `TestIsShallow` (`internal/git/git_integration_test.go`) | integration | NFR-006 | F-0001#8 | guardrails.md §1.1 Inv-8 | ADR-0009 | T-6x75 | green |
| `TestSetupVerify` (`internal/cli/verify_test.go`) | unit | REQ-002, NFR-004 | F-0003#42 | guardrails.md §1.1 Inv-4 | ADR-0016 | T-6x75 | green |
| `TestSetupVerifyOnAStandInWorkArea`, `TestSetupVerifyInputErrors` (`cmd/layup/verify_e2e_test.go`) | e2e | REQ-002, NFR-004, NFR-005 | F-0003#42 | guardrails.md §1.1 Inv-4 | ADR-0016 | T-6x75 | green |
| `TestARunResumesAfterItsDoneSteps`, `TestTheOutcomes`, `TestTheHandOff`, `TestTheProseGroupStopsOnce`, `TestTheStopTableOrder`, `TestTheTables`, `TestTheCommandsFile`, `TestTheCommitOfAStep`, `TestTheInputs`, `TestTheStubs` (`internal/setup/setup_test.go`) | unit | REQ-002, NFR-001, NFR-004 | F-0003#42 | guardrails.md §1.1 Inv-4 | ADR-0016 | T-79y7 | green |
| `TestTheAnswersRule`, `TestCheckAsked`, `TestTheAnswersOfADoneStep`, `TestADoneStepWithNoAnswersHash`, `TestMarkerID` (`internal/setup/setup_test.go`) | unit | REQ-002, NFR-003 | F-0003#42 | guardrails.md §1.1 Inv-4 | ADR-0016 | T-79y7 | green |
| `TestTheSchemasEqualTheirBlocks`, `TestARunOnARealWorkArea`, `TestNoCommitOffTheSetupBranch` (`internal/setup/setup_integration_test.go`) | integration | REQ-002, NFR-005 | F-0003#42 | guardrails.md §1.1 Inv-4 | ADR-0016 | T-79y7 | green |
| `TestBranch` (`internal/git/git_integration_test.go`) | integration | REQ-002 | F-0003#42 | guardrails.md §1.1 Inv-4 | ADR-0016 | T-79y7 | green |
| `TestSetupCommand` (`internal/cli/setup_test.go`) | unit | REQ-002, NFR-004 | F-0003#42 | guardrails.md §1.1 Inv-4 | ADR-0016 | T-79y7 | green |
| `TestSetupExitCodesOnAWorkArea` (`internal/cli/setup_integration_test.go`) | integration | REQ-002, NFR-004 | F-0003#42 | guardrails.md §1.1 Inv-4 | ADR-0016 | T-79y7 | green |
| `TestSetup` (`cmd/layup/setup_e2e_test.go`) | e2e | REQ-002, NFR-005 | F-0003#42 | guardrails.md §1.1 Inv-4 | ADR-0016 | T-79y7 | green |
| `T-c06a/integration/go-fixtures` | integration | REQ-004, REQ-007 | F-0003#44 | guardrails.md §1.1 Inv-7 | ADR-0016 | T-c06a | planned |
| `TestSetupVerifyWithNoBaselineScript` (`cmd/layup/verify_e2e_test.go`) | e2e | NFR-004, REQ-002 | F-0001#5 | guardrails.md §1.1 Inv-5 | ADR-0011 | T-8vpw | green |
| `T-7s0y/integration/target-pin` | integration | NFR-006, REQ-002 | F-0001#8 | guardrails.md §1.1 Inv-8 | ADR-0009 | T-7s0y | planned |
| `T-dep6/e2e/whole-setup` | e2e | REQ-002, NFR-003 | F-0003#42 | guardrails.md §1.1 Inv-4 | ADR-0016 | T-dep6 | planned |
| `T-dep6/e2e/records-in-git` | e2e | NFR-001 | F-0001#1 | guardrails.md §1.1 Inv-1 | ADR-0014 | T-dep6 | planned |
| `T-dep6/e2e/gate-without-layup` | e2e | NFR-002 | F-0001#2 | guardrails.md §1.1 Inv-2 | ADR-0013 | T-dep6 | planned |
| `T-tmhw/integration/telemetry-schema` | integration | REQ-011 | F-0003#50 | — | ADR-0024 | T-tmhw | planned |
| `T-dgy7/integration/stalls-schema` | integration | REQ-009 | F-0003#49 | — | ADR-0023 | T-dgy7 | planned |
| `T-efmy/uat/release-review` | uat | REQ-015, REQ-017 | F-0003#53 | — | — | T-efmy | planned |
| `T-evad/uat/first-pilot` | uat | REQ-001, REQ-002, REQ-004, REQ-007, REQ-016, REQ-018, NFR-001, NFR-002, NFR-003 | F-0003#42 | guardrails.md §1.1 Inv-7 | ADR-0012 | T-evad | planned |
| `M2f/uat/rule-protection` | uat | REQ-003 | F-0003#43 | guardrails.md §1.1 Inv-3 | ADR-0017 | M2f | planned |
| `M2e/uat/role-handoffs` | uat | REQ-005 | F-0003#45 | — | ADR-0019 | M2e | planned |
| `M3c/uat/clarification` | uat | REQ-006 | F-0003#46 | — | ADR-0019 | M3c | planned |
| `M3b/e2e/escalation` | e2e | REQ-008 | F-0003#48 | — | ADR-0022 | M3b | planned |
| `M3d/e2e/stall-procedure` | e2e | REQ-010 | F-0003#49 | — | ADR-0023 | M3d | planned |
| `M2c/uat/specification-trace` | uat | REQ-012 | F-0003#51 | — | ADR-0018 | M2c | planned |
| `M4a/uat/two-harnesses` | uat | REQ-013 | F-0003#52 | guardrails.md §1.1 Inv-9 | ADR-0020 | M4a | planned |
| `M4c/uat/two-stacks` | uat | REQ-014 | F-0003#67 | — | — | M4c | planned |
