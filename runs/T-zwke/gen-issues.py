#!/usr/bin/env python3
"""The issues and the plan table of the build tasks of M2a (task T-zwke, #124).

One list, ROWS, holds the seven rows. With --post, the script opens one issue per
row through the layup-agent App (GH_TOKEN from ~/.config/layup-agent/app-token.sh)
and writes runs/T-zwke/issues.tsv (row, task, issue). With --table, it prints the
table of docs/plan/README.md from ROWS and issues.tsv. A one-off evidence command,
not a check of the gate.

    python3 -I runs/T-zwke/gen-issues.py --post
    python3 -I runs/T-zwke/gen-issues.py --table
"""
import os
import subprocess
import sys
import tempfile

HERE = os.path.dirname(os.path.abspath(__file__))
ISSUES = os.path.join(HERE, "issues.tsv")
SPEC = "https://github.com/pharzam/layup/blob/main/docs/spec/"

# row, task, title, demo, spec parts, packages, requirements, items, tests, size, lines, after, cap
ROWS = [
    ("21", "T-8kqn", "The records of Start",
     "The Go schemas of `start`, `approvers`, `lease` and `copies` equal their blocks.",
     "`records.md`: the blocks `start`, `approvers`, `lease`, `copies`, and the row rules that they give in words; `packages.md`: the Job cell of `internal/records` in the table of phase 1",
     "`internal/records`", "NFR-001", "`later-p2-run-start`",
     "unit; the block test (the four names to `built`)", "large", "500", "—", "1"),
    ("22", "T-esfe", "The forge interface, the two host registers, and the package rules of M2a",
     "A forge register or a harness register that breaks its schema is refused with its line, by its owner package.",
     "`forge.md`: \"The six capabilities\" (the interface; the set of permissions of `M2a` and its check), the key-file checks of \"The App identity\"; "
     "`records.md`: the blocks `forge-register`, `harness-register`; `run.md`: rows 1 to 5 of \"Input states\"; "
     "`packages.md`: \"The table of M2a\", rule 5 by \"Connects\", read by `TestPackageRules` (the checker, its unit tests, the fixture `testdata/netimport`)",
     "`internal/forge`, `internal/route`, `cmd/layup` (the test)", "NFR-001, NFR-007", "`later-p2-run-start`",
     "unit; integration (the package rules); the block test (the two names to `built`)", "large", "600", "—",
     "1 (2 if its plan review reads the test as a gate)"),
    ("23", "T-xhgz", "The calls Fetch and Push of internal/git",
     "`Push` refuses a push to a local bare repository that is not a fast-forward.",
     "`packages.md`: \"The calls of `M2a`\"; `forge.md`: the token to `git` in one call's environment",
     "`internal/git`", "NFR-001", "`later-p2-run-start`",
     "unit; integration (real `git`; the fixed list of the environment)", "small", "300", "—", "1"),
    ("24", "T-6bq5", "The GitHub adapter",
     "Against a loopback server, the adapter plays each call of `M2a` with an installation token that it made from a JWT.",
     "`forge.md`: \"The App identity\" (the JWT and the installation token), \"The calls of M2a\", \"Forge errors\" (with a progress function that the adapter takes for the rate-limit wait), \"The test of the adapter\"",
     "`internal/forge/github`", "NFR-001, NFR-007", "`later-p2-run-start`",
     "unit; integration (`httptest`)", "large", "800", "22", "1"),
    ("25", "T-trej", "internal/run: Start and the restart",
     "Against the `httptest` forge and a local bare repository, Start makes the first records commit with the pin, the briefs, `approvers.tsv`, `start.tsv` and the lease row.",
     "`run.md`: \"The steps of `layup run --new`\", \"The README of a target that Start makes\", \"The restart\", \"The lease and fencing\", \"A human decision\", \"Copy before read\", "
     "\"The Intake and control issues\", rows 6 to 11 of \"Input states\", the table `run-steps`, NFR-006, NFR-002, REQ-002",
     "`internal/run`", "NFR-001, NFR-002, NFR-006, REQ-002", "`later-p2-run-start`",
     "unit (the lease with a stand-in clock and `git`; the two rules; the input states); integration (Start, and a restart with another LAYUP version); the block test (`run-steps` to `built`)",
     "large", "1200", "21, 22, 23, 24", "1"),
    ("26", "T-mqty", "The command layup run",
     "`layup run --new` then `layup run TARGET`, as the built binary against a local fake forge, give their tables, the same bytes on a repeat.",
     "`run.md`: \"The command\" (two rows, the selecting flag, the one flag that may be left out, each flag and its value rule, the printing of the table and the progress lines, the exit codes); `README.md`: the rules of the selecting flag and of a flag that may be left out; `packages.md`: the May import cell of `internal/cli`",
     "`internal/cli`", "NFR-001, NFR-002", "`later-p2-run-start`",
     "unit; e2e (the binary, no secret)", "large", "600", "25", "1"),
    ("27", "T-fnsr", "The demo of M2a (uat)",
     "On a real GitHub test target, the Operator reads with a plain `git clone` the records branch that Start wrote as the App's bot.",
     "`run.md`: the uat row of \"The acceptance tests of M2a\"", "—",
     "NFR-001, NFR-002, NFR-006, REQ-002", "`later-p2-run-start`, `gov-operator-setup-o112`",
     "uat, with its evidence under `runs/`", "small", "150", "26", "1"),
]

OPERATOR_INPUTS = """## The Operator's inputs (R6: severity normal; needed when row 26 merges)

- An empty repository of the Operator's account; public if the plan is `free`.
- The forge register row, `host:registers/forge.tsv`: `app_id`, `app_slug`, `key_file`, `watch_slug` (or `—`), `api`, `web` ([`records.md`]({spec}records.md#nfr-001--the-records-of-start)).
- A harness register, `host:registers/harnesses.tsv`; it may be empty.
- The App `layup-agent` installed on the repository, with contents and issues `write` and metadata `read`.
- The App's private key on the host, mode 0600 (O-112).
- The values of the command: the two logins, the plan, the intake cap, `lease.H` and `watch.T` ([`run.md`]({spec}run.md#the-command)).
""".format(spec=SPEC)


UAT_FROM = """- [ ] The demo above holds, by the tests of the Tests cell, written red first.
- [ ] Each part above is built as the specification says; a difference is a defect of the specification, fixed in the same pull request (R10).
- [ ] `go build ./...`, `go vet ./...`, `go test ./...` and `go test -tags=integration ./...` pass; `adr-lint`, `prd-lint`, `link-lint`, `setup-check` and `run-discipline-tests` exit 0.
"""
UAT_TO = """- [ ] The demo above holds on a real GitHub test target, with its evidence under `runs/` (the command, the table, the records branch read with a plain `git clone`).
- [ ] Each write of LAYUP on the target shows the App's bot as its author.
"""


def body(r):
    row, task, title, demo, parts, pkgs, reqs, items, tests, size, lines, after, cap = r
    b = f"""## Goal

Task `{task}`: row {row} of the [implementation plan](https://github.com/pharzam/layup/blob/main/docs/plan/README.md#the-tasks-of-m2a), {title}. Milestone `M2a`; the specification is [`docs/spec/run.md`]({SPEC}run.md), [`forge.md`]({SPEC}forge.md), [`records.md`]({SPEC}records.md#nfr-001--the-records-of-start) and [`packages.md`]({SPEC}packages.md#the-table-of-m2a). Fact: `F-0003#42`. Opened by `T-zwke` (#124). Refs #123, Refs #29.

**The demo** (R11): {demo}

## Scope

- **The parts of the specification:** {parts}.
- **Packages:** {pkgs}. **Requirements:** {reqs}. **Items:** {items}.
- **Tests:** {tests}.
- **Size:** {size}, about {lines} lines (an estimate; the plan review sets the budget). **After:** {after}. **Expected cap:** {cap}.

## Acceptance criteria

- [ ] The demo above holds, by the tests of the Tests cell, written red first.
- [ ] Each part above is built as the specification says; a difference is a defect of the specification, fixed in the same pull request (R10).
- [ ] `go build ./...`, `go vet ./...`, `go test ./...` and `go test -tags=integration ./...` pass; `adr-lint`, `prd-lint`, `link-lint`, `setup-check` and `run-discipline-tests` exit 0.
- [ ] Docs updated in the same PR (`docs/tests/traceability.md`, the Test cells of `PRD-0001` §12).
"""
    if row == "27":
        b = b.replace(UAT_FROM, UAT_TO) + "\n" + OPERATOR_INPUTS
    return b + "\n*Written by Claude Opus 5.5, the author of T-zwke (#124).*\n"


def post():
    token = subprocess.run([os.path.expanduser("~/.config/layup-agent/app-token.sh")],
                           capture_output=True, text=True, check=True).stdout.strip()
    env = dict(os.environ, GH_TOKEN=token, GH_PROMPT_DISABLED="1")
    lines = ["row\ttask\tissue"]
    for r in ROWS:
        with tempfile.NamedTemporaryFile("w", suffix=".md", delete=False) as f:
            f.write(body(r))
        title = f"{r[1]}: {r[2]} (row {r[0]}, M2a)"
        out = subprocess.run(["gh", "issue", "create", "--repo", "pharzam/layup", "--title", title,
                              "--body-file", f.name], capture_output=True, text=True, env=env,
                             stdin=subprocess.DEVNULL, timeout=60, check=True).stdout.strip()
        os.unlink(f.name)
        num = out.rsplit("/", 1)[1]
        print(r[0], r[1], out)
        lines.append(f"{r[0]}\t{r[1]}\t{num}")
    with open(ISSUES, "w") as f:
        f.write("\n".join(lines) + "\n")


def table():
    nums = {}
    with open(ISSUES) as f:
        for ln in f.read().splitlines()[1:]:
            row, task, num = ln.split("\t")
            nums[row] = num
    print("| # | Task ID | Issue | Task | Demo | Items | Requirements | Fact | Tests | Size | Lines | After | Cap |")
    print("| - | ------- | ----- | ---- | ---- | ----- | ------------ | ---- | ----- | ---- | ----- | ----- | --- |")
    for r in ROWS:
        row, task, title, demo, parts, pkgs, reqs, items, tests, size, lines, after, cap = r
        n = nums[row]
        print(f"| {row} | `{task}` | [#{n}](https://github.com/pharzam/layup/issues/{n}) | {title}: {parts} | {demo} | {items} | {reqs} | F-0003#42 | {tests} | {size} | {lines} | {after} | {cap} |")


if __name__ == "__main__":
    if sys.argv[1:] == ["--post"]:
        post()
    elif sys.argv[1:] == ["--table"]:
        table()
    else:
        sys.exit(__doc__)
