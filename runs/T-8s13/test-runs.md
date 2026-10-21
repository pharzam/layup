# T-8s13: test runs of the bound script (step T)

The test `brief-bound-test.sh` (nine cases) pins the terms and `L` that `brief-bound.sh` prints, on small records it builds in a temporary
plain git repository (`refs/remotes/origin/layup-records`). The formula is the one of `internal/run/betbrief.go` (`assemble`, with
`requirementsPart`, `solutionPart`, `briefOpenPart` and `briefFactsPart`) and the plan review (condition 1):
`L = 12 + R + max(M,1) + sum(1+O_d) + Q_io + E + S + A_io + C_op`.

## Red: the stub prints wrong terms

Command (from the root of the worktree at `76dfb684` plus these two files): `sh runs/T-8s13/brief-bound-test.sh`

- Start 2026-10-21T16:59:18Z, end 2026-10-21T16:59:18Z (read from `date -u` before and after the command).
- Exit 1. Nine cases, nine `FAIL`, no `ok`. Each `FAIL` is an assertion: the stub `brief-bound.sh` prints
  `R=0 M=0 S=0 A_io=0 C_op=0 Q_io=0 E=0 L0=0` and the case wants its own terms (the first: `R=2 M=0 S=2 A_io=0 C_op=1 Q_io=1 E=0 L0=19`
  and the line `G1 floor=22 limit=80 stop=no`). The reason is the intended one: the stub does not count.

## Green: the script counts

The implementation took four runs, all of the same command (`sh runs/T-8s13/brief-bound-test.sh`): the first (2026-10-21T17:00:03Z, exit 1) failed the case of the follow-ups on a
defect of the implementation (a quote compared as a string in awk, then a lookup in the wrong array); the fixtures were then corrected to be valid records (see the cross-check below).

Final run: start 2026-10-21T17:01:13Z, end 2026-10-21T17:01:13Z (`date -u`), exit 0, nine cases, nine `ok`, no `FAIL`:

```text
ok only an idea-owner question open: E=0, Q_io=1
ok no question open: E=1 for the line - none
ok an Operator and an idea-owner question open: E=1, Q_io=1
ok answered follow-up closes its question, asked one keeps it open; C_op by quote class
ok idea-owner question with an unanswered follow-up is open
ok G2 at L=80: room 20 passes
ok G2 at L=81: room 19 fails
ok G1 floor 79 (76 plus one decision of two options) passes
ok G1 floor 81 stops
```

The cases pin: the open-question case of plan condition 1 (E=1 when an Operator question is still open or none is open, E=0 when only idea-owner questions are open;
"still open" by `stillOpen` of `internal/run/shapequestions.go:25`, an answered follow-up closes its question and an asked one does not); C_op by the quote class of the
question an answer asks (a follow-up takes the class of the question it asks again); the 80/81 boundary of G2 (L=80 gives room 20, pass; L=81 gives room 19, fail); and the
G1 floor of 79 (L0=76 plus one decision of two options), with 81 stopping.

### Mutants (each a copy of the script, the test run on it)

- E always 1: 2 cases fail. C_op always 1: 7 cases fail. G2 room `> 20` instead of `>= 20`: the L=80 case fails.

### Cross-check against `internal/run/betbrief.go` (a throwaway Go test, not committed, deleted after the run)

The same fixtures (kept with `KEEP=DIR`) were read by `records.ReadDraft`, `ReadSpans`, `ReadQuestions` and `ReadAnswers` (which also proved them valid records), and
the brief's parts were built by `requirementsPart`, `solutionPart` (the two heads of `betbrief_test.go`, `filterAll` on the options table), `briefOpenPart` and `briefFactsPart`;
`strings.Count(assemble(parts, "", "", ""), "\n")` is the `L` that `overLimit` compares with `brief.lines`. The Go count and the script's `L` (with HEAD) or `L0` agree on every case:
followups 20, g1-79 76, g1-81 78, g2-80 80, g2-81 81, io-follow 17, io-open 19, none-open 20, op-open 19.

## Supplementary: the M2d records (read-only)

`git clone https://github.com/pharzam/layup-uat-m2d.git`, `origin/layup-records` at `d25f4df469d21d345c278af1eae73a6224efc92a`; `sh runs/T-8s13/brief-bound.sh CLONE` at 2026-10-21T17:00:51Z, exit 0:

```text
R=18 M=3 S=25 A_io=14 C_op=3 Q_io=0 E=1 L0=76
G1 floor=79 limit=80 stop=no
```

76 before decisions, as the plan predicted (R=18, M=3, S=25, A_io=14, C_op=3, Q_io=0, E=1); the floor is 79, so G1 passes with room 21. The Go functions above, run on the
same four files of that branch, also count 76. This is a comparison, not the test: the target of the demo is `pharzam/layup-uat-m2e`.
