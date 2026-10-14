# T-gkvi — the probe of the forge check on GitHub's real refusal

Issue: [#323](https://github.com/pharzam/layup/issues/323), a defect that the preparation of the demo of `M2d` (row 59b, `T-sdtb`, #284)
revealed, off its path, so a task of its own (gate step 5), as `T-pzgw` was for `M2c`. Serves `F-0003#43` through REQ-003. Base `6c9c314`.
Author: Claude Sonnet 5.5 on Claude Code. Evidence: [`runs/T-gkvi/`](../../runs/T-gkvi/).

## Plan and plan review

The plan (R12), its review and the author's answer are comments on #323. The plan review (GPT-6.1 Sol on Codex CLI, `high`, 4 min 24 s) gave
`approve-with-conditions`: 850 lines over 22 files against `6c9c314`, close-out inside; cap 1; no panel. Its five conditions are applied: every
regression test is written and run red first, in the order unit, job, adapter, world and e2e, with a skeleton that holds no fix; the match
requires the complete probe ref; the job's expectations that move are named (the push is still made on the forge's head); every stale statement
is listed and the missing administration read cites L-A4, not L-A1; and the e2e level is run as `go test -tags=e2e -timeout 10m ./...`. The plan
is reasoning-tier work that the session's model, an execution-tier one, wrote under the Operator's routing message of 2026-10-13: a routing
mismatch by choice. The goal count is 1.

## What was done

GitHub's refusal of a push for a rule violation names the ref and the violated rule and **no ruleset** (ten public reports, listed in the issue), so
the rule "standard error holds the name `layup: the default branch`" could not hold on a real forge, and the demo of `M2d` would have failed on a
correct target. The tests' hook had invented a line `- Rulesets: layup: the default branch`.

1. **The rule** (`internal/run/forgecheck.go`): a refusal passes when its standard error has the line `GH013: Repository rule violations found for
   refs/heads/layup-probe` with the whole ref (`layup-probe-other`, `layup-probe/child` and `layup-probe.x` are not it) and the refusal is
   **attributed**: the effective rules of `layup-probe`, read once, hold at least one `creation` rule and each is of a `Repository` ruleset whose ID is
   the `id` of the apply output of the default ruleset (`setup/rulesets.txt`). Another ruleset's `creation` rule, none, an unread answer, or an output
   with no positive integer `id`, or that is not one JSON object, is `fail`: "probe refused, but not attributed to the default branch's ruleset: ...".
2. **The data**: `forge.Rule` carries `RulesetSourceType`, `RulesetSource` and `RulesetID`, and the adapter reads `ruleset_source_type`,
   `ruleset_source` and `ruleset_id` (the zero value for an absent field or `null`; any other non-integer ID is an error of the call).
3. **The documents**: `scaffold.md` (the probe with a paragraph "Decided here", what fails what, the input states, the acceptance rows), `forge.md`,
   `packages.md`, `architecture.md` §3 item 2 (its parenthetical only; `docs/pdr/` untouched), the pitfall of `T-wphp`, the plan's row 57a, the
   traceability and the cells of REQ-003.

**Tests:** [`test-runs.md`](../../runs/T-gkvi/test-runs.md): the failing output of each red run, the guards named, 20 mutants each alone on a copy of
the tree (all caught), and the gate.

**Known limits:** the attribution reads a snapshot taken just before the push, and the bypass list is the apply output's (L-A4); an organization's
ruleset with a `creation` rule on `layup-probe` makes the row `fail` though the push is refused; GitHub's line is read by one pattern, so a change of
that sentence makes the probe "probe not run", never a pass. GitHub's own text is recorded first-hand by the wrapper of the demo (`T-sdtb`).

## Review rounds

Both rounds: GPT-6.1 Sol on Codex CLI (`xhigh`; ADR-0026 does not name it, the Operator allowed it).

Round 1 (`5446e8e`, cycle 0 of the cap of 1; comment 6168718075): `material`, four findings, each in the change: `rulesetID` read one JSON object and
never looked past it, so `{"id":77}{"id":99}` gave 77 and the probe passed; the contract of the adapter said that a non-integer `ruleset_id` is an error
while `null` reads as zero; two removed test names stayed in the traceability and the cell of REQ-003; the record lacked the failing output of the
red runs and the full e2e level. The reply is comment 6168772008: the first was red first (5 cases, the job's case) and fixed, with two mutants
caught; the contract now says what `null` is; the names are replaced; the record has the output and the gate.

Round 2 (`b4b4d51`, cycle 1 of the cap of 1; comment 6169277061): `nothing material in scope`, two notes: the comment of the adapter's test still
described the old null contract (applied), and the figure of the budget in the brief was low (the frozen diff is 737 lines over 17 files).

## Verdict

Delivered: the probe of the forge check passes on GitHub's real refusal, attributed to the default branch's ruleset by its ID, with the data, the tests
and the documents. The review ended at `nothing material in scope` after two rounds, at the cap of 1. Next: the demo of `M2d` (`T-sdtb`, #284), run
with the binary of this head.

## Resource record

Recorded, not budgeted (ADR-0007). UTC, 2026-10-14, from the clock and the comments of #323; the tokens of GPT-6.1 Sol are the "tokens used" of Codex
CLI; the author's are not reported.

| Part | Expected tier | Model | Effort | Tokens | Elapsed |
| ---- | ------------- | ----- | ------ | ------ | ------- |
| The plan | reasoning (a routing mismatch) | Claude Sonnet 5.5 | `max` | not reported | 11:51 to 11:53 |
| The plan review | reasoning | GPT-6.1 Sol, Codex CLI | `high` | 124,654 | 4 min 24 s, 11:53:50 |
| The answer, the skeleton, the tests, the code, the mutants, the documents, the gate | execution | Claude Sonnet 5.5 | `max` | not reported | 11:58 to 12:33 |
| Round 1 | reasoning | GPT-6.1 Sol, Codex CLI | `xhigh` | 167,499 | 14 min 50 s, 12:33:34 |
| The fix of round 1, its mutants, the gate again | execution | Claude Sonnet 5.5 | `max` | not reported | 12:49 to 13:07 |
| Round 2 | reasoning | GPT-6.1 Sol, Codex CLI | `xhigh` | 172,804 | 10 min 16 s, 13:07:15 |
| The note and the close-out | execution | Claude Sonnet 5.5 | `max` | not reported | 13:17 to 13:18 |
