# LAYUP architecture

LAYUP takes a problem statement from an idea owner and drives the delivery of
that product, from the first question to the last accepted requirement, in a
new repository that it sets up with the pinned baseline
([`setup/armature.pin`](setup/armature.pin)). It is a deterministic
orchestrator of the whole lifecycle (O-67): code runs the flow, harness sessions
do the work that needs judgement, a smart-if provider picks a branch at named
points, and humans decide only at the Human Decision Points of the PSB.

This document is written in slices (plan v2 of task `T-hbw8`); sections 5 to
13 come with slices B to H: 5 Intake and setup, 6 native gates and rule
protection, 7 the specification, 8 the phase loop, 9 squads and routing, 10
decisions, 11 stalls, 12 cost and the measures, 13 learning. Each section
names the ADR that decides it and the walkthrough that tests it. The walkthroughs
are in [`walkthroughs/`](walkthroughs/README.md). The Operator's decisions O-66
to O-92 are quoted in [`runs/T-hbw8/operator-decisions.md`](../runs/T-hbw8/operator-decisions.md)
and [`runs/T-hbw8/inputs-from-pr-69.md`](../runs/T-hbw8/inputs-from-pr-69.md).

## 1. LAYUP and a target

**LAYUP** is the `layup` binary, a work area on the host that runs it (the LAYUP
host), and LAYUP's own repository, which holds the stack catalog and the
defaults. A **target** is the repository of one product that LAYUP delivers.
LAYUP stays outside every target (O-76, [ADR-0013](adr/0013-orchestrate-a-target-from-outside-with-layup-run.md)).

Only three things go into a target:

1. **The setup output** that O-76 and `F-0003#42` name: the pinned baseline with
   its values and their evidence, the facts, and the target's own native stack
   gates with their CI job (section 6).
2. **The work of the role sessions**, through pull requests.
3. **LAYUP's records**, on the target's records branch (section 3).

No file goes in that the target needs LAYUP to build, test or pass its gates. A
target passes its own gates with LAYUP absent (Invariant 2).

**The LAYUP host** holds, per target: a clone of the target, one directory per
role session, the harness credentials, and the private key of the LAYUP GitHub
App. None of these is project state: each run rebuilds the clone and the
session directories from the forge and the records (Invariant 1).

**The forge** is GitHub for the pilot. LAYUP talks to it through one package with
a named set of calls: issues and comments, pull requests, commit statuses, the
effective rules of a branch, and the repository activity. A second forge is not
designed (known limit L-A2).

## 2. The components

| Component | What it is | Where it runs | What it may write |
| --------- | ---------- | ------------- | ----------------- |
| The engine checks | `layup psb check`, `layup setup verify`, `layup gate`, `layup spec check`, `layup audit`, `layup report`: pure commands; each reads files and prints a typed table | the LAYUP host | only standard output |
| The orchestrator | `layup run TARGET`: the phase loop (section 8), in the foreground, one run per target | the LAYUP host | the records branch, the forge, the session directories |
| The decision component | the smart-if: a package of `layup run` that asks a provider at named points (section 10) | inside `layup run` | a result that `layup run` records |
| The squad manager | the harness register, admission and routing (section 9) | inside `layup run` | the same |
| The learning loop | the reward and the routing update at each retrospective (section 13) | `layup run`, at the retrospective | the same |
| The runner | the LAYUP GitHub App and the LAYUP host; the dead-man job (section 11) | the forge; the host | through `layup run` only |
| The target | the product repository | the forge | — |

One run per target: `layup run` writes a lease row (run ID, host, start time) to
the records branch with a push that is not forced. A second run that finds a
live lease, or whose push is refused because the branch moved, stops.

## 3. Records and identities

[ADR-0014](adr/0014-keep-the-records-in-the-target-with-one-writer.md) decides
this section.

**Where the records live.** In the target, on the branch `layup-records`: tables
of tab-separated values with a header row, one file per record kind (an event
table only adds rows; a register is edited in place and Git history is its
log), and Markdown payloads beside them (ADR-0011 decision 2 keeps its form).
The target's README, written at setup, names the branch. A plain `git clone`
carries it as `origin/layup-records`, and a human reads every file with no tool
(Invariants 1, 2 and 9).

Records do not go on the default branch: a commit for each event would put every
open pull request out of date, and an agent's pull request could change a
record. ADR-0011 rejected a separate branch as "outside the tree a reviewer
reads"; the answer is that each pull request body links the records commit that
its task started from.

**One writer.** Only `layup run` commits to the records branch. Every other
producer hands it a typed result, which `layup run` checks and commits:

| Producer | What it hands over | The check before the commit |
| -------- | ------------------ | --------------------------- |
| An engine check | its table on standard output | the header and the columns of that record kind |
| A role session | a result file in its session directory, and its commits in its work tree | the schema of the result; the base commit and the records commit it started from are still the current ones for its task, or the result is refused |
| A human | an issue comment | the copy rule below |
| The dead-man job | nothing; it opens an issue (section 11) | — |

**Identities.** There are three kinds of actor, and the forge tells them apart
by credential, never by text that an agent writes or by a commit author:

- **The LAYUP App.** `layup run` alone holds the App's private key and acts with
  its installation token, as the App's bot. The App's permissions: contents,
  issues, pull requests and commit statuses (write), and workflows (write),
  because the setup writes the target's CI job. It has no administration
  permission. For the pilot, the App can be `layup-agent`, the App of O-92, with
  a private key added.
- **Humans.** The Operator and the idea owner, each named at Intake in
  `approvers.tsv` by the forge's numeric user ID (O-77: one GitHub account is
  enough; the idea owner and the Operator can be the same person).
- **Role sessions.** They hold no forge credential (section 4). Each comment that
  the App posts for a session starts with the session ID, and the records bind
  the session's commits to its harness and model.

**A human decision** is an issue comment on the target whose author ID is in
`approvers.tsv` and whose `performed_via_github_app` field is empty (O-77). A
review of a pull request, a review comment, a commit, a reaction or an edit is
never a decision: the review API has no App field (Sol-30), and a commit author
is text. Each of them is still recorded as human input (section 12).

**Copy before read.** Before any step acts on a comment, `layup run` copies it
to the records: the body, the author ID and login, the comment ID, the App
field, the time and a SHA-256 of the body. An edit of a comment is a new input
event; the first copy stays.

**The forge must enforce this, or LAYUP stops.** At the Scaffold phase and at
each start, `layup run` reads the effective rules of the default branch and of
the records branch (section 6 sets them) and stops with the difference when
they are not the expected ones (Invariant 5). `layup audit` reads the
repository activity of both branches and fails when an update of the records
branch came from another actor than the App, or an update of the default branch
was not a merge of a pull request by the App (FT3).

**Communication through issues** (O-73, vision 3.4). Every question to a human,
escalation, bet, acceptance, stall package and status is an issue comment that
the App posts. Between role sessions, the records are the channel, and
`layup run` also posts each handoff, question and answer on the task's issue, as
the target's issue workflow asks. A comment is never the only copy: the records
are the system of record (Invariant 1). This is where LAYUP differs from vision
3.4 ("exclusively through … issues"), and Invariant 1 is the reason.

## 4. Role sessions and model calls

[ADR-0015](adr/0015-keep-model-calls-out-of-the-engine-checks.md) decides this
section. It is the successor of ADR-0011 decision 8 (O-72).

**Where a model is called.**

1. The engine checks start no model process and open no connection to a model
   service.
2. The `layup` process opens a connection to a model service only in the
   smart-if provider client, at a named decision point, and each call writes one
   row to `decisions.tsv` (section 10).
3. A role session is a harness process that `layup run` starts. Its model calls
   belong to the harness, and the cost ledger records the session (section 12).

**How a session starts.** `layup run` makes a session directory with:

- `work/`: a Git work tree of LAYUP's clone of the target, at the base commit, on
  the branch `task/<task>/<attempt>`;
- `home/`: an empty home directory;
- `prompt.md`: the prompt file that code builds from the records and the target's
  own rule files (section 8);
- `result/`: where the session writes its typed result.

It starts the harness with the command line of its row in the harness register
(the permission flag, the model flag, the prompt file, the output format, the
session or resume flag) and with an environment that holds only a named list:
`PATH`, the locale, `HOME=home/`, `GIT_CONFIG_NOSYSTEM=1`, a Git configuration in
`home/` with the author name `layup session <ID>`, and the one variable that
carries the harness's own credential. It holds no `GH_TOKEN`, no SSH agent socket
and no forge credential. The session reads its rules only from the work tree: the
target's `AGENTS.md` and the documents it names; the harness's own entry file in
the target is a pointer to `AGENTS.md` (Invariant 9, Fable-M23).

**How a session ends.** When the harness process exits, `layup run` reads the
result file (status `completed`, `blocked`, `needs_context`, `decision_needed` or
`failed`, with a reason), the commits in `work/`, and the harness's usage output.
A session that exits with no valid result file has failed. `layup run` pushes
the session's branch with the App's token; the session never pushes.

**The limit.** Sessions run under the Operator's operating-system user. LAYUP
removes the credentials it knows, but a process can still read any file that the
user can read, the App's private key and the key chain included. The separation
holds against a session that uses its prompt and its environment; it does not
hold against one that searches the host (O-77: "by convention only"). Known
limit L-A1.

## 14. Coverage

Each row points to a walkthrough; a row with no walkthrough fails
(`runs/T-hbw8/root-cause-missed-solution.md`, fix 5). Slices B to H add their
rows. The evidence for each row is [`runs/T-hbw8/rewrite-checklist.md`](../runs/T-hbw8/rewrite-checklist.md).

| Item | Walkthrough | Sections | ADRs |
| ---- | ----------- | -------- | ---- |
| S12 Harness-Agent Neutrality (`F-0003#52`) | [W-12](walkthroughs/W-12-harness-agent-neutrality.md) | 1, 3, 4 | 0013, 0014, 0015 |
| `NFR-001`, Invariant 1 | W-12 steps 6, 10, 14 | 3 | 0014 |
| `NFR-002`, Invariant 2 | W-12 steps 12 to 16 | 1, 3 | 0013 |
| `NFR-005`, Invariant 6 (which process calls a model) | W-12 steps 2, 9 | 4 | 0015 |
| `NFR-007` (Go, standard library, `git`) | — (a build check: `go list -deps`) | 1 | 0013 |
| #69 B1 (nothing starts the role agents) | W-12 steps 2, 8 | 2, 4 | 0013, 0015 |
| #69 B6 (decisions on the forge) | W-12 step 6 | 3 | 0014 |
| #69 forge question | — | 1 | 0013 (known limit L-A2) |
| Table C, 3.4 (communication through issues) | W-12 step 6 | 3 | 0014 |

## 15. Known limits

Each limit is a finding that the design does not close, recorded here (O-66).

- **L-A1. The host is shared.** Role sessions run under the Operator's user, so
  a session that searches the host can reach the App's private key and the
  Operator's own credentials. Until sessions run in an isolated environment (a
  container or another operating-system user), the separation of O-77 is by
  convention for that case (K03, K04).
- **L-A2. One forge.** GitHub is the only forge for the pilot (#69 forge
  question).
- **L-A3. One host during delivery.** `layup run` runs in the foreground on one
  host; while the host is down, nothing moves (section 11 says how the stall is
  found).
