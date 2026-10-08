# T-mqty — row 26 of the plan, the command layup run

Issue: [#131](https://github.com/pharzam/layup/issues/131), row 26 of [the tasks
of M2a](../plan/README.md#the-tasks-of-m2a), opened by `T-zwke` (#124). Serves
`F-0003#42` through `NFR-001` and `NFR-002`. Base `075624f` (the merge of
#144). Author: Claude Opus 5.5 on Claude Code. Started under O-169. Evidence:
[`runs/T-mqty/`](../../runs/T-mqty/).

## Plan and plan review

The plan (R12, comment 6059146521) counted one goal before the plan (the lesson
of O-173): the surface of one command. Its review (Claude Fable 5.1, effort
`xhigh`, on the Claude Code CLI with stream output, a fresh read-only session
in a clone at `075624f`, 11 min 46 s; comment 6059346712):
`approve-with-conditions`; it counted four classes and asked no split, on the
wrong reading that O-173 counted this row (the author's answer corrected it:
the count is the author's, open to the Operator); six conditions (the `git`
check in `internal/run`, not `internal/cli`; the repeat of Start in a new
world; no proxy in the adapter's client; one line of a wait; a brief checked
for UTF-8 only; an absolute `--host`) and twelve notes. The author's answer
(6059347254) took all six and notes 2 to 9. Budget maximum 1,500 lines added
plus removed over 24 files against `075624f`, close-out inside; Cycle cap 1; no
panel.

## What was done

1. **`internal/cli`:** the command table takes a selecting flag and flags that
   may be left out; the two rows of `layup run` (Start by `--new`, the restart
   by `TARGET`); each input checked before the first step (exit 2, no table):
   `OWNER/NAME`, the briefs (UTF-8), the values of the flags, the two host
   registers, the key file, `git`; `--host` made absolute; the adapter built
   with its own client and the run's progress; a run ID, the host name, a real
   clock, Ctrl-C; the table on standard output, exit 0 or 1; LAYUP's pin as two
   variables held to `armature.pin` by a test. The progress prints a line of a
   wait with its step, and the next beat prints nothing.
2. **`internal/run`:** `Config.Step` at the start of each step (a defect
   against `README.md`, found by round 2 of #142); `CheckGit`; `CheckValues`
   with `records.CheckValue` (one home of the value forms).
3. **`internal/forge/github`:** `New` makes a client with no redirect and a
   transport with no proxy when its caller gives none (rule 5).
4. **The e2e test:** a binary built with the pin set to a local bare baseline,
   a fake forge on loopback, local bare targets, a stand-in Operator; Start in
   two worlds with the same arguments, the restart twice, a usage error and an
   input error.
5. **The documents:** `README.md` (the two table rules, the repeat of a command
   that changes its forge, exit 2 of `git`), `run.md` (the line of a wait, two
   rows of the input states), `forge.md` (the adapter's own client, a known
   limit of the certificate roots), `packages.md` (the May import cell of
   `internal/cli`, the Job of `internal/run`), two glossary rows, the
   traceability, the Test cells of `NFR-001`, `NFR-002` and a §13 line.
   [`docs.sh`](../../runs/T-mqty/docs.sh) checks them.

LAYUP's pin has two homes, both tested against `armature.pin`: the new
variables of `internal/cli` and `layupPin` of `internal/setup` (note 5 of the
plan review).

**Tests:** [`test-runs.md`](../../runs/T-mqty/test-runs.md): red, then green,
for each part; three mutations of `internal/cli`, each caught; the e2e red on
the commit before the command; `docs.sh` red, then green.

**The rejected alternatives:** an HTTP client built in `internal/cli` (rule 5);
the pin read from `armature.pin` at run time (the binary runs outside LAYUP's
checkout); a sub-word (`layup run new`) for Start.
