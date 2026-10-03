# T-dep6: the test runs

The runs of `T-dep6` (#93). Host: macOS, `go1.27.1`, `git` 2.54.0,
2026-10-03, UTC. "…" marks a cut part of an output.

## Test first, for tests of code that exists

This row adds e2e tests to the code of rows 9 to 15, which exists, so a test
cannot fail first on missing code. Step 1 of the plan reads "test first" as a
mutation run: each scenario runs against a change of the code that it guards,
and must fail on its own assertion (note 4 of the plan review). A script
changes one place, runs the five scenarios of `cmd/layup/whole_e2e_test.go`
through the binary that `TestMain` builds, and writes the file back from its
text before the change; `git status` after each run shows only the new test
file, so no mutation was committed.

| Mutation (the file and the line at `9a16ff8`) | The scenario that failed, and its first assertion |
| --------------------------------------------- | -------------------------------------------------- |
| M1 S10 names the last line of a marker (`internal/setup/scaffold.go:372`, `markerQuestions`) | `TestAWholeSetupOnAStandInBaseline`: "the stop S10: exit 3 … want 3 and" the table with `docs/ops.md:3` |
| M2 the runner forgets that S02 is done, so a run resolves the pin again (`internal/setup/setup.go:283`) | `TestAWholeSetupOnAStandInBaseline`: "the stop S05 on a second run: exit 2" (the run after the baseline moved) |
| M3 S15 starts `layup-records` at the setup head, not as an orphan (`internal/setup/final.go:509`, `commitRecords`) | `TestTheRecordsAreInTheTargetsGit`: "git merge-base main layup-records: \"f8c2…\", <nil>; want no common commit"; `TestThePushesOfCommandsSh`: "… want the records commit …, with no parent" |
| M4 `commands.sh` applies the ruleset before the push of `layup-setup` (`final.go:289`, order 4 to 1) | `TestTheRecordsAreInTheTargetsGit`: "the commands of commands.sh: … want, in this order (S03, S13, S15, S13)" |
| M5 `commands.sh` pushes `layup-setup` as a branch of its own, not onto `main` (`final.go:288`) | `TestTheRecordsAreInTheTargetsGit` (the commands); `TestThePushesOfCommandsSh`: "main of the bare repository 5fc8…; want the setup head 506f…" |
| M6 the ruleset requires a `layup/` check (`final.go:213`, `rulesetBody`) | `TestTheTargetPassesItsGateWithLAYUPAbsent`: "the required checks [… \"layup/gates\" …]; want the five kinds twice, and no layup/ check" |
| M7 the gate job of `static` runs `layup` (`internal/catalog/go/files/.github/workflows/gates.yml.tmpl:25`) | `TestTheTargetPassesItsGateWithLAYUPAbsent`: ".github/workflows/gates.yml:25: \"      - run: layup gate . --kind static\" names layup or setup-check" |
| M8 check `sources` takes an `answer` ref that `answers.tsv` does not have (`internal/verify/sources.go:36`) | `TestLayupSetupVerifyOnAWholeSetup`: "a value row whose answer is not a row of answers.tsv: <nil>, exit 0" |
| M9 check `jobs` gives no finding (`internal/verify/jobs.go:37`) | `TestLayupSetupVerifyOnAWholeSetup`: "a removed gate job: <nil>, exit 0" |
| M10 S15 takes a row `not-active` (`final.go:419`, `runS15`) | `TestLayupSetupVerifyOnAWholeSetup`: "S15 with a not-active row: exit 0" |
| M11 a pending kind is a pass (`internal/verify/gates.go:57`) | `TestLayupSetupVerifyOnAWholeSetup`: "layup setup verify at the stop O-verify: exit 0 … want 0 and" the table with the three pending rows `clear` |

M2 first failed through the shared run itself: the maker stopped with an
error when `layup-records` was not there, so each scenario failed on that
error and not on its own assertion. The maker now keeps an empty value there,
and the run above is the second one.

## What the first runs found

- **A bare repository with no default branch.** The first run of the five
  scenarios on the code of `main` passed four; the fifth failed at
  `git show HEAD:docs/gates.tsv: exit status 128`, as the bare repository that
  stands in for GitHub was made with no `-b`, so its `HEAD` named `master`, and
  the clone had no checkout. The test now makes it with `-b main`, the default
  branch of the target; the lesson is in `docs/guardrails.md` §2.
- **A name that the package had.** `commitAll` was the name of a helper of
  `gate_e2e_test.go`; the new helper is `commitIn`.

## The green runs

The five scenarios on the code of `main` (`go test -count=1 -tags=e2e
-timeout 10m -run '…' ./cmd/layup/`): the shared whole-setup run takes 6 s,
and the five scenarios pass in 19 s (`TestAWholeSetupOnAStandInBaseline` 5.9 s,
`TestLayupSetupVerifyOnAWholeSetup` 8.1 s, `TestTheRecordsAreInTheTargetsGit`
0.5 s, `TestThePushesOfCommandsSh` 1.1 s,
`TestTheTargetPassesItsGateWithLAYUPAbsent` 2.6 s).

On the freeze head `8647291`, 02:08:05Z to 02:10:26Z, the ladder: the local
checks of `AGENTS.md`, `go build`, `go vet` with each tag, `gofmt -l .`, the
three test levels (the e2e package of `cmd/layup` in 28 s), the race tests and
`run.sh` (44 passed): each of the 18 steps exit 0. CI of PR #116 on
`8647291` (Linux, `git` 2.55.0) passed its job `tests` before review round 1,
which ran the e2e level twice on it with the same result.

The demo, in the shared run: on a stand-in baseline by its `file://` URL, with
a problem statement with one gap, `layup setup` stops at S01 (the four
questions of S01 and `Q-001`), at S05 (the link that the deletion breaks), at
the prose step (six input files, one of them a file that check `adapted`
flags), at S10 (three markers: one twice, one that does not close on its line,
and one kept as a gap) and at S15 (`O-verify`), each stop the same bytes on a
second run; a commit of the baseline after S02 changes no `pin.commit`; with
the table of `layup setup verify` (exit 0, each row `pass` or `clear`), the
last run does S15 and exits 0, and a rerun gives the same table and no second
records commit.
