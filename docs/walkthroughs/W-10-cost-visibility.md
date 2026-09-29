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
| 1 | the idea owner; the Operator | the Intake form | answers `B` = USD 180, `U` = USD 240 (the money appetite) and a wall-clock appetite of 40 hours (Decision Point 1, O-68); the Operator's answer sets `harness.H3.wall` = 30 minutes and the Claude Code cap of USD 5 per session, each with its evidence | the answer comment | `human` | §12 |
| 2 | the idea owner | the bet of milestone 2, whose brief proposes a cap of USD 40 and 6 hours; code checked that the money caps so far stay within `U` and the wall-clock caps within 40 hours | bets | the bet comment; the cap in `budget.tsv` | `human` | §8, §12 |
| 3 | `layup run` | the start of a session of `T-20`; the ledger; the harness row (reports tokens; cap USD 5 per session) | the known spend of milestone 2 is USD 31; plus the running sessions' caps and this session's cap of USD 5 is USD 36, under the cap of USD 40; the project total is under `B`: go on | a start row with the cap applied | `code` | §12; ADR-0024 |
| 4 | `layup run` | the end of that session: the harness's usage output | tokens by class `observed`; money `reported` (USD 3.10); latency and duration from the times | a ledger row | `code` | §12; ADR-0024 |
| 5 | `layup run` | a session of `T-21` (no token report, no cap) | before its start it is an unknown part; the idea owner had accepted "up to USD 3 per session on H3", so it counts as USD 3 and starts; runs it under `harness.H3.wall` = 30 minutes; tokens `unavailable` ("the harness reports none"); money `unknown`; `T-21`'s telemetry is incomplete (O-80) | a ledger row | `code` | §12; ADR-0024 |
| 6 | `layup run` | the spend before the next start | a session on another harness with no cap and no accepted bound is next: an unknown part that no bound covers, so it escalates to the idea owner (no smart-if call); a budget escalation is planned input | an escalation (§10) | `code` | §12; ADR-0024 |
| 7 | `layup report` | the records | per task the tokens, latency and duration; `T-21` incomplete, so Telemetry Completeness fails for it (L-G1); Cost per Requirement `partial` | the report | `code` | §12 |

## Checklist rows

S10, R06, K50 to K58, P06, P10, P22, D02, D11, D16 (appetite); FT2. Known limit L-G1.
