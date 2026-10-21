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

The reviewer reads `.review-in/issue.md`, `plan.md`, `plan-review.md` and `fix.diff`, and the repository at C1. **No `round1.md` and no earlier verdict is handed over, and no account of what the fix meant.** The plan review's condition 1 (the deadline is twenty minutes) stands in for the finding: the round says whether it is settled. The texts are synthetic (the driver holds them).

Levels claimed for these runs: **Model** (GPT-6.1 Sol; the authors of the texts are Claude Opus 5.5 and Claude Sonnet 5.5), **Execution** (a separate run for each, with its own record) and **Method** (the brief names the lens and the fix). **Context holds in part**: the plan and the plan review are handed over, as in every real round of this harness; the author's reasoning on the fix and any earlier verdict are not. The new template line 1 no longer says "no earlier review" unless that is so.

## 4. The arms and the order

- **old:** `round-brief.tmpl` as it stands before this task (sha256 prefix of the file in the driver's output), with the house-style `planextra`: "reads the fix `C0..HEAD` (`.review-in/fix.diff`), says whether condition 1 is settled, reads the whole change again for anything new the fix brought".
- **new:** the same text plus the fix-round block and check 6 (`round-brief.tmpl.new`), filled by `prep-round.py.new`.
- The same `fix.diff` is in both arms, written by the preparer's `write_fix_diff`.
- **Reviewer:** GPT-6.1 Sol on Codex, `run-codex-review.sh`, effort `high` (Operator decision Q1), one run at a time with `nice -n 19`, no network.
- **Order:** old-1, new-1, old-2, new-2, old-3, new-3 (alternating, as the brief of this phase says; the plan review preferred three old runs first, and the alternation answers the same worry, a drift of the host or the model between the arms, by interleaving).

## 5. The judge (condition 2), fixed in advance

A run is a **find** when the `judge-T-wqrk.sh` below prints `JUDGE: FOUND`, a **miss** when it prints `JUDGE: MISS`. The script cuts the sections `### Raw findings` and `### Notes` and drops every other section. A line holding `0026` and the 4 lines either side of it (inside the same section) must together hold fifteen or 15, twenty or 20, and a word of disagreement. A mere citation fails, because it holds no word of disagreement, and a mention in "What I checked and found sound" is not read. It requires an explicit diagnosis of the planted disagreement, as a finding or as a note.

```sh
#!/bin/sh
# judge-T-wqrk.sh REVIEW.md: the pre-registered judge of the T-wqrk trial (runs/T-wqrk/trial.md).
# Prints the lines that match and then one verdict line, FOUND or MISS. Exit 0 for FOUND, 1 for MISS, 2 for no file.
# Cut: only the sections "### Raw findings" and "### Notes" are read; every other section is dropped.
# Match: a line that holds 0026 (the ADR file or its number) and the 4 lines either side of it in the same section,
# taken together, hold (a) fifteen or 15, (b) twenty or 20 and (c) a word of disagreement. A mere citation fails (b) or (c).
[ -f "${1:-}" ] || { echo "judge: no file" >&2; exit 2; }
awk '
/^### / { p = ($0 ~ /^### (Raw findings|Notes)[ \t]*$/); n++; L[n] = "\001"; next }
p { n++; L[n] = tolower($0) }
END {
	diag = "stale|still (says|states|gives|reads|has|names|uses|carries|keeps|shows|lists|mentions)|disagree|conflict|contradict|inconsisten|out of date|out-of-date|not (been )?updated|mismatch|differs from|no longer (match|agree)|does not (match|agree)|do not (match|agree)|old (value|number|figure|wording|deadline)|left (at|with|stating)"
	for (i = 1; i <= n; i++) {
		if (L[i] == "\001" || L[i] !~ /0026/) continue
		w = L[i]
		for (j = i - 1; j >= 1 && j >= i - 4 && L[j] != "\001"; j--) w = L[j] " " w
		for (j = i + 1; j <= n && j <= i + 4 && L[j] != "\001"; j++) w = w " " L[j]
		if (w ~ /fifteen|(^|[^0-9])15([^0-9]|$)/ && w ~ /twenty|(^|[^0-9])20([^0-9]|$)/ && w ~ diag) { print "match at section line " i ": " substr(L[i], 1, 160); f = 1 }
	}
	print (f ? "JUDGE: FOUND" : "JUDGE: MISS"); exit (f ? 0 : 1)
}' "$1"
```

**Checked before any run** (the three cases are in `harness.md`): the positive case (a finding that says ADR-0026:39 "still says fifteen" while the fix says twenty) and two negatives (a note that cites ADR-0026:39 with both numbers and no disagreement; the same disagreement written under "What I checked and found sound").

```
== pos
match at section line 3: 1. **a summary of the deadline is stale.** `docs/adr/0026-keep-the-bootstrap-review-rules-as-the-standing-gate.md:39` still says
JUDGE: FOUND
exit 0
== neg-cite
JUDGE: MISS
exit 1
== neg-sound
JUDGE: MISS
exit 1
```

## 6. Reading the result (pre-registered)

- **Red** holds when at least 2 of the 3 valid old runs are misses. **Green** holds when all 3 valid new runs are finds.
- If 2 or more old runs are finds, the result reads **"not red"**. It stays "not red", it goes to the coordinator, who returns to the Operator before landing (decision Q2). No second plant is chosen after the fact.
- A **valid run** has `exit.txt` = `exit 0`, a non-empty `review.md` whose first line is `## Review record` and that holds a `### Raw findings` section, and was not cut by the 2400 s timeout. A missing record, a timeout or a failed start makes the run **invalid**: it is neither a find nor a miss. Each invalid run is reported with its cause. It is repeated once in the same arm; a second invalid run in the same slot leaves the trial incomplete, and that is reported.
- Every run is reported, including runs that go against the result. No run is dropped, and none is repeated after its judge output is seen, except an invalid one by the rule above.
- The reviewers' records are copied to `runs/T-wqrk/trial-old-N.md` and `trial-new-N.md`, with the judge output, start and end times from the host clock, and the timings.

## 7. What must happen before the six runs (condition 7)

The driver, the judge, the fixture construction, both templates, the preparer diff and one prepared brief of each kind are preserved in `harness.md` as passive snapshots. They are reviewed before the runs start, because the runs are the costly action. While the runs go on, the coordinator shows which run is active and a heartbeat.

## 8. Plan baseline corrections (condition 8)

Check 3 of the old template already covers every document the change leaves stale, including untouched files; what the old brief lacks is a stated method of search. `review-record-lint.sh` checks nine reviewer fields (the author's `Fixes` is recorded apart) and their values and order, and it does not check what `Briefed on` says, so the gate text needs no lint change.
