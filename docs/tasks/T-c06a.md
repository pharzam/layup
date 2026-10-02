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
`bash`, the trap on `EXIT` could remove the temporary file before `tr` read
it (4 runs of 20 gave `clear`); the script removes the file in `result()`, and
a lesson is in `guardrails.md` §2.
