# The test runs of T-fsjp

The build ran from 2026-10-09T19:41Z (the claim of the branch) to 19:50Z (this
record), in the order below; the steps inside it were not stamped one by one.

## Red 1: the forge call before its code

`go vet -tags=integration ./internal/forge/github/` does not build the test,
for want of the method it names:

```
vet: internal/forge/github/github_integration_test.go:236:18: a.Comment undefined (type *Adapter has no field or method Comment)
```

The method then broke the two stand-ins of the interface (`internal/forge`'s
and `internal/run`'s tests), which each gained `Comment`.

## Red 2: the open attempt before its code

```
vet: internal/run/end_test.go:37:13: undefined: openAttempt
```

## The integration tests of the end

They were written after `end.go` (a deviation, as in row 36a), so their red is
each mutation below. Two cases were added before the mutations, as two rules
had none: an artifact that the head lacks, and a program that cannot start,
whose telemetry row needs an end. Row 36a's `TestTheStartRowBeforeTheProcess`
now expects the removal of the directory after the call.

## The mutations

A mutation of each rule, one at a time, run by a script that restores each file
from a copy after each; the unit rules run the unit tests, the others the
integration tests of the end. Each test line in full, cut at 400 characters
where a line is longer:

```
== open-closed: openAttempt: a closed or rebased of its attempt does not close it (exit 1)
end_test.go:38: a closed of its attempt after: true, want false
end_test.go:38: a rebased of its attempt after: true, want false
== open-other: openAttempt: an attempt of another does not close it (exit 1)
end_test.go:38: an attempt of another after: true, want false
== open-before: openAttempt: an event before the session's counts (exit 1)
end_test.go:38: a closed of another attempt after: false, want true
end_test.go:38: a closed of its attempt before: false, want true
end_test.go:38: no event session: true, want false
== open-none: openAttempt: no event session is open (exit 1)
end_test.go:38: no event session: true, want false
== end-open: endTask: no open-attempt check (exit 1)
end_integration_test.go:177: the events:
end_integration_test.go:180: a results/ file of a closed attempt
end_integration_test.go:183: the comment ["S-53ffaccd: developer session of T-ab12, attempt 1: done\n"]
== keep-done: endTask: the result file kept for every class (exit 1)
end_integration_test.go:127: a results/ file with no valid result file
== keep-closed: endTask: a closed attempt's result kept (exit 1)
end_integration_test.go:180: a results/ file of a closed attempt
== keep-branch: endTask: a head that cannot be read keeps the result (exit 1)
end_integration_test.go:147: a results/ file of a head that cannot be read
== artifact-refused: endTask: a differing artifact adds no refusal (exit 1)
end_integration_test.go:106: the events:
end_integration_test.go:112: the comment ["S-2d214685: developer session of T-ab12, attempt 1: done\n"]
end_integration_test.go:195: the events:
== artifact-sum: artifacts: the sum not compared (exit 1)
end_integration_test.go:106: the events:
end_integration_test.go:112: the comment ["S-89557cba: developer session of T-ab12, attempt 1: done\n"]
== artifact-missing: artifacts: a path that the head lacks passes (exit 1)
end_integration_test.go:195: the events:
== one-commit: endTask: the telemetry row in a commit of its own (exit 1)
end_integration_test.go:81: the last records commit changed ["tasks/T-ab12/events.tsv" "tasks/T-ab12/results/S-8ce823f1.tsv"]; want the result, the events and the telemetry row in one
== telemetry: endTask: no telemetry row (exit 1)
end_integration_test.go:81: the last records commit changed ["tasks/T-ab12/events.tsv" "tasks/T-ab12/results/S-951bc668.tsv"]; want the result, the events and the telemetry row in one
end_integration_test.go:83: git show layup-records:telemetry.tsv: exit status 128
end_integration_test.go:129: git show layup-records:telemetry.tsv: exit status 128
end_integration_test.go:212: git show layup-records:telemetry.tsv: exit status 128
== comment: endTask: no comment (exit 1)
end_integration_test.go:91: the comments on the control issue []
end_integration_test.go:112: the comment []
end_integration_test.go:150: the comment []
end_integration_test.go:183: the comment []
== comment-refused: endTask: the comment names no refusal (exit 1)
end_integration_test.go:112: the comment ["S-23b5753a: developer session of T-ab12, attempt 1: done\n"]
end_integration_test.go:150: the comment ["S-bb6af469: developer session of T-ab12, attempt 1: done\n"]
end_integration_test.go:183: the comment ["S-559a1a8f: developer session of T-ab12, attempt 1: done\n"]
== start-end: telemetry: a program that could not start has a zero end (exit 1)
end_integration_test.go:207: TaskSession: "S-6670eed1", the telemetry row of S-6670eed1: column end: the end is before the start
== remove-after: TaskSession: the directory kept after the call (exit 1)
end_integration_test.go:94: the session directory after the call: <nil>; want it removed
```

## Green (2026-10-09T19:51Z)

Each with exit 0: `go build ./...`, `go vet ./...`, `gofmt -l internal cmd`
(empty), `go test ./...`, `go test -tags=integration ./...` (with
`TestPackageRules` and `TestInputRule`); `adr-lint`, `prd-lint`,
`link-lint`, `setup-check`, `run-discipline-tests`, `git diff --check`;
`sh runs/T-fdaq/parts.sh` (each part in one row, after the move of the fetch).
