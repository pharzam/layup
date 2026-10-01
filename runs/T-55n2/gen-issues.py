#!/usr/bin/env python3
"""T-55n2 (plan step 7, O-125): write the issue texts of the plan's tasks by code.

Run by hand from the repository root:

    python3 runs/T-55n2/gen-issues.py FIRST_NUMBER

FIRST_NUMBER is the forge number that the first new issue gets (the parent task
`layup setup`); the rows follow in order. The script reads docs/plan/README.md and
runs/T-55n2/inventory.md, and writes one file per issue under runs/T-55n2/issues/:
the title on line 1, a blank line, then the body. It is evidence, not a check.
"""
import re
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parents[2]
OUT = ROOT / "runs/T-55n2/issues"
BLOB = "https://github.com/pharzam/layup/blob/main"
first = int(sys.argv[1])


def cells(line):
    s = line.strip()
    if not s.startswith("|"):
        return None
    return [p.strip().replace("\\|", "|") for p in re.split(r"(?<!\\)\|", s.strip("|"))]


def table_after(text, heading):
    lines = text.split("\n")
    start = next((i for i, l in enumerate(lines) if l.strip().startswith(heading)), None)
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


plan = (ROOT / "docs/plan/README.md").read_text()
inv = (ROOT / "runs/T-55n2/inventory.md").read_text()
tasks = {int(r["#"]): r for r in table_after(plan, "## The tasks of phase 1")}
register = table_after(plan, "## The defect register")
titles = dict(re.findall(r"^### `([a-z0-9-]+)` — (.+)$", inv, re.M))

number = {n: first + n for n in tasks}          # row n -> issue number
SETUP = first                                   # the parent task `layup setup`
PARENT = {"T-vk3k": ("#33", "layup gate"), "T-b97r": (f"#{SETUP}", "layup setup"), "T-stfn": ("#29", "the core engine")}
FACT = {"F-0003#41": "Problem Statement Quality", "F-0003#42": "Reproducible Discipline Setup",
        "F-0003#44": "Stack-Dependent Gates", "F-0003#47": "Verification on Every Change",
        "F-0003#49": "Stall Resolution", "F-0003#50": "Cost Visibility"}


def rowset(cell):
    """The row numbers that a register cell names: "row 7", "rows 7, 8, 13", "rows 1 to 5, 14"."""
    out = set()
    for seg in re.findall(r"\brows? ([0-9][0-9 ,toand]*)", cell or ""):
        nums = re.findall(r"\d+|to", seg)
        i = 0
        while i < len(nums):
            if i + 2 < len(nums) and nums[i + 1] == "to":
                out.update(range(int(nums[i]), int(nums[i + 2]) + 1))
                i += 3
            elif nums[i] != "to":
                out.add(int(nums[i]))
                i += 1
            else:
                i += 1
    return out


def settles(n):
    return [r for r in register if n in rowset(r.get("Settled by", ""))]


def reads(n):
    return [r for r in register if n in rowset(r.get("Read by", ""))]


def body(n, r):
    tid = r["Task ID"].strip("`")
    pnum, pname = PARENT[r["Parent"].strip("`")]
    keys = re.findall(r"`([a-z0-9][a-z0-9-]*)`", r["Items"])
    after = [int(x) for x in re.findall(r"\d+", r["After"])] if r["After"] not in ("—", "") else []
    starts = ("**It starts now** (O-126: the hold of #42 is lifted for rows 1 and 2; it needs no other task)."
              if n in (1, 2) else
              "**It is on hold until #42 (the PDR) closes**" + (", and until its predecessors merge." if after else "."))
    facts = re.findall(r"F-0003#\d+", r["Fact"])
    s = []
    s.append("## Goal\n")
    s.append(f"Task `{tid}`, row {n} of the [implementation plan]({BLOB}/docs/plan/README.md#the-tasks-of-phase-1) (#76). "
             f"Child of {pnum} (`{pname}`). The goal: {r['Task']}.\n")
    s.append(starts + "\n")
    s.append("## Requirements\n")
    s.append(f"- {r['Requirements']} ([`PRD-0001`]({BLOB}/docs/prd/PRD-0001-layup.md)).")
    s.append("- In-Scope fact (Bootstrap mode rule 1): " + ", ".join(f"`{f}` ({FACT.get(f, '')})" for f in facts) + ".\n")
    s.append("## Scope\n")
    s.append(f"The inventory items of this task. Each one has its specification sections, its present state, its tests and its open questions in [`runs/T-55n2/inventory.md`]({BLOB}/runs/T-55n2/inventory.md); the open questions are inputs of this task's plan.\n")
    for k in keys:
        s.append(f"- `{k}` — {titles.get(k, '')}")
    s.append("")
    st, rd = settles(n), [x for x in reads(n) if x not in settles(n)]
    if st:
        s.append("**The defects that this task settles** in `docs/spec/`, in its own pull request (R10; the [defect register]"
                 f"({BLOB}/docs/plan/README.md#the-defect-register)):\n")
        for x in st:
            s.append(f"- **{x['K']}** — {x['Note']}")
        s.append("")
    if rd:
        s.append("**The defects that this task reads**, settled before it:\n")
        for x in rd:
            s.append(f"- **{x['K']}** (settled by {x['Settled by']}) — {x['Note']}")
        s.append("")
    s.append("**Out:** the parts that the plan gives to other rows, and the parts \"Not in phase 1\" of its specification sections.\n")
    s.append("## Duplicate check (R2)\n")
    s.append("- [x] Searched the open and closed issues when the plan was written (#76): no duplicate. Related: #29, "
             f"{pnum}, and the predecessors below.\n")
    s.append("## Solution note (R3)\n")
    s.append("- **Chosen:** the contract of the specification sections of its items (`docs/spec/`); the design inside it is this task's own plan (R12).")
    s.append("- **Why:** the specification decides the commands, the records and the rules, and the plan decides the order.")
    s.append("- **Rejected:** this task's plan compares the alternatives.")
    s.append("- **Important tradeoffs:** none recorded by the plan.")
    s.append("- **Decision record:** this issue (the task's plan), or an ADR if its plan finds an architecturally significant decision.\n")
    s.append("## Acceptance criteria\n")
    s.append("- [ ] Each inventory item above is delivered as its specification sections say.")
    s.append(f"- [ ] The tests of the plan pass: {r['Tests']}.")
    if st:
        s.append("- [ ] Each defect that this task settles is settled in `docs/spec/`: a value \"decided here\" with its reason, or a marker with a row in `docs/setup/open-gaps.tsv`.")
    s.append(f"- [ ] The rows of [`docs/tests/traceability.md`]({BLOB}/docs/tests/traceability.md) for its tests are `green`, and `PRD-0001` §12 names its tests in the Test column.")
    s.append("- [ ] Tests cover the change and pass (R8).")
    s.append("- [ ] Docs updated in the same PR.\n")
    s.append("## Notes\n")
    s.append(f"- Size class: `{r['Size']}`; expected lines: about {r['Lines']} (an estimate for the plan review, not a budget).")
    s.append(f"- Cycle cap: {r['Cap']} (Bootstrap mode rule 3).")
    if after:
        s.append("- After: " + ", ".join(f"#{number[a]} (row {a})" for a in after) + ".")
    else:
        s.append("- After: no other task.")
    s.append(f"- Refs #29, {pnum}.")
    s.append("\n*Written by code from the plan ([`runs/T-55n2/gen-issues.py`]"
             f"({BLOB}/runs/T-55n2/gen-issues.py)), posted by the author (Claude Opus 5.5) through the `layup-agent` App.*")
    return "\n".join(s) + "\n"


OUT.mkdir(parents=True, exist_ok=True)
children = [n for n, r in tasks.items() if r["Parent"].strip("`") == "T-b97r"]
setup = ["## Goal\n",
         "Task `T-b97r`: `layup setup` and `layup setup verify`, the setup of a target repository (`REQ-002`), one of the four engine "
         "deliverables of ADR-0012 part 5 (O-121). It is the parent of its child tasks below; each child has its own plan and review. "
         f"Its specification is [`docs/spec/setup.md`]({BLOB}/docs/spec/setup.md); its order is the "
         f"[implementation plan]({BLOB}/docs/plan/README.md#the-tasks-of-phase-1) (#76).\n",
         "**It is on hold until #42 (the PDR) closes.**\n",
         "## Child tasks (in order; each opened with its own issue)\n",
         "| Row | Task | Issue | Goal |", "| --- | ---- | ----- | ---- |"]
for n in sorted(children):
    setup.append(f"| {n} | `{tasks[n]['Task ID'].strip('`')}` | #{number[n]} | {tasks[n]['Task']} |")
setup += ["",
          "## Duplicate check (R2)\n",
          "- [x] Searched when the plan was written (#76). `T-b97r` is the ID that #29's first children table gave to `layup setup`; this is its first issue.\n",
          "## Acceptance criteria\n",
          "- [ ] Each child task is closed with a merged pull request.",
          f"- [ ] The phase-1 parts of `REQ-002` hold, as the plan's table \"What phase 1 proves\" says ([plan]({BLOB}/docs/plan/README.md#what-phase-1-proves)).\n",
          "## Notes\n",
          "- In-Scope fact: `F-0003#42` (Reproducible Discipline Setup).",
          "- Refs #29.",
          "\n*Written by code from the plan ([`runs/T-55n2/gen-issues.py`]"
          f"({BLOB}/runs/T-55n2/gen-issues.py)), posted by the author (Claude Opus 5.5) through the `layup-agent` App.*"]
(OUT / f"{SETUP:03d}-T-b97r.md").write_text("T-b97r: `layup setup` and `layup setup verify` — the setup of a target (a deliverable of ADR-0012 part 5)\n\n" + "\n".join(setup) + "\n")
for n, r in sorted(tasks.items()):
    tid = r["Task ID"].strip("`")
    (OUT / f"{number[n]:03d}-{tid}.md").write_text(f"{tid}: {r['Task']}\n\n" + body(n, r))
print(f"wrote {len(tasks) + 1} files in {OUT.relative_to(ROOT)}")
