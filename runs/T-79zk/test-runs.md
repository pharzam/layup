# T-79zk — test runs

The step verify (#278, row 55). All times are UTC, 2026-10-13, read from the clock after each step.

## The order of the gate

The claim was pushed at 14:43:37Z and the plan posted at 14:47:51Z; its review (GPT-6.1 Sol on Codex CLI, `high`) ran from 14:47:57Z to
14:54:02Z and was posted at 14:54:42Z, with the author's answer at 14:55:06Z. The base was run at 14:55:11Z and the first red run is at
14:57:55Z: the tests were written after the confirmation.

## The base (14:55:13Z to 14:55:34Z)

At `6dec793`: build, vet and `go test ./...` exit 0, and the CI of `main` had passed.

## Red

The tests come first and the package does not build until the code exists (14:57:55Z and 14:58:55Z: `undefined: verifyJob`). With a skeleton that compiles and
does nothing, each fails on its assertion: the unit level at 14:59:09Z (the three tests of the step, each with an empty step, and `restartNames` without `verify`), the
integration level at 14:59:40Z (the block `scaffold-readme` empty, the progress counting 18 steps where 19 are wanted, the tree "not built", the step `verify` an empty row) and
S12 in phase 2 at 15:02:39Z (the module row still `answer S01-name`).

## Written after the code

The test of `setupOutputs`, the table that its schema refuses, the module row of a name that cites a fact and the progress lines of the world were written
after the code and before the mutants; the messages of a setup commit of another form and the name of the commit by its message, after the first run of the mutants.

## Turning green

The first real run found a defect of an earlier task: `layup setup verify` refused the record of phase 2, "sources is fail: source: S12 module: the answer
S01-name is not a row of inputs/answers.tsv" (15:01:37Z); S12 now gives its module row the source of the name, with its test (comment 6148707143 of #278). The
records commit is not the last of the branch (the lease commit follows it): the test finds it by its message. The gate then found two subtests of the demo of row 54 that read the
last row of a restart, which `verify` now is.

## Round 1: the fixes

Round 1 (comment 6149719082, three findings and a note): the guard of S15 was written red first (16:02:16Z: the rows `verify.sha256` and `rule-paths.sha256` missing, two hashes of another file
and a register with no path each passed); the case of a row that is `not-active` was green at once, as the code named any row that is not `pass` or `clear` and the test lacked, and the mutant
`row-only-fail` is its check; the recovery of the world test asserts that the second records commit changes `setup/record.tsv` alone; the note on the red-first claim of the S12 cases is corrected in its comment.

## Mutations

Each rule of the step, its tree, its files and its wiring, and the module row of S12, mutated alone on a copy of the tree and never on the worktree: a condition inverted or
dropped, a name, a path or a field changed, so each mutant builds. A mutant that does not build is counted as not run.

- **Run 1** (15:09:51Z to 15:13:42Z): 32 caught, 2 survived, 3 not run, of 37. Survived: `tree-step-S` (a step that does not begin with S passed the check of the writer) and `wire-message` (the test took the message of the commit from the constant that it tests). Three did not build (`go vet` refused a condition that is always true, a variable left unused), so none counted as caught: each is written again to build. The writer now takes the form of the block (`setupStepForm`), and the test names the commit by the text of its message.
- **Run 2** (15:15:35Z to 15:17:37Z): 11 caught, 0 survived, 0 not run, of 11. The two survivors and the three that did not build, and the mutants of the tree, which the refinement touched.
- **Run 3** (15:42:05Z to 15:46:26Z): 36 caught, 0 survived, 0 not run, of 36. All the 36 mutants on `7e6cfe2`, the head of round 1.
- **Run 4** (16:05:21Z to 16:10:13Z): 40 caught, 0 survived, 2 not run, of 42. The 42 mutants (the six of the guard of S15 and of the rows that are not pass or clear added after round 1) on the head that round 2 reads; two did not build (an import and a variable left unused).
- **Run 5** (16:21:39Z to 16:21:50Z): 2 caught, 0 survived, 0 not run, of 2. The two that did not build, written again to build, on the same head.

## Green

On `41867b0` (16:10:13Z to 16:21:34Z), each with exit 0, in this order (the time is when the command ended; `go test` with `-count=1`, then
`-tags=integration` and `-tags=e2e`, `TMPDIR` on the disk, `-timeout 30m`): go build ./... 16:10:15Z; go vet ./... 16:10:15Z; gofmt -l internal cmd (empty) 16:10:15Z; go test ./... 16:10:21Z; go test -tags=integration ./...
16:19:08Z; go test -tags=e2e ./... 16:20:39Z; adr-lint 16:20:41Z; prd-lint 16:20:41Z; link-lint 16:20:56Z; run-discipline-tests 16:21:04Z; nested-
checkout-check 16:21:04Z; setup-check 16:21:34Z; git diff --check origin/main HEAD 16:21:34Z.
The integration level of `internal/run` took 510 s of the 600 that CI's `-timeout 10m` gives a package (458 s at the end of T-eahr: each restart of the worlds of row 54 that passes `setup`
now runs S01 to S14 once more in the step `verify`, and the new world makes four restarts); CI's own job takes about 4 minutes. The first run of the gate on `8a8816b` (15:19:05Z to
15:28:00Z) exited 1 on two subtests of the demo of row 54 that read the last row of a restart (fixed in `7e6cfe2`, then green from 15:30:52Z to 15:41:57Z).
The second run of S01 to S14 that the step `verify` makes took 2.58 s, 2.67 s and 2.54 s in the three restarts of the demo world that run it (the loop of `setup` before it: 3.98 s with
the two prose sessions, then 2.54 s and 2.43 s; the loop that ends S15: 15 to 21 ms), measured on a scratch copy with the clock around each call of the loop; the estimate of about 15 seconds
in the plan was wrong. The commits after `0e1cd25` change only documents: this record, the task file and the two indexes.
