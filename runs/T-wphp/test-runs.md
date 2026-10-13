# T-wphp — test runs

The step forge-check (#280, row 57a). All times are UTC, 2026-10-13, read from the clock after each step.

## The order of the gate

The claim was pushed after its secret scan, which ended at 18:22:57Z, and the plan was posted at 18:27:41Z; its review (GPT-6.1 Sol on Codex CLI, `high`) ran from 18:27:44Z to
18:32:19Z and was posted at 18:32:52Z, with the author's answer at 18:32:54Z. The base was run at 18:33:47Z, and the specification was committed first (18:35:05Z by the clock of the
commit: the readings of conditions 1 and 2 of the plan review, an R10 correction) before the first test, which did not build until the code existed (18:36:49Z: `undefined: forgeCheckJob`).

## The base (18:33:47Z to 18:34:09Z)

At `be28498`: build, vet and `go test ./...` exit 0, and the CI of `main` had passed.

## Red

With a skeleton that compiles and does nothing (18:41:11Z), each test fails on its assertion: the unit level at 18:44:38Z (the test of the names of the steps, which wants `forge-check` after `scaffold-push`, and ten
tests of the pure parts and of the step with its seams, each with an empty step), and the integration level from 18:44:38Z to 18:45:28Z (the probe, whose skeleton makes no commit; the world, which has no row
`forge-check`; and the progress, which counts 20 steps where 21 are wanted).

## Turning green

The unit level passed at 18:46:54Z. The world passed at its first run on the code (18:46:59Z to 18:49:02Z). The probe test failed once, on an assertion of the plan that was false: the commit of an accepted probe is on the
remote-tracking ref `origin/layup-probe` of the clone, as `git push` to the URL of the clone's own `origin` makes it. The specification says "no branch" now, the test reads the refs after a refused push and an accepted one,
and a pitfall is written. The fake forge gives the rules of a branch and the activity of a ref, and the bare repository of the world refuses the push to `layup-probe` by a hook (the `-c core.hooksPath` of
`internal/git` does not stop a hook of the receiving repository).

## Written after the code

The exact detail of each failing row (a missing parameter and a different one, a record that is not in the records and one that its reader refuses, a ruleset file that is missing and one with no rule), the cases of the framing
that no marker or separator breaks, the probe made on the head that the forge gives, the work area made sure of for a last commit only, the first line and the valid UTF-8 of a detail, a table that its reader refuses, and the
scenarios of the world (the default branch one commit ahead, a ruleset file gone from the work area, a work area that cannot be built, one built again whose head is not the last commit, a push of the table that is refused once)
were written after the code and before the final run of the mutants, each for a condition of the plan review or for a mutant that no test could see.

## Mutations

Each rule of the step mutated alone on a copy of the tree and never on the worktree: a condition inverted or dropped, a name, a path or a field changed, so each mutant builds. A mutant that does not build is counted as not run.
The 106 cover the read-back (the type match, the first candidate, the extra types and their order, the parameters in the order of their names, the pairs of the required checks), the reader of the ruleset file, the framing of the applies, the two
bypass lists, the verdict of the probe, the probe's commit, the eight rows (each of the checks, its details, its clock), the step and its wiring.

- **Run 1** (18:55:29Z to 18:59:58Z): 77 caught, 1 survived, 3 not run, of 81. Unit-level mutants only, on a tree of the first form of the tests. Survived: `rr-json` (a file that is not JSON was refused with another reason, and the test did not read the reason); not run: `bypass-outs-err`, `rb-first` and `sp-fallback` (they left a variable unused, so they were written again). After this run the reasons of the reader of the ruleset file were asserted, the world gained the scenario of a work area built again whose head is not the last commit, and the code and the tests were made shorter.
- **Run 2** (19:02:42Z to 19:22:09Z): 105 caught, 1 survived, 0 not run, of 106. All the mutants, the unit level and the world, on the final code. Survived: `w-fenced` (the commit without fencing).
- **Run 3** (19:23:49Z to 19:24:54Z): 0 caught, 1 survived, 0 not run, of 1. `w-fenced` again, with a scenario in which another run takes the lease as the step starts: it survived, as the beat of the lease can see the loss before the push is refused, and the two give the same detail. The scenario was replaced.
- **Run 4** (19:29:07Z to 19:30:10Z): 1 caught, 0 survived, 0 not run, of 1. `w-fenced` with a scenario in which the forge refuses the push of the table once (a hook of the bare repository): fencing reads the lease, still the run's, and tries the push once more, so the table is committed; without it the step fails with "the records push was refused". Caught, every time. With it, all 106 are caught on the code of the head.

## Green

On `7c03bc2` (19:31:32Z to 19:44:42Z), each with exit 0, in this order (the time is when the command ended; `go test` with `-count=1`, then
`-tags=integration` and `-tags=e2e`, `TMPDIR` on the disk, `-timeout 30m`): go build ./... 19:31:33Z; go vet ./... 19:31:34Z; gofmt -l internal cmd (empty) 19:31:34Z; go test ./... 19:31:40Z; go test -tags=integration ./...
19:42:21Z; go test -tags=e2e ./... 19:43:50Z; adr-lint 19:43:52Z; prd-lint 19:43:52Z; link-lint 19:44:06Z; run-discipline-tests 19:44:13Z; nested-
checkout-check 19:44:13Z; setup-check 19:44:42Z; git diff --check origin/main HEAD 19:44:42Z.
The integration level of `internal/run` took 622 s on this host (573 s at the end of T-hfn1: the new world makes about ten restarts, several of which rebuild the work area); CI's own run of the
integration level of that package took 125 s on `main` at `720a149`, against the 600 s of its `-timeout 10m`. The run here used `-timeout 30m`.
