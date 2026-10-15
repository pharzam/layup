#!/usr/bin/env python3
"""The issues and the plan table of the build tasks of M2e (task T-57vx, #346).

Copied from runs/T-bng8/gen-issues.py and changed for M2e: one list, ROWS, holds
the 24 rows of the plan; the parts of the specification of each row come from
runs/T-57vx/parts.tsv, so the issue, the table and the check parts.sh read one
list. With --post, the script opens one issue per row through the layup-agent App
(the token from ~/.config/layup-agent/app-token.sh) and writes
runs/T-57vx/issues.tsv (row, task, issue). With --body ROW, it prints the body of
one row's issue; with --edit ROW..., it writes the body of each named row's issue
again. With --table, it prints the table of docs/plan/README.md from ROWS,
parts.tsv and issues.tsv. A one-off evidence command, not a check of the gate.

    python3 -I runs/T-57vx/gen-issues.py --body 61
    python3 -I runs/T-57vx/gen-issues.py --post
    python3 -I runs/T-57vx/gen-issues.py --table
"""
import os
import subprocess
import sys
import tempfile

HERE = os.path.dirname(os.path.abspath(__file__))
ISSUES = os.path.join(HERE, "issues.tsv")
PARTS = os.path.join(HERE, "parts.tsv")
SPEC = "https://github.com/pharzam/layup/blob/main/docs/spec/"
ITEM = "`later-p2-phase-loop-handoffs`"
ITEMS = {"61": "`later-p2-phase-loop-handoffs`, `gov-issue-68`", "64": "`later-gate-status`", "66d": "`later-p2-phase-loop-handoffs`, `gov-issue-68`",
         "66f": "`later-gate-status`, `later-p2-phase-loop-handoffs`", "66g": "`later-gate-status`, `later-p2-phase-loop-handoffs`",
         "65a": "`later-p2-phase-loop-handoffs` (F-22)", "65b": "`later-p2-phase-loop-handoffs` (F-22)"}


def item(row):
    return ITEMS.get(row, ITEM) + (", `gov-operator-setup-o112`" if row == "68b" else "")


# row, task, title, demo, packages, requirements, fact, tests, size, lines, after, classes
ROWS = [
    ('60',
     'T-fn5s',
     "The kinds of the task loop's records",
     'The Go values of the blocks `result` and `events` hold the new kinds (`question`, `decision`, `lesson`, `comment`; `handoff`, `merged`), and each reader refuses a row that breaks their rules.',
     '`internal/records`', 'NFR-001, REQ-005', 'F-0001#1, F-0003#59',
     'unit; the block test (the names stay in `built`)', 'small', '200', '—',
     '1: the Go values of two built blocks, changed in words and code'),
    ('61',
     'T-h17s',
     'The handoff',
     'The handoff check accepts a handoff only when its schema, each artifact with its hash (a comment by the three equal hashes) and each required item hold, refuses a pair with no row, and the thirteen default rows of the transition table are valid and a replacement is checked.',
     '`internal/handoff` (new)', 'REQ-005, NFR-004', 'F-0003#45, F-0003#59',
     'unit; the block test (`transitions` to `built`); `TestPackageRules` reads the table of `M2e`', 'large', '700', '60',
     '2: the table, its block and its replacement check; the check of a handoff'),
    ('62a',
     'T-0pt8',
     'The Shape filter',
     'The filter drops an option with a failed row, a missing row or a pass with no evidence, chooses one survivor, and lists the constraints that apply from the records, never from the session.',
     '`internal/spec`', 'REQ-012', 'F-0003#51',
     'unit; the block test (`options` to `built`)', 'small', '300', '—',
     '1: the filter and its block'),
    ('62b',
     'T-3q0b',
     'The checks of the bet and the rendering',
     "The brief check refuses a brief that lacks a heading, the line limit, a link or the appetite that fits `architecture.md` §12; the answer's lines are read; the rendering of the confirmed inventory and of the rendered commit is code's, and the quote's fact is checked.",
     '`internal/spec`', 'REQ-012', 'F-0003#51, F-0003#54',
     'unit; the block test (`inventory`, `inventory-index` to `built`)', 'large', '600', '—',
     '2: the checks of the brief and the answer; the rendering and the versions of the inventory'),
    ('62c',
     'T-4wyg',
     'The full spec check',
     'The full check gives each of its ten rules a row, before the first bet judges the numbered facts only (`classes` and `inventory` are `clear`), and accepts the head of a specification that waits for its bet only when each `prd-lint` failure line names an empty MoSCoW or Phase cell.',
     '`internal/spec` (starts `sh` for the pinned `prd-lint`)', 'REQ-012, NFR-004', 'F-0003#51, F-0003#62',
     'unit; integration (the pinned `prd-lint`, the real `sh`); the block test (`task-spec-check` to `built`)', 'large', '600', '62b',
     '2: the ten rules and the mode before the bet; the expected `prd-lint` failure'),
    ('63a',
     'T-9et9',
     'The basic layup/rules',
     "`layup/rules` fails a pull request that changes a path of the base's rule-path register unless it is a rendered record equal to code's rendering of the named copy.",
     '`internal/rules`', 'REQ-003, NFR-004', 'F-0003#43',
     'unit; integration (the real `git`)', 'small', '300', '62b',
     '1: the basic rule'),
    ('63b',
     'T-svn3',
     "The record kinds of a step's context",
     'The step table names, for each step, the record kinds of its prompt file, in the order of `loop.md`.',
     '`internal/session`', 'REQ-005', 'F-0003#45',
     'unit', 'small', '150', '—',
     "1: the step table's record kinds"),
    ('64',
     'T-cr8j',
     'The forge calls of M2e',
     'Against the `httptest` forge, `OpenDraft`, `MarkReady`, `SetStatus`, `Statuses`, `PullRequest` and `Merge` behave as `forge.md` says (a moved head is refused, a status is bound to its SHA), the stand-in gives the six, and a missing permission of the three fails the step `forge`.',
     '`internal/forge`, `internal/forge/github`', 'REQ-007, NFR-004, NFR-001', 'F-0003#47, F-0003#58',
     'integration (the `httptest` forge); unit (the stand-in)', 'large', '600', '—',
     '2: the six calls and the stand-in; the permission check'),
    ('65a',
     'T-hz4m',
     'The gate jobs read the base',
     "The catalog's workflow makes each native gate job take the manifest, the `config` paths and `gates.sh` from the base of the pull request, and from its own head for a `batch/<point>` branch.",
     '`internal/catalog`', 'REQ-007', 'F-0003#47',
     'unit; integration (the fixtures and the job script)', 'small', '300', '—',
     '1: the workflow of the catalog'),
    ('65b',
     'T-jew3',
     'The code-owner rule',
     "The setup writes `CODEOWNERS` for the four path classes and the rule `pull_request` of the default branch's ruleset sets `require_code_owner_review` with no required approval.",
     '`internal/setup`', 'REQ-007, REQ-003', 'F-0003#47, F-0003#43',
     'unit (the golden ruleset file, the `CODEOWNERS` text); the refusal of a pull request that weakens its own gate is shown on a real repository in the demo (row 68b)', 'small', '300', '65a',
     '2: the `CODEOWNERS` file; the code-owner rule of the ruleset'),
    ('66a',
     'T-cak3',
     'The task register',
     '`layup run` writes `tasks.tsv` and `transitions.tsv` from the default when the loop first runs and refuses a task whose predecessor is not a task, a cycle, or no requirement.',
     '`internal/run`', 'REQ-005, NFR-001', 'F-0003#45',
     'unit; integration (the `httptest` forge); the block test (`tasks` to `built`)', 'small', '400', '60, 61',
     '1: the register and the first writes'),
    ('66b',
     'T-dx1h',
     'The parameters',
     '`layup run` writes `parameters.tsv` with its defaults, copies a `set <name> <value> because <reason>` and applies it at the next step boundary only for an allowed role, a known name and a value inside its bound, and replaces the transition table by `set transitions`.',
     '`internal/run`', 'REQ-012, NFR-001', 'F-0003#51',
     'unit; the block test (`parameters` to `built`)', 'small', '400', '66a',
     '2: the first writer and the `set` command; the replacement of the transition table'),
    ('66c',
     'T-v2sv',
     'The Retrospective',
     'The Retrospective step writes its event and no session starts, and the loop goes to the next Bet.',
     '`internal/run`', 'REQ-012', 'F-0003#51',
     'unit', 'small', '100', '66a',
     '1: the step with no session'),
    ('66d',
     'T-q5k4',
     'The handoff in the loop',
     "The plan, the plan review and the review record are posted from the bytes of the session's file, copied with their hash and kept; the `handoff` event holds the comment ID; the end of a session checks a `comment` row against the session\'s `result/` file; a question goes to its owner role and an answer is accepted by its attempt; the items that need a program (`spec-sections`, `rendering-equal`) are computed by the caller and passed to the handoff check.",
     '`internal/run`', 'REQ-005', 'F-0003#45, F-0003#59',
     'unit; integration (the `httptest` forge)', 'large', '600', '61, 62b, 64, 66a',
     '2: the comment artifact; the `handoff` event and the route of a question'),
    ('66e',
     'T-6zj3',
     'The attempts, the independence and the frozen tests',
     'A gate failure starts a new attempt with its findings; the authors of a plan review and of a verification come from the ledger; the frozen test list freezes at the first valid handoff and a changed source is a finding; the tier and the prompt file are built from the plan and the records.',
     '`internal/run`', 'REQ-005, NFR-004', 'F-0003#45, F-0003#58',
     'unit; integration (the real `git`); the block test (`task-tests` to `built`)', 'large', '700', '61, 63b, 66d',
     '3: the attempts and the independence; the frozen test list; the tier and the prompt file'),
    ('66f',
     'T-q0kr',
     'The four statuses',
     '`layup run` posts `layup/gates`, `layup/spec`, `layup/rules` and `layup/verify` on each head, once for an unchanged state, each from its own records, with `not-active` a `failure`.',
     '`internal/run`', 'REQ-007, REQ-003, NFR-004', 'F-0003#47, F-0003#58',
     'unit; integration (the `httptest` forge); the block test (`task-gates` to `built`)', 'large', '700', '62c, 63a, 64, 66e',
     '2: the four statuses and their posting; the table of each head'),
    ('66g',
     'T-xvvd',
     'The pull request and the merge',
     'A draft pull request is opened with its body linking the issue and the records commit; it is marked ready only when the four statuses and each native job are `success` by a read; the close-out commit is checked; the merge is at the head SHA; a `merged` event is written and the tree check accepts a merge that is not yet recorded.',
     '`internal/run`', 'REQ-007, NFR-004, NFR-001', 'F-0003#58',
     'integration (the `httptest` forge, the real `git`)', 'large', '700', '64, 66f',
     '3: the pull request and the ready rule; the close-out and the merge; the tree rule'),
    ('66h',
     'T-3f94',
     'The step shape',
     '`layup run` runs the Shape tasks through the loop, writes the options, asks the questions left open (the intent answer writes `budget.tsv`) and stops with no survivor; `run-steps` gains `shape`.',
     '`internal/run`', 'REQ-012', 'F-0003#51',
     'unit (the rewrite and the stop); integration (the `httptest` forge, scripted harnesses)', 'large', '600', '62a, 66b, 66g',
     '2: the step shape; the questions left open and the band'),
    ('66i',
     'T-8rm6',
     'The approvers and the revision of the draft',
     "The success criteria become measure rows, an approver is tied to its point, and the revision session runs with the quote's-fact check.",
     '`internal/run`', 'REQ-012', 'F-0003#51, F-0003#54',
     'unit; integration (scripted harnesses)', 'small', '400', '62b, 66d',
     '2: the approvers and the success criteria; the revision session'),
    ('66j',
     'T-2vva',
     'The step bet',
     '`layup run` posts the brief, takes the answer of the approver, writes the inventory version, pushes the rendered commit with its carry-over of `layup/verify`, opens the architecture tasks, renders a later bet as a rendered task, and writes `bets.tsv` last; `run-steps` gains `bet`.',
     '`internal/run`', 'REQ-012, REQ-003', 'F-0003#51, F-0003#54',
     'unit (the option answer lines, the approver\'s role); integration (the `httptest` forge, the real `git`); the block test (`bets` to `built`)', 'large', '700', '62b, 66g, 66h, 66i',
     '2: the step bet; the rendered commit and the architecture tasks'),
    ('66k',
     'T-9nbf',
     'The step loop and the restart',
     "`layup run` resumes the loop from the records and the forge, binds a push with no result or refuses it as `push-lost`, posts a task session's comment on its issue, merges the Shape tasks in the merge order after a bet, and redoes a task whose base moved (the events `closed` and `rebased`); `run-steps` gains `loop`.",
     '`internal/run`', 'REQ-005, NFR-001', 'F-0003#45',
     'integration (the `httptest` forge, scripted harnesses)', 'large', '700', '66c, 66j',
     '2: the step loop and the restart; the redo of a task'),
    ('67',
     'T-3xvc',
     'The e2e of the specification task',
     'In CI with no secret, Shape, the first bet and the specification task run from a local fake forge and scripted harnesses to a merge with the four statuses `success`.',
     '`cmd/layup`', 'REQ-005, REQ-007, REQ-012, NFR-004', 'F-0003#45, F-0003#58',
     'e2e (the binary, a local fake forge, local bare repositories, scripted harness programs and approver comments)', 'large', '700', '66k, 65b',
     '1: the whole loop, end to end'),
    ('68a',
     'T-ym0k',
     'The review of the release of M2e',
     'A recorded code review of the release of `M2e` by a model that wrote none of it, with `release-check.sh` on the release commit, `holds` or its findings fixed or opened as issues.',
     '—', 'REQ-015, REQ-017', 'F-0003#51',
     'uat (a review)', 'small', '400', '67',
     '1: the review'),
    ('68b',
     'T-8s13',
     'The demo of M2e',
     'On a real GitHub test target that `M2d` set up, a specification task goes through the loop, Shape runs with its filter, the idea owner answers the brief with one comment, the rendered commit is added and the task merges at its head SHA with the four statuses `success`.',
     '—', 'REQ-005, REQ-007, REQ-012', 'F-0003#51, F-0003#58',
     'uat', 'small', '300', '68a',
     '1: the demo'),
]

CAPS = {}

EXTRA = {
    "65b": "The Operator's question (R6: severity normal; row 65b waits on it): where the setup writes `CODEOWNERS`, and how the Operator applies the code-owner rule of the ruleset to a target that `M2d` set up. It is the open gap \"where CODEOWNERS is written and how its ruleset is applied\" of `docs/setup/open-gaps.tsv`; ask it in this issue before the build, with the answer as the evidence.",
    "66g": "The forge calls to edit a pull request's body and to close a pull request or an issue are not yet in `forge.md`: issue #373 asks them; this row may settle them in its own pull request (R10).",
    "66k": "The forge call that closes a pull request or an issue (the redo of a task) is not yet in `forge.md`: issue #373; this row may settle it in its own pull request (R10).",
    "64": "The home of the stand-in forge (`forge.md` names `internal/standin`, which cannot hold it) is in issue #373; this row may settle it in its own pull request (R10).",
    "66c": "The record of the Retrospective's event is not yet named: the block `events` is the file of one task. Issue #372 asks it; this row may settle it in its own pull request (R10), and #372 then closes with it.",
}

OPERATOR_INPUTS = """## The Operator's inputs (R6: severity normal; needed when row 68a merges)

- A test target that `M2d` set up (the target of the demo of `M2d`, or a new one), public or on a plan that allows draft pull requests and enforces rulesets.
- The App's permissions on it: pull requests `write`, commit statuses `write` and checks `read` (checks `read` amends ADR-0014 decision 3, #336), and the Operator's apply of the code-owner rule (`CODEOWNERS` and the ruleset), in the SSH form.
- At least two admitted harnesses and the routing of the roles of the loop.
- The value of `brief.lines`, set on the idea owner's evidence, and the idea owner's account as an approver.
- The Operator's confirmation of the reading of `F-0003#58` (a pull request "reaches human review" when a review is requested or it is marked ready).
- The Operator's one-word acceptance of the scenario before the review of the demo ([`loop.md`]({spec}loop.md#the-demo-and-the-release-review)).
""".format(spec=SPEC)

REVIEW_NOTE = """## The release

The release of `M2e` is the code of `main` that the demo of row 68b runs (the non-test Go files of `cmd/` and `internal/`, `go.mod`, and the files that the binary embeds), at the merge of the last build row, 67; the review reads the whole release ([`loop.md`]({spec}loop.md#the-demo-and-the-release-review)).
""".format(spec=SPEC)

COMMON = """- [ ] The demo above holds, by the tests of the Tests cell, written red first.
- [ ] Each part above is built as the specification says; a difference is a defect of the specification, fixed in the same pull request (R10).
- [ ] `go build ./...`, `go vet ./...`, `go test ./...` and `go test -tags=integration ./...` pass; `adr-lint`, `prd-lint`, `link-lint`, `setup-check` and `run-discipline-tests` exit 0.
"""
REVIEW = """- [ ] The review is recorded under `runs/`, by a reviewer of a model that wrote none of the release, with the result of `release-check.sh` at the reviewed commit and a coverage table of the files read that a script compares with `git ls-files`.
- [ ] Each finding of the review is fixed, or opens an issue with its reason.
"""
UAT = """- [ ] The demo above holds on a test target, on the Operator's harnesses, with its evidence under `runs/`, and passes when each of the five rules of `loop.md` (The demo and the release review) holds.
- [ ] Each write of LAYUP on the target shows the App's bot as its author; the Operator's apply and the code-owner review are the Operator's.
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
    if row == "68a":
        crit = REVIEW + COMMON.split("\n", 1)[1]
    elif row == "68b":
        crit = UAT
    items = item(row)
    b = f"""## Goal

Task `{task}`: row {row} of the [implementation plan](https://github.com/pharzam/layup/blob/main/docs/plan/README.md#the-tasks-of-m2e), {title}. Milestone `M2e`; the specification is [`docs/spec/loop.md`]({SPEC}loop.md), [`records.md`]({SPEC}records.md#nfr-001--the-records-of-the-task-loop), [`forge.md`]({SPEC}forge.md#the-calls-of-m2e) and [`packages.md`]({SPEC}packages.md#the-table-of-m2e). Fact: {fact}. Opened by `T-57vx` (#346). Refs #333, Refs #29.

**The demo** (R11): {demo}

**The goal count** (R11), final by the Operator's decision O-214 of #346: {classes}.

## Scope

- **The parts of the specification:** {"; ".join(p[row])}.
- **Packages:** {pkgs}. **Requirements:** {reqs}. **Items:** {items}.
- **Tests:** {tests}.
- **Size:** {size}, about {lines} lines (an estimate; the plan review sets the budget). **After:** {after}. **Expected cap:** {cap}.

## Acceptance criteria

{crit}- [ ] Docs updated in the same PR (`docs/tests/traceability.md`, the Test cells of `PRD-0001` §12).
"""
    if row in EXTRA:
        b += "\n## Notes\n\n- " + EXTRA[row] + "\n"
    if row == "68a":
        b += "\n" + REVIEW_NOTE
    if row == "68b":
        b += "\n" + OPERATOR_INPUTS
    return b + "\n*Written by Claude Sonnet 5.5, the author of T-57vx (#346).*\n"

def dup(title, env, num=None):
    """The duplicate check of one child issue (R2), as the model of T-bng8 made it:
    the open and closed issues whose title holds the row's title. At an edit the
    issue exists, so num drops it from the found set."""
    q = f'repo:pharzam/layup in:title "{title}"'
    sel = f'[.items[] | select(.number != {num or 0}) | "#\\(.number) (\\(.title))"] | join("; ")'
    out = subprocess.run(["gh", "api", "-X", "GET", "search/issues", "-f", f"q={q}",
                          "--jq", sel],
                         capture_output=True, text=True, env=env, stdin=subprocess.DEVNULL,
                         timeout=60, check=True).stdout.strip()
    day = subprocess.run(["date", "-u", "+%Y-%m-%d"], capture_output=True, text=True).stdout.strip()
    found = out if out else "none"
    when = "before this issue opened" if num is None else "at this edit, this issue itself excluded"
    return (f"\n## Duplicate check (R2)\n\n- [x] Searched the open and closed issues of "
            f"`pharzam/layup` for the title \"{title}\" on {day}, {when}: "
            f"{found}. None of them is this task: each is a task of another "
            f"milestone or of the specification. Related: #333 (the specification), #346 (the slicing).\n")


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
        title = f"{r[1]}: {r[2]} (row {r[0]}, M2e)"
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
        b = body(r, p)
        sig = b.rindex("\n*Written by")
        b = b[:sig] + dup(r[2], env, nums[r[0]]) + b[sig:]  # as the model does
        with tempfile.NamedTemporaryFile("w", suffix=".md", delete=False) as f:
            f.write(b)
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
