# T-mqty — the test runs

Evidence for row 26 (#131): each test red first, for the right reason, then green.

## `internal/run` and `internal/records`: red, then green (2026-10-08T11:58Z)

`TestStartAndRestartCallStepAtTheStartOfEachStep`, `TestCheckGit`, `TestCheckValues` (`internal/run`) and `TestCheckValueGivesTheFormOfAName` (`internal/records`) before the code: neither package compiled (`cfg.Step undefined`, `undefined: CheckValue`). Then `Config.Step` (called at the start of each step, the defect of round 2 of #142), `CheckGit`, `CheckValues` and `records.CheckValue`: both packages pass.

## The adapter's own client: red, then green (2026-10-08T12:02Z)

`TestTheOwnClientFollowsNoRedirectAndNoProxy` (a server that answers `302` to another server; the transport's `Proxy`) on the adapter of `075624f`: `panic: runtime error: invalid memory address or nil pointer dereference` (no client). Then `New` makes a client with no redirect and a transport with no proxy: `go test -tags=integration ./internal/forge/github/` passes. A proxy of the environment is not shown by a request, as Go's proxy rule never takes a proxy for a loopback address; the test reads the transport.
