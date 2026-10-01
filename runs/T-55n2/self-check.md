# T-55n2 — the author's self-check of the plan (plan step 6)

The workflow `t55n2-selfcheck` (run `wf_625ca7fb-cb4`, 2026-10-01) ran four critics of Claude Opus 5.5, the author's model, each with one lens (order, coverage, claims, decisions), on the plan at `71c1bfb`. Each finding went to three refuters, and it survives with two confirmations. By the Operator's decision O-125, the self-check stopped after round 1; the workflow was stopped when 103 of its 113 agents had returned, so the refuters of the lenses `claims` and `decisions` had not finished, and their findings are marked `unverified`. The author read each finding against the text. This is the author's self-check, not a review: the review round of the gate comes later, by another model.

Result: 73 findings; 28 survive, 5 refuted, 40 unverified.

| Lens | # | Status (confirmed of votes) | Severity | Where | What the author did |
| ---- | - | --------------------------- | -------- | ----- | ------------------- |
| order | 1 | survives (3 of 3) | material | docs/plan/README.md:96 (task table, row 10, After column) | applied: row 10 after 12 |
| order | 2 | survives (3 of 3) | material | docs/plan/README.md:202 (defect register K17) and :95 (row 9, After column) | applied: row 9 after 10; K17 now settled by row 12 |
| order | 3 | survives (3 of 3) | material | docs/plan/README.md:204 (defect register K19) and :93 (row 7, After column) | applied: K19 is read by row 15 only |
| order | 4 | survives (3 of 3) | material | docs/plan/README.md:105-108 (rows 19 and 20, After column; the sentence after the table); docs/glossary.md:147 | applied: row 20 after 19; the pilot is the last task |
| order | 5 | survives (3 of 3) | material | docs/plan/README.md:92 (row 6, Items and After) and :20 | applied: psb-batch-api moved to row 9; the seam is decided in the plan |
| order | 6 | survives (3 of 3) | material | docs/plan/README.md:94 (row 8, Tests column) | applied: the stop-and-resume e2e test moved to row 9 |
| order | 7 | survives (3 of 3) | material | docs/plan/README.md:132 (What phase 1 proves, row NFR-004); docs/tests/traceability.md:32 | applied: NFR-004 names rows 5, 10 and 15; the handle moved to T-8vpw |
| order | 8 | survives (3 of 3) | note | docs/plan/README.md:89 (row 3, After column) | applied: row 3 after 2 only |
| order | 9 | survives (3 of 3) | note | docs/plan/README.md:101 (row 15, Task column) | applied: row 15 names S12, S13 and S15 |
| order | 10 | refuted (0 of 3) | note | docs/plan/README.md:196 (defect register K11) | refuted by the refuters; no change |
| order | 11 | refuted (0 of 3) | note | docs/plan/README.md:209 (defect register K24) | refuted by the refuters; no change |
| order | 12 | survives (3 of 3) | note | docs/plan/README.md:87 (row 1) and :193 (defect register K8) | applied: option A decided in the plan |
| order | 13 | survives (3 of 3) | note | docs/plan/README.md:129 (What phase 1 proves, row NFR-001) | applied |
| order | 14 | survives (3 of 3) | note | docs/plan/README.md:4-5 and :186 (defect register K1) | applied: the After column gives the order |
| order | 15 | survives (3 of 3) | note | runs/T-55n2/plan-check.py:158-187 | applied: two checks added to plan-check.py (edges; settled before read) |
| coverage | 1 | survives (3 of 3) | material | docs/plan/README.md:96 (row 10, Items and After cells); docs/plan/README.md:98 (row 12, After cell `7`) | applied (as order 1) |
| coverage | 2 | survives (3 of 3) | material | docs/plan/README.md:92 (row 6, Items and After cells) | applied (as order 5) |
| coverage | 3 | survives (3 of 3) | material | docs/plan/README.md:197 (defect register, K12) | applied: K12 split; row 13 decides the kit-history part |
| coverage | 4 | survives (3 of 3) | material | docs/plan/README.md:199 (K14); :98 (row 12 After `7`); :95 (row 9 After `4, 6, 8, 12`); :203 (K18) | applied: K14 read by rows 6, 9, 12, 13 |
| coverage | 5 | survives (3 of 3) | material | docs/plan/README.md:101 (row 15, Task cell); :99 (row 13); :201 (K16) | applied: row names; K16 read by row 8 |
| coverage | 6 | refuted (1 of 3) | material | docs/plan/README.md:147 (hosts table); :222 (K37) | refuted by the refuters; no change |
| coverage | 7 | survives (3 of 3) | material | docs/plan/README.md:228-229; rows 3, 13 and 14 (:89, :99, :100); :196 (K11) | applied: the duplicate resolutions in the plan |
| coverage | 8 | refuted (1 of 3) | material | docs/plan/README.md:108; :105-106 (rows 19 and 20, both After `16, 17, 18`); docs/glossary.md:147 ("It is the last task  | refuted by the refuters; no change |
| coverage | 9 | survives (3 of 3) | note | docs/plan/README.md:132; docs/tests/traceability.md:32 | applied (as order 7) |
| coverage | 10 | survives (3 of 3) | note | docs/plan/README.md:210 (K25) | applied: K25 is the plan's choice |
| coverage | 11 | survives (3 of 3) | note | docs/plan/README.md:192 (K7), :202 (K17), :209 (K24), :214 (K29), :224 (K39); also K8 (:193) and K31 (:216) | applied: readers added; K24 settled by rows 14, 15; K31 settled by row 2 |
| coverage | 12 | survives (3 of 3) | note | docs/prd/PRD-0001-layup.md:212-214, :218, :226, :230, :238, :240; docs/plan/README.md:46 (M2b) | applied: the §12 cells follow one rule; REQ-013 in M2b |
| coverage | 13 | survives (3 of 3) | note | docs/plan/README.md:149, :237 | applied: the sentence of docs/ci/README.md fixed in this task |
| coverage | 14 | survives (3 of 3) | note | docs/plan/README.md:129 | applied |
| coverage | 15 | survives (3 of 3) | note | docs/plan/README.md:197 (K12 note) and :106 (row 20 Requirements, REQ-018); docs/prd/PRD-0001-layup.md:235 | applied: new register row K42 (row 13) |
| coverage | 16 | refuted (0 of 3) | note | docs/plan/README.md:100 (row 14 Tests and After `4, 5`); docs/tests/traceability.md:31 | refuted by the refuters; no change |
| claims | 1 | survives (3 of 3) | material | docs/plan/README.md:96 (row 10), with :98 (row 12) | applied (as order 1) |
| claims | 2 | unverified (1 of 1) | material | docs/plan/README.md:202 (K17) and :204 (K19); rows at :95 and :93 | applied (as orders 2, 3) |
| claims | 3 | survives (2 of 2) | material | docs/plan/README.md:108; rows :105-106; docs/glossary.md:147 | applied (as order 4) |
| claims | 4 | unverified (0 of 0) | material | docs/plan/README.md:106 (row 20); docs/tests/traceability.md:40; docs/tasks/T-55n2.md:46-50 (O-123) | applied (K42) |
| claims | 5 | unverified (0 of 0) | material | docs/prd/PRD-0001-layup.md:218, :219, :221, :224, :236 (with the preamble at :212-214 and the §13 row at :251) | applied (§12 cells) |
| claims | 6 | unverified (0 of 0) | material | docs/prd/PRD-0001-layup.md:218, :221, :238, :240; docs/tests/traceability.md:40 | applied (§12 cells; the T-evad traceability row) |
| claims | 7 | unverified (0 of 0) | material | docs/plan/README.md:129, :124, :95 | applied |
| claims | 8 | unverified (0 of 0) | material | docs/plan/README.md:123 | applied: the F-0004 clause in the acceptance table |
| claims | 9 | unverified (0 of 0) | material | docs/glossary.md:147; docs/prd/PRD-0001-layup.md:144-146; docs/facts/F-0004-psb-gap-answers.md:28 | applied: the glossary collision names PRD §8 and F-0004#11; PRD §8 unchanged |
| claims | 10 | unverified (0 of 0) | material | docs/plan/README.md:197 (K12); also :210 (K25) | applied (as coverage 3) |
| claims | 11 | unverified (0 of 0) | material | docs/tests/dod-checklist.md:25, :73-74, :80; docs/tests/README.md:42; docs/tests/traceability.md:12-14 | applied: dod-checklist and docs/tests/README.md point to traceability.md |
| claims | 12 | unverified (0 of 0) | note | docs/tests/traceability.md:28, :29, :39, :40 | applied in part: REQ-009 and REQ-011 off the T-18v6 row; a multi-requirement row gives the first fact and ADR, as the table now says |
| claims | 13 | unverified (0 of 0) | note | docs/tests/traceability.md:4 | applied: the table says why those tests have no row yet |
| claims | 14 | unverified (0 of 0) | note | docs/plan/README.md:116 | applied: the label says the plan's reading until the PDR |
| claims | 15 | unverified (0 of 0) | note | docs/plan/README.md:224; :60-62 | applied: K39 read by M2a and M3a |
| claims | 16 | unverified (0 of 0) | note | README.md:58 (also AGENTS.md:185) | applied: "one milestone at a time" |
| claims | 17 | unverified (0 of 0) | note | docs/glossary.md:146; docs/plan/README.md:89 (row 3) | applied: "Size class (of LAYUP's plan)"; row 3 is large |
| claims | 18 | unverified (0 of 0) | note | docs/glossary.md:144 | applied |
| claims | 19 | unverified (0 of 0) | note | docs/glossary.md:147 | applied |
| claims | 20 | unverified (0 of 0) | note | docs/plan/README.md:149, :237 | applied (as coverage 13) |
| claims | 21 | unverified (0 of 0) | note | docs/plan/README.md:248-249 | applied: the known limit cites R12 |
| claims | 22 | unverified (0 of 0) | note | docs/plan/README.md:44 and :121-135 | applied: rows REQ-015 to REQ-018 and M1 |
| claims | 23 | unverified (0 of 0) | note | docs/plan/README.md:105 | applied |
| claims | 24 | unverified (0 of 0) | note | docs/prd/PRD-0001-layup.md:242 | applied: the NFR-007 Test cell says planned |
| claims | 25 | unverified (0 of 0) | note | docs/tasks/backlog.md:34-35; docs/plan/README.md:87-106 (Issue column); docs/tests/traceability-template.md:64-66 | applied in plan step 7: the issues, their numbers and the backlog lines |
| decisions | 1 | unverified (0 of 0) | material | docs/plan/README.md:199 (K14); :95 (row 9); :98 (row 12); :203 (K18); :129 (NFR-001 row of "What phase 1 proves") | applied (as coverage 4): the first reader of each text writes the reading |
| decisions | 2 | unverified (0 of 0) | material | docs/plan/README.md:101 (row 15); see also :99 (row 13) | applied (as order 9) |
| decisions | 3 | unverified (0 of 0) | material | docs/prd/PRD-0001-layup.md:143-146 (§8); docs/glossary.md:147; docs/plan/README.md:207 (K22) | applied in part: the glossary names PRD §8 and F-0004#11; PRD §8 stays as it is (it quotes the idea owner's fact) |
| decisions | 4 | unverified (0 of 0) | material | docs/plan/README.md:114-135 ("What phase 1 proves") | applied (as claims 22, with K42) |
| decisions | 5 | unverified (0 of 0) | material | docs/plan/README.md:197 (K12) | applied (as coverage 3) |
| decisions | 6 | unverified (0 of 0) | material | docs/plan/README.md:210 (K25), :108, :76-77, :152, :187 (K2); the After cells of rows 19 and 20 at :105-106 | applied: K25 the plan's choice; row 20 after 19; the ADR reads the numbers |
| decisions | 7 | unverified (0 of 0) | note | docs/plan/README.md:87-89 (rows 1 to 3), :92 (row 6), :101-102 (rows 15, 16) | declined: O-121 gives each child the goal count of one; the plan review of each row decides, and the Operator decides a reject on the count alone (Bootstrap mode rule 2) |
| decisions | 8 | unverified (0 of 0) | note | docs/plan/README.md:99 (row 13) and :101 (row 15); also :91 (row 5) and :93 (row 7) | declined: rows 13 and 15 stay whole; each plan review sets its budget, and a task that grows past it splits (R12) |
| decisions | 9 | unverified (0 of 0) | note | docs/plan/README.md:82-83 and the Cap column (:87-106) | applied |
| decisions | 10 | unverified (0 of 0) | note | docs/plan/README.md:248-249 (Known limits) | applied |
| decisions | 11 | unverified (0 of 0) | note | docs/tasks/T-55n2.md:18-19 | applied: the self-check stopped after round 1 (O-125) |
| decisions | 12 | unverified (0 of 0) | note | docs/glossary.md:147 | applied (as claims 19) |
| decisions | 13 | unverified (0 of 0) | note | README.md:58; AGENTS.md:185; docs/plan/README.md:60-61 | applied (as claims 16) |
| decisions | 14 | unverified (0 of 0) | note | docs/plan/README.md:72 (T-b97r); the Issue column at :87-106 | applied (as claims 25) |
| decisions | 15 | unverified (0 of 0) | note | docs/plan/README.md:237 and :149 (#24); :235 and :148 (#15) | applied: the #24 sentence now; #15 at close-out |
| decisions | 16 | unverified (0 of 0) | note | docs/prd/PRD-0001-layup.md:218, :219, :221, :224, :238, :240 (§12 Task cells) | applied (as claims 5) |
| decisions | 17 | unverified (0 of 0) | note | docs/glossary.md:146 | applied (as claims 17) |

## The findings in full

### order 1 — survives

- **Where:** docs/plan/README.md:96 (task table, row 10, After column)
- **Problem:** Row 10 holds `verify-sources`, and that item needs `verify-facts`, which is in row 12. The predecessors of row 10 are 4 and 7 only. With the rows they need, that gives 1, 2, 3, 4 and 7. Row 12 is not in that set. Check `sources` must find each `fact` ref in `docs/facts/`, and the F-NNNN#n resolver for this comes from row 12. So row 10 cannot be done in the order given, or it must write a second resolver.
- **Fix proposed:** Add 12 to the After cell of row 10 ("4, 7, 12"). This makes no cycle, because row 12 comes after row 7 only. Or move `verify-sources` to row 12 or to a later row, and keep `markers`, `discipline-tests` and `link-lint` in row 10.
- Refuter: C(material): I could not refute the finding. Each part of it is correct, and the defect is material.  1. **What the cited text says.** docs/plan/README.md:96 (row 10, `T-8vpw`) holds the item `verify-sources`, and its After cell is "4, 7". Row 12 is at docs/plan/README.md:98. Row 12 holds `verify-facts`, and its
- Refuter: C(material): I tried to refute the finding and could not. The cited text says what the finding claims, and the problem is real.  What the text says: - docs/plan/README.md:96, row 10 (`T-8vpw`) holds `verify-sources` and has the After cell "4, 7". With the rows that 4 and 7 need, its predecessors are 1, 2, 3, 4 a
- Refuter: C(material): I tried to refute the finding and could not. Each cited text says what the finding claims.  1. docs/plan/README.md:96 (row 10, `T-8vpw`) holds `verify-sources`, and its After cell is "4, 7". The plan gets the meaning of After from architecture.md §8 (README.md:6-8; architecture.md:744 says "A task w

### order 2 — survives

- **Where:** docs/plan/README.md:202 (defect register K17) and :95 (row 9, After column)
- **Problem:** The register says that row 9 reads the text that row 10 settles. The After cell of row 9 is "4, 6, 8, 12". With the rows they need, that gives 1, 2, 3, 4, 6, 7, 8 and 12. Row 10 is not in that set. Under O-124, row 9 (S04) writes the raw fact record of the `S01-` and `Q-` answers into `docs/facts/`, and that is the record that K17 is about. So row 9 can merge before row 10 decides how a marker in `docs/facts/` is written or exempted.
- **Fix proposed:** Add 10 to the After cell of row 9 ("4, 6, 8, 10, 12"). This makes no cycle, because row 10 comes after rows 4 and 7 only. Or give K17 to a settler that row 9 already comes after.
- Refuter: C(material): I could not refute it. The cited text says what the finding claims, and the problem is real.  1. The text matches the finding:    - docs/plan/README.md:202 reads "\| K17 \| row 10 \| rows 9, 13 \| Markers inside `docs/facts/` of a target. \|".    - The After cell of row 9 at docs/plan/README.md:95 is "4,
- Refuter: C(material): I could not refute the finding. The cited text says what the finding claims, and the problem is real.  1. docs/plan/README.md:202 says exactly "\| K17 \| row 10 \| rows 9, 13 \| Markers inside `docs/facts/` of a target. \|". Line 95 (row 9, T-7s0y) gives After = "4, 6, 8, 12".  2. I worked out the full s
- Refuter: C(material): Confirmed. docs/plan/README.md:202 says "\| K17 \| row 10 \| rows 9, 13 \| Markers inside `docs/facts/` of a target. \|". The register intro at :181-182 says that the "Read by" tasks are "the tasks that read the settled text". At :95, the After cell of row 9 is "4, 6, 8, 12".  I computed the closure of t

### order 3 — survives

- **Where:** docs/plan/README.md:204 (defect register K19) and :93 (row 7, After column)
- **Problem:** The register says that row 7 reads K19, which row 5 settles. The After cell of row 7 is "1, 2, 3" only, so row 7 does not come after row 5. Rows 5 and 7 can run at the same time, and row 7 can merge before row 5 changes the wording of NFR-004 item 3. NFR-004 applies to `layup setup verify` too.
- **Fix proposed:** Add 5 to the After cell of row 7. Or remove row 7 from the "Read by" cell of K19, because the K19 sources are the result table of row 5 and the `gate:<kind>` rows of row 15. The second fix agrees with the fix of the NFR-004 finding below.
- Refuter: C(material): CONFIRMED, material. The finding quotes the cited text correctly, and the problem is real.  Evidence: - docs/plan/README.md:204 says: "\| K19 \| row 5 \| rows 7, 15 \| A pending kind with a missing tool: the wording of `NFR-004` item 3. \|" - :181-182 define the "Read by" column as "the tasks that read t
- Refuter: C(material): I could not refute the finding. All the cited text agrees with it.  1. docs/plan/README.md:204 says: "\| K19 \| row 5 \| rows 7, 15 \| A pending kind with a missing tool: the wording of `NFR-004` item 3. \|". README.md:181-182 and the new glossary row (docs/glossary.md:148) define "Read by" as "the tasks
- Refuter: C(material): I could not refute the finding. The quotes are exact. docs/plan/README.md:204 reads "\| K19 \| row 5 \| rows 7, 15 \| A pending kind with a missing tool: the wording of `NFR-004` item 3. \|", and the After cell of row 7 at :93 is "1, 2, 3". The predecessors of row 7 are rows 1, 2 and 3 only, because row 

### order 4 — survives

- **Where:** docs/plan/README.md:105-108 (rows 19 and 20, After column; the sentence after the table); docs/glossary.md:147
- **Problem:** Rows 19 (release review) and 20 (first pilot) both have "16, 17, 18" as their predecessors, and no edge connects them. So the release review can finish after the pilot. But the plan says that the pilot is the last child of #29, and the new glossary row says "It is the last task of phase 1." Also, the plan dropped the inventory edge gov-release-review → gov-first-pilot (the review after the pilot) and gives no reason. O-121 as quoted only makes the two tasks "two more children of #29" and gives no order.
- **Fix proposed:** Add 19 to the After cell of row 20, so that the order agrees with line 108 and glossary.md:147. Then say that the review covers the code that the pilot runs. Or keep the order of the inventory (row 19 After 20), and change line 108 and the glossary row.
- Refuter: C(material): I could not refute this finding. Every part of it is true.  1. The plan's task table does not put row 19 before row 20. In docs/plan/README.md:105, row 19 (T-efmy, the release review) has the After cell "16, 17, 18". In docs/plan/README.md:106, row 20 (T-evad, the first pilot) has the same After cel
- Refuter: C(material): Confirmed. I read each cited line.  The plan text: - docs/plan/README.md:105 is row 19 (T-efmy, the release review) and :106 is row 20 (T-evad, the first pilot). Both have the After cell "16, 17, 18". No cell connects the two rows. - :108 says "The first pilot is the last child of #29 (O-121)." - do
- Refuter: C(material): Confirmed. The cited text says what the finding claims.  In /Users/farzam/projects/layup/.worktree/T-55n2/docs/plan/README.md, row 19 at line 105 (T-efmy, the release review) and row 20 at line 106 (T-evad, the first pilot) both have After "16, 17, 18". No cell connects the two rows. Line 108 says "

### order 5 — survives

- **Where:** docs/plan/README.md:92 (row 6, Items and After) and :20
- **Problem:** `psb-batch-api` is the seam where internal/cli gives the gap table to internal/setup. For one of the two possible designs (a type of internal/setup), the inventory makes setup-runner (row 8) its predecessor, and it tells the plan to choose the design. The plan does not choose. It leaves the design to each task. But row 6 does not come after row 8: it comes after rows 1, 2 and 3 only. If the task of row 6 chooses the internal/setup type, it cannot be done in the order given. With both designs, no internal/setup exists at row 6 to receive the table.
- **Fix proposed:** Decide in the plan that internal/cli gives the `psb-gaps` bytes and that internal/setup reads them with internal/tsv in row 9. Then row 6 needs no internal/setup, and row 6 must say this. Or move `psb-batch-api` to row 9, which already comes after rows 6 and 8.
- Refuter: C(material): I could not refute the finding. Its quotes are exact, and the problem is real. It is MATERIAL: an operative ambiguity in the order of the plan.  What I checked:  1. The cited text says what the finding says.    - docs/plan/README.md:92: row 6 holds `psb-batch-api`, and its After cell is "1, 3". Row 
- Refuter: C(material): Confirmed. The cited text says what the finding claims, and the problem is real.  What the text says: - In docs/plan/README.md:92, row 6 holds `psb-batch-api`. Its After cell is "1, 3", and row 3 comes after 1 and 2. So rows 1, 2 and 3 are its only predecessors. `setup-runner` is in row 8 (:94), and
- Refuter: C(material): Confirmed. The cited text says what the finding claims, and the problem is real. (1) docs/plan/README.md:92 puts `psb-batch-api` in row 6 with After "1, 3". Row 8 (docs/plan/README.md:94) holds `setup-runner` with After "3, 7". No chain of After cells puts row 6 after row 8. I compared each item edg

### order 6 — survives

- **Where:** docs/plan/README.md:94 (row 8, Tests column)
- **Problem:** An e2e test runs the built binary. A stop (exit 3) needs a step that asks for input. Row 8 holds no step. S01 is the first step that stops (for missing answers), and it is in row 9, which comes after row 8. So the e2e test of row 8 cannot be written in row 8.
- **Fix proposed:** Move "e2e (a stop and a resume)" to row 9: S01 stops, and the next run resumes at S02. Give row 8 the e2e test of the inventory (the exit 2 cases), and keep the resume test at the integration level with a stub step.
- Refuter: C(material): I could not refute the finding. The quote is exact. docs/plan/README.md:94 (row 8, T-79y7) has the Tests cell "unit; integration (resume from `out/record.tsv`; exit codes 0, 1, 2, 3); e2e (a stop and a resume)". Its items are only `setup-runner`, `setup-answers` and `setup-commands-file`, and it run
- Refuter: C(material): I could not refute the finding. The cited text says what the finding claims, and the problem is real.  1. The cited cell is correct. At HEAD, with a clean tree, docs/plan/README.md:94 (row 8, T-79y7, Tests) is: "unit; integration (resume from `out/record.tsv`; exit codes 0, 1, 2, 3); e2e (a stop and
- Refuter: C(material): I could not refute the finding. The cited text says what the finding says, and the problem is real.  1. docs/plan/README.md:94 (row 8, T-79y7) has this Tests cell: "unit; integration (resume from `out/record.tsv`; exit codes 0, 1, 2, 3); e2e (a stop and a resume)". The plan says that this column hol

### order 7 — survives

- **Where:** docs/plan/README.md:132 (What phase 1 proves, row NFR-004); docs/tests/traceability.md:32
- **Problem:** Row 7 holds the verify frame and the checks `kit-history`, `pin` and `identity`. The specification gives none of them a way to "not run": `not-active` is for a missing tool or a missing script. The verify fixtures that prove NFR-004 are in row 10 (an absent baseline script gives `not-active`) and in row 15 (a `gate:<kind>` with a missing tool gives `not-active`). Rows 10 and 15 both come after row 7, and both list NFR-004 in their own Requirements cells. So the test `T-6x75/e2e/verify-not-active`, which the traceability table gives to row 7, has no fixture in row 7. Also, the acceptance table names the wrong row for the proof.
- **Fix proposed:** Change "(rows 5 and 7)" to "(rows 5, 10 and 15)". Move `T-6x75/e2e/verify-not-active` to row 10 (`T-8vpw`) or row 15 (`T-d6q5`). Or name in row 7 the check that gives the not-active fixture, and say how that check fails to run.
- Refuter: C(material): I tried to refute the finding and could not. The cited text says what the finding claims, and the problem is real.  1. What the cited text says. - docs/plan/README.md:132: "\| NFR-004 \| `layup gate` and `layup setup verify` report a check that did not run as not passed, and a fixture proves it (rows 
- Refuter: C(material): CONFIRMED. I read the cited text and the sources.  1. The quotes are correct. docs/plan/README.md:132 says: "\| NFR-004 \| `layup gate` and `layup setup verify` report a check that did not run as not passed, and a fixture proves it (rows 5 and 7) \| — \|". docs/tests/traceability.md:32 says: "\| `T-6x75/
- Refuter: C(material): CONFIRMED, MATERIAL. (1) Row 7 (T-6x75, docs/plan/README.md:93) holds verify-frame, verify-harness, verify-test-baseline, verify-kit-history, verify-pin, verify-identity. docs/spec/setup.md:235-237 gives kit-history, pin and identity the rules of the sh functions of the same name. check_pin, check_k

### order 8 — survives

- **Where:** docs/plan/README.md:89 (row 3, After column)
- **Problem:** No item of row 3 needs row 1. The only predecessor of `found-cli` is `found-e2e-harness`, and `found-e2e-harness` has no predecessor. internal/cli may not import internal/tsv, and no register row goes from row 1 to row 3 (K32 goes to row 6). Row 3 needs row 2 for K8. With this edge, row 3, and therefore rows 5, 6 and 7, waits for row 1, the largest foundation task, although row 3 could run at the same time as row 1.
- **Fix proposed:** Change the After cell of row 3 to "2", or write down the reason for the edge to row 1.
- Refuter: C(note): The finding is correct. Its severity is note.  What the text says. In docs/plan/README.md:89, row 3 (`T-2yw7`) has the After cell "1, 2". In the inventory, `found-cli` has one predecessor, `found-e2e-harness` (runs/T-55n2/inventory.md:1158). `found-e2e-harness` has no predecessor (":1161 ... \| small
- Refuter: C(note): I could not refute the finding. The cited text and the sources say what the finding says.  1. The cell is correct as quoted. docs/plan/README.md:89, row 3 (`T-2yw7`), has After = "1, 2".  2. No item of row 3 needs row 1:    - runs/T-55n2/inventory.md:1158 gives `found-cli` the After "found-e2e-harne
- Refuter: C(note): I confirm the finding. It is a note.  The cited text: docs/plan/README.md:89 gives row 3 (`T-2yw7`) the After cell `1, 2`. I found no source that supports the edge from row 1 to row 3: - runs/T-55n2/inventory.md:1158: `found-cli` has After `found-e2e-harness`. Line 1161: `found-e2e-harness` has Afte

### order 9 — survives

- **Where:** docs/plan/README.md:101 (row 15, Task column)
- **Problem:** The range S12 to S15 includes S14. But S14 is part of the prose step, which O-123 runs before S10, and its item `setup-prose-inputs` ("S07, S08, S09, S14") is in row 13. A reader can think that S14 is work of row 15, after S13.
- **Fix proposed:** In row 15, write "Steps S12, S13 and S15", and name S14 in the Task cell of row 13.
- Refuter: C(note): I could not refute the finding. It is a real problem, but only a note.  What the text says: - docs/plan/README.md:101, row 15 (T-d6q5), Task cell: "Steps S12 to S15, and the checks `jobs` and `gate:<kind>`". The quote is exact. - docs/tasks/T-55n2.md:46 (O-123): "the prose step (S07 to S09 and S14, 
- Refuter: C(note): I confirm this finding as a NOTE. I could not refute it.  1. The cited text is correct. docs/plan/README.md:101 (row 15, Task cell) says "Steps S12 to S15, and the checks `jobs` and `gate:<kind>`".    - If you read the range by step ID, it includes S14.    - The step table of the specification still
- Refuter: C(note): The finding is real, and it is a note.  **The quote is correct.** docs/plan/README.md:101, the Task cell of row 15, says: "Steps S12 to S15, and the checks `jobs` and `gate:<kind>`".  **The problem is real.** - In the step table, S14 is between S13 and S15 (docs/spec/setup.md:97-100). - O-123 writes

### order 10 — refuted

- **Where:** docs/plan/README.md:196 (defect register K11)
- **Problem:** `found-e2e-harness` (row 3) is one of the sources of K11. Its helpers "make a git repository from a fixture directory under root tests/", and K11 settles against that checked-in form. Row 3 comes before row 7, the settler, and the register does not name row 3 as a reader. So row 3 can check in a fixture tree that LAYUP's own checks `markers` or `adapted`, or link-lint, then fail.
- **Fix proposed:** Add row 3 to the "Read by" cell and make row 3 the settler, because row 3 is the first row that makes a fixture repository. Or say in the K11 note that row 3 checks in no baseline tree under tests/.
- Refuter: R(none): Refuted. The quote is correct: docs/plan/README.md:196 says "\| K11 \| row 7 \| rows 10, 11, 12, 16 \| The stand-in baseline is built at test time; a marker in Go source is written as an escape. \|". `found-e2e-harness` is one of the sources of K11 (runs/T-55n2/inventory.md:2267). But the problem does no
- Refuter: R(none): Not confirmed. The quote is correct: docs/plan/README.md:196 is "\| K11 \| row 7 \| rows 10, 11, 12, 16 \| The stand-in baseline is built at test time; a marker in Go source is written as an escape. \|". found-e2e-harness (row 3) is a source of K11 (runs/T-55n2/inventory.md:2267), and row 3 comes before 
- Refuter: R(none): Not confirmed. The quoted facts are correct, but the problem is not real.  What is correct: docs/plan/README.md:196 reads "\| K11 \| row 7 \| rows 10, 11, 12, 16 \| The stand-in baseline is built at test time; a marker in Go source is written as an escape. \|". Row 3 (line 89) holds found-e2e-harness. K1

### order 11 — refuted

- **Where:** docs/plan/README.md:209 (defect register K24)
- **Problem:** The note makes row 14 act on K24 first: row 14 writes the catalog form of the floor, and that form cannot be a marker under internal/catalog/, because MK_EXEMPT does not exempt that path. But the settler, row 15, comes after row 14. The register does not say which row decides the catalog form.
- **Fix proposed:** Write "rows 14, 15" in the "Settled by" cell (row 14: the catalog form; row 15: the write at S12), and "rows 15, 16" in the "Read by" cell.
- Refuter: R(none): Refuted. The finding's main claim is false. It says "The register does not say which row decides the catalog form". But the same register row, docs/plan/README.md:209, says: "\| K24 \| row 15 \| row 16 \| S12 writes the coverage floor as an open gap (L-B2); row 14 gives its catalog form. \|". The note gi
- Refuter: R(none): Not confirmed. The finding says: "The register does not say which row decides the catalog form." The cited cell says the opposite. docs/plan/README.md:209, Note: "S12 writes the coverage floor as an open gap (L-B2); row 14 gives its catalog form." The register gives the catalog form to row 14 and th
- Refuter: R(none): Refuted. Some parts of the finding are correct. The quote at docs/plan/README.md:209 is exact. Row 15 has After "13, 14" (README.md:101). MK_EXEMPT (docs/setup/setup-check.sh:305) does not exempt internal/catalog/.  The central claim is false: "The register does not say which row decides the catalog

### order 12 — survives

- **Where:** docs/plan/README.md:87 (row 1) and :193 (defect register K8)
- **Problem:** Row 1 holds `found-spec-schema-test`, and the home of its registry is an open question. Option B (one test file in internal/cli) is possible only "if test files are exempt from the import table", and that is K8, which row 2 settles. Rows 1 and 2 have no predecessors, so row 1 can choose the home before K8 is settled.
- **Fix proposed:** Say that row 1 uses option A (a test in each owner package and one completeness test over the block names). Or add row 1 to the readers of K8 and put row 1 after row 2.
- Refuter: C(note): I could not refute the finding: the cited text is correct and the gap is real. But it is a note, not a material finding.  What the text says: - docs/plan/README.md:193 reads "\| K8 \| row 2 \| rows 3, 7, 16 \| Whether the package rules bind the test files. \|". Row 1 is not listed as a reader. - Row 1 (:
- Refuter: C(note): I could not refute the finding. The cited text says what the finding claims, and the gap is real. It is a note, not a material finding.  Facts I checked at HEAD 71c1bfb: - docs/plan/README.md:87 is row 1, `T-18v6`. Its Items cell is "`found-tsv`, `found-spec-schema-test`" and its After cell is "—". 
- Refuter: C(note): Confirmed as a note. I checked each claim at HEAD 71c1bfb. The work tree is clean.  1. docs/plan/README.md:87 gives row 1 the items `found-tsv` and `found-spec-schema-test`, with After "—". Line 88 (row 2) also has After "—". No edge puts row 1 after row 2. Line 5 says the table gives "the tasks of 

### order 13 — survives

- **Where:** docs/plan/README.md:129 (What phase 1 proves, row NFR-001)
- **Problem:** Under O-124, S04 (row 9) writes the raw fact record of the `S01-` and `Q-` answers in the first commit on the setup branch. The REQ-002 cell (line 124) names rows 9 and 13 for the answers. The NFR-001 cell does not name row 9.
- **Fix proposed:** Write "(rows 9, 13, 15, 16 and 20; O-115, O-124)".
- Refuter: C(note): I tried to refute this finding and could not. The finding is real, and it is a note.  The cited text, docs/plan/README.md:129, says exactly: "the setup record, the verify table, the rule-path register and the answers in the target's Git (rows 13, 15, 16 and 20; O-115)". The list does not name row 9.
- Refuter: C(note): Confirmed. The cited text says what the finding says, and the problem is real.  Text: docs/plan/README.md:129, NFR-001 row, cell 2: "the setup record, the verify table, the rule-path register and the answers in the target's Git (rows 13, 15, 16 and 20; O-115)". O-124 (docs/tasks/T-55n2.md:51-55): "S
- Refuter: C(note): I confirm the finding. It is a note.  **The quote is exact.** docs/plan/README.md:129 says: "the setup record, the verify table, the rule-path register and the answers in the target's Git (rows 13, 15, 16 and 20; O-115)".  **What the sources say:** - O-124 is at docs/tasks/T-55n2.md:51-55. It says: 

### order 14 — survives

- **Where:** docs/plan/README.md:4-5 and :186 (defect register K1)
- **Problem:** Not all row numbers follow the order that the After column sets. Row 9 needs row 12, which the table lists after it. After the K17 fix and the row 10 fix above, row 9 also needs rows 10 and 12, and row 10 needs row 12. A reader who does the tasks in row order starts S01 to S04 before check `facts` exists. K1 ("The order of the task table") makes the row order look like the build order.
- **Fix proposed:** In "How to read this plan", say that the After column gives the order and that the row number is only a handle. Or renumber the rows so that each row comes after all its predecessors, and update the row references in the register, PRD-0001 §12 and the traceability table.
- Refuter: C(note): I could not refute the facts in the finding. I rate the problem a note, not a material finding.  What I verified: - docs/plan/README.md:4-5 says "the tasks of phase 1 in their order". - docs/plan/README.md:186 is "\| K1 \| `T-55n2` \| all rows \| The order of the task table. \|". - docs/plan/README.md:95
- Refuter: C(note): The finding is correct. It is a NOTE, not a material finding.  What the cited text says (checked at HEAD 71c1bfb): - docs/plan/README.md:4-5: "the tasks of phase 1 in their order". README.md:59 and AGENTS.md:185 repeat "in their order". - docs/plan/README.md:186: "\| K1 \| `T-55n2` \| all rows \| The or
- Refuter: C(note): Confirmed as a note.  The cited text is correct: - docs/plan/README.md:4-5 says "the tasks of phase 1 in their order". - :186 says "\| K1 \| `T-55n2` \| all rows \| The order of the task table. \|". - :95 gives row 9 (`T-7s0y`, steps S01 to S04) the After cell "4, 6, 8, 12".  I ran a read-only script ove

### order 15 — survives

- **Where:** runs/T-55n2/plan-check.py:158-187
- **Problem:** The only order check in the script tests that each predecessor is a row and that the rows make no cycle. It does not test two things: that each inventory edge of an item (and each edge that K4 adds) is among the rows that its row needs, and that each "Settled by" row of the register comes before its "Read by" rows. These two tests find the row 6, row 10, row 19, K17 and K19 defects above. Now the script gives PASS on all checks and exit 0 with these defects present.
- **Fix proposed:** Add the two checks to plan-check.py and run it again as evidence of the plan.
- Refuter: C(note): I confirm this finding as a note.  (1) The cited text is correct. runs/T-55n2/plan-check.py:158-187 is check 3, "# 3. predecessors resolve, with no cycle". It tests only three things: each After number is a row, no row names itself, and the rows make no cycle. No other check in the script tests the 
- Refuter: C(note): I could not refute the finding. Each part of it is correct. It is a note, not material.  1. The cited text is correct. runs/T-55n2/plan-check.py:158 says "# 3. predecessors resolve, with no cycle". Lines 159-187 check only three things: each number in "After" is a row, no row needs itself, and the g
- Refuter: C(note): Confirmed as a NOTE.  (1) The cited text says what the finding claims. runs/T-55n2/plan-check.py:158 has "# 3. predecessors resolve, with no cycle". Lines 159-187 test only three things: that each After number is a row, that a row does not depend on itself, and that there is no cycle. No other check

### coverage 1 — survives

- **Where:** docs/plan/README.md:96 (row 10, Items and After cells); docs/plan/README.md:98 (row 12, After cell `7`)
- **Problem:** Row 10 hosts `verify-sources`. Check `sources` must resolve each `fact` ref with the F-NNNN#n resolver of `verify-facts`, and row 12 hosts `verify-facts`. Row 10's predecessors are 4 and 7, not 12. Row 12's only predecessor is 7. So row 10 can land before row 12, with no resolver for its `fact` rule. The task cannot be done in the order given. The plan drops this inventory edge, and K4 says that "The predecessors of the task table hold these edges".
- **Fix proposed:** Set row 10's After to "4, 7, 12". This makes no cycle, because row 12 needs only row 7. Another fix is to move `verify-sources` into row 12. Row 13 already comes after both rows.
- Refuter: C(material): I could not refute the finding. The cited text says what the finding claims, and the problem is real.  1. The quote is correct. docs/plan/README.md:96 (row 10, T-8vpw) has the Items `verify-markers`, `verify-sources`, `verify-baseline-scripts`, `gov-issue-21` and After "4, 7". docs/plan/README.md:98
- Refuter: C(material): I could not refute the finding. The problem is real and material.  What I checked: 1. Cited text is correct. The worktree is clean at HEAD 71c1bfb. docs/plan/README.md:96 is row 10 (`T-8vpw`). It holds `verify-markers`, `verify-sources`, `verify-baseline-scripts` and `gov-issue-21`, and its After ce
- Refuter: C(material): Confirmed. I read the cited text and the sources myself. I could not refute the finding.  1. The plan rows. docs/plan/README.md:96 is row 10 (`T-8vpw`). It holds `verify-sources`, and its After cell is `4, 7`. docs/plan/README.md:98 is row 12 (`T-9t1q`). It holds `verify-facts`, and its After cell i

### coverage 2 — survives

- **Where:** docs/plan/README.md:92 (row 6, Items and After cells)
- **Problem:** Row 6 hosts `psb-batch-api`. That item is the hand-off inside `layup setup`: internal/cli reads the brief once and gives the bytes and the gap table to internal/setup. internal/setup and the `setup` dispatch in internal/cli first exist in row 8 (`setup-runner`). Row 6's predecessors are only 1 and 3, so row 6 can land before row 8. Row 6 then cannot build or test the hand-off: the item's integration test compares the IDs "that cli hands on". If the seam is "a type of internal/setup", row 6 cannot build it at all. The plan says that it fixes order and not design, but this order is correct for only one of the two designs that the inventory leaves open. By O-124 the table is also needed at S04, not at S06.
- **Fix proposed:** Move `psb-batch-api` to row 9, with `setup-s01-questions`, which holds the same seam and test. Another fix is to add 8 to row 6's After. If row 6 keeps the item, record the seam (table bytes) in the plan. Change "S01 and S06" to "S01 and S04" in the scope of the item.
- Refuter: C(note): One part of the finding is real. The rest does not hold, so this is a note, not a material finding.  What is real: docs/plan/README.md:92 puts `psb-batch-api` in row 6 with After `1, 3`. Row 8 (`setup-runner`, :94) is the first row that makes internal/setup and the `layup setup WORK` dispatch. The i
- Refuter: C(material): The core of the finding is real. Two of its parts claim too much.  What is real: - docs/plan/README.md:92 puts `psb-batch-api` in row 6 with After "1, 3". - The item says that cli "hands the gap table to internal/setup" (runs/T-55n2/inventory.md:64). It also says "The plan chooses one" of two seam d
- Refuter: C(note): I confirm this finding only as a NOTE. The claim that it is material is not correct.  The finding quotes the plan correctly. docs/plan/README.md:92 puts `psb-batch-api` in row 6 with After "1, 3". Row 8 (`setup-runner`, :94) is the first row that makes internal/setup and the `setup` dispatch. The in

### coverage 3 — survives

- **Where:** docs/plan/README.md:197 (defect register, K12)
- **Problem:** O-123 settles only the `adapted` half of K12. The `kit-history` half is the 7 + 32 lines of backlog.md and completed.md that link the baseline repository and name no deleted T-ID. The register settles it with a new S05 rule, but neither O-123 nor any other recorded decision states that rule. The register gives O-124's authority to a value that the plan sets itself. The implementer of row 13 will then write the S05 rule into setup.md:90 with O-123 as its source, and not as "decided here" with a reason.
- **Fix proposed:** Split the Settled-by cell: "O-123 (check adapted); row 13 (S05: each line of backlog.md and completed.md that links the baseline's repository, decided here)". Another fix is to get an Operator decision for the S05 rule and cite it.
- Refuter: C(material): I could not refute the finding. The cited text says what the finding claims, and the problem is real.  1. The cited row. /Users/farzam/projects/layup/.worktree/T-55n2/docs/plan/README.md:197 is "\| K12 \| O-123 \| rows 11, 13, 20 \| The prose step stops for every flagged file; S05 also removes each line
- Refuter: C(material): I could not refute the finding. The problem is real and it is material.  What I checked: 1. docs/plan/README.md:197 says: "\| K12 \| O-123 \| rows 11, 13, 20 \| The prose step stops for every flagged file; S05 also removes each line that links the baseline's repository (row 13). \|". The Settled-by cell 
- Refuter: C(material): I could not refute the finding. The cited text says what the finding claims, and the problem is real.  1. docs/plan/README.md:197 says: "\| K12 \| O-123 \| rows 11, 13, 20 \| The prose step stops for every flagged file; S05 also removes each line that links the baseline's repository (row 13). \|". Lines 

### coverage 4 — survives

- **Where:** docs/plan/README.md:199 (K14); :98 (row 12 After `7`); :95 (row 9 After `4, 6, 8, 12`); :203 (K18)
- **Problem:** The reading of O-124 changes the rule that row 12 implements. Check `facts` (setup.md:246) names "the `S01-` and `Q-` IDs in the record of S06", and row 12 also settles K18, which says which records exist at each step. By O-124, S04 writes the first answers record, before S06 copies the briefs. But row 12 comes before row 9, and row 12 is not a reader of K14. So row 12 builds the record placement, its fixtures ("an M- ID in the S06 record", "a correct S06-only target") and the per-step rule from the text that O-124 replaces. The note also leaves out two more places that state the old reading: records.md:105 ("two raw fact records (S06, S11)") and psb-check.md:69 ("at setup step S06").
- **Fix proposed:** Add row 12 to the "Read by" cell of K14. Make row 12 write the reading of O-124, because it is the first task that implements text that the reading changes. Another fix is to write the reading now in T-55n2. In both cases, write it into setup.md:91, :96 and :246, records.md:105 and psb-check.md:68-71, not only into `spec/setup.md`.
- Refuter: C(material): I could not refute the finding. Each part of the cited text says what the finding claims.  1. O-124 moves the record. docs/tasks/T-55n2.md:52-55 says: "S04 also writes the raw fact record of the `S01-` and `Q-` answers, with its line in `facts.sha256` and its index row, in the first commit on the se
- Refuter: C(material): I could not refute the finding. Every citation is correct.  1. **The plan text.** docs/plan/README.md:199 says "\| K14 \| O-124 \| rows 9, 13 \| Row 9 writes the reading into `spec/setup.md`. \|". Row 12 (:98) has After `7`. Row 9 (:95) has After `4, 6, 8, 12`, so row 12 is done before row 9. K18 (:203) 
- Refuter: C(material): I could not refute the finding. Each cited text says what the finding claims. I checked the text at HEAD 71c1bfb, and the worktree has no uncommitted changes.  What the cited text says: - docs/plan/README.md:199 says: "\| K14 \| O-124 \| rows 9, 13 \| Row 9 writes the reading into `spec/setup.md`. \|". R

### coverage 5 — survives

- **Where:** docs/plan/README.md:101 (row 15, Task cell); :99 (row 13); :201 (K16)
- **Problem:** By O-123, S14 is part of the prose step that runs before S10, and row 13 hosts it (`setup-prose-inputs`: S07, S08, S09, S14). But the name of row 15 still claims S12 to S15. Thus S14 has two hosts in the text, and the issue title of row 15 includes S14. O-123 also changes runner behavior that row 8 builds. S14 moves before S10, and one stop table covers S07 to S09, S14 and the flagged kept files ("all in one table"). The setup-runner item "runs the next steps in order S01 to S15", and setup.md:107-108 stops for "every missing input of that step". But only row 13 reads K16, and row 13 comes after row 8.
- **Fix proposed:** Name row 15 "Steps S12, S13 and S15, and the checks `jobs` and `gate:<kind>`". Name row 13 "Steps S05 to S11 and S14 (the prose step S07 to S09 and S14, before S10, O-123) …". Add row 8 to the "Read by" cell of K16, for the step order and the stop that covers more than one step.
- Refuter: C(material): Confirmed. The cited text says what the finding says, and the problem is real.  1) Row 15 claims S14, but row 13 hosts it. docs/plan/README.md:101 names row 15 "Steps S12 to S15, and the checks `jobs` and `gate:<kind>`". Its Items cell has no S14 item: S12 is setup-s12-stack-files, S13 is setup-s13-
- Refuter: C(note): The quotes are correct, and the problem is real. But it is a note, not material.  What I checked: - docs/plan/README.md:101 (row 15, Task cell): "Steps S12 to S15, and the checks `jobs` and `gate:<kind>`". Its Items cell has no S14 item. - docs/plan/README.md:99 (row 13, Task cell): "Steps S05 to S1
- Refuter: C(material): The finding is correct and material for the names of rows 15 and 13. The part about row 8 is only a note.  1. The cited text says what the finding claims. docs/plan/README.md:101 (row 15, Task) reads: "Steps S12 to S15, and the checks `jobs` and `gate:<kind>`". The Items of row 15 do not include S14

### coverage 6 — refuted

- **Where:** docs/plan/README.md:147 (hosts table); :222 (K37)
- **Problem:** T-55n2 does not do the work of this item. The item is the specification work of each later milestone: the docs/spec sections and the packages.md rows. It is large, it comes after the first pilot, and the schema-block test on the new sections tests it. T-55n2 only puts it in the milestone schedule. The plan also contradicts itself: K37 says that this item is "out of phase 1 (the hosts table)", but the hosts table gives it to T-55n2, which is a phase-1 task. The inventory Host column, which the plan fills, now shows this future work as done by T-55n2.
- **Fix proposed:** Host the item at "the specification task of each milestone `M2a` to `M4c`", with the reason "the milestone table schedules it".
- Refuter: R(none): I refute the finding. The cited row exists (docs/plan/README.md:147), but it does not say what the finding claims, and no reader would act wrongly on it.  (1) The row does not claim that T-55n2 writes the later specification. Its Reason cell says the item is hosted through "The milestone table: each
- Refuter: C(note): The finding is real only as a small error in a label. It is a note, not material. Its claim of a contradiction with K37 is wrong.  1. The quotes are correct. docs/plan/README.md:147 says: "\| `gov-milestone-spec` \| `T-55n2` \| The milestone table: each later milestone starts with its own specification
- Refuter: R(none): Refuted. The quotes are accurate, but the problem is not real in any way that a reader would act on.  (1) The K37 contradiction depends on a false premise. The finding says T-55n2 is a phase-1 task. The plan says it is not: - docs/plan/README.md:25-26: "No task of phase 1 starts before #42 closes". 

### coverage 7 — survives

- **Where:** docs/plan/README.md:228-229; rows 3, 13 and 14 (:89, :99, :100); :196 (K11)
- **Problem:** The plan gives each task the open questions of its items. It does not give the duplicate resolutions of the inventory. Some duplicates are in different rows, so the same work is in the scope of two rows, and nothing tells the earlier row to leave it. (a) setup-s10-markers (row 13) builds its own marker scanner, a MK_EXEMPT copy and parity tests. verify-markers (row 10) builds the same scanner, and the inventory keeps only that one and gives S10 the list through internal/cli. (b) The unit test of gate-catalog-go-workflow (row 14) needs a job-name reader. The inventory moves that reader to verify-jobs (row 15), which comes after row 14. (c) found-e2e-harness (row 3) makes fixture repositories "from a fixture directory under root tests/", and verify-test-baseline (row 7) and setup-e2e-offline (row 16) also build the stand-in baseline. K11 (row 7) then decides that the baseline is built at test time, but row 3 comes before row 7 and does not read K11. (d) Both found-e2e-harness (row 3) and psb-cli-contract (row 6) contain the psb golden and repeat e2e scenarios. Thus some work is hosted twice.
- **Fix proposed:** Add "and the duplicate resolutions of its items" to that sentence, or give one line for each row. Row 13 uses the scanner of row 10 through internal/cli. The job-name test of row 14 moves to row 15. Row 3 keeps only the helpers and the version scenario, and it uses no baseline that is kept in Git (add row 3 to the "Read by" cell of K11). Row 6 owns the psb e2e scenarios.
- Refuter: C(material): I could not refute the finding. Each claim agrees with the files.  (1) docs/plan/README.md:228-229, "Each task also takes the open questions of its own items in the inventory as inputs of its plan.", is the only rule on what a task takes from the inventory. No line of the plan, docs/tests/traceabili
- Refuter: C(note): The cited text says what the finding claims. The problem is real, but it is a note, not material.  What is true. docs/plan/README.md:228-229 says: "Each task also takes the open questions of its own items in the inventory as inputs of its plan." The plan does not refer to the "Duplicates" part of th
- Refuter: C(material): Confirmed. The finding is material: an operative omission and an ambiguity of scope that a task author will act on.  **What the plan says.** docs/plan/README.md:228-229 says only: "Each task also takes the open questions of its own items in the inventory as inputs of its plan." The plan also says th

### coverage 8 — refuted

- **Where:** docs/plan/README.md:108; :105-106 (rows 19 and 20, both After `16, 17, 18`); docs/glossary.md:147 ("It is the last task of phase 1")
- **Problem:** The recorded O-121 says only that the release review and the first pilot "are two more children of #29". It gives no order. The inventory puts the release review after the pilot: gov-release-review has After gov-first-pilot, and its layer (22) comes after the pilot's (21). Rows 19 and 20 have the same predecessors, and the plan makes the pilot the last child on the authority of O-121. So the plan drops an inventory edge with no reason, and it gives the order to a decision that does not state it. If the pilot finds defects and they are fixed, row 19 reviewed a different release from the one that phase 1 ships. REQ-015 and REQ-017 say "a code review of each release records it".
- **Fix proposed:** One fix: add 20 to the After cell of row 19, and make the release review the last child. Another fix: keep the order, and write its reason and which release row 19 reviews. In both cases, remove "(O-121)" from that sentence and remove "It is the last task of phase 1" from the glossary row, unless the Operator's decision states it.
- Refuter: R(none): Refuted. The finding quotes the cited text correctly, but the problem it describes is not real.  (1) An Operator decision supports the order. O-122 (docs/tasks/T-55n2.md:38-39) says: "the first pilot is the one of ADR-0012 part 6, the end of phase 1 and of bootstrap mode". O-121 (:35-36) makes the r
- Refuter: C(note): Partly true, but not material. Do not apply the fix that the finding gives.  TRUE: docs/plan/README.md:108 says "The first pilot is the last child of #29 (O-121)." O-121 (docs/tasks/T-55n2.md:35-37) says only "The release review of phase 1 and the first pilot are two more children of #29". It gives 
- Refuter: R(none): REFUTED. The Operator's decision says that the pilot is the last child.  1. O-121 is the Operator's answer "Q-1 A". docs/tasks/T-55n2.md:27-28 says: "The options are those of the questions comment as edited after the plan review." The bullet at :30-37 only summarises that option. I read the text of 

### coverage 9 — survives

- **Where:** docs/plan/README.md:132; docs/tests/traceability.md:32
- **Problem:** For `layup setup verify`, the fixtures with a check that does not run are in other rows. Row 10 has verify-baseline-scripts (a missing script gives not-active and exit 1). Row 15 has verify-gate-fixtures (a missing tool gives `gate:<kind>` not-active). Row 16 has the verify-acceptance e2e test. No item of row 7 has such a fixture, but the traceability handle `T-6x75/e2e/verify-not-active` gives that e2e test to row 7.
- **Fix proposed:** Cite rows 5, 10, 15 and 16, and give the traceability handle to the row that holds the fixture. Another fix is to add the not-active scenario to the Tests cell of row 7.
- Refuter: C(material): CONFIRMED. The change adds both cited lines (git diff 7cdd346 HEAD). docs/plan/README.md:132 says "\| NFR-004 \| `layup gate` and `layup setup verify` report a check that did not run as not passed, and a fixture proves it (rows 5 and 7) \| — \|". docs/tests/traceability.md:32 says "\| `T-6x75/e2e/verify-
- Refuter: C(note): The cited text is exact. docs/plan/README.md:132 says: "\| NFR-004 \| `layup gate` and `layup setup verify` report a check that did not run as not passed, and a fixture proves it (rows 5 and 7) \| — \|". docs/tests/traceability.md:32 gives `T-6x75/e2e/verify-not-active` (e2e, NFR-004, REQ-002) to T-6x75
- Refuter: C(note): Confirmed. The cited text says what the finding claims, and the problem is real.  1) The cited text. docs/plan/README.md:132 says: "\| NFR-004 \| `layup gate` and `layup setup verify` report a check that did not run as not passed, and a fixture proves it (rows 5 and 7) \| — \|". docs/tests/traceability.

### coverage 10 — survives

- **Where:** docs/plan/README.md:210 (K25)
- **Problem:** O-121 makes the ADR a separate issue that the pilot task opens, so the count comes after the pilot. But O-121 does not say where the rule that reads the numbers is written. That half comes from guardrails §1 ("Write the pass/fail numbers first").
- **Fix proposed:** Write "O-121 (the count in the ADR task); guardrails §1 (the rule in row 20's plan review)".
- Refuter: C(note): I tried to refute this finding and could not. It is real, but it is a note, not material.  **What I checked** - **The cited cell.** docs/plan/README.md:210 has the exact text: "\| K25 \| O-121 \| row 20 \| The pilot rule in row 20's plan review; the count in the ADR task. \|". The register header (README
- Refuter: C(note): I could not refute the finding. It is real, and it is a note.  1. The cited cell is word for word at docs/plan/README.md:210: "\| K25 \| O-121 \| row 20 \| The pilot rule in row 20's plan review; the count in the ADR task. \|". Line 181 says that the "Settled by" column names "the task that settles it". 
- Refuter: C(note): I tried to refute the finding and could not. It is real, and it is a note.  **The quote is exact.** docs/plan/README.md:210 reads "\| K25 \| O-121 \| row 20 \| The pilot rule in row 20's plan review; the count in the ADR task. \|".  **O-121 does not cover the rule half.** - O-121 (docs/tasks/T-55n2.md:30

### coverage 11 — survives

- **Where:** docs/plan/README.md:192 (K7), :202 (K17), :209 (K24), :214 (K29), :224 (K39); also K8 (:193) and K31 (:216)
- **Problem:** Some rows implement settled text but are not in the "Read by" cell. K7: rows 10 and 11 list tracked files (`ls-files`) for checks `markers` and `adapted`. K17: row 12 sets the answers-record form, which must not show a raw marker. K24: the note says that row 14 "gives its catalog form", but row 14 is in neither column. K29: rows 10 and 15 print progress for the baseline scripts and the fixture runs. K39: M2a, because the forge adapter imports `net/http` in phase 2, and line 61 also makes K38 to K41 inputs of M2a. K8: row 14, whose fixture test starts git from test code. K31: row 2, which isolates git from the host before row 9.
- **Fix proposed:** Add those rows to the "Read by" cells. Add row 14 to the "Settled by" cell of K24.
- Refuter: C(note): The quoted cell is exact: docs/plan/README.md:192 "\| K7 \| row 2 \| rows 7, 8, 9, 13, 15 \| The verbs of `internal/git` in `spec/packages.md`. \|". Five of the seven parts are real. (1) K7 (:192): inventory.md:2263 needs "ls-files (S10, checks markers and adapted)". inventory.md:842-843 put verify-marke
- Refuter: C(note): Confirmed as a note. The cited cells match the quotes in docs/plan/README.md (:192, :193, :202, :209, :214, :216, :224). Line 182 says that "Read by" names "the tasks that read the settled text". The inventory shows that these cells leave out tasks: - K7 (:192): inventory.md:2263 needs `ls-files` fo
- Refuter: C(note): Confirmed in part. Severity: note. The cited cells say what the finding says: - docs/plan/README.md:192: "\| K7 \| row 2 \| rows 7, 8, 9, 13, 15 \|" - :193: "\| K8 \| row 2 \| rows 3, 7, 16 \|" - :209: K24 is settled by row 15 and read by row 16, with the note "row 14 gives its catalog form". - :214: "\| K29

### coverage 12 — survives

- **Where:** docs/prd/PRD-0001-layup.md:212-214, :218, :226, :230, :238, :240; docs/plan/README.md:46 (M2b)
- **Problem:** The phase-1 rows are not consistent. REQ-009, REQ-011, NFR-001, NFR-002 and NFR-005 name their later milestones. REQ-001 ("T-dq05, T-zmj6, T-5zmw"), REQ-002, REQ-004, REQ-007, NFR-003, NFR-004 and NFR-006 do not, but the column "A later milestone proves" of the plan names them. NFR-003 does not name T-evad (the audit, row 20). NFR-005 does not name T-b97r (the verify verdicts, row 7). M2b does not list REQ-013, which later-p2-sessions-ledger serves, so the REQ-013 cell names only `M4a`.
- **Fix proposed:** Name the later milestones in all phase-1 rows, from "What phase 1 proves", or in none of them. Add T-evad to NFR-003, T-b97r to NFR-005, and REQ-013 to M2b.
- Refuter: C(note): Partly confirmed; it is a NOTE. Two sub-claims are wrong and must not go into the fix.  TRUE: (1) No single rule. docs/prd/PRD-0001-layup.md:212-214 says the plan "fills the Task column with its tasks for phase 1 and its milestones (`M2a` to `M4c`) for the later phases". These phase-1 rows name late
- Refuter: C(note): I confirm this finding as a NOTE, with one correction. I could not refute it.  1. The phase-1 rows follow no single rule. PRD-0001:212-214 says the plan fills the Task column "with its tasks for phase 1 and its milestones (`M2a` to `M4c`) for the later phases". These rows name later milestones:    -
- Refuter: C(note): CONFIRMED, as a NOTE. I read each cited line and the sources. One sub-claim is wrong, but this does not change the result.  (1) Inconsistent phase-1 rows: real. PRD-0001:212-214 says the plan "fills the Task column with its tasks for phase 1 and its milestones (`M2a` to `M4c`) for the later phases".

### coverage 13 — survives

- **Where:** docs/plan/README.md:149, :237
- **Problem:** The item also holds the fix of one sentence: docs/ci/README.md:37 says "Every job restores its check scripts", but nested-checkout-check.sh is not restored, on purpose. The item puts that fix "into the first product task that edits docs/ci/README.md or ci.yml". When #24 closes, no issue and no host keep that fix.
- **Fix proposed:** Name the host, for example row 16 if it changes ci.yml, or a known limit. Another fix is to make the change of one sentence in T-55n2.
- Refuter: C(note): I could not refute the finding. The cited text says what the finding claims, and the problem is real. It is a note, not material.  What the plan says: - docs/plan/README.md:149 says "\| `gov-issue-24` \| `T-55n2` \| Checked and closed with its evidence. \|". - :237 says "Checked and closed by `T-55n2`."
- Refuter: C(note): CONFIRMED. I could not refute this finding. Its severity is NOTE.  (1) The item holds the fix. runs/T-55n2/inventory.md:1637 (gov-issue-24, What) says: "The fix of that sentence is batched into the first product task that edits docs/ci/README.md or ci.yml." At :1642 the item has its own test: "disci
- Refuter: C(note): Confirmed. I could not refute it.  What the cited text says: - docs/plan/README.md:149 says "\| `gov-issue-24` \| `T-55n2` \| Checked and closed with its evidence. \|". Line 237 says "Checked and closed by `T-55n2`." The K36 row at line 221 says only "#24 closed with its evidence". - runs/T-55n2/invento

### coverage 14 — survives

- **Where:** docs/plan/README.md:129
- **Problem:** By O-124, S04 writes the first answers record (the S01- and Q- answers), and S04 is in row 9. This cell names only row 13.
- **Fix proposed:** Write "rows 9, 13, 15, 16 and 20".
- Refuter: C(note): I could not refute the finding. The problem is real, but it is a note, not a material finding.  1. The cited cell is accurate. docs/plan/README.md:129 says: "\| NFR-001 \| the setup record, the verify table, the rule-path register and the answers in the target's Git (rows 13, 15, 16 and 20; O-115) \|".
- Refuter: C(note): I could not refute the finding. The problem is real, but it is a note, not a material finding.  Evidence (worktree clean, HEAD 71c1bfb): - docs/plan/README.md:129 says: "\| NFR-001 \| the setup record, the verify table, the rule-path register and the answers in the target's Git (rows 13, 15, 16 and 20
- Refuter: C(note): I could not refute the finding. It is real, and it is a note.  1. The cited cell, docs/plan/README.md:129, says: "\| NFR-001 \| the setup record, the verify table, the rule-path register and the answers in the target's Git (rows 13, 15, 16 and 20; O-115) \|".  2. O-124, at docs/tasks/T-55n2.md:51-55, s

### coverage 15 — survives

- **Where:** docs/plan/README.md:197 (K12 note) and :106 (row 20 Requirements, REQ-018); docs/prd/PRD-0001-layup.md:235
- **Problem:** By O-123, a person writes each flagged kept file again as an F-<path> input. The flagged files that the inventory measured include rule paths (docs/engineering-discipline.md, docs/issue-workflow.md, docs/ci/README.md, docs/tests/README.md). The setup record has no source kind for an input file (answer\|catalog\|fact\|computed\|gap\|step). So NFR-003 item 1 and REQ-018 ("except the adapted values the setup records with evidence") have no evidence for these files, and row 20 claims REQ-018. No register row records this.
- **Fix proposed:** Add a register row, hosted by row 13, for the record evidence of the F-<path> input files (prose files and flagged kept files).
- Refuter: C(note): CONFIRMED at HEAD 0de8f88. The new commit (O-125, O-126) does not change the cited lines.  The cited text says what the finding claims: - docs/plan/README.md:197 (K12): "The prose step stops for every flagged file". - docs/plan/README.md:106: row 20 lists REQ-018. - docs/prd/PRD-0001-layup.md:235: R
- Refuter: C(note): Confirmed as a note. I read the cited text and the sources myself.  The problem is real: (1) docs/plan/README.md:197 (K12) reads "The prose step stops for every flagged file". K12 is settled by O-123 and read by rows 11, 13 and 20. O-123 in docs/tasks/T-55n2.md says the step stops "with one `F-<path
- Refuter: C(note): I confirm the finding as a note. One part of it is refuted.  The cited text is correct: - docs/plan/README.md:197 (K12) says "The prose step stops for every flagged file". - Row 20 at :106 lists REQ-018, and docs/prd/PRD-0001-layup.md:235 gives REQ-018 to T-evad. - docs/spec/setup.md:159 gives the r

### coverage 16 — refuted

- **Where:** docs/plan/README.md:100 (row 14 Tests and After `4, 5`); docs/tests/traceability.md:31
- **Problem:** The fixture test runs the catalog with internal/gate. internal/catalog may import only internal/tsv, and only internal/verify may import both catalog and gate. internal/verify starts in row 7, which is not a predecessor of row 14. Thus an integration test has no legal home when row 14 lands, and the inventory gives this test at e2e level, through the binary.
- **Fix proposed:** Make the handle `T-c06a/e2e/go-fixtures`, or add 7 to row 14's After.
- Refuter: R(none): I do not confirm this finding.  The quotes are correct: - docs/tests/traceability.md:31 gives `T-c06a/integration/go-fixtures` at level integration. - docs/plan/README.md:100 (row 14) gives "integration (each active kind passes on a clean tree and fails on its fixture); e2e" and After "4, 5". - docs
- Refuter: R(none): Not confirmed. The finding quotes the cited text correctly, but the problem that it claims does not follow from the sources.  The facts are correct: - docs/plan/README.md:100 (row 14) gives "integration (each active kind passes on a clean tree and fails on its fixture); e2e", with After "4, 5". - do
- Refuter: R(none): Not confirmed. The parts of the finding are correct. The conclusion does not follow from them.  Correct parts: - packages.md:41 lets `internal/catalog` import only `internal/tsv`. - In packages.md:36-44, only `internal/verify` may import both `internal/catalog` and `internal/gate`. - Row 7 (`T-6x75`

### claims 1 — survives

- **Where:** docs/plan/README.md:96 (row 10), with :98 (row 12)
- **Problem:** Row 10 holds the item verify-sources. Check `sources` must resolve "each `fact` ref is a fact of `docs/facts/`". The fact resolver is in verify-facts, which is row 12. Row 10 has the predecessors 4 and 7 only, and row 12 has 7 only. Thus rows 10 and 12 have no order, and row 10 can start before the resolver exists. The plan does not say why it removed this edge of the inventory. plan-check.py passes because it does not compare the edges with the inventory.
- **Fix proposed:** Set the After cell of row 10 to "4, 7, 12". This makes no cycle, because row 12 depends only on row 7. As an alternative, move verify-sources to a row that comes after row 12.
- Refuter: C(material): I could not refute the finding. Every part of it checks out.  1. **The cited text says what the finding claims.**    - docs/plan/README.md:96 is row 10 (`T-8vpw`, Items `verify-markers`, `verify-sources`, `verify-baseline-scripts`, `gov-issue-21`). Its After cell is "4, 7".    - docs/plan/README.md:
- Refuter: C(material): I could not refute the finding. Each part of it is correct.  1. The cited text is correct. docs/plan/README.md:96 (row 10, T-8vpw) has the items `verify-markers`, `verify-sources`, `verify-baseline-scripts`, `gov-issue-21` and the After cell "4, 7". docs/plan/README.md:98 (row 12, T-9t1q) has the it
- Refuter: C(material): I tried to refute the finding and could not. Every claim in it is true, and the edge it names is missing.  1. The cited cells are correct.    - docs/plan/README.md:96 (row 10, `T-8vpw`) holds `verify-sources` and has After "4, 7".    - docs/plan/README.md:98 (row 12, `T-9t1q`) holds `verify-facts` a

### claims 2 — unverified

- **Where:** docs/plan/README.md:202 (K17) and :204 (K19); rows at :95 and :93
- **Problem:** The register names "the tasks that read the settled text". A reader must thus come after the task that settles. Row 9 reads K17, because at S04 (O-124) it writes the record of the S01- and Q- answers into docs/facts/. But the After cell of row 9 is "4, 6, 8, 12", and none of these rows depends on row 10. Thus row 9 can merge before row 10 settles K17. K19 has the same defect: "\| K19 \| row 5 \| rows 7, 15 \| …", but the After cell of row 7 is "1, 2, 3", not 5.
- **Fix proposed:** Set the After cell of row 9 to "4, 6, 8, 10, 12". For K19, set the After cell of row 7 to "1, 2, 3, 5", or remove row 7 from the readers of K19 if the verify frame does not read that text. Neither change makes a cycle.
- Refuter: C(material): I could not refute the finding. The cited text says what the finding claims. The only difference in the quote is the backticks that the file has around docs/facts/.  1. What the register says. docs/plan/README.md:181-182 says: "Each conflict of the inventory (K1 to K41), with the task that settles i

### claims 3 — survives

- **Where:** docs/plan/README.md:108; rows :105-106; docs/glossary.md:147
- **Problem:** The After cell of row 19 (the release review) is "16, 17, 18". The After cell of row 20 (the first pilot) is also "16, 17, 18". No edge puts them in an order, so the release review can run after the pilot. Thus the claim "last child" and the glossary sentence "It is the last task of phase 1" are not true of the order that a reader follows. O-121, as recorded, does not say "last". The inventory put gov-release-review after gov-first-pilot, and the plan removes that edge without a reason.
- **Fix proposed:** Select one order and put it in the After column. Either set the After cell of row 20 to "16, 17, 18, 19" and cite O-122 for "last", or set the After cell of row 19 to "20" and remove "last" from plan:108 and glossary:147.
- Refuter: C(material): The finding is real. In docs/plan/README.md, row 19 (T-efmy, the release review, line 105) has the After cell "16, 17, 18", and row 20 (T-evad, the first pilot, line 106) has the same cell. No edge puts the two rows in an order. The row number does not give the order: row 9 has the After cell "4, 6,
- Refuter: C(material): I could not refute it. I checked at HEAD 0de8f88. That new commit (O-125, O-126) does not change these lines, so the line numbers stay the same.  The cited text says what the finding claims: - docs/plan/README.md:105: row 19 `T-efmy` (the release review) has After "16, 17, 18". - docs/plan/README.md

### claims 4 — unverified

- **Where:** docs/plan/README.md:106 (row 20); docs/tests/traceability.md:40; docs/tasks/T-55n2.md:46-50 (O-123)
- **Problem:** The plan and the traceability table say that the first pilot proves REQ-018. Its criterion is: "A target's baseline rules are byte-identical to the pinned Armature's after setup, except the adapted values the setup records with evidence". By O-123, the prose step stops for "each flagged kept file that has no input file", and "a person or an agent writes those files outside layup". The measured flagged files include docs/engineering-discipline.md (69), docs/ci/README.md (37), docs/tests/README.md (21) and docs/issue-workflow.md (11). All of them are rule paths. The setup record has rows for values, not for whole replaced files. Thus the rule files of the pilot target will not be byte-identical, and no recorded evidence covers the change, so the claim fails as planned. No K row, no host and no known limit records this conflict.
- **Fix proposed:** Record the conflict as K42, with a host before row 13. For example, the prose step writes one setup-record row (source and ref) for each replaced file, and the REQ-018 audit of row 20 accepts only changes that are recorded in this way. As an alternative, get the Operator's reading of REQ-018 for the prose step. Or remove REQ-018 from row 20 and from traceability.md:40, and write it as a known limit.

### claims 5 — unverified

- **Where:** docs/prd/PRD-0001-layup.md:218, :219, :221, :224, :236 (with the preamble at :212-214 and the §13 row at :251)
- **Problem:** The preamble says that the plan fills the Task column with "its milestones (`M2a` to `M4c`) for the later phases", and the §13 row says the same. The table "What phase 1 proves" of the plan gives later milestones for these criteria: REQ-001 M2c and M4b; REQ-002 M2d and M4c; REQ-004 M2f and M4c; REQ-007 M2e and M2g; NFR-001 M4c. But the Task cells of REQ-001, REQ-002, REQ-004 and REQ-007 name no milestone, and NFR-001 names only M2a. Other phase-1 rows (REQ-009, REQ-011, NFR-002, NFR-005) do name their later milestones, so the column has no single rule. A reader of §12, "the one place a requirement is followed end to end", will think that the phase-1 tasks deliver these four criteria fully. This is the opposite of O-122.
- **Fix proposed:** Add the later milestones: REQ-001 `M2c`, `M4b`; REQ-002 `M2a`, `M2d`, `M4c`; REQ-004 `M2f`, `M2g`, `M4c`; REQ-007 `M2e`, `M2g`; NFR-001 add `M4c`. Then write the rule of the column in the preamble in one sentence.

### claims 6 — unverified

- **Where:** docs/prd/PRD-0001-layup.md:218, :221, :238, :240; docs/tests/traceability.md:40
- **Problem:** The plan says that only row 20 (T-evad) proves the audit half of NFR-003 ("the audit of the pilot's values (row 20)"). But the NFR-003 Task cell does not name T-evad. The same defect occurs in other cells. REQ-001 names T-5zmw but not T-evad (plan:123 gives row 20). REQ-004 does not name row 15 (T-d6q5, a child of T-b97r), but plan:125 gives "rows 14 and 15" and "(row 15)". NFR-005 does not name row 7 (T-6x75), but plan:133 gives "(rows 5 and 7)". Also, the Covers cell of the traceability row of T-evad does not name REQ-001 and NFR-001, which plan:123 and :129 give to row 20.
- **Fix proposed:** Add T-evad to REQ-001 and NFR-003. Add T-d6q5 (or T-b97r) to REQ-004. Add T-6x75 (or T-b97r) to NFR-005. Add REQ-001 and NFR-001 to the Covers cell at traceability.md:40.

### claims 7 — unverified

- **Where:** docs/plan/README.md:129, :124, :95
- **Problem:** By O-124, step S04 writes the raw fact record of the S01- and Q- answers, and S04 is in row 9. But the NFR-001 cell does not name row 9 for "the answers in the target's Git". The Requirements cell of row 9 ("NFR-003, NFR-006, REQ-001, REQ-002") also does not name NFR-001, although records.md makes the answers records item 3 of NFR-001. The inventory gave NFR-001 to setup-s06-facts, and O-124 moves this record from that item to S04. In the REQ-002 row, "the steps of the step table … (rows 9 and 13)" does not name row 15, which holds S12 to S15.
- **Fix proposed:** Change the NFR-001 cell to "rows 9, 13, 15, 16 and 20". Add NFR-001 to the Requirements cell of row 9. Change the REQ-002 text to "the steps of the step table (rows 9, 13 and 15)".

### claims 8 — unverified

- **Where:** docs/plan/README.md:123
- **Problem:** The first sentence of the §7.1 criterion also says: "and the idea owner's answers to it are the raw fact `F-0004`, one fact per question (check `facts`)". This part is phase 1 and is already green (T-zmj6; the §12 Test cell names "check facts (F-0004)"). But neither column of the table names it. By O-122, phase 1 is accepted only against the parts in this table. Thus the PDR would approve a phase-1 acceptance that does not include this clause. The NFR-006 row shows how to write a part that is already green.
- **Fix proposed:** Add to the phase-1 column: "the answers to LAYUP's own batch are the raw fact F-0004, one fact per question (check facts; already green, T-zmj6)".

### claims 9 — unverified

- **Where:** docs/glossary.md:147; docs/prd/PRD-0001-layup.md:144-146; docs/facts/F-0004-psb-gap-answers.md:28
- **Problem:** The new row defines "First pilot" as T-evad. By O-122, this pilot "does not measure the baseline of §8". But PRD-0001 §8 says "the first pilot measures the baseline with the current process", and it derives from the raw fact F-0004#11 ("The first pilot measures the baseline with the current process"). With the new term, these sentences give T-evad the duty to measure the baseline, which contradicts O-122. The collision line names the pilot of phase 4, but it does not say that §8 and F-0004#11 call that pilot "the first pilot". The register (plan:207) says that this row settles K22, but the PRD sentence that drives a decision still has two readings.
- **Fix proposed:** Add to the collision line: "PRD-0001 §8 and F-0004#11 use 'the first pilot' for the pilot that measures the baseline. By O-122, that is the pilot of phase 4, not this one." As an alternative, change PRD §8, which is not a §7.1 criterion and so is not barred by O-115, and add a §13 row. Do not change F-0004.

### claims 10 — unverified

- **Where:** docs/plan/README.md:197 (K12); also :210 (K25)
- **Problem:** K12 has two halves. First, check `adapted` fails on 21 kept files. Second, check `kit-history` fails after S05, because 32 lines of completed.md and 7 lines of backlog.md still link github.com/pharzam/armature/. O-123 decides only the prose step and check `adapted`. The S05 change changes the S05 row of setup.md ("their lines in `backlog.md` and `completed.md`"), and no Operator decision contains it. But the "Settled by" cell names only O-123. The author of row 13 will read the S05 change as the Operator's decision, and will not give the "decided here" reason that spec/README.md asks for. K25 has the same defect with less effect: O-121 does not say where the rule for the pilot numbers is written.
- **Fix proposed:** Write "Settled by: O-123 (check `adapted`); row 13 (check `kit-history`: the line rule of S05, decided here with its reason)". For K25, write "O-121; this plan (guardrails §1)".

### claims 11 — unverified

- **Where:** docs/tests/dod-checklist.md:25, :73-74, :80; docs/tests/README.md:42; docs/tests/traceability.md:12-14
- **Problem:** This change makes docs/tests/traceability.md the traceability table of LAYUP. The table also tells each delivering task to "put the real test name in its place, with its status, in its own pull request". But the close-out checklist, which is in force, still sends the reader to the template. The index of docs/tests/ ("What's here") does not list the new file. Only one Reason cell of the plan (plan:144) links to it. Thus a phase-1 task that runs the DoD checklist will open the template, find no rows, and not update traceability.md. Then the rows stay `planned`, and the DoD rule cannot close.
- **Fix proposed:** Point dod-checklist.md:25, :73-74 and :80 to traceability.md, and keep the template as the link for the format. Add a row for traceability.md to docs/tests/README.md. In "How to read this plan", say that each task updates its rows in traceability.md and its §12 Test cell.

### claims 12 — unverified

- **Where:** docs/tests/traceability.md:28, :29, :39, :40
- **Problem:** Each row that covers more than one requirement gives the Fact and ADR of only one requirement. The schema-blocks row does not give F-0003#49 and ADR-0023 (REQ-009) or F-0003#50 and ADR-0024 (REQ-011). The package-rules row (:29) does not give F-0004#1 and ADR-0010 (NFR-007). The release-review row (:39) does not give F-0003#55 (REQ-017). The first-pilot row (:40) names ADR-0012, which PRD §12 gives to none of the requirements that the row covers. Also, the telemetry and stall Go schemas come with rows 17 and 18, so the T-18v6 test cannot cover REQ-009 or REQ-011 when it merges.
- **Fix proposed:** Give the fact and the ADR of each covered requirement, or split the rows. Remove REQ-009 and REQ-011 from the T-18v6 row, because rows 37 and 38 cover them.

### claims 13 — unverified

- **Where:** docs/tests/traceability.md:4
- **Problem:** At 7cdd346 there are seven tests, but the table has rows for only three. TestBinaryPrintsItsVersion (cmd/layup/main_e2e_test.go:12), TestVersionPrintsTheVersionAndExitsZero, TestNoSubcommandPrintsUsageAndExitsTwo and TestUnknownSubcommandExitsTwo (internal/cli/cli_test.go:15, :28, :38) have no row and no reason.
- **Fix proposed:** Add a row for each of these tests (for example, they cover the exit-code rules of docs/spec/README.md "Commands" or a DoD item), or say why the table does not include them.

### claims 14 — unverified

- **Where:** docs/plan/README.md:116
- **Problem:** O-122 decides that phase 1 is accepted against "the plan's table, operative with the Operator's approval at the PDR". The author wrote the cells, and the Operator has not approved them yet. Thus the label gives the Operator's authority to text that the Operator has not seen.
- **Fix proposed:** Write: "The plan's reading of each phase-1 criterion …; it becomes the Operator's reading when the Operator approves the PDR (O-122)."

### claims 15 — unverified

- **Where:** docs/plan/README.md:224; :60-62
- **Problem:** packages.md:24-25 names two later exceptions: internal/smartif (phase 3) and the forge adapter (phase 2). The first later reader is thus M2a (the forge adapter of later-p2-run-start), not M3a. The Milestones paragraph also says that M2a takes "the conflicts K38 to K41" as inputs, and this does not agree with the register row.
- **Fix proposed:** Write "Read by: `M2a`, `M3a`" and add "`M2a` names the forge adapter".

### claims 16 — unverified

- **Where:** README.md:58 (also AGENTS.md:185)
- **Problem:** README.md:20 and docs/spec/README.md:9 say "one milestone at a time" (O-114). The plan gives phases 2 to 4 seven, four and three milestones, and each has its own specification task. The new row and AGENTS.md:185 ("of each phase") do not agree with this.
- **Fix proposed:** Write "one milestone at a time" in both rows.

### claims 17 — unverified

- **Where:** docs/glossary.md:146; docs/plan/README.md:89 (row 3)
- **Problem:** architecture.md §8 (:740) asks only for "a size class". The values `small` and `large` are the plan field `size` of §9 (:923), which selects the tier for a target and has no rule for packages, days or lines. The definition comes from the plan (plan:36-38), but the row does not limit it to LAYUP's plan, as the Milestone row does. Thus it reads as the meaning of the `size` field of a target. Also, row 3 is `small` but touches two packages (internal/cli and cmd/layup) and has 400 lines, which is outside its own definition.
- **Fix proposed:** Name the row "Size class (of LAYUP's plan)". Cite plan:36-38 for the values and §9 for the field of a target. Set row 3 to `large`, or split it.

### claims 18 — unverified

- **Where:** docs/glossary.md:144
- **Problem:** architecture.md:680 also uses "the preliminary design review", for a target: "the specification of step 6 with the architecture that the first bet approves". The collision line does not name this second meaning.
- **Fix proposed:** Add: "and not the preliminary design review of a target, which is its specification with the architecture that the first bet approves (architecture.md §7)".

### claims 19 — unverified

- **Where:** docs/glossary.md:147
- **Problem:** ADR-0012 part 6 says: "The mode ends when the ADR that supersedes this record is accepted, and not before". The task of the pilot only opens that ADR. The row gives the rule in different words, and a reader can think that the full gate applies again when T-evad merges.
- **Fix proposed:** Write "The pilot whose task opens the ADR that ends bootstrap mode (ADR-0012 part 6; O-122)".

### claims 20 — unverified

- **Where:** docs/plan/README.md:149, :237
- **Problem:** The inventory item closes #24 "with the evidence and one note": docs/ci/README.md:37 ("Every job restores its check scripts") is false for nested-checkout-check.sh, and "the fix of that sentence is batched into the first product task that edits docs/ci/README.md or ci.yml". The plan gives no host and no known limit for this sentence. Thus it is lost when #24 closes.
- **Fix proposed:** Give the sentence fix a host (for example row 16, if it changes ci.yml, or this task), or write it as a known limit.

### claims 21 — unverified

- **Where:** docs/plan/README.md:248-249
- **Problem:** R12 sends a step that grows past its scale back to the issue as a new slice or a child issue. The number of rounds is the cycle cap, and this plan gives a cap of 2 to rows 10 and 12 (and to row 16 if ci.yml changes). Thus "one review round" gives both rules in wrong words (Bootstrap mode rule 7).
- **Fix proposed:** Link R12 and Bootstrap mode rule 3, or write: "a task that grows past its budget goes back to its issue as a new slice (R12); a task that reaches its cycle cap with a material finding is split".

### claims 22 — unverified

- **Where:** docs/plan/README.md:44 and :121-135
- **Problem:** Rows 19 and 20 of M1 serve REQ-015 to REQ-018, and PRD §9 says these rows "hold in every phase". But the Requirements cell of M1 and the acceptance table name none of them. Thus the PDR approves a phase-1 acceptance with no line for the Won't rows.
- **Fix proposed:** Add REQ-015 to REQ-018 to M1. Add one row for each of them to "What phase 1 proves" (row 19 for REQ-015 and REQ-017, row 20 for REQ-016 and REQ-018).

### claims 23 — unverified

- **Where:** docs/plan/README.md:105
- **Problem:** The release that row 19 reviews includes rows 17 and 18 (After "16, 17, 18"), and these rows serve F-0003#50 and F-0003#49. The Fact cell does not name these two facts.
- **Fix proposed:** Add F-0003#49 and F-0003#50 to the Fact cell.

### claims 24 — unverified

- **Where:** docs/prd/PRD-0001-layup.md:242
- **Problem:** The criterion of NFR-007 is: `go list -deps ./...` names no module outside the standard library, and Git is called only as the `git` program. gofmt and go vet do not test this. The inventory item that this task hosts (gov-prd-task-column) asked to mark this cell. The change edits this row but keeps the false Test claim.
- **Fix proposed:** Change the Test cell to show that the criterion has no test yet, for example "— (T-2tc2/integration/package-rules, planned)".

### claims 25 — unverified

- **Where:** docs/tasks/backlog.md:34-35; docs/plan/README.md:87-106 (Issue column); docs/tests/traceability-template.md:64-66
- **Problem:** The 20 phase-1 tasks and T-b97r have the Issue "—", no backlog line and no task file. The plan does not say who opens each phase-1 issue or when, but for phases 2 to 4 it says "opens its own issues then". The template says that each task ID in a traceability table "matches the backlog line", and none of the 13 plan IDs in traceability.md has one. The T-stfn line (backlog.md:35) still names only the four deliverables, but by O-121 #29 now also holds the release review and the first pilot. Compare "the phase-1 issues" with the DoD of #76.
- **Fix proposed:** In "How to read this plan", say when each phase-1 issue is opened (for example after T-4wrw, in the order of the rows) and add the backlog lines then. If the issues are not opened by this task, change the summary to "the phase-1 task list". Update the T-stfn line.

### decisions 1 — unverified

- **Where:** docs/plan/README.md:199 (K14); :95 (row 9); :98 (row 12); :203 (K18); :129 (NFR-001 row of "What phase 1 proves")
- **Problem:** O-124 moves the record of the `S01-` and `Q-` answers to S04. It keeps only the briefs at S06. Three more operative texts still put this record at S06: the rule of check `facts` (setup.md:246: "each brief is in `docs/facts/` ...; the `S01-` and `Q-` IDs in the record of S06"), records.md:105 ("two raw fact records (S06, S11)") and psb-check.md:69 ("at setup step S06"). Row 12 builds check `facts` and settles K18 before row 9 (row 9 After: 4, 6, 8, 12). But K14 does not send O-124 to row 12. Row 9 then uses "check `pin` and check `facts` as evidence" at S04, when no brief is in `docs/facts/` yet. Thus a check that row 12 builds to setup.md:246 can fail a correct S04, and row 9 must open the merged work of row 12 again. K14 names only setup.md, so no task is sent to change records.md:105 or psb-check.md:69 (row 6 edits psb-check.md before row 9). Also: row 9 now needs row 12, which the table lists after it, so the row numbers are not an order. Row 9 does not list NFR-001. The NFR-001 cell gives "rows 13, 15, 16 and 20; O-115" for the answers in Git, not row 9 and O-124.
- **Fix proposed:** Best: write the readings of O-124 and O-123 into docs/spec/ in this pull request: the S04, S06 and S14 rows and the `facts` row of setup.md, records.md:105 and psb-check.md:69. Else: K14 "Read by: rows 6, 9, 12, 13", and give each text to the first task that reads it (row 12: setup.md:246 and records.md:105, with check `facts` at S04 before the briefs; row 6: psb-check.md:69; row 9: the S04 and S06 rows). In both cases: say that the After column, not the row number, gives the order; add NFR-001 to row 9; write "rows 9 and 13 (O-124)" in the NFR-001 cell.

### decisions 2 — unverified

- **Where:** docs/plan/README.md:101 (row 15); see also :99 (row 13)
- **Problem:** The range "S12 to S15" includes S14. O-123 puts S14 in the prose step, which runs before S10. Row 13 holds the item of S14 (`setup-prose-inputs`: S07, S08, S09, S14). Thus two rows claim S14, and the title of row 15 puts S14 after S13, against O-123. The author of row 15 can build S14 again, or in the old position.
- **Fix proposed:** Write row 15 as "Steps S12, S13 and S15, and the checks `jobs` and `gate:<kind>`". Write row 13 as "Steps S05 to S11 and S14: the prose step (S07 to S09 and S14) before S10 (O-123), and the record of the `M-` answers at S11 (O-124)".

### decisions 3 — unverified

- **Where:** docs/prd/PRD-0001-layup.md:143-146 (§8); docs/glossary.md:147; docs/plan/README.md:207 (K22)
- **Problem:** O-122 and the new glossary row make "first pilot" the pilot of ADR-0012 part 6. O-122 says that this pilot "does not measure the baseline of §8". But PRD §8 still says that "the first pilot measures the baseline". It takes these words from the idea owner's fact F-0004#11, which has the same words. With the new term, these texts have two readings (R13). The author of row 20 can read them as an instruction to measure the baseline. K22 cites PRD:144-145 as a source of "two meanings of 'pilot'" and is marked settled. But the PRD does not change, and the collision note of the glossary row does not name this use. The note also puts the baseline measurement under REQ-014, but §9 phase 4 names it, not REQ-014.
- **Fix proposed:** In PRD §8, write "the pilot of phase 4 measures the baseline", with a §13 row (O-115 and O-122 freeze only the §7.1 words; F-0004 stays as it is). In the collision note of the glossary row, name F-0004#11 and PRD §8, where "the first pilot" is the pilot that measures the baseline (phase 4: `M4b`, `M4c`). Cite §9 phase 4, not REQ-014, for the baseline measurement.

### decisions 4 — unverified

- **Where:** docs/plan/README.md:114-135 ("What phase 1 proves")
- **Problem:** O-122 accepts phase 1 "against the parts of each §7.1 criterion that phase 1 can prove (the plan's table ...)". The table has no row for REQ-015 to REQ-018. These requirements "hold in every phase", and phase-1 tasks prove parts of them: row 19 (REQ-015, REQ-017) and row 20 (REQ-016, REQ-018). Thus the basis that the Operator approves at the PDR does not include the release review or the byte comparison of the baseline rules. For REQ-018 the reading is not simple: by O-123 the pilot replaces flagged kept files, for example engineering-discipline.md and issue-workflow.md. No milestone holds the review of each later release (REQ-015, REQ-017: "a code review of each release"). Also, the table is the plan's reading until the PDR, not yet the Operator's.
- **Fix proposed:** Add rows REQ-015 to REQ-018. Phase 1 proves: the recorded code review of the phase-1 release (row 19); the intent decisions of the pilot, recorded as the idea owner's (row 20); the byte comparison of the pilot target's baseline rules, with the rule for the files that the prose step replaces (row 20). A later milestone proves: the review of each later release (give it a host); the pilot of phase 4 (`M4c`). Write "The plan's reading ... (O-122); it becomes the Operator's at the PDR".

### decisions 5 — unverified

- **Where:** docs/plan/README.md:197 (K12)
- **Problem:** O-123, as T-55n2.md quotes it, decides only the prose step, and it keeps the rule of check `adapted`. It says nothing about S05 or check `kit-history`. The second clause changes S05: setup.md:90 deletes only "their lines", the lines of the deleted T-IDs. Thus it changes the target's copy of backlog.md and completed.md. The inventory gives this change as "An Operator decision, or a spec fix". The register gives it the authority of the Operator. Row 13 then writes it into the specification with no "decided here" reason, and the plan review of row 13 cannot change it.
- **Fix proposed:** Split K12. The part about check `adapted`: settled by O-123. The part about check `kit-history`: settled by row 13, as a "decided here" value of S05 in setup.md, with its reason. If the option that the Operator chose held this clause, quote its text in T-55n2.md under O-123, with the comment link.

### decisions 6 — unverified

- **Where:** docs/plan/README.md:210 (K25), :108, :76-77, :152, :187 (K2); the After cells of rows 19 and 20 at :105-106
- **Problem:** O-121, as T-55n2.md quotes it, does not decide the rule or the count of the pilot's numbers. It does not make the pilot the last child of #29. It says only: "The release review of phase 1 and the first pilot are two more children of #29". ADR-0012 part 6 also does not say that the ADR task counts the numbers (:76-77, :152). Thus the plan gives its own choices the authority of the Operator or of the ADR, and the plan review of row 20 cannot change them. Also, "last" is not in the order: row 20 is "After 16, 17, 18", so row 19 (the release review of phase 1) can close after the pilot that ends phase 1 (O-122). K2 says "The four deliverables as parents of child tasks", but O-121 makes only `layup gate` and `layup setup` parents; the telemetry and stall records are one task each.
- **Fix proposed:** K25: "Settled by `T-55n2`" (the plan's own placement; the plan review of row 20 can change it). Line 108: cite O-122 ("the end of phase 1") for "last", and add 19 to the After cell of row 20. Lines 76-77 and 152: "the ADR reads the pilot's numbers (ADR-0012 part 6)"; give the count to the ADR task as the plan's choice. K2: "`layup gate` and `layup setup` as parents of child tasks; the telemetry and stall records one task each".

### decisions 7 — unverified

- **Where:** docs/plan/README.md:87-89 (rows 1 to 3), :92 (row 6), :101-102 (rows 15, 16)
- **Problem:** O-121 gives each child a goal count of one because "Each child is one part of one command". Rows 1 to 3 (`internal/tsv` and the schema-block test, `internal/git` and the package-rule test, the command frame and the e2e harness) serve all commands, but they are children of `layup gate`. Rows 15 and 16 join steps of `layup setup` and checks of `layup setup verify`. Row 6 holds all of `layup psb check` to its specification (`psb-cli-contract`, `psb-rule-edges`, `psb-rule-values`, `psb-utf8-input`), but O-121 gives T-b97r only "the S01 work of `layup psb check`". A plan reviewer of these rows can find the premise false and count goal classes; rule 2 then sends the count to the Operator.
- **Fix proposed:** In "How to read this plan", say that the count of one of O-121 applies to each row as written (if the option that the Operator chose showed these rows), or split the rows on command lines. For row 6, keep the S01 work (`psb-batch-api`, `psb-tsv`) under T-b97r and put the rest in a task that serves F-0003#41, or get the Operator's word that "the S01 work" means all of psb-check.md.

### decisions 8 — unverified

- **Where:** docs/plan/README.md:99 (row 13) and :101 (row 15); also :91 (row 5) and :93 (row 7)
- **Problem:** Rows 13 and 15 each expect 1,200 lines, and rows 5 and 7 expect 1,100. That is about three times the bound that the inventory gives for one review round in bootstrap mode ("about 400 lines of diff"), and the plan's own `small` class. Row 15 holds 8 items of two commands and #34. Row 13 holds five step groups and reads K10, K12 and K14 to K18. With one round and one more after a fix, these rows can stop at the cap. T-0drh (a 1,900-line diff) needed the Operator to raise its cap two times.
- **Fix proposed:** Split row 15 into (a) S12 with the checks `jobs` and `gate:<kind>`, and (b) S13, the rule-path register and S15, with #34. Split row 13 into (a) S05 and S06, and (b) the prose step, S10 and S11. Or write in the plan why 1,200 lines fit one review round.

### decisions 9 — unverified

- **Where:** docs/plan/README.md:82-83 and the Cap column (:87-106)
- **Problem:** The sentence restates rule 3 and makes it narrower. Rule 3 gives cap 2 when the change touches a gate: "a check script, a hook, CI or branch protection". The column also sets each cap now, but by rule 2 the plan review of each task records its `Cycle cap`. A task that touches a hook or branch protection, or a check script that the plan does not expect, gets the wrong cap from this column.
- **Fix proposed:** Write: "'Cap' is the expected cycle cap; the plan review of each task sets it by [Bootstrap mode] (../engineering-discipline.md#bootstrap-mode) rule 3."

### decisions 10 — unverified

- **Where:** docs/plan/README.md:248-249 (Known limits)
- **Problem:** This is a new rule, and its cited sources do not say it. R12 sends back a step that grows past the scale that the plan implied (a size). Rule 3 lets a task fix its findings and get a second round (cap 1). If a task is split when it needs a second round, the plan makes again the successor loop that ADR-0012 gives as root cause 5.
- **Fix proposed:** Write: "Each task's plan review sets its budget; a task that grows past its budget goes back to its issue as a new slice or a child issue (R12)."

### decisions 11 — unverified

- **Where:** docs/tasks/T-55n2.md:18-19
- **Problem:** This review runs as "round 1 of at most three". With `Cycle cap` 1, rule 3 allows this round and one round after a fix. A third round needs the Operator to raise the cap first, with the raised cap recorded on the issue, as O-116, O-117 and O-120 did for T-0drh on #74.
- **Fix proposed:** Stop after two rounds, or get the Operator's raise before a third round and record it in T-55n2.md and in a plan-review comment on #76.

### decisions 12 — unverified

- **Where:** docs/glossary.md:147
- **Problem:** The operative rule says that the mode ends when the ADR that supersedes ADR-0012 is accepted, "and not before". The task that closes the pilot only opens that ADR. A reader of this row can use the full gate when the pilot closes, before the ADR is accepted.
- **Fix proposed:** Write: "The pilot whose closing task opens the ADR that ends bootstrap mode (ADR-0012 part 6; O-122)".

### decisions 13 — unverified

- **Where:** README.md:58; AGENTS.md:185; docs/plan/README.md:60-61
- **Problem:** The new rows say "one phase at a time" (README.md:58) and "of each phase" (AGENTS.md:185). README.md:20, spec/README.md:9 (O-114) and the plan say "one milestone at a time", and the plan gives each of `M2a` to `M4c` its own specification task. Phase 2 has seven milestones, so the texts give one or seven specification tasks for phase 2. Line 60 adds a third reading: "The specification task of phase 2 (`M2a`)".
- **Fix proposed:** Write "one milestone at a time" in README.md:58 and AGENTS.md:185. At docs/plan/README.md:60, write "The specification task of `M2a`, the first milestone of phase 2, takes ...".

### decisions 14 — unverified

- **Where:** docs/plan/README.md:72 (T-b97r); the Issue column at :87-106
- **Problem:** The plan does not say who opens the parent issue of T-b97r and the 20 child issues, or when. Each Issue cell is "—". But the backlog line of T-55n2 promises "the phase-1 issues", and `gov-rescope-29` writes the children table of #29, which then holds a parent with no issue. R1 needs an issue before step 1, and R11 needs a parent issue for child issues. The inventory asked this question (`gov-rescope-29`, open question 3), and the plan gives no answer.
- **Fix proposed:** Add one sentence under "How to read this plan", for example: "T-55n2 opens the issue of T-b97r at close-out; the issue of each row is opened as a child of its parent's issue after #42 closes and when the rows in its After cell have merged." Or open the issues and fill the Issue column.

### decisions 15 — unverified

- **Where:** docs/plan/README.md:237 and :149 (#24); :235 and :148 (#15)
- **Problem:** The inventory keeps one residual of #24: docs/ci/README.md:37 ("Every job restores its check scripts from the default branch") is not true for nested-checkout-check.sh, and its fix goes to the first task that edits docs/ci/README.md or ci.yml. The plan closes #24 and gives this residual no host, so it is lost; ADR-0012 says that nothing is lost silently. Also, the head has no edit of guardrails.md §2 for #15, but the plan says that T-55n2 settles #15.
- **Fix proposed:** Fix the sentence of docs/ci/README.md:37 in this pull request (T-55n2 already edits AGENTS.md, the other file of #24), or give it a host or a known limit. At close-out, write the guardrails §2 lesson and settle the lines of #15.

### decisions 16 — unverified

- **Where:** docs/prd/PRD-0001-layup.md:218, :219, :221, :224, :238, :240 (§12 Task cells)
- **Problem:** The Task cells do not name all the tasks and milestones that "What phase 1 proves" names. REQ-004 omits row 15 (`T-d6q5`, under T-b97r) and `M2f`, `M4c`. REQ-001 omits `T-evad`, `M2c`, `M4b`. REQ-002 omits `M2d`, `M4c`. REQ-007 omits `M2e`, `M2g`. NFR-003 omits `T-evad`. NFR-005 omits `T-6x75`. Other rows (REQ-009, REQ-011, NFR-002, NFR-005) name their later milestones, so the column follows two rules.
- **Fix proposed:** In each phase-1 row of §12, name each row and milestone that "What phase 1 proves" names for that requirement.

### decisions 17 — unverified

- **Where:** docs/glossary.md:146
- **Problem:** architecture.md §8 asks for a size class, but it gives no numbers. The numbers are the plan's (docs/plan/README.md:36-38). The row defines the term for every plan, and architecture.md §9 uses the `size` of a target's task to select its model tier. Thus a reader can apply LAYUP's line estimate to the routing of a target.
- **Fix proposed:** Name the row "Size class (of LAYUP's plan)", say that the numbers are the plan's, with a link, and leave a target's `size` to architecture.md §8 and §9.
