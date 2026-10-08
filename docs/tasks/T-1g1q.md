# T-1g1q — row 22b of the plan, the two host registers and the key file

Issue: [#137](https://github.com/pharzam/layup/issues/137), row 22b of [the tasks
of M2a](../plan/README.md#the-tasks-of-m2a), opened by `T-esfe` (#127) after
O-170 (b) split row 22. Serves `F-0003#42` through `NFR-001` and `NFR-007`. Base
`10a5691` (the merge of #138). Author: Claude Opus 5.5 on Claude Code. Evidence:
[`runs/T-1g1q/`](../../runs/T-1g1q/).

## Plan and plan review

The plan (R12, comment 6054078662) listed three goal classes. Its review (Claude
Fable 5.1, effort `xhigh`, on the Claude Code CLI with stream output, a fresh
read-only session in a clone at `10a5691`, 12 min 24 s; comment 6054282765) gave
`reject` on the goal count alone (three), with six more conditions; the author's
answer (6054283098) applied them and asked **O-171**. **O-171** (b) (comment
6054413706; the text was written by a Devin session and pasted and adopted by the
Operator, comment 6054453470): the two host-register readers and the key-file
validation are one goal; the forge interface (class 3) moves to row 24, its
delivery scope only, its code staying in `internal/forge`, with conditions 4 to 6
and note 5 of the review carried to #129, open. Budget maximum 1,050 lines added
plus removed over 17 files against `10a5691`, close-out inside; Cycle cap 1; no
panel.

## What was done

1. **`internal/forge`** (new): `ForgeRegisterSchema` and `ReadForgeRegister` (one
   row; each column but `watch_slug` holds a value; an absolute `key_file`; each
   error names its line), and `CheckKeyFile`, which reads the mode, the owner and
   the bytes and calls the pure `checkKey`: mode exactly 0600, the owner the user
   of the run, one PEM block (`RSA PRIVATE KEY`, PKCS #1 of version 0, or
   `PRIVATE KEY`, PKCS #8 of the RSA algorithm), an RSA key that `Validate`
   accepts; it gives the key, so the adapter of row 24 parses nothing. It uses
   `encoding/pem`, `encoding/asn1` and `crypto/rsa`, not `crypto/x509` (rule 5).
2. **`internal/route`** (new, its `M2a` job only): `HarnessRegisterSchema` and
   `ReadHarnesses`, which gives the rows and the IDs in their order; `wall` is 1
   or more; `cap` may be empty; an empty register is allowed.
3. **Both owners** compare their schemas with their blocks; the two names moved
   to `built`; `TestPackageRules` passes with the two packages.
4. **The documents:** `forge.md` says "one PEM block" and its two forms; the Job
   cell of `internal/forge`; the plan rows 22b and 24 (class 3 to row 24 by
   O-171); the traceability rows; the Test cells of `NFR-001` and `NFR-007` and a
   §13 line. [`docs.sh`](../../runs/T-1g1q/docs.sh) checks them.

**Tests:** [`test-runs.md`](../../runs/T-1g1q/test-runs.md): red 1 (no package),
red 2 (a stub; each rule case failed for its reason), green; one comparison shown
to fail with a column changed on a backup copy; `docs.sh` red, then green. The
tests make one 2,048-bit key per test binary at run time: no key file enters the
tree.

**The rejected alternatives:** the key parse with `crypto/x509` (it depends on
`net`); the harness register in `internal/forge`; the interface in this row
(O-171).
