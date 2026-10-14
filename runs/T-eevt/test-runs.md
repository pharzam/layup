# T-eevt — test runs

The e2e of the Scaffold (#282, row 58). All times are UTC, 2026-10-14, read from the clock after each step.

## The order of the gate

The claim was pushed at 07:12:05Z, after its secret scan, and the plan was posted at 07:14:57Z; its review (GPT-6.1 Sol on Codex CLI, `high`) ran from 07:14:57Z to 07:20:49Z and was posted at 07:21:25Z, with the author's
answer. The base was run at 07:15:07Z, and the specification was committed first (07:22:15Z: the five stages and the stand-ins, "decided here") before the first test.

## The base (07:15:07Z to 07:15:29Z)

At `a880936`: build, vet and `go test ./...` exit 0, and the CI of `main` had passed; the e2e level of `cmd/layup` (07:15:29Z to 07:16:56Z) exit 0 in 84.6 s.

## Red, stage by stage

The scenario was written with the assertions of its five stages and a skeleton of the helpers that compiled (the binary and the world of the other scenarios of `layup run`), and run at each step, which is the e2e's red: a stage is red for what its
world lacks, and the cause of each run is the one read from its output. Stage 5 and the record witnesses that already held were written with their stages and passed on arrival: they are guards of behavior that the earlier rows built, recorded
as such, and not as red.

| Run | Time | The first failure | The support that followed |
| --- | ---- | ----------------- | ------------------------- |
| 1 | 07:25:13Z to 07:25:22Z | stage 1: the markers row says "no marker in the pinned tree" and the batch holds 5 questions where 7 are wanted: the skeleton's baseline is a README | the scaffold baseline, built into the pin of a second binary |
| 2 | 07:25:53Z to 07:26:12Z | stage 1: `markers` fails, "the session: refused file: line 1, column id: the header row has start in its place": the harness has no case for the session of the step | the case of `markers.tsv` in the world's harness |
| 3 | 07:27:01Z to 07:27:23Z | stage 1 passes; stage 2: `setup` fails, "the session: crash": the harness has no prose session | the prose session of the harness and the prose files |
| 4 | 07:27:44Z to 07:28:03Z | stage 2 passes (the setup, the real verification and the wait); stage 3: `commands.sh` exits 4, "gh auth login": the `gh` of the host answers | the stand-in `gh`, first on the `PATH` of the file |
| 5 | 07:28:14Z to 07:28:44Z | stage 3 passes; stage 4: `scaffold-push` waits again, as the forge gives no head of a branch | the head of each branch (`commit.sha`) in the fake forge |
| 6 | 07:28:55Z to 07:29:25Z | the applies are recorded; `forge-check` fails, "rules:main: EffectiveRules: status 404" | the rules of a branch and the activity in the fake forge |
| 7 | 07:29:36Z to 07:30:06Z | `forge-check` fails, "probe: probe accepted": nothing refuses the push | the `pre-receive` hook of the bare repository |
| 8 | 07:30:16Z to 07:30:47Z | stage 4 passes but for an assertion of the test: the commit of `setup/tree.tsv` is after a tab and not a line feed | the assertion, a defect of the test |
| 9 | 07:30:54Z to 07:31:27Z | none: the whole scenario passes in 30.3 s, stage 5 with it | |

## After the first green

Four things were written after the scenario's first green: the assertions of the records (the thirteen checks of the verification in order, the rule-paths register and the README of the records) that condition 3 of the plan review named,
each of which passed on arrival, as a guard of the earlier rows, and each shown sensitive by a mutant of its own (run 2); the second of waiting before the later start, because two runs of the scenario in succession differed in one assertion (the table of the
later start was strictly later than the first by the second of its clock, which a restart can fall inside: a pitfall of `guardrails.md` §2), after which the scenario passed twice in succession (30.8 s and 31.2 s); and the limit of three minutes for each run of
the binary, which the mutant `lease-kept` showed to be missing. The three are test support and guards, not tests of a behavior that was missing.

## Mutations

Each rule of the plan on a copy of the tree and never on the worktree: a condition inverted, a path or a name changed, a part dropped, so each mutant builds. A mutant that does not build is counted as not run. The mutants show the sensitivity of
the scenario, and say nothing of its order.

- **Run 1** (07:35:19Z to 07:59:15Z): 9 caught, 0 survived, 2 not run, of 11. The eleven mutants of the plan: the later chain off, the check before `lease`, the lease kept, the table not written, no rule answered, the probe accepted, the wrong last commit of the tree, the applies not recorded, the commands not printed, the wait ignored and the record of the verification misplaced. Eight were caught by an assertion; `lease-kept` hung (the next restart waits for the lease that the mutant never releases) and was stopped by hand, so it does not count here; `commands-not-printed` and `wait-ignored` did not build (an unused variable), so they were written again.
- **Run 2** (07:59:31Z to 08:05:43Z): 6 caught, 0 survived, 0 not run, of 6. The six mutants that run 1 did not settle: the three above, with the limit of three minutes for each run of the binary that turns the hang of `lease-kept` into a failure that says so, and the three of the witnesses that condition 3 of the plan review asked for (the order of the checks of the verification, the rule-paths register, the README of the records). With run 1, all 14 are caught.

## Green

On `80c2a51` (08:06:27Z to 08:20:43Z), each with exit 0, in this order (the time is when the command ended; `go test` with `-count=1`, then `-tags=integration` and `-tags=e2e`, `TMPDIR` on the disk,
`-timeout 30m`): go build ./... 08:06:28Z; go vet ./... 08:06:29Z; gofmt -l internal cmd (empty) 08:06:29Z; go test ./... 08:06:35Z; go test -tags=integration ./...
08:17:49Z; go test -tags=e2e ./... 08:19:50Z; adr-lint 08:19:52Z; prd-lint 08:19:52Z; link-lint 08:20:06Z; run-discipline-tests 08:20:13Z; nested-
checkout-check 08:20:14Z; setup-check 08:20:43Z; git diff --check origin/main HEAD 08:20:43Z.
The e2e level took 121 s on this host for `cmd/layup` (84.6 s at the base), against CI's `-timeout 10m`; the scenario alone takes about 31 s.
