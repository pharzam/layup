#!/bin/sh
# The check of step T of task T-57vx (#346): each part of the specification of
# M2e is in exactly one row of "The tasks of M2e" in docs/plan/README.md, or
# named with a reason (parts.tsv: file, heading, sub, row). Copied from
# runs/T-bng8/parts.sh and changed for M2e, with conditions 2 and 5 of the plan
# review of #346: each clause of a test, a heading, a block or a table row that
# the slicing cuts between rows has its own entry, a part named in a second Task
# cell fails, and the issues are read by issues.sh. A one-off evidence command of
# this task, not a check of the gate. Run from the root of the repository; exit
# 0 when each rule holds, 1 otherwise.
tab=runs/T-57vx/parts.tsv
plan=${PLAN:-docs/plan/README.md}
spec=docs/spec/loop.md
acc='## The acceptance tests of M2e'
out=$(mktemp)
fail=0
# 1. Each heading of loop.md outside a fenced block is in the list.
awk '/^```/{f=!f; next} !f && (/^## / || /^### /)' "$spec" | while IFS= read -r h; do
	grep -Fq "$spec	$h	" "$tab" || echo "FAIL a heading of $spec has no part: $h"
done >> "$out"
# 2. Each part (file, heading, sub) is listed once.
tail -n +2 "$tab" | cut -f1-3 | sort | uniq -d | sed 's/^/FAIL a part is listed twice: /' >> "$out"
# 3. Each row of the list is a row of the table "The tasks of M2e", and each row
#    of that table has a part.
rows=$(awk '/^## The tasks of M2e/{t=1; next} t && /^## /{t=0} t && /^\| [0-9]+[a-z]? \|/{print $2}' "$plan")
[ -n "$rows" ] || echo "FAIL $plan has no table The tasks of M2e" >> "$out"
listed=$(tail -n +2 "$tab" | cut -f4 | grep -v '^none' | sort -u)
printf '%s\n' "$listed" | while read -r r; do
	printf '%s\n' "$rows" | grep -qx "$r" || echo "FAIL row $r is not a row of The tasks of M2e"
done >> "$out"
printf '%s\n' "$rows" | grep . | while read -r r; do
	printf '%s\n' "$listed" | grep -qx "$r" || echo "FAIL row $r of The tasks of M2e has no part"
done >> "$out"
# 4. Each row of the acceptance tests of loop.md is a part, whole or by its
#    clauses ("<Part>: <clause>").
awk -v a="$acc" '$0==a{t=1; next} t && /^## /{t=0} t && /^\| / && !/^\| (Test|-)/' "$spec" |
	cut -d'|' -f2 | sed 's/^ *//; s/ *$//' | while IFS= read -r p; do
	awk -F'\t' -v f="$spec" -v h="$acc" -v p="$p" \
		'$1==f && $2==h && ($3==p || index($3, p ": ")==1){n++} END{exit !n}' "$tab" ||
		echo "FAIL a row of $acc has no part: $p"
done >> "$out"
# 5. Each clause that the slicing cuts between rows has its own entry, so that
#    one surviving clause does not cover another (conditions 2 of #346).
for want in \
	"docs/spec/loop.md	### The tier and the context of a task	the record kinds of the step table" \
	"docs/spec/loop.md	### The tier and the context of a task	the tier selection and the construction of the prompt file" \
	"docs/spec/loop.md	### The artifact that is a comment	the post, the copy and the kept file" \
	"docs/spec/loop.md	### The transition table	the default rows, the lists and the check of a table" \
	"docs/spec/loop.md	### The transition table	the replacement by set transitions and its history" \
	"docs/spec/loop.md	### The checks of a handoff	the schema, the artifacts, the items and the pair" \
	"docs/spec/loop.md	## REQ-007 — The four \`layup/\` statuses	the statuses, their sources and their posting" \
	"docs/spec/loop.md	## REQ-007 — The four \`layup/\` statuses	the permissions of the App" \
	"docs/spec/loop.md	### layup/rules	the check (the basic rule)" \
	"docs/spec/loop.md	### layup/rules	the status" \
	"docs/spec/loop.md	### The gate jobs read the base	part 1: the jobs read the base and the exception for a batch branch" \
	"docs/spec/loop.md	### The gate jobs read the base	part 2: the code-owner review" \
	"docs/spec/loop.md	### The gate jobs read the base	part 3: the test" \
	"docs/spec/loop.md	## REQ-012 — Shape	the filter and its list of constraints" \
	"docs/spec/loop.md	## REQ-012 — Shape	the step: the options, the rewrite, the stop with no survivor" \
	"docs/spec/loop.md	## REQ-012 — The bet	the check of the brief, the appetite rule and the answer lines" \
	"docs/spec/loop.md	## REQ-012 — The bet	the step: the brief, the answer, bets.tsv, the option choices, the architecture tasks" \
	"docs/spec/loop.md	### The revision of the draft	the check of the quote's fact" \
	"docs/spec/loop.md	### The revision of the draft	the session of the revision" \
	"docs/spec/loop.md	### The rendered commit and the confirmed inventory	the rendering and the versions of the inventory" \
	"docs/spec/loop.md	### The rendered commit and the confirmed inventory	the expected prd-lint failure" \
	"docs/spec/loop.md	### The rendered commit and the confirmed inventory	the push of the rendered commit" \
	"docs/spec/loop.md	## REQ-005 — The steps of the loop after the Scaffold	the step shape" \
	"docs/spec/loop.md	## REQ-005 — The steps of the loop after the Scaffold	the step bet" \
	"docs/spec/loop.md	## REQ-005 — The steps of the loop after the Scaffold	the step loop, the restart, the push with no result, the comments of a task session" \
	"docs/spec/loop.md	### The demo and the release review	the demo" \
	"docs/spec/loop.md	### The demo and the release review	the review of the release" \
	"docs/spec/loop.md	## The acceptance tests of M2e	A handoff is valid only by its table (inventory)" \
	"docs/spec/loop.md	## The acceptance tests of M2e	The artifact that is a comment (\`gov-issue-68\`): the unit part" \
	"docs/spec/loop.md	## The acceptance tests of M2e	The artifact that is a comment (\`gov-issue-68\`): the integration part" \
	"docs/spec/loop.md	## The acceptance tests of M2e	The transition table: the default rows and the check of a table" \
	"docs/spec/loop.md	## The acceptance tests of M2e	The transition table: the replacement by set transitions" \
	"docs/spec/loop.md	## The acceptance tests of M2e	\`layup spec check\`, the full check (inventory)" \
	"docs/spec/loop.md	## The acceptance tests of M2e	\`layup/gates\` (inventory, \`later-gate-status\`)" \
	"docs/spec/loop.md	## The acceptance tests of M2e	\`layup/rules\`, the basic rule" \
	"docs/spec/loop.md	## The acceptance tests of M2e	\`layup/verify\` and the independence: the authors and the independence" \
	"docs/spec/loop.md	## The acceptance tests of M2e	\`layup/verify\` and the independence: the status" \
	"docs/spec/loop.md	## The acceptance tests of M2e	The ready and merge rule (inventory, \`later-gate-status\`)" \
	"docs/spec/loop.md	## The acceptance tests of M2e	The close-out commit and the rendered commit: the close-out and the carry-over check" \
	"docs/spec/loop.md	## The acceptance tests of M2e	The close-out commit and the rendered commit: the rendered commit's push" \
	"docs/spec/loop.md	## The acceptance tests of M2e	The steps after the Scaffold: the restart, the push with no result and the comments of a session" \
	"docs/spec/loop.md	## The acceptance tests of M2e	The steps after the Scaffold: the step shape that is not done before" \
	"docs/spec/loop.md	## The acceptance tests of M2e	The bet step is atomic to its last act" \
	"docs/spec/loop.md	## The acceptance tests of M2e	The tier and the context of a task (\`session.md\` deferral): the record kinds of the step table" \
	"docs/spec/loop.md	## The acceptance tests of M2e	The tier and the context of a task (\`session.md\` deferral): the tier and the prompt file" \
	"docs/spec/loop.md	## The acceptance tests of M2e	The questions left open and the band (\`intake.md\` deferral)" \
	"docs/spec/loop.md	## The acceptance tests of M2e	The approvers and the success criteria (\`intake.md\` deferral)" \
	"docs/spec/loop.md	## The acceptance tests of M2e	The task register and the Retrospective: the register" \
	"docs/spec/loop.md	## The acceptance tests of M2e	The task register and the Retrospective: the Retrospective" \
	"docs/spec/loop.md	## The acceptance tests of M2e	The Shape filter (inventory): the filter" \
	"docs/spec/loop.md	## The acceptance tests of M2e	The Shape filter (inventory): the rewrite, the round with no progress and the stop" \
	"docs/spec/loop.md	## The acceptance tests of M2e	The Shape filter (inventory): the option answer lines" \
	"docs/spec/loop.md	## The acceptance tests of M2e	The bet and its brief (inventory): the checks of the brief and the appetite" \
	"docs/spec/loop.md	## The acceptance tests of M2e	The bet and its brief (inventory): the role of the approver and the copy of the answer" \
	"docs/spec/loop.md	## The acceptance tests of M2e	The rendered commit (inventory): the rendering and the changes allowed" \
	"docs/spec/loop.md	## The acceptance tests of M2e	The rendered commit (inventory): layup/rules passes a rendered record" \
	"docs/spec/loop.md	## The acceptance tests of M2e	The parameter change (inventory): the refusals and the application of a set" \
	"docs/spec/loop.md	## The acceptance tests of M2e	The parameter change (inventory): the prd-lint head" \
	"docs/spec/loop.md	## The acceptance tests of M2e	The frozen test list" \
	"docs/spec/loop.md	## The acceptance tests of M2e	The gate jobs read the base (F-22): the catalog" \
	"docs/spec/loop.md	## The acceptance tests of M2e	The gate jobs read the base (F-22): the setup" \
	"docs/spec/loop.md	## The acceptance tests of M2e	The task loop with the stand-in forge (inventory)" \
	"docs/spec/loop.md	## The acceptance tests of M2e	The specification task, end to end (inventory)" \
	"docs/spec/loop.md	## The acceptance tests of M2e	The demo (inventory)" \
	"docs/spec/records.md	## NFR-001 — The records of the task loop	the block tasks" \
	"docs/spec/records.md	## NFR-001 — The records of the task loop	the block transitions" \
	"docs/spec/records.md	## NFR-001 — The records of the task loop	the block task-gates" \
	"docs/spec/records.md	## NFR-001 — The records of the task loop	the block task-tests" \
	"docs/spec/records.md	## NFR-001 — The records of the task loop	the block task-spec-check" \
	"docs/spec/records.md	## NFR-001 — The records of the task loop	the block options" \
	"docs/spec/records.md	## NFR-001 — The records of the task loop	the block bets" \
	"docs/spec/records.md	## NFR-001 — The records of the task loop	the blocks inventory and inventory-index" \
	"docs/spec/records.md	## NFR-001 — The records of the task loop	the block parameters" \
	"docs/spec/records.md	## NFR-001 — The records of the task loop	result: the kinds question, decision, lesson and comment" \
	"docs/spec/records.md	## NFR-001 — The records of the task loop	events: the kinds handoff and merged" \
	"docs/spec/records.md	## NFR-001 — The records of the task loop	events: the details of closed and rebased" \
	"docs/spec/records.md	## NFR-001 — The records of the task loop	copies: the copy of a comment that is an artifact" \
	"docs/spec/loop.md	## REQ-005 — The steps of the loop after the Scaffold	run-steps: shape" \
	"docs/spec/loop.md	## REQ-005 — The steps of the loop after the Scaffold	run-steps: bet" \
	"docs/spec/loop.md	## REQ-005 — The steps of the loop after the Scaffold	run-steps: loop" \
	"docs/spec/records.md	### On the records branch (\`records:\`)	tasks/<task>/results/<session>/<file>" \
	"docs/spec/forge.md	## The calls of M2e	—" \
	"docs/spec/packages.md	### The table of M2e	internal/handoff: the transition table and the check of a handoff" \
	"docs/spec/packages.md	### The table of M2e	internal/spec: the Shape filter" \
	"docs/spec/packages.md	### The table of M2e	internal/spec: the checks of the bet and the rendering" \
	"docs/spec/packages.md	### The table of M2e	internal/spec: the full check and the pinned prd-lint, which starts sh" \
	"docs/spec/packages.md	### The table of M2e	internal/rules: the basic layup/rules" \
	"docs/spec/packages.md	### The table of M2e	internal/session: the record kinds of the step table" \
	"docs/spec/packages.md	### The table of M2e	internal/forge and internal/forge/github: the calls of M2e and the permission check" \
	"docs/spec/packages.md	### The table of M2e	internal/catalog: the workflow that reads the base" \
	"docs/spec/packages.md	### The table of M2e	internal/setup: CODEOWNERS and the code-owner rule" \
	"docs/spec/packages.md	### The table of M2e	internal/records: the kinds of result and events" \
	"docs/spec/packages.md	### The table of M2e	internal/run: the May import cell of the table of M2a" \
	"docs/spec/packages.md	### The table of M2e	TestPackageRules reads the table of M2e" \
	"docs/spec/loop.md	### The independence of the reviewers	the authors and the admitted pair" \
	"docs/spec/loop.md	### The independence of the reviewers	the status not-active" \
	"docs/spec/loop.md	### layup/verify	the status" \
	"docs/spec/loop.md	### layup/verify	the redo of a task whose base moved" \
	"docs/spec/loop.md	### The checks of a handoff	the event handoff, the items that need a program, the routing of a question and the accepted answer" \
	"docs/spec/loop.md	### The rendered commit and the confirmed inventory	the rendered task of a later bet" \
	"docs/spec/loop.md	## REQ-012 — The bet	the merge order after a bet" \
	"docs/spec/packages.md	### The table of M2e	internal/run: the register and the transition table" \
	"docs/spec/packages.md	### The table of M2e	internal/run: the parameters" \
	"docs/spec/packages.md	### The table of M2e	internal/run: the handoff in the loop" \
	"docs/spec/packages.md	### The table of M2e	internal/run: the attempts and the frozen test list" \
	"docs/spec/packages.md	### The table of M2e	internal/run: the statuses" \
	"docs/spec/packages.md	### The table of M2e	internal/run: the pull request and the merge" \
	"docs/spec/packages.md	### The table of M2e	internal/run: the step shape" \
	"docs/spec/packages.md	### The table of M2e	internal/run: the step bet" \
	"docs/spec/packages.md	### The table of M2e	internal/run: the step loop" \
	"docs/spec/packages.md	### The table of M2e	internal/handoff: the items that a program computes are the caller's" \
	"docs/spec/loop.md	### The artifact that is a comment	the three equal hashes (the handoff check)" \
	"docs/spec/loop.md	### The artifact that is a comment	the end check of a comment row of the result" \
	"docs/spec/loop.md	### The questions left open	asking the questions, the intent answer and the wait" \
	"docs/spec/loop.md	### The questions left open	a question still open shown in the brief" \
	"docs/spec/loop.md	## The acceptance tests of M2e	The steps after the Scaffold: the tree check that accepts a lost merge event" \
	"docs/spec/loop.md	## The acceptance tests of M2e	The steps after the Scaffold: the step loop that records the lost merge event" \
	"docs/spec/loop.md	## The acceptance tests of M2e	The Shape filter (inventory): two survivors give a brief with no recommendation" \
	"docs/spec/packages.md	### The table of M2e	internal/run: the Retrospective step"; do
	grep -Fq "$want	" "$tab" || echo "FAIL no part: $want" >> "$out"
done
# 6. Each listed part is named in the Task cell of its own row, in the form that
#    gen-issues.py writes it: `<file>`: <heading> (<sub>), and in no other Task cell.
lines=$(awk '/^## The tasks of M2e/{t=1; next} t && /^## /{t=0} t && /^\| [0-9]+[a-z]? \|/' "$plan")
cells=$(printf '%s\n' "$lines" | awk -F' [|] ' '{print $4}')
tail -n +2 "$tab" | grep -v '	none' | while IFS='	' read -r file heading sub row; do
	name=${file##*/}
	h=$(printf '%s' "$heading" | sed 's/^#* *//')
	if [ "$sub" = "—" ]; then text="\`$name\`: $h"; else text="\`$name\`: $h ($sub)"; fi
	line=$(printf '%s\n' "$lines" | awk -v r="$row" '$2==r')
	task=$(printf '%s' "$line" | awk -F' [|] ' '{print $4}')
	case "$task" in
	*"$text"*) ;;
	*) echo "FAIL row $row of The tasks of M2e does not name its part: $text" ;;
	esac
	n=$(printf '%s\n' "$cells" | grep -Fc "$text" || true)
	[ "$n" = "1" ] || echo "FAIL the part is named in $n Task cells of The tasks of M2e, not one: $text"
done >> "$out"
if [ -s "$out" ]; then cat "$out"; fail=1; else echo "ok   each part is in one row of The tasks of M2e, or named with its reason"; fi
rm -f "$out"
exit "$fail"
