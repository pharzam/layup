# T-ysph — row 30a of the plan, the host registers of M2b

Issue: [#155](https://github.com/pharzam/layup/issues/155), row 30a of [the tasks
of M2b](../plan/README.md#the-tasks-of-m2b), opened by `T-fdaq` (#152). Serves
`F-0003#52` through `REQ-013` and `NFR-001`. Base `ad8774d` (the merge of #173;
the plan was reviewed at `8d2fa2e`). Author: Claude Opus 5.5 on Claude Code.
Evidence: [`runs/T-ysph/`](../../runs/T-ysph/).

## Plan and plan review

The plan (R12, comment 6080340002), its review (6080580136) and the author's
answer (6080580748) are comments on #155. The plan review (Claude Fable 5.1,
effort `xhigh`, on the Claude Code CLI with stream output, a fresh read-only
session in a clone at `8d2fa2e`, 7 min 45 s) gave `approve-with-conditions`:
Budget maximum 1,100 lines added plus removed over 20 files against `8d2fa2e`,
close-out inside; Cycle cap 1; no panel. Its four conditions are applied: the
columns of `harness-register` that never hold `—`, one test per column; one
reading of a `PATH` with `..` and the form of `NAME`; a home in `records.md` for
each rule beyond the types; `packages.md` and the evidence file. The goal count
is 2, final by O-187 of #152.

## What was done

1. **D1:** the ten columns of `M2b` in the block `harness-register` and in
   `HarnessRegisterSchema`, after `wall`; `ReadHarnesses` checks their rules
   (the columns that never hold `—`; `{cap}` exactly when `cap`; one space
   between words; `credential`, `credential_to`, `policy`, `vars`); the three
   fixtures that write the register.
2. **D2:** `ModelsSchema` and `ReadModels`, `RoutingRegisterSchema` and
   `ReadRoutingRegister` (the positions of each role and tier 1 to k, in any
   order, as `routing`), `CheckRegisters` across the three files; `models` and
   `routing-register` moved to `built`.
3. **D3:** `CheckCredential`: an absolute path, a regular file of mode 0600,
   owned by the user of the run.
4. **D4:** `layup run` reads `models.tsv` and `routing.tsv`, checks the three
   registers across their files and each credential, and exits 2 naming the
   file. `Config` keeps no field for them yet: row 30b or 38 adds it.
5. **D5:** `records.md` (the block, and the rules of the three host registers
   of `M2b`), `session.md` (the form of `NAME`, the prohibitions of
   `var:NAME`, one reading of a `PATH` with `..`), `run.md` (the four
   registers), `packages.md` (the Job cells of `internal/cli` and
   `internal/route`), `traceability.md`, the Test cell of `REQ-013` and a §13
   line.

**Tests:** [`test-runs.md`](../../runs/T-ysph/test-runs.md): red before the
readers (no compile); red for `layup run` before it read the registers of
`M2b`; green after.

**Decided here** (note 6 of the plan review): the `NAME` of `var:NAME` takes
the prohibitions of a `vars` name, as a `var:PATH` would take the place of the
host's `PATH`.

**The rejected alternatives:** one check of a file's mode and owner shared by
`internal/forge` and `internal/route` (a new package or an import that neither
row allows); the credential check at the start of a session (`session.md` gives
it to the input states of `layup run`); a closed list of roles.
