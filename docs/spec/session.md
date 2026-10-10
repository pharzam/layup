# `layup run`: role sessions on registered harnesses

Milestone `M2b` of the [implementation plan](../plan/README.md#milestones):
`layup run` probes each registered harness, starts a role session on an
admitted pair of a harness and a model, checks what the session hands back
before anything leaves the host, pushes the session's commit, and writes one
ledger row per session. Requirements: `REQ-013`, `REQ-003`, `REQ-005`,
`REQ-011`, `NFR-005`, `NFR-001`. It derives from
[`architecture.md`](../architecture.md) §3 (one writer), §4 (role sessions),
§6 (rule protection, prevention layers 1 and 2), §9 (the harness register, the
probe, routing, the context of a session) and §12 (the ledger), and from
ADR-0014, ADR-0015, ADR-0017, ADR-0020 and ADR-0024. Conventions:
[`README.md`](README.md). The records:
[`records.md`](records.md#nfr-001--the-records-of-a-session). Task `T-ywk7`
(#147) wrote it, after a public-solution search and the survey of `T-w73g`
([`runs/T-ywk7/`](../../runs/T-ywk7/)); "survey row *n*" is row *n* of the
section `M2b` of that survey, as
[`survey-m2b.tsv`](../../runs/T-ywk7/survey-m2b.tsv) takes or rejects it.

**Decided by the Operator** (O-180, #147): `layup run` of `M2b` runs the probe
only. `internal/session` gives the call that starts a session for a task and
an attempt, which the task loop of `M2e` calls for each step; the developer
session of the demo runs from the uat test, with a small fixed task that the
test gives, and `layup`'s own code writes its records.

## REQ-013 — A role session

Requirement: "The rules, the context and the task state of a target stay in a
form that belongs to no harness agent, so that a second harness agent can do
the work and can check the work of the first." `M2b` gives a session that each
harness of the register runs from the same files; a second harness on one task
is `M4a`'s.

A session is started only by `layup run` (ADR-0015 decision 3), through one
call of `internal/session`, which takes the target, the task and its attempt,
the role, the base commit, the records commit that the prompt was built from,
the prompt's text and the admitted pair. In this order: the session ID; the
attempt check, the version check, the context check, the size check of the
prompt and the rule-file check (each can refuse the start); the start row, pushed; the process, under its limits; its end; and for a task session the
checks before a push, the push and the bind ([REQ-003](#req-003--before-a-push)).

### The session directory

`DIR/sessions/<session>/`, under the host directory of `--host`
([`run.md`](run.md#the-command)), named by the session ID, never reused:

- `target`: the target's `OWNER/NAME`, written first (**decided here**, task
  `T-vxdg`: it lets a run tell its own leftover directories from another
  target's, below).
- `repo/`: the session's clone of the target, made from `layup run`'s own
  clone with `CloneLocal`, which removes its remote before it returns, then the
  branch `task/<task>/<attempt>` at the base commit (`SwitchCreate`), with the
  files of the base in its work tree ([the calls of
  `M2b`](packages.md#the-calls-of-internalgit)). **Decided here:** the clone
  takes only the default branch, so it holds no records branch, and it has no
  remote, so no command of the session reaches a remote through it.
- `home/`: the session's `HOME`; at the start it holds `home/.gitconfig`
  (the author `layup session <session>`, the e-mail
  `<session>@sessions.layup.invalid`; `.invalid` is reserved, RFC 2606, so no
  mail goes anywhere) and a credential of the route `file:`, and no other file.
- `tmp/`: empty when the directory is made; the session's `TMPDIR`, and the
  `HOME` of the version check (step 2), so it may hold what the version
  command wrote when the process starts.
- `prompt.md`: the prompt, the caller's text (the fixed text of [the
  probe](#the-probe), or the task of a session).
- `result/`: where the session writes its result file.
- `stdout`, `stderr`: the two outputs of the process.

**Decided here:** `layup run` removes the directory after the session's records
are pushed (**decided here**, task `T-fsjp`: once its process has started, after
the session's call ends, whether its end and its push succeed or fail, so a
copied credential never waits for the next sweep), so the clone, the home and a copied credential are thrown away (a
session is not resumed across attempts: ADR-0015 decision 5, FT6); a directory
left by a run that stopped is removed before the next session starts.
**Decided here** (task `T-vxdg`, condition 1 of the plan review of #159): a run
removes only the directories of its own target (its file `target`), which the
target's lease guards, and a directory with no file `target`, which only a start
that stopped before its first write leaves; a directory of another target on the
same host is kept, as its own run may be live, and an entry that is not a
directory is skipped. **Known limit:** a start of another target that is between
its first two writes (its directory, then its file `target`) is removed, and
fails at once at its next write.

### The environment and the harness credential

The process gets these variables and no other (the named list of §4, with
`TMPDIR` and the fixed variables added; **decided here** the added names and the
values):

| Variable | Value |
| -------- | ----- |
| `PATH` | the host's `PATH`, which the harness needs to find its own tools |
| `LANG` | `C.UTF-8` |
| `HOME` | the session's `home/` |
| `TMPDIR` | the session's `tmp/` |
| `GIT_CONFIG_NOSYSTEM` | `1` |
| the credential's variable | for the route `var:NAME`, below |
| each fixed variable of the register row | its `NAME=VALUE` (`vars`), recorded in the start row |

No `GH_TOKEN`, no SSH agent socket, no forge credential and no other variable of
the host. **Decided here:** reading the host's `PATH` is the second read of an
environment variable that `TestInputRule` allows, after `environ` of
`internal/git`; the function that builds a session's environment reads only
`PATH` ([`README.md`](README.md#commands), Arguments). Reason: `PATH` is not an
input of `layup`'s result, and a harness that cannot find its tools cannot work.

**The credential** (**decided here**, K40). A register row names a credential
file on the host (`credential`: an absolute path, mode 0600, owned by the user
that runs `layup run`) and how the harness takes it (`credential_to`):
`var:NAME`, where the variable `NAME` is the file's content without its final
line feed; or `file:PATH`, where the file is copied to `home/PATH`, mode 0600.
**Decided here** (task `T-ysph`): `NAME`, and the name of each fixed variable,
is of the form `[A-Za-z_][A-Za-z0-9_]*`; `NAME` takes the prohibitions of a
`vars` name below (a `var:PATH` would take the place of the host's `PATH`); and
`PATH` is not empty, relative, with no part `..`, so `file:a..b` passes and `file:a/../b`
is refused.
Reason: harnesses take a key in a variable (`ANTHROPIC_API_KEY`,
`CODEX_API_KEY` and others,
`ruvnet-brain:tri-smart-skill/tri-smart/scripts/review.mjs:20-25` of the
survey) or a login file in their own directory (Devin's `credentials.toml`), and
the architecture names both forms (§4). The value comes from a file, so `layup`
reads no environment variable for it; and as the environment is the named list,
a variable of the host that would move the billing never reaches the session
(survey row 3). **Known limits:** a
session can read its own credential (L-A1); a copied login whose refresh token
rotates can be spent by the session, so that the host's own copy stops working
(found while `T-ywk7` set up a searcher; not tested).

**The fixed variables.** A row's `vars` give the session settings of the
harness itself, never a credential. **Decided here:** the prohibitions above
hold for `vars` too: the reader of the register refuses a `vars` name of the
named list, the credential's name, and each name of a forge credential or an
agent socket that `layup` knows (`GH_TOKEN`, `GITHUB_TOKEN`, `GH_ENTERPRISE_TOKEN`,
`GITHUB_ENTERPRISE_TOKEN`, `SSH_AUTH_SOCK`); a credential under another name is
the Operator's error, which no reader can see. Example: Claude Code's `WebFetch` reads a page with
`claude-haiku-5-5`, a model not to use, unless
`ANTHROPIC_DEFAULT_HAIKU_MODEL` names another; a row of Claude Code holds
`ANTHROPIC_DEFAULT_HAIKU_MODEL=<its model>` (seen in the runs of `T-ywk7`,
[`runs/T-ywk7/search/summary.md`](../../runs/T-ywk7/search/summary.md)).

### Rules only from the target

Before the start (ADR-0015 decision 4), code looks in each directory from the
session directory up to `/` for each rule-file name of the row's `rules`; an
entry of that name found (a file, a directory or a link) refuses the start (event `refused`, `rules` and the path), as the harness
would load it. Each path of the row's `policy` that exists is recorded in the
start row (`policy`), and the start goes on (L-A5). The rule files inside
`repo/` are the target's own: `AGENTS.md` and the harness's entry file that
points to it (Invariant 9). **Known limit:** a file that a harness loads from a
place that neither list names is not seen here; [the probe](#the-probe) asks
the harness which files it loaded.

### The start of a session

**The attempt** (**decided here**): `internal/run` gives the call that starts an
attempt of a task, a records commit with the event `attempt`, its number and its
base. The task loop of `M2e` calls it; in `M2b` the uat test calls it for the
task of the demo's developer session. A session's start needs the event of its
attempt ([Input states](#input-states)). **Decided here** (task `T-d8t9`,
condition 2 of the plan review of #164): the number is one more than the
largest attempt of an event `attempt` of the task, 1 for none, as an event of
another kind names an attempt that one of `attempt` started; a records table
that the records branch does not hold yet (`events.tsv` of a new task,
`sessions.tsv`) is an empty table.

1. **The ID**: `S-` and 8 random lowercase hexadecimal characters, new for each
   session ([`records.md`](records.md#req-011--the-telemetry-record)), drawn
   first, so that each refusal below names its session. **Decided here**
   (task `T-d8t9`, condition 1 of the plan review of #164): before the ID,
   `layup run` sweeps the directories that a stopped run of the target left;
   right after it, it makes the session directory, so that each check runs on
   the session's own files; a start that is refused, or that stops on an error
   before its start row is pushed, removes its directory at once, so a copied credential
   does not wait for the next sweep; a directory that the start made in part
   before it failed is removed too, and one of the ID that was there before
   the start is another session's and is kept. Then the attempt
   check: a task session whose attempt has no event `attempt` is refused
   (`attempt`).
2. **The version.** Code runs the row's `version` command, with the
   environment above and no credential (**decided here**, task `T-d8t9`: no
   credential variable, and `HOME` the session's `tmp/`, as `home/` holds a
   credential of `file:` once the directory is made), in `repo/`, and takes the first line of its
   standard output, with no white space at either end (**decided here**, task
   `T-5pxd`: a tab or a carriage return is no part of a version), as the version. A version whose
   last probe passed ([Admission](#admission)) goes on; any other runs
   [the probe](#the-probe) first, and a probe that fails refuses the start
   (`probe`). A probe session's own version check takes the version and runs
   no probe. A
   version command that exits non-zero, or whose first line is empty or white
   space only (**decided here**, task `T-5pxd`: such a line names no version),
   refuses it (`version`). No model is called (NEEDLE, mco: the search).
3. **[The context check](#the-context-of-a-start)**; the size check of the
   prompt: for a row whose `prompt` is `arg`, a `prompt.md` over 131,071 bytes
   (Linux's limit of one argument is 131,072 bytes with its final zero byte) refuses the start (`prompt` and its size);
   and **[the rule-file check](#rules-only-from-the-target)**.
4. **The start row**: a row of `sessions.tsv`, and for a task session the event
   `session` of the task, committed and pushed before the process starts
   (**decided here**, survey row 5: a crash leaves a start row with no end,
   which the next run finds), in one records commit (**decided here**, task
   `T-d8t9`: a crash leaves both or neither).
5. **The process**: the row's `command`, its words split at each space, with
   `{model}`, `{cap}` and `{prompt}` replaced in one pass (**decided here**,
   task `T-5pxd`: a prompt that holds `{model}` stays as it was written), started in `repo/`, in a process
   group of its own (**decided here**, so that the stop reaches each child). The
   row's `prompt` says how the prompt goes: `file`, `{prompt}` is the path of
   `prompt.md`; `arg`, `{prompt}` is its text as one word (its size was checked at
   step 3); `stdin`,
   the file is the standard input, and `{prompt}`, where the command holds
   it, is the path of `prompt.md` (**decided here**, task `T-5pxd`, condition
   3 of the plan review of #161; a `file` or `arg` command with no `{prompt}`
   is #186). With `file` and `arg` the standard input is
   empty (survey row 6: `codex exec` waits for its end). `stdout` and `stderr`
   receive the two outputs. A program that is missing or cannot run starts
   nothing (the class `start`).

**A refused start** (**decided here**): a refusal of a task session's start is
an event `refused` of its task, with its session ID and its reason, and no row
of `sessions.tsv`; a refusal of a probe's start is the probe's row of
`harnesses.tsv`, `failed`, with its reason, as a probe has no events. The
refusal `pair` comes before the call of `internal/session`, so its event has
`—` for the session. A start that is refused runs no process, so it has no
telemetry row.

### The context of a start

As §9: the bytes of `prompt.md`, over four and rounded up (**decided here**),
are the estimate; `context` of the model's row in `models.tsv` is the size. An
estimate over the size refuses the start (event `refused`, `context` and both
numbers, the estimate then the size: **decided here**, task `T-6sbe`); a
session that starts has both numbers in its start row.

### The limit of a session

- **Wall clock:** `wall` minutes of the row, timed by `layup run` from the
  start of the process (**decided here**: no harness has a time limit of its
  own, so `layup run` owns the timer; the search, searcher B). At the limit
  comes the stop.
- **Spend cap:** `cap` of the row, given to the harness by `{cap}` in its
  command (for example `--max-budget-usd` of Claude Code); a row with `cap` `—`
  runs under its wall-clock limit only, and its spend is an unknown part until
  the session ends (§12, ADR-0024 decision 5).
- **Output cap** (**decided here**): `stdout` and `stderr` are each cut at
  64 MiB, and a line of `stdout` (its bytes before its line feed) at 8 MiB;
  at either cap comes the stop, and the rest of the output is read and
  dropped (**decided here**, task `T-5pxd`: a reader that stops reading would
  let the program block on a full pipe before the stop reaches it). Reason:
  the two longest Claude Code sessions of `T-ywk7` printed 1.2 MB and 6.1 MB,
  with no line over 134 KB ([`test-runs.md`](../../runs/T-ywk7/test-runs.md),
  Measurements), so a cap is far above a real session and still keeps a harness
  that prints with no end from filling the host.
- **The stop** (survey row 1): `SIGINT` to the process group, `SIGTERM` 10 s
  later, `SIGKILL` 10 s after that (**decided here**: the two waits), each
  only while the process has not ended. **Decided here** (task `T-5pxd`): the
  process has ended when its first program has exited and both outputs are
  closed, so a child that holds an output runs under the same limits; after
  `SIGKILL` and the exit of the first program, the outputs are closed, so a
  program that left the group cannot hold the end. The first program is
  reaped only at the end, so the group's ID, its process ID, is taken by no
  other process while a signal can reach it.

### The end of a session

When the process has exited or been killed, code reads the result file
(`result/result.tsv` of a task session, [`result`](records.md#nfr-001--the-records-of-a-session);
`result/probe.tsv` of a probe) and the usage report in `stdout`
([REQ-011](#req-011--the-writer-of-the-telemetry-record)). A session that
exits with no valid result file has failed (§4). **Decided here**, the class of
the end (survey row 2), in the `detail` of the event `result`, or in the probe's
row:

| Class | When | Again? |
| ----- | ---- | ------ |
| `done` | exit 0 and a valid result file | — |
| `start` | the program is missing or cannot run | fails the same way again |
| `crash` | a non-zero exit, or a signal that `layup run` did not send | uncertain |
| `no-result` | exit 0 with no valid result file | uncertain |
| `wall`, `output` | stopped at the wall-clock limit or an output cap | uncertain |

**Decided here** (task `T-bpxg`, condition 3 of the plan review of #162): a
process may meet two rows, so the class is the first that holds in this order:
`start`; `wall` or `output`, a stop that `layup run` sent, whatever the exit and
the result file; `crash`; `no-result`; `done`. A process that traps `SIGINT`,
writes a valid result file and exits 0 after the stop at `wall` is `wall`. Every
other error of the process call, at the start or after it (an output file
that `layup run` cannot open, a write of an output that fails), is `layup
run`'s own and no class: the run reports it as its own. **Decided here** (task
`T-bpxg`): the result file is opened without following a link, and a file
that is a link or is not a regular file is no valid result file, as for
`repo/.git` (a session could link to any file of the host's user).

A retry is not `M2b`'s: a new attempt is the task loop's (`M2e`), and its limit
the stall procedure's (`M3d`).

## REQ-003 — Before a push

Requirement: "The rules and the gates of a target are protected from the agents
that they govern: an agent cannot change a rule path without a control that the
agents cannot pass by themselves (ADR-0017)." `M2b` gives prevention layer 1 (no
credential: [the environment](#the-environment-and-the-harness-credential)) and
layer 2 (no rule change leaves the host in a task); layer 3, the forge's rules,
is `M2d`'s and `M2f`'s. For a task session whose end is `done`: its event
`result` is `done` and it has no event `refused` (**decided here**, task
`T-z027`, condition 3 of the plan review of #166: a `done` session refused at
its end, `artifact` among them, is never checked or pushed). The checks run on
the head that the end fetched into `layup run`'s clone
([REQ-005](#req-005--the-result-of-a-session)). Each refusal below is one
records commit through fencing, and nothing is pushed:

### The fetch by SHA

**Decided here**, so that no configuration of the session's clone is read
(§4): no command of `layup` runs in `repo/` after the session starts. Code reads
the head SHA of `refs/heads/task/<task>/<attempt>` from the files of
`repo/.git`: the loose ref, a regular file of 40 lowercase hexadecimal
characters and a line feed, else the line of `packed-refs` whose ref field,
the text after the first space, is the ref, and whose first field is 40
lowercase hexadecimal characters (**decided here**, task `T-6sbe`; a line that
names the ref in another form refuses).
A `repo/.git` that is not a directory (a link to one is refused too); a loose ref that is a link, a symbolic
ref or a file of any other content (`packed-refs` is then not read); no such
ref; or a SHA that names no commit of the session's objects, which
`FetchSession` checks before it sets a ref: each refuses the result (`branch`).
`FetchSession` then makes a
scratch bare repository whose one alternate is the object directory of
`repo/.git`, sets a ref there to that SHA, fetches it into `layup run`'s own
clone with `core.hooksPath` set to an empty directory, and removes the scratch
repository: `git upload-pack` runs in the scratch repository, whose
configuration is `layup`'s, and reads the session's objects as data. The head
must descend from the base (`IsAncestor`), else the result is refused (`base`).
Each later check reads only that SHA in `layup run`'s clone, so a later change in
`repo/` cannot change what is checked and pushed. The integration test of
`FetchSession` writes into `repo/.git/config` each key that `git help --config` of the host's
`git` lists and that starts a program (task `T-z5dj`, O-188 of #157), and fails when one runs, as the test of `internal/git`
does for the host's configuration (`TestAHostileHostChangesNothing`).

### A workflow or rule-path change

`DiffNames` from the base to the head. A path under `.github/workflows/`
refuses the result (`workflow`; the App has no workflows permission, O-92). A
path that the rule-path register matches refuses it (`rule-path`; ADR-0017
layer 2). The register is `rule-paths.tsv` of the session's records commit
(`internal/rules`; [`setup.md`](setup.md#the-rule-path-register)); **decided
here**, the match: a pattern that ends with `/` matches each path under it,
another pattern the path itself, and any path whose name ends with `.sh`
matches, as the register's entry "each file of the tree whose name ends with
`.sh`" names a kind of file, not the files of one tree. The exception of
`docs/guardrails.md`: its diff (`DiffFile`) removes no line, and each added line
lies at the head between the line that starts `## 2.` and the next line that
starts `## `. **Decided here** (task `T-m1dx`): the exception holds only for a
diff read as hunks with at least one added line; a diff with no hunk, a mode
change, a binary diff, or other text is a rule-path change; "between" excludes both boundary
lines, so an added line that starts `## ` is refused and one under a `### `
sub-heading passes; the exception is that one file's, as the register's block
allows no other (`internal/rules` refuses another). A refused diff (`git diff --binary`) goes to the records as
`payloads/<sha256>` with the event `refused`, and the change becomes a proposal
for the next rule batch (`M2f`): the payload is the proposal, and `M2b` writes
no other record of it (**decided here**, task `T-z027`). The `detail` of the
event is the reason, one space, and the payload's SHA-256, for example
`rule-path 3f2a…` (**decided here**, condition 3 of the plan review of #166).
**Decided here** (condition 2): a records commit of the session with no
`rule-paths.tsv`, or one that its reader refuses, refuses the result
(`rule-paths`), fail closed (a records commit or a `git` that cannot be read
is the run's own error; round 1 of #166), as ADR-0017 decision 2 makes the register the
guard of layer 2; a target that Start made has no register until `layup
setup` writes it. **Known limit:** a target whose product holds
shell scripts changes them only in a rule batch.

### The push and the bind

For a result that passes: a records commit adds the event `push` (the SHA and the
branch), which announces the forge write (fencing, [`run.md`](run.md#the-lease-and-fencing));
then `Push` of the SHA to `task/<task>/<attempt>` with the installation token,
never with force; then, after the forge accepts it, a records commit adds the
event `bound`. A refused push binds nothing: an event `refused` with the detail
`push-refused`. The session never pushes. **Decided here** (task `T-e3sy`,
condition 1 of the plan review of #167): a refused push is one that `Push`
reports as its code 1, not a fast-forward or refused by the remote; a push that
fails otherwise (the remote cannot be reached) is the run's own error, and no
event follows `push`, as the forge's answer is unknown. A `push` with no
`bound` and no `refused` after it, from that error or from a lost lease at the
commit of `bound`, is a forge write whose result is not recorded: the next
reader of that session (the task loop of `M2e`) reads the target's branch
before it acts on it. The push runs from `layup run`'s clone, with the head
that the end fetched.

### A comment for a session

**Decided here:** after its records, each probe and each task session posts one
comment on the control issue (`issue.control` of `start.tsv`), whose first line
starts with the session ID and gives its result, for example
`S-1a2b3c4d: probe of claude 2.1.295: passed`. **Decided here** (task
`T-fsjp`): the first line of a task session's comment is `<session>: <role>
session of <task>, attempt <n>: <class>`, with `, refused <reason>` after it for
each refusal of its result, for example `S-1a2b3c4d: developer session of
T-ab12, attempt 1: done, refused artifact`. It is posted after every record of
the session, those of the checks before a push and of the push too, from what
the events hold, so it names the refusals of the end and of those checks
(**decided here**, round 1 of #165); a session whose end recorded nothing has no
comment, and a lost lease stops the run before it (round 2 of #165). The call is the forge's `Comment`
([`forge.md`](forge.md#the-calls-of-m2b), O-189). The task loop of `M2e` moves
a task's comments to the task's issue.

## REQ-005 — The result of a session

Requirement: "Information that passes between role agents is a record in the
target that a machine validates against a schema based on [the baseline]
conventions." A session writes `result/result.tsv` in the form of the block
[`result`](records.md#nfr-001--the-records-of-a-session); `layup run` reads it
by the block and its rules, checks each artifact's SHA-256 against the file at
the head (a mismatch adds the event `refused`, `artifact`), and commits it byte
for byte as
`tasks/<task>/results/<session>.tsv`, with the event `result`. The session ID,
the task, the role, the attempt and the base come from the start row, never
from the file (§3). The questions, the decisions and the lessons of a handoff,
and the transition table, are `M2e`'s.

**Decided here** (task `T-fsjp`, condition 2 of the plan review of #165): the
result file is read for every exit, to tell `done` from `no-result`, and is
checked against the head and committed for the class `done` only, as only a
session that ends `done` hands its work over; for another class the events and
the telemetry row are committed, and no `results/` file. The file is read once,
and the bytes committed are the bytes checked. A link as the last part of its
path is refused (no file of the host is followed); a link above it, `result/`
itself, is followed, and the file still has to pass its block, so it is one the
session could have written. **Decided here** (condition 1 of the plan review
of #165): the head of the session is read from the files of `repo/.git`
([The fetch by SHA](#the-fetch-by-sha)) and fetched into `layup run`'s clone
with `FetchSession` at the end of the session, and each artifact is read there,
so no command of `layup` runs in `repo/`; a head that cannot be read or fetched
adds the event `refused`, `branch`, and the result is not committed. An
artifact whose SHA-256 differs from its file at the head, or whose path the
head lacks, adds the event `refused`, `artifact`, and the result is still
committed (**decided here**: "adds"). The result file, the events `result` and
`refused`, and the telemetry row are one records commit.

### The open attempt

A result is refused (the event `refused`, `closed-attempt`) unless the attempt of its
start row is still the
task's open attempt: no row of `tasks/<task>/events.tsv` after the session's
event `session` is `closed` or `rebased` for that attempt, or `attempt` for
another one (§3). **Decided here:** the check reads the events at the run's own
last pushed records commit, as the run is the one writer. **Decided here**
(task `T-fsjp`): a refused result is not committed as `results/<session>.tsv`;
the event `result` with the class, the event `refused` and the telemetry row
are.

## REQ-013 — The probe, admission and routing

### The probe

**Decided here** (§9, ADR-0020 decision 3): a probe is a session of the role
`probe`, with a task ID of its own in the target's form (`T-` and four random
characters, drawn again until no row of `sessions.tsv` holds it; the task
register of `M2e` draws its IDs by the same rule), attempt 1, the head of the default branch as its base, and the
first model of the harness, in the order of `models.tsv`, whose row has `use`
`yes`. Its prompt,
with `{result}` the absolute path of `result/probe.tsv` and `{token}` 16 random
hexadecimal characters:

```text probe-prompt
This is a probe of LAYUP. Change no file in this repository and make no commit.

Write one file, {result}, of tab-separated values, with the header line
"kind<TAB>value" and these rows: one row "token<TAB>{token}"; and one row
"file<TAB><path>" for each instruction file that you loaded for this session,
with its path relative to this repository, or absolute when it is outside it.
Then end.
```

It passes when its end is `done`, `probe.tsv` passes the block
[`probe-result`](records.md#nfr-001--the-records-of-a-session), its token is
the prompt's (else `token`), a `file` row has the value `AGENTS.md` (else
`files`), no `file` row names a path outside the session directory other than
a `policy` path of the row (else `outside`), and its usage report names no model
whose `models.tsv` row has `use` `no` (else `not-used`). **Known limit:** the
`file` rows are the model's answer, as no output of a harness that this
specification reads lists the files it loaded; a model that leaves out a file
passes. A probe makes no commit and no
push; its records are a row of `records:harnesses.tsv`, its start row and its
telemetry row, or the row of `harnesses.tsv` alone when its start was refused.

It runs at the step `probe` of `layup run`, in the restart after `lease` and
before `phase` ([`run.md`](run.md#the-restart)), and before any session whose
version's last probe did not pass ([Admission](#admission)). **Decided here:** the step first skips each
harness of the register with no model of `use` `yes`, whatever its version, and
runs no command of it, as a probe needs a model; it gets no row, and the step
counts it as skipped. Then it probes each other harness whose last probe at
its version did not pass, or that has no probe at its version. The step is `done` when each harness that needs a probe was probed, and each
with no model skipped, with
three counts in `detail`: probed and passed, probed and failed, skipped; a
harness whose last probe at its version passed is in none; it is `fail` when a probe could not run
its records. **Decided here:** the block `run-steps` is built, so the build task
of the probe adds `probe` to its enum, with its Go schema, in one change
(condition 1 of the plan review of #147).

### Admission

A pair of a harness and a model is admitted when the harness's last probe, at
the version that the version check reads, passed (its last row of
`records:harnesses.tsv` at that version, in the order of the file, the one
home of this rule: task `T-cht1`); and the model's row in
`models.tsv` has `use` `yes`. A model with `use` `no` is on the "not used"
list (§9), which admission reads at each start, so a change of the list takes
effect at the next start, and each start row records the model that it
admitted. Code alone admits; no model is asked.

### The routing register

`records:routing.tsv` holds, for each role and tier, the ordered list of pairs
(K38: the order only; the weights of §13 come with the learning loop).
**Decided here:** its source is the Operator's `host:registers/routing.tsv`,
which `layup run` copies into the records at the step `probe` when the two
differ, as each value that a run uses is copied into the records (§1). The
first admitted pair of the role's list for the task's tier is the session's
pair; with none, the start is refused (`pair`).

## REQ-011 — The writer of the telemetry record

Requirement: "Every task has a record of its token count, its latency and its
wall-clock duration in the target." The writer is `internal/ledger`: one row of
[`telemetry.tsv`](records.md#req-011--the-telemetry-record) per session that
started, the probe's included, in the records commit of the session's end, checked by
`CheckTelemetry` before it is written. **Decided here**, the columns:

- `session`, `task`, `role`, `harness`, `model`: the start row's. `requirements`:
  the task's requirement IDs from the task register (`M2e`); `—` before it.
- `billing`: the row's `billing`.
- `start`, `first_output`, `end`: `layup run`'s own times: the start of the
  process, the first byte of `stdout`, and its end (**decided here**, task
  `T-fsjp`: a program that cannot start has no end of its own, so its `end`
  is its `start`), as the column `end` of
  [`telemetry.tsv`](records.md#req-011--the-telemetry-record) says.
- the tokens: from [the usage report](#the-usage-report-of-a-harness).
- the money: `reported` when the harness reports a cost and the billing is
  `api`; else `computed` from the tokens and `host:prices.tsv`: the rows of the
  session's harness and model, one per class, of one `date`, the newest `date`
  among those rows, with `price` naming the row of the class `in`; else
  `unknown`, when a class lacks a row, when one class has two rows of that
  date, when the rows differ in currency, when the usage report names more
  than one model, or (**decided here**, task `T-4c3q`) when it names one model
  that is not the start row's (the price would be of a model the session did
  not run), or when the tokens are not `observed` (a class with no count would
  sum as zero, the cost that FT2 forbids). Reason: the block holds one price ID (task `T-tmhw`), and the
  three rows of one read share the harness, the model and the date; a list of
  IDs would change a built block. A model priced only in a unit that is no
  currency (survey row 10: credits) has no row, so its money is `unknown`.
  **Decided here** (task `T-4c3q`): a
  computed money is the exact decimal of the sum, with no trailing zero and one
  decimal digit at least, computed with no float, so two writers give the same
  bytes; a reported cost is the text of the JSON number of the report, written
  in that form (`6` as `6.0`, `1e-7` as `0.0000001`).

No row holds the prompt's text (survey row 11).

### The usage report of a harness

The row's `usage` names one format of this closed list; a later section adds a
format here first, with its source, as the types are added
([`README.md`](README.md#the-types)):

| Format | What `layup run` reads |
| ------ | ---------------------- |
| `claude-result` | The last JSON object of `stdout` whose `type` is `result` (Claude Code's `--output-format stream-json`, [headless mode](https://code.claude.com/docs/en/headless)). The tokens are the sums over `modelUsage`: `inputTokens` to `tokens_in`, `outputTokens` to `tokens_out`, `cacheReadInputTokens` plus `cacheCreationInputTokens` to `tokens_cache`; `observed` when each model has the four, `partial` when one lacks a field and a class is still summed (below), `unavailable` with no `result` object. The money is `total_cost_usd`, in USD, Claude Code's own estimate. The models are the keys of `modelUsage`, sorted (**decided here**, task `T-bpxg`). |
| `none` | Nothing: tokens `unavailable` ("the harness reports none"), money `unknown`; the model is the start row's. |

**Decided here** (task `T-bpxg`, condition 1 of the plan review of #162): a
class of tokens is summed only when every model of `modelUsage` gives each of
its fields as an integer, else it is `—`, as a sum that leaves a model out would
undercount (FT2). Three classes summed are `observed`; one or two are
`partial`; none is `unavailable`. The reason names each model and the fields it
lacks (`m has no outputTokens`). A `result` object whose `modelUsage` is
missing or empty is `unavailable` ("the report names no model"), never 0.

**Decided here:** the tokens come from `modelUsage`, not from `usage`. Reason:
in two runs of `T-ywk7` on Claude Code 2.1.295, `usage` equals `modelUsage` in
a session with no subagent and no web tool, and covers only the main loop in a
session with five subagents and `WebFetch` (53,309 output tokens against
911,458 in `modelUsage`), whose costs add up to `total_cost_usd`. Devin's `-p`
gives no usage report (3000.11.3), so its rows use `none`.

## Input states

**Decided here**, for each state that the happy path never makes:

| State | Result |
| ----- | ------ |
| A host register (`harnesses.tsv`, `models.tsv`, `routing.tsv`) with two rows for a key, or a field that its type refuses | exit 2: its reader names the line |
| A harness row with `cap` and no `{cap}` in `command`, or `{cap}` and `cap` `—` | exit 2 |
| A `credential` that is not absolute, is missing, is not mode 0600 or has another owner | exit 2, naming the file |
| A `credential_to` that is not `var:NAME`, `file:PATH` or `—`, a `NAME` of another form or of [the list above](#the-environment-and-the-harness-credential), or a `PATH` that is empty, absolute or has a part `..` | exit 2 |
| A `vars` name of the named list, the credential's, or a name of a forge credential or an agent socket of [the list above](#the-environment-and-the-harness-credential) | exit 2 |
| A host directory with no `registers/models.tsv` or no `registers/routing.tsv` | exit 2, "a missing or unreadable file" ([`README.md`](README.md#commands)) |
| A `models.tsv` or `routing.tsv` row of a harness that the register lacks | exit 2 |
| A harness with no model of `use` `yes`, or with no row of `models.tsv` | not an error: the step `probe` skips it, and no pair of it is admitted |
| A version command that fails, or whose first line is empty or white space only | the start is refused (`version`) |
| A result file that is missing, a link or not a regular file, over 1 MiB, malformed, or has two `status` rows | `no-result`; the result is refused |
| A `probe.tsv` with no token, another token, or no `AGENTS.md` row | the probe fails, with the reason |
| A `stdout` that cannot be read | tokens `unavailable`, "stdout cannot be read" (**decided here**, task `T-fsjp`); a usage format of no list is the run's own error |
| `stdout` with no `result` object, or with lines that are not JSON | tokens `unavailable` (the lines that are not JSON, or not of the form of an object of the report, are skipped) |
| A task with no event `attempt` for the session's attempt | the start is refused (`attempt`) |
| A `prompt.md` over 131,071 bytes for a row whose `prompt` is `arg` | the start is refused (`prompt`), before the start row |
| A ref of `repo/.git` that is malformed, or whose SHA names no commit of the session's objects | the result is refused (`branch`) |
| A `host:prices.tsv` that is missing | an empty table (**decided here**, task `T-fsjp`): no money is `computed`, so a money that is not `reported` is `unknown` |
| A session's records commit with no `rule-paths.tsv`, or one that its reader refuses | the result is refused (`rule-paths`), and nothing is pushed |
| A push of a session's head that `git` refuses (not a fast-forward, or refused by the remote) | the event `refused`, `push-refused`; nothing is bound |
| A push of a session's head that fails otherwise (the remote cannot be reached) | the run's own error; no event follows `push` |
| A records push that is refused | the run stops, as [`run.md`](run.md#the-lease-and-fencing) |

## NFR-005 — No harness in the engine checks

Requirement: "A deterministic check is preferred to an LLM judgement wherever a
rule can be checked mechanically; the engine checks make no model call, and the
`layup` process calls a model only through the smart-if provider." From `M2b`,
`internal/session` starts a harness: the one model process that `layup` starts
(ADR-0015 decision 3). **Decided here**, the check that is added beside the
phase-1 import rule, which stays ([`gate.md`](gate.md#nfr-005--no-model-call-in-the-engine-checks)):
`TestPackageRules` reads a line of `packages.md`, "The packages of the engine
checks:", with one code span per package (`internal/psb`, `internal/verify`,
`internal/gate`; a later engine check adds its own), and checks that no package
of that line depends on `internal/session` or on a package of rule 5; and a form
of the cell "Starts a program" for a program that a register row names
([`packages.md`](packages.md#the-table-of-m2b); task `T-y10b`, row 29 of [the
plan](../plan/README.md#the-tasks-of-m2b)).

## REQ-015 and REQ-017 — The review of the release of M2b

Their criteria: "No code path of LAYUP modifies base LLM weights or trains a
custom foundation model; a code review of each release records it", and "No
code path of LAYUP calls a cloud provider or hosting platform to modify it; a
code review of each release records it". **Decided here:** a task of
its own, before the demo (row 39a of [the plan](../plan/README.md#the-tasks-of-m2b),
O-187 of #152), records a code review of the release of `M2b` (the non-test Go
files of `cmd/` and `internal/`, `go.mod` and the files that the binary embeds,
the code that the demo runs) by a reviewer of a model that wrote none of it, as
task `T-efmy` did for phase 1, with `runs/T-efmy/release-check.sh` adapted in
that task: its check 4 allows the verbs of `Fetch`, `Push` and `FetchSession`, and a
new list names each program that the code starts (`git`, `sh`, the command of a
harness register row). The review reads the whole release, so it covers the
code of `M2a` too (#148). A session runs under the Operator's user and can reach
what that user can (L-A1); the review records it as that limit.

## Not in M2b

The task loop: the plan session, the plan review, the verifier, the close-out,
the draft pull request, the transition table, and the questions, decisions and
lessons of a handoff (`M2e`). The tier of a task from its plan, and the authors
of a change for a plan review or a verification (`M2e`). The budget check
before each start and `budget.tsv` (`M2c`), and the circuit breaker (`M3d`).
The branch of a rule batch, `layup/rules` and the known-bad patches (`M2f`).
The context of a session by its links, the record kinds of each step's row
(§9, ADR-0020 decision 7): in `M2b` the prompt is the caller's text, and the step
table's record kinds come with the task loop (`M2e`). The learned weights of
routing (the learning loop). A harness's hook events and
the stall triggers (`M3d`). A second harness on one task (`M4a`). Telemetry
Completeness and Cost per Requirement (`M4b`). The usage formats of Codex,
Gemini CLI and OpenCode, added when a registered harness needs one.

## The acceptance tests of M2b

| Part | Level | Test |
| ---- | ----- | ---- |
| The session directory | integration | With the real `git`: `repo/` is a `--no-local` clone of the default branch only, with no remote and no records branch, on `task/<task>/<attempt>` at the base; `home/` holds only `.gitconfig` and a credential of `file:`, and `tmp/` and `result/` start empty; the directory is removed after the session's records are pushed, and one that a stopped run left is removed before the next session. |
| The environment | unit | The environment holds the named list only (no `GH_TOKEN`, no SSH agent socket, no variable of the host but `PATH`), the credential by `var:` and by `file:`, and the fixed variables; a `vars` name of the named list, the credential's or the list of forge credentials and agent sockets is refused (exit 2). |
| Rules only from the target | integration | In a real temporary tree, a rule-file name of the row above the session directory refuses the start; a policy path that exists is recorded. |
| The start, the limit and the end | integration | With a fake harness program: the version check; the start row is pushed before the process starts; a fake that waits for the end of its input ends; a fake that runs past `wall` is stopped (`SIGINT`, `SIGTERM`, `SIGKILL`) with the class `wall`; the output cap; each class of the end. |
| The context of a start | unit | An estimate over the model's context size refuses the start with both numbers. |
| The result of a session | integration | A result file of the block `result` is committed byte for byte with the event `result`; an artifact whose SHA-256 differs at the head adds the event `refused` (`artifact`); the session ID, the task, the role, the attempt and the base come from the start row. |
| The open attempt | unit | With a stand-in events table, a result whose attempt was closed, replaced or rebased is refused (`closed-attempt`); the attempt and the base come from the start row; a session with no event `attempt` of its attempt is refused at its start. |
| Before a push | integration | With the real `git` and a local bare repository: the head read from the files of `repo/.git` (a loose ref, a packed ref with `git pack-refs`; with a stand-in reader of the files, at unit: a link, a symbolic ref, a malformed loose ref, a packed line of another form and no ref refused; a SHA of no commit refused by `FetchSession`); the fetch through the scratch repository with hooks off; a session configuration that holds each key of git's documentation that starts a program runs none of them; a head that does not descend from the base is refused; a change of `.github/workflows/` and of a rule path is refused before any push, its diff a payload; added lines in §2 of `docs/guardrails.md` pass; the SHA is bound only after the push is accepted. |
| The probe and admission | unit | With a fake harness: a probe that passes, one that reports a `policy` path and passes, and one that fails for each reason; admission by the probe and `use`; a harness with no model of `use` `yes` is skipped. |
| A refused start | unit | A refused start of a task session is an event `refused` with its session ID, its reason (`attempt`, `version`, `probe`, `context`, `prompt`, `rules`) and no row of `sessions.tsv` or `telemetry.tsv`; the refusal `pair` has `—` for the session; a probe's refused start is its row of `harnesses.tsv`, `failed`, with `version` `—` when the version check refused it. |
| The routing register | unit | `host:registers/routing.tsv` is copied into `records:routing.tsv` at the step `probe` when the two differ, and not when they are equal; the session's pair is the first admitted pair of the role's list for the task's tier; with none, the start is refused (`pair`). |
| The usage report | unit | `claude-result` on two recorded `result` events of Claude Code 2.1.295 (one with subagents and a second model) sums `modelUsage`; `none` gives `unavailable` and `unknown`. |
| The writer | unit | One row per session that `CheckTelemetry` passes; money `reported`, `computed` or `unknown` by the billing, the prices and the models. |
| The records of a session | integration | Each new record kind is refused with a wrong header; its Go schema equals its block (`tsv.Compare`), and the owner moves its name from `notYetBuilt` to `built`. |
| The step `probe` | e2e | In CI with no secret: the binary, a local fake forge and a scripted harness: Start, then a restart whose step `probe` probes the harness, writes its rows, and posts one comment on the control issue whose first line starts with the probe's session ID; the same bytes on a repeat in a new world, as the first run changes the records ([`README.md`](README.md#commands), Determinism). |
| The review of the release | uat | A recorded code review of the release for `REQ-015` and `REQ-017`, by a reviewer of a model that wrote none of it, with the adapted `release-check.sh`, before the demo (row 39a of [the plan](../plan/README.md#the-tasks-of-m2b), O-187). |
| The demo | uat | With the Operator's credentials: the probe of each registered harness (at least two), then one developer session from the uat test whose commit lands on `task/<task>/1`, with its telemetry row. |
