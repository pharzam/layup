# T-cav3 compiled mutants

Rerun for round 1 (finding 2); every figure is from this rerun. Each mutant is one source edit of a copy of
the worktree (outside the tree), built (`go build ./...`) and vetted (`go vet ./internal/run/`); a failure of
either counts as not run (none occurred). The tests run in three stages, stopping at the first failure: the
unit level, the tree half of the integration level (with the tree-item tests), then the demo. `test` is the
exit status of the stage that failed (1 = killed; 0 would be not killed). The file is restored from a backup and
compared by SHA-256 (12-character prefixes below; `before` = `after` for every row, so each restore is
byte-exact). After the last mutant the unit and integration controls exit 0. The SHA-256 of the five files
(`specinput.go`, `treeitems.go`, `batch.go` and the two test files) as one `sha256sum` list is `139709889710c0c11b75df886b0cadc29a210396b153e0e1a62f82ff462e1c61`, equal in the
copy and in the worktree.

| # | build | vet | test | failing assertion (first line) | sha256 before | sha256 after |
| - | ----- | --- | ---- | ------------------------------ | ------------- | ------------ |
| M1 | 0 | 0 | 1 | specinput_test.go:124: the fields that the records do not hold: {Changed:[] PSBRecord | c375c220d6ec | c375c220d6ec |
| M2 | 0 | 0 | 1 | specinput_test.go:284: isNumberedFacts("F-0007-setup-answers.md") = true; want false | c375c220d6ec | c375c220d6ec |
| M3 | 0 | 0 | 1 | specinput_test.go:133: spec/spans.tsv: the input of the full spec check: Changed from | c375c220d6ec | c375c220d6ec |
| M4 | 0 | 0 | 1 | specinput_test.go:112: the input of the full spec check: Index from records:spec/inve | c375c220d6ec | c375c220d6ec |
| M5 | 0 | 0 | 1 | specinput_test.go:150: spec/inventory-2.tsv: the input of the full spec check: Change | c375c220d6ec | c375c220d6ec |
| M6 | 0 | 0 | 1 | specinput_test.go:224: Register "", the input of the full spec check: Register from r | c375c220d6ec | c375c220d6ec |
| M7 | 0 | 0 | 1 | specinput_test.go:120: the input: PSB 100 bytes, Spans [[0 10 1 need  s1] [10 20 2 me | c375c220d6ec | c375c220d6ec |
| M8 | 0 | 0 | 1 | specinput_integration_test.go:58: Changed []; want the paths of the diff [docs/facts/ | c375c220d6ec | c375c220d6ec |
| M9 | 0 | 0 | 1 | specinput_integration_test.go:52: SpecFiles []; want only a.md with its bytes | c375c220d6ec | c375c220d6ec |
| M10 | 0 | 0 | 1 | specinput_integration_test.go:55: PRDs [PRD-0001-x.md PRD-0002- | c375c220d6ec | c375c220d6ec |
| M11 | 0 | 0 | 1 | specinput_test.go:120: the input: PSB 0 bytes, Spans [[0 10 1 need  s1] [10 20 2 meas | c375c220d6ec | c375c220d6ec |
| M12 | 0 | 0 | 1 | specinput_integration_test.go:244: the copy of the bet: Copy "", Index [[1 8b866fbe94 | c375c220d6ec | c375c220d6ec |
| M13 | 0 | 0 | 1 | specinput_test.go:270: base: the input of the full spec check: Changed from target:di | c375c220d6ec | c375c220d6ec |
| M14 | 0 | 0 | 1 | specinput_test.go:150: tasks/T-bbbb/events.tsv: the input of the full spec check: Cha | c375c220d6ec | c375c220d6ec |
| M15 | 0 | 0 | 1 | specinput_test.go:301: psbRecordsOf([F-0007-numbered-facts.md F-0009-numbered-facts.m | c375c220d6ec | c375c220d6ec |
| M16 | 0 | 0 | 1 | specinput_integration_test.go:52: SpecFiles [a.md l.md]; want only a.md with its byte | c7d726dafcd3 | c7d726dafcd3 |
| M17 | 0 | 0 | 1 | specinput_integration_test.go:52: SpecFiles [a.md b.md]; want only a.md with its byte | c7d726dafcd3 | c7d726dafcd3 |
| M18 | 0 | 0 | 1 | specinput_test.go:334: an empty record that is there: the records have no empty; want | 82886162920e | 82886162920e |
| M19 | 0 | 0 | 1 | specinput_test.go:340: a record its reader refuses: line 2: a bad row; want the path  | 82886162920e | 82886162920e |
| M20 | 0 | 0 | 1 | specinput_test.go:284: isNumberedFacts("xF-0007-numbered-facts.md") = true; want fals | c375c220d6ec | c375c220d6ec |

The mutants: M1 `AnswersRecord` taken from the `answers.record` row (F-0001); M2 any `F-NNNN-*.md` counts as the numbered facts record; M3 a missing `spec/spans.tsv` gives no error; M4 a missing `spec/inventory.tsv` is an error; M5 a missing version is skipped; M6 `Register` read with `ReadRegister`; M7 any event counts as merged; M8 `Changed` is nil; M9 `SpecFiles` read at `base`; M10 the PRD prefix check dropped; M11 `PSB` not set; M12 `Copy` ignored; M13 the SHA check skipped; M14 a refused `events.tsv` ignored; M15 two numbered facts files give the first; M16 `filesAt` ignores the mode; M17 `filesAt` reads below the directory; M18 `readRequired`: an empty file counts as missing; M19 `readRequired`: the reader's error is not wrapped with the path; M20 the numbered facts name matches as a suffix.

M18 and M19 survived the first run of the harness (no test had a present empty record or a path in a reader error), and the test `TestReadRequiredKeepsTheStoreError...` was added; in this rerun both are killed. M10 now removes the `HasPrefix` of `isPRDName`. No mutant is claimed equivalent; M12 is reached only by the demo.
