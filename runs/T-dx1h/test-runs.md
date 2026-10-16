# T-dx1h — the runs of the tests

Row 66b of the plan of `M2e`, [#358](https://github.com/pharzam/layup/issues/358). Times are those of the host's clock,
UTC, 2026-10-16, read with `date -u` after each step. Each exit status is that of `go test` itself, never of a pipe.
The plan review is comment 6214506834 and the author's answer 6214542025 (11:31:02Z); every red run below is later.

## Red first

| Run | Command (focused) | Time | Exit | What failed |
| --- | ----------------- | ---- | ---- | ----------- |
| 1 | `go test -count=1 -run TestCheckTableRefuses… ./internal/handoff/` | 11:32:20–11:32:22 | 1 | by assertion: the three refusal cases (a rendered row that lists `schema`, `schema` alone, an extra pair of code's own handoff); the five controls passed |
| 2 | the 13 new unit tests of `internal/run`, against a stub that returned `errNotBuilt` | 11:35:20–11:35:26 | 1 | **not a genuine red**: most stopped at the setup with the stub's own error; the stub was replaced |
| 3 | the same, against a stub that returns empty values | 11:35:47–11:35:53 | 1 | assertion, but most stopped at the precondition "want the six rows"; the declarative table was added to the stub and the fixtures were made from the table, not from `DefaultParameters` |
| 4 | the same | 11:37:07–11:37:11 | 1 | 13 FAIL, each on its own behaviour, no panic |
| 5 | the 13 integration tests (`-tags=integration`) | 11:39:29–11:40:21 | 1 | each FAIL; the boundary tests stopped at the fixture (no first writer yet), so `loopWorld` seeds a table that the test makes |
| 6 | the 8 boundary integration tests | 11:40:59–11:41:29 | 1 | all 8 FAIL by their own assertions (`applied [], refused []`, no copy of an edit, no `LostError`, …) |

## Green

`go test -count=1 ./internal/run/ ./internal/handoff/` 11:44:25–11:44:31, exit 0 (one case of my own test was wrong and was
corrected before: `intake.cap 50.0, 8.0` is a form error, a value holds no space). The 13 focused integration tests 11:44:41–11:45:33,
exit 0. `go test -tags=integration ./internal/handoff/ ./internal/tsv/` and `-run TestPackageRules ./cmd/layup/` (after `packages.md`) exit 0.
The spec-table test failed once on a sentence that wraps in `loop.md` (my assertion did not normalise white space) and was corrected.

## The mutants

The helper (`/tmp/dx1h-mut/mutate.py`, reviewed before the loop and self-tested with a no-op, a mutant that does not build, one that
is not found and one that is caught) takes a fresh backup of the file from `git show HEAD:<file>` for each mutant, refuses to run if the
working file is not HEAD's, applies the change exactly once, runs the focused tests, reads the exit status, restores the file from HEAD,
verifies its SHA-256 against HEAD's, and after the list runs a green control on the restored tree. A mutant that does not build is
**not credited**, and nor is one that panics or is not applied.

**Sweep 1**, 11:54:25–11:59:06, 39 mutants: **30 killed by assertion, 5 survived, 4 did not build** (not credited). Controls after: exit 0.

| Mutant | Result |
| ------ | ------ |
| H1 rule on `to` instead of `from`; H2 rule off; H3 rule for the default pair only; H4 rule for every row | killed (`handoff_test.go`) |
| P1 bound 1 or more dropped; P2 `>=` to `>`; P3 `>=` to `<`; P4 superseded never; P6 role of the row unchecked | killed |
| P9 reason `—` accepted; P10 second line accepted; P11 first word `settle`; P12 bound words unchecked; P13 Start evidence unchecked; P14 set with no reason; P15 value with no comment | killed |
| P16 comment not recorded; P17 caller's table edited; P18 table check off; P19 table of the change dropped; P20 Start evidence dropped; P21 existing table not kept; P23 type unchecked; P24 bound unchecked | killed |
| I2 Start note dropped; I4 table not written with the row; I5 copy of another issue looked up; I6 first copy not required; I9, I10 no progress line | killed |
| P7 an edit is a decision; P8 an App's comment is a decision; I3 write without the fence; I11 an App's comment copied | **survived**: each showed a missing case, and the tests were added (see sweep 2) |
| P22 `len(f) < 4` for `< 5` in the form | **survived, equivalent**: a body of four fields has no reason, and `reason == ""` refuses the same body |
| P8b, I1 (`return false` leaves `id` unused), I7, I8 (`more`/`parameters` left unused) | **build failed**, not credited; redone in sweep 2 as variants that compile |

**Sweep 2**, 12:01:48–12:03:24, after the tests of the survivors were added (commit `0ca5ff1`), 8 mutants: **8 killed by assertion**
(P7, P8, P8b, I1 sorted downwards, I3, I7, I8, I11); the control exit follows below. I7 is killed through the helper's fatal on a
missing `parameters.tsv` (`git show` of a file the mutant did not write), not through a text assertion.
The integration control after sweep 2 (the 13 focused tests) exit 0, the git status clean (12:04:39).

## The checks at the head `a24e14d`

`go build ./...` exit 0, `go vet ./...` exit 0 and `go test -count=1 ./...` (the unit level, all packages) exit 0, 12:04:45–12:05:07;
`adr-lint`, `prd-lint`, `link-lint`, `run-discipline-tests` and `setup-check` exit 0, and `git diff --check` exit 0, to 12:06:44.
**Not run here, by the instruction of the root:** the full integration level (`go test -tags=integration ./...`), which the root runs
once on the stable head, the freeze and the review. The focused integration runs are the 13 tests of `internal/run`, `internal/handoff` and
`internal/tsv` (whole packages) and `TestPackageRules` of `cmd/layup`, each with its time above.
