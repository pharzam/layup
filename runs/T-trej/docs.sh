#!/bin/sh
# The check of the documents of task T-trej (#130, row 25a of the plan). A
# one-off evidence command of this task, not a check of the gate: no hook and
# no CI job runs it. Run from the root of the repository; exit 0 when each rule
# holds, 1 otherwise.
fail=0
ok() { printf 'ok   %s\n' "$1"; }
no() { printf 'FAIL %s\n' "$1"; fail=1; }
has() { grep -q -- "$2" "$1" && ok "$1: $3" || no "$1 does not say: $3"; }
row() { awk -v r="$1" -F'|' '/^## The tasks of M2a/{t=1; next} t && /^## /{t=0} t && $2 ~ "^ "r" $"' docs/plan/README.md; }
# 1. run.md: the R10 decisions (conditions 1, 3, 4 and note 1 of the plan review).
has docs/spec/run.md 'the takeover rule comes first' "the order of the two rules of the wait"
has docs/spec/run.md 'once more' "a refused push while the run holds the lease"
has docs/spec/run.md 'A refused takeover' "a refused takeover"
has docs/spec/run.md 'edited before its first copy' "a comment edited before its first copy"
# 2. run.md: the acceptance row of the lease (condition 2).
grep '^| The lease and fencing |' docs/spec/run.md | grep -q 'a stand-in records store' && ok "run.md: the acceptance row of the lease" || no "the acceptance row of the lease does not name a stand-in records store"
# 3. The plan: rows 25a and 25b, row 26 After 25b (O-173 a; condition 5).
case "$(row 25a)" in *'`T-trej`'*'#130'*) ok "plan: row 25a" ;; *) no "plan: no row 25a of T-trej" ;; esac
r25b=$(row 25b)
case "$r25b" in *'`T-ax3r`'*'#142'*'NFR-006'*'REQ-002'*'25a | 1 |') ok "plan: row 25b with its parts and After 25a" ;; *) no "plan: row 25b is missing a part" ;; esac
case "$(row 26)" in *'| 25b | 1 |') ok "plan: row 26 After 25b" ;; *) no "plan: row 26 is not After 25b" ;; esac
case "$(row 25)" in '') ok "plan: no row 25" ;; *) no "plan: row 25 is still there" ;; esac
has docs/plan/README.md 'O-173 of #130' "the sentence before the table names O-173"
# 4. guardrails.md §2: the lesson of the slicing (note 8).
has docs/guardrails.md 'A plan row sliced by package, with no goal count' "the slicing pitfall"
# 5. The traceability and the Test cell of NFR-001.
for t in TestAStoppedHeartbeatIsTakenOverAfterThreeTimesH TestADecisionIsAFirstCopyByAnApproverInTheRole TestCopyBeforeRead; do
	grep "\`$t\`" docs/tests/traceability.md | grep -q 'T-trej' && ok "traceability: $t" || no "traceability has no row of $t with T-trej"
done
grep '^| NFR-001 | F-0001#1 ' docs/prd/PRD-0001-layup.md | grep -q 'T-trej' && ok "PRD §12, NFR-001" || no "the Test cell of NFR-001 does not name T-trej"
# 6. The backlog: lines of 25a and 25b.
grep -q '^- \*\*T-ax3r\*\*' docs/tasks/backlog.md && ok "backlog: T-ax3r" || no "the backlog has no line of T-ax3r"
exit "$fail"
