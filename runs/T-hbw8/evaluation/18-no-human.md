# Evaluation: no_human

| Field | Value |
| ----- | ----- |
| Repository | https://github.com/no-human-ai/no_human |
| Pinned commit | 9b50bf6467e8beec62df7ebb625338aab152114e (2026-09-24T23:50:37+03:00) |
| License | MIT (LICENSE file). The repository also has CLA.md (for contributions) and TRADEMARK.md (for the name). |
| Language and needs | Python 3.12 or later, `uv`; `claude-agent-sdk` with a bundled Claude Code CLI (191 MB); `git`. The full loop also needs a local board server (`127.0.0.1:8420`), a SQLite database `~/.no_human/no_human.db`, and Node for the web board. Credential: a Claude subscription token (`CLAUDE_CODE_OAUTH_TOKEN` from `claude setup-token`) in `~/.no_human/.env`, or an API key in `api_key` mode; Codex and local backends exist. Telemetry to PostHog and a PyPI update check are on by default. |
| Evaluated by | Claude Opus 5.5 on Claude Code, 2026-09-28 |
| Elapsed | about 25 minutes of the session (shared with three other candidates) |
| Agent runs and cost | 4 runs of `nh gate` (one-shot review gate): 3 reviewer calls on `claude-sonnet-5`, 1 simulated failure with no model call. `nh gate` does not print a cost. |

## Verdict

Borrow the pattern. The one-shot gate did what it says on a Go diff: a fresh, single-turn, read-only reviewer on a different model found all four defects with file and line, passed the clean change, and failed closed (exit 2) when the reviewer could not answer. But its tamper guard does not see Go test files at all ("tests 0->0" on every run), its full loop keeps state in a local database, and it needs a Claude token in its own home directory; LAYUP should rebuild the reviewer contract in Go.

## What it is (from the code)

A local "coding factory": a ticket becomes a plan, a coder session, tests, an adversarial review, and a pull request that a human approves. We evaluated the review part, which also ships as a one-shot command (`nh gate`) and a GitHub Action.

- **The reviewer contract.** "A fresh-context Agent SDK session told to find faults and refute 'done'", on the configured `llm.review_model`, "a different model from the implementer", with file-edit tools, git and forge writes and subagents refused. It never returns a numeric score; a session with no parseable verdict "never passes" and raises `ReviewerUnavailable` after one retry ([src/no_human/review/reviewer.py#L1-L17](https://github.com/no-human-ai/no_human/blob/9b50bf6467e8beec62df7ebb625338aab152114e/src/no_human/review/reviewer.py#L1-L17)). A missing verdict block is "fail closed" ([#L2104-L2116](https://github.com/no-human-ai/no_human/blob/9b50bf6467e8beec62df7ebb625338aab152114e/src/no_human/review/reviewer.py#L2104-L2116)).
- **The one-shot gate.** `nh gate` ([src/no_human/cli/commands.py#L6092-L6133](https://github.com/no-human-ai/no_human/blob/9b50bf6467e8beec62df7ebb625338aab152114e/src/no_human/cli/commands.py#L6092-L6133)) calls `run_gate` ([src/no_human/review/oneshot.py#L512](https://github.com/no-human-ai/no_human/blob/9b50bf6467e8beec62df7ebb625338aab152114e/src/no_human/review/oneshot.py#L512)): it diffs the branch against its merge base, runs the tamper guard ([#L621](https://github.com/no-human-ai/no_human/blob/9b50bf6467e8beec62df7ebb625338aab152114e/src/no_human/review/oneshot.py#L621)), clones the head into a temporary directory, and runs the reviewer `single_turn` ([#L755](https://github.com/no-human-ai/no_human/blob/9b50bf6467e8beec62df7ebb625338aab152114e/src/no_human/review/oneshot.py#L755)). A reviewer with no verdict becomes `GateUnavailable` ("the gate did not run"), exit 2 ([#L655-L665](https://github.com/no-human-ai/no_human/blob/9b50bf6467e8beec62df7ebb625338aab152114e/src/no_human/review/oneshot.py#L655-L665)). The module states it only reads the user's checkout ([#L1-L49](https://github.com/no-human-ai/no_human/blob/9b50bf6467e8beec62df7ebb625338aab152114e/src/no_human/review/oneshot.py#L1-L49)). It does not run the tests.
- **The tamper guard.** It counts tests, assertions, skips, tautologies and faking fixtures in test files, before and after ([src/no_human/testing/tamper_guard.py#L1-L20](https://github.com/no-human-ai/no_human/blob/9b50bf6467e8beec62df7ebb625338aab152114e/src/no_human/testing/tamper_guard.py#L1-L20)). A test file is a Python, JavaScript/TypeScript or Java name, or any file under `test/`, `tests/`, `__tests__/`, `e2e/` ([#L32-L50](https://github.com/no-human-ai/no_human/blob/9b50bf6467e8beec62df7ebb625338aab152114e/src/no_human/testing/tamper_guard.py#L32-L50)). `*_test.go` is not in the list. Test and assertion patterns are Python, JavaScript and Java only ([#L76-L96](https://github.com/no-human-ai/no_human/blob/9b50bf6467e8beec62df7ebb625338aab152114e/src/no_human/testing/tamper_guard.py#L76-L96)).
- **Coder and reviewer on different models.** `_collides_with_the_review_gate` refuses `review_model == primary_model` ([src/no_human/core/model_catalog.py#L241-L257](https://github.com/no-human-ai/no_human/blob/9b50bf6467e8beec62df7ebb625338aab152114e/src/no_human/core/model_catalog.py#L241-L257)). This check is on the model-selection path; our `config.yaml` set the reviewer to `claude-sonnet-5`, which is also the default `primary_model`, and the gate ran without a warning.
- **Bounded loop.** `max_attempts: 3` per loop, `lifetime_attempts: 9`, a cost-weighted token cap; exceeding a lifetime cap "raises a BUDGET_EXHAUSTED blocker — an honest park" ([src/no_human/config.py#L1763-L1784](https://github.com/no-human-ai/no_human/blob/9b50bf6467e8beec62df7ebb625338aab152114e/src/no_human/config.py#L1763-L1784)). Several limits carry a measured rationale in the comment.
- **Planning fan-out.** "MoA" planning runs N independent plan proposals and one aggregator call that "synthesizes a single plan (evidence-based synthesis, never a numeric score)", gated by complexity signals ([src/no_human/config.py#L1583-L1596](https://github.com/no-human-ai/no_human/blob/9b50bf6467e8beec62df7ebb625338aab152114e/src/no_human/config.py#L1583-L1596)). We read this; we did not run it.
- **Telemetry.** "Opt-OUT usage telemetry + masked session replay ... default ON" ([src/no_human/config.py#L2407-L2415](https://github.com/no-human-ai/no_human/blob/9b50bf6467e8beec62df7ebb625338aab152114e/src/no_human/config.py#L2407-L2415)).

## What we ran

Install, in the work directory: `uv sync --frozen` with `UV_PROJECT_ENVIRONMENT` and `UV_CACHE_DIR` in the work directory. `HOME` was the work directory's `home/`, with `~/.no_human/config.yaml`: `telemetry.enabled: false`, `updates.enabled: false`, `llm.review_model: claude-sonnet-5`.

Credential: no_human demands a subscription token in `~/.no_human/.env`; we did not mint one. We wrote a dummy value there, and replaced the SDK-bundled CLI in our own venv with a wrapper that removes the dummy token, uses the host's Claude Code login, maps any opus model to `claude-sonnet-5`, and logs each call. The original bundled binary is kept in the work directory.

Target: toy Go module `calc` with `origin/HEAD` set on a local bare remote, and three branches.

| Branch | Command | Result |
| ------ | ------- | ------ |
| `flawed`: wrong `Abs`, `Div(1,0)` returns `ok=true`, deleted assertion; `go test` passes | `nh gate --title "Add Abs to calc" --description "... Div must keep returning ok=false ..."` | `## no_human gate — FAIL`, 12 s, exit 1. `Tamper guard — clean - tests 0->0, assertions 0->0, skips 0->0`. Four failed items, each with file:line: `Div regression ... calc/calc.go:10`, `Abs is wrong for negative input ... calc/calc.go:17`, `Div zero-divisor test deleted ... calc/calc_test.go:18`, `Abs test is trivial ... calc/calc_test.go:27`. |
| `skip`: correct `Abs` and table test, plus `t.Skip("flaky on CI")` in `TestDiv` | same | `FAIL`, exit 1. Tamper guard again `clean ... skips 0->0`. The reviewer: `Unrelated t.Skip disables the existing Div test — calc/calc_test.go:15`. |
| `clean`: correct `Abs` and table test | same | `PASS`, exit 0. Two passed notes: `Abs(math.MinInt)` overflow, and "No test output provided ... please paste a go test ./calc run". |
| `flawed`, reviewer forced to fail (wrapper exits 1, no model call) | same with `EVAL_FAIL=1` | exit 2: `cannot run the gate: the reviewer could not reach a verdict ... The review gate did not run, so this diff is unreviewed. Escalating rather than passing it`. |

The reviewer call used `--max-turns 1`, `--permission-mode bypassPermissions` with the read-only guard, and `--setting-sources=` (no user or project settings), in a temporary clone of the head.

## In-Scope items S1-S12

| Item | Mark | Evidence |
| ---- | ---- | -------- |
| S1 Problem Statement Quality | no | Intake "grill" asks questions per ticket; not a PSB gap batch. Not run. |
| S2 Reproducible Discipline Setup | no | `nh init` sets up the tool, not a baseline. |
| S3 Rule Protection | partly | The tamper guard counts weakened tests (not for Go); the reviewer session cannot edit files. No protection of rule files. |
| S4 Stack-Dependent Gates | no | The one-shot gate collects no lint or type evidence; the tamper guard is language-specific and skips Go. |
| S5 Role Handoffs | partly | The reviewer output is a parsed `REVIEW_JSON` block with a checklist; other handoffs not checked. |
| S6 Autonomous Clarification | no | Questions go to the human ("Needs answer" lane). |
| S7 Verification on Every Change | partly | Run: an LLM review of every diff, fail-closed. The one-shot gate runs no tests; the full loop runs them (README; not run). |
| S8 Human-on-the-Loop | partly | A human approves each PR (`nh approve`); no escalation rule for business-forking decisions. |
| S9 Stall Resolution | partly | Attempt and token caps, then an "honest park" with a question or a budget record (config.py#L1763-L1784). No fresh diagnosis. |
| S10 Cost Visibility | partly | Per-attempt tokens in the local database (`nh logs`); not in the repository. |
| S11 Specification Synthesis | partly | A plan with acceptance criteria per ticket (README); no trace to a problem statement. Not run. |
| S12 Harness-Agent Neutrality | partly | Claude, Codex and local backends; the reviewer can be a different backend from the coder. |

## Invariants I1-I9

| Invariant | Effect | Reason |
| --------- | ------ | ------ |
| I1 Git is the system of record | conflicts (full loop) | Tasks, attempts and reviews live in `~/.no_human/no_human.db`. The one-shot gate keeps no state; LAYUP could store its Markdown output. |
| I2 Independent repository | supports | `nh gate` writes nothing to the checkout; the target does not depend on it. |
| I3 Agents cannot change the gates | neutral | The reviewer is read-only and cannot edit; the gate config is in the operator's home, outside the target. |
| I4 No value without evidence | supports (partly) | Many limits cite a measurement (for example the 4M token cap: "1.6M parks 117/221 real tasks"). Other defaults (`max_files: 15`) do not. |
| I5 An inactive check is not a pass | conflicts (partly) | The reviewer fails closed (exit 2, run 4), and a PR with no test command reads "NOT RUN" (README). But the tamper guard reports `clean` when it counted zero tests in a Go repository. |
| I6 Deterministic first | supports (partly) | The tamper guard runs before the LLM; for Go it counts nothing, so only the LLM judges. |
| I7 Stack adds, never removes | neutral | Not related. |
| I8 Pinned Armature | neutral | Not related. |
| I9 Replaceable harness | supports | Coder and reviewer backends are separate settings; a different model is the rule. |

## Deep-check findings it answers

- **Sol-2** (a PR with a failed gate reaches review): pattern, from the README and code comments, not run. The full loop runs tests and the review before it opens the PR, and the human approves only after that. The GitHub Action runs on an open PR, so it does not answer Sol-2 alone.
- **Fable-M6** (a stall reaches the Operator with no diagnosis): partly, as a pattern. `GateUnavailable` says in plain words that the review did not run and why; a parked task carries a specific question or a budget record. There is no fresh diagnosis by another session.
- **Sol-24, Fable-M18** (panel execution path): partly. The MoA planning fan-out (N independent proposers, one aggregator, evidence-based synthesis, no score) is an execution path for a panel of options; the reviewer itself is one member, not a panel. Read, not run.
- **Sol-27** (no enforceable PR creation and review step): pattern. Ticket, plan, attempt, review, PR, human approval are separate states with a fail-closed review between them.
- **Sol-22** (a new record resets the stall clock): pattern. Lifetime attempt and token caps count across resumes, so new artifacts do not reset them.
- **Fable-M19, S7, S12**: Fable-M19 none; S7 and S12 partly, above.

## Parts and their verdicts

| Part | Verdict | How LAYUP would take it | Reason |
| ---- | ------- | ----------------------- | ------ |
| Reviewer contract (fresh session, single turn, read-only, different model, refute "done", file:line checklist, fail closed) | borrow the pattern | LAYUP's counterpart verification starts a fresh harness session with these rules and parses one verdict block; no verdict = "did not run" | Runs 1-4 showed each property. |
| `nh gate` as a separate program | reject (for now) | Possible: `nh gate --repo <checkout> --base <sha>`, exit 0/1/2 | Needs Python 3.12, a 191 MB bundled CLI, a token in its own `.env`, telemetry off by config; it runs no tests; its tamper guard is blind to Go. |
| Tamper guard | reject for Go | isitdone's scanner covers Go | `*_test.go` is not a test file to it; `clean` on all three Go branches. |
| Bounded attempts with an honest park | borrow the pattern | LAYUP stall rule: lifetime caps across resumes, then a record with a question or a budget reason | Answers part of Sol-22 and Fable-M6. |
| Measured defaults with the rationale next to the value | borrow the pattern | LAYUP parameter file: each default with its evidence line (I4) | A good I4 practice. |
| Full loop (board, database, `nh approve` squash-land) | reject | none | State outside Git (I1); `nh approve` squash-lands, which Armature forbids ("never squash"). |
| Overall | borrow the pattern | See above | Evidence: the runs above. |

## Where the searchers were wrong or incomplete

- Searcher B: "a second model in a fresh session reviews each diff and is told to refute done ... merging stays human". Correct, and the run confirmed it. Incomplete: `nh approve` merges the PR itself by a local squash under the operator's identity (README line 170); a human triggers it, but the tool does the merge.
- Searcher A: "keep: fresh independent review, test guard". The test guard ("tamper guard") does not cover Go test files; on a Go target it reports `clean` with 0 tests.
- Neither searcher noted that telemetry and session replay are on by default, or that the tool needs a Claude subscription token in `~/.no_human/.env` (it refuses to start with `ANTHROPIC_API_KEY` set in subscription mode).

## Limits of this evaluation

- Only the one-shot gate ran. The full loop (board, database, coder, tests, PR) and the GitHub Action were not run; statements about them come from the README and the code.
- We replaced the SDK-bundled CLI in our venv with a wrapper; no_human itself was unchanged. The reviewer model was `claude-sonnet-5`, not the default `claude-opus-4-8` on which the README's catch-rate was measured.
- No cost figure: `nh gate` does not print one, and the calls went through the SDK stream.
- One small toy diff per case; the review quality on large diffs was not tested.
