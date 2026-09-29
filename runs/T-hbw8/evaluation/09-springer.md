# Evaluation: Springer

| Field | Value |
| ----- | ----- |
| Repository | https://github.com/z4gunn/springer |
| Pinned commit | 902cd7298b604e8d1aaf00d6975c8bd4ae87bd8b (2026-09-07) |
| License | MIT (LICENSE file: "MIT License, Copyright (c) 2026 z4gunn") |
| Language and needs | Markdown agents and skills for Claude Code (27 agents, 197 skills), JSON Schema 2020-12 files, Python 3 scripts (`jsonschema` in a project venv). No server, no database, no daemon. Needs Claude Code. Optional: `gh` for the pr-merge gate, Node for Linear sync, Graphviz and PlantUML for diagrams. |
| Evaluated by | Claude Opus 5.5 on Claude Code, 2026-09-28 |
| Elapsed | about 62 minutes |
| Agent runs and cost | 4 runs of `claude -p` with `claude-sonnet-5` (the subagents used their pinned opus or sonnet models): 1.93 + 2.88 + 2.36 + 4.81 = 11.98 USD, 48 minutes of agent time |

**Model IDs (plan-review note 8, added 2026-09-29).** Springer's subagents ran on the aliases `opus` and `sonnet` that its own files pin. The evaluator did not record the full model ID behind each alias, and the run logs are not kept, so this file does not name one (Invariant 4). No verdict of this file rests on the quality of that model: each rests on what the candidate's code did.

## Verdict

Borrow the pattern. Springer is the best role-and-handoff design in this group: typed JSON artifacts that a script validates, a gap list batched at the first gate, an append-only cycle log, and bounded retries. But it is Claude Code only, it has no rule protection, its trace is free text, and it keeps telemetry outside Git. LAYUP cannot run it as a component.

## What it is (from the code)

- **Project copy.** `scripts/new-project.sh` copies the runtime (`.claude/skills`, `.claude/agents`, `.claude/references`, `.claude/hooks`, `schemas/`) into a new Git repository. It installs a project `CLAUDE.md` and makes `runs/` tracked ([new-project.sh:63-71](https://github.com/z4gunn/springer/blob/902cd7298b604e8d1aaf00d6975c8bd4ae87bd8b/scripts/new-project.sh#L63-L71), [:81](https://github.com/z4gunn/springer/blob/902cd7298b604e8d1aaf00d6975c8bd4ae87bd8b/scripts/new-project.sh#L81)). So the whole runtime becomes part of the target.
- **The harness.** `spgr-run-harness` is a skill (prose that the main Claude session follows). It runs Plan, Do, Check, Act per tick. It claims a lock, derives a readiness snapshot, asks the orchestrator agent for a batch, dispatches subagents, validates, appends a `pdca-cycle` record and pauses at a gate ([SKILL.md:40-65](https://github.com/z4gunn/springer/blob/902cd7298b604e8d1aaf00d6975c8bd4ae87bd8b/.claude/skills/spgr-run-harness/SKILL.md#L40-L65), [:100-116](https://github.com/z4gunn/springer/blob/902cd7298b604e8d1aaf00d6975c8bd4ae87bd8b/.claude/skills/spgr-run-harness/SKILL.md#L100-L116), [:140-163](https://github.com/z4gunn/springer/blob/902cd7298b604e8d1aaf00d6975c8bd4ae87bd8b/.claude/skills/spgr-run-harness/SKILL.md#L140-L163)). The loop itself is not code. The model follows the text.
- **Deterministic scripts.** `derive-ready-queue.py` computes open gates, open escalations and `blocked` from the artifact files ([derive-ready-queue.py:162-189](https://github.com/z4gunn/springer/blob/902cd7298b604e8d1aaf00d6975c8bd4ae87bd8b/.claude/skills/spgr-run-harness/scripts/derive-ready-queue.py#L162-L189)). `rebuild-projection.py` rebuilds `run-state.json` from the cycle log. `claim-run.py` holds a single-writer lock with a 30-minute heartbeat ([claim-run.py:27](https://github.com/z4gunn/springer/blob/902cd7298b604e8d1aaf00d6975c8bd4ae87bd8b/.claude/skills/spgr-run-harness/scripts/claim-run.py#L27)). `pin-learnings.py` hashes prior retrospectives.
- **Schemas.** 26 artifact types share an envelope with `confidence_map` (confirmed, proposed, needs-human-input), `decision_log`, `version_type` and `checksum` ([_envelope-v1.json:7-18](https://github.com/z4gunn/springer/blob/902cd7298b604e8d1aaf00d6975c8bd4ae87bd8b/schemas/_envelope-v1.json#L7-L18), [:60-71](https://github.com/z4gunn/springer/blob/902cd7298b604e8d1aaf00d6975c8bd4ae87bd8b/schemas/_envelope-v1.json#L60-L71)). `schemas/validate.py` runs a real `Draft202012Validator` ([validate.py:82-95](https://github.com/z4gunn/springer/blob/902cd7298b604e8d1aaf00d6975c8bd4ae87bd8b/schemas/validate.py#L82-L95)).
- **Gates.** The `hil-checkpoint` type has seven checkpoint kinds: architecture options, architecture confirmation, PRD approval, design direction, pr-merge, security flag, scope change ([hil-checkpoint-v1.json:23-33](https://github.com/z4gunn/springer/blob/902cd7298b604e8d1aaf00d6975c8bd4ae87bd8b/schemas/hil-checkpoint-v1.json#L23-L33)).
- **Escalations.** The `escalation` type has seven kinds and a `routing_target` (human, orchestrator or an agent) ([escalation-v1.json:28-49](https://github.com/z4gunn/springer/blob/902cd7298b604e8d1aaf00d6975c8bd4ae87bd8b/schemas/escalation-v1.json#L28-L49)). The orchestrator routes an ambiguity to the upstream agent and sends to the human only what changes a human constraint or confirmed scope ([spgr-agent-orchestrator.md:50-55](https://github.com/z4gunn/springer/blob/902cd7298b604e8d1aaf00d6975c8bd4ae87bd8b/.claude/agents/spgr-agent-orchestrator.md#L50-L55), [pdca-harness.md:407-422](https://github.com/z4gunn/springer/blob/902cd7298b604e8d1aaf00d6975c8bd4ae87bd8b/.claude/references/pdca-harness.md#L407-L422)).
- **Telemetry hook.** A PreToolUse, PostToolUse and SubagentStop hook on the Agent tool appends dispatch and completion events with token and duration counters to `runs/<run-id>/events.jsonl` ([log-agent-events.py:28-35](https://github.com/z4gunn/springer/blob/902cd7298b604e8d1aaf00d6975c8bd4ae87bd8b/.claude/hooks/log-agent-events.py#L28-L35), [:116-130](https://github.com/z4gunn/springer/blob/902cd7298b604e8d1aaf00d6975c8bd4ae87bd8b/.claude/hooks/log-agent-events.py#L116-L130)).
- **Retrospective.** At completion the harness writes a `run-retrospective`. Learnings are advisory, pinned by hash at the next run start, and a learning that changes a rule needs `requires_human_promotion` ([SKILL.md:169-173](https://github.com/z4gunn/springer/blob/902cd7298b604e8d1aaf00d6975c8bd4ae87bd8b/.claude/skills/spgr-run-harness/SKILL.md#L169-L173), [pdca-harness.md:424-435](https://github.com/z4gunn/springer/blob/902cd7298b604e8d1aaf00d6975c8bd4ae87bd8b/.claude/references/pdca-harness.md#L424-L435)).
- **Run profiles.** `brochure`, `small`, `saas`, `mobile` scale the phase set, the story cap and the PR unit ([pdca-harness.md:39-51](https://github.com/z4gunn/springer/blob/902cd7298b604e8d1aaf00d6975c8bd4ae87bd8b/.claude/references/pdca-harness.md#L39-L51)).

## What we ran

1. `git clone https://github.com/z4gunn/springer src` at the pinned SHA.
2. `sh ../_shared/make-toy.sh <springer dir>` made the Go toy (`go.mod`, `greet/`, one test, `Makefile` with `go vet` and `go test`) and a local bare remote.
3. `./src/scripts/new-project.sh app small` made the project copy with one commit. We copied the toy Go files and the 15-line PSB (`docs/inputs/PSB.md`) into it, and added a bare remote `app-remote.git`. We installed `jsonschema` into `app/.venv` with pip. `python3 schemas/validate.py --self-check` printed `SELF-CHECK OK: 27 schema files, 26 artifact types, gating verified.`
4. Safety: `GH_TOKEN` and `GITHUB_TOKEN` unset, a `gh` shim on `PATH` that exits 1, `GH_CONFIG_DIR` set to an empty directory, `SPGR_DASHBOARD=0`, no Linear config.
5. Four runs of `claude -p ... --model claude-sonnet-5 --max-budget-usd 5..8 --setting-sources project,local --permission-mode bypassPermissions --output-format json` in `app/`:
   - **Run 1** (start, 9 min, 1.93 USD): paused at `prd-approval` (`CKPT-toy-1-001`). The PM agent wrote a PRD, an NFR, a backlog, 3 user stories and 3 acceptance-criteria sets. All 13 files passed `validate.py`.
   - **Run 2** (our gate answer, 15 min, 2.88 USD): stamped the checkpoint resumed, set the requirements to v1.0 `approved` with a checksum, and the Architect wrote three options with scores. Paused at `architecture-options-selection`.
   - **Run 3** (6.5 min, 2.36 USD): the Architect wrote one note and four ADRs. The harness applied four "fold-ins" to the approved PRD and NFR. Paused at `architecture-confirmation`.
   - **Run 4** (19 min, 4.81 USD): red acceptance tests first, then three stories on `feature/toy-1-todocount`, a Security and a Performance audit, Code Review pass 1 REQUEST_CHANGES, one fresh fix agent, pass 2 APPROVE. It pushed the branch to the bare remote and paused at `pr-merge` (`CKPT-toy-1-004`). We checked: `make check` passes; `go run ./cmd/todocount` prints `example.com/toy/greet	1`.
6. Negative schema tests on copies of the run's own artifacts:
   - A criterion with no `then`: `INVALID ... 'then' is a required property` (rc 1).
   - A PRD whose `evidence_refs` names `docs/inputs/NOPE.md#L99`: `VALID` (rc 0).
   - A success metric with an empty `measurement`: `VALID` (rc 0).
   - An artifact with an unknown `artifact_type`: `VALID ... (envelope-only, no content schema registered for this type yet)` (rc 0).

**The PSB gaps.** Both gaps came back in one batch at the first gate. Line 9 ("stale") became OQ-1, OQ-2 and OQ-3 with options and a recommendation. Line 8 ("fast") became OQ-4 ("How big is a large repository?") with a proposed value (10,000 files, 5 s). The PM also listed ten "proposed rulings" (PR-01 to PR-10) for the human to confirm, and blocked the stale story. It did not guess: the decision log says "Defining it without the human fills a gap with an assumption." Later gates asked new questions (HI-1 to HI-5, NHI-1, NHI-2), so questions did not all come before delivery.

**Faults we saw.**
- The orchestrator named a non-existent agent `spgr-agent-pm`. The main session corrected it by hand.
- The run changed the approved PRD from v1.0 to v1.1 (fold-ins) and kept the old checksum `sha256:58e374dd...`. No code checks the checksum. `grep` for `checksum` in `*.py`, `*.ts`, `*.mjs` finds nothing. The "tamper detection" is prose only ([spgr-version-artifact SKILL.md:32-38](https://github.com/z4gunn/springer/blob/902cd7298b604e8d1aaf00d6975c8bd4ae87bd8b/.claude/skills/spgr-version-artifact/SKILL.md#L32-L38)).
- The harness ran the Architect in the background once, so the hook logged its completion at dispatch time. The harness said so itself.
- Nothing was committed to Git until run 4. For three runs the whole run store was untracked files in the work tree. The final `chore(run)` commit left `events.jsonl` (the telemetry) uncommitted.
- The commits used the host's global Git identity. Springer gives agents no identity of their own.

**Telemetry.** `runs/toy-1/events.jsonl` had 16 completion events with metrics, 890,242 tokens in total, with `duration_ms` per subagent. It does not record cost, and it does not record the main session's own tokens.

## In-Scope items S1-S12

| Item | Mark | Evidence |
| ---- | ---- | -------- |
| S1 Problem Statement Quality | partly | Run 1 found both planted gaps and asked them in one batch at `prd-approval`, with options. But the check is a model's judgement, not a gap check before delivery. The PM rule is "No vague NFRs" and a content-sources table ([spgr-agent-product-manager.md:37](https://github.com/z4gunn/springer/blob/902cd7298b604e8d1aaf00d6975c8bd4ae87bd8b/.claude/agents/spgr-agent-product-manager.md#L37), [:42](https://github.com/z4gunn/springer/blob/902cd7298b604e8d1aaf00d6975c8bd4ae87bd8b/.claude/agents/spgr-agent-product-manager.md#L42)). Later gates asked seven more questions. |
| S2 Reproducible Discipline Setup | partly | `new-project.sh` makes a repeatable copy at a Git SHA. It has no pinned version record, no Armature, and no evidence per value. The script stamps `saas` as the default profile whatever the argument, because its `sed` pattern does not match the template line ([new-project.sh:75-76](https://github.com/z4gunn/springer/blob/902cd7298b604e8d1aaf00d6975c8bd4ae87bd8b/scripts/new-project.sh#L75-L76); our `app/CLAUDE.md` says `saas` after `small`). |
| S3 Rule Protection | no | The only guard is a permission deny list for five destructive Git forms ([project-settings.json:27-33](https://github.com/z4gunn/springer/blob/902cd7298b604e8d1aaf00d6975c8bd4ae87bd8b/templates/project-settings.json#L27-L33)). The hooks only log. Agents have `Write` and `Edit` on `CLAUDE.md`, `.claude/` and `schemas/`. The approved-artifact checksum is not checked by any code. |
| S4 Stack-Dependent Gates | no | No layout, boundary or contract gate exists. The Code Reviewer (a model) reviews against ADRs. The only stack rule file is `typescript-standards.md`. For Go, the run used the toy's own `make check`. |
| S5 Role Handoffs | partly | Every handoff is a JSON artifact, and `validate.py` runs a real JSON Schema check (run 1: 13 of 13 valid; a missing `then` fails). But an unknown type passes as "envelope-only", content objects allow extra fields, and a trace to a missing file passes. The check is run by the model following the skill, not forced by code. |
| S6 Autonomous Clarification | partly | Escalations route to the upstream agent or the Architect, and reach the human only when a human constraint or confirmed scope changes ([pdca-harness.md:417-422](https://github.com/z4gunn/springer/blob/902cd7298b604e8d1aaf00d6975c8bd4ae87bd8b/.claude/references/pdca-harness.md#L417-L422)). In our run no agent-to-agent question happened; the PM put all open items to the human. No record of an "accepted answer" exists. |
| S7 Verification on Every Change | partly | Check runs `validate.py`, the project check and CI; a Code Reviewer and vertical audits run on code diffs; the review is bounded to two passes ([pdca-harness.md:98-121](https://github.com/z4gunn/springer/blob/902cd7298b604e8d1aaf00d6975c8bd4ae87bd8b/.claude/references/pdca-harness.md#L98-L121)). Run 4 did this. No deterministic layout, boundary or test-pyramid gate exists, and a change only to `runs/` or `docs/` gets no review. |
| S8 Human-on-the-Loop | partly | Seven fixed gate kinds in the schema; the harness pauses only there ([SKILL.md:188-189](https://github.com/z4gunn/springer/blob/902cd7298b604e8d1aaf00d6975c8bd4ae87bd8b/.claude/skills/spgr-run-harness/SKILL.md#L188-L189)). A human merges every PR on `saas`. Gates are phase approvals, not the acceptance of each requirement. No escalation rule for business-forking decisions exists in code. We answered four gates in four runs. |
| S9 Stall Resolution | partly | Two failed retries on one unit, then escalate to the human; one review plus one re-review, then open findings go to the human ([SKILL.md:136-139](https://github.com/z4gunn/springer/blob/902cd7298b604e8d1aaf00d6975c8bd4ae87bd8b/.claude/skills/spgr-run-harness/SKILL.md#L136-L139)). A fix goes to a fresh agent. No independent diagnosis, no stall record type, and the lock is only for a dead session. |
| S10 Cost Visibility | partly | The hook writes tokens and duration per subagent to `events.jsonl` (run 4: 890,242 tokens over 16 events). No cost, no main-session tokens, no per-task record, and the harness left the file uncommitted. The retrospective metrics have no token field ([run-retrospective-v1.json:52-61](https://github.com/z4gunn/springer/blob/902cd7298b604e8d1aaf00d6975c8bd4ae87bd8b/schemas/run-retrospective-v1.json#L52-L61)). |
| S11 Specification Synthesis | partly | IDs exist and are pattern-checked: `STORY-YYYY-n` ([user-story-v1.json:25-28](https://github.com/z4gunn/springer/blob/902cd7298b604e8d1aaf00d6975c8bd4ae87bd8b/schemas/user-story-v1.json#L25-L28)), `AC-<story>-n` with given, when, then ([acceptance-criteria-v1.json:24-31](https://github.com/z4gunn/springer/blob/902cd7298b604e8d1aaf00d6975c8bd4ae87bd8b/schemas/acceptance-criteria-v1.json#L24-L31)). Our run gave each criterion a `check` command. The trace to the PSB is free text: `evidence_refs` is an array of strings ([prd-v1.json:27](https://github.com/z4gunn/springer/blob/902cd7298b604e8d1aaf00d6975c8bd4ae87bd8b/schemas/prd-v1.json#L27)); the run wrote `docs/inputs/PSB.md#L3-L12`, and a line that does not exist passes. A story traces to a PRD goal (`prd_goal_ref`), not to a sentence. |
| S12 Harness-Agent Neutrality | no | Everything is in Claude Code form: `.claude/agents`, `.claude/skills`, `.claude/settings.json` hooks, the `Agent` tool, the `model:` field. No file for another harness exists. The rules are one copy, but in one product's format (Problem 7). |

## Invariants I1-I9

| Invariant | Effect | Reason |
| --------- | ------ | ------ |
| I1 Git is the system of record | supports (partly) | All run state is files in the project tree: `runs/<id>/artifacts`, `checkpoints`, `pdca-cycle` log, `run-brief.json`. `.gitignore` tracks `runs/`. But the harness commits them only at the pr-merge gate; for three runs they were untracked, and `events.jsonl` stayed uncommitted. The transcripts stay in `~/.claude/projects`. |
| I2 The project repository is independent | conflicts | The target cannot work without Springer: `new-project.sh` copies the whole runtime (27 agents, 197 skills, hooks, schemas) into the target. The produced Go code builds with `make check` alone, but the run store needs `schemas/validate.py` and Claude Code to go on. This also conflicts with O-76. |
| I3 Agents cannot change the rules | conflicts | The agents that do the work can write `CLAUDE.md`, `.claude/` and `schemas/`. Run 3 changed an approved PRD without a new human approval (fold-ins), and the checksum stayed stale with no code to see it. |
| I4 No value without evidence | supports (partly) | The PM refused to invent "stale" and marked items `needs-human-input`; the content-sources table cites a file per fact. But the Architect's timing values came from one host, and the WIP limit (2) and lock time (30 min) have no evidence. |
| I5 An inactive check is not a pass | supports (partly) | The skill says a claim no script checks is "unverified" ([SKILL.md:86-89](https://github.com/z4gunn/springer/blob/902cd7298b604e8d1aaf00d6975c8bd4ae87bd8b/.claude/skills/spgr-run-harness/SKILL.md#L86-L89)) and forbids a pause with a failing required check. But `validate.py` gives `VALID` for a type with no schema, and in our run with no CI the harness wrote "does not apply" and paused. |
| I6 Deterministic over LLM | supports | Check runs `validate.py` first and opens a model only on failure; readiness, projection and lock are scripts ([pdca-harness.md:480-490](https://github.com/z4gunn/springer/blob/902cd7298b604e8d1aaf00d6975c8bd4ae87bd8b/.claude/references/pdca-harness.md#L480-L490)). The loop and the routing are still model steps. |
| I7 Domain changes content, never rules | neutral | Profiles change phases and caps, not rules. Nothing stops a project from editing the copied rules. |
| I8 Armature at a pinned version | neutral | No Armature. The copy records no Springer version either. |
| I9 A harness is replaceable | conflicts | Only Claude Code can run the agents, skills and hooks. |

## Deep-check findings it answers

- **Sol-17** (fact numbering needs a semantic step): partly. A model role (the PM) extracts needs from prose, and the human confirms at `prd-approval`. This is the fix Sol-17 names. It does not number clauses or check byte-exact quotes.
- **Sol-18** (a trace check does not find omitted needs): none. No completeness check exists. The human reads the PRD at the gate.
- **Sol-19** (code sets `Must` without an intent decision): answers the pattern. Priority is set by the PM and approved by the human at `prd-approval`. The stale story stayed `needs-refinement` until the human answered.
- **Sol-20** (no check that each requirement has a specification): none in code. The phase gate needs every prior-phase artifact confirmed, but no script maps a story to a design section.
- **Sol-32** (the ambiguity-owner map has no evidence): partly. The escalation kinds and the default routing are fixed text in the orchestrator; the human cannot set them per project, and no evidence is given for them.
- **Fable-M16** (no actor for the role of each step, the task class or the prompt): answers it. Each phase has a named agent with a pinned model and tool set; the dispatch contract fixes what each unit prompt contains ([pdca-harness.md:273-285](https://github.com/z4gunn/springer/blob/902cd7298b604e8d1aaf00d6975c8bd4ae87bd8b/.claude/references/pdca-harness.md#L273-L285)); a dispatch-tier table sends mechanical units to haiku ([pdca-harness.md:257-271](https://github.com/z4gunn/springer/blob/902cd7298b604e8d1aaf00d6975c8bd4ae87bd8b/.claude/references/pdca-harness.md#L257-L271)); a per-story brief limits context.
- **Fable-M20** (a deterministic check claimed for meaning gaps): answers the pattern. Springer gives meaning gaps to a model role and states it as such. Our run confirmed that the model found both planted gaps.
- **Fable-M21** (the spec check knows In-Scope facts only from a fixed PSB form): answers the pattern. The PM reads free prose, so no fixed heading is needed. The cost is that no code can check completeness.
- **Author-5** (fact numbering overclaimed): same answer as Sol-17: a role session reads the prose; no code claims to.
- **Author-18** (code sets a priority): same answer as Sol-19.
- Also relevant: **Author-8 / Sol-21 / Fable-M14** (handoff schema not defined): Springer shows a working form, a JSON Schema per artifact type plus a validator; **Author-1** (disagreement stall): the review bound of one pass plus one re-review, then the human, is a concrete rule.

## Parts and their verdicts

| Part | Verdict | How LAYUP would take it | Reason |
| ---- | ------- | ----------------------- | ------ |
| Whole system as a component | reject | none | Claude Code only (I9, S12); the runtime is copied into the target (I2, O-76); no rule protection (I3). |
| Artifact envelope (confidence map, decision log, version type) and per-type JSON Schemas | borrow the pattern | Define LAYUP's own handoff records with the same fields, validated by Go standard library code | Worked in all four runs. JSON Schema needs a library; LAYUP is standard library only, so write a small Go validator or use TSV with column checks. Close the holes we found: unknown type must fail, trace targets must resolve. |
| `needs-human-input` batch at the first gate, with options and a recommendation | borrow the pattern | The Intake role writes one question record per gap with options; one batch comment | It found both planted gaps and did not guess. |
| Append-only cycle log plus a projection rebuilt by a script | borrow the pattern | LAYUP's records branch: append rows, derive state by code | Rehydration from files worked across four separate sessions. |
| Bounded retry and bounded review rules | borrow the pattern | Parameters (O-82) with Springer's values as a data point | Run 4 used the bound correctly (REQUEST_CHANGES, one fix, APPROVE). |
| Agent-tool hook for tokens and duration | borrow the pattern | LAYUP reads the harness's own usage output per session and writes it to Git | Works for subagents; misses cost and the main session; harness-specific. |
| Advisory learnings pinned by hash, `requires_human_promotion` | borrow the pattern | Vision 2.3 retrospective with a human gate for rule changes | Not seen in our run (the run did not complete). The design fits I3. |
| Approved-artifact checksum | reject | none | No code checks it; our run shows it goes stale. |

## Where the searchers were wrong or incomplete

- Searcher A, S5 `covers`: too strong. The validator is real, but it passes unknown types, allows extra content fields and does not check that a trace target exists. The model runs it, not code that the agent cannot skip. `partly`.
- Searcher A, S12 `no`: correct. Searcher A said "0 stars; last push 2026-09-07": the last commit at the pinned SHA is 2026-09-07; stars not checked.
- Searcher A, S10 `no`: incomplete. A hook writes per-subagent tokens and duration to the repository tree. It is not complete, so `partly`.
- Searcher A, S1 "one-batch idea-owner questions not stated": the run did batch the gap questions at the first gate. `partly`, not a guess.
- Searcher A, "Architecture immutable once approved is methodology, not proven file protection": confirmed by the run. The run also changed an approved PRD.
- Searcher A, I1 `P`, I2 `P`: I1 is fair. I2 is a conflict, because the Springer runtime is copied into every target.
- Searcher A named five human gates from the README. The schema has seven checkpoint kinds, and our small-profile run hit four gates before the first merge.
- Searcher B listed Springer only as "partial-no-discipline". That misses the typed handoffs and the bounded-retry rules, which are the best in this group.

## Limits of this evaluation

- We ran four of the five allowed agent runs and stopped at the first `pr-merge` gate. We did not see the merge, a completed run, or a `run-retrospective`.
- `gh` was blocked on purpose, so the "publish before pausing" and CI checks did not run. No CI runner existed.
- One PSB and one Go target. Springer aims at SaaS and mobile, but it took the Go CLI well on the `small` profile.
- Agents ran with `--permission-mode bypassPermissions`, so the permission deny list was not tested.
- Model output varies. Another run can find fewer gaps or batch them differently.
- Processes: all four `claude -p` runs exited on their own. No dashboard, server or tmux session was started. We stopped our own file-watch monitors. `pgrep` found no remaining process.
