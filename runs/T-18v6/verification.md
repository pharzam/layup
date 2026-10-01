# T-18v6 — the author's verification (not a review round)

A fresh session of Claude Opus 5.5, the author's model, verified the head of the implementer adversarially (the workflow `rows-1-2-implement`, run `wf_79f86cf4-afc`, 2026-10-01); a fixer then applied each finding that it found real. This is the author's own step: the review round of the gate came later, by another model.

14 findings, 1 of them material.

## The findings

1. **material** — docs/spec/README.md:55, :62-67; docs/spec/psb-check.md:3-5, :21-26
   - Problem: The new K32 rule documents exit 2 for a problem statement that is not valid UTF-8. The present `layup psb check` exits 1 and writes the raw bytes into the psb-gaps table. psb-check.md says that it states the contract of the existing code, and it records the one known gap of the present code (a lone CR). It does not record this gap, and its exit-code list (2 = usage error or unreadable file) now disagrees with the README rule for the same input. The README K32 text also drops the sentence of D5 'Row 6 applies it to `layup psb check FILE`'.
   - Fix proposed: In psb-check.md:21-24, add one sentence like the lone-CR one: 'the present code does not yet give exit 2 for a FILE that is not valid UTF-8 (task T-5zmw, row 6, does)'. Or put the D5 sentence 'Row 6 applies it to `layup psb check FILE`' into README:62-67. One line fits the budget.

2. **note** — internal/tsv/tsv_test.go:218-231 (the Read half, :227)
   - Problem: TestWriteAndReadRefuseASchemaThatIsNotValid does not prove that Read refuses 5 of its 6 bad schemas. For no column, a name that repeats, an empty name, a name with a tab and a name that is not UTF-8, Read("a\n") fails on the header comparison or the field count, not on the schema check. Only the 'float' case reaches the schema check through Read. The Write half is sound.
   - Fix proposed: For the Read half, use an input that the bad schema would otherwise accept, for example NoHeader schemas with the input "x\n", or "a\ta\n" for the repeated name. Or assert that the error text starts with "schema x:".

3. **note** — internal/tsv/tsv.go:126; internal/tsv/tsv_test.go:66-81; docs/spec/README.md:118-120
   - Problem: In the no-header form, the reader takes a line that equals the column names as a data row. The found-tsv inventory test says that the no-header form 'does not take a header row as data'. The README does not decide this case, and TestTheNoHeaderFormHasRowsOnly checks the opposite reading ('the reader takes every line as a row').
   - Fix proposed: In the no-header form, refuse a first line that equals the column names joined by one tab, and add a unit test. Or write the chosen reading into README:118-120 as decided here.

4. **note** — internal/tsv/tsv.go:122-123, :207-208; docs/spec/README.md:96-97; internal/tsv/tsv_test.go:124
   - Problem: README:96-97 says 'Each error names the line, and the column when one field is wrong'. For a carriage return inside one field ("a\rb" in excerpt), the error names no column, and the test requires Column "". Too many fields also gives no column; the reason gives the counts. Note 9 (the line and the fix) is met.
   - Fix proposed: Name the column for a CR inside a field. Or reword README:96-97 so that the line-level refusals (CR, byte-order mark, empty line, final line feed, field count) name the line only.

5. **note** — internal/tsv/block.go:50-96; docs/spec/README.md:130, :136-142, :151-159
   - Problem: The parser does not see a block in a blockquote or under a list item with 4 or more spaces of indent. So a new block of that kind passes the completeness test without a list entry, against README:158-159. The parser does read a block inside an HTML comment, which Markdown does not show as a block. It also accepts a closing fence with a tab ("```\t"), against README:130 and :136 ('refuses a block with a tab'). None of these occurs in docs/spec today, and the README does not record the limits.
   - Fix proposed: Write these limits into README:136-142 as known limits. Or refuse a line with 'tsv-schema' in a fence that the parser cannot place, and refuse a tab on the closing fence line.

6. **note** — internal/tsv/block.go:17
   - Problem: fs.Glob ignores I/O errors, so ReadBlocks returns an empty map and no error when the directory is missing or cannot be read. The present callers still fail loudly: the integration test lists 14 names, and Compare against the zero Schema reports differences. But an owner test that ranges over the result would pass with no block.
   - Fix proposed: Read the directory with fs.ReadDir and return its error, or return an error when no *.md file is found.

7. **note** — internal/tsv/tsv.go:149-155; internal/tsv/block.go:47, :58, :64, :88; internal/tsv/tsv_test.go:176; internal/tsv/block_test.go:12-23, :44-61
   - Problem: The code comments and the report state five behaviours that no test covers: JoinList refuses a CR; a closing fence can be longer than the opening fence; a fence line with an info string does not close a fence; a backtick fence with a backtick in its info string is not a fence; a fence with 4 or more spaces of indent is not a fence.
   - Fix proposed: Add "a\rb" to the JoinList cases at tsv_test.go:176. Add one parseBlocks case for each of the four fence rules.

8. **note** — internal/tsv/tsv.go:56-61, :148-160; docs/spec/README.md:177
   - Problem: D7 says 'the writer refuses' a list value with a space. Only JoinList refuses it. Write accepts a pre-joined field such as "docs/my file.md" in a list(path) column with no error, and that field reads back as two values. Neither the Write comment nor the README tells a caller to make a list field with JoinList.
   - Fix proposed: State in the Write comment and in README:177 that a list field is made with JoinList.

9. **note** — internal/tsv/types.go:100-102; docs/spec/README.md:176, :179-184
   - Problem: Because a run of N is a minimum count, id(Q-NNN) also takes Q-0001, which no writer makes (Q-%03d). So one ID has two forms, and the key check takes Q-001 and Q-0001 as two different keys. The README uses 'one form' as the reason for the path reading, but ids do not get the same rule.
   - Fix proposed: Optional: make a run of N match exactly N digits, or more digits with no leading zero. Or say in README:176 that the owner checks a zero beyond the run.

10. **note** — docs/spec/README.md:174
   - Problem: The path row says 'no `..`', which a reader can take as 'no `..` anywhere'. The code refuses only a `..` part and accepts 'a..b/c'.
   - Fix proposed: Write 'no `..` part'.

11. **note** — internal/tsv/blocks_integration_test.go:18
   - Problem: The comment says that notYetBuilt 'lists each block that no package writes or reads yet'. But internal/psb writes the psb-gaps table today, and psb-gaps is on the list.
   - Fix proposed: Write 'lists each block whose owner does not yet compare it with a Go schema'.

12. **note** — branch diff against origin/main b48764f; plan review of #78 (Budget maximum 1,200 lines over 18 files, base 7cdd346); docs/engineering-discipline.md:345-351
   - Problem: The branch diff is 1,191 lines (1,183 added, 8 removed) over 10 files before the close-out, so only 9 lines remain. The close-out cannot fit: earlier close-outs added 53 and 69 lines, and their task records have 119 and 139 lines. The plan review names base 7cdd346, but the branch now sits on b48764f, and the budget rule measures against a named base. An overrun that the Operator did not approve blocks the merge.
   - Fix proposed: Before the freeze, name b48764f as the base on #78. Then get the Operator's one-time approval of the landing figure, or move part of the work to a child issue.

13. **note** — docs/tests/traceability.md:34; docs/tasks/backlog.md:37
   - Problem: Some items of the author's answer are still open at the head. The traceability row T-18v6/integration/schema-blocks is still `planned` with no test name (note 6). Note 5 says this edit lands before the freeze; #76 has merged, and the branch is on b48764f. Condition 2(e) needs the task record, which does not exist yet. The issue's acceptance box 'PRD-0001 §12 names its tests' disagrees with note 6 (no §12 cell), so the task record must give the reason.
   - Fix proposed: Before the freeze: set the row to `internal/tsv TestEverySchemaBlockIsBuiltOrNotYetBuilt`, green. Write the task record with condition 2(e) and the §12 reason. Move the backlog line to completed.md.

14. **note** — runs/T-18v6/red.md:9, :35-40
   - Problem: red.md:9 says that '[...] marks where like lines are cut'. But the excerpts for TestParseTypeRefusesATypeOffTheList and TestTypeCheckTakesTheValuesOfItsType leave out the first lines of the run (parseType(""), and the text and int lines) with no [...] marker. Everything else in the red evidence reproduces.
   - Fix proposed: Put '[...]' before the lines shown, or show the first lines of the run.

## What the fixer did not apply, with its reasons

- Finding 12 (budget), not fixed. To name the base on #78, and to get the Operator's one-time approval or open a child issue, are the author's steps on the issue. The task rules forbid gh. This round adds 243 lines and removes 50 against 0db1aa1.
- Finding 13 (traceability row, task record, backlog to completed), not fixed. The task rules forbid edits to docs/tests/traceability.md and docs/tasks/; these items are the author's close-out.
- Finding 5, in part by documentation only. A tsv-schema block in an HTML comment is still read as a block. The README records this as a known limit, and the parser does not read HTML comments. Reason: the case fails loudly (the block must be listed or matched), docs/spec has no HTML comment, and a comment parser adds code for no present case. A side effect of the new refusal: a tsv-schema fence that is text in an indented code block (4 spaces or more) is also refused. That refusal is loud, and the README says 'four spaces or more'.
- Finding 4: I chose to name the column for a CR, not to reword the README only. A byte-order mark stays an error of line 1 with no column in the reader, because it is a fact of the file and not of a field. The README now says which errors name a column, so this choice is written down.
- Finding 8: no code change is possible without a new API, because Write gets a joined string. The fix is the stated contract (JoinList), as the finding proposed.
- Base: the worktree branch sits on b48764f (origin/main after #76 merged), not on 7cdd346 as the task text says. All diff figures are against b48764f.
