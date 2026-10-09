#!/usr/bin/env python3
"""The issues and the plan table of the build tasks of M2b (task T-fdaq, #152).

Copied from runs/T-zwke/gen-issues.py and changed for M2b: one list, ROWS, holds
the 18 rows of O-186 and O-187 (#152); the parts of the specification of each
row come from runs/T-fdaq/parts.tsv, so the issue, the table and the check
parts.sh read one list. With --post, the script opens one issue per row through
the layup-agent App (the token from ~/.config/layup-agent/app-token.sh) and
writes runs/T-fdaq/issues.tsv (row, task, issue). With --body ROW, it prints the
body of one row's issue; with --edit ROW..., it writes the body of each named
row's issue again (the fixes of round 1). With --table, it prints the table of
docs/plan/README.md from ROWS, parts.tsv and issues.tsv. A one-off evidence
command, not a check of the gate.

    python3 -I runs/T-fdaq/gen-issues.py --body 28
    python3 -I runs/T-fdaq/gen-issues.py --post
    python3 -I runs/T-fdaq/gen-issues.py --table
"""
import os
import subprocess
import sys
import tempfile

HERE = os.path.dirname(os.path.abspath(__file__))
ISSUES = os.path.join(HERE, "issues.tsv")
PARTS = os.path.join(HERE, "parts.tsv")
SPEC = "https://github.com/pharzam/layup/blob/main/docs/spec/"
ITEM = "`later-p2-sessions-ledger`"

# row, task, title, demo, packages, requirements, fact, tests, size, lines, after, classes
ROWS = [
    ("28", "T-3py1", "The records of a session",
     "The Go schemas of `sessions`, `harnesses`, `routing`, `events` and `result` equal their blocks.",
     "`internal/records`", "NFR-001, REQ-005, REQ-013", "F-0003#52, F-0003#45",
     "unit; the block test (the five names to `built`)", "large", "600", "—",
     "2: the five Go schemas equal their blocks; the row rules in words"),
    ("29", "T-y10b", "The package rules of M2b",
     "`TestPackageRules` refuses a package of the engine checks that depends on `internal/session`, with \"The table of M2b\" read.",
     "`cmd/layup` (the checker, its unit tests, a fixture)", "NFR-005, NFR-007", "F-0003#52",
     "unit (the checker); integration (`TestPackageRules`)", "large", "500", "—",
     "2: the rule of the engine checks; the form of \"Starts a program\" for the program of a register row"),
    ("30a", "T-ysph", "The host registers of M2b",
     "A harness register row whose `vars` names `GH_TOKEN` is refused with its line, by `internal/route`.",
     "`internal/route`", "REQ-013, NFR-001", "F-0003#52",
     "unit; the block test (`models`, `routing-register` to `built`; `harness-register` with its new columns and the fixtures that write it); `tsv.Compare`",
     "large", "600", "—",
     "2: the readers refuse a register that breaks its block or its row rules; the checks of the credential file"),
    ("30b", "T-cht1", "Admission and the routing order",
     "With a stand-in harnesses record, the pair of a session is the first pair of its role's list whose harness passed its probe at the version read.",
     "`internal/route`", "REQ-013", "F-0003#52",
     "unit", "small", "350", "28, 30a",
     "2: admission (the probe at the version, `use`); the order of the routing register and the refusal `pair`"),
    ("31", "T-z5dj", "The calls of M2b of internal/git",
     "`FetchSession` fetches a session's head into another clone while a session configuration that holds each key of git's documentation that starts a program runs none.",
     "`internal/git`", "REQ-003, NFR-001", "F-0003#43",
     "unit; integration (real `git`, the hostile session configuration)", "large", "600", "—",
     "3: `CloneLocal`; `FetchSession`; the three reads `IsAncestor`, `DiffFile`, `DiffBinary`"),
    ("32", "T-m1dx", "The rule-path check before a push",
     "`internal/rules` refuses a change of a rule path but passes added lines in §2 of `docs/guardrails.md`.",
     "`internal/rules`", "REQ-003", "F-0003#43",
     "unit", "small", "350", "29",
     "2: the match of a path (workflows, the rule paths, `.sh`); the exception of `docs/guardrails.md`"),
    ("33a", "T-vxdg", "The directory and the environment of a session",
     "In a real temporary tree, a session directory is a `--no-local` clone with no remote on `task/<task>/<attempt>`, whose environment holds the named list only.",
     "`internal/session`", "REQ-013, REQ-003", "F-0003#52",
     "unit (the environment and the credential); integration (the directory, real `git`)", "large", "500", "29, 30a, 31",
     "2: the directory; the environment and the credential"),
    ("33b", "T-6sbe", "The checks before the start of a session",
     "Before any process starts, a `CLAUDE.md` above the session directory refuses the start with `rules`.",
     "`internal/session`", "REQ-013, REQ-003", "F-0003#52",
     "unit (the context, the prompt size, the head); integration (the rule files in a real temporary tree)", "small", "400", "33a",
     "3: the rule files; the context and prompt-size checks; the head read from the files of `repo/.git`"),
    ("34a", "T-5pxd", "The process of a session and its stop",
     "With a fake harness program, a session that runs past `wall` is stopped by `SIGINT`, `SIGTERM` and `SIGKILL`.",
     "`internal/session`", "REQ-013", "F-0003#52",
     "integration (a fake harness: the version check, the input, the stop at `wall`, the output cap)", "large", "600", "30a, 33a",
     "2: the start (the words, the placeholders, the prompt, the process group); the limits and the stop"),
    ("34b", "T-bpxg", "The end of a session",
     "With a fake harness program, each way a session ends gets its class: `done`, `start`, `crash`, `no-result`, `wall` or `output`.",
     "`internal/session`", "REQ-013, REQ-005, REQ-011", "F-0003#52, F-0003#50",
     "unit (the usage report on two recorded `result` events); integration (each class of the end); the block test (`probe-result` to `built`)",
     "large", "600", "28, 34a",
     "2: the end (the classes, the result file); the usage report"),
    ("35", "T-4c3q", "The writer of the telemetry record",
     "One row per session passes `CheckTelemetry`, with its money `reported`, `computed` or `unknown` by the billing, the prices and the models.",
     "`internal/ledger`", "REQ-011", "F-0003#50",
     "unit", "small", "350", "29",
     "1: one row per session from the usage report and the prices"),
    ("36a", "T-d8t9", "The start of a task session and its refusals",
     "Against the `httptest` forge and a local bare repository, a task session's start row is pushed before its fake harness starts.",
     "`internal/run`", "REQ-013, NFR-001", "F-0003#52",
     "unit (each refused start; the attempt); integration (a task session with a fake harness, the `httptest` forge, a local bare repository)",
     "large", "600", "28, 30b, 33b, 34a",
     "2: the attempt and the start row before the process; a refused start"),
    ("36b", "T-fsjp", "The result and the end of a task session",
     "A task session's result file is committed byte for byte as `tasks/<task>/results/<session>.tsv`, in the records commit that holds its telemetry row.",
     "`internal/run`", "REQ-005, REQ-011, NFR-001", "F-0003#45, F-0003#50",
     "unit (the open attempt with a stand-in events table); integration (the result, the telemetry row, the comment, the removal of the directory)",
     "large", "500", "34b, 35, 36a",
     "2: the result (byte for byte, the artifact, the open attempt); the telemetry row, the comment and the removal at the end"),
    ("37a", "T-z027", "The checks before a push",
     "A session head that changes a rule path is refused before any push, with its diff as a payload.",
     "`internal/run`", "REQ-003", "F-0003#43",
     "integration (real `git`, a local bare repository)", "large", "500", "31, 32, 33b, 36b",
     "2: the fetch by SHA and the base; the refusal of a workflow or rule path, with its payload"),
    ("37b", "T-e3sy", "The push and the bind",
     "A clean session head is bound only after the forge accepts its push.",
     "`internal/run`", "REQ-003, NFR-001", "F-0003#43",
     "integration (real `git`, a local bare repository, the `httptest` forge)", "small", "300", "37a",
     "1: the push and the bind"),
    ("38", "T-nxe4", "The step probe",
     "In CI with no secret, a restart's step `probe` writes the rows of a scripted harness's probe, with the same bytes on a repeat.",
     "`internal/run`", "REQ-013, NFR-001", "F-0003#52",
     "unit (each reason of a failed probe; the copy of the routing register); e2e (the binary, a local fake forge, a scripted harness); the Go schema of `run-steps` with `probe`",
     "large", "800", "36b",
     "2: the probe and its pass rules; the step (the skip, the counts, the copy of the routing register, the comment)"),
    ("39a", "T-4tjy", "The review of the release of M2b",
     "The release review of `M2b` is a recorded review of the code for `REQ-015` and `REQ-017`, by a reviewer of a model that wrote none of it.",
     "—", "REQ-015, REQ-017", "F-0003#52 (the review keeps the out-of-scope facts F-0003#53 and F-0003#55 out)",
     "the review, with `runs/T-efmy/release-check.sh` adapted, its evidence under `runs/`", "small", "250", "37b, 38",
     "2: the deterministic check of `release-check.sh`, adapted; the review by a reviewer of another model"),
    ("39b", "T-x7cs", "The demo of M2b (uat)",
     "On real harnesses, after a probe of each of two registered harnesses, one developer session's commit lands on `task/<task>/1` with its telemetry row.",
     "—", "REQ-013, REQ-003, REQ-005, REQ-011, NFR-001", "F-0003#52, F-0003#50",
     "uat, with its evidence under `runs/`", "small", "250", "39a",
     "1: the demo"),
]

CAPS = {"29": "1 (2 if its plan review reads the test as a gate)"}

OPERATOR_INPUTS = """## The Operator's inputs (R6: severity normal; needed when row 39a merges)

- The test target of the demo of `M2a` (row 27, #132), started, with its records branch, and the App `layup-agent` installed on it.
- At least two harnesses on the host, each with its credential file, mode 0600 ([`session.md`]({spec}session.md#the-environment-and-the-harness-credential), K40).
- The host registers: `registers/harnesses.tsv` with the columns of `M2b`, `registers/models.tsv`, `registers/routing.tsv` and `prices.tsv` ([`records.md`]({spec}records.md#nfr-001--the-records-of-a-session)).
- The spend that the demo may use: the `cap` of each harness row.
""".format(spec=SPEC)

REVIEW_NOTE = """## #148

This row takes #148 (O-187 of #152): the review reads the whole release, the code that the demo of row 39b runs, so it covers the code of `M2a` too, and `release-check.sh` is adapted here (its check 4 allows the verbs of `Fetch`, `Push` and `FetchSession`, and a new list names each program that the code starts). Its pull request says `Closes #148`.
"""

COMMON = """- [ ] The demo above holds, by the tests of the Tests cell, written red first.
- [ ] Each part above is built as the specification says; a difference is a defect of the specification, fixed in the same pull request (R10).
- [ ] `go build ./...`, `go vet ./...`, `go test ./...` and `go test -tags=integration ./...` pass; `adr-lint`, `prd-lint`, `link-lint`, `setup-check` and `run-discipline-tests` exit 0.
"""
REVIEW = """- [ ] The review is recorded under `runs/`, by a reviewer of a model that wrote none of the release, with the result of the adapted `release-check.sh` at the reviewed commit.
- [ ] Each finding of the review is fixed, or opens an issue with its reason.
"""
UAT = """- [ ] The demo above holds on the Operator's harnesses, with its evidence under `runs/` (the probe rows, the session's commit on its task branch, the telemetry row).
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
    if row == "39a":
        crit = REVIEW + COMMON.split("\n", 2)[2]
    elif row == "39b":
        crit = UAT
    items = ITEM + (", `gov-operator-setup-o112`" if row == "39b" else "")
    b = f"""## Goal

Task `{task}`: row {row} of the [implementation plan](https://github.com/pharzam/layup/blob/main/docs/plan/README.md#the-tasks-of-m2b), {title}. Milestone `M2b`; the specification is [`docs/spec/session.md`]({SPEC}session.md), [`records.md`]({SPEC}records.md#nfr-001--the-records-of-a-session) and [`packages.md`]({SPEC}packages.md#the-table-of-m2b). Fact: {fact}. Opened by `T-fdaq` (#152). Refs #147, Refs #29.

**The demo** (R11): {demo}

**The goal count** (R11), final by the Operator's decision O-187 of #152: {classes}.

## Scope

- **The parts of the specification:** {"; ".join(p[row])}.
- **Packages:** {pkgs}. **Requirements:** {reqs}. **Items:** {items}.
- **Tests:** {tests}.
- **Size:** {size}, about {lines} lines (an estimate; the plan review sets the budget). **After:** {after}. **Expected cap:** {cap}.

## Acceptance criteria

{crit}- [ ] Docs updated in the same PR (`docs/tests/traceability.md`, the Test cells of `PRD-0001` §12).
"""
    if row == "39a":
        b += "\n" + REVIEW_NOTE
    if row == "39b":
        b += "\n" + OPERATOR_INPUTS
    return b + "\n*Written by Claude Opus 5.5, the author of T-fdaq (#152).*\n"


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
        title = f"{r[1]}: {r[2]} (row {r[0]}, M2b)"
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
        items = ITEM + (", `gov-operator-setup-o112`" if row == "39b" else "")
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
