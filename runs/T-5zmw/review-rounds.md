# T-5zmw: the review records

The two review rounds of `T-5zmw` (#83), copied from the issue as they stand there.

## Review record — round 1

| Field | Value |
| ----- | ----- |
| Commit reviewed | `edbeb6095ae85a8e351e9bcea9e6f5c6d8f0036a` |
| Reviewer | GPT-6 Sol (xhigh) on Devin CLI |
| Lens | correctness and acceptance criteria |
| Briefed on | `.review-in/brief.md`, issue body, ordered issue comments and O-131, six-item inventory; row 6 and K14/K32 of the plan; the whole diff against `4b47982`; `docs/spec/psb-check.md`, the named parts of `docs/spec/README.md`, `packages.md`, S01/S04/S06 of `setup.md`; test levels, traceability, PRD §§12–13, engineering discipline (Bootstrap mode and Testing), guardrails §2, `runs/T-5zmw/test-runs.md`, `docs/tasks/T-5zmw.md`, and the code and tests of `internal/psb`, `internal/tsv`, `internal/cli` and `cmd/layup`. |
| Barred from | nothing; the comments of #83 are the plan and its record |
| Independence claimed | Fresh read-only reviewer session, not the author; the author of all branch commits is Claude Opus 5.5, a different model. I read the published plan review as part of the contract, not an earlier review of this head. |
| Cycle | 0 |
| Verdict | `material` |

### Raw findings

1. **Material — a whitespace-only stack value incorrectly passes G1.** `docs/spec/psb-check.md:52,58-61` and O-131 (`.review-in/issue-comments.md:139-141`) require a value that is neither white space nor `*`; `internal/psb/check.go:40` uses Go regexp `[^\s*]`. Go's `\s` does not include U+00A0, U+2003 or vertical tab. With a file containing only `Technology stack: `, U+00A0 and LF, the built command gives exit **0** and only the header; the specified result is exit **1** with a G1 row. U+2003 and vertical tab give the same wrong result. `internal/psb/check_test.go:63-74` tests ASCII empty labels but none of these whitespace values. The example expression in O-131 itself has this limit; settle the intended reading without silently changing the Operator's decision.

2. **Material — the promised `edge` golden was not extended.** The scope says to add the unsettled cases to `edge.md` (`.review-in/issue-body.md:23`, `.review-in/inventory.md:49-56`), and approved D3 says to pin them in `edge.md` **and** `edge.tsv` (`.review-in/issue-comments.md:19`). Neither file changed against `4b47982`. Instead, `docs/spec/psb-check.md:63-66` now permits `TestEdgeCases` in place of that golden, and `docs/tasks/T-5zmw.md:74-79` records this as the author's own plan change. The Go cases do exercise the behavior, but the author's answer to the plan review did not amend D3's named deliverable (`.review-in/issue-comments.md:89-106`). Whole-file cases can have separate tests without silently dropping the named golden for the other cases. This is a scope/authority change, not a claim that the Go cases fail.

3. **Material — D4 does not pin every G4 word in mixed case.** D4 expressly calls for *each* word of G4 in mixed case (`.review-in/issue-comments.md:20`); the changed specification says the `values` golden does this (`docs/spec/psb-check.md:81-85`). But `internal/psb/testdata/values.md:14-20` uses mixed case for five words and only lowercase `fast`; `internal/psb/testdata/triggers.md:3` also uses lowercase `fast`. The unit golden therefore cannot catch a regression that makes `Fast` fail while lowercase `fast` still works. A manual run with `This is Fast.` gives the correct G4 row, but that behavior has no promised golden case. The inventory's earlier five-word list is narrower than D4's reviewed, later decision.

4. **Note — qualify a base-green statement in the evidence.** `runs/T-5zmw/test-runs.md:9-16` includes “the record check of each golden” among tests that “pass on the base”; the same file at `:46-51` correctly reports that its integration test does not compile before `GapsSchema` exists. This is an overbroad summary, not a missing red run: the detailed red result is present. State that the *golden contents* characterize the base, not that the new record/schema test passes there.

### Acceptance criteria

| Item | Status and basis |
| ---- | ---------------- |
| #83: each of six inventory items delivered as specified | **Not met:** G1 whitespace fails and D3's named `edge.md`/`edge.tsv` delivery is absent (findings 1–2). Other items are implemented. |
| #83: planned unit, integration and E2E tests pass, with each rule value and edge | **Not met:** all executed test levels pass, but the promised mixed-case `Fast` golden and the `edge` golden are absent (findings 2–3). |
| #83: traceability rows green and PRD §12 Test cells filled | **Met:** the new and moved tests have real names, levels and `green` status; the four specified Test cells and §13 row were changed. |
| #83: tests cover change and pass (R8) | **Not met:** commands pass, but G1 accepts whitespace-only values without a test, contrary to O-131 (finding 1). |
| #83: documents updated in same PR | **Not met:** the changed rule's reading is not matched by code, and the changed description of the D4 golden overstates its coverage. |
| D1 / D2 | **Met:** schema and `tsv.Write` produce the unchanged old goldens and propagate write errors; invalid UTF-8 is rejected before checking with a first-byte line diagnostic. |
| D3 / D4 | **Not met:** the cases run, but the D3 fixture contract and one D4 value case remain open (findings 2–3). |
| D5 / D6 | **Met:** unit goldens are embedded, file-using CLI cases moved to integration, and E2E covers the real file, repeat, clean file, input errors and empty environment. |
| D7 / D8 | **Met for row 6:** `psb-check.md` states S04 under O-124; the S04/S06 step-table changes are assigned to row 9 by K14 and plan-review note 3. Test rows and PRD cells were updated. |
| Plan-review conditions 1–4 | **1: authority obtained, semantics not fully met** (finding 1); **2–4: met** (write error contract/tests, schema and golden record checks, truly empty environment with closed stdin). |
| Plan-review notes 1–5 | **1, 3–5: met** (only G1 sought the Operator; K14 work stayed with later rows; levels and names updated; real PSB remains integration). **2: qualified** by evidence note 4; no baseline-green characterization is reported as a red test. |
| O-131 | **Not fully met:** the ASCII empty-label cases pass, but whitespace-only labels can still suppress G1 (finding 1). |
| Close-out | **At close-out:** boxes, verdict, resource record and completed-log move are not in this frozen head. |

### Runs

Every shell command used `HOME=<scratch>/home`; Go commands used `GOCACHE=/Users/farzam/Library/Caches/go-build`, and temporary inputs and the binary were under `<scratch>` (`rv83-scratch`), not the repository. No external network or forge command was run.

| Command / case | Result |
| -------------- | ------ |
| `git diff --shortstat 4b47982 HEAD` | `21 files changed, 768 insertions(+), 114 deletions(-)`: 882 lines added plus removed, within 1,250 lines and 28 files. |
| `go test -count=1 ./...` | Exit 0; seven packages passed. |
| `go test -count=1 -tags=integration ./...` | Exit 0; seven packages passed. |
| `go test -count=1 -tags=e2e -timeout 10m ./...` | Exit 0; seven packages passed. |
| `go test -race -count=1 ./internal/psb/ ./internal/cli/`; same with `-tags=integration` | Exit 0 for both commands and both packages. |
| `go vet ./...`; `go vet -tags=integration ./...`; `go vet -tags=e2e ./...`; `gofmt -l .` | Exit 0 for each; no Go file listed by `gofmt`. |
| `go build -o <scratch>/layup ./cmd/layup` | Exit 0. |
| Built binary: invalid bytes at start/end/after CRLF, overlong `\xc0\x80`, surrogate half | Each exit 2, empty stdout, file and correct first-invalid-byte line on stderr (line 1 or 2). |
| Built binary: BOM before first `# Terms`; zero-byte file; only lone-CR line ends; long excerpt with `é猫` at character 80 | Exit 0 with no G3 for the BOM case; exit 1 with G1 for empty input; exit 1 with one sanitized G4 excerpt on line 1 for lone-CR ends; exit 1 with an 80-character G4 excerpt retaining `é猫`. |
| Built binary: tab, lone CR, or pipe character inside a G4 excerpt; NUL | Exit 1; tab and CR each become one space, the pipe stays literal; NUL remains in the UTF-8 excerpt (no stated rule forbids it). |
| Built binary: ASCII G1 empty bold/star labels and tab before colon; bold/star labels with value and uppercase label with value | Empty labels and tab give exit 1 G1; named values give exit 0. |
| Built binary: `Technology stack: ` followed only by U+00A0, U+2003, or vertical tab | **Exit 0**, header only; expected exit 1 G1 under the white-space wording (finding 1). |
| Built binary: directory, symlink to directory, unreadable file (`000` mode); `./-x` | Each unreadable input exits 2 with no stdout; `./-x` resolves as a file and exits 0. |
| Built binary: stdout open for reading only; same real PSB twice; empty environment and null stdin | Output failure exits 2 with diagnostic; the two normal runs have equal output and code 1; the no-environment run matches them. |
| Built binary: a file containing only `—`; a mixed-case `Fast` line | `—` alone gives a valid G1 row (no `tsv.Write` rejection); `Fast` gives a G4 row and exit 1. Code inspection of `internal/psb/check.go:55-126` and `internal/tsv/tsv.go:63-99` found no reachable schema/field rejection for a UTF-8 CLI input; `—` is allowed in a non-key text field. |
| `sh docs/prd/prd-lint.sh`, `sh docs/links/link-lint.sh`, `sh docs/adr/adr-lint.sh`, `sh docs/setup/setup-check.sh`, `sh docs/tests/run-discipline-tests.sh`; `git diff --check 4b47982 HEAD` | All exit 0; 81 discipline tests pass, no whitespace error. |

## Review record — round 2

| Field | Value |
| ----- | ----- |
| Commit reviewed | `b6b52d5abd8fd7e16a864b5aaa4f0647e2ae1a61` |
| Reviewer | GPT-6 Sol (xhigh) on Devin CLI |
| Lens | the fixes of round 1, and correctness of what they touch |
| Briefed on | `.review-in/brief.md`, `issue-body.md`, the ordered plan, plan review, author's answer and O-131 in `contract.md`, `round-1.md`, `inventory.md`; the full fix diff and the branch diff for touched files; `docs/spec/psb-check.md`, `docs/spec/README.md` (Commands, Records), `docs/tests/test-levels.md`, `docs/tests/traceability.md`, `docs/prd/PRD-0001-layup.md` §§12–13, `docs/guardrails.md` §2, `runs/T-5zmw/test-runs.md`, `docs/tasks/T-5zmw.md`, the code and tests of `internal/psb`, and the related CLI, writer and E2E code. |
| Barred from | the comments of #83 after the record of O-131 |
| Independence claimed | Fresh read-only reviewer session, not the author; the author of the branch is Claude Opus 5.5, a different model. I read round 1 only as the assigned input and did not use later issue comments or the network. |
| Cycle | 1 |
| Verdict | `nothing material in scope` |

### Raw findings

None. No material finding or note remains in this lens.

### The findings of round 1

1. **Fixed.** `internal/psb/check.go:37-42` excludes Unicode `White_Space` from the value and permits it as a prefix after `:` except for tab, form feed and carriage return. This matches `docs/spec/psb-check.md:52,58-72`: O-131 requires a non-white-space, non-`*` value; the extra prefix reading keeps the old result when such a value follows Unicode white space. The Operator decided the empty-label result (`.review-in/contract.md:115-141`); the author did not reverse that choice. All 30 built-binary boundary inputs gave the specified exit and table, including the seven named white-space characters alone and before `Go`, tab/form feed/carriage return, mixed `*` and spaces, U+00A0 before `:`, U+200B and the pinned template forms. A separate Go check found zero differences between the expression's white-space class and `unicode.IsSpace` over all runes; its prefix class also had zero differences from the documented rule. The red run for U+00A0 and the code before the fix agree (`runs/T-5zmw/test-runs.md:115-151`).
2. **Fixed.** `internal/psb/testdata/edge.md:14-29` and `edge.tsv:3-4` now pin the shared-file cases: a bold header and an escaped pipe must give no G2 row; the two-column table must give a G2 row only for line 21, so choosing the wrong column changes the golden; the fenced line must give a G4 row. G1 needs a file without the named stack already in `edge.md`; G3 needs separate whole-file heading and terms-table states. Those cases remain in `internal/psb/check_test.go:66-105`, including the fenced heading's G3 effect. This division retains D3's cases and makes each shared case distinguish the intended reading.
3. **Fixed.** `internal/psb/testdata/values.md:14-19` now has `Fast`, as well as the other five G4 words in case-varied forms. `values.tsv:7-12` expects a row for each and the lowercase question word, including `fast` on line 19.
4. **Note applied.** `runs/T-5zmw/test-runs.md:9-19` says that the golden *contents* characterize the base, but `TestEveryGoldenIsARecordOfTheBlock` could not compile there before `GapsSchema` existed. The red runs and the two lessons in `docs/guardrails.md:355-372` agree with the inspected code and the reproduced output-error and whitespace cases. `docs/tasks/T-5zmw.md:88-111` records the fixes. The remaining close-out (verdict, resource record, logs and review copies) follows this review and is not in this frozen head.

### Runs

Every command used `HOME=<scratch>/home`; Go commands used `GOCACHE=/Users/farzam/Library/Caches/go-build` and disabled module network access. The built binary and probe inputs were in `<scratch>` (`rv83b-scratch`).

| Command | Result |
| ------- | ------ |
| `git rev-parse HEAD` | `b6b52d5abd8fd7e16a864b5aaa4f0647e2ae1a61`. |
| `go build -o <scratch>/layup ./cmd/layup` | Exit 0. |
| Built binary on 30 G1 boundary inputs | Each exit code and stdout table matched the rule; stderr empty. Seven Unicode white-space values alone gave exit 1 and one G1 row; each followed by `Go` gave exit 0. Tab, form feed and carriage return before `Go` gave exit 1; U+200B as value gave exit 0. |
| Scratch Go expression probe over every rune | Exit 0; white-space class and after-colon prefix class each had zero mismatches with `unicode.IsSpace` and the stated exceptions. This independently reproduces one of the two measurements in `test-runs.md:142-150`, not its 465,010-line comparison. |
| Built binary on the real PSB with closed pipe; then read-only stdout | Closed pipe: signal 13 (`returncode -13`), empty stderr. Read-only stdout: exit 2, `layup: write /dev/stdout: bad file descriptor`; matches `test-runs.md:94-108`. |
| `go test -count=1 ./...` | Exit 0; seven packages pass. |
| `go test -count=1 -tags=integration ./...` | Exit 0; seven packages pass. |
| `go test -count=1 -tags=e2e -timeout 10m ./...` | Exit 0; seven packages pass. |
| `go test -race -count=1 ./internal/psb/ ./internal/cli/` | Exit 0; both packages pass. |
| `go test -race -count=1 -tags=integration ./internal/psb/ ./internal/cli/` | Exit 0; both packages pass. |
| `go vet ./...`; `go vet -tags=integration ./...`; `go vet -tags=e2e ./...` | Exit 0 for each. |
| `gofmt -l .` | Exit 0; no file listed. |
| `sh docs/adr/adr-lint.sh`; `sh docs/prd/prd-lint.sh`; `sh docs/links/link-lint.sh` | Exit 0 for each; links: 1,318 resolved. |
| `sh docs/tests/run-discipline-tests.sh`; `sh docs/setup/setup-check.sh` | Exit 0 for each; 81 discipline tests pass, setup checks pass. |
| `git diff --check 4b47982 HEAD` | Exit 0; no whitespace errors. |
| `git diff --shortstat 4b47982 HEAD` | 23 files, 893 insertions and 114 deletions: 1,007 of 1,250 lines and 23 of 28 files. Close-out has 243 lines and five files of headroom; it has not yet happened. |
