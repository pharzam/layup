# 0015. Roles are the kit's gate roles, and a handoff is a TSV row that points to a kit record

Date: 2026-09-25

## Status

Accepted

## Context

The PSB requires that information which passes between role agents has a form a machine can validate, based on Armature conventions (`F-0003#45`, REQ-005), and it sets no list of roles (`F-0001#21`). REQ-005's criterion requires that every role transition of a pilot task carries a record that validates against its schema, and that the schema encodes the Armature conventions it is based on, so that a schema which accepts everything does not pass. [ADR-0011](0011-structure-the-core-engine-as-a-go-cli-over-repository-files.md) decision 2 fixes the shape of state: tab-separated tables with a fixed header, rows only added, beside Markdown records in the kit's shape. Question 3 of `PRD-0001` §11.

The panel (`runs/T-7qvc/`) compared: the kit's gate roles with one TSV event table that points at kit records (B 3A, C 3-A); four delivery roles by stage (A H1); artifact owners by question kind (A H2); the PSB's functional roles — specifier, architect, implementer, tester, reviewer, examiner — with a TSV event table whose columns name the harness, the model and the requirements (C 3-B); a Markdown record per handoff with a new linter in the kit copy (B 3B); the forge as carrier (B 3C); and C's 3-C, that the schema is a rule and the role names are content. The Operator decided (O-55, #66): the kit's gate roles — planner, implementer, reviewer, examiner, operator — with the handoff records being the plan, the review record, the verdict and the stall record, validated by the kit's linters plus a transition table. The `idea owner` is added as the human role of Decision Points 1, 2 and 4 (`F-0001#23`; C 3-A).

## Decision

We will model roles and handoffs as follows:

1. **The roles are the kit's own steps.** `planner` (the plan and its answers), `implementer` (the tests and the code), `reviewer` (a plan review or a review round, a different model from the author's), `examiner` (the fresh-context diagnosis of a stall, [ADR-0017](0017-stop-a-stall-at-a-counted-limit-and-examine-it-fresh.md)), `operator` (the human who starts, monitors and approves) and `idea owner` (the human of Decision Points 1, 2 and 4). A specialist function of PSB §2 (a QA engineer, a domain expert) acts through one of these roles; a target may name its specialists in a register, `docs/roles.tsv` (`role, specialism, harness_allowed`), which is **content**, never a rule (`F-0001#7`; C 3-C).
2. **A handoff is one row** of `runs/<task>/handoffs.tsv`, append-only, with the header `ts, task, from_role, to_role, kind, artifact_path, artifact_sha, harness, model, verdict`. The `kind` is one of the kit's record kinds, and each kind has its transition: `plan` (planner → reviewer for the review; planner → implementer once the plan review approves it); `plan-review` (reviewer → planner); `question` (any role → operator, or → the responsible role of [ADR-0016](0016-escalate-by-a-deterministic-floor-and-a-declared-class.md)'s clarification); `change`, a pull-request head (implementer → reviewer); `handover` (implementer → planner for the close-out); `review-record` (reviewer → implementer when it ends `material`, reviewer → operator when it is the last round); `fixes` (implementer → reviewer); `verdict` (planner → operator, the close-out); `decision`, an ADR or a decision note (planner → reviewer); `resource-record` (planner → operator); `escalation` (any role → idea owner, [ADR-0016](0016-escalate-by-a-deterministic-floor-and-a-declared-class.md)); `answer` (idea owner → the role that escalated; operator → the role that asked); `stall-record` (any role → examiner when the stall opens; examiner → operator when the diagnosis is written; operator → the role that continues). The artifact is the kit record itself, at the commit the row names; for `change` it is the pull-request head; for `answer` it is the record that holds the answer.
3. **The engine validates each row**, `layup handoff check <checkout>`: the header and the column count; the `kind` against the list; `from_role → to_role` against the transition of that kind (the pair is checked per kind, not alone); that `artifact_path` exists at `artifact_sha`; and the artifact's own kit linter where one exists (`review-record-lint` for a review record, `adr-lint` for a decision, the resource-record shape of ADR-0007). A row with a legal header but a bad transition, a missing artifact or an empty record fails: the schema does not accept everything.
4. **The transition table and the kind list are rule paths** ([ADR-0014](0014-protect-rule-paths-with-a-rule-guard-check.md)), pinned in LAYUP; a target's `docs/roles.tsv` is not.

We rejected: **four delivery roles** and **artifact owners** — role lists with no evidence before a pilot (Invariant 4), and the architect or an unowned question kind becomes the message bus in agent form; **a Markdown record per handoff with a new linter in the kit copy** — kit content that LAYUP would add to a target (O-10); **the forge as carrier** — a record that lives only on the forge until mirrored (Invariant 1).

## Consequences

- REQ-005's criterion is testable: validated transitions divided by all transitions in a pilot task, and a fixture row per failure cause (bad transition, missing artifact, empty record).
- A second harness agent continues a task from `handoffs.tsv` and the records it names, with no chat history (`F-0003#66`, Invariant 9); the pilot's NFR-002 criterion exercises it.
- The role list may grow with pilot evidence; a new role or transition is a change to a rule path, approved as Decision Point 3 (O-61).
- The engine gains `layup handoff check`; the two kit linters it calls stay the kit's (Invariant 2).
