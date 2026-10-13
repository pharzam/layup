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

The unit level passed at 18:46:54Z. The world passed at its first run on the code (18:46:59Z to 18:49:02Z). The probe test failed once, on an assertion of the plan that was false on this host: the commit of an accepted probe is on the
remote-tracking ref `origin/layup-probe` of the clone, as the `git` of this host (2.56) makes it for a push to the URL of the clone's own `origin`. The test then read the refs after a refused push and an accepted one, and
wanted that ref. The fake forge gives the rules of a branch and the activity of a ref, and the bare repository of the world refuses the push to `layup-probe` by a hook (the `-c core.hooksPath` of
`internal/git` does not stop a hook of the receiving repository).

## The first push: CI

The `tests` job of CI failed on `e7e818c`, in `TestTheProbeIsAPassOnlyForTheRefusalOfTheRulesetOfTheDefaultBranch` alone: the `git` of CI makes no remote-tracking ref for that push, so the commit was "on \"\" of the clone" where this host's
`git` has it on `origin/layup-probe`. The test accepts both forms now, the specification says that the ref may be left, and a pitfall of `guardrails.md` §2 says it. The integration level of `internal/run` took 282 s in CI (the
world makes about twelve restarts), against its limit of 600 s for a package.

## Written after the code

The exact detail of each failing row (a missing parameter and a different one, a record that is not in the records and one that its reader refuses, a ruleset file that is missing and one with no rule), the cases of the framing
that no marker or separator breaks, the probe made on the head that the forge gives, the work area made sure of for a last commit only, the first line and the valid UTF-8 of a detail, a table that its reader refuses, and the
scenarios of the world (the default branch one commit ahead, a ruleset file gone from the work area, a work area that cannot be built, one built again whose head is not the last commit, a push of the table that is refused once)
were written after the code and before the final run of the mutants, each for a condition of the plan review or for a mutant that no test could see. The scenarios that conditions 2 and 3 of the plan review named, a failed rebuild and a
rebuilt head that differs, should have been red first: round 1 took it as a failed criterion (finding 2). The red of each is against the mutant of its rule, which shows the right reason and not the order: the rebuilt head that
differs is red against `w-area-head` (run 8, at its scenario), the loop that fails against `w-area-loop`, the work area that is a link against `w-area-build`, the details against `area-prefix`, `tree-nil`, `tree-reader-prefix`,
`rules-file-prefix` and `outputs-nil`, the read-back against `rb-first`, `pd-missing`, `pd-differs`, `pd-sorted` and `sp-key-context`, the framing against `ra-prefix`, `ra-index` and `ra-end`, the verdict against `pv-ruleset-name`, the
details of a row against `row-clean`, `row-first-line` and `row-utf8`, the table's reader against `do-table-reader`, and the push that is refused once against `w-fenced`.

## Mutations

Each rule of the step mutated alone on a copy of the tree and never on the worktree: a condition inverted or dropped, a name, a path or a field changed, so each mutant builds. A mutant that does not build is counted as not run.
The 110 cover the read-back (the type match, the first candidate, the extra types and their order, the parameters in the order of their names, the pairs of the required checks), the reader of the ruleset file, the framing of the applies, the two
bypass lists, the verdict of the probe, the probe's commit, the eight rows (each of the checks, its details, its clock), the step and its wiring, and the work area of `scaffold-push` that the step shares.

- **Run 1** (18:55:29Z to 18:59:58Z): 77 caught, 1 survived, 3 not run, of 81. Unit-level mutants only, on a tree of the first form of the tests. Survived: `rr-json` (a file that is not JSON was refused with another reason, and the test did not read the reason); not run: `bypass-outs-err`, `rb-first` and `sp-fallback` (they left a variable unused, so they were written again). After this run the reasons of the reader of the ruleset file were asserted, the world gained the scenario of a work area built again whose head is not the last commit, and the code and the tests were made shorter.
- **Run 2** (19:02:42Z to 19:22:09Z): 105 caught, 1 survived, 0 not run, of 106. All the mutants, the unit level and the world, on the code of `7c03bc2`. Survived: `w-fenced` (the commit without fencing).
- **Run 3** (19:23:49Z to 19:24:54Z): 0 caught, 1 survived, 0 not run, of 1. `w-fenced` again, with a scenario in which another run takes the lease as the step starts: it survived, as the beat of the lease can see the loss before the push is refused, and the two give the same detail. The scenario was replaced.
- **Run 4** (19:29:07Z to 19:30:10Z): 1 caught, 0 survived, 0 not run, of 1. `w-fenced` with a scenario in which the forge refuses the push of the table once (a hook of the bare repository): fencing reads the lease, still the run's, and tries the push once more, so the table is committed; without it the step fails with "the records push was refused". Caught, every time.
- **Run 5** (19:59:01Z to 19:59:16Z): 5 caught, 0 survived, 0 not run, of 5. The pairs of the required checks after the fix of round 1 (finding 1): `sp-key-id`, `sp-key-context`, `sp-key-zero` (which brings the bug of the finding back), `sp-set` and `sp-fallback`.
- **Run 6** (20:00:59Z to 20:01:25Z): 5 caught, 0 survived, 0 not run, of 5. The five mutants of the probe on the probe test after the CI of the first push (the remote-tracking ref of an accepted push is accepted, and not wanted).
- **Run 7** (20:01:58Z to 20:04:06Z): 1 caught, 1 survived, 0 not run, of 2. Two mutants of the shared `workAreaAgain` on the world, to see the two scenarios of condition 3 that were written after the code red: `w-area-head` (the head of a work area built again is not compared) is caught at the scenario of the head that differs; `w-area-build` (the error of the build of the area is dropped) survived, as the scenario that was meant for it, a file in the place of the area's directory, ends in the `os.Stat` of the area, before the build.
- **Run 8** (20:08:05Z to 20:11:22Z): 2 caught, 1 survived, 0 not run, of 3. With the scenario of a link where the directory of the work area belongs, which does reach the build (the work area refuses a link), and the scenario of a loop that fails (a glossary whose citation resolves to no fact): `w-area-head` and `w-area-loop` are caught at their scenarios; `w-area-build` survived again, the run before the world had the link.
- **Run 9** (20:12:58Z to 20:31:30Z): 26 caught, 0 survived, 0 not run, of 26. The 26 mutants of the world and of the probe (the slow ones) on the final world, with the scenarios of the failed loop, the link and the head that differs: `w-area-head`, `w-area-loop` and `w-area-build` are caught at their scenarios, and `w-fenced` at the push that is refused once.
- **Run 10** (20:31:40Z to 20:35:43Z): 84 caught, 0 survived, 0 not run, of 84. The 84 unit-level mutants on the final tests, after the fix of round 1. With run 9, all 110 are caught on the code of the fix.

## Round 1 and its fix

Round 1 (comment 6154043779, `e7e818c`, GPT-6.1 Sol on Codex CLI `xhigh`, 9 min 37 s from 19:46:55Z): `material`, three findings in the change and one note.

- Finding 1, a required check whose `integration_id` the call lacks compared equal to one whose file `integration_id` is 0: the two cases were written first and were red at 19:57:54Z (`true` where a difference is wanted, in both directions); the key of a pair is now the context and the
  `integration_id` as they are, absent or not (green at 19:58:07Z; run 5 of the mutants).
- Finding 2, the cases written after the code: the paragraph "Written after the code" above, with the red of each against its mutant; a pitfall of `guardrails.md` §2.
- Finding 3, stale names of tests in the traceability and the PRD cells (three tests that were removed or merged) and a comment of `pushProbe`: fixed; the names are checked against the code.
- Finding 4, a note: the reviewer's sandbox cannot open a loopback socket, so it could not run the integration level; nothing to do.

## Green

On `7c03bc2` (19:31:32Z to 19:44:42Z), the code that round 1 read less its comment, each with exit 0, in this order (the time is when the command ended; `go test` with `-count=1`, then
`-tags=integration` and `-tags=e2e`, `TMPDIR` on the disk, `-timeout 30m`): go build ./... 19:31:33Z; go vet ./... 19:31:34Z; gofmt -l internal cmd (empty) 19:31:34Z; go test ./... 19:31:40Z; go test -tags=integration ./...
19:42:21Z; go test -tags=e2e ./... 19:43:50Z; adr-lint 19:43:52Z; prd-lint 19:43:52Z; link-lint 19:44:06Z; run-discipline-tests 19:44:13Z; nested-
checkout-check 19:44:13Z; setup-check 19:44:42Z; git diff --check origin/main HEAD 19:44:42Z.

On `3adb2b8` (20:35:49Z to 20:49:23Z), the head that round 2 reads, each with exit 0: go build ./... 20:35:51Z; go vet ./... 20:35:51Z; gofmt -l internal cmd (empty) 20:35:52Z; go test ./... 20:35:58Z; go test -tags=integration ./...
20:47:02Z; go test -tags=e2e ./... 20:48:30Z; adr-lint 20:48:32Z; prd-lint 20:48:32Z; link-lint 20:48:46Z; run-discipline-tests 20:48:54Z; nested-
checkout-check 20:48:55Z; setup-check 20:49:23Z; git diff --check origin/main HEAD 20:49:23Z.
The integration level of `internal/run` took 622 s on this host on the first and 645 s on the second (573 s at the end of T-hfn1: the new world makes about twelve restarts, several of which rebuild the work area); CI's own run took 282 s, against the 600 s of its
`-timeout 10m`. The run here used `-timeout 30m`.
