#!/bin/sh
# The evidence check of task T-8bwd (#387): the handoff-item mapping of
# docs/spec/loop.md, "### The items a caller computes", against anchors it does
# not write: the item lists of internal/handoff and internal/records, the
# default rows of internal/handoff, the transition table and "Not in M2e" of
# loop.md, the task rows of docs/plan/README.md, and the verdict words and
# field names of engineering-discipline.md, review-record-lint.sh,
# architecture.md and loop.md's own text. A one-off evidence command of this
# task, wired into no hook, no CI job and no run-discipline-tests. Run from the
# root of the repository; exit 0 when each rule holds, 1 otherwise. Each FAIL
# names the item and the anchor it disagrees with.

loop=docs/spec/loop.md
plan=docs/plan/README.md
handoff=internal/handoff/table.go
records=internal/records/session.go
disc=docs/engineering-discipline.md
lint=docs/ci/review-record-lint.sh
arch=docs/architecture.md
out=$(mktemp)
fail=0
say() { echo "$@" >> "$out"; }

# The mapping rows of the section: item, computed-by cell, milestone cell.
sec=$(mktemp)
awk '/^### The items a caller computes/{s=1; next} s && /^###? /{s=0} s' "$loop" > "$sec"
if [ ! -s "$sec" ]; then
	# The baseline failure: the section does not exist. One FAIL per expected
	# item and per anchor check, a legitimate assertion failure, not ENOENT.
	for it in schema artifacts plan-fields verdict-field tests-red-green spec-sections \
		patch-per-kind register-rows rendering-equal nothing-material findings-located \
		heading-ids question-row answer-row; do
		say "FAIL (1) the section does not exist: $it has no mapping row (anchor: $handoff, $records)"
		say "FAIL (2) the section does not exist: $it is unmapped for the default rows (anchor: $handoff Default())"
		say "FAIL (3) the section does not exist: $it has no owner row or deferral (anchor: $plan, $loop 'Not in M2e')"
	done
	say "FAIL (4) the section does not exist: the pinned forms and verdict words (anchors: $disc, $lint, $arch)"
	say "FAIL (5) the section does not exist: the no-marker rule cannot be checked"
	cat "$out"; rm -f "$out" "$sec"; exit 1
fi

# The mapping table rows: | `item` | <rule> | <computed by> | <milestone> |
map=$(mktemp)
awk -F'|' '/^\| `/{it=$2; gsub(/[` ]/,"",it); cb=$4; ms=$5; gsub(/^ +| +$/,"",cb); gsub(/^ +| +$/,"",ms); print it "\t" cb "\t" ms}' "$sec" > "$map"

# Check 1: the mapping's first column equals the items list of internal/handoff
# and the handoffItems of internal/records (minus pair).
items_of() { sed -n "s/.*$1 = \[\]string{\(.*\)}.*/\1/p" "$2" | tr ',' '\n' | tr -d ' "'; }
items_of 'items' "$handoff" | sort > /tmp/it-handoff.$$
items_of 'handoffItems' "$records" | grep -v '^pair$' | sort > /tmp/it-records.$$
diff -q /tmp/it-handoff.$$ /tmp/it-records.$$ > /dev/null || say "FAIL (1) the two code anchors differ: $handoff items vs $records handoffItems"
cut -f1 "$map" | sort > /tmp/it-map.$$
join -v1 /tmp/it-map.$$ /tmp/it-handoff.$$ | while read -r x; do say "FAIL (1) item in the mapping, not in $handoff: $x"; done
join -v2 /tmp/it-map.$$ /tmp/it-handoff.$$ | while read -r x; do say "FAIL (1) item of $handoff not in the mapping: $x"; done
join -v1 /tmp/it-map.$$ /tmp/it-records.$$ | while read -r x; do say "FAIL (1) item in the mapping, not in $records: $x"; done
join -v2 /tmp/it-map.$$ /tmp/it-records.$$ | grep -v '^pair$' | while read -r x; do say "FAIL (1) item of $records not in the mapping: $x"; done

# Check 2a: each item of each Default() row is in the mapping.
grep -o 'w("[^)]*"))' "$handoff" | sed 's/^w("//; s/"))$//' | tr ',' '\n' | tr -d ' "' | sort -u | while read -r x; do
	[ -z "$x" ] || grep -q "^$x	" "$map" || say "FAIL (2) item of a Default() row not in the mapping: $x"
done
# Check 2b: each (M2f)/(M2g) mark of the transition table agrees with the
# mapping's milestone column.
awk '/^\| From, to/{t=1; next} t && /^$/{t=0} t && /^\| /' "$loop" |
	grep -o '`[a-z-]*` *(`M2[fg]`)' | sed 's/`//g; s/ *(/ /; s/)//' |
	while read -r x m; do
		grep -q "^$x	.*	$m" "$map" || say "FAIL (2) transition-table mark ($m) of $x disagrees with the mapping's milestone"
	done
# Check 2c: the row posting a question and an answer is named (66k).
grep -q '66k' "$sec" || say "FAIL (2) the section names no row for the question and answer post (anchor: $loop)"

# Check 3: for each item the mapping gives to M2e, the owner row exists in the
# plan table and that row's own Task or Demo cell names the item; for each
# moved item, "Not in M2e" names it. Disclosed exception: tests-red-green/66e —
# the confirmation forbids editing the cells of rows 66e-66k (C1), so the check
# verifies the row exists and prints the exception with its reason.
notin=$(mktemp); awk '/^## Not in M2e/{s=1; next} s && /^## /{s=0} s' "$loop" > "$notin"
tasks=$(mktemp); awk '/^## The tasks of M2e/{s=1; next} s && /^## /{s=0} s && /^\| [0-9]+[a-z]* \|/' "$plan" > "$tasks"
while IFS='	' read -r item cb ms; do
	# The primary owner is the first "row N" token of the computed-by cell;
	# secondary citations (a row a parse also serves) are not owners.
	r=$(echo "$cb" | grep -o 'row [0-9a-z]*' | head -1 | cut -d' ' -f2)
	[ -n "$r" ] && ! grep -q "^| $r |" "$tasks" && say "FAIL (3) owner row $r of $item is not a row of The tasks of M2e"
	case "$ms" in
	M2e*)
		for r in $r; do
			if [ "$item:$r" = "tests-red-green:66e" ]; then
				say "NOTE (3) exception: $item's owner 66e keeps its cells (C1); the row exists and the item is named in the mapping"
				continue
			fi
			grep "^| $r |" "$tasks" | grep -q "$item" || say "FAIL (3) $item's owner row $r does not name it in its Task or Demo cell"
		done
		;;
	esac
	case "$ms" in
	*M2f*|*M2g*)
		grep -q "\`$item\`" "$notin" || say "FAIL (3) moved item $item is not named in 'Not in M2e'"
		;;
	esac
done < "$map"

# Check 4: the verdict words and the field names equal their anchors.
for w in "material" "nothing material in scope" "not mergeable, findings recorded"; do
	grep -q "$w" "$disc" || say "FAIL (4) verdict word not in $disc: $w"
	grep -q "$w" "$sec" || say "FAIL (4) verdict word not in the mapping section: $w"
done
for w in approve approve-with-conditions reject; do
	grep -q "$w" "$lint" || say "FAIL (4) plan-review word not in $lint: $w"
	grep -q "$w" "$sec" || say "FAIL (4) plan-review word not in the mapping section: $w"
done
tier=$(mktemp); awk '/^### The tier and the context of a task/{s=1; next} s && /^###? /{s=0} s' "$loop" > "$tier"
tierflat=$(tr '\n' ' ' < "$tier" | tr -d '`')
for f in size "open questions" "new interfaces" "new dependencies"; do
	grep -q "$f" "$arch" || say "FAIL (4) plan field not in $arch: $f"
	echo "$tierflat" | grep -q "$f" || say "FAIL (4) plan field not in loop.md 'The tier': $f"
	grep -q "$f:" "$sec" || say "FAIL (4) plan field not pinned in the section: $f"
done

# Check 5: no open-gap marker in the new text (the section and row 66l). The
# marker character is built from its UTF-8 bytes so this script, a tracked
# file, holds no unpaired marker itself.
lt=$(printf '\342\200\271')
grep -q "$lt" "$sec" && say "FAIL (5) an open-gap marker is in the mapping section"
grep '^| 66l ' "$plan" | grep -q "$lt" && say "FAIL (5) an open-gap marker is in row 66l of the plan table"

cat "$out"
grep -q '^FAIL' "$out" && fail=1
rm -f "$out" "$sec" "$map" "$notin" "$tasks" "$tier" /tmp/it-handoff.$$ /tmp/it-records.$$ /tmp/it-map.$$
exit $fail
