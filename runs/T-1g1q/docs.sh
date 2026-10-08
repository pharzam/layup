#!/bin/sh
# The check of the documents of task T-1g1q (#137, row 22b of the plan). A
# one-off evidence command of this task, not a check of the gate: no hook and
# no CI job runs it. Run from the root of the repository; exit 0 when each rule
# holds, 1 otherwise.
fail=0
ok() { printf 'ok   %s\n' "$1"; }
no() { printf 'FAIL %s\n' "$1"; fail=1; }
row() { awk -v r="$1" -F'|' '/^## The tasks of M2a/{t=1; next} t && /^## /{t=0} t && $2 ~ "^ "r" $"' docs/plan/README.md; }
# 1. forge.md: one PEM block (note 2 of the plan review).
grep -q 'one PEM block of an RSA private key' docs/spec/forge.md && ok "forge.md: one PEM block" || no "forge.md does not say one PEM block"
# 2. packages.md: the Job cell of internal/forge names the key check (note 6).
grep '^| `internal/forge` |' docs/spec/packages.md | grep -q 'the check of the key file' && ok "packages.md: the Job cell of internal/forge" || no "the Job cell of internal/forge does not name the check of the key file"
# 3. The plan: class 3 moved from row 22b to row 24 (O-171).
case "$(row 22b)" in *"The six capabilities"*) no "row 22b still holds The six capabilities" ;; *) ok "row 22b: no class 3" ;; esac
r24=$(row 24)
case "$r24" in *"O-171"*) case "$r24" in *"The six capabilities"*) case "$r24" in *'`internal/forge`'*) ok "row 24: class 3 by O-171, in internal/forge" ;; *) no "row 24 does not name internal/forge" ;; esac ;; *) no "row 24 does not hold The six capabilities" ;; esac ;; *) no "row 24 does not name O-171" ;; esac
# 4. docs/tests/traceability.md.
for t in TestForgeRegisterRefusesEachBrokenRule TestCheckKeyRefusesEachBrokenRule TestHarnessRegisterRefusesEachBrokenRule; do
	grep "\`$t\`" docs/tests/traceability.md | grep -q 'T-1g1q' && ok "traceability: $t" || no "traceability has no row of $t with T-1g1q"
done
# 5. The Test cells of NFR-001 and NFR-007.
grep '^| NFR-001 | F-0001#1 ' docs/prd/PRD-0001-layup.md | grep -q 'T-1g1q' && ok "PRD §12, NFR-001" || no "the Test cell of NFR-001 does not name T-1g1q"
grep '^| NFR-007 | F-0004#1 ' docs/prd/PRD-0001-layup.md | grep -q 'T-1g1q' && ok "PRD §12, NFR-007" || no "the Test cell of NFR-007 does not name T-1g1q"
exit "$fail"
