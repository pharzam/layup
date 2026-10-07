#!/bin/sh
# The check of D1 of task T-zwke (#124): each part of the specification of M2a is
# in exactly one row of "The tasks of M2a" in docs/plan/README.md, or named with
# a reason (parts.tsv: file, heading, sub, row). A one-off evidence command of
# this task, not a check of the gate. Run from the root of the repository; exit
# 0 when each rule holds, 1 otherwise.
tab=runs/T-zwke/parts.tsv
plan=docs/plan/README.md
out=$(mktemp)
fail=0
# 1. Each heading of the M2a files outside a fenced block is in the list.
for f in docs/spec/run.md docs/spec/forge.md; do
	awk '/^```/{f=!f; next} !f && /^#{2,3} /' "$f" | while IFS= read -r h; do
		grep -Fq "$f	$h	" "$tab" || echo "FAIL a heading of $f has no part: $h"
	done
done >> "$out"
# 2. Each part (file, heading, sub) is listed once.
tail -n +2 "$tab" | cut -f1-3 | sort | uniq -d | sed 's/^/FAIL a part is listed twice: /' >> "$out"
# 3. Each row of the list is a row of the table "The tasks of M2a".
rows=$(awk '/^## The tasks of M2a/{t=1; next} t && /^## /{t=0} t && /^\| [0-9]+ \|/{print $2}' "$plan")
tail -n +2 "$tab" | cut -f4 | grep -v '^none' | sort -u | while read -r r; do
	printf '%s\n' "$rows" | grep -qx "$r" || echo "FAIL row $r is not a row of The tasks of M2a"
done >> "$out"
# 4. Each heading of runs/T-zck8/sections.tsv maps to a part.
tail -n +2 runs/T-zck8/sections.tsv | cut -f2,3 | while IFS='	' read -r f h; do
	grep -Fq "$f	$h	" "$tab" || echo "FAIL a heading of sections.tsv has no part: $f $h"
done >> "$out"
if [ -s "$out" ]; then cat "$out"; fail=1; else echo "ok   each part is in one row of The tasks of M2a, or named with its reason"; fi
rm -f "$out"
exit "$fail"
