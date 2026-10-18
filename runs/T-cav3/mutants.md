# T-cav3 compiled mutants

Recorded at 2026-10-18T11:26Z. Each mutant is one source edit of a copy of the worktree at the frozen
test and code files (outside the tree), vetted (`go vet ./internal/run/`; a build failure
counts as not run), then run against the unit level, then the tree half of the integration
level (with the characterization tests of the tree items), then the demo. The file is
restored from a backup taken before the first mutant and compared by SHA-256 after each
one; at the end the four files are byte-identical (`cmp`) to the worktree's and the unit and
integration controls are green. No mutant failed to build.

| # | Mutant | Result, and the test that kills it |
| - | ------ | ----------------------------------- |
| M1 | `AnswersRecord` taken from the `answers.record` row (F-0001) | killed: `TestAnAnswersCitationIsAQuestionIDNotARecordLine`, `TestEachRecordFieldOfTheFullInputIsItsRecord` |
| M2 | any `F-NNNN-*.md` counts as the numbered facts record | killed: `TestOnlyTheNumberedFactsFileNamesAProblemStatementRecord` |
| M3 | a missing `spec/spans.tsv` gives no error | killed: `TestAMissingRecordOfTheFullInputIsAnError` |
| M4 | a missing `spec/inventory.tsv` is an error | killed: `TestNoInventoryIndexIsBeforeTheFirstBet` |
| M5 | a missing version is skipped | killed: `TestEachVersionOfTheIndexIsItsExactBytes` |
| M6 | `Register` read with `ReadRegister` | killed: `TestTheRegisterIsItsBytesAndNotJudged` |
| M7 | any event counts as merged | killed: `TestMergedIsTheTasksWithAMergedEvent` |
| M8 | `Changed` is nil | killed (integration): `TestTheTreeFieldsOfTheFullInputAreTheRegularFilesOfTheHead` |
| M9 | `SpecFiles` read at `base` | killed (integration): the same test |
| M10 | the PRD prefix check dropped | killed (integration): the same test (README.md) |
| M11 | `PSB` not set | killed: `TestEachRecordFieldOfTheFullInputIsItsRecord` |
| M12 | `Copy` ignored | killed (demo only, 82.9 s): `TestTheFullInputOfAHeadAfterIntakeAndTheScaffold` |
| M13 | the SHA check skipped | killed: `TestANonFullSHAIsRefusedBeforeAnyRead` |
| M14 | a refused `events.tsv` ignored | killed: `TestARecordThatItsReaderRefusesIsAnError` |
| M15 | two numbered facts files give the first | killed: `TestOnlyTheNumberedFactsFileNamesAProblemStatementRecord` |
| M16 | `filesAt` ignores the mode | killed (integration): the tree tests and `TestTheSpecFilesAreTheRegularMdBlobsOfDocsSpecAtTheHead` |
| M17 | `filesAt` reads below the directory | killed (integration): the same |
| M18 | `readRequired`: an empty file counts as missing | **not killed** at first (no test had a present empty record); killed after `TestReadRequiredKeepsTheStoreErrorTheMissingDiagnosticAndTheReaderPath` was added |
| M19 | `readRequired`: the reader's error is not wrapped with the path | **not killed** at first; killed by the same added test |
| M20 | the numbered facts name matches as a suffix | killed: `TestOnlyTheNumberedFactsFileNamesAProblemStatementRecord` |

No mutant is claimed equivalent. M12 is reached only by the demo, because the copy is
the one field that only a successful call carries (condition 3 of the plan review).
