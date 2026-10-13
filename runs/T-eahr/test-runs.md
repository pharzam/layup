# T-eahr — test runs

The step setup (#277, row 54). All times are UTC, 2026-10-13, read from the clock after each step.

## The order of the gate

The claim was pushed at 10:55:00Z and the plan posted at 11:03:55Z; its review (GPT-6.1 Sol on Codex CLI, `high`) ran from 11:04:21Z to
11:13:15Z, with the author's answer at 11:15:09Z. The base was run at 11:15:23Z and the first red run is at 11:18:20Z: the tests were written
after the confirmation.

## The base (11:15:23Z to 11:16:07Z)

At `b28171b`: build, vet and `go test ./...` exit 0; the integration run of four packages exited 1 on `TestProcessTimesTheFirstOutput`,
a timing race in a test that this task does not touch (a fourth cause on #291): three runs of it alone passed, and the four packages passed again.

## Red

Each slice has its tests first; the package does not build until its code exists. S1, the answers (11:18:20Z): `undefined: setupData`,
`setupAnswers`. S2, the work area (11:22:31Z): `undefined: AddRemote`, `workSource`, `buildWorkArea`. S3, the output of files and the removal
(11:25:54Z): `undefined: ReadOutputFiles`, `s.Files undefined (type PhaseSpec ...)`, `proseNames`. S4, the loop (11:34:49Z): `undefined: setupLoop`,
`proseDone`; at 11:39:34Z the four sites of the tests that end at `follow-up-answers`, edited first, fail (`want ... setup`).

## Written after the code

The test of `ReadFiles` of the store (`TestTheFilesOfTheRecordsAreReadByTheirPaths`) was written after its code, and the cases that the mutants
asked for (the answer `Gap`, the order of `Q-999` and `Q-1000`, an S13 handed off, a name `..` and `.`, a file where a directory is named, a path
two directories deep, a URL with a delete character, an `M-` row that names itself) after the first run of the mutants.

## Turning green

The run of the worlds found four things, put right in the tests: the glossary that was meant to fail its check had no citation to fail on
(`F-0099#1` resolves to no fact); the pair of `product-owner` is taken out of the host's routing register, not of the records; the stand-in
baseline had two commits and made `read-back` fail in one run of four (the pitfall of `guardrails.md`); and Intake ends in a fail at S04 for a file the baseline lacks.

## Round 1 and round 2: the fixes

Each test was written and run red before its code. Round 1 (comment 6146226165, five findings): 1, a session that ends after its commit
(13:07:05Z, 13:07:49Z: `removed []`, and "docs/guide.md is in the records after the fail"); 2, a link at `result/` (13:07:58Z: no error); 3, a link
above the work area or in the records base (13:08:14Z, 13:08:32Z: no error, and the file outside is gone); 4, the file-path test moved to the
integration level (green at once); 5, words of `records.md` only. Round 2 (comment 6147313832, one finding; the cap raised to 2 at 14:04:30Z):
the prompt of an adapted file (14:05:14Z: the prompt lacks "A row whose ask says Adapt" and the three phrases after it; green at 14:05:23Z).

## Mutations

Each rule of the answers, the loop, the prose session and its prompt, the work area, `AddRemote`, `ReadOutputFiles`, the session of a phase
with files, the store's removal and reading, the steps and the step setup mutated alone, on a copy of the tree and never on the worktree: a
condition inverted or dropped, a name, a path or a field changed, so each mutant builds. A mutant that does not build is counted as not run.

- **Run 1** (12:07:25Z to 12:27:26Z): 87 caught, 3 survived, 13 not run, of 103. Survived: `out-above-dir`, `out-deep` (the error of a file where a directory is named was that of a missing file, and the directories above a name were seen through another name) and `run-steps-enum` (its command was the unit level, which does not compare the schema with its block). Thirteen did not build (an import or a variable left unused), so none counted as caught: each was written again to build.
- **Run 2** (12:29:22Z to 12:31:14Z): 15 caught, 1 survived, 0 not run, of 16. The thirteen and the three. `remove-valid` survived: git refuses the empty commit that follows a removal of nothing, so the call failed all the same; a test now names a file two directories above the work tree.
- **Run 3** (12:31:15Z to 12:31:24Z): 1 caught, 0 survived, 0 not run, of 1. `remove-valid`, with the new test.
- **Run 4** (13:14:06Z to 13:17:13Z): 22 caught, 0 survived, 1 not run, of 23. The rules of the fixes of round 1 and the five of the file path, now at the integration level. `remove-tree-check` did not build (`entries` left unused) and was written again; `remove-missing` is gone, as the tree check refuses the absent file first.
- **Run 5** (13:20:01Z to 13:22:50Z): 23 caught, 0 survived, 0 not run, of 23. The same 23, on `1086261`.
- **Run 6** (14:05:58Z to 14:06:20Z): 8 caught, 0 survived, 0 not run, of 8. The prompt of an adapted file (round 2) and the three mutants of its other lines.

## Green

On `89d5c0c` (14:09:15Z to 14:19:29Z), in this order, each with exit 0 but the integration level (the time is when the command ended; `go test` with `-count=1`, then
`-tags=integration` and `-tags=e2e`, `TMPDIR` on the disk, `-timeout 30m`): go build ./... 14:09:16Z; go vet ./... 14:09:16Z; gofmt -l internal cmd (empty) 14:09:17Z; go test ./... 14:09:23Z; go test -tags=integration ./...
14:17:19Z EXIT 1; go test -tags=e2e ./... 14:18:36Z; adr-lint 14:18:37Z; prd-lint 14:18:37Z; link-lint 14:18:52Z; run-discipline-tests 14:18:59Z;
nested-checkout-check 14:19:00Z; setup-check 14:19:29Z; git diff --check origin/main HEAD 14:19:29Z.
The integration level exited 1 on `TestProcessTimesTheFirstOutput` alone (`internal/session`, the timing race of #291: the end 0.19951 s after the first byte, where the test wants
0.2 s; the test of the base too): that package run again alone exited 0 at 14:19:48Z. `internal/run` took 453.9 s of the 600 that CI's `-timeout 10m` gives a package (355 s before
this task); CI's own job takes 4 minutes.
