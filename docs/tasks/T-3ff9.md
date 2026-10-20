# T-3ff9 — put ADR-0028 into effect

Issue [#447](https://github.com/pharzam/layup/issues/447), child of the M2e throughput investigation [#438](https://github.com/pharzam/layup/issues/438) and follow-up of [#441](https://github.com/pharzam/layup/issues/441). Serves F-0003#48; goal count 1, the operative text of one accepted record.
Plan [6297702690](https://github.com/pharzam/layup/issues/447#issuecomment-6297702690); plan review [6298181777](https://github.com/pharzam/layup/issues/447#issuecomment-6298181777) (Claude Opus 5.5 on OpenCode, variant max, approve-with-conditions); amendment 1 [6298337523](https://github.com/pharzam/layup/issues/447#issuecomment-6298337523).
The Operator chose option 1, execution data only, for "input" in [6298369117](https://github.com/pharzam/layup/issues/447#issuecomment-6298369117). Maximum 650 lines added plus removed over 20 files, evidence and close-out included, against current `main` at landing; cap 1; no panel.
Evidence: [test runs](../../runs/T-3ff9/test-runs.md), [conditions](../../runs/T-3ff9/conditions.md), [semantic comparison](../../runs/T-3ff9/semantic.md), [reference sweep](../../runs/T-3ff9/mentions.tsv), [task check](../../runs/T-3ff9/checks.sh).

## Delivered and boundaries

- [x] Reviewed plan and amendment; all three review conditions met, as mapped in the conditions table.
- [x] The canonical [review protocol](../engineering-discipline.md#reviewing-until-findings-decay) puts ADR-0028 into effect: evidence records and close-out bookkeeping leave the budget count, executable tests and inputs stay, and mixed files retain their non-bookkeeping lines and touched-file count.
- [x] The accepted tolerance (20 %, at least 60 and at most 500 lines, plus 4 files), the post and the Operator's overrule, and one exact approval beyond the tolerance are operative there.
- [x] Default end 3 applies only when every open material finding is a records finding; its exclusions, disputed classification, unread-fix disclosure and retained verdict remain explicit. Ends 1, 2 and 4 are preserved.
- [x] "Input" is defined once for both the count and end 3: execution data, excluding records read only to audit form, citations or evidence; dual-use files remain inputs. The pitfall is in `guardrails.md` §2.
- [x] Entry-point, plan, glossary and DoD pointers, and the Status lines and index cells of ADR-0007 and ADR-0026, follow the canonical home. Accepted ADR bodies remain unchanged. Round 1 supplies the independent clause-by-clause semantic pass; its remaining notes are disposed of below.
- K5 implementation, fixtures, runner support, mutation evidence and real-task replays moved before S1 to [T-c826, #449](https://github.com/pharzam/layup/issues/449). K5 still counts the whole diff and adds **70 lines and 2 files**; until #449 lands, it counts more than the new rule. No K5 or review-record parser change is delivered here.
- The figures of ADR-0028 and the M2e rows are outside this task. T-3ff9 uses its baseline 650/20 maximum and cap 1, counting evidence and close-out; its own proposed delegation does not authorize its landing.

## Verification and verdict

**Red.** S1, `9c19eff`: `sh runs/T-3ff9/checks.sh` exited 1 at 18:43:41Z on 2026-10-20, before operative edits. New rules and pointers were absent and retired wording remained. Existing unit, ends 1/2/4, immutable ADR bodies and parser/K5 scope passed as regression controls. The earlier commit attempt exited 128 because `.git` was read-only; it ran neither hooks nor the red check. Commits were then delegated to the coordinating session.

**Green.** The unchanged S1 check passed all 58 assertions in Stage B at 18:48:04Z and on committed head `43c1bd4` in Stage C at 18:54:11Z. The latter follows the clean merge of `origin/main` (`df9bfac`).

**Gates.** Stage C's ADR, PRD and link linters, 138 discipline tests, 62 freeze-check fixtures, full setup check and `git diff --check origin/main HEAD` all exited 0. The command table in the test-runs record gives measured UTC endpoints and results.

**Pre-freeze.** On `43c1bd4`, all K1–K6 passed, exit 0 at 18:54:55Z. K5 projected **516/16** from a whole diff of **446/14** plus 70/2; margin to 650/20 was 134/4. The 33-line Stage C record then made the frozen head's whole diff **479/14**. These projections do not replace the final landing measurement.

**Round 1** ([6299443653](https://github.com/pharzam/layup/issues/447#issuecomment-6299443653)): Claude Opus 5.5 (`anthropic/claude-opus-5-5`, variant max) on OpenCode, cycle 0, **nothing material in scope** on frozen head `a6dd03a37274ee9a6e2a97be2b4b970876449e2b`. Five notes, no material finding; no further round follows at cap 1. The review includes the clause-by-clause semantic pass against ADR-0028 and the Operator's input decision.
The reviewer reports green text gates, 10 nested-checkout cases, 62 freeze-check fixtures, full setup, whitespace, Go build, vet and unit tests. K5 on the frozen head projected 549/16 (479/14 plus 70/2), margin 101/4. The reviewer's shell permission refused `checks.sh`; it did not claim to run it, instead replicating its text assertions and inspecting the ADR/index diff. Stage C supplies the actual green script run.

**Close-out validation.** On 2026-10-20, `sh runs/T-3ff9/checks.sh` passed all 58 checks (19:31:03.831660Z–19:31:04.046055Z, exit 0). `sh .githooks/pre-commit` passed its text gates, formatting, `go vet ./...` and `go test ./...` (19:31:04.046162Z–19:31:33.881272Z, exit 0). This invokes the hook without claiming a commit; the final evidence-only paragraphs were written afterward.

**Landing measurement.** Against `origin/main` at `df9bfac`, `git diff --numstat origin/main` plus the untracked task file's `git diff --no-index --numstat -- /dev/null docs/tasks/T-3ff9.md` gives **501 added + 46 removed = 547 lines over 16 files**, evidence and close-out included. Margin to the baseline 650/20 maximum: **103 lines and 4 files**. No tolerance or evidence exclusion is used for this task.

**Note dispositions.** N1–N4 are declined in this close-out and routed to [T-ebm4, #451](https://github.com/pharzam/layup/issues/451). Applying any of them changes files outside the task's own detail file and indexes, beyond the close-out allowed by "A fix re-freezes".

- **N1, R12's scale sentence:** declined here because changing `issue-workflow.md` is an operative edit. #451 carries the pointer fix reconciling the sentence with the canonical tolerance; the conflicting restatement remains a known note at this landing.
- **N2, Size class calls the maximum a ceiling:** declined here because the glossary edit is outside close-out. #451 carries the pointer to the canonical rule rather than leaving the maximum described as the absolute bound.
- **N3, glossary terms and examples:** declined here because adding rows for "records finding", "tolerance" and "input", and concrete examples for Cycle cap and Ceiling, changes the glossary. #451 carries them; the terms already have definitions in their operative home.
- **N4, curly quotes:** declined here because normalizing the canonical section and its matching task assertions changes both `engineering-discipline.md` and `checks.sh`. #451 carries the normalization, including matching assertions; the straight-character absence checks' stated limitation remains.
- **N5, margin and time to overrule:** referred to the Operator on [#447](https://github.com/pharzam/layup/issues/447#issuecomment-6299443653), outside #451. The margin could mean maximum minus figure (K5 today) or tolerance limit minus figure (#449's proposed design); "before the merge" supplies no waiting interval. Both readings come from ADR-0028. Settling them needs the Operator's definition or figure; this close-out invents neither. #449 can implement the margin once settled.

**Verdict: nothing material in scope.** The reviewed head stands, with N1–N4 routed to #451 and N5 to the Operator on #447. Only close-out bookkeeping follows the round; no unread implementation fix is claimed.

**Author and harnesses.** The Operator's choice of GPT-6.1 Sol on Codex CLI for the plan, implementation and close-out was recorded in the claim [6297465432](https://github.com/pharzam/layup/issues/447#issuecomment-6297465432) on 2026-10-20. Harness exception: ADR-0026's routing rule names `claude`, `devin` and `opencode`; Codex is used on the Operator's message of 2026-10-13. OpenCode is named by ADR-0026 and needs no harness exception. The reason for passing over earlier reviewer harnesses is not reported in the supplied record.

## Resource record

[ADR-0007](../adr/0007-record-task-resource-use.md), **recorded, not budgeted**. All intervals below are on 2026-10-20, UTC; elapsed durations are derived from the supplied measured endpoints. Tokens are the harness reports; OpenCode includes cache tokens and is kept separate in the totals. Costs are reported only for OpenCode.

| Part | Expected tier | Model and harness | Effort | Tokens | Elapsed |
| --- | --- | --- | --- | --- | --- |
| Plan | reasoning | GPT-6.1 Sol on Codex CLI | xhigh | 135,595 | 17:54:57Z–18:04:35Z; 9m 38s |
| Plan review | reasoning | Claude Opus 5.5 on OpenCode | variant max | 11,872,965 with cache; USD 6.24 | 18:05:07Z–18:29:11Z; 24m 04s |
| Amendment 1 | reasoning | GPT-6.1 Sol on Codex CLI | xhigh | 79,874 | 18:29:50Z–18:35:59Z; 6m 09s |
| First build attempt, blocked by read-only `.git` | execution | GPT-6.1 Sol on Codex CLI | xhigh | 77,177 | 18:38:14Z–18:41:10Z; 2m 56s |
| Stage A: S1 check and red evidence | execution | GPT-6.1 Sol on Codex CLI | xhigh | 81,439 | 18:41:41Z–18:44:48Z; 3m 07s |
| Stage B: canonical text, pointers and evidence | execution | GPT-6.1 Sol on Codex CLI | xhigh | 123,696 | 18:45:37Z–18:52:06Z; 6m 29s |
| Stage C: committed-head gates and pre-freeze record | — | GPT-6.1 Sol on Codex CLI | high | 55,531 | 18:53:42Z–18:55:43Z; 2m 01s |
| Round 1 | reasoning | Claude Opus 5.5 on OpenCode | variant max | 15,904,209 with cache; USD 7.61 | 19:00:27Z–19:25:57Z; 25m 30s |
| Stage D: close-out | — | GPT-6.1 Sol on Codex CLI | xhigh | not reported at writing | not reported at writing |
| Coordinating session: claim, posting, commit-plan commits, merge of `main`, push | — | Claude Opus 5.5 on Claude Code | not reported | not kept | not reported |
| **Total reported** | | | | **553,312 Codex; 27,777,174 OpenCode with cache; USD 13.85** | **79m 54s summed measured parts; Stage D and coordination not reported** |

The coordinating session wrote none of the reviewed change. This stage leaves the two close-out files for its one commit; no commit or push is claimed from this read-only-`.git` session.
