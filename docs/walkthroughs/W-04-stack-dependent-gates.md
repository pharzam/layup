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
| 1 | the setup | the Go entry of the stack catalog | open | open | `code` | open |
| 2 | an architect session | the approved architecture | open | open | `model` | open |
| 3 | the developer session | its task | adds the import; commits | commits in `work/` | `model` | §4 |
| 4 | `layup run` | the result | later: slice D | later: slice D | `code` | — |
| 5 | the target's CI | the pull request head | open | open | `code` | open |
| 6 | `layup gate` | the base branch's gate definitions; the head | open | open | `code` | open |
| 7 | `layup run` | the two results | later: slice D | later: slice D | `code` | — |
