# The task loop: handoffs, the four `layup/` statuses, Shape and the first bet

The specification of milestone `M2e` of [the plan](../plan/README.md#milestones):
what `layup run` adds after the Scaffold, so that each change to a target is a
task that goes through a plan, a plan review, a developer, a verifier and a
close-out, each handoff is checked by code against a table, the forge shows four
`layup/` statuses for each head, and the first bet approves the milestone. Its
requirements are `REQ-005` (role handoffs), the status part of `REQ-007`,
`REQ-012` (the rest of the specification synthesis: Shape, the bet, the full
`layup spec check`), the basic rule of `REQ-003` (`layup/rules`), `NFR-004` and
`NFR-001`; the design is [`architecture.md`](../architecture.md) §6 (the stack
gates and rule protection), §7 (the specification), §8 (the phase loop and its
handoffs), §9 (squads and routing) and §10 (the parameters), with ADR-0016,
ADR-0017, ADR-0019 and ADR-0020. Task `T-xs4y` (#333) wrote it, after a
public-solution search ([`runs/T-xs4y/search/`](../../runs/T-xs4y/search/)),
whose result is that no public solution can be its base; the decisions are
comment 6191483681 of #333. The records are in
[`records.md`](records.md#nfr-001--the-records-of-the-task-loop). The sessions
are those of [`session.md`](session.md), the forge calls those of
[`forge.md`](forge.md), the gate that `layup/gates` reports is
[`gate.md`](gate.md).

**What the search gives this file.** GitHub's commit-status API gives a status
for a head SHA with a `context` and one of four states; a required check can be
pinned to the App as its source; the merge call takes the head SHA and refuses a
head that moved; and a required check may be met by a `neutral` or `skipped`
state, which is why this file counts only `success` (searcher A); at most 1,000
statuses are kept for one SHA and context (searcher B). The best public handoff
tools validate a handoff record against a schema and a role-to-role rule, as
this file does, but none keeps the table of transitions where the agent cannot
write it, and none gives an artifact that exists only as an issue comment a place
and a hash (both searchers). Nothing public filters architecture options by code,
checks a brief to the idea owner, or fails a requirement that has no criterion
by four rules together.

## REQ-005 — The task loop

Requirement `REQ-005`: information that passes between role agents has a form
that a machine can validate (`F-0003#45`, `#59`). **Every pull request is a
task** (`architecture.md` §8): the specification, the target's architecture, each rule batch, each rendered
record, the milestone plan and each build task. Each task has a row in
`tasks.tsv` ([the records](records.md#nfr-001--the-records-of-the-task-loop)),
written when `layup run` opens its issue, and goes through the steps below. The
steps are the eight of `architecture.md` §8, "The task loop"; this table gives
their actor, their input and the check of code at each, and the section that
holds the rule.

| # | Step | Actor | Input | Checked by code | Rule in |
| - | ---- | ----- | ----- | --------------- | ------- |
| 1 | the issue | `layup run` | the requirement IDs and the Definition of Done of the task's row | the row of `tasks.tsv` names each requirement ID and each predecessor | [The task register](#the-task-register) |
| 2 | the plan and its review | a plan session; a plan-review session on another harness | the issue, the specification sections of the requirement IDs | the plan has its four fixed fields and lists its tests by ID; the review has the verdict field `approve` or `approve-with-conditions` | [The transition table](#the-transition-table) |
| 3 | the developer | a developer session | the plan, the findings of the last attempt | the typed result ([`session.md`](session.md#req-005--the-result-of-a-session)); the tests of the plan fail at the base and pass at the head | [The checks of a handoff](#the-checks-of-a-handoff) |
| 4 | the draft pull request | `layup run` | the branch of the attempt | no rule-path change ([`architecture.md`](../architecture.md) §6, prevention layer 2); the title in the target's form; the body links the issue and the records commit the task started from | [The pull request](#the-pull-request) |
| 5 | the checks | the target's CI, `layup gate`, `layup spec check`, `layup/rules` | each new head | the four statuses | [The four `layup/` statuses](#req-007--the-four-layup-statuses) |
| 6 | the verifier | a verifier session on a harness that is not an author | the head, read-only | the review record in the target's form, with a `file:line` for each finding | [The independence of the reviewers](#the-independence-of-the-reviewers) |
| 7 | the close-out | a close-out session | the verdict, the ledger | the commit changes only the task file and the completed log | [The merge at the head SHA and the close-out commit](#the-merge-at-the-head-sha-and-the-close-out-commit) |
| 8 | the merge | `layup run` | the head, the four statuses | the merge at the head SHA | [The ready and merge rule](#the-ready-and-merge-rule) |

### The task register

`tasks.tsv` has one row for each task, written by `layup
run` and never changed (a change is a new row of a later `tasks.tsv` commit that
Git history keeps, as for each record, ADR-0014): the task ID, the issue, the
kind (`specification`, `architecture`, `batch`, `rendered`, `plan`, `build`), the requirement IDs
it serves, its size class, its predecessors, and its pull request when it has
one. **Decided here:** a task starts only when each predecessor is a task of the
register and has merged, and the predecessors form no cycle; code checks both
when it writes the row, and refuses a row that breaks one. At most
`build.parallel` build tasks run at once (a parameter, `architecture.md` §10);
the merge order of build tasks is `M2g`'s.

### The pull request

**Decided here.** `layup run` opens the pull request as a draft, with a title
of the form `<type>: <task ID> <description>` (the target's own convention,
read from the setup record), a body that links the issue (`Closes #N`, or
`Refs #N` for a task that does not close it) and links the records commit that
the task started from, and requests no review. A draft is only visible; a
pull request "reaches human review" when a review is requested or it is marked
ready (the reading of `F-0003#58`, which the approval brief asks the Operator to
confirm). The review record of the verifier is posted on the **issue**, and the
body of the pull request is edited to link it, which makes the target's own
`review-record` check run again. **Reason:** the target's checks parse the
issue's comments and the body, so the loop writes them in the forms that those
checks read.

### The attempts

A new attempt starts from the base commit, with the findings of the last attempt
and any diagnosis in its prompt file and not its diff, so that its commits have
one author harness (`architecture.md` §8, "Attempts"). **Decided here:**

- A failure of a check at step 5 goes back to the developer as a finding and
  starts a new attempt (the two expected failures of `architecture.md` §8 step 5, the
  rule batch before its approval and the specification that waits for its bet, do
  not: [The rendered commit and the confirmed inventory](#the-rendered-commit-and-the-confirmed-inventory)): an `attempt` row of `tasks/<task>/events.tsv` with the
  base SHA, then the sessions of that attempt ([`session.md`](session.md)).
- **Decided here** (task `T-6zj3`, #361, row 66e): the findings of a gate failure are
  the rows of the table of `layup gate` ([`gate.md`](gate.md): `gate.Table`) whose
  result is neither `pass` nor `clear`, in the order of the manifest, each a finding
  with the rule `gate:<kind>`, the row's reason as its text, and no file and no line (a
  row has none); a `not-active` row is a finding, as it never counts as a pass
  (`NFR-004`). `internal/run`'s `NextAttempt` starts the attempt (`StartAttempt`: the
  next number, the base that the caller gives) only when a row failed, and gives the
  findings back (`ErrNoFailure` otherwise, and nothing is written); the prompt file of the
  attempt carries them and no diff ([The tier and the context](#the-tier-and-the-context-of-a-task)).
  **Later wiring:** the call that notices the failure at step 5 and the rows of
  `gates.tsv` that keep the findings across a restart are #362's and #367's; until then
  the caller gives the table, and the findings are the caller's to give again.
- Once the task's test list is frozen ([The frozen test list](#the-frozen-test-list)),
  code puts the frozen test files at the start of each new attempt, so that their
  hashes stay; the harness that wrote them is an author of every later attempt.
  **Decided here** (task `T-6zj3`): `FrozenFiles` gives, for each listed test, the bytes
  whose SHA-256 is the frozen one, found among the heads that were bound to the task's
  sessions, the newest first (a test whose bytes no bound head holds is an error, never a
  guess); `TaskSession` writes them in the session's clone, right after the clone is made,
  and commits them as one commit `layup: the frozen tests` (so the head's tree holds them,
  as the base lacks the tests that an earlier attempt wrote), its author the session's own
  identity (`layup session <session>`) at the run's time; a file that is already as frozen
  at the base needs no commit. The file's mode is not kept: a known limit. A symbolic link in the
  checkout, in a parent of a path or as the file itself, whatever it points to, is **not followed**:
  the call refuses with an error naming the path, after checking the whole list and before writing
  any file, so nothing is written inside or outside the checkout (`guardrails.md` §2, a path that is
  clean by its spelling).
- When the Operator passes the last diff to the new attempt (a stall answer,
  `M3d`), the harness that wrote it is an author of the new attempt. No record names
  that diff before `M3d`, so the authors of [the independence](#the-independence-of-the-reviewers)
  do not read it in `M2e`, and `M3d` adds it.
- A question that ends an attempt (status `needs_context`) is a wait for the
  owner role, not a round without progress ([The checks of a
  handoff](#the-checks-of-a-handoff)).
- **Known limit:** no bound stops a loop of failures of step 5 before `M3d` (the
  stall triggers) but the budget cap of `M2c` (`intake.md`, the budget of Intake)
  and the cycle cap of the target for the review rounds of step 6; an attempt
  that ends with no commit is an `attempt` row with no `push` row, and the next
  attempt starts as above.

### The independence of the reviewers

The authors of a change decide which harness may review it
(`architecture.md` §9). **Decided here:** code computes the authors from the
records and the ledger, never from a field that an agent writes:

- the **authors for a plan review** are the harnesses of the task's plan
  sessions (`sessions.tsv`);
- the **authors for a verification** are every harness whose session is bound,
  in the ledger, to a commit in the diff from the base, every harness whose diff
  went into an attempt's prompt, and the harness that wrote the task's frozen
  tests (Invariant 9).

**Decided here** (task `T-6zj3`, #361, row 66e; `internal/run`'s `Authors`): a session is
bound, in the ledger, to a commit of the diff from the base when the event `bound` of its
task holds a SHA that is an ancestor of the head and not an ancestor of the base (the two
calls `IsAncestor` of the run's clone, into which the end of each session fetched its head);
its harness is that of its row of `sessions.tsv`, a row that `layup run` wrote, never a
field of the session's result. A session that pushed and was never bound is no author, nor
is one bound to a commit of an earlier attempt or of the base. The authors of a
verification are those harnesses and the `harness` of each row of the frozen list
(`tests.tsv`); the authors of a plan review are the harnesses of the task's sessions whose
role is the plan role of the step table (`plan-of-a-task`), and as `sessions.tsv` has a role
and no step, a session of that role that answered a question is counted too: a known limit,
on the safe side (fewer admitted reviewers). The result is sorted, each harness once. The
third source (a diff that an Operator's stall answer passed to an attempt) has no record
before `M3d`, which adds it.

Code admits a pair of harness and model for a plan review or a verification only
when its harness passed the probe ([`session.md`](session.md#admission)), its
model is not on the target's "not used" list, and its harness is not an author.
When no admitted harness is outside the authors, the verification is
`not-active`: `layup/verify` is posted `failure` with the description
`not-active: no harness outside the authors`, and the pull request does not merge
(`F-0001#5`, Invariant 5). **Reason:** the forge's four states have no
`not-active`, and `failure` is the one that stays a refusal until a harness
exists; `pending` is kept for a verifier that waits for its turn (below). **Decided here** (task
`T-6zj3`): `PairFor` gives the first pair of `routing.tsv` for the role and the tier of the
step (`plan-review` and `verification` of the step table: `qa-engineer`, `reasoning`) that is
admitted by the rule of [Admission](session.md#admission) and whose harness is none of the
authors, and `ReviewerPair` gives it with the authors read; with none the error holds both
`ErrNoHarnessOutsideAuthors`, whose text `no harness outside the authors` is the description
after `not-active: `, and `route.ErrNoPair`. The status that posts it is row 66f's.

### The tier and the context of a task

**The tier** of an implementing task starts from its plan, which has four fixed
fields: `size` (`small` or `large`), `open questions`, `new interfaces` and `new
dependencies`. The reasoning tier when `size` is `large` or any of the three
lists is not empty; the execution tier otherwise (`architecture.md` §9). The
progress position of a task, uphill while it has open unknowns and before its
test list is frozen, moves it to the reasoning tier, and downhill to the
execution tier (`architecture.md` §11). **Decided here:** the tier is read when
the plan is accepted (step 2) and again at each attempt, from the plan and the
open unknowns that the records give; the result is a column of the routing
request, not a new record.

**Decided here** (task `T-6zj3`, #361, row 66e), the reading of the tier:

- **The plan is a typed value.** `Plan` holds the four fixed fields (`size`, and the three
  lists); the text form of a plan, the reader of its comment and the check `plan-fields` are
  #387's, so the caller gives the value. `PlanTier` is the tier of the plan alone: reasoning
  when `size` is not `small` (a size that is neither `small` nor `large` counts as large: the
  field has no other value, and a typo must not choose the cheaper tier) or any list is not
  empty, else execution.
- **The position moves it at an attempt.** `TierOf(plan, position)` is reasoning when the task
  is uphill and execution when it is downhill; with no position the plan alone decides, as at
  its acceptance (step 2). By this reading the plan's class decides at the acceptance and the
  position at each attempt; both are the text of `architecture.md` §9. A task is downhill when
  its test list is frozen and it has no open unknown, uphill otherwise (`Position`).
- **The open unknowns that the records give in `M2e`** are the question rows of the last attempt
  that ended `needs_context`, until an answer that started after it is accepted by the next
  attempt (`AnswerAccepted`, with the rule of the owner role in [The checks of a
  handoff](#the-checks-of-a-handoff), "A question during the work"). The open material findings (their
  form is #387's) and the open escalations (`M3`) are added when they exist.
- **`by-class` and `owner-of-its-kind`** ([`session.md`](session.md#the-step-table-and-the-kinds-of-a-steps-context))
  are resolved by `ResolveStep`: `by-class` becomes the tier of `TierOf`, `owner-of-its-kind`
  becomes the routing role of the owner map's role. The key of `routing.tsv` is the role ID of
  the step table (`systems-architect`), the owner map's roles are names (`Systems Architect`),
  so the one rule is lower case with each run of spaces a hyphen (`RoutingRole`): the three
  default owners that the step table has as roles come out as exactly those IDs, and `Domain
  Expert`, which no step has, becomes `domain-expert`, for which the Operator's register has a
  pair or the start is refused (`pair`). A human has no routing role.
- **The prompt file** is the list of links of the five kinds of the step table in their order
  (`PromptFile`), each item a path with the SHA-256 of the bytes that code read, a file once
  (the first kind that names it keeps it): `issue-and-plan` is `tasks.tsv` (the issue's own text
  is no record in `M2e`: a known limit; its row has the number, the requirement IDs and the
  class) and the kept copy of the last valid handoff with a comment from a session of the plan
  role; `requirements` are the files of `docs/prd/` that have a table row, and the files of
  `docs/spec/` that have a heading, that holds a requirement ID of the task's row as a whole
  word (`REQ-0050` is not `REQ-005`); `decision-records` are the `docs/adr/<NNNN>-*.md` of each
  `ADR-NNNN` that a section of those headings names (the heading's text down to the next heading
  of its level or above), and a decision record that a section names and the target lacks is an
  error; `handoffs-and-answers` are the copies of each valid handoff with a comment, and the
  result files that hold an `answer` row, then, as text, the findings of the last attempt and no
  diff; `rule-files` are the files of the target at the base that the rule-path register of the
  records matches. **Decided here** (task `T-c6sx`, #416, row 66s): the frozen list `tests.tsv` of the task, when it exists,
  joins `issue-and-plan` after the plan; the revised draft and its facts join `requirements` and the
  options tables join `decision-records` for a Shape task, with the survivors as text under them; the
  ledger joins `handoffs-and-answers` for the close-out ([The instruction of a task
  step](#the-instruction-of-a-task-step) holds the table); the order of the items that exist does not change. A target item is read at the base; a records item is copied under the
  session's `input/` by its file name (two records of one name are refused), and the text names
  it at `<session directory>/input/<name>`. The SHA-256 that the ledger keeps, as the paragraph below says, is
  the one in the prompt file itself, which the start row's `records` commit rebuilds. **Later wiring:**
  the call of `PromptFile` and of `TierOf` at each step is the step `loop` (#367).

**The context of a session.** The step table has a column of record kinds for
each step ([`session.md`](session.md#the-step-table-and-the-kinds-of-a-steps-context) holds the
table and `M2b` the prompt as the caller's text). **Decided here:** from `M2e` the
prompt file of a step of a task holds, in this order, the task's issue and plan, the rows and the
specification sections of the task's requirement IDs, the decision records that
those sections name, the task's handoffs and the answers to its questions, and
the target's rule files; code builds the list from the records and writes the
links, so a session reads the files and not a summary. Each item is a path of the
target or of the records branch with its SHA-256, which the ledger keeps.

### The instruction of a task step

**Decided here** (task `T-c6sx`, #416, row 66s; the plan review is the comment 6249930540 of #416, whose conditions 2 to 7 are taken as written; the Operator set the goal count at 2, comment 6249938546): the prompt file of a step of a task is the list of links above and, after it, the step's fixed instruction. The lists give the harness its inputs; the instruction gives it the work, the outputs and the form of its files, so a model with no other context can write a result that the loop accepts.

- **The frame.** `Text` keeps the header line, the sentence that tells the harness to read, and the five lists, in their order, each existing item in its place relative to the others; the new items below join their own kind. After the lists, code writes a blank line, the line `The work of this step:`, the block of the step with its placeholders filled, a blank line and the block `instruction-result` filled. A target item or a records item that is new makes the bytes of a list longer, so the golden files, not a byte prefix, are the oracle of the composition.
- **The home of each text.** The eight blocks below are the text of record, with one home each. `internal/run` keeps each as a constant and a test compares it with its block byte for byte (as it does for `probe-prompt`); the golden files of the prompt of each step compare the composed bytes, and no flag rewrites a golden file.
- **The placeholders** are filled by code and never by a session; a prompt that keeps a `{` placeholder is a defect. `{task}`, `{kind}` and `{issue}` are the row of the task in `tasks.tsv`. `{result}` is the absolute `<session directory>/result/result.tsv`, and `{comment}` the absolute `<session directory>/result/<file>`, the file being `plan.md`, `plan-review.md` or `review-record.md` (the constants `PlanFile`, `PlanReviewFile` and `ReviewRecordFile`, which row 66r's driver passes as the file of the handoff). `{head}` is the base of the session's start row, the head that the item `nothing-material` reads. `{from}` is the base of the task's last `attempt` event. `{labels}` are the keys of the default owner map, sorted, joined by `, `. `{record}` and `{facts}` (the specification only) are the highest `F-NNNN` of `docs/facts/` at the base and the next number after it.
- **The steps.** A step has an instruction exactly when this table has a row; `shape-activation-batch` and `plan-of-a-milestone` are the milestones M2f and M2g, and `answer-to-a-question`, its posting and its tier are row 66k's (the boundary is in the plan): `PromptFile` refuses the three by name, before every read, with `the step <step> has no instruction in M2e`.

| Step | Role | Kind of the task | Writes in the repository | Comment file, read by | Required outputs for the result |
| ---- | ---- | ---------------- | ------------------------ | --------------------- | ------------------------------- |
| `plan-of-a-task` | `software-engineer` | any | nothing, no commit | `plan.md`, `ParsePlanComment` | one `comment` row |
| `plan-review` | `qa-engineer` | any | nothing, no commit | `plan-review.md`, `ParsePlanReview` | one `comment` row |
| `specification` | `product-owner` | `specification` | `docs/prd/`, `docs/spec/`, the task file, commits | none | the task file and each file written |
| `shape-architecture` | `systems-architect` | `architecture` | `docs/spec/`, `docs/adr/`, the task file, commits | none | the task file and each file written |
| `implement` | `software-developer` | `build` | tests and change, the task file, commits | none | the task file and each file written |
| `verification` | `qa-engineer` | any | nothing, no commit | `review-record.md`, `ParseReviewRecord` and `NothingMaterialOK` | one `comment` row |
| `close-out` | `software-developer` | any | one commit of the task file and `completed.md` | none | those two files |

The step `specification` takes a task of kind `specification`, `shape-architecture` one of kind `architecture` and `implement` one of kind `build`; any other step takes any kind, and a pair that disagrees is an error that names the step, the task and the kind. A comment file is written in the pinned form of [The items a caller computes](#the-items-a-caller-computes) (the target's own forms, never LAYUP's `| Field | Value |` table), and the instruction shows it with its placeholders so that a test fills them and runs the real parser on it. The shared text names every rule of the block `result` that a model must meet: the actual tab, UTF-8 with no byte-order mark, the line feed at the end of each line, no blank row, the empty cell, the numbering without gap, a status row with a non-empty reason, a hash of 64 lowercase hexadecimal characters, a question that names a record, and the artifact rows for the required outputs only.

```text instruction-plan-of-a-task
Plan the task {task} (kind {kind}, issue {issue}) from the files listed above: the issue row, the requirements, the decision records, the handoffs and findings, and the rule files. Change no file of the repository and make no commit.
Write the plan as the file {comment}, in this form, each line starting at the first column and ended by a line feed:
size: <small or large>
open questions: <one line, or —>
new interfaces: <one line, or —>
new dependencies: <one line, or —>
tests:
- <id> <path>
The word after size: is exactly small or large. The character — is for the three list fields only, and a field with nothing after the colon and the space is wrong. The line tests: stands alone, once. Under it, write one line "- <id> <path>" for each test of the plan: exactly two words, each ID unique in the list, each path relative to the repository, clean and inside it (no space, no .., no leading /). A task of kind specification or architecture lists no test: write no line under tests:. A task of kind build lists at least one test.
In the result file write the row "comment<TAB>1<TAB>result/plan.md<TAB><the SHA-256 of {comment}><TAB>—".
```

```text instruction-plan-review
Review the plan of the task {task} (kind {kind}, issue {issue}), the newest plan among the handoffs listed above, against the issue row, the requirements, the decision records and the rule files: its coverage of the requirements, its list of tests, and its fit with the decision records and the rule files. Change no file of the repository and make no commit.
Write the review as the file {comment}, in this form, the first line starting at the first column and ended by a line feed:
verdict: <approve, approve-with-conditions or reject>
The line verdict: is written once, and its value is exactly one of the three words approve, approve-with-conditions and reject. After it, write the conditions or the reasons as plain lines that do not start with "verdict: ". This is the form the target reads; it is not the table "| Field | Value |" of LAYUP's own reviews, which you must not write.
In the result file write the row "comment<TAB>1<TAB>result/plan-review.md<TAB><the SHA-256 of {comment}><TAB>—".
```

```text instruction-specification
Write the specification of the task {task} (kind {kind}, issue {issue}) from the revised draft and its facts listed above: spec/draft.tsv, spec/spans.tsv, start/problem-statement.md and intake/answers.tsv. Commit your work in this repository in one or more commits; push nothing.
The IDs. The row R-NNN of kind requirement in spec/draft.tsv is the requirement REQ-NNN, and the row R-NNN of kind non-functional is NFR-NNN, with the same three digits. A row of kind success-criterion is no requirement and has no such ID. Keep the IDs as the draft listed above numbers them; never renumber one.
In docs/prd/ write the product requirements document of the target, with one row for each requirement ID, in these forms (every line starts at the first column):

| REQ | Statement | MoSCoW | Phase | Facts |
| --- | --------- | ------ | ----- | ----- |
| REQ-<NNN> | <the requirement> |  |  | F-{facts}#<n> F-{record} |

| Requirement | Acceptance criterion |
| ----------- | -------------------- |
| `REQ-<NNN>` | <the acceptance criterion> |

## 9. Success metrics
| Metric | Fact |
| ------ | ---- |
| <the success criterion> | F-{facts}#<n> |

| REQ | Facts | Guardrail | ADR | Task | Test |
| --- | ----- | --------- | --- | ---- | ---- |
| REQ-<NNN> | F-{facts}#<n> F-{record} | — | — | — | — |

Write NFR-<NNN> in place of REQ-<NNN> for a row of kind non-functional, and keep the backticks around the ID in the acceptance table. Leave the MoSCoW and Phase cells empty (not —): code fills them at the bet. In a Facts cell, cite each fact that the row covers: F-{facts}#<n> for the fact n of spec/spans.tsv, and Q-<NNN> for a fact of intake/answers.tsv; always cite F-{record}, the existing record, too. Write one acceptance row and one matrix row for each requirement ID, and one success row for each success-criterion row of the draft, citing the facts it covers.
In docs/spec/ write, for each requirement ID, a section whose heading holds the ID as a whole word, for example "## REQ-<NNN> — <title>", with a non-empty text under it.
Write docs/tasks/{task}.md, the task's file, in the form of the target's docs/tasks/. Write nothing under docs/facts/: code renders it at the first bet.
The required outputs for the result file are docs/tasks/{task}.md and each file of docs/prd/ and docs/spec/ that you wrote.
```

```text instruction-shape-architecture
Write the architecture of the target for the task {task} (kind {kind}, issue {issue}) from the options tables of its decisions and the survivors of the filter, both listed above under decision-records, with the revised draft and its facts. Commit your work in this repository in one or more commits; push nothing.
For each decision that has one survivor, write that option as decided. For each decision that has more than one survivor, write a section with the heading "Decision pending the bet" that lists its survivors and decides nothing. Never choose among survivors, and never write an option that is not a survivor.
Write the architecture under docs/spec/ and the decision records under docs/adr/, in the form of the rule files listed above, and write docs/tasks/{task}.md, the task's file, in the form of the target's docs/tasks/.
The required outputs for the result file are docs/tasks/{task}.md and each file of docs/spec/ and docs/adr/ that you wrote.
```

```text instruction-implement
Implement the task {task} (kind {kind}, issue {issue}) as its plan says: the plan and its review are among the handoffs listed above. For each test the plan lists, write the failing test first, watch it fail for the right reason, then write the change. Answer each finding listed under the findings of the last attempt. When tests.tsv is listed above, change none of the test files it names. Keep docs/tasks/{task}.md, the task's file, current. Commit your work in this repository in one or more commits; push nothing.
The required outputs for the result file are docs/tasks/{task}.md and each file you wrote or changed.
```

```text instruction-verification
Verify the change of the task {task} (kind {kind}, issue {issue}) from {from} to {head}, the commit your repository is at: review the changes from {from} to {head} against the plan, the requirements, the decision records, the rule files and the frozen list tests.tsv, as listed above. Change no file of the repository and make no commit.
Write the review record as the file {comment}, in this form, each line starting at the first column and ended by a line feed:
head: {head}
verdict: <material, nothing material in scope or not mergeable, findings recorded>
findings:
- <path>:<line> <rule>: <text>
The line head: is written as shown, once. The value of verdict: is exactly one of the three strings material, nothing material in scope and not mergeable, findings recorded. The line findings: stands alone, once. Under it write one line for each material finding: <path> is relative to the repository, clean and inside it (no space, no .., no leading /), <line> is a positive integer, <rule> is one word with no colon, and <text> is not empty. For nothing material in scope write no line under findings:. This is the form the target reads; it is not the table "| Field | Value |" of LAYUP's own reviews, which you must not write.
In the result file write the row "comment<TAB>1<TAB>result/review-record.md<TAB><the SHA-256 of {comment}><TAB>—".
```

```text instruction-close-out
Close out the task {task} (kind {kind}, issue {issue}). From the verdict among the handoffs listed above and the ledger (sessions.tsv and telemetry.tsv, listed above), write the task's verdict, its line of the completed log and its resource record. Make one commit that changes only docs/tasks/{task}.md and docs/tasks/completed.md; push nothing.
The required outputs for the result file are those two files.
```

```text instruction-result
Write the result file {result}: UTF-8 text with no byte-order mark, one row per line, each line ended by a line feed, the last too, and no blank line. Separate the cells of a row with one tab character each, never with spaces. The first line is "kind<TAB>n<TAB>value<TAB>sha256<TAB>reason". The single character — (U+2014) is the empty cell; no cell holds a tab or a line break. The rows of each kind are numbered 1, 2, 3 and so on, with no gap.
Write exactly one status row, "status<TAB>1<TAB><status><TAB>—<TAB><reason>": <status> is completed when the work of this step is done, needs_context when a question stops it, and blocked or failed otherwise; <reason> is one line and not empty.
For each required output that the work above names, write the row "artifact<TAB><n><TAB><path><TAB><sha256><TAB>—": <path> is relative to the repository, and <sha256> is 64 lowercase hexadecimal characters, the SHA-256 of that file at your last commit. A file that this prompt only asks you to read is no output.
When a question stops the work, write one row for each question, "question<TAB><n><TAB><text><TAB>—<TAB><label> <needs_human> <record>": <label> is one of {labels}; <needs_human> is yes or no; at least one <record> names a file the question concerns, more are separated by one space, and there is no space at either end of the cell.
A comment file is plain Markdown with no code fence. Push nothing. Then end.
```

**The records each step receives.** Every new item is a byte copy under `input/` with its SHA-256, inside one of the five kinds (no sixth kind: `ContextKinds` does not change), after the items that exist; a file is read by its reader, which fails closed.

| Input | Records path | Kind, after | Steps | Reader |
| ----- | ------------ | ----------- | ----- | ------ |
| the frozen list | `tasks/<task>/tests.tsv`, when it exists | `issue-and-plan`, after the plan | all | `ReadTests`; absent is no item |
| the revised draft and its facts | `spec/draft.tsv`, `spec/spans.tsv`, `start/problem-statement.md`, `intake/answers.tsv`, in this order | `requirements`, after the target items | all but `close-out`, for a task of kind `specification` or `architecture` | `ReadDraft`, `ReadSpans`, `ReadAnswers`; the brief as bytes; a missing file is an error |
| the options | `spec/options/<decision>.tsv` for each decision of the caller's map, sorted | `decision-records`, after the ADRs | all but `close-out`, for a task of kind `architecture` | `spec.ReadOptions` |
| the survivors | text: `the survivors of the filter:` and a line `- <decision>: <option> ...` for each decision, under the decision records | `decision-records` | as the options | the caller's map, checked against the tables |
| the ledger | `sessions.tsv` and `telemetry.tsv` | `handoffs-and-answers`, after the answers | `close-out` | `ReadSessions`, `ReadTelemetry`; both tables whole (a known limit: code keeps no table of one task) |

The Intake records go to both Shape kinds because the `covers` column of the draft names span numbers and `Q-NNN` facts, and the constraints of an options table name `R-NNN`, `F-<n>` and `Q-NNN`; neither can be read without the four files. **The options (R10, binds 66h, 66r and 66j).** `PromptIn.Options` maps each decision to the `spec.Result` of its table and is required, and not empty, for a step that takes the options (an architecture task, any step but `close-out`), and refused for any other. *Before the bet* the caller (row 66h's filter, passed by row 66r's driver) gives each decision's result unchanged: every option of the table and every survivor, the survivors listed in the prompt and none chosen. *After the bet* the caller (row 66j opens the task of the chosen option; row 66r passes it) gives the same options and only the option that the bet chose as the one survivor; that narrowed choice is no filter result. `PromptFile` reads each table and refuses a result whose option names are not the table's, in order, a result with no survivor and a survivor that is not an option; whether the set of decisions is complete is the caller's, as `SessionStore` has no call that lists a directory. **The base names.** The existing refusal of two records of one base name applies to the records that one prompt holds, and names both; a decision named `draft`, `spans`, `answers`, `tests` or `tasks` collides with a record of its own prompt and is refused. That no real run meets it is a contract of row 66h, whose step must refuse those decision names (`optionsKeep` is not on this main: a required contract, not a reuse).

**The requirement IDs (R10, binds 66r; condition 3 of the plan review of #364, comment 6249724875).** The correspondence is by kind, from the pinned revised draft: the row `R-NNN` of kind `requirement` is `REQ-NNN` and the row of kind `non-functional` is `NFR-NNN`, the same three digits; a row of kind `success-criterion` is no requirement ID. The `requirements` of **both** initial rows of `tasks.tsv` (the specification task and the target's architecture task) are exactly the set of those IDs over all the rows of the pinned draft, in the order of the draft. The number a row has in the pinned draft is the number of its PRD identity: a later revision of the draft changes no ID that a PRD row holds, and the PRD, not the draft, is the identity from the merge. Row 66r writes the rows by this correspondence and tests the item `spec-sections` for them, with a missing-section refusal; this branch writes the sentence and opens no task. The later architecture tasks of row 66j, and the option constraints, use the same IDs.

**The specification output (R10, binds 66j and 66k).** The instruction asks for what the existing checks read, so a head that waits for its bet is accepted: one five-cell row for each requirement ID with empty MoSCoW and Phase cells, one acceptance row, one matrix row and one section of `docs/spec/` for each ID, one success row for each `success-criterion` row of the draft, and Facts cells whose tokens the checks resolve. A fact is cited as `F-<record>#<n>`, the record being the numbered facts record of the problem statement, `docs/facts/F-NNNN-numbered-facts.md` of the decision of #412 (`T-cav3`), and `n` the number of the fact in `spec/spans.tsv`, which `RenderFacts` writes as the line `n.`; a fact of the answers is cited as `Q-NNN`. The record does not exist until the rendered commit of the first bet (row 66j), and its `NNNN` is one more than the highest `F-NNNN` of `docs/facts/` at the base of that commit, which equals the base of the specification task; so the prompt gives that number as `{facts}`, and, as the pinned `prd-lint` resolves a cell only when one of its tokens names an existing record, the cell also cites `{record}`. The consumers bind to it: row 66j renders the record under `{facts}`, and the call of the full check on a head before the bet (rows 66k and 66n; `T-cav3` for the record that exists) gives `{facts}` among `PSBRecords` and `{record}` as `AnswersRecord`. Until it does, `covers` fails such a head: a limit of the consumers, stated here and not hidden. The rendering stays row 66j's.

**The errors.** `PromptFile` writes nothing at any stage, and every error returns no context, so no session starts on a partial prompt.

| Stage | Case | Error |
| ----- | ---- | ----- |
| 0, before every read | a step with no row in the step table; a step with no instruction | `the step table has no step "<s>"`; `the step <s> has no instruction in M2e` |
| 1 | no row of the task; the step and the kind disagree; options missing for a step that takes them, or given for one that does not | `<tasks.tsv> has no row of <t>`; `the step <s> is not a step of the task <t>, of kind <k>`; `the options of the task <t> are ...` |
| 2 | `verification` with no `attempt` event | `the task <t> has no attempt` |
| 3 | a Shape input, the ledger or an options table missing or refused; a result that disagrees with its table; the specification of a target with no record in `docs/facts/` | `the records have no <path>`; `<path>: <reader error>`; `decision <d>: ...`; `the target has no record in docs/facts/` |
| 5 | two records of one base name | `two records are named <name>: <a> and <b>` |

**The session base** (that the context was built for the base and the task that the session starts with) is checked at the start of a session, and is row 66r's (issue #415), not this task's.

## REQ-005 — The handoff

A **handoff** is the typed result of one session for the next step: its session
ID, the task, the role, its status (`completed`, `blocked`, `needs_context`,
`decision_needed`, `failed`), the reason, each artifact with its path and
SHA-256, the questions, the proposed decisions and the lessons
(`architecture.md` §8). `M2b` built the typed result with the kinds `status` and
`artifact` ([`records.md`](records.md#nfr-001--the-records-of-a-session), the
block `result`); `M2e` adds the rows of the other fields as kinds of the same
file, `question`, `decision` and `lesson`, and a kind `comment` for an artifact that
`layup run` posts as a comment (**decided here**; the block changes in
the build task, which moves nothing in `built`, as the plan review of #333
asked: a built block changes in words here and in its owner's code there).
The Inter-Role Communication Format of `F-0003#59` is the share of valid
handoffs among all handoffs, which `layup report` (`M4b`) reads from the records.

### The artifact that is a comment

The plan, the plan review and the review record of a task are comments of its
issue, so no commit holds them (#68). **Decided here** (decision 1 of #333, as
corrected by the answer to round 1):

- **The session writes a file, code writes the comment.** The session that
  writes the plan, the plan review or the review record leaves its text as a file
  of its `result/` directory, which its typed result names as a **`comment` row**
  (not an `artifact` row: the check of `M2b` refuses an `artifact` that the session's
  head lacks, and this file is in no commit) with its path in `result/` and its
  SHA-256, which the end of the session checks against that file
  ([`session.md`](session.md#req-005--the-result-of-a-session)); a `comment` row
  never names a path of the target. **Decided here** (task `T-q5k4`, row 66d): a row
  whose file is missing, is a link as its last part, is not a regular file, is over 1 MiB
  or has another SHA-256 than the row's adds the event `refused` with the reason `artifact`
  (one refusal, even when an `artifact` row fails too), the result is still committed, as for an
  artifact, and the file of that row is not kept; a result that holds a comment ID in a
  `comment` row is no valid result (the block refuses it: class `no-result`). **The file is kept:** at the records commit of the
  result, `layup run` copies it byte for byte as
  `tasks/<task>/results/<session>/<file>` ([the records](records.md#nfr-001--the-records-of-the-task-loop)),
  because the session's directory is removed after its records are pushed
  ([`session.md`](session.md#the-session-directory)); a run that stops after that
  commit and before the post posts from the copy at the next start, and the posted
  bytes are those of the copy. **Confirmed at the build** (task `T-q5k4`, 2026-10-16,
  `runs/T-q5k4/`): GitHub gave back, byte for byte, seven bodies posted as comments of
  `pharzam/layup` #360 by the App and read by `GET /issues/comments/{id}`: one with CRLF,
  trailing spaces and tabs, a letter outside the Basic Multilingual Plane and no final
  newline, and bodies of 65,536 and 65,537 characters (ASCII and three-byte letters), of
  131,072 bytes and of 262,143 bytes. The size limit is **not** what its own message says:
  that message ("maximum is 65536 characters") came with the refusal of a body of 131,072
  three-byte letters (393,216 bytes), and bodies of 65,537 characters and of 262,143
  bytes were accepted. **Decided here:** code keeps the documented limit, the safe side:
  a file that is empty, not valid UTF-8 (JSON would replace its bad bytes, and the three
  hashes would differ) or over 65,536 characters (counted in characters, not bytes) is
  refused as `refused` `handoff` before any forge call, once for the session, with no
  post, no copy and no event `handoff`. A body that the forge gives back with other
  bytes is not trusted on that account: it is caught by the three hashes below.
  The session posts nothing and holds no forge credential. After the result is
  valid, `layup run` posts that file's bytes, unchanged, as the comment of the
  task's issue (`architecture.md` §8 step 2 and step 6), and copies the comment
  into the records branch as `copies/<comment-id>-<seen>.md`, with its row in
  `copies.tsv` (the author's numeric ID and login, the comment ID, whether an App
  made it, the time, the SHA-256 of the body), by the mechanism of `M2a`
  ([`records.md`](records.md#nfr-001--the-records-of-start), the block `copies`).
  **The order of the post** (**decided here**, task `T-q5k4`; the call is `internal/run`'s
  `PostComment`, a stop between two of its acts is resumed from the records and the forge,
  never from memory): the kept copy is read from the records; a row of `copies.tsv` on the
  task's issue with the file's SHA-256 ends it (no post); else a comment of the App on the
  issue whose body is the file's bytes (a post whose copy was lost) is copied and no
  second is posted; else the lease is read from the records (a lease of another run, or
  not held, is a lost lease and nothing is posted) and the one `Comment` call is made;
  the forge is read again for that comment, as `Comment` gives only the ID, and the
  comment is copied in one fenced records commit. The window between the lease read and
  the post is one call: a takeover inside it could post once, and the next run adopts that
  comment as the copy.
- **The link is written by code.** The comment ID is never in the session's
  result. The `handoff` event ([The checks of a handoff](#the-checks-of-a-handoff))
  holds the comment ID that `layup run` read from the row of `copies.tsv` whose
  `sha256` equals the artifact's SHA-256 and whose `issue` is the task's issue (the last
  such row of the file, **decided here**, task `T-q5k4`, as a comment edited back to its
  first text is a new copy of those bytes). At each handoff `layup run` reads the comment
  again (copy before read), and the copy of an edit is written in the same records commit
  as the event. A comment never posted has no row: the event is `invalid artifacts`, with
  no comment named.
- **Code accepts the handoff only when three SHA-256 are equal:** the session's
  file, the copy, and the body that the forge gives for that comment at the step.
  A comment edited after the copy is a new copy with a new `seen`, the forge's
  body then differs from the first copy, and the handoff is invalid until a new
  handoff names the new copy.
- An artifact of the target's files (the task file, the tests) is named by path
  and SHA-256 at the commit, as before.

**Rejected:** a file in the task's branch (the plan would enter the target's tree
and its diff); a comment ID or a hash that a session writes (it proves only what
the agent says, `F-0003#59`). **Reason:** a copy by the controller with a hash is
the one form that code owns, and the forge keeps the original. The search found no
tool that does it.

### The transition table

`transitions.tsv` on the records branch is the table of the allowed handoffs; its
default rows are those of `architecture.md` §8, and the table is a register that
the Operator can replace. **Decided here:**

- `layup run` writes `transitions.tsv` from the default rows, which the binary
  embeds, in the same commit that opens the first task of the target; a later run
  reads it and never rewrites it from the default.
- **The kinds of artifact** that a row may name are: `plan`, `plan-review`,
  `review-record`, `task-file`, `commits`, `rows`, `patch`, `rendering`, `question`
  and `answer`.
- **The names of the items** that a row may name are in the table below, one word
  each, with what code checks for it. The block `transitions` holds the words.
- A step of the task loop that is `layup run`'s own act and no handoff (the
  issue, the draft pull request, the checks, the merge: steps 1, 4, 5 and 8) has
  no row. **The check of a table** is that it has a row for each default pair
  `(from, to)` of the table below, and for no other pair than a pair that the
  Operator's table lists with kinds and items of these lists; and that the row of
  code's own handoff, `code (rendered task)`, lists no `schema` (**decided here**,
  R10, task `T-dx1h`, the comment of #358): code writes that handoff and a rendered
  task has no session result, so a row that asks for the `schema` could never be
  satisfied and no rendered task would be valid under it. This is the one row-specific
  rule: every other row keeps its `schema` (the row `owner role, asking role` becomes
  satisfiable when #377 gives the kind `answer` of the block `result`), and a row
  that leaves out an item is the Operator's to leave out. The function
  `handoff.CheckTable` makes the check, at the `set` and at each read of the table.
- The Operator replaces the table with `set transitions <sha256> because
  <reason>` on the control issue ([The parameter change](#req-012--the-parameter-change)),
  where the value is the SHA-256 of a copied comment **of the control issue** (a
  restriction of "a copied comment": the `set` is on that issue) whose whole body,
  byte for byte, is the table with its header row (no code fence, no carriage
  return); code reads it with the same schema and refuses a table that fails the
  check above. The new `transitions.tsv` is the bytes of that copy, written in the
  one records commit that also sets the row `transitions` of `parameters.tsv`, so the
  SHA-256 of the file equals the value of the row. The author of the comment that
  holds the table is not checked: the Operator's `set` names it by its hash. The
  old table stays in Git history.

**Rejected:** a table in the code only (it could not be replaced, `architecture.md`
§8); a free handoff file (the public tools validate a record against a schema but
leave the transition to the agent).

| Item | What code checks |
| ---- | ---------------- |
| `schema` | the typed result passes the block `result` and its rules |
| `artifacts` | each named artifact is in the commits (or in the copies) with its hash |
| `plan-fields` | the plan has its four fixed fields and lists its tests by ID |
| `verdict-field` | the review has the verdict `approve` or `approve-with-conditions` in its fixed field |
| `tests-red-green` | each listed test fails at the base and passes at the head, by typed per-test results (the producing run is `M2g`'s) |
| `spec-sections` | a section of `docs/spec/` for each requirement ID the task names |
| `patch-per-kind` | a known-bad patch for each kind a rule batch touches; `layup gate` passes on its head and fails on each patch (`M2f`) |
| `register-rows` | the task register rows of a milestone plan: each `Must` requirement has a task, each predecessor is a task, no cycle (`M2g`) |
| `rendering-equal` | the change equals code's rendering of the named copy |
| `nothing-material` | the review record says `nothing material in scope` at the head |
| `findings-located` | each finding of the record is at a `file:line` |
| `heading-ids` | a heading of `docs/spec/` holds each requirement ID of the milestone (`M2g`) |
| `question-row` | a question row names the records it concerns |
| `answer-row` | an answer row cites the records it rests on |

### The checks of a handoff

A handoff is valid only when code finds the schema right, each named artifact in
the commits (or in the copies, for a comment) with its hash, and each required
item true; a pair with no row in the table is not valid; and a field that an
agent writes never makes it valid (`F-0003#59`, Sol-21, Fable-M14).
**Decided here:**

1. **The schema.** The typed result of `M2b` and the rows of the kinds above; a
   result that fails it is `refused` with the reason `handoff` and the line of the
   row, as a refused result of `M2b` is.
2. **Each artifact.** A path and a SHA-256 of the file at the session's head; a
   comment by the three equal hashes of [the artifact that is a
   comment](#the-artifact-that-is-a-comment). A path that is not in the commits,
   or a hash that differs, makes the handoff invalid.
3. **The required items** of the row, by the table below.
4. **The pair.** The `from` and `to` of the handoff are a row of `transitions.tsv`; a pair with no row is `invalid pair` (the one detail that is not an item of the table below).
5. **The result.** `layup run` writes an event row of the task, `handoff`, in
   `tasks/<task>/events.tsv` with the detail `valid` or `invalid <item>`, and, for
   an artifact that is a comment, after a space `comment <ID>-<seen>` (the kind is
   added in prose here and in the block by the build task).
6. **What the check is given** (**decided here**, task `T-q5k4`, row 66d; the call is
   `internal/run`'s `RecordHandoff`). The schema is the result file of the records,
   `tasks/<task>/results/<session>.tsv`, read by the block `result` (a handoff of a
   session with no such file is `invalid schema`). The items `question-row` and
   `answer-row` are computed from that result by code: `question-row` is true when the
   result has at least one `question` row and each names at least one record after its
   label and its `needs_human` word, and `answer-row` when it has at least one `answer`
   row (the block gives each one or more records); neither resolves a record. A value of
   the caller's for `schema`, `artifacts`, `question-row` or `answer-row` is never read,
   nor is one for `spec-sections` or `rendering-equal` (**decided here**, task `T-yeav`,
   #404: they are computed from the tree at the head, in the row of the pair that lists
   them). Every other item is the caller's result, passed in the value it hands the check
   (`plan-fields`, `verdict-field`, `nothing-material`, `findings-located`, `heading-ids`
   and `register-rows` from the text of a comment or a table; `tests-red-green`
   and `patch-per-kind` from a program, a gate run
   or a read of a commit): an item it does not pass is false. The rows that compute
   them, and the rule text of each, are [the items a caller
   computes](#the-items-a-caller-computes). A handoff that code writes (the rendered task) has no session:
   the event's `session` is `—` and its attempt is the caller's. A handoff whose event
   is already in `events.tsv` with the same session and detail is not written again;
   a different outcome of the same session is a new event.

| From, to | Artifacts | Items |
| -------- | --------- | ----- |
| task plan, plan review | `plan` | `schema`, `artifacts`, `plan-fields` |
| plan review, developer | `plan-review` | `schema`, `artifacts`, `verdict-field` |
| developer (build task), verifier | `commits`, `task-file` | `schema`, `artifacts`, `tests-red-green` |
| developer (specification), verifier | `commits`, `task-file` | `schema`, `artifacts`, `spec-sections` |
| developer (architecture), verifier | `commits`, `task-file` | `schema`, `artifacts` |
| developer (rule batch), verifier | `commits`, `task-file`, `patch` | `schema`, `artifacts`, `patch-per-kind` (`M2f`) |
| developer (milestone plan), verifier | `rows` | `schema`, `artifacts`, `register-rows` (`M2g`) |
| code (rendered task), verifier | `rendering` | `artifacts`, `rendering-equal` |
| verifier, close-out | `review-record` | `schema`, `artifacts`, `nothing-material` |
| verifier, developer | `review-record` | `schema`, `artifacts`, `findings-located` |
| specification, milestone plan | `task-file` | `artifacts`, `heading-ids` (`M2g`) |
| asking role, owner role | `question` | `schema`, `question-row` |
| owner role, asking role | `answer` | `schema`, `answer-row` |

**A question during the work.** A session that meets a question ends with the
status `needs_context` and the question: its text, its own label for the kind,
its `needs_human` field and the records it concerns. `layup run` records it with
its time, posts it on the task's issue and, in `M2e`, sends it to the owner role
of its kind (`architecture.md` §9, the default owner map: domain or business, the
Domain Expert; architecture boundary, the Systems Architect; interface contract, the
Software Architect; environment or infrastructure, the Software Engineer; the
confirmed map, `owners.tsv`, is phase 3) when it needs no human and its label is one
of those four kinds, and to the Operator
when it needs one, or when its label is none of the four (**decided here**, task `T-q5k4`,
which settles the words "the Operator or the idea owner": `architecture.md` §10 names the
idea owner for the points P1 and P5 and "the Operator for the others", and for P2, the
point of this question, it has the Operator classify the kind when the provider fails;
no clause names the idea owner for a question that needs a human); the smart-if that gives the kind and
`needs_human` is `M3a`'s, so in `M2e` the session's own label and field are
used, and a question whose `needs_human` is not `yes` or `no` goes to a human
(**decided here**, the conservative reading). The answer's pass through the
escalation screen (`architecture.md` §8, §10) is `M3`'s ([Not in M2e](#not-in-m2e)).
The owner role's answer cites the records it rests on; the asking role starts its
next attempt with the answer in its prompt file. The answer is **accepted** when
that attempt ends `completed`, cites the answer ID, and none of its questions
cites the answer ID; the accepting actor is that session. An attempt that a
question ends counts only toward `stall.attempts` (`M3d`).
**Decided here** (task `T-q5k4`, row 66d, which settles #377 by R10): the session of the
owner role writes its answer as a row of the kind `answer` of its typed result
([`records.md`](records.md#nfr-001--the-records-of-a-session): the text, and the records it rests
on in `reason`, at least one); the **answer ID** is `<S-xxxxxxxx>.<n>`, that session and
the row's `n`; a result **cites** an ID when it is a word of the `reason` of its
`status` row, and a question cites it when it is a word of the `reason` of its
`question` row. The session of the owner role is recorded under the attempt it answers
(no attempt starts for it); the answer is the answer to the attempt that last ended
`needs_context` before that session started, and the **next attempt** is the attempt
with the next number, started after that session ended and neither `closed` nor
`rebased`; one of its sessions that ended `done` with the status `completed` and the
cite accepts it, the first in the order of the events. The acceptance is a function of
the records (`internal/run`'s `AnswerAccepted`), not a record of its own: the columns of
the accepted answer and its actor are `M3c`'s (`records.md`, the block `questions`).
**Decided here** (task `T-b3s3`, #401, under the plan review of #401): the rule of the owner
role, which is the one place that states it. The session that wrote the answer has to be of the
owner role of **every** question of the attempt that it answers — the question rows of every
session of that attempt that ended `done` with `needs_context` — by the role of its start row in
`sessions.tsv`, which `layup run` writes from the role it starts the session as and no result of
a session can claim. The owner role of a question is `RoutingRole` of the owner of its label in
the owner map (`internal/run`'s `Owners`, the default map in `M2e`), the role of a session is
read by `Roles`, and the answer is accepted by the next attempt only when both are the same role
ID. A session whose start row is missing, a question that goes to a human, a label of no owner
and an attempt whose questions have more than one owner role are never closed by a session's
answer: they fail closed (the task stays uphill), because one accepted answer closes every
question of its attempt and no field binds an answer row to one question. Who resolves those
attempts (the Operator's answer is no session's `answer` row) is `T-9nbf`'s (#367) or later
text's, not decided here. A known limit: the `From` and `To` of a handoff are role-class labels
of the transition table, and `RecordHandoff` reads no session's role, so the handoff of an
answer is `valid` whatever role wrote it; the handoff closes no question, and the acceptance is
what reads the role.
**Decided here** (task `T-8bwd`, #387, under the plan review of #387): the question's
comment and the answer's comment are composed by code from their rows — first line
`<session>: <role> session of <task>, attempt <n>: question <q-id>`, and for the
answer `answer <a-id> of <q-id>` in the place of `question <q-id>`, where `<q-id>`
and `<a-id>` are `<S-xxxxxxxx>.<n>`, the session that wrote the row and the row's
`n` (the form of the answer ID), then, after one blank line, the row's text with
no final newline — posted on the task's issue
and copied by the same machinery as a session's comment: a plain comment of the
App, **not** an artifact of a handoff, and `commentKind` keeps its three kinds. A
restart recomposes the same bytes from the row, so the adoption of a comment of
the App that exists applies. **Decided here** (task `T-gyvo`, #405, row 66o, under
its plan review): the text follows one blank line and the body has no final
newline, because the pinned form named a first line and the text but not the
separator or the end; `n` is the result's column `n`, never the text or the
reason; the answer's `<q-id>` is named by the caller and checked by the form of
an ID, neither read from the row's reason nor resolved against the records. The
composers take the session's start row and one row of its result, as the
records' readers give them, and give an error and no bytes for a start row or a
row that fails the reader's schema and types or its rules (with no repair of a
byte), for a row of another kind, and for a body that `CheckCommentBody` refuses
(a `HandoffRefusal`). [The items a caller computes](#the-items-a-caller-computes)
names the row that posts them.

### The items a caller computes

`internal/handoff` computes `schema` and `artifacts`, and `internal/run` computes
`question-row` and `answer-row` from the result (`RecordHandoff`, row 66d);
every other item is the caller's result, and **decided here** (task `T-8bwd`,
#387; the assignment of rows is the decision of the plan review of #387, an
`approve-with-conditions`) each item below is computed by the named package in
the named row — rows 66l (`T-4pje`, #396: the four comment items), 66n (`T-yeav`, #404: the
tree items) and 66o (`T-gyvo`, #405: the composers), the item computers of the task loop —
by the rule text of the table. The forms the loop's sessions write in the three
comment artifacts are pinned here in one grammar, so two implementers read one
form.

**The fields and the blocks of a comment.** A plan, a plan review or a review
record is parsed by lines, and two line forms count — a field and a block
marker; every other line is ignored. A field is a line whose key starts at
column 0, made of lowercase words separated by one space, and is followed by
`: `, and whose value runs to the end of the line with trailing spaces and one
final `\r` ignored (the forge and a checkout can carry CRLF); `—` (U+2014) is
the empty value of a field, the `key: ` line of a required field with nothing after
the space is malformed (it is not `—`), and a `key:` line without the space is no
field. A
block marker is a line that is exactly a marker word and `:` — `tests:` in a
plan, `findings:` in a review record, never a value after the colon: `tests:\r`
and `tests: x` are not markers (the second is a field) — and its block is the
run of `- ` lines directly under it, ending at the first line that does not
start `- `; a `- ` line inside the run that does not match the block's bullet
form is malformed and makes the item false, as is one whose text holds a
`\r` (a path that is clean can still hold one), and a `- ` line outside every
block is ignored prose. A line inside a fenced code block is not read: a fence
opens at a line at column 0 of three or more backticks or three or more
tildes (the rest of the line is an info string), and closes at a line of only
that character, as many as the opener or more — a closing line with a `\r`
does not close — and a fence never closed runs to the end of the comment.
Each required field and each required marker appears exactly once — twice or
not at all makes the item false — other field lines are ignored, and field and
block order is free. The parse is of the whole comment: a comment that fails
any part of its kind's form makes every item read from it false (a bad `head:`
makes `findings-located` false as well), a body that is empty or not valid
UTF-8 parses to nothing, and a comment is read only for the items of its own
kind. **The plan** holds each of `size:`, `open questions:`,
`new interfaces:` and `new dependencies:` once — `size:` is `small` or `large`
and each of the other three is a one-line text or `—`, read as an empty list for
`—` and otherwise as one element, the whole text, with no separator invented
(the tier reads emptiness only; a consumer that needs the elements pins its
separator then) — and the `tests:` marker once, its block holding a line
`- <id> <path>` for each listed test: exactly two words, the ID unique in the
list, and the path relative, clean and inside the tree (the register's own
`path` rule of `internal/rules`, exported as `ValidPath`; a path with a space
is malformed here). The
block may hold no line — that is the form of a specification task's plan,
which lists no tests — and `tests: —` is not the marker. **The plan review**
holds `verdict:` once, one of `approve`, `approve-with-conditions` and `reject`
(the words `review-record-lint.sh` reads). **The review record** holds `head:`
once — 40 or 64 lowercase hexadecimal characters, as the check already accepts
for commits — `verdict:` once — one of `material`, `nothing material in scope`
and `not mergeable, findings recorded` (the words `engineering-discipline.md`
fixes for a round) — and the `findings:` marker once, its block holding a line
`- <path>:<line> <rule>: <text>` per finding — the first space ends
`<path>:<line>`, the last `:` of that word separates the `line`, an integer of
1 or more, the path passes the same `path` rule, the `rule` is one word without
a `:`, and the `<text>` is not empty; the count of material findings is the number of bullets
of the block. These are the forms the loop's sessions write for a target;
`review-record-lint.sh` reads a different, LAYUP-own form (the
`| Field | Value |` table), and neither document is read as the other's check.

| Item | The input and the rule | Computed by | Milestone |
| ---- | ---------------------- | ----------- | --------- |
| `schema` | the result file and its records: the typed result passes the block `result` and its rules | `internal/handoff` `Check` (row 61, done) | M2e |
| `artifacts` | the commits and the copies: each named artifact is there with its hash, a comment by its three equal hashes | `internal/handoff` `Check` (row 61, done) | M2e |
| `plan-fields` | the plan comment's kept copy: it parses in the pinned form — the four fields each once, the `tests:` marker once, and each `- ` line of its block a `- <id> <path>`; an empty block is valid | `internal/run` (row 66l, `T-4pje`) | M2e |
| `verdict-field` | the plan-review copy: its `verdict:` is `approve` or `approve-with-conditions` (`reject` is false) | `internal/run` (row 66l, `T-4pje`) | M2e |
| `tests-red-green` | typed per-test results — the ID, the base verdict, the head verdict and the source hash of each listed test: each fails at the base and passes at the head, and its source hash is the frozen list's | `internal/run` from typed values (row 66e — #361's typed boundary); the producer — the gate run of the task's test kind and the mapping of an ID to its file — is the specification task of `M2g`'s | M2e item, `M2g` producer |
| `spec-sections` | the `docs/spec/` files at the head and the task's requirement IDs: a non-empty section whose heading holds each ID — the `sections` predicate of the full check over the task's IDs. **Decided here** (task `T-yeav`, #404): the files are the regular files (mode `100644`, type blob) whose names end in `.md` directly in `docs/spec/` of the head's tree, read with `git ls-tree` and `git show` in the run's clone — a nested file, a link, a gitlink, an executable file and a file of another extension are not read; the IDs are the `requirements` of the task's row of `tasks.tsv`; no ID, no file, no head or one ID without its section is false; the head is code's: the SHA of the `bound` event of the handoff's session, never the clone's `HEAD` and never a field of the caller | `internal/spec` (`HasSection`, the one predicate that `sections` also calls), wired by `internal/run` (row 66n, `T-yeav`, #404) | M2e |
| `patch-per-kind` | the batch's known-bad patches: `layup gate` passes on the head and fails on each patch of a kind the batch touches | the rule-batch tasks of `M2f` | M2f |
| `register-rows` | the milestone plan's rows artifact: each `Must` requirement has a task, each predecessor is a task, no cycle | the milestone-plan task of `M2g` | M2g |
| `rendering-equal` | the tree at the head and code's rendering of the named copy: byte-equal — `internal/spec`'s `RenderPRD` and `RenderFacts` for a specification copy, `internal/rules`' `Render` for a record copy. **Decided here** (task `T-yeav`, #404): the caller hands typed `rules.Record` values and the name of the revision (`PreRev`) that a PRD record is rendered from, never bytes; the paths are those of the artifacts of kind `rendering`, and they match the records one to one (no artifact, a repeated path, a path with no record or a record with no path is false); the head is the `Head` that code sets on a code-written handoff (none is false); the rule for each record is `rules.RecordRendered`, the one that `Judge` applies — a regular file (`100644`) at the head whose bytes equal `Render`'s, with no finding. A path that fails `rules.ValidPath` or is the guardrails file, a kind that is neither `prd` nor `facts`, and a PRD whose `PreRev` is empty, is the head, is not an ancestor of it or does not hold the PRD as a regular file are refused (false) before any later read. **Every failed call of `git` is an error, never false** — the resolution of the head or of a `PreRev` too — and the handoff is then not recorded | `internal/rules` (`RenderingEqual` and `RecordRendered`), wired by `internal/run` (row 66n, `T-yeav`, #404) | M2e |
| `nothing-material` | the review-record copy and the head code holds for the verifier's session — the `base` of that session's start row in `sessions.tsv` (`records.md`; its clone is made at that commit, `session.md`), never a field of the result and not the session's `bound` event, which exists only after a passed push: `verdict:` is `nothing material in scope`, `head:` equals that head, and the `findings:` marker is present with an empty block | `internal/run` (row 66l, `T-4pje`; the same parse serves `layup/verify` of row 66f) | M2e |
| `findings-located` | the review-record copy: the `findings:` marker once, every `- ` line of its block a `- <path>:<line> <rule>: <text>` finding, and at least one finding — an empty block is false, an empty list is never a pass | `internal/run` (row 66l, `T-4pje`) | M2e |
| `heading-ids` | the `docs/spec/` files at the head and the milestone's requirement IDs: a heading holds each ID — the `sections` predicate over the milestone's IDs | the milestone-plan task of `M2g` | M2g |
| `question-row` | the result: at least one `question` row and each names at least one record | `internal/run` (row 66d, done) | M2e |
| `answer-row` | the result: at least one `answer` row and each cites at least one record | `internal/run` (row 66d, done) | M2e |

The edge cases read one way: `tests:` with no bullet is the valid plan of a
specification task, and `tests-red-green` reads an empty list as false for a
build task; `tests:` absent or twice, or `findings:` absent or twice, makes
`plan-fields` or each review-record item false; a `- ` line that is no bullet
of the block's form inside a block makes the item false, while the same line
after the block's end is ignored prose; `verdict: nothing material in scope`
with a non-empty `findings:` block makes `nothing-material` false, and
`not mergeable, findings recorded` with an empty one makes `findings-located`
false. The parsed plan fields likewise feed the task's tier and context read:
the parse is row 66l's, the consuming wiring the step loop's (row 66k) — a
declaration of the boundary, nothing of it claimed built. An item is true only
for the comment kind it reads (a `plan-fields` listed on a review-record row is
false), the four comment items are computed from the kept copy and overwrite the
caller's value whether it is `true` or `false`, and a handoff with no kept copy,
or that code writes, has them false.

**The question and answer comments** are composed (composition only, no post) by
`internal/run` (row 66o, `T-gyvo`, #405) from the result's `question` and `answer` rows — first line and body as [a question during
the work](#the-checks-of-a-handoff) pins — and posted on the task's issue by the
post→copy machinery of the step loop (**row 66k**, #367, whose comments of a task
session it is); they are plain comments of the App, not artifacts of a handoff.
The parsed review record is also the input of the caller boundary that keeps
the loop uphill while a material finding is open (`architecture.md` §11):
row 66l's parse exposes the review's verdict, head and located findings and a
predicate that fails uphill — a record that does not parse, a verdict other than
`nothing material in scope`, or a bullet in the block is an open material
finding, so `material` or `not mergeable, findings recorded` with an empty block
never reads as none open — the step loop's wiring of it stays row 66k's (#367),
and nothing of either is claimed built here.

**The deferrals are recorded, not hidden.** `patch-per-kind`, `register-rows` and
`heading-ids` name checks whose inputs — a rule batch's patches, a milestone
plan's rows, a milestone's requirement IDs — exist only where [Not in
M2e](#not-in-m2e) puts them (`M2f`, `M2g`), and the transition rows carry the
same marks, so no `M2e` run makes one and none is silently passed.
`tests-red-green`'s producer is `M2g`'s for the same reason: the item's `M2e`
contract is row 66e's typed values, a build-task handoff cannot be valid in
`M2e` (no producer exists), and the demo and the end-to-end row exercise a
specification task, whose transition row does not list the item. An empty
`tests:` list makes `tests-red-green` false for a build task.

### The frozen test list

The plan lists its tests by ID (`architecture.md` §8 step 2). **Decided here**
(`architecture.md` §11, "Progress is computed"):

- The list is **frozen** at the task's first valid handoff to the verifier, when
  each listed test fails at the base and passes at the head: `tasks/<task>/tests.tsv`
  gets one row for each test, its ID, the SHA-256 of its source file at the head
  and the harness that wrote it, and the time.
- A later head where the SHA-256 of a listed test's source differs makes that
  test not counted as passed, and the change is a finding `test changed`.
- An **amendment** is a new plan, which the plan session writes on that finding,
  and a new plan review; the list is frozen again at the next valid handoff, and
  the finding is closed.
- A test that fails at the head, or passes at the base, makes the handoff to the
  verifier invalid, so the list is never frozen on a test that proves nothing.

**Decided here** (task `T-6zj3`, #361, row 66e; `internal/run`'s `RecordHandoff` with
`HandoffInput.Freeze`, and `FrozenFindings`):

- **The text of the list is #387's.** The caller gives the tests as typed values: the ID, the
  file that holds it, and whether it fails at the base and passes at the head (the producer of
  that result is a gate or test run that is not in `M2e`).
- **The freeze is in the commit of the event.** When a `Freeze` is given to a handoff to the
  verifier, `tests-red-green` is computed from it and the caller's own value is never read:
  it is true for a list of at least one test, each with an ID and a source, none repeated,
  each failing at the base and passing at the head, and each source a file of the head that
  the session was bound to (its event `bound`; a session with no bound head is an error, and
  nothing is written). The first valid handoff writes `tasks/<task>/tests.tsv` in the records
  commit of its event, with the SHA-256 of each source at that head, the harness of the
  session (from `sessions.tsv`) and one time; no other handoff, no invalid one and none to
  another role freezes. An Operator's table that leaves `tests-red-green` out of the row is the
  Operator's: the handoff is as the table says, and the list is still never frozen on a test
  that proves nothing.
- **A changed source is a finding and an invalid handoff.** With a list frozen and no amendment
  pending, a handoff to the verifier is `invalid tests-red-green` unless its tests are the
  frozen ones, each at its frozen hash; `FrozenFindings` gives, for a head, the finding `test
  changed` (the rule, the source as its file, the test's ID as its text) of each listed test
  whose source there has another hash or is missing.
- **An amendment is read from the events.** After the last valid handoff of a session of the
  developer's role (`software-developer`), a valid handoff **with a comment** from a session of
  the plan role (`software-engineer`), then one from a session of the plan-review role
  (`qa-engineer`) **whose event `session` comes after it**: the order is the events', never the
  clock's (the `at` of two events can be one second), and a verifier's review record, which
  shares the plan reviewer's role, is no review of a plan that began after its session. With an
  amendment pending the results alone decide, the next valid handoff to the verifier writes the
  list again whole (when it differs; a test whose source and hash did not change keeps its
  harness, any other takes the session's) and uses the amendment up, as it is that developer's
  handoff.

**Execution of the named tests** — the gate run of the task's test kind that
produces the per-test results, and the mapping of an ID to its file — is `M2g`'s
(the build task's test run). In `M2e` the item `tests-red-green` is computed
from the typed per-test results a caller holds (row 66e), and no caller produces
them: a build-task handoff cannot be valid in `M2e`, and no `M2e` run makes one
(the demo and the end-to-end exercise a specification task). An empty `tests:`
list makes `tests-red-green` false for a build task; a specification task's
transition row does not list the item.

## REQ-007 — The four `layup/` statuses

Requirement `REQ-007`: each change passes the gates before review or merge
(`F-0003#47`), and no pull request with a failed gate reaches human review or
merges (`F-0003#58`). **Decided here** (decision 5 of #333). `layup run` posts
four commit statuses on each head of a task's pull request, each with the call
`SetStatus(sha, context, state, description)` of `internal/forge`
([`forge.md`](forge.md#the-calls-of-m2e)):

| Context | From | Posted when |
| ------- | ---- | ----------- |
| `layup/gates` | the table of `layup gate` at the head | after each new head |
| `layup/spec` | the full `layup spec check` at the head | after each new head |
| `layup/rules` | the basic rule below | after each new head |
| `layup/verify` | the review record of the verifier | after each review, and `pending` while the verifier waits |

The state is `success`, `failure`, `pending` or `error`. `error` means that the
check itself could not run; it is never a pass. **The source** of each status is
the App: the ruleset of `M2d` requires the four contexts pinned to the App's
`integration_id`, the `app_id` of the forge register
([`scaffold.md`](scaffold.md#the-ruleset-of-the-default-branch)), so a status of
the same name from another source does not count. **A status is posted again only
when its state or its description changes**, as GitHub keeps at most 1,000
statuses for one SHA and context. The description is one line of at most 140
characters that names the reason (`pass`, `clear: no product path`, `not-active:
tool not found: <tool>`, `waits for its turn`).

**The permissions** of the App beyond those of `M2a` are pull requests `write` and
commit statuses `write`, both in `architecture.md` §3 and ADR-0014 decision 3, and
checks `read`, which the native jobs need (a job of GitHub Actions is a check run, and
the `mergeable_state` of a draft is `draft`, never `clean`): checks `read` **amends
ADR-0014 decision 3**, asked in #336 (`gov-operator-setup-o112`: the Operator's step,
asked in the demo row of the slicing task, as the first demo of `M2d` asked its
own) ([`forge.md`](forge.md#the-calls-of-m2e)). The step
`forge` of `layup run` reads the granted permissions: Start and the restart before the
loop record them, as for any later milestone's, and the loop's first run (row 66a of
the plan) calls `forge.CheckLoopPermissions` before its first write, so a missing one
is then `fail`, naming it, as it is for the `M2a` permissions (`forge.md`).

### layup/gates

`layup/gates` is the table of `layup gate` at the head, posted from its rows
(`architecture.md` §6, the table; [`gate.md`](gate.md#the-table)):

| Result of a kind | `layup/gates` |
| ---------------- | ------------- |
| `pass` | counts as a pass |
| `fail`, or `not-active` | `failure` |
| `clear` (a `pending` kind with no product path changed, or an `active` kind with no product path) | counts as a pass, with the reason from the row |
| a `pending` kind and the head changes a product path | `failure` |

The status is `success` only when every kind counts as a pass. **Decided here:**
`layup/gates` is read from the table that `layup gate` gives for the head, run by
`layup run` in a scratch work tree with the base's gate files (FT4), and the
table is committed as `tasks/<task>/gates.tsv` with the head SHA as a key (the
block `task-gates`); the head's own gate files never judge the head. A run that
the tool could not make (git, the scratch tree) is `error`.

### layup/spec

`layup/spec` is the full `layup spec check` at the head
([The full spec check](#req-012--the-full-spec-check)): `success` when no row of its
table is `fail` (a row `clear` counts as a pass, whether the whole table is `clear`
because the task touches neither the PRD nor `docs/spec/` before the first bet, or
only the rules `classes` and `inventory` are, before it, on the specification
task's head); `failure` with the first failing rule in the description otherwise;
`error` when the check could not run.

### layup/rules

The **basic rule** (`architecture.md` §6; the batch part is `M2f`'s): `layup/rules`
fails a pull request that changes a path of the rule-path register
(`rule-paths.tsv`, [`setup.md`](setup.md#the-rule-path-register)) unless it is a
**rendered record**: a change to the PRD or the facts record that equals code's
rendering of the named copy (an inventory version of a bet, a decision copy or the
start-values copy), checked by rendering again from the copy and comparing bytes.
**Decided here:** the status is `success` when the head changes no rule path
(the reason `no rule path`) or only a rendered record that equals; `failure`
otherwise; before `M2f` no rule batch is approved, so a change to a rule path that
is not a rendered record fails. The check reads the register of the **base**, not
of the head. The refusal before the push (prevention layer 2) is the first line
of defence and `layup/rules` the second; both read the same register.

**Decided here** (R10, task `T-9et9`, #352; the functions `Basic` and `Judge` of
`internal/rules`):

- **The inputs.** The caller names the base, the head, the commit of the register
  (the records commit) and the rendered records, and, for a PRD, the commit that it is
  rendered from. A record is a path, a kind (`prd` or `facts`) and the copy as parsed
  values: the answer of the bet for the PRD, the version of the inventory and the text
  of each fact for the facts record. `internal/rules` reads no copy and checks no hash
  of a copy: both are the caller's.
- **The register** is the file `rule-paths.tsv` at the root of the commit that the
  caller names, as the check before a push reads it ([`session.md`](session.md#a-workflow-or-rule-path-change));
  the head's tree is never read for it, and a commit that is the head is refused. A
  register with no file, one that is no regular file, one that its reader refuses and
  one with no row is `error`: a register with no row passes every path, and S15 refuses
  to write one ([`scaffold.md`](scaffold.md#the-step-verify-and-the-records-of-the-setup)).
  The refusal of a register with no row is `layup/rules`'s; the reader before a push is
  not changed by this task.
- **The paths** are those that differ between the base and the head, both ends of a
  rename; the base must be an ancestor of the head, else `error`.
- **A rule path** is a path that `Check` refuses (the register, the kind `.sh` and the
  workflow rule): the exception of `docs/guardrails.md` holds here as before the push
  (`architecture.md` §6), and a path under `.github/workflows/` is a rule path (the
  reason `rule path`, not one of its own).
- **A rendered record** is judged only for a changed path that the rule-path check
  refuses. It passes when the head holds a regular file there (mode `100644`: a link, an
  executable, a submodule and a directory are not one) whose bytes equal `Render`, byte
  for byte, with no trimming: the rendering of `internal/spec`, of the PRD from its
  pre-image or of the facts record from the version. A copy that cannot be rendered is
  a `failure`, never a pass. An unchanged path or a path that `Check` allows is not
  judged: under the register of phase 1 the PRD is no rule path
  ([`setup.md`](setup.md#the-rule-path-register)), and its cells are judged by the rule
  `classes` of `layup/spec`.
- **The pre-image of a PRD** is its file at the commit that the caller names: at the
  first bet the verified head of the specification (the parent of the rendered commit),
  at a later bet the base. The commit is the base or lies between the base and the head,
  and is not the head, else `error`.
- **The result** is a state, a reason and the first path in the order of git. `success`
  with `no rule path` or `rendered record`. `failure` with `rule path`, `not equal` (the
  bytes differ, or the record is deleted), `not a regular file` or `no rendering`.
  `error` with `register`, `input`, `base` or `git`; it is a check that could not run
  and never a pass. A record is `input` when its path is not relative and clean, holds
  a `..` or `.git` part, a backslash or a NUL, or is the path of the guardrails; when two
  records name one path, a kind is neither of the two, or a PRD has no pre-image. The
  status, its description and its posting once for an unchanged state are row 66f's.

### layup/verify

`layup/verify` follows the verifier (`architecture.md` §8 step 6). **Decided here:**

- `pending` with the description `waits for its turn` while a build task is not
  next in the merge order or not up to date with the base (a wait state of
  `architecture.md` §11; the merge order is `M2g`'s, so in `M2e` the only tasks are
  one at a time and the wait does not occur); `pending` with `verifier running`
  while its session runs.
- `success` at the head SHA the verifier reviewed, only for a review record with
  `nothing material in scope`; `failure` for a material finding (one more round, up
  to the target's cycle cap, then the stall of `M3d`), and for `not-active` (above).
- After the rendered commit of the first bet, and after the close-out commit,
  `layup/verify` is set on the new head by the check of [The merge at the head SHA
  and the close-out commit](#the-merge-at-the-head-sha-and-the-close-out-commit),
  not by a second review: the verdict carries over to a head that differs from the
  head that carries it by that one commit. A head **carries the verdict** when it is
  the verified head, a head that came from it by the rendered commit, or a head that
  came from it by a clean merge of the base that changes no file the task touches
  (the next bullet), also when the verified head is the close-out head of a Shape
  task that the loop brought up to date with the base.
- A verified task whose base then moves carries its verdict over a clean merge of
  the base that changes no file it touches (the target's own rule); otherwise the
  task is closed and redone as a new task with its own issue, which takes the old
  task's row and predecessor links (`architecture.md` §8). The redone task has an
  `closed` event on the old one and a `rebased` event on the new one
  ([the records](records.md#nfr-001--the-records-of-the-task-loop)).

### The ready and merge rule

**Decided here.** No pull request is marked ready or merges while a selected kind
did not run (`F-0003#58`, `NFR-004`):

1. `layup run` marks a pull request ready (`MarkReady`) only when, at the head,
   each of the four contexts is `success` **by a read of the forge**, not by its
   own earlier post, and each native gate job of the target is green: each required
   check of the default branch's ruleset whose source is GitHub Actions (read back
   by `EffectiveRules`) has a check run at the head with the conclusion `success`, taking the newest run of
   each name (`filter=latest` of the list call; a later `failure` beside an older
   `success` is a failure) (`skipped` and `neutral` do not count).
2. A state other than `success` of any of the four, including a context that has
   no status at the head, a status of another source and an `error`, blocks the
   marking and the merge. GitHub accepts a `neutral` or `skipped` check as met
   for a required check; code counts only `success`.
3. A kind that is `not-active` is never a pass: it makes `layup/gates` a
   `failure`, which blocks both.
4. The merge is the call of the next section.

### The merge at the head SHA and the close-out commit

**Decided here** (F-37; decision 6 of #333):

- **The close-out commit.** After `nothing material in scope`, a close-out session
  writes the target's close-out (the verdict, the line of the completed log, the
  resource record from the ledger) as one last commit. Code checks that this
  commit changes only the task file `docs/tasks/<task>.md` and
  `docs/tasks/completed.md` of the target, the paths from the setup record, and
  that its parent is a head that carries the verdict (defined in
  [layup/verify](#layupverify): the verified head, the head after the rendered commit
  of the first bet, or the head after a clean merge of the base). A commit that changes another
  file is refused as a rule-path or scope refusal, as a developer's would be.
- **The rendered commit.** The commit that code adds at the first bet
  ([The rendered commit and the confirmed inventory](#the-rendered-commit-and-the-confirmed-inventory))
  gets the same check with the rendered files in place of the task file: it
  changes only them, they equal code's rendering of the bet's copy, and its parent
  is a head that carries the verdict. The close-out of the specification task then follows it.
- **The verdict at the new head.** `layup run` sets `layup/verify` `success` on
  the new head with the description `verified head <sha> plus the close-out` (or
  `plus the rendered commit`); the other three contexts run on it as on any head.
- **The merge.** When the four contexts are `success` at that head and the native
  jobs are green, `layup run` calls `Merge(pull request, head SHA, method)` with
  the head SHA it read the statuses for; the forge refuses a head that moved, and
  code then reads the statuses again, as it never trusts the earlier read. The
  method is a plain merge, never a squash (the target's rule). A batch that
  changes `.github/workflows/`, which the Operator pushed, is merged by its
  approver (`M2f`).
- **A tested tool** (F-37): the close-out check and the merge call are one
  function each in `internal/run`, with a stand-in forge, so the steps that land a
  target task are a Go test, not a script.

### The gate jobs read the base

F-22 (the Operator's direction, task `T-evad`): a pull request could weaken its
own gate, as each native gate job runs the head's `gates.sh` with the head's
`docs/gates.tsv`. **Decided here**, in three parts that `M2e` builds:

1. **The jobs read the base.** The catalog workflow makes each
   gate job take the manifest, the `config` paths and `gates.sh` from the base of
   the pull request before it runs them, as `layup gate` does (`gate.md`, The run).
   That workflow is the setup's file, so each target has it from its setup; a target
   set up before `M2e` gets it by a rule batch (`M2f`). **The exception, decided
   here:** a pull request whose head branch is `batch/<point>` takes the gate files
   from its own head, because the activation batch changes the manifest and writes
   Go files in the scope of the `pending` kinds, and with the base's manifest its
   `pending` jobs would fail it for ever (`architecture.md` §6, ADR-0016 decision 2).
   **The match** (decided in task `T-hz4m`, R13; no source gives a form for
   `<point>`, and a point is named once, in a line `<point>: <login>`): the head
   ref is `batch/` and exactly one nonempty segment that holds no `/`, matched on
   the whole ref. `batch/activation` takes its own files; `batch/`, `batch/a/b`,
   `batch`, `xbatch/a` and `feature/batch/x` do not. The reading that fails
   closed holds, as the exception lets a head judge itself; a point that needs
   another form is a change of this text (`M2f`, #345). The step is
   [`setup.md`](setup.md#the-stack-catalog) (the mode `--overlay` of the job script).
   The control for such a pull request is part 2: a change of a gate path waits for
   the Operator as code owner. Whether this exception stands, or another way, is the
   specification of `M2f` (issue #345).
2. **Code-owner review.** A change to `.github/**`, `docs/gates.tsv`, `docs/gates/**`
   or a `config` path of the manifest (the setup writes the file from the manifest)
   needs the Operator's review as code owner: the setup writes a
   `CODEOWNERS` file for those paths with the Operator's account, and the rule
   `pull_request` of the default branch's ruleset sets `require_code_owner_review`
   `true` with `required_approving_review_count` `0`, so that the rule binds a pull
   request that touches a listed path and no other
   ([`scaffold.md`](scaffold.md#the-code-owner-rule), task `T-jew3`: `.github/CODEOWNERS` is in the
   tree of the target, and the rule is in the ruleset file of the default branch). **Unconfirmed:** whether GitHub binds the
   code-owner rule when the count of approvals is 0 is not in a source of the
   search; the build task confirms it on a real repository in the demo. If GitHub
   needs one approval for the rule to bind, every pull request would wait for a
   human, which `architecture.md` §8 step 8 forbids, and the Operator decides
   (a question asked then, severity normal).
3. **The test.** A pull request that weakens its own gate (it edits `docs/gates.tsv`
   to make a kind pass) is refused: its jobs read the base and fail, and the
   change of a gate file waits for the code owner.

**The contract changes** the sentence of `architecture.md` §6 that the target's
own CI runs the workflow and the manifest of the pull request's head: with part
1, the manifest and the gate files come from the base (except for a pull request of a
`batch/<point>` branch, whose gate files come from its head: the exception of part 1),
and the workflow file of a `pull_request` still comes from the head, which part 2 covers
with the code owner's review. `architecture.md` is edited in the same change, with the
exception in it.

**How it fits the ADRs.** F-22 is the Operator's direction (comment 6002406785 of
#119, in `runs/T-evad/findings.md`), later than ADR-0017. ADR-0017 "We reject" code
owners whose approval comes from a second human account (O-77), as its means of
preventing a rule change by an agent; here the code
owner's review is a condition of the ruleset on four path classes (`.github/**`, `docs/gates.tsv`, `docs/gates/**` and the `config` paths of the manifest), for a human
push or a rule batch, and a task of `M2e` never touches them (the refusal before
the push, prevention layer 2, stops an agent's change there). It is therefore no
approval point of the milestone (`architecture.md` §8 and ADR-0019 decision 2 allow
none inside it but escalations and stalls): in `M2e` no pull request waits for it.
A rule batch that touches those paths (the activation of the first bet) waits for
the code owner too: that is an exception to ADR-0019 decision 2 and a use that
ADR-0017 rejects as a means against agents, so **both ADRs need a note**; it is the
issue #336, and the batch part is `M2f`'s (both in [Not in M2e](#not-in-m2e)). The
comment with the decisions of S3 said that no decision changes an accepted ADR; the
correction of round 2 withdraws that for decision 6 and points to #336. The account
of the code owner is the Operator's row of `approvers.tsv`; where `CODEOWNERS` is
written and how the ruleset reaches a target that `M2d` set up is decided by the Operator (comment 6209820437 of #356): the file in the tree of
the target, the ruleset files and the commands for the Operator to review and run, and, for a target set up before the rule, the update of its
ruleset by the recorded ID with its protections kept ([`scaffold.md`](scaffold.md#the-code-owner-rule-of-a-target-that-m2d-set-up)). **Decided
(the child issue #384, task `T-4cep`; the Operator's ruling 6217133467):** the Operator's one push of the commit that lands `CODEOWNERS` on the base of such a
target is an act of the Operator completing that target's setup at the planned point "the setup", and [the tree check](scaffold.md#the-tree-check) admits that commit and no other
head above the setup head (a non-operative summary; the conditions are there). The note of #336 names the act and its actor classification for the audit of `M2f`.

### The tree after the first merge

The tree check of a later start compares the default branch with the head of
`setup/tree.tsv` (and with the code-owner commit of a target set up before the code-owner rule, [`scaffold.md`](scaffold.md#the-tree-check))
only until the first merge after the setup
([`scaffold.md`](scaffold.md#at-each-later-start)). The check of this section is row 66g's and is specified here, not built. **Decided here:** from the
first merge of a task, the check is that the head of the default branch is the
commit of a `merged` event of `tasks/<task>/events.tsv` (the event that `layup run`
writes with the merge commit's SHA when the merge call answers), or a descendant
of one through such commits, **or the merge commit of the pull request of a task of
the register that the forge shows merged and whose `merged` event is not yet
recorded**: a lease lost, or a process killed, between the merge call and the records
commit leaves such a head, and the row `tree` of `forge-check` is `pass` (its detail
stays `—`, as the block says) while the step's own detail line says
`tree: <sha> merge of <task>, not yet recorded`, read from the forge with the call
`PullRequest` ([`forge.md`](forge.md#the-calls-of-m2e)); the step `loop` then writes the
missing event (the restart rule below). A head with no such explanation fails; a commit that reaches the default
branch in any other way (the ruleset forbids it, so this is a read-back's
second line) makes `forge-check` `fail` with the detail `tree: <sha> is not a
recorded merge`. **Reason:** after the first merge the head is no longer the
verified setup tree, and the one thing that may move the branch is a task.

## REQ-012 — Shape

Requirement `REQ-012` (`F-0003#51`): Shape writes the specification, the target's
architecture and the gate-activation batch (`architecture.md` §8). **Decided
here** (decision 7 of #333): for each architecture decision an **architect
session** writes at least two options, each with a fixed table, the block
`options` ([the records](records.md#nfr-001--the-records-of-the-task-loop)): for
every constraint that applies a `pass` or `fail` with its evidence. **The
constraints that apply** are listed by code, never by the session: each requirement
of kind `non-functional` of the draft or the PRD, by its ID `R-NNN` and its criterion (a
requirement of another kind that the architect names in the column `constraint` of its
rows is judged too, but no row is demanded for it); each fact of class `constraint` in `spec/spans.tsv` (before the first
bet; the confirmed constraints of the inventory after it), by `F-<n>`, and a fact
of the answers, which the inventory holds by its ID `Q-NNN`, by that ID (the same
constraint facts, named so); each
invariant of the PSB, by `Inv-<n>`; and each accepted decision record of the target,
by `ADR-NNNN` (its status is `Accepted`, or `Accepted. Amended by ADR-NNNN`, which
the baseline's record rules write for an extended record whose old decision still
holds; a proposed, deprecated or superseded record does not apply).

**The filter is code.** An option is **dropped** when a row is `fail`, a row is
missing (a constraint that applies and has no row), or a `pass` has no evidence
(the evidence column holds `—`). The filter reads the table, never the session's
own judgement of it.

- **One survivor** is chosen.
- **More than one survivor:** in phase 2 no panel sits (the panel is `M3`'s,
  `REQ-010`), so the bet brief shows the survivors, the dropped options with their
  failed rows, and **no recommendation**; the bet decides (`architecture.md` §8,
  "The brief then shows the survivors with no recommendation"). **How the bet
  decides, decided here:** an option's name is a word of lowercase letters, digits and
  `-`, so that `<decision>=<option>` holds no space; the answer to the brief has one
  more line for each such decision, `option <decision>: <option>`, naming one survivor,
  and an answer with a `bet` and no such line for a decision with more than one
  survivor is recorded as human input and answers nothing (the brief stays open); code checks that the
  option is a survivor of that decision's table, and records the choices in the column
  `options` of the row of `bets.tsv`. **What the architecture task writes before the
  bet:** the decisions that have one survivor, written as decided, and, for each
  decision with more survivors, a section `Decision pending the bet` that lists the
  survivors and decides nothing. Its prompt file holds each options table and the
  survivors that the caller's filter gives ([The instruction of a task step](#the-instruction-of-a-task-step), task `T-c6sx`). **After the bet**, code opens one task of kind
  `architecture` for each settled decision, with the chosen option in its prompt
  file (the options table and the chosen option as the only survivor, [The instruction of a task
  step](#the-instruction-of-a-task-step)), which writes it into the target's architecture; it merges in the merge order
  below, after the specification.
- **Every option dropped:** the architect session writes new options once; if they
  are dropped too, the round made no progress. `layup run` posts a comment on the
  milestone's issue that names each dropped option with its failed rows, and stops
  the run as a wait for the Operator (the step's detail says `shape: no survivor`).
  It writes no row of `stalls.tsv`: that row holds the SHA-256 of a stall package,
  which with the stall triggers and their procedure is `M3d`'s ([Not in
  M2e](#not-in-m2e)). **Reason:** open question
  2 of the inventory; the architecture allows no panel with fewer than two
  harnesses, and the filter is the part of the solution routing that code does.

The decisions and their options are tasks' outputs: the specification task, the
target's architecture task and the activation batch are tasks like any other.

### The questions left open

A question still `open` or still `asked` in `questions.tsv` after Intake
([`intake.md`](intake.md#the-review-of-the-answers-and-the-follow-up)) is read by Shape. **Decided here:** before
the options of a decision are written, `layup run` asks each such question of the
owner role of its kind (the idea owner for the intent form, as `intake.md` sets):
one comment on the milestone's issue, copied into the records, and the answer sets
the question's state. **An answer to a question of the intent form** is parsed as the
step `answers` parses it (`B` below `U`, the hours a number): a value that parses is
written to `budget.tsv` with `source` the question's ID, and replaces the value there
when the file exists (the answer to Shape's question is the source, as the follow-up
was in `M2c`); a value that does not parse leaves the question `open` and the run
in its wait. A question that stays `open` when the brief is made is shown
in the brief's "rabbit holes" part. The intent form's budget fields (`B`, `U` and
the hours) are never defaulted: with none of them Shape does not start, and the
run stops as a wait.

## REQ-012 — The bet

**The Bet** is Decision Points 1 and 3, one bet for each milestone
(`architecture.md` §8). **Decided here** (decision 8 of #333): a session of the
Product Owner role writes the **brief**, in five parts:

1. **The problem:** the requirements of the milestone.
2. **The appetite:** its cap of money and wall-clock, in the fixed line `appetite: USD
   <n>, <n> hours`, which fits the rule below (`architecture.md` §12).
3. **The solution:** the architecture, or a proposed change of it, which sends the
   milestone back to Shape.
4. **The rabbit holes and the no-gos**, with the fixed line `architecture change:
   yes` or `architecture change: no`.
5. At the first bet, **every fact with its class** (`architecture.md` §7 step 7).

Code checks the five headings, the line limit, that at least one link is
there and that each link resolves; the brief names the head SHA and the
rule-file hash of each verified batch that it asks to approve, one line
`approve batch <point>: <head-sha> <rule-sha256>` for each. `layup run` posts
it on the milestone's issue.

- **The line limit** is the parameter `brief.lines`, a row of the Operator
  (`architecture.md` §10: "the brief's line limit", set at Intake). No source
  gives its value: `M2c` did not set it. It is a **marker** (row in
  [`open-gaps.tsv`](../setup/open-gaps.tsv)): ‹the line limit of the bet brief›,
  and the check refuses a brief while the parameter is unset, with the message
  `brief.lines is not set`. The Operator sets it, on the evidence that the idea owner gives, with a
  `set brief.lines <n> because <evidence>` ([The parameter
  change](#req-012--the-parameter-change)).
- **The caps.** The appetite of a bet is two columns of its row of `bets.tsv`,
  `appetite_usd` and `appetite_hours`, read by code from the part 2 of the brief's
  copy (the brief states them in a fixed line `appetite: USD <n>, <n> hours`);
  code refuses a brief that does not fit, by `architecture.md` §12: the money
  caps of the bets so far, this one included, are at most `band.U` of `budget.tsv`
  ([`records.md`](records.md#nfr-001--the-records-of-intake)), which stays as `M2c`
  made it; their wall-clock caps are at most `appetite.hours` of the same table; and
  the bets' caps plus the spend of milestone 0 (`telemetry.tsv`) stay within
  `band.U`. A cap below `band.B` is allowed. The idea owner can raise `U` or the
  appetite, or the brief is reshaped. `milestones.tsv` stays phase 3. **Reason:** open question
  3 of the inventory; the block `budget` takes its `source` from a question, and a
  bet's answer is not one.
- **The answer.** The approver of that planned point, the idea owner by default,
  answers by one comment: a line for each requirement, a line for each fact to
  change, a line `option <decision>: <option>` for each decision with more than one
  survivor, and `bet` or `no bet`. **Decided here** (R10, task `T-3q0b`): a
  requirement's line is `<REQ-or-NFR-ID>: <Must|Should|Could|Won't>, phase <tag|—>`,
  with `phase —` exactly for `Won't` and the tag one of the target's §9 phase tags,
  and a fact's line is `fact <number-or-Q-NNN>: out of scope` or `fact <id>:
  <need|constraint|measure|context>`. Code takes it only from an account whose role in
  `approvers.tsv` allows the point; `layup run` copies it, and records the bet in
  `bets.tsv` (the block `bets`). The bet approves the milestone, the architecture
  and its batches; no other approval sits inside the milestone, apart from
  escalations (`M3`) and stalls (`M3d`). A requirement changes only at a bet, or by
  an escalation decision.

**The merge order after a bet:** the specification (or the rendered task), the target's
architecture tasks, the rule batches, then the build tasks (`layup run` merges the
base into the branch of each next task, a Shape task included, as `architecture.md` §8
says for the build tasks, because the branch rule is "up to date"); a bet that changes a class or a priority does not
touch the approved batches.

### The approvers and the success criteria

Deferred by `M2c` ([`intake.md`](intake.md#not-in-m2c)). **Decided here**
(decision 12 of #333):

- The **success criteria** of the intent form are read from the answer to its
  question and kept as the measure rows of the requirement draft
  ([`records.md`](records.md#nfr-001--the-records-of-intake), the block `draft`):
  one row for each criterion, with the idea owner's words, and a measure fact that
  has none is refused by the draft check of `M2c`. **Decided here** (task
  `T-8rm6`): the form's fourth field is `success-criteria`, its criteria
  separated by `;`, each segment trimmed and an empty one dropped; the answer
  is that of the follow-up that asked the field again when the follow-up has
  one, else the batch's (the rule of `budget.tsv`'s fields, task `T-tmm2`),
  and the revision appends one `success-criterion` row for each, in the
  answer's order, covering the question whose answer gave them.
- The **approver of each planned point** comes from the Intake batch and is a row
  of `approvers.tsv` with `source` the ID of the answering comment (the column that
  `M2a` left for it). Code finds the point from the copy of that comment (the
  batch's line `<point>: <login>`), so that a row of `approvers.tsv` and the point
  it serves are tied by `source`, as `architecture.md` §3 asks. The **role `approver`** is the one that an approval point
  names for the account that is neither the Operator nor the idea owner; a point
  that names no approver is approved by the idea owner. **Decided here** (task
  `T-8rm6`): the planned points are `setup`, `bet`, `activation`,
  `retrospective`, `early-retrospective` and `audit`; only a counted answer of
  the idea owner's block names an approver, the last counted line of a point
  binds it, the earliest comment that named a login gives the row's `source`
  and `since`, and whether a login's account already holds a row is decided by
  the numeric ID the forge resolves — a login is never a key, whatever its
  spelling. The same resolved ID decides which existing row serves a point: a
  point whose line resolves to a recorded account is served by that account's
  row, not the idea owner's, however the stored login is spelled. A login the
  forge does not know is input only and gets no row, and the point falls to
  the idea owner; a lookup that fails otherwise resolves nothing and is the
  call's error, never a choice of another account.
- A comment by an account with no row, or in a role the point does not allow, is
  recorded as human input and answers nothing.

### The revision of the draft

After the answers, a session revises the requirement draft (`architecture.md` §7
step 5). **Decided here:** the step `revision` runs it on the role
`product-owner`, after `follow-up-answers` and before the Scaffold
([`intake.md`](intake.md#the-steps-of-intake); task `T-8rm6`). Code checks that a
requirement that covers an answer's
fact also names the fact of the answer's quote (the link kept by `M2c`:
[`intake.md`](intake.md#the-raw-fact-of-the-answers)); a requirement that covers an
answer and not the quote's fact is refused, and the session gets the finding in a
new attempt — a refused revision is run again with the check's finding in its
prompt, until a draft keeps or the cap of milestone 0 stops the step
([`intake.md`](intake.md#input-states), task `T-8rm6`). **Decided here** (task `T-8rm6`): the check keeps the coverage
rules of the draft of `M2c` — a `need` or `constraint` fact must stay covered, a
`measure` fact must keep a success criterion — and adds that a fact covered by
both drafts whose set of classes changed is kept only when a row that covers it
names a fact of the answers; a fact newly covered or dropped is no class
change. A row that covers the answer of `success-criteria` is refused — the
question of the batch's form and the follow-up that asked it again alike: its
criteria are the code's rows, which the keep appends as `success-criterion`
rows after the session's.
An `R-NNN` is the position of its row in the draft it belongs to, never an
identity kept from the draft before it: a revision that reorders or drops
rows writes the IDs of the new order, and a record written before the
revision — a row of `intake/coverage.tsv`, an answer's text that names one —
names that earlier draft's row, the draft the revision session reads
([`records.md`](records.md#nfr-001--the-records-of-intake), the block `draft`).
**Done before** (task `T-8rm6`): the revised draft is
there when it has at least one row and every row's `session` is a session that
started on a records base holding `intake/answers.tsv`, and the step is done
before, with the session's comment reposted from the record when the control
issue lacks it, the batch's precedent. An empty draft (every fact of the problem
statement is `context`) holds no session row, so it is not revised until a
session row exists: the step runs its session and appends the criteria rows.
Known limit of task `T-8rm6` (rule: #409): an empty draft whose revision keeps no row
(no need, constraint or measure fact and no criterion) is not told from the draft
session's own empty draft, so the step runs its session again at each restart until
the milestone cap.
The revised draft, with its facts, is the input of the specification task and of the target's architecture task ([The instruction of a task step](#the-instruction-of-a-task-step)).

### The rendered commit and the confirmed inventory

At the first bet, code adds the **rendered commit** to the specification's head
(`architecture.md` §7 step 7). **Decided here:**

- The commit changes only the MoSCoW and Phase columns of the requirement tables,
  the non-goals and the numbered facts record, from the idea owner's answer;
  nothing else of the PRD or of the facts. **Decided here** (R10, task `T-3q0b`):
  the non-goals change is the bounded portion `Out of scope by the bet: <ids>.`
  appended to the `**Non-goals.**` paragraph, the ids sorted numbers first then
  the answer facts, and a fact's line is `<ID>. **<class>:** <text>` with `(out of
  scope)` in the class text for a mark of out. `layup/rules` passes it only when it
  equals code's rendering of the copy of the bet's answer, by [the basic
  rule](#layuprules).
- The **confirmed inventory** is what `architecture.md` §7 step 8 says: at each
  bet code writes a new version from that bet's copy, with the numbered facts, their
  classes and the out-of-scope marks (`spec/inventory-<n>.tsv`, the block
  `inventory`: one row for each fact, with its class and its mark), and the SHA-256
  of the version in `spec/inventory.tsv`; `n` counts from 1, and a later bet writes
  a new version, never a change of one that exists. **The numbered facts record** of
  the target, `docs/facts/`, is rendered by code from the version, with the
  confirmed classes: at the first bet on the specification's head, as the rendered
  commit above. The priorities (MoSCoW and Phase) are not in the version: they are
  the lines of the bet's copy, which the PRD's columns are rendered from and the rule
  `classes` of the full check compares. The specification's `layup/spec` check
  reads the hash.
- The specification task waits for its bet: its head fails `prd-lint` on the empty
  MoSCoW and Phase cells, which code expects: it runs the pinned `prd-lint` itself
  and accepts the head only when each failure line names an empty MoSCoW or Phase
  cell (`architecture.md` §8 step 5). **Decided here** (row 62c, #351): `prd-lint`
  runs by `sh <script> <dir>` and exit 0 accepts; exit 1 accepts only when at
  least one output line stands and each non-empty line of both streams is one of
  the two forms the pinned producer prints for an empty cell —
  `FAIL  <PRD-NNNN-kebab>.md: <REQ|NFR>-NNN MoSCoW "" not in Must|Should|Could|Won't`
  and `FAIL  <PRD-NNNN-kebab>.md: <REQ|NFR>-NNN has no phase value`. A `—`
  MoSCoW is not an empty cell; a line of any other form on either stream — a real
  violation, the filename refusal `err()` sends to stderr, a stray line — or an
  empty output refuses the head, and any other exit is an error, not a verdict.
- At a later bet, code renders the change as a rendered task of its own, with a
  plan ("render from copy N") that a plan-review and a verifier session check.

## REQ-012 — The full spec check

`layup spec check` is the command of the whole specification step
(`architecture.md` §7); `M2c` built the part on the numbered spans
([`intake.md`](intake.md#the-span-check)). **Decided here** (decision 9 of #333): the
full check adds these rules, in this order, each a row of the table
`tasks/<task>/spec-check.tsv` (the block `task-spec-check`: rule, result, detail),
one table for each head:

| Rule | Fails when |
| ---- | ---------- |
| `covers` | a `covers` ID does not resolve to a fact of the confirmed inventory or of the answers |
| `needs` | a `need` or `constraint` fact has no requirement, and it has no out-of-scope mark and no requirement marked "waits for a bet" |
| `measures` | a `measure` fact has no success-criterion row that cites its ID, and it has no out-of-scope mark and no requirement marked "waits for a bet" |
| `criteria` | a requirement has no criterion |
| `classes` | a MoSCoW or Phase value in the PRD differs from the latest copy of the bet |
| `inventory` | the SHA-256 of a version of the inventory is wrong |
| `quotes` | a requirement covers an answer's fact and not its quote's fact |
| `sections` | a delivered requirement has no non-empty `docs/spec/` section whose heading holds its ID |
| `tasks` | a task of the register names no requirement (the column `requirements` of `tasks.tsv` is never `—`) |
| `ran` | a rule did not run (`did not run`, never a pass, `NFR-004`) |

**Before the first bet** no inventory is confirmed (`architecture.md` §7, the end).
**Decided here:** a task that changes the PRD or `docs/spec/` is checked against the
numbered facts of Intake (`spec/spans.tsv` and the answers) by the rules `covers`,
`needs`, `measures`, `criteria`, `quotes`, `sections`, `tasks` and `ran`, and the
rules `classes` and `inventory`, which have nothing to judge against, are `clear` on
that table, each with the reason `before the first bet`; the status is then `success`
only when each other rule is `pass`. Any other task is `clear` for the whole table
with the reason `no specification yet`, and the status is `success`. After the first
bet each rule is `pass` or `fail`. The rows are committed as the table `tasks/<task>/spec-check.tsv`, with the head SHA
in the first column that `layup run` adds (the block `task-spec-check`); they are not
the table of `layup spec check --facts`, whose rule `classes` is another rule of
another block `spec-check` (`M2c`). The full check is a function of `internal/spec`
that `layup run` calls in process for the status `layup/spec`; **no command of the
binary runs it in `M2e`**: `layup spec check` keeps its flag `--facts`, and a CLI form
for the PRD, `docs/spec/` and the register is a later task. The rule `ran` makes a
missing row a failure. The function returns the table; it has no exit code, as no
command of the binary runs it (the status `layup/spec` is `success` when no row is
`fail`).

**Decided here** (row 62c, #351): the mode is source-bound — the caller's list of
the head's changed paths, repo-relative against the base, decides it before the
first bet: a head that touches `docs/prd/` or `docs/spec/` is judged, any other
head takes the whole-table `clear`; after the first bet all ten rules judge and
the changed paths decide nothing. A fact token is record-bound:
`F-<a problem-statement record>#n` resolves to the numbered fact `n` (LAYUP's
records are `F-0001` and `F-0003`, two records of the one file, so the input is
a list), `F-<the answers' record>#N` resolves to `Q-00N`, a bare `Q-NNN` is the
answers' fact itself, and a bare `F-NNNN` names a record, resolves to no
numbered fact and is no failure of `covers` — the pinned `prd-lint` already
requires the record file to exist; any other token fails `covers`. `measures`
reads only a `Fact` cell of a table under a `## … Success metrics` heading, never
an acceptance criterion's prose; the exemption of `needs` and `measures` is a
fact a token under a `waits for a bet` heading names — the heading text matches
case-insensitively at any heading level, the section runs to the next heading of
any level, and a token is one of the citation forms above
(`F-<record>#n`, `F-<the answers' record>#N`, `Q-NNN`), resolved record-bound —
or one the confirmed version marks out of scope. `classes` compares the PRD's values with the latest
copy's lines both directions, parsed by the one requirement-line grammar of the
bet's answer; a Phase cell `—` reads as empty, and a PRD row with an empty
MoSCoW and an empty or `—` Phase that the copy does not name waits for its bet
and agrees. The `inventory` rule verifies each index row's `sha256` against the
exact bytes the caller supplies and does not repair or infer one. `ran` fails a
missing row, and a rule that did not run — whose own row is a `fail` with the
detail `did not run: <what>` — fails `ran` as `did not run: <rule>` (NFR-004).

## REQ-012 — The parameter change

Every value that the Operator or the idea owner can set is a row of
`parameters.tsv` (the block `parameters`; `architecture.md` §10). **Decided here**
(decision 10 of #333, narrowed by R10 in task `T-dx1h`, #358): **the first writer** is
`layup run`, which writes the rows of [the table below](#the-parameters-of-m2e), each
with its default (a value that has none is `—`, and the parameters that Start sets,
`lease.H`, `watch.T` and the intake cap, take their values from `start.tsv`, the
evidence being the command: the word `command` and no comment; the band `B`, `U` and the
hours of the intent form are **not** parameter rows: `budget.tsv` is their one home
and a bet changes them as a rendered change, `architecture.md` §10) in the commit that writes
`transitions.tsv`, when the loop first runs: the commit of the first task, in which each
of `tasks.tsv`, `transitions.tsv` and `parameters.tsv` that is absent is written, one that
exists is read, checked and kept (never written again from the default), and one that its
reader refuses refuses the task and writes nothing. Each row has **one role**: the Operator
for most rows and for `brief.lines`, `build.parallel` and `transitions`; the idea
owner for the allowed dependencies. A change is a comment on the control issue (`layup run`
opens it at Start) in the form

```
set <name> <value> because <reason>
```

`layup run` copies the comment, and judges its first copy (an edit of a comment is never
a decision, and an App's comment is none, [`run.md`](run.md#a-human-decision)). A comment
whose first word is not `set` is no set and is ignored. A set is judged by five items, **the
form first**, then in this order (**decided here**, R10, task `T-dx1h`: `form` is a fifth item of
refusal, judged before the four that this section first gave):

0. the **form** (`form`): the body is one line, with its trailing white space trimmed,
   `set`, a name and a value of no space, `because`, and a reason that is not empty and
   not `—` (the writer of a table reads `—` as empty), each part after one space;
1. the **author** (`author`): an account whose role in `approvers.tsv` the row allows (the
   Operator for most rows; the idea owner for the allowed dependencies); before the name is
   known, any account of `approvers.tsv`;
2. the **name** (`name`): a row of `parameters.tsv`;
3. the **value** (`value`): of the row's type (and not `—`);
4. the **bound** (`bound`) of the row, which says what the value cannot do (no parameter
   removes a PSB rule, O-84).

The type and the bound of each name are held by the binary, in the table below in code;
the column `bound` of `parameters.tsv` is that bound written in words for a reader, and
the two are compared by the owner's test (`TestTheParameterTableEqualsTheSpecification`
reads the table below).

A change that passes is a new row of `parameters.tsv` and takes effect **at the
next step boundary**; a change that takes effect on an open task is input to that
task. The **boundary** is a call of `(*state).parametersBoundary` between two steps
of the loop: it copies the comments of the control issue, judges the first copy of
each in the order of the comment IDs, and writes each change as one records commit, through the fence of the
lease: `parameters.tsv` with the row replaced (its value, the reason as evidence, the
comment, the time) and, for `transitions`, `transitions.tsv` in the same commit. Until the
call, no read of the table sees the change. **The call is row 66k's** (`T-9nbf`, #367, the
step `loop` and the order of the steps; rows 66h and 66j read the table after it): row 66b
builds the function and shows it on the state of a restart, and no step calls it yet.
**A set is applied at most once and a superseded one never again** (the plan review of
#358, C2): it is applied only when its comment ID is greater than the `comment` that its
row records (`—` is none), the IDs of the comments of one forge growing with time; a set
at or below that comment is no change, no refusal and no report. A refused change is
recorded as human input, names the refused item, and changes nothing: in `M2e` the
record is the copy of the comment in `copies.tsv` (`human-inputs.tsv` is phase 3), and the
boundary reports each refusal (the comment, the item `form`, `author`, `name`, `value` or
`bound`, and the reason) as a progress line and in its result; it writes no row, no file and no
comment, so a refused set is judged again at each boundary, with the same result unless
what refused it changed (an account that gains its role, a name that a milestone adds),
and then it is applied if no later set of that row was. A set of `lease.H`, `watch.T` or the intake cap is judged like
any other and recorded, but nothing of `M2e` reads those rows from `parameters.tsv` (the
run reads `start.tsv`), so the boundary says `recorded; this run reads start.tsv, so M2e
still acts on the old value`. `architecture.md` §10 lets the idea owner change "the intent rows": in `M2e` they are not parameter rows (the band and the hours are `budget.tsv`, changed by a bet), so the idea owner's rows are the allowed dependencies. The parameters that `M2e` reads are `brief.lines`, `build.parallel`,
`transitions` and the allowed dependencies; the others are read by their
milestones.

### The parameters of M2e

**Decided here** (R10, task `T-dx1h`; the plan review of #358, C3). `architecture.md` §10
gives most of its rows no name, and no row other than these six has a complete source: a
name, a type and a bound (`stall.T`'s "10 minutes (a maximum)" is a direction, not a bound).
The table has exactly the six rows that a source gives completely, in this order, each
cell with its source in the last column; nothing in it is made up but the lower bound 1 of
two counts, which no source gives and which is the least that a use can hold.

| Name | Type | Bound | Role | Default | Source |
| ---- | ---- | ----- | ---- | ------- | ------ |
| `brief.lines` | `int` | a whole number of lines, 1 or more | `operator` | `—` | name: `architecture.md` §10 ("the brief's line limit"), The bet; type: a count of lines; bound: decided here, 1 the least that holds a brief of five headings; role: this section (the Operator for `brief.lines`); default: none, no source gives a value (the open gap of The bet: the Operator sets it, on the idea owner's evidence) |
| `build.parallel` | `int` | a whole number of tasks, 1 or more | `operator` | `—` | name: `architecture.md` §10, The task register; type: a count of tasks; bound: decided here, 1 the least that lets a task run; role: this section (the Operator for `build.parallel`); default: none, §10 says "set at Intake" and no step of Intake sets it |
| `transitions` | `sha256` | the SHA-256 of a copied comment of the control issue whose whole body is a table that passes the check of a table: a row for each default pair, only the kinds and items of the lists, and no `schema` in the row of code's own handoff | `operator` | `—` | name: The transition table; type: The transition table (the value is a SHA-256); bound: The transition table (a table that fails the check of a table is refused); role: this section (the Operator for `transitions`); default: none, the table that the binary embeds is in force while the value is `—` (decided here) |
| `lease.H` | `int` | whole minutes, 1 or more, as the flag `--lease-h` | `operator` | `—` | name: `records.md` (the names of `start.tsv`), `architecture.md` §10; type: the block `start`; bound: the flag's check (`run.md`, The command: whole minutes, 1 or more); role: this section (the Operator for most rows); default: none, the first value is `start.tsv`'s with the evidence `command` |
| `watch.T` | `int` | whole minutes, 1 or more, as the flag `--watch-t` | `operator` | `—` | name: `records.md` (the names of `start.tsv`), `architecture.md` §10; type: the block `start`; bound: the flag's check (`run.md`, The command: whole minutes, 1 or more); role: this section (the Operator for most rows); default: none, the first value is `start.tsv`'s with the evidence `command` |
| `intake.cap` | `MONEY,HOURS` | MONEY,HOURS, each a decimal with a fraction, as the flag `--intake-cap` | `operator` | `—` | name: `records.md` (the names of `start.tsv`), `architecture.md` §10; type: the block `start` (two decimals joined by a comma); bound: the flag's check (`run.md`, The command), no more; role: this section (the Operator for most rows); default: none, the first value is `start.tsv`'s with the evidence `command` |

**No row of the idea owner is written:** the allowed dependencies, which `architecture.md`
§10 gives to the idea owner, have no name in any source, so no row exists and no name is
made up; the role check of such a row is tested on a table that the test supplies. **The
other rows of §10** (the provider, the model version, the authority and threshold per
point, the escalation frequency, the owner map, the role matrix, the step table,
`stall.T`, `stall.N`, `stall.attempts`, `ci.T`, `panel.K`, the reward's term weights,
`learn.*`, the learning authority, `audit.n`, `harness.<id>.cap` and `harness.<id>.wall`,
the paid work of a harness with no spend cap) are read by their milestones, which write
each row with its type and bound; a `set` of one of them is refused as an unknown name
until then. The questions that remain (which of these rows are written now, and the name
of the idea owner's row) are asked on #358 and block no code.

**The owner checks** of `parameters` are in [`records.md`](records.md#nfr-001--the-records-of-the-task-loop):
a row has no comment when it is its own origin (a Start row with the value of `start.tsv`
and the evidence `command`; any other row unset) and a comment when a set gave it, the
reason being its evidence.

### The Retrospective

`architecture.md` §8 puts a Retrospective between the Accept and the next Bet. In
phases 2 and 3 **decided here** (open question 4 of the inventory) it is a step
with **no session**: code writes its event — a row of `retrospectives.tsv`, the
block `retrospectives` of [`records.md`](records.md#nfr-001--the-records-of-the-task-loop) — and goes to the next Bet. Proposals for
a rule batch that came from refused rule-path changes are collected for `M2f`; the
lessons and `layup learn` are phase 4. **Reason:** `architecture.md` §13 gives the
retrospective its content in phase 4, and an empty step that is a row keeps the
loop's order and no more.

## REQ-005 — The steps of the loop after the Scaffold

**Decided here.** `layup run` goes on, after the steps of the Scaffold
([`scaffold.md`](scaffold.md#the-steps-of-the-scaffold)), with the steps `shape`,
`bet` and `loop`; the block `run-steps` takes each name into its enum with its Go
schema in the build task, as `M2c` and `M2d` did.

- **`shape`** runs the Shape tasks that `M2e` can verify (the specification and the
  target's architecture; the activation batch is `M2f`'s and joins the step then)
  through the task loop of this file, up to the verified head of each, which
  `architecture.md` §8 asks before the bet; it also writes the
  options of each decision and asks the questions left open. It is **done** when each
  decision's table in `spec/options/` has a survivor and each Shape task has a
  verified head; a decision with no survivor leaves the step with the detail
  `shape: no survivor`, which is not "done before", so a restart stops there again
  and never goes on to `bet`.
- **`bet`** posts the brief, waits for the answer, copies it, writes the version of
  the inventory, pushes the rendered commit to the specification's branch (with its
  carry-over of `layup/verify`), opens the `architecture` tasks of the settled
  decisions (rows of `tasks.tsv`; each with `requirements` the requirement IDs among
  the `constraint` values of the chosen option's rows, and a decision whose table
  names none is refused, as an architecture decision serves a requirement), and writes
  the row of `bets.tsv` **last**. It is done when that row exists with its verdict, so
  a stop before it resumes inside `bet`, and no later step runs on a half-made bet.
- **`loop`** merges the Shape tasks in the merge order after a bet (the specification
  first), then runs the tasks that the register holds; the build tasks are `M2g`'s.

The record of each step: for `shape`, the files of `spec/options/`, the rows and the
`tasks/<task>/` files of its tasks; for `bet`, `bets.tsv`, the copies of the brief and
of the answer, and `spec/inventory.tsv`; for `loop`, the rows of `tasks.tsv`,
`transitions.tsv` and each `tasks/<task>/` file.

- **A restart resumes from the records and the forge, never from memory.** A step
  whose record exists is `done` with "done before", but `shape` and `loop`, which
  have their own rule above and below (a record of `shape` that holds a decision with
  no survivor, or of `loop` that holds an open task, is not "done before"). The `loop` step reads
  `tasks.tsv`, the events of each task and the forge (the pull request of the
  branch `task/<task>/<attempt>`, the statuses of its head) and goes on from the
  first step of the task loop (the table of [The task loop](#req-005--the-task-loop))
  that its records show not done.
- **A write is read before it is repeated.** Before `layup run` posts a comment,
  opens a pull request, posts a status or merges, it reads the forge: a comment whose
  SHA-256 is in `copies.tsv` is not posted again, an open pull request of the branch
  is not opened again, a status already at the state is not posted again, and a
  merged pull request is recorded with its `merged` event.
- **The branch of a session that has a `push` and no result.** A `push` event with
  no `bound` and no `refused` after it is a forge write whose result is not
  recorded ([`session.md`](session.md#the-push-and-the-bind)). **Decided here:** the
  `loop` step, before it acts on that session, reads the branch `task/<task>/<attempt>`
  of the forge: when it holds the SHA of the `push` event, the step writes the missing
  `bound` event; when it does not, it writes `refused` with the detail `push-lost`
  and starts a new attempt.
- **The comments of a session** ([`session.md`](session.md#a-comment-for-a-session)):
  from `M2e` the comment of a task session is posted on the task's issue and no
  longer on the control issue, with the same first line; the control issue keeps the
  comments of the probes and of the phase sessions.
- The lease, the fencing and the exit codes are those of [`run.md`](run.md); a lost
  lease stops the run before any write of a step.

## NFR-004 — A check that did not run

`NFR-004` (`F-0001#5`): a check that is not active does not count as passed. In
`M2e` this holds in four places, each tested:

- `layup/gates`: a `not-active` kind and a `pending` kind with a product path give
  `failure`; a `clear` is a pass for the reason it names.
- `layup/verify`: no harness outside the authors is `failure` with `not-active`.
- `layup/spec`: the rule `ran` of the full check.
- The merge: a context with no status at the head, a status of another source or a
  state other than `success` blocks the marking and the merge; the code counts only
  `success`.

## NFR-001 — The records

Git is the system of record (`F-0001#1`). The records of `M2e`, each a file of
the records branch that only the App may update (the rulesets of `M2d`), are in
[`records.md`](records.md#nfr-001--the-records-of-the-task-loop): the task
register, the transition table, the gates of each head, the frozen test list, the
spec check of each head, the options of a decision, the bets, the inventory
versions and the parameters. Each artifact that is a comment is a copy with its
hash. Nothing of the task loop is kept only on the forge.

## The acceptance tests of M2e

Each test names its level ([`test-levels.md`](../tests/test-levels.md)). The nine
rows of the inventory item `later-p2-phase-loop-handoffs` (six unit rows, one
integration row, one e2e row and one uat row; the table below splits the unit rows
into several) and the two of `later-gate-status` are the rows marked "inventory"; each is
kept at its level, and the e2e row of `later-gate-status` ("a pull request whose head
has a failed kind is not marked ready and does not merge") is kept as the integration
row "The ready and merge rule", because the stand-in forge shows it and the e2e row of
the milestone covers it end to end.

| Test | Level | What it shows |
| ---- | ----- | ------------- |
| A handoff is valid only by its table (inventory) | unit | A handoff with a wrong schema, an artifact that is not in the commits, a hash that differs, a required item that is false, and a pair with no row is each invalid, naming the item; a field that an agent writes never makes one valid; a comment artifact whose body at the forge differs from its copy is invalid. |
| The artifact that is a comment (`gov-issue-68`) | unit; integration | With the stand-in forge (the `httptest` forge for the integration part), for each of the kinds `plan`, `plan-review` and `review-record`: the session's file, the comment that code posts from its bytes, the copy and its row, the `handoff` event with the comment ID that code read, the three equal hashes; an edited comment is a new copy and the first handoff is invalid; a result that holds a comment ID is refused; a `comment` row is checked against the session's `result/` file, its file is kept byte for byte in the commit of the result, and an `artifact` row against the head; a run that stops after the end, between the post and the copy, or before the post, posts once from the records; a lost lease posts nothing. |
| The route of a question and the acceptance of an answer | unit; integration | A question goes to the owner role of its kind when its `needs_human` word is `no` and its label is one of the four, and to the Operator for `yes`, for any other word and for any other label; a replaced owner map routes by its roles; the items `question-row` and `answer-row` are computed from the result, whatever the caller passes; an answer is accepted by the next attempt that ends `completed` and cites its ID in the status row, with no question of that attempt citing it, and by no other attempt, and only when its session is of the owner role (the rule of "A question during the work": an answer of a developer, of the owner of another kind, of a session with no start row, to a human question or to questions of more than one owner role is refused, and the owner role's is accepted, from `sessions.tsv` and through `Position`); with the real `git`, a question session, the session of the owner role and the next attempt give the records from which the route and the acceptance are read. |
| The comment items of the task loop (row 66l, `T-4pje`) | unit; integration (the kept copy, the records and the forge) | `plan-fields`, `verdict-field`, `findings-located` and `nothing-material` are computed from the kept copy of a comment and, for the last, the verifier's recorded `base`: each item is true or false by its own rule in the table of the items (a plan-review copy that parses with `verdict: reject` is false), and a handoff listing the item is `valid` or `invalid <item>` by the computed result, never by a caller's literal of either value. |
| The tree items of the task loop (row 66n, `T-yeav`) | unit; integration (a real `git` tree) | `spec-sections` and `rendering-equal` are computed from the `docs/spec/` files or the tree at the head, never from the caller's literal of either value: a handoff of a developer (specification) or of the rendered task is `valid` or `invalid <item>` by the tree at the `bound` head or at the head that code sets, a head that a session handoff supplies is not read, an item that the row of the pair does not list is not computed, and a failed call of `git` or of the register is an error with no event. |
| The question and answer comments (row 66o, `T-gyvo`) | unit; integration (the committed rows of a real end) | The question and answer comments are composed from their rows, byte for byte, with the row's `n` in the ID, and refused with no bytes for a row or a start row off the reader's schema, types or rules, a row of another kind, an ID off its form, or a body that `CheckCommentBody` refuses; posting them is the step loop's (row 66k). |
| The transition table | unit | The default rows are the thirteen of the table (the twelve of `architecture.md` §8 and the row of the target's architecture task, which §8 names a task and gives no row); a replacement table is refused for a missing default pair, a kind or an item not in the lists, and a row of code's own handoff (`code (rendered task)`) that lists `schema`, while another row keeps its `schema`, and the steps that are `layup run`'s own acts need no row. |
| The replacement of the transition table at the boundary | integration | With the stand-in forge and real git: `set transitions <sha256>` names a copy of a comment of the control issue; the new table is the bytes of that copy, written with the row `transitions` in one records commit, so its SHA-256 is the row's value; the old table stays in Git history and the next task keeps the new one; a hash of no copy or of a copy of another issue, a body that is no table (a code fence, a carriage return) and a table that fails the check each refuse with the item `bound` (`value` for a value that is no SHA-256) and write nothing. |
| `layup spec check`, the full check (inventory) | unit | Each rule of the table fails on its input and passes on its fixture; before the first bet a task that touches neither the PRD nor `docs/spec/` is `clear`; `did not run` is a failure. |
| `layup/gates` (inventory, `later-gate-status`) | unit; integration | Each row combination maps to the status by the table, `not-active` is never a success, a `pending` kind with a product path is `failure`; with the stand-in forge, the status is posted once for an unchanged state. |
| `layup/rules`, the basic rule | unit | A change to a rule path fails; a rendered record equal to code's rendering passes; one that differs by a byte fails; the register is read at the base. |
| `layup/verify` and the independence | unit | The authors of a plan review and of a verification from the ledger; no harness outside them is `failure` `not-active`; a waiting verifier is `pending`. |
| The ready and merge rule (inventory, `later-gate-status`) | integration | With the stand-in forge, a pull request whose head has one context not `success`, no status, a status of another source, or an `error`, or a required job of GitHub Actions with no check run or a conclusion other than `success`, is not marked ready and does not merge; with the four `success` it is marked ready and merged at the head SHA, and a moved head is refused. |
| The close-out commit and the rendered commit | integration | A close-out that changes only the task file and the completed log, on a parent that carries the verdict, is accepted and `layup/verify` is set on the new head; one that changes another file, or whose parent does not carry the verdict, is refused; the rendered commit of the first bet gets the same carry-over, and the close-out after it is accepted. |
| The steps after the Scaffold | integration | A restart resumes the `loop` step from the records and the forge: a comment, a pull request, a status or a merge that exists is not made again; a `push` with no result is bound when the branch holds its SHA and `refused` with `push-lost` when it does not; a task session's comment is on the task's issue; a merge whose `merged` event was lost passes `forge-check` with its detail and `loop` records the event; `shape` with a decision that has no survivor is not "done before" and a restart stops there. |
| The bet step is atomic to its last act | integration | A stop after the inventory version and before the row of `bets.tsv` resumes inside `bet`: the rendered commit and each `architecture` task are not made twice, and the row is written last; a decision whose chosen option names no requirement is refused. |
| The tier and the context of a task (`session.md` deferral) | unit | The tier from the four fields of a plan (`large`, any non-empty list: reasoning; else execution), moved by the progress position (uphill to reasoning, downhill to execution); the prompt file lists the issue and plan, the rows and sections of the requirement IDs, the decision records they name, the handoffs and answers, and the rule files, in that order, each with its SHA-256. |
| The questions left open and the band (`intake.md` deferral) | integration | With the stand-in forge: an `open` and an `asked` question are asked of the owner role; an intent answer that parses (`B` below `U`) writes `budget.tsv` with the question as `source` and replaces a value that exists; one that does not parse leaves the question `open` and the run in its wait; a question still `open` at the brief is shown in "rabbit holes". |
| The approvers and the success criteria (`intake.md` deferral) | unit | The measure rows read from the intent answer; an approver row with `source` the answering comment ties to its point; a role `approver`; a point with no approver is the idea owner's; an account named by another spelling gets no second row; the revision's check that a requirement that covers an answer's fact names the quote's fact, and that a class change of a fact covered by both drafts names a fact of the answers. |
| The task register and the Retrospective | unit | A row with a predecessor that is not a task, a cycle of predecessors, or `—` in `requirements` is refused; the Retrospective step writes its event and no session starts. |
| The Shape filter (inventory) | unit | An option with a failed row, a missing row or a pass with no evidence is dropped; one survivor is chosen; two survivors give a brief with no recommendation, and an answer line `option <decision>: <option>` naming a dropped option or an unknown decision is refused and one naming a survivor is recorded in `bets.tsv`; every option dropped asks one rewrite, and a second full drop is a round with no progress. |
| The bet and its brief (inventory) | unit | The five headings, `brief.lines` unset is refused, a link that does not resolve is refused, two bets of 0.6 `U` each (the second brief), a cap above `band.U`, hours above `appetite.hours` and caps plus the spend of milestone 0 above `band.U` are each refused, and a cap below `band.B` passes, an answer by an account with no allowed role answers nothing, and the answer is copied. |
| The rendered commit (inventory) | unit | The commit may change only the MoSCoW and Phase columns, the non-goals and the numbered facts record; `layup/rules` passes a rendered record only when it equals code's rendering of the named copy. |
| The parameter change (inventory) | unit | A `set` is refused for a form that is not one line, a wrong role, an unknown name, a value of no type or out of bound, and the one that passes is applied as one replaced row, its default, bound and role kept; a set at or below the comment that its row records is superseded and is no change; a `prd-lint` head is accepted only when each failure line names an empty MoSCoW or Phase cell. |
| The parameter change at the boundary | integration | With the stand-in forge and real git: `layup run` writes `parameters.tsv` with its six rows in the commit of the first task (a table that exists is kept, one that its reader refuses refuses the task); a `set` is copied and, at the boundary, applied as one records commit when its author has the role of the row, and changes nothing before it; each refusal (`form`, `author`, `name`, `value`, `bound`) is reported with its item, changes nothing and is judged again at the next boundary; a set that is applied is not applied again, the later of two wins in the order of the comment IDs, and an edit of a comment is no set; a set of a Start row is recorded with the note that the run reads `start.tsv`; the write is fenced by the lease. The call between steps is row 66k's. |
| The frozen test list | unit | The list freezes at the first valid handoff to the verifier, a changed source hash is a finding, an amendment freezes it again, and a test that does not fail at the base makes the handoff invalid. |
| The gate jobs read the base (F-22) | integration | A pull request that edits `docs/gates.tsv` to pass a failing kind is judged by the base's file and fails; the catalog workflow reads the base; a pull request of a `batch/<point>` branch reads its own head's gate files and a pull request of any other branch does not; a change of a `config` path of the manifest waits for the code owner (shown on a real repository in the demo). The setup (task `T-jew3`, unit and integration): `.github/CODEOWNERS` of the four path classes with the Operator's login is in the tree of a new target, the ruleset sets `require_code_owner_review` with no required approval, and a target set up before the rule gets the update by its recorded ruleset ID, which keeps what the file never held ([`scaffold.md`](scaffold.md#the-code-owner-rule)). |
| The task loop with the stand-in forge (inventory) | integration | With the httptest forge, fake harness programs and the real git: a gate failure goes back to the developer as a new attempt; the close-out commit is refused when it changes a file other than the task file and the completed log; the pull request body links the records commit; the merge is at the head SHA. |
| The specification task, end to end (inventory) | e2e | In CI on a local fake forge, with scripted harnesses and scripted approver comments: Shape, the first bet, and the specification task merged with the four `layup/` statuses green. |
| The demo (inventory) | uat | On a real test target, the idea owner reads a bet brief of five parts within the line limit and answers it with one comment; the specification task merges after the first bet, with each handoff valid, the rendered commit holding the MoSCoW and Phase columns, and the four `layup/` checks green. |

### The demo and the release review

Defined here, as acceptance contracts; the rows, their IDs, issues and the inputs
that each asks are the slicing task's.

**The demo** (REQ-015) is the last row of the table: on a real GitHub test target
that `M2d` set up, a specification task goes through the loop, Shape runs with its
filter, the idea owner answers the brief with one comment, the rendered commit is
added, and the task merges at its head SHA with the four statuses `success`. As the
demo of `M2d` held five rules (`runs/T-sdtb/README.md`), this one passes when each of
five rules holds on the target and fails when one is false: (1) each handoff of the
task is a `handoff` event `valid`; (2) the rendered commit changes only the MoSCoW and
Phase columns, the non-goals and the facts record; (3) the four `layup/` statuses are
`success` at the merged head, read from the forge; (4) the merge is the recorded
`merged` event at that head SHA; (5) a pull request that weakens its own gate is
refused. The evidence is kept under `runs/` by the slicing task's demo row, the Operator
accepts the scenario in one word before the review, and a defect that the demo
finds is a task of its own. **The inputs** it needs from the Operator (for the
slicing task to ask): the commit-statuses, pull-requests and checks `read` permissions of the App, at least
two admitted harnesses (`architecture.md` §9), the confirmation of the reading of
`F-0003#58` (a pull request "reaches human review" when a review is requested or it is
marked ready), a target whose plan allows draft pull requests (GitHub refuses a draft
pull request in a private repository of the Free plan, so the target is public or on a
plan that allows it),
the `CODEOWNERS` apply, the value of `brief.lines`, and the idea owner's account.

**The review of the release** (REQ-015 and REQ-017) is a recorded code review of the
release of `M2e` by a reviewer of a model that wrote none of it, with the release
check of `runs/T-efmy/release-check.sh` on the release commit, before the demo, with
a coverage table of the files read that a script compares with `git ls-files`.

## Not in M2e

The batch part of `layup/rules` (an approved rule batch, its approval hash, the
known-bad patches, the activation of the `pending` kinds at the first bet, and
with it the item `patch-per-kind` of a rule-batch handoff) and
`layup audit` (`M2f`); the Operator's code-owner review of a rule batch that touches
the gate paths and the amendment of ADR-0017 and ADR-0019 decision 2 that records the use of code
owners (#336, `M2f` or earlier). A bet that raises `U` or the appetite, which `architecture.md` §10 makes a rendered change of `budget.tsv` (a later milestone; in `M2e` a brief that does not fit is reshaped), and a command form of the full check for the PRD, `docs/spec/` and the register (a later task). The escalation screen of an answer (`M3`) and the stall package, the stall
triggers and their procedure, which Shape's round with no progress and the loops of
failures wait for (`M3d`). The milestone plan, the build tasks and their merge order,
the execution of the named tests of a plan (the producer of `tests-red-green`), the
items `heading-ids` and `register-rows` of its handoffs, and Accept (`M2g`). The smart-if that
gives the kind and `needs_human` of a question (`M3a`). The panel of Shape and the
stall triggers (`M3`, `M3d`). The lessons and `layup learn`, the Retrospective's
content (phase 4). `layup report` and the Inter-Role Communication Format share
(`M4b`). The adoption of a target that phase 1 set up (none, O-200).
