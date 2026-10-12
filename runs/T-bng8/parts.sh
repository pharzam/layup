#!/bin/sh
# The check of step T of task T-bng8 (#267): each part of the specification of
# M2d is in exactly one row of "The tasks of M2d" in docs/plan/README.md, or
# named with a reason (parts.tsv: file, heading, sub, row). Copied from
# runs/T-2r8u/parts.sh and changed for M2d, with condition 2 of the plan review
# of #267: the input states, the records of the Scaffold, the calls of M2d, the
# table of M2d, CommitTree and each built block changed in words. A one-off
# evidence command of this task, not a check of the gate. Run from the root of
# the repository; exit 0 when each rule holds, 1 otherwise.
tab=runs/T-bng8/parts.tsv
plan=docs/plan/README.md
spec=docs/spec/scaffold.md
acc='## The acceptance tests of M2d'
inp='## Input states'
out=$(mktemp)
fail=0
# 1. Each heading of scaffold.md outside a fenced block is in the list.
awk '/^```/{f=!f; next} !f && (/^## / || /^### /)' "$spec" | while IFS= read -r h; do
	grep -Fq "$spec	$h	" "$tab" || echo "FAIL a heading of $spec has no part: $h"
done >> "$out"
# 2. Each part (file, heading, sub) is listed once.
tail -n +2 "$tab" | cut -f1-3 | sort | uniq -d | sed 's/^/FAIL a part is listed twice: /' >> "$out"
# 3. Each row of the list is a row of the table "The tasks of M2d", and each row
#    of that table has a part.
rows=$(awk '/^## The tasks of M2d/{t=1; next} t && /^## /{t=0} t && /^\| [0-9]+[a-z]? \|/{print $2}' "$plan")
[ -n "$rows" ] || echo "FAIL $plan has no table The tasks of M2d" >> "$out"
listed=$(tail -n +2 "$tab" | cut -f4 | grep -v '^none' | sort -u)
printf '%s\n' "$listed" | while read -r r; do
	printf '%s\n' "$rows" | grep -qx "$r" || echo "FAIL row $r is not a row of The tasks of M2d"
done >> "$out"
printf '%s\n' "$rows" | grep . | while read -r r; do
	printf '%s\n' "$listed" | grep -qx "$r" || echo "FAIL row $r of The tasks of M2d has no part"
done >> "$out"
# 4. Each row of the acceptance tests and of the input states of scaffold.md is
#    a part, whole or by its clauses ("<Part>: <clause>").
for sect in "$acc" "$inp"; do
	awk -v a="$sect" '$0==a{t=1; next} t && /^## /{t=0} t && /^\| / && !/^\| (Test|State|-)/' "$spec" |
		cut -d'|' -f2 | sed 's/^ *//; s/ *$//' | while IFS= read -r p; do
		awk -F'\t' -v f="$spec" -v h="$sect" -v p="$p" \
			'$1==f && $2==h && ($3==p || index($3, p ": ")==1){n++} END{exit !n}' "$tab" ||
			echo "FAIL a row of $sect has no part: $p"
	done >> "$out"
done
# 5. The parts outside scaffold.md (condition 2): each kind has its part.
for want in \
	"docs/spec/records.md	## NFR-001 — The records of the Scaffold	the block markers" \
	"docs/spec/records.md	## NFR-001 — The records of the Scaffold	the blocks setup-tree and forge-check" \
	"docs/spec/forge.md	## The calls of M2d	—" \
	"docs/spec/packages.md	### The table of M2d	—" \
	"docs/spec/packages.md	### The calls of \`internal/git\`	CommitTree"; do
	grep -Fq "$want	" "$tab" || echo "FAIL no part: $want" >> "$out"
done
for b in questions run-steps setup-record; do
	grep -q "^docs/spec/records.md	## NFR-001 — The records of the Scaffold	$b: " "$tab" ||
		echo "FAIL the built block $b, changed in words, has no part" >> "$out"
done
if [ -s "$out" ]; then cat "$out"; fail=1; else echo "ok   each part is in one row of The tasks of M2d, or named with its reason"; fi
rm -f "$out"
exit "$fail"
