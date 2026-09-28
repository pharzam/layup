# 0016. Keep the run records in Git and tell agent from human

Date: 2026-09-28

## Status

Proposed

## Context

Git is the system of record (Invariant 1, `F-0001#1`). In the previous
architecture the plan, the plan review and the review record stayed as issue
comments with no path into Git (finding B6, issue #68); requirement acceptance
had no record (finding B3, Decision Point 3, `F-0001#12`); and telemetry was
written after the work, with no requirement and no price (finding B7). Every
human decision is an issue comment (O-73). The agents and the Operator act with
one GitHub account today (O-9), and the Operator does not want a second account
(O-77): agents act under the Operator's account through a GitHub App, so that
GitHub shows the Operator's avatar with the App's badge and the API field
`performed_via_github_app` names the App (GitHub documentation, "Authenticating
with a GitHub App on behalf of a user"). ADR-0011 decision 2 rejected state on
"an orphan branch (outside the tree a reviewer reads)". The options and the
selection are in [`runs/T-hbw8/selection.md`](../../runs/T-hbw8/selection.md),
section 6; the decisions in
[`runs/T-hbw8/operator-decisions.md`](../../runs/T-hbw8/operator-decisions.md).

## Decision

We will keep the records of a run in the target's Git, as follows:

1. **The records branch.** The target has one branch, `layup-records`, that
   holds only the run records: the registers (`parameters.tsv`,
   `approvers.tsv`, `harnesses.tsv`, `roles.tsv`, `routing.tsv`, `prices.tsv`,
   `budget.tsv`, `stack.tsv`), the event tables (`phase.tsv`,
   `heartbeat.tsv`, `acceptance.tsv`) and one directory per task,
   `runs/<task>/`, with `handoffs.tsv`, `questions.tsv`, `decisions.tsv`,
   `telemetry.tsv` and the typed Markdown payloads they point at. Event tables
   only add rows; a register is edited in place and Git history is its log
   ([ADR-0011](0011-structure-the-core-engine-as-a-go-cli-over-repository-files.md),
   decision 2). The branch is fetched by its name; a clean clone of any harness
   reads it with `git fetch origin layup-records` and no LAYUP tool.
2. **One writer.** Only `layup run` writes the records branch, with the
   installation token of the LAYUP GitHub App (its author is the App's bot). A
   ruleset on the branch blocks every other update and lists only that App as a
   bypass actor. Setup proves the ruleset with a probe: a push with the agents'
   token must fail. If the forge cannot give this, the run records say so, and
   a deterministic audit of the branch history (every commit by the App's bot)
   is the complement.
3. **Three identities, one account.** A human decision is a comment or review by
   the Operator's or the idea owner's account (listed in `approvers.tsv`,
   written at Intake) with `performed_via_github_app` empty. An agent acts under
   the Operator's account through the App's user access token, so its comments
   carry the App's badge. The orchestrator acts as the App's bot. `layup run`
   starts each harness session with the App's user access token only, never with
   the Operator's own token or SSH key (a separate operating-system user, a
   container or a sandbox, O-77); a run that cannot do this records that the
   separation is by convention only.
4. **The forge copy.** Before any action reads a comment or a review, `layup run`
   copies it into the records branch: the body word for word, the author, the
   comment ID, the `performed_via_github_app` value, the time and a SHA-256 of
   the body. A later edit or deletion is a new correction event and pauses the
   work that depends on it. A decision that is only on the forge authorizes
   nothing (O-73).
5. **The record kinds.** `handoffs.tsv` (from role, to role, kind, artifact path
   and hash, condition, validity); `questions.tsv` (kind, how the kind was
   found, owner role, asked, answered, accepted, human yes or no, planned yes or
   no); `acceptance.tsv` (requirement, milestone, delivered commit, the idea
   owner's accept or reject, the comment ID, the review count, the time), which
   closes a requirement only on an accept (Decision Point 3); `telemetry.tsv`
   per action (task, requirement or requirements with a cost share, role,
   harness, model, input and output tokens, latency, wall-clock time, price and
   its source in `prices.tsv`). A value that a harness does not give is
   `not reported`, which counts as incomplete, never as zero (`F-0003#60`).
6. **Merges.** Only the orchestrator merges a change into the target's default
   branch, after its gates and its verification
   ([ADR-0015](0015-route-role-sessions-over-registered-harnesses.md)). The
   audit reports a merge by any other actor as unplanned input.

This decision amends ADR-0011 decision 2 on one point: the run records live on
the records branch, not on the default branch. The reason of the rejection in
ADR-0011 ("outside the tree a reviewer reads") is answered by fetching the
branch by name, and by citing the records commit in each check result and each
handoff; the new reason for a separate branch is that no product pull request
and no agent can change a record.

We reject: a protected ref per task (a default clone misses the task refs); the
records on the default branch through agent commits (an agent can write a row);
a human decision as a commit in a pull request (human load, and the PSB counts
a pull-request approval as unplanned input); and a text tag in a comment as the
agent marker (anyone can type it).

## Consequences

- Each decision, acceptance and cost of a run is in Git with its source, so a
  second harness on a clean clone continues without the first (Invariants 1 and
  9), and the measures of PSB §7.2 can be computed from the records.
- The separation of agent and human is as strong as the isolation of the agent
  sessions from the Operator's own credentials; this record states that limit
  and does not hide it.
- The Operator creates the LAYUP GitHub App under the Operator's account and
  installs it on each target: one setup step with a human decision.
- On a repository owned by one user, GitHub cannot stop an actor with the
  Operator's rights from merging; the orchestrator's merge rule and the audit
  detect a merge by another actor, they do not prevent it.
