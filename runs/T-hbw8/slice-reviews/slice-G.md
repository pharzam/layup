# Slice G review record

Reviewer: Claude Fable 5.1 (`claude-fable-5-1`), a fresh `claude -p` session in a disposable clone, read-only tools, no forge token, `--setting-sources project,local`. Not a gate round.

## Round 1

2026-09-29, head `4a58197`; 6 min 11 s; 28,122 output tokens; USD 3.46 at list price. Word for word:

~~~~text
# Slice G review — round 1
Reviewer: claude-fable-5-1. Head: 4a58197. Verdict: `material findings: 6`.

## Material findings

### M1. The spend check has no rule for "may reach `B`", and a session with no cap starts unchecked
- Where: `docs/architecture.md:1149`: "the project total has an unknown part and may reach `B` | escalate to the idea owner". Also `:1146`: "the known total plus the new session's cap would pass the milestone's cap". And W-10 step 6, tagged `code`.
- Source: FT2; lens 1 (a `code` step names its rule); K50, P10.
- Why it matters: the money of an unknown amount has no bound, so code cannot decide "may reach". Every unknown part either always escalates or never does. After `T-21` ends, its money stays `unknown` for good. So every later session start escalates again, and nothing says what the idea owner's answer records to stop the repeat. Also, for a new session with no cap, "the new session's cap" does not exist. W-10 step 5 starts `T-21` with no check, so at its own start its spend counts as nothing.
- Fix: say that an unknown part always escalates. Say what the answer writes (for example an accepted bound per `harness.<id>.wall` period, recorded in `budget.tsv`) and that it clears the escalation. Say that a new session with no cap counts as an unknown part before it starts.

### M2. Missed Escalations passes when there is no sample, and it uses the wrong window
- Where: `:1160`: "confirmed misses in the sample; it must be zero". `:1167-1169`: "agent answers merged at least 30 days before …, and screen rows and agent decisions of the milestone". `:1173`: "A later reversal writes its time and actor into `questions.tsv`."
- Source: FT1; `F-0003#57` ("the audit uses … the window of the Reversal Rate"); `#72`; K56.
- Why it matters: before the first retrospective, or with `audit.n` = 0, the sample is empty. Zero misses then meets the target "0", so a check that did not run counts as a pass. The Missed Escalations sample reads decisions "of the milestone", not those 30 days after the merge that #57 names. Retrospectives only sample answers merged at least 30 days earlier, so answers from the last month of a delivery are never audited. No actor or mechanism detects "a later reversal".
- Fix: an empty or missing sample prints "not measured", never a pass. Use the 30-day window for both samples, and add a final audit 30 days after the last merge. Name the mechanism that writes a reversal.

### M3. The Task Intervention Rate population has no source for restarts and several other human actions
- Where: `:1156`: "every human action on the target (each copied comment, each review, each push or merge by a human from the repository activity, each parameter change)", with the kind "restart". And `:136`: "a reaction or an edit … Each of them is still recorded as human input (section 12)."
- Source: `F-0001#28` ("an answer, a correction, a restart, or a change to a gate"); K53; the GitHub REST "List repository activities" endpoint, which returns only `push`, `force_push`, `branch_creation`, `branch_deletion`, `pr_merge` and `merge_queue_merge`.
- Why it matters: when a human re-runs a failed CI job to clear `ci.T`, that is a restart. It is visible only in the Actions run's `triggering_actor`, which no rule reads. The same gap covers a human who closes or reopens a PR or issue, marks a draft ready, or kills a session on the LAYUP host. §3 copies a comment only "before any step acts on" it, so a human comment that no step acts on is never copied. Reactions, which §3 sends to §12, are missing from §12's list. So "restart" is answered by name only, and the rate reads low.
- Fix: list every source (all issue and PR comments, the issue and PR events, workflow runs by `triggering_actor`, reactions, and host-side commands through `layup run`), each with its API. Record anything that cannot be seen as a known limit.

### M4. The wall-clock side of the budget has no clock rule, and the appetite has no role
- Where: `:1130`: "the appetite of the delivery in money and wall-clock, the budget `B` and the upper edge `U`". `:1146`: "or its wall-clock cap is reached". W-10 step 1: appetite "USD 200 and 40 hours", `B` = 180, `U` = 240.
- Source: O-68; D16 (appetite); §11 "Waits are not stalls" (the clock stops for a human).
- Why it matters: the text does not say when the milestone's wall-clock starts, or whether human waits count. In W-10 step 6, the escalation waits overnight for the idea owner. The 6-hour cap of milestone 2 then trips the breaker at the next start, only because a human took time to answer. The 40-hour appetite is checked nowhere: the bets' caps are checked "at most `U`", which is money only. The appetite's USD 200 differs from `B` and is used by no rule.
- Fix: define the milestone clock (its start, and that waits for a human pause it, as in §11). Check the bets' wall-clock caps against the delivery's wall-clock appetite. Say whether the appetite is `B` or drop the third number.

### M5. Values the stop reads have no home, no default and no evidence
- Where: `:1141`: "`harness.<id>.wall`, a parameter". `:1169`: "`audit.n` (a parameter set with evidence)". W-10 step 3: "the harness row (reports tokens; cap USD 5 per session)".
- Source: §10 "Every value that the Operator or the idea owner can set is a row of `parameters.tsv`" (its table has neither row); Invariant 4; FT6. §9 puts the harness register "on the LAYUP host" and gives it only "whether it can enforce a spend cap".
- Why it matters: the per-session cap drives the spend check, but nothing says who sets the amount, with what evidence, or where it lives. It is on the host, so a clone of the target cannot replay the stop decision. The ledger row in `:1113-1119` does not record the cap. W-10's 30 minutes has no source.
- Fix: add `harness.<id>.wall`, `harness.<id>.cap` and `audit.n` to the §10 table, each with a default or "set at Intake with evidence". Record the cap applied in each session's start row.

### M6. Cost per Requirement changes the PSB's definition without a recorded reading
- Where: `:1163`: "a task's money split equally over its requirement IDs".
- Source: `F-0003#74`: "Sum of the token cost of all tasks of a requirement"; `F-0004#18` (the target is the median).
- Why it matters: take a task of USD 10 that serves REQ-1 and REQ-2. The PSB's sum gives each requirement USD 10; the design gives USD 5. The idea owner's baseline, measured "with the current process", may use the PSB's rule, so the comparison is skewed. The rule also never says "median", and it never says how smart-if costs in `decisions.tsv` join a task. §8 records its reading of `#58` for the Operator to confirm, but this change of meaning has no such record.
- Fix: follow #74, or record the split as a reading for the approval brief. State the median and how smart-if costs are attributed.

## Notes
- N1. Clarification Turnaround measures "accepted minus asked", which takes in the asker's whole next attempt. L-D1 and L-D3 (slice D) already record this. "Answered minus asked, over answers later accepted" is closer to `#71`.
- N2. `billing type` is recorded, but no rule says what money a subscription session has. ev03's lesson was that a subscription recorded as 0 made the stop never trip. Claude Code's reported cost is an estimate, not a bill.
- N3. Claude Code's `--max-budget-usd` works in print mode only and is checked between steps, so one session can pass its cap a little. "Known total plus caps" is therefore not a hard ceiling. Say so in the ADR's consequences.
- N4. The §14 row for K53–K57 cites ADR-0014, which does not decide these records or the audit. ADR-0024's decision does not name them either, so the audit design has no ADR.
- N5. Early Question Share: budget escalations and floor escalations to the idea owner are not question rows. Say whether they count as "human questions".
- N6. First-Review Acceptance: "a rejection becomes a new need". If that need gets a new REQ ID, its next review is review 1 again. Keep the count on the original requirement.
- N7. A budget escalation answer is planned only if it is "confirmed as business-forking", but `F-0001#13` already makes budget business-forking. State that it is planned.
- N8. Nothing says what happens when "the caps of the bets so far are at most `U`" fails at a bet.

## Checklist rows
- S10: answered.
- R06: answered; its pilot criterion of 1 is limited by L-G1 (known limit, acceptable).
- K11: answered.
- K50: answered in mechanism, but the rule is undefined (M1).
- K51: known limit L-G1 (acceptable).
- K52: answered.
- K53: not answered for restarts and non-comment actions (M3).
- K54: answered (N5).
- K55: answered (N1, L-D1).
- K56: not answered in part (M2).
- K57: answered (N6).
- K58: answered (the session is the unit).
- P06: answered (`acceptance.tsv`).
- P10: answered, with defects in M1 and M5.
- P22: answered (a row per session; a row per call only where the harness streams it).
- D02: answered.
- D11: answered.
- D16 (appetite): answered for money; the wall-clock side is not answered (M4).
- FT2: not fully answered (M1).

## Existing solutions
- Paperclip's budget policies (ev03): the scope, a warning threshold, a hard stop, a check before each run and an override approval. The design repeats them well but has no equivalent of the override record that clears an incident (M1).
- LiteLLM's soft and hard budgets per key match `B` and `U`. Shape Up gives the appetite and the circuit breaker.
- OpenTelemetry's GenAI semantic conventions (`gen_ai.usage.*`) give the token field names.
- Claude Code has a spend cap (`--max-budget-usd`); Codex's `exec --json` reports tokens only, so its money is `computed`.
- For audit sampling, PCAOB AS 2315 (a recorded seed and a population hash) is followed.

Sources: [GitHub REST: List repository activities](https://docs.github.com/en/rest/repos/repos#list-repository-activities); [Claude Code budget cap investigation](https://linuxjedi.co.uk/when-the-docs-fall-short-investigating-claude-codes-budget-cap/); [Claude CLI startup flags](https://github.com/shanraisshan/claude-code-best-practice/blob/main/best-practice/claude-cli-startup-flags.md)~~~~

### The author's answer to round 1

All six material findings and notes N1 to N8 are applied in one commit.

| Finding | Fix |
| ------- | --- |
| M1 | An unknown part always escalates unless a bound that the idea owner accepted covers it; a session with no cap is an unknown part before it starts. |
| M2 | An empty sample is "not measured"; both samples use the 30-day window; a final audit 30 days after the last merge; a confirmed reversal writes its time and auditor. |
| M3 | The sources of human actions are listed with their APIs (comments, reviews, reactions, issue and PR events, workflow re-runs by `triggering_actor`, pushes, parameter changes, host commands); the rest is L-G2. |
| M4 | `U` is the money appetite; a wall-clock appetite; bets checked against both; a milestone's clock starts at its bet and pauses for human waits. |
| M5 | `harness.<id>.cap`, `harness.<id>.wall` and `audit.n` are §10 rows set with evidence; each start row records the cap and limit applied. |
| M6 | Cost per Requirement follows `F-0003#74` (a shared task counts in full), with the median; decision rows name their task. |
| N1 to N8 | Turnaround is answered minus asked over accepted answers; subscription money is `computed`, never 0; L-G3 for caps checked between steps; ADR-0024 decides the measures and the audit; escalations are question rows; the original REQ ID keeps the review count; a budget escalation is planned; a bet that does not fit is refused. |
