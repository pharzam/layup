# W-02 — Reproducible Discipline Setup

**Item.** S2, Reproducible Discipline Setup (`F-0003#42`): "A repeatable,
evidence-based setup of the … baseline for a new project, for the domain
and technology stack of that project." Also Invariants 4 and 8, `REQ-002`,
`NFR-003`, `NFR-006`, and ADR-0012 part 6.

**The case.** The Intake of W-01 is answered: the stack is Go, the repository is
public. `layup run` sets up the target.

Sections are those of [`architecture.md`](../architecture.md).

| # | Actor | Input | Mechanism | Record | Tag | Where |
| - | ----- | ----- | --------- | ------ | --- | ----- |
| 1 | `layup run` | `docs/setup/armature.pin` of LAYUP; the pinned commit | open | open | `code` | open |
| 2 | `layup run` | the copy; `docs/setup/steps.tsv` | open | open | `code` | open |
| 3 | a role session | the prose rows of the steps | open | open | `model` | open |
| 4 | `layup run` | the Intake answers; the stack catalog entry for Go | open | open | `code` | open |
| 5 | `layup setup verify` | a checkout of the setup | open | open | `code` | open |
| 6 | `layup run` | the verified setup | open | open | `code` | open |
| 7 | the Operator | the ruleset request | open | open | `human` | open |
| 8 | `layup run` | the forge's effective rules | open | open | `code` | open |
