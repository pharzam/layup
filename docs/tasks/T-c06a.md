# T-c06a — the Go entry of the stack catalog

Issue: [#91](https://github.com/pharzam/layup/issues/91), row 14 of the
[implementation plan](../plan/README.md), before row 9 by O-137 (#86); child of
`layup gate` (#33). Serves `F-0003#44`. Base `4c50e81` (the merge of row 10,
#111). Author: Claude Opus 5.5 on Claude Code. Evidence:
[`runs/T-c06a/`](../../runs/T-c06a/).

## Plan and plan review

The plan (R12), its review and the author's answer are comments on #91. The
reviewer order is the Operator's (Devin, then OpenCode, then Claude), copied to
#91 in the answer. Devin (its usage quota) gave no record; OpenCode (Grok 4.7,
effort `xhigh`) stopped with "Go usage limit exceeded" after ten minutes and
gave no record. The plan review (Claude Fable 5.1, effort `xhigh`, on the
Claude Code CLI, a fresh session) gave `approve-with-conditions`: Budget
maximum 2,400 lines added plus removed over 40 files against the base,
close-out inside; Cycle cap 1 (the author proposed 2; no check script, hook, CI
file or protection file of LAYUP changes). Its two conditions (three stale
sentences: §6 of the architecture twice and `NFR-002` item 2; the rule for the
activation as a rule of the Go entry) are applied, and its twelve notes too,
but note 7: by the Operator's standing instruction, a finding is fixed in its
pull request or becomes a known limit, never a follow-up issue, so the bare
`gofmt` form of LAYUP's own job `lint` is a known limit in `guardrails.md` §2
(`ci.yml` cannot change in this branch: K28, O-92).

## What was done

1. **The Go entry** (`internal/catalog/go/`): five kinds (K20: `static` and
   `test` active with the version and the documentation of `go`; `layout`,
   `boundary` and `contract` pending with `—`); `go.mod` on the `go` line of
   LAYUP's own; the workflow of one job per kind, whose id and name are the
   kind (K30); one job script with the rules of the run of `layup gate`; the
   file of the coverage floor, whose one line is a gap token, with its question
   in `gaps.tsv` (K24); the two fixtures, each only new files under
   `gatefixture/`.
2. **`internal/catalog`:** the rules of a pending and an active kind (note 1),
   the gap list (the block `catalog-gaps`), the gap token and `Gaps`, the
   marker rule, the module check of `Files`, `Has` of `gaps.tsv`, and
   `Embedded`, the root of the entries of the binary. The test entry follows
   the rules (the version and the documentation of `sh`).
3. **`internal/verify`:** check `sources` resolves a `catalog` ref in the
   entry of the binary.
4. **The tests:** the unit tests of the entry, of its workflow and of the
   rules; the integration tests of each fixture of each entry, of good code,
   and of the job script against `layup gate` on each line of the table of the
   run and on the input errors; the e2e test of the fixtures through the
   binary.
5. **The documents:** "The stack catalog" of `setup.md` (the blocks, the
   rules, the decided values of the Go entry), §6 of the architecture,
   `gate.md` (K30, a known limit), `records.md` (`NFR-002`), W-04 (K20), the
   known limit of #21 in `docs/setup/README.md` (O-138), a lesson in
   `guardrails.md` §2, the traceability rows (the planned row named by its
   tests) and the `PRD-0001` cells with a §13 row.

**The rejected alternatives:** the bare form `test -z "$(gofmt -l .)"` (a pass
on nothing when `gofmt` cannot run); the planned tools and commands of the
pending kinds at setup (values with no source); a job script in each job (five
copies of one rule); one job with a matrix (a static reader sees one job, and
`gate.md` asks one job per kind); a marker character in the entry (LAYUP's own
check `markers` would read it); fixtures that change a file of the baseline
(they stop applying when the baseline changes); a test of the fixtures on the
clean commit only (a command that fails on each tree would be `clear` there).

## Review round 1 and its fix (cycle 1)

Devin (its usage quota) and OpenCode ("Go usage limit exceeded") gave no
record. Round 1 (Claude Fable 5.1, effort `xhigh`, on the Claude Code CLI, on
`b00d5c9`; the record is on #91) gave `material`, with one finding and eleven
notes. The fix has its red run ([`test-runs.md`](../../runs/T-c06a/test-runs.md)):

1. **The job passed on a manifest that `layup gate` refuses.** With an empty
   `command` or `config` field, a `config` value that is not a path, a carriage
   return in a row, or no line feed at the end, the job passed where `layup
   gate` gives exit 2; with an empty command, it passed on nothing. The job
   script now refuses each input that `layup gate` refuses: each form that the
   reader of `internal/tsv` refuses (with `iconv` for UTF-8 and `tail` for the
   last line feed), and a base or a head that is not a commit (note 5). The
   parity test has these cases, and it runs the script with each `sh` of the
   host.

The notes: note 2 (a builtin as the tool) is fixed, as the tool of an active
kind is a program, found by its path; note 3 (the programs that the script
starts) is fixed in its comment and in `setup.md`; note 4 (a tab in a reason)
is fixed, as the reason writes a space; note 5 is fixed (above). Notes 6 to 12
confirm the change. One more defect, found by the author during the fix: with
`bash` 5.3, a pending case gave `clear` in 4 runs of 20. The author took the
trap on `EXIT` for its cause; that was wrong (below).

## The late run of round 2, skipped, and the rest of the fix of cycle 1

The first run of round 2 (Claude Fable 5.1, effort `xhigh`, on `b9e1b05`)
wrote its record at 15 min 9 s, after the fifteen minutes of Bootstrap mode
rule 4, so it is skipped and is not a round. Its text gave `not mergeable,
findings recorded`: a subshell of the script crashed under `bash` 5.3 (5 runs
of 100, measured by the author), and the script read each failure of its
scope check as `clear`, a pass with the command not run. A known defect is
not set aside, so the fix of cycle 1 holds it too: the scope check gives three
answers and a failure is `not-active`; each read checks its status; the tool
is found as `exec.LookPath` finds it (its notes 3 and 4); a last byte NUL is no
line feed (its note 2); the text of `setup.md` names a manifest with no row,
and the checkout rule apart (its note 5). The lesson of `guardrails.md` §2 is
now "A failed check that reads as a pass". Round 2 then ran on the new head,
with a fresh session.

## Review round 2

Devin (its usage quota) and OpenCode ("Go usage limit exceeded") gave no record
again. Round 2 (Claude Fable 5.1, effort `xhigh`, on the Claude Code CLI, a
fresh session with another prompt, on `7b14e95`, its record at 9 min 10 s; the
record is on #91) gave `nothing material in scope`, with five notes. It ran the
script with five shells on 27 manifests against the binary, 800 runs of a
pending and an active kind (each right), and failures of `awk` and `tr` (each
`not-active` or `fail`, never a pass). Notes 1 to 3 are applied as text at the
close-out, as known limits in `setup.md`: a `tr` that fails in the check of the
last line gives `fail` with no reason; a NUL byte in another row ends the line
for the `awk` of a host, so the job fails where `layup gate` runs; a relative
directory of `PATH` before the tool's absolute one. Notes 4 and 5 confirm the
change.

## Verdict

Delivered: the Go entry of the stack catalog, embedded in the binary: five
kinds (`static` and `test` active, with the version and the documentation of
`go`; `layout`, `boundary` and `contract` pending, K20); `go.mod`; the workflow
of one job per kind, whose id and name are the kind (K30); the job script with
the rules of the run of `layup gate`, which fails on each input that `layup
gate` refuses and never reads a failed check as a pass; the coverage floor as
a gap token with its question (K24); two fixtures; and LAYUP's CI run of each
fixture at the integration and the e2e levels, with no change of `ci.yml`.
Check `sources` resolves a `catalog` ref in the entry of the binary. O-137 put
this row before row 9, and O-138 closed #21.

The plan review (Claude Fable 5.1) gave `approve-with-conditions`, with two
conditions, applied; round 1 (`b00d5c9`) gave `material`, one finding, fixed
in cycle 1 with notes 2 to 5; a first run of round 2 on `b9e1b05` came after
fifteen minutes and is skipped, and the defect that its text reported is fixed
in the same cycle (`7b14e95`); round 2 (`7b14e95`) gave `nothing material in
scope`. The records are on #91. At `7b14e95`, `go build`, `go vet` with each
tag, `gofmt`, the three test levels, `go test -race` on `internal/catalog`,
`internal/verify`, `internal/work` and `internal/standin`, `run.sh` and the
discipline tests pass; at the head, all local checks pass, and
`review-record-lint` passes on the comments of #91 (2 rounds, cap 1). The diff
against `origin/main` is 1,785 lines over 32 files with the close-out,
inside the Budget maximum of 2,400 lines over 40 files.

Next: row 9 of the plan (`T-7s0y`, #86), which O-137 put after this row.

## Resource record

Recorded, not budgeted (ADR-0007). Times are 2026-10-02, UTC. A token count is
the `result` event of the Claude Code CLI (input, output, cache creation and
cache read tokens, and its cost) where that harness gave one; `not reported`
where the harness or the author's session gives none.

| Part | Expected tier | Model | Effort | Tokens | Elapsed |
| ---- | ------------- | ----- | ------ | ------ | ------- |
| The records of O-137 and O-138 on #86 and #21, and #21 closed | reasoning | Claude Opus 5.5 | max | not reported | 18:47 to 18:49 |
| The plan, with the two inventory items, the specification, the evidence URLs and the measurements of D2 and D9 | reasoning | Claude Opus 5.5 | max | not reported | 18:49 to 18:59 |
| The plan review, first and second harness: skipped (Devin's usage quota; OpenCode "Go usage limit exceeded" during its run) | reasoning | GPT-6 Sol on the Devin CLI; Grok 4.7 on the OpenCode CLI | `xhigh`; `xhigh` | not reported | 19:00:04 to 19:11:22 |
| The plan review | reasoning | Claude Fable 5.1 on the Claude Code CLI | `xhigh` | 1,069,741 (USD 5.17) | 8 min 43 s, 19:11:31 to 19:20:14 |
| Drafts of the job script, the workflow and the fixtures, while the plan review ran | execution | Claude Opus 5.5 | max | not reported | 19:02 to 19:11 |
| The answer to the plan review | reasoning | Claude Opus 5.5 | max | not reported | 19:20 to 19:23 |
| The tests first, the entry, the code, the integration and e2e tests, the mutations, the documents and the evidence; the freeze | execution | Claude Opus 5.5, a reasoning-tier model on an execution part | max | not reported | 19:23 to 19:43 |
| Review round 1, first and second harness: skipped (the same) | reasoning | GPT-6 Sol on the Devin CLI; Grok 4.7 on the OpenCode CLI | `xhigh`; `low` | not reported | 72 s, 19:43:38 to 19:44:50 |
| Review round 1 | reasoning | Claude Fable 5.1 on the Claude Code CLI | `xhigh` | 2,123,257 (USD 6.18) | 9 min 20 s, 19:44:56 to 19:54:16 |
| The fix of round 1, test first, and the freeze | execution | Claude Opus 5.5 | max | not reported | 19:54 to 20:04 |
| Review round 2, first try: Devin and OpenCode skipped; a Fable run skipped (its record at 15 min 9 s, rule 4) | reasoning | GPT-6 Sol; Grok 4.7; Claude Fable 5.1 on the Claude Code CLI | `xhigh`; `low`; `xhigh` | 1,174,711 (USD 5.61) for the Fable run | 20:04:16 to 20:20:58 |
| The fix of the defect that the skipped run reported, test first, and the freeze | execution | Claude Opus 5.5 | max | not reported | 20:21 to 20:30 |
| Review round 2, first and second harness: skipped (the same) | reasoning | GPT-6 Sol on the Devin CLI; Grok 4.7 on the OpenCode CLI | `xhigh`; `low` | not reported | 72 s, 20:29:39 to 20:30:51 |
| Review round 2 | reasoning | Claude Fable 5.1 on the Claude Code CLI | `xhigh` | 1,852,105 (USD 4.44) | 9 min 21 s, 20:30:58 to 20:40:19 |
| The close-out, with notes 1 to 3 of round 2 | execution | Claude Opus 5.5 | max | not reported | 20:40 to 20:43 |
