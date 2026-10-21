# The packages

The Go packages of `layup`, their jobs, and the import rules between them. The
boundary gate of a Go target is "configured from the package table of the
approved architecture" ([`architecture.md`](../architecture.md) §6); this is
that table for LAYUP itself. Conventions: [`README.md`](README.md).

## NFR-007 — Go, the standard library only, and Git as the `git` program

Requirement: "LAYUP is written in Go with the standard library only, and calls
Git as the `git` program." Derives from `architecture.md` §1 and §2, ADR-0010,
ADR-0011 decision 1 and ADR-0013.

1. The module is `github.com/pharzam/layup` (the present `go.mod`), with one
   binary, `cmd/layup`, and packages under `internal/`.
2. `go.mod` has no `require` line. `go list -deps ./...` names only packages
   of the standard library and of this module (the criterion of `NFR-007`).
3. Only the package `internal/git` starts the `git` program, with
   `os/exec`; no other package imports a Git library or starts `git`. Its
   calls are [the calls of `internal/git`](#the-calls-of-internalgit).
4. Only the packages that the table below marks "starts a program" import
   `os/exec`.
5. No package of phase 1 depends on `net`, `net/http` or `crypto/tls`, by its
   own imports or through another package: the engine checks open no
   connection (`NFR-005`, [`gate.md`](gate.md#nfr-005--no-model-call-in-the-engine-checks)).
   From milestone `M2a`, only the packages that the column "Connects" of
   [the table of M2a](#the-table-of-m2a) names depend on them (by D8 below,
   through the adapter too), and one package alone imports them itself; `M3a`
   gives `internal/smartif` (phase 3) its rule
   ([the milestones](../plan/README.md#milestones)).
   The one package that imports them: `internal/forge/github`.

**Decided here** (task `T-esfe`, #127, row 22a of the plan): the test of the
package rules reads the one package of rule 5 from the line that starts with
"The one package that imports them:", its one code span, so the test holds no
list of its own (D7 of #79); no such line, two such lines, or a line with no
code span or with two, is an error. The three packages of rule 5 (`net`,
`net/http`, `crypto/tls`) stay a constant of the checker, as before.

**Decided here** (D6 of #79, K8): the rules bind the non-test Go files only. A
test file (`_test.go`) can import any package of the standard library and of
this module, and can start any program. Reason: `go list -deps -test` names
`os/exec` in each test binary, through `testing`, and the tests start `git` and
`go`.

**Decided here** (D8 of #79, K39): rule 5 counts each package that
`go list -deps` names for a package, not only its own imports. Reason: a
connection that a dependency opens is a connection of the package that imports
it.

### The table of phase 1

"May import" lists the packages of this module that a package may import; any
package of the standard library is allowed, except as rules 3 to 5 say. An
import that the table does not allow is a defect, and the boundary rule of the
future Go gate of LAYUP reads this table.

**The form of a cell** (decided here, D7 of #79): "Package" is one code span.
"May import" is one or more code spans with a comma between two; for none, the
cell is the character — and no code span. "Starts a program" is exactly one
code span at the start of the cell, whose first word is the program, and text
with no code span after it (a register row's cell is its own form, [below](#the-table-of-m2b)); for none, the cell is the word no and no code
span. Reason: the [test of the package rules](#the-test-of-the-package-rules)
reads this table, and holds no copy of it, which could differ from it.

| Package | Job | May import | Starts a program |
| ------- | --- | ---------- | ---------------- |
| `cmd/layup` | `main`: passes the arguments to `internal/cli` and exits with its code | `internal/cli` | no |
| `internal/cli` | parses the arguments, runs one command, maps its result to an exit code ([`README.md`](README.md#commands)); for `layup run`, reads the two host registers and the key file, builds the adapter and hands it to `internal/run` (task `T-mqty`), and from `M2b` reads `models.tsv` and `routing.tsv` and checks each harness credential (task `T-ysph`), and hands their rows and those of `host:prices.tsv` to `internal/run` for the step `probe` (task `T-nxe4`) | `internal/psb`, `internal/setup`, `internal/verify`, `internal/gate`, `internal/run`, `internal/forge`, `internal/forge/github`, `internal/route`, `internal/spec` | no |
| `internal/tsv` | reads and writes a record: checks the header row against a schema, the field count of each row, the types and the key; parses the `tsv-schema` blocks of `docs/spec/`, and compares a block with the Go schema of its record ([`README.md`](README.md#the-schema-block)) | — | no |
| `internal/git` | the one caller of the `git` program: [its calls](#the-calls-of-internalgit); in `M2d`, `CommitTree` (task `T-j4sd`); in `M2e`, `RemoteHead`, the read of a branch (task `T-syxm`), and `Merge` and `MergeBase`, the merge of the default branch into a task's head (task `T-89xf`) | — | `git` |
| `internal/psb` | the rules G1 to G5 and the gap table ([`psb-check.md`](psb-check.md)) | `internal/tsv` | no |
| `internal/catalog` | the stack catalog, embedded with `embed` ([`setup.md`](setup.md#the-stack-catalog)) | `internal/tsv` | no |
| `internal/work` | the work area of a target: its paths, and the schemas and the readers of `answers.tsv` and `record.tsv` ([`setup.md`](setup.md#the-answers)), which `internal/setup` and `internal/verify` share (K9); the record row `answers.sha256`, which `internal/setup` writes and `internal/standin` gives its stand-in; the texts of the four questions of S01, which the stop table, the answers record of S04 and the stand-in use (task `T-7s0y`); the schema of the table of `layup setup verify`, which `internal/verify` writes and S15 reads (task `T-d6q5`) | `internal/tsv` | no |
| `internal/gate` | runs the kinds of a gate manifest on a head ([`gate.md`](gate.md)) | `internal/tsv`, `internal/git` | `sh -c`: the gate commands |
| `internal/setup` | the step runner of `layup setup` ([`setup.md`](setup.md)); in `M2d`, its inputs of phase 2, which `internal/run` gives in Go (`Phase2` and `StepsPhase2`, task `T-ktgx`), the two rulesets of S13 (task `T-9v92`) and the commands of S13 for the Operator, with `set -eu` as the first line of `commands.sh` (task `T-56dj`; [`scaffold.md`](scaffold.md#req-002--the-setup-under-layup-run)) | `internal/tsv`, `internal/git`, `internal/catalog`, `internal/work` | no |
| `internal/verify` | the checks of `layup setup verify` ([`setup.md`](setup.md)), and `InvariantIDs`, the `Inv-<n>` entries of a `guardrails.md`, which the step `shape` reads (task `T-3f94`) | `internal/tsv`, `internal/git`, `internal/catalog`, `internal/gate`, `internal/work` | `sh`: the baseline's own check scripts |
| `internal/records` | the schemas and the row rules of the record kinds that phase 1 defines and later phases write: `telemetry.tsv` and `prices.tsv` of `REQ-011` ([`records.md`](records.md#req-011--the-telemetry-record), task `T-tmhw`), and `stalls.tsv` of `REQ-009` ([`records.md`](records.md#req-009--the-stall-record), task `T-dgy7`); and the records of Start of `M2a`, `start.tsv`, `approvers.tsv`, `lease.tsv` and `copies.tsv` ([`records.md`](records.md#nfr-001--the-records-of-start), task `T-8kqn`), which `internal/run` writes; and the records of a session of `M2b`, `sessions.tsv`, `harnesses.tsv`, `routing.tsv`, the events and the result of a task ([`records.md`](records.md#nfr-001--the-records-of-a-session), task `T-ywk7`; built by task `T-3py1`); and the records of Intake of `M2c`, `intake/meaning.tsv`, `spec/spans.tsv`, `spec/draft.tsv`, `intake/coverage.tsv`, `questions.tsv`, `intake/answers.tsv`, `intake/answers-review.tsv` and `budget.tsv` ([`records.md`](records.md#nfr-001--the-records-of-intake), task `T-g49f`; built by task `T-529i`); and the block `markers` of `M2d`, `intake/markers.tsv` ([`records.md`](records.md#nfr-001--the-records-of-the-scaffold), task `T-msrd`; built by task `T-tetz`), and the source `setup` of `questions` (built by task `T-4c5c`) | `internal/tsv` | no |
| `internal/standin` | for the tests only: a stand-in of the pinned baseline, and a work area set up from it, built at test time (K11; [`setup.md`](setup.md#the-checks-of-layup-setup-verify)); only test files import it | `internal/tsv`, `internal/git`, `internal/work` | no |

**Decided here:** the split of `internal/setup`, `internal/verify` and
`internal/gate` (ADR-0011 decision 1 names "setup, gates, … state files, Git
access" as concerns and gives no names). Reason: `layup setup verify` runs the
gates of a kind on a clean tree and on its known-bad fixture
([`setup.md`](setup.md)), so it imports `internal/gate`; `internal/setup` writes
a tree and must not depend on the checks that judge it, so `internal/cli`
runs the check of each step from `internal/verify` after `internal/setup` did
the step. For the same reason `internal/cli` hands `internal/setup` the lists
that `internal/verify` gives from a tree: the flagged files of the prose step
(task `T-8ya0`), the markers of S10 and the files whose links break at S05
(K10, task `T-8vpw`; [`setup.md`](setup.md#the-checks-of-layup-setup-verify)),
and the evidence call of a step: a function that runs the one-check call of
`internal/verify` on the work area (D12 of #86, task `T-7s0y`). Since task
`T-b3r1` (D9 of #90) these are one value, `setup.Calls`, which `internal/cli`
fills from `internal/verify`: the one-check call, the markers of a tree (each
with the byte column of its open quote, so S11 fills a marker where the scanner
found it, D7), the files whose links break (S05), the files that check
`adapted` flags (S14), and the link rule of check `kit-history` (S05).
`internal/cli` also refuses a brief that holds a marker, by the same scanner. In the same way
`internal/cli` reads the problem statement once per run, runs `internal/psb`,
and hands `internal/setup` the gap table, as `layup psb check` writes it, and
the SHA-256 of the file; `internal/setup` reads the table by its own Go value
of the block `psb-gaps`, which its block test compares with the block (D5 of
#86), so the IDs `Q-NNN` keep one home and `internal/setup` imports no
`internal/psb`.

### The table of M2a

Milestone `M2a` (task `T-zck8`, #123; [`run.md`](run.md), [`forge.md`](forge.md)).
The columns of the table of phase 1, and one more: **"Connects"** is the
character — for a package that may not depend on a package of rule 5, or one
code span per package of rule 5 that it may depend on, by its own imports or
through another package (D8). `cmd/layup` and `internal/cli` depend on them
through the adapter, so they have a row here too, with only their "Connects"
cell; their other cells stay in the table of phase 1. **Decided
here** (O-170 of #127; task `T-esfe`, row 22a of the plan): `TestPackageRules`
reads this table too, and rule 5 from this column, before any package of this
table exists, so the test holds no list of its own (D7 of #79). A cell of this
table that starts with "(the row of phase 1" stands for that package's cell of
the table of phase 1; a row whose cells Job, May import and Starts a program
each start so adds only its cell Connects to that row of phase 1.

| Package | Job | May import | Starts a program | Connects |
| ------- | --- | ---------- | ---------------- | -------- |
| `internal/run` | `layup run`: the steps of Start, the restart, the lease and fencing, the copy of a comment and the rule of a decision ([`run.md`](run.md)); the Go schema of the table `run-steps`; the checks of `git` and of the values of the flags, which `internal/cli` calls before the first step (task `T-mqty`); in `M2b`, the step `probe` and the probe before a task session (task `T-nxe4`), the call that starts an attempt (the event `attempt`) and the start of a task session with its refusals (task `T-d8t9`), the end of a task session: its result, the open attempt, its telemetry row and its comment (task `T-fsjp`), the records of a session, the checks before a push (with the head that `internal/session` reads; task `T-z027`), the push and the bind (task `T-e3sy`) ([`session.md`](session.md)); in `M2c`, the steps of Intake in place of `phase` when Start is done, the first `rule-gaps` (task `T-qk66`), and the session of a phase, its pair, its records commit and its comment (task `T-gmey`), with the budget of milestone 0 asked before it (task `T-ybd5`), and the step `meaning`, the review of meaning with its quote check (task `T-1u2v`), and the steps `numbering` and `spans`, the span table of a session and the span check of the records (task `T-hiv9`), and the step `draft` with its check of the facts' coverage (task `T-2et0`), and the step `completeness` on a harness other than the draft's (task `T-0mey`), and the step `batch`, the merge of the questions and the comment of the batch (task `T-8r3u`), and the step `answers`, the wait, the answer lines and `budget.tsv` (task `T-hxb9`), and the steps `answers-review` and `follow-up` (task `T-vwu3`), and the step `follow-up-answers`, the last of Intake (task `T-tmm2`) ([`intake.md`](intake.md#the-steps-of-intake)); in `M2d`, the step `markers` of Intake, with the setup questions of the batch and its comment (task `T-4c5c`), and the steps of the Scaffold, `setup` (task `T-eahr`: the answers of the setup, the work area, the loop over the stops of the runner and the prose session), `verify` (task `T-79zk`: `layup setup verify` in process, S15 and the one records commit of the setup), `scaffold-push` (task `T-hfn1`: the print of `commands.sh`, the wait and the record of the Operator's applies) and `forge-check` (task `T-wphp`: the eight checks of the forge, the probe among them, and the one records commit of their table; between `lease` and `probe` once the Scaffold is done, task `T-2j13`), with the work area of the setup and the prose sessions ([`scaffold.md`](scaffold.md#the-steps-of-the-scaffold)), and the Go schemas and the readers of the blocks `setup-tree` and `forge-check` ([`records.md`](records.md#nfr-001--the-records-of-the-scaffold), task `T-tetz`); in `M2e`, the task register: the Go schema and the checks of `tasks`, the one records commit that opens the first task with the default `transitions.tsv`, and the check of the permissions of the loop, a `fail` of the step `forge` (task `T-cak3`, [`loop.md`](loop.md#the-task-register)), and the handoff in the loop: the end check of a `comment` row and the kept copy of its file, the post of a comment from that copy with its copy and its lease fence, the record of a handoff with the items that the caller passes, the route of a question to its owner role, and the acceptance of an answer by the next attempt (task `T-q5k4`, [`loop.md`](loop.md#the-artifact-that-is-a-comment); the step `loop` that calls them is row 66k's), and the item computers of the task loop: the parsers of the pinned comment forms and the items `plan-fields`, `verdict-field`, `findings-located` and `nothing-material` (task `T-4pje`, row 66l), the `spec-sections` and `rendering-equal` reads (task `T-yeav`, row 66n) and the composer of the question and answer comments (task `T-gyvo`, row 66o) ([`loop.md`](loop.md#the-items-a-caller-computes)), and the input of the full spec check at a head, `specInput`, each field from its record or the tree (task `T-cav3`, row 66p, [`loop.md`](loop.md#req-012--the-full-spec-check)), and the parameters: the Go schema and owner checks of `parameters`, the table of its six rows, `parameters.tsv` written in the commit of the first task, and the function that a boundary between two steps calls to copy and judge a `set` and write its change, a replacement of `transitions.tsv` included (task `T-dx1h`, [`loop.md`](loop.md#req-012--the-parameter-change)), and the attempts, the independence and the frozen tests: the new attempt after a gate failure with its findings, the authors of a plan review and of a verification and the reviewer's pair, the Go schema of `task-tests` and the freezing of its list in the commit of a handoff to the verifier, the finding `test changed` and the frozen files that an attempt starts with, the tier of a task and the resolution of `by-class` and `owner-of-its-kind`, and the prompt file of a step (task `T-6zj3`, [`loop.md`](loop.md#the-attempts); the step `loop` that calls them is row 66k's), and the four statuses `layup/gates`, `layup/spec`, `layup/rules` and `layup/verify`: their descriptions, their posting once for an unchanged state, the Go schema, the reader and the writer of `task-gates` with the table of a head, and the sources of each status (task `T-q0kr`, [`loop.md`](loop.md#req-007--the-four-layup-statuses); the step `loop` that posts them is row 66k's), and the repairs of a task session at a restart: `RepairTask`, which binds a push with no result or refuses it as `push-lost` and starts one new attempt, and the comment of a task session on the task's issue, read before it is posted (task `T-syxm`, row 66v, [`loop.md`](loop.md#req-005--the-steps-of-the-loop-after-the-scaffold)), and the question in the loop: the post of the question and answer comments, the handoffs to and from the owner role, the answer session `answer-to-a-question` and the next attempt as a position of the driver of a Shape task, and the open material finding of `Position` (task `T-r7b9`, row 66z, [`loop.md`](loop.md#the-checks-of-a-handoff)), and the findings of `layup/spec` and `layup/rules` that start a new attempt (`SpecFindings`, `RulesFindings`, `CheckFindings`, task `T-tnxn`, row 66aa, [`loop.md`](loop.md#the-attempts)), and a task whose base moved: `BringUpToDate`, the clean merge of the default branch into a verified head with the carry of `layup/verify` over it (`mergeEdge`, `CarryInput.Merge`), the redo with its intent and one records commit, `UpToDate` and `TipOf`, `SessionStore.FetchBranch`, and the close of an older attempt's pull request by `RepairTask` (task `T-89xf`, row 66x, [`loop.md`](loop.md#a-task-whose-base-moved-the-carry-and-the-redo)), and the close-out and the merge of a verified task (`FinishTask`, `PostChecks`, task `T-ehb4`, row 66y, [`loop.md`](loop.md#the-merge-at-the-head-sha-and-the-close-out-commit)), and the specification task that a changed class reopens (`openReopenings`, `reopenedTasks`, the record `spec/reopenings.tsv`, the variant of the instruction of the specification and the exemption it feeds to the full check; task `T-4naq`, row 66ab, [`loop.md`](loop.md#the-specification-task-that-a-changed-class-reopens)), and the step loop: the merge order, the scheduler that acts on one task at a time, the permission check and the parameter boundary of `shape`, `bet` and `loop`, the done rule and the restart's last row `loop` (task `T-9nbf`, row 66k, [`loop.md`](loop.md#the-step-loop)) | `internal/tsv`, `internal/git`, `internal/records`, `internal/forge`, `internal/route`, `internal/session`, `internal/ledger`, `internal/rules`, `internal/psb`, `internal/spec`, `internal/setup`, `internal/verify`, `internal/work`, `internal/catalog`, `internal/handoff`, `internal/gate` | no | — |
| `internal/forge` | the forge interface: the six capabilities and their types ([`forge.md`](forge.md)), `Missing` and `CheckPermissions` (task `T-6bq5`), and the types of the calls of `M2d`, `Rule` and `ActivityEntry`, with the head of each branch in `Repository` (task `T-jpkc`); in `M2e` (task `T-cr8j`, [the table of M2e](#the-table-of-m2e)), the nine calls `OpenDraft`, `MarkReady`, `SetStatus`, `Statuses`, `PullRequest`, `Merge`, `EditBody` (task `T-08rj`), `ClosePullRequest` and `CloseIssue` (task `T-xyws`), with their types, `ValidStatus` and `FullSHA`, `CheckLoopPermissions`, and the stand-in forge `forge.Fake`, a non-test file for the tests of other packages; the reader of the forge register (`host:registers/forge.tsv`), with the Go schema of its block; the check of the key file (task `T-1g1q`) | `internal/tsv` | no | — |
| `internal/forge/github` | the GitHub adapter: the JWT, the installation token, the calls of [`forge.md`](forge.md#the-calls-of-m2a), and `Comment` of [`M2b`](forge.md#the-calls-of-m2b) (task `T-fsjp`), and `EffectiveRules` and `Activity` of [`M2d`](forge.md#the-calls-of-m2d) (task `T-jpkc`), and the six calls of [`M2e`](forge.md#the-calls-of-m2e) (task `T-cr8j`), `EditBody` (task `T-08rj`), and `ClosePullRequest` and `CloseIssue` (task `T-xyws`) | `internal/forge` | no | `net`, `net/http`, `crypto/tls` |
| `internal/route` | the reader of the harness register (`host:registers/harnesses.tsv`), with the Go schema of its block; in `M2b`, the columns that `M2b` adds, the readers of `models.tsv` and `routing.tsv`, the checks across the three and of a credential file (task `T-ysph`), the version that needs a probe, admission and the order of the routing register ([`session.md`](session.md#req-013--the-probe-admission-and-routing); task `T-cht1`) | `internal/tsv`, `internal/records` | no | — |
| `cmd/layup` | (the row of phase 1) | (the row of phase 1) | (the row of phase 1) | `net`, `net/http`, `crypto/tls` |
| `internal/cli` | (the row of phase 1, with the change below) | (the row of phase 1, with the change below) | (the row of phase 1) | `net`, `net/http`, `crypto/tls` |

In `M2a` two rows of the table of phase 1 change, in the build task that needs
each: `internal/cli` may also import `internal/run`, `internal/forge`, `internal/forge/github` and
`internal/route` (it reads the forge register through `internal/forge`, builds
the adapter from it, and hands it to `internal/run`); `internal/records` also holds the schemas of the records of
Start ([`records.md`](records.md#nfr-001--the-records-of-start)), and still
imports `internal/tsv` only. `internal/run` commits and pushes the records with
`internal/git`; `internal/records` gives the rows. The later table below keeps
`internal/route` for its jobs of `M2b`.

### The table of M2b

Milestone `M2b` (task `T-ywk7`, #147; [`session.md`](session.md)). The columns
of [the table of M2a](#the-table-of-m2a). **Decided here** (condition 2 of the
plan review of #147): a package that a table already has keeps its one row,
changed in place (`internal/route` and `internal/run` in the table of `M2a`,
`internal/records` in the table of phase 1, and the calls of `internal/git`
[below](#the-calls-of-internalgit)); this table holds only the packages that no
table has. `TestPackageRules` reads this table (task `T-y10b`, row 29 of [the
plan](../plan/README.md#the-tasks-of-m2b)), and the line of the engine checks
below it ([`session.md`](session.md#nfr-005--no-harness-in-the-engine-checks)).

| Package | Job | May import | Starts a program | Connects |
| ------- | --- | ---------- | ---------------- | -------- |
| `internal/session` | a role session: its directory, its environment and credential (task `T-vxdg`), the ID, the context and prompt-size checks, the rule-file check and the head read from the files of `repo/.git` (task `T-6sbe`), the version check, the words of the command, the process, its limits and its stop, and the one call that orders the steps (task `T-5pxd`), its end, the result file, the usage report and the Go schema of `probe-result` (task `T-bpxg`) ([`session.md`](session.md#req-013--a-role-session)); in `M2c`, the `input/` of a session of a phase and the read of its output file by the rules of a result file (task `T-gmey`, [`intake.md`](intake.md#the-sessions-of-intake)); in `M2d`, `ReadOutputFiles`, the read of the output files of a prose session (task `T-eahr`, [`scaffold.md`](scaffold.md#the-prose-rows)); in `M2e`, the record kinds of each step's context in the step table, `ContextKinds` and `StepTable` (task `T-svn3`, [`session.md`](session.md#the-step-table-and-the-kinds-of-a-steps-context)), and `CommentFile`, the read of the file of a `comment` row (task `T-q5k4`, [`loop.md`](loop.md#the-artifact-that-is-a-comment)) | `internal/tsv`, `internal/git`, `internal/records` | `harness` (the command of a harness register row) | — |
| `internal/ledger` | the writer of `telemetry.tsv`: one row per session, from the usage report and `host:prices.tsv` ([`session.md`](session.md#req-011--the-writer-of-the-telemetry-record); task `T-4c3q`); in `M2c`, the rule of the budget of milestone 0, which `internal/run` asks before each session of Intake ([`intake.md`](intake.md#the-budget-of-intake); task `T-ybd5`); in `M2e`, `KnownMoney`, the exact sum of the known money of the rows of `telemetry.tsv`, which that rule and the appetite of the bet's brief read ([`loop.md`](loop.md#the-brief-of-the-first-bet); task `T-h65e`) | `internal/tsv`, `internal/records` | no | — |
| `internal/rules` | in `M2b`, the reader of the rule-path register and the check of a diff before a push ([`session.md`](session.md#a-workflow-or-rule-path-change)), with a second Go value of the block `rule-paths`, whose owner is `internal/setup`, compared with the block by its own test (task `T-m1dx`); the check `layup/rules` and rule batches come in `M2f`, except its basic rule, which is `M2e`'s ([the table of M2e](#the-table-of-m2e)), built by task `T-9et9` (#352) as `Basic`, `Judge` and `Render` of `basic.go` | `internal/tsv`, `internal/records`, `internal/spec`, `internal/git` | no | — |

The packages of the engine checks: `internal/psb`, `internal/verify`, `internal/gate`, `internal/spec`.

**The rule of the engine checks** (NFR-005, task `T-y10b`): no package of that
line depends on `internal/session` or on a package of rule 5, by any path; the
rule reads no cell, so a cell Connects cannot lift it. A later engine check adds
its package to the line. No such line, two such lines, a line with no code span,
or a span that names no package of the module is an error of the test, so the
rule never passes with a package left out.

**A register row.** A cell "Starts a program" that is exactly one code span, one
space and the words "(the command of a harness register row)" makes the row a
register row: the span names the register, not a program. Its package may
import `os/exec` (rule 4), and may start, by a call of `exec.Command` or
`exec.CommandContext`, a program that is not a string literal: the command of a
harness register row. A start by a string literal, a use of `exec.Command` that
is not a call, and each other start are findings there, as elsewhere. A cell that
holds those words and any other text cannot be read. **Known limit:** in a
register row the scan cannot tell `git` from a harness, so rule 3 rests on review
in that package.

### The table of M2c

Milestone `M2c` (task `T-g49f`, #206; [`intake.md`](intake.md)). The columns of
[the table of M2a](#the-table-of-m2a), and the rule of the table of `M2b`: a
package that a table already has keeps its one row, changed in place. **Changed
in place** by this task: `internal/cli` may import `internal/spec` (the command
`layup spec check`); `internal/run` may import `internal/psb` (the rule gaps, in
process) and `internal/spec` (the span check, in process). The other changes of
`M2c` are inside rows that already allow them: `internal/run` writes the records
of Intake, `internal/records` owns their schemas but the block `spec-check`,
which `internal/spec` owns with a second Go value of `spans` (task `T-4ywt`),
`internal/session` builds the
directory of a session of a phase with its `input/`, and `internal/forge` posts
and reads the comments of the Intake issue.

| Package | Job | May import | Starts a program | Connects |
| ------- | --- | ---------- | ---------------- | -------- |
| `internal/spec` | the span check of the facts of a problem statement and its table ([`intake.md`](intake.md#the-span-check)); in `M2e`, the full check, the filter of options and the rendering of a bet ([the table of M2e](#the-table-of-m2e)) | `internal/tsv`, `internal/records`, `internal/git` | `sh`: the pinned prd-lint script of a target | — |

**Decided here** (decision 11 of #206): `internal/spec` is an engine check
(code, no model call), so the build task that made the package (row 41a, task
`T-4ywt`) added it to the line of the engine checks of the table of `M2b`, and
the heading of this table to `TestPackageRules`, in one change, as a span that
names no package of the module is an error of the test. `TestPackageRules`
reads this table in the form of the table of `M2b`: each row is a package of
its own.

### The table of M2d

Milestone `M2d` (task `T-msrd`, #265; [`scaffold.md`](scaffold.md)). It adds no
package. **Changed in place** by this task: `internal/run` may import
`internal/setup`, `internal/verify`, `internal/work` and `internal/catalog`
(decision 4 of #265: the step runner and `layup setup verify` in process, with
`setup.Calls` filled from `internal/verify` as `internal/cli` fills it), and
its job gains the steps `markers`, `setup`, `verify`, `scaffold-push` and
`forge-check`; `internal/setup` gains the inputs of phase 2 and the two
rulesets of S13; `internal/forge` and `internal/forge/github` gain
`EffectiveRules` and `Activity` ([`forge.md`](forge.md#the-calls-of-m2d));
`internal/records` owns the block `markers`, and `internal/run` the blocks
`setup-tree` and `forge-check`; `internal/git` gains `CommitTree` (task `T-j4sd`)
([above](#the-calls-of-internalgit)). `internal/run` starts no program for the
setup: the cell "Starts a program" of its row does not change, and the
programs that `internal/verify` starts stay in its own row.

### The table of M2e

Milestone `M2e` (task `T-xs4y`, #333; [`loop.md`](loop.md)). The columns of
[the table of M2a](#the-table-of-m2a), and the rule of the table of `M2b`: a
package that a table already has keeps its one row, changed in place; this table
holds only the package that no table has. `TestPackageRules` reads the tables
of `M2a`, `M2b`, `M2c` and `M2e`; the build task that made `internal/handoff` (task
`T-h17s`, row 61) added the heading of this table to it in the same change, as task
`T-4ywt` did for `internal/spec`.

| Package | Job | May import | Starts a program | Connects |
| ------- | --- | ---------- | ---------------- | -------- |
| `internal/handoff` | the transition table and the check of a handoff ([`loop.md`](loop.md#req-005--the-handoff)): the schema, each artifact with its hash (a comment by its copy), each required item; the Go schemas of `transitions`; it calls no forge and no session | `internal/tsv`, `internal/records` | no | — |

**Changed in place** by this task (each row keeps its one row in the table that has it):

- `internal/spec`: the full `layup spec check` ([`loop.md`](loop.md#req-012--the-full-spec-check); its exemption of a re-classed fact is `FullInput.Reopening`, task `T-4naq`),
  the filter of options, the rendering of a bet (the rendered commit and the
  inventory), the brief check, and the Go schemas of `inventory`, `inventory-index`,
  `options` and `task-spec-check`; it also holds a second Go value of `tasks`,
  to read the register for the full check — the `spans` precedent — which its
  task compares with the block (`internal/run` owns the writer's). It may now
  import `internal/tsv`,
  `internal/records` and `internal/git`, and it **starts a program**: `sh`, for the
  pinned `prd-lint` of the target (the check that accepts a head only when each
  failure line names an empty MoSCoW or Phase cell), by `os/exec` with the
  target's own script, as `internal/gate` starts `sh -c`. It stays an engine check
  (no model call, no connection), so its row stays on the line of the engine
  checks.
- `internal/rules` (its row in [the table of M2b](#the-table-of-m2b) is changed in place by this task: the job names the basic rule, and the May import cell gains `internal/spec` and `internal/git`): the check `layup/rules` of the basic rule
  ([`loop.md`](loop.md#layuprules)): the base's register against the changed paths
  and the comparison of a rendered record with code's rendering; it may import
  `internal/tsv`, `internal/records`, `internal/spec` and `internal/git`. The batch
  part stays `M2f`'s.
- `internal/run`: the task loop: the register and the transition table, the
  handoff and its events, the attempts, the statuses, the ready and merge rule, the
  close-out check and the merge at the head SHA (task `T-xvvd`, #363, row 66g; it exports `setup.ActionsApp`, the one home of the ID of GitHub Actions, from `internal/setup`), the link of the review record from the body of the pull request (task `T-08rj`, #413, row 66q), the Shape tasks through the task loop to their verified heads (task `T-jy6f`, #415, row 66r), the step `shape` (task `T-3f94`, #364, row 66h: the questions left open, the options session, the filter, the rewrite and the stop, and the two Shape tasks in order), the brief of the first bet (task `T-h65e`, #418, row 66u: `betBrief`, the inputs, the session `bet-brief`, the checks, the appetite after the telemetry, the promotion and the basis `bets/<milestone>/basis.tsv`, [`loop.md`](loop.md#the-brief-of-the-first-bet)), the step bet (task `T-2vva`, #366, row 66j: the post of the brief and the answer of the approver, the inventory version, the rendered commit with its code push and carry-over, the architecture tasks of the settled decisions and `bets.tsv` last, [`loop.md`](loop.md#the-step-bet)), the specification task that a changed class reopens (task `T-4naq`, #428, row 66ab: `spec/reopenings.tsv` and `openReopenings`, [`loop.md`](loop.md#the-specification-task-that-a-changed-class-reopens)), the parameter
  change and the Retrospective step, the instructions of the task steps and their
  inputs ([`loop.md`](loop.md#the-instruction-of-a-task-step)), and the Go schemas of `tasks`,
  `task-gates`, `task-tests`, `bets`, `parameters` and `retrospectives` (`internal/handoff` owns that of
  `transitions`). Its row in [the table of M2a](#the-table-of-m2a) is changed in place:
  its May import cell gains `internal/handoff` and `internal/gate`. It starts
  no program: the cell "Starts a program" does not change.
- `internal/forge` and `internal/forge/github`: `OpenDraft`, `MarkReady`,
  `SetStatus`, `Statuses`, `PullRequest` and `Merge`, `EditBody` (task `T-08rj`), `ClosePullRequest` and `CloseIssue` (task `T-xyws`) ([`forge.md`](forge.md#the-calls-of-m2e)),
  their types, the permission check for pull requests, statuses and checks `read`
  (the last amends ADR-0014 decision 3, #336), and the stand-in forge `forge.Fake`, a
  non-test file of `internal/forge` (standard library only; tests construct it); the
  adapter alone connects, as before.
- `internal/session`: the record kinds of each step's context in the step table
  ([`loop.md`](loop.md#the-tier-and-the-context-of-a-task)), and `CommentFile`, the
  reader of the file that a `comment` row names, by the rules of a result file
  ([`loop.md`](loop.md#the-artifact-that-is-a-comment)).
- `internal/setup` and `internal/catalog`: the `CODEOWNERS` file and the code-owner
  rule of the default branch's ruleset, and the catalog workflow whose gate jobs
  read the base ([`loop.md`](loop.md#the-gate-jobs-read-the-base)); that workflow is
  the catalog's file, the ruleset the setup's. Task `T-jew3` (row 65b): `internal/setup` writes
  `.github/CODEOWNERS`, the rule and, for a target set up before it, the files of its migration
  ([`scaffold.md`](scaffold.md#the-code-owner-rule)); `internal/run` gives it the Operator's login (the row of
  `approvers.tsv`) and the mode (from `setup/record.tsv` and `setup/rulesets.txt`) in `setup.Phase2`. No new
  package, import, command, flag, step or record kind. Task `T-4cep` (#384, the child of row 65b): `internal/setup` owns
  `CheckCodeOwnerLanding`, the judgment of the code-owner commit that [the tree check](scaffold.md#the-tree-check) admits, and
  the landing script it agrees with; `internal/run` asks it of a head that is not the last commit of `setup/tree.tsv`, with the
  Operator's login and the record of the setup, in the step `forge-check`; `internal/git` gives it `Parents` and `BlobID`, with
  `LsTree` for the whole tree. No new package, import, command, flag, step, record kind or schema.
- `internal/records`: the kinds `question`, `decision`, `lesson`, `comment` and `answer` of the block
  `result` (the `answer`, task `T-q5k4`, which settles #377: [`records.md`](records.md#nfr-001--the-records-of-a-session)), and the kinds `handoff` and `merged` of the block `events` (`closed` and
  `rebased` exist there; `M2e` sets their detail).

**No new command** and no flag: the handoff check is a function that `layup run`
calls, the full check is a function of `internal/spec` that `layup run` calls in
process, `layup spec check` keeps its flag `--facts` ([`loop.md`](loop.md#req-012--the-full-spec-check)),
and `internal/cli` does not change.

`internal/handoff` holds the form and the table. Every other item is computed by
the caller (`internal/run`): the comment items `plan-fields`, `verdict-field`,
`nothing-material` and `findings-located` from the kept copy of a comment,
`spec-sections` and `rendering-equal` with `internal/spec` (`HasSection`) and
`internal/rules` (`RenderingEqual`, which shares `RecordRendered` with `Judge`),
and `tests-red-green` from typed per-test results whose producing run is `M2g`'s.
The comment computers are row 66l's (task `T-4pje`), the tree reads row 66n's (task
`T-yeav`) and the composers row 66o's (task `T-gyvo`); the wire that calls them is
the step loop's (row 66k and its size children 66v to 66ab). The caller passes their results in the value it hands to the
check; the check verifies that each required item has a result and that it is
true.

### The test of the package rules

**Decided here** (D7 and D9 of #79):

- The test is `TestPackageRules` in `cmd/layup/rules_integration_test.go`, an
  integration test: it starts `go` and reads files. The checker and its unit
  tests are untagged test files of package `main`: the binary holds no checker
  code, and the hook runs the unit part.
- It reads the table of phase 1, [the table of M2a](#the-table-of-m2a) and the
  line of rule 5 (task `T-esfe`), and [the table of M2b](#the-table-of-m2b) and
  the line of the engine checks (task `T-y10b`); the tables of `M2c` and of `M2e`
  (task `T-h17s`) have the form of the table of `M2b`. A row of the table of `M2b` is a
  package of its own: a row whose cell starts with "(the row of phase 1", or
  whose package another table has, is an error. It fails on a cell that it cannot read, on a
  missing heading or column, on a table with no row, on a row of the table of
  `M2a` that mixes cells of phase 1 and cells of its own, on such a row of
  phase 1 that the table of phase 1 lacks, and on a package in two rows. A
  package that a table names and that does not exist yet is not an error; a
  package that exists and has no row is.
- Rule 5 holds two checks: an own import of a package of rule 5 in a package
  other than the one of its line, and a dependency on one that the cell
  Connects of the package does not name.
- Rules 1, 2, 4, 5 and "May import" come from `go mod edit -json` and
  `go list -deps -json ./...` at the module root, and from the imports of each
  non-test Go file. `go list` gives no import of a file behind a build
  constraint, so the test lists each import that only such a file has, with
  `go list -e -deps -json`.
- Rule 3 and "Starts a program" come from a `go/ast` scan of each non-test Go
  file. A program starts with `exec.Command` or `exec.CommandContext` and a
  string literal that is the program of the row, or, in a register row, a
  program that is not a string literal ([the table of M2b](#the-table-of-m2b)).
  Each other start that the scan finds is a defect: `os.StartProcess`, `syscall.Exec`,
  `syscall.ForkExec`, `syscall.StartProcess`, and an `exec.Cmd` that the code
  makes itself, whose program the scan cannot read. The scan does not read
  cgo code or a raw system call.
- The same checker must find the breaches of rule 5 and of the rule of the
  engine checks, and no other, in the fixture module
  `cmd/layup/testdata/netimport`: an import of `net/http` in `cmd/layup`,
  outside the one package of rule 5, and an import of `net/smtp` in a file
  behind a build constraint, whose dependency on `net` no cell Connects of
  `internal/psb` names, and which breaks the rule of the engine checks in
  `internal/psb` too. In the fixture module `cmd/layup/testdata/enginedep`, it
  must find that `internal/gate`, a package of the engine checks, depends on
  `internal/session` (and imports it, which its row does not allow), and no other
  breach: the start of `internal/session`, by a variable, is its register row's.

Reason: `cmd/layup` is the entry of the module, and its tests already start
programs; a new package for the test needs a row of its own, and the root
`tests/` holds end-to-end fixtures only.

### The calls of `internal/git`

The calls of phase 1 as task `T-2tc2` (#79) leaves them; a later task adds
each call that it needs here first. **Decided here** (D1 of #79, K7): the calls
that the steps, the checks and `layup gate` name.

| Call | The command, after the `-c` values below | Used by |
| ---- | ---------------------------------------- | ------- |
| `Version` | `git --version` | the minimum version (below) |
| `LsRemote` | `git ls-remote --exit-code -- URL REF` | S02; `layup run`, step 3 |
| `Clone` | `git clone --no-checkout -- URL DIR`; with a token, as `Fetch` and `Push` (task `T-ax3r`) | S02; `layup run`, step 3 and the restart (a private target needs the token) |
| `CheckoutDetach` | `git checkout --detach COMMIT` | S02 |
| `Init` | `git init -b main -- DIR` | S03; `layup run`, step 5 |
| `Add` | `git add --all -- PATH…`; no path is the whole tree | S03 to S15; a fixture run of `gate:<kind>` |
| `Commit` | `git commit -m MESSAGE` | S03 to S15; a fixture run |
| `SwitchCreate` | `git switch -c BRANCH COMMIT` | S04: the branch `layup-setup`; `M2b`: the branch of a session's clone ([`session.md`](session.md#the-session-directory)) |
| `SwitchOrphan` | `git switch --orphan BRANCH` | S15: the branch `layup-records` |
| `Branch` | `git symbolic-ref --quiet HEAD` | the step runner: a commit of S04 to S14 only on `layup-setup` (task `T-79y7`) |
| `RevParse` | `git rev-parse --verify --end-of-options REV` | S02, S03: the tree of a commit; S04: the root commit and the branch `layup-setup`; the step runner: the head of `layup-setup` (task `T-7s0y`); `layup gate`: `--base`, `--head`; S15: the heads of `layup-setup` and `layup-records`; a fixture run: its commit (task `T-d6q5`) |
| `ResetSoft` | `git reset --soft COMMIT` | the step runner: the undo of a step whose evidence fails (task `T-7s0y`, D12 of #86) |
| `ResetHard` | `git reset --hard --quiet HEAD` | the step runner: before each step from S04 to S14, the target back to its head (task `T-b3r1`, D11 of #90) |
| `Staged` | `git diff --cached --quiet --exit-code`; exit 1 is a staged change | the step runner: a commit only of a staged change (task `T-b3r1`, D12 of #90) |
| `RootCommits` | `git rev-list --max-parents=0 --end-of-options REV --` | check `pin`; S03: the root commit of a run that stopped (task `T-7s0y`); S15: the records commit of a run that stopped (task `T-d6q5`) |
| `Message` | `git log -1 --format=%B --end-of-options REV --` | S03: the message of the root commit of a run that stopped (task `T-7s0y`); S15: the same for the records commit (task `T-d6q5`) |
| `CommitTime` | `git log -1 --format=%ct --end-of-options REV --`, read as seconds since the epoch | `layup run`: the time of Start, the committer time of the first commit of the records branch, for the budget of milestone 0 (task `T-ybd5`) |
| `IsShallow` | `git rev-parse --is-shallow-repository` | check `pin` (task `T-6x75`) |
| `LsFiles` | `git ls-files -z` | S10; checks `markers` and `adapted` |
| `WorktreeAdd` | `git worktree add --detach -- PATH REV` | `layup gate`, step 2 of the run; `layup setup verify`: the scratch tree; a fixture run; S15: the scratch tree of the records commit (task `T-d6q5`) |
| `WorktreeRemove` | `git worktree remove --force -- PATH` | `layup gate`, step 4 of the run; `layup setup verify`; a fixture run; S15 |
| `Show` | `git show --end-of-options REV:PATH --` | `layup gate`, steps 1 and 2 of the run; `layup setup verify`: the manifest at the head of `layup-setup`; S04: the two index files and `docs/setup/facts.sha256` of the root commit (task `T-7s0y`); S05, S06 and S11: the task indexes, the facts index, `facts.sha256` and `open-gaps.tsv` of the head (task `T-b3r1`); S12: the paths of the entry and `open-gaps.tsv` of the head; S13 and S15: the manifest of the head; S15: the files of the records commit of a run that stopped (task `T-d6q5`) |
| `DiffNames` | `git diff --name-only --no-renames -z --end-of-options BASE HEAD --` | `layup gate`: a `pending` kind; `M2b`: the check before a push ([`session.md`](session.md#a-workflow-or-rule-path-change)) |
| `Apply` | `git apply -- PATCH` | check `gate:<kind>`: the known-bad fixture |
| `LsTree` | `git ls-tree -r -z --full-tree --end-of-options REV -- PATH` | `layup gate`, step 2 of the run: the files of a `config` path at the base, with their modes (task `T-5sgt`); S04: the records of `docs/adr/` and `docs/facts/` at the root commit (task `T-7s0y`); S05: the history at the head; S11: the records of `docs/facts/` at the head (task `T-b3r1`); S12: the tree of the head; S15: the `.sh` files of `layup-setup`, and the tree of the records commit of a run that stopped (task `T-d6q5`) |
| `Parents` | `git rev-list --parents -n 1 --end-of-options REV --`; the fields after the first of the line | `M2e`, task `T-4cep` (#384): the parents of the head in the tree check of the code-owner commit ([`scaffold.md`](scaffold.md#the-tree-check)) |
| `BlobID` | none: the SHA-1 of `blob`, the length, a NUL and the bytes, as `git hash-object` gives, with no process and no object written | `M2e`, task `T-4cep`: the object of the bytes that S13 writes for `.github/CODEOWNERS`, compared with the object of the head's tree |

**The calls of `M2a`** (task `T-zck8`; [`run.md`](run.md)). Each takes the
token of [`forge.md`](forge.md#the-app-identity) in its environment when the
remote needs it; no other value enters the fixed list.

| Call | The command, after the `-c` values below | Used by |
| ---- | ---------------------------------------- | ------- |
| `Fetch` | `git fetch --no-tags --update-head-ok -- URL REF:REF` | `layup run`: the read-back of the root commit (step 5 of [`run.md`](run.md#the-steps-of-layup-run---new)) |
| `Push` | `git push --porcelain -- URL COMMIT:refs/heads/BRANCH`; never `--force` | `layup run`: each records commit (fencing: a push that is not a fast-forward is refused); `M2b`: the push of a session's SHA ([`session.md`](session.md#the-push-and-the-bind)) |

- `Fetch` and `Push` refuse their input before `git` starts unless `REF`
  starts with `refs/` and holds no `:`, `COMMIT` is a full object ID and
  `BRANCH` is not empty and holds no `:`, so no `+` forces an update, and,
  with a token, unless `web` is not empty and has no final slash, as the key
  `http.<web>/.extraHeader` of [`forge.md`](forge.md#the-app-identity) needs
  (task `T-xhgz`). `--update-head-ok` lets `Fetch` set the branch of a repository that
  `Init` made: without it, `git` 2.54.0 refuses to fetch into the branch that is
  checked out (exit 128); the result is the state that `Clone --no-checkout`
  gives. A push that `git` refuses (not a fast-forward, or refused by the remote)
  is a `FailedError` with `Code` 1; its line `!` is on the standard output of
  `--porcelain`, which the error does not keep. A remote that cannot be reached
  gives another code (128, measured with `git` 2.54.0). `Code` 1 is also a
  commit that is not in the local repository, which a records commit that the
  run has just made never is.

**The calls of `M2b`** (task `T-ywk7`; [`session.md`](session.md)). A session's
clone and `layup run`'s clone are both local directories, so no call of `M2b`
takes a token.

| Call | The command, after the `-c` values below | Used by |
| ---- | ---------------------------------------- | ------- |
| `CloneLocal` | `git clone --no-local --no-checkout --single-branch --no-tags --branch BRANCH -- SRC DIR`, then, in `DIR`, `git remote remove origin` | the clone of a session, `repo/` |
| `FetchSession` | `git init --bare -- TMP`; the object directory of the session's `.git` as the one line of `TMP/objects/info/alternates`; in `TMP`, `git cat-file -t SHA`, which must print `commit`, else the result is refused (`branch`); in `TMP`, `git update-ref refs/heads/session SHA`; `git -c core.hooksPath=EMPTY fetch --no-tags --no-write-fetch-head -- TMP refs/heads/session:DST`, where `EMPTY` is an empty directory; `TMP` removed | the fetch of a session's head into `layup run`'s clone, with no command in the session's clone |
| `IsAncestor` | `git merge-base --is-ancestor BASE HEAD`; exit 1 is "no", not an error | the check that a session's head descends from its base |
| `DiffFile` | `git diff -U0 --no-color --no-renames --end-of-options BASE HEAD -- PATH` | the exception of `docs/guardrails.md` before a push |
| `DiffBinary` | `git diff --binary --no-renames --end-of-options BASE HEAD --` | the payload of a refused diff |

- `CloneLocal` copies the objects and shares no ref, configuration or hook with
  `SRC` (`--no-local`, §4); `IsAncestor` takes full object IDs only, as
  `CheckoutDetach` does, and refuses any other text before `git` starts.
- `FetchSession` takes the SHA that `internal/session` read from the files of
  the session's `.git`, a full object ID; `git upload-pack` runs in `TMP`, whose
  configuration is `layup`'s, so no configuration of the session's clone is read
  ([`session.md`](session.md#the-fetch-by-sha)).
- A SHA that names no commit, whether `cat-file -t` prints another type or exits
  non-zero (an absent object), is `ErrNotACommit` in the `Err` of a
  `*FailedError` (task `T-z5dj`); `TMP` and `EMPTY` are removed on each return,
  a refusal included. `FetchSession` refuses, before `git` starts, a `SHA` that is
  not a full object ID, as `Push` does, and a `DST` that does not start with
  `refs/` or holds a `:`, as `Fetch` does. `CloneLocal` refuses an
  empty `BRANCH` or one that starts with `-` before `git` starts. **Known
  limit:** `FetchSession` has no `--update-head-ok`, so a `DST` that is the
  branch of `HEAD` of `dir` is refused by `git`; `layup run`'s clone is on
  `main` and `DST` is a task branch.

**The calls of `M2d`** (task `T-msrd`; [`scaffold.md`](scaffold.md)).

| Call | The command, after the `-c` values below | Used by |
| ---- | ---------------------------------------- | ------- |
| `CommitTree` | `git commit-tree TREE -p PARENT -m MESSAGE`, with the author and the committer in the fixed list of the environment, as `Commit` takes them; it gives the ID that `git` prints. An empty commit is the call with the tree of its parent, which the caller gives | the probe of the forge check: an empty commit on top of the default branch's head, on no ref ([`scaffold.md`](scaffold.md#the-probe)) |
| `AddRemote` | `git remote add -- NAME URL` | the work area's `target/` gets its `origin`, the forge's URL with no credential ([`scaffold.md`](scaffold.md#the-work-area)) |

- The probe is pushed with `Push` and the token, to `refs/heads/layup-probe`;
  a refusal by the remote is a `FailedError` with `Code` 1 whose `Stderr` holds
  the remote's lines, where the forge check looks for GitHub's line of `GH013` for the probe's ref.
- `CommitTree` refuses, before `git` starts, a `TREE` or a `PARENT` that is not
  a full object ID (a `FailedError` with `Code` -1); an output of `git` that is no
  object ID is a `FailedError` with `Code` 0.
- An empty `MESSAGE` is the argument `-m ''` (task `T-j4sd`): `git` finds no
  message and reads it from its standard input, which `call` leaves at the null
  device, so the commit has an empty message and the call returns.
- `AddRemote` refuses, before `git` starts, a `NAME` that is not a plain remote
  name (a letter or a digit, then letters, digits, dots, underscores and dashes)
  and a `URL` that is empty, starts with a dash or holds a control character.
- The setup under `layup run` uses the calls of phase 1 through
  `internal/setup`, and `layup run` builds the work area's `target/` with
  `CloneLocal` from its own clone, then `ResetHard`, `ResetSoft`, `SwitchCreate`
  and `AddRemote` (task `T-eahr`).

**The calls of `M2e`** (task `T-syxm`; [`loop.md`](loop.md#req-005--the-steps-of-the-loop-after-the-scaffold)).

| Call | The command, after the `-c` values below | Used by |
| ---- | ---------------------------------------- | ------- |
| `RemoteHead` | `git ls-remote --exit-code -- URL refs/heads/BRANCH`; with a token, as `Fetch` and `Push`; it gives the SHA of the exact line `refs/heads/BRANCH` | the read of the branch of a `push` at a restart (`RepairTask`) |
| `MergeBase` | `git merge-base -- A B`; the SHA of the merge base, both full object IDs | the path sets of `BringUpToDate` and of `mergeEdge` (task `T-89xf`) |
| `Merge` | in a scratch work tree at `X`: `git merge --no-ff --no-edit -m MESSAGE --end-of-options D`, with the identity and the time of `Identity` as the author and the committer; then `git ls-files -u -z` after an exit of 1 | the merge of the default branch into a verified head (task `T-89xf`) |
| `UpdateRef` | `git update-ref REF SHA` | the scratch branch `layup/base/<task>/<attempt>` of `layup run`'s clone, at the base of a session that is not on the default branch: a verifier's bound head, or the architecture task's base, the specification's verified head (tasks `T-jy6f`, `T-3f94`) |
| `Refs` | `git for-each-ref --format=%(refname) PREFIX` | the refs `refs/layup/sessions/` of the old clone, which `cloneStep` of `layup run` fetches into the new one, and the check that the kept head of an ended session is in the clone (task `T-ehb4`); a `PREFIX` that does not start with `refs/`, or holds `:` or a space, is refused before `git` starts, and a failed read is an error, never an empty list |

- Exit 2 of `git`, and an output with no line for exactly `refs/heads/BRANCH`
  (`ls-remote` matches a pattern against the tail of a ref name), are "no such
  branch", not an error. The SHA of the exact line must be a full object ID, else a
  `FailedError` with `Code` 0. An empty `BRANCH`, one that holds `:` and one that starts
  with `-` are refused before `git` starts (`Code` -1). Any other failure is a
  `FailedError`. The call keeps the token in the environment, as `Push` does.
- `MergeBase` and `Merge` refuse, before `git` starts, an argument that is not a full object ID (`Code` -1); `Merge` also refuses an empty message. `Merge` gives its SHA and `clean` true for a merge commit; exit 1 of `git merge` with unmerged entries is `clean` false with no commit and the work tree removed; exit 1 with none, and every other failure, including a failed read of the unmerged entries, is a `FailedError`. The commit has the two parents `X` and `D` in that order and one SHA for one input. `SessionStore.FetchBranch` is `Fetch` of `refs/heads/BRANCH` of the target under the store's lock and token, then the SHA of that ref of the clone.
- `UpdateRef` refuses, before `git` starts (`Code` -1), a `REF` that does not start with `refs/heads/layup/` or holds a character outside letters, digits, dot, underscore, dash and slash, a `..`, a `//`, an end in `/` or `.lock`, and a `SHA` that is not a full object ID; it moves a ref of `layup`'s own namespace only, and it gives no ref of the target.

Start makes its clone with `Init` and `Fetch`, and the restart with `Clone`. The
first records commit is an orphan commit in a scratch work tree, as S15 makes it
(`WorktreeAdd`, `SwitchOrphan`, `Add`, `Commit`, `RevParse`, `WorktreeRemove`);
each later one uses the same calls on the last records commit, with no
`SwitchOrphan`.

- `--end-of-options` or `--` comes before each revision, URL and path, so an
  input is never an option (`layup gate` takes revisions from its arguments).
  `checkout` and `switch` may read `--end-of-options` as a revision before
  `git` 2.44 (a reading of git's option parser; not measured, the Operator's
  host has 2.54.0). So `CheckoutDetach`, `SwitchCreate` and `ResetSoft` get
  none: `COMMIT` is a full object ID (40 or 64 hexadecimal characters), and the
  call refuses any other text before `git` starts. `CommitTree` gets none for the
  same reason (task `T-j4sd`): its `TREE` and `PARENT` are full object IDs, so no
  input can be an option on any version of `git` from the minimum, and its
  `MESSAGE` is the argument of `-m`.
- `DiffNames` names a renamed path at both ends, so a renamed product path
  counts as changed; `-z` gives each path unchanged.
- S15 makes the orphan commit in a scratch work tree (task `T-d6q5`, #92), as
  `switch --orphan` removes the tracked files from the work tree.

**No default identity** (decided here, D2 of #79). `Commit` takes a name, an
e-mail address and a time: the author and the committer, and both dates
(`GIT_AUTHOR_DATE`, `GIT_COMMITTER_DATE`). Reason: one input then gives one
commit ID (`NFR-005`), and a CI runner has no `user.name`. The identity of the
setup's commits is a decision of the Operator: O-136 (#85) makes it the LAYUP
App's bot, at the date `pin.time` ([`setup.md`](setup.md#the-steps)).

**A call reads no configuration of the host** (decided here, D3 of #79, with
condition 1 of its plan review and findings 1 and 7 of its verification).
Reason: no file and no variable of the host
may change a branch, a byte or a mode of a tree, a file list, a hook, an author
or a signature. `TestAHostileHostChangesNothing` of `internal/git` seeds each
input that the plan review measured with `git` 2.54.0, and first shows that it
changes a plain `git` run. `Commit` sets its identity after the fixed list, so
the host's `GIT_AUTHOR_NAME` cannot change a commit in that test: the unit test
of the environment, `TestTheEnvironmentIsAFixedList`, proves that it does not
pass.

- Each call starts with `-c core.hooksPath=/dev/null`,
  `-c core.attributesFile=/dev/null`, `-c core.excludesFile=/dev/null`,
  `-c core.autocrlf=false`, `-c core.precomposeUnicode=false`,
  `-c commit.gpgsign=false`, `-c http.emptyAuth=false` and
  `-c maintenance.auto=false`. Without the second and the third, the per-user
  attributes file changed the bytes of a staged file, and the per-user ignore
  file dropped a file. Without the last, a commit runs `git maintenance run
  --auto`, whose tasks go on in the background after the call ends (since
  `git` 2.47); in CI of the pull request #115 (`git` 2.55.0) such a task
  repacked into `.git/objects` of a work area while a test removed it (task
  `T-d6q5`). Which version first repacks after a commit is not shown: `git`
  2.54 made the "geometric" strategy the default of manual maintenance, and a
  hand run on 2.54.0 saw no repack after a commit (review round 2 of #92). So
  no call leaves a process of `git` in a repository of LAYUP.
- On macOS, `git init` writes `core.precomposeunicode = true` into the
  repository; `Add` then stored a decomposed (NFD) file name composed (NFC),
  so a tree with such a path differed on a macOS host only. With `false`,
  "file names are handled fully transparent by Git" (`git-config(1)`), and on
  APFS both forms kept their bytes (measured with `git` 2.54.0). A file system
  that changes a name itself ("the unicode decomposition of filenames done by
  Mac OS", the same manual; not measured) still changes the tree, and the S03
  check `root tree = pin.tree` stops the setup.
- The environment of a call is a fixed list: `PATH` and `TMPDIR` of the host
  when they are set; `LC_ALL=C`, `HOME=/dev/null`, `GIT_CONFIG_NOSYSTEM=1`,
  `GIT_CONFIG_GLOBAL=/dev/null`, `GIT_ATTR_NOSYSTEM=1`, `GIT_TERMINAL_PROMPT=0`,
  `GIT_ALLOW_PROTOCOL=file:git:http:https`; for a commit, the six variables of
  its identity; with a token, `Fetch` and `Push` add `GIT_CONFIG_COUNT`,
  `GIT_CONFIG_KEY_0` and `GIT_CONFIG_VALUE_0` of
  [`forge.md`](forge.md#the-app-identity), with LAYUP's values. No other
  `GIT_*` variable, no `GIT_ASKPASS` and no `SSH_ASKPASS` of the host reaches
  `git`: `GIT_CONFIG_COUNT` set a value,
  `GIT_DIR` sent a commit to another repository, and `GIT_AUTHOR_NAME` changed
  the author. A call has no standard input.
- `GIT_ALLOW_PROTOCOL` allows `file` (the stand-in baselines of the tests),
  `https` (a public baseline), `http` (a baseline on a plain host, and the
  loopback server of the test of the no-prompt rule) and `git` (the
  unauthenticated `git` protocol). It refuses `ssh` and each remote helper.
- The list also leaves out the host's proxy and certificate variables
  (`http_proxy`, `https_proxy`, `no_proxy` and their upper-case forms,
  `SSL_CERT_FILE`, `SSL_CERT_DIR`, `GIT_SSL_CAINFO`), because a proxy address can
  carry a credential. So a host behind a proxy, or with a private certificate
  authority, cannot clone the baseline: part of the known limit
  [L-A7](../architecture.md#15-known-limits).
- The `git(1)` manual names `/dev/null` for `GIT_CONFIG_GLOBAL`, and the
  `git-config(1)` manual names `core.hooksPath=/dev/null`; the plan gave
  `core.hooksPath` an empty value, which no manual names. `GIT_ATTR_NOSYSTEM`
  is not in the manual of 2.54.0. Measured on the LAYUP host: its `git` has a
  system attributes file, which `GIT_CONFIG_NOSYSTEM=1` and
  `core.attributesFile=/dev/null` do not skip, and `GIT_ATTR_NOSYSTEM=1` does.
- No credential of the host reaches `git` (K31). With no configuration, no
  credential helper runs. libcurl, under `git`, reads `.netrc` in `HOME`, or in
  the home of the password database when `HOME` is not set; under
  `HOME=/dev/null` no file can be. `GIT_ALLOW_PROTOCOL` (`git(1)`) refuses
  `ssh` and each remote helper, so neither starts: `ssh` finds the user's home
  through the password database, not through `HOME`, and reads its keys there.
  `http.emptyAuth=false` stops GSS-Negotiate with no user name
  (`git-config(1)`; not measured, the LAYUP host has no Kerberos ticket).
  `TestNoCallUsesACredentialOfTheHost` seeds a `.netrc`, an `ssh` and a remote
  helper; under the host's `HOME` and with no `GIT_ALLOW_PROTOCOL`, the server
  got an `Authorization` header and both programs started
  ([`runs/T-2tc2/red.md`](../../runs/T-2tc2/red.md)).
- So `git` asks no question and uses no credential: a clone that needs one
  fails at once. The baseline's repository is public; a private baseline is
  known limit [L-A7](../architecture.md#15-known-limits).

**Two kinds of error** (decided here, D4 of #79, with condition 2 of its plan
review), as Go types for `errors.As`. `NotFoundError`: `git` is not on the
`PATH`. `FailedError`: each other failure, with the arguments, the exit code
and the standard error; the code is -1 when `git` did not start (also when
the call refused its input) or was stopped, and 0 when `git` exited 0 with an
output that the call cannot read.
The packages that may import `internal/git` (`internal/gate`, `internal/setup`,
`internal/verify`) tell the kinds apart and give `internal/cli` a result;
`internal/cli` maps results, never a `git` error, to exit codes. Reason: a
check that could not run is `not-active`, never `fail` (`NFR-004`).

**The minimum version of `git` is 2.32.0** (decided here, D5 of #79): the
first version with `GIT_CONFIG_GLOBAL` (the release notes of git 2.32.0, in
[`Documentation/RelNotes`](https://github.com/git/git/tree/master/Documentation/RelNotes)).
It also covers `switch` (2.23), `init -b` (2.28) and `--end-of-options`
(2.24; for `rev-parse`, 2.30). `internal/git` gives the
version (`Version`) and its test (`Supported`); the packages that import it
check it. The Operator's LAYUP host has 2.54.0, and the host of the build tasks of `M2b` has 2.47.3 (task `T-z5dj`); no test runs 2.32.0, so the minimum rests
on the release notes.

### The components of phase 1

[`architecture.md`](../architecture.md) §2 has a row per component and no row
for `layup setup`, which is neither an engine check (it writes files) nor
`layup run`. **Decided here:**

| Component | What it is | Where it runs | What it may write |
| --------- | ---------- | ------------- | ----------------- |
| The setup runner | `layup setup`: the steps of a target's setup that need no judgement ([`setup.md`](setup.md)) | the LAYUP host | its work area: the baseline copy and its setup branch, the files it prints commands for; and its table on standard output |

### The packages of later phases

Name, job and phase only. Each later milestone gives a package its row in the
table above, with its import rules, in the same change that adds its
requirement's section.

| Package | Job | Phase |
| ------- | --- | ----- |
| `internal/records` | the schemas and rows of the records ([`records.md`](records.md)); the package has its row in the table of phase 1 since task `T-tmhw` (K23); `internal/run` commits the records ([the table of M2a](#the-table-of-m2a)); in `M2e`, `CheckEvent` admits a `push` and a `bound` with no session, the rendered commit of the first bet that code makes (task `T-2vva`) | 2 |
| `internal/forge` | the forge interface: the six capabilities of `architecture.md` §1; its row is in [the table of M2a](#the-table-of-m2a) | 2 |
| `internal/forge/github` | the GitHub adapter of the forge interface; its row is in [the table of M2a](#the-table-of-m2a) | 2 |
| `internal/run` | `layup run`: Start, the phase loop, the lease, fencing; its row is in [the table of M2a](#the-table-of-m2a) | 2 |
| `internal/session` | a role session: its directory, its start, its result; its row is in [the table of M2b](#the-table-of-m2b) | 2 |
| `internal/handoff` | the transition table and the check of a handoff; its row is in [the table of M2e](#the-table-of-m2e) | 2 |
| `internal/spec` | `layup spec check`; its row is in [the table of M2c](#the-table-of-m2c) (task `T-4ywt`) | 2 |
| `internal/rules` | the rule-path register, the check `layup/rules`, rule batches; its part before a push is in [the table of M2b](#the-table-of-m2b), the basic `layup/rules` in [the table of M2e](#the-table-of-m2e) | 2 |
| `internal/audit` | `layup audit` | 2 |
| `internal/ledger` | the cost ledger (`telemetry.tsv`) and the budget; the writer is in [the table of M2b](#the-table-of-m2b) and the budget of milestone 0 in its row there (task `T-ybd5`); the known money that the appetite of the first bet reads is in its row there (task `T-h65e`); the budget of a later bet comes later | 2 |
| `internal/smartif` | the smart-if client: the only package that connects to a model service | 3 |
| `internal/escalate` | the escalation screen | 3 |
| `internal/stall` | progress, the stall triggers and the procedure | 3 |
| `internal/route` | the harness register, the probe, admission and routing; its row is in [the table of M2a](#the-table-of-m2a), with its jobs of `M2b` | 2 |
| `internal/report` | `layup report`: the measures | 4 |
| `internal/learn` | `layup learn`: the reward and the routing update | 4 |

**Decided here** (task `T-tmhw`, #94, K23): `internal/records` has its row in
the table of phase 1, as the section of `REQ-011` is in phase 1 and this table
gives the package "the schemas of every record kind"; one home serves the
telemetry record and its price table (row 17) and the stall record (row 18),
and `internal/ledger` and `internal/stall` import the schemas from it in their
phases.

**Decided here:** the names and the phases. The phase of a package is the
earliest `PRD-0001` phase whose requirement needs it: `layup run` and the
ledger are first needed for the handoffs of `REQ-005` (phase 2), which start
role sessions.
