# The test runs of T-w89c

The host: the LAYUP host of 2026-10-09, Linux on amd64, `go version go1.26.9
linux/amd64`. The red downloads `go1.26.8`; each run fetches govulncheck and
its vulnerability database, so it needs the network.

## Red, at the base `ad3c400`

`GOTOOLCHAIN=go1.26.8 go run golang.org/x/vuln/cmd/govulncheck@v1.8.0 ./...`:
govulncheck exits 3 (`go run` reports `exit status 3` and exits 1), with the
ten advisories of #150. The block is cut: its first line lists the IDs of
the ten entries that govulncheck prints, each with its "Found in: …@go1.26.8"
and "Fixed in: …@go1.26.9" lines, and its last lines are govulncheck's own:

```
GO-2026-6603 GO-2026-6604 GO-2026-6605 GO-2026-6607 GO-2026-6608 GO-2026-6610 GO-2026-6611 GO-2026-6612 GO-2026-6613 GO-2026-6617
Your code is affected by 10 vulnerabilities from the Go standard library.
This scan also found 1 vulnerability in packages you import and 2
vulnerabilities in modules you require, but your code doesn't appear to call
these vulnerabilities.
Use '-show verbose' for more details.
exit status 3
```

## Green, with `go 1.26.9` in `go.mod`

`GOTOOLCHAIN=go1.26.9 go run golang.org/x/vuln/cmd/govulncheck@v1.8.0 ./...`:
exit 0, `No vulnerabilities found.` With `GOTOOLCHAIN=go1.26.9`: `go build ./...`,
`go vet ./...`, `go test ./...`, `go test -tags=integration ./...` and
`go test -tags=e2e ./...` pass.

## Green in CI

Run 37912054398 of #151 at `b595c34`: the job `security` passes; its setup-go
step logs `Setup go version spec 1.26.9`, then `go version go1.26.9
linux/amd64` and `GOVERSION='go1.26.9'`. Each other job passes except
`review-record`, which waits for the review round.
