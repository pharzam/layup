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
plan=${PLAN:-docs/plan/README.md}
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
#    a part, whole or by its clauses ("<Part>: <clause>"), and each clause of a
#    row split between task rows — of a test, an input state or a heading — has
#    its own entry: one surviving clause does not cover another (findings 1 of
#    rounds 5 and 6 of #267).
for sect in "$acc" "$inp"; do
	awk -v a="$sect" '$0==a{t=1; next} t && /^## /{t=0} t && /^\| / && !/^\| (Test|State|-)/' "$spec" |
		cut -d'|' -f2 | sed 's/^ *//; s/ *$//' | while IFS= read -r p; do
		awk -F'\t' -v f="$spec" -v h="$sect" -v p="$p" \
			'$1==f && $2==h && ($3==p || index($3, p ": ")==1){n++} END{exit !n}' "$tab" ||
			echo "FAIL a row of $sect has no part: $p"
	done >> "$out"
done
for want in \
	"docs/spec/scaffold.md	## Input states	No admitted pair of \`systems-architect\` (for \`markers\`) or \`product-owner\` (for a prose session): for markers" \
	"docs/spec/scaffold.md	## Input states	No admitted pair of \`systems-architect\` (for \`markers\`) or \`product-owner\` (for a prose session): for a prose session" \
	"docs/spec/scaffold.md	## The acceptance tests of M2d	The answers of the setup: the rows of answers.tsv" \
	"docs/spec/scaffold.md	### The answers of the setup	the merge into the batch and the comment" \
	"docs/spec/scaffold.md	### The answers of the setup	the rows of answers.tsv" \
	"docs/spec/scaffold.md	### The commands for the Operator	the text of commands.sh" \
	"docs/spec/scaffold.md	### The commands for the Operator	the print, the wait, the apply outputs and setup/rulesets.txt, the work area built again"; do
	grep -Fq "$want	" "$tab" || echo "FAIL no part: $want" >> "$out"
done
# 5. The parts outside scaffold.md (condition 2): each kind has its part, and
#    each clause of a built block or of the table of M2d changed in words has
#    its own entry — one surviving clause of a block does not cover another
#    (finding 2 of round 2 of #267).
for want in \
	"docs/spec/records.md	## NFR-001 — The records of the Scaffold	the block markers" \
	"docs/spec/records.md	## NFR-001 — The records of the Scaffold	the blocks setup-tree and forge-check" \
	"docs/spec/records.md	## NFR-001 — The records of the Scaffold	questions: source setup and the setup refs" \
	"docs/spec/records.md	## NFR-001 — The records of the Scaffold	run-steps: markers" \
	"docs/spec/records.md	## NFR-001 — The records of the Scaffold	run-steps: setup" \
	"docs/spec/records.md	## NFR-001 — The records of the Scaffold	run-steps: verify" \
	"docs/spec/records.md	## NFR-001 — The records of the Scaffold	run-steps: scaffold-push" \
	"docs/spec/records.md	## NFR-001 — The records of the Scaffold	run-steps: forge-check" \
	"docs/spec/records.md	## NFR-001 — The records of the Scaffold	setup-record: the computed refs of S01 and S02" \
	"docs/spec/records.md	## NFR-001 — The records of the Scaffold	setup-record: the row without-layup" \
	"docs/spec/records.md	## NFR-001 — The records of the Scaffold	setup/inputs/files" \
	"docs/spec/records.md	## NFR-001 — The records of the Scaffold	setup/rulesets.txt" \
	"docs/spec/forge.md	## The calls of M2d	—" \
	"docs/spec/packages.md	### The table of M2d	internal/run: the imports of the step runner" \
	"docs/spec/packages.md	### The table of M2d	internal/run: the step markers" \
	"docs/spec/packages.md	### The table of M2d	internal/run: the step setup" \
	"docs/spec/packages.md	### The table of M2d	internal/run: the step verify" \
	"docs/spec/packages.md	### The table of M2d	internal/run: the step scaffold-push" \
	"docs/spec/packages.md	### The table of M2d	internal/run: the step forge-check" \
	"docs/spec/packages.md	### The table of M2d	internal/setup: the inputs of phase 2" \
	"docs/spec/packages.md	### The table of M2d	internal/setup: the two rulesets of S13" \
	"docs/spec/packages.md	### The table of M2d	internal/forge and internal/forge/github: EffectiveRules and Activity" \
	"docs/spec/packages.md	### The table of M2d	internal/records: the block markers" \
	"docs/spec/packages.md	### The table of M2d	internal/run: the blocks setup-tree and forge-check" \
	"docs/spec/packages.md	### The table of M2d	internal/git: CommitTree" \
	"docs/spec/packages.md	### The table of M2d	internal/run: starts no program for the setup" \
	"docs/spec/packages.md	### The calls of \`internal/git\`	CommitTree"; do
	grep -Fq "$want	" "$tab" || echo "FAIL no part: $want" >> "$out"
done
# 6. Each listed part is named in the Task cell of its own row, in the form that
#    gen-issues.py writes it: `<file>`: <heading> (<sub>) (finding 1 of round 1
#    of #267: a part dropped from its row fails here, not only a missing row),
#    and in no other Task cell — "exactly one" (finding 2 of round 3 of #267:
#    a copy in a second row's cell must fail).
lines=$(awk '/^## The tasks of M2d/{t=1; next} t && /^## /{t=0} t && /^\| [0-9]+[a-z]? \|/' "$plan")
cells=$(printf '%s\n' "$lines" | awk -F' [|] ' '{print $4}')
tail -n +2 "$tab" | grep -v '	none' | while IFS='	' read -r file heading sub row; do
	name=${file##*/}
	h=$(printf '%s' "$heading" | sed 's/^#* *//')
	if [ "$sub" = "—" ]; then text="\`$name\`: $h"; else text="\`$name\`: $h ($sub)"; fi
	line=$(printf '%s\n' "$lines" | awk -v r="$row" '$2==r')
	task=$(printf '%s' "$line" | awk -F' [|] ' '{print $4}')
	case "$task" in
	*"$text"*) ;;
	*) echo "FAIL row $row of The tasks of M2d does not name its part: $text" ;;
	esac
	n=$(printf '%s\n' "$cells" | grep -Fc "$text" || true)
	[ "$n" = "1" ] || echo "FAIL the part is named in $n Task cells of The tasks of M2d, not one: $text"
done >> "$out"
if [ -s "$out" ]; then cat "$out"; fail=1; else echo "ok   each part is in one row of The tasks of M2d, or named with its reason"; fi
rm -f "$out"
exit "$fail"
