# 0013. Run stack gates through a LAYUP GitHub App that posts required check runs

Date: 2026-09-25

## Status

Accepted

## Context

[ADR-0011](0011-structure-the-core-engine-as-a-go-cli-over-repository-files.md) decides that `layup gate` runs a target's stack-dependent gates from outside the target (decision 4) and writes no LAYUP code into it (decision 3), and its Consequences leave open which runner reports on a target's pull requests. `PRD-0001` requires that no PR with a failed gate reaches human review or merges (REQ-007, `F-0003#58`), that a gate that did not run counts as not passed (NFR-004, `F-0001#5`), and that the target passes its gates without LAYUP (NFR-002, `F-0003#65`). Question 1 of `PRD-0001` §11.

A panel of three (ADR-0006; the records under `runs/T-7qvc/`) compared four options: a GitHub App that posts required check runs, started by the pull-request event (members A, B, C); the same App with a polling runner instead of events (A); a workflow in LAYUP's repository that posts a commit status through a token (B, C); a local pre-merge step with a `layup merge` command (B, C). The falsifiable tests the panel named: stop the runner and try to merge (the forge must block); post a status with the agents' token (the forge must refuse); merge through the forge UI after a local `fail` (the merge must not go through). The Operator decided on 2026-09-25 (O-52 on #66): a GitHub App with check runs; and (O-54) the verdict lives in the target's `runs/`, with the forge's protection setting mirrored and verified in the target's Git.

## Decision

We will run a target's stack gates through **a LAYUP GitHub App**, the runner, which is a component outside the engine:

1. **The App is installed on each target.** On each `pull_request` event it starts a job on a LAYUP-owned runner (a small receiver, or a workflow in LAYUP's own repository that the App's events reach). The job checks out the pull request's head with a read token, runs `layup gate <checkout>` at a pinned LAYUP commit, and posts **one check run per gate kind** under the App's identity: `layup/layout`, `layup/interfaces`, `layup/contracts`, `layup/tests`, and the guard of [ADR-0014](0014-protect-rule-paths-with-a-rule-guard-check.md), `layup/rule-guard`. A `pass` posts `success`; a `fail` and a `not-active` post `failure` (NFR-004 allows no neutral result).
2. **The target's protection requires each check by name and pins it to the LAYUP App's identity** (the App's application ID), in the branch-protection body that the kit's own setup step S13 writes and `layup setup verify` compares with the live setting (a kit record of the target, O-54). S13 pins the kit's own checks to the GitHub Actions app; the `layup/*` checks are pinned to the LAYUP App instead, and the task that builds `layup setup` records that departure in `steps.tsv`. A check run posted by any other identity, and a commit status of the same name, does not satisfy the requirement; a required check that reports nothing keeps the pull request at "expected", so a gate that did not run blocks the merge by the forge's rule, not by a rule of ours.
3. **The verdict is a row in the target's Git before it is a forge status.** The job checks the pull request's head `X`, appends one row per gate to `runs/<task>/gates.tsv` (`ts, pr, head_sha, gate, verdict, layup_commit`, `head_sha` = `X`) as a commit `Y` on the pull request's branch under the App identity, asserts that `Y` differs from `X` only under `runs/*/gates.tsv`, and posts the check runs **on `Y`**, the new head, with a summary that names `X` and the row. A commit by the App that touches only `runs/*/gates.tsv` does not start a new run; any other commit does, and the cycle repeats on the new head. `<task>` is the task ID in the branch name (the kit's `T-` scheme); a pull request from a fork, whose branch the App cannot write, gets `failure` with that reason and is out of the pilot's scope.
4. **The engine stays a CLI** (ADR-0011 decision 1): the App and the receiver hold no gate logic; they call the binary at a pinned commit and relay its verdicts. A target's gate set comes from its recorded stack, `docs/setup/stack.tsv` (`gate, kind, command, active`), which `layup setup` writes from the answer to setup step S01, never from the pull request.

We rejected: **a commit status posted by a workflow with a token** — any identity with write access can post a status, so the shared account (O-9) and, after O-53, any agent with `repo:status` can forge a `success`; **a local pre-merge step and `layup merge`** — nothing stops a merge through the forge; **polling instead of events** — a delay with no gain once the App exists (the App may fall back to a poll when an event is lost).

## Consequences

- REQ-007 and NFR-004 are met by construction of the forge's required checks; the falsifiers are the tests of `PRD-0001` §7.1 for REQ-004 and REQ-007.
- **A network component exists outside the engine.** ADR-0011 decision 1 keeps the engine free of a daemon and a network service; the receiver is the runner, not the engine. ADR-0011's Status becomes `Accepted. Amended by ADR-0013` (the runner its Consequences left open).
- **Independence (NFR-002) has a stated limit:** a target whose protection requires the App's checks cannot merge while the App is removed until its Operator removes those names from the protection body — a kit record, editable without LAYUP. `layup setup verify` reports the difference.
- **Key custody:** the App's private key is held by LAYUP's Operator, never by an agent (O-53); no agent session ever acts as the LAYUP App ([ADR-0014](0014-protect-rule-paths-with-a-rule-guard-check.md) decision 1); a leaked key can post `success` on any gate, the main risk the panel named. The setup record of a target names who holds it (Invariant 4).
- The App's identity is also the identity of the rule guard (ADR-0014) and of the checks that the escalation and stall states raise ([ADR-0016](0016-escalate-by-a-deterministic-floor-and-a-declared-class.md), [ADR-0017](0017-stop-a-stall-at-a-counted-limit-and-examine-it-fresh.md)).
- The pilot measures the delay from a pull-request update to the check result, and the Operator decisions that runner outages cause (the cost ADR-0012 names).
