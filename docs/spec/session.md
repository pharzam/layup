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
the prompt's text and the admitted pair. In this order: the version check, the
context check and the rule-file check (each can refuse the start); the start
row, pushed; the process, under its limits; its end; and for a task session the
checks before a push, the push and the bind ([REQ-003](#req-003--before-a-push)).

### The session directory

`DIR/sessions/<session>/`, under the host directory of `--host`
([`run.md`](run.md#the-command)), named by the session ID, never reused:

- `repo/`: the session's clone of the target, made with `CloneLocal` from
  `layup run`'s own clone, then the branch `task/<task>/<attempt>` at the base
  commit (`SwitchCreate`), then its remote removed ([the calls of
  `M2b`](packages.md#the-calls-of-internalgit)). **Decided here:** the clone
  takes only the default branch, so it holds no records branch, and it has no
  remote, so no command of the session reaches a remote through it.
- `home/`: empty at the start; the session's `HOME`. It holds `home/.gitconfig`
  (the author `layup session <session>`, the e-mail
  `<session>@sessions.layup.invalid`; `.invalid` is reserved, RFC 2606, so no
  mail goes anywhere) and a credential of the route `file:`.
- `tmp/`: empty at the start; the session's `TMPDIR`.
- `prompt.md`: the prompt, the caller's text (the fixed text of [the
  probe](#the-probe), or the task of a session).
- `result/`: where the session writes its result file.
- `stdout`, `stderr`: the two outputs of the process.

**Decided here:** `layup run` removes the directory after the session's records
are pushed, so the clone, the home and a copied credential are thrown away (a
session is not resumed across attempts: ADR-0015 decision 5, FT6); a directory
left by a run that stopped is removed before the next session starts.

### The environment and the harness credential

The process gets these variables and no other (the named list of §4; **decided
here** the values):

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
harness itself. Example: Claude Code's `WebFetch` reads a page with
`claude-haiku-5-5`, a model not to use, unless
`ANTHROPIC_DEFAULT_HAIKU_MODEL` names another; a row of Claude Code holds
`ANTHROPIC_DEFAULT_HAIKU_MODEL=<its model>` (seen in the runs of `T-ywk7`,
[`runs/T-ywk7/search/summary.md`](../../runs/T-ywk7/search/summary.md)).

### Rules only from the target

Before the start (ADR-0015 decision 4), code looks in each directory from the
session directory up to `/` for each rule-file name of the row's `rules`; a
file found refuses the start (event `refused`, with the path), as the harness
would load it. Each path of the row's `policy` that exists is recorded in the
start row (`policy`), and the start goes on (L-A5). The rule files inside
`repo/` are the target's own: `AGENTS.md` and the harness's entry file that
points to it (Invariant 9). **Known limit:** a file that a harness loads from a
place that neither list names is not seen here; [the probe](#the-probe) asks
the harness which files it loaded.

### The start of a session

1. **The version.** Code runs the row's `version` command, with the
   environment above and no credential, and takes the first line of its
   standard output, with no space at either end, as the version. A version that
   `records:harnesses.tsv` holds with a passed probe goes on; any other runs
   [the probe](#the-probe) first, and a probe that fails refuses the start. A
   version command that exits non-zero or prints nothing refuses it
   (`version`). No model is called (NEEDLE, mco: the search).
2. **The ID**: `S-` and 8 random lowercase hexadecimal characters, new for each
   session ([`records.md`](records.md#req-011--the-telemetry-record)).
3. **[The context check](#the-context-of-a-start)** and **[the rule-file
   check](#rules-only-from-the-target)**.
4. **The start row**: a row of `sessions.tsv`, and for a task session the event
   `session` of the task, committed and pushed before the process starts
   (**decided here**, survey row 5: a crash leaves a start row with no end,
   which the next run finds).
5. **The process**: the row's `command`, its words split at each space, with
   `{model}`, `{cap}` and `{prompt}` replaced, started in `repo/`, in a process
   group of its own (**decided here**, so that the stop reaches each child). The
   row's `prompt` says how the prompt goes: `file`, `{prompt}` is the path of
   `prompt.md`; `arg`, `{prompt}` is its text as one word, and a text over
   131,072 bytes refuses the start (Linux's limit of one argument); `stdin`,
   the file is the standard input. With `file` and `arg` the standard input is
   empty (survey row 6: `codex exec` waits for its end). `stdout` and `stderr`
   receive the two outputs.

### The context of a start

As §9: the bytes of `prompt.md`, over four and rounded up (**decided here**),
are the estimate; `context` of the model's row in `models.tsv` is the size. An
estimate over the size refuses the start (event `refused`, both numbers); a
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
  64 MiB, and a line of `stdout` at 8 MiB; at either cap comes the stop. Reason:
  the two longest Claude Code sessions of `T-ywk7` printed 1.2 MB and 6.1 MB,
  with no line over 134 KB, so a cap is far above a real session and still keeps
  a harness that prints with no end from filling the host.
- **The stop** (survey row 1): `SIGINT` to the process group, `SIGTERM` 10 s
  later, `SIGKILL` 10 s after that (**decided here**: the two waits).

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

A retry is not `M2b`'s: a new attempt is the task loop's (`M2e`), and its limit
the stall procedure's (`M3d`).

## REQ-003 — Before a push

Requirement: "The rules and the gates of a target are protected from the agents
that they govern: an agent cannot change a rule path without a control that the
agents cannot pass by themselves (ADR-0017)." `M2b` gives prevention layer 1 (no
credential: [the environment](#the-environment-and-the-harness-credential)) and
layer 2 (no rule change leaves the host in a task); layer 3, the forge's rules,
is `M2d`'s and `M2f`'s. For a task session whose end is `done`:

### The fetch by SHA

`layup run` fetches `refs/heads/task/<task>/<attempt>` of `repo/` into its own
clone with `FetchLocal`, with `core.hooksPath` set to an empty directory, and
reads only the head SHA of that fetch: a later change in `repo/` cannot change
what is checked and pushed. The head must descend from the base
(`IsAncestor`), else the result is refused (`base`). **Known limit:** the fetch
runs `git upload-pack` in `repo/`, which reads that clone's configuration, a file
that the session can write. git honors `uploadpack.packObjectsHook` only in
protected configuration (git-config(1)); **decided here**, the integration test
of `FetchLocal` writes each key of git's documentation that starts a program
into the session's configuration, and shows that none runs during the fetch, as
the test of `internal/git` does for the host's configuration
(`TestAHostileHostChangesNothing`).

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
starts `## `. A refused diff (`git diff --binary`) goes to the records as
`payloads/<sha256>` with the event `refused`, and the change becomes a proposal
for the next rule batch (`M2f`). **Known limit:** a target whose product holds
shell scripts changes them only in a rule batch.

### The push and the bind

For a result that passes: a records commit adds the event `push` (the SHA and the
branch), which announces the forge write (fencing, [`run.md`](run.md#the-lease-and-fencing));
then `Push` of the SHA to `task/<task>/<attempt>` with the installation token,
never with force; then, after the forge accepts it, a records commit adds the
event `bound`. A refused push binds nothing: the event `result` says
`push-refused`. The session never pushes.

### A comment for a session

**Decided here:** after its records, each probe and each task session posts one
comment on the control issue (`issue.control` of `start.tsv`), whose first line
starts with the session ID and gives its result, for example
`S-1a2b3c4d: probe of claude 2.1.295: passed`. The task loop of `M2e` moves a
task's comments to the task's issue.

## REQ-005 — The result of a session

Requirement: "Information that passes between role agents is a record in the
target that a machine validates against a schema based on [the baseline]
conventions." A session writes `result/result.tsv` in the form of the block
[`result`](records.md#nfr-001--the-records-of-a-session); `layup run` reads it
by the block and its rules, checks each artifact's SHA-256 against the file at
the head (a mismatch refuses it: `artifact`), and commits it byte for byte as
`tasks/<task>/results/<session>.tsv`, with the event `result`. The session ID,
the task, the role, the attempt and the base come from the start row, never
from the file (§3). The questions, the decisions and the lessons of a handoff,
and the transition table, are `M2e`'s.

### The open attempt

A result is refused (`attempt`) unless the attempt of its start row is still the
task's open attempt: no row of `tasks/<task>/events.tsv` after the session's
event `session` is `closed` or `rebased` for that attempt, or `attempt` for
another one (§3). **Decided here:** the check reads the events at the run's own
last pushed records commit, as the run is the one writer.

## REQ-013 — The probe, admission and routing

### The probe

**Decided here** (§9, ADR-0020 decision 3): a probe is a session of the role
`probe`, with a task ID of its own in the target's form (`T-` and four random
characters), attempt 1, the head of the default branch as its base, and the
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
`files`), no `file` row names a path outside the session directory (else
`outside`), and its usage report names no model whose `models.tsv` row has
`use` `no` (else `not-used`). A probe makes no commit and no
push; its records are a row of `records:harnesses.tsv`, its start row and its
telemetry row.

It runs at the step `probe` of `layup run`, in the restart after `lease` and
before `phase` ([`run.md`](run.md#the-restart)), for each harness of the
register whose version has no passed probe, and before any session whose version
check finds none. The step is `done` when each such harness was probed, passed
or failed, with the counts in `detail`; it is `fail` when a probe could not run
its records. **Decided here:** the block `run-steps` is built, so the build task
of the probe adds `probe` to its enum, with its Go schema, in one change
(condition 1 of the plan review of #147).

### Admission

A pair of a harness and a model is admitted when the harness's last probe, at
the version that the version check reads, passed; and the model's row in
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
[`telemetry.tsv`](records.md#req-011--the-telemetry-record) per session, the
probe's included, in the records commit of the session's end, checked by
`CheckTelemetry` before it is written. **Decided here**, the columns:

- `session`, `task`, `role`, `harness`, `model`: the start row's. `requirements`:
  the task's requirement IDs from the task register (`M2e`); `—` before it.
- `billing`: the row's `billing`.
- `start`, `first_output`, `end`: `layup run`'s own times: the start of the
  process, the first byte of `stdout`, the exit or the kill.
- the tokens: from [the usage report](#the-usage-report-of-a-harness).
- the money: `reported` when the harness reports a cost and the billing is
  `api`; else `computed` from the tokens and `host:prices.tsv`: the rows of the
  session's model, one per class, of one `date` (the newest), with `price`
  naming the row of the class `in`; else `unknown`, when a class lacks a row,
  when the rows differ in currency, or when the usage report names more than
  one model. Reason: the block holds one price ID (task `T-tmhw`), and the
  three rows of one read share the harness, the model and the date; a list of
  IDs would change a built block. A model priced only in a unit that is no
  currency (survey row 10: credits) has no row, so its money is `unknown`.

No row holds the prompt's text (survey row 11).

### The usage report of a harness

The row's `usage` names one format of this closed list; a later section adds a
format here first, with its source, as the types are added
([`README.md`](README.md#the-types)):

| Format | What `layup run` reads |
| ------ | ---------------------- |
| `claude-result` | The last JSON object of `stdout` whose `type` is `result` (Claude Code's `--output-format stream-json`, [headless mode](https://code.claude.com/docs/en/headless)). The tokens are the sums over `modelUsage`: `inputTokens` to `tokens_in`, `outputTokens` to `tokens_out`, `cacheReadInputTokens` plus `cacheCreationInputTokens` to `tokens_cache`; `observed` when each model has the four, `partial` when one lacks a field, `unavailable` with no `result` object. The money is `total_cost_usd`, in USD, Claude Code's own estimate. The models are the keys of `modelUsage`. |
| `none` | Nothing: tokens `unavailable` ("the harness reports none"), money `unknown`; the model is the start row's. |

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
| A `credential_to` that is not `var:NAME`, `file:PATH` or `—`, or a `PATH` that is absolute or holds `..` | exit 2 |
| A `vars` name that is a name of the named list, or the credential's | exit 2 |
| A `models.tsv` or `routing.tsv` row of a harness that the register lacks | exit 2 |
| A harness with no model of `use` `yes` | not an error: no pair of it is admitted |
| A version command that fails, or prints nothing | the start is refused (`version`) |
| A result file that is missing, over 1 MiB, malformed, or has two `status` rows | `no-result`; the result is refused |
| A `probe.tsv` with no token, another token, or no `AGENTS.md` row | the probe fails, with the reason |
| `stdout` with no `result` object, or with lines that are not JSON | tokens `unavailable` (the lines that are not JSON are skipped) |
| A task with no event `attempt` for the session's attempt | the start is refused (`attempt`) |
| A records push that is refused | the run stops, as [`run.md`](run.md#the-lease-and-fencing) |

## NFR-005 — No harness in the engine checks

Requirement: "A deterministic check is preferred to an LLM judgement wherever a
rule can be checked mechanically; the engine checks make no model call, and the
`layup` process calls a model only through the smart-if provider." From `M2b`,
`internal/session` starts a harness: the one model process that `layup` starts
(ADR-0015 decision 3). **Decided here**, the check that will replace the
phase-1 import rule ([`gate.md`](gate.md#nfr-005--no-model-call-in-the-engine-checks)):
the build task of `internal/session` adds to `TestPackageRules` a line of
`packages.md`, "The packages of the engine checks:", with one code span per
package (`internal/psb`, `internal/verify`, `internal/gate`; a later engine
check adds its own), and the rule that no package of that line depends on
`internal/session` or on a package of rule 5; and a form of the cell "Starts a
program" for a program that a register row names. Until that task, the
phase-1 import rule stays the check of the phase-1 commands.

## REQ-015 and REQ-017 — The review of the release of M2b

Their criteria: "No code path of LAYUP modifies base LLM weights or trains a
custom foundation model; a code review of each release records it", and "No
code path of LAYUP calls a cloud provider or hosting platform to modify it; a
code review of each release records it". **Decided here:** the demo task
of `M2b` records a code review of its release (the non-test Go files of `cmd/`
and `internal/`, `go.mod` and the files that the binary embeds, at the commit
that the demo runs) by a reviewer of a model that wrote none of it, as task
`T-efmy` did for phase 1, with `runs/T-efmy/release-check.sh` adapted in that
task: its check 4 allows the verbs of `Fetch`, `Push` and `FetchLocal`, and a
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
The learned weights of routing (the learning loop). A harness's hook events and
the stall triggers (`M3d`). A second harness on one task (`M4a`). Telemetry
Completeness and Cost per Requirement (`M4b`). The usage formats of Codex,
Gemini CLI and OpenCode, added when a registered harness needs one.

## The acceptance tests of M2b

| Part | Level | Test |
| ---- | ----- | ---- |
| The environment | unit | The environment holds the named list only (no `GH_TOKEN`, no SSH agent socket, no variable of the host but `PATH`), the credential by `var:` and by `file:`, and the fixed variables. |
| Rules only from the target | integration | In a real temporary tree, a rule-file name of the row above the session directory refuses the start; a policy path that exists is recorded. |
| The start, the limit and the end | integration | With a fake harness program: the start row is pushed before the process starts; a fake that waits for the end of its input ends; a fake that runs past `wall` is stopped (`SIGINT`, `SIGTERM`, `SIGKILL`) with the class `wall`; the output cap; each class of the end. |
| The context of a start | unit | An estimate over the model's context size refuses the start with both numbers. |
| The open attempt | unit | With a stand-in events table, a result whose attempt was closed, replaced or rebased is refused; the attempt and the base come from the start row. |
| Before a push | integration | With the real `git` and a local bare repository: the fetch by SHA with hooks off; a session configuration that holds each key of git's documentation that starts a program runs none of them; a head that does not descend from the base is refused; a change of `.github/workflows/` and of a rule path is refused before any push, its diff a payload; added lines in §2 of `docs/guardrails.md` pass; the SHA is bound only after the push is accepted. |
| The probe and admission | unit | With a fake harness: the version check; a probe that passes, and one that fails for each reason; admission by the probe and `use`. |
| The usage report | unit | `claude-result` on two recorded `result` events of Claude Code 2.1.295 (one with subagents and a second model) sums `modelUsage`; `none` gives `unavailable` and `unknown`. |
| The writer | unit | One row per session that `CheckTelemetry` passes; money `reported`, `computed` or `unknown` by the billing, the prices and the models. |
| The records of a session | unit | Each new record kind is refused with a wrong header; its Go schema equals its block (`tsv.Compare`), and the owner moves its name from `notYetBuilt` to `built`. |
| The step `probe` | e2e | In CI with no secret: the binary, a local fake forge and a scripted harness: Start, then a restart whose step `probe` probes the harness and writes its rows; the same bytes on a repeat in a new world, as the first run changes the records ([`README.md`](README.md#commands), Determinism). |
| The demo | uat | With the Operator's credentials: the probe of each registered harness (at least two), then one developer session from the uat test whose commit lands on `task/<task>/1`, with its telemetry row; and the review of the release (`REQ-015`, `REQ-017`). |
