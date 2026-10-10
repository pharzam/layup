# T-x7cs — test runs

The uat test is [`internal/run/demo_uat_test.go`](../../internal/run/demo_uat_test.go),
under the build tag `uat`, run by hand on the host `hetzam` with the host
directory `~/layup-host` (plan of #170, D1 and D5).

## Run 1: the restart refused at `clone` (2026-10-10, 09:06Z to 09:07Z)

`go test -tags=uat -run TestTheDemoOfM2b -timeout 30m -v ./internal/run -args -host ~/layup-host -target pharzam/layup-uat`,
with the registers of D5 and no credential on either row (the red of condition
2 of the plan review). Exit 1, before the step `probe` and before any write:

```
step forge: done, the installation token, the six capabilities and the permissions of M2a; a public repository, with layup-records
step clone: fail, start/start.tsv: line 13, column "name": the next name is harness.claude.cap (the names of the block, in its order, with the harnesses of the register)
```

The target was started in `M2a` with a harness register of no row; the restart
reads `start.tsv` against the register now, which has two. This is not the
planned red: it is the defect #199, revealed here, off the path of this task.
The demo waits for the Operator's choice (comment 6095966713 of #170).
