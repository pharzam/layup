# 0014. Keep the records in the target with one writer

Date: 2026-09-29

## Status

Accepted

## Context

Invariant 1 says that no project state and no decision is kept only outside the
project repository (`F-0001#1`). The Operator decided that agents act under the
Operator's account through a GitHub App, so that the App's badge tells an agent
from a human, and that agent sessions get only the App's token (O-77, O-92); that
a person decides by an issue comment, which the orchestrator copies into Git
(O-73).
On 2026-09-29 the Operator moved the agents' writes to the App's installation
token, for the bot's logo and a clear split of agent from human (comment
5885656174 on #72), and chose the same for LAYUP itself (O-95).

The deep check found four defects in the first draft: six producers wrote records
that "only `layup run`" was said to write (Sol-3, Fable-M12); a commit author was
taken as proof of the actor (Sol-5, Fable-M3); a review of a pull request has no
`performed_via_github_app` field, so it cannot be told from a human's (Fable-M3,
Sol-30); and one App token gave agent sessions every right of the orchestrator
(Fable-M1). The evaluation found the actor taken from text that an agent writes
in five candidates (summary, "Failures that repeat", 3) and gave the pattern of a
typed result that one committer admits (AgentPlane, `evaluation/11`). The options
compared are in [`runs/T-hbw8/selection-v2.md`](../../runs/T-hbw8/selection-v2.md),
row 0014.

## Decision

We will keep LAYUP's records in the target, commit them with one writer, and take
each actor from a credential:

1. **The records branch.** The records live on the orphan branch
   `layup-records` of the target, as tab-separated tables with a header row and Markdown payloads (the
   form of ADR-0011 decision 2). The target's README, written at setup, names
   the branch; each pull request body links the records commit its task started
   from.
2. **One writer.** Only `layup run` commits to the records branch, from its own
   clone, with pushes that are never forced. An engine check hands it a table; a
   role session hands it a result file, and the attempt and base commit that
   bind it come from the session start row that `layup run` wrote; a result
   whose attempt is no longer the task's open attempt is refused; a human hands it
   an issue comment.
3. **Three kinds of actor** (O-77, O-95). LAYUP acts as the LAYUP App's bot,
   with an installation token made from the App's private key, which only
   `layup run` holds; the App is the only bypass actor of the records branch; it
   has the permissions of O-92 plus commit statuses, and no workflows and no
   administration permission. Humans are the accounts that `approvers.tsv` names by numeric user ID, each
   with a role (Operator, idea owner, approver); each rule names the role it takes. Role sessions hold no forge credential
   ([ADR-0015](0015-keep-model-calls-out-of-the-engine-checks.md)).
4. **A human decision** is an issue comment whose author ID is in
   `approvers.tsv` and whose `performed_via_github_app` field is empty. A review,
   a review comment, a commit, a reaction or an edit is never a decision, and is
   recorded as human input.
5. **Copy before read.** Before any step acts on a comment, `layup run` copies
   it into the records with its author ID, comment ID, App field, time and
   SHA-256.
6. **Fail closed.** At setup and at each start, `layup run` reads the effective
   rules of the default and the records branch and stops when a rule is missing;
   it also pushes an empty probe commit to a probe ref that the default
   branch's ruleset covers, with the App's token, and stops unless the forge
   refuses it for a rule violation (any other failure is "probe not run", never
   a pass). `layup audit` checks the actor of
   each update of both branches from the forge's repository activity, from the
   setup commits and the records' first commit (a `branch_creation`) onward,
   never from a commit author.

This amends ADR-0011 decision 2 (records on a branch of the target, not in the
tree of the default branch) and follows O-95, which changes the token of O-92:
every GitHub write of LAYUP is the App bot's, and a role session holds no GitHub
credential at all, which is stronger than "only the App's token".

We reject: the records on the default branch (a commit for each event puts every
open pull request out of date, and an agent's pull request could change a
record); a protected ref per task (a default clone does not fetch it); the commit
author as proof (FT3); a review of a pull request as a human decision (no App
field); the App's user access token (slice A as written after review round 1):
under it, pushes and merges show only the Operator's account, and the records
branch must let the admin role through (Q-15 on #72; O-95).

## Consequences

- An agent session cannot write a record or a decision with any token it holds.
- A clone of the target carries every record, and a human reads it with no tool.
- The Operator installs the App on each target; `layup run` keeps the App's
  private key on the LAYUP host, and the key does not expire until the Operator
  revokes it (known limit L-A1 of [`architecture.md`](../architecture.md)).
- The forge shows every write of LAYUP as the bot's, apart from the Operator's
  own; the App cannot read a ruleset's bypass list (known limit L-A4).
- The rulesets that the setup applies are a precondition of every run.
