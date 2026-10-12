#!/bin/sh
# The check of task T-bgvq (#288): the two lessons of the review rounds of T-bng8 are
# the last two entries of the first list of section 2 of docs/guardrails.md, between
# its heading and "Writing a lesson back", each in the form that "Writing a lesson
# back" asks (the trap, why it is silent, the check) and each ending "Learned in
# `T-bng8`". A one-off evidence command of this task, not a check of the gate. Run
# from the root of the repository; exit 0 when each rule holds, 1 otherwise. GUARDRAILS
# names another file, for a copy that a mutant changes.
LC_ALL=C
export LC_ALL
doc=${GUARDRAILS:-docs/guardrails.md}
t1='A review record without its heading.'
t2='An allow rule whose glob reaches a script with a write mode.'
out=$(mktemp)
# The entries of the first list: one line each, `lines<TAB>unindented<TAB>text`, the lines
# of a top-level bullet joined by a space. A line after the first must start with two spaces.
awk '
	function flush() { if (cur != "") print n "\t" bad "\t" cur; cur = ""; n = 0; bad = 0 }
	/^## 2\. Known pitfalls/ { on = 1; next }
	/^### Writing a lesson back/ { flush(); on = 0 }
	!on { next }
	/^- / { flush(); cur = $0; n = 1; next }
	cur != "" && /^  / { cur = cur " " substr($0, 3); n++; next }
	cur != "" && /^$/ { flush(); next }
	cur != "" { bad++; cur = cur " " $0; n++ }
	END { flush() }' "$doc" > "$out"
[ -s "$out" ] || { echo "FAIL $doc has no first list of section 2 before \"Writing a lesson back\""; rm -f "$out"; exit 1; }
fail=0
entry() { grep "^[0-9]*	[0-9]*	- ❌ \*\*$1\*\*" "$out"; }
present() { # present TITLE: 0 when exactly one entry has the title, else one FAIL line
	n=$(entry "$1" | wc -l | tr -d ' ')
	[ "$n" = 1 ] && return 0
	echo "FAIL the first list of section 2 holds $n entries titled: $1"; fail=1; return 1
}
has() { # has TITLE TEXT: the entry holds TEXT, else one FAIL line
	entry "$1" | grep -Fq -- "$2" && return 0
	echo "FAIL the entry \"$1\" does not hold: $2"; fail=1; return 1
}
form() { # form TITLE: the rules of one entry that is present
	bad=$(entry "$1" | cut -f2)
	[ "$bad" = 0 ] || { echo "FAIL the entry \"$1\" has $bad lines that are not indented by two spaces"; fail=1; }
	has "$1" 'It is silent' && has "$1" '**The check:**' && has "$1" 'Learned in `T-bng8`'
	# The parts come in the order of the form: why it is silent, the check, the task.
	entry "$1" | awk -F'\t' '{ s = index($3, "It is silent"); c = index($3, "**The check:**"); l = index($3, "Learned in `T-bng8`"); exit !(s > 0 && c > s && l > c) }' ||
		{ echo "FAIL the entry \"$1\" does not hold its parts in the order: silent, check, Learned in"; fail=1; }
}
present "$t1"; p1=$?
present "$t2"; p2=$?
[ "$p1" = 0 ] && form "$t1"
[ "$p2" = 0 ] && form "$t2"
# Placement: the two lessons are the last two entries of the first list, in this order.
if [ "$p1" = 0 ] && [ "$p2" = 0 ]; then
	last2=$(cut -f3 "$out" | tail -n 2 | cut -c1-200)
	case "$last2" in
	"- ❌ **$t1**"*"
- ❌ **$t2**"*) ;;
	*) echo "FAIL the two lessons are not the last two entries of the first list of section 2"; fail=1 ;;
	esac
fi
# The conditions of the plan review: the local form of the check, and the harness outside the repository.
[ "$p1" = 0 ] && has "$t1" 'docs/ci/review-record-lint.sh' && has "$t1" '--paginate' && has "$t1" 'next push'
[ "$p2" = 0 ] && has "$t2" 'outside the repository' && has "$t2" 'gen-issues.py'
rm -f "$out"
[ "$fail" = 0 ] && echo "ok   the two lessons of T-bng8 are the last entries of the first list of section 2, in the form of their neighbours"
exit "$fail"
