# T-wqrk: the harness, as snapshots for review before the runs

Written 2026-10-21T12:05:25Z (`date -u`). Passive copies of the local harness files, none of them in the repository, preserved so that the producing code of the trial is reviewed before the six paid runs (plan review condition 7). The live files are `/home/layup/.cache/layup-session/m2e-scratchpad/rounds`; the `.new` files are not yet in place (the shared-harness protocol: the coordinator gives notice, keeps a `.bak`, and moves them). The sections: the driver (which builds the plant and the run directories), the judge and its three cases, both templates and the preparer with its diff, one prepared brief of each kind, and the sweep of the old and new wording.

### A. trial-T-wqrk.py (the driver and the fixture construction)

````python
#!/usr/bin/env python3
# trial-T-wqrk.py prepare [--worktree DIR] [--prep PATH] [--old TMPL] [--new TMPL]
# Builds the T-wqrk trial (runs/T-wqrk/trial.md): the plant (a disposable bare repository with a synthetic round-1 head C0 and its
# fix C1), the six run directories (old, new, old, new, old, new) with .review-in/ and brief.md, and the run order. It launches nothing.
# The brief is filled by render_brief of the preparer (--prep), the same replacement and '@' assertion as prep-round.py.
# It writes only under /home/layup/.cache/layup-session/reviews/T-wqrk-trial-*; no forge, no network.
import argparse, hashlib, importlib.machinery, importlib.util, json, os, shutil, subprocess, sys
S = "/home/layup/.cache/layup-session/m2e-scratchpad/rounds"
RV = "/home/layup/.cache/layup-session/reviews/T-wqrk-trial-"
TASK, ISSUE = "T-9qz4", 9999
ORDER = ["old-1", "new-1", "old-2", "new-2", "old-3", "new-3"]
DATE = ["2026-10-21T12:00:00Z", "2026-10-21T12:10:00Z"]
def under(p):
    p = os.path.abspath(p); assert p.startswith(RV), "write outside " + RV + ": " + p; return p
def sh(*a, env=None, cwd=None): return subprocess.run(a, check=True, capture_output=True, text=True, env=env, cwd=cwd).stdout
ISSUE_MD = """## Goal

Task `%s`: a reviewer harness that is slow to answer is skipped too early. Raise the deadline for a record from fifteen to twenty minutes, and say from when the clocks run.

## Scope

- The skip deadlines of a reviewer harness, as the project's documents state them.

## Acceptance criteria

- [ ] A reviewed plan.
- [ ] The text says that the deadlines run from the start of the harness's session.
- [ ] The text gives twenty minutes for a record.
""" % TASK
PLAN_MD = """## Plan

1. Say in the rule of reviewer selection that both deadlines run from the start of the harness's session.
2. Raise the deadline for a record from fifteen to twenty minutes.
3. Run the lints that read text (`sh docs/adr/adr-lint.sh`, `sh docs/prd/prd-lint.sh`, `sh docs/links/link-lint.sh`) and `git diff --check`.

Size: about 4 lines in one document. Cycle cap 1.
"""
REVIEW_MD = """## Plan review

| Field | Value |
| ----- | ----- |
| Verdict | approve-with-conditions |
| Budget maximum | 10 lines added plus removed over 2 files |
| Cycle cap | 1 |

### Conditions

1. The deadline for a record is twenty minutes, not fifteen: the slowest round of the last week took 17 minutes.
2. The clock starts when the harness session starts, and the text says so.
"""
PLANEXTRA = ("`.review-in/fix.diff` is the fix `C0..HEAD`: this round reads it, says whether condition 1 of the plan review is settled, and reads the whole change again "
    "for anything new the fix brought. This is the last round the cap allows: the verdict is `nothing material in scope` or `not mergeable, findings recorded`, never `material`.")
def load_prep(path):
    ld = importlib.machinery.SourceFileLoader("prep_round", path); sp = importlib.util.spec_from_loader("prep_round", ld)
    m = importlib.util.module_from_spec(sp); ld.exec_module(m); return m
def build_plant(wt, main):
    P = under(RV + "plant"); shutil.rmtree(P, ignore_errors=True); os.makedirs(P)
    sh("git", "init", "-q", "--bare", "plant.git", cwd=P)
    sh("git", "fetch", "-q", wt, main + ":refs/heads/main", cwd=P + "/plant.git")
    w = P + "/work"; sh("git", "clone", "-q", "--no-tags", "-b", "main", P + "/plant.git", w)
    f = w + "/docs/engineering-discipline.md"; t = open(f).read()
    a = "the resource record; no decision is asked for it."
    assert t.count(a) == 1 and t.count("within fifteen") == 1
    for i, (old, new, msg) in enumerate([(a, "the resource record; both deadlines run from the start of that harness's session; no decision is asked for it.", "docs: %s the deadlines run from the start of the session" % TASK),
                                         ("within fifteen", "within twenty", "docs: %s twenty minutes for a record" % TASK)]):
        t = open(f).read().replace(old, new); open(f, "w").write(t)
        env = dict(os.environ, GIT_AUTHOR_NAME="Trial Author", GIT_AUTHOR_EMAIL="trial@example.invalid", GIT_COMMITTER_NAME="Trial Author", GIT_COMMITTER_EMAIL="trial@example.invalid",
                   GIT_AUTHOR_DATE=DATE[i], GIT_COMMITTER_DATE=DATE[i])
        sh("git", "commit", "-q", "-am", msg, env=env, cwd=w)
        if i == 0: c0 = sh("git", "rev-parse", "HEAD", cwd=w).strip()
    c1 = sh("git", "rev-parse", "HEAD", cwd=w).strip()
    sh("git", "push", "-q", "origin", "HEAD:refs/heads/trial", cwd=w); shutil.rmtree(w)
    return P, c0, c1
def main():
    ap = argparse.ArgumentParser(); ap.add_argument("cmd", choices=["prepare"])
    ap.add_argument("--worktree", default="/home/layup/projects/layup/.worktree/T-wqrk")
    ap.add_argument("--prep", default=S + "/prep-round.py.new"); ap.add_argument("--old", default=S + "/round-brief.tmpl"); ap.add_argument("--new", default=S + "/round-brief.tmpl.new")
    a = ap.parse_args(); prep = load_prep(a.prep)
    main_sha = sh("git", "rev-parse", "refs/remotes/origin/main", cwd=a.worktree).strip()
    P, c0, c1 = build_plant(a.worktree, main_sha)
    os.makedirs(P + "/templates"); tm = {}
    for k, src in (("old", a.old), ("new", a.new)):
        shutil.copy(src, P + "/templates/" + k + ".tmpl"); tm[k] = open(src).read()
        print(k, "template", hashlib.sha256(tm[k].encode()).hexdigest()[:16], src)
    print("main", main_sha, "C0", c0, "C1", c1)
    for name in ORDER:
        arm, n = name.split("-"); R = under(RV + name); shutil.rmtree(R, ignore_errors=True); os.makedirs(R); L = R + "/layup"
        sh("git", "clone", "-q", "--no-tags", P + "/plant.git", L)
        sh("git", "checkout", "-q", "--detach", c1, cwd=L); sh("git", "remote", "set-url", "origin", "https://github.com/pharzam/layup.git", cwd=L)
        assert sh("git", "rev-parse", "origin/main", cwd=L).strip() == main_sha and sh("git", "rev-parse", "HEAD", cwd=L).strip() == c1
        RI = L + "/.review-in"; os.makedirs(RI)
        for d in (".gocache", ".tmp"): os.makedirs(L + "/" + d)
        open(L + "/.git/info/exclude", "a").write(".review-in\n.gocache\n.tmp\nopencode.json\n")
        for fn, body in (("issue.md", ISSUE_MD), ("plan.md", PLAN_MD), ("plan-review.md", REVIEW_MD)): open(RI + "/" + fn, "w").write(body)
        c = dict(task=TASK, issue=ISSUE, round=2, cycle=1, cap=1, what="the skip deadlines of a reviewer harness", reviewer="GPT-6.1 Sol (`gpt-6.1-sol`, effort high)",
                 harness="Codex CLI (codex exec)", exc=" Harness exception: ADR-0026's routing rule names `claude`, `devin` and `opencode`; Codex is used on the Operator's message of 2026-10-13.",
                 authors="Claude Opus 5.5 (the plan) and Claude Sonnet 5.5 (the documents)", planextra=PLANEXTRA, head=c1, main=main_sha, budget="10 lines added plus removed over 2 files", fix_base=c0)
        prep.write_fix_diff(L, RI, c0, [])
        b = prep.render_brief(tm[arm], c)  # the old template has no fix-round placeholder, so the same call renders it as before
        open(RI + "/brief.md", "w").write(b)
        assert not sh("git", "status", "--short", cwd=L)
        print(name, sorted(os.listdir(RI)), "brief", hashlib.sha256(b.encode()).hexdigest()[:16])
    print("RUN ORDER (one at a time, nice 19; NOT launched by this script):")
    for name in ORDER: print("  nice -n 19 sh /home/layup/.cache/layup-session/harness/2026-10-16-m2e/run-codex-review.sh %s high" % (RV + name))
main()
````

### B. judge-T-wqrk.sh

````sh
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
````

### B. judge case pos.md

````markdown
## Review record — round 2

| Field | Value |
| --- | --- |
| Verdict | material |

### Raw findings

1. **A summary of the deadline is stale.** `docs/adr/0026-keep-the-bootstrap-review-rules-as-the-standing-gate.md:39` still says
   "fifteen with no record", but the fix sets the deadline to twenty minutes in `docs/engineering-discipline.md:457`.
   Falsifier: the ADR line reads twenty. Fix: say twenty there, or replace the summary with a link.

### Notes

None.

### What I checked and found sound

The fix is one word.
````

### B. judge case neg-cite.md

````markdown
## Review record — round 2

### Raw findings

None.

### Notes

1. I read `docs/adr/0026-keep-the-bootstrap-review-rules-as-the-standing-gate.md:39`. It summarises reviewer selection:
   five minutes with no output, fifteen with no record. The fix changes the row in `docs/engineering-discipline.md:457` to twenty.
   The summary links to its home.

### What I checked and found sound

Condition 1 is met.
````

### B. judge case neg-sound.md

````markdown
## Review record — round 2

### Raw findings

None.

### Notes

None.

### What I checked and found sound

`docs/adr/0026-keep-the-bootstrap-review-rules-as-the-standing-gate.md:39` still says fifteen, while the fix says twenty, and that is stale.
````

### C. round-brief.tmpl (the old brief)

````
You are the reviewer of round @ROUND@ of task `@TASK@`, issue #@ISSUE@ of `pharzam/layup` (@WHAT@), in a fresh session. You are @REVIEWER@, a reasoning-tier model. The authors of the change are @AUTHORS@; you are none of them, and you have read no earlier review of this change.

## What you may read

- `.review-in/issue.md`: the issue on the forge.
- `.review-in/plan.md` and `.review-in/plan-review.md`: the plan and the plan review that approved it with conditions. The plan review's conditions bind the change. @PLANEXTRA@
- The repository in this directory at the frozen head `@HEAD@` (the branch), with `origin/main` at `@MAIN@`. Start with `AGENTS.md`. The change is `git diff @MAIN@..HEAD` (the run records under `runs/@TASK@/` are the author's claims, not proof).

You have no forge access and need none. Do not read outside this directory. You work in a sandbox with no network; the Go cache is `.gocache` and the temporary files are `.tmp`, which are not part of the change. Write nothing else to the tree. The sandbox has no local socket, so `httptest` integration tests may fail there with the base's own code: the integration level is the author's run and CI's; say which you ran.

## What to check

1. **Correctness** against the issue's acceptance criteria, the plan review's conditions and the specs (`docs/spec/loop.md`, `records.md`, `packages.md`): each behaviour the change claims is true in the code, and each refusal named in the plan has a test that fails if the refusal is removed.
2. **Test first and test quality** (R8, `docs/guardrails.md` §2): the named tests assert semantics, not only the absence of an error; the evidence of red runs is credible against the code.
3. **R10 and gate step 7:** every document the change leaves stale is fixed in the branch (spec, plan row, packages, traceability, PRD cells).
4. **Scope and budget:** nothing outside the issue's group entered the branch. Measure `git diff --shortstat @MAIN@..HEAD` against the plan review's `Budget maximum` (@BUDGET@, close-out inside; the close-out is not written yet).
5. **Material has a test** (`docs/engineering-discipline.md`): a finding is material only by that test; a note causes no round. A defect the change revealed off its path is a note naming a new issue, not a finding.

## The form of your answer

Your final message is the record, in Markdown, and nothing else. Its first line is `## Review record — round @ROUND@`; write no sentence before it. Then a table `| Field | Value |` with these fields, one row each: `Commit reviewed` (the full SHA `@HEAD@`), `Reviewer`, `Lens`, `Briefed on`, `Barred from`, `Independence claimed`, `Cycle` (`@CYCLE@`), `Verdict` (`material` if any finding is material and the cap allows another round, else `nothing material in scope`, or `not mergeable, findings recorded` at the cap). Then the sections `### Raw findings` (each with its operative clause, the evidence as `path:line`, a falsifier and a bounded fix), `### Notes`, `### What I checked and found sound`, `### Commands and results` (each command with its exit status). The cap is @CAP@. End with the line `*Written by @REVIEWER@ on @HARNESS@, the reviewer of round @ROUND@ of @TASK@ (#@ISSUE@).@EXC@*`

Write plain English with short sentences, in the style of the repository's documents.
````

### C. round-brief.tmpl.new (the new brief)

````
You are the reviewer of round @ROUND@ of task `@TASK@`, issue #@ISSUE@ of `pharzam/layup` (@WHAT@), in a fresh session. You are @REVIEWER@, a reasoning-tier model. The authors of the change are @AUTHORS@; you are none of them, and @EARLIER@

## What you may read

- `.review-in/issue.md`: the issue on the forge.
- `.review-in/plan.md` and `.review-in/plan-review.md`: the plan and the plan review that approved it with conditions. The plan review's conditions bind the change. @PLANEXTRA@
- The repository in this directory at the frozen head `@HEAD@` (the branch), with `origin/main` at `@MAIN@`. Start with `AGENTS.md`. The change is `git diff @MAIN@..HEAD` (the run records under `runs/@TASK@/` are the author's claims, not proof).
@FIXREAD@
You have no forge access and need none. Do not read outside this directory. You work in a sandbox with no network; the Go cache is `.gocache` and the temporary files are `.tmp`, which are not part of the change. Write nothing else to the tree. The sandbox has no local socket, so `httptest` integration tests may fail there with the base's own code: the integration level is the author's run and CI's; say which you ran.

## What to check

1. **Correctness** against the issue's acceptance criteria, the plan review's conditions and the specs (`docs/spec/loop.md`, `records.md`, `packages.md`): each behaviour the change claims is true in the code, and each refusal named in the plan has a test that fails if the refusal is removed.
2. **Test first and test quality** (R8, `docs/guardrails.md` §2): the named tests assert semantics, not only the absence of an error; the evidence of red runs is credible against the code.
3. **R10 and gate step 7:** every document the change leaves stale is fixed in the branch (spec, plan row, packages, traceability, PRD cells).
4. **Scope and budget:** nothing outside the issue's group entered the branch. Measure `git diff --shortstat @MAIN@..HEAD` against the plan review's `Budget maximum` (@BUDGET@, close-out inside; the close-out is not written yet).
5. **Material has a test** (`docs/engineering-discipline.md`): a finding is material only by that test; a note causes no round. A defect the change revealed off its path is a note naming a new issue, not a finding.
@FIXCHECK@
## The form of your answer

Your final message is the record, in Markdown, and nothing else. Its first line is `## Review record — round @ROUND@`; write no sentence before it. Then a table `| Field | Value |` with these fields, one row each: `Commit reviewed` (the full SHA `@HEAD@`), `Reviewer`, `Lens`, `Briefed on`, `Barred from`, `Independence claimed`, `Cycle` (`@CYCLE@`), `Verdict` (`material` if any finding is material and the cap allows another round, else `nothing material in scope`, or `not mergeable, findings recorded` at the cap). Then the sections `### Raw findings` (each with its operative clause, the evidence as `path:line`, a falsifier and a bounded fix), `### Notes`, `### What I checked and found sound`, `### Commands and results` (each command with its exit status). The cap is @CAP@. End with the line `*Written by @REVIEWER@ on @HARNESS@, the reviewer of round @ROUND@ of @TASK@ (#@ISSUE@).@EXC@*`

Write plain English with short sentences, in the style of the repository's documents.
````

### D. prep-round.py.new

````python
# prep-round.py CONFIG.json: clone at the frozen head, inputs and brief of a review round. Writes RUNDIR/layup/.review-in/.
# T-wqrk: a round after a fix (round 2 or later) needs "fix_base" (a full SHA, an ancestor of "head"). The script then writes
# .review-in/fix.diff and fills the fix-round block and check 6 of the template. "fix_exclude" lists passive run records
# (runs/TASK/NAME.md only) left out of that diff; they are named in the block. A "round*.md" in "extra" is an earlier verdict:
# it needs "earlier_record": true, and the brief then says that the reviewer was handed it. PREP_TEMPLATE overrides the template path.
import json, os, re, shutil, subprocess, sys
S = "/home/layup/.cache/layup-session/m2e-scratchpad"
KEYS = ("ROUND", "TASK", "ISSUE", "WHAT", "REVIEWER", "AUTHORS", "PLANEXTRA", "HEAD", "MAIN", "BUDGET", "CYCLE", "CAP", "HARNESS", "EXC")
NO_EARLIER = "you have read no earlier review of this change."
FIXREAD = ("- This round follows a fix. The fix is `git diff FIXBASE HEAD` (saved as `.review-in/fix.diff`EXCLUDED); read it first. "
    "For each changed line that states a rule (a sentence, a number, a name, a code constant), choose its key terms, in the old wording and in the new, "
    "with the variants a reference takes (a number spelled out and in digits, an abbreviation, a section number, a link target, a constant). "
    "Run `git grep -n -i -F -e TERM -- . ':!runs/*.md'` at this head for each term. The pathspec leaves out only the passive run records; "
    "code, scripts, indexes, task files and decision records stay in. Every file that states the same rule is a home of it. "
    "Read each home and say whether it agrees with the fix. Write the terms, the commands and the homes you read in `Briefed on` and in `### Commands and results`.\n")
FIXCHECK = ("6. **Homes of a rule the fix changed:** for each home that disagrees with the fix, say whether it is an operative statement, a marked summary, a fixture or history. "
    "Judge it by the material test and by One home per rule: a conflicting restatement is a note unless it also meets the material test.\n")
def sh(*a, **k): return subprocess.run(a, check=True, capture_output=True, text=True, **k).stdout
def fail(m): sys.exit("prep-round: " + m)
def earlier_text(c):
    names = [n for n in c.get("extra", {}) if re.fullmatch(r"round\d+\.md", n)]
    if not names: return NO_EARLIER
    if not c.get("earlier_record"): fail("extra " + ", ".join(names) + " holds an earlier verdict; set earlier_record to true to hand it over")
    return "you were handed an earlier round's record (" + ", ".join("`.review-in/" + n + "`" for n in names) + "), so your Context independence holds in part; your record must say so."
def fix_parts(c):
    """(fix_base, fix_exclude, FIXREAD text, FIXCHECK text) for a round after a fix; empty for round 1."""
    base = c.get("fix_base")
    if not base:
        if int(c["round"]) > 1: fail("round " + str(c["round"]) + " follows a fix and needs fix_base, the previous frozen head")
        return None, [], "", ""
    if not re.fullmatch(r"[0-9a-f]{40}", base): fail("fix_base must be a full SHA")
    ex = c.get("fix_exclude", [])
    for p in ex:
        if not re.fullmatch("runs/" + re.escape(c["task"]) + r"/[A-Za-z0-9._-]+\.md", p): fail("fix_exclude takes only runs/" + c["task"] + "/NAME.md, not " + p)
    note = "; left out as passive run records: " + ", ".join("`" + p + "`" for p in ex) if ex else ""
    return base, ex, FIXREAD.replace("FIXBASE", base).replace("EXCLUDED", note), FIXCHECK
def write_fix_diff(L, RI, base, ex):
    if sh("git", "merge-base", base, "HEAD", cwd=L).strip() != base or sh("git", "rev-parse", "HEAD", cwd=L).strip() == base: fail("fix_base is not an earlier commit of the head")
    open(RI + "/fix.diff", "w").write(sh("git", "diff", base, "HEAD", "--", ".", *[":!" + p for p in ex], cwd=L))
def render_brief(tmpl, c):
    b = tmpl.replace("@EARLIER@", earlier_text(c)).replace("@FIXREAD@", fix_parts(c)[2]).replace("@FIXCHECK@", fix_parts(c)[3])
    for k in KEYS: b = b.replace("@" + k + "@", str(c[k.lower()]))
    assert "@" not in b.replace("@anthropic", ""), [l for l in b.splitlines() if "@" in l]
    return b
def main():
    c = json.load(open(sys.argv[1]))
    base, ex, _, _ = fix_parts(c); earlier_text(c)
    R = c["rundir"]; os.makedirs(R, exist_ok=True); L = R + "/layup"
    if not os.path.isdir(L): sh("git", "clone", "-q", "--no-tags", "--single-branch", "-b", "main", "https://github.com/pharzam/layup.git", "layup", cwd=R)
    sh("git", "fetch", "-q", "origin", f"+refs/heads/{c['task']}:refs/remotes/origin/{c['task']}", "+refs/heads/main:refs/remotes/origin/main", cwd=L)
    sh("git", "checkout", "-q", "--detach", c["head"], cwd=L)
    assert sh("git", "rev-parse", "HEAD", cwd=L).strip() == c["head"]
    assert sh("git", "rev-parse", "origin/main", cwd=L).strip() == c["main"], "main moved"
    RI = L + "/.review-in"; os.makedirs(RI, exist_ok=True)
    for d in (".gocache", ".tmp"): os.makedirs(L + "/" + d, exist_ok=True)
    with open(L + "/.git/info/exclude", "a") as f: f.write(".review-in\n.gocache\n.tmp\nopencode.json\n")
    def body(n): return sh("gh", "api", f"repos/pharzam/layup/issues/comments/{n}", "--jq", ".body")
    open(RI + "/issue.md", "w").write(sh("gh", "api", f"repos/pharzam/layup/issues/{c['issue']}", "--jq", ".body"))
    open(RI + "/plan.md", "w").write(body(c["plan"])); open(RI + "/plan-review.md", "w").write(body(c["plan_review"]))
    for name, n in c.get("extra", {}).items(): open(RI + "/" + name, "w").write(body(n))
    if base: write_fix_diff(L, RI, base, ex)
    open(RI + "/brief.md", "w").write(render_brief(open(os.environ.get("PREP_TEMPLATE", S + "/rounds/round-brief.tmpl")).read(), c))
    if c.get("opencode"): shutil.copy(S + "/opencode-review.json", L + "/opencode.json")
    print(c["task"], "round", c["round"], sorted(os.listdir(RI)), sh("git", "status", "--short", cwd=L) or "clean")
if __name__ == "__main__": main()
````

### D. prep-round.py: diff of the old file to the new

````diff
--- /home/layup/.cache/layup-session/m2e-scratchpad/rounds/prep-round.py	2026-10-20 13:07:46.473367370 +0300
+++ /home/layup/.cache/layup-session/m2e-scratchpad/rounds/prep-round.py.new	2026-10-21 15:00:52.029827344 +0300
@@ -1,25 +1,65 @@
 # prep-round.py CONFIG.json: clone at the frozen head, inputs and brief of a review round. Writes RUNDIR/layup/.review-in/.
-import json, os, shutil, subprocess, sys
+# T-wqrk: a round after a fix (round 2 or later) needs "fix_base" (a full SHA, an ancestor of "head"). The script then writes
+# .review-in/fix.diff and fills the fix-round block and check 6 of the template. "fix_exclude" lists passive run records
+# (runs/TASK/NAME.md only) left out of that diff; they are named in the block. A "round*.md" in "extra" is an earlier verdict:
+# it needs "earlier_record": true, and the brief then says that the reviewer was handed it. PREP_TEMPLATE overrides the template path.
+import json, os, re, shutil, subprocess, sys
 S = "/home/layup/.cache/layup-session/m2e-scratchpad"
-c = json.load(open(sys.argv[1]))
-R = c["rundir"]; os.makedirs(R, exist_ok=True); L = R + "/layup"
+KEYS = ("ROUND", "TASK", "ISSUE", "WHAT", "REVIEWER", "AUTHORS", "PLANEXTRA", "HEAD", "MAIN", "BUDGET", "CYCLE", "CAP", "HARNESS", "EXC")
+NO_EARLIER = "you have read no earlier review of this change."
+FIXREAD = ("- This round follows a fix. The fix is `git diff FIXBASE HEAD` (saved as `.review-in/fix.diff`EXCLUDED); read it first. "
+    "For each changed line that states a rule (a sentence, a number, a name, a code constant), choose its key terms, in the old wording and in the new, "
+    "with the variants a reference takes (a number spelled out and in digits, an abbreviation, a section number, a link target, a constant). "
+    "Run `git grep -n -i -F -e TERM -- . ':!runs/*.md'` at this head for each term. The pathspec leaves out only the passive run records; "
+    "code, scripts, indexes, task files and decision records stay in. Every file that states the same rule is a home of it. "
+    "Read each home and say whether it agrees with the fix. Write the terms, the commands and the homes you read in `Briefed on` and in `### Commands and results`.\n")
+FIXCHECK = ("6. **Homes of a rule the fix changed:** for each home that disagrees with the fix, say whether it is an operative statement, a marked summary, a fixture or history. "
+    "Judge it by the material test and by One home per rule: a conflicting restatement is a note unless it also meets the material test.\n")
 def sh(*a, **k): return subprocess.run(a, check=True, capture_output=True, text=True, **k).stdout
-if not os.path.isdir(L): sh("git", "clone", "-q", "--no-tags", "--single-branch", "-b", "main", "https://github.com/pharzam/layup.git", "layup", cwd=R)
-sh("git", "fetch", "-q", "origin", f"+refs/heads/{c['task']}:refs/remotes/origin/{c['task']}", "+refs/heads/main:refs/remotes/origin/main", cwd=L)
-sh("git", "checkout", "-q", "--detach", c["head"], cwd=L)
-assert sh("git", "rev-parse", "HEAD", cwd=L).strip() == c["head"]
-assert sh("git", "rev-parse", "origin/main", cwd=L).strip() == c["main"], "main moved"
-RI = L + "/.review-in"; os.makedirs(RI, exist_ok=True)
-for d in (".gocache", ".tmp"): os.makedirs(L + "/" + d, exist_ok=True)
-with open(L + "/.git/info/exclude", "a") as f: f.write(".review-in\n.gocache\n.tmp\nopencode.json\n")
-def body(n): return sh("gh", "api", f"repos/pharzam/layup/issues/comments/{n}", "--jq", ".body")
-open(RI + "/issue.md", "w").write(sh("gh", "api", f"repos/pharzam/layup/issues/{c['issue']}", "--jq", ".body"))
-open(RI + "/plan.md", "w").write(body(c["plan"])); open(RI + "/plan-review.md", "w").write(body(c["plan_review"]))
-for name, n in c.get("extra", {}).items(): open(RI + "/" + name, "w").write(body(n))
-b = open(S + "/rounds/round-brief.tmpl").read()
-for k in ("ROUND", "TASK", "ISSUE", "WHAT", "REVIEWER", "AUTHORS", "PLANEXTRA", "HEAD", "MAIN", "BUDGET", "CYCLE", "CAP", "HARNESS", "EXC"):
-    b = b.replace("@" + k + "@", str(c[k.lower()]))
-assert "@" not in b.replace("@anthropic", ""), [l for l in b.splitlines() if "@" in l]
-open(RI + "/brief.md", "w").write(b)
-if c.get("opencode"): shutil.copy(S + "/opencode-review.json", L + "/opencode.json")
-print(c["task"], "round", c["round"], sorted(os.listdir(RI)), sh("git", "status", "--short", cwd=L) or "clean")
+def fail(m): sys.exit("prep-round: " + m)
+def earlier_text(c):
+    names = [n for n in c.get("extra", {}) if re.fullmatch(r"round\d+\.md", n)]
+    if not names: return NO_EARLIER
+    if not c.get("earlier_record"): fail("extra " + ", ".join(names) + " holds an earlier verdict; set earlier_record to true to hand it over")
+    return "you were handed an earlier round's record (" + ", ".join("`.review-in/" + n + "`" for n in names) + "), so your Context independence holds in part; your record must say so."
+def fix_parts(c):
+    """(fix_base, fix_exclude, FIXREAD text, FIXCHECK text) for a round after a fix; empty for round 1."""
+    base = c.get("fix_base")
+    if not base:
+        if int(c["round"]) > 1: fail("round " + str(c["round"]) + " follows a fix and needs fix_base, the previous frozen head")
+        return None, [], "", ""
+    if not re.fullmatch(r"[0-9a-f]{40}", base): fail("fix_base must be a full SHA")
+    ex = c.get("fix_exclude", [])
+    for p in ex:
+        if not re.fullmatch("runs/" + re.escape(c["task"]) + r"/[A-Za-z0-9._-]+\.md", p): fail("fix_exclude takes only runs/" + c["task"] + "/NAME.md, not " + p)
+    note = "; left out as passive run records: " + ", ".join("`" + p + "`" for p in ex) if ex else ""
+    return base, ex, FIXREAD.replace("FIXBASE", base).replace("EXCLUDED", note), FIXCHECK
+def write_fix_diff(L, RI, base, ex):
+    if sh("git", "merge-base", base, "HEAD", cwd=L).strip() != base or sh("git", "rev-parse", "HEAD", cwd=L).strip() == base: fail("fix_base is not an earlier commit of the head")
+    open(RI + "/fix.diff", "w").write(sh("git", "diff", base, "HEAD", "--", ".", *[":!" + p for p in ex], cwd=L))
+def render_brief(tmpl, c):
+    b = tmpl.replace("@EARLIER@", earlier_text(c)).replace("@FIXREAD@", fix_parts(c)[2]).replace("@FIXCHECK@", fix_parts(c)[3])
+    for k in KEYS: b = b.replace("@" + k + "@", str(c[k.lower()]))
+    assert "@" not in b.replace("@anthropic", ""), [l for l in b.splitlines() if "@" in l]
+    return b
+def main():
+    c = json.load(open(sys.argv[1]))
+    base, ex, _, _ = fix_parts(c); earlier_text(c)
+    R = c["rundir"]; os.makedirs(R, exist_ok=True); L = R + "/layup"
+    if not os.path.isdir(L): sh("git", "clone", "-q", "--no-tags", "--single-branch", "-b", "main", "https://github.com/pharzam/layup.git", "layup", cwd=R)
+    sh("git", "fetch", "-q", "origin", f"+refs/heads/{c['task']}:refs/remotes/origin/{c['task']}", "+refs/heads/main:refs/remotes/origin/main", cwd=L)
+    sh("git", "checkout", "-q", "--detach", c["head"], cwd=L)
+    assert sh("git", "rev-parse", "HEAD", cwd=L).strip() == c["head"]
+    assert sh("git", "rev-parse", "origin/main", cwd=L).strip() == c["main"], "main moved"
+    RI = L + "/.review-in"; os.makedirs(RI, exist_ok=True)
+    for d in (".gocache", ".tmp"): os.makedirs(L + "/" + d, exist_ok=True)
+    with open(L + "/.git/info/exclude", "a") as f: f.write(".review-in\n.gocache\n.tmp\nopencode.json\n")
+    def body(n): return sh("gh", "api", f"repos/pharzam/layup/issues/comments/{n}", "--jq", ".body")
+    open(RI + "/issue.md", "w").write(sh("gh", "api", f"repos/pharzam/layup/issues/{c['issue']}", "--jq", ".body"))
+    open(RI + "/plan.md", "w").write(body(c["plan"])); open(RI + "/plan-review.md", "w").write(body(c["plan_review"]))
+    for name, n in c.get("extra", {}).items(): open(RI + "/" + name, "w").write(body(n))
+    if base: write_fix_diff(L, RI, base, ex)
+    open(RI + "/brief.md", "w").write(render_brief(open(os.environ.get("PREP_TEMPLATE", S + "/rounds/round-brief.tmpl")).read(), c))
+    if c.get("opencode"): shutil.copy(S + "/opencode-review.json", L + "/opencode.json")
+    print(c["task"], "round", c["round"], sorted(os.listdir(RI)), sh("git", "status", "--short", cwd=L) or "clean")
+if __name__ == "__main__": main()
````

### E. a prepared old brief (old-1)

````
You are the reviewer of round 2 of task `T-9qz4`, issue #9999 of `pharzam/layup` (the skip deadlines of a reviewer harness), in a fresh session. You are GPT-6.1 Sol (`gpt-6.1-sol`, effort high), a reasoning-tier model. The authors of the change are Claude Opus 5.5 (the plan) and Claude Sonnet 5.5 (the documents); you are none of them, and you have read no earlier review of this change.

## What you may read

- `.review-in/issue.md`: the issue on the forge.
- `.review-in/plan.md` and `.review-in/plan-review.md`: the plan and the plan review that approved it with conditions. The plan review's conditions bind the change. `.review-in/fix.diff` is the fix `C0..HEAD`: this round reads it, says whether condition 1 of the plan review is settled, and reads the whole change again for anything new the fix brought. This is the last round the cap allows: the verdict is `nothing material in scope` or `not mergeable, findings recorded`, never `material`.
- The repository in this directory at the frozen head `a5906ee579bdfe2c7ec92da7fa18f41c0ea8cc3d` (the branch), with `origin/main` at `01d7da3082333efadb227afda14fd98715edf5c4`. Start with `AGENTS.md`. The change is `git diff 01d7da3082333efadb227afda14fd98715edf5c4..HEAD` (the run records under `runs/T-9qz4/` are the author's claims, not proof).

You have no forge access and need none. Do not read outside this directory. You work in a sandbox with no network; the Go cache is `.gocache` and the temporary files are `.tmp`, which are not part of the change. Write nothing else to the tree. The sandbox has no local socket, so `httptest` integration tests may fail there with the base's own code: the integration level is the author's run and CI's; say which you ran.

## What to check

1. **Correctness** against the issue's acceptance criteria, the plan review's conditions and the specs (`docs/spec/loop.md`, `records.md`, `packages.md`): each behaviour the change claims is true in the code, and each refusal named in the plan has a test that fails if the refusal is removed.
2. **Test first and test quality** (R8, `docs/guardrails.md` §2): the named tests assert semantics, not only the absence of an error; the evidence of red runs is credible against the code.
3. **R10 and gate step 7:** every document the change leaves stale is fixed in the branch (spec, plan row, packages, traceability, PRD cells).
4. **Scope and budget:** nothing outside the issue's group entered the branch. Measure `git diff --shortstat 01d7da3082333efadb227afda14fd98715edf5c4..HEAD` against the plan review's `Budget maximum` (10 lines added plus removed over 2 files, close-out inside; the close-out is not written yet).
5. **Material has a test** (`docs/engineering-discipline.md`): a finding is material only by that test; a note causes no round. A defect the change revealed off its path is a note naming a new issue, not a finding.

## The form of your answer

Your final message is the record, in Markdown, and nothing else. Its first line is `## Review record — round 2`; write no sentence before it. Then a table `| Field | Value |` with these fields, one row each: `Commit reviewed` (the full SHA `a5906ee579bdfe2c7ec92da7fa18f41c0ea8cc3d`), `Reviewer`, `Lens`, `Briefed on`, `Barred from`, `Independence claimed`, `Cycle` (`1`), `Verdict` (`material` if any finding is material and the cap allows another round, else `nothing material in scope`, or `not mergeable, findings recorded` at the cap). Then the sections `### Raw findings` (each with its operative clause, the evidence as `path:line`, a falsifier and a bounded fix), `### Notes`, `### What I checked and found sound`, `### Commands and results` (each command with its exit status). The cap is 1. End with the line `*Written by GPT-6.1 Sol (`gpt-6.1-sol`, effort high) on Codex CLI (codex exec), the reviewer of round 2 of T-9qz4 (#9999). Harness exception: ADR-0026's routing rule names `claude`, `devin` and `opencode`; Codex is used on the Operator's message of 2026-10-13.*`

Write plain English with short sentences, in the style of the repository's documents.
````

### E. a prepared new brief (new-1)

````
You are the reviewer of round 2 of task `T-9qz4`, issue #9999 of `pharzam/layup` (the skip deadlines of a reviewer harness), in a fresh session. You are GPT-6.1 Sol (`gpt-6.1-sol`, effort high), a reasoning-tier model. The authors of the change are Claude Opus 5.5 (the plan) and Claude Sonnet 5.5 (the documents); you are none of them, and you have read no earlier review of this change.

## What you may read

- `.review-in/issue.md`: the issue on the forge.
- `.review-in/plan.md` and `.review-in/plan-review.md`: the plan and the plan review that approved it with conditions. The plan review's conditions bind the change. `.review-in/fix.diff` is the fix `C0..HEAD`: this round reads it, says whether condition 1 of the plan review is settled, and reads the whole change again for anything new the fix brought. This is the last round the cap allows: the verdict is `nothing material in scope` or `not mergeable, findings recorded`, never `material`.
- The repository in this directory at the frozen head `a5906ee579bdfe2c7ec92da7fa18f41c0ea8cc3d` (the branch), with `origin/main` at `01d7da3082333efadb227afda14fd98715edf5c4`. Start with `AGENTS.md`. The change is `git diff 01d7da3082333efadb227afda14fd98715edf5c4..HEAD` (the run records under `runs/T-9qz4/` are the author's claims, not proof).
- This round follows a fix. The fix is `git diff a0fed3c442031e86a70e6b460c6c0df88a1fe93f HEAD` (saved as `.review-in/fix.diff`); read it first. For each changed line that states a rule (a sentence, a number, a name, a code constant), choose its key terms, in the old wording and in the new, with the variants a reference takes (a number spelled out and in digits, an abbreviation, a section number, a link target, a constant). Run `git grep -n -i -F -e TERM -- . ':!runs/*.md'` at this head for each term. The pathspec leaves out only the passive run records; code, scripts, indexes, task files and decision records stay in. Every file that states the same rule is a home of it. Read each home and say whether it agrees with the fix. Write the terms, the commands and the homes you read in `Briefed on` and in `### Commands and results`.

You have no forge access and need none. Do not read outside this directory. You work in a sandbox with no network; the Go cache is `.gocache` and the temporary files are `.tmp`, which are not part of the change. Write nothing else to the tree. The sandbox has no local socket, so `httptest` integration tests may fail there with the base's own code: the integration level is the author's run and CI's; say which you ran.

## What to check

1. **Correctness** against the issue's acceptance criteria, the plan review's conditions and the specs (`docs/spec/loop.md`, `records.md`, `packages.md`): each behaviour the change claims is true in the code, and each refusal named in the plan has a test that fails if the refusal is removed.
2. **Test first and test quality** (R8, `docs/guardrails.md` §2): the named tests assert semantics, not only the absence of an error; the evidence of red runs is credible against the code.
3. **R10 and gate step 7:** every document the change leaves stale is fixed in the branch (spec, plan row, packages, traceability, PRD cells).
4. **Scope and budget:** nothing outside the issue's group entered the branch. Measure `git diff --shortstat 01d7da3082333efadb227afda14fd98715edf5c4..HEAD` against the plan review's `Budget maximum` (10 lines added plus removed over 2 files, close-out inside; the close-out is not written yet).
5. **Material has a test** (`docs/engineering-discipline.md`): a finding is material only by that test; a note causes no round. A defect the change revealed off its path is a note naming a new issue, not a finding.
6. **Homes of a rule the fix changed:** for each home that disagrees with the fix, say whether it is an operative statement, a marked summary, a fixture or history. Judge it by the material test and by One home per rule: a conflicting restatement is a note unless it also meets the material test.

## The form of your answer

Your final message is the record, in Markdown, and nothing else. Its first line is `## Review record — round 2`; write no sentence before it. Then a table `| Field | Value |` with these fields, one row each: `Commit reviewed` (the full SHA `a5906ee579bdfe2c7ec92da7fa18f41c0ea8cc3d`), `Reviewer`, `Lens`, `Briefed on`, `Barred from`, `Independence claimed`, `Cycle` (`1`), `Verdict` (`material` if any finding is material and the cap allows another round, else `nothing material in scope`, or `not mergeable, findings recorded` at the cap). Then the sections `### Raw findings` (each with its operative clause, the evidence as `path:line`, a falsifier and a bounded fix), `### Notes`, `### What I checked and found sound`, `### Commands and results` (each command with its exit status). The cap is 1. End with the line `*Written by GPT-6.1 Sol (`gpt-6.1-sol`, effort high) on Codex CLI (codex exec), the reviewer of round 2 of T-9qz4 (#9999). Harness exception: ADR-0026's routing rule names `claude`, `devin` and `opencode`; Codex is used on the Operator's message of 2026-10-13.*`

Write plain English with short sentences, in the style of the repository's documents.
````

### E. the fix.diff of both

````diff
diff --git a/docs/engineering-discipline.md b/docs/engineering-discipline.md
index 2a8e306c..9b889062 100644
--- a/docs/engineering-discipline.md
+++ b/docs/engineering-discipline.md
@@ -454,7 +454,7 @@ model the reviewer must differ from, and the session that freezes the head names
 the reviewer order and answers the notes. The reviewer is any model that differs
 from every author's, on `claude`, `devin` or `opencode`, tried one at a time in
 that order, until one returns a record. A harness that fails, gives no output
-within five minutes, or gives no record within fifteen, is skipped and named in
+within five minutes, or gives no record within twenty, is skipped and named in
 the resource record; both deadlines run from the start of that harness's session; no decision is asked for it. The [Model tiers](#model-tiers)
 name the models and the models not to use.

````

## F. The sweep of the old and new wording (plan step 4, condition 5)

Pre-registered search, run at the head of the branch and over the local harness. The pathspec leaves out only the passive run records (`runs/*.md`) and the per-task history files, so scripts and other operative files under `runs/` stay visible. The terms are the old wording ("no earlier review", "the fix it reads", "round after a fix"), the new wording (`fix.diff`, `FIXREAD`, `fix_base`) and the surrounding references (`prep-round`, `round-brief`, `planextra`).

```
git grep -n -i -E 'no earlier review|earlier review|fix\.diff|the fix it reads|reads the fix|round after a fix|after a fix|prep-round|round-brief|planextra' -- . ':!runs/*.md' ':!docs/tasks/T-*'
grep -rIl -E 'no earlier review|fix\.diff|round-brief\.tmpl|prep-round' /home/layup/.cache/layup-session/m2e-scratchpad /home/layup/.cache/layup-session/harness
```

| Match | Class | Action |
| --- | --- | --- |
| `docs/engineering-discipline.md:439` (Method row), `:499` (`Briefed on`) | operative statement | changed in this task (step 6), in one home |
| `docs/engineering-discipline.md:242`, `:318`, `README.md:37`, `docs/onboarding-for-engineers.md:179`, `docs/adr/0012-...:45`, `docs/adr/0026-...:11` and `:37`, `docs/glossary.md:171` | say a round runs after a fix, or that the next round "reads the fix" | not the rule of what a fix round reads; marked summaries and history; unchanged (`:318` is a reference, as the plan says) |
| `docs/guardrails.md:723` | lesson on reading an earlier review | agrees with the new refusal of `round*.md` extras without `earlier_record`; unchanged |
| `rounds/round-brief.tmpl`, `rounds/prep-round.py` (live) | operative code | replaced by the `.new` files by the coordinator, per the shared-harness protocol |
| `harness/prep-round.sh`, `harness/2026-10-13-codex/prep-round-codex.py`, `m2e-scratchpad/T-8rm6/prep-r4.sh` | older preparers from 2026-10-12 and 2026-10-13, with their own templates (`@POSITION@`, `@FIXDIFF@`); not on the live path of the host (`prep-round.py` with `rounds/*.json`) | unchanged; named here as legacy |
| `m2e-scratchpad/rounds/r*-N.json` | past configs; a round-2 config has no `fix_base`, so the new preparer refuses it until one is added | by design (condition 3); the coordinator adds `fix_base` to a new config |
| copies of `docs/tasks/T-8bwd.md` and `runs/T-4wrw/review-round-1.md` in scratch trees | history | unchanged |

(The blank context line of the diff above lost its one space of padding so that `git diff --check` passes; no other byte of a snapshot differs from its file.)
