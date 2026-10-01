# Evaluation: Agent Context Store (ACS)

| Field | Value |
| ----- | ----- |
| Repository | https://github.com/max9159/agent-context-store |
| Pinned commit | ae90049b87d1a529a1e0c1e595879c53c03b104f (2026-07-06T02:28:07Z, "chore: release v0.3.1"; 107 commits) |
| License | MIT, from the LICENSE file |
| Language and needs | TypeScript, Node.js 20 or later, npm package `agent-context-store` (CLI `acs`, 0.3.1), dependency `ajv`. No server and no database. Optional local web server (`acs site`) and optional MkDocs. Store modes: `in-repo` (`.acs/`), `local` (user data directory), `dedicated` (another repository). 6 stars, last push 2026-07-06 (API GET). |
| Evaluated by | Claude Opus 5.5 on Claude Code, 2026-09-28 |
| Elapsed | About 35 minutes of wall clock for all three candidates in one session (12:31 to 13:06 UTC); I did not time each candidate apart |
| Agent runs and cost | 1 run: Claude Code `claude -p` with the `acs-dev` skill, claude-haiku-4-5-20251001, `--max-budget-usd 0.60`; cost 0.1295 USD, 74 s, 23 turns |

## Verdict

Reject as a tool and as a format; borrow two small patterns (the handoff rule table and the per-task audit log). ACS is a document relay with structural validation. Its approvals are fields that the producing agent can set, its rules and schemas are files in the store that any agent can edit, two writers on different tasks conflict in shared files, a fresh clone fails its own `validate`, and it has no telemetry.

## What it is (from the code)

Paths are at the pinned SHA; `C` = `src/packages/core/src/index.ts` (3,021 lines), `CLI` = `src/packages/cli/src/index.ts`.

- **Store layout** (in-repo mode): `.acs/` with `artifacts/<task>/<type>/*.md` (YAML front matter plus Markdown), `handoffs/<task>/HOFF-<task>-<FROM>-<TO>.yaml`, `packages/<task>/<role>.context.md`, `audit/<YYYY-MM-DD>.log`, `audit/tasks/<task>.jsonl`, `index.json`, and the rules: `roles/*.yaml`, `workflows/default-sdlc.yaml`, `artifact-types/*.yaml`, `schemas/*.json`, `docs/*.md`.
- **Handoff rule table**: `from`, `to`, `requiredArtifacts`, `requiredState`, for example `ba -> sa` needs `srs`, `user-story`, `acceptance-criteria` in state `approved`; `system -> <any role>` needs nothing ([C#L476-L485](https://github.com/max9159/agent-context-store/blob/ae90049b87d1a529a1e0c1e595879c53c03b104f/src/packages/core/src/index.ts#L476-L485)). `acs handoff check` applies it ([C#L1842-L1852](https://github.com/max9159/agent-context-store/blob/ae90049b87d1a529a1e0c1e595879c53c03b104f/src/packages/core/src/index.ts#L1842-L1852)).
- **Schemas**: [handoff.schema.json](https://github.com/max9159/agent-context-store/blob/ae90049b87d1a529a1e0c1e595879c53c03b104f/src/assets/schemas/handoff.schema.json) requires `id, task_id, from_role, to_role, handoff_type, status, approval_status, artifacts, context_summary, open_questions, readiness`; `status` and `approval_status` are free strings; `additionalProperties: true`. The project's own guide calls four schemas "currently permissive" ([docs/ACS_CONTEXT_STORE_INTRO_en.md#L132-L135](https://github.com/max9159/agent-context-store/blob/ae90049b87d1a529a1e0c1e595879c53c03b104f/docs/ACS_CONTEXT_STORE_INTRO_en.md#L132-L135)). The validator loads each schema from the store first, so the store's copy wins ([C#L2469-L2476](https://github.com/max9159/agent-context-store/blob/ae90049b87d1a529a1e0c1e595879c53c03b104f/src/packages/core/src/index.ts#L2469-L2476)).
- **Modes**: `strict` (default) and `relaxed`; in relaxed mode missing upstream inputs are warnings ([C#L17-L22](https://github.com/max9159/agent-context-store/blob/ae90049b87d1a529a1e0c1e595879c53c03b104f/src/packages/core/src/index.ts#L17-L22), [C#L872-L890](https://github.com/max9159/agent-context-store/blob/ae90049b87d1a529a1e0c1e595879c53c03b104f/src/packages/core/src/index.ts#L872-L890)). Any caller can pass `--mode relaxed`.
- **Approval**: `acs handoff approve ... --reviewer <NAME>` sets `approval_status: approved` ([C#L1258](https://github.com/max9159/agent-context-store/blob/ae90049b87d1a529a1e0c1e595879c53c03b104f/src/packages/core/src/index.ts#L1258)); the reviewer is a free string. Artifact approval is a front matter field.
- **Audit**: every command appends one JSON line to `audit/<day>.log` and to `audit/tasks/<task>.jsonl` with `action`, `role`, `task_id`, `ts`, `session_id` ([C#L2862-L2890](https://github.com/max9159/agent-context-store/blob/ae90049b87d1a529a1e0c1e595879c53c03b104f/src/packages/core/src/index.ts#L2862-L2890)). `role` is the role the caller named; `session_id` is per CLI call unless `ACS_SESSION_ID` is set.
- **Index**: `acs index` rewrites the whole `index.json` ([C#L1542](https://github.com/max9159/agent-context-store/blob/ae90049b87d1a529a1e0c1e595879c53c03b104f/src/packages/core/src/index.ts#L1542)).
- **Tokens**: only a "conservative" estimate of the size of a context package against a budget ([C#L1954](https://github.com/max9159/agent-context-store/blob/ae90049b87d1a529a1e0c1e595879c53c03b104f/src/packages/core/src/index.ts#L1954)). No usage, latency or duration.
- **Skill install**: `acs install-skills` writes to `os.homedir()/.claude/skills` (and `.cursor`, `.codex`) and appends a block to `~/.claude/CLAUDE.md` ([CLI#L90](https://github.com/max9159/agent-context-store/blob/ae90049b87d1a529a1e0c1e595879c53c03b104f/src/packages/cli/src/index.ts#L90), [CLI#L428](https://github.com/max9159/agent-context-store/blob/ae90049b87d1a529a1e0c1e595879c53c03b104f/src/packages/cli/src/index.ts#L428)).

## What we ran

`npm install --prefix eval-work/agent-context-store/npm agent-context-store@0.3.1` (matches the pinned release), `HOME=eval-work/agent-context-store/home`, `GH_TOKEN` and `GITHUB_TOKEN` unset. Toy target: the same Go module, with a local bare remote.

1. `acs init --mode in-repo` (no TTY, so no wizard): about 55 files under `.acs/` (roles, workflow, 16 artifact types, 9 schemas, 14 templates, 4 rule documents, `index.json`). No skill files were installed.
2. BA leg by hand: `acs ba new srs`, `acs ba new acceptance-criteria` created files with the template text "Describe the artifact content here." Then:
   ```
   acs validate --role ba --task TOY-0001    ->  mode strict / OK validation passed
   acs dev new srs ...                       ->  exit 1 (dev may not create srs)
   acs handoff create --from ba --to sa      ->  OK Created handoff HOFF-TOY-0001-BA-SA (mode=strict)
   acs handoff check --from ba --to sa       ->  ERROR ... status or approval_status must be approved
                                                  Missing required artifact for ba -> sa: user-story
   ```
   Validation passed on empty template content. `handoff create` succeeded before `handoff check` could pass.
3. Self-approval: as the same BA agent, I created the user story and changed `status: draft` to `approved` and `approval_status: pending` to `approved` in the front matter with `sed`. `acs handoff check --from ba --to sa` then gave `OK validation passed`.
4. SA leg by hand (sdd, adr, api-design, self-approved), `handoff create --from sa --to dev`, commit, push.
5. `acs install-skills --agent claude --path .` wrote `home/.claude/skills/acs-*` and `home/.claude/CLAUDE.md`, not the project. A second run appended the same block to `CLAUDE.md` again. Without the `HOME` redirect this would change the user's global Claude Code setup.
6. One agent run: I copied the skills into the toy's `.claude/skills` (excluded from Git) and ran `claude -p "/acs-dev Pick up TOY-0001 from SA. Implement it in Go ..."`. The agent added `Sub` and a table test, ran `make check`, created `implementation-note` and `unit-test-note`, ran `acs validate`, `acs handoff create --from dev --to qa` and `acs package --role qa`. Audit lines:
   ```
   {"action":"artifact.create","role":"dev","task_id":"TOY-0001","artifact":".../IMPL-TOY-0001.md","ts":"2026-09-28T12:58:31.037Z","session_id":"cli-mul96yfw-g2ygfq"}
   {"action":"handoff.create","from":"dev","to":"qa","task_id":"TOY-0001",...,"mode":"strict",...}
   ```
   `make check` was the agent's own choice; ACS did not run a check on the code. No token, cost or duration was recorded anywhere; only the Claude Code JSON output had them.
7. Two writers on different tasks: clones `c1` and `c2` of the bare remote each created one task (`TOY-0002`, `TOY-0003`) with `acs ba new srs` and `acs index`, committed and pushed. `c2` was rejected; both `git pull --rebase` and `git pull --no-rebase` gave:
   ```
   CONFLICT (content): Merge conflict in .acs/audit/2026-09-28.log
   CONFLICT (content): Merge conflict in .acs/index.json
   ```
   ACS ships no `.gitattributes` (for example `merge=union` for the log).
8. Rule change by a working agent: in `c1`, I added `dev` to `allowed_roles.create` in `.acs/artifact-types/srs.yaml`. `acs dev new srs` then succeeded (exit 0), where it had failed before (step 2).
9. Fresh clone: `git clone` of the bare remote, then `acs validate` gave `ERROR validation failed ... error Missing directory: summaries` (exit 1). Git does not keep the empty `summaries/` directory that `init` creates ([C#L346](https://github.com/max9159/agent-context-store/blob/ae90049b87d1a529a1e0c1e595879c53c03b104f/src/packages/core/src/index.ts#L346), [C#L977](https://github.com/max9159/agent-context-store/blob/ae90049b87d1a529a1e0c1e595879c53c03b104f/src/packages/core/src/index.ts#L977)).

Processes: I did not start `acs site`. No process stayed; `ps` at 12:59:49Z showed none.

## In-Scope items S1-S12

| Item | Mark | Evidence |
| ---- | ---- | -------- |
| S1 Problem Statement Quality | no | No gap check; the BA writes an SRS from a template. |
| S2 Reproducible Discipline Setup | partly | `acs init` writes the same default SDLC pack each time (step 1). Not Armature; no evidence per value. |
| S3 Rule Protection | no | Roles, workflow, artifact types and schemas are store files; step 8 shows an agent granting itself a right. |
| S4 Stack-Dependent Gates | no | No code gates. |
| S5 Role Handoffs | partly | A handoff YAML with required fields and a from/to rule table with required artifacts and states (steps 2, 3). But the schemas are permissive, the states are self-set (step 3), `--mode relaxed` weakens checks, and it is not based on Armature. |
| S6 Autonomous Clarification | no | `open_questions` and `blocking_questions` are string lists. No routing, no answer record. |
| S7 Verification on Every Change | no | `acs validate` checks store structure and metadata only (step 2 passed on template text). |
| S8 Human-on-the-Loop | no | A human starts each role by hand. Approval is a field or a `--reviewer` string that the agent can set (step 3). No escalation rule. |
| S9 Stall Resolution | no | Nothing. |
| S10 Cost Visibility | no | Only a size estimate of a context package. Step 6 recorded no tokens, latency or duration. |
| S11 Specification Synthesis | partly | SRS, user story and acceptance-criteria templates, and a `source_refs` field. No trace check to the text of a problem statement. |
| S12 Harness-Agent Neutrality | partly | Skills for Claude Code, Cursor, Codex and OpenClaw read one in-repo store. The skills go into the user's home directory (step 5), not the repository. |

## Invariants I1-I9

| Invariant | Effect | Reason |
| --------- | ------ | ------ |
| I1 Git is the system of record | supports | In `in-repo` mode, all records are in `.acs/` (my run). The `local` and `dedicated` modes put them outside the project repository ([CLI#L374](https://github.com/max9159/agent-context-store/blob/ae90049b87d1a529a1e0c1e595879c53c03b104f/src/packages/cli/src/index.ts#L374)); LAYUP would need to forbid them. |
| I2 Independent repository | neutral | The records are plain Markdown, YAML and JSON; a human continues without ACS. But a fresh clone fails `acs validate` (step 9). |
| I3 Agents cannot change rules or gates | conflicts | Rules are store files that the working agent can edit (step 8); the schema in the store overrides the packaged one. |
| I4 No value without evidence | neutral | Default roles and rules with no evidence; they are ACS's own process, not configuration of the target's build. |
| I5 Inactive check is not a pass | conflicts | `--mode relaxed` turns missing-input errors into warnings, so a check that is not applied still gives "OK". Strict `validate` passed with no content (step 2). |
| I6 Deterministic over LLM | supports | All checks are code. They are weak, but deterministic. |
| I7 Domain changes content, not rules | neutral | No such notion. |
| I8 Pinned baseline | neutral | `.acs/acs.yaml` has `version: 1`; the CLI version is not pinned. |
| I9 Replaceable harness | neutral | The store is neutral; the agent instructions are installed per harness into the user's home directory, which is the host-level rule source of Fable-M23. |

## Deep-check findings it answers

- **Author-8 / Fable-M14 / Sol-21 (handoff schema and transition table).** Partly, as a pattern: the rule table `from`, `to`, `requiredArtifacts`, `requiredState` is the "transition table" that the fixes ask for. It does not answer the core of Sol-21: validity rests on status fields that the agent writes (step 3).
- **Sol-14 / Fable-M15 (event records).** Weak pattern: one append-only JSONL file per task with action, role, time and session. The role is self-declared.
- **Fable-M12 / Sol-3 (many writers).** Not answered; step 7 shows the opposite: shared aggregate files (`index.json`, the daily log) make two writers on different tasks conflict. The per-task log alone would not conflict.
- **Fable-M23 (host user-level rule source).** Not answered; `install-skills` creates this defect (step 5).
- **Sol-15 (project-wide question index), Author-9 / Sol-13 (reversal), Sol-11 / Sol-12 (telemetry), Sol-4, Sol-5, Author-6.** None.

## Parts and their verdicts

| Part | Verdict | How LAYUP would take it | Reason |
| ---- | ------- | ----------------------- | ------ |
| Handoff rule table (`from`, `to`, required artifacts, required state) | borrow the pattern | A transition table in LAYUP's handoff check, with the required state computed by code (gate results, verifier record), not read from a field the author sets | Gives Fable-M14 its table; ACS's own version trusts self-set fields. |
| Append-only JSONL audit per task | borrow the pattern | One event file per task (or per session) in LAYUP's records | Append-only per-task files do not conflict between tasks. |
| Artifact front matter (`id`, `type`, `owner`, `version`, `source_refs`, `depends_on`) | borrow the pattern, small | Field names for requirement and design records, where `source_refs` becomes a checked trace to a PSB fact | Close to what S11 needs; ACS does not check the trace. |
| Daily shared log and rebuilt `index.json` | reject | Not used | Conflict between any two writers (step 7). |
| Status and approval as author-editable fields | reject | Not used | Self-approval passes the handoff check (step 3). |
| ACS as a program next to LAYUP | reject | Not used | Node tool, writes into the user's home directory, adds a second rule system in `.acs/` that agents can edit, fails on a fresh clone, no telemetry. |
| Its file format as LAYUP's records format | reject | Not used | No telemetry, question-answer, decision or reversal records; permissive schemas. |

How it deals with many writers to one Git branch: it does not. Each CLI call writes files with plain `writeFile` and `appendFile`, with no lock. Per-task files are safe across tasks; the daily log and `index.json` conflict every time two writers work on the same day (step 7).

## Where the searchers were wrong or incomplete

- Searcher A: "Git in-repo `.acs/`, a dedicated repo, or local-only mode; no server/database required." Correct. It did not see that `install-skills` writes to the user's home directory.
- Searcher A: S3 `partly` ("role profiles define permissions"). The permissions are editable files, and step 8 bypassed them in one line: `no`.
- Searcher A: S7 `partly` ("`acs validate --role ...`"). `validate` never looks at code; step 2 passed on template text: `no`.
- Searcher A: S8 `partly` ("handoff confirmation by next role"). The confirmation is a status field that the sender can set (step 3): `no`.
- Searcher A: "some handoff and approval schemas are currently permissive". Confirmed in the code, and worse: the store's own copy of a schema overrides the packaged one, so an agent can weaken it.
- Searcher A: I5 `P` (partial support). Step 2 and the `relaxed` mode show a conflict.
- Searcher A did not test many writers; step 7 shows conflicts between different tasks.

## Limits of this evaluation

- One agent run (DEV leg). The BA, SA and QA legs were run by hand with the CLI.
- I did not run the `local` or `dedicated` modes, `acs site`, or the MkDocs engine.
- I did not read all 3,021 lines of the core; I read the parts behind each command I ran and each claim above.
