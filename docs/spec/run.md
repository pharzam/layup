# `layup run`: the Start of a target

Milestone `M2a` of the [implementation plan](../plan/README.md#milestones), the
first milestone of phase 2: `layup run --new` makes the first records commit of an
empty target, and `layup run TARGET` restarts on it. Requirements: `NFR-001`,
`NFR-002`, `NFR-006`, `REQ-002`. It derives from
[`architecture.md`](../architecture.md) §1 (the host, the forge), §2 (one run per
target, the lease, fencing), §3 (one writer, identities, a human decision, copy
before read) and §5 (Start 1 to 3), and from ADR-0013 and ADR-0014. Conventions:
[`README.md`](README.md). The records: [`records.md`](records.md#nfr-001--the-records-of-start).
The forge: [`forge.md`](forge.md). Task `T-zck8` (#123) wrote it.

**Decided by the Operator** (O-163, #123): `M2a` specifies an empty repository
only (`--new`). A target that phase 1 set up, whose records branch already holds
the Operator's commit (K41), is adopted in milestone `M2d`, where `layup run`
drives the Scaffold.

## The command

```text
layup run --new OWNER/NAME --host DIR --psb FILE [--vision FILE]
          --operator LOGIN --idea-owner LOGIN --plan PLAN
          --intake-cap MONEY,HOURS --lease-h MINUTES --watch-t MINUTES
layup run TARGET --host DIR
```

- **Two rows of the command table** (decided here). The command table of
  `internal/cli` ([`README.md`](README.md#commands)) holds two rows with the word
  `run`: Start, which has the flag `--new`, and the restart, which has one
  positional argument. **Decided here:** a row may name one selecting flag; when
  two rows have the same words, the row whose selecting flag is given is the
  command, else the row with none. So `layup run --new …` is Start, and
  `layup run TARGET` is the restart; a restart with a flag of Start, or Start with
  a positional argument, is a usage error (exit 2). Reason: the architecture
  names both forms (§5), and the one-table rule of `README.md` stays.
- `OWNER/NAME` and `TARGET` name a repository on the forge, in the form
  `OWNER/NAME`. `--new` takes it as its value; the restart takes it as its one
  positional argument.
- **`--host DIR`** (decided here, K40): the directory of the LAYUP host, the root
  of each `host:` location. `layup run` reads
  `DIR/registers/forge.tsv` and `DIR/registers/harnesses.tsv`, and keeps its own
  clone of each target in `DIR/targets/OWNER/NAME/`. Reason: one flag names every
  file of the host, so a restart reads the same values; a positional argument
  would make `layup run TARGET` and `layup run --new` take it in two places.
- `--psb FILE` and `--vision FILE` are the two briefs. **Decided here:** each flag
  of `layup run` is required, except `--vision`, which may be left out
  (`architecture.md` §5 gives it in brackets). A missing
  `--vision` writes `vision.sha256` as `—` and no `start/vision.md`.
- `--plan` is the forge plan that the Operator names (`free`, `pro`, `team`,
  `enterprise`); `--intake-cap` is the cap of Intake and Shape (milestone 0,
  §12), a decimal amount of US dollars and a decimal number of hours;
  `--lease-h` and `--watch-t` are `lease.H` and `watch.T` (§10), whole minutes,
  1 or more.
- The [rules of `README.md`](README.md#commands) hold: no environment variable,
  no standard input, no question on the terminal.

**Output** (decided here): `layup run` prints one table, `run-steps`, on standard
output at its exit, one row per step that it reached, in the order of the steps;
and its progress on standard error, by the progress rule of `README.md`, with
`[<i>/<n>]` over the rows of the table. A wait for a human (the root push) prints
one progress line every ten seconds. Reason: a Start ends; the phase loop of later
milestones adds its own rows, so the one-table rule holds.

```tsv-schema run-steps stdout
step enum(forge|plan|baseline|root-push|read-back|records|issues|watch|clone|version|lease|phase) key the step, in the order of the two lists below
result enum(done|fail) - `done`, or `fail` with its reason in `detail`
detail text - one line: what the step found, with no time, no path of the host and no SHA of a records commit (a records commit holds its time, so its SHA differs on a repeat)
```

**Exit codes** (decided here): 0 when each row is `done`; 1 when a row is `fail`
(the run stops at that row, and prints the table of the rows it reached); 2 for a
usage or input error that the run finds before its first step, with no table: a
flag, a brief that is not a readable UTF-8 file, a register that its reader
refuses, and a key file that is missing, not PEM, not mode 0600 or owned by
another user. Each state that a step finds on the forge or in the records is
`fail` on that step, with exit 1: a repository that does not exist, `--new` on a
repository that is not empty (O-163), and a record of the records branch that its
reader refuses. Code 3 is not used: the wait for the root push is in the
foreground.

## NFR-001 — The Start of a target

From Start on, `layup run` is the one writer of the records branch
`layup-records` (ADR-0014; [`records.md`](records.md#nfr-001--git-is-the-system-of-record),
item 1). It commits from its own clone, `DIR/targets/OWNER/NAME/`, which no
session uses.

### The steps of `layup run --new`

1. **`forge`.** It makes an installation token for `OWNER/NAME`
   ([`forge.md`](forge.md#the-app-identity)), checks that the adapter gives the
   six capabilities and that the installation has the permissions that `M2a`
   uses (contents, issues: write; metadata: read), and records each permission of
   the installation in `app.permissions`. A missing permission of a later
   milestone (commit statuses, `M2e`) is recorded, not a `fail`. The repository
   must exist and have no commit (O-163). It reads the numeric IDs of
   `--operator` and `--idea-owner`; an unknown login is `fail`.
2. **`plan`**: the plan check of §5 (K15), on the visibility that step 1 read.
   `public` passes on each plan; a private repository passes on `team` and
   `enterprise` only (GitHub Free has no rulesets and no draft pull requests on
   it, and GitHub Pro no drafts). A failed check names the plan and the
   visibility.
3. **`baseline`** (Start 2, S02 of [`setup.md`](setup.md#the-steps)): it resolves
   the latest commit of the baseline's default branch (`LsRemote`), clones that
   commit, and removes `.git`. The pin (source, commit, tree, time) is held in
   the run's memory until step 6 writes it.
4. **`root-push`** (Start 2, S03): it makes the root commit as S03 does, and
   prints on standard error the one command that pushes it, with the resolved
   commit and its difference from LAYUP's own pin. The Operator runs it with the
   Operator's own login (the copy holds CI files, and the App has no workflows
   permission, O-92). `layup run` reads the list of branches every ten seconds
   (`Repository`, [`forge.md`](forge.md#the-calls-of-m2a)) until the default branch exists. **Decided here:** the wait has no limit (Ctrl-C stops it), as
   a wait for a human is not a stall (ADR-0023); a run stopped here leaves a root
   commit with no records, and the Operator starts again with an empty
   repository (§5).
5. **`read-back`**: `DIR/targets/OWNER/NAME/` is made with `Init` and the
   default branch fetched into it (`Fetch`). The default branch is the
   repository's only branch, its head is one commit with no parent, its root tree
   equals `pin.tree`, and the visibility equals that of step 1. Any difference is
   `fail`.
6. **`records`**: the first records commit, an orphan commit made in a scratch
   work tree as S15 makes it (`SwitchOrphan`), pushed to `layup-records` with the
   installation token: `start/start.tsv`, `start/problem-statement.md` and
   `start/vision.md` byte for byte, `approvers.tsv` (the two rows of the Start
   command), `lease.tsv` (this run, `held`), and the README of the records branch
   in [the text of Start](#the-readme-of-a-target-that-start-makes). In `start.tsv`,
   `issue.intake`, `issue.control` and `watch` are `—` until their steps write
   them. The author and the committer are the App's bot (`<slug>[bot]`, the
   e-mail `<bot-id>+<slug>[bot]@users.noreply.github.com`), at the time of the
   commit. From this commit on, only a comment by an ID in `approvers.tsv`, in the
   role that a rule names, is an answer or a decision.
7. **`issues`**: for the Intake issue, then the control issue: a records commit
   sets its row (`issue.intake`, `issue.control`) to `opening`, which announces
   the write; the run opens the issue; a records commit sets the row to the
   issue's number, so a restart finds it with no search.
8. **`watch`**: it reads the comments of the control issue for up to `watch.T`
   for a notice from the App `watch_slug` (the dead-man job, §11). Before it acts
   on a comment, it pushes the copy ([copy before read](#copy-before-read)); then a
   records commit writes `watch`: `confirmed` when such a comment came, else
   `not-confirmed` (L-F1); with no `watch_slug`, `not-confirmed` at once. Both
   are `done`.
9. **`lease`**: the last records commit sets the lease row to `released`, and the
   run exits.

**Decided here:** steps 1 and 2 put every check that needs no root commit before
the root push, so the Operator pushes only to a target that can be held; the
architecture gives the plan check after the push (§5, Start 3), and the push does
not change the plan or the visibility, which step 5 reads again. In `M2a`, a
Start ends after step 9; the Intake of `M2c` adds its rows after `watch`.

### The README of a target that Start makes

The text of `README.md` of `layup-records` that step 6 writes (**decided here**:
the text of [`setup.md`](setup.md#the-readme-of-the-records-branch) names the
files of the setup, which a target that Start makes has only after the Scaffold;
the Scaffold of `M2d` replaces this text with that one when it adds them):

```text start-readme
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
```

### The restart

`layup run TARGET --host DIR`:

1. **`forge`**: the installation token, the six capabilities and the permissions,
   as step 1 above. The repository must exist and hold `layup-records`. No login
   is looked up: the approvers are read from `approvers.tsv`.
2. **`clone`**: it rebuilds `DIR/targets/OWNER/NAME/` from the forge (a fresh
   `Clone` of the default branch and `layup-records`); nothing in the old clone is
   read. It reads `start.tsv`, `approvers.tsv` and `lease.tsv`.
3. **`version`**: the LAYUP version of the run equals `layup.version` of
   `start.tsv`; else `fail` (the inventory: "stops when its own version differs").
4. **`lease`**: [the lease](#the-lease-and-fencing) is taken: at once when it is
   `released`; else after the takeover rule. A records commit sets this run as
   the holder.
5. **`phase`**: the first step of Start that is not done, run again. A row
   `opening` means that a run stopped between the announcement and the record
   of an issue: the restart runs step 7 for that issue again (a duplicate issue
   is possible, known limit below). When each step of Start is done, the next
   phase is Intake (`M2c`): the row is `done`, and the run releases the lease and
   exits. `M2c` replaces this row with its steps. **Decided here:** `phase` is one
   row of the restart's table, whatever steps of Start it runs again; its
   `detail` names them, so the key `step` holds no second `lease`.

### The lease and fencing

As `architecture.md` §2, with these values decided here:

- The heartbeat is a records commit that adds 1 to `heartbeat`, every `lease.H`,
  while the run holds the lease, the wait of step 8 included. Before step 6 no
  lease exists, so nothing beats during the wait of step 4.
- A second run reads the lease row every ten seconds and times it by its own
  clock from the moment it last saw `heartbeat` change. When `heartbeat` has not
  changed for `3 × lease.H`, it takes the lease over: a records commit that
  writes its own row, with the old run ID in the commit message. While
  `heartbeat` moves, it does not take the lease; after `3 × lease.H` with a
  moving counter it exits 1, `lease`: `fail`, naming the run that holds it. At
  each read the takeover rule comes first: no change for `3 × lease.H` since
  the last change it saw takes the lease over; else `3 × lease.H` since the
  start of its wait exits 1. So a counter that moves once and then stops gives
  exit 1, and the next run takes the lease over (task `T-trej`). A refused takeover
  push means that another run took the lease first: `fail`, naming it.
- **Fencing.** A run pushes the records branch only on top of its own last pushed
  records commit, never after a fetch and a rebase, and never with force. Each
  forge write (an issue, a comment) comes after the records push that announces
  it. A refused records push stops the run: it reads the lease again, goes on
  only while it still holds it, and else exits 1 with `fail` on the step. To go
  on is to try the same push once more; a second refusal is `fail` (task
  `T-trej`). The heartbeat and the release go through this rule.

### A human decision

A comment is a human decision only when its author ID is in `approvers.tsv`, in
the role that the rule names, and its App field is empty (§3, O-77). A review, a
review comment, a commit, a reaction or an edit is never a decision. The
decision is the first copy of the comment (`seen` 1 of `copies.tsv`), so a
comment edited before its first copy is a decision as copied (task `T-trej`). The rules
of `M2a` name no decision; the function and its test come with `M2a`, so that
`M2c` uses them: the Intake answers take the roles `operator` and `idea-owner`,
an acceptance the role `idea-owner` only.

### Copy before read

Before any step acts on a comment, `layup run` writes a row of `copies.tsv` and
its body file, and pushes them. A comment that a later read finds with another
SHA-256 of its body is an edit: a new row with the next `seen`, and the first
copy stays. A comment whose SHA-256 is the same as its last copy gets no new
row. In `M2a`, the comments of the control issue at step 8 are copied.

### The Intake and control issues

**Decided here:** the titles are `LAYUP Intake` and `LAYUP control`. Each body
names the records branch and says that a comment there is recorded before LAYUP
acts on it, and that only the accounts of `approvers.tsv` decide. The App's bot
opens both. In `M2a` the Intake issue receives no comment from LAYUP; the gap
check of `M2c` posts there.

### Input states

**Decided here**, for each state that the steps above never make:

| State | Result |
| ----- | ------ |
| A register with two rows for one key, an unknown or missing column, or a field that its type refuses | exit 2: the reader of `internal/tsv` refuses it and names the line |
| A forge register with no row | exit 2 |
| A harness register with no row | allowed: `start.tsv` has no `harness.<id>.*` row; `M2b` admits no session until it has one |
| A harness row with `wall` empty | exit 2 |
| A key file that is missing, not PEM, or not mode 0600, or owned by another user | exit 2, naming the file and its mode |
| `--new` on a repository with a commit, or with the branch `layup-records` | `forge`: `fail` (O-163) |
| A restart on a repository with no `layup-records` | `forge`: `fail` |
| A restart whose `start.tsv`, `approvers.tsv` or `lease.tsv` its reader refuses | `clone`: `fail` |
| A lease row from another LAYUP version | the `version` step fails first |
| A lease table with no row, or two rows | `clone`: `fail` |
| A forge error during a step | `fail` on that step, with the error ([`forge.md`](forge.md#forge-errors)) |

**Known limit:** a run that stops after it opens an issue and before the records
commit of its number leaves an issue that the records do not name; the restart
opens a second one. The Operator closes the first, which the App's bot opened
with the same title.

### Not in M2a

The adoption of a phase-1 target (`M2d`, O-163); the Scaffold, the rulesets, the
probes and the forge read-back of the rules (`M2d`); the gap check and the
answers (`M2c`); role sessions, the harness probe, admission and the harness
credential (`M2b`, K40); the commit status `layup/gates` (`M2e`); the dead-man
job itself (`M3d`; `M2a` only reads its notice).

## NFR-006 — The target's pin at Start

The baseline commit that step 3 resolves is the target's own pin (O-101): its
source, commit, tree and time go into the first records commit (`pin.*` of
`start.tsv`). Once that commit exists, no run resolves the commit again; the
Scaffold (`M2d`, step S04) writes `docs/setup/armature.pin` of the target from
that record. LAYUP's own pin does not bind a target. A run that stops before
step 6 leaves no pin, and the Operator starts again with an empty repository.

## NFR-002 — The records with no tool

Each file of the records of Start is a table of [`README.md`](README.md#records)
or a brief byte for byte, on a branch that a plain `git clone` carries as
`origin/layup-records`. The uat of `M2a` reads them so.

## REQ-002 — The Start before the Scaffold

`REQ-002` ("`layup setup` creates a target repository from the pinned
[the baseline] …") is the Scaffold's. `M2a` gives it the steps before it: the
resolved pin and the root commit come from the same calls as S02 and S03 of
`layup setup`, so the Scaffold of `M2d` starts on the commit that Start read back.

## The acceptance tests of M2a

| Part | Level | Test |
| ---- | ----- | ---- |
| The lease and fencing | unit | With a stand-in clock and a stand-in records store, a second run takes the lease only after `heartbeat` has not moved for `3 × lease.H` by its own clock; with a moving counter it exits 1; a refused records push stops the run until it reads the lease again. |
| A human decision | unit | An author ID in `approvers.tsv` with an empty App field, in the rule's role, is a decision; an App comment, a review, a review comment, a commit, a reaction and an edit are not; Intake takes `operator` and `idea-owner`, an acceptance `idea-owner` only. |
| Copy before read | unit | A copy writes the body, the author ID and login, the comment ID, the App field, the times and the SHA-256; an edit adds a row with the next `seen`, and the first copy stays; an unchanged comment adds none. |
| The records of Start | unit | Each new record kind is refused with a wrong header; its Go schema equals its block (`tsv.Compare`), and the owner moves its name from `notYetBuilt` to `built`. |
| Input states | unit | Each row of the table of input states gives its result. |
| The GitHub adapter | integration | Against an `httptest` server on loopback: the JWT and the installation token, each call of [`forge.md`](forge.md#the-calls-of-m2a), and a server that lacks a capability or a permission, which stops `forge` and names it: the permission, or the call and the status ([`forge.md`](forge.md#the-test-of-the-adapter)). |
| Start | integration | Against the `httptest` forge and a real local bare repository with the real `internal/git`: the first records commit holds the pin, the briefs byte for byte, `approvers.tsv`, `start.tsv` and the lease row; a restart with another LAYUP version fails at `version`. |
| `layup run --new`, then `layup run TARGET` | e2e | In CI with no secret: the binary against a local fake forge (the register's `api`) and local bare repositories (`web`), from the command line to the records branch; the table of each run, and the same bytes on a repeat. |
| The demo | uat | On a real GitHub test target, the Operator runs Start and the printed root push, sees each LAYUP write as the App's bot, and reads the records branch with a plain `git clone` and no tool. |

## The findings of the first pilot

The plan hosts nine findings here ([`docs/plan/README.md`](../plan/README.md#the-findings-of-the-first-pilot-that-are-not-fixed)).
`M2a` changes none of them: F-2, F-6, F-13, F-16, F-19, F-31 and F-35 are in
`M2d`, F-21 in `M2c`, and F-37 in `M2e` and `M2g`, each with its reason in the
plan.
