# W-10 — Cost Visibility

**Item.** S10, Cost Visibility (`F-0003#50`): "Each task gets a record of its
token count, its latency, and its wall-clock duration in the project repository."
Also Problem 5, Telemetry Completeness (`F-0003#60`), Cost per Requirement
(`#74`), `REQ-011`, O-68 and O-80, and vision 3.4 (telemetry per action).

**The cases.** Build task `T-20` runs on Claude Code, which reports tokens and
cost and takes a spend cap. Build task `T-21` runs on a harness that reports no
tokens and cannot enforce a cap. Milestone 2 has a cap of USD 40 and 6 hours.

Sections are those of [`architecture.md`](../architecture.md).

| # | Actor | Input | Mechanism | Record | Tag | Where |
| - | ----- | ----- | --------- | ------ | --- | ----- |
| 1 | the idea owner | the Intake form | open | open | `human` | open |
| 2 | the idea owner | the bet of milestone 2 | open | open | `human` | open |
| 3 | `layup run` | the start of a session of `T-20` | open | open | `code` | open |
| 4 | `layup run` | the end of that session | open | open | `code` | open |
| 5 | `layup run` | a session of `T-21` | open | open | `code` | open |
| 6 | `layup run` | the spend before each start | open | open | `code` | open |
| 7 | `layup report` | the records | open | open | `code` | open |
