# The Scaffold: `layup run` drives the setup, the rulesets and the forge check

The specification of milestone `M2d` of [the plan](../plan/README.md#milestones):
what `layup run` adds after Intake, so that a target that its own Start made is
set up by the phase-1 step runner with the answers of the one batch, its
rulesets are applied by the Operator from files that it writes, and it proves
from outside that the rules hold. Its requirements are `REQ-002` in phase 2,
`REQ-003` at the setup (the forge refuses), `NFR-001`, `NFR-002`, `NFR-003` and
`NFR-006` at the Scaffold; the design is [`architecture.md`](../architecture.md)
§3 (What the forge enforces, and what LAYUP checks) and §5 (the gap check step
4, Scaffold 1 to 7), with ADR-0014 decision 6 and ADR-0017 decision 3. Task
`T-msrd` (#265) wrote it, after a public-solution search
([`runs/T-msrd/search/`](../../runs/T-msrd/search/)), whose result is that no
public solution can be its base; the decisions are comment 6110224778 of #265.
The records are in [`records.md`](records.md#nfr-001--the-records-of-the-scaffold).
The phase-1 setup that it drives is [`setup.md`](setup.md).

**What the search gives this file.** GitHub's rules API gives the body of a
ruleset and the effective rules of a branch, and its documentation gives the
limits that this file takes as rules: the read-back returns only the active
rules, so a ruleset in `evaluate` mode is a missing rule; a required check
may leave out its `integration_id`, so the read-back compares it; the bypass list is read only from the
ruleset itself; and a push to a ref that does not exist is a creation, which
only the rule `creation` refuses (searcher A). No public tool compares a
read-back with the file that a program wrote, or proves the rules by a push
that the forge must refuse (both searchers).

## REQ-002 — The setup questions in the one batch

Requirement `REQ-002`: the setup "stops at each human decision, writes the
answers to Git". In phase 1 the decisions of steps S01 and S10 are rows of
`answers.tsv` that the Operator writes ([`setup.md`](setup.md#the-answers)). In
phase 2 they are questions of the one batch of Intake (`architecture.md` §5,
the gap check step 4), as the Operator's block of the batch
([`intake.md`](intake.md#the-comment-of-the-batch)) was left for them.
**Decided here** (decision 5 of #265).

### The step markers

A step of Intake, `markers`, runs after `completeness` and before `batch`
([The steps of the Scaffold](#the-steps-of-the-scaffold)):

1. **The scan.** Code lists each marker of the pinned tree, the tree of the
   default branch's root commit that Start read back (`pin.tree`), by the
   scanner of S10 with LAYUP's `MK_EXEMPT` at the engine's version
   ([`setup.md`](setup.md#the-steps), S10): each by its key `M-<x8>`, its file,
   the line of its open quote and the marker text, in the order of S10's stop
   table (by file, then line). A marker over more lines is one row. **Decided
   here** (finding 1 of round 3 of #265): the key is that of S10, the file and
   the marker text, so a text that a file holds on more lines is one key and
   one row, with the line of its first occurrence, and one question, as S10
   asks it once (`internal/setup`, the stop table of S10). **Decided here**
   (task `T-4c5c`): code reads the one root commit of the default branch in a
   scratch work tree, and the step is `fail` when its tree is not `pin.tree`
   (after Scaffold 6 the head is the setup head, not the pinned tree).
2. **The proposals.** When the scan lists one marker or more, one session of a
   phase ([`intake.md`](intake.md#the-sessions-of-intake)) of the role
   `systems-architect` (the gap review of `architecture.md` §9) reads
   `start/problem-statement.md` and the list, as `input/markers.tsv` in the
   form of the block `markers` with the columns `value`, `start`, `end` and
   `session` as `—`. It writes `result/markers.tsv` (**decided here**, finding 2
   of round 3 of #265), a table of three columns, `id`, `value` and `quote`, one
   row per marker of the list that it proposes a value for, in the order of the
   list: the marker's key, the proposed value, and the quote of the problem
   statement that states it, byte for byte; a marker with no row has no
   proposal. Code reads it by the rules of a session's output file
   ([`intake.md`](intake.md#the-sessions-of-intake)), and refuses a key that the
   list does not hold, a key twice, and an empty value or quote. Code checks each
   quote by the rule of the review of meaning
   ([`intake.md`](intake.md#the-review-of-meaning)): the quote occurs exactly
   once in the file, byte for byte, and holds no tab; a file with a quote that
   fails is refused, and the step is `fail`, naming the session. The offsets
   that code finds are `start` and `end` of the record; the quote itself is
   not kept, as the record's offsets give it from the file. **Decided here**
   (task `T-4c5c`): the rows of the file may be in any order, as code refuses
   only what is named here; the record is in the scan's order.
3. **The record.** One records commit writes `records:intake/markers.tsv`
   ([`records.md`](records.md#nfr-001--the-records-of-the-scaffold)), one row per
   marker of the scan, in the scan's order, with the proposal or `—`. With no
   marker, the table has no row, and no session starts.

### The answers of the setup

The step `batch` merges, after the questions of the idea owner's block and in
the Operator's block, first the question `S01-stack` (the stack of the target,
from the catalog), then one question per row of `intake/markers.tsv`, in its
order. Each is a row of `questions.tsv` with `source` `setup` and `ref` the
setup's own ID (`S01-stack`, `M-<x8>`); its `kind` is `setup`, its `start` and
`end` are the quote's for a marker with a proposal, else `—`. **Decided here**
(task `T-4c5c`): its `block` is `operator`; its `asker` is `layup run` for
`S01-stack`, which has no source row, and the `session` of its row of
`intake/markers.tsv` for a marker; the text of `S01-stack` is that of S01 in
phase 1 (`work.S01Questions`), and a marker's is `the value of <marker> at
<file>:<line>`, the marker as the key of S10 gives it; and a setup question is
never dropped as a repeat of quote and kind ([`intake.md`](intake.md#the-merge-and-the-ids)). **Decided here:**
the enum of `source` of the block `questions` gains `setup`, which the build
task of the step `markers` adds to the block and its Go schema in one change
(condition 4 of the plan review of #265: the block is built). `S01-stack` is
asked whatever rule G1 of `layup psb check` found: a gap of G1 is a question of
the idea owner, and the stack is the Operator's (F-21, decision 4 of #206).

The comment of the batch writes a marker's question as `Q-NNN: the value of
<marker> at <file>:<line>`, with its proposal, when it has one, on the next
lines: `proposed: <value>` and the quote in backquotes. The rule of the answers
at the end of the comment adds the three forms of a marker's answer: `accept`
(the proposed value), a value, or `gap` (the marker stays as an open gap).

The answer of each setup question is read by the rules of the answers and the
follow-up ([`intake.md`](intake.md#the-answers)): its last answer line, or that
of the follow-up that asks it again. **Decided here**, what `layup run` writes
from it into the setup's `answers.tsv` ([The work area](#the-work-area)):

| The answer | The row of `answers.tsv` |
| ---------- | ------------------------ |
| `accept`, to a marker with a proposal | `answer` the proposed value |
| `accept`, to a question with no proposal | read as no answer (the rows below): there is no value to accept |
| `gap`, to a marker | `answer` `gap`, `question_text` the question of the batch |
| a value | `answer` the value |
| no answer, or `open` after the follow-up, of a marker | `answer` `gap`, `question_text` the question of the batch (Invariant 4: a value with no source keeps its marker) |
| no answer, or `open` after the follow-up, of `S01-stack` | no row; the step `setup` is `fail` ([The stops of the runner](#the-stops-of-the-runner)) |

Each row has `by` `operator` and `source` the URL of the comment that holds
the answer, `<web>/OWNER/NAME/issues/<n>#issuecomment-<id>` of its copy, and
for a `gap` with no answer the URL of the comment of the batch. **Decided here**
(task `T-eahr`): the answer of a setup question is that of the follow-up that asks
it again, when that has one (whatever the state of the question it asks again,
which keeps its own); else its own, unless its state is `open` (the review found it
unclear and nothing cleared it), which is no answer. `accept` and `gap`, in lower
case, are the whole answer; any other text, with the spaces round it, is a value,
byte for byte. `by` is the block of the question's row, and the comment of the batch
is the App's comment of the Intake issue that starts with the batch's header line. An accepted
proposal is an answer of the Operator, not a fact citation: the numbered facts
of the target come only with the first bet (`architecture.md` §7), so the quote
stays the evidence of the proposal in `intake/markers.tsv`.

**Known limit:** a marker of a file that the setup then deletes (S05) or
replaces (the prose step) is asked, and its answer is not used: S10 does not
list it, so `layup run` writes no row of it. It stays in `intake/answers.tsv`
and `questions.tsv` of the records, not in the raw fact of S11. A batch that a run
of this version posted before the step `markers` existed (a development build)
gets the record and no setup question: the batch is "done before".

## REQ-002 — The setup under `layup run`

`layup setup WORK` of phase 1 does not change. **Decided here** (decisions 3
and 4 of #265): `layup run` calls the same step runner of `internal/setup` in
its own process, with the inputs of phase 2 given in Go, not as a file;
`internal/run` fills the runner's calls of `layup setup verify` from
`internal/verify`, as `internal/cli` does for phase 1. Rejected: the binary as a
program, which starts a second process whose version must equal the run's and
gives its result as exit codes and text.

### The steps of the setup in phase 2

One row per step that differs from [the step table](setup.md#the-steps); each
other step is as in phase 1. Each row names this section as its decision.

| Step | In phase 2 |
| ---- | ---------- |
| S01 | Reads only `S01-stack`. The name is the target's `OWNER/NAME`, the visibility `forge.visibility` of `start.tsv`, and the baseline `pin.source`; the record rows `name`, `visibility` and `baseline` have source `computed` and ref `records:start/start.tsv`. It runs no `layup psb check`: Intake did, and the `Q-` rows of `answers.tsv` are the questions of the batch ([The work area](#the-work-area)), which the runner takes as given. The checks of phase 1 apply to these values (the stack has an entry in the catalog, the name is `OWNER/NAME`, the visibility is `public` or `private`, the baseline is a URL), and a value that fails is a `fail`, not an input error. `brief.sha256` is `psb.sha256` of `start.tsv`, with the ref of phase 1, `sha256 inputs/briefs/problem-statement.md`. Its evidence line is `S01-stack answered; the stack has a catalog entry`. |
| S02 | Resolves nothing: the record rows `pin.source`, `pin.commit`, `pin.tree` and `pin.time` are those of `start.tsv`, with source `computed` and ref `records:start/start.tsv` (NFR-006: no run resolves the commit again). Its evidence line is `the commit and the tree of start.tsv`. |
| S03 | Makes no commit and writes no command: `WORK/target` is a clone of the default branch, whose one commit is the root commit that Start read back; the step checks that its tree equals `pin.tree`. |
| S04 | As in phase 1, from the pin rows of S02. Its raw fact of the answers holds the `S01-` and `Q-` rows of `answers.tsv` (O-124), each `Q-` row with the `question` of its row of `questions.tsv` as its question text and `S01-stack` with the text of S01 (a `Q-` row that the batch does not hold is a `fail`), so the answers of the batch become the file of the raw fact under the target's `docs/facts/` (`architecture.md` §5, Scaffold 3). Its evidence is that of phase 1, the checks `pin` and `facts`. |
| S12 | As in phase 1, but the record row `module` has the source and the ref of the row `name` of S01 (`computed`, `records:start/start.tsv`): the answer `S01-name` is not a row of `answers.tsv`, and check `sources` of `layup setup verify` refuses a row that cites it (found by the real runner, task `T-79zk`). |
| S13 | Writes the two rulesets of [REQ-003](#req-003--the-rulesets-of-scaffold-6) and their commands: `WORK/out/ruleset-default.json` with the four `layup/` checks, `WORK/out/ruleset-records.json`, and the record row `without-layup` ([NFR-002](#nfr-002--a-target-without-layup)). `layup run` gives the `app_id` of the forge register in Go, and a value below 1 (the register can hold 0) is a `fail`, not an input error. Its evidence line is `the ruleset files, and their commands in commands.sh`. |
| S15 | Writes the record's last rows and the rule-path register, and makes no commit of `layup-records`: `layup run` commits the records ([The step verify and the records of the setup](#the-step-verify-and-the-records-of-the-setup)). |

The actor of S02, S03 and S15 in the step table `setup-steps` is `layup-run`,
a value that the block has. Rejected: new rows of LAYUP's `steps.tsv`, which is
LAYUP's own record and does not change for a target.

### The work area

The work area of a target under `layup run` is `DIR/work/OWNER/NAME/`, in the
layout of [`setup.md`](setup.md#the-command-layup-setup). **Decided here:** the
steps `setup` and `verify` remove it and build it again from the records at
their start, so nothing in it is the only copy of a value (`records.md`, NFR-001
item 4):

- `inputs/briefs/problem-statement.md` and `inputs/briefs/vision.md` from
  `records:start/`, byte for byte;
- `inputs/answers.tsv`: the row of `S01-stack` of [The answers of the setup](#the-answers-of-the-setup)
  (no `M-` row: `layup run` adds those at the stop of S10, [The stops of the runner](#the-stops-of-the-runner)),
  then one row per row of `intake/answers.tsv` whose question is not a setup
  question (**decided here**: a setup question is a row of `questions.tsv` with
  `source` `setup`, or a follow-up whose `ref` is the `Q-NNN` of one; the
  answer of such a follow-up is the setup row's, and it is not a `Q-` row), with `question` its `Q-NNN`, `by` the role of its block, `source`
  the URL of its comment and `question_text` `—`; in the order of the question
  IDs;
- `inputs/files/<path>` from `records:setup/inputs/files/<path>`, the texts of
  [the prose rows](#the-prose-rows);
- `target/`: a clone of the default branch from the run's clone, with `main`
  at the root commit: the row of order 0 of `setup/tree.tsv` once it exists,
  else the default branch's one commit, which Start read back (**decided here**,
  finding 2 of round 2 of #265: after the Operator's push the default branch
  holds the setup commits, and S03 in phase 2 needs the root commit).

**Decided here** (task `T-eahr`): `target/` is made by `git clone` of the run's
clone (one branch, with its objects, no remote), then its work tree is put at the
branch, `main` is put at the root commit (the branch of the clone is `main` whatever
the forge's default branch is called, as the runner reads `refs/heads/main`), and
`origin` is the URL of the forge's repository with no credential: S03 in phase 2
writes no command, and the commands of S13 read and push through `origin` with the
Operator's own login. The row of order 0 of `setup/tree.tsv` must be in the clone,
and with no such file the default branch must hold the one commit that Start read
back, else the step is `fail`. A path of `setup/inputs/files/` that leaves
`inputs/files/` is refused before the work area is touched, and so is an entry of
that directory that is not a regular file. A link at `DIR/work` or at
`DIR/work/OWNER` is refused too, as the removal of the area would reach outside `DIR`.

The setup is deterministic (O-136: the author, the committer and the date of
each setup commit are fixed), so a rebuilt work area gives the same commits.
The outputs of the Operator's apply commands are not in the work area: they are
in `DIR/applies/OWNER/NAME/` ([The commands for the Operator](#the-commands-for-the-operator)).

### The stops of the runner

A run of the step runner that stops (exit 3 of phase 1) gives its stop table
([`setup.md`](setup.md#the-stop-table)). **Decided here:** `layup run` answers a
stop from the records, and runs the runner again, until a run does not stop:

| Rows of the stop table | What `layup run` does |
| ---------------------- | --------------------- |
| `M-<x8>` rows (S10) whose marker the batch asked | `layup run` adds one row per such marker to `inputs/answers.tsv` from [the answers of the setup](#the-answers-of-the-setup) (an answer, or `gap`), and runs the runner again |
| `F-<path>` rows (S05, and the prose step) | one session of [the prose rows](#the-prose-rows) for that stop |
| the one row `O-verify` of S15 | the runner has done S01 to S14: the step `setup` is `done`, and the step `verify` runs it again once `out/verify.tsv` exists |
| any other row (`S01-stack` with no answer, an `M-` row that the batch did not ask, any other `O-` row, `O-verify` of another step, `O-verify` with another row) | the step `setup` is `fail`; its `detail` names the rows |

**Decided here** (finding 1 of round 2 of #265): `inputs/answers.tsv` holds an
`M-` row only for a marker that S10 lists, as the runner refuses an `M-` row
whose marker is not in the tree at S10 (an input error, `setup.md`, The answers;
`internal/setup`, `CheckAsked`). So `layup run` writes no `M-` row before S10,
and S10 stops once with each marker of the tree; a marker that the batch asked
then gets its row, and the next run of S10 is done. A step that ends `fail` or `not-active`
makes the step `setup` `fail`, with the runner's reason; an input error (exit 2
of phase 1) does the same.

**Decided here** (task `T-eahr`, #277, with the conditions of its plan review): the
runner always reaches S15 after S14 and S15 stops on `O-verify` while `out/verify.tsv`
is missing, so that stop, alone and of S15, is the end of the step `setup`
(`M-` rows are answered one run at a time, so a table that holds them never holds
it). The same stop twice in a row, the rows equal in order in `Step`, `Question`,
`Ask` and `Where`, is `fail`, naming the rows: the loop has no other bound, and it
stops after twelve runs, `fail` too. A marker that the batch asked and that no stop
lists gets no row and fails nothing. A prose session ends `fail` of the step with
its own text, which names the session (a pair that is not admitted, a refused output,
a class that is not `done`, the cap of milestone 0); nothing of it is committed, but
for a session that ends in an error after its commit (its comment is refused, its
directory is not removed, the copy into the work area fails): its files are in the
records, and it is the last session of the run.
**The files that a `fail` removes** are those of the last prose session of the run,
which the loop keeps across the runs that follow it with `M-` rows only: whichever run
ends in the `fail`, one records commit through fencing removes them, its message naming
the session and the reason, and it refuses a path that is no regular file of the records
base, one below a link of it too. A later stop with `F-` rows means that the runner has passed
the rows of the session before it, whose files stay; a `fail` with no session in the
run removes none.

**Decided here** (finding 3 of round 1 of #265): the runner stops with `F-`
rows only for a file that is missing (`setup.md`, S05 and S07), so a prose
text that is in the records is never asked again. When the step `setup` is
`fail` after a prose session of the same run (a check that fails, an `M-` row
of S10 that the batch did not ask, as the session's text holds a new marker, or
any other `fail` that follows the copy), one records commit removes each file
of that session from `setup/inputs/files/` before the step ends, its message
naming the session and the reason. So the next restart builds the work area
without them, the runner stops on their `F-` rows again, and a new session
writes them. The files of an earlier session whose rows the runner passed
stay.

### The prose rows

Scaffold 2: the texts of S07 (`docs/onboarding-for-engineers.md`), S08
(`docs/glossary.md`), S09 (`docs/guardrails.md`) and S14 (`README.md`,
`AGENTS.md` and each file that check `adapted` flags), and the fixed text of a
file whose link S05's deletion breaks. **Decided here** (decision 6 of #265):

- one session of a phase per stop with `F-` rows, of the role `product-owner`.
  The §9 table has no row for these steps; these texts restate the problem
  statement for the target's readers (the onboarding links it, the glossary
  defines its terms), which is the work that §9 gives the Product Owner (the
  draft, the specification); one session keeps the prose step one group
  (O-123);
- its `input/`: `problem-statement.md`, `vision.md` when there is one, and
  `stop.tsv`, the stop table of the runner (the block `setup-stop`); its
  `repo/` is `WORK/target` at the head of `layup-setup` at that stop;
- its prompt (**decided here**, task `T-eahr`, round 2): no placeholder and no
  angle quote, but a row of an adapted file (its ask is the `adaptAsk` of
  `internal/setup`, which the prompt names by its first word) names a file that
  `repo/` has, which the session gives with its changes and each marker it holds
  kept, in every place, byte for byte, as S14 requires ([`setup.md`](setup.md));
- its output: one file per `F-` row, at `result/files/<path>` of its session
  directory (`result/` is where a session writes, [`session.md`](session.md#the-session-directory)); a missing file,
  a link, a file over 1 MiB, a file that is not valid UTF-8, or a file at a
  path that the stop does not name refuses the output, and the step is `fail`,
  naming the session. **Decided here** (task `T-eahr`): the whole set is read and
  checked before any file is committed: a link is refused (as a file, as a
  directory above a file, and as `result/` or its output directory), an empty
  directory or a nested file that no row names
  refuses the output, 1 MiB (1,048,576 bytes) passes and one byte more does not,
  and a path that is not clean (`..`, an empty part, absolute) is a defect of the
  runner's row, which the step reports as its own `fail`; the session's `repo/` is
  `WORK/target` cloned on its one branch `layup-setup`, on its task branch at the
  head of `layup-setup` that the stop was made at;
- code commits each file to `records:setup/inputs/files/<path>` in one records
  commit (the commit of the end of the session, with its telemetry row), writes
  the same bytes to `inputs/files/<path>` of the work area, which the runner reads
  at its next run, and the next build of the work area copies it to
  `inputs/files/<path>`; the runner then runs each row's check (`kit-history`
  and `link-lint` for S05; `onboarding`, `glossary`, `guardrails`, `identity`
  and `adapted` for the prose step). A check that fails is the runner's `fail`.

So no role session edits a path of the target: code copies its text, and the
check of the row decides. A restart after a `fail` builds the work area again
and starts a new session for the stop ([The stops of the runner](#the-stops-of-the-runner)):
the files of the failed session are removed from the head of the records, and
stay in their history.

### The step verify and the records of the setup

When the runner has done S01 to S14, the step `verify` runs `layup setup
verify` on the work area ([`setup.md`](setup.md#the-checks-of-layup-setup-verify))
and writes its table to `WORK/out/verify.tsv`. A row that is not `pass` or
`clear` makes the step `fail`, naming the check and its reason. Otherwise the
runner does S15, and one records commit writes:

- `setup/record.tsv`, `setup/verify.tsv` and `rule-paths.tsv`, equal to the
  files of `WORK/out/`;
- `setup/tree.tsv`, the record of the verified tree
  ([`records.md`](records.md#nfr-001--the-records-of-the-scaffold)): the root
  commit, then each setup commit of `layup-setup` in order, each with its tree
  and its step, so its last row is the head that the Operator then pushes;
- `README.md` in [the text of the Scaffold](#the-readme-of-the-scaffold).

**Decided here** (task `T-79zk`, #278, with the conditions of its plan review): the
step is "done before" when `setup/record.tsv` is in the records, as `setup` is, and
else removes the work area, builds it again from the records and runs the loop of
[the stops of the runner](#the-stops-of-the-runner) as `setup` does (the runner is
deterministic, so the commits are the same); a loop that does not end done is the
`fail` of the step with its detail. It then runs `layup setup verify` in process,
writes its table to `out/verify.tsv` and runs the loop again, which ends S15. The step
is `fail`, and no setup records commit is published, when `layup setup verify` gives an
error whatever its table holds (a scratch tree that is not removed, with every row
`pass`, too), when `out/verify.tsv` cannot be written, when a row is not `pass` or
`clear` (its detail is `layup setup verify: <check> is <result>: <reason>` for the
first, in the table's order, whatever its result), when the second loop does not end
done, when S15 did not complete (`out/record.tsv` lacks a row `setup.head`,
`verify.sha256` or `rule-paths.sha256` of S15, or holds a hash that is not that of its
file, or `out/rule-paths.tsv` holds no path), or when the tree or the commit fails. The
records commit (message `layup run: the records of the setup`) writes exactly the five
paths (a path whose bytes the base holds already is no change of the commit) and keeps
the other records. `setup/tree.tsv` has the commit of
`main` of `WORK/target` with the step `root`, then the commits of `layup-setup` above
it, oldest first by their first parents (S14 comes before S11), each with the step of
its message `chore: setup <step>`; a message of another form, a `main` that is not an
ancestor of `layup-setup` and no commit above `main` are a `fail` (the reader of the
block still takes a table of the root alone).

The step `verify` is the only caller of `layup setup verify` under `layup
run`, after S14, so F-6 (verify before S12) is not reached in phase 2.

### The README of the Scaffold

The text of `README.md` of `layup-records` that the step `verify` writes in
place of [the text of Start](run.md#the-readme-of-a-target-that-start-makes)
(**decided here**: it names the files of Start and of the setup):

```text scaffold-readme
# The records of this repository

This branch, `layup-records`, holds the records that LAYUP keeps for this
repository. It is an orphan branch: it shares no commit with the default
branch.

Only `layup run` writes this branch from Start on. Its first commit holds the
records of the Start.

Each file is a table of tab-separated values with a header row, or Markdown,
so a person reads it with no tool. A plain `git clone` carries the branch as
`origin/layup-records`.

- `start/start.tsv`: each value of the Start, with its source.
- `start/problem-statement.md`, `start/vision.md`: the two briefs, byte for byte.
- `approvers.tsv`: each account whose comment can decide, by its numeric ID and
  role.
- `lease.tsv`: the run that holds this target.
- `questions.tsv`, `intake/`: the questions of Intake and their answers.
- `setup/record.tsv`: each value of the setup, with its source, and one done
  row per step.
- `setup/verify.tsv`: the table of `layup setup verify` before the records
  commit of the setup.
- `setup/tree.tsv`: the commits of the setup that `layup setup verify` read,
  which the Operator then pushes.
- `setup/inputs/files/`: the texts of the setup that a role session wrote.
- `setup/rulesets.txt`: the Operator's apply of the two rulesets, with their
  bypass lists.
- `setup/forge-check.tsv`: the rules read back and the probes of the last start.
- `rule-paths.tsv`: the rule-path register: the paths whose change is a change
  of the rules.
```

## REQ-003 — The rulesets of Scaffold 6

`REQ-003` at the setup: "an agent cannot change a rule path without a control
that the agents cannot pass by themselves"; ADR-0017 decision 3 gives the
control of the forge. **Decided here** (decision 2 of #265): the Operator
applies both rulesets at Scaffold 6, after the push of the setup commits, as
ADR-0014 decision 6 and ADR-0017 decision 3 say ("The Operator applies both at
setup"); `records.md` (NFR-001, "Not in phase 1": "applied at Start") is
corrected to it. Before Scaffold 7 no start reads back or probes the rules,
as there are none.

S13 writes both files in phase 2 (decision 7 of #265, with S13 as the one
writer of a ruleset, as in phase 1). Each file is the body of `POST
/repos/{owner}/{repo}/rulesets`, in the form of the file of phase 1
(`internal/setup`, `rulesetBody`): JSON with an indent of two spaces, the target
`branch` and the enforcement `active`. The name of the default branch's
ruleset is that of phase 1, `layup: the default branch` (`setup.md`, S13); the
records branch's is `layup: the records branch` (**decided here**, in the form of
the first, so the Operator tells the two apart). Every parameter of
each rule is written, so the read-back compares each.

### The ruleset of the default branch

`WORK/out/ruleset-default.json`, the ruleset of `~DEFAULT_BRANCH` and
`refs/heads/layup-probe`, with an empty `bypass_actors`:

| Rule | Parameters |
| ---- | ---------- |
| `creation` | none. **Decided here** (searcher A): the probe pushes to `layup-probe`, a ref that does not exist, which is a creation; with no bypass actor, no one creates it |
| `deletion` | none |
| `non_fast_forward` | none |
| `pull_request` | those of phase 1 (no required review, thread resolution required), with `allowed_merge_methods` `["merge"]` (F-19: the target's rules forbid a squash and a rebase) and `require_extra_approval_for_unattributed_changes` `false` (F-19: no rule of the target asks for it, and the pilot saw that it did not block a pull request) |
| `required_status_checks` | `strict_required_status_checks_policy` `true`; one check per gate job, its `integration_id` that of GitHub Actions, 15368, as in phase 1 (V-21 of `docs/setup/record-T-n1hp.md`); and `layup/gates`, `layup/spec`, `layup/verify` and `layup/rules`, each with `integration_id` the `app_id` of the forge register (§6: each `layup/` check pinned to the App) |

**Decided here:** `creation` and the two parameters of `pull_request` are
written in phase 1 too, as `WORK/out/ruleset-default.json` is a file of the
work area, not of the target's tree, so no target commit changes;
`docs/setup/branch-protection.json` does not change. The four `layup/` checks
are phase 2 only: no phase-1 command posts them (`records.md`, NFR-002 item 3).
**Decided here** (task `T-9v92`): the rules are written in the order of the table,
so the file reads as the table does; the read-back compares them by type.

### The ruleset of the records branch

`WORK/out/ruleset-records.json`, the ruleset of `refs/heads/layup-records`:
the rules `update`, `deletion` and `non_fast_forward`, in this order, with one
bypass actor, `actor_type` `Integration`, `actor_id` the `app_id` of the forge
register, `bypass_mode` `always`. The rule `update` lets only a bypass actor update
the branch (searcher A), so only the App writes it, and the ruleset of the default
branch keeps its empty bypass list: the two scopes are two rulesets. **Decided
here** (task `T-9v92`, condition 1 of its plan review): the rule `update` has one
parameter in GitHub's description of it, `update_allows_fetch_and_merge`, and it
is written `false` (no fork pulls into the branch), so the read-back compares it;
`deletion` and `non_fast_forward` have none. The setup record has no hash row for
this file: it is a function of the App's ID, which the forge register holds.

### The commands for the Operator

S13 writes to `commands.sh`, in order: the read of the default branch's head on
the forge (`git ls-remote origin refs/heads/<default>`); the push of
`layup-setup` onto the default branch (`git push origin
layup-setup:<default>`, a fast-forward from the root commit) when that head is
the root commit of `WORK/target`, the row of order 0 of `setup/tree.tsv`
(`pin.commit` is the baseline's commit, not the target's), no push when it is
already the head of `layup-setup`, and a stop with exit 1, naming the head, when
it is neither; then the apply of each ruleset with `gh api --method POST
/repos/OWNER/NAME/rulesets --input <file>`, its output written to
`DIR/applies/OWNER/NAME/default.json` and `records.json`, each apply only when
its output file does not exist yet. The output goes first to a temporary file
beside it, which `mv` renames to the output file only after `gh api` exits 0
(**decided here**, finding 1 of round 5 of #265), so an apply that fails, also
one that wrote part of its output, leaves no output file and is done again. **Decided here** (finding 1 of round 4 of
#265): so a run of `commands.sh` that stopped after the push, or after one
apply, is run again as it is, and does each command that is not done yet. **Decided here** (F-13): `commands.sh` starts with `set -eu`,
so the first command that fails stops it, in phase 1 and phase 2; a command of
phase 1 that does not hold under `set -eu` is a defect of its step.

**Decided here** (task `T-56dj`, #275, with the conditions of its plan review): in
phase 2 the text of S13 is four commands, each with a comment line of its own: the
reads, the guard on the head and the push, in order 2, and the directory of the
apply outputs and the apply of each ruleset, in order 4. The paths are absolute, so
the Operator may run the file from any directory; `layup run` gives `DIR` and the
default branch in Go. Each read of `git` is an assignment of its own, so a read that
fails stops the file with its status, as a pipe would hide it; the head is the first
field of the answer of `ls-remote`, cut in a second assignment, and a head that the
forge does not give is empty, which is neither commit: the file stops with exit 1 and
the message, on standard error, names the empty head, the root commit and the head of
`layup-setup`. The root commit and the head of `layup-setup` are read from the work
area when the file runs and are not written into it, so one file serves a work area
that is built again, which gives the same commits ([The work area](#the-work-area)).
The temporary file of an apply that fails stays, and the next run writes over it:
only the output file says that an apply is done. S13 is a `fail` when the host
directory is not an absolute path or the default branch is empty.

When `DIR/work/OWNER/NAME/` is missing (a restart on another host, or after
it was removed), `scaffold-push` builds it again as `setup` does, with the
prose files of the records, and runs the step runner; the head of
`layup-setup` must equal the last commit of `setup/tree.tsv`, else the step is
`fail`, naming both. The step `scaffold-push` prints `commands.sh` on standard error with the
Operator's rule (with the Operator's own login, from a clone with no hooks
installed, `architecture.md` §5 Scaffold 6) and is `done` with "waiting for the
Operator's pushes and ruleset apply". A later restart goes on when the head of
the default branch on the forge is the last commit of `setup/tree.tsv` and both
files of `DIR/applies/OWNER/NAME/` exist; else it is `done` with the same
"waiting". When it goes on, one records commit writes `setup/rulesets.txt`:
each apply output, byte for byte, after a line `== <file>` that names it. Each
output must be a JSON object whose `source` is `OWNER/NAME`; else the step is
`fail`, naming the file. Each apply output holds the ruleset's `bypass_actors`,
the list that the App's permissions cannot read (`architecture.md` §3; known
limit L-A4).

**Decided here** (task `T-hfn1`, #279, with the conditions of its plan review): the
step is "done before" when `setup/rulesets.txt` is in the records. Else a missing
work area is built again, and its head checked, before either outcome below, the
waiting one and the going-on one; an existing work area is used as it is, and its
head is not read. The step then reads the last commit of `setup/tree.tsv`, the head
of the default branch on the forge (a branch that the forge lacks has no head, which
is no commit, so the step waits) and the files `default.json` and `records.json` of
`DIR/applies/OWNER/NAME/`; a name that is a directory or cannot be read is a `fail`,
naming it. While it waits, the step prints a note that says to run `out/commands.sh`
of the work area with the Operator's own login, from a clone with no hooks installed,
then each line of that file as a note, and it ends the restart, which releases the
lease. The file `setup/rulesets.txt` is, for `default.json` and then `records.json`,
the line `== <file>` and a line feed, the output byte for byte, and one line feed that
LAYUP adds as a separator. The marker and the separator are not bytes of the output:
an output is what lies between its marker line and the line feed before the next
marker line, or before the end of the file, whether or not it ends with a line feed.
A loop of the rebuild that does not end done is the `fail` of the step, with its detail.

## REQ-003 — Scaffold 7, the forge check

The step `forge-check` (decision 8 of #265). Each rule below is the exact rule;
its result is one row of `records:setup/forge-check.tsv`
([`records.md`](records.md#nfr-001--the-records-of-the-scaffold)), in the order
of the subsections, and every check runs, so the table shows each result. A row
that is not `pass` makes the step `fail`, its `detail` naming the first such
row; the table is committed in both cases.

### The tree check

`tree`: the head of the default branch on the forge (`Repository`, which gains
the head of each branch in `M2d`, [`forge.md`](forge.md#the-calls-of-m2d)) equals, by SHA, the last commit of `setup/tree.tsv`. The push is a
fast-forward of the verified commits, so their SHA, not only their tree, is
on the forge. Any other head is `fail`, naming the head.

### The read-back

`rules:<branch>` for the default branch, `layup-probe` and `layup-records`:
`EffectiveRules(branch)` ([`forge.md`](forge.md#the-calls-of-m2d)) holds each
rule of the branch's ruleset file, with each parameter equal; the required
checks are compared as a set of pairs of a context and an `integration_id`. A
rule that is missing (also one of a ruleset in `evaluate` or `disabled`, which
the call does not return) or a parameter that differs is `fail`, naming it. A
rule that the call returns and the file does not hold (a ruleset of an
organization) is not a failure; its type is named in the `detail`.

`bypass:<ruleset>` for `default` and `records`: the `bypass_actors` of the
apply output in `setup/rulesets.txt` are empty for `default`, and exactly the
App's (`Integration`, `app_id`, `always`) for `records`. A later change of a
bypass list is not seen (L-A4).

### The probe

`probe`: `layup run` makes an empty commit on top of the default branch's head
(the tree of the head, the head as its parent, the App's bot as author and
committer, at the run's time) in its own clone, and pushes it to
`refs/heads/layup-probe` with the installation token (`internal/git`). It
passes when the push exits non-zero and its standard error holds both `GH013`
(GitHub's "Repository rule violations found") and the name of the default
branch's ruleset, `layup: the default branch` (**decided here**, finding 3 of
round 3 of #265: §3 counts only a refusal "naming the ruleset", so a refusal by
another ruleset, which a bypass of LAYUP's own could hide, is no pass). Else it is `fail`: a push that
exits 0, as "probe accepted: <SHA>" (the Operator then deletes the ref and
checks any workflow run that it started, §3); any other non-zero exit, as
"probe not run: <the first line of standard error>" (an expired token, a
network or server error), never a pass (FT1). The commit is on no branch of the
run's clone after the push (**decided here**, task `T-wphp`: a push that the forge accepts, which is a
`fail`, may leave the remote-tracking ref `origin/layup-probe`, as the `git` of the host of the tests
does and that of CI does not, because the push is to the URL of the clone's own `origin`; the next start
clones afresh).

### The activity read

`activity`: `Activity(ref)` of the default branch ([`forge.md`](forge.md#the-calls-of-m2d)),
one page, with the installation token, returns without an error. An error is
`fail` (`layup audit`, `M2f`, needs the read). The rows that it returns are not
judged here.

**Decided here** (task `T-wphp`, #280, with the conditions of its plan review): the
step runs after `scaffold-push` is `done` and not waiting, and is not "done before"
(its later starts are [below](#at-each-later-start)). Its table always has the eight
rows, and the records commit of `setup/forge-check.tsv` is made in every case; the
step is `done` "the forge check passed", or `fail` naming the first row that is not
`pass` as `<check>: <detail>`.

- **The files** of the read-back are `out/ruleset-default.json`, for the default branch and for
  `layup-probe`, which its conditions name, and `out/ruleset-records.json`, for `layup-records`, in the
  work area, which the step makes sure of as `scaffold-push` does (a work area that is there is used as
  it is, and its head is not read).
- **The comparison:** each rule of the file needs one returned rule of its type whose parameters equal
  each of the file's own; the parameters of two returned rules are never combined, so the answer's order
  does not matter. A parameter that the call returns and the file lacks is ignored; one that the file
  holds and the call lacks is a difference, whatever its value (`false` and `0` too). The required checks
  are the set of the pairs `context` and `integration_id`, in any order, with their strict flag. The
  `detail` of a `pass` is empty (`—`), except a rules row with a rule type that the call returns and the
  file lacks: `extra: <types>`, sorted, a comma and a space.
- **The bypass lists** are read from `setup/rulesets.txt` by its framing (an output is what lies between
  its marker line and the line feed before the next marker line, or before the end of the file): the
  `bypass_actors` of `default.json` is an empty list, and that of `records.json` is one actor,
  `Integration`, `always`, with the `app_id` of the register.
- **What fails what:** a failure of a check's own call (the forge, the push, a file, the token, the
  commit of the probe) is that row's `fail`; a work area that cannot be made sure of fails the three
  rules rows; a `setup/tree.tsv` that is absent or refused fails the `tree` row and those three, and a
  `setup/rulesets.txt` that is absent or refused the two bypass rows; a read of the repository that fails
  fails the `tree` row and the `probe` row ("the head was not read": the expected commit is never put in
  its place). Only records that cannot be read, a table that its reader (`ReadForgeCheck`) refuses (a
  run ID that is none) or a records commit that fails end the step with no table, and its detail says so.
- **The probe's commit** has the message `layup: the probe of the forge check`, and is made in the run's
  own clone, which the step `clone` made afresh, so the objects of the head are there. The `run` and `at`
  columns are the run's ID and the run's clock when the check ends: one reading for each row.
- **A detail** is the first line of its reason, in valid UTF-8, with `DIR` for the host's directory as a
  detail of `run-steps` has it. The `activity` read is of the default branch by its name.

### At each later start

**Decided here:** once the Scaffold is done, each restart runs `forge-check`
after `lease` and before `probe`, and writes the table again in a records
commit (Git history keeps each); a row that is not `pass` stops the run at that
step, as `architecture.md` §3 says ("stops when a rule is missing"), and the
steps after it have no row. The tree check of a later start compares the default
branch with the head of `setup/tree.tsv` only until the first merge after the
setup (`M2e`), which adds its own rule; in `M2d` no merge follows the setup.

**The Scaffold is done** (**decided here**, task `T-2j13`, #281, with the
conditions of its plan review) when `setup/forge-check.tsv` is in the records,
`ReadForgeCheck` accepts it and each of its eight rows is `pass`; the step `clone`
reads it with the other records. The specification fixes this trigger and leaves
a refused table unstated, so: an absent table, a table that its reader refuses
and a table with a row that is not `pass` give the chain of the first start, with
`forge-check` last, after `scaffold-push`, which writes the table again, and a
refused table is replaced, not a failure of its own. With Start not done the chain
is that of Start, with `phase`, and has no `forge-check`, whatever the table.

At a later start the steps of Intake and of the Scaffold are each "done before",
or give the detail that their own rule gives for a finished Intake (an R10
correction of "are then done before", which the rules of
[`intake.md`](intake.md#the-steps-of-intake) contradict): `follow-up` says
"nothing to ask again; Intake is done" and `follow-up-answers` "no follow-up;
Intake is done". The progress of a restart with Start done counts 21 steps in
both chains, after the `1/6 forge` and `2/6 clone` of a restart with Start not
done: `5/21 forge-check` at a later start and `21/21 forge-check` at the first
start.

## NFR-002 — A target without LAYUP

The four `layup/` checks of the default branch's ruleset are posted only by
LAYUP, so a target whose ruleset requires them merges nothing when LAYUP is
absent. **Decided here** (decision 7 of #265): S13 in phase 2 writes the record
row `without-layup` (source `computed`, ref `out/ruleset-default.json`), whose
value is the line "To merge with LAYUP absent, remove the required checks
layup/gates, layup/spec, layup/verify and layup/rules from the ruleset "layup: the
default branch"." The gate jobs stay required, so the target still passes
its own gates with LAYUP absent. The name `without-layup` is a new row name of
the built block `setup-record`, whose `name` column is text: its build task
adds it with its test.

## NFR-003 — No value without evidence at the Scaffold

Every value of the setup under `layup run` has its source in the records:

- the name, the visibility, the baseline and the pin: `start.tsv`, by the ref
  `records:start/start.tsv`;
- the stack and each marker's value: the Operator's answer in the batch, by
  the URL of its comment, whose copy is in `copies/`; an accepted proposal
  keeps its quote of the problem statement in `intake/markers.tsv`, which code
  checked;
- a marker with no answer: `gap`, kept as an open gap with its row in
  `docs/setup/open-gaps.tsv` (Invariant 4), never a guessed value;
- a prose text: the session that wrote it (`sessions.tsv`) and its file in
  `records:setup/inputs/files/`, with the check of its row;
- the `integration_id` of the `layup/` checks: `app_id` of the forge register;
  of the gate jobs: GitHub Actions' ID, as in phase 1.

Check `sources` of `layup setup verify` holds as in phase 1: each `answer` ref
is a row of `answers.tsv`.

## NFR-006 — The pin at the Scaffold

The pin of a target is that of `start.tsv` (`run.md`, NFR-006). S02 in phase 2
takes it from there and resolves nothing; S04 writes
`docs/setup/armature.pin` from it on `layup-setup`, so the root commit stays
the unmodified copy; check `pin` compares the file with the record rows, as in
phase 1.

## The adoption of a phase-1 target

K41 (O-163): a target that phase 1 set up, whose records branch holds the
Operator's commit, was given to `M2d`. **Decided by the Operator** (O-200,
comment 6110498133 of #265, decision 1 (a)): `M2d` adopts no phase-1 target.
`layup run` drives the Scaffold only on a target that its own Start made
(`--new`, O-163); the restart refuses any other target at its step `clone`, as
the records of Start (`start/start.tsv`) are missing (`run.md`, The restart).

A target that phase 1 set up stays as phase 1 left it: its default branch's
ruleset requires only its own gate jobs, so it merges with LAYUP absent
(`records.md`, NFR-002 item 3); its records branch has no ruleset of the App.
To bring it under `layup run`, the Operator starts again on an empty
repository with `layup run --new`. Reason (decision 1 of #265): an adoption is a
second path through each step of Start, Intake and the Scaffold, as no phase-1
run made `start.tsv`, `approvers.tsv`, the copied answers or the records of
Intake; and the first pilot showed what a second run on a set-up target costs
(F-31, F-35, F-37, F-38).

## The findings of the first pilot

The plan hosts these findings in `M2d`
([`docs/plan/README.md`](../plan/README.md#the-findings-of-the-first-pilot-that-are-not-fixed)).
**Decided here** (decision 9 of #265):

| Finding | Settled |
| ------- | ------- |
| F-2: a changed input of a done step needs a new work area | Under `layup run`, the inputs are copies of records that do not change, and the work area is built again from them ([The work area](#the-work-area)), so no input of a done step changes. Phase 1 keeps its rule. |
| F-6: `layup setup verify` before S12 exits 2 | Not reached under `layup run`: the step `verify` runs after S14. A call by hand before S12 stays exit 2, a known limit of phase 1. |
| F-13: `commands.sh` is not fail-fast | `set -eu` ([The commands for the Operator](#the-commands-for-the-operator)), in both phases. |
| F-16: a stale line that the setup does not flag | A known limit: no rule finds a stale line that check `adapted` does not flag without a guess. Its fix is a change of the baseline, or of the list of `adapted`, which is LAYUP's rule. |
| F-19: two ruleset parameters that S13 does not set | Both set ([The ruleset of the default branch](#the-ruleset-of-the-default-branch)), in both phases. |
| F-31, F-35: a second setup run on a target | No path in phase 2 ([the adoption](#the-adoption-of-a-phase-1-target), O-200): under `layup run` the Scaffold runs once per target (its steps are "done before" once their records exist), and no phase-1 target is adopted. |

## The steps of the Scaffold

**Decided here** (decision 4 of #265): the restart runs, after
`follow-up-answers`, the steps `setup`, `verify`, `scaffold-push` and
`forge-check`, and the step `markers` of Intake runs between `completeness` and
`batch`. The record of each step: `intake/markers.tsv`; the files of
`setup/inputs/files/` (by `setup`, which writes no other record: its result is
the work area); `setup/record.tsv` with the files of
[the step verify](#the-step-verify-and-the-records-of-the-setup);
`setup/rulesets.txt`; `setup/forge-check.tsv`. Each step whose record exists is
`done` with "done before", as the steps of Intake are
([`intake.md`](intake.md#the-steps-of-intake)), but `setup`, which is "done
before" when `setup/record.tsv` exists, and `forge-check` at a later start
([At each later start](#at-each-later-start)). `scaffold-push` that waits is
`done` with "waiting", the last row of the restart, which then releases the
lease and exits 0. When `forge-check` at the tail is done, the run releases the
lease and exits 0; at a later start the check is not the last row, and the run goes
on through the steps that are done before and releases the lease after its last row
([At each later start](#at-each-later-start)). The next phase is Shape (`M2e`).

The block `run-steps` is built: each build task that adds a step adds its name
to the enum, with its Go schema, and the step to the restart's order, in one
change, as `M2c` did.

## Input states

**Decided here**, for each state that the happy path never makes (the lesson
"a specification that fixes one input state at a time"):

| State | Result |
| ----- | ------ |
| A restart on a target that phase 1 set up (a `layup-records` with no `start/start.tsv`) | `clone` is `fail`, as today (`run.md`, The restart); no adoption (O-200) |
| No marker in the pinned tree | `markers` writes a table with no row and starts no session; the batch asks `S01-stack` only in the Operator's block |
| No admitted pair of `systems-architect` (for `markers`) or `product-owner` (for a prose session) | the step is `fail` (`pair`), as in Intake |
| A proposal's quote that occurs nowhere, more than once, or holds a tab | the output of the session is refused; `markers` is `fail` |
| `S01-stack` with no answer, or `open`, after the follow-up | `setup` is `fail`, naming `S01-stack`; no setup commit is kept |
| A stack that the catalog does not have | the runner's own result (S01 in phase 2, "the stack has no entry in the catalog"; S12 in phase 1), and `setup` is `fail` with it, and no setup commit is kept |
| A prose session's output that is missing, a link, over 1 MiB, not UTF-8, or at a path that the stop does not name | `setup` is `fail`, naming the session; nothing of it is committed |
| A prose text whose check fails | `setup` is `fail`, naming the check; the session's files are removed from the records, and a restart starts a new session |
| A prose text that holds a marker that the pinned tree does not | S10 stops on an `M-` row that the batch did not ask: `setup` is `fail`, and the session's files are removed, as above |
| A row of `layup setup verify` that is not `pass` or `clear` | `verify` is `fail`, naming the check and its reason; no records commit of the setup is published |
| The default branch's head is neither the root commit nor the verified head when the Operator's commands run | the first command of `commands.sh` stops it with exit 1, naming the head; `scaffold-push` keeps "waiting" |
| `commands.sh` stopped after the push, or after one apply | it is run again as it is: no second push, and only the apply whose output file is missing |
| The default branch's head is the verified head, and an apply output is missing | `scaffold-push` keeps "waiting" |
| An apply output that is not a JSON object, or whose `source` is another repository | `scaffold-push` is `fail`, naming the file |
| A ruleset in `evaluate` or `disabled` | its rules are missing in the read-back: `forge-check` is `fail` |
| A probe that the forge accepts | `forge-check` is `fail`, "probe accepted" |
| A probe that fails for another reason (token, network, server) | `forge-check` is `fail`, "probe not run" |
| The activity read fails | `forge-check` is `fail` |
| A private repository on GitHub Free (no enforced rulesets) | the read-back finds no rule: `forge-check` is `fail` (K15; the plan check of Start refuses it before) |
| A restart while `scaffold-push` waits | it runs the steps again from the records: `setup` and `verify` are "done before", and `scaffold-push` reads the forge again |
| A work area that is missing when `scaffold-push` runs | it is built again from the records; a head of `layup-setup` other than the last commit of `setup/tree.tsv` is `fail` |

## REQ-015 and REQ-017 — The review of the release of M2d

As for `M2c` ([`intake.md`](intake.md#req-015-and-req-017--the-review-of-the-release-of-m2c)).
**Decided here** (decision 10 of #265): a row of the slicing task of `M2d`,
before the demo, records a code review of the release of `M2d` by a reviewer
of a model that wrote none of it, with `runs/T-efmy/release-check.sh` run on
the release commit and adapted where the code of `M2d` needs it (the push of
the probe reaches a remote). The demo is the next row; its inputs (the test
target, the Operator's admin login for the apply) are asked in it.

## Not in M2d

The adoption of a target that phase 1 set up (none, O-200). The commit statuses `layup/gates`, `layup/spec`, `layup/verify` and
`layup/rules` themselves, the first merge after the setup and the tree rule
that follows it (`M2e`); `layup audit` and its judgement of the activity
(`M2f`); the activation of the `pending` kinds at the first bet (`M2f`); the
pilot baseline block of the Intake answer (`baseline.tsv`, phase 4); the
dead-man job (`M3d`).

## The acceptance tests of M2d

Each test names its level ([`test-levels.md`](../tests/test-levels.md)). The
six rows of the inventory item `later-p2-scaffold` are the rows marked
"inventory"; each is kept.

| Test | Level | What it shows |
| ---- | ----- | ------------- |
| The stops of the runner under `layup run` (inventory) | unit | A run that does not stop goes on; a stop with `M-` rows adds the rows of the markers that the batch asked and runs again, and one with an `M-` row that it did not ask is `fail`; a marker of the pinned tree that S05 deletes gets no row and fails nothing; a stop with `F-` rows starts a prose session; the one stop `O-verify` of S15 ends the step `done`, and a table that holds it with another row, or an `O-verify` of another step, is `fail`; a stop with any other row, an input error, a `fail` or a `not-active` make `setup` `fail`; the same stop twice in a row is `fail`; a `fail` after a prose session of the run, also after the answers of the markers that followed it, removes that session's files from the records, one before it removes none, and a session whose rows the runner passed keeps its files when a later one fails. |
| The steps of the setup in phase 2, and the pin | unit | S01 reads only `S01-stack` and runs no `layup psb check`; S02 writes the pin rows of `start.tsv` with ref `records:start/start.tsv` and resolves nothing; S03 makes no commit and fails on a tree other than `pin.tree`; S04 writes `docs/setup/armature.pin` from those rows and the raw fact of the `S01-` and `Q-` rows; S15 makes no commit of `layup-records`. |
| The row `without-layup` | unit | S13 in phase 2 writes it with its fixed line and source `computed`; in phase 1 it writes none. |
| The apply outputs | unit | Both files present and each a JSON object of `OWNER/NAME`: `setup/rulesets.txt` holds both after their lines; a file missing: "waiting"; a file that is not a JSON object, or of another repository: `fail`, naming it. |
| The forge check at each later start | integration | With the records of a done Scaffold ([the rule](#at-each-later-start)), a restart runs `forge-check` after `lease` and before `probe` (`5/21`), writes the table again, and stops at a row that is not `pass`, with no row after it; the steps of Intake and of the Scaffold are "done before", but `follow-up` and `follow-up-answers`, which keep their own details; a table that is absent, refused or with a row that is not `pass` gives the first-start chain (`21/21 forge-check`, the table written again). |
| The answers of the setup | unit | Each row of the table of [The answers of the setup](#the-answers-of-the-setup): `accept` with and without a proposal, `gap`, a value, no answer and `open`, for a marker and for `S01-stack`; the answer of a follow-up that asks it again; the source URL and `by`; the `Q-` rows. |
| The step markers | unit; integration | The scan of a fixture tree gives the keys of S10, one row for a text that a file holds on two lines; the session's `result/markers.tsv` with an unknown key, a key twice or an empty quote is refused; a quote that occurs once passes, and one that occurs twice, nowhere or with a tab refuses the file; with no marker no session starts (the integration level). |
| The setup questions in the batch and the comment | integration | With `intake/markers.tsv` of two markers, one with a proposal and one with none, and the httptest forge: the step `batch` merges, after the questions of the idea owner's block and in the Operator's block, first `S01-stack` and then one question per row of `intake/markers.tsv`, in its order, each a row of `questions.tsv` with `source` `setup`, `ref` its ID (`S01-stack`, `M-<x8>`), `kind` `setup`, and `start` and `end` the quote's for a marker with a proposal, else `—`; `S01-stack` is asked whatever rule G1 found; the comment writes a marker's question as `Q-NNN: the value of <marker> at <file>:<line>`, with `proposed: <value>` and the quote in backquotes on the next lines when it has a proposal, and its rule of the answers adds the three forms of a marker's answer, `accept`, a value and `gap`. |
| The ruleset files (inventory) | unit | The default branch's file names each gate job pinned to 15368, the four `layup/` checks pinned to the forge register's `app_id`, `strict` true, `creation`, `allowed_merge_methods` `["merge"]` and an empty bypass list; the records branch's file has the App as its only bypass actor; phase 1 gets `creation` and the two parameters, and no `layup/` check. |
| `commands.sh` is fail-fast and can be run again | unit | Its first line is `set -eu`, and a stand-in command that fails stops the rest (F-13); with stand-ins of `git` and `gh`, a run again after the push or after one apply pushes nothing and applies only the missing ruleset, and a `gh` that writes part of its output and exits 1 leaves no output file, so the run again applies that ruleset; and a head that is neither the root nor the setup head stops it. |
| The probe (inventory) | integration | With a real local bare repository whose pre-receive hook refuses with the text of `GH013` and the ruleset's name, and the real `internal/git`: the refusal is `pass`; a hook that refuses with `GH013` and another ruleset's name, one that refuses with another text, and a remote that cannot be reached, are "probe not run"; an accepted push is "probe accepted". |
| The read-back | unit | A read-back equal to the files passes; a missing rule, a different parameter, an `integration_id` that differs, and a bypass list that is not empty (default) or not exactly the App (records) each fail, naming it; an extra rule is named and passes. |
| `layup run` with the real step runner (inventory) | integration | On a fixture baseline: `answers.tsv` is written from the copied answers, a prose stop starts a scripted session, the runner finishes, and `layup run` commits `setup/record.tsv`, `setup/verify.tsv`, `rule-paths.tsv`, `setup/tree.tsv` and the README of the Scaffold; a second run is "done before"; a rebuilt work area gives the same commits, also after the default branch holds the setup commits (its `main` is the root commit). |
| The steps of the Scaffold, end to end (inventory) | e2e | In CI on a local fake forge and local bare repositories: Start, scripted Intake answers with the setup questions, the setup, the printed commands run by the test in place of the Operator, and the forge check, whose fake forge answers the read-back and the activity, and whose bare repository refuses the probe with `GH013`. |
| The demo (inventory) | uat | On a real GitHub test target, the Operator runs the printed pushes and the ruleset apply, and sees the push to `layup-probe` refused and `setup/forge-check.tsv` show each rule and probe as passed. |
