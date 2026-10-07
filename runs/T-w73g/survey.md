# A survey of five public repositories, as an input of the later specification tasks

Task `T-w73g` ([#133](https://github.com/pharzam/layup/issues/133)), 2026-10-07.
Author: Claude Opus 5.5 on Claude Code.

## What this survey is, and what it is not

The Operator named five repositories and asked for ready components and patterns
for LAYUP. This file lists what was found, by the milestone whose specification
task reads it. It makes no decision: each specification task takes a pattern or
rejects it, with its own reasons.

This survey is **not** the public-solution search of
[Solution selection](../../docs/engineering-discipline.md#solution-selection).
It reads five repositories that the Operator chose; it did not search from the
terms of the problem statement, read no list in full, and had one model, not two
blind searchers. Each specification task still runs its own search
([`guardrails.md` §2](../../docs/guardrails.md), "A public-solution search that
confirms the design").

**Method.** Each repository was cloned with `--depth 1` into a scratch directory
outside the checkout, at the commit in the table below. One read-only agent
(Claude Opus 5.5, a fresh Explore subagent) read each repository, with the brief
in [`brief.md`](brief.md). No code of a repository was run. The five reports are
here word for word: [`report-agnostic-ai.md`](report-agnostic-ai.md),
[`report-ruflo.md`](report-ruflo.md),
[`report-cc-multi-cli-plugin.md`](report-cc-multi-cli-plugin.md),
[`report-ruvnet-brain.md`](report-ruvnet-brain.md),
[`report-self-improvement-loop.md`](report-self-improvement-loop.md).

**Limits.** A pattern below is what the code says at one commit, read by one
model; nothing was run, so no claim of behaviour is tested. The author checked by
hand only these claims: each license, where the agents' runtime state lives (the
table below), the Rust files of `ruflo`, the Jev numbers (M3a), and the `stdin`
note of the `codex exec` launcher (M2b). [`cites.sh`](cites.sh) checks that each
source path below exists at its commit; it does not check what the file says.

**Citation form.** A source path is written `` `<repo>:<path>` ``, where `<repo>`
is the first column of the table below. A directory ends with `/`.

## The repositories

None gives a component that LAYUP can use as it is. Each is mainly Node or
TypeScript (`ruflo` also has 39 Rust files under `ruflo:v3/crates/`), LAYUP is Go
([ADR-0010](../../docs/adr/0010-use-go-as-the-technology-stack.md)), and each
keeps the runtime state of its agents (sessions, memory, locks, routing counts,
lessons) outside Git, where LAYUP keeps every record in Git (Invariant 1). So
each row below is a pattern to write again in Go, never a dependency.

| Repo | Repository | Commit | License | Runtime state of the agents |
| ---- | ---------- | ------ | ------- | --------------------------- |
| `agnostic-ai` | [ucsandman/Agnostic-AI](https://github.com/ucsandman/Agnostic-AI) | `aab09ff0d9533967045e7f2b04e478ec4e05528e` | MIT | `storage/*.jsonl` and `~/.claude`, ignored by Git (`agnostic-ai:.gitignore`) |
| `ruflo` | [ruvnet/ruflo](https://github.com/ruvnet/ruflo) | `9fa701a466159242dfec06379ce1278c849a3e34` | MIT | SQLite and vector memory under `.claude-flow/` and `.ruv-swarm/`, ignored by Git (`ruflo:.gitignore`, `ruflo:v3/@claude-flow/memory/`) |
| `cc-multi-cli-plugin` | [greenpolo/cc-multi-cli-plugin](https://github.com/greenpolo/cc-multi-cli-plugin) | `3dcb057c10b16d3ed911e8ce509e91398712c906` | Apache-2.0 (one directory copied from openai/codex) | the user's `~/.claude` and `~/.codex`; receipts only to a file that an environment variable names (`cc-multi-cli-plugin:plugins/multi-core/src/launcher.ts`) |
| `ruvnet-brain` | [stuinfla/ruvnet-brain](https://github.com/stuinfla/ruvnet-brain) | `b4590d469f52937550d8d89afa86a394a1d72875` | MIT | `~/.config/ruvnet-brain` (`ruvnet-brain:plugin/scripts/continuation-gate.mjs`); the committed data files under `ruvnet-brain:data/` are the exception |
| `self-improvement-loop` | [kolezka/self-improvement-loop](https://github.com/kolezka/self-improvement-loop) | `ba245adbb5be8075fa01044dd82aac636155870f` | **PolyForm Noncommercial 1.0.0**, with a separate paid commercial license | XDG data and state directories (`self-improvement-loop:packages/core/src/layout.ts`) |

**The license of `self-improvement-loop`.** Its license allows no commercial use
of the code, and its `COMMERCIAL-LICENSE.md` counts internal tooling of a
for-profit company as commercial use. A port of its code to Go is a derived work.
Take its ideas only, and write the code fresh. The other four allow copying with
the license notice.

## Patterns by milestone

### `M2b` — sessions on registered harnesses, and the telemetry row

Concerns: [ADR-0020](../../docs/adr/0020-route-role-sessions-over-registered-harnesses.md),
[ADR-0024](../../docs/adr/0024-record-every-sessions-cost-and-stop-at-the-milestone-cap.md);
the demo of `M2b` is a probe session on each registered harness and one complete
telemetry row.

| Pattern | Source | What fits | What does not fit |
| ------- | ------ | --------- | ----------------- |
| A process runner: one JSON event per line with a size cap on each line, the end of stderr kept with tokens removed, the process group stopped with SIGINT, then SIGTERM, then SIGKILL; a run counts as done only when the harness sends its final result event | `cc-multi-cli-plugin:plugins/multi-core/src/gateway/harness-process.ts` | The done rule and the stop sequence | It has no wall-clock deadline and no idle timeout; a hung harness runs until someone stops it. LAYUP needs both |
| Failure classes: "fails the same way again" (no binary, policy, busy) against "uncertain" (stream ended early, crash) | `cc-multi-cli-plugin:plugins/multi-core/src/gateway/harness-failure.ts` | A retry only for the uncertain class | — |
| The environment cleaned: the API keys of other providers removed, and one key that would switch the billing | `cc-multi-cli-plugin:plugins/multi-core/src/gateway/harness-process.ts`, `ruvnet-brain:tri-smart-skill/tri-smart/scripts/review.mjs` | The credential of a session comes only from its register row (K40, open for `M2b`) | Neither gives a session its own HOME; LAYUP's evaluation found tools that write into the home directory (`guardrails.md` §2) |
| The tool policy checked after launch: the tool list that the harness announces is compared with the policy, and a forbidden tool stops the run with `policy` | `cc-multi-cli-plugin:plugins/multi-grok/src/cli.ts`, `cc-multi-cli-plugin:plugins/multi-grok/src/permissions.ts` | A harness that ignores an unknown tool name fails closed | — |
| The session ID chosen before launch and recorded as `interrupted` at the first event; a resume after a crash is then deterministic | `cc-multi-cli-plugin:plugins/multi-grok/src/harness.ts` | Restart of a session | — |
| The launch of `claude -p` and `codex exec`, one worktree for each writer, a timeout and an output cap; `codex exec` waits for the end of stdin, so stdin is closed at launch (checked by hand) | `ruflo:v3/@claude-flow/codex/src/dual-mode/orchestrator.ts` | The flags and the stdin rule for two harnesses | — |
| The argument lists of `claude`, `codex` and `grok` in a read-only mode with JSON output, with process-group timeouts | `ruvnet-brain:tri-smart-skill/tri-smart/scripts/review.mjs` | A third source for the same flags | Its acceptance is a model's text matched by a pattern |
| The headless flags of `grok` and `agy` (Antigravity): the table in the report | [`report-cc-multi-cli-plugin.md`](report-cc-multi-cli-plugin.md) | The flags, the events and the resume of two more harnesses | `agy` runs with its own permissions switched off and depends on a global hook |
| A token ledger from the Claude and Codex transcripts: Claude duplicates removed by message and request ID; Codex gives running totals, so the ledger takes their differences and reads a drop as a reset; cache reads apart from cache writes | `ruflo:plugins/ruflo-cost-tracker/scripts/_ledger.mjs` | The telemetry row of a session from the harness's own transcript | — |
| A price table of one row per model, each with the URL of its source | `ruflo:plugins/ruflo-cost-tracker/data/prices.json` | The form of `prices.tsv` (Invariant 4: a value with its evidence) | — |
| Receipts of one line per run with no prompt text, written atomically | `cc-multi-cli-plugin:plugins/multi-core/src/gateway/receipts.ts`, `cc-multi-cli-plugin:plugins/multi-core/src/gateway/atomic-write.ts` | The content of a telemetry row | LAYUP already has one writer of the records in Git; its lock file is not needed |

### `M2f` — rule protection

Concern: [ADR-0017](../../docs/adr/0017-prevent-rule-changes-by-agents.md); the
demo of `M2f` is that each known-bad patch fails its kind and that no agent
writes to a rule path.

| Pattern | Source | What fits | What does not fit |
| ------- | ------ | --------- | ----------------- |
| A hash lock over the gate files, with an entry for one key of a JSON file | `agnostic-ai:tools/gates/freeze.cjs`, `agnostic-ai:tools/gates/gate-manifest.json`, `agnostic-ai:engine/hooks/gate-freeze.cjs` | A committed lock that the audit compares | The lock file there is ignored by Git and local to one machine |
| A gate counts only after a proof: break the file that it checks, expect red, restore, expect green | `agnostic-ai:tools/prove/prove.cjs` | The known-bad patch of each kind | — |
| A fix kept as a file and a substring that must stay; counts of a bad pattern that may only go down | `ruflo:verification/witness-fixes.json`, `ruflo:verification/mcp-tool-baseline.json` | Cheap checks on files in Git | The signatures on them |
| Deterministic checks of a command: destructive commands, secrets, a tool allowlist, a diff-size limit | `ruflo:v3/@claude-flow/guidance/src/gates.ts` | — | — |
| A small predicate language for a gate (all, any, not, a pattern with a check against slow patterns), replayed on recorded inputs in its own process with a timeout | `self-improvement-loop:packages/nudges/src/gates.ts`, `self-improvement-loop:packages/nudges/src/gate-runner.ts` | Idea only (license) | — |

### `M2g` — a build task merges, and the idea owner accepts

Concerns: `REQ-007` and [`docs/spec/gate.md`](../../docs/spec/gate.md); the
demo of `M2g` is a build task that merges with every gate green, and the idea
owner's "accept" as a row of `acceptance.tsv`.

| Pattern | Source | What fits | What does not fit |
| ------- | ------ | --------- | ----------------- |
| Four verdicts by exit code: 0 pass, 1 fail, 2 `insufficient_spec` (a human decides), 3 malformed; a check that cannot run fails; the "needs a human" tag is written in the specification, never asked of a model | `agnostic-ai:tools/task-contract/task-contract.mjs`, `agnostic-ai:tools/task-contract/parse.mjs` | The verdict of an acceptance check | — |
| An accept bound to a digest: the human accepts the SHA-256 of the exact diff shown, and the accept is refused if the branch, the commit or the base moved | `self-improvement-loop:packages/review/src/snapshot.ts`, `self-improvement-loop:packages/review/src/index.ts` | The "accept" row names what was accepted. Idea only (license) | — |
| Many policies, one refusal point, with an order of precedence and one combined reason | `ruvnet-brain:plugin/scripts/decision-gate.mjs` | — | — |
| The fix of review findings: findings in groups that share no file, the fixer and the reviewer on different models, one fix, then one fresh read of the whole diff | `agnostic-ai:workflows/fix-findings.js` | The review of a build task | The finders are models; this is a review, not a gate |

### `M3a` — the smart-if provider

Concern: [ADR-0021](../../docs/adr/0021-branch-at-named-points-through-a-smart-if-provider.md).

| Pattern | Source | What fits | What does not fit |
| ------- | ------ | --------- | ----------------- |
| A measured result of Jev: over 407 options the top-1 recall was 5 % and top-3 8 %; over the 22 options that were ever correct it was 18 % and 42 % (checked by hand) | `agnostic-ai:labs/claude-mods/experiments/jev/FINDINGS.md` | The option set of a decision point: its size is a parameter with this evidence | One experiment on one machine, on a different question |
| Five named yes/no questions to a judge (contradicts, vague, unsupported, unsafe, unrelated) with a probability threshold | `self-improvement-loop:packages/curriculum/src/prompts.ts`, `self-improvement-loop:packages/providers/src/index.ts` | A decision point, never an engine check. Idea only (license) | — |

### `M3d` — stalls and the budget cap

Concerns: [ADR-0023](../../docs/adr/0023-stop-a-stall-at-a-limit-and-diagnose-it-with-a-fresh-context.md),
[ADR-0024](../../docs/adr/0024-record-every-sessions-cost-and-stop-at-the-milestone-cap.md).

| Pattern | Source | What fits | What does not fit |
| ------- | ------ | --------- | ----------------- |
| Numeric stop checks: a step limit, the ratio of repeated steps, the slope of the cost over the last 10 steps, forced checkpoints; the outcomes continue, checkpoint, pause (a human looks) or stop | `ruflo:v3/@claude-flow/guidance/src/continue-gate.ts` | The limit of the stall procedure | Two of its inputs come from a model; drop them |
| The same tool call, hashed, counted while it repeats; the limits 3, 5 and 8 are logged so that they can be tuned | `agnostic-ai:engine/hooks/repeat-tool-guard.cjs` | Raw data for "a step that repeats without progress" | Advisory only |
| A stop after the third machine wake-up with no human turn | `agnostic-ai:engine/hooks/wakeup-guard.cjs` | — | — |
| Reserve the budget before the spend, not after: the race of check, spend, record | `ruflo:v3/docs/adr/ADR-164.1-budget-tracker-atomicity.md` | The milestone cap | Its own code only warns (`ruflo:plugins/ruflo-cost-tracker/scripts/budget.mjs`) |
| Escalate when rewording does not help: after a set number of revisions, if the rate since the last revision is still near the rate since promotion, a human decides | `self-improvement-loop:packages/curriculum/src/plan.ts` | A limit with no model. Idea only (license) | — |

### The learning loop — waits for the Operator's decision

Concern: [ADR-0025](../../docs/adr/0025-learn-routing-from-the-records-at-each-retrospective.md).
The plan does not make the learning loop a milestone; these rows are for the
Operator's decision and for the task that follows it.

| Pattern | Source | What fits | What does not fit |
| ------- | ------ | --------- | ----------------- |
| Route qualification: a candidate model passes exact-answer checks first; then a reviewer of another provider grades the old and the new output blind, in random order; one role changes per run; the policy file is written only if its hash did not change | `ruvnet-brain:scripts/model-weekly-qualification.mjs`, `ruvnet-brain:scripts/model-routing-policy-promotion.mjs`, `ruvnet-brain:config/model-router/qualification-contract.json` | A routing change at a retrospective | The reviewer is a model; it sits behind the exact-answer checks |
| Model classes that a task's facts can raise but never lower; no route gives no model, never a weaker one in silence; eval cases with a minimum and a maximum class | `ruvnet-brain:config/model-router/policy.default.mjs`, `ruvnet-brain:config/model-router/routing-eval-cases.json` | — | Its classifier of free text is keyword matching |
| Success and failure counts per role, size class and model, as Beta(α, β) | `ruflo:v3/@claude-flow/cli/src/ruvector/model-router.ts` | Counts from the TSV records, with a deterministic choice | Its random sampling and its keyword score |
| A watermark for a lesson: a pattern needs a set number of new cases past the count at its last promotion or rejection, so a rejected lesson cannot come back at once | `self-improvement-loop:packages/curriculum/src/plan.ts` | Deterministic over records. Idea only (license) | — |
| Rates before promotion, since promotion and since revision, with a minimum number of sessions | `self-improvement-loop:packages/feedback/src/rates.ts` | Whether a lesson worked, per milestone. Idea only (license) | — |
| A lesson promoted only when two or more projects teach it independently | `ruvnet-brain:plugin/scripts/lesson-promote.mjs`, `ruvnet-brain:plugin/scripts/lesson-gate.mjs` | A count, not a model judgement | — |
| A frozen eval set: the cases pinned by hash, the Wilson lower bound as a ratchet, a missing baseline is a failure, and a baseline is written only on request | `ruvnet-brain:scripts/eval-brain.mjs`, `ruvnet-brain:evals/held-out.json`, `ruvnet-brain:evals/baseline.json` | The measure of a routing change | — |
| A model claim counts only with a word-for-word quote of its source of a set minimum length | `self-improvement-loop:packages/curriculum/src/router.ts` | A deterministic check on a model's lesson. Idea only (license) | — |
| A declared estimate against the actual count, logged, and the constants fitted again from the data | `agnostic-ai:engine/hooks/subagent-budget-guard.cjs`, `agnostic-ai:tools/subagent-budget/calibrate.cjs` | — | — |

## What not to take

| What | Source | Why |
| ---- | ------ | --- |
| Agent memory in SQLite and vectors | `ruflo:v3/@claude-flow/memory/` | State outside Git (Invariant 1) |
| "Consensus" between agents | `ruflo:v3/@claude-flow/swarm/src/consensus/` | In one process; one node approves its own proposal; it is not an independent review |
| Hooks that pass when they fail | `ruflo:plugin/hooks/hooks.json` | Most end in `\|\| true`, and one approves the plugin's own tools; the opposite of ADR-0017 |
| The path to OpenAI through a private backend | `cc-multi-cli-plugin:plugins/multi-openai/src/auth.ts` | It reads the Codex login and calls a private endpoint |
| A global hook written into the user's home | `cc-multi-cli-plugin:plugins/multi-antigravity/src/hooks.ts` | Outside the target; the policy depends on it |
| A spend summary from transcripts with prices in the code | `agnostic-ai:tools/spend/spend.cjs` | Prices with no source |

## The claims of the READMEs that the code does not support

These are the agents' findings, not checked by hand; each report gives its
evidence.

- `agnostic-ai`: the stop gate blocks only when an environment variable sets it
  (`agnostic-ai:engine/hooks/forced-verify-stop-gate.cjs`); the pre-commit hook
  that the README describes is not in the repository.
- `ruflo`: the "89 % routing accuracy" comes from a hand-written table
  (`ruflo:.claude/helpers/router.cjs`); the budget warns and never stops; the
  README counts 314 tools and its own baseline file lists 285.
- `cc-multi-cli-plugin`: no timeout on a harness run; the README says each run
  isolates its identity, but the harnesses use the real home directory.
- `ruvnet-brain`: the baseline numbers in the README differ from
  `ruvnet-brain:evals/baseline.json`; the threshold of its proof script is -3, so
  a negative relevance passes (`ruvnet-brain:scripts/prove.mjs`); its review is
  not blind.
- `self-improvement-loop`: not every model call runs at temperature 0; a lesson
  can come back in up to five sessions; an option merges without a human (off by
  default).

## What this survey does not decide

- No pattern is chosen. The specification task of each milestone above takes or
  rejects each row, and runs its own public-solution search.
- The order of the plan does not change, and the learning loop stays the
  Operator's decision.
- No code from these repositories enters LAYUP by this task.
