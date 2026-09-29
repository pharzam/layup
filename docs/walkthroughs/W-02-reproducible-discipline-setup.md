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
| 1 | `layup run` | LAYUP's `docs/setup/armature.pin`; the pinned commit | `git clone` of the pinned commit into LAYUP's work area; removes `.git` (step S02 as changed) | the pinned commit and tree in the setup record | `code` | §5 |
| 2 | `layup run` | the copy; `docs/setup/steps.tsv` | does every row that needs no judgement; takes each value from an Intake answer (by comment ID), a fact line or a catalog entry; a value with no source keeps its marker as an open gap | the target's setup record: each value with its source | `code` | §5 |
| 3 | a role session | the prose rows S07, S08, S09, S14 | writes them from the facts; code runs each row's check | the session result; the rows' check output | `model` | §5; ADR-0015 |
| 4 | `layup run` | the Intake answer "Go"; the Go entry of the stack catalog | writes the tools' configuration, `docs/gates.tsv` and the CI job with one step per gate kind (§6) | files in the setup | `code` | §6; ADR-0016 |
| 5 | `layup setup verify` | a checkout of the setup | the pinned baseline's discipline tests; zero values without a source; each source resolves; each known-bad fixture makes its gate fail | a result table, committed by `layup run` | `code` | §5, §6; ADR-0016 |
| 6 | the Operator | the verified setup; the ruleset file | pushes the setup as the first commit of the default branch and applies the two rulesets, with the Operator's own login and the two commands that `layup run` prints (the planned point "setup") | the commit and the rulesets on the forge | `human` | §5; ADR-0017 |
| 7 | `layup run` | the default branch | checks that the pushed tree equals the verified one | a check row | `code` | §5 |
| 8 | `layup run` | the forge's effective rules of both branches | compares them with the expected rulesets; pushes an empty probe commit to `layup-probe` with the App's token and expects a refusal; opens a probe pull request, posts the `layup/` statuses and reads back that they count for the pinned checks; stops when a rule is missing or a probe gives the wrong result | a check row with the rules read back | `code` | §3, §5; ADR-0014 |

## Checklist rows

S2, R02, R09, R12, I4, I8, K08, K11, K15, K16, D07, D18; FT1, FT5.
