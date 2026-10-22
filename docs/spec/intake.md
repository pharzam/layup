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
  the tier `reasoning` ([`session.md`](session.md#the-routing-register)), but
  that of the completeness review (below); **decided here** (task `T-gmey`),
  a harness is read at the version of its last row of `records:harnesses.tsv`,
  as no version command runs before the session; the session's own version
  check then reads the version, and probes it when its last probe did not
  pass; its
  role is that of the step table of `architecture.md` §9: `systems-architect`
  for the review of meaning and of the answers, `qa-engineer` for the
  completeness review, and `product-owner` for the numbering and the draft;
- its directory ([`session.md`](session.md#the-session-directory)) gains
  `input/`, which `internal/session` makes right after the rest of the
  directory (`prompt.md` its last file), and before the start row (one
  reading, task `T-gmey`): the files that
  its step reads, byte for byte from the records (`problem-statement.md` from
  `records:start/problem-statement.md` for each step; `spans.tsv` for the
  draft, **decided here** (task `T-2et0`), as its rows name fact numbers;
  `spans.tsv` and `draft.tsv` for the completeness review; `questions.tsv` and `answers.tsv` for
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
  then posts the session's comment, as a probe does; **decided here** (task
  `T-gmey`), the first line of that comment is `<session>: <role> session of
  <step>: <class>`, with `, refused file` when its file is refused, for
  example `S-1a2b3c4d: systems-architect session of meaning: done`;
- a session whose end is not `done`, or whose file its block or its check
  refuses, writes no record of its file; the step is `fail`, naming the session
  and the reason, and the run stops (a restart runs the step again with a new
  session). **Decided here** (task `T-gmey`): its records commit holds its
  telemetry row alone, as the telemetry record has a row for each session that
  started ([`records.md`](records.md#req-011--the-telemetry-record)).

**Reason:** a task, an attempt and its events are the task loop's (`M2e`);
Intake has none, and a probe has the same frame. **Rejected:** a task session
for each step; one session for every step (the completeness review must run on
another harness, below).

### The budget of Intake

**Decided here** (decision 3 of #206): Intake and Shape are milestone 0, whose
cap is `intake.cap` of `start.tsv` (`MONEY,HOURS`). Before each session of
Intake starts, code sums the milestone's known money (each telemetry row since
Start whose `money_status` is `reported` or `computed`, the probes' rows
included: **decided here**, as milestone 0 runs from Start and a probe is
spend before the first bet) and adds
the spend cap of the new session, or nothing for a session whose register row
has no cap (`—`; its spend is an unknown part, [`session.md`](session.md#the-limit-of-a-session),
whose actions are not in `M2c`); when the sum passes `MONEY`, or the time from
Start (the time of the first commit of the records branch, as milestone 0's
clock starts at Start, `architecture.md` §12) to now passes `HOURS`, the step is `fail`
("the cap of milestone 0") and no session starts. **Decided here** (note 5 of
#210): the clock of milestone 0 does not pause while the batch waits for a human,
as milestone 0 has no task whose wait could pause it; the Operator raises the cap
on the control issue (`architecture.md` §12). A session whose money is not
known after it ends adds nothing to the known money and is recorded `unknown`,
as `session.md` writes it. **Decided here** (task `T-ybd5`): `MONEY` is US
dollars (`run.md`, `--intake-cap`), and a sum or a time **passes** its cap when
it is above it (equal is inside); the time of Start is the committer time of the
root commit of the run's records base; the rule is asked at the start of each
session of a phase, before its task ID is drawn, so a refused session leaves no
record; a probe that the session's own version check runs after the rule
counts from the next session. **Reason:** the cap of milestone 0 is the one budget
value that exists before the answers (`architecture.md` §12). The band `B` and
`U` (the money appetite) and the wall-clock appetite come with the answers and go to `budget.tsv` (below); no
rule of `M2c` reads them.

### The rule gaps

The first step of Intake runs the rule check of `layup psb check`
([`psb-check.md`](psb-check.md)) on `records:start/problem-statement.md`, in
process (`internal/psb`), and commits its table byte for byte as
`records:intake/psb.tsv`, in the form of the block `psb-gaps`. **Decided here**
(decision 4 of #206, F-21 of the first pilot): rule G1 is not changed; a stack
gap that G1 finds is a question of the batch like any rule gap, so a false
positive of G1, as F-21 was, reaches the idea owner as a question that the
answer closes; a second text pattern for G1 would be a guess at the next brief's
wording. The stack question of setup step S01 joins the batch with the Scaffold
(`M2d`).

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
**Decided here** (task `T-1u2v`): the form of the session's file is one Go
value derived from the block `meaning` by dropping the columns that code adds,
with no block of its own; its header is the remaining columns in the block's
order (`kind`, `quote`, `question`), and a key column that is kept stays a key.
A `fail` of the step releases the lease, by the rule of
[`run.md`](run.md#the-lease-and-fencing), The release at the end of a restart.

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
order of the spans. **Decided here** (task `T-hiv9`): the session's file is the
block `spans` without `text`, which code adds from the file at the span's
offsets; a span partly outside the file gets the text of its part in the file,
and the rule Bounds of the check reports it; a span with no byte in the file has
no text, and a repeated `start` repeats the block's key, so the block refuses
either, and with it the file.

### The span check

The span check is code, in `internal/spec`. It passes when:

1. **Cover:** every byte that is not white space (ASCII space, tab, CR, LF) lies
   in exactly one span, and no span holds only white space (**decided here**: a
   span of white space only states nothing);
2. **Bounds:** each span has `0 <= start < end <= size`, and the spans are in
   the order of `start` and do not overlap (**decided here**: one order, so
   that the numbers of the facts follow the text);
3. **Text:** each span's text is the file's text at its offsets (the session
   writes no text; code reads it, so this rule holds by construction of the
   committed record, and the check reads the file again and compares after the
   field rule of the column `text`, each tab, CR and LF as one space);
4. **IDs:** the fact numbers are 1 to n, each once, in the order of the spans
   (**decided here**: ADR-0018 asks unique IDs; a number in the order of the
   text gives one form, as `F-0003` numbers its facts);
5. **Classes:** each fact's class is one of the four and it has no reason, and a
   not-a-fact span has a reason and no class (the two pairings of the block
   `spans`, each with `fact`; task `T-4ywt`).

A rule reports only its own first broken place, so one table can fail more
than one row. **Decided here** (task `T-4ywt`): the rule Text skips a span that
breaks Bounds, which has no text of the file to compare, and Bounds reports it;
a field `text` of `—`, which reads back as the empty value, is compared as the
text `—`.

It prints its table (the block `spec-check`): one row per rule, `pass` or `fail`,
with the first broken place in `detail`.

**The command** (**decided here**, decision 11 of #206):
`layup spec check --facts SPANS --psb FILE`. `--facts` names the span table and
`--psb` the problem statement; both are required and take a value, as each
flag of "Commands" does ([`README.md`](README.md#commands)). Exit 0 when each
row is `pass`, 1 when a row is `fail`, 2 on a usage or input error (a file that
is missing, unreadable, not valid UTF-8, or a span table that its block
refuses, or a table that it cannot write to standard output: **decided here**,
task `T-hpid`, as for `layup psb check`). **Reason:** a flag with no value is a
usage error in the frame of `layup`, and the check needs both files. The exit-code map of `internal/cli`
([`README.md`](README.md#commands)) gains the table `spec-check`, whose result
words are `pass` and `fail`. `internal/cli` gains the row of the
command; `internal/spec` joins the line of the engine checks
([`packages.md`](packages.md#the-table-of-m2c)): code, no model call.

**The steps:** the step `numbering` commits the session's rows as
`records:spec/spans.tsv` (with the text of each span read from the file) when
its block reads them; the step `spans` then runs the same check in process on
that record and the problem statement, commits the table as
`records:spec/check.tsv`, and posts no batch and stops when a row is `fail`
(`architecture.md` §7 step 2). **Decided here:** a `spec/check.tsv` that has a
`fail` row is the record of neither step: a restart runs `numbering` again with
a new session, whose rows replace `spec/spans.tsv`, then `spans` again, whose
table replaces `spec/check.tsv` (a record of the records branch is edited in
place, and Git history is its log, `architecture.md` §3); no step after `spans`
runs while the last `spec/check.tsv` has a `fail` row, so no batch is posted on a
failed check (ADR-0018 decision 2). **Decided here** (task `T-hiv9`): the
detail of `spans` is "the span check passed, in spec/check.tsv", or "the span
check: <rule> fail: <detail>" of its first `fail` row; a table equal to the
last `spec/check.tsv` (a rerun whose session wrote the same spans) is kept as it
is, with no records commit, as a records commit changes a file.

### The draft

A session of the role `product-owner` (routing puts the numbering and the draft
on one harness: the first admitted pair of that role) writes `result/draft.tsv`:
one or more requirements for each `need` fact and one non-functional requirement
for each `constraint` fact, each with a statement, `covers` (the fact numbers)
and an acceptance criterion; one success criterion for each `measure` fact; no
priority. **Decided here** (task `T-2et0`): the session's file is the block
`draft` without `id` and `session`, which code adds (`R-NNN` in the order of the
rows), so its header is `kind`, `statement`, `covers`, `criterion`; a row whose
`covers` is empty names no fact. **The check** (code, **decided here**: the draft's coverage
of the classes that §7 step 3 gives, checked byte for byte before a session reads
it): each `covers` number is a fact of
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
its pair is the first admitted pair of the role's list for the tier `reasoning`
whose harness is not the harness of the numbering and the draft sessions, a
rule of this step that the routing rule of
[`session.md`](session.md#the-routing-register) does not have; with no such
pair (no admitted pair of the role at all, one admitted harness, or a list with
none other: one detail for each, **decided here**, task `T-0mey`), the review
does not run; the step is `fail` ("not run: no admitted pair on another
harness") and no batch is posted. **Decided here** (task `T-0mey`): the harness
passed over is the draft's session's (the `session` column of `spec/draft.tsv`,
its row of `sessions.tsv`), as the records name the draft's session and the span
table names none, and the routing rule of The draft puts the numbering on the
same harness; the session's file is the block `coverage` without `id` and
`session`, which code adds (`C-NNN` in the order of the rows), and each `fact` of
it must be a fact of `spec/spans.tsv` and each `requirement` a row of
`spec/draft.tsv`, else the file is refused, as each row becomes a question of
the batch. **Reason:** an independent
review is the point of the step (`F-0003#52`), and `M2b` admits several
harnesses. Each row is `records:intake/coverage.tsv`, and each becomes a
question of the batch.

## REQ-001 — The one batch

### The merge and the IDs

Code merges into `records:questions.tsv`, in this order: the rule gaps
(`intake/psb.tsv`), the gaps of meaning (`intake/meaning.tsv`), the items of the
completeness review (`intake/coverage.tsv`; the quote of an `uncovered` or
`class` item is its fact's span, of an `unsupported` item none), the intent
form (below), and, from `M2d`, the setup questions of the Operator's block:
`S01-stack`, and one question for each row of `intake/markers.tsv`, in its
order ([`scaffold.md`](scaffold.md#the-answers-of-the-setup)). A
question whose quote and kind repeat those of a question before it is dropped
from the batch, and its source row stays in its own record (**decided here**:
the quote's offsets and the kind are compared; a rule gap's quote is its line,
trimmed of its line break). Each kept question gets an ID `Q-NNN` from `Q-001`, in the order of the
merge; a question of the batch is never renumbered. **Decided here** (task
`T-8r3u`): a question with no quote is never dropped, as the four fields of the
intent form share the kind `intent` and have none; a rule gap's line is trimmed
of a CR before its line feed too. In `M2c` a repeat of quote and kind is within
one source, as the kinds of the sources differ. **Decided here** (task `T-4c5c`):
a setup question is never dropped as a repeat: two markers may share the quote of
their proposals, they share the kind `setup`, and each is asked, as a marker with
no question would stop S10. **Known limit:** two rule gaps
of one line and one rule (two undefined abbreviations of one line, G3) have one
quote and one kind, so the second is dropped and its abbreviation is not asked.

### The comment of the batch

One comment on the Intake issue holds the batch, in two blocks (**decided here**,
decisions 6, 9 and 10 of #206):

1. **For the idea owner:** each question, one line each:
   `Q-NNN: <question>` and, on the next line, its quote in backquotes; then the
   intent form, one line per field, as `architecture.md` §12 gives it: the band
   `B` (go on without asking) and `U` (the upper edge, which is the money
   appetite of the delivery), decimals in USD with `B` below `U`, the
   delivery's wall-clock appetite in hours, a decimal, and the success criteria
   in the idea owner's own words, separated by `;` (`M2e`, task `T-8rm6`). The
   block ends with the naming of the approvers: one line `<point>: <login>`
   for each planned approval point the idea owner names — `setup`, `bet`,
   `activation`, `retrospective`, `early-retrospective`, `audit` — and a point
   not named is the idea owner's ([`loop.md`](loop.md#the-approvers-and-the-success-criteria),
   task `T-8rm6`).
2. **For the Operator:** the setup questions: `S01-stack`, and one for each
   marker, as `Q-NNN: the value of <marker> at <file>:<line>` and, when the
   marker has a proposal, the line `proposed: <value>` and its quote in
   backquotes ([`scaffold.md`](scaffold.md#the-answers-of-the-setup)).

**Decided here** (task `T-4c5c`): a block with no question says `No question for
the Operator.` or `No question for the idea owner.`, as a follow-up may ask one
block only. The comment ends with the rule of the answers: one comment per block,
one line per question ID, `Q-NNN: <answer>`; when a marker is asked, the rule adds
the three forms of its answer, `accept` (the proposed value), a value, and `gap`.

**Decided here** (task `T-8r3u`): the comment's first line is `LAYUP Intake:
the batch of questions`, by which a restart knows the batch was posted; a quote
is written in a run of backquotes one longer than its longest run of
backquotes, with a space inside when it begins or ends with one; the intent
questions are fixed text. The step `batch` reads the time of `asked` from the
run's clock, writes `questions.tsv` in one records commit, then posts the
comment on the Intake issue (fencing). A restart that finds `questions.tsv` and
the App's comment of the batch is "done before"; one that finds the record and
no such comment (a run that stopped between the commit and the post) posts the
comment from the record, with no second commit. The step `answers` (row 47a)
compares "a comment after the batch" with that time of `asked`, which comes a
little before the post.

### The answers

`layup run` copies each comment of the Intake issue (`copies.tsv`,
[`records.md`](records.md#nfr-001--the-records-of-start)). An answer is a line
`Q-NNN: <answer>` of a comment whose author's numeric ID is in `approvers.tsv`
in the role of the question's block (the idea owner, or the Operator); any other
comment is recorded as input and answers nothing (`architecture.md` §3).
**Decided here:** the batch's answers are in hand when each block that asks a
question has at least one comment after the batch by its approver that holds at
least one answer line (note 3 of #210: an acknowledgement with no answer line
does not count) (`architecture.md` §5 step 6:
"each in one comment"); until then the step `answers` waits. Then code reads
every such comment after the batch, and the last answer line of a question ID,
in the order of the copies, is its answer (**decided here**: a later line
corrects an earlier one). Each answered question is a row of
`records:intake/answers.tsv`, and its state in `questions.tsv` is `answered`;
a question with no line stays `asked`. **Decided here:** a field of the intent
form parses when its answer is a `decimal`, and `B` and `U` parse only when `B`
is below `U` (so a `B` not below `U` is an unparsed answer of both). When the
three parse, code writes `records:budget.tsv`.

**Decided here** (task `T-hxb9`): the step `answers` copies each comment of the
Intake issue but the App's own ([`run.md`](run.md#copy-before-read)). An answer
comes from a human decision only ([`run.md`](run.md#a-human-decision)): a first
copy (`seen` 1, so an edit answers nothing and stays as input), made by no App,
by an approver in the role of the question's block, whose `created` is strictly
after the batch's `asked`. An answer line is read after a CR at its end is
trimmed: it starts with `Q-` and the question's number, then `: ` exactly, and
the rest is the answer, its trailing spaces kept (and a tab written as a space
by the field rule of the record); an empty answer, and one of exactly `—` (which
reads back as the empty value), answer nothing. The rows of
`intake/answers.tsv` are in the order of the questions. A field of the intent
form is parsed after the spaces around its answer are trimmed, and `B` is
compared with `U` as numbers. A block that waits makes the step `done` with
"waiting" in its detail, the last row of the restart.

### The review of the answers and the follow-up

When the answers are in hand, a session of the role `systems-architect` reads
the answered questions and their answers and writes `result/answers-review.tsv`:
for each, `closed` or `open` with a reason. The rows are
`records:intake/answers-review.tsv`, and each sets the question's state; with no
answered question the session does not start and the step writes no row
(**decided here**, note 3 of #210). Then
code posts **one** follow-up comment on the Intake issue that asks again only
for the questions still `asked`, a field of the intent form that does not parse,
and the `open` ones, in the form of the batch (its two blocks, each question
with its quote, and the rule of the answers); each is a new row of `questions.tsv` with the source
`follow-up`, whose `ref` names the question it asks again. With none, Intake is
done.

**Decided here** (task `T-vwu3`): with no answered question, the step
`answers-review` commits `intake/answers-review.tsv` with its header alone, so a
restart knows that it ran. The review's file is the block `answers-review`
without `session`, its questions the answered questions, each once; its states
go to `questions.tsv` in the session's one records commit. The follow-up asks
again one row per question, whatever the number of rules that name it; a
follow-up row has `kind` `follow-up`, the block of the question it asks again,
no quote of its own, the question's text, `asked` from the run's clock and
`asker` `layup run`; its comment's first line is `LAYUP Intake: the follow-up`,
not the batch's, and each question shows the quote of the question it asks
again. **Decided here** (task `T-4c5c`): each question is in the block of the
question it asks again, so a follow-up of a setup question is in the Operator's
block, shows the line `proposed:` and the quote of its marker when it has one, and
an empty block says `No question for the Operator.` or `No question for the idea
owner.`. The follow-up rows are one records commit before the comment; a restart
that finds the rows and the App's comment of the follow-up is "done before", one
that finds the rows and no comment posts it from the record (as `batch`), and
with nothing to ask again the step writes nothing.

**The answers to the follow-up** (the step `follow-up-answers`): they are in
hand when each block that the follow-up asks has at least one answer comment after it,
by the same rule; until then the step waits. Each answered follow-up is a row of
`intake/answers.tsv` and its state is `answered`; the question that it asks
again keeps its own state. When `budget.tsv` is not written yet, code takes each
of `B`, `U` and the hours from the answer to its follow-up, else from the batch,
and writes `budget.tsv` when the three parse with `B` below `U` (**decided
here**: the follow-up of an intent value asks the same field); a value still
missing or unparsed leaves `budget.tsv` unwritten and sets the state of its
question of the intent form to `open`, and Shape (`M2e`, [`loop.md`](loop.md#the-questions-left-open)) asks it. When
`budget.tsv` is written already and the follow-up asks a field of the intent
form again (the review left it `open`), a parsed answer replaces that value in
`budget.tsv`, its one home, with the follow-up as its `source` (**decided here**,
note 6 of #210). No review runs on the answers to the follow-up, and
Intake is done.
**Decided here** (task `T-tmm2`): the step `follow-up-answers` with no
follow-up row (nothing was asked again) is `done` with "no follow-up; Intake is
done" and writes nothing; it is "done before" when an answer to a follow-up is
in `intake/answers.tsv`; else it copies the comments of the Intake issue (but
the App's own) and reads the answers to the follow-up by the rule of the step
`answers`, from first copies after the follow-up's `asked`. With `budget.tsv`
not yet written, a `B` not below `U` after the follow-up leaves both questions
`open`, whatever the state of the hours. With `budget.tsv` written, a replacement holds only when the three values
still parse with `B` below `U`; else `budget.tsv` stays as it was, and the
questions of the fields asked again are `open`, as is a field whose answer to
the follow-up does not parse. The answers, their states, the questions set
`open` and `budget.tsv` are one records commit, and the step's detail says that
Intake is done.
**Decided here:** one follow-up round in `M2c`; a question still `asked` or
`open` after it stays so in `questions.tsv`, and Shape (`M2e`, [`loop.md`](loop.md#the-questions-left-open)) reads it. Its rows (source `follow-up`, kind `shape`) are not Intake's follow-up: the steps `follow-up` and `follow-up-answers` ignore them (task `T-3f94`).
**Reason:** `architecture.md` §5 step 6 asks again "in one comment", and a loop
with no bound waits on a human with no end.

### The raw fact of the answers

The answers are a raw fact of the target, one fact per question ID
(`architecture.md` §7 step 5). In `M2c` the fact is `records:intake/answers.tsv`
with the copied comments that hold it (`copies/`), byte for byte; the file of
the fact under the target's `docs/facts/` is written by the Scaffold (`M2d`,
`architecture.md` §7 step 8; S04 in phase 2,
[`scaffold.md`](scaffold.md#the-steps-of-the-setup-in-phase-2)), as `docs/facts/` is a rule path that no role
session may edit. **Decided here** (decision 6 of #206). Each question keeps its
quote, so a requirement that later covers an answer's fact can name the quote's
fact (`M2e`).

## The steps of Intake

**Decided here** (decision 11 of #206): when each step of Start is done, the
restart runs the steps of Intake after `probe`, in place of the row `phase`:
`rule-gaps`, `meaning`, `numbering`, `spans`, `draft`, `completeness`, `batch`,
`answers`, `answers-review`, `follow-up`, `follow-up-answers`, `revision`. The record of each
step is: `intake/psb.tsv`; `intake/meaning.tsv`; `spec/spans.tsv` (by `numbering`);
`spec/check.tsv` (by `spans`); `spec/draft.tsv`; `intake/coverage.tsv`;
the comment of the batch (its rows of `questions.tsv`); the rows of
`intake/answers.tsv` of the batch's questions, with the rows of `approvers.tsv`
the batch's point lines named; `intake/answers-review.tsv`; the
follow-up comment (its rows of `questions.tsv`), or, when nothing is left to
ask (no question `asked` or `open`, and the three fields of the intent form that
`budget.tsv` records parse), that state of the records, from which the step posts nothing; the rows of `intake/answers.tsv` of the follow-up's questions, with the rows of `approvers.tsv` its point lines named; the revised `spec/draft.tsv`.
Each step whose record exists — the record of `revision` being the *revised*
`spec/draft.tsv`, told from the session's own draft by the done-before rule of
[`loop.md`](loop.md#the-revision-of-the-draft) — is
`done` with "done before" and runs nothing, but `numbering` and `spans` while the
last `spec/check.tsv` has a `fail` row: a restart runs `numbering` again (The
span check; note 1 of #210); a step whose input is not there yet
(a block of the batch, or of the follow-up, with no answer comment yet, by the
rule of The answers) is `done` with "waiting" in `detail`, and the run releases
the lease and exits 0. A restart goes on from
the records. The block `run-steps` is built, so each build task that adds a step
adds its name to the enum, with its Go schema, and the step to the restart's
order, in one change, as `M2b` did for `probe` ([`run.md`](run.md#the-command));
the first, `rule-gaps`, replaces the row `phase` (the tasks of `M2c`,
[`plan/README.md`](../plan/README.md#the-tasks-of-m2c); task `T-2r8u`). `M2d` adds the step `markers` between `completeness`
and `batch`, `M2e` adds `revision` after `follow-up-answers` (task `T-8rm6`,
[`loop.md`](loop.md#the-revision-of-the-draft)), and the steps of the Scaffold follow it
([`scaffold.md`](scaffold.md#the-steps-of-the-scaffold)).

## NFR-003 — No value without evidence at Intake

Every value that Intake writes has its evidence in the records: each quote by
its offsets in `start/problem-statement.md`; each span by its offsets; each
answer by the comment ID of its copy; each value of `budget.tsv` by the question
ID of its answer. A value of the intent form that does not parse is asked again,
never guessed.

## REQ-015 and REQ-017 — The review of the release of M2c

Their criteria, as for `M2b`
([`session.md`](session.md#req-015-and-req-017--the-review-of-the-release-of-m2b)).
**Decided here** (decision 8 of #206; task `T-2r8u`): a task of its own, before
the demo (row 49a of [the plan](../plan/README.md#the-tasks-of-m2c)), records a
code review of the release of `M2c` (the non-test Go files of `cmd/` and
`internal/`, `go.mod` and the files that the binary embeds, the code that the
demo runs) by a reviewer of a model that wrote none of it, as row 39a did for
`M2b`, with `runs/T-efmy/release-check.sh` run on the release commit, and
adapted in that task where the code of `M2c` needs it (a program that the code
starts, a verb of `git` that reaches a remote). The review reads the whole
release, the code of `M2b` and `M2a` included.

## Input states

**Decided here**, for each state that the happy path never makes:

| State | Result |
| ----- | ------ |
| A routing register with no admitted pair for a role of Intake (`systems-architect` for `meaning` and `answers-review`, `product-owner` for `numbering` and `draft`; `qa-engineer` is the next row's) | the step is `fail` (`pair`), no session starts |
| No admitted pair of `qa-engineer` on a harness other than that of the numbering and the draft, or none of the role at all | `completeness` is `fail` ("not run: no admitted pair on another harness"); no batch |
| A session's output file that is missing, a link, over 1 MiB, or that its block refuses | the step is `fail`, naming the session; nothing of the file is committed |
| A refused start of a session of Intake (`session.md`, The start of a session) | the step is `fail`, naming the session and the refusal; no record (**decided here**, task `T-gmey`: no event holds it, and the row of a refused probe is the probe's own) |
| A quote that occurs nowhere, more than once, or holds a tab | the file of the review is refused (as above) |
| A span table that fails the span check | `spans` is `fail`; `spec/check.tsv` is committed; no batch; a restart runs `numbering` again |
| A draft row whose `covers` names no fact, a `need` or `constraint` fact with no requirement, or a `measure` fact with no success criterion | the draft is refused |
| The cap of milestone 0 reached | the step is `fail` ("the cap of milestone 0"); no session starts |
| A known row of `telemetry.tsv` (`reported` or `computed`) whose `currency` is not `USD` | the step is `fail` ("the cap of milestone 0"), the detail naming the row's session and its currency, as the sum cannot be compared with `MONEY`; no session starts (**decided here**, task `T-ybd5`) |
| A comment with answer lines by an account that is not an approver, or in the other block's role | recorded as input; it answers nothing |
| An output of the revision that its check refuses | the finding is posted on the session, and a new attempt runs with it in the prompt, until the output keeps or the cap of milestone 0 stops the step (task `T-8rm6`) |
| A line `<point>: <login>` in a comment that is no counted answer of the idea owner's block, or whose login the forge does not know | recorded as input; it names no approver (task `T-8rm6`) |
| An answer of the intent form that does not parse (a `B` not below `U` is an unparsed answer of `B` and of `U`) | the field is asked again in the follow-up; still unparsed after it, its question is `open` and `budget.tsv` is not written |
| A question still open after the follow-up | open in `questions.tsv`; Intake is done |

## Not in M2c

The setup questions of steps S01 and S10, the proposed marker sources with their
quotes and their check (`M2d`, decision 9 of #206;
[`scaffold.md`](scaffold.md#req-002--the-setup-questions-in-the-one-batch)), with
the step `markers` between `completeness` and `batch`; the file of the raw fact
under `docs/facts/` (`M2d`);
the first bet, the rendered commit and the confirmed inventory (`M2e`, [`loop.md`](loop.md#the-revision-of-the-draft)); the
owner map (phase 3); the smart-if provider and its authority (`M3a`); the
actions of the band `B` and `U` and the accepted bounds of an unknown part
(`M3a`, `M3b`); the caps of later milestones (`M2e`, [`loop.md`](loop.md#req-012--the-bet)); the pilot baseline block of
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
| The draft | unit | The draft is refused when `covers` names no fact, a `need` or `constraint` fact is not covered, or a `measure` fact has no success criterion, and kept otherwise. |
| The completeness review | integration | With two admitted harnesses, the review runs on the first admitted pair of `qa-engineer` whose harness is not that of the numbering and the draft, even when the role's first pair is on that harness; with none, the step is `fail` and no batch is posted. |
| The merge and the IDs | unit | A question whose quote and kind repeat is dropped; each kept question gets one ID, in the order of the merge. |
| The comment and the answers | integration | With the httptest forge: one comment with two blocks; the step waits until each block that asks has a comment by its approver; answers by an approver in its role are kept, others recorded as input; a later line of an ID corrects an earlier one; a missing, unparsed or `open` answer goes to the follow-up; the answers and `budget.tsv` are committed. |
| The steps of Intake | e2e | In CI with no secret, the binary, a local fake forge and scripted harness programs: a restart runs the steps to `batch` and exits "waiting"; after a scripted answer comment that leaves one question unanswered, the next restart reads it, runs the review, posts one follow-up for that question and exits "waiting"; after a scripted answer to the follow-up, the next restart records it and ends Intake. |
| The review of the release | uat | A recorded code review of the release for `REQ-015` and `REQ-017`, by a reviewer of a model that wrote none of it, with `release-check.sh` on the release commit, before the demo (row 49a of [the plan](../plan/README.md#the-tasks-of-m2c)). |
| The demo | uat | On a pilot problem statement that the idea owner chooses, on real harnesses: one Intake comment holds the rule gaps and the gaps of meaning, and `layup spec check --facts` passes on the numbered spans. |
