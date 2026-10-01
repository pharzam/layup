# The implementation plan of LAYUP

This is LAYUP's implementation plan: the phases of
[`PRD-0001`](../prd/PRD-0001-layup.md) §9 with their milestones, the tasks of
phase 1 in their order, and where each part of the work goes. It is written for
LAYUP itself in the form that [`architecture.md`](../architecture.md) §8 asks of a
target's milestone plan: each task with the requirement IDs it serves, the tests
that will show it done, a size class, and its predecessors.

It comes from task `T-55n2` ([#76](https://github.com/pharzam/layup/issues/76)),
child 5 of the PDR ([#42](https://github.com/pharzam/layup/issues/42)). Its
evidence is the inventory of the phase-1 work,
[`runs/T-55n2/inventory.md`](../../runs/T-55n2/inventory.md): 103 work items and
41 cross-area conflicts (K1 to K41), found by 22 read-only agents at `7cdd346`.
The Operator's decisions O-121 to O-124 set its structure; they are quoted in
[`tasks/T-55n2.md`](../tasks/T-55n2.md).

## How to read this plan

- The plan fixes the order and the scope of each task, not its design. Each task
  still writes its own plan and gets its own plan review and review round under
  the gate ([`engineering-discipline.md`](../engineering-discipline.md), and
  [Bootstrap mode](../engineering-discipline.md#bootstrap-mode) while ADR-0012 is
  in force).
- **No task of phase 1 starts before #42 closes**, except rows 1 and 2 (O-126):
  #42 holds the core engine until the Operator approves the PDR (`T-4wrw`, O-14).
- Each task settles the defects of the specification sections it implements, in
  its own pull request ([R10](../issue-workflow.md#r10--sync-with-governance)):
  the defect register below names them. A value that the task sets is marked
  "decided here" with its reason, and a value that needs the Operator or the idea
  owner stays a marker with a row in
  [`open-gaps.tsv`](../setup/open-gaps.tsv), as
  [`spec/README.md`](../spec/README.md) says.
- Phases 2 to 4 are at milestone level only. Each milestone starts with its own
  specification task (O-114) and opens its own issues then.
- A size class is `small` (one package or one step group: about one day of agent
  work and a diff under about 400 lines) or `large`. The expected lines of a task
  are an estimate for its plan review, not a budget; each task's plan review sets
  its budget and its cycle cap (Bootstrap mode rules 2 and 3).
- The column "After" gives the order of the tasks; the row number is only a
  handle. A task starts when the tasks in its "After" cell have merged.
- `T-55n2` opens the issue of `T-b97r` and of each row; the numbers are in the
  column "Issue". Each task, in its own pull request, puts its real test names and
  their status in [`tests/traceability.md`](../tests/traceability.md) and in the
  Test column of `PRD-0001` §12.

## Milestones

| Milestone | Phase | Requirements | First demo | After | Starts with |
| --------- | ----- | ------------ | ---------- | ----- | ----------- |
| `M1` | 1 | REQ-001, REQ-002, REQ-004, REQ-007, REQ-009, REQ-011, NFR-001 to NFR-007; the `Won't` rows REQ-015 to REQ-018 hold in every phase | The first pilot (O-122): `layup` sets up one target repository from a problem statement with a Go stack, and runs that target's gate from outside (ADR-0012 part 6). | the PDR (`T-4wrw`) | the specification of phase 1 (`T-0drh`, merged) |
| `M2a` | 2 | NFR-001, NFR-002, NFR-006, REQ-002 | On an empty repository, after the Operator's root push, the records branch shows the first commit by the LAYUP App, and the Intake and control issues exist. | `M1` | its specification task |
| `M2b` | 2 | REQ-011, REQ-005, REQ-003, REQ-013, NFR-005, NFR-001 | A probe session on each registered harness, then one developer session whose commit lands on its task branch, with one complete telemetry row. | `M2a` | its specification task |
| `M2c` | 2 | REQ-001, REQ-012, NFR-003 | On a pilot problem statement, one Intake comment holds the rule gaps and the gaps of meaning, and `layup spec check --facts` passes on the numbered spans. | `M2a`, `M2b` | its specification task |
| `M2d` | 2 | REQ-002, NFR-001, NFR-002, NFR-003, NFR-006, REQ-003 | After the Operator's pushes and ruleset apply, the push to `layup-probe` is refused, and `setup/forge-check.tsv` shows each rule and probe as passed. | `M2a`, `M2b`, `M2c` | its specification task |
| `M2e` | 2 | REQ-005, REQ-007, REQ-012, REQ-003, NFR-004, NFR-001 | The specification task of a target merges after its first bet: each handoff is valid, the rendered commit holds the MoSCoW and Phase columns, and the four `layup/` checks are green. | `M2b`, `M2c`, `M2d` | its specification task |
| `M2f` | 2 | REQ-003, REQ-004, NFR-001, NFR-004 | The activation batch of a first bet merges only after each known-bad patch fails its kind, and `layup audit` shows zero agent writes to rule paths. | `M2d`, `M2e` | its specification task |
| `M2g` | 2 | REQ-005, REQ-007, REQ-004, REQ-012, NFR-004 | One build task that changes a product path merges after the activation with every gate and the four `layup/` checks green, and the idea owner's "accept" is a row of `acceptance.tsv`. | `M2e`, `M2f` | its specification task |
| `M3a` | 3 | NFR-005, REQ-006, REQ-008 | At P2 under `shadow`, a question goes to the registered provider; the answer is recorded with its model version, tokens and price; the deterministic branch decides. | `M2a`, `M2e` | its specification task |
| `M3b` | 3 | REQ-008, NFR-001 | A seeded plan that adds a dependency stops, one brief reaches the idea owner, and the answer is classed as planned input. | `M2c`, `M2g`, `M3a` | its specification task |
| `M3c` | 3 | REQ-006 | A seeded interface-contract question gets an answer from a Software Architect session, accepted without a human, with its times in `questions.tsv`. | `M2b`, `M2g`, `M3a`, `M3b` | its specification task |
| `M3d` | 3 | REQ-010, REQ-009, REQ-011 | A seeded no-progress task opens a stall with a stall row, a diagnosis row from another harness, and an outcome row. | `M2b`, `M2g`, `M3a`, `M3b` | its specification task |
| `M4a` | 4 | REQ-013, NFR-002 | Walkthrough W-12 on a test target: a task built on one harness, verified on another, and continued by a third with LAYUP removed. | `M2b`, `M2g`, `M3d` | its specification task |
| `M4b` | 4 | REQ-011, REQ-005, REQ-008, REQ-006, REQ-009, REQ-010, REQ-001 | `layup report` on a test target's records prints Telemetry Completeness and the Stall Rate, and "not comparable" where no start value exists. | `M2g`, `M3b`, `M3c`, `M3d` | its specification task |
| `M4c` | 4 | REQ-014, and each requirement that its criterion measures "in the pilot" | The pilot of `PRD-0001` phase 4: a target of a second stack passes `layup setup verify`, and after its activation batch a seeded violation of each of the four gate kinds gives `fail`. | `M2f`, `M3b`, `M4a`, `M4b` | its specification task |

The order of phases 2 to 4 is the inventory's (area `later`). The specification
task of `M2a`, the first milestone of phase 2, takes as its inputs the conflicts K38 to K41 and the open
questions of the `later` items. The learning loop (`later-p4-learning`: the
retrospective, `layup learn` and the lessons) serves no requirement of
`PRD-0001`, so it is not a milestone: it waits for the Operator's decision.

## The parent tasks

| Task ID | Issue | Task | Parent | Holds |
| ------- | ----- | ---- | ------ | ----- |
| `T-stfn` | [#29](https://github.com/pharzam/layup/issues/29) | Step 2, the LAYUP core engine | — | the four deliverables of ADR-0012 part 5, the release review and the first pilot (O-121) |
| `T-vk3k` | [#33](https://github.com/pharzam/layup/issues/33) | `layup gate`: the stack-dependent gates of a target, run from outside (ADR-0016) | #29 | rows 1 to 5 and 14 |
| `T-b97r` | [#77](https://github.com/pharzam/layup/issues/77) | `layup setup` and `layup setup verify` | #29 | rows 6 to 13, 15 and 16 |

The telemetry record (row 17, `T-tmhw`) and the stall record (row 18) are the other
two deliverables of ADR-0012 part 5, each one task. The ADR that supersedes
ADR-0012 is a separate issue that row 20 opens, with `Refs #29` (O-121). That ADR
reads the pilot's numbers (ADR-0012 part 6); this plan gives their count to the
ADR task as its first step.

## The tasks of phase 1

"Demo" is the one thing that a reader is shown when the task is done, in one
sentence ([R11](../issue-workflow.md#r11--single-goal-issues)). "Items" are the
keys of the inventory. "Fact" is the In-Scope fact that the task serves
(Bootstrap mode rule 1). "Cap" is the expected cycle cap; the plan review of
each task sets it by [Bootstrap
mode](../engineering-discipline.md#bootstrap-mode) rule 3.

| # | Task ID | Issue | Task | Demo | Parent | Items | Requirements | Fact | Tests | Size | Lines | After | Cap |
| - | ------- | ----- | ---- | ---- | ------ | ----- | ------------ | ---- | ----- | ---- | ----- | ----- | --- |
| 1 | `T-18v6` | [#78](https://github.com/pharzam/layup/issues/78) | Records: `internal/tsv`, the field rule, and the test that reads every schema block of `docs/spec/` | `go test -tags=integration ./internal/tsv/` passes: the schema parser of the new record package reads every schema block of `docs/spec/`. | `T-vk3k` | `found-tsv`, `found-spec-schema-test` | NFR-001, NFR-002, NFR-003, NFR-005, NFR-007, REQ-001, REQ-002, REQ-004, REQ-007, REQ-009, REQ-011 | F-0003#44 | unit (the writer, the reader, the types); integration (each schema block against its Go schema) | large | 900 | — | 1 |
| 2 | `T-2tc2` | [#79](https://github.com/pharzam/layup/issues/79) | `internal/git`, the one caller of `git`, and the test of the package rules | The test of the package rules proves, on the real module, that only `internal/git` starts `git`. | `T-vk3k` | `found-git`, `found-boundary-test`, `gate-nfr005-imports` | NFR-001, NFR-005, NFR-007 | F-0003#44 | unit; integration (real `git` on temporary repositories; `go list -deps` and the import rules) | large | 700 | — | 1 |
| 3 | `T-2yw7` | [#80](https://github.com/pharzam/layup/issues/80) | The command frame (usage, arguments, exit codes, progress lines) and the end-to-end harness | Each usage error of the built `layup` binary gives exit code 2, by an end-to-end test. | `T-vk3k` | `found-cli`, `found-e2e-harness` | NFR-004, NFR-005, REQ-001, REQ-002, REQ-004 | F-0003#44 | unit; integration; e2e (the binary run as a user) | large | 400 | 2 | 1 |
| 4 | `T-3jpx` | [#81](https://github.com/pharzam/layup/issues/81) | The stack catalog package: the embedded reader and a test entry | The embedded test entry of the catalog is read by its rules, with its `go.mod` and its `.github/` files. | `T-vk3k` | `gate-catalog-package`, `setup-catalog` | NFR-003, NFR-007, REQ-002, REQ-004 | F-0003#44 | unit; integration (the embedded tree read by its rules) | small | 350 | 1 | 1 |
| 5 | `T-5sgt` | [#82](https://github.com/pharzam/layup/issues/82) | `layup gate REPO --base REV --head REV` | `layup gate` on a Go repository prints one verdict per gate kind, with exit code 0 only when each kind passes or is clear. | `T-vk3k` | `gate-manifest`, `gate-scratch-tree`, `gate-results`, `gate-run`, `gate-command` | NFR-004, NFR-005, REQ-004, REQ-007 | F-0003#44, F-0003#47 | unit (the result rules in their order); integration (a scratch tree with the base's gate files); e2e (the command, exit codes 0, 1, 2; the repeat rule) | large | 1100 | 1, 2, 3 | 1 |
| 6 | `T-5zmw` | [#83](https://github.com/pharzam/layup/issues/83) | `layup psb check` to its specification: the table of S01 | `layup psb check` passes each case of its rule table, with its table written through the record package. | `T-b97r` | `psb-field-rule-cr`, `psb-tsv`, `psb-cli-contract`, `psb-rule-edges`, `psb-rule-values`, `psb-utf8-input` | NFR-004, NFR-005, NFR-007, REQ-001 | F-0003#41 | unit (each rule value and edge); integration (the `psb-gaps` block; `TestGoldenRealPSB`); e2e (the command) | large | 600 | 1, 3 | 1 |
| 7 | `T-6x75` | [#84](https://github.com/pharzam/layup/issues/84) | The frame of `layup setup verify`, the fixture harness, and the checks `kit-history`, `pin` and `identity` | `layup setup verify` on a stand-in baseline prints the rows of the checks `kit-history`, `pin` and `identity`. | `T-b97r` | `verify-frame`, `verify-harness`, `verify-test-baseline`, `verify-kit-history`, `verify-pin`, `verify-identity` | NFR-004, NFR-005, NFR-006, REQ-002 | F-0003#42 | unit; integration (the fixtures of `docs/setup/tests/` through the Go checks); e2e (the command on a stand-in baseline) | large | 1100 | 1, 2, 3 | 1 |
| 8 | `T-79y7` | [#85](https://github.com/pharzam/layup/issues/85) | The step runner: the work area, the answers rule, the stop and step tables, `commands.sh` | `layup setup WORK` resumes from its setup record, with stub steps, in the order of the step table. | `T-b97r` | `setup-runner`, `setup-answers`, `setup-commands-file` | NFR-001, NFR-003, REQ-002 | F-0003#42 | unit; integration (resume from `out/record.tsv` with a stub step; exit codes 0, 1, 2, 3); e2e (the usage and input errors, exit 2) | large | 800 | 3, 7 | 1 |
| 9 | `T-7s0y` | [#86](https://github.com/pharzam/layup/issues/86) | Steps S01 to S04, with the gap table of S01 and the record of the `S01-` and `Q-` answers at S04 (O-124) | `layup setup WORK` runs S01 to S04 on a local baseline repository, with one stop table for the missing answers of S01. | `T-b97r` | `setup-s01-questions`, `psb-batch-api`, `setup-s02-s03-baseline`, `setup-s04-pin` | NFR-001, NFR-003, NFR-006, REQ-001, REQ-002 | F-0003#42 | unit; integration (a local baseline repository by a file URL; the pin; check `pin` and check `facts` as evidence); e2e (S01 stops, and the next run goes on at S02) | large | 900 | 4, 6, 8, 10, 12 | 1 |
| 10 | `T-8vpw` | [#87](https://github.com/pharzam/layup/issues/87) | The checks `markers`, `sources`, `discipline-tests` and `link-lint` | `layup setup verify` reports the checks `markers`, `sources`, `discipline-tests` and `link-lint` on the shared fixtures. | `T-b97r` | `verify-markers`, `verify-sources`, `verify-baseline-scripts`, `gov-issue-21` | NFR-003, NFR-004, REQ-002 | F-0003#42 | unit; integration (the shared fixtures; the baseline's own scripts with `sh`); e2e (a missing baseline script gives `not-active`, exit 1); discipline (the fix of #21 in `setup-check.sh`) | large | 700 | 4, 7, 12 | 2 |
| 11 | `T-8ya0` | [#88](https://github.com/pharzam/layup/issues/88) | The check `adapted` | `layup setup verify` reports check `adapted` on the shared fixtures, with the list of the flagged files for the prose step. | `T-b97r` | `verify-adapted` | REQ-002 | F-0003#42 | unit; integration (the shared fixtures; the list of flagged files for the prose step, O-123) | large | 700 | 7 | 1 |
| 12 | `T-9t1q` | [#89](https://github.com/pharzam/layup/issues/89) | The checks `facts`, `onboarding`, `glossary` and `guardrails`, in a target's form | `layup setup verify` reports the checks `facts`, `onboarding`, `glossary` and `guardrails`, in a target's form, on the shared fixtures. | `T-b97r` | `verify-facts`, `verify-onboarding-glossary`, `verify-guardrails`, `gov-issue-48`, `gov-issue-61` | NFR-003, REQ-001, REQ-002 | F-0003#42 | unit; integration (the shared fixtures); discipline (the fixes of #48 and #61) | large | 900 | 7 | 2 |
| 13 | `T-b3r1` | [#90](https://github.com/pharzam/layup/issues/90) | Steps S05 to S11 and S14: the prose step (S07 to S09 and S14) before S10 (O-123), and the record of the `M-` answers at S11 (O-124) | `layup setup WORK` runs S05 to S11 and S14 on a local baseline, with the check of each step as its evidence. | `T-b97r` | `setup-s05-history`, `setup-s06-facts`, `setup-prose-inputs`, `setup-s10-markers`, `setup-s11-fill` | NFR-001, NFR-003, REQ-002 | F-0003#42 | unit; integration (each step with its check as evidence; one stop table per stop) | large | 1200 | 9, 10, 11, 12 | 1 |
| 14 | `T-c06a` | [#91](https://github.com/pharzam/layup/issues/91) | The Go entry of the stack catalog: kinds, tools, versions, the CI workflow, the fixtures | The fixture test of the Go entry proves each active kind on a clean tree, against its known-bad fixture. | `T-vk3k` | `gate-catalog-go-entry`, `gate-catalog-go-workflow` | NFR-002, NFR-003, NFR-004, REQ-002, REQ-004, REQ-007 | F-0003#44 | unit; integration (each active kind passes on a clean tree and fails on its fixture); e2e | large | 600 | 4, 5 | 1 |
| 15 | `T-d6q5` | [#92](https://github.com/pharzam/layup/issues/92) | Steps S12, S13 and S15, and the checks `jobs` and `gate:<kind>` | `layup setup WORK` runs S12, S13 and S15 on a local baseline, with the orphan records commit in the work area. | `T-b97r` | `setup-s12-stack-files`, `setup-s13-rulesets`, `setup-rule-paths`, `setup-s15-records-commit`, `verify-jobs`, `verify-gate-fixtures`, `gate-verify-kind`, `gov-issue-34` | NFR-001, NFR-002, NFR-003, NFR-004, NFR-005, REQ-002, REQ-004, REQ-007 | F-0003#42, F-0003#47 | unit; integration (the ruleset file, the orphan records commit, the fixture runs on no ref) | large | 1200 | 13, 14 | 1 |
| 16 | `T-dep6` | [#93](https://github.com/pharzam/layup/issues/93) | A whole setup, end to end, with no network | On a stand-in baseline with no network, a whole setup ends with `layup setup verify` at exit code 0. | `T-b97r` | `setup-e2e-offline`, `verify-acceptance`, `found-nfr001-tests`, `found-nfr002-tests` | NFR-001, NFR-002, NFR-003, NFR-004, NFR-006, REQ-002 | F-0003#42 | e2e (a whole setup and `layup setup verify` exit 0 on a stand-in baseline; the records in Git; the target passes its own gate with LAYUP absent) | large | 900 | 15 | 1 (2 if `ci.yml` changes) |
| 17 | `T-tmhw` | [#94](https://github.com/pharzam/layup/issues/94) | The telemetry record: the schemas of `telemetry.tsv` and `prices.tsv` in code | The Go schemas of `telemetry.tsv` and `prices.tsv` match their schema blocks. | `T-stfn` | `found-telemetry-schema` | REQ-011 | F-0003#50 | unit (the validator); integration (the schema blocks) | small | 300 | 1 | 1 |
| 18 | `T-dgy7` | [#95](https://github.com/pharzam/layup/issues/95) | The stall record: the schema of `stalls.tsv` in code | The Go schema of `stalls.tsv` matches its schema block. | `T-stfn` | `found-stalls-schema` | REQ-009 | F-0003#49 | unit (the validator); integration (the schema block) | small | 250 | 1 | 1 |
| 19 | `T-efmy` | [#96](https://github.com/pharzam/layup/issues/96) | The release review of phase 1 for the `Won't` rows | The release review of phase 1 is a recorded review of the code for `REQ-015` and `REQ-017`. | `T-stfn` | `gov-release-review` | REQ-015, REQ-017 | F-0003#41, F-0003#42, F-0003#44, F-0003#47, F-0003#49, F-0003#50 | uat (a recorded code review of the phase-1 release) | small | 150 | 16, 17, 18 | 1 |
| 20 | `T-evad` | [#97](https://github.com/pharzam/layup/issues/97) | The first pilot (O-122): set up one target from a Go problem statement, and run its gate from outside | `layup` sets up one target repository from a Go problem statement, whose gate it then runs from outside. | `T-stfn` | `gov-pilot-inputs`, `gov-first-pilot`, `setup-pilot-acceptance` | NFR-001, NFR-002, NFR-003, NFR-004, NFR-006, REQ-001, REQ-002, REQ-004, REQ-007, REQ-016, REQ-018 | F-0003#41, F-0003#42, F-0003#44, F-0003#47 | uat (the pilot run, with its evidence under `runs/`); e2e | large | 600 | 16, 17, 18, 19 | 1 |

The first pilot is the last task of phase 1 (it ends phase 1, O-122), so it comes
after the release review, which reviews the code that the pilot runs. Its pass and
fail rule is written in its own plan-review comment before the pilot runs
([guardrails §1](../guardrails.md)); its inputs, a problem statement with a Go
stack and an empty public repository of the Operator's account, come at its start
(O-122).

**Decided by this plan, where the order needs it** (each task's plan review can
change it):

- Row 1 keeps the schemas of each record in its owner package, with one test of
  the schema blocks per owner and one test that lists every block name (option A
  of `found-spec-schema-test`), so row 1 does not wait for K8.
- For S01, `internal/cli` gives the bytes of the `psb-gaps` table to
  `internal/setup`, which reads them with `internal/tsv` (the seam of
  `psb-batch-api`, in row 9).
- Row 13 uses the marker scanner of row 10 through `internal/cli`; it writes no
  second scanner. The test of the job names of row 14 moves to row 15, with check
  `jobs`. Row 3 keeps the helpers and the scenario of `layup version`; the e2e
  scenarios of `layup psb check` are row 6's. No test keeps a stand-in baseline in
  Git (K11).

## The edges of the inventory that the plan drops

| From | To | Reason |
| ---- | -- | ------ |
| `setup-runner` | `setup-s15-records-commit` | Row 8 builds the runner with stub steps; the behaviour of S15 lands with row 15. |
| `gov-release-review` | `gov-first-pilot` | The first pilot is the last task of phase 1 (O-122), so the release review comes before it. |

## What phase 1 proves

The plan's reading of each phase-1 criterion of `PRD-0001` §7.1. It becomes the
Operator's reading when the Operator approves the PDR (`T-4wrw`; O-122). It does not
restate a criterion: the §7.1 words stay as they are (O-115), and this table only
says which part phase 1 proves, and where.

| Requirement | Phase 1 proves (task) | A later milestone proves |
| ----------- | --------------------- | ------------------------ |
| REQ-001 | `layup psb check` writes the batch of a problem statement byte for byte (`TestGoldenRealPSB`; row 6); the answers to LAYUP's own batch are the raw fact `F-0004`, one fact per question (check `facts`; already green, `T-zmj6`); the first pilot's S01 stop table holds the rule gaps of the pilot's problem statement (row 20) | one Intake batch with the gaps of meaning (`M2c`); the Early Question Share (`M4b`) |
| REQ-002 | a setup that passes the baseline's discipline tests and `layup setup verify` with zero values without a source (rows 16 and 20); the audit of the setup values of the first pilot (row 20); the steps of the step table (rows 9, 13 and 15), and the answers in the target's Git by the reading of O-124 (rows 9 and 13) | the Scaffold driven by `layup run` (`M2d`); the pilot of phase 4 (`M4c`) |
| REQ-004 | one verdict per kind on a Go target (row 5); a known-bad fixture of each active kind fails (rows 14 and 15); a pending kind fails on a changed product path (row 5); `not-active` never counts as a pass (row 5); the baseline's own jobs stay unchanged (row 15); the gate run from outside on the pilot's target (row 20) | the fixtures of the pending kinds at activation (`M2f`); a second stack (`M4c`) |
| REQ-007 | the target's own CI job of each kind is a required check of its ruleset, and a failed job blocks the merge (rows 15 and 20); the same verdict twice on one input (row 5) | the status `layup/gates` and the rule that no pull request is ready or merges with a kind that did not run (`M2e`, `M2g`) |
| REQ-009 | the schema of `stalls.tsv` in code (row 18) | the writer and the count of the stalls with no diagnosis (`M3d`) |
| REQ-011 | the schemas of `telemetry.tsv` and `prices.tsv` in code (row 17) | the writer (`M2b`); Telemetry Completeness (`M4b`) |
| NFR-001 | the setup record, the verify table, the rule-path register and the answers in the target's Git (rows 9, 13, 15, 16 and 20; O-115, O-124) | the one writer `layup run` (`M2a`); the audit of a delivery task (`M4c`) |
| NFR-002 | a target set up in phase 1 passes its own gate with LAYUP absent (rows 16 and 20) | a fresh session of another harness continues a task (`M4a`) |
| NFR-003 | zero values without a source, by check `sources` (rows 10 and 16); the audit of the pilot's values (row 20) | — |
| NFR-004 | `layup gate` and `layup setup verify` report a check that did not run as not passed, and a fixture proves it (rows 5, 10 and 15) | — |
| NFR-005 | each verdict reproducible (rows 5 and 7); the engine checks start no model process and open no network package (row 2) | the smart-if client and `decisions.tsv` (`M3a`) |
| NFR-006 | LAYUP's own pin and check `pin` (already green); the target's pin written once and checked (rows 7 and 9) | — |
| NFR-007 | the standard library only, by `go list -deps`, and `git` only as a program (row 2) | — |
| REQ-015, REQ-017 | the recorded code review of the phase-1 release (row 19) | the review of each later release, by the specification task of each milestone |
| REQ-016 | the intent decisions of the first pilot are recorded as the idea owner's (row 20) | the pilot of phase 4 (`M4c`) |
| REQ-018 | the baseline rules of the pilot's target equal the pinned baseline's, except the values and the replaced files that the setup record names with their source (row 20; K42) | the pilot of phase 4 (`M4c`) |

## The hosts of the other items

Each inventory item that is not in a row of the task table has one host here.

| Item | Host | Reason |
| ---- | ---- | ------ |
| `gov-prd-task-column` | `T-55n2` | This task fills the Task column of `PRD-0001` §12. |
| `gov-test-traceability` | `T-55n2` | This task writes [`tests/traceability.md`](../tests/traceability.md). |
| `gov-rescope-29` | `T-55n2` | This task rewrites the children table of #29. |
| `gov-rescope-33` | `T-55n2` | This task rewrites #33 as the parent task `layup gate` (ADR-0016). |
| `gov-milestone-spec` | `T-55n2` | The milestone table: each later milestone starts with its own specification task. |
| `gov-issue-15` | `T-55n2` | This task writes a lesson in `docs/guardrails.md` §2, so it also settles the lines of #15 in that file. |
| `gov-issue-24` | `T-55n2` | Checked and closed with its evidence; this task also fixes the sentence of `docs/ci/README.md` on the restore of the check scripts, its one open part. |
| `gov-pdr-approval` | `T-4wrw` | The PDR record and the Operator's approval (#42, child 6). |
| `gov-pilot-definition` | O-122 | The Operator's reading of the first pilot. |
| `gov-pilot-numbers` | the ADR that supersedes ADR-0012 | Its rule is written in row 20's plan review before the pilot (guardrails §1); the count is the first step of the ADR task, by this plan's choice. |
| `gov-supersede-adr-0012` | the ADR that supersedes ADR-0012 | A separate issue that row 20 opens (O-121). |
| `gov-operator-setup-o112` | `M2a` | The App, key and status permissions of O-112 are needed when `layup run` starts. |
| `gov-operator-ci-protection` | out | Phase 1 adds no CI job; a task that changes `ci.yml` names the Operator's step in its own plan (K28). |
| `gov-issue-68` | `M2e` | The handoff records of the task loop. |
| `gov-issue-49` | out | It waits for the end of bootstrap mode (ADR-0012 part 4). |
| `later-psb-meaning-review` | `M2c` | The review of meaning needs a role session. |
| `later-psb-intake-batch` | `M2c` | The one Intake batch needs `layup run`. |
| `later-psb-early-question-share` | `M4b` | A measure of `layup report`. |
| `later-gate-status` | `M2e` | The status `layup/gates` needs `layup run` and the forge. |
| `later-gate-rule-batch` | `M2f` | Rule batches and activation. |
| `later-p2-run-start` | `M2a` | The milestone itself. |
| `later-p2-sessions-ledger` | `M2b` | The milestone itself. |
| `later-p2-intake-spec` | `M2c` | The milestone itself. |
| `later-p2-scaffold` | `M2d` | The milestone itself. |
| `later-p2-phase-loop-handoffs` | `M2e` | The milestone itself. |
| `later-p2-rule-protection` | `M2f` | The milestone itself. |
| `later-p2-build-accept` | `M2g` | The milestone itself. |
| `later-p3-smartif` | `M3a` | The milestone itself. |
| `later-p3-escalation` | `M3b` | The milestone itself. |
| `later-p3-clarification` | `M3c` | The milestone itself. |
| `later-p3-stalls-budget` | `M3d` | The milestone itself. |
| `later-p4-neutrality` | `M4a` | The milestone itself. |
| `later-p4-measures` | `M4b` | The milestone itself. |
| `later-p4-pilot` | `M4c` | The milestone itself. |
| `later-p4-learning` | a decision of the Operator | It serves no requirement of `PRD-0001`, so it is not a milestone. |

## The defect register

Each conflict of the inventory (K1 to K41), and K42 from the self-check, with the
task that settles it and the tasks that read the settled text. "Row *n*" is a row of the task table.

| K | Settled by | Read by | Note |
| - | ---------- | ------- | ---- |
| K1 | `T-55n2` | all rows | The order of the task table. |
| K2 | O-121 | `T-55n2` | The four deliverables as parents of child tasks. |
| K3 | `T-55n2` | — | The traceability table names the tests of row 2 by name; it waits for no build task. |
| K4 | `T-55n2` | all rows | The predecessors of the task table hold these edges. |
| K5 | `T-55n2` | rows 16, 20 | The edges name rows 9, 13 and 15, not the runner alone. |
| K6 | row 15 | row 16 | The scenarios that need a WORK through S12 move to row 16, or row 15 builds its trees by hand. |
| K7 | row 2 | rows 7, 8, 9, 10, 11, 13, 15 | The verbs of `internal/git` in `spec/packages.md`. |
| K8 | row 2 | rows 3, 7, 14, 16 | Whether the package rules bind the test files. |
| K9 | row 7 | rows 8, 10, 12, 13, 15 | The one home of the schemas that `internal/setup` and `internal/verify` share. |
| K10 | row 10 | row 13 | The call that gives S05 the files whose links break. |
| K11 | row 7 | rows 10, 11, 12, 16 | The stand-in baseline is built at test time; a marker in Go source is written as an escape. |
| K12 | O-123 (check `adapted`); row 13 (check `kit-history`) | rows 13, 20 | O-123: the prose step stops for every flagged file; row 11 builds the list of flagged files that the prose step reads. Row 13 decides here, with its reason, that S05 also removes each line of `backlog.md` and `completed.md` that links the baseline's repository. |
| K13 | row 7 | row 15 | What the evidence "checks `jobs` and `gate:<kind>`" of S12 reads. |
| K14 | O-124 | rows 6, 9, 12, 13 | The first task that implements a text that the reading changes writes the reading into it: row 12 (check `facts` in `spec/setup.md`; `spec/records.md`), row 6 (`spec/psb-check.md`), row 9 (the S04 and S06 rows). |
| K15 | row 12 | rows 9, 13 | The answer facts hold `by` and `source`. |
| K16 | O-123 | rows 8, 13 | The prose step runs before S10, and its one stop table covers more than one step. |
| K17 | row 12 | rows 9, 10, 13 | How an answers record in `docs/facts/` holds a marker, and how check `markers` reads it. |
| K18 | row 12 | rows 9, 13 | Check `facts` reads only the records that exist at its step. |
| K19 | row 5 | row 15 | A pending kind with a missing tool: the wording of `NFR-004` item 3. |
| K20 | row 14 | row 15 | The state of the test kind at setup (§6). |
| K21 | O-122 | rows 19, 20 | The table "What phase 1 proves". |
| K22 | O-122 | `T-55n2` | The glossary tells the first pilot and the pilot of phase 4 apart. |
| K23 | rows 17, 18 | `M2b`, `M3d` | The package rows of the phase-1 types. |
| K24 | rows 14, 15 | row 16 | Row 14 gives the coverage floor its catalog form; S12 (row 15) writes it as an open gap (L-B2). |
| K25 | `T-55n2` | row 20 | This plan's choice: the pilot rule in row 20's plan review (guardrails §1); the count in the ADR task. Row 20's plan review can change it. |
| K26 | `T-55n2` | rows 1 to 5, 14 | #33 rewritten as the parent `layup gate`. |
| K27 | row 16 | — | One shared whole-setup run, inside the e2e time limit. |
| K28 | row 16 | — | No CI change by default; a change gets cap 2 and the Operator's step. |
| K29 | row 3 | rows 5, 7, 8, 9, 10, 15 | One rule for progress lines on standard error. |
| K30 | row 14 | row 15 | The name and the ID of a CI job are the kind. |
| K31 | row 2 | row 9 | `git` reads no credential helper of the host; the baseline's repository is public, and a private baseline is a known limit. |
| K32 | row 1 | row 6 | One rule for input that is not UTF-8. |
| K33 | row 8 | row 9 | The stop table gives the line of a gap. |
| K34 | row 3 | — | Extra words after `layup version`. |
| K35 | row 4 | rows 14, 15 | The catalog's rename rule for `go.mod` and `.github/` under `embed`. |
| K36 | `T-55n2` | — | #24 closed with its evidence; #15 settled with this task's lesson. |
| K37 | `T-55n2` | — | Those items are out of phase 1 (the hosts table). |
| K38 | `M2a` | — | An input of the specification task of phase 2. |
| K39 | row 2 | `M2a`, `M3a` | Phase 1: no package imports a network package; `M2a` names the forge adapter and `M3a` names `internal/smartif`. |
| K40 | `M2a` | — | An input of the specification task of phase 2. |
| K41 | `M2a` | — | An input of the specification task of phase 2. |
| K42 | row 13 | rows 16, 20 | Found by the self-check: the setup record names a source for each value, but no row for a kept file that the prose step replaces (O-123). The prose step writes one record row for each replaced file, so `REQ-018` has its evidence. |

Each task also takes the open questions and the duplicate resolutions of its own
items in the inventory as inputs of its plan.

## The open issues

| Issue | Disposition |
| ----- | ----------- |
| [#15](https://github.com/pharzam/layup/issues/15) | Settled by `T-55n2` (its lines in `docs/guardrails.md`). |
| [#21](https://github.com/pharzam/layup/issues/21) | Row 10 (`T-8vpw`), with check `markers`; cap 2. |
| [#24](https://github.com/pharzam/layup/issues/24) | Checked and closed by `T-55n2`, which also fixes its one open sentence in `docs/ci/README.md`. |
| [#34](https://github.com/pharzam/layup/issues/34) | Row 15 (`T-d6q5`), with the rule-path register. |
| [#48](https://github.com/pharzam/layup/issues/48) | Row 12 (`T-9t1q`), with check `facts`; cap 2. |
| [#49](https://github.com/pharzam/layup/issues/49) | Waits for the end of bootstrap mode (ADR-0012 part 4). |
| [#61](https://github.com/pharzam/layup/issues/61) | Row 12 (`T-9t1q`), with check `facts`; cap 2. |
| [#68](https://github.com/pharzam/layup/issues/68) | Milestone `M2e`, with the handoff records. |

## Known limits of this plan

- **The inventory is not complete.** A fresh review of a task's sections can find
  more defects; the task settles them in its own pull request.
- **The expected lines are estimates.** Each task's plan review sets its budget;
  a task that grows past its budget goes back to its issue as a new slice or a
  child issue (R12).
- **The milestones of phases 2 to 4 are coarse.** Each one is fixed by its own
  specification task, which can change its scope and its order.
