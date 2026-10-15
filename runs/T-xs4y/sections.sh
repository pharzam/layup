#!/bin/sh
# The check of step T of task T-xs4y (#333): each part of the inventory item
# later-p2-phase-loop-handoffs and later-gate-status, each item that M2a to M2d deferred to M2e, and each settlement of acceptance criterion 2, has its
# heading in docs/spec/ (sections.tsv: part, file, heading). A one-off
# evidence command of this task, not a check of the gate, in the form of
# runs/T-zck8/sections.sh. Run from the root of the repository; exit 0 when
# each heading exists as a whole line of its file, 1 otherwise.
tab=runs/T-xs4y/sections.tsv
out=$(mktemp) || exit 1
trap 'rm -f "$out"' EXIT
tail -n +2 "$tab" | while IFS='	' read -r part file heading; do
	if [ -f "$file" ] && grep -Fqx -- "$heading" "$file"; then
		printf 'ok   %s: %s\n' "$file" "$heading"
	else
		printf 'FAIL %s: %s (%s)\n' "$file" "$heading" "$part"
	fi
done > "$out"
cat "$out"
if grep -q '^FAIL' "$out"; then exit 1; fi
exit 0
