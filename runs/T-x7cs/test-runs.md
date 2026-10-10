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

## Run 2: the demo, green (2026-10-10, 14:58:11Z to 14:59:05Z)

On the new target `pharzam/layup-uat-m2b`, started at 14:52Z to 14:58Z with
the registers of D5 and both credentials: exit 0, `--- PASS:
TestTheDemoOfM2b (52.41s)`; the output and the records are in
[`README.md`](README.md).

**A deviation:** the red of condition 2 of the plan review (each probe
`failed` with no credential) was not run. Run 1 stopped at `clone` (#199)
before any probe, and on the new target a red after the green would leave a
`failed` row as the last row of each harness at its version, which admits no
session. The red of the probe's rules is that of row 38 (`runs/T-nxe4/`); this
run is the demo on real harnesses.
