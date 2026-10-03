# T-d6q5: the test runs

The runs of `T-d6q5` (#92). Host: macOS, `go1.27.1`, `git` 2.54.0,
2026-10-03, UTC. "…" marks a cut part of an output.

## The red runs

Each red run is a run of the new tests on a skeleton of the code: the steps
S12, S13 and S15 `not-active`, `not built yet`; the register, the commands and
the hook of S15 empty; no wiring of the three steps in `Steps`; the runner
with no call of the hook; the job reader and check `jobs` with no result; each
row `gate:<kind>` `not-active`, `not built yet`. A script wrote the skeleton
over the code and wrote the code back from its saved text after the run.

1. **The unit tests**, ten of them, each for its part:

   ```text
   $ go test -count=1 -run 'TestS12|TestS13|TestTheRulePathRegister|TestS15|TestTheRecordsCommit|TestTheRecordsHook|TestTheLastSteps' ./internal/setup/
   --- FAIL: TestS12                  (S12: not-active "not built yet", values []; want done, and the values [["module" "github.com/acme/widget" …
   --- FAIL: TestS13                  (S13: not-active "not built yet" false, values []; want done with the rows [["branch-protection.sha256" …
   --- FAIL: TestTheRulePathRegister  (the register got []; want [[".github/" "" "baseline"] …)
   --- FAIL: TestS15                  (S15: not-active "not built yet" false, values []; want done, no commit of the runner, and [["setup.head" …
   --- FAIL: TestTheRecordsCommit     (the records commit: files map[], calls []; want … "/tmp/s/tree/README.md" …)
   --- FAIL: TestTheRecordsHook       (a commit: "done records", the hook got 0 rows ending []; want … the value row and the done row)
   --- FAIL: TestTheLastSteps         (the evidence of S12 []; want jobs and gates, S13 a hand-off with commands, and S15 a command and the hook)
   $ go test -count=1 -run 'TestJobNames|TestCheckJobs|TestTheGateRows' ./internal/verify/
   --- FAIL: TestTheGateRows          (detected: {"gate:static" "not-active" "not built yet"} after 0 runs; want {"gate:static" "pass" ""} after 2)
   --- FAIL: TestJobNames             (an id and a name: []; want ["static" "test"])
   --- FAIL: TestCheckJobs            (no job of layout: []; want ["no CI job for the kind layout"])
   ```

   The first skeleton run did not build (imports with no use), and the second
   panicked in `TestTheRecordsHook` and `TestTheGateRows` at an index of an
   empty slice, which stopped the other tests of the package. The tests now
   guard each index, so a wrong result fails a test and does not stop the
   run.

2. **The integration and e2e tests**, on the same skeleton:

   ```text
   $ go test -count=1 -tags=integration -run TestSetupRunsS05ToS15 ./internal/cli/
   --- FAIL: TestSetupRunsS05ToS15     (the run with the answers: exit 1 … S12 not-active not built yet …)
   $ go test -count=1 -tags=integration -run 'TestRunOnAStandInWorkArea|TestTheGateRowsOnARealTarget' ./internal/verify/
   --- FAIL: TestRunOnAStandInWorkArea      (the rows … "gate:static not-active not built yet" …)
   --- FAIL: TestTheGateRowsOnARealTarget   (Check(gates): <nil>, [{"gate:static" "not-active" "not built yet"} …]; want gate:static pass)
   $ go test -count=1 -tags=e2e -run 'TestSetupRunsS01ToS15|TestSetupVerifyOnAStandInWorkArea' ./cmd/layup/
   --- FAIL: TestSetupRunsS01ToS15               (the run with the inputs: exit 1 …)
   --- FAIL: TestSetupVerifyOnAStandInWorkArea   (exit 1 …)
   ```

3. **The text of the records README**, before `setup.md` held it:

   ```text
   $ go test -count=1 -tags=integration -run TestTheRecordsReadme ./internal/setup/
   --- FAIL: TestTheRecordsReadmeIsTheTextOfSetupMd   (the block records-readme of setup.md: an empty text; want # The records of this repository …)
   ```

## What the first green runs found

- **The stand-in fixture failed for another reason.** Before the first run of
  the rows `gate:<kind>` on the stand-in work area, the author read its
  manifest: `go vet ./...` on a tree with no `go.mod`, so the static fixture of
  the Go entry (a file that `gofmt` changes) would fail with "go.mod file not
  found", and `gate:static` would pass for that reason. The stand-in now has a
  `go.mod`, the command of the Go entry, and a workflow with a job per kind;
  the lesson is in `docs/guardrails.md` §2.
- **The git date format.** The check of the records commit compared
  `%aI` with `+00:00`; git 2.54 gives `Z`. The test now compares `%at` and
  `%ct`; the lesson is in `docs/guardrails.md` §2.
- **The author's own expected values**, each fixed in the test, not the code:
  the step numbers of the progress lines of the e2e test (`[7/9] S12`, not
  `[6/9]`); the row of S12 in `open-gaps.tsv` at the head, which the check of
  S11 in the CLI test now holds; two scratch directories in the unit test of
  the gate rows, as each row makes its own fixture commit.
- **A link to an anchor with no heading.** LAYUP's `link-lint.sh` refused the
  link to an `<a id>` anchor in `setup.md`; the text of the README has its own
  heading now.

## The mutations

Each mutation is one change of the code, run against the tests of its package
at the named levels; the file is written back from its text before the
change. All 33 are detected. M17 did not build in the first run (a variable
with no use), and a mutation that keeps it in use is detected; M25 was not
detected in the first run, and a new case of `TestJobNames` (a `name:` of a
deeper map) detects it; M20 was detected by the unit test only, and the CLI
test now checks that no scratch work tree stays.

| Mutation | Detected by |
| -------- | ----------- |
| M1 S12 writes the module path with no host | `TestS12`; `TestSetupRunsS05ToS15` (cli, integration) |
| M2 S12 takes a file of the head with other bytes | `TestS12` |
| M3 S12 writes the gap row again | `TestS12` |
| M4 S12 gives the manifest the ref of a file | `TestS12`; `TestSetupRunsS05ToS15` |
| M5 S12 writes no gap row | `TestS12`; `TestSetupRunsS05ToS15` |
| M6 S13 pins the checks to another app | `TestS13` |
| M7 S13 leaves out the ref `layup-probe` | `TestS13` |
| M8 S13 makes the checks not strict | `TestS13` |
| M9 S13 applies the ruleset before the push of `layup-setup` | `TestS13`; `TestSetupRunsS05ToS15` |
| M10 S13 requires one approval | `TestS13` |
| M11 the register keeps the first source of a path listed twice | `TestTheRulePathRegister` |
| M12 the register has no exception for the guardrails | `TestTheRulePathRegister`, `TestS15`; `TestSetupRunsS05ToS15` |
| M13 the register gives a `.sh` file of the catalog the source `baseline` | `TestTheRulePathRegister`, `TestS15`; `TestSetupRunsS05ToS15` |
| M14 S15 takes a row `not-active` | `TestS15` |
| M15 S15 fails a `verify.tsv` of another form | `TestS15` |
| M16 S15 commits a `verify.tsv` that changed | `TestTheRecordsCommit` |
| M17 S15 takes a records branch with another message | `TestTheRecordsCommit` |
| M18 S15 takes a records branch with an executable file | `TestTheRecordsCommit` |
| M19 S15 makes the records commit on `layup-setup` | `TestTheRecordsCommit`; `TestSetupRunsS05ToS15` |
| M20 S15 does not remove its scratch tree | `TestTheRecordsCommit`; `TestSetupRunsS05ToS15` |
| M21 the runner gives the hook the record with no done row | `TestTheRecordsHook`; `TestSetupRunsS05ToS15` |
| M22 the runner writes the done row when the hook fails | `TestTheRecordsHook` |
| M23 the job reader takes no name | `TestJobNames`, `TestCheckJobs` |
| M24 the job reader keeps the carriage return | `TestJobNames` |
| M25 the job reader takes a name of a deeper map | `TestJobNames` |
| M26 check `jobs` passes with no manifest | `TestCheckJobs` |
| M27 a fixture run that passes is a pass | `TestTheGateRows`; `TestTheGateRowsOnARealTarget` (verify, integration) |
| M28 a fixture that does not apply is not `not-active` with its reason | `TestTheGateRows`; `TestTheGateRowsOnARealTarget` |
| M29 the clean run is made for each kind | `TestTheGateRows` |
| M30 a clean fail does not decide | `TestTheGateRows` |
| M31 the fixture commit is dated by no `pin.time` | `TestTheGateRows` |
| M32 the fixture work tree stays | `TestTheGateRows`, `TestRunGivesEachRowInTheOrderOfTheTable`; `TestRunOnAStandInWorkArea` (verify, integration) |
| M33 a pending kind runs the gate | `TestTheGateRows`, `TestRunGivesEachRowInTheOrderOfTheTable` |

## The green runs

On `dde2e66`, 00:41:41Z to 00:43:46Z; the commits after it add the check of
the scratch trees to `TestSetupRunsS05ToS15` and text:

| Command | Result |
| ------- | ------ |
| The local checks of `AGENTS.md`, with `git diff --check origin/main HEAD`; `go build ./...`; `go vet` with no tag, `-tags=integration` and `-tags=e2e`; `gofmt -l .` | exit 0; no file |
| `go test -count=1 ./...`; with `-tags=integration`; with `-tags=e2e -timeout 10m` | exit 0 each |
| `go test -race -count=1` on `internal/catalog`, `internal/verify`, `internal/work`, `internal/standin`, `internal/setup` and `internal/cli`, and on `internal/catalog`, `internal/verify`, `internal/setup` and `internal/cli` with `-tags=integration` | exit 0 |
| `sh docs/setup/tests/run.sh`; `sh docs/tests/run-discipline-tests.sh` | 44 passed, 0 failed; 81 passed, 0 failed |

On the freeze head `a7a64c9`, 00:47:41Z to 00:49:44Z, the same 18 steps, each
exit 0 (`run.sh` 44 passed, the discipline tests 81 passed); review round 1
ran the three test levels, the race tests and the local checks on it again,
with the same results.

The demo of the plan, in `TestSetupRunsS05ToS15` through `internal/cli`: on a
stand-in baseline, S12 writes the files of the Go entry and passes checks
`jobs` and `gates` (the gate gives `clear` on the setup head, and `fail` on the
commit of each fixture: `gofmt` lists the static fixture, and the test of the
test fixture fails); S13 commits `docs/setup/branch-protection.json` and
writes `out/ruleset-default.json` (`operator`); S15 stops for `O-verify`; with
the table of `layup setup verify` (exit 0), S15 makes the records commit (no
parent, the four files, `setup/record.tsv` equal to `out/record.tsv`), `main`
and `layup-setup` do not move, no ref and no work tree stays, `commands.sh`
holds the four commands in the order of `setup.md`, and a rerun gives the
same table with no second records commit.
