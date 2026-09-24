# T-7wqc — ADR-0012 (Proposed): successor of T-9tgn and T-w79d after two issue splits

Issue: [#54](https://github.com/pharzam/layup/issues/54). Split predecessors: [#53](https://github.com/pharzam/layup/issues/53) (task [`T-9tgn`](T-9tgn.md)) and [#46](https://github.com/pharzam/layup/issues/46) (task [`T-w79d`](T-w79d.md)); each ended `not mergeable, findings recorded` at its cycle cap, and each detail file stays as its issue's record. This task carries #53's work from `2592624` and fixes the two findings of #53's last round. Record: [ADR-0012](../adr/0012-route-work-across-harness-agents.md). Evidence: `runs/T-w79d/`.

## Operator decisions (quoted; given in the session on 2026-09-24)

- **O-32**: the Operator chose "Successor takes both (Recommended)" — "Close #53 as split. Open a successor from 2592624 that fixes both: rows read up to a blank line (GFM), with a new self-test case; D5 tries fallback rows of the step's tier first, then the other tier with the limit recorded. It gets its own plan review and rounds (cap 2)."
- **O-33**: the Operator chose "Yes, to every successor (Recommended)" — "Record O-33: O-15 applies to each successor that carries #46's goal (ADR-0012) after a split. No new goal-count question for this chain. #49 still fixes the R11 text for all issues."

Earlier decisions that bind this task are quoted in `T-w79d.md` (O-15 to O-18) and `T-9tgn.md` (O-26 to O-29); O-30 and O-31 are on #53.

## The panel's iteration bound was set before the panel ran

#46's [Plan (R12), revision 2](https://github.com/pharzam/layup/issues/46#issuecomment-5809399188) (07:03 UTC) and [revision 3](https://github.com/pharzam/layup/issues/46#issuecomment-5809891207) (07:38 UTC) state "one pass each and no second pass (the iteration bound)". The panel ran from 08:10 UTC; its comments were posted at 08:39 UTC ([A](https://github.com/pharzam/layup/issues/46#issuecomment-5810780834), [B](https://github.com/pharzam/layup/issues/46#issuecomment-5810781152), [C](https://github.com/pharzam/layup/issues/46#issuecomment-5810781436)).

## Which binding ran each role

| Role | Binding | Result |
| -- | -- | -- |
| Plan review | GPT-6 Sol (xhigh) on Devin — outside O-3; a reasoning-class model by its vendor's listing | 364 s, `approve-with-conditions` |

## Affected documents (step 5)

`git grep -n -E "D5|fallback|leading pipe|delimiter row" -- docs/ ':!docs/facts/'`:

- **Stale, fixed:** none outside ADR-0012's D1 and D5, which this task edits.
- **Accurate history, kept (counted by class):** 3 lines in `T-9tgn.md` (the quote of O-26 and O-27, and #53's selection table, which describes #53's D5); 1 line in `T-w79d.md` (the panel's option C1); 6 lines in `docs/links/link-lint.sh` and `docs/tests/nested-checkout-check.sh` ("fallback walk", another sense).
