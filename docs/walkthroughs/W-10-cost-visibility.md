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
| 1 | the idea owner | the Intake form | answers the appetite (USD 200 and 40 hours), `B` = USD 180 and `U` = USD 240 (Decision Point 1, O-68) | the answer comment | `human` | §12 |
| 2 | the idea owner | the bet of milestone 2, whose brief proposes a cap of USD 40 and 6 hours; code checked that the caps so far stay within `U` | bets | the bet comment; the cap in `budget.tsv` | `human` | §8, §12 |
| 3 | `layup run` | the start of a session of `T-20`; the ledger; the harness row (reports tokens; cap USD 5 per session) | the known spend of milestone 2 is USD 31; plus the running sessions' caps and this session's cap of USD 5 is USD 36, under the cap of USD 40; the project total is under `B`: go on | a start row | `code` | §12; ADR-0024 |
| 4 | `layup run` | the end of that session: the harness's usage output | tokens by class `observed`; money `reported` (USD 3.10); latency and duration from the times | a ledger row | `code` | §12; ADR-0024 |
| 5 | `layup run` | a session of `T-21` (no token report, no cap) | runs it under `harness.H3.wall` = 30 minutes; tokens `unavailable` ("the harness reports none"); money `unknown`; `T-21`'s telemetry is incomplete (O-80) | a ledger row | `code` | §12; ADR-0024 |
| 6 | `layup run` | the spend before the next start | the project total has an unknown part (`T-21`) and may reach `B`; an unknown amount never counts as below `B`, so it escalates to the idea owner (no smart-if call) | an escalation (§10) | `code` | §12; ADR-0024 |
| 7 | `layup report` | the records | per task the tokens, latency and duration; `T-21` incomplete, so Telemetry Completeness fails for it (L-G1); Cost per Requirement `partial` | the report | `code` | §12 |

## Checklist rows

S10, R06, K50 to K58, P06, P10, P22, D02, D11, D16 (appetite); FT2. Known limit L-G1.
