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
separate clone per role session, the harness credentials, and the App's private
key, in a file only the Operator's user can read (mode 0600), from which
`layup run` makes installation tokens. None of these is project state: each run rebuilds the clones from
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
comment, a status, a merge) comes after the records push that announces it. A
run pushes its records only on top of its own last pushed records commit, never
after a fetch and a rebase, so a takeover makes its next records push fail. Any
refused push of the records branch stops the run: it re-reads the lease and goes
on only while it still holds it. So a run whose lease was taken over, for example
after the host slept, writes nothing more, except the one forge write that it
had already announced (L-A3).

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
| A role session | a result file in its session directory, and commits in its own clone | the schema of the result; the attempt and the base commit come from the session start row that `layup run` wrote, never from the result; the result is refused unless that attempt is still the task's open attempt: no later row of the task's event table (`tasks/<task>/events.tsv`) closes it, starts another attempt, or changes its base commit |
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
   as the default branch covers, with the App's token. Only a refusal for a rule
   violation (GitHub's `GH013` "Repository rule violations found", naming the
   ruleset) counts; any other failure (an expired token, a network or server
   error) stops the run as "probe not run", never as a pass (FT1). A bypass is set per ruleset, so this proves that no bypass
   covers LAYUP's own actor on the default branch, and a probe that the forge
   accepts lands on a ref that nothing uses (the Operator then deletes the ref and
   checks any workflow run that it started);
3. reads the repository activity with the App's token, and stops when it cannot
   (the audit below needs it).

The list of bypass actors is not readable with the App's permissions (the API
returns it only to a token with write access to the ruleset). The Operator's
apply command at setup prints it, and `layup run` records that output; a later
change to it is not seen (known limit L-A4).

**The actor on the forge.** `layup audit` reads the repository activity of both
branches, from the setup commits that the records name onward. Each update of
the records branch must be a `push` by the App's bot, apart from one
`branch_creation` by the App's bot whose SHA is the records' first commit. Each update of the default
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
an approved batch, which its human approver pushes and merges (O-93). The diff of the
refused change goes to the records branch as a payload, where the approver finds
it. The records bind a commit SHA to the session only after the forge accepted
the push. The
session never pushes. A session is not resumed across attempts: harness state in `home/` is
thrown away, and the next attempt starts from the records (FT6).

**The limit.** Sessions run under the Operator's operating-system user. LAYUP
removes the credentials it knows, but a process can still read any file that the
user can read: the App's private key, and the Operator's own `gh` login, whose
comments have an empty App field and would pass as human decisions. The
separation holds against a session that uses its prompt, its environment and its
own clone; it does not hold against one that searches the host (O-77: "by
convention only"). Known limit L-A1.

## 5. Intake and setup

[ADR-0016](adr/0016-put-the-native-stack-gates-in-the-target.md) and
[ADR-0017](adr/0017-prevent-rule-changes-by-agents.md) decide the gates and the
rules; this section gives the order. The public tools make the repository first
and record the intake in it; they apply the forge settings later
([`selection-v2.md`](../runs/T-hbw8/selection-v2.md) §1.5). LAYUP does the same.

**Start.**

1. The Operator makes an empty repository on the forge, installs the LAYUP App
   on it, and runs `layup run --new OWNER/NAME --psb FILE [--vision FILE]
   --operator LOGIN --idea-owner LOGIN`.
2. `layup run` copies the pinned baseline with `git clone` of the pinned commit
   and removes `.git` (ADR-0011 decision 7; step S02 changes from `npx degit`).
   The Operator pushes this unmodified copy as the root commit of the default
   branch, with the Operator's own login and one command that `layup run`
   prints: the copy holds CI files, and the App has no workflows permission
   (O-92). So the default branch is the first branch of the repository, and its
   root tree equals the pinned tree, as step S03 asks. The command pushes from a
   plain clone with no hooks installed.
3. `layup run` reads back the repository's default branch, its root tree and its
   visibility and plan, and stops when one differs or when the plan does not
   enforce rulesets on it (on GitHub Free, a private repository: the Operator
   makes it public or moves it to a plan that does, K15). It then pushes the records branch (an orphan, §3)
   with its first commit: the problem statement and the vision brief byte for
   byte, each with its SHA-256; `approvers.tsv` with the numeric IDs of the two
   logins of step 1; and the lease row. From this commit on, only a comment by
   one of those IDs is an answer or a decision (§3); any other comment is
   recorded as input. It opens the Intake issue.

**The gap check, in one batch** (Decision Points 1 and 2, `F-0001#10`, `#11`):

1. `layup psb check` applies rules G1 to G5 (`internal/psb`) to the problem
   statement. They find text patterns: no named stack, a metric row with no
   method, an undefined abbreviation, a vague word, an uncalibrated target. They
   do not find a gap of meaning (`REQ-001` changes with this, plan step 10).
2. A review session (the roles and harness of the step table, §9) reads the
   problem statement and writes the gaps of meaning as typed rows: a byte-exact
   quote, the kind (a term with two readings, a metric that cannot be measured as
   written, a missing intent value, and the four questions of Shape Up's "rabbit
   holes": new technical work, an assumption about how parts fit, a design that
   may not exist, a hard decision to settle before the work), and the question.
   Code checks each quote against the file.
3. The specification sessions of §7 add the needs that the requirement draft
   leaves open (§7 steps 1 to 4).
4. The Operator's questions of setup steps S01 and S10 of LAYUP's
   `setup/steps.tsv` (the stack, and each marker of the pinned baseline; the
   name and the visibility are already fixed) come from LAYUP's steps and the
   pinned commit. For a marker, a session proposes a source from the facts with a
   byte-exact quote, and code checks the quote; a marker with no proposed source
   becomes a question.
5. Code merges the rows, removes a question whose quote and kind repeat, gives
   each an ID, and posts **one** comment on the Intake issue. It has two blocks:
   the questions and the intent form for the idea owner (the success criteria,
   the appetite and the band of §12, the approver of each planned approval
   point; a login named there is added to `approvers.tsv`, because the idea
   owner named it), and the setup questions, the proposed marker sources and the
   parameters for the Operator (the owner map of §9, the smart-if provider and
   authority of §10).
6. The idea owner and the Operator answer, each in one comment, one line per
   question ID. `layup run` copies each comment (§3) and checks that every ID has
   an answer and that each number and login parses. A review session then reads
   the answers and lists each one that leaves its gap open. `layup run` asks
   again, in one comment, only for the missing and the open answers. Each
   question and each follow-up is a row of the project question table with the
   flag "before delivery" (Early Question Share, §12).

**Scaffold** (`layup setup`, LAYUP's `setup/steps.tsv` with the changes named
here). `setup/steps.tsv` and `setup/setup-check.sh` are LAYUP's own files, not
the pinned baseline's: neither goes into a target (FT5). `layup setup verify`
does each check of `setup-check.sh` from outside, for this target (its name, its
facts, its pin), and the evidence of steps S04 to S15 becomes "`layup setup
verify <check>` OK".

1. On a branch from the root commit, code does each row that needs no judgement.
   A value comes only from an Intake answer (by question ID), a catalog entry
   (§6), or a fact source that the Operator accepted in the Intake answer. The
   setup record names the source of each value; a value with no source keeps its
   marker and becomes an open gap of the target (Invariant 4).
2. A role session writes the prose rows (S07, S08, S09, S14), and code runs
   each row's check.
3. The facts go into the target's `docs/facts/`: the problem statement byte for
   byte, and the answers as a raw fact; the numbered facts come with the first
   bet (§7).
4. The stack gates come from the catalog (§6). Step S12 adds no
   `setup-check` job, and the target's CI runs only the pinned baseline's own
   jobs and the gate jobs (K16).
5. `layup setup verify` checks the result from outside: the pinned baseline's own
   discipline tests on a checkout of it; its own checks (pin: the root tree
   equals the pinned tree; markers; facts; and the others of `setup-check.sh`);
   zero values without a source; each source resolving; and each active gate,
   which must pass on the clean tree and fail on its known-bad fixture (§6).
6. The Operator pushes the setup commits on top of the root commit and applies
   the two rulesets of §6 from a file that `layup run` writes, with the
   Operator's own login and the commands it prints: the setup changes CI files,
   and the App has neither the workflows nor the administration permission
   (O-92). Step S13 still writes `docs/setup/branch-protection.json`; the
   Operator applies the rulesets, which hold the same required checks and
   LAYUP's, in place of the classic protection. The rulesets name each required
   check with its source App by ID (`integration_id`), so they need no earlier
   run of the check. This is the planned approval
   point "setup", listed at Intake.
7. `layup run` checks that the pushed tree equals the verified one, then reads
   back the effective rules of both branches and makes the probes of §3: the push
   to `layup-probe` must be refused; a read of the repository activity with the
   App's token must work. When a rule is missing or a probe gives the wrong
   result, it stops. On GitHub Free, a private repository has no enforced
   rulesets, so it stops there too; the Operator makes the repository public or
   moves it to a plan that enforces them (K15).

**The pilot baseline** (`F-0004#11`). For a pilot target, the idea owner posts
the baseline numbers measured with the current process, each with its evidence,
as a block of the Intake answer; `layup run` records them. `layup report`
compares a measure with its start value only when both exist; otherwise it
prints "not comparable", never a pass. The idea owner's batch of start values,
and where it is recorded, are later: slice G (§12).

## 6. Native gates and rule protection

### The stack gates

[ADR-0016](adr/0016-put-the-native-stack-gates-in-the-target.md) decides this.

**The stack catalog** is in LAYUP's repository, one directory per stack. An entry
lists, per gate kind (layout, interface boundary, contract, test quality): the
tool and its version, the command, the paths in scope, the configuration it
writes, a known-bad fixture that must make the gate fail, and the evidence for
the tool (its documentation at the version). LAYUP's own CI runs every fixture.
A new stack gets an entry through a LAYUP change under LAYUP's gate, by a person
or a LAYUP task, never by a session that works on a target (#69 B9).

**In the target**, the setup writes the tools' configuration, the gate manifest
`docs/gates.tsv` (gate kind, command, scope, state), and **one CI job per gate
kind**, each its own required check, pinned in the ruleset to GitHub Actions as
its source (so a status of the same name from another source does not count).
These are the target's own files; they run with LAYUP absent (Invariant 2,
O-76). A gate kind whose rules depend on the architecture (in the Go entry:
layout, boundary and contract) has the state `pending` until its activation
(below). Its job then fails a pull request that changes a path in the product's
scope, and passes one that changes none, with that reason in its output: no
product code merges before its gates exist, and a pending gate never passes code
(Invariant 5). An active kind treats a tree with no path in its scope the same
way (for Go, `./...` that matches no package), so it never passes silently on
nothing. The jobs of all kinds exist from the setup, so the activation changes no
CI file.

**The Go entry**: the setup writes `go.mod` (the module path from the repository
name); `test -z "$(gofmt -l .)"` (bare `gofmt -l` exits 0 on a bad file) and
`go vet ./...` are active from the setup; at activation, an
import-rule tool (golangci-lint's `depguard`, configured from the package table
of the approved architecture) for the boundary, a layout test as a Go test, and a
contract test for each interface that the architecture names; test quality is
`go test -count=1 ./...`, with a coverage floor that stays an open gap until the
idea owner or the Operator sets it with evidence (L-B2).

**`layup gate`** runs the same commands from outside: it checks out the base
branch's manifest and gate files, applies them to the head of the pull request
in a scratch work tree, and gives one result per kind. `layup run` posts the
commit status `layup/gates` from those results:

| Result of a kind | `layup/gates` |
| ---------------- | ------------- |
| active, and it passed | counts as a pass |
| active, and it failed or did not run (`not-active`) | failure |
| `pending`, and the head changes no path in the product's scope | `clear`: counts as a pass, with the reason "pending: no product path" |
| active, and the tree has no path in its scope (for Go: no package) | `clear`: counts as a pass, with the reason "no product path" |
| `pending`, and the head changes such a path | failure |

The status is a success only when every kind counts as a pass. It never runs the
head's gate files (FT4), with one exception: the activation batch below, whose
rule files are the approved ones. The LAYUP host needs each stack's toolchain
(L-B1). With LAYUP absent, the target's own CI runs the workflow and the manifest
of the pull request's head; FT4 then rests on the rulesets and on review, not on
LAYUP.

**Activation.** At the first bet (§8), an architect session writes the
boundary configuration, the layout test, the contract tests, the manifest with
those kinds `active` (in the batch head itself), and, for each kind it
activates, one known-bad patch. The patches go to the records branch as
payloads; they never merge (FT6: a clone carries them). This is a rule batch
(below). When its approval is recorded, `layup gate` runs the batch's own gate
files on the batch head, which must pass, and on the head with each known-bad
patch applied, which must fail its kind (`F-0003#64`).

### Rule protection

[ADR-0017](adr/0017-prevent-rule-changes-by-agents.md) decides this.

**The rule paths** of a target are a register on the records branch, written at
setup from the catalog and the pinned baseline: `.github/`, `.githooks/`, the gate
manifest, the tools' configuration, the gate tests, `AGENTS.md` and the harness
entry files, the baseline's rule documents and check scripts, `docs/setup/`, and
`docs/facts/` (§7). `docs/guardrails.md` is a rule path except added lines inside its
section 2 (Known pitfalls): the baseline's step 7 writes a lesson there in the
same pull request, and a line added to §2 cannot remove a rule (K13).

**A rule batch** is the one way a rule changes. It is proposed at a planned
point: the setup, the gate activation of the first bet, and each retrospective
(O-69). `layup run` pushes it to a branch `batch/<point>` and opens its pull
request, so the approver sees the change before deciding. The approval request
names the batch head's SHA and the tree hash of its rule files. The approver
approves by an issue comment (§3); code refuses the approval when the batch head
has moved since the request, and records the hash with the approval. The approved batch merges by
`layup run`, or, when it changes `.github/workflows/`, by its approver, who also
pushes it, because the App has no workflows permission (O-93).

**Prevention, in three layers:**

1. **No credential.** A role session holds no forge credential (§4), so it
   cannot push, merge, approve or change a ruleset with a token.
2. **No rule change leaves the host in a task.** Before `layup run` pushes a
   task branch, it compares the changed paths with the rule-path register of the
   base branch. A change to a rule path refuses the result: the task gets the
   finding, and the change becomes a proposal for the next batch.
3. **The forge refuses.** The default branch and the probe ref `layup-probe`
   have one ruleset with an empty bypass list: a pull request is required; the
   required checks are the native gate jobs, pinned to GitHub Actions, and
   `layup/gates`, `layup/spec`, `layup/verify` and `layup/rules`, pinned to the
   LAYUP App; force
   pushes and deletion are blocked. The records branch has a ruleset that
   restricts updates and deletion to the LAYUP App, its only bypass actor (§3).
   `layup/rules` fails a pull request that changes a rule path unless it is a
   batch whose head's rule-file hash equals the hash recorded with its approval.
   (`layup/verify` comes from the counterpart verification of §8.)

**Reading of `F-0003#64`.** "0 agent writes to rule paths" counts the rule-path
changes on the default branch that did not land in an approved batch; the
approval makes the change the human's decision (the approval brief asks the
Operator to confirm this reading).

**Detection, as the complement.** `layup audit` lists every change to a rule path
on the default branch with its merge, its actor and its approval (§3). The
known-bad fixtures run at setup for the active kinds, and at activation for the
others; a `pending` kind's fixture is recorded as not run, never as a
detection.

**When LAYUP is absent,** the `layup/` checks never report. The target's setup
record tells the Operator to remove them from the ruleset; the native gate jobs
stay (Invariant 2).

**The limit.** The Operator's own account can change the rulesets, and a process
on the host can reach the App's key (L-A1). Prevention holds against every
credential that a role session holds.

## 7. The specification

[ADR-0018](adr/0018-derive-the-specification-from-numbered-source-lines.md)
decides this section. It serves `F-0003#51` and Problem 6. A problem statement
is prose, so every step that reads its meaning is a session or a human; code
checks what can be checked byte for byte. The vision brief is solution input
(`F-0002`), not a trace source.

**During Intake, before the batch** (§5, gap check step 3):

1. **Numbering** (a session). It splits the problem statement into spans, each
   given by its byte offsets: a **fact** has an ID and a class (`need`,
   `constraint`, `measure`, `context`); a **"not a fact"** span (a heading, a
   rule line) has a reason. A span may be part of a line, so one sentence of a
   paragraph can be a fact. Offsets count bytes from 0, and the end is
   exclusive.
2. **The span check** (code, `layup spec check --facts`). Every byte that is not
   white space (ASCII space, tab, CR, LF) lies in exactly one span, each span's
   text is the file's text at its offsets, each fact ID is unique, and each class
   is one of the four. When the check fails or does not run, `layup run` posts no
   batch and stops (FT1, Fable-M21, Author-5).
3. **The draft** (a session). One or more requirements for each `need` fact, and
   one non-functional requirement for each `constraint` fact: an ID, a statement,
   `covers` (the fact IDs), and an acceptance criterion. It sets no priority.
   Each `measure` fact becomes a success criterion row of the target's PRD, which
   §12 reads; a `context` fact gets nothing.
4. **The completeness review** (a session on a harness other than that of steps
   1 and 3, which routing puts on one harness; with only one admitted harness it does not run, and `layup run`
   records "not run" and stops). It reads the problem statement and the draft,
   and lists each need, constraint or measure that no requirement or success
   criterion covers, each requirement
   that no fact supports, and each fact whose class it doubts (Sol-18). Each item
   becomes a question of the one Intake batch.

**After the answers:**

5. Code writes the answers as a raw fact of the target, one fact per question ID.
   Each question carries the byte-exact quote of the problem statement that it
   asks about (§5), so a requirement that covers an answer's fact also names that
   quote's fact, and code checks that it does. An answer that states a new need
   is classed as a fact of the answers, and the rules below apply to it.
   **Reading of `F-0003#62`:** the trace reaches the problem
   statement's text through the question's quote; a requirement that traces only
   to a question with no quote (a setup question, or a review item "a
   requirement that no fact supports") traces to the answer alone, and
   the approval brief asks the Operator to confirm this reading.
6. A session revises the draft with the answers. A change of class must cite the
   answer ID that causes it; code checks that the ID exists (not that it causes
   the change). The session writes the
   target's PRD and, per requirement, a technical specification section: a
   heading that holds the requirement ID, in a file under `docs/spec/`. The PRD's
   MoSCoW and Phase columns stay empty; the pull request waits.
7. **The first bet** (§8) shows the idea owner every fact with its class and
   text, and the requirements and success criteria. The idea owner answers one
   line per requirement ID: its priority (`Must`, `Should`, `Could`, `Won't`) and
   its milestone; one line per fact to change: "out of scope", or a new class
   (for example `need`). A new class reopens the draft for that fact (step 6),
   and its requirement waits for the next bet line. Priority is
   intent (`F-0003#54`), so code copies the comment (§3), writes the MoSCoW and
   Phase columns from that copy, and lists the out-of-scope fact IDs under the
   PRD's non-goals, and only
   then does the pull request pass the baseline's own `prd-lint`. A requirement
   added later gets its priority at the next bet. Inside a requirement, a session
   may mark a task nice-to-have; a cut that touches an acceptance criterion is
   business-forking (§10).
8. **The confirmed inventory.** At each bet, code writes a new version of it
   from that bet's copy: the numbered facts, their classes and the out-of-scope
   marks, on the records branch, with the SHA-256 of the version. The numbered
   facts record of the target, with the confirmed classes, lands in
   `docs/facts/` through that bet's rule batch (§6); the Scaffold writes only the
   problem statement and the answers there. `docs/facts/` is a rule path from the
   setup, so a role session cannot edit it.

**`layup spec check`** (code, no model call) posts its own required status,
`layup/spec`. It reads the latest confirmed inventory, whose hash it checks, and
the bet copy (records), the
PRD and `docs/spec/` (the head), the task register (§8) and the accepted and
delivered requirements (§12). It fails when: a `covers` ID does not resolve to a
fact of the confirmed inventory or of the answers; a need or constraint has no
requirement, or a measure has no success-criterion row that cites its ID, and it
has no out-of-scope mark; a requirement has no acceptance criterion;
a MoSCoW or Phase value in the PRD differs from the bet copy; a delivered
requirement has no non-empty section in `docs/spec/` whose heading holds its ID
(Sol-20); a task names no requirement; or it did not run. It proves the links. It
does not prove that the meaning agrees: the counterpart verification of each
change judges that (§8), and the completeness review of step 4 is the check
against a missed need.

**Success criteria.** The measure rows and the success criteria of the Intake
form are two sources; how §12 merges them is later: slice G.

**Vision 2.1.** The preliminary design review is the specification of step 6 with
the architecture that the first bet approves (§8); the phased plan is the
milestone plan of the bets (later: slice D, §8).

## 8. The phase loop

[ADR-0019](adr/0019-run-the-lifecycle-as-a-phase-loop-with-one-bet-per-milestone.md)
decides this section.

**The phases.** Intake and Scaffold (§5) → **Shape** → **Bet** → **Build** →
**Accept** → **Retrospective** (§13) → the next **Bet**, until every `Must`
requirement is accepted. A later milestone goes back to Shape only when its bet
brief proposes a change of the architecture.

- **Shape.** Sessions write the specification (§7 step 6) and the architecture of
  the target: its document and its decision records, by the target's own rules;
  a decision that those rules give to a panel gets the blind panel of §11. The
  architect session writes the gate-activation batch (§6).
- **Bet** (Decision Points 1 and 3; Shape Up's betting table, one bet per
  milestone). `layup run` posts one brief on the milestone's issue, one screen
  long, in five parts: the problem (the requirements of the milestone), the
  appetite (its cap of money and wall-clock, inside the band, §12), the solution
  (the architecture, or its change), the rabbit holes, and the no-gos; and, for
  the first bet, every fact with its class (§7 step 7). The idea owner answers by
  one comment: a line per requirement, a line per fact to change, and "bet" or
  "no bet". Each bet is a planned approval point, listed at Intake. It is also
  the architecture's approval: there is no other approval inside the milestone,
  apart from escalations (§10) and stalls (§11). Merge order after a bet: the
  specification pull request (its MoSCoW and Phase written from the bet), then
  the rule batch, then the tasks. A requirement changes only at a bet.
- **Build.** A plan session splits the milestone into tasks, each with the
  requirement IDs it serves and the tests that will show it done; code checks
  that each `Must` requirement of the milestone has a task and each task a
  requirement, and writes the task register. Then the task loop below runs, one
  task at a time or in parallel as the plan's order allows.
- **Accept** (Decision Point 3). For each delivered requirement, `layup run`
  posts the requirement, its criterion, the merged changes and the gate and
  verification results; the idea owner answers "accept" or "reject" with a
  reason. A rejection becomes a new need for the next bet (Shape Up: a new bet,
  not a patch). The record is in §12.

**The task loop.** Each step of the target's own quality gate has an actor
(Fable-M9):

1. `layup run` opens the task's issue (its requirement IDs, its Definition of
   Done).
2. A plan session writes the task's plan; a plan-review session on another
   harness reviews it. `layup run` posts both as issue comments in the forms that
   the target's own `review-record-lint` parses, and reads the review's verdict
   from its fixed field.
3. A developer session writes the failing test first, then the code, and the
   task file and close-out that the target's rules ask for; it ends with its
   typed result (§4).
4. The handoff check (below); the rule-path and workflow checks (§4, §6); then
   `layup run` pushes the branch and opens a **draft** pull request that links
   the issue, and requests no review (Sol-2).
5. The target's CI, `layup gate` and `layup spec check` report their statuses.
   A failure goes back to the developer as a finding, in a new attempt.
6. When they pass, a verifier session on a harness other than the author's
   (§9) reviews the head: fresh, one turn, read-only, told to refute "done",
   with a `file:line` checklist. `layup run` posts its record as a pull-request
   comment in the target's review-record form, and sets `layup/verify` at that
   head SHA: success only for the verdict `nothing material`. A material finding
   is one more round, up to the target's cycle cap; past the cap the round limit
   of §11 opens a stall.
7. When every required check is green at the head, `layup run` marks the pull
   request ready and merges it at that head SHA. No human review is requested: a
   routine review of each pull request is not a planned point (PSB §6). A human
   who reviews anyway is recorded as human input (§12).

**Handoffs** (`F-0003#45`, `#59`). A handoff is the typed result of one session
to the next role: session ID, task, role, status (`completed`, `blocked`,
`needs_context`, `decision_needed`, `failed`), reason, each artifact with its
path and SHA-256, the questions, the proposed decisions, the lessons (the fields
of Spec Kitty's handoff packet v1, as a pattern). The transition table, a
register on the records branch that the Operator can replace (O-81), says for
each pair of steps which artifact kinds are required and which state code must
compute (for example: the plan's tests fail at the base and pass at the head, as
`layup gate` runs them). A handoff is valid only when code finds the schema
right, each named artifact in the commits with its hash, and the computed state
true; a field that an agent writes never makes it valid (Sol-21, Fable-M14). The
Inter-Role Communication Format is the valid handoffs over all handoffs, from the
records.

**A question during the work** (Problem 1). A session that meets a question ends
with status `needs_context` and the question: its text, its own label for the
kind, and the records it concerns. `layup run` records it with its time, posts
it on the task's issue, and gets its kind and whether it needs a human from the
smart-if (§10). A question that needs no human goes to a session of the owner role
of its kind (§9), which answers with the records it cites. The answer passes the
escalation screen (§10), is recorded and posted, and the asking role starts its
next attempt with the answer in its prompt file. The answer is **accepted** when
that attempt ends with a result that cites the answer ID and does not ask the
same question again; the accepting actor is that session. An attempt that a
question ends does not count toward the attempt limits of §11.

## 9. Squads and routing

[ADR-0020](adr/0020-route-role-sessions-over-registered-harnesses.md) decides
this section.

**Roles and steps** (O-81). The roles are the seven functions of PSB §2 by
default, and the Operator can replace the matrix per target. The default step
table:

| Step | Role | Tier |
| ---- | ---- | ---- |
| Intake gap review; completeness review (§7 step 4) | Systems Architect; QA Engineer | reasoning |
| Numbering, draft, specification (§7) | Product Owner | reasoning |
| Shape: the architecture; the gate-activation batch | Systems Architect; Software Architect | reasoning |
| Plan of a milestone; plan of a task | Software Engineer | reasoning |
| Implement | Software Developer | by the task's class |
| Plan review; verification | QA Engineer, on another harness | reasoning |
| Answer to a question | the owner of its kind | reasoning |

**The owner map** (Problem 1, Sol-32). Default: domain or business → Domain
Expert; architecture boundary → Systems Architect; interface contract → Software
Architect; environment or infrastructure → Software Engineer. No fact supports a
default, so the Intake form shows it and the Operator confirms or changes it; the
confirmed map, with the comment ID as its evidence, is a register on the records
branch (Invariant 4).

**The harness register** is on the LAYUP host, one row per harness: its command
template, its credential route, its rule-file names and policy paths (§4),
whether it reports tokens and whether it can enforce a spend cap (§12), and its
models. At Intake, `layup run` runs a fixed probe session on each harness (its
version, a result file, the instruction files it reports it loaded); the admitted
harnesses and their probe results go to the records.

**Routing** is a register on the records branch: for each role and tier, an
ordered list of harness and model pairs. Code admits a pair only when its harness
passed the probe, its model is not on the target's "not used" list, and, for a
plan review or a verification, its harness differs from the author's harness as
the session rows of the ledger record it (Invariant 9; never a field an agent
writes). With fewer than two admitted harnesses, a verification is `not-active`,
and the change does not merge (Invariant 5). Among the admitted pairs, the
learned weight of §13 ranks them; with no weight yet, or a tie, the smart-if's
fit point (§10) or the table's order picks one.

**The tier** of an implementing task comes from its class: the plan session
gives each task a size class, and the computed position of §11 (a task with open
unknowns is uphill) sends an uphill task to the reasoning tier and a downhill one
to the execution tier (vision 2.2, Sol-23).

**The context of a session** (vision 3.1). The step table names, per step, the
record kinds that go into the prompt file: the task's issue and plan, the rows and
specification sections of its requirement IDs, the decision records those
sections name, the task's handoffs and answers, and the target's rule files. Code
selects them by these links, estimates the size (bytes divided by four) against
the model's context size in the register, and refuses a start that does not fit:
the step fails with "context too large", and the plan session splits the task.

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
| S1 Problem Statement Quality (`F-0003#41`), `REQ-001` | [W-01](walkthroughs/W-01-problem-statement-quality.md) | 5 | 0014, 0015 |
| S2 Reproducible Discipline Setup (`#42`), `REQ-002`, `NFR-003`, `NFR-006`, Invariants 4 and 8 | [W-02](walkthroughs/W-02-reproducible-discipline-setup.md) | 5, 6 | 0016, 0017 |
| S3 Rule Protection (`#43`), Gate Integrity (`#64`), Invariant 3 | [W-03](walkthroughs/W-03-rule-protection.md) | 6 | 0017 |
| S4 Stack-Dependent Gates (`#44`), `REQ-004`, `NFR-004`, Invariants 5 and 7 | [W-04](walkthroughs/W-04-stack-dependent-gates.md) | 6 | 0016 |
| #69 A2 (the target cannot pass its gates without LAYUP) | W-04 step 8; W-12 step 15 | 6 | 0016 |
| #69 B9 (the source of the stack gates) | W-02 step 4 | 6 | 0016 |
| S11 Specification Synthesis (`F-0003#51`), Specification Traceability (`#62`), `REQ-012` | [W-11](walkthroughs/W-11-specification-synthesis.md) | 7 | 0018 |
| Table C 2.1 (PDR, PRD) | W-11 | 7 | 0018 |
| #69 B4 (specification synthesis has no component) | W-11 | 7 | 0018 |
| S5 Role Handoffs (`F-0003#45`), Inter-Role Communication Format (`#59`), `REQ-005` | [W-05](walkthroughs/W-05-role-handoffs.md) | 8 | 0019 |
| S6 Autonomous Clarification (`#46`), `REQ-006` | [W-06](walkthroughs/W-06-autonomous-clarification.md) | 8, 9 | 0019, 0020 |
| S7 Verification on Every Change (`#47`), Structural Conformance (`#58`), `REQ-007`, the independent verification of `#66` | [W-07](walkthroughs/W-07-verification-on-every-change.md) | 8, 9 | 0019, 0020 |
| #69 A3 (every question goes to the Operator) | W-06 | 8, 9 | 0019, 0020 |
| #69 B5, table C 2.2 (squads, counterpart harness) | W-07 step 7 | 9 | 0020 |
| Table C 3.1 (model and context routing) | the routing and context rules of §9; the smart-if part is later: slice E | 9 | 0020 |
| Table C 2.1, the phased plan | the bets and the Build plan of §8 | 8 | 0019 |
| ADR-0012 part 6 (set up a target, run its gate from outside) | W-02; W-04 steps 5, 9 | 5, 6 | 0016 |

## 15. Known limits

Each limit is a finding that the design does not close, recorded here (O-66).

- **L-A1. The host is shared.** Role sessions run under the Operator's user, so
  a session that searches the host can reach the App's private key, which does
  not expire until the Operator revokes it, makes tokens for every target where
  the App is installed, and passes the records rulesets; and the Operator's own
  `gh` login; a comment made with that login passes as a human decision.
  Until sessions run in an isolated environment (a container or another
  operating-system user), the separation of O-77 is by convention for that case
  (K03, K04).
- **L-A2. One forge.** GitHub is the only forge for the pilot (#69 forge
  question).
- **L-A3. One host during delivery.** `layup run` runs in the foreground on one
  host; while the host is down, nothing moves (section 11 says how the stall is
  found). A takeover on another host needs the App's private key on that host. A
  run that paused between the records push that announces a forge write and the
  write itself still makes that one write after a takeover.
- **L-A4. The bypass list is read once.** The App cannot read a ruleset's
  bypass list, so it is read only at setup, from the Operator's command; a later
  change to it is not seen. The forge's rule-suites API, which reports a bypass
  after the fact, needs the administration permission, which the App does not
  have.
- **L-A5. A policy file of the host.** A harness that always loads a system-wide
  policy file gives its sessions rules that another harness does not get; code
  records the file, and does not remove it.
- **L-B1. The toolchains on the host.** `layup gate` needs each target stack's
  toolchain on the LAYUP host (K17).
- **L-B2. The coverage floor.** The test-quality gate of the Go entry has no
  coverage floor until the idea owner or the Operator sets one with evidence;
  until then it runs the tests and checks no floor, and the setup record lists
  the floor as an open gap.
- **L-B3. Gates with LAYUP absent.** Once the `layup/` checks are removed, the
  target's own CI runs the workflow and `docs/gates.tsv` of each pull request's
  head, so a pull request can weaken a gate that it is checked by; review is then
  the only guard (FT4). An organisation ruleset that requires a workflow from a
  protected repository would close this, and needs an organisation.
- **L-C1. A fact classed wrongly.** The numbering session can class a need as
  context. The completeness review on another harness and the idea owner's
  confirmation at each bet, which shows every fact with its class and text and
  takes a new class by one line, are the two checks; a need that both miss has no requirement, and nothing
  after them finds it.
- **L-D1. Clarification Turnaround.** A question costs a new session of the owner
  role and a new attempt of the asking role; while a session start takes minutes,
  the start value of 120 seconds at the 95th percentile (`F-0004#15`) is out of
  reach. The measure is still recorded (§12).
- **L-D2. Context by link.** A session gets the records that the step table's
  links reach; a relevant record that no link reaches is not in its prompt.
