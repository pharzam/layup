#!/bin/sh
# Tests of tree-equal.sh (task T-evad, #97): condition 2 of O-146, with the base that the plan review of the target
# task T-vu2j asked for (O-148 to O-150). The fixture repositories have the shape of the target:
#   first-setup    the first setup head: a rule, the onboarding and the completed log
#   base           main before the target task: the first setup head with the change of T-a0rt (two files changed,
#                  two files added)
#   layup-setup-2  the second setup head: the first setup head with a changed rule
#   main           the head of the target task: the base with the tree of layup-setup-2 copied in, and the one change
#                  that VARIANT names
#
# Usage: sh tree-equal-test.sh [TREE-EQUAL-SCRIPT]
# Exit 0 when each case behaves as written, else 1.
here=$(cd "$(dirname "$0")" && pwd)
TE=${1:-$here/tree-equal.sh}
T=$(mktemp -d) || exit 2
trap 'rm -rf "$T"' EXIT INT TERM
g() { git -C "$repo" -c user.name=fixture -c user.email=fixture@example.invalid -c commit.gpgsign=false \
	-c core.hooksPath=/dev/null -c maintenance.auto=false "$@"; }

LOGHEAD='# completed

<!-- Most recent first. -->

'
A0RTLINE='- **2026-10-05** — **T-a0rt** — the ruleset is applied'
TASKLINE='- **2026-10-06** — **T-vu2j** — the tree of the second setup run'
OTHERLINE='- **2026-10-06** — **T-zz99** — another task'

# log LINE...: the completed log with these lines, most recent first
log() { { printf '%s' "$LOGHEAD"; for l in "$@"; do printf '%s\n' "$l"; done; } > "$repo/docs/tasks/completed.md"; }
task() { printf '# T-vu2j\n' > "$repo/docs/tasks/T-vu2j.md"; }
evidence() { mkdir -p "$repo/runs/T-vu2j" && printf 'red, then green\n' > "$repo/runs/T-vu2j/evidence.md"; }

build() {
	repo=$T/$1; v=$2
	mkdir -p "$repo/docs/tasks" && git init -q -b main "$repo"
	printf 'a rule\n' > "$repo/docs/rule.md"
	printf '# onboarding\n| gate | no |\n' > "$repo/docs/onboarding-for-engineers.md"
	log
	[ "$v" = task-in-setup ] && printf '# an older file\n' > "$repo/docs/tasks/T-vu2j.md"
	g add -A && g commit -q -m "chore: setup S14"
	g branch -q first-setup
	# the base: the change of T-a0rt
	printf '# onboarding\n| gate | yes |\n' > "$repo/docs/onboarding-for-engineers.md"
	log "$A0RTLINE"
	mkdir -p "$repo/runs/T-a0rt" && printf 'evidence\n' > "$repo/runs/T-a0rt/evidence.md"
	printf '# T-a0rt\n' > "$repo/docs/tasks/T-a0rt.md"
	[ "$v" = task-in-base ] && task
	g add -A && g commit -q -m "docs: T-a0rt"
	g branch -q base
	# the second setup head: the first one with a changed rule
	g switch -q -c layup-setup-2 first-setup
	printf 'a rule of the second run\n' > "$repo/docs/rule.md"
	g add -A && g commit -q -m "chore: setup S14, second run"
	# the head: the copy, then VARIANT
	g switch -q main
	[ "$v" = not-copied ] || printf 'a rule of the second run\n' > "$repo/docs/rule.md"
	case "$v" in
	other-file) printf 'a changed rule\n' > "$repo/docs/rule.md" ;;
	extra-file) printf 'more\n' > "$repo/docs/extra.md" ;;
	lost-file) rm "$repo/docs/rule.md" ;;
	mode) chmod +x "$repo/docs/rule.md" ;;
	a0rt-text) printf '# onboarding\n| gate | no |\n' > "$repo/docs/onboarding-for-engineers.md" ;;
	a0rt-lost) rm "$repo/docs/tasks/T-a0rt.md" ;;
	a0rt-evidence) printf 'other evidence\n' > "$repo/runs/T-a0rt/evidence.md" ;;
	a0rt-mode) chmod +x "$repo/runs/T-a0rt/evidence.md" ;;
	log-lost) log ;;
	log-mode) chmod +x "$repo/docs/tasks/completed.md" ;;
	log-other) log "$OTHERLINE" "$A0RTLINE" ;;
	task) task ;;
	task-line) task; log "$TASKLINE" "$A0RTLINE" ;;
	task-two-lines) task; log "$TASKLINE" "$TASKLINE" "$A0RTLINE" ;;
	task-line-edited) task; log "$TASKLINE" '- **2026-10-05** — **T-a0rt** — the ruleset is not applied' ;;
	task-in-setup) printf '# T-vu2j\n' > "$repo/docs/tasks/T-vu2j.md" ;;
	task-evidence) task; evidence; log "$TASKLINE" "$A0RTLINE" ;;
	esac
	g add -A && g commit -q -m "chore: T-vu2j main"
}

n=0; bad=0; HEADREF=main
# run NAME VARIANT WANT_EXIT WANT_MESSAGE [ARGUMENT...]: the arguments after REPO HEAD SETUP2
run() {
	n=$((n + 1)); build "$1" "$2"; name=$1; we=$3; wm=$4; shift 4
	out=$(sh "$TE" "$repo" "$HEADREF" layup-setup-2 "$@" 2>&1); code=$?
	if [ "$code" = "$we" ] && printf '%s' "$out" | grep -Fq -- "$wm"; then echo "ok   $n $name"; else
		bad=$((bad + 1)); echo "FAIL $n $name: exit $code (want $we); want \"$wm\""; printf '%s\n' "$out" | sed 's/^/       | /' | head -12
	fi
}
TF='--task-file docs/tasks/T-vu2j.md'
EF='--evidence-file runs/T-vu2j/evidence.md'

# The copy of the second run's tree (the parity of the tree).
run equal          none          0 'tree-equal: PASS'                  --base base
run not-copied     not-copied    1 'FAIL: docs/rule.md (M)'            --base base
run other-file     other-file    1 'FAIL: docs/rule.md (M)'            --base base
run extra-file     extra-file    1 'FAIL: docs/extra.md (A)'           --base base
run lost-file      lost-file     1 'FAIL: docs/rule.md (D)'            --base base
run mode           mode          1 'FAIL: docs/rule.md (M)'            --base base
# The change of T-a0rt keeps the text of the base.
HEADREF=layup-setup-2
run setup2-head    none          1 'tree-equal: FAIL (findings: 4)'    --base base
run setup2-head-a  none          1 'FAIL: T-a0rt: docs/onboarding-for-engineers.md differs from the base' --base base
HEADREF=main
run a0rt-text      a0rt-text     1 'FAIL: T-a0rt: docs/onboarding-for-engineers.md differs from the base' --base base
run a0rt-lost      a0rt-lost     1 'FAIL: T-a0rt: docs/tasks/T-a0rt.md is not in the head' --base base
run a0rt-evidence  a0rt-evidence 1 'FAIL: T-a0rt: runs/T-a0rt/evidence.md differs from the base' --base base
run a0rt-mode      a0rt-mode     1 'FAIL: T-a0rt: runs/T-a0rt/evidence.md differs from the base' --base base
run log-lost       log-lost      1 'FAIL: T-a0rt: docs/tasks/completed.md differs from the base' --base base
run log-mode       log-mode      1 'FAIL: T-a0rt: docs/tasks/completed.md: its mode differs from the base' --base base
run log-other      log-other     1 'FAIL: T-a0rt: docs/tasks/completed.md differs from the base, and it has no line of T-vu2j' --base base $TF
# The files of the task, and its line in the completed log.
run task           task          0 'task file: docs/tasks/T-vu2j.md (new)' --base base $TF
run task-no-line   task          0 'T-a0rt: docs/tasks/completed.md as in the base' --base base $TF
run task-log-line  task          1 'FAIL: T-a0rt: docs/tasks/completed.md has no line of T-vu2j' --base base $TF --log-line
run task-line      task-line     0 'T-a0rt: docs/tasks/completed.md as in the base, with the line of T-vu2j' --base base $TF --log-line
run task-line-pass task-line     0 'tree-equal: PASS'                  --base base $TF --log-line
run task-two-lines task-two-lines 1 'FAIL: T-a0rt: docs/tasks/completed.md has 2 lines of T-vu2j' --base base $TF
run task-edited    task-line-edited 1 'FAIL: T-a0rt: docs/tasks/completed.md differs from the base in more than the line of T-vu2j' --base base $TF
run task-unnamed   task          1 'FAIL: docs/tasks/T-vu2j.md (A)'    --base base
run task-absent    none          1 'FAIL: docs/tasks/T-vu2j.md: the task file is not in the head' --base base $TF
run task-in-setup  task-in-setup 1 'FAIL: docs/tasks/T-vu2j.md: the task file is in layup-setup-2, so it is not a new file' --base base $TF
run task-in-base   task-in-base  1 'FAIL: docs/tasks/T-vu2j.md: the task file is in the base, so it is not a new file' --base base $TF
run evidence       task-evidence 0 'evidence file: runs/T-vu2j/evidence.md (new)' --base base $TF $EF --log-line
run evidence-pass  task-evidence 0 'tree-equal: PASS'                  --base base $TF $EF --log-line
run evid-unnamed   task-evidence 1 'FAIL: runs/T-vu2j/evidence.md (A)' --base base $TF --log-line
run evid-absent    task-line     1 'FAIL: runs/T-vu2j/evidence.md: the evidence file is not in the head' --base base $TF $EF
# Input errors.
run no-base        none          2 'input error: --base is required'
run base-no-a0rt   none          2 'input error: the base has no docs/tasks/T-a0rt.md' --base first-setup
run task-a0rt      none          2 'input error: the task file is one of the four paths of T-a0rt' --base base --task-file docs/tasks/T-a0rt.md
run task-form      none          2 'input error: the task file must be docs/tasks/T-' --base base --task-file docs/rule.md
run task-noval     none          2 'input error: --task-file needs a path' --base base --task-file
run evid-other     none          2 'input error: the evidence file must be runs/T-vu2j/evidence.md' --base base $TF --evidence-file runs/T-zz99/evidence.md
run evid-no-task   none          2 'input error: --evidence-file needs --task-file' --base base $EF
run line-no-task   none          2 'input error: --log-line needs --task-file' --base base --log-line
run unknown        none          2 'input error: an unknown argument: --other' --base base --other
HEADREF=no-such-ref
run head-unknown   none          2 'input error: no-such-ref, layup-setup-2 or base is not a commit' --base base
HEADREF=main

echo "tree-equal-test: $((n - bad)) of $n cases behave as written"
[ "$bad" = 0 ]
