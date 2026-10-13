# T-hfn1 — test runs

The step scaffold-push (#279, row 56). All times are UTC, 2026-10-13, read from the clock after each step.

## The order of the gate

The claim was pushed at 16:45:49Z and the plan posted at 16:48:59Z; its review (GPT-6.1 Sol on Codex CLI, `high`) ran from 16:49:02Z to
16:53:00Z and was posted at 16:53:43Z, with the author's answer at 16:53:45Z. The base was run at 16:53:53Z, and the specification was
corrected first (16:55:54Z: the framing of `setup/rulesets.txt`, an R10 correction of the plan review) before the first red run, at 16:56:29Z.

## The base (16:53:53Z to 16:54:15Z)

At `720a149`: build, vet and `go test ./...` exit 0, and the CI of `main` had passed.

## Red

The tests come first and the package does not build until the code exists (16:56:29Z: `undefined: scaffoldJob`). With a skeleton that compiles and does nothing, each fails on its
assertion: the unit level at 16:57:39Z (`restartNames` without `scaffold-push` and the five tests of the record and the step, each with an empty step), and the integration level at
16:59:10Z (the progress counting 19 steps where 20 are wanted, `readApplies` "not built", and the world, whose last row is `verify`).

## Written after the code

The phases of the world that rebuild a work area whose loop fails or whose head is not the last commit, the check that a work area that is there stays, and the case of a null output were
written after the code and before the first run of the mutants; the details of the step as literals and the message of each refusal after it.

## Turning green

The world passed at its first run on the code. Two things of the plan were needed in the tests: the fake forge gives the head of each branch (`commit.sha`), and the world of row 55
no longer asserts that a restart builds no work area, as the step `scaffold-push` builds one that is missing.

## Round 1: the fix

Round 1 (comment 6151696485, one finding): condition 4 of the plan review asked for an output that is a directory and one that cannot be read, and only the directory had a case. The cases of the
file with no permission (in `readApplies` and in the world, with the check that no record of the applies is published) were green at once, as the code names any name that it cannot read; they are skipped for a
user that may read any file. The code is that of the second run of the mutants.

## Mutations

Each rule of the step, its record, its reading of the outputs, its work area and its wiring mutated alone on a copy of the tree and never on the worktree: a condition inverted or dropped,
a name, a path or a field changed, so each mutant builds. A mutant that does not build is counted as not run.

- **Run 1** (17:06:43Z to 17:17:34Z): 31 caught, 4 survived, 3 not run, of 38. Survived: `waiting-detail` and `done-detail` (the tests took the two details from the constants that they test), `record-null` (a null output was refused by the next check with the same error) and `area-rebuild-existing` (the branch of an existing work area returned the same as its neighbour, so it could be dropped with no change of result). Three did not build (a variable left unused). The tests now say the details and the message of each refusal, the second branch of the check is its own error, and the three mutants are written again to build.
- **Run 2** (17:18:58Z to 17:29:02Z): 38 caught, 0 survived, 0 not run, of 38. All 38 on `decab3b`, the code that round 1 reads.

## Green

On `5cb9c94` (17:52:02Z to 18:04:10Z), each with exit 0, in this order (the time is when the command ended; `go test` with `-count=1`, then
`-tags=integration` and `-tags=e2e`, `TMPDIR` on the disk, `-timeout 30m`): go build ./... 17:52:03Z; go vet ./... 17:52:04Z; gofmt -l internal cmd (empty) 17:52:04Z; go test ./... 17:52:11Z; go test -tags=integration ./...
18:02:01Z; go test -tags=e2e ./... 18:03:18Z; adr-lint 18:03:20Z; prd-lint 18:03:20Z; link-lint 18:03:34Z; run-discipline-tests 18:03:42Z; nested-
checkout-check 18:03:43Z; setup-check 18:04:10Z; git diff --check origin/main HEAD 18:04:10Z.
The integration level of `internal/run` took 573 s on this host (510 s at the end of T-79zk: the new world makes about twelve restarts, several of which rebuild the work area); CI's own run of the
integration level of that package took 125 s on `main` at `720a149`, against the 600 s of its `-timeout 10m`, and CI's whole job about 4.5 minutes. The run here used `-timeout 30m`. The commits
after `5cb9c94` change only a comment of a test (the note of round 2) and documents: this record, the task file and the two indexes.
