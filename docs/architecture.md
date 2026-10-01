# LAYUP architecture

LAYUP takes a problem statement from an idea owner and drives the delivery of
that product, from the first question to the last accepted requirement, in a new
repository that it sets up with the discipline baseline, pinned at its latest
state when the target is set up (§5; LAYUP's own pin is
[`setup/armature.pin`](setup/armature.pin)). It carries the work from the accepted requirements to the produced solution
architecture, technical specifications, and features—verified by a human at each
bet (§8)—all the way to working, verified software (O-101, O-104, O-105). It is a deterministic
orchestrator of the whole lifecycle (O-67): code runs the flow, harness sessions
do the work that needs judgement, a smart-if provider picks a branch at named
points, and humans decide only at the Human Decision Points of the PSB.

This document was written in slices (plan v2 of task `T-hbw8`; approved, O-112);
sections 5 to 13 are slices B to H: 5 Intake and setup, 6 native gates and rule
protection, 7 the specification, 8 the phase loop, 9 squads and routing, 10
decisions, 11 stalls, 12 cost and the measures, 13 learning. Each section
names the ADR that decides it and the walkthrough that tests it. The walkthroughs
are in [`walkthroughs/`](walkthroughs/README.md). The Operator's decisions O-66 to O-113 are quoted in [`runs/T-hbw8/operator-decisions.md`](../runs/T-hbw8/operator-decisions.md)
and [`runs/T-hbw8/inputs-from-pr-69.md`](../runs/T-hbw8/inputs-from-pr-69.md).
LAYUP's own technical specification, in the form that §7 asks of a target, is
[`spec/`](spec/README.md).

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
target passes its own gates with LAYUP absent (Invariant 2). The issues, comments,
rulesets and the ref `layup-probe` are forge objects, not repository content. The
harness register and the price list stay on the LAYUP host; each value that a run
uses is copied into the records (the session start row, the ledger row).

**The LAYUP host** holds, per target: LAYUP's own clone of the target, one
separate clone per role session, the harness and provider registers and their
credentials, and the App's private key, in a file only the Operator's user can read (mode 0600), from which
`layup run` makes installation tokens. None of these is project state: each run rebuilds the clones from
the forge and the records (Invariant 1).

**The forge** is behind a forge interface (O-102): the core engine names
capabilities, never one platform's API. The capabilities: issues and comments
with the actor and whether an App made it; pull requests with a draft state;
commit statuses bound to a source; branch rules with bypass actors, read back;
the repository activity with its actors; an App identity for LAYUP with scoped
permissions. GitHub is the default forge for the pilot and its only adapter so
far (its calls: the REST API, and the GraphQL call that marks a draft ready); the
GitHub terms in this document (rulesets, `performed_via_github_app`, the activity
API) are that adapter's. An adapter for another forge (GitLab, Bitbucket, Forgejo)
maps each capability. All six are required: an adapter that cannot give one
stops the setup at Start or at the probes of §3, and the stop names the
capability. No part of LAYUP runs on a forge without them. No other adapter is designed (known limit L-A2).

## 2. The components

| Component | What it is | Where it runs | What it may write |
| --------- | ---------- | ------------- | ----------------- |
| The engine checks | `layup psb check`, `layup setup verify`, `layup gate`, `layup spec check`, `layup report`, `layup learn` read files; `layup audit` reads files and the forge's read-only API; each prints a typed table | the LAYUP host | only standard output |
| The orchestrator | `layup run TARGET`: the phase loop (section 8), in the foreground, one run per target | the LAYUP host | the records branch, the forge, the session clones |
| The decision component | the smart-if: a package of `layup run` that asks a registered provider at named points (section 10); the providers are declared in the provider register, as harnesses are (O-106) | inside `layup run` | a result that `layup run` records |
| The squad manager | the harness register, admission and routing (section 9) | inside `layup run` | the same |
| The learning loop | the reward and the routing update at each retrospective (section 13) | `layup run`, at the retrospective | the same |
| The runner | the LAYUP App of the forge adapter (for the pilot, the GitHub App) and the LAYUP host | the forge; the host | through `layup run` only |
| The dead-man job | a scheduled workflow in the Operator's control repository, as the App `layup-watch` (section 11) | the forge | a notice on a target's control issue |
| The target | the product repository | the forge | — |

**One run per target** (O-103) means one orchestrator process per target at a
time: no two `layup run` instances drive the same target. Inside that run, the
sessions of different tasks can run in parallel as the milestone plan's order
allows (§8: tasks whose predecessors have merged, up to `build.parallel`); only
the merges are serial, first ready first. `layup run` holds a lease row on the records branch: run
ID, host, start time and a heartbeat counter that it increases every `lease.H`
(a parameter, section 10). A second run watches the branch and times the lease
by its own clock from the moment it last saw the counter change, so the two
hosts' clocks are never compared (the pattern of Kubernetes leader election). It
stops while the counter moves; when the counter has not moved for `3 × lease.H`
of its own time, it takes the lease over and records the takeover. The time with
no run is a stall of the orchestrator (§11).

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
of tab-separated values with a header row, one file per record kind (per task,
one events table in the task's directory; an event
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
| The dead-man job | nothing; it adds a notice on the control issue (section 11) | — |

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
- **Humans.** The accounts in `approvers.tsv`, each by the forge's numeric user
  ID and a role: Operator, idea owner (both from Start; they can be the same
  person), or approver (named by the idea owner at Intake for a planned point).
  Each rule that takes a human's comment names the role it needs: the Intake
  answers only from the Operator and the idea owner, an acceptance only from the
  idea owner, a planned point's answer from its named approver.
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
to the LAYUP App, its only bypass actor. At the Scaffold phase and at each start after it, `layup run`:

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
pushes that SHA to `task/<task>/<attempt>` with the App's token. Before the push, it refuses a branch whose diff from the base touches
`.github/workflows/`: the App has no workflows permission (O-92), so LAYUP never
delivers a workflow change; the result fails with that reason, and the change
becomes a proposal for a rule batch. A rule batch task's branch is the exception:
it is not refused; `layup run` posts its diff, the Operator pushes it (O-93), and
code checks that the pushed head's tree equals the session's recorded tree before
it binds the SHA. The diff of the
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
   on it and on the dead-man job's App (§11), adds the target to that job's list,
   and runs `layup run --new OWNER/NAME --psb FILE [--vision FILE] --operator
   LOGIN --idea-owner LOGIN --plan PLAN --intake-cap MONEY,HOURS --lease-h MINUTES --watch-t MINUTES`. The
   command's values are the evidence of the values that must exist before the
   Intake answers: the cap of Intake and Shape (milestone 0, §12), `lease.H`, and
   each harness's cap and wall-clock limit, which the harness register gives as
   the Operator set them there.
2. `layup run` resolves the latest commit of the baseline's default branch at
   that moment (`git ls-remote`), clones that commit, removes `.git` (ADR-0011
   decision 7; step S02 changes from `npx degit`), and records it as the target's
      own pin (O-101): source, commit, tree and time go into the first records
   commit (step 3); step S04 of the Scaffold writes the target's
   `docs/setup/armature.pin` from that record, on the setup branch, so the root
   commit stays unmodified. Once that record exists, no run resolves the commit again; a run that stops
   before it (for example at the plan check of step 3) leaves a root commit with no
   record, and the Operator starts again with an empty repository. LAYUP's own pin does not bind a target. The printed push command shows
   the resolved commit and its difference from LAYUP's own pin, so the Operator
   sees what the push brings.
   The Operator pushes this unmodified copy as the root commit of the default
   branch, with the Operator's own login and one command that `layup run`
   prints: the copy holds CI files, and the App has no workflows permission
   (O-92). So the default branch is the first branch of the repository, and its
   root tree equals the pinned tree, as step S03 asks. The command pushes from a
   plain clone with no hooks installed.
3. `layup run` reads back the repository's default branch, its root tree and its
      visibility, records the LAYUP version and the plan that the Operator names with
   `--plan`, and stops when one differs, when the forge adapter lacks one of the
   six capabilities of §1, or when the plan does not enforce rulesets or offer
   draft pull requests on it (a private repository on GitHub
   Free has neither, and on GitHub Pro no drafts: the Operator makes it public or
   moves it to a plan that has both, K15); the probes of §3 at Scaffold prove it. It then pushes the records branch (an orphan, §3)
   with its first commit: the target's pin (step 2); the problem statement and the
   vision brief byte for byte, each with its SHA-256; `approvers.tsv` with the
   numeric IDs and roles of the two logins of step 1; and the lease row. From this
   commit on, only a comment by an ID in `approvers.tsv`, in the role that the
   rule names, is an answer or a decision (§3); any other comment is recorded as
   input. It opens the Intake issue and the control issue (§10), and
   waits up to `watch.T` (a Start value) for the dead-man job's first notice there
   ("watch started"), which shows that its App and list work; without it, it goes
   on and records "watch not confirmed" (L-F1).

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
   the appetite and the band of §12, the approver of each planned approval point
   except the acceptance of a requirement, which is always the idea owner's (PSB
   §8); a login named there is added to `approvers.tsv` with the role
   "approver", never "idea owner"), and the setup questions, the proposed marker sources and the
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
the pinned baseline's: neither goes into a target (FT5). `layup setup verify` does, from outside and for this target, the checks of
`setup-check.sh` that apply to one: `pin`, `kit-history`, `facts`, `onboarding`,
`glossary`, `guardrails`, `markers`, `adapted` and `identity` (with the target's
name and facts), and the pinned baseline's own `link-lint.sh`. Three do not
apply: `ci` and `procedure` read LAYUP's own CI and `steps.tsv`, and `protection`
reads the classic protection that the rulesets replace. The evidence of each step
of S04 to S14 becomes "`layup setup verify <check>` OK", except S12 (the native gate jobs exist,
one per kind) and S13 (the rulesets read back, §3); step S15 writes the setup
record into the records branch, not `steps.tsv` into the target.

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
      discipline tests on a checkout of it; the checks that apply to a target (pin:
      the root tree equals the pinned tree; and the others listed above);
   zero values without a source; each source resolving; and each active gate,
   which must pass on the clean tree and fail on its known-bad fixture (§6).
6. The Operator pushes the setup commits on top of the root commit and applies
   the two rulesets of §6 from a file that `layup run` writes,    with the Operator's own login and the commands it prints, from a clone with no
   hooks installed: the setup changes CI files,
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
prints "not comparable", never a pass. The idea owner's batch of start values follows the baseline (§12).

## 6. Native gates and rule protection

### The stack gates

[ADR-0016](adr/0016-put-the-native-stack-gates-in-the-target.md) decides this.

**The stack catalog** is in LAYUP's repository, one directory per stack. An entry
lists, per gate kind (layout, interface boundary, contract, test quality, and any
further kind of the stack, such as the Go entry's static checks): the
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
name); `test -z "$(gofmt -l .)"` (bare `gofmt -l` exits 0 on a bad file) and `go vet
./...` are a fifth kind, "static checks", active from the setup (a stack gate
adds rules, Invariant 7); at activation, an
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
| `pending`, and the head changes no path in the product's scope | `clear`: counts as a pass, with the reason "pending: no product path" (the kind's own gate did not run; the rule that ran is "no product path may change while this kind is pending") |
| active, and the tree has no path in its scope (for Go: no package) | `clear`: counts as a pass, with the reason "no product path" |
| `pending`, and the head changes such a path | failure |

The status is a success only when every kind counts as a pass. It never runs the
head's gate files for a status (FT4), with one exception: an approved rule batch,
whose rule files are the approved ones (below). Before its approval, a batch's own
gate files run only in a scratch work tree, as evidence for its verifier and its
brief, never as a status. The LAYUP host needs each stack's toolchain
(L-B1). With LAYUP absent, the target's own CI runs the workflow and the manifest
of the pull request's head; FT4 then rests on the rulesets and on review, not on
LAYUP.

**Activation.** Each rule batch that adds or changes a gate kind carries one new
known-bad patch per kind it touches, and may replace a recorded patch of such a
kind that no longer applies (the old one is retired with the approval). Before
its approval, the scratch run above shows the verifier and the brief that the
head passes and each patch fails. After it, `layup gate` runs the batch's own gate
files on its head (which must pass) and on the head with every recorded known-bad
patch of each kind it touches (each must fail) for the status; a patch that
passes, or that no longer applies and was not replaced, refuses the merge (§13). At the
first bet (§8), an architect session writes the
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
request, so the approver sees the change before deciding. A batch that changes
`.github/workflows/` cannot be pushed by the App (O-92): its sessions write it,
`layup run` posts its diff as a payload, and the Operator pushes the batch branch
(a human step of that planned point); then the target's CI and the verifier run
on it as on any task. The approval request
names the batch head's SHA and the tree hash of its rule files. The approver
approves by an issue comment (§3); code refuses the approval when the batch head
has moved since the request, and records the hash with the approval. The approved batch merges by
`layup run`, or, when it changes `.github/workflows/`, by its approver, after the Operator pushed its branch
(below), because the App has no workflows permission (O-93).

**Prevention, in three layers:**

1. **No credential.** A role session holds no forge credential (§4), so it
   cannot push, merge, approve or change a ruleset with a token.
2. **No rule change leaves the host in a task.** Before `layup run` pushes a task
   branch, it compares the changed paths with the rule-path register of the base
   branch. A change to a rule path refuses the result: the task gets the finding,
   and the change becomes a proposal for the next batch. The two exceptions are
   those of `layup/rules` below: the branch of a rule batch task, and a rendered
   change; `layup/rules` then gates their merge.
3. **The forge refuses.** The default branch and the probe ref `layup-probe`
   have one ruleset with an empty bypass list: a pull request is required; the
   required checks are the native gate jobs, pinned to GitHub Actions, and
   `layup/gates`, `layup/spec`, `layup/verify` and `layup/rules`, pinned to the
   LAYUP App; force
   pushes and deletion are blocked. The records branch has a ruleset that
   restricts updates and deletion to the LAYUP App, its only bypass actor (§3).
   `layup/rules` fails a pull request that changes a rule path unless it is an
   approved batch (the hash of the rule files that the batch itself changes, at
   its head, equals the hash recorded with its approval), or a rendered record
   that equals code's rendering of an inventory version that a bet or a decision
   wrote (§7, §8).
   (`layup/verify` comes from the counterpart verification of §8.)

**Reading of `F-0003#64`.** "0 agent writes to rule paths" counts the rule-path
changes on the default branch that did not land in an approved batch or as a
rendered record of an approved bet or escalation decision; the approval makes the change the human's
decision (the approval brief asks the
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
   (for example `need`). A new class reopens the draft for that fact in a new specification task after this
   one merges (step 6); its requirement is written under a PRD section "waits for
   a bet", not as a table row, so `prd-lint` does
   not read it, and it gets its row and its line at the next bet. Priority is intent
   (`F-0003#54`), so code copies the comment (§3) and adds one **rendered commit**
   to the verified specification's head: the MoSCoW and Phase columns from that
   copy, the out-of-scope fact IDs under the PRD's non-goals, and the numbered
   facts record (step 8). Code checks that the commit changes only those, and
   carries `layup/verify` to it, as for a close-out (§8); only then does the pull
   request pass the baseline's own `prd-lint`. A requirement
   added later gets its priority at the next bet. Inside a requirement, a session
   may mark a task nice-to-have; a cut that touches an acceptance criterion is
   business-forking (§10).
8. **The confirmed inventory.** At each bet, code writes a new version of it
   from that bet's copy: the numbered facts, their classes and the out-of-scope
   marks, on the records branch, with the SHA-256 of the version. An escalation decision that changes a requirement or a priority (§10) also
   writes a version, from its own copy. The numbered facts record of the
   target, with the confirmed classes, is rendered by code from the version: at
   the first bet on the specification's head (step 7), later as a rendered task
   (§8); `layup/rules` passes it when the file equals the rendering. The Scaffold
   writes only the problem statement and the answers into `docs/facts/`. `docs/facts/` is a rule path from the
   setup, so a role session cannot edit it.

**`layup spec check`** (code, no model call) posts its own required status,
`layup/spec`. It reads the latest confirmed inventory, whose hash it checks, and
the bet and decision copies (records), the
PRD and `docs/spec/` (the head), the task register (§8) and the accepted and
delivered requirements (§12). It fails when: a `covers` ID does not resolve to a
fact of the confirmed inventory or of the answers; a need or constraint has no
requirement, or a measure has no success-criterion row that cites its ID, and it
has no out-of-scope mark and no requirement marked "waits for a bet"; a
requirement has no acceptance criterion; a MoSCoW or Phase value in the PRD
differs from the latest copy; a delivered
requirement has no non-empty section in `docs/spec/` whose heading holds its ID
(Sol-20); a task names no requirement; or it did not run. **Before the first
bet**, no inventory is confirmed: a task that changes the PRD or `docs/spec/` is
checked against the Intake's numbered facts without the MoSCoW, Phase and
out-of-scope rules; any other task is `clear` ("no specification yet"). It proves the links. It
does not prove that the meaning agrees: the counterpart verification of each
change judges that (§8), and the completeness review of step 4 is the check
against a missed need.

**Success criteria.** The PRD's success-criterion rows are the one source (§12).

**Vision 2.1.** The preliminary design review is the specification of step 6 with
the architecture that the first bet approves (§8); the phased plan is the
milestone plan of the bets (§8).

## 8. The phase loop

[ADR-0019](adr/0019-run-the-lifecycle-as-a-phase-loop-with-one-bet-per-milestone.md)
decides this section.

**The phases.** Intake and Scaffold (§5) → **Shape** → **Bet** → **Build** →
**Accept** → **Retrospective** (§13) → the next **Bet**, until every `Must`
requirement is accepted. A later milestone goes back to Shape only when its bet
brief proposes a change of the architecture.

**Every pull request is a task.** The specification, each rule batch, each
rendered record, the milestone plan and each build task are tasks, and each goes
through the task loop below: its issue, its plan and plan review, its
verification and its close-out. So each one gets the target's own records and
`layup/verify` (Sol-27, Fable-M9). For a **rendered** task (a later bet's or an
escalation decision's change to the PRD or the facts record, §7, or the pilot's
start values, §12), code writes the plan ("render from copy N") and the change; a plan-review session and
a verifier session check them as for any task.

- **Shape.** Tasks write the specification (§7 step 6), the target's
  architecture (its document and decision records, by the target's own rules),
  and the gate-activation batch (§6). **Solution routing** (vision 3.1): for each
  architecture decision, an architect session writes at least two options, each
  with a fixed table: every constraint that applies (the requirements' criteria,
  the confirmed constraints, the PSB invariants, the accepted decision records),
  and for each a pass or fail with its evidence. The table is a **filter**, applied by
  code: an option with a failed row, a missing row, or a pass without evidence is
  dropped. With one survivor, it is chosen. With more, the blind panel of §11
  compares them (the trade-offs that the constraints do not decide) and its
  synthesis recommends one. When every option is dropped, the architect session writes new options once; if
  they are dropped too, the round made no progress, so trigger 1 of §11 opens a
  stall. With fewer than two harnesses the panel does not run, and a panel that ends
  `insufficient panel` gives none; the brief then shows the survivors with no
  recommendation, and the bet decides. The bet
  brief shows the survivors, the dropped options with the failed rows, and the
  recommendation; the bet decides.
- **Bet** (Decision Points 1 and 3; Shape Up's betting table, one bet per
  milestone). A session of the Product Owner role writes the brief in five parts:
  the problem (the requirements of the milestone), the appetite (its cap of money
  and wall-clock, inside the band, §12), the solution (the architecture, or a
  proposed change of it, which sends the milestone back to Shape), the rabbit
  holes and the no-gos, with the fixed line "architecture change: yes or no"; and,
  at the first bet, every fact with its class (§7 step 7). Code checks the five headings, a line limit (a parameter), and that each
  link resolves; the brief names the head SHA and the rule-file hash of each
  verified batch that it asks to approve. `layup run` posts it on the milestone's
  issue. The approver of that planned point (the idea owner by default, from
  Intake) answers by one comment: a line per requirement, a line per fact to
  change, and "bet" or "no bet". The bet approves the milestone, the architecture
  and its batches; there is no other approval inside the milestone, apart from
  escalations (§10) and stalls (§11). A bet that changes a class or a priority
  does not touch the approved batches: at the first bet, code adds the rendered
  commit to the specification's head (§7 step 7); at a later bet, it renders the
  change as a rendered task. **Merge order after a bet:** the specification (or
  the rendered task), the rule batches, then the build tasks. A requirement changes only
  at a bet, or by an escalation decision (§10).
- **Build.** A plan task splits the milestone into build tasks, each with the
  requirement IDs it serves, the tests that will show it done, a size class, and
  its predecessors;
  code checks that each `Must` requirement of the milestone has a task and each
  task a requirement, and writes the task register. Code also checks that each predecessor is a task of the register and that the
  predecessors form no cycle. A task whose predecessors have merged may start; at
  most `build.parallel` tasks (a parameter) run at once. The **merge order** among
  build tasks is first ready: the first task to pass step 5 of the loop is next,
  with ties by the register's row order; a task is ready only while its current
  head has passed step 5, and a task that waits (a question, an escalation) does
  not hold up a ready one. Tasks merge one at a time (the branch rule is "up to date"): after a
  merge, `layup run` merges the base into the next task's branch; a conflict is
  resolved by a new attempt. **A verdict and a moving base.** A build task is
  verified when it is next in the merge order and up to date with the base (step
  6); until then its `layup/verify` is posted `pending` with the reason "waits for
  its turn", a wait state (§11). Shape tasks and rule batches are verified before
  their bet, and rendered tasks when they are made. A verified task whose base
  then moves carries its verdict over a clean merge of the base that changes no
  file it touches, as the target's own rule allows; otherwise it is closed and redone as a new task, with its own issue and records,
  which takes the old task's row and predecessor links in the register, so the target's review
  record never needs a round after `nothing material in scope`.
- **Accept** (Decision Point 3). A requirement is **delivered** when every task
  that names it in the task register has merged. `layup run` then posts the
  requirement, its criterion, the merged changes and the gate and verification
  results; the idea owner answers "accept" or "reject" with a reason. Code takes
  the answer only from an account whose role in `approvers.tsv` is idea owner
  (PSB §8: the idea owner accepts the delivered requirements); an approver of
  another planned point cannot accept a requirement. A rejection
  becomes a new need for the next bet (Shape Up: a new bet, not a patch). The
  record is in §12.

**The task loop.** Each step of the target's own quality gate has an actor:

1. `layup run` opens the task's issue (its requirement IDs, its Definition of
   Done).
2. A plan session writes the task's plan; a plan-review session on another
      harness reviews it. `layup run` posts both on the issue, in the forms that the
   target's own checks parse, and reads the verdict from its fixed field. The plan lists its tests by ID; the list is frozen at the task's first valid
   handoff to the verifier (§11).
3. A developer session writes the failing test first, then the change; it ends
   with its typed result (§4).
4. The handoff check (below), the rule-path and workflow checks (§4, §6); then
   `layup run` pushes the branch and opens a **draft** pull request that links
   the issue, with a title in the target's form (`<type>: <task ID>
   <description>`), and requests no review (Sol-2).
5. The target's CI, `layup gate`, `layup spec check` and `layup/rules` report, on
   each new head (the close-out's included). A failure goes back to the developer
   as a finding, in a new attempt. Two failures are expected and do not go back:
   on a rule batch before its approval, `layup/rules` and the pending kinds of
   `layup/gates` (the scratch run of §6 is its gate evidence instead); on a
      specification that waits for its bet, `prd-lint` on the empty MoSCoW and Phase
   columns only: code runs the pinned `prd-lint` itself and accepts the head only
   when each of its failure lines names an empty MoSCoW or Phase cell.
6. When they pass (apart from the expected failures of step 5), and, for a build
   task, it is next in the merge order and up to date with the base (§8 Build), a
   verifier session reviews the head: on a harness that wrote
   no commit of the change (§9), fresh, one turn, read-only, told to refute
   "done", with a `file:line` checklist. `layup run` posts its review record on
   the task's **issue**, in the target's form and with the target's verdict
   values, then edits the pull request's body to add the link to it, which makes
   the target's `review-record` check run again. It sets `layup/verify` at that
   head SHA: success only for `nothing material in scope`. A material finding is
   one more round, up to the target's cycle cap; past the cap the round limit of
   §11 opens a stall.
7. **Close-out.** A close-out session writes the target's close-out (the
   verdict, the completed-log line, and the resource record from the ledger,
   §12) as one last commit. Code checks that this commit changes only the task
   file and the completed log; `layup run` then sets `layup/verify` on the new
   head, which differs from the verified one by that commit alone.
8. When every required check is green at the head, `layup run` marks the pull
   request ready and merges it at that head SHA (a batch that changes
   `.github/workflows/`, which the Operator pushed, is merged by its approver,
   §6). No human review is requested: a
   routine review of each pull request is not a planned point (PSB §6). A human
   who reviews anyway is recorded as human input (§12). **Reading of
   `F-0003#58`:** a pull request "reaches human review" when a review is requested
   or it is marked ready; a draft is only visible (the approval brief asks the
   Operator to confirm this reading).

**Handoffs** (`F-0003#45`, `#59`). A handoff is the typed result of one session
for the next step: session ID, task, role, status (`completed`, `blocked`,
`needs_context`, `decision_needed`, `failed`), reason, each artifact with its
path and SHA-256, the questions, the proposed decisions, the lessons (Spec Kitty's
handoff packet v1 as a pattern for the fields). The **transition table** is a
register on the records branch; its default rows use the target's own record
kinds (the plan, the plan review, the review record, the task file):

| From | To | Required, checked by code |
| ---- | -- | ------------------------- |
| task plan | plan review | a plan in the target's form that lists its tests by ID |
| plan review | developer | a review in the target's form with the verdict `approve` or `approve-with-conditions` |
| developer (build task) | verifier | commits; the task file; the plan's tests fail at the base and pass at the head (`layup gate` runs them) |
| developer (specification) | verifier | commits; the task file; a section in `docs/spec/` for each requirement it names |
| developer (rule batch) | verifier | commits; the task file; a known-bad patch per kind it touches; `layup gate` passes on its head and fails on each patch |
| developer (milestone plan) | verifier | the task register rows; each `Must` requirement of the milestone has a task; each predecessor is a task of the register; no cycle |
| code (rendered task) | verifier | the change equals code's rendering of the named copy: an inventory version, a decision copy, or the start-values copy (§12) |
| verifier | close-out | a review record with `nothing material in scope` at the head |
| verifier | developer | a review record with each finding at a `file:line` |
| specification | milestone plan | a section in `docs/spec/` whose heading holds each requirement ID of the milestone |
| asking role | owner role | a question row with the records it concerns |
| owner role | asking role | an answer row that cites the records it rests on |

A handoff is valid only when code finds the schema right, each named artifact in
the commits with its hash, and each required item true; a pair with no row is not
valid; a field that an agent writes never makes it valid (Sol-21, Fable-M14). The
Operator can replace the table (a parameter, §10). The Inter-Role Communication
Format is the valid handoffs over all handoffs, from the records.

**A question during the work** (Problem 1). A session that meets a question ends
with status `needs_context` and the question: its text, its own label for the
kind, its `needs_human` field, and the records it concerns. `layup run` records
it with its time, posts
it on the task's issue, and gets its kind and whether it needs a human from the
smart-if (§10). A question that needs no human goes to a session of the owner role
of its kind (§9), which answers with the records it cites. The answer passes the
escalation screen (§10), is recorded and posted, and the asking role starts its
next attempt with the answer in its prompt file. The answer is **accepted** when
that attempt ends with the status `completed`, its result cites the answer ID, and
none of its questions cites the answer ID; the accepting actor is that session.
An attempt that a question ends counts toward `stall.attempts` (§11); it is a
wait for the owner role, not a round without progress.

**Attempts.** A new attempt starts from the base commit, with the last attempt's
findings and any diagnosis in its prompt file, not its diff, so that its commits
have one author harness. Once the task's test list is frozen (§11), code applies
the frozen test files to each new attempt's start, so their hashes stay; the
harness that wrote them counts as an author of every later attempt (§9); a change
to a frozen test needs an amended plan, which the plan session writes on the
finding "test changed", and a new plan review, which freezes the list again.
When the Operator passes the diff too, by a stall answer, the harness that wrote
it counts as an author of the new attempt (§9).

## 9. Squads and routing

[ADR-0020](adr/0020-route-role-sessions-over-registered-harnesses.md) decides
this section.

**Roles and steps** (O-81). The roles are the seven functions of PSB §2 by
default, and the Operator can replace the matrix per target. The default step
table:

| Step | Role | Tier |
| ---- | ---- | ---- |
| Intake gap review; completeness review (§7 step 4) | Systems Architect; QA Engineer | reasoning |
| Numbering, draft, specification (§7); numbering and draft on one harness | Product Owner | reasoning |
| Bet brief | Product Owner | reasoning |
| Shape: the architecture; the gate-activation batch | Systems Architect; Software Architect | reasoning |
| Plan of a milestone; plan of a task | Software Engineer | reasoning |
| Implement; close-out | Software Developer | by the task's class |
| Plan review; verification | QA Engineer, on a harness outside the authors | reasoning |
| Answer to a question | the owner of its kind | reasoning |
| Retrospective lessons | Systems Architect, on a harness outside the milestone's authors where one exists | reasoning |

**The owner map** (Problem 1, Sol-32). Default: domain or business → Domain
Expert; architecture boundary → Systems Architect; interface contract → Software
Architect; environment or infrastructure → Software Engineer. No fact supports a
default, so the Intake form shows it and the Operator confirms or changes it; the
confirmed map, with the comment ID as its evidence, is a register on the records
branch (Invariant 4).

**The harness register** is on the LAYUP host, one row per harness: its command
template, its credential route, its rule-file names and policy paths (§4),
whether it reports tokens, whether it can enforce a spend cap (§12), whether it sends hook events (§11), its models, and each model's context size with its source (the model's
documentation). At Intake, `layup run` runs a fixed probe session on each harness
(its version, a result file, the instruction files it reports it loaded); the
admitted harnesses, their versions and their probe results go to the records.
Before each session, code reads the harness version; a new version is probed
again before it runs.

**Routing** is a register on the records branch: for each role and tier, an
ordered list of harness and model pairs. Code admits a pair only when its harness
passed the probe, its model is not on the target's "not used" list, and, for a
plan review or a verification, its harness is not an author: for a plan review,
the **authors** are the harnesses of the task's plan sessions; for a
verification, every harness whose session is bound, in the ledger, to a commit in
the diff from the base, every harness whose diff went into an attempt's prompt,
and the harness that wrote the task's frozen tests (Invariant 9; never a field an
agent writes). When no
admitted harness is outside the authors, the verification is `not-active`, and
the change does not merge (Invariant 5). Among the admitted pairs, the learned
weight of §13 ranks them; with no weight yet, or a tie, the smart-if's fit point
(§10) or the table's order picks one.

**The tier** of an implementing task starts from its plan, which has four fixed
fields: `size` (`small` or `large`), `open questions`, `new interfaces`, and `new
dependencies` (read by the escalation floor, §10). The
reasoning tier when `size` is `large` or any of the three lists is not empty; the execution
tier otherwise. The progress position of §11 moves it later: an uphill task (open unknowns) to the reasoning tier, a downhill one to the execution tier.

**The context of a session** (vision 3.1). The step table names, per step, the
record kinds that go into the prompt file: the task's issue and plan, the rows and
specification sections of its requirement IDs, the decision records those
sections name, the task's handoffs and answers, and the target's rule files. Code
selects them by these links, estimates the size (bytes divided by four) against
the model's context size, and writes both numbers in the session start row. A
start that does not fit is refused: for a build task, the plan session splits the
task; for any other step, the step fails, and §11 handles it.

## 10. Decisions: the smart-if, the escalation screen, the parameters

[ADR-0021](adr/0021-branch-at-named-points-through-a-smart-if-provider.md) and
[ADR-0022](adr/0022-screen-for-business-forking-decisions-before-the-work.md)
decide this section. The order of decisions is fixed (Invariant 6, O-67): a
deterministic check first; the smart-if only where the question is one of
meaning; a human only at a decision point.

### The smart-if

The smart-if is a conditional branch (O-84): at a named point, the flow takes one
branch or another by the provider's answer. It writes no text and plans no work.

**The provider register** (O-106) is on the LAYUP host, in the same form as the
harness register (§9): one row per provider, with its kind (Jev, TypeSafe's
System One API; Laya, an open-weight decision model with a Jev-compatible API,
run on the host or at an endpoint; or another provider behind the same
interface), its endpoint, its credential route, its model versions, its size
limit, and its price source. Before a provider is used, a fixed probe request
checks that it answers, which model version answered, and its response form (a
paid call, with its ledger row); the
admitted providers and their probe results go to the records, as the harnesses'
do. **The provider** of a target (O-78) is chosen by the Operator at Intake from
the admitted providers, or `none`. Laya needs its own calibration: its confidence
is not Jev's, and its authors say it needs fine-tuning
([`selection-v2.md`](../runs/T-hbw8/selection-v2.md) §1.1). The model version is pinned per target; a new version resets
every `delegate` point to `shadow` until it is calibrated again. A request
carries the point, the literal questions of the point's fixed pack, the options,
and a state text that code builds from named records and checks against the
provider's size limit (a state that does not fit is a failure); code keeps all
arithmetic, dates and counts (the provider's documented weak points). The
response carries each answer (a choice, a probability of yes, or a score), the
provider, the model version, and the tokens and price where the provider reports
them, else `unknown`, never zero (FT2).

**The named points** are a closed list; a call at any other point is a defect
(Fable-M13):

| Point | The question | The deterministic branch | A safer branch | Ground truth for calibration |
| ----- | ------------ | ------------------------ | -------------- | ---------------------------- |
| P1 escalation | four questions, one per PSB axis | the floor and the session's own `decision_needed` | escalate | the idea owner's confirmations; the Missed Escalations sample (§12) |
| P2 question | the kind of a question (four kinds); whether it needs a human | the asker's label and its `needs_human` field | "needs a human" = yes | the owner session's disposition; the Reversal sample (§12) |
| P3 stall action | after the diagnosis: retry, panel, or the Operator (§11) | the fixed ladder of §11 | the next rung up | the stall's outcome (§11) |
| P4 fit | which admitted pair fits the task (§9) | the routing order | none | the verification's first-round verdict |
| P5 over budget | continue inside the band (§12) | stop and escalate | stop | the idea owner's decisions on budget escalations |

**How an answer decides.** A yes-or-no question has a threshold `t` per point: yes
when p ≥ t, no when p ≤ 1 − t, undecided in between. A choice or a score decides
when its top option's probability is ≥ t, and is undecided otherwise. P1 only
ever adds a candidate: no answer of P1 removes a floor candidate or a session's
own declaration.

**The authority** of each point (O-78, O-79) is set by the Operator per target:
`off` (no call; the deterministic branch), `shadow` (the provider answers and is
recorded; the deterministic branch decides), `cautious` (a decided answer may only
move the flow to the safer branch; not for P4), or `delegate` (a decided answer
decides). The Intake form offers `shadow` for every point (Sol-6: an offer, not a
rule), and says what `shadow` leaves open (below); the Operator may choose another
level and a threshold.

**The bounds** (O-71, O-84). Under `off`, `shadow` and `cautious`, a failure (an
error, a time-out, a `429` after its backoff, a state too large) takes the
deterministic branch and is recorded. At a `delegate` point, a failure or an
undecided answer goes to a human (the idea owner for P1 and P5, the Operator for
the others), never to a pass. No parameter turns off a PSB rule: the floor always
runs, and a business-forking decision always goes to the idea owner.

**Calibration** (Invariant 4). A threshold's evidence is the Operator's comment
that set it, or a calibration record: over the rows of a point, the share where
the provider's decided answer agreed with that point's ground truth (the table
above), never with the deterministic default.

**The record.** Each call writes one row to `decisions.tsv`: the point, the
questions, the SHA-256 of the state text, the provider, the model version, each
answer and probability, the tokens and price, the authority, the threshold, the
branch taken, and who decided (the provider, code, or a human).

**What `shadow` leaves open.** Under `off` or `shadow`, P1 adds no candidate and
P2 follows the asker's label: a business-forking choice that a session does not
declare, in prose, is not screened. Each selection row records whether P1 and P2
were active. The PSB's escalation rule (`F-0003#48`) then rests on the floor and
on the declarations (known limit L-E1).

### The escalation screen

A **candidate** business-forking decision comes from three sources (Decision
Point 4, `F-0001#13`):

1. **The session declares it.** Each session's prompt says: before a choice that
   may change the budget, the legal or compliance position, or the approved
   intent, or that trades approved goals against each other (scope against date),
   stop with `decision_needed` and the options.
2. **The floor** (code) reads fixed fields and proposed diffs, never prose: the
   plan's fixed field `new dependencies` (each by its manifest identifier, with
   its stated cost and the alternatives the plan considered); a diff that adds a
   dependency to the stack's manifest (for Go, each `require` in `go.mod`,
   indirect ones too), or touches a licence file, the PRD's requirement rows,
   their priorities or criteria, or a path that the Intake answers named as
   intent. Every new dependency is a candidate unless its identifier is on the
   allowed-dependency list that the idea owner set (at Intake or in an earlier
   decision); nothing that a session writes removes a candidate. For the brief,
   code reads the licence from the module itself (the catalog names the tool).
   Exempt: a diff that equals code's rendering of an approved bet or decision,
   and a specification task of the Shape phase, which a bet approves. The floor runs on every plan and every handoff diff,
   whatever the frequency parameter says.
3. **The four questions** (P1), on prose: each plan, each handoff and each answer
   to a question, as the frequency parameter (O-79) says. One literal question
   per axis: does it change the budget; the legal or compliance position; the
   approved intent; does it trade approved goals. Code combines the answers with
   OR (a deterministic floor that a model may raise and never lower, as in
   Governed APA, `selection-v2.md` §1.1).

**A question that needs a human.** When P2 decides "needs a human", or a session
sets the question's fixed field `needs_human` (§8), the question becomes a
candidate of this screen. At a `delegate` P2 whose answer fails, the Operator
classifies only the kind; a question that needs a human still goes to this
screen.

**One row per screen.** Each run of the screen (a plan, a handoff, an answer)
writes one row to `screens.tsv`, selected or not: the input, the floor result,
whether P1 and P2 were active, and their answers. The Missed Escalations audit
samples these rows (§12).

**When one is selected**, the task stops in a wait state (§11 does not count the
wait), and the candidate goes to `candidates.tsv`. The options come from the
declaration or the plan's `alternatives`; for a candidate that has none, a
Product Owner session writes them. `layup run` posts one escalation brief on the
task's issue for the idea owner: the decision, the numbered options, the
evidence, the cost. The idea owner answers in a fixed form: `option <N>;
business-forking: yes` or `no`, with free text after it. `layup run` copies the
comment (§3) and writes `escalations.tsv`: the chosen option and each rejected
option by name, and an option that is a dependency by its manifest identifier.
An option that changes a requirement, a priority or the band carries lines in
the bet's form (`REQ-7 Won't`; for the band, `band B <amount> U <amount>`), which
the Product Owner session writes and code checks (the IDs exist, the values are
valid); the decision's copy writes a new inventory version (§7) for a requirement
or a priority, and `budget.tsv` for the band. A decision confirmed as business-forking is planned input
(Decision Point 4); one confirmed as not is unplanned input (`F-0001#28`, §12). A
decision that changes a requirement, its priority or the band is written at once,
as a rendered task (§8), like a bet line. The task starts its next attempt with
the decision in its prompt file. A later floor candidate whose manifest identifier equals that of an option that a decision
rejected, anywhere in the target, gives the task a finding, not a new brief; an
option written in prose never matches; any other candidate is a new brief.

### The parameters

Every value that the Operator or the idea owner can set is a row of
`parameters.tsv` on the records branch: its name, value, default, bound, the role
that may change it, the evidence of the value, and the comment that set it
(Fable-M22). `layup run` opens the target's control issue at Start. A change is a
comment there, in the form `set <name> <value> because <reason>`, by an account
whose role in `approvers.tsv` the row allows (the Operator for most rows, the idea
owner for the allowed dependencies and the intent rows; `approvers.tsv` holds a
role per row); `layup run` copies it, checks the
author, the name, the value and the bound, and applies it at the next step
boundary. A change that takes effect on an open task is input to that task
(§12). The bound of each row says what the value cannot do: no parameter removes
a PSB rule (O-84).

| Parameter | Default | Source |
| --------- | ------- | ------ |
| provider (from the provider register); model version; authority and threshold per point | set at Intake; `shadow` offered | O-78, O-79, O-106 |
| escalation frequency (P1 only) | each plan, each handoff, each answer | O-79 |
| owner map; role matrix; step table; transition table | §8, §9 | O-81 |
| harness paid work without tokens or a spend cap | allowed, with a wall-clock limit (§12) | O-80 |
| `stall.T`, `stall.N` | 10 minutes (a maximum), 1 | O-82 |
| `stall.attempts`, `ci.T`, `panel.K` | set at Intake with evidence | O-82, Invariant 4 |
| the reward's term weights; `learn.step`, `learn.min_weight`, `learn.max_weight`, `learn.min`, `learn.explore`; `learn.trigger`; the learning authority | set at Intake, the evidence is the Operator's comment; `learn.trigger` each end-of-milestone retrospective; the authority `propose` | O-83, Invariant 4 |
| allowed dependencies | set by the idea owner at Intake or by a decision | `F-0001#13` |
| `build.parallel` | set at Intake; the evidence is the Operator's setting (the host's and the harnesses' limits) | O-103, Invariant 4 |
| `lease.H`, `watch.T`, the intake cap; the brief's line limit | `lease.H`, `watch.T` and the intake cap at Start (the command is the evidence); the line limit at Intake | Invariant 4 |
| `harness.<id>.cap`, `harness.<id>.wall`; `audit.n` | the harness values in the harness register, which the Operator sets before Start (recorded at Start); `audit.n` at Intake; the evidence is the Operator's setting | Invariant 4 |

The band has one home, `budget.tsv` (§12); a bet or an escalation decision
changes it as a rendered change, never a parameter row.

## 11. Stalls

[ADR-0023](adr/0023-stop-a-stall-at-a-limit-and-diagnose-it-with-a-fresh-context.md)
decides this section. A stall is "a task that does not reach its goal and does
not fail cleanly, because role agents do not agree, or because a step repeats
without progress" (PSB §8). The limits are parameters (O-82, §10); each one that
has no default here is set at Intake with its evidence (Invariant 4).

**Progress is computed, never reported** (Shape Up's hill, computed; Sol-22,
Fable-M4, Author-1). Code keeps two sets per task from the records: the open
**unknowns** (its open questions, open material findings and open escalations;
a finding has a fixed form of file, line and rule, so the same finding in two
rounds is the same unknown) and the **passed tests** of its frozen test list. The list is frozen at the task's first valid handoff to the verifier, when each
listed test fails at the base and passes at the head (§8): each test by its ID and
the SHA-256 of its source then. Later, a test whose source hash differs is no
longer counted as passed, and the change is a finding; an amended plan with a new
plan review freezes the list again and closes that finding. A round (an attempt, or a review round)
makes **progress** when at least one unknown that was open at the end of the last
round is closed, or a test of the list passes that did not pass before. The first
review of a change sets the baseline and is never a round without progress; a
round that ends with a question is a wait for the owner role (§8), not a round,
and counts only toward `stall.attempts`. A new commit, handoff or record alone is not progress. A task is **uphill** until its test list is frozen, and while it has an open
unknown; it is **downhill** otherwise (§9 uses this for the tier).

**The triggers** (all recorded with their evidence):

1. **No progress.** `stall.N` rounds in a row without progress (default 1, O-82).
      A disagreement between two roles is this case: a rejection of the same finding
   closes nothing.
2. **Too many rounds.** Review rounds past the target's own cycle cap, or
   attempts past `stall.attempts` (attempts that a question ended included), even
   with progress.
3. **A hang.** A running session writes no output and sends no hook event for
   `stall.T` (default 10 minutes, a maximum, O-82); the first session of a task
   gets no more. `layup run` kills it. A harness whose register row says it sends
   no hook events is judged by its output alone.
4. **A check that does not report.** A required check with no result for `ci.T`
   after the push; a `layup/verify` posted `pending` with a wait reason is a wait,
   not a missing result.
5. **The orchestrator.** A takeover of the lease (§2) opens a stall for the time
   with no run; it counts per project, not per task.

**Not stalls.** A step that fails cleanly (for example a context that does not
fit, §9) writes a "failed" row and a notice to the Operator on the task's issue;
the circuit breaker (below) writes a "milestone stopped" row. Neither counts in the
Stall Rate (`F-0003#73`).

**Waits are not stalls** (Fable-M5). The clock and the round count stop while a
task waits for a human: a bet, an escalation, an acceptance, a stall package.

**The procedure** (Decision Point 5, `F-0001#14`):

1. **The package.** Code builds the evidence from the records: the task, its
   plan, its findings and their history, its sets per round, the diffs as
   payloads, the gate outputs, the harness's exit and logs. It leaves out the
   sessions' own reasoning.
2. **The diagnosis.** An examiner session with a fresh context (an admitted
   harness other than one whose failure the package shows; never the session that
   stalled) reads the package and writes a diagnosis in a fixed form: the cause
   (disagreement, missing information, a wrong gate, a harness failure, a task
   too large, or other), each open unknown with its evidence, and the rung it
   recommends. If the examiner fails (no output for `stall.T`, an invalid form,
   or no admitted harness), `layup run` writes a row "diagnosis failed" and sends
   the package to the Operator at once; Stall Diagnosis counts that stall as one
   with no diagnosis. No stall goes on without one of the two rows (Fable-M6).
3. **The action** (smart-if P3, §10: `retry`, `panel` or `Operator`). The
   deterministic branch is the ladder: a **retry** with the diagnosis in the
   prompt, then a **panel**, then the **Operator**; a new trigger in the same task
   moves one rung up, with its own stall row, package and diagnosis. The panel rung is skipped when fewer than two admitted
   harnesses are free of the failure that the diagnosis names. A stall of the orchestrator (trigger 5) has only the diagnosis and the Operator.
A stall of an owner-role session (§8, a question) is a stall of that step: the
retry is a new owner session, and the asking role still waits; a panel's path or
an `external:` answer goes into the owner's next session.
4. **The panel** (vision 3.2; AgentJury's quorum). `panel.K` members, each a
   fresh session on the reasoning tier, from at least two harnesses, with the same
   sealed input (the package and the diagnosis), none seeing another's output, and
   a prompt that asks for hypotheses and trade-offs from first principles, not a
   hunt for fault. Each returns options, each with its evidence, and one
   recommendation. A synthesis session that sees only the members' outputs merges
   them into numbered options with one recommended path; unlike AgentJury, whose
   members vote on a fixed question, the options here are text, so a session, not
   code, merges them. Code checks the quorum: a majority of members returned a
   valid output, from at least two harnesses; otherwise the result is
   `insufficient panel`, and the next rung is the Operator. The recommended path
   goes into the task's next attempt. The Shape phase uses the same panel to recommend among surviving options (§8).
5. **The Operator** (an account whose role in `approvers.tsv` is Operator) gets
   one comment on the task's issue: the stall, the diagnosis, the panel's result,
   the evidence, and the answer form, one line:    `answer: <text>` (for a stalled question, it is the question's answer, by a
   human: the row is closed as answered by a human, so it is a human question of
   the delivery phase, §12, and the asking role gets it in its next attempt),
   `reroute <task> [<role>] to <harness>` (with a role, it moves that role's steps,
   for example the owner of a question), `stop <task>`, or `external: <text>` for
   a candidate solution from another source (vision 3.5); a change of a parameter goes to the control
   issue in the form of §10. `layup run` copies the reply (§3); it is planned
   input (Decision Point 5). A reroute writes a harness override for the task, and
   the next attempt starts from the base on that harness (§8). A stopped task's
   requirements go to the next bet; dropping one is the idea owner's decision.
6. **The outcome.** Each stall ends with one outcome row: closed without a human
   (by a retry or a panel), closed by the Operator, or the task stopped. Stall
   Diagnosis (`F-0003#61`) counts the stalls with no diagnosis row, which must be
   zero.

**The circuit breaker** (Shape Up; vision 3.5). At a milestone's cap, `layup run`
stops the milestone, whatever records arrive. Code computes each open task's
position (uphill or downhill); an examiner gives each uphill task's cause. When
every open task is downhill and one more round fits inside the band, P5 (§10) may grant **one** extension; otherwise the open work
goes to the next bet (for milestone 0, the Operator raises its cap or stops,
§12), and a
change of the band goes to the idea owner. No requirement is dropped by the
breaker: that is business-forking.

**A wrong gate** (Fable-M19). A stall whose diagnosis names a gate as its cause
opens an **early retrospective**, a planned point that Intake lists, whose rule
batch may fix the gate (§13).

**The dead-man job** (Sol-4, Author-10). A scheduled workflow in a control
repository that the Operator owns (never a target, FT5) acts as a separate
GitHub App, `layup-watch`, whose only permissions are metadata and contents
(read, for the repository activity) and issues (write); its key is that
repository's secret. A key of the LAYUP App cannot be limited in scope, so the job
never holds one. For each target on its list, it reads the forge's time of the last
update of the records branch from the repository activity, and when that time is
older than `3 × lease.H` it adds a notice ("no LAYUP run") to the target's
control issue. It writes nothing else on a target. The next run of `layup run` takes the lease
over, opens the orchestrator's stall, and builds its package. The forge can start
a scheduled workflow late, can drop a queued run, and disables the schedule of a
public repository with no activity for 60 days; the job therefore also commits
its last run time to its own repository, which keeps the schedule alive, and a
silence of the job is itself a known limit (L-F1).

## 12. Cost, budget and the measures

[ADR-0024](adr/0024-record-every-sessions-cost-and-stop-at-the-milestone-cap.md)
decides the cost and the budget. The records of the measures follow ADR-0014 and
the PSB's own definitions (§7.1, §7.2).

### The ledger

`telemetry.tsv` has one row per session, the unit of an "action" (vision 3.4;
a harness hides its inner calls, so a row per model call exists only where the
harness streams its usage): the session ID, the task and its requirement IDs (from
the task register), the role, the harness, the model, the billing type (an API
account or a subscription), the start, the first output (latency = first output −
start), the end (duration), the tokens by class with a status (`observed`,
`partial` or `unavailable`, with the reason), and the money with a status
(`reported` by the harness; `computed` from the tokens and a row of `prices.tsv`
that cites its source; or `unknown`). A session on a subscription has its money
`computed` from its tokens and the list price, marked as a subscription. An
unknown value is never zero (FT2, Sol-11). The start row of each session records
the spend cap and the wall-clock limit applied to it.
The smart-if's calls carry their own cost in `decisions.tsv` (§10); the report adds
both.

**Telemetry Completeness** (`F-0003#60`) counts a task as complete when each of its
sessions has `observed` tokens, a latency and a duration. The Operator's default
lets a harness with no token report do paid work under a wall-clock limit (O-80);
its tasks are then incomplete, and the check fails for them. The report says so;
for the pilot, the Operator can route every task to harnesses that report tokens
(known limit L-G1).

### The budget

**At Intake** (Decision Point 1, O-68), the idea owner answers the band, `B` (go
on without asking) and `U` (the upper edge, which is the money appetite of the
delivery), and the delivery's wall-clock appetite. `budget.tsv` is their one home
(§10). There is no per-task budget: tasks do not
exist yet, and most are found during the work (Shape Up; Author-7, Fable-M11).

**Milestone 0** is Intake and Shape, from Start to the first bet. Its clock starts
at Start, and its cap is the Start command's `--intake-cap`. At its cap, at any time before the first bet, the circuit breaker stops the
work, and the Operator raises the cap on the control issue or stops the target
(there is no earlier bet to send the work to). After the answers, its spend also
counts in the project total, and the bets' caps plus milestone 0's spend stay
within `U`.

**At each bet**, the brief proposes the milestone's cap in money and wall-clock;
code checks that the money caps of the bets so far are at most `U` and their
wall-clock caps at most the wall-clock appetite, and refuses a brief that does not
fit (the idea owner can raise `U` or the appetite, or the brief is reshaped); the
bet fixes the cap. A milestone's clock starts at its bet; it runs while at least one
of its open tasks has a running session or a step ready to run, and pauses only
while every open task waits for a human. The **project total** is the sum over all
milestones.

**Before each session start**, code sums the milestone's spend: the known money
(`reported` and `computed`), plus, for each running session, its spend cap when its
harness enforces one. A session on a harness with neither a token report nor a cap
has a wall-clock limit (`harness.<id>.wall`), and so has any session whose harness
enforces no spend cap, even one that reports tokens (O-80). Such a session counts
as an unknown part before it starts and while it runs; its money is known only
after it ends, if its harness reports tokens. The idea owner can accept a
bound for it (for example "up to USD 3 per session on H3"), recorded in
`budget.tsv`; with that bound, its part is known as that bound. Then:

| The total | The action |
| --------- | ---------- |
| the known total plus the new session's cap would pass the milestone's cap, or its wall-clock cap is reached | the circuit breaker (§11): stop; running sessions are killed |
| the project total is below `B`, and has no unknown part | go on |
| the project total is between `B` and `U`, and has no unknown part | smart-if P5 (§10); its deterministic branch stops and escalates to the idea owner |
| the project total, or the new session, has an unknown part that no accepted bound covers | escalate to the idea owner; an unknown amount never counts as below `B` |
| the project total reaches `U` | stop; the idea owner decides; no smart-if call |

### The records of the measures

| Measure (`F-0003`) | Record | The rule in code |
| ------------------ | ------ | ---------------- |
| Task Intervention Rate (#70) | `human-inputs.tsv`: every human action, from these sources: all issue and pull-request comments by a human (read at each step, copied whether a step acts on them or not) and their edits, reviews and review comments, reactions, the issue and pull-request events (closed, reopened, ready for review, labelled), workflow runs whose `triggering_actor` is a human, read per attempt (`…/runs/{run_id}/attempts/{n}`; a re-run is a restart), pushes and merges by a human from the repository activity, parameter changes, and each command given to `layup run` on the host; each with its time, actor, kind (answer, correction, restart, gate change, approval, review, parameter, stall answer, bet, acceptance, other), task, and its class; a reaction or a close is "other" and unplanned, so the rate reads high rather than low | planned only when it is the answer at a planned point listed at Intake (Intake, a
bet, an acceptance, a retrospective, a setup) or a push that such a point asks of
the Operator (the setup, a workflows batch), a stall answer, or an escalation answer confirmed as business-forking; every other input is unplanned (`F-0001#28`); a parameter change counts for each open task it reaches (§10). Tasks with an unplanned input over all tasks |
| Early Question Share (#75) | `questions.tsv`, one table for the project: each question with its phase ("before delivery" until the first task starts), its time, its asker, and whether a human was asked; each escalation to the idea owner is a row too | human questions before delivery over all human questions |
| Clarification Turnaround (#71) | the same rows: asked, answered, accepted, and the accepting actor (§8) | the 95th percentile of accepted minus asked, for questions resolved without a human (`F-0003#71`: "from the question to the accepted answer") |
| Reversal Rate (#72) | `audit.tsv` | overturned answers over audited answers |
| Missed Escalations (#57) | `audit.tsv`, over the rows of `screens.tsv` and the agents' decisions | confirmed misses in the sample; it must be zero; an empty or missing sample is "not measured", never a pass |
| First-Review Acceptance (#69) | `acceptance.tsv`: one row per review of a requirement, with the set of merged tasks it reviews; the review number is counted from the rows before it, on the original requirement ID even when a rejection's new need gets a new one | requirements accepted at review number 1 over delivered requirements |
| Delivery Lead Time (#68) | `telemetry.tsv` and `acceptance.tsv` | the median of acceptance time minus the start of the requirement's first task |
| Cost per Requirement (#74) | `telemetry.tsv` and `decisions.tsv` (a decision row names its task) | per requirement, the sum of the money of all tasks of that requirement (a shared task counts in full for each, as `F-0003#74` says); the median over requirements; `partial` when a part is unknown |
| Stall Rate and Resolution (#73) | the stall and outcome rows of §11 | tasks with a stall over all tasks; stalls closed without a human over all stalls |

**The audit** (`F-0003#57`, `#72`). At each end-of-milestone retrospective, and once more 30 days
after the last merge of the delivery (a planned point that Intake lists; the
Operator starts `layup run` for it, and the idea owner confirms its positives
there), code
draws a sample with a recorded seed and a hash of its population: the agents'
answers, screen rows and decisions whose task merged at least 30 days before and
that no earlier audit sampled (the window of the Reversal Rate, for both
measures). Its size is `audit.n` (a parameter set with evidence).
An auditor session on the reasoning tier, on a harness outside each item's authors,
judges each item: was the answer overturned; was a business-forking decision made
without an escalation. The idea owner confirms each positive at the retrospective (a planned point); a
confirmed reversal writes its time and the auditor into `questions.tsv`.

**Success criteria.** The PRD's success-criterion rows (§7: one per measure fact)
are the one source. A criterion that the idea owner adds in the Intake form is a
fact of the answers, classed `measure`, so it gets its row the same way.

**The baseline and the start values** (`F-0004#11`). For a pilot, the idea owner
posts the baseline, measured with the current process, in the Intake answer (§5),
and after it the start values in one batch comment on the control issue, one line
per measure: `start <measure> <value>`. `layup run` records `baseline.tsv` and
`start-values.tsv`, and renders the start values into the pilot's PRD as a
rendered task (§8).

**`layup report`** computes every measure from the records, with its population,
and compares it with its start value only when both exist; otherwise it prints "not
comparable", never a pass. A measure with an unknown input is `partial`.

## 13. Retrospective and learning

[ADR-0025](adr/0025-learn-routing-from-the-records-at-each-retrospective.md)
decides this section. The learning loop changes LAYUP's routing and defaults from
the records; it never changes a model's weights (PSB §6 Out of Scope, `REQ-015`).

**When.** After the Accept phase of each milestone, and early when a stall's
diagnosis names a gate (§11). Each retrospective is a planned point that Intake
lists, with two approvers: the idea owner confirms the audit's positives, and the
Operator (by default) approves the batch. An early retrospective carries only a
rule batch; `layup learn` runs only as `learn.trigger` says (a parameter, O-83:
each end-of-milestone retrospective by default, or every N milestones, or never),
and each run reads, per route, the records since that route's last adopted
update (a route's records carry over while it gets no change).

**The order.** First the audit of §12 and the idea owner's confirmations; then
`layup learn`; then the lessons; then one brief to the Operator.

**The reward** (vision 3.3, O-83). `layup learn`, a pure command (§2), computes a
reward per **implementing route** (the role, tier, harness and model of a task's
developer sessions); the terms of a task are charged to its implementing route.
Per task: plus when each requirement it served was accepted at its first review;
plus when its verification passed in the first round; minus its material
findings, its stalls, its unplanned human inputs, and each reversal of its answers
that an audit confirmed (charged to the route of the task that the answer served,
at the retrospective where it is confirmed; the final audit's are recorded and
charged to no route); minus its cost. The weights of the terms are parameters, whose evidence is the Operator's
comment (no fact supports a reward weight). A route's reward is the mean of its
tasks' values, compared with the mean over the routes of the same role and tier.
Cost is compared by unit: money against the median of the tasks with known money,
wall-clock against the median wall-clock. A task of unknown money is left out of
the money term, for its route and for the median; a route with no task of known
money gets no money term and no upward step, so an unknown cost never counts as a
low one (FT2; L-H2). The gates' first-attempt pass rate is not a term: the
PSB keeps it for monitoring only (`F-0003#58`, Author-12). The routes of the other
roles (plan, review, verification) get no update (known limit L-H2). A route with fewer than `learn.min` tasks since its last adopted update gets no
change.

**The update.** Code proposes, per route, a new weight: the old weight plus
`learn.step` times the route's reward minus the mean, kept inside
`[learn.min_weight, learn.max_weight]`. The first weight of a route is LAYUP's prior for it; with no prior, the mean of
the weighted routes of its role and tier, or the middle of the bounds when none
has a weight. So every admitted pair has a weight once one route of its role and
tier has one. The authority is a parameter (O-83):
`propose` (the default: the Operator adopts the proposal at the retrospective, and
can revert to the previous weights at a later one) or `apply` (applied within the
bounds and recorded). The weight ranks the admitted pairs in routing (§9), so learning acts at every
smart-if level; the fit point, or the table's order, decides only when no pair of
the role and tier has a weight, or on a tie (Fable-M17). **Exploration** (vision 3.3): a share `learn.explore` of the implementing
role's tasks (a parameter set with evidence) goes to the admitted pairs other than
the leader, in turn by the table's order, so that each route keeps evidence; with a share of zero,
the weights stop moving once one route leads (known limit L-H1).

**The lessons.** A retrospective session (the step table's row: Systems Architect,
reasoning tier, on a harness outside the milestone's authors where one exists, §9)
reads the milestone's stalls, diagnoses, findings, escalations and audit, and
writes each lesson with the records it rests on and its scope: `target` (a new
pitfall of the target's `docs/guardrails.md` §2, or a rule change) or `layup` (a
default, a routing prior, a catalog entry, a prompt of the step table).

**The batch.** Each proposed rule change (from a lesson, or a change that §6
refused in a task, marked with its source) becomes a rule batch task (§6, §8): its
sessions write, review and verify it before the brief. `layup run` builds one
brief for the Operator: the reward table and the routing proposal; each rule batch
by its head and hash, with its source; and the lessons. The Operator answers by one
comment: adopt the routing proposal or not, approve each batch or not, keep each
lesson or not (O-69). `layup run` copies the comment and writes the routing
register. A batch that changes a gate kind is run, before its merge, on its head
(which must pass) and on every recorded known-bad patch of each kind it touches,
and its own new one (each must fail), with the batch's own gate files (§6); a
patch that passes, or that no longer applies, refuses the merge, and the batch
goes back as a task.
`layup run` merges an approved batch; the Operator pushes and merges one that
changes `.github/workflows/` (O-93). A §2 pitfall is added lines, so it lands as a
normal task.

**The next project** (PSB §5, Sol-25). For each kept `layup` lesson, `layup run`
opens an issue on LAYUP's own repository with the lesson's text and links to its
records, where the LAYUP App is installed there; otherwise it records the lesson
for the Operator to file. The issue copies no content of the target, and links to
its records only when the target is public. A LAYUP task
decides it under LAYUP's gate, and a change lands in LAYUP's defaults: the routing
priors, a parameter default, a catalog entry or a prompt. Intake records the LAYUP version (§5); `layup run` stops at each start when its
own version differs from the recorded one; a bet brief may propose a new version
by one line, and the bet that adopts it records the new version. The host keeps
the binary of each version that a running target records, or that target's run
stops. A new target
starts from the defaults of its version. So a lesson reaches the next project through a
reviewed LAYUP release, never directly. Reading of O-69 ("learned … as RL"):
inside a target, lessons change rules only through a batch, and routing only
through the reward.

## 14. Coverage

Each row points to a walkthrough, or names the check or the known limit that
stands in its place; a row with none of the three fails
(`runs/T-hbw8/root-cause-missed-solution.md`, fix 5). Slices B to H added their
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
| #69 forge question; O-102 | the forge interface of §1; known limit L-A2 | 1 | 0013 |
| Table C, 3.4 (communication through issues) | W-12 step 6 | 3 | 0014 |
| S1 Problem Statement Quality (`F-0003#41`), `REQ-001` | [W-01](walkthroughs/W-01-problem-statement-quality.md) | 5 | 0014, 0015 |
| S2 Reproducible Discipline Setup (`#42`), `REQ-002`, `NFR-003`, `NFR-006`, Invariants 4 and 8 | [W-02](walkthroughs/W-02-reproducible-discipline-setup.md) | 5, 6 | 0016, 0017 |
| S3 Rule Protection (`#43`), `REQ-003`, Gate Integrity (`#64`), Invariant 3 | [W-03](walkthroughs/W-03-rule-protection.md) | 6 | 0017 |
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
| Table C 3.1 (model, context and solution routing) | W-12 step 1 and W-06 step 4 (routing); W-13 step 10 (weights); for solution routing, the named check of §8 Shape (options with a constraint table, filtered by code; the panel recommends among the survivors); for context, the named check of §9: `layup run` refuses a start whose estimated size passes the model's context (the estimate: L-D2) | 9, 10 | 0020, 0021 |
| S8 Human-on-the-Loop (`F-0003#48`), Missed Escalations (`#57`), `REQ-008`, Invariant 6 | [W-08](walkthroughs/W-08-human-on-the-loop.md) | 10 | 0021, 0022 |
| #69 B8 (an escalation after the work) | W-08 steps 2 to 5 | 10 | 0022 |
| O-78, O-79, O-84 (provider, authority, parameters, bounds) | W-08 steps 3, 4; W-06 step 3 | 10 | 0021 |
| Table C 2.1, the phased plan | W-05 (specification → milestone plan); the named check of §8: each `Must` requirement of the milestone has a task | 8 | 0019 |
| S9 Stall Resolution (`F-0003#49`), Stall Diagnosis (`#61`), `REQ-009`, `REQ-010`, Decision Point 5 | [W-09](walkthroughs/W-09-stall-resolution.md) | 11 | 0023 |
| #69 B2 (the time-limited stall), table C 3.2 (blind panel), 3.5 (circuit breaker, package, external answers) | W-09 | 11 | 0023 |
| S10 Cost Visibility (`F-0003#50`), Telemetry Completeness (`#60`), `REQ-011`, O-68, O-80 | [W-10](walkthroughs/W-10-cost-visibility.md) | 12 | 0024 |
| #69 B3 (requirement acceptance has no record) | W-11 steps 11, 12 | 8, 12 | 0019, 0024 |
| #69 B7 (no stop before the spend), table C 3.4 (telemetry per action) | W-10 steps 3 to 6 | 12 | 0024 |
| The PSB §7.2 measures and the audit (K53 to K57) | W-10 step 7; W-11 step 12; W-13 step 2 | 12 | 0024, 0014 |
| The pilot baseline and start values (K11, `F-0004#11`) | W-02 step 8 | 5, 12 | 0024 |
| Vision 2.3 and 3.3; #69 A1 (by O-69), table C 2.3, 3.3; O-69, O-83 | [W-13](walkthroughs/W-13-retrospective-and-learning.md) | 13 | 0025 |
| ADR-0012 part 6 (set up a target, run its gate from outside) | W-02; W-04 steps 5, 9 | 5, 6 | 0016 |

## 15. Known limits

Each limit is a finding that the design does not close, recorded here (O-66).
A "Close when" sentence names the event that closes the limit (O-113); an ID
such as ISO-01 is a scenario of the [research review](../runs/T-hbw8/slice-reviews/operator-codex-research.md), §8.

- **L-A1. The host is shared.** Role sessions run under the Operator's user, so
      a session that searches the host can reach the App's private key, which does
  not expire until the Operator revokes it, makes tokens for every target where
  the App is installed, and passes the records rulesets; and the Operator's own
  `gh` login; a comment made with that login passes as a human decision.
  Until sessions run in an isolated environment (a container or another
  operating-system user), the separation of O-77 is by convention for that case
  (K03, K04). The same is true for the target's own code: `layup gate` runs the
  product's tests and build scripts, and the gate files of a rule batch before its
  approval, in a scratch work tree on the host, under the same user, with the same
  network access; no section names the environment that a gate run gets.
  Close when: before LAYUP runs on a target that is not the Operator's; each
  session and each gate run gets a container or another user (ISO-01, ISO-02).
- **L-A2. One forge adapter.** The engine is forge-neutral (§1), but only the
  GitHub adapter is designed; another forge needs its adapter, and a forge that
  lacks one of the six capabilities cannot hold a target (#69 forge question,
  O-102).
  Close when: a target needs a forge other than GitHub.
- **L-A6. A newer baseline.** A target pins the baseline's latest state (O-101);
  LAYUP's setup steps and checks are written against the baseline's structure,
  so a change there can break a step. `layup setup verify` then fails on that
  step, and the setup stops until a LAYUP change follows the baseline.
- **L-A3. One host during delivery.** `layup run` runs in the foreground on one
  host; while the host is down, nothing moves (section 11 says how the stall is
  found). A takeover on another host needs the App's private key on that host. A
  run that paused between the records push that announces a forge write and the
  write itself still makes that one write after a takeover. A forge write whose
  result is not known (a time-out, or a stop after the write and before the next
  records push) has no receipt, and no rule reads the forge before the write is
  made again, so a retry can make it twice (for example a second comment), and a
  run that takes the write as done can leave it not made. A paid session start
  whose acknowledgement is lost has the same gap.
  Close when: before LAYUP runs on a target that is not the Operator's; each
  forge write gets an ID, and a retry reads the forge first (REC-01 to REC-04).
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
  reach; an accepted time is the end of the asker's next attempt. The measure is
  still recorded (§12); a question with no accepted time is counted as open,
  never left out.
  Close when: a session start takes less than 120 seconds; until then, record.
- **L-D2. Context by link, and its size.** A session gets the records that the step table's
  links reach; a relevant record that no link reaches is not in its prompt. The
  size check of §9 (bytes divided by four) is an estimate that is not measured
  against a tokenizer, and it counts only the prompt file, not the harness's own
  instructions, its tool definitions or the space for the output, so a start
  that passes the check can still not fit.
  Close when: the size check is built; it keeps a margin, and is compared with
  the token count that a harness reports (CTX-01).
- **L-D3. Acceptance of an answer.** The accepted time is mechanical: an attempt
  that asks the same question again in other words, without citing the answer
  ID, still counts the answer as accepted; the Clarification Turnaround can read
  too short.
- **L-D4. Independence by harness.** A check that must run on another harness
  than the author's is independent by harness only; two harnesses that run the
  same model, or read the same inputs, can make the same error, and no record
  measures how often a verifier misses what its author missed.
  Close when: the records show how often a verifier misses a finding that a
  later step finds; if often, a verifier needs another model family.
- **L-E1. An undeclared decision.** A session that makes a business-forking
  choice without declaring it, in prose that the four questions miss and in no
  file that the floor reads, is found only by the Missed Escalations audit (§12).
  Under `off` or `shadow` at P1 and P2, this is every undeclared choice made in
  prose: the checklist rows K32 and K35 are answered only at `cautious` or
  `delegate`.
  Close when: the Missed Escalations audit finds a missed decision; then P1 and
  P2 go to `cautious`.
- **L-E2. Calibrating the thresholds.** A threshold's evidence can be only the
  Operator's comment (§10), at every point. A calibration record is not required; at
  P1 it is one agreement share over the `screens.tsv` rows that the Missed
  Escalations sample checks, which are not held out and do not show the misses
  apart, and at P2 no record shows how often the provider misses "needs a human". Under `shadow`, the branch that the provider
  picks at P3 or P4 is not run, so its outcome is unknown; a `delegate` threshold
    for these two points rests only on the Operator's comment. P4 is asked only when
  no weight decides (§13), so it gets few rows; that is accepted.
  Close when: `screens.tsv` holds enough rows for a held-out set; each
  threshold then rests on a calibration record (DEC-01 to DEC-03).
- **L-F1. The dead-man job can be late or silent.** `layup run` checks at Start
  that the job's App and list work (its first notice), but not later. A scheduled workflow can
  start late or be dropped, so the notice of a dead host can come later than
  `3 × lease.H` or not at all; the package is built only at the next run. A
  push-based monitor outside the forge would not have this limit, and is not
  designed.
- **L-G1. Telemetry Completeness with O-80.** A task that a harness without a token
  report runs has incomplete telemetry, so Telemetry Completeness (`F-0003#60`)
  fails for it; only routing every task to harnesses that report tokens avoids it.
  Close when: routing is built; a target that requires the measure routes only
  to harnesses that report tokens.
- **L-G2. Human actions the forge does not show.** A human action on the LAYUP
  host outside `layup run` (for example killing a session) and an action on the
  forge with no API event are not recorded, so the Task Intervention Rate can read
  low.
- **L-G3. A spend cap between steps.** A harness that checks its spend cap between
  steps (as Claude Code's is) can pass it inside one step, so the known total plus
  the caps is not a hard ceiling.
- **L-H1. Learning and exploration.** With `learn.explore` at zero, every task
  of a role goes to the leading route, so the other routes get no new evidence and
  the weights stop moving; a route that became worse is not found by the reward.
  With exploration, the turns go by the table's order, not by chance, and
  inside a tier the reward does not take a task's difficulty into account, so a route that gets the
  hard tasks can read worse than a weaker route that gets the easy ones.
  Close when: the reward has data from real tasks; `learn.explore` stays above
  zero, and routes are compared only inside one tier (LEARN-01).
- **L-H2. Only some routes learn.** The routes of planning, review and
    verification keep their weights; their quality reaches the reward only through
  the implementing route's findings and acceptances. A route whose harness
  reports no tokens has no known money, so it never steps up.
