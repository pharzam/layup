# The count of the first pilot's numbers

Task `T-gd8q` ([#120](https://github.com/pharzam/layup/issues/120)), the first step
of the ADR that ends bootstrap mode (`gov-pilot-numbers`,
[`runs/T-55n2/inventory.md`](../T-55n2/inventory.md)). Read at `34dd858`. Each value
cites a `file:line`; [`cites.sh`](cites.sh) prints the text of each cited line into
[`cites.txt`](cites.txt), so a reader can check that the line holds the value.

**The rule was written after the numbers were known.** The inventory asked the plan
of the pilot to pre-register the thresholds (`runs/T-55n2/inventory.md:1544`); the
plan of #97 did not. So this count is not a blind test
(`docs/guardrails.md:21`). The Operator accepted the rule on that basis (O-159
part 1, [#120](https://github.com/pharzam/layup/issues/120)).

## The rule

- **Populations.** P1: the 20 rows of the plan of phase 1
  (`docs/plan/README.md:96`). P2: each other task of `docs/tasks/completed.md`
  dated 2026-09-25 (the start of bootstrap mode) or later. Each is read from its
  `docs/tasks/<ID>.md`.
- **A round** is a blind review on a frozen head that returned a record. A run that
  was skipped or stopped is not a round. A "fresh verification" before the freeze
  (rows 1 and 2) is not a round.
- **A material finding** is one that the record classes as material, after any edit
  of the verdict.
- **A review ends by decay** when its last round is `nothing material in scope` and
  no Operator decision raised the cap, ruled a finding or closed out a fix with no
  further round. Otherwise it **ends by decision**.
- **Two units of a stall.** A *stall event* is each `not mergeable, findings
  recorded` verdict, each Operator decision that raised a cap, ruled a finding a
  note or a known limit, or closed out a fix with no further round, each split, and
  each plan review with `reject`. A *stall task* is a task with at least one stall
  event. A budget approval is an Operator decision, not a stall.
- **An Operator decision of a task** is an O- number that the task asked and that
  was answered inside the time of the task. A cited O- number is not counted.
- **Product or process.** A task is product when it is a row of the plan, a parent
  of rows, or a task whose deliverable is a fact, the PRD, the architecture, the
  specification, the plan, the PDR or LAYUP's code. Every other task is process. A
  task with both is product, and the table names it.
- **A pilot defect's row** is the row of the plan whose task delivered the code of
  the finding, by the plan table and the task file. A finding can have more than
  one row.

## P1: the 20 rows of phase 1

| # | Task | Plan review | Cap | Rounds: verdict (material) | End | Decisions asked | Source |
| - | ---- | ----------- | --- | -------------------------- | --- | --------------- | ------ |
| 1 | `T-18v6` | awc | 1 | 1: nothing (0) | decay | O-127 (budget) | `docs/tasks/T-18v6.md:13`, `:81`, `:55` |
| 2 | `T-2tc2` | awc | 1 | 1: nothing (0) | decay | O-128 (budget) | `docs/tasks/T-2tc2.md:13`, `:87`, `:58` |
| 3 | `T-2yw7` | awc | 1 | 1: nothing (0) | decay | — | `docs/tasks/T-2yw7.md:11`, `:70` |
| 4 | `T-3jpx` | awc | 1 | 1: nothing (0) | decay | — | `docs/tasks/T-3jpx.md:11`, `:66` |
| 5 | `T-5sgt` | awc | 1 → 2 | 3: material (3), not mergeable → material (1), nothing (0) | decision: raise | O-130 | `docs/tasks/T-5sgt.md:11`, `:61`, `:78`, `:84`, `:114` |
| 6 | `T-5zmw` | awc | 1 | 2: material (3), nothing (0) | decay | O-131 | `docs/tasks/T-5zmw.md:10`, `:90`, `:115`, `:16` |
| 7 | `T-6x75` | awc | 1 → 2 → 3 → 5 | 5: material (5), not mergeable → material (2), (1), (1), nothing (0) | decision: raise ×3 | O-132, O-133, O-134, O-135 | `docs/tasks/T-6x75.md:11`, `:64`, `:89`, `:103`, `:115`, `:129` |
| 8 | `T-79y7` | awc | 1 | 2: material (5), nothing (0) | decay | O-136 | `docs/tasks/T-79y7.md:11`, `:56`, `:83`, `:18` |
| 9 | `T-7s0y` | awc | 1 | 2: material (1), nothing (0) | decay | O-137 | `docs/tasks/T-7s0y.md:15`, `:88`, `:119`, `:11` |
| 10 | `T-8vpw` | awc | 2 | 2: material (1), nothing (0) | decay | O-138 | `docs/tasks/T-8vpw.md:15`, `:63`, `:91`, `:20` |
| 11 | `T-8ya0` | awc | 1 | 1: nothing (0) | decay | — | `docs/tasks/T-8ya0.md:14`, `:55` |
| 12 | `T-9t1q` | awc | 2 | 1: nothing (0) | decay | — | `docs/tasks/T-9t1q.md:16`, `:65` |
| 13 | `T-b3r1` | awc | 1 | 2: material (1), nothing (0) | decay | — | `docs/tasks/T-b3r1.md:13`, `:82`, `:114` |
| 14 | `T-c06a` | awc | 1 | 2: material (1), nothing (0) | decay | — | `docs/tasks/T-c06a.md:15`, `:68`, `:109` |
| 15 | `T-d6q5` | awc | 1 | 2: nothing → material (0 by the reviewer, 1 by CI), nothing (0) | decay | — | `docs/tasks/T-d6q5.md:14`, `:104`, `:129`, `:155` |
| 16 | `T-dep6` | awc | 1 | 1: nothing (0) | decay | — | `docs/tasks/T-dep6.md:13`, `:81` |
| 17 | `T-tmhw` | awc | 1 | 2: material (2), nothing (0) | decay | — | `docs/tasks/T-tmhw.md:13`, `:63`, `:89` |
| 18 | `T-dgy7` | awc | 1 | 1: nothing (0) | decay | — | `docs/tasks/T-dgy7.md:13`, `:75` |
| 19 | `T-efmy` | awc | 1 | 1: nothing (0) | decay | — | `docs/tasks/T-efmy.md:14`, `:91` |
| 20 | `T-evad` | awc | 2 | 2: material (1), nothing (0) | decay | O-139 to O-158 (20; O-146 a budget) | `docs/tasks/T-evad.md:22`, `:159`, `:161`, `:14` |

"awc" is `approve-with-conditions`; "nothing" is `nothing material in scope`.

| Measure | P1 |
| ------- | -- |
| Plan reviews | 20, all `approve-with-conditions`; 0 `reject` |
| Rounds | 35: 9 tasks with 1, 9 with 2, 1 with 3, 1 with 5 |
| Material findings | 28 (29 with the CI failure of row 15) |
| Tasks with a material finding in round 1 | 10 |
| Tasks where a round after a fix found a new material finding | 2 (rows 5 and 7) |
| Ends | 18 by decay; 2 by decision (rows 5 and 7) |
| Stall tasks | 2 |
| Stall events | 8: 4 `not mergeable` verdicts (row 5 round 2; row 7 rounds 2 to 4), each edited to `material` when the cap was raised, and 4 raises of the cap (O-130, O-133, O-134, O-135) |
| Operator decisions asked | rows 1 to 19: 11 (O-127, O-128, O-130 to O-138), of which 2 budget approvals and 4 raises; row 20: 20 (O-139 to O-158), of which 1 budget approval (O-146) |
| Rows whose record names a skipped reviewer harness | 15 of 20 (rows 1, 2 and 8 to 20); most records do not count the runs |

## P2: the other tasks under bootstrap mode

| Task | Class | Plan review | Cap | Rounds: verdict (material) | End | Stall events | Source |
| ---- | ----- | ----------- | --- | -------------------------- | --- | ------------ | ------ |
| `T-stfn`, `T-b97r`, `T-vk3k` | product (parents of rows) | no task file | — | — | — | — | `docs/tasks/completed.md:24` to `:26` |
| `T-meh2` | product (the PDR parent) | no task file | — | — | — | — | `docs/tasks/completed.md:45` |
| `T-4wrw` | product (PDR) | awc | 1 | 1: nothing (0) | decay | — | `docs/tasks/T-4wrw.md:12`, `:67` |
| `T-55n2` | product (plan) | awc | 1 | 2: material (1), nothing (0) | decay | — | `docs/tasks/T-55n2.md:16`, `:109` |
| `T-0drh` | product (specification) | awc | 1 → 2 → 3 | 4: material (5), not mergeable (1), not mergeable (1), not mergeable (1) | decision: fix with no round (O-119) | 3 `not mergeable`; raise O-116, O-117; fix with no round O-119 | `docs/tasks/T-0drh.md:43`, `:26`, `:30`, `:34`, `:91` |
| `T-hbw8` | product (architecture) | awc | 1 | 1: nothing (0) | decay (the gate; its design path ended by O-112) | the design path only: 6 raises, 1 known limit (O-110), 1 fix with no round (O-112) | `docs/tasks/T-hbw8.md:12`, `:132`, `:150` |
| `T-84r5` | product (PRD) | by the Operator | 1 | 1: nothing (0) | decay | — | `docs/tasks/T-84r5.md:18`, `:42` |
| `T-wjq4` | product (PRD; also a linter fix) | by the Operator | 2 | 3: material (5), material (13), not mergeable (8) | decision: split (O-49) | 1 `not mergeable`; 1 split | `docs/tasks/T-wjq4.md:19`, `:23`, `:27` |
| `T-zmj6` | product (fact `F-0004`; also check `facts`) | not recorded | 1 | 2: material (4), not mergeable (2) | decision: ruled notes (O-46) | 1 `not mergeable`; 1 ruling | `docs/tasks/T-zmj6.md:20`, `:22` |
| `T-745n` | process | awc | 2 | 2: material (4), nothing (0) | decay | — | `docs/tasks/T-745n.md:11`, `:130` |
| `T-7sbn` | process (ADR-0012) | by the Operator | 1 | 2: material (3), not mergeable (2) | decision: ruled notes (O-44) | 1 `not mergeable`; 1 ruling | `docs/tasks/T-7sbn.md:21`, `:23` |
| `T-8ywj` | process (ADR-0012, the revert of `F-0005`) | by the Operator | 1 | 2: material (13), not mergeable (5) | decision: split | 1 `not mergeable`; 1 split | `docs/tasks/T-8ywj.md:20`, `:22`, `:26` |

| Measure | P2 (the 10 tasks with a task file) |
| ------- | ---------------------------------- |
| Rounds | 20 (`T-hbw8` also had 24 design rounds, not counted) |
| Material findings | 68 (the 5 of `T-8ywj` round 2 counted; its record does not class them) |
| Ends | 5 by decay; 5 by decision: 2 splits, 2 ruled notes, 1 fix with no round |
| Stall tasks | 5 on the gate path (`T-0drh`, `T-wjq4`, `T-zmj6`, `T-7sbn`, `T-8ywj`); `T-hbw8` on its design path |
| Stall events, the gate path | 7 `not mergeable` verdicts, 2 raises, 2 ruled notes, 1 fix with no round, 2 splits, 0 `reject` |

**Every end by decision under bootstrap mode** (P1 and P2, the gate path) was one of
four kinds: raise the cap (rows 5 and 7, `T-0drh`), rule the findings notes
(`T-7sbn`, `T-zmj6`), close out a fix with no further round (`T-0drh`), or split
(`T-8ywj`, `T-wjq4`). The text of the gate named only the split
(`docs/engineering-discipline.md:319`).

## Product and process

| Lines of `docs/tasks/completed.md` | Product | Process |
| ---------------------------------- | ------- | ------- |
| 2026-09-25 or later (lines 23 to 56) | 31 (P1: 20; P2: 11) | 3 (`T-745n`, `T-7sbn`, `T-8ywj`) |
| Before 2026-09-25 (lines 57 to 73) | 5 | 12 |

## The pilot's defects in phase-1 code, by row

Each finding of [`runs/T-evad/findings.md`](../T-evad/findings.md) in the code (or
the specification of the code) of a plan row. "Rounds" is the number of review
rounds of each row.

| Finding | findings.md | Row (rounds) | Differs from `runs/T-evad/numbers.md:13` |
| ------- | ----------- | ------------ | ----------------------------------------- |
| F-1 | `:59` | 14 (2); 15 (2) for S12 | names row 14 only |
| F-2 | `:77` | 8 (2); 9 (2) | — (by specification) |
| F-3, F-24 | `:30`, `:33` | 13 (2); 10 (2) for the scanner | names row 13 only |
| F-6 | `:80` | 7 (5); 15 (2) | names row 15 only |
| F-8 (a) | `:25` | 13 (2) | agrees |
| F-8 (b) | `:27` | 10 (2) | agrees |
| F-13 | `:81` | 8 (2); 15 (2) | names row 15 only |
| F-16 | `:86` | 11 (1); 13 (2) | names row 13 only |
| F-19 | `:110` | 15 (2) | agrees |
| F-21 | `:91` | 6 (2); rule G1 first came from `T-dq05`, before bootstrap mode | — |
| F-22, F-23 | `:93`, `:103` | 14 (2) | — |
| F-25 | `:33` | 13 (2); 11 (1), whose check took the fix | — |
| F-30 | `:49` | 8 (2); the defect is in the specification text | — |
| F-31 | `:115` | 9 (2); 13 (2); a gap of the specification | — |

**The reading.** Each of these findings is in at least one row with two or more
rounds, and row 7, with five rounds, holds F-6. One row with one round, row 11,
shares two findings with row 13 (F-16, F-25). Most are of integration or of the
specification (a marker over more lines, a second setup run, a parameter that
GitHub adds to a ruleset, a gate that a pull request can weaken). This shows a limit
of a review of one diff, not that one round is enough: the rows with more rounds are
also the larger rows, which the pilot used most.

## The full gate in the pilot's target

The target task `T-vu2j` ran under the gate of the kit as written, with one round of
four lenses on each frozen head (a workaround for F-15): 3 plan reviews (`reject`,
`reject`, `approve-with-conditions`) and 5 rounds with 4, 8, 3, 2 and 1 material
findings; three Operator decisions ended it (O-156 to O-158; F-38)
(`runs/T-evad/numbers.md:131`, `runs/T-evad/findings.md:133`,
`docs/tasks/T-evad.md:49`).
