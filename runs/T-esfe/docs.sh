#!/bin/sh
# The check of the documents of task T-esfe (#127, row 22a of the plan). A
# one-off evidence command of this task, not a check of the gate: no hook and
# no CI job runs it. Run from the root of the repository; exit 0 when each rule
# holds, 1 otherwise.
fail=0
ok() { printf 'ok   %s\n' "$1"; }
no() { printf 'FAIL %s\n' "$1"; fail=1; }
plan=docs/plan/README.md
cell() { awk -v r="$1" -F'|' '/^## The tasks of M2a/{t=1; next} t && /^## /{t=0} t && $2 ~ "^ "r" $" {print $(NF-2)}' "$plan" | sed 's/^ *//; s/ *$//'; }
# 1. packages.md names row 22a as the task that reads the table of M2a.
grep -q 'task `T-esfe`, row 22a of the plan' docs/spec/packages.md && ok "packages.md names row 22a" || no "packages.md does not name row 22a"
# 2. The plan: rows 22a and 22b, and the edges.
grep -q '^| 22a | `T-esfe` |' "$plan" && ok "the plan has row 22a, T-esfe" || no "the plan has no row 22a of T-esfe"
grep -q '^| 22b | `T-1g1q` |' "$plan" && ok "the plan has row 22b, T-1g1q" || no "the plan has no row 22b of T-1g1q"
[ "$(cell 22b)" = "22a" ] && ok "22b is after 22a" || no "22b is not after 22a: $(cell 22b)"
[ "$(cell 24)" = "22b" ] && ok "24 is after 22b" || no "24 is not after 22b: $(cell 24)"
[ "$(cell 25)" = "21, 22b, 23, 24" ] && ok "25 is after 21, 22b, 23, 24" || no "25 is not after 21, 22b, 23, 24: $(cell 25)"
grep -q '^| 22 |' "$plan" && no "the plan still has a row 22" || ok "the plan has no row 22"
# 3. docs/tests/traceability.md.
grep '`TestPackageRules` (`cmd/layup/rules_integration_test.go`)' docs/tests/traceability.md | grep -q 'T-esfe' && ok "traceability: TestPackageRules names T-esfe" || no "the row of TestPackageRules does not name T-esfe"
grep '`TestReadAdapter`' docs/tests/traceability.md | grep -q 'T-esfe' && ok "traceability: the unit tests of the checker" || no "traceability has no row of the unit tests of the checker with T-esfe"
# 4. The Test cell of NFR-007.
grep '^| NFR-007 | F-0004#1 ' docs/prd/PRD-0001-layup.md | grep -q 'TestReadAdapter.*T-esfe' && ok "PRD-0001 §12, NFR-007: the checker of M2a" || no "the Test cell of NFR-007 does not name the checker of M2a (T-esfe)"
exit "$fail"
