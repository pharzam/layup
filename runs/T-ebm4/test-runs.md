# T-ebm4 test runs

Stage A: write the amended text check and run it red on unchanged documents.
Baseline and checked head: `abdbe996282b33080406223ae8ba19e5a3634360`.
No documentation edits, task file, completed-log line or commit belong to this stage.

Command times below come from `date -u +%FT%TZ`; statuses are captured directly,
before printing or recording output. No pipeline supplies a test's status.
Startup read-only inspection preceded the recorder; its full timing was not
captured. Those inspection invocations and their reported statuses are listed
below rather than assigning them invented timestamps. They were not test runs.

## Evidence consumers

The new task check reads the four planned source documents and compares
protected text and existing gates against the named baseline through git.
It does not read any of the four declared evidence records.
The conditions and run records will be audited as evidence by freeze-check,
including K6's table/path checks. That is an audit of form and citations, not
execution data. The other two planned records are created in later stages.

## Commands

### Startup inspection ledger

These are the exact command groups submitted before the recorder was installed.
Each group exited 0 as reported by the command tool. An individual command's
status inside a group was not separately captured. Missing times below are
unavailable, not estimates. The observed timestamps are outputs of `date -u
+%FT%TZ` inside the group, not reconstructed command boundaries.

Group 1. Start: not reported. End: not reported. Group exit status: 0.

```sh
pwd && git status --short && cat AGENTS.md && cat docs/engineering-discipline.md && cat docs/issue-workflow.md
```

Group 2. Start observed: 2026-10-21T07:22:48Z. End: not reported. Group exit status: 0.

```sh
date -u +%FT%TZ
rg --files -g AGENTS.md -g '*.md' -g checks.sh runs/T-ebm4 .review-in docs/adr 2>/dev/null
sed -n '130,460p' docs/engineering-discipline.md
sed -n '460,830p' docs/engineering-discipline.md
```

Group 3. Start observed: 2026-10-21T07:22:58Z. End: not reported. Group exit status: 0.

```sh
date -u +%FT%TZ
sed -n '375,415p' docs/engineering-discipline.md
sed -n '1,220p' docs/glossary.md
sed -n '1,180p' runs/T-3ff9/checks.sh
cat docs/adr/0028-delegate-two-bounded-landing-decisions.md
rg --files -g AGENTS.md -g '!vendor' . docs runs
sed -n '1,120p' docs/guardrails.md
sed -n '625,680p' docs/guardrails.md
sed -n '1235,1270p' docs/guardrails.md
sed -n '1460,1510p' docs/guardrails.md
rg -n 'evidence|conditions|[Kk]6' docs/tests/freeze-check.sh
```

Group 4. Start observed: 2026-10-21T07:23:08Z. End observed: 2026-10-21T07:23:08Z. Group exit status: 0.

```sh
date -u +%FT%TZ
cat runs/T-3ff9/checks.sh
cat docs/adr/0028-delegate-two-bounded-landing-decisions.md
cat docs/adr/0026-keep-the-bootstrap-review-rules-as-the-standing-gate.md
sed -n '86,102p' docs/glossary.md
sed -n '155,170p' docs/glossary.md
sed -n '573,625p' docs/tests/freeze-check.sh
git rev-parse HEAD
find .. -name AGENTS.md -not -path '*/.git/*' -print
date -u +%FT%TZ
```

Group 5. Start observed: 2026-10-21T07:23:34Z. End observed: 2026-10-21T07:23:34Z. Group exit status: 0.

```sh
date -u +%FT%TZ
sed -n '88,100p' docs/glossary.md
sed -n '1,100p' docs/adr/0026-keep-the-bootstrap-review-rules-as-the-standing-gate.md
sed -n '35,95p' docs/adr/0028-delegate-two-bounded-landing-decisions.md
sed -n '100,135p' docs/tests/freeze-check.sh
sed -n '414,445p' docs/engineering-discipline.md
sed -n '1,90p' docs/tests/README.md
date -u +%FT%TZ
```

Group 6. Start observed: 2026-10-21T07:23:56Z. End observed: 2026-10-21T07:23:56Z. Group exit status: 0.

```sh
date -u +%FT%TZ
sed -n '120,360p' docs/guardrails.md
date -u +%FT%TZ
```

Group 7. Start observed: 2026-10-21T07:23:56Z. End observed: 2026-10-21T07:23:56Z. Group exit status: 0.

```sh
date -u +%FT%TZ
sed -n '360,625p' docs/guardrails.md
date -u +%FT%TZ
```

Group 8. Start observed: 2026-10-21T07:24:01Z. End observed: 2026-10-21T07:24:01Z. Group exit status: 0.

```sh
date -u +%FT%TZ
sed -n '680,950p' docs/guardrails.md
date -u +%FT%TZ
```

Group 9. Start observed: 2026-10-21T07:24:01Z. End observed: 2026-10-21T07:24:01Z. Group exit status: 0.

```sh
date -u +%FT%TZ
sed -n '950,1235p' docs/guardrails.md
sed -n '1270,1460p' docs/guardrails.md
date -u +%FT%TZ
```

### `sh -n runs/T-ebm4/checks.sh`

Start: 2026-10-21T07:27:14Z. End: 2026-10-21T07:27:14Z. Exit status: 0.

```text
```

### `sh runs/T-ebm4/checks.sh`

Start: 2026-10-21T07:27:19Z. End: 2026-10-21T07:27:19Z. Exit status: 1.

```text
ok: budget-unit
ok: evidence-exclusion
ok: executable-inputs
ok: input-definition
ok: audit-and-dual-use
ok: bookkeeping
ok: mixed-files
ok: advance-approval
FAIL: tolerance (missing required text)
ok: margin
ok: budget-post
ok: fixed-tolerance
ok: exact-approval
ok: overrun-block
ok: end-3-eligibility
FAIL: records-definition (missing required text)
FAIL: records-exclusions (missing required text)
FAIL: disputed-end (missing required text)
ok: end-3-post
FAIL: end-3-authority (missing required text)
FAIL: end-3-exception (missing required text)
FAIL: canonical-straight-quotes (curly quotation marks remain)
FAIL: C1-new-comment-wait (missing required text)
ok: C1-retired-wait
FAIL: C1-one-shared-wait (protected text changed or absent)
ok: end-1-preserved
ok: end-2-preserved
ok: end-4-preserved
FAIL: gate-pointer (missing required text)
FAIL: gate-retired-next-round (retired text remains)
FAIL: R12-budget-pointer (missing required text)
FAIL: R12-retired-slicing (retired text remains)
FAIL: named-base-pointer (missing required text)
FAIL: named-base-retired-count (retired text remains)
ok: C2-validation-preserved
FAIL: size-class-pointer (missing required text)
FAIL: size-class-retired-ceiling (retired text remains)
FAIL: issue-split-pointer (missing required text)
FAIL: issue-split-example (missing required text)
FAIL: issue-split-retired-overrun (retired text remains)
FAIL: issue-split-retired-example (retired text remains)
FAIL: material-pointer (missing required text)
FAIL: material-retired-stop (retired text remains)
FAIL: frozen-head-pointer (missing required text)
FAIL: frozen-head-retired-next-round (retired text remains)
FAIL: C3-Records finding-row (missing required text)
FAIL: C3-records-bound (missing required text)
FAIL: C3-records-collision (missing required text)
FAIL: records-example (missing required text)
FAIL: C3-Tolerance-row (missing required text)
FAIL: tolerance-pointer (missing required text)
FAIL: tolerance-collision (missing required text)
FAIL: tolerance-example (missing required text)
FAIL: C3-Input-row (missing required text)
FAIL: input-pointer (missing required text)
FAIL: input-collision (missing required text)
FAIL: input-example (missing required text)
FAIL: C3-Margin-row (missing required text)
FAIL: margin-pointer (missing required text)
FAIL: margin-collision (missing required text)
FAIL: margin-example (missing required text)
FAIL: C3-Coordinated Universal Time-row (missing required text)
FAIL: C3-UTC-abbreviation (missing required text)
FAIL: C3-UTC-description (missing required text)
FAIL: C3-UTC-collision (missing required text)
FAIL: C3-UTC-example (missing required text)
FAIL: cap-example (missing required text)
FAIL: cap-retired-example (retired text remains)
FAIL: ceiling-example (missing required text)
FAIL: ceiling-retired-example (retired text remains)
ok: existing-checks-and-gates-preserved
```

The red run reached every assertion and exited 1, not a setup or shell error.
The failures identify the unchanged summaries, the missing five glossary rows
and examples, the missing new-comment waiting sentence and its single-copy
assertion, and the curly quotation characters. The tolerance and records/end-3
sentence failures are the expected straight-quote replacements. Count
exclusions, the input and margin definitions, end-3 eligibility, ends 1, 2 and 4,
the full validation section (condition 2), and existing checks/gates passed.
No source document was changed before this run or in this stage.

### `git diff --exit-code abdbe996282b33080406223ae8ba19e5a3634360 -- . :!runs/T-ebm4`

Start: 2026-10-21T07:29:29Z. End: 2026-10-21T07:29:29Z. Exit status: 0.

```text
```

### `git status --short --untracked-files=all`

Start: 2026-10-21T07:29:29Z. End: 2026-10-21T07:29:29Z. Exit status: 0.

```text
?? runs/T-ebm4/checks.sh
?? runs/T-ebm4/conditions.md
?? runs/T-ebm4/test-runs.md
```

### `git grep -n -F -e runs/T-ebm4/mentions.tsv -e runs/T-ebm4/semantic.md -e runs/T-ebm4/conditions.md -e runs/T-ebm4/test-runs.md -- .githooks docs cmd internal runs`

Start: 2026-10-21T07:29:29Z. End: 2026-10-21T07:29:29Z. Exit status: 1.

```text
```

### `git diff --check`

Start: 2026-10-21T07:29:29Z. End: 2026-10-21T07:29:29Z. Exit status: 0.

```text
```

### `git diff --no-index --check /dev/null runs/T-ebm4/checks.sh`

Start: 2026-10-21T07:29:29Z. End: 2026-10-21T07:29:29Z. Exit status: 1.

```text
```

### `git diff --no-index --check /dev/null runs/T-ebm4/conditions.md`

Start: 2026-10-21T07:29:29Z. End: 2026-10-21T07:29:29Z. Exit status: 1.

```text
```

### `git diff --no-index --check /dev/null runs/T-ebm4/test-runs.md`

Start: 2026-10-21T07:29:30Z. End: 2026-10-21T07:29:30Z. Exit status: 1.

```text
```

The consumer search returned 1 because no tracked file names any of these new
evidence paths. The generic freeze-check consumers were inspected separately:
the evidence-path validation and K5 exclusion in `docs/tests/freeze-check.sh`,
and K6's conditions-table/citation audit in the same file. The new task check
contains no evidence-record consumer.

Each `git diff --no-index --check /dev/null` command returned 1 for an added file,
with no whitespace diagnostic. These statuses are retained rather than called
exit 0. They supplement `git diff --check`, which does not inspect untracked
files. The unchanged-documents comparison returned 0; status lists exactly the
three files named by this stage's commit plan.

### `python3 /home/layup/.cache/layup-session/tmp-green/T-ebm4-stage-a-audit.py`

Start: 2026-10-21T07:30:22Z. End: 2026-10-21T07:30:22Z. Exit status: 0.

```text
ok: only the three Stage A files are changed
ok: final newline and whitespace in runs/T-ebm4/checks.sh
ok: final newline and whitespace in runs/T-ebm4/conditions.md
ok: final newline and whitespace in runs/T-ebm4/test-runs.md
ok: Markdown evidence uses straight quotes
ok: the task check does not consume its evidence records
Counted Stage A addition: 164 lines over 1 file (checks.sh).
```

### `sed -n 535,573p docs/tests/freeze-check.sh`

Start: 2026-10-21T07:31:16Z. End: 2026-10-21T07:31:16Z. Exit status: 0.

```text
check_k5() {
	max_l=${budget%/*}
	max_f=${budget#*/}
	# Disable rename folding so each declared endpoint contributes independently.
	git diff --no-renames --numstat "$base" HEAD > "$tmp/k5_gross"
	: > "$tmp/k5_evidence"
	while IFS= read -r ev; do
		git diff --no-renames --numstat "$base" HEAD -- ":(literal)$ev" >> "$tmp/k5_evidence"
	done < "$tmp/evidence"
	set -- $(awk '
		{ a = ($1 == "-" ? 0 : $1); r = ($2 == "-" ? 0 : $2); l += a + r; f += 1 }
		END { printf "%d %d\n", l + 0, f + 0 }
	' "$tmp/k5_gross")
	fig_l=$1
	fig_f=$2
	set -- $(awk '
		{ a = ($1 == "-" ? 0 : $1); r = ($2 == "-" ? 0 : $2); l += a + r; f += 1 }
		END { printf "%d %d\n", l + 0, f + 0 }
	' "$tmp/k5_evidence")
	fig_l=$((fig_l - $1))
	fig_f=$((fig_f - $2))
	tol_l=$((max_l / 5))
	[ "$tol_l" -ge 60 ] || tol_l=60
	[ "$tol_l" -le 500 ] || tol_l=500
	tol_f=4
	lim_l=$((max_l + tol_l))
	lim_f=$((max_f + tol_f))
	mar_l=$((lim_l - fig_l))
	mar_f=$((lim_f - fig_f))
	k5='ok K5: projection'
	if [ "$fig_l" -gt "$lim_l" ] || [ "$fig_f" -gt "$lim_f" ]; then
		k5='FAIL K5: -:0: projection'
		failed=1
	fi
	printf '%s: figure %d/%d, maximum %d/%d, tolerance %d/%d, limit %d/%d, margin %d/%d (the figure at landing is still measured at landing)\n' \
		"$k5" "$fig_l" "$fig_f" "$max_l" "$max_f" "$tol_l" "$tol_f" "$lim_l" "$lim_f" "$mar_l" "$mar_f"
}

check_k6() {
```

## Stage B: changes and green

This stage implements the plan as amended on #451. The starting committed head
is cff8900583b5fe02d9b58c87c10dbd22ccd7619e; the comparison base remains
abdbe996282b33080406223ae8ba19e5a3634360. The task check is unchanged from Stage A.
The coordinator will commit the Stage B files; this session changes no repository
git metadata. No task file or completed-log line belongs to this stage.

The recorder calls `date -u +%FT%TZ` before and after each command, keeps the
subprocess's actual exit status, and writes output for test runs. Inspection
commands retain command, times and status without duplicating their source text.
Scratch scripts and raw search output are under
/home/layup/.cache/layup-session/tmp-green. Acceptance-check subprocesses use that
directory as TMPDIR, including any temporary fixture repositories. Commands
running longer than ten seconds print a heartbeat.

### Stage B startup inspection ledger

The following read-only commands preceded the Stage B recorder. Each exited 0
as reported by the command tool; UTC start and end were not captured and are
not reported. They are inspection, not acceptance runs. No times are inferred
from wall durations or file dates.

```sh
pwd
git status --short
cat AGENTS.md
cat docs/engineering-discipline.md
cat docs/issue-workflow.md
rg --files -g AGENTS.md -g 'runs/T-ebm4/*' -g 'docs/tasks/T-ebm4.md'
cat runs/T-ebm4/checks.sh runs/T-ebm4/test-runs.md runs/T-ebm4/conditions.md
sed -n '260,470p' docs/engineering-discipline.md
sed -n '470,800p' docs/engineering-discipline.md
cat docs/glossary.md
rg -n '^#|^##|T-ebm4|review|summar|budget|named base|named-base' docs/guardrails.md
cat docs/adr/0028-*.md docs/adr/0026-*.md docs/adr/0007-*.md
git log -3 --format=fuller
cat runs/T-ebm4/conditions.md
sed -n '285,440p' docs/engineering-discipline.md
```

One attempted recorder invocation used incorrect shell quoting and exited 2
before the recorder ran. Start: not reported. End: not reported. No test ran and
no source edit was made by that invocation. The corrected primary search below
is the measured successful search used for the inventory.

```sh
python3 /home/layup/.cache/layup-session/tmp-green/T-ebm4-record.py \"git grep -n -iE 'budget|tolerance|ceiling|end 3|overrun' -- docs/ AGENTS.md > /home/layup/.cache/layup-session/tmp-green/T-ebm4-budget-mentions.txt\"
```

The four evidence records remain the declared mentions.tsv, semantic.md,
conditions.md and test-runs.md. They record occurrences, comparisons, conditions
and actual runs; none is expected output or configuration. The task check's
consumer inspection and the K6 source inspection are recorded below. The
path-specific consumer search exited 1 for no matches outside these records.
K6 will remain enabled at the later clean committed-head pre-freeze run.

### Stage B command

```sh
sed -n '1,360p' docs/guardrails.md
```

Start: 2026-10-21T07:33:54Z. End: 2026-10-21T07:33:54Z. Exit status: 0.

### Stage B command

```sh
sed -n '711,1050p' docs/guardrails.md
```

Start: 2026-10-21T07:33:59Z. End: 2026-10-21T07:33:59Z. Exit status: 0.

### Stage B command

```sh
sed -n '361,710p' docs/guardrails.md
```

Start: 2026-10-21T07:33:59Z. End: 2026-10-21T07:33:59Z. Exit status: 0.

### Stage B command

```sh
sed -n '620,710p' docs/guardrails.md
```

Start: 2026-10-21T07:34:06Z. End: 2026-10-21T07:34:06Z. Exit status: 0.

### Stage B command

```sh
sed -n '1051,1269p' docs/guardrails.md
```

Start: 2026-10-21T07:34:06Z. End: 2026-10-21T07:34:06Z. Exit status: 0.

### Stage B command

```sh
sed -n '1270,1550p' docs/guardrails.md
```

Start: 2026-10-21T07:34:11Z. End: 2026-10-21T07:34:11Z. Exit status: 0.

### Stage B command

```sh
sed -n '1,120p' docs/glossary.md
```

Start: 2026-10-21T07:34:11Z. End: 2026-10-21T07:34:11Z. Exit status: 0.

### Stage B command

```sh
sed -n '1458,1510p' docs/guardrails.md
```

Start: 2026-10-21T07:34:28Z. End: 2026-10-21T07:34:28Z. Exit status: 0.

### Stage B command

```sh
"git
```

Start: 2026-10-21T07:34:28Z. End: 2026-10-21T07:34:28Z. Exit status: 2.

### Stage B command

```sh
sed -n '145,175p' docs/glossary.md
```

Start: 2026-10-21T07:34:28Z. End: 2026-10-21T07:34:28Z. Exit status: 0.

### Stage B command

```sh
sed -n '1,110p' docs/glossary.md
```

Start: 2026-10-21T07:34:28Z. End: 2026-10-21T07:34:28Z. Exit status: 0.

### Stage B command

```sh
git grep -n -iE "budget|tolerance|ceiling|end 3|overrun" -- docs/ AGENTS.md > /home/layup/.cache/layup-session/tmp-green/T-ebm4-budget-mentions.txt
```

Start: 2026-10-21T07:34:42Z. End: 2026-10-21T07:34:42Z. Exit status: 0.

### Stage B command

```sh
sed -n '573,630p' docs/tests/freeze-check.sh
```

Start: 2026-10-21T07:34:42Z. End: 2026-10-21T07:34:42Z. Exit status: 0.

### Stage B command

```sh
python3 /home/layup/.cache/layup-session/tmp-green/T-ebm4-audit.py
```

Start: 2026-10-21T07:34:57Z. End: 2026-10-21T07:34:57Z. Exit status: 0.

### Stage B command

```sh
sh runs/T-ebm4/checks.sh
```

Start: 2026-10-21T07:36:20Z. End: 2026-10-21T07:36:20Z. Exit status: 0.

```text
ok: budget-unit
ok: evidence-exclusion
ok: executable-inputs
ok: input-definition
ok: audit-and-dual-use
ok: bookkeeping
ok: mixed-files
ok: advance-approval
ok: tolerance
ok: margin
ok: budget-post
ok: fixed-tolerance
ok: exact-approval
ok: overrun-block
ok: end-3-eligibility
ok: records-definition
ok: records-exclusions
ok: disputed-end
ok: end-3-post
ok: end-3-authority
ok: end-3-exception
ok: canonical-straight-quotes
ok: C1-new-comment-wait
ok: C1-retired-wait
ok: C1-one-shared-wait
ok: end-1-preserved
ok: end-2-preserved
ok: end-4-preserved
ok: gate-pointer
ok: gate-retired-next-round
ok: R12-budget-pointer
ok: R12-retired-slicing
ok: named-base-pointer
ok: named-base-retired-count
ok: C2-validation-preserved
ok: size-class-pointer
ok: size-class-retired-ceiling
ok: issue-split-pointer
ok: issue-split-example
ok: issue-split-retired-overrun
ok: issue-split-retired-example
ok: material-pointer
ok: material-retired-stop
ok: frozen-head-pointer
ok: frozen-head-retired-next-round
ok: C3-Records finding-row
ok: C3-records-bound
ok: C3-records-collision
ok: records-example
ok: C3-Tolerance-row
ok: tolerance-pointer
ok: tolerance-collision
ok: tolerance-example
ok: C3-Input-row
ok: input-pointer
ok: input-collision
ok: input-example
ok: C3-Margin-row
ok: margin-pointer
ok: margin-collision
ok: margin-example
ok: C3-Coordinated Universal Time-row
ok: C3-UTC-abbreviation
ok: C3-UTC-description
ok: C3-UTC-collision
ok: C3-UTC-example
ok: cap-example
ok: cap-retired-example
ok: ceiling-example
ok: ceiling-retired-example
ok: existing-checks-and-gates-preserved
```

### Stage B command

```sh
sed -n '167,186p' docs/onboarding-for-engineers.md
```

Start: 2026-10-21T07:36:40Z. End: 2026-10-21T07:36:40Z. Exit status: 0.

### Stage B command

```sh
sed -n '750,790p' docs/engineering-discipline.md
```

Start: 2026-10-21T07:36:40Z. End: 2026-10-21T07:36:40Z. Exit status: 0.

### Stage B command

```sh
rg -n 'T-ebm4/(mentions.tsv|semantic.md|conditions.md|test-runs.md)' --glob '!runs/T-ebm4/**' .
```

Start: 2026-10-21T07:36:40Z. End: 2026-10-21T07:36:40Z. Exit status: 1.

### Stage B command

```sh
sed -n '630,670p' docs/tests/freeze-check.sh
```

Start: 2026-10-21T07:36:40Z. End: 2026-10-21T07:36:40Z. Exit status: 0.

### Stage B command

```sh
sed -n '25,43p' README.md
```

Start: 2026-10-21T07:36:40Z. End: 2026-10-21T07:36:40Z. Exit status: 0.

### Stage B command

```sh
git grep -n -iE "cycle cap|frozen head|verdict|re-freez|one more after|scale the plan implied" abdbe996282b33080406223ae8ba19e5a3634360 -- docs/ AGENTS.md README.md > /home/layup/.cache/layup-session/tmp-green/T-ebm4-review-mentions.txt
```

Start: 2026-10-21T07:36:40Z. End: 2026-10-21T07:36:40Z. Exit status: 0.

### Stage B command

```sh
python3 /home/layup/.cache/layup-session/tmp-green/T-ebm4-supplement-summary.py
```

Start: 2026-10-21T07:36:59Z. End: 2026-10-21T07:36:59Z. Exit status: 0.

### Stage B command

```sh
python3 /home/layup/.cache/layup-session/tmp-green/T-ebm4-inventory.py
```

Start: 2026-10-21T07:38:11Z. End: 2026-10-21T07:38:11Z. Exit status: 0.

```text
Wrote 1699 distinct baseline occurrences (913 primary, 789 supplementary, plus named context).
add-after: 1
keep-accepted-ADR: 55
keep-canonical: 30
keep-compatible: 57
keep-executable-or-description: 338
keep-fact: 14
keep-gap: 1
keep-historical: 982
keep-other-meaning: 6
keep-pitfall: 14
keep-product: 149
keep-product-context: 1
keep-product-plan: 12
keep-summary: 10
keep-validation: 1
normalize-quotes: 8
replace: 18
replace-example: 2
```

### Stage B command

```sh
python3 /home/layup/.cache/layup-session/tmp-green/T-ebm4-inventory.py
```

Start: 2026-10-21T07:41:31Z. End: 2026-10-21T07:41:31Z. Exit status: 0.

```text
Wrote 1699 distinct baseline occurrences (913 primary, 789 supplementary, plus named context).
add-after: 1
keep-accepted-ADR: 55
keep-canonical: 30
keep-compatible: 47
keep-executable-or-description: 348
keep-fact: 14
keep-gap: 1
keep-historical: 982
keep-other-meaning: 6
keep-pitfall: 14
keep-product: 149
keep-product-context: 1
keep-product-plan: 12
keep-summary: 10
keep-validation: 1
normalize-quotes: 8
replace: 18
replace-example: 2
```

### Stage B command

```sh
git diff -- docs/engineering-discipline.md docs/issue-workflow.md docs/guardrails.md docs/glossary.md
```

Start: 2026-10-21T07:41:31Z. End: 2026-10-21T07:41:31Z. Exit status: 0.

### Stage B command

```sh
sh docs/prd/prd-lint.sh
```

Start: 2026-10-21T07:42:02Z. End: 2026-10-21T07:42:02Z. Exit status: 0.

```text
prd-lint: OK
```

### Stage B command

```sh
sh docs/adr/adr-lint.sh
```

Start: 2026-10-21T07:42:02Z. End: 2026-10-21T07:42:02Z. Exit status: 0.

```text
adr-lint: OK
```

### Stage B command

```sh
sh docs/tests/run-discipline-tests.sh
```

Start: 2026-10-21T07:42:02Z. End: 2026-10-21T07:42:04Z. Exit status: 0.

```text

run-discipline-tests: 138 passed, 0 failed
```

### Stage B command

```sh
sh docs/links/link-lint.sh
```

Start: 2026-10-21T07:42:02Z. End: 2026-10-21T07:42:15Z. Exit status: 0.

```text
link-lint: OK  3047 links resolved
```

### Stage B command

```sh
sh docs/setup/setup-check.sh
```

Start: 2026-10-21T07:42:02Z. End: 2026-10-21T07:42:24Z. Exit status: 0.

```text
setup-check: pin OK
setup-check: kit-history OK
setup-check: facts OK
setup-check: onboarding OK
setup-check: glossary OK
setup-check: guardrails OK
setup-check: markers OK
setup-check: adapted OK
setup-check: ci OK
setup-check: protection OK
setup-check: identity OK
setup-check: procedure OK

run-discipline-tests: 138 passed, 0 failed
adr-lint: OK
prd-lint: OK
link-lint: OK  3047 links resolved
setup-check: kit-linters OK
```

### Stage B command

```sh
git diff --check abdbe996282b33080406223ae8ba19e5a3634360 HEAD
```

Start: 2026-10-21T07:42:33Z. End: 2026-10-21T07:42:33Z. Exit status: 0.

```text

```

### Stage B command

```sh
git diff --check
```

Start: 2026-10-21T07:42:33Z. End: 2026-10-21T07:42:33Z. Exit status: 0.

```text

```

### Stage B command

```sh
sh docs/tests/nested-checkout-check.sh
```

Start: 2026-10-21T07:42:33Z. End: 2026-10-21T07:42:34Z. Exit status: 0.

```text
ok    1 control: clean tree
ok    2 plain directory: its dead link is found
ok    3 dead link inside the nested checkout is not read
ok    4 dead link inside a linked worktree is not read
ok    5 no checkout: nested dead link still not read
ok    6 no checkout: plain directory still walked
ok    7 vendored and gitignored: link-lint
ok    8 symlink into the nested checkout is not read
ok    9 a quoted filename is still read
warning: adding embedded git repository: kit/.claude/worktrees/W
hint: You've added another git repository inside your current repository.
hint: Clones of the outer repository will not contain the contents of
hint: the embedded repository and will not know how to obtain it.
hint: If you meant to add a submodule, use:
hint:
hint: 	git submodule add <url> kit/.claude/worktrees/W
hint:
hint: If you added this path by mistake, you can remove it from the
hint: index with:
hint:
hint: 	git rm --cached kit/.claude/worktrees/W
hint:
hint: See "git help submodule" for more information.
hint: Disable this message with "git config set advice.addEmbeddedRepo false"
warning: adding embedded git repository: kit/nested
ok    10 vendored, a partial outer ignore: link-lint
nested-checkout-check: OK  10 cases behaved
```

### Stage B command

```sh
sh docs/tests/freeze-check/run-all.sh
```

Start: 2026-10-21T07:42:33Z. End: 2026-10-21T07:42:42Z. Exit status: 0.

```text
ok    bad-K1-blank-eof
ok    bad-K1-trailing-space
ok    bad-K2-missing-file
ok    bad-K2-missing-prd
ok    bad-K2-missing-traceability
ok    bad-K2-no-prd-frame-test
ok    bad-K2-no-prd-mapped-test
ok    bad-K2-no-s12-frame-test
ok    bad-K2-no-s12-mapped-test
ok    bad-K2-prd-other-file
ok    bad-K2-two-missing
ok    bad-K3-backtick-outside
ok    bad-K3-double-quoted
ok    bad-K3-no-match
ok    bad-K3-other-package
ok    bad-K3-other-tag
ok    bad-K3-span-across-lines
ok    bad-K3-subtest
ok    bad-K3-unclosed-span-eof
ok    bad-K3-unclosed-span-paragraph
ok    bad-K3-unquoted
ok    bad-K3-unsupported-form
ok    bad-K4-dropped-in-description
ok    bad-K4-form-a-different
ok    bad-K4-form-b-equal-hashes
ok    bad-K4-result-only
ok    bad-K4-two-hashes-no-restore
ok    bad-K5-over-files
ok    bad-K5-over-lines
ok    bad-K6-empty-table
ok    bad-K6-line-past-end
ok    bad-K6-line-zero
ok    bad-K6-missing-commit
ok    bad-K6-missing-file
ok    bad-K6-missing-path
ok    bad-K6-missing-test
ok    bad-K6-mixed-valid-invalid
ok    bad-K6-no-span
ok    bad-K6-range-zero
ok    bad-base-not-a-commit
ok    bad-crlf-conditions
ok    bad-env-K2-unreadable-prd
ok    bad-env-K2-unreadable
ok    bad-env-K3-unreadable-test
ok    bad-env-dirty-worktree
ok    bad-env-untracked-evidence
ok    bad-not-a-worktree
ok    bad-usage-bad-budget
ok    bad-usage-closeout
ok    bad-usage-evidence-directory-base
ok    bad-usage-evidence-directory-head
ok    bad-usage-evidence-missing
ok    bad-usage-evidence-no-value
ok    bad-usage-evidence-other-task
ok    bad-usage-evidence-suffix
ok    bad-usage-evidence-task-prefix
ok    bad-usage-evidence-trailing-slash
ok    bad-usage-evidence-traversal
ok    bad-usage-missing-budget
ok    bad-usage-no-args
ok    bad-usage-unknown-flag
ok    good-K1-clean
ok    good-K2-prd-by-file
ok    good-K2-registered
ok    good-K3-alternatives
ok    good-K3-span-across-lines
ok    good-K4-forms
ok    good-K5-at-max
ok    good-K5-ceiling
ok    good-K5-evidence-added
ok    good-K5-evidence-deleted
ok    good-K5-evidence-duplicate
ok    good-K5-evidence-inputs
ok    good-K5-evidence-kinds
ok    good-K5-evidence-literal
ok    good-K5-evidence-modified
ok    good-K5-evidence-renamed
ok    good-K5-file-limit
ok    good-K5-line-limit
ok    good-K5-no-closeout
ok    good-K5-percentage
ok    good-K5-truncation
ok    good-K6-no-conditions
ok    good-K6-valid
ok    good-crlf-all
freeze-check: 85 passed, 0 failed
```

### Stage B command

```sh
sh runs/T-ebm4/checks.sh
```

Start: 2026-10-21T07:44:17Z. End: 2026-10-21T07:44:17Z. Exit status: 0.

```text
ok: budget-unit
ok: evidence-exclusion
ok: executable-inputs
ok: input-definition
ok: audit-and-dual-use
ok: bookkeeping
ok: mixed-files
ok: advance-approval
ok: tolerance
ok: margin
ok: budget-post
ok: fixed-tolerance
ok: exact-approval
ok: overrun-block
ok: end-3-eligibility
ok: records-definition
ok: records-exclusions
ok: disputed-end
ok: end-3-post
ok: end-3-authority
ok: end-3-exception
ok: canonical-straight-quotes
ok: C1-new-comment-wait
ok: C1-retired-wait
ok: C1-one-shared-wait
ok: end-1-preserved
ok: end-2-preserved
ok: end-4-preserved
ok: gate-pointer
ok: gate-retired-next-round
ok: R12-budget-pointer
ok: R12-retired-slicing
ok: named-base-pointer
ok: named-base-retired-count
ok: C2-validation-preserved
ok: size-class-pointer
ok: size-class-retired-ceiling
ok: issue-split-pointer
ok: issue-split-example
ok: issue-split-retired-overrun
ok: issue-split-retired-example
ok: material-pointer
ok: material-retired-stop
ok: frozen-head-pointer
ok: frozen-head-retired-next-round
ok: C3-Records finding-row
ok: C3-records-bound
ok: C3-records-collision
ok: records-example
ok: C3-Tolerance-row
ok: tolerance-pointer
ok: tolerance-collision
ok: tolerance-example
ok: C3-Input-row
ok: input-pointer
ok: input-collision
ok: input-example
ok: C3-Margin-row
ok: margin-pointer
ok: margin-collision
ok: margin-example
ok: C3-Coordinated Universal Time-row
ok: C3-UTC-abbreviation
ok: C3-UTC-description
ok: C3-UTC-collision
ok: C3-UTC-example
ok: cap-example
ok: cap-retired-example
ok: ceiling-example
ok: ceiling-retired-example
ok: existing-checks-and-gates-preserved
```

### Stage B command

```sh
python3 /home/layup/.cache/layup-session/tmp-green/T-ebm4-final-audit.py
```

Start: 2026-10-21T07:44:17Z. End: 2026-10-21T07:44:18Z. Exit status: 0.

```text
Only the eight planned Stage B files differ; checks.sh is unchanged; close-out is deferred.
All four evidence records have straight quotes, no trailing white space and one final newline.
All 1,699 inventory locations exist at the baseline; both complete search outputs are covered.
Counted: docs/engineering-discipline.md: +19 -13
Counted: docs/glossary.md: +11 -6
Counted: docs/guardrails.md: +4 -3
Counted: docs/issue-workflow.md: +6 -6
Counted: runs/T-ebm4/checks.sh: +164 -0
Current working-tree branch projection: 232 lines over 5 files.
Maximum 400/10; tolerance 80/4; limit 480/14; projected margin 248/9.
Exactly the four declared audit evidence records leave this projection; the task check stays counted.
Committed starting head: cff8900583b5fe02d9b58c87c10dbd22ccd7619e
Comparison base: abdbe996282b33080406223ae8ba19e5a3634360
No clean-head pre-freeze or final landing audit is claimed by this working-tree projection.
```

## Stage C: merge of main and clean-head pre-freeze check

The coordinating session merged origin/main at
be6a0c7c199c4036c62933b45d38055e16e80681 (including PR #458) into T-ebm4.
The merge is 1c3ae2f1ac5019adc560e6c455064899727caf31, with parents
a18c047 and be6a0c7. This stage starts on that clean committed head; this
session made no commit or merge and changed no git metadata.

The plan-review condition names current main at landing as the budget base.
The pre-freeze run therefore uses origin/main, resolved below, rather than
silently retaining the original abdbe996282b33080406223ae8ba19e5a3634360 base.
The prior Stage B projection against that original base was 232/5; the
clean-head K5 figure against current main is also 232/5. The maximum remains
400/10, tolerance 80/4, limit 480/14 and margin 248/9. These are pre-freeze
figures; the final landing audit and issue post belong to later stages.

All four requested checks ran before this record was edited. The recorder
called `date -u +%FT%TZ` immediately before and after each command, retained
its actual exit status without a pipeline, and printed a heartbeat every ten
seconds during longer runs. TMPDIR and all scratch output were under
/home/layup/.cache/layup-session/tmp-green. The freeze-check command used the
requested redirection; its complete K1-K6 output is reproduced below and
remains in /home/layup/.cache/layup-session/tmp-green/ebm4-fc.txt.
The four evidence declarations are unchanged, K6 is enabled, and checks.sh
remains counted and unchanged. No task file or completed-log line is added
in this stage. Only this run record is prepared for the coordinator's commit.

### Stage C startup inspection ledger

These read-only inspection groups preceded the recorder. Their UTC start and
end were not captured and are not reported; no timestamps are inferred.
They were not acceptance runs. The statuses below are the command tool's
group statuses, not separately captured statuses for every command in a group.

Group 1. Start: not reported. End: not reported. Exit status: 0.

```sh
pwd && git status --short && git log -5 --oneline && cat AGENTS.md && cat docs/engineering-discipline.md && cat docs/issue-workflow.md
```

Group 2. Start: not reported. End: not reported. Exit status: 0.

```sh
cat runs/T-ebm4/test-runs.md && cat runs/T-ebm4/conditions.md && cat runs/T-ebm4/checks.sh
```

Group 3. Start: not reported. End: not reported. Exit status: 1.

```sh
sed -n '290,470p' docs/engineering-discipline.md && sed -n '1,210p' docs/tests/freeze-check.sh && cat docs/adr/0028-delegate-records-only-ends-and-budget-tolerance.md
```

The last cat failed because that ADR path does not exist. The correct ADR
path was discovered in group 4 and read in group 5. This inspection failure
was not a freeze-check failure and changed nothing.

Group 4. Start: not reported. End: not reported. Exit status: 0.

```sh
rg --files -g AGENTS.md -g '*0028*' && git rev-parse HEAD origin/main && git show -s --format=fuller 1c3ae2f && sed -n '1,180p' docs/links/link-lint.sh && sed -n '1,180p' docs/setup/setup-check.sh
```

Group 5. Start: not reported. End: not reported. Exit status: 0.

```sh
sed -n '290,455p' docs/engineering-discipline.md
cat docs/adr/0028-delegate-two-bounded-landing-decisions.md
cat runs/T-ebm4/conditions.md
sed -n '1,42p' runs/T-ebm4/checks.sh
sed -n '1,142p' docs/tests/freeze-check.sh
sed -n '535,660p' docs/tests/freeze-check.sh
tail -65 runs/T-ebm4/test-runs.md
```

### Stage C measured commands

```sh
git status --short --untracked-files=all
```

Start: 2026-10-21T07:47:10Z. End: 2026-10-21T07:47:10Z. Exit status: 0.

```text
```

```sh
git rev-parse HEAD origin/main
```

Start: 2026-10-21T07:47:10Z. End: 2026-10-21T07:47:10Z. Exit status: 0.

```text
1c3ae2f1ac5019adc560e6c455064899727caf31
be6a0c7c199c4036c62933b45d38055e16e80681
```

```sh
sh runs/T-ebm4/checks.sh
```

Start: 2026-10-21T07:47:10Z. End: 2026-10-21T07:47:10Z. Exit status: 0.

```text
ok: budget-unit
ok: evidence-exclusion
ok: executable-inputs
ok: input-definition
ok: audit-and-dual-use
ok: bookkeeping
ok: mixed-files
ok: advance-approval
ok: tolerance
ok: margin
ok: budget-post
ok: fixed-tolerance
ok: exact-approval
ok: overrun-block
ok: end-3-eligibility
ok: records-definition
ok: records-exclusions
ok: disputed-end
ok: end-3-post
ok: end-3-authority
ok: end-3-exception
ok: canonical-straight-quotes
ok: C1-new-comment-wait
ok: C1-retired-wait
ok: C1-one-shared-wait
ok: end-1-preserved
ok: end-2-preserved
ok: end-4-preserved
ok: gate-pointer
ok: gate-retired-next-round
ok: R12-budget-pointer
ok: R12-retired-slicing
ok: named-base-pointer
ok: named-base-retired-count
ok: C2-validation-preserved
ok: size-class-pointer
ok: size-class-retired-ceiling
ok: issue-split-pointer
ok: issue-split-example
ok: issue-split-retired-overrun
ok: issue-split-retired-example
ok: material-pointer
ok: material-retired-stop
ok: frozen-head-pointer
ok: frozen-head-retired-next-round
ok: C3-Records finding-row
ok: C3-records-bound
ok: C3-records-collision
ok: records-example
ok: C3-Tolerance-row
ok: tolerance-pointer
ok: tolerance-collision
ok: tolerance-example
ok: C3-Input-row
ok: input-pointer
ok: input-collision
ok: input-example
ok: C3-Margin-row
ok: margin-pointer
ok: margin-collision
ok: margin-example
ok: C3-Coordinated Universal Time-row
ok: C3-UTC-abbreviation
ok: C3-UTC-description
ok: C3-UTC-collision
ok: C3-UTC-example
ok: cap-example
ok: cap-retired-example
ok: ceiling-example
ok: ceiling-retired-example
ok: existing-checks-and-gates-preserved
```

```sh
sh docs/links/link-lint.sh
```

Start: 2026-10-21T07:47:10Z. End: 2026-10-21T07:47:24Z. Exit status: 0.

```text
link-lint: OK  3064 links resolved
```

```sh
sh docs/setup/setup-check.sh
```

Start: 2026-10-21T07:47:24Z. End: 2026-10-21T07:47:47Z. Exit status: 0.

```text
setup-check: pin OK
setup-check: kit-history OK
setup-check: facts OK
setup-check: onboarding OK
setup-check: glossary OK
setup-check: guardrails OK
setup-check: markers OK
setup-check: adapted OK
setup-check: ci OK
setup-check: protection OK
setup-check: identity OK
setup-check: procedure OK

run-discipline-tests: 138 passed, 0 failed
adr-lint: OK
prd-lint: OK
link-lint: OK  3064 links resolved
setup-check: kit-linters OK
```

```sh
sh docs/tests/freeze-check.sh --task T-ebm4 --base origin/main --budget 400/10 --evidence runs/T-ebm4/conditions.md --evidence runs/T-ebm4/semantic.md --evidence runs/T-ebm4/test-runs.md --evidence runs/T-ebm4/mentions.tsv > /home/layup/.cache/layup-session/tmp-green/ebm4-fc.txt 2>&1
```

Start: 2026-10-21T07:47:47Z. End: 2026-10-21T07:47:48Z. Exit status: 0.

```text
ok K1: white space
ok K2: registration of new tests
ok K3: test selectors
ok K4: restore evidence of mutants
ok K5: projection: figure 232/5, maximum 400/10, tolerance 80/4, limit 480/14, margin 248/9 (the figure at landing is still measured at landing)
ok K6: conditions table
```

All four requested checks exited 0. K6 audited the committed conditions table
and its cited paths. The successful freeze-check read merge head 1c3ae2f;
it does not claim to check the later evidence-only commit that this stage
prepares. Independent review and close-out remain for subsequent stages.

```sh
git diff --check
```

Start: 2026-10-21T07:49:00Z. End: 2026-10-21T07:49:00Z. Exit status: 0.

```text
```

```sh
git status --short --untracked-files=all
```

Start: 2026-10-21T07:49:00Z. End: 2026-10-21T07:49:00Z. Exit status: 0.

```text
 M runs/T-ebm4/test-runs.md
```

```sh
python3 /home/layup/.cache/layup-session/tmp-green/T-ebm4-stage-c-audit.py
```

Start: 2026-10-21T07:49:00Z. End: 2026-10-21T07:49:00Z. Exit status: 0.

```text
Only runs/T-ebm4/test-runs.md changed; checks.sh is unchanged.
The run record uses straight quotes, no trailing white space and one final newline.
No task file, completed-log line or scratch file was added to the worktree.
```
