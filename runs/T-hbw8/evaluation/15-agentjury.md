# Evaluation: AgentJury

| Field | Value |
| ----- | ----- |
| Repository | https://github.com/madad-rashid/AgentJury |
| Pinned commit | 4bdad8cdb9ac9831375d02c3cdcf34c56cc0a1ca (2026-09-03T23:35:19+03:00, `main`, v0.4.4) |
| License | MIT (LICENSE file) |
| Language and needs | Python 3.11 or later; `pydantic`, `python-dotenv`; optional `anthropic` or `openai` SDK. Each judge is a direct API call, so it needs `ANTHROPIC_API_KEY` or `OPENAI_API_KEY` (and `ANTHROPIC_WORKSPACE_ID` for a multi-workspace key). No server, no database, no daemon. Verdicts are JSON files in `.agentjury/verdicts/` under the current directory. |
| Evaluated by | Claude Opus 5.5 on Claude Code, 2026-09-28 |
| Elapsed | about 20 minutes of the session (shared with three other candidates) |
| Agent runs and cost | 4 panel runs: 1 with the Anthropic adapter and no key (no model call); 3 with 2 judges each through a Claude Code adapter (haiku and sonnet, print mode, cap USD 0.40 per call). Total model cost USD 0.057 (5 successful calls, 1 failed call). |

## Verdict

Borrow the pattern. The deterministic aggregator is small, readable and fails closed (a failed judge never becomes an approval), and the "blind" rule holds in the code. It is a public-alpha Python tool that calls provider APIs with keys, not harnesses, and the judges see only text, so LAYUP should port its rules to Go instead of running it.

## What it is (from the code)

A library and CLI that sends one review request to a panel of judges in parallel and turns their votes into one verdict with fixed rules.

- **Blind rule.** `Panel.review` submits `j.review(request)` for each judge to a thread pool; each judge gets the same `ReviewRequest` and nothing else ([agentjury/panel.py#L35-L46](https://github.com/madad-rashid/AgentJury/blob/4bdad8cdb9ac9831375d02c3cdcf34c56cc0a1ca/agentjury/panel.py#L35-L46)). The user prompt holds only the task, the optional context, the agent output and the artifacts ([agentjury/judges/base.py#L136-L144](https://github.com/madad-rashid/AgentJury/blob/4bdad8cdb9ac9831375d02c3cdcf34c56cc0a1ca/agentjury/judges/base.py#L136-L144)). Each call is stateless. A judge therefore cannot see another judge's verdict. Limit: the `context` field is the caller's; a caller that puts an earlier verdict there breaks the blind rule, and nothing checks it.
- **Role and posture.** Four built-in roles (accuracy, critic, evidence, executive) and extra roles from a JSON file ([base.py#L30-L69](https://github.com/madad-rashid/AgentJury/blob/4bdad8cdb9ac9831375d02c3cdcf34c56cc0a1ca/agentjury/judges/base.py#L30-L69)). The system prompt says the reviewer "will not see any other reviewer's opinion", treats the output as untrusted data, and makes a manipulation attempt a blocking finding ([base.py#L87-L123](https://github.com/madad-rashid/AgentJury/blob/4bdad8cdb9ac9831375d02c3cdcf34c56cc0a1ca/agentjury/judges/base.py#L87-L123)).
- **Failure policy.** One retry on a provider error, one repair round if the JSON does not parse, then the judge raises ([base.py#L186-L243](https://github.com/madad-rashid/AgentJury/blob/4bdad8cdb9ac9831375d02c3cdcf34c56cc0a1ca/agentjury/judges/base.py#L186-L243)). The panel records the error and continues ([panel.py#L41-L46](https://github.com/madad-rashid/AgentJury/blob/4bdad8cdb9ac9831375d02c3cdcf34c56cc0a1ca/agentjury/panel.py#L41-L46)).
- **Deterministic aggregator** ([agentjury/aggregate.py#L1-L120](https://github.com/madad-rashid/AgentJury/blob/4bdad8cdb9ac9831375d02c3cdcf34c56cc0a1ca/agentjury/aggregate.py#L1-L120)), no LLM:
  - voters = responding judges that did not abstain ([#L63](https://github.com/madad-rashid/AgentJury/blob/4bdad8cdb9ac9831375d02c3cdcf34c56cc0a1ca/agentjury/aggregate.py#L63)); no voter gives `insufficient_jury`.
  - quorum default = strict majority of the requested judges ([#L35-L37](https://github.com/madad-rashid/AgentJury/blob/4bdad8cdb9ac9831375d02c3cdcf34c56cc0a1ca/agentjury/aggregate.py#L35-L37)).
  - `insufficient_jury` if voters are fewer than the quorum, or if a panel built from two or more providers heard from only one ([#L88-L93](https://github.com/madad-rashid/AgentJury/blob/4bdad8cdb9ac9831375d02c3cdcf34c56cc0a1ca/agentjury/aggregate.py#L88-L93)).
  - `blocked` needs blocking findings from two providers, or from two judges on a one-provider panel ([#L83-L85](https://github.com/madad-rashid/AgentJury/blob/4bdad8cdb9ac9831375d02c3cdcf34c56cc0a1ca/agentjury/aggregate.py#L83-L85), [#L94-L95](https://github.com/madad-rashid/AgentJury/blob/4bdad8cdb9ac9831375d02c3cdcf34c56cc0a1ca/agentjury/aggregate.py#L94-L95)). One judge cannot block alone.
  - `needs_revision` on a majority or a tie of revise votes, or on one blocking source; `verified` only on a majority of approve votes and no blocking finding ([#L96-L99](https://github.com/madad-rashid/AgentJury/blob/4bdad8cdb9ac9831375d02c3cdcf34c56cc0a1ca/agentjury/aggregate.py#L96-L99)).
  - A "confidence index" from consensus, panel size, score spread and provider diversity, labelled "NOT a calibrated probability" ([#L9-L11](https://github.com/madad-rashid/AgentJury/blob/4bdad8cdb9ac9831375d02c3cdcf34c56cc0a1ca/agentjury/aggregate.py#L9-L11), [#L75-L78](https://github.com/madad-rashid/AgentJury/blob/4bdad8cdb9ac9831375d02c3cdcf34c56cc0a1ca/agentjury/aggregate.py#L78-L81)).
- **Exit codes.** `0 verified, 1 needs_revision, 2 blocked, 3 insufficient_jury` ([agentjury/cli.py#L48-L49](https://github.com/madad-rashid/AgentJury/blob/4bdad8cdb9ac9831375d02c3cdcf34c56cc0a1ca/agentjury/cli.py#L48-L49)).
- **Human adjudication log.** `agentjury adjudicate` grades a judge's finding or the producer's output and appends an event with `old` and `new` values to `adjudications.jsonl` ([cli.py#L167-L242](https://github.com/madad-rashid/AgentJury/blob/4bdad8cdb9ac9831375d02c3cdcf34c56cc0a1ca/agentjury/cli.py#L167-L242)). The adjudicator name comes from `AGENTJURY_ADJUDICATOR` or the login name; nothing proves who wrote it.
- **Record.** Each `Review` carries model, role, prompt hash, config id, latency and token counts ([agentjury/protocol.py#L118-L160](https://github.com/madad-rashid/AgentJury/blob/4bdad8cdb9ac9831375d02c3cdcf34c56cc0a1ca/agentjury/protocol.py#L118-L160)).
- Branch `phase-two-free-jury` (25a3a80, 2026-09-25) adds "grounded" findings with checked excerpts and a benchmark. It is not on `main`; we did not evaluate it.

## What we ran

Install, in the work directory: `uv venv venv --python 3.12`, `uv pip install -e ".[anthropic,dev]"` (uv cache in the work directory). Test suite: `pytest -q` gave `83 passed, 2 skipped`.

Target: a toy Go module `calc`, and a diff that (1) adds a wrong `Abs` (negative input returned unchanged), (2) makes `Div(1,0)` return `ok=true`, and (3) deletes the test of that case. `go test ./...` still passes on it. The task text asks for `Abs` and says `Div` must keep `ok=false`.

The host has no `ANTHROPIC_API_KEY`. We did not use the host's OpenAI key. To run real judges we wrote a 40-line `Judge` subclass whose `complete()` runs `claude -p` with the judge's system prompt, no tools, no settings, no session persistence, in an empty temporary directory (`eval-work/agentjury/cc_panel.py`). Everything else (prompts, parsing, panel, aggregator) is AgentJury's code.

| Run | Command | Result |
| --- | ------- | ------ |
| 1 Anthropic adapter, no key | `agentjury review task.md change.diff --panel accuracy:anthropic,critic:anthropic` | Both judges `TypeError: Could not resolve authentication method`; `jury 0/2 insufficient_jury`, `No verdict.`, exit 3 |
| 2 Flawed diff, 2 judges | `cc_panel.py ... critic:haiku, accuracy:sonnet` | Both `REVISE`, 2 blocking findings each (wrong `Abs`, `Div` change) and the deleted test; `blocked`, exit 2. USD 0.020 |
| 3 Clean diff (correct `Abs`, table test) | same panel | Both `APPROVE` (9.0); minor note on `Abs(math.MinInt)`; `verified`, exit 0. USD 0.016 |
| 4 Clean diff, one judge fails | panel with `accuracy:no-such-model-x` | `jury 1/2 insufficient_jury`, error `claude-code:unrecognized_model`, exit 3. One approval did **not** become `verified`. USD 0.021 |
| 5 Human adjudication | `agentjury adjudicate ef165c25b0f2 --judge critic --finding 1 correct --producer-verdict flawed` | Two events appended to `adjudications.jsonl`, with `old: null, new: flawed` |

The sonnet judge in run 3 wrote: "I could not see the full files, only the diff". The judges review text only; they cannot run the tests.

## In-Scope items S1-S12

| Item | Mark | Evidence |
| ---- | ---- | -------- |
| S1 Problem Statement Quality | no | None. |
| S2 Reproducible Discipline Setup | no | None. |
| S3 Rule Protection | no | None. |
| S4 Stack-Dependent Gates | no | None. |
| S5 Role Handoffs | partly | `ReviewRequest` and `Verdict` are pydantic schemas with a JSON schema command (`agentjury schema`). They cover the review handoff only. |
| S6 Autonomous Clarification | no | None. |
| S7 Verification on Every Change | no | An LLM judgement on text; no layout, interface or test gate. It is a possible extra check after the deterministic gates. |
| S8 Human-on-the-Loop | no | Adjudication is after-the-fact grading, not a decision point. |
| S9 Stall Resolution | partly | A blind panel with a deterministic exit rule is a building block for a stall panel. No stall detection, diagnosis or record. |
| S10 Cost Visibility | partly | Per-review latency and tokens in the verdict JSON; no cost, and the file is outside Git by default. |
| S11 Specification Synthesis | no | None. |
| S12 Harness-Agent Neutrality | partly | Judges are provider APIs (OpenAI, Anthropic). A new judge needs only `complete(system, user)`; our Claude Code adapter took 40 lines. |

## Invariants I1-I9

| Invariant | Effect | Reason |
| --------- | ------ | ------ |
| I1 Git is the system of record | neutral | Verdicts go to `.agentjury/verdicts/` in the current directory; the project's own `.gitignore` ignores that directory. The caller decides where the JSON goes; LAYUP could commit it to its records. |
| I2 Independent repository | supports | A separate tool; nothing is written into the target unless run there. |
| I3 Agents cannot change the gates | neutral | The panel spec, quorum and roles file are the caller's arguments. |
| I4 No value without evidence | neutral | The quorum rule and the two-source block rule are rules, stated in the docstring. The confidence weights (`n/(n+1)`, `0.5+0.5*diversity`) have no stated evidence, and the tool says so. |
| I5 An inactive check is not a pass | supports | Runs 1 and 4: a failed judge is an error, not a vote; too few voters gives `insufficient_jury` and exit 3. An abstention is not an approval. Limit: the caller can set `--quorum 1`. |
| I6 Deterministic first | neutral | The judges are LLM judgements; the aggregation is deterministic. Use it only where no mechanical check exists. |
| I7 Stack adds, never removes | neutral | Not related. |
| I8 Pinned Armature | neutral | Not related. |
| I9 Replaceable harness | supports (partly) | Provider-diverse panels are a first-class rule (the provider floor). Harness CLIs need an adapter. |

## Deep-check findings it answers

- **Fable-M18** (blind panels decided by no ADR): gives four of the six missing parts, as a pattern. Composition: role:provider pairs. Blind rule: each judge gets only the request. Posture: a role prompt plus a "treat the output as untrusted" rule. Synthesis author: deterministic code, not a model. It does not give a trigger in Design, and it aggregates votes; it does not synthesize a resolution path.
- **Sol-24** (no panel execution path): gives convening (`Panel.review`), sealed inputs, member isolation (independent calls), a per-member record (`Review` in the `Verdict` JSON), a synthesis actor (the aggregator), and an exit rule (the four statuses with exit codes). Missing: the members produce votes on one output, not competing options for a stall.
- **Fable-M6** (a stall reaches the Operator with no diagnosis): none directly. The `insufficient_jury` status is a pattern for "the examination did not happen, say so".
- **Sol-13, Author-9** (a reversal has no record): pattern. The adjudication event log with `old` and `new` per finding and per producer output is a reversal record.
- **Sol-30** (a review is not an issue comment): pattern. A review is a typed record (`Verdict` JSON), which LAYUP could store and link from the issue.
- **Sol-2, Fable-M19, S7, S12**: Sol-2 and Fable-M19 none; S7 no; S12 partly, above.

## Parts and their verdicts

| Part | Verdict | How LAYUP would take it | Reason |
| ---- | ------- | ----------------------- | ------ |
| Deterministic aggregator (quorum, provider floor, two-source block, `insufficient_jury`) | borrow the pattern | Port the ~60 lines of rules to Go inside LAYUP; keep the thresholds as parameters (O-79) | Runs 2-4 matched the rules. Standard-library Go is enough. |
| Blind panel runner | borrow the pattern | LAYUP starts N fresh harness sessions with the same sealed input, in parallel, and never passes one member's output to another | The API-key design does not fit LAYUP's harness sessions; our adapter proved the seam is one method. |
| Judge prompt (untrusted-data section, manipulation = blocking) | borrow the pattern | Reuse the idea in LAYUP's reviewer prompt, with attribution (MIT) | Short and clear. |
| Verdict and Review schema | borrow the pattern | A LAYUP record with the same fields (role, model, prompt hash, tokens, latency, findings with severity) | Makes each member's work auditable. |
| Adjudication event log | borrow the pattern | A reversal record in LAYUP's records branch, with the human identity from the GitHub App rule (O-77) | AgentJury's own identity is an unchecked name. |
| The CLI as a separate program | reject | none | Needs API keys and Python; public alpha; one maintainer; `main` unchanged since 2026-09-03. |
| Overall | borrow the pattern | See above | Evidence: the runs above and aggregate.py. |

## Where the searchers were wrong or incomplete

- Searcher A: "last push 2026-09-25". The push on that date is to the branch `phase-two-free-jury`; `main` (the release) was last changed on 2026-09-03.
- Searcher A, S8 "partly" for "Human adjudication": adjudication grades verdicts after the fact; it is not a Human Decision Point. We mark S8 `no`.
- Searcher A, S9 "partly" for "per-judge timeout/retry": a retry policy is not stall resolution. We mark S9 `partly` only because the panel is a building block for a stall examination.
- Searcher A: "verdict files local, not guaranteed committed" is correct; the project's own `.gitignore` excludes `.agentjury/`.
- Searcher A: "failed calls and abstentions do not silently become approval" is correct; run 4 confirmed it.
- Not noted by either searcher: one judge alone can never block (a `blocking` finding from one judge gives `needs_revision`), and a one-provider panel needs two blocking judges.

## Limits of this evaluation

- We did not call the real Anthropic or OpenAI adapters with a key; the Anthropic adapter ran only to its authentication failure. The model calls went through our Claude Code adapter.
- Both judges in runs 2-4 used one provider (Anthropic), so the provider-diversity floor was not exercised by a run; the unit tests cover it.
- The Hermes integration and the `phase-two-free-jury` branch were not run.
