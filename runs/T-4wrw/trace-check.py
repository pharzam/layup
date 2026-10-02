#!/usr/bin/env python3
"""T-4wrw: the one-off check of the PDR (plan step 1 red, step 4 green).

Evidence of #99, run by hand from the repository root:

    python3 runs/T-4wrw/trace-check.py [--prd PATH] [--record PATH]

No hook and no CI job runs it: it is not a check of the gate (Bootstrap mode
rule 1). It prints one line per check (PASS or FAIL, then the details), and exits
0 only when every check passes.

What it proves: the form and the resolution of each fact citation of PRD-0001;
that the record names each approved document by a commit of main that holds it;
that the record holds the approval as a link and a quote; that PRD-0001 has the
Status `Accepted`.
What it cannot prove: that a cited fact supports its requirement (the review
rounds of T-wjq4 and T-84r5 judged that); that the approval comment exists on the
forge, says what the record quotes, or was written by the Operator (the review
round reads the link).
"""
import re
import subprocess
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parents[2]
args = sys.argv[1:]
PRD = Path(args[args.index("--prd") + 1]) if "--prd" in args else ROOT / "docs/prd/PRD-0001-layup.md"
REC = Path(args[args.index("--record") + 1]) if "--record" in args else ROOT / "docs/pdr/PDR-0001.md"
DOCS = {  # name: (the path that the record shows, a path that git tests in the commit)
    "PRD-0001": ("docs/prd/PRD-0001-layup.md", "docs/prd/PRD-0001-layup.md"),
    "the architecture": ("docs/architecture.md", "docs/architecture.md"),
    "the specification": ("docs/spec/", "docs/spec/README.md"),
    "the plan": ("docs/plan/README.md", "docs/plan/README.md"),
}
results = []


def report(name, problems):
    results.append(problems)
    print(f"{'PASS' if not problems else 'FAIL'}  {name}")
    for p in problems[:30]:
        print(f"      {p}")


def cells(line):
    s = line.strip()
    if not s.startswith("|"):
        return None
    return [p.strip() for p in re.split(r"(?<!\\)\|", s.strip("|"))]


# the numbered facts of each record, read as check `facts` reads them
facts = {}
for f in sorted((ROOT / "docs/facts").glob("F-[0-9][0-9][0-9][0-9]-*.md")):
    rid = f.name[:6]
    nums = set()
    for l in f.read_text().split("\n"):
        m = re.match(r"^0*(\d+)\. ", l)
        if m:
            nums.add(int(m.group(1)))
    facts[rid] = nums
print("info  numbered facts: " + ", ".join(f"{k} {len(v)}" for k, v in facts.items()))

prd = PRD.read_text()
TOK = re.compile(r"F-\d{4}#\d+")

# 1. each REQ and NFR row of §6 and §7 cites at least one fact
rows, probs, sect = {}, [], None
for l in prd.split("\n"):
    m = re.match(r"^## (\d+)\.", l)
    if m:
        sect = m.group(1)
        continue
    c = cells(l)
    if sect in ("6", "7") and c and len(c) == 5 and re.fullmatch(r"(REQ|NFR)-\d{3}", c[0]):
        rows[c[0]] = TOK.findall(c[4])
        if not rows[c[0]] and c[2] != "Won't":
            probs.append(f"{c[0]}: no fact in its Facts cell")
        if not rows[c[0]] and c[2] == "Won't":
            probs.append(f"{c[0]} (Won't): no fact in its Facts cell")
report(f"each requirement of §6 and §7 cites a fact ({len(rows)} rows)", probs)

# 2. each fact token of the whole PRD resolves to a numbered fact
probs = []
alltok = TOK.findall(prd)
for t in sorted(set(alltok)):
    rid, n = t[:6], int(t.split("#")[1])
    if rid not in facts:
        probs.append(f"{t}: no record {rid} in docs/facts/")
    elif n not in facts[rid]:
        probs.append(f"{t}: {rid} has no numbered fact {n} ({len(facts[rid])} numbered facts)")
report(f"each fact token resolves ({len(set(alltok))} distinct tokens)", probs)

# 3. the §12 Facts cell of each row equals its §6 or §7 cell
probs, sect = [], None
for l in prd.split("\n"):
    m = re.match(r"^## (\d+)\.", l)
    if m:
        sect = m.group(1)
        continue
    c = cells(l)
    if sect == "12" and c and len(c) == 6 and re.fullmatch(r"(REQ|NFR)-\d{3}", c[0]):
        if TOK.findall(c[1]) != rows.get(c[0]):
            probs.append(f"{c[0]}: §12 cites {TOK.findall(c[1])}, §6/§7 cite {rows.get(c[0])}")
report("each §12 Facts cell equals its §6 or §7 cell", probs)

# 4. the record names each document by a commit of main that holds it
probs = []
if not REC.exists():
    probs.append(f"{REC} is missing")
else:
    rec = REC.read_text()
    for name, (shown, path) in DOCS.items():
        line = next((l for l in rec.split("\n") if l.startswith("|") and f"`{shown}`" in l), None)
        if line is None:
            probs.append(f"{name}: no table row names `{shown}`")
            continue
        shas = re.findall(r"`([0-9a-f]{7,40})`", line)
        if not shas:
            probs.append(f"{name}: its row names no commit")
            continue
        sha = shas[0]
        if subprocess.run(["git", "-C", str(ROOT), "merge-base", "--is-ancestor", sha, "origin/main"], capture_output=True).returncode != 0:
            probs.append(f"{name}: {sha} is not a commit of origin/main")
        elif subprocess.run(["git", "-C", str(ROOT), "cat-file", "-e", f"{sha}:{path}"], capture_output=True).returncode != 0:
            probs.append(f"{name}: {path} is not in {sha}")
report("the record names each document by a commit of main that holds it", probs)

# 5. the record holds the approval: the link and a quote
probs = []
if not REC.exists():
    probs.append("no record")
else:
    rec = REC.read_text()
    part = rec.split("## The approval", 1)[1] if "## The approval" in rec else ""
    if not re.search(r"https://github\.com/pharzam/layup/issues/99#issuecomment-\d+", part):
        probs.append("no link to the approval comment on #99")
    if not any(l.startswith("> ") and l[2:].strip() for l in part.split("\n")):
        probs.append("no quote of the approval")
    if "pending" in part.lower():
        probs.append("the approval is still pending")
report("the record holds the approval", probs)

# 6. PRD-0001 has the Status Accepted
st = re.search(r"^\| Status\s*\| `([^`]+)`", prd, re.M)
report("PRD-0001 has the Status `Accepted`", [] if st and st.group(1) == "Accepted" else [f"Status is {st.group(1) if st else 'missing'}"])

sys.exit(0 if all(not p for p in results) else 1)
