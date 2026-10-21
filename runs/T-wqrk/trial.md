# T-wqrk: the trial, written down before any run

Written 2026-10-21T12:03:29Z (host clock, `date -u`), at the end of the phase 1 build and before any of the six model runs. The plan is comment 6314041425 of #462, its review is 6314145194 (eight conditions), and the Operator's decisions are 6315292816. This file is an evidence record. It states the plant, the judge, the run order and the rules for reading the result, so none of them can be chosen after a result is known.

## 1. What is tested, and what is not

A fix round's brief with a fix-round block (the "new" brief) against the same brief without it (the "old" brief). Only that treatment differs between the arms: the inputs, the head, the model and its effort are the same. The model's randomness still shows in the observations. The trial tests this planted case. It does not show that the block helps review in general (Operator decision Q3).

## 2. The plant (condition 1)

**I chose the second option of the brief: a rule whose living homes I could list exactly.** The task-ID scheme of the plan has more than two living homes (`docs/issue-workflow.md:30` and `:257`, `docs/tasks/backlog.md:16` and `:21`, `docs/engineering-discipline.md:790`, the alphabet constant in `internal/run/probe.go:164`, and the fixture `internal/verify/markers_test.go:53`), so it is not used.

**The rule:** the deadlines after which a reviewer harness is skipped. `docs/engineering-discipline.md:456-457` ("gives no output within five minutes, or gives no record within fifteen") is its operative home. The fix of the plant changes "fifteen" to "twenty" there. The other living home is the marked summary at `docs/adr/0026-keep-the-bootstrap-review-rules-as-the-standing-gate.md:39` ("the skip after five minutes with no output or fifteen with no record. Home: Who may review"). So the plant has exactly **two living homes**; the fix changes one and leaves the other.

**The commits** (a disposable bare repository built by `trial-T-wqrk.py`, with fixed author, committer and dates, so every run directory clones the same SHAs):

- C0 = the synthetic round 1 head, on `origin/main` of the date of the build: one docs commit in `docs/engineering-discipline.md` that adds "both deadlines run from the start of that harness's session" and leaves "fifteen".
- C1 = the fix: one commit that changes only "fifteen" to "twenty" in `docs/engineering-discipline.md`. The ADR line is in neither commit.

**Ground truth, pre-registered.** The home list is the result of this `git grep` at C0 and at C1 (the same command), taken before any run:

```
git grep -n -i -E 'no record within|with no record' -- . ':!runs' ':!docs/tasks/T-*'
```

and of the broader `git grep -n -i -E 'fifteen|(^|[^0-9])15 ?min|five minutes' -- . ':!runs' ':!docs/tasks/T-*'`. The matches are classified as follows (the code was searched: no Go file, script, hook or CI file states the deadlines; the `internal/` and `docs/spec/` matches say "no record" in other senses, and the `five minutes` matches are the forge token's):

| Match | Class |
| --- | --- |
| `docs/engineering-discipline.md:457` | the operative home; C0 says fifteen, C1 says twenty |
| `docs/adr/0026-keep-the-bootstrap-review-rules-as-the-standing-gate.md:39` | the second living home, a marked summary with a link to the first; says fifteen at C0 (agrees) and at C1 (stale) |
| `docs/adr/0012-build-layup-in-bootstrap-mode.md:46` | history: the ADR is "Superseded by ADR-0026"; says fifteen at both |
| `docs/tasks/T-*.md`, `runs/**` | history and passive records; many say "fifteen minutes" of past rounds |
| the other matches of the two greps | unrelated (a fixture dir name, the lease and token windows, a test comment, "no record" in other senses) |

The judge below accepts the second living home. A run that names only ADR-0012 (history) does not find the stale home.

## 3. The inputs, and the independence claim (condition 4)

The reviewer reads `.review-in/issue.md`, `plan.md`, `plan-review.md` and `fix.diff`, and the repository at C1. **No `round1.md` and no earlier verdict is handed over, and no account of what the fix meant.** The new preparer refuses any `round*.md` in `extra`, with no opt-in (round 1 finding M2), so the fix-round path itself carries none. The plan review's condition 1 (the deadline is twenty minutes) stands in for the finding: the round says whether it is settled. The texts are synthetic (the driver holds them).

Levels claimed for these runs: **Model** (GPT-6.1 Sol; the authors of the texts are Claude Opus 5.5 and Claude Sonnet 5.5), **Execution** (a separate run for each, with its own record) and **Method** (the brief names the lens and the fix). **Context holds in part**: the plan and the plan review are handed over, as in every real round of this harness; the author's reasoning on the fix and any earlier verdict are not. Line 1 of the template ("you have read no earlier review of this change") stays, because it is now always true.

## 4. The arms and the order

- **old:** `round-brief.tmpl` as it stands before this task (sha256 prefix of the file in the driver's output), with the house-style `planextra`: "reads the fix `C0..HEAD` (`.review-in/fix.diff`), says whether condition 1 is settled, reads the whole change again for anything new the fix brought".
- **new:** the same text plus the fix-round block and check 6 (`round-brief.tmpl.new`), filled by `prep-round.py.new`.
- The same `fix.diff` is in both arms, written by the preparer's `write_fix_diff`. The trial names no passive record, so the new block's pathspec is `.` and nothing is left out. In real rounds only named records (`runs/TASK/NAME.md`, a tracked plain file that no script reads) can be left out of the diff and the search; never by suffix (round 1 finding M4).
- **Reviewer:** GPT-6.1 Sol on Codex, `run-codex-review.sh`, effort `high` (Operator decision Q1), one run at a time with `nice -n 19`, no network.
- **Order:** old-1, new-1, old-2, new-2, old-3, new-3 (alternating, as the brief of this phase says; the plan review preferred three old runs first, and the alternation answers the same worry, a drift of the host or the model between the arms, by interleaving).

## 5. The judge (condition 2), frozen

**Revised 2026-10-21T12:34:07Z after round 1 of #462 (findings M1 and M5), before any run.** The first judge (sha256 `927893dfa9409c5649a4721965ed766af2290f83ba0ead4897d2f4e9c1fd94ef`) counted any disagreement word within 4 lines of a line that held `0026`. A note that denied an inconsistency ("I found no inconsistency in this summary") and a citation beside an unrelated stale item both returned FOUND. It is replaced by the judge below. **Its sha256 is `77e89964c9ce9bd160b49b26b9f01e05f92f0bc4f6ac70c8bfce4e518ebd48e1`** (file `judge-T-wqrk.sh`, text below); the runs use this file and no other.

**Validity (one check for both arms).** A record is valid when its first line is exactly `## Review record — round 2` (the heading the brief prescribes for round 2; no prefix test, no trailing text), it holds a `### Raw findings` line, and its `exit.txt` reads exactly `exit 0`. Otherwise the judge prints `JUDGE: INVALID <reason>` and exits 3: the run is neither a find nor a miss (section 6).

**The match.** Only the sections `### Raw findings` and `### Notes` are read. The unit is one finding or note: a list item at column 0 with its indented lines, or a paragraph. A unit is a find when it holds `0026`, fifteen or 15 and twenty or 20, and has one sentence that holds a word of disagreement, a deadline anchor (fifteen, twenty, 15, 20, deadline, no record) and no denial (no, nothing, none, never, without, n't, or "not stale" and its kind). A mere citation fails, a denial fails, an unrelated stale item in another sentence fails, and a mention in "What I checked and found sound" is not read. **Stated limit:** the unit is cut by layout, so an item whose parts sit at column 0 after blank lines is split and can give a miss; this is the same for both arms. A reviewer who diagnoses the planted disagreement in prose that the sentence test cannot read also gives a miss.

```sh
#!/bin/sh
# judge-T-wqrk.sh REVIEW.md [EXIT.TXT]: the pre-registered judge of the T-wqrk trial (runs/T-wqrk/trial.md), the same for both arms.
# Validity first: the first line of REVIEW.md is exactly "## Review record — round 2", it holds a "### Raw findings" line, and EXIT.TXT, if given,
# reads "exit 0". Otherwise: "JUDGE: INVALID <reason>", exit 3 (neither a find nor a miss). Exit 2 for no file.
# Cut: only the sections "### Raw findings" and "### Notes" are read. Unit: one finding or note (a list item at column 0 with its indented
# lines, or a paragraph). A unit is a FIND when it holds 0026, fifteen/15 and twenty/20, and one sentence that holds a disagreement word
# and fifteen/15/twenty/20/deadline/no record and no denial (no, nothing, none, never, without, n't, "not stale" and its kind).
# Prints FOUND lines and then "JUDGE: FOUND" (exit 0) or "JUDGE: MISS" (exit 1).
[ -f "${1:-}" ] || { echo "judge: no file" >&2; exit 2; }
inv() { echo "JUDGE: INVALID $1"; exit 3; }
[ "$(head -n 1 "$1")" = "## Review record — round 2" ] || inv "first line is not exactly '## Review record — round 2'"
grep -q '^### Raw findings' "$1" || inv "no '### Raw findings' line"
if [ -n "${2:-}" ]; then [ "$(cat "$2" 2>/dev/null)" = "exit 0" ] || inv "exit file is not 'exit 0'"; fi
awk '
function flush(   i, n, s, t, u, w) {
	if (blk == "") return
	t = tolower(blk); blk = ""
	if (t !~ /0026/ || t !~ /fifteen|(^|[^0-9])15([^0-9]|$)/ || t !~ /twenty|(^|[^0-9])20([^0-9]|$)/) return
	n = split(t, s, /[.!?;] +/)
	for (i = 1; i <= n; i++) {
		u = s[i]; w = u; gsub(/no longer|no record/, " ", w)
		if (u ~ diag && u ~ /fifteen|twenty|(^|[^0-9])(15|20)([^0-9]|$)|deadline|no record/ && w !~ /(^|[^a-z])(no|nothing|none|never|without)([^a-z]|$)|n.t / && w !~ neg) { print "FOUND in a unit at: " substr(u, 1, 150); f = 1; return }
	}
}
BEGIN {
	diag = "stale|still (says|states|gives|reads|has|names|uses|carries|keeps|shows|lists|mentions)|disagree|conflict|contradict|inconsisten|out of date|out-of-date|not (been )?updated|mismatch|differs from|no longer (match|agree)|does not (match|agree)|do not (match|agree)|old (value|number|figure|wording|deadline)|left (at|with|stating)"
	neg = "not (stale|inconsisten|in conflict|conflict|contradict|out of date|a mismatch)"
}
/^### / { flush(); p = ($0 ~ /^### (Raw findings|Notes)[ \t]*$/); prevblank = 1; next }
!p { next }
{
	if ($0 ~ /^[ \t]*$/) { prevblank = 1; blk = blk " "; next }
	if ($0 !~ /^[ \t]/ && (prevblank || $0 ~ /^([0-9]+[.)]|[-*]) /)) flush()
	blk = blk " " $0; prevblank = 0
}
END { flush(); print (f ? "JUDGE: FOUND" : "JUDGE: MISS"); exit (f ? 0 : 1) }' "$1"
```

**Checked before any run** (the cases are in `harness.md`): the registered positive (a finding that says ADR-0026:39 "still says fifteen" while the fix says twenty); the earlier negatives (a mere citation; the same diagnosis under "found sound"); the two negative controls of the round 1 review (a note that says "The deadlines mentioned are fifteen and twenty minutes. I found no inconsistency in this summary"; a citation beside an unrelated stale item); a record with the heading of round 1; and the exit-file check.

```
== pos
FOUND in a unit at: **a summary of the deadline is stale.** `docs/adr/0026-keep-the-bootstrap-review-rules-as-the-standing-gate.md:39` still says    "fifteen with no reco
JUDGE: FOUND
exit 0
== neg-cite
JUDGE: MISS
exit 1
== neg-sound
JUDGE: MISS
exit 1
== neg-denial
JUDGE: MISS
exit 1
== neg-unrelated
JUDGE: MISS
exit 1
== bad-heading
JUDGE: INVALID first line is not exactly '## Review record — round 2'
exit 3
== pos with exit file 'exit 124'
JUDGE: INVALID exit file is not 'exit 0'
exit 3
== pos with exit file 'exit 0'
FOUND in a unit at: **a summary of the deadline is stale.** `docs/adr/0026-keep-the-bootstrap-review-rules-as-the-standing-gate.md:39` still says    "fifteen with no reco
JUDGE: FOUND
exit 0
```

## 6. Reading the result (pre-registered)

- **Red** holds when at least 2 of the 3 valid old runs are misses. **Green** holds when all 3 valid new runs are finds.
- If 2 or more old runs are finds, the result reads **"not red"**. It stays "not red", it goes to the coordinator, who returns to the Operator before landing (decision Q2). No second plant is chosen after the fact.
- A **valid run** is one the judge does not print `INVALID` for (section 5): the exact first line, a `### Raw findings` line and `exit 0`. A timeout (exit 124), a missing record or a failed start makes the run **invalid**: it is neither a find nor a miss. Each invalid run is reported with its cause. It is repeated once in the same arm; a second invalid run in the same slot leaves the trial incomplete, and that is reported.
- Every run is reported, including runs that go against the result. No run is dropped, and none is repeated after its judge output is seen, except an invalid one by the rule above.
- The reviewers' records are copied to `runs/T-wqrk/trial-old-N.md` and `trial-new-N.md`, with the judge output, start and end times from the host clock, and the timings.

## 7. What must happen before the six runs (condition 7)

The driver, the judge, the fixture construction, both templates, the preparer diff and one prepared brief of each kind are preserved in `harness.md` as passive snapshots. They are reviewed before the runs start, because the runs are the costly action. While the runs go on, the coordinator shows which run is active and a heartbeat.

## 8. Plan baseline corrections (condition 8)

Check 3 of the old template already covers every document the change leaves stale, including untouched files; what the old brief lacks is a stated method of search. `review-record-lint.sh` checks nine reviewer fields (the author's `Fixes` is recorded apart) and their values and order, and it does not check what `Briefed on` says, so the gate text needs no lint change.
