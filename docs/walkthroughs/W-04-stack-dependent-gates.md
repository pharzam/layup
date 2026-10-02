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
| 1 | the setup | the Go entry of the stack catalog | writes one CI job per gate kind; `gofmt` and `go vet` (the static-checks kind) and `go test` (the test kind) are active; layout, boundary and contract are `pending` in `docs/gates.tsv` | the target's CI jobs and manifest | `code` | §6; ADR-0016 |
| 2 | an architect session | the architecture approved at the first bet | writes the `depguard` rule "`internal/store` is imported only by `internal/app`", a layout test, the manifest with those kinds `active`, and one known-bad patch per kind it activates | the session result | `model` | §6; ADR-0016 |
| 3 | `layup run` | the result | the batch is a task (§8): its issue, plan and plan review, then pushes it to `batch/bet-1`, opens its pull request, and runs its verification before the bet | the batch's task records and pull request | `code` | §6, §8; ADR-0017, ADR-0019 |
| 4 | the approver | the bet brief, which names the verified batch head and its rule hash | bets by one comment, which approves the batch (the planned point of the first bet) | the comment, copied (§3) | `human` | §6, §8 |
| 5 | `layup run`, `layup gate` | the approved batch | refuses the approval if the batch head moved since the request; records the rule-file hash with the approval and the known-bad patches on the records branch; runs the batch's gate files on the head (must pass) and on the head with each patch (each must fail); `layup/rules` passes; merges the batch | the hash, the detection results, the merge | `code` | §6; ADR-0016, ADR-0017 |
| 6 | the developer session | its task | adds the import; commits | commits in `repo/` | `model` | §4 |
| 7 | `layup run` | the result | the handoff, rule-path and workflow checks; pushes the branch; opens a draft pull request | the draft pull request | `code` | §8; ADR-0019 |
| 8 | the target's CI | the pull request head | the boundary job (`depguard`) fails: `cmd/server` imports `internal/store` | the failed required check on the forge | `code` | §6; ADR-0016 |
| 9 | `layup gate` | the base branch's gate files; the head | runs the base branch's boundary rule on the head in a scratch work tree: `fail` | the status `layup/gates` = failure, and a result row | `code` | §6; ADR-0016 |
| 10 | `layup run` | the two results | gives the developer a new attempt with the findings; the pull request stays a draft; no verifier starts and no review is requested (W-07 step 5) | an attempt row | `code` | §8; ADR-0019 |

## Checklist rows

S4, R03, R10, I5, I7, K10, K14, K17, P02, P12; FT1, FT4. Known limits L-B1, L-B2.
