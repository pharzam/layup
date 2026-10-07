#!/bin/sh
# The check of D1 and D2 of task T-zck8 (#123): each part of the inventory item
# later-p2-run-start has its heading in docs/spec/ (sections.tsv: part, file,
# heading). A one-off evidence command of this task, not a check of the gate.
# Run from the root of the repository; exit 0 when each heading exists as a
# whole line of its file, 1 otherwise.
tab=runs/T-zck8/sections.tsv
fail=0
tail -n +2 "$tab" | while IFS='	' read -r part file heading; do
	if [ -f "$file" ] && grep -Fqx -- "$heading" "$file"; then
		printf 'ok   %s: %s\n' "$file" "$heading"
	else
		printf 'FAIL %s: %s (%s)\n' "$file" "$heading" "$part"
	fi
done > /tmp/sections.$$
cat /tmp/sections.$$
grep -q '^FAIL' /tmp/sections.$$ && fail=1
rm -f /tmp/sections.$$
exit "$fail"
