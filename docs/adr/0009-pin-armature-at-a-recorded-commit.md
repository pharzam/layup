# 0009. Pin Armature at a recorded commit

Date: 2026-09-23

## Status

Accepted

## Context

LAYUP is a software product, so it obeys the discipline that it enforces on
other projects. Its discipline system is the Armature kit. The PSB
([`F-0001`](../facts/README.md)) sets Invariant 8: "Armature is used at a pinned,
recorded version." It also puts a change of the Armature rules out of scope
(PSB §6 Out of Scope).

The kit README says that a project is a one-time copy with no upstream link: no
`upstream` remote and no fork relationship. `npx degit` copies the files of one
commit and no Git history. So after the copy, the Armature commit is not in this
repository, and nothing records which version was copied unless the project
records it.

Measured on 2026-09-23: `main` of `pharzam/armature` was
`a95965534b14b0bf14ad74da0c9a45b5f4aedf88`. The copy was taken with
`npx degit pharzam/armature#a95965534b14b0bf14ad74da0c9a45b5f4aedf88`, and it
became the root commit `d2516fd` of this repository, with no change. The tree of
`d2516fd` is `8ffb250afd584da8b418bc220fb6d72e802924ce`; the tree of the Armature
commit is the same (GitHub API, checked by the plan review of issue #1).

## Decision

We will record the Armature version in `docs/setup/armature.pin`: the source, the
commit, the tree, the copy method, and the date. The check `pin` in
`docs/setup/setup-check.sh` fails when the file is absent, when the commit is not
40 hexadecimal characters, when the clone is shallow, or when the tree differs
from the tree of the one root commit.

We will keep the root commit as the unmodified copy. Adaptation lands in later
commits, so the diff from the root commit is the whole adaptation.

We reject: a Git submodule or subtree, and a fork (each is an upstream link the kit
forbids); a synced copy (the operator forbids it, and it would change the baseline
rules without a decision); a tag or a branch name as the pin (a name can move, a
SHA cannot — see `SHA` in the [glossary](../glossary.md)).

## Consequences

- A reader can check the pin with no network: `git rev-parse <root>^{tree}`.
- Armature changes after `a959655` come into LAYUP only by hand, each as its own
  task that updates the pin in the same change. The tree check then needs a new
  rule, because the root commit stays the old copy; that task writes it.
- The check needs the full history. A shallow clone fails with a named cause, so
  CI must fetch the full history for this check.
