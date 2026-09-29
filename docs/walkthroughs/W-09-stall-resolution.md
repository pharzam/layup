# W-09 — Stall Resolution

**Item.** S9, Stall Resolution (`F-0003#49`): "A disagreement between role
agents, or a step that repeats without progress, gets a procedure with a limit.
An independent examination with a fresh context gives a diagnosis. The project
repository keeps the diagnosis and the outcome." Also Decision Point 5
(`F-0001#14`), Stall Diagnosis (`F-0003#61`), `REQ-009`, `REQ-010`, O-82, and
vision 3.2 and 3.5.

**The cases.** Part 1: in task `T-9`, the verifier on Codex rejects the change
because the invoice total excludes tax; the developer on Claude Code changes the
test, not the code, and the verifier rejects it again for the same finding.
Part 2: in task `T-11`, a developer session stops writing any output.

Sections are those of [`architecture.md`](../architecture.md).

## Part 1 — a disagreement

| # | Actor | Input | Mechanism | Record | Tag | Where |
| - | ----- | ----- | --------- | ------ | --- | ----- |
| 1 | the verifier | the head | rejects: a finding at `invoice/total.go:40`, rule "the total includes tax (REQ-3)" | a review record | `model` | §8 |
| 2 | the developer session | the finding | changes the test, not the code; commits | a result | `model` | §8 |
| 3 | the verifier | the new head | rejects again: the same file, line and rule | a review record | `model` | §8 |
| 4 | `layup run` | the task's hill numbers | the same finding is still open, and no test of the frozen list passed that did not before: a round without progress; with `stall.N` = 1, a stall | a stall row: trigger "no progress" | `code` | §11; ADR-0023 |
| 5 | `layup run` | the stall | builds the package: the plan, both review records, the hill history, the diffs as payloads; no session's own reasoning | the package payload | `code` | §11 |
| 6 | an examiner session | the package | cause "disagreement": the plan's test contradicts the specification's rule on tax; recommends the panel, since a retry repeats the error | the diagnosis | `model` | §11 |
| 7 | the smart-if (P3) | the diagnosis | under `shadow`, the ladder decides: the retry rung is used up by step 2's round, so the panel | a `decisions.tsv` row | `model` | §10, §11; ADR-0021 |
| 8 | three panel sessions, on two harnesses; a synthesis session | the sealed package and diagnosis | each writes hypotheses and options; the synthesis merges them: "correct the test to the specification, then the code" | the members' outputs; the synthesis | `model` | §11; ADR-0023 |
| 9 | `layup run` | the panel's outputs | quorum: three valid outputs from two harnesses | the panel row | `code` | §11 |
| 10 | `layup run` | the resolution | starts the next attempt with it; the verifier passes; the stall closes | the outcome row: "closed without a human" | `code` | §11 |

## Part 2 — a session that hangs

| # | Actor | Input | Mechanism | Record | Tag | Where |
| - | ----- | ----- | --------- | ------ | --- | ----- |
| 11 | `layup run` | the running session's output and hook events | none for `stall.T` = 10 minutes (not a wait state): kills it; a stall | a stall row: trigger "hang" | `code` | §11; ADR-0023 |
| 12 | an examiner session | the package: the harness's exit, its log | cause "a harness failure"; recommends a retry | the diagnosis | `model` | §11 |
| 13 | `layup run` | the retry, which hangs again | one rung up; with one other harness only, the panel's quorum cannot be met (`insufficient panel`), so the Operator | the ladder row | `code` | §11 |
| 14 | the Operator | the stall package | answers `reroute T-11 to H2` (Decision Point 5) | the comment | `human` | §11 |
| 15 | `layup run` | the Operator's answer | copies it (§3); writes the harness override; the next attempt starts from the base on H2 (W-12) | the outcome row: "closed by the Operator" | `code` | §3, §11 |

## Checklist rows

S9, R05, K42 to K49, P05, P19, P23, C8 (the circuit breaker), D13, D16 (the breaker and the hill); FT1, FT6. Known limit L-F1.
