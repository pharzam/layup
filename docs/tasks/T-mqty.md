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
panel. **O-175 (a)** (the Operator's answer in the session to comment
6059557984, recorded in 6059638358): 1,500 lines over 28 files.

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

## Review rounds

Round 1 (Claude Fable 5.1, effort `xhigh`, on the Claude Code CLI with stream
output, a fresh read-only session in a clone at `09a0377`, 12 min 14 s; comment
6059868590): `material`. Finding 1: a wait printed two lines every ten seconds
(the old `note` started a new beat just before the wait's own ten seconds
ended). Fixed in `eb913f5` (comment 6059995056): a line of a wait takes the
place of the next beat, test first; notes 3, 4, 5, 6, 7 and 9 applied.

Round 2 (the same set-up, at `eb913f5`, 9 min 9 s; comment 6060160752):
`material`, the last round under cap 1: two texts left stale by that fix (the
sentence of `run.md` on the beat; a test name in the traceability). **O-176
(a)** (the Operator's answer in the session to comment 6060161092, recorded in
6060243206): one more cycle, cap 2. Fixed in `7363056` (comment 6060275922),
with a rule of `docs.sh` for each, red on `eb913f5`.

Round 3 (the same set-up, at `7363056`, 7 min 52 s; comment 6060431989): `nothing material in scope`, four notes, not applied to the
reviewed head.

**Known limits:**

- A line of a wait repeats its step name (`root-push: root-push: waiting …`;
  round 1, note 2).
- The wait line and the beat drift apart over many cycles, so one extra beat
  prints now and then (round 3, note 1).
- A rate-limit line of the heartbeat goroutine goes through the same progress
  as the main goroutine, a data race on a rare path (round 2, note 6).
- A `run-steps` table that cannot be written gives exit 2, which `run.md` does
  not name (round 2, note 7).
- `stepNames` keys on the count of steps; a later milestone that adds a step to
  Start changes it (round 3, note 3).
- LAYUP's pin has two homes, each tested against `armature.pin` (note 5 of the
  plan review).

## Verdict

Delivered: `layup run --new` then `layup run TARGET`, as the built binary
against a local fake forge, give their tables, the same bytes on a repeat. The
review ended by decay at cycle 2 (cap 2 by O-176 a). The diff against
`075624f` is inside 1,500 lines over 28 files (O-175 a).

Next: row 27 (`T-fnsr`, #132), the uat of `M2a`, which needs the Operator's
inputs of #132.

## Resource record

Recorded, not budgeted (ADR-0007). Times are UTC on 2026-10-08; tokens are the
`result` event of the Claude Code CLI (stream runs); `not reported` otherwise.

| Part | Expected tier | Model | Effort | Tokens | Elapsed |
| ---- | ------------- | ----- | ------ | ------ | ------- |
| The reading; the goal count; the plan | reasoning | Claude Opus 5.5 | max | not reported | 11:30 to 11:45 |
| Its plan review | reasoning | Claude Fable 5.1, Claude Code CLI | `xhigh` | 1,075,363 (USD 7.62) | 11 min 46 s, from 11:45 |
| The answer; the work, test first; the documents | execution | Claude Opus 5.5, a reasoning-tier model on an execution part | max | not reported | 11:57 to 12:13 |
| O-175 | reasoning | Claude Opus 5.5 | max | not reported | 12:13 to 12:15 |
| Round 1 | reasoning | Claude Fable 5.1, Claude Code CLI | `xhigh` | 1,096,213 (USD 5.96) | 12 min 14 s, from 12:15 |
| The fix of round 1 | execution | Claude Opus 5.5 | max | not reported | 12:27 to 12:35 |
| Round 2 | reasoning | Claude Fable 5.1, Claude Code CLI | `xhigh` | 1,242,119 (USD 5.48) | 9 min 9 s, from 12:35 |
| O-176; the fix of round 2 | execution | Claude Opus 5.5 | max | not reported | 12:44 to 12:52 |
| Round 3 | reasoning | Claude Fable 5.1, Claude Code CLI | `xhigh` | 1,178,703 (USD 4.56) | 7 min 52 s, from 12:52 |
| The close-out | execution | Claude Opus 5.5 | max | not reported | 13:00 to 13:04 |
