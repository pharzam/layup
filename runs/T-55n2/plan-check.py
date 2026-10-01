#!/usr/bin/env python3
"""T-55n2: the one-off check of LAYUP's implementation plan (plan step 1 red, step 4 green).

Evidence of #76, run by hand from the repository root:

    python3 runs/T-55n2/plan-check.py

No hook and no CI job runs it: it is not a check of the gate (Bootstrap mode
rule 1). It reads docs/plan/README.md, runs/T-55n2/inventory.md,
docs/prd/PRD-0001-layup.md and docs/tests/traceability.md, prints one line per
check (PASS or FAIL, then the details), and exits 0 only when every check
passes.

The tables it reads in docs/plan/README.md, by their headings:
  ## Milestones            | Milestone | Phase | Requirements | ... |
  ## The tasks of phase 1  | # | Task ID | Issue | Task | Parent | Items | Requirements | Fact | Tests | Size | Lines | After | Cap |
  ## The hosts of the other items  | Item | Host | Reason |
  ## The defect register   | K | Settled by | Read by | Note |
"""
import re
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parents[2]
PLAN = ROOT / "docs/plan/README.md"
INV = ROOT / "runs/T-55n2/inventory.md"
PRD = ROOT / "docs/prd/PRD-0001-layup.md"
TRACE = ROOT / "docs/tests/traceability.md"

PHASE1 = ["REQ-001", "REQ-002", "REQ-004", "REQ-007", "REQ-009", "REQ-011"] + [
    f"NFR-{n:03d}" for n in range(1, 8)
]
results = []


def report(name, problems, note=""):
    results.append((name, problems))
    head = "PASS" if not problems else "FAIL"
    print(f"{head}  {name}{(' — ' + note) if note else ''}")
    for p in problems[:40]:
        print(f"      {p}")
    if len(problems) > 40:
        print(f"      … and {len(problems) - 40} more")


def cells(line):
    s = line.strip()
    if not s.startswith("|"):
        return None
    s = s.strip("|")
    parts = re.split(r"(?<!\\)\|", s)
    return [p.strip().replace("\\|", "|") for p in parts]


def table_after(text, heading):
    """The first Markdown table after a heading line that starts with `heading`."""
    lines = text.split("\n")
    start = None
    for i, l in enumerate(lines):
        if l.strip().startswith(heading):
            start = i
            break
    if start is None:
        return None
    rows, header = [], None
    for l in lines[start + 1:]:
        if l.startswith("#"):
            break
        c = cells(l)
        if c is None:
            if header is not None and rows:
                break
            continue
        if header is None:
            header = c
            continue
        if all(re.fullmatch(r":?-+:?", x or "-") for x in c):
            continue
        rows.append(dict(zip(header, c)))
    return rows


def keys_in(cell):
    return re.findall(r"`([a-z0-9][a-z0-9-]*)`", cell or "")


def ids_in(cell, pat=r"(?:REQ|NFR)-\d{3}"):
    return re.findall(pat, cell or "")


# --- the inventory -------------------------------------------------------------
if not INV.exists():
    print(f"FAIL  inputs — {INV.relative_to(ROOT)} is missing")
    sys.exit(1)
inv = INV.read_text()
items, area = {}, None
for l in inv.split("\n"):
    m = re.match(r"## Area `([a-z]+)`", l)
    if m:
        area = m.group(1)
        continue
    m = re.match(r"### `([a-z0-9-]+)` — ", l)
    if m and area:
        items[m.group(1)] = area
conflicts = sorted({int(n) for n in re.findall(r"\*\*K(\d+)\.\*\*", inv)})
inv_hosts = {}
for l in inv.split("\n"):
    c = cells(l)
    if c and len(c) == 6 and re.fullmatch(r"`[a-z0-9-]+`", c[0]):
        inv_hosts[c[0].strip("`")] = c[5]
print(f"info  inventory: {len(items)} items, {len(conflicts)} conflicts")

# --- the plan ------------------------------------------------------------------
if not PLAN.exists():
    report("plan exists", [f"{PLAN.relative_to(ROOT)} is missing"])
    for name in ["phase-1 requirements have a task", "each task names a requirement",
                 "predecessors resolve with no cycle", "each inventory item has one host",
                 "each conflict has a host", "PRD-0001 section 12 Task cells",
                 "traceability table", "inventory Host column agrees with the plan"]:
        report(name, ["no plan"])
    sys.exit(1)
plan = PLAN.read_text()
report("plan holds no marker", [f"line {i}" for i, l in enumerate(plan.split("\n"), 1) if "\u2039" in l])

milestones = table_after(plan, "## Milestones") or []
mids = {r.get("Milestone", "").strip("`") for r in milestones}
tasks = table_after(plan, "## The tasks of phase 1") or []
others = table_after(plan, "## The hosts of the other items") or []
register = table_after(plan, "## The defect register") or []
print(f"info  plan: {len(milestones)} milestones, {len(tasks)} tasks, {len(others)} other hosts, {len(register)} register rows")

rows = {}
for r in tasks:
    n = r.get("#", "").strip()
    if n.isdigit():
        rows[int(n)] = r

# 1. each phase-1 requirement has a task
have = set()
for r in rows.values():
    have.update(ids_in(r.get("Requirements")))
report("phase-1 requirements have a task", [f"{q} has no task" for q in PHASE1 if q not in have])

# 2. each task names a requirement, a task ID and an In-Scope fact
probs = []
for n, r in sorted(rows.items()):
    if not ids_in(r.get("Requirements")):
        probs.append(f"row {n}: no requirement ID")
    if not re.fullmatch(r"`?T-[0-9a-hj-km-np-tv-z]{4}`?", r.get("Task ID", "")):
        probs.append(f"row {n}: task ID {r.get('Task ID')!r} is not of the form T-xxxx")
    if not re.search(r"F-0003#\d+", r.get("Fact", "")):
        probs.append(f"row {n}: no In-Scope fact F-0003#n")
    if r.get("Size") not in ("small", "large"):
        probs.append(f"row {n}: size {r.get('Size')!r}")
report("each task names a requirement", probs)

# 3. predecessors resolve, with no cycle
probs, graph = [], {}
for n, r in rows.items():
    after = r.get("After", "")
    deps = [int(x) for x in re.findall(r"\d+", after)] if after not in ("—", "-", "") else []
    graph[n] = deps
    for d in deps:
        if d not in rows:
            probs.append(f"row {n}: predecessor {d} is not a row")
        elif d == n:
            probs.append(f"row {n}: depends on itself")
state = {}


def visit(n, path):
    if state.get(n) == 1:
        probs.append("cycle: " + " -> ".join(map(str, path + [n])))
        return
    if state.get(n) == 2:
        return
    state[n] = 1
    for d in graph.get(n, []):
        if d in graph:
            visit(d, path + [n])
    state[n] = 2


for n in sorted(graph):
    visit(n, [])
report("predecessors resolve with no cycle", probs)

# 4. each inventory item has exactly one host
host = {}
probs = []
for n, r in rows.items():
    for k in keys_in(r.get("Items")):
        host.setdefault(k, []).append(f"row {n}")
for r in others:
    for k in keys_in(r.get("Item")):
        h = r.get("Host", "").strip()
        if not h or not r.get("Reason", "").strip():
            probs.append(f"{k}: other host without a host or a reason")
        host.setdefault(k, []).append(h)
for k in items:
    if k not in host:
        probs.append(f"{k} ({items[k]}): no host")
    elif len(host[k]) > 1:
        probs.append(f"{k}: {len(host[k])} hosts: {', '.join(host[k])}")
for k in host:
    if k not in items:
        probs.append(f"{k}: named in the plan, but not an inventory item")
report("each inventory item has one host", probs)

# 5. each conflict has a host
seen, probs = {}, []
for r in register:
    m = re.fullmatch(r"K(\d+)", r.get("K", "").strip())
    if not m:
        continue
    k = int(m.group(1))
    if not r.get("Settled by", "").strip() or r.get("Settled by", "").strip() in ("—", "-"):
        probs.append(f"K{k}: no host")
    seen[k] = seen.get(k, 0) + 1
for k in conflicts:
    if k not in seen:
        probs.append(f"K{k}: not in the register")
    elif seen[k] > 1:
        probs.append(f"K{k}: {seen[k]} register rows")
report("each conflict has a host", probs)

# 6. PRD-0001 section 12 Task cells
probs = []
tids = {r.get("Task ID", "").strip("`") for r in rows.values()}
prd = PRD.read_text().split("\n")
in12 = False
for l in prd:
    if l.startswith("## 12."):
        in12 = True
        continue
    if in12 and l.startswith("## "):
        break
    c = cells(l) if in12 else None
    if c and len(c) == 6 and re.fullmatch(r"(?:REQ|NFR)-\d{3}", c[0]):
        task = c[4]
        named = set(re.findall(r"T-[0-9a-z]{4}", task)) | set(re.findall(r"`?(M\d+[a-z]?)`?", task))
        if not named:
            probs.append(f"{c[0]}: Task cell {task!r} names no task and no milestone")
        for t in re.findall(r"T-[0-9a-z]{4}", task):
            if t not in tids and not (ROOT / f"docs/tasks/{t}.md").exists():
                probs.append(f"{c[0]}: {t} is not a task of the plan and has no task file")
        for m in re.findall(r"`(M\d+[a-z]?)`", task):
            if m not in mids:
                probs.append(f"{c[0]}: milestone {m} is not in the plan")
report("PRD-0001 section 12 Task cells", probs)

# 7. the traceability table
probs = []
if not TRACE.exists():
    probs.append("docs/tests/traceability.md is missing")
else:
    t = TRACE.read_text()
    if "\u2039" in t:
        probs.append("it holds a marker")
    req = set(re.findall(r"^\| ((?:REQ|NFR)-\d{3}) \|", PRD.read_text(), re.M))
    covered = set()
    for l in t.split("\n"):
        c = cells(l)
        if not c or len(c) < 8 or c[0].startswith("Test ID") or set(c[0]) <= {"-", ":"}:
            continue
        level, status = c[1].strip("`"), c[7].strip("`")
        ids = set(ids_in(c[2]))
        covered |= ids
        if level not in ("unit", "integration", "e2e", "uat"):
            probs.append(f"{c[0]}: level {level!r}")
        if status not in ("planned", "red", "green", "frozen"):
            probs.append(f"{c[0]}: status {status!r}")
        for i in ids:
            if i not in req:
                probs.append(f"{c[0]}: {i} is not a requirement of PRD-0001")
    for q in sorted(req - covered):
        probs.append(f"{q} has no row")
report("traceability table", probs)

# 8. the Host column of the inventory agrees with the plan
probs = []
for k, h in inv_hosts.items():
    want = ", ".join(host.get(k, []))
    if h != want:
        probs.append(f"{k}: inventory says {h!r}, plan says {want!r}")
report("inventory Host column agrees with the plan", probs)

sys.exit(0 if all(not p for _, p in results) else 1)
