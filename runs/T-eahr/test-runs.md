# T-eahr — test runs

The step setup (#277, row 54). All times are UTC, 2026-10-13, read from the clock after each step.

## The order of the gate

The claim was pushed at 10:55:00Z and the plan posted at 11:03:55Z. Its review (GPT-6.1 Sol on Codex CLI, `high`) ran from
11:04:21Z to 11:13:15Z and was posted at 11:15:04Z, with the author's answer at 11:15:09Z. The base was run at 11:15:23Z and the
first red run is at 11:18:20Z: the tests were written after the confirmation.

## The base (11:15:23Z to 11:16:07Z)

At `b28171b`. The integration run of four packages exited 1 on `TestProcessTimesTheFirstOutput`, a timing race in a
test that this task does not touch (a fourth cause on #291); three runs of it alone passed, and the four packages passed again:

```
go build ./...: exit 0 (11:15:30Z)
go vet ./...: exit 0 (11:15:33Z)
go test ./...: exit 0 (11:15:54Z)
go test -tags=integration ./internal/session/ ./internal/git/ ./internal/setup/ ./internal/work/: exit 1 (11:16:07Z)
go test -tags=integration ./internal/session/ ./internal/git/ ./internal/setup/ ./internal/work/ (again): exit 0 (11:17:42Z)
```

## Red

Each slice has its tests first; the package does not build until its code exists:

- **S1, the answers** (11:18:20Z): `undefined: setupData`, `undefined: setupAnswers`.
- **S2, the work area** (11:22:31Z to 11:22:34Z): `undefined: AddRemote`, `workSource`, `buildWorkArea`.
- **S3, the output of files and the removal** (11:25:54Z to 11:25:56Z): `undefined: ReadOutputFiles`, `s.Files undefined (type PhaseSpec ...)`,
  `undefined: proseNames`.
- **S4, the loop** (11:34:49Z): `undefined: setupLoop`, `undefined: proseDone`; at 11:39:34Z the four sites of the tests that end at
  `follow-up-answers`, edited first, fail (`want ... setup`; `undefined: setupRunner`).

## Written after the code

The test of `ReadFiles` of the store (`TestTheFilesOfTheRecordsAreReadByTheirPaths`) was written after its code, and the cases that
the mutants asked for (the answer `Gap`, the order of `Q-999` and `Q-1000`, an S13 handed off, a name `..` and `.`, a file where a
directory is named, a path two directories deep, a URL with a delete character, an `M-` row that names itself) after the first run of the mutants.

## Turning green

Four things were found by the run of the tests of the world and put right in the tests: the glossary that was meant to fail its check had no
citation to fail on (`F-0099#1` resolves to no fact); the pair of `product-owner` is taken out of the host's routing register, which the step
`probe` copies into the records at each restart; the stand-in baseline had two commits and made `read-back` fail in one run of four (a baseline
of one commit, as the pitfall of `guardrails.md` says); and the world that Intake ends in is fail at S04 for the file `docs/adr/README.md` that its baseline lacks.

## Round 1: the fixes

Round 1 (comment 6146226165, five findings, posted 13:06:06Z): each test was written and run red before its code.

- **1, a session that ends after its commit** (13:07:05Z, 13:07:49Z): the two new rows of `TestAFailAfterAProseSession...` fail with
  `removed []`, and `TestAProseSessionWhoseCommentFailsLeavesNoFileInTheRecords` with "docs/guide.md is in the records after the fail".
- **2, a link at `result/`** (13:07:58Z): `TestReadOutputFiles` gets no error for a link to a directory that holds the output.
- **3, a link above the work area or in the records base** (13:08:14Z, 13:08:32Z): `TestTheWorkAreaIsNotBuiltThroughALink` gets no error and the
  file outside is gone; so is it in `TestARecordsCommitRemovesNoFileThroughALink`.
- **4**, the test of a file path moved to the integration level, green at once (13:08:14Z); **5**, words of `records.md` only.

## Mutations

Each rule of the answers, the loop, the prose session, the work area, `AddRemote`, `ReadOutputFiles`, the session of a phase with files, the
store's removal and reading, the steps and the step setup mutated alone, on a copy of the tree and never on the worktree: a condition inverted or
dropped, a name, a path or a field changed, so each mutant builds; then the tests that can see it. A mutant that does not build is counted as
not run, never as caught.

- **Run 1** (12:07:25Z to 12:27:26Z): 87 caught, 3 survived, 13 not run, of 103. Survived: `out-above-dir` and `out-deep` (the error of a file where a directory is named was that of a missing file, and the directories above a name were seen through another name) and `run-steps-enum` (its command was the unit level, which does not compare the schema with its block). Thirteen did not build (an import or a variable left unused), so none counted as caught: each was written again to build.
- **Run 2** (12:29:22Z to 12:31:14Z): 15 caught, 1 survived, 0 not run, of 16. The thirteen and the three. `remove-valid` survived: git refuses the empty commit that follows a removal of nothing, so the call failed all the same; a test now names a file two directories above the work tree.
- **Run 3** (12:31:15Z to 12:31:24Z): 1 caught, 0 survived, 0 not run, of 1. `remove-valid`, with the new test.
- **Run 4** (13:14:06Z to 13:17:13Z): 22 caught, 0 survived, 1 not run, of 23. The rules of the fixes, and the five of the file path, now at the integration level. `remove-tree-check` did not build (`entries` left unused) and was written again; `remove-missing` is gone, as the tree check refuses the absent file first.
- **Run 5** (13:20:01Z to 13:22:50Z): 23 caught, 0 survived, 0 not run, of 23. The same 23 on `1086261`, the code that round 2 reads.

## Green

On `10c2ac9` (13:35:09Z to 13:45:36Z), each with exit 0, in this order; the time of a line is when its command ended. The runs of
`go test` are with `-count=1`, then with `-tags=integration` and with `-tags=e2e`, `TMPDIR` is on the disk, and `git diff --check` reads the commits against `origin/main`.
An earlier run on `1086261` (13:23:03Z to 13:33:37Z) had the same exits; `10c2ac9` after it changes a comment of `store.go` and two cells of the PRD, and the code of the mutants is that of `1086261`.

```
go build ./...: exit 0 (13:35:10Z)
go vet ./...: exit 0 (13:35:11Z)
gofmt -l internal cmd (empty): exit 0 (13:35:11Z)
go test ./...: exit 0 (13:35:17Z)
go test -tags=integration ./...: exit 0 (13:43:17Z)
go test -tags=e2e ./...: exit 0 (13:44:44Z)
adr-lint: exit 0 (13:44:45Z)
prd-lint: exit 0 (13:44:45Z)
link-lint: exit 0 (13:45:00Z)
run-discipline-tests: exit 0 (13:45:07Z)
nested-checkout-check: exit 0 (13:45:08Z)
setup-check: exit 0 (13:45:35Z)
git diff --check origin/main HEAD: exit 0 (13:45:36Z)
```

The integration level of `internal/run` took 458 s of the 600 that CI's `-timeout 10m` gives a package (355 s before this task: the worlds
of the step setup cost about 100 s on this host, which was quiet); CI's own `tests` job runs the three levels in 4 minutes. The run here used `-timeout 30m`.