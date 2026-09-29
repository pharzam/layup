# 0016. Put the native stack gates in the target

Date: 2026-09-29

## Status

Proposed

## Context

The PSB asks for gates for "repository layout, interface boundaries, contract
checks, and test quality", selected by the stack, that add rules and change no
baseline rule (`F-0003#44`, Invariant 7), and a target that "passes its gates
without the automation" (`F-0003#65`, Invariant 2). The Operator decided that the
stack gates of a target are the target's own native tools and its own CI job,
written at setup, which run with LAYUP absent; `layup gate` runs the same native
commands from outside; this supersedes O-10 and O-11 where they put the gates
outside the target (O-76).
[ADR-0011](0011-structure-the-core-engine-as-a-go-cli-over-repository-files.md)
decision 4 still put them outside (Sol-26). The review of #69 asked who writes
the gates of a new stack (B9).

The evaluation found gates that pass when they did not run (Gas Town, BMAD,
isitdone) and gates that run the pull request's own copy of the gate scripts
(sdlc-gh, AI-SDLC); it gave the pattern of a catalog of native commands per stack
and of a verifier that reads its rules from the base branch
([`runs/T-hbw8/evaluation/summary.md`](../../runs/T-hbw8/evaluation/summary.md)).
The options compared are in
[`runs/T-hbw8/selection-v2.md`](../../runs/T-hbw8/selection-v2.md), row 0016.

## Decision

We will put each target's stack gates in the target, as its own native tools,
from a stack catalog that LAYUP keeps:

1. **The catalog.** LAYUP's repository has one entry per stack. Per gate kind
      (layout, interface boundary, contract, test quality, and any further kind of
   the stack, such as Go's static checks), an entry names the tool
   and its version, the command, the paths in scope, the configuration it writes,
   a known-bad fixture that must make the gate fail, and the tool's
   documentation as evidence. A new entry is a LAYUP change under LAYUP's gate,
   never the work of a session on a target.
2. **In the target.** The setup writes the tools' configuration, the gate
   manifest `docs/gates.tsv`, and one CI job per gate kind, each its own
   required check, pinned to GitHub Actions as its source. A gate kind whose
   rules depend on the architecture is `pending` until its activation at the
   first bet; while it is pending, its job fails a pull request that changes a
   path in the product's scope and passes one that changes none, with that
   reason. The jobs of all kinds exist from the setup.
3. **From outside.** `layup gate` checks out the base branch's manifest and gate
   files, applies them to the head of a pull request, and reports `pass`, `fail`
   or `not-active` per kind. It never runs the head's gate files for a status,
   except those of an approved rule batch; before its approval, a batch's own
   gate files run only in a scratch tree, as evidence. A kind whose tree or change has no path in
   its scope is `clear`, with that reason. The status `layup/gates` is a success
   only when every kind passed or is `clear`; `not-active` is never a pass.
4. **Detection.** The setup verification runs each active kind on the clean tree
   (it must pass) and on its known-bad fixture (it must fail); the activation
   batch, and every later rule batch that changes a gate kind, carries one new
   known-bad patch per kind it touches, kept on the records branch; before it
   merges, every recorded patch of each touched kind and its new one must fail,
   and a patch that no longer applies refuses the merge. A pending kind's fixture is not run, and never
   counts as a detection.
5. **No LAYUP check in the target.** LAYUP's `setup-check.sh` and its job stay
   out of the target; `layup setup verify` does those checks from outside.

This amends ADR-0011 decision 4 (the stack gates are the target's own and run
inside it; `layup gate` runs the same commands from outside) and decision 7
(setup step S12 no longer adds LAYUP's `setup-check` job to a target; step S02
copies with `git clone`; the Operator pushes the unmodified copy as the root
commit, as step S03 asks, because the App has no workflows permission). ADR-0012 part 6 stays true: `layup` sets up the target
and runs its gate from outside.

We reject: the gates only outside the target (finding A2 of #69; O-76); the pull
request's own gate files as the only gate (FT4); a gate that passes when it finds
nothing to check without saying so (FT1).

## Consequences

- A target passes its own gates with LAYUP absent.
- The LAYUP host needs the toolchain of each stack it runs `layup gate` for.
- A second stack needs a catalog entry with a fixture for each gate kind before
  the pilot can show Generality (`F-0003#67`).
- The coverage floor of a test-quality gate is an open gap of the target until
  the idea owner or the Operator sets it with evidence.
