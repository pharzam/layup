#!/bin/sh
# The check of D3 and D4 of task T-8kqn (#126): the documents name the records
# of Start and their tests. A one-off evidence command of this task, not a
# check of the gate. Run from the root of the repository; exit 0 when each
# rule holds, 1 otherwise.
fail=0
ok() { printf 'ok   %s\n' "$1"; }
no() { printf 'FAIL %s\n' "$1"; fail=1; }
# 1. The Job cell of internal/records in the table of phase 1 names them.
cell=$(awk '/^### The table of phase 1/{t=1} /^### The table of M2a/{t=0} t && /^\| `internal\/records` \|/' docs/spec/packages.md)
case "$cell" in *"the records of Start"*) ok "the Job cell of internal/records names the records of Start" ;; *) no "the Job cell of internal/records does not name the records of Start" ;; esac
# 2. docs/tests/traceability.md has a row for each new test, with the task.
for t in TestAValidStartIsRead TestStartRefusesEachBrokenRule TestApproversRefusesEachBrokenRule TestLeaseRefusesEachBrokenRule TestCopiesRefusesEachBrokenRule; do
	if grep "\`$t\`" docs/tests/traceability.md | grep -q 'T-8kqn'; then ok "traceability: $t"; else no "traceability has no row of $t with T-8kqn"; fi
done
if grep '`TestTheSchemasEqualTheirBlocks` (`internal/records/records_integration_test.go`)' docs/tests/traceability.md | grep -q 'NFR-001.*T-8kqn'; then ok "traceability: TestTheSchemasEqualTheirBlocks of internal/records names NFR-001 and T-8kqn"; else no "the row of TestTheSchemasEqualTheirBlocks of internal/records does not name NFR-001 and T-8kqn"; fi
# 3. The Test cell of NFR-001 in PRD-0001 §12 names the tests of the records of Start.
row=$(grep '^| NFR-001 | F-0001#1 ' docs/prd/PRD-0001-layup.md)
case "$row" in *"TestStartRefusesEachBrokenRule"*"T-8kqn"*) ok "PRD-0001 §12, NFR-001: the tests of the records of Start" ;; *) no "the Test cell of NFR-001 does not name the tests of the records of Start (T-8kqn)" ;; esac
exit "$fail"
