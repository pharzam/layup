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
