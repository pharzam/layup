# Evaluation: GNAP (Git-Native Agent Protocol)

| Field | Value |
| ----- | ----- |
| Repository | https://github.com/farol-team/gnap |
| Pinned commit | 2571f5574a0e28001638030222ad7f8bb65b57f7 (2026-03-17T21:12:21Z, 24 commits in total) |
| License | MIT, from the LICENSE file ("Copyright (c) 2026 Farol Labs") |
| Language and needs | No code. The protocol is prose and JSON examples in README.md. It needs Git and a shared remote. No server, no database, no daemon. 86 stars, not archived, last push 2026-03-17 (API GET, 2026-09-28). |
| Evaluated by | Claude Opus 5.5 on Claude Code, 2026-09-28 |
| Elapsed | About 35 minutes of wall clock for all three candidates in one session (12:31 to 13:06 UTC); I did not time each candidate apart |
| Agent runs and cost | 1 run: Claude Code `claude -p`, model claude-haiku-4-5-20251001, `--max-budget-usd 0.50`; cost 0.040556 USD, 22 s wall clock |

## Verdict

Borrow the pattern. The idea "one JSON file per record, in the project repository, written by any harness" fits Invariants 1, 2 and 9. But the protocol does not prevent a double claim, it loses the loser's claim with no trace, its run record has no latency and only optional fields, and a commit author proves nothing. LAYUP must not depend on it as a format.

## What it is (from the code)

GNAP has no program. The whole protocol is the README at the pinned SHA. The repository ships no JSON Schema file, no CLI and no test.

- Four entities in `.gnap/`: `agents.json`, `tasks/<id>.json`, `runs/<task>-<attempt>.json`, `messages/<n>.json`, and a `version` file with `4`. [README.md#L95-L117](https://github.com/farol-team/gnap/blob/2571f5574a0e28001638030222ad7f8bb65b57f7/README.md#L95-L117)
- Agent fields: `id`, `name`, `role`, `type` (`ai` or `human`), `status`. [README.md#L152-L173](https://github.com/farol-team/gnap/blob/2571f5574a0e28001638030222ad7f8bb65b57f7/README.md#L152-L173)
- Task fields and the state machine `backlog, ready, in_progress, review, done, blocked, cancelled`; `blocked` and `blocked_reason` are optional fields. [README.md#L195-L245](https://github.com/farol-team/gnap/blob/2571f5574a0e28001638030222ad7f8bb65b57f7/README.md#L195-L245)
- Run fields: required `id`, `task`, `agent`, `state`, `started_at`; optional `finished_at`, `tokens {input, output}`, `cost_usd`, `result`, `error`, `commits`, `artifacts`. The one sentence "cost tracking (budget = sum of runs)" is at [README.md#L296-L298](https://github.com/farol-team/gnap/blob/2571f5574a0e28001638030222ad7f8bb65b57f7/README.md#L273-L298).
- Message fields: `id`, `from`, `to`, `at`, `text`; optional `type` (`directive`, `status`, `request`, `info`, `alert`), `channel`, `thread`, `read_by`. [README.md#L319-L336](https://github.com/farol-team/gnap/blob/2571f5574a0e28001638030222ad7f8bb65b57f7/README.md#L319-L336)
- Transport: commit subject `<agent-id>: <action>`; "Eventual consistency"; "Conflicts: Standard git merge. If conflict, pull + rebase + retry push". [README.md#L342-L359](https://github.com/farol-team/gnap/blob/2571f5574a0e28001638030222ad7f8bb65b57f7/README.md#L342-L359)
- The agent loop: "Set state -> in_progress -> commit + push" and "If `git push` fails: `git pull --rebase`, re-check, retry (max 3)". [ONBOARDING.md#L117-L131](https://github.com/farol-team/gnap/blob/2571f5574a0e28001638030222ad7f8bb65b57f7/ONBOARDING.md#L117-L131)
- Out of scope by design: "Operational concerns (budget enforcement, stall detection, concurrency control) belong in AgentHQ, not the protocol." [CLAUDE.md#L47](https://github.com/farol-team/gnap/blob/2571f5574a0e28001638030222ad7f8bb65b57f7/CLAUDE.md#L47), and [CONTRIBUTING.md#L23](https://github.com/farol-team/gnap/blob/2571f5574a0e28001638030222ad7f8bb65b57f7/CONTRIBUTING.md#L23).

History that the README does not tell: versions 2 and 3 had a `gnap.sh` that wrote through the GitHub contents API with the blob `sha` (optimistic locking: "On 409 Conflict: re-read, re-apply your change, retry"), a `budget.json` with `limit_usd` and `spent_usd`, and spending-limit policies. Commits `e01b31d` ("Remove gnap.sh"), `2429434` and `121dda1` removed them. Version 4 keeps only `git pull --rebase`.

## What we ran

Work directory: `eval-work/gnap/`. `HOME` was `eval-work/gnap/home` for all Git commands. No push to any network remote.

1. Toy target: a Go module `toy/` (`go.mod`, package `calc` with `Add`, one test, a `Makefile` with `go vet ./...` and `go test ./...`), a local bare remote `remote.git`. `make check` passed.
2. Schemas: GNAP has none, so I wrote four JSON Schemas (draft 2020-12) from the README tables (`eval-work/gnap/schemas/`) and a validator (`validate.py`, `jsonschema` 4.26.0 in a venv). It also checks the rule "id matches filename". GNAP's own examples pass:
   ```
   ok   agents.json / tasks/FA-1.json / runs/FA-1-1.json / messages/1.json   version: 4
   ```
3. By hand, as the Operator `op`: `.gnap/version`, `agents.json` (op human; dev1, dev2, qa as ai), `tasks/FA-1.json` (assigned to dev1 and dev2, reviewer qa, state `ready`). Commit `op: create FA-1 add calc.Sub`, push to the bare remote. Two clones `a1` (dev1) and `a2` (dev2) of the same bare remote.
4. Double claim, scenario A (the ONBOARDING loop exactly: only `state` changes). Both clones read `ready`, both commit `in_progress`, dev1 pushes first:
   ```
   dev2 push exit: 1                     (non-fast-forward)
   hint: use --reapply-cherry-picks to include skipped commits
   Successfully rebased and updated refs/heads/main.
   Everything up-to-date                 dev2 retry push exit: 0
   eee66ae dev1 dev1: checkout FA-1      (dev2's commit is gone)
   {"state":"in_progress","assigned_to":["dev1","dev2"]}
   ```
   Result: Git dropped dev2's claim as an identical patch. Both of dev2's Git commands then succeeded. The task file has no "claimed by" field, so dev2 cannot see from the file that the claim is not its own. `git log` shows only dev1. The race left no trace.
5. Scenario B (the claim also writes `updated_at`): the rebase stops with `CONFLICT (content): Merge conflict in .gnap/tasks/FA-2.json`, on the `updated_at` line only. The `state` lines merged with no conflict. The protocol gives no rule for this conflict.
6. Scenario C (both agents write run attempt 1 and the next message number):
   ```
   CONFLICT (add/add): Merge conflict in .gnap/messages/1.json
   CONFLICT (add/add): Merge conflict in .gnap/runs/FA-3-1.json
   ```
   Sequential message IDs and `<task>-<attempt>` run IDs collide when two writers work at the same time.
7. One real run (`run-agent.sh`): the wrapper wrote `runs/FA-1-1.json` with `state: running`, pushed, ran Claude Code on the task, and then filled the run record from the harness JSON output. The agent itself cannot know its token count; the wrapper must write it.
   ```
   "started_at": "2026-09-28T12:36:06Z", "finished_at": "2026-09-28T12:36:28Z",
   "tokens": {"input": 138348, "output": 1422}, "cost_usd": 0.040556,
   "duration_ms": 16561, "duration_api_ms": 15464, "commits": ["b21dd46..."]
   ```
   `duration_ms` and `duration_api_ms` are my additions. The spec neither allows nor forbids extra fields. The agent added `Sub` and a table test; `make check` passed. The wrapper set FA-1 to `review` and wrote message 2 (`type: request`, to `qa`). `qa` (by hand, in `a2`) set FA-1 to `done` with a comment, and set FA-2 to `blocked` with `blocked_reason`. All files passed my validator.
8. Forgery: in clone `a2`, the AI agent dev2 set FA-3 from `ready` straight to `done`, paused the reviewer `qa` in `agents.json`, and committed as `op: done FA-3` with author `op`. The push succeeded. The per-file validator passed. Only my history walk (`transitions.py`, every commit, every task state change against the README state machine) found it:
   ```
   BAD  8c0f48d .gnap/tasks/FA-3.json: ready -> done by op
   ```
   The check names `op`, the forged author.

What a human sees in `git log` at the end: one line per action, in the GNAP convention (`dev1: start run FA-1-1`, `dev1: finish run FA-1-1`, `dev1: move FA-1 -> review`, `qa: done FA-1 - make check passes`, `op: done FA-3`). The log is readable. It does not show the lost claim of dev2, and it cannot show that `op: done FA-3` came from dev2.

Processes: no server or daemon. The one `claude -p` process exited. Nothing to stop.

## In-Scope items S1-S12

| Item | Mark | Evidence |
| ---- | ---- | -------- |
| S1 Problem Statement Quality | no | No entity or step for gaps in a problem statement. |
| S2 Reproducible Discipline Setup | no | The protocol sets up `.gnap/` only. No Armature, no evidence per value. |
| S3 Rule Protection | no | ONBOARDING step 2 gives each agent Git write access. In step 8, dev2 paused the reviewer and marked a task done as `op`. Nothing refused it. |
| S4 Stack-Dependent Gates | no | No gates. |
| S5 Role Handoffs | partly | Tasks, runs and messages have required fields, and I validated them with schemas that I wrote from the README. There is no handoff entity, no artifact schema per role, and no transition check. Transition validity needed my own history walk. Not based on Armature. |
| S6 Autonomous Clarification | no | A message with `type: request` and `thread` can carry a question. Nothing classifies it, routes it to an owner, or records an accepted answer. |
| S7 Verification on Every Change | no | No check runs. `review` is only a state. |
| S8 Human-on-the-Loop | no | Humans are `type: human` participants. There are no planned decision points and no escalation rule. A human can edit any file at any time. |
| S9 Stall Resolution | no | `blocked`, `blocked_reason` and many runs per task exist. There is no retry limit ("max 3" is for push retries only), no stall detection and no diagnosis. CLAUDE.md#L47 puts stall detection outside the protocol. |
| S10 Cost Visibility | partly | A run can hold `tokens`, `cost_usd`, `started_at`, `finished_at` in the repository (step 7). All of them except `started_at` are optional. There is no latency field. The input and output split has no cache tokens. |
| S11 Specification Synthesis | no | Not in the protocol. |
| S12 Harness-Agent Neutrality | partly | Any process that can `git push` takes part; I used a shell wrapper and Claude Code. Rules and gates are not part of it, and nothing makes a second harness verify the first. |

## Invariants I1-I9

| Invariant | Effect | Reason |
| --------- | ------ | ------ |
| I1 Git is the system of record | supports | All state is in `.gnap/` in the project repository. No other store exists. |
| I2 Independent repository | supports | Plain JSON files. A human with Git and an editor continues the work. The target needs no GNAP tool. |
| I3 Agents cannot change rules or gates | neutral | GNAP has no rules or gates of its own, and it gives every agent full write access. Step 8 shows that an agent can change the team and the task states. Protection must come from outside the protocol. |
| I4 No value without evidence | neutral | The only default is `heartbeat_sec` 300, with no evidence. It is not written into the target unless the author copies it. |
| I5 Inactive check is not a pass | neutral | No checks. |
| I6 Deterministic over LLM | neutral | No judgement of any kind. A deterministic transition check is possible (step 8) but not supplied. |
| I7 Domain changes content, not rules | neutral | No rules. |
| I8 Pinned baseline | neutral | `.gnap/version` pins the protocol version, a useful pattern. No Armature pin. |
| I9 Replaceable harness | supports | The format belongs to no harness. |

## Deep-check findings it answers

- **Fable-M12 / Sol-3 (one writer has six writers).** GNAP is the opposite model: every producer writes and pushes its own records, and no single writer exists. As a pattern this answers the finding only when each record has a unique path. GNAP's own IDs (`messages/<n>.json`, `runs/<task>-<attempt>.json`) collide (step 6). LAYUP would need collision-free names (for example a time-ordered unique ID per record) and append-only files.
- **Author-6 / Fable-M10 (records before Intake).** Weak pattern: `.gnap/` is committed first, before any task. It says nothing about where Intake decisions live.
- **Sol-15 (Early Question Share has no project-wide index).** Pattern: `messages/` is one project-wide directory with `at`, `from`, `to` and `thread`. It lacks a "human" flag and a delivery-start boundary.
- **Sol-11 / Sol-12 (telemetry).** Pattern only: one run file per attempt, so a retry loop is visible run by run, and "budget = sum of runs". It does not answer the findings: missing tokens are an absent field, not a recorded "not reported", and latency does not exist.
- **Author-8 / Fable-M14 / Sol-21 (handoff schema).** Not answered. GNAP gives required fields and a state table, not an artifact schema or a validator.
- **Sol-5 (commit author does not prove the pusher).** Not answered. GNAP rests on the commit author convention, and step 8 shows the forgery that Sol-5 describes.
- **Sol-4 (dead-man job), Author-9 / Sol-13 (reversal record), Sol-14 / Fable-M15 (measures with no record), Fable-M4 / Author-1 (disagreement stall).** None.

## Parts and their verdicts

| Part | Verdict | How LAYUP would take it | Reason |
| ---- | ------- | ----------------------- | ------ |
| One JSON file per record under one directory, with a `version` file | borrow the pattern | LAYUP's own records directory, written with Go `encoding/json` (standard library) | Fits I1, I2, I9 and F-0004 fact 1. The idea costs nothing. |
| Run record per attempt (`tokens`, `cost_usd`, `started_at`, `finished_at`, `commits`) | borrow the pattern, extended | LAYUP's telemetry row per session: make token, latency and wall-clock fields required, add cache tokens, model, harness, and an explicit `not reported` value | GNAP's fields are optional and have no latency (PSB S10, section 7.1 Telemetry Completeness). |
| Message with `thread` and `type: request` | borrow the pattern | A question record and an answer record linked by ID | Needs owner, kind, human flag, acceptance and reversal fields that GNAP lacks. |
| Claim by "set `in_progress`, push, rebase, retry" | reject | Not used | It does not prevent a double claim. Step 4 lost a claim with no error and no trace. |
| Commit subject `<agent-id>: <action>` as identity | reject as proof; keep as a readable log style | Readable commit subjects only | Step 8: the author field is self-declared (Sol-5). |
| The protocol as a dependency or file format | reject | Not used | Draft v4, no schema file, no code, no activity since 2026-03-17, and the authors reject new entities (CLAUDE.md#L45), but LAYUP needs handoff, decision, question, answer, reversal and stall records. |

How GNAP deals with many writers to one Git branch: by optimistic push races only. The loser of a push race rebases. An identical edit disappears silently (step 4). A different edit to the same file gives a text conflict with no resolution rule (step 5). Two new files with the same name give an add/add conflict (step 6). The earlier SHA-based compare-and-swap through the GitHub API was removed in version 4. For LAYUP, a safe form is: each writer only adds new files with unique names, never edits a shared file, and a claim counts only if the claim commit is on the remote branch after the push (checked by reading the branch, not by the push exit code).

Could its file format be LAYUP's records format? Not as it is. It has no record for a handoff, a decision, an answer, a reversal or a stall, and its telemetry record lacks latency and required fields. The directory-of-small-JSON-files shape is a good base for LAYUP's own format.

## Where the searchers were wrong or incomplete

- Searcher B: "A run has `tokens` and `cost_usd`." True, but both are optional (README#L285-L290). The searcher does not say so; searcher A does.
- Searcher B: "cost tracking (budget = sum of runs)." That is one README sentence. No budget entity, limit or check exists at the pinned SHA. The budget file and spending limits existed in version 3 and were removed (commit `121dda1`).
- Searcher B: `blocked` and `blocked_reason` exist. True, as optional task fields. Nothing sets or reads them; they are free text.
- Searcher B marks S12 `covers`. Too high: the protocol carries no rules and no cross-harness verification, so the PSB item is only partly met.
- Searcher B marks S5 and S8 `partly`; searcher A marks S6, S8 and S9 `partly`. From the files: S5 partly is fair; S6, S8 and S9 have a field or a participant type, but no mechanism.
- Both searchers take "Git history IS the audit log" at face value. The run shows that the log omits lost claims and cannot prove who made a commit.
- GNAP's own article says "When two agents race to claim the same task, git's non-fast-forward rejection handles it. The loser retries." ([docs/article.md#L51](https://github.com/farol-team/gnap/blob/2571f5574a0e28001638030222ad7f8bb65b57f7/docs/article.md#L51)). Step 4 disproves this for the claim the protocol describes: the loser's retry succeeds and both agents hold the task.
- Summary: "a Git-only task, run and cost protocol". Correct, with the limits above; the "cost" part is two optional fields.

## Limits of this evaluation

- GNAP has no implementation, so "what we ran" is my own implementation of the README. Another reading of the prose can differ (for example, an agent that writes a "claimed by" comment would turn scenario A into a conflict like scenario B).
- The schemas are mine. They allow extra fields because the spec says nothing about them.
- I did not test the heartbeat loop over time, or a remote on a forge with branch protection.
- One agent run only. The race tests used shell scripts in place of agents, which is enough to show the Git behaviour.
