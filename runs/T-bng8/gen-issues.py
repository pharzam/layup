#!/usr/bin/env python3
"""The issues and the plan table of the build tasks of M2d (task T-bng8, #267).

Copied from runs/T-2r8u/gen-issues.py and changed for M2d: one list, ROWS, holds
the 15 rows of O-209 (#267); the parts of the specification of each row come
from runs/T-bng8/parts.tsv, so the issue, the table and the check
parts.sh read one list. With --post, the script opens one issue per row through
the layup-agent App (the token from ~/.config/layup-agent/app-token.sh) and
writes runs/T-bng8/issues.tsv (row, task, issue). With --body ROW, it prints the
body of one row's issue; with --edit ROW..., it writes the body of each named
row's issue again (the fixes of round 1). With --table, it prints the table of
docs/plan/README.md from ROWS, parts.tsv and issues.tsv. A one-off evidence
command, not a check of the gate.

    python3 -I runs/T-bng8/gen-issues.py --body 40
    python3 -I runs/T-bng8/gen-issues.py --post
    python3 -I runs/T-bng8/gen-issues.py --table
"""
import os
import subprocess
import sys
import tempfile

HERE = os.path.dirname(os.path.abspath(__file__))
ISSUES = os.path.join(HERE, "issues.tsv")
PARTS = os.path.join(HERE, "parts.tsv")
SPEC = "https://github.com/pharzam/layup/blob/main/docs/spec/"
ITEM = "`later-p2-scaffold`"
ITEMS = {}


def item(row):
    """The item of a row (note 4 of the plan review of #211)."""
    return ITEMS.get(row, ITEM) + (", `gov-operator-setup-o112`" if row == "59b" else "")

# row, task, title, demo, packages, requirements, fact, tests, size, lines, after, classes
ROWS = [
    ("50", "T-tetz", "The records of the Scaffold",
     "The Go schemas of `markers`, `setup-tree` and `forge-check` equal their blocks, and each reader refuses a row that breaks its rule.",
     "`internal/records` (`markers`), `internal/run` (`setup-tree`, `forge-check`)", "NFR-001, REQ-002, REQ-003", "F-0001#1, F-0003#42",
     "unit; the block test (the three names to `built`)", "small", "400", "—",
     "2: the Go schemas equal their blocks; the row rules in words"),
    ("51a", "T-jpkc", "The forge calls of M2d",
     "Against the `httptest` forge, `EffectiveRules` gives each active rule of a branch with its parameters, `Activity` gives the activities of a ref, and `Repository` gives the head of each branch.",
     "`internal/forge`, `internal/forge/github`", "REQ-003, NFR-001", "F-0003#43",
     "integration (the `httptest` forge)", "small", "400", "—",
     "1: the calls of `M2d`"),
    ("51b", "T-j4sd", "The call CommitTree",
     "`CommitTree` makes an empty commit on a parent with the fixed identity, and refuses a tree or a parent that is not a full object ID before `git` starts.",
     "`internal/git`", "REQ-003", "F-0003#43",
     "unit; integration (the real `git`)", "small", "200", "—",
     "1: the call `CommitTree`"),
    ("52a", "T-ktgx", "The setup in phase 2",
     "The step runner with the inputs of phase 2 reads only `S01-stack`, takes the pin rows from `start.tsv` with resolving nothing, makes no commit at S03 and none of `layup-records` at S15, and writes the target's pin file from the rows.",
     "`internal/setup`", "REQ-002, NFR-006, NFR-003", "F-0003#42, F-0001#8",
     "unit (each step of phase 2)", "large", "600", "—",
     "1: the setup in phase 2"),
    ("52b", "T-9v92", "The rulesets of S13",
     "S13 writes the ruleset of the default branch with `creation`, the merge method and, in phase 2, the four `layup/` checks pinned to the App, and in phase 2 the ruleset of the records branch and the row `without-layup`.",
     "`internal/setup`", "REQ-003, NFR-002", "F-0003#43, F-0003#65",
     "unit (the two files in both phases; the record row)", "small", "400", "52a",
     "2: the ruleset files; the row `without-layup`"),
    ("52c", "T-56dj", "The commands of S13",
     "`commands.sh` stops at the first command that fails, and a run of it again after the push or after one apply pushes nothing and applies only the missing ruleset.",
     "`internal/setup`", "REQ-003", "F-0003#43",
     "unit (stand-ins of `git` and `gh`)", "small", "300", "52b",
     "1: `commands.sh`, fail-fast and safe to run again"),
    ("53", "T-4c5c", "The step markers",
     "A restart runs `markers` before `batch`: the markers of the pinned tree with their checked proposals are committed, and the batch asks `S01-stack` and each marker in the Operator's block.",
     "`internal/run`, `internal/records` (the enum of `questions`)", "REQ-002, NFR-003", "F-0003#42, F-0003#63",
     "unit (the scan, the result file of the proposals); integration (the step, the batch, a scripted harness)", "large", "700", "50",
     "2: the scan, the session and the record; the merge and the comment"),
    ("54", "T-eahr", "The step setup",
     "After Intake, a restart builds the work area from the records, answers each stop of the runner from them or with a prose session, and runs S01 to S14.",
     "`internal/run`", "REQ-002, NFR-001, NFR-003", "F-0003#42, F-0001#1",
     "unit (the answers, the stops); integration (the real step runner, a scripted harness)", "large", "800", "52a, 52c, 53",
     "2: the work area and the answers; the stops and the prose rows"),
    ("55", "T-79zk", "The step verify",
     "After S14, `layup run` runs `layup setup verify` and S15 and commits the setup's records, `setup/tree.tsv` and the README of the Scaffold in one records commit.",
     "`internal/run`", "REQ-002, NFR-001", "F-0003#42, F-0001#1",
     "integration (`layup run` with the real step runner)", "small", "400", "54",
     "1: the step `verify` and its records commit"),
    ("56", "T-hfn1", "The step scaffold-push",
     "The step prints the Operator's commands and waits, commits `setup/rulesets.txt` from both apply outputs, and builds a missing work area again from the records.",
     "`internal/run`", "REQ-003, NFR-001", "F-0003#43",
     "unit (the apply outputs); integration (the wait, the work area built again)", "small", "500", "52c, 55",
     "2: the wait and the apply outputs; the work area built again"),
    ("57a", "T-wphp", "The step forge-check",
     "The forge check commits one row per check: the tree, the rules read back, the bypass lists, the probe refused with `GH013` and the ruleset's name, and the activity read.",
     "`internal/run`", "REQ-003, NFR-001", "F-0003#43, F-0003#64",
     "unit (the read-back); integration (the probe against a bare repository with a refusing hook; the `httptest` forge)", "large", "700", "50, 51a, 51b, 56",
     "2: the read-back and the tree; the probe and the activity"),
    ("57b", "T-2j13", "The forge check at each later start",
     "Once the Scaffold is done, each restart runs the forge check after `lease` and before `probe`, and stops at a row that is not `pass`.",
     "`internal/run`", "REQ-003", "F-0003#43",
     "integration", "small", "300", "57a",
     "1: the forge check at each later start"),
    ("58", "T-eevt", "The e2e of the Scaffold",
     "In CI with no secret, the restarts run from Intake's answers through the setup, the printed commands run by the test in place of the Operator, and the forge check.",
     "`cmd/layup` (the e2e test)", "REQ-002, REQ-003, NFR-001", "F-0003#42, F-0003#43",
     "e2e (the binary, a local fake forge, local bare repositories, scripted harness programs)", "large", "600", "57b",
     "1: the steps of the Scaffold, end to end"),
    ("59a", "T-5xjc", "The review of the release of M2d",
     "The release review of `M2d` is a recorded review of the code for `REQ-015` and `REQ-017`, by a reviewer of a model that wrote none of it.",
     "—", "REQ-015, REQ-017", "F-0003#43 (the review keeps the out-of-scope facts F-0003#53 and F-0003#55 out)",
     "the review, with `runs/T-efmy/release-check.sh` on the release commit, adapted where the code of `M2d` needs it (the push of the probe reaches a remote); its evidence under `runs/`", "small", "300", "58",
     "2: the deterministic check of `release-check.sh`; the review by a reviewer of another model"),
    ("59b", "T-sdtb", "The demo of M2d (uat)",
     "After the Operator's pushes and ruleset apply, the push to `layup-probe` is refused, and `setup/forge-check.tsv` shows each rule and probe as passed.",
     "—", "REQ-002, REQ-003", "F-0003#42, F-0003#43",
     "uat, with its evidence under `runs/`", "small", "300", "59a",
     "1: the demo"),
]

CAPS = {}

OPERATOR_INPUTS = """## The Operator's inputs (R6: severity normal; needed when row 59a merges)

- A test target that Start makes from the host (`layup run --new`), with the App installed on it, and its Intake answered, the setup questions included ([`scaffold.md`]({spec}scaffold.md#req-002--the-setup-questions-in-the-one-batch)).
- The admitted harnesses and the routing of `systems-architect` (the proposals of the markers) and `product-owner` (the prose rows), at the tier `reasoning`.
- The Operator's admin login on the target, which runs the printed pushes and the apply of the two rulesets (`architecture.md` §5, Scaffold 6), in the SSH form.
- The plan of the target: public, or a plan that enforces rulesets (K15).
""".format(spec=SPEC)

REVIEW_NOTE = """## The release

The release of `M2d` is the code of `main` that the demo of row 59b runs (the non-test Go files of `cmd/` and `internal/`, `go.mod`, and the files that the binary embeds), at the merge of the last build row, 58; the review reads the whole release ([`scaffold.md`]({spec}scaffold.md#req-015-and-req-017--the-review-of-the-release-of-m2d)).
""".format(spec=SPEC)


COMMON = """- [ ] The demo above holds, by the tests of the Tests cell, written red first.
- [ ] Each part above is built as the specification says; a difference is a defect of the specification, fixed in the same pull request (R10).
- [ ] `go build ./...`, `go vet ./...`, `go test ./...` and `go test -tags=integration ./...` pass; `adr-lint`, `prd-lint`, `link-lint`, `setup-check` and `run-discipline-tests` exit 0.
"""
REVIEW = """- [ ] The review is recorded under `runs/`, by a reviewer of a model that wrote none of the release, with the result of `release-check.sh` at the reviewed commit.
- [ ] Each finding of the review is fixed, or opens an issue with its reason.
"""
UAT = """- [ ] The demo above holds on a test target that Start made, on the Operator's harnesses, with its evidence under `runs/` (`setup/forge-check.tsv`, `setup/rulesets.txt`, the refused push of the probe).
- [ ] Each write of LAYUP on the target shows the App's bot as its author; the setup push and the apply are the Operator's.
"""


def parts():
    out = {}
    with open(PARTS) as f:
        for ln in f.read().splitlines()[1:]:
            file, heading, sub, row = ln.split("\t")
            name = file if file.endswith("README.md") else file.rsplit("/", 1)[1]
            h = heading.lstrip("#").strip().strip("*")
            text = f"`{name}`: {h}" + ("" if sub == "—" else f" ({sub})")
            out.setdefault(row, []).append(text)
    return out


def body(r, p):
    row, task, title, demo, pkgs, reqs, fact, tests, size, lines, after, classes = r
    cap = CAPS.get(row, "1")
    crit = COMMON
    if row == "59a":
        crit = REVIEW + COMMON.split("\n", 2)[2]
    elif row == "59b":
        crit = UAT
    items = item(row)
    b = f"""## Goal

Task `{task}`: row {row} of the [implementation plan](https://github.com/pharzam/layup/blob/main/docs/plan/README.md#the-tasks-of-m2d), {title}. Milestone `M2d`; the specification is [`docs/spec/scaffold.md`]({SPEC}scaffold.md), [`records.md`]({SPEC}records.md#nfr-001--the-records-of-the-scaffold), [`forge.md`]({SPEC}forge.md#the-calls-of-m2d) and [`packages.md`]({SPEC}packages.md#the-table-of-m2d). Fact: {fact}. Opened by `T-bng8` (#267). Refs #265, Refs #29.

**The demo** (R11): {demo}

**The goal count** (R11), final by the Operator's decision O-209 of #267: {classes}.

## Scope

- **The parts of the specification:** {"; ".join(p[row])}.
- **Packages:** {pkgs}. **Requirements:** {reqs}. **Items:** {items}.
- **Tests:** {tests}.
- **Size:** {size}, about {lines} lines (an estimate; the plan review sets the budget). **After:** {after}. **Expected cap:** {cap}.

## Acceptance criteria

{crit}- [ ] Docs updated in the same PR (`docs/tests/traceability.md`, the Test cells of `PRD-0001` §12).
"""
    if row == "59a":
        b += "\n" + REVIEW_NOTE
    if row == "59b":
        b += "\n" + OPERATOR_INPUTS
    return b + "\n*Written by Claude Opus 5.5, the author of T-bng8 (#267).*\n"


def dup(title, env):
    """The duplicate check of one child issue (R2; condition 4 of the plan review
    of #267): the open and closed issues whose title holds the row's title."""
    q = f'repo:pharzam/layup in:title "{title}"'
    out = subprocess.run(["gh", "api", "-X", "GET", "search/issues", "-f", f"q={q}",
                          "--jq", '[.items[] | "#\\(.number) (\\(.title))"] | join("; ")'],
                         capture_output=True, text=True, env=env, stdin=subprocess.DEVNULL,
                         timeout=60, check=True).stdout.strip()
    day = subprocess.run(["date", "-u", "+%Y-%m-%d"], capture_output=True, text=True).stdout.strip()
    found = out if out else "none"
    return (f"\n## Duplicate check (R2)\n\n- [x] Searched the open and closed issues of "
            f"`pharzam/layup` for the title \"{title}\" on {day}, before this issue "
            f"opened: {found}. None of them is this task: each is a task of another "
            f"milestone or of the specification. Related: #265 (the specification), #267 (the slicing).\n")


def post():
    p = parts()
    token = subprocess.run([os.path.expanduser("~/.config/layup-agent/app-token.sh")],
                           capture_output=True, text=True, check=True).stdout.strip()
    env = dict(os.environ)
    env["GH_TOKEN"] = token
    env["GH_PROMPT_DISABLED"] = "1"
    lines = ["row\ttask\tissue"]
    for r in ROWS:
        b = body(r, p)
        sig = b.rindex("\n*Written by")
        b = b[:sig] + dup(r[2], env) + b[sig:]
        with tempfile.NamedTemporaryFile("w", suffix=".md", delete=False) as f:
            f.write(b)
        title = f"{r[1]}: {r[2]} (row {r[0]}, M2d)"
        out = subprocess.run(["gh", "issue", "create", "--repo", "pharzam/layup", "--title", title,
                              "--body-file", f.name], capture_output=True, text=True, env=env,
                             stdin=subprocess.DEVNULL, timeout=60, check=True).stdout.strip()
        os.unlink(f.name)
        num = out.rsplit("/", 1)[1]
        print(r[0], r[1], out)
        lines.append(f"{r[0]}\t{r[1]}\t{num}")
        with open(ISSUES, "w") as f:
            f.write("\n".join(lines) + "\n")


def edit(rows):
    """Write the body of each named row's issue again, as the App."""
    p = parts()
    nums = {}
    with open(ISSUES) as f:
        for ln in f.read().splitlines()[1:]:
            row, task, num = ln.split("\t")
            nums[row] = num
    token = subprocess.run([os.path.expanduser("~/.config/layup-agent/app-token.sh")],
                           capture_output=True, text=True, check=True).stdout.strip()
    env = dict(os.environ)
    env["GH_TOKEN"] = token
    env["GH_PROMPT_DISABLED"] = "1"
    for r in ROWS:
        if r[0] not in rows:
            continue
        with tempfile.NamedTemporaryFile("w", suffix=".md", delete=False) as f:
            f.write(body(r, p))
        subprocess.run(["gh", "issue", "edit", nums[r[0]], "--repo", "pharzam/layup",
                        "--body-file", f.name], env=env, stdin=subprocess.DEVNULL,
                       timeout=60, check=True)
        os.unlink(f.name)


def table():
    p = parts()
    nums = {}
    with open(ISSUES) as f:
        for ln in f.read().splitlines()[1:]:
            row, task, num = ln.split("\t")
            nums[row] = num
    print("| # | Task ID | Issue | Task | Demo | Items | Requirements | Fact | Tests | Size | Lines | After | Cap |")
    print("| - | ------- | ----- | ---- | ---- | ----- | ------------ | ---- | ----- | ---- | ----- | ----- | --- |")
    for r in ROWS:
        row, task, title, demo, pkgs, reqs, fact, tests, size, lines, after, classes = r
        n = nums[row]
        items = item(row)
        print(f"| {row} | `{task}` | [#{n}](https://github.com/pharzam/layup/issues/{n}) | {title}: {'; '.join(p[row])} | {demo} | {items} | {reqs} | {fact} | {tests} | {size} | {lines} | {after} | {CAPS.get(row, '1')} |")


if __name__ == "__main__":
    if sys.argv[1:] == ["--post"]:
        post()
    elif sys.argv[1:] == ["--table"]:
        table()
    elif len(sys.argv) >= 3 and sys.argv[1] == "--edit":
        edit(sys.argv[2:])
    elif len(sys.argv) == 3 and sys.argv[1] == "--body":
        r = [x for x in ROWS if x[0] == sys.argv[2]]
        if not r:
            sys.exit("no such row")
        print(body(r[0], parts()))
    else:
        sys.exit(__doc__)
