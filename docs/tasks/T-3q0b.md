# T-3q0b — row 62b of the plan of M2e, the checks of the bet and the rendering

Issue: [#350](https://github.com/pharzam/layup/issues/350), row 62b of [the tasks of
M2e](../plan/README.md#the-tasks-of-m2e), opened by `T-57vx` (#346). Serves
`F-0003#51`, `F-0003#54` (`REQ-012`). Branched off `2aafa64`; `main` was merged
three times: `f12da6a` (`2e1e9b4`), `8e26eca` (`db3330f`), `0e2f06c` (`c9a31c7`).
Pull request [#382](https://github.com/pharzam/layup/pull/382), ready. Author of the plan and
of record: Claude Sonnet 5.5; implementer: Devin SWE-2 Max; design clarification:
GPT-6 Sol on Devin. Evidence: [`runs/T-3q0b/`](../../runs/T-3q0b/).

## Plan and plan review

The plan, its review and the author's complete answer are comments on #350
(6208792576, 6208870269, 6209279856). The plan review (GPT-6.1 Sol on Codex CLI,
`codex exec`) gave `approve-with-conditions`: budget maximum 1,600 lines added
plus removed over 20 files, against `6bc1051`, close-out inside — the base later
confirmed by the Operator as current `main` at landing (6209447438); cycle cap 1;
no panel. Its conditions are applied: the named facts and text inputs are
caller-supplied (`RenderFacts` invents none); the `covers` column was widened to
`list(text)` for `Q-NNN` tokens with the M2c numeric-only regression kept in
`internal/run`; the quote contract states that required requirement IDs enter the
pure checks through `ReqAnswer`; `option <decision>: <option>`, the verdict
spelling and `approve batch` are the decided-here forms.

## What was done

The five pure modules of `internal/spec` — `bet.go` (`Brief`, `CheckBrief`,
`CheckAppetite`, `VerifiedBatch`), `answer.go` (`Answer`, `ReqAnswer`,
`ParseAnswer`), `inventory.go` (the `inventory` and `inventory-index` readers,
`NextInventory`, the SHA-256 of exact version bytes), `render.go` (`RenderPRD`,
`RenderFacts`) and `quote.go` (`CheckQuoteFacts`, `QuoteFactOf`) — their tests,
the schema registration, `loop.md`/`records.md`/`traceability.md`/`PRD-0001` §12
and the evidence file. Round-1 fixes on `8c8e85a`: unordered `covers` tokens no
longer refuse (`covers` is a set, judged after it is filled); option lines are
required only of decisions with more than one survivor; later bets union prior
out-of-scope marks with new ones; the appetite amount is validated by the
explicit grammar `^[0-9]+(\.[0-9]+)?$` before `big.Rat`; an empty index row
reports an `index` finding instead of panicking; the hash-verification caller
contract is documented at its boundary. First-finding-versus-all was declined as
a note — no acceptance criterion requires all findings (finding 6). Finding 1,
the budget, is recorded unresolved below. Mutation evidence: the root's
deterministic replay of the implementation sweep at `2e1e9b4` caught 59 of 68
unique patches by assertion, with 5 build-broken invalids and 4 not-applicable
patches not credited and zero valid survivors; a first post-fix sweep was
discarded after its harness restored stale pre-fix backups, and the valid
re-run at 07:46:15Z–07:47:05Z killed six fix-path mutants — five assertion
failures and one expected panic. **Off the path:** round 2's notes F2–F5 (a lost
out-of-scope mark when the PRD lacks a `**Non-goals.**` paragraph, a
`NextInventory` panic on rows that bypass `ReadInventory`, non-deterministic
finding order from map ranging, and smaller edges) are recorded here for later
issues; none is on the Definition of Done and none was fixed in this branch.

## Review rounds

| Round | Head | Cycle | Cap | Material | Then |
| ----- | ---- | ----- | --- | -------- | ---- |
| 1 | `72528bd` | 0 | 1 | 8 findings, 6 required fixes | fixed on `8c8e85a`, merge `db3330f`, evidence `2e8b944` |
| 2 | `2e8b944` | 1 | 1 | 1 (the budget) | `not mergeable, findings recorded` — budget settled by 6213987433 |

Round 1 (comment 6210579018) was Claude Sonnet 5.5 on Claude Code, a fresh
read-only session, independent of the implementer. Round 2 (comment 6213844910)
was Claude Sonnet 5.5 on Claude Code at effort `high`, context-, method- and
model-independent of the author (Devin SWE-2 Max), the design clarifier (GPT-6
Sol) and the coordinator (GPT-6 family); it reviewed the PRD changelog merge
`db3330f` explicitly — a same-date append keeping both dated entries, a proposed
resolution under the coordinator's ruling (6210747940), not one of O-219's named
§12 exceptions — and found it correct. Its verdict rests on the single material
finding F1, the budget; it found no material code finding and its notes F2–F5 are
the off-path notes above.

## Verdict

Delivered, as code and test result: the checks of the bet and the rendering, with
the widened `covers` schema, the specification decisions and the evidence. The
full integration run on the round-1 head `db3330f` exited 0 (07:51:08.801645Z to
08:03:48.871209Z, 21 packages ok, root's direct subprocess capture); the earlier
run at `2e1e9b4` exited 0 (07:08:36Z to 07:27:36Z); this task's Go is unchanged.
The merge `c9a31c7` imports #357's independently tested register; the combined
block test passes. Its paired-index resolution is covered by O-219 as extended
by 6214113095. **Budget:** maximum 1,600 lines / 20 files, current `main` `0e2f06c`, close-out
inside. The complete landing is **2,001 changed lines / 21 files**, approved
once by the Operator in comment 6213987433 (answer to 6213978483), including
this close-out. This settles the sole material budget finding; the last review
records the then-unapproved budget. No further growth or cycle is approved.

## Resource record

Recorded, not budgeted (ADR-0007). UTC; "not reported" is a value the harness did
not give. Elapsed figures are the recorded windows between dated events, not
session boundaries, where no start was measured.

| Part | Expected tier | Model | Effort | Tokens | Elapsed |
| ---- | ------------- | ----- | ------ | ------ | ------- |
| The plan | reasoning | Claude Sonnet 5.5 (the author of record) | not reported | not reported | not reported (comment 05:29:34Z) |
| The plan review | reasoning | GPT-6.1 Sol on Codex CLI (`codex exec`) | `xhigh` | 173,360 | 5 min 8 s, comment 05:36:07Z |
| The author's answer and design clarification | reasoning | GPT-6 Sol on Devin | not reported | not reported by its own footer | not reported (comment 06:09:29Z); SWE-2 read-only source investigation before it |
| The implementation (T and D), the merge of `f12da6a`, R7, the budget question | execution | Devin SWE-2 Max | not reported | not reported by the CLI | first red run 06:25:44Z to the R6 comment 06:54:01Z |
| The audit replay of the mutation calls | — (deterministic) | not applicable | not applicable | not applicable | 07:02:30Z–07:03:37Z, at `2e1e9b4` |
| The evidence correction (the audit figures, the interrupted run, the piped-status pitfall) | — | Devin SWE-2 Max | not reported | not reported | to commit `e70c910` 07:28:12Z; root's wording correction `72528bd` 07:31:00Z (coordinator, GPT-6 family) |
| The full integration run at `2e1e9b4` | — (deterministic) | not applicable | not applicable | not applicable | 07:08:36Z–07:27:36Z, exit 0 |
| Round 1 | reasoning | Claude Sonnet 5.5 on Claude Code | high | input 30; output 16,381 (thinking 7,354 included); cache creation 71,780; cache read 821,551 | 325,417 ms |
| The round-1 fixes, their mutants, the merge of `8e26eca` | execution | Devin SWE-2 Max | not reported | not reported | triage comment 07:42:38Z to the push of `db3330f` (merge 07:49:41Z) |
| The discarded first fix-mutation sweep | execution | Devin SWE-2 Max | not reported | not reported | inside the window above; discarded, credited nothing |
| The full integration run at `db3330f` | — (deterministic) | not applicable | not applicable | not applicable | 07:51:08.801645Z–08:03:48.871209Z, exit 0 |
| The docs-only evidence correction (the full run, the discarded sweep, the stale-restore pitfall, the changelog ruling) | — | Devin SWE-2 Max | not reported | not reported | R7 comment 08:07:33Z to commit `2e8b944` 08:09:06Z |
| Round 2 | reasoning | Claude Sonnet 5.5 on Claude Code | high | input 44; output 18,933 (thinking 9,890 included); cache creation 90,203; cache read 1,444,112 | 302,455 ms |
| This close-out and merge bookkeeping | — | Devin SWE-2 Max; Codex GPT-6 family | not reported | not reported | not reported by the harness; completed before staging |
| **Total** | | | | not fully reported; reported subtotal 2,636,394 tokens including review cache | not fully reported |

The figures are what each harness reported. The retained Devin CLI outputs give
no token figures; `not reported` is recorded, not guessed. The Sonnet counts
come from the final JSON result, with thinking included in output. The elapsed
windows cite dated events (test runs, commits, comments) where no invocation
start was measured.
