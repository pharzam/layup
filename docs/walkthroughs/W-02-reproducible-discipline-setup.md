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
| 1 | `layup run`, the Operator | LAYUP's `docs/setup/armature.pin`; the pinned commit | `git clone` of the pinned commit, `.git` removed (step S02 as changed); the Operator pushes the unmodified copy as the root commit (step S03), before the Intake of W-01 | the root commit; its tree equals the pinned tree | `code`, `human` | §5 |
| 2 | `layup run` | the copy; `docs/setup/steps.tsv`; the Intake answers | on a branch from the root commit, does each row that needs no judgement; takes each value from an Intake answer by question ID, a catalog entry, or a fact source that the Operator accepted in the Intake answer; a value with no source keeps its marker as an open gap | the target's setup record: each value with its source | `code` | §5 |
| 3 | a role session | the prose rows S07, S08, S09, S14 | writes them from the facts; code runs each row's check | the session result; the rows' check output | `model` | §5; ADR-0015 |
| 4 | `layup run` | the Intake answer "Go"; the Go entry of the stack catalog | writes `go.mod`, the tools' configuration, `docs/gates.tsv` and one CI job per gate kind; layout, boundary and contract are `pending` | files in the setup | `code` | §6; ADR-0016 |
| 5 | `layup setup verify` | a checkout of the setup | the pinned baseline's discipline tests; the checks of LAYUP's `setup-check.sh` from outside, for this target (`pin`: the root tree equals the pinned tree); zero values without a source; each source resolves; each active kind (`gofmt`, `go vet`) passes on the clean tree (`clear`: no package yet) and fails on its known-bad fixture; the pending kinds' fixtures are not run | a result table, committed by `layup run` | `code` | §5, §6; ADR-0016 |
| 6 | the Operator | the verified setup; the ruleset file | pushes the setup commits on top of the root commit and applies the two rulesets, with the Operator's own login and the commands that `layup run` prints (the planned point "setup") | the commits and the rulesets on the forge | `human` | §5; ADR-0017 |
| 7 | `layup run` | the default branch; the forge's effective rules | checks that the pushed tree equals the verified one; reads the rules of both branches back; pushes to `layup-probe` and expects a refusal for a rule violation (any other failure: "probe not run"); reads the repository activity with the App's token; stops when a rule is missing or a probe gives the wrong result | a check row with the rules read back | `code` | §3, §5; ADR-0014 |
| 8 | the idea owner; `layup run` | for a pilot target: the baseline, measured with the current process, in the Intake answer; later the start values, one line per measure, on the control issue | records `baseline.tsv` and `start-values.tsv`; renders the start values into the pilot's PRD as a rendered task (`F-0004#11`) | the two registers; the rendered task | `human`, `code` | §5, §12 |

## Checklist rows

S2, R02, R09, R12, I4, I8, K08, K11, K15, K16, D07, D18; FT1, FT5.
