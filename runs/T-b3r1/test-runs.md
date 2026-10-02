# T-b3r1: the test runs

The runs of `T-b3r1` (#90). Host: macOS, `go1.27.1`, `git` 2.54.0,
2026-10-02, UTC. "…" marks a cut part of an output; `\u2039` and `\u203a` are the
two angle quotes of a marker, written as escapes so that check `markers` reads
no marker in this file.

## The measurements of the plan

On LAYUP's root commit `d2516fd` (the baseline at its pin): `docs/tasks/` holds
19 task files; the 7 lines of `backlog.md` that link `github.com/pharzam/armature`
name none of the 19 IDs; of the 56 such lines of `completed.md`, 32 name none
and 24 name one; the one blockquote with such a line is the pivot note of
`backlog.md` (lines 28 to 31). The baseline's `docs/links/link-lint.sh` is the
same file in LAYUP (`git diff d2516fd HEAD -- docs/links/link-lint.sh` is
empty). On a stand-in baseline, that script refuses a `file://` link ("links
file:///…/issues/1, but that path does not exist"), so the lines of the task
indexes that link the stand-in's URL break it until S05 removes them, and a
guide that links `decisions/D-0001-stand-in.md` breaks after the deletion.

## The red runs

1. **`internal/git`**, on skeletons of `ResetHard` and `Staged` (no start of
   `git`; `false`, no error):

   ```text
   $ go test -count=1 ./internal/git/
   --- FAIL: TestEachCallRunsItsVerb/reset_--hard            (0 starts of git, want 1)
   --- FAIL: TestEachCallRunsItsVerb/diff_--cached_--quiet   (0 starts of git, want 1)
   --- FAIL: TestStagedReadsTheExitCode   (exit 1: false, <nil>; want true …; exit 128: false, <nil>; want false and an error true)
   ```

2. **`internal/verify`**, on skeletons (the column 0, `TextMarkers` nil,
   `LinksBaseline` false, no line of the records branch):

   ```text
   $ go test -count=1 ./internal/verify/
   --- FAIL: TestLinksBaseline        (LinksBaseline("https://github.com/pharzam/armature", "- **T-1** ([#1](…/issues/1))") = false, want true; …)
   --- FAIL: TestTheChecksOfATarget   (identity with no records branch: …)
   --- FAIL: TestScanMarkers          (the markers: {docs/a.md 1 0 \u2039x\u203a} …; want {docs/a.md 1 4 \u2039x\u203a} …)
   --- FAIL: TestTheMarkersOfAText    (TextMarkers: []; want 4 markers with their columns)
   ```

3. **`internal/setup`**, on skeletons of S05, S06, the prose step, S14, S10 and
   S11 (each `not-active`, `not built yet`): eight tests fail, each for its
   step.

   ```text
   $ go test -count=1 ./internal/setup/
   --- FAIL: TestS05, TestS05StopsForTheLinksThatBreak, TestS06, TestTheProseStep, TestS14, TestS10, TestS11, TestS11ChecksBeforeItWrites
   ```

   `TestTheRunnerResetsTheTargetBeforeAStep` came with the code of the runner
   (D11); its evidence is the mutation M20 below.

4. **`internal/cli`**, on a skeleton of `calls` (the one-check call only) and of
   the brief reader (no check of a marker or of the vision brief):

   ```text
   $ go test -count=1 ./internal/cli/
   --- FAIL: TestSetupReadsTheProblemStatement   (a marker in the problem statement, a vision brief that is not UTF-8,
                                                  a marker in the vision brief: exit 0 …; want 2 and the reason)
   --- FAIL: TestTheCallsOfTheSteps              (panic: a nil call)
   ```

5. **The integration and e2e tests**, on skeletons of the six steps:

   ```text
   $ go test -count=1 -tags=integration -run TestSetupRunsS05ToS14 ./internal/cli/
   --- FAIL: TestSetupRunsS05ToS14   (the stop of S05: exit 1 …)
   $ go test -count=1 -tags=e2e -timeout 10m -run TestSetupRunsS01ToS14 ./cmd/layup/
   --- FAIL: TestSetupRunsS01ToS14   (the run with the answers: exit 1, … S05 not-active not built yet …)
   ```

   The first green run of `TestSetupRunsS05ToS14` failed on the author's own
   expected table: the runner sorts the `F-` rows by path, and the test had
   `docs/how-to.md` before `docs/glossary.md`; the test was fixed, not the
   code.

6. **Two angle quotes as characters.** Before the first commit, a byte search
   found the two characters in `internal/setup/scaffold.go` (the check of S11)
   and in `internal/cli/setup_integration_test.go`: the tool call had turned
   their escapes into the characters (a `grep` with the escape found nothing).
   Each is an escape now, and LAYUP's check `markers` passes; the lesson is in
   `docs/guardrails.md` §2.

## The mutations

Each mutation is one change of the code, run against the tests of its package
at the named levels; the file is written back from its text before the
change, and `cmp` showed each file unchanged after the run. All 27 are
detected.

| Mutation | Detected by |
| -------- | ----------- |
| M1 S05 keeps a line that names a deleted task | `TestS05` |
| M2 S05 keeps a line that links the baseline | `TestS05`; `TestSetupRunsS01ToS04`, `TestABrokenPinFailsTheEvidenceOfS04`, `TestSetupRunsS05ToS14` (cli, integration) |
| M3 S05 removes one line of a blockquote, not the whole blockquote | `TestS05`; `TestSetupRunsS05ToS14` |
| M4 S05 reads a task ID inside a longer word | `TestS05` |
| M5 S05 stops for a file that has its input | `TestS05`, `TestS05StopsForTheLinksThatBreak`, `TestTheProseStep`, `TestS14`; `TestSetupRunsS05ToS14` |
| M6 S05 does not ask again after the copy | `TestS05StopsForTheLinksThatBreak` |
| M7 S06 does not copy the vision brief | `TestS06` |
| M8 S06 dates its index rows by the clock | `TestS06` |
| M9 S14 asks again for the files of S07 to S09 | `TestS14` |
| M10 S14 does not ask for a flagged file | `TestS14`; `TestSetupRunsS05ToS14` |
| M11 a copied file gets no record row | `TestS05`, `TestTheProseStep`, `TestS14`; `TestSetupRunsS05ToS14` |
| M12 a step copies its inputs before it knows that one is missing | `TestS14` |
| M13 S10 names the last line of a marker | `TestS10`; `TestSetupRunsS05ToS14` |
| M14 S10 takes an answer to no marker | `TestS10` |
| M15 S11 fills a marker one byte after its column | `TestS11`; `TestSetupRunsS05ToS14` |
| M16 S11 writes a row of `open-gaps.tsv` per place | `TestS11` |
| M17 S11 takes a value with an angle quote | `TestS11ChecksBeforeItWrites` |
| M18 S11 writes no second answers record | `TestS11`; `TestSetupRunsS05ToS14` |
| M19 S11 writes the `Source` of S04's record | `TestS11` |
| M20 the runner does not reset the target before a step | `TestTheRunnerResetsTheTargetBeforeAStep`; `TestEachStepStartsFromTheHead` (setup, integration) |
| M21 a step that changes no file makes a commit | `TestEachStepStartsFromTheHead` |
| M22 `internal/cli` takes a brief that holds a marker | `TestSetupReadsTheProblemStatement` |
| M23 `internal/cli` does not read the vision brief | `TestSetupReadsTheProblemStatement` |
| M24 `internal/cli` drops the column of a marker | `TestTheCallsOfTheSteps`; `TestSetupRunsS05ToS14` |
| M25 the scanner gives the column of the end of the open quote | `TestScanMarkers`, `TestTheMarkersOfAText`; `TestSetupRunsS05ToS14` |
| M26 check `identity` does not read the records branch | `TestTheChecksOfATarget`; `TestEachFindingOfATarget` (verify, integration); `TestSetupRunsS05ToS14` |
| M27 the link rule of the baseline links nothing | `TestLinksBaseline`; `TestSetupRunsS01ToS04`, `TestABrokenPinFailsTheEvidenceOfS04`, `TestSetupRunsS05ToS14` |

M9 was not detected by `TestSetupRunsS05ToS14`: its stand-in baseline has no
file of S07 to S09, so check `adapted` flags none of them, and no row was twice
in the table; the unit test `TestS14` holds the case.

## The green runs

On the tree of the commit before the one that adds this file (`af7b218`; this
commit adds only text), 23:02:19Z to 23:04:13Z:

| Command | Result |
| ------- | ------ |
| The local checks of `AGENTS.md`, with `git diff --check b812943 HEAD`; `go build ./...`; `go vet` with no tag, `-tags=integration` and `-tags=e2e`; `gofmt -l .` | exit 0; no file |
| `go test -count=1 ./...` | `ok` × 11 packages |
| `go test -count=1 -tags=integration ./...` | `ok` × 11 |
| `go test -count=1 -tags=e2e -timeout 10m ./...` | `ok` × 11 |
| `go test -race -count=1` on `internal/catalog`, `internal/verify`, `internal/work`, `internal/standin`, `internal/setup` and `internal/cli`, and on `internal/catalog`, `internal/verify`, `internal/setup` and `internal/cli` with `-tags=integration` | `ok` |
| `sh docs/setup/tests/run.sh`; `sh docs/tests/run-discipline-tests.sh` | 44 passed, 0 failed; 81 passed, 0 failed |
| The first commit (`344ecce`) alone, in a scratch work tree: `go vet` with each tag; the unit tests; the integration tests of `internal/cli`, `internal/setup`, `internal/verify` and `internal/git`; the e2e tests of `cmd/layup` | `ok` |
