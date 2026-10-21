#!/bin/sh
# T-ebm4 (#451): scoped text contract; semantic agreement still needs review.
set -eu
LC_ALL=C
export LC_ALL
cd "$(dirname "$0")/../.."
base=abdbe996282b33080406223ae8ba19e5a3634360
ed=docs/engineering-discipline.md
gl=docs/glossary.md
anchor='engineering-discipline.md#reviewing-until-findings-decay'
failed=0
prose() { awk 'NF { $1=$1; printf "%s%s", sep, $0; sep=" " } END { print "" }'; }
section() {
    awk -v start="$2" -v stop="$3" '
        $0 == stop { active=0 }
        $0 == start { active=1; found++ }
        active { print }
        END { if (found != 1) exit 2 }
    ' "$1"
}
paragraph() {
    awk -v lead="$2" 'BEGIN { RS="" }
        index($0, lead) { gsub(/[[:space:]]+/, " "); print; found++ }
        END { if (found != 1) exit 2 }
    ' "$1"
}
row() {
    awk -F'|' -v term="$1" '
        { name=$2; gsub(/^ +| +$/, "", name) }
        name == term { print; found++ }
        END { if (found != 1) exit 2 }
    ' "$gl"
}
has() {
    if printf '%s\n' "$scope" | grep -Fq -- "$2"; then
        printf 'ok: %s\n' "$1"
    else
        printf 'FAIL: %s (missing required text)\n' "$1"; failed=1
    fi
}
absent() {
    if printf '%s\n' "$scope" | grep -Fq -- "$2"; then
        printf 'FAIL: %s (retired text remains)\n' "$1"; failed=1
    else
        printf 'ok: %s\n' "$1"
    fi
}
same() {
    if [ "$2" = "$3" ] && [ -n "$2" ]; then
        printf 'ok: %s\n' "$1"
    else
        printf 'FAIL: %s (protected text changed or absent)\n' "$1"; failed=1
    fi
}
end_body() { awk -v n="$1" '/^  [1-4]\. |^- \*\*/ { active=($1 == n ".") } active' | prose; }
# Fail on missing source files before any pipeline can conceal a read failure.
for path in "$ed" "$gl" docs/issue-workflow.md docs/guardrails.md; do
    test -r "$path" || { printf 'ERROR: cannot read %s\n' "$path"; exit 2; }
done
canonical=$(section "$ed" '## Reviewing until findings decay' '### Who may review')
scope=$(printf '%s\n' "$canonical" | prose)
has budget-unit 'unit is lines added plus lines removed on the whole branch diff, plus files touched, measured against a named base SHA.'
has evidence-exclusion 'Evidence records under `runs/<task>/` leave it, lines and files: Markdown, TSV and text files that record runs, times and outputs.'
has executable-inputs 'A file under `runs/` that a check, a test or a script executes or reads as its input (for example `runs/T-gd8q/checks.sh`, or a file that such a script reads) still counts, lines and file. Executable tests and their inputs stay in the count.'
has input-definition '"Input" means the data that a script or a test runs on: executable content, fixtures, expected outputs and configuration used to run or exercise behaviour.'
has audit-and-dual-use 'A record that a check only reads to audit its form, its citations or its evidence is not an input. A file serving both purposes remains an input.'
has bookkeeping "The close-out bookkeeping also leaves the count: the completed-log line, the backlog line it removes, and the verdict and resource record of the task's own file."
has mixed-files 'Where a file holds bookkeeping and other changes, only the bookkeeping lines leave the count. The file still counts as touched.'
has advance-approval 'A landing figure above the `Budget maximum` by no more than the tolerance is approved in advance.'
has tolerance "The tolerance is 20 % of the maximum's lines, at least 60 lines and at most 500, plus 4 files."
has margin 'The margin is the maximum plus the tolerance, minus the landing figure, calculated separately in lines and in files.'
has budget-post 'Before the merge, the author posts the figure, the maximum and the margin on the issue. The Operator can overrule that delegated decision there before the merge.'
has fixed-tolerance 'The tolerance is one fixed bound, the same for every task. It is not a multiple, does not grow with the findings, and is not a second approval.'
has exact-approval 'A figure past the tolerance is a finding reported on the issue, never a revision. The growth becomes a child issue unless the Operator approves one exact landing figure, once, on the issue. That approval reaches the named figure and no further.'
has overrun-block 'An unapproved overrun blocks the merge: the last round carries it as a finding and returns `not mergeable, findings recorded`. This route does not run through the cap.'
has end-3-eligibility 'At the cap with material findings still open, this end applies without an Operator comment when every open material finding is a records finding.'
has records-definition "A records finding is one whose fix restores a named source or a record's form, and changes nothing else."
has records-exclusions "A fix that changes an exit code, an assertion, a verdict, a check's input (as defined in the budget rule below), a test, product code, the specification, or a sentence that grants or limits an authority is not a records finding."
has disputed-end "Its end stays the Operator's. A disputed classification also stays the Operator's."
has end-3-post 'The author records `end 3 by ADR-0028` with the findings and their fixes on the issue.'
has end-3-authority "Either this default or the Operator's comment permits the merge."
has end-3-exception 'This end is the one exception to "A fix re-freezes", "The last round carries the verdict", and the `Verdict` field of [What a round records](#what-a-round-records). After that fix, no further fix or round follows; the close-out bookkeeping and a merge of `main` land as before.'
if printf '%s\n' "$canonical" | grep -Eq '‘|’|“|”'; then
    printf 'FAIL: canonical-straight-quotes (curly quotation marks remain)\n'; failed=1
else
    printf 'ok: canonical-straight-quotes\n'
fi
# The shared wait belongs to the budget bullet, not a second end-3 rule.
scope=$(printf '%s\n' "$canonical" | awk '/^- \*\*The budget is R12/ { active=1 } active' | prose)
has C1-new-comment-wait "For a delegated figure or end, the author posts a new issue comment that names it; the merge waits 10 minutes from that new comment's UTC posting time when that posting time is from 05:00 inclusive to 19:00 exclusive UTC (08:00 to 22:00 UTC+3), and may follow the new comment with no wait when its posting time is from 19:00 inclusive to 05:00 exclusive UTC (22:00 to 08:00 UTC+3); an edit of an earlier comment posts nothing."
absent C1-retired-wait "For a delegated figure or end posted in an issue comment that names it, the merge waits 10 minutes from that comment's UTC posting time"
waits=$(prose < "$ed" | awk -v needle='For a delegated figure or end,' '{ while (p=index($0, needle)) { n++; $0=substr($0, p+length(needle)) } } END { print n+0 }')
same C1-one-shared-wait 1 "$waits"
for n in 1 2 4; do
    original=$(git show "$base:$ed")
    same "end-$n-preserved" "$(printf '%s\n' "$original" | end_body "$n")" "$(printf '%s\n' "$canonical" | end_body "$n")"
done
scope=$(paragraph "$ed" '5. **Review until findings decay.**')
has gate-pointer 'After the code works, freeze the head and run one independent blind review round on it. The [review protocol](#reviewing-until-findings-decay) defines the cycles, verdicts and ends at the cap.'
absent gate-retired-next-round 'A fix re-freezes, and one more round reads it;'
scope=$(section docs/issue-workflow.md '## R12 — Slice and prioritize' '## R13 — One reading, not two' | prose)
has R12-budget-pointer "growth past the scale the plan implied follows [the budget]($anchor), not the strength of the earlier confirmation."
absent R12-retired-slicing 'grows past the scale the plan implied goes back to the issue as a new slice or a child sub-issue'
scope=$(paragraph docs/guardrails.md '**A merge of `main` moves the base of the budget.**')
has named-base-pointer "The plan review names a base SHA; [the budget]($anchor) defines the count against it."
absent named-base-retired-count 'the budget counts the lines added plus removed on the whole branch diff against it.'
# C2: protect the entire validation section, not merely the old pass-cell token.
original=$(git show "$base:docs/guardrails.md")
validation=$(printf '%s\n' "$original" | sed -n '/^## 3\. Validation /,/^## 4\. Mechanics$/p')
same C2-validation-preserved "$validation" "$(sed -n '/^## 3\. Validation /,/^## 4\. Mechanics$/p' docs/guardrails.md)"
scope=$(row "Size class (of LAYUP's plan)")
has size-class-pointer "not the \`Budget maximum\` of a plan review ([Reviewing until findings decay]($anchor))."
absent size-class-retired-ceiling 'of a plan review, which is a ceiling.'
scope=$(row 'Issue split')
has issue-split-pointer "An end at the cap defined in [Reviewing until findings decay]($anchor), which sets what the successor and any child issues carry."
has issue-split-example "A successor takes the stopped branch's work and finding F1; a child issue takes finding F2."
absent issue-split-retired-overrun 'an unapproved overrun is always among the findings it takes'
absent issue-split-retired-example 'A branch stopped at its cap is not abandoned'
scope=$(row Material)
has material-pointer "The material test and review outcomes are defined in [Reviewing until findings decay]($anchor); this row is a non-operative summary."
absent material-retired-stop 'which stop on `nothing material in scope`'
scope=$(row 'Frozen head')
has frozen-head-pointer "The changes allowed after it and their effect on review are defined in [Reviewing until findings decay]($anchor)."
absent frozen-head-retired-next-round 'the first two produce a frozen head a next round names'
# Missing rows are acceptance failures (exit 1), not extraction errors (exit 2).
for term in 'Records finding' Tolerance Input Margin 'Coordinated Universal Time'; do
    scope=$(row "$term") || scope=''
    has "C3-$term-row" "| $term |"
    case "$term" in
    'Records finding')
        has C3-records-bound "A finding whose fix restores a named source or a record's form, and changes nothing else, defined with its exclusions in [Reviewing until findings decay]($anchor)."
        has C3-records-collision 'a finding in a record is not necessarily a records finding.'
        has records-example 'A traceability row corrected to name a test that exists, without changing the test or any other behaviour.' ;;
    Tolerance)
        has tolerance-pointer "The allowance above a task's \`Budget maximum\`, defined in [Reviewing until findings decay]($anchor)."
        has tolerance-collision "this is a size allowance, not a judgement about a finding's materiality."
        has tolerance-example 'A maximum of 650 lines over 20 files has a tolerance of 130 lines and 4 files.' ;;
    Input)
        has input-pointer "Execution data in the budget and records-finding rules; its definition is in [Reviewing until findings decay]($anchor)."
        has input-collision 'reading a record to audit it does not alone make it execution data.'
        has input-example 'A test fixture is input; a log read only to audit evidence is not; a file used for both remains input.' ;;
    Margin)
        has margin-pointer "The remaining room to the budget's limit, defined in [Reviewing until findings decay]($anchor)."
        has margin-collision 'it is measured from the maximum plus tolerance.'
        has margin-example 'For maximum 650/20, tolerance 130/4 and figure 700/20, the margin is 80 lines and 4 files.' ;;
    'Coordinated Universal Time')
        has C3-UTC-abbreviation '| `UTC` |'
        has C3-UTC-description "The time standard used for the posting time and waiting window in [Reviewing until findings decay]($anchor)."
        has C3-UTC-collision 'UTC+3 is a fixed offset three hours ahead of UTC, not UTC itself.'
        has C3-UTC-example '05:00 UTC is 08:00 UTC+3; 19:00 UTC is 22:00 UTC+3.' ;;
    esac
done
scope=$(row 'Cycle cap')
has cap-example 'With `Cycle cap` 1, the first round has `Cycle` 0 and a round after one fix has `Cycle` 1.'
absent cap-retired-example 'The cap answers how many cycles remain.'
scope=$(row Ceiling)
has ceiling-example 'With a maximum of 650 lines over 20 files, a figure of 700 lines over 20 files is within the tolerance.'
absent ceiling-retired-example 'The budget answers how large the change may grow.'
if git diff --quiet "$base" -- runs/T-3ff9/checks.sh docs/tests/freeze-check.sh docs/tests/freeze-check docs/ci/review-record-lint.sh .githooks .github; then
    printf 'ok: existing-checks-and-gates-preserved\n'
else
    printf 'FAIL: existing-checks-and-gates-preserved\n'; failed=1
fi
exit "$failed"
