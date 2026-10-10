# The implementation plan of LAYUP

This is LAYUP's implementation plan: the phases of
[`PRD-0001`](../prd/PRD-0001-layup.md) §9 with their milestones, the tasks of
phase 1 and of the milestones `M2a` and `M2b` in their order, and where each part of the work goes. It is written for
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
  the gate ([`engineering-discipline.md`](../engineering-discipline.md), as
  [ADR-0026](../adr/0026-keep-the-bootstrap-review-rules-as-the-standing-gate.md) sets it; phase 1 ran under
  [Bootstrap mode](../engineering-discipline.md#bootstrap-mode)).
- **No task of phase 1 starts before #42 closes**, except rows 1 and 2 (O-126):
  #42 holds the core engine until the Operator approves the PDR (`T-4wrw`, O-14).
- Each task settles the defects of the specification sections it implements, in
  its own pull request ([R10](../issue-workflow.md#r10--sync-with-governance)):
  the defect register below names them. A value that the task sets is marked
  "decided here" with its reason, and a value that needs the Operator or the idea
  owner stays a marker with a row in
  [`open-gaps.tsv`](../setup/open-gaps.tsv), as
  [`spec/README.md`](../spec/README.md) says.
- Phases 2 to 4 are at milestone level only, except `M2a` and `M2b`, whose tasks
  are below ([`M2a`](#the-tasks-of-m2a), [`M2b`](#the-tasks-of-m2b)). Each milestone starts with its own specification
  task (O-114), and a slicing task then writes its rows and opens their issues
  (O-164: `T-zwke` for `M2a`, `T-fdaq` for `M2b`).
- A size class is `small` (one package or one step group: about one day of agent
  work and a diff under about 400 lines) or `large`. The expected lines of a task
  are an estimate for its plan review, not a budget; each task's plan review sets
  its budget and its cycle cap ([R12](../issue-workflow.md#r12--slice-and-prioritize)
  and [the cycle cap](../engineering-discipline.md#reviewing-until-findings-decay)).
- The column "After" gives the order of the tasks; the row number is only a
  handle. A task starts when the tasks in its "After" cell have merged.
- `T-55n2` opens the issue of `T-b97r` and of each row of phase 1, `T-zwke`
  of each row of `M2a`, and `T-fdaq` of each row of `M2b`; the numbers are in the column "Issue". Each task, in its own pull request, puts its real test names and
  their status in [`tests/traceability.md`](../tests/traceability.md) and in the
  Test column of `PRD-0001` §12.

## Milestones

| Milestone | Phase | Requirements | First demo | After | Starts with |
| --------- | ----- | ------------ | ---------- | ----- | ----------- |
| `M1` | 1 | REQ-001, REQ-002, REQ-004, REQ-007, REQ-009, REQ-011, NFR-001 to NFR-007; the `Won't` rows REQ-015 to REQ-018 hold in every phase | The first pilot (O-122): `layup` sets up one target repository from a problem statement with a Go stack, and runs that target's gate from outside (ADR-0012 part 6). | the PDR (`T-4wrw`) | the specification of phase 1 (`T-0drh`, merged) |
| `M2a` | 2 | NFR-001, NFR-002, NFR-006, REQ-002 | On an empty repository, after the Operator's root push, the records branch shows the first commit by the LAYUP App, and the Intake and control issues exist. | `M1` | its specification task (`T-zck8`); its tasks: [rows 21 to 27](#the-tasks-of-m2a) |
| `M2b` | 2 | REQ-011, REQ-005, REQ-003, REQ-013, NFR-005, NFR-001 | A probe session on each registered harness, then one developer session whose commit lands on its task branch, with one complete telemetry row. | `M2a` | its specification task (`T-ywk7`, #147); its tasks: [rows 28 to 39b](#the-tasks-of-m2b) |
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

The specification tasks of `M2b`, `M2f`, `M2g`, `M3a` and `M3d`, and the
Operator's decision on the learning loop, also read
[`runs/T-w73g/survey.md`](../../runs/T-w73g/survey.md) (`T-w73g`,
[#133](https://github.com/pharzam/layup/issues/133)): the patterns of five public
repositories that the Operator named, each with its source path and the
milestone that reads it. The survey decides nothing: each task takes a pattern
or rejects it, and still runs its own public-solution search.

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
([the scope of a task](../engineering-discipline.md#working-a-task-under-the-quality-gate)).
"Cap" is the expected cycle cap; the plan review of each task sets it by
[the cycle cap rule](../engineering-discipline.md#reviewing-until-findings-decay).

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
| 9 | `T-7s0y` | [#86](https://github.com/pharzam/layup/issues/86) | Steps S01 to S04, with the gap table of S01 and the record of the `S01-` and `Q-` answers at S04 (O-124) | `layup setup WORK` runs S01 to S04 on a local baseline repository, with one stop table for the missing answers of S01. | `T-b97r` | `setup-s01-questions`, `psb-batch-api`, `setup-s02-s03-baseline`, `setup-s04-pin` | NFR-001, NFR-003, NFR-006, REQ-001, REQ-002 | F-0003#42 | unit; integration (a local baseline repository by a file URL; the pin; check `pin` and check `facts` as evidence); e2e (S01 stops, and the next run goes on at S02) | large | 900 | 4, 6, 8, 10, 12, 14 | 1 |
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

The phase-1 release is the code of `main` that the first pilot runs (the
non-test Go files of `cmd/` and `internal/`, `go.mod`, and the files that the
binary embeds), at the commit that the release review reads; it is not a build,
as phase 1 ships no binary (task `T-efmy`, #96). A task of phase 1 that changes
that code after the release review checks its own diff for `REQ-015` and
`REQ-017` in its review round, and runs
[`release-check.sh`](../../runs/T-efmy/release-check.sh) with the full commit ID
of its head.

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

## The tasks of M2a

The build tasks of milestone `M2a`, the Start of `layup run` (task `T-zwke`, #124;
O-164 of #123). Its specification is [`spec/run.md`](../spec/run.md),
[`spec/forge.md`](../spec/forge.md), [the records of
Start](../spec/records.md#nfr-001--the-records-of-start) and [the table of
M2a](../spec/packages.md#the-table-of-m2a). The rule of the slicing: one task per
package boundary of the table of `M2a`, in the order of its imports; each part of
the specification is in one row ([`runs/T-zwke/parts.tsv`](../../runs/T-zwke/parts.tsv)),
and each acceptance test of `run.md` in the row that delivers its part. The
columns are those of the tasks of phase 1 without Parent, which is `—` for each
row; "Task" also names the parts of the specification; each issue says `Refs #123` and `Refs #29`. Rows
21, 22a and 23 can start at once (row 22 split into 22a and 22b by O-170 of
#127); rows 21 and 22b both change the lists of the block test, so the second to
merge takes the first's change by a merge of `origin/main`. Row 25 split into
25a and 25b by O-173 of #130.

| # | Task ID | Issue | Task | Demo | Items | Requirements | Fact | Tests | Size | Lines | After | Cap |
| - | ------- | ----- | ---- | ---- | ----- | ------------ | ---- | ----- | ---- | ----- | ----- | --- |
| 21 | `T-8kqn` | [#126](https://github.com/pharzam/layup/issues/126) | The records of Start: `records.md`: the blocks `start`, `approvers`, `lease`, `copies`, and the row rules that they give in words; `packages.md`: the Job cell of `internal/records` in the table of phase 1 | The Go schemas of `start`, `approvers`, `lease` and `copies` equal their blocks. | `later-p2-run-start` | NFR-001 | F-0003#42 | unit; the block test (the four names to `built`) | large | 500 | — | 1 |
| 22a | `T-esfe` | [#127](https://github.com/pharzam/layup/issues/127) | The package rules of M2a: `packages.md`: "The table of M2a", rule 5 by "Connects" and the line of rule 5, read by `TestPackageRules` (the checker, its unit tests, the fixture `testdata/netimport`), before any package of that table exists (O-170) | `TestPackageRules` refuses a direct import of a network package outside the package that rule 5 names, with "The table of M2a" read. | `later-p2-run-start` | NFR-007 | F-0003#42 | unit (the checker); integration (`TestPackageRules`) | small | 500 | — | 1 |
| 22b | `T-1g1q` | [#137](https://github.com/pharzam/layup/issues/137) | The two host registers and the key file (O-171 of #137 moved the forge interface to row 24): `forge.md`: the key-file checks of "The App identity"; `records.md`: the blocks `forge-register`, `harness-register`; `run.md`: rows 1 to 5 of "Input states" | A forge register or a harness register that breaks its schema is refused with its line, by its owner package. | `later-p2-run-start` | NFR-001, NFR-007 | F-0003#42 | unit; the block test (the two names to `built`); each owner's `tsv.Compare` | large | 700 | 22a | 1 |
| 23 | `T-xhgz` | [#128](https://github.com/pharzam/layup/issues/128) | The calls Fetch and Push of internal/git: `packages.md`: "The calls of `M2a`"; `forge.md`: the token to `git` in one call's environment | `Push` refuses a push to a local bare repository that is not a fast-forward. | `later-p2-run-start` | NFR-001 | F-0003#42 | unit; integration (real `git`; the fixed list of the environment) | small | 300 | — | 1 |
| 24 | `T-6bq5` | [#129](https://github.com/pharzam/layup/issues/129) | The GitHub adapter, and by O-171 of #137 the forge interface of row 22b, in `internal/forge`: `forge.md`: "The six capabilities" (the neutral interface of the calls of `M2a`, its types, the six capabilities with their declaration, the permissions of `M2a` and their check), with conditions 4 to 6 and note 5 of comment 6054282765 of #137, open; `forge.md`: "The App identity" (the JWT and the installation token), "The calls of M2a", "Forge errors" (with a progress function that the adapter takes for the rate-limit wait), "The test of the adapter" | Against a loopback server, the adapter plays each call of `M2a` with an installation token that it made from a JWT. | `later-p2-run-start` | NFR-001, NFR-007 | F-0003#42 | unit; integration (`httptest`) | large | 1000 | 22b | 1 |
| 25a | `T-trej` | [#130](https://github.com/pharzam/layup/issues/130) | The rules of a run (O-173 of #130 split row 25): `run.md`: "The lease and fencing", "A human decision", "Copy before read", and their acceptance rows | With a stand-in clock and a stand-in records store, a second run takes the lease only after the heartbeat has not moved for `3 × lease.H` by its own clock. | `later-p2-run-start` | NFR-001 | F-0003#42 | unit (the lease with a stand-in clock and a stand-in records store; the two rules) | large | 500 | 21, 24 | 1 |
| 25b | `T-ax3r` | [#142](https://github.com/pharzam/layup/issues/142) | internal/run: Start and the restart, on the rules of 25a: `run.md`: the one-writer sentence of NFR-001, "The steps of `layup run --new`", "The README of a target that Start makes", "The restart", the heartbeat from step 6, "The Intake and control issues", rows 6 to 11 of "Input states" and their known limit, the table `run-steps`, NFR-006, NFR-002, REQ-002 | Against the `httptest` forge and a local bare repository, Start makes the first records commit with the pin, the briefs, `approvers.tsv`, `start.tsv` and the lease row. | `later-p2-run-start` | NFR-001, NFR-002, NFR-006, REQ-002 | F-0003#42 | unit (the input states); integration (Start, and a restart with another LAYUP version); the block test (`run-steps` to `built`) | large | 800 | 22b, 23, 25a | 1 |
| 26 | `T-mqty` | [#131](https://github.com/pharzam/layup/issues/131) | The command layup run: `run.md`: "The command" (two rows, the selecting flag, the one flag that may be left out, each flag and its value rule, the printing of the table and the progress lines, the exit codes); `README.md`: the rules of the selecting flag and of a flag that may be left out; `packages.md`: the May import cell of `internal/cli` | `layup run --new` then `layup run TARGET`, as the built binary against a local fake forge, give their tables, the same bytes on a repeat. | `later-p2-run-start` | NFR-001, NFR-002 | F-0003#42 | unit; e2e (the binary, no secret) | large | 600 | 25b | 1 |
| 27 | `T-fnsr` | [#132](https://github.com/pharzam/layup/issues/132) | The demo of M2a (uat): `run.md`: the uat row of "The acceptance tests of M2a" | On a real GitHub test target, the Operator reads with a plain `git clone` the records branch that Start wrote as the App's bot. | `later-p2-run-start`, `gov-operator-setup-o112` | NFR-001, NFR-002, NFR-006, REQ-002 | F-0003#42 | uat, with its evidence under `runs/` | small | 150 | 26 | 1 |

Row 27 needs the Operator's inputs, which its issue lists (O-112).

## The tasks of M2b

The build tasks of milestone `M2b`, role sessions on registered harnesses (task
`T-fdaq`, #152; O-164 of #123). Its specification is
[`spec/session.md`](../spec/session.md), [the records of a
session](../spec/records.md#nfr-001--the-records-of-a-session), [the table of
M2b](../spec/packages.md#the-table-of-m2b) and [the calls of
`M2b`](../spec/packages.md#the-calls-of-internalgit). The rule of the slicing,
as for `M2a`: one task per package boundary, in the order of the imports; each
part of the specification in one row, and each acceptance test of `session.md`,
or each clause of one that the slicing cuts, in the row that delivers its part
([`runs/T-fdaq/parts.tsv`](../../runs/T-fdaq/parts.tsv)). Two lessons of `M2a`
are applied before any issue opened: the package rules come first, before any
package of their table exists (row 29, as O-170 did for row 22a), and the goal
classes of each row were counted (`docs/guardrails.md` §2). The Operator split
rows 30, 33, 34, 36, 37 and 39 by their classes (O-186 and O-187 of #152), and
the count of each row is final (O-187); its issue states it. The columns are
those of the tasks of `M2a`; each issue says `Refs #147` and `Refs #29`. Rows
28, 29, 30a and 31 can start at once.

| # | Task ID | Issue | Task | Demo | Items | Requirements | Fact | Tests | Size | Lines | After | Cap |
| - | ------- | ----- | ---- | ---- | ----- | ------------ | ---- | ----- | ---- | ----- | ----- | --- |
| 28 | `T-3py1` | [#153](https://github.com/pharzam/layup/issues/153) | The records of a session: `session.md`: The acceptance tests of M2b (The records of a session: sessions, harnesses, routing, events, result); `records.md`: NFR-001 — The records of a session (the blocks sessions, harnesses, routing, events, result, and the row rules in words) | The Go schemas of `sessions`, `harnesses`, `routing`, `events` and `result` equal their blocks. | `later-p2-sessions-ledger` | NFR-001, REQ-005, REQ-013 | F-0003#52, F-0003#45 | unit; the block test (the five names to `built`) | large | 600 | — | 1 |
| 29 | `T-y10b` | [#154](https://github.com/pharzam/layup/issues/154) | The package rules of M2b: `session.md`: NFR-005 — No harness in the engine checks; `packages.md`: The table of M2b (the table read by TestPackageRules, the line of the engine checks, the form of Starts a program for a register row) | `TestPackageRules` refuses a package of the engine checks that depends on `internal/session`, with "The table of M2b" read. | `later-p2-sessions-ledger` | NFR-005, NFR-007 | F-0003#52 | unit (the checker); integration (`TestPackageRules`) | large | 500 | — | 1 (2 if its plan review reads the test as a gate) |
| 30a | `T-ysph` | [#155](https://github.com/pharzam/layup/issues/155) | The host registers of M2b: `session.md`: Input states (rows 1 to 7: the host registers, the credential, credential_to, vars, the missing files, a row of an unknown harness); `session.md`: The acceptance tests of M2b (The records of a session: models, routing-register); `records.md`: NFR-001 — The records of a session (the blocks models, routing-register, and the columns that M2b adds to harness-register) | A harness register row whose `vars` names `GH_TOKEN` is refused with its line, by `internal/route`. | `later-p2-sessions-ledger` | REQ-013, NFR-001 | F-0003#52 | unit; the block test (`models`, `routing-register` to `built`; `harness-register` with its new columns and the fixtures that write it); `tsv.Compare` | large | 600 | — | 1 |
| 30b | `T-cht1` | [#156](https://github.com/pharzam/layup/issues/156) | Admission and the routing order: `session.md`: The start of a session (the refusal pair); `session.md`: Admission; `session.md`: The routing register (the first admitted pair of the role's list for the task's tier, and the refusal pair); `session.md`: Input states (row 8: a harness with no model of use yes); `session.md`: The acceptance tests of M2b (The probe and admission: admission by the probe and use); `session.md`: The acceptance tests of M2b (A refused start: the refusal pair); `session.md`: The acceptance tests of M2b (The routing register: the pair and its refusal) | With a stand-in harnesses record, the pair of a session is the first pair of its role's list whose harness passed its probe at the version read. | `later-p2-sessions-ledger` | REQ-013 | F-0003#52 | unit | small | 350 | 28, 30a | 1 |
| 31 | `T-z5dj` | [#157](https://github.com/pharzam/layup/issues/157) | The calls of M2b of internal/git: `session.md`: The acceptance tests of M2b (Before a push: the session configuration that holds each key that starts a program); `packages.md`: The calls of `M2b` (CloneLocal, FetchSession, IsAncestor, DiffFile, DiffBinary) | `FetchSession` fetches a session's head into another clone while a session configuration that holds each key of git's documentation that starts a program runs none. | `later-p2-sessions-ledger` | REQ-003, NFR-001 | F-0003#43 | unit; integration (real `git`, the hostile session configuration) | large | 600 | — | 1 |
| 32 | `T-m1dx` | [#158](https://github.com/pharzam/layup/issues/158) | The rule-path check before a push: `session.md`: A workflow or rule-path change (the reader of rule-paths.tsv, the match, the exception of docs/guardrails.md, the refusal of .github/workflows/); `session.md`: The acceptance tests of M2b (Before a push: added lines in §2 of docs/guardrails.md pass) | `internal/rules` refuses a change of a rule path but passes added lines in §2 of `docs/guardrails.md`. | `later-p2-sessions-ledger` | REQ-003 | F-0003#43 | unit | small | 350 | 29 | 1 |
| 33a | `T-vxdg` | [#159](https://github.com/pharzam/layup/issues/159) | The directory and the environment of a session: `session.md`: The session directory (the directory, its clone, home/ and tmp/, and the removal of a directory that a stopped run left); `session.md`: The environment and the harness credential; `session.md`: The acceptance tests of M2b (The session directory: the clone, home/ and tmp/, and the removal of a directory that a stopped run left); `session.md`: The acceptance tests of M2b (The environment); `records.md`: NFR-001 — The records of a session (the session directory is not a record); `docs/spec/README.md`: Commands (the read of PATH that TestInputRule allows) | In a real temporary tree, a session directory is a `--no-local` clone with no remote on `task/<task>/<attempt>`, whose environment holds the named list only. | `later-p2-sessions-ledger` | REQ-013, REQ-003 | F-0003#52 | unit (the environment and the credential); integration (the directory, real `git`) | large | 500 | 29, 30a, 31 | 1 |
| 33b | `T-6sbe` | [#160](https://github.com/pharzam/layup/issues/160) | The checks before the start of a session: `session.md`: Rules only from the target; `session.md`: The start of a session (step 1: the session ID); `session.md`: The start of a session (step 3: the context check, the size check of the prompt and the rule-file check); `session.md`: The context of a start; `session.md`: The fetch by SHA (the head read from the files of repo/.git); `session.md`: Input states (row 14: a prompt.md over 131,071 bytes); `session.md`: Input states (row 15: a malformed ref of repo/.git); `session.md`: The acceptance tests of M2b (Rules only from the target); `session.md`: The acceptance tests of M2b (The context of a start); `session.md`: The acceptance tests of M2b (Before a push: the head read from the files of repo/.git) | Before any process starts, a `CLAUDE.md` above the session directory refuses the start with `rules`. | `later-p2-sessions-ledger` | REQ-013, REQ-003 | F-0003#52 | unit (the context, the prompt size, the head); integration (the rule files in a real temporary tree) | small | 400 | 33a | 1 |
| 34a | `T-5pxd` | [#161](https://github.com/pharzam/layup/issues/161) | The process of a session and its stop: `session.md`: REQ-013 — A role session (the one call of internal/session and the order of its steps); `session.md`: The start of a session (step 2: the version check); `session.md`: The start of a session (step 5: the process); `session.md`: The limit of a session; `session.md`: Input states (row 9: a version command that fails); `session.md`: The acceptance tests of M2b (The start, the limit and the end: the version check, the fake that waits for its input, the stop at wall, the output cap; the version check was in The probe and admission until task `T-5pxd`) | With a fake harness program, a session that runs past `wall` is stopped by `SIGINT`, `SIGTERM` and `SIGKILL`. | `later-p2-sessions-ledger` | REQ-013 | F-0003#52 | integration (a fake harness: the version check, the input, the stop at `wall`, the output cap) | large | 600 | 30a, 33a | 1 |
| 34b | `T-bpxg` | [#162](https://github.com/pharzam/layup/issues/162) | The end of a session: `session.md`: The end of a session; `session.md`: The usage report of a harness; `session.md`: Input states (row 10: a result file that is missing, too large or malformed); `session.md`: Input states (row 12: stdout with no result object); `session.md`: The acceptance tests of M2b (The start, the limit and the end: each class of the end); `session.md`: The acceptance tests of M2b (The usage report); `session.md`: The acceptance tests of M2b (The records of a session: probe-result); `records.md`: NFR-001 — The records of a session (the block probe-result) | With a fake harness program, each way a session ends gets its class: `done`, `start`, `crash`, `no-result`, `wall` or `output`. | `later-p2-sessions-ledger` | REQ-013, REQ-005, REQ-011 | F-0003#52, F-0003#50 | unit (the usage report on two recorded `result` events); integration (each class of the end); the block test (`probe-result` to `built`) | large | 600 | 28, 34a | 1 |
| 35 | `T-4c3q` | [#163](https://github.com/pharzam/layup/issues/163) | The writer of the telemetry record: `session.md`: REQ-011 — The writer of the telemetry record; `session.md`: The acceptance tests of M2b (The writer) | One row per session passes `CheckTelemetry`, with its money `reported`, `computed` or `unknown` by the billing, the prices and the models. | `later-p2-sessions-ledger` | REQ-011 | F-0003#50 | unit | small | 350 | 29 | 1 |
| 36a | `T-d8t9` | [#164](https://github.com/pharzam/layup/issues/164) | The start of a task session and its refusals: `session.md`: The start of a session (the call that starts an attempt (the event attempt), and the attempt check of step 1); `session.md`: The start of a session (step 4: the start row and the event session, pushed before the process); `session.md`: The start of a session (a refused start of a task session); `session.md`: Input states (row 13: a task with no event attempt); `session.md`: Input states (row 16: a records push that is refused); `session.md`: The acceptance tests of M2b (The start, the limit and the end: the start row pushed before the process); `session.md`: The acceptance tests of M2b (The open attempt: a session with no event attempt of its attempt); `session.md`: The acceptance tests of M2b (A refused start: a task session's refusals) | Against the `httptest` forge and a local bare repository, a task session's start row is pushed before its fake harness starts. | `later-p2-sessions-ledger` | REQ-013, NFR-001 | F-0003#52 | unit (each refused start; the attempt); integration (a task session with a fake harness, the `httptest` forge, a local bare repository) | large | 600 | 28, 30b, 33b, 34a | 1 |
| 36b | `T-fsjp` | [#165](https://github.com/pharzam/layup/issues/165) | The result and the end of a task session: `session.md`: The fetch by SHA (the fetch through FetchSession, from row 37a by task `T-fsjp`, condition 1 of the plan review of #165); `forge.md`: The calls of M2b (`Comment`, O-189); `session.md`: The session directory (the removal after the session's records are pushed); `session.md`: A comment for a session (the comment of a task session); `session.md`: REQ-005 — The result of a session; `session.md`: The open attempt; `session.md`: The acceptance tests of M2b (The session directory: the removal after the session's records are pushed); `session.md`: The acceptance tests of M2b (The result of a session); `session.md`: The acceptance tests of M2b (The open attempt: a result of a closed, replaced or rebased attempt) | A task session's result file is committed byte for byte as `tasks/<task>/results/<session>.tsv`, in the records commit that holds its telemetry row. | `later-p2-sessions-ledger` | REQ-005, REQ-011, NFR-001 | F-0003#45, F-0003#50 | unit (the open attempt with a stand-in events table); integration (the result, the telemetry row, the comment, the removal of the directory) | large | 500 | 34b, 35, 36a | 1 |
| 37a | `T-z027` | [#166](https://github.com/pharzam/layup/issues/166) | The checks before a push: `session.md`: REQ-003 — Before a push (the checks of a task session whose end is done); `session.md`: The fetch by SHA (the check of the base; the fetch through FetchSession is row 36b's since task `T-fsjp`); `session.md`: A workflow or rule-path change (the call before a push, the refused diff as a payload, the proposal); `session.md`: The acceptance tests of M2b (Before a push: the fetch with hooks off, and a head that does not descend from the base); `session.md`: The acceptance tests of M2b (Before a push: a change of .github/workflows/ or of a rule path refused before any push, its diff a payload); `records.md`: NFR-001 — The records of a session (the rule of payloads/<sha256>) | A session head that changes a rule path is refused before any push, with its diff as a payload. | `later-p2-sessions-ledger` | REQ-003 | F-0003#43 | integration (real `git`, a local bare repository) | large | 500 | 31, 32, 33b, 36b | 1 |
| 37b | `T-e3sy` | [#167](https://github.com/pharzam/layup/issues/167) | The push and the bind: `session.md`: The push and the bind; `session.md`: The acceptance tests of M2b (Before a push: the SHA bound only after the push is accepted) | A clean session head is bound only after the forge accepts its push. | `later-p2-sessions-ledger` | REQ-003, NFR-001 | F-0003#43 | integration (real `git`, a local bare repository, the `httptest` forge) | small | 300 | 37a | 1 |
| 38 | `T-nxe4` | [#168](https://github.com/pharzam/layup/issues/168) | The step probe: `session.md`: The start of a session (a probe's refused start); `session.md`: A comment for a session (the comment of a probe); `session.md`: The probe; `session.md`: The routing register (the copy into records:routing.tsv at the step probe); `session.md`: Input states (row 11: a probe.tsv that fails); `session.md`: The acceptance tests of M2b (The probe and admission: a probe that passes, one with a policy path, one that fails for each reason; the skip was in it until task `T-nxe4`); `session.md`: The acceptance tests of M2b (The skip, a refused probe and the early end: the skip, a probe's refused start, which was in A refused start until task `T-nxe4`, and the early end); `session.md`: The acceptance tests of M2b (The routing register: the copy at the step probe); `session.md`: The acceptance tests of M2b (The step `probe`); `run.md`: The command (probe in the enum of the block run-steps, with its Go schema) | In CI with no secret, a restart's step `probe` writes the rows of a scripted harness's probe, with the same bytes on a repeat. | `later-p2-sessions-ledger` | REQ-013, NFR-001 | F-0003#52 | unit (each reason of a failed probe; the copy of the routing register); e2e (the binary, a local fake forge, a scripted harness); the Go schema of `run-steps` with `probe` | large | 800 | 36b | 1 |
| 39a | `T-4tjy` | [#169](https://github.com/pharzam/layup/issues/169) | The review of the release of M2b: `session.md`: REQ-015 and REQ-017 — The review of the release of M2b; `session.md`: The acceptance tests of M2b (The review of the release) | The release review of `M2b` is a recorded review of the code for `REQ-015` and `REQ-017`, by a reviewer of a model that wrote none of it. | `later-p2-sessions-ledger` | REQ-015, REQ-017 | F-0003#52 (the review keeps the out-of-scope facts F-0003#53 and F-0003#55 out) | the review, with `runs/T-efmy/release-check.sh` adapted, its evidence under `runs/` | small | 250 | 37b, 38 | 1 |
| 39b | `T-x7cs` | [#170](https://github.com/pharzam/layup/issues/170) | The demo of M2b (uat): `session.md`: The acceptance tests of M2b (The demo) | On real harnesses, after a probe of each of two registered harnesses, one developer session's commit lands on `task/<task>/1` with its telemetry row. | `later-p2-sessions-ledger`, `gov-operator-setup-o112` | REQ-013, REQ-003, REQ-005, REQ-011, NFR-001 | F-0003#52, F-0003#50 | uat, with its evidence under `runs/` | small | 250 | 39a | 1 |

Rows 33b and 34a both change the one call of `internal/session` and neither
is after the other, so the second to merge takes the first's change by a merge
of `origin/main`.

Row 39a reviews the release before the demo, as row 19 did before row 20 in
phase 1, and takes #148. The release of `M2b` is the code of `main` that the
demo runs (the non-test Go files of `cmd/` and `internal/`, `go.mod`, and the
files that the binary embeds), at `9490de9`, the merge of row 38, the last
build row; it holds the code of phase 1 and `M2a` too (task `T-4tjy`, #169).
A task of `M2b` that changes the release code after row 39a checks its own
diff for `REQ-015` and `REQ-017` in its review round, and runs the adapted
`release-check.sh` with the full commit ID of its head. Row 39b needs the
Operator's inputs, which its issue lists.

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
| NFR-005 | each verdict reproducible (rows 5 and 7); the engine checks start no model process and open no network package (row 2) | no package of an engine check depends on the package that starts a harness (`M2b`); the smart-if client and `decisions.tsv` (`M3a`) |
| NFR-006 | LAYUP's own pin and check `pin` (already green); the target's pin written once and checked (rows 7 and 9) | — |
| NFR-007 | the standard library only, by `go list -deps`, and `git` only as a program (row 2) | — |
| REQ-015, REQ-017 | the recorded code review of the phase-1 release (row 19) | the review of each later release, by the specification task of each milestone |
| REQ-016 | the intent decisions of the first pilot are recorded as the idea owner's (row 20) | the pilot of phase 4 (`M4c`) |
| REQ-018 | the baseline rules of the pilot's target equal the pinned baseline's, except the values and the replaced files that the setup record names with their source (row 20; K42) | the pilot of phase 4 (`M4c`) |

## The hosts of the other items

Each inventory item that is not in a row of a task table has one host here; `later-p2-run-start` and `gov-operator-setup-o112` keep a row, which points to their rows and gives their other hosts.

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
| `gov-pilot-numbers` | the ADR that supersedes ADR-0012 | The count is the first step of the ADR task, by this plan's choice. Row 20's plan review did not pre-register its rule; task `T-gd8q` wrote the rule after the numbers, and the Operator accepted it (O-159, [`runs/T-gd8q/numbers.md`](../../runs/T-gd8q/numbers.md)). |
| `gov-supersede-adr-0012` | the ADR that supersedes ADR-0012 | A separate issue that row 20 opens (O-121). |
| `gov-operator-setup-o112` | row 27 (`M2a`), row 39b (`M2b`), `M2e`, `M3d` | The App's private key on the host (mode 0600) and the App installed on the target: Start of `M2a`, asked in the issue of row 27. The harness credentials and the host registers of `M2b`: the demo of `M2b`, asked in the issue of row 39b. The commit-statuses permission: `M2e` (`later-gate-status`, the four `layup/` checks). `layup-watch`: `M3d`; until then Start records `watch` as `not-confirmed`. |
| `gov-operator-ci-protection` | out | Phase 1 adds no CI job; a task that changes `ci.yml` names the Operator's step in its own plan (K28). |
| `gov-issue-68` | `M2e` | The handoff records of the task loop. |
| `gov-issue-49` | out | [R11](../issue-workflow.md#r11--single-goal-issues) now holds the answer ("A task of one artifact", [ADR-0026](../adr/0026-keep-the-bootstrap-review-rules-as-the-standing-gate.md)); the Operator closes #49 or keeps it open. |
| `later-psb-meaning-review` | `M2c` | The review of meaning needs a role session. |
| `later-psb-intake-batch` | `M2c` | The one Intake batch needs `layup run`. |
| `later-psb-early-question-share` | `M4b` | A measure of `layup report`. |
| `later-gate-status` | `M2e` | The status `layup/gates` needs `layup run` and the forge. |
| `later-gate-rule-batch` | `M2f` | Rule batches and activation. |
| `later-p2-run-start` | rows 21 to 27 | [The tasks of M2a](#the-tasks-of-m2a). |
| `later-p2-sessions-ledger` | rows 28 to 39b | [The tasks of M2b](#the-tasks-of-m2b). |
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
| K38 | `T-zck8` (`M2a`) | `M2b`, `M2c`, `M2e`, `M2g`, `M3d` | Settled by one rule: a record is specified in the milestone whose code first writes it ([`records.md`](../spec/records.md#the-layout-of-the-records-branch)). `questions.tsv` and `budget.tsv` go to `M2c` (Intake writes them; the cap of milestone 0 is `intake.cap` of `start.tsv`), `routing.tsv` and the second admitted harness to `M2b`, the frozen test list to `M2e` (frozen at the task's first valid handoff to the verifier, `architecture.md` §8; read by `M2g`), the dead-man job and its record to `M3d`; `M2a` only reads its first notice. |
| K39 | row 2 | `M2a`, `M3a` | Phase 1: no package imports a network package; `M2a` names the forge adapter and `M3a` names `internal/smartif`. |
| K40 | `T-zck8` (`M2a`): the App key; `T-ywk7` (`M2b`): the harness credential | `M2b` | The App ID, the key file and the URLs are the forge register that `--host DIR` names (decided in [`run.md`](../spec/run.md), [`records.md`](../spec/records.md#nfr-001--the-records-of-start)). The harness credential is a file that the harness register's row names, given to the session as a variable or a copied file, so `layup` reads no environment variable for it ([`session.md`](../spec/session.md#the-environment-and-the-harness-credential)). |
| K41 | O-163 (#123) | `M2d` | `M2a` specifies an empty repository only (`--new`); the adoption of a target that phase 1 set up goes to `M2d`, where `layup run` drives the Scaffold. |
| K42 | row 13 | rows 16, 20 | Found by the self-check: the setup record names a source for each value, but no row for a kept file that the prose step replaces (O-123). The prose step writes one record row for each replaced file, so `REQ-018` has its evidence. |

Each task also takes the open questions and the duplicate resolutions of its own
items in the inventory as inputs of its plan.

### The findings of the first pilot that are not fixed

The first pilot (row 20) left these findings of
[`runs/T-evad/findings.md`](../../runs/T-evad/findings.md) on its path, not fixed.
The specification task of `M2a` (`T-zck8`, #123) moved the nine that it hosted to
the milestone whose code each one changes. The ADR that ended bootstrap mode
([ADR-0026](../adr/0026-keep-the-bootstrap-review-rules-as-the-standing-gate.md),
task `T-gd8q`) gives each one host. F-15 is settled by that ADR, and F-1 and F-11
are known limits by the Operator's rulings.

| Finding | Host | Note |
| ------- | ---- | ---- |
| F-2, F-31, F-35 | `M2d` | A second setup run on a target: a changed input of a done step, the immutable records, the names of the branches. `layup run` drives the setup in the Scaffold, and adopts a phase-1 target there (K41, O-163). |
| F-6 | `M2d` | `layup setup verify` before S12; the Scaffold reads the verify table. |
| F-13, F-19 | `M2d` | S13: `commands.sh` is not fail-fast; the ruleset parameters that S13 does not set. The Scaffold re-specifies the files and the printed commands of S13 (the Operator still applies the rulesets with the Operator's own login, `architecture.md` §5, Scaffold 6); the fail-fast rule and the two parameters belong to that rewrite. |
| F-16 | `M2d` | A stale line that the setup does not flag; the adaptation step of the Scaffold. |
| F-21 | `M2c` | Rule G1 of `layup psb check` reads only `technology stack:`; the gap check of Intake. |
| F-22 | `M2e` | The gate jobs read the base's gate files (the Operator's direction). |
| F-23 | `M2f` | The catalog gates aligned with the target's own gate script. |
| F-37 | `M2e`, `M2g` | A tested tool for the steps that land a target task: the merge at the head SHA and the close-out commit (`M2e`), and the merge order of the build tasks (`M2g`). |

## The open issues

| Issue | Disposition |
| ----- | ----------- |
| [#15](https://github.com/pharzam/layup/issues/15) | Settled by `T-55n2` (its lines in `docs/guardrails.md`). |
| [#21](https://github.com/pharzam/layup/issues/21) | Row 10 (`T-8vpw`), with check `markers`; cap 2. |
| [#24](https://github.com/pharzam/layup/issues/24) | Checked and closed by `T-55n2`, which also fixes its one open sentence in `docs/ci/README.md`. |
| [#34](https://github.com/pharzam/layup/issues/34) | Row 15 (`T-d6q5`), with the rule-path register. |
| [#48](https://github.com/pharzam/layup/issues/48) | Row 12 (`T-9t1q`), with check `facts`; cap 2. |
| [#49](https://github.com/pharzam/layup/issues/49) | Answered by [R11](../issue-workflow.md#r11--single-goal-issues) ("A task of one artifact", [ADR-0026](../adr/0026-keep-the-bootstrap-review-rules-as-the-standing-gate.md)); the Operator closes it or keeps it open. |
| [#61](https://github.com/pharzam/layup/issues/61) | Row 12 (`T-9t1q`), with check `facts`; cap 2. |
| [#68](https://github.com/pharzam/layup/issues/68) | Milestone `M2e`, with the handoff records. |
| [#148](https://github.com/pharzam/layup/issues/148) | Row 39a (`T-4tjy`, O-187 of #152): the review of the `M2b` release, before the demo, reads the whole release, the code of `M2a` included ([`session.md`](../spec/session.md#req-015-and-req-017--the-review-of-the-release-of-m2b)). |

## Known limits of this plan

- **The inventory is not complete.** A fresh review of a task's sections can find
  more defects; the task settles them in its own pull request.
- **The expected lines are estimates.** Each task's plan review sets its budget;
  a task that grows past its budget goes back to its issue as a new slice or a
  child issue (R12).
- **The milestones of phases 2 to 4 are coarse.** Each one is fixed by its own
  specification task, which can change its scope and its order.
