# `layup setup`, `layup setup verify` and the stack catalog

Conventions: [`README.md`](README.md). The order of the setup is
[`architecture.md`](../architecture.md) §5 (Start, the gap check, Scaffold);
the gates are §6.

## REQ-002 — Set up a target, and prove the setup

Requirement: "`layup setup` creates a target repository from [the pinned
baseline] for the target's domain and stack, stops at each human decision,
writes the answers to Git, and `layup setup verify` proves the setup with
evidence for every value." Derives from `architecture.md` §5 and §6, ADR-0011
decisions 6 and 7, ADR-0016 and ADR-0017.

### The boundary of phase 1

In phase 1, `layup setup` is a **step runner**: it does each step of a target's
setup that needs no judgement, on the LAYUP host, with `git` as its only
external program. It makes no forge call and starts no role session (decision
D1 of the plan of #74). What a step needs from a human or a session is an
input file; what it needs on the forge is a command that it prints for the
Operator.

| Part of §3 and §5 | In phase 1 | Who, and when not in phase 1 |
| ----------------- | ---------- | ---------------------------- |
| Start 1: the empty repository on the forge, the App installed | no | the Operator, by hand, as today |
| Start 2: resolve the baseline's latest commit, clone it, remove `.git`, record the pin | **yes** (S02) | — |
| Start 2: the push of the root commit | **yes**, as a printed command (S03) | the Operator runs it |
| Start 3: read back the default branch, its root tree and visibility; the plan check; the six capabilities | no | `layup run`, phase 2 |
| Start 3: the first records commit, `approvers.tsv`, the lease row | no | `layup run`, phase 2 |
| Start 3: the Intake issue, the control issue, the dead-man job's first notice | no | `layup run`, phase 2 |
| Gap check 1: `layup psb check` | **yes** ([`psb-check.md`](psb-check.md)) | — |
| Gap check 2 to 4: the review of meaning, the specification sessions, the proposed marker sources | no | sessions of `layup run`, phase 2 |
| Gap check 5 and 6: one comment on the Intake issue; the answers copied; the follow-up | no; the answers are the input file `answers.tsv` | `layup run`, phase 2 |
| Scaffold 1: the rows that need no judgement | **yes** | — |
| Scaffold 2: the prose rows S07, S08, S09, S14, by a role session | no; the text is an input file, and `layup setup` runs the row's check | a session of `layup run`, phase 2 |
| Scaffold 3 and 4: the facts, the stack gates | **yes** (S06, S12) | — |
| Scaffold 5: `layup setup verify` | **yes** | — |
| Scaffold 6: the push of the setup commits; the rulesets | **yes**, as printed commands (S13) | the Operator runs them |
| Scaffold 7: the pushed tree equals the verified one; the rules read back; the probes of §3 | no | `layup run`, phase 2 |
| The pilot baseline | no | `layup run`, phase 4 |

### The command `layup setup`

```text
layup setup WORK
```

- `WORK`: the work area of one target on the LAYUP host. The Operator makes it
  and puts the inputs in it; `layup setup` writes the rest.
- No flag. Every value comes from a file of `WORK` (an argument would be a value
  with no record).
- Exit codes: 0 every step of phase 1 is done or handed to the Operator; 1 a
  step's check failed; 2 a usage or input error (an input that does not match
  its schema, or an input that changed after a step read it); 3 the run stopped for human input (the table lists every missing
  input of that step).
- A run is **resumable**: it reads `WORK/out/record.tsv` and skips each step
  that has a `done` row there. So a run after a stop goes on from the stopped
  step, and the pin is resolved once (§5 Start 2: "no run resolves the commit
  again").

**Decided here** (task `T-79y7`, #85), the rules of the step runner:

- **The run order** (K16, O-123): S01 to S06; then the prose step, S07, S08,
  S09 and S14, as one group whose missing inputs make one stop table; then S10,
  S11, S12, S13 and S15. The step table lists S01 to S15 in their order; the
  progress lines (`layup setup: [<i>/<n>] <step>`) follow the run order.
- **The outcome of a step:** `done` (the runner writes its value rows and its
  `done` row with the evidence, and goes on), a stop (the run ends with the
  stop table, exit 3), `fail` or `not-active` (the run ends, exit 1), or a
  hand-off (S13: its own part is done, and its commands are in `commands.sh`;
  its `done` row says `handed to the Operator: …`, never that the Operator ran
  a command, and its result is `operator`). A hand-off with no command is
  `fail`: the runner asks the step for its commands with the record and the
  step's value rows, before its `done` row. After a step that did not pass,
  each later step is `not-active`, evidence `not run: <step> did not pass`.
  When the commit of a step fails, the step is `fail`, with no `done` row, so
  a rerun does it again. Before the commit, the runner checks that the target
  is on the branch `layup-setup`; when it is not, the runner makes no commit,
  and the step is `fail`. The actor of each step of phase 1 is `layup-setup`.
- **The inputs:** a `WORK` that is not a directory, or an `answers.tsv` or a
  `record.tsv` that does not match its schema, is an input error (exit 2). A
  missing `out/record.tsv` is a new work area, and a missing
  `inputs/answers.tsv` is no answer; what a step does with no answer is the
  step's rule (S01 stops for its questions).
- **`commands.sh`** is written again at the end of each run that gives a table
  (exit 0, 1 or 3), from the commands of the done steps, each a function of the
  record, in the order of
  [Where the records go in phase 1](#where-the-records-go-in-phase-1), each
  with a comment line before it, and empty when no done step has a command, so
  no command of an earlier record stays; a run that ends with exit 2 does not
  write it. The texts of the commands, the apply of the ruleset included, are
  the steps' (S03, S13, S15).
- **The present code** (task `T-79y7`): each step is a stub, `not-active`,
  evidence `not built yet`, so `layup setup WORK` gives exit 1 until the rows 9,
  13 and 15 of the [plan](../plan/README.md#the-tasks-of-phase-1) build the
  steps (`NFR-004`).

**The work area** (decided here: §1 names the work area on the host, and no section gives its layout; one directory per target keeps a run resumable from files alone):

| Path | What | Written by |
| ---- | ---- | ---------- |
| `WORK/inputs/answers.tsv` | the answers, by question ID ([below](#the-answers)) | the Operator (phase 1); `layup run` from the Intake comments (phase 2) |
| `WORK/inputs/briefs/problem-statement.md` | the problem statement | the Operator |
| `WORK/inputs/briefs/vision.md` | the vision brief, when there is one | the Operator |
| `WORK/inputs/files/<path>` | the text of a prose row, at its path in the target, for example `WORK/inputs/files/docs/glossary.md` | a human or an agent (phase 1); a role session (phase 2) |
| `WORK/target/` | the baseline copy: a Git repository with the root commit on `main` and the setup commits on the branch `layup-setup` | `layup setup` |
| `WORK/out/record.tsv` | the setup record ([below](#the-setup-record)) | `layup setup` |
| `WORK/out/rule-paths.tsv` | the rule-path register ([below](#the-rule-path-register)) | `layup setup` |
| `WORK/out/verify.tsv` | the table of `layup setup verify`, which S15 commits | the Operator, by redirecting its output |
| `WORK/out/ruleset-default.json` | the ruleset of the default branch ([S13](#the-steps)) | `layup setup` |
| `WORK/out/commands.sh` | each command for the Operator, in order, with a comment line before each | `layup setup` |

### The steps

One row per step of LAYUP's [`steps.tsv`](../setup/steps.tsv). That file is the
record of LAYUP's own setup, and it does not change for a target; this table is
the one home of the phase-1 steps of a target. A row that differs from
`steps.tsv` names the decision that changes it.

| Step | Actor in phase 1 | Inputs | What it does | Output | Evidence (the `done` row) |
| ---- | ---------------- | ------ | ------------ | ------ | ------------------------- |
| S01 | `layup setup`; stops for the Operator and the idea owner | `answers.tsv`; the problem statement | Reads the answers of the questions `S01-stack`, `S01-name`, `S01-visibility`, `S01-baseline` (the baseline's repository URL), and runs `layup psb check` on `inputs/briefs/problem-statement.md`: each gap is a question whose ID is the gap's `id` (`Q-NNN`; [`psb-check.md`](psb-check.md#the-table)), for the idea owner. A missing answer stops the run, with all missing ones, of both kinds, in one table (Decision Point 2). `internal/cli` runs the rules of `internal/psb` and hands the gap table to `internal/setup`, as it does for the step checks. A record row `brief.sha256` (`computed`) holds the SHA-256 of the problem statement that the rules read. In phase 1 the name and the visibility are answers, because Start is not in phase 1; §5 (gap check 4) asks the Operator only the stack at S01. No gate-mode question: **decided here**, as §5 names none, and the target runs its own full gate (§8). | record rows `stack`, `name`, `visibility`, `baseline`, `brief.sha256` | every answer present; the stack has a catalog entry |
| S02 | `layup setup` | `S01-baseline` | `git ls-remote <url> HEAD` gives the commit; `git clone` and `git checkout` of it; `git rev-parse <commit>^{tree}` gives the tree; removes `.git`. The clone reads no credential of the host and starts no `ssh`, so the baseline's repository is public (known limit [L-A7](../architecture.md#15-known-limits)). Changes `steps.tsv` S02 (`npx degit`): ADR-0011 decision 7, §5 Start 2. | the copy; record rows `pin.source`, `pin.commit`, `pin.tree`, `pin.time` | the commit and the tree |
| S03 | `layup setup`; the push by the Operator | the copy | `git init`, one commit of the unmodified copy on `main` (message `chore: the unmodified baseline at <commit>`), checks that its tree equals `pin.tree`, writes the push command to `commands.sh`. Creating the remote repository is the Operator's (Start 1). | the root commit; a command | root tree = `pin.tree` |
| S04 | `layup setup` | the pin rows | On the branch `layup-setup` from the root commit: writes `docs/setup/armature.pin` from the pin rows ([`NFR-006`](#nfr-006--the-baseline-at-a-pinned-recorded-version)). Installing the hooks is a setting of a clone, not of the tree: not done (**decided here**, as §5 names no hook for a target; the Operator's clone has none, §5 Scaffold 6). The decision record of the pin is the baseline's own ADR form, written from a fixed text with the pin values. | the pin file; the ADR | check `pin` of `layup setup verify` |
| S05 | `layup setup` | the copy | Deletes the baseline's own history: the paths that check `kit-history` reads (`docs/decisions/`, `docs/audit/`, each `docs/tasks/T-*.md`, their lines in `backlog.md` and `completed.md`). A link that the deletion breaks is a missing input: the run stops and lists each one, and the fixed text of that file comes as `inputs/files/<path>`. | the commit | checks `kit-history` and `link-lint` |
| S06 | `layup setup` | the briefs; `answers.tsv` | Refuses a problem statement whose SHA-256 differs from the row `brief.sha256` of S01 (exit 2), so the `Q-NNN` IDs stay true; the Operator restores the file, or starts again in a new work area. Copies each brief byte for byte into `docs/facts/`; writes the `S01-` and `Q-` answers of `answers.tsv` as a raw fact record, one fact per question ID, each with the question ID, the question text and the answer; writes `docs/setup/facts.sha256` and the index rows. The numbered facts of the problem statement come with the first bet (§7), not here. | the facts | check `facts` |
| S07 | `layup setup`; the text is an input | `inputs/files/docs/onboarding-for-engineers.md` | Copies the file into the tree; a missing file stops the run. | the file | check `onboarding` |
| S08 | the same | `inputs/files/docs/glossary.md` | the same | the file | check `glossary` |
| S09 | the same | `inputs/files/docs/guardrails.md` | the same | the file | check `guardrails` |
| S10 | `layup setup`; stops for the Operator | the tree; `answers.tsv` | Lists every marker of the tree outside the exemptions of LAYUP's `MK_EXEMPT` ([`setup-check.sh`](../setup/setup-check.sh); the baseline has no setup check, §5), which the engine embeds at its version. A path of that pattern that a target does not have matches nothing. A marker with no answer row stops the run; the table lists all of them at once (§5 gap check, "one batch"). | — | every marker has an answer row |
| S11 | `layup setup` | the answers of S10 | Replaces each marker whose answer has a value with that value, and writes its record row with the source. A marker whose answer is `gap` keeps its marker and gets a row in `docs/setup/open-gaps.tsv` with the answer's question (Invariant 4). When S10 listed at least one marker, writes the answer of each marker that S10 listed as a second raw fact record, in the same form, with its own index row and its own line in `facts.sha256`; with no marker, it writes no second record; the record of S06 does not change (a raw facts record is immutable). | the tree; the second answers record; record rows `marker:<file>:<line>` | checks `markers`, `sources` and `facts` |
| S12 | `layup setup` | the catalog entry of the stack | Writes the files of the entry (for Go: `go.mod` with the module path from `name`, the tools' configuration), `docs/gates.tsv`, and one CI job per gate kind, whose id and name are the kind (K30). For each gap of the entry (the coverage floor of Go, K24), writes its row of `docs/setup/open-gaps.tsv` and its record row. The baseline's own workflows stay byte for byte (`REQ-018`). Adds no `setup-check` job. Changes `steps.tsv` S12: ADR-0011 decision 7, ADR-0016, §5 Scaffold 4. | the gate files; the gap rows of the entry | checks `jobs` and `gates` (each row `gate:<kind>`) |
| S13 | `layup setup`; applied by the Operator | the job names | Writes `docs/setup/branch-protection.json` (in the form of LAYUP's own file of that name; the baseline has no `docs/setup/`) and `WORK/out/ruleset-default.json`: the default branch and the ref `layup-probe`; a pull request required; each gate job a required check, pinned to GitHub Actions; no force push, no deletion; an empty bypass list. In phase 1 it requires no `layup/` check, because no phase-1 command posts one ([`records.md`](records.md#nfr-002--a-target-is-independent-of-layup)). Writes to `commands.sh` the push of `layup-setup` onto the default branch (`git push origin layup-setup:main`, a fast-forward from the root commit; §5 Scaffold 6: "pushes the setup commits on top of the root commit") and, after it, the apply command of the ruleset. | the ruleset file; commands | the ruleset file, and its commands in `commands.sh` for the Operator (a hand-off: phase 1 does not see the Operator's run) |
| S14 | `layup setup`; the text is an input | `inputs/files/README.md`, `inputs/files/AGENTS.md` | Copies the files; a missing file stops the run. | the files | check `identity` |
| S15 | `layup setup`; the push by the Operator | the record; `out/verify.tsv` | Writes the record's last rows and the rule-path register, and the first commit of `layup-records` ([below](#where-the-records-go-in-phase-1), O-115). Not `steps.tsv` into the target (§5: LAYUP's own file). | `out/record.tsv`, `out/rule-paths.tsv`; the records commit; a command | every row of `verify.tsv` is `pass` or `clear` |

Each step from S04 to S14 that changes the tree is one commit on `layup-setup`,
with the message `chore: setup <step>`, so the history shows each step (the
runner first checks that the target is on that branch); S03 makes
the root commit on `main` and S15 the first commit of `layup-records`, each with
its own message. **Decided by the Operator** (O-136, #85): the author and the
committer of each setup commit and of the records commit is
`layup-agent[bot] <335371832+layup-agent[bot]@users.noreply.github.com>`, the
LAYUP App's bot, and the date of each is `pin.time` of S02, so one input gives
one commit ID (`NFR-005`).

### The stop table

When a step stops (exit 3), `layup setup` prints every missing input of that
step, and nothing else:

```tsv-schema setup-stop stdout
step      id(SNN)    -    the step that stopped
question  text       key  the question ID: `S01-<name>`; `Q-NNN` for a gap of `layup psb check`; `M-<x8>` for a marker; `F-<path>` for an input file of the target; `O-<name>` for an output the Operator makes
ask       text       -    the question in words
where     text       -    for a marker, `<file>:<line> <marker>`; for a file, its path in the target; for a gap `Q-NNN`, `inputs/briefs/problem-statement.md:<line>`, `—` for line 0 (K33); `—` otherwise
```

**Decided here:** the ID of a marker's question is `M-` and the first 8
hexadecimal characters of the SHA-256 of `<file>`, a tab, and the marker text,
so it stays the same while the marker stays. **Decided here** (task `T-79y7`,
#85): the rows are in a fixed order: the `S01-` rows in the order of S01 (stack,
name, visibility, baseline), the `Q-` rows by the number of their ID (`Q-999`
before `Q-1000`), the `F-` rows by path, the `M-` rows by file and line, and the
other rows last.

### The step table

When a run does not stop, `layup setup` prints one row per step:

```tsv-schema setup-steps stdout
step      id(SNN)                         key  S01 to S15
actor     enum(layup-setup|operator|layup-run)  -  who does the step in phase 1; `layup-run` for a part of a later phase
result    enum(done|fail|not-active|operator)  -  `done`: done and checked; `fail`: its check failed; `not-active`: its check could not run (exit 1); `operator`: a command in `commands.sh` waits for the Operator
evidence  text                            -    the evidence line of the step table above, or the failed check's reason
```

### The answers

Each row of `answers.tsv` answers a question that a step of this run asked: an
`S01-` question, a gap `Q-NNN` of S01, or a marker `M-<x8>` that S10 listed.
Any other row (a stale marker of an earlier baseline, an `F-` or `O-` ID, whose
answer is a file or a command) is an input error: the run stops with exit 2 and
names the row, so no answer becomes a fact of the target without a question
(**decided here**, Invariant 4). The runner checks the prefixes and the `S01-`
rows, whose four names are fixed, before any step; S01 checks the `Q-` rows
against its gap table, and S10 the `M-` rows against its markers.

**Decided here** (task `T-79y7`, #85): before any step, the runner refuses a
row whose question no step asks (`F-`, `O-`, another prefix, or an `S01-`
question that is not one of the four of S01), an answer `gap` to a question
that is not a marker, and an answer `gap` with no `question_text`; a question
twice is an error of the key. When a step that reads answers is done (S01 its
`S01-` and `Q-` rows, S10 its `M-` rows), the runner writes the record row
`<step> answers.sha256` (source `computed`, ref `sha256 inputs/answers.tsv
<prefix>…`), the SHA-256 of those rows as the step read them; each run first
compares it with the rows as they are, so a row that changed, went or came is
exit 2, "an input that changed after a step read it", before S04 or S06 writes
an answer as a fact. A done step that reads answers and has no such row is exit
2 too.

```tsv-schema setup-answers host:<work>/inputs/answers.tsv
question  text                       key  the question ID, as the stop table gives it
answer    text                       -    the value; `gap` keeps a marker as an open gap
by        enum(operator|idea-owner)  -    who answered
source    text                       -    the evidence: the URL of the comment that holds the answer, or a fact citation `F-NNNN#n` that the Operator accepted
question_text  text                  -    for `gap`: the question that goes into `open-gaps.tsv`; `—` otherwise
```

### The setup record

Each value that the setup sets, with its source (`NFR-003`), and one `done` row
per finished step.

```tsv-schema setup-record records:setup/record.tsv
step    id(SNN)                                   key  the step that set the value
name    text                                      key  the value's name: `stack`, `pin.commit`, `marker:<file>:<line>`, …; `done` for a step's evidence row
value   text                                      -    the value; for `done`, the evidence line
source  enum(answer|catalog|fact|computed|gap|step)  -  (`computed` is decided here: §5 names three sources, and the pin values come from `git`, not from a person) `answer`: an answer row; `catalog`: a catalog file; `fact`: a fact citation the Operator accepted; `computed`: a command's output, or a hash that the engine computes; `gap`: kept as an open gap; `step`: a `done` row
ref     text                                      -    `answer`: the question ID; `catalog`: `<stack>/<path>`; `fact`: `F-NNNN#n`; `computed`: the command, `sha256 <path>` for the hash of a file, or `sha256 <path> <prefix>…` for the hash of the rows of the file whose question has one of the prefixes (`answers.sha256`); `gap`: `docs/setup/open-gaps.tsv`; `step`: `—`
```

### The rule-path register

Written at setup from the catalog and the baseline (§6, Rule protection); read
by `layup/rules` in phase 2 (`REQ-003`).

```tsv-schema rule-paths records:rule-paths.tsv
pattern    text       key  a path, or a directory ending with `/`
exception  text       -    the part that is not a rule path; for `docs/guardrails.md`: `added lines in section 2`; `—` otherwise
source     enum(baseline|catalog|architecture)  -  where the entry comes from
```

The entries of phase 1 (**decided here** from the baseline at LAYUP's pin, by §6's
list): `.github/`, `.githooks/`, `.gitattributes`, `AGENTS.md`, `CLAUDE.md`,
`docs/engineering-discipline.md`, `docs/issue-workflow.md`, `docs/ci/`,
`docs/tests/`, each file of the tree whose name ends with `.sh`, `docs/gates.tsv`,
each `config` path of the manifest, `docs/setup/`, `docs/facts/`, and
`docs/guardrails.md` with its exception. A baseline whose rule files differ
is known limit L-A6 of the architecture. **Known limit of phase 1:** a new rule
document that a newer baseline adds, and that this list does not name, breaks
no step and fails no check, so it is not in the register; a LAYUP change that
follows the baseline adds it.

### The files of `docs/setup/` in a target

The baseline has no `docs/setup/`; S04, S06, S11 and S13 write it in the form of
LAYUP's own files of the same names, so that the same check rules read them.

```tsv-schema open-gaps target:docs/setup/open-gaps.tsv no-header
file      path  key  the file that holds the marker
marker    text  key  the marker, with its two angle quotes
question  text  -    the question that its value needs; never `—`
```

`docs/setup/facts.sha256` is not a table: one line per raw facts file (the key is the path), in the
form of `sha256sum`: the SHA-256, two spaces, and the path from the root, for
example `<64 hex>  docs/facts/problem-statement-brief.md`. Check `facts` reads
it as `check_facts` does.

### Where the records go in phase 1

The Operator's decision O-115 (comment 5928163366 on #74): `layup setup` writes
the records as files and prints one command; the Operator runs it, as for the
root commit. The order of a setup in phase 1:

1. `layup setup WORK` runs S01 to S14, with a stop (exit 3) at each step that
   needs an input, until S15.
2. S15 needs `WORK/out/verify.tsv`. When it is missing, the run stops with the
   question `O-verify`: run `layup setup verify WORK >
   WORK/out/verify.tsv`.
3. The next run of `layup setup WORK` does S15: it refuses a `verify.tsv` with a
   row that is not `pass` or `clear` (exit 1); otherwise it makes, in
   `WORK/target`, the first commit of the orphan branch `layup-records` with
   `README.md` (a fixed text: what the branch is, that only `layup run` writes
   it from Start on), `setup/record.tsv`, `setup/verify.tsv` and
   `rule-paths.tsv`, and adds its push to `commands.sh`.
4. The Operator runs `commands.sh`: the push of the root commit to `main`
   (S03), the push of `layup-setup` onto `main` (S13), the push of
   `layup-records` (S15), and the apply of the default branch's ruleset (S13),
   in that order, so the setup commits are on the default branch before the
   ruleset requires a pull request, with the Operator's own login,
   from a clone with no hooks installed.

### The checks of `layup setup verify`

```text
layup setup verify WORK
```

No flag. It reads `WORK/target` at the head of `layup-setup`, `WORK/out/record.tsv` and `WORK/inputs/answers.tsv`,
in a scratch work tree. It changes no ref of `WORK/target` and no file of `WORK/out` or `WORK/inputs`; the commit of a fixture run is an object on no ref. Exit codes as in
[`README.md`](README.md#commands): 0 when each row is `pass` or `clear`.

Each check of `setup-check.sh` below has the rule of the function of the same
name in LAYUP's [`setup-check.sh`](../setup/setup-check.sh), run from outside on
the target, with the target's values in place of LAYUP's. The engine's version
of each such check passes and fails on the same fixtures as
`docs/setup/tests/run.sh` (the defence that ADR-0011's Consequences name).

| Check | In phase 1 | The rule for a target |
| ----- | ---------- | --------------------- |
| `discipline-tests` | yes | the baseline's own `sh docs/tests/run-discipline-tests.sh` exits 0, run as a gate command (below, task `T-8vpw`) |
| `pin` | yes | `check_pin`; and the pin file's `source`, `commit` and `tree` equal the record rows `pin.source`, `pin.commit` and `pin.tree` of S02, its `date` is the date of `pin.time`, and `method` equals the text of [`NFR-006`](#nfr-006--the-baseline-at-a-pinned-recorded-version) with the source and the commit |
| `kit-history` | yes | `check_kit_history`, whose repository is the record row `pin.source` with no scheme and no `/` or `.git` at its end (for LAYUP's baseline, `github.com/pharzam/armature`); a task index links it when it holds that text with the start of the file or a separator before it, and after it, with or without `.git` or `.` in between, the end of the file or a separator; a separator is a character that is not a letter or a digit of any script, `.`, `-` or `_` (review rounds 1 to 4 of #84) |
| `facts` | yes, in a target's form | each brief is in `docs/facts/` and in `facts.sha256` with its hash; each question ID of `answers.tsv` is a fact in exactly one answers record: the `S01-` and `Q-` IDs in the record of S04 (O-124), the `M-` IDs in the record of S11 (which exists only when S10 listed a marker); the index has a row per record; what exists is read, and the `done` rows decide what must exist (below, task `T-9t1q`). LAYUP's counts (39, 75, 19) are LAYUP's and do not apply. |
| `onboarding` | yes, in a target's form | the file exists, holds no marker, links the problem statement (`](facts/problem-statement-brief.md)`); each `F-NNNN#n` it cites is a fact of the one record `F-NNNN-*.md` in `docs/facts/` |
| `glossary` | yes, in a target's form | each `F-NNNN#n` it cites resolves as above; LAYUP's heading and its count of 25 do not apply |
| `guardrails` | yes, in a target's form | each entry (a bullet `- **Inv-N**`) has a `Check:` value that is `no check yet` or a file and a gate, as `check_guardrails`; each citation resolves; LAYUP's heading and its count of 9 do not apply |
| `markers` | yes | `check_markers`, with LAYUP's `MK_EXEMPT` (S10) |
| `adapted` | yes | `check_adapted`, with LAYUP's `AD_EXCLUDE` and `ad_allowed`, which the engine embeds at its version; a path that a target does not have matches nothing |
| `identity` | yes | `check_identity`, with the target's name: `README.md` also holds the record row `name` of S01 |
| `link-lint` | yes | the baseline's own `sh docs/links/link-lint.sh` exits 0, run as a gate command (below, task `T-8vpw`) |
| `sources` | yes | every value row of the record has a source; each `answer` ref is a row of `answers.tsv`; each `catalog` ref is a file of the catalog entry; each `fact` ref is a fact of `docs/facts/`; each `gap` row has its marker in the tree and its row in `open-gaps.tsv`; each `computed` ref is not empty, and a hash names a file (below, task `T-8vpw`) |
| `jobs` | yes | each kind of `docs/gates.tsv` has a CI job with the kind's name (the evidence of S12) |
| `gate:<kind>`, one per kind | yes | an `active` kind: `layup gate` with `--base` and `--head` the setup head gives `pass` or `clear`, and with `--head` a commit of the kind's known-bad fixture applied (`git apply`) on the setup head gives `fail`; then `pass`. The first rule that matches decides: the clean run `fail` gives `fail`; either run `not-active` gives `not-active`; a fixture run that is not `fail` gives `fail`, reason `fixture not detected`. A `pending` kind: `clear`, reason `pending: fixture not run` (§6: "recorded as not run, never as a detection"). |
| `ci`, `procedure`, `protection` | not a check for a target | they read LAYUP's own CI and `steps.tsv`, and the classic protection that the rulesets replace (§5) |
| the rulesets read back (S13); the probes of §3 | no | `layup run`, phase 2: not in the list, so a correct setup exits 0 |

**Decided here:** the target's form of `facts`, `onboarding`, `glossary` and
`guardrails`. LAYUP's functions count LAYUP's own facts (39 facts of `F-0001`,
25 glossary rows, 9 invariants); a target's counts come from its own problem
statement, and its numbered facts come only with the first bet (§7). So in phase
1 these checks hold the rules that do not depend on LAYUP's numbers. The
prose of S07 to S09 cites the raw problem statement by its path until the
numbered facts exist.

```tsv-schema setup-verify stdout
check   text                             key  the check name of the table above
result  enum(pass|fail|not-active|clear)  -   `not-active` when the check could not run (for example a missing tool or script)
reason  text                             -    the first failure, or the `clear` reason; `—` for `pass`
```

**Decided here** (task `T-6x75`, #84):

- **The rows** are in the order of the table above, then one row
  `gate:<kind>` per kind of `docs/gates.tsv` at the head of `layup-setup`, in
  the order of the manifest. A missing or malformed manifest, or one with no
  row, is an input error (exit 2), as for `layup gate`
  ([`gate.md`](gate.md#nfr-004--a-check-that-is-not-active-is-not-passed) item
  4); so are `git` older than 2.32, an `answers.tsv` or a `record.tsv` that is
  missing or does not match its schema, and no commit at the branch
  `layup-setup`. Reason: the run cannot name its rows, or cannot read its input.
- **A check that this version of `layup` does not have yet** is `not-active`,
  reason `not built yet`, so the command gives exit 1 until rows 10 to 15 of the
  [plan](../plan/README.md#the-tasks-of-phase-1) add each check (`NFR-004`
  item 1). The present code has `discipline-tests`, `link-lint`, `markers` and
  `sources` (task `T-8vpw`), `pin`, `kit-history`, `facts`, `onboarding`,
  `glossary` and `guardrails` (task `T-9t1q`), `adapted` (task `T-8ya0`) and
  `identity`.
- **The scratch tree** is `git worktree add --detach` of the head of
  `layup-setup`, in a new temporary directory outside `WORK`: a temporary
  directory (`TMPDIR`) in `WORK` is an input error. The run adds the tree in the
  step of the first built check, and removes it in the step of the last row, so
  the progress lines cover both; a call with no built check makes no tree. A
  tree that the run cannot add makes each row of a built check `not-active`,
  reason `scratch tree: add failed`; a tree or a directory that it cannot remove
  gives the whole table, the path on standard error and exit 2, as for
  `layup gate` (review round 1 of #84, findings 2, 3 and 5).
- **A row** is `fail` with the first finding of its check as its reason. A
  check reads the scratch tree, and the record as it stands: a record row that
  it needs and does not find, or finds with no value, is the finding
  `record: no value at <step> <name>`.
- **The one-check call** that `internal/cli` makes as the evidence of a step
  (S04 to S14) runs the named checks in the order of the table; a name that is
  not a check is an error, never a row. The name `gates` gives each row
  `gate:<kind>`, and is the only name that reads the manifest (K13: the evidence
  "checks `jobs` and `gates`" of S12). `layup setup verify` takes no check name:
  "`layup setup verify <check>` OK" of `architecture.md` §5 names a check, not
  a form of the command.
- **The fixtures.** The core of a check is the rule of its sh function on a
  repository root; it gives each finding as the text after `FAIL ` of the sh
  line, in the order of the function. `TestTheFixturesOfSetupCheck` of
  `internal/verify` builds each case of `docs/setup/tests/` as `run.sh` does;
  for each built check that `EXPECT` has lines of, the lines of the core are
  those lines, as a set; in a case of `frame`, a built check that `EXPECT` has
  no line of gives `OK`, because each overlay of `frame` holds a good setup for
  the checks that it runs (finding 4). Each group of `docs/setup/tests/` is in
  one list: built, not built yet, not a check for a target (`ci`, `procedure`,
  `protection`), or `frame`. **Known limit:** CI gives the sh runner the
  fixtures of the default branch and the Go test those of the pull request, so a
  change of a fixture reaches the sh runner only after the merge.
- **The tests build their baseline at test time** (K11): `internal/standin`,
  which only test files import, writes a stand-in baseline with its own task
  history and the phrases that check `identity` refuses, gives its `file://`
  URL (a valid `S01-baseline` in a test: `git ls-remote` and `git clone` take
  it), and makes a work area from it with a setup by hand. No file of it is in
  Git; a marker in Go source is written as an escape.

**Decided here** (task `T-8vpw`, #87), the checks `markers`, `sources`,
`discipline-tests` and `link-lint`, and the calls that S05 and S10 read:

- **One scanner of markers** for check `markers` and for S10: the files of
  `git ls-files -z` outside LAYUP's `MK_EXEMPT` (embedded at the engine's
  version; `TestTheExemptionsOfMarkersEqualTheSh` compares it with
  `setup-check.sh`), each a regular file; on a line, each open angle quote
  starts a marker to the first close quote after it, or to the line end; an
  open quote between two backticks is the mention and is skipped; the text of
  the convention (the two quotes around an ellipsis) is not a marker. The
  scanner gives each occurrence with its file, its line and its text; check
  `markers` keys it by file and text, as `check_markers` does, and check
  `sources` reads the line. S10 (row 13) gets the list through `internal/cli`.
- **The open gaps:** `docs/setup/open-gaps.tsv` is read by tabs, by the columns
  of its block, as `check_markers` reads it: a missing column is empty, and an
  absent file is no rows. Each marker needs its row (`unlisted:`), each row its
  marker (`stale:`), and each row a question that is not empty, the empty mark
  or only blanks (`question:`), with the lines of `check_markers`. The Go schema
  of the block is in `internal/work`, for the block test and the writer of S11.
- **#21 in `check_markers`:** it lists the files with `git ls-files -z` (a name
  with `"`, `\` or a control character is no longer skipped; the fixture
  `markers/bad-quoted-name` has a DEL character in its name, which git quotes
  and a Windows checkout accepts), and it refuses a blank question or the empty
  mark (`markers/bad-blank-question`). Its `awk` gets each name through the
  environment, as `-v` reads a `\` in it as an escape, and runs in the C
  locale, which reads bytes as the engine does (macOS `awk` in a UTF-8 locale
  misses a marker at the end of a line). Round 1 of #87 found both;
  `TestTheShAndTheGoFormOfMarkersAgree` runs the two forms on a scratch
  repository in two locales, as no fixture can hold a name with `\` (a Windows
  checkout refuses it). **Known limit of the sh function:** a name with a line
  feed is still split by `tr`. **Known limit of both forms** (note 1 of round 2
  of #87): a marker in a file whose name holds a tab cannot be listed in
  `open-gaps.tsv`, as a cell of a TSV file holds no tab, so it stays
  `unlisted:`.
- **Check `sources`:** each row of the setup record but a `done` row, by its
  source: `answer`, the ref is a row of `answers.tsv`; `catalog`, the entry of
  the row `S01 stack` has the file of the ref (in the entry of the binary,
  task `T-c06a`); `fact`, the ref resolves;
  `gap`, the row `marker:<file>:<line>` names a line that holds the marker of
  its value, and `open-gaps.tsv` has the row; `computed`, the ref is not empty,
  and a ref `sha256 <path>` (with or without prefixes) names a file of the
  tree or of the work area (the step that wrote a hash checks its value);
  `step`, a value row with it fails. Each finding is
  `source: <step> <name>: <reason>`.
- **The two scripts:** `sh docs/tests/run-discipline-tests.sh` and
  `sh docs/links/link-lint.sh` of the target run in the scratch tree, with no
  argument and with the environment of `layup`, by the rule of a gate command
  (`gate.md`): exit 0 is `pass`; another exit is `fail`, reason `exit <code>`;
  a signal is `fail`, reason `signal <name>`; a script that is not a file of the
  tree is `not-active`, reason `missing: <path>`; an `sh` that is not found is
  `not-active`, reason `tool not found: sh`, and one that is found and then
  does not start is `fail`, reason `exit -1`, as for a gate command. The output
  of a script goes to standard error, never into the table; the progress line
  of the row names its check, which runs the one script named here. The scratch
  checkout keeps the target's `.gitattributes`, and every call of
  `internal/git` sets `core.autocrlf=false` (K7), so the `eol=crlf` fixtures
  keep their bytes. **Known limit:** no timeout in phase 1.
- **The call of S05** (K10): `internal/verify` gives the files whose links
  break. It takes the directory of a checked-out tree (S05 gives
  `WORK/target` after the deletion and before its commit), runs
  `sh docs/links/link-lint.sh` of that tree as above, and reads the script's
  standard error: each line `FAIL  L<n>: <path>:<line> …` names a file; the
  files, once each, in byte order. A line `FAIL` that names no file (for
  example `L5`, no in-tree link resolved, or a root not found), or an exit
  other than 0 with no file, is an error, and S05 fails. S05 (row 13) gets the
  list through `internal/cli`, and stops with one `F-<path>` row per file.
- **The shared fixtures:** check `markers` passes and fails on the cases of
  `docs/setup/tests/markers/` with the lines of `check_markers`. The two script
  checks have no fixture group of their own: `setup-check.sh` runs the scripts
  inside `kit-linters`, a check of LAYUP's own form with four linters, whose
  frame cases the harness does not compare; their tests run stub scripts and
  LAYUP's own scripts on a clone of its `HEAD`.

**Decided here** (task `T-9t1q`, #89), the checks `facts`, `onboarding`,
`glossary` and `guardrails` in a target's form, and the records that they read:

- **The two answers records** (K14, O-124): the record of the `S01-` and `Q-`
  answers is S04's, and the record of the `M-` answers is S11's. The rows S04
  and S06 of the step table say it in row 9 of the
  [plan](../plan/README.md#the-tasks-of-phase-1) (K14).
- **The form of an answers record** (K15): the form of the baseline's
  `docs/facts/template.md`: its header table, and under `## Facts as collected`
  one numbered fact per row of `answers.tsv` that its step reads, in the order
  of the stop table, each ``N. `<question ID>` <answer> — by <by>; source
  <source>; the question: <question text>``, with the question text of the stop
  table's `ask` (for a `Q-` question, the writer of S04 computes the gap table
  again from the brief that `brief.sha256` holds). The answer and the question
  text are never in a code span. Its path is
  `docs/facts/F-NNNN-setup-answers.md` (S04) or
  `docs/facts/F-NNNN-marker-answers.md` (S11), `F-NNNN` the next free ID
  (`F-0001` and `F-0002` on the baseline at its pin, whose index has no row),
  with its index row and its line in `facts.sha256`. Check `facts` reads the
  question ID of a fact: the code span that opens the text after `N. `.
- **No marker in an answers record** (K17): the writer writes each left and
  right angle quote of a recorded text as `&lsaquo;` and `&rsaquo;`, which
  render the same, and the record's "Notes on capture" says so. Check
  `markers` reads `docs/facts/` (it is not in `MK_EXEMPT`), so it finds no
  marker in a record; a marker that stays (the answer `gap`) is only in its own
  file and in `open-gaps.tsv`. The answer `gap` is the word `gap`, with the
  question text.
- **What check `facts` reads** (K18): a record or a brief that exists is read
  in full whenever it exists: its question IDs (each one of its kind, a row of
  `answers.tsv`, once; and each row of its kind has a fact in it), its line in
  `facts.sha256`, its index row. The `done` rows of the setup record decide
  only what must exist: `facts.sha256` and the record of S04 once S04 is done;
  the problem statement once S06 is done; the record of S11 once S11 is done
  and `answers.tsv` has an `M-` row. So the evidence call of a step, which
  comes before its `done` row, reads what the step wrote, and an `M-` row given
  early is no finding until S11 is done.
- **The briefs:** S06 copies the problem statement to
  `docs/facts/problem-statement-brief.md` and the vision brief to
  `docs/facts/architectural-vision-brief.md`, the names of LAYUP's own raw
  files, so check `facts` and check `onboarding` read the paths of their sh
  functions. The vision brief is listed when its file exists.
- **The lines of check `facts`:** the hash loop of `check_facts` with its
  lines (`hash:`), which hashes the file of the scratch tree; the engine's git
  calls set `core.autocrlf=false` (K7), so a host's setting does not change the
  bytes. Then `brief:`, `listed:`, `record:`, `answers:` and `index:`, as the
  tests of `internal/verify` name them.
- **The fact resolver:** `F-NNNN#n` resolves when `docs/facts/` holds exactly
  one `F-NNNN-*.md` with a numbered fact `n` (a line `^0*n\. `); the finding is
  `fact: F-NNNN#n is not a fact of the F-NNNN record`. `check_onboarding` reads
  the first `F-0001-*.md`; the engine wants exactly one record per ID.
- **Onboarding, glossary and guardrails:** each file must exist (`missing:`).
  Onboarding: the marker and link rules of `check_onboarding`, and the
  resolver. Glossary: the resolver for each citation; LAYUP's heading, its 25
  rows and `F-0001#15` to `#39` do not apply. Guardrails: an entry is a bullet
  that starts with `- **Inv-N**`, up to the next entry, a heading or the end of
  the file; a file with no entry has no finding; an entry's `Check:` value is
  read as `check_guardrails` reads it (after the last `Check: ` of the first
  line that has one), with its `check:` lines, and its path must be a file of
  the tree (`check_guardrails` also takes a path outside it, such as `../x`);
  each citation of the file resolves.
- **The shared fixtures:** for these four checks, the engine passes and fails
  on the same fixtures as `docs/setup/tests/run.sh` in the lines of the kinds
  that a target keeps with the same text: `hash` (facts); `missing`, `marker`,
  `link` and `fact` (onboarding); none (glossary); `check` (guardrails). The
  harness compares a case, its lines of those kinds and its exit, only when its
  `EXPECT` has such a line; the cases with none are LAYUP's form and are not
  compared: `facts/bad-answers-blank`, `bad-answers-missing`,
  `bad-answers-repeat`, `bad-batch-absent`, `bad-batch-rows` and
  `bad-blank-tab`, `glossary/bad-no-section` and `bad-rows`, and
  `guardrails/bad-no-section`; nor is `facts/good-autocrlf`, whose mode clones
  LAYUP's own repository (#48).

**Decided here** (task `T-8ya0`, #88), check `adapted`:

- **The rule is `check_adapted` as it is** (O-123): a target's check reads
  LAYUP's `AD_EXCLUDE` and `ad_allowed`, which the engine embeds at its version,
  as S10 reads LAYUP's `MK_EXEMPT`; a path that a target does not have matches
  nothing. `TestTheListsOfAdaptedEqualTheSh` compares the two lists with
  `setup-check.sh`.
- **The files** are the paths of `git ls-files -z` whose name ends in `.md` (in
  that case, as the pathspec `'*.md'` of the sh function reads it), that
  `AD_EXCLUDE` does not match, and that are regular files when links are
  followed (`[ -f ]`). The sh function lists the paths without `-z`, so it skips
  a path that git quotes (one with `"`, `\` or a control character in it); the
  engine reads it, because a tracked file is read, not skipped. A failed list is
  the finding `git: cannot list the tracked files`.
- **The text is bytes,** as `awk` reads it in the C locale; a NUL byte is one
  more byte, where macOS `awk` ends the line at it. A line loses a last
  carriage return; each run of spaces and tabs becomes one space; a space at
  the start or at the end of a line is removed; a line with no other character
  ends the paragraph; a marker from `‹` to the first `›`
  after it, or to the line end when it does not close, is one unit; the lines
  of a paragraph are joined with one space. The lower-case text changes only
  `A` to `Z`, so a position in it is the same position in the text.
- **A hit** of one of the 16 patterns of the sh function is at the first letter
  of its word, the end of the boundary before it, so a boundary of more than one
  byte (for example `é`) does not move it; its line is the line where its word
  starts. A hit of the first pattern of rule 1 whose eleven lower-case bytes
  from its first letter are `kit-history` or `kit-linters` is not a finding.
- **The findings** are `<rule> <name>: <path>:<line>`, each once, in byte order;
  the first is the reason of the row. The sh function sorts in the locale of its
  host; byte order makes the reason one value on every host (`NFR-005`).
- **The list for the prose step** (K12, O-123): `Flagged` of `internal/verify`
  gives the path of each file of a work tree that the check flags, once, in byte
  order; the prose step (row 13) stops for each one, through `internal/cli`.
- **Measured:** on the tree of LAYUP's root commit `d2516fd` (the unchanged
  baseline), the engine and the sh function in the C locale give the same 441
  lines in 38 files. **Known limit of the sh function:** macOS `awk` 20200816
  in a UTF-8 locale stops on a hit whose word follows a character of more than
  one byte ("towc: multibyte conversion failure"), and the sh check then passes
  that file; the engine's check does not read the locale.

### The stack catalog

One directory per stack in LAYUP's repository, `internal/catalog/<stack>/`,
embedded in the binary with `embed`. **Decided here:** the path, because `embed`
reads only files under the package's own directory; and the embedding, because
a target starts from the defaults of the LAYUP version that sets it up (§13).
This is the form of an entry; the Go entry is below (task `T-c06a`, #91).

| Path | What |
| ---- | ---- |
| `internal/catalog/<stack>/kinds.tsv` | one row per gate kind (the schema below) |
| `internal/catalog/<stack>/files/<path>.tmpl` | each file that the setup writes into the target at `<path>`: the tools' configuration, the CI workflow of the gate jobs; `{{module}}` in a file is replaced with the module path from `name`, and a gap token `{{gap:<text>}}` with the marker of that text |
| `internal/catalog/<stack>/fixtures/<kind>.patch` | the known-bad fixture of a kind: a patch that `git apply` applies to the setup head and that must make the kind fail |
| `internal/catalog/<stack>/gaps.tsv` | the gaps of the entry, when it has one: for each gap token, its file and its question (the block below; task `T-c06a`, #91) |

```tsv-schema catalog-kinds layup:internal/catalog/<stack>/kinds.tsv
kind      id(<word>)            key  the gate kind, as in the manifest
state     enum(active|pending)  -    the state at setup; `pending` for a kind whose rules depend on the architecture (§6)
tool      text                  -    the program, as in the manifest; `—` for a `pending` kind until its activation
version   text                  -    the tool's version that the evidence documents; `—` for a `pending` kind until its activation
command   text                  -    the command, as in the manifest; `—` for a `pending` kind until its activation
scope     list(text)            -    the scope patterns, as in the manifest
config    list(path)            -    the gate files of the kind in the target; `—` when none, and for a `pending` kind until its activation
fixture   text                  -    `fixtures/<kind>.patch`, relative to the entry's directory; `—` for a `pending` kind until its activation
evidence  text                  -    the URL of the tool's documentation at that version, `https`; `—` for a `pending` kind until its activation
```

```tsv-schema catalog-gaps layup:internal/catalog/<stack>/gaps.tsv
path      path  key  the file of the target that holds the gap token: a file of `files/`, without `.tmpl`
marker    text  key  the text of the marker, as the gap token `{{gap:<text>}}` holds it
question  text  -    the question of the gap, for the row of `docs/setup/open-gaps.tsv` in the target
```

The manifest that S12 writes is this table without the columns `version`,
`fixture` and `evidence`. LAYUP's own CI runs each fixture of each entry (§6),
in its job `tests` (task `T-c06a`, below).

**Decided here** (K35 of the [defect register](../plan/README.md#the-defect-register),
task `T-3jpx`, #81):

- **The suffix `.tmpl`.** Each file of `files/` ends with `.tmpl`, and its
  path in the target is its path without `.tmpl`: `files/go.mod.tmpl` is the
  target's `go.mod`. The `embed` pattern of an entry has the prefix `all:`
  (`//go:embed all:<stack>`), so that `embed` keeps `.github/`. Reason: `embed`
  skips a directory that holds a `go.mod`, the root of another module, with
  every file in it and with no error, and leaves out `.github/` without `all:`
  (measured, `runs/T-3jpx/`); a `.go` file under `files/` would be a package of
  LAYUP's module for `go vet`, `go test` and `go list`, and `gofmt` reads it.
  One suffix for every file is one rule with no exception. Because the skip is
  silent, the test of each entry checks the list of its files.
- **The rules of an entry**, which `internal/catalog` checks when it reads one:
  `kinds.tsv` matches its schema and has at least one row; an `active` kind
  names `fixtures/<kind>.patch`, and that file exists; a `pending` kind has `—`
  as its fixture; each file of `fixtures/` is the fixture of an `active` kind;
  each file of `files/` ends with `.tmpl`. Task `T-c06a` (#91) adds: an
  `active` kind has a version, and an `https` URL as its evidence; a `pending`
  kind has `—` in each column but its kind, its state and its scope; `gaps.tsv`,
  when there is one, matches its schema, and each of its rows names a file of
  `files/` that holds its gap token once, with a question; each gap token of a
  file has its row; no file of the entry holds a marker character (LAYUP's own
  check `markers` would read it, so an entry writes a gap token in its place).
  Reason: an entry that breaks one of
  them sets up a target whose gate cannot be proven (a fixture that no kind
  runs, or a kind with no fixture), so the error comes when LAYUP reads its own
  catalog, not at a target's setup. A rule for the authors of an entry, which
  no reader checks (task `T-5sgt`, #82): a fixture does not change a `config`
  path of its own kind, because `layup gate` puts the base's gate files back
  in its scratch tree ([`gate.md`](gate.md#the-command)), so such a fixture
  would test nothing.
- **The manifest** is written by `internal/catalog` through `internal/tsv`,
  by the block `gate-manifest` of [`gate.md`](gate.md#the-gate-manifest).
- **A `catalog` ref** of the setup record (`<stack>/<path>`) names a file by
  its path in the entry's directory, for example `go/kinds.tsv` or
  `go/files/go.mod.tmpl`, so a reader finds it at `internal/catalog/<ref>` in
  LAYUP's repository; check `sources` resolves it there. Reason: the record
  says "a catalog file", and the file that holds the value is the file of the
  catalog.
- **The test entry** of the package's tests is
  `internal/catalog/testdata/test/`, embedded only by those tests, with the same
  pattern form (`all:`). The binary embeds the Go entry (task `T-c06a`, #91),
  and never the test entry, so S01 never accepts a stack named `test`.
- **No conversion of the bytes:** `.gitattributes` holds
  `internal/catalog/*/** -text`, so a checkout never changes the line endings
  of a file of an entry or of the test entry. Reason: a fixture is a patch that
  `git apply` reads, and each file goes into a target byte for byte.

**Decided here** (task `T-c06a`, #91), the Go entry, `internal/catalog/go/`:

- **The kinds** (K20 of the [defect register](../plan/README.md#the-defect-register)):
  `static`, `layout`, `boundary`, `contract` and `test`, in the order of the
  kind list of [`gate.md`](gate.md#the-gate-manifest). `static` and `test` are
  `active`; `layout`, `boundary` and `contract` are `pending`, because §6 makes
  `pending` only a kind whose rules depend on the architecture, and the tests
  of `test` need none. A `pending` kind has `—` in each column but its kind,
  its state and its scope: its command never runs and its tool is never looked
  up ([`gate.md`](gate.md#the-command)), and the activation batch writes its
  tool, version, command, configuration and evidence (§6, Activation); a value
  with no source is not written (Invariant 4). So ADR-0016 decision 1 ("an
  entry names the tool and its version, the command") holds for each `active`
  kind. Each kind has the scope `./*.go`.
- **The commands:** `static` is `out=$(gofmt -l .) && test -z "$out" && go vet
  ./...`, and `test` is `go test -count=1 ./...`, each with the tool `go`.
  Reason: one row has one command, and §6 names the two static checks as one
  kind. In the form `test -z "$(gofmt -l .)"` the shell loses the exit status
  of `gofmt`, so a `gofmt` that is not found, or that cannot parse a file,
  gives no output, and the kind passes on nothing; the assignment keeps the
  status (measured: 127 when `gofmt` is not found, 2 on a file that does not
  parse; `fail`, `NFR-004`). `gofmt` comes with `go` in a Go distribution.
- **The version and the evidence** of each `active` kind: `1.26`, the `go`
  line of LAYUP's own `go.mod`, and `https://pkg.go.dev/cmd/go@go1.26.0`, the
  documentation of the `go` command at the first release of that line (it
  documents `go vet`, `go test -count=1`, and `go fmt`, which runs `gofmt`),
  read on 2026-10-02. The target's `go.mod` has `go 1.26`. Reason: LAYUP's CI
  installs Go from its own `go.mod`, so the fixture test runs on that toolchain
  and downloads none.
- **The tool of a kind** of the Go entry is `go` or a program of the runner,
  because its workflow installs only Go and the activation changes no CI file
  (§6); a kind that the activation makes active fetches any other tool in its
  command, for example `go run <module>@<version>`.
- **The workflow** (K30) is `.github/workflows/gates.yml` in the target, a name
  that no workflow of the baseline at LAYUP's pin has. It runs on
  `pull_request`, with `permissions: contents: read`, and has one job per kind,
  in the order of `kinds.tsv`. The id and the `name:` of each job are the kind,
  so a reader that takes the name, and a ruleset that takes the context of a
  check from the name, both get the kind. Each job checks out the head commit
  of the pull request with its history (`actions/checkout@v4`, `fetch-depth:
  0`), installs Go from `go.mod` (`actions/setup-go@v5`, `go-version-file`), and
  runs `sh .github/gates.sh <kind>` with the base and the head commits in
  `GATE_BASE` and `GATE_HEAD`. No job has `paths:` or a job-level `if:`, because
  a skipped required check counts as passed; no job starts `layup`, fetches a
  LAYUP file or runs `setup-check.sh` (`NFR-002`, ADR-0016 decision 5). Each
  `uses:` line ends with its evidence, the action's repository at that version
  (`https://github.com/actions/checkout/tree/v4`,
  `https://github.com/actions/setup-go/tree/v5`, read on 2026-10-02): the
  versions that LAYUP's own CI runs.
- **The job script,** `.github/gates.sh` in the target, holds the rules of the
  run of [`layup gate`](gate.md#the-command) for one kind, in POSIX `sh`, `awk`
  and `git`. It reads the kind's row of the head's `docs/gates.tsv` (known
  limit L-B3), prints the kind, the result and the reason of `gate.md`, and
  exits 0 for `pass` and `clear`, 1 otherwise. A manifest that `layup gate`
  refuses with exit 2 (none, no row for the kind, a form of another kind, a
  scope pattern of another form) is `fail` in the job: a job has no third
  state, and each of the two is not a pass. **Known limit:** a file name with a
  line feed is split by `tr`.
- **The coverage floor** (K24, L-B2): the file `docs/gates/coverage-floor.txt`
  of the target, the `config` of `test`, has one line, the marker of the floor,
  and `gaps.tsv` gives its question. S12 (row 15 of the plan) writes, for each
  gap of the entry (`Entry.Gaps`), its row of `docs/setup/open-gaps.tsv` and
  its record row `marker:<path>:<line>` (source `gap`). The command checks no
  floor until the idea owner or the Operator sets it with evidence.
- **The fixtures** add only new files, under `gatefixture/` (no baseline at
  LAYUP's pin has that directory), and change no `config` path of their kind:
  `static.patch` adds a file that `gofmt` changes, and `test.patch` a test that
  fails. Reason: a patch that changes a file of the baseline stops applying
  when a later baseline changes that file (O-101). A target rendered from the
  entry has no `.go` file, so on its clean commit `static` and `test` are
  `clear` (`no product path`), and each fixture adds the first `.go` file.
- **The tests in LAYUP's CI**, in the job `tests`, with no change of `ci.yml`
  (K28): `TestTheFixturesOfEachEntry` of `internal/catalog` (integration:
  `gate.Run` on a target rendered from each entry of the binary, on its clean
  commit and on each fixture), `TestTheJobScriptOfEachEntry` (integration: the
  job script and `gate.Run` on the same base and head, for each line of the
  table of the run and for the input errors, give the same pass or fail), and
  `TestGateOnEachEntryOfTheCatalog` of `cmd/layup` (e2e: the built binary).

### Not in phase 1

- Each part of §3 and §5 that the boundary table gives to `layup run` or to a
  session.
- The ruleset of the records branch (§6: updates only by the LAYUP App), and the
  `layup/` required checks of the default branch (`layup/gates`, `layup/spec`,
  `layup/verify`, `layup/rules`, §6): the Operator applies both with the
  commands that `layup run` prints when it first starts on the target (phase 2).
- The activation of the `pending` kinds at the first bet (§6): phase 2. The
  tool of a kind of the Go entry is `go` or a program of the runner (The stack
  catalog, task `T-c06a`).

## NFR-003 — No value without evidence

Requirement: "No configuration value is set without evidence; a value that
comes from a guess is a defect." Derives from `architecture.md` §5 (Scaffold
1); no ADR; [`guardrails.md`](../guardrails.md) §1.1, Inv-4.

1. Each value that `layup setup` writes into the target has a row in the setup
   record, with a source of one of the kinds `answer`, `catalog`, `fact`,
   `computed` or `gap` ([the setup record](#the-setup-record)).
2. A marker with no answer stops the run (S10); an answer `gap` keeps the
   marker, with a row in `open-gaps.tsv`. `layup setup` never writes a value
   that it did not read from one of these sources.
3. Check `sources` of `layup setup verify` fails on a value row with no
   source, or a source that does not resolve; check `markers` fails on a
   marker with no `open-gaps.tsv` row.
4. Whether a cited source supports its value is a judgement: the audit that
   the criterion names is a review, not a phase-1 check.

## NFR-006 — The baseline at a pinned, recorded version

Requirement: "LAYUP uses [the baseline] at a pinned, recorded version." Derives from
`architecture.md` §5 (Start 2), ADR-0009.

1. **LAYUP's own pin** is [`docs/setup/armature.pin`](../setup/armature.pin);
   LAYUP's check `pin` passes (the criterion of `NFR-006`). It does not bind a
   target.
2. **A target's pin** is the latest commit of the baseline when the target is
   set up (O-101), resolved once (S02) and written at S04:

```text
source=<the baseline's repository URL>
commit=<40 hexadecimal characters>
tree=<the tree of that commit>
method=git clone <source>, checkout <commit>, .git removed
date=<the date of pin.time, YYYY-MM-DD>
```

   The keys and their order are those of LAYUP's own pin file, so the same
   check `pin` reads both. Each key appears once. **Decided here** (task
   `T-6x75`, #84): S02 writes `pin.time` in UTC, in the form
   `YYYY-MM-DDTHH:MM:SSZ`, and the `date` of the pin is its first ten
   characters.
