# W-04 — Stack-Dependent Gates

**Item.** S4, Stack-Dependent Gates (`F-0003#44`): "Gates for repository layout,
interface boundaries, contract checks, and test quality, which the domain-free
baseline does not give. The technology stack of a project selects these gates.
They add rules; they do not change the baseline rules (Invariant 7)." Also
Invariants 5 and 7, `REQ-004`, `NFR-004`, and finding A2 of #69.

**The case.** The Go target's architecture puts the storage code in
`internal/store`, which only `internal/app` may import. A developer session adds
an import of `internal/store` in `cmd/server`.

Sections are those of [`architecture.md`](../architecture.md).

| # | Actor | Input | Mechanism | Record | Tag | Where |
| - | ----- | ----- | --------- | ------ | --- | ----- |
| 1 | the setup | the Go entry of the stack catalog | writes `gofmt` and `go vet` steps; the layout, boundary and contract gates are `pending` in `docs/gates.tsv` | the target's CI job and manifest | `code` | §6; ADR-0016 |
| 2 | an architect session | the architecture approved at the first bet | writes a layout test and an import-boundary test as Go tests, from its package table ("`internal/store` is imported only by `internal/app`"); the batch lands at that planned point; the gates become active | a batch pull request; the manifest states `active` | `model` | §6; ADR-0016, ADR-0017 |
| 3 | the developer session | its task | adds the import; commits | commits in `work/` | `model` | §4 |
| 4 | `layup run` | the result | later: slice D | later: slice D | `code` | — |
| 5 | the target's CI | the pull request head | the boundary test fails: `cmd/server` imports `internal/store` | the failed required check on the forge | `code` | §6; ADR-0016 |
| 6 | `layup gate` | the base branch's gate files; the head | runs the base branch's boundary test on the head in a scratch work tree: `fail` | the status `layup/gates` = failure, and a result row | `code` | §6; ADR-0016 |
| 7 | `layup run` | the two results | later: slice D | later: slice D | `code` | — |

## Checklist rows

S4, R03, R10, I5, I7, K10, K14, K17, P02, P12; FT1, FT4. Known limits L-B1, L-B2.
