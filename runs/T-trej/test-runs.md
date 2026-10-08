# T-trej — the test runs

Evidence for row 25a (#130, O-173 a): each test red first, for the right reason, then green.

## Red 1 (2026-10-08T10:02Z)

`lease_test.go` (a stand-in clock and a stand-in records store) and `rules_test.go` before any code: `go test ./internal/run/` did not compile (`undefined: LeaseRow`).

## Red 2 (2026-10-08T10:04Z)

With stubs that do nothing: each case failed for its rule: `Take = {…}` (not taken), `taken over after 0s, want 3 × H`, `the messages []; want one that names the old run`, `progress lines []`, `Take with a moving counter = <nil>; want a HeldError`, `a refused takeover = <nil>; want a LostError`, `Take with a cancelled context = <nil>`, `Beat = <nil>; … want 4`, `heartbeat 0 with 0 commits after 3 × H`, `Decision = false, want true` (four cases), `the roles: Intake [], acceptance []`, `the first copy = [], "", false`. The cases whose answer is "no" passed on the stub, as a stub must.

## Green, one fix of a stand-in, two mutations (2026-10-08T10:08Z)

The real code: one failure, `heartbeat 4 with 4 commits after 3 × H; want 3`: the stand-in clock checked the context only before it moved, so a sleep in which the context ended did not end. A real sleep ends with its context (`Clock`: "Sleep ends early with the context"); the stand-in now gives `ctx.Err()` after it moves. Then `go test ./internal/run/` passes, and `TestPackageRules` (`go test -tags=integration ./cmd/layup/`) passes with `internal/run`.

Mutations on backup copies, each put back (`cmp` equal):
- the "held" check before the takeover check → `the lease is held by the run aaaa…` in the demo test, and the refused-takeover case;
- `Decision` with no check of `seen` → `an edit (seen 2): Decision = true, want false`.

## The documents (2026-10-08T10:12Z)

Red before the edits: `sh runs/T-trej/docs.sh` gave sixteen `FAIL` lines, exit 1. After the edits, two rules of the check were wrong, not the documents: the rule of row 25b asked for an After cell of `25a` alone (the row's After is `22b, 23, 25a`, as 25b uses `internal/git` and the registers), and the rule of the pitfall asked for a phrase that the text wraps over two lines. Each now tests what the row and the pitfall hold. Then sixteen `ok`, exit 0.
