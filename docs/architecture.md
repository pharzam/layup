# The software architecture of LAYUP

The technical specification of LAYUP (`F-0003#51`): the components, their
interfaces, the state data model, the runner, the role model, the escalation
rule, the stall procedure and the telemetry record. It rests on
[`PRD-0001`](prd/PRD-0001-layup.md) and on the decision records that each
section names (ADR-0002, ADR-0007, ADR-0009, and ADR-0010 to ADR-0017); each
section names the requirements it designs and the record it rests on. Written by task `T-7qvc` (#66) after the panel and the Operator decisions
O-52 to O-62 recorded there and under [`runs/T-7qvc/`](../runs/T-7qvc/selection.md).

## In plain terms

> LAYUP is one command-line program, `layup`, that works on plain files in a Git
> repository, plus one small service, the LAYUP App, that reports the program's
> gate verdicts to the forge so that a pull request cannot merge past a failed or
> missing gate. Everything the program knows about a project is in that project's
> repository as tables and Markdown records that a person can read without a tool.

## 1. Components

| Component | What it is | Requirements | Rests on |
| --------- | ---------- | ------------ | -------- |
| **The engine**, `layup` | One Go binary, standard library only, Git called as `git`; no daemon, no database, no network service, no model call | NFR-005, NFR-007 | ADR-0010, ADR-0011 |
| **The runner**, the LAYUP App | A GitHub App installed on each target; a small receiver (or a workflow in LAYUP's repository) that runs the engine at a pinned commit on each pull-request event and posts required check runs under the App's identity | REQ-004, REQ-007, NFR-004 | ADR-0013 |
| **The rule guard** | The check `layup/rule-guard` of the App: red on a rule-path change until the Operator's approval | REQ-003 | ADR-0014 |
| **The target** | A project repository that LAYUP set up: the adapted Armature kit, the facts, the product the role agents deliver, and the records under `runs/` and `docs/`; never LAYUP's code | REQ-002, NFR-002 | ADR-0011 (O-10, O-11, O-13) |
| **The role agents** | Sessions of harness agents in the kit's roles, orchestrated by the engine's records, never by a LAYUP model call | REQ-005, REQ-013 | ADR-0015, ADR-0012 |

The engine and the runner are the only LAYUP software. The runner holds no gate
logic: it calls the engine and relays. A target passes its own kit gates without
either (NFR-002); it needs the runner only for the stack gates that the kit does
not carry, and its Operator can remove the runner's checks from the protection
body without LAYUP (ADR-0013).

## 2. The engine's commands and interfaces

Every repository command takes a repository root as an explicit argument, reads
and writes files under it, prints one line per result in the form `layup:
<command> OK|FAIL <cause>: <detail>`, and exits `0` (pass), `1` (fail) or `2`
(usage error, or a file it cannot read). The two commands that exist today are
the exceptions of shape: `layup version` prints one line, and `layup psb check
FILE` takes a file and prints the gap batch as TSV (`internal/cli/cli.go`). No
command reads global state as project state (NFR-001).

| Command | Phase | Reads | Writes | Requirement |
| ------- | ----- | ----- | ------ | ----------- |
| `layup version` | 1 (done) | — | stdout | — |
| `layup psb check FILE` | 1 (done) | a problem statement | the gap batch as TSV; exit 1 with gaps | REQ-001 |
| `layup setup TARGET` | 1 | `docs/setup/steps.tsv`, the pinned kit (`armature.pin`), the facts | the target: the kit copy at the pinned commit, the facts, the setup record, `branch-protection.json`, `stack.tsv`; stops at each `human_decision = yes` row and writes the answers before the next step | REQ-002, NFR-006 |
| `layup setup verify TARGET` | 1 | the target, the live protection | a report; `not-active` for a check that does not run | REQ-002, NFR-003, NFR-004 |
| `layup gate TARGET` | 1 | the target's recorded stack, its tree | one verdict per gate kind (`layout`, `interfaces`, `contracts`, `tests`): `pass`, `fail`, `not-active`; a row in `runs/<task>/gates.tsv` | REQ-004, REQ-007 |
| `layup handoff check TARGET` | 2 | `runs/<task>/handoffs.tsv`, the records it names | a report | REQ-005 |
| `layup escalation check TARGET CHANGE` | 3 | the intent set, the decision records, the telemetry | `runs/<task>/escalations.tsv` | REQ-008 |
| `layup stall open|examine|close TARGET` | 1 (record), 3 (procedure) | `runs/<task>/`, `handoffs.tsv` | `runs/<task>/stalls.tsv` (`open` the row; `examine` the examiner's harness, model and diagnosis path; `close` the outcome); the examiner, not the engine, writes `stall-<id>.md` | REQ-009, REQ-010 |
| `layup telemetry record|check TARGET` | 1 | the values a harness or a person gives | `runs/<task>/telemetry.tsv`; the count of incomplete records | REQ-011 |
| `layup audit rules TARGET` | 2 | the default branch's history | `runs/audit/rule-changes.tsv` | REQ-003 |

The engine cannot measure tokens: a harness or a person supplies them, and a
value a harness does not give is recorded as `not reported` and counts as
incomplete (REQ-011, ADR-0007).

## 3. The state data model

All state is plain files in the target (ADR-0011 decision 2): tab-separated
tables with a fixed header, rows only added, beside Markdown records in the
kit's shape. A deterministic check validates each table's header and columns
(`layup setup verify`, and each command for its own table).

| File | Kind | Header (or shape) | Written by | Record |
| ---- | ---- | ----------------- | ---------- | ------ |
| `docs/facts/F-*.md`, `docs/facts/*.md` | the facts, immutable, hashed in `docs/setup/facts.sha256` | the kit's facts record | `layup setup`, then the kit's rule | ADR-0011 |
| `docs/prd/PRD-*.md` | the requirements | the kit's PRD template | the role agents (REQ-012) | ADR-0002 |
| `docs/setup/steps.tsv`, `docs/setup/record-*.md`, `docs/setup/branch-protection.json`, `docs/setup/armature.pin` | the setup record and the mirrored forge setting | the kit's setup records | `layup setup` | ADR-0009, ADR-0011, ADR-0013 |
| `docs/setup/stack.tsv` | the target's recorded stack and gate set | `gate, kind, command, active` | `layup setup`, from the answer to step S01 | ADR-0013 |
| `docs/setup/budget.md` | the idea owner's limits: tokens and money per requirement and per task, committed dates | a Markdown record | the idea owner, before the pilot | ADR-0016 |
| `docs/roles.tsv` | the target's specialists (content, not a rule) | `role, specialism, harness_allowed` | the target's Operator | ADR-0015 |
| `runs/<task>/handoffs.tsv` | the handoff events | `ts, task, from_role, to_role, kind, artifact_path, artifact_sha, harness, model, verdict` | each role's session, validated by the engine | ADR-0015 |
| `runs/<task>/gates.tsv` | the gate verdicts | `ts, pr, head_sha, gate, verdict, layup_commit` | the runner | ADR-0013 |
| `runs/<task>/escalations.tsv` | the escalations | `id, trigger, decision_path, stopped_at_sha, answer_path, confirmed_business_forking, answered_at` | the engine; the idea owner's answer | ADR-0016 |
| `runs/<task>/stalls.tsv`, `stall-<id>.md` | the stall state and the diagnosis | `id, kind, limit, limit_value, opened_sha, examiner_harness, examiner_model, session, diagnosis_path, outcome, closed_at`; headings Trigger, Evidence, Diagnosis, Outcome, Operator package | the engine; the examiner | ADR-0017 |
| `runs/<task>/telemetry.tsv` | the token count, the latency and the wall-clock duration per gate part | `ts, task, part, model, effort, tokens, latency_ms, elapsed_s, source` | a harness or a person through `layup telemetry record` | ADR-0007, REQ-011 |
| `runs/audit/rule-changes.tsv` | the rule-path changes on the default branch | `sha, author_identity, signer, pr, paths, approval` | `layup audit rules` | ADR-0014 |
| `docs/tasks/*.md`, `docs/tasks/backlog.md`, `docs/tasks/completed.md`, `docs/adr/*.md` | the kit's task and decision records | the kit's shapes | the role agents | the kit |

The tables that the App's checks read (`gates`, `escalations`, `stalls`) live
in the target (O-54); the forge holds only the check runs and the protection
setting, which `docs/setup/branch-protection.json` mirrors and `layup setup
verify` compares (NFR-001 with its stated limit, ADR-0013).

## 4. The runner: from a pull request to a verdict (ADR-0013, ADR-0014; REQ-004, REQ-007, NFR-004)

1. A role agent, under the agents' App or the machine identity — never the LAYUP App (ADR-0014) — opens or updates a pull request on the target; its head is `X`.
2. The forge sends the event (a pull-request event, or a review by the Operator's identity) to the LAYUP App; the receiver starts a job that checks out `X` with a read token and runs `layup gate` at a pinned LAYUP commit, with the gate set of `docs/setup/stack.tsv` as it stands on the target's base branch.
3. The job appends the verdict rows (`head_sha` = `X`) to `runs/<task>/gates.tsv` as a commit `Y` on the pull request's branch under the LAYUP App's identity, and asserts that `Y` differs from `X` only under `runs/*/gates.tsv`; `<task>` is the task ID of the branch name under the target's recorded scheme.
4. The job posts, on `Y`, one check run per gate kind and the rule guard, escalation and stall checks, under the LAYUP App's identity: `pass` → `success`; `fail` and `not-active` → `failure`; an open escalation or stall row → `failure`; the rule guard reads the Operator's approval against `X`. A commit by the LAYUP App that touches only `runs/*/gates.tsv` neither starts a new run nor dismisses an approval; any other commit does both.
5. The target's protection requires each `layup/*` check by name, pinned to the LAYUP App's identity (the kit's setup step S13 body, with that departure recorded); a check with no report keeps the pull request at "expected", so a gate that did not run blocks the merge (NFR-004), and a status of the same name from another identity does not count. Review assignment is the target's own rule; the pilot audits early reviews (REQ-007).

## 5. The role model and the handoffs (ADR-0015; REQ-005, REQ-013, NFR-002)

Roles: `planner`, `implementer`, `reviewer`, `examiner`, `operator`, `idea
owner`. A handoff is a row of `handoffs.tsv` that points at a kit record at a
commit; each kind has its transition (ADR-0015 decision 2), and `layup handoff
check` validates the transition per kind, the artifact's existence and the
artifact's own kit linter (REQ-005). A second harness agent continues a task
from the table and the records alone (REQ-013, NFR-002). The transition table
and the kind list are rule paths; a target's `docs/roles.tsv` is content.

## 6. The escalation rule (ADR-0016; REQ-008)

The intent set — the facts of the problem statement and its answers, the PRD, `docs/setup/budget.md` — and the trigger
list of the target's PDR are the deterministic floor: a change to an intent
path, a dropped or deferred `Must`, a moved phase or date, a licence class
change, or a telemetry cost past a budget value selects an escalation. Every
decision record carries `class` and the four axes; `yes` or `unknown` stops the
agent. The row in `escalations.tsv` and the `layup/escalation` check hold the
pull request until the idea owner's answer is in Git; a not-confirmed escalation
is unplanned input (`F-0001#28`).

## 7. The stall procedure (ADR-0017; REQ-009, REQ-010)

`no-progress`: the same fingerprint fails a second time (N = 1 retry) or 15
minutes pass without a progress event (T); `no-agreement`: the cycle cap is
reached with a material finding open. The row in `stalls.tsv` and the
`layup/stall` check hold the pull request; the examiner — a fresh session, a
different model, a different harness where available — writes `stall-<id>.md`;
the outcome is a verified recovery or the package for the Operator, who may give
the task to a different harness agent. N and T are hypotheses the pilot resets
(O-58).

## 8. The telemetry record (ADR-0007, REQ-011)

One row per gate part per task in `runs/<task>/telemetry.tsv`, with the model,
the effort, the token count, the latency and the elapsed time, and its source (a
harness's report or a person). `layup telemetry check` counts tasks with an
incomplete record, the §7.1 measure "Telemetry Completeness" (`F-0003#60`). This
repository's task records carry the same figures in their resource-record
tables until the command exists.

## 9. Phase-1 components (the core engine)

| Requirement | Component | Status |
| ----------- | --------- | ------ |
| REQ-001 gap batch | `layup psb check` | delivered (`T-dq05`, `T-zmj6`) |
| REQ-002 setup | `layup setup`, `layup setup verify`, `docs/setup/steps.tsv` | next (#29) |
| REQ-004 stack gates | `layup gate`, the runner | next (#29); the runner after the App exists |
| REQ-007 verification on every change | the runner's required checks | with REQ-004 |
| REQ-009 stall record | `layup stall open|close`, `stalls.tsv` | next (#29) |
| REQ-011 cost visibility | `layup telemetry record|check`, `telemetry.tsv` | next (#29) |
| NFR-001 to NFR-007 | the engine's form and the state model | hold from the first release |

## 10. Open items for the implementation plan (#42 child 5)

- The setup step for the identities of O-53 (the App installation or the machine account; the Operator's identity as the approver of Decision Point 3) in `docs/setup/steps.tsv` (ADR-0014).
- The receiver of the App: a workflow in LAYUP's repository that the App's events reach, or a small hosted service; the choice is the implementation's, within ADR-0013 (the engine holds no service). If the receiver is a workflow in LAYUP's own repository, that repository must keep the App's private key from its own agents (a secret the agents' identity cannot read); ADR-0013's key-custody consequence depends on it.
- `docs/setup/budget.md` for this repository and for each target, written by the idea owner (ADR-0016, O-57).
- The two pilot problem statements and their stacks: the idea owner's Decision Point 1 (`PRD-0001` §11 question 6), recorded as an Operator decision before the pilot.
