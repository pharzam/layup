# Intake: the gap check in one batch, and the fact spans

The specification of milestone `M2c` of [the plan](../plan/README.md#milestones):
what `layup run` adds after Start and the probe, so that the gaps of a problem
statement reach the idea owner as one batch, and its facts are numbered spans
that code checks. Its requirements are the rest of `REQ-001` (the review of
meaning, the one batch), the Intake part of `REQ-012` (the numbering, the span
check, the draft, the completeness review), and `NFR-003` at Intake; the design
is [`architecture.md`](../architecture.md) §5 (The gap check, in one batch) and
§7 (steps 1 to 5), and [ADR-0018](../adr/0018-derive-the-specification-from-numbered-source-lines.md).
Task `T-g49f` (#206) wrote it, after a public-solution search
([`runs/T-g49f/search/`](../../runs/T-g49f/search/)), whose result is that no
public solution can be its base; the decisions are comment 6100128462 of #206.
The records are in [`records.md`](records.md#nfr-001--the-records-of-intake).

**What the search gives this file.** The span of a fact is a start and an end
offset with its text, which code checks against the file, as the brat standoff
format and the W3C `TextPositionSelector` and `TextQuoteSelector` do; no public
tool checks that each non-blank byte is in exactly one span, so the cover rule
below is this file's. The kinds of a gap of meaning take the smells of ISO/IEC/IEEE
29148 and the ambiguity handbook of Berry, Kamsties and Krieger as their
source for the kinds that name the text, and Shape Up's "rabbit holes" for the
rest (`architecture.md` §5).

## REQ-001 — The gap check of Intake, in one batch

Requirement: "Before delivery starts, `layup psb check` finds the rule gaps of a
problem statement, a review session finds its gaps of meaning, and both go to
the idea owner as one batch of questions, whose answers are stored as a raw
fact." Phase 1 delivered the rule gaps ([`psb-check.md`](psb-check.md)); this
section delivers the rest.

### The sessions of Intake

**Decided here** (decision 5 of #206): a session of Intake is a role session of
`M2b` ([`session.md`](session.md#req-013--a-role-session)) with no attempt and
no event, a **session of a phase**:

- its task ID is drawn as a probe's (`T-` and four random characters, drawn
  again until no row of `sessions.tsv` holds it), attempt 1, the head of the
  default branch as its base, and its start row in `sessions.tsv`;
- its pair is the first admitted pair of the routing register for its role and
  the tier `reasoning` ([`session.md`](session.md#the-routing-register)); its
  role is that of the step table of `architecture.md` §9: `systems-architect`
  for the review of meaning and of the answers, `qa-engineer` for the
  completeness review, and `product-owner` for the numbering and the draft;
- its directory ([`session.md`](session.md#the-session-directory)) gains
  `input/`, made with the directory and before the start row: the files that
  its step reads, byte for byte from the records (`problem-statement.md` from
  `records:start/problem-statement.md` for each step; `spans.tsv` and
  `draft.tsv` for the completeness review; `questions.tsv` and `answers.tsv` for
  the review of the answers); **decided here**, as `repo/` is the target's
  default branch, which holds none of them;
- its prompt is fixed text built by code, with the absolute paths of its input
  files and of its output file under `result/`;
- it writes one output file in `result/` (`meaning.tsv`, `spans.tsv`,
  `draft.tsv`, `coverage.tsv` or `answers-review.tsv`) in the form of its block
  in [`records.md`](records.md#nfr-001--the-records-of-intake) without the
  columns that code adds; `layup run` reads the file with the rules of a result
  file (a link as the last part of its path is refused; at most 1 MiB;
  [`session.md`](session.md#req-005--the-result-of-a-session)), checks it, and
  commits its records with the session's telemetry row in one records commit,
  then posts the session's comment, as a probe does;
- a session whose end is not `done`, or whose file its block or its check
  refuses, writes no record of its file; the step is `fail`, naming the session
  and the reason, and the run stops (a restart runs the step again with a new
  session).

**Reason:** a task, an attempt and its events are the task loop's (`M2e`);
Intake has none, and a probe has the same frame. **Rejected:** a task session
for each step; one session for every step (the completeness review must run on
another harness, below).

### The budget of Intake

**Decided here** (decision 3 of #206): Intake and Shape are milestone 0, whose
cap is `intake.cap` of `start.tsv` (`MONEY,HOURS`). Before each session of
Intake starts, code sums the milestone's known money (each telemetry row of a
session of Intake whose `money_status` is `reported` or `computed`) and adds
the spend cap of the new session; when the sum passes `MONEY`, or the time from
the first records commit of Intake to now passes `HOURS`, the step is `fail`
("the cap of milestone 0") and no session starts. A session whose money is not
known after it ends adds nothing to the known money and is recorded `unknown`,
as `session.md` writes it. **Reason:** the cap of milestone 0 is the one budget
value that exists before the answers (`architecture.md` §12). The band `B` and
`U` and the appetite come with the answers and go to `budget.tsv` (below); no
rule of `M2c` reads them.

### The rule gaps

The first step of Intake runs the rule check of `layup psb check`
([`psb-check.md`](psb-check.md)) on `records:start/problem-statement.md`, in
process (`internal/psb`), and commits its table byte for byte as
`records:intake/psb.tsv`, in the form of the block `psb-gaps`. **Decided here**
(decision 4 of #206, F-21 of the first pilot): rule G1 is not changed; a stack
gap that G1 finds is a question of the batch like any rule gap, and the review
of meaning and the merge's duplicate removal reach the case where the brief
names the stack in other words. The stack question of setup step S01 joins the
batch with the Scaffold (`M2d`).

### The review of meaning

A session of the role `systems-architect` reads the problem statement and writes
`result/meaning.tsv`: one row per gap of meaning, with its kind, its quote and
its question. The kinds (`architecture.md` §5, gap check step 2): `two-readings`
(a term with two readings), `unmeasurable` (a metric that cannot be measured as
written), `missing-intent` (a missing intent value), and the four questions of
Shape Up's "rabbit holes": `new-work` (new technical work), `fit` (an assumption
about how parts fit), `no-design` (a design that may not exist), `decision` (a
hard decision to settle before the work).

**The quote check** (code): a quote is text of one line of the problem
statement, byte for byte, with no line break; code finds it in the file and
writes its byte offsets. **Decided here:** a quote that occurs nowhere, or more
than once, refuses the file (the session must quote enough to be unique), and
so does a quote that holds a tab. **Reason:** a TSV field holds no line break,
and a quote that occurs twice names no one place. A gap whose text spans lines
quotes the part of one line that names it best. The committed rows are
`records:intake/meaning.tsv`, each with an ID `M-NNN` in the order of the rows.

## REQ-012 — The facts and the draft at Intake

Requirement: "LAYUP derives identified requirements and technical
specifications from an approved problem statement; each requirement has a trace
to the text that states it and an acceptance criterion." At Intake, before the
batch, the problem statement is numbered into spans, checked, and drafted into
requirements, whose gaps join the batch (`architecture.md` §7 steps 1 to 4). The
revision after the answers, the first bet and the confirmed inventory are
`M2e`'s.

### The numbering

A session of the role `product-owner` splits the problem statement into spans and
writes `result/spans.tsv`: each span by its byte offsets (from 0; the end is
exclusive), and for a **fact** its number and its class (`need`, `constraint`,
`measure`, `context`), for a **not a fact** span (a heading, a rule line) its
reason. A span may be part of a line. The numbers of the facts are 1 to n in the
order of the spans.

### The span check

The span check is code, in `internal/spec`. It passes when:

1. **Cover:** every byte that is not white space (ASCII space, tab, CR, LF) lies
   in exactly one span, and no span holds only white space;
2. **Bounds:** each span has `0 <= start < end <= size`, and the spans are in
   the order of `start` and do not overlap;
3. **Text:** each span's text is the file's text at its offsets (the session
   writes no text; code reads it, so this rule holds by construction of the
   committed record, and the check reads the file again);
4. **IDs:** the fact numbers are 1 to n, each once, in the order of the spans;
5. **Classes:** each fact's class is one of the four, and a not-a-fact span has a
   reason and no class.

It prints its table (the block `spec-check`): one row per rule, `pass` or `fail`,
with the first broken place in `detail`.

**The command** (**decided here**, decision 11 of #206):
`layup spec check --facts SPANS --psb FILE`. `--facts` names the span table and
`--psb` the problem statement; both are required and take a value, as each
flag of "Commands" does ([`README.md`](README.md#commands)). Exit 0 when each
row is `pass`, 1 when a row is `fail`, 2 on a usage or input error (a file that
is missing, unreadable, not valid UTF-8, or a span table that its block
refuses). **Reason:** a flag with no value is a usage error in the frame of
`layup`, and the check needs both files. `internal/cli` gains the row of the
command; `internal/spec` joins the line of the engine checks
([`packages.md`](packages.md#the-table-of-m2c)): code, no model call.

**The step:** `layup run` calls the same check in process after the numbering,
commits `records:spec/spans.tsv` (with the text of each span read from the
file) and the table as `records:spec/check.tsv`, and posts no batch and stops
when a row is `fail` (`architecture.md` §7 step 2).

### The draft

A session of the role `product-owner` (routing puts the numbering and the draft
on one harness: the first admitted pair of that role) writes `result/draft.tsv`:
one or more requirements for each `need` fact and one non-functional requirement
for each `constraint` fact, each with an ID, a statement, `covers` (the fact
numbers) and an acceptance criterion; one success criterion for each `measure`
fact; no priority. **The check** (code): each `covers` number is a fact of
`spec/spans.tsv`, each `need` and `constraint` fact is covered at least once,
and each `measure` fact has a success criterion; a row that breaks it refuses
the file. The rows are `records:spec/draft.tsv`.

### The completeness review

A session of the role `qa-engineer` on an admitted harness **other than that of
the numbering and the draft** reads the problem statement, `spec/spans.tsv` and
`spec/draft.tsv`, and writes `result/coverage.tsv`: each need, constraint or
measure that no requirement or success criterion covers (`uncovered`), each
requirement that no fact supports (`unsupported`), and each fact whose class it
doubts (`class`), each with a question. **Decided here** (decision 1 of #206):
with one admitted harness the review does not run; the step is `fail` ("not
run: one admitted harness") and no batch is posted. **Reason:** an independent
review is the point of the step (`F-0003#52`), and `M2b` admits several
harnesses. Each row is `records:intake/coverage.tsv`, and each becomes a
question of the batch.

## REQ-001 — The one batch

### The merge and the IDs

Code merges into `records:questions.tsv`, in this order: the rule gaps
(`intake/psb.tsv`), the gaps of meaning (`intake/meaning.tsv`), the items of the
completeness review (`intake/coverage.tsv`), and the intent form (below). A
question whose quote and kind repeat those of a question before it is dropped
from the batch, and its source row stays in its own record (**decided here**:
the quote's offsets and the kind are compared; a rule gap's quote is its line,
trimmed of its line break). Each kept question gets an ID `Q-NNN` from `Q-001`, in the order of the
merge; a question of the batch is never renumbered.

### The comment of the batch

One comment on the Intake issue holds the batch, in two blocks (**decided here**,
decisions 6, 9 and 10 of #206):

1. **For the idea owner:** each question, one line each:
   `Q-NNN: <question>` and, on the next line, its quote in backquotes; then the
   intent form, one line per field: the appetite in money (`Q-NNN: appetite,
   money, a decimal in USD`), the appetite in wall-clock hours, and the band `B`
   and `U` (decimals in USD, `B` below `U`).
2. **For the Operator:** the setup questions. In `M2c` the block is empty and
   says so; the Scaffold (`M2d`) fills it with the questions of setup steps S01
   and S10 and the proposed marker sources, in the same merge.

The comment ends with the rule of the answers: one comment per block, one line
per question ID, `Q-NNN: <answer>`.

### The answers

`layup run` copies each comment of the Intake issue (`copies.tsv`,
[`records.md`](records.md#nfr-001--the-records-of-start)). An answer is a line
`Q-NNN: <answer>` of a comment whose author's numeric ID is in `approvers.tsv`
in the role of the question's block (the idea owner, or the Operator); any other
comment is recorded as input and answers nothing (`architecture.md` §3). The
last answer line of a question ID, in the order of the copies, is its answer.
Code checks that each ID of the batch has an answer, and that the answer of each
field of the intent form parses (a `decimal`, and `B` below `U`). The answers are
`records:intake/answers.tsv`; when every field of the intent form is answered,
code writes `records:budget.tsv`.

### The review of the answers and the follow-up

When each question has an answer, a session of the role `systems-architect`
reads the questions and the answers and writes `result/answers-review.tsv`: for
each question, `closed` or `open` with a reason. The rows are
`records:intake/answers-review.tsv`. Then code posts **one** follow-up comment on
the Intake issue that asks again only for the questions with no answer, an
answer that does not parse, or an `open` answer; each is a new row of
`questions.tsv` with the source `follow-up`, which names the question it asks
again. With none, Intake is done. **Decided here:** one follow-up round in `M2c`;
a question still open after it is open in `questions.tsv`, and Shape (`M2e`)
reads it. **Reason:** `architecture.md` §5 step 6 asks again "in one comment",
and a loop with no bound waits on a human with no end.

### The raw fact of the answers

The answers are a raw fact of the target, one fact per question ID
(`architecture.md` §7 step 5). In `M2c` the fact is `records:intake/answers.tsv`
with the copied comments that hold it (`copies/`), byte for byte; the file of
the fact under the target's `docs/facts/` is written by the Scaffold (`M2d`,
`architecture.md` §7 step 8), as `docs/facts/` is a rule path that no role
session may edit. **Decided here** (decision 6 of #206). Each question keeps its
quote, so a requirement that later covers an answer's fact can name the quote's
fact (`M2e`).

## The steps of Intake

**Decided here** (decision 11 of #206): when each step of Start is done, the
restart runs the steps of Intake after `probe`, in place of the row `phase`:
`rule-gaps`, `meaning`, `numbering`, `spans`, `draft`, `completeness`, `batch`,
`answers`, `answers-review`, `follow-up`. Each step whose record exists is
`done` with "done before" and runs nothing; a step whose input is not there yet
(no answer to the batch, or to the follow-up) is `done` with "waiting" in
`detail`, and the run releases the lease and exits 0. A restart goes on from
the records. The block `run-steps` is built, so the build task that adds these
steps adds them to its enum, with its Go schema, in one change, as `M2b` did
for `probe` ([`run.md`](run.md#the-command)).

## NFR-003 — No value without evidence at Intake

Every value that Intake writes has its evidence in the records: each quote by
its offsets in `start/problem-statement.md`; each span by its offsets; each
answer by the comment ID of its copy; each value of `budget.tsv` by the question
ID of its answer. A value of the intent form that does not parse is asked again,
never guessed.

## Input states

**Decided here**, for each state that the happy path never makes:

| State | Result |
| ----- | ------ |
| A routing register with no admitted pair for a role of Intake | the step is `fail` (`pair`), no session starts |
| One admitted harness | `completeness` is `fail` ("not run: one admitted harness"); no batch |
| A session's output file that is missing, a link, over 1 MiB, or that its block refuses | the step is `fail`, naming the session; nothing of the file is committed |
| A quote that occurs nowhere, more than once, or holds a tab | the file of the review is refused (as above) |
| A span table that fails the span check | `spans` is `fail`; `spec/check.tsv` is committed; no batch |
| A draft row whose `covers` names no fact, or a `need` or `constraint` fact with no requirement | the draft is refused |
| The cap of milestone 0 reached | the step is `fail` ("the cap of milestone 0"); no session starts |
| A comment with answer lines by an account that is not an approver, or in the other block's role | recorded as input; it answers nothing |
| An answer of the intent form that does not parse, or `B` not below `U` | the question is asked again in the follow-up |
| A question still open after the follow-up | open in `questions.tsv`; Intake is done |

## Not in M2c

The setup questions of steps S01 and S10, the proposed marker sources with their
quotes and their check (`M2d`, decision 9 of #206); the file of the raw fact
under `docs/facts/` (`M2d`); the success criteria in the intent form, the
approver of each approval point and the role `approver` (`M2e`, decision 10);
the revision of the draft with the answers, the check that a requirement that
covers an answer's fact names the quote's fact (`architecture.md` §7 step 5),
the first bet, the rendered commit and the confirmed inventory (`M2e`); the
owner map (phase 3); the smart-if provider and its authority (`M3a`); the
actions of the band `B` and `U` and the accepted bounds of an unknown part
(`M3a`, `M3b`); the caps of later milestones (`M2e`); the pilot baseline block of
the Intake answer (`architecture.md` §5; `baseline.tsv`, phase 4); the Early
Question Share (`M4b`).

## The acceptance tests of M2c

| Part | Level | Test |
| ---- | ----- | ---- |
| The sessions of Intake | integration | With scripted harness programs and the real `git`: a session of a phase has its own task ID, a start row and no event, its output read by the rules of a result file, its records and its telemetry row in one commit, then its comment. |
| The budget of Intake | unit | The sum of the known money and the new session's cap against `MONEY`, and the hours against `HOURS`: the step fails at the cap and starts no session. |
| The rule gaps | integration | `intake/psb.tsv` is the table of `layup psb check` on the problem statement, byte for byte. |
| The review of meaning | unit | The quote check refuses a quote that occurs nowhere, twice, or holds a tab, and gives the offsets of one that occurs once. |
| The span check | unit | It fails when a non-blank byte is in no span or two, a span is out of bounds or out of order, a span's text differs from the file, a fact number repeats or skips, or a class is not one of the four; it passes on a span table that covers the file. |
| The command `layup spec check --facts` | e2e | The binary prints the table and exits 0, 1 or 2 as above. |
| The draft | unit | The draft is refused when `covers` names no fact or a `need` or `constraint` fact is not covered, and kept otherwise. |
| The completeness review | integration | With two admitted harnesses, the review runs on the other one; with one, the step is `fail` and no batch is posted. |
| The merge and the IDs | unit | A question whose quote and kind repeat is dropped; each kept question gets one ID, in the order of the merge. |
| The comment and the answers | integration | With the httptest forge: one comment with two blocks; answers by an approver in its role are kept, others recorded as input; a missing or unparsed answer goes to the follow-up; the answers and `budget.tsv` are committed. |
| The steps of Intake | e2e | In CI with no secret, the binary, a local fake forge and scripted harness programs: a restart runs the steps to `batch` and exits "waiting"; after scripted answer comments, the next restart reads them, runs the review, posts one follow-up for the missing answer, and the next one ends Intake. |
| The demo | uat | On a pilot problem statement that the idea owner chooses, on real harnesses: one Intake comment holds the rule gaps and the gaps of meaning, and `layup spec check --facts` passes on the numbered spans. |
