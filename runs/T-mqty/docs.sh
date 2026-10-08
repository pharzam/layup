#!/bin/sh
# The check of the documents of task T-mqty (#131, row 26 of the plan). A
# one-off evidence command of this task, not a check of the gate: no hook and
# no CI job runs it. Run from the root of the repository; exit 0 when each rule
# holds, 1 otherwise.
fail=0
ok() { printf 'ok   %s\n' "$1"; }
no() { printf 'FAIL %s\n' "$1"; fail=1; }
has() { grep -q -- "$2" "$1" && ok "$1: $3" || no "$1 does not say: $3"; }
# 1. README.md: the two rules of the table, the repeat, the exit code 2 of git.
has docs/spec/README.md 'a row may name one \*selecting flag\*' "the selecting flag"
has docs/spec/README.md '\*flag that may be left out\*' "the flag that may be left out"
has docs/spec/README.md 'repeats in a new world with the same arguments' "condition 2, the repeat of Start"
has docs/spec/README.md 'For `layup run` also `git` not found' "decision 1, git"
# 2. run.md: the form of a line of a wait (condition 4), the two rows of the input states.
has docs/spec/run.md 'layup run: \[<i>/<n>\] <step>: <line>' "condition 4, a line of a wait"
has docs/spec/run.md 'next ten-second beat' "finding 1 of round 2, the beat after a line of a wait"
has docs/spec/run.md '^| `git` not found, or older than 2.32 |' "the input state of git"
has docs/spec/run.md '^| A host directory with no `registers/harnesses.tsv` |' "the input state of a missing harness register"
# 3. forge.md: the adapter's own client (decision 4, condition 3).
has docs/spec/forge.md 'reads no proxy of' "the own client, no proxy"
has docs/spec/forge.md 'SSL_CERT_FILE' "the known limit of the certificate roots"
# 4. packages.md: the May import cell of internal/cli; the Job of internal/run.
grep '^| `internal/cli` | parses' docs/spec/packages.md | grep -q '`internal/run`, `internal/forge`, `internal/forge/github`, `internal/route`' && ok "packages.md: the May import cell of internal/cli" || no "packages.md: the May import cell of internal/cli"
grep '^| `internal/run` |' docs/spec/packages.md | grep -q 'the checks of `git`' && ok "packages.md: the Job of internal/run" || no "packages.md: the Job of internal/run"
# 5. The glossary (note 9).
has docs/glossary.md '^| Selecting flag |' "the glossary: selecting flag"
has docs/glossary.md '^| Flag that may be left out |' "the glossary: flag that may be left out"
# 6. The traceability and the Test cells.
for t in TestRunNewThenRestart TestTheSelectingFlagAndTheOptionalFlag TestANoteOfAWaitTakesThePlaceOfTheNextBeat TestRunNewChecksEachInputBeforeTheFirstStep TestTheOwnClientFollowsNoRedirectAndNoProxy TestCheckGit; do
	grep "\`$t\`" docs/tests/traceability.md | grep -q 'T-mqty' && ok "traceability: $t" || no "traceability has no row of $t with T-mqty"
done
for r in 'NFR-001 | F-0001#1 ' 'NFR-002 | F-0001#2'; do
	grep "^| $r" docs/prd/PRD-0001-layup.md | grep -q 'T-mqty' && ok "PRD §12: $r" || no "the Test cell of $r does not name T-mqty"
done
exit "$fail"
