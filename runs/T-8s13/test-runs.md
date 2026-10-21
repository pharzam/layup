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
