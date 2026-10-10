#!/bin/sh
# The check of D1 of task T-2r8u (#211): each part of the specification of M2c is
# in exactly one row of "The tasks of M2c" in docs/plan/README.md, or named with
# a reason (parts.tsv: file, heading, sub, row). A one-off evidence command of
# this task, not a check of the gate. Run from the root of the repository; exit
# 0 when each rule holds, 1 otherwise.
tab=runs/T-2r8u/parts.tsv
plan=docs/plan/README.md
spec=docs/spec/intake.md
acc='## The acceptance tests of M2c'
out=$(mktemp)
fail=0
# 1. Each heading of intake.md outside a fenced block is in the list.
awk '/^```/{f=!f; next} !f && (/^## / || /^### /)' "$spec" | while IFS= read -r h; do
	grep -Fq "$spec	$h	" "$tab" || echo "FAIL a heading of $spec has no part: $h"
done >> "$out"
# 2. Each part (file, heading, sub) is listed once.
tail -n +2 "$tab" | cut -f1-3 | sort | uniq -d | sed 's/^/FAIL a part is listed twice: /' >> "$out"
# 3. Each row of the list is a row of the table "The tasks of M2c", and each row
#    of that table has a part.
rows=$(awk '/^## The tasks of M2c/{t=1; next} t && /^## /{t=0} t && /^\| [0-9]+[a-z]? \|/{print $2}' "$plan")
listed=$(tail -n +2 "$tab" | cut -f4 | grep -v '^none' | sort -u)
printf '%s\n' "$listed" | while read -r r; do
	printf '%s\n' "$rows" | grep -qx "$r" || echo "FAIL row $r is not a row of The tasks of M2c"
done >> "$out"
printf '%s\n' "$rows" | grep . | while read -r r; do
	printf '%s\n' "$listed" | grep -qx "$r" || echo "FAIL row $r of The tasks of M2c has no part"
done >> "$out"
# 4. Each heading of runs/T-g49f/sections.tsv maps to a part.
tail -n +2 runs/T-g49f/sections.tsv | cut -f2,3 | while IFS='	' read -r f h; do
	grep -Fq "$f	$h	" "$tab" || echo "FAIL a heading of sections.tsv has no part: $f $h"
done >> "$out"
# 5. Each row of "The acceptance tests of M2c" is a part, whole or by its
#    clauses ("<Part>: <clause>"), so the test of each part has its row.
awk -v a="$acc" '$0==a{t=1; next} t && /^## /{t=0} t && /^\| / && !/^\| (Part|-)/' "$spec" |
	cut -d'|' -f2 | sed 's/^ *//; s/ *$//' | while IFS= read -r p; do
	awk -F'\t' -v f="$spec" -v h="$acc" -v p="$p" \
		'$1==f && $2==h && ($3==p || index($3, p ": ")==1){n++} END{exit !n}' "$tab" ||
		echo "FAIL an acceptance row of $spec has no part: $p"
done >> "$out"
if [ -s "$out" ]; then cat "$out"; fail=1; else echo "ok   each part is in one row of The tasks of M2c, or named with its reason"; fi
rm -f "$out"
exit "$fail"
