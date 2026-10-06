#!/bin/sh
# tree-equal.sh: condition 2 of O-146 (task T-evad, #97), with the base that the plan review of the target task asked
# for (O-148 to O-150). The tree of HEAD (the head of the target task, or main after its merge) equals the tree of the
# head of layup-setup-2, the second setup run, except two sets of paths:
#   - the four paths of task T-a0rt, the target's own change of D7: they keep the text of BASE (main before the
#     target task). The three files keep their object and mode. The completed log keeps the text of BASE, with at
#     most one more line, the line of the target task (with --log-line, exactly one);
#   - the target task's own new files, when they are named: its task file (O-147 = a) and its evidence file
#     (O-150 = a). A named file is in HEAD, and in neither layup-setup-2 nor BASE.
# The line of the target task, when the log has it, has the form of the log and is its first entry (an example in
# an HTML comment is not an entry): - **YYYY-MM-DD** — **T-xxxx** — summary
# ([#N](https://github.com/OWNER/NAME/issues/N); [detail](T-xxxx.md)), with the repository and the number of --issue.
# One mechanism and one test: HEAD must equal the tree that these rules make from layup-setup-2 and BASE.
#
# Usage: sh tree-equal.sh REPO HEAD SETUP2 --base BASE [--task-file docs/tasks/T-xxxx.md
#        [--evidence-file runs/T-xxxx/evidence.md]] [--log-line --issue OWNER/NAME#N]
#   REPO             a clone of the target
#   HEAD             the head of the target task, or main after its merge
#   SETUP2           the head of layup-setup-2
#   --base           main before the target task; it holds the change of T-a0rt
#   --task-file      the file of the target task (gate step 8 of the target)
#   --evidence-file  the evidence file of the same task (gate step 6 of the target)
#   --log-line       the completed log must hold the line of the task
#   --issue          the task's issue, OWNER/NAME#N, which its line links (required with --log-line)
# Exit 0 and `tree-equal: PASS` when HEAD keeps each rule; else exit 1, with one `FAIL:` line for each finding.
# Exit 2 on an input error.
set -u
usage='usage: sh tree-equal.sh REPO HEAD SETUP2 --base BASE [--task-file docs/tasks/T-xxxx.md [--evidence-file runs/T-xxxx/evidence.md]] [--log-line --issue OWNER/NAME#N]'
R=${1:?$usage}
H=${2:?$usage}
S=${3:?$usage}
shift 3
LOG=docs/tasks/completed.md
FILES3="docs/onboarding-for-engineers.md docs/tasks/T-a0rt.md runs/T-a0rt/evidence.md"
A0RT="$FILES3 $LOG"
BASE= TASK= EVID= NEEDLINE=0 ISSUE=
inerr() { echo "tree-equal: input error: $*" >&2; exit 2; }
while [ $# -gt 0 ]; do
	case $1 in
	--base) [ $# -ge 2 ] || inerr "--base needs a commit"; BASE=$2; shift 2 ;;
	--task-file) [ $# -ge 2 ] || inerr "--task-file needs a path"; TASK=$2; shift 2 ;;
	--evidence-file) [ $# -ge 2 ] || inerr "--evidence-file needs a path"; EVID=$2; shift 2 ;;
	--log-line) NEEDLINE=1; shift ;;
	--issue) [ $# -ge 2 ] || inerr "--issue needs OWNER/NAME#N"; ISSUE=$2; shift 2 ;;
	*) inerr "an unknown argument: $1" ;;
	esac
done
[ -n "$BASE" ] || inerr "--base is required: the paths of T-a0rt are compared with the base"
ID=
if [ -n "$TASK" ]; then
	case " $A0RT " in *" $TASK "*) inerr "the task file is one of the four paths of T-a0rt: $TASK" ;; esac
	case $TASK in
	docs/tasks/T-[0-9a-z][0-9a-z][0-9a-z][0-9a-z].md) ID=${TASK#docs/tasks/}; ID=${ID%.md} ;;
	*) inerr "the task file must be docs/tasks/T-xxxx.md: $TASK" ;;
	esac
fi
if [ -n "$EVID" ]; then
	[ -n "$ID" ] || inerr "--evidence-file needs --task-file"
	[ "$EVID" = "runs/$ID/evidence.md" ] || inerr "the evidence file must be runs/$ID/evidence.md: $EVID"
fi
[ "$NEEDLINE" = 0 ] || [ -n "$ID" ] || inerr "--log-line needs --task-file"
[ "$NEEDLINE" = 0 ] || [ -n "$ISSUE" ] || inerr "--log-line needs --issue"
IREPO= INUM=
if [ -n "$ISSUE" ]; then
	[ -n "$ID" ] || inerr "--issue needs --task-file"
	IREPO=${ISSUE%#*} INUM=${ISSUE##*#}
	case $INUM in ''|*[!0-9]*) inerr "--issue needs OWNER/NAME#N: $ISSUE" ;; esac
	case $IREPO in [A-Za-z0-9]*/[A-Za-z0-9._-]*) ;; *) inerr "--issue needs OWNER/NAME#N: $ISSUE" ;; esac
	case $IREPO in */*/*|*[!A-Za-z0-9._/-]*) inerr "--issue needs OWNER/NAME#N: $ISSUE" ;; esac
fi

gt() { git -C "$R" -c core.quotepath=off "$@"; }
h=$(gt rev-parse --verify -q "$H^{commit}") && s=$(gt rev-parse --verify -q "$S^{commit}") \
	&& b=$(gt rev-parse --verify -q "$BASE^{commit}") || inerr "$H, $S or $BASE is not a commit of $R"
for p in $A0RT; do
	gt cat-file -e "$b:$p" 2>/dev/null || inerr "the base has no $p: the base must be main after T-a0rt"
done
TMP=$(mktemp -d) || inerr "no temporary directory"
trap 'rm -rf "$TMP"' EXIT INT TERM
echo "tree-equal: head $h against layup-setup-2 $s, base $b"
[ -z "$TASK" ] || echo "tree-equal: the files of the task: $TASK${EVID:+ and $EVID}"
fail=0
F() { echo "FAIL: $*"; fail=$((fail + 1)); }

# 1. Each path that differs from layup-setup-2 is a path of T-a0rt (step 2) or a named file of the task (step 3).
out=$(gt diff --no-renames --name-status "$s" "$h") || inerr "git diff failed"
while IFS='	' read -r st p; do
	[ -n "$p" ] || continue
	case " $A0RT " in *" $p "*) continue ;; esac
	if [ -n "$TASK" ] && { [ "$p" = "$TASK" ] || [ "$p" = "$EVID" ]; }; then continue; fi
	F "$p ($st)"
done <<EOF
$out
EOF

# 2. The change of T-a0rt keeps the text of the base.
for p in $FILES3; do
	eb=$(gt ls-tree "$b" -- "$p") eh=$(gt ls-tree "$h" -- "$p")
	if [ -z "$eh" ]; then F "T-a0rt: $p is not in the head"
	elif [ "$eh" != "$eb" ]; then F "T-a0rt: $p differs from the base"
	else echo "T-a0rt: $p as in the base"; fi
done
lb=$(gt ls-tree "$b" -- "$LOG") lh=$(gt ls-tree "$h" -- "$LOG")
if [ -z "$lh" ]; then F "T-a0rt: $LOG is not in the head"
else
	[ "${lh%% *}" = "${lb%% *}" ] || F "T-a0rt: $LOG: its mode differs from the base"
	gt cat-file blob "$b:$LOG" > "$TMP/base" && gt cat-file blob "$h:$LOG" > "$TMP/head" || inerr "cannot read $LOG"
	mark=; [ -z "$ID" ] || mark="**$ID**"
	nb=0 nh=0
	if [ -n "$mark" ]; then nb=$(grep -cF -- "$mark" "$TMP/base"); nh=$(grep -cF -- "$mark" "$TMP/head"); fi
	if [ "$nb" -gt 0 ]; then F "T-a0rt: $LOG: the base already has a line of $ID"
	elif cmp -s "$TMP/base" "$TMP/head"; then
		if [ "$NEEDLINE" = 1 ]; then F "T-a0rt: $LOG has no line of $ID"; else echo "T-a0rt: $LOG as in the base"; fi
	elif [ -z "$mark" ]; then F "T-a0rt: $LOG differs from the base"
	elif [ "$nh" -eq 0 ]; then F "T-a0rt: $LOG differs from the base, and it has no line of $ID"
	elif [ "$nh" -gt 1 ]; then F "T-a0rt: $LOG has $nh lines of $ID"
	else
		k=$(grep -nF -- "$mark" "$TMP/head" | cut -d: -f1)
		sed "${k}d" "$TMP/head" > "$TMP/rest"
		if cmp -s "$TMP/base" "$TMP/rest"; then
			# The form of the log, and the first entry outside an HTML comment (most recent first).
			nre=${INUM:-[0-9]+} rre='[^/ ]+/[^/ ]+'
			[ -z "$IREPO" ] || rre=$(printf '%s' "$IREPO" | sed 's/\./\\./g')
			pat='^- \*\*[0-9]{4}-[0-9]{2}-[0-9]{2}\*\* — \*\*'"$ID"'\*\* — .+ \(\[#'"$nre"'\]\(https://github\.com/'"$rre"'/issues/'"$nre"'\); \[detail\]\('"$ID"'\.md\)\)$'
			first=$(LC_ALL=C awk '/<!--/ { c = 1 } !c && /^- \*\*/ { print NR; exit } /-->/ { c = 0 }' "$TMP/head")
			if ! sed -n "${k}p" "$TMP/head" | LC_ALL=C grep -Eq -- "$pat"; then
				F "T-a0rt: $LOG: the line of $ID does not have the form of the log: - **YYYY-MM-DD** — **$ID** — a summary ([#${INUM:-N}](https://github.com/${IREPO:-OWNER/NAME}/issues/${INUM:-N}); [detail]($ID.md))"
			elif [ "$first" != "$k" ]; then F "T-a0rt: $LOG: the line of $ID is not the first entry of the log"
			else echo "T-a0rt: $LOG as in the base, with the line of $ID"; fi
		else F "T-a0rt: $LOG differs from the base in more than the line of $ID"; fi
	fi
fi

# 3. The named files of the task are new: in the head, and in neither layup-setup-2 nor the base.
for f in $TASK $EVID; do
	if [ "$f" = "$TASK" ]; then kind=task; else kind=evidence; fi
	if ! gt cat-file -e "$h:$f" 2>/dev/null; then F "$f: the $kind file is not in the head"
	elif gt cat-file -e "$s:$f" 2>/dev/null; then F "$f: the $kind file is in layup-setup-2, so it is not a new file"
	elif gt cat-file -e "$b:$f" 2>/dev/null; then F "$f: the $kind file is in the base, so it is not a new file"
	else echo "$kind file: $f (new)"; fi
done

if [ "$fail" -eq 0 ]; then
	echo "tree-equal: PASS (the head equals layup-setup-2, except the four paths of T-a0rt, which keep the text of the base${TASK:+, and the new files of $ID})"
	exit 0
fi
echo "tree-equal: FAIL (findings: $fail)"
exit 1
