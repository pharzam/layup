#!/usr/bin/env python3
"""The issues and the plan table of the build tasks of M2c (task T-2r8u, #211).

Copied from runs/T-2r8u/gen-issues.py and changed for M2c: one list, ROWS, holds
the 17 rows of O-196 (#211); the parts of the specification of each row come
from runs/T-2r8u/parts.tsv, so the issue, the table and the check
parts.sh read one list. With --post, the script opens one issue per row through
the layup-agent App (the token from ~/.config/layup-agent/app-token.sh) and
writes runs/T-2r8u/issues.tsv (row, task, issue). With --body ROW, it prints the
body of one row's issue; with --edit ROW..., it writes the body of each named
row's issue again (the fixes of round 1). With --table, it prints the table of
docs/plan/README.md from ROWS, parts.tsv and issues.tsv. A one-off evidence
command, not a check of the gate.

    python3 -I runs/T-2r8u/gen-issues.py --body 40
    python3 -I runs/T-2r8u/gen-issues.py --post
    python3 -I runs/T-2r8u/gen-issues.py --table
"""
import os
import subprocess
import sys
import tempfile

HERE = os.path.dirname(os.path.abspath(__file__))
ISSUES = os.path.join(HERE, "issues.tsv")
PARTS = os.path.join(HERE, "parts.tsv")
SPEC = "https://github.com/pharzam/layup/blob/main/docs/spec/"
ITEM = "`later-p2-intake-spec`"
ITEMS = {"43b": "`later-psb-meaning-review`", "46": "`later-psb-intake-batch`", "47a": "`later-psb-intake-batch`",
         "47b": "`later-psb-intake-batch`", "47c": "`later-psb-intake-batch`"}


def item(row):
    """The item of a row (note 4 of the plan review of #211)."""
    return ITEMS.get(row, ITEM) + (", `gov-operator-setup-o112`" if row == "49b" else "")

# row, task, title, demo, packages, requirements, fact, tests, size, lines, after, classes
ROWS = [
    ("40", "T-529i", "The records of Intake",
     "The Go schemas of `meaning`, `spans`, `draft`, `coverage`, `questions`, `answers`, `answers-review` and `budget` equal their blocks, and each reader refuses a row that breaks its rule.",
     "`internal/records`", "NFR-001, REQ-001, REQ-012", "F-0003#41, F-0003#51",
     "unit; the block test (the eight names to `built`)", "large", "800", "—",
     "2: the eight Go schemas equal their blocks; the row rules in words"),
    ("41a", "T-4ywt", "The span check and the package rules of M2c",
     "The span check refuses a span table with a byte in two spans, and `TestPackageRules` reads \"The table of M2c\" with `internal/spec` on the line of the engine checks.",
     "`internal/spec`, `cmd/layup` (the package rules)", "REQ-012, NFR-005, NFR-007", "F-0003#51",
     "unit (each rule of the check); integration (`TestPackageRules`); the block test (`spec-check` to `built`; the second Go value of `spans` compared)", "large", "600", "—",
     "2: the five rules of the span check and its table; the package rules of `M2c`"),
    ("41b", "T-hpid", "The command layup spec check --facts",
     "`layup spec check --facts SPANS --psb FILE` exits 1 on a span table with a byte in two spans, 0 on one that covers the file, and 2 on a missing file.",
     "`internal/cli`", "REQ-012", "F-0003#51",
     "e2e (the binary)", "small", "300", "41a",
     "1: the command and its exit codes"),
    ("42a", "T-gmey", "The session of a phase",
     "A session of a phase on a scripted harness commits its rows and its telemetry row in one records commit, from a directory that holds `input/`.",
     "`internal/session`, `internal/run`", "REQ-001, REQ-012, NFR-001", "F-0003#41, F-0003#51",
     "integration (scripted harness programs, the real `git`)", "large", "700", "40",
     "1: the session of a phase"),
    ("42b", "T-ybd5", "The budget of milestone 0",
     "No session of Intake starts when the known money since Start plus the new session's cap passes the money cap of milestone 0.",
     "`internal/run`", "REQ-011, NFR-001", "F-0003#50",
     "unit", "small", "300", "42a",
     "1: the cap of milestone 0"),
    ("43a", "T-qk66", "The step rule-gaps",
     "A restart after Start commits the table of `layup psb check` as `intake/psb.tsv`, byte for byte, in place of the row `phase`.",
     "`internal/run`", "REQ-001", "F-0003#41",
     "integration; the Go schema of `run-steps` with `rule-gaps`", "small", "300", "—",
     "1: the step `rule-gaps` and the rule \"done before\""),
    ("43b", "T-1u2v", "The review of meaning",
     "A quote that occurs twice refuses the session's file; one that occurs once is committed with its offsets.",
     "`internal/run`", "REQ-001", "F-0003#41",
     "unit (the quote check); integration (the step `meaning`, a scripted harness)", "large", "500", "42a, 43a",
     "1: the review of meaning and its quote check"),
    ("44", "T-hiv9", "The steps numbering and spans",
     "After a failed span check no batch is posted, and a restart runs `numbering` again.",
     "`internal/run`", "REQ-012", "F-0003#51",
     "integration (scripted harness programs)", "large", "500", "41a, 42a, 43a",
     "2: the two steps and their records; the rerun after a failed check"),
    ("45a", "T-2et0", "The draft",
     "A draft whose `covers` names no fact is refused, and one that covers each need and constraint is committed.",
     "`internal/run`", "REQ-012", "F-0003#51",
     "unit (the check); integration (the step `draft`)", "small", "400", "44",
     "1: the draft and its check"),
    ("45b", "T-0mey", "The completeness review",
     "With two harnesses the review runs on the one that did not number or draft; with none, the step is `fail` and no batch is posted.",
     "`internal/run`", "REQ-012, REQ-013", "F-0003#51, F-0003#52",
     "integration (two scripted harnesses)", "small", "400", "45a",
     "1: the completeness review on another harness"),
    ("46", "T-8r3u", "The batch",
     "One comment on the Intake issue holds the merged questions, with no repeat of a quote and kind, in two blocks.",
     "`internal/run`", "REQ-001", "F-0003#41",
     "unit (the merge and the IDs); integration (the comment, the `httptest` forge)", "large", "600", "43b, 45b",
     "2: the merge and the IDs; the comment, its two blocks and the intent form"),
    ("47a", "T-hxb9", "The answers",
     "The step waits until each block has an approver's comment with an answer line, then writes the answers and `budget.tsv`.",
     "`internal/run`", "REQ-001, NFR-003", "F-0003#41",
     "integration (scripted answer comments, the `httptest` forge)", "large", "500", "46",
     "2: the wait and the answer lines; the parse of the intent form and `budget.tsv`"),
    ("47b", "T-vwu3", "The review of the answers and the follow-up",
     "One follow-up comment asks only the missing, unparsed and open questions.",
     "`internal/run`", "REQ-001", "F-0003#41",
     "integration (a scripted harness, the `httptest` forge)", "large", "500", "47a",
     "1: the review of the answers and the follow-up comment"),
    ("47c", "T-tmm2", "The answers to the follow-up",
     "The answers to the follow-up are recorded, an intent value among them replaces its value in `budget.tsv`, and Intake ends.",
     "`internal/run`", "REQ-001", "F-0003#41",
     "integration", "small", "350", "47b",
     "1: the answers to the follow-up"),
    ("48", "T-mo9w", "The e2e of the steps of Intake",
     "In CI with no secret, a restart runs Intake to the batch, the next reads a scripted answer and posts one follow-up, and the next ends Intake.",
     "`cmd/layup` (the e2e test)", "REQ-001, REQ-012, NFR-001", "F-0003#41, F-0003#51",
     "e2e (the binary, a local fake forge, scripted harness programs)", "large", "500", "41b, 42b, 47c",
     "1: the steps of Intake, end to end"),
    ("49a", "T-8byd", "The review of the release of M2c",
     "The release review of `M2c` is a recorded review of the code for `REQ-015` and `REQ-017`, by a reviewer of a model that wrote none of it.",
     "—", "REQ-015, REQ-017", "F-0003#51 (the review keeps the out-of-scope facts F-0003#53 and F-0003#55 out)",
     "the review, with `runs/T-efmy/release-check.sh` on the release commit, adapted where the code of `M2c` needs it; its evidence under `runs/`", "small", "300", "48",
     "2: the deterministic check of `release-check.sh`; the review by a reviewer of another model"),
    ("49b", "T-ekh5", "The demo of M2c (uat)",
     "On a pilot problem statement, one Intake comment holds the rule gaps and the gaps of meaning, and `layup spec check --facts` passes on the numbered spans.",
     "—", "REQ-001, REQ-012", "F-0003#41, F-0003#51",
     "uat, with its evidence under `runs/`", "small", "300", "49a",
     "1: the demo"),
]

CAPS = {"41a": "1 (2 if its plan review reads the package rules as a gate)"}

OPERATOR_INPUTS = """## The Operator's inputs (R6: severity normal; needed when row 49a merges)

- The pilot problem statement of the demo, chosen by the idea owner (`PRD-0001` §11 item 6; decision 7 of #206).
- A target started from the host, with the App installed on it (a register that gains a harness after Start is allowed since #199).
- At least two admitted harnesses, each with its credential file, mode 0600, with a pair of the roles `systems-architect`, `product-owner` and `qa-engineer` at the tier `reasoning` in `registers/routing.tsv` ([`intake.md`]({spec}intake.md#the-sessions-of-intake)).
- The cap of milestone 0, `--intake-cap`, of the Start.
- The idea owner's and the Operator's answers to the batch, as approvers.
""".format(spec=SPEC)

REVIEW_NOTE = """## The release

The release of `M2c` is the code of `main` that the demo of row 49b runs (the non-test Go files of `cmd/` and `internal/`, `go.mod`, and the files that the binary embeds), at the merge of the last build row; the review reads the whole release ([`intake.md`]({spec}intake.md#req-015-and-req-017--the-review-of-the-release-of-m2c)).
""".format(spec=SPEC)


COMMON = """- [ ] The demo above holds, by the tests of the Tests cell, written red first.
- [ ] Each part above is built as the specification says; a difference is a defect of the specification, fixed in the same pull request (R10).
- [ ] `go build ./...`, `go vet ./...`, `go test ./...` and `go test -tags=integration ./...` pass; `adr-lint`, `prd-lint`, `link-lint`, `setup-check` and `run-discipline-tests` exit 0.
"""
REVIEW = """- [ ] The review is recorded under `runs/`, by a reviewer of a model that wrote none of the release, with the result of `release-check.sh` at the reviewed commit.
- [ ] Each finding of the review is fixed, or opens an issue with its reason.
"""
UAT = """- [ ] The demo above holds on a pilot problem statement that the idea owner chose, on the Operator's harnesses, with its evidence under `runs/` (the Intake comment, `spec/spans.tsv`, `spec/check.tsv` and the output of `layup spec check --facts`).
- [ ] Each write of LAYUP on the target shows the App's bot as its author.
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
    if row == "49a":
        crit = REVIEW + COMMON.split("\n", 2)[2]
    elif row == "49b":
        crit = UAT
    items = item(row)
    b = f"""## Goal

Task `{task}`: row {row} of the [implementation plan](https://github.com/pharzam/layup/blob/main/docs/plan/README.md#the-tasks-of-m2c), {title}. Milestone `M2c`; the specification is [`docs/spec/intake.md`]({SPEC}intake.md), [`records.md`]({SPEC}records.md#nfr-001--the-records-of-intake) and [`packages.md`]({SPEC}packages.md#the-table-of-m2c). Fact: {fact}. Opened by `T-2r8u` (#211). Refs #206, Refs #29.

**The demo** (R11): {demo}

**The goal count** (R11), final by the Operator's decision O-196 of #211: {classes}.

## Scope

- **The parts of the specification:** {"; ".join(p[row])}.
- **Packages:** {pkgs}. **Requirements:** {reqs}. **Items:** {items}.
- **Tests:** {tests}.
- **Size:** {size}, about {lines} lines (an estimate; the plan review sets the budget). **After:** {after}. **Expected cap:** {cap}.

## Acceptance criteria

{crit}- [ ] Docs updated in the same PR (`docs/tests/traceability.md`, the Test cells of `PRD-0001` §12).
"""
    if row == "49a":
        b += "\n" + REVIEW_NOTE
    if row == "49b":
        b += "\n" + OPERATOR_INPUTS
    return b + "\n*Written by Claude Opus 5.5, the author of T-2r8u (#211).*\n"


def post():
    p = parts()
    token = subprocess.run([os.path.expanduser("~/.config/layup-agent/app-token.sh")],
                           capture_output=True, text=True, check=True).stdout.strip()
    env = dict(os.environ)
    env["GH_TOKEN"] = token
    env["GH_PROMPT_DISABLED"] = "1"
    lines = ["row\ttask\tissue"]
    for r in ROWS:
        with tempfile.NamedTemporaryFile("w", suffix=".md", delete=False) as f:
            f.write(body(r, p))
        title = f"{r[1]}: {r[2]} (row {r[0]}, M2c)"
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
