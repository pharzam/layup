# T-9t1q: the test runs

The runs of `T-9t1q` (#89). Host: macOS, `go1.27.1`, `git` 2.54.0, 2026-10-02.
A path of a temporary directory shows as `<tmp>/`; "…" marks a cut part of an
output.

## The red runs

Each red run is of a test in its final form, on the base `bea865f`, on a
skeleton of the code (the names, with no behaviour), on a scratch commit, or
on a mutation of the code that breaks one rule.

1. **#61, the fixture first** (`facts/bad-blank-tab`, on the script of the
   base):

   ```text
   $ sh docs/setup/tests/run.sh
   FAIL  facts/bad-blank-tab: no output line: setup-check: facts FAIL verbatim: F-0001 fact 2 is empty
   setup-check tests: 40 passed, 1 failed
   ```

   Then `tr -d ' \t'` in the loop of `F-0001` and `F-0003`: 41 passed.

2. **#48.** A first draft held a Go test; on the base's `.gitattributes` it
   failed (`setup-check.sh: line 26: : command not found`, exit 2). The plan
   review asked for the discipline level (condition 3), so the case
   `facts/good-autocrlf` of `run.sh` replaced it. Its red runs are on scratch
   commits that hold the new runner and case:

   ```text
   # on e5d4727 (no line of #48 in .gitattributes)
   FAIL  facts/good-autocrlf: exit 2, want 0
         | <tmp>/clone/docs/setup/setup-check.sh: line 26: : command not found
         | <tmp>/clone/docs/setup/setup-check.sh: line 27: set: -: invalid option
   # on 81d7a35 (the script and docs/facts/ pinned, the list not yet)
   FAIL  facts/good-autocrlf: exit 1, want 0
         | setup-check: facts FAIL hash: docs/facts/problem-statement-brief.md does not match docs/setup/facts.sha256
         | setup-check: facts FAIL listed: docs/facts/problem-statement-brief.md is not in docs/setup/facts.sha256
   ```

   The second run found a third file that #48 did not name: the list
   `docs/setup/facts.sha256` gains carriage returns, and its paths do not
   match. It got `text eol=lf` (`bcfbeaa`); then 42 passed.

3. **The unit tests of the four checks, on the skeleton:** all six fail.

   ```text
   $ go test -count=1 -run 'Facts|Resolver|OfATarget' ./internal/verify/
   --- FAIL: TestTheHashesOfFacts
   --- FAIL: TestTheFactsOfATarget
   --- FAIL: TestTheFactResolver
   --- FAIL: TestTheOnboardingOfATarget
   --- FAIL: TestTheGlossaryOfATarget
   --- FAIL: TestTheGuardrailsOfATarget
   ```

4. **The integration tests, on the skeleton:** the harness compares the four
   cases of a target's form, and the new cases on the stand-in fail.

   ```text
   $ go test -count=1 -tags=integration -run 'TestTheFixturesOfSetupCheck|TestEachFindingOfATarget' ./internal/verify/
       harness_integration_test.go:141: facts/bad-hash: the lines of facts …
       harness_integration_test.go:141: facts/bad-not-verbatim: the lines of facts …
       harness_integration_test.go:141: onboarding/bad-all: the lines of onboarding …
       harness_integration_test.go:141: guardrails/bad-entries: the lines of guardrails …
       verify_integration_test.go:152: a brief that does not match its hash: [{"facts" "pass" ""}], <nil>; want facts fail, …
       verify_integration_test.go:152: no onboarding file: [{"onboarding" "pass" ""}], <nil>; want onboarding fail, …
   ```

   The rows of the built checks (`TestRunGivesEachRowInTheOrderOfTheTable`,
   `TestTheScratchTree`, `TestRunOnAStandInWorkArea`) failed on the first
   draft, before the check table, the stand-in and the tests were changed
   together.

5. **Mutations,** each a change of one rule of the head:

   | Mutation | Result |
   | -------- | ------ |
   | a record is read only once its step is done (condition 2 of the plan review) | fail: each case that reads a record before its step, among them a record of S11 that exists before S11 is done |
   | the list of hashes is required before S04 | fail: "a new tree: nothing is required" |
   | the first `Check: ` of a line, not the last | fail: the entry `Inv-9` |
   | a heading does not end an entry | fail: the entry `Inv-8` takes the `Check:` of a later paragraph |
   | the resolver reads the first of two records | fail: "two records" |
   | the harness compares all kinds of lines | fail: each case of LAYUP's form, for example `facts/bad-answers-blank` |

## The green runs

On the tree of the commit that adds this file.

| Command | Result |
| ------- | ------ |
| `go build ./...`; `go vet` with no tag, `-tags=integration` and `-tags=e2e`; `gofmt -l .` | exit 0; no file |
| `go test -count=1 ./...` | `ok` × 11 packages |
| `go test -count=1 -tags=integration ./...` | `ok` × 11; the harness compares `facts/bad-hash`, `facts/bad-not-verbatim`, `onboarding/bad-all`, `guardrails/bad-entries` and the frame case, and names the ten cases that it does not compare |
| `go test -count=1 -tags=e2e -timeout 10m ./...` | `ok` × 11 |
| `go test -race -count=1 ./internal/verify/ ./internal/standin/`, also with `-tags=integration` | `ok` |
| `sh docs/setup/tests/run.sh` | 42 passed, 0 failed |
