#!/bin/sh
# Tests of rules-diff.sh on fixture repositories (task T-evad, #97: D8 (b) of its plan and
# condition 4 of its plan review).
#
# Usage: sh rules-diff-test.sh [RULES-DIFF-SCRIPT]
#
# Each case builds a small target repository: the root commit is the unmodified baseline, and
# each next commit is one setup step ("chore: setup S05", ...). It runs rules-diff.sh on it and
# checks the exit code and one message of the output. A case that must fail also names a message
# that must NOT be in the output, so that each case fails for its own reason and no other.
# Exit 0 when each case behaves as written, else 1.

here=$(cd "$(dirname "$0")" && pwd)
RD=${1:-$here/rules-diff.sh}
SC=${SETUP_CHECK:-$here/../../docs/setup/setup-check.sh}
if [ -n "${KEEP_DIR:-}" ]; then T=$KEEP_DIR; mkdir -p "$T"; else T=$(mktemp -d) || exit 2; trap 'rm -rf "$T"' EXIT INT TERM; fi
LQ=$(printf '\342\200\271')   # the marker characters, never typed in this file
RQ=$(printf '\342\200\272')

sha256() { if command -v sha256sum >/dev/null 2>&1; then sha256sum "$1" | cut -d' ' -f1; else shasum -a 256 "$1" | cut -d' ' -f1; fi; }
g() { git -C "$repo" -c user.name=fixture -c user.email=fixture@example.invalid -c commit.gpgsign=false \
	-c core.hooksPath=/dev/null -c maintenance.auto=false "$@"; }
step() { g add -A && g commit -q -m "chore: setup $1"; }
# row STEP NAME VALUE: one row of out/record.tsv
row() { printf '%s\t%s\t%s\t%s\t%s\n' "$1" "$2" "$3" computed "sha256 inputs/files/$2" >> "$rec"; }
filerow() { row "$1" "file:$2" "$(sha256 "$repo/$2")"; }

# build NAME VARIANT: the clean fixture, with the one change that VARIANT names.
build() {
	repo=$T/$1; v=$2; rec=$T/$1.record.tsv
	mkdir -p "$repo/docs/decisions" "$repo/docs/tasks" "$repo/docs/facts"
	git init -q -b main "$repo"
	printf '# A\nThis kit is for adopters.\nA rule line.\nPort: %sport%s\n' "$LQ" "$RQ" > "$repo/docs/a.md"
	printf '# B\nRule B.\n' > "$repo/docs/b.md"
	printf '# C\nTest runner: %stest runner%s\n' "$LQ" "$RQ" > "$repo/docs/c.md"
	printf '# D\n## How to adapt this kit\n\nFour things need doing.\n\nSecond paragraph.\n## Next section\nKeep this rule.\n' > "$repo/docs/d.md"
	# (two marker lines that check `adapted` does not flag, as in the guardrails of the first pilot)
	if [ "$v" = written-markers ]; then printf '# Project readme\nName: %sname%s\nOwner: %sowner%s\n' "$LQ" "$RQ" "$LQ" "$RQ" > "$repo/README.md"
	else printf '# kit readme\ntext\n' > "$repo/README.md"; fi
	printf 'history\n' > "$repo/docs/decisions/D1.md"
	printf -- '- T-aaaa the task\n- keep this line\nScheme: %sscheme%s\n' "$LQ" "$RQ" > "$repo/docs/tasks/backlog.md"
	printf '# facts\n| index |\n| _none yet_ | | | |\n' > "$repo/docs/facts/README.md"
	g add -A && g commit -q -m "chore: the unmodified baseline at fixture"
	printf 'step\tname\tvalue\tsource\tref\n' > "$rec"

	# S05: the history goes (a deleted rule file when the variant says so)
	git -C "$repo" rm -q docs/decisions/D1.md
	printf -- '- keep this line\nScheme: %sscheme%s\n' "$LQ" "$RQ" > "$repo/docs/tasks/backlog.md"
	[ "$v" = deleted-rule ] && git -C "$repo" rm -q docs/b.md
	[ "$v" = index-grows ] && printf -- '- a new line\n' >> "$repo/docs/tasks/backlog.md"
	step S05
	# S06: the brief and its index row (a row in a rule file when the variant says so)
	printf 'brief\n' > "$repo/docs/facts/brief.md"
	# the placeholder row of the index is replaced by the first row (a line of the index changed in another way when the variant says so)
	if [ "$v" = index-line-changed ]; then printf '# facts changed\n| index |\n| row |\n' > "$repo/docs/facts/README.md"
	else printf '# facts\n| index |\n| row |\n' > "$repo/docs/facts/README.md"; fi
	[ "$v" = insert-in-rule ] && printf 'An inserted rule.\n' >> "$repo/docs/b.md"
	step S06
	# S07: a file that is written for the target
	# (the written text replaces two baseline marker lines, and keeps one marker of its own, which S11 records as a gap)
	if [ "$v" = written-markers ]; then printf '# my project\nnew text\nTimeout: %stimeout%s\n' "$LQ" "$RQ" > "$repo/README.md"
	else printf '# my project\nnew text\n' > "$repo/README.md"; fi
	[ "$v" != no-row-written ] && filerow S07 README.md
	step S07
	# S14: the flagged lines are adapted; the section "How to adapt" is replaced whole (the new text shares the blank line
	# after the heading with the old, so that diff gives the heading and the body as two hunks)
	case "$v" in
	heading-renamed-body-changed)
		printf '# D\n## How this project was set up\n\nFour things need doing.\n\nSecond paragraph changed.\n## Next section\nKeep this rule.\n' > "$repo/docs/d.md" ;;
	section-too-big)
		printf '# D\n## How this project was set up\n\nSee the pin.\n' > "$repo/docs/d.md" ;;
	*)
		printf '# D\n## How this project was set up\n\nSee the pin.\n## Next section\nKeep this rule.\n' > "$repo/docs/d.md" ;;
	esac
	printf '# A\nThis repository is for projects.\nA rule line.\nPort: %sport%s\n' "$LQ" "$RQ" > "$repo/docs/a.md"
	case "$v" in
	rule-changed|rule-changed-with-row) printf '# B\nRule B changed.\n' > "$repo/docs/b.md" ;;
	esac
	[ "$v" != no-row-adapted ] && filerow S14 docs/a.md
	filerow S14 docs/d.md
	[ "$v" = rule-changed-with-row ] && filerow S14 docs/b.md
	if [ "$v" = row-hash-differs ]; then
		printf 'S14\tfile:docs/a.md\t%s\tcomputed\tx\n' "0000000000000000000000000000000000000000000000000000000000000000" >> "$rec"
		awk -F'\t' 'BEGIN{OFS="\t"} $2=="file:docs/a.md" && $3!~/^0+$/ {next} {print}' "$rec" > "$rec.new" && mv "$rec.new" "$rec"
	fi
	step S14
	# S11: a marker is filled
	case "$v" in
	marker-missing) printf '# C\nTest runner: make\n' > "$repo/docs/c.md" ;;
	*) printf '# C\nTest runner: go test\n' > "$repo/docs/c.md" ;;
	esac
	# the row of a marker: marker:<file>:<line>, or marker:<file>:<line>:<column> when the line holds more markers than one
	if [ "$v" = marker-column ]; then row S11 "marker:docs/c.md:2:13" "go test"; else row S11 "marker:docs/c.md:2" "go test"; fi
	# S11 fills the markers of the files that S07 to S14 wrote, so the hash of a file row is the hash at its own step
	printf '# A\nThis repository is for projects.\nA rule line.\nPort: 8080\n' > "$repo/docs/a.md"
	row S11 "marker:docs/a.md:4" "8080"
	printf -- '- keep this line\nScheme: T-xxxx\n' > "$repo/docs/tasks/backlog.md"
	[ "$v" = index-grows ] && printf -- '- a new line\n' >> "$repo/docs/tasks/backlog.md"
	row S11 "marker:docs/tasks/backlog.md:2" "T-xxxx"
	[ "$v" = written-markers ] && printf 'S11\tmarker:README.md:3\t%stimeout%s\tgap\tdocs/setup/open-gaps.tsv\n' "$LQ" "$RQ" >> "$rec"
	if [ "$v" = row-at-head ]; then
		h=$(sha256 "$repo/docs/a.md")
		awk -F'\t' -v h="$h" 'BEGIN{OFS="\t"} $2=="file:docs/a.md" {$3=h} {print}' "$rec" > "$rec.new" && mv "$rec.new" "$rec"
	fi
	step S11
	[ "$v" = foreign-commit ] && { printf 'x\n' >> "$repo/docs/a.md"; g add -A; g commit -q -m "fix: a change that no step made"; }
	return 0
}

n=0; bad=0
# run NAME VARIANT WANT_EXIT WANT_MESSAGE [NOT_WANTED_MESSAGE]
run() {
	n=$((n + 1)); name=$1; variant=$2; wantexit=$3; want=$4; notwant=${5:-}
	build "$name" "$variant"
	out=$(sh "$RD" "$repo" HEAD "$rec" "$SC" 2>&1); code=$?
	ok=1
	[ "$code" = "$wantexit" ] || ok=0
	printf '%s' "$out" | grep -Fq -- "$want" || ok=0
	if [ -n "$notwant" ] && printf '%s' "$out" | grep -Fq -- "$notwant"; then ok=0; fi
	if [ "$ok" = 1 ]; then echo "ok   $n $name"; else
		bad=$((bad + 1)); echo "FAIL $n $name: exit $code (want $wantexit); want \"$want\"${notwant:+, not \"$notwant\"}"
		printf '%s\n' "$out" | sed 's/^/       | /' | head -12
	fi
}

run clean               clean               0 'rules-diff: PASS'
run marker-column       marker-column       0 'rules-diff: PASS'
run row-at-head         row-at-head         1 'record row differs from the file at its step: docs/a.md'
run rule-changed        rule-changed        1 'unflagged baseline line changed: docs/b.md:2' 'rules-diff: PASS'
run rule-changed-row    rule-changed-with-row 1 'unflagged baseline line changed: docs/b.md:2' 'no record row for docs/b.md'
run deleted-rule        deleted-rule        1 'deleted baseline path: docs/b.md'
run no-row-adapted      no-row-adapted      1 'no record row for docs/a.md' 'unflagged baseline line changed'
run no-row-written      no-row-written      1 'no record row for README.md'
run row-hash-differs    row-hash-differs    1 'record row differs from the file at its step: docs/a.md'
run foreign-commit      foreign-commit      1 'a commit that no step made'
run index-grows         index-grows         1 'task index changed other than by removing lines: docs/tasks/backlog.md'
run insert-in-rule      insert-in-rule      1 'lines added to a baseline rule file: docs/b.md'
run marker-missing      marker-missing      1 'recorded marker value is not in the file: docs/c.md'
run section-too-big     section-too-big     1 'unflagged baseline line changed: docs/d.md:7'
run heading-renamed     heading-renamed-body-changed 1 'unflagged baseline line changed: docs/d.md:6'
run index-line-changed  index-line-changed  1 'unflagged baseline line changed: docs/facts/README.md:1'
# the report counts the marker rows of S11, not the baseline marker lines that a written file replaced (finding F-18)
run written-markers     written-markers     0 'MARKERS README.md: S11 marker rows 1 (filled 0, gaps 1); each recorded value is in the file' 'marker lines filled'

echo "rules-diff-test: $((n - bad)) of $n cases behave as written"
[ "$bad" = 0 ]
