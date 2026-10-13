# T-eahr — test runs

The step setup (#277, row 54): the answers of the setup, the work area, the loop over the stops of the runner and
the prose session, in `internal/run`, with the output of files of `internal/session`, the removal of files in the
records store and `AddRemote` of `internal/git`. All times are UTC, 2026-10-13, read from the clock after each step.

## The order of the gate

The claim was pushed at 10:55:00Z, the plan was posted at 11:03:55Z, its review (GPT-6.1 Sol on Codex CLI, effort
`high`) ran from 11:04:21Z to 11:13:15Z and was posted at 11:15:04Z, and the author's answer at 11:15:09Z. The base
was run at 11:15:23Z and the first red run is at 11:18:20Z: the tests were written after the confirmation.

## The base (11:15:23Z to 11:16:07Z)

At `b28171b`. The integration run of four packages exited 1 on `TestProcessTimesTheFirstOutput`, a
timing race in a test that this task does not touch (the end 0.19995 s after the first byte where the test wants
0.2 s; a fourth cause on #291); three runs of it alone passed, and the four packages passed again:

```
go build ./...: exit 0 (11:15:30Z)
go vet ./...: exit 0 (11:15:33Z)
go test ./...: exit 0 (11:15:54Z)
go test -tags=integration ./internal/session/ ./internal/git/ ./internal/setup/ ./internal/work/: exit 1 (11:16:07Z)
go test -tags=integration ./internal/session/ ./internal/git/ ./internal/setup/ ./internal/work/ (again): exit 0 (11:17:42Z)
```

## Red

Each slice has its tests first; the package does not build until its code exists, and the first error is the name
of what is missing. The files of a slice, and its run:

- **S1, the answers** (11:18:20Z): `undefined: setupData`, `undefined: setupAnswers` (`go vet ./internal/run/`).
- **S2, the work area** (11:22:31Z to 11:22:34Z): `undefined: AddRemote` (`internal/git`, integration), `undefined: workSource` and
  `undefined: buildWorkArea` (`internal/run`).
- **S3, the output of files and the removal** (11:25:54Z to 11:25:56Z): `undefined: ReadOutputFiles` (`internal/session`,
  integration), `s.Files undefined (type PhaseSpec ...)`, and in the unit level `undefined: proseNames`.
- **S4, the loop** (11:34:49Z): `undefined: setupLoop`, `undefined: proseDone`; and at 11:39:34Z the four sites of the tests that
  end at `follow-up-answers`, edited first: `TestStartAndRestartCallStepAtTheStartOfEachStep` fails with "the steps of a
  restart with Start done [... follow-up-answers] want ... setup", and the integration tests of the step do not build
  (`undefined: setupRunner`).

## Written after the code

The test of `ReadFiles` of the store (`TestTheFilesOfTheRecordsAreReadByTheirPaths`) was written after its code, and
the cases that the mutants asked for (the answer `Gap`, the order of `Q-999` and `Q-1000`, an S13 handed off, a name
`..` and `.`, a file where a directory is named, a path two directories deep, a URL with a delete character, an `M-`
row that names itself) after the first run of the mutants.

## Turning green

Four things were found by the run of the tests of the world and put right in the tests: the glossary that was meant to
fail its check had no citation to fail on (the check reads citations: `F-0099#1` resolves to no fact); the pair of
`product-owner` is taken out of the host's routing register, which the step `probe` copies into the records at each
restart, and not out of the records; the stand-in baseline had two commits and made `read-back` fail in one run of
four (a baseline of one commit, as the pitfall of `guardrails.md` says); and the world that Intake ends in is fail at
S04 for the file `docs/adr/README.md` that its baseline lacks.

## Mutations

Each rule of the answers, the loop, the prose session, the work area, `AddRemote`, `ReadOutputFiles`, the session of a
phase with files, the store's removal and reading, the steps and the step setup mutated alone, on a copy of the tree
and never on the worktree: a condition inverted or dropped, a name, a path or a field changed, so each mutant builds;
then the tests that can see it (the unit level of `internal/run`, or the named tests of the integration level; the
mutants of the step itself run its test of the world). A mutant that does not build is counted as not run, never as caught.

- **Run 1** (12:07:25Z to 12:27:26Z): 87 caught, 3 survived, 13 not run, of 103. Three survived: `out-above-dir` and `out-deep` (the error of a file where a directory is named was the error of the file that its absence made missing, and the directories above a name were seen through another name) and `run-steps-enum` (its command was the unit level, which does not compare the schema with its block; the integration test does). Thirteen did not build (an import or a variable that the mutation left unused), so none was counted as caught: each is written again to build, with a fresh copy of the merged tree.
- **Run 2** (12:29:22Z to 12:31:14Z): 15 caught, 1 survived, 0 not run, of 16. The sixteen are the thirteen and the three of run 1. One survived: `remove-valid`, the check of the paths of the removal of files (git refuses the empty commit that follows a removal of nothing, so the call failed all the same); a test now names a file two directories above the work tree.
- **Run 3** (12:31:15Z to 12:31:24Z): 1 caught, 0 survived, 0 not run, of 1. The one mutant of run 2 that survived, with the new test.

- **20** caught by `setupanswers_test.go`
- **18** caught by `setup_test.go`
- **12** caught by `prose_test.go`
- **12** caught by `setup_integration_test.go`
- **11** caught by `phase_integration_test.go`
- **10** caught by `workarea_integration_test.go`
- **5** caught by `git_integration_test.go`
- **5** caught by `workarea_test.go`
- **4** caught by `a panic or a failing test of the level`
- **4** caught by `store_integration_test.go`
- **1** caught by `command_test.go`
- **1** caught by `run_integration_test.go`

## Green

On `5076512` (12:32:03Z to 12:42:26Z), each with exit 0, in this order; the time of a line is when its command ended. The three runs of
`go test` are with `-count=1`, then with `-tags=integration` and with `-tags=e2e`, `TMPDIR` is on the disk, and
`git diff --check` reads the commits against `origin/main`:

```
go build ./...: exit 0 (12:32:04Z)
go vet ./...: exit 0 (12:32:05Z)
gofmt -l internal cmd (empty): exit 0 (12:32:05Z)
go test ./...: exit 0 (12:32:11Z)
go test -tags=integration ./...: exit 0 (12:40:03Z)
go test -tags=e2e ./...: exit 0 (12:41:33Z)
adr-lint: exit 0 (12:41:35Z)
prd-lint: exit 0 (12:41:35Z)
link-lint: exit 0 (12:41:49Z)
run-discipline-tests: exit 0 (12:41:57Z)
nested-checkout-check: exit 0 (12:41:57Z)
setup-check: exit 0 (12:42:26Z)
git diff --check origin/main HEAD: exit 0 (12:42:26Z)
```

The integration level of `internal/run` took 451 s of the 600 that CI's `-timeout 10m` gives a package (355 s before this
task: the worlds of the step setup and the extra row of the restart cost about 100 s on this host, which was quiet); CI's
own `tests` job runs the three levels in less than 3 minutes. The run here used `-timeout 30m`.