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
| 4 | the draft pull request | `layup run` | the branch of the attempt | no rule-path change ([`architecture.md`](../architecture.md) §6, prevention layer 2); the title in the form of the `pr-title` check; the body links the issue and the records commit the task started from | [The pull request](#the-pull-request) |
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

**The opening of a Shape task** (**decided here**, R10, task `T-jy6f`, #415, row 66r; plan review conditions 2, 5 and 7).
`(*state).openShapeTask(ctx, kind)` opens the specification task or the target's architecture task, once, in this order, and
every refusal before the issue writes nothing:

1. `tasks.tsv` and `transitions.tsv` are read; a row of that kind is returned as it is (no issue, no draw, no commit), and a file
   its reader refuses is an error naming the path.
2. The permissions of the loop are checked before any write of the forge, the issue included (`forge.md`, the calls of M2e).
3. `spec/draft.tsv` is read by its reader, and the requirement IDs are derived from it (below); none is an error.
4. The candidate task is made with its drawn ID (`T-` and four characters that no row of `sessions.tsv` and no row of `tasks.tsv`
   holds), and the commit that will register it is built and refused in a dry run, with `registerFiles` and `parameterFiles`, the
   functions that `registerTask` calls, before the forge is written.
5. The lease is read, then the issue is opened with the fixed text below, then `registerTask` writes the row.

**Reason:** a refusal that a deterministic check can foresee must not leave an issue behind. **Known limit:** no call of the forge
lists issues, so a stop between the issue and the row opens a second issue at the next run; the first stays on the forge.
The first attempt is not started here: it is the first position of the driver.

**The requirement IDs of a Shape task** (**decided here**, R10; narrows "never changed" above to what the opening needs). The
correspondence of the draft's rows to the IDs, `REQ-NNN` and `NFR-NNN`, and the exact set of both rows are those of [The instruction of a
task step](#the-instruction-of-a-task-step), the one home of that rule; `finalIDs` is its code, and rows 66h and 66j use it. The
architecture task is not given the `non-functional` rows alone, since a draft with no constraint would give an empty list, which
`tasks.tsv` refuses, and the decisions of the options are not known at the opening. The `sections` rule counts a requirement as delivered
when every task that names it has merged, so naming it in both rows is safe. The issue names the SHA-256 of `spec/draft.tsv` and the
records commit it was read at, so the set is pinned. What the opening needs is that a later opening reuses the existing row and keeps its
requirement IDs; a later records commit may still replace the row to set `pr`.

The issue is titled `LAYUP task <ID>: the specification` or `LAYUP task <ID>: the target's architecture`, and its body is the
block of the kind with the placeholders filled by code (a test compares each constant with its block byte for byte). **Decided here:**
the Definition of Done is the text of the block itself, since no reusable text of it exists, and it holds no path of a session.

```text shape-issue-specification
LAYUP task {task}: the specification of the target.

Requirement IDs: {ids}
Pinned draft: spec/draft.tsv, SHA-256 {sha}, read at the records commit {records}.

Definition of Done: docs/prd/ holds a row of the requirement table, an acceptance row and a matrix row for each requirement ID, with the MoSCoW and Phase cells empty; docs/spec/ holds a section for each requirement ID; docs/tasks/{task}.md is the task's file. The head passes the four layup/ statuses, and a verifier on a harness outside the authors finds nothing material in scope.
```

```text shape-issue-architecture
LAYUP task {task}: the target's architecture.

Requirement IDs: {ids}
Pinned draft: spec/draft.tsv, SHA-256 {sha}, read at the records commit {records}.

Definition of Done: docs/spec/ and docs/adr/ hold each decision that has one survivor as decided, and a section "Decision pending the bet" for each decision that has more than one; docs/tasks/{task}.md is the task's file. The head passes the four layup/ statuses, and a verifier on a harness outside the authors finds nothing material in scope.
```

### The pull request

**Decided here.** `layup run` opens the pull request as a draft, with a title
of the form `<type>: <description>` (the form of the checked-in `pr-title` check,
decided below), a body that links the issue (`Closes #N`, or
`Refs #N` for a task that does not close it) and links the records commit that
the task started from, and requests no review. A draft is only visible; a
pull request "reaches human review" when a review is requested or it is marked
ready (the reading of `F-0003#58`, which the approval brief asks the Operator to
confirm). The review record of the verifier is posted on the **issue**, and the
body of the pull request is edited to link it (the edit: row 66q, task `T-08rj`, [below](#the-link-of-the-review-record-in-the-body)), which makes the target's own
`review-record` check run again. **Reason:** the target's checks parse the
issue's comments and the body, so the loop writes them in the forms that those
checks read.

**Decided here** (R10, task `T-xvvd`, #363, row 66g; plan review conditions 1 and 2; `OpenPullRequest`
of `internal/run`):

- **The inputs are from the records.** The head is the SHA of the last event `bound` of the attempt that has a session; the head that the pull
  request must have is the **tip** of the branch, `TipOf` (the last `bound` of any session): the same, or a merge of the base that code made above it
  (task `T-89xf`: every commit from the tip down is the commit of an event `bound` with no session and a merge whose first parent is the
  next, down to the head or to the rendered commit of the first bet, the one commit of code with one parent, which sits directly on the head), so a pull request carried over a merge is adopted at the merge, again and again, with no draft and no change of the register, and any other
  move of the branch is the error below;
  the session of that event gives the start row, and **the records commit that the task started
  from** is its `records` column (the snapshot of the push check and of `layup/rules`), never the
  run's changing base and never a later session's. The branch is `task/<task>/<attempt>`.
- **Read before write.** The pull request of the branch is read first: an open one is used, a merged
  one is returned untouched, a closed one that is not merged is a Go error (the wait and the redo are `BringUpToDate`'s, row 66x). The
  head of the open or the new draft must be the bound SHA: if the branch moved after the event `bound`,
  the call is a Go error (the head moved) and registers nothing (a draft just opened stays, as above).
  Only a branch with none is opened, under the lease. The number is then written by one fenced
  records commit that replaces the task's row in `tasks.tsv` with `pr` set; a stop after the draft
  opened leaves a draft, which the next call finds by its branch.
- **The title** is the form of the checked-in `pr-title` workflow, the only check of a title that
  exists: `<type>: <description>`, `<type>` one of `feat`, `fix`, `docs`, `refactor`, `perf`, `test`,
  `build`, `ci`, `chore`, `revert`, and a description of one line that matches `^[a-z].+[^.]$` with no
  control character. **The task ID is not in the title** (it begins with a capital, which the check
  refuses, and the workflow keeps the ID in the commit subjects, with the issue reference and the records link in the body); this replaces the
  earlier `<type>: <task ID> <description>` read from a setup record that does not exist. The caller
  (rows 66h, 66k) gives the two words; a title off the form is refused before any call.
- **The body** is `Closes #<issue>`, a blank line and `The records commit this task started from:
  <web>/<owner>/<name>/commit/<records>`, with no final line feed (`<web>` is the forge register's
  `web`). `Closes #N` comes first, so the issue that the target's `review-record` job reads is the
  task's. Every M2e task has its own issue, so M2e writes `Closes`. The composer takes further lines.
- **The edit of the body to link the review record is row 66q's** (task `T-08rj`, #413), decided
  [below](#the-link-of-the-review-record-in-the-body). It uses the seventh method of `internal/forge`
  ([the calls of M2e](forge.md#the-calls-of-m2e)) and this task's composer, whose further lines it fills.
  The step that calls it is row 66k's, after row 66h.

#### The link of the review record in the body

**Decided here** (R10, task `T-08rj`, #413, row 66q; plan review conditions 1 to 6;
`(*Sessions).LinkReviewRecord(ctx, task, verifier)` of `internal/run`, which gives whether it edited):

- **The body** is the composer's body with its further lines `""` and `The review record:
  <web>/<owner>/<name>/issues/<issue>#issuecomment-<id>`, and no final line feed: the lines of the draft
  unchanged, then the blank line, then the link (the form of `setup.go`'s link of a comment). `Closes #N`
  stays the first line, so the issue that the target's `review-record` job reads is the task's, and the
  job runs again on the event `edited`. **Reason:** one composer, so the draft and the edited body cannot drift.
- **The inputs are from the records, never from the caller.** The verifier's start row is the session's
  (its task, its role the one of the step `verification`, its attempt the newest `attempt` event of the
  task). **The comment** is the one of the **latest** event `handoff` of that session, which must be
  `valid comment <ID>-<seen>`: the latest is chosen first and then required valid, so an earlier valid
  handoff is never used after a later invalid one (a different outcome of the same session is a new event, and
  an edit of the comment can make it invalid). A valid handoff to the developer, for a material review, links as
  well as a handoff to the close-out: both are valid review records. The copy is the row of `copies.tsv` with that
  ID, that `seen` and the task's issue; its file must parse as a review record whose `head` is the verifier's
  `base`, which tells it from a plan review of the same role. **Reason:** the link fails closed (NFR-004);
  whether a record's verdict allows a merge is not decided here.
- **The records commit of the body is the one of the opening, read from the body.** `OpenPullRequest` writes the
  records commit of the last bound session into the draft's body, and no later bind changes that line: a verifier or
  a close-out that leaves the head unchanged binds the same SHA again at a snapshot of its own, so no choice among the
  events (the first, the last, the last before the verifier) names the opening for every verifier, a second one and a
  restart included. The edit therefore keeps the commit that the body it replaces links (`bodyRecords`, the inverse
  of `PullBody`), provided that some session which bound the reviewed head in this attempt started at that commit. A
  body with no such line, or a commit that no such session holds, is a Go error and nothing is edited: the call never
  publishes a snapshot that it chose itself. **Reason:** the body is the one stable record of the snapshot of the
  opening (a new record or event was not needed), and the check against the sessions of the head keeps a mistyped or
  invented commit from being published. A person's edit that keeps the line is replaced whole; one that removes it is
  refused. **Known limit** (ruled by the Operator, #413): a person's edit that names another genuine records commit of a
  session that bound the same reviewed head, or that puts such a line before the opening one, is accepted and that commit
  is published; it is still a real commit of a session of the head. A durable association (the snapshot written at the
  registration of the pull request, in row 66g's code) is not built in M2e.
- **The pull request is the attempt's.** `Task.PR` has no attempt, and a new attempt leaves it on the older
  pull request until `OpenPullRequest` replaces it. The call reads the pull request of the branch
  `task/<task>/<attempt>`, and edits only when it is found, its number is the register's `pr` and its head SHA
  is the verifier's `base`: a newer attempt that exists while the register names the older pull request is a
  Go error, and the older pull request is never edited. A pull request that is not `open` is a Go error naming
  the state (what a closed one means is `BringUpToDate`'s, row 66x).
- **Read before write.** A body equal to the composed one is no edit (`false`, no call). Any other body is
  replaced whole, also one that a person edited and that keeps the line of the records commit: the body is code's
  artifact in the form that the target's checks read. **The lease** is read immediately before the edit, `*LostError` when it is not held, and a failed
  read is an error; as for a post, a takeover between the read and the edit remains, and the edit is harmless
  there, because both runs write the same bytes from the same records.
- **The answer is read, never trusted.** The body that the forge gives back must be the composed one, else a
  Go error. An edit that took effect and gave an error is repeated by no one: the next call reads the body,
  finds it equal and makes no second edit.
- **No record is written**, on a retry either: the body is a function of records that exist, so the edit
  decides nothing (NFR-001).

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
  **Later wiring:** the rows of `gates.tsv` that keep the table across a restart are
  written and read by task `T-q0kr` (#362: `GatesTable` gives the `gate.Table` back, from
  which `GateFindings` is computed); the call that notices the failure at step 5 and
  resumes the loop from it stays #367's, and until then the caller gives the table.
- **Decided here** (R10, task `T-tnxn`, #427, row 66aa; plan review conditions 1 to 7 of 6262459454): the failures of
  `layup/spec` and `layup/rules` start a new attempt as the failure of `layup/gates` does. Three classes, each
  able to fail while the others pass: the spec findings, the rules finding and the expected-failure exclusion.
  - **The spec findings** are the `fail` rows of the head's table (`tasks/<task>/spec-check.tsv`, read by `SpecTable`),
    in the order of the rules, the row `ran` included (a failed `ran` is a recorded failure, not a status error). Each has the
    rule `spec:<rule>`, the row's detail as its text, and no file and no line. A `clear` row is none.
  - **The rules finding** is the result of the basic rule when its state is `failure`: the rule `rules:<reason>`, the
    path as its file with line 0, and the description of `RulesStatus` as its text (the same words as the status; the
    complete path stays in `File`, as the description can be cut). `success` gives none.
  - **Errors come first.** An `error` of any of the four statuses, `not-active` of `layup/verify`, and a Go error
    are decided before any retry, and no attempt starts even when another check has findings. `CheckFindings` gives
    a stop for a head whose `layup/spec` has no row (the status `error` that `SpecAt` posts for a head whose own tree
    is refused, which writes none) and for a `layup/rules` of state `error`; a table that has rows but not the ten rules
    in order is refused by its reader: a Go error, never a stop or a finding.
  - **The inputs of a restart.** The findings are rebuilt from the records, and no record is added: the spec rows
    from `spec-check.tsv`; the rules result by `rules.Basic` over the attempt's base and the head, with the register
    commit that `registerRevOf` names, which is the `records` column of the start row of the session bound to the head, never
    the run's current records commit. A Shape head has no rendering input (`PreRev` and `Records` are empty); a later
    caller that gives either must take it from the same recorded source on a restart.
  - **One attempt for all failures.** `NextAttempt` takes the gate table and the given findings, the gates first,
    then the spec, then the rules (the order of the contexts), and `ErrNoFailure` is returned only when both are empty. The attempt starts
    at the base of the attempt that failed.
  - **The expected failures of `architecture.md` §8 step 5** start no attempt. (a) **A rule batch before its approval**,
    in `M2e` (no approval exists before `M2f`): for a task of the register kind `batch`, a failure of `layup/rules` is no
    finding and the rows of `layup/gates` whose state is `pending` are none; the rows of an `active` kind stay findings,
    and so does every row of any other kind. Its evidence is the scratch run of `M2f`. The exclusion is made at the retry
    boundary, over `GateFindings`' translation, which is not repeated. (b) **A specification that waits for its bet**:
    its failure is the pinned `prd-lint`, which no `layup/` status runs. Before the first bet `layup/spec` clears `classes` and
    `inventory`, and the PRD is no rule path in phase 1 (see `layup/rules` below), so that failure yields no finding by
    construction. Its acceptance before the verifier is the driver's ([The driver of a Shape task](#the-driver-of-a-shape-task)),
    the Operator's ruling of 2026-10-19 on #427; the native job's readiness is the close-out's (row 66y).
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
  The loops of the failures of `layup/spec` and `layup/rules` are covered by the same limit (decided here, task `T-tnxn`):
  they are bounded only by the budget cap and by the acts of one call of the driver, and count no review round.

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
after `not-active: `, and `route.ErrNoPair`. **Decided here** (task `T-q0kr`, #362): `layup/verify` is posted `failure` at the
head with that description, and never `pending` or `success`.

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
  handoff](#the-checks-of-a-handoff), "A question during the work"). The open escalations (`M3`) are
  added when they exist. **Decided here** (R10, task `T-r7b9`, #426, row 66z; plan review 6262463499, condition 5):
  an open material finding keeps the task uphill: `Position` is `Downhill` only when the list is frozen, there is no
  open question **and** `openMaterial` is false. `openMaterial` selects the **last verifier** by the order of the events of its
  start, and reads the output of that one only after the selection (round 1 of `T-r7b9`, finding 1, comment 6289084862; plan review
  6262463499, condition 5). A verifier is told by its position and its head binding, the classification of `loopView.qaKind`
  (the one of `layup/verify`, "The independence of the reviewers"): a session of the role `qa-engineer` that started after an event
  `bound` of a developer's session of its own attempt, at the SHA that is the base of its start row. A plan review has no such
  `bound`, and a session of the role that is bound to no head is no verifier, so neither is ever read as a review record. With **no**
  verifier nothing is open. The selected verifier clears a finding only when its output is qualified: it ended `done` with the status
  `completed`, its result holds a `comment` row of `result/review-record.md`, and its handoff is recorded, `valid` and names a comment
  (a verifier that is still running, ended `needs_context` or any other status, whose result has no such row, or whose handoff is not yet
  recorded, is open); its kept copy exists, has the SHA-256 that the result names (`keptComment`), parses (`ParseReviewRecord`) and has a
  `head:` that is the base of the session's start row; otherwise it is open; and then `OpenMaterial` of the parsed record decides.
  The output of the latest verifier that is missing or invalid therefore never clears a finding of an earlier record, and no older clean
  record is used in its place. A failed read is an error, not an answer.
  **Known limit (decided on #426, comment 6289096719, finding 2 of round 1):** every step of a Shape task has the fixed tier
  `reasoning` in the step table, so the driver of a Shape task asks no tier of `TierOf`, and no production code yet calls `Position`
  and then `TierOf`. The production consumer of `Position` → `TierOf`, and the semantic test that fails when that consumer ignores an
  open material finding, belong to **row 66k** (#367, `T-9nbf`), where the routed tiers of the step loop first appear; this task
  tests `Position` and the tier it moves through the production constructor (I5) and releases the caller-wiring clause of condition 5 by
  that decision alone.
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

- **The frame.** `Text` keeps the header line, the sentence that tells the harness to read, and the five lists, in their order, each existing item in its place relative to the others; the new items below join their own kind. After the lists, code writes a blank line, the line `The work of this step:`, the block of the step with its placeholders filled, a blank line and the block `instruction-result` filled; **decided here** (task `T-r7b9`, row 66z) code then writes, only when the attempt before the open one has an accepted-eligible answer (the answer session of the owner role of [A question during the work](#the-checks-of-a-handoff), ended `done` with the status `completed`, with its answers), a blank line and the block `instruction-answers` filled, so the golden file of every other prompt is unchanged byte for byte. A target item or a records item that is new makes the bytes of a list longer, so the golden files, not a byte prefix, are the oracle of the composition.
- **The home of each text.** The blocks below (the seven steps with a fixed instruction, `answer-to-a-question`, the shared `instruction-result` and the conditional `instruction-answers`) are the text of record, with one home each. `internal/run` keeps each as a constant and a test compares it with its block byte for byte (as it does for `probe-prompt`); the golden files of the prompt of each step compare the composed bytes, and no flag rewrites a golden file.
- **The placeholders** are filled by code and never by a session; a prompt that keeps a `{` placeholder is a defect. `{task}`, `{kind}` and `{issue}` are the row of the task in `tasks.tsv`. `{result}` is the absolute `<session directory>/result/result.tsv`, and `{comment}` the absolute `<session directory>/result/<file>`, the file being `plan.md`, `plan-review.md` or `review-record.md` (the constants `PlanFile`, `PlanReviewFile` and `ReviewRecordFile`, which row 66r's driver passes as the file of the handoff). `{head}` is the base of the session's start row, the head that the item `nothing-material` reads. `{from}` is the base of the task's last `attempt` event. `{labels}` are the keys of the default owner map, sorted, joined by `, `. `{record}` and `{facts}` (the specification only) are the highest `F-NNNN` of `docs/facts/` at the base and the next number after it. `{questions}` (`answer-to-a-question` only) are the IDs `<session>.<n>` of the question rows of the sessions that asked in the open attempt, in the order of [A question during the work](#the-checks-of-a-handoff), joined by one space; `{answers}` (`instruction-answers` only) are the IDs of the answer rows of the answer session of the attempt before, in the order of their rows, joined by one space.
- **The steps.** A step has an instruction exactly when this table has a row; `shape-activation-batch` and `plan-of-a-milestone` are the milestones M2f and M2g: `PromptFile` refuses the two by name, before every read, with `the step <step> has no instruction in M2e`. **Decided here** (R10, task `T-r7b9`, #426, row 66z): `answer-to-a-question` has an instruction, the block `instruction-answer-to-a-question` below, and it leaves the refused steps; its posting and its place in the loop are in [A question during the work](#the-checks-of-a-handoff), and its tier is the step table's `reasoning`.

| Step | Role | Kind of the task | Writes in the repository | Comment file, read by | Required outputs for the result |
| ---- | ---- | ---------------- | ------------------------ | --------------------- | ------------------------------- |
| `plan-of-a-task` | `software-engineer` | any | nothing, no commit | `plan.md`, `ParsePlanComment` | one `comment` row |
| `plan-review` | `qa-engineer` | any | nothing, no commit | `plan-review.md`, `ParsePlanReview` | one `comment` row |
| `specification` | `product-owner` | `specification` | `docs/prd/`, `docs/spec/`, the task file, commits | none | the task file and each file written |
| `shape-architecture` | `systems-architect` | `architecture` | `docs/spec/`, `docs/adr/`, the task file, commits | none | the task file and each file written |
| `implement` | `software-developer` | `build` | tests and change, the task file, commits | none | the task file and each file written |
| `verification` | `qa-engineer` | any | nothing, no commit | `review-record.md`, `ParseReviewRecord` and `NothingMaterialOK` | one `comment` row |
| `close-out` | `software-developer` | any | one commit of the task file and `completed.md` | none | the task file and `completed.md` |
| `answer-to-a-question` | the owner role of the question | any | nothing, no commit | none | none: the `answer` rows |

The step `specification` takes a task of kind `specification`, `shape-architecture` one of kind `architecture` and `implement` one of kind `build`; any other step takes any kind, and a pair that disagrees is an error that names the step, the task and the kind. A comment file is written in the pinned form of [The items a caller computes](#the-items-a-caller-computes) (the target's own forms, never LAYUP's `| Field | Value |` table), and the instruction shows it with its placeholders so that a test fills them and runs the real parser on it. The shared text names every rule of the block `result` that a model must meet: the actual tab, UTF-8 with no byte-order mark, the line feed at the end of each line, no blank row, the empty cell, the numbering without gap, a status row with a non-empty reason, a hash of 64 lowercase hexadecimal characters, a question that names a record, and the artifact rows for the required outputs only.

**Decided here** (task `T-4naq`, #428, row 66ab; [the specification task that a changed class reopens](#the-specification-task-that-a-changed-class-reopens), item 9): the step `specification` of a task
that has a row in `spec/reopenings.tsv` has the block below in place of `instruction-specification`, with `{fact}`, `{class}`, `{token}`, `{requirement}` and `{record}` filled by code from the row and from the tree at the base. The oracle of the
composed bytes is a golden file for a `need` and one for a `measure`, and the test fills the placeholders of the lines it shows and runs the real full check and the pinned `prd-lint` on the result.

```text instruction-specification-reopened
Write the requirement that waits for a bet, for the task {task} (kind {kind}, issue {issue}), from the files listed above: the row of this task in spec/reopenings.tsv, the facts spec/spans.tsv, start/problem-statement.md and intake/answers.tsv, the revised draft spec/draft.tsv, the product requirements document in docs/prd/ and the specification in docs/spec/. The bet changed the class of the fact {fact} to {class}, and no requirement names it yet; the token {token} cites it. Commit your work in this repository in one or more commits; push nothing.
Keep every line of the files that already exist, except the additions below, renumber nothing, and put no new row in any table. In the file of docs/prd/ write, under a heading "## Waits for a bet" (add it as the last heading of the file when the file has none), the line of your class, as plain text starting at the first column, never in a table:
- need or constraint: "{requirement}: <the requirement, one sentence> Facts: {token} {record}. Acceptance criterion: <the acceptance criterion>." In docs/spec/ write a section whose heading holds {requirement} as a whole word, for example "## {requirement} — <title>", with a non-empty text under it.
- measure: "Success criterion: <the success criterion> Fact: {token}." Write no section in docs/spec/.
The class of this task is {class}: write the line of that class only.
Write docs/tasks/{task}.md, the task's file, in the form of the target's docs/tasks/. Write nothing under docs/facts/.
The required outputs for the result file are docs/tasks/{task}.md, the file of docs/prd/ that you wrote and, for a need or a constraint, each file of docs/spec/ that you wrote.
```

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

## 6. Functional requirements
| REQ | Statement | MoSCoW | Phase | Facts |
| --- | --------- | ------ | ----- | ----- |
| REQ-<NNN> | <the requirement> |  |  | F-{facts}#<n> F-{record} |

| Requirement | Acceptance criterion |
| ----------- | -------------------- |
| `REQ-<NNN>` | <the acceptance criterion> |

## 8. Success metrics
| Metric | Fact |
| ------ | ---- |
| <the success criterion> | F-{facts}#<n> |

| REQ | Facts | Guardrail | ADR | Task | Test |
| --- | ----- | --------- | --- | ---- | ---- |
| REQ-<NNN> | F-{facts}#<n> F-{record} | — | — | — | — |

## 9. Rollout & phases
| Phase | Name | Ships |
| ----- | ---- | ----- |
| <the phase tag> | <the name of the phase> | <what ships in it> |

Put the row of a requirement of kind non-functional, written NFR-<NNN> in place of REQ-<NNN>, in the same form under the heading "## 7. Non-functional requirements", and write a paragraph that starts with **Non-goals.** (the non-goals of the target, or none), to which code appends the facts that the bet marks out of scope. Keep the backticks around the ID in the acceptance table. Leave the MoSCoW and Phase cells empty (not —): code fills them at the bet. In the phases table write one row for each phase of the target (at least one): its tag, the tag that the bet's answer gives a requirement, its name and what ships in it. In a Facts cell, cite each fact that the row covers: F-{facts}#<n> for the fact n of spec/spans.tsv, and Q-<NNN> for a fact of intake/answers.tsv; always cite F-{record}, the existing record, too. Write one acceptance row and one matrix row for each requirement ID, and one success row for each success-criterion row of the draft, citing the facts it covers.
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
The required outputs for the result file are docs/tasks/{task}.md and docs/tasks/completed.md.
```

```text instruction-answer-to-a-question
Answer the questions that the sessions of the last attempt of the task {task} (kind {kind}, issue {issue}) asked: {questions}. The question rows are in the result files of those sessions, listed above under handoffs-and-answers; the IDs above are in the order of the questions, each a session and the n of its row. Read the files each question concerns, the requirements, the decision records and the rule files. Change no file of the repository and make no commit.
Write one row of the kind answer for each question, and no other row of that kind: "answer<TAB><n><TAB><the answer, one line><TAB>—<TAB><record>". The answer row of column n = 1 answers the first question of the list above, the row of column n = 2 the second, and so on. At least one <record> names a file the answer rests on, more are separated by one space, and there is no space at either end of the cell. Write the status row completed when every question is answered.
```

```text instruction-result
Write the result file {result}: UTF-8 text with no byte-order mark, one row per line, each line ended by a line feed, the last too, and no blank line. Separate the cells of a row with one tab character each, never with spaces. The first line is "kind<TAB>n<TAB>value<TAB>sha256<TAB>reason". The single character — (U+2014) is the empty cell; no cell holds a tab or a line break. The rows of each kind are numbered 1, 2, 3 and so on, with no gap.
Write exactly one status row, "status<TAB>1<TAB><status><TAB>—<TAB><reason>": <status> is completed when the work of this step is done, needs_context when a question stops it, and blocked or failed otherwise; <reason> is one line and not empty.
For each required output that the work above names, write the row "artifact<TAB><n><TAB><path><TAB><sha256><TAB>—": <path> is relative to the repository, and <sha256> is 64 lowercase hexadecimal characters, the SHA-256 of that file at your last commit. A file that this prompt only asks you to read is no output.
When a question stops the work, write one row for each question, "question<TAB><n><TAB><text><TAB>—<TAB><label> <needs_human> <record>": <label> is one of {labels}; <needs_human> is yes or no; at least one <record> names a file the question concerns, more are separated by one space, and there is no space at either end of the cell.
A comment file is plain Markdown with no code fence. Push nothing. Then end.
```

```text instruction-answers
The sessions of the last attempt asked questions, and the owner role answered them; the answers are listed above under handoffs-and-answers. Their IDs are {answers}. Write each of these IDs as a word of the reason of the status row, and write none of them in a question row.
```

**The records each step receives.** Every new item is a byte copy under `input/` with its SHA-256, inside one of the five kinds (no sixth kind: `ContextKinds` does not change), after the items that exist; a file is read by its reader, which fails closed.

| Input | Records path | Kind, after | Steps | Reader |
| ----- | ------------ | ----------- | ----- | ------ |
| the frozen list | `tasks/<task>/tests.tsv`, when it exists | `issue-and-plan`, after the plan | all | `ReadTests`; absent is no item |
| the revised draft and its facts | `spec/draft.tsv`, `spec/spans.tsv`, `start/problem-statement.md`, `intake/answers.tsv`, in this order | `requirements`, after the target items | all but `close-out`, for a task of kind `specification` or `architecture` | `ReadDraft`, `ReadSpans`, `ReadAnswers`; the brief as bytes; a missing file is an error |
| the options | `spec/options/<decision>.tsv` for each decision of the caller's map, sorted | `decision-records`, after the ADRs | all but `close-out`, for a task of kind `architecture` | `spec.ReadOptions` |
| the survivors | text: `the survivors of the filter:` and a line `- <decision>: <option> ...` for each decision, under the decision records | `decision-records` | as the options | the caller's map, checked against the tables |
| the ledger | `sessions.tsv` and `telemetry.tsv` | `handoffs-and-answers`, after the answers | `close-out` | `ReadSessions`, `ReadTelemetry`; both tables whole (a known limit: code keeps no table of one task) |
| the asking results | `tasks/<task>/results/<session>.tsv` of each asking session of the open attempt (**decided here**, task `T-r7b9`) | `handoffs-and-answers`, after the answers | `answer-to-a-question` | `records.ReadResult` |

The Intake records go to both Shape kinds because the `covers` column of the draft names span numbers and `Q-NNN` facts, and the constraints of an options table name `NFR-NNN` (the final ID of a draft row `R-NNN`, see the requirement IDs of a Shape task), `F-<n>` and `Q-NNN`; neither can be read without the four files. **The options (R10, binds 66h, 66r and 66j).** `PromptIn.Options` maps each decision to the `spec.Result` of its table and is required, and not empty, for a step that takes the options (an architecture task, any step but `close-out`), and refused for any other. *Before the bet* the caller (row 66h's filter, passed by row 66r's driver) gives each decision's result unchanged: every option of the table and every survivor, the survivors listed in the prompt and none chosen. *After the bet* the caller (row 66j opens the task of the chosen option; row 66r passes it) gives the same options and only the option that the bet chose as the one survivor; that narrowed choice is no filter result. `PromptFile` reads each table and refuses a result whose option names are not the table's, in order, a result with no survivor and a survivor that is not an option; whether the set of decisions is complete is the caller's, as `SessionStore` has no call that lists a directory. **The base names.** The existing refusal of two records of one base name applies to the records that one prompt holds, and names both; a decision named `draft`, `spans`, `answers`, `tests` or `tasks` collides with a record of its own prompt and is refused. That no real run meets it is a contract of row 66h, whose step must refuse those decision names (`optionsKeep` is not on this main: a required contract, not a reuse).

**The requirement IDs (R10, binds 66r; condition 3 of the plan review of #364, comment 6249724875).** The correspondence is by kind, from the pinned revised draft: the row `R-NNN` of kind `requirement` is `REQ-NNN` and the row of kind `non-functional` is `NFR-NNN`, the same three digits; a row of kind `success-criterion` is no requirement ID. The `requirements` of **both** initial rows of `tasks.tsv` (the specification task and the target's architecture task) are exactly the set of those IDs over all the rows of the pinned draft, in the order of the draft. The number a row has in the pinned draft is the number of its PRD identity: a later revision of the draft changes no ID that a PRD row holds, and the PRD, not the draft, is the identity from the merge. Task `T-jy6f` (row 66r) writes the rows by this correspondence and tests the item `spec-sections` for them, with a missing-section refusal; task `T-c6sx` (row 66s) wrote the sentence and opened no task. The later architecture tasks of row 66j, and the option constraints, use the same IDs.

**The specification output (R10, binds 66j and 66k).** The instruction asks for what the existing checks read, so a head that waits for its bet is accepted: one five-cell row for each requirement ID with empty MoSCoW and Phase cells, one acceptance row, one matrix row and one section of `docs/spec/` for each ID, one success row for each `success-criterion` row of the draft, and Facts cells whose tokens the checks resolve. A fact is cited as `F-<record>#<n>`, the record being the numbered facts record of the problem statement, `docs/facts/F-NNNN-numbered-facts.md` of the decision of #412 (`T-cav3`), and `n` the number of the fact in `spec/spans.tsv`, which `RenderFacts` writes as the line `n.`; a fact of the answers is cited as `Q-NNN`. The record does not exist until the rendered commit of the first bet (row 66j), and its `NNNN` is one more than the highest `F-NNNN` of `docs/facts/` at the base of that commit, which equals the base of the specification task; so the prompt gives that number as `{facts}`, and, as the pinned `prd-lint` resolves a cell only when one of its tokens names an existing record, the cell also cites `{record}`. The consumers bind to it: row 66j renders the record under `{facts}` with the same function, and the input of the full check gives `PSBRecords` the predicted ID before the bet and `Records` the records of `docs/facts/`, so a bare `{record}` and `F-{facts}#n` both resolve ([The full spec check](#req-012--the-full-spec-check)). The specimen test builds its input through `specRecords` and `specTree` and runs the full check and the pinned `prd-lint` on it. The rendering stays row 66j's.

**The errors.** `PromptFile` writes nothing at any stage, and every error returns no context, so no session starts on a partial prompt.

| Stage | Case | Error |
| ----- | ---- | ----- |
| 0, before every read | a step with no row in the step table; a step with no instruction | `the step table has no step "<s>"`; `the step <s> has no instruction in M2e` |
| 1 | no row of the task; the step and the kind disagree; options missing for a step that takes them, or given for one that does not | `<tasks.tsv> has no row of <t>`; `the step <s> is not a step of the task <t>, of kind <k>`; `the options of the task <t> are ...` |
| 2 | `verification` with no `attempt` event | `the task <t> has no attempt` |
| 3 | a Shape input, the ledger or an options table missing or refused; a result that disagrees with its table; the specification of a target with no record in `docs/facts/` | `the records have no <path>`; `<path>: <reader error>`; `decision <d>: ...`; `the target has no record in docs/facts/` |
| 5 | two records of one base name | `two records are named <name>: <a> and <b>` |

**The session base** (that the context was built for the base and the task that the session starts with) is checked at the start of a session by `TaskSession`, before the sweep and before anything is made (**decided here**, task `T-jy6f`, #415, row 66r; [`session.md`](session.md#the-start-of-a-session)): a context of another task, or built at another base, is refused, so the `head` of a verifier's prompt is the base of its start row, which `nothing-material` reads.

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

**The question in the loop.** **Decided here** (R10, task `T-r7b9`, #426, row 66z; plan review 6262463499, conditions 2 to 5 and
Q2 to Q10; the Operator's goal count 4, comment 6262588622). The post, the route, the answer and the next attempt are a position of
the driver of a Shape task ([The driver of a Shape task](#the-driver-of-a-shape-task)), in the records and the forge alone and never
`bets.tsv`, so a question is reachable during Shape before any bet. The driver calls `RepairTask` first, as before.

- **The asking sessions.** The *asking sessions* of the latest attempt A are the sessions of A that ended `done` with the status
  `needs_context` and started before the first such result of A. Any session of A that started after that result is the *answer
  session*: the driver starts nothing else in A, so an owner session that itself ends `needs_context` is no asking session, and it is
  the stop of P5 below. The *questions* are one list over the asking sessions in the order of their start events, and over the rows of
  each result in the order of its file; the *question ID* is `<session>.<n>` of the row, so two sessions whose rows are both `n` 1 are
  two questions. A status other than `completed` and `needs_context` (`blocked`, `failed`, `decision_needed`) keeps the stop of the
  driver (`the session <id> of <step> ended <status>`).
- **Which stages can ask, and how each resumes.** Any task session of A can ask: the plan, the plan review, the developer and the
  verifier. The answer is not accepted by A: after the stages below the driver starts the attempt A+1 at the base of A, and the
  stage that asked is then the first stage the records show not done. For the plan and the plan review the sessions that ended
  `needs_context` in an earlier attempt are not the last session of their stage (the stage runs again, with the answer in its prompt
  file); for the developer and the verifier the stage is tied to its attempt, so the developer runs again at A+1 (a verifier's
  question discards the push of A, which A+1 makes again, **decided here**: the verifier of A+1 starts only after the new bound head,
  as before). A plan question is therefore followed by a completed plan and its review, never by the developer at once.
- **The stages**, each derived from the records and repeated safely; after each act the driver reads the records again:

  | # | The records show | The act |
  | - | ---------------- | ------- |
  | P1 | an asking session with no event `handoff` | `RecordHandoff{asking role → owner role}`; one with an event `handoff` that is not `valid` is a **stop** `<task>: the handoff asking role → owner role is <detail>`, on every restart |
  | P2 | a question row whose composed comment (`QuestionComment`) has no copy on the task's issue | post it and copy it (the read before the write of [the comments of a session](#req-005--the-steps-of-the-loop-after-the-scaffold): the copy, then the App's comment of the same bytes, then the lease, then the post and its copy) |
  | P3 | no answer session, and the questions are not all of one owner role (`answerOwner`: a human route, a label with no role, more than one owner role) | `Advance{Wait: "waiting for the Operator: <task>: question <id>"}`, the first question; no session starts |
  | P4 | no answer session, or the last one has no event `result` (a stopped run) | the milestone cap, `ResolveStep("answer-to-a-question", "", owner)`, `PairFor`, the prompt file, `TaskSession` under attempt A, read-only |
  | P5 | the answer session ended other than `done` with `completed`, or its answer rows are not one for each question | **stop** `<task>: the session <id> of answer-to-a-question ended <status>`, or `<task>: the answer <id> answers no question`, or `<task>: the question <id> has no answer` |
  | P6 | an answer row whose composed comment (`AnswerComment`, the question ID of the row's column `n`) has no copy | post it and copy it |
  | P7 | the answer session with no event `handoff` | `RecordHandoff{owner role → asking role}`; one that is not `valid` is a **stop** |
  | P8 | all of the above done | `StartAttempt` at the base of A |

- **The base of the answer session** is the base of the start row of the **last asking session** of A (a verifier's is its bound head),
  both for the prompt file and for the session, so the owner reads the tree the asker read. The asking results are inputs of the
  prompt file (the records table above), and the placeholder `{questions}` is the question IDs.
- **The answers.** Each answer row of the answer session selects its question by its column `n` (1-based into the ordered list of the
  questions, whatever the file order of the rows): the answer row of column `n` answers the `n`-th question of the list, and
  `AnswerComment` gets that question's ID. Every question of the list has one answer row and no answer row has an `n` outside the list:
  a missing answer is the stop `the question <id> has no answer` (the first unanswered in the list), and an `n` outside the list is
  `the answer <id> answers no question`. The correspondence is for the comment's first line only; it creates no acceptance binding,
  which stays the attempt-level `AnswerAccepted` (above, and `M3c`'s columns).
- **A wait and a stop repeat on the next call**, which reads the same records, and nothing retries by itself. The resolution of a
  question that goes to a human or to more than one owner role, and the Operator's answer, are later text's and are not decided here;
  the clause of `M3c` about acceptance columns is not a resolver of them.

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
the parse is row 66l's, the consuming wiring the step loop's (row 66k and its size children) — a
declaration of the boundary, nothing of it claimed built. An item is true only
for the comment kind it reads (a `plan-fields` listed on a review-record row is
false), the four comment items are computed from the kept copy and overwrite the
caller's value whether it is `true` or `false`, and a handoff with no kept copy,
or that code writes, has them false.

**The question and answer comments** are composed (composition only, no post) by
`internal/run` (row 66o, `T-gyvo`, #405) from the result's `question` and `answer` rows — first line and body as [a question during
the work](#the-checks-of-a-handoff) pins — and posted on the task's issue by the
post→copy machinery of the step loop (**row 66z**, `T-r7b9`, #426, a size child of row
66k, #367, split by its plan review); they are plain comments of the App, not artifacts of a handoff.
The parsed review record is also the input of the caller boundary that keeps
the loop uphill while a material finding is open (`architecture.md` §11):
row 66l's parse exposes the review's verdict, head and located findings and a
predicate that fails uphill — a record that does not parse, a verdict other than
`nothing material in scope`, or a bullet in the block is an open material
finding, so `material` or `not mergeable, findings recorded` with an empty block
never reads as none open — the step loop's wiring of it is row 66z's (#426, a size child of row 66k),
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

**Decided here** (R10, task `T-q0kr`, #362, row 66f; `internal/run`):

- **The description is cut by one function.** A CR or LF becomes a space, and a text of more
  than 140 characters (counted in characters, as `SetStatus` counts) is cut to 139 and ends in
  `…`, so `SetStatus` never refuses a description. The forms are under each context below.
- **Unchanged** means the state **and** the description. `layup run` reads the statuses of the
  SHA (`Statuses`, newest first), takes the newest of the context whose creator is the App's
  bot (the numeric ID that Start read; a bot ID that is not known is an error, never a match
  of any creator), and posts only when there is none or it differs in state or description.
  A status of another creator is ignored. The lease is read just before the post, and a lost
  lease posts nothing.
- **What is recorded.** The tables of `layup/gates` and `layup/spec` of a head are records
  (`gates.tsv`, `spec-check.tsv`), written once, before the post; after a stop the status is
  rebuilt from the table with no new run. `layup/rules` and `layup/verify` are computed again
  from git and the records. That a status was posted is read from the forge and never
  recorded, since a record of our own post would trust it.
- **Failures, by stage.** A failure of the check's input (`*gate.InputError`, a `target:` source
  of the full check, a rule check that gives `error`) is posted as `error`, with nothing
  recorded. A failure of the run's own plumbing (the records, a record that its reader refuses,
  the lease, the forge) is a Go error. The guarantee is that of the stage reached, with no
  rollback: a table committed stays committed when its post fails, and the next call reads
  that table, runs nothing and posts; a call that posts several contexts stops at the first
  failure and leaves the earlier posts as they are, and the next call posts the rest.

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
the tool could not make because of its input (a `*gate.InputError`: a revision, the
manifest at the base, git or `sh`) is `error`, with no row written.

**Decided here** (R10, task `T-q0kr`, #362):

- **The description.** `pass` when every row is `pass`; `clear: <reason>` of the first `clear`
  row when every row is `pass` or `clear`; otherwise `failure` with
  `<kind>: <result>: <reason>` of the first row, in the order of the manifest, that is neither
  of the two (this covers `not-active` and the `fail` of a `pending` kind with a product
  path); `error: <first line>` for an input error. A table with no row is `error`, never a pass.
- **A scratch tree that cannot be made** is not `error`: `layup gate` gives the rows that need
  the tree the result `not-active` and keeps the diff result of a `pending` kind
  ([`gate.md`](gate.md)), so the rows are written and the status is `failure`. A scratch tree
  that cannot be removed (`*gate.CleanupError`) comes with the whole table: the rows are
  written and the status is posted, and the error is then returned to the caller.
- **The rows of a head are the table's rows.** They are written once, in one fenced records
  commit, and a head that has rows is never run again. The reader refuses (a Go error, and
  nothing is overwritten): a head whose rows are not contiguous, two bases in the rows of a
  head, a `reason` that is `—` for a row that is not `pass` or that is not `—` for a `pass`,
  and a table that is not complete, namely rows whose kinds and states are not exactly those
  of `docs/gates.tsv` at the table's base, in its order. A head's rows are given to
  `GateFindings` by `GatesTable`.
- **The output and the progress of the run are kept.** `layup run` passes the step and the
  output of `gate.Run` through: the kind that runs and that the run lives are shown, and
  the held output of a command is written after its step, as `layup gate` writes it. The
  output has no record: the table is the durable record.

### layup/spec

`layup/spec` is the full `layup spec check` at the head
([The full spec check](#req-012--the-full-spec-check)): `success` when no row of its
table is `fail` (a row `clear` counts as a pass, whether the whole table is `clear`
because the task touches neither the PRD nor `docs/spec/` before the first bet, or
only the rules `classes` and `inventory` are, before it, on the specification
task's head); `failure` with the first failing rule in the description otherwise;
`error` when the check could not run.

**Decided here** (R10, task `T-q0kr`, #362; the sources of the input are decided in
[The full spec check](#req-012--the-full-spec-check) by task `T-cav3`, #412):

- **The description.** `pass` when no row is `fail` or `clear`; `clear: <detail>` of the first
  `clear` row (`no specification yet`, `before the first bet`) when none is `fail`;
  otherwise `failure` with `<rule>: <detail>` of the first `fail` row, in the order of
  the rules (`ran` last); `error: <what>` when the input of the check could not be read.
- **The table of a head** is the ten rows of `CheckFull`, written once with
  `spec.WriteSpecCheck`, by the schema `spec.TaskSpecCheckSchema`, in one fenced records
  commit before the post. The reader refuses what that of `gates.tsv` refuses (rows that are
  not contiguous, and a table that does not hold the ten rules in order, a `clear` head's
  too), and a table that is not complete never gives `success`.
- **A source that cannot be read.** The input builder gives a `target:` source error when the
  head's own tree or files cannot be used (a git read of the head): it is posted as `error`
  and writes no row, and a later call reads again. A `records:` source error (a record that
  the run keeps, or its reader) is a Go error and posts nothing.
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
  records name one path, a kind is neither of the two, or a PRD has no pre-image.

**Decided here** (R10, task `T-q0kr`, #362; the status):

- **The status** is the state of the result, and its description is the reason, or
  `<reason>: <path>` when the result has a path.
- **The register commit** is never the head's and never the run's changing records base. For
  a head that a session pushed, it is the `records` column of the start row of the session
  that has the `bound` event of that head (a session of the task; the sessions of one head
  must agree on the commit, else it is an error), as the check before the push reads it. For a
  head that code produced (the rendered commit has no bound session), the caller names the commit,
  as a full SHA, and a head that has a bound session refuses the caller's. A head with neither is an error.
- **A check that could not run** (`error`) is posted as `error` and is never a pass.

### layup/verify

`layup/verify` follows the verifier (`architecture.md` §8 step 6). **Decided here:**

- **Decided here** (R10, task `T-q0kr`, #362): the caller names the verifier's session; the
  order of the checks and the descriptions are these.
  1. No admitted harness outside the authors (`ReviewerPair`): `failure`,
     `not-active: no harness outside the authors`, at the head. Any other error of the
     pair is a Go error.
  2. The caller says that the verifier waits for its turn: `pending`, `waits for its
     turn`, at the head. In `M2e` no producer says so.
  3. The session must be a session of the task, of the role of the step `verification`, and
     the head it reviewed is the `base` of its start row; else it is a Go error, a caller's
     fault. A session whose recorded harness is an author of the change, or whose
     harness and model are not admitted (the probe passed, the model is used), is not
     independent: `failure`, `not-active: the verifier is not independent`, at the
     reviewed head, though another admitted harness exists.
  4. No event `result` or `refused` yet: `pending`, `verifier running`, at the reviewed head.
  5. A `refused` event, a result that is not `done`, a result with no kept copy whose
     SHA-256 is the one its `comment` row names, or a copy that is no review record:
     `failure`, `no review record`.
  6. `success`, `nothing material in scope`, only for a copy that passes `nothing-material`
     at the reviewed head; a record of another head is `failure`, `review record of another
     head`; otherwise `failure`, `material: <n> findings` (n the number of its bullets).
  A status of a verifier is posted on the head it reviewed.
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
  task is closed and redone (the closes: [`forge.md`](forge.md#the-calls-of-m2e), row 66w, task `T-xyws`) as a new task with its own issue, which takes the old
  task's row and predecessor links (`architecture.md` §8). The redone task has an
  `closed` event on the old one and a `rebased` event on the new one
  ([the records](records.md#nfr-001--the-records-of-the-task-loop)). The files a task touches are those of `git diff --name-only` from
  the merge base of its head and the default branch to its head (the architecture task's attempt starts at the specification's verified head, REQ-012 — Shape).
  The carry and the redo are [A task whose base moved](#a-task-whose-base-moved-the-carry-and-the-redo) (row 66x, task `T-89xf`).

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

**Decided here** (R10, task `T-xvvd`, #363; plan review conditions 5 and 6; `Ready`, `MarkReady` and
`MergeTask` of `internal/run`):

- **`Ready` is a pure read** of a full SHA: `Statuses` of the head and `EffectiveRules` of the default
  branch, read anew at each call. A read that fails or a rule that is malformed is a Go error, never a pass.
- **The four contexts**, in the order `layup/gates`, `layup/spec`, `layup/rules`, `layup/verify`, each
  need the **newest status of the context at the head, of any creator**, to be the App's and `success`.
  This is a second selection, apart from the one of posting (`PostStatus` compares with the newest
  status of the *App*, so that it never posts twice for one state): a newer status of another source
  blocks and can never cause a merge, at the cost that another writer can delay one. The reasons are
  `<context>: no status`, `<context>: a status of another source is newer` and `<context>: <state>`.
- **The native jobs** are the entries of the rules `required_status_checks` whose `integration_id` is
  that of GitHub Actions, `setup.ActionsApp` (15368, one home). A rule whose parameter is not a list of
  objects with a non-empty string `context` and a whole `integration_id` of 1 or more is malformed: a Go
  error, whatever the other entries. For each job, the check runs of that name by that App (the newest
  of each name, `filter=latest`; no second reading of "latest") must be at least one and each
  `completed` with `success`: `job <name>: no check run`, `job <name>: <status> <conclusion>`. A ruleset
  that requires **no** job of GitHub Actions is not ready (`no required job of GitHub Actions`: it
  fails closed, `NFR-004`; the setup always requires one of each kind).
- **`MarkReady`** reads the pull request by the register's number; a head other than the pull
  request's is `the head moved to <sha>`, a pull request that is not open is not ready. It calls the
  forge only for a ready draft, under the lease.
- **`MergeTask`**, the one function of the merge, first reads the pull request (a merged one only gets
  its event), marks ready as above, then reads the statuses and the jobs **again** just before the
  merge, so a status that changed after the marking blocks it, and merges at the head SHA it read, under
  the lease.
- **A 409 performs its read.** The pull request is read again and, for the new head, `Ready` (its
  statuses); the call returns `the head moved to <sha>` (with the reason that the new head is not ready,
  when so), sends **no second merge** from the old read and writes no event. A 409 whose pull request
  has the head that was sent, and any other error, are Go errors; the next call reads the pull request first.

### The merge at the head SHA and the close-out commit

**Decided here** (F-37; decision 6 of #333):

- **The close-out commit.** After `nothing material in scope`, a close-out session
  writes the target's close-out (the verdict, the line of the completed log, the
  resource record from the ledger) as one last commit. Code checks that this
  commit changes only the task file `docs/tasks/<task>.md` and
  `docs/tasks/completed.md` of the target, the two literal paths (**decided here**, task `T-xvvd`,
  #363: no setup record holds them, and the paths do not widen: a backlog edit that the
  target needs is in the head that the verifier reviewed, before the close-out; admitting
  `docs/tasks/backlog.md` to the last commit is the Operator's decision), and
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

**Decided here** (R10, task `T-xvvd`, #363; plan review conditions 3, 4 and 9; `PushCloseOut`,
`CheckRendered` and `CarryVerify` of `internal/run`):

- **The push hook of the close-out session** runs the checks of any task session, then two more,
  refused as a developer's would be, with the diff as a payload: **`scope`** when the head is not one
  commit on the session's base, changes a path other than the two, or holds one that is not a regular
  file (mode other than `100644`, a link, a deletion); **`verdict`** when the base does not carry the
  verdict. A Go error of a git read or of `VerifyAt` writes nothing.
- **A base carries the verdict** when it is the *verified head* (the `base` of the named verifier's
  start row, for which `VerifyAt` gives `success`) or the rendered commit above it that `CheckRendered`
  accepts. A verifier of another task is a Go error. A verified head that is itself a merge commit is
  no reason to refuse: what counts is how it got its verdict. **Decided here** (R10, task `T-89xf`, row 66x): a head that came from
  the verified head by a clean merge of the base that `BringUpToDate` made is carried, by the rule of
  [A task whose base moved](#a-task-whose-base-moved-the-carry-and-the-redo) (`mergeEdge`, which `carries` calls for the parent of a
  close-out and `CarryVerify` for the merge); any other merge commit is refused as `verdict`.
- **`CheckRendered`** accepts the rendered commit when its only parent is the verified head, it changes
  exactly the paths of the rendered records as regular files, and `RenderingEqual` holds with the
  **pre-image of a PRD pinned to the verified head** (`RenderingEqual` alone accepts any ancestor). It
  needs no session: code made the commit, which has no `bound` event. Row 66j calls it before its push.
- **`CarryVerify`** posts `layup/verify` `success` through `PostStatus` with the description `verified
  head <sha> plus the close-out` (or `plus the rendered commit`), after checking the whole edge again:
  the head is the close-out session's `bound` SHA, never the caller's (a session of the task and the
  role), or the accepted rendered commit; the commit count, paths, modes and the parent's carry. The
  other three contexts run on the new head as on any head.
- **The `merged` event** is written by `MergeTask`, or by `RecordMerged` (row 66y's `FinishTask` calls it on a restart),
  as one fenced records commit: the merge commit's SHA, the pull request's number, session `—`. Its
  attempt is that of the `bound` events of the merged head; none, or several attempts, is a Go error
  and no event (never guessed). A second call writes nothing.

**Decided here** (R10, task `T-ehb4`, #425, row 66y; plan review conditions 1 to 7 of 6262459073; `(*Sessions).FinishTask` of `internal/run`).
`FinishTask(ctx, FinishInput{Task, Records, Harnesses, Models, Bet, Copy})` takes a verified task to its recorded merge, in the records and the forge alone: it
reads before each write, holds nothing and returns at the first merge, wait, stop or handback. It gives `Advance{Task, Merged, Head, Wait, Stop, Attempt}`.
`Merged` is true with `Head` the merged head; `Attempt` is the attempt that the repair started (D7). Its inputs are the landed ones: `loopPosition` for the verified head
(row 66r), `PostStatus`, `CarryVerify`, `PushCloseOut`, `MergeTask` and `RecordMerged` (rows 66f, 66g), `TierOf` and `ResolveStep` (row 66e), and the scratch base of a session
(plan amendment 1 of #415). **D1, the positions**, each read before its write: (1) `A` is the largest event `attempt`, and the pull request is read by the branch
`task/<task>/<A>` (D8); (2) a pull request closed and not merged is the wait of D2, with no write of any kind and before any repair; (3) a merged one with its event `merged`
is `Merged` and writes nothing, and one without it gets its event and no repair that could start an attempt (D7; the lost event, O16); (4) `RepairTask` (D7); (5) the verified head
`V` and its verifier `X` come from `loopPosition` over the events up to the last handoff of a session of the verification role (the close-out's own sessions, their
bound and their refusals come after it), and the position must be the verified one, else the Go error `<task>: attempt <A> has no verified head`; (6) the carried head `C` of D3; (7) the up-to-date
check of D9 on `C`; (8) for the rendered commit, `CarryVerify(in, "")` and the three checks at `R`; (9) the close-out session (D4, D5); (10)
`CarryVerify(in, id)` and the three checks at `H'`, the close-out's bound head; (11) `MergeTask(task, H')` (D6).
- **D2, a pull request closed and not merged** (Q11 of #367, Q7 of the plan) is `Advance{Wait: "waiting for the Operator: <task>: the pull request #<n> of task/<task>/<A> is closed
  and not merged"}`. It writes nothing: no repair, no records commit, no forge call but reads; no redo is inferred from it (the redo is row 66x's). The wait repeats at each restart,
  and a pull request that a fresh read finds closed during the call (after the repair, or when `MergeTask` says it is not open) gives the same wait.
- **D3, the carried head and the three checks.** `C`, the head that carries the verdict, is the last event `bound` with session `—` of `A` (a bound of code: the rendered commit `R` pushed by the step bet, or the merge of the default branch that `BringUpToDate` of [row 66x](#a-task-whose-base-moved-the-carry-and-the-redo) makes above `V` or above `R`), and `V` when `A` has none; it is checked to carry the verdict (`carries`, which walks the merges) before anything is asked of it. **`R` is found by its own record**: the bound of code whose one parent is `V` (a merge has two parents and is never taken for it); the specification task must have it (a Go error otherwise: the step bet finishes it), and `CarryVerify` checks it again before it posts (`renderedOK`, the check of `CheckRendered`) with the typed inputs of D10. `V` and `R` keep their pinned inputs when `C` is a merge.
  `PostChecks` posts `layup/gates`, `layup/spec` and `layup/rules` of a head, in that order, as `PostStatuses` does (an extraction: the same order and the same delayed error of a scratch tree
  that cannot be removed), and `PostStatuses` calls it before `layup/verify`. `FinishTask` posts at `R` and at `H'` with `PostChecks` and never with `VerifyAt`: the authors of `base..H'` include the close-out session's
  harness, so `VerifyAt` could post `not-active` at `H'` after the `success` of `CarryVerify`; `layup/verify` of a carried head comes from `CarryVerify` alone. **The inputs of a head are pinned
  apart:** at `R` the base is the attempt's, `PreRev` is `V`, `Copy` and `Records` are D10's, and `RegisterRev` is the `records` column of the start row of the session bound to `V` (`registerRevOf` at `V`);
  at `H'` the same `Base`, `PreRev`, `Copy` and `Records`, and no `RegisterRev`, since the close-out session's start row is the snapshot and `registerRevOf` refuses a caller's commit for a head with a bound session.
- **D4, the close-out session.** Its step is `close-out` (role `software-developer`, tier `by-class`). The plan is the kept copy of the last valid comment handoff of a session of the plan role
  (`planHandoff` and `copyPathOf`, the selection that `PromptFile` makes, extracted), parsed by `ParsePlanComment`; `Position` is read at the call, and its error and `ResolveStep`'s are the call's errors; the tier is
  `TierOf(plan, position)` and the pair is `PairFor(role, tier, …, nil)`, with no exclusion (the close-out needs no independence). The prompt file is built at `C`, and `C` is not on the
  default branch: the session clones the scratch branch `layup/base/<task>/<A>` that `git.UpdateRef` set to `C` (plan amendment 1 of #415), through the one `startAs` of the driver, whose `Push` hook is `PushCloseOut`.
- **D5, how a close-out session ends.** The sessions of the position are those of the close-out role in `A` whose event `session` comes after the last handoff of `X` and whose start-row base is `C`; the last of them
  by the order of the events decides (never the clock), so a session that was rerun is not the one read. None, or one with no `result` and no `refused`, starts a session (a rerun; the milestone cap is asked first). A
  `refused` event of a session with no start row after that handoff is the stop `<task>: the start of the session <id> was refused (<reason>)`; a `refused` of one with a start row is `<task>: the session <id> was refused (<reason>)`;
  a result other than `done`, or a status other than `completed`, is `<task>: the session <id> of close-out ended <status>`. None of them is tried again. **A session that ended `done`, `completed` and has no `push`, `bound`
  or `refused`** (the end commits before the push hook runs) resumes its checked push from the kept head `refs/layup/sessions/<id>` of the run's clone, with `PushCloseOut`, and starts no session: **a restart makes the clone anew, so `cloneStep` fetches the refs `refs/layup/sessions/` of the old clone, set aside in `DIR/recovery/OWNER/NAME/` (no target's clone; [`run.md`](run.md#the-restart) step 2 holds the same contract, with the mark `DIR/recovery.imported/OWNER/NAME` that a restart writes and syncs after the import and before the removal, so that an interrupted removal is finished and never read again), into the new one** (`git.Refs`, `git.Fetch`; an old clone that cannot be read is a failure of the step, not an absence), and a kept head that is in no clone is the Go error `<task>: the close-out session <id> ended with no push, and the head it kept is not in the run's clone` and never a second session (condition 3 of the plan review binds; a read of the refs that fails is the error of the read). A `push` with no `bound` is `RepairTask`'s (D7) and, if it is still there, a stop.
  A `bound` is `H'`.
- **D6, not ready.** `MergeTask` is called with `H'`. A merge gives `Advance{Merged, Head: H'}`. A head that moved (`the head moved to <sha>`) is the stop `<task>: not ready at <H'>: <reason>`, and a pull request that is no longer open, also when the refusal of the merge (409) is followed by a read that finds it closed at the same head (`afterRefusal` gives "the pull request is not open" before it maps the head), is D2's
  wait or, merged, the recorded merge; every other reason (a context not `success`, a job pending or failed, no required job) is the wait `waiting: <task>: not ready at <H'>: <reason>`, which a restart reads again. **Rejected:** a wait only for
  `pending` and a stop for a failure (a failed job may be run again, and only the Operator knows).
- **D7, the repair.** `RepairTask` comes after the read-only checks of positions 1 to 3, so that D2 writes nothing even when a session comment is missing or a push is unresolved. A repair that starts an attempt (a lost
  push) gives `Advance{Task, Attempt: n}` and nothing else, **but never for a merged pull request**: that one is terminal before the repair (`terminalPull`), and only the resolution of the pushes of the task (`resolvePushes`: the `bound` or the `refused push-lost` of each, the first part of `RepairPushes`, with no attempt) runs, to find the attempt of the merged head for `RecordMerged`: the attempt that was read is no longer the latest, the task has no verified head, and row 66k calls the driver of a Shape task for it; no error string stands for it.
- **D8, the pull request of the attempt** is read by its branch and its number must be the register's `pr`: a missing pull request, a register without one and a different number are Go errors (`<task>: the pull request <n> of <branch> is not the register's <m>`).
- **D9, the up-to-date check** has one home, `(*Sessions).UpToDate(ctx, head)` of `baseup.go` (row 66x, #424), which `BringUpToDate` and `FinishTask` both ask: it fetches the default branch and says whether its head is an ancestor of `head`. `FinishTask` asks it of `C` of D3, the tip that row 66x's section names ("Row 66y takes the head to close out from the tip and asks `UpToDate`"), never of `V` alone and never of "the last bound of code" taken for `R`, so a task that `BringUpToDate` brought current is closed out from its merge. A failed fetch or read is a Go error, never evidence of a moved base; a head that is not current is the Go error `<task>: the base moved (row 66x's)`. This task's own `upToDate` was deleted when row 66x landed (condition 4 of the plan review: one home in both landing orders).
- **D10, the typed inputs of the heads after the bet.** `FinishInput.Copy(ctx)` gives the first copy of the answer that the latest version of the inventory names (`(*state).betCopy`; nil before the first bet) for **every** task, since the full check of every head after the confirmed bet reads it (`classes` and `ran` fail with no copy); the driver takes the same bytes as `ShapeTaskInput.Copy` for its own statuses, and row 66k passes both. `FinishInput.Bet(ctx, verified)` gives `RenderedBet{Copy, Records}` for the specification task (nil for any other kind): the answer's copy and the rendered records, rebuilt from `bets.tsv`, the latest
  version of the inventory and the verified head, as the step bet built them (`(*state).renderedBet`, the extraction of `renderedStage`'s inputs). Row 66j planned it as `BetInputs` and landed the step without it; row 66k passes it.

### A task whose base moved: the carry and the redo

**Decided here** (R10, task `T-89xf`, #424, row 66x; plan review conditions 1 to 8 of 6262560771, the Operator's decisions of
6262587821, and conditions 1, 5 and 8 of the plan review of #367). `(*Sessions).BringUpToDate(ctx, CarryInput) (Advance, error)`
leaves a verified task in exactly one recoverable end state: its head current with the default branch, carrying `layup/verify`, or
the task redone. The loop-permission check is its caller's (row 66k). It gives `Advance{Task, Verified, Head}` for the first,
`Advance{Task, Replaced}` for the second, and a `Wait` or a `Stop` for the cases below; a git read that fails is an error that writes
nothing, never a redo.

1. **The tip and the default head.** The tip `X` of the open attempt `A` (the largest `attempt` event) is the SHA of its last event
   `bound`, of any session: the `bound` of code has none (`TipOf`; `lastBound` still skips them). The default head `D` is fetched by
   `SessionStore.FetchBranch` (`git.Fetch` under the store's lock and token; `refs/heads/<default>` of the run's clone). The base of
   `A` must be an ancestor of `D`; else it is an error (fail closed).
2. **The up-to-date predicate has one home** (condition 4 of the plan review of #425, row 66y):
   `(*Sessions).UpToDate(ctx, head) (current bool, def string, err error)` fetches `D` and gives whether `D` is an ancestor of `head`.
   A head brought current by a merge of this section is current; testing ancestry against the verified head alone is wrong, and so is
   selecting only the verified head or the rendered commit. Row 66y takes the head to close out from `TipOf` and asks `UpToDate`.
3. **Up to date** (`D` is an ancestor of `X`): no write and no request. When `X` is a merge commit that this section made, the carry
   is posted again (`CarryVerify` is idempotent), since a stop between the `bound` and the status leaves it undone.
4. **Not up to date.** With `mb = MergeBase(X, D)`, `T = DiffNames(mb, X)` (the files the task touches) and
   `P = DiffNames(mb, D)`: `T ∩ P = ∅` leads to `git.Merge(X, D)`; a clean merge is the carry, a conflict or an overlap is the redo.
   Any other failure of git (including the read that classifies a conflict) is an error. A closed pull request that is not merged,
   with no redo decided, is `Wait` `waiting for the Operator: <task>: the pull request <n> is closed and not merged`; a merged one is
   a `Stop` (`RecordMerged`, row 66y, records it). The verdict must be `success` (`verified`) before either.
5. **The merge commit.** `git.Merge(dir, x, d, message, who)` runs `git merge --no-ff --no-edit` in a scratch work tree at `x`,
   with the author and the committer of the App and, as the time, the later committer time of `x` and `d`, so one input gives one SHA.
   Exit 1 with unmerged entries is `clean=false` (the work tree is removed); every other failure is a `FailedError`. The commit is
   kept by the ref `refs/heads/layup/merge/<task>-<attempt>`. Before it is used, the commit is checked by `mergeEdge`.
6. **The push of the merge is announced.** The events `push <M> <branch>` and `bound <M> <branch>` with session `—` (the form of the
   rendered commit's push, row 66j; `RepairPushes` skips them: their step resolves them). The order: the event `push` (one fenced records
   commit), the lease, `PushHead`, the event `bound`, then `CarryVerify`. The branch must hold `X` (to push) or `M` (pushed, to bind), else
   it is `errPushRefused`. A restart makes `M` again (the same `X` and `D` give the same SHA) and finds the events and the branch.
   **An outstanding push is resolved first** (`adoptPushed`, round 1 of #424): before another merge is computed, an event `push` of code with
   no `bound` is read against the branch of the forge. A branch that holds the announced `M` (the push was made and its `bound` was lost) is
   adopted whatever `D` has become: `M` must be a merge whose first parent is the tip and that the carry rule accepts, its `bound` is written and
   its status posted, and the comparison then goes on from `M`, so a default branch that moved again gets a merge of `D2` over `M`. A branch that
   still holds the tip leaves the announced `M` unpushed: the next merge `M'` gets its own `push`, and a `push` with no `bound` is harmless to every
   reader, which looks only at `bound`. Any other head of the branch (a push from outside) is `errPushRefused`, before any comparison, merge or redo,
   and writes nothing.
7. **The carry rule (`mergeEdge`).** A merge commit `M` with the parents `[p1, p2]` is carried when: `p1` carries (the verified head, the
   accepted rendered commit, an earlier accepted merge, or the one-commit close-out head of such a head that `scopeOK` accepts, the
   close-out only under a merge), `p2` is an ancestor of the default head, and with `mb = MergeBase(p1, p2)`,
   `T = DiffNames(mb, p1)` and `P = DiffNames(mb, p2)`: `T ∩ P = ∅`, `DiffNames(p1, M) ⊆ P` and `DiffNames(p2, M) ⊆ T`. Disjoint
   *result* diffs are not enough: a merge that takes the base's version of a file both parents changed has the diffs `{a}` and `{}` and
   changes a task path. The rule proves the combined tree from the commit alone. `CarryInput` gains `Merge`; `CarryVerify` with no session
   and `Merge` posts `verified head <sha> plus a clean merge of the base`, after checking the edge again. `PushCloseOut` on `M` composes
   unchanged, since `carries` accepts `M`.
8. **The redo.** A task whose branch holds the rendered commit (a `bound` of code with one parent) is a `Stop`
   `waiting for the Operator: <task>: a redo would drop the rendered commit` in `M2e` (the Operator's decision on #424; the follow-on is
   row 66t's, `M2g`). Otherwise, in this order:
   1. **Preflight, no write.** The pull requests of the attempts' branches are read (a merged one is an error), the new ID is drawn by
      `taskID`, and every file of the combined commit is built and read by its own reader (`ReadRegister`, `ReadEvents`, `ReadOptions`).
   2. **The intent.** One fenced commit writes `tasks/<old>/redo.txt`: `<new ID> <D>` and a line feed, before any act of the forge. A
      restart reads it before any comparison, so a decision that flips afterwards cannot leave the old issue closed with no redo.
   3. Each pull request of the old task that is open is closed, the lease read right before each close; one that `ClosePullRequest`
      gives merged is an error `<old>: the pull request <n> merged after its redo`. Then `CloseIssue(old)`, then `OpenIssue`, each after
      the lease is read; a lost lease stops before the act.
   4. **One fenced records commit** holds `tasks.tsv` (the row replaced in place), the event `closed` of the old task (attempt `A`, detail
      the new ID), the event `rebased` of the new task (attempt 1, base `D`, detail the old ID) and the bindings. Last, so no act waits for a
      repair after it.
   Recovery: every act before the commit repeats from the intent (a repeated close is a `200`, `forge.md`). **Known limit:** a stop after
   `OpenIssue` and before the commit opens a second issue at the restart; the records hold one redo, and the first issue stays open on the
   forge, as in [The opening of a Shape task](#the-task-register).
9. **What the new row takes.** The row stays in place; the kind, the requirement IDs, the size and the predecessors are kept, the `pr` is
   cleared, the issue and `opened` are new, and each task that named the old one names the new one. The bytes of each
   `tasks/<old>/options/<decision>.tsv` are copied to the new ID. The issue text is the composer of the kind with the new ID and the
   inputs the records keep: the decision and option of the binding for a bet's architecture task; for a Shape task its **opening pin**: the SHA-256 of
   `spec/draft.tsv` and the records commit that its issue was opened with (`openShapeTask` keeps them in `tasks/<task>/pin.txt`, in the commit that
   registers the task: `<sha256> <records commit>`, a 64 and a 40 hexadecimal characters, and a line feed, captured once for the issue and the pin together), which the combined commit of a redo copies to the new task by value, so any
   number of redos keeps the pin of the original issue, and a later draft, a later snapshot or the first session's records change nothing. A Shape task with no pin, or with one that is not two full hashes,
   is an error before any write, never a pin guessed from another record. Attempts, sessions and the frozen list are not carried.
   The first attempt of the new task starts at the `base` of its event `rebased`, in the driver.
10. **The older attempt's pull request** (O22). `RepairTask` closes, for an open attempt `A > 1`, the open pull request of each older
    attempt's branch, only when the register's `pr` is the pull request of the branch of `A` (so a new attempt that has no pull request
    leaves the older one open). A task without such an attempt needs no forge.

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
([`scaffold.md`](scaffold.md#at-each-later-start)). The check of this section is built by row 66g (task `T-xvvd`, #363). **Decided here:** from the
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

**Decided here** (R10, task `T-xvvd`, #363; plan review conditions 7 and 9; `explainHead` and the
reading of the merges in `internal/run`):

- **The switch.** The rule below applies only once a merge of a task is known: an event `merged` of a
  task of the register, or a task with a pull request that the forge shows merged with a merge commit.
  With none known the rule of the setup is unchanged, byte for byte. A register or a task's `events.tsv`
  that is absent is no merge; one that is present and that its reader refuses, and a read of the
  forge that fails, are a `fail` of the row `tree` that names it (these reads are inside the check, so
  the table is still committed; the records of the setup, read before it, end the step with no table, as before).
- **The walk** goes from the head along the first parents: each commit must be a recorded merge or a
  merge that the forge shows for a task with no event yet. It passes at the bottom, the last commit of
  `setup/tree.tsv` or the code-owner commit above it (the setup's own judgment, unchanged), **only**
  after at least one task merge was accepted: the setup head alone after the first merge is a rollback
  and fails. Any other commit, a walk longer than the known merges plus one, or an unreadable parent
  fails; ancestry alone would admit a push to the branch.
- **The detail.** A pass keeps `—`, and the step prints `tree: <sha> merge of <task>, not yet recorded`
  for each unrecorded merge on the walk; a fail is `<head> is not a recorded merge`.

## REQ-012 — Shape

Requirement `REQ-012` (`F-0003#51`): Shape writes the specification, the target's
architecture and the gate-activation batch (`architecture.md` §8). **Decided
here** (decision 7 of #333): for each architecture decision an **architect
session** writes at least two options, each with a fixed table, the block
`options` ([the records](records.md#nfr-001--the-records-of-the-task-loop)): for
every constraint that applies a `pass` or `fail` with its evidence. **The
constraints that apply** are listed by code, never by the session: each requirement
of kind `non-functional` of the draft or the PRD, by its final ID `NFR-NNN` (the requirement IDs of a Shape task, above, are the one home of the correspondence) and its criterion (a
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
  survivors that the caller's filter gives ([The instruction of a task step](#the-instruction-of-a-task-step), task `T-c6sx`). **Its base, decided here** (R10, task `T-3f94`, #364, plan amendment 1, #432): the architecture task starts at the verified head of the specification task, never at the default branch. It is opened when the specification task has a verified head; its first attempt starts at that head, and each later attempt keeps that base. The specification is not its predecessor in `tasks.tsv`: a predecessor has merged, and nothing merges before the bet. **Reason:** the full spec check judges a head that touches `docs/spec/` against the PRD of that head, so the architecture is judged on the tree that the merge order gives it, and its `Changed`, its statuses and its verifier's range hold only its own changes. **Rejected:** a PRD read from another task's head (two trees in one input); `covers` or the PRD rules cleared for an architecture head (six rules would have to be cleared, and NFR-004 forbids it). **After the bet**, code opens one task of kind
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

**The options session** (**decided here**, R10, task `T-3f94`, #364, row 66h; plan review conditions 3 and 5). The decisions and their
options are not a task's output: the filter must run before the architecture task writes, so they are written by a **session of the
Shape phase**, role `systems-architect`, tier `reasoning`, as a session of a phase (no attempt, no event; the pair as for every session of
a phase, and the cap of milestone 0, which Intake and Shape share). The specification task, the target's architecture task and the
activation batch are tasks like any other.

- **Its input** is `spec/draft.tsv`, `spec/spans.tsv`, `start/problem-statement.md` and `intake/answers.tsv` of the records, and the list of
  the constraints that apply, which code writes into the prompt with their IDs and criteria (it is computed, not a record). The target's
  `docs/` is in the session's clone.
- **Its output** is one file, `options.tsv`, read by the one reader of the output file of a phase (a link, a file that is not text, a size
  over 1 MiB or a text that is not UTF-8 refuses it): the header `decision`, `option`, `constraint`, `result`, `evidence` and a row for each
  row of the tables of the block `options`, with the decision first. **Reason:** a session of a phase names its output files before it
  starts and the decisions are not known then, so a directory of files cannot be read; one bounded file can, and code splits it by decision.
- **Its checks.** The whole output is refused when it holds no decision, when a decision's name is not a word of lowercase letters, digits
  and `-`, when a decision's table fails the reader of the block `options`, or when a decision has fewer than two options (at least two
  options, above). A refused output is a failed step, and the next run starts a new session. The evidence of a pass ("a path or a link")
  is checked for its presence only: the block names no base to resolve a path against, a known limit.
- **What is kept** is `spec/options/<decision>.tsv` for each decision, in the one records commit of the session's end.

**The sources of the constraints** (**decided here**, R10, task `T-3f94`; the notes of round 1 of #349). At each run `layup run` builds the
list that `spec.Standing` takes, from the records and from the target at the head of its default branch, and never from the session:

- the requirements of kind `non-functional` of `spec/draft.tsv`, by their final ID `NFR-NNN` (the requirement IDs of a Shape task, above:
  the options table, the rows of the Shape tasks and the bet's `constraint` values agree), with their criteria;
- the facts of class `constraint` of `spec/spans.tsv`, by `F-<n>`. A fact that the confirmed inventory marks with the scope `out` does not
  apply: it is out of the delivery. M2e runs Shape once, before the first bet, so the PRD's requirements and the inventory's facts, which
  the Shape after a bet reads, have their reader in row 66j and none here;
- the `Inv-<n>` entries of the target's `docs/guardrails.md` (`verify.InvariantIDs`, the matcher of the check `guardrails`);
- the decision records `docs/adr/NNNN-*.md` (regular files directly in the directory) whose first non-empty line after `## Status` is
  `Accepted`, or begins `Accepted. Amended by ` (the rule of the baseline's linter). A record with no status line is an error, and the
  step fails; it is never "not accepted".

**The round, the rewrite and the stop** (**decided here**, R10, task `T-3f94`; plan review note 2). The round of a decision is read from
the records alone: when `spec/options/round-1/<decision>.tsv` exists, the table `spec/options/<decision>.tsv` is the rewrite; otherwise it
is round 1. When a table of round 1 has no survivor, **one** rewrite session runs for every such decision at once, with each dropped table
and its reasons as input; its output has the form above and holds exactly those decisions (a decision it adds or omits refuses it). Its one
records commit keeps the new `spec/options/<decision>.tsv` and the table it replaces as `spec/options/round-1/<decision>.tsv`. **One
rewrite is kept for each decision**; a rewrite session that stopped before its commit left nothing and may run again. The filter runs on
the rewrite, and a decision with no survivor then is the round with no progress: `layup run` reads the lease, then posts on the milestone's
issue, which is the Intake issue (`issue.intake` of `start.tsv`; Intake and Shape are milestone 0, and no other record names an issue),
one comment whose first line is `LAYUP Shape: no survivor` and that lists, for each decision, each dropped option of both rounds with its
reasons, in the filter's order. It posts it only when no comment of its App has the same body. The row of the step is `done` with the
detail `waiting for the Operator: shape: no survivor`, a wait: the last row, the lease released, the exit status 0. A restart judges the
tables again and gives the same row.

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

**The questions asked, the rows and the wait** (**decided here**, R10, task `T-3f94`, #364; plan review conditions 6 and 7).

- **Which.** A question is asked when its source is not `follow-up`, no Shape row that is still `asked` asks it, and either its state is
  `open`, or it is `asked` and no `follow-up` row that asks it again is `answered`. An `open` question was judged open by the review or by
  code, so it is always asked again; an `asked` one whose follow-up was answered has its answer there (the rule of `budget.tsv`).
- **The rows.** One row of `questions.tsv` for each: source `follow-up`, ref the question, block its block, kind `shape` (the block
  `questions` admits that kind for that source; the two steps of the follow-up of Intake ignore such rows), asker `layup run`, state
  `asked`. One records commit holds the rows. Then one comment on the milestone's issue, the Intake issue, is posted, in the form of the
  batch, with the first line `LAYUP Shape: the questions left open`, and copied into the records. A restart that finds the rows and no such
  comment posts it from the records; one that finds the comment and no copy copies it; neither posts a second.
- **The answers** are read from the copies of the comments of that issue made after the round's `asked`, by the rule of the step `answers`.
  An answer is recorded once, the first time it is read: the Shape row becomes `answered`, and so does the question it asks again, with the
  same time. An answer to a field of the intent form that does not parse (a decimal, `B` below `U`) leaves that question `open`, and the next
  run asks it again in a new row, so a later correction is the answer to that row. A value of `budget.tsv` takes the last Shape row of its
  field as its `source`, and replaces the value that is there only when the three values still parse with `B` below `U`.
- **The wait.** Only `band.B`, `band.U` and `appetite.hours` block: the step is `done` with the detail `waiting: the intent form: <IDs>`
  while `budget.tsv` is missing, or a Shape row asks one of those fields and no answer to it parses (`success-criteria` is a list, not a
  number, and never blocks). With `budget.tsv` missing and no question of those fields open, the step fails: `budget.tsv is not written and
  no question of the intent form is open`. The other questions never block; they stay for the brief (row 66j).

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

**The merge order after a bet:** the specification (or the rendered task), the tasks that a changed class reopens ([row 66ab](#the-specification-task-that-a-changed-class-reopens), in the order of `spec/reopenings.tsv`), the target's
architecture tasks, the rule batches, then the build tasks (`layup run` merges the
base into the branch of each next task, a Shape task included, as `architecture.md` §8
says for the build tasks, because the branch rule is "up to date"); a bet that changes a class or a priority does not
touch the approved batches. An architecture task whose base is a specification's head merges only after that specification has merged. A
specification that is closed and redone also closes and redoes the architecture task, which starts again at the new verified head,
because its branch holds the old head's commits.

### The brief of the first bet

Row 66u (task `T-h65e`, #418, a size child of row 66j). The entry point `betBrief(milestone, kind)` of `internal/run` keeps
`bets/<milestone>/brief.md` only when the brief passes every check of its bet; row 66j (the step `bet`) calls it, posts the brief and
parses the answer. Everything below is **decided here** (R10, note 6198919867 of #333, which the step table of
[`session.md`](session.md#the-step-table-and-the-kinds-of-a-steps-context) leaves open: it names no record for the step `bet-brief`).

- **The first bet only.** The kind is `first`; `later` is an error, because no row of M2e builds the brief of a later bet. The milestone is
  one path part, `^[A-Za-z0-9._-]+$` and not `.` or `..`. These two checks come first and read nothing. Then a kept `brief.md` is returned
  as done before, with no check and no write; this comes before every precondition of a fresh build. A fresh build refuses a `bets.tsv`
  that exists, because its reader is row 66j's.
- **The inputs**, read in this order by their readers before any session; a record that is missing or refused is an error that names the
  path and the line, and nothing is written. The records go to the session byte for byte as files of `input/`: `problem-statement.md`
  (`start/problem-statement.md`), `spans.tsv`, `draft.tsv`, `questions.tsv`, `answers.tsv` (`intake/answers.tsv`) and `budget.tsv`.
  `bets.tsv` is not an input. `brief.lines` is read from `parameters.tsv`; `—` is unset.
- **The verified specification.** The Shape tasks are the first register row of the kind `specification` and the first of the kind
  `architecture` (a later row of the kind is another task). Each must have a verified head: the position `posVerified` of the driver of a
  Shape task, read and never driven (`shapeView` and `loopPosition`), whose head is the bound head of the developer. Its row has a pull
  request, which the forge gives as found with that head; a pull request not found, or whose head moved, is an error before any session.
- **The pinned draft.** `spec/draft.tsv` of the records is the draft that the Shape tasks were opened on: Intake is done, and the row of
  the specification task lists `finalIDs(draft)` (the requirements column), which is checked equal. The content of the requirements is
  **the PRD at the verified head of the specification task** (`docs/prd/PRD-*.md`, read by `filesAt`): the statement and the acceptance
  criterion of each requirement ID of the register row, and each success metric. The PRD must hold exactly the requirement IDs of the
  register row, each with a statement and a criterion. The brief shows the PRD's words, so a statement or a criterion that differs from the
  draft's cannot reach the approver under an ID that agrees. **Reason:** the approver answers one line per requirement ID of the
  specification that row 66j parses the answer against, and a session may word the PRD better than the draft.
- **The five texts that code writes**, also files of `input/`; each line is one line, and a `](` or `][` in a text that code copies is
  written `] (` or `] [` so that no copied text is a link:
  - `requirements.md`: `### The requirements`, then for each requirement ID in the order of the register row `- <ID>: <statement> Criterion:
    <criterion>`; then `### The success criteria` and one line `- <metric>` for each success metric of the PRD, or `- none`.
  - `solution.md`: `### The verified heads` and one line `- <task> (<kind>): [pull request #<n>](https://github.com/<owner>/<name>/pull/<n>),
    head <sha>` for the specification task and then the architecture task; then, for each decision in name order, `### Decision <d>` and one
    line for each option in the order of the filter: `- <option> (a survivor)` or `- <option> (dropped: <reason>; <reason>)`, and, for a
    decision with one survivor, `- <option> (the one survivor, decided at Shape)`. A decision with no survivor is an error.
  - `open.md`: `### The questions still open` and one line `- <Q-NNN>: <question>` for each question still open, or `- none`. A question
    is still open when it is not a follow-up and its state is `open`, or `asked` with no answered follow-up; this is `leftOpen` of
    [the questions left open](#the-questions-left-open) without its last clause (no Shape row asks it yet), which would hide the questions
    that Shape asked and that stay open.
  - `facts.md`: one line `- <number or Q-NNN> (<class>): <text>` per fact: the span facts first, in the order of `spec/spans.tsv`, with
    the `text` column; then one fact per answer, by its ID, with the answer's text and its class `answerFactClass`.
  - `links.md`: the allowed link targets, one for each line, in order.
- **The class of the fact of an answer** (`answerFactClass`): the class of the span fact that holds the first byte of the question's quote;
  a follow-up, and Shape's row, takes the quote of the question it asks again (its `ref`); a question with no quote, or a quote outside
  every fact, is `context`. The approver confirms or changes it with a line `fact <id>: <class>`.
- **The links.** A link of the brief resolves exactly when its target is in this set, computed and never invented: for each Shape task,
  the URL `https://github.com/<owner>/<name>/pull/<n>`, and, for each path of `git diff --name-only <base> <head>` (the base of the
  attempt) that `git ls-tree <head>` lists as a regular blob of `[A-Za-z0-9._/-]`, the URL
  `https://github.com/<owner>/<name>/blob/<head>/<path>`. The forge of phase 2 is GitHub.

**The session.** One `PhaseSession` of the step `bet-brief`, role `product-owner`, on `PhasePair`. Its prompt is the first block below;
the finding of a refusal is the second. Its output is the directory `brief` with the files `problem.md`, `appetite.md` and
`rabbit-holes.md`, read by `ReadOutputFiles`. **Code assembles the brief**; the session writes no heading, no option and no solution:

```text
## The problem
<requirements.md><problem.md>
## The appetite
<appetite.md>
## The solution
<solution.md>
## The rabbit holes and the no-gos
<open.md>architecture change: no
<rabbit-holes.md>
## The facts
<facts.md>
```

Each file ends in a line feed (code adds the last). The first bet has `architecture change: no`. **Reason:** a session that writes the
whole brief and is checked part by part against code's bytes re-copies text, and the prose of its solution part goes unchecked; the
presentation of every choice is code's.

```text bet-brief
You write three parts of the brief of the bet of milestone {milestone}, for LAYUP. Change no file in this repository and make no commit.

Code writes the other parts of the brief and puts yours between them. Read these files:
- The problem statement: {problem}
- The fact spans: {spans}
- The draft: {draft}
- The questions: {questions}
- The answers: {answers}
- The budget: {budget}
- The requirements and the success criteria, which code writes: {requirements}
- The solution, which code writes: {solution}
- The questions still open, which code writes: {open}
- The facts, which code writes: {facts}
- The links you may use, which code writes: {links}

Write these three files of Markdown in the directory {output}, and no other file there:
- problem.md: the problem of this milestone in plain prose for the idea owner, from the problem statement and the requirements.
- appetite.md: the money and the time you ask for, with the line "appetite: USD <n>, <n> hours" exactly once, each <n> a decimal such as 12.5, and the reasons in a few lines.
- rabbit-holes.md: the rabbit holes and the no-gos of this milestone in plain prose, one line each.

The whole brief is at most {lines} lines, and your three files together have {room} lines at most. Code checks the brief after you, refuses it and asks you again for each of these:
- A line that starts with "## " in any file: code writes the headings.
- A line that starts with "architecture change:" or "approve batch " in any file: code writes the first, and this bet has no batch.
- A link in any form but [text](target), or a target that is not a line of the links file.
- An option of a decision with more than one survivor, named anywhere in your files. Code lists the options in the solution part. You recommend none of them, not by name and not by a phrase such as "the first one": the idea owner chooses.
- No appetite line, or two. A money cap above U of the budget, hours above appetite.hours of the budget, or a money cap that with the spend of the milestone so far is above U.
- A brief of more than 65536 characters.

{refused}Then end.
```

```text bet-brief-refused
The check refused the last brief: {finding}
Write the three files again.
```

**The order of the checks** of `briefKeep`, which refuses the output as a `refused file` on the first finding of a stage, with the findings
of that stage joined as `rule: detail; …`: (1) each of the three files has no line that starts with `## `, `architecture change:` or
`approve batch `; (2) `spec.CheckBrief` on the assembled brief with `brief.lines`, the links and no verified batch; (3) the option-word
rule: no word of the brief outside code's parts equals, in lower case and bounded by characters outside `[a-z0-9-]`, an option of a
decision that has more than one survivor; (4) `CheckCommentBody`, so that a kept brief can always be posted. A pass keeps one file,
`bets/<milestone>/attempts/<session>.md`, the **candidate**, in the one records commit of the session's end with its telemetry row.

- **The requirement and what is checked.** The brief recommends no option, directly or indirectly (Shape and The bet, above: two survivors give a brief with no recommendation). Code's
  presentation of the options, and the option-word rule, are partial evidence for it. Prose that recommends without naming an option
  ("the first one") is not found by the rule and is **noncompliant** all the same; the rule is not the requirement. The rule checks only the
  decisions with more than one survivor.
- **Before the first session**, with no session started: `brief.lines` unset is `ErrBriefLinesUnset`; the parts of code alone (the five
  headings and the lines of code) above `brief.lines` is `ErrBriefOverLimit`. Row 66j maps both to `done waiting for the Operator: …`.
- **The appetite is judged after the session's telemetry**, which `endPhase` commits after `Keep` has run: so `Keep` cannot read it. After
  `PhaseSession` returns, the run reads `telemetry.tsv` of the records, which holds the session's own row, takes the known money of its rows
  (`ledger.KnownMoney`: every `reported` or `computed` row in exact decimals, probes included; unknown rows add nothing; a known row not in
  USD is an error that wraps `ErrMilestoneCap`) and calls `spec.CheckAppetite(nil, caps, band.U, appetite.hours, known)` on the line of the
  candidate; at the first bet every row since Start is milestone 0. A finding is a refusal, like a refused file. The cap of milestone 0 is
  asked before each attempt (`PhaseSession`), and bounds the retry. A cap above `band.U`, and a cap below `band.B`, are checked through
  this caller; the second bet of 0.6 `U` is the check's own case with prior caps, which this first-bet caller does not pass (its caller is
  row 66t's, the later bet, in M2g).
- **The promotion.** A candidate that passes both is kept by one fenced records commit of `bets/<milestone>/brief.md`, the candidate's
  bytes, and `bets/<milestone>/basis.tsv`: `layup run: the brief of the bet <milestone> (session <id>)`. The write is the run's `write`
  through `Fenced`: a lease that is lost is a `*LostError` and promotes nothing, and the candidate stands. The stop between the candidate and the
  promotion that the integration test makes is the control comment of the session's end failing (`endPhase` returns it after the candidate's commit).
- **The retry.** A `refused file` of the checks, or a finding of the appetite, starts the next attempt with `bet-brief-refused`. A session
  whose output files `ReadOutputFiles` refuses (a missing file, an unlisted entry) ends `no-result`, not `refused file`; the shared phase
  machinery is unchanged, so that session's telemetry stands, `betBrief` returns the `*PhaseError`, and the next run starts a new session.
  A comment that fails after the candidate's commit also leaves the candidate in the records for the recovery.
- **The recovery** reads before it writes. A `brief.md` that exists is done before. Else the latest candidate under
  `bets/<milestone>/attempts/` (the last row of `sessions.tsv` that has one) is judged again, from its stored bytes alone (the session's
  files are removed after it): the checks (2) to (4) on the assembled brief, with code's parts found in it by their exact bytes and removed
  before the option-word rule, each part in its place as the inputs stand now (a candidate made on other inputs is refused with the rule `fixed`), and the appetite on the telemetry as it now stands. A pass is promoted with no new session; a refusal starts
  a session with its finding. The control comment of a session whose candidate the appetite refused says `done`, because its own class was
  done; the finding is in the next prompt.
- **The basis.** `basis.tsv` (the block `bet-basis`, [`records.md`](records.md#nfr-001--the-records-of-the-task-loop)) keeps what row 66j
  needs after the inventory write: the requirement IDs, the survivors and the dropped options, the facts with their classes, the questions
  still open and the verified heads. `betBrief`'s companion `briefBasis(milestone)` reads it with the brief, and gives the basis
  without filtering again, and the appetite from the kept brief's line; a basis or a brief that is missing or malformed is an error (a brief is malformed by the rules of its headings, appetite line and architecture line; the limit and the links are not asked again).

### The step bet

Row 66j (task `T-2vva`, #366; the brief is row 66u's, the rendered task of a later bet is row 66t's, moved to `M2g` as a whole by the Operator on #417).
Everything below is **decided here** (R10) unless it cites another section. The step reads the records and the forge at each run, holds nothing, and each
stage keeps what an earlier stage wrote; it is `done before` when `bets.tsv` holds the row of its milestone with the verdict `bet`.

1. **Read.** `bets.tsv` is read by `ReadBets` (the block `bets`; a refusal names the line and the column, and nothing is written). The Shape tasks are the first
   register row of the kind `specification` and the first of the kind `architecture`; each must have a verified head (`verifiedStep`, which reads and drives
   nothing). **The milestone** is the first row's tag of the table under `## 9.` of the PRD at the verified head of the specification (`phaseTagsOf`), a
   temporary convention of `M2e` that does not say that a phase and a milestone are the same thing; it must be one path part, and it is pinned by `bets/<milestone>/`. **Plan amendment 1 of #366** (the coordinator, R10): the specification instruction asks for that table, `## 9. Rollout & phases` with `| Phase | Name | Ships |`, a separator row of three cells (hyphens, with a colon at either end or both) and one row or more of three cells with a tag (any other form gives no tag), and moves the success metrics to `## 8.` as the baseline's PRD form has them (the reader finds the table by the name of its heading); it asks for `## 6. Functional requirements`, `## 7. Non-functional requirements` and a paragraph that starts `**Non-goals.**` too, the places where the rendering of the bet reads and writes (decided here: without them `RenderPRD` finds no row and drops the marks out of scope).
   **The kind** is `first` when no row exists. A row of that milestone with the verdict `bet` is `done before`; one with `no-bet` is the row `done`, `waiting for the
   Operator: no bet`, at every restart, and no later step runs on it. Any other state of `bets.tsv` is the kind `later`, which `M2e` does not build: the step fails,
   `a later bet's rendered task is row 66t's (M2g)`, and writes nothing (the contract of `M2e`, fail closed).
2. **The brief.** `betBrief(milestone, "first")` of [the brief of the first bet](#the-brief-of-the-first-bet). `ErrBriefLinesUnset` and `ErrBriefOverLimit` are the row
   `done`, `waiting for the Operator: <the error>`, with no session; any other error is a `fail`.
3. **The post.** `CheckCommentBody` runs on the kept brief again, then `postOnIssue` posts it on the Intake issue (`issue.intake`, the milestone's issue) and copies it:
   the copy by hash, the comment of the App that has the bytes, the lease, the post, the read-back and the copy, the core that `PostComment` and `postBody` share.
4. **The answer.** Every comment of the issue that is not the App's is copied first. The approver is `approverFor(…, "bet", …)` over `pointBindings` (a lookup error is the
   step's error, never a fallback). **The answer is the earliest copy by numeric comment ID** that is of the issue, is the first copy (seen 1) and no App's, is by the
   approver's account, was created strictly after the brief's copy, and for which `ParseAnswer` gives no finding, with the requirement IDs, the facts, the survivors and the
   options of the kept basis (`briefBasis`) and the phase tags of §9. Any other comment is human input that answers nothing: the Operator's, the idea owner's when an
   approver is named, an edit, one before the brief, one with a finding. With no answer the row is `done`, `waiting: the answer of the bet on the brief <comment>`.
5. **The settled decisions** are the option lines of the answer for the decisions that have more than one survivor in the basis; an option line for a decision with one
   survivor answers nothing new and opens no task. **Validate before any write** (before the inventory, the rendering and the tasks; the copies already exist): for a `bet`,
   the constraint values of the chosen option's rows are mapped to requirement IDs (`REQ-NNN` and `NFR-NNN` kept when the register row names them, `R-NNN` mapped by
   `finalID` with its kind in `spec/draft.tsv`, `F-`, `Q-`, `Inv-` and `ADR-` values are no requirement), and a decision whose set is empty fails, `the decision <d>: its
   option <o> names no requirement` (an architecture decision serves a requirement), at every restart.
6. **The version.** A restart goes on with the answer that the latest version of the inventory names (when it is of this bet's milestone), its first copy and the account that wrote it, whatever comment or binding of the approver came after (decided here: round 2 of #366). For a `bet`, `spec/inventory-<n>.tsv` and `spec/inventory.tsv` are written in one fenced records commit by `NextInventory`, after the sha256 of
   version `n-1` is checked against the index. They are skipped when the latest version's rows name this answer's copy. The facts are those of the basis: a fact of
   `spec/spans.tsv` has the class of Intake and its text from the span, and **a fact of the answers (`Q-NNN`) has the class of its quote's fact, else `context`, and the
   answer as its text** (decided here; the approver may change the class with a line `fact <id>: <class>`). A class that the answer changes is recorded and nothing else:
   the new specification task that a changed class reopens, and the representation "waits for a bet", are row 66ab's (the Operator on #366; [the specification task that a changed class reopens](#the-specification-task-that-a-changed-class-reopens)).
7. **The rendered commit** (the first bet): the next record ID is `nextFactsID` over the records of `docs/facts/` at the verified head `H` (the one home, which the instruction
   of the specification uses for `{facts}`). `H` holds exactly one `docs/prd/PRD-*.md` and no numbered facts record, or the step fails. The commit is made on a work tree at `H`
   by `rules.Render`: the PRD (pre-image `H`) and `docs/facts/<ID>-numbered-facts.md` (`RenderFacts` exactly: no table, no index row and no line of `facts.sha256`, which
   would be a rule-path change that is no rendered record; a known limit), with the App's identity at the time the answer's copy was made, so one input gives one commit
   ID, and the message `layup run: the rendered commit of the bet <milestone> (copy <comment>)`. `CheckRendered` runs before the push. **The push is code's**: the event `push R
   task/<task>/<attempt>` with the session `—`, the push, the event `bound`, with the session `—` (`records.CheckEvent` admits that for these two kinds). A restart makes the same
   commit and finds the events: `bound` skips, a `push` with no `bound` pushes the same SHA again (an equal fast-forward is a no-op), a branch that holds neither `H` nor
   `R`, a missing branch and an ancestor of `H` included, is `errPushRefused`, a fail that keeps the tip, checked before the push of a pending commit (decided here: round 1 of #366). **The consumers of `bound`** that key on a session skip a `bound` with no session (the repairs of a push, `loopPosition`, the authors, the
   register snapshot of a head, the last bound of an attempt), and the driver of a Shape task does not settle a verified head again once the push of the rendered commit is pending or made (an event `push` or `bound` of code in the attempt), since the pull request's head can be `R` before the `bound`, since the
   pull request's head is `R`. Then `CarryVerify` posts `layup/verify` `success`, `verified head <H> plus the rendered commit`, once. The other three contexts on `R` and the
   close-out are row 66k's.
8. **The architecture tasks.** For each settled decision in name order, if no task of the kind `architecture` after the first holds the binding
   `tasks/<task>/options/<decision>.tsv`, one is opened by the opening of a Shape task, parameterized (the order of [the task register](#the-task-register), with no reuse by
   kind): `requirements` as step 5, `predecessors` the two Shape tasks, the issue `LAYUP task <ID>: the architecture of the decision <d>` and the block below. The register commit
   writes the row and the binding (the chosen option's rows, the block `options`) together, and the candidate row, the parameters and the binding bytes are checked before the
   issue. The known window between the issue and the row is that of the opening of a Shape task. Row 66k reads the binding when it runs the task (`PromptIn.Options` is the chosen option as the only survivor).
9. **`bets.tsv` last**, in one fenced records commit: `brief` and `brief_sha256` of the posted copy, the appetite of the kept brief (`—` for `no-bet`; the decimals of the brief's line, `5` written `5.0`), `answer`, `options`
   as `<decision>=<option>` sorted by decision (`—` for none), the verdict, the approver's account and `decided`, the time the answer's copy was made, so a restart writes the same row.
   For `no bet` stages 6 to 8 do not run.

```text bet-issue-architecture
LAYUP task {task}: the architecture of the decision {decision}.

Requirement IDs: {ids}
Chosen option: {option}, recorded in tasks/{task}/options/{decision}.tsv of the records.

Definition of Done: docs/spec/ and docs/adr/ hold the decision {decision} as decided with the option {option}, and nothing of the other options; docs/tasks/{task}.md is the task's file. The head passes the four layup/ statuses, and a verifier on a harness outside the authors finds nothing material in scope.
```

**Known limits.** The brief's "no recommendation" for two survivors rests on code's presentation and the option-word rule (row 66u). The milestone is the §9 tag convention above.
A bet that changes a class does not reopen the specification here (row 66ab, [below](#the-specification-task-that-a-changed-class-reopens)).

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
  plan ("render from copy N") that a plan-review and a verifier session check (row 66t, moved to `M2g`; `M2e` fails closed at a later bet, [the step bet](#the-step-bet)).

### The specification task that a changed class reopens

Row 66ab (task `T-4naq`, #428; the plan review is the comment 6262491026 and the Operator's decisions the comment 6262589337 of #428). **Decided here** (R10) unless it cites another
section. [The step bet](#the-step-bet) records a class that the answer changes and nothing else; this section opens the work that the class asks for
(`architecture.md` §7 step 7). It is a primitive of `internal/run`, `(*state).openReopenings(ctx, spec)`, and `reopenedTasks(rows)`, the set of the tasks that
`spec/reopenings.tsv` names. Row 66k calls the first after the event `merged` of the specification task, reads the register again, and puts the reopened tasks
after the specification and before the architecture tasks in the merge order, in the order of the rows of `spec/reopenings.tsv`; the call, the order and the drive of
a reopened task are 66k's.

1. **When.** `openReopenings` reads the events of the task `spec`, the first row of the register of the kind `specification`. With no event `merged` it opens nothing and
   writes nothing, not even a call of the forge. Then the index `spec/inventory.tsv` must hold exactly one version, whose sha256 is checked against the index first; an
   absent index is an error, and a second version fails closed, `the reopening of the classes of a later bet is not built in M2e`, as the step bet refuses a later bet.
2. **The selection** (Q2). The rows of version 1, in its order, whose scope is `in`, whose class is `need`, `constraint` or `measure`, and whose class **differs from
   the class of the fact at Intake**: for a fact of `spec/spans.tsv` its class there, for a `Q-NNN` the class of the fact of its quote, else `context` (`factsOf`, the
   inventory's own source). Comparing the version with Intake **selects**; it never **deduplicates**, which is the table of item 3. A fact changed from `need` to
   `constraint` is selected though the PRD already covers it, because `architecture.md` §7 step 3 asks for a non-functional requirement for a constraint; a fact marked out of scope,
   and a fact whose class the answer does not change, are not selected.
3. **The record `spec/reopenings.tsv`** (the block `reopenings` of `records.md`; owner `internal/run`: `ReopeningsSchema`, `ReadReopenings`). One row for each selected fact;
   **all the rows of a version are written in one fenced records commit before any call of the forge**, and a row is never changed. The columns are `n` (the
   version, always 1 in `M2e`), `fact`, `class`, `token`, `requirement`, `task` and `copy`, and the reader's checks are in `records.md`. At every call, **the call with an empty selection included**, the rows are
   **compared with the selection**: the same facts, in the same order, with the version, the class of version 1, its `copy` and the token of item 4; any difference is an error and
   nothing is opened, so a row cannot be read as an identity of another version, fact, class or copy.
4. **The token.** `F-NNNN#<n>` for the fact `n` of the spans, where `F-NNNN` is the one numbered facts record, a regular file `docs/facts/F-NNNN-numbered-facts.md`,
   at the merge commit that the event `merged` names (none or two is an error): the record the specification's rendered commit wrote (`PSBRecords` of the full check),
   pinned, never predicted as the instruction of the specification predicts it before the bet. A fact of the answers is its own token, `Q-NNN`.
5. **The identities** (Q3, Q4, Q5; decided here). The row is the identity: the allocation of every row of the selection, in its order, is made once, in the commit of
   item 3, and never again. A task ID is drawn as `drawTaskID` draws it, again until no row of `sessions.tsv`, of `tasks.tsv` or of `spec/reopenings.tsv` holds it and
   none of the draws of the same batch is it. The requirement IDs have **one number space** for `REQ-` and `NFR-`: the next number is one more than the highest
   `NNN` of every `REQ-NNN` and `NFR-NNN` among `finalIDs` of `spec/draft.tsv`, the `requirements` column of the register and the `requirement` column of
   `spec/reopenings.tsv`, written with three digits at least; a `need` takes `REQ-` and the next number, a `constraint` `NFR-` and the next number, in the order of the rows. A
   `measure` has **no ID**: its `requirement` is `—`, and its task names the requirements of the row of `spec` (the register row is never changed, so they are the same at every call).
6. **The opening.** For each row, in order, whose `task` has no row in the register, the opening core `openTask` runs with the reserved ID, in its own order
   ([the opening of a Shape task](#the-task-register), whose steps it keeps: the permissions, the dry `registerFiles` and `parameterFiles` on the real candidate, the lease, `OpenIssue` and
   `registerTask`): the kind `specification`, the `requirements` of item 5, the `predecessors` `[spec]`, the issue of the block `reopen-issue`, **no reuse by kind** (the Shape
   task of the kind is never returned for it) and no opening pin. A row whose task is registered is checked, not opened: its kind, requirements and predecessors are the
   reservation's, else the call is an error that names the task. It gives the reopened tasks in the order of the rows.
7. **Recovery.** A stop after the allocation commit and before an issue leaves the rows, and the restart opens one issue for each task. A stop after the registration
   of a task writes nothing more for it. **Known limit** (the Operator on #428, as for the Shape tasks): a stop between `OpenIssue` and the register commit leaves the
   first issue on the forge, and the restart opens a second under the same task ID and the same requirement IDs and writes one row; no call of the forge lists issues. A restart after the
   opening leaves **one registered task** for each row.
8. **The pending prose** (Q7). The requirement is written **under the heading `Waits for a bet`, as one line of plain text and never as a table row**, in the PRD of the target, so
   that the pinned `prd-lint`, which reads the lines that start with `|` after white space, does not read it ([the full spec check](#req-012--the-full-spec-check) reads the
   heading, whatever its level and case). A `need` or a `constraint`: `<ID>: <the requirement, one sentence> Facts: <token> F-<record>. Acceptance criterion: <the
   criterion>.`, with a section `## <ID> — <title>` in `docs/spec/` (the item `spec-sections` reads the ID of the task's row); a `measure`: `Success criterion: <the criterion> Fact:
   <token>.` and no section. It gets its row of the tables at the next bet (`architecture.md` §7 step 7).
9. **The instruction and its inputs.** A task that has a row in `spec/reopenings.tsv` and the step `specification` get the block `instruction-specification-reopened`
   in place of `instruction-specification`, filled with the row's `{fact}`, `{class}`, `{token}`, `{requirement}` (`—` for a `measure`) and `{record}`, the numbered
   facts record of item 4, **found at the base of the session and required to be the token's**. Every step of such a task has `spec/reopenings.tsv` among the
   `requirements`, after the draft items, and the **files of `docs/prd/` whether or not a table row holds the ID** (a pending requirement has none); the
   plan, the plan review and the verification see the same list. Code reads the row of the task from the file, the session never.
10. **The exemption of the full check** (Q1; the Operator on #428): [the full spec check](#req-012--the-full-spec-check).
11. **Known limits.** A reopened task is not redone by a moved base in `M2e`: it has no opening pin, so `redoIssue` is an error before any write (the task is brought current by the
    clean merge of row 66x, or it waits). The drive of a reopened task by `FinishTask` and by the driver of a Shape task, which both take every task of the kind `specification` for the first one
    (the rendered commit of the first bet, and the pinned `prd-lint` of a specification that waits for its bet), is row 66k's, which tells the two apart by `reopenedTasks`. A later
    bet's reopening, and the new row of a table for a pending requirement, are `M2g`'s.

```text reopen-issue
LAYUP task {task}: the {what} of the fact {fact}, re-classed by the bet as {class}.

Requirement IDs: {ids}
Reopened after the merge of the specification task {spec}: the row of the fact {fact} in spec/reopenings.tsv of the records, the token {token}.

Definition of Done: docs/prd/ holds, under the heading "Waits for a bet" and not in a table, {line}; {sections}docs/tasks/{task}.md is the task's file. The head passes the four layup/ statuses, and a verifier on a harness outside the authors finds nothing material in scope.
```

Code fills `{what}` (`requirement` for a `need` or a `constraint`, `success criterion` for a `measure`), `{line}` (`one line that starts with {requirement}: and cites {token}`, or `one line that starts with Success
criterion: and cites {token}`) and `{sections}` (`docs/spec/ holds a section whose heading holds {requirement}; ` or nothing); `{ids}` are the task's requirement IDs, `{spec}` the specification task, `{fact}`, `{class}` and `{token}` the row's.

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
| `needs` | a `need` or `constraint` fact has no requirement, and it has no out-of-scope mark, no requirement marked "waits for a bet" and no exemption (below) |
| `measures` | a `measure` fact has no success-criterion row that cites its ID, and it has no out-of-scope mark, no requirement marked "waits for a bet" and no exemption (below) |
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
with the reason `no specification yet`, and the status is `success`. A head is judged against its own tree, whatever
the kind of its task; the architecture task's head holds the specification's PRD through its base (REQ-012 — Shape). After the first
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
numbered fact and is no failure of `covers` when the head holds the record as a
direct file `F-NNNN*.md` of `docs/facts/` (the set `prd-lint` reads, `FullInput.Records`);
`F-NNNN#n` of such a record that numbers no fact, and a record the head lacks, fail;
any other token fails `covers`. **Decided here** (task `T-c6sx`, #416, round 1, R10): this is
what lets a specification that waits for its bet pass `covers` while the pinned `prd-lint`
needs an existing `F-NNNN` in each Facts cell. `measures`
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

**Decided here** (task `T-4naq`, #428, row 66ab; the Operator's decision 6262589337 on Q1 and condition 1 of the plan review): after the first bet every rule judges the latest
version, so a fact that the bet re-classed into `need`, `constraint` or `measure` has no requirement at the rendered head of the specification, `needs` or `measures` fails
there, and the specification could never merge, nor reopen anything. `needs` and `measures` therefore skip the facts of `FullInput.Reopening`, which `specInput` computes for the
task whose head is judged (`SpecInput.Task`): **the re-classed facts** (the selection of [the specification task that a changed class
reopens](#the-specification-task-that-a-changed-class-reopens), item 2: version 1 against Intake) **except those whose row of `spec/reopenings.tsv` names a task that has merged or
the task whose head is judged**. The states: before the bet nothing is re-classed, and the set is empty; from the bet to the allocation every re-classed fact is exempt at every head, so
the specification's rendered head passes and merges; from the allocation to the merge of a fact's task the fact is exempt at every head but the head of its own task, which **fails**
`needs` or `measures` unless it holds the pending line of the fact under `waits for a bet` (a missing output is refused, never exempted); with two re-classed facts, the head of the first task holds only
its own line and passes while the second's fact is exempt, and the head of the second, judged after the first has merged, has the first's line in its base and holds its own; from the merge of a fact's
task the fact is no longer exempt, so a head that removes its line fails. The names are those of the universe (a number, or `Q-NNN`). A version that cannot be read gives an empty set, and
the rules do not run, as for any head. **The own output is read by the stored identity** (round 1 of #428, F1): at the head of a reopened task `needs` (for a `need` or a
`constraint`) and `measures` (for a `measure`) also require, from `FullInput.Pending` (the fact and the requirement ID of the row of the task in `spec/reopenings.tsv`), a line under the heading that
cites the fact and, for a requirement, holds the allocated ID as a whole word; a requirement of the table that already covers the fact (a need changed to a constraint) or the old requirement's ID
is not that output. **The pending output of a merged reopening stays required** (round 2 of #428, F6): `FullInput.Pending` also holds the identity of every reopening whose task has merged, so a later head of any task that removes its line fails, and one that keeps it passes. The detail is `the reopened <class> fact <f> has no pending requirement <ID> under waits for a bet` (`pending success criterion` for a measure).

**Decided here** (task `T-cav3`, #412, row 66p; the input of the full check): code
builds the `FullInput` of a head, `Sessions.specInput(base, head, copy)` in
`internal/run` (row 66ab adds the task whose head it is, for `Reopening`), each field from one named source, and a caller's literal of a field is
never read. The records are read at the run's records base at the call, not at the
snapshot of the head's session; saving the table once for a head, and so not
rebuilding the input at a restart, is row 66f's, which calls the builder.

| Field | Source | A missing source |
| ----- | ------ | ---------------- |
| `Changed` | `target:` the paths of `git diff --name-only` between `base` and `head` in the run's clone | an error |
| `SpecFiles` | `target:` the regular `.md` files (type blob, mode `100644`) directly in `docs/spec/` at `head`, by base name (the selector of `spec-sections`, row 66n) | an empty map |
| `PRDs` | `target:` the regular files directly in `docs/prd/` at `head` whose base name starts with `PRD-` and ends in `.md`, by base name | an empty map |
| `PSB` | `records:start/problem-statement.md`, the bytes that the spans index | an error |
| `Spans`, `Questions`, `Answers` | `records:spec/spans.tsv`, `records:questions.tsv`, `records:intake/answers.tsv`, each by its reader | an error |
| `Index` | `records:spec/inventory.tsv` by `ReadInventoryIndex`; written at the first bet | not an error: `nil` is "before the first bet"; a file with a header and no row is an error |
| `Versions` | `records:spec/inventory-<n>.tsv`, the exact bytes, for each `n` of `Index`; the builder never verifies a hash (the rule `inventory` does) | an error |
| `Register` | `records:tasks.tsv`, the raw bytes; not read by `ReadRegister`, which refuses a task with no requirement, the case the rule `tasks` reports | an error |
| `Merged` | the tasks of the register whose `records:tasks/<task>/events.tsv` holds an event `merged` | no file is not merged; a refused file is an error |
| `PSBRecords` | the numbered facts record (below); before the first bet and with no such file, the predicted ID | none is the predicted ID before the bet, else empty |
| `Records` | `target:` the direct regular files of `docs/facts/` at `head` named `F-NNNN*.md`, by record ID (`T-c6sx`) | an empty list |
| `AnswersRecord` | none: `""` | not applicable |
| `Copy` | the caller's (the bet's copy, row 66j) and passed through unchanged | `nil` keeps the `fail` row of `classes` and the `ran` row |
| `Pending` | the row of the task whose head is judged in `records:spec/reopenings.tsv`: its fact and allocated requirement ID (the output that the head must hold) | none for a task that no class reopens |
| `Reopening` | `records:spec/reopenings.tsv` with the Intake records, the latest version and `Merged`, for the task whose head is judged (`SpecInput.Task`; the exemption above) | an absent `spec/reopenings.tsv` is no reopened task; a file its reader refuses is an error |

- **`PSBRecords`.** The numbered facts record of the problem statement is the file
  `docs/facts/F-NNNN-numbered-facts.md` at `head`, a regular blob (mode `100644`)
  directly in `docs/facts/`. Its producer is row 66j: the rendered commit of the first
  bet writes it, and `NNNN` is one more than the highest `F-NNNN` of `docs/facts/` at
  the base of the rendered commit. The builder reads it: no such file at `head` gives
  an empty `PSBRecords`; exactly one gives `["F-NNNN"]`; two or more are an error whose
  source is `target:`. The correspondence: the fact `n` of `F-NNNN#n` is the line
  `n. ` of that record, which `RenderFacts` writes with the fact ID of the confirmed
  inventory as its number, and that ID is the number `n` of the span in
  `spec/spans.tsv`. **The predicted ID** (`T-c6sx`, #416): with an empty `Index` and no numbered facts file
  at `head`, `PSBRecords` is the ID one more than the highest `F-NNNN` of `docs/facts/` at
  `head`, so the citations `F-NNNN#n` of the specification that waits for its bet resolve. The
  prediction holds: no task writes `docs/facts/` (a rule path), nothing merges before the bet,
  and row 66j computes the ID of its record at the same head with the same function
  (`nextFactsID`, which the instruction of the specification also uses for `{facts}`). After
  the first bet the file exists or the list stays empty, as `T-cav3` decided.
- **`AnswersRecord`.** It stays `""`. The answers record of the setup (S04) numbers
  `S01-stack` first and only the answered questions after it, so its fact `N` is not
  `Q-00N`; with it the check would resolve a wrong citation. A target PRD cites an
  answer as a bare `Q-NNN`, which resolves through `intake/answers.tsv`; the pinned
  `prd-lint` still needs an `F-NNNN` citation for the requirement.
- **Errors.** The builder writes nothing, and on any error returns the zero input. A
  `base` or `head` that is not a full SHA is refused before any read, as a `SourceError`
  whose `Source` is `base` or `head`, a caller's error; any other `SourceError` has the
  source `records:<path>` (a Go error of the run) or `target:<what>@<sha>` (a
  posted `error` of the status), and the status is row 66f's. A failed read of a
  file that may be absent (`spec/inventory.tsv`, an `events.tsv`) is an error; only its
  absence is not.

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
  through the task loop of this file, up to the verified head of each, the
  specification first, then the architecture task, opened at that point at the
  specification's verified head (REQ-012 — Shape; a stop or a wait of the specification
  task is the step's stop or wait, and the architecture task is not opened), which
  `architecture.md` §8 asks before the bet; it also writes the
  options of each decision and asks the questions left open. It is **done** when each
  decision's table in `spec/options/` has a survivor and each Shape task has a
  verified head; a decision with no survivor leaves the step with the detail
  `shape: no survivor`, which is not "done before", so a restart stops there again
  and never goes on to `bet`. **The order inside the step, decided here** (R10, task
  `T-3f94`): (1) the questions left open and the intent wait; (2) the options session, the
  filter, the one rewrite and the stop; (3) the specification task, then the architecture
  task, each to its verified head; (4) done. The step reads its records and the forge at each
  run and holds nothing; it stops at the first stage that waits or fails, and each stage keeps
  what an earlier one wrote. It is `done` with `<k> decisions, each with a survivor; <T-x> and
  <T-y> verified` when this run wrote something, and with `done before` when it wrote nothing.
  A wait of a Shape task is the row `done` with the driver's `waiting for the Operator: …`; a
  stop is `fail`. With several survivors nothing is chosen (the bet decides, row 66j). The
  Shape tasks stop at the verified head: the ready rule, the close-out and the merge are row
  66k's.
- **`bet`** posts the brief, waits for the answer, copies it, writes the version of
  the inventory, pushes the rendered commit to the specification's branch (with its
  carry-over of `layup/verify`), opens the `architecture` tasks of the settled
  decisions (rows of `tasks.tsv`; each with `requirements` the requirement IDs among
  the `constraint` values of the chosen option's rows, and a decision whose table
  names none is refused, as an architecture decision serves a requirement), and writes
  the row of `bets.tsv` **last**. It is done when that row exists with its verdict, so
  a stop before it resumes inside `bet`, and no later step runs on a half-made bet.
- **`loop`** merges the Shape tasks in the merge order after a bet (the specification
  first, then the tasks it reopens with `openReopenings`, after which it reads the register again), then runs the tasks that the register holds; the build tasks are `M2g`'s.

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
  is not opened again, a status already at the state and the description is not posted again, and a
  merged pull request is recorded with its `merged` event.
- **The branch of a session that has a `push` and no result.** A `push` event with
  no `bound` and no `refused` after it is a forge write whose result is not
  recorded ([`session.md`](session.md#the-push-and-the-bind)). **Decided here**
  (task `T-syxm`, row 66v): `(*Sessions).RepairTask(ctx, task)`, in the records and
  the forge alone (never `bets.tsv`), repairs a task's sessions before a driver reads
  the task's position, and each driver of a task calls it first: the driver of the
  Shape tasks (`AdvanceShapeTask`, so the step `shape`, with no bet) and the act of
  `loop`. A stop of the driver on an unresolved push stays as a guard that the hook
  makes unreachable. For each such `push`, it reads the branch of the push event
  (`task/<task>/<attempt>`) of the forge. A head equal to the SHA of the `push`, byte
  for byte, writes the missing `bound`. Any other result writes `refused` with the
  detail `push-lost`: the branch is absent, or holds another SHA, a descendant of the
  pushed one included (code also advances branches, so being the only writer proves
  no equality). A read that fails is an error and writes nothing. A `push` of an
  attempt that is no longer the latest gets the same `bound` or `refused` and no new
  attempt. **The replacement rule:** after those writes, when the latest attempt (the
  largest `attempt` event) is still open (no `closed` and no `rebased` of it) and holds
  an event `refused` with the detail `push-lost`, one new attempt starts at that
  attempt's base (`StartAttempt`). The new attempt holds no such event, so the rule
  starts one attempt however many pushes were lost and however often the repair runs,
  and a restart between the commit of `refused` and the commit of `attempt` still
  starts it. The hook then closes the older attempts' pull requests and reposts the comments of the task's sessions (next item;
  the close is O22 of [A task whose base moved](#a-task-whose-base-moved-the-carry-and-the-redo), row 66x).
  Each stage keeps what an earlier stage committed.
- **The comments of a session** ([`session.md`](session.md#a-comment-for-a-session)):
  from `M2e` the comment of a task session is posted on the task's issue and no
  longer on the control issue, with the same first line; the control issue keeps the
  comments of the probes and of the phase sessions. The issue is the register's
  (`tasks.tsv`): a missing register, a missing row or a register that its reader
  refuses is an error, and no comment is posted anywhere (no fallback to the control
  issue). Before it posts, `layup run` reads the issue's comments, and a comment of
  its own App with the same bytes is not posted again; the lease is read right before
  the post. The hook posts the missing line of each session of the task that has a
  `result`, from the events. A line changed by a repair (`…: done, refused
  push-lost`) is a new comment, and the old one stays.
- The lease, the fencing and the exit codes are those of [`run.md`](run.md); a lost
  lease stops the run before any write of a step.

### The driver of a Shape task

**Decided here** (R10, task `T-jy6f`, #415, row 66r; plan review conditions 3 and 4 of 6250262119 and plan amendment 1).
`(*Sessions).AdvanceShapeTask(ctx, ShapeTaskInput)` takes one Shape task from its register row to a verified head, in the records and
the forge alone. It calls `RepairTask` first, then reads the task's records and acts on the first position of the table of
[The task loop](#req-005--the-task-loop) that they show not done, then reads again. It gives `Advance{Task, Verified, Head, Wait, Stop}`: a
verified head (`Head` is the bound head that the verifier found nothing material in), a wait for the Operator, or a stop. A stop and a wait repeat on the next call, which reads the same records: there is no
automatic retry. `loopPosition` is a pure function of the events, the rows of `sessions.tsv`, the status of each kept result and the
parsed review records of the valid handoffs, so each position is told from the records alone, and a stage whose session has ended is
**finished from its kept files** and never run again:

| Position | Records show | Act |
| -------- | ------------ | --- |
| attempt | no event `attempt` | `StartAttempt` at the base the caller gives: the base of the run for the specification, the specification's verified head for the architecture; the `base` of the event `rebased` of a task that a redo made ([A task whose base moved](#a-task-whose-base-moved-the-carry-and-the-redo)) |
| plan, plan review, developer | no session of the stage, or its last session has no event `result` (a stopped run) | the milestone cap, the pair, the prompt file, `TaskSession` |
| finish | the last session of the stage ended `done` with the status `completed` and has no event `handoff` | `PostComment`, then `RecordHandoff`; both read before they write, so an interrupted post or handoff resumes |
| to the verifier | the developer's session is bound and has no handoff | `RecordHandoff` with the commits and the task file |
| ship | the handoff is valid and no verifier session of the attempt has the bound head as its base | the draft pull request, the four statuses, the acceptance of a specification that waits for its bet, then the verifier |
| verdict | the verifier's handoff is valid | `LinkReviewRecord`, the four statuses with `layup/verify`, then the verified head, the next attempt or a wait |
| question | the latest attempt has a session that ended `done` with the status `needs_context` and no attempt follows | the stage P1 to P8 of [A question during the work](#the-checks-of-a-handoff): one act, then a fresh read; the position comes before the stages of the plan, the plan review, the developer and the verifier, and replaces the stop for that status (**decided here**, task `T-r7b9`, row 66z) |

- **Which sessions push.** A task session of a step that reads (`plan-of-a-task`, `plan-review`, `verification` and, from row 66z, `answer-to-a-question`) has no push: a
  session that commits nothing has its head equal to its base, which would pass every check and be bound, and a first `bound` would then
  say nothing about the developer. The rule is in `TaskSession`, by the step of the prompt file, so a caller cannot forget it.
  **Rejected:** a no-op push supplied by the driver (not the production path); the role as the key (the plan role is also the owner role
  of an environment question). Only the developer's session is ever bound; a session that has a `push` and no `bound`, or none of them
  after it ended `completed`, is a stop (the repair of a push is `RepairTask`'s). The close-out session of row 66y is bound too, after the verifier's last handoff: `FinishTask` tells the verified head over the events up to that handoff.
- **The plan review and the verifier**, both of the role `qa-engineer`, are told apart by the order of the events, never the clock: a
  session that starts after the last valid plan handoff (a handoff of a comment: **decided here**, task `T-r7b9`, a question or an answer handoff of
  a session of the plan role is no plan handoff) and before the first `push` of the task is the plan review; one that starts after a
  `bound` of the developer's session of its attempt, with that SHA as its base, is a verifier. Any other `qa-engineer` session is a stop.
- **The base of a session** may be off the default branch: the bound head that a verifier reviews, or the specification's verified head
  that is the architecture task's base (plan amendment 1 of #364, #432). Before **each** session of a Shape task the driver sets the
  branch `layup/base/<task>/<attempt>` of the run's clone to the base of the session (`git.UpdateRef`, which refuses a ref outside
  `refs/heads/layup/`), and the session clones that branch instead of the default one; the default branch is never moved.
- **The cycle cap of a Shape task is 1.** A Shape task cannot touch a gate: `.github/` and the check scripts are rule paths, refused
  before any push, and the cap is 1 unless the change touches a gate (`docs/engineering-discipline.md`). The rounds with findings are the
  valid handoffs `verifier → developer` of the task; a gate failure starts an attempt and counts no round. `nothing material in scope` is
  the verified head. `material` posts the status `layup/verify` as a failure at the reviewed head and, if the rounds number the cap or
  fewer, starts the next attempt at the base of this one, with the findings of that record; otherwise it is a wait
  (`waiting for the Operator: <task>: the review-cycle cap 1 is reached with findings open`). `not mergeable, findings recorded` is
  always a wait, and never another cycle. The findings of an attempt are those of the record of the attempt before it, else the rows of its
  head's gate table; never an older record.
- **After the four statuses of a head,** an `error` of any of them, or `not-active` of `layup/verify` (posted before the stop,
  `NFR-004`), is a stop, and no attempt starts. Otherwise a failure of `layup/gates`, `layup/spec` or `layup/rules` starts the next
  attempt at the base of this one, with the findings of the gate's table, then those of `layup/spec`, then that of `layup/rules`
  ([The attempts](#the-attempts)); a stop or a Go error of the check findings comes before it. **Decided here** (R10, task `T-tnxn`, #427,
  row 66aa): the decision is the act `ship`, the single function that decides what follows the statuses of a bound head. It is
  not a position of its own: a stop between the post of the statuses and `StartAttempt` leaves the records as they were (the
  verifier has not started and the attempt is still the latest), the next call reads the same position (`ship`) and the statuses
  already posted, and starts exactly one attempt. The next developer's prompt file is built from the records of the failed
  head, the gate table and then the check findings (`CheckFindings`), so the findings are never kept apart from them.
  **Reasons:** retry classification belongs with the checks of the current bound head; `RepairTask` runs once at the entry of
  `AdvanceShapeTask`, before its first position is selected, and repairs the sessions of the task (it can start an attempt and post
  comments), so the retry is not placed there. The driver is reached during Shape through the step `shape` (row 66h) and after
  the bet through the step `loop` (row 66k), with no `bets.tsv` read.
- **A specification that waits for its bet** is accepted by code before the verifier runs: the pinned `prd-lint.sh` of the target at the base is
  run on the `docs/prd/` and `docs/facts/` of the bound head by `spec.WaitingSpecAccepted`; a head that is not accepted is a stop that names
  the lines that are not an empty MoSCoW or Phase cell. **Reason:** that failure is the target's own job, not a `layup/` status, so no
  status carries it to an attempt, and the head must be judged before the verifier is asked.
- **A question is a position of the driver, not a stop** (**decided here**, task `T-r7b9`, #415 and #426): a session of the latest attempt that ended `done`
  with the status `needs_context` is the position *question*, which the stages of [A question during the work](#the-checks-of-a-handoff) take, in the
  same loop and after the same `RepairTask`; the sessions that asked in an attempt that a later attempt follows are not the last session of the plan
  or of the plan review, and the answer session of an attempt is a session of none of the four stages. Only the statuses `blocked`, `failed` and `decision_needed` are the stop
  `the session <id> of <step> ended <status>`.
- **A developer session that ended with the status `completed` and no `push`** is a stop (`<task>: the session <id> ended with no push`):
  the push hook did not run, and its resume is the step loop's. Every other end that is not `completed` (a refusal, another class or
  status) is a stop that names the session.

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
| The input of the full spec check (row 66p, `T-cav3`) | unit; integration (a real `git` tree and the records of a target after Intake and the Scaffold) | `specInput` builds each field of `FullInput` from its named record or from the tree at the head: a missing required record, a record its reader refuses or a failed `git` call is an error with the zero input and no write, an absent `spec/inventory.tsv` is "before the first bet", `PSBRecords` is read from the one numbered facts file at the head (none is empty, two are an error), `AnswersRecord` is `""`, and the copy is the caller's. |
| The specification task that a changed class reopens (row 66ab, `T-4naq`) | unit; integration (the records, the `httptest` forge, real git, scripted harnesses, `state.sessions()`) | After the specification's event `merged`, one specification task is opened for each fact that the bet re-classed into `need`, `constraint` or `measure` (a fact changed from `need` to `constraint` included, a fact out of scope or unchanged not), with the requirement IDs of one number space allocated once in `spec/reopenings.tsv` before any call of the forge, no reuse by kind and the specification as its predecessor; nothing is opened before the event, and a restart after the allocation or after the registration leaves the same IDs and one registered task for each row (the known window between the issue and the row leaves a second issue); a lost lease opens no issue; the reopened prompt carries the PRD, the row and the pinned numbered record, and a session that follows it gives a valid handoff and a green `layup/spec`; and a fact whose task has not merged is exempt from `needs` and `measures` at every head but its own task's, so the rendered head of the specification passes and a reopened task without its line fails. |
| The question and answer comments (row 66o, `T-gyvo`) | unit; integration (the committed rows of a real end) | The question and answer comments are composed from their rows, byte for byte, with the row's `n` in the ID, and refused with no bytes for a row or a start row off the reader's schema, types or rules, a row of another kind, an ID off its form, or a body that `CheckCommentBody` refuses; posting them is the step loop's (row 66z, a size child of row 66k). |
| A question during the work, in the loop (row 66z, `T-r7b9`) | unit; integration (the records, the `httptest` forge, real `git`, scripted harnesses, `state.sessions()`) | A session that ends `needs_context` during Shape, with no bet, has its question posted and copied on the task's issue, its handoff recorded, and, when all its questions go to one owner role, the answer session `answer-to-a-question` of that role (read-only, at the base of the last asking session) and the answer comment posted and copied (each answer row selecting its question by its column `n`), the handoff back recorded and the next attempt started; a human-routed question and questions of more than one owner role end in the wait `waiting for the Operator: <task>: question <id>` with no answer session; the stages are told from the records and repeated safely (a restart after each act writes each comment, handoff, session and attempt once), an invalid recorded handoff and a question or answer row that cannot be composed stay stops, an answer session that fails, asks again or answers too few or too many rows is a stop, and the plan and plan review that asked run again after the answer; an open material finding keeps `Position` uphill: the source is the last verifier, told by its position and head binding and qualified after its selection, so an unbound session of the role, a verifier that asked or has not ended, and a verifier with no review record, no valid comment handoff or a record of another head never clear a finding and never let an older clean record stand in their place; a verifier that asked at the head it was bound to (not the base of the attempt) is answered at that head. |
| The transition table | unit | The default rows are the thirteen of the table (the twelve of `architecture.md` §8 and the row of the target's architecture task, which §8 names a task and gives no row); a replacement table is refused for a missing default pair, a kind or an item not in the lists, and a row of code's own handoff (`code (rendered task)`) that lists `schema`, while another row keeps its `schema`, and the steps that are `layup run`'s own acts need no row. |
| The replacement of the transition table at the boundary | integration | With the stand-in forge and real git: `set transitions <sha256>` names a copy of a comment of the control issue; the new table is the bytes of that copy, written with the row `transitions` in one records commit, so its SHA-256 is the row's value; the old table stays in Git history and the next task keeps the new one; a hash of no copy or of a copy of another issue, a body that is no table (a code fence, a carriage return) and a table that fails the check each refuse with the item `bound` (`value` for a value that is no SHA-256) and write nothing. |
| `layup spec check`, the full check (inventory) | unit | Each rule of the table fails on its input and passes on its fixture; before the first bet a task that touches neither the PRD nor `docs/spec/` is `clear`; `did not run` is a failure. |
| `layup/gates` (inventory, `later-gate-status`) | unit; integration | Each row combination maps to the status by the table, `not-active` is never a success, a `pending` kind with a product path is `failure`; with the stand-in forge, the status is posted once for an unchanged state. |
| `layup/rules`, the basic rule | unit | A change to a rule path fails; a rendered record equal to code's rendering passes; one that differs by a byte fails; the register is read at the base. |
| The findings of `layup/spec` and `layup/rules` start a new attempt (row 66aa, `T-tnxn`) | unit; integration (the `httptest` forge, real git, scripted harnesses) | The fail rows of a head's table are findings `spec:<rule>` in the order of the rules, the failure of the basic rule is `rules:<reason>` with its path, a pass, a clear row and an error are none; the next attempt takes the gate findings, then the given ones; an expected failure of a batch task is none and an active gate row is one; a stop or a Go error comes before any retry. Through `state.sessions()`, a failed `layup/spec` or `layup/rules` at the head of a Shape task starts one new attempt at the same base whose prompt file holds the findings, a restart after the failed statuses starts exactly one, and the specification that waits for its bet starts none. |
| `layup/verify` and the independence | unit | The authors of a plan review and of a verification from the ledger; no harness outside them is `failure` `not-active`; a waiting verifier is `pending`. |
| The ready and merge rule (inventory, `later-gate-status`) | integration | With the stand-in forge, a pull request whose head has one context not `success`, no status, a status of another source, or an `error`, or a required job of GitHub Actions with no check run or a conclusion other than `success` (a ruleset that requires none, or a malformed required check, included), is not marked ready and does not merge; with the four `success` it is marked ready and merged at the head SHA, and a moved head is refused, its statuses are read again in the same call and no second merge is sent from the old read. |
| The body of the pull request links the review record (row 66q, `T-08rj`) | unit; integration (the `httptest` forge, the real `git`) | After a review record is posted and its handoff recorded, one edit replaces the body with the composer's body and the link of the comment, from the records only: the latest handoff of the verifier must be valid (an earlier valid one after a later invalid one is refused), the snapshot is the producing session's when the verifier bound the same head with another, a newer attempt with the register on the older pull request edits nothing, an equal body, a lost lease, a closed pull request and a forge error edit nothing or give an error, an edit that took effect and gave an error is not repeated, and a second call after a close-out and a restart makes no second edit. |
| The close-out commit and the rendered commit | integration | A close-out that changes only the task file and the completed log, on a parent that carries the verdict, is accepted and `layup/verify` is set on the new head; one that changes another file, or whose parent does not carry the verdict, is refused (`scope` or `verdict`, with the diff as a payload); an ordinary child of the verified head, a session of another task and a rendered commit above an unverified head do not carry; the rendered commit of the first bet gets the same carry-over, with no bound session and its pre-image pinned to the verified head, and the close-out after it is accepted. |
| The close-out and the merge of a verified task (row 66y, `T-ehb4`) | unit; integration (the `httptest` forge, real git, scripted harnesses, `state.sessions()`) | A verified task goes through the close-out session at the tier of `TierOf(plan, position)`, its checked push, `layup/verify` carried to the new head, the three other statuses of the head (also at the rendered commit of the specification task, with the inputs pinned apart for each head) and the ready mark to the merge at the head SHA and the event `merged`, and a second call is `Merged` and writes nothing; a pull request closed and not merged is the wait `waiting for the Operator: …` with no write of any kind, even with a missing session comment and an unresolved push, and a repair that starts an attempt is the handback `Attempt`; a merge that the forge holds gets its lost event; a close-out that is refused, or failed, is a stop and is not tried again, one that ended with no push resumes its checked push, and one with no result runs again; not ready is a wait and a moved head a stop; a restart after each write adopts it and sends no second session, status, mark or merge; no `layup/verify` is computed at the close-out head. |
| The steps after the Scaffold | integration | A restart resumes the `loop` step from the records and the forge: a comment, a pull request, a status or a merge that exists is not made again; a `push` with no result is bound when the branch holds its SHA and `refused` with `push-lost` when it does not; a task session's comment is on the task's issue; a merge whose `merged` event was lost passes `forge-check` with its detail (the tree check of task `T-xvvd`: a walk of the first parents through recorded and shown merges, and the setup head alone is no pass after the first merge) and `loop` records the event; `shape` with a decision that has no survivor is not "done before" and a restart stops there. |
| The bet step is atomic to its last act | integration | A stop after the inventory version and before the row of `bets.tsv` resumes inside `bet`: the rendered commit and each `architecture` task are not made twice, and the row is written last; a decision whose chosen option names no requirement is refused. |
| The tier and the context of a task (`session.md` deferral) | unit | The tier from the four fields of a plan (`large`, any non-empty list: reasoning; else execution), moved by the progress position (uphill to reasoning, downhill to execution); the prompt file lists the issue and plan, the rows and sections of the requirement IDs, the decision records they name, the handoffs and answers, and the rule files, in that order, each with its SHA-256. |
| The questions left open and the band (`intake.md` deferral) | integration | With the stand-in forge: an `open` and an `asked` question are asked of the owner role; an intent answer that parses (`B` below `U`) writes `budget.tsv` with the question as `source` and replaces a value that exists; one that does not parse leaves the question `open` and the run in its wait; a question still `open` at the brief is shown in "rabbit holes". |
| The approvers and the success criteria (`intake.md` deferral) | unit | The measure rows read from the intent answer; an approver row with `source` the answering comment ties to its point; a role `approver`; a point with no approver is the idea owner's; an account named by another spelling gets no second row; the revision's check that a requirement that covers an answer's fact names the quote's fact, and that a class change of a fact covered by both drafts names a fact of the answers. |
| The task register and the Retrospective | unit | A row with a predecessor that is not a task, a cycle of predecessors, or `—` in `requirements` is refused; the Retrospective step writes its event and no session starts. |
| The Shape filter (inventory) | unit | An option with a failed row, a missing row or a pass with no evidence is dropped; one survivor is chosen; two survivors give a brief with no recommendation, and an answer line `option <decision>: <option>` naming a dropped option or an unknown decision is refused and one naming a survivor is recorded in `bets.tsv`; every option dropped asks one rewrite, and a second full drop is a round with no progress. |
| The bet and its brief (inventory) | unit | The five headings, `brief.lines` unset is refused, a link that does not resolve is refused, two bets of 0.6 `U` each (the second brief), a cap above `band.U`, hours above `appetite.hours` and caps plus the spend of milestone 0 above `band.U` are each refused, and a cap below `band.B` passes, an answer by an account with no allowed role answers nothing, and the answer is copied. |
| The brief of the first bet (row 66u) | unit; integration | Code assembles the brief from the PRD at the verified head, the survivors and dropped options, the facts with their class and text, and the questions still open; the brief is refused for each rule of its checks (a heading, a link outside the set, an option word, a cap above `band.U`, the session's own spend) and kept when it passes; a refused brief is tried again with its finding; a stop before the promotion promotes the candidate with no new session; `brief.lines` unset starts no session; the basis is read back without filtering again. |
| The rendered commit (inventory) | unit | The commit may change only the MoSCoW and Phase columns, the non-goals and the numbered facts record; `layup/rules` passes a rendered record only when it equals code's rendering of the named copy. |
| The parameter change (inventory) | unit | A `set` is refused for a form that is not one line, a wrong role, an unknown name, a value of no type or out of bound, and the one that passes is applied as one replaced row, its default, bound and role kept; a set at or below the comment that its row records is superseded and is no change; a `prd-lint` head is accepted only when each failure line names an empty MoSCoW or Phase cell. |
| The parameter change at the boundary | integration | With the stand-in forge and real git: `layup run` writes `parameters.tsv` with its six rows in the commit of the first task (a table that exists is kept, one that its reader refuses refuses the task); a `set` is copied and, at the boundary, applied as one records commit when its author has the role of the row, and changes nothing before it; each refusal (`form`, `author`, `name`, `value`, `bound`) is reported with its item, changes nothing and is judged again at the next boundary; a set that is applied is not applied again, the later of two wins in the order of the comment IDs, and an edit of a comment is no set; a set of a Start row is recorded with the note that the run reads `start.tsv`; the write is fenced by the lease. The call between steps is row 66k's. |
| The frozen test list | unit | The list freezes at the first valid handoff to the verifier, a changed source hash is a finding, an amendment freezes it again, and a test that does not fail at the base makes the handoff invalid. |
| The gate jobs read the base (F-22) | integration | A pull request that edits `docs/gates.tsv` to pass a failing kind is judged by the base's file and fails; the catalog workflow reads the base; a pull request of a `batch/<point>` branch reads its own head's gate files and a pull request of any other branch does not; a change of a `config` path of the manifest waits for the code owner (shown on a real repository in the demo). The setup (task `T-jew3`, unit and integration): `.github/CODEOWNERS` of the four path classes with the Operator's login is in the tree of a new target, the ruleset sets `require_code_owner_review` with no required approval, and a target set up before the rule gets the update by its recorded ruleset ID, which keeps what the file never held ([`scaffold.md`](scaffold.md#the-code-owner-rule)). |
| A task whose base moved (row 66x, `T-89xf`) | unit; integration (the `httptest` forge, real git) | A task whose base moved with no overlap is merged by code, the merge is pushed as code, and `layup/verify` is carried over it (a fabricated merge that takes the base's version of a file both changed, or whose second parent is off the default branch, is not carried); a task that overlaps the base, or whose merge conflicts, is redone as one new task under the old row with the order close the pull request, close the issue, open the issue, one records commit, and a restart after each act completes it once; a failed git read, a lost lease and a merged pull request are errors; the older attempt's pull request is closed once. |
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
