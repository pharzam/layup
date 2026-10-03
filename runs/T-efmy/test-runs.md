# T-efmy: the test runs

The runs of `T-efmy` (#96). Host: macOS, `go1.27.1`, `git` 2.54.0,
2026-10-03, UTC. "…" marks a cut part of an output.

## The red runs of the release check

[`release-check.sh`](release-check.sh) is the deterministic part of the release
review (D3 of the plan, with condition 1 and notes 2 and 4 of its plan review).
Three mutation runs, each on its own commit of a scratch clone of `621af09`, so
that check (1) (the commit and a clean tree) passes and each run fails for its
own reason. Each run gives exit 1:

1. **An import of `net/http`** in a non-test file (`internal/psb/zz_mutation.go`,
   `import _ "net/http"`), at 04:25:00Z: `TestPackageRules` fails, naming its
   findings, and the check gives `FAIL: TestPackageRules`:

   ```text
           rule 5: cmd/layup depends on crypto/tls
           rule 5: cmd/layup depends on net
           rule 5: cmd/layup depends on net/http
           … (the same three for internal/cli and internal/psb)
   FAIL: TestPackageRules
   == FAIL
   ```

2. **A program other than the row's** (`internal/gate/zz_mutation.go`,
   `exec.Command("curl", "https://example.com")`), at 04:25:02Z:

   ```text
           Starts a program: internal/gate/zz_mutation.go:5 starts curl; its row says sh
   FAIL: TestPackageRules
   == FAIL
   ```

3. **A `git` verb that reaches a remote** (`internal/git/zz_mutation.go`,
   `func Push(dir string) error { return do(dir, "push") }`), at 04:25:04Z.
   `TestPackageRules` passes, as `push` is a call of `git` in `internal/git`;
   only check (4) fails:

   ```text
   == (4) The git verbs that reach a remote, in the non-test Go files of internal/git
     internal/git/zz_mutation.go:4:func Push(dir string) error { return do(dir, "push") }
   FAIL: a git verb that reaches a remote
   == FAIL
   ```

## The green run of the release check

At `621af09`, in a fresh clone with `origin` removed, 04:25:16Z to 04:25:17Z:
`sh release-check.sh 621af09bce4e3ef897a9fa655cabb59229f700ba` gives exit 0
and `== PASS`. Its whole output is
[`release-check.txt`](release-check.txt): `TestPackageRules` passes; `go list
-deps ./...` names only the packages of this module out of the standard
library, and none of `net`, `net/http` and `crypto/tls`; three calls start a
program (`git` in `internal/git/git.go:81`, `sh` in
`internal/gate/scratch.go:190` and `internal/verify/scripts.go:61`); no string
literal of a `git` verb that reaches a remote is in `internal/git` (the comment
"no remote helper" on line 65 and the literal `"ls-remote"` on line 129 do not
match, condition 1 of the plan review); the binary embeds `internal/catalog/go`,
whose two active kinds run `gofmt` with `go vet`, and `go test`; and
`internal/verify` runs two scripts of the baseline, `run-discipline-tests.sh`
and `link-lint.sh`.

## The release review (the uat)

The release review of D4 ran in the reviewer order. Devin (its usage quota,
04:26:30 to 04:27:44) and OpenCode (Grok 4.7: no output in five minutes,
04:27:49 to 04:32:49) gave no record. Claude Fable 5.1 (effort `xhigh`, on the
Claude Code CLI, a fresh read-only session in a clone at `621af09`, from
04:32:54Z, its record at 4 min 22 s) recorded `holds` for `REQ-015` and for
`REQ-017`, with ten notes, and accepted the scenario. It ran the release check
again from a scratch path (exit 0, the same lines), and read the 36 non-test Go
files in full. Its record is on #96 (comment 5965588985), and the same text is
[`release-review.md`](release-review.md).

Its heading, `## Release review — phase 1`, is not one that
`review-record-lint` reads. The lint reads only the round records of this
task's pull request, and those meet the inventory's test "The review record
passes review-record-lint in CI" (note 6 of the plan review).

## The freeze

On the frozen head `5eb1794`, 04:38:54Z to 04:41:16Z, the 18 steps of the
ladder, each exit 0: the eight local checks of `AGENTS.md` (with `git diff
--check` against `origin/main`), `go build`, `go vet` with each tag, `gofmt
-l` (no file), the three test levels, the two runs with `-race`, and the
harness of the fixtures of `setup-check.sh` (44 passed). CI of PR #119 passed
its job `tests` on it before review round 1.

## The notes of round 1 (close-out)

- **Note 1, a raw string literal.** Round 1 ran `` do(dir, `push`) `` in
  `internal/git` through the check: exit 0, as check (4) matched only the form
  `"push"`. The pattern of check (4) now takes a raw string literal too. The
  same mutation, on its own commit of a scratch clone of `621af09`, at
  16:12:22Z:

  ```text
    internal/git/zz_mutation.go:4:func Push(dir string) error { return do(dir, `push`) }
  FAIL: a git verb that reaches a remote
  == FAIL
  ```

  The changed check at `621af09`, in a fresh clone, 16:12:35Z to 16:12:37Z:
  exit 0 and `== PASS`, with the same lines as the run of 04:25:16Z except the
  time of the test; [`release-check.txt`](release-check.txt) is the output of
  this run.
- **Note 3, the copy of the release review.** The body of comment 5965588985 on
  #96 and [`release-review.md`](release-review.md) are the same text, 10,782
  characters (a comparison of the body that the API gives, at 16:11:50Z).
- **Note 5, the argument `HEAD`.** The header of the check and the sentence of
  D1 in `docs/plan/README.md` name the full commit ID.
