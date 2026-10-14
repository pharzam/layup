# T-eevt — test runs

The e2e of the Scaffold (#282, row 58). All times are UTC, 2026-10-14, read from the clock after each step.

## The order of the gate

The claim was pushed at 07:12:05Z, after its secret scan, and the plan was posted at 07:14:57Z; its review (GPT-6.1 Sol on Codex CLI, `high`) ran from 07:14:57Z to 07:20:49Z and was posted at 07:21:23Z, with the author's
answer at 07:21:25Z. The base was run at 07:15:07Z, and the specification was committed first (07:21:42Z by the clock of the commit: the five stages and the stand-ins, "decided here") before the first test.

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

- **Run 1** (07:35:19Z to 07:59:15Z): 8 caught by an assertion, 0 survived, 3 not counted, of 11. The eleven mutants of the plan: the later chain off, the check before `lease`, the lease kept, the table not written, no rule answered, the probe accepted, the wrong last commit of the tree, the applies not recorded, the commands not printed, the wait ignored and the record of the verification misplaced. Eight were caught by an assertion; `lease-kept` hung (the next restart waits for the lease that the mutant never releases) and was stopped by hand, so it does not count here; `commands-not-printed` and `wait-ignored` did not build (an unused variable), so they were written again.
- **Run 2** (07:59:31Z to 08:05:43Z): 6 caught, 0 survived, 0 not run, of 6. The six mutants that run 1 did not settle: the three above, with the limit of three minutes for each run of the binary that turns the hang of `lease-kept` into a failure that says so, and the three of the witnesses that condition 3 of the plan review asked for (the order of the checks of the verification, the rule-paths register, the README of the records). With run 1, all 14 are caught.
- **Run 3** (08:39:30Z to 08:48:28Z), after round 1: the eight mutants of its findings against the strengthened scenario (new: 8 caught, 0 survived, 0 not run, of 8) and against the scenario that round 1 read (old: 0 caught, 8 survived, 0 not run, of 8). The strengthened assertions are red against each of the eight defects, and the old ones let each pass: `order-swap` (the steps `setup` and `verify` swapped), `later-drops-setup` (the later chain without `setup`), `leak-path` (the clean of the details off and a host path in the detail of `markers`), `leak-token` (a token in a note of the wait), `gates-dropped` (a gate row of one kind only), `record-bytes` (a byte added to `setup/record.tsv` at its publication), `readme-text` (a word of the README of the records changed) and `apply-endpoint` (the apply of a ruleset to another endpoint).

## Round 1 and its fix

Round 1 (comment 6164621188, `a631df1`, GPT-6.1 Sol `xhigh`) gave `material`: five findings, each of an assertion that the scenario lacked or made too weak, three notes. Each fix is the assertion first: it was written, run on the code as it is (green, 08:38:14Z
to 08:38:47Z, 30 s), and shown red against the mutant that the old assertion let pass (run 3, above), the new scenario against the old one. Finding 1: the whole sequence of rows of each restart (`order-swap`, `later-drops-setup`). Finding 2: no host path in a
detail and no token in an output, for Start and each restart (`leak-path`, `leak-token`). Finding 3: the whole ordered list of the rows of the verification, with the result of each. Finding 4: the three published files equal their files of `out/`, and
the README is the block of the specification. Finding 5: the stand-in `gh` takes only the one call of `commands.sh` (`apply-endpoint`). Notes: the two times and the tally of run 1 of this record are corrected (the comments of the plan review at 07:21:23Z and 07:21:25Z, the specification
commit at 07:21:42Z), the specification and the scenario say what the stages are, and the Go environment of the real verification is the cache, `GOTOOLCHAIN=local`, `GOPROXY=off` and `GOFLAGS=` that `goEnv` gives, with a new `HOME` and `TMPDIR` (the author's answer to the plan review said `GOENV=off` too, which was wrong).

## Green

On `47645e8` (08:53:59Z to 09:08:08Z), each with exit 0, in this order (the time is when the command ended; `go test` with `-count=1`, then `-tags=integration` and `-tags=e2e`, `TMPDIR` on the disk,
`-timeout 30m`): go build ./... 08:54:01Z; go vet ./... 08:54:01Z; gofmt -l internal cmd (empty) 08:54:02Z; go test ./... 08:54:08Z; go test -tags=integration ./...
09:05:14Z; go test -tags=e2e ./... 09:07:15Z; adr-lint 09:07:16Z; prd-lint 09:07:16Z; link-lint 09:07:31Z; run-discipline-tests 09:07:38Z; nested-
checkout-check 09:07:39Z; setup-check 09:08:08Z; git diff --check origin/main HEAD 09:08:08Z.
An earlier run of the whole gate, on `80c2a51` (08:06:27Z to 08:20:43Z), the head before round 1's fix, was also green in each step.
The e2e level took 121 s on this host for `cmd/layup` (84.6 s at the base), against CI's `-timeout 10m`; the scenario alone takes about 31 s.
