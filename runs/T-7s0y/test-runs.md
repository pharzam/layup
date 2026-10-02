# T-7s0y: the test runs

The runs of `T-7s0y` (#86). Host: macOS, `go1.27.1`, `git` 2.54.0,
2026-10-02, UTC. "…" marks a cut part of an output.

## The red runs

1. **The unit tests of S01, of the gap table, of the outcome `input`, of the
   check of a done step's problem statement (D6), of the evidence call with its
   undo, and of the texts of S04,** on a skeleton (`Steps` gave the stubs; the
   reader of the gap table, the texts and the new runner parts gave nothing):
   thirteen tests fail, each for its rule.

   ```text
   $ go test -count=1 ./internal/setup/
   --- FAIL: TestTheGapTable                  (readGaps = [], <nil>; want [{Q-001 0 …} {Q-002 3 …}]; readGaps(""): no error, …)
   --- FAIL: TestS01StopsForEachMissingAnswer (no answer: not-active []; want a stop with the four S01- rows, Q-001 and Q-002 at line 3)
   --- FAIL: TestS01RefusesAnAnswerToNoGap    (an answer to Q-009: not-active "not built yet"; want an input error that names Q-009)
   --- FAIL: TestS01ChecksTheAnswers          (S01-stack rust, S01-name widget, acme/widget/x, -acme/widget, acme-/widget, ac--me/widget,
                                               acme/.., an OWNER of 40, a NAME of 101, S01-visibility internal, S01-baseline ssh://…,
                                               git@…:…, github.com/…: not-active; want fail with the reason; and seven good answers)
   --- FAIL: TestS01RowOfAFact                (a fact source: not-active []; want the row stack go fact F-0003#44)
   --- FAIL: TestS01UnchangedBrief            (S01 has no check of its inputs)
   --- FAIL: TestTheRunnerGivesTheInputErrorOfAStep    (the row S01 input …; want the input error and nothing written)
   --- FAIL: TestTheRunnerChecksTheInputsOfADoneStep   (calls ["S02" … "S15"]; want an input error before any step)
   --- FAIL: TestTheEvidenceOfAStep           (fails: done true, resets []; want done false, resets ["p"])
   --- FAIL: TestThePinText, TestTheDecisionRecordOfThePin, TestAddIndexRow, TestTheAnswersRecord   (the empty text)
   ```

2. **The unit tests of S02, S03, S04 and the commands of S03,** with a fake of
   `internal/git` and of the files of the work area, and a fixed clock, on a
   skeleton of the four functions (each step `not-active`, `not built yet`; no
   command): five tests fail.

   ```text
   $ go test -count=1 -run 'TestS02|TestS03|TestTheCommandsOfS03|TestS04' ./internal/setup/
   --- FAIL: TestS02   (S02: not-active "not built yet" [], calls []; want done, the four pin rows with pin.time 2026-10-02T09:30:00Z, and
                        the calls remove, ls-remote, clone, checkout, rev-parse, remove .git, rename; each failed call: want fail with its first line;
                        a target that exists: want an input error and no call)
   --- FAIL: TestS03   (want init, the commit by layup-agent[bot] at pin.time, rev-parse; another tree: want fail; a target with a history: …)
   --- FAIL: TestTheCommandsOfS03   (commandsS03 = []; want the remote origin and the push of main)
   --- FAIL: TestS04   (want done, the three value rows and the branch made; S04 wrote 0 files, want 6)
   --- FAIL: TestS04OnABranchThatExists
   ```

3. **The integration and e2e tests,** on a skeleton of `Steps` (each step a
   stub, `not built yet`): each fails at S01.

   ```text
   $ go test -count=1 -tags=integration -run '…' ./internal/setup/
   --- FAIL: TestS01ToS04OnABaseline, TestALoginURLFailsWithNoPrompt, TestAStepThatStoppedInItsMiddle, TestTheUndoOfAnEvidenceThatFails
             (the run: [{S01 layup-setup not-active not built yet} {S02 … not run: S01 did not pass} …])
   $ go test -count=1 -tags=integration -run 'TestSetupRunsS01ToS04|TestABrokenPinFailsTheEvidenceOfS04' ./internal/cli/
   --- FAIL: TestSetupRunsS01ToS04   (the first run: exit 1 …; want 3 and the stop table: the four S01- rows, then
                                      Q-001 (G1, line 0, —) and Q-002 ("fast", inputs/briefs/problem-statement.md:3))
   --- FAIL: TestABrokenPinFailsTheEvidenceOfS04   (the run to S03: exit 1, S01 not built yet)
   $ go test -count=1 -tags=e2e -timeout 10m -run TestSetupRunsS01ToS04 ./cmd/layup/
   --- FAIL: TestSetupRunsS01ToS04   (a new work area: exit 1 …; want 3 and the stop table)
   ```

4. **The unit tests of `internal/cli`** (`TestSetupReadsTheProblemStatement`,
   `TestTheEvidenceCall`) came after their code, in the same step; their
   evidence is the mutations M22 to M24 below.

5. **The message of S03's own commit** (D8): a self-review against the plan
   found that S03 took a root commit with the tree `pin.tree` and any message.
   The new case "another message" of `TestS03`, and the case "one commit with
   the tree pin.tree and another message" of `TestAStepThatStoppedInItsMiddle`
   (real `git`), are the mutation M26 below.

## The mutations

Each mutation is one change of the code, run against the tests of its package
at the named levels; the file is written back from its text before the
change (the change is not committed yet), and `cmp` showed each file unchanged
after the run. All 26 are detected.

| Mutation | Detected by |
| -------- | ----------- |
| M1 S01 asks no `Q-` question | `TestS01StopsForEachMissingAnswer`; `TestSetupRunsS01ToS04` (cli, integration) |
| M2 a gap at line 0 gets a line | the same two |
| M3 `OWNER/NAME` takes a third part | `TestS01ChecksTheAnswers` |
| M4 a visibility of another value passes | `TestS01ChecksTheAnswers` |
| M5 a fact source is an answer | `TestS01RowOfAFact` |
| M6 a changed problem statement passes | `TestS01UnchangedBrief` |
| M7 S02 keeps the `.git` of the clone | `TestS02`; `TestS01ToS04OnABaseline`, `TestAStepThatStoppedInItsMiddle`, `TestTheUndoOfAnEvidenceThatFails` |
| M8 S02 replaces a target that exists | `TestS02`; `TestAStepThatStoppedInItsMiddle` |
| M9 S03 commits at the time of the clock | `TestS03` |
| M10 S03 makes a second root commit after a stop | `TestS03`; `TestAStepThatStoppedInItsMiddle` |
| M11 S04 takes the highest number | `TestS04`; `TestSetupRunsS01ToS04` (cli, integration: `adr-lint` and the files) |
| M12 S04 dates the records by the clock | `TestS04` (`pin.time` 2026-09-30, the clock 2026-10-02) |
| M13 S04 takes `layup-setup` at any commit | `TestS04OnABranchThatExists` |
| M14 the row `_none yet_` stays | `TestAddIndexRow`, `TestS04` |
| M15 the answers record keeps the angle quotes | `TestTheAnswersRecord` |
| M16 the remote of S03 is not quoted | `TestTheCommandsOfS03`; `TestSetupRunsS01ToS04` (cli) |
| M17 the runner does no undo | `TestTheEvidenceOfAStep`; `TestTheUndoOfAnEvidenceThatFails`; `TestABrokenPinFailsTheEvidenceOfS04` |
| M18 the undo runs after the head moved | `TestTheEvidenceOfAStep` |
| M19 the evidence call does not run | `TestTheEvidenceOfAStep`; `TestTheUndoOfAnEvidenceThatFails`; `TestABrokenPinFailsTheEvidenceOfS04` |
| M20 an input error of a step is a row | `TestTheRunnerGivesTheInputErrorOfAStep`; `TestAStepThatStoppedInItsMiddle` |
| M21 the check of a done step's problem statement does not run | `TestTheRunnerChecksTheInputsOfADoneStep` |
| M22 `internal/cli` hands the SHA-256 of the table | `TestSetupReadsTheProblemStatement` |
| M23 `internal/cli` runs the steps with no problem statement | `TestSetupReadsTheProblemStatement` |
| M24 the evidence call reads `clear` as a fail | `TestTheEvidenceCall` |
| M25 S04 gets no evidence call | `TestTheUndoOfAnEvidenceThatFails`; `TestABrokenPinFailsTheEvidenceOfS04` |
| M26 S03 takes a commit with another message | `TestS03`; `TestAStepThatStoppedInItsMiddle` |

## A measurement: a URL that asks for a login

`TestALoginURLFailsWithNoPrompt` serves `401` with `WWW-Authenticate: Basic`
from `httptest`. S02 fails at once, with the reason (one run):

```text
git ls-remote http://127.0.0.1:53469/target.git HEAD: git ls-remote --exit-code -- http://127.0.0.1:53469/target.git HEAD:
exit status 128: fatal: could not read Username for 'http://127.0.0.1:53469': terminal prompts disabled
```

## The green runs

On the tree of the commit before the one that adds this file (`0c00064`; this
commit adds only text), 21:35:57Z to 21:37:49Z:

| Command | Result |
| ------- | ------ |
| The local checks of `AGENTS.md`, with `git diff --check 94d1d71 HEAD`; `go build ./...`; `go vet` with no tag, `-tags=integration` and `-tags=e2e`; `gofmt -l .` | exit 0; no file |
| `go test -count=1 ./...` | `ok` × 11 packages |
| `go test -count=1 -tags=integration ./...` | `ok` × 11 (`internal/setup` 4.5 s, `internal/cli` 4.4 s) |
| `go test -count=1 -tags=e2e -timeout 10m ./...` | `ok` × 11 (`cmd/layup` 6.7 s) |
| `go test -race -count=1` on `internal/catalog`, `internal/verify`, `internal/work`, `internal/standin`, `internal/setup` and `internal/cli`, and on `internal/catalog`, `internal/verify`, `internal/setup` and `internal/cli` with `-tags=integration` | `ok` |
| `sh docs/setup/tests/run.sh`; `sh docs/tests/run-discipline-tests.sh` | 44 passed, 0 failed; 81 passed, 0 failed |
| The first commit (`1ec7e7e`) alone, in a scratch work tree: `go vet` with each tag; the integration tests of `internal/cli` and `internal/setup`; `TestSetup` (e2e) | `ok` |

## The fix of review round 1 (cycle 1)

Round 1 (Claude Fable 5.1, on `5dadab0`; the record is on #86) gave
`material`: finding 1, the `source` of an answer kept its angle quotes in the
answers record, against `setup.md`. Note 2 (b) is one of the three defects
that the author found in a self-review before the record came; the fix holds
the three. Notes 2 (a) and 5 are known limits in `setup.md`, note 3 is text
(the order of the checks of a done step's inputs), and notes 4 and 6 to 9
confirm the change. The red runs, on the
code of `5dadab0` with the new cases (`\u2039` and `\u203a` are the two angle
quotes, written as escapes so that check `markers` reads no marker in this file):

```text
$ go test -count=1 -run TestTheAnswersRecord ./internal/setup/          (finding 1)
    answersRecord = … 5. `Q-001` Go &lsaquo;1.26&rsaquo; — by idea-owner; source said \u2039here\u203a; the question: …
    want            … 5. `Q-001` Go &lsaquo;1.26&rsaquo; — by idea-owner; source said &lsaquo;here&rsaquo;; the question: …
$ go test -count=1 -run 'TestAddIndexRow|TestS04$' ./internal/setup/     (a file that ends with no line feed; note 2 (b))
    a table at the end of a file with no line feed:
        "… | [0001](0001-a.md) | A | Accepted || [0002](0002-b.md) | B | Accepted |\n"
    a list of hashes in the root commit, "abc  docs/facts/x.md": done "abc  docs/facts/x.md7fb3702e…  docs/facts/F-0001-setup-answers.md\n"
$ go test -count=1 -run TestS04OnABranchThatExists ./internal/setup/      (a target on main whose layup-setup exists)
    main, and layup-setup exists: input "target of the work area is on refs/heads/main, not on main or layup-setup" …
```

The mutations of the fix, each detected:

| Mutation | Detected by |
| -------- | ----------- |
| M27 the answers record keeps the angle quotes of a source | `TestTheAnswersRecord` |
| M28 a table at the end of a file joins the new row | `TestAddIndexRow` |
| M29 a list of hashes with no last line feed joins the new line | `TestS04` |
| M30 a target on main with the branch `layup-setup` gets the old reason | `TestS04OnABranchThatExists` |

The green runs, on the tree of the fix (the code and the tests of the fix
commit), 21:52:07Z to 21:53:58Z: the local checks of `AGENTS.md`; `go build`;
`go vet` with each tag; `gofmt -l .` (no file); the three test levels (`ok` ×
11 each); the race tests of the six packages, and of four with
`-tags=integration`; `run.sh` (44 passed) and the discipline tests (81
passed): each exit 0.
