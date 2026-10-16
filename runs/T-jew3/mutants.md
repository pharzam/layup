# T-jew3: the mutation audit

Each variant changes `internal/setup/codeowner.go` or `internal/run/codeowner.go` alone, runs the tests, and is restored. A variant counts only if
`go build` and `go vet` pass on it and its kill is a named assertion in its own output.

## First pass (helper `/tmp/mut/mutants.py`, 2026-10-16T11:20:19Z to 11:25:06Z): 22 variants, 19 valid, 3 invalid

The helper read the lines `--- FAIL` of standard output and ignored the exit status. **Three of its variants did not build and were counted as kills; they are not
credited:** `name guard missing`, `source guard missing` and `one pull_request rule not required` removed a `%s` of a `fmt.Sprintf` and left its argument, so `go vet` refused
the package (`BUILD FAILED`). The other 19 were killed by a test, by its own `--- FAIL` line (the first column is the variant, the second the test that named it):

| Variant | Killed by |
| ------- | --------- |
| no code-owner flag | `TestTheRulesetRequiresTheCodeOwner`, `TestTheRulesAreThoseOfTheTables`, `TestTheRulesetOfTheDefaultBranchInPhase2`, `TestS13InMigrateModeWritesTheFilesOfM2d` |
| CODEOWNERS left out of the tree | `TestS13WritesCodeOwnersInTheTree`, `TestPhase2S13WritesBothRulesetsAndTheRow` |
| CODEOWNERS in the tree in migrate mode | `TestS13InMigrateModeWritesTheFilesOfM2d` |
| POST instead of PUT in the update | `TestS13InMigrateModeWritesTheFilesOfM2d` |
| a parameter of pull_request dropped | `TestTheUpdateKeepsWhatIsUnknown` |
| bypass list dropped from the write fields | `TestTheUpdateKeepsWhatIsUnknown`, `TestTheUpdateRefusals`, `TestTheLandingScript` |
| read back not compared | `TestTheUpdateRefusals` |
| trap missing | `TestTheLandingScript` |
| restore puts the BEFORE object (an unbound `$before`, so a weak variant) | `TestTheLandingScript` |
| a failing restore swallowed | `TestTheLandingScript` |
| the head guard missing in the update | `TestTheUpdateRefusals` |
| landing script applicable | `TestS13InMigrateModeWritesTheFilesOfM2d` |
| login unchecked | `TestCodeOwnersRefusals`, `TestS13WritesCodeOwnersInTheTree` |
| config path unchecked | `TestCodeOwnersRefusals` |
| prefix cover without the slash | `TestCodeOwnersText` |
| the update pushes | `TestS13InMigrateModeWritesTheFilesOfM2d` |
| the mode ignores the row | `TestSetupAreaChoosesTheMode` |
| rulesets never applied: no error | `TestSetupAreaChoosesTheMode` |
| the owner is the idea owner | `TestOperatorLoginIsTheRowOfTheRole` |

That is 19 valid kills of the first pass, credited from the helper's log. An unrecorded rerun of three variants in the same session, 11:25Z, is not credited either.

## Second pass (`mutate.py`, 2026-10-16T11:31:07Z to 11:33:01Z): 3 re-expressed variants and one added, all valid, all killed

The three invalid variants were written again with every placeholder kept (`X and false`, the same formatted guard), and one variant was added for the rule that the landing puts back the
`AFTER` object (it restores `default.before.json`). The harness keeps the exit status and the whole standard output and error of each run, a fresh backup outside the tree, and the
SHA-256 of the file before and after. All four: build 0, vet 0, unit 0, then the integration tests of the update and the landing exit 1:

| Variant | Build, vet, unit | Killed by (integration) |
| ------- | ---------------- | ----------------------- |
| name guard (`(.name != "%s" and false)`) | 0, 0, 0 | `TestTheUpdateRefusals`: `another name: exit 0, 1 PUT … want a refusal that names "name other" and no write` |
| source guard | 0, 0, 0 | `TestTheUpdateRefusals`: `another source: exit 0, 1 PUT … names "source acme/other"` |
| one pull_request rule (`((%s) != 1 and false)`) | 0, 0, 0 | `TestTheUpdateRefusals`: `no pull_request rule: exit 0, 1 PUT … names "0 pull_request rules"` (and the case of two) |
| landing restores the object before the update | 0, 0, 0 | `TestTheLandingScript`: `the landing exited 1` (the read back of the object put back differs from the `AFTER` object) |

Hash of `internal/setup/codeowner.go` before and after all four: `71fb9f6d5a70e7e9651e07d21d3a9dc2231b8f7f2f9bc58e3e5dca2e823ce89e`, equal; no `.orig` file in the tree.
**Control** (the restored file): `go build ./...` 0, `go vet ./internal/setup` 0, `go test ./internal/setup` 0, the three integration tests 0 (18 s).

**Count:** 22 variants of the first pass, of which 19 valid and killed and 3 invalid (not credited); 4 variants of the second pass, valid and killed, of which 3 stand for the invalid
three. So 22 valid variants are killed in all (19 + 3), plus one added; the first pass alone does not give 22.
