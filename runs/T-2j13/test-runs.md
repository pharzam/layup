# T-2j13 — test runs

The forge check at each later start (#281, row 57b). All times are UTC, 2026-10-14, read from the clock after each step.

## The order of the gate

The claim was pushed at 06:09:07Z, after its secret scan, and the plan was posted at 06:11:23Z; its review (GPT-6.1 Sol on Codex CLI, `high`) ran from 06:11:23Z to 06:16:50Z and was posted at 06:17:21Z, with the author's
answer at 06:17:41Z. The base was run at 06:11:31Z. The specification was committed first (06:18:43Z by the clock of the commit: the readings of the plan, with the R10 correction that
condition 3 of the plan review asked for) and no test or code came before it.

## The base (06:11:31Z to 06:11:53Z)

At `18565ce`: build, vet and `go test ./...` exit 0, and the CI of `main` had passed.

## Red

All the tests were written and run red before any line of the behavior (condition 1 of the plan review). With a skeleton that has the new parameter of `restartNames` and ignores it, the unit test was red at 06:19:40Z in its
one case that the skeleton cannot satisfy, the Scaffold done; its three other chains pass at the base, as they must (guards). The world was red from 06:20:30Z to 06:21:39Z, at the two restarts that are later starts (the rows and the
notes of progress: `forge-check` last, `21/21`, where the fifth row and `5/21` are wanted); the details of the steps that are done before at a later start, the first-start chains of its later scenarios and the table that its reader
refuses pass at the base: they are guards of behavior that the base has, and are recorded as such, not as red.

## Turning green

The unit level passed at 06:22:11Z and the world at its first run on the code (06:22:23Z to 06:23:36Z), with `TestTheStepRuleGaps`, whose notes of a restart with the Scaffold not done are unchanged.

## Mutations

Each rule of the change mutated alone on a copy of the tree and never on the worktree: the reading of the table (a table that its reader refuses, a row that is not a pass, a row that is a pass), where `clone` reads it and its path, the names
that `Restart` gives, the place of the step in each chain (before `lease`, never early, a step of Intake dropped), the names of the chains and the guard of Start not done. A mutant that does not build is counted as not run.
Run 1 (06:24:10Z to 06:35:08Z): 13 caught, 0 survived, 0 not run, of 13. The mutants show the sensitivity of the tests, and say nothing of their order.

## Green

On `9a42eb0` (06:36:31Z to 06:50:17Z), each with exit 0, in this order (the time is when the command ended; `go test` with `-count=1`, then `-tags=integration` and `-tags=e2e`, `TMPDIR` on the disk,
`-timeout 30m`): go build ./... 06:36:32Z; go vet ./... 06:36:33Z; gofmt -l internal cmd (empty) 06:36:33Z; go test ./... 06:36:39Z; go test -tags=integration ./...
06:47:56Z; go test -tags=e2e ./... 06:49:25Z; adr-lint 06:49:27Z; prd-lint 06:49:27Z; link-lint 06:49:41Z; run-discipline-tests 06:49:49Z; nested-
checkout-check 06:49:49Z; setup-check 06:50:17Z; git diff --check origin/main HEAD 06:50:17Z.
The integration level of `internal/run` took 655 s on this host (645 s at the last run of T-wphp); the change adds no world.
