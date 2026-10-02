# T-4wrw — the PDR of LAYUP

Issue: [#99](https://github.com/pharzam/layup/issues/99). Parent
[#42](https://github.com/pharzam/layup/issues/42), child 6, the last one: after
the implementation plan (`T-55n2`, #76). Serves `F-0003#51`. Base `b48764f`.
Author: Claude Opus 5.5 on Claude Code. The record is
[`docs/pdr/PDR-0001.md`](../pdr/PDR-0001.md); the evidence is under
[`runs/T-4wrw/`](../../runs/T-4wrw/).

## Plan and plan review

The plan (R12) and its review are comments on #99. The plan review (Claude Fable
5.1 on Claude Code, a fresh session) gave `approve-with-conditions`: Budget
maximum 800 lines added plus removed over 16 files on the whole diff against
`b48764f`, close-out inside; Cycle cap 1. The author applied the three
conditions and the thirteen notes (comment "Plan: the author's answer to the
plan review" on #99).

## What was done

1. **Red** (plan step 1, condition 3): the one-off check
   [`trace-check.py`](../../runs/T-4wrw/trace-check.py), run on a copy of
   `PRD-0001` with one seeded defect per part, then on the real tree before the
   record existed; each run failed for its reason
   ([`trace-check.md`](../../runs/T-4wrw/trace-check.md)).
2. **The record** (plan step 3): [`docs/pdr/PDR-0001.md`](../pdr/PDR-0001.md),
   with the approval part pending, and the request for the approval on #99. Rows
   1 and 2 of the plan merged after `b48764f` (#100, #101); the branch merged
   `origin/main` in (`a432c8a`), and the record says what those rows changed in
   the approved documents (a note on #99, before the approval).
3. **The approval copied** (plan step 4): O-129 below, in the record and here;
   `PRD-0001` has the Status `Accepted`, with a row in its change log; the index
   row of `docs/prd/README.md`; the onboarding document (the `prd/` row, a new
   `pdr/` row, and the sentence on #42 under "Where the project stands"); the
   glossary row "Preliminary Design Review"; a `docs/pdr/` row in `README.md` and
   in the table "Sources of truth" of `AGENTS.md` (note 11); the check green, 6 of
   6.
4. **The children of #42** on the forge (note 13): #43 (PR #44), #45 (PR #62),
   #63 split at its cycle cap to #64 (PR #65), #66 closed as not planned and
   replaced by #72 (PR #73), #74 (PR #75) and #76 (PR #98) are closed with
   merged pull requests; #99 closes with this one. #29 lists the plan's issues
   (re-scoped by #76).

## The Operator's decisions

**O-129.** The Operator approved the PDR on 2026-10-02 at 06:37:30 UTC, in
[a comment on #99](https://github.com/pharzam/layup/issues/99#issuecomment-5946822016)
from the Operator's own account, with no App: "The PDR (#99): approved". The
approval covers the four documents at `b48764f`, asks for no change, and answers
none of the open items of the record. The number was taken at copy time (note 8
of the plan review): O-127 and O-128 went to rows 1 and 2.
