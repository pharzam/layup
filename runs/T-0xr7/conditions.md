# T-0xr7: the conditions of the plan review of #493, and their evidence

One row per condition of the plan review (6348315001).

| Condition | Evidence |
| --- | --- |
| 1. The first slice is the red tests: the three rendered prompts, the copied-command flow failing on the old behaviour, a binary scenario, and the updated exact wait-message assertion | `runs/T-0xr7/test-runs.md:3`, `internal/run/instructions_test.go:718`, `internal/run/shapeloop_integration_test.go:1032`, `internal/run/shapeloop_integration_test.go:170`, `cmd/layup/loop_e2e_test.go:452` |
| 2. The exemption is limited to the deferral itself; a pending section is still checked; the single-survivor check is explicit in all three summaries | `internal/run/instructions.go:58`, `internal/run/instructions.go:71`, `internal/run/instructions.go:120`, `docs/spec/loop.md:1804`, `internal/run/instructions_test.go:718` |
| 3. Eligibility from the recorded history: a not-mergeable cap round grants nothing and consumes no grant; an accepted grant stays accepted; a command before the cap stays refused | `internal/run/cycle.go:73`, `internal/run/cycle.go:97`, `internal/run/cycle_test.go:119`, `internal/run/cycle_test.go:147`, `internal/run/cycle_test.go:180`, `docs/spec/loop.md:2844` |
| 4. Every refusal is reachable: candidates are examined before the task filter, the order is stated, the `id(T-xxxx)` form is used, the wrong issue and the IDs below, equal to and above the cap comment are tested, and records.md says a valid copy is the grant | `internal/run/cycle.go:19`, `internal/run/cycle.go:51`, `internal/run/cycle.go:155`, `internal/run/cycle_test.go:65`, `internal/run/cycle_test.go:119`, `internal/run/shapeloop_integration_test.go:1068`, `docs/spec/records.md:269` |
| 5. The demo fixture is pinned at layup-records e65db3e30027cee5e74f466c9fc88897fd60a8b3; old records are an unchanged prefix; attempt 4 starts at d03c8650303fbdff49986ef5478ed53d11219964 and its developer receives the attempt-3 findings | `internal/run/cycledemo_test.go:49`, `internal/run/cycledemo_test.go:27`, `internal/run/testdata/t-am8b.txtar:1` |
| 6. The budget excludes only run evidence; fixtures and expected outputs count | 1042 lines over 20 files against `3efd52b8`, excluding `runs/T-0xr7/` (maximum 900/21, tolerance 180 lines and 4 files); `internal/run/testdata/t-am8b.txtar:1` |
