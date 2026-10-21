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

**Revised three times before any run.** The first judge (sha256 `927893dfa9409c5649a4721965ed766af2290f83ba0ead4897d2f4e9c1fd94ef`) counted any disagreement word near a line holding `0026`; a denial and an unrelated stale item gave FOUND (round 1 of #462, M1 and M5). The second (sha256 `77e89964c9ce9bd160b49b26b9f01e05f92f0bc4f6ac70c8bfce4e518ebd48e1`, revised 2026-10-21 after round 1) tied the unit to one finding, but accepted a sentence about another deadline, and a hedged denial gave FOUND (round 2 of #462, F1). The third (sha256 `1f0c49e941f92294349fef6dd127d6588563970475eab2d8c6810c65eaddde59`, revised 12:58:26Z) bound the sentence to ADR-0026 but not the clause, and let hedged forms through (round 3 of #462, F1: two notes scored FOUND that diagnose nothing). **The fourth, below, is the frozen judge, repaired 2026-10-21 under the Operator's end 3 (comment 6316981250, term 1); its sha256 is `d9e1e291489ae7e8257a1b6b7814fbc94bc36572f625cab16c06b5ad21fc111e`** (file `judge-T-wqrk.sh`). It has had no review round (end 3). The runs use this file and no other.

**Validity (one check for both arms).** A record is valid when its first line is exactly `## Review record — round 2` (the heading the brief prescribes for round 2; no prefix test, no trailing text), it holds a `### Raw findings` line, and its `exit.txt` reads exactly `exit 0`. Otherwise the judge prints `JUDGE: INVALID <reason>` and exits 3: the run is neither a find nor a miss (section 6).

**The match.** Only the sections `### Raw findings` and `### Notes` are read. The unit is one finding or note: a list item at column 0 with its indented lines, or a paragraph. A unit is a find when it holds `0026`, fifteen or 15 and twenty or 20, and **one clause of one sentence holds all of** (a clause is cut at while, whilst, whereas, but, although, though, however, separately, meanwhile): `0026`, fifteen or 15 or "no record", a word of disagreement, and no other document (README, a token, a lease, another ADR); and **the sentence holds** no denial (no, nothing, none, never, without, n't, or "not stale" and its kind) and no hedge (incorrect, wrong, would be, could be, may, might, could, possible, possibly, maybe, perhaps, unclear, needs checking, should be checked, to check, verify, if, whether, unless, a question mark, is fine, is correct, is consistent, agrees with). The affirmative clause must name ADR-0026 itself with its record deadline and the disagreement word; a neighbouring sentence or clause, another item, another document or another deadline does not count. A mere citation fails, a denial fails, a hedged denial fails, an unrelated stale item fails, and a mention in "What I checked and found sound" is not read. **Stated limits:** the unit is cut by layout, so an item whose parts sit at column 0 after blank lines is split; a diagnosis whose ADR-0026 mention sits in another sentence or clause than the disagreement word gives a miss, and a hedge anywhere in the sentence rejects it. The judge can under-count and cannot be trusted to confirm: see the hand-confirmation rule in section 6. Both limits are the same for both arms.

```sh
#!/bin/sh
# judge-T-wqrk.sh REVIEW.md [EXIT.TXT]: the pre-registered judge of the T-wqrk trial (runs/T-wqrk/trial.md), the same for both arms.
# Validity first: the first line of REVIEW.md is exactly "## Review record — round 2", it holds a "### Raw findings" line, and EXIT.TXT, if given,
# reads "exit 0". Otherwise: "JUDGE: INVALID <reason>", exit 3 (neither a find nor a miss). Exit 2 for no file.
# Cut: only the sections "### Raw findings" and "### Notes" are read. Unit: one finding or note (a list item at column 0 with its indented
# lines, or a paragraph). A unit is a FIND when it holds 0026, fifteen/15 and twenty/20, and ONE sentence of it holds all of: 0026,
# fifteen/15/no record, a word of disagreement, and no denial (no, nothing, none, never, without, n't, "not stale" and its kind) and no hedge
# (incorrect, wrong, would be, could be, may/might/could, possible/possibly, needs checking, should be checked, verify whether, if, whether,
# a question mark, is fine, is correct, is consistent, agrees with). Binding: ONE clause of that sentence (cut at while, but, whereas, although,
# however, separately...) holds 0026, the deadline token AND the disagreement word, and names no other document (README, a token, a lease, another
# ADR); a diagnosis in a neighbouring sentence or clause, or of another item or document, does not count. Third repair, after round 3 of #462 (F1).
# Prints FOUND lines and then "JUDGE: FOUND" (exit 0) or "JUDGE: MISS" (exit 1).
[ -f "${1:-}" ] || { echo "judge: no file" >&2; exit 2; }
inv() { echo "JUDGE: INVALID $1"; exit 3; }
[ "$(head -n 1 "$1")" = "## Review record — round 2" ] || inv "first line is not exactly '## Review record — round 2'"
grep -q '^### Raw findings' "$1" || inv "no '### Raw findings' line"
if [ -n "${2:-}" ]; then [ "$(cat "$2" 2>/dev/null)" = "exit 0" ] || inv "exit file is not 'exit 0'"; fi
awk '
function bound(u,   k, m, c, ok) {
	m = split(u, c, sep)
	for (k = 1; k <= m; k++)
		if (c[k] ~ /0026/ && c[k] ~ /fifteen|(^|[^0-9])15([^0-9]|$)|no record/ && c[k] ~ diag && c[k] !~ other) return 1
	return 0
}
function flush(   i, n, s, t, u, w) {
	if (blk == "") return
	t = tolower(blk); blk = ""; gsub(/\?/, " qmark?", t)
	if (t !~ /0026/ || t !~ /fifteen|(^|[^0-9])15([^0-9]|$)/ || t !~ /twenty|(^|[^0-9])20([^0-9]|$)/) return
	n = split(t, s, /[.!?;] +/)
	for (i = 1; i <= n; i++) {
		u = s[i]; w = u; gsub(/no longer|no record/, " ", w)
		if (bound(u) && w !~ /(^|[^a-z])(no|nothing|none|never|without)([^a-z]|$)|n.t / && w !~ neg && w !~ hedge) { print "FOUND in a unit at: " substr(u, 1, 150); f = 1; return }
	}
}
BEGIN {
	diag = "stale|still (says|states|gives|reads|has|names|uses|carries|keeps|shows|lists|mentions)|disagree|conflict|contradict|inconsisten|out of date|out-of-date|not (been )?updated|mismatch|differs from|no longer (match|agree)|does not (match|agree)|do not (match|agree)|old (value|number|figure|wording|deadline)|left (at|with|stating)"
	hedge = "incorrect|wrong|would be|could be|possib|qmark|needs? (to be )?(check|verif|confirm)|should be (check|verif|confirm)|to (check|verify|confirm)|(^|[^a-z])(verify|maybe|perhaps|may|might|could|if|whether|unless|unclear)([^a-z]|$)|(is|are) (fine|correct|consistent)|(^|[^a-z])(agrees|agree with)"
	sep = "(^|[^a-z])(while|whilst|whereas|but|although|though|however|separately|meanwhile)([^a-z]|$)"
	other = "readme|token|lease|adr-00(0[0-9]|1[0-9]|2[0-57-9])"
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

**Checked before any run** (the cases are in `harness.md`): the registered positive (a finding that says ADR-0026:39 "still says fifteen" while the fix says twenty); the earlier negatives (`neg-cite`, `neg-sound`, `neg-denial`, `neg-unrelated`, `neg-token`, `neg-hedged`); the four controls registered for round 3 F1 (`neg-token-same` and `neg-possible`, the two notes of the round verbatim, `neg-question`, `neg-may`); a record with the heading of round 1; and the exit-file check. **Red first:** the third judge returned FOUND, exit 0, on all four new controls (13:15:18Z, recorded in `test-runs.md` section 8). **Green** (the fourth judge, run 2026-10-21T2026-10-21T13:16:21Z to 2026-10-21T13:16:21Z, `date -u`): the positive FOUND with exit 0, every negative MISS with exit 1, the invalid cases exit 3.

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
== neg-token
JUDGE: MISS
exit 1
== neg-hedged
JUDGE: MISS
exit 1
== neg-token-same
JUDGE: MISS
exit 1
== neg-possible
JUDGE: MISS
exit 1
== neg-question
JUDGE: MISS
exit 1
== neg-may
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
- **Hand confirmation (the Operator's end 3, comment 6316981250, term 2).** A judge FOUND is necessary, not sufficient. A run counts as a find only when the coordinator's hand reading of that run's record confirms an explicit diagnosis of the planted disagreement in ADR-0026's record deadline (the fifteen-minute deadline it still gives, against the twenty of `docs/engineering-discipline.md`). The reading is recorded per run, with the quoted sentence, before the arms are compared. A judge MISS stays a miss. The judge can then only under-count, never invent a find. Red and green below are read on the confirmed finds.
- If 2 or more old runs are finds, the result reads **"not red"**. It stays "not red", it goes to the coordinator, who returns to the Operator before landing (decision Q2). No second plant is chosen after the fact.
- A **valid run** is one the judge does not print `INVALID` for (section 5): the exact first line, a `### Raw findings` line and `exit 0`. A timeout (exit 124), a missing record or a failed start makes the run **invalid**: it is neither a find nor a miss. Each invalid run is reported with its cause. It is repeated once in the same arm; a second invalid run in the same slot leaves the trial incomplete, and that is reported.
- Every run is reported, including runs that go against the result. No run is dropped, and none is repeated after its judge output is seen, except an invalid one by the rule above.
- The reviewers' records are copied to `runs/T-wqrk/trial-old-N.md` and `trial-new-N.md`, with the judge output, start and end times from the host clock, and the timings.

## 7. What must happen before the six runs (condition 7)

The driver, the judge, the fixture construction, both templates, the preparer diff and one prepared brief of each kind are preserved in `harness.md` as passive snapshots. They are reviewed before the runs start, because the runs are the costly action. While the runs go on, the coordinator shows which run is active and a heartbeat.

## 8. Plan baseline corrections (condition 8)

Check 3 of the old template already covers every document the change leaves stale, including untouched files; what the old brief lacks is a stated method of search. `review-record-lint.sh` checks nine reviewer fields (the author's `Fixes` is recorded apart) and their values and order, and it does not check what `Briefed on` says, so the gate text needs no lint change.

## 9. The results (the trial ran 2026-10-21T13:19:27Z to 13:40:23Z; the Operator's decision b, comment 6317693581)

Six valid runs (GPT-6.1 Sol, effort high, Codex, one at a time, alternating), all `exit 0`. The driver's log and the judge's raw output are in `test-runs.md` section 9. The judge is the frozen one (sha256 `d9e1e291489ae7e8257a1b6b7814fbc94bc36572f625cab16c06b5ad21fc111e`).

| Run | Start to end (UTC) | Elapsed | Tokens | Judge | Hand reading (the coordinator) | Counts as a find |
| --- | --- | --- | --- | --- | --- | --- |
| old-1 | 13:19:27Z to 13:21:23Z | 1 min 56 s | 66,503 | MISS | diagnosed, confirmed | no (judge MISS, term 2) |
| new-1 | 13:21:23Z to 13:26:16Z | 4 min 53 s | 128,623 | MISS | diagnosed, confirmed | no (judge MISS, term 2) |
| old-2 | 13:26:16Z to 13:28:38Z | 2 min 22 s | 80,483 | MISS | diagnosed, confirmed | no (judge MISS, term 2) |
| new-2 | 13:28:38Z to 13:31:39Z | 3 min 1 s | 111,817 | MISS | diagnosed, confirmed | no (judge MISS, term 2) |
| old-3 | 13:31:39Z to 13:34:43Z | 3 min 4 s | 88,764 | FOUND | yes, confirmed | yes |
| new-3 | 13:34:43Z to 13:40:23Z | 5 min 40 s | 130,901 | MISS | diagnosed, confirmed | no (judge MISS, term 2) |

Tokens are the `tokens used` line of each run's `err.txt`.

**The hand reading,** per run, with the quoted sentence (term 2 of comment 6316981250). Every run diagnoses ADR-0026:39 still giving fifteen minutes against the fix's twenty:

- **old-1**, a note (line 21): "ADR-0026 still records “fifteen with no record” at `docs/adr/0026-keep-the-bootstrap-review-rules-as-the-standing-gate.md:39`. It names the canonical home that this branch updates." Its disposition is a close-out record, so it is a note, not material.
- **new-1**, a note (finding 1): "ADR-0026 retains an unmarked restatement of the old deadline. Its accepted Decision §5 says fifteen minutes at `docs/adr/0026-keep-the-bootstrap-review-rules-as-the-standing-gate.md:39` … This disagrees with the fix, but does not establish a separate operative deadline."
- **old-2**, a note (line 23): "`docs/adr/0026-keep-the-bootstrap-review-rules-as-the-standing-gate.md:39` retains fifteen minutes. It points to the canonical home changed here." It is not treated as a second operative deadline.
- **new-2**, material (finding 1): "The fix gives twenty minutes at `docs/engineering-discipline.md:457`. ADR-0026’s accepted Decision, item 5, still specifies “fifteen with no record” at `docs/adr/0026-keep-the-bootstrap-review-rules-as-the-standing-gate.md:39`." Finding 1 calls it "an operative contradiction".
- **old-3**, material (finding 1): "Accepted ADR-0026 still states “fifteen with no record” in its Decision (`docs/adr/0026-keep-the-bootstrap-review-rules-as-the-standing-gate.md:39`)." The heading reads "the deadline conflicts with accepted ADR-0026".
- **new-3**, material (F1): "`docs/adr/0026-keep-the-bootstrap-review-rules-as-the-standing-gate.md:39` still states “the skip after five minutes with no output or fifteen with no record.” … This is an operative contradiction on the issue’s stated path."

**The registered result (section 6).** Old: 1 of 3 finds (old-3), so **red holds**. New: 0 of 3, so **green fails**. The finds are the judge's FOUND results that the hand reading confirmed. The hand reading found 6 of 6 against 6 of 6 (old: note, note, material; new: note, material, material); it is reported, and it does **not** count (term 2: a judge MISS stays a miss). The judge under-counted five real diagnoses.

**Conclusion.** No gain is shown. The plant was too easy: the old brief finds it as well, in all three runs. The new brief's runs cost more tokens (111,817 to 130,901 against 66,503 to 88,764) and more time (3 min 1 s to 5 min 40 s against 1 min 56 s to 3 min 4 s). The trial tests one planted case and says nothing on review in general (decision Q3).

**The Operator's decision b (comment 6317693581).** The fix-round preparer and brief are delivered (installed in the shared harness at 2026-10-21T13:18:59Z, backups `*.bak-20261021T131859Z`). The gate-text change (the Method row of "Who may review" and the `Briefed on` clause) is **dropped**: `docs/engineering-discipline.md` is back at the text of `origin/main`.
