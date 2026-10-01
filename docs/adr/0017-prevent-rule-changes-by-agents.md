# 0017. Prevent rule changes by agents

Date: 2026-09-29

## Status

Accepted

## Context

Invariant 3 says: "The agents that do the work cannot change the rules or the
gates that check the work" (`F-0001#3`); Gate Integrity asks for "100% detection
of known-bad commits. 0 agent writes to rule paths." (`F-0003#64`). The Operator
decided that rule-path changes happen at the retrospective, approved there (O-69),
that no second GitHub account is made (O-77), and that the App has no workflows
or administration permission for the agents (O-92).

The deep check found that the first draft detected a rule change by an audit
after the merge, while a ruleset that could prevent it was left unused (Sol-1,
Fable-M2), and that a lesson that the baseline's step 7 writes into
`docs/guardrails.md` in the same pull request would be refused (Fable-N9). The
evaluation gave the patterns of a ruleset applied and read back, and of a hash of
the rule files recorded at approval (sdlc-gh, AgentPlane). The options compared
are in [`runs/T-hbw8/selection-v2.md`](../../runs/T-hbw8/selection-v2.md), row
0017.

## Decision

We will prevent rule changes by agents with three layers, and keep an audit as
the complement:

1. **No credential.** Role sessions hold no forge credential
   ([ADR-0015](0015-keep-model-calls-out-of-the-engine-checks.md)).
2. **No rule change leaves the host in a task.** The rule paths of a target are a
   register on its records branch. Before `layup run` pushes a task branch, it
   refuses a change to a rule path; the change becomes a proposal for the next
   batch. Lines added inside
   section 2 of `docs/guardrails.md` are not a rule change.
3. **The forge refuses.** The default branch and the probe ref `layup-probe` have one ruleset with no bypass actor:
   a pull request is required; the required checks are the native gate kinds and
   `layup/gates`, `layup/spec`, `layup/verify` and `layup/rules`, each `layup/` check pinned to
   the LAYUP App as its source; force pushes and deletion are blocked. The records branch
   restricts updates and deletion to the LAYUP App, its only bypass actor
   (O-95). The Operator applies both at setup,
   and `layup run` reads them back and probes them before it goes on
   ([ADR-0014](0014-keep-the-records-in-the-target-with-one-writer.md)).
4. **A rule change lands** only in a batch, proposed at a planned point (the
   setup, the gate activation of the first bet, each retrospective) on a branch
   `batch/<point>` with its own pull request, so the approver sees it before
   deciding. The approval request names the batch head and the hash of the rule
   files that the batch changes; code refuses an approval when the head has
   moved since, and records the hash with the approval; `layup/rules` passes a
   pull request that changes a rule path only when that hash still matches, or
   when it is a record that code renders from an inventory version that a bet or
   an escalation decision wrote. `layup run` merges the approved batch; for one that changes
   `.github/workflows/`, the Operator pushes the branch and the approver merges it,
   because the App has no workflows permission (O-92, O-93).
5. **The audit.** `layup audit` lists each rule-path change on the default branch
   with its merge actor, from the forge's repository activity, and its approval.
   `F-0003#64` counts a rule-path change that landed outside an approved batch,
   and not as a rendered record of an approved bet or escalation decision, as an
   agent write. The merge actor of each change must be the App's bot, or, for a
   batch that changes `.github/workflows/`, the account in `approvers.tsv` whose
   role is that batch's approver (O-95, O-93).

We reject: code owners whose approval comes from a second human account (O-77);
detection by audit alone (Sol-1); a rule-path change at any time other than an
approved batch (O-69).

## Consequences

- An agent cannot merge a change to a gate with any credential it holds.
- The gate activation at the first bet is a rule change at a planned point, the
  second one after the setup; this reads O-69 ("only at the retrospective") with
  O-76 (the gates are written at setup) for the gates that need the
  architecture. The approval brief asks the Operator to confirm this reading.
- When LAYUP is absent, the Operator removes the `layup/` checks from the
  ruleset; the target's setup record says how.
- The Operator's own account and a process that reaches the App's key on the host
  can still change a rule (known limit L-A1 of
  [`architecture.md`](../architecture.md)).
