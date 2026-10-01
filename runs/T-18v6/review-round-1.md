## Review record — round 1

| Field | Value |
| ----- | ----- |
| Commit reviewed | `6e0b4c29fb5694ecc4b89a1c1929036487051c71` |
| Reviewer | Claude Fable 5.1 on Claude Code |
| Lens | correctness and acceptance criteria |
| Briefed on | `.review-in/brief.md`; `.review-in/issue-body.md`; `.review-in/issue-comments.md` (the plan, the plan review, the author's answer, the status comment); the whole diff `b48764f..6e0b4c2` (14 files, 1,441 lines added, 10 removed); `docs/plan/README.md` (row 1, the defect register, the table "What phase 1 proves"); `docs/spec/README.md` in full; `docs/spec/packages.md`; `docs/spec/psb-check.md`; the 14 schema blocks of `gate.md`, `psb-check.md`, `records.md` and `setup.md`; `docs/spec/setup.md` (The answers); `docs/tests/test-levels.md`; `docs/tests/traceability.md`; `docs/engineering-discipline.md` (Bootstrap mode, Testing, Who may review, What a round records); `docs/prd/PRD-0001-layup.md` §7 and §12; `runs/T-18v6/red.md` and `test-runs.md`; `runs/T-55n2/inventory.md` (`found-tsv`, `found-spec-schema-test`); `internal/psb/check.go`; `internal/cli/cli.go`; `go.mod`; `.github/workflows/ci.yml`; the commit log of the branch with its dates; `docs/ci/review-record-lint.sh` (the record fields). |
| Barred from | nothing; the comments of #78 are the plan and its record |
| Independence claimed | **Execution:** held; this record is a separate run. **Model:** held; the reviewer is Claude Fable 5.1, and each author is Claude Opus 5.5 (the author's session, the implementer, the verifier and the fixer agents). The plan reviewer was also Claude Fable 5.1; a plan reviewer is not an author. **Context:** held in part; a fresh read-only session in a clone at the frozen head, with no earlier round's verdict, but the brief handed me the plan, the plan review, the author's answer and the status comment, which hold the author's reasoning. **Method:** not claimed; this is the first round. |
| Cycle | 0 |
| Verdict | `nothing material in scope` |

### Raw findings

Each finding gives where, the cited sentence or code, the basis, and its class.
No finding is material. The tests pass at the head, the probes found no bug,
each condition and each note of the author's answer is applied, and each
declared departure is recorded with a reason.

1. **A block in a subdirectory of `docs/spec/` escapes the completeness test.**
   Where: `docs/spec/README.md:164-165`, `:153-154`, `:171-172`;
   `internal/tsv/block.go:17-26`. The sentences: "One test of `internal/tsv`
   [...] parses each block of `docs/spec/`", "so that no block escapes the
   test", and "A block that a later section adds makes that test fail until its
   name is in one of the two lists." The code: `ReadBlocks` reads the `.md` files
   at the root of the directory only, and skips an entry that is a directory
   (`block.go:25`). A probe with a block in `sub/a.md` gave no block and no
   error. Basis: the sentences hold for the layout that the README's own table
   "The files" gives (flat), and `docs/spec/` has no subdirectory today, so no
   block escapes at this head. A later section in a subdirectory would escape
   with no error. Fix: write "the Markdown files at the root of `docs/spec/`" in
   the README, or walk the tree. Class: **note**.

2. **"A byte-order mark" has no place.** Where: `docs/spec/README.md:93`;
   `internal/tsv/tsv.go:115-116`. The sentence: "The reader refuses a byte-order
   mark, a carriage return anywhere". The code refuses a byte-order mark at the
   start of the record only. A probe with U+FEFF inside a `text` field read it
   back as part of the value, with no error. Basis: the words "anywhere" and
   "at its start" are given for the carriage return and not for the byte-order
   mark, so the sentence has two readings. Fix: "a byte-order mark at its start".
   Class: **note**.

3. **The red evidence shows a stub run, not the order of the work.** Where:
   `runs/T-18v6/red.md:3-7`. The sentence: "each step ran its new tests before
   its code, and one commit holds the tests and the code of a step." The times:
   the step-1 red run is at 16:50:24Z, and the commit that holds the tests and
   the code of step 1 (`6a94efc`, 372 lines of code) is at 16:51:13Z, 49 seconds
   later; the step-2 red run is at 16:51:36Z and its commit (`c639eba`, 199
   lines) at 16:51:52Z, 16 seconds later; the step-3 red run is at 16:52:14Z and
   its commit (`ba906cc`) at 16:52:33Z. Basis: the record proves what gate step 3
   asks a reader to see, that each test fails on a stub for the right reason
   (the output names each assertion, and the reasons match the stubs). It does
   not prove that the tests existed before the code, and the intervals make that
   order improbable. Gate step 6 asks the record to say what happened, for
   example that the stub run was made with the code set aside. Class: **note**.

4. **The §12 Test cell is not filled; the departure is recorded.** Where: the
   issue's acceptance criterion 4, second half; `docs/prd/PRD-0001-layup.md:237`;
   `docs/tests/traceability.md:34`; the author's answer, note 6. The criterion:
   "`PRD-0001` §12 names its tests in the Test column." The NFR-001 Test cell
   stays `—`, and the diff does not touch the PRD. The answer declines the cell
   with a reason: the test proves the form of a record, not the §7.1 criterion.
   Basis: §12's own rule (`PRD-0001-layup.md:211-215`) is "Each delivering task
   fills the Test column", the Task cell of NFR-001 names `T-b97r, T-evad` and
   not `T-18v6`, and the plan's table "What phase 1 proves" names rows 9, 13,
   15, 16 and 20 for NFR-001. So the departure is sound and consistent with §12.
   But `traceability.md:34` says the same test covers NFR-001, so the two tables
   give two readings of one test, and `docs/tasks/T-18v6.md` does not name this
   departure. Fix at close-out: say in the task record that the §12 half of the
   criterion is declined, with the reason, so the ticked box has one reading.
   Class: **note**.

5. **The budget is over.** Where: the status comment; the diff. The maximum of
   the plan review: 1,200 lines added plus removed over 18 files against the
   base, close-out inside. The diff at the head: 1,451 lines (1,441 added, 10
   removed) over 14 files, before the review record, the resource record and the
   completed-log line. Basis: the status comment reports it to the Operator,
   whose decision it is. Noted, not judged. Class: **note**.

6. **A sentence about the present code sits inside the specification.** Where:
   `docs/spec/psb-check.md:26-30`. The sentence: "The present code does not yet
   give 2 for a `FILE` that is not valid UTF-8: it reads the bytes as they are,
   and can copy them into `excerpt` (task `T-5zmw` adds the check)." Verified
   true at the head: `internal/cli/cli.go:48-53` gives 2 only for a file that it
   cannot read, and `internal/psb/check.go:110-115` copies the bytes of the line
   into `excerpt` (it makes runes only when the line is over 80 runes). Basis:
   this is not a command whose exit code differs from the documented result
   without a record, because the document itself names the difference and the
   task that closes it, and the defect register (K32: settled by row 1, read by
   row 6) agrees. The sentence goes stale when `T-5zmw` lands; that task must
   remove it in the same change (R10). Class: **note**.

7. **The refusal of a misplaced fence is wider than its sentence.** Where:
   `internal/tsv/block.go:116-120`; `docs/spec/README.md:151-153`. The
   sentence: "The parser refuses a `tsv-schema` fence after a blockquote mark, a
   list marker, or four spaces or more." The code trims any run of the
   characters space, tab, `>`, `*`, `+`, `-`, `.`, `)` and the digits before the
   fence, so it also refuses a `tsv-schema` fence after a tab or after a date. A
   probe with the line `2026-10-01 ```tsv-schema a b` was refused with the
   list-marker message. Basis: each such case is a refusal and not an escape, so
   no block is lost; only the message can mislead. Class: **note**.

8. **No lesson is written back to `guardrails.md` §2.** Where: gate step 7; the
   plan review's note 11 (which counts a guardrails §2 lesson in the close-out);
   the diff (no change to `docs/guardrails.md`). The verification taught one
   that the next reader can hit (`runs/T-18v6/red.md`, step 4): `fs.Glob` drops
   the error of a directory that it cannot read, so a wrong directory gives no
   block and no error. Basis: a close-out item, not a defect of the head.
   Class: **note**.

9. **The close-out is still open at the head.** Where: `docs/tasks/backlog.md:37`
   (the task line is still in the backlog); `docs/tasks/completed.md` (no line);
   `docs/tasks/T-18v6.md:43` ("The Operator's answer goes here"); no resource
   record and no review record under `runs/T-18v6/` yet. Basis: the plan puts
   these in step 6, after this round; they are not findings against the head.
   Class: **note**.

### Acceptance criteria

1. **Each inventory item is delivered as its specification sections say** — met.
   `found-tsv`: `internal/tsv` reads and writes a record by its schema with the
   field rule, the empty mark, the 11 types, the key as the tuple of the `key`
   columns, and each refusal that README Records names; each of the item's seven
   open questions has a value "decided here" with its reason
   (`docs/spec/README.md:85-112`, `:141-144`, `:189-198`). `found-spec-schema-test`:
   the parser, the comparer and the completeness test are in, and the 14 blocks
   parse; zero blocks are compared with a Go schema, as condition 2(e) says and
   `docs/tasks/T-18v6.md:36-39` records.
2. **The tests of the plan pass** — met. At the head: `go build ./...`,
   `go vet ./...`, `go test -count=1 ./...`,
   `go test -count=1 -tags=integration ./...` and
   `go test -count=1 -tags=e2e -timeout 10m ./...` each exit 0;
   `git diff --check b48764f HEAD` exits 0; the six text checks exit 0 (with
   `HOME` set to a temporary directory). The integration level compares zero
   blocks (declared).
3. **K32 is settled in `docs/spec/`** — met. `docs/spec/README.md:62-69` gives
   the rule "decided here" with its reason, bounded as condition 3 asks;
   `docs/spec/psb-check.md:26-30` reads it and names `T-5zmw`; the register row
   K32 agrees (finding 6).
4. **The traceability rows are `green`, and `PRD-0001` §12 names its tests** —
   met in part. `docs/tests/traceability.md:34` holds the real test name
   (`TestEverySchemaBlockIsBuiltOrNotYetBuilt`) and `green`, and the test passes.
   The §12 Test cell of NFR-001 stays `—`, by the answer's note 6, a recorded
   departure with a reason that §12's own rule supports (finding 4).
5. **Tests cover the change and pass (R8)** — met. Each public function
   (`Read`, `Write`, `JoinList`, `ReadBlocks`, `Compare`) and each refusal that
   the README names has a test; the three levels pass; the red evidence shows
   each test fails on a stub for the right reason (finding 3 on what it does not
   show). Probes in a copy of the package under `/tmp/rv-T-18v6/` found no bug:
   the real `docs/setup/open-gaps.tsv` reads with the `open-gaps` schema (18
   rows, the two-column key holds), and the output of `layup psb check` on
   `F-0001` reads with the `psb-gaps` schema (1 row).
6. **Docs updated in the same PR** — met for the head: `docs/spec/README.md`,
   `packages.md`, `psb-check.md`, the task record, `traceability.md`, `red.md`
   and `test-runs.md` are in. The close-out items of findings 8 and 9 are still
   open, as the plan's step 6 says.
