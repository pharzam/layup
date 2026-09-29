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

**The LAYUP host** holds, per target: LAYUP's own clone of the target, one
separate clone per role session, the harness credentials, and the App's user
access token. None of these is project state: each run rebuilds the clones from
the forge and the records (Invariant 1).

**The forge** is GitHub for the pilot. LAYUP talks to it through one package with
a named set of calls: issues and comments, pull requests, commit statuses, the
effective rules of a branch, and the repository activity. A second forge is not
designed (known limit L-A2).

## 2. The components

| Component | What it is | Where it runs | What it may write |
| --------- | ---------- | ------------- | ----------------- |
| The engine checks | `layup psb check`, `layup setup verify`, `layup gate`, `layup spec check`, `layup report` read files; `layup audit` reads files and the forge's read-only API; each prints a typed table | the LAYUP host | only standard output |
| The orchestrator | `layup run TARGET`: the phase loop (section 8), in the foreground, one run per target | the LAYUP host | the records branch, the forge, the session clones |
| The decision component | the smart-if: a package of `layup run` that asks a provider at named points (section 10) | inside `layup run` | a result that `layup run` records |
| The squad manager | the harness register, admission and routing (section 9) | inside `layup run` | the same |
| The learning loop | the reward and the routing update at each retrospective (section 13) | `layup run`, at the retrospective | the same |
| The runner | the LAYUP GitHub App and the LAYUP host; the dead-man job (section 11) | the forge; the host | through `layup run` only |
| The target | the product repository | the forge | — |

**One run per target.** `layup run` holds a lease row on the records branch: run
ID, host, start time and a heartbeat counter that it increases every `lease.H`
(a parameter, section 10). A second run watches the branch and times the lease
by its own clock from the moment it last saw the counter change, so the two
hosts' clocks are never compared (the pattern of Kubernetes leader election). It
stops while the counter moves; when the counter has not moved for `3 × lease.H`
of its own time, it takes the lease over and records the takeover. The time with
no run is a stall of the orchestrator (later: slice F).

**Fencing.** Every write of `layup run` to the forge (a push of a task branch, a
comment, a status, a merge) comes after the records push that announces it. Any
refused push of the records branch stops the run: it re-reads the lease and goes
on only while it still holds it. So a run whose lease was taken over, for example
after the host slept, writes nothing more.

## 3. Records and identities

[ADR-0014](adr/0014-keep-the-records-in-the-target-with-one-writer.md) decides
this section.

**Where the records live.** In the target, on the branch `layup-records`, an
orphan branch that shares no history with the default branch (so it carries no
CI file, and a records push runs no CI): tables
of tab-separated values with a header row, one file per record kind (an event
table only adds rows; a register is edited in place and Git history is its
log), and Markdown payloads beside them (ADR-0011 decision 2 keeps its form).
The target's README, written at setup, names the branch. A plain `git clone`
carries it as `origin/layup-records`, and a human reads every file with no tool
(Invariants 1, 2 and 9). When LAYUP stops for good, the records stop growing; to add to
them, the Operator first lifts the records ruleset, as an admin can.

Records do not go on the default branch: a commit for each event would put every
open pull request out of date, and an agent's pull request could change a
record. ADR-0011 rejected a separate branch as "outside the tree a reviewer
reads"; the answer is that each pull request body links the records commit that
its task started from.

**One writer.** Only `layup run` commits to the records branch, from its own
clone, which no session uses. Every other producer hands it a typed result,
which `layup run` checks and commits:

| Producer | What it hands over | The check before the commit |
| -------- | ------------------ | --------------------------- |
| An engine check | its table on standard output | the header and the columns of that record kind |
| A role session | a result file in its session directory, and commits in its own clone | the schema of the result; the attempt and the base commit come from the session start row that `layup run` wrote, never from the result; the result is refused unless that attempt is still the task's open attempt: no later row of the task's event table closes it, starts another attempt, or changes its base commit |
| A human | an issue comment | the copy rule below |
| The dead-man job | nothing; it opens an issue (section 11) | — |

**Identities.** O-77 and O-95 decide them: no second GitHub account; LAYUP acts
as the LAYUP GitHub App's bot, with an installation token, so the forge shows
`<app>[bot]` with the App's logo, and the API field `performed_via_github_app`
names the App. The Operator owns and installs the App and stays the responsible
actor.

- **LAYUP.** `layup run` alone holds the App's private key and makes an
  installation token from it for each hour of work. The App's permissions are
  those of O-92 (contents, issues and pull requests: write; metadata: read) plus
  commit statuses (write); no workflows and no administration. For the pilot
  this is the App `layup-agent`, installed on each target and given the
  commit-statuses permission; both changes are the Operator's.
- **Humans.** The Operator and the idea owner, each named at Intake in
  `approvers.tsv` by the forge's numeric user ID. They can be the same person.
- **Role sessions.** They hold no forge credential (section 4). Each comment that
  LAYUP posts for a session starts with the session ID, and the records bind the
  session's commits to its harness and model.

**A human decision** is an issue comment on the target whose author ID is in
`approvers.tsv` and whose `performed_via_github_app` field is empty (O-77). A
review of a pull request, a review comment, a commit, a reaction or an edit is
never a decision: the review API has no App field (Sol-30), and a commit author
is text. Each of them is still recorded as human input (section 12).

**Copy before read.** Before any step acts on a comment, `layup run` copies it
to the records: the body, the author ID and login, the comment ID, the App
field, the time and a SHA-256 of the body. An edit of a comment is a new input
event; the first copy stays.

**What the forge enforces, and what LAYUP checks.** Section 6 sets two rulesets:
the default branch has no bypass actor, and the records branch restricts updates
to the LAYUP App, its only bypass actor. At the Scaffold phase and
at each start, `layup run`:

1. reads the effective rules of both branches (`GET
   /repos/{owner}/{repo}/rules/branches/{branch}`) and stops when a rule is
   missing (Invariant 5);
2. pushes an empty probe commit to the ref `layup-probe`, which the same ruleset
   as the default branch covers, with the App's token, and stops unless the
   forge refuses it. A bypass is set per ruleset, so this proves that no bypass
   covers LAYUP's own actor on the default branch, and a probe that the forge
   accepts lands on a ref that nothing uses.

The list of bypass actors is not readable with the App's permissions (the API
returns it only to a token with write access to the ruleset). The Operator's
apply command at setup prints it, and `layup run` records that output; a later
change to it is not seen (known limit L-A4).

**The actor on the forge.** `layup audit` reads the repository activity of both
branches, from the setup commits that the records name onward. Each update of
the records branch must be a push by the App's bot. Each update of the default
branch after the setup must be a pull-request merge (`pr_merge` or
`merge_queue_merge`) by the App's bot or, for a batch that changes
`.github/workflows/` (O-93), by an account in `approvers.tsv`; a `push` there
means a bypass, and an update with no actor fails (FT1). The setup commits are
allowed by their SHA. The audit never reads a commit author (FT3).

**Communication through issues** (O-73, vision 3.4). Every question to a human,
escalation, bet, acceptance, stall package and status is an issue comment that
LAYUP posts. Between role sessions, the records are the channel, and `layup run`
also posts each handoff, question and answer on the task's issue, as the
target's issue workflow asks. A comment is never the only copy: the records are
the system of record (Invariant 1). This is where LAYUP differs from vision 3.4
("exclusively through … issues"), and Invariant 1 is the reason.

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

**How a session starts.** `layup run` makes a session directory under its work
area with:

- `repo/`: a separate clone of the target, made with `--no-local` (not a linked
  work tree, which would share refs, configuration and hooks with LAYUP's clone,
  and not hard links to its objects), at the base commit, on the branch
  `task/<task>/<attempt>`; it has no remote with a credential and no records
  branch;
- `home/`: an empty home directory;
- `prompt.md`: the prompt file that code builds from the records and the target's
  own rule files (section 8);
- `result/`: where the session writes its typed result.

It starts the harness with the command line of its row in the harness register
(the permission flag, the model flag, the prompt file, the output format) and
with an environment that holds only a named list: `PATH`, the locale,
`HOME=home/`, `GIT_CONFIG_NOSYSTEM=1`, a Git configuration in `home/` with the
author name `layup session <ID>`, and the harness's own credential: a variable,
or a file that the register row names and that `layup run` copies into `home/`
(for example a harness whose login is a file in its configuration directory,
which then also lies under `home/`). It holds no `GH_TOKEN`, no SSH agent socket and no
forge credential.

**Rules only from the target.** The session reads its rules from its clone: the
target's `AGENTS.md` and the documents it names; the harness's entry file in the
target is a pointer to `AGENTS.md` (Invariant 9, Fable-M23). A harness can also
load rule files from the directories above its working directory and from a
system-wide policy path. So, before each start, code checks every directory from
the session directory up to the root for the rule-file names that the harness's
register row lists, and each policy path that the row lists; it refuses to start
when it finds one above the session, and records each policy file it finds. The
harness probe (section 9) records which instruction files the harness reports
that it loaded, where the harness can report them. A system-wide policy file that
a harness always loads is known limit L-A5.

**How a session ends.** When the harness process exits, `layup run` reads the
result file (status `completed`, `blocked`, `needs_context`, `decision_needed` or
`failed`, with a reason) and the harness's usage output. A session that exits
with no valid result file has failed. `layup run` then fetches the session's
branch into its own clone by commit SHA, with hooks turned off
(`core.hooksPath` set to an empty directory) and no configuration read from the
session's clone; it checks that the branch descends from the base commit; and it
pushes that SHA to `task/<task>/<attempt>` with the App's token. Before the push,
it refuses a branch whose diff from the base touches `.github/workflows/`: the
App has no workflows permission (O-92), so LAYUP never delivers a workflow
change; the result fails with that reason, and the change becomes a proposal for
an approved batch, which its human approver pushes and merges (O-93). The records
bind a commit SHA to the session only after the forge accepted the push. The
session never pushes. A session is not resumed across attempts: harness state in `home/` is
thrown away, and the next attempt starts from the records (FT6).

**The limit.** Sessions run under the Operator's operating-system user. LAYUP
removes the credentials it knows, but a process can still read any file that the
user can read: the App's private key, and the Operator's own `gh` login, whose
comments have an empty App field and would pass as human decisions. The
separation holds against a session that uses its prompt, its environment and its
own clone; it does not hold against one that searches the host (O-77: "by
convention only"). Known limit L-A1.

## 14. Coverage

Each row points to a walkthrough, or names the check or the known limit that
stands in its place; a row with none of the three fails
(`runs/T-hbw8/root-cause-missed-solution.md`, fix 5). Slices B to H add their
rows. The evidence for each row is [`runs/T-hbw8/rewrite-checklist.md`](../runs/T-hbw8/rewrite-checklist.md).

| Item | Walkthrough, check or limit | Sections | ADRs |
| ---- | --------------------------- | -------- | ---- |
| S12 Harness-Agent Neutrality (`F-0003#52`) | [W-12](walkthroughs/W-12-harness-agent-neutrality.md) | 1, 3, 4 | 0013, 0014, 0015 |
| `NFR-001`, Invariant 1 | W-12 steps 6, 10, 14 | 3 | 0014 |
| `NFR-002`, Invariant 2 | W-12 steps 12 to 16 | 1, 3 | 0013 |
| `NFR-005`, Invariant 6 (which process calls a model) | W-12 steps 2, 9 | 4 | 0015 |
| `NFR-007` (Go, standard library, `git`) | the check of its `PRD-0001` criterion: `go list -deps ./...` names no package outside the standard library and the `layup` module | 1 | 0013 |
| #69 B1 (nothing starts the role agents) | W-12 steps 2, 8 | 2, 4 | 0013, 0015 |
| #69 B6 (decisions on the forge) | W-12 step 6 | 3 | 0014 |
| #69 forge question | known limit L-A2 | 1 | 0013 |
| Table C, 3.4 (communication through issues) | W-12 step 6 | 3 | 0014 |

## 15. Known limits

Each limit is a finding that the design does not close, recorded here (O-66).

- **L-A1. The host is shared.** Role sessions run under the Operator's user, so
  a session that searches the host can reach the App's private key, which does
  not expire until the Operator revokes it, and the Operator's own `gh` login; a comment made with that login passes as a human decision.
  Until sessions run in an isolated environment (a container or another
  operating-system user), the separation of O-77 is by convention for that case
  (K03, K04).
- **L-A2. One forge.** GitHub is the only forge for the pilot (#69 forge
  question).
- **L-A3. One host during delivery.** `layup run` runs in the foreground on one
  host; while the host is down, nothing moves (section 11 says how the stall is
  found). A takeover on another host needs the App's private key on that host.
- **L-A4. The bypass list is read once.** The App cannot read a ruleset's
  bypass list, so it is read only at setup, from the Operator's command; a later
  change to it is not seen. The forge's rule-suites API, which reports a bypass after the fact, may narrow
  this; the permission it needs is not yet checked.
- **L-A5. A policy file of the host.** A harness that always loads a system-wide
  policy file gives its sessions rules that another harness does not get; code
  records the file, and does not remove it.
