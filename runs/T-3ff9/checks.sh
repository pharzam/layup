#!/bin/sh
# Task-only text checks for T-3ff9 (#447); semantic agreement needs review.
set -eu
cd "$(dirname "$0")/../.."
base=f45fffc78070d96c9ab3aab15bbbf1e4e31272bc
failed=0
prose() { awk 'NF { $1=$1; printf "%s%s", sep, $0; sep=" " } END { print "" }'; }
text_of() {
    if [ "$1" = docs/engineering-discipline.md ]; then
        sed -n '/^## Reviewing until findings decay$/,/^### Who may review$/p' "$1" | prose
    else
        prose < "$1"
    fi
}
has() {
    if text_of "$1" | grep -Fq -- "$2"; then
        printf 'ok: %s\n' "$3"
    else
        printf 'FAIL: %s\n' "$3"
        failed=1
    fi
}
absent() {
    if text_of "$1" | grep -Fq -- "$2"; then
        printf 'FAIL: retired %s\n' "$3"
        failed=1
    else
        printf 'ok: retired %s\n' "$3"
    fi
}
same() {
    if [ "$1" = "$2" ]; then
        printf 'ok: %s\n' "$3"
    else
        printf 'FAIL: %s\n' "$3"; failed=1
    fi
}
end_text() { awk -v end="$1" '/^  [1-4]\. |^- \*\*/ { selected=($1 == end ".") } selected' | prose; }
ed=docs/engineering-discipline.md
anchor='engineering-discipline.md#reviewing-until-findings-decay'
has "$ed" 'unit is lines added plus lines removed on the whole branch diff, plus files touched, measured against a named base SHA.' unit
has "$ed" 'Evidence records under `runs/<task>/` leave it, lines and files: Markdown, TSV and text files that record runs, times and outputs.' evidence
has "$ed" 'A file under `runs/` that a check, a test or a script executes or reads as its input (for example `runs/T-gd8q/checks.sh`, or a file that such a script reads) still counts, lines and file.' executable-inputs
has "$ed" 'Executable tests and their inputs stay in the count.' tests-count
has "$ed" '"Input" means the data that a script or a test runs on: executable content, fixtures, expected outputs and configuration used to run or exercise behaviour.' input-definition
has "$ed" 'A record that a check only reads to audit its form, its citations or its evidence is not an input. A file serving both purposes remains an input.' inspection-and-dual-use
# The positive check scopes the definition to the canonical section; this counts
# it across the whole document so a second copy cannot silently pass.
definitions=$(prose < "$ed" | awk -v needle='"Input" means' '{ while (i=index($0, needle)) { n++; $0=substr($0, i+length(needle)) } } END { print n+0 }')
same "$definitions" 1 input-defined-once
has "$ed" "The close-out bookkeeping also leaves the count: the completed-log line, the backlog line it removes, and the verdict and resource record of the task's own file." bookkeeping
has "$ed" 'Where a file holds bookkeeping and other changes, only the bookkeeping lines leave the count. The file still counts as touched.' mixed-files
has "$ed" 'A landing figure above the `Budget maximum` by no more than the tolerance is approved in advance.' delegation
has "$ed" 'The tolerance is 20 % of the maximum’s lines, at least 60 lines and at most 500, plus 4 files.' tolerance
has "$ed" 'Before the merge, the author posts the figure, the maximum and the margin on the issue. The Operator can overrule that delegated decision there before the merge.' budget-post-and-overrule
has "$ed" 'The tolerance is one fixed bound, the same for every task. It is not a multiple, does not grow with the findings, and is not a second approval.' fixed-bound
has "$ed" 'A figure past the tolerance is a finding reported on the issue, never a revision. The growth becomes a child issue unless the Operator approves one exact landing figure, once, on the issue. That approval reaches the named figure and no further.' exact-approval
has "$ed" 'An unapproved overrun blocks the merge: the last round carries it as a finding and returns `not mergeable, findings recorded`. This route does not run through the cap.' budget-block
has "$ed" '**Once** counts per issue. A successor issue starts with an approval of its own, and its plan review sets its maximum with the carried size already measured. This is a discipline, not a mechanism.' successor
has "$ed" 'Growth past the approved figure is absorbed, routed to a child issue where the finding is non-material, or ends at the cap.' growth
has "$ed" 'or the landing figure exceeds the `Budget maximum` plus its tolerance in either dimension and the Operator has not approved that exact landing figure on the issue.' verdict-budget
has "$ed" 'is fixed here when the resulting landing figure is permitted by the budget rule below; otherwise it becomes a child issue.' finding-budget
has "$ed" '**The ends at the cap.** With material findings still open at the cap, the author posts them on the issue. The Operator chooses one of four ends, except where end 3 applies by default below.' cap-intro
has "$ed" 'task record names it as a fix that no round read. The last verdict stays `not mergeable, findings recorded`.' unread-fix-and-verdict
has "$ed" 'At the cap with material findings still open, this end applies without an Operator comment when every open material finding is a records finding.' default-end-3
has "$ed" 'A records finding is one whose fix restores a named source or a record’s form, and changes nothing else.' records-definition
has "$ed" 'Examples are a traceability row or Test cell made to name a test that exists, an index line, record columns made to match their schema, and trailing white space.' records-examples
has "$ed" 'A fix that changes an exit code, an assertion, a verdict, a check’s input (as defined in the budget rule below), a test, product code, the specification, or a sentence that grants or limits an authority is not a records finding.' records-exclusions
has "$ed" 'Its end stays the Operator’s. A disputed classification also stays the Operator’s.' operator-end
has "$ed" 'The author records `end 3 by ADR-0028` with the findings and their fixes on the issue. The Operator can overrule it there before the merge. Either this default or the Operator’s comment permits the merge.' end-post-and-overrule
has "$ed" 'This end is the one exception to “A fix re-freezes”, “The last round carries the verdict”, and the `Verdict` field of [What a round records](#what-a-round-records). After that fix, no further fix or round follows; the close-out bookkeeping and a merge of `main` land as before.' end-exception
absent "$ed" 'overran without the operator' old-verdict
absent "$ed" 'fixed here inside the budget' old-finding
absent "$ed" 'no kind of growth an approval covers in advance' old-budget
absent "$ed" '**The Operator ends the work at the cap.**' old-cap
absent "$ed" 'An approval is **one number, named once,' old-approval
absent "$ed" "only the Operator's comment makes it" old-end-3
for end in 1 2 4; do
    same "$(git show "$base:$ed" | end_text "$end")" "$(end_text "$end" < "$ed")" "end-$end-preserved"
done
has AGENTS.md "The [review protocol](docs/$anchor) defines the cycles, verdicts and ends at the cap." entry-point
absent AGENTS.md 'the Operator chooses the end at the cap' old-entry-point
has docs/tests/dod-checklist.md "The landing figure and any overrun follow [Reviewing until findings decay](../$anchor)." dod-pointer
has docs/tests/dod-checklist.md "A budget that was never declared is the plan review's gap, not this task's." dod-gap
absent docs/tests/dod-checklist.md 'The budget declared in the plan-review confirmation was **held**' old-dod-budget
has docs/plan/README.md "[Reviewing until findings decay](../$anchor) defines the budget and the authority over a landing figure." plan-pointer
absent docs/plan/README.md 'a task that grows past its budget goes back to its issue as a new slice or a child issue' old-plan-budget
has docs/glossary.md "| Cycle cap | \`—\` | The bound on fix-and-review cycles after a task's first frozen head. Its value, verdicts and ends are defined in [Reviewing until findings decay]($anchor). Collision to watch for: a cap bounds revisions; a **ceiling** bounds size. | The cap answers how many cycles remain. |" cap-pointer
has docs/glossary.md "| Ceiling | \`—\` | An upper bound on a change's size. The budget, its tolerance and the authority over a landing figure are defined in [Reviewing until findings decay]($anchor). Collision to watch for: the **cycle cap** bounds revisions rather than size. | The budget answers how large the change may grow. |" ceiling-pointer
has docs/glossary.md "Its counterpart is a defect *in the change*, whose disposition follows [Reviewing until findings decay]($anchor)." revealed-pointer
absent docs/glossary.md "**A thing LAYUP's gate does not have.**" old-ceiling
absent docs/glossary.md 'The maximum number of fix-and-review cycles a task may run' old-cap-row
absent docs/glossary.md 'fixed inside the budget' old-revealed-budget
for adr in 0007-record-task-resource-use 0026-keep-the-bootstrap-review-rules-as-the-standing-gate; do
    path=docs/adr/$adr.md
    has "$path" 'Accepted. Amended by [ADR-0028](0028-delegate-two-bounded-landing-decisions.md)' "$adr-status"
    # Undo only the permitted Status edit; compare bytes, including the body.
    restored=$(sed '/^## Status$/,/^## Context$/s@^Accepted\. Amended by \[ADR-0028\](0028-delegate-two-bounded-landing-decisions.md)$@Accepted@' "$path" | git hash-object --stdin)
    same "$(git rev-parse "$base:$path")" "$restored" "$adr-only-status-changed"
    row=$(grep -F "| [${adr%%-*}](" docs/adr/README.md)
    case "$row" in *'| Accepted; amended |'*) printf 'ok: %s index\n' "$adr" ;; *) printf 'FAIL: %s index\n' "$adr"; failed=1 ;; esac
done
restored=$(sed -e '/^| \[0007\](/s/| Accepted; amended |$/| Accepted |/' -e '/^| \[0026\](/s/| Accepted; amended |$/| Accepted |/' docs/adr/README.md | git hash-object --stdin)
same "$(git rev-parse "$base:docs/adr/README.md")" "$restored" index-only-status-cells-changed
if git diff --quiet "$base" -- docs/ci/review-record-lint.sh docs/tests/freeze-check.sh docs/tests/freeze-check; then
    printf 'ok: parser and K5 scope unchanged\n'
else
    printf 'FAIL: parser or K5 scope changed\n'; failed=1
fi
exit "$failed"
